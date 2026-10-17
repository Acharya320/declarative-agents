---
id: srd042-otlp-receiver-tools
title: OTLP Receiver Tools
depends_on:
  - srd_id: srd003-command-dispatch
    symbols_used: [Command, Result, Signal, CommandError]
  - srd_id: srd059-objectstore-tools
    symbols_used: [Opener, ConnectionConfig, bucket URL scheme selection, explicit endpoint]
  - srd_id: srd008-telemetry
    symbols_used: [OTLP exporter, stdouttrace file shape, persistent integration ingress]
  - srd_id: srd022-builtin-tools
    symbols_used: [BuiltinRegistry, ToolDef, DecodeToolConfig]
  - srd_id: srd024-blocking-input-tools
    symbols_used: [blocking input, timeout, cancellation]
  - srd_id: srd038-command-state-store
    symbols_used: [batch_source selector]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Declarative agents export spans through OTLP, but agent-core has no inbound
OTLP vocabulary. Integration deployments therefore require a collector
process whose receive, retain, and relay decisions are outside MachineSpec.
That process cannot participate in the agent monitor, lifecycle, scenario
rig, or declared filtering rules.

OTLP transport belongs in Go because a machine does not parse protobuf or
own sockets. Intake policy belongs in MachineSpec because accepting, testing,
spooling, and relaying a batch are workflow decisions. The runtime needs
small boundary words that expose transport while leaving their order and
branching visible in the grammar.

<!-- field: goals -->
## Goals

- G1: Receive OTLP trace and metric batches through a declared server-binding family.
- G2: Deliver bounded batches to a persistent machine through a blocking await word.
- G3: Write deterministic trace and metric evidence in NDJSON shapes (stdouttrace for spans).
- G4: Relay a selected batch to an upstream OTLP endpoint without changing span identity.
- G5: Prevent the component implementing ingress from exporting telemetry to itself.

<!-- req: R0 -->
## R0 - Tool-language boundary

<!-- req: R0.1 -->
### R0.1

Each receiver tool SHALL follow cfg-tool-authoring-standard and srd022 for category, typed config, signals, side effects, errors, and undo.

<!-- req: R0.2 -->
### R0.2

OTLP protobuf decoding, socket ownership, queue synchronization, NDJSON encoding, and gRPC export SHALL live in the focused `internal/tools/otlp` package.

<!-- req: R0.3 -->
### R0.3

MachineSpec SHALL own launch, await, offline load, filter, spool, relay, retry, and stop order.

<!-- req: R0.4 -->
### R0.4

The package SHALL expose separate `otlp_receiver_launch`, `await_spans`, `load_otlp_batch`, `spool_spans`, `relay_spans`, and `otlp_receiver_stop` words; no word SHALL hide a complete ingress or replay sentence.

<!-- req: R0.5 -->
### R0.5

Internal protobuf validation, stdouttrace-compatible span projection, synchronized complete-batch append, and bounded rotation SHALL remain in process unless a replacement proves field-level compatibility, queue/backpressure preservation, atomic batch writes, rotation semantics, runtime provisioning, and the same declared signals. The existence of otelcol or a stdout exporter alone is not equivalence (GH-1382).

<!-- req: R1 -->
## R1 - Receiver launch

<!-- req: R1.1 -->
### R1.1

`otlp_receiver_launch` SHALL bind OTLP/gRPC on a configured address and register the OTLP TraceService Export operation.

<!-- req: R1.2 -->
### R1.2

The default address SHALL be `0.0.0.0:4317`; automated tests SHALL configure an isolated loopback address.

<!-- req: R1.3 -->
### R1.3

Launch SHALL validate receiver name, address, queue capacity, and overflow policy before listening.

<!-- req: R1.4 -->
### R1.4

The receiver SHALL decode `ExportTraceServiceRequest` batches with `go.opentelemetry.io/proto/otlp` and SHALL NOT translate transport payloads in the gRPC handler.

<!-- req: R1.5 -->
### R1.5

Launch SHALL return without waiting for a batch and SHALL emit `ReceiverLaunched` or `CommandError`.

<!-- req: R1.6 -->
### R1.6

Launch output SHALL include receiver name, bound address, queue capacity, and overflow policy.

<!-- req: R1.7 -->
### R1.7

Launch side effects SHALL be `network_listen` and receiver-state mutation; undo SHALL stop the listener only when this launch owns it.

<!-- req: R2 -->
## R2 - Bounded batch queue

<!-- req: R2.1 -->
### R2.1

Receiver state SHALL hold complete OTLP trace export requests as batches in a bounded FIFO queue.

<!-- req: R2.2 -->
### R2.2

Queue capacity SHALL be at least one and SHALL be fixed at launch.

<!-- req: R2.3 -->
### R2.3

The supported overflow policies SHALL be `reject`, `drop_oldest`, and `drop_newest`; any other value SHALL fail registration.

<!-- req: R2.4 -->
### R2.4

Under `reject`, a full queue SHALL return an OTLP partial-success response whose rejected span count equals the incoming batch span count.

<!-- req: R2.5 -->
### R2.5

Under a drop policy, the receiver SHALL report the dropped batch and span counts in bounded diagnostics and SHALL NOT block the gRPC handler indefinitely.

<!-- req: R2.6 -->
### R2.6

Receiver state SHALL be safe for concurrent exporters, one or more awaiters, and stop.

<!-- req: R3 -->
## R3 - Batch await

<!-- req: R3.1 -->
### R3.1

`await_spans` SHALL block until one queued batch is available, its configured timeout expires, the receiver stops, or context cancellation occurs.

<!-- req: R3.2 -->
### R3.2

Await config SHALL name the receiver and MAY set a timeout; no request body or model output SHALL select the receiver.

<!-- req: R3.3 -->
### R3.3

Await SHALL consume one FIFO batch and emit `SpansReceived`, `AwaitTimedOut`, `ReceiverStopped`, or `CommandError`.

<!-- req: R3.4 -->
### R3.4

Successful output SHALL include the OTLP batch under `batch`, plus `batch_id`, `span_count`, service names, trace IDs, and resource attributes needed by later declared decisions.

<!-- req: R3.5 -->
### R3.5

`batch_id` SHALL be generated by receiver state and SHALL be bounded opaque correlation metadata, not a telemetry resource identity.

<!-- req: R3.6 -->
### R3.6

Await side effects SHALL be limited to queue consumption and command-state output; undo SHALL be noop.

<!-- req: R3.7 -->
### R3.7

`load_otlp_batch` SHALL read one trusted OTLP protobuf-JSON file, decode it as an ExportTraceServiceRequest, and emit `BatchLoaded` or `CommandError` without exporting it.

<!-- req: R3.8 -->
### R3.8

Load output SHALL include path, span count, and the decoded batch so a following declared `relay_spans` word can export it.

<!-- req: R3.9 -->
### R3.9

Read and protobuf-JSON decode failures SHALL remain distinct and SHALL prevent relay dispatch.

<!-- req: R3.10 -->
### R3.10

`load_otlp_batch` SHALL NOT be replaced by an unvalidated byte copy: protobuf-JSON decoding is the trust-boundary check that makes its typed batch safe for relay.

<!-- req: R4 -->
## R4 - Deterministic spool

<!-- req: R4.1 -->
### R4.1

`spool_spans` SHALL resolve an OTLP batch from a declared `batch_source` selector and, with the `filesystem` storage backend, append one NDJSON object per span to a configured file; the `object` backend persists the batch as specified in R10.

<!-- req: R4.2 -->
### R4.2

Each object SHALL use the field names and value shapes emitted by the existing stdouttrace exporter, including `Name`, `SpanContext`, `Parent`, `StartTime`, `EndTime`, `Attributes`, `Events`, `Links`, `Status`, `Resource`, and `InstrumentationScope`.

<!-- req: R4.3 -->
### R4.3

Resource attributes SHALL remain under `Resource` and SHALL preserve `service.name`, integration run identity, and every other received attribute.

<!-- req: R4.4 -->
### R4.4

Appends SHALL serialize concurrent writers, write complete lines, and sync before reporting success.

<!-- req: R4.5 -->
### R4.5

With the `filesystem` backend, config MAY set `max_bytes` and `max_files`; a positive `max_bytes` SHALL require `max_files` of at least two. When the next complete batch would exceed `max_bytes`, rotation SHALL rename older files before writing the new batch and SHALL retain at most `max_files`. The `object` backend has no rotation; its local evidence is bounded by R10.6.

<!-- req: R4.6 -->
### R4.6

A failed append or rotation SHALL emit `CommandError`, SHALL leave the received batch addressable through command state, and SHALL NOT trigger relay implicitly.

<!-- req: R4.7 -->
### R4.7

Successful output SHALL include path, batch id, span count, bytes written, and the selected batch so a following declared word can relay it.

<!-- req: R4.8 -->
### R4.8

Spool words SHALL declare a `filesystem_write` side effect (the NDJSON spool, or the write-ahead log and stage of the `object` backend) and an `external_api` side effect against the configured application object bucket for the `object` backend. Reversibility SHALL be `irreversible` with an `irreversible` undo strategy, because appended evidence and persisted objects are retained audit records.

<!-- req: R5 -->
## R5 - Declared relay

<!-- req: R5.1 -->
### R5.1

`relay_spans` SHALL resolve a batch from a declared `batch_source` selector and export that unchanged request to a configured OTLP/gRPC endpoint.

<!-- req: R5.2 -->
### R5.2

Relay SHALL preserve resource spans, scope spans, trace IDs, span IDs, parent IDs, timestamps, events, links, status, and attributes byte-for-byte at the protobuf field level.

<!-- req: R5.3 -->
### R5.3

Relay config SHALL set an endpoint and MAY set a timeout; TLS, credentials, and endpoint selection from runtime input are outside this release.

<!-- req: R5.4 -->
### R5.4

A partial-success response with rejected spans SHALL emit `CommandError` and report the bounded rejected count and upstream message.

<!-- req: R5.5 -->
### R5.5

Relay SHALL emit `SpansRelayed` or `CommandError`; retry routing SHALL remain a MachineSpec decision.

<!-- req: R5.6 -->
### R5.6

Relay output SHALL include endpoint, batch id, span count, and accepted span count.

<!-- req: R5.7 -->
### R5.7

Relay side effects SHALL be `external_api`; undo SHALL be noop because an OTLP export cannot retract accepted spans.

<!-- req: R6 -->
## R6 - Stop and lifecycle

<!-- req: R6.1 -->
### R6.1

`otlp_receiver_stop` SHALL gracefully stop the named receiver, reject new exports, and unblock all awaiters.

<!-- req: R6.2 -->
### R6.2

Stop config SHALL include a shutdown timeout and a drain policy of `preserve` or `drop`.

<!-- req: R6.3 -->
### R6.3

Stop SHALL emit `ReceiverStopped` or `CommandError` and report queued, drained, and dropped batch counts.

<!-- req: R6.4 -->
### R6.4

Stop SHALL be idempotent for an already stopped receiver owned by the same runtime.

<!-- req: R6.5 -->
### R6.5

No receiver operation SHALL leak a listener or goroutine after stop or failed launch.

<!-- req: R7 -->
## R7 - Self-loop rule

<!-- req: R7.1 -->
### R7.1

A profile that implements shared ingress SHALL configure its own agent telemetry with file export only and SHALL NOT set its OTLP endpoint to the receiver it owns (srd008 R9.6).

<!-- req: R7.2 -->
### R7.2

Registration SHALL reject a statically identical receiver bind address and relay endpoint after normalizing wildcard and loopback aliases.

<!-- req: R7.3 -->
### R7.3

The receiver SHALL preserve a batch whose resource `service.name` equals the collector role when that batch came from another process; service-name filtering alone SHALL NOT implement the self-loop rule.

<!-- req: R7.4 -->
### R7.4

Tests SHALL prove that collector lifecycle and machine spans do not arrive at its receiver while client-agent spans do.

<!-- req: R8 -->
## R8 - Scope and signal coverage

<!-- req: R8.1 -->
### R8.1

The receiver SHALL accept OTLP trace and metric exports on the same gRPC listener.

<!-- req: R8.2 -->
### R8.2

OTLP logs and direct Prometheus scraping SHALL remain outside these words; Dolt metrics reach the receiver as OTLP metric exports forwarded by the in-cluster collector, not by direct scraping.

<!-- req: R8.3 -->
### R8.3

The receiver SHALL satisfy the srd008 R9 shared-ingress contract for both trace and metric evidence, so no separate metrics backend is required.

<!-- req: R9 -->
## R9 - Metric intake, spool, and query

<!-- req: R9.1 -->
### R9.1

`otlp_receiver_launch` SHALL register the OTLP MetricsService Export operation on the same gRPC listener as the TraceService, and SHALL decode `ExportMetricsServiceRequest` batches with `go.opentelemetry.io/proto/otlp` without translating transport payloads in the handler.

<!-- req: R9.2 -->
### R9.2

Receiver state SHALL hold complete OTLP metric export requests as batches under the same bounded FIFO and drop discipline as trace batches, and the drop policy SHALL report dropped batch and data-point counts in bounded diagnostics.

<!-- req: R9.3 -->
### R9.3

With the `filesystem` backend, a declared spool word SHALL append received metric batches as complete NDJSON lines carrying resource attributes, scope, metric name, and data points, under the same lock and bounded rotation as the trace spool, and SHALL NOT leave a partial line; the `object` backend persists metric batches as specified in R10.

<!-- req: R9.4 -->
### R9.4

Declared read-only query words SHALL list and fetch spooled metrics by metric name through config-capped pagination, and each fetched detail record SHALL return its decoded resource attributes unchanged, including `service.name` and `test.run.id`, so callers can identify the originating integration run while preserving the service and metric payload without an unbounded response from a persistent spool.

<!-- req: R9.5 -->
### R9.5

Trace intake, spool, relay, and query behavior SHALL remain unchanged; metric words SHALL NOT alter span records or the stdouttrace shape.

<!-- req: R10 -->
## R10 - Durable object-bucket spool storage

<!-- req: R10.1 -->
### R10.1

The trace and metric spool words SHALL select a storage backend through declared configuration: `filesystem` (the default, the existing bounded NDJSON append, unchanged) or `object`. The backend, bucket URL, endpoint, and in-bucket prefix SHALL be declared configuration only, and no code path SHALL read an environment variable for backend selection (srd008 R8, R9).

<!-- req: R10.2 -->
### R10.2

The object backend SHALL persist one versioned immutable envelope per batch containing at least schema version, signal, application, namespace, run, collector instance, stable batch identity, receive time, payload format, payload checksum, and payload; traces and metrics SHALL persist as separate object families (srd008 R3, R5, R6.3).

<!-- req: R10.3 -->
### R10.3

Stable batch identity SHALL be derived from batch content, independent of the receiver's process-local sequence, so a retry resolves to one object and two process lifetimes whose sequence both restart at one still yield distinct identities for distinct batches (srd008 R4).

<!-- req: R10.4 -->
### R10.4

The object backend SHALL record the final object key, envelope metadata, checksum, and pending/committed upload state in a local write-ahead log, and SHALL stage the envelope bytes and sync that record before the first remote attempt, so bucket latency or outage leaves replayable evidence and a restart resumes the exact pre-recorded keys (srd008 R2, R4, R6).

<!-- req: R10.5 -->
### R10.5

Object creation SHALL be idempotent: the same key with the same checksum is success, and the same key with different bytes is a named integrity error and never an overwrite. A committed WAL record SHALL be compacted only after the remote object is verified (srd008 R7, AC5).

<!-- req: R10.6 -->
### R10.6

The write-ahead log SHALL bound its pending bytes and pending record count; overflow SHALL be a named, observable error rather than a silent drop of uncommitted evidence (srd008 R10).

<!-- req: R10.7 -->
### R10.7

Every object envelope SHALL stamp the declared retention class, optional retain-until instant derived from a non-negative retention-days policy, and the fixed destructive authority `app:purge`. Retention metadata SHALL NOT cause automatic deletion during collector, application, or platform teardown.

<!-- req: R11 -->
## R11 - Object-backed query read

<!-- req: R11.1 -->
### R11.1

The trace and metric query words SHALL select a read backend through declared configuration: the filesystem NDJSON spool (default, unchanged) or the object backend, which reads the configured application bucket and in-bucket prefix. The bucket, endpoint, and prefix SHALL be declared configuration only; no request field SHALL redirect a query to another bucket (srd008 R2.2).

<!-- req: R11.2 -->
### R11.2

The object read SHALL project each immutable envelope's payload back into the existing spooled span and metric record shapes, so the trace, metric, span-stats, and breakdown response shapes are unchanged (R4.2).

<!-- req: R11.3 -->
### R11.3

The object read SHALL merge committed objects with the pending local WAL by the R10 stable batch identity, producing no duplicate logical batch.

<!-- req: R11.4 -->
### R11.4

The object read SHALL be bounded per request by object count, byte budget, and deadline, and SHALL never scan the whole bucket; a listing truncated by budget SHALL report partial.

<!-- req: R11.5 -->
### R11.5

The query SHALL report storage status as complete, partial, or unavailable. A bucket outage SHALL report unavailable and SHALL NOT present WAL-only data as complete history.

<!-- field: non_goals -->
## Non-goals

- Receiving OTLP logs.
- Acting as a general-purpose, firehose-rate telemetry collector.
- Sampling languages beyond machine composition from existing predicate words.
- Querying Jaeger or serving a trace search API.
- TLS, authentication, compression negotiation beyond gRPC defaults, or multi-tenant isolation.
- Replacing deterministic per-run file export from producing agents.
- Requiring otelcol as a runtime dependency for the receiver's declared spool contract without a compatibility and provisioning decision.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

A declared launch word accepts an OTLP/gRPC trace export into a bounded queue and a declared await word returns that complete batch to a machine, while transport/spool internals retain their compatibility and provisioning gate.

<!-- field: traces -->
#### Traces

- [R0.1](#r01)
- [R0.2](#r02)
- [R0.3](#r03)
- [R0.4](#r04)
- [R0.5](#r05)
- [R1.1](#r11)
- [R1.2](#r12)
- [R1.3](#r13)
- [R1.4](#r14)
- [R1.5](#r15)
- [R1.6](#r16)
- [R2.1](#r21)
- [R2.2](#r22)
- [R2.3](#r23)
- [R2.4](#r24)
- [R2.5](#r25)
- [R2.6](#r26)
- [R3.1](#r31)
- [R3.2](#r32)
- [R3.3](#r33)
- [R3.4](#r34)
- [R3.5](#r35)
- [R3.6](#r36)

<!-- ac: AC2 -->
### AC2

A received batch is appended as complete stdouttrace-compatible NDJSON lines with span and resource identity intact, and bounded rotation never leaves a partial line.

<!-- field: traces -->
#### Traces

- [R4.1](#r41)
- [R4.2](#r42)
- [R4.3](#r43)
- [R4.4](#r44)
- [R4.5](#r45)
- [R4.6](#r46)
- [R4.7](#r47)
- [R4.8](#r48)

<!-- ac: AC3 -->
### AC3

A declared relay exports the selected batch to an upstream OTLP endpoint without changing protobuf span identity, and upstream rejection is a machine-visible fault.

<!-- field: traces -->
#### Traces

- [R5.1](#r51)
- [R5.2](#r52)
- [R5.3](#r53)
- [R5.4](#r54)
- [R5.5](#r55)
- [R5.6](#r56)
- [R5.7](#r57)

<!-- ac: AC4 -->
### AC4

Stop rejects new exports, unblocks awaiters, obeys the drain policy, and leaves no listener or goroutine behind.

<!-- field: traces -->
#### Traces

- [R1.7](#r17)
- [R6.1](#r61)
- [R6.2](#r62)
- [R6.3](#r63)
- [R6.4](#r64)
- [R6.5](#r65)

<!-- ac: AC5 -->
### AC5

The ingress agent exports its own telemetry to file only, rejects a static relay self-loop, and receives client spans without receiving its own lifecycle spans.

<!-- field: traces -->
#### Traces

- [R7.1](#r71)
- [R7.2](#r72)
- [R7.3](#r73)
- [R7.4](#r74)

<!-- ac: AC6 -->
### AC6

The receiver accepts OTLP trace and metric exports on one gRPC listener, makes no OTLP logs claim, and satisfies the persistent integration metrics contract without a separate metrics backend.

<!-- field: traces -->
#### Traces

- [R8.1](#r81)
- [R8.2](#r82)
- [R8.3](#r83)

<!-- ac: AC7 -->
### AC7

An offline replay machine validates and loads one trusted OTLP protobuf-JSON batch and relays it unchanged through separate visible words, with distinct read and decode failures and no unvalidated byte-copy substitute.

<!-- field: traces -->
#### Traces

- [R0.2](#r02)
- [R0.3](#r03)
- [R0.4](#r04)
- [R3.7](#r37)
- [R3.8](#r38)
- [R3.9](#r39)
- [R3.10](#r310)
- [R5.1](#r51)
- [R5.2](#r52)
- [R5.4](#r54)
- [R5.5](#r55)

<!-- ac: AC8 -->
### AC8

The receiver accepts an OTLP/gRPC metric export into the same bounded queue, spools it as NDJSON with resource and metric identity intact under bounded rotation, and serves it back through declared, config-capped paginated query words whose detail records preserve resource attributes such as `service.name` and `test.run.id`, leaving trace intake, spool, relay, and query behavior unchanged.

<!-- field: traces -->
#### Traces

- [R9.1](#r91)
- [R9.2](#r92)
- [R9.3](#r93)
- [R9.4](#r94)
- [R9.5](#r95)

<!-- ac: AC9 -->
### AC9

A spool word with the default filesystem backend behaves exactly as before, while the object backend persists one versioned immutable envelope per batch under a content-derived stable identity, records the object in a bounded write-ahead log before the first upload, creates the object idempotently, refuses a different-content collision, resumes pending keys after a restart, keeps traces and metrics in separate object families, and stamps retention class, optional retain-until, and fixed app:purge authority without automatic teardown deletion. Unit tests prove this on mem:// and file:// with no container; a fake-GCS integration proves it end to end over gs:// and skips naming docker where docker is absent.

<!-- field: traces -->
#### Traces

- [R10.1](#r101)
- [R10.2](#r102)
- [R10.3](#r103)
- [R10.4](#r104)
- [R10.5](#r105)
- [R10.6](#r106)
- [R10.7](#r107)

<!-- ac: AC10 -->
### AC10

A query word with the object backend reads the configured application bucket (never a request-supplied bucket), projects immutable envelopes into the existing response shapes, merges committed objects with the pending WAL without duplicates, bounds each read by object count and byte budget, and reports complete, partial, or unavailable storage status; a bucket outage never labels WAL-only data complete. The filesystem backend is unchanged. Unit tests prove this on file:// with no container.

<!-- field: traces -->
#### Traces

- [R11.1](#r111)
- [R11.2](#r112)
- [R11.3](#r113)
- [R11.4](#r114)
- [R11.5](#r115)

<!-- field: related_documents -->
## Related documents

- cfg-tool-authoring-standard
- sm-tool-language
- srd029-rest-server-tools

<!-- field: implementation_status -->
## Implementation status

<!-- field: implemented -->
### Implemented

- `internal/tools/otlp` owns receiver state, launch, await, offline batch load, spool, relay, stop, and focused tests for the trace signal.
- `tools/builtin/otlp/all.yaml` declares the six OTLP trace words and registry composition.
- `go test ./internal/tools/otlp` proves bounded receive, await, offline load-to-relay, spool compatibility, relay identity, lifecycle cleanup, and static self-loop rejection.
- R10 durable object-bucket spool storage: `spoolstore.go`, `envelope.go`, and `wal.go` add the storage port, the versioned immutable envelope, content-derived stable batch identity, and the bounded write-ahead log; the trace and metric spool words select the backend through declared `storage` configuration and reuse the srd059 objectstore opener. `go test ./internal/tools/otlp` proves object persistence, WAL-first replay, idempotence, integrity refusal, and WAL bounds on mem:// and file://, and the fake-GCS integration proves the gs:// path (GH-2475).

<!-- field: planned -->
### Planned

- R8/R9 metric intake, spool, and query (the MetricsService registration, `spool_metrics`, and metric query words) are specified here and implemented under GH-1204; the receiver is trace-only until that lands.
