// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package specmd

// FormatMarkdown rewrites an SRD Markdown document in canonical form: the
// formatter is the parser followed by the canonical printer, so anything the
// parser accepts converges on one spelling.
func FormatMarkdown(data []byte, source string) ([]byte, error) {
	srd, err := ParseMarkdown(data, source)
	if err != nil {
		return nil, err
	}
	return PrintMarkdown(&srd)
}
