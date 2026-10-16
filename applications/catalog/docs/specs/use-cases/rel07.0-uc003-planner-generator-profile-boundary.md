---
id: rel07.0-uc003-planner-generator-profile-boundary
title: Planner to Generator Profile Boundary
status: planned
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

This library-owned tracer is retired. It synthesized planner work in Mage and
delegated to a shell child that fabricated the workspace change. The
replacement belongs to the coding application and runs the canonical planner
and executor through the real built agent binary.

Profile-owned fixtures define the requirement graph, issue/materialization
shape, workspace, and expected terminal states. agent-core supplies planning
and nested execution capabilities, but profile-owned configuration selects how
planner delegates to generator.

<!-- field: actor -->
## Actor

Profile maintainer verifying planner orchestration.

<!-- field: trigger -->
## Trigger

A maintainer runs the application-owned planner integration target:

  cd applications/coding-agent && mage integration:plannerDelegation

<!-- field: preconditions -->
## Preconditions

- agents/planner/profile.yaml and agents/executor/profile.yaml validate.
- testdata/integration/rel07-planner-generator contains a minimal graph, workspace, and expected outcome records.
- External issue tracker writes are disabled, mocked, or directed at local fixtures.

<!-- field: flow -->
## Flow

- F1: The target loads a small requirement graph fixture for planner.
- F2: Planner extracts or materializes one executable task.
- F3: Planner invokes generator as the child profile for the materialized task.
- F4: Generator runs in a fixture workspace and returns a child run outcome.
- F5: Planner records graph terminal state, task outcome, and validation evidence.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | sm-agent-role-taxonomy planner | Planner owns requirement graph, issue materialization, and child generator execution. |
| T2 | sm-agent-role-taxonomy generator | Generator executes the materialized task in the child workspace. |
| T3 | sm-multi-agent-profile-system planner->generator | Planner delegates task execution to generator. |
| T4 | test-rel07.0-profile-integrations planner_generator_boundary | The suite records target, fixtures, and terminal graph evidence. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The planner fixture records a graph or task input with no external side effects.

<!-- field: traces -->
#### Traces

- sm-agent-role-taxonomy planner inputs

<!-- sc: S2 -->
### S2

Planner invokes generator as a child profile through profile-owned configuration.

<!-- field: traces -->
#### Traces

- sm-multi-agent-profile-system dependency_directions planner->generator

<!-- sc: S3 -->
### S3

The tracer records explicit graph terminal state and child run outcome.

<!-- field: traces -->
#### Traces

- sm-agent-role-taxonomy planner outputs

<!-- field: out_of_scope -->
## Out of scope

- Creating real GitHub issues from the tracer.
- Exercising a multi-task planning campaign.
- Changing planner tool implementations in agent-core.

<!-- field: test_suite_trace -->
## Test suite trace

test-rel07.0-profile-integrations
