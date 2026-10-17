---
id: test-rel02.0-tool-contracts
title: Test Suite for Release 02.0 Tool Contracts
release: "02.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel02.0-uc003-tool-authoring
- rel02.0-uc004-sentence-tool-migration

<!-- field: preconditions -->
## Preconditions

- Tool declarations can be loaded from YAML into ToolDef values.
- MachineSpec validation can check selected action names and declared emitted signals.
- Contract validation can run in authoring review or future strict mode.
- Legacy sentence tools can be represented as deprecated compatibility declarations while replacement words are selected.

<!-- field: test_cases -->
## Test cases

<!-- case: Planner tracker sentence keeps CLI policy outside state adapters | status: done -->
### Planner tracker sentence keeps CLI policy outside state adapters

Format an implementation plan into parameters consumed by separate write
and exec words, then record a configured tracker's JSON result. Verify the
adapters contain no tracker invocation, update only planner state, and
restore that state on undo or malformed output.

<!-- field: use_case -->
#### Use case

rel02.0-uc004-sentence-tool-migration

<!-- field: go_test -->
#### Go test

TestFormatIssueEmitsTrackerAgnosticParameters, TestRecordTrackerIssueMapsConfiguredOutputAndUndo, TestRecordTrackerIssueRejectsInvalidOutputWithoutMutation

<!-- field: inputs | form: "tree" -->
#### Inputs

- current_plan: implementation_plan
- configured_tracker_result: "{\"id\":\"tracker-42\"}"

<!-- field: expected | form: "tree" -->
#### Expected

- tracker_neutral_parameters: true
- issue_id_recorded: tracker-42
- dependency_mapping_recorded: true
- invalid_output_preserves_state: true

<!-- field: traces -->
#### Traces

- srd013-standard-tool-library AC11

<!-- case: Declared provider address is an environment reference the loader expands | status: done -->
### Declared provider address is an environment reference the loader expands

Load a declaration whose provider_url is an environment reference and
verify the loader expands it before parsing. A set variable wins, an unset
one falls back to the declared local default so the shipped program runs
unconfigured, and $-prefixed selectors survive expansion unchanged.

The address matters because a chart co-generates the REST client base_url
but mounts tool declarations verbatim, so a literal loopback address is
correct on a developer machine and unreachable in a pod (GH-728).

<!-- field: use_case -->
#### Use case

rel02.0-uc003-tool-authoring

<!-- field: go_test -->
#### Go test

TestLoadToolDefsExpandsEnvReference, TestLoadToolDefsUsesDefaultWhenUnset, TestLoadToolDefsLeavesSelectorsUntouched

<!-- field: inputs | form: "tree" -->
#### Inputs

- declaration_field: provider_url
- declared_value: ${OLLAMA_URL:-http://localhost:11434}
- environment_value: http://ollama.svc:11434
- selectors_present:
  1. $from(compose).prompt
  1. $.output.answer

<!-- field: expected | form: "tree" -->
#### Expected

- expands_to_environment_value: true
- falls_back_to_declared_default: true
- selectors_unchanged: true

<!-- field: traces -->
#### Traces

- srd013-standard-tool-library AC10

<!-- case: Tool definition includes all required contract sections | status: done -->
### Tool definition includes all required contract sections

Load a new builtin ToolDef and verify it includes every required
authoring section named in cfg-tool-authoring-standard.

<!-- field: use_case -->
#### Use case

rel02.0-uc003-tool-authoring

<!-- field: go_test -->
#### Go test

TestValidate_SelectedToolContractCompletenessErrorsForActiveMigratedTool, TestValidate_SelectedToolContractCompletenessPassesCompleteTool

<!-- field: inputs | form: "tree" -->
#### Inputs

- tool_category: word
- declaration: complete_builtin_tooldef

<!-- field: expected | form: "tree" -->
#### Expected

- problem_present: true
- goals_present: true
- requirements_present: true
- non_goals_present: true
- output_schema_present: true
- side_effects_present: true
- reversibility_present: true
- undo_present: true
- errors_present: true
- relationships_present: true

<!-- field: traces -->
#### Traces

- srd022-builtin-tools AC1
- srd025-rollback-lifecycle AC6

<!-- case: One-verb rule rejects workflow-shaped ordinary tools | status: planned -->
### One-verb rule rejects workflow-shaped ordinary tools

Review a proposed ordinary word whose description performs multiple
observable verbs. The authoring check should reject it unless it is split
into separate words or classified as an explicit boundary.

No live test. The authoring reviewer that held this check had no caller
and was removed with GH-2071; the expectation stands for whatever
surfaces authoring review next.

<!-- field: use_case -->
#### Use case

rel02.0-uc003-tool-authoring

<!-- field: inputs | form: "tree" -->
#### Inputs

- proposed_tool: load config and discover samples and initialize session
- category: word

<!-- field: expected | form: "tree" -->
#### Expected

- authoring_result: rejected
- remediation: split_into_words_or_mark_boundary

<!-- field: traces -->
#### Traces

- srd022-builtin-tools AC1
- srd023-cli-exec-tools AC5

<!-- case: Typed input and output contracts are machine-consumable | status: done -->
### Typed input and output contracts are machine-consumable

Load a tool declaration with JSON Schema parameters and output.schema.
Verify parameter and output contracts are precise enough for another
tool or machine transition to consume without prose interpretation.

<!-- field: use_case -->
#### Use case

rel02.0-uc003-tool-authoring

<!-- field: go_test -->
#### Go test

TestParseToolDefs_ContractFields, TestValidate_SelectedToolContractCompletenessPassesCompleteTool

<!-- field: inputs | form: "tree" -->
#### Inputs

- parameters_schema: object_with_required_fields
- output_schema: object_with_result_fields

<!-- field: expected | form: "tree" -->
#### Expected

- parameters_schema_valid: true
- output_schema_valid: true
- llm_visible_schema_excludes_implementation_extensions: true

<!-- field: traces -->
#### Traces

- srd022-builtin-tools AC1
- srd023-cli-exec-tools AC2

<!-- case: Side effects and undo strategy are required for mutating tools | status: done -->
### Side effects and undo strategy are required for mutating tools

Check a filesystem-mutating tool and a command/domain-state-mutating
tool. Both must declare their rollback contract and
state_after_failure error semantics.

<!-- field: use_case -->
#### Use case

rel02.0-uc003-tool-authoring

<!-- field: go_test -->
#### Go test

TestSignatureNeverWaivesBoundarySideEffects, TestValidate_ToolUndoConsistency

<!-- field: inputs | form: "tree" -->
#### Inputs

- tools:
  1. filesystem_write_word
  1. conversation_state_word

<!-- field: expected | form: "tree" -->
#### Expected

- filesystem_tool_declares_workspace_restore: true
- state_tool_declares_command_state_restore: true
- errors_include_state_after_failure: true

<!-- field: traces -->
#### Traces

- srd025-rollback-lifecycle AC6

<!-- case: Irreversible tools require confirmation or compensation | status: planned -->
### Irreversible tools require confirmation or compensation

Assess a boundary tool that affects external user data. If its effects
are irreversible or compensatable, the ToolDef must declare confirmation
or a compensating action before it is selected by a user-facing machine.

No live test. The authoring reviewer that held this check had no caller
and was removed with GH-2071; the expectation stands for whatever
surfaces authoring review next.

<!-- field: use_case -->
#### Use case

rel02.0-uc003-tool-authoring

<!-- field: inputs | form: "tree" -->
#### Inputs

- tool_category: boundary
- reversibility: irreversible
- side_effect_target: external_service

<!-- field: expected | form: "tree" -->
#### Expected

- requires_confirmation: true
- compensation_or_irreversible_reason_documented: true

<!-- field: traces -->
#### Traces

- srd024-blocking-input-tools AC1
- srd025-rollback-lifecycle AC6

<!-- case: Relationships distinguish overlapping vocabulary words | status: planned -->
### Relationships distinguish overlapping vocabulary words

Review a new tool that overlaps with existing words. The ToolDef must
document likely before/after neighbors and overlap differences so agents
can choose the correct word.

No live test. The authoring reviewer that held this check had no caller
and was removed with GH-2071; the expectation stands for whatever
surfaces authoring review next.

<!-- field: use_case -->
#### Use case

rel02.0-uc003-tool-authoring

<!-- field: inputs | form: "tree" -->
#### Inputs

- new_tool: collect_artifact_metadata
- overlaps_with:
  1. collect_metrics
  1. dump_config

<!-- field: expected | form: "tree" -->
#### Expected

- relationships_before_present: true
- relationships_after_present: true
- overlaps_explain_differences: true

<!-- field: traces -->
#### Traces

- srd022-builtin-tools AC1

<!-- case: Sentence tool migration creates replacement word declarations | status: done -->
### Sentence tool migration creates replacement word declarations

Given a legacy aggregate tool, identify each observable verb and create
one ToolDef per word. Replacement declarations must include the
required contract metadata.

The declarations are core-owned, so the proof reads them from
tools/builtin. That a shipped profile selects them is proven in
applications/catalog by test-rel07.1-profile-conformance (GH-512).

<!-- field: use_case -->
#### Use case

rel02.0-uc004-sentence-tool-migration

<!-- field: go_test -->
#### Go test

TestSessionWordsDeclareFullContract, TestPointWordsDeclareFullContract, TestSessionWordsReplaceLoadSuite

<!-- field: inputs | form: "tree" -->
#### Inputs

- legacy_tool: load_suite
- identified_words:
  1. parse_suite_config
  1. discover_suite_samples
  1. expand_eval_grid
  1. init_eval_session
  1. report_suite_summary

<!-- field: expected | form: "tree" -->
#### Expected

- replacement_tooldefs_created: true
- each_replacement_has_contract_metadata: true
- aggregate_not_selected_by_default: true

<!-- field: traces -->
#### Traces

- srd013-standard-tool-library AC9
- srd022-builtin-tools AC1
- srd025-rollback-lifecycle AC6

<!-- case: MachineSpec owns migrated workflow sequence | status: done -->
### MachineSpec owns migrated workflow sequence

Rewire a migrated workflow so machine.yaml has one state/action step per
replacement word. Startup validation should confirm each named action is
registered and each emitted signal has a transition.

This is the reusable contract, so it is proven on a synthetic
session-shaped machine under internal/evaluation/testdata rather than on a
shipped profile, including the rejections: an emit the machine never
declares, and a selected word no declaration provides (GH-512).

<!-- field: use_case -->
#### Use case

rel02.0-uc004-sentence-tool-migration

<!-- field: go_test -->
#### Go test

TestSessionPipelineEmitsValidate, TestUndeclaredEmitIsRejected, TestSelectionRejectsUndeclaredWord

<!-- field: inputs | form: "tree" -->
#### Inputs

- machine: internal/evaluation/testdata/session_machine.yaml
- replacement_words:
  1. load_input
  1. prepare_work
  1. report_result

<!-- field: expected | form: "tree" -->
#### Expected

- named_actions_registered: true
- emitted_signals_declared: true
- emitted_signals_have_transitions: true
- sequence_visible_in_machine_yaml: true

<!-- field: traces -->
#### Traces

- srd006-yaml-machine-config AC3
- srd022-builtin-tools AC4

<!-- case: Boundary and nested-machine tools declare their boundary effects | status: done -->
### Boundary and nested-machine tools declare their boundary effects

Exercise a tool that invokes a child process, nested machine, model, or
human actor. The declaration must classify it as boundary and describe
the external state and workspace effects rollback may need to compensate.

<!-- field: use_case -->
#### Use case

rel02.0-uc004-sentence-tool-migration

<!-- field: go_test -->
#### Go test

TestValidate_ToolBoundaryCategory

<!-- field: inputs | form: "tree" -->
#### Inputs

- boundary_tools:
  1. run_agent
  1. rest_await_event
  1. self_invoke
  1. invoke_llm

<!-- field: expected | form: "tree" -->
#### Expected

- category_boundary: true
- emitted_signals_declared: true
- child_or_external_effects_declared: true
- undo_or_compensation_declared: true

<!-- field: traces -->
#### Traces

- srd024-blocking-input-tools AC1
- srd025-rollback-lifecycle AC6

<!-- case: Deprecated compatibility aggregate is retired | status: done -->
### Deprecated compatibility aggregate is retired

Remove the legacy aggregate implementation, declaration, and factory
registration from Agent Core.

<!-- field: use_case -->
#### Use case

rel02.0-uc004-sentence-tool-migration

<!-- field: go_test -->
#### Go test

TestBuiltinFactoryCatalogCoversSelectedActiveInits

<!-- field: inputs | form: "tree" -->
#### Inputs

- retired_aggregate: validate
- replacement_sequence:
  1. build
  1. lint
  1. test

<!-- field: expected | form: "tree" -->
#### Expected

- implementation_removed: true
- declaration_removed: true
- factory_registration_removed: true
- validation_uses_normal_machine_dispatch: true

<!-- field: traces -->
#### Traces

- srd013-standard-tool-library AC8
- srd022-builtin-tools AC1

<!-- case: Migration preserves intermediate failure signals | status: done -->
### Migration preserves intermediate failure signals

Split a workflow-shaped tool into multiple words and verify the migrated
machine exposes intermediate failure transitions that the aggregate had
hidden or collapsed.

Core proves the property: every stage word declares a failure emit, and
every failure emit has a transition. The shipped critic point sequence that
exercises it is proven in applications/catalog by
test-rel07.1-profile-conformance (GH-512).

<!-- field: use_case -->
#### Use case

rel02.0-uc004-sentence-tool-migration

<!-- field: go_test -->
#### Go test

TestSessionPipelineFailureSignalsRoute, TestPointWordsDeclareFullContract

<!-- field: inputs | form: "tree" -->
#### Inputs

- migrated_sequence:
  1. load_input
  1. prepare_work
  1. report_result
- point_words:
  1. record_oracle_result
  1. collect_trace_tokens
  1. check_agent_version
  1. summarize_point_results
  1. record_agent_commit
- failure_points:
  1. step_failed

<!-- field: expected | form: "tree" -->
#### Expected

- each_stage_failure_transition: explicit
- point_words_declare_failure_emits: true
- failure_routes_to_terminal: true

<!-- field: traces -->
#### Traces

- srd006-yaml-machine-config AC3
- srd013-standard-tool-library AC9

<!-- case: Lifecycle builtin declarations expose full adapter contracts | status: planned -->
### Lifecycle builtin declarations expose full adapter contracts

Inspect the checkpoint_history and checkpoint_rollback ToolDefs. Both
declarations should expose the required adapter contract fields and the
correct builtin init value.

Partial proof only. TestBuiltinFactoryCatalogSelectsEntriesByInit asserts
the tools are builtin adapters selected by their init (type_builtin,
init_present); TestValidate_SelectedToolContractCompletenessKeepsLegacyWarningOnly
holds a legacy contract to warning level while a migrated one errors. No current test asserts the full
nine-field adapter contract set for these ToolDefs, so the remaining
field-presence checks are listed under planned rather than expected.

<!-- field: planned_note -->
#### Planned note

Claimed by this case but not asserted by the go_test above for the checkpoint_history and checkpoint_rollback ToolDefs; planned, not proven.

<!-- field: use_case -->
#### Use case

rel02.0-uc003-tool-authoring

<!-- field: go_test -->
#### Go test

TestBuiltinFactoryCatalogSelectsEntriesByInit, TestValidate_SelectedToolContractCompletenessKeepsLegacyWarningOnly

<!-- field: inputs | form: "tree" -->
#### Inputs

- lifecycle_tools:
  1. checkpoint_history
  1. checkpoint_rollback

<!-- field: expected | form: "tree" -->
#### Expected

- type_builtin: true
- init_present: true

<!-- field: planned -->
#### Planned

- emits_present
- output_present
- side_effects_present
- reversibility_present
- undo_present
- errors_present
- relationships_present

<!-- field: traces -->
#### Traces

- srd022-builtin-tools AC1
- srd025-rollback-lifecycle AC6
- srd026-lifecycle-tools AC5

<!-- case: Lifecycle tool contracts reject exec-wrapper behavior | status: done -->
### Lifecycle tool contracts reject exec-wrapper behavior

Check lifecycle ToolDefs for hidden shelling into the agent binary. Desired
declarations should be builtin adapters, selected by lifecycle
profiles, and invoked by MachineSpec actions.

<!-- field: use_case -->
#### Use case

rel02.0-uc003-tool-authoring

<!-- field: go_test -->
#### Go test

TestRootCommandHasNoLifecycleSubcommands, TestBuiltinFactoryCatalogSelectsEntriesByInit

<!-- field: inputs | form: "tree" -->
#### Inputs

- lifecycle_tools:
  1. checkpoint_history
  1. checkpoint_rollback
- forbidden_shape: exec_wrapper_around_agent_binary

<!-- field: expected | form: "tree" -->
#### Expected

- shells_into_agent_binary: false
- selected_by_lifecycle_profile: true
- machine_action_visible: true
- core_operation_boundary_documented: true

<!-- field: traces -->
#### Traces

- srd025-rollback-lifecycle AC7
- srd026-lifecycle-tools AC1
- srd026-lifecycle-tools AC5
