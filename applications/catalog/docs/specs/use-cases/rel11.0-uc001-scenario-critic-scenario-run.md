---
id: rel11.0-uc001-scenario-critic-scenario-run
title: The Scenario Critic Runs an Agent's Scenarios to a Verdict Without Live Services
status: in-progress
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: historical_aliases -->
## Historical aliases

- rel11.0-uc001-assembler-scenario-run

<!-- field: summary -->
## Summary

A profile author adds a tests/<scenario>/ directory beside an agent and runs
the scenario critic. The scenario critic discovers the scenario, starts the mocks its
fixtures declare, starts the subject with its REST dependency pointed at a
mock, awaits health, dispatches each validator machine visibly, and reduces
the ordered join to a verdict. The subject gets a fresh instance per
scenario, every child is torn down, and no live Chroma, Ollama, or
Kubernetes is involved. A second scenario over the same subject scripts a
dependency failure and asserts the subject's clean failure path.

<!-- field: actor -->
## Actor

A profile author testing an agent's machine and REST tools.

<!-- field: trigger -->
## Trigger

A profile author runs the scenario critic over discovery roots containing an agent
that carries one or more scenario directories.

<!-- field: preconditions -->
## Preconditions

- The scenario critic and mock profiles are shipped in applications/catalog.
- agent-core provides the service-child, validator-execution, and discovery words (srd040) and the mock endpoint binding (srd039).
- The subject agent's REST dependency base_url is environment-expandable.

<!-- field: flow -->
## Flow

- F1: The scenario critic enumerates scenarios under its configured roots and takes the first, which names the subject, its validators, and its fixture files.
- F2: A labeled command-state fixture array drives sequential MachineSpec for_each dispatch; each item starts one mock on its own address, records its own receipt-bearing Entry, and joins in input order.
- F3: The scenario critic starts the subject with environment additions pointing its dependency base_url at the mock addresses, and awaits its health.
- F4: A labeled validator array drives sequential MachineSpec for_each dispatch; each child execution drives the subject over HTTP and records its own receipt-bearing Entry before an ordered join.
- F5: A thin word records joined validator outcomes and derives the verdict; the scenario critic then lists the subject before its mocks and iterates stop_service once per child before joining teardown.
- F6: The scenario critic repeats for the subject's remaining scenarios with a fresh subject each time, then reports the aggregate verdict as its terminal state.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd018-scenario-critic R1 | The pipeline composes, verdicts, and tears down each scenario with a fresh subject and no model call. |
| T2 | srd018-scenario-critic R2 | One subject carries several self-contained scenarios as profile-owned data. |
| T3 | srd018-scenario-critic R3 | Scenarios are discovered by convention under configurable roots. |
| T4 | srd018-scenario-critic R4 | The subject reaches its mocks through declared environment additions over base_url expansion, with its behavior unchanged. |
| T5 | srd018-scenario-critic R6 | The verdict is the conjunction of validator terminals, and a failure names its cause. |
| T6 | srd019-mock R2 | Each mock serves the scenario's fixture selected at start, so several mocks run from one profile. |
| T7 | sm-agent-role-taxonomy | The scenario-critic joins the realization inventory as a Critic-bound test harness, while mock remains a non-role test double. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

A scenario runs end to end with no live external service; fixture starts, validator executions, and child stops are ordered machine-visible item dispatches with separate checkpoint, receipt, and rollback boundaries.

<!-- field: traces -->
#### Traces

- srd018-scenario-critic AC1
- srd019-mock AC2

<!-- sc: S2 -->
### S2

Adding a second scenario to the same subject requires only new files in a tests/ directory, with no change to the scenario critic profile.

<!-- field: traces -->
#### Traces

- srd018-scenario-critic AC2

<!-- sc: S3 -->
### S3

A scenario whose fixture scripts a dependency failure drives the subject's failure path, and the scenario passes because the validator's expectation was met.

<!-- field: traces -->
#### Traces

- srd018-scenario-critic AC6

<!-- sc: S4 -->
### S4

A scenario whose validator expectation is deliberately wrong makes the scenario critic reach Failed and name the failing validator.

<!-- field: traces -->
#### Traces

- srd018-scenario-critic AC6

<!-- sc: S5 -->
### S5

No child process survives the ordered subject-first stop iterator, and a second consecutive run passes without port or state leakage.

<!-- field: traces -->
#### Traces

- srd018-scenario-critic AC1

<!-- sc: S6 -->
### S6

The same scenario critic discovers and runs scenarios under the applications tree with no special-casing.

<!-- field: traces -->
#### Traces

- srd018-scenario-critic AC3

<!-- field: failure_scenarios -->
## Failure scenarios

| id | scenario | expected |
| --- | --- | --- |
| E1 | A scenario's mock fails to start and the scenario critic proceeds to run validators against a subject with a dead dependency. | Rejected; a failed start or health timeout routes to teardown and the scenario reports Failed with the cause. |
| E2 | Two scenarios share one subject instance and the second observes state the first left behind. | Rejected; each scenario gets a fresh subject (srd018 R1.3). |
| E3 | A scenario patches the subject's machine or tool selections to make a test pass. | Rejected; only addresses and other harness-controlled fields may differ (srd018 R4.3). |

<!-- field: out_of_scope -->
## Out of scope

- The agent-core words themselves and the mock fixture schema; those are srd040 and srd039.
- Kubernetes, containers, or a hierarchical control plane; children are local processes at this release.
- Span-based state assertion; responses and the mock request log are the witnesses.
- Authoring scenarios for the existing mesh agents, which is application-owned work.
- Bounded parallel validator iteration; sequential dispatch is intentional until MachineSpec supports a bounded parallel mode.

<!-- field: evidence_gaps -->
## Evidence gaps

- Current rig evidence asserts ordered scenario and aggregate verdicts but not the failing validator's name or reason.
- Child-stop item spans and two consecutive runs show dispatch and repeatability but do not inspect operating-system process survival.
- No current test induces mock-start failure and then verifies failure-path teardown.

<!-- field: test_suite_trace -->
## Test suite trace

test-rel11.0-scenario-critic-rig
