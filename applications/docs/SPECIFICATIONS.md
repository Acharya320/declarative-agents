---
id: applications-shared-contract-specifications
title: Shared Applications Specification Index
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

This corpus defines the contract for runnable modules under applications/.
srd001 is the baseline for every runnable application. srd002 is conditional:
its requirement groups apply only to managed-service, packaged, Helm-managed,
or kind-demo capabilities. applications/catalog is indexed as a related
foundation but is not a runnable application.

srd003 implements the Release 14.0 consistency target: one
agents/application.yaml composition authority, application-owned actor and UI
grammar, manifest-derived packages, canonical promotion and shared-token
rules, common README and status vocabulary, and explicit catalog/application/
agent-core boundaries. Its capability gates do not require empty optional
directories.

The release 00.0 use case and suite trace the baseline SRDs reciprocally.
Release 14.0 adds a reciprocal cross-application use case and suite.
Release 15.0 adds srd004, the presentation contract and panel composition
model for application UIs, with its use case and suite.
Release 19.0 adds srd008, the persistent local/cloud platform contract:
separate platform, application-compute, and retained-data ownership; one
namespace, collector, and bucket per application; WAL-first immutable
telemetry; explicit purge and reset; and one reusable application-rig
boundary.
Structural
checks parse documents, verify module classification and root orchestration,
enforce ownership-aware stats, and invoke application-owned capability
evidence. The runnable set is chatbot-mesh, coding-agent,
and agent-architecture. All three expose implemented apprig lifecycle,
namespaced serving, retained telemetry, and shared-platform evidence. ENG01
remains the sole kind-rig contract; authenticated live GKE execution is
dependency-gated where project credentials are unavailable.

<!-- field: roadmap_summary | form: "tree" -->
## Roadmap summary

1. - version: "00.0"
   - name: Shared runnable-application contract
   - use_cases_done: 1
   - use_cases_total: 1
   - status: implemented
   - evidence: Root Go tests and application Mage gates are executable.
1. - version: "12.0"
   - name: Agent role-realization alignment
   - use_cases_done: 1
   - use_cases_total: 1
   - status: implemented
   - evidence: Catalog migration, role governance, consumer closure, Helm package, tag-semantic, and scoped zero-hit tests are executable.
1. - version: "13.0"
   - name: Multi-agent application vision
   - use_cases_done: 0
   - use_cases_total: 0
   - status: implemented
   - evidence: The shared vision, root guide, architecture index, and roadmap parse and pass the repository audit.
1. - version: "14.0"
   - name: Application consistency and conformance
   - use_cases_done: 1
   - use_cases_total: 1
   - status: implemented
   - evidence:
     The SRD, cross-application use case, and eight-case test suite are
     reciprocal and executable. Root conformance tests validate all three
     manifests, ownership, package assets, canonical applier promotion, shared
     UI tokens, and README/status classification.
1. - version: "15.0"
   - name: Declarative UX
   - use_cases_done: 1
   - use_cases_total: 1
   - status: implemented
   - evidence: srd004, the use case, and the nine-case suite are reciprocal and executable; the kit, shell, embedded observer, and monorepo migrations landed in GH-2156 through GH-2160.
1. - version: "16.0"
   - name: Chart conformance gate
   - use_cases_done: 1
   - use_cases_total: 1
   - status: implemented
   - evidence:
     ENG01 carries the numbered chart rules, srd005 owns them, and all twelve
     suite cases are executable. The gate renders every chart under every
     checked-in overlay; the baseline is empty, so the charts conform by
     construction rather than by exception.
1. - version: "17.0"
   - name: Pin currency report
   - use_cases_done: 1
   - use_cases_total: 1
   - status: implemented
   - evidence:
     ENG01 carries the pin currency section and the bump verb, srd006 owns
     the two properties, and all nine suite cases are executable. The target
     resolves eleven pins against their registries without a daemon, and its
     own suite reaches no network.
1. - version: "18.0"
   - name: Release annotation gate
   - use_cases_done: 1
   - use_cases_total: 1
   - status: implemented
   - evidence:
     ENG01 carries the release-annotation section and the tag:verify verb,
     srd007 owns the four rules, and all six suite cases are executable. The
     target passes the most recent real-summary tag and refuses the
     degenerate one naming its rule; the runbook carries the step.
1. - version: "19.0"
   - name: Persistent application platform
   - use_cases_done: 1
   - use_cases_total: 1
   - status: implemented
   - evidence:
     srd008, its use case, and the eight-case suite define the ownership,
     routing, lifecycle, storage, query, threat, and reusable-runner contract.
     TC1 through TC8 are implemented by GH-2474: three concurrent application
     migrations, durable collector/query/WAL evidence, explicit purge,
     platform failure separation, external-module apprig use, and cloud IAM
     policy fixtures.

<!-- field: foundation_document_index | form: "tree" -->
## Foundation document index

1. - id: application-composition-pattern-language
   - title: A Pattern Language for Declarative Agent Applications
   - summary:
     Recurring application-composition, package-time closure, deployment
     topology, and module-governance patterns linked to the single-agent
     language and the role-realization semantic model.
   - path: docs/pattern-language.md
1. - id: applications-shared-contract
   - title: Vision for Declarative Multi-Agent Applications
   - summary:
     Multi-agent rationale, intent and assurance flows, role realization,
     human governance, construction sequence, runnable-module contract, and
     capability-aware success criteria.
   - path: docs/VISION.md
1. - id: architecture-applications-shared-contract
   - title: Shared Applications Architecture
   - summary: Ownership layers, module classes, capability layers, and Mage interfaces.
   - path: docs/ARCHITECTURE.md
1. - id: applications-shared-contract-roadmap
   - title: Shared Applications Contract Road Map
   - summary: Release status tied to executable structural evidence.
   - path: docs/road-map.md
1. - id: applications-readme
   - title: Applications
   - summary: Human entry point, terminology, and current classifications.
   - path: README.md
1. - id: eng01-kind-test-demo-rig
   - title: Standard kind Rig for Agent Integration Tests and Demos
   - summary: External engineering contract used only by kind-demo applications.
   - path: ../docs/engineering/eng01-kind-test-demo-rig.md
1. - id: architecture-application-catalog
   - title: Application Catalog Architecture
   - summary: Canonical reusable catalog ownership and consumer boundary.
   - path: catalog/docs/ARCHITECTURE.md
1. - id: architecture-agent-architecture
   - title: Agent Architecture Application Architecture
   - summary: Composition-only ownership, capability status, and local evidence boundary.
   - path: agent-architecture/docs/ARCHITECTURE.md

<!-- field: srd_index | form: "tree" -->
## SRD index

1. - id: srd001-application-module-contract
   - title: Runnable Application Module Contract
   - status: implemented
   - applicability: Every runnable application; excludes applications/catalog.
   - path: docs/specs/software-requirements/srd001-application-module-contract.md
   - use_cases:
     1. rel00.0-uc001-application-contract-conformance
   - test_suites:
     1. test-rel00.0-application-contract-conformance
1. - id: srd002-managed-application-services
   - title: Capability-Dependent Application Services
   - status: implemented
   - applicability: Requirement groups apply only to declared managed_service, packaged, helm_managed, or kind_demo capabilities.
   - path: docs/specs/software-requirements/srd002-managed-application-services.md
   - use_cases:
     1. rel00.0-uc001-application-contract-conformance
   - test_suites:
     1. test-rel00.0-application-contract-conformance
1. - id: srd003-application-consistency
   - title: Application Composition Consistency and Conformance
   - status: implemented
   - applicability:
     Manifest, actor, naming, README, status, and ownership requirements apply
     to runnable or audit-only applications; package, deployment, kind, and UI
     requirements apply only when the matching capability or asset exists.
   - path: docs/specs/software-requirements/srd003-application-consistency.md
   - use_cases:
     1. rel14.0-uc001-application-consistency
   - test_suites:
     1. test-rel14.0-application-consistency
1. - id: srd004-declarative-ux
   - title: Presentation Contract and Panel Composition
   - status: implemented
   - applicability:
     Applications and actors that ship UI assets; R9 applies only to REST
     definitions that select an embedded bundle.
   - path: docs/specs/software-requirements/srd004-declarative-ux.md
   - use_cases:
     1. rel15.0-uc001-declarative-ux
   - test_suites:
     1. test-rel15.0-declarative-ux
1. - id: srd005-chart-conformance
   - title: Chart Conformance Against the Rig Guideline
   - status: implemented
   - applicability:
     R1 through R5, R8, and R11 apply to every application that ships a chart; R6
     applies to every checked-in kind configuration; R7 governs the
     exception baseline.
   - path: docs/specs/software-requirements/srd005-chart-conformance.md
   - use_cases:
     1. rel16.0-uc001-chart-conformance
   - test_suites:
     1. test-rel16.0-chart-conformance
1. - id: srd006-pin-currency
   - title: Pin Integrity and Currency
   - status: implemented
   - applicability:
     Applies to every pinned upstream dependency this repository resolves
     from a registry; Go module versions stay outside it.
   - path: docs/specs/software-requirements/srd006-pin-currency.md
   - use_cases:
     1. rel17.0-uc001-pin-currency
   - test_suites:
     1. test-rel17.0-pin-currency
1. - id: srd007-release-annotation
   - title: Release Annotation Verification
   - status: implemented
   - applicability:
     Applies to every tag on the v0 release line, on demand at the push
     boundary; gates no build and joins no audit.
   - path: docs/specs/software-requirements/srd007-release-annotation.md
   - use_cases:
     1. rel18.0-uc001-release-annotation
   - test_suites:
     1. test-rel18.0-release-annotation
1. - id: srd008-persistent-application-platform
   - title: Persistent Platform and Namespaced Application Contract
   - status: planned
   - applicability:
     Helm-managed applications on the shared platform; cluster-mutating
     integration scenarios retain ENG01's owned-cluster exception.
   - path: docs/specs/software-requirements/srd008-persistent-application-platform.md
   - use_cases:
     1. rel19.0-uc001-persistent-application-platform
   - test_suites:
     1. test-rel19.0-persistent-application-platform

<!-- field: config_format_index | form: "scalar" -->
## Config format index

[]

<!-- field: semantic_model_index | form: "tree" -->
## Semantic model index

1. - id: sm-agent-role-realizations
   - title: Agent Role Realizations and Bindings
   - summary:
     Pinned canonical role and sub-role authority, realization naming grammar,
     actor-role bindings, wrapper inheritance, classifications, and migration
     aliases shared by the catalog and runnable applications.
   - status: implemented
   - path: docs/specs/semantic-models/agent-role-realizations.md
   - reciprocal_indexes:
     1. catalog/docs/SPECIFICATIONS.md#semantic_model_index
     1. catalog/docs/specs/semantic-models/agent-role-taxonomy.md#shared_authority

<!-- field: use_case_index | form: "typed table" -->
## Use case index

| id | title | release | status | srds | test_suite | path |
| --- | --- | --- | --- | --- | --- | --- |
| rel00.0-uc001-application-contract-conformance | Verify a Runnable Application Against the Shared Contract | "00.0" | implemented | [srd001-application-module-contract, srd002-managed-application-services] | test-rel00.0-application-contract-conformance | docs/specs/use-cases/rel00.0-uc001-application-contract-conformance.md |
| rel14.0-uc001-application-consistency | Normalize and Verify a Showcase Application Composition | "14.0" | implemented | [srd003-application-consistency] | test-rel14.0-application-consistency | docs/specs/use-cases/rel14.0-uc001-application-consistency.md |
| rel15.0-uc001-declarative-ux | Compose an Application UI from Declared Panels | "15.0" | implemented | [srd004-declarative-ux] | test-rel15.0-declarative-ux | docs/specs/use-cases/rel15.0-uc001-declarative-ux.md |
| rel16.0-uc001-chart-conformance | Add a Workload and Learn at Audit Time That Its Image Floats | "16.0" | implemented | [srd005-chart-conformance] | test-rel16.0-chart-conformance | docs/specs/use-cases/rel16.0-uc001-chart-conformance.md |
| rel17.0-uc001-pin-currency | Read Which Pins Moved and Which Aged, Then Decide | "17.0" | implemented | [srd006-pin-currency] | test-rel17.0-pin-currency | docs/specs/use-cases/rel17.0-uc001-pin-currency.md |
| rel18.0-uc001-release-annotation | A Degenerate Release Annotation Is Refused Before the Push | "18.0" | implemented | [srd007-release-annotation] | test-rel18.0-release-annotation | docs/specs/use-cases/rel18.0-uc001-release-annotation.md |
| rel19.0-uc001-persistent-application-platform | Operate Two Namespaced Applications Without Losing Retained Telemetry | "19.0" | implemented | [srd008-persistent-application-platform] | test-rel19.0-persistent-application-platform | docs/specs/use-cases/rel19.0-uc001-persistent-application-platform.md |

<!-- field: test_suite_index | form: "typed table" -->
## Test suite index

| id | title | release | status | test_case_count | traces | srds | path |
| --- | --- | --- | --- | --- | --- | --- | --- |
| test-rel00.0-application-contract-conformance | Shared Application Contract Conformance Suite | "00.0" | implemented | 8 | [rel00.0-uc001-application-contract-conformance] | [srd001-application-module-contract, srd002-managed-application-services] | docs/specs/test-suites/test-rel00.0-application-contract-conformance.md |
| test-rel14.0-application-consistency | Release 14.0 Application Consistency Test Suite | "14.0" | implemented | 8 | [rel14.0-uc001-application-consistency] | [srd003-application-consistency] | docs/specs/test-suites/test-rel14.0-application-consistency.md |
| test-rel15.0-declarative-ux | Release 15.0 Declarative UX Test Suite | "15.0" | implemented | 9 | [rel15.0-uc001-declarative-ux] | [srd004-declarative-ux] | docs/specs/test-suites/test-rel15.0-declarative-ux.md |
| test-rel16.0-chart-conformance | Release 16.0 Chart Conformance Test Suite | "16.0" | implemented | 12 | [rel16.0-uc001-chart-conformance] | [srd005-chart-conformance] | docs/specs/test-suites/test-rel16.0-chart-conformance.md |
| test-rel17.0-pin-currency | Release 17.0 Pin Currency Test Suite | "17.0" | implemented | 9 | [rel17.0-uc001-pin-currency] | [srd006-pin-currency] | docs/specs/test-suites/test-rel17.0-pin-currency.md |
| test-rel18.0-release-annotation | Release 18.0 Release Annotation Test Suite | "18.0" | implemented | 6 | [rel18.0-uc001-release-annotation] | [srd007-release-annotation] | docs/specs/test-suites/test-rel18.0-release-annotation.md |
| test-rel19.0-persistent-application-platform | Release 19.0 Persistent Application Platform Test Suite | "19.0" | implemented | 8 | [rel19.0-uc001-persistent-application-platform] | [srd008-persistent-application-platform] | docs/specs/test-suites/test-rel19.0-persistent-application-platform.md |

<!-- field: coverage_gaps | form: "tree" -->
## Coverage gaps

1. - id: release-19-persistent-application-platform
   - status: implemented
   - description:
     GH-2474 implements the contract and all eight cases. Retained objects are
     intentionally not a live query service while the application collector
     is absent; a global query service and authenticated hostile local
     multi-tenancy remain outside this release. The approved GH-2481 decision
     retains application-owned vector stores.
1. - id: release-18-release-annotation
   - status: implemented
   - description:
     The check runs and the runbook carries it between annotation and push.
     The degenerate tags already pushed stay as they are; they are the
     evidence, not the workload. The skill half lives in the coding-skills
     repository as its issue 395 and lands there.
1. - id: release-17-pin-currency
   - status: implemented
   - description:
     The survey runs and reports. On its first run no pin had moved and
     seven were behind, among them three nobody had counted: busybox, chroma
     and the kubectl proxy image. Taking those bumps is separate work, and
     the report exists so the decision is informed rather than overdue.
     Classifying an update from its release notes stays out of this release
     and remains gated in GH-2049.
1. - id: release-16-chart-conformance
   - status: implemented
   - description:
     The gate found fifty violations on its first run, across five upstream
     images and the init containers that stated no pull policy. All fifty
     are closed and the baseline is empty. Rules C4 and C5 found nothing
     when the gate first ran and remain ratchets rather than cleanups.
     Classifying an application template as an un-migrated library candidate
     stays outside this release and is tracked in GH-2410.
1. - id: release-15-declarative-ux
   - status: implemented
   - description:
     TC1 through TC9 are executable. Open kit gaps (same-origin trace backend,
     URL-driven TraceView, span status, shell styles) are tracked in GH-2282.
     The cohere-demo and agentic-wiki-mesh migrations are owned by those
     repositories and consume release ui-kit/v1.0.0.
1. - id: release-14-application-normalization
   - status: implemented
   - description:
     The three showcase applications have validated manifests and normalized
     ownership. Coding Agent and Agent Architecture are composition-only with
     local composition wrappers; Chatbot Mesh is agent-owning while reporting
     wrappers separately. TC1 through TC8 are executable.
1. - id: knowledge-manager-live-managed-service-evidence
   - status: dependency_gated
   - description:
     The demo has implemented deterministic command, configuration, audit,
     test, stats, and root-dispatch evidence. Live lifecycle, readiness,
     liveness, failure, diagnostics, telemetry, shutdown, in-flight, and
     listener-release observations remain required before it satisfies srd002
     AC1 through AC3.
1. - id: future-application-registration
   - status: by_design
   - description:
     The shared suite cannot name future application paths. A future runnable
     module must add itself to root orchestration and add path-specific evidence
     while retaining the generic shared requirement IDs.
1. - id: optional-capability-behavior
   - status: delegated
   - description:
     This suite verifies that capability evidence exists and is executable.
     Application-local SRDs and suites own business behavior, topology, and
     runtime outcomes.
