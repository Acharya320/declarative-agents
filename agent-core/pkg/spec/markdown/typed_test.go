// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"strings"
	"testing"
)

// TestTypedTableKeepsTypesAndOptionalKeys prints a list of records holding
// numbers, booleans, a list, a string that looks like a number, and keys some
// records lack as a typed table, and reads every type back (GH-2641).
func TestTypedTableKeepsTypesAndOptionalKeys(t *testing.T) {
	d := toyDeclaration(t)
	original := yamlDoc(t, recipeWith(`index:
  - version: "00.0"
    done: 2
    shipped: true
    path: docs/a.md
  - version: "01.0"
    done: 0
    status: planned
    shipped: false
    tags: [hot, "1", "a, b"]
    path: docs/b.md
  - version: "02.0"
    note: ""
    path: "x | y"
`))
	md, back := roundTrip(t, d, original)
	if diff := Diff(original, back); diff != "" {
		t.Fatalf("round trip changed %s\n%s", diff, md)
	}
	if !strings.Contains(md, `<!-- field: index | form: "typed table" -->`) {
		t.Fatalf("index did not print as a typed table:\n%s", md)
	}
	section := fieldSection(md, "index")
	for _, want := range []string{
		"| version | done | status | shipped | tags | note | path |",
		`| "00.0" | 2 |  | true |  |  | docs/a.md |`,
		`| "01.0" | 0 | planned | false | [hot, "1", "a, b"] |  | docs/b.md |`,
		`| "02.0" |  |  |  |  | "" | x \| y |`,
	} {
		if !strings.Contains(section, want) {
			t.Errorf("typed table lacks %q:\n%s", want, section)
		}
	}
}

// TestScalarFormKeepsLoneValues prints a number, a null, and a boolean as one
// flow value each, and reads them back; a string stays text.
func TestScalarFormKeepsLoneValues(t *testing.T) {
	d := toyDeclaration(t)
	original := yamlDoc(t, recipeWith("servings: 12\nfollow_up: \"#601\"\nversion: \"1.0\"\nretired: null\nflag: true\n"))
	md, back := roundTrip(t, d, original)
	if diff := Diff(original, back); diff != "" {
		t.Fatalf("round trip changed %s\n%s", diff, md)
	}
	for field, want := range map[string]string{"servings": "12", "retired": "null", "flag": "true"} {
		if !strings.Contains(md, "<!-- field: "+field+` | form: "scalar" -->`) || fieldSection(md, field) != want {
			t.Errorf("%s printed as %q, want scalar %q", field, fieldSection(md, field), want)
		}
	}
	// Strings that read back as text print as text, with no form: "#601" is
	// not a heading in CommonMark, which needs a space after the hash.
	for field, want := range map[string]string{"version": "1.0", "follow_up": "#601"} {
		if fieldSection(md, field) != want || strings.Contains(md, "<!-- field: "+field+" | form") {
			t.Errorf("%s printed as %q with a form, want text %q", field, fieldSection(md, field), want)
		}
	}
}

// TestScalarListsPrintAsLists keeps a list of plain strings as a bullet list
// with no form, and a list holding a number as a typed scalar list in a
// table-free fallback that still round-trips.
func TestScalarListsPrintAsLists(t *testing.T) {
	d := toyDeclaration(t)
	original := yamlDoc(t, recipeWith("sources:\n  - docs/a.md\n  - docs/b.md\nsizes:\n  - 1\n  - 2\n"))
	md, back := roundTrip(t, d, original)
	if diff := Diff(original, back); diff != "" {
		t.Fatalf("round trip changed %s\n%s", diff, md)
	}
	if fieldSection(md, "sources") != "- docs/a.md\n- docs/b.md" {
		t.Errorf("sources printed as %q", fieldSection(md, "sources"))
	}
}

// TestUnknownFormIsAFinding rejects a form name the format does not define.
func TestUnknownFormIsAFinding(t *testing.T) {
	d := toyDeclaration(t)
	md, _ := roundTrip(t, d, yamlDoc(t, recipeWith("servings: 12\n")))
	broken := strings.Replace(md, `form: "scalar"`, `form: "number"`, 1)
	if _, findings, _ := d.Parse([]byte(broken), toyName); len(findings) == 0 {
		t.Error("an unknown form raised no finding")
	}
}
