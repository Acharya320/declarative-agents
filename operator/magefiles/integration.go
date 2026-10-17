// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"bufio"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Nokia-Bell-Labs/declarative-agents/magefiles/kindrig"
	"github.com/magefile/mage/mg"
)

// Integration holds the kind proofs (applications srd009 R5).
type Integration mg.Namespace

// The proof's namespace, the collector profile it mounts, and its bounds.
const (
	proofNamespace   = "agent-operator-proof"
	proofConfigMap   = "collector-profile"
	proofAgent       = "collector"
	collectorProfile = "applications/catalog/agents/collector"
	// The shared units file the collector profile imports relative to its directory.
	collectorUnitsFile = "types-collector.yaml"
	proofReadyWait     = "180s"
	proofHTTPTimeout   = 60 * time.Second
	proofGoneTimeout   = 90 * time.Second
	proofCommandWait   = 3 * time.Minute
)

//go:embed testdata/collector-agent.yaml
var collectorAgentManifest string

var forwardedPort = regexp.MustCompile(`Forwarding from 127\.0\.0\.1:(\d+) ->`)

// All runs the collector proof on the kind platform (srd009 R5.1, R5.2):
// the profile directory becomes a ConfigMap, an Agent is applied, the
// operator makes it Ready, the health endpoint answers through the Service,
// and deleting the Agent removes the Deployment and the Service.
func (Integration) All() error {
	root, err := repositoryRoot()
	if err != nil {
		return err
	}
	platform, err := kindrig.AcquirePlatform(kindrig.PlatformOptions{
		EvidenceDirectory: filepath.Join(root, "build", "kind-evidence",
			kindrig.PlatformClusterName+"-"+time.Now().UTC().Format("20060102T150405Z")),
	})
	if err != nil {
		return fmt.Errorf("acquire platform: %w", err)
	}
	failed := true
	defer func() { platform.Stop(failed) }()
	if err := collectorProof(root, platform); err != nil {
		return err
	}
	failed = false
	fmt.Println("integration:all PASS - the collector Agent became Ready, answered its health endpoint, and was removed with its Deployment and Service")
	return nil
}

func collectorProof(root string, platform *kindrig.Platform) error {
	cluster := platform.Cluster.Name
	coreRoot := filepath.Join(root, "agent-core")
	image, err := kindrig.AgentCoreRuntimeImage(coreRoot)
	if err != nil {
		return err
	}
	if _, err := kindrig.EnsureAgentCoreImage(coreRoot, image); err != nil {
		return fmt.Errorf("agent-core image: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), proofCommandWait)
	defer cancel()
	if err := kindrig.LoadImage(ctx, kindrig.DefaultRunContext, cluster, image); err != nil {
		return err
	}
	commands, unbind, err := kindrig.ClusterCommands(kindrig.CaptureRun, cluster)
	if err != nil {
		return err
	}
	defer unbind()
	kubectl := func(args ...string) ([]byte, error) {
		ctx, cancel := context.WithTimeout(context.Background(), proofCommandWait)
		defer cancel()
		return commands.CommandContext(ctx, "kubectl", args...).CombinedOutput()
	}
	step := func(name string, args ...string) error {
		started := time.Now()
		out, err := kubectl(args...)
		outcome := "passed"
		if err != nil {
			outcome = "failed"
		}
		kindrig.LogPhase(cluster, name, outcome, started, "component=agent-operator-proof")
		if err != nil {
			return fmt.Errorf("%s: kubectl %s: %w: %s", name, strings.Join(args, " "), err, strings.TrimSpace(string(out)))
		}
		return nil
	}
	_, _ = kubectl("delete", "namespace", proofNamespace, "--ignore-not-found", "--wait=true")
	if err := step("namespace", "create", "namespace", proofNamespace); err != nil {
		return err
	}
	defer func() { _, _ = kubectl("delete", "namespace", proofNamespace, "--ignore-not-found", "--wait=false") }()
	// The profile's top-level files and the units file it imports make the
	// ConfigMap, keyed so the operator lays them out as directories. Its
	// ui/dist is served per request and absent here, so the UI route answers
	// 404 while the control server answers health, which is what the proof
	// reads.
	files, err := profileFiles(filepath.Join(root, collectorProfile))
	if err != nil {
		return err
	}
	if err := step("profile-configmap", append([]string{"create", "configmap", proofConfigMap, "--namespace", proofNamespace}, files...)...); err != nil {
		return err
	}
	manifest := strings.ReplaceAll(collectorAgentManifest, "AGENT_CORE_IMAGE", image)
	path, err := stageManifest(manifest)
	if err != nil {
		return err
	}
	defer os.Remove(path)
	if err := step("agent-apply", "apply", "--namespace", proofNamespace, "-f", path); err != nil {
		return err
	}
	if err := step("agent-ready", "wait", "--namespace", proofNamespace, "agent/"+proofAgent,
		"--for=condition=Ready", "--timeout="+proofReadyWait); err != nil {
		_ = captureProofEvidence(platform, cluster)
		return err
	}
	if err := healthThroughService(commands, cluster); err != nil {
		_ = captureProofEvidence(platform, cluster)
		return err
	}
	if err := step("agent-delete", "delete", "--namespace", proofNamespace, "agent/"+proofAgent, "--wait=true"); err != nil {
		return err
	}
	return ownedObjectsGone(kubectl, cluster)
}

// healthThroughService forwards the Service's control port and reads the
// health endpoint, the readiness signal the probe uses (srd009 R1.2).
func healthThroughService(commands kindrig.Commands, cluster string) error {
	started := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), proofHTTPTimeout+10*time.Second)
	defer cancel()
	forward := commands.CommandContext(ctx, "kubectl", "port-forward", "--namespace", proofNamespace, "svc/"+proofAgent, "0:18191")
	stdout, err := forward.StdoutPipe()
	if err != nil {
		return err
	}
	forward.Stderr = os.Stderr
	if err := forward.Start(); err != nil {
		return fmt.Errorf("port-forward: %w", err)
	}
	defer func() { _ = forward.Process.Kill(); _ = forward.Wait() }()
	scanner := bufio.NewScanner(stdout)
	port := ""
	for scanner.Scan() {
		if m := forwardedPort.FindStringSubmatch(scanner.Text()); m != nil {
			port = m[1]
			break
		}
	}
	if port == "" {
		return errors.New("port-forward printed no forwarded port")
	}
	url := "http://127.0.0.1:" + port + "/api/lifecycle/health"
	deadline := time.Now().Add(proofHTTPTimeout)
	var last error
	for time.Now().Before(deadline) {
		resp, err := http.Get(url) //nolint:gosec // the proof's own loopback forward
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				kindrig.LogPhase(cluster, "health", "passed", started, "url="+url)
				return nil
			}
			last = fmt.Errorf("GET %s: status %d", url, resp.StatusCode)
		} else {
			last = err
		}
		time.Sleep(time.Second)
	}
	kindrig.LogPhase(cluster, "health", "failed", started, "url="+url)
	return fmt.Errorf("no HTTP 200 from the collector's health endpoint within %s: %w", proofHTTPTimeout, last)
}

// ownedObjectsGone waits for the garbage collector to remove what the Agent
// owned (srd009 R2.3).
func ownedObjectsGone(kubectl func(...string) ([]byte, error), cluster string) error {
	started := time.Now()
	deadline := time.Now().Add(proofGoneTimeout)
	for time.Now().Before(deadline) {
		gone := true
		for _, kind := range []string{"deployment", "service"} {
			out, err := kubectl("get", kind, proofAgent, "--namespace", proofNamespace)
			if err == nil || !strings.Contains(string(out), "NotFound") {
				gone = false
			}
		}
		if gone {
			kindrig.LogPhase(cluster, "owned-objects-gone", "passed", started, "component=agent-operator-proof")
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	kindrig.LogPhase(cluster, "owned-objects-gone", "failed", started, "component=agent-operator-proof")
	return fmt.Errorf("the Deployment or Service outlived the Agent by %s", proofGoneTimeout)
}

func captureProofEvidence(platform *kindrig.Platform, cluster string) error {
	evidence := kindrig.FailureEvidence{Directory: filepath.Join("build", "kind-evidence", "proof-"+time.Now().UTC().Format("20060102T150405Z")),
		Namespaces: []string{proofNamespace}, Run: platform.Run}
	return evidence.Capture(kindrig.DefaultRun, cluster)
}

// profileFiles lists the profile directory's regular files, plus the shared
// units file the profile imports as ../units/types-collector.yaml, as
// --from-file arguments keyed <dir>__<file>: the operator reads "__" as a
// directory separator (srd009 R2.1), so the mount reproduces the catalog's
// agents/<name> and agents/units layout. Subdirectories (tests, ui) are not
// part of the profile contract.
func profileFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read collector profile: %w", err)
	}
	name := filepath.Base(dir)
	var args []string
	for _, entry := range entries {
		if entry.Type().IsRegular() {
			args = append(args, fmt.Sprintf("--from-file=%s__%s=%s", name, entry.Name(), filepath.Join(dir, entry.Name())))
		}
	}
	if len(args) == 0 {
		return nil, fmt.Errorf("collector profile %s holds no files", dir)
	}
	units := filepath.Join(filepath.Dir(dir), "units", collectorUnitsFile)
	if _, err := os.Stat(units); err != nil {
		return nil, fmt.Errorf("collector profile imports %s: %w", units, err)
	}
	return append(args, fmt.Sprintf("--from-file=units__%s=%s", collectorUnitsFile, units)), nil
}

func stageManifest(manifest string) (string, error) {
	file, err := os.CreateTemp("", "collector-agent-*.yaml")
	if err != nil {
		return "", err
	}
	if _, err := file.WriteString(manifest); err != nil {
		_ = file.Close()
		return "", err
	}
	return file.Name(), file.Close()
}

func repositoryRoot() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("repository root: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
