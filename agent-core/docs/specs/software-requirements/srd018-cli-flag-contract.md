---
id: srd018-cli-flag-contract
title: CLI Standard Flag Contract
depends_on:
  - srd_id: srd016-traceparent
    symbols_used: [FormatTraceparent, ParseTraceparent]
  - srd_id: srd009-ollama-adapter
    symbols_used: [Adapter]
  - srd_id: srd035-checkpoint-port
    symbols_used: [Checkpoint, NoopCheckpoint]
  - srd_id: srd036-dolt-state-persistence
    symbols_used: [DoltCheckpoint, OpenDoltCheckpoint]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Agent CLIs share a common set of flags for telemetry propagation, model
selection, timeouts, and working directory. Without a shared contract,
each CLI binary reinvents flag names and argument-building logic,
leading to inconsistencies and broken parent→child span propagation.

<!-- field: goals -->
## Goals

- G1: Define canonical flag name constants so all agent binaries use identical CLI surfaces.
- G2: Provide PropagateArgs for building child-process argument slices with traceparent and budget propagation.
- G3: Provide helpers for remaining-budget calculation and ad-hoc flag formatting.
- G4: Define the persistent checkpoint-backend flag and default-backend behavior for the agent binary.
- G5: Define explicit OTLP destination and service identity flags for direct and child agent runs.

<!-- req: R1 -->
## R1 - Standard flag name constants

<!-- req: R1.1 -->
### R1.1

FlagOTelLogFile ("otel-log-file") — path to the OTLP JSON log file shared across the process tree.

<!-- req: R1.2 -->
### R1.2

FlagOTelParent ("otel-parent-span") — W3C traceparent string linking child spans to the parent trace.

<!-- req: R1.3 -->
### R1.3

FlagModel ("model") — LLM model identifier passed to the Ollama adapter.

<!-- req: R1.4 -->
### R1.4

FlagDirectory ("directory") — working directory the agent operates in.

<!-- req: R1.5 -->
### R1.5

FlagMaxTime ("max-time") — total wall-clock budget for the agent process.

<!-- req: R1.6 -->
### R1.6

FlagLLMTimeout ("llm-timeout") — per-LLM-call timeout.

<!-- req: R1.7 -->
### R1.7

FlagPrompt ("prompt") — path to the prompt/task file.

<!-- req: R1.8 -->
### R1.8

FlagOllamaURL ("ollama-url") — base URL of the Ollama API server.

<!-- req: R1.9 -->
### R1.9

FlagOTelOTLPEndpoint ("otel-otlp-endpoint") — OTLP gRPC endpoint used for live trace and metric export.

<!-- req: R1.10 -->
### R1.10

FlagOTelServiceName ("otel-service-name") — stable OpenTelemetry service.name for the agent process.

<!-- req: R2 -->
## R2 - PropagateArgs struct and Build method

<!-- req: R2.1 -->
### R2.1

PropagateArgs fields: SpanContext (trace.SpanContext), TraceFile (string), OTLPEndpoint (string), OTelServiceName (string), Model (string), Directory (string), Prompt (string), OllamaURL (string), MaxTime (time.Duration), LLMTimeout (time.Duration).

<!-- req: R2.2 -->
### R2.2

Build() serialises SpanContext via telemetry.FormatTraceparent and emits --otel-parent-span only when the span context is valid.

<!-- req: R2.3 -->
### R2.3

Build() omits flags whose values are zero/empty, producing a minimal argument slice.

<!-- req: R2.4 -->
### R2.4

Build() emits durations using Go's Duration.String() format (e.g. "5m0s").

<!-- req: R2.5 -->
### R2.5

Build() SHALL propagate non-empty OTLPEndpoint and OTelServiceName values unchanged so a child joins the same integration ingress under its own declared service identity.

<!-- req: R3 -->
## R3 - RemainingBudget helper

<!-- req: R3.1 -->
### R3.1

RemainingBudget(total, elapsed time.Duration) time.Duration computes the remaining wall-clock budget.

<!-- req: R3.2 -->
### R3.2

Returns 0 when total is zero (unlimited budget) or when elapsed >= total (exhausted).

<!-- req: R4 -->
## R4 - FormatFlag helper

<!-- req: R4.1 -->
### R4.1

FormatFlag(name, value string) []string returns ["--name", "value"] for ad-hoc flags not covered by PropagateArgs.

<!-- req: R5 -->
## R5 - Persistent checkpoint-backend flag

<!-- req: R5.1 -->
### R5.1

The agent binary SHALL expose --dolt-dsn, a MySQL-wire DSN naming a running dolt sql-server, selecting the Dolt-backed persistent checkpoint backend (srd036-dolt-state-persistence R1.4).

<!-- req: R5.2 -->
### R5.2

When --dolt-dsn is empty the agent SHALL default to NoopCheckpoint so a run without persistence keeps disabled-mode behavior (srd035-checkpoint-port R5.1).

<!-- req: R5.3 -->
### R5.3

The "dolt" database/sql driver SHALL be registered at the composition root (cmd/agent); core reaches the database only through the database/sql seam and never imports Dolt (srd036-dolt-state-persistence R1.3).

<!-- req: R5.4 -->
### R5.4

--dolt-dsn SHALL supersede the legacy FileStore state-store directory as the persistence surface; runtime state is versioned by Dolt rather than a state-store directory.

<!-- req: R6 -->
## R6 - Process exit code contract

<!-- req: R6.1 -->
### R6.1

A run whose machine reaches a success terminal state SHALL exit 0, and a run that reaches a failure terminal state SHALL exit non-zero, so a caller running the agent as a child process can read the outcome from the exit status.

<!-- req: R6.2 -->
### R6.2

The failure-terminal exit code SHALL be distinct from the code used when the binary could not complete a run at all (bad configuration, unreadable profile, transport failure), so a caller can tell a clean domain failure from a crash.

<!-- req: R6.3 -->
### R6.3

Success and failure SHALL be derived from the run's terminal RunStatus, not from a hardcoded terminal state name, so a machine may name its terminal states freely.

<!-- req: R6.4 -->
### R6.4

A suspended run SHALL exit 0, because a suspension is a deliberate pause with a persisted checkpoint that a caller resumes rather than a failure (srd025-rollback-lifecycle R5.3).

<!-- req: R6.5 -->
### R6.5

--validate-config SHALL keep its documented preflight codes, 0 for a valid corpus and non-zero for an invalid one, independently of any run.

<!-- req: R7 -->
## R7 - Live telemetry flags

<!-- req: R7.1 -->
### R7.1

--otel-otlp-endpoint SHALL enable live OTLP trace and metric export to the named gRPC endpoint and SHALL compose with --otel-log-file for dual export.

<!-- req: R7.2 -->
### R7.2

--otel-service-name SHALL set the OpenTelemetry service.name resource attribute and SHALL NOT replace the profile as the source of agent program identity.

<!-- req: R7.3 -->
### R7.3

Integration launchers SHALL pass both flags explicitly. Stable names SHALL identify an agent role and, for scaled roles, its instance.

<!-- req: R7.4 -->
### R7.4

An empty --otel-otlp-endpoint SHALL leave live export disabled; an empty --otel-service-name SHALL use the runtime's documented default service name.

<!-- field: non_goals -->
## Non-goals

- Flag registration with cobra/pflag (that is the responsibility of each CLI binary's main).
- Parsing inbound flags (each binary uses its own flag set).
- Environment variable fallbacks.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

Ten standard flag name constants are exported and match the canonical CLI surface.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R1.4
- R1.5
- R1.6
- R1.7
- R1.8
- R1.9
- R1.10

<!-- ac: AC2 -->
### AC2

PropagateArgs.Build() produces a correct flag slice, omitting flags with zero values, formatting traceparent via telemetry.FormatTraceparent, and preserving explicit OTLP endpoint and service identity values.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3
- R2.4
- R2.5

<!-- ac: AC3 -->
### AC3

RemainingBudget returns 0 for unlimited or exhausted budgets, and the correct remainder otherwise.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2

<!-- ac: AC6 -->
### AC6

An agent whose machine reaches a failure terminal exits non-zero with a code distinct from the run-error code, one that reaches a success or suspended terminal exits zero, the mapping is taken from the terminal RunStatus, and --validate-config keeps its own preflight codes independently of any run.

<!-- field: traces -->
#### Traces

- R6.1
- R6.2
- R6.3
- R6.4
- R6.5

<!-- ac: AC4 -->
### AC4

FormatFlag returns a two-element slice with the prefixed flag name and value.

<!-- field: traces -->
#### Traces

- R4.1

<!-- ac: AC5 -->
### AC5

The agent documents --dolt-dsn as the MySQL-wire DSN selecting the Dolt-backed persistent backend, defaults to NoopCheckpoint when empty, registers the "dolt" driver at the composition root, and supersedes the legacy FileStore state-store directory.

<!-- field: traces -->
#### Traces

- R5.1
- R5.2
- R5.3
- R5.4

<!-- ac: AC7 -->
### AC7

The live telemetry flags compose file and OTLP export, set stable service identity without changing profile identity, propagate to child agents, and preserve disabled defaults when omitted.

<!-- field: traces -->
#### Traces

- R7.1
- R7.2
- R7.3
- R7.4

<!-- field: diagrams -->
## Diagrams

````yaml
requirements_trace: |
    ```mermaid
    flowchart LR
      Problem["Problem: Agent CLIs share a common set of flags for telemetry..."] --> Goals["Goals: CLI Standard Flag Contract"]
      Goals --> R1["R1: Standard flag name constants"]
      Goals --> R2["R2: PropagateArgs struct and Build method"]
      Goals --> R3["R3: RemainingBudget helper"]
      Goals --> R4["R4: FormatFlag helper"]
      Goals --> R5["R5: Persistent checkpoint-backend flag"]
      Goals --> R6["R6: Process exit code contract"]
      Goals --> R7["R7: Live telemetry flags"]
      AC1["AC1: Eight standard flag name constants are exported and match..."]
      R1 --> AC1
      AC2["AC2: PropagateArgs.Build() produces a correct flag slice, omit..."]
      R2 --> AC2
      AC3["AC3: RemainingBudget returns 0 for unlimited or exhausted budg..."]
      R3 --> AC3
      AC4["AC4: FormatFlag returns a two-element slice with the prefixed..."]
      R4 --> AC4
      AC5["AC5: --dolt-dsn selects the Dolt backend; Noop default; driver at root..."]
      R5 --> AC5
      AC7["AC7: Live telemetry flags compose dual export and stable service identity..."]
      R7 --> AC7
    ```
````
