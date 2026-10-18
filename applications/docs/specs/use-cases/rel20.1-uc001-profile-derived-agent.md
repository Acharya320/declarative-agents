---
id: rel20.1-uc001-profile-derived-agent
title: Declare an Agent from Its Profile Alone
status: implemented
test_suite: test-rel20.1-profile-derived-agent
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: release -->
## Release

20.1

<!-- field: summary -->
## Summary

An application operator applies an Agent for the catalog collector that
names only its profile ConfigMap, the agent-core image, and a work volume.
The agent operator reads the profile's REST definitions and derives the
container and Service ports, the readiness probe on the control server,
the monitor port label, and the bind-host environment from them. Editing a
port default in the ConfigMap rolls the Deployment; a profile without a
health server is reported as invalid before any pod exists.

<!-- field: actor -->
## Actor

An application operator using kubectl against the kind platform.

<!-- field: trigger -->
## Trigger

The operator wants an Agent to repeat nothing its profile already declares,
so that one change to a profile is one change on the cluster.

<!-- field: preconditions -->
## Preconditions

- The kind platform satisfies ENG01 and boots with the agent operator installed.
- The agent-core image is built from the checkout and loaded into the platform.
- The catalog collector profile and the units file it imports are available under applications/catalog/agents.

<!-- field: flow -->
## Flow

- F1: The operator creates a ConfigMap from the collector's profile
  directory and the units file it imports, keyed `<dir>__<file>`, and
  applies an Agent naming that ConfigMap, the agent-core image, and a
  work volume, with no ports and no control port.
  - traces:
    - srd009-agent-operator R1.1
    - srd009-agent-operator R1.2
- F2: The agent operator reads the profile's REST definitions through
  agent-core's loader, finds the control, monitor, and query servers with
  their default ports, and creates a Deployment whose container exposes
  the three ports, probes `/api/lifecycle/health` on the control port,
  carries the `monitorPort` label, and sets `COLLECTOR_BIND_HOST` to every
  interface, and a Service for the three ports, both owned by the Agent.
  - traces:
    - srd009-agent-operator R2.1
    - srd009-agent-operator R2.2
    - srd009-agent-operator R2.3
- F3: The Deployment becomes available; the Agent reports Ready and
  `GET /api/lifecycle/health` answers through the Service.
  - traces:
    - srd009-agent-operator R3.1
- F4: The operator edits the ConfigMap, changing the control port's
  default; the agent operator reconciles the Agent, the Deployment rolls
  with the new port and probe, and the Service follows.
  - traces:
    - srd009-agent-operator R2.5
    - srd009-agent-operator R2.3
- F5: The operator applies an Agent whose profile declares no server with
  a health endpoint; the Agent reports Ready False with reason
  ProfileInvalid and the loader's message, and no Deployment or Service
  is created.
  - traces:
    - srd009-agent-operator R2.4

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC5](../software-requirements/srd009-agent-operator.md#ac5) | Derivation from the profile and reconciliation on a profile change. |
| T2 | [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC6](../software-requirements/srd009-agent-operator.md#ac6) | An invalid profile reported before any pod exists. |
| T3 | [rel20.0-uc001-agent-operator](rel20.0-uc001-agent-operator.md) | The release 20.0 walk this one shortens: the Agent no longer declares what the profile declares. |
| T4 | [test-rel20.1-profile-derived-agent](../test-suites/test-rel20.1-profile-derived-agent.md) | Executable evidence for every success criterion. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The ports, the readiness probe, the monitor label, and the bind-host
environment of the Deployment and Service come from the profile, read
through agent-core's own loader.

<!-- field: traces -->
#### Traces

- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC5](../software-requirements/srd009-agent-operator.md#ac5)

<!-- sc: S2 -->
### S2

A profile that fails validation or has no health server is reported with
reason ProfileInvalid and creates nothing, and a profile change rolls the
Deployment.

<!-- field: traces -->
#### Traces

- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC6](../software-requirements/srd009-agent-operator.md#ac6)
- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC5](../software-requirements/srd009-agent-operator.md#ac5)

<!-- sc: S3 -->
### S3

The collector proof passes with an Agent of profile, image, and work
volume only, on kind and in the release gates.

<!-- field: traces -->
#### Traces

- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC5](../software-requirements/srd009-agent-operator.md#ac5)
- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC4](../software-requirements/srd009-agent-operator.md#ac4)

<!-- field: out_of_scope -->
## Out of scope

- Moving the applications' charts onto Agent resources.
- Init containers and sidecars in the Agent spec.
- Autoscaling, task agents as Jobs, package images, and cloud installation.
