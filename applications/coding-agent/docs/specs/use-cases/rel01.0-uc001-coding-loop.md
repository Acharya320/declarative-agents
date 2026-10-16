---
id: rel01.0-uc001-coding-loop
title: Planner Materializes, Executor Changes, and Critic Gates
test_suite: test-rel01.0-coding-loop
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An operator supplies an SRD, requirement graph, and isolated workspace to the
coding application. The canonical planner profile materializes one task and
delegates it to the canonical executor profile. The executor changes the
workspace and runs declared validation. The canonical critic profile evaluates
the produced change with a suite and oracle. Only a successful executor result
and an accepting critic report produce the application's successful terminal
outcome.

<!-- field: actor -->
## Actor

A developer or delivery system requesting a validated implementation of an SRD requirement.

<!-- field: trigger -->
## Trigger

The actor submits a coding request containing the SRD and requirement graph,
an isolated workspace, model configuration, and critic evaluation inputs.

<!-- field: preconditions -->
## Preconditions

- The application package records the configured applications/catalog checkout revision and compatible release and contains the resolved planner, executor, and critic profile closure.
- The workspace contains the baseline source and its validation command.
- The planner and executor model endpoints are reachable.
- The canonical changed-workspace critic profile and its declared oracle are available.
- Planner, executor, and critic run as separate agent-core processes with their selected profiles mounted from the application package.

<!-- field: flow -->
## Flow

- F1: The planner receives the requirement graph, SRD context, planning model config, and workspace reference.
- F2: The planner traverses the requirement graph, assembles and parses a plan, and materializes a task that names the requested workspace change and validation.
- F3: The planner delegates the materialized task, workspace, and execution model config to the real executor profile through child-profile execution.
- F4: The executor invokes its model boundary, dispatches filesystem tools to change the workspace, runs its declared validation tools, and returns terminal status, tool results, validation results, and a task summary.
- F5: The planner records the executor outcome and makes the changed workspace and evidence available to the application evaluation stage.
- F6: The application supplies the produced workspace directly to the canonical changed-workspace critic, which runs its declared oracle in place without launching another coding agent.
- F7: The critic selects its accepted or rejected machine branch from the oracle signal and writes critic-verdict.json with the explicit verdict and oracle evidence.
- F8: The application succeeds only when executor validation succeeds and the critic gate returns an accepting Done report; executor failure, budget exhaustion, or critic rejection produces a failed outcome with evidence.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | applications/catalog srd004-planner R1, R2, R3, AC1, AC2, AC3 | The planner traverses requirements, materializes the task, delegates through child-profile execution, and records the outcome. |
| T2 | applications/catalog srd002-executor R1, R2, R3, AC1, AC2, AC3 | The executor receives the task and workspace, edits through declared tools, validates, and returns its explicit outcome. |
| T3 | applications/catalog srd003-critic R1, R2, R3, AC1, AC2, AC3 | The critic evaluates the produced change through suite points and oracles and returns the gating session report. |
| T4 | architecture-coding-agent-example components | The application architecture owns the planner to executor to critic sequence and one-container-per-agent topology. |
| T5 | architecture-coding-agent-example design_decisions 1, 2, 3 | The flow references canonical library agents, uses package-time profile closure, and keeps the runtime image profile-free. |
| T6 | test-rel01.0-coding-loop | The suite proves the loop in live-executor, real-planner-delegation, and critic-gate stages. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The real executor changes the isolated workspace for the supplied task, its declared validation command succeeds, and it returns a traceable successful result.

<!-- field: traces -->
#### Traces

- applications/catalog srd002-executor R1.1
- applications/catalog srd002-executor R3.1
- applications/catalog srd002-executor R3.2
- applications/catalog srd002-executor AC1
- applications/catalog srd002-executor AC3
- srd001-serving-contract AC4

<!-- sc: S2 -->
### S2

The real planner materializes the task from SRD context and delegates it to the real executor profile without a shell-script or fake agent binary.

<!-- field: traces -->
#### Traces

- applications/catalog srd004-planner R1.1
- applications/catalog srd004-planner R1.2
- applications/catalog srd004-planner R2.1
- applications/catalog srd004-planner R3.1
- applications/catalog srd004-planner AC1
- applications/catalog srd004-planner AC3
- srd001-serving-contract AC4

<!-- sc: S3 -->
### S3

The canonical critic evaluates the produced change in place and emits a machine-readable accepted or rejected verdict whose value gates the application outcome.

<!-- field: traces -->
#### Traces

- applications/catalog srd003-critic R1.1
- applications/catalog srd003-critic R1.2
- applications/catalog srd003-critic R3.1
- applications/catalog srd003-critic R3.2
- applications/catalog srd003-critic AC1
- applications/catalog srd003-critic AC3
- srd001-serving-contract AC4

<!-- sc: S4 -->
### S4

All three stages use profiles from one checkout-derived package-time closure while every agent process runs the same canonical profile-free agent-core image; the executor mounts its pinned Go and golangci-lint v2 donor volume, and the critic mounts a separate Go-only donor volume for its independent test oracle.

<!-- field: traces -->
#### Traces

- architecture-coding-agent-example design_decisions 1
- architecture-coding-agent-example design_decisions 2
- architecture-coding-agent-example design_decisions 3
- srd002-deployment-profile-package AC1

<!-- field: out_of_scope -->
## Out of scope

- Defining or changing planner, executor, or critic family behavior.
- Copying library profiles or SRDs into the application.
- Implementing the profile reference resolver, package manifest, or Helm chart.
- Proving multi-task scheduling, retries across failed critic sessions, or human approval.
- Replacing library conformance tests with this application integration suite.
