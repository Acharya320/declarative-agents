---
id: srd001-application-module-contract
title: Runnable Application Module Contract
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

A runnable module needs a common ownership, documentation, orchestration, and
accounting contract. Without one, an application can copy catalog behavior,
disappear from root gates, or claim conformance from prose rather than
executable evidence.

<!-- field: goals -->
## Goals

- G1: Define the baseline contract for every runnable application module.
- G2: Preserve canonical catalog ownership and application-owned composition.
- G3: Make documentation, root Mage participation, release participation, and agent accounting structurally observable.
- G4: Let application-local documents extend the shared contract by reference.

<!-- field: applicability -->
## Applicability

<!-- field: capability -->
### Capability

runnable_module

<!-- field: applies_to -->
### Applies to

Every runnable application under applications/.

<!-- field: excludes -->
### Excludes

- applications/catalog, because it is a reusable catalog module.
- Directories that contain only shared documentation or fixtures and expose no runnable application composition.

<!-- req: R1 -->
## R1 - Module identity and ownership

<!-- field: applicability -->
### Applicability

runnable_module

<!-- req: R1.1 -->
### R1.1

A runnable application shall own an executable composition under one
applications/<application>/ root.

<!-- field: evidence -->
#### Evidence

The application root contains the composition entry point named by
its README or architecture and an application-local execution or
validation target.

<!-- req: R1.2 -->
### R1.2

The application shall classify itself as agent-owning or
composition-only. An agent-owning application owns only behavior that
is specific to that application. A composition-only application
contributes no canonical agent implementation.

<!-- field: evidence -->
#### Evidence

Application architecture names ownership. Its stats output includes
local agents only when it owns them and reports canonical references
separately.

<!-- req: R1.3 -->
### R1.3

applications/catalog shall remain a reusable catalog module and shall
not appear in the runnable application module list.

<!-- field: evidence -->
#### Evidence

magefiles/build.go lists applications/catalog in subModules and omits
it from applicationModules and auditOnlyApplicationModules.

<!-- req: R2 -->
## R2 - Canonical references and composition

<!-- field: applicability -->
### Applicability

runnable_module when reusable catalog behavior is consumed

<!-- req: R2.1 -->
### R2.1

An application shall reference a catalog member from its canonical
applications/catalog path. It shall not copy the member's machine,
declarations, reusable behavior, or family SRD.

<!-- field: evidence -->
#### Evidence

Composition manifests and architecture name canonical source paths;
structural closure or ownership tests reject copied canonical assets.

<!-- req: R2.2 -->
### R2.2

The application shall own topology, sequencing between members,
application-only wrappers, endpoint choices, configuration,
packaging, deployment, operator UX, and end-to-end tests that are
specific to the composition.

<!-- field: evidence -->
#### Evidence

Application architecture maps each app-owned artifact to a path
beneath the application root.

<!-- req: R2.3 -->
### R2.3

Packaging that resolves catalog references shall record compatible
release and actual source provenance separately and shall assemble a
closed runtime tree before execution.

<!-- field: evidence -->
#### Evidence

A packaging manifest and deterministic package tests show compatible
release, source revision, dirty state when applicable, and closure.

<!-- field: applies_when -->
#### Applies when

The application declares the packaged capability.

<!-- req: R3 -->
## R3 - Application documentation

<!-- field: applicability -->
### Applicability

runnable_module

<!-- req: R3.1 -->
### R3.1

The application shall own docs/VISION.md, docs/ARCHITECTURE.md,
docs/road-map.md, and docs/SPECIFICATIONS.md.

<!-- field: evidence -->
#### Evidence

All four paths exist, parse as YAML, and contain the required fields
from agent-core/docs/constitutions/design.yaml.

<!-- req: R3.2 -->
### R3.2

The application shall own SRDs, use cases, and test suites for its
application-specific behavior. SPECIFICATIONS shall index each file
and preserve reciprocal use-case and test-suite traces.

<!-- field: evidence -->
#### Evidence

Indexed paths resolve; each use case names or is named by its suite;
each suite traces its use case; tests trace SRD acceptance criteria.

<!-- req: R3.3 -->
### R3.3

Application-local architecture and SRDs shall cite this SRD and
applicable requirement groups from srd002. They shall define only
application behavior, stricter constraints, or local evidence and
shall not copy shared obligations.

<!-- field: evidence -->
#### Evidence

Related-document or dependency entries name shared IDs, while local
requirements name app-owned behavior and paths.

<!-- req: R4 -->
## R4 - Root Mage participation

<!-- field: applicability -->
### Applicability

runnable_module

<!-- req: R4.1 -->
### R4.1

Every runnable application shall expose a local mage audit target and
participate in root mage audit through applicationModules or
auditOnlyApplicationModules.

<!-- field: evidence -->
#### Evidence

The module has magefiles/magefile.go or magefile.go with Audit, is
present in the applicable root list, and root mage audit dispatches it.

<!-- req: R4.2 -->
### R4.2

An application with Go tests shall participate in root mage test.
An application with no Go tests shall not claim Go-test evidence.

<!-- field: evidence -->
#### Evidence

Root test dispatch includes applicationModules with discoverable Go
tests, and the application's test suite names existing tests.

<!-- req: R4.3 -->
### R4.3

A runnable application that reports repository metrics shall expose
mage stats and participate in root mage stats.

<!-- field: evidence -->
#### Evidence

statsParticipants contains the application and its local stats output
parses as the root aggregator's JSON shape.

<!-- req: R4.4 -->
### R4.4

A release-relevant runnable application shall participate in the root
release gate and module tag set. A documentation-only transitional
application may remain audit-only until executable surfaces exist.

<!-- field: evidence -->
#### Evidence

magefiles/tag.go derives release tags from releaseModules and release
gates execute root audit and test; the module's declared release state
matches its orchestration list.

<!-- req: R5 -->
## R5 - Zero duplicate agent accounting

<!-- field: applicability -->
### Applicability

runnable_module

<!-- req: R5.1 -->
### R5.1

Root agent totals shall count each canonical agent implementation
exactly once.

<!-- field: evidence -->
#### Evidence

Root stats aggregates only agents sections and validates module totals.

<!-- req: R5.2 -->
### R5.2

A composition-only application shall report canonical references and
serving composition outside an agents section.

<!-- field: evidence -->
#### Evidence

Application stats JSON has application or composition data and no
agents key; unit tests pin the output.

<!-- req: R5.3 -->
### R5.3

An agent-owning application shall count only locally owned agents and
shall report catalog wrappers or references separately.

<!-- field: evidence -->
#### Evidence

Ownership scans distinguish local implementations from composition;
root totals ignore the composition section.

<!-- req: R6 -->
## R6 - Baseline Mage surface

<!-- field: applicability -->
### Applicability

runnable_module

<!-- req: R6.1 -->
### R6.1

The baseline application Mage surface shall include audit. test and
stats shall exist when the module has corresponding executable tests
and metrics participation.

<!-- field: evidence -->
#### Evidence

mage -l in the application root lists the claimed targets, and root
orchestration invokes them through the matching participant lists.

<!-- req: R6.2 -->
### R6.2

Build, package, image, Helm, integration, doctor, and demo targets
shall not be baseline requirements. They apply only through the
capabilities in srd002 or stricter application-local requirements.

<!-- field: evidence -->
#### Evidence

The application capability declaration and mage -l agree; conformance
does not fail for an undeclared optional surface.

<!-- field: non_goals -->
## Non-goals

- Define managed-service lifecycle, packaging, Helm, or kind-demo behavior; srd002 owns those capability-dependent requirements.
- Define catalog member behavior or agent-core runtime behavior.
- Require an application to contribute a new agent implementation.
- Require Build or the root default target to dispatch runnable applications.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

Every runnable application has an app-owned executable composition and
explicit agent-owning or composition-only classification; the catalog
remains outside the runnable list.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3

<!-- field: applicability -->
#### Applicability

runnable_module

<!-- ac: AC2 -->
### AC2

Catalog consumers use canonical references while the application owns
composition and, when packaged, records a deterministic closed package
with truthful compatibility and source provenance.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3

<!-- field: applicability -->
#### Applicability

runnable_module; R2.3 only when packaged

<!-- ac: AC3 -->
### AC3

Required local documents parse, index reciprocal traces, and extend the
shared contract by reference.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3

<!-- field: applicability -->
#### Applicability

runnable_module

<!-- ac: AC4 -->
### AC4

Root audit always dispatches the application, root test and stats dispatch
it when applicable, and release-relevant modules enter root release gates
and tags.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
- R4.3
- R4.4

<!-- field: applicability -->
#### Applicability

runnable_module with each conditional surface as stated

<!-- ac: AC5 -->
### AC5

Stats count canonical implementations once and report application
composition separately from locally owned agents.

<!-- field: traces -->
#### Traces

- R5.1
- R5.2
- R5.3

<!-- field: applicability -->
#### Applicability

runnable_module

<!-- ac: AC6 -->
### AC6

audit is the only universal local Mage target; test, stats, build,
package, image, Helm, integration, doctor, and demo targets are required
only by actual module content or declared capabilities.

<!-- field: traces -->
#### Traces

- R6.1
- R6.2

<!-- field: applicability -->
#### Applicability

runnable_module

<!-- field: use_case_traces -->
## Use case traces

- rel00.0-uc001-application-contract-conformance

<!-- field: test_suite_traces -->
## Test suite traces

- test-rel00.0-application-contract-conformance
