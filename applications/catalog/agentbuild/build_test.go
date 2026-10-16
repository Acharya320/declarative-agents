// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package agentbuild

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildTagSelectsProduction(t *testing.T) {
	t.Parallel()
	if BuildTag != "production" {
		t.Fatalf("BuildTag = %q, want production", BuildTag)
	}
}

func TestBuildReturnsErrorForMissingTarget(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	var log bytes.Buffer
	out, err := Build(dir, filepath.Join(dir, "agent-bin"), &log)
	if err == nil {
		t.Fatal("expected build error for a root without ./cmd/agent")
	}
	if out != "" {
		t.Fatalf("output path = %q, want empty on failure", out)
	}
	if !strings.Contains(err.Error(), dir) {
		t.Fatalf("error %q does not mention coreRoot %q", err, dir)
	}
}

// TestOutputPathSeparatesCheckouts names one binary per agent-core root, so
// audits of two checkouts at once cannot run each other's agent (GH-2593).
func TestOutputPathSeparatesCheckouts(t *testing.T) {
	t.Parallel()
	first := OutputPath("agent", "/work/one/agent-core")
	second := OutputPath("agent", "/work/two/agent-core")
	if first == second {
		t.Fatalf("two checkouts share binary %s", first)
	}
	if again := OutputPath("agent", "/work/one/../one/agent-core"); again != first {
		t.Fatalf("one checkout names two binaries: %s and %s", first, again)
	}
	if filepath.Dir(first) != filepath.Clean(os.TempDir()) || !strings.HasPrefix(filepath.Base(first), "agent-") {
		t.Fatalf("binary %s is not agent-<digest> under the temp directory", first)
	}
}

// TestBuildLeavesNoPartialBinary fails a build and requires neither the
// output nor its staging file to exist.
func TestBuildLeavesNoPartialBinary(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	out := filepath.Join(dir, "agent-bin")
	if _, err := Build(dir, out, nil); err == nil {
		t.Fatal("expected build error for a root without ./cmd/agent")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "agent-bin") {
			t.Fatalf("failed build left %s", entry.Name())
		}
	}
}
