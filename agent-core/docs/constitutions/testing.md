---
id: testing
title: Testing Constitution
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: purpose -->
## Purpose

We use this constitution to keep YAML test suites, Go tests, integration
targets, and acceptance criteria aligned.

<!-- field: articles -->
## Articles

<!-- article: T1 -->
### T1 - Test the contract

Tests should validate observable behavior, public contracts, and
specification promises. Do not couple tests to private implementation
details unless the package is itself the implementation boundary.

<!-- article: T2 -->
### T2 - Trace YAML suites to Go tests

Every formal YAML test case should name a Go test, mage target, or
integration target that covers it. Every release use case should either
point to a formal test suite or explain why Go-only coverage is the
current state.

<!-- article: T3 -->
### T3 - Preserve release coverage

Adding a use case requires adding or updating the release test suite.
Closing a release requires matching SPECIFICATIONS.md counts, use case
status, and test-suite entries.

<!-- article: T4 -->
### T4 - Keep tests isolated

Tests write only to t.TempDir unless a fixture path is read-only. Tests
that change process state, environment, current directory, global
variables, or shared files must not run in parallel.

<!-- article: T5 -->
### T5 - Gate by task risk

Code issues run the full code gate unless the issue states a narrower
gate. Documentation issues run mage audit and de-AI checks on touched
docs. Integration targets run when the issue changes profile behavior,
machine routing, evaluator flow, bench behavior, lifecycle operations, or
model-boundary config.

<!-- field: go_test_style | form: "tree" -->
## Go test style

- table_driven_tests:
  - rule: Use named rows and t.Run for repeated cases.
  - example:
    tests := []struct {
      name string
      input string
      want string
    }{
      {name: "empty", input: "", want: ""},
    }
    for _, tc := range tests {
      t.Run(tc.name, func(t *testing.T) {
        t.Parallel()
        require.Equal(t, tc.want, transform(tc.input))
      })
    }
- parallelism:
  - use_t_parallel_when:
    1. The test uses only local values.
    1. The test writes only to its own t.TempDir.
    1. The test does not mutate global state.
    1. The test does not call os.Chdir.
  - do_not_parallelize_when:
    1. The test changes process working directory.
    1. The test mutates environment variables without isolated cleanup.
    1. The test uses shared package state.
    1. The test binds a fixed network port.
- helpers:
  - rules:
    1. Every helper calls t.Helper() first.
    1. Helpers create clear failure messages.
    1. Shared fixtures should be read-only.
    1. Use t.Cleanup for temporary process or filesystem state.
- assertions:
  - rules:
    1. Use require for setup failures and fatal preconditions.
    1. Use assert for independent checks after setup succeeds.
    1. Check error messages when the contract promises specific context.

<!-- field: yaml_test_suites | form: "tree" -->
## YAML test suites

- required_fields:
  1. id
  1. title
  1. release
  1. traces
  1. preconditions
  1. test_cases
- test_case_fields:
  1. use_case
  1. name
  1. go_test
  1. inputs
  1. expected
  1. traces
- trace_formats:
  - srd_acceptance: srdNNN-name ACN
  - srd_requirement: srdNNN-name RN.N
  - use_case_success: relNN.N-ucNNN-name SN

<!-- field: mage_targets | form: "tree" -->
## Mage targets

- unit: go test ./...
- build: go build ./...
- vet: go vet ./...
- lint: mage lint
- audit: mage audit
- integration:
  1. mage integration:monitor
  1. mage integration:ollamaRest
  1. mage integration:ollamaMonitor
  1. mage integration:dolt
  1. mage integration:doltWord
  1. mage integration:all

<!-- field: quality_gates | form: "tree" -->
## Quality gates

- code:
  1. go build ./...
  1. go vet ./...
  1. mage lint
  1. go test ./...
- documentation:
  1. mage audit
  1. bash .cursor/skills/de-ai/scripts/detect-lexical.sh <touched-doc>
  1. python3 .cursor/skills/de-ai/scripts/detect-structural.py <touched-doc>

<!-- field: documentation_de_ai_scope | form: "tree" -->
## Documentation de ai scope

- lexical: Full touched-document lexical scans are blocking for documentation issues.
- structural:
  Structural scans are blocking for touched prose blocks in Markdown and YAML
  block scalars when they report suspicious or likely-ai prose that can be
  rewritten without changing formal meaning.

  Formal reference docs may keep whole-file structural warnings caused by
  required YAML shape. Required lists, field inventories, trace matrices, and
  examples can trigger those metrics. Compact status sections may show the
  same pattern. Accept residual warnings about list density and punctuation
  density when lexical scans are clean.
  The same exception covers inline inventories, parenthetical text, and
  paragraph length when three focused rewrite passes would make the reference
  less precise.
- close_comment: Record any accepted formal-document residuals in the issue stats comment.
- config_or_profile:
  1. mage audit
  1. go test ./...
  1. relevant mage integration target when profile behavior changes

<!-- field: coverage_notes -->
## Coverage notes

- Release 00.0 has formal YAML coverage through test-rel00.0.
- Release 01.0 currently relies on Go and integration tests until a formal suite is added.
- Release 02.0 has formal lifecycle and tool-contract suites, plus Go tests for implemented runtime pieces.
- Release 14.0 has implemented request-signal unit, HTTP, lifecycle, mock-model compatibility, and real-Dolt evidence in test-rel14.0-agent-input-and-dolt.

<!-- field: references -->
## References

- docs/SPECIFICATIONS.md
- docs/specs/test-suites/test-rel00.0.md
- docs/specs/test-suites/test-rel02.0-lifecycle.md
- docs/specs/test-suites/test-rel02.0-tool-contracts.md
- docs/constitutions/go-style.md
