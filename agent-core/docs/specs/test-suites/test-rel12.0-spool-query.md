---
id: test-rel12.0-spool-query
title: Test Suite for Release 12.0 Spool Query Tools
release: "12.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- [rel12.0-uc001-spool-query](../use-cases/rel12.0-uc001-spool-query.md)

<!-- field: preconditions -->
## Preconditions

- srd044-spool-query-tools defines spool_list_traces and spool_get_trace behavior.
- srd042-otlp-receiver-tools defines the NDJSON spool shape and rotation naming that the query tools read.
- Tests create spool files in temporary directories and do not require a running OTLP receiver.

<!-- field: test_cases -->
## Test cases

<!-- case: List traces aggregates across active and rotated spool files | status: implemented -->
### List traces aggregates across active and rotated spool files

One span in a rotated file and two spans (same trace) in the active file
produce two trace summaries. The newer trace appears first. Root service
and root span name match the root span of each trace. Skipped lines is
zero.

<!-- field: use_case -->
#### Use case

[rel12.0-uc001-spool-query](../use-cases/rel12.0-uc001-spool-query.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestListTracesAcrossActiveAndRotated

<!-- field: inputs | form: "tree" -->
#### Inputs

- rotated_file_spans: 1
- active_file_spans: 2
- traces: 2

<!-- field: expected | form: "tree" -->
#### Expected

- signal: TracesListed
- total: 2
- first_trace_span_count: 2
- skipped_lines: 0

<!-- field: traces -->
#### Traces

- [rel12.0-uc001-spool-query](../use-cases/rel12.0-uc001-spool-query.md) [S1](../use-cases/rel12.0-uc001-spool-query.md#s1)
- [srd044-spool-query-tools](../software-requirements/srd044-spool-query-tools.md) [AC1](../software-requirements/srd044-spool-query-tools.md#ac1)

<!-- case: Pagination returns the correct page with offset and page size | status: implemented -->
### Pagination returns the correct page with offset and page size

Five single-span traces with increasing start times. Page size 2 and
offset 1 returns two traces starting from the second-newest. Total
remains 5.

<!-- field: use_case -->
#### Use case

[rel12.0-uc001-spool-query](../use-cases/rel12.0-uc001-spool-query.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestListTracesPagination

<!-- field: inputs | form: "tree" -->
#### Inputs

- total_traces: 5
- page_size: 2
- offset: 1

<!-- field: expected | form: "tree" -->
#### Expected

- signal: TracesListed
- total: 5
- returned: 2

<!-- field: traces -->
#### Traces

- [rel12.0-uc001-spool-query](../use-cases/rel12.0-uc001-spool-query.md) [S2](../use-cases/rel12.0-uc001-spool-query.md#s2)
- [srd044-spool-query-tools](../software-requirements/srd044-spool-query-tools.md) [AC2](../software-requirements/srd044-spool-query-tools.md#ac2)

<!-- case: Page size is capped at max_page_size | status: implemented -->
### Page size is capped at max_page_size

A declared page_size of 50 with max_page_size 5 results in an effective
page size of 5.

<!-- field: use_case -->
#### Use case

[rel12.0-uc001-spool-query](../use-cases/rel12.0-uc001-spool-query.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestListTracesPageSizeCap

<!-- field: inputs | form: "tree" -->
#### Inputs

- page_size: 50
- max_page_size: 5

<!-- field: expected | form: "tree" -->
#### Expected

- signal: TracesListed
- effective_page_size: 5

<!-- field: traces -->
#### Traces

- [srd044-spool-query-tools](../software-requirements/srd044-spool-query-tools.md) [AC2](../software-requirements/srd044-spool-query-tools.md#ac2)

<!-- case: Empty or absent spool returns an empty result | status: implemented -->
### Empty or absent spool returns an empty result

A spool path that does not exist produces TracesListed with zero traces
and zero total, not an error.

<!-- field: use_case -->
#### Use case

[rel12.0-uc001-spool-query](../use-cases/rel12.0-uc001-spool-query.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestListTracesEmptySpool

<!-- field: inputs | form: "tree" -->
#### Inputs

- spool_exists: false

<!-- field: expected | form: "tree" -->
#### Expected

- signal: TracesListed
- total: 0

<!-- field: traces -->
#### Traces

- [rel12.0-uc001-spool-query](../use-cases/rel12.0-uc001-spool-query.md) [S4](../use-cases/rel12.0-uc001-spool-query.md#s4)
- [srd044-spool-query-tools](../software-requirements/srd044-spool-query-tools.md) [AC4](../software-requirements/srd044-spool-query-tools.md#ac4)

<!-- case: Malformed NDJSON lines are skipped and counted | status: implemented -->
### Malformed NDJSON lines are skipped and counted

Two malformed lines followed by one valid span line produce one trace
summary with skipped_lines equal to 2.

<!-- field: use_case -->
#### Use case

[rel12.0-uc001-spool-query](../use-cases/rel12.0-uc001-spool-query.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestListTracesMalformedLineTolerance

<!-- field: inputs | form: "tree" -->
#### Inputs

- malformed_lines: 2
- valid_spans: 1

<!-- field: expected | form: "tree" -->
#### Expected

- signal: TracesListed
- total: 1
- skipped_lines: 2

<!-- field: traces -->
#### Traces

- [srd044-spool-query-tools](../software-requirements/srd044-spool-query-tools.md) [AC4](../software-requirements/srd044-spool-query-tools.md#ac4)

<!-- case: Get trace returns span details with waterfall fields | status: implemented -->
### Get trace returns span details with waterfall fields

Three spans across two traces. Querying trace t1 returns two spans
sorted by start time with correct span id, parent span id, service,
name, and timing. The span from trace t2 is excluded.

<!-- field: use_case -->
#### Use case

[rel12.0-uc001-spool-query](../use-cases/rel12.0-uc001-spool-query.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestGetTraceReturnsSpanDetails

<!-- field: inputs | form: "tree" -->
#### Inputs

- total_spans: 3
- traces: 2
- query_trace: t1

<!-- field: expected | form: "tree" -->
#### Expected

- signal: TraceRetrieved
- span_count: 2
- first_span_name: root
- first_span_parent: ""

<!-- field: traces -->
#### Traces

- [rel12.0-uc001-spool-query](../use-cases/rel12.0-uc001-spool-query.md) [S3](../use-cases/rel12.0-uc001-spool-query.md#s3)
- [srd044-spool-query-tools](../software-requirements/srd044-spool-query-tools.md) [AC3](../software-requirements/srd044-spool-query-tools.md#ac3)

<!-- case: Get trace reads across rotated files | status: implemented -->
### Get trace reads across rotated files

One span in a rotated file and one span in the active file for the same
trace. Both appear in the result sorted by start time.

<!-- field: use_case -->
#### Use case

[rel12.0-uc001-spool-query](../use-cases/rel12.0-uc001-spool-query.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestGetTraceAcrossRotatedFiles

<!-- field: inputs | form: "tree" -->
#### Inputs

- rotated_spans: 1
- active_spans: 1
- same_trace: true

<!-- field: expected | form: "tree" -->
#### Expected

- signal: TraceRetrieved
- span_count: 2

<!-- field: traces -->
#### Traces

- [srd044-spool-query-tools](../software-requirements/srd044-spool-query-tools.md) [AC1](../software-requirements/srd044-spool-query-tools.md#ac1)
- [srd044-spool-query-tools](../software-requirements/srd044-spool-query-tools.md) [AC3](../software-requirements/srd044-spool-query-tools.md#ac3)

<!-- case: Unknown trace ID returns empty spans | status: implemented -->
### Unknown trace ID returns empty spans

Querying a trace id that does not exist in the spool returns
TraceRetrieved with zero spans.

<!-- field: use_case -->
#### Use case

[rel12.0-uc001-spool-query](../use-cases/rel12.0-uc001-spool-query.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestGetTraceUnknownTraceID

<!-- field: inputs | form: "tree" -->
#### Inputs

- query_trace: nonexistent

<!-- field: expected | form: "tree" -->
#### Expected

- signal: TraceRetrieved
- span_count: 0

<!-- field: traces -->
#### Traces

- [srd044-spool-query-tools](../software-requirements/srd044-spool-query-tools.md) [AC3](../software-requirements/srd044-spool-query-tools.md#ac3)

<!-- case: Empty trace ID emits CommandError | status: implemented -->
### Empty trace ID emits CommandError

A spool_get_trace with an empty trace_id config emits CommandError
before reading any file.

<!-- field: use_case -->
#### Use case

[rel12.0-uc001-spool-query](../use-cases/rel12.0-uc001-spool-query.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestGetTraceEmptyTraceIDError

<!-- field: inputs | form: "tree" -->
#### Inputs

- trace_id: ""

<!-- field: expected | form: "tree" -->
#### Expected

- signal: CommandError

<!-- field: traces -->
#### Traces

- [srd044-spool-query-tools](../software-requirements/srd044-spool-query-tools.md) [AC5](../software-requirements/srd044-spool-query-tools.md#ac5)

<!-- case: Empty spool path emits CommandError | status: implemented -->
### Empty spool path emits CommandError

A spool_list_traces with an empty path config emits CommandError.

<!-- field: use_case -->
#### Use case

[rel12.0-uc001-spool-query](../use-cases/rel12.0-uc001-spool-query.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/otlp -run TestListTracesEmptyPathError

<!-- field: inputs | form: "tree" -->
#### Inputs

- path: ""

<!-- field: expected | form: "tree" -->
#### Expected

- signal: CommandError

<!-- field: traces -->
#### Traces

- [srd044-spool-query-tools](../software-requirements/srd044-spool-query-tools.md) [AC5](../software-requirements/srd044-spool-query-tools.md#ac5)
