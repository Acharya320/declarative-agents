// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const extRefParentSRD = `id: srd013-contract
requirements:
  R2:
    title: Group
    items:
      - R2.4: Expose the applier contract.
`

const extRefParentSRDRetired = `id: srd013-contract
requirements:
  R2:
    title: Group
    items:
      - R2.4: Expose the applier contract.
        status: retired
        retired_in: v0.20270102.0
        replaced_by: [R2.5]
      - R2.5: Expose the lease-holding applier contract.
`

// swapExtGit replaces the external-reference git reader for one test.
func swapExtGit(t *testing.T, fn func(dir string, args ...string) (string, error)) {
	t.Helper()
	orig := extGitOutput
	extGitOutput = fn
	t.Cleanup(func() { extGitOutput = orig })
}

// fakeExtCache answers toplevel discovery and one cache's tag/tree surface.
type fakeExtCache struct {
	toplevel string
	cacheDir string
	tags     string // newline list, newest first
	srd      string // served for the srd013-contract file at any ref
}

func (f fakeExtCache) output(dir string, args ...string) (string, error) {
	if args[0] == "rev-parse" && args[1] == "--show-toplevel" {
		return f.toplevel, nil
	}
	if dir != f.cacheDir {
		return "", errors.New("not a repository: " + dir)
	}
	switch args[0] {
	case "rev-parse":
		return ".git", nil
	case "tag":
		return f.tags, nil
	case "ls-tree":
		return "agent-core/docs/specs/software-requirements/srd013-contract.yaml", nil
	case "show":
		return f.srd, nil
	case "ls-remote":
		var lines []string
		for _, tag := range strings.Split(f.tags, "\n") {
			if tag != "" {
				lines = append(lines, "deadbeef\trefs/tags/"+tag)
			}
		}
		return strings.Join(lines, "\n"), nil
	}
	return "", errors.New("unexpected git call: " + strings.Join(args, " "))
}

// writeExtRefModule lays out one module with an SRD carrying a derived_from
// entry and returns the repo-root directory the test chdirs into.
func writeExtRefModule(t *testing.T, release string, withPins bool) string {
	t.Helper()
	parent := t.TempDir()
	root := filepath.Join(parent, "child")
	srdDir := filepath.Join(root, "mod-a", "docs", "specs", "software-requirements")
	if err := os.MkdirAll(srdDir, 0o755); err != nil {
		t.Fatal(err)
	}
	releaseLine := ""
	if release != "" {
		releaseLine = "\n            release: " + release
	}
	srd := `id: srd050-child
requirements:
  R1:
    title: Group
    items:
      - R1.1: Implement the contract.
        derived_from:
          - repo: github.com/example/parent
            srd: srd013-contract
            item: R2.4` + releaseLine + "\n"
	if err := os.WriteFile(filepath.Join(srdDir, "srd050-child.yaml"), []byte(srd), 0o644); err != nil {
		t.Fatal(err)
	}
	if withPins {
		pins := "pins:\n  github.com/example/parent: v0.20261001.0\n"
		if err := os.WriteFile(filepath.Join(root, "mod-a", "docs", "external-pins.yaml"), []byte(pins), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	return root
}

func TestExternalReferenceGateRefusesUnpinned(t *testing.T) {
	root := writeExtRefModule(t, "", false)
	swapExtGit(t, fakeExtCache{
		toplevel: root,
		cacheDir: filepath.Join(filepath.Dir(root), "parent"),
		tags:     "v0.20261001.0",
		srd:      extRefParentSRD,
	}.output)
	err := externalReferenceReleaseGate([]string{"mod-a"})
	if err == nil || !strings.Contains(err.Error(), "unpinned") {
		t.Fatalf("err = %v, want unpinned refusal", err)
	}
}

func TestExternalReferenceGateRefusesRetired(t *testing.T) {
	root := writeExtRefModule(t, "v0.20261001.0", true)
	swapExtGit(t, fakeExtCache{
		toplevel: root,
		cacheDir: filepath.Join(filepath.Dir(root), "parent"),
		tags:     "v0.20270102.0\nv0.20261001.0",
		srd:      extRefParentSRDRetired,
	}.output)
	err := externalReferenceReleaseGate([]string{"mod-a"})
	if err == nil || !strings.Contains(err.Error(), "retired at v0.20270102.0") {
		t.Fatalf("err = %v, want retired refusal", err)
	}
}

func TestExternalReferenceGateAcceptsPinnedLive(t *testing.T) {
	root := writeExtRefModule(t, "v0.20261001.0", true)
	swapExtGit(t, fakeExtCache{
		toplevel: root,
		cacheDir: filepath.Join(filepath.Dir(root), "parent"),
		tags:     "v0.20270102.0\nv0.20261001.0",
		srd:      extRefParentSRD,
	}.output)
	if err := externalReferenceReleaseGate([]string{"mod-a"}); err != nil {
		t.Fatalf("err = %v, want nil: a stale-but-live entry is a warning, not a refusal", err)
	}
}

func TestDoctorReportsMissingExternalCache(t *testing.T) {
	root := writeExtRefModule(t, "v0.20261001.0", true)
	swapExtGit(t, fakeExtCache{
		toplevel: root,
		cacheDir: filepath.Join(filepath.Dir(root), "somewhere-else"),
	}.output)
	var out strings.Builder
	if err := reportExternalCaches([]string{"mod-a"}, &out); err != nil {
		t.Fatal(err)
	}
	want := "git clone https://github.com/example/parent " +
		filepath.Join(filepath.Dir(root), "parent")
	if !strings.Contains(out.String(), want) {
		t.Fatalf("output %q, want it to contain %q", out.String(), want)
	}
}

func TestBumpReportsStaleSpecPin(t *testing.T) {
	root := writeExtRefModule(t, "v0.20261001.0", true)
	swapExtGit(t, fakeExtCache{
		toplevel: root,
		cacheDir: filepath.Join(filepath.Dir(root), "parent"),
		tags:     "v0.20270102.0\nv0.20261001.0",
		srd:      extRefParentSRD,
	}.output)
	var out strings.Builder
	reportExternalPinCurrency([]string{"mod-a"}, &out)
	if !strings.Contains(out.String(), "pinned v0.20261001.0, newest release v0.20270102.0") {
		t.Fatalf("output %q, want the stale pin line", out.String())
	}
}
