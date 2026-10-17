---
id: rel06.0-uc001-declarative-applier
title: The Applier Agent Realizes the Deployment API and Rolls Out a Values Change
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

The creator (srd005) submits a decided, values-plane patch to the applier
(srd006) -- the declarative agent that replaces the standard-library Go
provisioner as the deployment API (srd003 R4). The applier validates the patch
against the chart's values schema with a helm dry-run, rejecting a non-conforming
patch without a rollout. For a valid patch it materializes the values as a file and
applies helm upgrade --reuse-values --atomic, re-rendering the co-generated
topology, then verifies the affected Deployment's rollout via kubectl rollout
status and reports the phase, rolling the release back on an apply or verify
failure. Throughout, the applier edits deployment values and triggers rollouts
only; the patch carries no host, URL, method, or credential for a running agent,
the apply path is reached only by the creator, and authority stays set by declared,
validated config (srd003 R4, agent-core srd028 V30). No agent-core change is
required: the actuation is a state machine binding the helm and kubectl CLIs in the
applier's container.

The live cluster tier proves both sides of the apply path. A successful patch
moves the real Helm release. A test-only post-upgrade hook then creates a
deterministic regression after Helm accepts a second upgrade; the real kubectl
verification fails, the applier returns 500 rolled_back, and real helm rollback
restores the prior computed values and a ready chatbot Deployment. The hook and
its permissions exist only in the temporary chart staged by the live test.

<!-- field: actor -->
## Actor

The creator agent, realizing a decided instance operation behind the deployment API.

<!-- field: trigger -->
## Trigger

The provisioning-workflow-orchestrator decides a values change (add a RAG source, reconfigure a
parameter) and the creator submits it to the applier's apply endpoint.

<!-- field: preconditions -->
## Preconditions

- The applier profile is deployed with its state, apply, and rollout servers on the Service name and port the provisioner served.
- The applier runs the agent-core image with the helm and kubectl CLIs delivered by the pinned CLI donor, and its ServiceAccount can upgrade the release and read the Deployment.
- The chart carries a values.schema.json encoding the deployment values constraints.
- The apply path is reached only by the creator, authenticated apart from the chat path.

<!-- field: flow -->
## Flow

- F1: The creator posts the decided values patch to the applier's apply endpoint through its declared client; it submits no host, URL, method, or credential for a running agent.
- F2: The applier validates the patch against the chart's values schema with a helm dry-run and rejects a non-conforming patch with a client error and no rollout.
- F3: The applier materializes the validated patch as a values file and applies helm upgrade --reuse-values --atomic, re-rendering the co-generated topology, rather than constructing per-field --set arguments.
- F4: The applier verifies the affected Deployment's rollout via kubectl rollout status and reports the phase with the Deployment's ready and desired replicas and revision, rolling the release back on an apply or verify failure.
- F5: The applier reports health to the creator, which reports it to the provisioning-workflow-orchestrator; a subsequent chat turn is grounded in the applied change.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd006-applier](../software-requirements/srd006-applier.md) [AC1](../software-requirements/srd006-applier.md#ac1) | The applier exposes the state, apply, and rollout surface on the provisioner's Service name and port with a request-scoped apply endpoint. |
| T2 | [srd006-applier](../software-requirements/srd006-applier.md) [AC2](../software-requirements/srd006-applier.md#ac2) | The applier validates the patch against the chart schema by a dry-run and applies it as a values-file rollout, values-plane only. |
| T3 | [srd006-applier](../software-requirements/srd006-applier.md) [AC3](../software-requirements/srd006-applier.md#ac3) | The applier verifies the rollout and rolls back on an apply or verify failure. |
| T4 | [srd006-applier](../software-requirements/srd006-applier.md) [AC4](../software-requirements/srd006-applier.md#ac4) | The actuation surface is authenticated apart from chat, reached only by the creator, and authority stays set by declared validated config. |
| T5 | [srd006-applier](../software-requirements/srd006-applier.md) [AC5](../software-requirements/srd006-applier.md#ac5) | Deciding the values and owning lifecycle are non-goals; the applier receives helm and kubectl from the pinned CLI donor, so no image diverges. |
| T6 | [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC4](../software-requirements/srd003-chatbot-deployment.md#ac4) | The applier realizes the separately-authenticated deployment API that never submits an endpoint to a running agent. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The applier exposes the state, apply, and rollout surface on the provisioner's Service name and port, so the creator's declared client is unchanged. The panel does not call this surface; it submits its intent to the provisioning-workflow-orchestrator, which drives the creator (srd003 R4.1, srd004 R1.5).

<!-- field: traces -->
#### Traces

- [srd006-applier](../software-requirements/srd006-applier.md) [R1.1](../software-requirements/srd006-applier.md#r11)
- [srd006-applier](../software-requirements/srd006-applier.md) [R1.3](../software-requirements/srd006-applier.md#r13)
- [srd006-applier](../software-requirements/srd006-applier.md) [AC1](../software-requirements/srd006-applier.md#ac1)

<!-- sc: S2 -->
### S2

A non-conforming values patch is rejected by the chart-schema dry-run with a client error and no rollout.

<!-- field: traces -->
#### Traces

- [srd006-applier](../software-requirements/srd006-applier.md) [R2.1](../software-requirements/srd006-applier.md#r21)
- [srd006-applier](../software-requirements/srd006-applier.md) [AC2](../software-requirements/srd006-applier.md#ac2)

<!-- sc: S3 -->
### S3

A valid patch is materialized as a values file and applied with helm upgrade --reuse-values --atomic, re-rendering the co-generated topology, with no per-field --set construction.

<!-- field: traces -->
#### Traces

- [srd006-applier](../software-requirements/srd006-applier.md) [R2.2](../software-requirements/srd006-applier.md#r22)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [R2.1](../software-requirements/srd003-chatbot-deployment.md#r21)
- [srd006-applier](../software-requirements/srd006-applier.md) [AC2](../software-requirements/srd006-applier.md#ac2)

<!-- sc: S4 -->
### S4

The rollout is verified via kubectl rollout status and its phase reported, and an apply or verify failure rolls the release back. The live cluster tier forces a deterministic post-upgrade verification failure and proves the compensating helm rollback restores the prior computed release values and a ready Deployment.

<!-- field: traces -->
#### Traces

- [srd006-applier](../software-requirements/srd006-applier.md) [R3.1](../software-requirements/srd006-applier.md#r31)
- [srd006-applier](../software-requirements/srd006-applier.md) [R3.2](../software-requirements/srd006-applier.md#r32)
- [srd006-applier](../software-requirements/srd006-applier.md) [AC3](../software-requirements/srd006-applier.md#ac3)

<!-- sc: S5 -->
### S5

The patch submits no endpoint or credential to a running agent, the apply path is reached only by the creator, and authority stays set by declared validated config.

<!-- field: traces -->
#### Traces

- [srd006-applier](../software-requirements/srd006-applier.md) [R4.1](../software-requirements/srd006-applier.md#r41)
- [srd006-applier](../software-requirements/srd006-applier.md) [R4.2](../software-requirements/srd006-applier.md#r42)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [R4.2](../software-requirements/srd003-chatbot-deployment.md#r42)
- [srd006-applier](../software-requirements/srd006-applier.md) [AC4](../software-requirements/srd006-applier.md#ac4)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC4](../software-requirements/srd003-chatbot-deployment.md#ac4)

<!-- sc: S6 -->
### S6

The rendered chart routes browser provisioning traffic to the provisioning-workflow-orchestrator, not the apply surface, and the applier apply port stays reachable only from creator-labelled pods, so route and NetworkPolicy agree on one authority path.

<!-- field: traces -->
#### Traces

- [srd006-applier](../software-requirements/srd006-applier.md) [R4.4](../software-requirements/srd006-applier.md#r44)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [R4.4](../software-requirements/srd003-chatbot-deployment.md#r44)
- [srd004-provisioning-workflow-orchestrator](../software-requirements/srd004-provisioning-workflow-orchestrator.md) [R1.5](../software-requirements/srd004-provisioning-workflow-orchestrator.md#r15)
- [srd006-applier](../software-requirements/srd006-applier.md) [AC6](../software-requirements/srd006-applier.md#ac6)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC7](../software-requirements/srd003-chatbot-deployment.md#ac7)

<!-- sc: S7 -->
### S7

The rollout read reports the phase with the Deployment's ready and desired replicas and revision, and a read the cluster cannot serve -- an absent Deployment, a wrong kubeconfig, a denied RBAC read, an unreachable API server -- returns a gateway error at every hop rather than a phase, so a broken read is never rendered as a progressing rollout.

<!-- field: traces -->
#### Traces

- [srd006-applier](../software-requirements/srd006-applier.md) [R3.3](../software-requirements/srd006-applier.md#r33)
- [srd006-applier](../software-requirements/srd006-applier.md) [AC3](../software-requirements/srd006-applier.md#ac3)

<!-- field: out_of_scope -->
## Out of scope

- The decision of what to change, which is the provisioning-workflow-orchestrator's (srd004 R3.1).
- Template selection and agent lifecycle, which are the creator's (srd005).
- The chat turn path itself (route, embed, fan-out, compose), which is srd002.

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel06.0-applier](../test-suites/test-rel06.0-applier.md)
