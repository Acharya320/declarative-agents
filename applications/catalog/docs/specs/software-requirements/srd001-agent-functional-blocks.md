---
id: srd001-agent-functional-blocks
title: Agent Functional Blocks
depends_on: []
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: family_srds -->
## Family srds

- srd002-executor
- srd003-critic
- srd004-planner
- srd005-specification-critic
- srd006-bench
- srd008-runtime-state-reader
- srd011-knowledge-manager
- srd012-chroma-corpus-agents
- srd018-scenario-critic
- srd019-mock
- srd020-collector
- srd022-applier
- srd023-rig-doctor

<!-- field: conformance_fixture_srds -->
## Conformance fixture srds

- srd007-rest
- srd009-lifecycle
- srd010-control

<!-- field: context -->
## Context

This SRD holds the cross-cutting requirements that every migrated profile
family shares: the profile ownership boundary between applications/catalog and
agent-core, and the dependency-direction invariants that connect the families.

Each profile family owns its own SRD for actor boundary, inputs, outputs, tool
selections, runtime dependencies, and terminal outcomes:

  - srd002-executor             (executor)
  - srd003-critic               (critic)
  - srd004-planner              (planner)
  - srd005-specification-critic (specification-critic)
  - srd006-bench                (bench)
  - srd008-runtime-state-reader (runtime-state-reader infrastructure adapter)
  - srd011-knowledge-manager    (knowledge-manager)
  - srd012-chroma-corpus-agents (corpus-ingest and corpus-reader)
  - srd018-scenario-critic            (scenario critic)
  - srd019-mock                 (mock)
  - srd020-collector            (collector)
  - srd022-applier              (deployment applier)

The REST, lifecycle, and control contracts (srd007, srd009, and srd010) cover
internal conformance fixtures under testdata/conformance rather than shipped
library members.

A functional block is the contract a profile presents to the rest of the
multi-agent system. The block does not own the Go code that executes a word.
That runtime remains in agent-core.

<!-- field: problem -->
## Problem

Requirements that describe individual profile families were originally held as
ten blocks in this single SRD. Any change to one family churned a shared file,
and the rel07.0 use cases already treat families individually while the
requirements layer did not.

Each family now owns its own SRD. This document keeps only what is genuinely
cross-cutting: the ownership boundary that separates profile-owned configured
behavior from agent-core-owned runtime capability, and the dependency
directions that must survive as the families evolve.

<!-- field: goals -->
## Goals

- G1: Fix the profile ownership boundary between applications/catalog and agent-core.
- G2: Record the cross-profile dependency directions as invariants the families must preserve.
- G3: Route each profile family to its own dedicated SRD.

<!-- req: R1 -->
## R1 - Profile ownership boundary

<!-- req: R1.1 -->
### R1.1

Each migrated profile must keep profile.yaml, machine.yaml, tools.yaml, profile-local declarations, REST definitions, and UI assets under applications/catalog when those files define agent behavior.

<!-- req: R1.2 -->
### R1.2

Migrated profiles must reference agent-core runtime assets through documented core-owned roots, mounted paths, or release bundle paths.

<!-- req: R1.3 -->
### R1.3

Migrated profiles must not copy agent-core Go runtime packages, builtin tool implementations, or container build logic.

<!-- req: R2 -->
## R2 - Cross-profile dependency direction

<!-- req: R2.1 -->
### R2.1

Planner must invoke executor through child profile execution; executor must not depend on planner.

<!-- req: R2.2 -->
### R2.2

Critic point runs may invoke executor or another configured profile as the benchmark subject; the subject profile must not depend on the critic.

<!-- req: R2.3 -->
### R2.3

Bench must launch critic from a human operator action; critic must not depend on bench.

<!-- req: R2.4 -->
### R2.4

Monitor must read cached runtime state supplied by agent-core, and control must route REST events to visible lifecycle words such as exit_agent; the documentation-curator must embed control flow and wait for explicit exit events.

<!-- req: R2.5 -->
### R2.5

Dependency edges must point from configured profiles toward agent-core runtime ports and never in the reverse direction.

<!-- req: R3 -->
## R3 - Per-family requirement ownership

<!-- req: R3.1 -->
### R3.1

Each shipped profile family's actor boundary, inputs, outputs, tool selections, runtime dependencies, and terminal outcomes must live in one of the SRDs listed by family_srds.

<!-- req: R3.2 -->
### R3.2

Per-family SRDs must depend on this SRD for the ownership boundary and dependency-direction invariants.

<!-- req: R3.3 -->
### R3.3

Fixture-only contracts must be listed separately by conformance_fixture_srds and must not be represented as shipped library roles.

<!-- field: non_goals -->
## Non-goals

- Define Go package APIs for agent-core.
- Reimplement builtin tool behavior in applications/catalog.
- Define release container build mechanics.
- Treat migrated profiles as separate binaries.
- Restate per-family functional block requirements; those live in the SRDs listed by family_srds.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

The SRD defines the profile ownership boundary that keeps configured behavior in applications/catalog and runtime capability in agent-core.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3

<!-- ac: AC2 -->
### AC2

The SRD records the cross-profile dependency directions as invariants that the per-family SRDs must preserve.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3
- R2.4
- R2.5

<!-- ac: AC3 -->
### AC3

Each profile family has a dedicated SRD that depends on this document for the shared boundary and invariants.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3

<!-- field: diagrams | form: "tree" -->
## Diagrams

- requirements_trace:
  ```mermaid
  flowchart LR
    Problem["Problem: Per-family requirements churned a shared file..."] --> Goals["Goals: Agent Functional Blocks"]
    Goals --> R1["R1: Profile ownership boundary"]
    Goals --> R2["R2: Cross-profile dependency direction"]
    Goals --> R3["R3: Per-family requirement ownership"]
    R1 --> AC1["AC1: Ownership boundary defined"]
    R2 --> AC2["AC2: Dependency directions recorded"]
    R3 --> AC3["AC3: Per-family SRDs depend on this document"]
  ```
