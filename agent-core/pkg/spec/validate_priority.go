// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package spec

import (
	"fmt"
	"path"
)

// validPriorities is the use-case priority vocabulary (design.yaml
// use_case priority): P1 must ship for its release, P2 should, P3 could.
// Absent means unranked.
var validPriorities = map[string]bool{"": true, "P1": true, "P2": true, "P3": true}

// checkInvalidPriorities reports every use-case priority outside the
// P1/P2/P3 vocabulary (invalid-priority, srd014 R3.15).
func checkInvalidPriorities(corpus *Corpus) []Finding {
	var findings []Finding
	for _, ucID := range sortedUseCaseIDs(corpus) {
		uc := corpus.UseCases[ucID]
		if validPriorities[uc.Priority] {
			continue
		}
		findings = append(findings, Finding{
			Check: "invalid-priority",
			Level: "error",
			Message: fmt.Sprintf(
				"use case %s has priority %q; valid values are P1, P2, and P3, and an absent priority means unranked (design.yaml use_case priority)",
				ucID, uc.Priority),
			File: path.Join(UCSubdir, ucID+".yaml"),
		})
	}
	return findings
}
