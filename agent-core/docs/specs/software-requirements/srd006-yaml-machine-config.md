---
id: srd006-yaml-machine-config
title: YAML State Machine Configuration
depends_on:
  - srd_id: srd001-core-types
    symbols_used: [State, Signal, Result, Command]
  - srd_id: srd002-state-machine
    symbols_used: [TransitionTable, TransitionInput, TransitionValue, ActionFunc, TerminalFunc]
  - srd_id: srd004-tool-registry
    symbols_used: [Registry]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Hardcoded transition tables require Go code changes to define new
agent types. A declarative YAML format lets users define state
machines by editing a configuration file. Agent-core provides the
schema, parser, validator, and table builder.

<!-- field: goals -->
## Goals

- G1: YAML schema for declarative state machine definition.
- G2: Validation with referential integrity checks.
- G3: BuildTransitionTable converting YAML to a core.TransitionTable.
- G4: Trusted data-driven iteration with deterministic join semantics.
- G5: Workload-derived command timeout envelopes for shipped machines.

<!-- req: R1 -->
## R1 - YAML schema (MachineSpec)

<!-- req: R1.1 -->
### R1.1

Top-level fields: name (string), initial_state (string), states (list), terminal_states (list), signals (list), transitions (list of TransitionSpec), budget (required BudgetSpec).

<!-- req: R1.2 -->
### R1.2

TransitionSpec fields: state (string), signal (string), next (string), action (string, optional), label (string, optional), and for_each (optional ForEachSpec).

<!-- req: R1.3 -->
### R1.3

action is either a tool name resolvable from the Registry, the special sentinel $tool for dynamic dispatch, or empty for terminal transitions.

<!-- req: R1.4 -->
### R1.4

$tool is resolved using a caller-supplied ActionFunc for runtime tool lookup.

<!-- req: R1.5 -->
### R1.5

label is the stable, authored command-state address for the action executed by that transition. It is distinct from action and metric_labels.

<!-- req: R1.6 -->
### R1.6

When label is absent, the completed step remains addressable by its executed command name. This preserves existing unlabeled machines and programmatic TransitionTable runs with no MachineSpec label.

<!-- req: R1.7 -->
### R1.7

A labeled $tool transition records label as the stable workflow-step address and records the dynamically selected tool separately as the command name.

<!-- req: R1.8 -->
### R1.8

external_labels (list of strings, optional) names command-state labels the runtime seeds rather than any transition publishing them. Each entry satisfies the transition-label grammar and widens the machine's label universe for selector validation. The REST machine_request seed label is the first user.

<!-- req: R1.9 -->
### R1.9

expand (list, optional) names stage fragments the machine splices in (srd052 R4). Each entry carries fragment, a path relative to the machine file, and args. A stage fragment declares params and a stage body of states, signals, and transitions; the loader fills the arguments into scalar values, appends the produced states, signals, and transitions to the machine's lists, and validates the whole as one machine, so a spliced transition is held to R2 exactly as a hand-written one is. A machine declares no imports: a stage is expanded, and a machine has no other unit kind to import. An expand entry whose fragment has a machine body instead makes the file an instance of a machine template, which replaces the machine rather than splicing into it (srd054).

<!-- req: R2 -->
## R2 - Validation

<!-- req: R2.1 -->
### R2.1

initial_state must appear in the states list.

<!-- req: R2.2 -->
### R2.2

Every terminal_state must appear in the states list.

<!-- req: R2.3 -->
### R2.3

Every transition's state and next fields must appear in the states list.

<!-- req: R2.4 -->
### R2.4

Every transition's signal must appear in the signals list.

<!-- req: R2.5 -->
### R2.5

At least one state, one terminal_state, and one transition are required.

<!-- req: R2.6 -->
### R2.6

Validation errors are collected and reported together, not one at a time.

<!-- req: R2.7 -->
### R2.7

A non-empty transition label must satisfy the canonical $from(label).path label grammar: no dots, parentheses, whitespace, or control characters. The error identifies the transition index and label field.

<!-- req: R2.8 -->
### R2.8

Duplicate transition labels are valid because command-state lookup resolves the most recent matching completed step.

<!-- req: R2.9 -->
### R2.9

Every $from(label) selector reachable from the machine resolves its label against the machine's label universe at load: transition labels, for_each join and item labels, and external_labels. Selectors are taken from the machine's for_each items and from the config of every tool a transition names as an action. An unresolved label is a load error naming the selector, the unknown label, the tool or state carrying it, and the nearest declared label.

<!-- req: R2.10 -->
### R2.10

Every transition waits on a signal something entering its state can produce. The producible set is what the actions entering that state emit, the signals a for_each join publishes into its own next state, the runtime seed, the machine's declared resume_signal and summary_signal, and the signals a request signal source can inject. A transition outside that set is a load error naming the transition and its state, because it can never fire.

<!-- req: R2.11 -->
### R2.11

An action declaring no emitted signals contributes nothing to the producible set rather than being read as emitting nothing. The producible set is therefore what the declarations state, never what their absence implies.

<!-- req: R3 -->
## R3 - BuildTransitionTable

<!-- req: R3.1 -->
### R3.1

BuildTransitionTable(spec MachineSpec, reg *Registry, toolAction ActionFunc) returns (TransitionTable, TerminalFunc, error).

<!-- req: R3.2 -->
### R3.2

For named actions, reg.Resolve(action) obtains the Builder. Resolution failure is an error.

<!-- req: R3.3 -->
### R3.3

For "$tool" action, the caller-supplied toolAction is used. Nil toolAction with "$tool" is an error.

<!-- req: R3.4 -->
### R3.4

Empty action produces nil ActionFunc (terminal transition).

<!-- req: R3.5 -->
### R3.5

TerminalFunc is constructed from the terminal_states list.

<!-- req: R4 -->
## R4 - File loading

<!-- req: R4.1 -->
### R4.1

LoadMachineSpec(path string) reads and parses a YAML file.

<!-- req: R4.2 -->
### R4.2

ParseMachineSpec(data []byte) parses YAML bytes and validates.

<!-- req: R4.3 -->
### R4.3

File read errors and YAML parse errors are wrapped with context.

<!-- req: R4.4 -->
### R4.4

LoadMachineClosure(path, visit) reads the machine file and every stage fragment it expands through visit, splices, and validates. LoadMachineSpec is that with no visitor. ParseMachineSpec refuses a document that expands, because its fragments are named relative to a file it does not have.

<!-- req: R5 -->
## R5 - Data-driven iteration and join

<!-- req: R5.1 -->
### R5.1

for_each.items must be a $from(label).path selector that resolves at runtime to a JSON array; for_each.as is a valid temporary command-state label.

<!-- req: R5.2 -->
### R5.2

for_each supports sequential dispatch and bounded parallel dispatch with a positive max_concurrency. Each item is built through the named Registry action and executed through normal core dispatch, producing its own trace, budget iteration, checkpoint Entry, and opaque receipt.

<!-- req: R5.3 -->
### R5.3

The iterator supports fail_fast and collect_all policies and emits one declared all_success, partial, failed, or empty join signal with outcomes retained in input order.

<!-- req: R5.4 -->
### R5.4

The join itself is recorded as an engine-owned Entry whose labeled output reports ordered outcomes and succeeded and failed counts.

<!-- req: R5.5 -->
### R5.5

AgentSnapshot persists selected items, the next item index, outcomes, and iterator grammar so resume does not replay completed items.

<!-- req: R5.6 -->
### R5.6

Cancellation stops scheduling items, cancels contextual children, waits for all started work, and leaves receipt-bearing item Entries visible to the existing reverse rollback walk.

<!-- req: R5.7 -->
### R5.7

for_each must reject dynamic $tool actions, malformed selectors and labels, undeclared signals or states, unsupported modes, invalid concurrency bounds, and invalid policies during startup validation.

<!-- req: R5.8 -->
### R5.8

Runtime data may supply items only; it cannot supply actions or executable expressions. The iterator does not provide general scripting.

<!-- req: R5.9 -->
### R5.9

Each join outcome SHALL retain result.output as a string and SHALL add result.structured_output only when the already-redacted ResultDigest output has recognized applied redaction metadata and is valid JSON. The projection SHALL preserve object, array, and scalar JSON shapes, SHALL be absent for plain text and output_omitted results, and SHALL be shaped at join time without changing the persisted IteratorOutcome format.

<!-- req: R6 -->
## R6 - Command timeout envelope

<!-- req: R6.1 -->
### R6.1

Every shipped MachineSpec SHALL declare a finite, positive budget.command_timeout.

<!-- req: R6.2 -->
### R6.2

The declared command_timeout SHALL exceed the longest legitimate duration of every action reachable in that machine, including scheduling margin.

<!-- req: R6.3 -->
### R6.3

Calibration SHALL use the selected ToolDef timeout, REST limit or machine_request timeout, model or child-process limit, and blocking await duration as authorities; one copied repository-wide value SHALL NOT substitute for workload review.

<!-- req: R6.4 -->
### R6.4

A selected operation with an environment-overridable timeout SHALL resolve to a finite value strictly below the machine command_timeout.

<!-- req: R6.5 -->
### R6.5

Each owning module SHALL statically test its shipped machine-to-authority mapping without sleeping for the declared wall-clock duration.

<!-- req: R6.6 -->
### R6.6

A reusable profile inspector SHALL resolve declaration includes and override precedence, selected named and dynamic actions, REST definitions, machine_request profile or machine references, child profiles, evaluator point machines, and compatibility wrappers before comparing authorities.

<!-- req: R6.7 -->
### R6.7

The inspector SHALL classify ToolDef, REST, model, child/service, evaluator-point, and await authorities; it SHALL fail closed when a reachable timeout-bearing boundary action has no supported finite authority.

<!-- req: R6.8 -->
### R6.8

Environment expansion SHALL precede comparison. The active MachineSpec command_timeout SHALL be the finite runtime maximum, and an expanded operation timeout that is malformed, non-positive, equal to, or greater than that maximum SHALL fail startup before dispatch.

<!-- req: R6.9 -->
### R6.9

Diagnostics SHALL deterministically identify the profile, machine, selected action, authority source, operation duration, and command timeout.

<!-- req: R6.10 -->
### R6.10

Agent Core SHALL inspect only a caller-supplied profile closure. Catalog and application owners SHALL discover their own shipped profile inventories and invoke the shared inspector without Agent Core runtime tests discovering product source trees.

<!-- field: non_goals -->
## Non-goals

- Dynamic table modification during a run.
- Tool registration from YAML (tools are registered in Go and resolved by name).
- CLI integration (domain agent responsibility).
- Unbounded parallel iteration.
- A general expression, mutation, or scripting language in machine YAML.
- Propagating request-machine timeout policy into nested runners; that runtime wiring is tracked separately.
- Treating a Helm, JSON Schema, or other deployment-input limit as the runtime timeout safety boundary.
- Inventing a second profile registry or compatibility-alias inventory for timeout inspection.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

MachineSpec has name, initial_state, states, terminal_states, signals, transitions, and a required budget with a positive command timeout. TransitionSpec has state, signal, next, optional action, and optional command-state label. An unlabeled transition falls back to its executed command name, including programmatic tables without MachineSpec labels, while a labeled $tool transition keeps its stable label distinct from the selected command.

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

<!-- ac: AC2 -->
### AC2

Validation checks referential integrity for initial_state, terminal_states, transitions, minimum counts, and transition-label grammar. Duplicate valid labels remain legal. Errors are collected together.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3
- R2.4
- R2.5
- R2.6
- R2.7
- R2.8
- R2.9
- R2.10
- R2.11

<!-- ac: AC3 -->
### AC3

BuildTransitionTable resolves named actions from Registry, uses toolAction for "$tool", nil for empty. Returns TransitionTable, TerminalFunc, and error.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3
- R3.4
- R3.5

<!-- ac: AC4 -->
### AC4

LoadMachineSpec reads and parses YAML files, splicing any stage fragments the machine expands and validating the result as one machine. ParseMachineSpec parses bytes with validation and refuses a document that expands. Errors include context, naming the fragment where one is involved.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
- R4.3
- R4.4

<!-- ac: AC5 -->
### AC5

A trusted command-state array drives sequential or positively bounded parallel named-action dispatches with per-item trace, checkpoint, receipt, cancellation, rollback visibility, deterministic input-order join output, and checkpoint-resume progress; fail-fast waits for all started children and collect-all retains every outcome.

<!-- field: traces -->
#### Traces

- R5.1
- R5.2
- R5.3
- R5.4
- R5.5
- R5.6
- R5.7
- R5.8
- R5.9

<!-- ac: AC6 -->
### AC6

Every shipped machine declares a finite command timeout strictly above each closure-resolved selected operation. The reusable inspector resolves includes, overrides, aliases, request and nested machines, classifies all supported authority families, rejects unresolved or oversized environment-expanded values before dispatch, emits deterministic evidence, and is driven by each owner's profile inventory without cross-owner discovery.

<!-- field: traces -->
#### Traces

- R6.1
- R6.2
- R6.3
- R6.4
- R6.5
- R6.6
- R6.7
- R6.8
- R6.9
- R6.10

<!-- field: diagrams | form: "tree" -->
## Diagrams

- requirements_trace:
  ```mermaid
  flowchart LR
    Problem["Problem: Hardcoded transition tables require Go code changes t..."] --> Goals["Goals: YAML State Machine Configuration"]
    Goals --> R1["R1: YAML schema (MachineSpec)"]
    Goals --> R2["R2: Validation"]
    Goals --> R3["R3: BuildTransitionTable"]
    Goals --> R4["R4: File loading"]
    Goals --> R5["R5: Data-driven iteration and join"]
    Goals --> R6["R6: Command timeout envelope"]
    AC1["AC1: MachineSpec has name, initial_state, states, terminal_sta..."]
    R1 --> AC1
    AC2["AC2: Validation checks referential integrity for initial_state..."]
    R2 --> AC2
    AC3["AC3: BuildTransitionTable resolves named actions from Registry..."]
    R3 --> AC3
    AC4["AC4: LoadMachineSpec reads and parses YAML files. ParseMachine..."]
    R4 --> AC4
    AC5["AC5: trusted bounded iteration, join, and resume"]
    R5 --> AC5
    AC6["AC6: workload-derived finite command timeout envelopes"]
    R6 --> AC6
  ```
