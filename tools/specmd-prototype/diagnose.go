// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package specmd

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"gopkg.in/yaml.v3"
)

// A Finding is one diagnostic: a stable kebab-case name, where it is, what
// the rule expected and what the document has, the rule it enforces, and how
// to fix it. Tier 0 (canonical form) is never a finding: the formatter fixes
// it.
type Finding struct {
	Name     string
	Tier     int
	File     string
	Line     int
	Expected string
	Found    string
	Rule     string
	Hint     string
}

func (f Finding) String() string {
	return fmt.Sprintf("%s:%d: %s: expected %s, found %s [rules/srd.yaml %s]; %s",
		f.File, f.Line, f.Name, f.Expected, f.Found, f.Rule, f.Hint)
}

// dnode is a node the diagnostics have opened from a marker and heading.
type dnode struct {
	kind       string
	id         string
	line       int
	markerLine int
	fields     map[string]bool
	content    []block
}

// diagnosis walks one document collecting findings; it never stops early.
type diagnosis struct {
	rules    *rules
	file     string
	lines    []string
	offset   int
	findings []Finding
	pending  *marker
	pendLine int
	l2, l3   *dnode
	open     *dnode
	ids      map[string]int
	docKeys  map[string]bool
	groups   int
	status   string
}

// findingCatalog is every finding name the diagnostics can report, by tier.
var findingCatalog = map[string]int{
	"frontmatter-missing": 1, "frontmatter-invalid": 1, "frontmatter-missing-key": 1,
	"frontmatter-unknown-key": 1, "frontmatter-id-mismatch": 1, "marker-malformed": 1,
	"marker-unknown-tag": 1, "marker-unknown-attribute": 1, "marker-orphan": 1,
	"heading-untagged": 2, "heading-misplaced": 2, "id-malformed": 2, "item-wrong-group": 2,
	"heading-id-mismatch": 2, "duplicate-id": 2, "field-disallowed": 2, "field-duplicate": 2,
	"content-untagged": 2, "untagged-structure": 2, "missing-content": 2,
	"missing-required-field": 2, "missing-requirement-groups": 2,
	"field-shape-mismatch": 3, "list-item-prefix": 3, "table-columns-mismatch": 3,
	"attribute-shape": 3, "criterion-traces-missing": 3, "criterion-trace-invalid": 3,
	"prose-containment": 4,
}

// Diagnose reports every finding in an SRD Markdown document.
func Diagnose(data []byte, source string) []Finding {
	d := &diagnosis{rules: srdRules, file: source, ids: map[string]int{}, docKeys: map[string]bool{}}
	body, offset := d.frontmatter(data)
	d.offset = offset
	blocks, lines := topLevelBlocks(body)
	d.lines = lines
	for _, b := range blocks {
		d.step(b)
	}
	d.closeOpen()
	if d.pending != nil {
		d.orphan()
	}
	d.documentLevel()
	sort.SliceStable(d.findings, func(i, j int) bool { return d.findings[i].Line < d.findings[j].Line })
	return d.findings
}

func (d *diagnosis) add(name string, tier, line int, expected, found, rule, hint string) {
	d.findings = append(d.findings, Finding{
		Name: name, Tier: tier, File: d.file, Line: line,
		Expected: expected, Found: found, Rule: rule, Hint: hint,
	})
}

func (d *diagnosis) fileLine(b block) int { return b.start + d.offset + 1 }

func (d *diagnosis) step(b block) {
	if text, ok := commentLine(b, d.lines); ok {
		d.comment(b, text)
		return
	}
	if level, heading, ok := headingOf(b, d.lines); ok {
		d.heading(b, level, heading)
		return
	}
	if d.open == nil {
		d.add("content-untagged", 2, d.fileLine(b), "content under a tagged heading",
			"content before any marker and heading", "nodes",
			"put a marker and heading above it, or move it under the section it belongs to")
		return
	}
	d.open.content = append(d.open.content, b)
}

func (d *diagnosis) comment(b block, text string) {
	m, isMarker, err := parseMarker(text)
	switch {
	case err != nil:
		d.add("marker-malformed", 1, d.fileLine(b), "<!-- <tag>: <value> | <key>: <value> -->",
			text, "nodes.*.tag", "write the tag and its value separated by a colon, attributes after \" | \"")
		return
	case !isMarker:
		return
	}
	if !knownTag(d.rules, m.Tag) {
		d.add("marker-unknown-tag", 1, d.fileLine(b), "a tag the rules declare (ac, field, req)",
			m.Tag, "nodes.*.tag", suggest(m.Tag, []string{"ac", "field", "req"}))
		return
	}
	if d.pending != nil {
		d.orphan()
	}
	d.pending, d.pendLine = &m, d.fileLine(b)
}

func knownTag(r *rules, tag string) bool {
	for _, node := range r.Nodes {
		if node.Tag == tag {
			return true
		}
	}
	return false
}

func (d *diagnosis) orphan() {
	d.add("marker-orphan", 1, d.pendLine, "a heading directly below the marker",
		"no heading before the next marker or the end of the document", "nodes",
		"add the heading the marker tags, or delete the marker")
	d.pending = nil
}

// suggest names the nearest declared value, when one is close.
func suggest(found string, known []string) string {
	best, bestDistance := "", len(found)/2+2
	for _, candidate := range known {
		if distance := editDistance(found, candidate); distance < bestDistance {
			best, bestDistance = candidate, distance
		}
	}
	if best == "" {
		return "use one of " + strings.Join(known, ", ")
	}
	return fmt.Sprintf("did you mean %q?", best)
}

func editDistance(a, b string) int {
	previous := make([]int, len(b)+1)
	for j := range previous {
		previous[j] = j
	}
	for i := 1; i <= len(a); i++ {
		current := make([]int, len(b)+1)
		current[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			current[j] = min(previous[j]+1, current[j-1]+1, previous[j-1]+cost)
		}
		previous = current
	}
	return previous[len(b)]
}

// frontmatter checks the frontmatter and returns the body and its line
// offset in the file.
func (d *diagnosis) frontmatter(data []byte) ([]byte, int) {
	text := string(data)
	if !strings.HasPrefix(text, "---\n") {
		d.add("frontmatter-missing", 1, 1, "a \"---\" frontmatter block first", "no frontmatter",
			"frontmatter", "start the file with ---, id and title, then ---")
		return data, 0
	}
	closing := strings.Index(text[4:], "\n---\n")
	if closing < 0 {
		d.add("frontmatter-missing", 1, 1, "a closing \"---\"", "an unterminated frontmatter block",
			"frontmatter", "close the frontmatter with a line holding only ---")
		return data, 0
	}
	head := text[4 : 4+closing]
	d.frontmatterKeys(head)
	return []byte(text[4+closing+5:]), strings.Count(head, "\n") + 3
}

func (d *diagnosis) frontmatterKeys(head string) {
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(head), &doc); err != nil || len(doc.Content) == 0 ||
		doc.Content[0].Kind != yaml.MappingNode {
		d.add("frontmatter-invalid", 1, 2, "a YAML mapping", "text that is not a YAML mapping",
			"frontmatter", "write the frontmatter as key: value lines")
		return
	}
	root := doc.Content[0]
	seen := map[string]string{}
	for i := 0; i+1 < len(root.Content); i += 2 {
		key := root.Content[i]
		seen[key.Value] = root.Content[i+1].Value
		if !contains(d.rules.Frontmatter.Allowed, key.Value) {
			d.add("frontmatter-unknown-key", 1, key.Line+1, "one of "+strings.Join(d.rules.Frontmatter.Allowed, ", "),
				key.Value, "frontmatter.allowed", suggest(key.Value, d.rules.Frontmatter.Allowed))
		}
	}
	d.status = seen["status"]
	for _, key := range d.rules.Frontmatter.Required {
		if _, ok := seen[key]; !ok {
			d.add("frontmatter-missing-key", 1, 2, "frontmatter key "+key, "none", "frontmatter.required",
				"add \""+key+": ...\" to the frontmatter")
		}
	}
	want := strings.TrimSuffix(filepath.Base(d.file), filepath.Ext(d.file))
	if id, ok := seen["id"]; ok && d.rules.Frontmatter.IDMatchesFilename && id != want {
		d.add("frontmatter-id-mismatch", 1, frontmatterLine(root, "id"), "id "+want+" (the file name)", id,
			"frontmatter.id_matches_filename", "rename the file or change the id so they agree")
	}
}

func frontmatterLine(root *yaml.Node, key string) int {
	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Value == key {
			return root.Content[i].Line + 1
		}
	}
	return 2
}

func contains(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

// documentLevel checks what the whole document must carry.
func (d *diagnosis) documentLevel() {
	if contains(d.rules.RequiredFieldsExemptStatus, d.status) {
		return
	}
	end := len(d.lines) + d.offset
	if len(d.lines) > 0 && d.lines[len(d.lines)-1] == "" {
		end--
	}
	for _, key := range d.rules.RequiredFields {
		if !d.docKeys[key] {
			d.add("missing-required-field", 2, end, "a "+key+" field", "none", "required_fields",
				"add <!-- field: "+key+" --> with a level-2 heading and its content")
		}
	}
	if d.rules.RequiresRequirementGroups && d.groups == 0 {
		d.add("missing-requirement-groups", 2, end, "at least one requirement group", "none",
			"requires_requirement_groups", "add <!-- req: R1 --> above \"## R1 - <title>\"")
	}
}

// containsStructure reports whether any block is a table or a list: structure a
// prose node must carry under its own tagged field.
func containsStructure(blocks []block) (block, string, bool) {
	for _, b := range blocks {
		switch b.node.(type) {
		case *east.Table:
			return b, "table", true
		case *ast.List:
			return b, "list", true
		}
	}
	return block{}, "", false
}
