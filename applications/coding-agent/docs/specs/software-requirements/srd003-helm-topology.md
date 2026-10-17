---
id: srd003-helm-topology
title: Core Coding-Agent Helm Topology
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

The packaged serving roles need a Kubernetes topology that preserves their
mounted-profile, shared-workspace, fixed-service, lifecycle, and tracing
contracts without baking application programs into the agent-core image or
duplicating #875 closure membership in Helm values.

<!-- req: R1 -->
## R1 - Real role workloads

<!-- req: R1.1 -->
### R1.1

The chart shall render one Deployment and ClusterIP Service for planner, executor, and critic.

<!-- req: R1.2 -->
### R1.2

Every role shall use the same canonical, profile-free agent-core image and select the serving profile named by its prepared role manifest.

<!-- req: R1.3 -->
### R1.3

Role ConfigMaps and projected volumes shall derive file and partition membership only from prepared

<!-- req: R1.4 -->
### R1.4

A checksum of only the role's reachable profile files shall annotate its pod template.

<!-- req: R1.5 -->
### R1.5

The executor shall receive Go 1.26 and golangci-lint v2.12.2 from separate digest-pinned upstream init containers through one read-only tools volume so it can complete its mandatory build, lint, and test sequence without extending the agent image; the changed-workspace critic shall receive only Go through a separate read-only donor volume for its independent test oracle.

<!-- req: R2 -->
## R2 - Serving contract wiring

<!-- req: R2.1 -->
### R2.1

All roles shall mount one trusted shared PVC at the configured runtime directory.

<!-- req: R2.2 -->
### R2.2

Planner shall receive fixed executor and critic Service hosts and URLs matching ports 18210 and 18220.

<!-- req: R2.3 -->
### R2.3

Readiness and liveness shall call each existing /api/lifecycle/health control endpoint.

<!-- req: R2.4 -->
### R2.4

Role pods shall use non-root, no-privilege-escalation, dropped-capability, read-only-root, and disabled service-account-token defaults.

<!-- req: R3 -->
## R3 - Model and observability tiers

<!-- req: R3.1 -->
### R3.1

Planner and executor shall receive one values-selected OLLAMA_URL, from either an external endpoint or the optional chart-owned Ollama Service.

<!-- req: R3.2 -->
### R3.2

The optional Ollama tier shall provide persistent model storage, configured preload, and model-presence startup gating.

<!-- req: R3.3 -->
### R3.3

Every role shall use a distinct OTel service name and export OTLP to the optional collector.

<!-- req: R3.4 -->
### R3.4

The collector agent shall receive OTLP, spool traces, and serve a query surface.

<!-- req: R4 -->
## R4 - Reproducible preparation

<!-- req: R4.1 -->
### R4.1

mage helmPrepare shall regenerate

<!-- req: R4.2 -->
### R4.2

Default and small-footprint CI values shall lint and render without requiring Helm schema work from the next issue.

<!-- field: goals -->
## Goals

- G1: Real role workloads.
- G2: Serving contract wiring.
- G3: Model and observability tiers.
- G4: Reproducible preparation.

<!-- field: non_goals -->
## Non-goals

- values.schema.json and negative schema fixtures.
- Packaged-chart archive testing.
- Live kind installation and deployment documentation.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

Helm renders exactly three role Deployments with real profiles, the canonical agent-core image, role-minimal read-only tool donors (executor Go plus linter; critic Go only), shared workspace, fixed role Services, and truthful probes.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R1.4
- R1.5
- R2.1
- R2.2
- R2.3
- R2.4

<!-- ac: AC2 -->
### AC2

Editing one prepared role profile changes only that role's rollout checksum.

<!-- field: traces -->
#### Traces

- R1.3
- R1.4

<!-- ac: AC3 -->
### AC3

External and in-cluster model modes plus collector agent tracing render coherently.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3
- R3.4

<!-- ac: AC4 -->
### AC4

helm lint and default and CI helm template commands pass after mage helmPrepare.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
