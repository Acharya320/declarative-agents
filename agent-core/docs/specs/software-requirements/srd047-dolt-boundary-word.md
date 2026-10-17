---
id: srd047-dolt-boundary-word
title: Configured Dolt Boundary Words
depends_on:
  - srd_id: srd001-core-types
    symbols_used: [Command, Result, Signal, Receipt]
  - srd_id: srd003-command-dispatch
    symbols_used: [Dispatch]
  - srd_id: srd013-standard-tool-library
    symbols_used: [internal/tools package family]
  - srd_id: srd022-builtin-tools
    symbols_used: [BuiltinRegistry, ToolDef]
  - srd_id: srd025-rollback-lifecycle
    symbols_used: [receipt, compensation]
  - srd_id: srd036-dolt-state-persistence
    symbols_used: [DoltCheckpoint, checkpoint database identity]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Declarative agents need versioned relational storage for domain data without
receiving authority to invent SQL, database names, or connection targets.
The existing DoltCheckpoint adapter is a lifecycle backend with a fixed
runtime-state schema. It is not a general vocabulary word and must not become
one.

We need a general-purpose Dolt boundary family that follows the configured
REST tool safety model. Trusted ToolDef configuration owns the connection,
database, statements, parameter schemas, schema provisioning, and commit
messages. A caller selects no transport authority and supplies only declared
parameter values. Separate provision, query, and write words keep database
lifecycle, reads, and commit-producing mutations explicit in MachineSpec.

<!-- field: goals -->
## Goals

- G1: Define general-purpose Dolt provision, query, and write boundary words.
- G2: Keep connection, database, SQL statements, schemas, and commit messages under trusted configuration.
- G3: Let callers fill typed declared parameters without supplying SQL or connection authority.
- G4: Provision a database and schema idempotently from a declared word.
- G5: Commit every successful write exactly once with a configured message template and return structured results.
- G6: Keep domain-word databases separate from the Dolt checkpoint backend and prove behavior against a real Dolt server.

<!-- req: R1 -->
## R1 - Word family and ownership

<!-- req: R1.1 -->
### R1.1

Agent Core SHALL implement the family under internal/tools/dolt and SHALL register it through the standard catalog and registry path.

<!-- req: R1.2 -->
### R1.2

Shared declarations SHALL live under tools/builtin/dolt and SHALL expose distinct dolt_provision, dolt_query, and dolt_write init entries.

<!-- req: R1.3 -->
### R1.3

Each configured ToolDef SHALL bind exactly one provision, query, or write operation; MachineSpec SHALL own sequencing among operations.

<!-- req: R1.4 -->
### R1.4

Every entry SHALL be a boundary word with typed input and output, emitted signals, database side effects, reversibility, undo or compensation, errors, and relationships.

<!-- req: R1.5 -->
### R1.5

The family SHALL be general-purpose and SHALL NOT name or encode any task, blackboard, application, or other consumer domain.

<!-- req: R2 -->
## R2 - Trusted configuration and parameter authority

<!-- req: R2.1 -->
### R2.1

Trusted ToolDef configuration SHALL own a connection reference, literal database name, operation id, operation kind, SQL statement or schema statement list, parameter schema, result limits, timeout, and applicable commit-message template.

<!-- req: R2.2 -->
### R2.2

Connection credentials SHALL resolve from a named trusted runtime reference that cmd/agent supplies through FactoryDeps from --dolt-connection NAME=DSN. Process environment is not a resolution channel. Callers and model output SHALL NOT supply or override DSNs, hosts, ports, users, passwords, TLS policy, database names, SQL text, schema text, operation kinds, or commit templates.

<!-- req: R2.3 -->
### R2.3

Runtime input SHALL contain only params whose keys and scalar or structured types are declared by the operation's parameter schema; unknown, missing, or invalid params SHALL fail before database execution.

<!-- req: R2.4 -->
### R2.4

SQL SHALL use driver parameter binding for caller values. Runtime SHALL NOT construct SQL by concatenating, formatting, or template-evaluating caller strings.

<!-- req: R2.5 -->
### R2.5

A runtime input containing statement, sql, query, dsn, connection, database, schema, commit_message, or another authority field SHALL be refused as schema-invalid rather than ignored.

<!-- req: R2.6 -->
### R2.6

Startup validation SHALL reject empty statements, unknown config fields, undeclared placeholders, unused required parameters, multi-statement query or write entries, and operation-kind mismatches.

<!-- req: R3 -->
## R3 - Database provisioning

<!-- req: R3.1 -->
### R3.1

dolt_provision SHALL connect to the configured Dolt server without accepting a caller-selected database and SHALL create the configured database when it does not exist.

<!-- req: R3.2 -->
### R3.2

dolt_provision SHALL apply only the ordered DDL statements listed in trusted schema configuration and SHALL accept no runtime params.

<!-- req: R3.3 -->
### R3.3

Provisioning SHALL be idempotent: a second invocation against the same compatible database and schema SHALL succeed without destructive recreation or an empty duplicate commit.

<!-- req: R3.4 -->
### R3.4

After creating or changing schema, dolt_provision SHALL stage the database changes and create exactly one Dolt commit using its configured schema commit-message template.

<!-- req: R3.5 -->
### R3.5

A schema conflict, incompatible existing object, permission failure, or failed statement SHALL roll back the current SQL transaction when supported, create no Dolt commit, and return a structured DoltCommandError.

<!-- req: R3.6 -->
### R3.6

Provision output SHALL include operation, database, created, schema_changed, commit_hash when changed, and bounded diagnostics without connection or credential material.

<!-- req: R4 -->
## R4 - Read-only query word

<!-- req: R4.1 -->
### R4.1

dolt_query SHALL accept only a single configured read statement whose validated statement kind is SELECT, SHOW, DESCRIBE, EXPLAIN, or a read-only common-table expression.

<!-- req: R4.2 -->
### R4.2

dolt_query SHALL execute with bound declared params and SHALL never stage changes or call DOLT_COMMIT.

<!-- req: R4.3 -->
### R4.3

Query output SHALL include operation, database, ordered column metadata, rows as structured values, row_count, truncated, and elapsed duration.

<!-- req: R4.4 -->
### R4.4

Configuration SHALL set maximum rows and maximum encoded bytes; reaching either limit SHALL return a deterministic bounded result with truncated true.

<!-- req: R4.5 -->
### R4.5

A statement that mutates data, schema, branches, commits, or server configuration SHALL fail startup validation for dolt_query.

<!-- req: R5 -->
## R5 - Write and commit word

<!-- req: R5.1 -->
### R5.1

dolt_write SHALL accept one configured INSERT, UPDATE, DELETE, or other explicitly allowed data-mutation statement and execute it with bound declared params inside one SQL transaction.

<!-- req: R5.2 -->
### R5.2

On a successful mutation, dolt_write SHALL stage the resulting Dolt changes and call DOLT_COMMIT exactly once before reporting success.

<!-- req: R5.3 -->
### R5.3

The commit message SHALL be rendered from a trusted template using only configured literal text and declared, schema-valid parameter values. Unknown placeholders or caller-supplied message text SHALL be rejected.

<!-- req: R5.4 -->
### R5.4

Write output SHALL include operation, database, rows_affected, commit_hash, committed true, and elapsed duration; it SHALL NOT expose the DSN, credentials, raw statement, or secret params.

<!-- req: R5.5 -->
### R5.5

A zero-row successful write SHALL follow explicit trusted config, commit_on_no_change. The default SHALL be false and SHALL return committed false with no empty commit.

<!-- req: R5.6 -->
### R5.6

A statement failure, staging failure, commit failure, timeout, or cancellation SHALL return DoltCommandError, report its stage, and SHALL NOT report committed true. SQL failure before DOLT_COMMIT SHALL roll back and create no commit.

<!-- req: R6 -->
## R6 - Signals, lifecycle, and observability

<!-- req: R6.1 -->
### R6.1

Successful provision, query, and write words SHALL emit DoltProvisioned, DoltRowsRead, and DoltCommitted respectively; expected configured no-change writes MAY emit DoltNoChange; all infrastructure and validation failures SHALL emit CommandError.

<!-- req: R6.2 -->
### R6.2

dolt_query SHALL declare no side effects and use noop Undo. dolt_provision and dolt_write SHALL declare external database mutation and classify it as compensatable or irreversible according to configured compensation policy.

<!-- req: R6.3 -->
### R6.3

A mutating word marked compensatable SHALL encode an opaque receipt with operation identity, database identity, and commit hash for a separately declared compensation action. The checkpoint adapter and engine SHALL not parse it.

<!-- req: R6.4 -->
### R6.4

Dispatch traces SHALL record operation id, word kind, configured server identity, database, rows read or affected, commit presence and hash, elapsed time, and resulting signal; they SHALL redact credentials, SQL text, commit-message values marked sensitive, and secret params.

<!-- req: R6.5 -->
### R6.5

Errors SHALL distinguish config_validation, connection_resolution, database_create, schema_apply, parameter_validation, statement_execute, result_scan, stage, commit, timeout, cancellation, and compensation stages.

<!-- req: R7 -->
## R7 - Checkpoint separation and real integration proof

<!-- req: R7.1 -->
### R7.1

The Dolt word connection and database SHALL be configured independently from --dolt-dsn and the srd036 DoltCheckpoint backend.

<!-- req: R7.2 -->
### R7.2

When both are enabled, startup validation SHALL reject an exact normalized server-and-database identity shared by a Dolt word and DoltCheckpoint. Separate databases on the same configured server MAY be used.

<!-- req: R7.3 -->
### R7.3

Dolt words SHALL NOT read or write the checkpoint tables, branches, run commits, Position, Execution, tool_outputs, or receipts owned by srd036.

<!-- req: R7.4 -->
### R7.4

GH-1732 SHALL add mage integration:doltWord, or an equivalently named standard integration target, that launches or connects to a real Dolt instance rather than a mock SQL driver.

<!-- req: R7.5 -->
### R7.5

The real-Dolt gate SHALL prove configured database creation, ordered schema application, idempotent reprovisioning, bound parameter query and write, one visible declared-message commit per successful write, structured results, and separate checkpoint and word databases.

<!-- req: R7.6 -->
### R7.6

The real-Dolt gate SHALL prove that caller-supplied SQL, connection, database, and commit-message fields are refused before execution and that a failed write creates no commit.

<!-- field: non_goals -->
## Non-goals

- Replacing or extending DoltCheckpoint, its five-table schema, branch-per-run history, or Revert behavior.
- Letting a model or runtime caller submit arbitrary SQL, connection targets, database names, schema, or commit messages.
- Providing a general migration engine, ORM, query builder, branch manager, or Dolt administration shell.
- Encoding a particular downstream consumer's schema or lifecycle.
- Hiding provision, query, write, or compensation order inside one sentence-shaped tool.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

The general-purpose internal/tools/dolt family exposes separate provision, query, and write boundary entries through normal declarations and registration, with MachineSpec owning composition.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R1.4
- R1.5

<!-- ac: AC2 -->
### AC2

Trusted config owns every connection and SQL authority field, callers provide only schema-valid bound params, authority-shaped runtime fields are refused, and startup rejects unsafe statement definitions.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3
- R2.4
- R2.5
- R2.6

<!-- ac: AC3 -->
### AC3

Provisioning creates the configured database, applies ordered schema idempotently, commits changed schema once, creates no empty duplicate commit, and returns structured redacted output or a staged error.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3
- R3.4
- R3.5
- R3.6

<!-- ac: AC4 -->
### AC4

Query accepts only validated read statements, binds params, returns bounded ordered structured rows, and never stages or commits.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
- R4.3
- R4.4
- R4.5

<!-- ac: AC5 -->
### AC5

Write executes one configured mutation transaction, creates exactly one declared-template Dolt commit on change, handles configured no-change behavior, returns structured commit data, and reports failures without false commit success.

<!-- field: traces -->
#### Traces

- R5.1
- R5.2
- R5.3
- R5.4
- R5.5
- R5.6

<!-- ac: AC6 -->
### AC6

Signals, side effects, receipts, compensation, traces, redaction, and staged error taxonomy follow the boundary-word and lifecycle contracts.

<!-- field: traces -->
#### Traces

- R6.1
- R6.2
- R6.3
- R6.4
- R6.5

<!-- ac: AC7 -->
### AC7

Word and checkpoint databases are independently configured and cannot share one normalized database, and a real-Dolt integration gate proves provisioning, schema, parameterized reads and writes, commit history, isolation, refusal of caller authority, and no commit after failure.

<!-- field: traces -->
#### Traces

- R7.1
- R7.2
- R7.3
- R7.4
- R7.5
- R7.6

<!-- field: test_trace_plan | form: "typed table" -->
## Test trace plan

| test_suite | covers |
| --- | --- |
| test-rel14.0-agent-input-and-dolt | [AC1, AC2, AC3, AC4, AC5, AC6, AC7] |

<!-- field: related_documents -->
## Related documents

- docs/ARCHITECTURE.md
- docs/specs/config-formats/tool-authoring-standard.md
- docs/specs/software-requirements/srd027-rest-tool-runtime.md
- docs/guides/rest-tool.md

<!-- field: implementation_status -->
## Implementation status

```yaml
specified:
    - Trusted operation config, provision/query/write contracts, commit semantics, checkpoint separation, and real-Dolt proof obligations.
implemented:
    - "GH-1732: runtime words, declarations, unit tests, and mage integration:doltWord real-Dolt gate."
```
