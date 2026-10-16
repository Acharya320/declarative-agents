---
id: rel02.0-uc003-tool-authoring
title: Developer authors a new tool as a vocabulary word
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A developer needs to expose a new capability to an agent. Before writing the
implementation, the developer authors a ToolDef using the tool authoring
standard: the new tool has a clear problem, goals, requirements, non-goals,
typed input/output, declared side effects, reversibility, undo, errors, and
relationships. The tool is classified as a word, stateful_internal,
boundary, or legacy sentence, and the MachineSpec composes it with other
tools through explicit signals.

<!-- field: actor -->
## Actor

Developer adding an agent capability

<!-- field: trigger -->
## Trigger

The developer identifies a vocabulary gap in an agent workflow, such as a
missing parser, checker, artifact collector, or external CLI verb.

<!-- field: flow -->
## Flow

- F1: The developer states the problem and non-goals to prove the tool fills one vocabulary gap rather than hiding a workflow (cfg-tool-authoring-standard).
- F2: The developer classifies the tool as word, stateful_internal, boundary, or sentence migration candidate (srd022 R0, srd023 R0, srd024 R0).
- F3: The developer defines parameters and output.schema so the next machine transition or tool can consume the result without prose interpretation (cfg-tool-authoring-standard).
- F4: The developer declares emits so named MachineSpec actions can validate returned signals (srd022 R4, srd023 R3).
- F5: The developer declares side_effects, reversibility, undo, errors, and state_after_failure using the current rollback lifecycle model (srd025 R3, R8).
- F6: If the tool is irreversible or compensatable, the developer documents confirmation or compensation behavior before selecting it in user-affecting workflows (srd025 R8).
- F7: The developer documents relationships: likely preceding tools, likely consumers, and overlaps with similar words (cfg-tool-authoring-standard).
- F8: The developer wires the tool into machine.yaml as one action in a visible sequence rather than embedding sequencing inside the tool (srd006, sm-tool-language).

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd006-yaml-machine-config R3.1, R3.2, R3.3, R3.4, R3.5, R4.1, R4.2, R4.3, R4.4 | MachineSpec composes declared action names and signals |
| T2 | srd022-builtin-tools R0.1, R0.2, R0.3, R1.1, R1.2, R1.3, R1.4, R1.5, R1.6, R1.7, R1.8, R4.1, R4.2, R4.3, R4.4, R4.5, R4.6, R4.7 | builtin declarations follow the authoring standard and integrate with machines |
| T3 | srd023-cli-exec-tools R0.1, R0.2, R1.1, R1.2, R1.3, R1.4, R1.5, R1.6, R1.7, R1.8, R3.1, R3.2, R3.3, R3.4, R3.5, R3.6, R3.7, R3.8, R3.9, R5.1, R5.2, R5.3, R5.4 | exec declarations model one CLI verb with declared side effects |
| T4 | srd024-blocking-input-tools R0.1, R0.2, R0.3, R2.1, R2.2, R2.3, R2.4, R3.1, R3.2, R3.3, R3.4, R3.5 | boundary/blocking tools declare external input and signals |
| T5 | srd025-rollback-lifecycle R3.1, R3.2, R3.3, R3.4, R3.5, R7.1, R7.2, R7.3, R8.1, R8.2, R8.3, R8.4 | receipts, opt-in lifecycle behavior, reversibility, and safety behavior |
| T6 | sm-tool-language | ToolDef, ToolSpec, Command, and MachineSpec relationships |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The tool declaration includes problem, goals, requirements, non_goals, output, side_effects, reversibility, undo, errors, and relationships.

<!-- field: traces -->
#### Traces

- srd022-builtin-tools AC1

<!-- sc: S2 -->
### S2

The tool is classified correctly as a word, stateful_internal, boundary, or legacy sentence migration candidate.

<!-- field: traces -->
#### Traces

- srd023-cli-exec-tools AC5

<!-- sc: S3 -->
### S3

Irreversible or compensatable side effects are documented and require confirmation where they affect user data or shared infrastructure.

<!-- field: traces -->
#### Traces

- srd025-rollback-lifecycle AC6

<!-- sc: S4 -->
### S4

MachineSpec owns sequencing; the new tool implementation does not hide a multi-step workflow.

<!-- field: traces -->
#### Traces

- srd006-yaml-machine-config AC3

<!-- field: failure_scenarios -->
## Failure scenarios

| id | scenario | expected |
| --- | --- | --- |
| E1 | The tool description contains multiple observable verbs joined by "and". | The review rejects the ToolDef until it is split or classified as a documented boundary. |
| E2 | The tool writes files but omits side_effects or undo. | The declaration fails authoring review because rollback cannot reason about it. |
| E3 | The tool output requires LLM judgment to know whether it succeeded. | The developer adds structured output and declared signals, or classifies the tool as boundary. |

<!-- field: out_of_scope -->
## Out of scope

- Implementing the tool body.
- Designing new MachineSpec syntax.
- Completing contract metadata for every existing tool.

<!-- field: semantic_model -->
## Semantic model

sm-tool-language

<!-- field: test_suite_trace -->
## Test suite trace

test-rel02.0-tool-contracts
