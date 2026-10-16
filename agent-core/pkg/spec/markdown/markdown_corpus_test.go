// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec/markdown"
)

// trackedMarkdownDocuments lists every tracked Markdown document at a
// declared corpus location: the documents already converted, and those
// written in Markdown.
func trackedMarkdownDocuments(t *testing.T, declarations markdown.Declarations) map[string]*markdown.Declaration {
	t.Helper()
	out, err := exec.Command("git", "-C", repositoryRoot, "ls-files", "*.md").Output()
	if err != nil {
		t.Fatal(err)
	}
	documents := map[string]*markdown.Declaration{}
	for _, path := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if strings.Contains(path, "/testdata/") {
			continue
		}
		if d, ok := declarations.ForMarkdown(path); ok {
			documents[path] = d
		}
	}
	return documents
}

// TestEveryMarkdownCorpusDocumentIsCanonical holds the corpus already in
// Markdown to the laws: it parses with no finding of tiers 1 to 3, it is
// in canonical form, and rolling it back to YAML and converting again
// reproduces it.
func TestEveryMarkdownCorpusDocumentIsCanonical(t *testing.T) {
	for path, d := range trackedMarkdownDocuments(t, repositoryDeclarations(t)) {
		t.Run(path, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(repositoryRoot, path))
			if err != nil {
				t.Fatal(err)
			}
			formatted, findings, err := d.Format(data, path)
			if err != nil || formatted == nil {
				t.Fatalf("does not parse: %v %v", err, findings)
			}
			if string(formatted) != string(data) {
				t.Errorf("not in canonical form; run mage fmt")
			}
			back, _, err := d.ConvertMarkdown(data, path)
			if err != nil {
				t.Fatal(err)
			}
			again, _, err := d.ConvertYAML(back, d.SourcePath(path))
			if err != nil {
				t.Fatal(err)
			}
			if string(again) != string(formatted) {
				t.Errorf("Markdown -> YAML -> Markdown changed the document")
			}
		})
	}
}
