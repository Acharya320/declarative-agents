---
id: rel02.0-uc001-render-packaged-chart
title: Render and Validate the Packaged Agent Architecture Chart
status: planned
test_suite: test-rel02.0-render-packaged-chart
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An operator prepares the catalog-owned curator and collector profile closures
and renders a chart that runs the documentation-curator and the collector as
persistent workloads. The workloads share the profile-free agent-core image
(GH-1368), mount only their prepared closures, report their real lifecycle health, and
wire the curator's OTLP export to the in-cluster collector. A values schema
and template guards reject unsafe or incoherent values, and a self-contained
archive lints and renders with no source checkout present.

<!-- field: actor -->
## Actor

An operator preparing an Agent Architecture Kubernetes release.

<!-- field: trigger -->
## Trigger

Run the mage preparation target followed by helm lint, helm template, or the package target.

<!-- field: preconditions -->
## Preconditions

- The demo.yaml catalog_root resolves to a catalog checkout containing the documentation-curator and collector profiles.
- Helm is available.

<!-- field: flow -->
## Flow

- F1: The preparation target regenerates the catalog-owned curator and collector profile closures and atomically stages them as the chart's only profile source.
- F2: Helm renders one curator and one collector Deployment, each projecting only its prepared ConfigMaps and selecting its mounted canonical profile.
- F3: Each workload exposes its request, control, and monitor or query ports through a ClusterIP Service, and its readiness and liveness probes call its real lifecycle health endpoint.
- F4: The curator receives a values-selected OTLP endpoint and service name and exports traces to the collector, which spools them and serves its query surface.
- F5: The values schema rejects unknown keys and malformed values, template guards retain mount, fixed-port, and telemetry invariants, and the package target produces one self-contained archive.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [AC1](../software-requirements/srd001-helm-deployment.md#ac1) | The rendered topology runs the curator and collector as truthful workloads from mounted closures. |
| T2 | [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [AC2](../software-requirements/srd001-helm-deployment.md#ac2) | Preparation stages and checksum-validates the catalog-owned closures. |
| T3 | [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [AC3](../software-requirements/srd001-helm-deployment.md#ac3) | The schema, guards, and self-contained archive hold. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The rendered topology contains exactly one curator and one collector workload with the profile-free image, mounted closures, fixed Services, truthful lifecycle probes, and hardened pod security.

<!-- field: traces -->
#### Traces

- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [R1](../software-requirements/srd001-helm-deployment.md#r1---deployed-topology)
- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [AC1](../software-requirements/srd001-helm-deployment.md#ac1)

<!-- sc: S2 -->
### S2

Preparation regenerates and checksum-validates the catalog-owned curator and collector closures and rejects stale, tampered, or out-of-bound package inputs before staging.

<!-- field: traces -->
#### Traces

- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [R2.1](../software-requirements/srd001-helm-deployment.md#r21)
- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [R2.2](../software-requirements/srd001-helm-deployment.md#r22)
- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [AC2](../software-requirements/srd001-helm-deployment.md#ac2)

<!-- sc: S3 -->
### S3

The values schema and template guards accept every valid fixture and reject each invalid fixture by its offending field, and the self-contained archive lints and renders without a source checkout.

<!-- field: traces -->
#### Traces

- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [R2.3](../software-requirements/srd001-helm-deployment.md#r23)
- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [R2.4](../software-requirements/srd001-helm-deployment.md#r24)
- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [R2.5](../software-requirements/srd001-helm-deployment.md#r25)
- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [AC3](../software-requirements/srd001-helm-deployment.md#ac3)

<!-- sc: S4 -->
### S4

The curator's OTLP export to the in-cluster collector renders coherently under default and CI values.

<!-- field: traces -->
#### Traces

- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [R1.5](../software-requirements/srd001-helm-deployment.md#r15)
- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [AC1](../software-requirements/srd001-helm-deployment.md#ac1)

<!-- field: out_of_scope -->
## Out of scope

- Live kind installation and on-cluster behavior.
- The applier actuation surface.
- Changing the canonical curator or collector profiles.
