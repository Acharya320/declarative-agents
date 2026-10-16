---
id: rel11.0-uc001-otlp-receiver-await
title: A persistent machine receives, records, and relays OTLP trace and metric batches
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An operator starts a profile whose machine binds an OTLP/gRPC receiver and
waits for trace batches. A client agent exports one batch. The receiver queues
it, the await word returns it to the machine, and visible transitions decide
whether to spool and relay it. The spool uses the same NDJSON shape as
deterministic agent trace files. Relay preserves the request's resource and
span identity for a persistent upstream backend.

The same listener also accepts OTLP metric exports on a separate bounded
queue: `await_metrics` hands a metric batch to the machine, `spool_metrics`
records one NDJSON line per metric, and `spool_list_metrics` and
`spool_get_metric` read that evidence back by name and run identity. The
receiver process exports its own runtime telemetry to a file, not to the
endpoint it owns.

<!-- field: actor -->
## Actor

Operator running a declarative agent as an OTLP trace and metric ingress

<!-- field: trigger -->
## Trigger

The operator starts the ingress profile and a client agent exports an
`ExportTraceServiceRequest` to its configured OTLP/gRPC endpoint.

<!-- field: flow -->
## Flow

- F1: The profile machine dispatches `otlp_receiver_launch` with a receiver name, isolated bind address, queue capacity, and overflow policy (srd042 R1).
- F2: Launch returns `ReceiverLaunched`; the machine enters its persistent await loop.
- F3: A client exports one OTLP trace request. The gRPC handler decodes and enqueues the complete request under bounded backpressure rules (srd042 R2).
- F4: `await_spans` consumes the oldest batch and returns its protobuf payload plus bounded metadata as `SpansReceived` (srd042 R3).
- F5: The machine evaluates declared policy with ordinary vocabulary and routes the accepted batch to `spool_spans`; rejected policy remains a visible transition.
- F6: `spool_spans` resolves the batch from declared command state and appends complete stdouttrace-compatible NDJSON lines (srd042 R4).
- F7: In relay mode, the machine dispatches `relay_spans`; the upstream receives the unchanged trace request and returns success (srd042 R5).
- F8: The machine returns to `await_spans` for the next batch.
- F9: A lifecycle control event routes to `otlp_receiver_stop`, which stops intake, unblocks awaiters, and reaches a success terminal (srd042 R6).
- F10: A client exports one OTLP metric request to the same listener; the MetricsService handler decodes and enqueues the complete request on a separate bounded metric queue under the same backpressure rules (srd042 R9.1, R9.2).
- F11: `await_metrics` consumes the oldest metric batch and returns its payload plus data-point, service, and metric-name metadata as `MetricsReceived` (srd042 R9.2).
- F12: `spool_metrics` appends one NDJSON record per metric carrying resource, scope, and data-point identity under the same bounded rotation as the trace spool (srd042 R9.3).
- F13: `spool_list_metrics` and `spool_get_metric` read the metric spool back by name and run identity without mutation (srd042 R9.4).

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd042-otlp-receiver-tools R0.1, R0.2, R0.3, R0.4, R0.5 | transport words stay separate and MachineSpec owns their sentence |
| T2 | srd042-otlp-receiver-tools R1.1, R1.2, R1.3, R1.4, R1.5, R1.6, R1.7 | declared launch binds OTLP/gRPC and returns without waiting |
| T3 | srd042-otlp-receiver-tools R2.1, R2.2, R2.3, R2.4, R2.5, R2.6 | complete requests enter a bounded concurrent FIFO queue |
| T4 | srd042-otlp-receiver-tools R3.1, R3.2, R3.3, R3.4, R3.5, R3.6, R3.7, R3.8, R3.9, R3.10 | await returns one batch or a distinct timeout, stop, or fault signal |
| T5 | srd042-otlp-receiver-tools R4.1, R4.2, R4.3, R4.4, R4.5, R4.6, R4.7, R4.8 | deterministic spool matches existing stdouttrace NDJSON readers |
| T6 | srd042-otlp-receiver-tools R5.1, R5.2, R5.3, R5.4, R5.5, R5.6, R5.7 | declared relay preserves protobuf identity |
| T7 | srd042-otlp-receiver-tools R6.1, R6.2, R6.3, R6.4, R6.5 | stop owns listener and waiter cleanup |
| T8 | srd042-otlp-receiver-tools R7.1, R7.2, R7.3, R7.4 | file-only process telemetry prevents ingress self-export |
| T9 | srd042-otlp-receiver-tools R8.1, R8.2, R8.3 | the receiver accepts trace and metric exports on one gRPC listener |
| T10 | srd008-telemetry R9.1, R9.2, R9.3, R9.4, R9.5, R9.6, R9.7 | integration evidence keeps stable service and run identity |
| T11 | srd042-otlp-receiver-tools R9.1, R9.2, R9.3, R9.4, R9.5 | metric intake, spool, and query mirror the trace words without a separate backend |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

A real OTLP/gRPC client exports a complete trace batch that the declared await word returns to the machine.

<!-- field: traces -->
#### Traces

- srd042-otlp-receiver-tools AC1

<!-- sc: S2 -->
### S2

The batch is appended in the existing stdouttrace-compatible NDJSON shape with trace, parent, service, and run identity intact.

<!-- field: traces -->
#### Traces

- srd042-otlp-receiver-tools AC2

<!-- sc: S3 -->
### S3

Relay delivers the unchanged batch to an upstream OTLP server, and an upstream rejection reaches the machine as CommandError.

<!-- field: traces -->
#### Traces

- srd042-otlp-receiver-tools AC3

<!-- sc: S4 -->
### S4

Lifecycle stop releases the listener, unblocks waits, and leaves no background goroutine.

<!-- field: traces -->
#### Traces

- srd042-otlp-receiver-tools AC4

<!-- sc: S5 -->
### S5

Client spans arrive while the receiver process's own spans remain only in its deterministic file.

<!-- field: traces -->
#### Traces

- srd042-otlp-receiver-tools AC5

<!-- sc: S6 -->
### S6

The receiver accepts trace and metric exports on one gRPC listener and makes no OTLP logs claim; no separate metrics backend is required.

<!-- field: traces -->
#### Traces

- srd042-otlp-receiver-tools AC6

<!-- sc: S7 -->
### S7

A real OTLP/gRPC client exports a metric batch that await returns, spool records as NDJSON with resource and metric identity, and the query words read back by name, all while trace intake stays unchanged.

<!-- field: traces -->
#### Traces

- srd042-otlp-receiver-tools AC8

<!-- field: failure_scenarios -->
## Failure scenarios

| id | scenario | expected |
| --- | --- | --- |
| E1 | A full queue silently drops a batch. | Rejected; configured overflow behavior reports rejected or dropped span counts and never blocks indefinitely. |
| E2 | A spool write fails and the tool relays the batch anyway. | Rejected; spool emits CommandError and relay can occur only through a later explicit machine transition. |
| E3 | Relay rewrites service.name or a trace identifier. | Rejected; protobuf span and resource fields arrive unchanged at the upstream. |
| E4 | The receiver exports its own lifecycle spans to its listening endpoint. | Rejected; the profile uses file-only telemetry and the self-loop proof observes no collector-process batch. |
| E5 | A metric export reaches the trace queue, or a metric word alters a spooled span. | Rejected; the metric signal uses a separate queue and spool, and trace intake, spool, relay, and query behavior stay unchanged. |

<!-- field: out_of_scope -->
## Out of scope

- OTLP logs.
- General-purpose collector throughput or multi-tenant isolation.
- Jaeger query and trace search behavior.
- TLS and authenticated transport.

<!-- field: test_suite_trace -->
## Test suite trace

test-rel11.0-otlp-receiver
