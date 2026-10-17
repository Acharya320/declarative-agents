---
id: rel04.0-uc001-in-cluster-llm-tier
title: A Fresh-Cluster Install Answers a Turn With No External LLM Dependency
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An operator installs the chatbot-mesh chart with defaults on a fresh Kubernetes
cluster. The in-cluster LLM tier (srd003 R6) -- an optional Ollama subchart
enabled by default -- deploys Ollama with PVC-backed model storage and a
model-preload Job that pulls the embedding model, every declared chat model, and
the tier-selector model before the tier reports ready. The agent pods gate their
readiness on those models being present, so no turn fails at runtime for a
missing model. Once ready, a chat turn is answered end to end with zero external
LLM dependency. An operator running a shared Ollama instead sets
ollama.enabled=false with an endpoint override, reproducing the external-endpoint
behavior; the co-generated client entries (R2) are identical either way.

<!-- field: actor -->
## Actor

An operator installing and running the mesh on a fresh cluster.

<!-- field: trigger -->
## Trigger

The operator runs helm install with default values on a fresh cluster, or with
ollama.enabled=false and an endpoint override.

<!-- field: preconditions -->
## Preconditions

- The chart ships the Ollama subchart enabled by default with the model names declared once in values.
- A CPU-only kind cluster is available for the smoke path (GPU scheduling is values-driven for production).
- The embedding, chat, and tier-selector model names are the same values that render the agents' client config.

<!-- field: flow -->
## Flow

- F1: helm install with defaults deploys the Ollama subchart with PVC-backed model storage alongside the chatbot, RAG units, Dolt, and collector.
- F2: A model-preload Job pulls the embedding model, every declared chat model, and the tier-selector model into the PVC before the LLM tier reports ready.
- F3: The agent pods gate their readiness on the LLM tier reporting the required models present (Ollama /api/tags), so a missing model is a deploy-time gate rather than a runtime turn failure.
- F4: Once every workload is ready, a chat turn is answered end to end with zero external LLM dependency.
- F5: Setting ollama.enabled=false with an endpoint override points the embedding and chat client values at an operator-supplied service, reproducing the external-endpoint behavior; the co-generated client entries are identical.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC6](../software-requirements/srd003-chatbot-deployment.md#ac6) | The in-cluster LLM tier, model preload, readiness gating, and enabled/external switch are specified. |
| T2 | [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC2](../software-requirements/srd003-chatbot-deployment.md#ac2) | Both the in-cluster and external paths render into the same co-generated client entries (R2 unchanged). |
| T3 | test-rel04.0-llm-tier llm_tier | The suite records the default-install turn, the disabled-tier override, and the preload-gate ordering. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

A default helm install on a fresh cluster answers a chat turn with no external LLM dependency; the Ollama subchart is deployed with PVC model storage.

<!-- field: traces -->
#### Traces

- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [R6.1](../software-requirements/srd003-chatbot-deployment.md#r61)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC6](../software-requirements/srd003-chatbot-deployment.md#ac6)

<!-- sc: S2 -->
### S2

The model-preload Job pulls the embedding, chat, and tier-selector models -- named once in values -- into the PVC before the tier reports ready, and the same names render the agents' config.

<!-- field: traces -->
#### Traces

- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [R6.2](../software-requirements/srd003-chatbot-deployment.md#r62)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC6](../software-requirements/srd003-chatbot-deployment.md#ac6)

<!-- sc: S3 -->
### S3

Agent pods are not ready until the required models are present, so a missing model is a deploy-time gate rather than a runtime turn failure.

<!-- field: traces -->
#### Traces

- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [R6.3](../software-requirements/srd003-chatbot-deployment.md#r63)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC6](../software-requirements/srd003-chatbot-deployment.md#ac6)

<!-- sc: S4 -->
### S4

Setting ollama.enabled=false with an endpoint override reproduces the external-endpoint behavior with identical co-generated client entries.

<!-- field: traces -->
#### Traces

- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [R6.1](../software-requirements/srd003-chatbot-deployment.md#r61)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [R2.1](../software-requirements/srd003-chatbot-deployment.md#r21)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC6](../software-requirements/srd003-chatbot-deployment.md#ac6)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC2](../software-requirements/srd003-chatbot-deployment.md#ac2)

<!-- field: out_of_scope -->
## Out of scope

- GPU production sizing; the kind smoke path runs CPU-only small models and diverges from production sizing by design (R6.4).
- The chat turn machine itself (releases 00.0/01.0) and the deployment topology (03.0).

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel04.0-llm-tier](../test-suites/test-rel04.0-llm-tier.md)
