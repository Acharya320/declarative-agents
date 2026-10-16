// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package filesystem

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestParsedModeReadsAMarkdownCorpusDocumentAsItsTree serves the road map in
// the parsed_yaml mode: a Markdown corpus document parses into the tree its
// declaration gives it, releases included, not into its frontmatter alone
// (GH-2629).
func TestParsedModeReadsAMarkdownCorpusDocumentAsItsTree(t *testing.T) {
	t.Parallel()
	abs, err := filepath.Abs(filepath.Join("..", "..", "..", "docs", "road-map.md"))
	require.NoError(t, err)
	data, err := os.ReadFile(abs)
	require.NoError(t, err)
	detail, err := buildResourceDetail(abs, "road-map.md", data, int64(len(data)),
		ResourceDefinition{Modes: []string{"raw_markdown", "parsed_yaml"}})
	require.NoError(t, err)
	parsed, ok := detail.Parsed.(map[string]interface{})
	require.True(t, ok, "parsed = %T", detail.Parsed)
	require.NotEmpty(t, parsed["id"])
	releases, ok := parsed["releases"].([]interface{})
	require.True(t, ok && len(releases) > 0, "releases = %#v", parsed["releases"])
	require.Equal(t, string(data), detail.Raw)
}

// TestParsedModeRejectsABrokenCorpusDocument reports a Markdown corpus
// document whose structure does not parse instead of serving its frontmatter.
func TestParsedModeRejectsABrokenCorpusDocument(t *testing.T) {
	t.Parallel()
	abs := filepath.Join(t.TempDir(), "docs", "road-map.md")
	data := []byte("---\nid: roadmap\ntitle: Road Map\n---\n\n## Untagged heading\n\nText.\n")
	_, err := buildResourceDetail(abs, "road-map.md", data, int64(len(data)),
		ResourceDefinition{Modes: []string{"parsed_yaml"}})
	require.Error(t, err)
}

// TestParsedModeServesNativeMarkdownRaw serves a guide, which no declaration
// covers, with its raw text and an empty parsed mapping, rather than failing
// to parse prose as YAML.
func TestParsedModeServesNativeMarkdownRaw(t *testing.T) {
	t.Parallel()
	abs := filepath.Join(t.TempDir(), "docs", "guides", "operator.md")
	data := []byte("# Operator guide\n\nRun: the agent with --profile.\n- a: b\nc\n")
	detail, err := buildResourceDetail(abs, "guides/operator.md", data, int64(len(data)),
		ResourceDefinition{Modes: []string{"raw_markdown", "parsed_yaml"}})
	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{}, detail.Parsed)
	require.Equal(t, string(data), detail.Raw)
}
