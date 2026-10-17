---
id: test-rel07.0-collector
title: Test Suite for Release 07.0 Collector Agent
release: "07.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel07.0-uc001-collector-ingress

<!-- field: preconditions -->
## Preconditions

- agent-core srd042 receiver, await, spool, relay, and stop words are available.
- The collector profile exports its process telemetry to file only.
- The default production chart routes traces and OTLP metrics to the local collector agent. Integration overlays keep traces on that intake and its declared relay, but export agent metrics directly to the persistent external collector because the collector agent does not relay its metric spool. No contrib metrics gateway or in-cluster Dolt scrape is deployed (GH-1207, GH-1366).

<!-- field: test_cases -->
## Test cases

<!-- case: Collector profile validates in test and relay modes | status: implemented -->
### Collector profile validates in test and relay modes

Both mode overlays load the same profile and every action, signal, declaration, control route, and monitor route validates.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-collector-ingress

<!-- field: expected | form: "tree" -->
#### Expected

- test_mode: spool without relay dependency
- relay_mode: spool precedes relay
- hidden_transport_policy: false

<!-- field: traces -->
#### Traces

- rel07.0-uc001-collector-ingress S1
- rel07.0-uc001-collector-ingress S3
- srd007-collector-agent AC1
- srd007-collector-agent AC3

<!-- case: A real client batch reaches stdouttrace spool unchanged | status: implemented -->
### A real client batch reaches stdouttrace spool unchanged

A real OTLP client exports parent and child spans with full run identity; the spool preserves every identity field.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-collector-ingress

<!-- field: integration_target -->
#### Integration target

mage integration:helmSmoke

<!-- field: expected | form: "tree" -->
#### Expected

- trace_id_match: true
- span_id_match: true
- parent_id_match: true
- service_and_run_identity_match: true
- complete_ndjson_lines: true

<!-- field: traces -->
#### Traces

- rel07.0-uc001-collector-ingress S2
- srd007-collector-agent AC2

<!-- case: Canned positive-span batch takes the collector filter-to-spool leg | status: implemented -->
### Canned positive-span batch takes the collector filter-to-spool leg

The hermetic rig launches the shipped collector profile, exports a canned parent/child OTLP batch over gRPC, and requires two attributable stdouttrace spool lines.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-collector-ingress

<!-- field: integration_target -->
#### Integration target

mage integration:rig

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestCollectorIntakeFilterScenario$'

<!-- field: expected | form: "tree" -->
#### Expected

- positive_filter_decision: true
- spooled_span_count: 2
- service_identity_preserved: collector-rig-fixture

<!-- field: traces -->
#### Traces

- rel07.0-uc001-collector-ingress S2
- srd007-collector-agent AC2

<!-- case: Collector lifecycle releases receiver control and monitor listeners | status: implemented -->
### Collector lifecycle releases receiver control and monitor listeners

The hermetic rig launches the shipped collector profile, verifies monitor reachability, sends lifecycle exit, asserts terminal state succeeded, and rebinds all listener addresses.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-collector-ingress

<!-- field: integration_target -->
#### Integration target

mage integration:rig

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestCollectorLifecycleRebindAndTerminalState$'

<!-- field: expected | form: "tree" -->
#### Expected

- blocked_await_released: true
- all_addresses_rebind: true
- final_monitor_state_bounded: true

<!-- field: traces -->
#### Traces

- rel07.0-uc001-collector-ingress S4
- srd007-collector-agent AC4

<!-- case: Collector process telemetry cannot enter its receiver | status: implemented -->
### Collector process telemetry cannot enter its receiver

Client spans enter while collector lifecycle and tool spans appear only in the collector file; equivalent relay self-endpoints fail validation.

<!-- field: use_case -->
#### Use case

rel07.0-uc001-collector-ingress

<!-- field: integration_target -->
#### Integration target

mage integration:helmSmoke

<!-- field: expected | form: "tree" -->
#### Expected

- client_batch_received: true
- collector_batch_received: false
- equivalent_self_endpoint_rejected: true

<!-- field: traces -->
#### Traces

- rel07.0-uc001-collector-ingress S5
- srd007-collector-agent AC5

<!-- case: Agent mode separates integration metric export from trace intake | status: implemented -->
### Agent mode separates integration metric export from trace intake

Production agents send traces and OTLP metrics to the local collector. The integration overlay keeps traces on that collector and its declared relay while agents export metrics directly to the persistent host collector, preserving evidence without a contrib gateway or in-cluster Dolt Prometheus scrape (GH-1366).

<!-- field: use_case -->
#### Use case

rel07.0-uc001-collector-ingress

<!-- field: integration_target -->
#### Integration target

mage integration:helmSmoke

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestCollectorDefaultRenderStaysSelfContained$|^TestCollectorKindOverlayExportsMetricsDirectlyAndTracesLocally$'

<!-- field: expected | form: "tree" -->
#### Expected

- trace_ingress: collector_agent
- production_metric_ingress: collector_agent
- integration_metric_ingress: persistent_host_collector
- trace_relay: collector_agent_declared_relay
- agent_metrics_available: true
- dolt_metrics_available: false
- combined_health_required: true

<!-- field: traces -->
#### Traces

- rel07.0-uc001-collector-ingress S6
- srd007-collector-agent AC6
- srd003-chatbot-deployment AC5

<!-- case: Persistent evidence survives mesh teardown | status: implemented -->
### Persistent evidence survives mesh teardown

After chart and cluster teardown, representative trace ids relayed by the local collector and run-correlated agent metrics exported directly to the persistent host collector remain queryable; Dolt metrics apply only under the opt-in contrib collector (GH-1366).

<!-- field: use_case -->
#### Use case

rel07.0-uc001-collector-ingress

<!-- field: integration_target -->
#### Integration target

mage integration:observabilityEvidence

<!-- field: expected | form: "tree" -->
#### Expected

- trace_evidence_retained: true
- metric_evidence_retained: true
- run_identity_filterable: true

<!-- field: traces -->
#### Traces

- rel07.0-uc001-collector-ingress S6
- srd007-collector-agent AC6
- srd003-chatbot-deployment AC5
