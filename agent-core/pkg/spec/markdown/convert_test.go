// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestConvertReportsWhatItChanges converts a recipe whose YAML carries a
// note comment, a literal block, and a record in non-canonical order, and
// checks that the report names each and the rollback restores the tree.
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
	back, _, err := d.ConvertMarkdown(md, report.Target)
	if err != nil {
		t.Fatal(err)
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(back, &doc); err != nil {
		t.Fatal(err)
	}
	canonical, _ := d.Canonicalize(yamlDoc(t, source))
	if diff := Diff(canonical, &doc); diff != "" {
		t.Fatalf("rollback changed %s", diff)
	}
	if !strings.HasPrefix(string(back), "# Copyright (c) 2026 Nokia\n# SPDX-License-Identifier: BSD-3-Clause\n") {
		t.Errorf("rollback lost the license header:\n%s", back)
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
		if got := d.SourcePath(tc.markdown); got != tc.source {
			t.Errorf("%s: SourcePath = %s, want %s", tc.kind, got, tc.source)
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
