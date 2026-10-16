// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec/markdown"
	"github.com/stretchr/testify/require"
)

// validFixture is the fixture corpus in Markdown, the form LoadCorpus reads.
const validFixture = "testdata/valid"

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

// writeCorpusDocument writes a corpus document, given as YAML at its former
// path, as the Markdown mage convert makes of it, and returns the Markdown
// path. Tests author fixtures in YAML; the corpus reads only Markdown.
func writeCorpusDocument(t *testing.T, root, rel, yamlText string) string {
	t.Helper()
	target, err := markdown.ConvertFile(filepath.Join(root, filepath.FromSlash(rel)), []byte(yamlText))
	require.NoError(t, err, rel)
	return target
}

// parseSRDYAMLFile decodes an SRD YAML form fixture, as a release tag that
// predates the migration holds it.
func parseSRDYAMLFile(path string) (SRD, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return SRD{}, err
	}
	return parseSRDBytes(data, path)
}

// parseUseCaseYAMLFile decodes a use case YAML form fixture.
func parseUseCaseYAMLFile(path string) (UseCase, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return UseCase{}, err
	}
	return parseUseCaseBytes(data, path)
}

// TestCorpusLoadsMarkdownOnly loads the fixture corpus and requires every
// document to come from its Markdown file with no format finding, and a
// YAML document beside it to be ignored: the working tree has no YAML front
// end, and the placement gate rejects the stray file (GH-2626).
func TestCorpusLoadsMarkdownOnly(t *testing.T) {
	root := copyTree(t, validFixture)
	corpus, err := LoadCorpus(root)
	require.NoError(t, err)
	require.Empty(t, corpus.FormatFindings)
	require.NotEmpty(t, corpus.Files)
	for file := range corpus.Files {
		require.True(t, strings.HasSuffix(corpus.Files[file], markdownExt), file)
	}

	require.NoError(t, os.WriteFile(filepath.Join(root, SRDSubdir, "srd009-stray.yaml"), []byte(taggedSRDComplete), 0o600))
	again, err := LoadCorpus(root)
	require.NoError(t, err)
	require.Equal(t, len(corpus.SRDs), len(again.SRDs))
	require.NotContains(t, again.SRDs, "srd090-example")

	_, err = ParseSRD(filepath.Join(root, SRDSubdir, "srd009-stray.yaml"))
	require.ErrorContains(t, err, "corpus documents are Markdown")
}

// TestFormatFindingsAreCollected seeds the corpus with a document out of
// canonical form and a document whose structure is broken, and requires one
// finding each while the rest loads.
func TestFormatFindingsAreCollected(t *testing.T) {
	root := copyTree(t, validFixture)
	srdDir := filepath.Join(root, SRDSubdir)
	appendTo(t, filepath.Join(srdDir, "srd002-api.md"), "\n")
	breakMarkdown(t, filepath.Join(root, CFSubdir, "test-format.md"), "<!-- field: purpose -->\n")

	corpus, err := LoadCorpus(root, WithOptionalCorpus())
	require.NoError(t, err)
	byCheck := map[string][]Finding{}
	for _, f := range corpus.FormatFindings {
		byCheck[f.Check] = append(byCheck[f.Check], f)
	}
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

// TestParentSRDFilesMatchEitherFormat keeps YAML for pinned parent trees,
// which may predate the migration (GH-2626 R3).
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
	root := copyTree(t, validFixture)
	breakMarkdown(t, filepath.Join(root, SRDSubdir, "srd003-storage.md"), "<!-- field: problem -->\n")
	_, err := LoadCorpus(root)
	require.Error(t, err)
	require.Contains(t, err.Error(), "srd003-storage")
	require.Contains(t, err.Error(), "heading-untagged")
}
