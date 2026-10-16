---
id: rel04.0-uc001-render-helm-topology
title: Render the Persistent Coding-Agent Topology
status: implemented
test_suite: test-rel04.0-helm-topology
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An operator prepares #875 role packages and renders a chart containing one
persistent planner, executor, and critic process. The roles share workspace
storage, communicate over fixed Services, report their real lifecycle health,
and export connected traces through the chart observability tier.

<!-- field: actor -->
## Actor

An operator preparing a coding-agent Kubernetes release.

<!-- field: trigger -->
## Trigger

Run mage helmPrepare followed by helm lint or helm template.

<!-- field: preconditions -->
## Preconditions

- The #875 package can resolve all serving role closures.
- Helm is available.

<!-- field: flow -->
## Flow

- F1: helmPrepare regenerates and stages deployment and role package manifests and profile files.
- F2: Helm renders one or more ConfigMaps per role from the manifest partitions.
- F3: Each role Deployment projects only its ConfigMaps and selects its packaged serving profile.
- F4: Every role mounts the shared workspace claim and exposes request and control ports through a ClusterIP Service.
- F5: Planner resolves executor and critic through fixed Service URLs; W3C trace context flows through the serving contract.
- F6: Role OTLP exporters send to the collector agent, which spools traces and serves a query surface; model clients use the selected Ollama endpoint.

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The rendered topology contains exactly one real workload for each application role.

<!-- field: traces -->
#### Traces

- srd003-helm-topology R1.1
- srd003-helm-topology R1.2
- srd003-helm-topology R1.3
- srd003-helm-topology R1.4
- srd003-helm-topology R1.5
- srd003-helm-topology AC1

<!-- sc: S2 -->
### S2

Role profile changes produce role-scoped rollout checksums.

<!-- field: traces -->
#### Traces

- srd003-helm-topology AC2

<!-- sc: S3 -->
### S3

Default, CI, and optional in-cluster Ollama renders are coherent.

<!-- field: traces -->
#### Traces

- srd003-helm-topology AC3
- srd003-helm-topology AC4

<!-- field: out_of_scope -->
## Out of scope

- values.schema.json and negative schema fixtures.
- Packaged-chart archive testing.
- Live kind installation and deployment documentation.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd003-helm-topology R1.1, R1.2, R1.3, R1.4, R1.5 | S1 traces this requirement in the existing success criteria. |
| T2 | srd003-helm-topology AC1 | S1 traces this requirement in the existing success criteria. |
| T3 | srd003-helm-topology AC2 | S2 traces this requirement in the existing success criteria. |
| T4 | srd003-helm-topology AC3 | S3 traces this requirement in the existing success criteria. |
| T5 | srd003-helm-topology AC4 | S3 traces this requirement in the existing success criteria. |
