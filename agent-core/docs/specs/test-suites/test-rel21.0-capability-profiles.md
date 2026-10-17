---
id: test-rel21.0-capability-profiles
title: Test Suite for Release 21.0 Capability Profiles
release: "21.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel21.0-uc001-capability-profiles

<!-- field: preconditions -->
## Preconditions

- srd057-capability-profiles defines the capability profile, the consumption forms it supports, and the binding discipline the consistency check reports against.
- srd022 R6 self_invoke, srd027 machine_request bindings, and srd055 blueprints are available and unchanged by this release.
- A capability fixture exists whose machine reaches a terminal state and which declares no serve lifecycle, control or monitor server, or static assets.

<!-- field: test_cases -->
## Test cases

<!-- case: A capability fixture loads, terminates, and carries no hosting | status: planned -->
### A capability fixture loads, terminates, and carries no hosting

Load the capability fixture and run it to a terminal state; then add a control server, a monitor server, and a static_assets binding to copies of it.

<!-- field: use_case -->
#### Use case

rel21.0-uc001-capability-profiles

<!-- field: go_test -->
#### Go test

go test ./internal/load ./internal/tools/catalog -run 'TestCapabilityProfileLoadsAndTerminates|TestCapabilityProfileDeclaresNoServingAssets'

<!-- field: expected | form: "tree" -->
#### Expected

- loads_and_terminates: true
- serving_assets_reported: true
- no_extra_profile_field_required: true

<!-- field: traces -->
#### Traces

- rel21.0-uc001-capability-profiles S1
- srd057-capability-profiles AC1
- srd057-capability-profiles R1.1

<!-- case: A capability is reached only by a declared path | status: planned -->
### A capability is reached only by a declared path

Reference the fixture by a machine_request profile path, by a child-agent config.profile, and by a serving wrapper's fields; then by a bare name with no path.

<!-- field: use_case -->
#### Use case

rel21.0-uc001-capability-profiles

<!-- field: go_test -->
#### Go test

go test ./internal/load -run 'TestCapabilityProfileIsReachedOnlyByPath'

<!-- field: expected | form: "tree" -->
#### Expected

- path_references_resolve: true
- bare_name_rejected: true

<!-- field: traces -->
#### Traces

- rel21.0-uc001-capability-profiles S1
- srd057-capability-profiles AC1
- srd057-capability-profiles R1.3

<!-- case: One unchanged capability serves the route, the child, and the wrapper | status: planned -->
### One unchanged capability serves the route, the child, and the wrapper

Bind the fixture by machine_request naming profile with its terminal states mapped to responses, invoke the same file through self_invoke with request and output selectors, and host it from a serving wrapper; compare the fixture's bytes across the three.

<!-- field: use_case -->
#### Use case

rel21.0-uc001-capability-profiles

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest ./internal/tools/control ./internal/load -run 'TestMachineRequestBindsACapabilityProfile|TestSelfInvokeRunsACapabilityProfile|TestServingWrapperHostsACapabilityProfile'

<!-- field: expected | form: "tree" -->
#### Expected

- terminal_states_map_to_responses: true
- selectors_undo_and_trace_propagate: true
- hosted_without_edit: true
- capability_bytes_identical: true

<!-- field: traces -->
#### Traces

- rel21.0-uc001-capability-profiles S2
- srd057-capability-profiles AC2
- srd057-capability-profiles R2.1
- srd057-capability-profiles R2.2
- srd057-capability-profiles R2.3

<!-- case: The in-process boundary stays evaluation-scoped | status: planned -->
### The in-process boundary stays evaluation-scoped

Point run_point at the capability fixture, which declares no point_machine, point_tools, or point_tool_declarations.

<!-- field: use_case -->
#### Use case

rel21.0-uc001-capability-profiles

<!-- field: go_test -->
#### Go test

go test ./internal/tools/evaluation -run 'TestRunPointRequiresItsPointMachine'

<!-- field: expected | form: "tree" -->
#### Expected

- registration_fails_naming_missing_fields: true

<!-- field: traces -->
#### Traces

- rel21.0-uc001-capability-profiles S2
- srd057-capability-profiles AC2
- srd057-capability-profiles R2.4

<!-- case: Shared machines are reported and private ones are not | status: planned -->
### Shared machines are reported and private ones are not

Point two profiles at one machine file; add a machine_request override whose machine another profile also reaches; then leave a machine reached by exactly one profile.

<!-- field: use_case -->
#### Use case

rel21.0-uc001-capability-profiles

<!-- field: go_test -->
#### Go test

go test ./pkg/profileaudit -run 'TestAuditReportsAMachineReachedByTwoProfiles|TestAuditReportsAnOverrideWhoseMachineIsShared|TestAuditIsSilentForAPrivateMachine'

<!-- field: expected | form: "tree" -->
#### Expected

- shared_machine_reported_with_every_reaching_profile: true
- shared_override_reported: true
- reported_as_finding_not_load_error: true
- private_machine_silent: true

<!-- field: traces -->
#### Traces

- rel21.0-uc001-capability-profiles S3
- srd057-capability-profiles AC3
- srd057-capability-profiles R3.1
- srd057-capability-profiles R3.2
- srd057-capability-profiles R3.3
