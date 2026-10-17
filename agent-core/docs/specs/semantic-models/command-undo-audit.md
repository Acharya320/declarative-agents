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

<!-- field: reversibility_kinds | form: "tree" -->
## Reversibility kinds

- reversible: The receipt holds enough for the tool's Undo to fully reverse the effect (prior content, snapshot, or a managed commit hash).
- compensatable: The receipt holds metadata for operator or higher-level compensation; external audit history, subprocess behavior, or human/model decisions cannot be erased.
- irreversible: No receipt can reverse the effect; the tool requires confirmation and the reverse walk skips and logs it.

<!-- field: compensation_boundaries | form: "tree" -->
## Compensation boundaries

- description:
  These commands encode compensatable receipts. Rollback can preserve the
  child run, artifact, issue, workspace, or approval metadata needed for
  operator or higher-level compensation, but it does not pretend external
  audit history, subprocess behavior, or human/model decisions can be erased.
- commands:
  1. - command: create_tracker_issue
     - receipt_payload: configured_tracker_issue_id
  1. - command: self_invoke
     - receipt_payload: child_run_id_machine_tools_and_trace_artifact
  1. - command: run_agent
     - receipt_payload: point_context_workspace_and_artifact_paths
  1. - command: run_point
     - receipt_payload: eval_session_snapshot_and_nested_machine
  1. - command: suspend
     - receipt_payload: approval_reason_and_checkpoint_requirement

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

<!-- field: receipt_undo | form: "tree" -->
## Receipt undo

- conversation_and_retry_state:
  1. - command: invoke_llm
     - strategy: conversation_truncate
     - receipt: reversible_conversation_snapshot
  1. - command: reset_history
     - strategy: conversation_restore
     - receipt: reversible_conversation_snapshot
  1. - command: parse_response
     - strategy: parse_retry_counter_restore_when_tracker_enabled
     - receipt: reversible_parse_retry_counter_snapshot
  1. - command: report_parse_error
     - strategy: parse_retry_counter_restore_when_tracker_enabled
     - receipt: reversible_parse_retry_counter_snapshot
- planner_state:
  1. - command: load_graph
     - strategy: pipeline_state_restore
     - receipt: reversible_pipeline_snapshot
  1. - command: extract_task
     - strategy: pipeline_state_restore
     - receipt: reversible_pipeline_snapshot
  1. - command: select_all_ready
     - strategy: pipeline_state_restore
     - receipt: reversible_pipeline_snapshot
  1. - command: seed_passthrough_plan
     - strategy: pipeline_state_restore
     - receipt: reversible_pipeline_snapshot
  1. - command: mark_nodes_planning
     - strategy: pipeline_state_restore
     - receipt: reversible_pipeline_snapshot
  1. - command: parse_plan
     - strategy: pipeline_state_restore
     - receipt: reversible_pipeline_snapshot
  1. - command: record_tracker_issue
     - strategy: pipeline_state_restore
     - receipt: reversible_pipeline_snapshot
  1. - command: mark_task_done
     - strategy: pipeline_state_restore
     - receipt: reversible_pipeline_snapshot
  1. - command: mark_task_failed
     - strategy: pipeline_state_restore
     - receipt: reversible_pipeline_snapshot
- workspace_state:
  1. - command: write
     - strategy: file_snapshot_restore
     - receipt: reversible_prior_content_or_commit_hash
  1. - command: edit
     - strategy: file_snapshot_restore
     - receipt: reversible_prior_content_or_commit_hash
  1. - command: copy_dir
     - strategy: tool_owned_commit_restore
     - receipt: reversible_commit_hash
  1. - command: make_dir
     - strategy: tool_owned_commit_restore
     - receipt: reversible_commit_hash
  1. - command: git_init
     - strategy: tool_owned_commit_restore
     - receipt: reversible_commit_hash
  1. - command: stage_all
     - strategy: tool_owned_commit_restore
     - receipt: reversible_commit_hash
  1. - command: commit
     - strategy: compensating_action
     - receipt: compensatable_boundary_compensation
  1. - command: issue_create
     - strategy: compensating_action
     - receipt: compensatable_boundary_compensation
  1. - command: issue_close
     - strategy: compensating_action
     - receipt: compensatable_boundary_compensation
- evaluator_state:
  1. - command: parse_suite_config
     - strategy: evaluator_session_restore
     - receipt: reversible_eval_session_snapshot
  1. - command: discover_suite_samples
     - strategy: evaluator_session_restore
     - receipt: reversible_eval_session_snapshot
  1. - command: expand_eval_grid
     - strategy: evaluator_session_restore
     - receipt: reversible_eval_session_snapshot
  1. - command: init_eval_session
     - strategy: owned_artifact_removal_and_evaluator_session_restore
     - receipt: reversible_eval_session_snapshot_and_owned_artifact_paths
  1. - command: materialize_eval_points
     - strategy: noop
     - receipt: none_read_only_projection
  1. - command: run_point
     - strategy: compensating_action
     - receipt: compensatable_eval_session_snapshot_and_boundary_compensation
  1. - command: report_session
     - strategy: evaluator_session_restore
     - receipt: reversible_eval_session_snapshot
  1. - command: create_point_dir
     - strategy: owned_artifact_removal_and_point_context_restore
     - receipt: reversible_point_context_snapshot_and_owned_artifact_paths
  1. - command: sample_docs
     - strategy: noop
     - receipt: none_read_only_projection
  1. - command: dump_config
     - strategy: owned_artifact_removal_and_point_context_restore
     - receipt: reversible_point_context_snapshot_and_owned_artifact_paths
  1. - command: collect_metrics
     - strategy: owned_artifact_removal_and_point_context_restore
     - receipt: reversible_point_context_snapshot_and_owned_artifact_paths
  1. - command: record_point_failure
     - strategy: point_context_restore
     - receipt: reversible_point_context_snapshot
  1. - command: record_oracle_result
     - strategy: point_context_restore
     - receipt: reversible_point_context_snapshot
  1. - command: record_agent_commit
     - strategy: point_context_restore
     - receipt: reversible_point_context_snapshot
  1. - command: collect_trace_tokens
     - strategy: point_context_restore
     - receipt: reversible_point_context_snapshot
  1. - command: check_agent_version
     - strategy: point_context_restore
     - receipt: reversible_point_context_snapshot
- validation_state:
  1. - command: load_corpus
     - strategy: validation_state_restore
     - receipt: reversible_validation_state_snapshot
  1. - command: validate_specs
     - strategy: validation_state_restore
     - receipt: reversible_validation_state_snapshot
- boundary_compensation:
  1. - command: suspend
     - strategy: resume_or_checkpoint_rollback
     - receipt: compensatable_boundary_compensation
     - captures:
       1. reason
       1. checkpoint_required
  1. - command: create_tracker_issue
     - strategy: close_or_delete_created_issue
     - receipt: compensatable_boundary_compensation
     - captures:
       1. issue_id
  1. - command: self_invoke
     - strategy: child_agent_workspace_restore
     - receipt: compensatable_boundary_compensation
     - captures:
       1. child_run_id
       1. child_machine
       1. child_tools
       1. trace_artifact
  1. - command: run_point
     - strategy: nested_machine_rollback
     - receipt: compensatable_eval_session_snapshot_and_boundary_compensation
     - captures:
       1. eval_session_snapshot
       1. nested_machine_history
       1. child_machine
  1. - command: run_agent
     - strategy: point_workspace_restore_and_child_process_compensation
     - receipt: compensatable_boundary_compensation
     - captures:
       1. point_context_snapshot
       1. workspace_paths
       1. artifact_paths

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
