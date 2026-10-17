---
id: test-rel15.0-explore
title: Test Suite for Release 15.0 Collector Explore Panel
release: "15.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- [rel15.0-uc001-explore-panel](../use-cases/rel15.0-uc001-explore-panel.md)

<!-- field: preconditions -->
## Preconditions

- srd020-collector R8 defines the /query/spans/stats and /query/spans/breakdown routes.
- srd020-collector R9 defines the Explore UI page.
- srd045-span-analytics-tools provides the spool_span_heatmap, spool_span_group_by, and spool_span_breakdown tools the routes bind.
- The agent-core binary is built and available via RequireCoreRoot.

<!-- field: test_cases -->
## Test cases

<!-- case: Stats route returns the pinned heatmap and group-by contract | status: implemented -->
### Stats route returns the pinned heatmap and group-by contract

Seed a spool, GET /query/spans/stats with a filter and heatmap request,
and assert the pinned JSON keys: time bucket boundaries, duration bucket
boundaries, heatmap cells, matched total, skipped lines, and groups.

<!-- field: use_case -->
#### Use case

[rel15.0-uc001-explore-panel](../use-cases/rel15.0-uc001-explore-panel.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCollectorSpanStatsContract

<!-- field: inputs | form: "tree" -->
#### Inputs

- spool: seeded fixture
- request: GET /query/spans/stats

<!-- field: expected | form: "tree" -->
#### Expected

- status: 200
- cell_sum_equals_matched: true

<!-- field: traces -->
#### Traces

- [rel15.0-uc001-explore-panel](../use-cases/rel15.0-uc001-explore-panel.md) [S1](../use-cases/rel15.0-uc001-explore-panel.md#s1)
- [srd020-collector](../software-requirements/srd020-collector.md) [AC6](../software-requirements/srd020-collector.md#ac6)

<!-- case: Breakdown route returns the pinned divergence contract | status: implemented -->
### Breakdown route returns the pinned divergence contract

Seed a spool where one attribute value concentrates in a slow region. GET
/query/spans/breakdown with a selection filter returns the ranked
divergence with pinned keys: value, inside count, outside count,
proportions, score, and the dropped count.

<!-- field: use_case -->
#### Use case

[rel15.0-uc001-explore-panel](../use-cases/rel15.0-uc001-explore-panel.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCollectorSpanBreakdownContract

<!-- field: inputs | form: "tree" -->
#### Inputs

- spool: seeded fixture
- request: GET /query/spans/breakdown

<!-- field: expected | form: "tree" -->
#### Expected

- status: 200
- top_ranked: the concentrated value

<!-- field: traces -->
#### Traces

- [rel15.0-uc001-explore-panel](../use-cases/rel15.0-uc001-explore-panel.md) [S2](../use-cases/rel15.0-uc001-explore-panel.md#s2)
- [srd020-collector](../software-requirements/srd020-collector.md) [AC6](../software-requirements/srd020-collector.md#ac6)

<!-- case: Result caps come from tool config, not the request | status: implemented -->
### Result caps come from tool config, not the request

With max_top_n patched to 1, a request for top_n 1000 over the two-service
fixture still returns one group and reports one dropped group, proving the
configured cap beats the request value.

<!-- field: use_case -->
#### Use case

[rel15.0-uc001-explore-panel](../use-cases/rel15.0-uc001-explore-panel.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCollectorExploreCapsFromConfig

<!-- field: inputs | form: "tree" -->
#### Inputs

- requested_top_n: 1000
- configured_cap: 1

<!-- field: expected | form: "tree" -->
#### Expected

- status: 200
- returned_groups: 1
- dropped_groups: 1

<!-- field: traces -->
#### Traces

- [srd020-collector](../software-requirements/srd020-collector.md) [AC6](../software-requirements/srd020-collector.md#ac6)

<!-- case: Empty spool yields well-formed empty Explore responses | status: implemented -->
### Empty spool yields well-formed empty Explore responses

With no spooled spans, GET /query/spans/stats returns an empty heatmap
(zero cells, zero matched) and GET /query/spans/breakdown returns an empty
ranking, both 200, not an error.

<!-- field: use_case -->
#### Use case

[rel15.0-uc001-explore-panel](../use-cases/rel15.0-uc001-explore-panel.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCollectorExploreEmptySpool

<!-- field: inputs | form: "tree" -->
#### Inputs

- spool: absent

<!-- field: expected | form: "tree" -->
#### Expected

- status: 200
- matched: 0
- ranked: 0

<!-- field: traces -->
#### Traces

- [rel15.0-uc001-explore-panel](../use-cases/rel15.0-uc001-explore-panel.md) [S4](../use-cases/rel15.0-uc001-explore-panel.md#s4)
- [srd020-collector](../software-requirements/srd020-collector.md) [AC7](../software-requirements/srd020-collector.md#ac7)

<!-- case: Explore page is wired into the UI descriptor | status: implemented -->
### Explore page is wired into the UI descriptor

The collector ui/ui.yaml declares an /explore panel and the actions that
bind the stats and breakdown machine requests, so the served SPA exposes
the Explore page.

<!-- field: use_case -->
#### Use case

[rel15.0-uc001-explore-panel](../use-cases/rel15.0-uc001-explore-panel.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCollectorExploreRouteDeclared

<!-- field: inputs | form: "tree" -->
#### Inputs

- file: agents/collector/ui/ui.yaml

<!-- field: expected | form: "tree" -->
#### Expected

- explore_route_declared: true

<!-- field: traces -->
#### Traces

- [rel15.0-uc001-explore-panel](../use-cases/rel15.0-uc001-explore-panel.md) [S1](../use-cases/rel15.0-uc001-explore-panel.md#s1)
- [srd020-collector](../software-requirements/srd020-collector.md) [AC8](../software-requirements/srd020-collector.md#ac8)

<!-- case: Explore UI reproduces from a clean source build | status: implemented -->
### Explore UI reproduces from a clean source build

The tracked collector ui/dist, including the Explore page, matches the
output of a clean npm ci + build from the lockfile-pinned source.

<!-- field: use_case -->
#### Use case

[rel15.0-uc001-explore-panel](../use-cases/rel15.0-uc001-explore-panel.md)

<!-- field: integration_target -->
#### Integration target

mage uiDist

<!-- field: inputs | form: "tree" -->
#### Inputs

- ui_dir: agents/collector/ui

<!-- field: expected | form: "tree" -->
#### Expected

- dist_diff: none

<!-- field: traces -->
#### Traces

- [srd020-collector](../software-requirements/srd020-collector.md) [AC8](../software-requirements/srd020-collector.md#ac8)
