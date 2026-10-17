---
id: rel00.0-uc002-spec-corpus-validation
title: Maintainer validates the specification corpus graph
test_suite: test-rel00.0
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A maintainer runs the spec graph loader and validator against the repository
documentation corpus. The validator parses SRDs, use cases, test suites,
the roadmap, and SPECIFICATIONS.md into a graph, then reports broken
references, orphaned SRDs, uncovered criteria, and release coverage gaps.

<!-- field: actor -->
## Actor

Specification maintainer

<!-- field: trigger -->
## Trigger

A maintainer changes docs/specs or runs mage audit before closing a
documentation task.

<!-- field: flow -->
## Flow

- F1: LoadCorpus discovers SRDs, use cases, test suites, road-map.md, and SPECIFICATIONS.md from the standard documentation layout (srd014 R4).
- F2: The parser reads each artifact into typed structs for SRDs, use cases, test suites, roadmap releases, and the spec index (srd014 R1).
- F3: BuildGraph creates release, SRD, requirement, acceptance criterion, use case, test suite, and test case nodes with contains, depends-on, traces, covers, and tests edges (srd014 R2).
- F4: Validate reports orphaned SRDs, broken touchpoints, invalid test citations, uncovered requirement items, uncovered acceptance criteria, dependency errors, and releases without test suites (srd014 R3).
- F5: The maintainer fixes the corpus or files follow-up issues for warnings that are intentionally deferred.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd014-spec-graph](../software-requirements/srd014-spec-graph.md) [R1.1](../software-requirements/srd014-spec-graph.md#r11), [R1.2](../software-requirements/srd014-spec-graph.md#r12), [R1.3](../software-requirements/srd014-spec-graph.md#r13), [R1.4](../software-requirements/srd014-spec-graph.md#r14), [R1.5](../software-requirements/srd014-spec-graph.md#r15), [R1.6](../software-requirements/srd014-spec-graph.md#r16) | typed parsers for SRDs, use cases, test suites, roadmap, and SPECIFICATIONS.md |
| T2 | [srd014-spec-graph](../software-requirements/srd014-spec-graph.md) [R2.1](../software-requirements/srd014-spec-graph.md#r21), [R2.2](../software-requirements/srd014-spec-graph.md#r22), [R2.3](../software-requirements/srd014-spec-graph.md#r23) | graph node kinds and relationship edges |
| T3 | [srd014-spec-graph](../software-requirements/srd014-spec-graph.md) [R3.1](../software-requirements/srd014-spec-graph.md#r31), [R3.2](../software-requirements/srd014-spec-graph.md#r32), [R3.3](../software-requirements/srd014-spec-graph.md#r33), [R3.4](../software-requirements/srd014-spec-graph.md#r34), [R3.5](../software-requirements/srd014-spec-graph.md#r35), [R3.6](../software-requirements/srd014-spec-graph.md#r36), [R3.7](../software-requirements/srd014-spec-graph.md#r37), [R3.8](../software-requirements/srd014-spec-graph.md#r38), [R3.9](../software-requirements/srd014-spec-graph.md#r39), [R3.10](../software-requirements/srd014-spec-graph.md#r310) | validation findings for broken, orphaned, and uncovered documentation links |
| T4 | [srd014-spec-graph](../software-requirements/srd014-spec-graph.md) [R4.1](../software-requirements/srd014-spec-graph.md#r41), [R4.2](../software-requirements/srd014-spec-graph.md#r42), [R4.3](../software-requirements/srd014-spec-graph.md#r43), [R4.4](../software-requirements/srd014-spec-graph.md#r44) | standard corpus discovery paths |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The loader parses all standard documentation artifacts into typed structs.

<!-- field: traces -->
#### Traces

- [srd014-spec-graph](../software-requirements/srd014-spec-graph.md) [AC1](../software-requirements/srd014-spec-graph.md#ac1)
- [srd014-spec-graph](../software-requirements/srd014-spec-graph.md) [AC4](../software-requirements/srd014-spec-graph.md#ac4)

<!-- sc: S2 -->
### S2

The graph contains the expected artifact, requirement, acceptance, use-case, suite, and test-case relationships.

<!-- field: traces -->
#### Traces

- [srd014-spec-graph](../software-requirements/srd014-spec-graph.md) [AC2](../software-requirements/srd014-spec-graph.md#ac2)

<!-- sc: S3 -->
### S3

The validator reports actionable findings with check identifiers and messages.

<!-- field: traces -->
#### Traces

- [srd014-spec-graph](../software-requirements/srd014-spec-graph.md) [AC3](../software-requirements/srd014-spec-graph.md#ac3)

<!-- field: out_of_scope -->
## Out of scope

- Generating new specifications.
- Changing audit warning policy.
- Interactive graph visualization.
