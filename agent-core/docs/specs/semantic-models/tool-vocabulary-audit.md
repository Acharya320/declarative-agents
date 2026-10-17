---
id: sm-tool-vocabulary-audit
title: Tool Vocabulary Audit
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: purpose -->
## Purpose

Classifies current agent-core tools using the vocabulary rule:
tools are words, machines compose sentences. A good tool is a single
deterministic verb with typed input, typed output, declared side effects,
explicit reversibility, and clear composition relationships.

<!-- field: scope | form: "tree" -->
## Scope

- shared_declarations:
  1. tools/builtin/
  1. tools/exec/
- mode_declarations:
  1. $AGENT_CATALOG_ROOT/agents/*/builtin.yaml
  1. $AGENT_CATALOG_ROOT/agents/*/llm/*.yaml
- selections:
  1. $AGENT_CATALOG_ROOT/agents/executor/tools.yaml
  1. $AGENT_CATALOG_ROOT/agents/planner/tools.yaml
  1. $AGENT_CATALOG_ROOT/agents/critic/tools.yaml
  1. $AGENT_CATALOG_ROOT/agents/critic/tools-point.yaml
  1. $AGENT_CATALOG_ROOT/agents/bench/tools.yaml
  1. $AGENT_CATALOG_ROOT/agents/specification-critic/tools.yaml
  1. $AGENT_CATALOG_ROOT/testdata/conformance/lifecycle/tools.yaml

<!-- field: role_legend | form: "tree" -->
## Role legend

- word:
  A single deterministic verb. It can be composed freely by a state
  machine or by LLM-selected $tool dispatch.
- sentence:
  A workflow hidden behind one tool name. These should usually be split
  so machine.yaml owns the sequence.
- boundary:
  A tool that crosses a program boundary: LLM call, child agent,
  nested machine, web UI event source, or external CLI with broad
  side effects. Boundary tools can remain, but require explicit contracts.
- stateful_internal:
  An internal machine action that mutates agent/domain state. It may be
  one word, but undo must be command-state undo rather than workspace undo.

<!-- field: overall_findings -->
## Overall findings

- Most file and simple exec tools are close to vocabulary words.
- Tool declarations now include structured side_effects (14 kind values),
  reversibility classification, and undo strategy for most tools. Contract
  metadata (problem, goals, non_goals) is populated for many but not all
  tools. output.schema is declared for file tools and test but is not yet
  broadly populated.
- All boundary tools now have UndoMementoProvider implementations that emit
  compensatable mementos with boundary_compensation payloads. Undo contracts
  are declared in tools/builtin/ and tools/exec/ declaration files with
  payload and captures fields.
- The remaining high-priority sentence-shaped tools are now mostly planner or lifecycle follow-ups; evaluator session and point workflows have been split into machine-composed words.
- The desired lifecycle history and rollback tools are first-class STL
  builtin adapters selected by lifecycle profiles. Legacy exec wrappers and
  hidden Cobra history/rollback commands are migration artifacts only.
- Generator validation is machine-declared through explicit build, lint, and test states; the deprecated validate aggregator has been retired.
- Boundary tools carry explicit category (boundary) classification in their declarations.
- Runtime output is still mostly human-readable core.Result.Output strings. Existing output.schema declarations are contract targets, not yet guaranteed serialized Result.Data.

<!-- field: tools_note -->
## Tools note

The tools list below runs in these groups, in order:

- File tools: read, write, edit, find, list_files
- LLM and conversation tools: invoke_llm, parse_response, report_parse_error, reset_history, nudge_reread, done, suspend, checkpoint_history, checkpoint_rollback, self_invoke
- Validation and build exec tools: build, vet, lint, test
- Planner tools: extract_task, select_all_ready, seed_passthrough_plan, mark_nodes_planning, project_planner_context, capture_planner_failure, compose_planner_prompt, parse_plan, format_issue, create_tracker_issue, record_tracker_issue, mark_nodes_executing, format_task_file, invoke_executor, reset_retry_count, increment_retry, check_retry_limit, mark_task_done, mark_task_failed, remaining_work
- Evaluator session tools: parse_suite_config, discover_suite_samples, expand_eval_grid, init_eval_session, report_suite_summary, materialize_eval_points, run_point, report_session
- Evaluator point tools: create_point_dir, sample_docs, commit_workspace_baseline, record_agent_commit, dump_config, run_agent, run_oracle_check, record_oracle_result, collect_trace_tokens, check_agent_version, summarize_point_results, record_point_failure, collect_metrics
- Spec validation tools: load_corpus, validate_specs, format_report
- Workspace and git exec tools: copy_dir, make_dir, git_init, stage_all, workspace_status, commit, rev_parse, diff_stat, log_oneline
- Issue exec tools: issue_create, issue_close, issue_list

Bench owns no core-specific tools. Its profile aliases generic REST lifecycle/event, evaluation-query, predicate, compose, exec, filesystem, and self_invoke words.

<!-- field: tools | form: "tree" -->
## Tools

1. - name: read
   - role: word
   - selected_by:
     1. generator
   - inputs_from:
     1. find
     1. list_files
     1. prior_tool_output
     1. user_task
   - feeds_into:
     1. edit
     1. write
     1. LLM_reasoning
   - side_effects: none
   - output_contract: declared
   - undo: noop
   - overlaps:
     1. - tool: list_files
        - difference: read returns file content; list_files returns directory shape
     1. - tool: find
        - difference: read returns known file content; find discovers matching paths
   - issues:
     1. needs full problem/goals/non_goals/reversibility metadata
   - migration_priority: low
1. - name: write
   - role: word
   - selected_by:
     1. generator
   - inputs_from:
     1. LLM_tool_call
   - feeds_into:
     1. read
     1. build
     1. lint
     1. test
     1. nudge_reread
   - side_effects: filesystem_write
   - output_contract: declared
   - undo: workspace_restore
   - overlaps:
     1. - tool: edit
        - difference: write replaces an entire file; edit replaces one exact occurrence
   - issues:
     1. needs structured side_effects and reversibility contract
     1. runtime output is still prose, not Result.Data
   - migration_priority: medium
1. - name: edit
   - role: word
   - selected_by:
     1. generator
   - inputs_from:
     1. read
     1. LLM_tool_call
   - feeds_into:
     1. read
     1. build
     1. lint
     1. test
     1. nudge_reread
   - side_effects: filesystem_write
   - output_contract: declared
   - undo: workspace_restore
   - overlaps:
     1. - tool: write
        - difference: edit is safer for localized changes and fails on ambiguous matches
   - issues:
     1. needs structured side_effects and reversibility contract
     1. runtime output is still prose, not Result.Data
   - migration_priority: medium
1. - name: find
   - role: word
   - selected_by:
     1. generator
   - inputs_from:
     1. LLM_tool_call
   - feeds_into:
     1. read
     1. edit
   - side_effects: none
   - output_contract: missing
   - undo: noop
   - overlaps:
     1. - tool: list_files
        - difference: find searches file contents with regex; list_files lists paths
   - issues:
     1. needs output.schema for path/match/line tuples
     1. needs contract metadata
   - migration_priority: low
1. - name: list_files
   - role: word
   - selected_by:
     1. generator
   - inputs_from:
     1. LLM_tool_call
   - feeds_into:
     1. read
     1. find
   - side_effects: none
   - output_contract: missing
   - undo: noop
   - overlaps:
     1. - tool: find
        - difference: list_files explores directory structure; find searches text
   - issues:
     1. needs output.schema for tree entries
     1. implementation silently skips unreadable subdirectories in walkTree
   - migration_priority: low
1. - name: invoke_llm
   - role: boundary
   - selected_by:
     1. generator
     1. planner
   - inputs_from:
     1. conversation_history
     1. registry_manifest
     1. prompt_assembler
   - feeds_into:
     1. parse_response
     1. parse_plan
   - side_effects: conversation_append, external_model_call
   - output_contract: missing
   - undo: conversation_truncate
   - overlaps: []
   - issues:
     1. not deterministic despite temperature/seed controls
     1. should declare output.schema for raw response plus token usage
   - migration_priority: medium
1. - name: parse_response
   - role: word
   - selected_by:
     1. generator
   - inputs_from:
     1. invoke_llm
   - feeds_into:
     1. $tool_dispatch
     1. done
     1. report_parse_error
   - side_effects: none
   - output_contract: missing
   - undo: noop
   - overlaps:
     1. - tool: parse_plan
        - difference: parse_response extracts a tool call or done signal; parse_plan extracts planner YAML
   - issues:
     1. needs output.schema for parsed tool request and completion summary
   - migration_priority: low
1. - name: report_parse_error
   - role: stateful_internal
   - selected_by:
     1. generator
   - inputs_from:
     1. parse_response
   - feeds_into:
     1. invoke_llm
     1. BudgetExceeded
   - side_effects: parse_retry_counter_update
   - output_contract: missing
   - undo: parse_retry_counter_restore
   - overlaps: []
   - issues:
     1. output.schema should distinguish retry feedback from retry-limit exhaustion
   - migration_priority: medium
1. - name: reset_history
   - role: stateful_internal
   - selected_by:
     1. generator
     1. planner
   - inputs_from:
     1. machine_recovery_path
   - feeds_into:
     1. invoke_llm
   - side_effects: conversation_replace
   - output_contract: missing
   - undo: conversation_restore
   - overlaps: []
   - issues:
     1. irreversible today unless prior messages are captured by command undo
   - migration_priority: medium
1. - name: nudge_reread
   - role: stateful_internal
   - selected_by:
     1. generator
   - inputs_from:
     1. edit
   - feeds_into:
     1. invoke_llm
   - side_effects: conversation_append
   - output_contract: missing
   - undo: conversation_truncate
   - overlaps: []
   - issues:
     1. should declare why this is state mutation rather than workflow logic
   - migration_priority: medium
1. - name: done
   - role: word
   - selected_by:
     1. generator
   - inputs_from:
     1. LLM_tool_call
   - feeds_into:
     1. generator_validation_sequence
   - side_effects: none
   - output_contract: missing
   - undo: noop
   - overlaps: []
   - issues:
     1. needs output.schema for completion summary
   - migration_priority: low
1. - name: suspend
   - role: boundary
   - selected_by: []
   - inputs_from:
     1. machine_transition
   - feeds_into:
     1. resume_or_checkpoint
   - side_effects: state_mutation
   - output_contract: missing
   - undo: boundary_compensation (compensating_action, captures reason + checkpoint_required)
   - overlaps: []
   - issues:
     1. output.schema for suspension reason and checkpoint reference
   - migration_priority: low
1. - name: checkpoint_history
   - role: stateful_internal
   - status: desired_lifecycle_adapter
   - selected_by:
     1. lifecycle_history
   - inputs_from:
     1. lifecycle_profile
     1. run_selector
     1. checkpoint_backend
   - feeds_into:
     1. history_machine_terminal
     1. operator_review
   - side_effects: none
   - output_contract: declared_by_srd026
   - undo: noop
   - overlaps:
     1. - tool: legacy_history_exec_wrapper
        - difference: checkpoint_history reads the checkpoint backend directly; the legacy wrapper shells through compatibility CLI behavior
     1. - tool: checkpoint_rollback
        - difference: checkpoint_history reads Dolt commit history without reverting the run branch or replaying receipts
   - issues:
     1. ToolDef declaration still needs to replace legacy exec wrapper in lifecycle config
   - migration_priority: high
1. - name: checkpoint_rollback
   - role: boundary
   - status: desired_lifecycle_adapter
   - selected_by:
     1. lifecycle_rollback
   - inputs_from:
     1. lifecycle_profile
     1. run_selector
     1. to_step
     1. checkpoint_backend
   - feeds_into:
     1. rollback_machine_terminal
     1. resume_decision
   - side_effects: dolt_branch_revert, receipt_walk_undo
   - output_contract: declared_by_srd026
   - undo: boundary_compensation (operator_step_selection, captures run + operator_decision)
   - overlaps:
     1. - tool: legacy_rollback_exec_wrapper
        - difference: checkpoint_rollback reverts the Dolt run branch and replays receipts directly; the legacy wrapper shells through compatibility CLI behavior
     1. - tool: checkpoint_history
        - difference: checkpoint_rollback reverts DB state (DB-only) and reverses external effects through tool receipts; checkpoint_history only reads
   - issues:
     1. ToolDef declaration still needs to replace legacy exec wrapper in lifecycle config
     1. Requires lifecycle profile wiring before it is the desired operator path
   - migration_priority: high
1. - name: self_invoke
   - role: boundary
   - selected_by: []
   - inputs_from:
     1. machine_transition
   - feeds_into:
     1. child_agent_result
   - side_effects: child_agent_execution
   - output_contract: missing
   - undo: boundary_compensation (child_agent_workspace_restore, captures child_machine + child_tools + child_run_id + trace_path)
   - overlaps:
     1. - tool: invoke_executor
        - difference: invoke_executor is a planner profile alias of self_invoke
   - issues:
     1. output.schema for child agent result
   - migration_priority: low
1. - name: build
   - role: boundary
   - selected_by:
     1. generator
     1. planner
   - inputs_from:
     1. workspace_files
   - feeds_into:
     1. lint
     1. test
     1. increment_retry
   - side_effects: subprocess, possible_build_cache_or_binary
   - output_contract: missing
   - undo: noop_or_workspace_restore_if_artifacts
   - overlaps:
     1. - tool: test
        - difference: build compiles packages; test compiles and runs tests
     1. - tool: vet
        - difference: vet reports suspicious code patterns
     1. - tool: lint
        - difference: lint applies configured static analysis rules
   - issues:
     1. external CLI boundary needs structured side_effects and output.schema
   - migration_priority: medium
1. - name: vet
   - role: boundary
   - selected_by:
     1. generator
     1. planner
   - inputs_from:
     1. workspace_files
   - feeds_into:
     1. LLM_reasoning
     1. increment_retry
   - side_effects: subprocess
   - output_contract: missing
   - undo: noop
   - overlaps:
     1. - tool: lint
        - difference: vet is Go toolchain analysis; lint is external linter rule set
   - issues:
     1. needs output.schema for diagnostics
   - migration_priority: low
1. - name: lint
   - role: boundary
   - selected_by:
     1. generator
   - inputs_from:
     1. workspace_files
   - feeds_into:
     1. test
   - side_effects: subprocess
   - output_contract: missing
   - undo: noop
   - overlaps:
     1. - tool: vet
        - difference: lint is golangci-lint; vet is Go toolchain analysis
   - issues:
     1. needs output.schema for diagnostics
   - migration_priority: low
1. - name: test
   - role: boundary
   - selected_by:
     1. generator
     1. planner
   - inputs_from:
     1. workspace_files
   - feeds_into:
     1. done
     1. increment_retry
     1. mark_task_done
   - side_effects: subprocess, possible_test_cache
   - output_contract: declared
   - undo: noop
   - overlaps:
     1. - tool: build
        - difference: test executes package tests; build only compiles
   - issues:
     1. runtime output is not yet structured according to output.schema
   - migration_priority: medium
1. - name: extract_task
   - role: stateful_internal
   - selected_by:
     1. planner
   - inputs_from:
     1. spec_graph
   - feeds_into:
     1. project_planner_context
   - side_effects: pipeline_state_mutation
   - output_contract: missing
   - undo: graph_snapshot_restore
   - overlaps:
     1. - tool: select_all_ready
        - difference: extract_task selects one weighted task; select_all_ready selects every ready node without plan or graph mutation
   - issues:
     1. needs command-state undo contract
   - migration_priority: medium
1. - name: select_all_ready
   - role: stateful_internal
   - selected_by:
     1. planner
   - inputs_from:
     1. spec_graph
   - feeds_into:
     1. seed_passthrough_plan
   - side_effects: pipeline_state_mutation
   - output_contract: declared
   - undo: pipeline_state_restore
   - overlaps:
     1. - tool: extract_task
        - difference: select_all_ready selects every ready node; extract_task selects one weighted task
   - issues: []
   - migration_priority: complete
1. - name: seed_passthrough_plan
   - role: stateful_internal
   - selected_by:
     1. planner
   - inputs_from:
     1. select_all_ready
     1. profile_config
   - feeds_into:
     1. mark_nodes_planning
   - side_effects: pipeline_state_mutation
   - output_contract: declared
   - undo: pipeline_state_restore
   - overlaps: []
   - issues: []
   - migration_priority: complete
1. - name: mark_nodes_planning
   - role: stateful_internal
   - selected_by:
     1. planner
   - inputs_from:
     1. seed_passthrough_plan
   - feeds_into:
     1. mark_nodes_executing
   - side_effects: requirement_graph_mutation
   - output_contract: declared
   - undo: pipeline_state_restore
   - overlaps: []
   - issues: []
   - migration_priority: complete
1. - name: project_planner_context
   - role: word
   - selected_by:
     1. planner
   - inputs_from:
     1. extract_task
     1. spec_corpus
   - feeds_into:
     1. compose_planner_prompt
   - side_effects: none
   - output_contract: declared structured context
   - undo: noop
   - overlaps: []
   - issues: []
   - migration_priority: complete
1. - name: capture_planner_failure
   - role: word
   - selected_by:
     1. planner
   - inputs_from:
     1. vet
     1. build
     1. test
   - feeds_into:
     1. increment_retry
     1. compose_planner_prompt
   - side_effects: none
   - output_contract: labelled failure context
   - undo: noop
   - overlaps: []
   - issues: []
   - migration_priority: complete
1. - name: compose_planner_prompt
   - role: word
   - selected_by:
     1. planner
   - inputs_from:
     1. project_planner_context
     1. failure_context
   - feeds_into:
     1. invoke_llm
   - side_effects: none
   - output_contract: profile-owned prompt text
   - undo: noop
   - overlaps:
     1. - tool: compose
        - difference: compose_planner_prompt is the planner profile's configured compose alias
   - issues: []
   - migration_priority: complete
1. - name: parse_plan
   - role: word
   - selected_by:
     1. planner
   - inputs_from:
     1. invoke_llm
   - feeds_into:
     1. format_issue
   - side_effects: pipeline_state_mutation
   - output_contract: missing
   - undo: restore_previous_plan
   - overlaps:
     1. - tool: parse_response
        - difference: parse_plan parses planner YAML; parse_response parses generator tool calls
   - issues:
     1. mutates CurrentPlan; needs undo contract
   - migration_priority: medium
1. - name: format_issue
   - role: stateful_internal
   - selected_by:
     1. planner
   - inputs_from:
     1. parse_plan
   - feeds_into:
     1. write
   - side_effects: planner_state_read
   - output_contract: structured tracker-neutral parameters
   - undo: restore_previous_pipeline_snapshot
   - overlaps: []
   - issues: []
   - migration_priority: low
1. - name: create_tracker_issue
   - role: boundary
   - selected_by:
     1. planner
   - inputs_from:
     1. write
     1. format_issue
   - feeds_into:
     1. record_tracker_issue
   - side_effects: configured_tracker_write
   - output_contract: profile-declared JSON id
   - undo: boundary_compensation (compensating_action, captures issue_id)
   - overlaps:
     1. - tool: issue_create
        - difference: create_tracker_issue is the planner profile's replaceable tracker boundary; issue_create is a shared direct bd CLI word
   - issues: []
   - migration_priority: low
1. - name: record_tracker_issue
   - role: stateful_internal
   - selected_by:
     1. planner
   - inputs_from:
     1. create_tracker_issue
   - feeds_into:
     1. mark_nodes_executing
     1. extract_task
   - side_effects: pipeline_state_mutation
   - output_contract: structured issue_id
   - undo: restore_previous_pipeline_snapshot
   - overlaps: []
   - issues: []
   - migration_priority: low
1. - name: mark_nodes_executing
   - role: stateful_internal
   - selected_by:
     1. planner
   - inputs_from:
     1. record_tracker_issue
     1. mark_nodes_planning
   - feeds_into:
     1. format_task_file
   - side_effects: requirement_graph_mutation
   - output_contract: declared
   - undo: pipeline_state_restore
   - overlaps: []
   - issues: []
   - migration_priority: complete
1. - name: format_task_file
   - role: word
   - selected_by:
     1. planner
   - inputs_from:
     1. current_plan
   - feeds_into:
     1. write
   - side_effects: none
   - output_contract: declared write parameters
   - undo: noop
   - overlaps: []
   - issues: []
   - migration_priority: complete
1. - name: invoke_executor
   - role: boundary
   - selected_by:
     1. planner
   - inputs_from:
     1. write
   - feeds_into:
     1. vet
     1. build
     1. test
   - side_effects: child_agent_process, workspace_mutation
   - output_contract: shared self_invoke result
   - undo: boundary_compensation
   - overlaps:
     1. - tool: self_invoke
        - difference: invoke_executor is the planner profile alias configured for the executor profile
   - issues: []
   - migration_priority: complete
1. - name: reset_retry_count
   - role: word
   - selected_by:
     1. planner
   - inputs_from:
     1. TaskExtracted
   - feeds_into:
     1. increment_retry
   - side_effects: none
   - output_contract: labelled command-state zero
   - undo: noop (history truncation removes output)
   - overlaps: []
   - issues: []
   - migration_priority: low
1. - name: increment_retry
   - role: word
   - selected_by:
     1. planner
   - inputs_from:
     1. labelled_command_state
   - feeds_into:
     1. check_retry_limit
   - side_effects: none
   - output_contract: structured retry count
   - undo: noop (history truncation removes output)
   - overlaps: []
   - issues: []
   - migration_priority: low
1. - name: check_retry_limit
   - role: word
   - selected_by:
     1. planner
   - inputs_from:
     1. increment_retry
   - feeds_into:
     1. invoke_executor
     1. remaining_work
   - side_effects: none
   - output_contract: declared predicate signals
   - undo: noop
   - overlaps: []
   - issues: []
   - migration_priority: low
1. - name: mark_task_done
   - role: stateful_internal
   - selected_by:
     1. planner
   - inputs_from:
     1. test
   - feeds_into:
     1. remaining_work
   - side_effects: pipeline_graph_completion
   - output_contract: declared TaskCompleted signal
   - undo: restore_pipeline_graph_snapshot
   - overlaps: []
   - issues: []
   - migration_priority: low
1. - name: mark_task_failed
   - role: stateful_internal
   - selected_by:
     1. planner
   - inputs_from:
     1. check_retry_limit
   - feeds_into:
     1. remaining_work
   - side_effects: pipeline_graph_failure
   - output_contract: declared TaskFailed signal
   - undo: restore_pipeline_graph_snapshot
   - overlaps: []
   - issues: []
   - migration_priority: complete
1. - name: remaining_work
   - role: word
   - selected_by:
     1. planner
   - inputs_from:
     1. mark_task_done
     1. mark_task_failed
   - feeds_into:
     1. extract_task
     1. Completed
     1. Stalled
   - side_effects: none
   - output_contract: declared WorkRemaining, AllDone, or Blocked signal
   - undo: noop
   - overlaps: []
   - issues: []
   - migration_priority: low
1. - name: parse_suite_config
   - role: stateful_internal
   - selected_by:
     1. evaluator_session
   - inputs_from:
     1. suite_yaml
     1. cli_input
     1. tool_config
   - feeds_into:
     1. discover_suite_samples
   - side_effects: eval_session_suite_metadata_mutation
   - output_contract: declared_by_signal
   - undo: session_state_restore
   - overlaps: []
   - issues:
     1. needs output.schema for suite metadata summary
   - migration_priority: medium
1. - name: discover_suite_samples
   - role: stateful_internal
   - selected_by:
     1. evaluator_session
   - inputs_from:
     1. parse_suite_config
   - feeds_into:
     1. expand_eval_grid
   - side_effects: eval_session_sample_mutation
   - output_contract: declared_by_signal
   - undo: session_state_restore
   - overlaps:
     1. - tool: load_corpus
        - difference: discover_suite_samples discovers evaluator sample workspaces; load_corpus loads project specs
   - issues:
     1. needs output.schema for discovered sample count and names
   - migration_priority: medium
1. - name: expand_eval_grid
   - role: stateful_internal
   - selected_by:
     1. evaluator_session
   - inputs_from:
     1. discover_suite_samples
   - feeds_into:
     1. init_eval_session
   - side_effects: eval_session_grid_mutation
   - output_contract: declared_by_signal
   - undo: session_state_restore
   - overlaps: []
   - issues:
     1. needs output.schema for grid point count
   - migration_priority: medium
1. - name: init_eval_session
   - role: stateful_internal
   - selected_by:
     1. evaluator_session
   - inputs_from:
     1. expand_eval_grid
     1. tool_config
   - feeds_into:
     1. report_suite_summary
   - side_effects: filesystem_directory_create, eval_session_runtime_state_mutation
   - output_contract: declared_by_signal
   - undo: session_state_restore_and_workspace_restore
   - overlaps: []
   - issues:
     1. creates session directory and should eventually participate in workspace rollback
   - migration_priority: medium
1. - name: report_suite_summary
   - role: word
   - selected_by:
     1. evaluator_session
   - inputs_from:
     1. init_eval_session
   - feeds_into:
     1. materialize_eval_points
   - side_effects: stderr_write
   - output_contract: declared_by_signal
   - undo: noop
   - overlaps:
     1. - tool: report_session
        - difference: report_suite_summary reports setup dimensions; report_session reports final outcomes
   - issues:
     1. needs output.schema for total point count
   - migration_priority: low
1. - name: materialize_eval_points
   - role: word
   - selected_by:
     1. evaluator_session
   - inputs_from:
     1. report_suite_summary
   - feeds_into:
     1. run_point
   - side_effects: none
   - output_contract: ordered evaluation point array
   - undo: noop
   - overlaps: []
   - issues: []
   - migration_priority: low
1. - name: run_point
   - role: boundary
   - selected_by:
     1. evaluator_session
   - inputs_from:
     1. materialize_eval_points
   - feeds_into:
     1. run_point
     1. report_session
   - side_effects: nested_machine_execution, point_result_recording, filesystem_mutation
   - output_contract: missing
   - undo: boundary_compensation (child_command_undo, captures eval_session_snapshot + nested_machine_history)
   - overlaps:
     1. - tool: invoke_executor
        - difference: run_point invokes an in-process nested machine; invoke_executor uses self_invoke
   - issues:
     1. output.schema for point result summary
   - migration_priority: medium
1. - name: report_session
   - role: word
   - selected_by:
     1. evaluator_session
   - inputs_from:
     1. run_point
   - feeds_into:
     1. terminal
   - side_effects: stderr_output
   - output_contract: missing
   - undo: noop
   - overlaps:
     1. - tool: format_report
        - difference: report_session summarizes eval run; format_report formats spec validation findings
   - issues:
     1. should declare whether stderr output is a side effect
   - migration_priority: low
1. - name: create_point_dir
   - role: word
   - selected_by:
     1. evaluator_point
   - inputs_from:
     1. run_point
   - feeds_into:
     1. copy_dir
   - side_effects: filesystem_write, point_context_mutation
   - output_contract: copy_dir source and destination parameters
   - undo: workspace_restore
   - overlaps: []
   - issues: []
   - migration_priority: low
1. - name: sample_docs
   - role: word
   - selected_by:
     1. evaluator_point
   - inputs_from:
     1. copy_dir
   - feeds_into:
     1. copy_dir
     1. git_init
   - side_effects: none
   - output_contract: optional-docs presence and copy_dir parameters
   - undo: noop
   - overlaps: []
   - issues: []
   - migration_priority: low
1. - name: commit_workspace_baseline
   - role: boundary
   - selected_by:
     1. evaluator_point
   - inputs_from:
     1. stage_all
   - feeds_into:
     1. rev_parse
   - side_effects: git_commit_create
   - output_contract: missing
   - undo: workspace_restore
   - overlaps:
     1. - tool: commit
        - difference: commit_workspace_baseline is a profile-configured alias with deterministic local identity and message; commit accepts caller input
   - issues:
     1. needs output.schema for commit hash
   - migration_priority: medium
1. - name: record_agent_commit
   - role: stateful_internal
   - selected_by:
     1. evaluator_point
   - inputs_from:
     1. rev_parse
   - feeds_into:
     1. dump_config
   - side_effects: point_context_commit_mutation
   - output_contract: commit provenance string
   - undo: restore_point_context
   - overlaps: []
   - issues: []
   - migration_priority: low
1. - name: dump_config
   - role: word
   - selected_by:
     1. evaluator_point
   - inputs_from:
     1. record_agent_commit
     1. point_context
   - feeds_into:
     1. run_agent
   - side_effects: filesystem_write
   - output_contract: missing
   - undo: workspace_restore
   - overlaps: []
   - issues:
     1. needs output.schema for experiment.yaml path and bytes
     1. needs structured filesystem side effect
   - migration_priority: medium
1. - name: run_agent
   - role: boundary
   - selected_by:
     1. evaluator_point
   - inputs_from:
     1. dump_config
     1. point_context
   - feeds_into:
     1. run_oracle_check
   - side_effects: subprocess, workspace_mutation, trace_file_write
   - output_contract: missing
   - undo: boundary_compensation (compensating_action, captures point_context_snapshot + point_workspace + trace_path + result_path)
   - overlaps:
     1. - tool: invoke_executor
        - difference: run_agent runs a suite harness binary; invoke_executor uses self_invoke for the planner executor
   - issues:
     1. output.schema for agent exit code and trace path
   - migration_priority: medium
1. - name: run_oracle_check
   - role: boundary
   - selected_by:
     1. evaluator_point
   - inputs_from:
     1. run_agent
   - feeds_into:
     1. record_oracle_result
   - side_effects: configured_child_process
   - output_contract: raw configured oracle output
   - undo: noop
   - overlaps:
     1. - tool: test
        - difference: run_oracle_check is a critic profile alias with profile-owned command arguments; generic test accepts a package parameter
   - issues: []
   - migration_priority: low
1. - name: record_oracle_result
   - role: stateful_internal
   - selected_by:
     1. evaluator_point
   - inputs_from:
     1. run_oracle_check
   - feeds_into:
     1. collect_trace_tokens
   - side_effects: point_context_test_result_mutation
   - output_contract: configured oracle output
   - undo: restore_point_context
   - overlaps: []
   - issues: []
   - migration_priority: low
1. - name: collect_trace_tokens
   - role: word
   - selected_by:
     1. evaluator_point
   - inputs_from:
     1. record_oracle_result
     1. trace_ndjson
   - feeds_into:
     1. check_agent_version
   - side_effects: point_context_token_mutation
   - output_contract: declared_by_signal
   - undo: restore_point_context
   - overlaps:
     1. - tool: collect_metrics
        - difference: collect_trace_tokens extracts token state; collect_metrics persists meta.json
   - issues:
     1. output.schema still needs full structured token summary
   - migration_priority: medium
1. - name: check_agent_version
   - role: word
   - selected_by:
     1. evaluator_point
   - inputs_from:
     1. collect_trace_tokens
     1. trace_ndjson
     1. harness_config
   - feeds_into:
     1. summarize_point_results
   - side_effects: point_context_version_mutation_and_warning_output
   - output_contract: declared_by_signal
   - undo: restore_point_context
   - overlaps: []
   - issues:
     1. version mismatch is a warning signal that continues the point sequence
   - migration_priority: medium
1. - name: summarize_point_results
   - role: word
   - selected_by:
     1. evaluator_point
   - inputs_from:
     1. record_oracle_result
     1. collect_trace_tokens
     1. check_agent_version
   - feeds_into:
     1. collect_metrics
   - side_effects: none
   - output_contract: missing
   - undo: noop
   - overlaps:
     1. - tool: collect_metrics
        - difference: summarize_point_results emits aggregate status; collect_metrics writes meta.json
   - issues:
     1. output.schema still needs full structured point summary
   - migration_priority: medium
1. - name: record_point_failure
   - role: stateful_internal
   - selected_by:
     1. evaluator_point
   - inputs_from:
     1. failed_point_command
   - feeds_into:
     1. collect_metrics
   - side_effects: point_context_mutation
   - output_contract: failure stage and cause summary
   - undo: restore_point_context
   - overlaps: []
   - issues: []
   - migration_priority: low
1. - name: collect_metrics
   - role: word
   - selected_by:
     1. evaluator_point
   - inputs_from:
     1. summarize_point_results
     1. record_point_failure
   - feeds_into:
     1. terminal
   - side_effects: filesystem_write
   - output_contract: missing
   - undo: workspace_restore
   - overlaps:
     1. - tool: collect_trace_tokens
        - difference: collect_metrics writes meta.json; collect_trace_tokens computes token state
     1. - tool: summarize_point_results
        - difference: collect_metrics persists metadata; summarize_point_results emits aggregate status
   - issues:
     1. needs output.schema for meta path and summary fields
   - migration_priority: medium
1. - name: load_corpus
   - role: stateful_internal
   - selected_by:
     1. jurist
   - inputs_from:
     1. workspace_specs
   - feeds_into:
     1. validate_specs
   - side_effects: validation_state_mutation
   - output_contract: missing
   - undo: restore_validation_state
   - overlaps:
     1. - tool: parse_suite_config
        - difference: load_corpus loads project specs; parse_suite_config loads evaluator suite metadata
     1. - tool: discover_suite_samples
        - difference: load_corpus loads project specs; discover_suite_samples finds evaluator sample workspaces
   - issues:
     1. needs output.schema for corpus counts
   - migration_priority: medium
1. - name: validate_specs
   - role: word
   - selected_by:
     1. jurist
   - inputs_from:
     1. load_corpus
   - feeds_into:
     1. format_report
   - side_effects: validation_state_mutation
   - output_contract: missing
   - undo: restore_validation_findings
   - issues:
     1. mutates validation findings, so undo contract must be explicit
   - migration_priority: medium
1. - name: format_report
   - role: word
   - selected_by:
     1. jurist
   - inputs_from:
     1. validate_specs
   - feeds_into:
     1. terminal
   - side_effects: none
   - output_contract: missing
   - undo: noop
   - overlaps:
     1. - tool: report_session
        - difference: format_report formats spec findings; report_session summarizes eval session
   - issues:
     1. needs output.schema for formatted report and error count
   - migration_priority: low
1. - name: copy_dir
   - role: word
   - selected_by:
     1. evaluator_point
   - inputs_from:
     1. create_point_dir
     1. tool_call
   - feeds_into:
     1. sample_docs
     1. copy_dir
     1. git_init
     1. build
     1. test
   - side_effects: filesystem_write
   - output_contract: missing
   - undo: workspace_restore
   - overlaps: []
   - issues:
     1. legacy scalar side_effects
   - migration_priority: medium
1. - name: make_dir
   - role: word
   - selected_by: []
   - inputs_from:
     1. tool_call
   - feeds_into:
     1. write
     1. copy_dir
   - side_effects: filesystem_write
   - output_contract: missing
   - undo: workspace_restore
   - overlaps: []
   - issues:
     1. legacy scalar side_effects
   - migration_priority: medium
1. - name: git_init
   - role: word
   - selected_by:
     1. evaluator_point
   - inputs_from:
     1. sample_docs
     1. copy_dir
     1. make_dir
     1. workspace
   - feeds_into:
     1. stage_all
     1. commit
   - side_effects: filesystem_write_git_metadata
   - output_contract: missing
   - undo: workspace_restore
   - overlaps: []
   - issues:
     1. legacy scalar side_effects
   - migration_priority: medium
1. - name: stage_all
   - role: word
   - selected_by:
     1. evaluator_point
   - inputs_from:
     1. git_init
     1. workspace_changes
   - feeds_into:
     1. commit
     1. commit_workspace_baseline
   - side_effects: git_index_mutation
   - output_contract: missing
   - undo: git_reset_index_or_workspace_restore
   - overlaps: []
   - issues:
     1. irreversible risk if used broadly without explicit selection
   - migration_priority: medium
1. - name: workspace_status
   - role: word
   - selected_by: []
   - inputs_from:
     1. workspace
   - feeds_into:
     1. stage_all
     1. commit
     1. LLM_reasoning
   - side_effects: none
   - output_contract: declared
   - undo: noop
   - overlaps:
     1. - tool: diff_stat
        - difference: workspace_status lists porcelain file status; diff_stat summarizes diff size
   - issues:
     1. runtime output is not yet structured according to output.schema
   - migration_priority: low
1. - name: commit
   - role: boundary
   - selected_by: []
   - inputs_from:
     1. stage_all
   - feeds_into:
     1. rev_parse
     1. log_oneline
   - side_effects: git_commit_create
   - output_contract: missing
   - undo: boundary_compensation (compensating_action, captures workspace_paths + git_ref)
   - overlaps:
     1. - tool: issue_close
        - difference: commit mutates git history; issue_close mutates bd issue state
   - issues:
     1. output.schema for commit hash
   - migration_priority: medium
1. - name: rev_parse
   - role: word
   - selected_by:
     1. evaluator_point
   - inputs_from:
     1. commit
     1. commit_workspace_baseline
     1. workspace
   - feeds_into:
     1. logs
     1. metadata
     1. record_agent_commit
   - side_effects: none
   - output_contract: missing
   - undo: noop
   - overlaps:
     1. - tool: log_oneline
        - difference: rev_parse returns one ref; log_oneline returns commit history
   - issues:
     1. needs output.schema for ref string
   - migration_priority: low
1. - name: diff_stat
   - role: word
   - selected_by: []
   - inputs_from:
     1. workspace_changes
   - feeds_into:
     1. LLM_reasoning
     1. commit_decision
   - side_effects: none
   - output_contract: missing
   - undo: noop
   - overlaps:
     1. - tool: workspace_status
        - difference: diff_stat summarizes diff size; workspace_status lists file statuses
   - issues:
     1. needs output.schema for file/insertion/deletion counts
   - migration_priority: low
1. - name: log_oneline
   - role: word
   - selected_by: []
   - inputs_from:
     1. git_repo
   - feeds_into:
     1. LLM_reasoning
   - side_effects: none
   - output_contract: missing
   - undo: noop
   - overlaps:
     1. - tool: rev_parse
        - difference: log_oneline returns history; rev_parse returns one ref
   - issues:
     1. needs output.schema for commit list
   - migration_priority: low
1. - name: issue_create
   - role: boundary
   - selected_by: []
   - inputs_from:
     1. planner_output
     1. user_request
   - feeds_into:
     1. issue_list
     1. issue_close
   - side_effects: bd_database_write
   - output_contract: missing
   - undo: boundary_compensation (compensating_action, captures created_issue_id + workspace_paths)
   - overlaps:
     1. - tool: create_tracker_issue
        - difference: issue_create is a shared direct bd CLI wrapper; create_tracker_issue is selected and configured by the planner profile
   - issues:
     1. output.schema for created issue ID
   - migration_priority: low
1. - name: issue_close
   - role: boundary
   - selected_by: []
   - inputs_from:
     1. issue_create
     1. user_request
   - feeds_into:
     1. issue_list
   - side_effects: bd_database_write
   - output_contract: missing
   - undo: boundary_compensation (compensating_action, captures issue_id + previous_issue_status + workspace_paths)
   - overlaps: []
   - issues:
     1. output.schema for closed issue confirmation
   - migration_priority: low
1. - name: issue_list
   - role: word
   - selected_by: []
   - inputs_from:
     1. bd_database
   - feeds_into:
     1. issue_create
     1. issue_close
     1. planner
   - side_effects: none
   - output_contract: missing
   - undo: noop
   - overlaps:
     1. - tool: workspace_status
        - difference: issue_list reports tracker state; workspace_status reports git state
   - issues:
     1. needs output.schema for issue summaries
   - migration_priority: low

<!-- field: sentence_tools | form: "tree" -->
## Sentence tools

- highest_priority: []

<!-- field: boundary_tools_with_compensation | form: "tree" -->
## Boundary tools with compensation

- note:
  All boundary tools now implement UndoMementoProvider and declare
  payload: boundary_compensation with captures in their tool YAML.
- tools:
  1. suspend
  1. checkpoint_rollback
  1. self_invoke
  1. create_tracker_issue
  1. invoke_executor
  1. run_point
  1. run_agent
  1. commit
  1. issue_create
  1. issue_close

<!-- field: remaining_contract_gaps | form: "tree" -->
## Remaining contract gaps

- typed_outputs:
  - declared:
    1. read
    1. write
    1. edit
    1. test
    1. workspace_status
  - still_missing:
    1. find
    1. list_files
    1. invoke_llm
    1. parse_response
    1. checkpoint_history
    1. checkpoint_rollback
    1. planner tools
    1. evaluator session tools
    1. evaluator point tools
    1. spec validation tools
    1. most exec tools
- structured_side_effects:
  - status: partially_populated
  - note:
    Structured side_effects entries (kind/target/description) are populated
    for many tools. Some exec tools still use legacy scalar form.
- undo_metadata:
  - status: populated
  - note:
    All boundary tools declare undo strategy, payload, and captures in their
    YAML declarations. UndoMementoProvider is implemented for boundary tools
    and workspace-mutating builtins. Simple read-only tools use noop.
- relationships:
  - status: partially_populated
  - note:
    Composition relationships are documented in this audit. Migration into
    individual ToolDef entries is incremental.
