---
id: rel02.0-uc002-history-rollback
title: User inspects history and rolls back after a bad tool action
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A user notices that an agent took a bad action, such as an incorrect edit or
a child process that generated unwanted artifacts. With persistence enabled
(the Dolt checkpoint backend), the user inspects the run's Dolt commit
history, chooses a target step, and asks the runtime to roll back.
checkpoint_rollback rewinds DB-persisted state with the Dolt adapter's
git-style Revert, then walks the reverted execution entries in reverse,
handing each restored, receipt-bearing Result to its tool's receipt-consuming
Undo so external effects are reversed too. A rolled-back file write leaves the
file at its prior content, and the run returns a resumable position.

<!-- field: actor -->
## Actor

Developer reviewing an agent run

<!-- field: trigger -->
## Trigger

The developer finds an undesirable result in the workspace or conversation
after a run whose persistence used the default Dolt checkpoint backend.

<!-- field: flow -->
## Flow

- F1: During execution, the loop appends one Execution Entry per dispatched command and calls Save(Position, Execution) after each cycle through the srd035 checkpoint port; the Dolt backend commits one step on the run branch (srd035 R1, srd036 R4).
- F2: State-mutating commands encode an opaque receipt during Execute that captures what a later Undo needs, such as prior file content or a resource id; read-only commands return an empty receipt (srd025 R3, cfg-tool-authoring-standard).
- F3: The developer inspects the run's Dolt commit history through the checkpoint_history lifecycle adapter when a profile selects it, reading step index, command, from/to state, signal, and receipt presence (srd026 R2).
- F4: The developer selects the last known-good step as the rollback target.
- F5: checkpoint_rollback calls the Dolt adapter Revert(run_id, step_index) to rewind DB-persisted state (Position, Execution, and the folded conversation) to the target step (srd036 R6, srd025 R4).
- F6: checkpoint_rollback walks the reverted entries in reverse dispatch order and hands each restored, receipt-bearing Result to its tool's receipt-consuming Undo to reverse external effects such as file writes (srd025 R4, srd026 R3).
- F7: Irreversible entries are skipped and logged, and the rollback report names them rather than claiming a clean rollback (srd025 R4, srd026 R3).
- F8: Rollback returns the resumable machine position; on retry the run re-enters with a clean conversation and a clean worktree, restored through receipts rather than a workspace layer (srd025 R4, R6).

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd035-checkpoint-port](../software-requirements/srd035-checkpoint-port.md) [R1.1](../software-requirements/srd035-checkpoint-port.md#r11), [R1.2](../software-requirements/srd035-checkpoint-port.md#r12), [R1.3](../software-requirements/srd035-checkpoint-port.md#r13), [R1.4](../software-requirements/srd035-checkpoint-port.md#r14), [R1.5](../software-requirements/srd035-checkpoint-port.md#r15), [R3.1](../software-requirements/srd035-checkpoint-port.md#r31), [R3.2](../software-requirements/srd035-checkpoint-port.md#r32), [R3.3](../software-requirements/srd035-checkpoint-port.md#r33), [R3.4](../software-requirements/srd035-checkpoint-port.md#r34), [R3.5](../software-requirements/srd035-checkpoint-port.md#r35), [R3.6](../software-requirements/srd035-checkpoint-port.md#r36), [R3.7](../software-requirements/srd035-checkpoint-port.md#r37), [R3.8](../software-requirements/srd035-checkpoint-port.md#r38), [R3.9](../software-requirements/srd035-checkpoint-port.md#r39), [R3.10](../software-requirements/srd035-checkpoint-port.md#r310) | checkpoint port Save/Load and the opaque receipt on each entry |
| T2 | [srd036-dolt-state-persistence](../software-requirements/srd036-dolt-state-persistence.md) [R4.1](../software-requirements/srd036-dolt-state-persistence.md#r41), [R4.2](../software-requirements/srd036-dolt-state-persistence.md#r42), [R4.3](../software-requirements/srd036-dolt-state-persistence.md#r43), [R4.4](../software-requirements/srd036-dolt-state-persistence.md#r44), [R4.5](../software-requirements/srd036-dolt-state-persistence.md#r45), [R4.6](../software-requirements/srd036-dolt-state-persistence.md#r46), [R6.1](../software-requirements/srd036-dolt-state-persistence.md#r61), [R6.2](../software-requirements/srd036-dolt-state-persistence.md#r62), [R6.3](../software-requirements/srd036-dolt-state-persistence.md#r63), [R6.4](../software-requirements/srd036-dolt-state-persistence.md#r64), [R6.5](../software-requirements/srd036-dolt-state-persistence.md#r65), [R6.6](../software-requirements/srd036-dolt-state-persistence.md#r66) | commit-per-step, branch-per-run, and git-style Revert |
| T3 | [srd025-rollback-lifecycle](../software-requirements/srd025-rollback-lifecycle.md) [R2.1](../software-requirements/srd025-rollback-lifecycle.md#r21), [R2.2](../software-requirements/srd025-rollback-lifecycle.md#r22), [R2.3](../software-requirements/srd025-rollback-lifecycle.md#r23), [R2.4](../software-requirements/srd025-rollback-lifecycle.md#r24), [R3.1](../software-requirements/srd025-rollback-lifecycle.md#r31), [R3.2](../software-requirements/srd025-rollback-lifecycle.md#r32), [R3.3](../software-requirements/srd025-rollback-lifecycle.md#r33), [R3.4](../software-requirements/srd025-rollback-lifecycle.md#r34), [R3.5](../software-requirements/srd025-rollback-lifecycle.md#r35), [R4.1](../software-requirements/srd025-rollback-lifecycle.md#r41), [R4.2](../software-requirements/srd025-rollback-lifecycle.md#r42), [R4.3](../software-requirements/srd025-rollback-lifecycle.md#r43), [R4.4](../software-requirements/srd025-rollback-lifecycle.md#r44), [R4.5](../software-requirements/srd025-rollback-lifecycle.md#r45), [R4.6](../software-requirements/srd025-rollback-lifecycle.md#r46) | execution history, receipt-driven reversibility, and two-part rollback |
| T4 | [sm-rollback-lifecycle](../semantic-models/rollback-lifecycle.md) | Dolt commit history plus receipts as the reverse edges for external effects |
| T5 | [sm-command-undo-audit](../semantic-models/command-undo-audit.md) | receipt-consuming Undo inventory across commands |
| T6 | [srd026-lifecycle-tools](../software-requirements/srd026-lifecycle-tools.md) [R2.1](../software-requirements/srd026-lifecycle-tools.md#r21), [R2.2](../software-requirements/srd026-lifecycle-tools.md#r22), [R2.3](../software-requirements/srd026-lifecycle-tools.md#r23), [R2.4](../software-requirements/srd026-lifecycle-tools.md#r24), [R2.5](../software-requirements/srd026-lifecycle-tools.md#r25), [R2.6](../software-requirements/srd026-lifecycle-tools.md#r26), [R2.7](../software-requirements/srd026-lifecycle-tools.md#r27), [R2.8](../software-requirements/srd026-lifecycle-tools.md#r28), [R2.9](../software-requirements/srd026-lifecycle-tools.md#r29), [R3.1](../software-requirements/srd026-lifecycle-tools.md#r31), [R3.2](../software-requirements/srd026-lifecycle-tools.md#r32), [R3.3](../software-requirements/srd026-lifecycle-tools.md#r33), [R3.4](../software-requirements/srd026-lifecycle-tools.md#r34), [R3.5](../software-requirements/srd026-lifecycle-tools.md#r35), [R3.6](../software-requirements/srd026-lifecycle-tools.md#r36), [R3.7](../software-requirements/srd026-lifecycle-tools.md#r37), [R3.8](../software-requirements/srd026-lifecycle-tools.md#r38), [R3.9](../software-requirements/srd026-lifecycle-tools.md#r39), [R3.10](../software-requirements/srd026-lifecycle-tools.md#r310), [R3.11](../software-requirements/srd026-lifecycle-tools.md#r311), [R3.12](../software-requirements/srd026-lifecycle-tools.md#r312) | checkpoint_history and checkpoint_rollback adapters |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Dolt commit history exposes enough per step (step index, command, signal, receipt presence) to choose a rollback target without inspecting private command objects.

<!-- field: traces -->
#### Traces

- [srd025-rollback-lifecycle](../software-requirements/srd025-rollback-lifecycle.md) [AC2](../software-requirements/srd025-rollback-lifecycle.md#ac2)
- [srd026-lifecycle-tools](../software-requirements/srd026-lifecycle-tools.md) [AC2](../software-requirements/srd026-lifecycle-tools.md#ac2)

<!-- sc: S2 -->
### S2

Rollback rewinds DB-persisted state (Position, Execution, conversation) with the Dolt adapter Revert to the target step.

<!-- field: traces -->
#### Traces

- [srd036-dolt-state-persistence](../software-requirements/srd036-dolt-state-persistence.md) [AC6](../software-requirements/srd036-dolt-state-persistence.md#ac6)
- [srd025-rollback-lifecycle](../software-requirements/srd025-rollback-lifecycle.md) [AC4](../software-requirements/srd025-rollback-lifecycle.md#ac4)

<!-- sc: S3 -->
### S3

A rolled-back file write leaves the file at its prior content, restored by a fresh command instance decoding its persisted receipt; the retry sees a clean context and a clean worktree through receipts rather than a workspace layer.

<!-- field: traces -->
#### Traces

- [srd025-rollback-lifecycle](../software-requirements/srd025-rollback-lifecycle.md) [AC4](../software-requirements/srd025-rollback-lifecycle.md#ac4)
- [srd035-checkpoint-port](../software-requirements/srd035-checkpoint-port.md) [AC3](../software-requirements/srd035-checkpoint-port.md#ac3)

<!-- sc: S4 -->
### S4

Irreversible entries are skipped and named, and partial failures identify which reverted entry's receipt-consuming Undo failed.

<!-- field: traces -->
#### Traces

- [srd025-rollback-lifecycle](../software-requirements/srd025-rollback-lifecycle.md) [AC4](../software-requirements/srd025-rollback-lifecycle.md#ac4)
- [srd026-lifecycle-tools](../software-requirements/srd026-lifecycle-tools.md) [AC3](../software-requirements/srd026-lifecycle-tools.md#ac3)

<!-- field: failure_scenarios -->
## Failure scenarios

| id | scenario | expected |
| --- | --- | --- |
| E1 | A traversed entry's tool returns an empty receipt but mutated external state. | The rollback report identifies the missing receipt as a residual risk; the validator's receipt-presence check for reversible tools should have caught it during authoring. |
| E2 | A receipt is present but insufficient to fully reverse the effect. | The tool's own round-trip test is the guard; rollback reports the receipt-consuming Undo failure and the state left behind. |
| E3 | A tool is irreversible or compensatable rather than reversible. | Rollback skips and names the entry, reporting the compensation or confirmation requirement declared in the ToolDef. |

<!-- field: out_of_scope -->
## Out of scope

- Replaying a full run forward from scratch.
- Designing the final rollback CLI syntax.
- Reverting external systems that only support manual compensation.
- Reverting files or HTTP effects inside the Dolt adapter; Revert is DB-only and tools own external compensation through receipts.

<!-- field: semantic_model -->
## Semantic model

sm-rollback-lifecycle

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel02.0-lifecycle](../test-suites/test-rel02.0-lifecycle.md)
