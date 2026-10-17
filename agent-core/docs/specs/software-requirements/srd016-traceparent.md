---
id: srd016-traceparent
title: W3C Traceparent Parse and Format
depends_on: []
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Agents that receive a parent trace context from an external caller need
to parse the W3C traceparent header and inject it into the Go context
so that OTel spans are linked to the parent trace. Without a shared
parser, each agent implements its own traceparent handling.

A multi-agent mesh also needs the trace to continue across REST calls between
running agents. Today propagation exists only for subprocess-spawned children
through --otel-parent-span; cross-agent REST calls start a fresh root span, so
one logical request fragments into disconnected traces. This SRD therefore also
governs injecting the active span's trace context into outbound REST client
requests and extracting it on inbound REST server endpoints, reusing the same
parser. HTTP propagation is a uniform transport concern with no per-operation
authoring surface.

<!-- field: goals -->
## Goals

- G1: Parse W3C traceparent version-00 headers into OTel SpanContext.
- G2: Format SpanContext back to a W3C traceparent string.
- G3: Convenience function to create a context with a remote parent span.
- G4: Inject W3C trace context into outbound REST client requests from the active dispatch span.
- G5: Extract W3C trace context on inbound REST server endpoints so a cross-agent trace continues.

<!-- req: R1 -->
## R1 - ParseTraceparent

<!-- req: R1.1 -->
### R1.1

ParseTraceparent(tp string) parses a W3C traceparent header with format "00-<traceID>-<spanID>-<flags>".

<!-- req: R1.2 -->
### R1.2

Returns a remote SpanContext on success, error on malformed input.

<!-- req: R1.3 -->
### R1.3

Rejects non-00 version prefixes.

<!-- req: R1.4 -->
### R1.4

Validates trace-id, span-id, and trace-flags hex encoding.

<!-- req: R2 -->
## R2 - FormatTraceparent

<!-- req: R2.1 -->
### R2.1

FormatTraceparent(sc SpanContext) builds a W3C traceparent string from a SpanContext.

<!-- req: R2.2 -->
### R2.2

Returns empty string for invalid SpanContext.

<!-- req: R3 -->
## R3 - ParseParentSpan

<!-- req: R3.1 -->
### R3.1

ParseParentSpan(raw string) parses a traceparent and returns a context carrying the remote span.

<!-- req: R3.2 -->
### R3.2

Empty input returns context.Background() with no error.

<!-- req: R3.3 -->
### R3.3

Invalid input returns a descriptive error.

<!-- req: R4 -->
## R4 - Outbound REST traceparent injection

<!-- req: R4.1 -->
### R4.1

Outbound REST client requests SHALL inject a traceparent header formatted (R2) from the active dispatch span's SpanContext into the request headers.

<!-- req: R4.2 -->
### R4.2

When the active span carries tracestate, REST client requests SHALL inject a tracestate header alongside traceparent, passed through opaquely.

<!-- req: R4.3 -->
### R4.3

Injection SHALL apply uniformly to every outbound REST client request as a transport concern; it SHALL NOT be author-configurable per operation and SHALL NOT add any rest-tool-format surface.

<!-- req: R4.4 -->
### R4.4

When no recording span is active, REST client requests SHALL omit the traceparent header rather than emit an invalid one.

<!-- req: R5 -->
## R5 - Inbound REST traceparent extraction

<!-- req: R5.1 -->
### R5.1

Inbound REST server endpoints SHALL extract an incoming traceparent header and parent the request span, and any machine_request span, on the extracted remote SpanContext (R1, R3).

<!-- req: R5.2 -->
### R5.2

When present, an incoming tracestate header SHALL be carried onto the extracted context opaquely.

<!-- req: R5.3 -->
### R5.3

An absent or malformed traceparent header SHALL fall back to current behavior, starting a new root span rather than failing the request.

<!-- req: R5.4 -->
### R5.4

Extraction SHALL reuse the ParseTraceparent and ParseParentSpan primitives and SHALL NOT introduce a second parser.

<!-- field: non_goals -->
## Non-goals

- W3C tracestate parsing or validation beyond opaque pass-through of the header value.
- Version negotiation beyond version 00.
- Propagation over transports other than the built-in REST client and server (for example message queues).

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

ParseTraceparent correctly parses valid version-00 traceparent headers and rejects malformed input.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R1.4

<!-- ac: AC2 -->
### AC2

FormatTraceparent round-trips with ParseTraceparent. Invalid SpanContext returns empty string.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2

<!-- ac: AC3 -->
### AC3

ParseParentSpan returns a context with remote span, handles empty input, and wraps errors.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3

<!-- ac: AC4 -->
### AC4

Outbound REST client requests inject traceparent (and tracestate when present) from the active dispatch span uniformly, with no per-operation authoring surface, and omit the header when no recording span is active.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
- R4.3
- R4.4

<!-- ac: AC5 -->
### AC5

Inbound REST server endpoints extract traceparent (and tracestate) to parent the request and machine_request spans, fall back to a new root on absent or malformed headers, and reuse the existing parser.

<!-- field: traces -->
#### Traces

- R5.1
- R5.2
- R5.3
- R5.4

<!-- field: implementation_status | form: "tree" -->
## Implementation status

- done:
  1. ParseTraceparent, FormatTraceparent, and ParseParentSpan primitives, and the subprocess --otel-parent-span propagation path, which remains unchanged and complementary to HTTP propagation.
  1. - "Outbound injection (R4)": the engine injects the active dispatch span into TraceContextAware commands (core.TraceContextAware, set in dispatchWithMonitor from the dispatch span context); the REST client's applyTraceContext writes traceparent (and tracestate when present) into every outbound request uniformly, omitting the header when no recording span is active. No rest-tool-format surface. Delivered by GH-310.
  1. - "Inbound extraction (R5)": the REST server's startMachineRequestSpan extracts an incoming traceparent (and tracestate), parents the machine_request span on the remote SpanContext through the global tracer, and falls back to a new root on absent or malformed headers using the same parser; traceparent/tracestate are allowed undeclared request headers. Delivered by GH-310.

<!-- field: diagrams | form: "tree" -->
## Diagrams

- requirements_trace:
  ```mermaid
  flowchart LR
    Problem["Problem: Agents that receive a parent trace context from an ex..."] --> Goals["Goals: W3C Traceparent Parse and Format"]
    Goals --> R1["R1: ParseTraceparent"]
    Goals --> R2["R2: FormatTraceparent"]
    Goals --> R3["R3: ParseParentSpan"]
    Goals --> R4["R4: Outbound REST traceparent injection"]
    Goals --> R5["R5: Inbound REST traceparent extraction"]
    AC1["AC1: ParseTraceparent correctly parses valid version-00 tracep..."]
    R1 --> AC1
    AC2["AC2: FormatTraceparent round-trips with ParseTraceparent. Inva..."]
    R2 --> AC2
    AC3["AC3: ParseParentSpan returns a context with remote span, handl..."]
    R3 --> AC3
    AC4["AC4: Outbound REST client injects traceparent from the active span"]
    R4 --> AC4
    AC5["AC5: Inbound REST server extracts traceparent to parent the span"]
    R5 --> AC5
  ```
