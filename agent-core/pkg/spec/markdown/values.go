// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/yuin/goldmark/ast"
	gtext "github.com/yuin/goldmark/text"
	"gopkg.in/yaml.v3"
)

// nestedKey matches an outline's nested item: "- <key>:" with a value or a
// nested list after it.
var nestedKey = regexp.MustCompile(`^- ([a-z_][a-z0-9_]*):(?: (.*))?$`)

// renderField converts a field's value under its declaration (eng10 Table 6):
// a keyed or outline field by its declared conversion, an empty list or table
// as nothing, and every other value by its own shape. A declared field whose
// value does not print in its declared shape is an error, not a fallback.
func (r *Declaration) renderField(slot, key string, value *yaml.Node) ([]string, string, error) {
	lines, err := r.renderUntypedForm(slot, key, value)
	if err != nil {
		return nil, "", err
	}
	field, typed := r.ruleFor(slot, key)
	if (typed && field.Shape != "") || classifyContent(strings.Join(lines, "\n")) != shapeYAML {
		return lines, "", nil
	}
	if form, typedLines, ok := typedFormFor(value); ok {
		return typedLines, form, nil
	}
	return lines, "", nil
}

// renderUntypedForm converts a field's value without a typed form: by its
// declared conversion, or by the first simple shape that holds it, else as
// a fenced yaml block.
func (r *Declaration) renderUntypedForm(slot, key string, value *yaml.Node) ([]string, error) {
	field, typed := r.ruleFor(slot, key)
	switch {
	case typed && field.Code != "":
		return codeLines(field.Code, value)
	case typed && field.Value == valueKeyed:
		return keyedLines(value)
	case typed && field.Value == valueOutline:
		return outlineLines(field, value)
	case typed && emptyCollection(value) && (field.Shape == shapeList || field.Shape == shapeTable):
		return nil, nil
	}
	lines, err := renderNode(value)
	if err != nil || !typed || field.Shape == "" || field.Shape == shapeEmpty {
		return lines, err
	}
	shape := classifyContent(strings.Join(lines, "\n"))
	if shape != field.Shape && (field.Shape != shapeProse || shape != shapeYAML) {
		return nil, fmt.Errorf("value prints as %s, declared %s", shape, field.Shape)
	}
	return lines, nil
}

// parseFieldForm reads a field in the typed form its marker names, or as an
// untyped field when it names none.
func (r *Declaration) parseFieldForm(slot, key, form, text string) (*yaml.Node, error) {
	if form != "" {
		return parseTyped(form, text)
	}
	return r.parseField(slot, key, text)
}

// ruleFor returns a field's rule where it applies: in the slots its in
// names, or anywhere when it names none. A nested field may share a name
// with a typed field elsewhere and is not held to that rule (GH-2657).
func (r *Declaration) ruleFor(slot, key string) (fieldRule, bool) {
	field, typed := r.Fields[key]
	return field, typed && (len(field.In) == 0 || contains(field.In, slot))
}

// markerForm returns the typed form a field marker names, unquoted.
func markerForm(m Marker) string {
	for _, a := range m.Attrs {
		if a.Key == attrForm {
			var form string
			if err := json.Unmarshal([]byte(a.Value), &form); err == nil {
				return form
			}
			return a.Value
		}
	}
	return ""
}

// parseField reads a field's Markdown back into the value renderField wrote.
func (r *Declaration) parseField(slot, key, text string) (*yaml.Node, error) {
	field, typed := r.ruleFor(slot, key)
	switch {
	case typed && field.Code != "":
		return parseCodeField(field.Code, text)
	case typed && field.Value == valueKeyed:
		return parseKeyed(text)
	case typed && field.Value == valueOutline:
		return parseOutline(field, text)
	case text == "" && typed && (field.Shape == shapeList || field.Shape == shapeTable):
		return &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}, nil
	case text == "":
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null"}, nil
	}
	return parseContent(text)
}

func emptyCollection(node *yaml.Node) bool {
	return (node.Kind == yaml.SequenceNode || node.Kind == yaml.MappingNode) && len(node.Content) == 0
}

// keyedLines prints a map of strings as a list of "name: text" items.
func keyedLines(value *yaml.Node) ([]string, error) {
	if value.Kind != yaml.MappingNode {
		return nil, errors.New("a keyed list holds a map")
	}
	items := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for i := 0; i+1 < len(value.Content); i += 2 {
		name, text := value.Content[i], value.Content[i+1]
		if !isString(text) || strings.ContainsAny(name.Value, ": \n") {
			return nil, fmt.Errorf("keyed entry %q is not a name and a text", name.Value)
		}
		items.Content = append(items.Content, stringNode(name.Value+": "+text.Value))
	}
	if len(items.Content) > 0 && !listSafe(items) {
		return nil, errors.New("a keyed entry is not prose-safe")
	}
	return listLines(items), nil
}

func parseKeyed(text string) (*yaml.Node, error) {
	keyed := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	if text == "" {
		return keyed, nil
	}
	for _, item := range parseList(text).Content {
		name, value, ok := strings.Cut(item.Value, ": ")
		if !ok {
			return nil, fmt.Errorf("keyed item %q has no name", firstLine(item.Value))
		}
		keyed.Content = append(keyed.Content, stringNode(name), stringNode(value))
	}
	return keyed, nil
}

// outlineLines prints a list of keyed records: each item reads id, colon,
// content, and every other record key is a nested item holding a line of
// text or a nested list of lines.
func outlineLines(field fieldRule, value *yaml.Node) ([]string, error) {
	if value.Kind != yaml.SequenceNode {
		return nil, errors.New("an outline holds a list")
	}
	m := nodeYAML{ID: field.ID, Content: field.ContentKey, Keyed: true}
	var lines []string
	for _, item := range value.Content {
		id, record, err := listRecord(m, item)
		if err != nil {
			return nil, err
		}
		itemLines, err := outlineItem(m, id, record)
		if err != nil {
			return nil, fmt.Errorf("outline item %s: %w", id, err)
		}
		lines = append(lines, itemLines...)
	}
	return lines, nil
}

func outlineItem(m nodeYAML, id string, record *yaml.Node) ([]string, error) {
	text := mapValue(record, m.Content)
	if text == nil || !isString(text) || !outlineSafe(id+": "+text.Value) {
		return nil, errors.New("content is not a line of prose-safe text")
	}
	lines := listLines(&yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{stringNode(id + ": " + text.Value)}})
	for i := 0; i+1 < len(record.Content); i += 2 {
		key, value := record.Content[i].Value, record.Content[i+1]
		if key == m.ID || key == m.Content {
			continue
		}
		nested, err := nestedLines(key, value)
		if err != nil {
			return nil, err
		}
		lines = append(lines, nested...)
	}
	return lines, nil
}

func nestedLines(key string, value *yaml.Node) ([]string, error) {
	switch {
	case isString(value) && cellSafe(value) && value.Value != "":
		return []string{"  - " + key + ": " + value.Value}, nil
	case isStringList(value):
		lines := []string{"  - " + key + ":"}
		for _, item := range value.Content {
			if !cellSafe(item) || item.Value == "" {
				return nil, fmt.Errorf("%s: item %q is not one line of text", key, item.Value)
			}
			lines = append(lines, "    - "+item.Value)
		}
		return lines, nil
	}
	return nil, fmt.Errorf("%s: an outline holds a line of text or a list of lines", key)
}

// outlineSafe reports text an outline item can hold: prose-safe, and no line
// that would read back as a nested item or a nested list.
func outlineSafe(text string) bool {
	if !proseSafe(text) {
		return false
	}
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimLeft(line, " ")
		if strings.HasPrefix(trimmed, "- ") || trimmed == "-" {
			return false
		}
	}
	return true
}

func parseOutline(field fieldRule, text string) (*yaml.Node, error) {
	outline := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	if text == "" {
		return outline, nil
	}
	for _, item := range parseList(text).Content {
		record, err := outlineRecord(field, item.Value)
		if err != nil {
			return nil, err
		}
		outline.Content = append(outline.Content, record)
	}
	return outline, nil
}

// outlineRecord reads one outline item: the keyed form when it has no nested
// items, the record form otherwise.
func outlineRecord(field fieldRule, value string) (*yaml.Node, error) {
	lines := strings.Split(value, "\n")
	split := len(lines)
	for i, line := range lines {
		if nestedKey.MatchString(line) {
			split = i
			break
		}
	}
	id, text, ok := strings.Cut(strings.Join(lines[:split], "\n"), ": ")
	if !ok {
		return nil, fmt.Errorf("outline item %q has no id", firstLine(value))
	}
	if split == len(lines) {
		return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map", Content: []*yaml.Node{stringNode(id), stringNode(text)}}, nil
	}
	record := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	setKey(record, field.ID, stringNode(id))
	setKey(record, field.ContentKey, stringNode(text))
	return record, addNested(record, lines[split:])
}

func addNested(record *yaml.Node, lines []string) error {
	var list *yaml.Node
	for _, line := range lines {
		if match := nestedKey.FindStringSubmatch(line); match != nil {
			list = nil
			if match[2] != "" {
				setKey(record, match[1], stringNode(match[2]))
				continue
			}
			list = &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
			setKey(record, match[1], list)
			continue
		}
		if list == nil || !strings.HasPrefix(line, "  - ") {
			return fmt.Errorf("outline line %q is not a nested item", line)
		}
		list.Content = append(list.Content, stringNode(strings.TrimPrefix(line, "  - ")))
	}
	return nil
}

// listRecord reads one entry of a record list: the keyed form ({<id>:
// <content>}) when the mapping allows it, the record form otherwise.
func listRecord(m nodeYAML, item *yaml.Node) (string, *yaml.Node, error) {
	if item.Kind != yaml.MappingNode {
		return "", nil, fmt.Errorf("a record is a map, found %s", item.ShortTag())
	}
	if m.Keyed && len(item.Content) == 2 && item.Content[0].Value != m.ID {
		record := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		setKey(record, m.Content, item.Content[1])
		return item.Content[0].Value, record, nil
	}
	id := mapValue(item, m.ID)
	if id == nil || id.Kind != yaml.ScalarNode || id.Value == "" {
		return "", nil, fmt.Errorf("a record has no %s", m.ID)
	}
	return id.Value, item, nil
}

// encodeNodeAttr writes a marker attribute value: a string raw when it reads
// back as the same string, every other value as JSON.
func encodeNodeAttr(node *yaml.Node) (string, error) {
	if isString(node) && rawAttrSafe(node.Value) && !strings.HasPrefix(node.Value, "-") {
		if back, err := decodeNodeAttr(node.Value); err == nil && isString(back) && back.Value == node.Value {
			return node.Value, nil
		}
	}
	data, err := nodeJSON(node)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(data, "|", `|`), nil
}

// decodeNodeAttr reads a marker attribute value. JSON is a subset of YAML's
// flow syntax, so the YAML decoder reads it with its key order; a value that
// is not JSON is a plain string.
func decodeNodeAttr(raw string) (*yaml.Node, error) {
	if raw == "" {
		return stringNode(""), nil
	}
	if !strings.ContainsAny(raw[:1], `"[{-0123456789`) && raw != "true" && raw != "false" && raw != "null" {
		return stringNode(raw), nil
	}
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(raw), &doc); err != nil || len(doc.Content) == 0 {
		return nil, fmt.Errorf("attribute value %q is not JSON", raw)
	}
	return doc.Content[0], nil
}

// nodeJSON writes a node as JSON, keeping map key order and scalar types.
func nodeJSON(node *yaml.Node) (string, error) {
	switch node.Kind {
	case yaml.ScalarNode:
		if node.ShortTag() == "!!str" {
			data, err := json.Marshal(node.Value)
			return string(data), err
		}
		var value any
		if err := node.Decode(&value); err != nil {
			return "", err
		}
		data, err := json.Marshal(value)
		return string(data), err
	case yaml.SequenceNode, yaml.MappingNode:
		return collectionJSON(node)
	}
	return "", fmt.Errorf("cannot write %s as JSON", node.ShortTag())
}

func collectionJSON(node *yaml.Node) (string, error) {
	open, closing, step := "[", "]", 1
	if node.Kind == yaml.MappingNode {
		open, closing, step = "{", "}", 2
	}
	parts := make([]string, 0, len(node.Content)/step)
	for i := 0; i < len(node.Content); i += step {
		value, err := nodeJSON(node.Content[i+step-1])
		if err != nil {
			return "", err
		}
		if step == 2 {
			key, _ := json.Marshal(node.Content[i].Value)
			value = string(key) + ": " + value
		}
		parts = append(parts, value)
	}
	return open + strings.Join(parts, ", ") + closing, nil
}

func stringNode(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value}
}

// mapValue returns the value under key in a mapping, or nil.
func mapValue(mapping *yaml.Node, key string) *yaml.Node {
	if mapping == nil || mapping.Kind != yaml.MappingNode || key == "" {
		return nil
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i+1]
		}
	}
	return nil
}

// setKey sets key in a mapping, replacing its value or appending the pair.
func setKey(mapping *yaml.Node, key string, value *yaml.Node) {
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			mapping.Content[i+1] = value
			return
		}
	}
	mapping.Content = append(mapping.Content, stringNode(key), value)
}

// parseCodeField reads a code field: its fence, or a form written before the
// declaration (plain text, or a fenced yaml block holding the string), which
// mage fmt rewrites into the fence.
func parseCodeField(language, text string) (*yaml.Node, error) {
	value, err := parseCode(language, text)
	switch {
	case err == nil:
		return value, nil
	case plainText(text):
		return stringNode(text), nil
	case classifyContent(text) == shapeYAML:
		if legacy, yamlErr := parseFencedYAML(text); yamlErr == nil && isString(legacy) {
			return legacy, nil
		}
	}
	return nil, err
}

// codeFieldReadable reports content a code field reads: its fence or a form
// mage fmt rewrites into it.
func codeFieldReadable(language, text string) bool {
	_, err := parseCodeField(language, text)
	return err == nil
}

// plainText reports field content that holds no fenced block, which a code
// field reads as its string so a document written before the declaration
// formats into the fence rather than failing.
func plainText(text string) bool {
	blocks, _ := topLevelBlocks([]byte(text))
	for _, b := range blocks {
		if _, fenced := b.Node.(*ast.FencedCodeBlock); fenced {
			return false
		}
	}
	return true
}

// codeLines prints a code field's string inside one fence of its declared
// language, longer than any backtick run in the text (GH-2640).
func codeLines(language string, value *yaml.Node) ([]string, error) {
	if !isString(value) {
		return nil, fmt.Errorf("a %s code field holds a string, found %s", language, value.ShortTag())
	}
	body := strings.TrimRight(value.Value, "\n")
	fence := "```"
	for strings.Contains(body, fence) {
		fence += "`"
	}
	lines := []string{fence + language}
	lines = append(lines, strings.Split(body, "\n")...)
	return append(lines, fence), nil
}

// parseCode reads a code field back: the content of its one fence, which
// must carry the declared language.
func parseCode(language, text string) (*yaml.Node, error) {
	src := []byte(text)
	doc := markdown.Parser().Parse(gtext.NewReader(src))
	fenced, ok := doc.FirstChild().(*ast.FencedCodeBlock)
	if !ok || fenced.NextSibling() != nil || string(fenced.Language(src)) != language {
		return nil, fmt.Errorf("a %s code field holds one fence of that language", language)
	}
	var body strings.Builder
	for i := 0; i < fenced.Lines().Len(); i++ {
		segment := fenced.Lines().At(i)
		body.Write(segment.Value(src))
	}
	return stringNode(strings.TrimRight(body.String(), "\n")), nil
}
