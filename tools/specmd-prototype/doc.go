// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

// Package specmd prototypes the Markdown front end for the specification
// corpus (GH-2578): a YAML printer, a Markdown parser, and a canonical
// Markdown printer around the shared spec.SRD model, so an SRD round-trips
// between the two formats through one model. The format engine it builds on
// (markers, shapes, declarations, diagnostics) moved to
// agent-core/pkg/spec/markdown in GH-2602; the SRD printer, parser, and the
// check command stay here until the generic converter replaces them (GH-2604).
package specmd
