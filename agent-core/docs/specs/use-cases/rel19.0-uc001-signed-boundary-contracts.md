---
id: rel19.0-uc001-signed-boundary-contracts
title: A boundary tool states its contract through a signature and its effect blocks
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A declaration author maintains a boundary or stateful_internal tool whose
descriptive prose -- problem, goals, requirements, non_goals, and
relationships -- restates the description and is checked only for presence.
They add a signature naming the tool's signals, keep its output schema and
its explicit side_effects, reversibility, and undo blocks, and delete the
prose. The loader and the corpus audit accept the tool as complete, and the
declaration legacy gate drops it from the unsigned list.

<!-- field: actor -->
## Actor

A declaration author reducing contract ceremony on boundary and stateful tools

<!-- field: trigger -->
## Trigger

The declaration legacy gate lists a boundary, stateful_internal, or
uncategorized tool, or an application migration converts its tools.

<!-- field: flow -->
## Flow

- F1: The author adds `signature` with `emits` equal to the tool's signals and deletes the problem, goals, requirements, non_goals, and relationships blocks.
- F2: The loader applies R6.6 and defaults the descriptive prose from the description; it defaults no effect block of a boundary tool and only reversibility and undo of a stateful_internal tool.
- F3: The corpus audit finds no missing contract field; a boundary tool that also dropped an effect block fails naming that block.
- F4: The declaration legacy gate no longer lists the tool as untyped, and the stale-entry check requires the author to delete its baseline line.
- F5: A tool without a category is listed as uncategorized until the author states one.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd051-declaration-types](../software-requirements/srd051-declaration-types.md) [R6.6](../software-requirements/srd051-declaration-types.md#r66) | a signature discharges the descriptive prose for every category |
| T2 | [srd051-declaration-types](../software-requirements/srd051-declaration-types.md) [R6.8](../software-requirements/srd051-declaration-types.md#r68), [R6.9](../software-requirements/srd051-declaration-types.md#r69) | stateful_internal and boundary effect blocks stay explicit |
| T3 | [srd051-declaration-types](../software-requirements/srd051-declaration-types.md) [R6.14](../software-requirements/srd051-declaration-types.md#r614) | boundary and stateful_internal tools declare a signature |
| T4 | [srd051-declaration-types](../software-requirements/srd051-declaration-types.md) [R6.15](../software-requirements/srd051-declaration-types.md#r615) | every tool states a category |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

A boundary tool with an emits-only signature, an output schema, and its three effect blocks is complete, and dropping any effect block fails naming it.

<!-- field: traces -->
#### Traces

- [srd051-declaration-types](../software-requirements/srd051-declaration-types.md) [AC8](../software-requirements/srd051-declaration-types.md#ac8)

<!-- sc: S2 -->
### S2

The legacy gate lists unsigned boundary and stateful_internal tools as untyped and tools without a category as uncategorized, and lists a signed boundary tool as neither.

<!-- field: traces -->
#### Traces

- [srd051-declaration-types](../software-requirements/srd051-declaration-types.md) [AC8](../software-requirements/srd051-declaration-types.md#ac8)

<!-- field: out_of_scope -->
## Out of scope

- Naming an output type on a boundary tool whose label the runtime result types (srd038 R2.19).
- Changing which effect blocks a category defaults.
- Promoting R6.14 and R6.15 into corpus audit errors before the baseline is empty.
