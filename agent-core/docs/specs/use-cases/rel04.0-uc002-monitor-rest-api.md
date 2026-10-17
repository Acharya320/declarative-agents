---
id: rel04.0-uc002-monitor-rest-api
title: Supervisor reads monitor state through REST
status: done
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

Supervisors sometimes need a narrow read surface while an agent is still
running. The monitor REST API provides that surface from configured
endpoints. Agent profiles own the server address and route paths. They own
authentication references too. Limits and OpenAPI exposure stay in the same
trusted profile layer, which keeps HTTP callers from expanding monitor
authority at request time.

Handlers read cached monitor snapshots only. Responses cover machine and tool
views plus metric summaries from the store. Recent events and streams come
from the same cache. No handler executes tools. No handler hides a
MachineSpec workflow step behind HTTP. Because a supervisor may reconnect
after a long run and ask for one endpoint before opening a stream, the same
cached snapshot must answer every view without asking a tool to refresh
itself. Read-only behavior is the point.

Monitor serving keeps the process alive by using the same shape as the
Knowledge Manager documentation curator. After launch, the REST listener runs
in the background. A blocking await word keeps the MachineSpec in a serving or
awaiting state. Later, a configured control exit request stops the owned
listener and routes the machine to Done or Succeeded. When the profile uses
127.0.0.1:0, the supervisor discovers the bound listener address from the
REST server launch output before it requests monitor routes.

Monitor OpenAPI covers read endpoints and profile-owned lifecycle control
routes such as POST /monitor/control/exit. JSON responses use stable
lower_snake_case API fields. Metric samples include the same runtime envelope
for dispatch and tool-owned specialized metrics: run_id, state, signal,
status, tool_name, and trusted workflow labels where present.
The OpenAPI schema uses precise JSON Schema types for these fields. It does
not hide scalar fields behind generic object schemas. The event stream route
is documented as text/event-stream. Stream frames use run_event and
metric_sample event names with JSON data payloads that match cached monitor
event and sample shapes.

A separate declared_machines read view exposes the profile's root
MachineSpec and every distinct machine referenced by its machine_request
endpoints. The collection is declaration-only: it preserves authored
view_tags and per-state tags for renderers but reports no live state for
request-scoped runs.

<!-- field: actor -->
## Actor

Supervisor reading a running agent

<!-- field: trigger -->
## Trigger

A running agent has monitor config enabled. A supervisor needs to inspect
current progress and metric state through the profile-owned REST API.

<!-- field: flow -->
## Flow

- F1: The profile declares monitor REST routes and OpenAPI exposure through trusted REST config.
- F2: Runtime starts the configured REST server using the REST server tool boundary and returns the bound listener address in launch output.
- F3: The monitor machine routes ServerLaunched to a non-terminal serving or awaiting state.
- F4: The machine waits on a configured control event while the REST listener serves monitor reads.
- F5: If the configured address uses 127.0.0.1:0, the supervisor discovers the actual base URL from launch output.
- F6: The supervisor requests the machine spec endpoint.
- F7: The supervisor requests current machine state and selected tool metadata.
- F8: The supervisor requests per-tool metrics and aggregate metrics.
- F9: The supervisor requests recent events or opens the configured event stream.
- F10: The supervisor receives SSE frames named run_event or metric_sample with JSON data payloads.
- F11: The supervisor requests OpenAPI and sees monitor read routes plus the configured control exit route.
- F12: The OpenAPI document describes snapshot, event, stream, and control routes with type-accurate schemas and text/event-stream content for the stream route.
- F13: Handlers return lower_snake_case fields and runtime envelopes on metric samples.
- F14: Handlers read cached monitor snapshots and apply redaction rules.
- F15: A configured control exit request stops the owned listener and routes the machine to Done or Succeeded.
- F16: Handler failure returns a structured REST error and does not change command or machine state.
- F17: The supervisor requests a profile-owned GET /monitor/machines endpoint whose read_state binding selects monitor_view declared_machines.
- F18: The supervisor receives the root machine first, followed by each distinct machine_request MachineSpec ordered by name, with authored presentation tags and no request-machine run state.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [R1.1](../software-requirements/srd033-monitor-rest-api.md#r11), [R1.2](../software-requirements/srd033-monitor-rest-api.md#r12), [R1.3](../software-requirements/srd033-monitor-rest-api.md#r13), [R1.4](../software-requirements/srd033-monitor-rest-api.md#r14), [R1.5](../software-requirements/srd033-monitor-rest-api.md#r15) | read-only API scope and profile-owned REST authority |
| T2 | [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [R2.1](../software-requirements/srd033-monitor-rest-api.md#r21), [R2.2](../software-requirements/srd033-monitor-rest-api.md#r22), [R2.3](../software-requirements/srd033-monitor-rest-api.md#r23), [R2.4](../software-requirements/srd033-monitor-rest-api.md#r24), [R2.5](../software-requirements/srd033-monitor-rest-api.md#r25), [R2.6](../software-requirements/srd033-monitor-rest-api.md#r26), [R2.7](../software-requirements/srd033-monitor-rest-api.md#r27), [R2.8](../software-requirements/srd033-monitor-rest-api.md#r28), [R2.9](../software-requirements/srd033-monitor-rest-api.md#r29), [R2.10](../software-requirements/srd033-monitor-rest-api.md#r210), [R2.11](../software-requirements/srd033-monitor-rest-api.md#r211), [R2.12](../software-requirements/srd033-monitor-rest-api.md#r212) | OpenAPI endpoint, control-route, redaction, precise schema documentation, and SSE content type |
| T3 | [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [R3.1](../software-requirements/srd033-monitor-rest-api.md#r31), [R3.2](../software-requirements/srd033-monitor-rest-api.md#r32), [R3.3](../software-requirements/srd033-monitor-rest-api.md#r33), [R3.4](../software-requirements/srd033-monitor-rest-api.md#r34), [R3.5](../software-requirements/srd033-monitor-rest-api.md#r35), [R3.6](../software-requirements/srd033-monitor-rest-api.md#r36), [R3.7](../software-requirements/srd033-monitor-rest-api.md#r37), [R3.8](../software-requirements/srd033-monitor-rest-api.md#r38), [R3.9](../software-requirements/srd033-monitor-rest-api.md#r39), [R3.10](../software-requirements/srd033-monitor-rest-api.md#r310), [R3.11](../software-requirements/srd033-monitor-rest-api.md#r311), [R3.12](../software-requirements/srd033-monitor-rest-api.md#r312) | monitor snapshot endpoints, response field shape, OpenAPI schema shape, and metric sample envelope |
| T4 | [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [R4.1](../software-requirements/srd033-monitor-rest-api.md#r41), [R4.2](../software-requirements/srd033-monitor-rest-api.md#r42), [R4.3](../software-requirements/srd033-monitor-rest-api.md#r43), [R4.4](../software-requirements/srd033-monitor-rest-api.md#r44), [R4.5](../software-requirements/srd033-monitor-rest-api.md#r45), [R4.6](../software-requirements/srd033-monitor-rest-api.md#r46), [R4.7](../software-requirements/srd033-monitor-rest-api.md#r47) | event stream reads, SSE event names, and frame payload shape |
| T5 | [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [R5.1](../software-requirements/srd033-monitor-rest-api.md#r51), [R5.2](../software-requirements/srd033-monitor-rest-api.md#r52), [R5.3](../software-requirements/srd033-monitor-rest-api.md#r53), [R5.4](../software-requirements/srd033-monitor-rest-api.md#r54), [R5.5](../software-requirements/srd033-monitor-rest-api.md#r55) | redaction and failure behavior |
| T6 | [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [R5.3](../software-requirements/srd029-rest-server-tools.md#r53), [R5.5](../software-requirements/srd029-rest-server-tools.md#r55), [R5.6](../software-requirements/srd029-rest-server-tools.md#r56), [R5.7](../software-requirements/srd029-rest-server-tools.md#r57) | read_state, stream_events, health, and static_metadata endpoint kinds |
| T7 | [srd032-tool-metrics-monitor-store](../software-requirements/srd032-tool-metrics-monitor-store.md) [R1.4](../software-requirements/srd032-tool-metrics-monitor-store.md#r14), [R2.6](../software-requirements/srd032-tool-metrics-monitor-store.md#r26), [R4.4](../software-requirements/srd032-tool-metrics-monitor-store.md#r44) | cached snapshots, recorder boundaries, and export diagnostics |
| T8 | cfg-rest-tool-format rest_server | profile-owned server routes, limits, auth, and OpenAPI binding |
| T9 | [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [R6.1](../software-requirements/srd033-monitor-rest-api.md#r61), [R6.2](../software-requirements/srd033-monitor-rest-api.md#r62), [R6.3](../software-requirements/srd033-monitor-rest-api.md#r63), [R6.4](../software-requirements/srd033-monitor-rest-api.md#r64), [R6.5](../software-requirements/srd033-monitor-rest-api.md#r65), [R6.6](../software-requirements/srd033-monitor-rest-api.md#r66), [R6.7](../software-requirements/srd033-monitor-rest-api.md#r67), [R6.8](../software-requirements/srd033-monitor-rest-api.md#r68), [R6.9](../software-requirements/srd033-monitor-rest-api.md#r69) | serving lifecycle, isolated binding, bound-address discovery, and real profile proof |
| T10 | [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [R8.1](../software-requirements/srd033-monitor-rest-api.md#r81), [R8.2](../software-requirements/srd033-monitor-rest-api.md#r82), [R8.3](../software-requirements/srd033-monitor-rest-api.md#r83), [R8.4](../software-requirements/srd033-monitor-rest-api.md#r84), [R8.5](../software-requirements/srd033-monitor-rest-api.md#r85), [R8.6](../software-requirements/srd033-monitor-rest-api.md#r86), [R8.7](../software-requirements/srd033-monitor-rest-api.md#r87), [R8.8](../software-requirements/srd033-monitor-rest-api.md#r88), [R8.9](../software-requirements/srd033-monitor-rest-api.md#r89), [R8.10](../software-requirements/srd033-monitor-rest-api.md#r810) | profile declaration collection, presentation metadata, stable ordering, OpenAPI, and no request-machine run state |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Monitor REST routes read cached monitor state and do not execute tools or MachineSpec actions.

<!-- field: traces -->
#### Traces

- [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [AC1](../software-requirements/srd033-monitor-rest-api.md#ac1)

<!-- sc: S2 -->
### S2

OpenAPI output documents every configured monitor read endpoint, the SSE stream content type, and profile-owned control endpoint without exposing secret or prompt data, and it uses precise JSON Schema types for runtime payload fields.

<!-- field: traces -->
#### Traces

- [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [AC2](../software-requirements/srd033-monitor-rest-api.md#ac2)

<!-- sc: S3 -->
### S3

Snapshot endpoints expose machine, tool, metric, aggregate, and recent-event views from the monitor store with lower_snake_case fields, matching OpenAPI schemas, and runtime envelopes on metric samples.

<!-- field: traces -->
#### Traces

- [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [AC3](../software-requirements/srd033-monitor-rest-api.md#ac3)

<!-- sc: S4 -->
### S4

Event streams emit cached monitor updates as named SSE frames with JSON data payloads and without creating workflow signals or keeping tools running.

<!-- field: traces -->
#### Traces

- [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [AC4](../software-requirements/srd033-monitor-rest-api.md#ac4)

<!-- sc: S5 -->
### S5

REST failures and redaction failures do not change command results, monitor store state, or MachineSpec state.

<!-- field: traces -->
#### Traces

- [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [AC5](../software-requirements/srd033-monitor-rest-api.md#ac5)

<!-- sc: S6 -->
### S6

The real cmd/agent runtime-state-reader profile stays running while the monitor REST API is queried, discovers the bound listener address from profile-owned launch output, uses isolated listener binding for automated proof, and exits successfully after a configured control request.

<!-- field: traces -->
#### Traces

- [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [AC6](../software-requirements/srd033-monitor-rest-api.md#ac6)

<!-- sc: S7 -->
### S7

A profile-owned declared_machines endpoint returns the root first and each distinct machine_request declaration by name, preserves authored view_tags and per-state tags in a typed response, leaves machine_spec unchanged, and reports no request-machine run state.

<!-- field: traces -->
#### Traces

- [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [AC8](../software-requirements/srd033-monitor-rest-api.md#ac8)

<!-- field: out_of_scope -->
## Out of scope

- Writing monitor state through REST.
- Dispatching tools from monitor HTTP handlers.
- Defining durable metric history.
- Reporting live, recent, or historical request-machine runs through declared_machines; spans remain the run-observation surface.

<!-- field: semantic_model -->
## Semantic model

sm-tool-metrics-monitoring

<!-- field: config_format -->
## Config format

cfg-rest-tool-format
