// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Clarify lists every open needs_clarification entry across the audit
// participants: the pre-planning worklist (design.yaml clarification_marker).
// It is a report, not a gate -- the unresolved-clarification audit check owns
// refusal -- so it always exits zero.
func Clarify() error {
	count, err := printOpenClarifications(auditParticipants(), os.Stdout)
	if err != nil {
		return err
	}
	if count == 0 {
		fmt.Println("no open clarifications")
		return nil
	}
	fmt.Printf("%d open clarification(s)\n", count)
	return nil
}

// clarifiedDocument is the minimal slice of a corpus document the worklist
// needs: its id and its open questions.
type clarifiedDocument struct {
	ID                 string `yaml:"id"`
	NeedsClarification []struct {
		ID       string `yaml:"id"`
		Question string `yaml:"question"`
	} `yaml:"needs_clarification"`
}

// printOpenClarifications walks each module's docs tree and prints one line
// per entry, sorted by file path within a module. Files that are not YAML
// mappings are skipped; the audit owns format enforcement.
func printOpenClarifications(modules []string, out io.Writer) (int, error) {
	total := 0
	for _, module := range modules {
		docsDir := filepath.Join(module, "docs")
		if _, err := os.Stat(docsDir); err != nil {
			continue
		}
		files, err := yamlDocFiles(docsDir)
		if err != nil {
			return total, err
		}
		for _, file := range files {
			total += printFileClarifications(file, out)
		}
	}
	return total, nil
}

// yamlDocFiles lists the module's docs YAML files in sorted order.
func yamlDocFiles(docsDir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(docsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".yaml") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

// printFileClarifications prints one worklist line per open entry in one
// document and returns how many it printed.
func printFileClarifications(file string, out io.Writer) int {
	data, err := os.ReadFile(file)
	if err != nil {
		return 0
	}
	var doc clarifiedDocument
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return 0
	}
	for _, entry := range doc.NeedsClarification {
		fmt.Fprintf(out, "%s: %s %s -- %s\n", file, doc.ID, entry.ID, entry.Question)
	}
	return len(doc.NeedsClarification)
}
