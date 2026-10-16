// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package specmd

import (
	"fmt"
	"strings"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec"
	"gopkg.in/yaml.v3"
)

// licenseHeader is the canonical header the copyright gate accepts for a
// Markdown file with frontmatter: two HTML comments after the closing fence.
var licenseHeader = []string{
	"<!-- Copyright (c) 2026 Nokia -->",
	"<!-- SPDX-License-Identifier: BSD-3-Clause -->",
}

// Field keys the printer and parser share for typed SRD content.
const (
	keyProblem            = "problem"
	keyGoals              = "goals"
	keyNonGoals           = "non_goals"
	keyAssumptions        = "assumptions"
	keyAcceptanceCriteria = "acceptance_criteria"
	keyNeedsClarification = "needs_clarification"
	keyEvidence           = "evidence"
	keyTraces             = "traces"
)

// fieldHeadings names the visible heading of each typed field; an untyped
// field's heading is derived from its key.
var fieldHeadings = map[string]string{
	keyProblem:            "Problem",
	keyGoals:              "Goals",
	keyNonGoals:           "Non-goals",
	keyAssumptions:        "Assumptions",
	keyAcceptanceCriteria: "Acceptance criteria",
	keyNeedsClarification: "Needs clarification",
	keyEvidence:           "Evidence",
	keyTraces:             "Traces",
}

func fieldHeading(key string) string {
	if heading, ok := fieldHeadings[key]; ok {
		return heading
	}
	words := strings.ReplaceAll(key, "_", " ")
	return strings.ToUpper(words[:1]) + words[1:]
}

// printer accumulates canonical Markdown lines.
type printer struct {
	lines []string
}

func (p *printer) add(lines ...string) { p.lines = append(p.lines, lines...) }

// section writes a marker, its heading, and the rendered value, followed by
// one blank line.
func (p *printer) section(m marker, level int, heading string, value *yaml.Node) error {
	p.add(m.String(), strings.Repeat("#", level)+" "+heading, "")
	if value == nil {
		return nil
	}
	content, err := renderNode(value)
	if err != nil {
		return fmt.Errorf("%s %s: %w", m.Tag, m.Value, err)
	}
	p.add(content...)
	p.add("")
	return nil
}

func (p *printer) field(level int, key string, value *yaml.Node) error {
	return p.section(marker{Tag: "field", Value: key}, level, fieldHeading(key), value)
}

// PrintMarkdown renders an SRD in canonical Markdown form.
func PrintMarkdown(srd *spec.SRD) ([]byte, error) {
	p := &printer{}
	if err := p.frontmatter(srd); err != nil {
		return nil, err
	}
	p.add(licenseHeader...)
	p.add("")
	steps := []func(*spec.SRD) error{p.opening, p.groups, p.closing, p.criteria, p.documentExtensions}
	for _, step := range steps {
		if err := step(srd); err != nil {
			return nil, err
		}
	}
	return []byte(strings.TrimRight(strings.Join(p.lines, "\n"), "\n") + "\n"), nil
}

// opening writes the problem and goals.
func (p *printer) opening(srd *spec.SRD) error {
	if srd.Problem != "" {
		if err := p.field(2, keyProblem, stringNode(srd.Problem)); err != nil {
			return err
		}
	}
	return p.stringList(2, keyGoals, srd.Goals)
}

// closing writes the non-goals, assumptions, and open clarifications.
func (p *printer) closing(srd *spec.SRD) error {
	if err := p.stringList(2, keyNonGoals, srd.NonGoals); err != nil {
		return err
	}
	if err := p.stringList(2, keyAssumptions, srd.Assumptions); err != nil {
		return err
	}
	if len(srd.NeedsClarification) == 0 {
		return nil
	}
	var node yaml.Node
	if err := node.Encode(srd.NeedsClarification); err != nil {
		return err
	}
	return p.field(2, keyNeedsClarification, &node)
}

func (p *printer) stringList(level int, key string, values []string) error {
	if len(values) == 0 {
		return nil
	}
	list := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for _, value := range values {
		list.Content = append(list.Content, stringNode(value))
	}
	return p.field(level, key, list)
}

func stringNode(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value}
}

func (p *printer) extensions(level int, extensions []spec.Extension) error {
	for i := range extensions {
		if err := p.field(level, extensions[i].Key, &extensions[i].Value); err != nil {
			return err
		}
	}
	return nil
}

func (p *printer) documentExtensions(srd *spec.SRD) error {
	return p.extensions(2, srd.Extensions)
}

// groups writes each requirement group with its own fields and items.
func (p *printer) groups(srd *spec.SRD) error {
	for _, groupID := range srd.OrderedGroups {
		group := srd.Requirements[groupID]
		heading := groupID
		if group.Title != "" {
			heading += " - " + group.Title
		}
		if err := p.section(marker{Tag: "req", Value: groupID}, 2, heading, nil); err != nil {
			return err
		}
		if err := p.extensions(3, group.Extensions); err != nil {
			return err
		}
		for i := range group.Items {
			if err := p.item(&group.Items[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *printer) item(item *spec.RequirementItem) error {
	m, err := itemMarker(item)
	if err != nil {
		return err
	}
	var text *yaml.Node
	if item.Text != "" {
		text = stringNode(item.Text)
	}
	if err := p.section(m, 3, item.ID, text); err != nil {
		return err
	}
	if item.Evidence != "" {
		if err := p.field(4, keyEvidence, stringNode(item.Evidence)); err != nil {
			return err
		}
	}
	return p.extensions(4, item.Extensions)
}

// criteria writes the acceptance criteria under one section heading.
func (p *printer) criteria(srd *spec.SRD) error {
	if len(srd.AcceptanceCriteria) == 0 {
		return nil
	}
	if err := p.field(2, keyAcceptanceCriteria, nil); err != nil {
		return err
	}
	for i := range srd.AcceptanceCriteria {
		criterion := &srd.AcceptanceCriteria[i]
		m, err := criterionMarker(criterion)
		if err != nil {
			return err
		}
		var text *yaml.Node
		if criterion.Criterion != "" {
			text = stringNode(criterion.Criterion)
		}
		if err := p.section(m, 3, criterion.ID, text); err != nil {
			return err
		}
		if err := p.stringList(4, keyTraces, criterion.Traces); err != nil {
			return err
		}
		if err := p.extensions(4, criterion.Extensions); err != nil {
			return err
		}
	}
	return nil
}
