// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package specmd

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec"
	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec/markdown"
)

// knownCorpusDefects lists SRDs whose source YAML is itself defective, so the
// round trip cannot be judged on them. Each entry names its fix; once the fix
// lands the entry must go, and the test fails until it does.
// GH-2590 fixed the last one, catalog srd021's duplicate mapping key.
var knownCorpusDefects = map[string]string{}

// TestEveryCorpusSRDRoundTrips runs the srd004 laws over the whole tracked
// SRD corpus: yaml -> model -> md -> model loses nothing, and the printed
// Markdown is a fixed point of parse-then-print.
func TestEveryCorpusSRDRoundTrips(t *testing.T) {
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	out, err := exec.Command("git", "-C", repo, "ls-files").Output()
	if err != nil {
		t.Fatal(err)
	}
	checked := 0
	for _, path := range strings.Split(string(out), "\n") {
		if !strings.Contains(path, "/docs/specs/software-requirements/srd") ||
			!strings.HasSuffix(path, ".yaml") || strings.Contains(path, "/testdata/") {
			continue
		}
		checked++
		failure := roundTripFailure(filepath.Join(repo, path), path)
		reason, known := knownCorpusDefects[path]
		switch {
		case failure != "" && !known:
			t.Errorf("%s: %s", path, failure)
		case failure == "" && known:
			t.Errorf("%s now round-trips; remove its known-defect entry (%s)", path, reason)
		}
	}
	if checked == 0 {
		t.Fatal("no corpus SRDs found")
	}
}

func roundTripFailure(file, name string) string {
	original, err := spec.ParseSRD(file)
	if err != nil {
		return "yaml parse: " + err.Error()
	}
	md, err := PrintMarkdown(&original)
	if err != nil {
		return "markdown print: " + err.Error()
	}
	back, err := ParseMarkdown(md, name)
	if err != nil {
		return "markdown parse: " + err.Error()
	}
	if diffs := modelDiffs(&original, &back, true); len(diffs) > 0 {
		return "model changed: " + diffs[0]
	}
	again, err := PrintMarkdown(&back)
	if err != nil || string(again) != string(md) {
		return "printed markdown is not a fixed point: " + firstDifference(md, again)
	}
	if _, err := PrintYAML(&original); err != nil {
		return "yaml print: " + err.Error()
	}
	return ""
}

// corpusAllowed are the findings the printed corpus may still carry: prose
// containment (normative wording in non-goals) is resolved file by file in
// the corpus migration (GH-2581 decision 2). Every other finding on a printed
// SRD is either a printer defect or a corpus fact the declaration should
// already admit, such as the widened id grammar (decision 1).
var corpusAllowed = map[string]bool{"prose-containment": true}

// TestPrintedCorpusIsClean diagnoses the printed form of every tracked SRD
// against design.yaml's srd declaration: nothing but prose containment may
// remain, so the printer is structurally sound, every id fits the declared
// grammar, and every criterion carries a visible traces field.
func TestPrintedCorpusIsClean(t *testing.T) {
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	out, err := exec.Command("git", "-C", repo, "ls-files").Output()
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range strings.Split(string(out), "\n") {
		if !strings.Contains(path, "/docs/specs/software-requirements/srd") ||
			!strings.HasSuffix(path, ".yaml") || strings.Contains(path, "/testdata/") {
			continue
		}
		srd, err := spec.ParseSRD(filepath.Join(repo, path))
		if err != nil {
			continue
		}
		md, err := PrintMarkdown(&srd)
		if err != nil {
			t.Errorf("%s: %v", path, err)
			continue
		}
		name := strings.TrimSuffix(path, ".yaml") + ".md"
		for _, f := range diagnose(t, md, name) {
			if !corpusAllowed[f.Name] {
				t.Errorf("%s: printed form has %s", path, f)
			}
		}
	}
}

// diagnose runs the promoted diagnostics against this repository's
// declarations.
func diagnose(t *testing.T, data []byte, source string) []markdown.Finding {
	t.Helper()
	declarations, err := RepositoryDeclarations(".")
	if err != nil {
		t.Fatal(err)
	}
	return markdown.Diagnose(declarations, data, source)
}
