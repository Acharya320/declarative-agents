---
id: agent-architecture-specifications
title: Agent Architecture Application Specification Index
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

This index governs the standalone Agent Architecture composition.
Release 00.0 establishes the guided demo with one use case and test suite.
Release 01.0 adds trace collection through the canonical collector agent.
Releases 02.0 and 03.0 add the planned Kubernetes deployment: a chart running
the documentation-curator and the collector from mounted catalog-owned
closures, a values schema, and a bounded kind smoke (srd001-helm-deployment).
Release 04.0 adds the planned standalone declarative applier that actuates the
chart for an authorized operator or CI caller (srd002-applier). The
application references the canonical documentation-curator and collector
requirements and profile UX evidence in applications/catalog; it does not
fork or re-index them as application-owned requirements.

The dependency direction is application -> applications/catalog ->
agent-core. The Go present tool, catalog checkout, and agent-core checkout are
explicit development or runtime inputs. Checked-in Mage tests validate
command construction, portable root resolution, lifecycle configuration,
tracing wiring, documentation traceability, and composition-only statistics.
Live launch, browser, readiness, shutdown, and listener-release observations
remain dependency-gated.

The application extends shared srd001 and srd002 by reference. runnable_module
and managed_service apply; srd002 R1 through R3 apply to both the canonical
documentation-curator server launched by the application and the Go present
server. packaged, helm_managed, and kind_demo are declared and planned, and
cite srd002 R4 through R6; the application owns srd001-helm-deployment and
srd002-applier for their local topology, package, schema, smoke, and actuation
requirements.

<!-- field: roadmap_summary | form: "typed table" -->
## Roadmap summary

| version | name | use_cases_done | use_cases_total | status |
| --- | --- | --- | --- | --- |
| "00.0" | Guided Knowledge Manager presentation | 0 | 1 | partial |
| "01.0" | Demo trace collection | 1 | 1 | implemented |
| "02.0" | Packaged chart render and validation | 0 | 1 | planned |
| "03.0" | Live Helm smoke | 0 | 1 | planned |
| "04.0" | Declarative applier agent | 0 | 1 | planned |

<!-- field: foundation_document_index -->
## Foundation document index

| id | title | summary | path |
| --- | --- | --- | --- |
| agent-architecture | Agent Architecture Application Vision | Application purpose, ownership boundary, dependencies, and success criteria. | docs/VISION.md |
| architecture-agent-architecture | Agent Architecture Application Architecture | Presentation flow, components, interfaces, dependency direction, and decisions. | docs/ARCHITECTURE.md |
| agent-architecture-roadmap | Agent Architecture Application Road Map | Release 00.0 contract, implemented deterministic evidence, and dependency-gated live observations. | docs/road-map.md |

<!-- field: srd_index -->
## Srd index

| id | title | status | path |
| --- | --- | --- | --- |
| srd001-helm-deployment | Agent Architecture Helm Deployment | planned | docs/specs/software-requirements/srd001-helm-deployment.md |
| srd002-applier | Declarative Applier Agent | implemented | docs/specs/software-requirements/srd002-applier.md |

<!-- field: external_requirement_references | form: "tree" -->
## External requirement references

1. - id: srd001-application-module-contract
   - owner: applications
   - path: ../docs/specs/software-requirements/srd001-application-module-contract.md
   - requirements:
     1. R1
     1. R2
     1. R3
     1. R4
     1. R5
     1. R6
   - usage:
     Governs runnable-module identity, canonical references, local
     documentation, Mage participation, and composition-only accounting.
     This application supplies local paths and evidence without copying those
     requirements.
1. - id: srd002-managed-application-services
   - owner: applications
   - path: ../docs/specs/software-requirements/srd002-managed-application-services.md
   - requirements:
     1. R1
     1. R2
     1. R3
     1. R4
     1. R5
     1. R6
   - applies_to:
     1. documentation-curator server launched by this application
     1. Go present server
     1. planned chart, package, and kind deployment
   - usage:
     Governs lifecycle, control and operator verbs, telemetry, and
     configuration for both Release 00.0 managed services (R1 through R3). The
     planned deployment cites the package (R4), Helm (R5), and kind (R6) groups,
     supplying their evidence through srd001-helm-deployment.
1. - id: srd011-knowledge-manager
   - owner: applications/catalog
   - path: ../catalog/docs/specs/software-requirements/srd011-knowledge-manager.md
   - usage:
     Defines the canonical documentation-curator profile, catalog-owned UI and
     resources, agent-core runtime dependencies, and lifecycle outcomes.
1. - id: rel07.0-uc001-documentation-curator-profile-ux
   - owner: applications/catalog
   - path: ../catalog/docs/specs/use-cases/rel07.0-uc001-documentation-curator-profile-ux.md
   - usage:
     Proves the catalog profile's browser, action-route, and lifecycle behavior.
     Release 00.0 composes that behavior into a guided presentation.

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
| rel00.0-uc001-guided-agent-architecture | Operator Runs the Guided Agent Architecture | 00.0 | partial | test-rel00.0-guided-agent-architecture | docs/specs/use-cases/rel00.0-uc001-guided-agent-architecture.md |
| rel01.0-uc001-demo-trace-collection | Demo Run Collects and Serves Traces | 01.0 | implemented | test-rel01.0-demo-trace-collection | docs/specs/use-cases/rel01.0-uc001-demo-trace-collection.md |
| rel02.0-uc001-render-packaged-chart | Render and Validate the Packaged Agent Architecture Chart | 02.0 | planned | test-rel02.0-render-packaged-chart | docs/specs/use-cases/rel02.0-uc001-render-packaged-chart.md |
| rel03.0-uc001-live-helm-smoke | Prove the Packaged Agent Architecture Chart on Kind | 03.0 | planned | test-rel03.0-live-helm-smoke | docs/specs/use-cases/rel03.0-uc001-live-helm-smoke.md |
| rel04.0-uc001-declarative-applier | The Applier Agent Validates, Applies, and Verifies an Agent Architecture Values Change | 04.0 | planned | test-rel04.0-declarative-applier | docs/specs/use-cases/rel04.0-uc001-declarative-applier.md |

<!-- field: test_suite_index | form: "typed table" -->
## Test suite index

| id | title | release | status | test_case_count | traces | path |
| --- | --- | --- | --- | --- | --- | --- |
| test-rel00.0-guided-agent-architecture | Test Suite for Release 00.0 Guided Agent Architecture | "00.0" | partial | 5 | [rel00.0-uc001-guided-agent-architecture] | docs/specs/test-suites/test-rel00.0-guided-agent-architecture.md |
| test-rel01.0-demo-trace-collection | Test Suite for Release 01.0 Demo Trace Collection | "01.0" | implemented | 3 | [rel01.0-uc001-demo-trace-collection] | docs/specs/test-suites/test-rel01.0-demo-trace-collection.md |
| test-rel02.0-render-packaged-chart | Test Suite for Release 02.0 Packaged Chart Render and Validation | "02.0" | planned | 4 | [rel02.0-uc001-render-packaged-chart] | docs/specs/test-suites/test-rel02.0-render-packaged-chart.md |
| test-rel03.0-live-helm-smoke | Test Suite for Release 03.0 Live Helm Smoke | "03.0" | planned | 4 | [rel03.0-uc001-live-helm-smoke] | docs/specs/test-suites/test-rel03.0-live-helm-smoke.md |
| test-rel04.0-declarative-applier | Test Suite for Release 04.0 Declarative Applier Agent | "04.0" | planned | 5 | [rel04.0-uc001-declarative-applier] | docs/specs/test-suites/test-rel04.0-declarative-applier.md |

<!-- field: coverage_gaps | form: "tree" -->
## Coverage gaps

1. - id: live-browser-observation
   - status: dependency_gated
   - description:
     Deterministic tests validate commands, configuration, reciprocal traces,
     root dispatch, and accounting without starting servers. Live launch,
     browser rendering, role readiness, exit, and listener-release evidence
     requires available ports, a buildable agent-core checkout, and operator
     execution. Command construction is not live-service proof.
1. - id: kubernetes-deployment
   - status: implemented
   - description:
     The chart, profile-closure package, values schema, kind smoke, and applier
     are implemented through srd001-helm-deployment and srd002-applier use
     cases and test suites. Package, Helm, kind, and canonical-applier wrapper
     evidence is executable; live runs retain their documented host
     prerequisite gates.
