// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package spec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func priorityCorpus(priorities map[string]string) *Corpus {
	corpus := &Corpus{
		SRDs:       map[string]SRD{},
		UseCases:   map[string]UseCase{},
		TestSuites: map[string]TestSuite{},
	}
	for id, priority := range priorities {
		corpus.UseCases[id] = UseCase{ID: id, Priority: priority}
	}
	return corpus
}

func TestInvalidPriorityValidAndAbsentStaySilent(t *testing.T) {
	corpus := priorityCorpus(map[string]string{
		"rel09.0-uc001-ranked":   "P1",
		"rel09.0-uc002-ranked":   "P2",
		"rel09.0-uc003-ranked":   "P3",
		"rel09.0-uc004-unranked": "",
	})
	require.Empty(t, checkInvalidPriorities(corpus))
}

func TestInvalidPriorityOutsideVocabularyFails(t *testing.T) {
	corpus := priorityCorpus(map[string]string{
		"rel09.0-uc001-bad":   "P4",
		"rel09.0-uc002-lower": "p1",
	})
	findings := checkInvalidPriorities(corpus)
	require.Len(t, findings, 2)
	for _, f := range findings {
		require.Equal(t, "invalid-priority", f.Check)
		require.Equal(t, "error", f.Level)
	}
	require.Contains(t, findings[0].Message, `priority "P4"`)
	require.Equal(t, UCSubdir+"/rel09.0-uc001-bad.yaml", findings[0].File)
	require.Contains(t, findings[1].Message, `priority "p1"`)
}

func TestParseUseCaseSpecQualityFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rel09.0-uc001-ex.yaml")
	require.NoError(t, os.WriteFile(path, []byte(`id: rel09.0-uc001-ex
title: Example
priority: P2
assumptions:
  - A1: The backend is reachable.
edge_cases:
  - E1: The backend drops the connection mid-transfer.
success_criteria:
  - id: S1
    criterion: The transfer completes within one retry.
`), 0o644))
	uc, err := ParseUseCase(path)
	require.NoError(t, err)
	require.Equal(t, "P2", uc.Priority)
	require.Equal(t, []string{"A1: The backend is reachable."}, uc.Assumptions)
	require.Equal(t, []string{"E1: The backend drops the connection mid-transfer."}, uc.EdgeCases)
}

func TestParseSRDAssumptions(t *testing.T) {
	srd, err := parseSRDBytes([]byte(`id: srd092-example
title: Example
assumptions:
  - A1: Exactly one writer holds the lease.
requirements:
  R1:
    title: Group
    items:
      - R1.1: Acquire the lease before writing.
`), "head")
	require.NoError(t, err)
	require.Equal(t, []string{"A1: Exactly one writer holds the lease."}, srd.Assumptions)
}
