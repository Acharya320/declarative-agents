// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Print writes a YAML document as canonical Markdown under its declaration.
// The frontmatter takes the top-level keys the declaration allows there, the
// license lines of the document's head comment become comment lines, and
// every other key becomes a field or a run of nodes, in document order. No
// branch depends on the document type: the declaration and its YAML mapping
// decide everything.
func (d *Declaration) Print(doc *yaml.Node) ([]byte, error) {
	out, _, err := d.PrintReport(doc)
	return out, err
}

func (p *printer) document(doc *yaml.Node) ([]byte, error) {
	d := p.d
	root, license := documentRoot(doc)
	if root == nil || root.Kind != yaml.MappingNode {
		return nil, errors.New("print: the document is not a mapping")
	}
	if err := p.frontmatter(root); err != nil {
		return nil, err
	}
	p.lines = append(p.lines, license...)
	for i := 0; i+1 < len(root.Content); i += 2 {
		key := root.Content[i].Value
		if contains(d.Frontmatter.Allowed, key) {
			continue
		}
		if err := p.member(2, "document", key, root.Content[i+1]); err != nil {
			return nil, fmt.Errorf("print: %w", err)
		}
	}
	return []byte(strings.Join(p.lines, "\n") + "\n"), nil
}

type printer struct {
	d      *Declaration
	lines  []string
	fences []Fence
}

// Fence names a value the printer still wrote as a fenced yaml block, the
// reported last resort (eng10), and why no other rendering held it.
type Fence struct {
	Field  string
	Reason string
}

func (f Fence) String() string { return f.Field + ": " + f.Reason }

// PrintReport prints a document like Print and also returns every value it
// wrote as a fenced yaml block (GH-2643).
func (d *Declaration) PrintReport(doc *yaml.Node) ([]byte, []Fence, error) {
	p := &printer{d: d}
	out, err := p.document(doc)
	return out, p.fences, err
}

// Fences reports the fenced yaml blocks a Markdown document prints with.
func (d *Declaration) Fences(data []byte, source string) ([]Fence, error) {
	doc, findings, err := d.Parse(data, source)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, fmt.Errorf("%s: %d finding(s)", source, len(findings))
	}
	_, fences, err := d.PrintReport(doc)
	return fences, err
}

// block appends lines as one Markdown block, a blank line after the last.
func (p *printer) block(lines ...string) {
	if len(p.lines) > 0 {
		p.lines = append(p.lines, "")
	}
	p.lines = append(p.lines, lines...)
}

func (p *printer) frontmatter(root *yaml.Node) error {
	head := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	for i := 0; i+1 < len(root.Content); i += 2 {
		if contains(p.d.Frontmatter.Allowed, root.Content[i].Value) {
			head.Content = append(head.Content, root.Content[i], root.Content[i+1])
		}
	}
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(explicitNulls(head)); err != nil {
		return fmt.Errorf("print frontmatter: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("print frontmatter: %w", err)
	}
	body := strings.TrimRight(buf.String(), "\n")
	if len(head.Content) == 0 {
		body = "{}"
	}
	p.lines = append(p.lines, "---", body, "---")
	return nil
}

// member prints one key of a record (or of the document): a run of nodes
// when a node kind takes its records straight from the key, a field
// otherwise.
func (p *printer) member(level int, parent, key string, value *yaml.Node) error {
	if kind, ok := p.d.directRecords(level, parent, key); ok {
		return p.records(kind, value)
	}
	if _, ok := p.d.fieldKind(level, parent); !ok {
		return fmt.Errorf("%s: no field is declared at level %d under %s", key, level, parent)
	}
	if kind, ok := p.d.fieldRecords(key); ok {
		p.block(Marker{Tag: tagField, Value: key}.String(), strings.Repeat("#", level)+" "+p.d.fieldHeading(key))
		return p.records(kind, value)
	}
	if value.Kind == yaml.MappingNode && len(value.Content) > 0 && p.d.nestedSlot(level, key) {
		return p.nested(level, key, value)
	}
	lines, form, err := p.d.renderField(key, value)
	if err != nil {
		return fmt.Errorf("field %s: %w", key, err)
	}
	if form == "" && classifyContent(strings.Join(lines, "\n")) == shapeYAML {
		p.fences = append(p.fences, Fence{Field: key, Reason: fenceReason(value)})
	}
	marker := Marker{Tag: tagField, Value: key}
	if form != "" {
		marker.Attrs = []Attr{{Key: attrForm, Value: jsonString(form)}}
	}
	p.block(marker.String(), strings.Repeat("#", level)+" "+p.d.fieldHeading(key))
	if len(lines) > 0 {
		p.block(lines...)
	}
	return nil
}

// nested prints a field whose value is a map as its heading, then each key
// as a nested field one level down (GH-2653).
func (p *printer) nested(level int, key string, value *yaml.Node) error {
	p.block(Marker{Tag: tagField, Value: key}.String(), strings.Repeat("#", level)+" "+p.d.fieldHeading(key))
	for i := 0; i+1 < len(value.Content); i += 2 {
		if err := p.member(level+1, "field:"+key, value.Content[i].Value, value.Content[i+1]); err != nil {
			return fmt.Errorf("field %s: %w", key, err)
		}
	}
	return nil
}

// records prints each record of a node kind, from a list or from a map
// keyed by id.
func (p *printer) records(kind string, value *yaml.Node) error {
	m := p.d.Nodes[kind].YAML
	switch {
	case m.Records == recordsMap && value.Kind == yaml.MappingNode:
		for i := 0; i+1 < len(value.Content); i += 2 {
			if err := p.record(kind, value.Content[i].Value, value.Content[i+1]); err != nil {
				return err
			}
		}
	case m.Records == recordsList && value.Kind == yaml.SequenceNode:
		for _, item := range value.Content {
			id, record, err := listRecord(m, item)
			if err != nil {
				return fmt.Errorf("%s: %w", kind, err)
			}
			if err := p.record(kind, id, record); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("%s: expected %s records, found %s", kind, m.Records, value.ShortTag())
	}
	return nil
}

// record prints one node: its marker with attributes, its heading, its
// prose content, then every other record key as a child, in record order.
func (p *printer) record(kind, id string, record *yaml.Node) error {
	rule := p.d.Nodes[kind]
	if record.Kind != yaml.MappingNode {
		return fmt.Errorf("%s %s: a record is a map", kind, id)
	}
	marker, err := recordMarker(rule, id, record)
	if err != nil {
		return fmt.Errorf("%s %s: %w", kind, id, err)
	}
	heading, err := recordHeading(rule, id, record)
	if err != nil {
		return fmt.Errorf("%s %s: %w", kind, id, err)
	}
	p.block(marker.String(), strings.Repeat("#", rule.Level)+" "+heading)
	if content := mapValue(record, rule.YAML.Content); content != nil {
		lines, err := contentLines(content, rule.Content == "any")
		if err != nil {
			return fmt.Errorf("%s %s: %w", kind, id, err)
		}
		if classifyContent(strings.Join(lines, "\n")) == shapeYAML {
			p.fences = append(p.fences, Fence{Field: kind + " " + id, Reason: "node text that would break the document structure"})
		}
		p.block(lines...)
	}
	for i := 0; i+1 < len(record.Content); i += 2 {
		key := record.Content[i].Value
		if consumedKey(rule, key) {
			continue
		}
		if err := p.member(rule.Level+1, kind, key, record.Content[i+1]); err != nil {
			return fmt.Errorf("%s %s: %w", kind, id, err)
		}
	}
	return nil
}

// contentLines prints a node's content as prose, as raw Markdown when the
// node's content may be any and the text reads back exactly, or else as a
// fenced yaml block.
func contentLines(content *yaml.Node, anyContent bool) ([]string, error) {
	switch {
	case isString(content) && (proseSafe(content.Value) || proseTextSafe(content.Value)):
		return proseLines(content.Value), nil
	case anyContent && isString(content) && rawSafe(content.Value):
		return proseLines(content.Value), nil
	}
	return fencedYAML(content)
}

// rawSafe reports Markdown text that reads back as itself under a node:
// no heading or comment line that would open or tag a node, no fenced yaml
// block standing alone (it would read back as a value), and no surrounding
// whitespace Markdown would drop.
func rawSafe(text string) bool {
	trimmed := strings.TrimRight(text, "\n")
	if trimmed == "" || strings.TrimSpace(trimmed) != trimmed || classifyContent(trimmed) == shapeYAML {
		return false
	}
	blocks, lines := topLevelBlocks([]byte(trimmed))
	for _, b := range blocks {
		if _, _, heading := headingOf(b, lines); heading {
			return false
		}
		if _, comment := commentLine(b, lines); comment {
			return false
		}
	}
	return len(blocks) > 0 && spanText(lines, blocks[0].Start, blocks[len(blocks)-1].End) == trimmed
}

// fieldTextSafe reports a field's text that reads back as itself: raw-safe,
// and not a lone list or table, which would read back as a list or a table.
// Fenced blocks such as Mermaid, block quotes, and indented lines stay in the
// text, so they render (GH-2640).
func fieldTextSafe(text string) bool {
	return rawSafe(text) && classifyContent(strings.TrimRight(text, "\n")) == shapeProse
}

// proseTextSafe reports prose node text that reads back as itself and keeps
// the prose content rule: raw-safe, with no untagged list or table.
func proseTextSafe(text string) bool {
	if !rawSafe(text) {
		return false
	}
	blocks, _ := topLevelBlocks([]byte(strings.TrimRight(text, "\n")))
	_, _, structured := containsStructure(blocks)
	return !structured
}

func recordMarker(rule nodeRule, id string, record *yaml.Node) (Marker, error) {
	marker := Marker{Tag: rule.Tag, Value: id}
	for i := 0; i+1 < len(record.Content); i += 2 {
		key := record.Content[i].Value
		if !contains(rule.Attributes, key) {
			continue
		}
		value, err := encodeNodeAttr(record.Content[i+1])
		if err != nil {
			return Marker{}, fmt.Errorf("attribute %s: %w", key, err)
		}
		marker.Attrs = append(marker.Attrs, Attr{Key: key, Value: value})
	}
	return marker, nil
}

func recordHeading(rule nodeRule, id string, record *yaml.Node) (string, error) {
	title := mapValue(record, rule.YAML.Title)
	if title != nil && (!isString(title) || !headingSafe(title.Value)) {
		return "", fmt.Errorf("%s is not one line of text", rule.YAML.Title)
	}
	switch {
	case rule.Heading == headingIDTitle && title != nil:
		return id + " - " + title.Value, nil
	case rule.Heading == headingFree && title != nil:
		return title.Value, nil
	case rule.Heading == headingFree:
		return "", fmt.Errorf("a free heading needs its %s", rule.YAML.Title)
	}
	return id, nil
}

func headingSafe(text string) bool {
	return text != "" && !strings.ContainsAny(text, "\n") && strings.TrimSpace(text) == text &&
		!strings.HasSuffix(text, "#")
}

// consumedKey reports the record keys the marker, heading, or content
// already printed.
func consumedKey(rule nodeRule, key string) bool {
	m := rule.YAML
	switch {
	case key == m.Content:
		return true
	case key == m.ID && m.Records == recordsList:
		return true
	case key == m.Title && rule.Heading != headingID:
		return true
	}
	return contains(rule.Attributes, key)
}

// documentRoot unwraps a document node and returns its license lines as
// Markdown comment lines.
func documentRoot(doc *yaml.Node) (*yaml.Node, []string) {
	if doc == nil {
		return nil, nil
	}
	root, comments := doc, []string{doc.HeadComment}
	if doc.Kind == yaml.DocumentNode && len(doc.Content) > 0 {
		root = doc.Content[0]
		comments = append(comments, root.HeadComment)
		if root.Kind == yaml.MappingNode && len(root.Content) > 0 {
			comments = append(comments, root.Content[0].HeadComment)
		}
	}
	var license []string
	for _, line := range strings.Split(strings.Join(comments, "\n"), "\n") {
		text := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "#"))
		if isLicenseComment(text) {
			license = append(license, markerOpen+" "+text+" "+markerClose)
		}
	}
	return root, license
}
