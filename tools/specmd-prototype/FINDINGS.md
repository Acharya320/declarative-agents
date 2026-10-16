<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

# Markdown spec compiler: prototype findings

This note records what the prototype in this directory established for the spec-format migration (GH-2578), so the format specification (GH-2581) is written from evidence. The processor came in GH-2588 and the diagnostics in GH-2589. GH-2597 moved the rules into design.yaml (the markdown block of each document type), added the engineering_guideline declaration, and applied GH-2581's decisions on the id grammar and criterion traces. GH-2581 declared the remaining document types, which made a file select its declaration by its whole path rather than its name, let a list or table field be empty, and let a table field leave its columns undeclared. GH-2602 moved the format engine and its diagnostics tests to agent-core/pkg/spec/markdown. The measurements below were taken on 2026-10-15.

## Verdict on the epic's decision metric

The epic asked whether Tier 3 shape rules stay around a dozen with clean failure modes, or whether the YAML-islands and render-only fallbacks need reconsidering. They stay around a dozen. The SRD declaration needs 14 field-content rules, and each fails with one named finding that quotes its rule. Marker attributes add 5 more, all checked mechanically (a JSON shape) rather than by judgment. A second document type, the engineering guideline, needed a declaration and no engine code. We see no case for the fallbacks.

## Round trip

The four functions around the shared model hold their laws on srd004: yaml to model to Markdown to model loses no field, Markdown is a fixed point of parse-then-print, model to yaml to model is the identity, and the formatter is idempotent and converges non-canonical spellings. Over the whole corpus, all 106 tracked SRDs pass the same laws. The sweep first failed on catalog srd021, whose source YAML had a duplicate mapping key, and it led to ten values that unquoted commas silently split in flow mappings. All eleven were defects in the corpus, not in the processor, and GH-2590 fixed them.

Table 1: Lossy spots and normalization

| Spot | Effect | Disposition |
| --- | --- | --- |
| Trailing newline of a YAML literal block | Markdown prose cannot express it; text compares with trailing newlines trimmed | The only whitespace normalization; recorded in the round-trip tests |
| Document extensions | Move after the typed sections in canonical order (srd004's applicability now follows the criteria) | Order of untyped keys among typed fields is not modeled; within extensions, source order holds |
| Values Markdown would misread or cannot shape | Render as a fenced yaml block: 159 blocks across the corpus | Lossless, less readable; nested values such as applicability, diagrams, implementation_status |
| Item form and group form | Tagged versus object items and mapped versus listed groups are not modeled; the printer chooses by content | No data lost; the YAML spelling may change |

## Finding catalog

Every finding has a stable name, a file line, expected and found text, the rule path in design.yaml, and a fix hint. Findings are collected, never fail-fast. Tier 0, canonical form, is never a finding: the formatter fixes it.

Table 2: Finding names by tier

| Tier | Names |
| --- | --- |
| 1, grammar | document-type-unknown, frontmatter-missing, frontmatter-invalid, frontmatter-missing-key, frontmatter-unknown-key, frontmatter-id-mismatch, marker-malformed, marker-unknown-tag, marker-unknown-attribute, marker-orphan |
| 2, structure | heading-untagged, heading-misplaced, id-malformed, id-outside-parent, heading-id-mismatch, duplicate-id, field-disallowed, field-duplicate, content-untagged, untagged-structure, missing-content, missing-required-field, missing-required-node |
| 3, shapes | field-shape-mismatch, list-item-pattern, table-columns-mismatch, attribute-shape |
| 4, prose | prose-containment |

## Tier 3 shape rules

Table 3: Tier 3 rule inventory

| Rule | Count | Rules |
| --- | --- | --- |
| Field content shape | 9 | problem prose, goals list, non_goals list, assumptions list, acceptance_criteria empty, needs_clarification table, evidence prose, endpoints table, traces list |
| List item pattern | 3 | goals G<n>:, assumptions A<n>:, traces a requirement id |
| Table columns | 2 | needs_clarification id, question; endpoints path, owner |
| Marker attribute shape | 5 | weight integer, status text, retired_in text, replaced_by list, derived_from references |

The first three rows are the 14 field-content rules the metric counts; the last is the 5 attribute rules. Criterion traces moved from a marker attribute to a visible field (GH-2581 decision 3), so their two former checks are now a required field and a list item pattern.

## What the diagnostics found in the corpus

On the printed form of all 106 tracked SRDs, the only findings left are prose containment (TestPrintedCorpusIsClean). The printer is structurally sound, every id fits the declared grammar, and all 583 acceptance criteria carry a visible traces field.

The declared R1.1 numbering first rejected 18 ids in 3 SRDs: agent-core srd008 uses R3.1b, srd013 uses three-level ids R10.0.1 through R10.1.6, and srd019 uses R3.1.1 through R3.1.6. D6 forbids renumbering an id a release contains, so GH-2581 widened the grammar instead (decision 1): an item id has two or three dotted levels and an optional lowercase suffix, and R1.1 stays canonical for new items.

Prose containment reports 15 hits in 5 SRDs, 14 of them in non-goals phrased as obligations ("the applier SHALL NOT ..."). A non-goal written with shall states a requirement no test traces. The check stays on non-goals (decision 2); each hit becomes a requirement item or is reworded during the corpus migration (GH-2584).

Five catalog SRDs are relocated stubs with empty required sections. The rules exempt status relocated and retired from the required-content checks rather than demand empty sections.

## Rendering

On GitHub (observed on the converted srd004): the frontmatter renders as a two-column table of names and values; no marker or license comment is visible; the outline reads as a document, with 14 level-2 sections, 39 level-3 items and criteria, and level-4 Evidence and Endpoints; the endpoint list renders as a table and the one fallback as a highlighted yaml block. Criterion traces were first written only in their markers, which left a reader unable to see what a criterion verifies; they now render as a visible list under each criterion (decision 3).

On GitLab: not observed; this repository has no GitLab mirror. GitLab documents that it hides HTML comments and shows front matter in a box above the rendered document, which would keep every marker invisible and show the frontmatter as code rather than a table. A rendering check on a GitLab instance remains open.
