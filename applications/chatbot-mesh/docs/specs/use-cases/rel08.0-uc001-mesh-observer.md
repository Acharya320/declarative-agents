---
id: rel08.0-uc001-mesh-observer
title: The observer agent discovers the mesh and shows a fleet view
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An operator brings up the demo cluster and opens the observer UX. The observer
discovers every mesh agent pod from the Kubernetes API, polls each agent's
monitor surface, and displays a fleet cast list: each agent's declared machine,
tools, current state, and recent events. While the operator runs a chat turn,
the fleet view updates to show the chatbot and RAG server state changes. The
observer's own profile is itself a declarative agent on the same runtime.

<!-- field: actor -->
## Actor

Demo presenter

<!-- field: trigger -->
## Trigger

The demo cluster is running (`mage demo:up`) and the presenter opens the observer UX at its .localhost URL.

<!-- field: flow -->
## Flow

- F1: The Helm chart deploys the observer alongside the mesh agents, with a ServiceAccount that grants namespace-scoped pod listing.
- F2: The observer launches its monitor and control servers, then starts discovery by polling the Kubernetes API for pods matching the release label selector.
- F3: The observer also lists the release Deployments and Services and polls the metrics.k8s.io PodMetrics endpoint for per-pod CPU and memory.
- F4: For each discovered pod, the observer polls `/monitor/machine`, `/monitor/state`, `/monitor/tools`, and `/monitor/events`.
- F5: The UX renders the deployment topology -- pods grouped by owning Deployment and fronting Service -- with a card per agent showing declared machine name and purpose, tool list, current state, recent events, and the pod's CPU and memory.
- F6: The presenter seeds a corpus via the corpus-ingest agent; the observer's next poll shows the ingest agent's state transition.
- F7: The presenter asks a question in the chatbot SPA; the observer's fleet view shows the chatbot entering its turn states and the RAG servers processing queries.
- F8: The presenter opens the observer's own profile.yaml to close the narrative -- the fleet instrument is itself a declarative agent.
- F9: Lifecycle exit stops the observer cleanly.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd008-observer-agent](../software-requirements/srd008-observer-agent.md) [R1.1](../software-requirements/srd008-observer-agent.md#r11), [R1.2](../software-requirements/srd008-observer-agent.md#r12), [R1.3](../software-requirements/srd008-observer-agent.md#r13), [R1.4](../software-requirements/srd008-observer-agent.md#r14), [R1.5](../software-requirements/srd008-observer-agent.md#r15), [R1.6](../software-requirements/srd008-observer-agent.md#r16) | Kubernetes API discovery of pods, Deployments, and Services with label selector and service-account auth |
| T2 | [srd008-observer-agent](../software-requirements/srd008-observer-agent.md) [R2.1](../software-requirements/srd008-observer-agent.md#r21), [R2.2](../software-requirements/srd008-observer-agent.md#r22), [R2.3](../software-requirements/srd008-observer-agent.md#r23), [R2.4](../software-requirements/srd008-observer-agent.md#r24) | monitor endpoint polling and fan-in aggregation |
| T3 | [srd008-observer-agent](../software-requirements/srd008-observer-agent.md) [R3.1](../software-requirements/srd008-observer-agent.md#r31), [R3.2](../software-requirements/srd008-observer-agent.md#r32), [R3.3](../software-requirements/srd008-observer-agent.md#r33), [R3.4](../software-requirements/srd008-observer-agent.md#r34), [R3.5](../software-requirements/srd008-observer-agent.md#r35), [R3.6](../software-requirements/srd008-observer-agent.md#r36) | fleet UX with auto-refresh, per-agent cards, and topology layout |
| T4 | [srd008-observer-agent](../software-requirements/srd008-observer-agent.md) [R4.1](../software-requirements/srd008-observer-agent.md#r41), [R4.2](../software-requirements/srd008-observer-agent.md#r42), [R4.3](../software-requirements/srd008-observer-agent.md#r43) | self-monitoring and lifecycle exit |
| T5 | [srd008-observer-agent](../software-requirements/srd008-observer-agent.md) [R5.1](../software-requirements/srd008-observer-agent.md#r51), [R5.2](../software-requirements/srd008-observer-agent.md#r52), [R5.3](../software-requirements/srd008-observer-agent.md#r53) | namespace-scoped RBAC over pods, Deployments, Services, and pod metrics |
| T6 | [srd008-observer-agent](../software-requirements/srd008-observer-agent.md) [R6.1](../software-requirements/srd008-observer-agent.md#r61), [R6.2](../software-requirements/srd008-observer-agent.md#r62), [R6.3](../software-requirements/srd008-observer-agent.md#r63) | metrics.k8s.io resource-usage polling with explicit unavailable state |
| T7 | [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [R5.1](../software-requirements/srd003-chatbot-deployment.md#r51), [R5.2](../software-requirements/srd003-chatbot-deployment.md#r52), [R5.3](../software-requirements/srd003-chatbot-deployment.md#r53), [R5.4](../software-requirements/srd003-chatbot-deployment.md#r54), [R5.5](../software-requirements/srd003-chatbot-deployment.md#r55), [R5.6](../software-requirements/srd003-chatbot-deployment.md#r56), [R5.7](../software-requirements/srd003-chatbot-deployment.md#r57), [R5.8](../software-requirements/srd003-chatbot-deployment.md#r58), [R5.9](../software-requirements/srd003-chatbot-deployment.md#r59), [R5.10](../software-requirements/srd003-chatbot-deployment.md#r510), [R5.11](../software-requirements/srd003-chatbot-deployment.md#r511) | observer deployed as part of the mesh chart |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The observer discovers every mesh agent pod in the release namespace by label selector.

<!-- field: traces -->
#### Traces

- [srd008-observer-agent](../software-requirements/srd008-observer-agent.md) [AC1](../software-requirements/srd008-observer-agent.md#ac1)

<!-- sc: S2 -->
### S2

The fleet view shows each agent's declared machine, tools, current state, and recent events from polled monitor data.

<!-- field: traces -->
#### Traces

- [srd008-observer-agent](../software-requirements/srd008-observer-agent.md) [AC2](../software-requirements/srd008-observer-agent.md#ac2)
- [srd008-observer-agent](../software-requirements/srd008-observer-agent.md) [AC3](../software-requirements/srd008-observer-agent.md#ac3)

<!-- sc: S3 -->
### S3

During a live chat turn, the fleet view updates to reflect chatbot and RAG server state changes on the next poll.

<!-- field: traces -->
#### Traces

- [srd008-observer-agent](../software-requirements/srd008-observer-agent.md) [AC2](../software-requirements/srd008-observer-agent.md#ac2)
- [srd008-observer-agent](../software-requirements/srd008-observer-agent.md) [AC3](../software-requirements/srd008-observer-agent.md#ac3)

<!-- sc: S4 -->
### S4

The observer's own monitor surfaces report fleet-level metrics and lifecycle exit stops cleanly.

<!-- field: traces -->
#### Traces

- [srd008-observer-agent](../software-requirements/srd008-observer-agent.md) [AC4](../software-requirements/srd008-observer-agent.md#ac4)

<!-- sc: S5 -->
### S5

The observer's RBAC is namespace-scoped and read-only; it holds no write authority.

<!-- field: traces -->
#### Traces

- [srd008-observer-agent](../software-requirements/srd008-observer-agent.md) [AC5](../software-requirements/srd008-observer-agent.md#ac5)

<!-- sc: S6 -->
### S6

The fleet view's topology matches the running Deployments, Services, and pods within one poll, grouping pods by owning Deployment and fronting Service.

<!-- field: traces -->
#### Traces

- [srd008-observer-agent](../software-requirements/srd008-observer-agent.md) [AC6](../software-requirements/srd008-observer-agent.md#ac6)

<!-- sc: S7 -->
### S7

Each pod card shows CPU and memory from metrics.k8s.io, or an explicit "metrics unavailable" state when the metrics API is absent.

<!-- field: traces -->
#### Traces

- [srd008-observer-agent](../software-requirements/srd008-observer-agent.md) [AC7](../software-requirements/srd008-observer-agent.md#ac7)

<!-- field: failure_scenarios -->
## Failure scenarios

| id | scenario | expected |
| --- | --- | --- |
| E1 | A mesh agent pod is unreachable (crashed or not yet ready). | The observer shows the agent as degraded in the fleet view; other agents remain unaffected. |
| E2 | The observer is deployed without the ServiceAccount or Role. | Discovery fails with a clear error; the observer enters a failed state, not a silent empty view. |
| E3 | A new RAG server is added to the mesh after the observer starts. | The next discovery poll picks up the new pod and the fleet view adds a card for it. |
| E4 | The cluster has no metrics-server (for example a bare kind cluster), so metrics.k8s.io is absent. | Topology and monitor fan-in still render; each card shows "metrics unavailable" rather than the observer failing. |

<!-- field: out_of_scope -->
## Out of scope

- SSE streaming from agent monitors (polling only per srd008 R2.4).
- Cross-namespace or cluster-wide observation.
- Controlling or reconfiguring agents from the observer UX.

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel08.0-observer](../test-suites/test-rel08.0-observer.md)
