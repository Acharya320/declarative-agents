---
id: test-rel11.0-otlp-receiver
title: Test Suite for Release 11.0 OTLP Receiver Tools
release: "11.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel11.0-uc001-otlp-receiver-await

<!-- field: preconditions -->
## Preconditions

- srd042-otlp-receiver-tools defines launch, bounded queue, await, spool, relay, stop, self-loop, and metric intake/spool/query behavior.
- srd008-telemetry defines OTLP export, stdouttrace evidence, stable integration identity, and the ingress self-loop rule.
- The tests bind isolated loopback ports and do not require a persistent external collector.
- Focused Go tests provide executable package evidence; deployment-only scope assertions remain identified as non-executable.

<!-- field: test_cases -->
## Test cases

<!-- case: OTLP export enters the bounded queue and reaches await | status: implemented -->
### OTLP export enters the bounded queue and reaches await

A real OTLP TraceService client exports a request containing two resource
spans to a receiver bound on an isolated address. Launch returns before
export. Await consumes exactly one complete FIFO batch and reports its
span count, service names, trace IDs, and resource attributes.

<!-- field: use_case -->
#### Use case

rel11.0-uc001-otlp-receiver-await

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestReceiverExportAwait

<!-- field: inputs | form: "tree" -->
#### Inputs

- queue_capacity: 2
- overflow_policy: reject
- batches:
  1. - services:
       1. chatbot
       1. rag0
     - spans: 3

<!-- field: expected | form: "tree" -->
#### Expected

- launch_signal: ReceiverLaunched
- await_signal: SpansReceived
- batch_count: 1
- span_count: 3
- services:
  1. chatbot
  1. rag0

<!-- field: traces -->
#### Traces

- rel11.0-uc001-otlp-receiver-await S1
- srd042-otlp-receiver-tools AC1

<!-- case: Queue overflow follows each declared policy | status: implemented -->
### Queue overflow follows each declared policy

A capacity-one queue receives two batches before an await. Reject reports
the incoming span count as partial success. Drop-oldest retains the newer
batch, and drop-newest retains the older batch. Concurrent exporters do
not race or exceed capacity.

<!-- field: use_case -->
#### Use case

rel11.0-uc001-otlp-receiver-await

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestReceiverOverflowPolicies

<!-- field: inputs | form: "tree" -->
#### Inputs

- queue_capacity: 1
- policies:
  1. reject
  1. drop_oldest
  1. drop_newest
- concurrent_exporters: 8

<!-- field: expected | form: "tree" -->
#### Expected

- queue_depth_max: 1
- reject: incoming batch reported rejected
- drop_oldest: newest batch awaits
- drop_newest: oldest batch awaits

<!-- field: traces -->
#### Traces

- rel11.0-uc001-otlp-receiver-await S1
- srd042-otlp-receiver-tools AC1

<!-- case: Spool output matches stdouttrace readers and rotates by complete batch | status: implemented -->
### Spool output matches stdouttrace readers and rotates by complete batch

Spool converts a received request into one NDJSON object per span. The
existing trace-reader shape decodes every line. Trace IDs, parent IDs,
timestamps, span attributes, service.name, and integration run identity
equal the source request. Rotation occurs before a batch and leaves no
partial JSON line.

<!-- field: use_case -->
#### Use case

rel11.0-uc001-otlp-receiver-await

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestSpoolStdoutTraceAndRotation

<!-- field: inputs | form: "tree" -->
#### Inputs

- max_files: 3
- max_bytes: one complete batch plus one byte
- batches: 3

<!-- field: expected | form: "tree" -->
#### Expected

- parse_errors: 0
- retained_files_max: 3
- identity: unchanged
- partial_lines: 0

<!-- field: traces -->
#### Traces

- rel11.0-uc001-otlp-receiver-await S2
- srd042-otlp-receiver-tools AC2

<!-- case: Bounded spools reject a missing retained generation | status: implemented -->
### Bounded spools reject a missing retained generation

Both trace and metric factories reject a positive byte cap when
max_files is omitted or one. The runtime guard leaves existing evidence
unchanged even when a builder bypasses factory validation.

<!-- field: use_case -->
#### Use case

rel11.0-uc001-otlp-receiver-await

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run 'Test(SpoolFactoriesRejectByteCap|RotateRejectsSingleGeneration)'

<!-- field: inputs | form: "tree" -->
#### Inputs

- max_bytes: positive
- max_files:
  1. 0
  1. 1

<!-- field: expected | form: "tree" -->
#### Expected

- load_rejected: true
- accepted_evidence_retained: true
- error_names_max_bytes_and_max_files: true

<!-- field: traces -->
#### Traces

- srd042-otlp-receiver-tools AC2

<!-- case: Relay preserves the complete OTLP request and reports rejection | status: implemented -->
### Relay preserves the complete OTLP request and reports rejection

Relay sends the batch selected from command state to an in-process
TraceService server. The upstream request is protobuf-equal to the
received request. A second upstream returns partial success and relay
emits CommandError without modifying or deleting local spool evidence.

<!-- field: use_case -->
#### Use case

rel11.0-uc001-otlp-receiver-await

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run 'TestRelay(PreservesBatch|ReportsPartialRejection)'

<!-- field: inputs | form: "tree" -->
#### Inputs

- upstreams:
  1. accept_all
  1. reject_one
- batch_source: declared command-state selector

<!-- field: expected | form: "tree" -->
#### Expected

- accept_all_signal: SpansRelayed
- request_fields: protobuf-equal
- reject_one_signal: CommandError
- spool_after_rejection: unchanged

<!-- field: traces -->
#### Traces

- rel11.0-uc001-otlp-receiver-await S3
- srd042-otlp-receiver-tools AC3

<!-- case: Stop releases the listener and unblocks every awaiter | status: implemented -->
### Stop releases the listener and unblocks every awaiter

Several awaits block on an empty receiver. Stop unblocks all of them with
ReceiverStopped, rejects later exports, and releases the address so a new
listener can bind it. Repeated stop is safe and the goroutine count
returns to its bounded baseline.

<!-- field: use_case -->
#### Use case

rel11.0-uc001-otlp-receiver-await

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run 'TestReceiverStop(UnblocksAndReleases|AppliesDrainPolicy)'

<!-- field: inputs | form: "tree" -->
#### Inputs

- blocked_awaiters: 4
- drain_policies:
  1. preserve
  1. drop

<!-- field: expected | form: "tree" -->
#### Expected

- await_signals:
  1. ReceiverStopped
  1. ReceiverStopped
  1. ReceiverStopped
  1. ReceiverStopped
- rebind: succeeds
- repeated_stop: ReceiverStopped
- leaked_goroutines: 0

<!-- field: traces -->
#### Traces

- rel11.0-uc001-otlp-receiver-await S4
- srd042-otlp-receiver-tools AC4

<!-- case: Ingress process telemetry cannot loop into its receiver | status: implemented -->
### Ingress process telemetry cannot loop into its receiver

Registration normalizes equivalent wildcard and loopback endpoints and
rejects a receiver whose declared relay points back to itself.

<!-- field: use_case -->
#### Use case

rel11.0-uc001-otlp-receiver-await

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestReceiverRejectsSelfLoopConfiguration

<!-- field: non_executable_scope_assertion -->
#### Non executable scope assertion

An ingress deployment configures file-only process telemetry so client
spans enter the queue while collector lifecycle spans remain in the file.
This deployment constraint is specified by srd042 R7.1 and is not package
test evidence.

<!-- field: inputs | form: "tree" -->
#### Inputs

- equivalent_endpoints:
  1. 1. 0.0.0.0:4317
     1. 127.0.0.1:4317
  1. 1. "[::]:4317"
     1. "[::1]:4317"
- collector_telemetry: file_only

<!-- field: expected | form: "tree" -->
#### Expected

- static_loop: registration failure

<!-- field: traces -->
#### Traces

- rel11.0-uc001-otlp-receiver-await S5
- srd042-otlp-receiver-tools AC5

<!-- case: Receiver accepts trace and metric exports on one listener | status: implemented -->
### Receiver accepts trace and metric exports on one listener

A real OTLP TraceService client and a real OTLP MetricsService client
export to the same bound listener. Both batches are accepted; the metric
export lands on a separate queue and does not disturb the trace queue.

<!-- field: use_case -->
#### Use case

rel11.0-uc001-otlp-receiver-await

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestReceiverExportAwaitMetrics

<!-- field: non_executable_scope_assertion -->
#### Non executable scope assertion

No OTLP logs service is registered, and the deployment needs no separate
metrics backend once the collector agent retains both signals. The
logs-absence and backend-independence are scope assertions specified by
srd042 R8/R9, not package test evidence.

<!-- field: inputs | form: "tree" -->
#### Inputs

- protocols:
  1. trace
  1. metric
  1. log

<!-- field: expected | form: "tree" -->
#### Expected

- trace: accepted
- metric: accepted
- log: not_declared
- metrics_backend: none_required

<!-- field: traces -->
#### Traces

- rel11.0-uc001-otlp-receiver-await S6
- srd042-otlp-receiver-tools AC6

<!-- case: Metric spool records NDJSON and rotates by complete batch | status: implemented -->
### Metric spool records NDJSON and rotates by complete batch

Spool converts a received metric request into one NDJSON record per
metric carrying name, type, data-point count, and resource identity.
Rotation occurs before a complete batch and leaves no partial line. The
metric spool is a separate file from the trace spool.

<!-- field: use_case -->
#### Use case

rel11.0-uc001-otlp-receiver-await

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestSpoolMetricsRoundTripAndRotation

<!-- field: inputs | form: "tree" -->
#### Inputs

- max_files: 3
- max_bytes: one complete batch plus one byte
- batches: 3

<!-- field: expected | form: "tree" -->
#### Expected

- record_name: dispatch_count
- record_type: Sum
- data_point_count: 2
- partial_lines: 0

<!-- field: traces -->
#### Traces

- rel11.0-uc001-otlp-receiver-await S7
- srd042-otlp-receiver-tools AC8

<!-- case: Metric query lists by name and reads back one metric's records | status: implemented -->
### Metric query lists by name and reads back one metric's records

spool_list_metrics aggregates spooled records into per-name summaries with
record and data-point counts and the contributing services, paginated
deterministically. spool_get_metric returns a config-capped page for one
metric name with total, offset, and page size; each record preserves its
decoded resource attributes unchanged, including service.name and
test.run.id, plus the service projection and data points, read-only.

<!-- field: use_case -->
#### Use case

rel11.0-uc001-otlp-receiver-await

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run 'TestSpool(ListMetricsSummaries|GetMetric(ByName|Pagination|SeedCapsAndFinalPage|OffsetPastEndIsEmpty))'

<!-- field: inputs | form: "tree" -->
#### Inputs

- metrics:
  1. - name: dispatch_count
     - services:
       1. chatbot
       1. rag0
     - data_points:
       1. 3
       1. 2
  1. - name: dss_rows
     - services:
       1. dolt
     - data_points:
       1. 1

<!-- field: expected | form: "tree" -->
#### Expected

- list_total: 2
- dispatch_count_records: 2
- dispatch_count_data_points: 5
- detail_page_size_capped: true
- final_and_empty_pages_bounded: true
- dispatch_count_services:
  1. chatbot
  1. rag0
- dispatch_count_resource:
  - service.name: chatbot
  - test.run.id: run-detail-42

<!-- field: traces -->
#### Traces

- rel11.0-uc001-otlp-receiver-await S7
- srd042-otlp-receiver-tools AC8
