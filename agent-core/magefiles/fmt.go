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
	changed, err := formatScope(root, scope)
	for _, path := range changed {
		fmt.Println("formatted " + path)
	}
	return err
}

// formatScope formats every Markdown corpus document in scope and returns
// the ones it changed. A document with findings is reported and skipped; the
// rest are still formatted.
func formatScope(root, scope string) ([]string, error) {
	declarations, err := loadSpecDeclarationsFrom(root)
	if err != nil {
		return nil, err
	}
	documents, err := documentsInScope(root, declarations, scope)
	if err != nil {
		return nil, err
	}
	var changed []string
	var failures []error
	for _, document := range documents {
		if !document.markdown {
			continue
		}
		path := filepath.Join(root, document.path)
		data, err := os.ReadFile(path)
		if err != nil {
			return changed, err
		}
		formatted, findings, err := document.declaration.Format(data, document.path)
		switch {
		case err != nil:
			failures = append(failures, err)
		case formatted == nil:
			failures = append(failures, fmt.Errorf("%s: %d finding(s), first %s", document.path, len(findings), findings[0]))
		case string(formatted) != string(data):
			if err := os.WriteFile(path, formatted, 0o644); err != nil {
				return changed, err
			}
			changed = append(changed, document.path)
		}
	}
	return changed, errors.Join(failures...)
}
