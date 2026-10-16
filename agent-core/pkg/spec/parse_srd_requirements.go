// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package spec

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func parseRequirements(node *yaml.Node) (map[string]RequirementGroup, []string, error) {
	if node == nil || node.Kind == 0 {
		return nil, nil, nil
	}

	n := node
	if n.Kind == yaml.DocumentNode && len(n.Content) > 0 {
		n = n.Content[0]
	}
	switch n.Kind {
	case yaml.MappingNode:
		return parseMappedGroups(n)
	case yaml.SequenceNode:
		return parseListedGroups(n)
	default:
		return nil, nil, fmt.Errorf("requirements: expected mapping or sequence, got kind %d", n.Kind)
	}
}

// parseMappedGroups reads requirements written as a mapping from group id
// to group (R1: {title, items}).
func parseMappedGroups(n *yaml.Node) (map[string]RequirementGroup, []string, error) {
	groups := make(map[string]RequirementGroup)
	var orderedKeys []string
	for i := 0; i+1 < len(n.Content); i += 2 {
		groupID := n.Content[i].Value
		orderedKeys = append(orderedKeys, groupID)
		group, err := parseRequirementGroup(groupID, n.Content[i+1], mappedGroupKeys)
		if err != nil {
			return nil, nil, fmt.Errorf("group %s: %w", groupID, err)
		}
		groups[groupID] = group
	}
	return groups, orderedKeys, nil
}

// parseListedGroups reads requirements written as a list of groups that
// carry their own id ([{id: R1, title, items}]), keeping list order.
func parseListedGroups(n *yaml.Node) (map[string]RequirementGroup, []string, error) {
	groups := make(map[string]RequirementGroup)
	var orderedKeys []string
	for index, groupNode := range n.Content {
		groupID := mappingScalar(groupNode, "id")
		if groupID == "" {
			return nil, nil, fmt.Errorf(
				"requirements: expected mapping of groups or a list of groups that carry an id; element %d has no id", index)
		}
		orderedKeys = append(orderedKeys, groupID)
		group, err := parseRequirementGroup(groupID, groupNode, listedGroupKeys)
		if err != nil {
			return nil, nil, fmt.Errorf("group %s: %w", groupID, err)
		}
		groups[groupID] = group
	}
	return groups, orderedKeys, nil
}

// mappedGroupKeys and listedGroupKeys are the group keys the model types;
// any other group key is kept as an extension.
var (
	mappedGroupKeys = map[string]bool{"title": true, "items": true}
	listedGroupKeys = map[string]bool{"id": true, "title": true, "items": true}
)

// mappingScalar returns the scalar value of key in a mapping node, or "".
func mappingScalar(node *yaml.Node, key string) string {
	if node == nil || node.Kind != yaml.MappingNode {
		return ""
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value == key && node.Content[i+1].Kind == yaml.ScalarNode {
			return node.Content[i+1].Value
		}
	}
	return ""
}

type rawGroup struct {
	Title string    `yaml:"title"`
	Items yaml.Node `yaml:"items"`
}

func parseRequirementGroup(groupID string, node *yaml.Node, typed map[string]bool) (RequirementGroup, error) {
	var rg rawGroup
	if err := node.Decode(&rg); err != nil {
		return RequirementGroup{}, fmt.Errorf("decode group: %w", err)
	}

	items, err := parseRequirementItems(groupID, &rg.Items)
	if err != nil {
		return RequirementGroup{}, err
	}

	return RequirementGroup{
		Title:      rg.Title,
		Items:      items,
		Extensions: untypedKeys(node, typed),
	}, nil
}

func parseRequirementItems(groupID string, node *yaml.Node) ([]RequirementItem, error) {
	if node == nil || node.Kind == 0 {
		return nil, nil
	}
	if node.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("items: expected sequence, got kind %d", node.Kind)
	}

	var items []RequirementItem
	for _, itemNode := range node.Content {
		if itemNode.Kind != yaml.MappingNode {
			return nil, fmt.Errorf("items: expected mapping entry, got kind %d", itemNode.Kind)
		}

		item, err := parseOneItem(itemNode)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// parseOneItem reads one requirement item in either form. An item with an
// id key is object form ({id, text, evidence, ...}); otherwise its single
// non-metadata key is the id and that key's value the text (R1.1: text).
func parseOneItem(node *yaml.Node) (RequirementItem, error) {
	if mappingScalar(node, "id") != "" {
		return parseObjectItem(node)
	}
	return parseTaggedItem(node)
}

func parseTaggedItem(node *yaml.Node) (RequirementItem, error) {
	item := RequirementItem{Weight: 1}
	for i := 0; i+1 < len(node.Content); i += 2 {
		key, val := node.Content[i].Value, node.Content[i+1]
		handled, err := parseItemMetadata(&item, key, val)
		if err != nil {
			return RequirementItem{}, err
		}
		switch {
		case handled:
		case item.ID == "":
			item.ID = key
			item.Text = val.Value
		case key == "evidence":
			item.Evidence = val.Value
		default:
			item.Extensions = append(item.Extensions, Extension{Key: key, Value: *val})
		}
	}
	if item.ID == "" {
		return RequirementItem{}, fmt.Errorf("requirement item has no ID key")
	}
	return item, nil
}

func parseObjectItem(node *yaml.Node) (RequirementItem, error) {
	item := RequirementItem{Weight: 1}
	for i := 0; i+1 < len(node.Content); i += 2 {
		key, val := node.Content[i].Value, node.Content[i+1]
		switch key {
		case "id":
			item.ID = val.Value
		case "text":
			item.Text = val.Value
		case "evidence":
			item.Evidence = val.Value
		default:
			handled, err := parseItemMetadata(&item, key, val)
			if err != nil {
				return RequirementItem{}, fmt.Errorf("item %s: %w", mappingScalar(node, "id"), err)
			}
			if !handled {
				item.Extensions = append(item.Extensions, Extension{Key: key, Value: *val})
			}
		}
	}
	return item, nil
}

// parseItemMetadata decodes a requirement item's metadata keys and reports
// whether key was one; any other key stays with the caller.
func parseItemMetadata(item *RequirementItem, key string, val *yaml.Node) (bool, error) {
	switch key {
	case "weight":
		if err := val.Decode(&item.Weight); err != nil {
			return false, fmt.Errorf("decode weight: %w", err)
		}
	case "status":
		item.Status = val.Value
	case "retired_in":
		item.RetiredIn = val.Value
	case "replaced_by":
		if err := val.Decode(&item.ReplacedBy); err != nil {
			return false, fmt.Errorf("decode replaced_by: %w", err)
		}
	case "derived_from":
		if err := val.Decode(&item.DerivedFrom); err != nil {
			return false, fmt.Errorf("decode derived_from: %w", err)
		}
	default:
		return false, nil
	}
	return true, nil
}
