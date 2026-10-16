// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Fmt rewrites Markdown corpus documents in canonical form in place, for the
// whole repository: it runs agent-core's mage fmt, which owns the format
// (eng10). The scope is a document path from the repository root, a document
// type, or all.
func Fmt(scope string) error {
	root, err := absoluteRepositoryRoot()
	if err != nil {
		return err
	}
	command := exec.Command("mage", "fmt", scope)
	command.Dir = filepath.Join(root, "agent-core")
	command.Stdout, command.Stderr = os.Stdout, os.Stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("fmt %s: %w", scope, err)
	}
	return nil
}
