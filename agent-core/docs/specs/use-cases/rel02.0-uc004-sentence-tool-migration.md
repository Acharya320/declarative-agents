---
id: rel02.0-uc004-sentence-tool-migration
title: Developer migrates a sentence-shaped tool into machine-composed words
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A developer finds a legacy tool that hides a workflow behind one name. The
developer decomposes it into single-purpose tools, declares each word's
contract, rewires machine.yaml to express the sequence through states and
signals, keeps any aggregate only while compatibility callers remain, and
removes it once callers migrate.

<!-- field: actor -->
## Actor

Developer maintaining tool declarations and machines

<!-- field: trigger -->
## Trigger

A tool audit or implementation review identifies a tool that performs
multiple observable steps, such as loading config and discovering samples, or
running checks and summarizing results.

<!-- field: flow -->
## Flow

- F1: The developer classifies the existing tool as a legacy sentence because it hides sequencing that should be visible in MachineSpec (cfg-tool-authoring-standard, sm-tool-language).
- F2: The developer identifies each observable verb and creates one ToolDef per word with name, type, init/binary, parameters, emits, output, side_effects, reversibility, undo, errors, and relationships (srd022 R0, srd023 R0).
- F3: The developer moves order into machine.yaml by adding states and transitions for each word's emitted signals (srd006 R3, R4).
- F4: The developer updates tool selection files so the machine selects the new words instead of the aggregate (cfg-tool-selection-format).
- F5: Boundary or nested-machine tools remain allowed only when the boundary itself is the word and side effects are explicit (srd024 R0, srd025 R8).
- F6: The developer updates semantic/audit docs to mark the aggregate as deprecated_compatibility while callers remain, then removes its implementation, registration, and declaration after migration (sm-tool-vocabulary-audit).
- F7: The developer adds tests or future test cases proving the new machine sequence preserves behavior and exposes intermediate failure signals.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd006-yaml-machine-config R3.1, R3.2, R3.3, R3.4, R3.5, R4.1, R4.2, R4.3, R4.4 | machine.yaml owns visible sequencing and transition validation |
| T2 | srd013-standard-tool-library R8.1, R8.2, R8.3, R10.0, R10.0.1, R10.0.2, R10.0.3, R10.0.4, R10.0.5, R10.1, R10.1.1, R10.1.2, R10.1.3, R10.1.4, R10.1.5, R10.1.6, R10.2, R10.3, R10.4, R10.5, R10.6, R10.7, R10.8 | validation and evaluator workflows use machine-declared atomic tools |
| T3 | srd022-builtin-tools R0.1, R0.2, R0.3, R4.1, R4.2, R4.3, R4.4, R4.5, R4.6, R4.7, R5.1, R5.2, R5.3, R5.4 | builtin words declare emits and reusable metadata |
| T4 | srd023-cli-exec-tools R0.1, R0.2, R5.1, R5.2, R5.3, R5.4 | exec tools model one CLI verb before bespoke wrappers |
| T5 | srd025-rollback-lifecycle R3.1, R3.2, R3.3, R3.4, R3.5, R8.1, R8.2, R8.3, R8.4 | split tools declare receipt-aware undo and side effects for rollback |
| T6 | sm-tool-language | migration states for legacy_sentence and deprecated_compatibility |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The old aggregate workflow is represented as multiple declared word tools.

<!-- field: traces -->
#### Traces

- srd013-standard-tool-library AC9

<!-- sc: S2 -->
### S2

machine.yaml expresses the sequence and handles each declared success/failure signal.

<!-- field: traces -->
#### Traces

- srd006-yaml-machine-config AC3

<!-- sc: S3 -->
### S3

The aggregate tool is removed after all active selections and fixtures migrate, or documented as deprecated compatibility only while callers remain.

<!-- field: traces -->
#### Traces

- srd022-builtin-tools AC1

<!-- sc: S4 -->
### S4

Each replacement word declares side effects, reversibility, undo, and relationships for rollback and composition.

<!-- field: traces -->
#### Traces

- srd025-rollback-lifecycle AC6

<!-- field: failure_scenarios -->
## Failure scenarios

| id | scenario | expected |
| --- | --- | --- |
| E1 | A split word lacks an emitted signal used by machine.yaml. | Startup validation fails before the machine runs. |
| E2 | The migration changes behavior by dropping an intermediate error path. | Tests or formal cases fail until the missing transition is restored. |
| E3 | The legacy aggregate remains registered after every production machine and fixture has migrated. | Tool vocabulary audit flags the dead compatibility runtime for removal. |

<!-- field: out_of_scope -->
## Out of scope

- Removing compatibility tools before all callers migrate.
- Refactoring helper functions that are not user-visible tools.
- Changing the semantics of the original workflow beyond exposing its steps.

<!-- field: semantic_model -->
## Semantic model

sm-tool-language

<!-- field: test_suite_trace -->
## Test suite trace

test-rel02.0-tool-contracts
