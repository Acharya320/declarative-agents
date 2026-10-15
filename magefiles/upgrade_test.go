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

const upgradeParentV1 = `id: srd013-contract
requirements:
  R2:
    title: Group
    items:
      - R2.4: Expose the applier contract.
      - R2.5: Report the applier lease.
`

// R2.4 unchanged, R2.5 changed.
const upgradeParentV2 = `id: srd013-contract
requirements:
  R2:
    title: Group
    items:
      - R2.4: Expose the applier contract.
      - R2.5: Report the applier lease and its expiry.
`

// fakeUpgradeCache serves per-ref SRD content plus fetch/tag/ls-tree calls.
type fakeUpgradeCache struct {
	toplevel string
	cacheDir string
	tags     string            // newline list, newest first
	srds     map[string]string // ref -> SRD content
}

func (f fakeUpgradeCache) output(dir string, args ...string) (string, error) {
	if args[0] == "rev-parse" && args[1] == "--show-toplevel" {
		return f.toplevel, nil
	}
	if dir != f.cacheDir {
		return "", errors.New("not a repository: " + dir)
	}
	switch args[0] {
	case "rev-parse":
		if args[1] == "--verify" {
			tag := strings.TrimPrefix(args[2], "refs/tags/")
			for _, known := range strings.Split(f.tags, "\n") {
				if known == tag {
					return tag, nil
				}
			}
			return "", errors.New("unknown tag " + tag)
		}
		return ".git", nil
	case "fetch":
		return "", nil
	case "tag":
		return f.tags, nil
	case "ls-tree":
		return "agent-core/docs/specs/software-requirements/srd013-contract.yaml", nil
	case "show":
		ref, _, _ := strings.Cut(args[1], ":")
		content, ok := f.srds[ref]
		if !ok {
			return "", errors.New("unknown ref " + ref)
		}
		return content, nil
	}
	return "", errors.New("unexpected git call: " + strings.Join(args, " "))
}

const upgradeChildSRD = `# Copyright (c) 2026 Example. All rights reserved.
id: srd050-child
requirements:
  R1:
    title: Group
    items:
      - R1.1: Implement the contract.
        derived_from:
          - repo: github.com/example/parent
            release: v0.20261001.0
            srd: srd013-contract
            item: R2.4
      - R1.2: Report the lease here.
        derived_from:
          - repo: github.com/example/parent
            release: v0.20261001.0
            srd: srd013-contract
            item: R2.5
`

// writeUpgradeModule lays out the child module and chdirs into its root.
func writeUpgradeModule(t *testing.T) (string, string) {
	t.Helper()
	parent := t.TempDir()
	root := filepath.Join(parent, "child")
	srdDir := filepath.Join(root, "mod-a", "docs", "specs", "software-requirements")
	if err := os.MkdirAll(srdDir, 0o755); err != nil {
		t.Fatal(err)
	}
	srdPath := filepath.Join(srdDir, "srd050-child.yaml")
	if err := os.WriteFile(srdPath, []byte(upgradeChildSRD), 0o644); err != nil {
		t.Fatal(err)
	}
	pins := "# Tracked target pins.\npins:\n  github.com/example/parent: v0.20261001.0\n"
	if err := os.WriteFile(filepath.Join(root, "mod-a", "docs", "external-pins.yaml"), []byte(pins), 0o644); err != nil {
		t.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	return root, srdPath
}

func upgradeFixture(t *testing.T) (string, string) {
	t.Helper()
	root, srdPath := writeUpgradeModule(t)
	swapExtGit(t, fakeUpgradeCache{
		toplevel: root,
		cacheDir: filepath.Join(filepath.Dir(root), "parent"),
		tags:     "v0.20270102.0\nv0.20261001.0",
		srds: map[string]string{
			"v0.20261001.0": upgradeParentV1,
			"v0.20270102.0": upgradeParentV2,
		},
	}.output)
	return root, srdPath
}

func TestUpgradeRepinsUnchangedAndListsStale(t *testing.T) {
	root, srdPath := upgradeFixture(t)
	var out strings.Builder
	cache := filepath.Join(filepath.Dir(root), "parent")
	if err := upgradeToRelease([]string{"mod-a"}, "github.com/example/parent", cache, "v0.20270102.0", &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "re-pinned") || !strings.Contains(out.String(), "R2.4 unchanged") {
		t.Fatalf("output %q, want the R2.4 re-pin line", out.String())
	}
	if !strings.Contains(out.String(), "stale:") || !strings.Contains(out.String(), "R2.5 validated at v0.20261001.0") {
		t.Fatalf("output %q, want the R2.5 stale line", out.String())
	}
	if !strings.Contains(out.String(), "1 entry re-pinned, 1 stale") {
		t.Fatalf("output %q, want the summary line", out.String())
	}

	document, err := os.ReadFile(srdPath)
	if err != nil {
		t.Fatal(err)
	}
	text := string(document)
	if strings.Count(text, "release: v0.20270102.0") != 1 || strings.Count(text, "release: v0.20261001.0") != 1 {
		t.Fatalf("document releases wrong:\n%s", text)
	}
	// Everything but the one release scalar survives byte-for-byte.
	want := strings.Replace(upgradeChildSRD,
		"item: R2.4", "item: R2.4", 1)
	want = strings.Replace(want,
		"release: v0.20261001.0\n            srd: srd013-contract\n            item: R2.4",
		"release: v0.20270102.0\n            srd: srd013-contract\n            item: R2.4", 1)
	if text != want {
		t.Fatalf("document diverged beyond the edited scalar:\n--- got ---\n%s\n--- want ---\n%s", text, want)
	}

	pins, err := os.ReadFile(filepath.Join("mod-a", "docs", "external-pins.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(pins), "github.com/example/parent: v0.20270102.0") {
		t.Fatalf("pins not raised:\n%s", pins)
	}
	if !strings.HasPrefix(string(pins), "# Tracked target pins.") {
		t.Fatalf("pins header lost:\n%s", pins)
	}
}

func TestUpgradeStatusListsLaggingEntries(t *testing.T) {
	upgradeFixture(t)
	// Raise the pin without re-pinning: both entries lag.
	pinsPath := filepath.Join("mod-a", "docs", "external-pins.yaml")
	if err := os.WriteFile(pinsPath,
		[]byte("pins:\n  github.com/example/parent: v0.20270102.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	if err := printUpgradeStatus([]string{"mod-a"}, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "R2.4 validated at v0.20261001.0, pin v0.20270102.0") ||
		!strings.Contains(out.String(), "2 entries behind their pins") {
		t.Fatalf("output %q, want both lagging entries and the count", out.String())
	}
}

func TestUpgradeStatusQuietWhenCurrent(t *testing.T) {
	upgradeFixture(t)
	var out strings.Builder
	if err := printUpgradeStatus([]string{"mod-a"}, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "no entries behind their pins") {
		t.Fatalf("output %q, want the quiet line", out.String())
	}
}

func TestUpgradeToNamedRelease(t *testing.T) {
	_, srdPath := upgradeFixture(t)
	var out strings.Builder
	if err := upgradeToNamedRelease([]string{"mod-a"}, "github.com/example/parent", "v0.20270102.0", &out); err != nil {
		t.Fatal(err)
	}
	document, err := os.ReadFile(srdPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(document), "release: v0.20270102.0") != 1 {
		t.Fatalf("want the unchanged entry re-pinned to the named release:\n%s", document)
	}
}

func TestUpgradeToRefusesUnknownRelease(t *testing.T) {
	upgradeFixture(t)
	var out strings.Builder
	err := upgradeToNamedRelease([]string{"mod-a"}, "github.com/example/parent", "v0.20991231.0", &out)
	if err == nil || !strings.Contains(err.Error(), "release v0.20991231.0 is not in the") {
		t.Fatalf("err = %v, want the unknown-release refusal", err)
	}
}

func TestUpgradeRefusesMissingCache(t *testing.T) {
	root, _ := writeUpgradeModule(t)
	swapExtGit(t, fakeUpgradeCache{
		toplevel: root,
		cacheDir: filepath.Join(filepath.Dir(root), "elsewhere"),
	}.output)
	err := Upgrade("github.com/example/parent")
	if err == nil || !strings.Contains(err.Error(), "git clone https://github.com/example/parent") {
		t.Fatalf("err = %v, want the clone instruction", err)
	}
}
