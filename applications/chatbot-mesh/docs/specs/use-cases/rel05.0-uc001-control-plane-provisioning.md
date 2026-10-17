---
id: rel05.0-uc001-control-plane-provisioning
title: A Provisioning Intent Flows Chatbot to Provisioning Workflow Orchestrator to Creator and Serves
status: partial
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A user asks the UX agent (chatbot, srd002) to add a directory as a knowledge
source. The chatbot delegates the provisioning intent to the provisioning-workflow-orchestrator
(srd004) through a declared client rather than editing any config itself. The
provisioning-workflow-orchestrator sequences the workflow: it requests a corpus-ingest run per directory
from the creator (srd005), verifies the collection counts, decides the values
changes that add the new RAG source, and requests a rollout of the affected
agents. The creator (srd005) realizes each operation -- selecting the profile and
rendered config, spawning the ingest child run, executing the request-draining
rollout through the deployment API (srd003 R3.2), and verifying post-rollout
health -- and reports health back to the provisioning-workflow-orchestrator, which aggregates status for
the provisioning panel. The new source then serves a grounded turn. All
deployment-API authority sits with the creator; no running agent receives an
endpoint or credential (srd003 R4, agent-core srd028 V30).

<!-- field: actor -->
## Actor

A user adding a knowledge source through the UX agent.

<!-- field: trigger -->
## Trigger

The user tells the chatbot to add a directory as a knowledge source; the chatbot
delegates the intent to the provisioning-workflow-orchestrator.

<!-- field: preconditions -->
## Preconditions

- The chatbot, provisioning-workflow-orchestrator, and creator profiles are deployed with declared clients between them.
- The provisioning-workflow-orchestrator and creator authenticate apart from the chat path; the creator alone holds the deployment-API credentials.
- The deployment API (srd003 R4) and the chart's values-to-config co-generation are available.

<!-- field: flow -->
## Flow

- F1: The chatbot delegates the user's provisioning intent to the provisioning-workflow-orchestrator through a declared client; it submits no host, URL, or credential.
- F2: The provisioning-workflow-orchestrator requests a corpus-ingest run per directory from the creator and verifies each ingested collection's count before proceeding.
- F3: The provisioning-workflow-orchestrator decides the values changes that add the new RAG source and requests a rollout of the affected agents from the creator.
- F4: The creator selects each instance's profile and rendered config, spawns the ingest child run, executes the request-draining rollout through the deployment API, and verifies post-rollout health.
- F5: The creator reports health to the provisioning-workflow-orchestrator, which aggregates status for the provisioning panel; a subsequent chat turn is grounded in the newly added source.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd004-provisioning-workflow-orchestrator](../software-requirements/srd004-provisioning-workflow-orchestrator.md) [AC1](../software-requirements/srd004-provisioning-workflow-orchestrator.md#ac1) | The provisioning-workflow-orchestrator binds the UX agent's intent to a workflow run and aggregates status. |
| T2 | [srd004-provisioning-workflow-orchestrator](../software-requirements/srd004-provisioning-workflow-orchestrator.md) [AC2](../software-requirements/srd004-provisioning-workflow-orchestrator.md#ac2) | The provisioning-workflow-orchestrator sequences the ingest per directory through the creator and verifies counts. |
| T3 | [srd004-provisioning-workflow-orchestrator](../software-requirements/srd004-provisioning-workflow-orchestrator.md) [AC3](../software-requirements/srd004-provisioning-workflow-orchestrator.md#ac3) | The provisioning-workflow-orchestrator is the sole author of the values decision and requests the rollout. |
| T4 | [srd005-creator](../software-requirements/srd005-creator.md) [AC3](../software-requirements/srd005-creator.md#ac3) | The creator spawns the ingest run and executes the request-draining rollout through the deployment API. |
| T5 | [srd005-creator](../software-requirements/srd005-creator.md) [AC5](../software-requirements/srd005-creator.md#ac5) | The creator alone holds the deployment-API credentials and passes no endpoint into a running agent. |
| T6 | test-rel05.0-control-plane provisioning_intent | The suite records the intent -> ingest -> reconfigure -> restart -> serve flow across the three agents. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The chatbot delegates the intent to the provisioning-workflow-orchestrator through a declared client and submits no endpoint or credential to any agent.

<!-- field: traces -->
#### Traces

- [srd004-provisioning-workflow-orchestrator](../software-requirements/srd004-provisioning-workflow-orchestrator.md) [R1.1](../software-requirements/srd004-provisioning-workflow-orchestrator.md#r11)
- [srd004-provisioning-workflow-orchestrator](../software-requirements/srd004-provisioning-workflow-orchestrator.md) [R1.2](../software-requirements/srd004-provisioning-workflow-orchestrator.md#r12)
- [srd004-provisioning-workflow-orchestrator](../software-requirements/srd004-provisioning-workflow-orchestrator.md) [AC1](../software-requirements/srd004-provisioning-workflow-orchestrator.md#ac1)

<!-- sc: S2 -->
### S2

The provisioning-workflow-orchestrator sequences a corpus-ingest run per directory through the creator and verifies collection counts before reconfiguring.

<!-- field: traces -->
#### Traces

- [srd004-provisioning-workflow-orchestrator](../software-requirements/srd004-provisioning-workflow-orchestrator.md) [R2.1](../software-requirements/srd004-provisioning-workflow-orchestrator.md#r21)
- [srd004-provisioning-workflow-orchestrator](../software-requirements/srd004-provisioning-workflow-orchestrator.md) [R2.2](../software-requirements/srd004-provisioning-workflow-orchestrator.md#r22)
- [srd004-provisioning-workflow-orchestrator](../software-requirements/srd004-provisioning-workflow-orchestrator.md) [AC2](../software-requirements/srd004-provisioning-workflow-orchestrator.md#ac2)

<!-- sc: S3 -->
### S3

The provisioning-workflow-orchestrator decides the values changes and requests a rollout; the creator executes the request-draining rollout through the deployment API and verifies health.

<!-- field: traces -->
#### Traces

- [srd004-provisioning-workflow-orchestrator](../software-requirements/srd004-provisioning-workflow-orchestrator.md) [R3.1](../software-requirements/srd004-provisioning-workflow-orchestrator.md#r31)
- [srd005-creator](../software-requirements/srd005-creator.md) [R3.2](../software-requirements/srd005-creator.md#r32)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [R3.2](../software-requirements/srd003-chatbot-deployment.md#r32)
- [srd004-provisioning-workflow-orchestrator](../software-requirements/srd004-provisioning-workflow-orchestrator.md) [AC3](../software-requirements/srd004-provisioning-workflow-orchestrator.md#ac3)
- [srd005-creator](../software-requirements/srd005-creator.md) [AC3](../software-requirements/srd005-creator.md#ac3)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC3](../software-requirements/srd003-chatbot-deployment.md#ac3)

<!-- sc: S4 -->
### S4

The creator alone holds the deployment-API authority, both control-plane agents authenticate apart from chat, and no running agent receives an endpoint or credential.

<!-- field: traces -->
#### Traces

- [srd005-creator](../software-requirements/srd005-creator.md) [R5.1](../software-requirements/srd005-creator.md#r51)
- [srd005-creator](../software-requirements/srd005-creator.md) [R5.3](../software-requirements/srd005-creator.md#r53)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [R4.2](../software-requirements/srd003-chatbot-deployment.md#r42)
- [srd005-creator](../software-requirements/srd005-creator.md) [AC5](../software-requirements/srd005-creator.md#ac5)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC4](../software-requirements/srd003-chatbot-deployment.md#ac4)

<!-- sc: S5 -->
### S5

After the rollout the newly added source serves a grounded chat turn, and the provisioning-workflow-orchestrator's aggregated status renders in the provisioning panel.

<!-- field: traces -->
#### Traces

- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [R2.1](../software-requirements/srd003-chatbot-deployment.md#r21)
- [srd004-provisioning-workflow-orchestrator](../software-requirements/srd004-provisioning-workflow-orchestrator.md) [R1.3](../software-requirements/srd004-provisioning-workflow-orchestrator.md#r13)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC2](../software-requirements/srd003-chatbot-deployment.md#ac2)
- [srd004-provisioning-workflow-orchestrator](../software-requirements/srd004-provisioning-workflow-orchestrator.md) [AC1](../software-requirements/srd004-provisioning-workflow-orchestrator.md#ac1)

<!-- field: out_of_scope -->
## Out of scope

- The chat turn path itself (route, embed, fan-out, compose), which is srd002 and releases rel00.0/rel01.0.
- In-process client-registry mutation of a running agent; reconfiguration is a rollout (srd003 R3.3).
- The mesh deployment topology and tiered swap themselves, which are srd003 and release rel03.0; this use case adds the control-plane flow on top of them.

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel05.0-control-plane](../test-suites/test-rel05.0-control-plane.md)
