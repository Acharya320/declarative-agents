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
| T1 | [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [R1.1](../software-requirements/srd047-dolt-boundary-word.md#r11), [R1.2](../software-requirements/srd047-dolt-boundary-word.md#r12), [R1.3](../software-requirements/srd047-dolt-boundary-word.md#r13), [R1.4](../software-requirements/srd047-dolt-boundary-word.md#r14), [R1.5](../software-requirements/srd047-dolt-boundary-word.md#r15) | separate general-purpose provision, query, and write words |
| T2 | [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [R2.1](../software-requirements/srd047-dolt-boundary-word.md#r21), [R2.2](../software-requirements/srd047-dolt-boundary-word.md#r22), [R2.3](../software-requirements/srd047-dolt-boundary-word.md#r23), [R2.4](../software-requirements/srd047-dolt-boundary-word.md#r24), [R2.5](../software-requirements/srd047-dolt-boundary-word.md#r25), [R2.6](../software-requirements/srd047-dolt-boundary-word.md#r26) | trusted SQL authority and caller-bound declared params |
| T3 | [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [R3.1](../software-requirements/srd047-dolt-boundary-word.md#r31), [R3.2](../software-requirements/srd047-dolt-boundary-word.md#r32), [R3.3](../software-requirements/srd047-dolt-boundary-word.md#r33), [R3.4](../software-requirements/srd047-dolt-boundary-word.md#r34), [R3.5](../software-requirements/srd047-dolt-boundary-word.md#r35), [R3.6](../software-requirements/srd047-dolt-boundary-word.md#r36) | idempotent database and schema provisioning |
| T4 | [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [R4.1](../software-requirements/srd047-dolt-boundary-word.md#r41), [R4.2](../software-requirements/srd047-dolt-boundary-word.md#r42), [R4.3](../software-requirements/srd047-dolt-boundary-word.md#r43), [R4.4](../software-requirements/srd047-dolt-boundary-word.md#r44), [R4.5](../software-requirements/srd047-dolt-boundary-word.md#r45) | bounded read-only query |
| T5 | [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [R5.1](../software-requirements/srd047-dolt-boundary-word.md#r51), [R5.2](../software-requirements/srd047-dolt-boundary-word.md#r52), [R5.3](../software-requirements/srd047-dolt-boundary-word.md#r53), [R5.4](../software-requirements/srd047-dolt-boundary-word.md#r54), [R5.5](../software-requirements/srd047-dolt-boundary-word.md#r55), [R5.6](../software-requirements/srd047-dolt-boundary-word.md#r56) | one commit per successful configured write |
| T6 | [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [R6.1](../software-requirements/srd047-dolt-boundary-word.md#r61), [R6.2](../software-requirements/srd047-dolt-boundary-word.md#r62), [R6.3](../software-requirements/srd047-dolt-boundary-word.md#r63), [R6.4](../software-requirements/srd047-dolt-boundary-word.md#r64), [R6.5](../software-requirements/srd047-dolt-boundary-word.md#r65) | signals, lifecycle declarations, traces, and errors |
| T7 | [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [R7.1](../software-requirements/srd047-dolt-boundary-word.md#r71), [R7.2](../software-requirements/srd047-dolt-boundary-word.md#r72), [R7.3](../software-requirements/srd047-dolt-boundary-word.md#r73), [R7.4](../software-requirements/srd047-dolt-boundary-word.md#r74), [R7.5](../software-requirements/srd047-dolt-boundary-word.md#r75), [R7.6](../software-requirements/srd047-dolt-boundary-word.md#r76) | checkpoint isolation and real-Dolt integration proof |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Separate declared words remain general-purpose and MachineSpec owns provision, write, and query order.

<!-- field: traces -->
#### Traces

- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC1](../software-requirements/srd047-dolt-boundary-word.md#ac1)

<!-- sc: S2 -->
### S2

The caller supplies only declared params, and all SQL and connection authority remains in trusted configuration.

<!-- field: traces -->
#### Traces

- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC2](../software-requirements/srd047-dolt-boundary-word.md#ac2)

<!-- sc: S3 -->
### S3

A real Dolt server gains the configured database and schema once, and compatible reprovisioning creates no empty commit.

<!-- field: traces -->
#### Traces

- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC3](../software-requirements/srd047-dolt-boundary-word.md#ac3)
- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC7](../software-requirements/srd047-dolt-boundary-word.md#ac7)

<!-- sc: S4 -->
### S4

A configured write creates exactly one visible declared-message commit and a configured query returns the written row without committing.

<!-- field: traces -->
#### Traces

- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC4](../software-requirements/srd047-dolt-boundary-word.md#ac4)
- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC5](../software-requirements/srd047-dolt-boundary-word.md#ac5)
- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC7](../software-requirements/srd047-dolt-boundary-word.md#ac7)

<!-- sc: S5 -->
### S5

Structured signals, results, traces, errors, and compensation metadata expose no secret or caller-controlled authority.

<!-- field: traces -->
#### Traces

- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC6](../software-requirements/srd047-dolt-boundary-word.md#ac6)

<!-- sc: S6 -->
### S6

The word database and checkpoint database are distinct, caller authority fields are refused before execution, and a failed write creates no commit.

<!-- field: traces -->
#### Traces

- [srd047-dolt-boundary-word](../software-requirements/srd047-dolt-boundary-word.md) [AC7](../software-requirements/srd047-dolt-boundary-word.md#ac7)

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

[test-rel14.0-agent-input-and-dolt](../test-suites/test-rel14.0-agent-input-and-dolt.md)
