// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// nodeYAML maps a node kind to the YAML records it comes from (eng10 Table 5):
// the key holding them, whether they are a list or a map keyed by id, and the
// record keys holding the id, the heading title, and the prose content. A
// keyed record may be written as a one-entry map from id to content.
type nodeYAML struct {
	From    string `yaml:"from"`
	Records string `yaml:"records"`
	ID      string `yaml:"id"`
	Title   string `yaml:"title"`
	Content string `yaml:"content"`
	Keyed   bool   `yaml:"keyed"`
	// DeriveID names how a record without an id gets one when it is
	// canonicalized: "slug" makes a kebab-case id from its title.
	DeriveID string `yaml:"derive_id"`
}

// Record forms, field value conversions, and the field tag.
const (
	recordsList  = "list"
	recordsMap   = "map"
	valueKeyed   = "keyed"
	valueOutline = "outline"
	tagField     = "field"
	headingFree  = "free"
	deriveSlug   = "slug"
)

// checkMapping rejects a YAML mapping the converter could not follow.
func (r *Declaration) checkMapping() error {
	for _, name := range sortedKeys(r.Nodes) {
		m := r.Nodes[name].YAML
		switch {
		case m == (nodeYAML{}):
		case m.From == "":
			return fmt.Errorf("nodes.%s.yaml.from: required", name)
		case m.Records != recordsList && m.Records != recordsMap:
			return fmt.Errorf("nodes.%s.yaml.records: %q is neither list nor map", name, m.Records)
		case m.Records == recordsList && m.ID == "":
			return fmt.Errorf("nodes.%s.yaml.id: list records name their id key", name)
		case m.Keyed && m.Content == "":
			return fmt.Errorf("nodes.%s.yaml.content: keyed records name their content key", name)
		case m.DeriveID != "" && (m.DeriveID != deriveSlug || m.Title == ""):
			return fmt.Errorf("nodes.%s.yaml.derive_id: only slug, from a title", name)
		}
	}
	if err := r.checkSlots(); err != nil {
		return err
	}
	for _, name := range sortedKeys(r.Fields) {
		field := r.Fields[name]
		switch {
		case field.Split != "" && (field.Shape != shapeTable || len(field.Columns) != 3):
			return fmt.Errorf("fields.%s.split: splits a keyed string into a three-column table", name)
		case field.Value == "" || field.Value == valueKeyed:
		case field.Value != valueOutline:
			return fmt.Errorf("fields.%s.value: %q is neither keyed nor outline", name, field.Value)
		case field.ID == "" || field.ContentKey == "":
			return fmt.Errorf("fields.%s: an outline names its id and content keys", name)
		}
	}
	return nil
}

// checkSlots rejects a fields list on a node kind that is not a field slot,
// a field listed twice, and a field rule naming a slot whose list lacks it.
func (r *Declaration) checkSlots() error {
	for _, name := range sortedKeys(r.Nodes) {
		node := r.Nodes[name]
		if len(node.Fields) > 0 && node.Tag != tagField {
			return fmt.Errorf("nodes.%s.fields: only a field slot (tag field) lists fields", name)
		}
		for i, field := range node.Fields {
			if indexOf(node.Fields[:i], field) >= 0 {
				return fmt.Errorf("nodes.%s.fields: %s is listed twice", name, field)
			}
		}
	}
	for _, name := range sortedKeys(r.Fields) {
		for _, slot := range r.Fields[name].In {
			if listed := r.Nodes[slot].Fields; len(listed) > 0 && indexOf(listed, name) < 0 {
				return fmt.Errorf("fields.%s.in: names %s, whose fields list lacks it", name, slot)
			}
		}
	}
	return nil
}

// fieldKind finds the field node kind that opens at level under parent.
func (r *Declaration) fieldKind(level int, parent string) (string, bool) {
	return r.nodeFor(tagField, level, parent)
}

// directRecords finds the node kind whose records sit under key with no
// field heading of their own, as requirement groups sit under requirements.
func (r *Declaration) directRecords(level int, parent, key string) (string, bool) {
	for _, name := range sortedKeys(r.Nodes) {
		node := r.Nodes[name]
		if node.Tag != tagField && node.Level == level && node.Parent == parent && node.YAML.From == key {
			return name, true
		}
	}
	return "", false
}

// fieldRecords finds the node kind whose records a field at level holds,
// one level down, as acceptance criteria sit under the acceptance_criteria
// field. The level tells apart one key holding records in one place and
// something else in another (GH-2663).
func (r *Declaration) fieldRecords(level int, key string) (string, bool) {
	for _, name := range sortedKeys(r.Nodes) {
		if node := r.Nodes[name]; node.Parent == "field:"+key && node.Tag != tagField && node.Level == level+1 {
			return name, true
		}
	}
	return "", false
}

// recordsFit reports a value that holds a node kind's records: a list for
// list records, a map for map records.
func (r *Declaration) recordsFit(kind string, value *yaml.Node) bool {
	if r.Nodes[kind].YAML.Records == recordsMap {
		return value.Kind == yaml.MappingNode
	}
	return value.Kind == yaml.SequenceNode
}

// nestedSlot reports whether a field at level declares a slot for nested
// fields one level down (eng10, declared fields and order; GH-2653).
func (r *Declaration) nestedSlot(level int, key string) bool {
	_, ok := r.fieldKind(level+1, "field:"+key)
	return ok
}

// headingAcronyms are the words a generated heading prints in their
// conventional case (GH-2643).
var headingAcronyms = map[string]string{
	"srd": "SRD", "id": "ID", "ids": "IDs", "api": "API", "ui": "UI", "ux": "UX", "llm": "LLM",
	"cli": "CLI", "rest": "REST", "yaml": "YAML", "url": "URL", "json": "JSON", "http": "HTTP",
	"otel": "OTel",
}

// fieldHeading is the heading a field prints with: the declared one, or the
// name with underscores as spaces, its first letter capitalized, and known
// acronyms in their conventional case.
func (r *Declaration) fieldHeading(key string) string {
	if heading := r.Fields[key].Heading; heading != "" {
		return heading
	}
	if key == "" {
		return ""
	}
	words := strings.Split(key, "_")
	for i, word := range words {
		if acronym, ok := headingAcronyms[word]; ok {
			words[i] = acronym
		}
	}
	heading := strings.Join(words, " ")
	return strings.ToUpper(heading[:1]) + heading[1:]
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
