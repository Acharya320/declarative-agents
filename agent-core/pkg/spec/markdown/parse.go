// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Parse reads a Markdown document of this declaration's type back into the
// YAML tree Print writes it from, with any prose-containment findings: those
// concern wording, not structure. A document with a finding of tiers 1 to 3
// yields no tree; the findings say what to fix. A record's keys come back in canonical
// order: id, title, content, attributes, then its children in document order.
func (r *Declaration) Parse(data []byte, source string) (*yaml.Node, []Finding, error) {
	findings := r.Diagnose(data, source)
	for _, f := range findings {
		if f.Tier < 4 {
			return nil, findings, nil
		}
	}
	head, body := splitFrontmatter(data)
	root := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	if err := frontmatterInto(root, head); err != nil {
		return nil, nil, fmt.Errorf("%s: %w", source, err)
	}
	rd := &reader{d: r, root: root}
	blocks, lines := topLevelBlocks(body)
	rd.lines = lines
	for _, b := range blocks {
		if err := rd.step(b); err != nil {
			return nil, nil, fmt.Errorf("%s: %w", source, err)
		}
	}
	if err := rd.closeFrom(2); err != nil {
		return nil, nil, fmt.Errorf("%s: %w", source, err)
	}
	doc := &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{root}, HeadComment: strings.Join(rd.license, "\n")}
	return doc, findings, nil
}

// Format rewrites a Markdown document in canonical form: Parse, then Print.
func (r *Declaration) Format(data []byte, source string) ([]byte, []Finding, error) {
	doc, findings, err := r.Parse(data, source)
	if err != nil || doc == nil {
		return nil, findings, err
	}
	out, err := r.Print(doc)
	return out, findings, err
}

// frame is a field or node the reader has opened from a marker and heading.
type frame struct {
	kind      string
	field     bool
	key       string     // the field name, or the record id
	value     *yaml.Node // a field's value, or a node's record
	holder    *yaml.Node // the mapping a field belongs to
	container *yaml.Node // the list or map a record belongs to
	blocks    []Block
	form      string // a field's typed form, from its marker
}

type reader struct {
	d       *Declaration
	root    *yaml.Node
	lines   []string
	license []string
	pending *Marker
	stack   [maxHeadingLevel + 1]*frame
	open    *frame
}

func (rd *reader) step(b Block) error {
	if text, ok := commentLine(b, rd.lines); ok {
		m, isMarker, err := parseMarker(text)
		switch {
		case err != nil:
			return err
		case isMarker:
			rd.pending = &m
		default:
			rd.keepLicense(text)
		}
		return nil
	}
	if level, heading, ok := headingOf(b, rd.lines); ok {
		return rd.heading(level, heading)
	}
	if rd.open != nil {
		rd.open.blocks = append(rd.open.blocks, b)
	}
	return nil
}

func (rd *reader) keepLicense(comment string) {
	body := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(comment, markerOpen), markerClose))
	if isLicenseComment(body) {
		rd.license = append(rd.license, "# "+body)
	}
}

func (rd *reader) heading(level int, heading string) error {
	if err := rd.closeFrom(level); err != nil {
		return err
	}
	m := rd.pending
	rd.pending = nil
	parent := rd.parentAt(level)
	kind, ok := rd.d.nodeFor(m.Tag, level, parent)
	if !ok {
		return fmt.Errorf("%s at level %d under %s is not declared", m.Tag, level, parent)
	}
	f, err := rd.openFrame(level, kind, *m, heading)
	if err != nil {
		return err
	}
	f.kind = kind
	rd.stack[level] = f
	rd.open = f
	return nil
}

// openFrame starts a field or a record frame for the marker's tag.
func (rd *reader) openFrame(level int, kind string, m Marker, heading string) (*frame, error) {
	if m.Tag == tagField {
		f := rd.openField(level, m.Value)
		f.form = markerForm(m)
		return f, nil
	}
	return rd.openRecord(level, kind, m, heading)
}

func (rd *reader) parentAt(level int) string {
	if level == 2 {
		return "document"
	}
	switch parent := rd.stack[level-1]; {
	case parent == nil:
		return "none"
	case parent.field:
		return "field:" + parent.key
	default:
		return parent.kind
	}
}

// holder is the mapping a field at level belongs to: the document, or the
// record one level up.
func (rd *reader) holder(level int) *yaml.Node {
	if level == 2 || rd.stack[level-1] == nil {
		return rd.root
	}
	return rd.stack[level-1].value
}

func (rd *reader) openField(level int, key string) *frame {
	f := &frame{field: true, key: key, holder: rd.holder(level)}
	f.value = &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null"}
	if kind, ok := rd.d.fieldRecords(key); ok {
		f.value = newContainer(rd.d.Nodes[kind].YAML.Records)
	}
	setKey(f.holder, key, f.value)
	return f
}

func (rd *reader) openRecord(level int, kind string, m Marker, heading string) (*frame, error) {
	rule := rd.d.Nodes[kind]
	container := rd.recordContainer(level, rule)
	record := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	if rule.YAML.Records == recordsMap {
		setKey(container, m.Value, record)
	} else {
		container.Content = append(container.Content, record)
		setKey(record, rule.YAML.ID, stringNode(m.Value))
	}
	if title, ok := headingTitle(rule, m.Value, heading); ok {
		setKey(record, rule.YAML.Title, stringNode(title))
	}
	if rule.YAML.Content != "" {
		setKey(record, rule.YAML.Content, stringNode("")) // keeps content ahead of children
	}
	for _, a := range m.Attrs {
		value, err := decodeNodeAttr(a.Value)
		if err != nil {
			return nil, fmt.Errorf("%s %s: %w", kind, m.Value, err)
		}
		setKey(record, a.Key, value)
	}
	return &frame{key: m.Value, value: record, container: container}, nil
}

// recordContainer is the list or map a record joins: its field's, or the
// one under the key its YAML mapping names, created on first use.
func (rd *reader) recordContainer(level int, rule nodeRule) *yaml.Node {
	if parent := rd.stack[level-1]; level > 2 && parent != nil && parent.field {
		return parent.value
	}
	holder := rd.holder(level)
	if existing := mapValue(holder, rule.YAML.From); existing != nil {
		return existing
	}
	container := newContainer(rule.YAML.Records)
	setKey(holder, rule.YAML.From, container)
	return container
}

func headingTitle(rule nodeRule, id, heading string) (string, bool) {
	switch rule.Heading {
	case headingIDTitle:
		return strings.CutPrefix(heading, id+" - ")
	case headingFree:
		return heading, rule.YAML.Title != ""
	}
	return "", false
}

// closeFrom closes every frame at level or deeper, deepest first.
func (rd *reader) closeFrom(level int) error {
	for l := maxHeadingLevel; l >= level; l-- {
		if f := rd.stack[l]; f != nil {
			if err := rd.close(f); err != nil {
				return err
			}
			rd.stack[l] = nil
		}
	}
	rd.open = nil
	return nil
}

func (rd *reader) close(f *frame) error {
	text := ""
	if len(f.blocks) > 0 {
		text = spanText(rd.lines, f.blocks[0].Start, f.blocks[len(f.blocks)-1].End)
	}
	if f.field {
		if _, records := rd.d.fieldRecords(f.key); records {
			return nil
		}
		value, err := rd.d.parseFieldForm(f.key, f.form, text)
		if err != nil {
			return fmt.Errorf("field %s: %w", f.key, err)
		}
		setKey(f.holder, f.key, value)
		return nil
	}
	return rd.closeRecord(f, text)
}

// closeRecord fills in the record's content and, when the record holds only
// its id and content and the mapping allows it, rewrites it in keyed form.
func (rd *reader) closeRecord(f *frame, text string) error {
	m := rd.d.Nodes[f.kind].YAML
	if m.Content != "" {
		if err := setContent(f.value, m.Content, text); err != nil {
			return fmt.Errorf("%s %s: %w", f.kind, f.key, err)
		}
	}
	if !m.Keyed || m.Records != recordsList || len(f.value.Content) != 4 || mapValue(f.value, m.Content) == nil {
		return nil
	}
	keyed := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map",
		Content: []*yaml.Node{stringNode(f.key), mapValue(f.value, m.Content)}}
	for i, item := range f.container.Content {
		if item == f.value {
			f.container.Content[i] = keyed
		}
	}
	return nil
}

// setContent stores a node's content: prose as a string, a fenced yaml
// block as the value it holds, and nothing when the node has no content.
func setContent(record *yaml.Node, key, text string) error {
	switch {
	case text == "":
		deleteKey(record, key)
	case classifyContent(text) == shapeYAML:
		value, err := parseFencedYAML(text)
		if err != nil {
			return err
		}
		setKey(record, key, value)
	default:
		setKey(record, key, stringNode(text))
	}
	return nil
}

func newContainer(records string) *yaml.Node {
	if records == recordsMap {
		return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	}
	return &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
}

func deleteKey(mapping *yaml.Node, key string) {
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			mapping.Content = append(mapping.Content[:i], mapping.Content[i+2:]...)
			return
		}
	}
}

// splitFrontmatter separates the frontmatter text from the body; Parse runs
// it only on documents the diagnostics passed, so both delimiters exist.
func splitFrontmatter(data []byte) (string, []byte) {
	text := string(data)
	closing := strings.Index(text[4:], "\n---\n")
	// Keep the newline before the closing delimiter: a block scalar that
	// ends the frontmatter keeps its own final newline only with it.
	return text[4 : 4+closing+1], []byte(text[4+closing+5:])
}

func frontmatterInto(root *yaml.Node, head string) error {
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(head), &doc); err != nil {
		return fmt.Errorf("frontmatter: %w", err)
	}
	if len(doc.Content) > 0 && doc.Content[0].Kind == yaml.MappingNode {
		root.Content = append(root.Content, doc.Content[0].Content...)
	}
	return nil
}
