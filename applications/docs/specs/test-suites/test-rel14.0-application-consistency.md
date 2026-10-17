---
id: test-rel14.0-application-consistency
title: Release 14.0 Application Consistency Test Suite
release: "14.0"
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: test_case_count | form: "scalar" -->
## Test case count

8

<!-- field: overview -->
## Overview

All eight cases are executable. TC1 and TC8 parse and audit the shared
Release 14.0 specification. TC2 through TC7 validate the three application
manifests, actor grammar, manifest-derived package assets, canonical
promotion and UI-token consumption, README/status/root agreement, ownership
accounting, and the current application matrix. Optional capability cases run
only for applications that declare the matching capability.

<!-- field: traces -->
## Traces

- rel14.0-uc001-application-consistency

<!-- field: srds -->
## Srds

- srd003-application-consistency

<!-- field: preconditions -->
## Preconditions

- Run commands from the repository root.
- Python 3 with PyYAML, Go, and Mage are available.
- applications/catalog and agent-core are present for ownership and runtime validation.
- A planned test is not passing evidence until its named command and expected observations exist.

<!-- field: test_cases -->
## Test cases

<!-- case: Release 14.0 shared documents parse and trace reciprocally | status: implemented -->
### Release 14.0 shared documents parse and trace reciprocally

<!-- field: id -->
#### Id

TC1

<!-- field: use_case -->
#### Use case

rel14.0-uc001-application-consistency

<!-- field: evidence_type -->
#### Evidence type

executable_structural

<!-- field: command -->
#### Command

python3 - <<'PY'
from pathlib import Path
import subprocess
import yaml

root = Path("applications/docs")
paths = {
    "srd": root / "specs/software-requirements/srd003-application-consistency.md",
    "use_case": root / "specs/use-cases/rel14.0-uc001-application-consistency.md",
    "suite": root / "specs/test-suites/test-rel14.0-application-consistency.md",
    "index": root / "SPECIFICATIONS.md",
    "roadmap": root / "road-map.md",
    "architecture": root / "ARCHITECTURE.md",
}
def load(path):
    # A corpus document is Markdown; mage specYAML prints the YAML tree
    # its declaration parses it into (GH-2627).
    out = subprocess.run(["mage", "specYAML", str(path)], cwd="agent-core",
                         check=True, capture_output=True, text=True).stdout
    return yaml.safe_load(out)

docs = {name: load(path) for name, path in paths.items()}
required = {
    "srd": {"id", "title", "problem", "goals", "requirements",
            "non_goals", "acceptance_criteria"},
    "use_case": {"id", "title", "summary", "actor", "trigger", "flow",
                 "touchpoints", "success_criteria", "out_of_scope"},
    "suite": {"id", "title", "release", "traces", "preconditions",
              "test_cases"},
}
for name, fields in required.items():
    assert not fields - set(docs[name]), (name, fields - set(docs[name]))

srd, use_case, suite = docs["srd"], docs["use_case"], docs["suite"]
assert srd["id"] == "srd003-application-consistency"
assert use_case["id"] == "rel14.0-uc001-application-consistency"
assert suite["id"] == "test-rel14.0-application-consistency"
assert use_case["test_suite"] == suite["id"]
assert use_case["id"] in suite["traces"]
assert use_case["id"] in srd["use_case_traces"]
assert suite["id"] in srd["test_suite_traces"]

item_ids = {
    item["id"]
    for group in srd["requirements"].values()
    for item in group["items"]
}
ac_traces = {
    trace
    for ac in srd["acceptance_criteria"]
    for trace in ac["traces"]
}
assert item_ids <= ac_traces, sorted(item_ids - ac_traces)
ac_ids = {ac["id"] for ac in srd["acceptance_criteria"]}
case_traces = {
    trace.rsplit(" ", 1)[-1]
    for case in suite["test_cases"]
    for trace in case.get("traces", [])
    if trace.startswith(srd["id"] + " ")
}
assert ac_ids <= case_traces, sorted(ac_ids - case_traces)

index = docs["index"]
assert any(entry["id"] == srd["id"] for entry in index["srd_index"])
assert any(entry["id"] == use_case["id"] for entry in index["use_case_index"])
assert any(entry["id"] == suite["id"] for entry in index["test_suite_index"])
release = next(r for r in docs["roadmap"]["releases"] if r["version"] == "14.0")
assert any(entry["id"] == use_case["id"] for entry in release["use_cases"])
assert any(entry["id"] == suite["id"] for entry in release["test_suites"])
print("validated Release 14.0 application consistency trace closure")
PY

<!-- field: expected -->
#### Expected

```yaml
yaml_parses: true
design_constitution_fields_present: true
requirement_to_acceptance_criterion_closure: true
reciprocal_srd_use_case_suite_traces: true
specification_and_roadmap_indexes_present: true
```

<!-- field: traces -->
#### Traces

- rel14.0-uc001-application-consistency S1
- rel14.0-uc001-application-consistency S5
- srd003-application-consistency AC1
- srd003-application-consistency AC6

<!-- case: Composition manifests are complete and capability-aware | status: implemented -->
### Composition manifests are complete and capability-aware

<!-- field: id -->
#### Id

TC2

<!-- field: use_case -->
#### Use case

rel14.0-uc001-application-consistency

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: command -->
#### Command

(cd magefiles && go test -run '^TestApplicationConsistencyManifests$')

<!-- field: expected -->
#### Expected

```yaml
all_runnable_and_audit_only_manifests_parse: true
identities_and_statuses_match_roots: true
direct_roots_and_runtime_paths_unique: true
catalog_compatibility_explicit: true
deployment_and_ui_entries_capability_gated: true
absent_optional_sections_accepted: true
```

<!-- field: traces -->
#### Traces

- rel14.0-uc001-application-consistency S1
- srd003-application-consistency AC1

<!-- case: Actor ownership and program names conform | status: implemented -->
### Actor ownership and program names conform

<!-- field: id -->
#### Id

TC3

<!-- field: use_case -->
#### Use case

rel14.0-uc001-application-consistency

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: command -->
#### Command

(cd magefiles && go test -run '^TestApplicationActorGrammar$')

<!-- field: expected -->
#### Expected

```yaml
local_actors_below_actor_roots: true
canonical_implementations_catalog_owned: true
fixtures_excluded_from_production: true
primary_and_operation_names_aligned: true
tools_are_name_only: true
declarations_own_tooldefs: true
request_and_ui_files_required_only_when_present: true
```

<!-- field: traces -->
#### Traces

- rel14.0-uc001-application-consistency S2
- srd003-application-consistency AC2
- srd003-application-consistency AC3

<!-- case: Packaged closure and provenance derive from the manifest | status: implemented -->
### Packaged closure and provenance derive from the manifest

<!-- field: id -->
#### Id

TC4

<!-- field: use_case -->
#### Use case

rel14.0-uc001-application-consistency

<!-- field: evidence_type -->
#### Evidence type

application_tests

<!-- field: applicability -->
#### Applicability

packaged

<!-- field: command -->
#### Command

(cd magefiles && go test -run '^TestManifestDerivedPackageAssets$')
(cd applications/chatbot-mesh && go test ./magefiles -run '^TestManifestDerivedPackageClosure$')
(cd applications/coding-agent && go test ./magefiles -run '^Test(ManifestMutationControlsCollectorAssets|UndeclaredCollectorUIIsNotPackaged|UndeclaredCollectorRootIsNotPackaged)$')
(cd applications/agent-architecture && go test ./magefiles -run '^Test(ManifestMutationControlsCatalogDocsArchive|UndeclaredCatalogDocsAreNotArchived)$')

<!-- field: expected -->
#### Expected

```yaml
manifest_is_only_direct_root_inventory: true
all_reference_classes_traversed: true
generated_inventory_not_authority: true
helm_consumes_generated_closure: true
compatibility_and_actual_revision_distinct: true
deterministic_sorted_checksummed_output: true
```

<!-- field: traces -->
#### Traces

- rel14.0-uc001-application-consistency S3
- srd003-application-consistency AC4

<!-- case: Promotion and shared UI-token claims have real consumers | status: implemented -->
### Promotion and shared UI-token claims have real consumers

<!-- field: id -->
#### Id

TC5

<!-- field: use_case -->
#### Use case

rel14.0-uc001-application-consistency

<!-- field: evidence_type -->
#### Evidence type

go_test_and_ui_build

<!-- field: applicability -->
#### Applicability

promotion candidates or ui

<!-- field: command -->
#### Command

(cd magefiles && go test -run '^TestApplicationPromotionAndUITokens$')
(cd applications/catalog && go test ./conformance -run '^Test(DesignTokens|ApplierPromotion)')

<!-- field: expected -->
#### Expected

```yaml
promotion_has_independent_or_multiple_consumers: true
promoted_source_has_one_owner: true
consumers_reference_canonical_source: true
application_specific_bindings_remain_local: true
shared_token_claims_match_canonical_source: true
```

<!-- field: traces -->
#### Traces

- rel14.0-uc001-application-consistency S4
- srd003-application-consistency AC5

<!-- case: README, status, capabilities, and root lists agree | status: implemented -->
### README, status, capabilities, and root lists agree

<!-- field: id -->
#### Id

TC6

<!-- field: use_case -->
#### Use case

rel14.0-uc001-application-consistency

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: command -->
#### Command

(cd magefiles && go test -run '^TestApplicationReadmeStatusAndRootClassification$')

<!-- field: expected -->
#### Expected

```yaml
common_readme_sections_present: true
status_vocabulary_closed: true
capability_statuses_independent: true
partial_and_dependency_gated_gaps_named: true
implemented_claims_have_evidence: true
root_membership_matches_module_status: true
```

<!-- field: traces -->
#### Traces

- rel14.0-uc001-application-consistency S5
- srd003-application-consistency AC6

<!-- case: Three-application ownership and capability matrix is truthful | status: implemented -->
### Three-application ownership and capability matrix is truthful

<!-- field: id -->
#### Id

TC7

<!-- field: use_case -->
#### Use case

rel14.0-uc001-application-consistency

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: command -->
#### Command

(cd magefiles && go test -run '^TestRelease14ApplicationMatrix$')

<!-- field: expected -->
#### Expected

```yaml
chatbot_mesh_normalized: true
coding_agent_ownership_reconciled: true
agent_architecture_status_reconciled: true
```

<!-- field: traces -->
#### Traces

- rel14.0-uc001-application-consistency S7
- srd003-application-consistency AC1
- srd003-application-consistency AC2
- srd003-application-consistency AC6

<!-- case: Repository audit preserves ownership and current evidence | status: implemented -->
### Repository audit preserves ownership and current evidence

<!-- field: id -->
#### Id

TC8

<!-- field: use_case -->
#### Use case

rel14.0-uc001-application-consistency

<!-- field: evidence_type -->
#### Evidence type

mage_gate

<!-- field: command -->
#### Command

mage audit

<!-- field: expected -->
#### Expected

```yaml
root_audit_green: true
shared_release_14_yaml_has_no_new_warnings: true
catalog_application_runtime_boundaries_preserved: true
completed_normalization_reported_implemented: true
```

<!-- field: traces -->
#### Traces

- rel14.0-uc001-application-consistency S6
- srd003-application-consistency AC7

<!-- field: requirement_coverage -->
## Requirement coverage

```yaml
R1:
    requirement_items: [R1.1, R1.2, R1.3, R1.4, R1.5, R1.6]
    acceptance_criterion: AC1
    use_case_steps: [F1, F2]
    test_cases: [TC1, TC2, TC7]
R2:
    requirement_items: [R2.1, R2.2, R2.3, R2.4]
    acceptance_criterion: AC2
    use_case_steps: [F3]
    test_cases: [TC3, TC7]
R3:
    requirement_items: [R3.1, R3.2, R3.3, R3.4, R3.5]
    acceptance_criterion: AC3
    use_case_steps: [F4]
    test_cases: [TC3]
R4:
    requirement_items: [R4.1, R4.2, R4.3, R4.4]
    acceptance_criterion: AC4
    use_case_steps: [F5]
    test_cases: [TC4]
    applicability: packaged
R5:
    requirement_items: [R5.1, R5.2, R5.3]
    acceptance_criterion: AC5
    use_case_steps: [F6]
    test_cases: [TC5]
R6:
    requirement_items: [R6.1, R6.2, R6.3, R6.4]
    acceptance_criterion: AC6
    use_case_steps: [F1, F7]
    test_cases: [TC1, TC6, TC7]
R7:
    requirement_items: [R7.1, R7.2, R7.3]
    acceptance_criterion: AC7
    use_case_steps: [F1, F8]
    test_cases: [TC8]
```

<!-- field: coverage_notes -->
## Coverage notes

```yaml
implemented_now:
    - TC1 through TC8 are executable and validate the completed Release 14 normalization, reciprocal traces, and repository audit.
remaining:
    - A capability with no declaration remains not_applicable and does not need an empty directory or a passing no-op test.
```
