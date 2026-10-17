---
id: rel99.0-uc001-tool-dispatch-simplification
title: Maintainer simplifies tool dispatch wiring without changing core dispatch
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A maintainer reviews tool dispatch wiring after the focused internal/tools
package split. They keep MachineSpec sequencing, core.Registry lookup,
core.Builder, core.Command, and Dispatch intact while reducing startup and
availability-resolution duplication around the dispatch model.

<!-- field: actor -->
## Actor

Maintainer of agent-core runtime and tool infrastructure

<!-- field: trigger -->
## Trigger

Tool registration or dynamic $tool review finds duplicated availability
checks, panic-based startup failures for profile-owned configuration, or
parallel builtin registration paths that should share one helper.

<!-- field: flow -->
## Flow

- F1: The maintainer confirms that core dispatch remains a small state-machine-to-command execution model (srd037 R1).
- F2: The maintainer changes configuration-originated registry registration failures to returned errors while keeping compatibility wrappers where callers still expect panic behavior (srd037 R2).
- F3: The maintainer makes single and batch builtin registration share one builtin helper, with exec declaration handling still separate (srd037 R3).
- F4: The maintainer introduces one state-scoped external-tool availability resolver and uses it from manifest, parse_response, and dynamic $tool dispatch paths (srd037 R4).
- F5: The maintainer adds regression tests for registry errors and shared parse-time/dispatch-time availability behavior without changing MachineSpec execution or response envelope syntax (srd037 R5).

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd037-tool-dispatch-simplification](../software-requirements/srd037-tool-dispatch-simplification.md) [R1.1](../software-requirements/srd037-tool-dispatch-simplification.md#r11), [R1.2](../software-requirements/srd037-tool-dispatch-simplification.md#r12), [R1.3](../software-requirements/srd037-tool-dispatch-simplification.md#r13), [R1.4](../software-requirements/srd037-tool-dispatch-simplification.md#r14) | MachineSpec, Registry, Builder, Command, and Dispatch semantics stay intact |
| T2 | [srd037-tool-dispatch-simplification](../software-requirements/srd037-tool-dispatch-simplification.md) [R2.1](../software-requirements/srd037-tool-dispatch-simplification.md#r21), [R2.2](../software-requirements/srd037-tool-dispatch-simplification.md#r22), [R2.3](../software-requirements/srd037-tool-dispatch-simplification.md#r23), [R2.4](../software-requirements/srd037-tool-dispatch-simplification.md#r24) | registry registration failures become returned startup errors for configuration-originated mistakes |
| T3 | [srd037-tool-dispatch-simplification](../software-requirements/srd037-tool-dispatch-simplification.md) [R3.1](../software-requirements/srd037-tool-dispatch-simplification.md#r31), [R3.2](../software-requirements/srd037-tool-dispatch-simplification.md#r32), [R3.3](../software-requirements/srd037-tool-dispatch-simplification.md#r33), [R3.4](../software-requirements/srd037-tool-dispatch-simplification.md#r34) | builtin registration shares one helper while exec registration remains separate |
| T4 | [srd037-tool-dispatch-simplification](../software-requirements/srd037-tool-dispatch-simplification.md) [R4.1](../software-requirements/srd037-tool-dispatch-simplification.md#r41), [R4.2](../software-requirements/srd037-tool-dispatch-simplification.md#r42), [R4.3](../software-requirements/srd037-tool-dispatch-simplification.md#r43), [R4.4](../software-requirements/srd037-tool-dispatch-simplification.md#r44), [R4.5](../software-requirements/srd037-tool-dispatch-simplification.md#r45), [R4.6](../software-requirements/srd037-tool-dispatch-simplification.md#r46), [R4.7](../software-requirements/srd037-tool-dispatch-simplification.md#r47) | manifest, parse_response, and dynamic $tool dispatch use one availability rule |
| T5 | [srd037-tool-dispatch-simplification](../software-requirements/srd037-tool-dispatch-simplification.md) [R5.1](../software-requirements/srd037-tool-dispatch-simplification.md#r51), [R5.2](../software-requirements/srd037-tool-dispatch-simplification.md#r52), [R5.3](../software-requirements/srd037-tool-dispatch-simplification.md#r53) | code-task tests prove the narrowed behavior change |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The design names the unchanged core dispatch model and the allowed extension points.

<!-- field: traces -->
#### Traces

- [srd037-tool-dispatch-simplification](../software-requirements/srd037-tool-dispatch-simplification.md) [AC1](../software-requirements/srd037-tool-dispatch-simplification.md#ac1)

<!-- sc: S2 -->
### S2

Configuration-originated startup failures are specified as returned errors, not panics.

<!-- field: traces -->
#### Traces

- [srd037-tool-dispatch-simplification](../software-requirements/srd037-tool-dispatch-simplification.md) [AC2](../software-requirements/srd037-tool-dispatch-simplification.md#ac2)

<!-- sc: S3 -->
### S3

Single and batch builtin registration converge without folding exec declarations into builtin handling.

<!-- field: traces -->
#### Traces

- [srd037-tool-dispatch-simplification](../software-requirements/srd037-tool-dispatch-simplification.md) [AC3](../software-requirements/srd037-tool-dispatch-simplification.md#ac3)

<!-- sc: S4 -->
### S4

Manifest, parse_response, and dynamic $tool dispatch share one state-scoped external-tool availability rule.

<!-- field: traces -->
#### Traces

- [srd037-tool-dispatch-simplification](../software-requirements/srd037-tool-dispatch-simplification.md) [AC4](../software-requirements/srd037-tool-dispatch-simplification.md#ac4)

<!-- sc: S5 -->
### S5

The follow-up code task has concrete tests and excludes MachineSpec, command execution, receipt, and parser-syntax changes.

<!-- field: traces -->
#### Traces

- [srd037-tool-dispatch-simplification](../software-requirements/srd037-tool-dispatch-simplification.md) [AC5](../software-requirements/srd037-tool-dispatch-simplification.md#ac5)

<!-- field: failure_scenarios -->
## Failure scenarios

| id | scenario | expected |
| --- | --- | --- |
| E1 | A refactor changes MachineSpec transition semantics while touching availability checks. | Review rejects the change as outside the SRD scope. |
| E2 | parse_response and dynamic $tool dispatch disagree about whether a tool is available in a state. | The sibling code-task regression test fails. |
| E3 | A builtin registration failure from profile YAML still panics at startup. | Registry or registration tests fail until the error is returned. |

<!-- field: out_of_scope -->
## Out of scope

- Changing LLM tool-call JSON syntax.
- Replacing core.Command, core.Builder, core.Result, or Dispatch.
- Moving workflow order from machine.yaml into Go code.
- Adding reflection-based registration or generated factory code.

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel99.0-tool-dispatch-simplification](../test-suites/test-rel99.0-tool-dispatch-simplification.md)
