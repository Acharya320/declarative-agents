---
id: rel14.0-uc002-dolt-boundary-word
title: Declared Dolt words provision and update a versioned store
status: specified
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A profile provisions a domain database and schema from trusted ToolDef
configuration, writes one row with caller-supplied declared parameters, and
reads it through a separate query word. The write produces exactly one Dolt
commit whose message comes from trusted configuration. The caller cannot
provide SQL, connection, database, schema, or commit-message authority. The
domain database remains separate from the runtime checkpoint database.

<!-- field: actor -->
## Actor

A declarative machine invoking configured general-purpose Dolt boundary words

<!-- field: trigger -->
## Trigger

A composition starts against a real Dolt server where its configured domain
database does not yet exist.

<!-- field: flow -->
## Flow

- F1: MachineSpec dispatches dolt_provision, which resolves a trusted connection, creates the configured database, applies the ordered schema, and commits changed schema once (srd047 R1, R3).
- F2: A repeated dolt_provision sees the compatible schema and succeeds without destructive recreation or an empty duplicate commit (srd047 R3).
- F3: MachineSpec dispatches dolt_write with only schema-declared params; the word binds them to its configured mutation and creates one declared-template Dolt commit (srd047 R2, R5).
- F4: MachineSpec dispatches dolt_query with declared params; the word runs its configured read statement and returns bounded ordered structured rows without a commit (srd047 R4).
- F5: Traces and Results expose operation, database, row counts, and commit identity while omitting connection, credentials, raw SQL, and sensitive params (srd047 R6).
- F6: Startup rejects a Dolt word configured for the exact checkpoint server-and-database identity, while a separate database on the same server remains valid (srd047 R7).
- F7: A runtime request containing SQL, connection, database, schema, or commit_message is refused before database execution, and a failed mutation creates no commit (srd047 R2, R7).

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd047-dolt-boundary-word R1.1, R1.2, R1.3, R1.4, R1.5 | separate general-purpose provision, query, and write words |
| T2 | srd047-dolt-boundary-word R2.1, R2.2, R2.3, R2.4, R2.5, R2.6 | trusted SQL authority and caller-bound declared params |
| T3 | srd047-dolt-boundary-word R3.1, R3.2, R3.3, R3.4, R3.5, R3.6 | idempotent database and schema provisioning |
| T4 | srd047-dolt-boundary-word R4.1, R4.2, R4.3, R4.4, R4.5 | bounded read-only query |
| T5 | srd047-dolt-boundary-word R5.1, R5.2, R5.3, R5.4, R5.5, R5.6 | one commit per successful configured write |
| T6 | srd047-dolt-boundary-word R6.1, R6.2, R6.3, R6.4, R6.5 | signals, lifecycle declarations, traces, and errors |
| T7 | srd047-dolt-boundary-word R7.1, R7.2, R7.3, R7.4, R7.5, R7.6 | checkpoint isolation and real-Dolt integration proof |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Separate declared words remain general-purpose and MachineSpec owns provision, write, and query order.

<!-- field: traces -->
#### Traces

- srd047-dolt-boundary-word AC1

<!-- sc: S2 -->
### S2

The caller supplies only declared params, and all SQL and connection authority remains in trusted configuration.

<!-- field: traces -->
#### Traces

- srd047-dolt-boundary-word AC2

<!-- sc: S3 -->
### S3

A real Dolt server gains the configured database and schema once, and compatible reprovisioning creates no empty commit.

<!-- field: traces -->
#### Traces

- srd047-dolt-boundary-word AC3
- srd047-dolt-boundary-word AC7

<!-- sc: S4 -->
### S4

A configured write creates exactly one visible declared-message commit and a configured query returns the written row without committing.

<!-- field: traces -->
#### Traces

- srd047-dolt-boundary-word AC4
- srd047-dolt-boundary-word AC5
- srd047-dolt-boundary-word AC7

<!-- sc: S5 -->
### S5

Structured signals, results, traces, errors, and compensation metadata expose no secret or caller-controlled authority.

<!-- field: traces -->
#### Traces

- srd047-dolt-boundary-word AC6

<!-- sc: S6 -->
### S6

The word database and checkpoint database are distinct, caller authority fields are refused before execution, and a failed write creates no commit.

<!-- field: traces -->
#### Traces

- srd047-dolt-boundary-word AC7

<!-- field: failure_scenarios -->
## Failure scenarios

| id | scenario | expected |
| --- | --- | --- |
| E1 | A caller includes sql or database beside valid params. | Parameter validation emits CommandError before opening an operation against that authority. |
| E2 | A configured write statement fails after its transaction begins. | The transaction rolls back, no commit appears in Dolt history, and the error names statement_execute. |
| E3 | Word and checkpoint configuration resolve to the same server and database. | Startup validation rejects the profile before either backend mutates data. |

<!-- field: out_of_scope -->
## Out of scope

- A downstream consumer schema or task lifecycle.
- Arbitrary SQL consoles, ORMs, migration frameworks, and branch administration.
- Changes to the DoltCheckpoint schema or lifecycle contract.

<!-- field: test_suite_trace -->
## Test suite trace

test-rel14.0-agent-input-and-dolt
