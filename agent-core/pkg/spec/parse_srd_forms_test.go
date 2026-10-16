// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package spec

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func extensionKeys(extensions []Extension) []string {
	keys := make([]string, 0, len(extensions))
	for _, extension := range extensions {
		keys = append(keys, extension.Key)
	}
	return keys
}

func TestParseSRDObjectFormItems(t *testing.T) {
	t.Parallel()
	srd, err := parseSRDBytes([]byte(`id: srd093-objects
title: Objects
requirements:
  R1:
    title: Contract
    applicability: Applies to every panel.
    items:
      - id: R1.1
        text: |
          The contract shall be the closed set below.
        endpoints:
          - path: GET /monitor/state
            owner: srd033 R3.9
        evidence: |
          A kit test fails on an endpoint outside the set.
      - id: R1.2
        text: Panels read only owner-defined fields.
        status: retired
        retired_in: v0.20261001.0
acceptance_criteria:
  - id: AC1
    criterion: Works.
    traces: [R1.1]
    tests: [TestContract]
`), "objects")
	require.NoError(t, err)

	group := srd.Requirements["R1"]
	require.Equal(t, "Contract", group.Title)
	require.Equal(t, []string{"applicability"}, extensionKeys(group.Extensions))
	require.Len(t, group.Items, 2)

	first := group.Items[0]
	require.Equal(t, "R1.1", first.ID)
	require.Equal(t, "The contract shall be the closed set below.\n", first.Text)
	require.Equal(t, "A kit test fails on an endpoint outside the set.\n", first.Evidence)
	require.Equal(t, []string{"endpoints"}, extensionKeys(first.Extensions))
	require.Equal(t, 1, first.Weight)

	second := group.Items[1]
	require.Equal(t, "R1.2", second.ID)
	require.Equal(t, "retired", second.Status)
	require.Equal(t, "v0.20261001.0", second.RetiredIn)
	require.Empty(t, second.Extensions)

	require.Equal(t, []string{"tests"}, extensionKeys(srd.AcceptanceCriteria[0].Extensions))
	require.Equal(t, []string{"R1.1"}, srd.AcceptanceCriteria[0].Traces)
}

func TestParseSRDListedGroupsKeepListOrder(t *testing.T) {
	t.Parallel()
	srd, err := parseSRDBytes([]byte(`id: srd094-listed
title: Listed
requirements:
  - id: R2
    title: Second
    items:
      - id: R2.1
        text: Comes first in the file.
  - id: R1
    title: First
    items:
      - R1.1: Tagged form inside a listed group.
`), "listed")
	require.NoError(t, err)
	require.Equal(t, []string{"R2", "R1"}, srd.OrderedGroups)
	require.Equal(t, "R2.1", srd.Requirements["R2"].Items[0].ID)
	require.Equal(t, "R1.1", srd.Requirements["R1"].Items[0].ID)
	require.Empty(t, srd.Requirements["R2"].Extensions, "a listed group's id is typed, not an extension")
}

func TestParseSRDTaggedItemsAreUnchanged(t *testing.T) {
	t.Parallel()
	srd, err := parseSRDBytes([]byte(`id: srd095-tagged
title: Tagged
requirements:
  R1:
    title: Group
    items:
      - R1.1: Plain tagged requirement.
      - R1.2: Weighted requirement.
        weight: 3
`), "tagged")
	require.NoError(t, err)
	items := srd.Requirements["R1"].Items
	require.Equal(t, "R1.1", items[0].ID)
	require.Equal(t, "Plain tagged requirement.", items[0].Text)
	require.Equal(t, "R1.2", items[1].ID)
	require.Equal(t, 3, items[1].Weight)
	require.Empty(t, items[0].Extensions)
	require.Empty(t, items[1].Extensions)
}

func TestParseSRDKeepsUntypedTopLevelKeysInOrder(t *testing.T) {
	t.Parallel()
	srd, err := parseSRDBytes([]byte(`id: srd096-extensions
title: Extensions
context: |
  Why this exists.
problem: A problem.
implementation_status: partial
related_documents:
  - docs/ARCHITECTURE.yaml
`), "extensions")
	require.NoError(t, err)
	require.Equal(t, "A problem.", srd.Problem)
	require.Equal(t, []string{"context", "implementation_status", "related_documents"},
		extensionKeys(srd.Extensions))
	require.Equal(t, "partial", srd.Extensions[1].Value.Value)
}

func TestParseSRDReadsSRD004ObjectForm(t *testing.T) {
	t.Parallel()
	// The corpus SRD moved to Markdown (GH-2584); this fixture keeps the YAML
	// object form it was written in.
	srd, err := parseSRDYAMLFile(filepath.Join("testdata", "forms", "srd004-declarative-ux.yaml"))
	require.NoError(t, err)
	items := srd.Requirements["R1"].Items
	require.Equal(t, "R1.1", items[0].ID)
	require.Contains(t, items[0].Text, "closed set of endpoints")
	require.Contains(t, items[0].Evidence, "one constant")
	require.Equal(t, []string{"endpoints"}, extensionKeys(items[0].Extensions))
	require.Equal(t, []string{"applicability"}, extensionKeys(srd.Extensions)[:1])
}

func TestParseSRDReadsAgentArchitectureListedGroups(t *testing.T) {
	t.Parallel()
	// The corpus SRD moved to Markdown (GH-2584); this fixture keeps the
	// listed-group YAML form it was written in.
	srd, err := parseSRDYAMLFile(filepath.Join("testdata", "forms", "srd001-helm-deployment.yaml"))
	require.NoError(t, err)
	require.Equal(t, []string{"R1", "R2", "R3"}, srd.OrderedGroups)
	require.Equal(t, "R1.1", srd.Requirements["R1"].Items[0].ID)
}
