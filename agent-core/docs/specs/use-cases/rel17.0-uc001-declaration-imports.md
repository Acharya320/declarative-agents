---
id: rel17.0-uc001-declaration-imports
title: A profile composes named tool and REST declaration units
status: specified
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A profile author extracts shared tool and REST declarations into named units
and imports them by relative path. The closure loader resolves one explicit
dependency graph, preserves source provenance, rejects ambiguous or unused
edges, selects imported tools, compiles merged REST bindings, and exposes the
resolved result through validation and the canonical dump.

<!-- field: actor -->
## Actor

A profile author extracting reusable declaration units

<!-- field: trigger -->
## Trigger

A profile declaration root names a unit whose imports reference shared tool
or REST declaration files.

<!-- field: flow -->
## Flow

- F1: The loader strictly decodes the root and each imported source and validates unit names and relative import paths.
- F2: It resolves canonical sources in first-discovery depth-first postorder, reusing one source across diamond paths and rejecting cycles or unit-name collisions.
- F3: Tool imports reject duplicate names before applying an explicit local override, then tools.yaml selects the merged vocabulary.
- F4: REST imports merge raw map families with provenance, reject collisions, resolve declaration-relative OpenAPI sources, and compile the merged definition.
- F5: Closure validation traverses selected-tool and REST reference roots, marks every contributing import edge, and rejects unused edges deterministically.
- F6: Validation, dump, startup, and profile audit consume the same closure; the dump lists every imported source once.
- F7: A legacy-only profile continues to load includes with one warning per including file, while a closure mixing imports and includes fails.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd050-declaration-imports](../software-requirements/srd050-declaration-imports.md) [R1.1](../software-requirements/srd050-declaration-imports.md#r11), [R1.2](../software-requirements/srd050-declaration-imports.md#r12), [R1.3](../software-requirements/srd050-declaration-imports.md#r13), [R1.4](../software-requirements/srd050-declaration-imports.md#r14), [R1.6](../software-requirements/srd050-declaration-imports.md#r16), [R2.1](../software-requirements/srd050-declaration-imports.md#r21), [R2.2](../software-requirements/srd050-declaration-imports.md#r22), [R2.3](../software-requirements/srd050-declaration-imports.md#r23), [R2.4](../software-requirements/srd050-declaration-imports.md#r24), [R2.5](../software-requirements/srd050-declaration-imports.md#r25), [R2.6](../software-requirements/srd050-declaration-imports.md#r26), [R2.7](../software-requirements/srd050-declaration-imports.md#r27), [R2.8](../software-requirements/srd050-declaration-imports.md#r28), [R2.9](../software-requirements/srd050-declaration-imports.md#r29) | syntax, graph resolution, provenance |
| T2 | [srd050-declaration-imports](../software-requirements/srd050-declaration-imports.md) [R3.1](../software-requirements/srd050-declaration-imports.md#r31), [R3.2](../software-requirements/srd050-declaration-imports.md#r32), [R3.3](../software-requirements/srd050-declaration-imports.md#r33), [R3.4](../software-requirements/srd050-declaration-imports.md#r34), [R3.5](../software-requirements/srd050-declaration-imports.md#r35), [R3.6](../software-requirements/srd050-declaration-imports.md#r36) | tool merge and override |
| T3 | [srd050-declaration-imports](../software-requirements/srd050-declaration-imports.md) [R4.1](../software-requirements/srd050-declaration-imports.md#r41), [R4.2](../software-requirements/srd050-declaration-imports.md#r42), [R4.3](../software-requirements/srd050-declaration-imports.md#r43), [R4.4](../software-requirements/srd050-declaration-imports.md#r44), [R4.5](../software-requirements/srd050-declaration-imports.md#r45), [R4.6](../software-requirements/srd050-declaration-imports.md#r46), [R4.7](../software-requirements/srd050-declaration-imports.md#r47) | REST merge and OpenAPI compilation |
| T4 | [srd050-declaration-imports](../software-requirements/srd050-declaration-imports.md) [R5.1](../software-requirements/srd050-declaration-imports.md#r51), [R5.2](../software-requirements/srd050-declaration-imports.md#r52), [R5.3](../software-requirements/srd050-declaration-imports.md#r53), [R5.4](../software-requirements/srd050-declaration-imports.md#r54), [R5.5](../software-requirements/srd050-declaration-imports.md#r55), [R5.6](../software-requirements/srd050-declaration-imports.md#r56), [R5.7](../software-requirements/srd050-declaration-imports.md#r57), [R5.8](../software-requirements/srd050-declaration-imports.md#r58) | import-edge usedness |
| T5 | [srd050-declaration-imports](../software-requirements/srd050-declaration-imports.md) [R6.3](../software-requirements/srd050-declaration-imports.md#r63), [R6.4](../software-requirements/srd050-declaration-imports.md#r64), [R6.5](../software-requirements/srd050-declaration-imports.md#r65) | compatibility and shared diagnostics |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Two-level tool and REST imports resolve deterministically and every source appears once in repeated canonical dumps.

<!-- field: traces -->
#### Traces

- [srd050-declaration-imports](../software-requirements/srd050-declaration-imports.md) [AC1](../software-requirements/srd050-declaration-imports.md#ac1)

<!-- sc: S2 -->
### S2

Duplicate units, declarations, cycles, and invalid tool overrides fail with complete unit/path evidence.

<!-- field: traces -->
#### Traces

- [srd050-declaration-imports](../software-requirements/srd050-declaration-imports.md) [AC2](../software-requirements/srd050-declaration-imports.md#ac2)
- [srd050-declaration-imports](../software-requirements/srd050-declaration-imports.md) [AC3](../software-requirements/srd050-declaration-imports.md#ac3)

<!-- sc: S3 -->
### S3

REST map families and OpenAPI bindings merge with declaring-source path semantics and duplicate/version checks.

<!-- field: traces -->
#### Traces

- [srd050-declaration-imports](../software-requirements/srd050-declaration-imports.md) [AC4](../software-requirements/srd050-declaration-imports.md#ac4)

<!-- sc: S4 -->
### S4

Used transitive and diamond edges pass while unused tool and REST edges fail in deterministic order.

<!-- field: traces -->
#### Traces

- [srd050-declaration-imports](../software-requirements/srd050-declaration-imports.md) [AC5](../software-requirements/srd050-declaration-imports.md#ac5)

<!-- sc: S5 -->
### S5

Legacy includes remain compatible but cannot coexist with imports in one closure.

<!-- field: traces -->
#### Traces

- [srd050-declaration-imports](../software-requirements/srd050-declaration-imports.md) AC6

<!-- sc: S6 -->
### S6

Every entry path reports the same import error and the control-imports proof matches control except identity and files.

<!-- field: traces -->
#### Traces

- [srd050-declaration-imports](../software-requirements/srd050-declaration-imports.md) [AC7](../software-requirements/srd050-declaration-imports.md#ac7)
- [srd050-declaration-imports](../software-requirements/srd050-declaration-imports.md) [AC8](../software-requirements/srd050-declaration-imports.md#ac8)

<!-- field: out_of_scope -->
## Out of scope

- Named types, signatures, typed labels, and machine exhaustiveness.
- Conditional, versioned, remote, or generated imports.
- Removal of legacy includes.
