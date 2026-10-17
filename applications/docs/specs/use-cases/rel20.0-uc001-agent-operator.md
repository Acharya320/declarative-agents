---
id: rel20.0-uc001-agent-operator
title: Apply an Agent for the Catalog Collector and Watch It Become Ready
status: planned
test_suite: test-rel20.0-agent-operator
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: release -->
## Release

20.0

<!-- field: summary -->
## Summary

An application operator boots the kind platform, which now carries the
agent operator. They create a ConfigMap from the catalog collector's
profile directory and apply one Agent resource naming it. The operator
reconciles the Agent into a Deployment and a Service, the Agent's status
turns Ready, and the collector's health endpoint answers through the
Service. Deleting the Agent removes the Deployment and the Service.

<!-- field: actor -->
## Actor

An application operator using the root Mage platform verbs and kubectl.

<!-- field: trigger -->
## Trigger

The operator wants to run an agent on the platform by declaring it, rather
than by writing its Deployment and Service by hand.

<!-- field: preconditions -->
## Preconditions

- The kind platform satisfies ENG01 and boots with the agent operator installed.
- The agent-core image is built from the checkout and loaded into the platform.
- The catalog collector profile directory is available at applications/catalog/agents/collector.

<!-- field: flow -->
## Flow

- F1: The operator runs `platform:up`. The platform boots with ingress,
  metrics-server, fake-gcs, and the agent operator, whose Deployment rolls
  out and whose CRD is applied.
  - traces:
    - srd009-agent-operator R1.3
    - srd009-agent-operator R4.2
- F2: The operator creates a ConfigMap from the collector's profile
  directory and applies an Agent naming that ConfigMap, the agent-core
  image, the control, monitor, and query ports, and a bind host of every
  interface.
  - traces:
    - srd009-agent-operator R1.1
    - srd009-agent-operator R1.2
- F3: The agent operator creates a Deployment running the image with
  `--profile` at the mounted path and a readiness probe on the control
  port, and a Service for the three ports, both owned by the Agent.
  - traces:
    - srd009-agent-operator R2.1
    - srd009-agent-operator R2.2
    - srd009-agent-operator R2.3
- F4: The Deployment becomes available; the Agent's status reports the
  ready replica and a Ready condition of True; `GET /api/lifecycle/health`
  through the Service answers.
  - traces:
    - srd009-agent-operator R3.1
- F5: The operator changes the Agent's image; the Deployment rolls and the
  status follows the rollout.
  - traces:
    - srd009-agent-operator R2.3
    - srd009-agent-operator R3.1
- F6: The operator deletes the Agent; the Deployment and the Service are
  garbage-collected with it.
  - traces:
    - srd009-agent-operator R2.3

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | docs/engineering/eng01-kind-test-demo-rig.md | The platform the operator installs into. |
| T2 | [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC1](../software-requirements/srd009-agent-operator.md#ac1) | The Agent resource and its schema. |
| T3 | [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC2](../software-requirements/srd009-agent-operator.md#ac2) | Reconciliation and status. |
| T4 | [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC3](../software-requirements/srd009-agent-operator.md#ac3) | Artifacts and platform install. |
| T5 | [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC4](../software-requirements/srd009-agent-operator.md#ac4) | The collector proof. |
| T6 | [test-rel20.0-agent-operator](../test-suites/test-rel20.0-agent-operator.md) | Executable evidence for every success criterion. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The Agent is admitted against the schema the operator applied at boot.

<!-- field: traces -->
#### Traces

- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC1](../software-requirements/srd009-agent-operator.md#ac1)

<!-- sc: S2 -->
### S2

The operator owns one Deployment and one Service per Agent, rolls them on
change, removes them on delete, and reports Ready truthfully.

<!-- field: traces -->
#### Traces

- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC2](../software-requirements/srd009-agent-operator.md#ac2)

<!-- sc: S3 -->
### S3

The platform boots with a ready operator built from the checkout.

<!-- field: traces -->
#### Traces

- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC3](../software-requirements/srd009-agent-operator.md#ac3)

<!-- sc: S4 -->
### S4

The collector answers its health endpoint through the Service the operator
made, on kind and in the release gates.

<!-- field: traces -->
#### Traces

- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC4](../software-requirements/srd009-agent-operator.md#ac4)

<!-- field: out_of_scope -->
## Out of scope

- Autoscaling, task agents as Jobs, and package images.
- Moving the applications' charts onto Agent resources.
- Cloud installation of the operator.
