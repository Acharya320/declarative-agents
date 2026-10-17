---
id: test-rel01.0
title: Test Suite for Release 01.0 Unified Agent Runtime
release: "01.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel01.0-uc001-generator-coding
- rel01.0-uc002-evaluator-benchmark
- rel01.0-uc003-bench-visualization

<!-- field: preconditions -->
## Preconditions

- Go toolchain installed and repository packages buildable.
- bin/agent is buildable through mage build or the integration target helper.
- Ollama is available for live generator and evaluator integration targets.
- Coding workflow proofs run from applications/coding-agent.
- Bench evaluator proofs run from applications/catalog.

<!-- field: test_cases -->
## Test cases

<!-- case: list_files executes bounded find within the workspace | status: implemented -->
### list_files executes bounded find within the workspace

<!-- field: use_case -->
#### Use case

rel01.0-uc001-generator-coding

<!-- field: go_test -->
#### Go test

TestListFilesDeclarationUsesBoundedFindExec, TestListFilesExecPreservesTreeContract, TestListFilesExecHonorsDepthAndSubpath, TestListFilesExecRejectsMissingAndEscapingPaths, TestListFilesExecCapsOutput

<!-- field: inputs | form: "tree" -->
#### Inputs

- declaration: tools/builtin/list-files.yaml
- executable_dependency: GNU find (findutils)
- cases:
  1. hidden_entries
  1. symlinks
  1. deterministic_order
  1. depth_and_subpath
  1. missing_and_escaping_paths
  1. output_cap

<!-- field: expected | form: "tree" -->
#### Expected

- production_go_walk: absent
- workspace_escape: rejected
- symlinks_followed: false
- output_line_cap: 200

<!-- field: traces -->
#### Traces

- srd013-standard-tool-library AC1
- srd023-cli-exec-tools AC2
- srd023-cli-exec-tools AC3

<!-- case: Qwen profile strips thinking blocks and extracts tool envelope -->
### Qwen profile strips thinking blocks and extracts tool envelope

<!-- field: use_case -->
#### Use case

rel01.0-uc001-generator-coding

<!-- field: go_test -->
#### Go test

applications/coding-agent mage integration:executorLive

<!-- field: inputs | form: "tree" -->
#### Inputs

- profile: ${AGENT_CATALOG_ROOT}/agents/executor/profile.yaml
- sample: testdata/integration/uc001-generator-coding
- model: qwen3.6:35b-mlx

<!-- field: expected | form: "tree" -->
#### Expected

- qwen_profile_selected: true
- thinking_blocks_stripped: true
- tool_call_envelope_extracted: true

<!-- field: traces -->
#### Traces

- srd009-ollama-adapter AC1
- srd017-llm-extensions AC1
- srd017-llm-extensions AC2
- srd017-llm-extensions AC3

<!-- case: Ollama adapter construction performs no readiness I/O -->
### Ollama adapter construction performs no readiness I/O

<!-- field: use_case -->
#### Use case

rel01.0-uc001-generator-coding

<!-- field: go_test -->
#### Go test

TestInvokeLLMBuildsWithoutNetworkIO

<!-- field: inputs | form: "tree" -->
#### Inputs

- provider_url: unserved loopback endpoint

<!-- field: expected | form: "tree" -->
#### Expected

- construction_succeeds: true
- startup_network_requests: 0

<!-- field: traces -->
#### Traces

- srd009-ollama-adapter AC2

<!-- case: Ollama adapter emits tracing around chat -->
### Ollama adapter emits tracing around chat

<!-- field: use_case -->
#### Use case

rel01.0-uc001-generator-coding

<!-- field: go_test -->
#### Go test

TestInvokeLLMRunsTheDialectPipeline

<!-- field: inputs | form: "tree" -->
#### Inputs

- tracer: recording tracer
- chat_request: non-streaming /api/chat

<!-- field: expected | form: "tree" -->
#### Expected

- inference_span_attrs: true
- nil_tracer_safe: true

<!-- field: traces -->
#### Traces

- srd009-ollama-adapter AC3

<!-- case: Ollama adapter reports HTTP and JSON failures -->
### Ollama adapter reports HTTP and JSON failures

<!-- field: use_case -->
#### Use case

rel01.0-uc001-generator-coding

<!-- field: go_test -->
#### Go test

TestProviderDialectFailureKeepsCommandError

<!-- field: inputs | form: "tree" -->
#### Inputs

- failures:
  1. connection_refused
  1. bad_http_status
  1. malformed_json
  1. cancelled_context

<!-- field: expected | form: "tree" -->
#### Expected

- distinct_errors: true
- no_retry: true

<!-- field: traces -->
#### Traces

- srd009-ollama-adapter AC4

<!-- case: Generator solves sample and passes validation -->
### Generator solves sample and passes validation

<!-- field: use_case -->
#### Use case

rel01.0-uc001-generator-coding

<!-- field: go_test -->
#### Go test

applications/coding-agent mage integration:executorLive

<!-- field: inputs | form: "tree" -->
#### Inputs

- profile: ${AGENT_CATALOG_ROOT}/agents/executor/profile.yaml
- workspace: testdata/integration/uc001-generator-coding
- validation: build_lint_test

<!-- field: expected | form: "tree" -->
#### Expected

- run_status: succeeded
- oracle_tests_pass: true
- workspace_changes_validated: true

<!-- field: traces -->
#### Traces

- srd005-agentic-loop AC5
- srd013-standard-tool-library AC2

<!-- case: Evaluator setup tools expand suite grid and initialize session -->
### Evaluator setup tools expand suite grid and initialize session

<!-- field: use_case -->
#### Use case

rel01.0-uc002-evaluator-benchmark

<!-- field: go_test -->
#### Go test

TestEvaluatorSessionSetupToolSequence

<!-- field: inputs | form: "tree" -->
#### Inputs

- suite: synthetic_profile_suite
- grid: effort_low_high
- repetitions: 2

<!-- field: expected | form: "tree" -->
#### Expected

- suite_parsed: true
- samples_discovered: true
- grid_points: 2
- reported_points: 4

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC0
- srd019-eval-harness AC1

<!-- case: Evaluator session defaults are declared and grid expansion is mandatory -->
### Evaluator session defaults are declared and grid expansion is mandatory

<!-- field: use_case -->
#### Use case

rel01.0-uc002-evaluator-benchmark

<!-- field: go_test -->
#### Go test

go test ./internal/evaluation -run 'Test(InitEvalSessionRequiresExpandedGrid|InitEvalSessionDefaultsComeFromDeclaration)'

<!-- field: inputs | form: "tree" -->
#### Inputs

- init_eval_session_config:
  - output_dir: eval-results
  - reps: 1
  - timeout: 600
- grid_points: 0

<!-- field: expected | form: "tree" -->
#### Expected

- defaults_loaded_from_tooldef: true
- missing_grid_signal: CommandError
- session_directory_created_on_failure: false

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC0
- srd019-eval-harness AC1

<!-- case: Evaluator rejects invalid point and aggregate duration envelopes | status: implemented -->
### Evaluator rejects invalid point and aggregate duration envelopes

<!-- field: use_case -->
#### Use case

rel01.0-uc002-evaluator-benchmark

<!-- field: go_test -->
#### Go test

go test ./internal/evaluation -run 'TestInitEvalSessionRejects(InvalidPointTimeout|ExcessiveAggregateDuration|AggregateDurationOverflow)$'

<!-- field: inputs | form: "tree" -->
#### Inputs

- point_timeout_authority: 15m
- aggregate_session_authority: 24h
- point_timeout_cases:
  1. malformed
  1. zero
  1. negative
  1. above_authority
- aggregate_cases:
  1. at_authority
  1. above_authority
  1. multiplication_overflow

<!-- field: expected | form: "tree" -->
#### Expected

- malformed_and_non_positive_rejected: true
- point_timeout_above_15m_rejected: true
- aggregate_duration_above_24h_rejected: true
- overflow_rejected: true
- session_directory_created_on_failure: false

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC5

<!-- case: Nested evaluator dispatch uses the point machine command timeout | status: implemented -->
### Nested evaluator dispatch uses the point machine command timeout

<!-- field: use_case -->
#### Use case

rel01.0-uc002-evaluator-benchmark

<!-- field: go_test -->
#### Go test

go test ./internal/evaluation -run TestRunPointPropagatesMachineCommandTimeout$

<!-- field: inputs | form: "tree" -->
#### Inputs

- point_machine: synthetic machine with a finite command_timeout
- point_command: context-aware blocking command

<!-- field: expected | form: "tree" -->
#### Expected

- nested_loop_uses_declared_command_timeout: true
- deadline_observed_without_wall_clock_sleep: true
- omitted_or_copied_fallback_timeout: false

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC5

<!-- case: Critic and bench profiles reserve nested evaluator timeout margins | status: implemented -->
### Critic and bench profiles reserve nested evaluator timeout margins

<!-- field: use_case -->
#### Use case

rel01.0-uc002-evaluator-benchmark

<!-- field: go_test -->
#### Go test

applications/catalog go test ./magefiles -run TestCatalogEvaluatorTimeoutEnvelopes$

<!-- field: inputs | form: "tree" -->
#### Inputs

- point_runtime_authority: 15m
- aggregate_session_authority: 24h
- critic_profiles:
  1. agents/critic/point.yaml
  1. agents/critic/machine.yaml
- bench_profile: agents/bench/machine.yaml

<!-- field: expected | form: "tree" -->
#### Expected

- critic_point_command_margin_at_least: 1m
- critic_run_point_dispatch_margin_at_least: 1m
- bench_child_session_margin_at_least: 1h
- all_bounds_finite_and_profile_owned: true

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC5

<!-- case: Evaluator receipts remove only owned artifacts -->
### Evaluator receipts remove only owned artifacts

<!-- field: use_case -->
#### Use case

rel01.0-uc002-evaluator-benchmark

<!-- field: go_test -->
#### Go test

go test ./internal/evaluation -run 'Test(EvaluatorReceiptRejectsArtifactOutsideOwnedRoot|CollectMetricsRejectsMissingPointDirectory)'

<!-- field: inputs | form: "tree" -->
#### Inputs

- point_dir: ""
- tampered_remove_path: outside owned root

<!-- field: expected | form: "tree" -->
#### Expected

- empty_point_dir_rejected: true
- outside_artifact_retained: true
- undo_signal: CommandError

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC2

<!-- case: Evaluator profile suite iterates generator profiles -->
### Evaluator profile suite iterates generator profiles

<!-- field: use_case -->
#### Use case

rel01.0-uc002-evaluator-benchmark

<!-- field: go_test -->
#### Go test

TestMaterializeEvalPointsPreservesProfileOrder

<!-- field: inputs | form: "tree" -->
#### Inputs

- profiles:
  1. agents/executor/profile-qwen35b.yaml
  1. agents/executor/profile-qwen27b.yaml
- samples: discovered_samples

<!-- field: expected | form: "tree" -->
#### Expected

- first_profile_selected: true
- second_profile_selected: true
- all_points_done_signal: true

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC1
- srd019-eval-harness AC2

<!-- case: Evaluator sample layout is declaration-driven -->
### Evaluator sample layout is declaration-driven

<!-- field: use_case -->
#### Use case

rel01.0-uc002-evaluator-benchmark

<!-- field: go_test -->
#### Go test

TestDiscoverSamplesUsesDeclaredLayout

<!-- field: inputs | form: "tree" -->
#### Inputs

- workspace_dir: project
- doc_dir: guidance
- prompt_file: task.txt

<!-- field: expected | form: "tree" -->
#### Expected

- custom_layout_discovered: true
- shared_prompt_resolved: true

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC0

<!-- case: Convergence accepts any structured metric tool name -->
### Convergence accepts any structured metric tool name

<!-- field: use_case -->
#### Use case

rel01.0-uc002-evaluator-benchmark

<!-- field: go_test -->
#### Go test

TestExtractToolSnapshots_WithMetrics

<!-- field: inputs | form: "tree" -->
#### Inputs

- tool_name: run_checks
- metrics:
  - total: 4
  - passed: 3
  - failed: 1

<!-- field: expected | form: "tree" -->
#### Expected

- snapshot_preserved: true
- go_output_regex_required: false

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC3

<!-- case: Evaluator integration writes result directory -->
### Evaluator integration writes result directory

<!-- field: use_case -->
#### Use case

rel01.0-uc002-evaluator-benchmark

<!-- field: go_test -->
#### Go test

applications/coding-agent mage integration:criticGate

<!-- field: inputs | form: "tree" -->
#### Inputs

- profile: ${AGENT_CATALOG_ROOT}/agents/critic/profile.yaml
- suite: testdata/integration/uc002-evaluator-benchmark/suite.yaml
- output: eval-results

<!-- field: expected | form: "tree" -->
#### Expected

- generator_profiles_run: true
- result_entries_created: true
- meta_and_experiment_files_available: true

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC1
- srd019-eval-harness AC2
- srd016-traceparent AC1

<!-- case: Point result tools classify output and collect metrics -->
### Point result tools classify output and collect metrics

<!-- field: use_case -->
#### Use case

rel01.0-uc002-evaluator-benchmark

<!-- field: go_test -->
#### Go test

TestSummarizePointResultsUsesCollectedState

<!-- field: inputs | form: "tree" -->
#### Inputs

- point_context: collected_trace_and_oracle_state
- metrics: duration_and_outcome

<!-- field: expected | form: "tree" -->
#### Expected

- summary_uses_collected_state: true
- outcome_classified: true
- metrics_available_on_failure: true

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC2
- srd019-eval-harness AC4

<!-- case: Failed point preserves diagnostic metadata | status: implemented -->
### Failed point preserves diagnostic metadata

<!-- field: use_case -->
#### Use case

rel01.0-uc002-evaluator-benchmark

<!-- field: go_test -->
#### Go test

TestPointMachineRoutesFailureThroughMetricsBeforeRunPointReturns

<!-- field: inputs | form: "tree" -->
#### Inputs

- point_machine: deterministic failure before collect_metrics
- release_results: mixed passing and failed points

<!-- field: expected | form: "tree" -->
#### Expected

- failure_metadata_preserved: true
- point_stage_and_cause_reported: true
- failed_point_rejects_gate: true
- nested_dispatch_span_present: true

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC2
- srd019-eval-harness AC3

<!-- case: Evaluation artifact queries preserve bench data -->
### Evaluation artifact queries preserve bench data

<!-- field: use_case -->
#### Use case

rel01.0-uc003-bench-visualization

<!-- field: go_test -->
#### Go test

TestEvaluationArtifactQueriesPreserveBenchData

<!-- field: inputs | form: "tree" -->
#### Inputs

- tools:
  1. list_evaluation_sessions
  1. analyze_evaluation_session
  1. list_evaluation_points
  1. read_evaluation_trace
- data: eval-results

<!-- field: expected | form: "tree" -->
#### Expected

- http_independent_statistics: true
- traversal_rejected: true

<!-- field: traces -->
#### Traces

- srd020-bench-web-ui AC4
- srd021-agent-cli-exec-tool AC2
- srd022-builtin-tools AC6

<!-- case: Bench integration target records manual UI preconditions -->
### Bench integration target records manual UI preconditions

<!-- field: use_case -->
#### Use case

rel01.0-uc003-bench-visualization

<!-- field: go_test -->
#### Go test

applications/catalog mage integration:benchEvaluator

<!-- field: inputs | form: "tree" -->
#### Inputs

- eval_results_dir: required
- free_port: required

<!-- field: expected | form: "tree" -->
#### Expected

- current_target_status: skipped
- skip_reason_documents_manual_preconditions: false
- skip_reason_documents_browser_preconditions: true
- no_automated_ui_success_claim: true

<!-- field: traces -->
#### Traces

- srd020-bench-web-ui AC3
- srd024-blocking-input-tools AC1

<!-- case: Evaluator model statistics aggregate run outcomes -->
### Evaluator model statistics aggregate run outcomes

<!-- field: use_case -->
#### Use case

rel01.0-uc002-evaluator-benchmark

<!-- field: go_test -->
#### Go test

TestComputeModelStatsAggregatesRuns

<!-- field: inputs | form: "tree" -->
#### Inputs

- grouped_runs: passed and failed runs for one model

<!-- field: expected | form: "tree" -->
#### Expected

- success_rate_computed: true
- convergence_rates_computed: true
- mean_tokens_and_duration_computed: true

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC3

<!-- case: Evaluator detailed rows aggregate sample and model results -->
### Evaluator detailed rows aggregate sample and model results

<!-- field: use_case -->
#### Use case

rel01.0-uc002-evaluator-benchmark

<!-- field: go_test -->
#### Go test

TestComputeDetailedAggregatesSampleModelRows

<!-- field: inputs | form: "tree" -->
#### Inputs

- grouped_runs: one sample and model result group

<!-- field: expected | form: "tree" -->
#### Expected

- sample_model_row_created: true
- convergence_counts_computed: true
- means_computed: true

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC3

<!-- case: Bench session APIs return evaluation artifacts -->
### Bench session APIs return evaluation artifacts

<!-- field: use_case -->
#### Use case

rel01.0-uc003-bench-visualization

<!-- field: go_test -->
#### Go test

TestEvaluationArtifactQueriesPreserveBenchData

<!-- field: inputs | form: "tree" -->
#### Inputs

- eval_results: meta.json, trace.ndjson, experiment.yaml

<!-- field: expected | form: "tree" -->
#### Expected

- sessions_returned: true
- points_returned: true
- trace_returned: true
- experiment_returned: true

<!-- field: traces -->
#### Traces

- srd020-bench-web-ui AC2

<!-- case: Agent child process args stay profile-first -->
### Agent child process args stay profile-first

<!-- field: use_case -->
#### Use case

rel01.0-uc002-evaluator-benchmark

<!-- field: go_test -->
#### Go test

TestBuildArgs_ProfileOnly

<!-- field: inputs | form: "tree" -->
#### Inputs

- profile: ${AGENT_CATALOG_ROOT}/agents/executor/profile.yaml

<!-- field: expected | form: "tree" -->
#### Expected

- profile_flag_present: true
- machine_and_model_flags_absent: true

<!-- field: traces -->
#### Traces

- srd021-agent-cli-exec-tool AC1

<!-- case: Self invoke maps child process success and failure -->
### Self invoke maps child process success and failure

<!-- field: use_case -->
#### Use case

rel01.0-uc002-evaluator-benchmark

<!-- field: go_test -->
#### Go test

go test ./internal/tools/control -run 'TestSelfInvokeBuilder_Execute_(Success|Failure)$'

<!-- field: inputs | form: "tree" -->
#### Inputs

- child_binary: echo
- failure_binary: false

<!-- field: expected | form: "tree" -->
#### Expected

- success_maps_to_tool_done: true
- failure_maps_to_tool_failed: true
- duration_recorded: true

<!-- field: traces -->
#### Traces

- srd021-agent-cli-exec-tool AC4

<!-- case: invoke_llm handles call timeout as command error -->
### invoke_llm handles call timeout as command error

<!-- field: use_case -->
#### Use case

rel01.0-uc001-generator-coding

<!-- field: go_test -->
#### Go test

TestInvokeLLM_CallTimeout

<!-- field: inputs | form: "tree" -->
#### Inputs

- client: waits for context cancellation
- call_timeout: 1ms

<!-- field: expected | form: "tree" -->
#### Expected

- deadline_error_returned: true
- command_error_signal: true
- duration_recorded: true

<!-- field: traces -->
#### Traces

- srd024-blocking-input-tools AC3

<!-- case: Bench action endpoint applies single pending action backpressure -->
### Bench action endpoint applies single pending action backpressure

<!-- field: use_case -->
#### Use case

rel01.0-uc003-bench-visualization

<!-- field: go_test -->
#### Go test

TestRESTServer_QueueOverflowPolicies

<!-- field: inputs | form: "tree" -->
#### Inputs

- action_queue_capacity: 1

<!-- field: expected | form: "tree" -->
#### Expected

- first_action_accepted: true
- second_action_rejected: true

<!-- field: traces -->
#### Traces

- srd024-blocking-input-tools AC4

<!-- case: invoke_llm records duration and token usage -->
### invoke_llm records duration and token usage

<!-- field: use_case -->
#### Use case

rel01.0-uc001-generator-coding

<!-- field: go_test -->
#### Go test

TestInvokeLLM_Success

<!-- field: inputs | form: "tree" -->
#### Inputs

- chat_response: content plus token counts

<!-- field: expected | form: "tree" -->
#### Expected

- duration_recorded: true
- token_usage_recorded: true
- result_signal_returned: true

<!-- field: traces -->
#### Traces

- srd024-blocking-input-tools AC5

<!-- case: Builtin registry resolves factories and rejects unknown init -->
### Builtin registry resolves factories and rejects unknown init

<!-- field: use_case -->
#### Use case

rel01.0-uc001-generator-coding

<!-- field: go_test -->
#### Go test

go test ./internal/tools/registry -run 'Test(BuiltinRegistryRegisterResolveOverride|RegisterUnifiedToolsUnknownInit|RegisterSingleBuiltinOverrides)'

<!-- field: inputs | form: "tree" -->
#### Inputs

- builtin_defs: registered and missing init names

<!-- field: expected | form: "tree" -->
#### Expected

- init_resolved: true
- unknown_init_rejected: true
- registration_overrides_existing_tool: true

<!-- field: traces -->
#### Traces

- srd022-builtin-tools AC2

<!-- case: Builtin tool config decodes typed structs and errors -->
### Builtin tool config decodes typed structs and errors

<!-- field: use_case -->
#### Use case

rel01.0-uc001-generator-coding

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog -run 'TestDecodeToolConfig'

<!-- field: inputs | form: "tree" -->
#### Inputs

- configs: child agent, LLM, load suite, run point, evaluation artifacts

<!-- field: expected | form: "tree" -->
#### Expected

- typed_configs_decoded: true
- type_mismatch_mentions_tool: true
- extra_fields_ignored: true

<!-- field: traces -->
#### Traces

- srd022-builtin-tools AC3

<!-- case: Exec declarations parse schema and contract fields -->
### Exec declarations parse schema and contract fields

<!-- field: use_case -->
#### Use case

rel01.0-uc001-generator-coding

<!-- field: go_test -->
#### Go test

TestParseToolDefs

<!-- field: inputs | form: "tree" -->
#### Inputs

- declaration: exec tool with args, dir, emits, output cap, precondition, parameters, and side effects

<!-- field: expected | form: "tree" -->
#### Expected

- schema_fields_loaded: true
- emits_loaded: true
- side_effects_loaded: true
- visibility_loaded: true

<!-- field: traces -->
#### Traces

- srd023-cli-exec-tools AC1

<!-- case: Exec tools map argv, process outcomes, caps, and timeouts -->
### Exec tools map argv, process outcomes, caps, and timeouts

<!-- field: use_case -->
#### Use case

rel01.0-uc001-generator-coding

<!-- field: go_test -->
#### Go test

go test ./internal/tools/exec -run 'Test(ExecCmd_(BuildArgs|Execute)|ExecBuilder_|RunProcGroup)'

<!-- field: inputs | form: "tree" -->
#### Inputs

- exec_defs: argv mapping, working dir, output cap, precondition, process group timeout

<!-- field: expected | form: "tree" -->
#### Expected

- args_mapped: true
- output_and_exit_mapped: true
- output_cap_applied: true
- timeout_handled: true

<!-- field: traces -->
#### Traces

- srd023-cli-exec-tools AC3

<!-- case: Ollama adapter options configure client without hidden checks -->
### Ollama adapter options configure client without hidden checks

<!-- field: use_case -->
#### Use case

rel01.0-uc001-generator-coding

<!-- field: go_test -->
#### Go test

TestInvokeLLMBuildsWithoutNetworkIO

<!-- field: inputs | form: "tree" -->
#### Inputs

- http_client_timeout: 37s
- model_check: absent

<!-- field: expected | form: "tree" -->
#### Expected

- custom_client_used: true
- startup_network_requests: 0

<!-- field: traces -->
#### Traces

- srd009-ollama-adapter AC5

<!-- case: Ollama migration removes legacy adapter paths -->
### Ollama migration removes legacy adapter paths

<!-- field: use_case -->
#### Use case

rel01.0-uc001-generator-coding

<!-- field: go_test -->
#### Go test

TestOllamaMigrationRemovesLegacyAdapterPaths

<!-- field: inputs | form: "tree" -->
#### Inputs

- legacy_paths:
  1. internal/llm/adapter.go
  1. cmd/planner/ollama.go
- retired_packages:
  1. internal/model/llm/ollama
  1. internal/model/llm/cohere

<!-- field: expected | form: "tree" -->
#### Expected

- legacy_generator_adapter_absent: true
- legacy_planner_adapter_absent: true
- current_invoke_imports_adapter: true

<!-- field: traces -->
#### Traces

- srd009-ollama-adapter AC6

<!-- case: LLM extraction functions handle fences and native tokens -->
### LLM extraction functions handle fences and native tokens

<!-- field: use_case -->
#### Use case

rel01.0-uc001-generator-coding

<!-- field: go_test -->
#### Go test

go test ./internal/model/llm -run 'Test(StripCodeFences|ExtractWithEnvelope|StripThinkingBlocks|ExtractBraces|MakeNativeTokenExtractor)'

<!-- field: inputs | form: "tree" -->
#### Inputs

- extraction_inputs: code fences, envelopes, thinking blocks, braces, native tokens

<!-- field: expected | form: "tree" -->
#### Expected

- code_fences_removed: true
- envelope_extracted: true
- thinking_blocks_stripped: true
- braces_extracted: true
- native_token_extracted: true

<!-- field: traces -->
#### Traces

- srd017-llm-extensions AC4
