---
id: test-rel02.0-serving-contract
title: Persistent Coding-Agent Serving Contract Suite
release: "02.0"
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

The always-runnable stage starts all three persistent profiles and proves real
lifecycle health plus a deterministic canonical critic request. The remote
stage uses the same real binary and profiles for planner to executor to critic,
with deterministic Ollama-compatible responses at the model boundary.

<!-- field: traces -->
## Traces

- rel02.0-uc001-remote-coding-request

<!-- field: preconditions -->
## Preconditions

- agent-core and applications/catalog checkouts are available.
- The real agent binary can be built.
- The remote stage starts a deterministic Ollama-compatible boundary; no external model is required.

<!-- field: test_cases -->
## Test cases

<!-- case: Persistent health and deterministic critic request | status: implemented -->
### Persistent health and deterministic critic request

<!-- field: use_case -->
#### Use case

rel02.0-uc001-remote-coding-request

<!-- field: integration_target -->
#### Integration target

mage integration:servingHealth

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^Test(ServingProfilesUseRealLifecycleAndCanonicalRoleProfiles|StageServingProfileTreeCombinesClosureAndApplicationProfiles)$'

<!-- field: expected_outputs | form: "tree" -->
#### Expected outputs

- planner_health: HTTP 200
- executor_health: HTTP 200
- critic_health: HTTP 200
- critic_verdict: accepted

<!-- field: traces -->
#### Traces

- srd001-serving-contract AC1

<!-- case: Connected remote coding request | status: implemented -->
### Connected remote coding request

<!-- field: use_case -->
#### Use case

rel02.0-uc001-remote-coding-request

<!-- field: integration_target -->
#### Integration target

mage integration:servingRemote

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestPlannerServingFlowUsesDeclaredRemoteBoundaries$'

<!-- field: expected_outputs | form: "tree" -->
#### Expected outputs

- planner_response: accepted
- workspace_validation: go test ./... passes
- remote_order: planner then executor then critic
- trace_connection: caller trace id appears in all three role trace files

<!-- field: traces -->
#### Traces

- srd001-serving-contract AC2
- srd001-serving-contract AC3

<!-- case: Local batch coding loop remains valid | status: implemented -->
### Local batch coding loop remains valid

<!-- field: use_case -->
#### Use case

rel02.0-uc001-remote-coding-request

<!-- field: integration_target -->
#### Integration target

mage integration:codingLoop

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestCodingLoopStagesGenerateOnceAndPreserveExactWorkspace$'

<!-- field: expected_outputs | form: "tree" -->
#### Expected outputs

- loop_outcome: the batch planner, executor, and critic stages pass unchanged alongside the serving profiles

<!-- field: traces -->
#### Traces

- srd001-serving-contract AC4
