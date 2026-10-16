<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

# Markdown spec compiler: prototype findings

This note records what the prototype in this directory established for the spec-format migration (GH-2578), so the format specification (GH-2581) is written from evidence. The processor came in GH-2588 and the diagnostics in GH-2589; the measurements below were taken on 2026-10-15 against main at a36ed4e4.

## Verdict on the epic's decision metric

The epic asked whether Tier 3 shape rules stay around a dozen with clean failure modes, or whether the YAML-islands and render-only fallbacks need reconsidering. They stay around a dozen. The SRD declaration needs 12 field-content rules, and each fails with one named finding that quotes its rule. Marker attributes add 8 more rules, all checked mechanically (a JSON shape or a count) rather than by judgment. We see no case for the fallbacks.

## Round trip

The four functions around the shared model hold their laws on srd004: yaml to model to Markdown to model loses no field, Markdown is a fixed point of parse-then-print, model to yaml to model is the identity, and the formatter is idempotent and converges non-canonical spellings. Over the whole corpus, 105 of 106 tracked SRDs pass the same laws. The exception, catalog srd021, has a duplicate mapping key in its source YAML; GH-2590 tracks it with ten values that unquoted commas silently split in flow mappings. All eleven are defects in the corpus, not in the processor.

Table 1: Lossy spots and normalization

| Spot | Effect | Disposition |
| --- | --- | --- |
| Trailing newline of a YAML literal block | Markdown prose cannot express it; text compares with trailing newlines trimmed | The only whitespace normalization; recorded in the round-trip tests |
| Document extensions | Move after the typed sections in canonical order (srd004's applicability now follows the criteria) | Order of untyped keys among typed fields is not modeled; within extensions, source order holds |
| Values Markdown would misread or cannot shape | Render as a fenced yaml block: 159 blocks across the corpus | Lossless, less readable; nested values such as applicability, diagrams, implementation_status |
| Item form and group form | Tagged versus object items and mapped versus listed groups are not modeled; the printer chooses by content | No data lost; the YAML spelling may change |

## Finding catalog

Every finding has a stable name, a file line, expected and found text, the rule path in rules/srd.yaml, and a fix hint. Findings are collected, never fail-fast. Tier 0, canonical form, is never a finding: the formatter fixes it.

Table 2: Finding names by tier

| Tier | Names |
| --- | --- |
| 1, grammar | frontmatter-missing, frontmatter-invalid, frontmatter-missing-key, frontmatter-unknown-key, frontmatter-id-mismatch, marker-malformed, marker-unknown-tag, marker-unknown-attribute, marker-orphan |
| 2, structure | heading-untagged, heading-misplaced, id-malformed, item-wrong-group, heading-id-mismatch, duplicate-id, field-disallowed, field-duplicate, content-untagged, untagged-structure, missing-content, missing-required-field, missing-requirement-groups |
| 3, shapes | field-shape-mismatch, list-item-prefix, table-columns-mismatch, attribute-shape, criterion-traces-missing, criterion-trace-invalid |
| 4, prose | prose-containment |

## Tier 3 shape rules

Table 3: Tier 3 rule inventory

| Rule | Count | Rules |
| --- | --- | --- |
| Field content shape | 8 | problem prose, goals list, non_goals list, assumptions list, acceptance_criteria empty, needs_clarification table, evidence prose, endpoints table |
| List item prefix | 2 | goals G<n>:, assumptions A<n>: |
| Table columns | 2 | needs_clarification id, question; endpoints path, owner |
| Marker attribute shape | 6 | weight integer, status text, retired_in text, replaced_by list, derived_from references, traces list |
| Criterion traces | 2 | at least one trace; each trace a requirement id |

The first three rows are the 12 field-content rules the metric counts; the last two are the 8 attribute rules.

## What the diagnostics found in the corpus

The printer never produces a structural finding: no marker, frontmatter, heading, or placement finding on the printed form of any tracked SRD (TestPrintedCorpusIsStructurallyClean). Two findings report the corpus's own content.

The id grammar rejects 18 ids in 3 SRDs: agent-core srd008 uses R3.1b, srd013 uses three-level ids R10.0.1 through R10.1.6, and srd019 uses R3.1.1 through R3.1.6. Design.yaml declares R1.1 numbering, and D6 forbids renumbering an id a release contains, so GH-2581 must either widen the id grammar to these forms or retire and replace them with pointers. Renumbering in place is not open.

Prose containment reports 15 hits in 5 SRDs, 14 of them in non-goals phrased as obligations ("the applier SHALL NOT ..."). A non-goal written with shall states a requirement no test traces. GH-2581 decides whether these become requirement items or are reworded; exempting non-goals would hide them.

Five catalog SRDs are relocated stubs with empty required sections. The rules exempt status relocated and retired from the required-content checks rather than demand empty sections.

## Rendering

On GitHub (observed on the converted srd004): the frontmatter renders as a two-column table of names and values; no marker or license comment is visible; the outline reads as a document, with 14 level-2 sections, 39 level-3 items and criteria, and level-4 Evidence and Endpoints; the endpoint list renders as a table and the one fallback as a highlighted yaml block. Criterion traces live only in their markers, so a reader cannot see which requirements a criterion verifies; GH-2581 should decide whether traces belong in the visible layer.

On GitLab: not observed; this repository has no GitLab mirror. GitLab documents that it hides HTML comments and shows front matter in a box above the rendered document, which would keep every marker invisible and show the frontmatter as code rather than a table. A rendering check on a GitLab instance remains open.
