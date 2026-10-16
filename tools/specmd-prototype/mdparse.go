// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package specmd

import (
	"fmt"
	"strings"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec"
	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec/markdown"
	"gopkg.in/yaml.v3"
)

// Node kinds a heading opens, named by marker tag and heading depth.
const (
	nodeDocField       = "document field"
	nodeGroup          = "requirement group"
	nodeGroupField     = "group field"
	nodeItem           = "requirement item"
	nodeItemField      = "item field"
	nodeCriterion      = "acceptance criterion"
	nodeCriterionField = "criterion field"
)

// openNode is the node whose content the reader is collecting.
type openNode struct {
	kind  string
	key   string
	line  int
	start int
	end   int
}

// reader walks top-level blocks: a marker is pending until the heading below
// it opens a node, and the node collects content blocks until the next
// marker. Parsing stops at the first structural error; GH-2589 collects
// findings instead.
type reader struct {
	srd        *spec.SRD
	lines      []string
	offset     int
	pending    *markdown.Marker
	open       *openNode
	groups     []*spec.RequirementGroup
	criteria   bool
	leafIsItem bool
}

// ParseMarkdown reads an SRD from Markdown into the shared model.
func ParseMarkdown(data []byte, source string) (spec.SRD, error) {
	fm, body, offset, err := splitFrontmatter(data)
	if err != nil {
		return spec.SRD{}, fmt.Errorf("%s: %w", source, err)
	}
	srd := spec.SRD{
		ID: fm.ID, Title: fm.Title, Status: fm.Status, RetiredIn: fm.RetiredIn,
		ReplacedBy: fm.ReplacedBy, DependsOn: fm.DependsOn, DerivedFrom: fm.DerivedFrom,
		Requirements: map[string]spec.RequirementGroup{},
	}
	blocks, lines := markdown.TopLevelBlocks(body)
	r := &reader{srd: &srd, lines: lines, offset: offset}
	for _, b := range blocks {
		if err := r.step(b); err != nil {
			return spec.SRD{}, fmt.Errorf("%s:%d: %w", source, b.Start+offset+1, err)
		}
	}
	if err := r.close(); err != nil {
		return spec.SRD{}, fmt.Errorf("%s: %w", source, err)
	}
	if r.pending != nil {
		return spec.SRD{}, fmt.Errorf("%s: marker %q has no heading below it", source, r.pending.String())
	}
	for i, group := range r.groups {
		srd.Requirements[srd.OrderedGroups[i]] = *group
	}
	return srd, nil
}

func (r *reader) step(b markdown.Block) error {
	if text, ok := markdown.CommentLine(b, r.lines); ok {
		m, isMarker, err := markdown.ParseMarker(text)
		if err != nil || !isMarker {
			return err
		}
		if r.pending != nil {
			return fmt.Errorf("marker %q has no heading below it", r.pending.String())
		}
		if err := r.close(); err != nil {
			return err
		}
		r.pending = &m
		return nil
	}
	if level, heading, ok := markdown.HeadingOf(b, r.lines); ok {
		return r.heading(b, level, heading)
	}
	if r.open == nil {
		return fmt.Errorf("content outside any tagged node")
	}
	if r.open.start < 0 {
		r.open.start = b.Start
	}
	r.open.end = b.End
	return nil
}

func (r *reader) heading(b markdown.Block, level int, heading string) error {
	if r.pending == nil {
		return fmt.Errorf("heading %q has no marker above it", heading)
	}
	m := *r.pending
	r.pending = nil
	kind, err := r.nodeKind(m, level)
	if err != nil {
		return fmt.Errorf("heading %q: %w", heading, err)
	}
	r.open = &openNode{kind: kind, key: m.Value, line: b.Start, start: -1, end: -1}
	return r.begin(m, heading)
}

// nodeKind resolves what a marker and heading depth open in the current
// context.
func (r *reader) nodeKind(m markdown.Marker, level int) (string, error) {
	switch {
	case m.Tag == "field" && level == 2:
		return nodeDocField, nil
	case m.Tag == "req" && level == 2 && !strings.Contains(m.Value, "."):
		return nodeGroup, nil
	case m.Tag == "field" && level == 3 && len(r.groups) > 0 && !r.criteria:
		return nodeGroupField, nil
	case m.Tag == "req" && level == 3 && len(r.groups) > 0 && !r.criteria:
		return nodeItem, nil
	case m.Tag == "ac" && level == 3 && r.criteria:
		return nodeCriterion, nil
	case m.Tag == "field" && level == 4 && r.leafIsItem:
		return nodeItemField, nil
	case m.Tag == "field" && level == 4 && r.criteria && len(r.srd.AcceptanceCriteria) > 0:
		return nodeCriterionField, nil
	}
	return "", fmt.Errorf("marker %q cannot open a level-%d heading here", m.String(), level)
}

// begin records what a heading opens before its content arrives.
func (r *reader) begin(m markdown.Marker, heading string) error {
	switch r.open.kind {
	case nodeDocField:
		r.criteria = m.Value == keyAcceptanceCriteria
		r.leafIsItem = false
	case nodeGroup:
		r.srd.OrderedGroups = append(r.srd.OrderedGroups, m.Value)
		title := strings.TrimPrefix(strings.TrimPrefix(heading, m.Value), " - ")
		r.groups = append(r.groups, &spec.RequirementGroup{Title: title})
		r.criteria, r.leafIsItem = false, false
	case nodeItem:
		item := spec.RequirementItem{ID: m.Value}
		if err := applyItemAttrs(m, &item); err != nil {
			return err
		}
		group := r.groups[len(r.groups)-1]
		group.Items = append(group.Items, item)
		r.leafIsItem = true
	case nodeCriterion:
		criterion := spec.AcceptanceCriterion{ID: m.Value}
		if err := applyCriterionAttrs(m, &criterion); err != nil {
			return err
		}
		r.srd.AcceptanceCriteria = append(r.srd.AcceptanceCriteria, criterion)
		r.leafIsItem = false
	}
	return nil
}

// close parses the open node's content and stores it on the model.
func (r *reader) close() error {
	open := r.open
	r.open = nil
	if open == nil {
		return nil
	}
	if open.start < 0 {
		return r.store(open, nil)
	}
	value, err := markdown.ParseContent(markdown.SpanText(r.lines, open.start, open.end))
	if err != nil {
		return fmt.Errorf("line %d: %s %s: %w", open.start+r.offset+1, open.kind, open.key, err)
	}
	return r.store(open, value)
}

func (r *reader) store(open *openNode, value *yaml.Node) error {
	switch open.kind {
	case nodeDocField:
		return r.storeDocField(open.key, value)
	case nodeGroup:
		if value != nil {
			return fmt.Errorf("group %s: content belongs under a tagged field or item", open.key)
		}
	case nodeGroupField:
		group := r.groups[len(r.groups)-1]
		group.Extensions = append(group.Extensions, newExtension(open.key, value))
	case nodeItem:
		return assignString(value, &r.lastItem().Text)
	case nodeItemField:
		item := r.lastItem()
		if open.key == keyEvidence {
			return assignString(value, &item.Evidence)
		}
		item.Extensions = append(item.Extensions, newExtension(open.key, value))
	case nodeCriterion:
		return assignString(value, &r.lastCriterion().Criterion)
	case nodeCriterionField:
		criterion := r.lastCriterion()
		if open.key == keyTraces {
			return assignStrings(value, &criterion.Traces)
		}
		criterion.Extensions = append(criterion.Extensions, newExtension(open.key, value))
	}
	return nil
}

func (r *reader) lastItem() *spec.RequirementItem {
	group := r.groups[len(r.groups)-1]
	return &group.Items[len(group.Items)-1]
}

func (r *reader) lastCriterion() *spec.AcceptanceCriterion {
	return &r.srd.AcceptanceCriteria[len(r.srd.AcceptanceCriteria)-1]
}

func (r *reader) storeDocField(key string, value *yaml.Node) error {
	switch key {
	case keyProblem:
		return assignString(value, &r.srd.Problem)
	case keyGoals:
		return assignStrings(value, &r.srd.Goals)
	case keyNonGoals:
		return assignStrings(value, &r.srd.NonGoals)
	case keyAssumptions:
		return assignStrings(value, &r.srd.Assumptions)
	case keyNeedsClarification:
		if value == nil {
			return nil
		}
		return value.Decode(&r.srd.NeedsClarification)
	case keyAcceptanceCriteria:
		if value != nil {
			return fmt.Errorf("acceptance criteria: content belongs under a tagged criterion")
		}
		return nil
	}
	r.srd.Extensions = append(r.srd.Extensions, newExtension(key, value))
	return nil
}

func newExtension(key string, value *yaml.Node) spec.Extension {
	if value == nil {
		return spec.Extension{Key: key, Value: yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"}}
	}
	return spec.Extension{Key: key, Value: *value}
}

func assignString(value *yaml.Node, target *string) error {
	if value == nil {
		return nil
	}
	if value.Kind != yaml.ScalarNode {
		return fmt.Errorf("expected text, found a %s", nodeKindName(value))
	}
	*target = value.Value
	return nil
}

func assignStrings(value *yaml.Node, target *[]string) error {
	if value == nil {
		return nil
	}
	if !markdown.IsStringList(value) {
		return fmt.Errorf("expected a list of text items, found a %s", nodeKindName(value))
	}
	for _, item := range value.Content {
		*target = append(*target, item.Value)
	}
	return nil
}

func nodeKindName(node *yaml.Node) string {
	switch node.Kind {
	case yaml.SequenceNode:
		return "list"
	case yaml.MappingNode:
		return "mapping"
	default:
		return "scalar"
	}
}
