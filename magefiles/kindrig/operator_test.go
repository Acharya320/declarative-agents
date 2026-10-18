// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package kindrig

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// operatorRunner scripts the commands InstallAgentOperator issues: the
// checkout root and revision, and whether the image already exists.
func operatorRunner(present bool) (CommandRunner, *[]string) {
	var calls []string
	return func(name string, args ...string) ([]byte, error) {
		command := strings.Join(append([]string{name}, args...), " ")
		calls = append(calls, command)
		switch {
		case command == "git rev-parse --show-toplevel":
			return []byte("/repo\n"), nil
		case strings.HasSuffix(command, "rev-parse HEAD"):
			return []byte("0123456789abcdef0123456789abcdef01234567\n"), nil
		case name == "docker" && len(args) > 1 && args[1] == "inspect" && !present:
			return []byte("No such image"), errors.New("absent")
		}
		return nil, nil
	}, &calls
}

// TestInstallAgentOperatorBuildsLoadsAndInstalls covers srd009 R4.1 and
// R4.2: the binary builds for the node platform, the image carries the
// first-party git-revision name, kind loads it, helm installs the module's
// chart with that image, and the rollout is awaited.
func TestInstallAgentOperatorBuildsLoadsAndInstalls(t *testing.T) {
	run, calls := operatorRunner(false)
	cleanup, err := InstallAgentOperator(run, "da-example")
	if err != nil {
		t.Fatal(err)
	}
	image := "localhost/declarative-agents/runtime/agent-operator:git-0123456789ab-linux-" + runtime.GOARCH
	for _, want := range []string{
		"env CGO_ENABLED=0 GOOS=linux GOARCH=" + runtime.GOARCH + " go build -C /repo/operator -trimpath -ldflags=-s -w -o ",
		"cp -R /repo/agent-core/tools ",
		"docker build --platform linux/" + runtime.GOARCH + " -t " + image + " ",
		"kind load docker-image " + image + " --name da-example",
		"helm upgrade --install agent-operator /repo/operator/deploy/chart/agent-operator --namespace agent-operator --create-namespace --kube-context kind-da-example --set image.repository=localhost/declarative-agents/runtime/agent-operator --set image.tag=git-0123456789ab-linux-" + runtime.GOARCH + " --wait --timeout 180s",
		"kubectl --context kind-da-example rollout status deployment/agent-operator --namespace agent-operator --timeout=180s",
	} {
		if !containsPrefix(*calls, want) {
			t.Errorf("no command starting %q in\n%s", want, strings.Join(*calls, "\n"))
		}
	}
	if err := cleanup(); err != nil {
		t.Fatal(err)
	}
	if last := (*calls)[len(*calls)-1]; !strings.HasPrefix(last, "helm uninstall agent-operator --namespace agent-operator --kube-context kind-da-example") {
		t.Errorf("cleanup ran %q", last)
	}
}

// TestInstallAgentOperatorReusesAPresentImage skips the build when the host
// already holds the revision's image, and still installs the chart.
func TestInstallAgentOperatorReusesAPresentImage(t *testing.T) {
	run, calls := operatorRunner(true)
	if _, err := InstallAgentOperator(run, "da-example"); err != nil {
		t.Fatal(err)
	}
	for _, c := range *calls {
		if strings.HasPrefix(c, "env ") || strings.HasPrefix(c, "docker build") {
			t.Errorf("built although the image is present: %s", c)
		}
	}
	if !containsPrefix(*calls, "helm upgrade --install agent-operator ") {
		t.Error("the chart was not installed")
	}
}

// TestAgentOperatorChartDeclaresItsRBAC pins the chart's RBAC to what the
// reconciler touches (srd009 R4.1).
func TestAgentOperatorChartDeclaresItsRBAC(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", operatorChart, "templates", "rbac.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	rbac := string(data)
	for _, want := range []string{"resources: [agents]", "resources: [agents/status]", "resources: [deployments]",
		"resources: [services]", "resources: [configmaps]", "resources: [customresourcedefinitions]"} {
		if !strings.Contains(rbac, want) {
			t.Errorf("rbac.yaml lacks %q", want)
		}
	}
	if strings.Contains(rbac, "resources: [\"*\"]") || strings.Contains(rbac, "verbs: [\"*\"]") {
		t.Error("rbac.yaml grants a wildcard")
	}
}

// TestAgentOperatorChartGivesTheLoadersAWritableTmp pins applications srd009
// R1.2's operational need: the operator inspects profiles through agent-core's
// loaders, which write a temporary tree, while the chart keeps the root
// filesystem read-only and points --core-root at the tools the image carries.
func TestAgentOperatorChartGivesTheLoadersAWritableTmp(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", operatorChart, "templates", "deployment.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	deployment := string(data)
	for _, want := range []string{"readOnlyRootFilesystem: true", "{name: tmp, mountPath: /tmp}", "{name: tmp, emptyDir: {}}", "--core-root=/opt/agent-core"} {
		if !strings.Contains(deployment, want) {
			t.Errorf("deployment.yaml lacks %q", want)
		}
	}
}

func containsPrefix(calls []string, prefix string) bool {
	for _, c := range calls {
		if strings.HasPrefix(c, prefix) {
			return true
		}
	}
	return false
}
