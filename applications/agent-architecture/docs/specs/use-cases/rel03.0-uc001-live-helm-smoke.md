---
id: rel03.0-uc001-live-helm-smoke
title: Prove the Packaged Agent Architecture Chart on Kind
status: planned
test_suite: test-rel03.0-live-helm-smoke
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A developer installs only the validated chart archive into a bounded kind
cluster, builds and loads the shared profile-free agent-core image the chart
references for both workloads, mounts the packaged curator and collector
closures, and provisions the curator UI shard ConfigMaps out-of-release so the
curator unpacks them. The
curator and collector Deployments and their real lifecycle health endpoints
become ready, the curator serves its documentation interface, a declarative
lifecycle-exit request stops the curator cleanly, and the collector retains a
trace the curator exported. The target skips with a named missing prerequisite
when its toolchain is unavailable and cleans up on every terminal path.

<!-- field: actor -->
## Actor

A developer proving the packaged Agent Architecture topology on kind.

<!-- field: trigger -->
## Trigger

Run the mage kind smoke target.

<!-- field: preconditions -->
## Preconditions

- Docker, kind, Helm, and kubectl are available and Docker is healthy, or the target skips.
- The validated chart archive and the catalog and agent-core checkouts are available.

<!-- field: flow -->
## Flow

- F1: The target checks its prerequisites and skips with a named missing tool when Docker, kind, Helm, kubectl, a required checkout, or Docker health is unavailable.
- F2: The target creates a kind cluster through a temporary kubeconfig, builds and loads the shared profile-free agent-core image, and installs only the validated archive.
- F3: The curator and collector Deployments and their real lifecycle health endpoints become ready.
- F4: The curator serves its documentation interface, and a declarative lifecycle-exit request stops the curator cleanly.
- F5: The collector retains a trace the curator exported under the curator service name.
- F6: The target cleans its namespace, release, forwards, temporary files, and any owned cluster on every terminal path, and failure diagnostics share one deadline and include Helm status, workloads, events, and targeted curator and collector logs.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [AC4](../software-requirements/srd001-helm-deployment.md#ac4) | The bounded kind install of the validated archive serves and exits cleanly. |
| T2 | [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [R3](../software-requirements/srd001-helm-deployment.md#r3---live-kind-smoke) | Prerequisite skips, mounted closures, readiness, trace retention, and cleanup follow the live-smoke requirements. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The target installs only the validated archive, builds and loads the shared profile-free agent-core image the chart references, mounts the packaged curator and collector closures, and provisions the curator UI shard ConfigMaps out-of-release for the curator to unpack.

<!-- field: traces -->
#### Traces

- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [R3.2](../software-requirements/srd001-helm-deployment.md#r32)
- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [AC4](../software-requirements/srd001-helm-deployment.md#ac4)

<!-- sc: S2 -->
### S2

The curator and collector Deployments and their real lifecycle health endpoints become ready, the curator serves its documentation interface, and a declarative lifecycle-exit request stops the curator cleanly.

<!-- field: traces -->
#### Traces

- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [R3.3](../software-requirements/srd001-helm-deployment.md#r33)
- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [AC4](../software-requirements/srd001-helm-deployment.md#ac4)

<!-- sc: S3 -->
### S3

The collector retains a trace the curator exported under the curator service name, proving the in-cluster trace path end to end.

<!-- field: traces -->
#### Traces

- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [R3.4](../software-requirements/srd001-helm-deployment.md#r34)
- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [AC4](../software-requirements/srd001-helm-deployment.md#ac4)

<!-- sc: S4 -->
### S4

The target skips with a named missing prerequisite without mutating a cluster, and cleans its owned resources on every terminal path with bounded failure diagnostics.

<!-- field: traces -->
#### Traces

- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [R3.1](../software-requirements/srd001-helm-deployment.md#r31)
- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [R3.5](../software-requirements/srd001-helm-deployment.md#r35)
- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [AC4](../software-requirements/srd001-helm-deployment.md#ac4)

<!-- field: out_of_scope -->
## Out of scope

- The applier actuation surface and rollback.
- Multi-node or production operational evidence.
- Changing the canonical curator or collector profiles.
