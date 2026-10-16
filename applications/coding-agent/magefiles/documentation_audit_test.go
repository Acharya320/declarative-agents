// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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

// TestAuditDocumentationOnTheCorpus audits the module's corpus as tracked,
// in Markdown since GH-2584.
func TestAuditDocumentationOnTheCorpus(t *testing.T) {
	if err := auditDocumentation(".."); err != nil {
		t.Fatal(err)
	}
}

// TestAuditDocumentationReportsAMissingField deletes a required section
// from the vision and requires the audit to name the document and the field.
func TestAuditDocumentationReportsAMissingField(t *testing.T) {
	root := copyDocs(t)
	vision := filepath.Join(root, "docs", "VISION.md")
	data, err := os.ReadFile(vision)
	if err != nil {
		t.Fatal(err)
	}
	head, _, found := strings.Cut(string(data), "<!-- field: not -->\n")
	if !found {
		t.Fatal("VISION.md has no not field")
	}
	if err := os.WriteFile(vision, []byte(head), 0o644); err != nil {
		t.Fatal(err)
	}
	err = auditDocumentation(root)
	if err == nil || !strings.Contains(err.Error(), "docs/VISION.md") || !strings.Contains(err.Error(), "not") {
		t.Fatalf("error = %v, want docs/VISION.md reported missing not", err)
	}
}
