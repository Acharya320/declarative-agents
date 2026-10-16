---
id: srd003-command-dispatch
title: Command Dispatch
depends_on:
  - srd_id: srd001-core-types
    symbols_used: [Command, Result, Cost, Signal, CommandError, SpanOverride]
  - srd_id: srd007-tracer-interface
    symbols_used: [Tracer]
  - srd_id: srd015-genai-semconv
    symbols_used: [ToolSpanName, ToolAttrs, ErrorAttrs]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Every command execution needs the same cross-cutting concerns:
tracing, timeout enforcement, panic recovery, duration measurement,
and signal normalization. Without a shared dispatch function, each
domain agent implements these independently and inconsistently.

<!-- field: goals -->
## Goals

- G1: Wrap Command.Execute with tracing, timeout, and panic recovery.
- G2: Normalize signals and durations on all command outcomes.
- G3: Export SafeExecute, FillDuration, and ForceErrorSignal as reusable building blocks.

<!-- req: R1 -->
## R1 - Dispatch function

<!-- req: R1.1 -->
### R1.1

Dispatch(cmd Command, tr Tracer, timeout time.Duration) returns Result.

<!-- req: R1.2 -->
### R1.2

Dispatch creates a child span using GenAI semconv naming: execute_tool {name} by default. Commands implementing SpanOverride customize the span name and creation attributes.

<!-- req: R1.3 -->
### R1.3

The child span records command.name, command.signal, command.duration_ms, gen_ai.usage.input_tokens, and gen_ai.usage.output_tokens as attributes.

<!-- req: R1.4 -->
### R1.4

If Result.Err is non-nil, Dispatch sets error.type attributes and calls RecordError on the span.

<!-- req: R1.5 -->
### R1.5

When an optional monitor recorder is configured, Dispatch SHALL emit standard dispatch metric samples without changing the tracing span contract. Runtime-owned fields include command name, signal, duration, status, error class, and token counts already present in Result.Cost.

<!-- req: R2 -->
## R2 - SafeExecute function

<!-- req: R2.1 -->
### R2.1

SafeExecute(cmd Command, timeout time.Duration) returns Result.

<!-- req: R2.2 -->
### R2.2

SafeExecute runs cmd.Execute() in a goroutine and recovers panics.

<!-- req: R2.3 -->
### R2.3

Panics produce a Result with CommandError signal, a formatted error, and panic output.

<!-- req: R2.4 -->
### R2.4

When timeout is positive and exceeded, SafeExecute returns a timeout Result with CommandError signal.

<!-- req: R2.5 -->
### R2.5

When timeout is zero or negative, SafeExecute waits indefinitely.

<!-- req: R2.6 -->
### R2.6

After normal completion, SafeExecute calls FillDuration and ForceErrorSignal.

<!-- req: R3 -->
## R3 - Helper functions

<!-- req: R3.1 -->
### R3.1

FillDuration(res *Result, start time.Time) sets Cost.Duration from wall clock if the command did not set it.

<!-- req: R3.2 -->
### R3.2

ForceErrorSignal(res *Result) overrides Signal to CommandError when Err is non-nil.

<!-- field: non_goals -->
## Non-goals

- Retry logic (domain agent responsibility).
- Command queuing or batching.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

Dispatch creates a child span, calls SafeExecute, sets CommandName, stamps span attributes.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R1.4
- R1.5

<!-- ac: AC2 -->
### AC2

SafeExecute recovers panics, enforces timeout, fills duration, and normalizes error signals.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3
- R2.4
- R2.5
- R2.6

<!-- ac: AC3 -->
### AC3

FillDuration and ForceErrorSignal are exported and operate on Result pointers.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2

<!-- field: related_documents -->
## Related documents

- docs/specs/semantic-models/tool-metrics-monitoring.yaml

<!-- field: diagrams -->
## Diagrams

````yaml
requirements_trace: |
    ```mermaid
    flowchart LR
      Problem["Problem: Every command execution needs the same cross-cutting..."] --> Goals["Goals: Command Dispatch"]
      Goals --> R1["R1: Dispatch function"]
      Goals --> R2["R2: SafeExecute function"]
      Goals --> R3["R3: Helper functions"]
      AC1["AC1: Dispatch creates a child span, calls SafeExecute, sets Co..."]
      R1 --> AC1
      AC2["AC2: SafeExecute recovers panics, enforces timeout, fills dura..."]
      R2 --> AC2
      AC3["AC3: FillDuration and ForceErrorSignal are exported and operat..."]
      R3 --> AC3
    ```
````
