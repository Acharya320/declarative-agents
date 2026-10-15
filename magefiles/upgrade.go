// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
)

// Upgrade moves the tracked target pin for one parent repository to its
// newest release and performs the mechanical half of the upgrade: every
// derived_from entry whose cited parent item is byte-identical between its
// validated release and the new pin is re-pinned in place. What remains is
// exactly the set of requirements whose upstream changed; each carries a
// stale-external-reference audit finding, and upgrade:status prints them as
// the worklist (design.yaml external_references, GH-2558).
// DA_UPGRADE_RELEASE overrides the target tag.
func Upgrade(repo string) error {
	cache, err := externalCacheDir(repo)
	if err != nil {
		return err
	}
	if _, err := extGitOutput(cache, "rev-parse", "--git-dir"); err != nil {
		return fmt.Errorf("parent cache for %s is missing; clone it: git clone https://%s %s",
			repo, repo, cache)
	}
	if _, err := extGitOutput(cache, "fetch", "--tags", "origin"); err != nil {
		return fmt.Errorf("fetch %s tags through %s: %w", repo, cache, err)
	}
	target := os.Getenv("DA_UPGRADE_RELEASE")
	if target == "" {
		if target, err = cacheLatestReleaseTag(cache); err != nil {
			return err
		}
	}
	if target == "" {
		return fmt.Errorf("no release tag in the %s cache to upgrade to", repo)
	}
	return upgradeToRelease(auditParticipants(), repo, cache, target, os.Stdout)
}

type UPGRADE mg.Namespace

// Status prints every derived_from entry still validated against a release
// older than its repository's tracked pin: the pre-planning worklist after
// an upgrade. A report, not a gate -- it changes nothing and exits zero.
func (UPGRADE) Status() error {
	return printUpgradeStatus(auditParticipants(), os.Stdout)
}

// cacheLatestReleaseTag returns the cache's newest local v0.* tag.
func cacheLatestReleaseTag(cache string) (string, error) {
	tags, err := extGitOutput(cache, "tag", "--list", "v0.*", "--sort=-version:refname")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.SplitN(tags, "\n", 2)[0]), nil
}

// upgradeToRelease rewrites the pin and auto-re-pins unchanged entries.
func upgradeToRelease(modules []string, repo, cache, target string, out io.Writer) error {
	for _, module := range modules {
		if err := writeExternalPin(module, repo, target); err != nil {
			return err
		}
	}
	fmt.Fprintf(out, "pin %s -> %s\n", repo, target)
	entries, err := collectExternalRefEntries(modules)
	if err != nil {
		return err
	}
	repinned, stale := 0, 0
	for _, entry := range entries {
		if entry.Repo != repo || entry.Release == "" || entry.Release == target {
			continue
		}
		identical, err := citedItemUnchanged(cache, entry, target)
		if err != nil {
			return err
		}
		if !identical {
			stale++
			fmt.Fprintf(out, "stale: %s: %s %s %s validated at %s, pin %s; revalidate and re-pin\n",
				entry.File, entry.Repo, entry.SRD, entry.Item, entry.Release, target)
			continue
		}
		if err := repinEntry(entry, target); err != nil {
			return err
		}
		repinned++
		fmt.Fprintf(out, "re-pinned: %s: %s %s unchanged between %s and %s\n",
			entry.File, entry.SRD, entry.Item, entry.Release, target)
	}
	fmt.Fprintf(out, "%d entr%s re-pinned, %d stale\n", repinned, pluralY(repinned), stale)
	return nil
}

func pluralY(n int) string {
	if n == 1 {
		return "y"
	}
	return "ies"
}

// writeExternalPin sets pins[repo] in the module's external-pins file with a
// scalar-only edit, creating the file when the module carries entries for
// the repo but no pins yet. Modules without docs/ are skipped.
func writeExternalPin(module, repo, target string) error {
	if _, err := os.Stat(filepath.Join(module, "docs")); err != nil {
		return nil
	}
	pinsPath := filepath.Join(module, "docs", "external-pins.yaml")
	data, err := os.ReadFile(pinsPath)
	if os.IsNotExist(err) {
		return createExternalPins(module, pinsPath, repo, target)
	}
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	prefix := "  " + repo + ":"
	for i, line := range lines {
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		lines[i] = prefix + " " + target
		return os.WriteFile(pinsPath, []byte(strings.Join(lines, "\n")), 0o644)
	}
	for i, line := range lines {
		if strings.TrimRight(line, " ") != "pins:" {
			continue
		}
		lines = append(lines[:i+1], append([]string{prefix + " " + target}, lines[i+1:]...)...)
		return os.WriteFile(pinsPath, []byte(strings.Join(lines, "\n")), 0o644)
	}
	return fmt.Errorf("%s: no pins mapping to update", pinsPath)
}

// createExternalPins writes a fresh pins file when the module references the
// repo without one.
func createExternalPins(module, pinsPath, repo, target string) error {
	entries, err := collectExternalRefEntries([]string{module})
	if err != nil {
		return err
	}
	referenced := false
	for _, entry := range entries {
		if entry.Repo == repo {
			referenced = true
			break
		}
	}
	if !referenced {
		return nil
	}
	content := "# Tracked target pins for parent repositories" +
		" (design.yaml external_references).\npins:\n  " + repo + ": " + target + "\n"
	return os.WriteFile(pinsPath, []byte(content), 0o644)
}

// citedItemUnchanged reports whether the cited parent item's text is
// byte-identical between the entry's validated release and the target.
func citedItemUnchanged(cache string, entry externalRefEntry, target string) (bool, error) {
	before, okBefore, err := parentItemTextAt(cache, entry.Release, entry.SRD, entry.Item)
	if err != nil {
		return false, err
	}
	after, okAfter, err := parentItemTextAt(cache, target, entry.SRD, entry.Item)
	if err != nil {
		return false, err
	}
	return okBefore && okAfter && before == after, nil
}

// parentItemTextAt reads the cited id's text in the parent SRD at ref.
func parentItemTextAt(cache, ref, srdID, id string) (string, bool, error) {
	file, err := parentSRDFileAt(cache, ref, srdID)
	if err != nil || file == "" {
		return "", false, err
	}
	content, err := extGitOutput(cache, "show", ref+":"+file)
	if err != nil {
		return "", false, nil
	}
	text, ok := parentItemText([]byte(content), id)
	return text, ok, nil
}

// repinEntry rewrites the entry's release scalar in place: within the
// derived_from block that names the entry's repo, srd, and item, the
// release line alone changes; everything else in the file survives
// byte-for-byte.
func repinEntry(entry externalRefEntry, target string) error {
	data, err := os.ReadFile(entry.File)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	inEntry := false
	matched := map[string]bool{}
	releaseLine := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "- "))
		if strings.HasPrefix(trimmed, "repo:") {
			inEntry = strings.TrimSpace(strings.TrimPrefix(trimmed, "repo:")) == entry.Repo
			matched = map[string]bool{}
			releaseLine = -1
		}
		if !inEntry {
			continue
		}
		switch {
		case strings.HasPrefix(trimmed, "srd:"):
			matched["srd"] = strings.TrimSpace(strings.TrimPrefix(trimmed, "srd:")) == entry.SRD
		case strings.HasPrefix(trimmed, "item:"):
			matched["item"] = strings.TrimSpace(strings.TrimPrefix(trimmed, "item:")) == entry.Item
		case strings.HasPrefix(trimmed, "release:"):
			if strings.TrimSpace(strings.TrimPrefix(trimmed, "release:")) == entry.Release {
				releaseLine = i
			}
		}
		if matched["srd"] && matched["item"] && releaseLine >= 0 {
			indent := lines[releaseLine][:len(lines[releaseLine])-len(strings.TrimLeft(lines[releaseLine], " "))]
			lines[releaseLine] = indent + "release: " + target
			return os.WriteFile(entry.File, []byte(strings.Join(lines, "\n")), 0o644)
		}
	}
	return fmt.Errorf("%s: derived_from entry for %s %s %s at %s not found for re-pinning",
		entry.File, entry.Repo, entry.SRD, entry.Item, entry.Release)
}

// printUpgradeStatus lists every entry lagging its repository's tracked pin.
func printUpgradeStatus(modules []string, out io.Writer) error {
	pins := collectExternalPins(modules)
	entries, err := collectExternalRefEntries(modules)
	if err != nil {
		return err
	}
	open := 0
	for _, entry := range entries {
		pin := pins[entry.Repo]
		if entry.Repo == "" || pin == "" || entry.Release == "" || entry.Release == pin {
			continue
		}
		open++
		fmt.Fprintf(out, "%s: %s %s %s validated at %s, pin %s\n",
			entry.File, entry.Repo, entry.SRD, entry.Item, entry.Release, pin)
	}
	if open == 0 {
		fmt.Fprintln(out, "no entries behind their pins")
		return nil
	}
	fmt.Fprintf(out, "%d entr%s behind their pins\n", open, pluralY(open))
	return nil
}
