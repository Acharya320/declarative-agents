---
id: rel14.0-uc001-served-trace-ui
title: The collector serves its own trace UI at the query origin
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An operator opens the collector's query URL in a browser and lands on the
trace list, then clicks through to a per-trace timeline. The UI is served by
the collector_query server from the shipped ui/dist tree, same-origin with
/query/*, so no separate UI workload is deployed. A collector built without
its UI answers an explicit error rather than a blank page. The UI is composed
on the ui-kit shell from ui/ui.yaml and reads /query/* through the kit client
(applications srd004 R5, R6.1).

<!-- field: actor -->
## Actor

Operator browsing a deployed collector, or the conformance harness

<!-- field: trigger -->
## Trigger

The collector profile is running with its query server bound, and a client
issues GET / (or a deep link such as /traces/{trace_id}) against the query
origin.

<!-- field: flow -->
## Flow

- F1: The client issues GET / against the collector_query origin; the server returns the SPA index document from the configured ui/dist root.
- F2: The SPA loads, the kit shell mounts the traces panel as the default, and the panel calls GET /query/traces on the same origin through the kit client; the server returns the paginated trace summaries.
- F3: The operator selects a trace; the traces panel pushes /traces/{trace_id}.
- F4: The client issues GET /traces/{trace_id}; the static_assets endpoint resolves the deep link to the SPA index (spa true), and the traces panel reads the trace id from the path, calls GET /query/traces/{trace_id}, and draws it with the kit trace timeline.
- F5: The /query/* endpoints keep answering their JSON contracts unchanged; the static_assets catch-all never intercepts a query path.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd020-collector R7 | UI served by the collector_query server via static_assets, same-origin with the query API |
| T2 | srd020-collector R5 | the query surface the UI reads for trace list and detail |
| T3 | srd020-collector R6 | the shipped ui/ SPA and its uiDist reproducibility gate |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

GET / on the query origin returns the SPA index with an HTML content type.

<!-- field: traces -->
#### Traces

- srd020-collector AC4

<!-- sc: S2 -->
### S2

A deep link such as GET /traces/{trace_id} resolves to the SPA index rather than a 404, so the traces panel can render it.

<!-- field: traces -->
#### Traces

- srd020-collector AC4

<!-- sc: S3 -->
### S3

Serving the UI leaves every /query/* response body and status unchanged.

<!-- field: traces -->
#### Traces

- srd020-collector AC4

<!-- sc: S4 -->
### S4

The UI root is a literal profile-declared path with no new environment variable; a missing UI root yields an explicit error status, not a silent empty body.

<!-- field: traces -->
#### Traces

- srd020-collector AC5

<!-- sc: S5 -->
### S5

The UI composes on the ui-kit shell from ui/ui.yaml, registers every local panel it declares, and reads /query/* through the kit client with no raw fetch or EventSource.

<!-- field: traces -->
#### Traces

- srd020-collector R6.1

<!-- field: out_of_scope -->
## Out of scope

- The Explore/analytics UI and span-stats query routes (a later release).
- Helm delivery of the UI to application charts (a sibling sub-issue).
- Metrics UI or Dolt scraping.

<!-- field: test_suite_trace -->
## Test suite trace

test-rel14.0-served-ui
