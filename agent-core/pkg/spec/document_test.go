// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package spec

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec/markdown"
	"github.com/stretchr/testify/require"
)

// formatsFixture is the valid fixture corpus completed to satisfy the
// Markdown declarations, so every document converts.
const formatsFixture = "testdata/formats"

// taggedSRDComplete is taggedSRD with the fields the SRD declaration
// requires of a Markdown document.
const taggedSRDComplete = `id: srd090-example
title: Example
problem: Something.
goals:
  - G1: A goal.
requirements:
  R1:
    title: Group
    items:
      - R1.1: First requirement.
      - R1.2: Second requirement.
non_goals:
  - Nothing else.
acceptance_criteria:
  - id: AC1
    criterion: Works.
    traces: [R1.1]
`

// copyTree copies a fixture directory so a test can convert it in place.
func copyTree(t *testing.T, from string) string {
	t.Helper()
	to := t.TempDir()
	require.NoError(t, filepath.WalkDir(from, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(from, path)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		target := filepath.Join(to, rel)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o600)
	}))
	return to
}

// convertCorpus converts every declared corpus document under root to
// Markdown the way mage convert does, removing the YAML.
func convertCorpus(t *testing.T, root string) int {
	t.Helper()
	declarations, err := markdown.Default()
	require.NoError(t, err)
	converted := 0
	require.NoError(t, filepath.WalkDir(filepath.Join(root, DocsDir), func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !strings.HasSuffix(path, yamlExt) {
			return err
		}
		d, ok := declarations.ForSource(path)
		if !ok {
			return nil
		}
		data, err := os.ReadFile(path)
		require.NoError(t, err)
		md, report, err := d.ConvertYAML(data, path)
		require.NoError(t, err, path)
		require.NoError(t, os.WriteFile(report.Target, md, 0o600))
		converted++
		return os.Remove(path)
	}))
	return converted
}

func modelJSON(t *testing.T, value any) string {
	t.Helper()
	data, err := json.Marshal(value)
	require.NoError(t, err)
	return strings.ReplaceAll(string(data), `\n"`, `"`)
}

// TestMarkdownCorpusLoadsLikeYAML loads the valid fixture corpus and the same
// corpus converted to Markdown, and requires equal models and equal findings,
// with findings naming the Markdown files.
func TestMarkdownCorpusLoadsLikeYAML(t *testing.T) {
	yamlRoot := copyTree(t, formatsFixture)
	mdRoot := copyTree(t, formatsFixture)
	require.Positive(t, convertCorpus(t, mdRoot))

	fromYAML, err := LoadCorpus(yamlRoot)
	require.NoError(t, err)
	fromMarkdown, err := LoadCorpus(mdRoot)
	require.NoError(t, err)
	require.Empty(t, fromMarkdown.FormatFindings)

	for name, pair := range map[string][2]any{
		"srds": {fromYAML.SRDs, fromMarkdown.SRDs}, "use cases": {fromYAML.UseCases, fromMarkdown.UseCases},
		"test suites": {fromYAML.TestSuites, fromMarkdown.TestSuites}, "roadmap": {fromYAML.Roadmap, fromMarkdown.Roadmap},
		"spec index": {fromYAML.SpecIndex, fromMarkdown.SpecIndex},
	} {
		require.Equal(t, modelJSON(t, pair[0]), modelJSON(t, pair[1]), name)
	}
	require.Equal(t, len(fromYAML.DocSpecs), len(fromMarkdown.DocSpecs))
	for file := range fromMarkdown.Files {
		require.True(t, strings.HasSuffix(fromMarkdown.Files[file], markdownExt), file)
	}

	yamlFindings := Validate(mustGraph(t, fromYAML), fromYAML)
	mdFindings := Validate(mustGraph(t, fromMarkdown), fromMarkdown)
	require.Equal(t, len(yamlFindings), len(mdFindings))
	for i := range yamlFindings {
		require.Equal(t, yamlFindings[i].Check, mdFindings[i].Check)
		require.Equal(t, strings.TrimSuffix(yamlFindings[i].File, yamlExt),
			strings.TrimSuffix(mdFindings[i].File, markdownExt))
	}
}

// TestFormatFindingsAreCollected seeds a converted corpus with a document in
// both formats, a document out of canonical form, and a document whose
// structure is broken, and requires one finding each while the rest loads.
func TestFormatFindingsAreCollected(t *testing.T) {
	root := copyTree(t, formatsFixture)
	convertCorpus(t, root)
	srdDir := filepath.Join(root, SRDSubdir)
	yamlTwin, err := os.ReadFile(filepath.Join(formatsFixture, SRDSubdir, "srd001-auth.yaml"))
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(srdDir, "srd001-auth.yaml"), yamlTwin, 0o600))

	appendTo(t, filepath.Join(srdDir, "srd002-api.md"), "\n")
	breakMarkdown(t, filepath.Join(root, CFSubdir, "test-format.md"), "<!-- field: purpose -->\n")

	corpus, err := LoadCorpus(root, WithOptionalCorpus())
	require.NoError(t, err)
	byCheck := map[string][]Finding{}
	for _, f := range corpus.FormatFindings {
		byCheck[f.Check] = append(byCheck[f.Check], f)
	}
	require.Len(t, byCheck[checkDualFormat], 1)
	require.Equal(t, SRDSubdir+"/srd001-auth.md", byCheck[checkDualFormat][0].File)
	require.Len(t, byCheck[checkNotCanonical], 1)
	require.Equal(t, SRDSubdir+"/srd002-api.md", byCheck[checkNotCanonical][0].File)
	require.NotEmpty(t, byCheck[checkMarkdownFormat])
	require.Equal(t, "heading-untagged", byCheck[checkMarkdownFormat][0].CheckID)
	require.Contains(t, byCheck[checkMarkdownFormat][0].Message, "add the marker")
	require.Equal(t, CFSubdir+"/test-format.md", byCheck[checkMarkdownFormat][0].File)
	require.Len(t, corpus.SRDs, 3)
	require.Contains(t, corpus.DocSpecs, "sm-test-model")
	require.NotContains(t, corpus.DocSpecs, "cfg-test-format")
	require.NotEmpty(t, findingsFor(Validate(mustGraph(t, corpus), corpus), checkNotCanonical))
}

func appendTo(t *testing.T, path, text string) {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, append(data, text...), 0o600))
}

// TestRetiredIDsReadMarkdownAtATag freezes an SRD in Markdown at a tag and
// checks the head corpus against it, as after the migration.
func TestRetiredIDsReadMarkdownAtATag(t *testing.T) {
	declarations, err := markdown.Default()
	require.NoError(t, err)
	file := SRDSubdir + "/srd090-example.yaml"
	d, ok := declarations.ForSource(file)
	require.True(t, ok)
	md, report, err := d.ConvertYAML([]byte(taggedSRDComplete), file)
	require.NoError(t, err)
	tag := "v0.20261016.0"
	swapGitOutput(t, fakeTaggedTree(tag, map[string]string{report.Target: string(md)}))

	corpus := corpusFromSRDYAML(t, strings.Replace(taggedSRDComplete, "R1.1: First", "R1.3: First", 1))
	findings := findingsFor(checkRetiredIDs(corpus), "retired-id-missing")
	require.Len(t, findings, 1)
	require.Contains(t, findings[0].Message, "srd090-example R1.1")
}

func TestParentSRDFilesMatchEitherFormat(t *testing.T) {
	for file, want := range map[string]bool{
		SRDSubdir + "/srd033-monitor.yaml":               true,
		"agent-core/" + SRDSubdir + "/srd033-monitor.md": true,
		SRDSubdir + "/srd033-monitor-rest.md":            false,
		"docs/specs/use-cases/srd033-monitor.md":         false,
	} {
		require.Equal(t, want, isParentSRDFile(file, "srd033-monitor"), file)
	}
}

func mustGraph(t *testing.T, corpus *Corpus) *Graph {
	t.Helper()
	graph, err := BuildGraph(corpus)
	require.NoError(t, err)
	return graph
}

func breakMarkdown(t *testing.T, path, marker string) {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Contains(t, string(data), marker)
	require.NoError(t, os.WriteFile(path, []byte(strings.Replace(string(data), marker, "", 1)), 0o600))
}

// TestLoadErrorNamesTheFormatFinding breaks an indexed SRD: the index check
// fails the load, and the error names the Markdown finding behind it.
func TestLoadErrorNamesTheFormatFinding(t *testing.T) {
	root := copyTree(t, formatsFixture)
	convertCorpus(t, root)
	breakMarkdown(t, filepath.Join(root, SRDSubdir, "srd003-storage.md"), "<!-- field: problem -->\n")
	_, err := LoadCorpus(root)
	require.Error(t, err)
	require.Contains(t, err.Error(), "srd003-storage")
	require.Contains(t, err.Error(), "heading-untagged")
}
