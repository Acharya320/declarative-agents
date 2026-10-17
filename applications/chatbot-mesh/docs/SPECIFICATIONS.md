---
id: chatbot-mesh-example-specifications
title: Chatbot Mesh Example Specification Index
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

Specification index for the standalone chatbot mesh example. Eight SRDs define the
program: three data-plane -- the RAG server agent (srd001), the chatbot agent
(srd002), and the Helm deployment model (srd003) -- two control-plane -- the
provisioning-workflow-orchestrator (srd004) and the creator (srd005) -- the deployment-plane applier
(srd006), the collector agent (srd007), and the fleet observer agent (srd008).
Ten use cases exercise them across nine releases (rel00.0 through rel08.0), each
with a formal test suite.

The application runs on the agent-core runtime and cites agent-core platform SRDs
(srd021, srd028, srd029, srd030, srd033, srd034, srd035, srd036, and srd016
traceparent) as external references rather than indexing them here. The control
plane (provisioning-workflow-orchestrator, creator, and the provisioning surface) is release rel05.0; the
applier (srd006) that replaces the Go provisioning surface with a declarative
agent is release rel06.0.

The corpus is validated by the specification critic (mage audit): it loads and validates with no
errors. Advisory findings remain visible, including tool-undo-mismatch warnings;
the gate does not report a warning-free corpus. The traceability chain is complete:
every use case touchpoint builds an SRD edge (the object-format {id, target, reason} touchpoints
parse since GH-448), every success criterion traces an acceptance criterion
alongside its requirement groups, and every acceptance criterion is covered by a
test case. Every SRD is indexed and placed in ARCHITECTURE, every use case targets
its SRDs through touchpoints, and every test suite traces its use case.

<!-- field: role_alignment_index | form: "tree" -->
## Role alignment index

- authority: ../../docs/specs/semantic-models/agent-role-realizations.md#bindings.chatbot_mesh
- architecture: docs/ARCHITECTURE.md#role_realization_alignment
- rule: Actors assume roles; profile, wrapper, Helm workload, and values-selected instance names are realization identities, not canonical role definitions.
- decisions:
  1. - current_actor: chatbot
     - canonical_role_sub_role: Planner / Intent Handler; Executor / Workflow Orchestrator; Knowledge Manager / Recommendation Generator
     - retained_or_proposed_name: chatbot
     - evidence: srd002
     - migration_impact: none
  1. - current_actor: rag-server
     - canonical_role_sub_role: Knowledge Manager / Recommendation Generator
     - retained_or_proposed_name: rag-server
     - evidence: srd001
     - migration_impact: none
  1. - current_actor: corpus-ingest
     - canonical_role_sub_role: inherited Knowledge Manager / Ontology Manager
     - retained_or_proposed_name: corpus-ingest
     - evidence: srd001 ingest boundary and wrapper profile
     - migration_impact: none
  1. - current_actor: provisioning-workflow-orchestrator
     - canonical_role_sub_role: Executor / Workflow Orchestrator; not Coordinator
     - retained_or_proposed_name: provisioning-workflow-orchestrator
     - evidence: srd004
     - migration_impact: GH-1077
  1. - current_actor: creator
     - canonical_role_sub_role: Creator / Instance Manager; Agent Owner / Agent Lifecycle Manager
     - retained_or_proposed_name: creator
     - evidence: srd005
     - migration_impact: none
  1. - current_actor: applier
     - canonical_role_sub_role: Executor / Actuation Agent and Change Manager
     - retained_or_proposed_name: applier
     - evidence: srd006
     - migration_impact: documentation-only
  1. - current_actor: collector
     - canonical_role_sub_role: Monitor / Telemetry Collector
     - retained_or_proposed_name: collector
     - evidence: srd007
     - migration_impact: none
  1. - current_actor: observer
     - canonical_role_sub_role: Monitor / Fleet Observer
     - retained_or_proposed_name: observer
     - evidence: srd008
     - migration_impact: none
- deployment_evidence:
  1. helm/templates/chatbot.yaml
  1. helm/templates/rag-units.yaml
  1. helm/templates/provisioning-workflow-orchestrator.yaml
  1. helm/templates/creator.yaml
  1. helm/templates/applier.yaml
  1. helm/templates/collector.yaml
  1. helm/templates/observer.yaml

<!-- field: roadmap_summary | form: "typed table" -->
## Roadmap summary

| version | name | use_cases_done | use_cases_total | status |
| --- | --- | --- | --- | --- |
| "00.0" | RAG server and single-RAG chatbot turn | 2 | 2 | implemented |
| "01.0" | Tier selection and multi-RAG fan-out | 1 | 1 | implemented |
| "02.0" | Observability | 1 | 1 | implemented |
| "03.0" | Mesh deployment and tiered swap | 1 | 1 | implemented |
| "04.0" | In-cluster LLM tier | 1 | 1 | implemented |
| "05.0" | Control plane | 0 | 1 | partial |
| "06.0" | Declarative applier agent | 1 | 1 | implemented |
| "07.0" | Declarative collector agent | 0 | 1 | partial |
| "08.0" | Fleet observer agent | 0 | 1 | planned |
| "09.0" | Canonical collector and Jaeger retirement | 1 | 1 | done |

<!-- field: srd_index -->
## SRD index

| id | title | summary | path |
| --- | --- | --- | --- |
| srd001-rag-server-agent | RAG Server Agent Functional Block | Persistent vector-in retrieval service behind a machine_request query endpoint. | docs/specs/software-requirements/srd001-rag-server-agent.md |
| srd002-chatbot-agent | Chatbot Agent Functional Block | Browser-facing agent whose request-scoped turn selects a model tier, fans out to RAG servers, composes, and answers. | docs/specs/software-requirements/srd002-chatbot-agent.md |
| srd003-chatbot-deployment | Chatbot Mesh Deployment Model | One Helm chart with values-driven RAG pairs, co-generated client config, tiered hot swap, and an in-cluster LLM tier. | docs/specs/software-requirements/srd003-chatbot-deployment.md |
| srd004-provisioning-workflow-orchestrator | Provisioning Workflow Orchestrator Agent Functional Block | The current provisioning-workflow-orchestrator actor assumes Executor / Workflow Orchestrator, binds a provisioning intent, and sequences ingestion and reconfiguration through the creator; GH-1077 completed its semantic rename. | docs/specs/software-requirements/srd004-provisioning-workflow-orchestrator.md |
| srd005-creator | Creator Agent Functional Block | Actor assuming Creator / Instance Manager and Agent Owner / Agent Lifecycle Manager that realizes decided operations and alone holds deployment-API authority. | docs/specs/software-requirements/srd005-creator.md |
| srd006-applier | Applier Agent Functional Block | Executor / Actuation Agent realization that validates a values patch, applies and verifies a rollout, and rolls back, binding helm and kubectl as exec tools. | docs/specs/software-requirements/srd006-applier.md |
| srd007-collector-agent | Collector Agent Functional Block | Monitor / Telemetry Collector realization composing bounded OTLP receive, deterministic spool, unchanged relay, lifecycle surfaces, and self-loop prevention while contrib retains metrics. | docs/specs/software-requirements/srd007-collector-agent.md |
| srd008-observer-agent | Observer Agent Functional Block | Monitor / Fleet Observer realization that discovers mesh agents and the deployment topology from the Kubernetes API, polls each agent's monitor surface and per-pod resource usage (metrics.k8s.io), and serves a fleet-level view. | docs/specs/software-requirements/srd008-observer-agent.md |

<!-- field: use_case_index -->
## Use case index

| id | title | release | status | test_suite | path |
| --- | --- | --- | --- | --- | --- |
| rel00.0-uc001-rag-server-query | Persistent RAG Server Answers a Vector-In Query | 00.0 | implemented | test-rel00.0-chatbot-mesh | docs/specs/use-cases/rel00.0-uc001-rag-server-query.md |
| rel00.0-uc002-chatbot-turn | Chatbot Turn Embeds, Retrieves From One RAG, and Answers | 00.0 | implemented | test-rel00.0-chatbot-mesh | docs/specs/use-cases/rel00.0-uc002-chatbot-turn.md |
| rel01.0-uc001-routed-fanout | Chatbot Turn Selects Sources and a Model Tier Before Trusted RAG Fan-Out | 01.0 | implemented | test-rel01.0-routing-fanout | docs/specs/use-cases/rel01.0-uc001-routed-fanout.md |
| rel02.0-uc001-observability | Observability Panel Shows Per-Agent State and a Cross-Agent Trace Waterfall | 02.0 | implemented | test-rel02.0-observability | docs/specs/use-cases/rel02.0-uc001-observability.md |
| rel03.0-uc001-mesh-deployment | Helm Chart Deploys the Mesh and Hot-Swaps a RAG Source | 03.0 | implemented | test-rel03.0-mesh-deployment | docs/specs/use-cases/rel03.0-uc001-mesh-deployment.md |
| rel04.0-uc001-in-cluster-llm-tier | A Fresh-Cluster Install Answers a Turn With No External LLM Dependency | 04.0 | implemented | test-rel04.0-llm-tier | docs/specs/use-cases/rel04.0-uc001-in-cluster-llm-tier.md |
| rel05.0-uc001-control-plane-provisioning | A Provisioning Intent Flows Chatbot to Provisioning Workflow Orchestrator to Creator and Serves | 05.0 | partial | test-rel05.0-control-plane | docs/specs/use-cases/rel05.0-uc001-control-plane-provisioning.md |
| rel06.0-uc001-declarative-applier | The Applier Agent Realizes the Deployment API and Rolls Out a Values Change | 06.0 | implemented | test-rel06.0-applier | docs/specs/use-cases/rel06.0-uc001-declarative-applier.md |
| rel07.0-uc001-collector-ingress | The collector agent retains and relays an attributable trace batch | 07.0 | partial | test-rel07.0-collector | docs/specs/use-cases/rel07.0-uc001-collector-ingress.md |
| rel08.0-uc001-mesh-observer | The observer agent discovers the mesh and shows a fleet view | 08.0 | implemented | test-rel08.0-observer | docs/specs/use-cases/rel08.0-uc001-mesh-observer.md |

<!-- field: test_suite_index | form: "typed table" -->
## Test suite index

| id | title | release | traces | path |
| --- | --- | --- | --- | --- |
| test-rel00.0-chatbot-mesh | Test Suite for Release 00.0 RAG Server and Single-RAG Chatbot Turn | "00.0" | [rel00.0-uc001-rag-server-query, rel00.0-uc002-chatbot-turn] | docs/specs/test-suites/test-rel00.0-chatbot-mesh.md |
| test-rel01.0-routing-fanout | Test Suite for Release 01.0 Source Routing, Tier Selection, and Multi-RAG Fan-Out | "01.0" | [rel01.0-uc001-routed-fanout] | docs/specs/test-suites/test-rel01.0-routing-fanout.md |
| test-rel02.0-observability | Test Suite for Release 02.0 Observability | "02.0" | [rel02.0-uc001-observability] | docs/specs/test-suites/test-rel02.0-observability.md |
| test-rel03.0-mesh-deployment | Test Suite for Release 03.0 Mesh Deployment and Tiered Swap | "03.0" | [rel03.0-uc001-mesh-deployment] | docs/specs/test-suites/test-rel03.0-mesh-deployment.md |
| test-rel04.0-llm-tier | Test Suite for Release 04.0 In-Cluster LLM Tier | "04.0" | [rel04.0-uc001-in-cluster-llm-tier] | docs/specs/test-suites/test-rel04.0-llm-tier.md |
| test-rel05.0-control-plane | Test Suite for Release 05.0 Control Plane | "05.0" | [rel05.0-uc001-control-plane-provisioning] | docs/specs/test-suites/test-rel05.0-control-plane.md |
| test-rel06.0-applier | Test Suite for Release 06.0 Declarative Applier Agent | "06.0" | [rel06.0-uc001-declarative-applier] | docs/specs/test-suites/test-rel06.0-applier.md |
| test-rel07.0-collector | Test Suite for Release 07.0 Collector Agent | "07.0" | [rel07.0-uc001-collector-ingress] | docs/specs/test-suites/test-rel07.0-collector.md |
| test-rel08.0-observer | Test Suite for Release 08.0 Observer Agent | "08.0" | [rel08.0-uc001-mesh-observer] | docs/specs/test-suites/test-rel08.0-observer.md |

<!-- field: foundation_document_index -->
## Foundation document index

| id | title | summary | path |
| --- | --- | --- | --- |
| chatbot-mesh-example | Chatbot Mesh: A Copyable Declarative Agent Application | Goals, boundaries, and success criteria for the chatbot mesh. | docs/VISION.md |
| architecture-chatbot-mesh-example | Chatbot Mesh Architecture | Components, ownership boundaries, deployment, and runtime flows. | docs/ARCHITECTURE.md |
| chatbot-mesh-example-roadmap | Chatbot Mesh Example Road Map | Release plan and implementation status for the application. | docs/road-map.md |
| chatbot-mesh-how-it-works | How the Chatbot Mesh Works | Operator and developer guide to the application flow. | docs/guides/how-it-works.md |

<!-- field: config_format_index | form: "empty list" -->
## Config format index

None

<!-- field: semantic_model_index -->
## Semantic model index

| id | title | summary | path |
| --- | --- | --- | --- |
| sm-agent-role-realizations | Agent Role Realizations and Bindings | Canonical role bindings for chatbot-mesh actors and profiles. | ../../docs/specs/semantic-models/agent-role-realizations.md |

<!-- field: coverage_gaps -->
## Coverage gaps

| release | status | gap |
| --- | --- | --- |
| 05.0 | partial | Control-plane provisioning remains partially implemented. |
| 07.0 | partial | Collector-agent evidence remains partial. |
| 08.0 | planned | Fleet-observer evidence remains planned. |
