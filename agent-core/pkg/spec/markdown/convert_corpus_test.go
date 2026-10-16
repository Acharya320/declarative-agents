// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown_test

import (
	"os/exec"
	"strings"
	"testing"
)

// TestTheTrackedCorpusIsMarkdown requires every tracked corpus document to
// be Markdown: no file outside testdata sits at a type's former YAML location,
// where only mage convert still reads one (GH-2626), and the Markdown corpus
// holds at least the 368 documents the migration converted.
func TestTheTrackedCorpusIsMarkdown(t *testing.T) {
	declarations := repositoryDeclarations(t)
	out, err := exec.Command("git", "-C", repositoryRoot, "ls-files", "*.yaml").Output()
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if _, ok := declarations.ForSource(path); ok && !strings.Contains(path, "/testdata/") {
			t.Errorf("%s: a corpus document in YAML; run mage convert %s", path, path)
		}
	}
	if n := len(trackedMarkdownDocuments(t, declarations)); n < 368 {
		t.Errorf("found %d Markdown corpus documents, want at least the 368 the migration converted", n)
	}
}
