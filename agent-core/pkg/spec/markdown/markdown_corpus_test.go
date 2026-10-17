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
	"gopkg.in/yaml.v3"
)

// trackedMarkdownDocuments lists every tracked Markdown document at a
// declared corpus location: the documents the migration converted, and those
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

// TestEveryMarkdownCorpusDocumentIsCanonical holds the corpus to the laws:
// each document parses with no finding of tiers 1 to 3, it is in canonical
// form, and converting the YAML tree it parses into reproduces it, so mage
// convert and the readers agree on every document.
func TestEveryMarkdownCorpusDocumentIsCanonical(t *testing.T) {
	documents := trackedMarkdownDocuments(t, repositoryDeclarations(t))
	indexes := corpusLinkIndexes(t, documents)
	for path, d := range documents {
		t.Run(path, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(repositoryRoot, path))
			if err != nil {
				t.Fatal(err)
			}
			formatted, findings, err := d.FormatLinked(data, path, indexes.For(path))
			if err != nil || formatted == nil {
				t.Fatalf("does not parse: %v %v", err, findings)
			}
			if string(formatted) != string(data) {
				t.Errorf("not in canonical form; run mage fmt")
			}
			tree, _, err := d.Parse(data, path)
			if err != nil || tree == nil {
				t.Fatalf("does not parse: %v", err)
			}
			back, err := yaml.Marshal(tree)
			if err != nil {
				t.Fatal(err)
			}
			var reread yaml.Node
			if err := yaml.Unmarshal(back, &reread); err != nil {
				t.Fatal(err)
			}
			canonical, _ := d.Canonicalize(&reread)
			again, err := d.PrintLinked(canonical, path, indexes.For(path))
			if err != nil {
				t.Fatal(err)
			}
			if string(again) != string(formatted) {
				t.Errorf("Markdown -> YAML -> Markdown changed the document")
			}
		})
	}
}

// corpusLinkIndexes indexes the tracked corpus documents as mage fmt does.
func corpusLinkIndexes(t *testing.T, documents map[string]*markdown.Declaration) markdown.LinkIndexes {
	t.Helper()
	indexes := markdown.LinkIndexes{}
	for path := range documents {
		data, err := os.ReadFile(filepath.Join(repositoryRoot, path))
		if err != nil {
			t.Fatal(err)
		}
		indexes.Add(path, data)
	}
	return indexes
}
