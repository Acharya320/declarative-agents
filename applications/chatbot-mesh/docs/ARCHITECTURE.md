---
id: architecture-chatbot-mesh-example
title: "Chatbot Mesh Architecture"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: capability_classification | form: "tree" -->
## Capability classification

- ownership: agent-owning
- module_status: implemented
- agents_contributed: 5
- composition_wrappers:
  1. corpus-ingest
  1. applier

<!-- field: overview | form: "tree" -->
## Overview

- summary:
  The chatbot mesh is a program of declarative agents that run on the agent-core
  runtime. A browser-facing chatbot agent runs a request-scoped turn: it embeds
  the message once, routes the question to a chat model, fans the embedding out
  to one or more RAG server agents, composes the answer prompt from the surviving
  sources through command-state selectors, and answers. Each RAG server agent
  serves a vector-in query against a Chroma collection. A corpus-ingest agent
  seeds those collections from a document directory under model control.

  Runtime ownership stays with agent-core. That project provides the agent
  binary, the state-machine engine, the REST client and server tools, the
  machine_request binding, the monitor, the checkpoint store, and the container
  image. The application provides the YAML programs, the user interface, the Helm
  chart, and the documents. Operators run each agent as the agent-core image with
  a mounted profile: `agent --profile /profiles/agents/<name>/profile.yaml`.

  This architecture covers the data plane and the control plane. The control plane
  -- the `provisioning-workflow-orchestrator` actor assuming Executor / Workflow Orchestrator, the
  `creator` actor assuming Creator / Instance Manager and Agent Owner / Agent
  Lifecycle Manager, and the provisioning surface -- is release rel05.0. Release
  rel06.0 realizes the provisioning surface through the `applier` actor, which
  assumes Executor / Actuation Agent (srd006).
- turn_flow_diagram:
  ```mermaid
  flowchart LR
    U[Browser SPA] -->|chat request| CB[Chatbot agent]
    CB -->|embed once| OLL[Ollama embed]
    CB -->|select model tier| TS[$tool tier selector -> chat-LLM word]
    CB -->|fan out one embedding| R0[RAG server 0]
    CB -->|fan out one embedding| R1[RAG server 1]
    R0 -->|query_embeddings| CH0[(Chroma collection 0)]
    R1 -->|query_embeddings| CH1[(Chroma collection 1)]
    CB -->|compose surviving chunks| ANS[Answer from chunks]
    ING[Corpus-ingest agent] -->|seed collections| CH0
    CB -.->|monitor SSE / traceparent| OBS[Observability + collector query]
  ```

<!-- field: interfaces -->
## Interfaces

<!-- interface: --profile -->
### --profile

agent-core CLI argument pointing at a profile.yaml under this application.
Each mesh agent runs as the agent-core image with a mounted profile path.

<!-- interface: chatbot chat REST endpoint -->
### chatbot chat REST endpoint

The browser SPA posts a chat request to the chatbot's machine_request chat
endpoint; the turn machine runs one request-scoped sentence and maps the
terminal output to the HTTP response.

<!-- interface: RAG server query endpoint -->
### RAG server query endpoint

A vector-in machine_request endpoint that accepts query_embeddings and
returns chunks with embedding-model metadata. The chatbot is its client.

<!-- interface: monitor and trace surfaces -->
### monitor and trace surfaces

Each agent exposes a monitor read-state and SSE stream (agent-core srd033);
the chatbot proxies them same-origin. Request-scoped spans propagate
W3C traceparent (agent-core srd016) and export to an OTel collector and
the collector query surface for the cross-agent waterfall.

<!-- field: role_realization_alignment | form: "tree" -->
## Role realization alignment

- authority: applications/docs/specs/semantic-models/agent-role-realizations.md
- language_rule:
  Actors and profiles assume or realize canonical roles; they are not role
  definitions. Helm names identify workloads and deployment instances and
  inherit the mounted profile's bindings.
- decisions:
  1. - current_actor: chatbot
     - implemented_responsibility: Route, retrieve, compose, answer, and host the browser request boundary.
     - canonical_role_sub_role: Planner / Intent Handler; Executor / Workflow Orchestrator; Knowledge Manager / Recommendation Generator.
     - retained_or_proposed_name: Retain descriptive multi-role actor `chatbot`.
     - authority_boundary: Selects declared words and trusted topology only; no lifecycle or deployment authority.
     - evidence:
       1. agents/chatbot/request-machine.yaml
       1. docs/specs/software-requirements/srd002-chatbot-agent.md
     - migration_impact: None.
  1. - current_actor: rag-server
     - implemented_responsibility: Serve ranked chunks from one configured corpus through a vector-in query boundary.
     - canonical_role_sub_role: Knowledge Manager / Recommendation Generator.
     - retained_or_proposed_name: Retain descriptive serving actor `rag-server`.
     - authority_boundary: Read-only retrieval; no ingest, answer, lifecycle, or transport-authority decision.
     - evidence:
       1. agents/rag-server/request-machine.yaml
       1. docs/specs/software-requirements/srd001-rag-server-agent.md
     - migration_impact: None; values-named RAG units are deployment instances.
  1. - current_actor: corpus-ingest
     - implemented_responsibility: Parameterize and package the catalog corpus-ingest closure for the mesh.
     - canonical_role_sub_role: Inherited Knowledge Manager / Ontology Manager.
     - retained_or_proposed_name: Retain descriptive wrapper `corpus-ingest`.
     - authority_boundary: Mesh owns parameters and packaging; catalog owns the canonical implementation.
     - evidence:
       1. agents/corpus-ingest/profile.yaml
       1. agents/creator/request-declarations.yaml
     - migration_impact: None; it runs as a creator-launched child and has no Helm workload.
  1. - current_actor: provisioning-workflow-orchestrator
     - implemented_responsibility: Sequence optional ingestion and rollout requests and aggregate their terminal result.
     - canonical_role_sub_role: Executor / Workflow Orchestrator; not Coordinator.
     - retained_or_proposed_name: Use `provisioning-workflow-orchestrator`; implemented in GH-1077.
     - authority_boundary: No deployment, child-process, synchronization, conflict-resolution, negotiation, or consensus authority.
     - evidence:
       1. agents/provisioning-workflow-orchestrator/request-machine.yaml
       1. helm/templates/provisioning-workflow-orchestrator.yaml
     - migration_impact: GH-1077 migrated profile, Helm, SRD, integration, packaging, release, and compatibility paths together.
  1. - current_actor: creator
     - implemented_responsibility: Launch an ingest child, drive decided rollout operations, and verify health and state.
     - canonical_role_sub_role: Creator / Instance Manager and Agent Owner / Agent Lifecycle Manager.
     - retained_or_proposed_name: Retain `creator`; instance realization is an implemented primary responsibility.
     - authority_boundary: Sole deployment-API client; acts on decided values and does not orchestrate the end-to-end intent.
     - evidence:
       1. agents/creator/request-machine.yaml
       1. helm/templates/creator.yaml
     - migration_impact: None.
  1. - current_actor: applier
     - implemented_responsibility: Validate, actuate, verify, and roll back a decided deployment values document.
     - canonical_role_sub_role: Executor / Actuation Agent and Change Manager.
     - retained_or_proposed_name: Retain descriptive actor/workload `applier`; Applier is not a canonical role.
     - authority_boundary: Creator-only caller; authors no desired state and receives no runtime transport authority.
     - evidence:
       1. ../catalog/agents/applier/apply-machine.yaml
       1. agents/applier/exec-declarations.yaml
       1. helm/templates/applier.yaml
     - migration_impact: Documentation alignment only.
  1. - current_actor: collector
     - implemented_responsibility: Receive, inspect, spool, and optionally relay bounded OTLP trace batches.
     - canonical_role_sub_role: Monitor / Telemetry Collector.
     - retained_or_proposed_name: Retain compatible descriptive actor `collector`.
     - authority_boundary: Trace intake only; contrib retains metrics and Dolt scraping during migration.
     - evidence:
       1. agents/collector/machine.yaml
       1. helm/templates/collector.yaml
     - migration_impact: None; distinguish profile-backed, contrib, and metrics workloads.
  1. - current_actor: observer
     - implemented_responsibility: Discover mesh agents from the Kubernetes API and poll each agent's monitor surface to serve a fleet-level view.
     - canonical_role_sub_role: Monitor / Fleet Observer.
     - retained_or_proposed_name: Retain descriptive actor `observer`.
     - authority_boundary: Read-only over kube pod listing and agent monitor surfaces; no lifecycle, deployment, configuration, or trace-intake authority.
     - evidence:
       1. agents/observer/machine.yaml
       1. agents/observer/declarations.yaml
       1. helm/templates/observer.yaml
     - migration_impact: None.

<!-- field: components -->
## Components

<!-- component: chatbot agent -->
### chatbot agent

Run the request-scoped turn: embed once, select a chat-LLM word
word through a $tool classifier, fan the embedding out to the configured RAG
servers, thread each source through command-state $from() selectors, compose
the answer prompt from the surviving chunks, and answer. Degrade to a
mapped 200 when a source fails or rejects the embedding.

<!-- field: path -->
#### Path

agents/chatbot/

<!-- component: RAG server agent -->
### RAG server agent

Serve a persistent vector-in retrieval service. Accept a query_embeddings
request, query the configured Chroma collection, and return chunks with
embedding-model metadata. The Chroma REST config is owned inline here.

<!-- field: path -->
#### Path

agents/rag-server/

<!-- component: scenario Critic realizations -->
### scenario Critic realizations

Evaluate scenario outputs through the application-local
`chatbot-turn-critic` and `rag-query-critic` profiles. The single-turn and
degraded-RAG directories configure different scopes of the same chatbot
actor; they do not define separate roles. Both actors bind Critic / Output
Evaluator. Only the isolated single-turn topology also binds Sandbox
Validator because it proves an undeclared URL-shaped source cannot escape
the controlled mock boundary. The catalog `scenario-critic` discovers,
composes, and reduces these critics without becoming their implementation
owner.

<!-- field: path -->
#### Path

agents/chatbot/tests/ and agents/rag-server/tests/

<!-- component: corpus-ingest agent -->
### corpus-ingest agent

Seed a Chroma collection from a document directory: the canonical library
machine lists trusted paths, the model selects among them, and the machine
embeds each document and threads its vector into Chroma. The mesh owns only
the wrapper and corpus-rest endpoint/model/collection parameterization;
packaging embeds the resolved canonical closure for runtime use.

<!-- field: path -->
#### Path

agents/corpus-ingest/ wrapper; <catalog_root>/agents/knowledge-manager/corpus-ingest/ canonical build input

<!-- component: provisioning-workflow-orchestrator actor -->
### provisioning-workflow-orchestrator actor

Assumes Executor / Workflow Orchestrator for the control plane (srd004). Bind a provisioning intent
the chatbot delegates to a request-scoped run, sequence the ingestion and
reconfiguration through declared requests to the creator, decide the values
patch, and aggregate status. Hold no deployment-API authority.

<!-- field: path -->
#### Path

agents/provisioning-workflow-orchestrator/

<!-- component: creator actor -->
### creator actor

Assumes Creator / Instance Manager and Agent Owner / Agent Lifecycle Manager
for the control plane (srd005). Realize a decided
instance operation by driving the deployment API, verify the instance's health
through a rollout read, and report back. The only agent that declares a
deployment-API client; it alone holds that authority.

<!-- field: path -->
#### Path

agents/creator/

<!-- component: applier actor -->
### applier actor

Assumes Executor / Actuation Agent for the deployment plane (srd006): the declarative agent that
realizes the deployment API (srd003 R4) the creator drives and the panel reads.
Its request-machine writes the decided values patch, validates it against the
chart schema with a helm dry-run, applies it as an atomic rollout, verifies the
rollout, and rolls back on failure -- binding the helm and kubectl CLIs in its
container as exec words. Authenticated apart from chat by a NetworkPolicy. It
replaces the standard-library Go provisioner.

<!-- field: path -->
#### Path

agents/applier/

<!-- component: collector actor -->
### collector actor

Assumes Monitor / Telemetry Collector for the staged declarative OTLP trace
path (srd007). The persistent
machine receives bounded trace batches, spools them before optional relay,
preserves integration identity, and exports its own telemetry only to file.
The contrib collector remains for metrics and Dolt scraping. Profile, intake,
relay, self-loop, and Helm evidence are implemented; complete listener and
monitor lifecycle proof remains planned.

<!-- field: path -->
#### Path

agents/collector/

<!-- component: observer actor -->
### observer actor

Assumes Monitor / Fleet Observer for the demo and operational plane
(srd008). The persistent machine discovers mesh agent pods from the
Kubernetes API by label selector, polls each agent's monitor endpoints
(machine spec, state, tools, events), and serves a fleet-level UX showing
the cast list and live state. Holds read-only authority; no lifecycle,
deployment, or configuration commands.

<!-- field: path -->
#### Path

agents/observer/

<!-- component: user interface -->
### user interface

The single-page application (chat, observability, trace, and provisioning
panels) and its UX config. Rendered into the deployed agents from the Helm
chart. The provisioning panel delegates an intent to the provisioning-workflow-orchestrator.

<!-- field: path -->
#### Path

agents/chatbot/ui/

<!-- component: Helm chart -->
### Helm chart

Deploy the mesh as one chart: the chatbot Deployment and Service, a
values-driven list of RAG pairs each rendering a RAG server and a Chroma
store, an in-cluster Ollama model tier, a Dolt checkpoint backend, and an
the canonical catalog collector with a spool query surface. Co-generate the chatbot client config,
user interface, and N-RAG fan-out from the values.

<!-- field: path -->
#### Path

helm/

<!-- component: documentation -->
### documentation

Own the application's vision, architecture, road map, and specification corpus
(software requirements, use cases, test suites).

<!-- field: path -->
#### Path

docs/

<!-- field: design_decisions -->
## Design decisions

<!-- decision: 1 -->
### 1 - Application composition on a shared runtime

The application runs on the published agent-core image. Agents are YAML
profiles the runtime loads. Corpus ingest references the canonical
applications/catalog program; other agents remain application-local.
Platform requirements are cited externally.

<!-- field: benefits -->
#### Benefits

- Reusable corpus-ingest fixes have one canonical home.
- The application demonstrates that agent behavior is configuration on a reusable runtime.
- Runtime changes stay covered by agent-core's own gates.

<!-- decision: 2 -->
### 2 - The mesh owns its Chroma retrieval config

The RAG server owns its Chroma REST config inline. The corpus-ingest
wrapper supplies application REST/model/collection values while selecting
the canonical library machine, tools, and declarations.

<!-- field: benefits -->
#### Benefits

- Application endpoints and collection identity stay mesh-owned.
- Trusted discovery and ingest behavior do not fork from the library.

<!-- decision: 3 -->
### 3 - A served UI lives under its serving agent

Each user interface the mesh ships lives beneath the agent whose
`static_assets` binding serves it -- the single-page application at
`agents/chatbot/ui/` -- per the eng02-agent-ui-placement guideline
(GH-1316). The observer ships none: its fleet view is the observer bundle
compiled into agent-core, selected and configured by
`agents/observer/monitor-rest.yaml` (applications srd004 R9). The
deployment chart stays top-level under `helm/`, rather than under
`deploy/`.

<!-- field: benefits -->
#### Benefits

- A reader finds an interface beside the agent that owns and serves it.
- The chart and the profiles it deploys live in one project.

<!-- decision: 4 -->
### 4 - Co-generation stays, for now

The Helm chart renders the chatbot client config, user interface, and
N-RAG fan-out, and the observer monitor server's proxy upstreams and UI
config, from the values. The packaged profile copies are the local
integration source; the render overrides them in the cluster. Inverting this
so the profile is the source is a separate follow-up.

<!-- field: benefits -->
#### Benefits

- The move is an extraction, not a redesign.
- One source (the values) drives the deployed topology and client config.

<!-- field: technology_choices -->
## Technology choices

| component | technology | purpose |
| --- | --- | --- |
| Agent programs | YAML profiles on agent-core | The runtime loads machines, tool selections, declarations, and REST config; agent behavior stays configuration. |
| Retrieval | Chroma | Vector store the RAG servers query and the ingest agent seeds. |
| Models | Ollama | Embedding and chat models, external by default and promotable to an in-cluster tier by the chart. |
| Checkpoint backend | Dolt | Versioned-SQL checkpoint persistence for resume across restarts. |
| Observability | Canonical collector agent with spool query surface | Per-agent monitor streams and the cross-agent trace waterfall. |
| Deployment | Helm | One chart deploys the mesh with values-driven RAG pairs and an in-cluster model tier. |

<!-- field: project_structure | form: "tree" -->
## Project structure

- configuration:
  1. - path: agents/
     - role: The chatbot, RAG server, and corpus-ingest agent programs.
  1. - path: agents/chatbot/ui/
     - role: The single-page application and UX config.
  1. - path: helm/
     - role: The deployment chart.
- documentation:
  1. - path: docs/VISION.md
     - role: Goals, boundaries, and success criteria.
  1. - path: docs/ARCHITECTURE.md
     - role: Components, interfaces, and decisions.
  1. - path: docs/road-map.md
     - role: The application's own release plan.
  1. - path: docs/specs/software-requirements/
     - role: srd001 rag-server, srd002 chatbot, srd003 deployment.
  1. - path: docs/specs/use-cases/
     - role: The data-plane use cases the releases deliver.
  1. - path: docs/specs/test-suites/
     - role: The test suites that validate the use cases.
  1. - path: README.md
     - role: The decisions record and the standalone-program framing.

<!-- field: implementation_status | form: "tree" -->
## Implementation status

- current_focus: The application is a self-governing program whose canonical deployed mode is app-chatbot-mesh on shared da-platform. The data plane, application-owned Chroma/RAG behavior, shared ingress, and durable application collector are implemented here, while the control plane remains partial pending one live ingest-to-grounded-turn lifecycle proof.
- progress:
  1. - done: The data-plane mesh (RAG server, chat turn, routing, fan-out, observability, Helm deployment, in-cluster model tier) is implemented, with the agents, UX, chart, and specs in this directory.
  1. - done: The application's own magefile audit and integration entries (chroma, rag-server, chatbot, helm smoke/swap/llm-tier, control plane).
  1. - done: The scenario-critic rig discovers application-local chatbot-turn-critic and rag-query-critic realizations, retains scenario scope in their paths/configuration, and records those identities in runtime verdicts.
  1. - done: The control-plane provisioning-workflow-orchestrator (srd004) and creator (srd005), realized as mage integration:controlPlane (rel05.0, deterministic, driving a fake deployment API).
  1. - done: The creator corpus-ingest child endpoint and count-delta checks (srd005 R3.1), covered independently from integration:controlPlane.
  1. - planned: One live integration that joins an ingest-carrying provisioning intent, rollout, and a grounded turn from the newly added source.
  1. - done: The rel07.0 collector profile, bounded intake and spool rig, unchanged relay identity, self-loop exclusion, staged Helm trace switch, and retained contrib metrics evidence.
  1. - planned: Collector lifecycle proof that exits a blocked await, releases and rebinds receiver, control, and monitor listeners, and records final monitor state.
  1. - done: The applier deployment surface (srd006) that replaces the Go provisioner. mage integration:applier drives every declared apply and rollout terminal against recording helm and kubectl, while mage integration:applierLive proves real on-cluster apply and rollback.
  1. - done: apprig app:up|status|seed|verify|diagnose|down in app-chatbot-mesh, with unique chatbot/observer/telemetry hosts, application-owned Chroma, an application bucket, persistent collector WAL, concurrent operation with Coding Agent and Agent Architecture, and retained query history across down/redeploy (GH-2482).

<!-- field: related_documents -->
## Related documents

| doc | purpose |
| --- | --- |
| docs/VISION.md | Goals, boundaries, and success criteria. |
| docs/road-map.md | The application's release plan. |
| docs/specs/software-requirements/srd001-rag-server-agent.md | The RAG server agent requirements (relocated, renumbered). |
| docs/specs/software-requirements/srd002-chatbot-agent.md | The chatbot agent requirements (relocated, renumbered). |
| docs/specs/software-requirements/srd003-chatbot-deployment.md | The deployment requirements (relocated, renumbered), including the provisioning authority boundary (R4). |
| docs/specs/software-requirements/srd004-provisioning-workflow-orchestrator.md | The provisioning-workflow-orchestrator agent requirements (control plane). |
| docs/specs/software-requirements/srd005-creator.md | The creator agent requirements (control plane). |
| docs/specs/software-requirements/srd006-applier.md | The applier agent requirements (deployment plane) that replace the Go provisioner. |
| README.md | The decisions record and standalone-program framing. |
