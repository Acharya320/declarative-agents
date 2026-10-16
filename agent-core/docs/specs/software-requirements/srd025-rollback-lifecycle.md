---
id: srd025-rollback-lifecycle
title: Rollback, Checkpoint, Suspend, and Resume Lifecycle
depends_on:
  - srd_id: srd001-core-types
    symbols_used: [Command, Result, State, Signal]
  - srd_id: srd005-agentic-loop
    symbols_used: [Loop, LoopParams, RunResult, RunEvent]
  - srd_id: srd010-conversation
    symbols_used: [Conversation]
  - srd_id: srd022-builtin-tools
    symbols_used: [ToolDef]
  - srd_id: srd035-checkpoint-port
    symbols_used: [Checkpoint, Position, Execution, Entry, Receipt, AgentSnapshot]
  - srd_id: srd036-dolt-state-persistence
    symbols_used: [DoltCheckpoint, Revert]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Long-running agent programs need a safe way to record what happened, stop at
approval gates, resume later, and roll back from bad actions. Persistence and
resume rest on the srd035 checkpoint port: Save records the resumable position
and the ordered execution log; Load restores them. The srd036 Dolt adapter is
the default persistent backend, and NoopCheckpoint keeps disabled-mode
execution unchanged.

Rollback has two halves that this SRD keeps separate. DB-persisted state is
rewound by the Dolt adapter's git-style Revert(run_id, step_index) (srd036).
External effects -- files, external resources -- are reversed by a lifecycle
tool that walks the reverted execution entries in reverse and hands each
restored, receipt-bearing Result back to its tool's receipt-consuming Undo.
The receipt is the memento: the originating tool encodes it during Execute and
is the only decoder, so the engine and the checkpoint adapters never interpret
it. This SRD defines the requirements for that lifecycle -- checkpoint,
history, suspend, resume, and the two-part rollback -- without a three-layer
state model, a key/blob StateStore, a Workspace layer, a CheckpointPolicy, or
an engine-owned undo-walk.

<!-- field: goals -->
## Goals

- G1: Rest persistence and resume on the srd035 checkpoint port with the srd036 Dolt adapter as the default persistent backend.
- G2: Define the ordered execution history that history inspection and rollback traverse.
- G3: Make reversibility receipt-driven -- tools encode a receipt on Execute and consume it on Undo -- with the engine and adapters treating receipts as opaque.
- G4: Define rollback as Dolt Revert for DB-persisted state plus a lifecycle-tool reverse receipt walk for external effects.
- G5: Preserve suspend, approval gates, and a single resume contract.
- G6: Keep lifecycle behavior opt-in with zero persistence overhead when disabled.
- G7: Keep history and rollback as declarative, profile-first lifecycle tool actions.

<!-- req: R1 -->
## R1 - Checkpoint-port-backed persistence

<!-- req: R1.1 -->
### R1.1

Agent state, command/domain state, and conversation SHALL be persisted through the srd035 checkpoint port (Save(Position, Execution)), not a key/blob StateStore.

<!-- req: R1.2 -->
### R1.2

The default persistent backend SHALL be the srd036 Dolt adapter; NoopCheckpoint SHALL be the default when persistence is disabled.

<!-- req: R1.3 -->
### R1.3

Position SHALL carry the resumable machine position, folded conversation snapshot, and verified identity of the declarative program that registered the run's tools; Execution SHALL carry the ordered entries with their opaque receipts.

<!-- req: R1.4 -->
### R1.4

The runtime SHALL NOT model workspace or filesystem state as a separate persisted state layer; external effects are reversed through receipts rather than Workspace refs.

<!-- req: R2 -->
## R2 - Execution history

<!-- req: R2.1 -->
### R2.1

Each dispatched step SHALL append one Execution Entry (iteration, timestamp, command name, from state, to state, signal, result digest, and receipt) per srd035.

<!-- req: R2.2 -->
### R2.2

Execution SHALL preserve dispatch order for forward inspection and reverse traversal.

<!-- req: R2.3 -->
### R2.3

History inspection SHALL be expressible as a MachineSpec transition to a lifecycle tool such as checkpoint_history, which reads Dolt commit history (srd026).

<!-- req: R2.4 -->
### R2.4

History recording SHALL be absent when persistence is disabled (NoopCheckpoint).

<!-- req: R3 -->
## R3 - Receipt-driven reversibility

<!-- req: R3.1 -->
### R3.1

State-mutating commands SHALL encode an opaque Receipt during Execute that captures what a later Undo needs, such as prior file content, external resource ids, or a workspace commit hash.

<!-- req: R3.2 -->
### R3.2

Command Undo SHALL be receipt-consuming: it takes the restored receipt-bearing Result and reverses the external effect.

<!-- req: R3.3 -->
### R3.3

Read-only commands SHALL return an empty receipt and a no-op Undo.

<!-- req: R3.4 -->
### R3.4

The engine and the checkpoint adapters SHALL treat receipts as opaque and SHALL NOT interpret them; only the originating tool decodes its receipt.

<!-- req: R3.5 -->
### R3.5

The lifecycle validator SHALL check receipt presence for reversible state-mutating tools; receipt sufficiency SHALL be each tool's own round-trip test.

<!-- req: R4 -->
## R4 - Two-part rollback

<!-- req: R4.1 -->
### R4.1

Rollback of DB-persisted state SHALL be the git-style Revert(run_id, step_index) provided by the srd036 Dolt adapter.

<!-- req: R4.2 -->
### R4.2

Rollback of external effects SHALL be a reverse receipt walk that, after Revert, verifies and rebuilds the target run's declarative tool program and hands each reverted entry's restored Result to its originating tool's receipt-consuming Undo in reverse dispatch order.

<!-- req: R4.3 -->
### R4.3

The reverse receipt walk SHALL live in the lifecycle tool (checkpoint_rollback), not in the engine.

<!-- req: R4.4 -->
### R4.4

Irreversible entries SHALL be skipped by their declared reversibility tier and logged; compensatable entries whose Undo requires later operator or machine action SHALL be reported as pending compensation rather than as Undo failures.

<!-- req: R4.5 -->
### R4.5

Rollback SHALL return the machine position from which execution can resume.

<!-- req: R4.6 -->
### R4.6

Rollback SHALL NOT require machine.yaml changes; it is lifecycle behavior around the machine.

<!-- req: R5 -->
## R5 - Suspend and approval gates

<!-- req: R5.1 -->
### R5.1

The runtime SHALL support a suspended run status when lifecycle features are enabled.

<!-- req: R5.2 -->
### R5.2

A suspend or approval-gate tool SHALL emit a machine-declared signal that ends the run without failure.

<!-- req: R5.3 -->
### R5.3

Suspend SHALL persist a checkpoint through the srd035 port before returning StatusSuspended.

<!-- req: R5.4 -->
### R5.4

Approval gates SHALL be expressible as machine.yaml transitions to an internal tool.

<!-- req: R5.5 -->
### R5.5

Approval outcomes such as approved and rejected SHALL map to explicit machine signals on resume.

<!-- req: R6 -->
## R6 - Resume

<!-- req: R6.1 -->
### R6.1

Resume SHALL Load Position and Execution through the srd035 port, backed by the Dolt adapter by default.

<!-- req: R6.2 -->
### R6.2

Resume SHALL restore agent state and conversation from the typed snapshot and re-enter the same machine.

<!-- req: R6.3 -->
### R6.3

Resume SHALL be a single path with no restore-hook fan-out.

<!-- req: R6.4 -->
### R6.4

Resume SHALL validate the current machine and tool declarations before re-entering the loop; a separate rollback process SHALL verify the target run's persisted program digest before rebuilding Undo.

<!-- req: R6.5 -->
### R6.5

Resume errors SHALL distinguish missing checkpoint, incompatible checkpoint, and load failure.

<!-- req: R7 -->
## R7 - Opt-in and zero-overhead defaults

<!-- req: R7.1 -->
### R7.1

Rollback, checkpointing, suspend, resume, and approval-gate behavior SHALL be opt-in.

<!-- req: R7.2 -->
### R7.2

When persistence is disabled (NoopCheckpoint), existing machines SHALL run with unchanged behavior and no persistence.

<!-- req: R7.3 -->
### R7.3

Disabled lifecycle features SHALL not persist checkpoints or record history.

<!-- req: R8 -->
## R8 - Tool contract integration

<!-- req: R8.1 -->
### R8.1

Tool declarations SHOULD follow cfg-tool-authoring-standard for side_effects, reversibility, receipt-based undo, errors, and relationships.

<!-- req: R8.2 -->
### R8.2

Tools that mutate state SHOULD declare a receipt-encoding Execute and a receipt-consuming Undo.

<!-- req: R8.3 -->
### R8.3

Irreversible tools SHOULD require confirmation and declare why no receipt can reverse the effect.

<!-- req: R8.4 -->
### R8.4

Lifecycle tools SHALL declare emitted signals, typed config, output shape, side effects, reversibility, undo strategy, errors, and relationships.

<!-- req: R9 -->
## R9 - Declarative lifecycle tools and profiles

<!-- req: R9.1 -->
### R9.1

History and rollback SHALL be profile-first lifecycle machines that dispatch declared lifecycle tools such as checkpoint_history and checkpoint_rollback.

<!-- req: R9.2 -->
### R9.2

Persistence and Revert operations SHALL remain in the checkpoint port and the Dolt adapter and SHALL NOT depend on Cobra, ToolDef, agent profiles, or STL factory wiring.

<!-- req: R9.3 -->
### R9.3

The reverse receipt walk SHALL be owned by the checkpoint_rollback lifecycle tool, which hands restored Results to tool Undo.

<!-- req: R9.4 -->
### R9.4

Lifecycle tools SHALL NOT shell back into the agent binary to perform history or rollback.

<!-- field: non_goals -->
## Non-goals

- Retaining the three-layer state model, the key/blob StateStore, Workspace/GitWorkspace, or CheckpointPolicy.
- An engine-owned reverse undo-walk; the walk lives in the checkpoint_rollback lifecycle tool.
- Reverting files, HTTP, or workspace state inside the Dolt adapter; Revert is DB-only and tools own external compensation through receipts.
- Defining the Dolt schema or revert mechanics (srd036) or the checkpoint port wire format and Receipt shape (srd035).
- Making core.Loop concurrent.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

Persistence and resume rest on the srd035 checkpoint port with the srd036 Dolt default backend and NoopCheckpoint when disabled; no three-layer state model, StateStore, Workspace, or CheckpointPolicy remains.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R1.4
- R6.1
- R6.2
- R6.3
- R6.4
- R6.5
- R7.1
- R7.2
- R7.3

<!-- ac: AC2 -->
### AC2

Execution history is the ordered srd035 Entry log with receipts, preserving dispatch order and inspectable through checkpoint_history over Dolt commit history.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3
- R2.4

<!-- ac: AC3 -->
### AC3

Reversibility is receipt-driven -- tools encode receipts on Execute and consume them on Undo, read-only tools return empty receipts, the engine treats receipts as opaque, and the validator checks presence while tools test sufficiency.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3
- R3.4
- R3.5
- R8.1
- R8.2
- R8.3
- R8.4

<!-- ac: AC4 -->
### AC4

Rollback is Dolt Revert for DB-persisted state plus a lifecycle-tool reverse receipt walk for external effects, skipping declared irreversible entries, reporting pending compensation separately from genuine Undo failures, returning a resumable position, with no engine undo-walk and no machine.yaml changes.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
- R4.3
- R4.4
- R4.5
- R4.6
- R9.1
- R9.2
- R9.3
- R9.4

<!-- ac: AC5 -->
### AC5

Suspend and approval gates end the run through machine-declared signals with a checkpoint persisted via the srd035 port, and resume is a single Load-and-re-enter path.

<!-- field: traces -->
#### Traces

- R5.1
- R5.2
- R5.3
- R5.4
- R5.5

<!-- field: test_trace_plan -->
## Test trace plan

```yaml
- test_suite: test-rel02.0-lifecycle
  covers: [AC1, AC2, AC3, AC4, AC5]
- test_suite: test-rel02.0-tool-contracts
  covers: [AC3]
```

<!-- field: historical_background -->
## Historical background

- docs/guides/roll-backs.md
- docs/guides/tools-as-dsl.md

<!-- field: related_documents -->
## Related documents

- cfg-tool-authoring-standard
- cfg-runtime-contract
- srd035-checkpoint-port
- srd036-dolt-state-persistence
- srd026-lifecycle-tools

<!-- field: implementation_status -->
## Implementation status

```yaml
done:
    - Checkpoint, history, suspend, and resume use the srd035 port; cmd/agent wires DoltCheckpoint when --dolt-dsn is set and NoopCheckpoint otherwise.
    - checkpoint_rollback uses Dolt Revert plus a reverse receipt walk over restored Results.
    - State-mutating tools encode receipts for persisted rollback; StateStore, Workspace/GitWorkspace, CheckpointPolicy, and the engine undo-walk are retired implementation history.
    - Suspend and resume operate around explicit approval-gate machine signals.
```

<!-- field: diagrams -->
## Diagrams

````yaml
requirements_trace: |
    ```mermaid
    flowchart LR
      Problem["Problem: Safe record, suspend, resume, and rollback on the checkpoint port"] --> Goals["Goals: Rollback, Checkpoint, Suspend, and Resume Lifecycle"]
      Goals --> R1["R1: Checkpoint-port-backed persistence"]
      Goals --> R2["R2: Execution history"]
      Goals --> R3["R3: Receipt-driven reversibility"]
      Goals --> R4["R4: Two-part rollback"]
      Goals --> R5["R5: Suspend and approval gates"]
      Goals --> R6["R6: Resume"]
      Goals --> R7["R7: Opt-in and zero-overhead defaults"]
      Goals --> R8["R8: Tool contract integration"]
      Goals --> R9["R9: Declarative lifecycle tools and profiles"]
      AC1["AC1: Port-backed persistence and resume; no legacy layers"]
      R1 --> AC1
      R6 --> AC1
      R7 --> AC1
      AC2["AC2: Ordered execution history with receipts"]
      R2 --> AC2
      AC3["AC3: Receipt-driven reversibility"]
      R3 --> AC3
      R8 --> AC3
      AC4["AC4: Dolt Revert + lifecycle-tool receipt walk"]
      R4 --> AC4
      R9 --> AC4
      AC5["AC5: Suspend/approval gates and single resume"]
      R5 --> AC5
    ```
````
