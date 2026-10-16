// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package spec

import (
	"fmt"
	"os/exec"
	"path"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// statusRetired is the retirement marker design.yaml D6 defines. An entry
// that no longer applies keeps its id with this status, the release that
// retired it, and the ids that replace it.
const statusRetired = "retired"

// gitOutput runs one read-only git command in dir and returns its trimmed
// stdout. The checks below stay deterministic, local, and offline (srd014
// R5.3): every call reads the checkout, none writes. Tests swap it to serve
// fixture trees without a repository.
var gitOutput = func(dir string, args ...string) (string, error) {
	out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// releaseTagPattern matches the repository-wide release tags mage tag
// creates (v0.YYYYMMDD.N).
const releaseTagPattern = "v0.*"

// taggedCorpusIDs is the identifier inventory of the module's corpus at one
// release tag: every id design.yaml D6 freezes, with the module-relative
// file that declared it.
type taggedCorpusIDs struct {
	srds     map[string]string            // srd id -> file
	items    map[string]map[string]string // srd id -> requirement item id -> file
	acs      map[string]map[string]string // srd id -> acceptance criterion id -> file
	useCases map[string]string            // use case id -> file
	suites   map[string]string            // test suite id -> file
}

func missingIDFinding(id, tag, file string) Finding {
	return Finding{
		Check: "retired-id-missing",
		Level: "error",
		Message: fmt.Sprintf(
			"id %s present at release tag %s is absent at HEAD without status: retired (design.yaml D6)",
			id, tag),
		File: file,
	}
}

func bareRetirementFinding(id, file string) Finding {
	return Finding{
		Check: "retired-id-missing",
		Level: "error",
		Message: fmt.Sprintf(
			"id %s is retired without retired_in naming the release that retired it (design.yaml D6)",
			id),
		File: file,
	}
}

// checkRetiredIDs enforces design.yaml D6: an id contained in a release tag
// is never renumbered or reused. It reads the module's corpus at the latest
// release tag reachable from HEAD^ through local git and reports every id
// that is absent at HEAD without a retirement entry, and every retired entry
// that does not name the release that retired it. A module without a
// repository, a release tag, or tagged specs is fresh and yields nothing.
func checkRetiredIDs(corpus *Corpus) []Finding {
	dir := corpus.RootDir
	if dir == "" {
		dir = "."
	}
	tag, err := gitOutput(dir, "describe", "--tags", "--match", releaseTagPattern, "--abbrev=0", "HEAD^")
	if err != nil || tag == "" {
		return nil
	}
	tagged, err := loadTaggedCorpusIDs(dir, tag)
	if err != nil || tagged == nil {
		return nil
	}
	var findings []Finding
	findings = append(findings, retiredSRDFindings(corpus, tagged, tag)...)
	findings = append(findings, bareSRDRetirements(corpus)...)
	findings = append(findings, retiredDocFindings(corpus, tagged, tag)...)
	return findings
}

// retiredSRDFindings compares each SRD the tag contains with its HEAD state.
func retiredSRDFindings(corpus *Corpus, tagged *taggedCorpusIDs, tag string) []Finding {
	var findings []Finding
	for _, srdID := range sortedFileKeys(tagged.srds) {
		tagFile := tagged.srds[srdID]
		srd, ok := corpus.SRDs[srdID]
		if !ok {
			findings = append(findings, missingIDFinding(srdID, tag, tagFile))
			continue
		}
		if srd.Status == statusRetired {
			// A retired SRD keeps its id; its items need no individual entries.
			continue
		}
		findings = append(findings, retiredSRDEntryFindings(&srd, srdID, corpus.documentFile("srd", SRDSubdir, srdID), tagged, tag, tagFile)...)
	}
	return findings
}

// retiredSRDEntryFindings checks one live SRD's requirement items and
// acceptance criteria against the ids the tag froze.
func retiredSRDEntryFindings(srd *SRD, srdID, headFile string, tagged *taggedCorpusIDs, tag, tagFile string) []Finding {
	var findings []Finding
	headItems := make(map[string]RequirementItem)
	for _, item := range srd.AllItems() {
		headItems[item.ID] = item
	}
	for _, itemID := range sortedFileKeys(tagged.items[srdID]) {
		item, present := headItems[itemID]
		if !present {
			findings = append(findings, missingIDFinding(srdID+" "+itemID, tag, tagFile))
			continue
		}
		if item.Status == statusRetired && item.RetiredIn == "" {
			findings = append(findings, bareRetirementFinding(srdID+" "+itemID, headFile))
		}
	}
	headACs := make(map[string]AcceptanceCriterion)
	for _, ac := range srd.AcceptanceCriteria {
		headACs[ac.ID] = ac
	}
	for _, acID := range sortedFileKeys(tagged.acs[srdID]) {
		ac, present := headACs[acID]
		if !present {
			findings = append(findings, missingIDFinding(srdID+" "+acID, tag, tagFile))
			continue
		}
		if ac.Status == statusRetired && ac.RetiredIn == "" {
			findings = append(findings, bareRetirementFinding(srdID+" "+acID, headFile))
		}
	}
	return findings
}

// bareSRDRetirements reports SRDs retired at HEAD without retired_in, which
// needs no tag baseline: the entry itself is incomplete.
func bareSRDRetirements(corpus *Corpus) []Finding {
	var findings []Finding
	for _, srdID := range sortedSRDIDs(corpus) {
		srd := corpus.SRDs[srdID]
		if srd.Status == statusRetired && srd.RetiredIn == "" {
			findings = append(findings, bareRetirementFinding(srdID, corpus.documentFile("srd", SRDSubdir, srdID)))
		}
	}
	return findings
}

// retiredDocFindings compares the use-case and test-suite ids the tag
// contains with their HEAD state.
func retiredDocFindings(corpus *Corpus, tagged *taggedCorpusIDs, tag string) []Finding {
	var findings []Finding
	for _, ucID := range sortedFileKeys(tagged.useCases) {
		uc, ok := corpus.UseCases[ucID]
		if !ok {
			findings = append(findings, missingIDFinding(ucID, tag, tagged.useCases[ucID]))
			continue
		}
		if uc.Status == statusRetired && uc.RetiredIn == "" {
			findings = append(findings, bareRetirementFinding(ucID, corpus.documentFile("use_case", UCSubdir, ucID)))
		}
	}
	for _, tsID := range sortedFileKeys(tagged.suites) {
		ts, ok := corpus.TestSuites[tsID]
		if !ok {
			findings = append(findings, missingIDFinding(tsID, tag, tagged.suites[tsID]))
			continue
		}
		if ts.Status == statusRetired && ts.RetiredIn == "" {
			findings = append(findings, bareRetirementFinding(tsID, corpus.documentFile("test_suite", TSSubdir, tsID)))
		}
	}
	return findings
}

// loadTaggedCorpusIDs reads the module's spec files at tag through local git
// and collects the frozen ids. Paths are module-relative: git -C resolves
// the ls-tree listing and the ./-prefixed show argument against the module
// directory inside the repository. A file that no longer parses is skipped;
// retired-id-missing reports renumbered ids, not historical format drift.
func loadTaggedCorpusIDs(dir, tag string) (*taggedCorpusIDs, error) {
	listing, err := gitOutput(dir, "ls-tree", "-r", "--name-only", tag, "--", "docs/specs")
	if err != nil || listing == "" {
		return nil, err
	}
	tagged := &taggedCorpusIDs{
		srds:     map[string]string{},
		items:    map[string]map[string]string{},
		acs:      map[string]map[string]string{},
		useCases: map[string]string{},
		suites:   map[string]string{},
	}
	for _, file := range strings.Split(listing, "\n") {
		file = strings.TrimSpace(file)
		if file == "" || !isCorpusDocument(file) {
			continue
		}
		base := path.Base(file)
		switch {
		case strings.HasPrefix(file, SRDSubdir+"/") && strings.HasPrefix(base, "srd"):
			tagged.addTaggedSRD(dir, tag, file)
		case strings.HasPrefix(file, UCSubdir+"/") && strings.HasPrefix(base, "rel"):
			if id := taggedDocumentID(dir, tag, file); id != "" {
				tagged.useCases[id] = file
			}
		case strings.HasPrefix(file, TSSubdir+"/") && strings.HasPrefix(base, "test-"):
			if id := taggedDocumentID(dir, tag, file); id != "" {
				tagged.suites[id] = file
			}
		}
	}
	return tagged, nil
}

// addTaggedSRD reads one SRD at tag and records its document, item, and
// acceptance-criterion ids.
func (t *taggedCorpusIDs) addTaggedSRD(dir, tag, file string) {
	content, err := gitOutput(dir, "show", tag+":./"+file)
	if err != nil {
		return
	}
	data, err := documentYAML(content, file)
	if err != nil {
		return
	}
	srd, err := parseSRDBytes(data, tag+":"+file)
	if err != nil || srd.ID == "" {
		return
	}
	t.srds[srd.ID] = file
	items := map[string]string{}
	for _, item := range srd.AllItems() {
		items[item.ID] = file
	}
	t.items[srd.ID] = items
	acs := map[string]string{}
	for _, ac := range srd.AcceptanceCriteria {
		acs[ac.ID] = file
	}
	t.acs[srd.ID] = acs
}

// taggedDocumentID reads one document at tag and returns its id field.
func taggedDocumentID(dir, tag, file string) string {
	content, err := gitOutput(dir, "show", tag+":./"+file)
	if err != nil {
		return ""
	}
	data, err := documentYAML(content, file)
	if err != nil {
		return ""
	}
	var doc struct {
		ID string `yaml:"id"`
	}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return ""
	}
	return doc.ID
}

// retiredInventory indexes every retired id at HEAD with its replacements.
type retiredInventory struct {
	srds     map[string][]string            // srd id -> replaced_by
	groups   map[string]map[string][]string // srd id -> item or ac id -> replaced_by
	useCases map[string][]string            // use case id -> replaced_by
	suites   map[string][]string            // test suite id -> replaced_by
}

func (r retiredInventory) empty() bool {
	return len(r.srds) == 0 && len(r.groups) == 0 &&
		len(r.useCases) == 0 && len(r.suites) == 0
}

// retiredCitationScan accumulates citation-of-retired-id warnings.
type retiredCitationScan struct {
	retired  retiredInventory
	findings []Finding
}

// checkCitationsOfRetiredIDs warns when a touchpoint, citation, or trace
// names a retired id; the message carries replaced_by so the fix is in the
// finding (design.yaml D6).
func checkCitationsOfRetiredIDs(corpus *Corpus) []Finding {
	scan := retiredCitationScan{retired: buildRetiredInventory(corpus)}
	if scan.retired.empty() {
		return nil
	}
	scan.scanUseCases(corpus)
	scan.scanTestSuites(corpus)
	return scan.findings
}

func (s *retiredCitationScan) warn(citer, id string, replacedBy []string) {
	replacement := "nothing"
	if len(replacedBy) > 0 {
		replacement = strings.Join(replacedBy, ", ")
	}
	s.findings = append(s.findings, Finding{
		Check: "citation-of-retired-id",
		Level: "warning",
		Message: fmt.Sprintf("%s cites retired id %s (replaced by %s)",
			citer, id, replacement),
	})
}

func (s *retiredCitationScan) checkRef(citer, srdID string, groups []string) {
	if replacedBy, ok := s.retired.srds[srdID]; ok {
		s.warn(citer, srdID, replacedBy)
		return
	}
	for _, group := range groups {
		if replacedBy, ok := s.retired.groups[srdID][group]; ok {
			s.warn(citer, srdID+" "+group, replacedBy)
		}
	}
}

func (s *retiredCitationScan) scanUseCases(corpus *Corpus) {
	for _, ucID := range sortedUseCaseIDs(corpus) {
		uc := corpus.UseCases[ucID]
		citer := "use case " + ucID
		for _, tp := range uc.Touchpoints {
			if srdID, groups := parseTouchpoint(tp); srdID != "" {
				s.checkRef(citer, srdID, groups)
			}
		}
		for _, sc := range uc.SuccessCriteria {
			for _, trace := range sc.Traces {
				if srdID, ref := parseACTrace(trace); srdID != "" {
					s.checkRef(citer+" "+sc.ID, srdID, []string{ref})
				}
			}
		}
		if replacedBy, ok := s.retired.suites[uc.TestSuite]; ok {
			s.warn(citer, uc.TestSuite, replacedBy)
		}
	}
}

func (s *retiredCitationScan) scanTestSuites(corpus *Corpus) {
	for _, tsID := range sortedTestSuiteIDs(corpus) {
		ts := corpus.TestSuites[tsID]
		citer := "test suite " + tsID
		for _, trace := range ts.Traces {
			if replacedBy, ok := s.retired.useCases[trace]; ok {
				s.warn(citer, trace, replacedBy)
				continue
			}
			if srdID, ref := parseACTrace(trace); srdID != "" {
				s.checkRef(citer, srdID, []string{ref})
			}
		}
		for _, tc := range ts.TestCases {
			caseCiter := citer + " case " + tc.Name
			if replacedBy, ok := s.retired.useCases[tc.UseCase]; ok {
				s.warn(caseCiter, tc.UseCase, replacedBy)
			}
			for _, trace := range tc.Traces {
				if srdID, ref := parseACTrace(trace); srdID != "" {
					s.checkRef(caseCiter, srdID, []string{ref})
				}
			}
		}
	}
}

func buildRetiredInventory(corpus *Corpus) retiredInventory {
	retired := retiredInventory{
		srds:     map[string][]string{},
		groups:   map[string]map[string][]string{},
		useCases: map[string][]string{},
		suites:   map[string][]string{},
	}
	for id, srd := range corpus.SRDs {
		if srd.Status == statusRetired {
			retired.srds[id] = srd.ReplacedBy
			continue
		}
		groups := map[string][]string{}
		for _, item := range srd.AllItems() {
			if item.Status == statusRetired {
				groups[item.ID] = item.ReplacedBy
			}
		}
		for _, ac := range srd.AcceptanceCriteria {
			if ac.Status == statusRetired {
				groups[ac.ID] = ac.ReplacedBy
			}
		}
		if len(groups) > 0 {
			retired.groups[id] = groups
		}
	}
	for id, uc := range corpus.UseCases {
		if uc.Status == statusRetired {
			retired.useCases[id] = uc.ReplacedBy
		}
	}
	for id, ts := range corpus.TestSuites {
		if ts.Status == statusRetired {
			retired.suites[id] = ts.ReplacedBy
		}
	}
	return retired
}

func sortedFileKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedSRDIDs(corpus *Corpus) []string {
	ids := make([]string, 0, len(corpus.SRDs))
	for id := range corpus.SRDs {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func sortedUseCaseIDs(corpus *Corpus) []string {
	ids := make([]string, 0, len(corpus.UseCases))
	for id := range corpus.UseCases {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func sortedTestSuiteIDs(corpus *Corpus) []string {
	ids := make([]string, 0, len(corpus.TestSuites))
	for id := range corpus.TestSuites {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
