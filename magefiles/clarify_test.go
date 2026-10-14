// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeClarifyFixture(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestPrintOpenClarifications(t *testing.T) {
	root := t.TempDir()
	module := filepath.Join(root, "mod-a")
	writeClarifyFixture(t, filepath.Join(module, "docs", "specs", "clean.yaml"),
		"id: srd001-clean\ntitle: Clean\n")
	writeClarifyFixture(t, filepath.Join(module, "docs", "specs", "blocked.yaml"),
		`id: srd002-blocked
needs_clarification:
  - id: R1.1
    question: How many retries?
  - id: AC1
    question: What latency bound?
`)
	writeClarifyFixture(t, filepath.Join(module, "docs", "notes.md"),
		"not yaml, skipped by extension\n")

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	var out strings.Builder
	count, err := printOpenClarifications([]string{"mod-a", "mod-missing"}, &out)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("count = %d, want 2", count)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("lines = %d, want 2:\n%s", len(lines), out.String())
	}
	want0 := "srd002-blocked R1.1 -- How many retries?"
	if !strings.Contains(lines[0], want0) {
		t.Errorf("line 0 = %q, want it to contain %q", lines[0], want0)
	}
	want1 := "srd002-blocked AC1 -- What latency bound?"
	if !strings.Contains(lines[1], want1) {
		t.Errorf("line 1 = %q, want it to contain %q", lines[1], want1)
	}
	if !strings.HasPrefix(lines[0], filepath.Join("mod-a", "docs", "specs", "blocked.yaml")+":") {
		t.Errorf("line 0 = %q, want the module-relative file prefix", lines[0])
	}
}

func TestPrintOpenClarificationsCleanCorpus(t *testing.T) {
	root := t.TempDir()
	writeClarifyFixture(t, filepath.Join(root, "mod-a", "docs", "specs", "clean.yaml"),
		"id: srd001-clean\ntitle: Clean\n")

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	var out strings.Builder
	count, err := printOpenClarifications([]string{"mod-a"}, &out)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("count = %d, want 0", count)
	}
	if out.String() != "" {
		t.Fatalf("output = %q, want empty", out.String())
	}
}
