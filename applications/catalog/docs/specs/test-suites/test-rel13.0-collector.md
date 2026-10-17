---
id: test-rel13.0-collector
title: Test Suite for Release 13.0 Collector Family
release: "13.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- [rel13.0-uc001-collector-family](../use-cases/rel13.0-uc001-collector-family.md)

<!-- field: preconditions -->
## Preconditions

- srd020-collector defines the canonical collector family with default spool mode.
- srd042-otlp-receiver-tools provides the OTLP builtins the collector exercises.
- The agent-core binary is built and available via RequireCoreRoot.

<!-- field: test_cases -->
## Test cases

<!-- case: Collector launches in spool mode and shuts down cleanly | status: implemented -->
### Collector launches in spool mode and shuts down cleanly

Launch the collector profile with default environment (spool mode, no
relay endpoint). Send an exit request via the control endpoint. Assert
the machine reaches Done with no error spans and the expected tool
vocabulary.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-collector-family](../use-cases/rel13.0-uc001-collector-family.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCollectorSpoolModeConformance

<!-- field: inputs | form: "tree" -->
#### Inputs

- collector_mode: spool (default)
- relay_endpoint: unset

<!-- field: expected | form: "tree" -->
#### Expected

- exit_code: 0
- terminal_state: Done
- tool_spans_present:
  1. otlp_receiver_launch
  1. launch_collector_control
  1. launch_collector_monitor
  1. launch_collector_query
  1. await_spans
  1. await_collector_control
  1. exit_agent
  1. otlp_receiver_stop
  1. stop_collector_query
  1. stop_collector_monitor
  1. stop_collector_control

<!-- field: traces -->
#### Traces

- [rel13.0-uc001-collector-family](../use-cases/rel13.0-uc001-collector-family.md) [S1](../use-cases/rel13.0-uc001-collector-family.md#s1)
- [rel13.0-uc001-collector-family](../use-cases/rel13.0-uc001-collector-family.md) [S3](../use-cases/rel13.0-uc001-collector-family.md#s3)
- [srd020-collector](../software-requirements/srd020-collector.md) [AC1](../software-requirements/srd020-collector.md#ac1)

<!-- case: Empty spool answers the query surface with zero traces | status: implemented -->
### Empty spool answers the query surface with zero traces

A collector with no traffic yet has no spool file; GET /query/traces
answers 200 with zero traces instead of an error, and a directory
misconfigured as the spool path is reported explicitly (GH-1168).

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-collector-family](../use-cases/rel13.0-uc001-collector-family.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCollectorQueryEmptySpool

<!-- field: inputs | form: "tree" -->
#### Inputs

- spool: absent (no batch received)

<!-- field: expected | form: "tree" -->
#### Expected

- status: 200
- total: 0

<!-- field: traces -->
#### Traces

- [rel13.0-uc001-collector-family](../use-cases/rel13.0-uc001-collector-family.md) [S2](../use-cases/rel13.0-uc001-collector-family.md#s2)
- [srd020-collector](../software-requirements/srd020-collector.md) [AC2](../software-requirements/srd020-collector.md#ac2)

<!-- case: Collector boots from the shipped profile with default environment | status: implemented -->
### Collector boots from the shipped profile with default environment

Launch the shipped profile without rewriting any file: only ports and
the spool path come from the environment, the relay endpoint stays
unset, and the machine reaches health then a clean Done (GH-1163). This
is the case the port-randomizing suites cannot cover, because they
patch the profile before launch.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-collector-family](../use-cases/rel13.0-uc001-collector-family.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCollectorDefaultEnvironmentBoot

<!-- field: inputs | form: "tree" -->
#### Inputs

- collector_mode: spool (default)
- relay_endpoint: unset
- profile: shipped, unpatched

<!-- field: expected | form: "tree" -->
#### Expected

- exit_code: 0
- terminal_state: Done

<!-- field: traces -->
#### Traces

- [rel13.0-uc001-collector-family](../use-cases/rel13.0-uc001-collector-family.md) [S1](../use-cases/rel13.0-uc001-collector-family.md#s1)
- [srd020-collector](../software-requirements/srd020-collector.md) [AC1](../software-requirements/srd020-collector.md#ac1)

<!-- case: Query response keys match the UI and smoke consumers | status: implemented -->
### Query response keys match the UI and smoke consumers

Pin the exact JSON keys of a trace summary and a span as emitted by the
query surface (GH-1164). The trace UI client and the coding-agent smoke
verdict decode these keys; drift on either side fails here first.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-collector-family](../use-cases/rel13.0-uc001-collector-family.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCollectorQueryResponseContract

<!-- field: inputs | form: "tree" -->
#### Inputs

- spool: seeded two-trace fixture

<!-- field: expected | form: "tree" -->
#### Expected

- summary_keys:
  1. duration_ms
  1. root_service
  1. root_span_name
  1. span_count
  1. start_time
  1. trace_id
- span_keys:
  1. attributes
  1. end_time
  1. name
  1. parent_span_id
  1. service
  1. span_id
  1. start_time
  1. status

<!-- field: traces -->
#### Traces

- [rel13.0-uc001-collector-family](../use-cases/rel13.0-uc001-collector-family.md) [S2](../use-cases/rel13.0-uc001-collector-family.md#s2)
- [srd020-collector](../software-requirements/srd020-collector.md) [AC2](../software-requirements/srd020-collector.md#ac2)

<!-- case: Monitor views carry no span payloads while the query surface serves them | status: implemented -->
### Monitor views carry no span payloads while the query surface serves them

Seed a spool, launch the collector, and read every monitor view. Each
answers 200 and none carries fixture span content, while the query
surface returns the same spool's traces. The two surfaces are scoped
apart rather than one merely being quiet.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-collector-family](../use-cases/rel13.0-uc001-collector-family.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCollectorMonitorServesNoSpanPayloads

<!-- field: inputs | form: "tree" -->
#### Inputs

- spool: seeded with 2 traces (3 spans)
- monitor_views: machine, machines, tools/declared, state, tools, metrics, events

<!-- field: expected | form: "tree" -->
#### Expected

- monitor_status: 200
- span_content_in_monitor: none
- query_returns_traces: true

<!-- field: traces -->
#### Traces

- [srd020-collector](../software-requirements/srd020-collector.md) [R2.1](../software-requirements/srd020-collector.md#r21)
- [srd020-collector](../software-requirements/srd020-collector.md) [R2.2](../software-requirements/srd020-collector.md#r22)
- [srd020-collector](../software-requirements/srd020-collector.md) [AC9](../software-requirements/srd020-collector.md#ac9)

<!-- case: Query list traces returns paginated summaries from seeded spool | status: implemented -->
### Query list traces returns paginated summaries from seeded spool

Seed a spool with two traces, launch the collector, and GET /query/traces
with page_size=1. Assert response contains one trace, total=2, and
page_size=1.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-collector-family](../use-cases/rel13.0-uc001-collector-family.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCollectorQueryListTraces

<!-- field: inputs | form: "tree" -->
#### Inputs

- spool: seeded with 2 traces (3 spans)
- page_size: 1

<!-- field: expected | form: "tree" -->
#### Expected

- status: 200
- total: 2
- page_size: 1
- traces_returned: 1

<!-- field: traces -->
#### Traces

- [srd020-collector](../software-requirements/srd020-collector.md) [R5.1](../software-requirements/srd020-collector.md#r51)
- [srd020-collector](../software-requirements/srd020-collector.md) [R5.3](../software-requirements/srd020-collector.md#r53)
- [srd020-collector](../software-requirements/srd020-collector.md) [AC9](../software-requirements/srd020-collector.md#ac9)

<!-- case: Query get trace returns spans for one trace from seeded spool | status: implemented -->
### Query get trace returns spans for one trace from seeded spool

Seed a spool with two traces, launch the collector, and GET
/query/traces/trace-aaa. Assert response contains trace_id=trace-aaa
and span_count=2.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-collector-family](../use-cases/rel13.0-uc001-collector-family.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCollectorQueryGetTrace

<!-- field: inputs | form: "tree" -->
#### Inputs

- spool: seeded with 2 traces (3 spans)
- trace_id: trace-aaa

<!-- field: expected | form: "tree" -->
#### Expected

- status: 200
- trace_id: trace-aaa
- span_count: 2

<!-- field: traces -->
#### Traces

- [srd020-collector](../software-requirements/srd020-collector.md) [R5.2](../software-requirements/srd020-collector.md#r52)
- [srd020-collector](../software-requirements/srd020-collector.md) [AC9](../software-requirements/srd020-collector.md#ac9)

<!-- case: Canonical design-token import covers collector UI App.css | status: implemented -->
### Canonical design-token import covers collector UI App.css

Collector App.css must import @declarative-agents/ui-kit/tokens.css,
which resolves through its file: dependency to
applications/ui-kit/src/tokens.css (srd004 R8.2), and contain no copied
canonical token declarations.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-collector-family](../use-cases/rel13.0-uc001-collector-family.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestDesignTokensImportsResolveCanonical/collector

<!-- field: inputs | form: "tree" -->
#### Inputs

- canonical: applications/ui-kit/src/tokens.css
- consumer: agents/collector/ui/src/App.css

<!-- field: expected | form: "tree" -->
#### Expected

- drift: none

<!-- field: traces -->
#### Traces

- [srd020-collector](../software-requirements/srd020-collector.md) [R6.1](../software-requirements/srd020-collector.md#r61)
- [srd020-collector](../software-requirements/srd020-collector.md) [AC10](../software-requirements/srd020-collector.md#ac10)

<!-- case: Collector UI dist reproduces from a clean source build | status: implemented -->
### Collector UI dist reproduces from a clean source build

The tracked dist directory matches the output of a clean npm ci + build
from the lockfile-pinned source. Verified by the root uiDist mage gate.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-collector-family](../use-cases/rel13.0-uc001-collector-family.md)

<!-- field: go_test -->
#### Go test

cd ../.. && mage uiDist

<!-- field: inputs | form: "tree" -->
#### Inputs

- ui_dir: agents/collector/ui

<!-- field: expected | form: "tree" -->
#### Expected

- dist_diff: none

<!-- field: traces -->
#### Traces

- [srd020-collector](../software-requirements/srd020-collector.md) [R6.2](../software-requirements/srd020-collector.md#r62)
- [srd020-collector](../software-requirements/srd020-collector.md) [AC10](../software-requirements/srd020-collector.md#ac10)

<!-- case: Durable object storage, retention metadata, and query completeness remain explicit | status: implemented -->
### Durable object storage, retention metadata, and query completeness remain explicit

Object-backed spool tests prove immutable envelopes, WAL replay, and
retention_class/retain_until/purge_authority metadata. Object query tests
prove bounded bucket/WAL merge and complete, partial, or unavailable
storage status without request-selected bucket authority.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-collector-family](../use-cases/rel13.0-uc001-collector-family.md)

<!-- field: go_test -->
#### Go test

cd ../../../../agent-core && go test ./internal/tools/otlp

<!-- field: expected | form: "tree" -->
#### Expected

- filesystem_default_unchanged: true
- purge_authority: app:purge
- automatic_teardown_deletion: false

<!-- field: traces -->
#### Traces

- [srd020-collector](../software-requirements/srd020-collector.md) [AC11](../software-requirements/srd020-collector.md#ac11)

<!-- case: Trace UI names incomplete retained history | status: implemented -->
### Trace UI names incomplete retained history

Pin the bounded status read and the complete, partial, unavailable, and
unknown trace-view banner, including explicit incomplete-history copy.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-collector-family](../use-cases/rel13.0-uc001-collector-family.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCollectorUIStorageStatusContract

<!-- field: expected | form: "tree" -->
#### Expected

- partial_is_not_empty_healthy: true
- unavailable_is_not_complete: true

<!-- field: traces -->
#### Traces

- [srd020-collector](../software-requirements/srd020-collector.md) [AC12](../software-requirements/srd020-collector.md#ac12)

<!-- case: Explicit purge refuses mismatched identity and isolates peer data | status: implemented -->
### Explicit purge refuses mismatched identity and isolates peer data

Run the approved model-free profile first with a mismatched token and
then with the exact application token. The former reaches Refused without
mutation; the latter reaches Purged, deletes only the rendered prefix,
verifies zero remaining, preserves a peer prefix, and writes the external
audit.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-collector-family](../use-cases/rel13.0-uc001-collector-family.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestApplicationPurgeRefusesThenDeletesOnlyRenderedPrefix

<!-- field: expected | form: "tree" -->
#### Expected

- refused_before_storage: true
- peer_prefix_unchanged: true
- external_audit_remaining: 0

<!-- field: traces -->
#### Traces

- [srd020-collector](../software-requirements/srd020-collector.md) [AC13](../software-requirements/srd020-collector.md#ac13)
