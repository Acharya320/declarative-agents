// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"strings"
	"testing"
)

// orderedDesign lists the fields of the document slot and of overview's
// nested slot; the layout slot lists none, so it is not checked (GH-2654).
const orderedDesign = `
document_types:
  plan:
    location: docs/plans/name.md
    markdown:
      frontmatter: {required: [id, title], allowed: [id, title]}
      nodes:
        document_field: {tag: field, level: 2, parent: document, heading: free, content: any, fields: [overview, notes]}
        overview_field: {tag: field, level: 3, parent: 'field:overview', heading: free, content: any, fields: [summary, layout]}
        layout_field: {tag: field, level: 4, parent: 'field:layout', heading: free, content: any}
      fields:
        notes: {in: [document_field], shape: prose}
`

func orderedDeclaration(t *testing.T) *Declaration {
	t.Helper()
	declarations, err := LoadDeclarations([]byte(orderedDesign))
	if err != nil {
		t.Fatal(err)
	}
	return declarations["plan"]
}

// planDoc assembles a plan from field sections, in the order given.
func planDoc(sections ...string) []byte {
	return []byte("---\nid: plan\ntitle: Plan\n---\n\n" + strings.Join(sections, "\n"))
}

const (
	overviewOpen  = "<!-- field: overview -->\n## Overview\n"
	summaryField  = "<!-- field: summary -->\n### Summary\n\nText.\n"
	layoutField   = "<!-- field: layout -->\n### Layout\n\n<!-- field: zeta -->\n#### Zeta\n\nz\n\n<!-- field: alpha -->\n#### Alpha\n\na\n"
	notesField    = "<!-- field: notes -->\n## Notes\n\nPlain.\n"
	strangerField = "<!-- field: stranger -->\n## Stranger\n\nOdd.\n"
)

func TestFieldOrderChecks(t *testing.T) {
	d := orderedDeclaration(t)
	for _, c := range []struct {
		name string
		doc  []byte
		want []string // finding names, in order
	}{
		{"declared order", planDoc(overviewOpen, summaryField, layoutField, notesField), nil},
		{"top level out of order", planDoc(notesField, overviewOpen, summaryField), []string{"field-out-of-order"}},
		{"nested out of order", planDoc(overviewOpen, layoutField, summaryField, notesField), []string{"field-out-of-order"}},
		{"undeclared top level", planDoc(overviewOpen, summaryField, strangerField, notesField), []string{"field-undeclared"}},
		{"undeclared nested", planDoc(overviewOpen, summaryField, "<!-- field: extra -->\n### Extra\n\nx\n"), []string{"field-undeclared"}},
	} {
		t.Run(c.name, func(t *testing.T) {
			var got []string
			for _, f := range d.Diagnose(c.doc, "docs/plans/plan.md") {
				got = append(got, f.Name)
			}
			if strings.Join(got, ",") != strings.Join(c.want, ",") {
				t.Fatalf("findings = %v, want %v\n%s", got, c.want, c.doc)
			}
		})
	}
}

// TestFieldOrderNamesTheField says which field is out of place and what it
// belongs ahead of.
func TestFieldOrderNamesTheField(t *testing.T) {
	findings := orderedDeclaration(t).Diagnose(planDoc(notesField, overviewOpen, summaryField), "docs/plans/plan.md")
	if len(findings) != 1 || findings[0].Expected != "overview ahead of notes" ||
		!strings.HasSuffix(findings[0].Rule, "nodes.document_field.fields") {
		t.Fatalf("findings = %v", findings)
	}
}

func TestSlotDeclarationChecks(t *testing.T) {
	for _, c := range []struct{ name, nodes, fields, want string }{
		{"fields on a record kind",
			"thing: {tag: thing, level: 2, parent: document, heading: id, content: prose, fields: [a], yaml: {from: things, records: list, id: id}}",
			"", "only a field slot"},
		{"field listed twice", "", "", "listed twice"},
		{"rule names a slot that lacks it", "", "stray: {in: [document_field], shape: prose}", "lacks it"},
	} {
		t.Run(c.name, func(t *testing.T) {
			slot := "document_field: {tag: field, level: 2, parent: document, heading: free, content: any, fields: [a, b]}"
			if c.want == "listed twice" {
				slot = "document_field: {tag: field, level: 2, parent: document, heading: free, content: any, fields: [a, a]}"
			}
			design := "document_types:\n  plan:\n    location: docs/plans/name.md\n    markdown:\n" +
				"      frontmatter: {required: [id], allowed: [id]}\n      nodes:\n        " + slot + "\n"
			if c.nodes != "" {
				design += "        " + c.nodes + "\n"
			}
			if c.fields != "" {
				design += "      fields:\n        " + c.fields + "\n"
			}
			if _, err := LoadDeclarations([]byte(design)); err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("LoadDeclarations error = %v, want %q", err, c.want)
			}
		})
	}
}
