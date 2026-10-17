---
id: srd004-tool-registry
title: Tool Registry
depends_on:
  - srd_id: srd001-core-types
    symbols_used: [State, Builder, CommandResolver]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Agentic agents need to register tools at startup and look them up at
runtime. Some tools are visible to the LLM (external), others are
internal. The set of available tools may vary by execution state. A
central registry must pair tool metadata with command builders and
support state-filtered manifest generation.

<!-- field: goals -->
## Goals

- G1: ToolSpec metadata type with visibility and phase filtering.
- G2: Registry pairing ToolSpecs with Builders.
- G3: State-filtered manifest generation.
- G4: Freeze semantics to prevent mutation during the loop.

<!-- req: R1 -->
## R1 - Visibility and ToolSpec

<!-- req: R1.1 -->
### R1.1

Visibility is an int enum with External (0) and Internal (1) constants.

<!-- req: R1.2 -->
### R1.2

ToolSpec has five fields: Name (string), Description (string), InputSchema (json.RawMessage), Visibility (Visibility), Phases ([]State).

<!-- req: R2 -->
## R2 - Registry operations

<!-- req: R2.1 -->
### R2.1

NewRegistry() creates an empty, mutable Registry.

<!-- req: R2.2 -->
### R2.2

Register(spec ToolSpec, builder Builder) adds a pair. Panics on empty name, duplicate name, or post-freeze call.

<!-- req: R2.3 -->
### R2.3

Freeze() marks the registry as immutable for the duration of the run.

<!-- req: R2.4 -->
### R2.4

Resolve(name string) returns (Builder, bool). False if not found.

<!-- req: R2.5 -->
### R2.5

SpecByName(name string) returns (ToolSpec, bool). False if not found.

<!-- req: R3 -->
## R3 - Manifest generation

<!-- req: R3.1 -->
### R3.1

Manifest(state State) returns a sorted copy of all External ToolSpecs whose Phases match the given state.

<!-- req: R3.2 -->
### R3.2

PhaseMatch returns true when Phases is empty (matches all states) or contains the given state.

<!-- req: R4 -->
## R4 - Name listing

<!-- req: R4.1 -->
### R4.1

AllToolNames() returns every registered tool name sorted alphabetically.

<!-- req: R4.2 -->
### R4.2

ExternalToolNames() returns only External tool names sorted alphabetically.

<!-- req: R5 -->
## R5 - Interface conformance

<!-- req: R5.1 -->
### R5.1

Registry satisfies the CommandResolver interface.

<!-- req: R6 -->
## R6 - Signature-derived contracts

<!-- req: R6.1 -->
### R6.1

A tool that declares a signature satisfies the contract blocks that default from its category, per srd051 R6.6 through R6.10. Registry consumers read the effective contract and SHALL NOT distinguish a defaulted block from an authored one.

<!-- req: R6.2 -->
### R6.2

A boundary tool SHALL carry explicit side_effects, reversibility, and undo whether or not it declares a signature (srd051 R6.9).

<!-- field: non_goals -->
## Non-goals

- Tool implementation (domain agent responsibility).
- Dynamic registration during the loop (prevented by freeze).
- Tool input schema validation (consumer responsibility).

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

Visibility enum with External and Internal. ToolSpec with Name, Description, InputSchema, Visibility, and Phases.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2

<!-- ac: AC2 -->
### AC2

Registry supports Register (with panic on violations), Freeze, Resolve, and SpecByName.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3
- R2.4
- R2.5

<!-- ac: AC3 -->
### AC3

Manifest returns sorted External tools filtered by PhaseMatch. PhaseMatch handles empty Phases as match-all.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2

<!-- ac: AC4 -->
### AC4

AllToolNames and ExternalToolNames return sorted name lists.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2

<!-- ac: AC5 -->
### AC5

Registry satisfies CommandResolver interface.

<!-- field: traces -->
#### Traces

- R5.1

<!-- ac: AC6 -->
### AC6

A tool declaring a signature presents an effective contract to registry consumers, and a boundary tool still carries explicit side_effects, reversibility, and undo.

<!-- field: traces -->
#### Traces

- R6.1
- R6.2

<!-- field: diagrams | form: "tree" -->
## Diagrams

- requirements_trace:
  ```mermaid
  flowchart LR
    Problem["Problem: Agentic agents need to register tools at startup and..."] --> Goals["Goals: Tool Registry"]
    Goals --> R1["R1: Visibility and ToolSpec"]
    Goals --> R2["R2: Registry operations"]
    Goals --> R3["R3: Manifest generation"]
    Goals --> R4["R4: Name listing"]
    Goals --> R5["R5: Interface conformance"]
    AC1["AC1: Visibility enum with External and Internal. ToolSpec with..."]
    R1 --> AC1
    AC2["AC2: Registry supports Register (with panic on violations), Fr..."]
    R2 --> AC2
    AC3["AC3: Manifest returns sorted External tools filtered by PhaseM..."]
    R3 --> AC3
    AC4["AC4: AllToolNames and ExternalToolNames return sorted name lis..."]
    R4 --> AC4
    AC5["AC5: Registry satisfies CommandResolver interface."]
    R5 --> AC5
  ```
