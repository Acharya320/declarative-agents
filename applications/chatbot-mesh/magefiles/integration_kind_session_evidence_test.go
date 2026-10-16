// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/Nokia-Bell-Labs/declarative-agents/magefiles/kindrig"
	"gopkg.in/yaml.v3"
)

// TestBatchedScenarioFailureCapturesNamespaceBeforeCleanup fails a scenario
// inside a concurrent batch, where the session defers its poison until the
// batch ends. The scenario's kubectl runner stops working once its cleanup
// removes the kubeconfig, as kubectl falls back to localhost:8080 in a real
// run; the namespace evidence must already be on disk by then (GH-2568).
func TestBatchedScenarioFailureCapturesNamespaceBeforeCleanup(t *testing.T) {
	session := newIntegrationKindSession(t.TempDir())
	deactivate, err := activateIntegrationKindSession(session)
	if err != nil {
		t.Fatal(err)
	}
	defer deactivate()
	if err := session.beginConcurrentBatch(); err != nil {
		t.Fatal(err)
	}

	var kubeconfigRemoved atomic.Bool
	kubectl := func(name string, args ...string) ([]byte, error) {
		if kubeconfigRemoved.Load() {
			return []byte("The connection to the server localhost:8080 was refused"), errors.New("exit status 1")
		}
		return []byte("ok " + strings.Join(args, " ")), nil
	}
	var kindCalls []string
	kindRun := func(args ...string) ([]byte, error) {
		kindCalls = append(kindCalls, strings.Join(args, " "))
		return nil, nil
	}
	evidenceDir := filepath.Join(t.TempDir(), "helm-llm-tier")
	evidence := kindrig.FailureEvidence{
		Directory: evidenceDir, Namespaces: []string{"da-helm-llm-tier"}, Run: kubectl,
	}
	cause := errors.New("scenario failed")
	cluster := kindrig.Cluster{Name: kindrig.PlatformClusterName}
	if !releaseAggregateKindCluster(cluster, kindRun, evidence, cause) {
		t.Fatal("session did not adopt the scenario's cluster")
	}
	kubeconfigRemoved.Store(true) // the scenario's deferred cleanup
	session.endConcurrentBatch(cause)

	manifest := readEvidenceManifest(t, evidenceDir)
	if len(manifest.CaptureErrors) != 0 {
		t.Fatalf("scenario evidence capture errors: %v", manifest.CaptureErrors)
	}
	if !manifestHasNamespaceFile(manifest, "da-helm-llm-tier") {
		t.Fatalf("scenario manifest has no namespace evidence: %+v", manifest.Files)
	}
	clusterManifest := readEvidenceManifest(t, filepath.Join(evidenceDir, "cluster"))
	if len(clusterManifest.Namespaces) != 0 || len(clusterManifest.CaptureErrors) != 0 {
		t.Fatalf("session evidence = %+v, want the kind export alone", clusterManifest)
	}
	if len(kindCalls) == 0 || !strings.HasPrefix(kindCalls[0], "export logs ") {
		t.Fatalf("session kind calls = %v, want the cluster-wide log export", kindCalls)
	}
}

func readEvidenceManifest(t *testing.T, dir string) kindrig.EvidenceManifest {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, kindrig.EvidenceManifestFile))
	if err != nil {
		t.Fatal(err)
	}
	var manifest kindrig.EvidenceManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	return manifest
}

func manifestHasNamespaceFile(manifest kindrig.EvidenceManifest, namespace string) bool {
	for _, file := range manifest.Files {
		if file.Namespace == namespace {
			return true
		}
	}
	return false
}
