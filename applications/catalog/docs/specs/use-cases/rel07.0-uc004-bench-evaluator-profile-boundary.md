---
id: rel07.0-uc004-bench-evaluator-profile-boundary
title: Bench to Critic Profile Boundary
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

The bench profile proof covers the human-operated evaluation boundary: a
local operator action reaches the shipped bench REST/machine surface, bench
launches critic through configured vocabulary, and the result records launch
status plus the critic child arguments. This is a launch tracer, not a full
critic benchmark.

The proof stays in applications/catalog because the bench profile, critic profile
reference, launch action, and expected dispatch are profile-owned behavior.
agent-core owns generic REST, evaluation artifact queries, and child-agent
execution; it does not own bench routes, assets, workflow, or the concrete
critic profile.

<!-- field: actor -->
## Actor

Human operator using the bench profile through a local harness.

<!-- field: trigger -->
## Trigger

A maintainer runs the bench-evaluator integration target:

  mage integration:benchEvaluator

<!-- field: preconditions -->
## Preconditions

- agents/bench/profile.yaml and agents/critic/profile.yaml validate.
- A compatible agent-core checkout is available.
- The harness reserves a loopback address and supplies a controlled child binary that records arguments.

<!-- field: flow -->
## Flow

- F1: The target launches agents/bench/profile.yaml through the built agent-core runtime with only its listen address patched.
- F2: A fixture operator POSTs a critic run request to /api/v1/experiments.
- F3: The experiments capability starts the critic profile and suite as a detached child through start_service and returns its handle while the child runs.
- F4: The controlled child records the profile, request, and output arguments supplied by bench.
- F5: OpenTelemetry spans prove the shipped serve lifecycle and the start_service and list_services path executed, and the host shuts down cleanly through its control exit route.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | sm-agent-role-taxonomy bench | Bench owns the human UI/action boundary and critic launch request. |
| T2 | sm-agent-role-taxonomy critic | Critic owns the launched session. |
| T3 | sm-multi-agent-profile-system bench->critic | Bench launches critic from the human-operated surface. |
| T4 | test-rel07.0-profile-integrations bench_evaluator_boundary | The suite records target, fixture request, and launch evidence. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The bench boundary accepts a fixture local launch action.

<!-- field: traces -->
#### Traces

- sm-agent-role-taxonomy bench checks

<!-- sc: S2 -->
### S2

Bench resolves critic through profile-owned configuration.

<!-- field: traces -->
#### Traces

- sm-multi-agent-profile-system dependency_directions bench->critic

<!-- sc: S3 -->
### S3

The trace records the real REST/machine path and the controlled child records critic dispatch arguments.

<!-- field: traces -->
#### Traces

- sm-agent-role-taxonomy bench outputs

<!-- sc: S4 -->
### S4

A second launch overlaps the first, the experiments listing reports both runs, and a launch at the declared bound is refused.

<!-- field: traces -->
#### Traces

- [srd006-bench](../software-requirements/srd006-bench.md) [AC3](../software-requirements/srd006-bench.md#ac3)

<!-- field: out_of_scope -->
## Out of scope

- Browser visual regression testing for the bench UI.
- Running a long critic suite.
- Replacing critic's own point/session tests.

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel07.0-profile-integrations](../test-suites/test-rel07.0-profile-integrations.md)
