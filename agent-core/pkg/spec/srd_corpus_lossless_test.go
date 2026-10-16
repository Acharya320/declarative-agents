// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package spec

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// itemMetadataKeys are the item keys parseItemMetadata types.
var itemMetadataKeys = map[string]bool{
	"weight": true, "status": true, "retired_in": true, "replaced_by": true, "derived_from": true,
}

// corpusSRDPaths lists every tracked SRD the model must carry losslessly:
// the corpus outside testdata, plus the valid parser fixtures.
func corpusSRDPaths(t *testing.T, repo string) []string {
	t.Helper()
	out, err := exec.Command("git", "-C", repo, "ls-files").Output()
	require.NoError(t, err)
	var paths []string
	for _, path := range strings.Split(string(out), "\n") {
		if !strings.Contains(path, "/docs/specs/software-requirements/srd") || !strings.HasSuffix(path, ".yaml") {
			continue
		}
		if strings.Contains(path, "/testdata/") && !strings.Contains(path, "pkg/spec/testdata/valid/") {
			continue
		}
		paths = append(paths, path)
	}
	return paths
}

// TestSRDModelIsLosslessOverCorpus proves yaml -> model loses nothing: every
// source key at document, group, item, and criterion level is either bound to
// a typed field or kept, in order and with an identical value, as an
// extension (GH-2586).
func TestSRDModelIsLosslessOverCorpus(t *testing.T) {
	t.Parallel()
	repo := filepath.Dir(repoAgentCoreRoot(t))
	paths := corpusSRDPaths(t, repo)
	require.NotEmpty(t, paths)
	var losses []string
	for _, path := range paths {
		data, err := os.ReadFile(filepath.Join(repo, path))
		require.NoError(t, err)
		srd, err := parseSRDBytes(data, path)
		if err != nil {
			losses = append(losses, fmt.Sprintf("%s: does not parse: %v", path, err))
			continue
		}
		var doc yaml.Node
		require.NoError(t, yaml.Unmarshal(data, &doc))
		losses = append(losses, documentLosses(path, doc.Content[0], &srd)...)
	}
	require.Empty(t, losses, "model loses source content:\n%s", strings.Join(losses, "\n"))
}

func documentLosses(path string, root *yaml.Node, srd *SRD) []string {
	losses := extensionLosses(path, "document", root, srdTypedKeys, srd.Extensions)
	requirements := mappingValue(root, "requirements")
	for _, groupID := range srd.OrderedGroups {
		groupNode := sourceGroup(requirements, groupID)
		group := srd.Requirements[groupID]
		typed := mappedGroupKeys
		if requirements.Kind == yaml.SequenceNode {
			typed = listedGroupKeys
		}
		losses = append(losses, extensionLosses(path, "group "+groupID, groupNode, typed, group.Extensions)...)
		items := mappingValue(groupNode, "items")
		for index, item := range group.Items {
			losses = append(losses, itemLosses(path, items.Content[index], item)...)
		}
	}
	acceptance := mappingValue(root, "acceptance_criteria")
	for index, criterion := range srd.AcceptanceCriteria {
		losses = append(losses, extensionLosses(path, "criterion "+criterion.ID,
			acceptance.Content[index], acceptanceCriterionKeys, criterion.Extensions)...)
	}
	return losses
}

func itemLosses(path string, node *yaml.Node, item RequirementItem) []string {
	where := "item " + item.ID
	if id := mappingScalar(node, "id"); id != "" {
		var losses []string
		if item.ID != id || item.Text != mappingScalar(node, "text") || item.Evidence != mappingScalar(node, "evidence") {
			losses = append(losses, fmt.Sprintf("%s: %s: id, text, or evidence differs from source", path, where))
		}
		typed := map[string]bool{"id": true, "text": true, "evidence": true}
		for key := range itemMetadataKeys {
			typed[key] = true
		}
		return append(losses, extensionLosses(path, where, node, typed, item.Extensions)...)
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		key := node.Content[i].Value
		if itemMetadataKeys[key] {
			continue
		}
		if key != item.ID || node.Content[i+1].Value != item.Text {
			return []string{fmt.Sprintf("%s: %s: tagged id or text differs from source key %s", path, where, key)}
		}
		break
	}
	return nil
}

// extensionLosses compares a source mapping's untyped keys, in order and by
// value, with the extensions the model kept for it.
func extensionLosses(path, where string, node *yaml.Node, typed map[string]bool, kept []Extension) []string {
	want := untypedKeys(node, typed)
	if len(want) != len(kept) {
		return []string{fmt.Sprintf("%s: %s: source has untyped keys %v, model kept %v",
			path, where, extensionKeys(want), extensionKeys(kept))}
	}
	var losses []string
	for index := range want {
		if want[index].Key != kept[index].Key || nodeText(&want[index].Value) != nodeText(&kept[index].Value) {
			losses = append(losses, fmt.Sprintf("%s: %s: key %s lost or changed", path, where, want[index].Key))
		}
	}
	return losses
}

func mappingValue(node *yaml.Node, key string) *yaml.Node {
	if node != nil && node.Kind == yaml.MappingNode {
		for i := 0; i+1 < len(node.Content); i += 2 {
			if node.Content[i].Value == key {
				return node.Content[i+1]
			}
		}
	}
	return &yaml.Node{}
}

func sourceGroup(requirements *yaml.Node, groupID string) *yaml.Node {
	if requirements.Kind == yaml.SequenceNode {
		for _, group := range requirements.Content {
			if mappingScalar(group, "id") == groupID {
				return group
			}
		}
		return &yaml.Node{}
	}
	return mappingValue(requirements, groupID)
}

func nodeText(node *yaml.Node) string {
	out, err := yaml.Marshal(node)
	if err != nil {
		return "unmarshalable: " + err.Error()
	}
	return string(out)
}
