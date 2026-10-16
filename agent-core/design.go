// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

// Package agentcore carries the Agent Core design constitution into the code
// that reads the specification corpus. The constitution's document_types
// registry declares every corpus document type and its Markdown format
// (docs/engineering/eng10-markdown-spec-format.md); embedding it means a
// binary, a container image, or a consumer repository reads the same
// declarations without locating the file.
package agentcore

import _ "embed"

//go:embed docs/constitutions/design.yaml
var designConstitution []byte

// DesignConstitution returns the bytes of docs/constitutions/design.yaml as
// this module was built with them.
func DesignConstitution() []byte {
	return designConstitution
}
