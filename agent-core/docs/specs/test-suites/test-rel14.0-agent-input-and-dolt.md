---
id: test-rel14.0-agent-input-and-dolt
title: Test Suite for Release 14.0 Request Signals and Dolt Words
release: "14.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- [rel14.0-uc001-request-driven-signal](../use-cases/rel14.0-uc001-request-driven-signal.md)
- [rel14.0-uc002-dolt-boundary-word](../use-cases/rel14.0-uc002-dolt-boundary-word.md)

<!-- field: preconditions -->
## Preconditions

- srd046-request-driven-signal-source defines model-free profiles, SignalEnvelope admission, request-driven lifecycle behavior, and model-path compatibility.
- srd047-dolt-boundary-word defines trusted Dolt provision, query, and write boundary contracts.
- GH-1733 implements the request-driven HTTP source and mock-model regression.
- GH-1732 implements unit tests plus a real-Dolt mage integration:doltWord target or equivalent standard integration target.

<!-- field: test_cases -->
## Test cases

<!-- case: Model-free profile starts without model construction | status: implemented -->
### Model-free profile starts without model construction

A profile selects an HTTP source and deterministic words but no invoke_llm
declaration. Validation and startup succeed, no model client factory runs,
and one request completes. A mock-model profile then runs through the
unchanged model path.

<!-- field: use_case -->
#### Use case

[rel14.0-uc001-request-driven-signal](../use-cases/rel14.0-uc001-request-driven-signal.md)

<!-- field: go_test -->
#### Go test

go test ./cmd/agent -run 'TestRequestSignalSource_(ModelFreeProfile|ModelRunOrResumeRegression)'

<!-- field: inputs | form: "tree" -->
#### Inputs

- profiles:
  1. model_free_http
  1. mock_model

<!-- field: expected | form: "tree" -->
#### Expected

- model_free_client_constructions: 0
- model_free_succeeds: true
- mock_model_succeeds: true

<!-- field: traces -->
#### Traces

- [rel14.0-uc001-request-driven-signal](../use-cases/rel14.0-uc001-request-driven-signal.md) [S1](../use-cases/rel14.0-uc001-request-driven-signal.md#s1)
- [rel14.0-uc001-request-driven-signal](../use-cases/rel14.0-uc001-request-driven-signal.md) [S5](../use-cases/rel14.0-uc001-request-driven-signal.md#s5)
- [srd046-request-driven-signal-source](../software-requirements/srd046-request-driven-signal-source.md) [AC1](../software-requirements/srd046-request-driven-signal-source.md#ac1)
- [srd046-request-driven-signal-source](../software-requirements/srd046-request-driven-signal-source.md) [AC6](../software-requirements/srd046-request-driven-signal-source.md#ac6)

<!-- case: Trusted mapping validates signal payload and redaction | status: implemented -->
### Trusted mapping validates signal payload and redaction

Config maps a closed discriminator to a declared signal and payload
fields. Valid input produces one SignalEnvelope. Unknown discriminators,
schema-invalid values, and raw signal names are refused before run
creation. Sensitive mapped values are absent from live and persisted
output.

<!-- field: use_case -->
#### Use case

[rel14.0-uc001-request-driven-signal](../use-cases/rel14.0-uc001-request-driven-signal.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest -run TestRequestSignalSource_TrustedMapping

<!-- field: inputs | form: "tree" -->
#### Inputs

- mappings:
  - create: WorkRequested
- invalid:
  1. unknown_discriminator
  1. raw_signal
  1. wrong_payload_type

<!-- field: expected | form: "tree" -->
#### Expected

- valid_signal: WorkRequested
- invalid_runs_created: 0
- sensitive_values_exposed: false

<!-- field: traces -->
#### Traces

- [rel14.0-uc001-request-driven-signal](../use-cases/rel14.0-uc001-request-driven-signal.md) [S1](../use-cases/rel14.0-uc001-request-driven-signal.md#s1)
- [srd046-request-driven-signal-source](../software-requirements/srd046-request-driven-signal-source.md) [AC2](../software-requirements/srd046-request-driven-signal-source.md#ac2)

<!-- case: Admission distinguishes accepted undeclared and conflict | status: implemented -->
### Admission distinguishes accepted undeclared and conflict

The current state accepts one configured signal, a globally unknown
signal is refused_undeclared, and a declared signal with no transition,
stale state guard, or active owner is refused_conflict. Refused cases
leave Position, Execution, and checkpoint bytes unchanged.

<!-- field: use_case -->
#### Use case

[rel14.0-uc001-request-driven-signal](../use-cases/rel14.0-uc001-request-driven-signal.md)

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run TestSignalAdmission_Outcomes

<!-- field: inputs | form: "tree" -->
#### Inputs

- cases:
  1. accepted
  1. unknown
  1. no_transition
  1. stale_state
  1. active_owner

<!-- field: expected | form: "tree" -->
#### Expected

- outcomes:
  1. accepted
  1. refused_undeclared
  1. refused_conflict
  1. refused_conflict
  1. refused_conflict
- refused_mutations: 0

<!-- field: traces -->
#### Traces

- [rel14.0-uc001-request-driven-signal](../use-cases/rel14.0-uc001-request-driven-signal.md) [S2](../use-cases/rel14.0-uc001-request-driven-signal.md#s2)
- [rel14.0-uc001-request-driven-signal](../use-cases/rel14.0-uc001-request-driven-signal.md) [S3](../use-cases/rel14.0-uc001-request-driven-signal.md#s3)
- [srd046-request-driven-signal-source](../software-requirements/srd046-request-driven-signal-source.md) [AC3](../software-requirements/srd046-request-driven-signal-source.md#ac3)

<!-- case: Accepted signal uses the ordinary machine execution path | status: implemented -->
### Accepted signal uses the ordinary machine execution path

An accepted envelope becomes the Result signal and payload consumed by
StateMachine.Step. The selected Registry builder dispatches normally. The
response reports admission separately from a successful and a failing
terminal outcome, including the configured HTTP 409 conflict mapping.

<!-- field: use_case -->
#### Use case

[rel14.0-uc001-request-driven-signal](../use-cases/rel14.0-uc001-request-driven-signal.md)

<!-- field: go_test -->
#### Go test

go test ./cmd/agent -run TestRequestSignalSource_OrdinaryLoopPath

<!-- field: inputs | form: "tree" -->
#### Inputs

- terminals:
  1. Succeeded
  1. Failed

<!-- field: expected | form: "tree" -->
#### Expected

- admission: accepted
- ordinary_run_events: true
- run_statuses:
  1. succeeded
  1. failed
- conflict_http_status: 409

<!-- field: traces -->
#### Traces

- [rel14.0-uc001-request-driven-signal](../use-cases/rel14.0-uc001-request-driven-signal.md) [S2](../use-cases/rel14.0-uc001-request-driven-signal.md#s2)
- [srd046-request-driven-signal-source](../software-requirements/srd046-request-driven-signal-source.md) [AC4](../software-requirements/srd046-request-driven-signal-source.md#ac4)

<!-- case: Suspended request-driven run resumes exactly once | status: implemented -->
### Suspended request-driven run resumes exactly once

The first accepted request reaches StatusSuspended and persists Position
and Execution. A second request loads that checkpoint, uses its mapped
signal and payload as the first resumed Result, and reaches success
without replaying the first request or duplicating a source entry.
NoopCheckpoint refuses the same suspended-run claim.

<!-- field: use_case -->
#### Use case

[rel14.0-uc001-request-driven-signal](../use-cases/rel14.0-uc001-request-driven-signal.md)

<!-- field: go_test -->
#### Go test

go test ./cmd/agent -run TestRequestSignalSource_SuspendResume

<!-- field: inputs | form: "tree" -->
#### Inputs

- requests:
  1. start
  1. continue
- backends:
  1. in_memory
  1. dolt
  1. noop

<!-- field: expected | form: "tree" -->
#### Expected

- accepted_source_entries: 2
- replayed_requests: 0
- noop_refused: true

<!-- field: traces -->
#### Traces

- [rel14.0-uc001-request-driven-signal](../use-cases/rel14.0-uc001-request-driven-signal.md) [S4](../use-cases/rel14.0-uc001-request-driven-signal.md#s4)
- [srd046-request-driven-signal-source](../software-requirements/srd046-request-driven-signal-source.md) [AC5](../software-requirements/srd046-request-driven-signal-source.md#ac5)

<!-- case: Cancellation and failures release request ownership | status: implemented -->
### Cancellation and failures release request ownership

Success, suspension, command failure, cancellation, timeout, and
checkpoint failure each release the run's admission owner. Trace records
bounded source, signal, state, outcome, and stage fields without sensitive
payload content.

<!-- field: use_case -->
#### Use case

[rel14.0-uc001-request-driven-signal](../use-cases/rel14.0-uc001-request-driven-signal.md)

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core ./internal/tools/rest -run TestSignalSource_ReleasesOwnership

<!-- field: inputs | form: "tree" -->
#### Inputs

- exits:
  1. success
  1. suspended
  1. command_error
  1. cancelled
  1. timeout
  1. checkpoint_error

<!-- field: expected | form: "tree" -->
#### Expected

- ownership_released: 6
- sensitive_trace_values: 0

<!-- field: traces -->
#### Traces

- [rel14.0-uc001-request-driven-signal](../use-cases/rel14.0-uc001-request-driven-signal.md) [S5](../use-cases/rel14.0-uc001-request-driven-signal.md#s5)
- [srd046-request-driven-signal-source](../software-requirements/srd046-request-driven-signal-source.md) [AC6](../software-requirements/srd046-request-driven-signal-source.md#ac6)

<!-- case: Dolt declarations keep word kinds and authority separate | status: implemented -->
### Dolt declarations keep word kinds and authority separate

Shared declarations register distinct provision, query, and write inits.
Startup accepts one fixed operation per ToolDef and rejects unknown config,
kind-mismatched or multi-statement SQL, undeclared placeholders, and
caller-authority input fields.

<!-- field: use_case -->
#### Use case

[rel14.0-uc002-dolt-boundary-word](../use-cases/rel14.0-uc002-dolt-boundary-word.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/dolt ./internal/tools/catalog -run 'TestDoltDeclarations|TestDecodeConfigStrictAndValid|TestPrepareConfigRejectsUnsafeStatements|TestFactoriesStrictlySeparateOperationKinds'

<!-- field: inputs | form: "tree" -->
#### Inputs

- kinds:
  1. provision
  1. query
  1. write
- invalid:
  1. unknown_field
  1. kind_mismatch
  1. multi_statement
  1. undeclared_placeholder
  1. runtime_sql

<!-- field: expected | form: "tree" -->
#### Expected

- factories_registered: 3
- invalid_rejected: 5

<!-- field: traces -->
#### Traces

- [rel14.0-uc002-dolt-boundary-word](../use-cases/rel14.0-uc002-dolt-boundary-word.md) [S1](../use-cases/rel14.0-uc002-dolt-boundary-word.md#s1)
- [rel14.0-uc002-dolt-boundary-word](../use-cases/rel14.0-uc002-dolt-boundary-word.md) [S2](../use-cases/rel14.0-uc002-dolt-boundary-word.md#s2)
- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC1](../software-requirements/srd047-dolt-boundary-word.md#ac1)
- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC2](../software-requirements/srd047-dolt-boundary-word.md#ac2)

<!-- case: Real Dolt provisioning creates database and schema idempotently | status: implemented -->
### Real Dolt provisioning creates database and schema idempotently

Against a real Dolt server, dolt_provision creates a missing configured
database, applies ordered schema statements, and creates one schema
commit. Repeating the operation reports no schema change and adds no empty
commit. An incompatible schema fails without a commit.

<!-- field: use_case -->
#### Use case

[rel14.0-uc002-dolt-boundary-word](../use-cases/rel14.0-uc002-dolt-boundary-word.md)

<!-- field: go_test -->
#### Go test

mage integration:doltWord

<!-- field: inputs | form: "tree" -->
#### Inputs

- database_exists: false
- provision_runs: 2

<!-- field: expected | form: "tree" -->
#### Expected

- database_created: true
- schema_applied: true
- schema_commits: 1
- incompatible_schema_commits: 0

<!-- field: traces -->
#### Traces

- [rel14.0-uc002-dolt-boundary-word](../use-cases/rel14.0-uc002-dolt-boundary-word.md) [S3](../use-cases/rel14.0-uc002-dolt-boundary-word.md#s3)
- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC3](../software-requirements/srd047-dolt-boundary-word.md#ac3)
- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC7](../software-requirements/srd047-dolt-boundary-word.md#ac7)

<!-- case: Real Dolt query binds params and never commits | status: implemented -->
### Real Dolt query binds params and never commits

A configured SELECT binds a declared value and returns ordered columns and
rows with row_count and limits. A low row or byte limit sets truncated.
Mutation-shaped query config is rejected, and query adds no Dolt commit.

<!-- field: use_case -->
#### Use case

[rel14.0-uc002-dolt-boundary-word](../use-cases/rel14.0-uc002-dolt-boundary-word.md)

<!-- field: go_test -->
#### Go test

mage integration:doltWord

<!-- field: inputs | form: "tree" -->
#### Inputs

- params:
  - status: open
- limits:
  1. rows
  1. bytes

<!-- field: expected | form: "tree" -->
#### Expected

- bound_value_used: true
- structured_rows: true
- truncation_deterministic: true
- commits_added: 0

<!-- field: traces -->
#### Traces

- [rel14.0-uc002-dolt-boundary-word](../use-cases/rel14.0-uc002-dolt-boundary-word.md) [S4](../use-cases/rel14.0-uc002-dolt-boundary-word.md#s4)
- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC4](../software-requirements/srd047-dolt-boundary-word.md#ac4)
- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC7](../software-requirements/srd047-dolt-boundary-word.md#ac7)

<!-- case: Real Dolt write creates one declared-message commit | status: implemented -->
### Real Dolt write creates one declared-message commit

A configured INSERT binds declared params in one transaction. Dolt history
gains exactly one commit whose message comes from the configured template,
and output reports rows_affected and its exact commit hash. A zero-row
update follows commit_on_no_change false and adds no commit.

<!-- field: use_case -->
#### Use case

[rel14.0-uc002-dolt-boundary-word](../use-cases/rel14.0-uc002-dolt-boundary-word.md)

<!-- field: go_test -->
#### Go test

mage integration:doltWord

<!-- field: inputs | form: "tree" -->
#### Inputs

- operations:
  1. insert_one
  1. update_zero
- commit_on_no_change: false

<!-- field: expected | form: "tree" -->
#### Expected

- insert_commits: 1
- message_matches_template: true
- hash_matches_history: true
- zero_row_commits: 0

<!-- field: traces -->
#### Traces

- [rel14.0-uc002-dolt-boundary-word](../use-cases/rel14.0-uc002-dolt-boundary-word.md) [S4](../use-cases/rel14.0-uc002-dolt-boundary-word.md#s4)
- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC5](../software-requirements/srd047-dolt-boundary-word.md#ac5)
- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC7](../software-requirements/srd047-dolt-boundary-word.md#ac7)

<!-- case: Caller authority and failed writes create no commit | status: implemented -->
### Caller authority and failed writes create no commit

Inputs containing SQL, connection, database, schema, or commit_message are each
refused before statement execution. A validly shaped request whose
configured statement violates a constraint rolls back and leaves Dolt
history unchanged.

<!-- field: use_case -->
#### Use case

[rel14.0-uc002-dolt-boundary-word](../use-cases/rel14.0-uc002-dolt-boundary-word.md)

<!-- field: go_test -->
#### Go test

mage integration:doltWord

<!-- field: inputs | form: "tree" -->
#### Inputs

- forbidden_fields:
  1. sql
  1. connection
  1. database
  1. schema
  1. commit_message
- failure: constraint_violation

<!-- field: expected | form: "tree" -->
#### Expected

- refused_before_execution: 5
- failed_write_commits: 0

<!-- field: traces -->
#### Traces

- [rel14.0-uc002-dolt-boundary-word](../use-cases/rel14.0-uc002-dolt-boundary-word.md) [S2](../use-cases/rel14.0-uc002-dolt-boundary-word.md#s2)
- [rel14.0-uc002-dolt-boundary-word](../use-cases/rel14.0-uc002-dolt-boundary-word.md) [S6](../use-cases/rel14.0-uc002-dolt-boundary-word.md#s6)
- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC2](../software-requirements/srd047-dolt-boundary-word.md#ac2)
- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC7](../software-requirements/srd047-dolt-boundary-word.md#ac7)

<!-- case: Word and checkpoint databases remain separate | status: implemented -->
### Word and checkpoint databases remain separate

One Dolt server hosts a domain-word database and a checkpoint database.
Both operate without reading each other's tables or branches. Config that
normalizes to the same server and database is rejected at startup.

<!-- field: use_case -->
#### Use case

[rel14.0-uc002-dolt-boundary-word](../use-cases/rel14.0-uc002-dolt-boundary-word.md)

<!-- field: go_test -->
#### Go test

mage integration:doltWord

<!-- field: inputs | form: "tree" -->
#### Inputs

- valid_databases:
  1. agent_domain
  1. agent_checkpoint
- invalid_databases:
  1. shared

<!-- field: expected | form: "tree" -->
#### Expected

- separate_databases_succeed: true
- shared_database_rejected: true
- cross_table_access: false

<!-- field: traces -->
#### Traces

- [rel14.0-uc002-dolt-boundary-word](../use-cases/rel14.0-uc002-dolt-boundary-word.md) [S6](../use-cases/rel14.0-uc002-dolt-boundary-word.md#s6)
- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC7](../software-requirements/srd047-dolt-boundary-word.md#ac7)

<!-- case: Dolt signals traces errors and receipts follow boundary contracts | status: implemented -->
### Dolt signals traces errors and receipts follow boundary contracts

Success emits the word-specific signals; infrastructure faults emit
CommandError with a precise stage. Query has noop Undo. Compensatable
writes encode an opaque operation, database, and commit receipt consumed
only by their configured compensation path. Traces redact SQL, secrets,
credentials, and sensitive params.

<!-- field: use_case -->
#### Use case

[rel14.0-uc002-dolt-boundary-word](../use-cases/rel14.0-uc002-dolt-boundary-word.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/dolt -run 'TestWriteStagesAndCommitsExactlyOnceWithTrustedMessage|TestWriteNoChangeAndFailuresNeverCreateCommit|TestRuntimeRejectsAuthorityAndCancellationBeforeOpening'

<!-- field: inputs | form: "tree" -->
#### Inputs

- stages:
  1. connection_resolution
  1. statement_execute
  1. stage
  1. commit
  1. timeout
- reversibility:
  1. noop
  1. compensatable
  1. irreversible

<!-- field: expected | form: "tree" -->
#### Expected

- errors_staged: true
- receipt_opaque_to_core: true
- sensitive_trace_values: 0

<!-- field: traces -->
#### Traces

- [rel14.0-uc002-dolt-boundary-word](../use-cases/rel14.0-uc002-dolt-boundary-word.md) [S5](../use-cases/rel14.0-uc002-dolt-boundary-word.md#s5)
- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC6](../software-requirements/srd047-dolt-boundary-word.md#ac6)
