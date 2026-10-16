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

// TestKeynameTableReadsAConvertedConfigFormat converts a copy of a config
// format to Markdown and requires the same table rows from its YAML path.
func TestKeynameTableReadsAConvertedConfigFormat(t *testing.T) {
	source := filepath.Join("..", "..", "agent-core", "docs", "specs", "config-formats", "machine-format.yaml")
	want, err := yamlSchemaRows(source)
	if err != nil {
		t.Fatal(err)
	}
	if len(want) == 0 {
		t.Fatal("machine-format.yaml yields no schema rows")
	}
	root := t.TempDir()
	target := filepath.Join(root, "docs", "specs", "config-formats", "machine-format.yaml")
	data, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, data, 0o644); err != nil {
		t.Fatal(err)
	}
	declarations, err := markdown.Default()
	if err != nil {
		t.Fatal(err)
	}
	if reports, err := declarations.ConvertTree(root); err != nil || len(reports) != 1 {
		t.Fatalf("convert: %v (%d reports)", err, len(reports))
	}
	got, err := yamlSchemaRows(target)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("rows from Markdown differ:\n%v\nwant\n%v", got, want)
	}
}
