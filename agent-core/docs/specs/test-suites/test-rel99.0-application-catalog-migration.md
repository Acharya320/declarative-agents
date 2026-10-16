---
id: test-rel99.0-application-catalog-migration
title: Test Suite for Release 99.0 Application and Catalog Migration
release: "99.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel99.0-uc002-application-catalog-migration

<!-- field: preconditions -->
## Preconditions

- srd043-application-catalog-layout is indexed before migration code lands.
- The repository root contains agent-core plus the source trees being moved.
- Each command runs from the owner root named in the test case.
- Live integrations retain their existing prerequisite checks and skip semantics.
- Generated node_modules, build output, Helm archives, and unrelated untracked files are excluded from migration diffs.

<!-- field: test_cases -->
## Test cases

<!-- case: Migration contract is indexed and complete | status: specified -->
### Migration contract is indexed and complete

Parse and audit the SRD, use case, and suite. Confirm the roadmap and
specification index include all three artifacts and all fourteen GH-1036 child
issues map to requirements and named cases.

<!-- field: use_case -->
#### Use case

rel99.0-uc002-application-catalog-migration

<!-- field: issue -->
#### Issue

GH-1040

<!-- field: go_test -->
#### Go test

mage audit

<!-- field: inputs -->
#### Inputs

```yaml
owner_root: repository_root
artifacts:
    - agent-core/docs/specs/software-requirements/srd043-application-catalog-layout.md
    - agent-core/docs/specs/use-cases/rel99.0-uc002-application-catalog-migration.md
    - agent-core/docs/specs/test-suites/test-rel99.0-application-catalog-migration.md
```

<!-- field: expected -->
#### Expected

```yaml
yaml_parses: true
artifacts_indexed: true
child_issue_count: 14
ownership_decisions_complete: true
compatibility_decisions_complete: true
```

<!-- field: evidence -->
#### Evidence

- mage audit output with no new specification errors
- Issue mapping in srd043 implementation_plan and this suite

<!-- field: traces -->
#### Traces

- srd043-application-catalog-layout AC1
- srd043-application-catalog-layout AC4
- srd043-application-catalog-layout AC5
- srd043-application-catalog-layout AC6
- srd043-application-catalog-layout AC8

<!-- case: Runnable application source trees move mechanically | status: implemented -->
### Runnable application source trees move mechanically

Move both runnable application directories, update root participation and
module declarations, and preserve their relative root Mage replacement.

<!-- field: use_case -->
#### Use case

rel99.0-uc002-application-catalog-migration

<!-- field: issue -->
#### Issue

GH-1037

<!-- field: go_test -->
#### Go test

mage audit

<!-- field: integration_target -->
#### Integration target

root Mage module-list regression; go test ./... in applications/coding-agent and applications/chatbot-mesh

<!-- field: inputs -->
#### Inputs

```yaml
owner_root: repository_root
historical_old_paths: [examples/coding-agent, examples/chatbot-mesh]
new_paths: [applications/coding-agent, applications/chatbot-mesh]
```

<!-- field: expected -->
#### Expected

```yaml
old_active_source_paths_absent: true
applications_discovered: true
module_paths_updated: true
relative_mage_replacements_valid: true
```

<!-- field: evidence -->
#### Evidence

- Root module-list regression output
- Unit-test output from both application module roots

<!-- field: traces -->
#### Traces

- srd043-application-catalog-layout AC1
- srd043-application-catalog-layout AC2
- srd043-application-catalog-layout AC5
- srd043-application-catalog-layout AC7

<!-- case: Core integrations use only core-owned profiles | status: implemented -->
### Core integrations use only core-owned profiles

Run service integrations from agent-core using minimal profiles under
agent-core/testdata/integration. The targets keep live-service preflight
behavior but never resolve a catalog root.

<!-- field: use_case -->
#### Use case

rel99.0-uc002-application-catalog-migration

<!-- field: issue -->
#### Issue

GH-1038

<!-- field: go_test -->
#### Go test

mage audit

<!-- field: integration_target -->
#### Integration target

cd magefiles && go test ./...; mage integration:monitor; mage integration:ollamaRest; mage integration:ollamaMonitor; mage integration:dolt; mage integration:doltWord

<!-- field: inputs -->
#### Inputs

```yaml
owner_root: agent-core
profile_source: agent-core/testdata/integration
forbidden_inputs: [AGENT_CATALOG_ROOT, AGENT_PROFILES_ROOT, applications/catalog]
```

<!-- field: expected -->
#### Expected

```yaml
core_profiles_loaded: true
catalog_discovery_called: false
live_skip_semantics_preserved: true
process_cwd_changed: false
```

<!-- field: evidence -->
#### Evidence

- Focused Mage package test output
- integration:all pass and skip summary from agent-core

<!-- field: traces -->
#### Traces

- srd043-application-catalog-layout AC2
- srd043-application-catalog-layout AC3
- srd043-application-catalog-layout AC7

<!-- case: Core tests and audit require no catalog | status: implemented -->
### Core tests and audit require no catalog

Run Agent Core's full Go and audit gates in a checkout with no
applications/catalog. Control, lifecycle, bootstrap, monitor, and Jurist
runtime coverage uses only core-owned fixtures; catalog membership checks
execute in catalog conformance instead.

<!-- field: use_case -->
#### Use case

rel99.0-uc002-application-catalog-migration

<!-- field: issue -->
#### Issue

GH-1039

<!-- field: go_test -->
#### Go test

go test ./...

<!-- field: integration_target -->
#### Integration target

mage audit

<!-- field: inputs -->
#### Inputs

```yaml
owner_root: agent-core
catalog_present: false
environment_unset: [AGENT_CATALOG_ROOT, AGENT_PROFILES_ROOT]
```

<!-- field: expected -->
#### Expected

```yaml
go_tests_pass: true
audit_passes: true
missing_catalog_skips: 0
catalog_membership_assertions_in_core: false
```

<!-- field: evidence -->
#### Evidence

- agent-core go test ./... output
- agent-core mage audit output
- Catalog conformance test names replacing moved assertions

<!-- field: traces -->
#### Traces

- srd043-application-catalog-layout AC1
- srd043-application-catalog-layout AC3
- srd043-application-catalog-layout AC7

<!-- case: Catalog module relocates without semantic reorganization | status: implemented -->
### Catalog module relocates without semantic reorganization

Move the complete agent-profiles module to applications/catalog, update
its module identity and purpose language, and preserve internal files and
executable behavior until the later taxonomy issue.

<!-- field: use_case -->
#### Use case

rel99.0-uc002-application-catalog-migration

<!-- field: issue -->
#### Issue

GH-1041

<!-- field: go_test -->
#### Go test

go test ./...

<!-- field: integration_target -->
#### Integration target

mage validate; mage conformance; mage integration:all

<!-- field: inputs -->
#### Inputs

```yaml
owner_root: applications/catalog
old_path: agent-profiles
new_path: applications/catalog
```

<!-- field: expected -->
#### Expected

```yaml
complete_tree_moved: true
old_active_source_path_absent: true
internal_reorganization_performed: false
catalog_gates_pass: true
generated_dependency_trees_moved: false
```

<!-- field: evidence -->
#### Evidence

- Git rename inventory
- Catalog Go, validation, conformance, and integration output

<!-- field: traces -->
#### Traces

- srd043-application-catalog-layout AC1
- srd043-application-catalog-layout AC5
- srd043-application-catalog-layout AC7

<!-- case: Catalog root and legacy alias resolve deterministically | status: implemented -->
### Catalog root and legacy alias resolve deterministically

Exercise canonical input, the legacy-only alias, equal dual inputs,
conflicting dual inputs, applications/catalog repository discovery,
installed-runtime behavior, relative-path normalization, and missing-root
diagnostics without changing process CWD.

<!-- field: use_case -->
#### Use case

rel99.0-uc002-application-catalog-migration

<!-- field: issue -->
#### Issue

GH-1043

<!-- field: go_test -->
#### Go test

mage audit

<!-- field: integration_target -->
#### Integration target

focused internal/support/profiles and magefiles catalog-root tests

<!-- field: inputs -->
#### Inputs

```yaml
owner_root: agent-core
cases:
    - AGENT_CATALOG_ROOT_only
    - AGENT_PROFILES_ROOT_only
    - both_same_root
    - both_conflicting_roots
    - repository_discovery
    - installed_runtime
    - missing_catalog
```

<!-- field: expected -->
#### Expected

```yaml
canonical_env_authoritative: true
legacy_only_supported_in_rel99: true
same_dual_input_accepted: true
conflicting_input_rejected: true
absolute_root_resolved_once: true
process_cwd_changed: false
diagnostic_names_canonical_env: true
```

<!-- field: evidence -->
#### Evidence

- Resolver unit-test matrix
- Diagnostic assertions naming both conflicting variables

<!-- field: traces -->
#### Traces

- srd043-application-catalog-layout AC2
- srd043-application-catalog-layout AC4

<!-- case: Root gates tags UIs and stats use applications | status: implemented -->
### Root gates tags UIs and stats use applications

Run root orchestration against agent-core and applications. Verify release
gate order, fail-fast behavior, UI scanning, stable stats keys, and
canonical plus compatibility catalog v0 tags.

<!-- field: use_case -->
#### Use case

rel99.0-uc002-application-catalog-migration

<!-- field: issue -->
#### Issue

GH-1046

<!-- field: go_test -->
#### Go test

mage audit

<!-- field: integration_target -->
#### Integration target

root magefiles tests; mage audit; mage test; mage stats

<!-- field: inputs -->
#### Inputs

```yaml
owner_root: repository_root
canonical_tag: applications/catalog/v0.YYYYMMDD.N
compatibility_tag: agent-profiles/v0.YYYYMMDD.N
```

<!-- field: expected -->
#### Expected

```yaml
root_paths_use_applications: true
child_processes_use_owner_roots: true
catalog_gates_ordered: true
application_uis_scanned: true
catalog_and_application_stats_distinct: true
v0_catalog_tags_share_commit: true
release_fail_fast_preserved: true
```

<!-- field: evidence -->
#### Evidence

- Root Mage package tests
- Root audit, test, and stats output
- Atomic tag-set unit-test output

<!-- field: traces -->
#### Traces

- srd043-application-catalog-layout AC2
- srd043-application-catalog-layout AC5
- srd043-application-catalog-layout AC6
- srd043-application-catalog-layout AC7

<!-- case: Coding agent packages catalog closures with provenance | status: implemented -->
### Coding agent packages catalog closures with provenance

Resolve planner, executor, and critic blocks through AGENT_CATALOG_ROOT
from the coding-agent root. Preserve role shards and /profiles mount
destinations, accept exact canonical or legacy v0 compatibility tags, and
distinguish compatible release from actual source provenance.

<!-- field: use_case -->
#### Use case

rel99.0-uc002-application-catalog-migration

<!-- field: issue -->
#### Issue

GH-1044

<!-- field: go_test -->
#### Go test

mage audit

<!-- field: integration_target -->
#### Integration target

coding-agent magefiles tests; mage package; mage packageValidate; mage integration:codingLoop; mage integration:servingHealth

<!-- field: inputs -->
#### Inputs

```yaml
owner_root: applications/coding-agent
catalog_root: applications/catalog
runtime_mount: /profiles
tag_prefixes: [applications/catalog/v0.*, agent-profiles/v0.*]
```

<!-- field: expected -->
#### Expected

```yaml
catalog_closure_exact: true
role_shards_preserved: true
runtime_mount_preserved: true
checkout_not_mislabeled_release: true
process_cwd_changed: false
```

<!-- field: evidence -->
#### Evidence

- Coding-agent package and packageValidate output
- Closure manifest and provenance assertions
- Serving and coding-loop integration output

<!-- field: traces -->
#### Traces

- srd043-application-catalog-layout AC1
- srd043-application-catalog-layout AC2
- srd043-application-catalog-layout AC5
- srd043-application-catalog-layout AC6
- srd043-application-catalog-layout AC7

<!-- case: Chatbot mesh resolves catalog blocks from its owner root | status: implemented -->
### Chatbot mesh resolves catalog blocks from its owner root

Resolve Jurist, corpus-ingest, assembler, and mock catalog blocks through
AGENT_CATALOG_ROOT while preserving chatbot-mesh-owned profiles, topology,
chart paths, packaged /profiles closure, and integration semantics.

<!-- field: use_case -->
#### Use case

rel99.0-uc002-application-catalog-migration

<!-- field: issue -->
#### Issue

GH-1042

<!-- field: go_test -->
#### Go test

mage audit

<!-- field: integration_target -->
#### Integration target

chatbot-mesh magefiles tests; mage audit; mage helm:package; mage integration:rig; mage integration:helmSmoke

<!-- field: inputs -->
#### Inputs

```yaml
owner_root: applications/chatbot-mesh
catalog_root: applications/catalog
runtime_mount: /profiles
```

<!-- field: expected -->
#### Expected

```yaml
reusable_blocks_from_catalog: true
application_agents_remain_local: true
helm_package_self_contained: true
stats_ownership_distinct: true
process_cwd_changed: false
```

<!-- field: evidence -->
#### Evidence

- Chatbot-mesh Mage package tests
- Audit and Helm package output
- Rig and Helm smoke evidence

<!-- field: traces -->
#### Traces

- srd043-application-catalog-layout AC1
- srd043-application-catalog-layout AC2
- srd043-application-catalog-layout AC5
- srd043-application-catalog-layout AC7

<!-- case: Obsolete core catalog discovery is absent | status: implemented -->
### Obsolete core catalog discovery is absent

Verify Agent Core tests, audit, and integrations use only core-owned
fixtures and contain no catalog-root discovery or product-membership
assertions.

<!-- field: use_case -->
#### Use case

rel99.0-uc002-application-catalog-migration

<!-- field: issue -->
#### Issue

GH-1047

<!-- field: go_test -->
#### Go test

go test ./...

<!-- field: integration_target -->
#### Integration target

mage audit; mage integration:all

<!-- field: inputs -->
#### Inputs

```yaml
owner_root: agent-core
environment_unset: [AGENT_CATALOG_ROOT, AGENT_PROFILES_ROOT]
```

<!-- field: expected -->
#### Expected

```yaml
catalog_discovery_called: false
external_product_assertions: false
core_profile_validation_preserved: true
```

<!-- field: evidence -->
#### Evidence

- Agent Core Go, audit, and integration output
- Focused root and profile-asset regression tests

<!-- field: traces -->
#### Traces

- srd043-application-catalog-layout AC1
- srd043-application-catalog-layout AC2
- srd043-application-catalog-layout AC3

<!-- case: Root module registry uses applications | status: implemented -->
### Root module registry uses applications

Verify root orchestration dispatches the catalog and both runnable
applications from their owner roots and uses application vocabulary.

<!-- field: use_case -->
#### Use case

rel99.0-uc002-application-catalog-migration

<!-- field: issue -->
#### Issue

GH-1048

<!-- field: go_test -->
#### Go test

cd magefiles && go test ./...

<!-- field: integration_target -->
#### Integration target

mage audit; mage test

<!-- field: inputs -->
#### Inputs

```yaml
owner_root: repository_root
modules: [agent-core, applications/catalog, applications/coding-agent, applications/chatbot-mesh]
```

<!-- field: expected -->
#### Expected

```yaml
current_module_paths_only: true
application_vocabulary_only: true
owner_root_dispatch: true
```

<!-- field: evidence -->
#### Evidence

- Root Mage package tests
- Root audit and test output

<!-- field: traces -->
#### Traces

- srd043-application-catalog-layout AC1
- srd043-application-catalog-layout AC2
- srd043-application-catalog-layout AC7

<!-- case: Chatbot integrations use the canonical catalog root | status: implemented -->
### Chatbot integrations use the canonical catalog root

Verify chatbot integrations resolve reusable blocks through
AGENT_CATALOG_ROOT without changing process CWD or source labels.

<!-- field: use_case -->
#### Use case

rel99.0-uc002-application-catalog-migration

<!-- field: issue -->
#### Issue

GH-1049

<!-- field: go_test -->
#### Go test

go test ./...

<!-- field: integration_target -->
#### Integration target

mage audit; mage integration:chroma; mage integration:rig; mage integration:embeddingExclusion

<!-- field: inputs -->
#### Inputs

```yaml
owner_root: applications/chatbot-mesh
catalog_root: applications/catalog
```

<!-- field: expected -->
#### Expected

```yaml
canonical_root_propagated: true
legacy_alias_precedence_tested: true
process_cwd_changed: false
stale_source_labels: 0
```

<!-- field: evidence -->
#### Evidence

- Chatbot Mage package and focused integration output

<!-- field: traces -->
#### Traces

- srd043-application-catalog-layout AC1
- srd043-application-catalog-layout AC2
- srd043-application-catalog-layout AC4
- srd043-application-catalog-layout AC5

<!-- case: Chatbot checkpoint database initializes idempotently | status: implemented -->
### Chatbot checkpoint database initializes idempotently

Verify the chatbot Helm deployment creates or selects its configured Dolt
checkpoint database before agents save checkpoints and repeated
initialization succeeds.

<!-- field: use_case -->
#### Use case

rel99.0-uc002-application-catalog-migration

<!-- field: issue -->
#### Issue

GH-1050

<!-- field: go_test -->
#### Go test

go test ./...

<!-- field: integration_target -->
#### Integration target

mage helm:package; mage integration:helmSmoke

<!-- field: inputs -->
#### Inputs

```yaml
owner_root: applications/chatbot-mesh
```

<!-- field: expected -->
#### Expected

```yaml
checkpoint_database_ready: true
initialization_idempotent: true
runtime_mounts_unchanged: true
```

<!-- field: evidence -->
#### Evidence

- Helm render and package tests
- Chatbot Helm smoke output

<!-- field: traces -->
#### Traces

- srd043-application-catalog-layout AC2
- srd043-application-catalog-layout AC5
- srd043-application-catalog-layout AC7

<!-- case: Taxonomy and active path references match the new ownership | status: implemented -->
### Taxonomy and active path references match the new ownership

Document core implementations and declarations, catalog tool blocks,
catalog agent blocks, tool selections, profile-local overrides, and
application-local composition. Sweep active path claims and command
examples after all moves complete.

<!-- field: use_case -->
#### Use case

rel99.0-uc002-application-catalog-migration

<!-- field: issue -->
#### Issue

GH-1045

<!-- field: go_test -->
#### Go test

mage audit

<!-- field: integration_target -->
#### Integration target

mage audit in repository root, agent-core, applications/catalog, applications/coding-agent, and applications/chatbot-mesh

<!-- field: inputs -->
#### Inputs

```yaml
owner_roots:
    - repository_root
    - agent-core
    - applications/catalog
    - applications/coding-agent
    - applications/chatbot-mesh
forbidden_active_source_paths: [examples/, top-level_agent-profiles/]
```

<!-- field: expected -->
#### Expected

```yaml
tool_first_taxonomy_documented: true
ownership_direction_documented: true
migration_guide_complete: true
stale_active_source_claims: 0
commands_use_owner_root_or_explicit_root: true
```

<!-- field: evidence -->
#### Evidence

- Repository-wide path-reference search
- Audit output from all five owner roots
- applications/catalog/tools/README.md taxonomy

<!-- field: traces -->
#### Traces

- srd043-application-catalog-layout AC1
- srd043-application-catalog-layout AC2
- srd043-application-catalog-layout AC4
- srd043-application-catalog-layout AC6
- srd043-application-catalog-layout AC7
- srd043-application-catalog-layout AC8
