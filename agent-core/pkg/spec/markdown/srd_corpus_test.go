// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec"
	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec/markdown"
	"gopkg.in/yaml.v3"
)

const (
	designPath      = "../../../docs/constitutions/design.yaml"
	repositoryRoot  = "../../../.."
	fixtureYAML     = "testdata/srd004-declarative-ux.yaml"
	fixtureMarkdown = "testdata/srd004-declarative-ux.md"
	fixtureName     = "docs/specs/software-requirements/srd004-declarative-ux.md"
)

func srdDeclaration(t *testing.T) *markdown.Declaration {
	t.Helper()
	data, err := os.ReadFile(designPath)
	if err != nil {
		t.Fatal(err)
	}
	declarations, err := markdown.LoadDeclarations(data)
	if err != nil {
		t.Fatal(err)
	}
	return declarations["srd"]
}

func readYAML(t *testing.T, path string) *yaml.Node {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatal(err)
	}
	return &doc
}

// trackedSRDs lists every tracked SRD in the repository, testdata excluded.
func trackedSRDs(t *testing.T) []string {
	t.Helper()
	out, err := exec.Command("git", "-C", repositoryRoot, "ls-files").Output()
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	for _, path := range strings.Split(string(out), "\n") {
		if strings.Contains(path, "/docs/specs/software-requirements/srd") &&
			strings.HasSuffix(path, ".yaml") && !strings.Contains(path, "/testdata/") {
			paths = append(paths, path)
		}
	}
	if len(paths) == 0 {
		t.Fatal("no tracked SRDs found")
	}
	return paths
}

// TestEveryCorpusSRDRoundTrips states the round-trip laws on every tracked
// SRD (eng10): the canonical YAML tree survives print and parse unchanged,
// the printed Markdown is a fixed point of Format, the only findings are
// prose containment, and ParseSRD reads the same model from either form.
func TestEveryCorpusSRDRoundTrips(t *testing.T) {
	d := srdDeclaration(t)
	for _, path := range trackedSRDs(t) {
		t.Run(path, func(t *testing.T) {
			srdRoundTrip(t, d, filepath.Join(repositoryRoot, path), path)
		})
	}
}

func srdRoundTrip(t *testing.T, d *markdown.Declaration, file, path string) {
	canonical, _ := d.Canonicalize(readYAML(t, file))
	name := strings.TrimSuffix(path, ".yaml") + ".md"
	md, err := d.Print(canonical)
	if err != nil {
		t.Fatal(err)
	}
	back, findings, err := d.Parse(md, name)
	if err != nil || back == nil {
		t.Fatalf("parse: %v %v", err, findings)
	}
	for _, f := range findings {
		if f.Name != "prose-containment" {
			t.Errorf("printed form has %s", f)
		}
	}
	if diff := markdown.Diff(canonical, back); diff != "" {
		t.Errorf("yaml -> md -> yaml changed %s", diff)
	}
	if formatted, _, err := d.Format(md, name); err != nil || string(formatted) != string(md) {
		t.Errorf("printed Markdown is not a fixed point: %v", err)
	}
	original, err := spec.ParseSRD(file)
	if err != nil {
		t.Fatal(err)
	}
	fromMarkdown := parseSRDNode(t, back, filepath.Base(file))
	if diffs := modelDiffs(&original, &fromMarkdown, true); len(diffs) > 0 {
		t.Errorf("model changed: %s", diffs[0])
	}
}

// parseSRDNode runs ParseSRD over a tree by writing it to a file.
func parseSRDNode(t *testing.T, doc *yaml.Node, base string) spec.SRD {
	t.Helper()
	data, err := yaml.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), base)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	srd, err := spec.ParseSRD(path)
	if err != nil {
		t.Fatal(err)
	}
	return srd
}

// TestFixtureMatchesTheConverter pins the canonical Markdown of srd004.
func TestFixtureMatchesTheConverter(t *testing.T) {
	d := srdDeclaration(t)
	canonical, _ := d.Canonicalize(readYAML(t, fixtureYAML))
	md, err := d.Print(canonical)
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(fixtureMarkdown)
	if err != nil {
		t.Fatal(err)
	}
	if string(md) != string(want) {
		t.Fatalf("printed srd004 differs from %s; regenerate it and review the diff", fixtureMarkdown)
	}
}
