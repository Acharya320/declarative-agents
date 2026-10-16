// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

// Package markdown is the format engine of the Markdown specification corpus
// (docs/engineering/eng10-markdown-spec-format.md): the marker grammar, the
// block scanner, the shape classifier and renderers, the declaration loader
// for design.yaml's document_types.<type>.markdown entries, and the tiered
// diagnostics with their stable finding catalog. Callers load the
// declarations and pass them in; the package never searches for design.yaml.
package markdown
