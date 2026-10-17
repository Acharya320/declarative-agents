---
id: test-rel06.0-profile-migration
title: Test Suite for Release 06.0 Profile Migration
release: "06.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

This suite is now the Release 06.0 migration status record. Through mage
audit, the SRD, use case, and suite parse and appear in the specification
index.

Implemented cases name the Go tests, Mage targets, Docker target, and
integration targets that prove external profile roots, installed core-owned
tool paths, runtime image profile exclusion, and mounted-profile execution.
Remaining cleanup work belongs to stale documentation references, not to the
container or resolver proof itself.

<!-- field: traces -->
## Traces

- [rel06.0-uc001-mounted-profile-container-runtime](../use-cases/rel06.0-uc001-mounted-profile-container-runtime.md)

<!-- field: preconditions -->
## Preconditions

- applications/catalog exists as a checkout, release asset, or mounted fixture with standard agent profiles.
- The Agent Core runtime image can be built from Dockerfile.
- Core-owned shared tool declaration roots are available from the installed Agent Core runtime path.
- Tests can mount an isolated workspace at /work.

<!-- field: coverage_notes -->
## Coverage notes

Release 06.0 has moved past draft status.

Agent profile assets, demos, and profile-owned integration fixtures now live
in applications/catalog. Agent Core keeps cmd/agent, Go runtime packages, shared
tool declarations, Docker packaging, Mage targets, and the release image
contract.

Executable proof is split by boundary. Package tests verify installed
core-owned path mapping and Dockerfile contents. Through AGENT_CATALOG_ROOT,
Mage integration targets run external profiles. For runtime images, mage
docker builds a target that copies /opt/agent-core/tools, refuses
profile-less startup, and documents the mounted /profiles command.
applications/catalog owns the matching mounted-image smoke target for its profile
bundle.

<!-- field: test_cases -->
## Test cases

<!-- case: Specification index includes external profile contracts | status: implemented -->
### Specification index includes external profile contracts

Parse the three Release 06.0 YAML documents through the specification
audit. After migration, the index and road map include the SRD, use case,
and suite.

<!-- field: use_case -->
#### Use case

[rel06.0-uc001-mounted-profile-container-runtime](../use-cases/rel06.0-uc001-mounted-profile-container-runtime.md)

<!-- field: go_test -->
#### Go test

mage audit

<!-- field: inputs | form: "tree" -->
#### Inputs

- srd: docs/specs/software-requirements/srd034-external-agent-profiles.md
- use_case: docs/specs/use-cases/rel06.0-uc001-mounted-profile-container-runtime.md
- test_suite: docs/specs/test-suites/test-rel06.0-profile-migration.md

<!-- field: expected | form: "tree" -->
#### Expected

- yaml_parses: true
- specifications_indexed: true
- roadmap_indexed: true
- architecture_contract_updated: true
- release_status_current: true

<!-- field: traces -->
#### Traces

- [rel06.0-uc001-mounted-profile-container-runtime](../use-cases/rel06.0-uc001-mounted-profile-container-runtime.md) [S1](../use-cases/rel06.0-uc001-mounted-profile-container-runtime.md#s1)
- [srd034-external-agent-profiles](../software-requirements/srd034-external-agent-profiles.md) [AC1](../software-requirements/srd034-external-agent-profiles.md#ac1)

<!-- case: Runtime resolves profiles from external roots | status: implemented -->
### Runtime resolves profiles from external roots

Configure AGENT_CATALOG_ROOT outside the agent-core checkout. The coding
application packages the external profile closure and starts cmd/agent
with a profile from that closure.
In pkg/spec, the path test verifies that /opt/agent-core/tools references
map to the installed core asset root when AGENT_CORE_HOME is set. Those
The application integration exercises normal machine assets and fixture roots.
Concrete profile UX tracer bullets, including Knowledge Manager assets,
belong to applications/catalog suites.

<!-- field: use_case -->
#### Use case

[rel06.0-uc001-mounted-profile-container-runtime](../use-cases/rel06.0-uc001-mounted-profile-container-runtime.md)

<!-- field: integration_target -->
#### Integration target

applications/coding-agent mage integration:executorLive

<!-- field: go_test -->
#### Go test

TestResolveProfilePathMapsInstalledCoreHome

<!-- field: inputs | form: "tree" -->
#### Inputs

- profile_root: external_agent_profiles_root
- profiles:
  1. ${AGENT_CATALOG_ROOT}/generator/profile.yaml
  1. ${AGENT_CATALOG_ROOT}/evaluator/profile.yaml
- workspace: isolated_temp_workspace

<!-- field: expected | form: "tree" -->
#### Expected

- external_profile_loaded: true
- profile_relative_paths_resolved: true
- core_tool_declarations_resolved: true
- installed_core_tool_root_mapped: true
- profile_owned_fixtures_consumed: true

<!-- field: traces -->
#### Traces

- [rel06.0-uc001-mounted-profile-container-runtime](../use-cases/rel06.0-uc001-mounted-profile-container-runtime.md) [S2](../use-cases/rel06.0-uc001-mounted-profile-container-runtime.md#s2)
- [rel06.0-uc001-mounted-profile-container-runtime](../use-cases/rel06.0-uc001-mounted-profile-container-runtime.md) [S3](../use-cases/rel06.0-uc001-mounted-profile-container-runtime.md#s3)
- [rel06.0-uc001-mounted-profile-container-runtime](../use-cases/rel06.0-uc001-mounted-profile-container-runtime.md) [S4](../use-cases/rel06.0-uc001-mounted-profile-container-runtime.md#s4)
- [srd034-external-agent-profiles](../software-requirements/srd034-external-agent-profiles.md) [AC2](../software-requirements/srd034-external-agent-profiles.md#ac2)

<!-- case: Runtime image excludes migrated profile assets | status: implemented -->
### Runtime image excludes migrated profile assets

Inspect Dockerfile and build-command behavior. The runtime stage copies
the agent binary and core-owned shared declaration roots under
/opt/agent-core/tools. During the Dockerfile test, standard profile
directories such as generator, evaluator, and Knowledge Manager must stay
absent from runtime-image copy rules. The mage docker target is the
release-image build wrapper for the same Dockerfile.

<!-- field: use_case -->
#### Use case

[rel06.0-uc001-mounted-profile-container-runtime](../use-cases/rel06.0-uc001-mounted-profile-container-runtime.md)

<!-- field: integration_target -->
#### Integration target

mage docker

<!-- field: go_test -->
#### Go test

cd magefiles && go test ./... -run 'TestDockerfileRuntimeExcludesAgentProfiles'

<!-- field: inputs | form: "tree" -->
#### Inputs

- image: agent-core:latest
- expected_core_paths:
  1. /usr/local/bin/agent
  1. /opt/agent-core/tools
- forbidden_profile_paths:
  1. migrated_generator_profile_dir
  1. migrated_evaluator_profile_dir
  1. migrated_knowledge_manager_profile_dir

<!-- field: expected | form: "tree" -->
#### Expected

- agent_binary_present: true
- core_tool_assets_present: true
- profile_directories_absent: true
- profile_required_entrypoint_present: true
- integration_stage_defined: TestDockerfileDefinesIntegrationTarget

<!-- field: traces -->
#### Traces

- [rel06.0-uc001-mounted-profile-container-runtime](../use-cases/rel06.0-uc001-mounted-profile-container-runtime.md) [S1](../use-cases/rel06.0-uc001-mounted-profile-container-runtime.md#s1)
- [srd034-external-agent-profiles](../software-requirements/srd034-external-agent-profiles.md) [AC3](../software-requirements/srd034-external-agent-profiles.md#ac3)

<!-- case: Container runs with mounted profile and workspace | status: implemented -->
### Container runs with mounted profile and workspace

Build output documents the mounted profile command. The runtime command
mounts AGENT_CATALOG_ROOT at /profiles/agents, mounts the workspace at
/work, and invokes agent with --profile pointing at the mounted generator
profile. Inside the integration Docker target, Mage runs with
AGENT_CATALOG_ROOT mounted. In applications/catalog, the matching container
smoke command runs an installed image against /profiles.

<!-- field: use_case -->
#### Use case

[rel06.0-uc001-mounted-profile-container-runtime](../use-cases/rel06.0-uc001-mounted-profile-container-runtime.md)

<!-- field: integration_target -->
#### Integration target

mage docker; Dockerfile target integration; applications/catalog mage containerSmoke

<!-- field: go_test -->
#### Go test

cd magefiles && go test ./... -run 'TestContainerBuildSummaryForDocker'

<!-- field: inputs | form: "tree" -->
#### Inputs

- profile_mount: /profiles
- workspace_mount: /work
- command:
  1. agent
  1. --profile
  1. /profiles/agents/<agent>/profile.yaml
  1. --directory
  1. /work

<!-- field: expected | form: "tree" -->
#### Expected

- mounted_profile_used: true
- workspace_used: true
- image_profile_copy_required: false
- command_documents_profile_path: true
- integration_target_mounts_profiles: true

<!-- field: traces -->
#### Traces

- [rel06.0-uc001-mounted-profile-container-runtime](../use-cases/rel06.0-uc001-mounted-profile-container-runtime.md) [S2](../use-cases/rel06.0-uc001-mounted-profile-container-runtime.md#s2)
- [srd034-external-agent-profiles](../software-requirements/srd034-external-agent-profiles.md) [AC4](../software-requirements/srd034-external-agent-profiles.md#ac4)
