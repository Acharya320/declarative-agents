---
id: rel21.0-uc001-capability-profiles
title: One capability profile is consumed by more than one form
status: specified
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A profile author has behaviour that one agent reaches through a route and a
second agent now wants. Today the route names the agent's own profile and
overrides its machine, so the behaviour cannot be reached from anywhere else
and the second agent copies it. The author instead writes the behaviour as a
capability profile that terminates and hosts nothing, points the original
route at it by path, and invokes the same profile from the second agent as a
child process. Neither consumption edits the capability. The consistency
check reports the machines that more than one profile still reaches, so the
remaining copies are visible rather than assumed.

<!-- field: actor -->
## Actor

A profile author reusing behaviour that a machine override had trapped inside one agent

<!-- field: trigger -->
## Trigger

A second profile needs behaviour that an existing `machine_request` binding
reaches through a `machine` override.

<!-- field: flow -->
## Flow

- F1: The author extracts the overridden machine and its declarations into a capability profile that reaches a terminal state and declares no serve lifecycle, no control or monitor server, and no static assets.
- F2: The original route's binding drops its machine override and names the capability by profile path, mapping the capability's terminal states to its responses.
- F3: The second agent invokes the same capability profile through self_invoke, with selectors carrying the payload, trace context propagating into the child, and the workspace undo contract applying.
- F4: A serving wrapper hosts the capability as a workload where one is wanted, as a blueprint instance where the wrappers are one agent with arguments, without editing the capability.
- F5: The consistency check reports each machine file that more than one profile reaches and each override whose machine another profile reaches, naming the machine and every reaching profile, as findings rather than load errors.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd057-capability-profiles](../software-requirements/srd057-capability-profiles.md) [R1.1](../software-requirements/srd057-capability-profiles.md#r11), [R1.2](../software-requirements/srd057-capability-profiles.md#r12), [R1.3](../software-requirements/srd057-capability-profiles.md#r13) | the capability profile |
| T2 | [srd057-capability-profiles](../software-requirements/srd057-capability-profiles.md) [R2.1](../software-requirements/srd057-capability-profiles.md#r21), [R2.2](../software-requirements/srd057-capability-profiles.md#r22), [R2.3](../software-requirements/srd057-capability-profiles.md#r23), [R2.4](../software-requirements/srd057-capability-profiles.md#r24) | consumption forms |
| T3 | [srd057-capability-profiles](../software-requirements/srd057-capability-profiles.md) [R3.1](../software-requirements/srd057-capability-profiles.md#r31), [R3.2](../software-requirements/srd057-capability-profiles.md#r32), [R3.3](../software-requirements/srd057-capability-profiles.md#r33) | binding discipline and visibility |
| T4 | [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [R1.1](../software-requirements/srd055-agent-blueprints.md#r11), [R1.2](../software-requirements/srd055-agent-blueprints.md#r12), [R1.3](../software-requirements/srd055-agent-blueprints.md#r13) | the serving wrapper as a blueprint instance |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

A capability profile runs to a terminal state carrying no serve lifecycle, no control or monitor server, no static assets, and no field a hand-written profile lacks, and is reached only by a declared path.

<!-- field: traces -->
#### Traces

- [srd057-capability-profiles](../software-requirements/srd057-capability-profiles.md) [AC1](../software-requirements/srd057-capability-profiles.md#ac1)

<!-- sc: S2 -->
### S2

The same unchanged capability profile is bound by a machine_request naming profile, invoked by self_invoke with selectors, undo, and trace propagation, and hosted by a serving wrapper; run_point rejects a profile that declares no point machine.

<!-- field: traces -->
#### Traces

- [srd057-capability-profiles](../software-requirements/srd057-capability-profiles.md) [AC2](../software-requirements/srd057-capability-profiles.md#ac2)

<!-- sc: S3 -->
### S3

The consistency check names a machine reached by two profiles and an override whose machine another profile reaches, reports neither as a load error, and stays silent for a machine reached by one profile.

<!-- field: traces -->
#### Traces

- [srd057-capability-profiles](../software-requirements/srd057-capability-profiles.md) [AC3](../software-requirements/srd057-capability-profiles.md#ac3)

<!-- field: out_of_scope -->
## Out of scope

- Converting the existing machine-override bindings; the use case covers one extraction, not the migration.
- A general nested-machine word, and any change to run_point.
- Behaviour that genuinely differs between applications.
- Any change to the loader, the profile schema, staging, or program identity.
