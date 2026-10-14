// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package spec

import (
	"fmt"
	"path"
)

// checkUnresolvedClarifications reports every needs_clarification entry in
// the corpus (design.yaml clarification_marker): each one is a decision
// nobody has made, and a document carrying one is not ready. One error per
// entry, naming the document, the blocked id, and the question; an entry
// whose id names nothing in the document still fires, saying the target is
// unknown (srd014 R3.14).
func checkUnresolvedClarifications(corpus *Corpus) []Finding {
	var findings []Finding
	for _, srdID := range sortedSRDIDs(corpus) {
		srd := corpus.SRDs[srdID]
		file := path.Join(SRDSubdir, srdID+".yaml")
		findings = append(findings, clarificationFindings(
			srdID, file, srd.NeedsClarification, srdClarificationTargets(&srd))...)
	}
	for _, ucID := range sortedUseCaseIDs(corpus) {
		uc := corpus.UseCases[ucID]
		file := path.Join(UCSubdir, ucID+".yaml")
		targets := map[string]bool{}
		for _, sc := range uc.SuccessCriteria {
			targets[sc.ID] = true
		}
		findings = append(findings, clarificationFindings(
			ucID, file, uc.NeedsClarification, targets)...)
	}
	for _, tsID := range sortedTestSuiteIDs(corpus) {
		ts := corpus.TestSuites[tsID]
		file := path.Join(TSSubdir, tsID+".yaml")
		targets := map[string]bool{}
		for _, tc := range ts.TestCases {
			targets[tc.Name] = true
		}
		findings = append(findings, clarificationFindings(
			tsID, file, ts.NeedsClarification, targets)...)
	}
	return findings
}

// srdClarificationTargets collects the ids a clarification entry may block
// in an SRD: requirement groups, requirement items, and acceptance criteria.
func srdClarificationTargets(srd *SRD) map[string]bool {
	targets := map[string]bool{}
	for _, group := range srd.OrderedGroups {
		targets[group] = true
	}
	for _, item := range srd.AllItems() {
		targets[item.ID] = true
	}
	for _, ac := range srd.AcceptanceCriteria {
		targets[ac.ID] = true
	}
	return targets
}

// clarificationFindings emits one finding per entry against one document.
func clarificationFindings(docID, file string, entries []ClarificationEntry, targets map[string]bool) []Finding {
	var findings []Finding
	for _, entry := range entries {
		message := fmt.Sprintf(
			"%s %s is blocked on an unmade decision: %s -- make the decision, edit the entry in place, and delete the needs_clarification entry (design.yaml clarification_marker)",
			docID, entry.ID, entry.Question)
		if !targets[entry.ID] {
			message = fmt.Sprintf(
				"%s carries a needs_clarification entry for unknown id %q: %s -- the entry names nothing in the document (design.yaml clarification_marker)",
				docID, entry.ID, entry.Question)
		}
		findings = append(findings, Finding{
			Check:   "unresolved-clarification",
			Level:   "error",
			Message: message,
			File:    file,
		})
	}
	return findings
}
