// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package spec

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// externalPinsFile is the tracked target pin per parent repository: the
// staleness baseline (design.yaml external_references). Module-relative.
const externalPinsFile = "docs/external-pins.yaml"

// externalCheckoutsFile optionally overrides where a parent's cache clone
// lives, for layouts where the sibling convention does not hold.
const externalCheckoutsFile = "docs/external-checkouts.yaml"

// externalRefCitation is one derived_from entry with its carrier named.
type externalRefCitation struct {
	citer string // "srdNNN", "srdNNN R1.1", or "srdNNN AC1"
	file  string // module-relative file of the carrier
	ref   ExternalReference
}

// checkExternalReferences validates every derived_from entry against the
// parent's sibling clone cache: the entry resolves at its validated release,
// the cited id exists and is not retired at the cache's latest release, and
// the validated release keeps up with the tracked target pin. Resolution is
// local git against the cache, never the network (srd014 R5.3).
func checkExternalReferences(corpus *Corpus) []Finding {
	citations := collectExternalRefCitations(corpus)
	if len(citations) == 0 {
		return nil
	}
	resolver := newExternalRefResolver(corpus.RootDir)
	var findings []Finding
	for _, citation := range citations {
		findings = append(findings, resolver.check(citation)...)
	}
	return findings
}

// collectExternalRefCitations gathers derived_from entries from SRD
// documents, requirement items, and acceptance criteria, in corpus order.
func collectExternalRefCitations(corpus *Corpus) []externalRefCitation {
	var citations []externalRefCitation
	add := func(citer, file string, refs []ExternalReference) {
		for _, ref := range refs {
			citations = append(citations, externalRefCitation{citer: citer, file: file, ref: ref})
		}
	}
	for _, srdID := range sortedSRDIDs(corpus) {
		srd := corpus.SRDs[srdID]
		file := path.Join(SRDSubdir, srdID+".yaml")
		add(srdID, file, srd.DerivedFrom)
		for _, item := range srd.AllItems() {
			add(srdID+" "+item.ID, file, item.DerivedFrom)
		}
		for _, ac := range srd.AcceptanceCriteria {
			add(srdID+" "+ac.ID, file, ac.DerivedFrom)
		}
	}
	return citations
}

// externalRefResolver resolves entries against sibling clone caches and
// memoizes per-cache state: the parent tree listing per ref and the parsed
// parent SRDs.
type externalRefResolver struct {
	moduleDir string
	pins      map[string]string // repo -> tracked target release
	checkouts map[string]string // repo -> cache path override
	parentDir string            // directory holding sibling caches
	srdCache  map[string]*SRD   // "<cache>|<ref>|<srdID>" -> parsed SRD, nil when absent
	latest    map[string]string // cache path -> latest release tag
}

func newExternalRefResolver(moduleDir string) *externalRefResolver {
	if moduleDir == "" {
		moduleDir = "."
	}
	resolver := &externalRefResolver{
		moduleDir: moduleDir,
		pins:      loadRepoMap(filepath.Join(moduleDir, externalPinsFile), "pins"),
		checkouts: loadRepoMap(filepath.Join(moduleDir, externalCheckoutsFile), "checkouts"),
		srdCache:  map[string]*SRD{},
		latest:    map[string]string{},
	}
	if toplevel, err := gitOutput(moduleDir, "rev-parse", "--show-toplevel"); err == nil && toplevel != "" {
		resolver.parentDir = filepath.Dir(toplevel)
	}
	return resolver
}

// loadRepoMap reads a {<key>: {repo: value}} YAML file; a missing file is an
// empty map.
func loadRepoMap(file, key string) map[string]string {
	data, err := os.ReadFile(file)
	if err != nil {
		return map[string]string{}
	}
	var doc map[string]map[string]string
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return map[string]string{}
	}
	if doc[key] == nil {
		return map[string]string{}
	}
	return doc[key]
}

// cacheDir names the parent's clone cache: the override when declared,
// otherwise the sibling of the current checkout named after the repository.
func (r *externalRefResolver) cacheDir(repo string) string {
	if override := r.checkouts[repo]; override != "" {
		if filepath.IsAbs(override) {
			return override
		}
		return filepath.Join(r.moduleDir, override)
	}
	if r.parentDir == "" {
		return ""
	}
	return filepath.Join(r.parentDir, path.Base(repo))
}

func externalFinding(check, level, citer, file, message string) Finding {
	return Finding{
		Check:   check,
		Level:   level,
		Message: citer + ": " + message,
		File:    file,
	}
}

// check validates one citation and returns its findings.
func (r *externalRefResolver) check(c externalRefCitation) []Finding {
	if c.ref.Repo == "" || c.ref.SRD == "" || c.ref.Item == "" {
		return []Finding{externalFinding("broken-external-reference", "error", c.citer, c.file,
			"derived_from entry needs repo, srd, and item")}
	}
	cache := r.cacheDir(c.ref.Repo)
	if cache == "" || !r.cacheUsable(cache) {
		return []Finding{externalFinding("unverifiable-external-reference", "error", c.citer, c.file,
			fmt.Sprintf("parent cache for %s is missing; clone it: git clone https://%s %s",
				c.ref.Repo, c.ref.Repo, cache))}
	}
	ref := c.ref.Release
	if ref == "" {
		ref = "HEAD"
	} else if _, err := gitOutput(cache, "rev-parse", "--verify", "refs/tags/"+ref); err != nil {
		return []Finding{externalFinding("unverifiable-external-reference", "error", c.citer, c.file,
			fmt.Sprintf("release %s is not in the parent cache; fetch it: git -C %s fetch --tags",
				ref, cache))}
	}
	srd := r.parentSRD(cache, ref, c.ref.SRD)
	if srd == nil {
		return []Finding{externalFinding("broken-external-reference", "error", c.citer, c.file,
			fmt.Sprintf("%s has no SRD %s at %s", c.ref.Repo, c.ref.SRD, ref))}
	}
	text, ok := citedItemText(srd, c.ref.Item)
	if !ok {
		return []Finding{externalFinding("broken-external-reference", "error", c.citer, c.file,
			fmt.Sprintf("%s %s has no id %s at %s", c.ref.Repo, c.ref.SRD, c.ref.Item, ref))}
	}
	var findings []Finding
	findings = append(findings, r.retiredFindings(c, cache)...)
	findings = append(findings, r.currencyFindings(c, cache, text)...)
	return findings
}

// cacheUsable reports whether the cache answers read-only git calls.
func (r *externalRefResolver) cacheUsable(cache string) bool {
	_, err := gitOutput(cache, "rev-parse", "--git-dir")
	return err == nil
}

// retiredFindings warns when the cited id is retired at the cache's latest
// release tag.
func (r *externalRefResolver) retiredFindings(c externalRefCitation, cache string) []Finding {
	latest := r.latestRelease(cache)
	if latest == "" {
		return nil
	}
	srd := r.parentSRD(cache, latest, c.ref.SRD)
	if srd == nil {
		return nil
	}
	retired, replacedBy := citedItemRetired(srd, c.ref.Item)
	if !retired {
		return nil
	}
	replacement := "nothing"
	if len(replacedBy) > 0 {
		replacement = strings.Join(replacedBy, ", ")
	}
	return []Finding{externalFinding("retired-external-reference", "warning", c.citer, c.file,
		fmt.Sprintf("%s %s %s is retired at %s (replaced by %s)",
			c.ref.Repo, c.ref.SRD, c.ref.Item, latest, replacement))}
}

// currencyFindings compares the entry's validated release with the tracked
// target pin: unpinned entries warn, and a lagging entry warns only when the
// cited item's text changed between the two releases -- a new parent release
// alone is not staleness.
func (r *externalRefResolver) currencyFindings(c externalRefCitation, cache, textAtValidated string) []Finding {
	if c.ref.Release == "" {
		return []Finding{externalFinding("unpinned-external-reference", "warning", c.citer, c.file,
			fmt.Sprintf("derived_from %s %s %s has no release; pin it before tagging (design.yaml external_references)",
				c.ref.Repo, c.ref.SRD, c.ref.Item))}
	}
	pin := r.pins[c.ref.Repo]
	if pin == "" || pin == c.ref.Release {
		return nil
	}
	if compareReleaseTags(c.ref.Release, pin) > 0 {
		return []Finding{externalFinding("stale-external-reference", "warning", c.citer, c.file,
			fmt.Sprintf("validated release %s is ahead of the tracked pin %s for %s; raise the pin",
				c.ref.Release, pin, c.ref.Repo))}
	}
	srdAtPin := r.parentSRD(cache, pin, c.ref.SRD)
	if srdAtPin != nil {
		if textAtPin, ok := citedItemText(srdAtPin, c.ref.Item); ok && textAtPin == textAtValidated {
			return nil
		}
	}
	return []Finding{externalFinding("stale-external-reference", "warning", c.citer, c.file,
		fmt.Sprintf("%s %s %s changed between validated release %s and the tracked pin %s; revalidate and re-pin in place",
			c.ref.Repo, c.ref.SRD, c.ref.Item, c.ref.Release, pin))}
}

// latestRelease memoizes the cache's newest release tag.
func (r *externalRefResolver) latestRelease(cache string) string {
	if latest, seen := r.latest[cache]; seen {
		return latest
	}
	out, err := gitOutput(cache, "tag", "--list", releaseTagPattern, "--sort=-version:refname")
	latest := ""
	if err == nil {
		latest = strings.TrimSpace(strings.SplitN(out, "\n", 2)[0])
	}
	r.latest[cache] = latest
	return latest
}

// parentSRD reads and memoizes one parent SRD at one ref, locating the file
// anywhere in the parent tree so module layouts need no declaration.
func (r *externalRefResolver) parentSRD(cache, ref, srdID string) *SRD {
	key := cache + "|" + ref + "|" + srdID
	if srd, seen := r.srdCache[key]; seen {
		return srd
	}
	r.srdCache[key] = nil
	listing, err := gitOutput(cache, "ls-tree", "-r", "--name-only", ref)
	if err != nil {
		return nil
	}
	suffix := "/" + SRDSubdir + "/" + srdID + ".yaml"
	for _, file := range strings.Split(listing, "\n") {
		file = strings.TrimSpace(file)
		if file != SRDSubdir+"/"+srdID+".yaml" && !strings.HasSuffix(file, suffix) {
			continue
		}
		content, err := gitOutput(cache, "show", ref+":"+file)
		if err != nil {
			return nil
		}
		srd, err := parseSRDBytes([]byte(content), ref+":"+file)
		if err != nil {
			return nil
		}
		r.srdCache[key] = &srd
		return &srd
	}
	return nil
}

// citedItemText returns the text behind a cited id: a requirement item's
// text, an acceptance criterion's criterion, or a group's title joined with
// its items.
func citedItemText(srd *SRD, id string) (string, bool) {
	for _, item := range srd.AllItems() {
		if item.ID == id {
			return item.Text, true
		}
	}
	for _, ac := range srd.AcceptanceCriteria {
		if ac.ID == id {
			return ac.Criterion, true
		}
	}
	group, ok := srd.Requirements[id]
	if !ok {
		return "", false
	}
	parts := []string{group.Title}
	for _, item := range group.Items {
		parts = append(parts, item.ID+": "+item.Text)
	}
	return strings.Join(parts, "\n"), true
}

// citedItemRetired reports whether the cited id is retired in this SRD
// state, and what replaces it. A retired SRD retires everything it holds.
func citedItemRetired(srd *SRD, id string) (bool, []string) {
	if srd.Status == statusRetired {
		return true, srd.ReplacedBy
	}
	for _, item := range srd.AllItems() {
		if item.ID == id {
			return item.Status == statusRetired, item.ReplacedBy
		}
	}
	for _, ac := range srd.AcceptanceCriteria {
		if ac.ID == id {
			return ac.Status == statusRetired, ac.ReplacedBy
		}
	}
	return false, nil
}

// compareReleaseTags orders two v0.YYYYMMDD.N tags: negative, zero, or
// positive as a sorts before, equal to, or after b. Unparseable tags compare
// as strings.
func compareReleaseTags(a, b string) int {
	return strings.Compare(releaseTagSortKey(a), releaseTagSortKey(b))
}

func releaseTagSortKey(tag string) string {
	parts := strings.Split(strings.TrimPrefix(tag, "v"), ".")
	if len(parts) != 3 {
		return tag
	}
	return fmt.Sprintf("%08s.%08s.%08s", parts[0], parts[1], parts[2])
}
