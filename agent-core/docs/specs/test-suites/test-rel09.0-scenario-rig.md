---
id: test-rel09.0-scenario-rig
title: Test Suite for Release 09.0 Scenario Rig Runtime
release: "09.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- [rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md)

<!-- field: preconditions -->
## Preconditions

- srd039-mock-endpoint-binding defines the mock binding, fixture format, request log, and validation.
- srd040-service-child-tools defines scenario-specific child starts, stop_service, one-validator child execution, and scenario-session discovery; MachineSpec owns readiness, validator iteration, joins, and teardown iteration.
- GH-666 delivered the mock binding and GH-667 delivered the service-child tools; their component cases carry executable go_test evidence.
- The full end-to-end scenario tracer remains planned and keeps this release partial.

<!-- field: test_cases -->
## Test cases

<!-- case: Mock binding serves scripted responses in order | status: implemented -->
### Mock binding serves scripted responses in order

An endpoint with the mock binding serves a route's ordered responses on
successive requests, repeats the last response after the script is
exhausted, and returns 404 for an unmatched route.

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest -run TestRESTServer_MockBindingServesScriptedResponses

<!-- field: inputs | form: "tree" -->
#### Inputs

- fixture_route: POST /deploy with responses [200, 500]
- requests:
  1. POST /deploy
  1. POST /deploy
  1. POST /deploy
  1. GET /nope

<!-- field: expected | form: "tree" -->
#### Expected

- statuses:
  1. 200
  1. 500
  1. 500
  1. 404

<!-- field: traces -->
#### Traces

- [rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md) [S1](../use-cases/rel09.0-uc001-scenario-rig-run.md#s1)
- [srd039-mock-endpoint-binding](../software-requirements/srd039-mock-endpoint-binding.md) [AC1](../software-requirements/srd039-mock-endpoint-binding.md#ac1)
- [srd039-mock-endpoint-binding](../software-requirements/srd039-mock-endpoint-binding.md) [AC2](../software-requirements/srd039-mock-endpoint-binding.md#ac2)

<!-- case: Request log returns received requests in order | status: implemented -->
### Request log returns received requests in order

After a scripted request sequence, the log endpoint returns every
received request in order with method, path, headers, and body, and
flags the unmatched request as a miss.

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest -run TestRESTServer_MockBindingRequestLog

<!-- field: inputs | form: "tree" -->
#### Inputs

- requests:
  1. POST /deploy body A
  1. GET /nope

<!-- field: expected | form: "tree" -->
#### Expected

- log_length: 2
- first_entry:
  - method: POST
  - path: /deploy
  - body: A
  - matched: true
- second_entry:
  - method: GET
  - path: /nope
  - matched: false

<!-- field: traces -->
#### Traces

- [rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md) [S2](../use-cases/rel09.0-uc001-scenario-rig-run.md#s2)
- [srd039-mock-endpoint-binding](../software-requirements/srd039-mock-endpoint-binding.md) [AC3](../software-requirements/srd039-mock-endpoint-binding.md#ac3)

<!-- case: Malformed fixture fails validate-config and startup | status: implemented -->
### Malformed fixture fails validate-config and startup

Fixtures with an unknown method, an empty response list, and a duplicate
route each fail validation with a named error under --validate-config
and at server startup, before any request is served.

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest -run TestRESTServer_MockBindingFixtureValidation

<!-- field: inputs | form: "tree" -->
#### Inputs

- fixtures:
  1. unknown_method
  1. empty_responses
  1. duplicate_route

<!-- field: expected | form: "tree" -->
#### Expected

- validate_config_fails: true
- server_start_fails: true
- error_named: true

<!-- field: traces -->
#### Traces

- [srd039-mock-endpoint-binding](../software-requirements/srd039-mock-endpoint-binding.md) [AC4](../software-requirements/srd039-mock-endpoint-binding.md#ac4)

<!-- case: Concurrent requests neither skip nor duplicate script positions | status: implemented -->
### Concurrent requests neither skip nor duplicate script positions

Concurrent requests to a scripted route consume script positions exactly
once each, and a serial sequence is deterministic.

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest -run TestRESTServer_MockBindingConcurrentSequencing

<!-- field: inputs | form: "tree" -->
#### Inputs

- concurrent_requests: 8
- script_length: 8

<!-- field: expected | form: "tree" -->
#### Expected

- each_position_served_once: true

<!-- field: traces -->
#### Traces

- [srd039-mock-endpoint-binding](../software-requirements/srd039-mock-endpoint-binding.md) [AC5](../software-requirements/srd039-mock-endpoint-binding.md#ac5)

<!-- case: Serve-mode child lifecycle leaves no orphans | status: implemented -->
### Serve-mode child lifecycle leaves no orphans

The scenario child lifecycle spawns a serve-mode child with injected
environment and labels its runtime address and start time, stop_service
stops it, and no process from the child's group survives; a repeated
cycle passes. Cancelling the owning agent context also kills and reaps a
child and its grandchild through the shared subprocess transport.

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/service -run 'TestServiceChild_(StartStopNoOrphans|ParentCancellationKillsProcessGroup)'

<!-- field: inputs | form: "tree" -->
#### Inputs

- child: serve-mode agent profile with env addition
- cycles: 2

<!-- field: expected | form: "tree" -->
#### Expected

- runtime_target_reported: true
- child_group_gone_after_stop: true
- child_group_gone_after_parent_cancel: true
- second_cycle_passes: true

<!-- field: traces -->
#### Traces

- [rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md) [S3](../use-cases/rel09.0-uc001-scenario-rig-run.md#s3)
- [srd040-service-child-tools](../software-requirements/srd040-service-child-tools.md) [AC1](../software-requirements/srd040-service-child-tools.md#ac1)

<!-- case: Declarative REST health retry times out within its bound | status: implemented -->
### Declarative REST health retry times out within its bound

Each REST dispatch probes the trusted labeled child URL once. The
cancellation-aware delay emits retry-ready until its deadline, then
HealthTimeout, so MachineSpec owns the visible bounded loop.

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/lifecycle ./internal/tools/rest/client -run 'TestDelay_ElapsedAndDeadlineSignals|TestRESTClient_CommandStateSelectsTrustedBaseURL'

<!-- field: inputs | form: "tree" -->
#### Inputs

- target: reserved unserved port
- timeout: declared

<!-- field: expected | form: "tree" -->
#### Expected

- timeout_signal: true
- bounded_duration: true

<!-- field: traces -->
#### Traces

- [srd040-service-child-tools](../software-requirements/srd040-service-child-tools.md) [AC2](../software-requirements/srd040-service-child-tools.md#ac2)

<!-- case: Validator item reports terminal state and timeout | status: implemented -->
### Validator item reports terminal state and timeout

One atomic child dispatch reports its validator terminal state. A
hung validator is stopped at its item timeout and reported
as timed out; MachineSpec iteration and join own collection behavior.

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/service -run TestRunOneValidator_TerminalStateAndTimeout

<!-- field: inputs | form: "tree" -->
#### Inputs

- validators:
  1. passing
  1. hung
- item_timeout: declared

<!-- field: expected | form: "tree" -->
#### Expected

- per_item_terminal_state: true
- hung_reported_timed_out: true

<!-- field: traces -->
#### Traces

- [rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md) [S4](../use-cases/rel09.0-uc001-scenario-rig-run.md#s4)
- [srd040-service-child-tools](../software-requirements/srd040-service-child-tools.md) [AC3](../software-requirements/srd040-service-child-tools.md#ac3)
- [srd018-cli-flag-contract](../software-requirements/srd018-cli-flag-contract.md) [AC6](../software-requirements/srd018-cli-flag-contract.md#ac6)

<!-- case: Scenario session discovery is deterministic over a fixture tree | status: implemented -->
### Scenario session discovery is deterministic over a fixture tree

init_scenario_session discovers a fixture tree with two subjects and
three scenarios; next_scenario visits the expected entries once in sorted
order with subject, scenario, validator profiles, and fixture files.

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/service -run 'TestListScenarios_DeterministicDiscovery|TestScenarioSession_IteratesEachScenarioOnce'

<!-- field: inputs | form: "tree" -->
#### Inputs

- tree: two subjects, three tests/<scenario>/ directories

<!-- field: expected | form: "tree" -->
#### Expected

- entries: 3
- sorted: true
- repeat_identical: true

<!-- field: traces -->
#### Traces

- [rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md) [S5](../use-cases/rel09.0-uc001-scenario-rig-run.md#s5)
- [srd040-service-child-tools](../software-requirements/srd040-service-child-tools.md) [AC4](../software-requirements/srd040-service-child-tools.md#ac4)

<!-- case: Service tool declarations state reversibility and undo honestly | status: implemented -->
### Service tool declarations state reversibility and undo honestly

Every registered service init has a ToolDef, scenario-specific start
words compensate from their per-child receipt, and stop_service is noop
on undo. A spawn failure surfaces as a tool error signal rather than a
panic, and no word invokes a model.

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/service -run 'TestServiceDeclarationsMatchStandardInits|TestUndoStartedChildValidatesAndConsumesReceipt|TestServiceChild_StartRejectsBadSpawn'

<!-- field: inputs | form: "tree" -->
#### Inputs

- words:
  1. init_scenario_session
  1. start_scenario_mock
  1. start_scenario_subject
  1. stop_service
  1. run_scenario_validator
  1. start_service
  1. list_services
- spawn_failure_case: nonexistent profile path

<!-- field: expected | form: "tree" -->
#### Expected

- reversibility_consistent: true
- no_undo_mismatch: true
- spawn_failure_is_error_signal: true
- no_model_invocation: true

<!-- field: traces -->
#### Traces

- [srd040-service-child-tools](../software-requirements/srd040-service-child-tools.md) [AC5](../software-requirements/srd040-service-child-tools.md#ac5)

<!-- case: Generic detached start returns a handle and lists the outcome | status: implemented -->
### Generic detached start returns a handle and lists the outcome

start_service returns while its child still runs, list_services
reports the child running and then exited with its exit code, a
declared max_running bound emits ServiceLimitReached and holds when
concurrent starts race for the last slot, and undo stops only the
child its receipt names. The shared host and request service state
is proved in cmd/agent by TestRequestLocalStateSharesHostServiceState.

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/service -run 'TestStartService_DetachedReturnsAndListsExit|TestStartService_MaxRunningEmitsLimitReached|TestStartServiceUndoConsumesReceipt'

<!-- field: inputs | form: "tree" -->
#### Inputs

- child_modes:
  1. exit3slow
  1. hang
- max_running:
  1. 1
  1. 2
- concurrent_starts: 8

<!-- field: expected | form: "tree" -->
#### Expected

- returns_before_exit: true
- listed_exit_code: 3
- limit_signal: ServiceLimitReached
- concurrent_started: 2
- undo_stops_only_receipt_child: true

<!-- field: traces -->
#### Traces

- [rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md) [S6](../use-cases/rel09.0-uc001-scenario-rig-run.md#s6)
- [srd040-service-child-tools](../software-requirements/srd040-service-child-tools.md) [AC6](../software-requirements/srd040-service-child-tools.md#ac6)

<!-- case: One scenario runs end-to-end to a verdict | status: planned -->
### One scenario runs end-to-end to a verdict

A tracer over the full flow: mock from fixture, subject with dependency
URL injected, health, one passing validator asserting the mock request
log, verdict derived from the terminal state, teardown. A variant whose
fixture scripts a 500 exercises the subject failure path and still
passes via the validator's expectation.

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/service -run TestScenarioRig_EndToEndVerdict

<!-- field: inputs | form: "tree" -->
#### Inputs

- scenario: mock + subject + validator, happy and scripted-failure variants

<!-- field: expected | form: "tree" -->
#### Expected

- happy_verdict: Passed
- scripted_failure_verdict: Passed
- no_processes_after: true

<!-- field: traces -->
#### Traces

- [rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md) [S1](../use-cases/rel09.0-uc001-scenario-rig-run.md#s1)
- [rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md) [S2](../use-cases/rel09.0-uc001-scenario-rig-run.md#s2)
- [rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md) [S3](../use-cases/rel09.0-uc001-scenario-rig-run.md#s3)

<!-- case: Scenario session iterates each scenario exactly once | status: implemented -->
### Scenario session iterates each scenario exactly once

The session seeds from discovery roots and walks every scenario once in
deterministic order, then reports the work list exhausted. This is the
cursor the assembler machine loops on (srd018 R1.1).

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/service -run TestScenarioSession_IteratesEachScenarioOnce

<!-- field: inputs | form: "tree" -->
#### Inputs

- tree: two subjects, three scenarios

<!-- field: expected | form: "tree" -->
#### Expected

- visited_once_each: true
- deterministic_order: true
- exhausted_stays_exhausted: true

<!-- field: traces -->
#### Traces

- [rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md) [S5](../use-cases/rel09.0-uc001-scenario-rig-run.md#s5)
- [srd040-service-child-tools](../software-requirements/srd040-service-child-tools.md) [AC4](../software-requirements/srd040-service-child-tools.md#ac4)

<!-- case: A mock's runtime address reaches the subject's environment | status: implemented -->
### A mock's runtime address reaches the subject's environment

A mock binds a port chosen at runtime, and the subject started for that
scenario observes the mock's base URL in its environment under the
variable derived from the fixture name. This is the dependency
redirection the pipeline depends on (srd018 R4.1); without it the
per-scenario steps cannot be expressed as machine transitions.

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/service -run TestScenarioSteps_ThreadsMockURLIntoSubject

<!-- field: inputs | form: "tree" -->
#### Inputs

- fixture: mocks/dep.yaml

<!-- field: expected | form: "tree" -->
#### Expected

- env_var_derived: DEP_URL
- subject_observes_mock_url: true
- subject_profile_resolved_from_scenario: true

<!-- field: traces -->
#### Traces

- [rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md) [S1](../use-cases/rel09.0-uc001-scenario-rig-run.md#s1)
- [srd040-service-child-tools](../software-requirements/srd040-service-child-tools.md) [AC1](../software-requirements/srd040-service-child-tools.md#ac1)

<!-- case: A failed scenario tears its children down and names its cause | status: implemented -->
### A failed scenario tears its children down and names its cause

A subject that never becomes healthy drives the health-timeout signal;
the machine routes to a forced verdict and teardown, leaving nothing
running, and the session report names the cause (srd018 R1.5, R6.2).

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/service -run TestScenarioSteps_TeardownRunsOnFailurePath

<!-- field: inputs | form: "tree" -->
#### Inputs

- subject: started but never serving
- health_timeout: declared

<!-- field: expected | form: "tree" -->
#### Expected

- health_timeout_signal: true
- nothing_running_after_teardown: true
- report_names_cause: true

<!-- field: traces -->
#### Traces

- [rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md) [S3](../use-cases/rel09.0-uc001-scenario-rig-run.md#s3)
- [srd040-service-child-tools](../software-requirements/srd040-service-child-tools.md) [AC1](../software-requirements/srd040-service-child-tools.md#ac1)

<!-- case: Verdict passes only when every validator passed | status: implemented -->
### Verdict passes only when every validator passed

A scenario passes only when all its validators passed; a failure names
the failing validator and why, and the aggregate names the first failing
scenario. An empty run does not pass (srd018 R6.1, R6.2).

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/service -run TestScenarioSession_VerdictNamesItsCause

<!-- field: inputs | form: "tree" -->
#### Inputs

- scenarios:
  1. all-pass
  1. one-failing
  1. forced-reason

<!-- field: expected | form: "tree" -->
#### Expected

- pass_requires_all: true
- failure_names_validator: true
- empty_run_does_not_pass: true

<!-- field: traces -->
#### Traces

- [rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md) [S4](../use-cases/rel09.0-uc001-scenario-rig-run.md#s4)
- [srd040-service-child-tools](../software-requirements/srd040-service-child-tools.md) [AC3](../software-requirements/srd040-service-child-tools.md#ac3)
- [srd018-cli-flag-contract](../software-requirements/srd018-cli-flag-contract.md) [AC6](../software-requirements/srd018-cli-flag-contract.md#ac6)

<!-- case: A failed terminal state exits non-zero and distinctly | status: implemented -->
### A failed terminal state exits non-zero and distinctly

A run whose machine reaches a failure terminal exits with a code
distinct from the one meaning the binary could not complete a run, and a
success or suspended terminal exits zero. The mapping comes from the
terminal RunStatus, not a terminal state name, so a validator is judged
on its exit code rather than parsed log output (srd018 R6).

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md)

<!-- field: go_test -->
#### Go test

go test ./cmd/agent -run TestExitCode

<!-- field: inputs | form: "tree" -->
#### Inputs

- statuses:
  1. succeeded
  1. suspended
  1. failed
  1. budget_exceeded
  1. cancelled

<!-- field: expected | form: "tree" -->
#### Expected

- succeeded_and_suspended_exit_zero: true
- failure_terminal_exit_distinct_from_run_error: true
- mapping_from_run_status: true

<!-- field: traces -->
#### Traces

- [rel09.0-uc001-scenario-rig-run](../use-cases/rel09.0-uc001-scenario-rig-run.md) [S4](../use-cases/rel09.0-uc001-scenario-rig-run.md#s4)
- [srd040-service-child-tools](../software-requirements/srd040-service-child-tools.md) [AC3](../software-requirements/srd040-service-child-tools.md#ac3)
- [srd018-cli-flag-contract](../software-requirements/srd018-cli-flag-contract.md) [AC6](../software-requirements/srd018-cli-flag-contract.md#ac6)
