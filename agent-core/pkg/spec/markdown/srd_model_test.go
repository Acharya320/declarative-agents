// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown_test

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec"
	"gopkg.in/yaml.v3"
)

// These helpers compare two SRD models field by field. They came from the
// prototype (GH-2588) with its round-trip tests.

// textNorm is the one whitespace normalization the Markdown front end
// applies: a YAML literal block ends text with a newline that Markdown prose
// cannot express, so text compares with trailing newlines trimmed.
func textNorm(s string) string { return strings.TrimRight(s, "\n") }

// modelDiffs compares two SRDs field by field and names every difference.
// With normalize set, text fields compare under textNorm; extension values
// always compare by decoded content, not YAML style.
func modelDiffs(a, b *spec.SRD, normalize bool) []string {
	var diffs []string
	text := func(s string) string {
		if normalize {
			return textNorm(s)
		}
		return s
	}
	check := func(where string, x, y any) {
		if !reflect.DeepEqual(x, y) {
			diffs = append(diffs, fmt.Sprintf("%s: %v != %v", where, x, y))
		}
	}
	check("id", a.ID, b.ID)
	check("title", a.Title, b.Title)
	check("problem", text(a.Problem), text(b.Problem))
	check("status/retired_in", a.Status+"|"+a.RetiredIn, b.Status+"|"+b.RetiredIn)
	check("replaced_by", fmt.Sprint(a.ReplacedBy), fmt.Sprint(b.ReplacedBy))
	check("goals", mapText(a.Goals, text), mapText(b.Goals, text))
	check("non_goals", mapText(a.NonGoals, text), mapText(b.NonGoals, text))
	check("assumptions", mapText(a.Assumptions, text), mapText(b.Assumptions, text))
	check("depends_on", fmt.Sprint(a.DependsOn), fmt.Sprint(b.DependsOn))
	check("derived_from", fmt.Sprint(a.DerivedFrom), fmt.Sprint(b.DerivedFrom))
	check("needs_clarification", fmt.Sprint(a.NeedsClarification), fmt.Sprint(b.NeedsClarification))
	check("groups", a.OrderedGroups, b.OrderedGroups)
	diffs = append(diffs, extensionDiffs("document", a.Extensions, b.Extensions, text)...)
	for _, id := range a.OrderedGroups {
		diffs = append(diffs, groupDiffs(id, a.Requirements[id], b.Requirements[id], text)...)
	}
	check("criteria count", len(a.AcceptanceCriteria), len(b.AcceptanceCriteria))
	for i := range a.AcceptanceCriteria {
		if i < len(b.AcceptanceCriteria) {
			diffs = append(diffs, criterionDiffs(&a.AcceptanceCriteria[i], &b.AcceptanceCriteria[i], text)...)
		}
	}
	return diffs
}

func groupDiffs(id string, a, b spec.RequirementGroup, text func(string) string) []string {
	var diffs []string
	if a.Title != b.Title {
		diffs = append(diffs, fmt.Sprintf("group %s title: %q != %q", id, a.Title, b.Title))
	}
	diffs = append(diffs, extensionDiffs("group "+id, a.Extensions, b.Extensions, text)...)
	if len(a.Items) != len(b.Items) {
		return append(diffs, fmt.Sprintf("group %s: %d items != %d", id, len(a.Items), len(b.Items)))
	}
	for i := range a.Items {
		x, y := a.Items[i], b.Items[i]
		if x.ID != y.ID || text(x.Text) != text(y.Text) || text(x.Evidence) != text(y.Evidence) ||
			x.Weight != y.Weight || x.Status != y.Status || x.RetiredIn != y.RetiredIn ||
			fmt.Sprint(x.ReplacedBy) != fmt.Sprint(y.ReplacedBy) || fmt.Sprint(x.DerivedFrom) != fmt.Sprint(y.DerivedFrom) {
			diffs = append(diffs, fmt.Sprintf("item %s differs: %+v != %+v", x.ID, x, y))
		}
		diffs = append(diffs, extensionDiffs("item "+x.ID, x.Extensions, y.Extensions, text)...)
	}
	return diffs
}

func criterionDiffs(a, b *spec.AcceptanceCriterion, text func(string) string) []string {
	var diffs []string
	if a.ID != b.ID || text(a.Criterion) != text(b.Criterion) || fmt.Sprint(a.Traces) != fmt.Sprint(b.Traces) ||
		a.Status != b.Status || a.RetiredIn != b.RetiredIn || fmt.Sprint(a.ReplacedBy) != fmt.Sprint(b.ReplacedBy) {
		diffs = append(diffs, fmt.Sprintf("criterion %s differs: %+v != %+v", a.ID, a, b))
	}
	return append(diffs, extensionDiffs("criterion "+a.ID, a.Extensions, b.Extensions, text)...)
}

func mapText(values []string, text func(string) string) []string {
	out := make([]string, len(values))
	for i, v := range values {
		out[i] = text(v)
	}
	return out
}

func extensionDiffs(where string, a, b []spec.Extension, text func(string) string) []string {
	if len(a) != len(b) {
		return []string{fmt.Sprintf("%s extensions: %v != %v", where, extKeys(a), extKeys(b))}
	}
	var diffs []string
	for i := range a {
		if a[i].Key != b[i].Key || !reflect.DeepEqual(decoded(&a[i].Value, text), decoded(&b[i].Value, text)) {
			diffs = append(diffs, fmt.Sprintf("%s extension %s differs", where, a[i].Key))
		}
	}
	return diffs
}

func extKeys(extensions []spec.Extension) []string {
	keys := make([]string, len(extensions))
	for i, e := range extensions {
		keys[i] = e.Key
	}
	return keys
}

// decoded reads a node's content, applying text to every string inside it.
func decoded(node *yaml.Node, text func(string) string) any {
	var value any
	if err := node.Decode(&value); err != nil {
		return "undecodable: " + err.Error()
	}
	return normalizeValue(value, text)
}

func normalizeValue(value any, text func(string) string) any {
	switch v := value.(type) {
	case string:
		return text(v)
	case []any:
		for i := range v {
			v[i] = normalizeValue(v[i], text)
		}
	case map[string]any:
		for k := range v {
			v[k] = normalizeValue(v[k], text)
		}
	}
	return value
}
