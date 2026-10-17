---
id: test-rel12.0-agent-role-realization-alignment
title: Test Suite for Release 12.0 Agent Role-Realization Alignment
release: "12.0"
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

Deterministic release-integration evidence for the completed naming
migration. The suite validates aliases and durations, exact release metadata,
package closures, Helm staging and generated renders, consumer provenance,
and scoped zero-hit policies. Tag publication remains a post-merge operation.

<!-- field: traces -->
## Traces

- rel12.0-uc001-agent-role-realization-alignment

<!-- field: preconditions -->
## Preconditions

- Issues 1068 through 1070 and 1072 through 1077 are complete.
- The shared role-realization semantic model and executable validator are present.
- Catalog releases 07.1 and 11.0 retain their independent in-progress evidence gaps.

<!-- field: test_cases -->
## Test cases

<!-- case: Migration table and non-forking aliases | status: implemented -->
### Migration table and non-forking aliases

<!-- field: go_test -->
#### Go test

go test ./conformance -run 'Test(CatalogMembershipUsesSharedRealizationAndAliasAuthority|AlignmentMigrationReleaseAndConsumerPins)'

<!-- field: expected | form: "tree" -->
#### Expected

- every_rename_mapped: true
- v0_duration_and_v1_removal_recorded: true
- wrappers_reference_canonical_assets: true

<!-- field: traces -->
#### Traces

- rel12.0-uc001-agent-role-realization-alignment S1
- rel12.0-uc001-agent-role-realization-alignment S2

<!-- case: Exact catalog tag pair and immutability | status: implemented -->
### Exact catalog tag pair and immutability

<!-- field: go_test -->
#### Go test

cd ../../magefiles && go test . -run 'Test(AlignmentReleaseCreatesExactCatalogCompatibilityPair|GitCreateTagSetIsAtomicOnConflict|CreateReleaseTagRejectsNonMainBranch)'

<!-- field: expected | form: "tree" -->
#### Expected

- root_canonical_and_legacy_tags_share_exact_suffix: true
- tag_set_is_atomic: true
- feature_branch_publication_rejected: true

<!-- field: traces -->
#### Traces

- rel12.0-uc001-agent-role-realization-alignment S5

<!-- case: Coding-agent closure and provenance | status: implemented -->
### Coding-agent closure and provenance

<!-- field: integration_target -->
#### Integration target

cd ../coding-agent && mage audit

<!-- field: go_test -->
#### Go test

cd ../coding-agent && go test ./magefiles

<!-- field: expected | form: "tree" -->
#### Expected

- exact_compatible_release: applications/catalog/v0.20260727.0
- closure_contains_every_canonical_reference: true
- checkout_does_not_claim_release_provenance: true

<!-- field: traces -->
#### Traces

- rel12.0-uc001-agent-role-realization-alignment S3

<!-- case: Chatbot-mesh generated Helm package | status: implemented -->
### Chatbot-mesh generated Helm package

<!-- field: integration_target -->
#### Integration target

cd ../chatbot-mesh && mage helm:package

<!-- field: go_test -->
#### Go test

cd ../chatbot-mesh && go test ./magefiles

<!-- field: expected | form: "tree" -->
#### Expected

- exact_catalog_annotation: applications/catalog/v0.20260727.0
- canonical_profile_paths_staged: true
- old_coordinator_paths_absent: true
- generated_source_inventory_has_no_drift: true
- staged_and_archived_charts_lint_and_render: true

<!-- field: traces -->
#### Traces

- rel12.0-uc001-agent-role-realization-alignment S3

<!-- case: Scoped stale-identity policy | status: implemented -->
### Scoped stale-identity policy

<!-- field: go_test -->
#### Go test

go test ./conformance -run 'Test(NoActiveAssemblerPathConsumers|NoActiveMonitorProfilePathConsumers|NoActiveJuristProfilePathConsumers|AlignmentHasNoStaleActiveRuntimeIdentities)'

<!-- field: expected | form: "tree" -->
#### Expected

- stale_active_paths: 0
- stale_active_runtime_identities: 0
- compatibility_and_history_allowlists_exact: true

<!-- field: traces -->
#### Traces

- rel12.0-uc001-agent-role-realization-alignment S4

<!-- case: Full coordinated gates | status: implemented -->
### Full coordinated gates

<!-- field: commands -->
#### Commands

- mage validate
- mage conformance
- mage integration:all
- cd ../.. && mage audit
- cd ../.. && mage test
- cd ../.. && mage stats

<!-- field: traces -->
#### Traces

- rel12.0-uc001-agent-role-realization-alignment S6
