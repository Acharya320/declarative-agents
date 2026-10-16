// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package spec

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// ParseSRD reads a single SRD document and returns the parsed struct.
func ParseSRD(path string) (SRD, error) {
	data, err := readDocument(path)
	if err != nil {
		return SRD{}, fmt.Errorf("read SRD %s: %w", path, err)
	}
	return parseSRDBytes(data, path)
}

// rawSRD is the intermediate representation for YAML unmarshalling.
type rawSRD struct {
	ID                 string                `yaml:"id"`
	Title              string                `yaml:"title"`
	Problem            string                `yaml:"problem"`
	Goals              yaml.Node             `yaml:"goals"`
	Requirements       yaml.Node             `yaml:"requirements"`
	NonGoals           yaml.Node             `yaml:"non_goals"`
	AcceptanceCriteria []AcceptanceCriterion `yaml:"acceptance_criteria"`
	DependsOn          []Dependency          `yaml:"depends_on"`
	Status             string                `yaml:"status"`
	RetiredIn          string                `yaml:"retired_in"`
	ReplacedBy         []string              `yaml:"replaced_by"`
	NeedsClarification []ClarificationEntry  `yaml:"needs_clarification"`
	Assumptions        yaml.Node             `yaml:"assumptions"`
	DerivedFrom        []ExternalReference   `yaml:"derived_from"`
}

// srdTypedKeys are the top-level keys rawSRD types; any other top-level key
// is kept as an extension.
var srdTypedKeys = map[string]bool{
	"id": true, "title": true, "problem": true, "goals": true,
	"requirements": true, "non_goals": true, "acceptance_criteria": true,
	"depends_on": true, "status": true, "retired_in": true, "replaced_by": true,
	"needs_clarification": true, "assumptions": true, "derived_from": true,
}

// srdExtensions returns the document's untyped top-level keys in source
// order.
func srdExtensions(data []byte) ([]Extension, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 {
		return nil, nil
	}
	return untypedKeys(doc.Content[0], srdTypedKeys), nil
}

func parseSRDBytes(data []byte, source string) (SRD, error) {
	var raw rawSRD
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return SRD{}, fmt.Errorf("parse SRD %s: %w", source, err)
	}
	groups, orderedKeys, err := parseRequirements(&raw.Requirements)
	if err != nil {
		return SRD{}, fmt.Errorf("parse SRD %s requirements: %w", source, err)
	}
	goals, nonGoals, assumptions, err := parseSRDLists(&raw, source)
	if err != nil {
		return SRD{}, err
	}
	extensions, err := srdExtensions(data)
	if err != nil {
		return SRD{}, fmt.Errorf("parse SRD %s: %w", source, err)
	}
	return SRD{
		ID:                 raw.ID,
		Title:              raw.Title,
		Problem:            raw.Problem,
		Goals:              goals,
		Requirements:       groups,
		NonGoals:           nonGoals,
		AcceptanceCriteria: raw.AcceptanceCriteria,
		DependsOn:          raw.DependsOn,
		Status:             raw.Status,
		RetiredIn:          raw.RetiredIn,
		ReplacedBy:         raw.ReplacedBy,
		NeedsClarification: raw.NeedsClarification,
		Assumptions:        assumptions,
		DerivedFrom:        raw.DerivedFrom,
		Extensions:         extensions,
		OrderedGroups:      orderedKeys,
	}, nil
}

// parseSRDLists parses the SRD's tagged lists: goals, non-goals, and
// assumptions.
func parseSRDLists(raw *rawSRD, source string) (goals, nonGoals, assumptions []string, err error) {
	goals, err = parseTaggedList(&raw.Goals, "goals")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parse SRD %s: %w", source, err)
	}
	nonGoals, err = parseTaggedList(&raw.NonGoals, "non_goals")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parse SRD %s: %w", source, err)
	}
	assumptions, err = parseTaggedList(&raw.Assumptions, "assumptions")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parse SRD %s: %w", source, err)
	}
	return goals, nonGoals, assumptions, nil
}

// parseTaggedList handles YAML lists of the form:
//
//   - G1: Some text.
//   - G2: More text.
func parseTaggedList(node *yaml.Node, field string) ([]string, error) {
	if node == nil || node.Kind == 0 {
		return nil, nil
	}
	n := node
	if n.Kind == yaml.DocumentNode && len(n.Content) > 0 {
		n = n.Content[0]
	}
	if n.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("%s: expected sequence", field)
	}

	var result []string
	for index, item := range n.Content {
		switch item.Kind {
		case yaml.ScalarNode:
			if item.Tag == "!!null" || item.Value == "" {
				return nil, fmt.Errorf("%s[%d]: expected nonempty scalar or one-entry mapping", field, index)
			}
			result = append(result, item.Value)
		case yaml.MappingNode:
			value, err := taggedListMapping(item)
			if err != nil {
				return nil, fmt.Errorf("%s[%d]: %w", field, index, err)
			}
			result = append(result, value)
		default:
			return nil, fmt.Errorf("%s[%d]: unsupported YAML kind %d", field, index, item.Kind)
		}
	}
	return result, nil
}

func taggedListMapping(item *yaml.Node) (string, error) {
	if len(item.Content) == 2 && item.Content[0].Kind == yaml.ScalarNode && item.Content[1].Kind == yaml.ScalarNode {
		return item.Content[0].Value + ": " + item.Content[1].Value, nil
	}
	fields := map[string]string{}
	for i := 0; i+1 < len(item.Content); i += 2 {
		key, value := item.Content[i], item.Content[i+1]
		if key.Kind != yaml.ScalarNode || value.Kind != yaml.ScalarNode || value.Tag == "!!null" {
			return "", fmt.Errorf("mapping keys and values must be scalars")
		}
		fields[key.Value] = value.Value
	}
	if len(fields) == 2 && fields["id"] != "" && fields["text"] != "" {
		return fields["id"] + ": " + fields["text"], nil
	}
	return "", fmt.Errorf("expected one scalar mapping entry or an id/text object")
}

// parseTouchpointList parses the use-case touchpoints field, which accepts three
// forms: a scalar string ("srdNNN R1 -- desc"), a tagged mapping
// ({T1: "srdNNN R1 -- desc"}), or an object ({id, target, reason}). The object
// form is folded into the canonical "<target> -- <reason>" string that
// parseTouchpoint already parses into an SRD touches edge and its citation
// groups; without this, the flattened "target: srdNNN AC1" string defeats
// parseTouchpoint and every SRD is falsely orphaned (GH-448). Scalar and tagged
// forms pass through unchanged.
func parseTouchpointList(node *yaml.Node, field string) ([]string, error) {
	if node == nil || node.Kind == 0 {
		return nil, nil
	}
	n := node
	if n.Kind == yaml.DocumentNode && len(n.Content) > 0 {
		n = n.Content[0]
	}
	if n.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("%s: expected sequence", field)
	}

	var result []string
	for index, item := range n.Content {
		switch item.Kind {
		case yaml.ScalarNode:
			if item.Tag == "!!null" || item.Value == "" {
				return nil, fmt.Errorf("%s[%d]: expected nonempty scalar or mapping", field, index)
			}
			result = append(result, item.Value)
		case yaml.MappingNode:
			value, err := touchpointFromMapping(item)
			if err != nil {
				return nil, fmt.Errorf("%s[%d]: %w", field, index, err)
			}
			result = append(result, value)
		default:
			return nil, fmt.Errorf("%s[%d]: unsupported YAML kind %d", field, index, item.Kind)
		}
	}
	return result, nil
}

// touchpointFromMapping folds an object-format touchpoint ({id, target, reason})
// into a single canonical "<target> -- <reason>" string. The id is a display
// label the graph does not need, so it is dropped rather than prefixed (a
// non-Tn label would otherwise defeat parseTouchpoint). A mapping without a
// target key is the tagged form ({T1: "srdNNN R1 -- desc"}) and keeps its
// "key: value" flattening.
func touchpointFromMapping(item *yaml.Node) (string, error) {
	fields, err := scalarFields(item)
	if err != nil {
		return "", err
	}
	target := fields["target"]
	if target == "" {
		if len(fields) != 1 {
			return "", fmt.Errorf("tagged mapping must contain exactly one entry")
		}
		for key, value := range fields {
			if key == "" || value == "" {
				return "", fmt.Errorf("tagged mapping key and value must be nonempty")
			}
			return key + ": " + value, nil
		}
	}
	for key := range fields {
		if key != "id" && key != "target" && key != "reason" {
			return "", fmt.Errorf("object mapping contains unknown key %q", key)
		}
	}
	folded := target
	if reason := fields["reason"]; reason != "" {
		folded += " -- " + reason
	}
	if id := fields["id"]; id != "" {
		// Keep the label as the tagged form ("T1: target -- reason") does, so
		// both forms fold to one string (eng10, Comments and two-form fields).
		folded = id + ": " + folded
	}
	return folded, nil
}

// scalarFields reads a mapping whose keys and values are non-null scalars.
func scalarFields(item *yaml.Node) (map[string]string, error) {
	fields := make(map[string]string, len(item.Content)/2)
	for i := 0; i+1 < len(item.Content); i += 2 {
		keyNode, valueNode := item.Content[i], item.Content[i+1]
		if keyNode.Kind != yaml.ScalarNode || valueNode.Kind != yaml.ScalarNode || valueNode.Tag == "!!null" {
			return nil, fmt.Errorf("mapping keys and values must be scalars")
		}
		if _, duplicate := fields[keyNode.Value]; duplicate {
			return nil, fmt.Errorf("duplicate key %q", keyNode.Value)
		}
		fields[keyNode.Value] = valueNode.Value
	}
	return fields, nil
}

// ParseRoadmap reads a road-map.md file.
func ParseRoadmap(path string) (Roadmap, error) {
	data, err := readDocument(path)
	if err != nil {
		return Roadmap{}, fmt.Errorf("read roadmap %s: %w", path, err)
	}
	return parseRoadmapBytes(data, path)
}

// parseRoadmapBytes decodes a road map from its YAML tree.
func parseRoadmapBytes(data []byte, path string) (Roadmap, error) {
	var rm Roadmap
	if err := yaml.Unmarshal(data, &rm); err != nil {
		return Roadmap{}, fmt.Errorf("parse roadmap %s: %w", path, err)
	}
	return rm, nil
}

// ParseSpecIndex reads a SPECIFICATIONS.md file.
func ParseSpecIndex(path string) (SpecIndex, error) {
	data, err := readDocument(path)
	if err != nil {
		return SpecIndex{}, fmt.Errorf("read spec index %s: %w", path, err)
	}
	return parseSpecIndexBytes(data, path)
}

// parseSpecIndexBytes decodes a specification index from its YAML tree.
func parseSpecIndexBytes(data []byte, path string) (SpecIndex, error) {
	var si SpecIndex
	if err := yaml.Unmarshal(data, &si); err != nil {
		return SpecIndex{}, fmt.Errorf("parse spec index %s: %w", path, err)
	}
	return si, nil
}

// rawUseCase is the intermediate representation for YAML unmarshalling.
type rawUseCase struct {
	ID                 string               `yaml:"id"`
	Title              string               `yaml:"title"`
	Summary            string               `yaml:"summary"`
	Actor              string               `yaml:"actor"`
	Trigger            string               `yaml:"trigger"`
	Flow               yaml.Node            `yaml:"flow"`
	Touchpoints        yaml.Node            `yaml:"touchpoints"`
	SuccessCriteria    []SuccessCriterion   `yaml:"success_criteria"`
	OutOfScope         yaml.Node            `yaml:"out_of_scope"`
	TestSuite          string               `yaml:"test_suite"`
	Status             string               `yaml:"status"`
	RetiredIn          string               `yaml:"retired_in"`
	ReplacedBy         []string             `yaml:"replaced_by"`
	NeedsClarification []ClarificationEntry `yaml:"needs_clarification"`
	Assumptions        yaml.Node            `yaml:"assumptions"`
	EdgeCases          yaml.Node            `yaml:"edge_cases"`
	Priority           string               `yaml:"priority"`
}

// ParseUseCase reads a use case document.
func ParseUseCase(path string) (UseCase, error) {
	data, err := readDocument(path)
	if err != nil {
		return UseCase{}, fmt.Errorf("read use case %s: %w", path, err)
	}
	return parseUseCaseBytes(data, path)
}

// parseUseCaseBytes decodes a use case from its YAML tree, in any of the
// forms a document from before the migration may hold.
func parseUseCaseBytes(data []byte, path string) (UseCase, error) {
	var raw rawUseCase
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return UseCase{}, fmt.Errorf("parse use case %s: %w", path, err)
	}
	flow, err := parseTaggedList(&raw.Flow, "flow")
	if err != nil {
		return UseCase{}, fmt.Errorf("parse use case %s: %w", path, err)
	}
	touchpoints, err := parseTouchpointList(&raw.Touchpoints, "touchpoints")
	if err != nil {
		return UseCase{}, fmt.Errorf("parse use case %s: %w", path, err)
	}
	outOfScope, err := parseTaggedList(&raw.OutOfScope, "out_of_scope")
	if err != nil {
		return UseCase{}, fmt.Errorf("parse use case %s: %w", path, err)
	}
	assumptions, err := parseTaggedList(&raw.Assumptions, "assumptions")
	if err != nil {
		return UseCase{}, fmt.Errorf("parse use case %s: %w", path, err)
	}
	edgeCases, err := parseTaggedList(&raw.EdgeCases, "edge_cases")
	if err != nil {
		return UseCase{}, fmt.Errorf("parse use case %s: %w", path, err)
	}
	uc := useCaseFromRaw(&raw, flow, touchpoints, outOfScope)
	uc.Assumptions = assumptions
	uc.EdgeCases = edgeCases
	return uc, nil
}

func useCaseFromRaw(raw *rawUseCase, flow, touchpoints, outOfScope []string) UseCase {
	uc := UseCase{
		ID:                 raw.ID,
		Title:              raw.Title,
		Summary:            raw.Summary,
		Actor:              raw.Actor,
		Trigger:            raw.Trigger,
		Flow:               flow,
		Touchpoints:        touchpoints,
		SuccessCriteria:    raw.SuccessCriteria,
		OutOfScope:         outOfScope,
		TestSuite:          raw.TestSuite,
		RetiredIn:          raw.RetiredIn,
		ReplacedBy:         raw.ReplacedBy,
		NeedsClarification: raw.NeedsClarification,
		Priority:           raw.Priority,
	}
	// Only the retirement marker is copied. Lifecycle statuses (planned,
	// done) stay unparsed so checkIndexDocumentAgreement keeps treating the
	// document side as silent, exactly as before retirement support (D6).
	if raw.Status == statusRetired {
		uc.Status = raw.Status
	}
	return uc
}

// ParseTestSuite reads a test suite document.
func ParseTestSuite(path string) (TestSuite, error) {
	data, err := readDocument(path)
	if err != nil {
		return TestSuite{}, fmt.Errorf("read test suite %s: %w", path, err)
	}
	return parseTestSuiteBytes(data, path)
}

// parseTestSuiteBytes decodes a test suite from its YAML tree.
func parseTestSuiteBytes(data []byte, path string) (TestSuite, error) {
	var ts TestSuite
	if err := yaml.Unmarshal(data, &ts); err != nil {
		return TestSuite{}, fmt.Errorf("parse test suite %s: %w", path, err)
	}
	return ts, nil
}
