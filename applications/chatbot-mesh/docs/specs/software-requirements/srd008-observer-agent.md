---
id: srd008-observer-agent
title: Observer Agent Functional Block
depends_on:
  - srd_id: srd003-chatbot-deployment
    symbols_used: [R5]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: role_realization_decision | form: "tree" -->
## Role realization decision

- current_actor: observer
- implemented_responsibility: Discover mesh agents from the Kubernetes API and poll each agent's monitor surface to serve a fleet-level view.
- canonical_role_sub_role: Monitor / Fleet Observer.
- retained_or_proposed_name: Retain descriptive actor `observer`.
- authority_boundary: Read-only over kube pod listing and agent monitor surfaces; no lifecycle, deployment, configuration, or trace-intake authority.
- evidence:
  1. agents/observer/machine.yaml
  1. agents/observer/declarations.yaml
  1. helm/templates/observer.yaml
- migration_impact: None.

<!-- field: context -->
## Context

The `observer` actor assumes Monitor / Fleet Observer and implements a fleet-level
view of the running mesh. Each agent in the mesh already self-serves its own
monitor surface (agent-core srd033): machine spec, current state, tools, metrics,
events, and an SSE event stream. The `runtime-state-reader` profile in
`applications/catalog` reads one agent's monitor surface. The observer generalizes
this to the full mesh: it discovers members from the Kubernetes API and fans into
every member's monitor endpoints.

The observer is the co-star of the demo deck (#1125). The audience meets the mesh
through the observer's UX — seeing every agent's declared machine, tools, and live
state — and the deck closes by opening the observer's own profile, proving the
instrument is itself a declarative agent on the same runtime.

<!-- field: problem -->
## Problem

The mesh has no fleet-level view. Each agent serves its own monitor surface, and an
operator must visit each one individually to understand what is running and what each
agent is doing. The demo deck needs one window that answers both questions for the
audience: what agents compose the mesh, and what is each one doing right now.

<!-- field: goals -->
## Goals

- G1: Discover mesh agents from the Kubernetes API by label selector within the release namespace.
- G2: Poll each discovered agent's monitor endpoints and aggregate the results into a fleet view.
- G3: Serve a UX that shows the fleet cast list, per-agent state, and recent events.
- G4: Expose standard control and monitor surfaces for the observer itself.
- G5: Run entirely as a declarative agent profile on the unmodified agent-core runtime.
- G6: Show the deployment topology and per-pod resource usage, so the fleet view answers what the deployment looks like and what it consumes.

<!-- req: R1 -->
## R1 - Mesh discovery

<!-- req: R1.1 -->
### R1.1

The observer SHALL discover mesh agent pods by issuing a REST client GET to the Kubernetes API pod listing endpoint with a label selector scoped to the Helm release.

<!-- req: R1.2 -->
### R1.2

Authentication SHALL use the in-cluster service-account bearer token mounted by the Kubernetes downward API.

<!-- req: R1.3 -->
### R1.3

Discovery SHALL run on a declared polling interval, not a one-time snapshot, so that pods added or removed after launch appear in the fleet view.

<!-- req: R1.4 -->
### R1.4

The label selector and namespace SHALL be declared REST client config values, not constants in tool config.

<!-- req: R1.5 -->
### R1.5

Discovery SHALL also list the release Deployments and Services in the namespace by label selector, so the fleet view can render the deployment topology, not only the pod set.

<!-- req: R1.6 -->
### R1.6

Observer configuration -- selectors, namespace, intervals, and endpoints -- SHALL come from REST client and tool declarations resolved at load, following the mesh convention; the topology and resource additions SHALL introduce no new environment variable.

<!-- req: R2 -->
## R2 - Monitor fan-in

<!-- req: R2.1 -->
### R2.1

For each discovered agent pod, the observer SHALL poll the agent's `/monitor/machine`, `/monitor/state`, `/monitor/tools`, and `/monitor/events` endpoints via REST client GET.

<!-- req: R2.2 -->
### R2.2

The polling interval SHALL be a declared REST client config value.

<!-- req: R2.3 -->
### R2.3

A fan-in request that fails or times out SHALL NOT prevent polling of other agents; the fleet view SHALL report the unreachable agent as degraded.

<!-- req: R2.4 -->
### R2.4

The observer SHALL NOT consume SSE streams (`/monitor/events/stream`). The agent-core REST client tooling provides no SSE consumer (stream_events is a server-side binding only). This is an explicit polling design, not a workaround; a future agent-core client streaming capability is a separate concern.

<!-- req: R3 -->
## R3 - Fleet UX

<!-- req: R3.1 -->
### R3.1

The observer SHALL serve a UX that displays a card per discovered agent showing its declared machine name and purpose, tool list, current state, and recent events.

<!-- req: R3.2 -->
### R3.2

The UX SHALL auto-refresh on the declared polling interval.

<!-- req: R3.3 -->
### R3.3

The UX SHALL distinguish agents that are reachable from agents that are unreachable or degraded.

<!-- req: R3.4 -->
### R3.4

UX assets SHALL be served by the observer's own monitor REST server from the observer bundle compiled into agent-core (applications srd004 R9), selected by a static_assets bundle binding whose config map names the title, monitored agents, and trace backend; the observer directory SHALL carry no UI files.

<!-- req: R3.5 -->
### R3.5

The UX SHALL render the deployment topology -- pods grouped by their owning Deployment and the Services fronting them -- and show each pod's CPU and memory usage on its card, or an explicit "metrics unavailable" state when resource data is absent.

<!-- req: R3.6 -->
### R3.6

The UX SHALL reach other agents' monitor surfaces and the trace backend only through the observer's same-origin GET /monitor-proxy/{agent}/{path...} route, whose upstreams name the same agents as the UX config (applications srd004 R2.1, R2.3).

<!-- req: R4 -->
## R4 - Self-monitoring and control

<!-- req: R4.1 -->
### R4.1

The observer SHALL run standard control and monitor servers with health, lifecycle exit, machine spec, current state, tools, metrics, events, event stream, and OpenAPI surfaces.

<!-- req: R4.2 -->
### R4.2

Lifecycle exit SHALL stop the discovery and polling loop, then stop the monitor and control servers.

<!-- req: R4.3 -->
### R4.3

Monitor state SHALL report discovery count, reachable count, unreachable count, total polls, and last poll timestamp.

<!-- req: R5 -->
## R5 - RBAC and security

<!-- req: R5.1 -->
### R5.1

The observer's Helm template SHALL create a ServiceAccount, Role, and RoleBinding granting only `get` and `list` verbs on `pods`, `services`, `deployments` (apps), and `pods` under `metrics.k8s.io` in the release namespace.

<!-- req: R5.2 -->
### R5.2

The observer SHALL NOT request cluster-wide RBAC; discovery is namespace-scoped.

<!-- req: R5.3 -->
### R5.3

The observer SHALL NOT write to the Kubernetes API, modify any agent's state, or issue lifecycle commands to other agents.

<!-- req: R6 -->
## R6 - Resource usage

<!-- req: R6.1 -->
### R6.1

The observer SHALL poll the namespaced `metrics.k8s.io` PodMetrics endpoint on the declared interval and associate each pod's CPU and memory usage with its fleet entry.

<!-- req: R6.2 -->
### R6.2

When the metrics API is absent (for example a kind cluster with no metrics-server), the observer SHALL report resource usage as explicitly unavailable rather than failing discovery or the fleet view.

<!-- req: R6.3 -->
### R6.3

Resource polling SHALL follow the R2.3 degradation rule -- a failed metrics read SHALL NOT prevent topology discovery or monitor fan-in.

<!-- field: non_goals -->
## Non-goals

- Consuming SSE event streams from agent monitor surfaces.
- Writing to the Kubernetes API or modifying agent state.
- Cluster-wide or cross-namespace discovery.
- Trace intake, relay, or observability pipeline functions (those belong to the collector).
- Replacing per-agent monitor UIs; the fleet view complements them.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

The observer discovers mesh agent pods by label selector within the release namespace using a declared REST client and in-cluster service-account auth.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R1.4

<!-- ac: AC2 -->
### AC2

The observer polls each discovered agent's monitor endpoints and aggregates machine spec, state, tools, and events into one fleet data set, degrading gracefully for unreachable agents.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3
- R2.4

<!-- ac: AC3 -->
### AC3

The served UX shows one card per agent with machine, tools, state, and events, auto-refreshes, and distinguishes reachable from unreachable agents.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3
- R3.4
- R3.6

<!-- ac: AC4 -->
### AC4

Standard control and monitor surfaces expose fleet-level metrics and lifecycle exit stops all loops and servers.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
- R4.3

<!-- ac: AC5 -->
### AC5

RBAC is namespace-scoped, read-only, and the observer holds no write or lifecycle authority over any other agent or the Kubernetes API.

<!-- field: traces -->
#### Traces

- R5.1
- R5.2
- R5.3

<!-- ac: AC6 -->
### AC6

The observer discovers Deployments and Services alongside pods and the UX renders the deployment topology grouped by owning Deployment and fronting Service; the topology and resource additions introduce no new environment variable.

<!-- field: traces -->
#### Traces

- R1.5
- R1.6
- R3.5

<!-- ac: AC7 -->
### AC7

The observer polls metrics.k8s.io PodMetrics and shows per-pod CPU and memory on each card, reporting resource usage as explicitly unavailable rather than failing when the metrics API is absent.

<!-- field: traces -->
#### Traces

- R6.1
- R6.2
- R6.3

<!-- field: platform_dependencies | form: "typed table" -->
## Platform dependencies

| srd_id | provider | symbols_used |
| --- | --- | --- |
| srd033-monitor | agent-core | [R1, R2, R3, R4, R5, R6, R7] |

<!-- field: implementation_status | form: "tree" -->
## Implementation status

- implemented:
  1. "The aggregated fleet payload is exposed on a readable monitor surface through the srd033 R7 command_state view: the observer serves a /monitor/fleet endpoint whose labels allowlist names the aggregation steps, and the UX unwraps that per-label envelope into the agents, deployments, services, and pod_metrics it renders, degrading to the empty view when a step has no output yet (GH-1270)."
  1. Monitor fan-in (R2.1, R2.3) is realized as four MachineSpec for_each passes, one per monitor endpoint, each expanding a single-request read word (read_agent_machine, read_agent_state, read_agent_tools, read_agent_events) over the pods discovered this cycle. Iteration lives in the machine rather than a compound word, so one REST dispatch performs one request (srd028 R5.8, R1.5). Each pass targets a pod by composing its address and its own monitor port out of the discovered pod object (srd028 R14.6, R14.7); every agent publishes that port as the dot-free monitorPort pod label because container ports sit in an array the selector grammar cannot address, and the observer's client port allowlist authorizes the resolved value. Each pass declares its own join label, so the four joins are separately addressable in the fleet allowlist. CommandError is listed in continue_on so one unreachable agent degrades only itself; each agent's reachability is therefore its own item's recorded result signal rather than the join's succeeded and failed counts. integration:observer proves the four join labels are served and that a cycle whose discovery degrades still exits cleanly (GH-1319).
  1. "The fleet UX renders from the four fan-in joins (R3.1, R3.3): fleetData zips them into one card per discovered pod, joining the joins on the pod name each item carries. A card's machine, tools, state, and events come from that pod's entry in the matching join, mapped to the monitor view shapes — machine is flat, state nests the run snapshot, tools carries tools, events carries recent_events. A card is keyed by pod name so its per-pod CPU and memory lookup resolves (R3.5). Reachability is each item's own recorded result signal rather than the join's counts, because the fan-in continues past CommandError and so counts every dispatched item as succeeded. An absent or unavailable join leaves its fields empty instead of failing the view (GH-1319)."
  1. Discovery selects agent pods, not every pod in the release (R1.1). The release selector alone also matched the infrastructure the mesh runs on -- the Dolt and Chroma StatefulSets -- which serve no monitor surface, so the fan-in reported two permanently unreachable agents. The pod query composes the release selector with the monitor-port label every agent publishes, so it selects exactly the pods that have a monitor surface. The requirement is composed in the pod operation rather than added to the shared selector value because the Deployment, Service, and pod-metrics queries reuse that value and carry no such label, and it introduces no environment variable (R1.6, GH-1474).
  1. "Discovery and monitor fan-in are proven on the kind demo (AC1, AC2). On a running demo cluster the observer discovers its agent pods through the kubectl-proxy sidecar and reads every one of them: measured across the four passes, four discovered agents on four distinct ports -- chatbot 18082, rag-server 18087, collector 18192, observer 18202 -- each returning HTTP 200 from a composed pod address. Discovery follows membership as it changes, tracking pod replacement across restarts and shrinking to three agents and back when a Deployment scales to zero and returns. A member whose monitor read fails degrades alone: with two unreadable members present the other four continued reporting and the cycle advanced past iteration 1462. Resource polling degrades as R6.2 requires -- the kind demo registers no metrics.k8s.io API, so poll_pod_metrics reports its 404 and the cycle proceeds to the next stage unaffected (GH-1474)."
  1. integration:helmSmoke provides automated in-cluster fan-in evidence (R2.1, R2.3). It reads coherent fleet cycles through the observer Service, matches discovery to the release pods carrying monitorPort, and proves each of four joins contains exactly one HTTP-200 AgentMonitorRead per discovered agent with the selected podIP:monitorPort authority. It then relabels only the observer to an authorized but unserved port and proves the other agents remain reachable through a later metrics step, restoring the label on every path (GH-1486).
  1. integration:helmSmoke installs a pinned, test-only metrics-server and proves populated PodMetrics for every discovered agent (R6.1). The UI sums CPU and memory quantities across all containers, including the observer's agent and kube-proxy containers, and renders the aggregate on the pod card. The persistent demo still proves the independent absent-API degradation path (R6.2). metrics-server remains cluster test infrastructure, not an application chart dependency (GH-1488, GH-1489).
  1. integration:helmSwap proves the fleet updates during a live chat turn (R3.2, AC3). Its deterministic model mock holds the actual answer request while a newer observer snapshot shows the chatbot run and event clock advance into the answer path and rag0 record its query transition, with both agents reachable. Request classification uses prompt content rather than call-order parity (GH-1487).
