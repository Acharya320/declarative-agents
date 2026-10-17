// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"sort"
	"strings"
	"testing"
)

// misDesign declares every structure a document can implement wrongly: a
// nested slot, a flat: table record kind, a code field, a linked reference
// field, a typed field, and an ordered document slot (GH-2678).
const misDesign = `
document_types:
  plan:
    location: docs/plans/name.md
    markdown:
      frontmatter: {required: [id, title], allowed: [id, title]}
      nodes:
        document_field: {tag: field, level: 2, parent: document, heading: free, content: any, fields: [overview, index, sample, traces, counts, parts]}
        overview_field: {tag: field, level: 3, parent: 'field:overview', heading: free, content: any, fields: [summary, layout]}
        index_entry:
          tag: index_entry
          level: 3
          parent: 'field:index'
          heading: id
          content: none
          flat: table
          yaml: {from: index, records: list, id: id}
        index_entry_field: {tag: field, level: 4, parent: index_entry, heading: free, content: any, fields: [title, summary]}
        req:
          tag: req
          level: 2
          parent: document
          heading: id
          content: prose
          yaml: {from: requirements, records: list, id: id, content: text}
      fields:
        sample: {in: [document_field], code: yaml}
        traces: {in: [document_field], shape: list, links: true}
`

// audit returns the names of the findings a document raises, plus
// not-canonical when it parses but differs from its canonical form, as the
// corpus audit reports them.
func audit(t *testing.T, d *Declaration, md string) []string {
	t.Helper()
	var names []string
	formatted, findings, err := d.Format([]byte(md), "docs/plans/plan.md")
	for _, f := range findings {
		names = append(names, f.Name)
	}
	if err != nil {
		names = append(names, "error: "+err.Error())
	}
	if err == nil && formatted != nil && string(formatted) != md {
		line, found, expected := FirstDifference([]byte(md), formatted)
		names = append(names, "not-canonical")
		t.Logf("not canonical at line %d: %q, canonical %q", line, found, expected)
	}
	sort.Strings(names)
	return names
}

const misHead = "---\nid: plan\ntitle: Plan\n---\n\n"

func TestAuditFindsEachWrongImplementation(t *testing.T) {
	declarations, err := LoadDeclarations([]byte(misDesign))
	if err != nil {
		t.Fatal(err)
	}
	d := declarations["plan"]
	correct := misHead +
		"<!-- req: R1 -->\n## R1\n\nKeep a budget.\n\n" +
		"<!-- field: overview -->\n## Overview\n\n<!-- field: summary -->\n### Summary\n\nText.\n\n" +
		"<!-- field: index -->\n## Index\n\n| id | title |\n| --- | --- |\n| a | A |\n\n" +
		"<!-- field: sample -->\n## Sample\n\n```yaml\nkey: value\n```\n\n" +
		"<!-- field: traces -->\n## Traces\n\n- [R1](#r1)\n\n" +
		"<!-- field: counts | form: \"typed table\" -->\n## Counts\n\n| name | n |\n| --- | --- |\n| a | 1 |\n\n" +
		"<!-- field: parts | form: \"keyed table\" -->\n## Parts\n\n| id | kind | unit |\n| --- | --- | --- |\n| p1 | gauge | ms |\n| p2 | count | n |\n"
	if got := audit(t, d, correct); len(got) != 0 {
		t.Fatalf("the correct document raised %v", got)
	}
	swap := func(old, new string) string {
		if !strings.Contains(correct, old) {
			t.Fatalf("fixture lacks %q", old)
		}
		return strings.Replace(correct, old, new, 1)
	}
	for _, c := range []struct {
		name string
		doc  string
		want []string
	}{
		{"nested slot written as a tree",
			swap("<!-- field: overview -->\n## Overview\n\n<!-- field: summary -->\n### Summary\n\nText.\n", "<!-- field: overview | form: \"tree\" -->\n## Overview\n\n- summary: Text.\n"),
			[]string{"not-canonical"}},
		{"flat record list written as a tree",
			swap("<!-- field: index -->\n## Index\n\n| id | title |\n| --- | --- |\n| a | A |\n", "<!-- field: index | form: \"tree\" -->\n## Index\n\n1. - id: a\n   - title: A\n"),
			[]string{"not-canonical"}},
		{"rich record list written as a tree",
			swap("<!-- field: index -->\n## Index\n\n| id | title |\n| --- | --- |\n| a | A |\n", "<!-- field: index | form: \"tree\" -->\n## Index\n\n1. - id: a\n   - summary:\n     One line.\n\n     Two.\n"),
			[]string{"not-canonical"}},
		{"code field written as prose",
			swap("```yaml\nkey: value\n```\n", "key: value\n"),
			[]string{"not-canonical"}}, // mage fmt fences it; a finding that stops parsing would block the fix
		{"code field written as a tree",
			swap("<!-- field: sample -->\n## Sample\n\n```yaml\nkey: value\n```\n", "<!-- field: sample | form: \"tree\" -->\n## Sample\n\n- key: value\n"),
			[]string{"attribute-shape"}},
		{"keyed table written as a tree",
			swap("<!-- field: parts | form: \"keyed table\" -->\n## Parts\n\n| id | kind | unit |\n| --- | --- | --- |\n| p1 | gauge | ms |\n| p2 | count | n |\n",
				"<!-- field: parts | form: \"tree\" -->\n## Parts\n\n- p1:\n  - kind: gauge\n  - unit: ms\n- p2:\n  - kind: count\n  - unit: n\n"),
			[]string{"not-canonical"}},
		{"generated heading misspelled",
			swap("## Overview", "## Overviews"),
			[]string{"not-canonical"}},
		{"reference link to the wrong anchor",
			swap("- [R1](#r1)", "- [R1](#r9)"),
			[]string{"not-canonical"}},
		{"reference link to a missing item",
			swap("- [R1](#r1)", "- [R7](#r7)"),
			[]string{"not-canonical"}},
		{"record heading is not its id",
			swap("<!-- field: index -->\n## Index\n\n| id | title |\n| --- | --- |\n| a | A |\n", "<!-- field: index -->\n## Index\n\n<!-- index_entry: a -->\n### b\n\n<!-- field: summary -->\n#### Summary\n\nOne.\n\nTwo.\n"),
			[]string{"heading-id-mismatch"}},
		{"field out of order",
			strings.Replace(correct, "<!-- field: traces -->\n## Traces\n\n- [R1](#r1)\n\n", "", 1) + "\n<!-- field: traces -->\n## Traces\n\n- [R1](#r1)\n",
			[]string{"field-out-of-order"}},
		{"undeclared field",
			swap("<!-- field: counts", "<!-- field: stranger -->\n## Stranger\n\nOdd.\n\n<!-- field: counts"),
			[]string{"field-undeclared"}},
		{"undeclared nested field",
			swap("### Summary\n\nText.\n", "### Summary\n\nText.\n\n<!-- field: extra -->\n### Extra\n\nMore.\n"),
			[]string{"field-undeclared"}},
	} {
		t.Run(c.name, func(t *testing.T) {
			got := audit(t, d, c.doc)
			if strings.Join(got, ",") != strings.Join(c.want, ",") {
				t.Fatalf("findings = %v, want %v\n%s", got, c.want, c.doc)
			}
		})
	}
}

func TestFirstDifferenceNamesTheLine(t *testing.T) {
	line, found, expected := FirstDifference([]byte("a\nb\nc\n"), []byte("a\nB\nc\n"))
	if line != 2 || found != "b" || expected != "B" {
		t.Fatalf("got %d %q %q", line, found, expected)
	}
	if line, _, _ := FirstDifference([]byte("same\n"), []byte("same\n")); line != 0 {
		t.Fatalf("equal documents differ at %d", line)
	}
}
