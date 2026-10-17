---
id: coding-agent-example-roadmap
title: Coding Agent Example Road Map
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

The coding application starts with a specified end-to-end loop in release
01.0. The release composes the canonical planner, executor, and critic library
agents. Documentation and the test-suite contract land before integration
implementation, reference resolution, packaging, and deployment.

<!-- release: 01.0 | status: implemented -->
## 01.0 - Planner, executor, and critic coding loop

An SRD and requirement graph enter the real planner. The planner
materializes a task and delegates to the real executor. The executor
changes an isolated workspace and returns successful declared validation.
The critic evaluates the resulting change and gates the terminal
application outcome. Three ordered integration stages prove those
boundaries without a fake agent binary.

<!-- field: library_srds -->
### Library SRDs

- srd002-executor
- srd003-critic
- srd004-planner

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| rel01.0-uc001-coding-loop | A real planner materializes an SRD task, delegates it to a real executor, and submits the validated change to a critic gate. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| test-rel01.0-coding-loop | Three stages prove a live executor run, real planner delegation, and the critic verdict with explicit inputs and outputs. | implemented |

<!-- release: 02.0 | status: implemented -->
## 02.0 - Persistent remote coding request

Application-owned serving profiles keep planner, executor, and critic
alive behind truthful lifecycle endpoints. A planner request reuses
canonical planning words, crosses declared executor and critic REST
boundaries over one trusted shared workspace, and retains one W3C trace.

<!-- field: application_srds -->
### Application SRDs

- srd001-serving-contract

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| rel02.0-uc001-remote-coding-request | A persistent planner delegates one coding request to remote executor and critic services. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| test-rel02.0-serving-contract | Real health, canonical critic, remote coding, and connected trace evidence. | implemented |

<!-- release: 03.0 | status: implemented -->
## 03.0 - Role-specific deployment profile package

The #847 resolver packages independent planner, executor, and critic
profile roots from one canonical source namespace. Deterministic manifests
retain provenance and partition each role for bounded ConfigMaps.

<!-- field: application_srds -->
### Application SRDs

- srd002-deployment-profile-package

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| rel03.0-uc001-package-serving-profiles | A release builder emits one exact, mountable profile closure per serving role. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| test-rel03.0-deployment-profile-package | Exact closures, reproducibility, safety failures, ConfigMap bounds, and real-binary boot validation. | implemented |

<!-- release: 04.0 | status: implemented -->
## 04.0 - Core Helm topology

A package-driven chart renders persistent planner, executor, and critic
workloads with scoped profile ConfigMaps, shared workspace storage,
truthful health probes, fixed role Services, model endpoint selection,
and collector-to-Jaeger tracing. Release 07.0 replaces the trace path
with the canonical collector agent.

<!-- field: application_srds -->
### Application SRDs

- srd003-helm-topology

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| rel04.0-uc001-render-helm-topology | An operator prepares role shards and renders the core Kubernetes topology. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| test-rel04.0-helm-topology | Deterministic topology, checksum, CI, model-tier, and trace-wiring renders. | implemented |

<!-- release: 05.0 | status: implemented -->
## 05.0 - Helm values and package validation

Strict values and semantic guards reject unsafe topology. Prepared role
packages are checksum-validated, and a self-contained chart archive is
inventory-checked, linted, and rendered independently.

<!-- field: application_srds -->
### Application SRDs

- srd004-helm-package-validation

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| rel05.0-uc001-package-validated-chart | A release builder produces one validated chart archive with every role asset. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| test-rel05.0-helm-package-validation | Schema matrix, semantic invariants, tamper rejection, and independent archive rendering. | implemented |

<!-- release: 06.0 | status: implemented -->
## 06.0 - Packaged-chart live Helm smoke

A bounded kind target installs only the packaged chart and proves truthful
role health, remote workspace mutation, critic acceptance, and one
planner/executor/critic trace with classified failures and cleanup.

<!-- field: application_srds -->
### Application SRDs

- srd005-live-helm-smoke

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| rel06.0-uc001-live-helm-smoke | A developer proves the complete packaged topology on a disposable kind cluster. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| test-rel06.0-live-helm-smoke | Skip, classification, diagnostic, fixture, cleanup, and live topology evidence. | implemented |

<!-- release: 07.0 | status: implemented -->
## 07.0 - Canonical collector agent tracing

The chart replaces the contrib collector and the Jaeger all-in-one with
the canonical catalog collector agent (GH-1152, GH-1156). Roles export
OTLP to the profile-backed collector, which spools trace evidence and
serves the bounded query surface; the utility-image wait chain and every
Jaeger workload leave the chart, and the compatibility pin moves to
applications/catalog/v0.20260730.0.

<!-- field: application_srds -->
### Application SRDs

- srd003-helm-topology
- srd004-helm-package-validation
- srd005-live-helm-smoke

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| rel07.0-uc001-collector-agent-tracing | An operator renders and proves the packaged topology with the canonical collector agent and no Jaeger workload. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| test-rel07.0-collector-agent-tracing | Topology, schema, catalog-root discovery, and live smoke evidence with collector agent rendering and spool query trace proof. | implemented |

<!-- release: 08.0 | status: implemented -->
## 08.0 - Declarative applier agent

A standalone declarative applier actuates the coding chart. An authorized
operator or CI caller submits a decided values-plane patch to the applier's
apply endpoint; the applier validates it against the chart's
values.schema.json by a helm dry-run, applies it with helm upgrade
--reuse-values returning without waiting, verifies the planner, executor, and
critic Deployments via kubectl rollout status, rolls back on a verify stall,
and reports an apply-command failure as failed. The applier alone holds helm
and kubectl bound as exec tools, the
serving roles carry none, and its apply surface is network-gated off the
serving and collector pods. The coding application runs no control plane, so
there is no creator, orchestrator, or panel; the applier runs the agent-core
image with helm and kubectl delivered by a pinned CLI donor (GH-2222). The
profile, chart resources, fake-boundary tests, and live integration target
are implemented.

<!-- field: application_srds -->
### Application SRDs

- srd006-applier

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| rel08.0-uc001-declarative-applier | An operator or CI caller submits a decided values patch that the applier validates, applies, verifies, and rolls back on failure. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| test-rel08.0-applier | Recording-fake terminal coverage, real chart-schema rejection, applier-only CLI authority, broken-read classification, and a live on-cluster rollout and rollback. | implemented |

<!-- release: 99.0 | status: planned -->
## 99.0 - Fleet observer adoption

Adopt the catalog-owned observer family (catalog srd021-observer,
monitor/fleet-observer) as a coding-agent workload pointed at this
application's release label selector and namespace, giving the coding mesh
a fleet view without forking the profile. Unscheduled; specification of the
catalog family lands first (GH-1247). No new environment variables.

<!-- field: planned_work -->
### Planned work

| id | summary | status |
| --- | --- | --- |
| gh-1247 | Specify the cross-application observer family this application adopts. | planned |
