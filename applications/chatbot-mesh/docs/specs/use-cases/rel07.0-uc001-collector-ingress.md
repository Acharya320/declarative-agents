---
id: rel07.0-uc001-collector-ingress
title: The collector agent retains and relays an attributable trace batch
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An operator enables agent trace ingress. The collector profile starts its
receiver, control, and monitor surfaces. A chatbot run exports an attributable
batch; the collector awaits it, records stdouttrace-compatible evidence, and
relays the unchanged request to Jaeger. The collector's own lifecycle spans
remain in its file exporter. The migration is now complete: the collector agent
owns OTLP metric intake (GH-1207) and the in-cluster Dolt Prometheus scrape is
retired (GH-1366); contrib remains selectable for the trace path only.

<!-- field: actor -->
## Actor

Mesh integration operator

<!-- field: trigger -->
## Trigger

A values overlay selects collector implementation agent and starts a chatbot scenario with unique run identity.

<!-- field: flow -->
## Flow

- F1: Helm renders the collector-agent profile and points both trace and OTLP metric producers at its OTLP/gRPC Service; no contrib metrics gateway is deployed.
- F2: The collector launches control, monitor, and bounded OTLP receiver listeners and enters await.
- F3: A chatbot scenario exports spans carrying repository, module, target, commit, run id, and stable service identity.
- F4: Await returns one complete batch and declared policy accepts it.
- F5: Spool appends complete stdouttrace-compatible NDJSON lines and syncs them.
- F6: Relay mode exports the unchanged batch to Jaeger and returns to await; test mode returns directly to await.
- F7: Evidence asserts matching trace, span, parent, service, and run identity in the spool and Jaeger.
- F8: Evidence asserts agent OTLP metrics through the collector agent; the in-cluster Dolt Prometheus scrape is retired and no longer required (GH-1366).
- F9: Lifecycle exit stops receiver, monitor, and control while persistent evidence remains queryable.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd007-collector-agent](../software-requirements/srd007-collector-agent.md) [R1.1](../software-requirements/srd007-collector-agent.md#r11), [R1.2](../software-requirements/srd007-collector-agent.md#r12), [R1.3](../software-requirements/srd007-collector-agent.md#r13), [R1.4](../software-requirements/srd007-collector-agent.md#r14) | persistent bounded intake and visible machine sentence |
| T2 | [srd007-collector-agent](../software-requirements/srd007-collector-agent.md) [R2.1](../software-requirements/srd007-collector-agent.md#r21), [R2.2](../software-requirements/srd007-collector-agent.md#r22), [R2.3](../software-requirements/srd007-collector-agent.md#r23), [R2.4](../software-requirements/srd007-collector-agent.md#r24) | complete integration identity in spool and relay |
| T3 | [srd007-collector-agent](../software-requirements/srd007-collector-agent.md) [R3.1](../software-requirements/srd007-collector-agent.md#r31), [R3.2](../software-requirements/srd007-collector-agent.md#r32), [R3.3](../software-requirements/srd007-collector-agent.md#r33), [R3.4](../software-requirements/srd007-collector-agent.md#r34) | test and relay modes |
| T4 | [srd007-collector-agent](../software-requirements/srd007-collector-agent.md) [R4.1](../software-requirements/srd007-collector-agent.md#r41), [R4.2](../software-requirements/srd007-collector-agent.md#r42), [R4.3](../software-requirements/srd007-collector-agent.md#r43) | control, monitor, and lifecycle cleanup |
| T5 | [srd007-collector-agent](../software-requirements/srd007-collector-agent.md) [R5.1](../software-requirements/srd007-collector-agent.md#r51), [R5.2](../software-requirements/srd007-collector-agent.md#r52), [R5.3](../software-requirements/srd007-collector-agent.md#r53), [R5.4](../software-requirements/srd007-collector-agent.md#r54) | file-only telemetry and self-loop proof |
| T6 | [srd007-collector-agent](../software-requirements/srd007-collector-agent.md) [R6.1](../software-requirements/srd007-collector-agent.md#r61), [R6.2](../software-requirements/srd007-collector-agent.md#r62), [R6.3](../software-requirements/srd007-collector-agent.md#r63), [R6.4](../software-requirements/srd007-collector-agent.md#r64), [R6.5](../software-requirements/srd007-collector-agent.md#r65) | agent owns trace and OTLP metric intake; contrib opt-in trace-only |
| T7 | [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [R5.1](../software-requirements/srd003-chatbot-deployment.md#r51), [R5.2](../software-requirements/srd003-chatbot-deployment.md#r52), [R5.3](../software-requirements/srd003-chatbot-deployment.md#r53), [R5.4](../software-requirements/srd003-chatbot-deployment.md#r54), [R5.5](../software-requirements/srd003-chatbot-deployment.md#r55), [R5.6](../software-requirements/srd003-chatbot-deployment.md#r56), [R5.7](../software-requirements/srd003-chatbot-deployment.md#r57), [R5.8](../software-requirements/srd003-chatbot-deployment.md#r58), [R5.9](../software-requirements/srd003-chatbot-deployment.md#r59), [R5.10](../software-requirements/srd003-chatbot-deployment.md#r510), [R5.11](../software-requirements/srd003-chatbot-deployment.md#r511) | retained persistent trace and metric evidence |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The collector profile receives a real chatbot OTLP batch through declared bounded intake.

<!-- field: traces -->
#### Traces

- [srd007-collector-agent](../software-requirements/srd007-collector-agent.md) [AC1](../software-requirements/srd007-collector-agent.md#ac1)

<!-- sc: S2 -->
### S2

Spool and Jaeger contain matching trace, span, parent, service, and unique run identity.

<!-- field: traces -->
#### Traces

- [srd007-collector-agent](../software-requirements/srd007-collector-agent.md) [AC2](../software-requirements/srd007-collector-agent.md#ac2)

<!-- sc: S3 -->
### S3

The same profile passes spool-only test mode and spool-first relay mode.

<!-- field: traces -->
#### Traces

- [srd007-collector-agent](../software-requirements/srd007-collector-agent.md) [AC3](../software-requirements/srd007-collector-agent.md#ac3)

<!-- sc: S4 -->
### S4

Control exit releases listeners and blocked waits while monitor reports bounded outcomes.

<!-- field: traces -->
#### Traces

- [srd007-collector-agent](../software-requirements/srd007-collector-agent.md) [AC4](../software-requirements/srd007-collector-agent.md#ac4)

<!-- sc: S5 -->
### S5

Client spans enter the receiver and collector-process spans remain only in its file.

<!-- field: traces -->
#### Traces

- [srd007-collector-agent](../software-requirements/srd007-collector-agent.md) [AC5](../software-requirements/srd007-collector-agent.md#ac5)

<!-- sc: S6 -->
### S6

Agent OTLP metrics remain available through the collector agent and retained evidence survives teardown; the in-cluster Dolt scrape is retired in agent mode (GH-1366).

<!-- field: traces -->
#### Traces

- [srd007-collector-agent](../software-requirements/srd007-collector-agent.md) [AC6](../software-requirements/srd007-collector-agent.md#ac6)

<!-- field: failure_scenarios -->
## Failure scenarios

| id | scenario | expected |
| --- | --- | --- |
| E1 | Agent mode drops its own OTLP metric intake. | Rejected; the collector agent owns OTLP metric intake (GH-1207) and metric evidence must remain available. The in-cluster Dolt Prometheus scrape, by contrast, is intentionally retired in agent mode (GH-1366). |
| E2 | Collector lifecycle spans arrive at its own receiver. | Rejected; collector telemetry is file-only and the proof observes no self batch. |
| E3 | Relay succeeds before deterministic spool. | Rejected; spool is the preceding visible state and failure cannot route implicitly to relay. |
| E4 | Trace ids match but run or service identity changes. | Rejected; all bounded integration resource identity must match end to end. |

<!-- field: out_of_scope -->
## Out of scope

- Collector-agent logs intake.
- Re-adding an in-cluster Dolt Prometheus scrape to the default agent collector.
- Jaeger query behavior beyond evidence assertions.

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel07.0-collector](../test-suites/test-rel07.0-collector.md)
