---
id: test-rel02.0-lifecycle
title: Test Suite for Release 02.0 Checkpoint Port and Dolt Persistence
release: "02.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel02.0-uc001-approval-suspend-resume
- rel02.0-uc002-history-rollback

<!-- field: preconditions -->
## Preconditions

- Persistence defaults to NoopCheckpoint so disabled-mode runs keep current loop behavior with no store configured.
- InMemoryCheckpoint is the reference adapter for in-process round-trip scenarios.
- A fake database/sql seam (fakeDB) stands in for Dolt so branch-per-run, commit-per-step, merge-on-terminal, and Revert are asserted without a live server.
- A local dolt binary enables gated end-to-end persistence tests whose harness starts an ephemeral sql-server; they skip when no binary is available.
- Test commands can encode opaque receipts during Execute and reverse them through a receipt-consuming Undo.

<!-- field: case_groups -->
## Case groups

- rel02.0-uc001: suspend and resume through the typed checkpoint port
- rel02.0-uc002: Dolt persistence, history, and git-style revert

<!-- field: test_cases -->
## Test cases

<!-- case: Disabled persistence defaults to NoopCheckpoint and persists nothing | status: done -->
### Disabled persistence defaults to NoopCheckpoint and persists nothing

Run a machine without a configured Checkpoint. The loop defaults to
NoopCheckpoint, whose Save is a no-op and whose Load reports
ErrNoCheckpoint, preserving disabled-mode behavior with no persistence
overhead.

<!-- field: use_case -->
#### Use case

rel02.0-uc001-approval-suspend-resume

<!-- field: go_test -->
#### Go test

TestLoop_NoopCheckpointDefaultPersistsNothing, TestNoopCheckpointSaveIsNoopAndLoadReportsNotFound

<!-- field: inputs | form: "tree" -->
#### Inputs

- checkpoint: nil
- machine: simple_success

<!-- field: expected | form: "tree" -->
#### Expected

- status: succeeded
- save: noop
- load_error: ErrNoCheckpoint

<!-- field: traces -->
#### Traces

- srd025-rollback-lifecycle AC1
- srd035-checkpoint-port AC5
- srd035-checkpoint-port AC6

<!-- case: Loop saves Position and Execution after each dispatch cycle | status: done -->
### Loop saves Position and Execution after each dispatch cycle

Run a machine with an InMemoryCheckpoint configured. The loop calls
Save after each dispatch, appending one Execution entry per cycle and
updating the resumable Position, so a later Load returns the final
position and the full ordered log.

<!-- field: use_case -->
#### Use case

rel02.0-uc001-approval-suspend-resume

<!-- field: go_test -->
#### Go test

TestLoop_SavesSnapshotAfterDispatchWithConfiguredAdapter

<!-- field: inputs | form: "tree" -->
#### Inputs

- checkpoint: InMemoryCheckpoint
- machine: simple_success

<!-- field: expected | form: "tree" -->
#### Expected

- save_timing: after_each_dispatch
- execution_len: iteration_count
- position_matches_last_dispatch: true

<!-- field: traces -->
#### Traces

- srd025-rollback-lifecycle AC1
- srd025-rollback-lifecycle AC2
- srd035-checkpoint-port AC6

<!-- case: Conversation folds into the typed snapshot on Save | status: done -->
### Conversation folds into the typed snapshot on Save

With a SnapshotConversation hook supplying the domain-owned conversation,
the loop folds it into AgentSnapshot and persists it through the port.
InMemoryCheckpoint round-trips the folded conversation and each entry's
opaque receipt verbatim, restoring an empty receipt for read-only steps.

<!-- field: use_case -->
#### Use case

rel02.0-uc001-approval-suspend-resume

<!-- field: go_test -->
#### Go test

TestLoop_PortSavePersistsConversation, TestInMemoryCheckpointRoundTripsConversationAndReceipts

<!-- field: inputs | form: "tree" -->
#### Inputs

- checkpoint: InMemoryCheckpoint
- conversation: user_message
- receipts:
  1. write_receipt
  1. empty

<!-- field: expected | form: "tree" -->
#### Expected

- conversation_restored: true
- receipt_round_trip: verbatim
- empty_receipt_restored_empty: true

<!-- field: traces -->
#### Traces

- srd025-rollback-lifecycle AC1
- srd025-rollback-lifecycle AC3
- srd035-checkpoint-port AC3
- srd035-checkpoint-port AC4
- srd035-checkpoint-port AC5

<!-- case: Position and Execution JSON wire format | status: done -->
### Position and Execution JSON wire format

Position serializes as an object with current_state, last_signal, and a
nested snapshot; an empty Execution serializes as JSON null; the opaque
receipt is omitted when empty and present as a JSON string when set. The
typed values compile and round-trip through JSON.

<!-- field: use_case -->
#### Use case

rel02.0-uc001-approval-suspend-resume

<!-- field: go_test -->
#### Go test

TestCheckpointWireFormat, TestInMemoryCheckpointRoundTripsConversationAndReceipts

<!-- field: inputs | form: "tree" -->
#### Inputs

- value:
  1. Position
  1. Execution
  1. Entry

<!-- field: expected | form: "tree" -->
#### Expected

- position_fields:
  1. current_state
  1. last_signal
  1. snapshot
- empty_execution: "null"
- receipt_omitted_when_empty: true
- receipt_present_as_json_string: true

<!-- field: traces -->
#### Traces

- srd035-checkpoint-port AC2
- srd035-checkpoint-port AC3

<!-- case: Conversation references survive restart and stay unavailable without persistence | status: done -->
### Conversation references survive restart and stay unavailable without persistence

A persistent checkpoint exposes the latest completed run/step/revision
reference and resolves its authoritative AgentSnapshot.Conversation
through a fresh adapter or command. Parsing rejects cross-run and
mismatched step/revision identities, including against a real gated Dolt
server, and Revert restores the target reference. cmd/agent supplies the loop provider and
lifecycle-target resolver to invoke_llm and reset_history. Noop and
request-local execution expose neither a usable reference nor synthetic
telemetry identity. Parallel for_each rejects serial-only commands before
execution or conversation mutation, while sequential dispatch remains supported.

<!-- field: use_case -->
#### Use case

rel02.0-uc001-approval-suspend-resume

<!-- field: go_test -->
#### Go test

TestDoltCheckpointConversationReferencesSurviveFreshAdapter, TestDoltCheckpointConversationReferenceRejectsWrongRunAndStep, TestDoltCheckpointConversationReferenceRejectsLaterRevisionClaimingEarlierStep, TestDoltCheckpointConversationReferencesAgainstRealDolt, TestNoopCheckpointReportsConversationReferencesUnavailable, TestResetHistoryFactoryRestartsFromCheckpointReference, TestInvokeLLMRestartsThroughCompositionCheckpointPorts, TestRequestLocalStateDoesNotExposeHostCheckpointReferences, TestLoopParallelForEachRejectsSerialDispatchBeforeExecution, TestLoopSequentialForEachAllowsSerialDispatch

<!-- field: inputs | form: "tree" -->
#### Inputs

- backends:
  1. DoltCheckpoint
  1. InMemoryCheckpoint
  1. NoopCheckpoint
- restart_boundary:
  1. fresh_adapter
  1. fresh_command
- request_machine_persistence: omitted

<!-- field: expected | form: "tree" -->
#### Expected

- historical_conversation_resolved: true
- run_isolation: strict
- restart_undo:
  1. invoke_llm
  1. reset_history
- disabled_reference_available: false
- request_machine_reference_available: false

<!-- field: traces -->
#### Traces

- srd035-checkpoint-port AC3
- srd035-checkpoint-port AC5
- srd035-checkpoint-port AC8
- srd036-dolt-state-persistence AC5

<!-- case: Ten-turn trace-plus-receipt content growth is linear | status: done -->
### Ten-turn trace-plus-receipt content growth is linear

Run ten zero-sleep invoke_llm turns with fixed-size prompts and responses.
Sum serialized trace content attribute values and v2 receipt bytes. The
ten-turn total stays within twice the five-turn total plus a constant,
remains below the quadratic doubling bound, and is smaller than the old
full-history-receipt baseline. Versioned AgentSnapshot.Conversation
storage is authoritative and excluded from the metric.

<!-- field: use_case -->
#### Use case

rel02.0-uc001-approval-suspend-resume

<!-- field: go_test -->
#### Go test

TestTenTurnDeltaTraceAndReceiptVolumeGrowsLinearly

<!-- field: inputs | form: "tree" -->
#### Inputs

- turns: 10
- telemetry_capture: delta
- sleep: zero
- excluded_storage: AgentSnapshot.Conversation

<!-- field: expected | form: "tree" -->
#### Expected

- growth_bound: linear
- below_quadratic_doubling: true
- below_full_history_receipt_baseline: true
- receipts_contain_conversation_text: false

<!-- field: traces -->
#### Traces

- srd035-checkpoint-port AC3
- srd035-checkpoint-port AC8

<!-- case: Resume re-enters the loop from the typed port and reports not-found | status: done -->
### Resume re-enters the loop from the typed port and reports not-found

A run suspended at an approval gate is resumed purely through the typed
Checkpoint port: Load restores the Position, Execution, and folded
conversation, and the loop re-enters at the checkpointed state carrying
the persisted counters forward. The resume signal defaults to Approved
but honors an explicit override; a not-found snapshot surfaces
ErrNoCheckpoint through the single resume path.

<!-- field: use_case -->
#### Use case

rel02.0-uc001-approval-suspend-resume

<!-- field: go_test -->
#### Go test

TestResumeReentersLoopFromTypedPort, TestLoadResumeExposesTypedSnapshotForDomainRestore, TestResumeHonorsExplicitResumeSignal, TestResumeReportsMissingCheckpoint

<!-- field: inputs | form: "tree" -->
#### Inputs

- checkpoint: InMemoryCheckpoint
- suspended_state: AwaitingApproval
- resume_signal:
  1. Approved
  1. Rejected

<!-- field: expected | form: "tree" -->
#### Expected

- resumed_status: succeeded
- counters_carried_forward: true
- explicit_signal_preserved: true
- missing_checkpoint_error: ErrNoCheckpoint

<!-- field: traces -->
#### Traces

- srd025-rollback-lifecycle AC1
- srd035-checkpoint-port AC4
- srd035-checkpoint-port AC6

<!-- case: Approval profile suspends and resolves the checkpoint backend at the CLI | status: done -->
### Approval profile suspends and resolves the checkpoint backend at the CLI

The approval lifecycle profile routes through a suspend tool and exits
as suspended. The composition root resolves NoopCheckpoint by default
and the Dolt adapter when --dolt-dsn is set, surfacing an unparseable
DSN as a typed ErrDolt; resume without a persistent backend reports
ErrNoCheckpoint.

<!-- field: use_case -->
#### Use case

rel02.0-uc001-approval-suspend-resume

<!-- field: go_test -->
#### Go test

TestApprovalLifecycleProfileSuspendsThroughCheckpointPort, TestOpenDefaultsToNoop, TestOpenWithDoltDSNOpensDoltBackend, TestResumeWithoutPersistentBackendReportsNoCheckpoint

<!-- field: inputs | form: "tree" -->
#### Inputs

- profile: lifecycle/approval/profile.yaml
- dolt_dsn:
  1. unset
  1. invalid

<!-- field: expected | form: "tree" -->
#### Expected

- first_run_terminal_state: suspended
- default_backend: NoopCheckpoint
- dolt_dsn_error: ErrDolt
- resume_without_backend_error: ErrNoCheckpoint

<!-- field: traces -->
#### Traces

- srd025-rollback-lifecycle AC1
- srd018-cli-flag-contract AC5
- srd024-blocking-input-tools AC2
- srd035-checkpoint-port AC5
- srd036-dolt-state-persistence AC1

<!-- case: End-to-end suspend and resume across processes through the Dolt backend | status: done -->
### End-to-end suspend and resume across processes through the Dolt backend

Against a reachable local dolt sql-server, start one composition-root
child process to save a suspended run and let that process exit. Start a
second child process to load the run and emit an artifact proving the
Position, folded conversation, and Execution (including receipts) match
the pre-suspend state. Skipped when no dolt sql-server is reachable.

<!-- field: use_case -->
#### Use case

rel02.0-uc001-approval-suspend-resume

<!-- field: integration_target -->
#### Integration target

mage integration:dolt

<!-- field: go_test -->
#### Go test

TestDoltCheckpointSuspendResumeAcrossProcesses

<!-- field: inputs | form: "tree" -->
#### Inputs

- backend: DoltCheckpoint
- dsn: root@tcp(127.0.0.1:3306)/agent_core_test

<!-- field: expected | form: "tree" -->
#### Expected

- fresh_process_load: equivalent_state
- skipped_without_dsn: true

<!-- field: traces -->
#### Traces

- srd036-dolt-state-persistence AC1
- srd036-dolt-state-persistence AC5

<!-- case: DoltCheckpoint implements the port behind the database/sql seam | status: done -->
### DoltCheckpoint implements the port behind the database/sql seam

DoltCheckpoint satisfies the srd035 Checkpoint port over the database/sql
Database/Transaction/Scanner/Rows seam, with no Dolt import in core. Load
reports ErrNoCheckpoint when the run branch or its rows do not exist.

<!-- field: use_case -->
#### Use case

rel02.0-uc002-history-rollback

<!-- field: go_test -->
#### Go test

TestDoltCheckpointImplementsPort, TestDoltCheckpointLoadNotFound

<!-- field: inputs | form: "tree" -->
#### Inputs

- backend: DoltCheckpoint
- seam: fakeDB

<!-- field: expected | form: "tree" -->
#### Expected

- implements_checkpoint_port: true
- missing_run_load_error: ErrNoCheckpoint

<!-- field: traces -->
#### Traces

- srd036-dolt-state-persistence AC1
- srd036-dolt-state-persistence AC5

<!-- case: Save reconstructs Position and Execution with receipts from the run branch | status: done -->
### Save reconstructs Position and Execution with receipts from the run branch

Save a run's Position and ordered Execution across two steps, then Load
reconstructs both from the latest commit on the run branch, including the
folded conversation and each entry's opaque receipt. A step with no
receipt stores SQL NULL and restores an empty receipt.

<!-- field: use_case -->
#### Use case

rel02.0-uc002-history-rollback

<!-- field: go_test -->
#### Go test

TestDoltCheckpointSaveLoadRoundTrip

<!-- field: inputs | form: "tree" -->
#### Inputs

- backend: DoltCheckpoint
- steps: 2
- receipts:
  1. file_receipt
  1. empty

<!-- field: expected | form: "tree" -->
#### Expected

- position_reconstructed: true
- execution_order_preserved: true
- receipt_round_trip: verbatim
- null_receipt_restored_empty: true

<!-- field: traces -->
#### Traces

- srd036-dolt-state-persistence AC3
- srd036-dolt-state-persistence AC2
- srd036-dolt-state-persistence AC5

<!-- case: Commit-per-step, branch-per-run, and merge-to-main on terminal state | status: done -->
### Commit-per-step, branch-per-run, and merge-to-main on terminal state

Each run executes on its own branch created once; each Save creates one
commit with step_index and signal in the message. A non-terminal Save
does not merge, while a terminal-state Save merges the run branch to main
and deletes the run branch.

<!-- field: use_case -->
#### Use case

rel02.0-uc002-history-rollback

<!-- field: go_test -->
#### Go test

TestDoltCheckpointCommitPerStepAndBranchPerRun, TestDoltCheckpointMergeOnTerminalState

<!-- field: inputs | form: "tree" -->
#### Inputs

- backend: DoltCheckpoint
- steps: 2
- terminal_state: Done

<!-- field: expected | form: "tree" -->
#### Expected

- branch_created: once_per_run
- commits: one_per_step
- terminal_merge_to_main: true
- run_branch_deleted: true

<!-- field: traces -->
#### Traces

- srd036-dolt-state-persistence AC4

<!-- case: Revert resets the run branch to the target step commit | status: done -->
### Revert resets the run branch to the target step commit

Revert(run_id, step_index) resets the run branch to the commit recorded
at that step so a subsequent Load returns the Position and Execution as of
the target step. An unresolved step reports ErrRevertUnresolved.

<!-- field: use_case -->
#### Use case

rel02.0-uc002-history-rollback

<!-- field: go_test -->
#### Go test

TestDoltCheckpointRevertResetsToStepCommit, TestDoltCheckpointRevertLeavesReferenceUnavailableWithoutTargetConversation, TestDoltCheckpointRevertLeavesReferenceUnavailableForDifferentAdapterRun, TestDoltCheckpointRevertUnresolvedStep

<!-- field: inputs | form: "tree" -->
#### Inputs

- backend: DoltCheckpoint
- steps: 2
- revert_to_step:
  1. 0
  1. 99

<!-- field: expected | form: "tree" -->
#### Expected

- branch_reset_to_target: true
- load_after_revert: target_step_state
- unresolved_step_error: ErrRevertUnresolved

<!-- field: traces -->
#### Traces

- srd036-dolt-state-persistence AC6

<!-- case: checkpoint_history formats the persisted execution log | status: done -->
### checkpoint_history formats the persisted execution log

checkpoint_history loads the persisted Execution through the port and
formats each step's index, command, from/to state, signal, and receipt
presence. It reports CommandError without a Checkpoint backend and
ErrNoCheckpoint when nothing is persisted; its Undo is a no-op.

<!-- field: use_case -->
#### Use case

rel02.0-uc002-history-rollback

<!-- field: go_test -->
#### Go test

TestCheckpointHistoryExecuteFormatsExecutionLog, TestCheckpointHistoryExecuteRequiresCheckpoint, TestCheckpointHistoryExecuteReportsNoCheckpoint, TestCheckpointHistoryUndoIsNoop

<!-- field: inputs | form: "tree" -->
#### Inputs

- tool: checkpoint_history
- checkpoint:
  1. InMemoryCheckpoint
  1. nil
  1. empty

<!-- field: expected | form: "tree" -->
#### Expected

- formatted_history: true
- missing_backend_error: CommandError
- no_checkpoint_error: ErrNoCheckpoint
- undo: noop

<!-- field: traces -->
#### Traces

- srd025-rollback-lifecycle AC2
- srd026-lifecycle-tools AC6
- srd026-lifecycle-tools AC2
- srd035-checkpoint-port AC1

<!-- case: checkpoint_rollback reverts DB state and validates its inputs | status: done -->
### checkpoint_rollback reverts DB state and validates its inputs

checkpoint_rollback resolves the target iteration to a step, calls the
Dolt adapter Revert to rewind DB-persisted state, and returns a summary
naming the run and target. It requires a revertible Checkpoint backend
and a target iteration, reports a missing target or reverse path after DB
Revert, and its Undo requests compensation rather than silently reversing.

<!-- field: use_case -->
#### Use case

rel02.0-uc002-history-rollback

<!-- field: go_test -->
#### Go test

TestCheckpointRollbackReportsMissingUndoAfterDBRevert, TestCheckpointRollbackRequiresRevertibleCheckpoint, TestCheckpointRollbackRequiresTargetIteration, TestCheckpointRollbackReportsMissingTargetIteration, TestCheckpointRollbackUndoRequestsCompensation

<!-- field: inputs | form: "tree" -->
#### Inputs

- tool: checkpoint_rollback
- checkpoint: fakeReverter
- to_iteration:
  1. present
  1. missing
  1. unresolved

<!-- field: expected | form: "tree" -->
#### Expected

- revert_delegated_to_adapter: true
- load_after_rollback: target_step_state
- missing_backend_error: requires_revertible_checkpoint
- missing_target_error: requires_to_iteration
- unresolved_target_error: target_not_found
- undo: requests_compensation

<!-- field: traces -->
#### Traces

- srd036-dolt-state-persistence AC6
- srd026-lifecycle-tools AC4
- srd026-lifecycle-tools AC6
- srd026-lifecycle-tools AC3

<!-- case: Rollback restores a file through the persisted opaque receipt | status: done -->
### Rollback restores a file through the persisted opaque receipt

A real write overwrites a file and records its opaque receipt in the
persisted Execution. After Revert rewinds DB state, checkpoint_rollback
walks the reverted entries in reverse and hands each restored,
receipt-bearing Result to a fresh command instance resolved through the
registry, whose receipt-consuming Undo restores the file to its prior
content.

<!-- field: use_case -->
#### Use case

rel02.0-uc002-history-rollback

<!-- field: go_test -->
#### Go test

TestCheckpointRollbackRestoresFileViaPersistedReceipt

<!-- field: inputs | form: "tree" -->
#### Inputs

- tool: checkpoint_rollback
- command: write
- receipt: file_before_content
- registry: resolves_write_to_reverser

<!-- field: expected | form: "tree" -->
#### Expected

- db_state_reverted: true
- file_restored_to_prior_content: true
- external_effect_reversed_via_receipt: true

<!-- field: traces -->
#### Traces

- srd025-rollback-lifecycle AC3
- srd025-rollback-lifecycle AC4
- srd035-checkpoint-port AC3
- srd036-dolt-state-persistence AC6

<!-- case: Reconstructed registry proves representative rollback families | status: done -->
### Reconstructed registry proves representative rollback families

Static family validation rejects a selected mutating declaration when its
reconstructed builder or fresh Reverser is absent, nil, or loses the
declaration alias. A mixed-family tracer verifies the persisted program
digest, rebuilds builders from declarations, and walks real LLM retry,
validation, SelfInvoke, suspend, checkpoint, exit, and REST receipts.
Receiptless no-op entries succeed, irreversible entries are explicit
skips, named compensation stays pending, malformed receipts remain named
partial failures, and unrelated entries continue reversing.

<!-- field: use_case -->
#### Use case

rel02.0-uc002-history-rollback

<!-- field: go_test -->
#### Go test

TestValidateReceiptFamiliesGuardsSelectedMutatorPlumbing, TestFreshRollbackRegistryTracesRepresentativeFamilies

<!-- field: inputs | form: "tree" -->
#### Inputs

- restart_boundary: reconstructed_registry
- families:
  1. llm
  1. validation
  1. child
  1. lifecycle
  1. rest
- controls:
  1. receiptless_noop
  1. irreversible
  1. malformed_receipt

<!-- field: expected | form: "tree" -->
#### Expected

- program_digest_verified: true
- aliases_preserved: true
- receipt_backed_effects_reversed: true
- pending_compensation_named: true
- irreversible_skips_explicit: true
- missing_reverse_plumbing_failures: false
- malformed_receipt_failures_named: true

<!-- field: traces -->
#### Traces

- srd025-rollback-lifecycle AC3
- srd025-rollback-lifecycle AC4
- srd035-checkpoint-port AC3
- srd036-dolt-state-persistence AC6

<!-- case: Lifecycle profiles wire checkpoint_history and checkpoint_rollback | status: done -->
### Lifecycle profiles wire checkpoint_history and checkpoint_rollback

The history and rollback lifecycle profiles load their MachineSpec,
select checkpoint_history and checkpoint_rollback respectively, and
register the builtin adapters before the machine dispatches the action.

<!-- field: use_case -->
#### Use case

rel02.0-uc002-history-rollback

<!-- field: go_test -->
#### Go test

cd ../applications/catalog && go test ./conformance -run '^TestLifecycleHistoryRollback$'

<!-- field: inputs | form: "tree" -->
#### Inputs

- profiles:
  1. lifecycle/history
  1. lifecycle/rollback

<!-- field: expected | form: "tree" -->
#### Expected

- profile_loaded: true
- selected_as_builtin: true

<!-- field: traces -->
#### Traces

- srd026-lifecycle-tools AC1
