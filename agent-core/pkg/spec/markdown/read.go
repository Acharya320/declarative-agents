// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ReadAsYAML reads a file as YAML for code written against the document
// model: a Markdown corpus document as the YAML tree its declaration parses
// it into, and any other YAML file, such as a machine or a profile, as
// written. A YAML file at a corpus document's former location is an error:
// the corpus is Markdown, and mage convert migrates such a file once
// (GH-2626).
func ReadAsYAML(path string) ([]byte, error) {
	declarations, err := Default()
	if err != nil {
		return nil, err
	}
	slashed := filepath.ToSlash(path)
	if strings.HasSuffix(slashed, ".md") {
		return markdownAsYAML(declarations, path)
	}
	if _, ok := declarations.ForSource(slashed); ok {
		return nil, fmt.Errorf("%s: corpus documents are Markdown (eng10); run mage convert to migrate it", path)
	}
	return os.ReadFile(path)
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
