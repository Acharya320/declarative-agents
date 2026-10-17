// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"path"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Reference links (eng10, Reference links; GH-2665). A field whose rule
// declares links renders each reference it holds as a Markdown link: a bare
// item of the same document to its heading anchor, a document id to the
// document's file, and each item after a document id to that file's anchor.
// Links are rendering only: the parser reads a link back as its text.

var (
	// referenceValue is a whole reference: a document id with optional
	// items, or a bare item.
	referenceValue = regexp.MustCompile(`^(?:([a-z0-9][a-z0-9.-]*[a-z0-9])(?: ((?:[A-Z]+[0-9]+(?:\.[0-9]+)*[a-z]?)(?:, [A-Z]+[0-9]+(?:\.[0-9]+)*[a-z]?)*))?|([A-Z]+[0-9]+(?:\.[0-9]+)*[a-z]?))$`)
	itemToken      = regexp.MustCompile(`^[A-Z]+[0-9]+(?:\.[0-9]+)*[a-z]?$`)
	idToken        = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`)
	markdownLink   = regexp.MustCompile(`\[([^\[\]]+)\]\(([^()\s]*)\)`)
)

// LinkIndex holds what cross-document links need: each corpus document's
// path from the repository root and its items' heading anchors.
type LinkIndex map[string]LinkTarget

// LinkTarget is one document in the link index.
type LinkTarget struct {
	Path    string
	Anchors map[string]string // item id -> heading anchor
}

// Add records a corpus document under its id.
func (ix LinkIndex) Add(id, docPath string, data []byte) {
	ix[id] = LinkTarget{Path: docPath, Anchors: pageAnchors(strings.Split(string(data), "\n"))}
}

// LinkIndexes holds one link index per corpus, keyed by the corpus root: a
// reference links only within its own corpus, so every tool that prints a
// corpus document agrees on its links whatever else it has loaded.
type LinkIndexes map[string]LinkIndex

// Add records a corpus document under its frontmatter id in its corpus's
// index; a document without an id is not a link target.
func (ix LinkIndexes) Add(docPath string, data []byte) {
	id := DocumentID(data)
	if id == "" {
		return
	}
	root := CorpusRoot(docPath)
	if ix[root] == nil {
		ix[root] = LinkIndex{}
	}
	ix[root].Add(id, docPath, data)
}

// For returns the index of the corpus a document belongs to.
func (ix LinkIndexes) For(docPath string) LinkIndex { return ix[CorpusRoot(docPath)] }

// CorpusRoot is the directory holding a corpus document's docs tree: "" for
// docs/..., agent-core/ for agent-core/docs/..., in slash form.
func CorpusRoot(docPath string) string {
	if strings.HasPrefix(docPath, "docs/") {
		return ""
	}
	if i := strings.Index(docPath, "/docs/"); i >= 0 {
		return docPath[:i+1]
	}
	return ""
}

// DocumentID reads a document's frontmatter id, or "" when it has none.
func DocumentID(data []byte) string {
	text := string(data)
	if !strings.HasPrefix(text, "---\n") {
		return ""
	}
	head, _, found := strings.Cut(text[4:], "\n---\n")
	if !found {
		return ""
	}
	for _, line := range strings.Split(head, "\n") {
		if id, ok := strings.CutPrefix(line, "id: "); ok {
			return strings.Trim(strings.TrimSpace(id), `"'`)
		}
	}
	return ""
}

// linkContext is what the printer links against: the document's own path
// and item anchors, and the corpus index when it has one.
type linkContext struct {
	path    string
	anchors map[string]string
	index   LinkIndex
}

// slug is GitHub's heading anchor: lowercase, every character but letters,
// digits, underscores, hyphens, and spaces dropped, spaces to hyphens.
func slug(heading string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(heading) {
		switch {
		case r == ' ':
			b.WriteRune('-')
		case r == '-' || r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsMark(r):
			b.WriteRune(r)
		}
	}
	return b.String()
}

// pageAnchors maps each item heading (a heading under a marker whose tag is
// not field) to its anchor, counting repeats over every heading on the page
// as GitHub does.
func pageAnchors(lines []string) map[string]string {
	anchors, seen := map[string]string{}, map[string]int{}
	item, fenced := "", false
	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			fenced = !fenced
			continue
		}
		if fenced {
			continue
		}
		if m, ok, err := parseMarker(line); ok && err == nil {
			item = ""
			if m.Tag != tagField {
				item = m.Value
			}
			continue
		}
		level := len(line) - len(strings.TrimLeft(line, "#"))
		if level < 1 || level > maxHeadingLevel || !strings.HasPrefix(line[level:], " ") {
			continue
		}
		base := slug(strings.TrimSpace(line[level:]))
		anchor := base
		if n := seen[base]; n > 0 {
			anchor = base + "-" + strconv.Itoa(n)
		}
		seen[base]++
		if item != "" {
			if _, ok := anchors[item]; !ok {
				anchors[item] = anchor
			}
		}
		item = ""
	}
	return anchors
}

// link renders one reference value as links, or returns it unchanged when
// it is not a reference or names nothing the context knows.
func (c *linkContext) link(value string) string {
	m := referenceValue.FindStringSubmatch(value)
	switch {
	case c == nil || m == nil:
		return value
	case m[3] != "": // a bare item on this page
		if anchor, ok := c.anchors[m[3]]; ok {
			return "[" + m[3] + "](#" + anchor + ")"
		}
		return value
	}
	target, ok := c.index[m[1]]
	if !ok {
		return value
	}
	file := relativePath(c.path, target.Path)
	out := "[" + m[1] + "](" + file + ")"
	if m[2] == "" {
		return out
	}
	items := strings.Split(m[2], ", ")
	for i, item := range items {
		if anchor, ok := target.Anchors[item]; ok {
			items[i] = "[" + item + "](" + file + "#" + anchor + ")"
		}
	}
	return out + " " + strings.Join(items, ", ")
}

// linkCell links a table cell: one reference, or each reference in a
// typed-table cell holding a flow list of them, [a, b], which the parser
// unlinks before it decodes the flow value (GH-2669).
func (c *linkContext) linkCell(cell string) string {
	if len(cell) < 3 || cell[0] != '[' || cell[len(cell)-1] != ']' {
		return c.link(cell)
	}
	items := strings.Split(cell[1:len(cell)-1], ", ")
	for i, item := range items {
		items[i] = c.link(item)
	}
	return "[" + strings.Join(items, ", ") + "]"
}

// relativePath is target's path from the directory holding from, with
// forward slashes.
func relativePath(from, target string) string {
	fromParts := strings.Split(path.Dir(from), "/")
	toParts := strings.Split(target, "/")
	i := 0
	for i < len(fromParts) && i < len(toParts)-1 && fromParts[i] == toParts[i] {
		i++
	}
	var up []string
	for range fromParts[i:] {
		if fromParts[0] != "." {
			up = append(up, "..")
		}
	}
	return strings.Join(append(up, toParts[i:]...), "/")
}

// unlink reads every generated link in text back as its text: a link whose
// text is a document id and whose file is named after it, or whose text is
// an item and whose anchor starts with the item's own anchor. Any other
// link stays as written.
func unlink(text string) string {
	if !strings.Contains(text, "](") {
		return text
	}
	return markdownLink.ReplaceAllStringFunc(text, func(link string) string {
		m := markdownLink.FindStringSubmatch(link)
		label, target := m[1], m[2]
		file, anchor, _ := strings.Cut(target, "#")
		switch {
		case idToken.MatchString(label) && anchor == "" && strings.TrimSuffix(path.Base(file), ".md") == label:
			return label
		case itemToken.MatchString(label) && anchor != "" && strings.HasPrefix(anchor, slug(label)):
			return label
		}
		return link
	})
}

// linkLines renders the references a field holds: each list item, a
// one-line value, or the cells of the declared table columns.
func (c *linkContext) linkLines(field fieldRule, lines []string) []string {
	if c == nil || (!field.Links && len(field.LinkColumns) == 0) {
		return lines
	}
	out := append([]string(nil), lines...)
	switch {
	case len(field.LinkColumns) > 0 && len(out) > 2 && strings.HasPrefix(out[0], "|"):
		header := splitRow(out[0])
		for i := 2; i < len(out); i++ {
			cells := splitRow(out[i])
			for j, column := range header {
				if j < len(cells) && contains(field.LinkColumns, column) {
					cells[j] = c.linkCell(cells[j])
				}
			}
			out[i] = tableRow(cells)
		}
	case field.Links && len(out) == 1 && !strings.HasPrefix(out[0], "- "):
		out[0] = c.link(out[0])
	case field.Links:
		for i, line := range out {
			if item, ok := strings.CutPrefix(line, "- "); ok {
				out[i] = "- " + c.link(item)
			}
		}
	}
	return out
}
