---
id: rel06.0-uc001-live-helm-smoke
title: Prove the Packaged Coding Topology on Kind
status: implemented
test_suite: test-rel06.0-live-helm-smoke
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A developer runs one bounded target. It verifies host prerequisites, acquires
a disposable kind cluster, builds and loads one commit-addressed canonical
profile-free agent image, installs the packaged chart, runs every agent and
the deterministic model boundary from that one image, drives one coding
request through all roles and the required build, v2 lint, and test sequence,
verifies the shared result and trace, and cleans up.

<!-- field: actor -->
## Actor

A developer or CI worker validating the deployable coding application.

<!-- field: trigger -->
## Trigger

Run mage integration:helmSmoke.

<!-- field: preconditions -->
## Preconditions

- Docker, kind, Helm, kubectl, and the required local images are available.
- The agent-core and applications/catalog checkouts are present.

<!-- field: flow -->
## Flow

- F1: Check prerequisites without mutating a cluster; skip with one explicit reason when absent.
- F2: Acquire a pinned kind cluster and verify Docker and Kubernetes readiness.
- F3: Build and load one commit-addressed canonical profile-free agent image, stage the executor's Go 1.26 and golangci-lint v2.12.2 donors, the critic's separate Go-only oracle donor, and the canonical model-mock profile and fixture, create the RWX fixture, and package the chart.
- F4: Install the .tgz and wait for model, collector agent, planner, executor, and critic readiness.
- F5: Seed the shared workspace and verify every role lifecycle endpoint.
- F6: Submit one traced planner request and verify executor mutation, build, v2 lint, Go tests, critic acceptance, and the connected collector agent trace.
- F7: On failure classify infrastructure versus semantics, emit bounded diagnostics, and clean up; on success clean up.

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

All three deployed role health endpoints are truthful and reachable.

<!-- field: traces -->
#### Traces

- srd005-live-helm-smoke R3.1
- srd005-live-helm-smoke AC1

<!-- sc: S2 -->
### S2

The shared workspace and critic verdict prove remote role composition.

<!-- field: traces -->
#### Traces

- srd005-live-helm-smoke R3.2
- srd005-live-helm-smoke AC1

<!-- sc: S3 -->
### S3

The collector agent contains all three role service names under the submitted trace ID.

<!-- field: traces -->
#### Traces

- srd005-live-helm-smoke R3.3
- srd005-live-helm-smoke AC1

<!-- sc: S4 -->
### S4

Skip, classification, diagnostics, and cleanup behavior are bounded and explicit.

<!-- field: traces -->
#### Traces

- srd005-live-helm-smoke R1.1
- srd005-live-helm-smoke R1.2
- srd005-live-helm-smoke R1.3
- srd005-live-helm-smoke R4.1
- srd005-live-helm-smoke R4.2
- srd005-live-helm-smoke R4.3
- srd005-live-helm-smoke AC2
- srd005-live-helm-smoke AC3

<!-- field: out_of_scope -->
## Out of scope

- Production model performance and GPU sizing.
- Multi-node RWX and high availability.
- Public ingress and authentication.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd005-live-helm-smoke R3.1 | S1 traces this requirement in the existing success criteria. |
| T2 | srd005-live-helm-smoke R3.2 | S2 traces this requirement in the existing success criteria. |
| T3 | srd005-live-helm-smoke R3.3 | S3 traces this requirement in the existing success criteria. |
| T4 | srd005-live-helm-smoke R1.1, R1.2, R1.3 | S4 traces this requirement in the existing success criteria. |
| T5 | srd005-live-helm-smoke R4.1, R4.2, R4.3 | S4 traces this requirement in the existing success criteria. |
