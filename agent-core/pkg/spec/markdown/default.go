// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"sync"

	agentcore "github.com/Nokia-Bell-Labs/declarative-agents/agent-core"
)

var (
	defaultOnce         sync.Once
	defaultDeclarations Declarations
	defaultErr          error
)

// Default returns the declarations of the design constitution this module
// was built with, loaded once. Readers that need no other declarations, such
// as the corpus loader, use these.
func Default() (Declarations, error) {
	defaultOnce.Do(func() {
		defaultDeclarations, defaultErr = LoadDeclarations(agentcore.DesignConstitution())
	})
	return defaultDeclarations, defaultErr
}
