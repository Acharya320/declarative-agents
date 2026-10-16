---
id: srd007-tracer-interface
title: Tracer Port Interface
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

The agentic loop engine needs to emit OpenTelemetry spans and events
for observability. Depending directly on the OTel SDK would couple
the engine to a specific tracer configuration. A port interface
abstracts span creation, event recording, attribute setting, error
recording, and context access behind five methods. Domain agents
supply the concrete implementation.

<!-- field: goals -->
## Goals

- G1: Minimal interface decoupling the engine from OTel SDK internals.
- G2: Support child span creation, event recording, attribute setting, error recording, and context access.

<!-- req: R1 -->
## R1 - Tracer interface

<!-- req: R1.1 -->
### R1.1

Tracer is an interface in internal/observability/tracing with five methods.

<!-- req: R1.2 -->
### R1.2

Push(name string, attrs ...attribute.KeyValue) returns (Tracer, func()). Creates a child span and returns a scoped Tracer plus a done function.

<!-- req: R1.3 -->
### R1.3

Event(name string, attrs ...attribute.KeyValue) records a point-in-time span event.

<!-- req: R1.4 -->
### R1.4

SetAttributes(attrs ...attribute.KeyValue) sets attributes on the current span.

<!-- req: R1.5 -->
### R1.5

RecordError(err error) records err on the current span and sets its status to error.

<!-- req: R1.6 -->
### R1.6

Context() context.Context returns the context carrying the current span.

<!-- req: R2 -->
## R2 - Dependencies

<!-- req: R2.1 -->
### R2.1

internal/observability/tracing depends on go.opentelemetry.io/otel/attribute for the KeyValue type and context for Context().

<!-- req: R2.2 -->
### R2.2

The interface does not import the OTel SDK (go.opentelemetry.io/otel/sdk).

<!-- req: R3 -->
## R3 - NoopTracer

<!-- req: R3.1 -->
### R3.1

NoopTracer is a concrete implementation that discards all operations.

<!-- req: R3.2 -->
### R3.2

NoopTracer satisfies the Tracer interface at compile time.

<!-- req: R3.3 -->
### R3.3

NoopTracer.Context() returns context.Background().

<!-- field: non_goals -->
## Non-goals

- Concrete Tracer implementation beyond NoopTracer (domain agent responsibility).
- OTel exporter configuration.
- Metrics or logging interfaces.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

Tracer interface has Push, Event, SetAttributes, RecordError, and Context methods.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R1.4
- R1.5
- R1.6

<!-- ac: AC2 -->
### AC2

internal/observability/tracing depends on go.opentelemetry.io/otel/attribute and context. No OTel SDK import.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2

<!-- ac: AC3 -->
### AC3

NoopTracer satisfies Tracer at compile time and discards all operations.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3

<!-- field: diagrams -->
## Diagrams

````yaml
requirements_trace: |
    ```mermaid
    flowchart LR
      Problem["Problem: The agentic loop engine needs to emit OpenTelemetry s..."] --> Goals["Goals: Tracer Port Interface"]
      Goals --> R1["R1: Tracer interface"]
      Goals --> R2["R2: Dependencies"]
      Goals --> R3["R3: NoopTracer"]
      AC1["AC1: Tracer interface has Push, Event, SetAttributes, RecordEr..."]
      R1 --> AC1
      AC2["AC2: internal/observability/tracing depends on go.opentelemetr..."]
      R2 --> AC2
      AC3["AC3: NoopTracer satisfies Tracer at compile time and discards..."]
      R3 --> AC3
    ```
````
