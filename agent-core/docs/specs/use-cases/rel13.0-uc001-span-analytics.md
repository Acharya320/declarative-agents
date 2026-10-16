---
id: rel13.0-uc001-span-analytics
title: Aggregate reads turn a spool into a latency heatmap and a divergence ranking
status: planned
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A consumer drives the span-analytics tools over a seeded NDJSON spool. It
filters spans by time, service, name, and attribute equality; buckets the
filtered spans into a duration-over-time heatmap; counts them grouped by an
attribute with a reported dropped remainder; and ranks the attributes that
most distinguish a selected slow region from its complement. Every read is
side-effect free, tolerant of malformed lines, and deterministic.

<!-- field: actor -->
## Actor

The collector query surface (or a test harness) invoking the builtin analytics tools

<!-- field: trigger -->
## Trigger

A batch of spans has been spooled, and a consumer requests aggregate
statistics or an attribute breakdown over a filtered region of the spool.

<!-- field: flow -->
## Flow

- F1: The consumer invokes spool_span_heatmap with a filter (time range, service, span name, attribute equality) and a heatmap request; the tool reads the active and rotated spool files and returns the time-by-duration grid with axis boundaries, matched total, and skipped-line count.
- F2: The consumer independently invokes spool_span_group_by with a group_by key; the tool returns per-value counts ordered by count then value, capped at top-N, with the dropped group and span totals.
- F3: The consumer invokes spool_span_breakdown with a baseline filter and a selection filter; the tool partitions baseline spans into inside and outside the selection and ranks attribute values by divergence, capped at top-N with a dropped count.
- F4: A malformed line in the spool is skipped and counted; an absent spool yields an empty result rather than an error.
- F5: A missing or unreadable spool path emits CommandError before any aggregation runs.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd045-span-analytics-tools R1.1, R1.2, R1.3, R1.4 | conjunctive span filter over the srd044 read path |
| T2 | srd045-span-analytics-tools R2.1, R2.2, R2.3, R2.4, R3.1, R3.2, R3.3 | independently dispatchable duration heatmap and group-by counts |
| T3 | srd045-span-analytics-tools R4.1, R4.2, R4.3, R4.4 | attribute divergence in spool_span_breakdown |
| T4 | srd044-spool-query-tools | spool discovery, NDJSON parsing, span shape reused |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

A conjunctive filter selects exactly the matching spans and the heatmap cell counts sum to the matched total.

<!-- field: traces -->
#### Traces

- srd045-span-analytics-tools AC1
- srd045-span-analytics-tools AC2

<!-- sc: S2 -->
### S2

Group-by counts are ordered and capped at top-N with the dropped group and span totals reported.

<!-- field: traces -->
#### Traces

- srd045-span-analytics-tools AC3

<!-- sc: S3 -->
### S3

The breakdown ranks attribute values by a deterministic divergence between the selection and its complement, capped with a dropped count.

<!-- field: traces -->
#### Traces

- srd045-span-analytics-tools AC4

<!-- sc: S4 -->
### S4

Malformed lines are skipped and counted, an absent spool is empty, and a missing or unreadable path emits CommandError.

<!-- field: traces -->
#### Traces

- srd045-span-analytics-tools AC5

<!-- field: failure_scenarios -->
## Failure scenarios

| id | scenario | expected |
| --- | --- | --- |
| E1 | The selection filter matches no span. | spool_span_breakdown returns an empty ranking, not an error. |
| E2 | The spool path config is empty. | The tool emits CommandError before reading any file. |

<!-- field: out_of_scope -->
## Out of scope

- The HTTP query routes that expose these tools (catalog srd020, rel15.0).
- The Explore UI page that renders the heatmap and ranking (catalog rel15.0).

<!-- field: test_suite_trace -->
## Test suite trace

test-rel13.0-span-analytics
