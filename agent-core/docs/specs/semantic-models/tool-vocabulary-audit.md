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

<!-- field: role_legend -->
## Role legend

<!-- field: word -->
### Word

A single deterministic verb. It can be composed freely by a state
machine or by LLM-selected $tool dispatch.

<!-- field: sentence -->
### Sentence

A workflow hidden behind one tool name. These should usually be split
so machine.yaml owns the sequence.

<!-- field: boundary -->
### Boundary

A tool that crosses a program boundary: LLM call, child agent,
nested machine, web UI event source, or external CLI with broad
side effects. Boundary tools can remain, but require explicit contracts.

<!-- field: stateful_internal -->
### Stateful internal

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

<!-- field: tools -->
## Tools

<!-- tool: read -->
### read

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- generator

<!-- field: inputs_from -->
#### Inputs from

- find
- list_files
- prior_tool_output
- user_task

<!-- field: feeds_into -->
#### Feeds into

- edit
- write
- LLM_reasoning

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

declared

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| list_files | read returns file content; list_files returns directory shape |
| find | read returns known file content; find discovers matching paths |

<!-- field: issues -->
#### Issues

- needs full problem/goals/non_goals/reversibility metadata

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: write -->
### write

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- generator

<!-- field: inputs_from -->
#### Inputs from

- LLM_tool_call

<!-- field: feeds_into -->
#### Feeds into

- read
- build
- lint
- test
- nudge_reread

<!-- field: side_effects -->
#### Side effects

filesystem_write

<!-- field: output_contract -->
#### Output contract

declared

<!-- field: undo -->
#### Undo

workspace_restore

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| edit | write replaces an entire file; edit replaces one exact occurrence |

<!-- field: issues -->
#### Issues

- needs structured side_effects and reversibility contract
- runtime output is still prose, not Result.Data

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: edit -->
### edit

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- generator

<!-- field: inputs_from -->
#### Inputs from

- read
- LLM_tool_call

<!-- field: feeds_into -->
#### Feeds into

- read
- build
- lint
- test
- nudge_reread

<!-- field: side_effects -->
#### Side effects

filesystem_write

<!-- field: output_contract -->
#### Output contract

declared

<!-- field: undo -->
#### Undo

workspace_restore

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| write | edit is safer for localized changes and fails on ambiguous matches |

<!-- field: issues -->
#### Issues

- needs structured side_effects and reversibility contract
- runtime output is still prose, not Result.Data

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: find -->
### find

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- generator

<!-- field: inputs_from -->
#### Inputs from

- LLM_tool_call

<!-- field: feeds_into -->
#### Feeds into

- read
- edit

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| list_files | find searches file contents with regex; list_files lists paths |

<!-- field: issues -->
#### Issues

- needs output.schema for path/match/line tuples
- needs contract metadata

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: list_files -->
### list_files

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- generator

<!-- field: inputs_from -->
#### Inputs from

- LLM_tool_call

<!-- field: feeds_into -->
#### Feeds into

- read
- find

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| find | list_files explores directory structure; find searches text |

<!-- field: issues -->
#### Issues

- needs output.schema for tree entries
- implementation silently skips unreadable subdirectories in walkTree

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: invoke_llm -->
### invoke_llm

<!-- field: role -->
#### Role

boundary

<!-- field: selected_by -->
#### Selected by

- generator
- planner

<!-- field: inputs_from -->
#### Inputs from

- conversation_history
- registry_manifest
- prompt_assembler

<!-- field: feeds_into -->
#### Feeds into

- parse_response
- parse_plan

<!-- field: side_effects -->
#### Side effects

conversation_append, external_model_call

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

conversation_truncate

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues -->
#### Issues

- not deterministic despite temperature/seed controls
- should declare output.schema for raw response plus token usage

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: parse_response -->
### parse_response

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- generator

<!-- field: inputs_from -->
#### Inputs from

- invoke_llm

<!-- field: feeds_into -->
#### Feeds into

- $tool_dispatch
- done
- report_parse_error

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| parse_plan | parse_response extracts a tool call or done signal; parse_plan extracts planner YAML |

<!-- field: issues -->
#### Issues

- needs output.schema for parsed tool request and completion summary

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: report_parse_error -->
### report_parse_error

<!-- field: role -->
#### Role

stateful_internal

<!-- field: selected_by -->
#### Selected by

- generator

<!-- field: inputs_from -->
#### Inputs from

- parse_response

<!-- field: feeds_into -->
#### Feeds into

- invoke_llm
- BudgetExceeded

<!-- field: side_effects -->
#### Side effects

parse_retry_counter_update

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

parse_retry_counter_restore

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues -->
#### Issues

- output.schema should distinguish retry feedback from retry-limit exhaustion

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: reset_history -->
### reset_history

<!-- field: role -->
#### Role

stateful_internal

<!-- field: selected_by -->
#### Selected by

- generator
- planner

<!-- field: inputs_from -->
#### Inputs from

- machine_recovery_path

<!-- field: feeds_into -->
#### Feeds into

- invoke_llm

<!-- field: side_effects -->
#### Side effects

conversation_replace

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

conversation_restore

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues -->
#### Issues

- irreversible today unless prior messages are captured by command undo

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: nudge_reread -->
### nudge_reread

<!-- field: role -->
#### Role

stateful_internal

<!-- field: selected_by -->
#### Selected by

- generator

<!-- field: inputs_from -->
#### Inputs from

- edit

<!-- field: feeds_into -->
#### Feeds into

- invoke_llm

<!-- field: side_effects -->
#### Side effects

conversation_append

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

conversation_truncate

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues -->
#### Issues

- should declare why this is state mutation rather than workflow logic

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: done -->
### done

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- generator

<!-- field: inputs_from -->
#### Inputs from

- LLM_tool_call

<!-- field: feeds_into -->
#### Feeds into

- generator_validation_sequence

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues -->
#### Issues

- needs output.schema for completion summary

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: suspend -->
### suspend

<!-- field: role -->
#### Role

boundary

<!-- field: selected_by | form: "empty list" -->
#### Selected by

None

<!-- field: inputs_from -->
#### Inputs from

- machine_transition

<!-- field: feeds_into -->
#### Feeds into

- resume_or_checkpoint

<!-- field: side_effects -->
#### Side effects

state_mutation

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

boundary_compensation (compensating_action, captures reason + checkpoint_required)

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues -->
#### Issues

- output.schema for suspension reason and checkpoint reference

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: checkpoint_history -->
### checkpoint_history

<!-- field: role -->
#### Role

stateful_internal

<!-- field: status -->
#### Status

desired_lifecycle_adapter

<!-- field: selected_by -->
#### Selected by

- lifecycle_history

<!-- field: inputs_from -->
#### Inputs from

- lifecycle_profile
- run_selector
- checkpoint_backend

<!-- field: feeds_into -->
#### Feeds into

- history_machine_terminal
- operator_review

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

declared_by_srd026

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| legacy_history_exec_wrapper | checkpoint_history reads the checkpoint backend directly; the legacy wrapper shells through compatibility CLI behavior |
| checkpoint_rollback | checkpoint_history reads Dolt commit history without reverting the run branch or replaying receipts |

<!-- field: issues -->
#### Issues

- ToolDef declaration still needs to replace legacy exec wrapper in lifecycle config

<!-- field: migration_priority -->
#### Migration priority

high

<!-- tool: checkpoint_rollback -->
### checkpoint_rollback

<!-- field: role -->
#### Role

boundary

<!-- field: status -->
#### Status

desired_lifecycle_adapter

<!-- field: selected_by -->
#### Selected by

- lifecycle_rollback

<!-- field: inputs_from -->
#### Inputs from

- lifecycle_profile
- run_selector
- to_step
- checkpoint_backend

<!-- field: feeds_into -->
#### Feeds into

- rollback_machine_terminal
- resume_decision

<!-- field: side_effects -->
#### Side effects

dolt_branch_revert, receipt_walk_undo

<!-- field: output_contract -->
#### Output contract

declared_by_srd026

<!-- field: undo -->
#### Undo

boundary_compensation (operator_step_selection, captures run + operator_decision)

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| legacy_rollback_exec_wrapper | checkpoint_rollback reverts the Dolt run branch and replays receipts directly; the legacy wrapper shells through compatibility CLI behavior |
| checkpoint_history | checkpoint_rollback reverts DB state (DB-only) and reverses external effects through tool receipts; checkpoint_history only reads |

<!-- field: issues -->
#### Issues

- ToolDef declaration still needs to replace legacy exec wrapper in lifecycle config
- Requires lifecycle profile wiring before it is the desired operator path

<!-- field: migration_priority -->
#### Migration priority

high

<!-- tool: self_invoke -->
### self_invoke

<!-- field: role -->
#### Role

boundary

<!-- field: selected_by | form: "empty list" -->
#### Selected by

None

<!-- field: inputs_from -->
#### Inputs from

- machine_transition

<!-- field: feeds_into -->
#### Feeds into

- child_agent_result

<!-- field: side_effects -->
#### Side effects

child_agent_execution

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

boundary_compensation (child_agent_workspace_restore, captures child_machine + child_tools + child_run_id + trace_path)

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| invoke_executor | invoke_executor is a planner profile alias of self_invoke |

<!-- field: issues -->
#### Issues

- output.schema for child agent result

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: build -->
### build

<!-- field: role -->
#### Role

boundary

<!-- field: selected_by -->
#### Selected by

- generator
- planner

<!-- field: inputs_from -->
#### Inputs from

- workspace_files

<!-- field: feeds_into -->
#### Feeds into

- lint
- test
- increment_retry

<!-- field: side_effects -->
#### Side effects

subprocess, possible_build_cache_or_binary

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

noop_or_workspace_restore_if_artifacts

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| test | build compiles packages; test compiles and runs tests |
| vet | vet reports suspicious code patterns |
| lint | lint applies configured static analysis rules |

<!-- field: issues -->
#### Issues

- external CLI boundary needs structured side_effects and output.schema

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: vet -->
### vet

<!-- field: role -->
#### Role

boundary

<!-- field: selected_by -->
#### Selected by

- generator
- planner

<!-- field: inputs_from -->
#### Inputs from

- workspace_files

<!-- field: feeds_into -->
#### Feeds into

- LLM_reasoning
- increment_retry

<!-- field: side_effects -->
#### Side effects

subprocess

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| lint | vet is Go toolchain analysis; lint is external linter rule set |

<!-- field: issues -->
#### Issues

- needs output.schema for diagnostics

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: lint -->
### lint

<!-- field: role -->
#### Role

boundary

<!-- field: selected_by -->
#### Selected by

- generator

<!-- field: inputs_from -->
#### Inputs from

- workspace_files

<!-- field: feeds_into -->
#### Feeds into

- test

<!-- field: side_effects -->
#### Side effects

subprocess

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| vet | lint is golangci-lint; vet is Go toolchain analysis |

<!-- field: issues -->
#### Issues

- needs output.schema for diagnostics

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: test -->
### test

<!-- field: role -->
#### Role

boundary

<!-- field: selected_by -->
#### Selected by

- generator
- planner

<!-- field: inputs_from -->
#### Inputs from

- workspace_files

<!-- field: feeds_into -->
#### Feeds into

- done
- increment_retry
- mark_task_done

<!-- field: side_effects -->
#### Side effects

subprocess, possible_test_cache

<!-- field: output_contract -->
#### Output contract

declared

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| build | test executes package tests; build only compiles |

<!-- field: issues -->
#### Issues

- runtime output is not yet structured according to output.schema

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: extract_task -->
### extract_task

<!-- field: role -->
#### Role

stateful_internal

<!-- field: selected_by -->
#### Selected by

- planner

<!-- field: inputs_from -->
#### Inputs from

- spec_graph

<!-- field: feeds_into -->
#### Feeds into

- project_planner_context

<!-- field: side_effects -->
#### Side effects

pipeline_state_mutation

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

graph_snapshot_restore

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| select_all_ready | extract_task selects one weighted task; select_all_ready selects every ready node without plan or graph mutation |

<!-- field: issues -->
#### Issues

- needs command-state undo contract

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: select_all_ready -->
### select_all_ready

<!-- field: role -->
#### Role

stateful_internal

<!-- field: selected_by -->
#### Selected by

- planner

<!-- field: inputs_from -->
#### Inputs from

- spec_graph

<!-- field: feeds_into -->
#### Feeds into

- seed_passthrough_plan

<!-- field: side_effects -->
#### Side effects

pipeline_state_mutation

<!-- field: output_contract -->
#### Output contract

declared

<!-- field: undo -->
#### Undo

pipeline_state_restore

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| extract_task | select_all_ready selects every ready node; extract_task selects one weighted task |

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

complete

<!-- tool: seed_passthrough_plan -->
### seed_passthrough_plan

<!-- field: role -->
#### Role

stateful_internal

<!-- field: selected_by -->
#### Selected by

- planner

<!-- field: inputs_from -->
#### Inputs from

- select_all_ready
- profile_config

<!-- field: feeds_into -->
#### Feeds into

- mark_nodes_planning

<!-- field: side_effects -->
#### Side effects

pipeline_state_mutation

<!-- field: output_contract -->
#### Output contract

declared

<!-- field: undo -->
#### Undo

pipeline_state_restore

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

complete

<!-- tool: mark_nodes_planning -->
### mark_nodes_planning

<!-- field: role -->
#### Role

stateful_internal

<!-- field: selected_by -->
#### Selected by

- planner

<!-- field: inputs_from -->
#### Inputs from

- seed_passthrough_plan

<!-- field: feeds_into -->
#### Feeds into

- mark_nodes_executing

<!-- field: side_effects -->
#### Side effects

requirement_graph_mutation

<!-- field: output_contract -->
#### Output contract

declared

<!-- field: undo -->
#### Undo

pipeline_state_restore

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

complete

<!-- tool: project_planner_context -->
### project_planner_context

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- planner

<!-- field: inputs_from -->
#### Inputs from

- extract_task
- spec_corpus

<!-- field: feeds_into -->
#### Feeds into

- compose_planner_prompt

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

declared structured context

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

complete

<!-- tool: capture_planner_failure -->
### capture_planner_failure

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- planner

<!-- field: inputs_from -->
#### Inputs from

- vet
- build
- test

<!-- field: feeds_into -->
#### Feeds into

- increment_retry
- compose_planner_prompt

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

labelled failure context

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

complete

<!-- tool: compose_planner_prompt -->
### compose_planner_prompt

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- planner

<!-- field: inputs_from -->
#### Inputs from

- project_planner_context
- failure_context

<!-- field: feeds_into -->
#### Feeds into

- invoke_llm

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

profile-owned prompt text

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| compose | compose_planner_prompt is the planner profile's configured compose alias |

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

complete

<!-- tool: parse_plan -->
### parse_plan

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- planner

<!-- field: inputs_from -->
#### Inputs from

- invoke_llm

<!-- field: feeds_into -->
#### Feeds into

- format_issue

<!-- field: side_effects -->
#### Side effects

pipeline_state_mutation

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

restore_previous_plan

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| parse_response | parse_plan parses planner YAML; parse_response parses generator tool calls |

<!-- field: issues -->
#### Issues

- mutates CurrentPlan; needs undo contract

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: format_issue -->
### format_issue

<!-- field: role -->
#### Role

stateful_internal

<!-- field: selected_by -->
#### Selected by

- planner

<!-- field: inputs_from -->
#### Inputs from

- parse_plan

<!-- field: feeds_into -->
#### Feeds into

- write

<!-- field: side_effects -->
#### Side effects

planner_state_read

<!-- field: output_contract -->
#### Output contract

structured tracker-neutral parameters

<!-- field: undo -->
#### Undo

restore_previous_pipeline_snapshot

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: create_tracker_issue -->
### create_tracker_issue

<!-- field: role -->
#### Role

boundary

<!-- field: selected_by -->
#### Selected by

- planner

<!-- field: inputs_from -->
#### Inputs from

- write
- format_issue

<!-- field: feeds_into -->
#### Feeds into

- record_tracker_issue

<!-- field: side_effects -->
#### Side effects

configured_tracker_write

<!-- field: output_contract -->
#### Output contract

profile-declared JSON id

<!-- field: undo -->
#### Undo

boundary_compensation (compensating_action, captures issue_id)

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| issue_create | create_tracker_issue is the planner profile's replaceable tracker boundary; issue_create is a shared direct bd CLI word |

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: record_tracker_issue -->
### record_tracker_issue

<!-- field: role -->
#### Role

stateful_internal

<!-- field: selected_by -->
#### Selected by

- planner

<!-- field: inputs_from -->
#### Inputs from

- create_tracker_issue

<!-- field: feeds_into -->
#### Feeds into

- mark_nodes_executing
- extract_task

<!-- field: side_effects -->
#### Side effects

pipeline_state_mutation

<!-- field: output_contract -->
#### Output contract

structured issue_id

<!-- field: undo -->
#### Undo

restore_previous_pipeline_snapshot

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: mark_nodes_executing -->
### mark_nodes_executing

<!-- field: role -->
#### Role

stateful_internal

<!-- field: selected_by -->
#### Selected by

- planner

<!-- field: inputs_from -->
#### Inputs from

- record_tracker_issue
- mark_nodes_planning

<!-- field: feeds_into -->
#### Feeds into

- format_task_file

<!-- field: side_effects -->
#### Side effects

requirement_graph_mutation

<!-- field: output_contract -->
#### Output contract

declared

<!-- field: undo -->
#### Undo

pipeline_state_restore

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

complete

<!-- tool: format_task_file -->
### format_task_file

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- planner

<!-- field: inputs_from -->
#### Inputs from

- current_plan

<!-- field: feeds_into -->
#### Feeds into

- write

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

declared write parameters

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

complete

<!-- tool: invoke_executor -->
### invoke_executor

<!-- field: role -->
#### Role

boundary

<!-- field: selected_by -->
#### Selected by

- planner

<!-- field: inputs_from -->
#### Inputs from

- write

<!-- field: feeds_into -->
#### Feeds into

- vet
- build
- test

<!-- field: side_effects -->
#### Side effects

child_agent_process, workspace_mutation

<!-- field: output_contract -->
#### Output contract

shared self_invoke result

<!-- field: undo -->
#### Undo

boundary_compensation

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| self_invoke | invoke_executor is the planner profile alias configured for the executor profile |

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

complete

<!-- tool: reset_retry_count -->
### reset_retry_count

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- planner

<!-- field: inputs_from -->
#### Inputs from

- TaskExtracted

<!-- field: feeds_into -->
#### Feeds into

- increment_retry

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

labelled command-state zero

<!-- field: undo -->
#### Undo

noop (history truncation removes output)

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: increment_retry -->
### increment_retry

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- planner

<!-- field: inputs_from -->
#### Inputs from

- labelled_command_state

<!-- field: feeds_into -->
#### Feeds into

- check_retry_limit

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

structured retry count

<!-- field: undo -->
#### Undo

noop (history truncation removes output)

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: check_retry_limit -->
### check_retry_limit

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- planner

<!-- field: inputs_from -->
#### Inputs from

- increment_retry

<!-- field: feeds_into -->
#### Feeds into

- invoke_executor
- remaining_work

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

declared predicate signals

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: mark_task_done -->
### mark_task_done

<!-- field: role -->
#### Role

stateful_internal

<!-- field: selected_by -->
#### Selected by

- planner

<!-- field: inputs_from -->
#### Inputs from

- test

<!-- field: feeds_into -->
#### Feeds into

- remaining_work

<!-- field: side_effects -->
#### Side effects

pipeline_graph_completion

<!-- field: output_contract -->
#### Output contract

declared TaskCompleted signal

<!-- field: undo -->
#### Undo

restore_pipeline_graph_snapshot

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: mark_task_failed -->
### mark_task_failed

<!-- field: role -->
#### Role

stateful_internal

<!-- field: selected_by -->
#### Selected by

- planner

<!-- field: inputs_from -->
#### Inputs from

- check_retry_limit

<!-- field: feeds_into -->
#### Feeds into

- remaining_work

<!-- field: side_effects -->
#### Side effects

pipeline_graph_failure

<!-- field: output_contract -->
#### Output contract

declared TaskFailed signal

<!-- field: undo -->
#### Undo

restore_pipeline_graph_snapshot

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

complete

<!-- tool: remaining_work -->
### remaining_work

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- planner

<!-- field: inputs_from -->
#### Inputs from

- mark_task_done
- mark_task_failed

<!-- field: feeds_into -->
#### Feeds into

- extract_task
- Completed
- Stalled

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

declared WorkRemaining, AllDone, or Blocked signal

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: parse_suite_config -->
### parse_suite_config

<!-- field: role -->
#### Role

stateful_internal

<!-- field: selected_by -->
#### Selected by

- evaluator_session

<!-- field: inputs_from -->
#### Inputs from

- suite_yaml
- cli_input
- tool_config

<!-- field: feeds_into -->
#### Feeds into

- discover_suite_samples

<!-- field: side_effects -->
#### Side effects

eval_session_suite_metadata_mutation

<!-- field: output_contract -->
#### Output contract

declared_by_signal

<!-- field: undo -->
#### Undo

session_state_restore

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues -->
#### Issues

- needs output.schema for suite metadata summary

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: discover_suite_samples -->
### discover_suite_samples

<!-- field: role -->
#### Role

stateful_internal

<!-- field: selected_by -->
#### Selected by

- evaluator_session

<!-- field: inputs_from -->
#### Inputs from

- parse_suite_config

<!-- field: feeds_into -->
#### Feeds into

- expand_eval_grid

<!-- field: side_effects -->
#### Side effects

eval_session_sample_mutation

<!-- field: output_contract -->
#### Output contract

declared_by_signal

<!-- field: undo -->
#### Undo

session_state_restore

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| load_corpus | discover_suite_samples discovers evaluator sample workspaces; load_corpus loads project specs |

<!-- field: issues -->
#### Issues

- needs output.schema for discovered sample count and names

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: expand_eval_grid -->
### expand_eval_grid

<!-- field: role -->
#### Role

stateful_internal

<!-- field: selected_by -->
#### Selected by

- evaluator_session

<!-- field: inputs_from -->
#### Inputs from

- discover_suite_samples

<!-- field: feeds_into -->
#### Feeds into

- init_eval_session

<!-- field: side_effects -->
#### Side effects

eval_session_grid_mutation

<!-- field: output_contract -->
#### Output contract

declared_by_signal

<!-- field: undo -->
#### Undo

session_state_restore

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues -->
#### Issues

- needs output.schema for grid point count

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: init_eval_session -->
### init_eval_session

<!-- field: role -->
#### Role

stateful_internal

<!-- field: selected_by -->
#### Selected by

- evaluator_session

<!-- field: inputs_from -->
#### Inputs from

- expand_eval_grid
- tool_config

<!-- field: feeds_into -->
#### Feeds into

- report_suite_summary

<!-- field: side_effects -->
#### Side effects

filesystem_directory_create, eval_session_runtime_state_mutation

<!-- field: output_contract -->
#### Output contract

declared_by_signal

<!-- field: undo -->
#### Undo

session_state_restore_and_workspace_restore

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues -->
#### Issues

- creates session directory and should eventually participate in workspace rollback

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: report_suite_summary -->
### report_suite_summary

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- evaluator_session

<!-- field: inputs_from -->
#### Inputs from

- init_eval_session

<!-- field: feeds_into -->
#### Feeds into

- materialize_eval_points

<!-- field: side_effects -->
#### Side effects

stderr_write

<!-- field: output_contract -->
#### Output contract

declared_by_signal

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| report_session | report_suite_summary reports setup dimensions; report_session reports final outcomes |

<!-- field: issues -->
#### Issues

- needs output.schema for total point count

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: materialize_eval_points -->
### materialize_eval_points

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- evaluator_session

<!-- field: inputs_from -->
#### Inputs from

- report_suite_summary

<!-- field: feeds_into -->
#### Feeds into

- run_point

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

ordered evaluation point array

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: run_point -->
### run_point

<!-- field: role -->
#### Role

boundary

<!-- field: selected_by -->
#### Selected by

- evaluator_session

<!-- field: inputs_from -->
#### Inputs from

- materialize_eval_points

<!-- field: feeds_into -->
#### Feeds into

- run_point
- report_session

<!-- field: side_effects -->
#### Side effects

nested_machine_execution, point_result_recording, filesystem_mutation

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

boundary_compensation (child_command_undo, captures eval_session_snapshot + nested_machine_history)

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| invoke_executor | run_point invokes an in-process nested machine; invoke_executor uses self_invoke |

<!-- field: issues -->
#### Issues

- output.schema for point result summary

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: report_session -->
### report_session

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- evaluator_session

<!-- field: inputs_from -->
#### Inputs from

- run_point

<!-- field: feeds_into -->
#### Feeds into

- terminal

<!-- field: side_effects -->
#### Side effects

stderr_output

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| format_report | report_session summarizes eval run; format_report formats spec validation findings |

<!-- field: issues -->
#### Issues

- should declare whether stderr output is a side effect

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: create_point_dir -->
### create_point_dir

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- evaluator_point

<!-- field: inputs_from -->
#### Inputs from

- run_point

<!-- field: feeds_into -->
#### Feeds into

- copy_dir

<!-- field: side_effects -->
#### Side effects

filesystem_write, point_context_mutation

<!-- field: output_contract -->
#### Output contract

copy_dir source and destination parameters

<!-- field: undo -->
#### Undo

workspace_restore

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: sample_docs -->
### sample_docs

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- evaluator_point

<!-- field: inputs_from -->
#### Inputs from

- copy_dir

<!-- field: feeds_into -->
#### Feeds into

- copy_dir
- git_init

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

optional-docs presence and copy_dir parameters

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: commit_workspace_baseline -->
### commit_workspace_baseline

<!-- field: role -->
#### Role

boundary

<!-- field: selected_by -->
#### Selected by

- evaluator_point

<!-- field: inputs_from -->
#### Inputs from

- stage_all

<!-- field: feeds_into -->
#### Feeds into

- rev_parse

<!-- field: side_effects -->
#### Side effects

git_commit_create

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

workspace_restore

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| commit | commit_workspace_baseline is a profile-configured alias with deterministic local identity and message; commit accepts caller input |

<!-- field: issues -->
#### Issues

- needs output.schema for commit hash

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: record_agent_commit -->
### record_agent_commit

<!-- field: role -->
#### Role

stateful_internal

<!-- field: selected_by -->
#### Selected by

- evaluator_point

<!-- field: inputs_from -->
#### Inputs from

- rev_parse

<!-- field: feeds_into -->
#### Feeds into

- dump_config

<!-- field: side_effects -->
#### Side effects

point_context_commit_mutation

<!-- field: output_contract -->
#### Output contract

commit provenance string

<!-- field: undo -->
#### Undo

restore_point_context

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: dump_config -->
### dump_config

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- evaluator_point

<!-- field: inputs_from -->
#### Inputs from

- record_agent_commit
- point_context

<!-- field: feeds_into -->
#### Feeds into

- run_agent

<!-- field: side_effects -->
#### Side effects

filesystem_write

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

workspace_restore

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues -->
#### Issues

- needs output.schema for experiment.yaml path and bytes
- needs structured filesystem side effect

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: run_agent -->
### run_agent

<!-- field: role -->
#### Role

boundary

<!-- field: selected_by -->
#### Selected by

- evaluator_point

<!-- field: inputs_from -->
#### Inputs from

- dump_config
- point_context

<!-- field: feeds_into -->
#### Feeds into

- run_oracle_check

<!-- field: side_effects -->
#### Side effects

subprocess, workspace_mutation, trace_file_write

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

boundary_compensation (compensating_action, captures point_context_snapshot + point_workspace + trace_path + result_path)

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| invoke_executor | run_agent runs a suite harness binary; invoke_executor uses self_invoke for the planner executor |

<!-- field: issues -->
#### Issues

- output.schema for agent exit code and trace path

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: run_oracle_check -->
### run_oracle_check

<!-- field: role -->
#### Role

boundary

<!-- field: selected_by -->
#### Selected by

- evaluator_point

<!-- field: inputs_from -->
#### Inputs from

- run_agent

<!-- field: feeds_into -->
#### Feeds into

- record_oracle_result

<!-- field: side_effects -->
#### Side effects

configured_child_process

<!-- field: output_contract -->
#### Output contract

raw configured oracle output

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| test | run_oracle_check is a critic profile alias with profile-owned command arguments; generic test accepts a package parameter |

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: record_oracle_result -->
### record_oracle_result

<!-- field: role -->
#### Role

stateful_internal

<!-- field: selected_by -->
#### Selected by

- evaluator_point

<!-- field: inputs_from -->
#### Inputs from

- run_oracle_check

<!-- field: feeds_into -->
#### Feeds into

- collect_trace_tokens

<!-- field: side_effects -->
#### Side effects

point_context_test_result_mutation

<!-- field: output_contract -->
#### Output contract

configured oracle output

<!-- field: undo -->
#### Undo

restore_point_context

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: collect_trace_tokens -->
### collect_trace_tokens

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- evaluator_point

<!-- field: inputs_from -->
#### Inputs from

- record_oracle_result
- trace_ndjson

<!-- field: feeds_into -->
#### Feeds into

- check_agent_version

<!-- field: side_effects -->
#### Side effects

point_context_token_mutation

<!-- field: output_contract -->
#### Output contract

declared_by_signal

<!-- field: undo -->
#### Undo

restore_point_context

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| collect_metrics | collect_trace_tokens extracts token state; collect_metrics persists meta.json |

<!-- field: issues -->
#### Issues

- output.schema still needs full structured token summary

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: check_agent_version -->
### check_agent_version

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- evaluator_point

<!-- field: inputs_from -->
#### Inputs from

- collect_trace_tokens
- trace_ndjson
- harness_config

<!-- field: feeds_into -->
#### Feeds into

- summarize_point_results

<!-- field: side_effects -->
#### Side effects

point_context_version_mutation_and_warning_output

<!-- field: output_contract -->
#### Output contract

declared_by_signal

<!-- field: undo -->
#### Undo

restore_point_context

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues -->
#### Issues

- version mismatch is a warning signal that continues the point sequence

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: summarize_point_results -->
### summarize_point_results

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- evaluator_point

<!-- field: inputs_from -->
#### Inputs from

- record_oracle_result
- collect_trace_tokens
- check_agent_version

<!-- field: feeds_into -->
#### Feeds into

- collect_metrics

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| collect_metrics | summarize_point_results emits aggregate status; collect_metrics writes meta.json |

<!-- field: issues -->
#### Issues

- output.schema still needs full structured point summary

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: record_point_failure -->
### record_point_failure

<!-- field: role -->
#### Role

stateful_internal

<!-- field: selected_by -->
#### Selected by

- evaluator_point

<!-- field: inputs_from -->
#### Inputs from

- failed_point_command

<!-- field: feeds_into -->
#### Feeds into

- collect_metrics

<!-- field: side_effects -->
#### Side effects

point_context_mutation

<!-- field: output_contract -->
#### Output contract

failure stage and cause summary

<!-- field: undo -->
#### Undo

restore_point_context

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues | form: "empty list" -->
#### Issues

None

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: collect_metrics -->
### collect_metrics

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- evaluator_point

<!-- field: inputs_from -->
#### Inputs from

- summarize_point_results
- record_point_failure

<!-- field: feeds_into -->
#### Feeds into

- terminal

<!-- field: side_effects -->
#### Side effects

filesystem_write

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

workspace_restore

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| collect_trace_tokens | collect_metrics writes meta.json; collect_trace_tokens computes token state |
| summarize_point_results | collect_metrics persists metadata; summarize_point_results emits aggregate status |

<!-- field: issues -->
#### Issues

- needs output.schema for meta path and summary fields

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: load_corpus -->
### load_corpus

<!-- field: role -->
#### Role

stateful_internal

<!-- field: selected_by -->
#### Selected by

- jurist

<!-- field: inputs_from -->
#### Inputs from

- workspace_specs

<!-- field: feeds_into -->
#### Feeds into

- validate_specs

<!-- field: side_effects -->
#### Side effects

validation_state_mutation

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

restore_validation_state

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| parse_suite_config | load_corpus loads project specs; parse_suite_config loads evaluator suite metadata |
| discover_suite_samples | load_corpus loads project specs; discover_suite_samples finds evaluator sample workspaces |

<!-- field: issues -->
#### Issues

- needs output.schema for corpus counts

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: validate_specs -->
### validate_specs

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- jurist

<!-- field: inputs_from -->
#### Inputs from

- load_corpus

<!-- field: feeds_into -->
#### Feeds into

- format_report

<!-- field: side_effects -->
#### Side effects

validation_state_mutation

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

restore_validation_findings

<!-- field: issues -->
#### Issues

- mutates validation findings, so undo contract must be explicit

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: format_report -->
### format_report

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- jurist

<!-- field: inputs_from -->
#### Inputs from

- validate_specs

<!-- field: feeds_into -->
#### Feeds into

- terminal

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| report_session | format_report formats spec findings; report_session summarizes eval session |

<!-- field: issues -->
#### Issues

- needs output.schema for formatted report and error count

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: copy_dir -->
### copy_dir

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- evaluator_point

<!-- field: inputs_from -->
#### Inputs from

- create_point_dir
- tool_call

<!-- field: feeds_into -->
#### Feeds into

- sample_docs
- copy_dir
- git_init
- build
- test

<!-- field: side_effects -->
#### Side effects

filesystem_write

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

workspace_restore

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues -->
#### Issues

- legacy scalar side_effects

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: make_dir -->
### make_dir

<!-- field: role -->
#### Role

word

<!-- field: selected_by | form: "empty list" -->
#### Selected by

None

<!-- field: inputs_from -->
#### Inputs from

- tool_call

<!-- field: feeds_into -->
#### Feeds into

- write
- copy_dir

<!-- field: side_effects -->
#### Side effects

filesystem_write

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

workspace_restore

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues -->
#### Issues

- legacy scalar side_effects

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: git_init -->
### git_init

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- evaluator_point

<!-- field: inputs_from -->
#### Inputs from

- sample_docs
- copy_dir
- make_dir
- workspace

<!-- field: feeds_into -->
#### Feeds into

- stage_all
- commit

<!-- field: side_effects -->
#### Side effects

filesystem_write_git_metadata

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

workspace_restore

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues -->
#### Issues

- legacy scalar side_effects

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: stage_all -->
### stage_all

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- evaluator_point

<!-- field: inputs_from -->
#### Inputs from

- git_init
- workspace_changes

<!-- field: feeds_into -->
#### Feeds into

- commit
- commit_workspace_baseline

<!-- field: side_effects -->
#### Side effects

git_index_mutation

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

git_reset_index_or_workspace_restore

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues -->
#### Issues

- irreversible risk if used broadly without explicit selection

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: workspace_status -->
### workspace_status

<!-- field: role -->
#### Role

word

<!-- field: selected_by | form: "empty list" -->
#### Selected by

None

<!-- field: inputs_from -->
#### Inputs from

- workspace

<!-- field: feeds_into -->
#### Feeds into

- stage_all
- commit
- LLM_reasoning

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

declared

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| diff_stat | workspace_status lists porcelain file status; diff_stat summarizes diff size |

<!-- field: issues -->
#### Issues

- runtime output is not yet structured according to output.schema

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: commit -->
### commit

<!-- field: role -->
#### Role

boundary

<!-- field: selected_by | form: "empty list" -->
#### Selected by

None

<!-- field: inputs_from -->
#### Inputs from

- stage_all

<!-- field: feeds_into -->
#### Feeds into

- rev_parse
- log_oneline

<!-- field: side_effects -->
#### Side effects

git_commit_create

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

boundary_compensation (compensating_action, captures workspace_paths + git_ref)

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| issue_close | commit mutates git history; issue_close mutates bd issue state |

<!-- field: issues -->
#### Issues

- output.schema for commit hash

<!-- field: migration_priority -->
#### Migration priority

medium

<!-- tool: rev_parse -->
### rev_parse

<!-- field: role -->
#### Role

word

<!-- field: selected_by -->
#### Selected by

- evaluator_point

<!-- field: inputs_from -->
#### Inputs from

- commit
- commit_workspace_baseline
- workspace

<!-- field: feeds_into -->
#### Feeds into

- logs
- metadata
- record_agent_commit

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| log_oneline | rev_parse returns one ref; log_oneline returns commit history |

<!-- field: issues -->
#### Issues

- needs output.schema for ref string

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: diff_stat -->
### diff_stat

<!-- field: role -->
#### Role

word

<!-- field: selected_by | form: "empty list" -->
#### Selected by

None

<!-- field: inputs_from -->
#### Inputs from

- workspace_changes

<!-- field: feeds_into -->
#### Feeds into

- LLM_reasoning
- commit_decision

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| workspace_status | diff_stat summarizes diff size; workspace_status lists file statuses |

<!-- field: issues -->
#### Issues

- needs output.schema for file/insertion/deletion counts

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: log_oneline -->
### log_oneline

<!-- field: role -->
#### Role

word

<!-- field: selected_by | form: "empty list" -->
#### Selected by

None

<!-- field: inputs_from -->
#### Inputs from

- git_repo

<!-- field: feeds_into -->
#### Feeds into

- LLM_reasoning

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| rev_parse | log_oneline returns history; rev_parse returns one ref |

<!-- field: issues -->
#### Issues

- needs output.schema for commit list

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: issue_create -->
### issue_create

<!-- field: role -->
#### Role

boundary

<!-- field: selected_by | form: "empty list" -->
#### Selected by

None

<!-- field: inputs_from -->
#### Inputs from

- planner_output
- user_request

<!-- field: feeds_into -->
#### Feeds into

- issue_list
- issue_close

<!-- field: side_effects -->
#### Side effects

bd_database_write

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

boundary_compensation (compensating_action, captures created_issue_id + workspace_paths)

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| create_tracker_issue | issue_create is a shared direct bd CLI wrapper; create_tracker_issue is selected and configured by the planner profile |

<!-- field: issues -->
#### Issues

- output.schema for created issue ID

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: issue_close -->
### issue_close

<!-- field: role -->
#### Role

boundary

<!-- field: selected_by | form: "empty list" -->
#### Selected by

None

<!-- field: inputs_from -->
#### Inputs from

- issue_create
- user_request

<!-- field: feeds_into -->
#### Feeds into

- issue_list

<!-- field: side_effects -->
#### Side effects

bd_database_write

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

boundary_compensation (compensating_action, captures issue_id + previous_issue_status + workspace_paths)

<!-- field: overlaps | form: "empty list" -->
#### Overlaps

None

<!-- field: issues -->
#### Issues

- output.schema for closed issue confirmation

<!-- field: migration_priority -->
#### Migration priority

low

<!-- tool: issue_list -->
### issue_list

<!-- field: role -->
#### Role

word

<!-- field: selected_by | form: "empty list" -->
#### Selected by

None

<!-- field: inputs_from -->
#### Inputs from

- bd_database

<!-- field: feeds_into -->
#### Feeds into

- issue_create
- issue_close
- planner

<!-- field: side_effects -->
#### Side effects

none

<!-- field: output_contract -->
#### Output contract

missing

<!-- field: undo -->
#### Undo

noop

<!-- field: overlaps -->
#### Overlaps

| tool | difference |
| --- | --- |
| workspace_status | issue_list reports tracker state; workspace_status reports git state |

<!-- field: issues -->
#### Issues

- needs output.schema for issue summaries

<!-- field: migration_priority -->
#### Migration priority

low

<!-- field: sentence_tools | form: "tree" -->
## Sentence tools

- highest_priority: []

<!-- field: boundary_tools_with_compensation -->
## Boundary tools with compensation

<!-- field: note -->
### Note

All boundary tools now implement UndoMementoProvider and declare
payload: boundary_compensation with captures in their tool YAML.

<!-- field: tools -->
### Tools

- suspend
- checkpoint_rollback
- self_invoke
- create_tracker_issue
- invoke_executor
- run_point
- run_agent
- commit
- issue_create
- issue_close

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
