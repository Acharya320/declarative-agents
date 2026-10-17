---
id: rel14.0-uc001-application-consistency
title: Normalize and Verify a Showcase Application Composition
status: implemented
test_suite: test-rel14.0-application-consistency
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: release -->
## Release

14.0

<!-- field: summary -->
## Summary

A maintainer normalizes one showcase application around
agents/application.yaml, preserves actor and runtime ownership, derives
applicable package and UI evidence, and verifies capability-aware
conformance. The same flow covers Chatbot Mesh, Coding Agent, Agent
Architecture without treating absent optional directories or planned
behavior as implemented.

<!-- field: actor -->
## Actor

A maintainer normalizing or reviewing a showcase application.

<!-- field: trigger -->
## Trigger

An application enters Release 14.0 normalization, changes a direct profile
root or capability, or moves between unregistered, audit-only, and runnable
root orchestration.

<!-- field: preconditions -->
## Preconditions

- The shared srd001, srd002, srd003, pattern language, and role-realization model are available.
- applications/catalog is available as the canonical implementation owner, and applications/ui-kit as the shared UI-token owner (srd004 R8.2).
- agent-core validators are available for referenced profile, machine, tool, ToolDef, and REST documents.
- The maintainer records current evidence before changing a status claim.

<!-- field: candidate_matrix | form: "tree" -->
## Candidate matrix

1. - application: chatbot-mesh
   - current_module_state: runnable
   - current_ownership: agent-owning
   - relevant_surfaces:
     1. managed_service
     1. packaged
     1. helm_managed
     1. kind_demo
     1. ui
   - normalization_focus:
     Add the manifest authority over local actors, the canonical
     corpus-ingest and collector consumers, deployment entries, and both
     actor-served UIs; prove canonical design-token consumption.
1. - application: coding-agent
   - current_module_state: runnable
   - current_ownership: composition-only
   - relevant_surfaces:
     1. managed_service
     1. packaged
     1. helm_managed
     1. kind_demo
   - normalization_focus:
     Retain canonical planner, executor, critic, collector, and applier
     implementations while classifying local serving and applier profiles as
     composition wrappers that contribute no duplicate agents.
1. - application: agent-architecture
   - current_module_state: runnable
   - current_ownership: composition-only
   - relevant_surfaces:
     1. managed_service
     1. packaged
     1. helm_managed
     1. kind_demo
     1. ui
   - normalization_focus:
     Derive curator, collector, applier, deployment, UI, and catalog-doc
     package assets from the manifest and classify the local applier as a
     composition wrapper over the canonical implementation.

<!-- field: flow -->
## Flow

- F1: The maintainer inventories direct profile roots, actor ownership,
  deployment entry profiles, served UI assets, capability claims, and
  current executable evidence. Planned paths and generated package
  inventories are recorded separately from authoritative inputs.
  - traces:
    - srd003-application-consistency R1
    - srd003-application-consistency R6
    - srd003-application-consistency R7
- F2: The maintainer creates or normalizes agents/application.yaml with one
  schema and application identity, module and ownership status, capability
  statuses, direct catalog and local roots, stable runtime paths,
  compatibility, and only applicable deployment and UI entries.
  - traces:
    - srd003-application-consistency R1.1
    - srd003-application-consistency R1.2
    - srd003-application-consistency R1.3
    - srd003-application-consistency R1.4
    - srd003-application-consistency R1.5
    - srd003-application-consistency R1.6
- F3: Local actors and wrappers are placed or mapped under agents/<actor>/.
  Catalog implementations remain canonical references. Fixtures remain
  outside production roots, and no optional directory is created only to
  satisfy the layout.
  - traces:
    - srd003-application-consistency R2.1
    - srd003-application-consistency R2.2
    - srd003-application-consistency R2.3
    - srd003-application-consistency R2.4
- F4: Within each shipped actor surface, the maintainer aligns primary and
  operation-prefixed profile, machine, name-only tools, ToolDef declaration,
  REST, request, and UI names. Absent request or UI variants remain absent.
  - traces:
    - srd003-application-consistency R3.1
    - srd003-application-consistency R3.2
    - srd003-application-consistency R3.3
    - srd003-application-consistency R3.4
    - srd003-application-consistency R3.5
- F5: For a packaged application, packaging starts from manifest roots, traverses
  all reachable profile and declared UI assets, records compatibility and
  actual provenance, and generates the inventory consumed by Helm. For an
  application without packaged, this step is not applicable.
  - traces:
    - srd003-application-consistency R4.1
    - srd003-application-consistency R4.2
    - srd003-application-consistency R4.3
    - srd003-application-consistency R4.4
- F6: The maintainer evaluates duplicated assets for canonical promotion. A
  demonstrated reusable asset migrates with all consumers and compatibility;
  otherwise it remains local. Each UI claiming shared tokens proves an import
  or a generated, drift-checked copy of the canonical token source.
  - traces:
    - srd003-application-consistency R5.1
    - srd003-application-consistency R5.2
    - srd003-application-consistency R5.3
- F7: The README is aligned to the common sections. Manifest, README,
  architecture, roadmap, Mage surfaces, and root lists use the common
  status vocabulary and identify missing or dependency-gated evidence for
  each capability independently.
  - traces:
    - srd003-application-consistency R6.1
    - srd003-application-consistency R6.2
    - srd003-application-consistency R6.3
    - srd003-application-consistency R6.4
- F8: Structural tests parse and trace the corpus, validate the manifest and
  naming grammar, compare root orchestration and status, and verify ownership.
  Applicable executable package, Helm, kind, UI, and application tests run.
  Every absent capability is reported not_applicable rather than passed.
  - traces:
    - srd003-application-consistency R7.1
    - srd003-application-consistency R7.2
    - srd003-application-consistency R7.3

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | applications/<application>/agents/application.yaml | Single application-owned composition and capability authority. |
| T2 | applications/<application>/agents/<actor>/ | Local actor, wrapper, request, and actor-served UI ownership. |
| T3 | applications/catalog/agents/ | Canonical reusable implementation ownership and compatibility source. |
| T4 | applications/ui-kit/src/tokens.css | Canonical shared UI-token source when a UI declares consumption. |
| T5 | applications/<application>/magefiles/ | Manifest validation, package derivation, local audit, and capability evidence. |
| T6 | magefiles/build.go | Runnable versus audit-only root membership. |
| T7 | agent-core | Runtime validation and execution semantics remain external authority. |
| T8 | srd003-application-consistency AC1 | Manifest authority and capability gating. |
| T9 | srd003-application-consistency AC2 | Actor, catalog, fixture, and optional-directory grammar. |
| T10 | srd003-application-consistency AC3 | Program and UI naming consistency. |
| T11 | srd003-application-consistency AC4 | Manifest-derived package closure and provenance. |
| T12 | srd003-application-consistency AC5 | Canonical promotion and actual shared-token consumption. |
| T13 | srd003-application-consistency AC6 | Common README, statuses, and root classification. |
| T14 | srd003-application-consistency AC7 | Catalog, application, and runtime ownership closure. |
| T15 | test-rel14.0-application-consistency | Structural and capability-gated executable evidence for this flow. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The candidate has one manifest authority whose direct roots, runtime
paths, compatibility, deployment, UI, capabilities, and status are
internally consistent.

<!-- field: traces -->
#### Traces

- srd003-application-consistency R1
- srd003-application-consistency AC1

<!-- sc: S2 -->
### S2

Local actors and variants follow the normalized grammar, canonical
implementations remain catalog-owned, and absent optional surfaces do not
require placeholders.

<!-- field: traces -->
#### Traces

- srd003-application-consistency R2
- srd003-application-consistency R3
- srd003-application-consistency AC2
- srd003-application-consistency AC3

<!-- sc: S3 -->
### S3

Every applicable package and Helm profile root and provenance record
derives from the manifest, while non-packaged applications skip the group
as not applicable.

<!-- field: traces -->
#### Traces

- srd003-application-consistency R4
- srd003-application-consistency AC4

<!-- sc: S4 -->
### S4

Reusable promotion has consumer evidence and one owner, and every
shared-token claim proves canonical consumption.

<!-- field: traces -->
#### Traces

- srd003-application-consistency R5
- srd003-application-consistency AC5

<!-- sc: S5 -->
### S5

README, manifest, local documents, executable evidence, and root
orchestration agree on module and per-capability status.

<!-- field: traces -->
#### Traces

- srd003-application-consistency R6
- srd003-application-consistency AC6

<!-- sc: S6 -->
### S6

The application owns composition and evidence, the catalog owns reusable
behavior, and agent-core owns runtime semantics with no duplicate
implementation or redefined format.

<!-- field: traces -->
#### Traces

- srd003-application-consistency R7
- srd003-application-consistency AC7

<!-- sc: S7 -->
### S7

All three applications have truthful normalization findings. Chatbot Mesh
reports only its local agents, while Coding Agent and Agent Architecture
remain composition-only consumers of canonical implementations.

<!-- field: traces -->
#### Traces

- srd003-application-consistency R1.2
- srd003-application-consistency R2.2

<!-- field: out_of_scope -->
## Out of scope

- Rename or move application assets in this documentation release.
- Require an optional directory or capability because another application has it.
- Change canonical catalog behavior or agent-core runtime formats.
- Mark Release 14.0 normalization complete before named executable evidence exists.
