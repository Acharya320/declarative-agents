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
| T1 | srd035-checkpoint-port R1.1, R1.2, R1.3, R1.4, R1.5, R3.1, R3.2, R3.3, R3.4, R3.5, R3.6, R3.7, R3.8, R3.9, R3.10 | checkpoint port Save/Load and the opaque receipt on each entry |
| T2 | srd036-dolt-state-persistence R4.1, R4.2, R4.3, R4.4, R4.5, R4.6, R6.1, R6.2, R6.3, R6.4, R6.5, R6.6 | commit-per-step, branch-per-run, and git-style Revert |
| T3 | srd025-rollback-lifecycle R2.1, R2.2, R2.3, R2.4, R3.1, R3.2, R3.3, R3.4, R3.5, R4.1, R4.2, R4.3, R4.4, R4.5, R4.6 | execution history, receipt-driven reversibility, and two-part rollback |
| T4 | sm-rollback-lifecycle | Dolt commit history plus receipts as the reverse edges for external effects |
| T5 | sm-command-undo-audit | receipt-consuming Undo inventory across commands |
| T6 | srd026-lifecycle-tools R2.1, R2.2, R2.3, R2.4, R2.5, R2.6, R2.7, R2.8, R2.9, R3.1, R3.2, R3.3, R3.4, R3.5, R3.6, R3.7, R3.8, R3.9, R3.10, R3.11, R3.12 | checkpoint_history and checkpoint_rollback adapters |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Dolt commit history exposes enough per step (step index, command, signal, receipt presence) to choose a rollback target without inspecting private command objects.

<!-- field: traces -->
#### Traces

- srd025-rollback-lifecycle AC2
- srd026-lifecycle-tools AC2

<!-- sc: S2 -->
### S2

Rollback rewinds DB-persisted state (Position, Execution, conversation) with the Dolt adapter Revert to the target step.

<!-- field: traces -->
#### Traces

- srd036-dolt-state-persistence AC6
- srd025-rollback-lifecycle AC4

<!-- sc: S3 -->
### S3

A rolled-back file write leaves the file at its prior content, restored by a fresh command instance decoding its persisted receipt; the retry sees a clean context and a clean worktree through receipts rather than a workspace layer.

<!-- field: traces -->
#### Traces

- srd025-rollback-lifecycle AC4
- srd035-checkpoint-port AC3

<!-- sc: S4 -->
### S4

Irreversible entries are skipped and named, and partial failures identify which reverted entry's receipt-consuming Undo failed.

<!-- field: traces -->
#### Traces

- srd025-rollback-lifecycle AC4
- srd026-lifecycle-tools AC3

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

test-rel02.0-lifecycle
