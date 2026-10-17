---
id: rel25.0-uc001-drain-outcome
title: A monitor exit on a loaded host succeeds and names its unfinished drain
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An operator sends the control exit request to a running monitor while the
host is busy enough that an in-flight connection outlives the shutdown
budget. The listener closes at once, so the profile has done what it was
asked to do, and the run reaches its success terminal. The stop output
carries a structured drain field saying the budget expired, so a machine
branching on the report sees the condition that a passing exit status
would otherwise hide. A stop that cannot release the listener, or that
fails for any other reason, still fails the command.

<!-- field: actor -->
## Actor

An operator stopping a serving profile, and the machine that branches on what the stop reported.

<!-- field: trigger -->
## Trigger

A control exit request arrives at a monitor whose in-flight connections
cannot all finish inside the configured shutdown budget.

<!-- field: flow -->
## Flow

- F1: The control route fires the stop step, which closes the owned listener immediately and waits out the shutdown budget for open connections (srd033 R6.4).
- F2: The wait expires with connections still open; the server is already deregistered and no longer accepting requests (srd033 R6.5).
- F3: The stop step reports the drain outcome as a structured field of its declared output rather than as prose alone (srd033 R6.10).
- F4: The MachineSpec routes the stop to Done or Succeeded, and the run status is succeeded (srd033 R6.5).
- F5: A stop naming an unlaunched server, or one whose shutdown fails for a reason other than the expired budget, stays a command failure and routes to the failure terminal (srd033 R6.5).
- F6: A drain that finishes inside the budget reports the completed outcome through the same field, so the two cases differ in the value and not in the shape (srd033 R6.10).

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [R6.4](../software-requirements/srd033-monitor-rest-api.md#r64) | the stop step releases the owned monitor listener |
| T2 | [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [R6.5](../software-requirements/srd033-monitor-rest-api.md#r65) | a released listener routes to the success terminal, and other faults stay failures |
| T3 | [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [R6.10](../software-requirements/srd033-monitor-rest-api.md#r610) | the drain outcome is a structured field a machine can branch on |
| T4 | [srd038-command-state-store](../software-requirements/srd038-command-state-store.md) [R3.3](../software-requirements/srd038-command-state-store.md#r33) | the reported field travels as declared output, never as a receipt |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

A stop whose drain exceeds its budget returns success, and the run reaches the success terminal instead of the failure terminal.

<!-- field: traces -->
#### Traces

- [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [AC6](../software-requirements/srd033-monitor-rest-api.md#ac6)

<!-- sc: S2 -->
### S2

The stop output carries a structured drain field whose value separates a completed drain from one the budget cut short, in both cases and with the same shape.

<!-- field: traces -->
#### Traces

- [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [AC6](../software-requirements/srd033-monitor-rest-api.md#ac6)

<!-- sc: S3 -->
### S3

An unlaunched server and a shutdown failing for any reason other than the expired budget remain command failures.

<!-- field: traces -->
#### Traces

- [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [AC6](../software-requirements/srd033-monitor-rest-api.md#ac6)

<!-- sc: S4 -->
### S4

The runtime-state-reader proof keeps its shipped shutdown budget and reports succeeded whether or not the drain finished.

<!-- field: traces -->
#### Traces

- [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [AC6](../software-requirements/srd033-monitor-rest-api.md#ac6)

<!-- field: out_of_scope -->
## Out of scope

- Changing any shipped shutdown budget, which stays a profile decision.
- Reporting per-connection detail; the field names the outcome, not the connections.
- Any monitor read endpoint, which R5 and R7 already govern.
