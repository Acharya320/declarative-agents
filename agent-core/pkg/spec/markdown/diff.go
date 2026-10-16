// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Diff names the first difference between two YAML trees under the
// round-trip laws (eng10): maps and lists compare in order, scalars by tag
// and value, comments and styles are ignored, and strings compare without
// their trailing newlines. It returns "" when the trees agree.
func Diff(a, b *yaml.Node) string {
	return nodeDiff(a, b, "$")
}

func nodeDiff(a, b *yaml.Node, path string) string {
	if a.Kind == yaml.DocumentNode && len(a.Content) > 0 {
		a = a.Content[0]
	}
	if b.Kind == yaml.DocumentNode && len(b.Content) > 0 {
		b = b.Content[0]
	}
	switch {
	case a.Kind != b.Kind:
		return path + ": kind differs"
	case a.Kind == yaml.ScalarNode:
		return scalarDiff(a, b, path)
	case len(a.Content) != len(b.Content):
		return path + ": length differs"
	}
	for i := range a.Content {
		if diff := nodeDiff(a.Content[i], b.Content[i], path+"/"+childName(a, i)); diff != "" {
			return diff
		}
	}
	return ""
}

func scalarDiff(a, b *yaml.Node, path string) string {
	if a.ShortTag() != b.ShortTag() {
		return path + ": " + a.ShortTag() + " " + a.Value + " became " + b.ShortTag() + " " + b.Value
	}
	if a.ShortTag() != "!!null" && strings.TrimRight(a.Value, "\n") != strings.TrimRight(b.Value, "\n") {
		return path + ": " + a.Value + " became " + b.Value
	}
	return ""
}

func childName(parent *yaml.Node, i int) string {
	if parent.Kind == yaml.MappingNode && i%2 == 1 {
		return parent.Content[i-1].Value
	}
	return fmt.Sprintf("[%d]", i)
}
