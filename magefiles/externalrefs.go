// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// External references (design.yaml external_references): derived_from
// entries in the spec corpora cite parent-repository requirements, resolved
// against sibling clone caches. The deep validation lives in agent-core's
// pkg/spec; the helpers here give the root tooling the three thin views it
// needs -- the entry inventory, the pins, and the caches -- without a
// cross-module import.

// extGitOutput runs one read-only git command for the external-reference
// helpers. Tests swap it.
var extGitOutput = func(dir string, args ...string) (string, error) {
	out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// externalRefEntry is one derived_from entry with its carrier file.
type externalRefEntry struct {
	File    string
	Repo    string `yaml:"repo"`
	Release string `yaml:"release"`
	SRD     string `yaml:"srd"`
	Item    string `yaml:"item"`
}

// collectExternalRefEntries walks the modules' docs trees and returns every
// derived_from entry, wherever it sits in a document.
func collectExternalRefEntries(modules []string) ([]externalRefEntry, error) {
	var entries []externalRefEntry
	for _, module := range modules {
		docsDir := filepath.Join(module, "docs")
		if _, err := os.Stat(docsDir); err != nil {
			continue
		}
		files, err := yamlDocFiles(docsDir)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			var doc any
			if err := yaml.Unmarshal(data, &doc); err != nil {
				continue
			}
			collectDerivedFrom(doc, file, &entries)
		}
	}
	return entries, nil
}

// collectDerivedFrom recurses a decoded YAML value and gathers every
// derived_from list it finds.
func collectDerivedFrom(value any, file string, entries *[]externalRefEntry) {
	switch node := value.(type) {
	case map[string]any:
		for key, child := range node {
			if key != "derived_from" {
				collectDerivedFrom(child, file, entries)
				continue
			}
			list, ok := child.([]any)
			if !ok {
				continue
			}
			for _, raw := range list {
				fields, ok := raw.(map[string]any)
				if !ok {
					continue
				}
				*entries = append(*entries, externalRefEntry{
					File:    file,
					Repo:    stringField(fields, "repo"),
					Release: stringField(fields, "release"),
					SRD:     stringField(fields, "srd"),
					Item:    stringField(fields, "item"),
				})
			}
		}
	case []any:
		for _, child := range node {
			collectDerivedFrom(child, file, entries)
		}
	}
}

func stringField(fields map[string]any, key string) string {
	s, _ := fields[key].(string)
	return s
}

// collectExternalPins merges the modules' docs/external-pins.yaml files into
// one repo -> release map. Pins agree across modules or the audit's own
// staleness findings surface the drift; last reader wins here.
func collectExternalPins(modules []string) map[string]string {
	pins := map[string]string{}
	for _, module := range modules {
		data, err := os.ReadFile(filepath.Join(module, "docs", "external-pins.yaml"))
		if err != nil {
			continue
		}
		var doc struct {
			Pins map[string]string `yaml:"pins"`
		}
		if err := yaml.Unmarshal(data, &doc); err != nil {
			continue
		}
		for repo, release := range doc.Pins {
			pins[repo] = release
		}
	}
	return pins
}

// externalCacheDir names a parent's sibling clone cache from the repository
// root's parent directory.
func externalCacheDir(repo string) (string, error) {
	toplevel, err := extGitOutput(".", "rev-parse", "--show-toplevel")
	if err != nil || toplevel == "" {
		return "", fmt.Errorf("resolve repository toplevel: %w", err)
	}
	return filepath.Join(filepath.Dir(toplevel), path.Base(repo)), nil
}

// externalRefRepos lists every parent repository the corpus references or
// pins, sorted.
func externalRefRepos(entries []externalRefEntry, pins map[string]string) []string {
	repos := map[string]bool{}
	for _, entry := range entries {
		if entry.Repo != "" {
			repos[entry.Repo] = true
		}
	}
	for repo := range pins {
		repos[repo] = true
	}
	names := make([]string, 0, len(repos))
	for repo := range repos {
		names = append(names, repo)
	}
	sort.Strings(names)
	return names
}

// reportExternalCaches tells doctor whether each referenced parent's cache
// exists, with the clone command when it does not.
func reportExternalCaches(modules []string, out io.Writer) error {
	entries, err := collectExternalRefEntries(modules)
	if err != nil {
		return err
	}
	repos := externalRefRepos(entries, collectExternalPins(modules))
	for _, repo := range repos {
		cache, err := externalCacheDir(repo)
		if err != nil {
			return err
		}
		if _, gitErr := extGitOutput(cache, "rev-parse", "--git-dir"); gitErr != nil {
			fmt.Fprintf(out, "external cache missing for %s; clone it: git clone https://%s %s\n",
				repo, repo, cache)
			continue
		}
		fmt.Fprintf(out, "external cache ok: %s at %s\n", repo, cache)
	}
	return nil
}

// reportExternalPinCurrency tells bump whether each tracked spec pin is
// still the parent's newest release. It reads the remote through ls-remote
// in the cache: a network read, no local mutation, matching bump's contract.
func reportExternalPinCurrency(modules []string, out io.Writer) {
	pins := collectExternalPins(modules)
	repos := make([]string, 0, len(pins))
	for repo := range pins {
		repos = append(repos, repo)
	}
	sort.Strings(repos)
	for _, repo := range repos {
		cache, err := externalCacheDir(repo)
		if err != nil {
			fmt.Fprintf(out, "spec pin %s: %v\n", repo, err)
			continue
		}
		latest, err := remoteLatestReleaseTag(cache)
		if err != nil {
			fmt.Fprintf(out, "spec pin %s: cannot read remote tags through %s (%v)\n", repo, cache, err)
			continue
		}
		pin := pins[repo]
		if latest == "" || latest == pin {
			fmt.Fprintf(out, "spec pin %s: %s is current\n", repo, pin)
			continue
		}
		fmt.Fprintf(out, "spec pin %s: pinned %s, newest release %s\n", repo, pin, latest)
	}
}

// remoteLatestReleaseTag reads the parent remote's newest v0.* tag through
// the cache's origin.
func remoteLatestReleaseTag(cache string) (string, error) {
	out, err := extGitOutput(cache, "ls-remote", "--tags", "origin", "v0.*")
	if err != nil {
		return "", err
	}
	latest := ""
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		tag := strings.TrimSuffix(strings.TrimPrefix(fields[1], "refs/tags/"), "^{}")
		if releaseTagAfter(tag, latest) {
			latest = tag
		}
	}
	return latest, nil
}

// releaseTagAfter reports whether a sorts after b as a v0.YYYYMMDD.N tag;
// an empty b loses to any tag.
func releaseTagAfter(a, b string) bool {
	if b == "" {
		return a != ""
	}
	return releaseTagKey(a) > releaseTagKey(b)
}

func releaseTagKey(tag string) string {
	parts := strings.Split(strings.TrimPrefix(tag, "v"), ".")
	if len(parts) != 3 {
		return tag
	}
	return fmt.Sprintf("%08s.%08s.%08s", parts[0], parts[1], parts[2])
}

// externalReferenceReleaseGate refuses a release while any derived_from
// entry is unpinned or cites an id the parent retired. Stale entries stay
// warnings a person judges (GH-2557).
func externalReferenceReleaseGate(modules []string) error {
	entries, err := collectExternalRefEntries(modules)
	if err != nil {
		return err
	}
	var refusals []string
	for _, entry := range entries {
		if entry.Repo == "" || entry.SRD == "" || entry.Item == "" {
			continue // the audit's broken-external-reference owns malformed entries
		}
		if entry.Release == "" {
			refusals = append(refusals,
				fmt.Sprintf("%s: %s %s %s is unpinned; pin a release before tagging",
					entry.File, entry.Repo, entry.SRD, entry.Item))
			continue
		}
		if reason := retiredAtParentLatest(entry); reason != "" {
			refusals = append(refusals, fmt.Sprintf("%s: %s", entry.File, reason))
		}
	}
	if len(refusals) == 0 {
		return nil
	}
	return fmt.Errorf("external-reference release gate:\n  %s", strings.Join(refusals, "\n  "))
}

// retiredAtParentLatest reports why an entry blocks the release when its
// cited id is retired at the parent cache's newest release tag, or "" when
// it does not block.
func retiredAtParentLatest(entry externalRefEntry) string {
	cache, err := externalCacheDir(entry.Repo)
	if err != nil {
		return ""
	}
	tags, err := extGitOutput(cache, "tag", "--list", "v0.*", "--sort=-version:refname")
	if err != nil || tags == "" {
		return "" // missing cache or no releases: the audit's checks own this
	}
	latest := strings.TrimSpace(strings.SplitN(tags, "\n", 2)[0])
	file, err := parentSRDFileAt(cache, latest, entry.SRD)
	if err != nil || file == "" {
		return ""
	}
	content, err := extGitOutput(cache, "show", latest+":"+file)
	if err != nil {
		return ""
	}
	if !parentItemRetired([]byte(content), entry.Item) {
		return ""
	}
	return fmt.Sprintf("%s %s %s is retired at %s; revalidate before tagging",
		entry.Repo, entry.SRD, entry.Item, latest)
}

// parentSRDFileAt finds the SRD's file anywhere in the parent tree at ref.
func parentSRDFileAt(cache, ref, srdID string) (string, error) {
	listing, err := extGitOutput(cache, "ls-tree", "-r", "--name-only", ref)
	if err != nil {
		return "", err
	}
	want := "docs/specs/software-requirements/" + srdID + ".yaml"
	for _, file := range strings.Split(listing, "\n") {
		file = strings.TrimSpace(file)
		if file == want || strings.HasSuffix(file, "/"+want) {
			return file, nil
		}
	}
	return "", nil
}

// parentItemRetired reports whether the id carries status retired in the
// SRD document bytes: the document itself, a requirement item, or an
// acceptance criterion.
func parentItemRetired(data []byte, id string) bool {
	var doc struct {
		Status       string `yaml:"status"`
		Requirements map[string]struct {
			Items []map[string]any `yaml:"items"`
		} `yaml:"requirements"`
		AcceptanceCriteria []map[string]any `yaml:"acceptance_criteria"`
	}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return false
	}
	if doc.Status == "retired" {
		return true
	}
	for _, group := range doc.Requirements {
		for _, item := range group.Items {
			if _, hasID := item[id]; hasID {
				return stringField(item, "status") == "retired"
			}
		}
	}
	for _, ac := range doc.AcceptanceCriteria {
		if stringField(ac, "id") == id {
			return stringField(ac, "status") == "retired"
		}
	}
	return false
}
