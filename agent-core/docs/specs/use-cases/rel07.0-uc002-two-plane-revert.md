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
| T1 | srd036-dolt-state-persistence R2.1, R2.2, R2.3, R2.4 | five-table schema; tool_outputs forward plane and receipts reverse plane |
| T2 | srd036-dolt-state-persistence R4.1, R4.2, R4.3, R4.4, R4.5, R4.6 | both plane rows written in one per-step transaction; terminal and revert reaping |
| T3 | srd036-dolt-state-persistence R6.1, R6.2, R6.3, R6.4, R6.5, R6.6 | Revert resets the branch; the receipt walk reverses surviving steps |
| T4 | srd038-command-state-store R4.1, R4.2, R4.3, R4.4 | branch-scoped entry lifetime; Revert invalidates later entries by construction |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

A three-step run persists a tool_outputs row and a receipts row per step and round-trips both planes on Load.

<!-- field: traces -->
#### Traces

- srd036-dolt-state-persistence AC2

<!-- sc: S2 -->
### S2

After Revert to step 1, the tool_outputs and receipts rows for steps 2 and 3 do not resolve and the store's later entries are invalid without a separate invalidation walk.

<!-- field: traces -->
#### Traces

- srd038-command-state-store AC4
- srd036-dolt-state-persistence AC4

<!-- sc: S3 -->
### S3

The receipt walk reverses step 1's external effect from the surviving receipts row.

<!-- field: traces -->
#### Traces

- srd036-dolt-state-persistence AC3
- srd036-dolt-state-persistence AC6

<!-- sc: S4 -->
### S4

No forward selector resolves a receipt before or after Revert.

<!-- field: traces -->
#### Traces

- srd036-dolt-state-persistence AC3
- srd038-command-state-store AC3

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

test-rel07.0-command-state
