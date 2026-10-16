// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package specmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec/markdown"
)

// designPath is where the document-type declarations live, relative to the
// repository root.
const designPath = "agent-core/docs/constitutions/design.yaml"

// RepositoryDeclarations loads the Markdown declarations from the design.yaml
// of the repository enclosing dir. The library takes declarations as input;
// finding the file is the tool's job.
func RepositoryDeclarations(dir string) (markdown.Declarations, error) {
	path, err := findDesign(dir)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return markdown.LoadDeclarations(data)
}

// findDesign walks up from dir to the repository directory holding
// design.yaml.
func findDesign(dir string) (string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	for {
		candidate := filepath.Join(dir, designPath)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("%s not found above %s", designPath, dir)
		}
		dir = parent
	}
}
