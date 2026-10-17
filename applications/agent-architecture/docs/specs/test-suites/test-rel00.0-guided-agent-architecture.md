---
id: test-rel00.0-guided-agent-architecture
title: Test Suite for Release 00.0 Guided Agent Architecture
release: "00.0"
status: partial
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

This suite defines executable evidence for the standalone presentation
composition. Deterministic Go tests inspect exact commands, working
directories, root resolution, lifecycle configuration, documentation
reciprocity, and composition-only statistics without starting servers.
Dependency-gated Mage commands provide live launch and presentation evidence.
A managed-service case applies shared srd002 R1 through R3 separately to the
curator and Go present servers. The deterministic command, configuration,
audit, and stats cases are implemented. The aggregate suite is partial
because live service and browser observations remain dependency-gated and
are not inferred from command-construction tests.

<!-- field: live_observation_status -->
## Live observation status

dependency_gated

<!-- field: traces -->
## Traces

- rel00.0-uc001-guided-agent-architecture

<!-- field: preconditions -->
## Preconditions

- The six Release 00.0 application documents parse as YAML.
- The application Go module records a version-pinned Go present tool dependency.
- The demo.yaml catalog_root and core_root may be explicit independent checkouts or their portable monorepo defaults.
- Live evidence has available local ports and a buildable agent-core checkout.
- applications/catalog remains the canonical owner of the documentation-curator profile and UI.

<!-- field: test_cases -->
## Test cases

<!-- case: Pinned presentation command | status: implemented -->
### Pinned presentation command

The deterministic test constructs the presentation command without
starting a server. It proves that Mage invokes the Go-module-pinned
present tool for agent-architecture.slide and uses the application root as
the working directory. The live target serves the deck when its port is
available.

<!-- field: id -->
#### ID

presentation_command

<!-- field: use_case -->
#### Use case

rel00.0-uc001-guided-agent-architecture

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestPresentation'

<!-- field: live_target -->
#### Live target

mage presentation

<!-- field: live_status -->
#### Live status

dependency_gated

<!-- field: inputs | form: "tree" -->
#### Inputs

- application_root: applications/agent-architecture
- slide: agent-architecture.slide
- tool: Go present dependency declared by go.mod

<!-- field: expected | form: "tree" -->
#### Expected

- globally_installed_present_required: false
- slide_path_is_application_local: true
- working_directory_is_application_root: true

<!-- field: traces -->
#### Traces

- rel00.0-uc001-guided-agent-architecture S1

<!-- case: Canonical curator launch command | status: implemented -->
### Canonical curator launch command

The deterministic tests prove declared demo.yaml roots, monorepo
defaults, exact command arguments, catalog-root working directory, and
actionable missing-root failures. The command points to the canonical
catalog profile and passes --profile, --directory, and development
--core-root to agent-core. The live target starts the curator when the
runtime checkout and ports are available.

<!-- field: id -->
#### ID

canonical_curator_command

<!-- field: use_case -->
#### Use case

rel00.0-uc001-guided-agent-architecture

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^Test(Run|ResolveRoots)'

<!-- field: live_target -->
#### Live target

mage run

<!-- field: live_status -->
#### Live status

dependency_gated

<!-- field: inputs | form: "tree" -->
#### Inputs

- catalog_root_config: demo.yaml catalog_root
- core_root_config: demo.yaml core_root
- canonical_profile: agents/knowledge-manager/documentation-curator/profile.yaml
- documentation_workspace: demo.yaml catalog_root

<!-- field: expected | form: "tree" -->
#### Expected

- profile_owner: applications/catalog
- profile_copy_in_application: false
- ui_copy_in_application: false
- explicit_profile_argument: true
- explicit_directory_argument: true
- development_core_root_argument: true

<!-- field: traces -->
#### Traces

- rel00.0-uc001-guided-agent-architecture S2
- rel00.0-uc001-guided-agent-architecture S3
- srd011-knowledge-manager AC1
- srd011-knowledge-manager AC2
- srd001-application-module-contract AC2

<!-- case: Declarative lifecycle-exit contract | status: implemented -->
### Declarative lifecycle-exit contract

The deterministic test parses the application profile and machine,
verifies the profile names the tool, declaration, and REST files, and
inspects the REST definition. It proves the declared REST operation posts
the fixed demo reason to /api/lifecycle/exit, treats HTTP 200 and 202 as
accepted, and transitions to Done. Dependency-gated live evidence
observes the canonical curator reaching Done and releasing listeners.

<!-- field: id -->
#### ID

lifecycle_exit_contract

<!-- field: use_case -->
#### Use case

rel00.0-uc001-guided-agent-architecture

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestLifecycleExit'

<!-- field: live_evidence -->
#### Live evidence

mage run followed by the documented lifecycle-exit command

<!-- field: live_status -->
#### Live status

dependency_gated

<!-- field: inputs | form: "tree" -->
#### Inputs

- profile: ../catalog/agents/lifecycle-exit/profile.yaml
- endpoint: http://127.0.0.1:18082/api/lifecycle/exit
- request_reason: demo presentation

<!-- field: expected | form: "tree" -->
#### Expected

- implementation: declarative_agent
- bespoke_go_http_client: false
- accepted_statuses:
  1. 200
  1. 202
- exit_profile_terminal: succeeded
- live_curator_terminal_observation: dependency_gated
- live_listener_release_observation: dependency_gated

<!-- field: traces -->
#### Traces

- rel00.0-uc001-guided-agent-architecture S4
- srd011-knowledge-manager AC3
- rel07.0-uc001-documentation-curator-profile-ux S4

<!-- case: Curator and presentation managed-service evidence | status: dependency_gated -->
### Curator and presentation managed-service evidence

The deterministic test classifies both long-running servers under shared
srd002 R1 through R3 and verifies that each records an owned deterministic
command/configuration check plus planned live observations. It does not
execute startup, health, failure, telemetry, configuration-rejection,
shutdown, or resource-release behavior. Dependency-gated evidence must
probe each declared role, exercise representative failures, inspect
bounded diagnostics and available telemetry, then stop each process and
observe in-flight outcomes and listener release. The present server has
no dedicated health or telemetry API, so evidence must record that signal
limitation and cannot treat process existence as readiness. This case
does not satisfy srd002 AC1 through AC3 until those observations run.

<!-- field: id -->
#### ID

managed_service_contract

<!-- field: use_case -->
#### Use case

rel00.0-uc001-guided-agent-architecture

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestManagedServices'

<!-- field: live_evidence -->
#### Live evidence

bounded launch, readiness, failure, interrupt, and listener-release checks

<!-- field: live_status -->
#### Live status

dependency_gated

<!-- field: inputs | form: "tree" -->
#### Inputs

- services:
  1. documentation-curator
  1. Go present
- shared_contract: srd002-managed-application-services
- applicable_requirement_groups:
  1. R1
  1. R2
  1. R3

<!-- field: expected | form: "tree" -->
#### Expected

- deterministic_service_classification: true
- deterministic_operator_contract: true
- curator_role_readiness_observation: dependency_gated
- curator_declared_exit_observation: dependency_gated
- present_http_readiness_observation: dependency_gated
- present_interrupt_release_observation: dependency_gated
- process_existence_used_as_readiness: false
- packaged_capability_inferred: false
- helm_capability_inferred: false
- kind_capability_inferred: false

<!-- field: traces -->
#### Traces

- rel00.0-uc001-guided-agent-architecture S6
- srd002-managed-application-services AC1
- srd002-managed-application-services AC2
- srd002-managed-application-services AC3

<!-- case: Documentation audit and composition-only stats | status: implemented -->
### Documentation audit and composition-only stats

The application audit parses the six documentation YAML files, checks
their required fields and reciprocal use-case/suite traces, scans
documentation and lifecycle files for stale or developer-specific paths,
and confirms profile references use the canonical path. Lifecycle YAML
has separate profile/machine/REST coverage in lifecycle_exit_contract;
audit does not claim to parse every lifecycle YAML file. Stats emits only
the application composition fields and no agents section.

<!-- field: id -->
#### ID

application_contract

<!-- field: use_case -->
#### Use case

rel00.0-uc001-guided-agent-architecture

<!-- field: go_test -->
#### Go test

go test ./magefiles

<!-- field: audit_target -->
#### Audit target

mage audit

<!-- field: stats_target -->
#### Stats target

mage stats

<!-- field: inputs | form: "tree" -->
#### Inputs

- documentation_yaml: docs/VISION.md, docs/ARCHITECTURE.md, docs/road-map.md, docs/SPECIFICATIONS.md, the indexed use case, and the indexed test suite
- lifecycle_yaml: ../catalog/agents/lifecycle-exit/profile.yaml, ../catalog/agents/lifecycle-exit/machine.yaml, and ../catalog/agents/lifecycle-exit/rest.yaml are inspected by lifecycle_exit_contract
- canonical_profile_owner: applications/catalog
- expected_canonical_references: 2

<!-- field: expected | form: "tree" -->
#### Expected

- yaml_parse_errors: 0
- reciprocal_use_case_traces: true
- portable_paths: true
- ownership: composition-only
- agents_contributed: 0
- canonical_references: 2
- agents_section_present: false

<!-- field: traces -->
#### Traces

- rel00.0-uc001-guided-agent-architecture S5
- srd001-application-module-contract AC3
- srd001-application-module-contract AC4
- srd001-application-module-contract AC5
- srd001-application-module-contract AC6
