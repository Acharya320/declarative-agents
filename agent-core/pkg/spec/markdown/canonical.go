// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Canonicalize returns a copy of a YAML document in the canonical form the
// converter reads back (eng10, Comments and two-form fields), with one line
// per change for the conversion report. Frontmatter keys come first; a
// record holding only its id and content takes the keyed form where its
// mapping allows; every other record orders its keys id, title, content,
// attributes, then the rest; map records written as a list become a map; an
// outline item written as "<id>: <text>" becomes keyed. No value changes.
func (d *Declaration) Canonicalize(doc *yaml.Node) (*yaml.Node, []string) {
	c := &canonicalizer{d: d}
	copied := explicitNullsKeep(doc)
	root, _ := documentRoot(copied)
	if root == nil || root.Kind != yaml.MappingNode {
		return copied, nil
	}
	c.frontmatterFirst(root)
	c.members(2, "document", root, "", func(key string) bool { return contains(d.Frontmatter.Allowed, key) })
	return copied, c.changes
}

type canonicalizer struct {
	d       *Declaration
	changes []string
}

func (c *canonicalizer) note(path, format string, args ...any) {
	c.changes = append(c.changes, path+": "+fmt.Sprintf(format, args...))
}

func (c *canonicalizer) frontmatterFirst(root *yaml.Node) {
	var head, rest []*yaml.Node
	moved := false
	for i := 0; i+1 < len(root.Content); i += 2 {
		if contains(c.d.Frontmatter.Allowed, root.Content[i].Value) {
			head = append(head, root.Content[i], root.Content[i+1])
			moved = moved || len(rest) > 0
		} else {
			rest = append(rest, root.Content[i], root.Content[i+1])
		}
	}
	if moved {
		c.note("$", "frontmatter keys moved first")
	}
	root.Content = append(head, rest...)
}

// members canonicalizes every key of a mapping that skip does not claim. A
// key whose records sit straight under it, with no field heading, and that
// holds none is dropped: Markdown has no heading to carry it.
func (c *canonicalizer) members(level int, parent string, mapping *yaml.Node, path string, skip func(string) bool) {
	kept := make([]*yaml.Node, 0, len(mapping.Content))
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		key, value := mapping.Content[i].Value, mapping.Content[i+1]
		where := strings.TrimPrefix(path+"/"+key, "/")
		if !skip(key) {
			if _, direct := c.d.directRecords(level, parent, key); direct && emptyCollection(value) {
				c.note(where, "empty %s dropped; Markdown has no heading for it", key)
				continue
			}
			value = c.member(level, parent, key, value, where)
		}
		kept = append(kept, mapping.Content[i], value)
	}
	mapping.Content = kept
}

func (c *canonicalizer) member(level int, parent, key string, value *yaml.Node, path string) *yaml.Node {
	if kind, ok := c.d.directRecords(level, parent, key); ok {
		return c.records(kind, value, path)
	}
	if kind, ok := c.d.fieldRecords(key); ok {
		return c.records(kind, value, path)
	}
	if field := c.d.Fields[key]; field.Value == valueOutline && value.Kind == yaml.SequenceNode {
		c.outline(field, value, path)
	}
	return value
}

func (c *canonicalizer) records(kind string, value *yaml.Node, path string) *yaml.Node {
	rule := c.d.Nodes[kind]
	if rule.YAML.Records == recordsMap && value.Kind == yaml.SequenceNode {
		value = c.listToMap(rule, value, path)
	}
	switch {
	case rule.YAML.Records == recordsMap && value.Kind == yaml.MappingNode:
		for i := 0; i+1 < len(value.Content); i += 2 {
			id := value.Content[i].Value
			value.Content[i+1] = c.record(kind, rule, id, value.Content[i+1], path+"/"+id)
		}
	case rule.YAML.Records == recordsList && value.Kind == yaml.SequenceNode:
		for i, item := range value.Content {
			if id, _, err := listRecord(rule.YAML, item); err == nil {
				value.Content[i] = c.record(kind, rule, id, item, path+"/"+id)
			}
		}
	}
	return value
}

// listToMap rewrites map records written as a list of records that carry
// their id into the map form, when every entry allows it.
func (c *canonicalizer) listToMap(rule nodeRule, list *yaml.Node, path string) *yaml.Node {
	mapped := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	for _, item := range list.Content {
		id := mapValue(item, rule.YAML.ID)
		if rule.YAML.ID == "" || id == nil || id.Kind != yaml.ScalarNode {
			return list
		}
		rest := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		for i := 0; i+1 < len(item.Content); i += 2 {
			if item.Content[i].Value != rule.YAML.ID {
				rest.Content = append(rest.Content, item.Content[i], item.Content[i+1])
			}
		}
		mapped.Content = append(mapped.Content, stringNode(id.Value), rest)
	}
	c.note(path, "records written as a list became a map keyed by %s", rule.YAML.ID)
	return mapped
}

// record puts one record in canonical form and recurses into its children.
func (c *canonicalizer) record(kind string, rule nodeRule, id string, item *yaml.Node, path string) *yaml.Node {
	if item.Kind != yaml.MappingNode {
		return item
	}
	m := rule.YAML
	keyedNow := m.Keyed && len(item.Content) == 2 && item.Content[0].Value != m.ID
	if keyedNow {
		return item
	}
	if m.Keyed && m.Records == recordsList && onlyKeys(item, m.ID, m.Content) {
		c.note(path, "record with only %s and %s written keyed", m.ID, m.Content)
		return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map", Content: []*yaml.Node{stringNode(id), mapValue(item, m.Content)}}
	}
	if c.reorder(rule, item) {
		c.note(path, "keys ordered %s", strings.Join(mapKeys(item), ", "))
	}
	c.members(rule.Level+1, kind, item, path, func(key string) bool { return consumedKey(rule, key) })
	return item
}

// reorder sorts a record's keys into id, title, content, attributes, rest,
// keeping record order within each group, and reports whether anything
// moved.
func (c *canonicalizer) reorder(rule nodeRule, item *yaml.Node) bool {
	rank := func(key string) int {
		switch {
		case key == rule.YAML.ID && rule.YAML.Records == recordsList:
			return 0
		case key == rule.YAML.Title && rule.Heading != headingID:
			return 1
		case key == rule.YAML.Content:
			return 2
		case contains(rule.Attributes, key):
			return 3
		}
		return 4
	}
	var ordered []*yaml.Node
	for group := 0; group <= 4; group++ {
		for i := 0; i+1 < len(item.Content); i += 2 {
			if rank(item.Content[i].Value) == group {
				ordered = append(ordered, item.Content[i], item.Content[i+1])
			}
		}
	}
	moved := false
	for i := range ordered {
		moved = moved || ordered[i] != item.Content[i]
	}
	item.Content = ordered
	return moved
}

// outline puts an outline field's items in canonical form.
func (c *canonicalizer) outline(field fieldRule, list *yaml.Node, path string) {
	m := nodeYAML{ID: field.ID, Content: field.ContentKey, Keyed: true, Records: recordsList}
	for i, item := range list.Content {
		where := fmt.Sprintf("%s[%d]", path, i)
		if id, text, ok := strings.Cut(item.Value, ": "); ok && isString(item) && !strings.ContainsAny(id, " \n") {
			c.note(where, "item %q written keyed", id)
			list.Content[i] = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map", Content: []*yaml.Node{stringNode(id), stringNode(text)}}
			continue
		}
		if item.Kind == yaml.MappingNode && onlyKeys(item, m.ID, m.Content) {
			c.note(where, "record with only %s and %s written keyed", m.ID, m.Content)
			list.Content[i] = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map",
				Content: []*yaml.Node{stringNode(mapValue(item, m.ID).Value), mapValue(item, m.Content)}}
			continue
		}
		if item.Kind == yaml.MappingNode && len(item.Content) > 2 && c.reorder(nodeRule{YAML: m, Heading: headingID}, item) {
			c.note(where, "keys ordered %s", strings.Join(mapKeys(item), ", "))
		}
	}
}

// onlyKeys reports a mapping holding exactly the two keys, in either order.
func onlyKeys(item *yaml.Node, a, b string) bool {
	return item.Kind == yaml.MappingNode && len(item.Content) == 4 &&
		mapValue(item, a) != nil && mapValue(item, b) != nil && a != b
}

// explicitNullsKeep copies a tree with its comments, so Canonicalize never
// changes the caller's document.
func explicitNullsKeep(node *yaml.Node) *yaml.Node {
	if node == nil {
		return nil
	}
	clone := *node
	if len(node.Content) > 0 {
		clone.Content = make([]*yaml.Node, len(node.Content))
		for i, child := range node.Content {
			clone.Content[i] = explicitNullsKeep(child)
		}
	}
	return &clone
}
