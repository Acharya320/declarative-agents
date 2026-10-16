// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec/markdown"
)

const specCheckDesign = "docs/constitutions/design.yaml"

// SpecCheck diagnoses Markdown corpus documents against the declarations in
// docs/constitutions/design.yaml (eng10) and fails when any finding remains.
// Paths are comma-separated and relative to agent-core, for example
// mage specCheck ../docs/engineering/eng10-markdown-spec-format.md.
func SpecCheck(paths string) error {
	data, err := os.ReadFile(specCheckDesign)
	if err != nil {
		return err
	}
	declarations, err := markdown.LoadDeclarations(data)
	if err != nil {
		return err
	}
	total := 0
	for _, path := range strings.Split(paths, ",") {
		document, err := os.ReadFile(strings.TrimSpace(path))
		if err != nil {
			return err
		}
		for _, finding := range markdown.Diagnose(declarations, document, strings.TrimSpace(path)) {
			fmt.Println(finding)
			total++
		}
	}
	if total > 0 {
		return fmt.Errorf("%d finding(s)", total)
	}
	return nil
}
