---
id: execution
title: Work Execution Constitution
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: purpose -->
## Purpose

We use this constitution when completing one issue end-to-end. It
governs issue selection, scope, reading order, quality gates, issue-tracker updates,
and commits inside the selected project directory.

<!-- field: articles -->
## Articles

<!-- article: E1 -->
### E1 - Work from the issue tracker

Use the issue tracker for task tracking. Find available work, read the
selected issue, then claim it before starting. Do not use markdown TODO
lists as the system of record.

<!-- article: E2 -->
### E2 - Read first

Before editing, read the issue, this constitution set, and every
required_reading path named by the issue. For code work, read the relevant
SRDs, architecture sections, test-suite YAML, and source files. For
documentation work, read the governing document type in design.yaml.

<!-- article: E3 -->
### E3 - Stay in scope

Touch only files listed by the issue unless the issue requirement cannot
be satisfied without a related metadata update. When that happens, keep
the extra edit minimal and mention it in the close comment.

<!-- article: E4 -->
### E4 - Preserve user changes

The working tree may already contain user changes. Do not revert changes
you did not make. If a user change touches the same file, read it and work
with it.

<!-- article: E5 -->
### E5 - Keep specs and code aligned

Code implements documented contracts. Documentation describes current
code, config, and known gaps. Do not mark planned behavior as done unless
code and tests support it.

<!-- article: E6 -->
### E6 - Quality gates

Verify every acceptance criterion. For code, run go build ./..., go vet
./..., mage lint, and go test ./... unless the issue narrows the gate.
For documentation, run mage audit. If a gate fails because of unrelated
existing drift, file a follow-up issue and close the current issue
only when its own changes are valid.

<!-- article: E7 -->
### E7 - Documentation de-AI gate

For documentation issues, run the lexical and structural de-AI scripts on
each created or changed documentation file. Fix reported issues up to
three passes unless the flagged text is a formal requirement that should
not be rewritten.

<!-- field: formal_reference_scope -->
#### Formal reference scope

Whole-file structural metrics are advisory for formal
reference documents whose required shape is YAML, tables, enumerated
requirements, field lists, or trace lists. This includes ARCHITECTURE,
SPECIFICATIONS, road-map, SRDs, config-format specs, semantic models, use
cases, and test suites. For those files, treat lexical findings and
suspicious or likely-ai findings in touched prose blocks as the blocking
gate. Record residual list-heavy, colon-heavy, tricolon-heavy,
exhaustive-list-heavy, and uniform-paragraph findings as accepted formal
document structure when three focused rewrite passes would weaken the
reference format.

<!-- article: E8 -->
### E8 - Function size limit

Keep Go functions at or below 40 lines unless existing code already
exceeds the limit and the issue does not touch that function. Split new
logic by validation, parsing, formatting, I/O, or branch-specific work.

<!-- article: E9 -->
### E9 - File size limit

Keep Go files near or below 500 lines of production code. Split by
concern when adding new behavior to a large file. Do not create files
named helpers.go, utils.go, or types.go.

<!-- article: E10 -->
### E10 - Duplication limit

The duplication threshold is two. If a second copy of logic appears,
extract a shared function, method, table, or config helper.

<!-- article: E11 -->
### E11 - Declaration file size limit

Keep declaration YAML files at or below 300 lines. One concept belongs
in one file; a longer file is several concepts or a repeated shape, so split
it or extract a fragment rather than extend it. internal/declstyle holds a
shrink-only baseline of the files over the limit (GH-2111).

<!-- field: workflows | form: "tree" -->
## Workflows

- documentation:
  - steps:
    1. Read design.yaml and the issue required_reading.
    1. Identify the output path, document type, required fields, and naming rule.
    1. Write the document in valid YAML or Markdown as required.
    1. Run mage audit.
    1. Run de-AI lexical and structural checks on touched docs.
    1. Add a stats comment, close the issue, and commit the work.
- code:
  - steps:
    1. Read the relevant SRD, ARCHITECTURE section, test suite, source files, and issue.
    1. Implement only the listed files.
    1. Run go build ./...
    1. Run go vet ./...
    1. Run mage lint
    1. Run go test ./...
    1. Run mage audit if docs changed.
    1. Add a stats comment, close the issue, and commit the work.

<!-- field: git | form: "tree" -->
## Git

- commit_per_issue: true
- push_automatically: false
- commit_message:
  Use a concise message that names the outcome and, when applicable, the
  issue ID or SRD. Do not batch unrelated issues into one commit.

<!-- field: references -->
## References

- docs/constitutions/design.yaml
- docs/constitutions/go-style.md
- docs/constitutions/testing.md
- AGENTS.md
