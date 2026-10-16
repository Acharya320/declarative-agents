// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

// markdown is the CommonMark + GFM-table parser every read goes through.
var markdown = goldmark.New(goldmark.WithExtensions(extension.Table))

// A block is one top-level Markdown block with its line span (zero-based,
// inclusive) in the parsed source. Spans run from a block's first line to the
// line before the next block, minus trailing blank lines, so content between
// two structural lines is reproduced exactly.
type Block struct {
	Node  ast.Node
	Start int
	End   int
}

// topLevelBlocks parses src and returns its top-level blocks in order.
func topLevelBlocks(src []byte) ([]Block, []string) {
	doc := markdown.Parser().Parse(text.NewReader(src))
	lines := strings.Split(string(src), "\n")
	starts := lineStarts(src)
	var blocks []Block
	for node := doc.FirstChild(); node != nil; node = node.NextSibling() {
		if start, ok := startLine(node, src, starts); ok {
			blocks = append(blocks, Block{Node: node, Start: start})
		}
	}
	for i := range blocks {
		end := len(lines) - 1
		if i+1 < len(blocks) {
			end = blocks[i+1].Start - 1
		}
		for end > blocks[i].Start && strings.TrimSpace(lines[end]) == "" {
			end--
		}
		blocks[i].End = end
	}
	return blocks, lines
}

func lineStarts(src []byte) []int {
	starts := []int{0}
	for i, c := range src {
		if c == '\n' {
			starts = append(starts, i+1)
		}
	}
	return starts
}

func lineOf(offset int, starts []int) int {
	return sort.SearchInts(starts, offset+1) - 1
}

// startLine finds the first source line of a block: its own first segment,
// a fenced block's fence (the line before its info or content), or the first
// segment found among its descendants.
func startLine(node ast.Node, src []byte, starts []int) (int, bool) {
	if fenced, ok := node.(*ast.FencedCodeBlock); ok {
		if fenced.Info != nil {
			return lineOf(fenced.Info.Segment.Start, starts), true
		}
		if fenced.Lines().Len() > 0 {
			return lineOf(fenced.Lines().At(0).Start, starts) - 1, true
		}
	}
	if node.Type() == ast.TypeBlock && node.Lines().Len() > 0 {
		return lineOf(node.Lines().At(0).Start, starts), true
	}
	if textNode, ok := node.(*ast.Text); ok {
		return lineOf(textNode.Segment.Start, starts), true
	}
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if line, ok := startLine(child, src, starts); ok {
			return line, true
		}
	}
	return 0, false
}

// classifyContent names the shape a field's Markdown takes: a lone fenced
// yaml block, a lone table, a lone dash list, or prose.
func classifyContent(content string) string {
	src := []byte(content)
	doc := markdown.Parser().Parse(text.NewReader(src))
	first := doc.FirstChild()
	if first == nil || first.NextSibling() != nil {
		return shapeProse
	}
	switch node := first.(type) {
	case *ast.FencedCodeBlock:
		if string(node.Language(src)) == "yaml" {
			return shapeYAML
		}
	case *east.Table:
		return shapeTable
	case *ast.List:
		if !node.IsOrdered() && node.Marker == '-' {
			return shapeList
		}
	}
	return shapeProse
}

// headingOf reports a block's ATX heading level and text.
func headingOf(b Block, lines []string) (int, string, bool) {
	heading, ok := b.Node.(*ast.Heading)
	if !ok {
		return 0, "", false
	}
	line := strings.TrimSpace(lines[b.Start])
	return heading.Level, strings.TrimSpace(strings.TrimLeft(line, "#")), true
}

// commentLine reports a single-line HTML comment block and its text.
func commentLine(b Block, lines []string) (string, bool) {
	if _, ok := b.Node.(*ast.HTMLBlock); !ok || b.Start != b.End {
		return "", false
	}
	line := strings.TrimSpace(lines[b.Start])
	return line, strings.HasPrefix(line, markerOpen) && strings.HasSuffix(line, markerClose)
}

// spanText returns the source lines of a run of blocks, verbatim.
func spanText(lines []string, start, end int) string {
	if start > end {
		return ""
	}
	return strings.Join(lines[start:end+1], "\n")
}
