---
id: rel03.0-uc008-rest-await-any-fan-in
title: Staged REST await-any fan-in adds event sources
status: done
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

Profiles can launch REST servers in stages and still keep event fan-in
visible in MachineSpec. First, the machine launches two REST server sources.
Next, rest_await_event waits across those queues and returns one signal. A
normal transition handles that signal and launches a third REST server only
after the first inbound event has changed machine state. Later, another
rest_await_event dispatch waits across all three launched sources.
This staged shape matters when a profile must avoid exposing a later route
before earlier control or domain input has arrived.

rest_await_event only selects one configured event. Hidden server launches,
workflow choices, and machine mutation stay outside the fan-in word.
MachineSpec owns the staged launch order and each post-event transition.

<!-- field: actor -->
## Actor

Agent profile author composing staged inbound REST workflows

<!-- field: trigger -->
## Trigger

A staged workflow needs first-event behavior before it exposes another REST
source to the same await-any word.

<!-- field: flow -->
## Flow

- F1: The author declares three REST server definitions with separate event queues and configured signal-producing endpoints.
- F2: MachineSpec dispatches rest_server_launch for the first server and routes ServerLaunched to the second launch step.
- F3: MachineSpec dispatches rest_server_launch for the second server and routes ServerLaunched to the first await-any step.
- F4: The first rest_await_event word uses source filters for the first two launched servers.
- F5: An external actor sends an accepted request to either launched server.
- F6: rest_await_event returns exactly one selected event with source, route, method, signal, payload, and request metadata.
- F7: MachineSpec routes the selected signal to a visible transition that launches the third REST server.
- F8: After the third server emits ServerLaunched, MachineSpec dispatches a second rest_await_event word.
- F9: The second rest_await_event word includes sources for all three launched servers and returns the first matching event.
- F10: MachineSpec routes the second selected signal to the next visible workflow action or shutdown path.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [R2.2](../software-requirements/srd027-rest-tool-runtime.md#r22), [R5.4](../software-requirements/srd027-rest-tool-runtime.md#r54), [R5.5](../software-requirements/srd027-rest-tool-runtime.md#r55), [R5.6](../software-requirements/srd027-rest-tool-runtime.md#r56), [R5.7](../software-requirements/srd027-rest-tool-runtime.md#r57) | MachineSpec owns staged launch and await sequencing |
| T2 | [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [R10.1](../software-requirements/srd027-rest-tool-runtime.md#r101), [R10.2](../software-requirements/srd027-rest-tool-runtime.md#r102), [R10.3](../software-requirements/srd027-rest-tool-runtime.md#r103), [R10.4](../software-requirements/srd027-rest-tool-runtime.md#r104), [R10.5](../software-requirements/srd027-rest-tool-runtime.md#r105) | rest_await_event waits across configured sources and returns one event Result |
| T3 | [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [R2.3](../software-requirements/srd029-rest-server-tools.md#r23), [R2.4](../software-requirements/srd029-rest-server-tools.md#r24), [R2.5](../software-requirements/srd029-rest-server-tools.md#r25) | rest_server_launch registers routes without blocking for inbound requests |
| T4 | [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [R11.1](../software-requirements/srd029-rest-server-tools.md#r111), [R11.2](../software-requirements/srd029-rest-server-tools.md#r112), [R11.3](../software-requirements/srd029-rest-server-tools.md#r113) | launched server queues expose named sources for rest_await_event |
| T5 | cfg-rest-tool-format rest_await_event_config and event_source | source filters, timeout, read policy, and stopped-source behavior are configured |
| T6 | sm-rest-tool-language RestEventSource | fan-in sources are declarative and reference REST server queues |
| T7 | [rel03.0-uc003-rest-server-signals](rel03.0-uc003-rest-server-signals.md) [S5](rel03.0-uc003-rest-server-signals.md#s5) | generic REST await selects one matching event across configured sources |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Initial REST server launches are visible MachineSpec actions and do not block on inbound requests.

<!-- field: traces -->
#### Traces

- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC2](../software-requirements/srd029-rest-server-tools.md#ac2)

<!-- sc: S2 -->
### S2

First rest_await_event dispatch waits across only the first two launched sources and returns one selected event.

<!-- field: traces -->
#### Traces

- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC10](../software-requirements/srd027-rest-tool-runtime.md#ac10)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC10](../software-requirements/srd029-rest-server-tools.md#ac10)

<!-- sc: S3 -->
### S3

Later REST server launch happens only after MachineSpec routes the first selected event signal.

<!-- field: traces -->
#### Traces

- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC5](../software-requirements/srd027-rest-tool-runtime.md#ac5)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC7](../software-requirements/srd029-rest-server-tools.md#ac7)

<!-- sc: S4 -->
### S4

Second rest_await_event dispatch includes all three launched sources and still returns one selected event.

<!-- field: traces -->
#### Traces

- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC10](../software-requirements/srd027-rest-tool-runtime.md#ac10)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC10](../software-requirements/srd029-rest-server-tools.md#ac10)

<!-- sc: S5 -->
### S5

Source, route, signal, timeout, and stopped-source policy remain configured fan-in inputs rather than runtime HTTP caller authority.

<!-- field: traces -->
#### Traces

- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC2](../software-requirements/srd027-rest-tool-runtime.md#ac2)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC6](../software-requirements/srd027-rest-tool-runtime.md#ac6)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC10](../software-requirements/srd029-rest-server-tools.md#ac10)

<!-- field: out_of_scope -->
## Out of scope

- Defining a fairness algorithm beyond the configured read policy.
- Letting HTTP handlers launch the third server directly.
- Adding hidden workflow actions inside rest_await_event.
- Replacing single-server rest_server_await coverage.

<!-- field: semantic_model -->
## Semantic model

sm-rest-tool-language

<!-- field: config_format -->
## Config format

cfg-rest-tool-format

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel03.0-rest-tools](../test-suites/test-rel03.0-rest-tools.md)
