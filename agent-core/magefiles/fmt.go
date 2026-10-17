// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Fmt rewrites Markdown corpus documents in canonical form in place (eng10,
// Tier 0), so the audit's not-canonical finding never needs a hand edit. The
// scope is as for Convert: a document path from the repository root, a
// document type, or all. YAML documents are left alone; formatting a
// formatted document changes nothing.
func Fmt(scope string) error {
	root, err := repositoryRoot()
	if err != nil {
		return err
	}
	changed, fences, err := formatScope(root, scope)
	for _, path := range changed {
		fmt.Println("formatted " + path)
	}
	for _, fence := range fences {
		fmt.Println("yaml fence " + fence)
	}
	return err
}

// formatScope formats every Markdown corpus document in scope and returns
// the ones it changed. A document with findings is reported and skipped; the
// rest are still formatted.
func formatScope(root, scope string) ([]string, []string, error) {
	declarations, err := loadSpecDeclarationsFrom(root)
	if err != nil {
		return nil, nil, err
	}
	documents, err := documentsInScope(root, declarations, scope)
	if err != nil {
		return nil, nil, err
	}
	var changed, fences []string
	var failures []error
	for _, document := range documents {
		if !document.markdown {
			continue
		}
		path := filepath.Join(root, document.path)
		data, err := os.ReadFile(path)
		if err != nil {
			return changed, fences, err
		}
		formatted, findings, err := document.declaration.Format(data, document.path)
		switch {
		case err != nil:
			failures = append(failures, err)
		case formatted == nil:
			failures = append(failures, fmt.Errorf("%s: %d finding(s), first %s", document.path, len(findings), findings[0]))
		case string(formatted) != string(data):
			if err := os.WriteFile(path, formatted, 0o644); err != nil {
				return changed, fences, err
			}
			changed = append(changed, document.path)
		}
		fences = append(fences, documentFences(document, formatted)...)
	}
	return changed, fences, errors.Join(failures...)
}

// documentFences lists the values a formatted document still prints as
// fenced yaml, the reported last resort (eng10, GH-2643).
func documentFences(document conversion, formatted []byte) []string {
	if formatted == nil {
		return nil
	}
	fences, err := document.declaration.Fences(formatted, document.path)
	if err != nil {
		return []string{document.path + ": " + err.Error()}
	}
	out := make([]string, len(fences))
	for i, fence := range fences {
		out[i] = document.path + ": " + fence.String()
	}
	return out
}
