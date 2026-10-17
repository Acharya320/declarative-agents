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
