// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCollectorAgentManifestNamesTheProfileAndPorts pins the proof's Agent
// to srd009 R5.1: the collector profile ConfigMap, the three ports with the
// control port as the readiness signal, and the bind host on every interface.
func TestCollectorAgentManifestNamesTheProfileAndPorts(t *testing.T) {
	for _, want := range []string{"kind: Agent", "configMap: collector-profile", "path: collector/profile.yaml", "image: AGENT_CORE_IMAGE", "mountPath: /work"} {
		if !strings.Contains(collectorAgentManifest, want) {
			t.Errorf("manifest lacks %q", want)
		}
	}
	for _, repeated := range []string{"ports:", "controlPort", "COLLECTOR_BIND_HOST", "1819"} {
		if strings.Contains(collectorAgentManifest, repeated) {
			t.Errorf("manifest repeats %q, which the profile declares (srd009 R1.2)", repeated)
		}
	}
}

// TestProfileFilesListsTheCollectorProfile reads the real profile directory
// and requires its profile, machine, tools, declarations, and REST files.
func TestProfileFilesListsTheCollectorProfile(t *testing.T) {
	args, err := profileFiles(filepath.Join("..", "..", collectorProfile))
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(args, " ")
	for _, want := range []string{"profile.yaml", "machine.yaml", "tools.yaml", "declarations.yaml", "rest.yaml", "query-declarations.yaml"} {
		if !strings.Contains(joined, "--from-file=collector__"+want+"=") {
			t.Errorf("profile files lack %s keyed under collector: %v", want, args)
		}
	}
	if !strings.Contains(joined, "--from-file=units__types-collector.yaml=") || !strings.Contains(joined, "/units/types-collector.yaml") {
		t.Errorf("profile files lack the imported units file: %v", args)
	}
	if strings.Contains(joined, "/ui") || strings.Contains(joined, "/tests") {
		t.Errorf("profile files include a subdirectory: %v", args)
	}
	if _, err := profileFiles(filepath.Join(os.TempDir(), "no-such-profile")); err == nil {
		t.Error("a missing profile directory was accepted")
	}
}

func TestForwardedPortParses(t *testing.T) {
	m := forwardedPort.FindStringSubmatch("Forwarding from 127.0.0.1:54321 -> 18191")
	if m == nil || m[1] != "54321" {
		t.Fatalf("port = %v", m)
	}
}
