---
id: srd003-critic
title: Critic Functional Block
depends_on:
  - srd_id: srd001-agent-functional-blocks
    symbols_used: [R1]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: context -->
## Context

The critic family provides two canonical modes. Its benchmark profile is an
evaluation session and point runner: it parses a suite, expands points, runs
a nested point machine, and reports session outcomes. Its changed-workspace
profile evaluates an existing candidate in place and writes an explicit gate
verdict without launching a coding agent. This SRD owns the critic family
requirements that srd001-agent-functional-blocks previously held as block R3.

The block names the contract the critic presents to the rest of the profile
system. Nested loop execution stays in agent-core; the suites and samples stay
profile-owned.

<!-- field: problem -->
## Problem

The critic family needs its own requirements document so that changes to
critic behavior do not churn a shared cross-family file. The requirements
must name profile-owned suites and samples separately from the agent-core
nested loop runtime.

<!-- field: goals -->
## Goals

- G1: Model the critic profile as a configured session and point runner functional block.
- G2: Record the critic actor boundary, inputs, tool selections, runtime dependencies, outputs, and terminal outcomes.
- G3: Keep suites and samples profile-owned while agent-core supplies nested loop execution.
- G4: Evaluate an existing changed workspace and expose an application-neutral accepted or rejected verdict.

<!-- req: R1 -->
## R1 - Behavior model

<!-- req: R1.1 -->
### R1.1

The critic profile must model a deterministic session pipeline that parses a suite, expands evaluation points, runs a nested point machine, and reports session outcomes.

<!-- req: R1.2 -->
### R1.2

The critic point machine must prepare an isolated point workspace, invoke a configured agent profile, run oracle checks, and collect trace, version, summary, and metric outputs.

<!-- req: R1.3 -->
### R1.3

The changed-workspace critic mode must evaluate the supplied workspace in place, branch on the declared oracle result, and must not copy a baseline sample or launch a child coding agent.

<!-- req: R2 -->
## R2 - Inputs, tools, and runtime dependencies

<!-- req: R2.1 -->
### R2.1

Benchmark-mode critic inputs must include a suite YAML, samples, a model or profile grid, and an oracle command.

<!-- req: R2.2 -->
### R2.2

The benchmark profile must select the critic session, critic point, and child profile execution tool families.

<!-- req: R2.3 -->
### R2.3

Critic dependencies must name profile-owned suites and samples as inputs while agent-core supplies nested loop execution.

<!-- req: R2.4 -->
### R2.4

Changed-workspace inputs must include the candidate through the runtime directory binding and a profile-declared oracle whose ToolDone and ToolFailed signals select acceptance and rejection respectively.

<!-- req: R3 -->
## R3 - Outputs and outcomes

<!-- req: R3.1 -->
### R3.1

Critic outputs must include point directories, oracle results, trace token counts, metrics, and a session report.

<!-- req: R3.2 -->
### R3.2

Critic terminal outcomes must be Done or Failed.

<!-- req: R3.3 -->
### R3.3

Changed-workspace output must include critic-verdict.json with a versioned, machine-readable accepted or rejected verdict and the oracle command and status that selected it; session completion alone must never imply acceptance.

<!-- field: non_goals -->
## Non-goals

- Define Go package APIs for the agent-core nested loop runtime.
- Reimplement the point loop or oracle execution in applications/catalog.
- Own cross-profile dependency-direction invariants; those stay in srd001.
- Select an application-specific Succeeded or Failed outcome; consuming applications map the canonical critic verdict.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

The critic profile is described as a configured session and point pipeline rather than a Go package.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2

<!-- ac: AC2 -->
### AC2

Critic inputs, selected tool families, and runtime dependencies are named and separated from agent-core runtime capability.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3

<!-- ac: AC3 -->
### AC3

Critic outputs and terminal outcomes are checkable by mounted execution or integration targets.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2

<!-- ac: AC4 -->
### AC4

The canonical changed-workspace profile accepts a conforming candidate, rejects a non-conforming candidate, and records each decision in critic-verdict.json without a child agent.

<!-- field: traces -->
#### Traces

- R1.3
- R2.4
- R3.3

<!-- field: diagrams | form: "tree" -->
## Diagrams

- requirements_trace:
  ```mermaid
  flowchart LR
    Problem["Problem: Critic family needs its own requirements..."] --> Goals["Goals: Critic Functional Block"]
    Goals --> R1["R1: Behavior model"]
    Goals --> R2["R2: Inputs, tools, and runtime dependencies"]
    Goals --> R3["R3: Outputs and outcomes"]
    R1 --> AC1["AC1: Configured session and point pipeline"]
    R2 --> AC2["AC2: Inputs, tools, dependencies named"]
    R3 --> AC3["AC3: Outputs and outcomes checkable"]
    R1 --> AC4["AC4: Existing candidate evaluated in place"]
    R2 --> AC4
    R3 --> AC4
  ```
