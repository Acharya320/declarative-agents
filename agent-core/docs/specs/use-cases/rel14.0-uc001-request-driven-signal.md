---
id: rel14.0-uc001-request-driven-signal
title: A model-free profile advances from a validated HTTP request
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An operator starts a profile whose vocabulary has no model word. Trusted HTTP
source configuration validates a request and maps its discriminator to a
declared signal plus structured payload. The runtime admits that envelope
against the current state, executes the ordinary MachineSpec transition, and
returns the admission and run outcome. Undeclared and conflicting requests
are refused without state or checkpoint mutation. A later request resumes a
suspended run through the same checkpoint and loop path.

<!-- field: actor -->
## Actor

An external caller invoking a trusted request source on a model-free agent profile

<!-- field: trigger -->
## Trigger

The profile is serving a configured HTTP source and receives a schema-valid
request whose configured mapping identifies one machine signal and run id.

<!-- field: flow -->
## Flow

- F1: The profile starts without invoke_llm configuration or model-client construction (srd046 R1).
- F2: The HTTP source validates method, route, body, discriminator, run id, and sensitive-field policy before mapping a SignalEnvelope (srd046 R2).
- F3: Admission atomically checks the envelope against the current Position and transition table and acquires single-dispatch ownership for the run (srd046 R3).
- F4: The accepted signal and payload enter the ordinary StateMachine.Step, Registry, Dispatch, and Loop path (srd046 R4).
- F5: The response reports accepted plus the state and run outcome; a later command failure remains distinct from admission (srd046 R4).
- F6: A declared signal that has no transition from the current state returns refused_conflict, while an unknown mapping or signal returns refused_undeclared; neither changes runtime or durable state (srd046 R3).
- F7: A request resumes a suspended run by loading its checkpoint and using the new signal and payload as the first resumed Result, without replay or duplication (srd046 R5).
- F8: A mock-model profile runs through the existing model-driven path as a compatibility regression (srd046 R6).

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd046-request-driven-signal-source](../software-requirements/srd046-request-driven-signal-source.md) [R1.1](../software-requirements/srd046-request-driven-signal-source.md#r11), [R1.2](../software-requirements/srd046-request-driven-signal-source.md#r12), [R1.3](../software-requirements/srd046-request-driven-signal-source.md#r13), [R1.4](../software-requirements/srd046-request-driven-signal-source.md#r14), [R1.5](../software-requirements/srd046-request-driven-signal-source.md#r15) | model-free profile and transport-neutral source envelope |
| T2 | [srd046-request-driven-signal-source](../software-requirements/srd046-request-driven-signal-source.md) [R2.1](../software-requirements/srd046-request-driven-signal-source.md#r21), [R2.2](../software-requirements/srd046-request-driven-signal-source.md#r22), [R2.3](../software-requirements/srd046-request-driven-signal-source.md#r23), [R2.4](../software-requirements/srd046-request-driven-signal-source.md#r24), [R2.5](../software-requirements/srd046-request-driven-signal-source.md#r25) | trusted request mapping and sensitive payload handling |
| T3 | [srd046-request-driven-signal-source](../software-requirements/srd046-request-driven-signal-source.md) [R3.1](../software-requirements/srd046-request-driven-signal-source.md#r31), [R3.2](../software-requirements/srd046-request-driven-signal-source.md#r32), [R3.3](../software-requirements/srd046-request-driven-signal-source.md#r33), [R3.4](../software-requirements/srd046-request-driven-signal-source.md#r34), [R3.5](../software-requirements/srd046-request-driven-signal-source.md#r35), [R3.6](../software-requirements/srd046-request-driven-signal-source.md#r36) | atomic accepted, refused_undeclared, and refused_conflict admission |
| T4 | [srd046-request-driven-signal-source](../software-requirements/srd046-request-driven-signal-source.md) [R4.1](../software-requirements/srd046-request-driven-signal-source.md#r41), [R4.2](../software-requirements/srd046-request-driven-signal-source.md#r42), [R4.3](../software-requirements/srd046-request-driven-signal-source.md#r43), [R4.4](../software-requirements/srd046-request-driven-signal-source.md#r44), [R4.5](../software-requirements/srd046-request-driven-signal-source.md#r45) | ordinary machine execution and response |
| T5 | [srd046-request-driven-signal-source](../software-requirements/srd046-request-driven-signal-source.md) [R5.1](../software-requirements/srd046-request-driven-signal-source.md#r51), [R5.2](../software-requirements/srd046-request-driven-signal-source.md#r52), [R5.3](../software-requirements/srd046-request-driven-signal-source.md#r53), [R5.4](../software-requirements/srd046-request-driven-signal-source.md#r54), [R5.5](../software-requirements/srd046-request-driven-signal-source.md#r55), [R5.6](../software-requirements/srd046-request-driven-signal-source.md#r56) | checkpoint, suspend, and request-driven resume |
| T6 | [srd046-request-driven-signal-source](../software-requirements/srd046-request-driven-signal-source.md) [R6.1](../software-requirements/srd046-request-driven-signal-source.md#r61), [R6.2](../software-requirements/srd046-request-driven-signal-source.md#r62), [R6.3](../software-requirements/srd046-request-driven-signal-source.md#r63), [R6.4](../software-requirements/srd046-request-driven-signal-source.md#r64), [R6.5](../software-requirements/srd046-request-driven-signal-source.md#r65) | source observability and unchanged model path |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

A model-free profile starts without constructing a model client, and one validated HTTP request maps to a transport-neutral SignalEnvelope.

<!-- field: traces -->
#### Traces

- [srd046-request-driven-signal-source](../software-requirements/srd046-request-driven-signal-source.md) [AC1](../software-requirements/srd046-request-driven-signal-source.md#ac1)
- [srd046-request-driven-signal-source](../software-requirements/srd046-request-driven-signal-source.md) [AC2](../software-requirements/srd046-request-driven-signal-source.md#ac2)

<!-- sc: S2 -->
### S2

Accepted input takes the ordinary machine path and returns both admission and resulting run state.

<!-- field: traces -->
#### Traces

- [srd046-request-driven-signal-source](../software-requirements/srd046-request-driven-signal-source.md) [AC3](../software-requirements/srd046-request-driven-signal-source.md#ac3)
- [srd046-request-driven-signal-source](../software-requirements/srd046-request-driven-signal-source.md) [AC4](../software-requirements/srd046-request-driven-signal-source.md#ac4)

<!-- sc: S3 -->
### S3

Undeclared and conflicting requests are distinguished and leave Position, Execution, and the durable checkpoint byte-equivalent.

<!-- field: traces -->
#### Traces

- [srd046-request-driven-signal-source](../software-requirements/srd046-request-driven-signal-source.md) [AC3](../software-requirements/srd046-request-driven-signal-source.md#ac3)

<!-- sc: S4 -->
### S4

A suspended run resumes from its checkpoint with the new accepted signal and payload exactly once.

<!-- field: traces -->
#### Traces

- [srd046-request-driven-signal-source](../software-requirements/srd046-request-driven-signal-source.md) [AC5](../software-requirements/srd046-request-driven-signal-source.md#ac5)

<!-- sc: S5 -->
### S5

Source outcomes are traced without sensitive payloads, request ownership is released, and the mock-model regression remains unchanged.

<!-- field: traces -->
#### Traces

- [srd046-request-driven-signal-source](../software-requirements/srd046-request-driven-signal-source.md) [AC6](../software-requirements/srd046-request-driven-signal-source.md#ac6)

<!-- field: failure_scenarios -->
## Failure scenarios

| id | scenario | expected |
| --- | --- | --- |
| E1 | The caller puts an arbitrary signal name in a request field. | The source refuses it as undeclared unless trusted config maps the field value to a declared signal. |
| E2 | Two requests target the same run while one owns dispatch. | Exactly one is admitted; the other receives refused_conflict and is not queued. |
| E3 | A request claims a suspended run while persistence is NoopCheckpoint. | The request is refused without starting a fresh replacement run. |

<!-- field: out_of_scope -->
## Out of scope

- A downstream task or blackboard profile.
- Timer and queue source adapters beyond their shared SignalEnvelope contract.
- Changes to invoke_llm or model response parsing.

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel14.0-agent-input-and-dolt](../test-suites/test-rel14.0-agent-input-and-dolt.md)
