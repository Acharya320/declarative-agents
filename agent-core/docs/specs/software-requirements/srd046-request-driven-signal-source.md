---
id: srd046-request-driven-signal-source
title: Request-Driven Signal Source
depends_on:
  - srd_id: srd001-core-types
    symbols_used: [State, Signal, Result]
  - srd_id: srd002-state-machine
    symbols_used: [StateMachine, TransitionTable]
  - srd_id: srd005-agentic-loop
    symbols_used: [Loop, LoopParams, RunResult]
  - srd_id: srd024-blocking-input-tools
    symbols_used: [external input, Result mapping]
  - srd_id: srd025-rollback-lifecycle
    symbols_used: [suspend, resume]
  - srd_id: srd035-checkpoint-port
    symbols_used: [Checkpoint, Position, Execution, OutputRedaction]
  - srd_id: srd029-rest-server-tools
    symbols_used: [configured HTTP request source]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Agent Core currently reaches useful machine signals through dispatched words,
most visibly model-boundary words. A profile that only reacts to a webhook,
timer, queue delivery, or other validated request should not need a model or
an invoke_llm declaration. It still needs the same MachineSpec transition
semantics, lifecycle behavior, validation, and observability as every other
configured program.

We need a transport-neutral signal-source seam plus an HTTP adapter. Trusted
configuration maps validated requests to declared signals and payloads. The
runtime admits that signal against the current machine state, runs the normal
loop, and returns an explicit admission outcome. Callers cannot invent signal
names or program authority. The existing model-driven path remains unchanged.

<!-- field: goals -->
## Goals

- G1: Let a profile run without a model when its selected vocabulary and machine do not use a model boundary.
- G2: Define a generic signal-source envelope and admission contract reusable by HTTP, timers, and queue consumers.
- G3: Map validated inbound requests to ordinary MachineSpec signals and Result payloads without a second transition engine.
- G4: Define accepted, refused_undeclared, and refused_conflict outcomes precisely.
- G5: Preserve checkpoint, suspend, resume, redaction, and observability behavior for request-driven runs.
- G6: Keep the model-driven startup and execution path behaviorally unchanged.

<!-- req: R1 -->
## R1 - Profile and source scope

<!-- req: R1.1 -->
### R1.1

A profile SHALL be valid without an invoke_llm ToolDef, model name, provider URL, response profile, or model client when no selected action requires the model boundary.

<!-- req: R1.2 -->
### R1.2

Model configuration SHALL remain tool-owned and SHALL be required only by selected model-boundary words.

<!-- req: R1.3 -->
### R1.3

Agent Core SHALL define a transport-neutral SignalSource port whose adapters map one validated source request to one SignalEnvelope.

<!-- req: R1.4 -->
### R1.4

SignalEnvelope SHALL carry source name, request id, run id, mapped signal, structured payload, and typed sensitive-payload paths; transport-specific request objects SHALL NOT enter internal/runtime/core.

<!-- req: R1.5 -->
### R1.5

The first adapter SHALL be an HTTP request source under the configured REST boundary; later timer and queue adapters SHALL use the same SignalEnvelope and admission contract.

<!-- req: R2 -->
## R2 - Trusted request mapping

<!-- req: R2.1 -->
### R2.1

Trusted profile or source configuration SHALL declare the source, route or input binding, request schema, run-id mapping, allowed signal mapping, payload mapping, sensitive fields, timeout, and response mapping.

<!-- req: R2.2 -->
### R2.2

A caller SHALL NOT provide a MachineSpec path, profile path, tool declaration path, model config, raw signal name, checkpoint connection, or other program authority.

<!-- req: R2.3 -->
### R2.3

A request discriminator MAY select among configured signals only through a closed mapping whose values are machine-declared signals; an unmapped discriminator SHALL be refused_undeclared.

<!-- req: R2.4 -->
### R2.4

The adapter SHALL validate and map the complete request before admission. Invalid transport or schema input SHALL return a source-validation error and SHALL NOT create, load, or mutate a run.

<!-- req: R2.5 -->
### R2.5

Sensitive mapped fields SHALL become typed Result output-redaction paths before the payload enters InitialResult or Execution; malformed redaction metadata SHALL omit the payload according to srd035.

<!-- req: R3 -->
## R3 - Atomic signal admission

<!-- req: R3.1 -->
### R3.1

Admission SHALL evaluate one SignalEnvelope atomically against the run's current Position and MachineSpec transition table.

<!-- req: R3.2 -->
### R3.2

accepted SHALL mean the mapped signal is declared and the current state has exactly one transition for that signal, and the runtime has acquired the run's single-dispatch ownership.

<!-- req: R3.3 -->
### R3.3

refused_undeclared SHALL mean the source mapping is unknown, the mapped signal is absent from MachineSpec.signals, or trusted config maps to a signal the loaded machine does not declare.

<!-- req: R3.4 -->
### R3.4

refused_conflict SHALL mean the signal is declared but the current state has no transition for it, the caller's configured expected-state guard is stale, or another request already owns dispatch for the same run id.

<!-- req: R3.5 -->
### R3.5

A refused request SHALL NOT call StateMachine.Step, dispatch a command, append Execution, save a checkpoint, or change Position.

<!-- req: R3.6 -->
### R3.6

Admission SHALL serialize requests per run id. It SHALL NOT queue, merge, reorder, or retry a conflicting request implicitly.

<!-- req: R4 -->
## R4 - Ordinary machine execution and response

<!-- req: R4.1 -->
### R4.1

For an accepted request, the mapped signal SHALL enter the existing StateMachine.Step path and the mapped payload SHALL be the corresponding core.Result.Output; no request-specific transition table or direct action dispatch SHALL exist.

<!-- req: R4.2 -->
### R4.2

The accepted transition SHALL resolve its action through the profile-selected Registry and execute through core.Loop and core.Dispatch with normal timeout, budget, trace, event, cost, and terminal-status behavior.

<!-- req: R4.3 -->
### R4.3

The source response SHALL include outcome, source, request_id, run_id, signal, state_before, state_after, run_status, and a redacted diagnostic or terminal output allowed by configured response mapping.

<!-- req: R4.4 -->
### R4.4

The HTTP adapter SHALL map accepted, refused_undeclared, refused_conflict, source-validation error, and machine-run failure to distinct configured HTTP responses; the default conflict response SHALL use HTTP 409.

<!-- req: R4.5 -->
### R4.5

accepted describes transition admission, not command success. A later CommandError or failed terminal SHALL remain visible in run_status and response diagnostics without rewriting the admission outcome.

<!-- req: R5 -->
## R5 - Checkpoint, suspend, and resume

<!-- req: R5.1 -->
### R5.1

A fresh request-driven run SHALL use the mapped signal and payload as InitialResult rather than requiring the ordinary empty Seed bootstrap.

<!-- req: R5.2 -->
### R5.2

A request targeting a suspended persisted run SHALL Load Position and Execution through the configured srd035 Checkpoint, validate declarative-program identity, and use the accepted request signal and payload as the first resumed Result.

<!-- req: R5.3 -->
### R5.3

Request-driven resume SHALL use the same core resume and loop path as other runs; it SHALL NOT replay the original request or prepend a duplicate source entry.

<!-- req: R5.4 -->
### R5.4

After admission, Save timing, Execution entries, receipts, Position, terminal persistence, and StatusSuspended behavior SHALL remain the contracts of srd005, srd025, srd035, and srd036.

<!-- req: R5.5 -->
### R5.5

refused_undeclared and refused_conflict SHALL leave the last durable checkpoint unchanged.

<!-- req: R5.6 -->
### R5.6

A profile that allows continuation after process restart SHALL configure a persistent checkpoint. With NoopCheckpoint, request-driven execution MAY start fresh runs but SHALL refuse a request that claims a suspended run.

<!-- req: R6 -->
## R6 - Safety, observability, and compatibility

<!-- req: R6.1 -->
### R6.1

Source traces SHALL record configured source identity, bounded route identity, request id, run id, admission outcome, mapped signal, state before and after, run status, and elapsed time without recording raw sensitive payloads.

<!-- req: R6.2 -->
### R6.2

Concurrent requests, validation failures, undeclared signals, conflicts, cancellation, timeout, checkpoint load/save failures, and command failures SHALL have distinct structured stages.

<!-- req: R6.3 -->
### R6.3

Signal source adapters SHALL release request ownership on every terminal, suspended, failed, cancelled, and error path.

<!-- req: R6.4 -->
### R6.4

The model-driven path SHALL retain its current profile loading, invoke_llm dispatch, signal mapping, checkpoint, and resume semantics.

<!-- req: R6.5 -->
### R6.5

A mock-model regression SHALL run an existing model-bound profile through the unchanged path, while a separate model-free profile proof SHALL start and complete with no model client construction.

<!-- field: non_goals -->
## Non-goals

- Defining a task, blackboard, webhook, timer, or queue consumer domain model.
- Letting callers choose arbitrary signals, machines, tools, model settings, or checkpoint targets.
- Adding concurrent command execution or hidden request queues to core.Loop.
- Replacing machine_request for short-lived REST sentences whose request and response are already fully owned by one configured endpoint binding.
- Changing invoke_llm, model response parsing, or the model-driven signal path.

<!-- field: test_trace_plan | form: "typed table" -->
## Test trace plan

| test_suite | covers |
| --- | --- |
| test-rel14.0-agent-input-and-dolt | [AC1, AC2, AC3, AC4, AC5, AC6] |

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

Profiles without selected model words require no model configuration or client, and the transport-neutral SignalSource and SignalEnvelope contract keeps HTTP objects outside core.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R1.4
- R1.5

<!-- ac: AC2 -->
### AC2

Trusted source configuration maps validated input to a closed signal and redacted payload contract while callers receive no program or signal authority.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3
- R2.4
- R2.5

<!-- ac: AC3 -->
### AC3

Atomic admission distinguishes accepted, refused_undeclared, and refused_conflict, serializes each run, and guarantees refusals have no runtime or durable mutation.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3
- R3.4
- R3.5
- R3.6

<!-- ac: AC4 -->
### AC4

An accepted envelope enters the ordinary StateMachine.Step, Registry, Dispatch, and Loop path, and the configured response distinguishes admission from later run success or failure.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
- R4.3
- R4.4
- R4.5

<!-- ac: AC5 -->
### AC5

Fresh and persisted request-driven runs use one resume path, preserve normal checkpoint and suspend semantics, never mutate checkpoints on refusal, and reject unavailable suspended runs under NoopCheckpoint.

<!-- field: traces -->
#### Traces

- R5.1
- R5.2
- R5.3
- R5.4
- R5.5
- R5.6

<!-- ac: AC6 -->
### AC6

Structured traces and errors cover every source outcome, ownership is always released, a model-free proof constructs no model client, and a mock-model regression proves the existing model path unchanged.

<!-- field: traces -->
#### Traces

- R6.1
- R6.2
- R6.3
- R6.4
- R6.5

<!-- field: related_documents -->
## Related documents

- docs/ARCHITECTURE.md
- docs/specs/config-formats/runtime-contract.md
- docs/specs/semantic-models/machine-request-http.md
- docs/guides/rest-tool.md

<!-- field: implementation_status -->
## Implementation status

<!-- field: specified -->
### Specified

- Signal-source envelope, trusted mapping, atomic admission outcomes, ordinary machine execution, lifecycle semantics, and compatibility proof.

<!-- field: implemented -->
### Implemented

- GH-1742: transport-neutral SignalEnvelope, admission ownership, ordinary Loop execution, checkpoint behavior, and lifecycle outcomes.
- GH-1741: strict REST signal_source definition, request mapping, response policy, redaction, and trace behavior.
- GH-1743: profile-detected model-free serving, per-RunID in-memory or Dolt checkpoints, host state restoration, owned server shutdown, core-owned fixtures, and mock-model runOrResume regression.
