// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package spec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func clarificationCorpus(t *testing.T, srdYAML string) *Corpus {
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

func TestUnresolvedClarificationOnRequirementItem(t *testing.T) {
	corpus := clarificationCorpus(t, `id: srd091-example
title: Example
requirements:
  R1:
    title: Group
    items:
      - R1.1: Retry the fetch N times.
acceptance_criteria:
  - id: AC1
    criterion: Works.
    traces: [R1.1]
needs_clarification:
  - id: R1.1
    question: How many retries before giving up?
`)
	findings := checkUnresolvedClarifications(corpus)
	require.Len(t, findings, 1)
	require.Equal(t, "unresolved-clarification", findings[0].Check)
	require.Equal(t, "error", findings[0].Level)
	require.Contains(t, findings[0].Message, "srd091-example R1.1")
	require.Contains(t, findings[0].Message, "How many retries before giving up?")
	require.Equal(t, SRDSubdir+"/srd091-example.md", findings[0].File)
}

func TestUnresolvedClarificationOnAcceptanceCriterion(t *testing.T) {
	corpus := clarificationCorpus(t, `id: srd091-example
title: Example
requirements:
  R1:
    title: Group
    items:
      - R1.1: Retry the fetch.
acceptance_criteria:
  - id: AC1
    criterion: Works within a bound nobody has set.
    traces: [R1.1]
needs_clarification:
  - id: AC1
    question: What latency bound does the criterion assert?
`)
	findings := checkUnresolvedClarifications(corpus)
	require.Len(t, findings, 1)
	require.Contains(t, findings[0].Message, "srd091-example AC1")
	require.Contains(t, findings[0].Message, "latency bound")
}

func TestUnresolvedClarificationUnknownTargetStillFires(t *testing.T) {
	corpus := clarificationCorpus(t, `id: srd091-example
title: Example
requirements:
  R1:
    title: Group
    items:
      - R1.1: Retry the fetch.
acceptance_criteria:
  - id: AC1
    criterion: Works.
    traces: [R1.1]
needs_clarification:
  - id: R9.9
    question: Does this block anything?
`)
	findings := checkUnresolvedClarifications(corpus)
	require.Len(t, findings, 1)
	require.Contains(t, findings[0].Message, `unknown id "R9.9"`)
	require.Contains(t, findings[0].Message, "Does this block anything?")
}

func TestUnresolvedClarificationOnUseCaseAndTestSuite(t *testing.T) {
	corpus := &Corpus{
		RootDir: t.TempDir(),
		SRDs:    map[string]SRD{},
		UseCases: map[string]UseCase{
			"rel09.0-uc001-ex": {
				ID:              "rel09.0-uc001-ex",
				SuccessCriteria: []SuccessCriterion{{ID: "S1"}},
				NeedsClarification: []ClarificationEntry{
					{ID: "S1", Question: "Which actor confirms completion?"},
				},
			},
		},
		TestSuites: map[string]TestSuite{
			"test-rel09.0-ex": {
				ID:        "test-rel09.0-ex",
				TestCases: []TestCase{{Name: "roundtrip"}},
				NeedsClarification: []ClarificationEntry{
					{ID: "roundtrip", Question: "Hermetic or live backend?"},
				},
			},
		},
	}
	findings := checkUnresolvedClarifications(corpus)
	require.Len(t, findings, 2)
	require.Contains(t, findings[0].Message, "rel09.0-uc001-ex S1")
	require.Equal(t, UCSubdir+"/rel09.0-uc001-ex.md", findings[0].File)
	require.Contains(t, findings[1].Message, "test-rel09.0-ex roundtrip")
	require.Equal(t, TSSubdir+"/test-rel09.0-ex.md", findings[1].File)
}

func TestUnresolvedClarificationCleanCorpusStaysSilent(t *testing.T) {
	corpus := clarificationCorpus(t, `id: srd091-example
title: Example
requirements:
  R1:
    title: Group
    items:
      - R1.1: Retry the fetch 3 times.
acceptance_criteria:
  - id: AC1
    criterion: Works.
    traces: [R1.1]
`)
	require.Empty(t, checkUnresolvedClarifications(corpus))
}
