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
| T1 | [srd007-collector-agent](../software-requirements/srd007-collector-agent.md) [R1.1](../software-requirements/srd007-collector-agent.md#r11), [R1.2](../software-requirements/srd007-collector-agent.md#r12), [R1.3](../software-requirements/srd007-collector-agent.md#r13), [R1.4](../software-requirements/srd007-collector-agent.md#r14) | persistent bounded intake with query surface |
| T2 | [srd007-collector-agent](../software-requirements/srd007-collector-agent.md) [R2.1](../software-requirements/srd007-collector-agent.md#r21), [R2.2](../software-requirements/srd007-collector-agent.md#r22), [R2.3](../software-requirements/srd007-collector-agent.md#r23), [R2.4](../software-requirements/srd007-collector-agent.md#r24) | complete integration identity in spool |
| T3 | [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [R5.1](../software-requirements/srd003-chatbot-deployment.md#r51), [R5.2](../software-requirements/srd003-chatbot-deployment.md#r52), [R5.3](../software-requirements/srd003-chatbot-deployment.md#r53), [R5.4](../software-requirements/srd003-chatbot-deployment.md#r54), [R5.5](../software-requirements/srd003-chatbot-deployment.md#r55), [R5.6](../software-requirements/srd003-chatbot-deployment.md#r56), [R5.7](../software-requirements/srd003-chatbot-deployment.md#r57), [R5.8](../software-requirements/srd003-chatbot-deployment.md#r58), [R5.9](../software-requirements/srd003-chatbot-deployment.md#r59), [R5.10](../software-requirements/srd003-chatbot-deployment.md#r510), [R5.11](../software-requirements/srd003-chatbot-deployment.md#r511) | trace evidence through collector query surface |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Helm render contains no Jaeger template, values reference, or schema entry.

<!-- field: traces -->
#### Traces

- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC5](../software-requirements/srd003-chatbot-deployment.md#ac5)

<!-- sc: S2 -->
### S2

The collector query port 18193 is exposed in the collector Service.

<!-- field: traces -->
#### Traces

- [srd007-collector-agent](../software-requirements/srd007-collector-agent.md) [AC1](../software-requirements/srd007-collector-agent.md#ac1)

<!-- sc: S3 -->
### S3

The UX waterfall transforms collector query format (ISO timestamps, flat spans) to the internal trace model.

<!-- field: traces -->
#### Traces

- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC5](../software-requirements/srd003-chatbot-deployment.md#ac5)

<!-- sc: S4 -->
### S4

Spool identity assertions pass without any Jaeger endpoint.

<!-- field: traces -->
#### Traces

- [srd007-collector-agent](../software-requirements/srd007-collector-agent.md) [AC2](../software-requirements/srd007-collector-agent.md#ac2)

<!-- sc: S5 -->
### S5

Metrics evidence is collected from Prometheus and retained after teardown.

<!-- field: traces -->
#### Traces

- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC5](../software-requirements/srd003-chatbot-deployment.md#ac5)

<!-- sc: S6 -->
### S6

The collector program is consumed from the catalog, not from a local directory.

<!-- field: traces -->
#### Traces

- [srd007-collector-agent](../software-requirements/srd007-collector-agent.md) [AC1](../software-requirements/srd007-collector-agent.md#ac1)

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

[test-rel09.0-canonical-collector](../test-suites/test-rel09.0-canonical-collector.md)
