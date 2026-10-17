---
id: rel08.0-uc001-http-trace-propagation
title: Two agents linked by a REST call produce one connected trace
status: specified
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A caller agent invokes a callee agent over the built-in REST client and server
words. With HTTP trace-context propagation, the caller injects its active
dispatch span's traceparent into the outbound request, the callee extracts it
and parents its request (and machine_request) span on the remote span, and the
two agents' spans join into one connected trace. A request that arrives with an
absent or malformed traceparent starts a new root span rather than failing.

<!-- field: actor -->
## Actor

Operator observing a multi-agent mesh through a trace backend

<!-- field: trigger -->
## Trigger

A caller agent's REST client word calls a callee agent's REST server endpoint
while both export spans to the same trace backend.

<!-- field: flow -->
## Flow

- F1: The caller runs a REST client word inside an active recording dispatch span; the client injects a traceparent header (and tracestate when present) formatted from that span into the outbound request (srd016 R4).
- F2: The callee's REST server endpoint extracts the incoming traceparent and parents its request span, and any machine_request span, on the extracted remote SpanContext (srd016 R5).
- F3: The backend shows one connected trace in which the callee's server span is a child of the caller's client span.
- F4: A separate request arriving with no traceparent, or a malformed one, starts a new root span and completes normally.
- F5: A subprocess-spawned child continues to receive its parent context through --otel-parent-span, unchanged by the HTTP path.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd016-traceparent](../software-requirements/srd016-traceparent.md) [R4.1](../software-requirements/srd016-traceparent.md#r41), [R4.2](../software-requirements/srd016-traceparent.md#r42), [R4.3](../software-requirements/srd016-traceparent.md#r43), [R4.4](../software-requirements/srd016-traceparent.md#r44) | outbound REST client injects traceparent/tracestate from the active span, uniform, no authoring surface |
| T2 | [srd016-traceparent](../software-requirements/srd016-traceparent.md) [R5.1](../software-requirements/srd016-traceparent.md#r51), [R5.2](../software-requirements/srd016-traceparent.md#r52), [R5.3](../software-requirements/srd016-traceparent.md#r53), [R5.4](../software-requirements/srd016-traceparent.md#r54) | inbound REST server extracts traceparent to parent the request and machine_request spans |
| T3 | [srd016-traceparent](../software-requirements/srd016-traceparent.md) [R5.3](../software-requirements/srd016-traceparent.md#r53) | absent or malformed header falls back to a new root without failing |
| T4 | [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [R9.1](../software-requirements/srd028-rest-client-tools.md#r91) | client traces carry the injected trace context |
| T5 | [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [R9.1](../software-requirements/srd029-rest-server-tools.md#r91) | server continues the cross-agent trace on extraction |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

An outbound REST client request carries a traceparent header formatted from the active dispatch span, and tracestate when the span carries it.

<!-- field: traces -->
#### Traces

- [srd016-traceparent](../software-requirements/srd016-traceparent.md) [AC4](../software-requirements/srd016-traceparent.md#ac4)

<!-- sc: S2 -->
### S2

An inbound REST server endpoint parents its request and machine_request span on the extracted remote span, so the caller's client span and the callee's server span form one connected trace.

<!-- field: traces -->
#### Traces

- [srd016-traceparent](../software-requirements/srd016-traceparent.md) [AC5](../software-requirements/srd016-traceparent.md#ac5)

<!-- sc: S3 -->
### S3

A request with an absent or malformed traceparent starts a new root span and completes without error.

<!-- field: traces -->
#### Traces

- [srd016-traceparent](../software-requirements/srd016-traceparent.md) [AC5](../software-requirements/srd016-traceparent.md#ac5)

<!-- sc: S4 -->
### S4

Injection is uniform across outbound requests with no per-operation authoring surface, and no traceparent is emitted when no recording span is active.

<!-- field: traces -->
#### Traces

- [srd016-traceparent](../software-requirements/srd016-traceparent.md) [AC4](../software-requirements/srd016-traceparent.md#ac4)

<!-- field: failure_scenarios -->
## Failure scenarios

| id | scenario | expected |
| --- | --- | --- |
| E1 | The callee receives a well-formed traceparent but starts a new root span anyway. | The extraction regression fails; the server span must be a child of the caller's span. |
| E2 | A malformed traceparent header fails the inbound request. | Rejected as a regression; malformed headers must fall back to a new root (S3). |
| E3 | An operation adds a per-operation traceparent toggle to rest-tool-format. | Rejected; injection is a uniform transport concern with no authoring surface (S4). |

<!-- field: out_of_scope -->
## Out of scope

- tracestate parsing or validation beyond opaque pass-through.
- Propagation over transports other than the built-in REST client and server.
- The subprocess --otel-parent-span path, which is unchanged and complementary.

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel08.0-trace-propagation](../test-suites/test-rel08.0-trace-propagation.md)
