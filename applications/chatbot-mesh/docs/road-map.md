---
id: chatbot-mesh-example-roadmap
title: Chatbot Mesh Example Road Map
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

These releases restage the data-plane chatbot mesh as the application's own plan.
They derive from the shared catalog's rel09.0-09.3 and rel09.5; the behavior is
implemented and is being relocated into this directory. Requirement and use-case
ids are renumbered to the application's own scheme (srd001-003, rel00.0+). The
control plane (provisioning-workflow-orchestrator srd004, creator srd005, and the provisioning surface)
is release rel05.0. Release rel06.0 replaces the standard-library Go provisioning
surface with the applier (srd006), a declarative agent that binds the helm and
kubectl CLIs in its container.

<!-- field: role_alignment -->
## Role alignment

<!-- field: status -->
### Status

specified

<!-- field: implemented_decisions -->
### Implemented decisions

- Retain `chatbot` as a descriptive multi-role actor with Planner, Executor, and Knowledge Manager bindings.
- Retain `rag-server` as a descriptive Knowledge Manager / Recommendation Generator realization.
- Retain the mesh `corpus-ingest` wrapper with its inherited catalog Knowledge Manager / Ontology Manager binding.
- Bind `creator` to Creator / Instance Manager and Agent Owner / Agent Lifecycle Manager and retain the compatible name.
- Bind `applier` to Executor / Actuation Agent and Change Manager; Applier is not a canonical role.
- Bind `collector` to Monitor / Telemetry Collector and distinguish its profile-backed workload from contrib metrics infrastructure.
- Treat all Helm names as deployment identities inheriting profile bindings, not as semantic role names.

<!-- field: implemented_migrations | form: "typed table" -->
### Implemented migrations

| issue | current_actor | implemented_responsibility | canonical_role_sub_role | retained_or_proposed_name | authority_boundary | evidence | migration_impact |
| --- | --- | --- | --- | --- | --- | --- | --- |
| GH-1077 | provisioning-workflow-orchestrator | Sequence optional ingest and rollout requests through the creator. | Executor / Workflow Orchestrator; not Coordinator. | provisioning-workflow-orchestrator | No implemented synchronization, conflict-resolution, negotiation, consensus, deployment, or child-process authority. | [agents/provisioning-workflow-orchestrator/request-machine.yaml, helm/templates/provisioning-workflow-orchestrator.yaml, srd004] | Profile, Helm, SRD, integration, packaging, release, and compatibility paths move together. |

<!-- release: 00.0 | status: implemented -->
## 00.0 - RAG server and single-RAG chatbot turn

The core loop: a persistent RAG server agent (srd001, vector-in) that answers
a query_embeddings request with chunks and embedding-model metadata, and a
browser-facing chatbot agent (srd002) whose request-scoped turn embeds the
message once, queries a single RAG, composes the answer prompt from the
retrieved chunks through command-state $from(), and answers, degrading to a
mapped 200 when the RAG
fails.

<!-- field: srds -->
### SRDs

- srd001-rag-server-agent
- srd002-chatbot-agent

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| [rel00.0-uc001-rag-server-query](specs/use-cases/rel00.0-uc001-rag-server-query.md) | Persistent RAG server answers a vector-in query_embeddings request and returns chunks with embedding-model metadata. | implemented |
| [rel00.0-uc002-chatbot-turn](specs/use-cases/rel00.0-uc002-chatbot-turn.md) | Single-RAG chat turn embeds once, queries one RAG, composes via $from(), answers, and degrades to a 200 when the RAG fails. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| [test-rel00.0-chatbot-mesh](specs/test-suites/test-rel00.0-chatbot-mesh.md) | Tracers for the RAG server query and the single-RAG chatbot turn with a degraded RAG. | implemented |

<!-- release: 01.0 | status: implemented -->
## 01.0 - Tier selection and multi-RAG fan-out

The chatbot turn (srd002 R2, R3) gains the LLM $tool tier selector that classifies a
question to one chat-LLM word, and the sequential multi-RAG fan-out that
embeds once and queries N RAGs, each degraded (CommandError) or
embedding-mismatch-excluded (QueryRejected) source routing on as a visible
machine transition, compose assembling the surviving per-source chunks.

<!-- field: srds -->
### SRDs

- srd002-chatbot-agent

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| [rel01.0-uc001-routed-fanout](specs/use-cases/rel01.0-uc001-routed-fanout.md) | A tier-selected chat turn dispatches one chat-LLM word via $tool and fans the one embedding out to two RAGs, citing both corpora and degrading when one fails. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| [test-rel01.0-routing-fanout](specs/test-suites/test-rel01.0-routing-fanout.md) | Tracers for the $tool tier-selector dispatch, the two-RAG fan-out over disjoint corpora, and the degraded and embedding-mismatch paths. | implemented |

<!-- release: 02.0 | status: implemented -->
## 02.0 - Observability

The observability surface (srd002 R5.1, R5.3): a per-agent monitor panel fed
by each agent's monitor read_state and SSE stream through the chatbot's
same-origin monitor_proxy, and a cross-agent trace waterfall that renders a
turn's connected trace from the Jaeger query API, real because the
request-scoped machine exports its spans and propagates traceparent to the
RAG servers (agent-core srd016).

<!-- field: srds -->
### SRDs

- srd002-chatbot-agent

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| [rel02.0-uc001-observability](specs/use-cases/rel02.0-uc001-observability.md) | The observability panel shows each agent's run state and event stream and renders a selected turn's connected cross-agent trace as a span waterfall. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| [test-rel02.0-observability](specs/test-suites/test-rel02.0-observability.md) | Tracers for the per-agent monitor panels over monitor_proxy and the cross-agent trace waterfall from the trace backend. | implemented |

<!-- release: 03.0 | status: implemented -->
## 03.0 - Mesh deployment and tiered swap

The mesh deploys as one Helm chart (srd003) with values-driven RAG pairs,
co-generated chatbot client config and N-RAG fan-out from one source, a
tiered hot swap (repoint is a Service change; add-RAG re-renders and rolls
the chatbot after draining active turns), and OTLP observability. The
provisioning deployment API (srd003 R4) is the control-plane release
(rel05.0); this release covers the topology and the swap paths.

<!-- field: srds -->
### SRDs

- srd003-chatbot-deployment

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| [rel03.0-uc001-mesh-deployment](specs/use-cases/rel03.0-uc001-mesh-deployment.md) | One Helm chart deploys the mesh on a kind cluster with values-to-config co-generation, a Service repoint, and an add-RAG rollout that drains active turns. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| [test-rel03.0-mesh-deployment](specs/test-suites/test-rel03.0-mesh-deployment.md) | Tracers for the kind deploy smoke, the co-generation drift lock, and the repoint and add-RAG swaps. | implemented |

<!-- release: 04.0 | status: implemented -->
## 04.0 - In-cluster LLM tier

The chart promotes Ollama from an external values-driven endpoint to a
first-class, chart-owned LLM tier (srd003 R6): an optional Ollama subchart,
enabled by default, with PVC-backed model storage, a model-preload Job, agent
readiness gated on model presence, and GPU scheduling values. A fresh-cluster
install is self-contained; setting ollama.enabled=false with an endpoint
override reproduces the external-endpoint behavior, and the co-generation
contract (R2) is unchanged either way.

<!-- field: srds -->
### SRDs

- srd003-chatbot-deployment

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| [rel04.0-uc001-in-cluster-llm-tier](specs/use-cases/rel04.0-uc001-in-cluster-llm-tier.md) | A fresh-cluster install with defaults answers a chat turn with zero external LLM dependency; the disabled-tier override reproduces the external endpoint, and agents are not ready until the models are pulled. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| [test-rel04.0-llm-tier](specs/test-suites/test-rel04.0-llm-tier.md) | Tracers for the default-install turn, the disabled-tier endpoint override, the preload-gate ordering, and the CPU-only kind smoke. | implemented |

<!-- release: 05.0 | status: partial -->
## 05.0 - Control plane

The mesh gains a control plane on top of the data-plane topology: a
`provisioning-workflow-orchestrator` actor (srd004) that currently assumes Executor / Workflow
Orchestrator, binds a user's provisioning intent, sequences
the ingestion and reconfiguration workflows, and decides the values changes,
and a `creator` actor (srd005) that assumes Creator / Instance Manager and
Agent Owner / Agent Lifecycle Manager, realizes each instance operation, and alone
holds the deployment-API authority. A provisioning intent flows chatbot ->
provisioning-workflow-orchestrator -> creator: the chatbot delegates through a declared client, the
provisioning-workflow-orchestrator sequences corpus-ingest runs and the rollout, the creator executes
the request-draining rollout through the deployment API (srd003 R4) and
reports health, and the new source serves a grounded turn.
The shipped provisioning-workflow-orchestrator request machine has no synchronization,
conflict-resolution, negotiation, or consensus state, so it does not assume
canonical Coordinator; GH-1077 completed its coordinated mechanical rename.

<!-- field: srds -->
### SRDs

- srd004-provisioning-workflow-orchestrator
- srd005-creator

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| [rel05.0-uc001-control-plane-provisioning](specs/use-cases/rel05.0-uc001-control-plane-provisioning.md) | A provisioning intent flows panel -> provisioning-workflow-orchestrator -> creator and reconfigures the mesh; only the creator holds the deployment-API authority. The values-only apply -- the panel's operator intent -- is realized end to end and proven by mage integration:controlPlane. The corpus-ingest child run and count-delta checks are implemented and covered independently, but no live integration drives one ingest-carrying intent through ingest, rollout, and a grounded turn from the newly added source. The complete use case therefore remains partial. | partial |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| [test-rel05.0-control-plane](specs/test-suites/test-rel05.0-control-plane.md) | Tracer for the apply -> reconfigure -> verify flow across the provisioning-workflow-orchestrator and creator, asserting the decided values document crosses each hop and the authority boundary holds. Realized as mage integration:controlPlane (deterministic, no cluster), which also reads the panel's mesh view back through both hops (srd006 R1.5). It does not invoke the implemented ingest leg or join it to a live rollout and grounded turn; integration:chatbot and the deploy swap prove those behaviors separately. | implemented |

<!-- release: 06.0 | status: implemented -->
## 06.0 - Declarative applier agent

The deployment API (srd003 R4) moves from the standard-library Go provisioner
to the applier (srd006), a declarative agent whose container carries the helm
and kubectl CLIs bound as exec tools. A decided values patch from the creator
is validated against the chart's values.schema.json by a helm dry-run, applied
as a values-file rollout (helm upgrade --reuse-values --atomic) rather than
per-field --set construction, verified via kubectl rollout status, and rolled
back on failure. The authority boundary is unchanged -- the apply path is
reached only by the creator and authority stays set by declared, validated
config -- and no agent-core change is required.

<!-- field: srds -->
### SRDs

- srd006-applier

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| [rel06.0-uc001-declarative-applier](specs/use-cases/rel06.0-uc001-declarative-applier.md) | The applier agent validates a values patch against the chart schema, applies it as a values-file rollout, verifies the rollout, and rolls back on failure, replacing the Go provisioner with zero agent-core changes. The live on-cluster rollout shipped with GH-735, proving a patch applied through the applier changes a real release. The provisioning panel's initial mesh-view load (GET /provisioning/api/state, srd006 R1.5) shipped with GH-750/752/753 -- the applier reads the deployed release's computed values, and the creator and provisioning-workflow-orchestrator carry the flat response through to the panel, which reassembles it client-side. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| [test-rel06.0-applier](specs/test-suites/test-rel06.0-applier.md) | Tracer for the validate -> apply -> verify -> rollback flow against the applier agent, asserting schema rejection, values-file apply, rollout verification, and the authority boundary. Realized as mage integration:applier, which drives the shipped profile against recording fake helm and kubectl; the chart schema is proven against real helm and the enabled chart against its render. The live on-cluster tier is mage integration:applierLive (GH-735). The state-read chain (GH-750/752/753) is covered by unit tests pinning every hop's field mapping and failure-status behavior, and mage integration:controlPlane automates the creator/provisioning-workflow-orchestrator state round trip after the GH-755 selector fix. | implemented |

<!-- release: 07.0 | status: implemented -->
## 07.0 - Declarative collector agent

A persistent collector profile replaces contrib on the OTLP trace path.
Its machine visibly composes bounded receive, await, policy, deterministic
stdouttrace spool, optional unchanged relay to Jaeger, retry, monitor,
control, and stop. Its own telemetry is file-only, with a direct self-loop
proof. Migration is staged: contrib retains agent metrics export and Dolt
Prometheus scraping until a later metrics-ingress replacement is specified.

<!-- field: srds -->
### SRDs

- srd007-collector-agent

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| [rel07.0-uc001-collector-ingress](specs/use-cases/rel07.0-uc001-collector-ingress.md) | A chatbot trace batch enters the collector agent, remains attributable in spool and Jaeger, while contrib retains metrics evidence. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| [test-rel07.0-collector](specs/test-suites/test-rel07.0-collector.md) | Profile, intake, self-loop, staged Helm, identity, retained-evidence proofs, and lifecycle listener rebind and terminal state for collector ingress. | implemented |

<!-- release: 09.0 | status: implemented -->
## 09.0 - Canonical collector and Jaeger retirement

The mesh removes its Jaeger dependency entirely. The collector agent's
spool query surface (port 18193) replaces the Jaeger query API as the
trace backend for the UX waterfall. Collector mode derives from
externalOTLPEndpoint presence (relay when set, spool otherwise) instead
of from jaeger.enabled. Helm drops the Jaeger Deployment, Service,
values block, and schema entry. The collector program is consumed from
the catalog rather than carried locally. Integration tests assert
spool identity and metrics evidence without Jaeger cross-checks.

<!-- field: srds -->
### SRDs

- srd007-collector-agent
- srd003-chatbot-deployment

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| [rel09.0-uc001-canonical-collector](specs/use-cases/rel09.0-uc001-canonical-collector.md) | The mesh deploys with no Jaeger workload; the UX waterfall reads traces from the collector query surface and integration tests verify spool identity and metrics evidence. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| [test-rel09.0-canonical-collector](specs/test-suites/test-rel09.0-canonical-collector.md) | Helm render purity, UX waterfall transform, spool identity, metrics evidence, and catalog consumption proofs for the Jaeger-free collector. | implemented |

<!-- release: 08.0 | status: done -->
## 08.0 - Fleet observer agent

A persistent observer agent discovers mesh members from the Kubernetes API
and polls each agent's monitor surface to serve a fleet-level view. Beyond
the per-agent cast list, the observer discovers the deployment topology
(Deployments and Services) and polls metrics.k8s.io for per-pod CPU and
memory, so the UX shows what the deployment looks like and what it consumes,
grouping pods by owning Deployment and fronting Service and degrading
explicitly when the metrics API is absent. The UX auto-refreshes on a
declared interval. The observer is itself a declarative agent on the same
runtime, and its profile is the punchline of the demo deck (#1125).

The agent-core REST client has no SSE consumer, so the observer polls
monitor snapshots rather than consuming event streams. This is an explicit
design, not a workaround; client-side streaming is a separate agent-core
concern.

Once the catalog observer family ships (catalog srd021-observer, release
99.0), this mesh observer rebinds to the catalog family and srd008 becomes
the mesh-local generalization source, retaining the descriptive observer
actor name via the historical_aliases pattern (GH-1247).

<!-- field: srds -->
### SRDs

- srd008-observer-agent

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| [rel08.0-uc001-mesh-observer](specs/use-cases/rel08.0-uc001-mesh-observer.md) | The observer discovers every mesh agent and the deployment topology, polls monitor surfaces and per-pod resource usage, and displays a fleet view that updates during a live chat turn. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| [test-rel08.0-observer](specs/test-suites/test-rel08.0-observer.md) | Profile validation, kube discovery, monitor fan-in, topology grouping, per-pod resource usage, fleet UX, RBAC proof, and lifecycle exit. | implemented |
