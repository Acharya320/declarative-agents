---
id: rel03.0-uc001-mesh-deployment
title: Helm Chart Deploys the Mesh and Hot-Swaps a RAG Source
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An operator deploys the chatbot mesh (srd003) with one Helm chart onto a kind
cluster: the chatbot agent, two RAG (agent + Chroma) pairs from a values list,
a chat LLM, a Dolt backend, and an OTel collector, all over rel06.0 binary-only
images with profiles from ConfigMaps. The values list co-generates both the RAG
Deployments/Services and the chatbot's topology data and REST allowlist from one
source of truth. The operator then exercises the tiered hot swap: repointing an
existing RAG is a Service change touching no agent config, and adding a RAG is a
values change that re-renders ConfigMaps and rolls out the chatbot with a
drained handoff: the old pod completes active HTTP turns before it exits, and
the replacement pod serves subsequent turns.

<!-- field: actor -->
## Actor

An operator deploying and evolving the mesh.

<!-- field: trigger -->
## Trigger

The operator installs the chart on a kind cluster and later edits the values
RAG list.

<!-- field: preconditions -->
## Preconditions

- A kind cluster is available and the rel06.0 binary-only runtime image is loadable.
- The chart, values schema, and the chatbot/RAG profiles are chart-owned artifacts.
- The Dolt backend is deployable for chatbot checkpoint persistence.

<!-- field: flow -->
## Flow

- F1: helm install deploys the chatbot, two values-driven RAG pairs, the LLM, the Dolt backend, and the collector, with ingress to the chatbot UX and internal Services for each component.
- F2: The chart renders the chatbot's ordered RAG topology, REST network allowlist, and monitor upstreams from the same values list that creates the RAG objects; one selected-target REST operation passes ValidateDefinition before serving.
- F3: The operator repoints an existing RAG by changing a Kubernetes Service selector or endpoint, with no agent config change and no agent restart.
- F4: The operator adds a RAG by editing the values list; the chart re-renders the ConfigMaps and rolls out the chatbot. Kubernetes asks the old pod to execute its declared lifecycle exit, the chat server drains each active machine_request response, and the replacement pod serves subsequent turns.
- F5: Every agent exports OTLP to the collector, and the trace backend is available for the observability panel to query.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC1](../software-requirements/srd003-chatbot-deployment.md#ac1) | One chart deploys the chatbot, values-driven RAG pairs, LLM, Dolt, and collector over rel06.0 images with ConfigMap profiles. |
| T2 | [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC2](../software-requirements/srd003-chatbot-deployment.md#ac2) | Co-generation renders the chatbot RAG clients from the RAG values list and validates the rendered rest.yaml. |
| T3 | [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC3](../software-requirements/srd003-chatbot-deployment.md#ac3) | The tiered swap is a Service repoint or a values re-render plus a request-draining rollout. |
| T4 | test-rel03.0-mesh-deployment mesh_deployment | The suite records the kind-cluster install, the repoint, and the add-RAG rollout with expected outcomes. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

One chart deploys the chatbot, two values-driven RAG pairs, the LLM, the Dolt backend, and the collector on a kind cluster over rel06.0 binary-only images with ConfigMap-mounted profiles.

<!-- field: traces -->
#### Traces

- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [R1.1](../software-requirements/srd003-chatbot-deployment.md#r11)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [R1.3](../software-requirements/srd003-chatbot-deployment.md#r13)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC1](../software-requirements/srd003-chatbot-deployment.md#ac1)

<!-- sc: S2 -->
### S2

The chatbot topology data and REST allowlist are co-generated from the RAG object list, while one selected-target operation validates and the machine remains source-count independent.

<!-- field: traces -->
#### Traces

- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [R2.1](../software-requirements/srd003-chatbot-deployment.md#r21)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [R2.2](../software-requirements/srd003-chatbot-deployment.md#r22)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC2](../software-requirements/srd003-chatbot-deployment.md#ac2)

<!-- sc: S3 -->
### S3

Repointing an existing RAG is a Service change with no agent config change; adding a RAG is a values re-render plus a rollout in which the old chatbot completes every active HTTP turn before exit and the replacement serves the next turn.

<!-- field: traces -->
#### Traces

- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [R3.1](../software-requirements/srd003-chatbot-deployment.md#r31)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [R3.2](../software-requirements/srd003-chatbot-deployment.md#r32)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC3](../software-requirements/srd003-chatbot-deployment.md#ac3)

<!-- field: out_of_scope -->
## Out of scope

- In-process client-registry mutation of a running agent.
- The provisioning surface (the deployment API and its panel), owned by the control-plane surface (provisioning-workflow-orchestrator, creator, applier) and out of scope for this data-plane use case.
- Production cluster hardening beyond the kind smoke deployment.

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel03.0-mesh-deployment](../test-suites/test-rel03.0-mesh-deployment.md)
