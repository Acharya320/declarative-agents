// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"os"
	"sync"
	"testing"
)

// designPath is agent-core's design constitution, relative to this package.
const designPath = "../../../docs/constitutions/design.yaml"

var (
	designOnce         sync.Once
	designDeclarations Declarations
	designErr          error
)

// repositoryDeclarations loads the declarations the tests check documents
// against, once per test binary.
func repositoryDeclarations(t *testing.T) Declarations {
	t.Helper()
	designOnce.Do(func() {
		data, err := os.ReadFile(designPath)
		if err != nil {
			designErr = err
			return
		}
		designDeclarations, designErr = LoadDeclarations(data)
	})
	if designErr != nil {
		t.Fatal(designErr)
	}
	return designDeclarations
}

// diagnose runs the diagnostics against the repository's declarations.
func diagnose(t *testing.T, data []byte, source string) []Finding {
	t.Helper()
	return Diagnose(repositoryDeclarations(t), data, source)
}
