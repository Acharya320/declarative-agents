---
id: srd018-scenario-critic
title: Scenario Critic Functional Block
depends_on:
  - srd_id: srd001-agent-functional-blocks
    symbols_used: [R1]
  - srd_id: srd019-mock
    symbols_used: [R1, R2]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: historical_aliases -->
## Historical aliases

- srd018-assembler

<!-- field: context -->
## Context

The scenario-critic profile is the repository's test rig. It composes other
machines: for each scenario it starts the mocks a subject
depends on, starts the subject agent under test, runs the scenario's
validator machines against it, and reduces their terminal states to a
verdict.

Scenario-critic names the actor's purpose and primary canonical role. Its
scenario composition is an implementation function, not evidence that this
actor realizes the canonical Assembler role. The reference Assembler composes
services and integrates resources for production readiness. This rig instead
performs the Critic's safe-testing activity and binds Critic with Sandbox
Validator and Output Evaluator. The environment it composes is a sandbox in
the catalog's sense, its per-scenario validators exercise expectations in
that sandbox, and its verdict reduction assesses outputs and decisions
against expected outcomes.

The stand-in components inside the sandbox are mocks (srd019-mock); the
catalog has no term for a scripted test double, and its digital twin, a
replica derived from a real referent, is a different concept.

Scenario data is profile-owned and lives beside the agent it tests. The
words the scenario critic machine calls are agent-core runtime capability
(srd040-service-child-tools), and the mock responses its mocks serve come
from the agent-core mock binding (srd039-mock-endpoint-binding).

<!-- field: problem -->
## Problem

Agent tests today are hand-written Go harnesses that restate each agent's
behavior in code, one rig per integration test. That does not reuse across
the profile families, does not extend to the agents under applications/, and
covers REST tools poorly because a subject's dependencies must be live.

One agent under test has several tests. The rig must therefore be reusable
and data-driven: adding a test means adding scenario files beside an agent,
never editing the rig, and never building a container image per test.

<!-- field: goals -->
## Goals

- G1: Model the scenario critic profile as a deterministic scenario composition and verdict pipeline.
- G2: Define the profile-owned scenario layout so one subject can carry several tests.
- G3: Define discovery over configurable roots so agents under agents/ and under applications/ are tested by one rig.
- G4: Record the scenario critic actor boundary, inputs, tool selections, runtime dependencies, outputs, and terminal outcomes.
- G5: Keep scenario data profile-owned while agent-core supplies child-process, validator, and discovery execution.

<!-- req: R1 -->
## R1 - Behavior model

<!-- req: R1.1 -->
### R1.1

The scenario critic profile must model a deterministic pipeline that discovers scenarios, composes each one, collects a verdict, and reports an aggregate outcome.

<!-- req: R1.2 -->
### R1.2

For each scenario the scenario critic must expose fixture and validator arrays as command state; dispatch one mock start per fixture through MachineSpec sequential for_each and one validator child execution per validator through bounded parallel for_each; join each outcome list in input order; start and label the subject; loop through one REST health probe and one delay dispatch per retry; list scenario children; and dispatch stop_service once per child through a final iterator.

<!-- req: R1.3 -->
### R1.3

The scenario critic must give each scenario a fresh subject instance, so scenarios stay independent of one another.

<!-- req: R1.4 -->
### R1.4

Validators declared within one scenario must run against that scenario's single subject instance.

<!-- req: R1.5 -->
### R1.5

The scenario critic must tear down every child it started, including on a failure path, so no child process outlives a scenario.

<!-- req: R1.6 -->
### R1.6

The scenario critic must not invoke a model; its transitions are driven by tool results.

<!-- req: R1.7 -->
### R1.7

Each mock start must have its own trace, checkpoint Entry, opaque receipt, and rollback boundary; the scenario critic must not hide fixture iteration inside a collection Go tool.

<!-- req: R1.8 -->
### R1.8

Health retries must be visible as ProbingHealth and WaitingToRetry states; one probe dispatch performs at most one request, and timeout or budget exhaustion routes through teardown.

<!-- req: R1.9 -->
### R1.9

Each validator execution and child stop must have its own trace, checkpoint Entry, receipt, and rollback boundary; no service word may hide validator or teardown collection traversal.

<!-- req: R2 -->
## R2 - Scenario layout

<!-- req: R2.1 -->
### R2.1

A scenario must live in a tests/<scenario-name>/ directory beside the agent it tests.

<!-- req: R2.2 -->
### R2.2

A scenario directory must hold one or more validator machine programs with their profiles, and the mock fixture files its subject's dependencies require.

<!-- req: R2.3 -->
### R2.3

A subject agent must be able to carry several scenario directories, each self-contained.

<!-- req: R2.4 -->
### R2.4

Scenario data must be profile-owned and must not require changes to the scenario critic profile when a scenario is added.

<!-- req: R3 -->
## R3 - Discovery

<!-- req: R3.1 -->
### R3.1

The scenario critic must enumerate scenarios by convention under roots supplied by configuration or run input.

<!-- req: R3.2 -->
### R3.2

Default roots must cover the repository's agent families and the agents under the applications tree.

<!-- req: R3.3 -->
### R3.3

Discovery must not require a per-repository-area rig, special-casing, or hardcoded absolute paths.

<!-- req: R4 -->
## R4 - Dependency redirection

<!-- req: R4.1 -->
### R4.1

A subject's REST dependencies must be pointed at the scenario's mock addresses through declared environment additions.

<!-- req: R4.2 -->
### R4.2

Redirection must use the existing base_url environment expansion rather than rewriting the subject's behavior.

<!-- req: R4.3 -->
### R4.3

A scenario must not alter the subject's machine, tool selections, or declared behavior; only addresses and other harness-controlled fields may differ.

<!-- req: R5 -->
## R5 - Inputs, tools, and runtime dependencies

<!-- req: R5.1 -->
### R5.1

Scenario Critic inputs must include the discovery roots and the discovered scenario data.

<!-- req: R5.2 -->
### R5.2

The scenario critic profile must select atomic mock and subject starts, one-validator child execution, the generic stop_service word, one REST health probe, one generic delay, and thin scenario discover, next, list, record, verdict, and report words.

<!-- req: R5.3 -->
### R5.3

Scenario Critic dependencies must name agent-core for child-process execution, validator execution, and discovery, and srd019-mock for the mock family it starts.

<!-- req: R5.4 -->
### R5.4

The scenario critic must run the subject and the mocks as the same agent binary over mounted profiles, so no per-test container image is required.

<!-- req: R6 -->
## R6 - Outputs and outcomes

<!-- req: R6.1 -->
### R6.1

A scenario must pass only when every validator in it reaches its passing terminal state.

<!-- req: R6.2 -->
### R6.2

Scenario Critic outputs must record the per-scenario verdict and, for a failure, which validator failed and why, so a failing run names its cause.

<!-- req: R6.3 -->
### R6.3

Scenario Critic terminal outcomes must be Passed or Failed, with the aggregate reflecting every scenario.

<!-- req: R6.4 -->
### R6.4

The scenario, fixture, validator, and teardown loops must be bounded by a declared iteration budget.

<!-- field: non_goals -->
## Non-goals

- Define Go package APIs for the agent-core child-process, validator, or discovery words.
- Own the mock fixture format or the request-log contract; those are srd039 and srd019.
- Kubernetes, container, or hierarchical control-plane deployment of the rig; children are local processes at this release.
- Span-based state assertion; validators witness through responses and the mock request log.
- Own cross-profile dependency-direction invariants; those stay in srd001.
- Parallel fixture dispatch; mock starts stay sequential so each fixture's address is recorded before the subject starts.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

The scenario critic profile is a configured, deterministic scenario composition pipeline; each fixture start and child stop is a machine-visible sequential iterator dispatch and each validator execution a bounded parallel one, each with its own trace, checkpoint Entry, receipt, and rollback boundary.

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

The scenario layout lets one subject carry several self-contained tests as profile-owned data, with no scenario critic change required to add one.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3
- R2.4

<!-- ac: AC3 -->
### AC3

Discovery enumerates scenarios under configurable roots covering both the agent families and the applications tree, with no special-casing or absolute paths.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3

<!-- ac: AC4 -->
### AC4

A subject reaches its scenario's mocks through declared environment additions over the existing base_url expansion, with its declared behavior unchanged.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
- R4.3

<!-- ac: AC5 -->
### AC5

Scenario Critic inputs, selected tool families, and runtime dependencies are named and separated from agent-core runtime capability, and every child runs the same binary over a mounted profile.

<!-- field: traces -->
#### Traces

- R5.1
- R5.2
- R5.3
- R5.4

<!-- ac: AC6 -->
### AC6

A scenario passes only when all its validators pass, a failing run names the failing validator, the terminal outcome is Passed or Failed, and the loop is budget-bounded.

<!-- field: traces -->
#### Traces

- R6.1
- R6.2
- R6.3
- R6.4

<!-- field: diagrams -->
## Diagrams

````yaml
requirements_trace: |
    ```mermaid
    flowchart LR
      Problem["Problem: one rig per test, no reuse, poor REST coverage"] --> Goals["Goals: Scenario Critic Functional Block"]
      Goals --> R1["R1: Behavior model"]
      Goals --> R2["R2: Scenario layout"]
      Goals --> R3["R3: Discovery"]
      Goals --> R4["R4: Dependency redirection"]
      Goals --> R5["R5: Inputs, tools, dependencies"]
      Goals --> R6["R6: Outputs and outcomes"]
      R1 --> AC1["AC1: Deterministic composition pipeline"]
      R2 --> AC2["AC2: Several tests per subject, as data"]
      R3 --> AC3["AC3: Configurable-root discovery"]
      R4 --> AC4["AC4: Mock redirection, behavior unchanged"]
      R5 --> AC5["AC5: Inputs, tools, dependencies named"]
      R6 --> AC6["AC6: Verdict, cause, terminal, budget"]
    ```
scenario_pipeline: |
    ```mermaid
    flowchart LR
      Discover["DiscoverScenarios"] --> Next["NextScenario"]
      Next --> Fixtures["Fixture array"]
      Fixtures --> Mocks["for_each StartMock"]
      Mocks --> Join["Join mock outcomes"]
      Join --> Subject["StartSubject"]
      Subject --> Probe["ProbingHealth: one REST GET"]
      Probe -->|unavailable/error| Delay["WaitingToRetry: one delay"]
      Delay -->|retry ready| Probe
      Probe -->|healthy| Validate["for_each RunValidator"]
      Validate --> VJoin["Join validator outcomes"]
      VJoin --> Record["Record outcomes"]
      Record --> Collect["CollectVerdict"]
      Collect --> Children["List scenario children"]
      Children --> Stop["for_each StopService"]
      Stop --> TJoin["Join child stops"]
      TJoin --> Next
      Mocks -->|failure| Children
      Subject -->|failure| Children
      Delay -->|HealthTimeout| Children
    ```
````

<!-- field: platform_dependencies_note -->
## Platform dependencies note

These SRDs live in the agent-core corpus, not this one; they are platform dependencies rather than in-corpus depends_on (GH-2310).

<!-- field: platform_dependencies -->
## Platform dependencies

```yaml
- srd_id: srd006-yaml-machine-config
  provider: agent-core
  symbols_used: [R5]
```
