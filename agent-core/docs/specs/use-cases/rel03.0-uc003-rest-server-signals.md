---
id: rel03.0-uc003-rest-server-signals
title: Grammar receives inbound REST server signals
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An agent profile exposes configured HTTP endpoints for external actors. A
REST server word launches routes without deciding the workflow. HTTP handlers
validate requests and enqueue InboundEvent messages. They stop there. When
MachineSpec later dispatches rest_server_await, the grammar receives
Approved, Rejected, PauseRequested, or a domain-specific signal. One result
can move the workflow. Another can stop it until an operator acts.

<!-- field: actor -->
## Actor

Agent profile author exposing inbound control endpoints

<!-- field: trigger -->
## Trigger

Outside input is required. The agent still owns sequencing through
MachineSpec transitions, even when that input is approval, pause control, or
a domain event from another system.

<!-- field: flow -->
## Flow

- F1: The author defines a rest.servers entry with address, queue policy, endpoint bindings, schemas, and shutdown behavior.
- F2: The author declares rest_server_launch, rest_server_await, and rest_server_stop ToolDefs that reference the server name.
- F3: MachineSpec dispatches rest_server_launch from LaunchControlApi.
- F4: rest_server_launch binds the configured address, registers endpoints, and emits ServerLaunched or CommandError.
- F5: External actors call configured endpoints such as approve, reject, pause, or domain_event.
- F6: Endpoint handlers validate method, path, headers, query, and body before enqueueing an InboundEvent.
- F7: MachineSpec dispatches rest_server_await from AwaitExternalSignal.
- F8: rest_server_await consumes one queued event and emits Approved, Rejected, PauseRequested, DomainEventReceived, AwaitTimedOut, ServerStopped, or CommandError.
- F9: When a profile waits across more than one server queue, MachineSpec dispatches rest_await_event with source, route, signal, timeout, and stopped-source filters.
- F10: MachineSpec routes each signal to a visible next action. Approval continues work, rejection revises or rolls back, pause enters a suspended path, lifecycle exit invokes exit_agent, and domain events invoke domain handling.
- F11: MachineSpec dispatches rest_server_stop before shutdown or when the route set is no longer needed.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [R5.4](../software-requirements/srd027-rest-tool-runtime.md#r54), [R5.5](../software-requirements/srd027-rest-tool-runtime.md#r55), [R5.6](../software-requirements/srd027-rest-tool-runtime.md#r56), [R5.7](../software-requirements/srd027-rest-tool-runtime.md#r57) | launch and await split server activation from inbound event receipt |
| T2 | [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [R2.1](../software-requirements/srd029-rest-server-tools.md#r21), [R2.2](../software-requirements/srd029-rest-server-tools.md#r22), [R2.3](../software-requirements/srd029-rest-server-tools.md#r23), [R2.4](../software-requirements/srd029-rest-server-tools.md#r24), [R2.5](../software-requirements/srd029-rest-server-tools.md#r25), [R2.6](../software-requirements/srd029-rest-server-tools.md#r26), [R2.7](../software-requirements/srd029-rest-server-tools.md#r27) | rest_server_launch starts configured routes and declares effects |
| T3 | [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [R3.1](../software-requirements/srd029-rest-server-tools.md#r31), [R3.2](../software-requirements/srd029-rest-server-tools.md#r32), [R3.3](../software-requirements/srd029-rest-server-tools.md#r33), [R3.4](../software-requirements/srd029-rest-server-tools.md#r34), [R3.5](../software-requirements/srd029-rest-server-tools.md#r35), [R3.6](../software-requirements/srd029-rest-server-tools.md#r36), [R3.7](../software-requirements/srd029-rest-server-tools.md#r37) | rest_server_await consumes queued inbound events |
| T4 | [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [R4.1](../software-requirements/srd029-rest-server-tools.md#r41), [R4.2](../software-requirements/srd029-rest-server-tools.md#r42), [R4.3](../software-requirements/srd029-rest-server-tools.md#r43), [R4.4](../software-requirements/srd029-rest-server-tools.md#r44), [R4.5](../software-requirements/srd029-rest-server-tools.md#r45), [R4.6](../software-requirements/srd029-rest-server-tools.md#r46), [R4.7](../software-requirements/srd029-rest-server-tools.md#r47) | rest_server_stop shuts down shared server state |
| T5 | [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [R5.1](../software-requirements/srd029-rest-server-tools.md#r51), [R5.2](../software-requirements/srd029-rest-server-tools.md#r52), [R5.3](../software-requirements/srd029-rest-server-tools.md#r53), [R5.4](../software-requirements/srd029-rest-server-tools.md#r54), [R5.5](../software-requirements/srd029-rest-server-tools.md#r55), [R5.6](../software-requirements/srd029-rest-server-tools.md#r56), [R5.7](../software-requirements/srd029-rest-server-tools.md#r57) | endpoint binding kinds cover signals, state reads, handlers, streams, health, and metadata |
| T6 | [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [R6.1](../software-requirements/srd029-rest-server-tools.md#r61), [R6.2](../software-requirements/srd029-rest-server-tools.md#r62), [R6.3](../software-requirements/srd029-rest-server-tools.md#r63), [R6.6](../software-requirements/srd029-rest-server-tools.md#r66) | validation precedes event enqueue and handlers do not choose MachineSpec actions |
| T7 | cfg-rest-tool-format queue_config and endpoint fields | queue limits, payload shape, timeout, binding, signal, and allowed_signals are configured |
| T8 | [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [R10.1](../software-requirements/srd027-rest-tool-runtime.md#r101), [R10.2](../software-requirements/srd027-rest-tool-runtime.md#r102), [R10.3](../software-requirements/srd027-rest-tool-runtime.md#r103), [R10.4](../software-requirements/srd027-rest-tool-runtime.md#r104), [R10.5](../software-requirements/srd027-rest-tool-runtime.md#r105), [R10.6](../software-requirements/srd027-rest-tool-runtime.md#r106), [R10.7](../software-requirements/srd027-rest-tool-runtime.md#r107), [R10.8](../software-requirements/srd027-rest-tool-runtime.md#r108) | rest_await_event fans in across configured server queues |
| T9 | [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [R11.1](../software-requirements/srd029-rest-server-tools.md#r111), [R11.2](../software-requirements/srd029-rest-server-tools.md#r112), [R11.3](../software-requirements/srd029-rest-server-tools.md#r113), [R11.4](../software-requirements/srd029-rest-server-tools.md#r114), [R11.5](../software-requirements/srd029-rest-server-tools.md#r115), [R11.6](../software-requirements/srd029-rest-server-tools.md#r116), [R11.7](../software-requirements/srd029-rest-server-tools.md#r117), [R11.8](../software-requirements/srd029-rest-server-tools.md#r118) | lifecycle routes stay visible through configured server event sources |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Server launch is a non-blocking MachineSpec action that emits ServerLaunched or CommandError.

<!-- field: traces -->
#### Traces

- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC2](../software-requirements/srd029-rest-server-tools.md#ac2)

<!-- sc: S2 -->
### S2

Inbound HTTP requests become grammar-visible signals only through rest_server_await.

<!-- field: traces -->
#### Traces

- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC3](../software-requirements/srd029-rest-server-tools.md#ac3)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC6](../software-requirements/srd029-rest-server-tools.md#ac6)

<!-- sc: S3 -->
### S3

Approved, Rejected, PauseRequested, and DomainEventReceived each have visible machine routes.

<!-- field: traces -->
#### Traces

- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC5](../software-requirements/srd027-rest-tool-runtime.md#ac5)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC5](../software-requirements/srd029-rest-server-tools.md#ac5)

<!-- sc: S4 -->
### S4

Server shutdown drains or fails queued events according to configured policy.

<!-- field: traces -->
#### Traces

- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC4](../software-requirements/srd029-rest-server-tools.md#ac4)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC7](../software-requirements/srd029-rest-server-tools.md#ac7)

<!-- sc: S5 -->
### S5

rest_await_event selects one matching event across configured sources and reports timeout or stopped-source outcomes by policy.

<!-- field: traces -->
#### Traces

- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC10](../software-requirements/srd027-rest-tool-runtime.md#ac10)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC10](../software-requirements/srd029-rest-server-tools.md#ac10)

<!-- sc: S6 -->
### S6

Lifecycle exit requests enqueue ExitRequested and route through exit_agent instead of HTTP handler code.

<!-- field: traces -->
#### Traces

- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC10](../software-requirements/srd027-rest-tool-runtime.md#ac10)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC10](../software-requirements/srd029-rest-server-tools.md#ac10)

<!-- field: out_of_scope -->
## Out of scope

- Defining a general web framework.
- Allowing HTTP handlers to run hidden workflows.
- Replacing bench UI behavior.
- Creating the release 03.0 test suite.

<!-- field: semantic_model -->
## Semantic model

sm-rest-tool-language

<!-- field: config_format -->
## Config format

cfg-rest-tool-format
