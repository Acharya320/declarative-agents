---
id: rel07.0-uc005-monitor-control-profile-boundary
title: Runtime State Reader and Control Profile Boundary
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

This historically named monitor/control proof records profile-owned
load-and-wiring evidence for declared boundary behavior. Runtime-state-reader
config exposes cached runtime state through monitor REST routes, accepts a
control exit request, and declares cleanup routing through MachineSpec.
Control config accepts a lifecycle exit request and declares visible routing
from ExitRequested to exit_agent and AgentExited.

Runtime listener cleanup is exercised by runtime-state-reader conformance.
This proof stays configuration-level because adapter and control behavior is
encoded in profile-owned YAML assets while agent-core owns the REST and
lifecycle builtin implementations.

<!-- field: actor -->
## Actor

Profile maintainer verifying runtime-state-reader/control boundary behavior.

<!-- field: trigger -->
## Trigger

A maintainer runs the monitor-control integration target:

  mage integration:monitorControl

<!-- field: preconditions -->
## Preconditions

- agents/runtime-state-reader/profile.yaml and testdata/conformance/control/profile.yaml validate.
- testdata/integration/rel07-monitor-control/expected/evidence.yaml defines expected boundary evidence.
- agent-core remains responsible for rest_server_launch, rest_await_event, rest_server_stop, and exit_agent implementations.

<!-- field: flow -->
## Flow

- F1: The target loads runtime-state-reader/control profile assets from applications/catalog.
- F2: The monitor REST config declares current state, metrics, tools, machine, and recent-event read routes.
- F3: Monitor and control configs accept exit requests through emit_signal REST bindings.
- F4: The runtime-state-reader MachineSpec routes ExitRequested to stop_monitor_rest, while control routes ExitRequested to exit_agent.
- F5: The tracer records the declared monitor stop transition, AgentExited routing, and target ownership evidence.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | sm-agent-role-taxonomy runtime-state-reader | Runtime-state-reader adapts cached monitor state for supervisor reads. |
| T2 | sm-agent-role-taxonomy control | Control owns REST event to lifecycle signal routing. |
| T3 | sm-multi-agent-profile-system runtime-state-reader/control | Runtime-state-reader reads cached state and control routes explicit lifecycle signals. |
| T4 | test-rel07.0-profile-integrations monitor_control_boundary | The suite records whether the tracer is implemented or blocked. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The tracer records implemented runtime-state-reader/control load-and-wiring evidence.

<!-- field: traces -->
#### Traces

- test-rel07.0-profile-integrations monitor_control_boundary

<!-- sc: S2 -->
### S2

Runtime-state-reader monitor routes expose cached runtime state without running tools from HTTP handlers.

<!-- field: traces -->
#### Traces

- sm-agent-role-taxonomy runtime-state-reader checks

<!-- sc: S3 -->
### S3

Control routes emit declared lifecycle signals and reach terminal status.

<!-- field: traces -->
#### Traces

- sm-agent-role-taxonomy control checks

<!-- field: out_of_scope -->
## Out of scope

- Building a new monitor UI.
- Adding agent-core monitor features.
- Full HTTP listener exercise for monitor/control; agent-core owns that runtime behavior.

<!-- field: test_suite_trace -->
## Test suite trace

test-rel07.0-profile-integrations
