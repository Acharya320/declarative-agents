// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec/markdown"
	"gopkg.in/yaml.v3"
)

// SpecYAML prints Markdown corpus documents as the YAML trees their
// design.yaml declarations parse them into (eng10), for a script that checks
// document fields: mage specYAML applications/docs/SPECIFICATIONS.md prints
// one tree, and a directory prints one mapping from each declared document's
// path to its tree. Paths are relative to the repository root, testdata is
// skipped, and a document with a format finding fails the target (GH-2627).
func SpecYAML(path string) error {
	root, err := repositoryRoot()
	if err != nil {
		return err
	}
	out, err := specYAML(root, filepath.Clean(path))
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(out)
	return err
}

func specYAML(root, path string) ([]byte, error) {
	full := filepath.Join(root, path)
	info, err := os.Stat(full)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return markdown.ReadAsYAML(full)
	}
	declarations, err := markdown.Default()
	if err != nil {
		return nil, err
	}
	trees := &yaml.Node{Kind: yaml.MappingNode}
	err = filepath.WalkDir(full, func(file string, entry fs.DirEntry, err error) error {
		switch {
		case err != nil:
			return err
		case entry.IsDir() && entry.Name() == "testdata":
			return filepath.SkipDir
		case entry.IsDir() || !strings.HasSuffix(file, ".md"):
			return nil
		}
		if _, ok := declarations.ForMarkdown(filepath.ToSlash(file)); !ok {
			return nil // native Markdown, such as a guide
		}
		return appendTree(trees, root, file)
	})
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(trees)
}

// appendTree adds one document's tree to trees under its repository path.
func appendTree(trees *yaml.Node, root, file string) error {
	data, err := markdown.ReadAsYAML(file)
	if err != nil {
		return err
	}
	var tree yaml.Node
	if err := yaml.Unmarshal(data, &tree); err != nil {
		return fmt.Errorf("%s: %w", file, err)
	}
	rel, err := filepath.Rel(root, file)
	if err != nil {
		return err
	}
	key := &yaml.Node{Kind: yaml.ScalarNode, Value: filepath.ToSlash(rel)}
	trees.Content = append(trees.Content, key, tree.Content[0])
	return nil
}
