---
id: applications-shared-contract-roadmap
title: Shared Applications Contract Road Map
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- release: 00.0 | status: implemented -->
## 00.0 - Shared runnable-application contract

Define the reusable application documentation layer, separate baseline
module obligations from managed-service, package, Helm, and kind-demo
capabilities, and connect the contract to current executable structural
evidence.

<!-- field: srds -->
### SRDs

- srd001-application-module-contract
- srd002-managed-application-services

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| [rel00.0-uc001-application-contract-conformance](specs/use-cases/rel00.0-uc001-application-contract-conformance.md) | Verify baseline and declared capability conformance through structural and application-owned evidence. | implemented |

<!-- field: applicability -->
### Applicability

<!-- field: baseline -->
#### Baseline

Every runnable application.

<!-- field: managed_service -->
#### Managed service

Only applications with application-owned long-running services.

<!-- field: packaged -->
#### Packaged

Only applications that produce distributable artifacts.

<!-- field: helm_managed -->
#### Helm managed

Only applications that ship a Helm chart.

<!-- field: kind_demo -->
#### Kind demo

Only applications that ship a kind integration or demo; ENG01 applies.

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| [test-rel00.0-application-contract-conformance](specs/test-suites/test-rel00.0-application-contract-conformance.md) | Parse shared YAML, verify reciprocal traces, run root orchestration/accounting tests, and invoke current application evidence. | implemented |

<!-- field: current_applications | form: "typed table" -->
### Current applications

| path | ownership | runnable_module | managed_service | live_managed_service_evidence | packaged | helm_managed | kind_demo |
| --- | --- | --- | --- | --- | --- | --- | --- |
| applications/chatbot-mesh | agent-owning | implemented |  |  |  |  |  |
| applications/coding-agent | composition-only | implemented |  |  |  |  |  |
| applications/agent-architecture | composition-only | implemented | partial | dependency_gated | planned | planned | planned |

<!-- field: implementation_evidence | form: "tree" -->
### Implementation evidence

1. - path: magefiles/example_modules_test.go
   - proves: Runnable application paths and root audit, test, and stats participation are executable checks.
1. - path: magefiles/stats_test.go
   - proves: Root agent totals sum canonical agents sections and ignore composition-only sections.
1. - path: applications/coding-agent/magefiles/stats_test.go
   - proves: A composition-only application reports references without contributing duplicate agents.
1. - path: applications/chatbot-mesh/magefiles/agent_stats_test.go
   - proves: An agent-owning application separates local implementations from canonical wrappers.
1. - path: applications/chatbot-mesh/magefiles/magefile.go
   - proves: The chatbot-mesh application owns an executable audit gate.
1. - path: applications/coding-agent/magefiles/magefile.go
   - proves: The coding-agent application owns an executable audit gate.
1. - path: applications/agent-architecture/magefiles/magefile.go
   - proves: The demo owns executable audit and composition-only stats targets.
1. - path: applications/agent-architecture/magefiles/magefile_test.go
   - proves:
     The demo's runnable composition, canonical reference, reciprocal
     traces, and zero-agent accounting have deterministic tests; its live
     managed-service observations remain dependency-gated.
1. - path: magefiles/build.go
   - proves: Root audit, test, stats, and release lists register all three runnable applications.
1. - path: magefiles/stats_test.go
   - proves:
     Root accounting accepts the demo's application section without an
     agents section and retains 15 canonical agents.
1. - path: docs/engineering/eng01-kind-test-demo-rig.md
   - proves: kind-demo obligations have one referenced engineering contract.

<!-- field: non_runnable_modules -->
### Non runnable modules

| path | classification |
| --- | --- |
| applications/catalog | reusable catalog |

<!-- release: 12.0 | status: implemented -->
## 12.0 - Agent role-realization alignment

Align catalog and application actor names with the pinned Agent View role
authority while preserving canonical implementation ownership. The shared
semantic model lands first; executable binding governance and profile
migrations follow without rewriting historical releases.

<!-- field: dependencies | form: "tree" -->
### Dependencies

1. - id: catalog-release-07.1
   - status: in-progress
   - reason: Preserve and update incomplete per-family boundary evidence during migration.
1. - id: catalog-release-11.0
   - status: in-progress
   - reason: Preserve assembler-rig evidence while resolving its canonical-name collision.
1. - id: gh-1076-semantic-model
   - status: implemented
   - reason: Defines pinned role authority, classifications, inheritance, and aliases.
1. - id: gh-1074-binding-validation
   - status: implemented
   - reason:
     Deterministic root-audit and catalog-conformance tests validate profile
     inventory, classifications, role/sub-role bindings, inheritance,
     aliases, repeated bindings, and canonical-name collisions.

<!-- field: use_cases | form: "tree" -->
### Use cases

1. - id: rel12.0-uc001-agent-role-realization-alignment
   - summary:
     Classify each shipped profile, validate many-to-many role bindings and
     wrapper inheritance, and migrate incompatible canonical-name collisions
     through versioned aliases.
   - status: implemented
   - semantic_models:
     1. sm-agent-role-realizations

<!-- field: implementation_sequence -->
### Implementation sequence

- Pin the shared semantic model and reciprocal indexes.
- Add executable profile discovery, binding, parent, inheritance, and collision validation.
- Migrate Critic realizations and catalog canonical-name collisions.
- Coordinate release aliases, consumers, packaging closures, and migration evidence.

<!-- field: release -->
### Release

<!-- field: canonical_catalog -->
#### Canonical catalog

applications/catalog/v0.20260727.0

<!-- field: legacy_alias -->
#### Legacy alias

agent-profiles/v0.20260727.0

<!-- field: publication -->
#### Publication

post-merge from clean main through mage tag

<!-- field: implementation_evidence | form: "tree" -->
### Implementation evidence

1. - path: magefiles/agent_role_realizations_test.go
   - proves:
     All repository application profiles are discovered deterministically
     and checked against the shared role, classification, binding,
     inheritance, exclusion, path, and alias authority.
1. - path: applications/catalog/conformance/family_inventory_test.go
   - proves: Catalog family inventory consumes shared actor classifications instead of treating every family as a canonical role.
1. - path: applications/catalog/conformance/membership_contract_test.go
   - proves: Catalog test-time members and incompatible historical aliases retain their documented realization semantics.
1. - path: magefiles/audit.go
   - proves: Root mage audit executes the actor-role realization validator before module audits.
1. - path: applications/catalog/docs/migrations/v0.20260727.0-agent-role-realization-alignment.md
   - proves: Every rename has an old-to-new path or identity entry, compatibility mechanism and duration, consumer pin, and post-merge publication boundary.
1. - path: applications/catalog/conformance/release_alignment_test.go
   - proves: Exact release metadata, consumer pins, migration coverage, and scoped zero-hit policy are executable.
1. - path: applications/coding-agent/magefiles/profiles_closure_test.go
   - proves: The coding application pins the alignment release and packages every canonical transitive profile asset.
1. - path: applications/chatbot-mesh/magefiles/helm_package_test.go
   - proves: The chart records the catalog pin and its generated archive contains canonical renamed paths without stale coordinator artifacts.
1. - path: magefiles/tag_test.go
   - proves: The root release creates the exact canonical catalog and legacy compatibility pair atomically without feature-branch publication.

<!-- release: 13.0 | status: implemented -->
## 13.0 - Multi-agent application vision

Explain when multiple agents are warranted and how a scenario becomes a
governed declarative-agent application. Join intent-to-actuation and
observation-to-feedback flows, role realizations, human governance,
canonical reuse, package closure, and executable evidence to the retained
shared capability contract.

<!-- field: dependencies -->
### Dependencies

| id | status | reason |
| --- | --- | --- |
| gh-1076-semantic-model | implemented | Defines canonical roles as semantic interfaces and records actor/profile realizations. |
| application-composition-pattern-language | implemented | Supplies the construction patterns referenced by the root guide. |

<!-- field: implementation_evidence | form: "tree" -->
### Implementation evidence

1. - path: applications/docs/VISION.md
   - proves:
     The canonical vision records the multi-agent thesis, both operating
     flows, human boundary, six-step construction sequence, retained
     capability contract, and pinned conceptual-source provenance.
1. - path: applications/README.md
   - proves:
     A root reader can determine why and when to use several actors and
     follow the ordered scenario-to-evidence construction guide.
1. - path: applications/docs/ARCHITECTURE.md
   - proves: The project index assigns the revised VISION and README responsibilities and records implementation status.
1. - path: applications/docs/SPECIFICATIONS.md
   - proves: The foundation index names and summarizes the broadened canonical vision.
1. - path: applications/docs/specs/semantic-models/agent-role-realizations.md
   - proves: Canonical roles, actor/profile realizations, inheritance, and ownership boundaries remain machine-readable.
1. - path: applications/docs/pattern-language.md
   - proves: Canonical reuse, composition, package closure, profile-free runtime, wrappers, and workload boundaries have construction patterns.

<!-- release: 14.0 | status: implemented -->
## 14.0 - Application consistency and conformance

Implement one capability-aware composition, actor, package, UI, README, and
status contract for Chatbot Mesh, Coding Agent, and Agent Architecture.
agents/application.yaml becomes the composition authority.
Packaged closures and provenance derive from it, while absent optional
capabilities require no empty directories.

<!-- field: parent -->
### Parent

gh-1414-normalize-showcase-applications

<!-- field: dependencies -->
### Dependencies

| id | status | reason |
| --- | --- | --- |
| srd001-application-module-contract | implemented | Supplies baseline runnable and audit-only module governance. |
| srd002-managed-application-services | implemented | Supplies managed-service, package, Helm, and kind capability gates. |
| sm-agent-role-realizations | implemented | Owns actor classification, role bindings, inheritance, and naming grammar. |
| application-composition-pattern-language | implemented | Defines canonical ownership, composition manifest, closure, wrapper, and workload patterns. |
| design-constitution | implemented | Defines structured document fields and reciprocal traceability. |

<!-- field: srds -->
### SRDs

- srd003-application-consistency

<!-- field: use_cases | form: "tree" -->
### Use cases

1. - id: rel14.0-uc001-application-consistency
   - summary:
     Normalize and verify one showcase application while preserving
     capability gating, ownership, and truthful audit-only or runnable status.
   - status: implemented

<!-- field: implementation_sequence | form: "typed table" -->
### Implementation sequence

| order | action | status |
| --- | --- | --- |
| 1 | Specify the shared SRD, use case, suite, architecture, roadmap, index, and root guide. | implemented |
| 2 | Add shared manifest and application-layout validation without requiring optional surfaces. | implemented |
| 3 | Normalize Chatbot Mesh actor roots, manifest, deployment entries, UIs, and shared-token consumption. | implemented |
| 4 | Normalize Coding Agent manifest, actor variants, ownership, and manifest-derived packaging. | implemented |
| 5 | Normalize Agent Architecture manifest, applier ownership, package/UI roots, and stale capability status. | implemented |
| 6 | Run the full Release 14.0 suite and close only evidence-backed capability states. | implemented |

<!-- field: test_suites | form: "tree" -->
### Test suites

1. - id: test-rel14.0-application-consistency
   - summary:
     Parse and trace shared documents, then validate manifests, actor
     grammar, package derivation, UI tokens, README/status agreement, the
     three-application matrix, and root audit.
   - status: implemented
   - implemented_cases:
     1. TC1
     1. TC2
     1. TC3
     1. TC4
     1. TC5
     1. TC6
     1. TC7
     1. TC8
   - planned_cases: []

<!-- field: current_application_evidence -->
### Current application evidence

<!-- current_application_evidence_entry: applications/chatbot-mesh -->
#### applications/chatbot-mesh

<!-- field: module_state -->
##### Module state

runnable

<!-- field: ownership -->
##### Ownership

agent-owning

<!-- field: release_14_status -->
##### Release 14 status

implemented

<!-- field: note -->
##### Note

The manifest, package closure, local actors, two composition wrappers,
deployment entries, and canonical token imports have executable checks.

<!-- current_application_evidence_entry: applications/coding-agent -->
#### applications/coding-agent

<!-- field: module_state -->
##### Module state

runnable

<!-- field: ownership -->
##### Ownership

composition-only

<!-- field: release_14_status -->
##### Release 14 status

implemented

<!-- field: note -->
##### Note

Local planner, executor, critic, and applier profiles are composition
wrappers; canonical implementations remain catalog-owned and stats
contribute zero agents.

<!-- current_application_evidence_entry: applications/agent-architecture -->
#### applications/agent-architecture

<!-- field: module_state -->
##### Module state

runnable

<!-- field: ownership -->
##### Ownership

composition-only

<!-- field: release_14_status -->
##### Release 14 status

implemented

<!-- field: note -->
##### Note

Canonical curator, collector, lifecycle-exit, and applier behavior
remains catalog-owned. The local applier is one composition wrapper,
and package/UI assets derive from the manifest.

<!-- field: evidence_status | form: "tree" -->
### Evidence status

- shared_documents: implemented
- reciprocal_traceability: implemented
- yaml_parse_case: implemented
- current_root_audit: implemented
- cross_application_manifest_validation: implemented
- application_normalization: implemented

<!-- field: implementation_evidence | form: "typed table" -->
### Implementation evidence

| path | command | proves |
| --- | --- | --- |
| applications/docs/specs/software-requirements/srd003-application-consistency.md |  | Numbered capability-aware composition consistency requirements and acceptance criteria. |
| applications/docs/specs/use-cases/rel14.0-uc001-application-consistency.md |  | Three-application flow and ownership touchpoints. |
| applications/docs/specs/test-suites/test-rel14.0-application-consistency.md |  | Requirement-to-AC-to-flow-to-test coverage and eight executable cases. |
| applications/docs/ARCHITECTURE.md |  | Normalized authority, ownership components, decisions, target structure, and status. |
| applications/docs/SPECIFICATIONS.md |  | Reciprocal SRD, use-case, suite, and coverage-gap indexes. |
| applications/README.md |  | Capability-gated directory guidance and truthful current classifications. |
|  | mage audit | Repository documentation, ownership, application audit, and root classification gates remain green. |
|  | cd magefiles && go test -run 'Test(ApplicationConsistency\|ApplicationActor\|ManifestDerived\|ApplicationPromotion\|ApplicationReadme\|Release14)' | Cross-application Release 14 manifests, grammar, package assets, promotion, tokens, statuses, and matrix remain executable. |

<!-- release: 15.0 | status: implemented -->
## 15.0 - Declarative UX

Pin the presentation contract shared panels read, make the panel the reuse
unit, compose application bundles from ui.yaml at build time, move design
tokens and observability panels into applications/ui-kit, and serve the
observer UI from a bundle embedded in the REST tool.

<!-- field: parent -->
### Parent

gh-2154-declarative-ux

<!-- field: dependencies -->
### Dependencies

| id | status | reason |
| --- | --- | --- |
| srd003-application-consistency | implemented | Owns the application manifest, UI assets, and shared-token policy this release amends. |
| catalog-srd021-observer | implemented | Owns the /monitor/fleet response the embedded observer reads (GH-2170). |

<!-- field: srds -->
### SRDs

- srd004-declarative-ux

<!-- field: use_cases | form: "tree" -->
### Use cases

1. - id: rel15.0-uc001-declarative-ux
   - summary:
     Compose an application UI from ui.yaml, kit panels, and domain panels,
     and serve the observer UI from an embedded bundle.
   - status: implemented

<!-- field: implementation_sequence | form: "typed table" -->
### Implementation sequence

| order | action | status |
| --- | --- | --- |
| 1 | Specify srd004, the use case, the suite, and the index entries (GH-2155). | implemented |
| 2 | Scaffold the kit, move tokens, and add the client, fixtures, shell primitives, and hooks (GH-2156; GH-2259 to GH-2262). | implemented |
| 3 | Extract the kit panels as a best-of merge (GH-2157; GH-2263 to GH-2266). | implemented |
| 4 | Serve the observer from an embedded bundle and delete its UI tree (GH-2158; GH-2267 to GH-2269). | implemented |
| 5 | Validate ui.yaml version 2 and run the chatbot-mesh chatbot on the generic shell (GH-2159; GH-2270 to GH-2272). | implemented |
| 6 | Migrate the remaining monorepo UIs, add the duplication sweep, and version the kit 1.0.0 (GH-2160; GH-2273 to GH-2275). | implemented |

<!-- field: test_suites | form: "tree" -->
### Test suites

1. - id: test-rel15.0-declarative-ux
   - summary:
     Trace the release documents, then verify the kit client, fixtures,
     panels, shell, ui.yaml version 2, token ownership, embedded bundle,
     and the absence of duplicated UI code.
   - status: implemented
   - implemented_cases:
     1. TC1
     1. TC2
     1. TC3
     1. TC4
     1. TC5
     1. TC6
     1. TC7
     1. TC8
     1. TC9
   - planned_cases: []

<!-- field: implementation_evidence -->
### Implementation evidence

| path | proves |
| --- | --- |
| applications/docs/specs/software-requirements/srd004-declarative-ux.md | Numbered contract, proxy, compatibility, panel, composition, client, ui.yaml, kit, and embedded-bundle requirements. |
| applications/docs/specs/use-cases/rel15.0-uc001-declarative-ux.md | Composition flow and the endpoint-usage review of all three repositories. |
| applications/docs/specs/test-suites/test-rel15.0-declarative-ux.md | Requirement-to-AC-to-test coverage with TC1 executable. |

<!-- release: 16.0 | status: implemented -->
## 16.0 - Chart conformance gate

Turn the ENG01 chart rules from prose a reviewer holds in their head
into numbered requirements a gate cites. Render every application
chart under every checked-in values overlay, check the rendered
manifests, and give an accepted violation one checked-in place that
shrinks.

<!-- field: parent -->
### Parent

gh-2047-chart-conformance

<!-- field: dependencies -->
### Dependencies

| id | status | reason |
| --- | --- | --- |
| srd003-application-consistency | implemented | Owns the application manifest and the helm_managed capability gate this release reads. |
| eng01-kind-test-demo-rig | implemented | States the rules this release numbers and gates; GH-2407 adds the rule table to it. |

<!-- field: srds -->
### SRDs

- srd005-chart-conformance

<!-- field: use_cases | form: "tree" -->
### Use cases

1. - id: rel16.0-uc001-chart-conformance
   - summary:
     An engineer adds a workload whose image floats, the audit fails
     naming the chart, the overlay, the resource and the rule, and the
     exception baseline carries what must stand.
   - status: implemented

<!-- field: implementation_sequence | form: "typed table" -->
### Implementation sequence

| order | action | status |
| --- | --- | --- |
| 1 | Add the ENG01 rule table and specify srd005, the use case, the suite, and the index entries (GH-2407). | planned |
| 2 | Add the chart conformance module, its exception baseline, and the root audit wiring (GH-2408). | implemented |
| 3 | Pin the floating image tags and shrink the baseline (GH-2409). | implemented |

<!-- field: test_suites | form: "tree" -->
### Test suites

1. - id: test-rel16.0-chart-conformance
   - summary:
     Trace the release documents, then verify render coverage, each
     numbered rule, the baseline's three behaviours, the shape of a
     finding, and the pinning of the seeded violations.
   - status: planned
   - implemented_cases:
     1. TC1
     1. TC2
     1. TC3
     1. TC4
     1. TC5
     1. TC6
     1. TC7
     1. TC8
     1. TC9
     1. TC10
   - planned_cases: []

<!-- field: implementation_evidence | form: "typed table" -->
### Implementation evidence

| path | command | proves |
| --- | --- | --- |
| docs/engineering/eng01-kind-test-demo-rig.md |  | Numbered chart conformance rules C1 through C6 and the exception-baseline contract. |
| applications/docs/specs/software-requirements/srd005-chart-conformance.md |  | Numbered rendering, pinning, policy, service, readiness, baseline, and reporting requirements. |
| applications/docs/specs/use-cases/rel16.0-uc001-chart-conformance.md |  | The flow from a changed chart to a named finding to a shrinking baseline. |
| applications/docs/specs/test-suites/test-rel16.0-chart-conformance.md |  | Requirement-to-AC-to-test coverage with every case executable. |
| applications/charts/agent-services/templates/_naming.tpl |  | One helper renders every image reference and appends the digest the values pin. |
|  | cd applications/chatbot-mesh && mage integration:helmSmoke | The mesh installs and serves on the pinned images against a real cluster, Dolt metrics included. |
| magefiles/chartconf/ |  | Every numbered rule, the shrink-only baseline, and the shape of a finding, checked from rendered YAML without helm. |
| magefiles/chart_conformance_test.go |  | Every chart rendered under its defaults and each checked-in overlay, every kind configuration read, and the baseline reconciled. |
| magefiles/chartconf/baseline.yaml |  | The fifty violations the charts carried when the gate first ran, each naming GH-2409. |
|  | cd magefiles && go test . -run TestChartConformance | The gate passes against the seeded baseline and fails on an unlisted or a stale entry. |

<!-- release: 17.0 | status: implemented -->
## 17.0 - Pin currency report

Make a pin's two remaining properties visible: that its recorded digest
is still what the tag resolves to, and that no newer upstream release
has appeared. Report both, change nothing, and let the operator decide
which bumps to take.

<!-- field: parent -->
### Parent

gh-2049-pin-currency

<!-- field: dependencies -->
### Dependencies

| id | status | reason |
| --- | --- | --- |
| srd005-chart-conformance | implemented | Pins every chart image and renders every chart, which this release reads its derived inventory from. |
| eng01-kind-test-demo-rig | implemented | Holds the pin currency section and the bump verb contract GH-2415 adds. |

<!-- field: srds -->
### SRDs

- srd006-pin-currency

<!-- field: use_cases | form: "tree" -->
### Use cases

1. - id: rel17.0-uc001-pin-currency
   - summary:
     An operator runs the bump target, reads which pins moved and which
     aged, and edits the ones they choose to take.
   - status: implemented

<!-- field: implementation_sequence | form: "typed table" -->
### Implementation sequence

| order | action | status |
| --- | --- | --- |
| 1 | Add the ENG01 pin currency section and the bump verb, and specify srd006, the use case, the suite, and the index entries (GH-2415). | planned |
| 2 | Add the pin survey, its registry client, the declared pin list, and the bump target (GH-2416). | implemented |

<!-- field: test_suites | form: "tree" -->
### Test suites

1. - id: test-rel17.0-pin-currency
   - summary:
     Trace the release documents, then verify the derived and declared
     inventory, both verdicts, version comparison, the unknown case, and
     that the target writes nothing.
   - status: implemented
   - implemented_cases:
     1. TC1
     1. TC2
     1. TC3
     1. TC4
     1. TC5
     1. TC6
     1. TC7
     1. TC8
     1. TC9
   - planned_cases: []

<!-- field: implementation_evidence | form: "typed table" -->
### Implementation evidence

| path | command | proves |
| --- | --- | --- |
| docs/engineering/eng01-kind-test-demo-rig.md |  | The integrity and currency properties, and the bump verb's report-never-mutate contract. |
| applications/docs/specs/software-requirements/srd006-pin-currency.md |  | Numbered inventory, integrity, currency, no-write, unknown, and reporting requirements. |
| applications/docs/specs/use-cases/rel17.0-uc001-pin-currency.md |  | The flow from one command to a decision the operator takes. |
| applications/docs/specs/test-suites/test-rel17.0-pin-currency.md |  | Requirement-to-AC-to-test coverage with every case executable. |
| magefiles/pinsurvey/ |  | The distribution client, both verdicts, version comparison, and the inventory's two halves. |
| magefiles/pinsurvey/pins.yaml |  | The pins no render exposes, each naming the file and symbol that holds it. |
| magefiles/bump.go |  | The target joining both halves, skipping images this checkout produces, and reporting without writing. |
|  | mage bump | Eleven pins resolved against their registries with no moved pin and no unknown. |

<!-- release: 18.0 | status: implemented -->
## 18.0 - Release annotation gate

The v0 line's annotation is its only release record, and the record
has been shipping degenerate: three of five recent tags say nothing
but their own name. Refuse that at the push boundary, name the rule
broken, and leave the words to the session that releases.

<!-- field: parent -->
### Parent

gh-2418-tag-verify

<!-- field: dependencies -->
### Dependencies

| id | status | reason |
| --- | --- | --- |
| eng01-kind-test-demo-rig | implemented | Holds the release-annotation section and the tag:verify verb contract GH-2419 adds. |

<!-- field: srds -->
### SRDs

- srd007-release-annotation

<!-- field: use_cases | form: "tree" -->
### Use cases

1. - id: rel18.0-uc001-release-annotation
   - summary:
     An operator about to push runs the check; a real summary passes
     and a degenerate annotation is refused with the rule it broke.
   - status: implemented

<!-- field: implementation_sequence | form: "typed table" -->
### Implementation sequence

| order | action | status |
| --- | --- | --- |
| 1 | Add the ENG01 release-annotation section and the tag:verify verb, and specify srd007, the use case, the suite, and the index entries (GH-2419). | planned |
| 2 | Add the tag:verify target, its scratch-repository tests, and the runbook step (GH-2420). | implemented |

<!-- field: test_suites | form: "tree" -->
### Test suites

1. - id: test-rel18.0-release-annotation
   - summary:
     Trace the release documents, then verify each refusal by name,
     reachability from main, that a run mutates nothing, and the
     placement in the runbook and the skill.
   - status: implemented
   - implemented_cases:
     1. TC1
     1. TC2
     1. TC3
     1. TC4
     1. TC5
     1. TC6
   - planned_cases: []

<!-- field: implementation_evidence | form: "typed table" -->
### Implementation evidence

| path | command | proves |
| --- | --- | --- |
| docs/engineering/eng01-kind-test-demo-rig.md |  | The annotation-is-the-record statement and the verb's refuse-and-mutate-nothing contract. |
| applications/docs/specs/software-requirements/srd007-release-annotation.md |  | Numbered object, degeneracy, reachability, and no-mutation requirements. |
| applications/docs/specs/use-cases/rel18.0-uc001-release-annotation.md |  | The flow from an annotated tag to a push that carries a real record. |
| applications/docs/specs/test-suites/test-rel18.0-release-annotation.md |  | Requirement-to-AC-to-test coverage with every case executable. |
| magefiles/tag_verify.go |  | The four rules with one named reason per failure, and the bar stated once. |
|  | mage tag:verify v0.20260919.1 | The most recent real-summary tag passes and the degenerate v0.20260919.0 is refused naming R2.1. |

<!-- release: 19.0 | status: implemented -->
## 19.0 - Persistent application platform

Replace competing demo, integration, and host-observability assumptions
with one persistent local/cloud platform contract. Share ingress and the
object-storage service while each application retains one namespace,
release, collector, application-specific bucket, lifecycle authority,
and immutable telemetry history.

<!-- field: parent -->
### Parent

gh-2474-shared-platform

<!-- field: dependencies -->
### Dependencies

| id | status | reason |
| --- | --- | --- |
| srd003-application-consistency | implemented | Supplies the normalized application manifest and one composition authority. |
| eng01-kind-test-demo-rig | implemented | Supplies the local platform, scenario, readiness, and evidence baseline amended by GH-2476. |
| eng08-gcp-demo-rig | implemented | Supplies the cloud bootstrap and workload-identity baseline amended by GH-2476. |
| agent-core-srd008-telemetry | implemented | Supplies OTLP export and persistent integration identity. |
| agent-core-srd059-objectstore-tools | implemented | Supplies provider-neutral bucket access across fake GCS, GCS, and S3. |

<!-- field: srds -->
### SRDs

- srd008-persistent-application-platform

<!-- field: use_cases | form: "tree" -->
### Use cases

1. - id: rel19.0-uc001-persistent-application-platform
   - summary:
     Run two applications on one platform, remove and redeploy one, and
     prove that compute, retained data, and query availability follow
     their separate lifecycle boundaries.
   - status: implemented

<!-- field: implementation_sequence | form: "typed table" -->
### Implementation sequence

| order | action | status | note |
| --- | --- | --- | --- |
| 1 | Specify ownership, lifecycle, storage, query, threat, routing, and apprig boundaries (GH-2476). | implemented |  |
| 2 | Add collector WAL and immutable object-bucket persistence (GH-2475). | implemented |  |
| 3 | Add application-scoped object-backed query and retention (GH-2478, split into GH-2489 object-read query engine and GH-2490 collector wiring/retention/purge/cloud isolation). | implemented | GH-2489 supplies object-backed query/WAL merge; GH-2490 supplies collector wiring, retention metadata, storage-status UI, isolated purge, and cloud IAM fixtures. |
| 4 | Wire durable collector configuration through the shared agent-services chart (GH-2484). | implemented |  |
| 5 | Make ingress and filesystem-backed fake GCS persistent platform services (GH-2477). | implemented |  |
| 6 | Add the manifest-driven magefiles/apprig lifecycle runner (GH-2479, split into GH-2492 deterministic core and GH-2493 live proofs/external fixture). | implemented | GH-2492 supplies the deterministic core; GH-2493 supplies live lifecycle, compensation, diagnosis, external-module, retained-history, and approved purge evidence. |
| 7 | Migrate Chatbot Mesh, Coding Agent, and Agent Architecture (GH-2482, GH-2483, GH-2480). | implemented |  |

<!-- field: test_suites | form: "tree" -->
### Test suites

1. - id: test-rel19.0-persistent-application-platform
   - summary:
     Trace the release documents, then prove protocol paths,
     two-application lifecycle, local/cloud parity, WAL-first immutable
     objects, query restoration, failure isolation, and reusable runner
     adoption.
   - status: implemented
   - implemented_cases:
     1. TC1
     1. TC2
     1. TC3
     1. TC4
     1. TC5
     1. TC6
     1. TC7
     1. TC8
   - planned_cases: []

<!-- field: implementation_evidence | form: "typed table" -->
### Implementation evidence

| path | command | proves |
| --- | --- | --- |
| applications/docs/specs/software-requirements/srd008-persistent-application-platform.md |  | Numbered ownership, platform, collector, lifecycle, parity, object, query, routing, and runner requirements. |
| applications/docs/specs/use-cases/rel19.0-uc001-persistent-application-platform.md |  | Complete two-application lifecycle from platform creation through teardown, recovery, purge, and reset. |
| applications/docs/specs/test-suites/test-rel19.0-persistent-application-platform.md |  | Requirement-to-acceptance and acceptance-to-test closure with planned implementation owners. |
| docs/engineering/eng01-kind-test-demo-rig.md |  | Canonical local platform, exception modes, application verbs, persistence, and query availability. |
| docs/engineering/eng08-gcp-demo-rig.md |  | Cloud application buckets, workload identity, bucket IAM, and cluster-independent retention. |
|  | mage audit | Structured documents parse and reciprocal traceability remains valid. |
