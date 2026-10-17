---
id: test-rel13.0-span-analytics
title: Test Suite for Release 13.0 Span Analytics Tools
release: "13.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- [rel13.0-uc001-span-analytics](../use-cases/rel13.0-uc001-span-analytics.md)

<!-- field: preconditions -->
## Preconditions

- srd045-span-analytics-tools defines spool_span_heatmap, spool_span_group_by, and spool_span_breakdown behavior.
- srd044-spool-query-tools provides the spool read path and span shape the analytics reuse.
- Tests create spool files in temporary directories and do not require a running OTLP receiver.

<!-- field: test_cases -->
## Test cases

<!-- case: Conjunctive filter selects exactly the matching spans | status: implemented -->
### Conjunctive filter selects exactly the matching spans

A spool of spans across two services and two span names. A filter naming
one service, one span name, and one attribute equality matches only the
spans satisfying every term; changing any term changes the matched set.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-span-analytics](../use-cases/rel13.0-uc001-span-analytics.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestSpanStatsFilterConjunction

<!-- field: inputs | form: "tree" -->
#### Inputs

- services: 2
- filter:
  - service: svc-a
  - span_name: handle
  - attribute:
    - phase: parse

<!-- field: expected | form: "tree" -->
#### Expected

- signal: SpanHeatmapReady
- matched: exact conjunction

<!-- field: traces -->
#### Traces

- [rel13.0-uc001-span-analytics](../use-cases/rel13.0-uc001-span-analytics.md) [S1](../use-cases/rel13.0-uc001-span-analytics.md#s1)
- [srd045-span-analytics-tools](../software-requirements/srd045-span-analytics-tools.md) [AC1](../software-requirements/srd045-span-analytics-tools.md#ac1)

<!-- case: Empty filter matches every span | status: implemented -->
### Empty filter matches every span

With no filter terms, the matched total equals the count of valid spans in
the spool.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-span-analytics](../use-cases/rel13.0-uc001-span-analytics.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestSpanStatsEmptyFilterMatchesAll

<!-- field: inputs | form: "tree" -->
#### Inputs

- filter: {}

<!-- field: expected | form: "tree" -->
#### Expected

- signal: SpanHeatmapReady
- matched: all valid spans

<!-- field: traces -->
#### Traces

- [srd045-span-analytics-tools](../software-requirements/srd045-span-analytics-tools.md) [AC1](../software-requirements/srd045-span-analytics-tools.md#ac1)

<!-- case: Heatmap grid cell counts sum to the matched total | status: implemented -->
### Heatmap grid cell counts sum to the matched total

Spans with known start times and durations bucket into a time-by-duration
grid. The sum of all cell counts equals the matched total, the reported
time and duration boundaries are monotonic, and a duration at or above the
last boundary lands in the overflow bucket.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-span-analytics](../use-cases/rel13.0-uc001-span-analytics.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestSpanStatsHeatmapGrid

<!-- field: inputs | form: "tree" -->
#### Inputs

- time_buckets: 4
- duration_boundaries_ms:
  1. 0
  1. 10
  1. 100
  1. 1000

<!-- field: expected | form: "tree" -->
#### Expected

- signal: SpanHeatmapReady
- cell_sum_equals_matched: true

<!-- field: traces -->
#### Traces

- [rel13.0-uc001-span-analytics](../use-cases/rel13.0-uc001-span-analytics.md) [S1](../use-cases/rel13.0-uc001-span-analytics.md#s1)
- [srd045-span-analytics-tools](../software-requirements/srd045-span-analytics-tools.md) [AC2](../software-requirements/srd045-span-analytics-tools.md#ac2)

<!-- case: Group-by counts are ordered and capped with a dropped remainder | status: implemented -->
### Group-by counts are ordered and capped with a dropped remainder

Spans across five distinct attribute values with different counts. A
group_by on that attribute with top_n 3 returns the three highest-count
values ordered by count then value, and reports two dropped groups plus
the span total they cover.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-span-analytics](../use-cases/rel13.0-uc001-span-analytics.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestSpanStatsGroupByTopN

<!-- field: inputs | form: "tree" -->
#### Inputs

- distinct_values: 5
- group_by: tool.name
- top_n: 3

<!-- field: expected | form: "tree" -->
#### Expected

- signal: SpanGroupByReady
- returned_groups: 3
- dropped_groups: 2

<!-- field: traces -->
#### Traces

- [rel13.0-uc001-span-analytics](../use-cases/rel13.0-uc001-span-analytics.md) [S2](../use-cases/rel13.0-uc001-span-analytics.md#s2)
- [srd045-span-analytics-tools](../software-requirements/srd045-span-analytics-tools.md) [AC3](../software-requirements/srd045-span-analytics-tools.md#ac3)

<!-- case: Heatmap and group-by are independently dispatchable | status: implemented -->
### Heatmap and group-by are independently dispatchable

The heatmap output contains no group table. A later group-by action reads
the original request seed from command state, returns no heatmap, and
preserves the requested grouping key.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-span-analytics](../use-cases/rel13.0-uc001-span-analytics.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestSpanAnalyticsWordsRemainIndependentAcrossMachineSteps

<!-- field: expected | form: "tree" -->
#### Expected

- heatmap_has_groups: false
- group_by_has_heatmap: false
- group_by_seed_recovered: true

<!-- field: traces -->
#### Traces

- [srd045-span-analytics-tools](../software-requirements/srd045-span-analytics-tools.md) [AC2](../software-requirements/srd045-span-analytics-tools.md#ac2)
- [srd045-span-analytics-tools](../software-requirements/srd045-span-analytics-tools.md) [AC3](../software-requirements/srd045-span-analytics-tools.md#ac3)

<!-- case: Breakdown ranks attributes that distinguish the selection | status: implemented -->
### Breakdown ranks attributes that distinguish the selection

A baseline of spans where one attribute value is concentrated in a slow
selection. spool_span_breakdown ranks that value first with inside and
outside counts, inside and outside proportions, and a deterministic score,
capped at top_n with a dropped count.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-span-analytics](../use-cases/rel13.0-uc001-span-analytics.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestSpanBreakdownDivergenceRanking

<!-- field: inputs | form: "tree" -->
#### Inputs

- baseline_filter: {}
- selection_filter:
  - min_duration_ms: 500
- top_n: 5

<!-- field: expected | form: "tree" -->
#### Expected

- signal: SpanBreakdownReady
- top_ranked: the concentrated value

<!-- field: traces -->
#### Traces

- [rel13.0-uc001-span-analytics](../use-cases/rel13.0-uc001-span-analytics.md) [S3](../use-cases/rel13.0-uc001-span-analytics.md#s3)
- [srd045-span-analytics-tools](../software-requirements/srd045-span-analytics-tools.md) [AC4](../software-requirements/srd045-span-analytics-tools.md#ac4)

<!-- case: Empty selection yields an empty ranking | status: implemented -->
### Empty selection yields an empty ranking

A selection filter that matches no span produces SpanBreakdownReady with
an empty ranking rather than an error.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-span-analytics](../use-cases/rel13.0-uc001-span-analytics.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestSpanBreakdownEmptySelection

<!-- field: inputs | form: "tree" -->
#### Inputs

- selection_filter:
  - service: nonexistent

<!-- field: expected | form: "tree" -->
#### Expected

- signal: SpanBreakdownReady
- ranked: 0

<!-- field: traces -->
#### Traces

- [srd045-span-analytics-tools](../software-requirements/srd045-span-analytics-tools.md) [AC4](../software-requirements/srd045-span-analytics-tools.md#ac4)

<!-- case: Malformed lines are skipped, absent spool is empty, output is deterministic | status: implemented -->
### Malformed lines are skipped, absent spool is empty, output is deterministic

A spool with two malformed lines and valid spans reports skipped_lines 2;
an absent spool returns an empty result; two runs over the same spool and
config return byte-identical output.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-span-analytics](../use-cases/rel13.0-uc001-span-analytics.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestSpanStatsMalformedAndDeterministic

<!-- field: inputs | form: "tree" -->
#### Inputs

- malformed_lines: 2

<!-- field: expected | form: "tree" -->
#### Expected

- signal: SpanHeatmapReady
- skipped_lines: 2
- deterministic: true

<!-- field: traces -->
#### Traces

- [srd045-span-analytics-tools](../software-requirements/srd045-span-analytics-tools.md) [AC5](../software-requirements/srd045-span-analytics-tools.md#ac5)

<!-- case: Missing spool path emits CommandError | status: implemented -->
### Missing spool path emits CommandError

Both spool_span_heatmap and spool_span_group_by with an empty path config
emit CommandError before aggregation runs.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-span-analytics](../use-cases/rel13.0-uc001-span-analytics.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestSpanStatsEmptyPathError

<!-- field: inputs | form: "tree" -->
#### Inputs

- path: ""

<!-- field: expected | form: "tree" -->
#### Expected

- signal: CommandError

<!-- field: traces -->
#### Traces

- [srd045-span-analytics-tools](../software-requirements/srd045-span-analytics-tools.md) [AC5](../software-requirements/srd045-span-analytics-tools.md#ac5)

<!-- case: Selection exemplars are capped, deterministic, and deduplicated | status: implemented -->
### Selection exemplars are capped, deterministic, and deduplicated

spool_span_heatmap, spool_span_group_by, and spool_span_breakdown return
exemplar_trace_ids drawn from matched and selection sets; the list is
deduplicated, ordered by ascending trace ID, capped at exemplar_cap, and
excludes spans with no trace ID.

<!-- field: use_case -->
#### Use case

[rel13.0-uc001-span-analytics](../use-cases/rel13.0-uc001-span-analytics.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run Exemplar

<!-- field: inputs | form: "tree" -->
#### Inputs

- exemplar_cap: 2

<!-- field: expected | form: "tree" -->
#### Expected

- exemplar_trace_ids_deterministic: true

<!-- field: traces -->
#### Traces

- [srd045-span-analytics-tools](../software-requirements/srd045-span-analytics-tools.md) [AC6](../software-requirements/srd045-span-analytics-tools.md#ac6)
