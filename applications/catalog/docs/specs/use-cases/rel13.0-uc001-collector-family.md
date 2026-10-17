---
id: rel13.0-uc001-collector-family
title: A standalone collector receives, spools, and shuts down in spool mode
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

The canonical catalog collector launches in default spool mode, binds its
OTLP receiver and REST servers, receives one trace batch, spools it, skips
relay (mode is spool), and shuts down cleanly when the exit signal arrives.

<!-- field: actor -->
## Actor

Deployment or conformance harness running the catalog collector profile

<!-- field: trigger -->
## Trigger

The harness launches the collector profile with default environment (no
COLLECTOR_MODE override, no COLLECTOR_RELAY_ENDPOINT).

<!-- field: flow -->
## Flow

- F1: The profile machine seeds and launches the OTLP receiver, control REST, and monitor REST.
- F2: The machine enters the AwaitingSpans loop and waits for a trace batch.
- F3: A batch arrives; the machine filters it (span_count > 0), spools the spans to the NDJSON file, and evaluates the relay predicate.
- F4: The relay predicate is Unsatisfied (default mode is spool), so the machine returns to the control-check loop without relaying.
- F5: The harness posts an exit request; the machine exits, stops the receiver, monitor, and control listeners, and reaches Done.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd020-collector](../software-requirements/srd020-collector.md) [R1](../software-requirements/srd020-collector.md#r1---generalized-defaults) | default spool mode, no relay |
| T2 | [srd020-collector](../software-requirements/srd020-collector.md) [R3](../software-requirements/srd020-collector.md#r3---conformance-coverage) | conformance coverage exercises OTLP builtins |
| T3 | [srd020-collector](../software-requirements/srd020-collector.md) [R3](../software-requirements/srd020-collector.md#r3---conformance-coverage) | the conformance coverage that exercises the receiver, await, spool, and stop words, declared as a platform dependency on agent-core srd042-otlp-receiver-tools |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The collector reaches AwaitingSpans and processes at least one batch through filtering and spooling.

<!-- field: traces -->
#### Traces

- [srd020-collector](../software-requirements/srd020-collector.md) [AC1](../software-requirements/srd020-collector.md#ac1)

<!-- sc: S2 -->
### S2

The relay predicate is Unsatisfied in default spool mode.

<!-- field: traces -->
#### Traces

- [srd020-collector](../software-requirements/srd020-collector.md) [AC3](../software-requirements/srd020-collector.md#ac3)

<!-- sc: S3 -->
### S3

The machine reaches terminal state Done after orderly shutdown.

<!-- field: traces -->
#### Traces

- [srd020-collector](../software-requirements/srd020-collector.md) [AC1](../software-requirements/srd020-collector.md#ac1)

<!-- field: failure_scenarios -->
## Failure scenarios

| id | scenario | expected |
| --- | --- | --- |
| E1 | The OTLP receiver fails to bind. | The machine transitions to Failed. |
| E2 | The exit POST arrives while awaiting spans. | The machine processes the exit on the next control check. |

<!-- field: out_of_scope -->
## Out of scope

- Relay mode testing (requires a live upstream endpoint).
- Spool query tool integration (separate REST surface concern).

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel13.0-collector](../test-suites/test-rel13.0-collector.md)
