// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// recipeWith returns the minimal recipe plus extra top-level YAML.
func recipeWith(extra string) string {
	return "id: soup\ntitle: Soup\ndescription: A soup.\nsteps:\n  - id: S1\n    text: Heat.\n    checks:\n      - hot\n" + extra
}

// fieldSection returns the Markdown under a field's heading, up to the next
// marker.
func fieldSection(md, field string) string {
	start := -1
	for _, open := range []string{"<!-- field: " + field + " -->\n", "<!-- field: " + field + " | "} {
		if i := strings.Index(md, open); i >= 0 {
			start = i
			break
		}
	}
	if start < 0 {
		return ""
	}
	_, after, _ := strings.Cut(md[start:], "\n") // past the marker
	_, body, _ := strings.Cut(after, "\n")       // past the heading
	if end := strings.Index(body, "<!-- "); end >= 0 {
		body = body[:end]
	}
	return strings.TrimSpace(body)
}

// TestTextFieldsKeepMarkdownVerbatim prints text holding a Mermaid fence, a
// block quote, and indented lines as the field's Markdown, so it renders,
// and reads it back unchanged (GH-2640).
func TestTextFieldsKeepMarkdownVerbatim(t *testing.T) {
	d := toyDeclaration(t)
	original := yamlDoc(t, recipeWith(`diagram: |
  How the soup is made.

  `+"```mermaid"+`
  flowchart LR
    # not a heading inside the fence
    Water --> Soup
  `+"```"+`
story: |
  > A quoted line.

  First line
      indented continuation
`))
	md, back := roundTrip(t, d, original)
	if diff := Diff(original, back); diff != "" {
		t.Fatalf("round trip changed %s\n%s", diff, md)
	}
	for _, field := range []string{"diagram", "story"} {
		if section := fieldSection(md, field); strings.Contains(section, "yaml") || !strings.HasPrefix(section, strings.SplitN(mapValue(original.Content[0], field).Value, "\n", 2)[0]) {
			t.Errorf("%s did not print as Markdown text:\n%s", field, section)
		}
	}
	if !strings.Contains(fieldSection(md, "diagram"), "```mermaid\nflowchart LR") {
		t.Errorf("Mermaid fence lost:\n%s", fieldSection(md, "diagram"))
	}
}

// TestTextThatWouldBreakStructureKeepsTheFence keeps a fenced yaml block for
// text holding a heading or a marker line, which would read back as
// structure.
func TestTextThatWouldBreakStructureKeepsTheFence(t *testing.T) {
	d := toyDeclaration(t)
	for name, text := range map[string]string{
		"heading": "Intro.\n\n# Not a field\n\nMore.",
		"marker":  "Intro.\n\n<!-- field: x -->\n\nMore.",
		"list":    "- only\n- a list",
	} {
		original := yamlDoc(t, recipeWith("risky: "+yamlQuote(text)+"\n"))
		md, back := roundTrip(t, d, original)
		if diff := Diff(original, back); diff != "" {
			t.Fatalf("%s: round trip changed %s\n%s", name, diff, md)
		}
		if name != "list" && !strings.Contains(fieldSection(md, "risky"), "```yaml") {
			t.Errorf("%s: text that would break structure printed raw:\n%s", name, fieldSection(md, "risky"))
		}
	}
}

// TestCodeFieldPrintsInItsFence prints a declared bash field in one bash
// fence, a comment line included, and reads the command back (GH-2640).
func TestCodeFieldPrintsInItsFence(t *testing.T) {
	d := toyDeclaration(t)
	original := yamlDoc(t, recipeWith("method: |\n  boil water\n  # then season\n  echo done\n"))
	md, back := roundTrip(t, d, original)
	if diff := Diff(original, back); diff != "" {
		t.Fatalf("round trip changed %s\n%s", diff, md)
	}
	if section := fieldSection(md, "method"); section != "```bash\nboil water\n# then season\necho done\n```" {
		t.Errorf("method printed as:\n%s", section)
	}
	nonString := yamlDoc(t, recipeWith("method: [boil, serve]\n"))
	if _, err := d.Print(nonString); err == nil || !strings.Contains(err.Error(), "bash code field holds a string") {
		t.Errorf("a list in a code field printed: %v", err)
	}
	broken := strings.Replace(md, "```bash", "```sh", 1)
	if _, findings, _ := d.Parse([]byte(broken), toyName); len(findings) == 0 {
		t.Error("a code field in the wrong language raised no finding")
	}
}

func yamlQuote(text string) string {
	data, _ := yaml.Marshal(text)
	return strings.TrimRight(string(data), "\n")
}
