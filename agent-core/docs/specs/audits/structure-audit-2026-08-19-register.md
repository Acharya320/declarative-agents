---
id: structure-audit-2026-08-19
title: Repository Structure Audit -- 2026-08-19
audit_date: "2026-08-19"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: scope | form: "tree" -->
## Scope

- issue: GH-1776
- purpose:
  Compare tracked documentation and production command packages with the
  repository's declared document locations and Go package roles.
- document_placement:
  - authority: docs/constitutions/design.yaml
  - method:
    List tracked Markdown and YAML files directly below a docs tree after
    excluding declared top-level documents, specs, constitutions,
    engineering guidelines, fixtures, dependencies, and generated output.
  - reconciliation:
    The original method excludes docs/engineering. This run includes those
    files as scope-reconciliation candidates because a repository-wide
    placement gate must classify them too.
- package_role:
  - authority: docs/constitutions/go-style.md
  - method:
    Rank production command packages by non-test size, then inspect functions
    for runtime behavior that belongs in an internal package. Size orders the
    review and is not itself a finding.
- finding_gate:
  1. Search tracked files and open GitHub issues for every candidate basename.
  1. Preserve a live register named by an open recurring issue.
  1. Verify that a newer register carries prior rejected candidates before removing an older one.
  1. Classify each document as schema lag, misfiled, spent, or live.
  1. Name specific functions and an internal destination for every package-role finding.

<!-- field: candidates | form: "tree" -->
## Candidates

- documents:
  1. - path: agent-core/docs/bench-documentation-ux-inventory.md
     - source: method_1
     - status: fixed
     - classification: spent
     - referrers:
       - tracked:
         1. this structure audit register only
       - open_issues:
         1. GH-1776
         1. GH-1794
       - closed_issues:
         1. GH-87
         1. GH-1026
         1. GH-1031
     - supersession:
       - status: absorbed
       - documents:
         1. applications/catalog/docs/specs/software-requirements/srd006-bench.md
         1. applications/catalog/docs/specs/software-requirements/srd011-knowledge-manager.md
         1. applications/catalog/docs/specs/use-cases/rel07.0-uc001-documentation-curator-profile-ux.md
         1. applications/catalog/docs/specs/use-cases/rel07.0-uc004-bench-evaluator-profile-boundary.md
     - resolution: Deleted by GH-1794 after the finding gate found no live dependency.
  1. - path: agent-core/docs/specs/audits/declarative-decomposition-audit-2026-08-13.md
     - previous_path: agent-core/docs/declarative-decomposition-audit-2026-08-13.md
     - source: method_1
     - status: fixed
     - classification: live_audit_report
     - document_type: audit_report
     - referrers:
       1. issue GH-1639
       1. this structure audit register
     - resolution: GH-1791 relocated the report and updated the recurring issue before removing the old path.
  1. - path: agent-core/docs/specs/audits/declarative-decomposition-audit-2026-08.md
     - previous_path: agent-core/docs/declarative-decomposition-audit-2026-08.md
     - source: method_1
     - status: fixed
     - classification: live_audit_report
     - document_type: audit_report
     - referrers:
       1. agent-core/docs/specs/audits/declarative-decomposition-audit-2026-08-13.md
       1. this structure audit register
     - resolution: GH-1791 relocated the prior report after verifying all 79 rejected-candidate rows were carried forward.
  1. - path: agent-core/docs/documentation-divergence-inventory.md
     - source: method_1
     - status: fixed
     - classification: spent
     - referrers:
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
     - supersession:
       - status: absorbed
       - documents:
         1. agent-core/README.md
         1. agent-core/docs/ARCHITECTURE.md
         1. agent-core/docs/specs/config-formats/runtime-contract.md
         1. agent-core/package-layout.md
     - resolution: Deleted by GH-1794 after every recorded drift section named completed work.
  1. - path: agent-core/docs/grammar-machine-conformance-inventory.md
     - source: method_1
     - status: fixed
     - classification: spent
     - referrers:
       - tracked:
         1. this structure audit register only
       - open_issues:
         1. GH-1776
         1. GH-1794
       - closed_issues:
         1. GH-87
         1. GH-1518
     - supersession:
       - status: absorbed
       - documents:
         1. agent-core/docs/ARCHITECTURE.md
         1. agent-core/docs/specs/semantic-models/tool-language.md
         1. agent-core/docs/specs/config-formats/tool-authoring-standard.md
     - resolution: Deleted by GH-1794 after its seven alignment units were complete.
  1. - path: agent-core/docs/specs/audits/rel03-rest-design-audit.md
     - previous_path: agent-core/docs/rel03-rest-design-audit.md
     - source: method_1
     - status: fixed
     - classification: live_audit_report
     - referrers:
       - tracked:
         1. agent-core/docs/SPECIFICATIONS.md
         1. this structure audit register
       - open_issues:
         1. GH-1776
         1. GH-1792
         1. GH-1794
     - supersession:
       - status: not_superseded
     - resolution: GH-1792 relocated the report and repaired its active specification reference.
  1. - path: agent-core/pkg/spec/testdata/tool-contract-completeness-baseline.yaml
     - previous_path: agent-core/docs/tool-contract-completeness-baseline.yaml
     - source: method_1
     - status: fixed
     - classification: code_adjacent_data
     - referrers:
       1. agent-core/pkg/spec/validate_contract_baseline.go
       1. agent-core/pkg/spec/tool_declaration_reach_test.go
       1. agent-core/pkg/spec/validate_contract_baseline_test.go
     - resolution: GH-1797 relocated the unchanged shrink-only ratchet beside its validator.
  1. - path: applications/catalog/docs/migrations/v0.20260727.0-agent-role-realization-alignment.md
     - source: method_1
     - status: fixed
     - classification: schema_lag
     - document_type: migration
     - resolution: GH-1795 registered the existing migration convention; no file move was required.
  1. - path: applications/catalog/docs/migrations/v0.20260730.0-collector-family.md
     - source: method_1
     - status: fixed
     - classification: schema_lag
     - document_type: migration
     - resolution: GH-1795 registered the existing migration convention; no file move was required.
  1. - path: applications/catalog/docs/migrations/v0.20260804.0-applier-family.md
     - source: method_1
     - status: fixed
     - classification: schema_lag
     - document_type: migration
     - resolution: GH-1795 registered the existing migration convention; no file move was required.
  1. - path: applications/chatbot-mesh/docs/guides/how-it-works.md
     - previous_path: applications/chatbot-mesh/docs/how-it-works.md
     - source: method_1
     - status: fixed
     - classification: misfiled
     - document_type: guide
     - resolution: GH-1790 moved the guide and repaired every active and relative link.
  1. - path: applications/coding-agent/docs/guides/deployment.md
     - previous_path: applications/coding-agent/docs/deployment.md
     - source: method_1
     - status: fixed
     - classification: misfiled
     - document_type: guide
     - resolution: GH-1790 moved the guide and repaired README and specification-index links.
  1. - path: applications/coding-agent/docs/migrations/rel14.0-actor-directories.md
     - source: method_1
     - status: fixed
     - classification: schema_lag
     - document_type: migration
     - resolution: GH-1795 registered the existing migration convention; no file move was required.
  1. - path: applications/docs/pattern-language.md
     - source: method_1
     - status: fixed
     - classification: schema_lag
     - document_type: pattern_language
     - resolution: GH-1795 registered its active, repository-shared location.
  1. - path: docs/engineering/eng01-kind-test-demo-rig.md
     - source: scope_reconciliation
     - status: fixed
     - classification: schema_lag
     - document_type: engineering_guideline
     - resolution: GH-1795 declared the repository-wide engineering location.
  1. - path: docs/engineering/eng02-agent-ui-placement.md
     - source: scope_reconciliation
     - status: fixed
     - classification: schema_lag
     - document_type: engineering_guideline
     - resolution: GH-1795 declared the repository-wide engineering location.
- packages:
  1. - path: agent-core/cmd/agent
     - production_lines: 2787
     - status: findings_filed
     - files:
       1. - path: agent-core/cmd/agent/main.go
          - production_lines: 1139
          - disposition: retained_composition_root
          - responsibilities: CLI, run assembly, request-signal hosting, checkpoint selection, telemetry, operator output, and exit status.
       1. - path: agent-core/cmd/agent/factories.go
          - production_lines: 478
          - disposition: retained_composition_root
          - responsibilities: Selected-init factory catalog and process dependency wiring.
       1. - path: agent-core/cmd/agent/runtime_config.go
          - production_lines: 470
          - disposition: finding
          - finding: GH-1799
          - responsibilities: Profile-to-runtime adaptation plus monitor recorder policy compilation.
       1. - path: agent-core/cmd/agent/program_ref.go
          - production_lines: 254
          - disposition: finding
          - finding: GH-1798
          - responsibilities: Program closure hashing and originating rollback resource reconstruction.
       1. - path: agent-core/cmd/agent/lifecycle_request.go
          - production_lines: 188
          - disposition: retained_composition_root
          - responsibilities: Request selector rewrite, typed lifecycle decode, and target checkpoint adapter selection.
       1. - path: agent-core/cmd/agent/factories_llm.go
          - production_lines: 112
          - disposition: retained_composition_root
          - responsibilities: LLM factory construction and process-state resolution hooks.
       1. - path: agent-core/cmd/agent/conversation_checkpoint.go
          - production_lines: 79
          - disposition: retained_composition_root
          - responsibilities: Conversation and domain checkpoint JSON provider/resolver adapters.
       1. - path: agent-core/cmd/agent/factory_report_parse_error.go
          - production_lines: 29
          - disposition: retained_composition_root
          - responsibilities: report_parse_error ToolDef-to-builder factory.
       1. - path: agent-core/cmd/agent/dolt_driver.go
          - production_lines: 20
          - disposition: retained_composition_root
          - responsibilities: Dolt SQL driver registration.
       1. - path: agent-core/cmd/agent/child_agent_config.go
          - production_lines: 18
          - disposition: retained_composition_root
          - responsibilities: Child-agent profile config to execute.Config adaptation.
  1. - path: applications/catalog/cmd/catalog-test-evidence
     - production_lines: 218
     - status: rejected
     - classification: build_test_support
     - files:
       1. - path: applications/catalog/cmd/catalog-test-evidence/runner.go
          - production_lines: 200
          - disposition: retained_command_adapter
          - responsibilities: Go JSON evidence forwarding and stable compiled conformance test execution.
       1. - path: applications/catalog/cmd/catalog-test-evidence/main.go
          - production_lines: 18
          - disposition: retained_entry_point
          - responsibilities: Process argv, error reporting, and exit status.

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

<!-- field: rejected_candidates | form: "tree" -->
## Rejected candidates

1. - id: cmd-agent-main-host-and-cli
   - gate_question: 5
   - files:
     1. agent-core/cmd/agent/main.go
   - reason:
     Request-signal listener lifecycle, per-run checkpoint construction,
     domain snapshot adapters, startup sequencing, telemetry setup, operator
     output, and exit mapping are assigned to cmd/agent by architecture,
     srd018, srd035, and srd046. Size is only a review lead; main.go and
     buildPreparedRun are already grandfathered by the Go size gate.
   - prior_evidence:
     1. GH-884
     1. GH-1553
     1. GH-1743
1. - id: cmd-agent-factory-wiring
   - gate_question: 5
   - files:
     1. agent-core/cmd/agent/factories.go
     1. agent-core/cmd/agent/factories_llm.go
     1. agent-core/cmd/agent/factory_report_parse_error.go
   - reason:
     These functions decode trusted ToolDef config, bind process dependencies,
     and return builders implemented by owning internal packages. They contain
     no independently selectable runtime operation and have no narrower
     internal destination than the packages already providing the builders.
1. - id: cmd-agent-checkpoint-and-child-adapters
   - gate_question: 5
   - files:
     1. agent-core/cmd/agent/lifecycle_request.go
     1. agent-core/cmd/agent/conversation_checkpoint.go
     1. agent-core/cmd/agent/dolt_driver.go
     1. agent-core/cmd/agent/child_agent_config.go
   - reason:
     The runtime contract assigns request selector rewrite, conversation/domain
     JSON provider and resolver ports, SQL driver registration, and child
     profile adaptation to the composition boundary. Their implementations
     already delegate checkpoint, Dolt, lifecycle, and child execution behavior
     to internal packages.
   - prior_evidence:
     1. GH-1559
     1. GH-1683
1. - id: catalog-test-evidence-command
   - gate_question: 5
   - files:
     1. applications/catalog/cmd/catalog-test-evidence/main.go
     1. applications/catalog/cmd/catalog-test-evidence/runner.go
   - classification: build_test_support
   - reason:
     The only production caller is a Mage target that builds this helper into
     a temporary directory and rewrites a staged specification-critic profile.
     No shipped image or agent selects it. Moving it behind a product word or
     into agent-core would weaken the independent Go JSON evidence boundary.
   - tests:
     1. applications/catalog/cmd/catalog-test-evidence/runner_test.go
     1. applications/catalog/magefiles/test_evidence_test.go
   - prior_evidence:
     1. GH-1395
     1. GH-1580
     1. GH-1636

<!-- field: results | form: "tree" -->
## Results

- document_placement:
  - before:
    - method_1_candidates: 14
    - scope_reconciliation_candidates: 2
    - total_candidates: 16
  - after:
    - method_1_output: 7
    - unclassified: 0
    - explanation:
      The remaining Method 1 output consists of four migration records, two
      guides, and the shared application pattern language. The command is a
      broad candidate generator; each remaining path now matches its declared
      type and has a fixed disposition above.
  - changes:
    - document_types_declared: 6
    - guides_relocated: 2
    - audit_reports_relocated: 3
    - spent_documents_deleted: 3
    - code_adjacent_data_relocated: 1
    - paths_typed_without_move: 7
- package_role:
  - before:
    - command_packages: 2
    - production_files: 12
    - production_lines: 3005
  - after:
    - unclassified_packages: 0
    - unclassified_files: 0
    - findings:
      1. GH-1798
      1. GH-1799
    - rejected_candidate_groups: 4
- unresolved_external_work:
  1. - issue: GH-1798
     - reason: Move program identity and rollback reconstruction into catalog and lifecycle packages.
  1. - issue: GH-1799
     - reason: Move monitor recorder policy compilation into the monitor package.
  1. - issue: GH-1781
     - reason: Rebase the placement-gate epic after this audit merges so it retains one document type registry.

<!-- field: verification | form: "tree" -->
## Verification

- inventory_commands:
  1. git ls-files tracked docs candidate command from GH-1776 Method 1
  1. production command package ranking command from GH-1776 Method 2
- yaml_parse: passed
- mage_audit: passed
- mage_lint: passed
- mage_test: passed
