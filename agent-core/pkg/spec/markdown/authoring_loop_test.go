// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

// seededSRD carries violations from every tier: an unknown frontmatter key
// (1), an untagged heading and a missing required field (2), a goal that
// does not match its item pattern (3), and a normative keyword in prose (4).
const seededSRD = `---
id: srd950-loop
title: Authoring loop
stauts: draft
---

## Problem

The store MUST keep every record.

<!-- field: goals -->
## Goals

- Keep every record.

<!-- req: R1 -->
## R1 - Storage

<!-- req: R1.1 -->
### R1.1

The store keeps every record.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

A record survives a restart.

<!-- field: traces -->
#### Traces

- R1.1
`

const seededName = "docs/specs/software-requirements/srd950-loop.md"

// TestAuthoringLoopFollowsTheHints drives the seeded document to zero
// findings the way an agent iterates write, audit, fix (#2583 R4): each
// round fixes the first finding using only what the finding says, its
// line, expected and found text, and hint. When a hint stops carrying what
// its fix needs, the round changes nothing and the test fails.
func TestAuthoringLoopFollowsTheHints(t *testing.T) {
	d := repositoryDeclarations(t)["srd"]
	lines := strings.Split(seededSRD, "\n")
	tiers := map[int]bool{}
	for round := 0; round < 20; round++ {
		findings := d.Diagnose([]byte(strings.Join(lines, "\n")), seededName)
		if len(findings) == 0 {
			for tier := 1; tier <= 4; tier++ {
				if !tiers[tier] {
					t.Errorf("the seeded document raised no tier %d finding", tier)
				}
			}
			return
		}
		f := findings[0]
		tiers[f.Tier] = true
		fixed, err := followHint(lines, f)
		if err != nil {
			t.Fatalf("round %d: %s: %v", round, f, err)
		}
		if strings.Join(fixed, "\n") == strings.Join(lines, "\n") {
			t.Fatalf("round %d: following the hint changed nothing: %s", round, f)
		}
		lines = fixed
	}
	t.Fatalf("still has findings after 20 rounds:\n%s", strings.Join(lines, "\n"))
}

var (
	didYouMean    = regexp.MustCompile(`did you mean "([^"]+)"\?`)
	fieldMarker   = regexp.MustCompile(`<!-- field: ([a-z_]+) -->`)
	quotedHeading = regexp.MustCompile(`^"(.+)" with no marker$`)
	itemPattern   = regexp.MustCompile(`^each item matching \^([A-Z]+)\[0-9\]\+: $`)
	quotedItem    = regexp.MustCompile(`^item \d+: "(.*)"$`)
	keyword       = regexp.MustCompile(`^"([A-Za-z ]+)" in `)
)

// followHint applies the fix one finding asks for. Line numbers are 1-based.
func followHint(lines []string, f Finding) ([]string, error) {
	switch f.Name {
	case "frontmatter-unknown-key":
		return renameKey(lines, f)
	case "heading-untagged":
		return tagHeading(lines, f)
	case "missing-required-field":
		return addField(lines, f)
	case "field-shape-mismatch":
		return listContent(lines, f)
	case "list-item-pattern":
		return prefixItem(lines, f)
	case "prose-containment":
		return reword(lines, f)
	}
	return nil, fmt.Errorf("no fix follows from finding %s", f.Name)
}

func renameKey(lines []string, f Finding) ([]string, error) {
	match := didYouMean.FindStringSubmatch(f.Hint)
	if match == nil {
		return nil, fmt.Errorf("hint names no key")
	}
	out := append([]string(nil), lines...)
	out[f.Line-1] = strings.Replace(out[f.Line-1], f.Found, match[1], 1)
	return out, nil
}

func tagHeading(lines []string, f Finding) ([]string, error) {
	match := quotedHeading.FindStringSubmatch(f.Found)
	if match == nil || !strings.Contains(f.Hint, "<!-- field: <key> -->") {
		return nil, fmt.Errorf("finding names no heading or no marker form")
	}
	key := strings.ReplaceAll(strings.ToLower(match[1]), " ", "_")
	return insert(lines, f.Line-1, "<!-- field: "+key+" -->"), nil
}

func addField(lines []string, f Finding) ([]string, error) {
	match := fieldMarker.FindStringSubmatch(f.Hint)
	if match == nil || !strings.Contains(f.Hint, "level-2 heading") {
		return nil, fmt.Errorf("hint names no marker or heading level")
	}
	heading := strings.ToUpper(match[1][:1]) + strings.ReplaceAll(match[1][1:], "_", " ")
	out := append([]string(nil), lines...)
	for len(out) > 0 && out[len(out)-1] == "" {
		out = out[:len(out)-1]
	}
	return append(out, "", match[0], "## "+heading, "", "To be written.", ""), nil
}

func listContent(lines []string, f Finding) ([]string, error) {
	if !strings.HasSuffix(f.Hint, " as list") {
		return nil, fmt.Errorf("hint names no list shape")
	}
	out := append([]string(nil), lines...)
	out[f.Line-1] = "- " + out[f.Line-1]
	return out, nil
}

func prefixItem(lines []string, f Finding) ([]string, error) {
	pattern := itemPattern.FindStringSubmatch(f.Expected)
	item := quotedItem.FindStringSubmatch(f.Found)
	if pattern == nil || item == nil {
		return nil, fmt.Errorf("finding names no item form or no item")
	}
	out := append([]string(nil), lines...)
	for i, line := range out {
		if line == "- "+item[1] {
			out[i] = "- " + pattern[1] + "1: " + item[1]
			return out, nil
		}
	}
	return nil, fmt.Errorf("item %q not found", item[1])
}

func reword(lines []string, f Finding) ([]string, error) {
	match := keyword.FindStringSubmatch(f.Found)
	if match == nil || !strings.Contains(f.Hint, "reword it as description") {
		return nil, fmt.Errorf("finding names no keyword or no rewording")
	}
	out := append([]string(nil), lines...)
	out[f.Line-1] = strings.Replace(out[f.Line-1], match[1]+" keep", "keeps", 1)
	return out, nil
}

func insert(lines []string, at int, line string) []string {
	out := append([]string(nil), lines[:at]...)
	out = append(out, line)
	return append(out, lines[at:]...)
}
