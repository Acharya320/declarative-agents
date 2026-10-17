---
id: srd008-runtime-state-reader
title: Runtime State Reader Functional Block
depends_on:
  - srd_id: srd001-agent-functional-blocks
    symbols_used: [R1]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: historical_aliases -->
## Historical aliases

- srd008-monitor

<!-- field: context -->
## Context

The runtime-state-reader profile is an infrastructure adapter over embedded
monitor REST. It serves cached runtime state read-only and exits through
explicit control. This SRD owns the profile-family requirements that
srd001-agent-functional-blocks previously held as block R8.

The adapter is adjacent to, but does not implement, the canonical Monitor
role: it does not collect telemetry, detect anomalies, or produce assurance
reports. The monitor store and REST runtime stay in agent-core; selected routes
and lifecycle control stay profile-owned.

<!-- field: problem -->
## Problem

The runtime-state-reader family needs its own requirements document so changes
to monitor serving behavior do not churn a shared cross-family file. The
requirements must distinguish this profile-owned infrastructure adapter from
both the canonical Monitor role and the agent-core monitor store and REST
runtime.

<!-- field: goals -->
## Goals

- G1: Model runtime-state-reader as a read-only monitor REST infrastructure adapter.
- G2: Record the runtime-state-reader actor boundary, inputs, tool selections, runtime dependencies, outputs, and terminal outcomes.
- G3: Keep monitor routes profile-owned while agent-core supplies the monitor store and REST runtime.

<!-- req: R1 -->
## R1 - Behavior model

<!-- req: R1.1 -->
### R1.1

The runtime-state-reader profile must model read-only monitor REST serving and explicit control-based shutdown.

<!-- req: R1.2 -->
### R1.2

Runtime-state-reader responsibilities must stay limited to profile-selected REST routes and lifecycle control behavior.

<!-- req: R1.3 -->
### R1.3

Runtime-state-reader must be classified as an infrastructure adapter adjacent to canonical Monitor, not as a Monitor realization.

<!-- req: R2 -->
## R2 - Inputs, tools, and runtime dependencies

<!-- req: R2.1 -->
### R2.1

Runtime-state-reader inputs must include live runtime monitor state supplied by agent-core.

<!-- req: R2.2 -->
### R2.2

The runtime-state-reader profile must select the monitor REST and lifecycle control tool families.

<!-- req: R2.3 -->
### R2.3

Runtime-state-reader runtime dependencies are the agent-core monitor store and the agent-core REST runtime.

<!-- req: R3 -->
## R3 - Outputs and outcomes

<!-- req: R3.1 -->
### R3.1

Runtime-state-reader outputs must include monitor REST responses and listener lifecycle signals.

<!-- req: R3.2 -->
### R3.2

Runtime-state-reader terminal outcomes must be Done or Failed.

<!-- field: non_goals -->
## Non-goals

- Define Go package APIs for the agent-core monitor store or REST runtime.
- Reimplement monitor state collection in applications/catalog.
- Claim telemetry collection, anomaly detection, or assurance reporting.
- Own cross-profile dependency-direction invariants; those stay in srd001.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

Runtime-state-reader is described as a read-only infrastructure adapter with control-based shutdown rather than a canonical Monitor realization or Go package.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3

<!-- ac: AC2 -->
### AC2

Runtime-state-reader inputs, selected tool families, and runtime dependencies are named and separated from agent-core runtime capability.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3

<!-- ac: AC3 -->
### AC3

Runtime-state-reader outputs and terminal outcomes are checkable by mounted execution or integration targets.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2

<!-- field: diagrams | form: "tree" -->
## Diagrams

- requirements_trace:
  ```mermaid
  flowchart LR
    Problem["Problem: Runtime-state reader needs a distinct adapter contract..."] --> Goals["Goals: Runtime State Reader Functional Block"]
    Goals --> R1["R1: Behavior model"]
    Goals --> R2["R2: Inputs, tools, and runtime dependencies"]
    Goals --> R3["R3: Outputs and outcomes"]
    R1 --> AC1["AC1: Infrastructure adapter with control shutdown"]
    R2 --> AC2["AC2: Inputs, tools, dependencies named"]
    R3 --> AC3["AC3: Outputs and outcomes checkable"]
  ```
