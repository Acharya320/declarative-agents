---
id: srd010-control
title: Control Functional Block
depends_on:
  - srd_id: srd001-agent-functional-blocks
    symbols_used: [R1]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: context -->
## Context

The control profile is a REST lifecycle control block. It intakes REST control
requests as signals and routes them to explicit lifecycle vocabulary such as
exit_agent. This SRD owns the control family requirements that
srd001-agent-functional-blocks previously held as block R10.

The block names the contract the control profile presents to the rest of the
profile system. The REST runtime and lifecycle vocabulary stay in agent-core;
the routing and machine wiring stay profile-owned.

<!-- field: problem -->
## Problem

The control family needs its own requirements document so that changes to
control routing behavior do not churn a shared cross-family file. The
requirements must name profile-owned signal routing separately from the
agent-core REST runtime and lifecycle vocabulary.

<!-- field: goals -->
## Goals

- G1: Model the control profile as a REST lifecycle control functional block.
- G2: Record the control actor boundary, inputs, tool selections, runtime dependencies, outputs, and terminal outcomes.
- G3: Keep control routing profile-owned while agent-core supplies the REST runtime and lifecycle vocabulary.

<!-- req: R1 -->
## R1 - Behavior model

<!-- req: R1.1 -->
### R1.1

The control profile must model REST lifecycle control as signal intake followed by explicit lifecycle vocabulary.

<!-- req: R1.2 -->
### R1.2

Control HTTP handlers must enqueue signals only; MachineSpec owns routing to exit_agent or unsupported terminal paths.

<!-- req: R2 -->
## R2 - Inputs, tools, and runtime dependencies

<!-- req: R2.1 -->
### R2.1

Control inputs must include REST control requests.

<!-- req: R2.2 -->
### R2.2

The control profile must select the control REST and exit_agent tool families.

<!-- req: R2.3 -->
### R2.3

Control runtime dependencies are the agent-core REST runtime and the agent-core lifecycle vocabulary.

<!-- req: R3 -->
## R3 - Outputs and outcomes

<!-- req: R3.1 -->
### R3.1

Control outputs must include accepted control events, AgentExited or failed states, and listener lifecycle behavior.

<!-- req: R3.2 -->
### R3.2

Control terminal outcomes must be Succeeded or Failed.

<!-- field: non_goals -->
## Non-goals

- Define Go package APIs for the agent-core REST runtime or lifecycle vocabulary.
- Reimplement lifecycle word execution in applications/catalog.
- Own cross-profile dependency-direction invariants; those stay in srd001.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

The control profile is described as signal intake routed to explicit lifecycle vocabulary rather than a Go package.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2

<!-- ac: AC2 -->
### AC2

Control inputs, selected tool families, and runtime dependencies are named and separated from agent-core runtime capability.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3

<!-- ac: AC3 -->
### AC3

Control outputs and terminal outcomes are checkable by mounted execution or integration targets.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2

<!-- field: diagrams -->
## Diagrams

````yaml
requirements_trace: |
    ```mermaid
    flowchart LR
      Problem["Problem: Control family needs its own requirements..."] --> Goals["Goals: Control Functional Block"]
      Goals --> R1["R1: Behavior model"]
      Goals --> R2["R2: Inputs, tools, and runtime dependencies"]
      Goals --> R3["R3: Outputs and outcomes"]
      R1 --> AC1["AC1: Signal intake routed to lifecycle vocabulary"]
      R2 --> AC2["AC2: Inputs, tools, dependencies named"]
      R3 --> AC3["AC3: Outputs and outcomes checkable"]
    ```
````
