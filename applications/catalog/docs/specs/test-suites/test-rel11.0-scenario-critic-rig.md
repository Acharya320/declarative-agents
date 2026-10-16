---
id: test-rel11.0-scenario-critic-rig
title: Test Suite for Release 11.0 Scenario Critic Test Rig
release: "11.0"
status: in-progress
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: historical_aliases -->
## Historical aliases

- test-rel11.0-assembler-rig

<!-- field: overview -->
## Overview

This suite validates the scenario critic test rig: one shipped rig profile that
composes a scenario from profile-owned data, and one shipped mock profile
that stands in for any REST dependency. GH-668 delivered the mock, GH-669
the scenario critic profile, and GH-670 the reference subject, its scenarios,
TestRigSelfProof, and mage rig. GH-671 added the chatbot-mesh rag-server
scenario and cross-root proof. Every cited command now exists.

Evidence remains narrower than the complete requirements. The current tests
inspect request-log contents directly rather than running a validator over
correct and incorrect calls; malformed fixtures are exercised only through
--validate-config; and repeatability plus stop-dispatch spans do not prove
process non-survival or mock-start-failure teardown. Those behaviors remain
specified, not claimed by the implemented cases below.

<!-- field: traces -->
## Traces

- rel11.0-uc001-scenario-critic-scenario-run
- rel11.0-uc002-mock-serves-fixtures

<!-- field: preconditions -->
## Preconditions

- agent-core provides the mock endpoint binding (srd039) and the service-child, validator-execution, and discovery words (srd040).
- agent-core is available as a sibling ../../agent-core checkout or through demo.yaml core_root for mage validate and the rig target.
- No live Chroma, Ollama, Docker, or Kubernetes is required by any case in this suite; a case that needs one is misfiled.

<!-- field: test_cases -->
## Test cases

<!-- case: One mock profile serves a fixture and scripts a sequence | status: implemented -->
### One mock profile serves a fixture and scripts a sequence

The shipped mock profile starts with a fixture path and address supplied
by environment, serves the fixture's routes, returns a route's ordered
responses in sequence, repeats the last once exhausted, and 404s an
unmatched route.

<!-- field: use_case -->
#### Use case

rel11.0-uc002-mock-serves-fixtures

<!-- field: integration_target -->
#### Integration target

mage rig

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestMockServesFixtureSequence

<!-- field: inputs -->
#### Inputs

```yaml
fixture_route: "POST /deploy with responses [200, 500]"
requests: [POST /deploy, POST /deploy, POST /deploy, GET /nope]
```

<!-- field: expected -->
#### Expected

```yaml
statuses: [200, 500, 500, 404]
profile_unchanged_between_dependencies: true
```

<!-- field: traces -->
#### Traces

- rel11.0-uc002-mock-serves-fixtures S2
- srd019-mock AC1

<!-- case: Two mocks from one profile serve different dependencies | status: implemented -->
### Two mocks from one profile serve different dependencies

Two mock instances start from the same shipped profile with different
fixtures and addresses, each serving its own dependency's routes, with
neither instance observing the other's requests.

<!-- field: use_case -->
#### Use case

rel11.0-uc002-mock-serves-fixtures

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestMockTwoInstancesDifferentFixtures

<!-- field: inputs -->
#### Inputs

```yaml
mocks: 2
fixtures: [dep_a.yaml, dep_b.yaml]
```

<!-- field: expected -->
#### Expected

```yaml
one_profile_two_dependencies: true
logs_are_per_instance: true
```

<!-- field: traces -->
#### Traces

- rel11.0-uc002-mock-serves-fixtures S1
- rel11.0-uc002-mock-serves-fixtures S5
- srd019-mock AC2
- srd019-mock AC3

<!-- case: Request log exposes outbound calls and misses | status: implemented -->
### Request log exposes outbound calls and misses

The test calls the mock directly and inspects its log route. It asserts
that received requests retain order, method, path, body, and matched
state, including an unmatched call recorded as a miss. It does not run a
validator or prove correct-versus-wrong validator outcomes.

<!-- field: use_case -->
#### Use case

rel11.0-uc002-mock-serves-fixtures

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestMockRequestLogRecordsCalls

<!-- field: inputs -->
#### Inputs

```yaml
subject_calls: [POST /deploy body A]
```

<!-- field: expected -->
#### Expected

```yaml
order_method_path_and_body_visible: true
miss_visible_in_log: true
```

<!-- field: traces -->
#### Traces

- srd019-mock R3.1
- srd019-mock R3.3

<!-- case: Configuration validation rejects malformed fixtures | status: implemented -->
### Configuration validation rejects malformed fixtures

The test runs the shipped mock profile with --validate-config and asserts
a non-zero exit for unknown methods, empty response lists, and duplicate
routes. It does not start a serving process or probe an endpoint.

<!-- field: use_case -->
#### Use case

rel11.0-uc002-mock-serves-fixtures

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestMockMalformedFixtureFailsStartup

<!-- field: inputs -->
#### Inputs

```yaml
fixtures: [unknown_method, empty_responses, duplicate_route]
```

<!-- field: expected -->
#### Expected

```yaml
validate_config_fails: true
```

<!-- field: traces -->
#### Traces

- srd019-mock R2.4

<!-- case: Explicit authority gates a Kubernetes-addressable Ollama mock | status: implemented -->
### Explicit authority gates a Kubernetes-addressable Ollama mock

The canonical mock stays on loopback by default. With an explicit
public-listener authorization it binds a wildcard address, serves native
health and request-log routes, and answers Ollama-compatible GET
/api/tags plus three ordered POST /api/chat responses. A wildcard address
with absent or malformed authorization fails configuration validation.

<!-- field: use_case -->
#### Use case

rel11.0-uc002-mock-serves-fixtures

<!-- field: go_test -->
#### Go test

go test ./conformance -run 'TestMockAuthorizedPublicListenerServesOllamaFixture|TestMockPublicListenerRequiresExplicitAuthority'

<!-- field: inputs -->
#### Inputs

```yaml
address: 0.0.0.0 with an ephemeral port
fixture_routes: [GET /api/tags, POST /api/chat x3]
authority_cases: [true, absent, malformed]
```

<!-- field: expected -->
#### Expected

```yaml
authorized_listener_serves: true
absent_authority_rejected: true
malformed_authority_rejected: true
request_log_entries: 4
```

<!-- field: traces -->
#### Traces

- srd019-mock R1.5
- srd019-mock R2.5
- srd039-mock-endpoint-binding R1.6

<!-- case: Scenario Critic runs one scenario end to end to a verdict | status: implemented -->
### Scenario Critic runs one scenario end to end to a verdict

The scenario critic discovers a scenario, starts its mocks, starts the subject
with its dependency base_url pointed at a mock through declared
environment additions, awaits health, runs the validators, collects the
verdict, and tears down. No live external service participates.

<!-- field: use_case -->
#### Use case

rel11.0-uc001-scenario-critic-scenario-run

<!-- field: integration_target -->
#### Integration target

mage rig

<!-- field: inputs -->
#### Inputs

```yaml
scenario: reference subject with one REST dependency, one mock, one validator
```

<!-- field: expected -->
#### Expected

```yaml
verdict: Passed
no_live_services_used: true
subject_behavior_unpatched: true
```

<!-- field: traces -->
#### Traces

- rel11.0-uc001-scenario-critic-scenario-run S1
- srd018-scenario-critic AC1
- srd018-scenario-critic AC4

<!-- case: Fixture array drives machine-visible sequential mock starts | status: implemented -->
### Fixture array drives machine-visible sequential mock starts

A labeled next_scenario output supplies fixture objects to MachineSpec.
The iterator starts one mock per object in input order through the normal
registry and dispatch path, records one receipt-bearing Entry per mock,
emits a labeled ordered join, and resumes after the last persisted item
without replay. The former start_scenario_mocks Go collection loop is gone.

<!-- field: use_case -->
#### Use case

rel11.0-uc001-scenario-critic-scenario-run

<!-- field: go_test -->
#### Go test

cd ../../agent-core && go test ./internal/runtime/core ./internal/tools/service -run '^(TestLoopForEach.*|TestScenarioSteps_ThreadsMockURLIntoSubject)$'

<!-- field: inputs -->
#### Inputs

```yaml
fixtures: [dep_a.yaml, dep_b.yaml]
mode: sequential
```

<!-- field: expected -->
#### Expected

```yaml
dispatch_order: [dep_a.yaml, dep_b.yaml]
entry_per_mock: true
receipt_per_mock: true
join_is_labeled: true
resume_skips_completed_items: true
```

<!-- field: traces -->
#### Traces

- rel11.0-uc001-scenario-critic-scenario-run S1
- srd018-scenario-critic AC1

<!-- case: One subject carries several scenarios added as data only | status: implemented -->
### One subject carries several scenarios added as data only

The reference subject carries a happy-path scenario and a
dependency-failure scenario. Adding the second required only files under
tests/, with no change to the scenario critic profile, and each scenario runs
against a fresh subject instance.

<!-- field: use_case -->
#### Use case

rel11.0-uc001-scenario-critic-scenario-run

<!-- field: integration_target -->
#### Integration target

mage rig

<!-- field: inputs -->
#### Inputs

```yaml
scenarios: [happy-path, dep-failure]
```

<!-- field: expected -->
#### Expected

```yaml
scenario_critic_profile_unchanged: true
fresh_subject_per_scenario: true
both_scenarios_pass: true
```

<!-- field: traces -->
#### Traces

- rel11.0-uc001-scenario-critic-scenario-run S2
- srd018-scenario-critic AC2

<!-- case: Validator and teardown collections are machine-visible | status: implemented -->
### Validator and teardown collections are machine-visible

The current scenario exposes ordered validator items. Sequential
MachineSpec iteration dispatches one child execution per item, joins and
records outcomes, then lists the subject before its mocks and dispatches
stop_service once per child. Bounded parallel iteration is deferred, so
no service tool restores hidden validator goroutines.

<!-- field: use_case -->
#### Use case

rel11.0-uc001-scenario-critic-scenario-run

<!-- field: go_test -->
#### Go test

cd ../../agent-core && go test ./internal/tools/service -run '^(TestRunOneValidator_TerminalStateAndTimeout|TestRunScenarioValidator_EmitsItemReceipt|TestRecordScenarioValidators_PreservesJoinOrder|TestScenarioSteps_TeardownRunsOnFailurePath)$'

<!-- field: integration_target -->
#### Integration target

mage rig

<!-- field: inputs -->
#### Inputs

```yaml
validators: [first, second]
teardown_order: [subject, mocks_in_start_order]
mode: sequential
validator_otlp_endpoint: declared run_scenario_validator tool configuration
```

<!-- field: expected -->
#### Expected

```yaml
entry_per_validator: true
entry_per_child_stop: true
joins_are_labeled_and_ordered: true
no_hidden_collection_loop: true
validator_spans_exported_to_declared_otlp_endpoint: true
```

<!-- field: traces -->
#### Traces

- rel11.0-uc001-scenario-critic-scenario-run S1
- rel11.0-uc001-scenario-critic-scenario-run S5
- srd018-scenario-critic AC1

<!-- case: Subject readiness is a machine-visible REST retry loop | status: implemented -->
### Subject readiness is a machine-visible REST retry loop

StartingSubject labels the runtime health authority, path, and start time.
ProbingHealth performs one configured REST request. Unavailable status or
transport error enters WaitingToRetry, where one cancellation-aware delay
either returns to ProbingHealth or emits HealthTimeout. Timeout and global
budget exhaustion route through the ordinary verdict and teardown words.

<!-- field: use_case -->
#### Use case

rel11.0-uc001-scenario-critic-scenario-run

<!-- field: go_test -->
#### Go test

cd ../../agent-core && go test ./internal/tools/lifecycle ./internal/tools/rest/client ./internal/tools/service -run '^(TestDelay_ElapsedAndDeadlineSignals|TestDelay_CancellationAndValidation|TestRESTClient_CommandStateSelectsTrustedBaseURL|TestSubjectHealthTarget_PreservesRelativeAndAlternateListenerPaths)$'

<!-- field: integration_target -->
#### Integration target

mage rig

<!-- field: inputs -->
#### Inputs

```yaml
target_source: command_state subject label
health_paths: [relative child route, absolute alternate listener]
retry_interval: 100ms
startup_deadline: 30s
```

<!-- field: expected -->
#### Expected

```yaml
requests_per_probe_dispatch: 1
retry_states_visible: true
healthy_timeout_and_error_signals_distinct: true
cancellation_bounded: true
teardown_on_exhaustion: true
```

<!-- field: traces -->
#### Traces

- rel11.0-uc001-scenario-critic-scenario-run S1
- srd018-scenario-critic R1.8
- srd040-service-child-tools AC2

<!-- case: Scripted dependency-failure scenario reaches a passing verdict | status: implemented -->
### Scripted dependency-failure scenario reaches a passing verdict

TestRigSelfProof asserts that the dep-failure scenario is the first of two
passing scenario verdicts in sorted discovery order. The test does not
separately assert the subject terminal or validator explanation.

<!-- field: use_case -->
#### Use case

rel11.0-uc001-scenario-critic-scenario-run

<!-- field: integration_target -->
#### Integration target

mage rig

<!-- field: inputs -->
#### Inputs

```yaml
fixture_sequence: [500]
```

<!-- field: expected -->
#### Expected

```yaml
scenario_verdict: Passed
```

<!-- field: traces -->
#### Traces

- srd018-scenario-critic R6.3

<!-- case: A wrong expectation makes the scenario and aggregate fail | status: implemented -->
### A wrong expectation makes the scenario and aggregate fail

A scenario carrying a deliberately wrong validator expectation drives the
scenario verdict and aggregate terminal to Failed. TestRigSelfProof
asserts those failures; it does not assert a validator name or reason.

<!-- field: use_case -->
#### Use case

rel11.0-uc001-scenario-critic-scenario-run

<!-- field: integration_target -->
#### Integration target

mage rig

<!-- field: inputs -->
#### Inputs

```yaml
scenario: broken, with an expectation the subject does not satisfy
```

<!-- field: expected -->
#### Expected

```yaml
scenario_critic_terminal: Failed
broken_scenario_verdict: Failed
```

<!-- field: traces -->
#### Traces

- srd018-scenario-critic R6.3

<!-- case: Stop dispatches are visible and consecutive runs are repeatable | status: implemented -->
### Stop dispatches are visible and consecutive runs are repeatable

TestRigSelfProof observes child-stop item spans and completes two
consecutive full runs without a port or state collision. It does not
inspect operating-system process survival or induce a mock-start failure.

<!-- field: use_case -->
#### Use case

rel11.0-uc001-scenario-critic-scenario-run

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestRigSelfProof

<!-- field: inputs -->
#### Inputs

```yaml
runs: 2
```

<!-- field: expected -->
#### Expected

```yaml
child_stop_dispatches_visible: true
second_run_passes: true
```

<!-- field: traces -->
#### Traces

- srd018-scenario-critic R1.9

<!-- case: Discovery covers the applications tree with no special-casing | status: implemented -->
### Discovery covers the applications tree with no special-casing

With roots configured to cover both the agent families and the examples
tree, discovery enumerates scenarios from both, sorted and deterministic,
with no absolute paths and no per-area rig.

<!-- field: use_case -->
#### Use case

rel11.0-uc001-scenario-critic-scenario-run

<!-- field: integration_target -->
#### Integration target

mage integration:rig (run from applications/chatbot-mesh)

<!-- field: inputs -->
#### Inputs

```yaml
roots: [agents, examples agents]
```

<!-- field: expected -->
#### Expected

```yaml
scenarios_found_in_both_roots: true
deterministic_sorted: true
no_absolute_paths: true
```

<!-- field: traces -->
#### Traces

- rel11.0-uc001-scenario-critic-scenario-run S6
- srd018-scenario-critic AC3

<!-- case: Scenario Critic and mock profiles validate and boot as shipped profiles | status: implemented -->
### Scenario Critic and mock profiles validate and boot as shipped profiles

The scenario critic and mock profiles load and validate against agent-core and
pass the shipped-profile boot smoke alongside the existing families, so
the role catalog grows without breaking profile wiring.

<!-- field: use_case -->
#### Use case

rel11.0-uc001-scenario-critic-scenario-run

<!-- field: integration_target -->
#### Integration target

mage validate

<!-- field: inputs -->
#### Inputs

```yaml
profiles: [agents/scenario-critic, agents/mock]
```

<!-- field: expected -->
#### Expected

```yaml
mage_validate_passes: true
boot_smoke_covers_new_profiles: true
```

<!-- field: traces -->
#### Traces

- srd018-scenario-critic AC5
- srd019-mock AC4
- srd019-mock AC5
