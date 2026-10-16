// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec/markdown"
	"gopkg.in/yaml.v3"
)

// corpusDocumentNode reads a corpus document as the YAML tree the gates
// inspect: a Markdown document as the tree its design.yaml declaration parses
// it into (eng10), and one of the YAML documents that stay YAML (the design
// constitution, a corpus charter) as written. A Markdown document no
// declaration covers, such as a guide, returns ok false: it has no fields.
func corpusDocumentNode(path, candidate string) (*yaml.Node, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false, err
	}
	if !strings.HasSuffix(candidate, ".md") {
		var document yaml.Node
		if err := yaml.Unmarshal(data, &document); err != nil {
			return nil, false, err
		}
		return &document, true, nil
	}
	declarations, err := markdown.Default()
	if err != nil {
		return nil, false, err
	}
	d, ok := declarations.ForMarkdown(candidate)
	if !ok {
		return nil, false, nil
	}
	document, findings, err := d.Parse(data, candidate)
	if err != nil {
		return nil, false, err
	}
	if document == nil {
		return nil, false, fmt.Errorf("%d Markdown format finding(s), first %s", len(findings), findings[0])
	}
	return document, true, nil
}

// readCorpusYAML reads a Markdown corpus document as the YAML tree its
// declaration parses it into.
func readCorpusYAML(path string) ([]byte, error) {
	node, ok, err := corpusDocumentNode(path, filepath.ToSlash(path))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if !ok {
		return nil, fmt.Errorf("%s: not a declared Markdown corpus document", path)
	}
	return yaml.Marshal(node)
}
