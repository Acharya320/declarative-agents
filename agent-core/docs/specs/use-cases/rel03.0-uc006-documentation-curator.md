---
id: rel03.0-uc006-documentation-curator
title: Knowledge Manager serves documentation curator UX and REST tools
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

Start the external documentation-curator profile to launch a generic REST
server. Profile-owned static assets and machine_request routes serve the React
app and documentation workflows. Atomic filesystem, rg/jq, jurist, compose,
predicate, and file-write words execute policy selected by request-machine.yaml.
No documentation-specific Go host or bench handler participates.

<!-- field: actor -->
## Actor

Developer browsing or curating Agent Core documentation

<!-- field: trigger -->
## Trigger

A developer starts the curator when they need to browse the corpus or review
proposed documentation patches through the Knowledge Manager boundary.

<!-- field: flow -->
## Flow

- F1: The agent loads the external documentation-curator profile and resolves machine, tool selection, builtin declarations, REST declarations, and REST definitions relative to that profile.
- F2: The documentation-curator machine starts in Idle and dispatches launch_curator_http on Seed.
- F3: launch_curator_http selects documentation_curator from rest.yaml, starts the generic server without blocking the machine, and emits ServerLaunched.
- F4: The machine launches the profile-owned lifecycle control server and then dispatches await_curator_control.
- F5: The generic server serves the profile-owned documentation SPA and executes each API route through request-machine.yaml.
- F6: List/read use configured resources; search uses rg and jq; validation uses jurist words; suggestions and guarded decisions persist rollback-capable review artifacts under .documentation-curator.
- F7: A POST to /api/lifecycle/exit enqueues ExitRequested. await_curator_control returns the event. MachineSpec routes it to exit_agent, then stops monitor and curator REST listeners before Done.
- F8: Startup, control launch, await, or lifecycle failure emits CommandError or another declared failure signal and moves the machine to Failed.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd020-bench-web-ui R6.1, R6.2 | bench does not mount documentation API routes |
| T2 | srd027-rest-tool-runtime R5.1, R5.2, R5.3, R5.4, R5.5, R5.6, R5.7 and R10.1, R10.2, R10.3, R10.4, R10.5, R10.6, R10.7, R10.8 | rest_server_launch starts configured listeners and rest_await_event waits for lifecycle control |
| T3 | srd027-rest-tool-runtime R2.1, R2.2, R2.3, R2.4, R2.5, R2.6, R3.1, R3.2, R3.3, R3.4, R3.5, R3.6 | documentation REST tooling is defined by ToolDef and REST YAML |
| T4 | srd028-rest-client-tools R1.1, R1.2, R1.3, R1.4, R1.5, R1.6, R2.1, R2.2, R2.3, R2.4, R2.5, R2.6, R2.7, R6.1, R6.2, R6.3, R6.4, R6.5, R6.6, R6.7, R6.8 | doc_* words call configured REST operations and emit REST result signals |
| T5 | knowledge-manager/documentation-curator/profile.yaml | profile-owned documentation application under AGENT_CATALOG_ROOT |
| T6 | knowledge-manager/documentation-curator/ui/docs | profile-owned React source and built assets |
| T7 | srd029-rest-server-tools R7.1, R7.2, R7.3, R7.4, R7.5, R7.6 and R11.1, R11.2, R11.3, R11.4, R11.5, R11.6, R11.7, R11.8 | lifecycle control endpoints enqueue events for MachineSpec routing and shutdown releases the owned documentation listener |

<!-- field: proof_owner -->
## Proof owner

applications/catalog

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The documentation-curator profile loads generic curator/control/monitor lifecycle words, exit_agent, and doc_* tool selections.

<!-- field: traces -->
#### Traces

- srd022-builtin-tools AC1
- srd027-rest-tool-runtime AC2

<!-- sc: S2 -->
### S2

The configured generic REST server serves profile-owned SPA assets and /api/v1/docs machine requests.

<!-- field: traces -->
#### Traces

- srd024-blocking-input-tools AC1
- srd027-rest-tool-runtime AC2

<!-- sc: S3 -->
### S3

Bench /api/v1/docs routes return 404.

<!-- field: traces -->
#### Traces

- srd020-bench-web-ui AC4

<!-- sc: S4 -->
### S4

Documentation curator REST declarations cover read, search, validation, suggestion, approval, rejection, and reopen words.

<!-- field: traces -->
#### Traces

- srd028-rest-client-tools AC1
- srd028-rest-client-tools AC7

<!-- sc: S5 -->
### S5

The generic lifecycle exit endpoint causes ExitRequested, AgentExited, listener shutdown, and Done through MachineSpec routing.

<!-- field: traces -->
#### Traces

- srd027-rest-tool-runtime AC10
- srd029-rest-server-tools AC10

<!-- field: out_of_scope -->
## Out of scope

- Editing documentation content without explicit patch approval.
- Routing documentation UX through bench/profile.yaml under AGENT_CATALOG_ROOT.
- Replacing the declarative REST runtime with generated clients.

<!-- field: machine_config -->
## Machine config

agents/knowledge-manager/documentation-curator/machine.yaml

<!-- field: rest_definitions -->
## REST definitions

agents/knowledge-manager/documentation-curator/rest.yaml

<!-- field: test_suite_trace -->
## Test suite trace

test-rel03.0-rest-tools
