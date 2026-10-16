// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ReadAsYAML reads a corpus document as YAML for code written against the
// YAML corpus. A YAML path reads the file as written; once the document has
// moved to Markdown during the migration window, the same path reads the
// Markdown form as the tree its declaration parses it into. A Markdown path
// reads that tree directly.
func ReadAsYAML(path string) ([]byte, error) {
	declarations, err := Default()
	if err != nil {
		return nil, err
	}
	slashed := filepath.ToSlash(path)
	if strings.HasSuffix(slashed, ".md") {
		return markdownAsYAML(declarations, path)
	}
	data, err := os.ReadFile(path)
	if err == nil || !errors.Is(err, os.ErrNotExist) {
		return data, err
	}
	d, ok := declarations.ForSource(slashed)
	if !ok {
		return nil, err
	}
	moved := d.MarkdownPath(path)
	if _, statErr := os.Stat(moved); statErr != nil {
		return nil, err
	}
	return markdownAsYAML(declarations, moved)
}

func markdownAsYAML(declarations Declarations, path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	d, ok := declarations.ForMarkdown(path)
	if !ok {
		return nil, fmt.Errorf("%s: no document type declares this Markdown location", path)
	}
	doc, findings, err := d.Parse(data, path)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, fmt.Errorf("%s: %d Markdown format finding(s), first %s", path, len(findings), findings[0])
	}
	return yaml.Marshal(doc)
}
