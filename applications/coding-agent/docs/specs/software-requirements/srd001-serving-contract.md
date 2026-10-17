---
id: srd001-serving-contract
title: Persistent Coding-Agent Serving Contract
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

The canonical planner, executor, and critic are finite batch machines. A
Kubernetes Deployment needs a persistent process with truthful lifecycle
health, while the application needs network boundaries between roles without
copying reusable library behavior or allowing requests to choose profiles,
workspace roots, endpoints, or credentials.

<!-- req: R1 -->
## R1 - Persistent role hosts

<!-- req: R1.1 -->
### R1.1

The application shall own one persistent serving profile for planner, executor, and critic.

<!-- req: R1.2 -->
### R1.2

Each host shall launch a declared role request server and lifecycle server, remain in a non-terminal awaiting-control state, and stop through the declared lifecycle exit path.

<!-- req: R1.3 -->
### R1.3

Each lifecycle server shall expose GET /api/lifecycle/health through the agent-core health binding; no sidecar, sleep process, or unmodeled idle loop shall report health.

<!-- req: R2 -->
## R2 - Canonical role composition

<!-- req: R2.1 -->
### R2.1

The executor endpoint shall run the canonical executor profile as a profile-backed machine_request.

<!-- req: R2.2 -->
### R2.2

The critic endpoint shall run the canonical changed-workspace critic profile as a profile-backed machine_request.

<!-- req: R2.3 -->
### R2.3

The planner request machine shall reuse canonical planner LLM configuration, planner builtin tools, and plan parsing while owning only application sequencing.

<!-- req: R2.4 -->
### R2.4

Existing canonical batch profiles and local coding-loop integrations shall remain unchanged.

<!-- req: R3 -->
## R3 - Remote application sequence

<!-- req: R3.1 -->
### R3.1

A planner coding request shall invoke the executor through a declared REST client and invoke the critic through a second declared REST client only after executor success.

<!-- req: R3.2 -->
### R3.2

Executor failure, critic rejection, and infrastructure failure shall map to distinct declared terminal outcomes and HTTP responses.

<!-- req: R3.3 -->
### R3.3

Runtime input shall provide a task and opaque workspace identity but shall not choose a filesystem root, profile, machine, URL, host, or credential.

<!-- req: R4 -->
## R4 - Shared workspace and trace

<!-- req: R4.1 -->
### R4.1

All three processes shall bind the same operator-supplied runtime directory, suitable for one Kubernetes shared workspace volume.

<!-- req: R4.2 -->
### R4.2

The workspace identity shall be carried as correlation metadata while the trusted process directory remains the filesystem authority.

<!-- req: R4.3 -->
### R4.3

The planner inbound machine_request and both outbound REST operations shall use agent-core W3C trace propagation so all role spans share one trace identifier.

<!-- field: goals -->
## Goals

- G1: Persistent role hosts.
- G2: Canonical role composition.
- G3: Remote application sequence.
- G4: Shared workspace and trace.

<!-- field: non_goals -->
## Non-goals

- Helm resources, profile ConfigMap sizing, and closure sharding.
- Per-request filesystem roots or dynamically supplied transport authority.
- Replacing the local canonical batch integrations.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

All three profiles boot with the real agent binary and return HTTP 200 from their declared lifecycle health endpoints.

<!-- field: traces -->
#### Traces

- [R1.1](#r11)
- [R1.2](#r12)
- [R1.3](#r13)

<!-- ac: AC2 -->
### AC2

A localhost planner request changes the shared workspace through the remote canonical executor and receives an accepting verdict from the remote canonical critic.

<!-- field: traces -->
#### Traces

- [R2.1](#r21)
- [R2.2](#r22)
- [R2.3](#r23)
- [R3.1](#r31)
- [R3.2](#r32)
- [R3.3](#r33)
- [R4.1](#r41)
- [R4.2](#r42)

<!-- ac: AC3 -->
### AC3

Planner, executor, and critic trace output for the remote request contains the caller's W3C trace identifier.

<!-- field: traces -->
#### Traces

- [R4.3](#r43)

<!-- ac: AC4 -->
### AC4

The existing local batch coding-loop tests and integrations remain valid.

<!-- field: traces -->
#### Traces

- [R2.4](#r24)
