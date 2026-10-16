// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const canonicalGuideline = `---
id: eng90-sample
title: Sample
---

<!-- field: introduction -->
## Introduction

We describe a sample practice.

<!-- section: the-practice -->
## The practice

Each module follows the practice.

<!-- field: references -->
## References

- https://www.rfc-editor.org/rfc/rfc2119
`

// newFormatRepository makes a git repository holding the design constitution
// and the given documents, tracked.
func newFormatRepository(t *testing.T, documents map[string]string) string {
	t.Helper()
	root := t.TempDir()
	design, err := os.ReadFile(filepath.Join("..", specCheckDesign))
	if err != nil {
		t.Fatal(err)
	}
	documents["agent-core/"+specCheckDesign] = string(design)
	for path, body := range documents {
		target := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(target, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, args := range [][]string{{"init", "-q"}, {"add", "-A"}} {
		command := exec.Command("git", args...)
		command.Dir = root
		if out, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	return root
}

// TestFmtCanonicalizesAndIsIdempotent formats a repository holding one
// canonical guideline and one with stray blank lines: only the second
// changes, it then equals the canonical form, and a second run changes
// nothing (#2583 AC3).
func TestFmtCanonicalizesAndIsIdempotent(t *testing.T) {
	root := newFormatRepository(t, map[string]string{
		"docs/engineering/eng90-sample.md": canonicalGuideline,
		"docs/engineering/eng91-sample.md": replaceID(canonicalGuideline, "eng91-sample") + "\n\n",
	})
	changed, err := formatScope(root, "all")
	if err != nil {
		t.Fatal(err)
	}
	if len(changed) != 1 || changed[0] != "docs/engineering/eng91-sample.md" {
		t.Fatalf("changed = %v, want only eng91-sample.md", changed)
	}
	data, err := os.ReadFile(filepath.Join(root, "docs/engineering/eng91-sample.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != replaceID(canonicalGuideline, "eng91-sample") {
		t.Fatalf("formatted document is not canonical:\n%s", data)
	}
	if again, err := formatScope(root, "all"); err != nil || len(again) != 0 {
		t.Fatalf("second run changed %v (err %v)", again, err)
	}
}

func TestFmtReportsADocumentItCannotRead(t *testing.T) {
	root := newFormatRepository(t, map[string]string{
		"docs/engineering/eng90-sample.md": "---\nid: eng90-sample\ntitle: Sample\n---\n\n## Untagged\n",
	})
	if _, err := formatScope(root, "engineering_guideline"); err == nil {
		t.Fatal("formatted a document with findings")
	}
}

func replaceID(document, id string) string {
	return "---\nid: " + id + document[len("---\nid: eng90-sample"):]
}
