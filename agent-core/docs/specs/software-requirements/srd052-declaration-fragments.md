---
id: srd052-declaration-fragments
title: Parameterized Declaration Fragments
depends_on:
  - srd_id: srd050-declaration-imports
    symbols_used: [unit, imports, closure usedness, duplicate-name error]
  - srd_id: srd051-declaration-types
    symbols_used: [schema subset, type references, signature]
  - srd_id: srd038-command-state-store
    symbols_used: [typed labels, selector path checks]
  - srd_id: srd006-yaml-machine-config
    symbols_used: [MachineSpec, exhaustiveness]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Imports (srd050) gave declarations headers, and type units with signatures
(srd051) gave them named shapes and defaulted contract blocks. Neither can
express the shape that produces most of the duplication that remains: one
block written out once per provider, per test case, or per agent, with two
or three values changed each time. Measured by `mage stats:reuse`, eighty-five
of ninety imported units have exactly one importer, which is an include by
another name, and the largest duplicated blocks are all instances of one
shape with different values.

In the terms of a programming language, the declarations have `#include`
and default members but no function templates. A fragment is that template:
a declaration unit with declared, typed parameters that an importer
expands with arguments. The loader performs the expansion before
strict decoding, so every existing validation runs on the expanded
result and nothing generated is ever committed.

<!-- field: goals -->
## Goals

- G1: Let a declaration unit declare typed parameters and be expanded by an importer with arguments.
- G2: Keep substitution hygienic, so an argument fills a hole in a shape and can never add structure to it.
- G3: Validate an expansion the way a constrained generic is validated, naming the fragment, the parameter, and the importer on failure.
- G4: Validate the expanded result through the same path as a hand-written unit, so a fragment can produce nothing a hand-written file could not.
- G5: Record provenance, so the dump and the usedness rule see expansions as they see imports.

<!-- req: R1 -->
## R1 - Fragment unit format

<!-- req: R1.1 -->
### R1.1

A fragment is a declaration file with `unit`, `params`, and exactly one body kind: `tools`, `rest`, or `stage`. `params` SHALL be a list of entries carrying `name` and `type`, and optionally `enum`, `default`, and `description`. `type` SHALL be one of string, integer, number, or boolean.

<!-- req: R1.2 -->
### R1.2

A fragment file MAY import type units and other non-fragment units under srd050 R1.2, and SHALL NOT expand other fragments. Nesting is outside this release.

<!-- req: R1.3 -->
### R1.3

A file carrying `params` and no body, or a body and no `params`, SHALL be a load error naming the file. A unit is a fragment or it is not; a file that is half of one declares nothing.

<!-- req: R2 -->
## R2 - Expansion

<!-- req: R2.1 -->
### R2.1

An importing file SHALL list expansions under `expand`, beside `imports`, each carrying `fragment` (a path relative to the importing file under the same rules as an import), `args` (a mapping from parameter name to value), and optionally `as` (a name prefix).

<!-- req: R2.2 -->
### R2.2

Every declared parameter without a `default` MUST be supplied. An argument naming no declared parameter SHALL be an error, and an argument whose value violates the parameter's `type` or `enum` SHALL be an error. Each error SHALL name the importing file, the fragment path, and the parameter.

<!-- req: R2.3 -->
### R2.3

A parameter is referenced as `$param(name)`. Substitution SHALL replace references only inside scalar values, never inside mapping names or list structure, so an argument cannot add a mapping, a list entry, or a second tool. A `$param(` remaining after substitution, or one naming an undeclared parameter, SHALL be a load error naming the fragment path and the line. The form cannot collide with environment expansion, whose pattern is `${NAME}`, nor with selectors, which begin `$.` or `$from(`.

<!-- req: R2.4 -->
### R2.4

Expansion SHALL happen in the loader before strict decoding, and the expanded document SHALL then be decoded and validated by the same path as a hand-written unit: strict decode, per-tool hooks, tool validation, REST definition validation. No expanded output SHALL be committed to the repository.

<!-- req: R2.5 -->
### R2.5

Expanded names MUST be unique in the closure under the existing duplicate rules. When `as` is given, every name the expansion produces SHALL be prefixed `<as>_`, so one fragment can be expanded more than once in one closure.

<!-- req: R2.6 -->
### R2.6

A fragment SHALL only be expanded. Listing a fragment under `imports` SHALL be a load error naming the importing file and the fragment, because a fragment's body holds unfilled references and is not a declaration until its arguments arrive.

<!-- req: R2.7 -->
### R2.7

Expanding one fragment twice with the same arguments and no distinguishing `as` SHALL fail with the existing duplicate-name error of srd050 R3, not one of its own.

<!-- req: R3 -->
## R3 - Usedness and provenance

<!-- req: R3.1 -->
### R3.1

An expansion whose produced definitions are all unselected SHALL be an unused import under srd050 R5, reported with the same evidence: the importing file, the fragment path, and the arguments.

<!-- req: R3.2 -->
### R3.2

Every produced definition SHALL record the fragment path and the arguments that produced it, and the canonical dump SHALL carry an `expansions` section listing, for each expansion, the fragment, the arguments, and the names produced.

<!-- req: R3.3 -->
### R3.3

The dump key is `expansions`, following the operation's name. The reuse-statistics JSON reports the same count under the key `instantiations` and SHALL keep it: that key is a metrics series with a committed baseline, and renaming it would make every historical figure incomparable for a vocabulary gain no reader of the JSON needs. The two surfaces differ because a diagnostic dump is read by a person against this specification, while a metrics key is read by a series. The divergence is deliberate, not an unfinished rename.

<!-- req: R4 -->
## R4 - Machine stage fragments

<!-- req: R4.1 -->
### R4.1

A `stage` body SHALL declare states, signals, transitions, and labels that splice into an importing machine under a prefix. The entry signal and the exit signals SHALL be parameters, so the importing machine decides how the stage is reached and left.

<!-- req: R4.2 -->
### R4.2

The spliced machine SHALL be validated as one machine by the label and exhaustiveness checks of srd038 and srd006, so an expanded stage is held to what a hand-written one is held to.

<!-- req: R4.3 -->
### R4.3

A machine file's `expand[].fragment` path MAY carry `${NAME:-default}` references, expanded at load time by the envexpand rules tool declarations use (srd013 R5.6), so a deployment selects one of several fixed stage variants whose topologies differ. Only the fragment path SHALL be expanded; states, transitions, and arguments stay literal, so the environment chooses among validated stages and cannot rewrite a machine. A reference without a default, a malformed reference, or a default variant whose file is missing SHALL fail the load naming the path as written, and the selected variant SHALL be spliced and validated under R4.2 and recorded in the dump by its resolved path (GH-2235).

<!-- field: non_goals -->
## Non-goals

- Fragments expanding other fragments; nesting waits for a release that needs it.
- Conditional or computed structure. A fragment is a fixed shape; arguments fill scalar holes in it.
- Loader words that resolve, typecheck, expand, or dump as builtin tools, and machines as tools for other machines. Both build on this and are specified separately.
- Committing any expanded or generated declaration.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

A fragment declaring params and one body loads, a fragment with params and no body fails naming the file, and a fragment expanding another fragment fails.

<!-- field: traces -->
#### Traces

- [R1.1](#r11)
- [R1.2](#r12)
- [R1.3](#r13)

<!-- ac: AC2 -->
### AC2

An expansion supplying every required argument produces definitions under the same validation as a hand-written unit; a missing argument, an unknown argument, a wrong type, and a value outside an enum each fail naming importer, fragment, and parameter; a `$param(` in a mapping name is not substituted and fails; a leftover reference fails naming the fragment and line; two expansions with distinct `as` coexist and two without fail as duplicates; a fragment listed under imports fails.

<!-- field: traces -->
#### Traces

- [R2.1](#r21)
- [R2.2](#r22)
- [R2.3](#r23)
- [R2.4](#r24)
- [R2.5](#r25)
- [R2.6](#r26)
- [R2.7](#r27)

<!-- ac: AC3 -->
### AC3

An expansion none of whose definitions is selected fails as an unused import with the fragment and arguments named, and the dump lists every expansion with fragment, arguments, and produced names.

<!-- field: traces -->
#### Traces

- [R3.1](#r31)
- [R3.2](#r32)

<!-- ac: AC6 -->
### AC6

The canonical dump names the section `expansions`, and the reuse-statistics JSON reports the same count under `instantiations`; neither key follows the other.

<!-- field: traces -->
#### Traces

- [R3.3](#r33)

<!-- ac: AC4 -->
### AC4

A stage fragment splices into a machine under a prefix with parameterized entry and exit signals, and the result passes or fails the label and exhaustiveness checks exactly as the equivalent hand-written machine does.

<!-- field: traces -->
#### Traces

- [R4.1](#r41)
- [R4.2](#r42)

<!-- ac: AC5 -->
### AC5

A machine expanding stages/s-${X:-a}.yaml beside s-a.yaml and s-b.yaml loads stage a with X unset and stage b with X=b, each validating as one machine with the dump naming the resolved path; a variant whose spliced transitions are dead fails naming the resolved fragment; a templated path without a default, or whose default file is missing, fails naming the path; packaging stages every variant.

<!-- field: traces -->
#### Traces

- [R4.3](#r43)

<!-- field: related_documents -->
## Related documents

- cfg-tool-declaration-format
- cfg-rest-tool-format
- cfg-machine-format
- srd050-declaration-imports
- srd051-declaration-types

<!-- field: implementation_status -->
## Implementation status

<!-- field: complete | form: "tree" -->
### Complete

1. GH-2081 implements tool-declaration fragments and the substitution engine (R1, R2, R3 for the tools body), with the control-fragments profile expanding one fragment twice and dumping deterministically.
1. GH-2082 implements REST-definition fragments through the REST import resolver (R1, R2, R3 for the rest body); a prefix renames every produced entry and the fragment's own references to them, and a fragment's body is decoded only after its arguments arrive.
1. - "GH-2083 implements machine stage fragments (R4)": a stage is spliced into the importing machine's lists and the whole is validated as one machine; control-fragments splices its patient control stage from a fragment.

<!-- field: planned -->
### Planned

- GH-2085 factors the top duplicated blocks into fragments and shared units, reporting mage stats:reuse before and after.
