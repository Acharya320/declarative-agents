---
id: rel02.0-uc001-remote-coding-request
title: Persistent Planner Delegates to Executor and Critic
status: implemented
test_suite: test-rel02.0-serving-contract
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A caller submits one task and opaque workspace identity to the persistent
planner. The planner reuses canonical planning words, delegates execution to
the executor service, and submits the same shared volume to the critic service.
Each role is a real persistent agent profile with lifecycle health.

<!-- field: actor -->
## Actor

A delivery system submitting a coding task to a deployed coding application.

<!-- field: trigger -->
## Trigger

POST /api/v1/coding with a task and workspace_id.

<!-- field: preconditions -->
## Preconditions

- All three serving processes mount the same workspace at their trusted runtime directory.
- The assembled profile tree contains the application serving profiles and canonical library closure.
- The canonical planner and executor model is reachable.

<!-- field: flow -->
## Flow

- F1: The planner host accepts the declared request and starts its trusted request profile.
- F2: Canonical planner model and parsing words produce a structured plan.
- F3: The planner invokes the declared executor REST client; agent-core injects traceparent.
- F4: The executor endpoint runs the canonical executor profile against its trusted shared runtime directory.
- F5: After executor success, the planner invokes the declared critic REST client with the carried workspace identity.
- F6: The critic endpoint runs the canonical changed-workspace critic and maps acceptance or rejection to the declared HTTP response.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [R1.1](../software-requirements/srd001-serving-contract.md#r11), [R1.2](../software-requirements/srd001-serving-contract.md#r12), [R1.3](../software-requirements/srd001-serving-contract.md#r13) | S1 traces this requirement in the existing success criteria. |
| T2 | [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [AC1](../software-requirements/srd001-serving-contract.md#ac1) | S1 traces this requirement in the existing success criteria. |
| T3 | [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [R2.1](../software-requirements/srd001-serving-contract.md#r21), [R2.2](../software-requirements/srd001-serving-contract.md#r22), [R2.3](../software-requirements/srd001-serving-contract.md#r23), [R2.4](../software-requirements/srd001-serving-contract.md#r24) | S2 traces this requirement in the existing success criteria. |
| T4 | [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [R3.1](../software-requirements/srd001-serving-contract.md#r31), [R3.2](../software-requirements/srd001-serving-contract.md#r32), [R3.3](../software-requirements/srd001-serving-contract.md#r33) | S2 traces this requirement in the existing success criteria. |
| T5 | [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [AC2](../software-requirements/srd001-serving-contract.md#ac2) | S2 traces this requirement in the existing success criteria. |
| T6 | [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [R4.3](../software-requirements/srd001-serving-contract.md#r43) | S3 traces this requirement in the existing success criteria. |
| T7 | [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [AC3](../software-requirements/srd001-serving-contract.md#ac3) | S3 traces this requirement in the existing success criteria. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Planner, executor, and critic health endpoints report from live profile-owned listeners.

<!-- field: traces -->
#### Traces

- [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [R1.1](../software-requirements/srd001-serving-contract.md#r11)
- [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [R1.2](../software-requirements/srd001-serving-contract.md#r12)
- [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [R1.3](../software-requirements/srd001-serving-contract.md#r13)
- [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [AC1](../software-requirements/srd001-serving-contract.md#ac1)

<!-- sc: S2 -->
### S2

The executor changes the shared workspace and the critic accepts it through two remote boundaries.

<!-- field: traces -->
#### Traces

- [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [R2.1](../software-requirements/srd001-serving-contract.md#r21)
- [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [R2.2](../software-requirements/srd001-serving-contract.md#r22)
- [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [R2.3](../software-requirements/srd001-serving-contract.md#r23)
- [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [R2.4](../software-requirements/srd001-serving-contract.md#r24)
- [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [R3.1](../software-requirements/srd001-serving-contract.md#r31)
- [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [R3.2](../software-requirements/srd001-serving-contract.md#r32)
- [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [R3.3](../software-requirements/srd001-serving-contract.md#r33)
- [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [AC2](../software-requirements/srd001-serving-contract.md#ac2)

<!-- sc: S3 -->
### S3

The three role trace files contain one caller-supplied W3C trace identifier.

<!-- field: traces -->
#### Traces

- [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [R4.3](../software-requirements/srd001-serving-contract.md#r43)
- [srd001-serving-contract](../software-requirements/srd001-serving-contract.md) [AC3](../software-requirements/srd001-serving-contract.md#ac3)

<!-- field: out_of_scope -->
## Out of scope

- Helm resources, profile ConfigMap sizing, and closure sharding.
- Per-request filesystem roots or dynamically supplied transport authority.
- Replacing the local canonical batch integrations.
