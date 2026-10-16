// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"path/filepath"
	"testing"
)

// TestKeynameTableReadsAMarkdownConfigFormat reads a config format's schema
// rows from its Markdown document; its former YAML path no longer reads
// (GH-2626).
func TestKeynameTableReadsAMarkdownConfigFormat(t *testing.T) {
	formats := filepath.Join("..", "..", "agent-core", "docs", "specs", "config-formats")
	rows, err := yamlSchemaRows(filepath.Join(formats, "machine-format.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) == 0 {
		t.Fatal("machine-format yields no schema rows")
	}
	if _, err := yamlSchemaRows(filepath.Join(formats, "machine-format.yaml")); err == nil {
		t.Fatal("read a config format by its former YAML path")
	}
}
