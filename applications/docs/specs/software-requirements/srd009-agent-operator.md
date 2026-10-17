---
id: srd009-agent-operator
title: The Agent Operator and the Agent Resource
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Every agent the applications run on the platform is deployed the same way,
by hand, in each chart: a Deployment that runs the agent-core image with
`--profile` pointing into a mounted profiles ConfigMap, a readiness probe on
`/api/lifecycle/health`, a Service for the ports the profile serves, and
volumes for its work. Each chart repeats that wiring, and nothing on the
platform owns the lifecycle of an agent as a unit. There is no resource to
apply, watch, or scale, and no status that says an agent is ready.

An operator that reconciles an Agent resource into exactly that Deployment
and Service gives the platform one owner for agent workloads without
changing how agent-core runs. This document fixes the resource, the
reconciliation it implies, the status it reports, how the operator ships
and installs, and the proof that it works on the kind platform.

<!-- field: goals -->
## Goals

- G1: One namespaced Agent resource declares an agent from a profile ConfigMap and the agent-core image.
- G2: The operator owns one Deployment and one Service per Agent and reports readiness in the Agent's status.
- G3: The CRD, the operator image, and its chart ship together and install at platform boot.
- G4: A kind proof with the catalog collector runs as a release gate.

<!-- req: R1 -->
## R1 - The Agent resource

<!-- req: R1.1 -->
### R1.1

The operator shall define a namespaced custom resource Agent in group
declarative-agents.nokia.com, version v1alpha1, whose spec declares the
profile (a ConfigMap name and the path of the profile file inside it), the
agent-core image, the ports the profile serves (each a name and a number),
the name of the control port, the replica count, container resources,
environment variables, secret references for the environment, and
optional work volumes.

<!-- field: evidence -->
#### Evidence

The generated CRD validates each field; an Agent missing the profile or
the control port is rejected at admission.

<!-- req: R1.2 -->
### R1.2

The control port named in the spec shall be the one the profile's REST
server answers `GET /api/lifecycle/health` on; the operator shall treat that
endpoint as the agent's readiness signal.

<!-- req: R1.3 -->
### R1.3

The operator binary shall carry its CustomResourceDefinition and apply it,
creating or updating, before it starts reconciling, so the schema and the
binary ship as one unit and no installer can leave them apart.

<!-- req: R2 -->
## R2 - Reconciliation

<!-- req: R2.1 -->
### R2.1

For each Agent the operator shall maintain one Deployment whose container
runs the declared image with `--profile` at the mounted profile path, mounts
the profile ConfigMap read-only, probes readiness with `GET
/api/lifecycle/health` on the control port, and carries the declared
resources, environment, secret references, and work volumes. A ConfigMap
key holding `__` shall mount at the path the separator spells
(`units__types.yaml` at `units/types.yaml`), because keys cannot hold a
slash while profiles import sibling directories.

<!-- req: R2.2 -->
### R2.2

For each Agent the operator shall maintain one Service that exposes the
declared ports by name and selects the Deployment's pods.

<!-- req: R2.3 -->
### R2.3

The Deployment and the Service shall carry the Agent as their controlling
owner, so deleting the Agent removes both, and a change to the Agent's spec
shall patch both, rolling the Deployment when its pod template changes.

<!-- req: R2.4 -->
### R2.4

When the profile ConfigMap is absent or the spec cannot be rendered, the
operator shall set the Ready condition to False with a reason that names
the cause, requeue, and change nothing it does not own.

<!-- req: R3 -->
## R3 - Status

<!-- req: R3.1 -->
### R3.1

The Agent status shall carry the observed generation, the ready replica
count, and a Ready condition, kept current from the Deployment; Ready is
True only when the available replicas equal the desired count and the
observed generation equals the Agent's generation.

<!-- req: R4 -->
## R4 - Artifacts and platform install

<!-- req: R4.1 -->
### R4.1

The operator image shall build from the checkout as
localhost/declarative-agents/runtime/agent-operator under the first-party
naming chart conformance enforces, and a Helm chart in the operator module
shall install it with a namespace, a service account, and RBAC limited to
Agents and their status, Deployments, Services, CustomResourceDefinitions,
and events.

<!-- req: R4.2 -->
### R4.2

The kind platform boot shall load the operator image and install the chart
after fake-gcs, waiting for its rollout and logging the same phase records
as the other platform components, so a booted platform always carries a
ready operator.

<!-- req: R5 -->
## R5 - Proof

<!-- req: R5.1 -->
### R5.1

A kind proof shall create a ConfigMap from applications/catalog/agents/
collector and the units file it imports, keyed as directories, apply an
Agent for the collector with its control, monitor, and
query ports and the bind host set to every interface, wait for the Ready
condition, read `/api/lifecycle/health` through the Service, delete the
Agent, and observe the Deployment and Service go.

<!-- req: R5.2 -->
### R5.2

The proof shall run under `mage integration:all` in the operator module and
as its own gate in the release gate table, on the platform lane.

<!-- field: non_goals -->
## Non-goals

- Autoscaling: the replica count is fixed in this release.
- Task agents run as Jobs.
- Package images and image volumes as a profile source.
- Moving the applications' charts onto Agent resources.
- Installing the operator on a cloud platform.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

An Agent applied to the kind platform is admitted against the generated
schema, and the operator applied that schema itself at boot.

<!-- field: traces -->
#### Traces

- [R1.1](#r11)
- [R1.2](#r12)
- [R1.3](#r13)

<!-- ac: AC2 -->
### AC2

The operator creates the Deployment and Service for an Agent, rolls them on
a spec change, removes them with the Agent, reports a False Ready
condition with a reason when the profile is missing, and reports Ready
True once the Deployment is available.

<!-- field: traces -->
#### Traces

- [R2.1](#r21)
- [R2.2](#r22)
- [R2.3](#r23)
- [R2.4](#r24)
- [R3.1](#r31)

<!-- ac: AC3 -->
### AC3

The operator image and chart build from the checkout under first-party
naming, and the kind platform boots with the operator ready.

<!-- field: traces -->
#### Traces

- [R4.1](#r41)
- [R4.2](#r42)

<!-- ac: AC4 -->
### AC4

The collector proof passes on the kind platform and in a release gate run.

<!-- field: traces -->
#### Traces

- [R5.1](#r51)
- [R5.2](#r52)

<!-- field: related_documents -->
## Related documents

- docs/engineering/eng01-kind-test-demo-rig.md
- applications/docs/specs/software-requirements/srd008-persistent-application-platform.md
- applications/docs/specs/software-requirements/srd005-chart-conformance.md
