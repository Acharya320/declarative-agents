// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec/markdown"
)

// TestKeynameTableReadsAConvertedConfigFormat reads a config format by its
// YAML path from its Markdown form, rolls a copy back to YAML, and requires
// the same table rows from both.
func TestKeynameTableReadsAConvertedConfigFormat(t *testing.T) {
	source := filepath.Join("..", "..", "agent-core", "docs", "specs", "config-formats", "machine-format.yaml")
	got, err := yamlSchemaRows(source)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) == 0 {
		t.Fatal("machine-format yields no schema rows")
	}
	declarations, err := markdown.Default()
	if err != nil {
		t.Fatal(err)
	}
	mdPath := filepath.Join("..", "..", "agent-core", "docs", "specs", "config-formats", "machine-format.md")
	d, ok := declarations.ForMarkdown(filepath.ToSlash(mdPath))
	if !ok {
		t.Fatal("machine-format.md is not a declared config format")
	}
	data, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatal(err)
	}
	rolled, _, err := d.ConvertMarkdown(data, mdPath)
	if err != nil {
		t.Fatal(err)
	}
	copyPath := filepath.Join(t.TempDir(), "docs", "specs", "config-formats", "machine-format.yaml")
	if err := os.MkdirAll(filepath.Dir(copyPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(copyPath, rolled, 0o644); err != nil {
		t.Fatal(err)
	}
	want, err := yamlSchemaRows(copyPath)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("rows from Markdown differ from the YAML rollback:\n%v\nwant\n%v", got, want)
	}
}
