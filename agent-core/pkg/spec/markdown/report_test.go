// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"strings"
	"testing"
)

// TestFieldHeadingsKeepAcronyms generates headings that keep SRD, ID, and
// the other acronyms in their conventional case (GH-2643).
func TestFieldHeadingsKeepAcronyms(t *testing.T) {
	d := toyDeclaration(t)
	for key, want := range map[string]string{
		"srd_index": "SRD index", "id": "ID", "rest_api_url": "REST API URL",
		"otel_export": "OTel export", "ui_ux_notes": "UI UX notes", "ollama_rest": "Ollama REST",
		"description": "Description",
	} {
		if got := d.fieldHeading(key); got != want {
			t.Errorf("fieldHeading(%q) = %q, want %q", key, got, want)
		}
	}
}

// TestOldHeadingTextStillParses reads a field printed with the old heading
// text, and Format rewrites the heading.
func TestOldHeadingTextStillParses(t *testing.T) {
	d := toyDeclaration(t)
	md, _ := roundTrip(t, d, yamlDoc(t, recipeWith("srd_index: [srd001]\n")))
	old := strings.Replace(md, "## SRD index", "## Srd index", 1)
	if old == md {
		t.Fatalf("no SRD index heading:\n%s", md)
	}
	formatted, findings, err := d.Format([]byte(old), toyName)
	if err != nil || len(findings) > 0 {
		t.Fatalf("old heading did not parse: %v %v", err, findings)
	}
	if string(formatted) != md {
		t.Errorf("Format did not restore the new heading:\n%s", formatted)
	}
}

// TestPrintReportsEveryRemainingFence names the field and the reason for a
// value only the fenced yaml block holds, and nothing for one a form holds.
func TestPrintReportsEveryRemainingFence(t *testing.T) {
	d := toyDeclaration(t)
	doc := yamlDoc(t, recipeWith("notes:\n  body: \"\\nopens blank\\nsecond\"\nsizes: {small: 1}\n"))
	md, fences, err := d.PrintReport(doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(fences) != 1 || fences[0].Field != "notes" || !strings.Contains(fences[0].Reason, "body") {
		t.Fatalf("fences = %v, want one for notes naming body", fences)
	}
	again, err := d.Fences(md, toyName)
	if err != nil || len(again) != 1 || again[0] != fences[0] {
		t.Errorf("Fences on the printed document = %v, %v", again, err)
	}
	_, report, err := d.ConvertYAML([]byte(recipeWith("notes:\n  body: \"\\nopens blank\\nsecond\"\n")), "docs/recipes/soup.yaml")
	if err != nil || len(report.Fences) != 1 || !strings.Contains(strings.Join(report.Lines(), "\n"), "yaml fence: notes:") {
		t.Errorf("conversion report lacks the fence: %v %v", report.Lines(), err)
	}
}
