// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"os"
	"strings"
	"testing"
)

// validGuideline is a minimal engineering guideline in the Markdown format.
// It is checked against the engineering_guideline declaration in design.yaml,
// which shares no node, field, or rule with the SRD declaration: a new
// document type is data, not parser code.
const validGuideline = `---
id: eng90-sample
title: Sample Guideline
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: introduction -->
## Introduction

We describe a sample practice.

<!-- section: the-practice -->
## The practice

Each module follows the practice.

Table 1: Sample

| a | b |
| --- | --- |
| 1 | 2 |

<!-- field: references -->
## References

- docs/engineering/eng09-spec-kit-correspondence.yaml
`

const guidelineName = "docs/engineering/eng90-sample.md"

func TestValidGuidelineHasNoFindings(t *testing.T) {
	if findings := diagnose(t, []byte(validGuideline), guidelineName); len(findings) > 0 {
		t.Fatalf("valid guideline reported findings:\n%v", findings)
	}
}

func TestGuidelineDiagnostics(t *testing.T) {
	section := validGuideline[strings.Index(validGuideline, "<!-- section:"):strings.Index(validGuideline, "<!-- field: references")]
	introduction := validGuideline[strings.Index(validGuideline, "<!-- field: introduction"):strings.Index(validGuideline, "<!-- section:")]
	cases := []diagnosticCase{
		{"missing-required-node", strings.Replace(validGuideline, section, "", 1), "", lastLine,
			"add a <!-- section: <id> --> marker above a level-2 heading"},
		{"heading-misplaced", strings.Replace(validGuideline, "## The practice", "### The practice", 1), "", 14,
			"change the heading depth"},
		{"missing-required-field", strings.Replace(validGuideline, introduction, "", 1), "", lastLine,
			"add <!-- field: introduction -->"},
		{"id-malformed", strings.Replace(validGuideline, "section: the-practice", "section: The Practice", 1), "", 14,
			"use the declared id form"},
		{"prose-containment", strings.Replace(validGuideline, "Each module follows", "Each module MUST follow", 1), "", 16,
			"move the obligation into a requirement item"},
		{"marker-unknown-tag", strings.Replace(validGuideline, "<!-- section: the-practice -->", "<!-- req: R1 -->", 1), "", 13,
			"use one of field, section"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			line := tc.line
			if line == lastLine {
				line = strings.Count(tc.doc, "\n")
			}
			for _, f := range diagnose(t, []byte(tc.doc), guidelineName) {
				if f.Name != tc.name {
					continue
				}
				if f.Line != line || !strings.Contains(f.Hint, tc.hint) {
					t.Fatalf("got %s, want line %d and hint containing %q", f, line, tc.hint)
				}
				return
			}
			t.Fatalf("no %s finding; got:\n%v", tc.name, diagnose(t, []byte(tc.doc), guidelineName))
		})
	}
}

// TestEng10IsClean keeps the format specification its own first conformance
// test: eng10 is written in the format it specifies and reports nothing.
func TestEng10IsClean(t *testing.T) {
	const path = "../../../../docs/engineering/eng10-markdown-spec-format.md"
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if findings := diagnose(t, data, path); len(findings) > 0 {
		t.Fatalf("eng10 reported findings:\n%v", findings)
	}
}
