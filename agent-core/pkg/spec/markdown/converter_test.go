// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// toyDesign declares a document type no repository document uses. It exists
// only here, so a round trip through it proves a new type is declaration
// only: the converter has no branch for it (#2582 AC3).
const toyDesign = `
document_types:
  recipe:
    location: docs/recipes/name.yaml
    markdown:
      location: docs/recipes/name.md
      frontmatter: {required: [id, title], allowed: [id, title, status], id_matches_filename: true}
      nodes:
        document_field: {tag: field, level: 2, parent: document, heading: free, content: any}
        course:
          tag: course
          level: 2
          parent: document
          id_pattern: '^C[0-9]+$'
          heading: id_title
          content: none
          unique_id: true
          yaml: {from: courses, records: map, title: title}
        course_field: {tag: field, level: 3, parent: course, heading: free, content: any}
        dish:
          tag: dish
          level: 3
          parent: course
          id_pattern: '^C[0-9]+\.[0-9]+$'
          id_extends_parent: true
          heading: id
          content: prose
          unique_id: true
          attributes: [minutes, tags]
          yaml: {from: dishes, records: list, id: id, content: text, keyed: true}
        dish_field: {tag: field, level: 4, parent: dish, heading: free, content: any}
        step:
          tag: step
          level: 3
          parent: 'field:steps'
          id_pattern: '^S[0-9]+$'
          heading: id
          content: prose
          unique_id: true
          required_fields: [checks]
          yaml: {from: steps, records: list, id: id, content: text}
        step_field: {tag: field, level: 4, parent: step, heading: free, content: any}
      fields:
        description: {in: [document_field], shape: prose}
        ingredients: {in: [document_field], shape: list, value: keyed}
        tools: {in: [document_field], shape: list}
        timeline: {in: [document_field], shape: list, value: outline, id: id, content: text, item_pattern: '^T[0-9]+: '}
        suppliers: {in: [document_field], shape: table, columns: [name, city]}
        steps: {in: [document_field], shape: empty}
        checks: {in: [step_field], shape: list}
      required_fields: [description, steps]
      attributes: {minutes: integer, tags: list}
`

// toyRecipe is a recipe in canonical YAML: each record's keys run id,
// title, content, attributes, then children.
const toyRecipe = `# Copyright (c) 2026 Nokia
# SPDX-License-Identifier: BSD-3-Clause
id: soup
title: Soup
status: draft
description: |
  A warm soup.

  Two paragraphs.
ingredients:
  water: 1 litre
  salt: a pinch
tools:
  - pot
  - spoon
timeline:
  - T1: Boil the water.
  - id: T2
    text: Add the salt.
    owner: cook
    checks:
      - taste
      - smell
suppliers:
  - name: Market
    city: Paris
courses:
  C1:
    title: Starter
    dishes:
      - C1.1: Bread.
      - id: C1.2
        text: Soup itself.
        minutes: 20
        tags: [warm, "1"]
        garnish: parsley
    serving: Hot.
steps:
  - id: S1
    text: Heat.
    checks:
      - hot
notes:
  temperature: {min: 80, max: 95}
  history: null
`

const toyName = "docs/recipes/soup.md"

func toyDeclaration(t *testing.T) *Declaration {
	t.Helper()
	declarations, err := LoadDeclarations([]byte(toyDesign))
	if err != nil {
		t.Fatal(err)
	}
	return declarations["recipe"]
}

func yamlDoc(t *testing.T, text string) *yaml.Node {
	t.Helper()
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(text), &doc); err != nil {
		t.Fatal(err)
	}
	return &doc
}

// roundTrip prints a YAML document, parses the Markdown back, and returns
// both; the Markdown must diagnose clean.
func roundTrip(t *testing.T, d *Declaration, doc *yaml.Node) (string, *yaml.Node) {
	t.Helper()
	md, err := d.Print(doc)
	if err != nil {
		t.Fatal(err)
	}
	back, findings, err := d.Parse(md, toyName)
	if err != nil || len(findings) > 0 {
		t.Fatalf("parse: %v %v\n%s", err, findings, md)
	}
	return string(md), back
}

func TestToyTypeRoundTrips(t *testing.T) {
	d := toyDeclaration(t)
	original := yamlDoc(t, toyRecipe)
	md, back := roundTrip(t, d, original)
	if diff := Diff(original, back); diff != "" {
		t.Fatalf("yaml -> md -> yaml changed %s\n%s", diff, md)
	}
	if !strings.Contains(back.HeadComment, "SPDX-License-Identifier: BSD-3-Clause") {
		t.Errorf("license lines lost: %q", back.HeadComment)
	}
}

func TestToyMarkdownIsAFixedPoint(t *testing.T) {
	d := toyDeclaration(t)
	md, _ := roundTrip(t, d, yamlDoc(t, toyRecipe))
	formatted, findings, err := d.Format([]byte(md), toyName)
	if err != nil || len(findings) > 0 {
		t.Fatalf("format: %v %v", err, findings)
	}
	if string(formatted) != md {
		t.Fatalf("md -> yaml -> md is not the identity:\n--- printed\n%s\n--- formatted\n%s", md, formatted)
	}
	again, _, err := d.Format(formatted, toyName)
	if err != nil || string(again) != string(formatted) {
		t.Fatalf("format is not idempotent: %v", err)
	}
}

// TestEachConversionPrintsItsShape covers each value conversion and each
// node source (eng10 Tables 4 and 5): the Markdown carries the expected
// form, and the value comes back unchanged.
func TestEachConversionPrintsItsShape(t *testing.T) {
	base := "id: soup\ntitle: Soup\ndescription: Text.\nsteps:\n  - id: S1\n    text: Heat.\n    checks: [hot]\n"
	cases := []struct{ name, member, want string }{
		{"string as prose", "", "Text."},
		{"string list", "tools: [pot, spoon]\n", "- pot\n- spoon"},
		{"keyed list", "ingredients:\n  water: 1 litre\n", "- water: 1 litre"},
		{"outline keyed item", "timeline:\n  - T1: Boil.\n", "- T1: Boil."},
		{"outline nested items", "timeline:\n  - id: T1\n    text: Boil.\n    owner: cook\n    checks: [a, b]\n",
			"- T1: Boil.\n  - owner: cook\n  - checks:\n    - a\n    - b"},
		{"record table", "suppliers:\n  - {name: Market, city: Paris}\n", "| name | city |"},
		{"fenced yaml", "notes:\n  deep: {a: 1}\n", "```yaml"},
		{"list records under a field", "", "<!-- step: S1 -->\n### S1"},
		{"map records", "courses:\n  C1:\n    title: Starter\n", "<!-- course: C1 -->\n## C1 - Starter"},
		{"keyed records", "courses:\n  C1:\n    dishes:\n      - C1.1: Bread.\n", "<!-- dish: C1.1 -->\n### C1.1\n\nBread."},
		{"record attributes", "courses:\n  C1:\n    dishes:\n      - {id: C1.1, text: Bread., minutes: 5, tags: [a]}\n",
			`<!-- dish: C1.1 | minutes: 5 | tags: ["a"] -->`},
	}
	d := toyDeclaration(t)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			original := yamlDoc(t, base+tc.member)
			md, back := roundTrip(t, d, original)
			if !strings.Contains(md, tc.want) {
				t.Errorf("Markdown lacks %q:\n%s", tc.want, md)
			}
			if diff := Diff(original, back); diff != "" {
				t.Errorf("round trip changed %s", diff)
			}
		})
	}
}

func TestPrintRefusesValuesTheDeclarationCannotHold(t *testing.T) {
	base := "id: soup\ntitle: Soup\ndescription: Text.\nsteps: []\n"
	cases := map[string]string{
		"keyed field holding a list":    "ingredients: [water]\n",
		"outline item with a deep map":  "timeline:\n  - {id: T1, text: Boil., deep: {a: 1}}\n",
		"list field holding a map":      "tools: {a: b}\n",
		"record without its id":         "steps:\n  - {text: Heat.}\n",
		"map records written as a list": "courses: [a]\n",
	}
	d := toyDeclaration(t)
	for name, member := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := d.Print(yamlDoc(t, base+member)); err == nil {
				t.Errorf("printed a value the declaration cannot hold")
			}
		})
	}
}

func TestParseReturnsFindingsInsteadOfATree(t *testing.T) {
	d := toyDeclaration(t)
	md, _ := roundTrip(t, d, yamlDoc(t, toyRecipe))
	broken := strings.Replace(md, "<!-- step: S1 -->", "<!-- step: X1 -->", 1)
	doc, findings, err := d.Parse([]byte(broken), toyName)
	if err != nil || doc != nil || len(findings) == 0 || findings[0].Name != "id-malformed" {
		t.Fatalf("got doc=%v findings=%v err=%v", doc != nil, findings, err)
	}
}

func TestMappingDeclarationsAreChecked(t *testing.T) {
	cases := map[string]string{
		"records neither list nor map":  "yaml: {from: steps, records: set, id: id}",
		"list records without an id":    "yaml: {from: steps, records: list}",
		"keyed records without content": "yaml: {from: steps, records: list, id: id, keyed: true}",
	}
	for name, mapping := range cases {
		t.Run(name, func(t *testing.T) {
			design := strings.Replace(toyDesign, "yaml: {from: steps, records: list, id: id, content: text}", mapping, 1)
			if _, err := LoadDeclarations([]byte(design)); err == nil {
				t.Errorf("accepted %s", mapping)
			}
		})
	}
	outline := strings.Replace(toyDesign, "value: outline, id: id, content: text,", "value: outline,", 1)
	if _, err := LoadDeclarations([]byte(outline)); err == nil {
		t.Error("accepted an outline without its id and content keys")
	}
}
