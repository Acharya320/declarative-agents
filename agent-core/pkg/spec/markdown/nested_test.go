// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"strings"
	"testing"
)

// TestTreeFormRoundTripsCorpusShapes prints the structures that printed as
// fenced yaml in the corpus (test-case inputs and expected, an SRD
// implementation_status, pattern consequences, a diagram map) as nested
// lists and reads them back through the whole document parser (GH-2642).
func TestTreeFormRoundTripsCorpusShapes(t *testing.T) {
	d := toyDeclaration(t)
	original := yamlDoc(t, recipeWith(`inputs:
  profile: agents/executor/profile.yaml
  retries: 3
  words: [read, edit]
  env: {}
expected:
  exit_code: 0
  signals:
    - name: Done
      count: 1
    - name: Retry
      count: 0
  note: "#601"
implementation_status:
  state: partial
  evidence:
    - path: magefiles/a_test.go
      proves: The audit runs.
    - path: magefiles/b_test.go
      proves: "Counts: two."
consequences:
  benefits:
    - Rules live in data.
  liabilities: []
diagrams:
  requirements_trace: |
    `+"```mermaid"+`
    flowchart LR
      R1 --> AC1
    `+"```"+`
  note: |
    Two lines of text,
    the second one here.
`))
	md, back := roundTrip(t, d, original)
	if diff := Diff(original, back); diff != "" {
		t.Fatalf("round trip changed %s\n%s", diff, md)
	}
	for field, want := range map[string]string{
		"inputs":                "- profile: agents/executor/profile.yaml\n- retries: 3\n- words:\n  1. read\n  1. edit\n- env: {}",
		"expected":              "- exit_code: 0\n- signals:\n  1. - name: Done\n     - count: 1\n  1. - name: Retry\n     - count: 0\n- note: \"#601\"",
		"implementation_status": "- state: partial\n- evidence:\n  1. - path: magefiles/a_test.go\n     - proves: The audit runs.\n  1. - path: magefiles/b_test.go\n     - proves: \"Counts: two.\"",
		"diagrams":              "- requirements_trace:\n  ```mermaid\n  flowchart LR\n    R1 --> AC1\n  ```\n- note:\n  Two lines of text,\n  the second one here.",
	} {
		if !strings.Contains(md, "<!-- field: "+field+` | form: "tree" -->`) {
			t.Errorf("%s did not print as a tree", field)
		}
		if section := fieldSection(md, field); section != want {
			t.Errorf("%s printed as:\n%s\nwant:\n%s", field, section, want)
		}
	}
}

// TestTreeFormQuotesKeysAndMarksLiteralText prints a prose key quoted and
// text whose first line reads as a list item or is indented after a literal
// marker, and reads both back (GH-2644).
func TestTreeFormQuotesKeysAndMarksLiteralText(t *testing.T) {
	d := toyDeclaration(t)
	original := yamlDoc(t, recipeWith(`notes:
  "Outbound injection (R4)": done
  "say \"hi\": now": 2
  body: |
    - looks like a list
    second line
  code: "    indented first line\nback out"
`))
	md, back := roundTrip(t, d, original)
	if diff := Diff(original, back); diff != "" {
		t.Fatalf("round trip changed %s\n%s", diff, md)
	}
	want := "- \"Outbound injection (R4)\": done\n- \"say \\\"hi\\\": now\": 2\n- body: |\n  - looks like a list\n  second line\n- code: |\n      indented first line\n  back out"
	if section := fieldSection(md, "notes"); section != want {
		t.Errorf("notes printed as:\n%s\nwant:\n%s", section, want)
	}
}

// TestTreeFormKeepsTheFenceForTextOpeningBlank keeps a fenced yaml block for
// text that opens with a blank line, which no rendering reads back.
func TestTreeFormKeepsTheFenceForTextOpeningBlank(t *testing.T) {
	d := toyDeclaration(t)
	original := yamlDoc(t, recipeWith("notes:\n  body: \"\\nopens blank\\nsecond\"\n"))
	md, back := roundTrip(t, d, original)
	if diff := Diff(original, back); diff != "" {
		t.Fatalf("round trip changed %s\n%s", diff, md)
	}
	if !strings.Contains(fieldSection(md, "notes"), "```yaml") {
		t.Errorf("printed without the fence:\n%s", fieldSection(md, "notes"))
	}
}

// TestBareMarkerItemsKeepTheirFirstLine reads a top-level numbered list whose
// first item opens with a bare "1." back from its marker line, not from the
// first line of its content.
func TestBareMarkerItemsKeepTheirFirstLine(t *testing.T) {
	blocks, lines := topLevelBlocks([]byte("Intro.\n\n1.\n   - a: 1\n"))
	if len(blocks) != 2 || blocks[1].Start != 2 || lines[blocks[1].Start] != "1." {
		t.Fatalf("blocks %+v, want the list to start at the bare marker line 2", blocks)
	}
}
