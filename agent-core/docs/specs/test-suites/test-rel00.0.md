---
id: test-rel00.0
title: Test Suite for Release 00.0 (Core Engine Import)
release: "00.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel00.0-uc001-domain-agent-integration
- rel00.0-uc002-spec-corpus-validation

<!-- field: preconditions -->
## Preconditions

- Go toolchain installed (go 1.26+)
- internal/runtime/core, internal/observability/tracing, and internal/observability/telemetry packages buildable

<!-- field: case_groups -->
## Case groups

- srd001: core types (verified by compilation and loop tests)
- srd002: state machine (verified by loop and machine tests)
- srd003: command dispatch
- srd004: tool registry (additional coverage)
- srd005: agentic loop (additional coverage)
- srd007: tracer interface
- srd008: telemetry

<!-- field: test_cases -->
## Test cases

<!-- case: Loop runs a two-step state machine to completion -->
### Loop runs a two-step state machine to completion

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_RunsToCompletion

<!-- field: inputs -->
#### Inputs

```yaml
table: "Start->Working->Finished (two commands)"
budget: "MaxIterations=100"
```

<!-- field: expected -->
#### Expected

```yaml
status: succeeded
final_state: Finished
iterations: 2
```

<!-- field: traces -->
#### Traces

- srd005-agentic-loop AC5

<!-- case: Loop stops when iteration budget is exhausted -->
### Loop stops when iteration budget is exhausted

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_BudgetExhausted

<!-- field: inputs -->
#### Inputs

```yaml
budget: "MaxIterations=1"
```

<!-- field: expected -->
#### Expected

```yaml
status: budget_exceeded
final_state: OverBudget
```

<!-- field: traces -->
#### Traces

- srd005-agentic-loop AC1
- srd005-agentic-loop AC5

<!-- case: Loop produces cancelled status on context cancellation -->
### Loop produces cancelled status on context cancellation

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_ContextCancellation, TestLoop_ActiveDispatchCancellation

<!-- field: inputs -->
#### Inputs

```yaml
contexts: [pre_cancelled, cancelled_during_active_dispatch]
```

<!-- field: expected -->
#### Expected

```yaml
status: cancelled
active_context_command_joined: true
```

<!-- field: traces -->
#### Traces

- srd005-agentic-loop AC5

<!-- case: OnResult hook is called after each dispatch -->
### OnResult hook is called after each dispatch

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_OnResultHook

<!-- field: inputs -->
#### Inputs

```yaml
hooks: "OnResult increments a counter"
```

<!-- field: expected -->
#### Expected

```yaml
call_count: 2
```

<!-- field: traces -->
#### Traces

- srd005-agentic-loop AC3

<!-- case: Custom BudgetExceeded hook triggers budget exhaustion -->
### Custom BudgetExceeded hook triggers budget exhaustion

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_CustomBudgetHook

<!-- field: inputs -->
#### Inputs

```yaml
hooks: "BudgetExceeded returns true at iteration >= 1"
```

<!-- field: expected -->
#### Expected

```yaml
status: budget_exceeded
```

<!-- field: traces -->
#### Traces

- srd005-agentic-loop AC1
- srd005-agentic-loop AC3

<!-- case: Core source contains no mode-specific tool policy literals -->
### Core source contains no mode-specific tool policy literals

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestCoreSourceDoesNotEmbedModeOrToolPolicyNames

<!-- field: inputs -->
#### Inputs

```yaml
source: "internal/runtime/core/*.go excluding tests"
forbidden_literals: "agent mode names and workflow tool init names"
```

<!-- field: expected -->
#### Expected

```yaml
policy_boundary: "workflow behavior stays in MachineSpec and tool config"
```

<!-- field: traces -->
#### Traces

- srd005-agentic-loop AC5

<!-- case: Agent runtime does not branch on agent mode names -->
### Agent runtime does not branch on agent mode names

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestMainRuntimeDoesNotBranchOnAgentModeNames

<!-- field: inputs -->
#### Inputs

```yaml
source: "cmd/agent/main.go"
forbidden_branches: "if/switch comparisons against executor, planner, critic, bench, jurist"
```

<!-- field: expected -->
#### Expected

```yaml
bootstrap_boundary: "selected tool init gates are allowed; mode branches are not"
```

<!-- field: traces -->
#### Traces

- srd005-agentic-loop AC5

<!-- case: Token counts are accumulated in RunResult -->
### Token counts are accumulated in RunResult

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_AccumulatesTokens

<!-- field: inputs -->
#### Inputs

```yaml
command_cost: "50 tokens in, 50 tokens out"
```

<!-- field: expected -->
#### Expected

```yaml
tokens_in: 50
tokens_out: 50
```

<!-- field: traces -->
#### Traces

- srd005-agentic-loop AC2
- srd005-agentic-loop AC5

<!-- case: Loop emits init.registry_frozen OTel event -->
### Loop emits init.registry_frozen OTel event

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_EmitsRegistryFrozenEvent

<!-- field: inputs -->
#### Inputs

```yaml
trace: "loopRecorder"
```

<!-- field: expected -->
#### Expected

```yaml
event: "init.registry_frozen"
```

<!-- field: traces -->
#### Traces

- srd005-agentic-loop AC5
- srd004-tool-registry AC2

<!-- case: RunResult.Duration is positive after completion -->
### RunResult.Duration is positive after completion

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_DurationTracked

<!-- field: inputs -->
#### Inputs

```yaml
table: "two-step state machine"
```

<!-- field: expected -->
#### Expected

```yaml
duration: "> 0"
```

<!-- field: traces -->
#### Traces

- srd005-agentic-loop AC2

<!-- case: Omitted terminal status fails closed regardless of state spelling -->
### Omitted terminal status fails closed regardless of state spelling

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_OmittedTerminalStatusFailsClosed

<!-- field: inputs -->
#### Inputs

```yaml
states: "[Succeeded, BudgetExceeded, Anything]"
```

<!-- field: expected -->
#### Expected

```yaml
succeeded_status: failed
budget_status: failed
other_status: failed
```

<!-- field: traces -->
#### Traces

- srd005-agentic-loop AC5

<!-- case: Core budget check on iteration and token limits -->
### Core budget check on iteration and token limits

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestCoreBudgetExceeded

<!-- field: inputs -->
#### Inputs

```yaml
budget_scenarios: "zero limits, iteration limits, token limits"
```

<!-- field: expected -->
#### Expected

```yaml
zero_unlimited: false
at_limit: true
below_limit: false
```

<!-- field: traces -->
#### Traces

- srd005-agentic-loop AC1

<!-- case: Token budget triggers BudgetExhausted signal -->
### Token budget triggers BudgetExhausted signal

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_TokenBudgetExhausted

<!-- field: inputs -->
#### Inputs

```yaml
budget: "MaxTokens=1"
command_cost: "50 tokens"
```

<!-- field: expected -->
#### Expected

```yaml
status: budget_exceeded
```

<!-- field: traces -->
#### Traces

- srd005-agentic-loop AC1
- srd005-agentic-loop AC5

<!-- case: Registered invoke limits do not replace the machine run budget -->
### Registered invoke limits do not replace the machine run budget

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./cmd/agent ./internal/tools/llm -run 'Test(LoopParamsRunBudgetIgnoresRegisteredInvokeLimits|InvokeCallTimeoutUsesMaxTimeOnlyAsLegacyCallAlias)'

<!-- field: inputs -->
#### Inputs

```yaml
machine_budget: {max_tokens: 321, max_duration: 2h}
invoke_orders: [fast_then_deep, deep_then_fast]
invoke_limits: [max_time, max_tokens, llm_timeout]
```

<!-- field: expected -->
#### Expected

```yaml
machine_budget_authoritative: true
declaration_order_independent: true
unbounded_machine_has_no_inherited_deadline: true
max_time_is_per_call_compatibility_alias: true
llm_timeout_takes_precedence: true
```

<!-- field: traces -->
#### Traces

- srd005-agentic-loop R1.4
- srd005-agentic-loop AC1

<!-- case: Valid YAML parses into correct MachineSpec -->
### Valid YAML parses into correct MachineSpec

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestParseMachineSpec_Valid

<!-- field: inputs -->
#### Inputs

```yaml
yaml: "4 states, 2 terminals, 3 transitions"
```

<!-- field: expected -->
#### Expected

```yaml
name: test-machine
states: 4
terminals: 2
transitions: 3
```

<!-- field: traces -->
#### Traces

- srd006-yaml-machine-config AC1

<!-- case: Missing initial_state produces validation error -->
### Missing initial_state produces validation error

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestParseMachineSpecValidationErrors

<!-- field: inputs -->
#### Inputs

```yaml
yaml: "no initial_state field"
```

<!-- field: expected -->
#### Expected

```yaml
error: "initial_state is required"
```

<!-- field: traces -->
#### Traces

- srd006-yaml-machine-config AC2

<!-- case: Unknown state in transition produces validation error -->
### Unknown state in transition produces validation error

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestParseMachineSpecValidationErrors

<!-- field: inputs -->
#### Inputs

```yaml
yaml: "transition references state 'C' not in states list"
```

<!-- field: expected -->
#### Expected

```yaml
error: "state \"C\" not in states list"
```

<!-- field: traces -->
#### Traces

- srd006-yaml-machine-config AC2

<!-- case: Unknown signal in transition produces validation error -->
### Unknown signal in transition produces validation error

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestParseMachineSpecValidationErrors

<!-- field: inputs -->
#### Inputs

```yaml
yaml: "transition references signal 'Unknown' not in signals list"
```

<!-- field: expected -->
#### Expected

```yaml
error: "signal \"Unknown\" not in signals list"
```

<!-- field: traces -->
#### Traces

- srd006-yaml-machine-config AC2

<!-- case: Terminal state not in states list produces validation error -->
### Terminal state not in states list produces validation error

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestParseMachineSpecValidationErrors

<!-- field: inputs -->
#### Inputs

```yaml
yaml: "terminal_states: [Z] but Z not in states"
```

<!-- field: expected -->
#### Expected

```yaml
error: "terminal_state \"Z\" not in states list"
```

<!-- field: traces -->
#### Traces

- srd006-yaml-machine-config AC2

<!-- case: BuildTransitionTable resolves named actions from registry -->
### BuildTransitionTable resolves named actions from registry

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestBuildTransitionTable

<!-- field: inputs -->
#### Inputs

```yaml
spec: "valid YAML with do_work action"
registry: "do_work builder registered"
```

<!-- field: expected -->
#### Expected

```yaml
table_size: 3
done_terminal: true
error_terminal: true
```

<!-- field: traces -->
#### Traces

- srd006-yaml-machine-config AC3

<!-- case: $tool sentinel invokes the caller-supplied toolAction -->
### $tool sentinel invokes the caller-supplied toolAction

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestBuildTransitionTable_ToolActionSentinel

<!-- field: inputs -->
#### Inputs

```yaml
spec: "transition with action: $tool"
tool_action: "sets called=true"
```

<!-- field: expected -->
#### Expected

```yaml
called: true
```

<!-- field: traces -->
#### Traces

- srd006-yaml-machine-config AC3

<!-- case: $tool with nil toolAction produces an error -->
### $tool with nil toolAction produces an error

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestBuildTransitionTable_ToolActionNil

<!-- field: inputs -->
#### Inputs

```yaml
spec: "transition with action: $tool"
tool_action: nil
```

<!-- field: expected -->
#### Expected

```yaml
error: "not nil"
```

<!-- field: traces -->
#### Traces

- srd006-yaml-machine-config AC3

<!-- case: Unknown action name produces an error -->
### Unknown action name produces an error

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestBuildTransitionTable_UnknownAction

<!-- field: inputs -->
#### Inputs

```yaml
spec: "transition with action: nonexistent"
registry: empty
```

<!-- field: expected -->
#### Expected

```yaml
error: "\"nonexistent\" not found in registry"
```

<!-- field: traces -->
#### Traces

- srd006-yaml-machine-config AC3

<!-- case: Loading a nonexistent file produces an error -->
### Loading a nonexistent file produces an error

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoadMachineSpec_FileNotFound

<!-- field: inputs -->
#### Inputs

```yaml
path: "/nonexistent/path/machine.yaml"
```

<!-- field: expected -->
#### Expected

```yaml
error: "not nil"
```

<!-- field: traces -->
#### Traces

- srd006-yaml-machine-config AC4

<!-- case: MachineSpec iterates trusted command-state arrays -->
### MachineSpec iterates trusted command-state arrays

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run 'TestLoopForEach(DispatchesItemsInOrderAndJoins|EmptyCollectionJoinsDirectlyToTerminal|ResumeContinuesAfterLastPersistedItem|CollectAllJoinsPartialInInputOrder|AbortSignalStopsCollectAll)|TestMachineSpecForEach(Validation|RoundTrips)|TestIteratorJoinStructuredOutput'

<!-- field: inputs -->
#### Inputs

```yaml
items: $from(items).items
mode: sequential
policies: [fail_fast, collect_all]
```

<!-- field: expected -->
#### Expected

```yaml
deterministic_input_order: true
one_dispatch_and_receipt_per_item: true
join_recorded: true
resume_does_not_replay_completed_items: true
structured_projection_uses_redacted_digest_only: true
structured_projection_preserves_json_shapes: true
plain_text_and_omitted_output_have_no_projection: true
persisted_iterator_outcomes_unchanged: true
unsupported_parallel_rejected: true
```

<!-- field: traces -->
#### Traces

- srd006-yaml-machine-config AC5

<!-- case: Builtin collection words transform and validate selected arrays and structured output -->
### Builtin collection words transform and validate selected arrays and structured output

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/tools/control ./internal/tools/compose ./internal/tools/llm -run 'Test(Partition|SelectSubset|ValidatePartition|ValidateSelectSubset|RenderEach|ValidateRenderEach|FlatMap|ReorderByIndex|ValidateArrayTransformConfig|ValidateParseStructuredConfig|ParseStructured|CompileStructured)'

<!-- field: inputs -->
#### Inputs

```yaml
words: [partition, select_subset, render_each, flat_map, reorder_by_index, parse_structured]
outcomes: [full, partial, empty, nested_arrays, mismatched_lengths, ranked_subset, invalid_index, adjacent_valid, adjacent_invalid, malformed, schema_invalid, unresolved_selector]
```

<!-- field: expected -->
#### Expected

```yaml
input_order_preserved: true
undeclared_candidates_excluded: true
trusted_vocabulary_values_selected_in_candidate_order: true
duplicate_name_policy_deterministic: true
iterator_structured_output_traversable: true
parallel_nested_arrays_flattened: true
parent_fields_carried: true
candidates_joined_in_result_order: true
invalid_and_duplicate_indexes_rejected: true
empty_array_renders_empty: true
parse_structured_accepts_adjacent_output: true
other_current_result_paths_rejected: true
content_failure_uses_declared_signal: true
malformed_config_rejected_before_dispatch: true
```

<!-- field: traces -->
#### Traces

- srd022-builtin-tools AC7

<!-- case: Project collection word extracts one field from each structured row | status: implemented -->
### Project collection word extracts one field from each structured row

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/tools/control ./internal/tools/compose -run 'Test(Project|ValidateProject|ProjectedStringsFeedSelectSubset|RegisterFactories)'

<!-- field: inputs -->
#### Inputs

```yaml
word: project
selector: $from(rows).items
item_paths: [id, metadata.source.id, embeddings.0]
outcomes: [ordered_projection, duplicate_values, non_string_projection, empty, missing_field, invalid_item_path, unresolved_selector]
```

<!-- field: expected -->
#### Expected

```yaml
one_field_projected_from_each_row: true
projected_duplicates_and_json_shapes_preserved: true
projected_string_values_feed_select_subset: true
empty_array_returns_empty_items_and_zero_count: true
missing_and_invalid_paths_rejected: true
selector_failures_preserve_typed_cause: true
canonical_selector_grammar_unchanged: true
```

<!-- field: traces -->
#### Traces

- srd022-builtin-tools AC7
- srd038-command-state-store AC2

<!-- case: Normalize-vector collection word stabilizes adjacent vector shapes | status: implemented -->
### Normalize-vector collection word stabilizes adjacent vector shapes

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./internal/tools/compose -run 'Test(NormalizeVector|ValidateNormalizeVector|RegisterFactories)'

<!-- field: inputs -->
#### Inputs

```yaml
word: normalize_vector
path: mapped.embedding
accepted_shapes: [flat_numeric_vector, single_numeric_row]
rejected_shapes: [empty_vector, empty_row, multiple_rows, mixed_nesting, non_numeric_member, missing_path, non_object]
```

<!-- field: expected -->
#### Expected

```yaml
flat_vector_preserved: true
single_row_unwrapped: true
surrounding_fields_preserved: true
output_redaction_preserved: true
invalid_shapes_emit_command_error: true
canonical_selector_grammar_unchanged: true
```

<!-- field: traces -->
#### Traces

- srd022-builtin-tools AC7

<!-- case: Profile closure inspection enforces finite operation authorities before dispatch -->
### Profile closure inspection enforces finite operation authorities before dispatch

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./pkg/profileaudit ./cmd/agent -run 'TestInspect|TestStartupEnforcesProfileTimeoutClosure'

<!-- field: inputs -->
#### Inputs

```yaml
missing_command_timeout: true
valid_command_timeout: 1m
closure_edges: [declaration_include_override, dynamic_selection, rest_machine_request, child_profile, evaluator_point_machine, compatibility_wrapper]
authorities: [tooldef_timeout, rest_limits, machine_request, model, child_service, evaluator_point, blocking_await]
```

<!-- field: expected -->
#### Expected

```yaml
missing_value_rejected: true
positive_finite_value_accepted: true
environment_expanded_before_comparison: true
malformed_non_positive_equal_and_larger_rejected: true
diagnostics_are_deterministic: true
normal_startup_and_validate_config_enforce_before_dispatch: true
core_fixtures_do_not_discover_application_catalog: true
```

<!-- field: traces -->
#### Traces

- srd006-yaml-machine-config AC6

<!-- case: State and Signal are named string types with three generic constants -->
### State and Signal are named string types with three generic constants

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_RunsToCompletion

<!-- field: inputs -->
#### Inputs

```yaml
signals: "Seed, BudgetExhausted, CommandError used in loop setup"
```

<!-- field: expected -->
#### Expected

```yaml
compiles: true
```

<!-- field: traces -->
#### Traces

- srd001-core-types AC1

<!-- case: Command interface has Name, Execute, and Undo methods -->
### Command interface has Name, Execute, and Undo methods

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_RunsToCompletion

<!-- field: inputs -->
#### Inputs

```yaml
commands: "fakeCmd implements Command with Name(), Execute(), and Undo()"
```

<!-- field: expected -->
#### Expected

```yaml
dispatched: true
```

<!-- field: traces -->
#### Traces

- srd001-core-types AC2

<!-- case: Cost struct has Duration, TokensIn, TokensOut, and Dollars -->
### Cost struct has Duration, TokensIn, TokensOut, and Dollars

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_AccumulatesTokens

<!-- field: inputs -->
#### Inputs

```yaml
cost: "tokenFakeCmd returns Cost{TokensIn: 50, TokensOut: 50}"
```

<!-- field: expected -->
#### Expected

```yaml
accumulated: true
```

<!-- field: traces -->
#### Traces

- srd001-core-types AC3

<!-- case: Result struct has Output, Signal, Cost, Err, and CommandName -->
### Result struct has Output, Signal, Cost, Err, and CommandName

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_RunsToCompletion

<!-- field: inputs -->
#### Inputs

```yaml
result: "fakeCmd returns Result{Signal, CommandName}"
```

<!-- field: expected -->
#### Expected

```yaml
fields_used: true
```

<!-- field: traces -->
#### Traces

- srd001-core-types AC4

<!-- case: Builder has Build(Result) Command, CommandResolver has Resolve -->
### Builder has Build(Result) Command, CommandResolver has Resolve

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestBuildTransitionTable

<!-- field: inputs -->
#### Inputs

```yaml
registry: "registers fakeBuilder, resolves by name"
```

<!-- field: expected -->
#### Expected

```yaml
builder_resolved: true
```

<!-- field: traces -->
#### Traces

- srd001-core-types AC5

<!-- case: TransitionTable maps state-signal pairs to transitions -->
### TransitionTable maps state-signal pairs to transitions

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestBuildTransitionTable

<!-- field: inputs -->
#### Inputs

```yaml
table: "3 transitions built from YAML"
```

<!-- field: expected -->
#### Expected

```yaml
table_size: 3
```

<!-- field: traces -->
#### Traces

- srd002-state-machine AC1

<!-- case: TerminalFunc is domain-supplied predicate -->
### TerminalFunc is domain-supplied predicate

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestBuildTransitionTable

<!-- field: inputs -->
#### Inputs

```yaml
terminals: "[Done, Error]"
```

<!-- field: expected -->
#### Expected

```yaml
done_terminal: true
idle_not_terminal: true
```

<!-- field: traces -->
#### Traces

- srd002-state-machine AC2

<!-- case: Step performs lookup, calls ActionFunc, handles missing entries -->
### Step performs lookup, calls ActionFunc, handles missing entries

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_RunsToCompletion

<!-- field: inputs -->
#### Inputs

```yaml
state_machine: "Start->Working->Finished via Step calls"
```

<!-- field: expected -->
#### Expected

```yaml
transitions: 3
```

<!-- field: traces -->
#### Traces

- srd002-state-machine AC3

<!-- case: IsTerminal delegates to TerminalFunc -->
### IsTerminal delegates to TerminalFunc

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_RunsToCompletion

<!-- field: inputs -->
#### Inputs

```yaml
terminal: "Finished, OverBudget, Broken"
```

<!-- field: expected -->
#### Expected

```yaml
loop_stops: true
```

<!-- field: traces -->
#### Traces

- srd002-state-machine AC4

<!-- case: Dispatch creates child span and stamps attributes -->
### Dispatch creates child span and stamps attributes

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run 'TestDispatch_(StampsSpanAttributesAndCommandName|RecordsCommandErrorsOnSpan)|TestDispatchWithMonitor_EmitsDispatchMetrics'

<!-- field: inputs -->
#### Inputs

```yaml
command: "dispatchResultCmd"
tracer: "dispatchTestTracer"
monitor: "dispatchRuntimeRecorder"
```

<!-- field: expected -->
#### Expected

```yaml
child_span_recorded: true
command_attributes_stamped: true
error_recorded: true
dispatch_metrics_emitted: true
```

<!-- field: traces -->
#### Traces

- srd003-command-dispatch AC1

<!-- case: SafeExecute recovers panics and enforces timeout -->
### SafeExecute recovers panics and enforces timeout

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run 'TestSafeExecute_(RecoversPanicAndForcesCommandError|TimeoutAndIndefiniteWait)'

<!-- field: inputs -->
#### Inputs

```yaml
commands: "panic command, timeout command, zero-timeout wait command"
```

<!-- field: expected -->
#### Expected

```yaml
panic_signal: CommandError
timeout_signal: CommandError
zero_timeout_waits: true
```

<!-- field: traces -->
#### Traces

- srd003-command-dispatch AC2

<!-- case: FillDuration and ForceErrorSignal operate on Result pointers -->
### FillDuration and ForceErrorSignal operate on Result pointers

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestFillDurationAndForceErrorSignal

<!-- field: inputs -->
#### Inputs

```yaml
result: "Result pointer with zero duration and non-nil error"
```

<!-- field: expected -->
#### Expected

```yaml
duration_positive: true
existing_duration_preserved: true
signal: CommandError
```

<!-- field: traces -->
#### Traces

- srd003-command-dispatch AC3

<!-- case: Manifest returns sorted External tools filtered by PhaseMatch -->
### Manifest returns sorted External tools filtered by PhaseMatch

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestBuildTransitionTable

<!-- field: inputs -->
#### Inputs

```yaml
registry: "do_work registered with Phases: [Running]"
```

<!-- field: expected -->
#### Expected

```yaml
phase_match: true
```

<!-- field: traces -->
#### Traces

- srd004-tool-registry AC3

<!-- case: AllToolNames and ExternalToolNames return sorted lists -->
### AllToolNames and ExternalToolNames return sorted lists

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_EmitsRegistryFrozenEvent

<!-- field: inputs -->
#### Inputs

```yaml
registry: "multiple tools registered"
```

<!-- field: expected -->
#### Expected

```yaml
tool_count: "> 0"
```

<!-- field: traces -->
#### Traces

- srd004-tool-registry AC4

<!-- case: Registry satisfies CommandResolver interface -->
### Registry satisfies CommandResolver interface

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestBuildTransitionTable

<!-- field: inputs -->
#### Inputs

```yaml
registry: "resolves builders by name"
```

<!-- field: expected -->
#### Expected

```yaml
resolves: true
```

<!-- field: traces -->
#### Traces

- srd004-tool-registry AC5

<!-- case: Visibility enum with External and Internal, ToolSpec fields -->
### Visibility enum with External and Internal, ToolSpec fields

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_RunsToCompletion

<!-- field: inputs -->
#### Inputs

```yaml
registry: "tools with Internal visibility registered"
```

<!-- field: expected -->
#### Expected

```yaml
visibility_set: true
```

<!-- field: traces -->
#### Traces

- srd004-tool-registry AC1

<!-- case: LoopParams bundles all inputs for Loop -->
### LoopParams bundles all inputs for Loop

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_RunsToCompletion

<!-- field: inputs -->
#### Inputs

```yaml
params: "InitialState, Prompt, Registry, Table, IsTerminal, Trace, Budget, Hooks"
```

<!-- field: expected -->
#### Expected

```yaml
loop_runs: true
```

<!-- field: traces -->
#### Traces

- srd005-agentic-loop AC4

<!-- case: Tracer interface has Push, Event, SetAttributes methods -->
### Tracer interface has Push, Event, SetAttributes methods

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_EmitsRegistryFrozenEvent

<!-- field: inputs -->
#### Inputs

```yaml
tracer: "loopRecorder implements Tracer with Push, Event, SetAttributes"
```

<!-- field: expected -->
#### Expected

```yaml
interface_satisfied: true
```

<!-- field: traces -->
#### Traces

- srd007-tracer-interface AC1

<!-- case: internal/observability/tracing depends only on otel/attribute -->
### internal/observability/tracing depends only on otel/attribute

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestLoop_RunsToCompletion

<!-- field: inputs -->
#### Inputs

```yaml
imports: "internal/observability/tracing imports only go.opentelemetry.io/otel/attribute"
```

<!-- field: expected -->
#### Expected

```yaml
no_sdk_import: true
```

<!-- field: traces -->
#### Traces

- srd007-tracer-interface AC2

<!-- case: NewRoot rejects empty ExporterConfig -->
### NewRoot rejects empty ExporterConfig

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestNewRoot_EmptyConfigError

<!-- field: inputs -->
#### Inputs

```yaml
config: "ExporterConfig{} (both fields empty)"
```

<!-- field: expected -->
#### Expected

```yaml
error: "at least one exporter required"
```

<!-- field: traces -->
#### Traces

- srd008-telemetry AC1

<!-- case: Trace Push creates child, Event records event, SetAttributes sets attributes -->
### Trace Push creates child, Event records event, SetAttributes sets attributes

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestNewRoot_PushAndEvent

<!-- field: inputs -->
#### Inputs

```yaml
trace: "Trace from NewRoot with file exporter"
```

<!-- field: expected -->
#### Expected

```yaml
child_created: true
event_recorded: true
attributes_set: true
```

<!-- field: traces -->
#### Traces

- srd008-telemetry AC2

<!-- case: Trace Context returns context, Meter returns meter -->
### Trace Context returns context, Meter returns meter

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestNewRoot_PushAndEvent

<!-- field: inputs -->
#### Inputs

```yaml
trace: "Trace from NewRoot"
```

<!-- field: expected -->
#### Expected

```yaml
context_not_nil: true
meter_not_nil: true
```

<!-- field: traces -->
#### Traces

- srd008-telemetry AC3

<!-- case: NewRoot with file exporter creates providers and root span -->
### NewRoot with file exporter creates providers and root span

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestNewRoot_FileExporter

<!-- field: inputs -->
#### Inputs

```yaml
config: "ExporterConfig{FilePath: tempfile}"
```

<!-- field: expected -->
#### Expected

```yaml
trace_not_zero: true
shutdown_not_nil: true
```

<!-- field: traces -->
#### Traces

- srd008-telemetry AC4

<!-- case: Shutdown is idempotent and flushes providers -->
### Shutdown is idempotent and flushes providers

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestNewRoot_ShutdownIdempotent

<!-- field: inputs -->
#### Inputs

```yaml
shutdown: "called twice"
```

<!-- field: expected -->
#### Expected

```yaml
no_panic: true
```

<!-- field: traces -->
#### Traces

- srd008-telemetry AC5

<!-- case: File exporter writes temp file and renames atomically -->
### File exporter writes temp file and renames atomically

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestNewRoot_TempFileUsesServiceName, TestNewRoot_FileExporter

<!-- field: inputs -->
#### Inputs

```yaml
config: "ExporterConfig{FilePath: target}"
```

<!-- field: expected -->
#### Expected

```yaml
temp_file_during: true
final_file_after: true
```

<!-- field: traces -->
#### Traces

- srd008-telemetry AC6

<!-- case: OTLP exporter configuration accepted | status: planned -->
### OTLP exporter configuration accepted

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestNewRoot_OTLPExporter

<!-- field: follow_up_issue -->
#### Follow up issue

```yaml
"#601"
```

<!-- field: inputs -->
#### Inputs

```yaml
config: "ExporterConfig{OTLPEndpoint: localhost:4317}"
```

<!-- field: expected -->
#### Expected

```yaml
trace_not_zero: true
```

<!-- field: traces -->
#### Traces

- srd008-telemetry AC7

<!-- case: TraceAdapter satisfies tracing.Tracer and wraps Push | status: planned -->
### TraceAdapter satisfies tracing.Tracer and wraps Push

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestTraceAdapter_SatisfiesInterface

<!-- field: follow_up_issue -->
#### Follow up issue

```yaml
"#601"
```

<!-- field: inputs -->
#### Inputs

```yaml
adapter: "TraceAdapter{T: trace}"
```

<!-- field: expected -->
#### Expected

```yaml
interface_check: "compile-time var _ tracing.Tracer = TraceAdapter{}"
push_returns_adapter: true
```

<!-- field: traces -->
#### Traces

- srd008-telemetry AC8

<!-- case: Declared replay loads OTLP JSON and relays an unchanged batch -->
### Declared replay loads OTLP JSON and relays an unchanged batch

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestOTLPReplayMachineLoadsThenRelaysUnchangedBatch

<!-- field: inputs -->
#### Inputs

```yaml
machine: testdata/integration/profiles/otlp-replay/machine.yaml
```

<!-- field: expected -->
#### Expected

```yaml
status: succeeded
batch_identity: unchanged
```

<!-- field: traces -->
#### Traces

- srd008-telemetry AC9

<!-- case: Checkpoint contracts compile and round-trip -->
### Checkpoint contracts compile and round-trip

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestInMemoryCheckpointRoundTripsConversationAndReceipts

<!-- field: inputs -->
#### Inputs

```yaml
checkpoint: "conversation and opaque persisted receipts (Execution entries), position snapshot"
```

<!-- field: expected -->
#### Expected

```yaml
state_layers_separate: true
json_round_trip: true
```

<!-- field: traces -->
#### Traces

- srd001-core-types AC8
- srd005-agentic-loop AC7

<!-- case: Conversation append mode avoids Chat calls -->
### Conversation append mode avoids Chat calls

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestConversationHistoryLifecycle

<!-- field: inputs -->
#### Inputs

```yaml
message: "manual user message"
```

<!-- field: expected -->
#### Expected

```yaml
chat_called: false
history_contains_message: true
```

<!-- field: traces -->
#### Traces

- srd010-conversation AC1

<!-- case: Conversation snapshot and truncate preserve order -->
### Conversation snapshot and truncate preserve order

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestConversationHistoryLifecycle

<!-- field: inputs -->
#### Inputs

```yaml
operations: "append, snapshot, truncate"
```

<!-- field: expected -->
#### Expected

```yaml
history_order_preserved: true
snapshot_restored: true
```

<!-- field: traces -->
#### Traces

- srd010-conversation AC2

<!-- case: Conversation history uses defensive copies -->
### Conversation history uses defensive copies

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestConversationHistoryReturnsCopies

<!-- field: inputs -->
#### Inputs

```yaml
mutation: "modify returned history"
```

<!-- field: expected -->
#### Expected

```yaml
internal_history_unchanged: true
```

<!-- field: traces -->
#### Traces

- srd010-conversation AC3

<!-- case: Default assembler combines prompt and history -->
### Default assembler combines prompt and history

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestDefaultAssembler_AppendsConversationMessages

<!-- field: inputs -->
#### Inputs

```yaml
system_prompt: "non-empty"
history: "user and assistant messages"
```

<!-- field: expected -->
#### Expected

```yaml
system_first: true
history_included: true
```

<!-- field: traces -->
#### Traces

- srd010-conversation AC4

<!-- case: Conversation compatibility methods support generator history migration -->
### Conversation compatibility methods support generator history migration

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestConversationHistoryLifecycle

<!-- field: inputs -->
#### Inputs

```yaml
methods: "Append, Messages, Reset"
```

<!-- field: expected -->
#### Expected

```yaml
messages_alias_history: true
reset_clears_history: true
```

<!-- field: traces -->
#### Traces

- srd010-conversation AC5

<!-- case: Conversation reset preserves append reuse -->
### Conversation reset preserves append reuse

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestConversationHistoryLifecycle

<!-- field: inputs -->
#### Inputs

```yaml
operations: "append, reset, append"
```

<!-- field: expected -->
#### Expected

```yaml
history_cleared: true
second_append_succeeds: true
```

<!-- field: traces -->
#### Traces

- srd010-conversation AC6

<!-- case: Prompt struct feeds deterministic system assembly -->
### Prompt struct feeds deterministic system assembly

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./internal/model/prompt ./internal/model/llm -run 'Test(Render|DefaultAssembler)'

<!-- field: inputs -->
#### Inputs

```yaml
prompt: "role, task, constraints, output_format"
```

<!-- field: expected -->
#### Expected

```yaml
sections_counted: true
assemble_deterministic: true
```

<!-- field: traces -->
#### Traces

- srd011-prompt-schema AC1

<!-- case: ToolDef prompt policy and runtime task remain distinct -->
### ToolDef prompt policy and runtime task remain distinct

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm ./internal/model/prompt

<!-- field: inputs -->
#### Inputs

```yaml
yaml_modes: "structured task and simple prompt"
```

<!-- field: expected -->
#### Expected

```yaml
defaults_applied: true
simple_maps_to_task: true
```

<!-- field: traces -->
#### Traces

- srd011-prompt-schema AC2

<!-- case: Prompt config rejects unknown authority -->
### Prompt config rejects unknown authority

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm ./internal/tools/catalog

<!-- field: inputs -->
#### Inputs

```yaml
yaml: "task and prompt both set"
```

<!-- field: expected -->
#### Expected

```yaml
error_reported: true
```

<!-- field: traces -->
#### Traces

- srd011-prompt-schema AC3

<!-- case: Invoke boundary constructs prompt without file loading -->
### Invoke boundary constructs prompt without file loading

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm ./internal/model/prompt

<!-- field: inputs -->
#### Inputs

```yaml
prompt_string: "inline task"
```

<!-- field: expected -->
#### Expected

```yaml
mode: string
task_length_recorded: true
```

<!-- field: traces -->
#### Traces

- srd011-prompt-schema AC4

<!-- case: Prompt defaults remain overridable by declared config -->
### Prompt defaults remain overridable by declared config

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm ./internal/model/prompt

<!-- field: inputs -->
#### Inputs

```yaml
yaml: "task plus role override"
```

<!-- field: expected -->
#### Expected

```yaml
compiled_defaults_used: true
override_applied: true
```

<!-- field: traces -->
#### Traces

- srd011-prompt-schema AC5

<!-- case: Invoke boundary rejects missing required prompt input -->
### Invoke boundary rejects missing required prompt input

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm

<!-- field: inputs -->
#### Inputs

```yaml
task: ""
```

<!-- field: expected -->
#### Expected

```yaml
error_reported: true
```

<!-- field: traces -->
#### Traces

- srd011-prompt-schema AC6

<!-- case: System prompt emits embedded template sections conditionally -->
### System prompt emits embedded template sections conditionally

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./internal/model/prompt

<!-- field: inputs -->
#### Inputs

```yaml
prompt_data: "role, task, tool manifest"
```

<!-- field: expected -->
#### Expected

```yaml
role_and_task_rendered: true
tool_manifest_rendered: true
```

<!-- field: traces -->
#### Traces

- srd012-system-prompt-template AC1

<!-- case: PromptData envelope and strict format control output -->
### PromptData envelope and strict format control output

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestRender_WithEnvelope

<!-- field: inputs -->
#### Inputs

```yaml
envelope: "<tool> tags"
strict_format: true
```

<!-- field: expected -->
#### Expected

```yaml
envelope_instructions_rendered: true
strict_format_instructions_rendered: true
```

<!-- field: traces -->
#### Traces

- srd012-system-prompt-template AC2

<!-- case: System prompt output has no trailing newline -->
### System prompt output has no trailing newline

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestRender_NoTrailingNewlines

<!-- field: inputs -->
#### Inputs

```yaml
prompt_data: "role, task, tool manifest"
```

<!-- field: expected -->
#### Expected

```yaml
trailing_newline: false
```

<!-- field: traces -->
#### Traces

- srd012-system-prompt-template AC3

<!-- case: System prompt uses the embedded reviewed template -->
### System prompt uses the embedded reviewed template

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./internal/model/prompt -run TestRender

<!-- field: inputs -->
#### Inputs

```yaml
template: embedded
```

<!-- field: expected -->
#### Expected

```yaml
rendered: true
runtime_template_compilation_absent: true
```

<!-- field: traces -->
#### Traces

- srd012-system-prompt-template AC4

<!-- case: SerializeManifest emits headings and schema fences -->
### SerializeManifest emits headings and schema fences

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestSerializeManifest_SingleToolWithSchema

<!-- field: inputs -->
#### Inputs

```yaml
tools: "one ToolSpec with JSON schema"
```

<!-- field: expected -->
#### Expected

```yaml
heading_present: true
json_fence_present: true
```

<!-- field: traces -->
#### Traces

- srd012-system-prompt-template AC5

<!-- case: SerializeManifest normalizes schema JSON -->
### SerializeManifest normalizes schema JSON

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestNormalizeSchema

<!-- field: inputs -->
#### Inputs

```yaml
schema: "compact JSON object"
```

<!-- field: expected -->
#### Expected

```yaml
schema_indented_two_spaces: true
invalid_json_preserved: true
```

<!-- field: traces -->
#### Traces

- srd012-system-prompt-template AC6

<!-- case: Spec parsers load typed YAML artifacts -->
### Spec parsers load typed YAML artifacts

<!-- field: use_case -->
#### Use case

rel00.0-uc002-spec-corpus-validation

<!-- field: go_test -->
#### Go test

TestParseSpecIndex

<!-- field: inputs -->
#### Inputs

```yaml
artifacts: "SRD, use case, test suite, roadmap, spec index"
```

<!-- field: expected -->
#### Expected

```yaml
typed_structs_loaded: true
```

<!-- field: traces -->
#### Traces

- srd014-spec-graph AC1

<!-- case: Spec graph builds node and edge relationships -->
### Spec graph builds node and edge relationships

<!-- field: use_case -->
#### Use case

rel00.0-uc002-spec-corpus-validation

<!-- field: go_test -->
#### Go test

TestBuildGraph_NodeCounts

<!-- field: inputs -->
#### Inputs

```yaml
corpus: "valid fixture corpus"
```

<!-- field: expected -->
#### Expected

```yaml
node_kinds_present: true
edges_present: true
```

<!-- field: traces -->
#### Traces

- srd014-spec-graph AC2

<!-- case: Spec validation reports graph findings -->
### Spec validation reports graph findings

<!-- field: use_case -->
#### Use case

rel00.0-uc002-spec-corpus-validation

<!-- field: go_test -->
#### Go test

TestValidate_ReturnsFindings

<!-- field: inputs -->
#### Inputs

```yaml
corpus: "fixture with validation findings"
```

<!-- field: expected -->
#### Expected

```yaml
findings_have_check_and_message: true
```

<!-- field: traces -->
#### Traces

- srd014-spec-graph AC3

<!-- case: Spec corpus loader discovers standard docs layout -->
### Spec corpus loader discovers standard docs layout

<!-- field: use_case -->
#### Use case

rel00.0-uc002-spec-corpus-validation

<!-- field: go_test -->
#### Go test

TestLoadCorpus_Valid

<!-- field: inputs -->
#### Inputs

```yaml
root: "pkg/spec/testdata/valid"
```

<!-- field: expected -->
#### Expected

```yaml
srds_loaded: true
use_cases_loaded: true
test_suites_loaded: true
```

<!-- field: traces -->
#### Traces

- srd014-spec-graph AC4

<!-- case: Built-in spec_corpus charter preserves mage audit behavior | status: implemented -->
### Built-in spec_corpus charter preserves mage audit behavior

<!-- field: use_case -->
#### Use case

rel00.0-uc002-spec-corpus-validation

<!-- field: go_test -->
#### Go test

go test ./pkg/spec -run TestExecuteChartersSpecCorpusAddsProvenance

<!-- field: integration_target -->
#### Integration target

mage audit

<!-- field: inputs -->
#### Inputs

```yaml
profile: testdata/integration/profiles/audit/profile.yaml
suite: tools/builtin/spec-validation/spec-corpus-charter.yaml
```

<!-- field: expected -->
#### Expected

```yaml
default_suite_loaded_from_yaml: true
spec_corpus_findings_have_provenance: true
existing_audit_warnings_preserved: true
terminal_state: succeeded
```

<!-- field: traces -->
#### Traces

- srd014-spec-graph AC6
- gh-123 AC2

<!-- case: GenAI attribute constants and output types are exported -->
### GenAI attribute constants and output types are exported

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestInferenceAttrs

<!-- field: inputs -->
#### Inputs

```yaml
provider: ollama
model: qwen
```

<!-- field: expected -->
#### Expected

```yaml
genai_attributes_present: true
provider_and_output_constants_present: true
```

<!-- field: traces -->
#### Traces

- srd015-genai-semconv AC1
- srd015-genai-semconv AC4

<!-- case: GenAI span names follow operation target format -->
### GenAI span names follow operation target format

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestAgentSpanName

<!-- field: inputs -->
#### Inputs

```yaml
names: "agent, tool, workflow, inference"
```

<!-- field: expected -->
#### Expected

```yaml
operation_target_format: true
empty_name_falls_back: true
```

<!-- field: traces -->
#### Traces

- srd015-genai-semconv AC2

<!-- case: GenAI attribute builders return expected slices -->
### GenAI attribute builders return expected slices

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestToolAttrs

<!-- field: inputs -->
#### Inputs

```yaml
attrs: "inference, agent, tool, workflow, usage, error"
```

<!-- field: expected -->
#### Expected

```yaml
operation_attrs_present: true
optional_empty_values_omitted: true
```

<!-- field: traces -->
#### Traces

- srd015-genai-semconv AC3

<!-- case: Delta capture isolates current-turn messages and prompt hashes -->
### Delta capture isolates current-turn messages and prompt hashes

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestInvokeLLMDeltaCaptureIsolatesTwoTurns

<!-- field: inputs -->
#### Inputs

```yaml
turns: 2
capture_level: delta
```

<!-- field: expected -->
#### Expected

```yaml
current_user_and_assistant_only: true
full_message_attributes_absent: true
rendered_system_prompt_hash_stable: true
```

<!-- field: traces -->
#### Traces

- srd015-genai-semconv AC5
- srd015-genai-semconv AC6

<!-- case: Full capture emits input and assistant message arrays -->
### Full capture emits input and assistant message arrays

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestInvokeLLMDispatchTracerIsolatesTwoTurns

<!-- field: inputs -->
#### Inputs

```yaml
turns: 2
capture_level: full
```

<!-- field: expected -->
#### Expected

```yaml
input_message_array_valid: true
output_message_array_valid: true
```

<!-- field: traces -->
#### Traces

- srd015-genai-semconv AC5

<!-- case: Off capture omits model and parser content -->
### Off capture omits model and parser content

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestCaptureOffRecordsNoContent

<!-- field: inputs -->
#### Inputs

```yaml
capture_level: off
```

<!-- field: expected -->
#### Expected

```yaml
model_content_absent: true
parser_raw_output_absent: true
```

<!-- field: traces -->
#### Traces

- srd015-genai-semconv AC5

<!-- case: Conversation references are emitted only when available -->
### Conversation references are emitted only when available

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestInvokeLLMConversationReferenceIsOptional

<!-- field: inputs -->
#### Inputs

```yaml
providers: "missing, unavailable, empty, available"
```

<!-- field: expected -->
#### Expected

```yaml
missing_reference_omitted: true
available_reference_emitted: true
```

<!-- field: traces -->
#### Traces

- srd015-genai-semconv AC6

<!-- case: Traceparent formats valid and invalid contexts -->
### Traceparent formats valid and invalid contexts

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestParseTraceparentProtocolBoundaries

<!-- field: inputs -->
#### Inputs

```yaml
span_context: valid
```

<!-- field: expected -->
#### Expected

```yaml
traceparent_round_trips: true
invalid_context_returns_empty: true
```

<!-- field: traces -->
#### Traces

- srd016-traceparent AC2

<!-- case: Traceparent parent span parsing builds remote context -->
### Traceparent parent span parsing builds remote context

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestParseParentSpan_Valid

<!-- field: inputs -->
#### Inputs

```yaml
traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
```

<!-- field: expected -->
#### Expected

```yaml
remote_context_created: true
empty_input_background: true
invalid_input_errors: true
```

<!-- field: traces -->
#### Traces

- srd016-traceparent AC3

<!-- case: Runtime request bootstrap, resume precedence, and process result contract -->
### Runtime request bootstrap, resume precedence, and process result contract

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./cmd/agent -run 'Test(SeedRequest|ResumeLoad|CLIResultReporter|RunPreparedPrints|ExitCodeForStatus)'

<!-- field: inputs -->
#### Inputs

```yaml
request_cases: [finite_bytes, absent, unreadable]
resume_case: checkpoint_load_with_request_seed_present
summary_limit_bytes: 1048576
run_statuses: [succeeded, suspended, failed, budget_exceeded, cancelled]
```

<!-- field: expected -->
#### Expected

```yaml
non_resume_request_becomes_seed_output: true
no_request_bootstrap_unchanged: true
unreadable_request_is_error: true
loaded_checkpoint_overrides_request_seed: true
stdout_summary_bounded: true
stderr_has_run_and_machine_state: true
exit_code_uses_run_status: true
```

<!-- field: traces -->
#### Traces

- cfg-runtime-contract
- srd018-cli-flag-contract AC1

<!-- case: Self invoke records bounded child process result attributes -->
### Self invoke records bounded child process result attributes

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./internal/tools/control -run TestSelfInvokeTraceRecordsBoundedChildResult

<!-- field: inputs -->
#### Inputs

```yaml
child_result: [profile, run_id, exit_code, stdout, stderr]
content_prefix_limit_bytes: 500
```

<!-- field: expected -->
#### Expected

```yaml
current_span_attributes_recorded: true
result_event_recorded: true
stdout_and_stderr_bounded: true
```

<!-- field: traces -->
#### Traces

- cfg-runtime-contract

<!-- case: Integration telemetry dual export and retained backend evidence | status: implemented -->
### Integration telemetry dual export and retained backend evidence

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: integration_target -->
#### Integration target

cd ../applications/chatbot-mesh && mage integration:rig

<!-- field: inputs -->
#### Inputs

```yaml
exporters: [file, otlp]
identity: [repository, module, target, git_commit, run_id, service_name]
backends: [collector-agent-traces, collector-agent-metrics]
```

<!-- field: expected -->
#### Expected

```yaml
child_endpoint_propagated: true
child_service_identity_explicit: true
file_evidence_preserved: true
backend_evidence_retained_after_run: true
ingress_self_loop_absent: true
```

<!-- field: traces -->
#### Traces

- rel00.0-uc001-domain-agent-integration S13
- srd008-telemetry AC10
- srd018-cli-flag-contract AC2
- srd018-cli-flag-contract AC7

<!-- case: LLM tools parse responses and manage retry state -->
### LLM tools parse responses and manage retry state

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm -run 'Test(InvokeLLM|ParseResponse|ReportParseError|ResetHistory)'

<!-- field: inputs -->
#### Inputs

```yaml
llm_tools: "invoke_llm, parse_response, report_parse_error, reset_history"
```

<!-- field: expected -->
#### Expected

```yaml
conversation_updated: true
parser_validates_tools: true
retry_counter_managed: true
timeout_and_context_limits_enforced: true
```

<!-- field: traces -->
#### Traces

- srd013-standard-tool-library AC3

<!-- case: Subprocess helpers map output and process groups -->
### Subprocess helpers map output and process groups

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./internal/tools/exec -run 'Test(SubprocessResult|CapOutput|ExtractStringParam|ProcGroupCmd|RunProcGroup)'

<!-- field: inputs -->
#### Inputs

```yaml
subprocesses: "success, non-zero exit, timeout, missing binary"
```

<!-- field: expected -->
#### Expected

```yaml
result_signals_mapped: true
output_capped: true
params_extracted: true
process_group_timeout_handled: true
```

<!-- field: traces -->
#### Traces

- srd013-standard-tool-library AC4

<!-- case: Tool declarations load and register selected tools -->
### Tool declarations load and register selected tools

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/tools/exec -run 'Test(LoadToolDeclarations|SelectTools|RegisterToolDefs|ParseToolDefs)'

<!-- field: inputs -->
#### Inputs

```yaml
declarations: "builtin and exec tool YAML"
selection: "selected tool names"
```

<!-- field: expected -->
#### Expected

```yaml
declarations_loaded: true
selection_validated: true
registry_populated: true
legacy_paths_available: true
```

<!-- field: traces -->
#### Traces

- srd013-standard-tool-library AC5

<!-- case: Exec metrics report declared live process values -->
### Exec metrics report declared live process values

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

go test ./internal/tools/exec -run '^TestExecCmdRecordsMonitorMetrics$'

<!-- field: inputs -->
#### Inputs

```yaml
command: "echo hello"
```

<!-- field: expected -->
#### Expected

```yaml
process_output_bytes_recorded: true
process_exit_code_recorded: true
command_output_not_used_as_metric_attribute: true
```

<!-- field: traces -->
#### Traces

- srd013-standard-tool-library AC7

<!-- case: SpanOverride customizes dispatch span -->
### SpanOverride customizes dispatch span

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestSpanOverrideCustomizesDispatchSpan

<!-- field: inputs -->
#### Inputs

```yaml
command: spanOverrideCmd
span_name: custom-span
span_attr: custom.attr
```

<!-- field: expected -->
#### Expected

```yaml
command_signal: ToolDone
command_name_stamped: true
custom_span_name_used: true
custom_attrs_used: true
```

<!-- field: traces -->
#### Traces

- srd001-core-types AC6

<!-- case: ToolMetrics JSON tags preserve structured counts -->
### ToolMetrics JSON tags preserve structured counts

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestToolMetricsJSONTagsAndOmitEmptyDetails

<!-- field: inputs -->
#### Inputs

```yaml
metrics: "total, passed, failed, details"
```

<!-- field: expected -->
#### Expected

```yaml
json_tags_match_contract: true
details_omitted_when_empty: true
```

<!-- field: traces -->
#### Traces

- srd001-core-types AC7

<!-- case: NoopTracer satisfies and discards tracer operations -->
### NoopTracer satisfies and discards tracer operations

<!-- field: use_case -->
#### Use case

rel00.0-uc001-domain-agent-integration

<!-- field: go_test -->
#### Go test

TestNoopTracerSatisfiesInterfaceAndDiscardsOperations

<!-- field: inputs -->
#### Inputs

```yaml
operations: "Push, Event, SetAttributes, RecordError, Context"
```

<!-- field: expected -->
#### Expected

```yaml
compile_time_interface: true
child_tracer_noop: true
context_background: true
operations_do_not_panic: true
```

<!-- field: traces -->
#### Traces

- srd007-tracer-interface AC3
