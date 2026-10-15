// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package spec

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const parentSRDv1 = `id: srd013-contract
title: Contract
requirements:
  R2:
    title: Group
    items:
      - R2.4: Expose the applier contract.
acceptance_criteria:
  - id: AC1
    criterion: Works.
    traces: [R2.4]
`

const parentSRDv2Changed = `id: srd013-contract
title: Contract
requirements:
  R2:
    title: Group
    items:
      - R2.4: Expose the applier contract, now with leases.
acceptance_criteria:
  - id: AC1
    criterion: Works.
    traces: [R2.4]
`

const parentSRDv2Retired = `id: srd013-contract
title: Contract
requirements:
  R2:
    title: Group
    items:
      - R2.4: Expose the applier contract.
        status: retired
        retired_in: v0.20270102.0
        replaced_by: [R2.5]
      - R2.5: Expose the lease-holding applier contract.
acceptance_criteria:
  - id: AC1
    criterion: Works.
    traces: [R2.5]
`

// fakeParentCache serves one parent repository's git surface: tags, each
// tag's tree, and the child checkout's toplevel discovery. Caches the test
// does not declare refuse every call, modelling a missing clone.
type fakeParentCache struct {
	toplevel string            // the child checkout's toplevel
	cacheDir string            // the one usable cache directory
	tags     []string          // release tags, any order
	files    map[string]string // "tag|path" -> content; "HEAD|path" works too
}

func (f fakeParentCache) output(dir string, args ...string) (string, error) {
	if args[0] == "rev-parse" && args[1] == "--show-toplevel" {
		return f.toplevel, nil
	}
	if dir != f.cacheDir {
		return "", errors.New("not a repository: " + dir)
	}
	switch args[0] {
	case "rev-parse":
		if args[1] == "--git-dir" {
			return ".git", nil
		}
		tag := strings.TrimPrefix(args[len(args)-1], "refs/tags/")
		for _, t := range f.tags {
			if t == tag {
				return tag, nil
			}
		}
		return "", errors.New("unknown ref " + args[1])
	case "tag":
		sorted := append([]string(nil), f.tags...)
		sort.Slice(sorted, func(i, j int) bool {
			return compareReleaseTags(sorted[i], sorted[j]) > 0
		})
		return strings.Join(sorted, "\n"), nil
	case "ls-tree":
		ref := args[3]
		var names []string
		for key := range f.files {
			if strings.HasPrefix(key, ref+"|") {
				names = append(names, strings.TrimPrefix(key, ref+"|"))
			}
		}
		sort.Strings(names)
		return strings.Join(names, "\n"), nil
	case "show":
		ref, path, _ := strings.Cut(args[1], ":")
		content, ok := f.files[ref+"|"+path]
		if !ok {
			return "", errors.New("path not at ref: " + args[1])
		}
		return content, nil
	}
	return "", errors.New("unexpected git call: " + strings.Join(args, " "))
}

// externalRefCorpus builds a child corpus whose srd090 R1.1 derives from the
// parent's srd013 R2.4 at the given release ("" = unpinned), rooted in a
// module directory the test controls.
func externalRefCorpus(t *testing.T, moduleDir, release string) *Corpus {
	t.Helper()
	srdYAML := `id: srd090-child
title: Child
requirements:
  R1:
    title: Group
    items:
      - R1.1: Implement the applier contract here.
        derived_from:
          - repo: github.com/example/parent
            srd: srd013-contract
            item: R2.4
`
	srd, err := parseSRDBytes([]byte(srdYAML), "head")
	require.NoError(t, err)
	if release != "" {
		item := srd.Requirements["R1"].Items[0]
		item.DerivedFrom[0].Release = release
		srd.Requirements["R1"] = RequirementGroup{
			Title: srd.Requirements["R1"].Title,
			Items: []RequirementItem{item},
		}
	}
	return &Corpus{
		RootDir:    moduleDir,
		SRDs:       map[string]SRD{srd.ID: srd},
		UseCases:   map[string]UseCase{},
		TestSuites: map[string]TestSuite{},
	}
}

// externalRefFixture lays out a child module under a root, a cache beside
// it, and the swapped git reader. Returns the module dir and the cache dir.
func externalRefFixture(t *testing.T, tags []string, files map[string]string) (string, string) {
	t.Helper()
	root := t.TempDir()
	moduleDir := filepath.Join(root, "child", "core")
	require.NoError(t, os.MkdirAll(filepath.Join(moduleDir, "docs"), 0o755))
	cacheDir := filepath.Join(root, "parent")
	swapGitOutput(t, fakeParentCache{
		toplevel: filepath.Join(root, "child"),
		cacheDir: cacheDir,
		tags:     tags,
		files:    files,
	}.output)
	return moduleDir, cacheDir
}

func writePins(t *testing.T, moduleDir, release string) {
	t.Helper()
	content := "pins:\n  github.com/example/parent: " + release + "\n"
	require.NoError(t, os.WriteFile(
		filepath.Join(moduleDir, "docs", "external-pins.yaml"), []byte(content), 0o644))
}

const parentSRDPath = "agent-core/docs/specs/software-requirements/srd013-contract.yaml"

func TestExternalReferenceCleanPinnedPasses(t *testing.T) {
	moduleDir, _ := externalRefFixture(t, []string{"v0.20261001.0"},
		map[string]string{"v0.20261001.0|" + parentSRDPath: parentSRDv1})
	writePins(t, moduleDir, "v0.20261001.0")
	corpus := externalRefCorpus(t, moduleDir, "v0.20261001.0")
	require.Empty(t, checkExternalReferences(corpus))
}

func TestExternalReferenceMissingCacheIsUnverifiable(t *testing.T) {
	root := t.TempDir()
	moduleDir := filepath.Join(root, "child", "core")
	require.NoError(t, os.MkdirAll(filepath.Join(moduleDir, "docs"), 0o755))
	swapGitOutput(t, fakeParentCache{
		toplevel: filepath.Join(root, "child"),
		cacheDir: filepath.Join(root, "elsewhere"),
	}.output)
	corpus := externalRefCorpus(t, moduleDir, "v0.20261001.0")
	findings := checkExternalReferences(corpus)
	require.Len(t, findings, 1)
	require.Equal(t, "unverifiable-external-reference", findings[0].Check)
	require.Equal(t, "error", findings[0].Level)
	require.Contains(t, findings[0].Message, "git clone https://github.com/example/parent")
}

func TestExternalReferenceMissingTagIsUnverifiable(t *testing.T) {
	moduleDir, cacheDir := externalRefFixture(t, []string{"v0.20261001.0"},
		map[string]string{"v0.20261001.0|" + parentSRDPath: parentSRDv1})
	corpus := externalRefCorpus(t, moduleDir, "v0.20270101.9")
	findings := checkExternalReferences(corpus)
	require.Len(t, findings, 1)
	require.Equal(t, "unverifiable-external-reference", findings[0].Check)
	require.Contains(t, findings[0].Message, "git -C "+cacheDir+" fetch --tags")
}

func TestExternalReferenceBrokenSRDAndItem(t *testing.T) {
	moduleDir, _ := externalRefFixture(t, []string{"v0.20261001.0"},
		map[string]string{"v0.20261001.0|" + parentSRDPath: parentSRDv1})
	corpus := externalRefCorpus(t, moduleDir, "v0.20261001.0")
	srd := corpus.SRDs["srd090-child"]
	item := srd.Requirements["R1"].Items[0]
	item.DerivedFrom[0].Item = "R9.9"
	srd.Requirements["R1"] = RequirementGroup{Title: "Group", Items: []RequirementItem{item}}
	corpus.SRDs["srd090-child"] = srd
	findings := checkExternalReferences(corpus)
	require.Len(t, findings, 1)
	require.Equal(t, "broken-external-reference", findings[0].Check)
	require.Contains(t, findings[0].Message, "has no id R9.9")

	item.DerivedFrom[0].Item = "R2.4"
	item.DerivedFrom[0].SRD = "srd099-ghost"
	srd.Requirements["R1"] = RequirementGroup{Title: "Group", Items: []RequirementItem{item}}
	corpus.SRDs["srd090-child"] = srd
	findings = checkExternalReferences(corpus)
	require.Len(t, findings, 1)
	require.Contains(t, findings[0].Message, "has no SRD srd099-ghost")
}

func TestExternalReferenceRetiredCitedItemWarns(t *testing.T) {
	moduleDir, _ := externalRefFixture(t,
		[]string{"v0.20261001.0", "v0.20270102.0"},
		map[string]string{
			"v0.20261001.0|" + parentSRDPath: parentSRDv1,
			"v0.20270102.0|" + parentSRDPath: parentSRDv2Retired,
		})
	writePins(t, moduleDir, "v0.20261001.0")
	corpus := externalRefCorpus(t, moduleDir, "v0.20261001.0")
	findings := checkExternalReferences(corpus)
	require.Len(t, findings, 1)
	require.Equal(t, "retired-external-reference", findings[0].Check)
	require.Equal(t, "warning", findings[0].Level)
	require.Contains(t, findings[0].Message, "retired at v0.20270102.0")
	require.Contains(t, findings[0].Message, "replaced by R2.5")
}

func TestExternalReferenceStaleWhenCitedItemChanged(t *testing.T) {
	moduleDir, _ := externalRefFixture(t,
		[]string{"v0.20261001.0", "v0.20270102.0"},
		map[string]string{
			"v0.20261001.0|" + parentSRDPath: parentSRDv1,
			"v0.20270102.0|" + parentSRDPath: parentSRDv2Changed,
		})
	writePins(t, moduleDir, "v0.20270102.0")
	corpus := externalRefCorpus(t, moduleDir, "v0.20261001.0")
	findings := checkExternalReferences(corpus)
	require.Len(t, findings, 1)
	require.Equal(t, "stale-external-reference", findings[0].Check)
	require.Equal(t, "warning", findings[0].Level)
	require.Contains(t, findings[0].Message, "changed between validated release v0.20261001.0 and the tracked pin v0.20270102.0")
}

func TestExternalReferenceSilentWhenCitedItemUnchanged(t *testing.T) {
	moduleDir, _ := externalRefFixture(t,
		[]string{"v0.20261001.0", "v0.20270102.0"},
		map[string]string{
			"v0.20261001.0|" + parentSRDPath: parentSRDv1,
			"v0.20270102.0|" + parentSRDPath: parentSRDv1,
		})
	writePins(t, moduleDir, "v0.20270102.0")
	corpus := externalRefCorpus(t, moduleDir, "v0.20261001.0")
	require.Empty(t, checkExternalReferences(corpus))
}

func TestExternalReferenceUnpinnedWarns(t *testing.T) {
	moduleDir, _ := externalRefFixture(t, []string{"v0.20261001.0"},
		map[string]string{
			"v0.20261001.0|" + parentSRDPath: parentSRDv1,
			"HEAD|" + parentSRDPath:          parentSRDv1,
		})
	corpus := externalRefCorpus(t, moduleDir, "")
	findings := checkExternalReferences(corpus)
	require.Len(t, findings, 1)
	require.Equal(t, "unpinned-external-reference", findings[0].Check)
	require.Equal(t, "warning", findings[0].Level)
}

func TestExternalReferenceCheckoutOverride(t *testing.T) {
	root := t.TempDir()
	moduleDir := filepath.Join(root, "child", "core")
	require.NoError(t, os.MkdirAll(filepath.Join(moduleDir, "docs"), 0o755))
	cacheDir := filepath.Join(root, "vendored", "parent-cache")
	require.NoError(t, os.WriteFile(
		filepath.Join(moduleDir, "docs", "external-checkouts.yaml"),
		[]byte("checkouts:\n  github.com/example/parent: "+cacheDir+"\n"), 0o644))
	swapGitOutput(t, fakeParentCache{
		toplevel: filepath.Join(root, "child"),
		cacheDir: cacheDir,
		tags:     []string{"v0.20261001.0"},
		files:    map[string]string{"v0.20261001.0|" + parentSRDPath: parentSRDv1},
	}.output)
	corpus := externalRefCorpus(t, moduleDir, "v0.20261001.0")
	require.Empty(t, checkExternalReferences(corpus))
}

func TestExternalReferenceAheadOfPinWarns(t *testing.T) {
	moduleDir, _ := externalRefFixture(t,
		[]string{"v0.20261001.0", "v0.20270102.0"},
		map[string]string{
			"v0.20261001.0|" + parentSRDPath: parentSRDv1,
			"v0.20270102.0|" + parentSRDPath: parentSRDv1,
		})
	writePins(t, moduleDir, "v0.20261001.0")
	corpus := externalRefCorpus(t, moduleDir, "v0.20270102.0")
	findings := checkExternalReferences(corpus)
	require.Len(t, findings, 1)
	require.Equal(t, "stale-external-reference", findings[0].Check)
	require.Contains(t, findings[0].Message, "ahead of the tracked pin")
}
