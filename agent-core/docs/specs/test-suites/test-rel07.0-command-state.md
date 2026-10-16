---
id: test-rel07.0-command-state
title: Test Suite for Release 07.0 Command-State Store and Two-Plane Persistence
release: "07.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel07.0-uc001-command-state-threading
- rel07.0-uc002-two-plane-revert

<!-- field: preconditions -->
## Preconditions

- srd038-command-state-store defines the core store contract, label addressing, authority boundary, lifetime, and redaction ordering before code changes.
- srd036-dolt-state-persistence (amended) defines the five-table schema splitting tool_results into tool_outputs and receipts.
- srd028-rest-client-tools R13 defines body_source command_state consumption through $from selectors.
- srd006-yaml-machine-config and cfg-machine-format define optional TransitionSpec label as the authored command-state address, distinct from command_name and metric_labels.
- Epics GH-269 (persistence split) and GH-270 (core store and REST consumption) established the original executable evidence.
- GH-781, GH-782, and GH-783 add authored-label capture, typed selector errors, and restart persistence evidence.
- GH-787 specifies exec parameter sources and the reserved machine_request seed; GH-786 adds their executable evidence.
- GH-895 adds command-state stdin on exec ToolDefs; GH-1884 adds the env overlay expanded from resolved command parameters.

<!-- field: test_cases -->
## Test cases

<!-- case: Explicit raw selector returns redacted labeled output | status: implemented -->
### Explicit raw selector returns redacted labeled output

A consumer explicitly selecting $from(label).$ receives the complete
versioned, redaction-applied output string while receipts remain absent.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-command-state-threading

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run TestResolveFromSelectorReturnsRawLabeledOutput

<!-- field: expected -->
#### Expected

```yaml
raw_output_preserved: true
receipt_exposed: false
```

<!-- field: traces -->
#### Traces

- srd038-command-state-store AC2
- srd038-command-state-store AC3

<!-- case: A decimal component indexes an array and stays a key on an object | status: implemented -->
### A decimal component indexes an array and stays a key on an object

Current-value, $from, and REST response selectors walk $.data.0.embedding
into an array's first member, keep "0" a key where the walk stands on an
object, and leave an index past the end, a signed index, or a key applied
to an array unresolved (srd038 R2.20, GH-2246).

<!-- field: use_case -->
#### Use case

rel07.0-uc001-command-state-threading

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core ./internal/tools/rest -run 'TestSelectorIndexesArraysAndKeepsDigitKeys|TestFromSelectorIndexesAPriorStepsArray|TestResponseSelectorsShareNestedMapSemantics'

<!-- field: expected -->
#### Expected

```yaml
array_indexed: true
digit_key_on_object_preserved: true
out_of_range_unresolved: true
signed_index_unresolved: true
```

<!-- field: traces -->
#### Traces

- srd038-command-state-store AC2
- srd038-command-state-store R2.20

<!-- case: Adjacent label selection threads the previous word's output | status: implemented -->
### Adjacent label selection threads the previous word's output

A query word directly after an embed word selects the embed word's output
through body_source command_state and a $from(<embed-label>).<path>
selector, and presents only declared params to the request builder.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-command-state-threading

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest/client -run TestRESTClient_CommandStateThreadsNonAdjacentStep

<!-- field: inputs -->
#### Inputs

```yaml
chain: [embed_query, query_collection]
selector: $from(embed_query).embedding
body_source: command_state
```

<!-- field: expected -->
#### Expected

```yaml
selected_param_present: true
undeclared_params_absent: true
```

<!-- field: traces -->
#### Traces

- rel07.0-uc001-command-state-threading S1
- srd038-command-state-store AC2
- srd028-rest-client-tools AC12

<!-- case: Non-adjacent label selection reaches past an intervening word | status: implemented -->
### Non-adjacent label selection reaches past an intervening word

In a three-word chain the third word selects the first word's vector by
label past the intervening word and renders it as a nested array through the
single-token raw passthrough. The TracerBulletNonAdjacentWithRestart variant
is component integration: it manually assembles the Execution and an
InMemoryCheckpoint and injects the CommandStateView to simulate the
post-restart state, rather than running core.Loop suspend/resume in a fresh
process; it proves the selector resolves against a rehydrated store, not the
end-to-end lifecycle path.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-command-state-threading

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest/client -run 'TestRESTClient_CommandState(ThreadsNonAdjacentStep|TracerBulletNonAdjacentWithRestart)'

<!-- field: inputs -->
#### Inputs

```yaml
chain: [embed_query, record_event, query_collection]
selector: $from(embed_query).embedding
```

<!-- field: expected -->
#### Expected

```yaml
request_body_query_embeddings: nested_array
intervening_word_ignored: true
```

<!-- field: traces -->
#### Traces

- rel07.0-uc001-command-state-threading S1
- srd038-command-state-store AC2
- srd028-rest-client-tools AC12

<!-- case: Exec arguments select a non-adjacent labelled result | status: implemented -->
### Exec arguments select a non-adjacent labelled result

An exec parameter whose declaration names source resolves from an earlier
labelled entry after another word has completed. A parameter without
source still resolves from the immediately previous Result.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-command-state-threading

<!-- field: go_test -->
#### Go test

go test ./internal/tools/exec -run TestExecBuilderCommandStateSourceAcrossInterveningStep

<!-- field: owner -->
#### Owner

GH-786

<!-- field: inputs -->
#### Inputs

```yaml
chain: [capture_request, count_before, run_corpus_ingest]
sourced_parameter: {name: directory, source: "$from(capture_request).parameters.directory"}
adjacent_parameter: {name: mode, previous_result_path: parameters.mode}
```

<!-- field: expected -->
#### Expected

```yaml
sourced_argument_present: true
intervening_word_ignored: true
adjacent_mapping_unchanged: true
```

<!-- field: traces -->
#### Traces

- rel07.0-uc001-command-state-threading S1
- srd038-command-state-store AC2
- srd023-cli-exec-tools AC2

<!-- case: Broken exec sources fail before subprocess launch | status: implemented -->
### Broken exec sources fail before subprocess launch

Missing labels and paths preserve their typed selector causes and emit
CommandError. The exec word does not launch its binary with an empty
argument. A malformed source fails declaration validation before the run.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-command-state-threading

<!-- field: go_test -->
#### Go test

go test ./internal/tools/exec -run TestExecBuilderCommandStateSourceFailures

<!-- field: owner -->
#### Owner

GH-786

<!-- field: inputs -->
#### Inputs

```yaml
missing_label: "$from(absent).parameters.directory"
missing_path: "$from(seed).parameters.absent"
malformed_source: "$from(seed"
```

<!-- field: expected -->
#### Expected

```yaml
signal: CommandError
typed_cause_preserved: true
subprocess_started: false
malformed_source_rejected_at_startup: true
```

<!-- field: traces -->
#### Traces

- rel07.0-uc001-command-state-threading E1
- srd038-command-state-store AC1
- srd023-cli-exec-tools AC3
- srd023-cli-exec-tools AC4

<!-- case: Bounded exec stdin connects separate declarative words | status: implemented -->
### Bounded exec stdin connects separate declarative words

One exec producer opts into structured output and emits two words. A
separate consumer binary selects the producer's redacted output field and
receives it over bounded stdin without a shell or authored pipeline.
Companion cases reject malformed, missing, empty, non-string, oversized,
and receipt-shaped inputs before launch; cancel and join a blocked child;
and select the same safe value after checkpoint save/load without placing
stdin in traces, metrics, or receipts.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-command-state-threading

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/tools/exec -run 'Test(ParseToolDefsExecIOContract|ParseToolDefsRejectsInvalidExecIO|ExecStdin)'

<!-- field: owner -->
#### Owner

GH-895

<!-- field: inputs -->
#### Inputs

```yaml
chain: [produce_words, count_words]
producer_binary: printf
consumer_binary: wc
selector: "$from(produce_words).output"
invalid_sources: [malformed, missing, empty, non_string, oversized, receipt]
lifecycle: [cancel, checkpoint_save_load]
```

<!-- field: expected -->
#### Expected

```yaml
consumer_word_count: 2
shell_invoked: false
invalid_input_launches: 0
cancellation_joins_process: true
safe_output_survives_resume: true
stdin_in_telemetry_or_receipt: false
```

<!-- field: traces -->
#### Traces

- rel07.0-uc001-command-state-threading S1
- rel07.0-uc001-command-state-threading S3
- rel07.0-uc001-command-state-threading S5
- srd038-command-state-store AC1
- srd038-command-state-store AC5
- srd023-cli-exec-tools AC6

<!-- case: Exec env overlays resolved parameters onto the inherited environment | status: implemented -->
### Exec env overlays resolved parameters onto the inherited environment

An exec word declares env KEY=VALUE entries whose values expand
"{{ params.NAME }}" tokens from resolved parameters, including values
sourced from a labelled command-state entry. Declared keys overlay the
inherited parent environment rather than replacing it. Companion cases
reject malformed entries, undeclared template names, builtin env, and
missing or empty expansions before launch, and they omit resolved env
values from traces, metrics, and receipts. Argv-only and stdin-only
declarations keep their prior behavior.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-command-state-threading

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/tools/exec -run 'Test(ParseToolDefsExecEnv|ParseToolDefsRejectsInvalidExecEnv|ParseToolDefsRejectsEnvWithoutBinary|ExecEnv)'

<!-- field: owner -->
#### Owner

GH-1884

<!-- field: inputs -->
#### Inputs

```yaml
env:
    - GH1884_COLLECTION={{ params.collection }}
sourced_parameter: {name: collection, source: "$from(seed).parameters.collection"}
invalid: [malformed, undeclared_token, missing, empty, builtin]
```

<!-- field: expected -->
#### Expected

```yaml
child_inherits_parent_environment: true
declared_keys_overlay: true
env_only_parameter_argv_entries: 0
invalid_input_launches: 0
env_in_telemetry_or_receipt: false
```

<!-- field: traces -->
#### Traces

- rel07.0-uc001-command-state-threading S1
- rel07.0-uc001-command-state-threading S5
- srd023-cli-exec-tools AC1
- srd023-cli-exec-tools AC7

<!-- case: Machine request seed stays safe and addressable after intervening words | status: implemented -->
### Machine request seed stays safe and addressable after intervening words

A fresh machine_request execution contains one synthetic entry labelled
seed. Core removes its configured sensitive field before the entry becomes
live or reaches checkpoint storage. An exec word after intervening reads
still resolves non-sensitive input through
$from(seed).parameters.directory, and resumed execution does not prepend a
second seed. Malformed metadata omits the whole seed output.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-command-state-threading

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest -run 'TestMachineRequest(Seed(CommandStateAddress|RedactionLiveAndReloaded|RedactionFailsClosed)|SensitiveFieldsMustBeMapped)'

<!-- field: owner -->
#### Owner

GH-489

<!-- field: inputs -->
#### Inputs

```yaml
selector: "$from(seed).parameters.directory"
sensitive_selector: "$from(seed).parameters.token"
intervening_words: 2
resume_with_seed_entry: true
```

<!-- field: expected -->
#### Expected

```yaml
directory_resolved: true
sensitive_field_resolved: false
sensitive_value_persisted: false
malformed_output_omitted: true
seed_entries_fresh_run: 1
seed_entries_resumed_run: 1
```

<!-- field: traces -->
#### Traces

- rel07.0-uc001-command-state-threading S1
- rel07.0-uc001-command-state-threading S3
- srd038-command-state-store AC1
- srd038-command-state-store AC2
- srd038-command-state-store AC5
- srd030-machine-request-rest-binding AC4

<!-- case: Duplicate authored labels resolve to the most recent step | status: implemented -->
### Duplicate authored labels resolve to the most recent step

Two synthetic entries with different command names share an authored
label; a later lookup resolves to the higher step_index.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-command-state-threading

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run TestCommandStateViewDuplicateLabelMostRecentWins

<!-- field: inputs -->
#### Inputs

```yaml
steps:
    - {label: fetch, command_name: rest_get, step_index: 0}
    - {label: fetch, command_name: rest_retry, step_index: 2}
selector: $from(fetch).id
```

<!-- field: expected -->
#### Expected

```yaml
resolved_step_index: 2
```

<!-- field: traces -->
#### Traces

- rel07.0-uc001-command-state-threading S4
- srd038-command-state-store AC2

<!-- case: MachineSpec captures a validated authored transition label | status: implemented -->
### MachineSpec captures a validated authored transition label

A real declarative loop executes a transition whose label differs from its
action. The resulting Entry records both values. Invalid label grammar
fails startup with the transition index and field name, while duplicate
valid labels remain legal.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-command-state-threading

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run TestCommandStateAuthoredLabelMachineExecution

<!-- field: owner -->
#### Owner

GH-782

<!-- field: inputs -->
#### Inputs

```yaml
transition:
    action: fetch_document
    label: source
invalid_labels: ["", "source.part", "source label", "source(label)"]
```

<!-- field: expected -->
#### Expected

```yaml
entry_label: source
entry_command_name: fetch_document
invalid_label_rejected: true
duplicate_valid_labels_accepted: true
```

<!-- field: traces -->
#### Traces

- rel07.0-uc001-command-state-threading S1
- srd038-command-state-store AC2
- srd006-yaml-machine-config AC1
- srd006-yaml-machine-config AC2

<!-- case: Authored labels and command names share one newest-match lookup | status: implemented -->
### Authored labels and command names share one newest-match lookup

Entries carry labels distinct from command names. Lookup matches either
address while scanning newest to oldest, so duplicate labels and a label
that collides with another step's command name resolve deterministically.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-command-state-threading

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run TestCommandStateViewCommandNameFallbackAndAddressCollision

<!-- field: owner -->
#### Owner

GH-783

<!-- field: inputs -->
#### Inputs

```yaml
steps:
    - {label: fetch, command_name: rest_get, step_index: 0}
    - {label: transform, command_name: fetch, step_index: 1}
    - {label: fetch, command_name: rest_get, step_index: 2}
```

<!-- field: expected -->
#### Expected

```yaml
authored_label_fetch_resolves_step_index: 2
command_name_rest_get_resolves_step_index: 2
label_command_name_collision_uses_newest_match: true
unlabeled_legacy_entry_uses_command_name: true
```

<!-- field: traces -->
#### Traces

- rel07.0-uc001-command-state-threading S4
- srd038-command-state-store AC2

<!-- case: Unresolved label and path failures retain typed causes | status: implemented -->
### Unresolved label and path failures retain typed causes

ResolveFromSelector returns exported structured errors for a missing label
and a missing path. Consumer wrapping preserves errors.As access to the
offending label and path while retaining selector context in the message.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-command-state-threading

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run TestResolveFromSelectorTypedErrors

<!-- field: owner -->
#### Owner

GH-783

<!-- field: inputs -->
#### Inputs

```yaml
unresolved_label: missing
unresolved_path: mapped.absent
```

<!-- field: expected -->
#### Expected

```yaml
unresolved_label_error_typed: true
unresolved_path_error_typed: true
consumer_wrap_preserves_cause: true
```

<!-- field: traces -->
#### Traces

- rel07.0-uc001-command-state-threading E1
- srd038-command-state-store AC1

<!-- case: Synthetic restart rehydration resolves the same command-name address | status: implemented -->
### Synthetic restart rehydration resolves the same command-name address

The store reconstructs synthetic command-name entries from a fresh
checkpoint load (InMemory and Dolt). This proves existing output
rehydration, not authored-label persistence or a real MachineSpec loop.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-command-state-threading

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core ./internal/runtime/checkpoint/dolt -run 'TestCommandStateViewRehydrates(AcrossInMemoryCheckpoint|FromDoltLoad)'

<!-- field: coverage -->
#### Coverage

component-integration

<!-- field: inputs -->
#### Inputs

```yaml
suspend_after: embed_query
resume_in: fresh_process
selector: $from(embed_query).embedding
```

<!-- field: expected -->
#### Expected

```yaml
label_resolved_after_restart: true
value_equivalent_to_live_run: true
```

<!-- field: traces -->
#### Traces

- rel07.0-uc001-command-state-threading S3
- srd038-command-state-store AC1

<!-- case: Real MachineSpec labels survive schema upgrade and restart rehydration | status: implemented -->
### Real MachineSpec labels survive schema upgrade and restart rehydration

A real MachineSpec loop records an authored label, a fresh checkpoint
instance loads the persisted execution, and a fresh CommandStateView
resolves the same label. The schema upgrade is idempotent, and an old row
with a null label remains addressable by command name.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-command-state-threading

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/checkpoint/dolt -run 'Test(DoltCheckpointSchemaUpgradeAddsNullableLabel|CommandStateAuthoredLabelPersistsAcrossDoltRestart)'

<!-- field: owner -->
#### Owner

GH-781

<!-- field: coverage -->
#### Coverage

component-integration

<!-- field: inputs -->
#### Inputs

```yaml
authored_label: embed
command_name: rest_embed
legacy_label: null
upgrade_runs: 2
```

<!-- field: expected -->
#### Expected

```yaml
authored_label_round_trips: true
fresh_view_resolves_label: true
legacy_row_uses_command_name: true
repeated_upgrade_succeeds: true
receipt_remains_unaddressable: true
```

<!-- field: traces -->
#### Traces

- rel07.0-uc001-command-state-threading S3
- rel07.0-uc001-command-state-threading S4
- srd038-command-state-store AC1
- srd038-command-state-store AC2

<!-- case: Selector rejection table blocks unsafe threading | status: implemented -->
### Selector rejection table blocks unsafe threading

The R13.3 rejection list applies to command_state threading: an undeclared
input_mapping target, a selector naming an authority field, and a malformed
$from selector are each rejected with a clear message.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-command-state-threading

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest/client -run 'TestRESTClient_CommandState(RejectionTable|MissNamesLabel|NoViewConfiguredRejected)'

<!-- field: inputs -->
#### Inputs

```yaml
cases:
    - {selector: $from(embed_query).embedding, target: undeclared_param, reject: true}
    - {selector: $from(embed_query).headers.authorization, target: auth_header, reject: true}
    - {selector: "$from(.broken", target: query_embeddings, reject: true}
```

<!-- field: expected -->
#### Expected

```yaml
all_rejected: true
transport_authority_unreachable: true
```

<!-- field: traces -->
#### Traces

- rel07.0-uc001-command-state-threading S2
- srd038-command-state-store AC3
- srd028-rest-client-tools AC12

<!-- case: REST response mapping emits typed redaction before persistence | status: implemented -->
### REST response mapping emits typed redaction before persistence

mapClientResponse applies the REST operation's configured selector before
Execute returns Result, preserving the immediate marker contract and
attaching typed paths. The test passes that Result through core.DigestResult,
constructs a live Execution, saves and reloads it through
InMemoryCheckpoint, and proves a later $from selector cannot recover the
field.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-command-state-threading

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest/client -run TestRESTClient_RedactionRunsBeforePersistence

<!-- field: coverage -->
#### Coverage

component-integration

<!-- field: inputs -->
#### Inputs

```yaml
configured_selector: body.secret
response_secret: body-secret
```

<!-- field: expected -->
#### Expected

```yaml
original_secret_absent_from_result: true
replacement_marker_present_in_result: true
field_absent_from_live_entry: true
field_absent_after_reload: true
selector_resolves_field: false
```

<!-- field: traces -->
#### Traces

- srd028-rest-client-tools AC7
- srd028-rest-client-tools AC8

<!-- case: Typed redaction denies sensitive output live and after persistence | status: implemented -->
### Typed redaction denies sensitive output live and after persistence

A REST Result translates its configured selector into versioned typed paths
covering direct and mapped output. Live CommandStateView and
InMemoryCheckpoint reload cannot recover the field. The Dolt adapter proof
inspects tool_outputs write arguments and fresh-load selection. Companion
core cases exercise malformed metadata, marked non-object output, and an
unversioned legacy row.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-command-state-threading

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core ./internal/runtime/checkpoint/dolt ./internal/tools/rest -run 'Test(DoltCheckpoint(ReappliesAndPersistsOutputRedaction|PersistsFailClosedRedactionRow)|OutputRedactionFailsClosed|CommandStateRejectsUnversionedLegacyOutput|RESTClient_RedactionRunsBeforePersistence)'

<!-- field: owner -->
#### Owner

GH-489

<!-- field: coverage -->
#### Coverage

component-integration

<!-- field: inputs -->
#### Inputs

```yaml
redacted_field: body.secret
selector: $from(fetch).body.secret
persistence_backends: [in_memory, dolt]
failure_cases: [invalid_path, non_object_output, unversioned_legacy]
```

<!-- field: expected -->
#### Expected

```yaml
field_absent_from_live_entry: true
field_absent_in_tool_outputs: true
selector_resolves_field: false
unsafe_output_omitted: true
legacy_selector_error_typed: true
receipt_unchanged: true
```

<!-- field: traces -->
#### Traces

- rel07.0-uc001-command-state-threading S5
- rel07.0-uc001-command-state-threading E4
- rel07.0-uc001-command-state-threading E5
- srd035-checkpoint-port AC7
- srd036-dolt-state-persistence AC3
- srd038-command-state-store AC5

<!-- case: Five-table round-trip persists and restores both planes | status: implemented -->
### Five-table round-trip persists and restores both planes

A three-step run writes a tool_outputs row and a receipts row per step in
one per-step transaction; Load restores conversation, each entry's output
from tool_outputs, and each entry's receipt from receipts.

<!-- field: use_case -->
#### Use case

rel07.0-uc002-two-plane-revert

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/checkpoint/dolt -run 'TestDoltCheckpoint(SaveLoadRoundTrip|SplitsToolOutputsAndReceipts)'

<!-- field: inputs -->
#### Inputs

```yaml
steps: 3
tables: [machines, transitions, execution_steps, tool_outputs, receipts]
```

<!-- field: expected -->
#### Expected

```yaml
tool_outputs_rows: 3
receipts_rows: 3
load_restores_both_planes: true
```

<!-- field: traces -->
#### Traces

- rel07.0-uc002-two-plane-revert S1
- srd036-dolt-state-persistence AC2

<!-- case: Revert reaps later two-plane rows and invalidates store entries | status: implemented -->
### Revert reaps later two-plane rows and invalidates store entries

Revert(run_id, 1) resets the run branch so the tool_outputs and receipts
rows for steps 2 and 3 no longer resolve, and the command-state store's
later entries are invalid without a separate invalidation walk.

<!-- field: use_case -->
#### Use case

rel07.0-uc002-two-plane-revert

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/checkpoint/dolt -run 'TestDoltCheckpoint(RevertReapsBothPlanes|TerminalReapsRunHistoryBeforeMerge)'

<!-- field: inputs -->
#### Inputs

```yaml
steps: 3
revert_to_step_index: 1
```

<!-- field: expected -->
#### Expected

```yaml
later_tool_outputs_rows_resolve: false
later_receipts_rows_resolve: false
store_later_entries_invalid: true
```

<!-- field: traces -->
#### Traces

- rel07.0-uc002-two-plane-revert S2
- srd038-command-state-store AC4
- srd036-dolt-state-persistence AC4

<!-- case: Dolt component restores surviving receipts after reaping later rows | status: implemented -->
### Dolt component restores surviving receipts after reaping later rows

This component test calls Revert and Load through the Dolt adapter, then
manually walks the restored entries and invokes a synthetic Reverser. It
proves that the reverse plane retains surviving receipts and reaps later
receipts; it is not evidence that the production checkpoint_rollback
adapter invoked the receipt walk. The production invocation is covered by
TestCheckpointRollbackRestoresFileViaPersistedReceipt in test-rel02.0.

<!-- field: use_case -->
#### Use case

rel07.0-uc002-two-plane-revert

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/checkpoint/dolt -run TestDoltCheckpointReceiptWalkReversesSurvivingStep

<!-- field: inputs -->
#### Inputs

```yaml
surviving_step: 1
receipt_source: receipts
```

<!-- field: expected -->
#### Expected

```yaml
surviving_step_reversed: true
receipt_forward_addressable: false
```

<!-- field: traces -->
#### Traces

- rel07.0-uc002-two-plane-revert S3
- srd036-dolt-state-persistence AC3
- srd036-dolt-state-persistence AC6
