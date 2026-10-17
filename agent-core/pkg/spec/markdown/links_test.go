// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"strings"
	"testing"
)

// linkedDesign declares requirement records and three reference fields: a
// list, a one-line value, and a table column (GH-2665).
const linkedDesign = `
document_types:
  plan:
    location: docs/plans/name.md
    markdown:
      frontmatter: {required: [id, title], allowed: [id, title]}
      nodes:
        document_field: {tag: field, level: 2, parent: document, heading: free, content: any}
        req:
          tag: req
          level: 2
          parent: document
          heading: id_title
          content: prose
          yaml: {from: requirements, records: list, id: id, title: title, content: text}
        req_field: {tag: field, level: 3, parent: req, heading: free, content: any}
      fields:
        traces: {in: [document_field], shape: list, links: true}
        use_case: {in: [document_field], links: true}
        touchpoints: {in: [document_field], shape: table, link_columns: [target]}
`

const linkedPlan = `id: plan
title: Plan
requirements:
  - id: R1
    title: Budget
    text: Keep a budget.
  - id: R2
    text: Spend it.
traces:
  - R1
  - R2
  - R9
  - not a reference
  - "[R1](https://example.com)"
use_case: other-doc
touchpoints:
  - id: T1
    target: other-doc AC1, AC2
    reason: why
`

func linkedDeclaration(t *testing.T) *Declaration {
	t.Helper()
	declarations, err := LoadDeclarations([]byte(linkedDesign))
	if err != nil {
		t.Fatal(err)
	}
	return declarations["plan"]
}

func TestSlugFollowsGitHub(t *testing.T) {
	for heading, want := range map[string]string{
		"R1.1": "r11", "R1 - Budget": "r1---budget", "AC5": "ac5", "Loop function": "loop-function",
		"R2 - RunStatus and RunResult": "r2---runstatus-and-runresult", "snake_case (v2)!": "snake_case-v2",
	} {
		if got := slug(heading); got != want {
			t.Errorf("slug(%q) = %q, want %q", heading, got, want)
		}
	}
}

// TestPageAnchorsCountRepeats gives an item whose heading repeats earlier
// headings GitHub's -N suffix.
func TestPageAnchorsCountRepeats(t *testing.T) {
	lines := strings.Split("<!-- field: traces -->\n### Traces\n\n<!-- field: traces -->\n### Traces\n\n"+
		"<!-- req: R7 -->\n## Traces\n\n```text\n# not a heading\n```\n<!-- ac: AC1 -->\n### AC1\n", "\n")
	anchors := pageAnchors(lines)
	if anchors["R7"] != "traces-2" || anchors["AC1"] != "ac1" || len(anchors) != 2 {
		t.Fatalf("anchors = %v", anchors)
	}
}

// TestSamePageReferencesLink links each bare item the page holds to its
// anchor, leaves unknown items, other text, and a hand-written link as
// written, and reads every value back exactly.
func TestSamePageReferencesLink(t *testing.T) {
	d := linkedDeclaration(t)
	original := yamlDoc(t, linkedPlan)
	md, err := d.Print(original)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"- [R1](#r1---budget)\n- [R2](#r2)\n- R9\n- not a reference\n- [R1](https://example.com)\n",
		"## Use case\n\nother-doc\n", // no index: other documents stay text
		"| T1 | other-doc AC1, AC2 | why |",
	} {
		if !strings.Contains(string(md), want) {
			t.Errorf("printed Markdown lacks %q\n%s", want, md)
		}
	}
	back, findings, err := d.Parse(md, "docs/plans/plan.md")
	if err != nil || len(findings) > 0 {
		t.Fatalf("parse: %v %v", err, findings)
	}
	if diff := Diff(original, back); diff != "" {
		t.Fatalf("round trip changed %s\n%s", diff, md)
	}
}

// TestUnlinkReadsBackOnlyGeneratedLinks keeps a link whose target does not
// agree with its text.
func TestUnlinkReadsBackOnlyGeneratedLinks(t *testing.T) {
	for in, want := range map[string]string{
		"[R1.1](#r11)": "R1.1",
		"[srd005-agentic-loop](../software-requirements/srd005-agentic-loop.md) [AC5](../software-requirements/srd005-agentic-loop.md#ac5)": "srd005-agentic-loop AC5",
		"[R1.1](#ac2)":                    "[R1.1](#ac2)",
		"[srd005-agentic-loop](other.md)": "[srd005-agentic-loop](other.md)",
		"[see this](#r11)":                "[see this](#r11)",
	} {
		if got := unlink(in); got != want {
			t.Errorf("unlink(%q) = %q, want %q", in, got, want)
		}
	}
}
