// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package specmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec"
	"gopkg.in/yaml.v3"
)

const yamlLicenseHeader = "# Copyright (c) 2026 Nokia\n# SPDX-License-Identifier: BSD-3-Clause\n"

// mapping builds a YAML mapping node one entry at a time, skipping unset
// values so the printed document carries only what the model holds.
type mapping struct{ node *yaml.Node }

func newMapping() mapping {
	return mapping{node: &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}}
}

func (m mapping) set(key string, value *yaml.Node) {
	if value == nil {
		return
	}
	m.node.Content = append(m.node.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key}, value)
}

func (m mapping) text(key, value string) {
	if value != "" {
		m.set(key, textNode(value))
	}
}

func (m mapping) encode(key string, value any, present bool) error {
	if !present {
		return nil
	}
	var node yaml.Node
	if err := node.Encode(value); err != nil {
		return fmt.Errorf("%s: %w", key, err)
	}
	m.set(key, &node)
	return nil
}

func (m mapping) extensions(extensions []spec.Extension) {
	for i := range extensions {
		value := extensions[i].Value
		m.set(extensions[i].Key, &value)
	}
}

// textNode prints multi-line text as a literal block.
func textNode(value string) *yaml.Node {
	node := stringNode(value)
	if strings.Contains(value, "\n") {
		node.Style = yaml.LiteralStyle
	}
	return node
}

func textList(values []string) *yaml.Node {
	if len(values) == 0 {
		return nil
	}
	list := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for _, value := range values {
		list.Content = append(list.Content, textNode(value))
	}
	return list
}

// PrintYAML renders an SRD as YAML in design.yaml's declared field order:
// the required fields, depends_on, the optional fields, then every untyped
// key in source order.
func PrintYAML(srd *spec.SRD) ([]byte, error) {
	doc := newMapping()
	doc.text("id", srd.ID)
	doc.text("title", srd.Title)
	doc.text("problem", srd.Problem)
	doc.set("goals", textList(srd.Goals))
	requirements, err := requirementsNode(srd)
	if err != nil {
		return nil, err
	}
	doc.set("requirements", requirements)
	doc.set("non_goals", textList(srd.NonGoals))
	criteria, err := criteriaNode(srd.AcceptanceCriteria)
	if err != nil {
		return nil, err
	}
	doc.set("acceptance_criteria", criteria)
	if err := optionalDocumentFields(doc, srd); err != nil {
		return nil, err
	}
	doc.extensions(srd.Extensions)
	var out bytes.Buffer
	encoder := yaml.NewEncoder(&out)
	encoder.SetIndent(2)
	if err := encoder.Encode(doc.node); err != nil {
		return nil, err
	}
	return append([]byte(yamlLicenseHeader), out.Bytes()...), nil
}

func optionalDocumentFields(doc mapping, srd *spec.SRD) error {
	if err := doc.encode("depends_on", srd.DependsOn, len(srd.DependsOn) > 0); err != nil {
		return err
	}
	doc.text("status", srd.Status)
	doc.text("retired_in", srd.RetiredIn)
	if err := doc.encode("replaced_by", srd.ReplacedBy, len(srd.ReplacedBy) > 0); err != nil {
		return err
	}
	if err := doc.encode("needs_clarification", srd.NeedsClarification, len(srd.NeedsClarification) > 0); err != nil {
		return err
	}
	doc.set("assumptions", textList(srd.Assumptions))
	return doc.encode("derived_from", srd.DerivedFrom, len(srd.DerivedFrom) > 0)
}

func requirementsNode(srd *spec.SRD) (*yaml.Node, error) {
	if len(srd.OrderedGroups) == 0 {
		return nil, nil
	}
	groups := newMapping()
	for _, groupID := range srd.OrderedGroups {
		group := srd.Requirements[groupID]
		node := newMapping()
		node.text("title", group.Title)
		node.extensions(group.Extensions)
		items := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for i := range group.Items {
			item, err := itemNode(&group.Items[i])
			if err != nil {
				return nil, err
			}
			items.Content = append(items.Content, item)
		}
		node.set("items", items)
		groups.set(groupID, node.node)
	}
	return groups.node, nil
}

// itemNode prints the tagged form (R1.1: text) when the item carries nothing
// beyond its text and metadata, and the object form otherwise.
func itemNode(item *spec.RequirementItem) (*yaml.Node, error) {
	node := newMapping()
	if item.Evidence == "" && len(item.Extensions) == 0 {
		node.set(item.ID, textNode(item.Text))
	} else {
		node.text("id", item.ID)
		node.set("text", textNode(item.Text))
		node.text("evidence", item.Evidence)
	}
	if item.Weight != 1 {
		if err := node.encode("weight", item.Weight, true); err != nil {
			return nil, err
		}
	}
	if err := retirementFields(node, item.Status, item.RetiredIn, item.ReplacedBy, item.DerivedFrom); err != nil {
		return nil, fmt.Errorf("item %s: %w", item.ID, err)
	}
	node.extensions(item.Extensions)
	return node.node, nil
}

func retirementFields(node mapping, status, retiredIn string, replacedBy []string, derived []spec.ExternalReference) error {
	node.text("status", status)
	node.text("retired_in", retiredIn)
	if err := node.encode("replaced_by", replacedBy, len(replacedBy) > 0); err != nil {
		return err
	}
	return node.encode("derived_from", derived, len(derived) > 0)
}

func criteriaNode(criteria []spec.AcceptanceCriterion) (*yaml.Node, error) {
	if len(criteria) == 0 {
		return nil, nil
	}
	list := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for i := range criteria {
		criterion := &criteria[i]
		node := newMapping()
		node.text("id", criterion.ID)
		node.set("criterion", textNode(criterion.Criterion))
		if err := node.encode("traces", criterion.Traces, len(criterion.Traces) > 0); err != nil {
			return nil, err
		}
		if err := retirementFields(node, criterion.Status, criterion.RetiredIn,
			criterion.ReplacedBy, criterion.DerivedFrom); err != nil {
			return nil, fmt.Errorf("criterion %s: %w", criterion.ID, err)
		}
		node.extensions(criterion.Extensions)
		list.Content = append(list.Content, node.node)
	}
	return list, nil
}
