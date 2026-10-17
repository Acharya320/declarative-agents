---
id: sm-command-undo-audit
title: Command Undo Audit
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: purpose -->
## Purpose

Tracks which commands encode a sufficient Receipt during Execute and consume
it in a receipt-consuming Undo, which return an empty receipt because they are
read-only, and which are compensatable or irreversible. The persisted opaque
Receipt is both the audit/rollback payload and the exact input a tool's Undo
decodes. It documents current rollback coverage and the places where external
compensation is explicit rather than automatic.

<!-- field: current_contract -->
## Current contract

Every core.Command exposes a receipt-consuming Undo that takes the restored,
receipt-bearing Result and reverses the external effect. Read-only words
intentionally return an empty receipt and core.NoopUndo. State-mutating words
encode a Receipt during Execute that captures what Undo needs: prior file
content, an external resource id, a workspace git commit hash the tool manages
itself, or an in-process snapshot of conversation, planner, evaluator, point,
or validation state.

Persisted rollback replays receipts. The srd036 Dolt adapter reverts
DB-persisted state to a target step, then the checkpoint_rollback lifecycle
tool walks the reverted execution entries in reverse and hands each restored
Result to its tool's Undo. The engine and checkpoint adapters never interpret
a receipt; only the originating tool decodes it. The validator checks receipt
presence for reversible state-mutating tools; receipt sufficiency is each
tool's own round-trip test.

<!-- field: noop_contract | form: "tree" -->
## Noop contract

- signal: ToolDone
- receipt: empty
- output: "undo: no-op"
- command_name: the command being undone

<!-- field: reversibility_kinds -->
## Reversibility kinds

<!-- field: reversible -->
### Reversible

The receipt holds enough for the tool's Undo to fully reverse the effect (prior content, snapshot, or a managed commit hash).

<!-- field: compensatable -->
### Compensatable

The receipt holds metadata for operator or higher-level compensation; external audit history, subprocess behavior, or human/model decisions cannot be erased.

<!-- field: irreversible -->
### Irreversible

No receipt can reverse the effect; the tool requires confirmation and the reverse walk skips and logs it.

<!-- field: compensation_boundaries -->
## Compensation boundaries

<!-- field: description -->
### Description

These commands encode compensatable receipts. Rollback can preserve the
child run, artifact, issue, workspace, or approval metadata needed for
operator or higher-level compensation, but it does not pretend external
audit history, subprocess behavior, or human/model decisions can be erased.

<!-- field: commands -->
### Commands

| command | receipt_payload |
| --- | --- |
| create_tracker_issue | configured_tracker_issue_id |
| self_invoke | child_run_id_machine_tools_and_trace_artifact |
| run_agent | point_context_workspace_and_artifact_paths |
| run_point | eval_session_snapshot_and_nested_machine |
| suspend | approval_reason_and_checkpoint_requirement |

<!-- field: expected_noop_long_term | form: "tree" -->
## Expected noop long term

1. read
1. find
1. list_files
1. done
1. project_planner_context
1. capture_planner_failure
1. nudge_reread
1. format_report
1. summarize_point_results
1. report_suite_summary
1. - read_only_exec_tools:
     1. build
     1. vet
     1. lint
     1. test
     1. workspace_status
     1. rev_parse
     1. diff_stat
     1. log_oneline
     1. issue_list
1. FailedParamCmd
1. failCmd

<!-- field: receipt_undo -->
## Receipt undo

<!-- field: conversation_and_retry_state -->
### Conversation and retry state

| command | strategy | receipt |
| --- | --- | --- |
| invoke_llm | conversation_truncate | reversible_conversation_snapshot |
| reset_history | conversation_restore | reversible_conversation_snapshot |
| parse_response | parse_retry_counter_restore_when_tracker_enabled | reversible_parse_retry_counter_snapshot |
| report_parse_error | parse_retry_counter_restore_when_tracker_enabled | reversible_parse_retry_counter_snapshot |

<!-- field: planner_state -->
### Planner state

| command | strategy | receipt |
| --- | --- | --- |
| load_graph | pipeline_state_restore | reversible_pipeline_snapshot |
| extract_task | pipeline_state_restore | reversible_pipeline_snapshot |
| select_all_ready | pipeline_state_restore | reversible_pipeline_snapshot |
| seed_passthrough_plan | pipeline_state_restore | reversible_pipeline_snapshot |
| mark_nodes_planning | pipeline_state_restore | reversible_pipeline_snapshot |
| parse_plan | pipeline_state_restore | reversible_pipeline_snapshot |
| record_tracker_issue | pipeline_state_restore | reversible_pipeline_snapshot |
| mark_task_done | pipeline_state_restore | reversible_pipeline_snapshot |
| mark_task_failed | pipeline_state_restore | reversible_pipeline_snapshot |

<!-- field: workspace_state -->
### Workspace state

| command | strategy | receipt |
| --- | --- | --- |
| write | file_snapshot_restore | reversible_prior_content_or_commit_hash |
| edit | file_snapshot_restore | reversible_prior_content_or_commit_hash |
| copy_dir | tool_owned_commit_restore | reversible_commit_hash |
| make_dir | tool_owned_commit_restore | reversible_commit_hash |
| git_init | tool_owned_commit_restore | reversible_commit_hash |
| stage_all | tool_owned_commit_restore | reversible_commit_hash |
| commit | compensating_action | compensatable_boundary_compensation |
| issue_create | compensating_action | compensatable_boundary_compensation |
| issue_close | compensating_action | compensatable_boundary_compensation |

<!-- field: evaluator_state -->
### Evaluator state

| command | strategy | receipt |
| --- | --- | --- |
| parse_suite_config | evaluator_session_restore | reversible_eval_session_snapshot |
| discover_suite_samples | evaluator_session_restore | reversible_eval_session_snapshot |
| expand_eval_grid | evaluator_session_restore | reversible_eval_session_snapshot |
| init_eval_session | owned_artifact_removal_and_evaluator_session_restore | reversible_eval_session_snapshot_and_owned_artifact_paths |
| materialize_eval_points | noop | none_read_only_projection |
| run_point | compensating_action | compensatable_eval_session_snapshot_and_boundary_compensation |
| report_session | evaluator_session_restore | reversible_eval_session_snapshot |
| create_point_dir | owned_artifact_removal_and_point_context_restore | reversible_point_context_snapshot_and_owned_artifact_paths |
| sample_docs | noop | none_read_only_projection |
| dump_config | owned_artifact_removal_and_point_context_restore | reversible_point_context_snapshot_and_owned_artifact_paths |
| collect_metrics | owned_artifact_removal_and_point_context_restore | reversible_point_context_snapshot_and_owned_artifact_paths |
| record_point_failure | point_context_restore | reversible_point_context_snapshot |
| record_oracle_result | point_context_restore | reversible_point_context_snapshot |
| record_agent_commit | point_context_restore | reversible_point_context_snapshot |
| collect_trace_tokens | point_context_restore | reversible_point_context_snapshot |
| check_agent_version | point_context_restore | reversible_point_context_snapshot |

<!-- field: validation_state -->
### Validation state

| command | strategy | receipt |
| --- | --- | --- |
| load_corpus | validation_state_restore | reversible_validation_state_snapshot |
| validate_specs | validation_state_restore | reversible_validation_state_snapshot |

<!-- field: boundary_compensation | form: "typed table" -->
### Boundary compensation

| command | strategy | receipt | captures |
| --- | --- | --- | --- |
| suspend | resume_or_checkpoint_rollback | compensatable_boundary_compensation | [reason, checkpoint_required] |
| create_tracker_issue | close_or_delete_created_issue | compensatable_boundary_compensation | [issue_id] |
| self_invoke | child_agent_workspace_restore | compensatable_boundary_compensation | [child_run_id, child_machine, child_tools, trace_artifact] |
| run_point | nested_machine_rollback | compensatable_eval_session_snapshot_and_boundary_compensation | [eval_session_snapshot, nested_machine_history, child_machine] |
| run_agent | point_workspace_restore_and_child_process_compensation | compensatable_boundary_compensation | [point_context_snapshot, workspace_paths, artifact_paths] |

<!-- field: declaration_mismatches | form: "empty list" -->
## Declaration mismatches

None

<!-- field: related_documents -->
## Related documents

- srd035-checkpoint-port
- srd036-dolt-state-persistence
- srd025-rollback-lifecycle
- sm-rollback-lifecycle
- cfg-tool-authoring-standard

<!-- field: related_issues -->
## Related issues

- agent-core-gm-rollback.5
- agent-core-gm-rollback.6
- agent-core-gm-rollback.7
- agent-core-gm-rollback.8
