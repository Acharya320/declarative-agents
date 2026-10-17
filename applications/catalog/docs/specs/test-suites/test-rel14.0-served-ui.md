---
id: test-rel14.0-served-ui
title: Test Suite for Release 14.0 Served Collector Trace UI
release: "14.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- [rel14.0-uc001-served-trace-ui](../use-cases/rel14.0-uc001-served-trace-ui.md)

<!-- field: preconditions -->
## Preconditions

- srd020-collector R7 defines the served UI on the collector_query server.
- The collector family ships a built ui/dist tree (srd020-collector R6).
- The agent-core binary is built and available via RequireCoreRoot.

<!-- field: test_cases -->
## Test cases

<!-- case: Query origin serves the SPA index at root | status: implemented -->
### Query origin serves the SPA index at root

Launch the collector and GET / on the query origin. Assert a 200 with an
HTML content type and the SPA index document body.

<!-- field: use_case -->
#### Use case

[rel14.0-uc001-served-trace-ui](../use-cases/rel14.0-uc001-served-trace-ui.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCollectorServesUIIndex

<!-- field: inputs | form: "tree" -->
#### Inputs

- request: GET /

<!-- field: expected | form: "tree" -->
#### Expected

- status: 200
- content_type: text/html

<!-- field: traces -->
#### Traces

- [rel14.0-uc001-served-trace-ui](../use-cases/rel14.0-uc001-served-trace-ui.md) [S1](../use-cases/rel14.0-uc001-served-trace-ui.md#s1)
- [srd020-collector](../software-requirements/srd020-collector.md) [AC4](../software-requirements/srd020-collector.md#ac4)

<!-- case: Deep link resolves to the SPA index | status: implemented -->
### Deep link resolves to the SPA index

GET /traces/trace-aaa (a traces-panel path, not a file) returns the SPA
index document with a 200, so the client can render the trace timeline.

<!-- field: use_case -->
#### Use case

[rel14.0-uc001-served-trace-ui](../use-cases/rel14.0-uc001-served-trace-ui.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCollectorServesUIDeepLink

<!-- field: inputs | form: "tree" -->
#### Inputs

- request: GET /traces/trace-aaa

<!-- field: expected | form: "tree" -->
#### Expected

- status: 200
- content_type: text/html
- body: SPA index

<!-- field: traces -->
#### Traces

- [rel14.0-uc001-served-trace-ui](../use-cases/rel14.0-uc001-served-trace-ui.md) [S2](../use-cases/rel14.0-uc001-served-trace-ui.md#s2)
- [srd020-collector](../software-requirements/srd020-collector.md) [AC4](../software-requirements/srd020-collector.md#ac4)

<!-- case: Query API is unaffected by the static_assets catch-all | status: implemented -->
### Query API is unaffected by the static_assets catch-all

With the UI served, GET /query/traces still returns the pinned JSON
contract and GET /query/traces/{trace_id} still returns span detail; the
static_assets catch-all never intercepts a /query/* path.

<!-- field: use_case -->
#### Use case

[rel14.0-uc001-served-trace-ui](../use-cases/rel14.0-uc001-served-trace-ui.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCollectorQueryUnaffectedByUI

<!-- field: inputs | form: "tree" -->
#### Inputs

- spool: seeded two-trace fixture
- requests:
  1. GET /query/traces
  1. GET /query/traces/trace-aaa

<!-- field: expected | form: "tree" -->
#### Expected

- status: 200
- summary_keys:
  1. duration_ms
  1. root_service
  1. root_span_name
  1. span_count
  1. start_time
  1. trace_id

<!-- field: traces -->
#### Traces

- [rel14.0-uc001-served-trace-ui](../use-cases/rel14.0-uc001-served-trace-ui.md) [S3](../use-cases/rel14.0-uc001-served-trace-ui.md#s3)
- [srd020-collector](../software-requirements/srd020-collector.md) [AC4](../software-requirements/srd020-collector.md#ac4)

<!-- case: Missing UI root yields an explicit error | status: implemented -->
### Missing UI root yields an explicit error

Launch the collector with the static_assets root pointing at an absent
directory. GET / answers an explicit HTTP error status, not a 200 with an
empty body, and /query/* still works.

<!-- field: use_case -->
#### Use case

[rel14.0-uc001-served-trace-ui](../use-cases/rel14.0-uc001-served-trace-ui.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCollectorMissingUIRoot

<!-- field: inputs | form: "tree" -->
#### Inputs

- ui_root: absent directory
- request: GET /

<!-- field: expected | form: "tree" -->
#### Expected

- status: explicit error (not 200, not empty body)

<!-- field: traces -->
#### Traces

- [rel14.0-uc001-served-trace-ui](../use-cases/rel14.0-uc001-served-trace-ui.md) [S4](../use-cases/rel14.0-uc001-served-trace-ui.md#s4)
- [srd020-collector](../software-requirements/srd020-collector.md) [AC5](../software-requirements/srd020-collector.md#ac5)

<!-- case: UI root is a literal profile path with no new environment variable | status: implemented -->
### UI root is a literal profile path with no new environment variable

The collector rest.yaml static_assets root is a literal path resolved
relative to the profile directory, with no ${VAR} expansion. Assert the
shipped rest.yaml declares no new environment variable for the UI root.

<!-- field: use_case -->
#### Use case

[rel14.0-uc001-served-trace-ui](../use-cases/rel14.0-uc001-served-trace-ui.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCollectorUIRootLiteral

<!-- field: inputs | form: "tree" -->
#### Inputs

- file: agents/collector/rest.yaml

<!-- field: expected | form: "tree" -->
#### Expected

- new_ui_env_var: none

<!-- field: traces -->
#### Traces

- [srd020-collector](../software-requirements/srd020-collector.md) [AC5](../software-requirements/srd020-collector.md#ac5)

<!-- case: Served UIs compose on the kit shell and read through the kit client | status: implemented -->
### Served UIs compose on the kit shell and read through the kit client

For the collector and bench UIs, ui/ui.yaml is version 2, every local
panel it declares has an entry in src/panels.ts, src/main.tsx mounts the
kit AppShell from virtual:ui-config, and src/ opens no EventSource and
issues no raw fetch beyond the bench's one POST launch, which the read-only
kit client cannot carry.

<!-- field: use_case -->
#### Use case

[rel14.0-uc001-served-trace-ui](../use-cases/rel14.0-uc001-served-trace-ui.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestServedUIsComposeOnKitShell

<!-- field: inputs | form: "tree" -->
#### Inputs

- ui_dirs:
  1. agents/collector/ui
  1. agents/bench/ui

<!-- field: expected | form: "tree" -->
#### Expected

- ui_yaml_version: 2
- unregistered_panels: none
- raw_fetch:
  - collector: 0
  - bench: 1

<!-- field: traces -->
#### Traces

- [rel14.0-uc001-served-trace-ui](../use-cases/rel14.0-uc001-served-trace-ui.md) [S5](../use-cases/rel14.0-uc001-served-trace-ui.md#s5)
- [srd020-collector](../software-requirements/srd020-collector.md) [R6.1](../software-requirements/srd020-collector.md#r61)
