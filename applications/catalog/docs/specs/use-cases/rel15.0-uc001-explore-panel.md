---
id: rel15.0-uc001-explore-panel
title: An operator narrows a latency band to a responsible attribute and lands on a trace
status: planned
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An operator opens the collector Explore page. A duration-over-time heatmap
renders from /query/spans/stats. The operator drag-selects a slow band; the
page fetches /query/spans/breakdown for that region and shows the ranked
attributes that distinguish the selection from the rest, plus exemplar spans.
Clicking an exemplar opens the existing per-trace waterfall. The whole loop is
aggregate to selection to explanation to a single trace, same-origin with the
query API.

<!-- field: actor -->
## Actor

Operator browsing a deployed collector

<!-- field: trigger -->
## Trigger

The collector is running with spooled spans, and the operator opens the
Explore page at the collector query origin.

<!-- field: flow -->
## Flow

- F1: The Explore page requests GET /query/spans/stats with an optional filter; the collector runs declared spool_span_heatmap and spool_span_group_by actions, composes their results, and returns a duration-over-time heatmap with axis boundaries and cell counts.
- F2: The page renders the heatmap; the operator drag-selects a rectangular region (a time span and a duration band).
- F3: The page requests GET /query/spans/breakdown with the selection as the selection filter and the current view as the baseline; the collector runs spool_span_breakdown and returns the ranked attribute divergence.
- F4: The page shows the ranked attributes (value, inside/outside counts, score) and a list of exemplar spans from the selection.
- F5: The operator clicks an exemplar span; the page navigates to the per-trace waterfall (R6) for that trace id.
- F6: With an empty spool, the heatmap and breakdown render an explicit empty state rather than an error.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd020-collector](../software-requirements/srd020-collector.md) [R8](../software-requirements/srd020-collector.md#r8---explore-query-surface) | /query/spans/stats and /query/spans/breakdown routes |
| T2 | [srd020-collector](../software-requirements/srd020-collector.md) [R9](../software-requirements/srd020-collector.md#r9---explore-ui) | the Explore page, drag-select, ranked divergence, and waterfall links |
| T3 | [srd020-collector](../software-requirements/srd020-collector.md) [R9](../software-requirements/srd020-collector.md#r9---explore-ui) | the span-analytics words the routes bind, declared as a platform dependency on agent-core srd045-span-analytics-tools |
| T4 | [srd020-collector](../software-requirements/srd020-collector.md) [R6](../software-requirements/srd020-collector.md#r6---trace-query-ui) | the per-trace waterfall the exemplars link into |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

GET /query/spans/stats returns the pinned heatmap and group-by contract; the Explore page renders the heatmap from it.

<!-- field: traces -->
#### Traces

- [srd020-collector](../software-requirements/srd020-collector.md) [AC6](../software-requirements/srd020-collector.md#ac6)
- [srd020-collector](../software-requirements/srd020-collector.md) [AC8](../software-requirements/srd020-collector.md#ac8)

<!-- sc: S2 -->
### S2

A drag-selected region drives GET /query/spans/breakdown and the page shows the ranked attribute divergence for that selection.

<!-- field: traces -->
#### Traces

- [srd020-collector](../software-requirements/srd020-collector.md) [AC6](../software-requirements/srd020-collector.md#ac6)
- [srd020-collector](../software-requirements/srd020-collector.md) [AC8](../software-requirements/srd020-collector.md#ac8)

<!-- sc: S3 -->
### S3

An exemplar span links into the existing per-trace waterfall.

<!-- field: traces -->
#### Traces

- [srd020-collector](../software-requirements/srd020-collector.md) [AC8](../software-requirements/srd020-collector.md#ac8)

<!-- sc: S4 -->
### S4

An empty or absent spool yields a well-formed empty stats and breakdown response and an explicit empty UI state.

<!-- field: traces -->
#### Traces

- [srd020-collector](../software-requirements/srd020-collector.md) [AC7](../software-requirements/srd020-collector.md#ac7)

<!-- field: failure_scenarios -->
## Failure scenarios

| id | scenario | expected |
| --- | --- | --- |
| E1 | The drag-selected region contains no span. | The breakdown returns an empty ranking and the page shows an explicit empty state. |
| E2 | The spool is absent when the page loads. | The stats route returns an empty heatmap and the page renders its empty state, not an error banner. |

<!-- field: out_of_scope -->
## Out of scope

- The agent-core span-analytics tools themselves (agent-core rel13.0).
- Metrics analytics or Dolt scraping.

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel15.0-explore](../test-suites/test-rel15.0-explore.md)
