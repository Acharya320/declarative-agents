---
id: srd019-mock
title: Mock Functional Block
depends_on:
  - srd_id: srd001-agent-functional-blocks
    symbols_used: [R1]
  - srd_id: srd007-rest
    symbols_used: [R1]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: context -->
## Context

The mock profile is the rig's stand-in family. One shipped profile stands
in for any REST dependency a subject agent calls: it serves the canned,
scripted responses a scenario declares and records the requests it received
so a validator can assert what the subject actually sent.

The name matches the runtime surface it configures, the agent-core mock
endpoint binding, so the profile and its binding agree about what the thing
is. The reference architecture has no term for a scripted test double:
mock, stub, emulator, fixture, and simulator do not appear in its corpus,
which is an autonomous-networks vocabulary rather than a testing one. Its
digital twin is a different concept and is not used here: a twin mirrors a
real referent it is derived from, either a physical asset ("the OTDR trace
becomes a digital twin of the physical fiber") or the live network as a
whole ("sandboxes or digital twins"), and a scripted responder mirrors
nothing. The reference term that does apply at the environment level is the
sandbox, which srd018-scenario-critic uses for the composed scenario the rig
builds; a mock is one component inside it.

The mock is configuration over the agent-core mock endpoint binding
(srd039-mock-endpoint-binding). It is one profile, not one per dependency:
what a mock serves comes from the fixture file a scenario mounts, so a new
dependency needs a fixture, not a new profile or a new container image.

<!-- field: problem -->
## Problem

Testing a subject's REST tools requires controlling what its dependencies
return. Standing up the real services makes tests slow, machine-dependent,
and non-deterministic, and it cannot script a failure on a chosen call. A
purpose-built mock server per dependency would be Go code outside the
profile boundary and would multiply with every dependency.

<!-- field: goals -->
## Goals

- G1: Model the mock profile as one generic, configured REST responder for any dependency.
- G2: Define fixture selection so a scenario supplies what the mock serves without editing the profile.
- G3: Define the request-log surface validators assert against.
- G4: Record the mock actor boundary, inputs, tool selections, runtime dependencies, outputs, and terminal outcomes.
- G5: Keep loopback as the safe listener default while allowing a rig to authorize Kubernetes reachability explicitly.

<!-- req: R1 -->
## R1 - Behavior model

<!-- req: R1.1 -->
### R1.1

The mock profile must model a serve-mode responder that launches a REST server, serves until a lifecycle exit signal, and stops cleanly.

<!-- req: R1.2 -->
### R1.2

The mock must serve its declared routes through the agent-core mock endpoint binding rather than profile-specific response logic.

<!-- req: R1.3 -->
### R1.3

The mock must expose a health route so a rig can await readiness before starting a subject.

<!-- req: R1.4 -->
### R1.4

The mock must not invoke a model.

<!-- req: R1.5 -->
### R1.5

The mock listener must default to loopback; a wildcard or public address must require an explicit, default-false public-listener authorization, and absent or malformed authorization must fail before serving.

<!-- req: R2 -->
## R2 - Fixture selection and scope

<!-- req: R2.1 -->
### R2.1

The fixture file the mock serves must be selected by environment at start, so one profile serves any dependency.

<!-- req: R2.2 -->
### R2.2

A scenario must supply its own fixture files; the mock profile must not carry per-dependency response content.

<!-- req: R2.3 -->
### R2.3

Starting several mocks for one scenario, each with its own fixture and address, must require no profile change.

<!-- req: R2.4 -->
### R2.4

A malformed fixture must prevent the mock from serving, surfacing at configuration validation and at startup rather than on the first request.

<!-- req: R2.5 -->
### R2.5

A scenario fixture must be able to represent Ollama-compatible model discovery and an ordered chat-response script without a model binary or application-specific mock profile.

<!-- req: R3 -->
## R3 - Request log

<!-- req: R3.1 -->
### R3.1

The mock must expose the recorded request log on a declared route so a validator can read what the subject sent.

<!-- req: R3.2 -->
### R3.2

The log must let a validator distinguish a subject that called a dependency correctly from one that did not, including unmatched calls.

<!-- req: R3.3 -->
### R3.3

The log must be scoped to the mock instance and reset when it starts, so a scenario never observes an earlier scenario's requests.

<!-- req: R4 -->
## R4 - Inputs, tools, and runtime dependencies

<!-- req: R4.1 -->
### R4.1

Mock inputs must include the fixture file path and the listen address, both supplied at start.

<!-- req: R4.2 -->
### R4.2

The mock profile must select the REST server tool family and the lifecycle control words it needs to serve and exit.

<!-- req: R4.3 -->
### R4.3

Mock dependencies must name agent-core for the mock endpoint binding and REST server runtime; the fixture content stays profile-owned or scenario-owned.

<!-- req: R4.4 -->
### R4.4

The mock must run as the same agent binary over a mounted profile, requiring no dedicated container image.

<!-- req: R4.5 -->
### R4.5

The mock must not acquire a dedicated CLI/runtime dependency unless it preserves the full srd039 scripted-response, request-log, validation, concurrency, lifecycle, and telemetry contract without requiring a separate image (GH-1386).

<!-- req: R5 -->
## R5 - Outputs and outcomes

<!-- req: R5.1 -->
### R5.1

Mock outputs must include the served responses and the readable request log.

<!-- req: R5.2 -->
### R5.2

Mock terminal outcomes must be Done or Failed.

<!-- field: non_goals -->
## Non-goals

- Define the fixture schema, matching rules, or request-log payload shape; those are srd039 in agent-core.
- Proxy to, record from, or replay a real upstream service.
- Assert anything; the mock serves and records, and the validator asserts.
- Model a client-side mock; the mock serves inbound requests only.
- Replace the shared agent-core binding with a partially compatible third-party mock process.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

The mock profile is described as one generic configured responder that serves through the agent-core mock binding, exposes health, stops on a lifecycle signal, invokes no model, defaults to loopback, and accepts a public listener only with explicit authority.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R1.4
- R1.5

<!-- ac: AC2 -->
### AC2

One mock profile serves any dependency because its fixture is selected by environment, scenarios own their fixtures, several mocks can run per scenario, malformed fixtures block serving, and Ollama-shaped model discovery plus ordered chat responses are representable as fixture data.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3
- R2.4
- R2.5

<!-- ac: AC3 -->
### AC3

The request log is readable on a declared route, is sufficient to verify a subject's outbound calls including unmatched ones, and is per-instance.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3

<!-- ac: AC4 -->
### AC4

Mock inputs, selected tool families, and runtime dependencies are named and separated from agent-core runtime capability; the mock needs no dedicated image or partially compatible CLI.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
- R4.3
- R4.4
- R4.5

<!-- ac: AC5 -->
### AC5

Mock outputs and terminal outcomes are checkable by mounted execution or conformance targets.

<!-- field: traces -->
#### Traces

- R5.1
- R5.2

<!-- field: diagrams | form: "tree" -->
## Diagrams

- requirements_trace:
  ```mermaid
  flowchart LR
    Problem["Problem: controlling dependency responses without live services"] --> Goals["Goals: Mock Functional Block"]
    Goals --> R1["R1: Behavior model"]
    Goals --> R2["R2: Fixture selection and scope"]
    Goals --> R3["R3: Request log"]
    Goals --> R4["R4: Inputs, tools, dependencies"]
    Goals --> R5["R5: Outputs and outcomes"]
    R1 --> AC1["AC1: Generic configured responder"]
    R2 --> AC2["AC2: One profile, any dependency"]
    R3 --> AC3["AC3: Readable per-instance request log"]
    R4 --> AC4["AC4: Inputs, tools, dependencies named"]
    R5 --> AC5["AC5: Outputs and outcomes checkable"]
  ```
