// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec/markdown"
)

// copyDocs copies the module's docs tree into a temporary application root.
func copyDocs(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	source := filepath.Join("..", "docs")
	err := filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(source, path)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		target := filepath.Join(root, "docs", rel)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func TestAuditDocumentationOnTheYAMLCorpus(t *testing.T) {
	if err := auditDocumentation(".."); err != nil {
		t.Fatal(err)
	}
}

// TestAuditDocumentationOnTheMarkdownCorpus converts a copy of the corpus to
// Markdown and requires the same audit to pass (eng10 migration window).
func TestAuditDocumentationOnTheMarkdownCorpus(t *testing.T) {
	root := copyDocs(t)
	declarations, err := markdown.Default()
	if err != nil {
		t.Fatal(err)
	}
	reports, err := declarations.ConvertTree(filepath.Join(root, "docs"))
	if err != nil {
		t.Fatal(err)
	}
	if len(reports) == 0 {
		t.Fatal("converted no documents")
	}
	if err := auditDocumentation(root); err != nil {
		t.Fatalf("audit of the Markdown corpus: %v", err)
	}
}

func TestAuditDocumentationReportsAMissingField(t *testing.T) {
	root := copyDocs(t)
	vision := filepath.Join(root, "docs", "VISION.yaml")
	data, err := os.ReadFile(vision)
	if err != nil {
		t.Fatal(err)
	}
	changed := strings.Replace(string(data), "\nnot:", "\nnot_any_more:", 1)
	if changed == string(data) {
		t.Fatal("VISION.yaml has no top-level not field")
	}
	if err := os.WriteFile(vision, []byte(changed), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := auditDocumentation(root); err == nil || !strings.Contains(err.Error(), "docs/VISION.yaml missing [not]") {
		t.Fatalf("error = %v, want docs/VISION.yaml missing [not]", err)
	}
}
