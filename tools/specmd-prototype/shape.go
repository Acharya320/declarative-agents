// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package specmd

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// A field's value reaches Markdown in one of four shapes, chosen from the
// value itself: a string is prose, a list of strings a bullet list, a list of
// flat string maps sharing one key order a table, and anything else -- or a
// value whose visible form Markdown would misread -- a fenced yaml block.
// The fenced block keeps every value lossless; the other shapes keep it
// readable.
const (
	shapeProse = "prose"
	shapeList  = "list"
	shapeTable = "table"
	shapeYAML  = "yaml"
)

// renderNode returns the Markdown lines for a field value.
func renderNode(node *yaml.Node) ([]string, error) {
	switch {
	case isString(node) && proseSafe(node.Value):
		return proseLines(node.Value), nil
	case isStringList(node) && listSafe(node):
		return listLines(node), nil
	case tableSafe(node):
		return tableLines(node), nil
	default:
		return fencedYAML(node)
	}
}

func isString(node *yaml.Node) bool {
	return node.Kind == yaml.ScalarNode && node.ShortTag() == "!!str"
}

func isStringList(node *yaml.Node) bool {
	if node.Kind != yaml.SequenceNode || len(node.Content) == 0 {
		return false
	}
	for _, item := range node.Content {
		if !isString(item) {
			return false
		}
	}
	return true
}

// proseLines prints text as-is, minus trailing newlines; the parser reads the
// same lines back verbatim.
func proseLines(text string) []string {
	return strings.Split(strings.TrimRight(text, "\n"), "\n")
}

// proseSafe reports whether text survives as raw Markdown: it must not open a
// structural block (heading, marker, fence, rule) nor read back as a list or
// table, and it must not start or end with whitespace Markdown would drop.
func proseSafe(text string) bool {
	trimmed := strings.TrimRight(text, "\n")
	if trimmed == "" || strings.TrimSpace(trimmed) != trimmed {
		return false
	}
	for _, line := range strings.Split(trimmed, "\n") {
		if structuralLine(line) {
			return false
		}
	}
	return classifyContent(trimmed) == shapeProse
}

// structuralLine reports lines that would open a block the parser treats as
// structure rather than content.
func structuralLine(line string) bool {
	text := strings.TrimLeft(line, " ")
	if len(line)-len(text) >= 4 {
		return true // indented code
	}
	for _, prefix := range []string{"#", "<", "```", "~~~", ">"} {
		if strings.HasPrefix(text, prefix) {
			return true
		}
	}
	stripped := strings.ReplaceAll(text, " ", "")
	return len(stripped) >= 3 && strings.Trim(stripped, "-") == "" ||
		len(stripped) >= 3 && strings.Trim(stripped, "=") == "" ||
		len(stripped) >= 3 && strings.Trim(stripped, "*") == "" ||
		len(stripped) >= 3 && strings.Trim(stripped, "_") == ""
}

func listSafe(node *yaml.Node) bool {
	for _, item := range node.Content {
		if !proseSafe(item.Value) {
			return false
		}
	}
	return true
}

// listLines prints "- " before an item's first line and indents its
// continuation lines by two spaces; blank lines stay blank.
func listLines(node *yaml.Node) []string {
	var lines []string
	for _, item := range node.Content {
		for index, line := range proseLines(item.Value) {
			switch {
			case index == 0:
				lines = append(lines, "- "+line)
			case line == "":
				lines = append(lines, "")
			default:
				lines = append(lines, "  "+line)
			}
		}
	}
	return lines
}

// tableSafe reports a non-empty list of maps whose keys appear in one order
// and whose values are single-line strings a table cell can hold exactly.
func tableSafe(node *yaml.Node) bool {
	if node.Kind != yaml.SequenceNode || len(node.Content) == 0 {
		return false
	}
	keys := mapKeys(node.Content[0])
	if len(keys) == 0 {
		return false
	}
	for _, row := range node.Content {
		if row.Kind != yaml.MappingNode || !sameStrings(mapKeys(row), keys) {
			return false
		}
		for i := 1; i < len(row.Content); i += 2 {
			if !cellSafe(row.Content[i]) {
				return false
			}
		}
	}
	for _, key := range keys {
		if !cellSafe(&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key}) {
			return false
		}
	}
	return true
}

func cellSafe(node *yaml.Node) bool {
	return isString(node) && !strings.ContainsAny(node.Value, "\n\\") &&
		strings.TrimSpace(node.Value) == node.Value
}

func mapKeys(node *yaml.Node) []string {
	if node.Kind != yaml.MappingNode {
		return nil
	}
	keys := make([]string, 0, len(node.Content)/2)
	for i := 0; i+1 < len(node.Content); i += 2 {
		keys = append(keys, node.Content[i].Value)
	}
	return keys
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func tableLines(node *yaml.Node) []string {
	keys := mapKeys(node.Content[0])
	separators := make([]string, len(keys))
	for i := range separators {
		separators[i] = "---"
	}
	lines := []string{tableRow(keys), tableRow(separators)}
	for _, row := range node.Content {
		cells := make([]string, 0, len(keys))
		for i := 1; i < len(row.Content); i += 2 {
			cells = append(cells, row.Content[i].Value)
		}
		lines = append(lines, tableRow(cells))
	}
	return lines
}

func tableRow(cells []string) string {
	escaped := make([]string, len(cells))
	for i, cell := range cells {
		escaped[i] = strings.ReplaceAll(cell, "|", `\|`)
	}
	return "| " + strings.Join(escaped, " | ") + " |"
}

// fencedYAML prints the value as a yaml code block, choosing a fence longer
// than any backtick run inside it.
func fencedYAML(node *yaml.Node) ([]string, error) {
	data, err := yaml.Marshal(explicitNulls(node))
	if err != nil {
		return nil, fmt.Errorf("marshal fenced value: %w", err)
	}
	body := strings.TrimRight(string(data), "\n")
	fence := "```"
	for strings.Contains(body, fence) {
		fence += "`"
	}
	lines := []string{fence + "yaml"}
	lines = append(lines, strings.Split(body, "\n")...)
	return append(lines, fence), nil
}

// parseContent reads a field's content back into a value, choosing the shape
// from the Markdown itself.
func parseContent(text string) (*yaml.Node, error) {
	switch classifyContent(text) {
	case shapeYAML:
		return parseFencedYAML(text)
	case shapeTable:
		return parseTable(text)
	case shapeList:
		return parseList(text), nil
	default:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: text}, nil
	}
}

func parseFencedYAML(text string) (*yaml.Node, error) {
	lines := strings.Split(text, "\n")
	// The printer trimmed the encoder's final newline; restore it, or a
	// literal block at the end of the value loses its own trailing newline.
	body := strings.Join(lines[1:len(lines)-1], "\n") + "\n"
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(body), &doc); err != nil {
		return nil, fmt.Errorf("fenced yaml: %w", err)
	}
	if len(doc.Content) == 0 {
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null"}, nil
	}
	return doc.Content[0], nil
}

func parseList(text string) *yaml.Node {
	list := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	var current []string
	flush := func() {
		if current != nil {
			value := strings.TrimRight(strings.Join(current, "\n"), "\n")
			list.Content = append(list.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value})
		}
	}
	for _, line := range strings.Split(text, "\n") {
		switch {
		case strings.HasPrefix(line, "- "):
			flush()
			current = []string{strings.TrimPrefix(line, "- ")}
		case line == "":
			current = append(current, "")
		default:
			current = append(current, strings.TrimPrefix(line, "  "))
		}
	}
	flush()
	return list
}

func parseTable(text string) (*yaml.Node, error) {
	lines := strings.Split(text, "\n")
	keys := splitRow(lines[0])
	table := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for _, line := range lines[2:] {
		cells := splitRow(line)
		if len(cells) != len(keys) {
			return nil, fmt.Errorf("table row %q: %d cells, header has %d", line, len(cells), len(keys))
		}
		row := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		for i, key := range keys {
			row.Content = append(row.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key},
				&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: cells[i]})
		}
		table.Content = append(table.Content, row)
	}
	return table, nil
}

// splitRow splits a pipe-table row on unescaped pipes.
func splitRow(line string) []string {
	inner := strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(line), "|"), "|")
	var cells []string
	var cell strings.Builder
	for i := 0; i < len(inner); i++ {
		switch {
		case inner[i] == '\\' && i+1 < len(inner) && inner[i+1] == '|':
			cell.WriteByte('|')
			i++
		case inner[i] == '|':
			cells = append(cells, strings.TrimSpace(cell.String()))
			cell.Reset()
		default:
			cell.WriteByte(inner[i])
		}
	}
	return append(cells, strings.TrimSpace(cell.String()))
}

// explicitNulls copies a node tree, spelling every implicit null (a "!!null"
// scalar with an empty value, as a key without a value parses) as "null":
// the encoder prints an empty null as "" otherwise, and it reads back as an
// empty string.
func explicitNulls(node *yaml.Node) *yaml.Node {
	clone := *node
	if clone.Kind == yaml.ScalarNode && clone.ShortTag() == "!!null" && clone.Value == "" {
		clone.Value = "null"
	}
	if len(node.Content) > 0 {
		clone.Content = make([]*yaml.Node, len(node.Content))
		for i, child := range node.Content {
			clone.Content[i] = explicitNulls(child)
		}
	}
	return &clone
}
