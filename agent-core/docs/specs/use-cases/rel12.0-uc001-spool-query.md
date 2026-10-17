---
id: rel12.0-uc001-spool-query
title: A collector reads trace summaries and span detail from its NDJSON spool
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An operator or UI queries the collector's spooled traces. The list view
returns paginated trace summaries aggregated from the active and rotated
spool files, newest first. Selecting a trace returns all its spans with the
fields a waterfall view needs. Malformed or truncated lines are skipped and
counted, and an empty spool returns an empty result.

<!-- field: actor -->
## Actor

Operator or downstream query surface reading the collector spool

<!-- field: trigger -->
## Trigger

The operator or a REST query surface dispatches spool_list_traces to browse
ingested traces, or spool_get_trace to inspect one trace's spans.

<!-- field: flow -->
## Flow

- F1: The profile machine dispatches `spool_list_traces` with the spool path, page size, max page size, and offset from trusted config.
- F2: The tool reads the active spool file and every rotated file discovered by the srd042 naming convention, skipping malformed lines.
- F3: Spans are aggregated by trace ID into summaries with root service, root span name, span count, start time, and duration.
- F4: Summaries are sorted newest first, paginated by offset and capped page size, and returned as `TracesListed`.
- F5: The operator selects a trace. The machine dispatches `spool_get_trace` with the trace ID and spool path.
- F6: The tool reads the spool files, collects all spans matching the trace ID, sorts them by start time, and returns them as `TraceRetrieved` with waterfall fields.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd044-spool-query-tools](../software-requirements/srd044-spool-query-tools.md) [R1.1](../software-requirements/srd044-spool-query-tools.md#r11), [R1.2](../software-requirements/srd044-spool-query-tools.md#r12), [R1.3](../software-requirements/srd044-spool-query-tools.md#r13), [R1.4](../software-requirements/srd044-spool-query-tools.md#r14), [R1.5](../software-requirements/srd044-spool-query-tools.md#r15) | paginated trace listing reads active and rotated files |
| T2 | [srd044-spool-query-tools](../software-requirements/srd044-spool-query-tools.md) [R2.1](../software-requirements/srd044-spool-query-tools.md#r21), [R2.2](../software-requirements/srd044-spool-query-tools.md#r22), [R2.3](../software-requirements/srd044-spool-query-tools.md#r23), [R2.4](../software-requirements/srd044-spool-query-tools.md#r24), [R2.5](../software-requirements/srd044-spool-query-tools.md#r25) | per-trace span detail with waterfall fields |
| T3 | [srd044-spool-query-tools](../software-requirements/srd044-spool-query-tools.md) [R3.1](../software-requirements/srd044-spool-query-tools.md#r31), [R3.2](../software-requirements/srd044-spool-query-tools.md#r32), [R3.3](../software-requirements/srd044-spool-query-tools.md#r33) | read-only with malformed-line tolerance |
| T4 | [srd042-otlp-receiver-tools](../software-requirements/srd042-otlp-receiver-tools.md) [R4.1](../software-requirements/srd042-otlp-receiver-tools.md#r41), [R4.2](../software-requirements/srd042-otlp-receiver-tools.md#r42), [R4.3](../software-requirements/srd042-otlp-receiver-tools.md#r43), [R4.4](../software-requirements/srd042-otlp-receiver-tools.md#r44), [R4.5](../software-requirements/srd042-otlp-receiver-tools.md#r45), [R4.6](../software-requirements/srd042-otlp-receiver-tools.md#r46), [R4.7](../software-requirements/srd042-otlp-receiver-tools.md#r47), [R4.8](../software-requirements/srd042-otlp-receiver-tools.md#r48) | the spool NDJSON shape these tools read is written by spool_spans |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Trace summaries aggregate spans from active and rotated spool files with correct root service, span count, timing, and newest-first order.

<!-- field: traces -->
#### Traces

- [srd044-spool-query-tools](../software-requirements/srd044-spool-query-tools.md) [AC1](../software-requirements/srd044-spool-query-tools.md#ac1)

<!-- sc: S2 -->
### S2

Pagination respects offset and page size cap.

<!-- field: traces -->
#### Traces

- [srd044-spool-query-tools](../software-requirements/srd044-spool-query-tools.md) [AC2](../software-requirements/srd044-spool-query-tools.md#ac2)

<!-- sc: S3 -->
### S3

Per-trace detail returns all matching spans with waterfall fields, and an unknown trace returns zero spans.

<!-- field: traces -->
#### Traces

- [srd044-spool-query-tools](../software-requirements/srd044-spool-query-tools.md) [AC3](../software-requirements/srd044-spool-query-tools.md#ac3)

<!-- sc: S4 -->
### S4

Malformed lines are skipped and counted; an empty spool returns an empty result.

<!-- field: traces -->
#### Traces

- [srd044-spool-query-tools](../software-requirements/srd044-spool-query-tools.md) [AC4](../software-requirements/srd044-spool-query-tools.md#ac4)

<!-- field: failure_scenarios -->
## Failure scenarios

| id | scenario | expected |
| --- | --- | --- |
| E1 | A rotated file is missing between two existing generations. | The tool reads whichever files exist without error; the gap is invisible in the output. |
| E2 | The spool path does not exist. | An empty result with zero traces and zero skipped lines. |
| E3 | The trace ID does not match any spooled span. | An empty spans array and span_count of zero. |

<!-- field: out_of_scope -->
## Out of scope

- Mutating or rotating spool files.
- Full-text search, attribute filtering, or streaming.
- HTTP query endpoints (separate REST surface concern).

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel12.0-spool-query](../test-suites/test-rel12.0-spool-query.md)
