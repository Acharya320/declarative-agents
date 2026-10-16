// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package specmd

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec"
)

// validSRD is a minimal canonical SRD; each case below mutates it to provoke
// exactly the finding it names. Line numbers in the cases refer to the
// mutated document.
const validSRD = `---
id: srd900-sample
title: Sample
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Something is missing.

<!-- field: goals -->
## Goals

- G1: Fix it.

<!-- req: R1 -->
## R1 - Contract

<!-- req: R1.1 -->
### R1.1

The kit shall export one constant.

<!-- field: endpoints -->
#### Endpoints

| path | owner |
| --- | --- |
| GET /state | srd033 R3.9 |

<!-- field: non_goals -->
## Non-goals

- Rewriting the shell.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 | traces: ["R1.1"] -->
### AC1

A kit test fails on an unknown endpoint.
`

const validName = "srd900-sample.md"

func TestValidSRDHasNoFindings(t *testing.T) {
	if findings := Diagnose([]byte(validSRD), validName); len(findings) > 0 {
		t.Fatalf("valid document reported findings:\n%v", findings)
	}
}

type diagnosticCase struct {
	name string
	doc  string
	file string
	line int
	hint string
}

func replaceOnce(old, replacement string) string {
	return strings.Replace(validSRD, old, replacement, 1)
}

func diagnosticCases() []diagnosticCase {
	return []diagnosticCase{
		{"frontmatter-missing", strings.SplitN(validSRD, "---\n", 3)[2], "", 1, "start the file with ---"},
		{"frontmatter-invalid", replaceOnce("title: Sample", "title: [unclosed"), "", 2, "key: value lines"},
		{"frontmatter-missing-key", replaceOnce("title: Sample\n", ""), "", 2, `add "title: ..."`},
		{"frontmatter-unknown-key", replaceOnce("title: Sample\n", "title: Sample\nstatuss: draft\n"), "", 4, `did you mean "status"?`},
		{"frontmatter-id-mismatch", validSRD, "srd901-other.md", 2, "rename the file or change the id"},
		{"marker-malformed", replaceOnce("<!-- field: goals -->", "<!-- field goals -->"), "", 13, "separated by a colon"},
		{"marker-unknown-tag", replaceOnce("<!-- field: goals -->", "<!-- feild: goals -->"), "", 13, `did you mean "field"?`},
		{"marker-unknown-attribute", replaceOnce("<!-- req: R1.1 -->", "<!-- req: R1.1 | statsu: retired -->"), "", 21, `did you mean "status"?`},
		{"marker-orphan", replaceOnce("<!-- field: goals -->", "<!-- field: notes -->\n<!-- field: goals -->"), "", 13, "add the heading the marker tags"},
		{"heading-untagged", replaceOnce("<!-- field: problem -->\n", ""), "", 8, "add the marker"},
		{"heading-misplaced", replaceOnce("#### Endpoints", "##### Endpoints"), "", 27, "change the heading depth"},
		{"id-malformed", replaceOnce("<!-- req: R1.1 -->\n### R1.1", "<!-- req: R1.x -->\n### R1.x"), "", 22, "renumber to the declared form"},
		{"item-wrong-group", replaceOnce("<!-- req: R1.1 -->\n### R1.1", "<!-- req: R2.1 -->\n### R2.1"), "", 22, "move the item under its own group"},
		{"heading-id-mismatch", replaceOnce("### R1.1", "### R1.2"), "", 22, "make the heading say exactly the id"},
		{"duplicate-id", replaceOnce("<!-- field: non_goals -->", "<!-- req: R1.1 -->\n### R1.1\n\nAgain.\n\n<!-- field: non_goals -->"), "", 34, "give the second node a new id"},
		{"field-disallowed", replaceOnce("<!-- field: endpoints -->", "<!-- field: goals -->"), "", 27, "move the section to where the field belongs"},
		{"field-duplicate", replaceOnce("<!-- req: R1 -->", "<!-- field: goals -->\n## Goals\n\n- G2: More.\n\n<!-- req: R1 -->"), "", 19, "merge the two sections"},
		{"content-untagged", replaceOnce("## R1 - Contract\n", "## R1 - Contract\n\nStray text.\n"), "", 21, "tag it with <!-- field: <key> -->"},
		{"untagged-structure", replaceOnce("The kit shall export one constant.\n", "The kit shall export one constant.\n\n- a bullet\n"), "", 26, "move the list under its own"},
		{"missing-content", replaceOnce("The kit shall export one constant.\n\n", ""), "", 22, "write the requirement_item text"},
		{"missing-required-field", replaceOnce("<!-- field: non_goals -->\n## Non-goals\n\n- Rewriting the shell.\n\n", ""), "", 39, "add <!-- field: non_goals -->"},
		{"missing-requirement-groups", strings.Replace(validSRD, validSRD[strings.Index(validSRD, "<!-- req: R1 -->"):strings.Index(validSRD, "<!-- field: non_goals -->")], "", 1), "", 29, `add <!-- req: R1 -->`},
		{"field-shape-mismatch", replaceOnce("- G1: Fix it.", "G1: Fix it."), "", 16, "write goals as list"},
		{"list-item-prefix", replaceOnce("- G1: Fix it.", "- Fix it."), "", 16, `number the item`},
		{"table-columns-mismatch", replaceOnce("| path | owner |", "| path | team |"), "", 29, "use exactly the declared columns"},
		{"attribute-shape", replaceOnce("<!-- req: R1.1 -->", "<!-- req: R1.1 | weight: heavy -->"), "", 21, "write the value in its declared shape"},
		{"criterion-traces-missing", replaceOnce(`<!-- ac: AC1 | traces: ["R1.1"] -->`, "<!-- ac: AC1 -->"), "", 41, `add | traces: ["R1.1"]`},
		{"criterion-trace-invalid", replaceOnce(`traces: ["R1.1"]`, `traces: ["X9"]`), "", 41, "trace a requirement id"},
		{"prose-containment", replaceOnce("Something is missing.", "The system SHALL do it."), "", 11, "move the obligation into a requirement item"},
	}
}

// TestEveryFindingIsDiagnosedAtItsLineWithAHint provokes each finding name
// from one fixture and checks its name, file line, and fix hint (GH-2589 AC2).
func TestEveryFindingIsDiagnosedAtItsLineWithAHint(t *testing.T) {
	for _, tc := range diagnosticCases() {
		t.Run(tc.name, func(t *testing.T) {
			file := tc.file
			if file == "" {
				file = validName
			}
			findings := Diagnose([]byte(tc.doc), file)
			for _, f := range findings {
				if f.Name != tc.name {
					continue
				}
				if f.Line != tc.line || !strings.Contains(f.Hint, tc.hint) || f.Tier != findingCatalog[tc.name] {
					t.Fatalf("got %s (tier %d), want line %d, tier %d, hint containing %q", f, f.Tier, tc.line, findingCatalog[tc.name], tc.hint)
				}
				return
			}
			t.Fatalf("no %s finding; got:\n%v", tc.name, findings)
		})
	}
}

func TestEveryCatalogFindingHasAFixture(t *testing.T) {
	covered := map[string]bool{}
	for _, tc := range diagnosticCases() {
		covered[tc.name] = true
	}
	for name := range findingCatalog {
		if !covered[name] {
			t.Errorf("finding %s has no diagnostic fixture", name)
		}
	}
	if len(covered) != len(findingCatalog) {
		t.Errorf("%d fixtures for %d catalog names", len(covered), len(findingCatalog))
	}
}

// printerFindings are the findings only a malformed printer could cause: the
// markers, frontmatter, headings, and placement it writes. Findings about a
// document's own content (its ids, traces, wording, missing sections) are
// not the printer's.
var printerFindings = map[string]bool{
	"frontmatter-missing": true, "frontmatter-invalid": true, "frontmatter-missing-key": true,
	"frontmatter-unknown-key": true, "marker-malformed": true, "marker-unknown-tag": true,
	"marker-unknown-attribute": true, "marker-orphan": true, "heading-untagged": true,
	"heading-misplaced": true, "heading-id-mismatch": true, "content-untagged": true,
	"untagged-structure": true, "attribute-shape": true, "field-duplicate": true, "field-disallowed": true,
}

// TestPrintedCorpusIsStructurallyClean diagnoses the printed form of every
// tracked SRD: the printer must never produce a structural finding.
func TestPrintedCorpusIsStructurallyClean(t *testing.T) {
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	out, err := exec.Command("git", "-C", repo, "ls-files").Output()
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range strings.Split(string(out), "\n") {
		if !strings.Contains(path, "/docs/specs/software-requirements/srd") ||
			!strings.HasSuffix(path, ".yaml") || strings.Contains(path, "/testdata/") {
			continue
		}
		srd, err := spec.ParseSRD(filepath.Join(repo, path))
		if err != nil {
			continue
		}
		md, err := PrintMarkdown(&srd)
		if err != nil {
			t.Errorf("%s: %v", path, err)
			continue
		}
		name := strings.TrimSuffix(filepath.Base(path), ".yaml") + ".md"
		for _, f := range Diagnose(md, name) {
			if printerFindings[f.Name] {
				t.Errorf("%s: printer produced %s", path, f)
			}
		}
	}
}
