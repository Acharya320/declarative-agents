---
id: rel11.0-uc002-mock-serves-fixtures
title: One Mock Profile Stands In for Any REST Dependency and Records What It Received
status: in-progress
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A scenario author needs a subject's REST dependency to return controlled
responses, including a failure on a chosen call. They write a fixture file
and start the shipped mock profile against it. The mock serves the fixture's
routes, scripts the declared sequence, and records every request it received
on a log route. A validator reads that log to assert the subject called the
dependency with the expected method, path, and body. Serving a second
dependency in the same scenario needs another fixture and another mock
instance, not another profile and not a container image.

<!-- field: actor -->
## Actor

A scenario author controlling a subject agent's REST dependencies.

<!-- field: trigger -->
## Trigger

A scenario declares a dependency fixture, and the scenario critic starts a mock
against it before starting the subject.

<!-- field: preconditions -->
## Preconditions

- agent-core provides the mock endpoint binding, fixture format, and request log (srd039).
- The mock profile is shipped in applications/catalog and selects its fixture by environment.

<!-- field: flow -->
## Flow

- F1: The mock starts with a fixture path and listen address supplied by environment, loads and validates the fixture, and serves its health route.
- F2: The subject calls the mock; the mock answers the matching route with the next response in that route's declared sequence, repeating the last once the sequence is exhausted.
- F3: A call matching no route receives a 404 and is recorded as a miss.
- F4: The mock records every received request with its method, path, headers, and body, and serves the ordered log on its log route.
- F5: A validator reads the log route and asserts which calls the subject made, in what order, with what bodies.
- F6: A second mock instance in the same scenario serves a different dependency from a different fixture on a different address, from the same profile.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd019-mock R1 | The mock serves through the agent-core mock binding, exposes health, and exits on a lifecycle signal without calling a model. |
| T2 | srd019-mock R2 | The fixture is selected by environment, so one profile serves any dependency and several mocks run per scenario. |
| T3 | srd019-mock R3 | The request log is readable, per-instance, and sufficient to verify outbound calls including misses. |
| T4 | srd007-rest R1 | The mock is a REST server profile over the shared server runtime. |
| T5 | sm-agent-role-taxonomy | The mock family enters the catalog as the mock stand-in for a dependency. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

One shipped mock profile serves two different dependencies in one scenario, differing only by fixture and address.

<!-- field: traces -->
#### Traces

- srd019-mock AC2

<!-- sc: S2 -->
### S2

A route declaring an ordered response sequence returns them in order and repeats the last, so a failure can be scripted on a chosen call.

<!-- field: traces -->
#### Traces

- srd019-mock AC1

<!-- sc: S3 -->
### S3

A validator reading the log route can assert the subject's calls, distinguishing a correct call from a wrong or missing one, and seeing unmatched calls.

<!-- field: traces -->
#### Traces

- srd019-mock AC3

<!-- sc: S4 -->
### S4

A malformed fixture prevents the mock from serving, failing at configuration validation and startup rather than mid-scenario.

<!-- field: traces -->
#### Traces

- srd019-mock AC2

<!-- sc: S5 -->
### S5

A fresh mock instance starts with an empty log, so one scenario never observes another's requests.

<!-- field: traces -->
#### Traces

- srd019-mock AC3

<!-- field: failure_scenarios -->
## Failure scenarios

| id | scenario | expected |
| --- | --- | --- |
| E1 | The mock profile grows per-dependency response content so a new dependency needs a profile change. | Rejected; response content lives in scenario fixtures (srd019 R2.2). |
| E2 | A malformed fixture is accepted at startup and fails on the first request, mid-scenario. | Rejected; fixture validation blocks serving (srd019 R2.4). |
| E3 | The mock asserts an expectation itself and fails the run. | Rejected; the mock serves and records, and the validator asserts (srd019 non-goals). |

<!-- field: out_of_scope -->
## Out of scope

- The fixture schema, matching rules, and log payload shape, which are srd039 in agent-core.
- Proxying to, recording from, or replaying a real upstream service.
- Refreshing fixtures from observed real-service responses, which real-stack verification owns.

<!-- field: evidence_gaps -->
## Evidence gaps

- Current conformance reads request-log entries directly; it does not execute a validator that passes for a correct call and fails for a wrong or missing call.
- Malformed fixtures are rejected through --validate-config, but no current test starts serve mode and asserts startup rejection before any request.

<!-- field: test_suite_trace -->
## Test suite trace

test-rel11.0-scenario-critic-rig
