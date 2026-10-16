---
id: rel01.0-uc001-demo-trace-collection
title: Demo Run Collects and Serves Traces
status: implemented
test_suite: test-rel01.0-demo-trace-collection
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

The demo launch path starts the canonical collector agent in spool mode
before the documentation-curator, exports OTLP traces from the curator to
the local collector, and stops the collector cleanly after the curator
exits. The presenter opens the collector query surface to show the trace
waterfall as part of the demo story.

<!-- field: actor -->
## Actor

A developer or presenter running the Agent Architecture.

<!-- field: trigger -->
## Trigger

Run mage run from the application root.

<!-- field: preconditions -->
## Preconditions

- The demo.yaml catalog_root resolves to a catalog checkout containing the collector and documentation-curator profiles.
- The demo.yaml core_root resolves to a buildable agent-core checkout.
- Four ephemeral loopback ports can be reserved for the local collector; curator ports (18081, 18082, 18084) are available.

<!-- field: flow -->
## Flow

- F1: mage run builds the agent-core binary and checks whether tracing is enabled (default true, disabled by tracing:false in demo.yaml).
- F2: When tracing is enabled, the run target reserves distinct ephemeral loopback receiver, control, monitor, and query ports, starts the collector agent in spool mode with the canonical collector profile, and prints the selected query URL.
- F3: The run target waits for the collector control endpoint to become healthy before proceeding.
- F4: The run target launches the documentation-curator with --otel-otlp-endpoint and --otel-service-name flags so the curator exports traces to the local collector.
- F5: The presenter opens the collector query surface to browse collected traces.
- F6: After the curator exits, the run target posts to the collector control exit route, waits for the collector process to terminate, and removes the spool directory.

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The collector agent starts before the curator and binds four reserved loopback endpoints that do not overlap another application collector.

<!-- field: traces -->
#### Traces

- architecture-agent-architecture design_decisions 7
- test-rel01.0-demo-trace-collection collector_command

<!-- sc: S2 -->
### S2

The curator run command includes OTLP export flags pointing at the local collector.

<!-- field: traces -->
#### Traces

- test-rel01.0-demo-trace-collection traced_curator_command

<!-- sc: S3 -->
### S3

Tracing is on by default and disabled by tracing:false in demo.yaml.

<!-- field: traces -->
#### Traces

- test-rel01.0-demo-trace-collection tracing_enabled

<!-- sc: S4 -->
### S4

The collector exits cleanly via its control route after the curator, and no listeners remain orphaned.

<!-- field: traces -->
#### Traces

- architecture-agent-architecture design_decisions 7
- test-rel01.0-demo-trace-collection collector_command

<!-- field: out_of_scope -->
## Out of scope

- Changing the canonical collector or documentation-curator profiles.
- Proving live spool content or query results (dependency-gated).
- Browser automation of the collector query surface.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | architecture-agent-architecture design_decisions 7 | S1 traces this requirement in the existing success criteria. |
| T2 | test-rel01.0-demo-trace-collection collector_command | S1 traces this requirement in the existing success criteria. |
| T3 | test-rel01.0-demo-trace-collection traced_curator_command | S2 traces this requirement in the existing success criteria. |
| T4 | test-rel01.0-demo-trace-collection tracing_enabled | S3 traces this requirement in the existing success criteria. |
