---
id: test-rel01.0-coding-loop
title: Test Suite for Release 01.0 Coding Loop
release: "01.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

This suite defines three ordered integration stages for the coding
application. Stage A proves the real executor can complete a bounded coding
task. Stage B proves the real planner materializes that task and delegates to
the real built agent binary. Stage C proves the critic reviews the resulting
change and gates the terminal application state. The composite retains Stage
B's workspace and passes that exact path to Stage C; it generates no second
planner candidate. All stages use canonical library profiles and the real
agent binary against a deterministic Ollama-compatible model boundary; none
substitutes a shell-script child or fake agent binary.

<!-- field: traces -->
## Traces

- rel01.0-uc001-coding-loop

<!-- field: deadline_hierarchy -->
## Deadline hierarchy

<!-- field: model_call -->
### Model call

30s

<!-- field: child_agent_run -->
### Child agent run

10m

<!-- field: outer_emergency -->
### Outer emergency

12m

<!-- field: cleanup_grace -->
### Cleanup grace

3s

<!-- field: invariant -->
### Invariant

model_call < child_agent_run and child_agent_run + cleanup_grace < outer_emergency

<!-- field: timeout_evidence -->
### Timeout evidence

- combined stdout and stderr
- bounded trace tail
- last observed state
- last observed tool
- classified active phase

<!-- field: cancellation -->
### Cancellation

Outer emergency cancellation sends TERM and then KILL to the full process group after cleanup_grace.

<!-- field: preconditions -->
## Preconditions

- agent-core is available through the demo.yaml core_root or the repository's agent-core directory and can build the real agent binary.
- The release harness can bind a loopback HTTP server implementing the Ollama tags and chat protocol.
- The application package snapshots planner, executor, and critic from one compatible applications/catalog checkout and records its revision.
- Independently addressed stages use fresh greet workspaces; the composite keeps Stage B's workspace alive through Stage C.
- The fixture carries its SRD, requirement graph, fallback candidate variants, and expected critic/application terminal outcomes.

<!-- field: test_cases -->
## Test cases

<!-- case: Stage A - live executor run | status: implemented -->
### Stage A - live executor run

Build the production agent-core binary, start the canonical executor
profile against the deterministic Ollama-compatible boundary, and give it
the greet task in a fresh isolated workspace. The executor must edit the
actual workspace, run its declared Go validation, and report a successful
terminal outcome with evidence. This stage does not invoke the planner or
critic.

<!-- field: use_case -->
#### Use case

rel01.0-uc001-coding-loop

<!-- field: stage -->
#### Stage

A

<!-- field: integration_target -->
#### Integration target

mage integration:executorLive

<!-- field: inputs | form: "tree" -->
#### Inputs

- agent_binary: production agent-core binary built from the demo.yaml core_root
- executor_profile: applications/catalog/agents/executor/profile.yaml resolved through the application package
- task_prompt: Implement the greeting behavior specified by srd001-greet.yaml.
- workspace: fresh testdata/integration/coding-loop/workspace copy
- model_config: canonical model identity served by the deterministic Ollama-compatible release boundary
- model_call_deadline: 30s
- run_deadline: 12m outer emergency bound
- validation_command: go test ./...

<!-- field: expected_outputs | form: "tree" -->
#### Expected outputs

- workspace_change: The greet implementation differs from the baseline and satisfies the SRD.
- validation_result: go test ./... exits successfully in the changed workspace.
- executor_status: Succeeded
- tool_results: Filesystem edits and Go validation results are present.
- task_summary: Names the implemented greeting behavior and successful validation.
- deterministic_model_sequence: edit, then done
- fake_agent_invocations: 0

<!-- field: traces -->
#### Traces

- rel01.0-uc001-coding-loop S1
- srd002-executor AC1
- srd002-executor AC2
- srd002-executor AC3

<!-- case: Stage B - real planner delegation | status: implemented -->
### Stage B - real planner delegation

Start the canonical planner profile with the fixture requirement graph and
SRD context. The planner must materialize the greet task and delegate it to
the canonical executor using the real built agent binary. The delegated
executor changes the workspace and returns successful validation; no
`writeGeneratorChildAgent`, shell-script child, or other fake binary is
permitted. Planner generation occurs exactly once.

<!-- field: use_case -->
#### Use case

rel01.0-uc001-coding-loop

<!-- field: stage -->
#### Stage

B

<!-- field: integration_target -->
#### Integration target

mage integration:plannerDelegation

<!-- field: inputs | form: "tree" -->
#### Inputs

- agent_binary: production agent-core binary used for both parent and child profile execution
- planner_profile: applications/catalog/agents/planner/profile.yaml resolved through the application package
- executor_profile: applications/catalog/agents/executor/profile.yaml resolved through the same package
- requirement_graph: testdata/integration/coding-loop requirement graph with one ready greet requirement
- srd_context: testdata/integration/coding-loop/srd001-greet.yaml
- workspace: fresh testdata/integration/coding-loop/workspace copy
- planning_model_config: canonical planner model identity served by the deterministic Ollama-compatible release boundary
- execution_model_config: canonical executor model identity served by the same deterministic boundary and inherited by the real child process
- model_call_deadline: 30s
- child_agent_deadline: 10m
- outer_emergency_deadline: 12m with 3s process-group cleanup grace

<!-- field: expected_outputs | form: "tree" -->
#### Expected outputs

- materialized_task: One issue or task names the greet requirement, workspace, and go test ./... validation.
- delegation: The planner invokes the canonical executor profile through the real agent binary.
- workspace_change: The delegated executor produces the greet implementation.
- validation_result: go test ./... exits successfully in the changed workspace.
- planner_outcome: Records created task, graph status change, executor outcome, and Completed terminal state.
- planner_model_invocations: 1
- fake_agent_invocations: 0

<!-- field: traces -->
#### Traces

- rel01.0-uc001-coding-loop S1
- rel01.0-uc001-coding-loop S2
- srd004-planner AC1
- srd004-planner AC2
- srd004-planner AC3
- srd002-executor AC3

<!-- case: Stage C - critic gate | status: implemented -->
### Stage C - critic gate

In the composite target, pass the exact retained workspace path produced
by Stage B directly to the canonical changed-workspace critic profile. The
independently addressable criticGate target uses the deterministic
conforming fixture and never invokes the planner. In either form, the
critic runs its declared oracle without launching a coding child, writes
an explicit machine-readable accepted verdict, and the application maps
it to Succeeded. A deliberately non-conforming candidate drives the
rejected branch and maps to Failed.

<!-- field: use_case -->
#### Use case

rel01.0-uc001-coding-loop

<!-- field: stage -->
#### Stage

C

<!-- field: integration_target -->
#### Integration target

mage integration:criticGate

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^Test(ApplicationOutcomeMapsOnlyCanonicalVerdicts|CriticCandidateFixturesHaveOppositeOracleResults|CodingLoopStagesGenerateOnceAndPreserveExactWorkspace|RunBuiltAgentTimeoutRetainsPhaseDiagnosticsAndKillsGroup)$'

<!-- field: inputs | form: "tree" -->
#### Inputs

- agent_binary: production agent-core binary used for critic point execution
- critic_profile: applications/catalog/agents/critic/profile-workspace.yaml resolved through the application package
- composite_accepted_candidate: exact workspace object produced by Stage B, retained until both critic verdicts complete
- standalone_accepted_candidate: deterministic conforming fixture
- rejected_candidate: fixture variant whose greeting violates srd001-greet.yaml
- oracle_command: canonical profile declaration runs go test ./...

<!-- field: expected_outputs | form: "tree" -->
#### Expected outputs

- oracle_results: The conforming candidate passes and the non-conforming candidate fails.
- accepted_verdict: critic-verdict.json records verdict accepted, accepted true, and oracle status passed.
- rejected_verdict: critic-verdict.json records verdict rejected, accepted false, and oracle status failed.
- accepted_session: The conforming candidate reaches Succeeded and maps to application Succeeded.
- rejected_session: The non-conforming candidate reaches Rejected and maps to application Failed.
- composite_workspace_handoff: Stage C receives the identical path returned by Stage B and cleanup occurs only after Stage C.
- standalone_planner_invocations: 0
- composite_planner_invocations: 1
- fake_agent_invocations: 0

<!-- field: traces -->
#### Traces

- rel01.0-uc001-coding-loop S3
- rel01.0-uc001-coding-loop S4
- srd003-critic AC1
- srd003-critic AC2
- srd003-critic AC3
