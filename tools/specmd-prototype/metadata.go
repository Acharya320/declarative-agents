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

// frontmatter is the document-level machine layer: identity, retirement, and
// the structured references a reader does not read as prose.
type frontmatter struct {
	ID          string                   `yaml:"id"`
	Title       string                   `yaml:"title"`
	Status      string                   `yaml:"status,omitempty"`
	RetiredIn   string                   `yaml:"retired_in,omitempty"`
	ReplacedBy  []string                 `yaml:"replaced_by,omitempty"`
	DependsOn   []spec.Dependency        `yaml:"depends_on,omitempty"`
	DerivedFrom []spec.ExternalReference `yaml:"derived_from,omitempty"`
}

func (p *printer) frontmatter(srd *spec.SRD) error {
	data, err := yaml.Marshal(frontmatter{
		ID: srd.ID, Title: srd.Title, Status: srd.Status, RetiredIn: srd.RetiredIn,
		ReplacedBy: srd.ReplacedBy, DependsOn: srd.DependsOn, DerivedFrom: srd.DerivedFrom,
	})
	if err != nil {
		return fmt.Errorf("frontmatter: %w", err)
	}
	p.add("---")
	p.add(strings.Split(strings.TrimRight(string(data), "\n"), "\n")...)
	p.add("---")
	return nil
}

// splitFrontmatter separates the frontmatter from the body and reports how
// many lines the frontmatter occupied, so body line numbers can be mapped
// back to the file.
func splitFrontmatter(data []byte) (frontmatter, []byte, int, error) {
	text := string(data)
	if !strings.HasPrefix(text, "---\n") {
		return frontmatter{}, nil, 0, fmt.Errorf("line 1: expected frontmatter opening \"---\"")
	}
	closing := strings.Index(text[4:], "\n---\n")
	if closing < 0 {
		return frontmatter{}, nil, 0, fmt.Errorf("frontmatter: no closing \"---\"")
	}
	head := text[4 : 4+closing]
	body := text[4+closing+5:]
	var fm frontmatter
	if err := yaml.Unmarshal([]byte(head), &fm); err != nil {
		return frontmatter{}, nil, 0, fmt.Errorf("frontmatter: %w", err)
	}
	return fm, []byte(body), strings.Count(head, "\n") + 3, nil
}

func refsToJSON(refs []spec.ExternalReference) []markdown.Reference {
	out := make([]markdown.Reference, len(refs))
	for i, r := range refs {
		out[i] = markdown.Reference(r)
	}
	return out
}

func refsFromJSON(refs []markdown.Reference) []spec.ExternalReference {
	out := make([]spec.ExternalReference, len(refs))
	for i, r := range refs {
		out[i] = spec.ExternalReference(r)
	}
	return out
}

// retirementAttrs encodes the retirement and provenance fields shared by
// items and criteria, omitting the ones that are unset.
func retirementAttrs(status, retiredIn string, replacedBy []string, derived []spec.ExternalReference) ([]markdown.Attr, error) {
	values := []struct {
		key   string
		value any
		set   bool
	}{
		{"status", status, status != ""},
		{"retired_in", retiredIn, retiredIn != ""},
		{"replaced_by", replacedBy, len(replacedBy) > 0},
		{"derived_from", refsToJSON(derived), len(derived) > 0},
	}
	var attrs []markdown.Attr
	for _, v := range values {
		if !v.set {
			continue
		}
		encoded, err := markdown.EncodeAttr(v.value)
		if err != nil {
			return nil, fmt.Errorf("attribute %s: %w", v.key, err)
		}
		attrs = append(attrs, markdown.Attr{Key: v.key, Value: encoded})
	}
	return attrs, nil
}

func itemMarker(item *spec.RequirementItem) (markdown.Marker, error) {
	m := markdown.Marker{Tag: "req", Value: item.ID}
	if item.Weight != 1 {
		m.Attrs = append(m.Attrs, markdown.Attr{Key: "weight", Value: fmt.Sprint(item.Weight)})
	}
	attrs, err := retirementAttrs(item.Status, item.RetiredIn, item.ReplacedBy, item.DerivedFrom)
	if err != nil {
		return markdown.Marker{}, fmt.Errorf("item %s: %w", item.ID, err)
	}
	m.Attrs = append(m.Attrs, attrs...)
	return m, nil
}

func criterionMarker(criterion *spec.AcceptanceCriterion) (markdown.Marker, error) {
	m := markdown.Marker{Tag: "ac", Value: criterion.ID}
	attrs, err := retirementAttrs(criterion.Status, criterion.RetiredIn, criterion.ReplacedBy, criterion.DerivedFrom)
	if err != nil {
		return markdown.Marker{}, fmt.Errorf("criterion %s: %w", criterion.ID, err)
	}
	m.Attrs = append(m.Attrs, attrs...)
	return m, nil
}

// applyItemAttrs decodes an item marker's attributes onto the item.
func applyItemAttrs(m markdown.Marker, item *spec.RequirementItem) error {
	item.Weight = 1
	for _, a := range m.Attrs {
		var err error
		switch a.Key {
		case "weight":
			err = markdown.DecodeAttr(a.Value, &item.Weight)
		default:
			err = applyRetirementAttr(a, &item.Status, &item.RetiredIn, &item.ReplacedBy, &item.DerivedFrom)
		}
		if err != nil {
			return fmt.Errorf("req %s: %w", m.Value, err)
		}
	}
	return nil
}

func applyCriterionAttrs(m markdown.Marker, criterion *spec.AcceptanceCriterion) error {
	for _, a := range m.Attrs {
		err := applyRetirementAttr(a, &criterion.Status, &criterion.RetiredIn,
			&criterion.ReplacedBy, &criterion.DerivedFrom)
		if err != nil {
			return fmt.Errorf("ac %s: %w", m.Value, err)
		}
	}
	return nil
}

func applyRetirementAttr(a markdown.Attr, status, retiredIn *string, replacedBy *[]string, derived *[]spec.ExternalReference) error {
	switch a.Key {
	case "status":
		return markdown.DecodeAttr(a.Value, status)
	case "retired_in":
		return markdown.DecodeAttr(a.Value, retiredIn)
	case "replaced_by":
		return markdown.DecodeAttr(a.Value, replacedBy)
	case "derived_from":
		var refs []markdown.Reference
		if err := markdown.DecodeAttr(a.Value, &refs); err != nil {
			return err
		}
		*derived = refsFromJSON(refs)
		return nil
	default:
		return fmt.Errorf("unknown attribute %q", a.Key)
	}
}
