---
id: structure-audit-2026-08-19
title: Repository Structure Audit -- 2026-08-19
audit_date: "2026-08-19"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: scope -->
## Scope

<!-- field: issue -->
### Issue

GH-1776

<!-- field: purpose -->
### Purpose

Compare tracked documentation and production command packages with the
repository's declared document locations and Go package roles.

<!-- field: document_placement -->
### Document placement

<!-- field: authority -->
#### Authority

docs/constitutions/design.yaml

<!-- field: method -->
#### Method

List tracked Markdown and YAML files directly below a docs tree after
excluding declared top-level documents, specs, constitutions,
engineering guidelines, fixtures, dependencies, and generated output.

<!-- field: reconciliation -->
#### Reconciliation

The original method excludes docs/engineering. This run includes those
files as scope-reconciliation candidates because a repository-wide
placement gate must classify them too.

<!-- field: package_role -->
### Package role

<!-- field: authority -->
#### Authority

docs/constitutions/go-style.md

<!-- field: method -->
#### Method

Rank production command packages by non-test size, then inspect functions
for runtime behavior that belongs in an internal package. Size orders the
review and is not itself a finding.

<!-- field: finding_gate -->
### Finding gate

- Search tracked files and open GitHub issues for every candidate basename.
- Preserve a live register named by an open recurring issue.
- Verify that a newer register carries prior rejected candidates before removing an older one.
- Classify each document as schema lag, misfiled, spent, or live.
- Name specific functions and an internal destination for every package-role finding.

<!-- field: candidates -->
## Candidates

<!-- field: documents -->
### Documents

<!-- document: agent-core/docs/bench-documentation-ux-inventory.md -->
#### agent-core/docs/bench-documentation-ux-inventory.md

<!-- field: source -->
##### Source

method_1

<!-- field: status -->
##### Status

fixed

<!-- field: classification -->
##### Classification

spent

<!-- field: referrers | form: "tree" -->
##### Referrers

- tracked:
  1. this structure audit register only
- open_issues:
  1. GH-1776
  1. GH-1794
- closed_issues:
  1. GH-87
  1. GH-1026
  1. GH-1031

<!-- field: supersession | form: "tree" -->
##### Supersession

- status: absorbed
- documents:
  1. applications/catalog/docs/specs/software-requirements/srd006-bench.md
  1. applications/catalog/docs/specs/software-requirements/srd011-knowledge-manager.md
  1. applications/catalog/docs/specs/use-cases/rel07.0-uc001-documentation-curator-profile-ux.md
  1. applications/catalog/docs/specs/use-cases/rel07.0-uc004-bench-evaluator-profile-boundary.md

<!-- field: resolution -->
##### Resolution

Deleted by GH-1794 after the finding gate found no live dependency.

<!-- document: agent-core/docs/specs/audits/declarative-decomposition-audit-2026-08-13.md -->
#### agent-core/docs/specs/audits/declarative-decomposition-audit-2026-08-13.md

<!-- field: previous_path -->
##### Previous path

agent-core/docs/declarative-decomposition-audit-2026-08-13.md

<!-- field: source -->
##### Source

method_1

<!-- field: status -->
##### Status

fixed

<!-- field: classification -->
##### Classification

live_audit_report

<!-- field: document_type -->
##### Document type

audit_report

<!-- field: referrers -->
##### Referrers

- issue GH-1639
- this structure audit register

<!-- field: resolution -->
##### Resolution

GH-1791 relocated the report and updated the recurring issue before removing the old path.

<!-- document: agent-core/docs/specs/audits/declarative-decomposition-audit-2026-08.md -->
#### agent-core/docs/specs/audits/declarative-decomposition-audit-2026-08.md

<!-- field: previous_path -->
##### Previous path

agent-core/docs/declarative-decomposition-audit-2026-08.md

<!-- field: source -->
##### Source

method_1

<!-- field: status -->
##### Status

fixed

<!-- field: classification -->
##### Classification

live_audit_report

<!-- field: document_type -->
##### Document type

audit_report

<!-- field: referrers -->
##### Referrers

- agent-core/docs/specs/audits/declarative-decomposition-audit-2026-08-13.md
- this structure audit register

<!-- field: resolution -->
##### Resolution

GH-1791 relocated the prior report after verifying all 79 rejected-candidate rows were carried forward.

<!-- document: agent-core/docs/documentation-divergence-inventory.md -->
#### agent-core/docs/documentation-divergence-inventory.md

<!-- field: source -->
##### Source

method_1

<!-- field: status -->
##### Status

fixed

<!-- field: classification -->
##### Classification

spent

<!-- field: referrers | form: "tree" -->
##### Referrers

- tracked:
  1. this structure audit register only
- open_issues:
  1. GH-1776
  1. GH-1794
- closed_issues:
  1. GH-87
  1. GH-134
  1. GH-325
  1. GH-326
  1. GH-1310

<!-- field: supersession | form: "tree" -->
##### Supersession

- status: absorbed
- documents:
  1. agent-core/README.md
  1. agent-core/docs/ARCHITECTURE.md
  1. agent-core/docs/specs/config-formats/runtime-contract.md
  1. agent-core/package-layout.md

<!-- field: resolution -->
##### Resolution

Deleted by GH-1794 after every recorded drift section named completed work.

<!-- document: agent-core/docs/grammar-machine-conformance-inventory.md -->
#### agent-core/docs/grammar-machine-conformance-inventory.md

<!-- field: source -->
##### Source

method_1

<!-- field: status -->
##### Status

fixed

<!-- field: classification -->
##### Classification

spent

<!-- field: referrers | form: "tree" -->
##### Referrers

- tracked:
  1. this structure audit register only
- open_issues:
  1. GH-1776
  1. GH-1794
- closed_issues:
  1. GH-87
  1. GH-1518

<!-- field: supersession | form: "tree" -->
##### Supersession

- status: absorbed
- documents:
  1. agent-core/docs/ARCHITECTURE.md
  1. agent-core/docs/specs/semantic-models/tool-language.md
  1. agent-core/docs/specs/config-formats/tool-authoring-standard.md

<!-- field: resolution -->
##### Resolution

Deleted by GH-1794 after its seven alignment units were complete.

<!-- document: agent-core/docs/specs/audits/rel03-rest-design-audit.md -->
#### agent-core/docs/specs/audits/rel03-rest-design-audit.md

<!-- field: previous_path -->
##### Previous path

agent-core/docs/rel03-rest-design-audit.md

<!-- field: source -->
##### Source

method_1

<!-- field: status -->
##### Status

fixed

<!-- field: classification -->
##### Classification

live_audit_report

<!-- field: referrers | form: "tree" -->
##### Referrers

- tracked:
  1. agent-core/docs/SPECIFICATIONS.md
  1. this structure audit register
- open_issues:
  1. GH-1776
  1. GH-1792
  1. GH-1794

<!-- field: supersession | form: "tree" -->
##### Supersession

- status: not_superseded

<!-- field: resolution -->
##### Resolution

GH-1792 relocated the report and repaired its active specification reference.

<!-- document: agent-core/pkg/spec/testdata/tool-contract-completeness-baseline.yaml -->
#### agent-core/pkg/spec/testdata/tool-contract-completeness-baseline.yaml

<!-- field: previous_path -->
##### Previous path

agent-core/docs/tool-contract-completeness-baseline.yaml

<!-- field: source -->
##### Source

method_1

<!-- field: status -->
##### Status

fixed

<!-- field: classification -->
##### Classification

code_adjacent_data

<!-- field: referrers -->
##### Referrers

- agent-core/pkg/spec/validate_contract_baseline.go
- agent-core/pkg/spec/tool_declaration_reach_test.go
- agent-core/pkg/spec/validate_contract_baseline_test.go

<!-- field: resolution -->
##### Resolution

GH-1797 relocated the unchanged shrink-only ratchet beside its validator.

<!-- document: applications/catalog/docs/migrations/v0.20260727.0-agent-role-realization-alignment.md -->
#### applications/catalog/docs/migrations/v0.20260727.0-agent-role-realization-alignment.md

<!-- field: source -->
##### Source

method_1

<!-- field: status -->
##### Status

fixed

<!-- field: classification -->
##### Classification

schema_lag

<!-- field: document_type -->
##### Document type

migration

<!-- field: resolution -->
##### Resolution

GH-1795 registered the existing migration convention; no file move was required.

<!-- document: applications/catalog/docs/migrations/v0.20260730.0-collector-family.md -->
#### applications/catalog/docs/migrations/v0.20260730.0-collector-family.md

<!-- field: source -->
##### Source

method_1

<!-- field: status -->
##### Status

fixed

<!-- field: classification -->
##### Classification

schema_lag

<!-- field: document_type -->
##### Document type

migration

<!-- field: resolution -->
##### Resolution

GH-1795 registered the existing migration convention; no file move was required.

<!-- document: applications/catalog/docs/migrations/v0.20260804.0-applier-family.md -->
#### applications/catalog/docs/migrations/v0.20260804.0-applier-family.md

<!-- field: source -->
##### Source

method_1

<!-- field: status -->
##### Status

fixed

<!-- field: classification -->
##### Classification

schema_lag

<!-- field: document_type -->
##### Document type

migration

<!-- field: resolution -->
##### Resolution

GH-1795 registered the existing migration convention; no file move was required.

<!-- document: applications/chatbot-mesh/docs/guides/how-it-works.md -->
#### applications/chatbot-mesh/docs/guides/how-it-works.md

<!-- field: previous_path -->
##### Previous path

applications/chatbot-mesh/docs/how-it-works.md

<!-- field: source -->
##### Source

method_1

<!-- field: status -->
##### Status

fixed

<!-- field: classification -->
##### Classification

misfiled

<!-- field: document_type -->
##### Document type

guide

<!-- field: resolution -->
##### Resolution

GH-1790 moved the guide and repaired every active and relative link.

<!-- document: applications/coding-agent/docs/guides/deployment.md -->
#### applications/coding-agent/docs/guides/deployment.md

<!-- field: previous_path -->
##### Previous path

applications/coding-agent/docs/deployment.md

<!-- field: source -->
##### Source

method_1

<!-- field: status -->
##### Status

fixed

<!-- field: classification -->
##### Classification

misfiled

<!-- field: document_type -->
##### Document type

guide

<!-- field: resolution -->
##### Resolution

GH-1790 moved the guide and repaired README and specification-index links.

<!-- document: applications/coding-agent/docs/migrations/rel14.0-actor-directories.md -->
#### applications/coding-agent/docs/migrations/rel14.0-actor-directories.md

<!-- field: source -->
##### Source

method_1

<!-- field: status -->
##### Status

fixed

<!-- field: classification -->
##### Classification

schema_lag

<!-- field: document_type -->
##### Document type

migration

<!-- field: resolution -->
##### Resolution

GH-1795 registered the existing migration convention; no file move was required.

<!-- document: applications/docs/pattern-language.md -->
#### applications/docs/pattern-language.md

<!-- field: source -->
##### Source

method_1

<!-- field: status -->
##### Status

fixed

<!-- field: classification -->
##### Classification

schema_lag

<!-- field: document_type -->
##### Document type

pattern_language

<!-- field: resolution -->
##### Resolution

GH-1795 registered its active, repository-shared location.

<!-- document: docs/engineering/eng01-kind-test-demo-rig.md -->
#### docs/engineering/eng01-kind-test-demo-rig.md

<!-- field: source -->
##### Source

scope_reconciliation

<!-- field: status -->
##### Status

fixed

<!-- field: classification -->
##### Classification

schema_lag

<!-- field: document_type -->
##### Document type

engineering_guideline

<!-- field: resolution -->
##### Resolution

GH-1795 declared the repository-wide engineering location.

<!-- document: docs/engineering/eng02-agent-ui-placement.md -->
#### docs/engineering/eng02-agent-ui-placement.md

<!-- field: source -->
##### Source

scope_reconciliation

<!-- field: status -->
##### Status

fixed

<!-- field: classification -->
##### Classification

schema_lag

<!-- field: document_type -->
##### Document type

engineering_guideline

<!-- field: resolution -->
##### Resolution

GH-1795 declared the repository-wide engineering location.

<!-- field: packages -->
### Packages

<!-- package: agent-core/cmd/agent -->
#### agent-core/cmd/agent

<!-- field: production_lines | form: "scalar" -->
##### Production lines

2787

<!-- field: status -->
##### Status

findings_filed

<!-- field: files | form: "typed table" -->
##### Files

| path | production_lines | disposition | finding | responsibilities |
| --- | --- | --- | --- | --- |
| agent-core/cmd/agent/main.go | 1139 | retained_composition_root |  | CLI, run assembly, request-signal hosting, checkpoint selection, telemetry, operator output, and exit status. |
| agent-core/cmd/agent/factories.go | 478 | retained_composition_root |  | Selected-init factory catalog and process dependency wiring. |
| agent-core/cmd/agent/runtime_config.go | 470 | finding | GH-1799 | Profile-to-runtime adaptation plus monitor recorder policy compilation. |
| agent-core/cmd/agent/program_ref.go | 254 | finding | GH-1798 | Program closure hashing and originating rollback resource reconstruction. |
| agent-core/cmd/agent/lifecycle_request.go | 188 | retained_composition_root |  | Request selector rewrite, typed lifecycle decode, and target checkpoint adapter selection. |
| agent-core/cmd/agent/factories_llm.go | 112 | retained_composition_root |  | LLM factory construction and process-state resolution hooks. |
| agent-core/cmd/agent/conversation_checkpoint.go | 79 | retained_composition_root |  | Conversation and domain checkpoint JSON provider/resolver adapters. |
| agent-core/cmd/agent/factory_report_parse_error.go | 29 | retained_composition_root |  | report_parse_error ToolDef-to-builder factory. |
| agent-core/cmd/agent/dolt_driver.go | 20 | retained_composition_root |  | Dolt SQL driver registration. |
| agent-core/cmd/agent/child_agent_config.go | 18 | retained_composition_root |  | Child-agent profile config to execute.Config adaptation. |

<!-- package: applications/catalog/cmd/catalog-test-evidence -->
#### applications/catalog/cmd/catalog-test-evidence

<!-- field: production_lines | form: "scalar" -->
##### Production lines

218

<!-- field: status -->
##### Status

rejected

<!-- field: classification -->
##### Classification

build_test_support

<!-- field: files | form: "typed table" -->
##### Files

| path | production_lines | disposition | responsibilities |
| --- | --- | --- | --- |
| applications/catalog/cmd/catalog-test-evidence/runner.go | 200 | retained_command_adapter | Go JSON evidence forwarding and stable compiled conformance test execution. |
| applications/catalog/cmd/catalog-test-evidence/main.go | 18 | retained_entry_point | Process argv, error reporting, and exit status. |

<!-- field: findings | form: "tree" -->
## Findings

1. - id: cmd-agent-program-reconstruction
   - issue: GH-1798
   - rule: go-style package_layout.cmd
   - evidence:
     - file: agent-core/cmd/agent/program_ref.go
     - functions:
       1. augmentRollbackResources (lines 28-51)
       1. definitionsForExecution (lines 53-77)
       1. buildProgramRef (lines 96-116)
       1. loadReferencedProgram (lines 118-143)
       1. runtimeConfigForProfile (lines 145-159)
       1. programAssetFiles (lines 161-190)
       1. addDeclarationClosure (lines 192-222)
       1. addProgramDirectory (lines 224-234)
       1. mergeRESTCollections and mergeMap (lines 236-254)
     - contract:
       1. srd026-lifecycle-tools R3.5
       1. srd026-lifecycle-tools R3.12
     - reason:
       Program closure hashing and verified originating-resource reconstruction
       are catalog and lifecycle algorithms. addDeclarationClosure duplicates
       the include walk already owned by internal/tools/catalog.
   - destination:
     1. internal/tools/catalog for asset closure and ProgramRef digest construction
     1. internal/tools/lifecycle for checkpoint verification and originating ToolDef/REST reconstruction
   - cmd_wiring_remaining:
     1. adapt runtimeConfig and profile paths
     1. select the checkpoint backend and rollback operator profile
     1. register standard factories and call the catalog/lifecycle APIs
   - tests:
     1. cmd/agent TestAugmentRollbackResourcesLoadsVerifiedOriginatingProgram
     1. cmd/agent TestAugmentRollbackResourcesRejectsLegacyCheckpointWithoutProgram
     1. cmd/agent TestAugmentRollbackResourcesRejectsMissingReceiptBuilder
     1. cmd/agent TestFreshRollbackRegistryTracesRepresentativeFamilies
     1. cmd/agent TestFreshRollbackRegistryRestoresOriginatingFileTool
     1. cmd/agent TestFreshRollbackRegistryReportsAliasedSelfInvokeReceipt
1. - id: cmd-agent-monitor-policy-compilation
   - issue: GH-1799
   - rule: go-style package_layout.cmd
   - evidence:
     - file: agent-core/cmd/agent/runtime_config.go
     - functions:
       1. monitorRecorderConfig (lines 237-256)
       1. monitorEnvelopePolicy (lines 258-286)
       1. sortedSetValues (lines 288-297)
       1. toolMetricBindings (lines 299-320)
       1. toolMetricBinding (lines 322-356)
       1. machineMetricLabelValues (lines 358-379)
       1. metricAttributeAllowedValues (lines 381-404)
     - contract:
       1. srd032-tool-metrics-monitor-store R3.1
       1. srd032-tool-metrics-monitor-store R3.2
       1. srd032-tool-metrics-monitor-store R3.6
       1. srd032-tool-metrics-monitor-store R3.8
     - reason:
       These functions interpret MachineSpec and ToolDef metric declarations
       into RecorderConfig, MetricBinding, AttributePolicy, and EnvelopePolicy,
       all types owned by internal/observability/monitor.
   - destination:
     1. internal/observability/monitor
   - cmd_wiring_remaining:
     1. detect REST monitor-view opt-in without creating a monitor-to-rest import cycle
     1. allocate the process store and select the telemetry meter and run identity
     1. pass Store and Recorder into REST dependencies and LoopParams
   - tests:
     1. cmd/agent TestMonitorRuntimeRejectsUnboundedToolAttributeSetup
     1. cmd/agent TestMonitorRuntimeBindsDeclaredToolAndWorkflowAttributes
     1. cmd/agent TestMonitorRuntimeOptInProfileSetsLoopRecorder
     1. cmd/agent TestMonitorRuntimeRecordsDispatchMetricsInStore
     1. cmd/agent TestMachineNameDrivesSpanAndMetricIdentity

<!-- field: rejected_candidates -->
## Rejected candidates

<!-- rejected_candidate: cmd-agent-main-host-and-cli -->
### cmd-agent-main-host-and-cli

<!-- field: gate_question | form: "scalar" -->
#### Gate question

5

<!-- field: files -->
#### Files

- agent-core/cmd/agent/main.go

<!-- field: reason -->
#### Reason

Request-signal listener lifecycle, per-run checkpoint construction,
domain snapshot adapters, startup sequencing, telemetry setup, operator
output, and exit mapping are assigned to cmd/agent by architecture,
srd018, srd035, and srd046. Size is only a review lead; main.go and
buildPreparedRun are already grandfathered by the Go size gate.

<!-- field: prior_evidence -->
#### Prior evidence

- GH-884
- GH-1553
- GH-1743

<!-- rejected_candidate: cmd-agent-factory-wiring -->
### cmd-agent-factory-wiring

<!-- field: gate_question | form: "scalar" -->
#### Gate question

5

<!-- field: files -->
#### Files

- agent-core/cmd/agent/factories.go
- agent-core/cmd/agent/factories_llm.go
- agent-core/cmd/agent/factory_report_parse_error.go

<!-- field: reason -->
#### Reason

These functions decode trusted ToolDef config, bind process dependencies,
and return builders implemented by owning internal packages. They contain
no independently selectable runtime operation and have no narrower
internal destination than the packages already providing the builders.

<!-- rejected_candidate: cmd-agent-checkpoint-and-child-adapters -->
### cmd-agent-checkpoint-and-child-adapters

<!-- field: gate_question | form: "scalar" -->
#### Gate question

5

<!-- field: files -->
#### Files

- agent-core/cmd/agent/lifecycle_request.go
- agent-core/cmd/agent/conversation_checkpoint.go
- agent-core/cmd/agent/dolt_driver.go
- agent-core/cmd/agent/child_agent_config.go

<!-- field: reason -->
#### Reason

The runtime contract assigns request selector rewrite, conversation/domain
JSON provider and resolver ports, SQL driver registration, and child
profile adaptation to the composition boundary. Their implementations
already delegate checkpoint, Dolt, lifecycle, and child execution behavior
to internal packages.

<!-- field: prior_evidence -->
#### Prior evidence

- GH-1559
- GH-1683

<!-- rejected_candidate: catalog-test-evidence-command -->
### catalog-test-evidence-command

<!-- field: gate_question | form: "scalar" -->
#### Gate question

5

<!-- field: files -->
#### Files

- applications/catalog/cmd/catalog-test-evidence/main.go
- applications/catalog/cmd/catalog-test-evidence/runner.go

<!-- field: classification -->
#### Classification

build_test_support

<!-- field: reason -->
#### Reason

The only production caller is a Mage target that builds this helper into
a temporary directory and rewrites a staged specification-critic profile.
No shipped image or agent selects it. Moving it behind a product word or
into agent-core would weaken the independent Go JSON evidence boundary.

<!-- field: tests -->
#### Tests

- applications/catalog/cmd/catalog-test-evidence/runner_test.go
- applications/catalog/magefiles/test_evidence_test.go

<!-- field: prior_evidence -->
#### Prior evidence

- GH-1395
- GH-1580
- GH-1636

<!-- field: results -->
## Results

<!-- field: document_placement -->
### Document placement

<!-- field: before | form: "tree" -->
#### Before

- method_1_candidates: 14
- scope_reconciliation_candidates: 2
- total_candidates: 16

<!-- field: after -->
#### After

<!-- field: method_1_output | form: "scalar" -->
##### Method 1 output

7

<!-- field: unclassified | form: "scalar" -->
##### Unclassified

0

<!-- field: explanation -->
##### Explanation

The remaining Method 1 output consists of four migration records, two
guides, and the shared application pattern language. The command is a
broad candidate generator; each remaining path now matches its declared
type and has a fixed disposition above.

<!-- field: changes | form: "tree" -->
#### Changes

- document_types_declared: 6
- guides_relocated: 2
- audit_reports_relocated: 3
- spent_documents_deleted: 3
- code_adjacent_data_relocated: 1
- paths_typed_without_move: 7

<!-- field: package_role -->
### Package role

<!-- field: before | form: "tree" -->
#### Before

- command_packages: 2
- production_files: 12
- production_lines: 3005

<!-- field: after -->
#### After

<!-- field: unclassified_packages | form: "scalar" -->
##### Unclassified packages

0

<!-- field: unclassified_files | form: "scalar" -->
##### Unclassified files

0

<!-- field: findings -->
##### Findings

- GH-1798
- GH-1799

<!-- field: rejected_candidate_groups | form: "scalar" -->
##### Rejected candidate groups

4

<!-- field: unresolved_external_work -->
### Unresolved external work

| issue | reason |
| --- | --- |
| GH-1798 | Move program identity and rollback reconstruction into catalog and lifecycle packages. |
| GH-1799 | Move monitor recorder policy compilation into the monitor package. |
| GH-1781 | Rebase the placement-gate epic after this audit merges so it retains one document type registry. |

<!-- field: verification -->
## Verification

<!-- field: inventory_commands -->
### Inventory commands

- git ls-files tracked docs candidate command from GH-1776 Method 1
- production command package ranking command from GH-1776 Method 2

<!-- field: yaml_parse -->
### YAML parse

passed

<!-- field: mage_audit -->
### Mage audit

passed

<!-- field: mage_lint -->
### Mage lint

passed

<!-- field: mage_test -->
### Mage test

passed
