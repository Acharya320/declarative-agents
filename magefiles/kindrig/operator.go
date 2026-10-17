// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package kindrig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// The agent operator is a platform boot component (applications srd009
// R4.2): the platform builds its image from the checkout's operator module,
// loads it into the node, and installs the module's chart, so a booted
// platform always carries a ready operator.
const (
	operatorComponent = "agent-operator"
	operatorNamespace = "agent-operator"
	operatorRelease   = "agent-operator"
	operatorModule    = "operator"
	operatorChart     = "operator/deploy/chart/agent-operator"
	operatorTimeout   = "180s"
)

// operatorDockerfile is the runtime layer the platform builds around the
// checkout's binary; operator/Dockerfile builds the same layer for
// publication, from the module's source.
const operatorDockerfile = "FROM docker.io/library/alpine:3.22@sha256:5291449c3df73caf6ed85e649dec1b9e818b39a5d8c871e97afc13e9cd5e8fa8\n" +
	"RUN apk add --no-cache ca-certificates\n" +
	"COPY agent-operator /usr/local/bin/agent-operator\n" +
	"USER 65532:65532\n" +
	"ENTRYPOINT [\"agent-operator\"]\n"

// InstallAgentOperator builds the operator image from the checkout when the
// node lacks it, loads it, and installs or upgrades the chart, waiting for
// the rollout and logging one phase per step like the other components.
// The returned cleanup uninstalls the release.
func InstallAgentOperator(run CommandRunner, cluster string) (func() error, error) {
	if strings.TrimSpace(cluster) == "" {
		return nil, fmt.Errorf("install %s: kind cluster name is required", operatorComponent)
	}
	root, err := checkoutRoot(run)
	if err != nil {
		return nil, err
	}
	image, err := OperatorRuntimeImage(run, root)
	if err != nil {
		return nil, err
	}
	if err := buildOperatorImage(run, cluster, root, image); err != nil {
		return nil, err
	}
	repository, tag, _ := strings.Cut(image, ":")
	steps := []installStep{
		{"image-load", []string{"kind", "load", "docker-image", image, "--name", cluster}},
		{"chart-install", []string{"helm", "upgrade", "--install", operatorRelease, filepath.Join(root, operatorChart),
			"--namespace", operatorNamespace, "--create-namespace", "--kube-context", KubeContext(cluster),
			"--set", "image.repository=" + repository, "--set", "image.tag=" + tag,
			"--wait", "--timeout", operatorTimeout}},
		{"rollout", []string{"kubectl", "rollout", "status", "deployment/" + operatorRelease,
			"--namespace", operatorNamespace, "--timeout=" + operatorTimeout}},
	}
	if err := runInstallSteps(run, cluster, operatorComponent, steps); err != nil {
		return nil, err
	}
	return func() error {
		_, err := run("helm", "uninstall", operatorRelease, "--namespace", operatorNamespace,
			"--kube-context", KubeContext(cluster), "--ignore-not-found")
		return err
	}, nil
}

// checkoutRoot is the repository root the operator module builds from.
func checkoutRoot(run CommandRunner) (string, error) {
	out, err := run("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("resolve checkout root: %w", err)
	}
	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", fmt.Errorf("resolve checkout root: git printed nothing")
	}
	return root, nil
}

// OperatorRuntimeImage is the immutable local name of the operator built
// from the checkout revision on the host platform (srd009 R4.1), under the
// first-party naming chart conformance enforces.
func OperatorRuntimeImage(run CommandRunner, root string) (string, error) {
	out, err := run("git", "-C", root, "rev-parse", "HEAD")
	if err != nil {
		return "", fmt.Errorf("resolve checkout revision: %w", err)
	}
	revision := strings.ToLower(strings.TrimSpace(string(out)))
	if !gitRevision.MatchString(revision) {
		return "", fmt.Errorf("checkout revision %q is not a commit hash", revision)
	}
	return FormatGitLocal(RuntimeRole, operatorComponent, revision, HostPlatform())
}

// buildOperatorImage compiles the operator for the node's platform into a
// staging directory and builds the runtime image there, unless the host
// already holds that image.
func buildOperatorImage(run CommandRunner, cluster, root, image string) error {
	started := time.Now()
	if dockerImagePresent(run, image) {
		LogPhase(cluster, "image-build", "skipped", started, "image="+image)
		return nil
	}
	ctxDir, err := os.MkdirTemp("", "agent-operator-image-")
	if err != nil {
		return fmt.Errorf("stage operator image: %w", err)
	}
	defer os.RemoveAll(ctxDir)
	if err := os.WriteFile(filepath.Join(ctxDir, "Dockerfile"), []byte(operatorDockerfile), 0o644); err != nil {
		return fmt.Errorf("write operator image Dockerfile: %w", err)
	}
	_, arch, err := SplitPlatform(HostPlatform())
	if err != nil {
		return err
	}
	steps := []installStep{
		{"binary-build", []string{"env", "CGO_ENABLED=0", "GOOS=linux", "GOARCH=" + arch,
			"go", "build", "-C", filepath.Join(root, operatorModule), "-trimpath", "-ldflags=-s -w",
			"-o", filepath.Join(ctxDir, "agent-operator"), "./cmd/agent-operator"}},
		{"image-build", []string{"docker", "build", "--platform", HostPlatform(), "-t", image, ctxDir}},
	}
	return runInstallSteps(run, cluster, operatorComponent, steps)
}
