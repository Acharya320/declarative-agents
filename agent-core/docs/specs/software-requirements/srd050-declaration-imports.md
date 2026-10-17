---
id: srd050-declaration-imports
title: Named Declaration Units and Explicit Imports
depends_on:
  - srd_id: srd006-yaml-machine-config
    symbols_used: [AgentProfile]
  - srd_id: srd013-standard-tool-library
    symbols_used: [ToolDef, ToolDefsFile, ToolSelectionFile]
  - srd_id: srd027-rest-tool-runtime
    symbols_used: [Definition, Collection]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Tool declaration reuse is limited to includes, a path-based textual merge
whose last definition wins. Includes do not name the reusable unit, reject
cross-file duplicates, or report imports that contribute nothing to the
selected program. REST declarations have no reuse mechanism, so profiles
repeat clients, authentication, limits, retry policy, and response mappings.

A profile closure now loads once, decodes strictly, and has a canonical dump.
The loader can therefore compile explicit declaration dependencies with one
set of graph, merge, provenance, and usedness rules. Reuse remains authored
YAML; the repository commits no generated expansion.

<!-- field: goals -->
## Goals

- G1: Let tool and REST declaration files expose named reusable units and import them by explicit relative path, or by an explicit library-root path (srd056).
- G2: Reject ambiguous declaration graphs before registration or execution.
- G3: Preserve declaring-file provenance for diagnostics and for import and OpenAPI source-path resolution without changing existing operational path bases.
- G4: Reject imports that contribute nothing to the resolved profile program.

<!-- req: R1 -->
## R1 - Unit and import syntax

<!-- req: R1.1 -->
### R1.1

An importable tool or REST declaration file SHALL declare `unit` as a non-empty lowercase-kebab name matching `^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`.

<!-- req: R1.2 -->
### R1.2

A declaration file MAY declare `imports` as an ordered list of paths. Each path SHALL be relative to that file's directory, or absolute under a library root, `/opt/<root>/...`, which is either the implicit agent-core root or a root the profile declares (srd056-unit-libraries R1, R2). Any other absolute path SHALL fail, and no import SHALL use a search path, package registry, directory convention, or process working directory. Amended by srd056.

<!-- req: R1.3 -->
### R1.3

Fields `unit` and `imports` SHALL be top-level siblings of the fields named `tools` and `rest`. One declaration file SHALL contain exactly one content kind: `tools` or `rest`, never both. Strict decoding SHALL reject misspelled fields and mixed content kinds.

<!-- req: R1.4 -->
### R1.4

Every file declaring `imports` and every target named by `imports` SHALL declare `unit`. A root file MAY omit `unit` only when it declares no imports.

<!-- req: R1.6 -->
### R1.6

Import expansion SHALL remain runtime loader behavior. No expanded or generated declaration file SHALL be committed.

<!-- req: R2 -->
## R2 - Graph resolution and provenance

<!-- req: R2.1 -->
### R2.1

The loader SHALL resolve imports depth first in authored order and SHALL merge an imported unit before the importing file's local declarations.

<!-- req: R2.2 -->
### R2.2

The loader SHALL canonicalize every imported path and read one canonical source at most once per profile closure.

<!-- req: R2.3 -->
### R2.3

Importing the same canonical source through a diamond SHALL reuse that unit and SHALL NOT report a duplicate.

<!-- req: R2.4 -->
### R2.4

Two different canonical files declaring the same unit name SHALL fail with the unit name and both paths.

<!-- req: R2.5 -->
### R2.5

An import cycle SHALL fail with the ordered canonical path chain from the first repeated ancestor through the repeated path.

<!-- req: R2.6 -->
### R2.6

Import paths and OpenAPI source paths SHALL resolve from the file that declares them. Other operational paths inside ToolDef config or REST values SHALL retain the base defined by their existing config format, such as workspace-relative document roots.

<!-- req: R2.7 -->
### R2.7

Closure file inventory, immutable captured bytes, hashes, and `--dump-config` output SHALL include every imported source exactly once.

<!-- req: R2.8 -->
### R2.8

Canonical source identity SHALL be the absolute cleaned path produced by the existing profile path mapper and filepath resolution; symlink evaluation SHALL not create a second identity rule.

<!-- req: R2.9 -->
### R2.9

Unit traversal SHALL use first-discovery depth-first postorder. A repeated canonical source in a diamond SHALL retain its first position and SHALL not be traversed again.

<!-- req: R3 -->
## R3 - Tool unit merge

<!-- req: R3.1 -->
### R3.1

Tools contributed by imported units SHALL join the profile's tool universe before tools.yaml selection.

<!-- req: R3.2 -->
### R3.2

Two different imported units defining the same tool name SHALL fail before any local override is considered, with the tool name and each unit/path source.

<!-- req: R3.3 -->
### R3.3

Two local definitions of the same tool name in one unit SHALL fail and name the unit/path.

<!-- req: R3.4 -->
### R3.4

Field `override: true` MAY appear only on a tool local to an importing file and SHALL replace exactly one imported definition of the same name in the complete transitive import closure.

<!-- req: R3.5 -->
### R3.5

A local collision without `override: true`, an override with no imported target, and `override: true` on a tool supplied by an imported unit SHALL each fail with the tool and unit/path. A successful override's provenance SHALL become the local unit.

<!-- req: R3.6 -->
### R3.6

Tool order after merge SHALL be deterministic: imported units in first-discovery depth-first postorder, then non-overriding local tools in authored order, with an override occupying its one imported target's position.

<!-- req: R4 -->
## R4 - REST unit merge

<!-- req: R4.1 -->
### R4.1

REST imports SHALL merge clients, servers, OpenAPI declarations, auth profiles, limits, retry policies, response mappings, and document resources before closure validation.

<!-- req: R4.2 -->
### R4.2

A duplicate name in any REST map family SHALL fail; REST declarations SHALL provide no override flag.

<!-- req: R4.3 -->
### R4.3

A duplicate diagnostic SHALL name the map family, entry name, and both unit/path sources.

<!-- req: R4.4 -->
### R4.4

Every imported REST source SHALL apply environment expansion and strict decoding before merge.

<!-- req: R4.5 -->
### R4.5

The loader SHALL strictly decode each source, annotate every raw REST entry with its declaring unit/path, merge raw maps, reject collisions, and then compile the merged closure. An `openapi.path` SHALL resolve from the REST file that declares that OpenAPI entry; a nested external reference SHALL resolve from the OpenAPI document that declares the reference.

<!-- req: R4.6 -->
### R4.6

Merged REST version values SHALL agree. A conflicting version SHALL fail with both unit/path sources.

<!-- req: R4.7 -->
### R4.7

REST merge and diagnostic order SHALL use this fixed family order: clients, servers, openapi, auth, limits, retry_policies, response_mappings, document_resources. Entry names within one family SHALL be sorted.

<!-- req: R5 -->
## R5 - Usedness

<!-- req: R5.1 -->
### R5.1

Usedness SHALL be evaluated after tool selection and REST reference validation over the complete profile closure.

<!-- req: R5.2 -->
### R5.2

An imported tool unit SHALL be used when at least one tool whose provenance is that unit survives tools.yaml selection.

<!-- req: R5.3 -->
### R5.3

REST usedness roots SHALL include every trusted REST reference declared by a selected ToolDef, including client, server, await-source, operation/resource, endpoint, and machine-request bindings.

<!-- req: R5.4 -->
### R5.4

REST usedness nodes SHALL use fully qualified declaration paths such as `(unit, clients, client-name, operations, operation-name)` and `(unit, servers, server-name, endpoints, endpoint-name)`. Traversal SHALL follow containment and every declared `*_ref` edge to auth, limits, retry policies, response mappings, OpenAPI-compiled bindings, operations, resources, endpoints, and document resources.

<!-- req: R5.5 -->
### R5.5

Usedness SHALL apply to import edges. An edge is used when its target unit contributes a reachable declaration or re-exports a unit that does. Usedness SHALL propagate backward through every import edge, so both arms of a diamond that re-export one used leaf are used.

<!-- req: R5.6 -->
### R5.6

The loader SHALL reject every unused import edge in one deterministic diagnostic ordered by importing unit, imported unit, then canonical target path.

<!-- req: R5.7 -->
### R5.7

Root files SHALL not be subject to import usedness because they were not imported units.

<!-- req: R5.8 -->
### R5.8

A successful local override SHALL mark the import edge that supplied its target as used. A binding compiled from an imported OpenAPI entry SHALL retain a usedness edge to that entry's unit/path even when the client or endpoint receiving the binding is local.

<!-- req: R6 -->
## R6 - Import error reporting and reach

<!-- req: R6.3 -->
### R6.3

Import duplicate, cycle, override, and unused errors SHALL name the responsible unit and canonical source path; an error involving two sources SHALL name both.

<!-- req: R6.4 -->
### R6.4

Import errors SHALL occur under `--validate-config`, `--dump-config`, ordinary startup, and profile audit through the same closure loader.

<!-- req: R6.5 -->
### R6.5

Import support SHALL not make directory contents ambient imports. Profile declaration directories retain their existing explicit role as closure roots.

<!-- field: non_goals -->
## Non-goals

- Named data types, tool signatures, typed labels, or machine exhaustiveness.
- Parameterized fragments, conditional imports, version constraints, remote declaration-unit imports, or package registries; existing trusted remote OpenAPI sources remain governed by the REST format.
- Changing tools.yaml selection semantics.
- Generating or committing flattened declarations.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

A two-level tool import and a two-level REST import resolve relative to each declaring file, appear once in the closure inventory, and produce byte-identical repeated dumps.

<!-- field: traces -->
#### Traces

- [R1.1](#r11)
- [R1.2](#r12)
- [R1.3](#r13)
- [R1.4](#r14)
- [R2.1](#r21)
- [R2.2](#r22)
- [R2.3](#r23)
- [R2.6](#r26)
- [R2.7](#r27)
- [R2.8](#r28)
- [R2.9](#r29)

<!-- ac: AC2 -->
### AC2

Tool unit-name collisions, sibling tool-name collisions, duplicate local tools, cycles, and each invalid override form fail with unit and path evidence.

<!-- field: traces -->
#### Traces

- [R2.4](#r24)
- [R2.5](#r25)
- [R3.2](#r32)
- [R3.3](#r33)
- [R3.4](#r34)
- [R3.5](#r35)
- [R3.6](#r36)
- [R6.3](#r63)

<!-- ac: AC3 -->
### AC3

A marked local tool override replaces one imported tool without changing its deterministic position, while the same collision without override fails.

<!-- field: traces -->
#### Traces

- [R2.1](#r21)
- [R3.4](#r34)
- [R3.5](#r35)
- [R3.6](#r36)

<!-- ac: AC4 -->
### AC4

REST imports merge every declared map family, compile an OpenAPI path relative to its source unit, retain one version, and reject duplicate names or conflicting versions with both sources.

<!-- field: traces -->
#### Traces

- [R4.1](#r41)
- [R4.2](#r42)
- [R4.3](#r43)
- [R4.4](#r44)
- [R4.5](#r45)
- [R4.6](#r46)
- [R4.7](#r47)

<!-- ac: AC5 -->
### AC5

Closure validation accepts every transitively used tool and REST import edge, including both arms of a used diamond, and rejects unused edges in deterministic source/target/path order.

<!-- field: traces -->
#### Traces

- [R5.1](#r51)
- [R5.2](#r52)
- [R5.3](#r53)
- [R5.4](#r54)
- [R5.5](#r55)
- [R5.6](#r56)
- [R5.7](#r57)
- [R5.8](#r58)

<!-- ac: AC7 -->
### AC7

The same invalid import closure fails through validation, dump, startup, and profile audit without a second import implementation.

<!-- field: traces -->
#### Traces

- [R6.4](#r64)
- [R6.5](#r65)

<!-- ac: AC8 -->
### AC8

A control-imports profile validates and its canonical dump equals the control profile except for profile identity and the files section.

<!-- field: traces -->
#### Traces

- [R1.6](#r16)
- [R2.7](#r27)
- [R3.1](#r31)
- [R4.1](#r41)
- [R5.1](#r51)

<!-- field: related_documents -->
## Related documents

- cfg-tool-declaration-format
- cfg-rest-tool-format
- cfg-runtime-contract
- rel17.0-uc001-declaration-imports
- srd056-unit-libraries
- test-rel17.0-declaration-imports

<!-- field: implementation_status -->
## Implementation status

<!-- field: implemented -->
### Implemented

- The single closure loader, immutable source inventory, canonical dump, and strict YAML decoder are the compiler foundation.

<!-- field: planned -->
### Planned

- GH-1977 implements tool units and imports.
- GH-1976 implements REST units and imports.
- GH-1978 enforces closure usedness and adds the control-imports proof.
