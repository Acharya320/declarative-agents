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

<!-- field: roadmap_summary -->
## Roadmap summary

```yaml
- version: "01.0"
  name: Planner, executor, and critic coding loop
  use_cases_done: 1
  use_cases_total: 1
  status: implemented
- version: "02.0"
  name: Persistent remote coding request
  use_cases_done: 1
  use_cases_total: 1
  status: implemented
- version: "03.0"
  name: Role-specific deployment profile package
  use_cases_done: 1
  use_cases_total: 1
  status: implemented
- version: "04.0"
  name: Core Helm topology
  use_cases_done: 1
  use_cases_total: 1
  status: implemented
- version: "05.0"
  name: Helm values and package validation
  use_cases_done: 1
  use_cases_total: 1
  status: implemented
- version: "06.0"
  name: Packaged-chart live Helm smoke
  use_cases_done: 1
  use_cases_total: 1
  status: implemented
- version: "07.0"
  name: Canonical collector agent tracing
  use_cases_done: 1
  use_cases_total: 1
  status: implemented
- version: "08.0"
  name: Declarative applier agent
  use_cases_done: 1
  use_cases_total: 1
  status: implemented
- version: "99.0"
  name: Fleet observer adoption
  use_cases_done: 0
  use_cases_total: 0
  status: planned
```

<!-- field: foundation_document_index -->
## Foundation document index

| id | title | summary | path |
| --- | --- | --- | --- |
| coding-agent-example | Coding Agent Application Vision | Goals, boundaries, and success criteria for the deployable coding composition. | docs/VISION.yaml |
| architecture-coding-agent-example | Coding Agent Architecture | Coding-loop components, package boundary, deployment intent, and design decisions. | docs/ARCHITECTURE.yaml |
| coding-agent-example-roadmap | Coding Agent Example Road Map | Specification-first release plan for the coding loop. | docs/road-map.yaml |
| coding-agent-deployment-guide | Coding-Agent Deployment and Operations | Packaged chart prerequisites, configuration, verification, troubleshooting, security, rollback, and teardown. | docs/guides/deployment.md |

<!-- field: srd_index -->
## Srd index

| id | title | status | path |
| --- | --- | --- | --- |
| srd001-serving-contract | Persistent Coding-Agent Serving Contract | implemented | docs/specs/software-requirements/srd001-serving-contract.yaml |
| srd002-deployment-profile-package | Role-Specific Deployment Profile Package | implemented | docs/specs/software-requirements/srd002-deployment-profile-package.yaml |
| srd003-helm-topology | Core Coding-Agent Helm Topology | implemented | docs/specs/software-requirements/srd003-helm-topology.yaml |
| srd004-helm-package-validation | Helm Values and Package Validation | implemented | docs/specs/software-requirements/srd004-helm-package-validation.yaml |
| srd005-live-helm-smoke | Packaged-Chart Live Helm Smoke | implemented | docs/specs/software-requirements/srd005-live-helm-smoke.yaml |
| srd006-applier | Declarative Applier Agent | implemented | docs/specs/software-requirements/srd006-applier.yaml |

<!-- field: library_srd_references -->
## Library srd references

| id | owner | title | path | usage |
| --- | --- | --- | --- | --- |
| srd002-executor | applications/catalog | Executor Functional Block | ../../catalog/docs/specs/software-requirements/srd002-executor.yaml | The live executor stage and planner delegation use its coding-agent inputs, outputs, validation, and terminal outcomes. |
| srd003-critic | applications/catalog | Critic Functional Block | ../../catalog/docs/specs/software-requirements/srd003-critic.yaml | The critic-gate stage uses its suite, point, oracle, report, and terminal-outcome contract. |
| srd004-planner | applications/catalog | Planner Functional Block | ../../catalog/docs/specs/software-requirements/srd004-planner.yaml | The planner-delegation stage uses its requirement traversal, task materialization, child-profile execution, and validation contract. |

<!-- field: config_format_index -->
## Config format index

```yaml
[]
```

<!-- field: semantic_model_index -->
## Semantic model index

```yaml
[]
```

<!-- field: use_case_index -->
## Use case index

| id | title | release | status | test_suite | path |
| --- | --- | --- | --- | --- | --- |
| rel01.0-uc001-coding-loop | Planner Materializes, Executor Changes, and Critic Gates | 01.0 | implemented | test-rel01.0-coding-loop | docs/specs/use-cases/rel01.0-uc001-coding-loop.yaml |
| rel02.0-uc001-remote-coding-request | Persistent Planner Delegates to Executor and Critic | 02.0 | implemented | test-rel02.0-serving-contract | docs/specs/use-cases/rel02.0-uc001-remote-coding-request.yaml |
| rel03.0-uc001-package-serving-profiles | Package Serving Profiles for Separate Containers | 03.0 | implemented | test-rel03.0-deployment-profile-package | docs/specs/use-cases/rel03.0-uc001-package-serving-profiles.yaml |
| rel04.0-uc001-render-helm-topology | Render the Persistent Coding-Agent Topology | 04.0 | implemented | test-rel04.0-helm-topology | docs/specs/use-cases/rel04.0-uc001-render-helm-topology.yaml |
| rel05.0-uc001-package-validated-chart | Produce a Validated Self-Contained Chart | 05.0 | implemented | test-rel05.0-helm-package-validation | docs/specs/use-cases/rel05.0-uc001-package-validated-chart.yaml |
| rel06.0-uc001-live-helm-smoke | Prove the Packaged Coding Topology on Kind | 06.0 | implemented | test-rel06.0-live-helm-smoke | docs/specs/use-cases/rel06.0-uc001-live-helm-smoke.yaml |
| rel07.0-uc001-collector-agent-tracing | Trace the Packaged Topology Through the Canonical Collector Agent | 07.0 | implemented | test-rel07.0-collector-agent-tracing | docs/specs/use-cases/rel07.0-uc001-collector-agent-tracing.yaml |
| rel08.0-uc001-declarative-applier | The Applier Agent Validates, Applies, and Verifies a Coding-Chart Values Change | 08.0 | implemented | test-rel08.0-applier | docs/specs/use-cases/rel08.0-uc001-declarative-applier.yaml |

<!-- field: test_suite_index -->
## Test suite index

```yaml
- id: test-rel01.0-coding-loop
  title: Test Suite for Release 01.0 Coding Loop
  release: "01.0"
  test_case_count: 3
  traces: [rel01.0-uc001-coding-loop]
  path: docs/specs/test-suites/test-rel01.0-coding-loop.yaml
- id: test-rel02.0-serving-contract
  title: Persistent Coding-Agent Serving Contract Suite
  release: "02.0"
  test_case_count: 2
  traces: [rel02.0-uc001-remote-coding-request]
  path: docs/specs/test-suites/test-rel02.0-serving-contract.yaml
- id: test-rel03.0-deployment-profile-package
  title: Role-Specific Deployment Profile Package Suite
  release: "03.0"
  test_case_count: 4
  traces: [rel03.0-uc001-package-serving-profiles]
  path: docs/specs/test-suites/test-rel03.0-deployment-profile-package.yaml
- id: test-rel04.0-helm-topology
  title: Core Coding-Agent Helm Topology Suite
  release: "04.0"
  test_case_count: 4
  traces: [rel04.0-uc001-render-helm-topology]
  path: docs/specs/test-suites/test-rel04.0-helm-topology.yaml
- id: test-rel05.0-helm-package-validation
  title: Helm Values and Package Validation Suite
  release: "05.0"
  test_case_count: 4
  traces: [rel05.0-uc001-package-validated-chart]
  path: docs/specs/test-suites/test-rel05.0-helm-package-validation.yaml
- id: test-rel06.0-live-helm-smoke
  title: Packaged-Chart Live Helm Smoke Suite
  release: "06.0"
  test_case_count: 4
  traces: [rel06.0-uc001-live-helm-smoke]
  path: docs/specs/test-suites/test-rel06.0-live-helm-smoke.yaml
- id: test-rel07.0-collector-agent-tracing
  title: Canonical Collector Agent Tracing Suite
  release: "07.0"
  test_case_count: 5
  traces: [rel07.0-uc001-collector-agent-tracing]
  path: docs/specs/test-suites/test-rel07.0-collector-agent-tracing.yaml
- id: test-rel08.0-applier
  title: Test Suite for Release 08.0 Declarative Applier Agent
  release: "08.0"
  test_case_count: 5
  traces: [rel08.0-uc001-declarative-applier]
  path: docs/specs/test-suites/test-rel08.0-applier.yaml
```

<!-- field: coverage_gaps -->
## Coverage gaps

```yaml
- id: kubernetes-deployment
  status: partial
  description: |
    The one-container-per-agent Helm topology is implemented. Strict values
    schema, packaged-chart validation, and bounded single-node kind evidence
    are implemented. Multi-node and production operational evidence remain
    outside this application.
```
