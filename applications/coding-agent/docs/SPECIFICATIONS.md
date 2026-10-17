---
id: coding-agent-example-specifications
title: Coding Agent Example Specification Index
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

This index governs the coding application's composition documents. Release
01.0 defines the local three-stage loop. Release 02.0 adds the
application-owned persistent serving contract while continuing to reference
the canonical executor (srd002), critic (srd003), and planner (srd004)
requirements in `applications/catalog`.

Library references are not forks. Packaging snapshots the configured checkout
compatible with the declared catalog `v0.*` release, records its
revision, and assembles the transitive profile closure before deployment. The
shared agent-core runtime image remains profile-free.

<!-- field: roadmap_summary | form: "typed table" -->
## Roadmap summary

| version | name | use_cases_done | use_cases_total | status |
| --- | --- | --- | --- | --- |
| "01.0" | Planner, executor, and critic coding loop | 1 | 1 | implemented |
| "02.0" | Persistent remote coding request | 1 | 1 | implemented |
| "03.0" | Role-specific deployment profile package | 1 | 1 | implemented |
| "04.0" | Core Helm topology | 1 | 1 | implemented |
| "05.0" | Helm values and package validation | 1 | 1 | implemented |
| "06.0" | Packaged-chart live Helm smoke | 1 | 1 | implemented |
| "07.0" | Canonical collector agent tracing | 1 | 1 | implemented |
| "08.0" | Declarative applier agent | 1 | 1 | implemented |
| "99.0" | Fleet observer adoption | 0 | 0 | planned |

<!-- field: foundation_document_index -->
## Foundation document index

| id | title | summary | path |
| --- | --- | --- | --- |
| [coding-agent-example](VISION.md) | Coding Agent Application Vision | Goals, boundaries, and success criteria for the deployable coding composition. | docs/VISION.md |
| [architecture-coding-agent-example](ARCHITECTURE.md) | Coding Agent Architecture | Coding-loop components, package boundary, deployment intent, and design decisions. | docs/ARCHITECTURE.md |
| [coding-agent-example-roadmap](road-map.md) | Coding Agent Example Road Map | Specification-first release plan for the coding loop. | docs/road-map.md |
| coding-agent-deployment-guide | Coding-Agent Deployment and Operations | Packaged chart prerequisites, configuration, verification, troubleshooting, security, rollback, and teardown. | docs/guides/deployment.md |

<!-- field: srd_index -->
## SRD index

| id | title | status | path |
| --- | --- | --- | --- |
| [srd001-serving-contract](specs/software-requirements/srd001-serving-contract.md) | Persistent Coding-Agent Serving Contract | implemented | docs/specs/software-requirements/srd001-serving-contract.md |
| [srd002-deployment-profile-package](specs/software-requirements/srd002-deployment-profile-package.md) | Role-Specific Deployment Profile Package | implemented | docs/specs/software-requirements/srd002-deployment-profile-package.md |
| [srd003-helm-topology](specs/software-requirements/srd003-helm-topology.md) | Core Coding-Agent Helm Topology | implemented | docs/specs/software-requirements/srd003-helm-topology.md |
| [srd004-helm-package-validation](specs/software-requirements/srd004-helm-package-validation.md) | Helm Values and Package Validation | implemented | docs/specs/software-requirements/srd004-helm-package-validation.md |
| [srd005-live-helm-smoke](specs/software-requirements/srd005-live-helm-smoke.md) | Packaged-Chart Live Helm Smoke | implemented | docs/specs/software-requirements/srd005-live-helm-smoke.md |
| [srd006-applier](specs/software-requirements/srd006-applier.md) | Declarative Applier Agent | implemented | docs/specs/software-requirements/srd006-applier.md |

<!-- field: library_srd_references -->
## Library SRD references

| id | owner | title | path | usage |
| --- | --- | --- | --- | --- |
| srd002-executor | applications/catalog | Executor Functional Block | ../../catalog/docs/specs/software-requirements/srd002-executor.md | The live executor stage and planner delegation use its coding-agent inputs, outputs, validation, and terminal outcomes. |
| srd003-critic | applications/catalog | Critic Functional Block | ../../catalog/docs/specs/software-requirements/srd003-critic.md | The critic-gate stage uses its suite, point, oracle, report, and terminal-outcome contract. |
| srd004-planner | applications/catalog | Planner Functional Block | ../../catalog/docs/specs/software-requirements/srd004-planner.md | The planner-delegation stage uses its requirement traversal, task materialization, child-profile execution, and validation contract. |

<!-- field: config_format_index | form: "empty list" -->
## Config format index

None

<!-- field: semantic_model_index | form: "empty list" -->
## Semantic model index

None

<!-- field: use_case_index -->
## Use case index

| id | title | release | status | test_suite | path |
| --- | --- | --- | --- | --- | --- |
| [rel01.0-uc001-coding-loop](specs/use-cases/rel01.0-uc001-coding-loop.md) | Planner Materializes, Executor Changes, and Critic Gates | 01.0 | implemented | [test-rel01.0-coding-loop](specs/test-suites/test-rel01.0-coding-loop.md) | docs/specs/use-cases/rel01.0-uc001-coding-loop.md |
| [rel02.0-uc001-remote-coding-request](specs/use-cases/rel02.0-uc001-remote-coding-request.md) | Persistent Planner Delegates to Executor and Critic | 02.0 | implemented | [test-rel02.0-serving-contract](specs/test-suites/test-rel02.0-serving-contract.md) | docs/specs/use-cases/rel02.0-uc001-remote-coding-request.md |
| [rel03.0-uc001-package-serving-profiles](specs/use-cases/rel03.0-uc001-package-serving-profiles.md) | Package Serving Profiles for Separate Containers | 03.0 | implemented | [test-rel03.0-deployment-profile-package](specs/test-suites/test-rel03.0-deployment-profile-package.md) | docs/specs/use-cases/rel03.0-uc001-package-serving-profiles.md |
| [rel04.0-uc001-render-helm-topology](specs/use-cases/rel04.0-uc001-render-helm-topology.md) | Render the Persistent Coding-Agent Topology | 04.0 | implemented | [test-rel04.0-helm-topology](specs/test-suites/test-rel04.0-helm-topology.md) | docs/specs/use-cases/rel04.0-uc001-render-helm-topology.md |
| [rel05.0-uc001-package-validated-chart](specs/use-cases/rel05.0-uc001-package-validated-chart.md) | Produce a Validated Self-Contained Chart | 05.0 | implemented | [test-rel05.0-helm-package-validation](specs/test-suites/test-rel05.0-helm-package-validation.md) | docs/specs/use-cases/rel05.0-uc001-package-validated-chart.md |
| [rel06.0-uc001-live-helm-smoke](specs/use-cases/rel06.0-uc001-live-helm-smoke.md) | Prove the Packaged Coding Topology on Kind | 06.0 | implemented | [test-rel06.0-live-helm-smoke](specs/test-suites/test-rel06.0-live-helm-smoke.md) | docs/specs/use-cases/rel06.0-uc001-live-helm-smoke.md |
| [rel07.0-uc001-collector-agent-tracing](specs/use-cases/rel07.0-uc001-collector-agent-tracing.md) | Trace the Packaged Topology Through the Canonical Collector Agent | 07.0 | implemented | [test-rel07.0-collector-agent-tracing](specs/test-suites/test-rel07.0-collector-agent-tracing.md) | docs/specs/use-cases/rel07.0-uc001-collector-agent-tracing.md |
| [rel08.0-uc001-declarative-applier](specs/use-cases/rel08.0-uc001-declarative-applier.md) | The Applier Agent Validates, Applies, and Verifies a Coding-Chart Values Change | 08.0 | implemented | [test-rel08.0-applier](specs/test-suites/test-rel08.0-applier.md) | docs/specs/use-cases/rel08.0-uc001-declarative-applier.md |

<!-- field: test_suite_index | form: "typed table" -->
## Test suite index

| id | title | release | test_case_count | traces | path |
| --- | --- | --- | --- | --- | --- |
| [test-rel01.0-coding-loop](specs/test-suites/test-rel01.0-coding-loop.md) | Test Suite for Release 01.0 Coding Loop | "01.0" | 3 | [[rel01.0-uc001-coding-loop](specs/use-cases/rel01.0-uc001-coding-loop.md)] | docs/specs/test-suites/test-rel01.0-coding-loop.md |
| [test-rel02.0-serving-contract](specs/test-suites/test-rel02.0-serving-contract.md) | Persistent Coding-Agent Serving Contract Suite | "02.0" | 2 | [[rel02.0-uc001-remote-coding-request](specs/use-cases/rel02.0-uc001-remote-coding-request.md)] | docs/specs/test-suites/test-rel02.0-serving-contract.md |
| [test-rel03.0-deployment-profile-package](specs/test-suites/test-rel03.0-deployment-profile-package.md) | Role-Specific Deployment Profile Package Suite | "03.0" | 4 | [[rel03.0-uc001-package-serving-profiles](specs/use-cases/rel03.0-uc001-package-serving-profiles.md)] | docs/specs/test-suites/test-rel03.0-deployment-profile-package.md |
| [test-rel04.0-helm-topology](specs/test-suites/test-rel04.0-helm-topology.md) | Core Coding-Agent Helm Topology Suite | "04.0" | 4 | [[rel04.0-uc001-render-helm-topology](specs/use-cases/rel04.0-uc001-render-helm-topology.md)] | docs/specs/test-suites/test-rel04.0-helm-topology.md |
| [test-rel05.0-helm-package-validation](specs/test-suites/test-rel05.0-helm-package-validation.md) | Helm Values and Package Validation Suite | "05.0" | 4 | [[rel05.0-uc001-package-validated-chart](specs/use-cases/rel05.0-uc001-package-validated-chart.md)] | docs/specs/test-suites/test-rel05.0-helm-package-validation.md |
| [test-rel06.0-live-helm-smoke](specs/test-suites/test-rel06.0-live-helm-smoke.md) | Packaged-Chart Live Helm Smoke Suite | "06.0" | 4 | [[rel06.0-uc001-live-helm-smoke](specs/use-cases/rel06.0-uc001-live-helm-smoke.md)] | docs/specs/test-suites/test-rel06.0-live-helm-smoke.md |
| [test-rel07.0-collector-agent-tracing](specs/test-suites/test-rel07.0-collector-agent-tracing.md) | Canonical Collector Agent Tracing Suite | "07.0" | 5 | [[rel07.0-uc001-collector-agent-tracing](specs/use-cases/rel07.0-uc001-collector-agent-tracing.md)] | docs/specs/test-suites/test-rel07.0-collector-agent-tracing.md |
| [test-rel08.0-applier](specs/test-suites/test-rel08.0-applier.md) | Test Suite for Release 08.0 Declarative Applier Agent | "08.0" | 5 | [[rel08.0-uc001-declarative-applier](specs/use-cases/rel08.0-uc001-declarative-applier.md)] | docs/specs/test-suites/test-rel08.0-applier.md |

<!-- field: coverage_gaps | form: "tree" -->
## Coverage gaps

1. - id: kubernetes-deployment
   - status: partial
   - description:
     The one-container-per-agent Helm topology is implemented. Strict values
     schema, packaged-chart validation, and bounded single-node kind evidence
     are implemented. Multi-node and production operational evidence remain
     outside this application.
