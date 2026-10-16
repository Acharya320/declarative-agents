---
id: rel01.0-uc003-bench-visualization
title: Bench serves web UI for evaluation result exploration
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

The agent, configured with the external bench profile, starts a
state-machine-driven web server for browsing experiment results, inspecting
source code, and launching new evaluations. A generic REST event wait blocks
on operator input, and self_invoke spawns the configured critic profile as a
child agent process. The critic profile owns the session machine it
runs.

<!-- field: actor -->
## Actor

Developer analyzing benchmark results

<!-- field: trigger -->
## Trigger

The developer runs:
  agent --profile "$AGENT_CATALOG_ROOT/agents/bench/profile.yaml"

The profile owns data roots, routes, static assets, listener configuration,
request MachineSpec, and generic tool aliases.

<!-- field: flow -->
## Flow

- F1: The agent parses CLI flags and loads the external bench profile, resolving the bench machine, tool selection, and tool declarations relative to the profile file (cfg-runtime-contract).
- F2: The bench state machine is loaded from the profile-resolved bench/machine.yaml path (srd006).
- F3: The profile loads rest.yaml, request-machine.yaml, resource roots, and static assets (srd020).
- F4: Generic REST, filesystem, evaluation-query, predicate, compose, exec, and self_invoke factories are selected.
- F5: core.Loop launches the generic REST server and waits on its bounded action queue (srd005, bench/machine.yaml).
- F6: The React SPA loads in the browser. The user browses sessions, inspects traces, opens launch controls, and reads source code via REST APIs (srd020).
- F7: GET /api/v1/sessions returns experiment sessions with pass/fail/timeout counts. The user drills into session details, point traces, and experiment configs (srd020).
- F8: The user opens the Launch page and submits an experiment configuration. POST /api/v1/actions sends an ExperimentRequested action (srd020).
- F9: await_bench_action returns ExperimentRequested. The state machine validates suite input and self_invoke starts the configured critic profile (srd021, bench/machine.yaml, bench/builtin.yaml).
- F10: ToolDone or ToolFailed transitions Launching → Serving and awaits the next operator action.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd005-agentic-loop R4.1, R4.2, R4.3, R4.4, R5.1, R5.2, R5.3, R5.4, R5.5, R5.6, R5.7, R5.8, R5.9, R5.10, R5.11, R5.12 | LoopParams and Loop drive the bench state machine |
| T2 | srd006-yaml-machine-config R1.1, R1.2, R1.3, R1.4, R1.5, R1.6, R1.7, R1.8, R1.9, R2.1, R2.2, R2.3, R2.4, R2.5, R2.6, R2.7, R2.8, R2.9, R2.10, R2.11, R3.1, R3.2, R3.3, R3.4, R3.5, R4.1, R4.2, R4.3, R4.4 | bench machine.yaml loading, validation, and table building |
| T3 | srd020-bench-web-ui R1.1, R1.2, R1.3, R1.4, R2.1, R2.2, R2.3, R2.4, R2.5, R2.6, R3.1, R3.2, R3.3, R4.1, R4.2, R4.3, R4.4, R4.5, R4.6, R5.1, R5.2, R6.1, R6.2, R7.1, R7.2, R7.3, R8.1, R8.2, R8.3 | generic HTTP server, bench REST APIs, SPA routing, profile assets, and documentation boundary |
| T4 | bench/machine.yaml | generic REST lifecycle and self_invoke sequencing |
| T5 | srd019-eval-harness R0.1, R0.2, R0.3, R0.4, R0.5, R0.6, R1.1, R1.2, R1.3, R1.4, R1.5, R1.6, R1.7, R2.1, R2.2, R2.3, R2.4, R2.5, R3.1, R3.1.1, R3.1.2, R3.1.3, R3.1.4, R3.1.5, R3.1.6, R3.2, R3.3, R3.4, R3.5, R3.6, R3.7, R3.8, R3.9 | self_invoke reuses the profile-driven evaluator pipeline |
| T6 | bench/profile.yaml | profile-relative bench runtime entry point under AGENT_CATALOG_ROOT |
| T7 | srd021-agent-cli-exec-tool R1.1, R1.2, R1.3, R1.4, R1.5, R1.6, R1.7, R2.1, R2.2, R2.3, R2.4, R2.5, R3.1, R3.2, R3.3, R3.4, R4.1, R4.2, R4.3, R4.4 | child agent profile invocation, trace propagation, process result mapping |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The web UI serves from profile-owned assets and all SPA routes resolve without 404.

<!-- field: traces -->
#### Traces

- srd020-bench-web-ui AC1

<!-- sc: S2 -->
### S2

Session, point, trace, and experiment APIs return correct data from the results directory.

<!-- field: traces -->
#### Traces

- srd020-bench-web-ui AC2

<!-- sc: S3 -->
### S3

Source code and config APIs serve content from disk, and bench docs routes are absent.

<!-- field: traces -->
#### Traces

- srd020-bench-web-ui AC4

<!-- sc: S4 -->
### S4

User actions from the web UI drive state machine transitions via the configured REST event queue.

<!-- field: traces -->
#### Traces

- srd020-bench-web-ui AC4

<!-- sc: S5 -->
### S5

Launching an experiment spawns the evaluator profile child process and returns to serving on completion.

<!-- field: traces -->
#### Traces

- srd020-bench-web-ui AC4

<!-- field: out_of_scope -->
## Out of scope

- User authentication or multi-tenancy.
- Editing configurations through the UI.
- Real-time streaming of in-progress evaluations.

<!-- field: machine_config -->
## Machine config

bench/machine.yaml

<!-- field: test_suite_trace -->
## Test suite trace

test-rel01.0
