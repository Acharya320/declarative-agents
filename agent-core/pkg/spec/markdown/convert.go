// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"bytes"
	"fmt"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

// Report says what a conversion did beyond moving the content: the
// canonical-form changes, the YAML comments it did not carry, the strings
// whose trailing newlines it dropped (the one whitespace normalization), and
// the findings the converted document still has.
type Report struct {
	Source     string
	Target     string
	Changes    []string
	Comments   []string
	Whitespace []string
	Findings   []Finding
}

// Lines prints the report, one line per entry, under the source path.
func (r Report) Lines() []string {
	lines := []string{r.Source + " -> " + r.Target}
	for _, group := range []struct {
		label   string
		entries []string
	}{{"canonical", r.Changes}, {"comment dropped", r.Comments}, {"whitespace", r.Whitespace}} {
		for _, entry := range group.entries {
			lines = append(lines, "  "+group.label+": "+entry)
		}
	}
	for _, f := range r.Findings {
		lines = append(lines, "  finding: "+f.String())
	}
	return lines
}

// MarkdownPath names the Markdown document a YAML source converts to: the
// same name with the .md extension, plus any literal suffix the Markdown
// location adds to the YAML one (an audit register's -register).
func (d *Declaration) MarkdownPath(source string) string {
	return strings.TrimSuffix(source, path.Ext(source)) + d.markdownSuffix() + ".md"
}

// SourcePath names the YAML document a Markdown document converts back to.
func (d *Declaration) SourcePath(markdownPath string) string {
	return strings.TrimSuffix(strings.TrimSuffix(markdownPath, ".md"), d.markdownSuffix()) + ".yaml"
}

func (d *Declaration) markdownSuffix() string {
	yamlBase := strings.TrimSuffix(path.Base(d.sourceLocation), path.Ext(d.sourceLocation))
	mdBase := strings.TrimSuffix(path.Base(d.Location), ".md")
	return strings.TrimPrefix(mdBase, yamlBase)
}

// ConvertYAML converts a YAML document to canonical Markdown and verifies
// the result: the Markdown must parse back to the canonical tree, differing
// at most in the trailing newlines it reports. A conversion that would lose
// anything else is an error, never a file.
func (d *Declaration) ConvertYAML(data []byte, source string) ([]byte, Report, error) {
	report := Report{Source: source, Target: d.MarkdownPath(source)}
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, report, fmt.Errorf("%s: %w", source, err)
	}
	report.Comments = droppedComments(&doc)
	canonical, changes := d.Canonicalize(&doc)
	report.Changes = changes
	md, err := d.Print(canonical)
	if err != nil {
		return nil, report, fmt.Errorf("%s: %w", source, err)
	}
	back, findings, err := d.Parse(md, report.Target)
	report.Findings = findings
	if err != nil || back == nil {
		return nil, report, fmt.Errorf("%s: printed Markdown does not parse: %v %v", source, err, findings)
	}
	if diff := Diff(canonical, back); diff != "" {
		return nil, report, fmt.Errorf("%s: conversion would change %s", source, diff)
	}
	report.Whitespace = trailingNewlines(canonical, back, "$")
	return md, report, nil
}

// ConvertMarkdown converts a Markdown document back to YAML for rollback;
// the license lines return as the head comment.
func (d *Declaration) ConvertMarkdown(data []byte, source string) ([]byte, Report, error) {
	report := Report{Source: source, Target: d.SourcePath(source)}
	doc, findings, err := d.Parse(data, source)
	report.Findings = findings
	if err != nil || doc == nil {
		return nil, report, fmt.Errorf("%s: %v %v", source, err, findings)
	}
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(doc); err != nil {
		return nil, report, fmt.Errorf("%s: %w", source, err)
	}
	if err := encoder.Close(); err != nil {
		return nil, report, fmt.Errorf("%s: %w", source, err)
	}
	return buf.Bytes(), report, nil
}

// droppedComments lists every non-license comment in a YAML tree with its
// line.
func droppedComments(node *yaml.Node) []string {
	var dropped []string
	for _, comment := range []string{node.HeadComment, node.LineComment, node.FootComment} {
		for _, line := range strings.Split(comment, "\n") {
			text := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "#"))
			if text != "" && !isLicenseComment(text) {
				dropped = append(dropped, fmt.Sprintf("line %d: %s", node.Line, text))
			}
		}
	}
	for _, child := range node.Content {
		dropped = append(dropped, droppedComments(child)...)
	}
	return dropped
}

// trailingNewlines names the strings of a that b holds without their
// trailing newlines; Diff has already shown the trees otherwise agree.
func trailingNewlines(a, b *yaml.Node, where string) []string {
	if a.Kind == yaml.DocumentNode && len(a.Content) > 0 && len(b.Content) > 0 {
		return trailingNewlines(a.Content[0], b.Content[0], where)
	}
	if a.Kind == yaml.ScalarNode {
		if a.Value != b.Value {
			return []string{where + ": trailing newline dropped"}
		}
		return nil
	}
	var found []string
	for i := range a.Content {
		if i < len(b.Content) {
			found = append(found, trailingNewlines(a.Content[i], b.Content[i], where+"/"+childName(a, i))...)
		}
	}
	return found
}
