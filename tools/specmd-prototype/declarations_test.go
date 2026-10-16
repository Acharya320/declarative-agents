// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package specmd

import (
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestEveryDocumentTypeDeclaresMarkdown holds design.yaml to its own
// registry: each document type carries a markdown entry, either a rule
// declaration that compiles or a note saying why the type takes none.
func TestEveryDocumentTypeDeclaresMarkdown(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	path, err := findDesign(wd)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var design struct {
		DocumentTypes map[string]map[string]yaml.Node `yaml:"document_types"`
	}
	if err := yaml.Unmarshal(data, &design); err != nil {
		t.Fatal(err)
	}
	declarations, err := loadDeclarations(data)
	if err != nil {
		t.Fatal(err)
	}
	for name, entry := range design.DocumentTypes {
		markdown, ok := entry["markdown"]
		switch {
		case !ok:
			t.Errorf("document_types.%s has no markdown entry", name)
		case markdown.Kind == yaml.ScalarNode && markdown.Value == "":
			t.Errorf("document_types.%s.markdown is empty", name)
		case markdown.Kind == yaml.MappingNode && declarations[name] == nil:
			t.Errorf("document_types.%s.markdown did not load", name)
		}
	}
}

// TestDeclaredLocationsPickTheirType checks that each declared Markdown
// location selects its own type and no other: the singleton names whose
// capitals hold an N that is not a digit placeholder, and the types that
// share the name placeholder and differ only by directory.
func TestDeclaredLocationsPickTheirType(t *testing.T) {
	declarations, err := declarationsFromRepository()
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string]string{
		"docs/VISION.md":                                        "vision",
		"docs/ARCHITECTURE.md":                                  "architecture",
		"docs/SPECIFICATIONS.md":                                "specifications",
		"docs/road-map.md":                                      "roadmap",
		"applications/docs/pattern-language.md":                 "pattern_language",
		"docs/engineering/eng10-markdown-spec-format.md":        "engineering_guideline",
		"docs/specs/software-requirements/srd004-ux.md":         "srd",
		"docs/specs/use-cases/rel01.0-uc003-bench.md":           "use_case",
		"docs/specs/test-suites/test-rel02.0-lifecycle.md":      "test_suite",
		"docs/specs/audits/structure-audit-2026-08-register.md": "audit_register",
		"docs/specs/semantic-models/role-ontology.md":           "semantic_model",
		"docs/specs/config-formats/machine-format.md":           "config_format",
		"agent-core/docs/constitutions/design.md":               "constitution",
		"docs/migrations/runtime-split.md":                      "migration",
	}
	for file, want := range cases {
		r, ok := declarationFor(declarations, file)
		if !ok || r.documentType != want {
			got := "none"
			if ok {
				got = r.documentType
			}
			t.Errorf("%s: got %s, want %s", file, got, want)
		}
	}
}

// TestEmptyListAndUndeclaredColumns covers the two shape relaxations every
// type relies on: a list field with no entries is still a list, and a table
// field that declares no columns accepts any header.
func TestEmptyListAndUndeclaredColumns(t *testing.T) {
	empty := strings.Replace(validSRD, "- Rewriting the shell.\n", "", 1)
	if findings := Diagnose([]byte(empty), validName); len(findings) > 0 {
		t.Fatalf("empty non_goals list reported findings:\n%v", findings)
	}
	index := "---\nid: sample-specifications\ntitle: Sample\n---\n"
	for _, key := range []string{"overview", "roadmap_summary", "foundation_document_index", "srd_index",
		"config_format_index", "semantic_model_index", "use_case_index", "test_suite_index", "coverage_gaps"} {
		index += "\n<!-- field: " + key + " -->\n## " + key + "\n\n"
		switch key {
		case "overview", "coverage_gaps":
			index += "Text.\n"
		case "srd_index":
			index += "| id | title | any column |\n| --- | --- | --- |\n| srd001 | One | x |\n"
		}
	}
	if findings := Diagnose([]byte(index), "docs/SPECIFICATIONS.md"); len(findings) > 0 {
		t.Fatalf("index with undeclared columns and empty tables reported findings:\n%v", findings)
	}
}
