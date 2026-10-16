// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec/markdown"
)

// Convert converts corpus documents from YAML to Markdown in place (eng10):
// it writes each canonical Markdown document, removes its YAML source, and
// prints what the conversion changed beyond the content. The scope is a
// document path relative to the repository root, a document type such as
// use_case, or all. A Markdown document already in scope is rewritten in
// canonical form, so converting twice changes nothing. Convert is the
// migration bridge (GH-2578) and leaves with it (GH-2585).
func Convert(scope string) error {
	return convertScope(scope, false)
}

// ConvertBack converts corpus documents from Markdown back to YAML in place,
// the rollback for the migration window. The scope is as for Convert.
func ConvertBack(scope string) error {
	return convertScope(scope, true)
}

type conversion struct {
	path        string
	declaration *markdown.Declaration
	markdown    bool
}

func convertScope(scope string, back bool) error {
	root, err := repositoryRoot()
	if err != nil {
		return err
	}
	declarations, err := loadSpecDeclarations()
	if err != nil {
		return err
	}
	documents, err := documentsInScope(root, declarations, scope)
	if err != nil {
		return err
	}
	if len(documents) == 0 {
		return fmt.Errorf("convert: no corpus document matches %q", scope)
	}
	for _, document := range documents {
		if err := convertDocument(root, document, back); err != nil {
			return err
		}
	}
	return nil
}

func loadSpecDeclarations() (markdown.Declarations, error) {
	data, err := os.ReadFile(specCheckDesign)
	if err != nil {
		return nil, err
	}
	return markdown.LoadDeclarations(data)
}

func repositoryRoot() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("convert: find the repository root: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// documentsInScope lists the tracked corpus documents a scope names, in both
// formats; testdata is never in scope.
func documentsInScope(root string, declarations markdown.Declarations, scope string) ([]conversion, error) {
	out, err := exec.Command("git", "-C", root, "ls-files", "*.yaml", "*.md").Output()
	if err != nil {
		return nil, fmt.Errorf("convert: list tracked files: %w", err)
	}
	var documents []conversion
	for _, path := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		document, ok := classify(declarations, path)
		if !ok || strings.Contains(path, "/testdata/") {
			continue
		}
		if scope == "all" || scope == document.declaration.Type() || scope == path {
			documents = append(documents, document)
		}
	}
	return documents, nil
}

func classify(declarations markdown.Declarations, path string) (conversion, bool) {
	if strings.HasSuffix(path, ".md") {
		d, ok := declarations.ForMarkdown(path)
		return conversion{path: path, declaration: d, markdown: true}, ok
	}
	d, ok := declarations.ForSource(path)
	return conversion{path: path, declaration: d}, ok
}

// convertDocument converts one document in the direction asked, or puts a
// document already in the target format in canonical form.
func convertDocument(root string, document conversion, back bool) error {
	source := filepath.Join(root, document.path)
	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	d := document.declaration
	var out []byte
	var report markdown.Report
	switch {
	case back && document.markdown:
		out, report, err = d.ConvertMarkdown(data, document.path)
	case !back && !document.markdown:
		out, report, err = d.ConvertYAML(data, document.path)
	case !back:
		var findings []markdown.Finding
		out, findings, err = d.Format(data, document.path)
		report = markdown.Report{Source: document.path, Target: document.path, Findings: findings}
		if err == nil && out == nil {
			err = fmt.Errorf("%s: %d finding(s) before it can be formatted", document.path, len(findings))
		}
	default:
		return nil // a YAML document is already in the rollback format
	}
	if err != nil {
		return err
	}
	return writeConversion(root, report, out)
}

func writeConversion(root string, report markdown.Report, out []byte) error {
	target := filepath.Join(root, report.Target)
	if err := os.WriteFile(target, out, 0o644); err != nil {
		return err
	}
	if report.Target != report.Source {
		if err := os.Remove(filepath.Join(root, report.Source)); err != nil {
			return err
		}
	}
	fmt.Println(strings.Join(report.Lines(), "\n"))
	return nil
}
