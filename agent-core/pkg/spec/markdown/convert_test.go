// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestConvertReportsWhatItChanges converts a recipe whose YAML carries a
// note comment, a literal block, and a record in non-canonical order, and
// checks that the report names each and the Markdown parses back to the
// canonical tree.
func TestConvertReportsWhatItChanges(t *testing.T) {
	d := toyDeclaration(t)
	source := strings.Replace(toyRecipe, "tools:\n", "# Tools the cook needs.\ntools:\n", 1)
	source = strings.Replace(source, "  - id: S1\n    text: Heat.\n", "  - text: Heat.\n    id: S1\n", 1)
	md, report, err := d.ConvertYAML([]byte(source), "docs/recipes/soup.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if report.Target != "docs/recipes/soup.md" {
		t.Errorf("target %s", report.Target)
	}
	want := map[string][]string{
		"comment dropped": report.Comments, "canonical": report.Changes, "whitespace": report.Whitespace,
	}
	for label, entries := range want {
		if len(entries) == 0 {
			t.Errorf("report has no %s entry:\n%s", label, strings.Join(report.Lines(), "\n"))
		}
	}
	back, _, err := d.Parse(md, report.Target)
	if err != nil || back == nil {
		t.Fatalf("converted Markdown does not parse: %v", err)
	}
	canonical, _ := d.Canonicalize(yamlDoc(t, source))
	if diff := Diff(canonical, back); diff != "" {
		t.Fatalf("conversion changed %s", diff)
	}
}

func TestConvertRefusesWhatItCannotCarry(t *testing.T) {
	d := toyDeclaration(t)
	source := "id: soup\ntitle: Soup\ndescription: Text.\nsteps: []\ningredients: [water]\n"
	if _, _, err := d.ConvertYAML([]byte(source), "docs/recipes/soup.yaml"); err == nil {
		t.Fatal("converted a keyed field holding a list")
	}
}

func TestPathsFollowTheDeclaredLocations(t *testing.T) {
	declarations := repositoryDeclarations(t)
	cases := []struct{ kind, source, markdown string }{
		{"srd", "docs/specs/software-requirements/srd004-declarative-ux.yaml", "docs/specs/software-requirements/srd004-declarative-ux.md"},
		{"audit_register", "docs/specs/audits/structure-audit-2026-08-19.yaml", "docs/specs/audits/structure-audit-2026-08-19-register.md"},
		{"vision", "docs/VISION.yaml", "docs/VISION.md"},
	}
	for _, tc := range cases {
		d := declarations[tc.kind]
		if got := d.MarkdownPath(tc.source); got != tc.markdown {
			t.Errorf("%s: MarkdownPath = %s, want %s", tc.kind, got, tc.markdown)
		}
		if found, ok := declarations.ForSource(tc.source); !ok || found != d {
			t.Errorf("%s: ForSource(%s) found the wrong type", tc.kind, tc.source)
		}
		if found, ok := declarations.ForMarkdown(tc.markdown); !ok || found != d {
			t.Errorf("%s: ForMarkdown(%s) found the wrong type", tc.kind, tc.markdown)
		}
	}
}

// TestDesignConstitutionStaysYAML keeps the registry every loader reads out of
// conversion, while the other constitutions convert.
func TestDesignConstitutionStaysYAML(t *testing.T) {
	declarations := repositoryDeclarations(t)
	if _, ok := declarations.ForSource("agent-core/docs/constitutions/design.yaml"); ok {
		t.Error("design.yaml is in conversion scope")
	}
	if d, ok := declarations.ForSource("agent-core/docs/constitutions/testing.yaml"); !ok || d.Type() != "constitution" {
		t.Error("testing.yaml is not a constitution in scope")
	}
}

// TestConvertTreeAndReadAsYAML converts a small tree in place, leaving
// testdata alone, and reads the converted document by its Markdown path; its
// former YAML path no longer reads (GH-2626).
func TestConvertTreeAndReadAsYAML(t *testing.T) {
	root := t.TempDir()
	guideline := "id: eng90-sample\ntitle: Sample\nintroduction: We describe a practice.\n" +
		"sections:\n  - title: The practice\n    content: Each module follows it.\nreferences:\n  - https://example.com\n"
	write := func(rel, body string) {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("docs/engineering/eng90-sample.yaml", guideline)
	write("testdata/docs/engineering/eng91-kept.yaml", guideline)
	reports, err := repositoryDeclarations(t).ConvertTree(root)
	if err != nil || len(reports) != 1 {
		t.Fatalf("ConvertTree: %v (%d reports)", err, len(reports))
	}
	if _, err := os.Stat(filepath.Join(root, "testdata/docs/engineering/eng91-kept.yaml")); err != nil {
		t.Fatalf("testdata was converted: %v", err)
	}
	data, err := ReadAsYAML(filepath.Join(root, "docs/engineering/eng90-sample.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"id: eng90-sample", "introduction: We describe a practice.", "title: The practice"} {
		if !strings.Contains(string(data), want) {
			t.Errorf("ReadAsYAML lacks %q:\n%s", want, data)
		}
	}
	if _, err := ReadAsYAML(filepath.Join(root, "docs/engineering/eng90-sample.yaml")); err == nil {
		t.Error("read a corpus document by its former YAML path")
	}
	if _, err := ReadAsYAML(filepath.Join(root, "docs/engineering/eng92-absent.md")); err == nil {
		t.Error("read a document that does not exist")
	}
}
