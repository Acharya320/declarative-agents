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
| T1 | srd001-serving-contract R1.1, R1.2, R1.3 | S1 traces this requirement in the existing success criteria. |
| T2 | srd001-serving-contract AC1 | S1 traces this requirement in the existing success criteria. |
| T3 | srd001-serving-contract R2.1, R2.2, R2.3, R2.4 | S2 traces this requirement in the existing success criteria. |
| T4 | srd001-serving-contract R3.1, R3.2, R3.3 | S2 traces this requirement in the existing success criteria. |
| T5 | srd001-serving-contract AC2 | S2 traces this requirement in the existing success criteria. |
| T6 | srd001-serving-contract R4.3 | S3 traces this requirement in the existing success criteria. |
| T7 | srd001-serving-contract AC3 | S3 traces this requirement in the existing success criteria. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Planner, executor, and critic health endpoints report from live profile-owned listeners.

<!-- field: traces -->
#### Traces

- srd001-serving-contract R1.1
- srd001-serving-contract R1.2
- srd001-serving-contract R1.3
- srd001-serving-contract AC1

<!-- sc: S2 -->
### S2

The executor changes the shared workspace and the critic accepts it through two remote boundaries.

<!-- field: traces -->
#### Traces

- srd001-serving-contract R2.1
- srd001-serving-contract R2.2
- srd001-serving-contract R2.3
- srd001-serving-contract R2.4
- srd001-serving-contract R3.1
- srd001-serving-contract R3.2
- srd001-serving-contract R3.3
- srd001-serving-contract AC2

<!-- sc: S3 -->
### S3

The three role trace files contain one caller-supplied W3C trace identifier.

<!-- field: traces -->
#### Traces

- srd001-serving-contract R4.3
- srd001-serving-contract AC3

<!-- field: out_of_scope -->
## Out of scope

- Helm resources, profile ConfigMap sizing, and closure sharding.
- Per-request filesystem roots or dynamically supplied transport authority.
- Replacing the local canonical batch integrations.
