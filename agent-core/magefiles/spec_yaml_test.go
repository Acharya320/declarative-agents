// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"testing"

	"gopkg.in/yaml.v3"
)

// TestSpecYAMLPrintsADocumentOrADirectory reads one corpus document's tree
// and a directory's mapping of trees, which leaves out native Markdown.
func TestSpecYAMLPrintsADocumentOrADirectory(t *testing.T) {
	root, err := repositoryRoot()
	if err != nil {
		t.Fatal(err)
	}
	data, err := specYAML(root, "applications/docs/road-map.md")
	if err != nil {
		t.Fatal(err)
	}
	var roadmap struct {
		ID       string `yaml:"id"`
		Releases []any  `yaml:"releases"`
	}
	if err := yaml.Unmarshal(data, &roadmap); err != nil || roadmap.ID == "" || len(roadmap.Releases) == 0 {
		t.Fatalf("road map tree = %+v, %v", roadmap, err)
	}
	data, err = specYAML(root, "applications/docs")
	if err != nil {
		t.Fatal(err)
	}
	var trees map[string]map[string]any
	if err := yaml.Unmarshal(data, &trees); err != nil {
		t.Fatal(err)
	}
	if trees["applications/docs/SPECIFICATIONS.md"]["id"] == nil {
		t.Errorf("directory mapping lacks the specification index: %d documents", len(trees))
	}
	if _, ok := trees["applications/docs/specs/test-suites/test-rel14.0-application-consistency.md"]; !ok {
		t.Errorf("directory mapping lacks a test suite: %d documents", len(trees))
	}
	if _, err := specYAML(root, "applications/docs/road-map.yaml"); err == nil {
		t.Error("read a corpus document by its former YAML path")
	}
}
