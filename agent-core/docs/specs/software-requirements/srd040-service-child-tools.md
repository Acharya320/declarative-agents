---
id: srd040-service-child-tools
title: Service-Child and Scenario Tools
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

A test-rig machine composes other machines: it starts a subject agent and its
mocks as serve-mode children, runs validator machines against
the subject, and reduces their terminal states to a verdict. The existing
child-agent path (execute.RunAgent, used by the critic) spawns children that
run to completion. Serve-mode children are different: they must start in the
background, report an address, stay up while validators run, and stop
cleanly without leaking processes.

The rig also needs to find its work. Scenarios live as data beside the agent
they test, in tests/<scenario>/ directories under configurable roots, so a
discovery word must enumerate them deterministically.

<!-- field: goals -->
## Goals

- G1: Define start_scenario_mock, start_scenario_subject, and stop_service for background serve-mode scenario children, with readiness composed from ordinary REST and delay words.
- G2: Define one-validator child execution for MachineSpec iteration and join.
- G3: Define init_scenario_session and next_scenario for deterministic scenario discovery and iteration under configured roots.
- G4: Keep every word deterministic, LLM-free, and bounded in time.
- G5: Declare honest reversibility and undo behavior for each word.
- G6: Define start_service and list_services so any profile can start a child detached, bound how many run at once, and read their outcomes.

<!-- req: R1 -->
## R1 - Scenario child start words

<!-- req: R1.1 -->
### R1.1

start_scenario_mock SHALL spawn the declared mock profile for one current-scenario fixture, and start_scenario_subject SHALL spawn the current scenario's subject profile, as background children serving until stopped.

<!-- req: R1.2 -->
### R1.2

start_scenario_mock and start_scenario_subject SHALL inject declared and scenario-derived environment additions into each child, so fixture data and dependency URLs reach the intended scenario processes.

<!-- req: R1.3 -->
### R1.3

Each scenario child start SHALL allocate or accept a listen address, and its output SHALL include the child identity and trusted runtime base URL.

<!-- req: R1.4 -->
### R1.4

Children SHALL use the shared subprocess start transport for process-group management; cancellation of the owning agent context SHALL kill each child's whole process group, and orderly parent shutdown SHALL reap every tracked child rather than orphan it. SIGKILL cannot run process cleanup and is outside this guarantee.

<!-- req: R1.5 -->
### R1.5

start_scenario_mock and start_scenario_subject SHALL declare compensating undo that stops only the child named by that dispatch's receipt.

<!-- req: R2 -->
## R2 - Declarative readiness

<!-- req: R2.1 -->
### R2.1

start_scenario_subject SHALL expose the trusted runtime base URL and start timestamp needed by later profile-owned readiness words.

<!-- req: R2.2 -->
### R2.2

One REST readiness dispatch SHALL perform at most one HTTP request; MachineSpec SHALL own the visible probe, delay, retry, and timeout branches.

<!-- req: R2.3 -->
### R2.3

The readiness machine SHALL emit distinct Healthy, HealthTimeout, and CommandError signals and route timeout, budget exhaustion, and teardown infrastructure errors through a declared cleanup action.

<!-- req: R2.4 -->
### R2.4

Service-child tools SHALL NOT implement HTTP transport or polling and SHALL NOT import net/http.

<!-- req: R3 -->
## R3 - stop_service

<!-- req: R3.1 -->
### R3.1

stop_service SHALL stop a child by handle with a graceful signal, a bounded wait, then a kill.

<!-- req: R3.2 -->
### R3.2

stop_service SHALL be idempotent; stopping an already-stopped handle SHALL succeed.

<!-- req: R3.3 -->
### R3.3

stop_service undo SHALL be noop and the word SHALL be declared accordingly.

<!-- req: R3.4 -->
### R3.4

stop_all_services SHALL be an atomic, idempotent recovery word over the bounded child set owned by the shared service state; it SHALL NOT replace ordinary per-child teardown.

<!-- req: R4 -->
## R4 - Validator child execution

<!-- req: R4.1 -->
### R4.1

One validator dispatch SHALL accept one profile path and its environment additions from the current MachineSpec iterator item.

<!-- req: R4.2 -->
### R4.2

Each dispatch SHALL run one validator to completion through the shared child-agent execution path.

<!-- req: R4.3 -->
### R4.3

Each output SHALL include that validator's terminal state and exit code.

<!-- req: R4.4 -->
### R4.4

Each dispatch SHALL enforce a declared timeout; a validator still running at timeout SHALL be stopped and reported as timed out.

<!-- req: R4.5 -->
### R4.5

MachineSpec bounded parallel for_each SHALL expose every validator dispatch and its deterministic input-order join; a thin record word may retain joined outcomes for reporting.

<!-- req: R5 -->
## R5 - Scenario session initialization and iteration

<!-- req: R5.1 -->
### R5.1

init_scenario_session SHALL enumerate scenario directories matching tests/<scenario>/machine.yaml under declared roots and seed the scenario work list.

<!-- req: R5.2 -->
### R5.2

Each seeded scenario SHALL retain the subject directory, scenario name, validator profile paths, and fixture files present; next_scenario SHALL expose the current scenario's subject, name, validators, and fixtures.

<!-- req: R5.3 -->
### R5.3

The seeded work list SHALL be sorted and deterministic for a given tree, and next_scenario SHALL visit each entry once in that order.

<!-- req: R5.4 -->
### R5.4

Roots SHALL come from declared configuration or run input, not from hardcoded repository paths.

<!-- req: R6 -->
## R6 - Determinism and boundaries

<!-- req: R6.1 -->
### R6.1

The words in this SRD SHALL NOT invoke a model.

<!-- req: R6.2 -->
### R6.2

The words SHALL be declared as boundary tools crossing a child process, per the tool language boundary rules.

<!-- req: R6.3 -->
### R6.3

Failure of a child to spawn SHALL surface as a tool error signal, not a panic.

<!-- req: R7 -->
## R7 - Generic detached service words

<!-- req: R7.1 -->
### R7.1

start_service SHALL spawn one declared profile as a background child through the shared service state and return its handle without waiting for the child to exit.

<!-- req: R7.2 -->
### R7.2

The host run and every machine_request run in the same process SHALL share one tracked child set, so a child started by a request run is reaped under R1.4 and visible to list_services and stop_service from any run.

<!-- req: R7.3 -->
### R7.3

When max_running is declared, start_service SHALL compare it against the live child count before spawning and SHALL emit ServiceLimitReached without starting a child at the bound.

<!-- req: R7.4 -->
### R7.4

list_services SHALL report every tracked child with its running or exited status, exit code, and start and finish times; an exited child SHALL remain listed until stopped or reaped.

<!-- req: R7.5 -->
### R7.5

start_service SHALL declare compensating undo that stops only the child named by its receipt; list_services undo SHALL be noop.

<!-- field: tool_words | form: "typed table" -->
## Tool words

| name | input | output | signals | side_effects | errors | undo |
| --- | --- | --- | --- | --- | --- | --- |
| init_scenario_session | configured scenario roots. | discovered scenario count and searched roots. | [SessionSeeded, NoScenarios, CommandError] | [scenario_session_seed] | [scenario_discovery] | irreversible; a new initialization replaces the session. |
| next_scenario | the seeded scenario session. | current subject, scenario, validators, fixtures, or exhaustion. | [ScenarioReady, AllScenariosDone, CommandError] | [scenario_cursor_advance] | [scenario_manifest] | irreversible; a consumed cursor position is not restored. |
| start_scenario_mock | current scenario and one fixture iterator binding. | mock identity, fixture, environment binding, base URL, and receipt. | [MockStarted, CommandError] | [child_process_start, scenario_session_update] | [selector_resolution, address_allocation, spawn_failed] | stop only the child identified by the dispatch receipt. |
| start_scenario_subject | current scenario, subject profile, and recorded mock addresses. | subject identity, profile, trusted base URL, health target, start time, and receipt. | [SubjectStarted, CommandError] | [child_process_start, scenario_session_update] | [subject_profile, health_target, address_allocation, spawn_failed] | stop only the child identified by the dispatch receipt. |
| run_scenario_validator | one validator iterator binding and the current subject URL. | validator terminal state, exit code, timeout state, and diagnostics. | [ValidatorCompleted, ValidatorIncomplete, CommandError] | [child_agent_execution] | [selector_resolution, missing_subject, execution_failed] | noop; a completed validator run cannot be unrun. |
| record_scenario_validators | the validator iterator's ordered join. | ordered validator outcomes and aggregate pass state. | [ValidatorsRecorded, CommandError] | [scenario_session_update] | [selector_resolution, outcome_decode] | irreversible; recorded outcomes remain in the scenario report. |
| collect_scenario_verdict | recorded validator outcomes or a declared forced-failure reason. | scenario identity, pass state, reason, and outcomes. | [ScenarioPassed, ScenarioFailed] | [scenario_session_update] | [scenario_failed] | irreversible; the verdict remains in the run report. |
| list_scenario_children | the current scenario session. | subject-first ordered child identities. | [ScenarioChildrenListed, CommandError] | [scenario_session_read] | [missing_current_scenario] | noop. |
| stop_service | one literal or selector-resolved tracked child identity. | child identity and stop result. | [ServiceStopped, CommandError] | [child_process_stop] | [selector_resolution] | noop; a stopped child is restarted only by a later scenario start. |
| stop_all_services | the shared bounded service state. | count and results for all children reaped. | [AllServicesStopped] | [child_process_stop] | [] | irreversible; recovery does not restart children. |
| report_scenario_session | all scenario verdicts retained in the session. | aggregate pass state, verdicts, and first failure. | [SessionPassed, SessionFailed] | [scenario_session_read] | [session_failed] | noop. |
| start_service | a declared profile, an optional literal or selector-resolved service name, request, and output, a working directory, environment additions, and an optional max_running bound. | child identity, pid, address, base URL, start time, and receipt. | [ServiceStarted, ServiceLimitReached, CommandError] | [child_process_start] | [selector_resolution, address_allocation, spawn_failed] | stop only the child identified by the dispatch receipt. |
| list_services | the shared service state. | every tracked child with running or exited status, exit code, and start and finish times. | [ChildrenListed] | [service_state_read] | [] | noop. |

<!-- field: non_goals -->
## Non-goals

- A Kubernetes or container backend for scenario child starts; children are local processes in this SRD.
- Span-based assertion words; validators assert over HTTP responses and the srd039 request log.
- Scenario scheduling policy beyond one rig machine's sequential scenario loop; cross-scenario parallelism belongs to the caller.
- Scheduling policy across detached children beyond a declared count bound; queueing or prioritising launches belongs to the caller.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

A machine starts scenario mock and subject children with injected environment, probes the subject's labeled runtime URL through the REST client, stops each child, and no process from either group survives an ordinary teardown, an emergency cleanup edge, or orderly parent shutdown.

<!-- field: traces -->
#### Traces

- [R1.1](#r11)
- [R1.2](#r12)
- [R1.3](#r13)
- [R1.4](#r14)
- [R2.1](#r21)
- [R3.1](#r31)
- [R3.2](#r32)
- [R3.4](#r34)

<!-- ac: AC2 -->
### AC2

A never-healthy child follows machine-visible ProbingHealth and WaitingToRetry states, performs one request per REST dispatch, and returns HealthTimeout within the declared deadline.

<!-- field: traces -->
#### Traces

- [R2.2](#r22)
- [R2.3](#r23)
- [R2.4](#r24)

<!-- ac: AC3 -->
### AC3

MachineSpec dispatches one child execution per validator, joins outcomes in input order, and a hung validator is stopped and reported as timed out within its declared bound.

<!-- field: traces -->
#### Traces

- [R4.1](#r41)
- [R4.2](#r42)
- [R4.3](#r43)
- [R4.4](#r44)
- [R4.5](#r45)

<!-- ac: AC4 -->
### AC4

init_scenario_session over a fixture tree seeds the expected entries in sorted order, next_scenario exposes each subject, scenario, validator set, and fixture set exactly once, and repeated runs return identical order.

<!-- field: traces -->
#### Traces

- [R5.1](#r51)
- [R5.2](#r52)
- [R5.3](#r53)
- [R5.4](#r54)

<!-- ac: AC5 -->
### AC5

Each word declares its reversibility and undo consistently, invokes no model, and surfaces a spawn failure as a tool error signal, adding no tool-undo-mismatch audit warnings.

<!-- field: traces -->
#### Traces

- [R1.5](#r15)
- [R3.3](#r33)
- [R6.1](#r61)
- [R6.2](#r62)
- [R6.3](#r63)

<!-- ac: AC6 -->
### AC6

start_service returns before its child exits, list_services shows the child running then exited with its exit code, a declared max_running bound emits ServiceLimitReached, a child started inside a machine_request run is reaped by orderly host shutdown, and both words declare undo consistently.

<!-- field: traces -->
#### Traces

- [R7.1](#r71)
- [R7.2](#r72)
- [R7.3](#r73)
- [R7.4](#r74)
- [R7.5](#r75)

<!-- field: implementation_status -->
## Implementation status

<!-- field: implemented -->
### Implemented

- GH-2197 adds start_service and list_services and shares the host service state with machine_request runs (R7); the bound is reserved under the state lock, so concurrent request runs cannot overshoot it.
- GH-2188 gives the snake_case start_service fields and address_env their json tags, so a declaration decoded by the strict JSON config decoder loads.

<!-- field: planned -->
### Planned

- GH-2188 is the first consumer, launching bench experiments as detached critic children behind a request capability.
