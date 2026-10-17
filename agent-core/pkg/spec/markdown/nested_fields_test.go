// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// nestedDesign declares the eng10 Table 4 example: nested field slots under
// overview and layout, records under a nested field, a record field slot at
// level 6, and keyed records under a nested field (GH-2653).
const nestedDesign = `
document_types:
  plan:
    location: docs/plans/name.md
    markdown:
      frontmatter: {required: [id, title], allowed: [id, title]}
      nodes:
        document_field: {tag: field, level: 2, parent: document, heading: free, content: any}
        overview_field: {tag: field, level: 3, parent: 'field:overview', heading: free, content: any}
        layout_field: {tag: field, level: 4, parent: 'field:layout', heading: free, content: any}
        package:
          tag: package
          level: 5
          parent: 'field:packages'
          heading: id
          content: none
          yaml: {from: packages, records: list, id: path}
        package_field: {tag: field, level: 6, parent: package, heading: free, content: any}
        criterion:
          tag: criterion
          level: 4
          parent: 'field:criteria'
          heading: id
          content: prose
          yaml: {from: criteria, records: map, content: text}
        index_entry:
          tag: index_entry
          level: 3
          parent: 'field:index'
          heading: id
          content: none
          flat: table
          yaml: {from: index, records: list, id: id}
        index_entry_field: {tag: field, level: 4, parent: index_entry, heading: free, content: any}
      fields:
        index: {in: [document_field], link_columns: [id]}
`

const nestedPlan = `id: plan
title: Plan
overview:
  summary: |
    Agent Core is a framework.

    It has layers.
  layer_diagram: |
    ` + "```mermaid" + `
    flowchart TB
      A --> B
    ` + "```" + `
  layout:
    status: current
    packages:
      - path: internal/tools/catalog/
        responsibility: Catalog owns ToolDef.
        dependencies:
          may_depend_on: [internal/runtime/core]
          must_not_depend_on: [internal/tools/exec]
      - path: internal/tools/exec/
        responsibility: Exec owns ExecBuilder.
    words: [read, edit]
  criteria:
    clear_layout:
      text: Each package has one owner.
    closed_flows:
      text: Every flow ends.
  limits: {}
notes: Plain text.
`

func nestedDeclaration(t *testing.T) *Declaration {
	t.Helper()
	declarations, err := LoadDeclarations([]byte(nestedDesign))
	if err != nil {
		t.Fatal(err)
	}
	return declarations["plan"]
}

// TestNestedFieldsPrintAsHeadings round-trips a two-level nested field with
// prose, a Mermaid diagram, a scalar, a list, records under a nested field
// with a level-6 tree field, keyed records under a nested field, and an
// empty map, and checks each prints under its own heading.
func TestNestedFieldsPrintAsHeadings(t *testing.T) {
	d := nestedDeclaration(t)
	original := yamlDoc(t, nestedPlan)
	md, back := nestedRoundTrip(t, d, original)
	if diff := Diff(original, back); diff != "" {
		t.Fatalf("round trip changed %s\n%s", diff, md)
	}
	for _, want := range []string{
		"<!-- field: overview -->\n## Overview\n\n<!-- field: summary -->\n### Summary\n\nAgent Core is a framework.\n\nIt has layers.\n",
		"<!-- field: layer_diagram -->\n### Layer diagram\n\n```mermaid\n",
		"<!-- field: layout -->\n### Layout\n\n<!-- field: status -->\n#### Status\n\ncurrent\n",
		"<!-- field: packages -->\n#### Packages\n\n<!-- package: internal/tools/catalog/ -->\n##### internal/tools/catalog/\n\n<!-- field: responsibility -->\n###### Responsibility\n\nCatalog owns ToolDef.\n",
		"<!-- field: dependencies | form: \"tree\" -->\n###### Dependencies\n\n- may_depend_on:\n  1. internal/runtime/core\n",
		"<!-- field: words -->\n#### Words\n\n- read\n- edit\n",
		"<!-- criterion: clear_layout -->\n#### clear_layout\n\nEach package has one owner.\n",
		"<!-- field: notes -->\n## Notes\n\nPlain text.\n",
	} {
		if !strings.Contains(md, want) {
			t.Errorf("printed Markdown lacks %q\n%s", want, md)
		}
	}
}

// TestNestedSlotKeepsOtherValues prints a field whose slot declares nested
// fields, but whose value is text, as text: the slot applies to maps only.
func TestNestedSlotKeepsOtherValues(t *testing.T) {
	d := nestedDeclaration(t)
	original := yamlDoc(t, "id: plan\ntitle: Plan\noverview: Just a summary.\n")
	md, back := nestedRoundTrip(t, d, original)
	if diff := Diff(original, back); diff != "" {
		t.Fatalf("round trip changed %s\n%s", diff, md)
	}
	if !strings.Contains(md, "## Overview\n\nJust a summary.\n") {
		t.Errorf("text overview did not print as text\n%s", md)
	}
}

// TestNestedFieldTextBeforeChildrenFails refuses text between a field's
// heading and its first nested field: it belongs to no key.
func TestNestedFieldTextBeforeChildrenFails(t *testing.T) {
	d := nestedDeclaration(t)
	md := "---\nid: plan\ntitle: Plan\n---\n\n<!-- field: overview -->\n## Overview\n\nStray text.\n\n" +
		"<!-- field: summary -->\n### Summary\n\nText.\n"
	if _, _, err := d.Parse([]byte(md), "docs/plans/plan.md"); err == nil || !strings.Contains(err.Error(), "nested field") {
		t.Fatalf("Parse error = %v, want text-before-nested-fields error", err)
	}
}

func nestedRoundTrip(t *testing.T, d *Declaration, doc *yaml.Node) (string, *yaml.Node) {
	t.Helper()
	md, err := d.Print(doc)
	if err != nil {
		t.Fatal(err)
	}
	back, findings, err := d.Parse(md, "docs/plans/plan.md")
	if err != nil || len(findings) > 0 {
		t.Fatalf("parse: %v %v\n%s", err, findings, md)
	}
	return string(md), back
}

// TestNestedRecordListsPrintAsTables prints a list of flat records under a
// nested field heading as a table, typed when a value is not a string, and
// an empty field as None, all reading back exactly (GH-2655).
func TestNestedRecordListsPrintAsTables(t *testing.T) {
	d := nestedDeclaration(t)
	original := yamlDoc(t, `id: plan
title: Plan
overview:
  workspace_state:
    - command: write
      strategy: file_snapshot_restore
    - command: copy_dir
      strategy: tool_owned_commit_restore
  counts:
    - command: write
      calls: 3
  mismatches: []
  limits: {}
  word: None
`)
	md, back := nestedRoundTrip(t, d, original)
	if diff := Diff(original, back); diff != "" {
		t.Fatalf("round trip changed %s\n%s", diff, md)
	}
	for _, want := range []string{
		"### Workspace state\n\n| command | strategy |\n| --- | --- |\n| write | file_snapshot_restore |\n",
		"<!-- field: counts | form: \"typed table\" -->\n### Counts\n\n| command | calls |\n| --- | --- |\n| write | 3 |\n",
		"<!-- field: mismatches | form: \"empty list\" -->\n### Mismatches\n\nNone\n",
		"<!-- field: limits | form: \"empty map\" -->\n### Limits\n\nNone\n",
		"<!-- field: word -->\n### Word\n\nNone\n",
	} {
		if !strings.Contains(md, want) {
			t.Errorf("printed Markdown lacks %q\n%s", want, md)
		}
	}
}

// TestRecordFieldKeepsOtherValues prints a field whose level declares a
// record kind, but whose value is text or a map, in its usual form and
// reads it back, and finds no record kind for the same key at another level
// (GH-2663).
func TestRecordFieldKeepsOtherValues(t *testing.T) {
	d := nestedDeclaration(t)
	original := yamlDoc(t, `id: plan
title: Plan
overview:
  packages: Packages are listed elsewhere.
  layout:
    status: current
    packages: To be decided.
notes: Plain text.
`)
	md, back := nestedRoundTrip(t, d, original)
	if diff := Diff(original, back); diff != "" {
		t.Fatalf("round trip changed %s\n%s", diff, md)
	}
	for _, want := range []string{
		"<!-- field: packages -->\n### Packages\n\nPackages are listed elsewhere.\n",
		"<!-- field: packages -->\n#### Packages\n\nTo be decided.\n",
	} {
		if !strings.Contains(md, want) {
			t.Errorf("printed Markdown lacks %q\n%s", want, md)
		}
	}
}

// TestRecordFieldKeepsAStringList prints a list of strings under a key whose
// level declares a record kind as a plain list.
func TestRecordFieldKeepsAStringList(t *testing.T) {
	d := nestedDeclaration(t)
	original := yamlDoc(t, "id: plan\ntitle: Plan\noverview:\n  layout:\n    packages: [internal/a, internal/b]\n")
	md, back := nestedRoundTrip(t, d, original)
	if diff := Diff(original, back); diff != "" {
		t.Fatalf("round trip changed %s\n%s", diff, md)
	}
	if !strings.Contains(md, "#### Packages\n\n- internal/a\n- internal/b\n") {
		t.Errorf("string list did not print as a list\n%s", md)
	}
}

// TestFlatTableKindPrintsTableOrRecords prints one record kind declared
// flat: table as a linked table when every record is flat, as a typed table
// when a value is not a string, and as tagged records when a record holds a
// paragraph; a kind without flat always prints records (GH-2674).
func TestFlatTableKindPrintsTableOrRecords(t *testing.T) {
	d := nestedDeclaration(t)
	indexes := LinkIndexes{}
	indexes.Add("docs/specs/srd001-a.md", []byte("---\nid: srd001-a\ntitle: A\n---\n"))
	for name, c := range map[string]struct{ yaml, want string }{
		"flat strings": {"index:\n  - id: srd001-a\n    title: A\n  - id: srd002-b\n    title: B\n",
			"<!-- field: index -->\n## Index\n\n| id | title |\n| --- | --- |\n| [srd001-a](specs/srd001-a.md) | A |\n| srd002-b | B |\n"},
		"flat typed": {"index:\n  - id: srd001-a\n    cases: 3\n",
			"<!-- field: index | form: \"typed table\" -->\n## Index\n\n| id | cases |\n| --- | --- |\n| [srd001-a](specs/srd001-a.md) | 3 |\n"},
		"rich": {"index:\n  - id: srd001-a\n    title: A\n    summary: |\n      One line.\n\n      Another paragraph.\n  - id: srd002-b\n    title: B\n",
			"<!-- field: index -->\n## Index\n\n<!-- index_entry: srd001-a -->\n### srd001-a\n\n<!-- field: title -->\n#### Title\n\nA\n\n<!-- field: summary -->\n#### Summary\n\nOne line.\n\nAnother paragraph.\n\n<!-- index_entry: srd002-b -->\n### srd002-b\n"},
		"no flat on packages": {"overview:\n  layout:\n    packages:\n      - path: a/\n        responsibility: Short.\n",
			"<!-- field: packages -->\n#### Packages\n\n<!-- package: a/ -->\n##### a/\n"},
	} {
		t.Run(name, func(t *testing.T) {
			original := yamlDoc(t, "id: plan\ntitle: Plan\n"+c.yaml)
			md, err := d.PrintLinked(original, "docs/plan.md", indexes.For("docs/plan.md"))
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(md), c.want) {
				t.Fatalf("printed Markdown lacks %q\n%s", c.want, md)
			}
			back, findings, err := d.Parse(md, "docs/plans/plan.md")
			if err != nil || len(findings) > 0 {
				t.Fatalf("parse: %v %v\n%s", err, findings, md)
			}
			if diff := Diff(original, back); diff != "" {
				t.Fatalf("round trip changed %s\n%s", diff, md)
			}
		})
	}
}

func TestFlatIsOnlyTableOnListRecords(t *testing.T) {
	for _, c := range []struct{ name, node, want string }{
		{"on a map kind", "x: {tag: x, level: 2, parent: document, heading: id, content: none, flat: table, yaml: {from: xs, records: map}}", "list records"},
		{"not table", "x: {tag: x, level: 2, parent: document, heading: id, content: none, flat: rows, yaml: {from: xs, records: list, id: id}}", "only table"},
	} {
		design := "document_types:\n  plan:\n    location: docs/plans/name.md\n    markdown:\n      frontmatter: {required: [id], allowed: [id]}\n      nodes:\n        " + c.node + "\n"
		if _, err := LoadDeclarations([]byte(design)); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: LoadDeclarations error = %v, want %q", c.name, err, c.want)
		}
	}
}
