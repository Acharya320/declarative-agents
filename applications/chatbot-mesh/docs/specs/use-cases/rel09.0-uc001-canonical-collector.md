---
id: rel09.0-uc001-canonical-collector
title: The mesh deploys with no Jaeger workload and traces from the collector query surface
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An operator installs the mesh chart with default values. No Jaeger
Deployment or Service renders. The collector agent spools traces and serves
a REST query surface on port 18193. The UX waterfall reads traces from the
collector query surface through the chatbot monitor proxy instead of from
Jaeger. Collector mode derives from externalOTLPEndpoint presence: relay
when set, spool otherwise. The collector program is consumed from the
catalog rather than carried locally. Integration tests verify spool
identity and metrics evidence without Jaeger cross-checks.

<!-- field: actor -->
## Actor

Mesh integration operator

<!-- field: trigger -->
## Trigger

A default Helm install renders the collector agent with no Jaeger workload and the UX reads traces from the collector query surface.

<!-- field: flow -->
## Flow

- F1: Helm renders the collector agent with query port 18193 exposed and no Jaeger template, values, or schema entry.
- F2: The chatbot rest.yaml monitor proxy routes /monitor-proxy/collector/ to the collector query surface.
- F3: The UX trace backend reads traces from /monitor-proxy/collector/query/traces/{trace_id} and transforms collector format to the waterfall model.
- F4: Integration tests assert spool identity by reading the collector spool directly rather than cross-checking Jaeger.
- F5: Integration tests assert metrics evidence through Prometheus without Jaeger trace evidence.
- F6: The collector program is sourced from the catalog root, not from a local agents/collector directory.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd007-collector-agent R1.1, R1.2, R1.3, R1.4 | persistent bounded intake with query surface |
| T2 | srd007-collector-agent R2.1, R2.2, R2.3, R2.4 | complete integration identity in spool |
| T3 | srd003-chatbot-deployment R5.1, R5.2, R5.3, R5.4, R5.5, R5.6, R5.7, R5.8, R5.9, R5.10, R5.11 | trace evidence through collector query surface |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Helm render contains no Jaeger template, values reference, or schema entry.

<!-- field: traces -->
#### Traces

- srd003-chatbot-deployment AC5

<!-- sc: S2 -->
### S2

The collector query port 18193 is exposed in the collector Service.

<!-- field: traces -->
#### Traces

- srd007-collector-agent AC1

<!-- sc: S3 -->
### S3

The UX waterfall transforms collector query format (ISO timestamps, flat spans) to the internal trace model.

<!-- field: traces -->
#### Traces

- srd003-chatbot-deployment AC5

<!-- sc: S4 -->
### S4

Spool identity assertions pass without any Jaeger endpoint.

<!-- field: traces -->
#### Traces

- srd007-collector-agent AC2

<!-- sc: S5 -->
### S5

Metrics evidence is collected from Prometheus and retained after teardown.

<!-- field: traces -->
#### Traces

- srd003-chatbot-deployment AC5

<!-- sc: S6 -->
### S6

The collector program is consumed from the catalog, not from a local directory.

<!-- field: traces -->
#### Traces

- srd007-collector-agent AC1

<!-- field: failure_scenarios -->
## Failure scenarios

| id | scenario | expected |
| --- | --- | --- |
| E1 | Helm render produces a Jaeger Deployment or Service. | Rejected; the Jaeger template is deleted and values block removed. |
| E2 | The UX waterfall fails to parse collector query format timestamps. | Rejected; isoToUs converts ISO 8601 strings to microsecond integers. |
| E3 | Integration tests depend on a Jaeger endpoint for evidence. | Rejected; all trace evidence reads from the collector spool directly. |

<!-- field: out_of_scope -->
## Out of scope

- Collector-agent metrics and logs intake.
- External OTLP relay endpoint configuration.
- Observer agent fleet discovery.

<!-- field: test_suite_trace -->
## Test suite trace

test-rel09.0-canonical-collector
