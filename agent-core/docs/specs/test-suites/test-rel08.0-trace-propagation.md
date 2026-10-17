---
id: test-rel08.0-trace-propagation
title: Test Suite for Release 08.0 HTTP Trace Propagation
release: "08.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel08.0-uc001-http-trace-propagation

<!-- field: preconditions -->
## Preconditions

- srd016-traceparent (amended) defines outbound injection (R4) and inbound extraction (R5) over REST.
- ParseTraceparent/FormatTraceparent/ParseParentSpan and spawned-child --otel-parent-span propagation are implemented.
- Epic #295 child GH-310 owns the executable inject/extract Go tests these cases name.

<!-- field: test_cases -->
## Test cases

<!-- case: Outbound client request injects traceparent from the active span | status: implemented -->
### Outbound client request injects traceparent from the active span

A REST client word executed inside a recording dispatch span emits an
outbound request whose traceparent header is FormatTraceparent of the active
SpanContext, with tracestate injected when the span carries it.

<!-- field: use_case -->
#### Use case

rel08.0-uc001-http-trace-propagation

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest/client -run TestRESTClient_InjectsTraceparentFromActiveSpan

<!-- field: inputs | form: "tree" -->
#### Inputs

- active_span: recording
- span_tracestate: present

<!-- field: expected | form: "tree" -->
#### Expected

- traceparent_header_present: true
- traceparent_matches_active_span: true
- tracestate_header_present: true

<!-- field: traces -->
#### Traces

- rel08.0-uc001-http-trace-propagation S1
- srd016-traceparent AC4

<!-- case: Outbound injection is uniform and omitted without a span | status: implemented -->
### Outbound injection is uniform and omitted without a span

Injection applies to every outbound client request with no per-operation
configuration, and no traceparent header is emitted when no recording span
is active.

<!-- field: use_case -->
#### Use case

rel08.0-uc001-http-trace-propagation

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest/client -run TestRESTClient_TraceparentUniformAndOmittedWithoutSpan

<!-- field: inputs | form: "tree" -->
#### Inputs

- operations:
  1. get
  1. invoke
  1. send
- active_span: none

<!-- field: expected | form: "tree" -->
#### Expected

- injected_on_all_operations: true
- no_authoring_surface: true
- header_absent_without_span: true

<!-- field: traces -->
#### Traces

- rel08.0-uc001-http-trace-propagation S4
- srd016-traceparent AC4

<!-- case: Inbound extraction parents the request and machine_request span | status: implemented -->
### Inbound extraction parents the request and machine_request span

A REST server endpoint receiving a well-formed traceparent parents its
request span, and the machine_request span, on the extracted remote span, so
the caller's client span and the callee's server span share one trace id.

<!-- field: use_case -->
#### Use case

rel08.0-uc001-http-trace-propagation

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest -run TestRESTServer_ExtractsTraceparentParentsMachineRequestSpan

<!-- field: inputs | form: "tree" -->
#### Inputs

- incoming_traceparent: 00-<trace>-<span>-01

<!-- field: expected | form: "tree" -->
#### Expected

- server_span_parent_is_remote: true
- machine_request_span_parented: true
- same_trace_id_as_client: true

<!-- field: traces -->
#### Traces

- rel08.0-uc001-http-trace-propagation S2
- srd016-traceparent AC5

<!-- case: Absent or malformed header falls back to a new root | status: implemented -->
### Absent or malformed header falls back to a new root

An inbound request with no traceparent, and one with a malformed
traceparent, each start a new root span and complete without error rather
than failing the request.

<!-- field: use_case -->
#### Use case

rel08.0-uc001-http-trace-propagation

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest -run TestRESTServer_TraceparentFallbackToNewRoot

<!-- field: inputs | form: "tree" -->
#### Inputs

- cases:
  1. - traceparent: absent
  1. - traceparent: not-a-traceparent

<!-- field: expected | form: "tree" -->
#### Expected

- new_root_span: true
- request_succeeds: true
- no_error: true

<!-- field: traces -->
#### Traces

- rel08.0-uc001-http-trace-propagation S3
- srd016-traceparent AC5

<!-- case: ParseParentSpan restores a formatted remote parent | status: implemented -->
### ParseParentSpan restores a formatted remote parent

Format one SpanContext, parse it back in-process, and verify trace ID,
span ID, and remote-parent identity. Empty input remains a background
context. This is parser evidence, not subprocess transport evidence.

<!-- field: use_case -->
#### Use case

rel08.0-uc001-http-trace-propagation

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest -run '^TestParseParentSpanRoundTripAndEmptyInput$'

<!-- field: inputs | form: "tree" -->
#### Inputs

- traceparent: formatted remote SpanContext

<!-- field: expected | form: "tree" -->
#### Expected

- parsed_trace_id_preserved: true
- parsed_span_id_preserved: true
- empty_input_has_no_parent: true

<!-- field: traces -->
#### Traces

- rel08.0-uc001-http-trace-propagation S1

<!-- case: Spawned child consumes otel-parent-span | status: implemented -->
### Spawned child consumes otel-parent-span

Start a controlled child through the shared subprocess path and verify its
exported span keeps the caller trace ID and parent span ID.

<!-- field: use_case -->
#### Use case

rel08.0-uc001-http-trace-propagation

<!-- field: go_test -->
#### Go test

go test ./internal/support/execute -run '^TestRunAgent.*OTelParent'

<!-- field: inputs | form: "tree" -->
#### Inputs

- child_launch: subprocess with --otel-parent-span

<!-- field: expected | form: "tree" -->
#### Expected

- child_parented_on_caller: true
- child_span_exported: true
- missing_parent_emits_root: true
- malformed_parent_reports_diagnostic: true
- child_failure_preserves_stderr: true

<!-- field: traces -->
#### Traces

- rel08.0-uc001-http-trace-propagation S1
- srd016-traceparent AC4
