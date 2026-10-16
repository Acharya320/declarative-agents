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

// rollBack converts every declared Markdown corpus document under the copy
// back to YAML, as mage convertBack does, so the audit can be proven on the
// YAML form the corpus had before GH-2584.
func rollBack(t *testing.T, root string) {
	t.Helper()
	declarations, err := markdown.Default()
	if err != nil {
		t.Fatal(err)
	}
	rolled := 0
	err = filepath.WalkDir(filepath.Join(root, "docs"), func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}
		d, ok := declarations.ForMarkdown(filepath.ToSlash(path))
		if !ok {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		out, report, err := d.ConvertMarkdown(data, path)
		if err != nil {
			return err
		}
		rolled++
		if err := os.WriteFile(report.Target, out, 0o644); err != nil {
			return err
		}
		return os.Remove(path)
	})
	if err != nil {
		t.Fatal(err)
	}
	if rolled == 0 {
		t.Fatal("rolled back no documents")
	}
}

// TestAuditDocumentationOnTheCorpus audits the module's corpus as tracked,
// in Markdown since GH-2584.
func TestAuditDocumentationOnTheCorpus(t *testing.T) {
	if err := auditDocumentation(".."); err != nil {
		t.Fatal(err)
	}
}

// TestAuditDocumentationOnTheYAMLCorpus rolls a copy back to YAML and
// requires the same audit to pass on it.
func TestAuditDocumentationOnTheYAMLCorpus(t *testing.T) {
	root := copyDocs(t)
	rollBack(t, root)
	if err := auditDocumentation(root); err != nil {
		t.Fatalf("audit of the YAML corpus: %v", err)
	}
}

func TestAuditDocumentationReportsAMissingField(t *testing.T) {
	root := copyDocs(t)
	rollBack(t, root)
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
