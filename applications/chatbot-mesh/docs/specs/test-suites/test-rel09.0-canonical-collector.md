---
id: test-rel09.0-canonical-collector
title: Test Suite for Release 09.0 Canonical Collector and Jaeger Retirement
release: "09.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- [rel09.0-uc001-canonical-collector](../use-cases/rel09.0-uc001-canonical-collector.md)

<!-- field: preconditions -->
## Preconditions

- The Jaeger Helm template is deleted and no jaeger values block or schema entry exists.
- The collector agent exposes a query surface on port 18193.
- The UX trace backend is configured for the collector query surface.
- The collector program is consumed from the catalog.

<!-- field: test_cases -->
## Test cases

<!-- case: Default render contains no Jaeger reference | status: implemented -->
### Default render contains no Jaeger reference

A default helm template render produces no string matching jaeger in any template output.

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-canonical-collector](../use-cases/rel09.0-uc001-canonical-collector.md)

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestCollectorDefaultRenderStaysSelfContained$'

<!-- field: expected | form: "tree" -->
#### Expected

- jaeger_in_render: false
- collector_query_port_rendered: true

<!-- field: traces -->
#### Traces

- [rel09.0-uc001-canonical-collector](../use-cases/rel09.0-uc001-canonical-collector.md) [S1](../use-cases/rel09.0-uc001-canonical-collector.md#s1)
- [rel09.0-uc001-canonical-collector](../use-cases/rel09.0-uc001-canonical-collector.md) [S2](../use-cases/rel09.0-uc001-canonical-collector.md#s2)

<!-- case: Spool identity holds without Jaeger cross-check | status: implemented -->
### Spool identity holds without Jaeger cross-check

The smoke test asserts collector spool identity by reading the spool directly via kubectl exec; no Jaeger endpoint is contacted.

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-canonical-collector](../use-cases/rel09.0-uc001-canonical-collector.md)

<!-- field: integration_target -->
#### Integration target

mage integration:helmSmoke

<!-- field: expected | form: "tree" -->
#### Expected

- spool_trace_id_present: true
- service_identity_match: true
- jaeger_endpoint_contacted: false

<!-- field: traces -->
#### Traces

- [rel09.0-uc001-canonical-collector](../use-cases/rel09.0-uc001-canonical-collector.md) [S4](../use-cases/rel09.0-uc001-canonical-collector.md#s4)
- [srd007-collector-agent](../software-requirements/srd007-collector-agent.md) [AC2](../software-requirements/srd007-collector-agent.md#ac2)

<!-- case: Agent-mode metric evidence needs no in-cluster Dolt scrape | status: implemented -->
### Agent-mode metric evidence needs no in-cluster Dolt scrape

In default agent mode the collector agent owns OTLP metric intake; agent metrics are collected and retained as evidence after teardown while the in-cluster Dolt Prometheus scrape is retired (GH-1366). No Jaeger trace evidence is required.

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-canonical-collector](../use-cases/rel09.0-uc001-canonical-collector.md)

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestCollectSharedMetricsEvidenceAgentModeNeedsNoDolt$'

<!-- field: expected | form: "tree" -->
#### Expected

- agent_metrics_present: true
- dolt_metrics_present: false
- trace_evidence_from_jaeger: false

<!-- field: traces -->
#### Traces

- [rel09.0-uc001-canonical-collector](../use-cases/rel09.0-uc001-canonical-collector.md) [S5](../use-cases/rel09.0-uc001-canonical-collector.md#s5)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC5](../software-requirements/srd003-chatbot-deployment.md#ac5)

<!-- case: Contrib collector still retains Dolt metric evidence when opted in | status: implemented -->
### Contrib collector still retains Dolt metric evidence when opted in

When the opt-in contrib collector is selected its Prometheus scrape still surfaces Dolt dss_* metrics, and both agent and Dolt metrics are retained as evidence (GH-1366).

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-canonical-collector](../use-cases/rel09.0-uc001-canonical-collector.md)

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestCollectSharedMetricsEvidenceRetainsAgentAndDolt$'

<!-- field: expected | form: "tree" -->
#### Expected

- agent_metrics_present: true
- dolt_metrics_present: true
- trace_evidence_from_jaeger: false

<!-- field: traces -->
#### Traces

- [rel09.0-uc001-canonical-collector](../use-cases/rel09.0-uc001-canonical-collector.md) [S5](../use-cases/rel09.0-uc001-canonical-collector.md#s5)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC5](../software-requirements/srd003-chatbot-deployment.md#ac5)

<!-- case: Collector program consumed from catalog | status: implemented -->
### Collector program consumed from catalog

The chart profile source redirects the collector program to the catalog root alongside the corpus-ingest program.

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-canonical-collector](../use-cases/rel09.0-uc001-canonical-collector.md)

<!-- field: expected | form: "tree" -->
#### Expected

- local_agents_collector_absent: true
- catalog_source_used: true

<!-- field: traces -->
#### Traces

- [rel09.0-uc001-canonical-collector](../use-cases/rel09.0-uc001-canonical-collector.md) [S6](../use-cases/rel09.0-uc001-canonical-collector.md#s6)

<!-- case: UX waterfall transforms collector query format | status: implemented -->
### UX waterfall transforms collector query format

The ui-kit traceApi module the chatbot UI imports converts collector query response format (ISO timestamps, span_id, parent_span_id) to the internal TraceModel (microsecond timestamps, hierarchical spans).

<!-- field: use_case -->
#### Use case

[rel09.0-uc001-canonical-collector](../use-cases/rel09.0-uc001-canonical-collector.md)

<!-- field: expected | form: "tree" -->
#### Expected

- iso_to_microsecond_conversion: true
- flat_span_structure_parsed: true
- trace_backend_name: collector

<!-- field: traces -->
#### Traces

- [rel09.0-uc001-canonical-collector](../use-cases/rel09.0-uc001-canonical-collector.md) [S3](../use-cases/rel09.0-uc001-canonical-collector.md#s3)
