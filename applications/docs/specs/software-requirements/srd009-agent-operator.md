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

The profile already declares what a chart repeats. Its REST definitions
name every server with an address such as
`${COLLECTOR_BIND_HOST:-127.0.0.1}:${COLLECTOR_CONTROL_PORT:-18191}`, the
server that answers `GET /api/lifecycle/health`, and the server that
serves `/monitor/*`. Release 20.1 makes the profile the deployment spec:
the operator reads it and derives the ports, the control port, the
monitor label, and the bind-host environment, so an Agent names its
profile, its image, and what is genuinely per deployment, and nothing
else.

<!-- field: goals -->
## Goals

- G1: One namespaced Agent resource declares an agent from a profile ConfigMap and the agent-core image.
- G2: The operator owns one Deployment and one Service per Agent and reports readiness in the Agent's status.
- G3: The CRD, the operator image, and its chart ship together and install at platform boot.
- G4: A kind proof with the catalog collector runs as a release gate.
- G5: An Agent repeats nothing its profile declares; the operator reads ports, the control server, the monitor server, and the bind host from the profile.

<!-- req: R1 -->
## R1 - The Agent resource

<!-- req: R1.1 -->
### R1.1

The operator shall define a namespaced custom resource Agent in group
declarative-agents.nokia.com, version v1alpha1, whose spec declares the
profile (a ConfigMap name and the path of the profile file inside it), the
agent-core image, the replica count, container resources, environment
variables, secret references for the environment, arguments appended to
the agent's command line, pod labels, a service account name, and
optional work volumes, each an emptyDir or a named existing
PersistentVolumeClaim. The spec shall not declare ports or a control
port: the profile declares them (R1.2).

<!-- field: evidence -->
#### Evidence

The generated CRD validates each field; an Agent missing the profile or
the image is rejected at admission, and one declaring ports is rejected
as an unknown field.

<!-- req: R1.2 -->
### R1.2

The operator shall read the profile's REST definitions from the
ConfigMap through agent-core's own loader and take from them every
server's port, the control server (the one server answering `GET
/api/lifecycle/health`, the agent's readiness signal), the monitor server
(the one serving `/monitor/` endpoints), and each server's bind and port
variables, expanding `${VAR}` and `${VAR:-default}` with the Agent's
environment over the declaration defaults. A port variable with neither
a default nor an Agent environment entry, two servers with a health
endpoint, or none, is an invalid profile (R2.4).

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
runs the declared image with `--profile` at the mounted profile path
followed by the declared arguments, mounts the profile ConfigMap
read-only, exposes a named container port for every server the profile
declares, probes readiness with `GET /api/lifecycle/health` on the control
server's port, and carries the declared resources, environment, secret
references, service account name, and work volumes. The pod shall carry
the declared pod labels under the operator's own labels, plus
`monitorPort` naming the monitor server's port when the profile has one,
and its environment shall gain `<bind variable>=0.0.0.0` for every server
whose bind variable the Agent's environment does not set. A ConfigMap
key holding `__` shall mount at the path the separator spells
(`units__types.yaml` at `units/types.yaml`), because keys cannot hold a
slash while profiles import sibling directories.

<!-- req: R2.2 -->
### R2.2

For each Agent the operator shall maintain one Service that exposes every
server's port by name and selects the Deployment's pods.

<!-- req: R2.3 -->
### R2.3

The Deployment and the Service shall carry the Agent as their controlling
owner, so deleting the Agent removes both, and a change to the Agent's spec
shall patch both, rolling the Deployment when its pod template changes.

<!-- req: R2.4 -->
### R2.4

When the profile ConfigMap is absent, the profile fails agent-core's
validation or the derivation of R1.2, or the spec cannot be rendered, the
operator shall set the Ready condition to False with a reason that names
the cause (ProfileMissing, ProfileInvalid, InvalidSpec) and the loader's
message, requeue, and create or change nothing: an invalid profile never
reaches a pod.

<!-- req: R2.5 -->
### R2.5

A change to a profile ConfigMap shall reconcile every Agent in its
namespace that names it, so an edited port default or endpoint rolls the
Deployment without a change to the Agent.

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
Agent for the collector naming only that ConfigMap, the agent-core image,
and a work volume, wait for the Ready condition, read
`/api/lifecycle/health` through the Service on the port the profile
declares, delete the Agent, and observe the Deployment and Service go.

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
- Init containers and sidecars in the Agent spec: the charts that use them copy tool binaries or unpack archives, which role images and keyed ConfigMaps replace.
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

<!-- ac: AC5 -->
### AC5

An Agent naming only its profile, image, and work volume yields a
Deployment and Service whose ports, readiness probe, monitor label, and
bind-host environment come from the profile through agent-core's loader;
the declared arguments, pod labels, service account, and claim reach the
pod; and editing a port default in the profile ConfigMap rolls the
Deployment.

<!-- field: traces -->
#### Traces

- [R1.1](#r11)
- [R1.2](#r12)
- [R2.1](#r21)
- [R2.2](#r22)
- [R2.5](#r25)
- [R5.1](#r51)

<!-- ac: AC6 -->
### AC6

An Agent whose profile fails agent-core's validation, or declares no
single server with a health endpoint, reports Ready False with reason
ProfileInvalid and the loader's message, and no Deployment or Service is
created.

<!-- field: traces -->
#### Traces

- [R1.2](#r12)
- [R2.4](#r24)

<!-- field: related_documents -->
## Related documents

- docs/engineering/eng01-kind-test-demo-rig.md
- applications/docs/specs/software-requirements/srd008-persistent-application-platform.md
- applications/docs/specs/software-requirements/srd005-chart-conformance.md
