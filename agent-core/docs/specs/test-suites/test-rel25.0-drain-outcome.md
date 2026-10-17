---
id: test-rel25.0-drain-outcome
title: Test Suite for Release 25.0 Drain Outcome
release: "25.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel25.0-uc001-drain-outcome

<!-- field: preconditions -->
## Preconditions

- srd033-monitor-rest-api R6.5 and R6.10 define the routing and the reported field.
- Unit cases hold a connection open past a short shutdown budget on a loopback listener, with no container and no external network.
- The release proof runs the shipped runtime-state-reader profile through cmd/agent.

<!-- field: test_cases -->
## Test cases

<!-- case: A drain cut short by the budget still returns success | status: implemented -->
### A drain cut short by the budget still returns success

A handler holds a request open past a shutdown budget of a few milliseconds; the stop returns without a command error and the listener is closed.

<!-- field: use_case -->
#### Use case

rel25.0-uc001-drain-outcome

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest/... -run 'TestServerStopIncompleteDrainSucceeds'

<!-- field: expected | form: "tree" -->
#### Expected

- stop_returns_no_command_error: true
- listener_refuses_new_connections: true
- server_deregistered: true

<!-- field: traces -->
#### Traces

- rel25.0-uc001-drain-outcome S1
- srd033-monitor-rest-api AC6
- srd033-monitor-rest-api R6.5

<!-- case: The stop output names the drain outcome in both directions | status: implemented -->
### The stop output names the drain outcome in both directions

An idle server and a server with a held connection stop through the same path; the declared output carries the drain field in both, differing only in its value.

<!-- field: use_case -->
#### Use case

rel25.0-uc001-drain-outcome

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest/... -run 'TestServerStopReportsDrainOutcome'

<!-- field: expected | form: "tree" -->
#### Expected

- completed_drain_reported: true
- incomplete_drain_reported: true
- field_shape_identical: true

<!-- field: traces -->
#### Traces

- rel25.0-uc001-drain-outcome S2
- srd033-monitor-rest-api AC6
- srd033-monitor-rest-api R6.10

<!-- case: Faults other than the expired budget stay command failures | status: implemented -->
### Faults other than the expired budget stay command failures

Stopping a server that was never launched, and a shutdown returning an error unrelated to the deadline, each return a CommandError naming the fault.

<!-- field: use_case -->
#### Use case

rel25.0-uc001-drain-outcome

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest/... -run 'TestServerStopNonDrainFaultsFail'

<!-- field: expected | form: "tree" -->
#### Expected

- unlaunched_server_fails: true
- unrelated_shutdown_error_fails: true
- fault_named_in_error: true

<!-- field: traces -->
#### Traces

- rel25.0-uc001-drain-outcome S3
- srd033-monitor-rest-api AC6
- srd033-monitor-rest-api R6.5

<!-- case: The declared stop output matches the runtime field | status: implemented -->
### The declared stop output matches the runtime field

The drain field is declared in the stop word's output schema, so the existing declaration-to-runtime gate covers it rather than a hand-written assertion.

<!-- field: use_case -->
#### Use case

rel25.0-uc001-drain-outcome

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest/... -run 'TestRESTDeclaredOutputPropertiesExistAtRuntime'

<!-- field: expected | form: "tree" -->
#### Expected

- drain_field_declared: true
- declaration_matches_runtime: true

<!-- field: traces -->
#### Traces

- rel25.0-uc001-drain-outcome S2
- srd033-monitor-rest-api AC6
- srd033-monitor-rest-api R6.10

<!-- case: The monitor release proof reports succeeded on its shipped budget | status: implemented -->
### The monitor release proof reports succeeded on its shipped budget

The runtime-state-reader profile runs end to end through cmd/agent with the shutdown budget it ships, and the run status is succeeded whether or not the drain finished.

<!-- field: use_case -->
#### Use case

rel25.0-uc001-drain-outcome

<!-- field: go_test -->
#### Go test

mage integration:monitor

<!-- field: expected | form: "tree" -->
#### Expected

- run_status_succeeded: true
- shipped_budget_unchanged: true

<!-- field: traces -->
#### Traces

- rel25.0-uc001-drain-outcome S4
- srd033-monitor-rest-api AC6
- srd033-monitor-rest-api R6.5
