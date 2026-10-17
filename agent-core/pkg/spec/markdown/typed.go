// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Typed forms (eng10, typed forms) carry values that are not strings, or
// that hold structure, as readable Markdown. The field marker names the form
// (form: "scalar"), so a field without it reads as before and no document
// printed earlier changes meaning. Inside a typed form every value is a YAML
// flow value on one line (GH-2641).
const (
	attrForm        = "form"
	formScalar      = "scalar"
	formTypedTable  = "typed table"
	formTree        = "tree"
	flowPlainReject = "[]{},#&*!|>'\"%@`"
)

// typedFormFor prints a value no simple shape holds in the first typed form
// that parses back to the same node, or reports that none does.
func typedFormFor(value *yaml.Node) (string, []string, bool) {
	candidates := []struct {
		form  string
		print func(*yaml.Node) ([]string, error)
	}{{formScalar, scalarLines}, {formTypedTable, typedTableLines}, {formTree, treeLines}}
	for _, c := range candidates {
		lines, err := c.print(value)
		if err != nil {
			continue
		}
		back, err := parseTyped(c.form, strings.Join(lines, "\n"))
		if err == nil && Diff(value, back) == "" {
			return c.form, lines, true
		}
	}
	return "", nil, false
}

// parseTyped reads a field's content in the form its marker names.
func parseTyped(form, text string) (*yaml.Node, error) {
	switch form {
	case formScalar:
		return decodeFlow(strings.TrimSpace(text))
	case formTypedTable:
		return parseTypedTable(text)
	case formTree:
		return parseTree(text)
	}
	return nil, fmt.Errorf("unknown form %q (scalar, typed table, or tree)", form)
}

// encodeFlow writes a scalar, or a list of scalars, as a one-line YAML flow
// value that reads back as the same node: a string bare when YAML reads it
// back unchanged, JSON-quoted otherwise.
func encodeFlow(node *yaml.Node) (string, error) {
	switch node.Kind {
	case yaml.ScalarNode:
		return encodeFlowScalar(node)
	case yaml.SequenceNode:
		items := make([]string, 0, len(node.Content))
		for _, item := range node.Content {
			if item.Kind != yaml.ScalarNode {
				return "", errors.New("a flow list holds scalars")
			}
			text, err := encodeFlowScalar(item)
			if err != nil {
				return "", err
			}
			if strings.ContainsAny(text, ",[]") && !strings.HasPrefix(text, `"`) {
				text = jsonString(item.Value)
			}
			items = append(items, text)
		}
		return "[" + strings.Join(items, ", ") + "]", nil
	}
	return "", fmt.Errorf("a flow value is a scalar or a list of scalars, found %s", node.ShortTag())
}

func encodeFlowScalar(node *yaml.Node) (string, error) {
	if node.ShortTag() != "!!str" {
		if back, err := decodeFlow(node.Value); err == nil && Diff(node, back) == "" {
			return node.Value, nil
		}
		return "", fmt.Errorf("scalar %q does not read back as %s", node.Value, node.ShortTag())
	}
	plain := node.Value
	if strings.Contains(plain, "\n") {
		return "", errors.New("multi-line text is not a flow value")
	}
	if plain != "" && !strings.ContainsAny(plain[:1], flowPlainReject) && !strings.Contains(plain, "\n") {
		if back, err := decodeFlow(plain); err == nil && back.ShortTag() == "!!str" && back.Value == plain {
			return plain, nil
		}
	}
	return jsonString(plain), nil
}

func jsonString(value string) string {
	data, _ := json.Marshal(value)
	return string(data)
}

// decodeFlow reads one YAML flow value.
func decodeFlow(text string) (*yaml.Node, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(text), &doc); err != nil {
		return nil, fmt.Errorf("flow value %q: %w", text, err)
	}
	if len(doc.Content) == 0 {
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"}, nil
	}
	value := doc.Content[0]
	value.Style = 0
	for _, item := range value.Content {
		item.Style = 0
	}
	return value, nil
}

// scalarLines prints a lone scalar as one flow value.
func scalarLines(value *yaml.Node) ([]string, error) {
	if value.Kind != yaml.ScalarNode {
		return nil, errors.New("the scalar form holds a scalar")
	}
	text, err := encodeFlowScalar(value)
	return []string{text}, err
}

// typedTableLines prints a list of records whose values are scalars or lists
// of scalars as a table: one column per key, in the first record's order
// extended by later keys, a flow value per cell, an empty cell for a key the
// record lacks. Each record's keys must follow the column order, so the
// table reads back with the same key order.
func typedTableLines(value *yaml.Node) ([]string, error) {
	if value.Kind != yaml.SequenceNode || len(value.Content) == 0 {
		return nil, errors.New("a typed table holds a non-empty list of records")
	}
	columns, err := tableColumns(value)
	if err != nil {
		return nil, err
	}
	separators := make([]string, len(columns))
	for i := range separators {
		separators[i] = "---"
	}
	lines := []string{tableRow(columns), tableRow(separators)}
	for _, row := range value.Content {
		cells := make([]string, len(columns))
		for i, column := range columns {
			cell := mapValue(row, column)
			if cell == nil {
				continue
			}
			if cells[i], err = encodeFlow(cell); err != nil {
				return nil, fmt.Errorf("%s: %w", column, err)
			}
		}
		lines = append(lines, tableRow(cells))
	}
	return lines, nil
}

// tableColumns returns a column order every record's keys follow: the
// records' key orders merged, ties broken by first appearance. It fails when
// a record is not a map, a key is not a bare name, or two records order the
// same keys differently.
func tableColumns(value *yaml.Node) ([]string, error) {
	var first []string
	after := map[string]map[string]bool{} // key -> keys that must follow it
	for _, row := range value.Content {
		if row.Kind != yaml.MappingNode || len(row.Content) == 0 {
			return nil, errors.New("a typed table row is a non-empty map")
		}
		keys := mapKeys(row)
		for i, key := range keys {
			if !bareName(key) {
				return nil, fmt.Errorf("key %q is not a bare name", key)
			}
			if _, seen := after[key]; !seen {
				after[key] = map[string]bool{}
				first = append(first, key)
			}
			if i > 0 {
				after[keys[i-1]][key] = true
			}
		}
	}
	return mergeOrder(first, after)
}

// mergeOrder sorts keys so each comes after every key that must precede it,
// keeping first-appearance order among the free ones.
func mergeOrder(first []string, after map[string]map[string]bool) ([]string, error) {
	before := map[string]int{}
	for _, next := range after {
		for key := range next {
			before[key]++
		}
	}
	placed := map[string]bool{}
	var order []string
	for len(order) < len(first) {
		progressed := false
		for _, key := range first {
			if placed[key] || before[key] > 0 {
				continue
			}
			placed[key] = true
			order = append(order, key)
			for next := range after[key] {
				before[next]--
			}
			progressed = true
			break
		}
		if !progressed {
			return nil, errors.New("records order their keys inconsistently")
		}
	}
	return order, nil
}

// bareName reports a key the typed forms print without quoting: letters,
// digits, underscores, hyphens, and dots, starting with a letter, digit, or
// underscore.
func bareName(key string) bool {
	if key == "" || strings.ContainsAny(key[:1], "-.") {
		return false
	}
	for _, r := range key {
		if !(r == '_' || r == '-' || r == '.' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9') {
			return false
		}
	}
	return true
}

func parseTypedTable(text string) (*yaml.Node, error) {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) < 3 {
		return nil, errors.New("a typed table has a header, a separator, and rows")
	}
	columns := splitRow(lines[0])
	table := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for _, line := range lines[2:] {
		cells := splitRow(line)
		if len(cells) != len(columns) {
			return nil, fmt.Errorf("table row %q: %d cells, header has %d", line, len(cells), len(columns))
		}
		row := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		for i, column := range columns {
			if cells[i] == "" {
				continue
			}
			cell, err := decodeFlow(cells[i])
			if err != nil {
				return nil, err
			}
			row.Content = append(row.Content, stringNode(column), cell)
		}
		table.Content = append(table.Content, row)
	}
	return table, nil
}
