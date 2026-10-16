---
id: test-rel00.0-application-contract-conformance
title: Shared Application Contract Conformance Suite
release: "00.0"
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

This suite validates structural conformance. It parses the shared and local
YAML corpora, checks reciprocal traces, executes root module-classification
and accounting tests, and relies on each application's own audit and test
suites for business and capability behavior. An undeclared capability is not
applicable; its absent targets are not treated as passing evidence.

<!-- field: traces -->
## Traces

- rel00.0-uc001-application-contract-conformance

<!-- field: preconditions -->
## Preconditions

- Run commands from the repository root.
- Go, Mage, and Python 3 with PyYAML are available.
- agent-core and applications/catalog are present in the monorepo checkout for application-local audit dependencies.

<!-- field: test_cases -->
## Test cases

<!-- case: Shared corpus parses and traces reciprocally | status: implemented -->
### Shared corpus parses and traces reciprocally

<!-- field: id -->
#### Id

TC1

<!-- field: use_case -->
#### Use case

rel00.0-uc001-application-contract-conformance

<!-- field: evidence_type -->
#### Evidence type

executable_structural

<!-- field: command -->
#### Command

```yaml
|
    python3 - <<'PY'
    from pathlib import Path
    import yaml

    root = Path("applications/docs")
    docs = {p: yaml.safe_load(p.read_text()) for p in root.rglob("*.yaml")}
    required = {
        root / "VISION.yaml": {
            "id", "title", "executive_summary", "problem", "what_this_does",
            "why_we_build_this", "success_criteria", "not",
        },
        root / "ARCHITECTURE.yaml": {
            "id", "title", "overview", "interfaces", "components",
            "design_decisions", "technology_choices", "project_structure",
            "implementation_status", "related_documents",
        },
        root / "road-map.yaml": {"id", "title", "releases"},
        root / "SPECIFICATIONS.yaml": {
            "id", "title", "overview", "roadmap_summary",
            "foundation_document_index", "srd_index", "config_format_index",
            "semantic_model_index", "use_case_index", "test_suite_index",
            "coverage_gaps",
        },
    }
    for path, fields in required.items():
        assert path in docs, path
        assert not fields - set(docs[path]), (path, fields - set(docs[path]))
    index = docs[root / "SPECIFICATIONS.yaml"]
    for section in ("foundation_document_index", "srd_index",
                    "use_case_index", "test_suite_index"):
        for entry in index[section]:
            path = Path("applications") / entry["path"]
            assert path.is_file(), path
    uc_path = root / "specs/use-cases/rel00.0-uc001-application-contract-conformance.yaml"
    ts_path = root / "specs/test-suites/test-rel00.0-application-contract-conformance.yaml"
    uc, suite = docs[uc_path], docs[ts_path]
    assert uc["test_suite"] == suite["id"]
    assert uc["id"] in suite["traces"]
    assert "test-rel00.0-application-contract-conformance" in {
        trace["target"] for trace in uc["touchpoints"]
    }
    for srd_name in (
        "srd001-application-module-contract.yaml",
        "srd002-managed-application-services.yaml",
    ):
        srd = docs[root / "specs/software-requirements" / srd_name]
        assert uc["id"] in srd["use_case_traces"]
        assert suite["id"] in srd["test_suite_traces"]
    print(f"validated {len(docs)} shared application YAML documents")
    PY
```

<!-- field: expected -->
#### Expected

```yaml
yaml_parses: true
format_fields_present: true
indexed_paths_resolve: true
reciprocal_traces_resolve: true
```

<!-- field: traces -->
#### Traces

- rel00.0-uc001-application-contract-conformance S3
- srd001-application-module-contract AC3

<!-- case: Root module classification and dispatch are pinned | status: implemented -->
### Root module classification and dispatch are pinned

<!-- field: id -->
#### Id

TC2

<!-- field: use_case -->
#### Use case

rel00.0-uc001-application-contract-conformance

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: command -->
#### Command

cd magefiles
go test . -run 'Test(ApplicationModulesParticipateInAudit|ChatbotMeshIsAnApplicationModule|CodingAgentParticipatesInAudit|AgentArchitectureIsCompositionApplicationModule|EveryOrchestratedModuleDirectoryExists|OrchestrationUsesStableApplicationPaths|ApplicationModulesExcludedFromSubModules|StatsParticipantsIncludeApplicationModules|TestSubModulesDispatchesApplicationModules)$'

<!-- field: evidence -->
#### Evidence

- magefiles/build.go
- magefiles/example_modules_test.go
- magefiles/test.go
- magefiles/tag.go

<!-- field: expected -->
#### Expected

```yaml
catalog_is_not_runnable_application: true
runnable_paths_exist: true
root_audit_dispatches_applications: true
root_test_dispatches_applications_with_tests: true
root_stats_dispatches_applications: true
applications_excluded_from_platform_build_default: true
agent_architecture_registered: true
catalog_is_reusable_submodule: true
```

<!-- field: traces -->
#### Traces

- rel00.0-uc001-application-contract-conformance S1
- rel00.0-uc001-application-contract-conformance S4
- srd001-application-module-contract AC1
- srd001-application-module-contract AC4
- srd001-application-module-contract AC6

<!-- case: Canonical agents are counted once | status: implemented -->
### Canonical agents are counted once

<!-- field: id -->
#### Id

TC3

<!-- field: use_case -->
#### Use case

rel00.0-uc001-application-contract-conformance

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: command -->
#### Command

(cd magefiles && go test . -run 'TestSumAgentsTotals')
(cd applications/coding-agent && go test ./magefiles -run 'TestCollectCodingApplicationStats')
(cd applications/chatbot-mesh && go test ./magefiles -run 'Test(ScanAgentOwnershipSeparatesCompositionWrapper|RepositoryMeshOwnershipClassifiesCorpusIngestAsComposition)')
(cd applications/agent-architecture && go test ./magefiles -run '^TestStatsCompositionOnlyJSON$')

<!-- field: evidence -->
#### Evidence

- magefiles/stats.go
- magefiles/stats_test.go
- applications/coding-agent/magefiles/stats.go
- applications/coding-agent/magefiles/stats_test.go
- applications/chatbot-mesh/magefiles/stats.go
- applications/chatbot-mesh/magefiles/agent_stats.go
- applications/chatbot-mesh/magefiles/agent_stats_test.go
- applications/agent-architecture/magefiles/magefile.go
- applications/agent-architecture/magefiles/magefile_test.go

<!-- field: expected -->
#### Expected

```yaml
composition_only_has_no_agents_section: true
local_agents_counted_by_owner: true
composition_references_reported_separately: true
root_total_sums_agents_sections_only: true
canonical_agent_total: 15
```

<!-- field: traces -->
#### Traces

- rel00.0-uc001-application-contract-conformance S2
- srd001-application-module-contract AC2
- srd001-application-module-contract AC5

<!-- case: Current runnable applications pass their owned audit gates | status: implemented -->
### Current runnable applications pass their owned audit gates

<!-- field: id -->
#### Id

TC4

<!-- field: use_case -->
#### Use case

rel00.0-uc001-application-contract-conformance

<!-- field: evidence_type -->
#### Evidence type

mage_gate

<!-- field: command -->
#### Command

mage audit

<!-- field: evidence -->
#### Evidence

- applications/chatbot-mesh/magefiles/magefile.go
- applications/coding-agent/magefiles/magefile.go
- applications/agent-architecture/magefiles/magefile.go
- applications/chatbot-mesh/docs/SPECIFICATIONS.yaml
- applications/coding-agent/docs/SPECIFICATIONS.yaml
- applications/agent-architecture/docs/SPECIFICATIONS.yaml

<!-- field: expected -->
#### Expected

```yaml
root_audit_green: true
chatbot_mesh_audit_executed: true
coding_agent_audit_executed: true
agent_architecture_audit_executed: true
local_evidence_owned_by_application: true
```

<!-- field: traces -->
#### Traces

- rel00.0-uc001-application-contract-conformance S3
- rel00.0-uc001-application-contract-conformance S4
- srd001-application-module-contract AC3
- srd001-application-module-contract AC4

<!-- case: Managed-service evidence status is application-owned and truthful | status: implemented -->
### Managed-service evidence status is application-owned and truthful

<!-- field: id -->
#### Id

TC5

<!-- field: use_case -->
#### Use case

rel00.0-uc001-application-contract-conformance

<!-- field: evidence_type -->
#### Evidence type

application_test_suites

<!-- field: command -->
#### Command

(cd applications/coding-agent && go test ./magefiles -run 'TestServingProfilesUseRealLifecycleAndCanonicalRoleProfiles|TestPlannerServingFlowUsesDeclaredRemoteBoundaries')
(cd applications/chatbot-mesh && go test ./magefiles -run 'TestChatbotRolloutDrainsActiveRequests|TestDetachedAgentDualExportsAndStampsResourceIdentity')
(cd applications/agent-architecture && go test ./magefiles -run '^TestManagedServicesEvidenceContract$')

<!-- field: evidence -->
#### Evidence

- applications/coding-agent/docs/specs/test-suites/test-rel02.0-serving-contract.yaml
- applications/chatbot-mesh/docs/specs/test-suites/test-rel03.0-mesh-deployment.yaml
- applications/chatbot-mesh/docs/specs/test-suites/test-rel02.0-observability.yaml
- applications/agent-architecture/docs/specs/test-suites/test-rel00.0-guided-agent-architecture.yaml

<!-- field: expected -->
#### Expected

```yaml
lifecycle_evidence_exists: true
health_is_role_owned: true
authority_is_declared: true
telemetry_identity_is_tested: true
graceful_shutdown_is_tested: true
knowledge_manager_deterministic_evidence: implemented
knowledge_manager_managed_service_status: partial
knowledge_manager_live_observations: dependency_gated
knowledge_manager_srd002_ac1_ac3_satisfied: false
```

<!-- field: traces -->
#### Traces

- rel00.0-uc001-application-contract-conformance S5
- srd002-managed-application-services AC1
- srd002-managed-application-services AC2
- srd002-managed-application-services AC3

<!-- case: Package and Helm claims have executable evidence | status: implemented -->
### Package and Helm claims have executable evidence

<!-- field: id -->
#### Id

TC6

<!-- field: use_case -->
#### Use case

rel00.0-uc001-application-contract-conformance

<!-- field: evidence_type -->
#### Evidence type

application_test_suites

<!-- field: command -->
#### Command

(cd applications/coding-agent && go test ./magefiles -run 'Test(ProfilePackageIsDeterministic|DeploymentPackageIsDeterministic|HelmPrepareStagesPackageAsOnlyProfileSource|HelmPackage)')
(cd applications/chatbot-mesh && go test ./magefiles -run 'Test(StagedProfilesCoverEnabledDeployments|HelmPackage)')

<!-- field: evidence -->
#### Evidence

- applications/coding-agent/docs/specs/test-suites/test-rel03.0-deployment-profile-package.yaml
- applications/coding-agent/docs/specs/test-suites/test-rel04.0-helm-topology.yaml
- applications/coding-agent/docs/specs/test-suites/test-rel05.0-helm-package-validation.yaml
- applications/chatbot-mesh/docs/specs/test-suites/test-rel03.0-mesh-deployment.yaml

<!-- field: expected -->
#### Expected

```yaml
deterministic_package_evidence_exists: true
unsafe_input_rejection_exists: true
helm_schema_and_render_evidence_exists: true
packaged_archive_validation_exists: true
```

<!-- field: traces -->
#### Traces

- rel00.0-uc001-application-contract-conformance S6
- srd002-managed-application-services AC4
- srd002-managed-application-services AC5

<!-- case: kind-demo claims defer to ENG01 evidence | status: implemented -->
### kind-demo claims defer to ENG01 evidence

<!-- field: id -->
#### Id

TC7

<!-- field: use_case -->
#### Use case

rel00.0-uc001-application-contract-conformance

<!-- field: evidence_type -->
#### Evidence type

referenced_engineering_contract

<!-- field: command -->
#### Command

(cd applications/coding-agent && go test ./magefiles -run 'TestCodingHelm(SmokeSkipsMissingExternalPrerequisite|InfrastructureProbeClassification|DiagnosticsUseOneBoundedContext)')
(cd applications/chatbot-mesh && go test ./magefiles -run 'TestHelmSmokeSkipReasonMissingBinary|TestCollectSharedTelemetryEvidenceIdentifiesSlowComponent')

<!-- field: evidence -->
#### Evidence

- docs/engineering/eng01-kind-test-demo-rig.yaml
- applications/coding-agent/docs/specs/test-suites/test-rel06.0-live-helm-smoke.yaml
- applications/chatbot-mesh/docs/specs/test-suites/test-rel03.0-mesh-deployment.yaml

<!-- field: expected -->
#### Expected

```yaml
eng01_is_referenced_not_copied: true
missing_prerequisite_classification_is_tested: true
bounded_failure_evidence_is_tested: true
undeclared_kind_demo_is_not_applicable: true
```

<!-- field: traces -->
#### Traces

- rel00.0-uc001-application-contract-conformance S7
- srd002-managed-application-services AC6

<!-- case: Root statistics complete after conformance | status: implemented -->
### Root statistics complete after conformance

<!-- field: id -->
#### Id

TC8

<!-- field: use_case -->
#### Use case

rel00.0-uc001-application-contract-conformance

<!-- field: evidence_type -->
#### Evidence type

mage_gate

<!-- field: command -->
#### Command

mage stats

<!-- field: evidence -->
#### Evidence

- magefiles/stats.go
- applications/chatbot-mesh/magefiles/stats.go
- applications/coding-agent/magefiles/stats.go
- applications/agent-architecture/magefiles/magefile.go

<!-- field: expected -->
#### Expected

```yaml
json_output: true
application_modules_present: true
repository_agent_total_present: true
repository_agent_total: 15
knowledge_manager_agents_contributed: 0
```

<!-- field: traces -->
#### Traces

- rel00.0-uc001-application-contract-conformance S2
- rel00.0-uc001-application-contract-conformance S4
- srd001-application-module-contract AC4
- srd001-application-module-contract AC5

<!-- case: Actor classifications and reference-role bindings conform | status: implemented -->
### Actor classifications and reference-role bindings conform

<!-- field: id -->
#### Id

TC9

<!-- field: use_case -->
#### Use case

rel00.0-uc001-application-contract-conformance

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: command -->
#### Command

(cd magefiles && go test . -run '^TestAgentRoleRealization')
(cd applications/catalog && go test ./conformance -run 'Test(ShippedProfilesMatchFoundationModels|CatalogMembershipUsesSharedRealizationAndAliasAuthority)$')

<!-- field: evidence -->
#### Evidence

- magefiles/agent_role_realizations_test.go
- magefiles/audit.go
- applications/docs/specs/semantic-models/agent-role-realizations.yaml
- applications/catalog/conformance/family_inventory_test.go
- applications/catalog/conformance/membership_contract_test.go

<!-- field: expected -->
#### Expected

```yaml
deterministic_profile_discovery: true
canonical_role_and_sub_role_authority_validated: true
every_production_profile_classified: true
repeated_role_bindings_allowed: true
multi_role_actor_bindings_allowed: true
inheritance_and_alias_rules_validated: true
canonical_name_collisions_rejected: true
root_mage_audit_invokes_validation: true
```

<!-- field: traces -->
#### Traces

- rel00.0-uc001-application-contract-conformance S3
- rel00.0-uc001-application-contract-conformance S4
- sm-agent-role-realizations
