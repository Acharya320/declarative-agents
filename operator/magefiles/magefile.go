// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

// Mage targets for the operator module: the root gates dispatch build,
// test, lint, and audit here like every sub-module.
package main

import (
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Default builds the operator.
var Default = Build

// run executes a command in the module root, where mage runs these targets.
func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}

// Build compiles every package.
func Build() error { return run("go", "build", "./...") }

// Test runs the unit tests.
func Test() error { return run("go", "test", "./...") }

// Audit vets and tests the module; the operator has no document corpus of its own.
func Audit() error {
	if err := run("go", "vet", "./..."); err != nil {
		return err
	}
	return Test()
}

// Stats prints the module's Go file and line counts as JSON for the root
// stats gate; the operator declares no agents.
func Stats() error {
	files, lines := 0, 0
	err := filepath.WalkDir(".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() && (entry.Name() == "magefiles" || entry.Name() == ".git") {
			return filepath.SkipDir
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		files++
		lines += strings.Count(string(data), "\n")
		return nil
	})
	if err != nil {
		return err
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(map[string]int{"go_files": files, "go_lines": lines})
}
