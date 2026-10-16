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

<!-- field: non_goals -->
## Non-goals

- A Kubernetes or container backend for scenario child starts; children are local processes in this SRD.
- Span-based assertion words; validators assert over HTTP responses and the srd039 request log.
- Scenario scheduling policy beyond one rig machine's sequential scenario loop; cross-scenario parallelism belongs to the caller.
- Scheduling policy across detached children beyond a declared count bound; queueing or prioritising launches belongs to the caller.

<!-- field: tool_words -->
## Tool words

```yaml
- name: init_scenario_session
  input: configured scenario roots.
  output: discovered scenario count and searched roots.
  signals: [SessionSeeded, NoScenarios, CommandError]
  side_effects: [scenario_session_seed]
  errors: [scenario_discovery]
  undo: irreversible; a new initialization replaces the session.
- name: next_scenario
  input: the seeded scenario session.
  output: current subject, scenario, validators, fixtures, or exhaustion.
  signals: [ScenarioReady, AllScenariosDone, CommandError]
  side_effects: [scenario_cursor_advance]
  errors: [scenario_manifest]
  undo: irreversible; a consumed cursor position is not restored.
- name: start_scenario_mock
  input: current scenario and one fixture iterator binding.
  output: mock identity, fixture, environment binding, base URL, and receipt.
  signals: [MockStarted, CommandError]
  side_effects: [child_process_start, scenario_session_update]
  errors: [selector_resolution, address_allocation, spawn_failed]
  undo: stop only the child identified by the dispatch receipt.
- name: start_scenario_subject
  input: current scenario, subject profile, and recorded mock addresses.
  output: subject identity, profile, trusted base URL, health target, start time, and receipt.
  signals: [SubjectStarted, CommandError]
  side_effects: [child_process_start, scenario_session_update]
  errors: [subject_profile, health_target, address_allocation, spawn_failed]
  undo: stop only the child identified by the dispatch receipt.
- name: run_scenario_validator
  input: one validator iterator binding and the current subject URL.
  output: validator terminal state, exit code, timeout state, and diagnostics.
  signals: [ValidatorCompleted, ValidatorIncomplete, CommandError]
  side_effects: [child_agent_execution]
  errors: [selector_resolution, missing_subject, execution_failed]
  undo: noop; a completed validator run cannot be unrun.
- name: record_scenario_validators
  input: the validator iterator's ordered join.
  output: ordered validator outcomes and aggregate pass state.
  signals: [ValidatorsRecorded, CommandError]
  side_effects: [scenario_session_update]
  errors: [selector_resolution, outcome_decode]
  undo: irreversible; recorded outcomes remain in the scenario report.
- name: collect_scenario_verdict
  input: recorded validator outcomes or a declared forced-failure reason.
  output: scenario identity, pass state, reason, and outcomes.
  signals: [ScenarioPassed, ScenarioFailed]
  side_effects: [scenario_session_update]
  errors: [scenario_failed]
  undo: irreversible; the verdict remains in the run report.
- name: list_scenario_children
  input: the current scenario session.
  output: subject-first ordered child identities.
  signals: [ScenarioChildrenListed, CommandError]
  side_effects: [scenario_session_read]
  errors: [missing_current_scenario]
  undo: noop.
- name: stop_service
  input: one literal or selector-resolved tracked child identity.
  output: child identity and stop result.
  signals: [ServiceStopped, CommandError]
  side_effects: [child_process_stop]
  errors: [selector_resolution]
  undo: noop; a stopped child is restarted only by a later scenario start.
- name: stop_all_services
  input: the shared bounded service state.
  output: count and results for all children reaped.
  signals: [AllServicesStopped]
  side_effects: [child_process_stop]
  errors: []
  undo: irreversible; recovery does not restart children.
- name: report_scenario_session
  input: all scenario verdicts retained in the session.
  output: aggregate pass state, verdicts, and first failure.
  signals: [SessionPassed, SessionFailed]
  side_effects: [scenario_session_read]
  errors: [session_failed]
  undo: noop.
- name: start_service
  input: a declared profile, an optional literal or selector-resolved service name, request, and output, a working directory, environment additions, and an optional max_running bound.
  output: child identity, pid, address, base URL, start time, and receipt.
  signals: [ServiceStarted, ServiceLimitReached, CommandError]
  side_effects: [child_process_start]
  errors: [selector_resolution, address_allocation, spawn_failed]
  undo: stop only the child identified by the dispatch receipt.
- name: list_services
  input: the shared service state.
  output: every tracked child with running or exited status, exit code, and start and finish times.
  signals: [ChildrenListed]
  side_effects: [service_state_read]
  errors: []
  undo: noop.
```

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

A machine starts scenario mock and subject children with injected environment, probes the subject's labeled runtime URL through the REST client, stops each child, and no process from either group survives an ordinary teardown, an emergency cleanup edge, or orderly parent shutdown.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R1.4
- R2.1
- R3.1
- R3.2
- R3.4

<!-- ac: AC2 -->
### AC2

A never-healthy child follows machine-visible ProbingHealth and WaitingToRetry states, performs one request per REST dispatch, and returns HealthTimeout within the declared deadline.

<!-- field: traces -->
#### Traces

- R2.2
- R2.3
- R2.4

<!-- ac: AC3 -->
### AC3

MachineSpec dispatches one child execution per validator, joins outcomes in input order, and a hung validator is stopped and reported as timed out within its declared bound.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
- R4.3
- R4.4
- R4.5

<!-- ac: AC4 -->
### AC4

init_scenario_session over a fixture tree seeds the expected entries in sorted order, next_scenario exposes each subject, scenario, validator set, and fixture set exactly once, and repeated runs return identical order.

<!-- field: traces -->
#### Traces

- R5.1
- R5.2
- R5.3
- R5.4

<!-- ac: AC5 -->
### AC5

Each word declares its reversibility and undo consistently, invokes no model, and surfaces a spawn failure as a tool error signal, adding no tool-undo-mismatch audit warnings.

<!-- field: traces -->
#### Traces

- R1.5
- R3.3
- R6.1
- R6.2
- R6.3

<!-- ac: AC6 -->
### AC6

start_service returns before its child exits, list_services shows the child running then exited with its exit code, a declared max_running bound emits ServiceLimitReached, a child started inside a machine_request run is reaped by orderly host shutdown, and both words declare undo consistently.

<!-- field: traces -->
#### Traces

- R7.1
- R7.2
- R7.3
- R7.4
- R7.5

<!-- field: implementation_status -->
## Implementation status

```yaml
implemented:
    - GH-2197 adds start_service and list_services and shares the host service state with machine_request runs (R7); the bound is reserved under the state lock, so concurrent request runs cannot overshoot it.
    - GH-2188 gives the snake_case start_service fields and address_env their json tags, so a declaration decoded by the strict JSON config decoder loads.
planned:
    - GH-2188 is the first consumer, launching bench experiments as detached critic children behind a request capability.
```
