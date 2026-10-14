// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package spec

import (
	"errors"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const taggedSRD = `id: srd090-example
title: Example
requirements:
  R1:
    title: Group
    items:
      - R1.1: First requirement.
      - R1.2: Second requirement.
acceptance_criteria:
  - id: AC1
    criterion: Works.
    traces: [R1.1]
`

// swapGitOutput replaces the package git reader for one test, the way the
// external-reference tests swap their reads, so no repository is needed.
func swapGitOutput(t *testing.T, fn func(dir string, args ...string) (string, error)) {
	t.Helper()
	orig := gitOutput
	gitOutput = fn
	t.Cleanup(func() { gitOutput = orig })
}

// fakeTaggedTree serves describe/ls-tree/show against an in-memory file set.
// An empty tag refuses describe, modelling a repository without releases.
func fakeTaggedTree(tag string, files map[string]string) func(string, ...string) (string, error) {
	return func(dir string, args ...string) (string, error) {
		if len(args) == 0 {
			return "", errors.New("no git subcommand")
		}
		switch args[0] {
		case "describe":
			if tag == "" {
				return "", errors.New("no release tag reachable")
			}
			return tag, nil
		case "ls-tree":
			names := make([]string, 0, len(files))
			for name := range files {
				names = append(names, name)
			}
			sort.Strings(names)
			return strings.Join(names, "\n"), nil
		case "show":
			name := strings.TrimPrefix(args[1], tag+":./")
			content, ok := files[name]
			if !ok {
				return "", errors.New("path not in tagged tree: " + name)
			}
			return content, nil
		}
		return "", errors.New("unexpected git call: " + args[0])
	}
}

func corpusFromSRDYAML(t *testing.T, srdYAML string) *Corpus {
	t.Helper()
	srd, err := parseSRDBytes([]byte(srdYAML), "head")
	require.NoError(t, err)
	return &Corpus{
		RootDir:    t.TempDir(),
		SRDs:       map[string]SRD{srd.ID: srd},
		UseCases:   map[string]UseCase{},
		TestSuites: map[string]TestSuite{},
	}
}

func findingsFor(findings []Finding, check string) []Finding {
	var matched []Finding
	for _, f := range findings {
		if f.Check == check {
			matched = append(matched, f)
		}
	}
	return matched
}

func TestRetiredIDsNoTagEmitsNothing(t *testing.T) {
	swapGitOutput(t, fakeTaggedTree("", nil))
	corpus := corpusFromSRDYAML(t, taggedSRD)
	require.Empty(t, checkRetiredIDs(corpus))
}

func TestRetiredIDsRenumberedItemFails(t *testing.T) {
	tag := "v0.20260923.3"
	swapGitOutput(t, fakeTaggedTree(tag, map[string]string{
		SRDSubdir + "/srd090-example.yaml": taggedSRD,
	}))
	// R1.1 was renumbered to R1.3: the id from the tag is gone without a
	// retirement entry.
	corpus := corpusFromSRDYAML(t, `id: srd090-example
title: Example
requirements:
  R1:
    title: Group
    items:
      - R1.3: First requirement, renumbered.
      - R1.2: Second requirement.
acceptance_criteria:
  - id: AC1
    criterion: Works.
    traces: [R1.3]
`)
	findings := findingsFor(checkRetiredIDs(corpus), "retired-id-missing")
	require.Len(t, findings, 1)
	require.Equal(t, "error", findings[0].Level)
	require.Contains(t, findings[0].Message, "srd090-example R1.1")
	require.Contains(t, findings[0].Message, tag)
	require.Equal(t, SRDSubdir+"/srd090-example.yaml", findings[0].File)
}

func TestRetiredIDsDroppedDocumentsFail(t *testing.T) {
	tag := "v0.20260923.3"
	swapGitOutput(t, fakeTaggedTree(tag, map[string]string{
		SRDSubdir + "/srd090-example.yaml":  taggedSRD,
		UCSubdir + "/rel09.0-uc001-ex.yaml": "id: rel09.0-uc001-ex\ntitle: Example\n",
		TSSubdir + "/test-rel09.0-ex.yaml":  "id: test-rel09.0-ex\ntitle: Example\n",
	}))
	corpus := &Corpus{
		RootDir:    t.TempDir(),
		SRDs:       map[string]SRD{},
		UseCases:   map[string]UseCase{},
		TestSuites: map[string]TestSuite{},
	}
	findings := findingsFor(checkRetiredIDs(corpus), "retired-id-missing")
	require.Len(t, findings, 3)
	var ids []string
	for _, f := range findings {
		ids = append(ids, f.Message)
	}
	joined := strings.Join(ids, "\n")
	require.Contains(t, joined, "srd090-example")
	require.Contains(t, joined, "rel09.0-uc001-ex")
	require.Contains(t, joined, "test-rel09.0-ex")
}

func TestRetiredIDsRetirementWithPointerPasses(t *testing.T) {
	tag := "v0.20260923.3"
	swapGitOutput(t, fakeTaggedTree(tag, map[string]string{
		SRDSubdir + "/srd090-example.yaml": taggedSRD,
	}))
	corpus := corpusFromSRDYAML(t, `id: srd090-example
title: Example
requirements:
  R1:
    title: Group
    items:
      - R1.1: First requirement, superseded.
        status: retired
        retired_in: v0.20260924.0
        replaced_by: [R1.2]
      - R1.2: Second requirement.
acceptance_criteria:
  - id: AC1
    criterion: Works.
    traces: [R1.2]
`)
	require.Empty(t, findingsFor(checkRetiredIDs(corpus), "retired-id-missing"))
}

func TestRetiredIDsRetirementWithoutRetiredInFails(t *testing.T) {
	tag := "v0.20260923.3"
	swapGitOutput(t, fakeTaggedTree(tag, map[string]string{
		SRDSubdir + "/srd090-example.yaml": taggedSRD,
	}))
	corpus := corpusFromSRDYAML(t, `id: srd090-example
title: Example
requirements:
  R1:
    title: Group
    items:
      - R1.1: First requirement, superseded.
        status: retired
        replaced_by: [R1.2]
      - R1.2: Second requirement.
acceptance_criteria:
  - id: AC1
    criterion: Works.
    traces: [R1.2]
`)
	findings := findingsFor(checkRetiredIDs(corpus), "retired-id-missing")
	require.Len(t, findings, 1)
	require.Contains(t, findings[0].Message, "srd090-example R1.1")
	require.Contains(t, findings[0].Message, "retired_in")
}

func TestCitationOfRetiredIDWarns(t *testing.T) {
	corpus := corpusFromSRDYAML(t, `id: srd090-example
title: Example
requirements:
  R1:
    title: Group
    items:
      - R1.1: First requirement, superseded.
        status: retired
        retired_in: v0.20260924.0
        replaced_by: [R1.2]
      - R1.2: Second requirement.
acceptance_criteria:
  - id: AC1
    criterion: Works.
    traces: [R1.2]
`)
	corpus.UseCases["rel09.0-uc001-ex"] = UseCase{
		ID:          "rel09.0-uc001-ex",
		Touchpoints: []string{"T1: srd090-example R1.1 -- exercises the retired item"},
	}
	findings := findingsFor(checkCitationsOfRetiredIDs(corpus), "citation-of-retired-id")
	require.Len(t, findings, 1)
	require.Equal(t, "warning", findings[0].Level)
	require.Contains(t, findings[0].Message, "use case rel09.0-uc001-ex")
	require.Contains(t, findings[0].Message, "srd090-example R1.1")
	require.Contains(t, findings[0].Message, "replaced by R1.2")
}

func TestCitationOfLiveIDsStaysSilent(t *testing.T) {
	corpus := corpusFromSRDYAML(t, taggedSRD)
	corpus.UseCases["rel09.0-uc001-ex"] = UseCase{
		ID:          "rel09.0-uc001-ex",
		Touchpoints: []string{"T1: srd090-example R1.1 -- exercises a live item"},
	}
	require.Empty(t, checkCitationsOfRetiredIDs(corpus))
}
