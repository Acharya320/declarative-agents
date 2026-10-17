---
id: rel07.0-uc002-two-plane-revert
title: Reverting to an earlier step reaps later two-plane rows while the receipt walk reverses the surviving step
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

The runtime runs three steps, each persisting a tool_outputs row and a receipts
row in one per-step transaction, then reverts to step 1. The tool_outputs and
receipts rows for the later steps are gone with the branch reset, the
command-state store's later entries are invalidated by construction, and the
rollback receipt walk still reverses step 1's external effect from the surviving
receipts row. No forward selector reaches any receipt at any point.

<!-- field: actor -->
## Actor

Agent runtime performing git-style rollback

<!-- field: trigger -->
## Trigger

A rollback request targets an earlier step after later steps have persisted
output and receipts.

<!-- field: flow -->
## Flow

- F1: Three steps run; each writes its tool_outputs row and its receipts row inside one per-step transaction (srd036 R3, R4.4).
- F2: Revert(run_id, 1) resets the run branch to the step-1 commit (srd036 R6).
- F3: The tool_outputs and receipts rows for steps 2 and 3 are gone by branch construction; the store's later entries are invalidated without a separate walk (srd038 R4, srd036 R4).
- F4: The rollback receipt walk reads the surviving step-1 receipt from the receipts plane and reverses that step's external effect (srd036 R3, R6).
- F5: A forward selector cannot reach any receipt at any point; only tool_outputs is store-addressable (srd036 R3.6, srd038 R3).

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd036-dolt-state-persistence](../software-requirements/srd036-dolt-state-persistence.md) [R2.1](../software-requirements/srd036-dolt-state-persistence.md#r21), [R2.2](../software-requirements/srd036-dolt-state-persistence.md#r22), [R2.3](../software-requirements/srd036-dolt-state-persistence.md#r23), [R2.4](../software-requirements/srd036-dolt-state-persistence.md#r24) | five-table schema; tool_outputs forward plane and receipts reverse plane |
| T2 | [srd036-dolt-state-persistence](../software-requirements/srd036-dolt-state-persistence.md) [R4.1](../software-requirements/srd036-dolt-state-persistence.md#r41), [R4.2](../software-requirements/srd036-dolt-state-persistence.md#r42), [R4.3](../software-requirements/srd036-dolt-state-persistence.md#r43), [R4.4](../software-requirements/srd036-dolt-state-persistence.md#r44), [R4.5](../software-requirements/srd036-dolt-state-persistence.md#r45), [R4.6](../software-requirements/srd036-dolt-state-persistence.md#r46) | both plane rows written in one per-step transaction; terminal and revert reaping |
| T3 | [srd036-dolt-state-persistence](../software-requirements/srd036-dolt-state-persistence.md) [R6.1](../software-requirements/srd036-dolt-state-persistence.md#r61), [R6.2](../software-requirements/srd036-dolt-state-persistence.md#r62), [R6.3](../software-requirements/srd036-dolt-state-persistence.md#r63), [R6.4](../software-requirements/srd036-dolt-state-persistence.md#r64), [R6.5](../software-requirements/srd036-dolt-state-persistence.md#r65), [R6.6](../software-requirements/srd036-dolt-state-persistence.md#r66) | Revert resets the branch; the receipt walk reverses surviving steps |
| T4 | [srd038-command-state-store](../software-requirements/srd038-command-state-store.md) [R4.1](../software-requirements/srd038-command-state-store.md#r41), [R4.2](../software-requirements/srd038-command-state-store.md#r42), [R4.3](../software-requirements/srd038-command-state-store.md#r43), [R4.4](../software-requirements/srd038-command-state-store.md#r44) | branch-scoped entry lifetime; Revert invalidates later entries by construction |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

A three-step run persists a tool_outputs row and a receipts row per step and round-trips both planes on Load.

<!-- field: traces -->
#### Traces

- [srd036-dolt-state-persistence](../software-requirements/srd036-dolt-state-persistence.md) [AC2](../software-requirements/srd036-dolt-state-persistence.md#ac2)

<!-- sc: S2 -->
### S2

After Revert to step 1, the tool_outputs and receipts rows for steps 2 and 3 do not resolve and the store's later entries are invalid without a separate invalidation walk.

<!-- field: traces -->
#### Traces

- [srd038-command-state-store](../software-requirements/srd038-command-state-store.md) [AC4](../software-requirements/srd038-command-state-store.md#ac4)
- [srd036-dolt-state-persistence](../software-requirements/srd036-dolt-state-persistence.md) [AC4](../software-requirements/srd036-dolt-state-persistence.md#ac4)

<!-- sc: S3 -->
### S3

The receipt walk reverses step 1's external effect from the surviving receipts row.

<!-- field: traces -->
#### Traces

- [srd036-dolt-state-persistence](../software-requirements/srd036-dolt-state-persistence.md) [AC3](../software-requirements/srd036-dolt-state-persistence.md#ac3)
- [srd036-dolt-state-persistence](../software-requirements/srd036-dolt-state-persistence.md) [AC6](../software-requirements/srd036-dolt-state-persistence.md#ac6)

<!-- sc: S4 -->
### S4

No forward selector resolves a receipt before or after Revert.

<!-- field: traces -->
#### Traces

- [srd036-dolt-state-persistence](../software-requirements/srd036-dolt-state-persistence.md) [AC3](../software-requirements/srd036-dolt-state-persistence.md#ac3)
- [srd038-command-state-store](../software-requirements/srd038-command-state-store.md) [AC3](../software-requirements/srd038-command-state-store.md#ac3)

<!-- field: failure_scenarios -->
## Failure scenarios

| id | scenario | expected |
| --- | --- | --- |
| E1 | A step commits a tool_outputs row without a matching receipts row. | The per-step transaction is rejected so a partial step is never committed (srd036 R4.4). |
| E2 | Revert targets a step_index with no recorded commit. | A typed error is reported (srd036 R6.5). |
| E3 | Rollback attempts to read a receipt through a forward selector. | Structurally impossible; receipts are not store-addressable (srd036 R3.6, srd038 R3.3). |

<!-- field: out_of_scope -->
## Out of scope

- The $from selector surface (covered by rel07.0-uc001).
- File and HTTP compensation semantics beyond the receipt walk (owned by srd025 and srd026).

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel07.0-command-state](../test-suites/test-rel07.0-command-state.md)
