---
id: rel07.0-uc001-collector-agent-tracing
title: Trace the Packaged Topology Through the Canonical Collector Agent
status: implemented
test_suite: test-rel07.0-collector-agent-tracing
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An operator renders and proves the packaged coding topology with the
canonical catalog collector agent as the only trace workload. Roles export
OTLP to the collector Service, the collector spools trace evidence and
serves the bounded query surface, and no Jaeger workload or utility-image
wait chain renders anywhere in the chart.

<!-- field: actor -->
## Actor

An operator or CI worker deploying the packaged coding application.

<!-- field: trigger -->
## Trigger

Render the chart or run mage integration:helmSmoke after mage helmPrepare.

<!-- field: preconditions -->
## Preconditions

- The applications/catalog checkout is discoverable or the demo.yaml catalog_root names it.
- The prepared role packages and the staged collector profile are present in the chart.

<!-- field: flow -->
## Flow

- F1: Preparation stages the canonical catalog collector profile into the chart alongside the role shards.
- F2: The chart renders the collector agent Deployment and Service with OTLP, query, monitor, and control ports and a lifecycle readiness probe.
- F3: Role workloads receive the collector OTLP endpoint and distinct OTel service names; no workload waits on a utility image.
- F4: Values and schema validation accept the collector block and reject unknown collector keys, including the removed utilityImage knob.
- F5: The live smoke drives one traced request and reads the connected trace back through the collector query surface.

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Helm renders the collector agent and no Jaeger manifest under default and CI values.

<!-- field: traces -->
#### Traces

- srd003-helm-topology R3.4
- srd003-helm-topology AC3

<!-- sc: S2 -->
### S2

Every role exports OTLP to the collector with a distinct service name.

<!-- field: traces -->
#### Traces

- srd003-helm-topology R3.3
- srd003-helm-topology AC3

<!-- sc: S3 -->
### S3

The schema fixture matrix proves valid spool default, valid overrides, and invalid telemetry rejection.

<!-- field: traces -->
#### Traces

- srd004-helm-package-validation R1.1
- srd004-helm-package-validation R1.2
- srd004-helm-package-validation R1.3
- srd004-helm-package-validation AC1

<!-- sc: S4 -->
### S4

The live smoke finds all role service names under the submitted trace ID via the collector query surface.

<!-- field: traces -->
#### Traces

- srd005-live-helm-smoke R3.3
- srd005-live-helm-smoke AC1

<!-- sc: S5 -->
### S5

Catalog-root discovery resolves applications/catalog from relative and absolute startup directories without environment overrides.

<!-- field: traces -->
#### Traces

- srd004-helm-package-validation R1.1
- srd004-helm-package-validation R1.2
- srd004-helm-package-validation R1.3
- srd004-helm-package-validation AC1

<!-- field: out_of_scope -->
## Out of scope

- Metrics ingress; the coding chart exports no metrics pipeline.
- Relay mode; the chart collector runs spool-only.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd003-helm-topology R3.4 | S1 traces this requirement in the existing success criteria. |
| T2 | srd003-helm-topology R3.3 | S2 traces this requirement in the existing success criteria. |
| T3 | srd004-helm-package-validation R1.1, R1.2, R1.3 | S3 traces this requirement in the existing success criteria. |
| T4 | srd005-live-helm-smoke R3.3 | S4 traces this requirement in the existing success criteria. |
