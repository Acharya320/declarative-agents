// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec/markdown"
	"gopkg.in/yaml.v3"
)

// TestConvertingTheCorpusAndBackKeepsEveryTree converts every tracked corpus
// document to Markdown and back to YAML, as mage convert and its inverse do,
// and requires each tree to come back equal to its canonical form; what
// differs from the original text is in the conversion report (#2582 AC2).
func TestConvertingTheCorpusAndBackKeepsEveryTree(t *testing.T) {
	declarations := repositoryDeclarations(t)
	out, err := exec.Command("git", "-C", repositoryRoot, "ls-files", "*.yaml").Output()
	if err != nil {
		t.Fatal(err)
	}
	converted := 0
	for _, path := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		d, ok := declarations.ForSource(path)
		if !ok || strings.Contains(path, "/testdata/") {
			continue
		}
		converted++
		data, err := os.ReadFile(filepath.Join(repositoryRoot, path))
		if err != nil {
			t.Fatal(err)
		}
		md, report, err := d.ConvertYAML(data, path)
		if err != nil {
			t.Errorf("%s: %v", path, err)
			continue
		}
		back, _, err := d.ConvertMarkdown(md, report.Target)
		if err != nil {
			t.Errorf("%s: rollback: %v", path, err)
			continue
		}
		var original, rolled yaml.Node
		if err := yaml.Unmarshal(data, &original); err != nil {
			t.Fatal(err)
		}
		if err := yaml.Unmarshal(back, &rolled); err != nil {
			t.Fatalf("%s: rollback is not YAML: %v", path, err)
		}
		canonical, _ := d.Canonicalize(&original)
		if diff := markdown.Diff(canonical, &rolled); diff != "" {
			t.Errorf("%s: convert and back changed %s", path, diff)
		}
	}
	if converted < 368 {
		t.Errorf("converted %d documents, want every one of the 368 converted", converted)
	}
}
