---
id: test-rel07.1-profile-boundaries
title: Test Suite for Release 07.1 Profile Boundaries
release: "07.1"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

This suite validates the per-family boundary use cases for the specification-critic,
lifecycle, and rest profiles. Release 07.0 gave dedicated boundary treatment to
evaluator, generator, planner, bench, monitor/control, and the
documentation-curator profiles. Release 07.1 closes the coverage gap for the
three families that were previously validated only by the generic rel06.0
mounted-profile and profile-validation use cases.

Each test case names the profile assets, machine word selection, tool
selection, inputs, and expected terminal states for one family, and traces to
the success criteria of its use case and to the family SRD.

The suite records evidence honestly. A test case is specified when the use
case and profile assets exist but no profile-owned Mage integration target
runs the flow yet. Profile-owned integration targets that exercise these
fixtures, mirroring magefiles/integration_*.go for the covered families, are a
deferred follow-up tracked in the road-map rel07.1 entry.

<!-- field: traces -->
## Traces

- rel07.1-uc001-specification-critic-profile-boundary
- rel07.1-uc002-lifecycle-profile-boundary
- rel07.1-uc003-rest-profile-boundary

<!-- field: preconditions -->
## Preconditions

- agent-core is available as a sibling ../../agent-core checkout or through demo.yaml core_root.
- applications/catalog owns all profile paths, machines, tool selections, suites, REST definitions, and OpenAPI assets named by each test case.
- Profiles load through the demo.yaml catalog_root or repository-relative paths.
- agent-core remains responsible for the specification validation tools, lifecycle ports and checkpoint store, and the reusable HTTP runtime.

<!-- field: test_cases -->
## Test cases

<!-- case: Specification-critic deterministic charter validation boundary | status: specified -->
### Specification-critic deterministic charter validation boundary

Run the specification-critic profile against a target directory and the demo charter
suite. Confirm the machine exposes load, validation, rg search, finding
reduction, and report steps with no LLM tool family and produces a
provenance-rich report for both pass and fail outcomes.

<!-- field: use_case -->
#### Use case

rel07.1-uc001-specification-critic-profile-boundary

<!-- field: inputs | form: "tree" -->
#### Inputs

- profile: agents/specification-critic/profile.yaml
- machine: agents/specification-critic/machine.yaml
- tools: agents/specification-critic/tools.yaml
- charter_suite: agents/specification-critic/suites/demo-charter.yaml
- target_directory: a specification corpus directory

<!-- field: expected | form: "tree" -->
#### Expected

- selected_tools:
  1. load_corpus
  1. validate_specs
  1. jurist_rg
  1. reduce_grep_checks
  1. format_report
- executable_dependency: rg
- visible_grep_steps:
  1. jurist_rg
  1. reduce_grep_checks
- llm_tool_family_selected: false
- report_formatted_on_pass: true
- report_formatted_on_fail: true
- finding_provenance_fields:
  1. suite
  1. check
  1. severity
  1. message
  1. file
  1. line
- terminal_states:
  1. Passed
  1. Failed

<!-- field: traces -->
#### Traces

- rel07.1-uc001-specification-critic-profile-boundary S1
- rel07.1-uc001-specification-critic-profile-boundary S2
- rel07.1-uc001-specification-critic-profile-boundary S3
- srd005-specification-critic AC1
- srd005-specification-critic AC3
- sm-agent-role-taxonomy specification-critic

<!-- case: Specification-critic declarative test-evidence audit boundary | status: implemented -->
### Specification-critic declarative test-evidence audit boundary

Confirm the explicit audit profile selects declared Go inventory and run
exec words, visible claim and result reducers, and the shared specification-critic report.
The runtime receives an ordinary profile invocation rather than a special
test-evidence CLI mode.

<!-- field: use_case -->
#### Use case

rel07.1-uc001-specification-critic-profile-boundary

<!-- field: go_test -->
#### Go test

TestSpecificationCriticAuditProfileDeclaresEvidencePipeline

<!-- field: inputs | form: "tree" -->
#### Inputs

- profile: agents/specification-critic/audit-profile.yaml
- machine: agents/specification-critic/audit-machine.yaml
- tools: agents/specification-critic/audit-tools.yaml
- declarations: agents/specification-critic/go-test.yaml

<!-- field: expected | form: "tree" -->
#### Expected

- claims_step: load_test_claims
- inventory_steps:
  1. go_module
  1. go_packages_raw
  1. go_packages
  1. go_test_inventory
- package_normalization_input: $from(go_packages_raw).output
- resolution_step: resolve_test_evidence
- inventory_failure_reduction: resolve_test_evidence
- inventory_failure_reported: true
- implicit_module_download: false
- execution_step: go_test_run
- conformance_tmpdir: private_staged_directory
- reduction_step: reduce_test_evidence_run
- report_step: format_report
- special_cli_mode: false

<!-- field: traces -->
#### Traces

- rel07.1-uc001-specification-critic-profile-boundary S4
- srd005-specification-critic R2.5
- srd005-specification-critic AC2

<!-- case: Lifecycle history, rollback, and approval boundary | status: specified -->
### Lifecycle history, rollback, and approval boundary

Run the history, rollback, and approval lifecycle programs. Confirm history
and rollback run as single declared machine actions, and confirm the
approval machine suspends and resumes through explicit Approved or Rejected
signals rather than hidden CLI branches.

<!-- field: use_case -->
#### Use case

rel07.1-uc002-lifecycle-profile-boundary

<!-- field: inputs | form: "tree" -->
#### Inputs

- history_machine: testdata/conformance/lifecycle/history.yaml
- rollback_machine: testdata/conformance/lifecycle/rollback.yaml
- approval_machine: testdata/conformance/lifecycle/approval/machine.yaml
- history_rollback_tools: testdata/conformance/lifecycle/tools.yaml
- approval_tools: testdata/conformance/lifecycle/approval/tools.yaml
- checkpoint_id: a saved checkpoint ID
- target_iteration: a rollback target iteration
- approval_signal: Approved or Rejected

<!-- field: expected | form: "tree" -->
#### Expected

- history_selected_tool: checkpoint_history
- rollback_selected_tool: checkpoint_rollback
- approval_selected_tools:
  1. suspend
  1. done
- suspend_is_machine_action: true
- approved_reaches: Succeeded
- rejected_reaches: Rejected
- approved_path_skipped_on_reject: true
- terminal_states:
  1. Done
  1. Succeeded
  1. Rejected
  1. Failed

<!-- field: traces -->
#### Traces

- rel07.1-uc002-lifecycle-profile-boundary S1
- rel07.1-uc002-lifecycle-profile-boundary S2
- rel07.1-uc002-lifecycle-profile-boundary S3
- srd009-lifecycle AC1
- srd009-lifecycle AC3
- sm-agent-role-taxonomy lifecycle

<!-- case: REST client and server boundary words | status: specified -->
### REST client and server boundary words

Run the sample REST flow and the ollama-* variant. Confirm client read,
async create and await, and server launch, await, and stop run as explicit
machine states, with REST authority in the profile-owned REST definitions
and OpenAPI imports that do not generate Go code.

<!-- field: use_case -->
#### Use case

rel07.1-uc003-rest-profile-boundary

<!-- field: inputs | form: "tree" -->
#### Inputs

- profile: testdata/conformance/rest/profile.yaml
- ollama_profile: testdata/conformance/rest/ollama-profile.yaml
- machine: testdata/conformance/rest/machine.yaml
- tools: testdata/conformance/rest/tools.yaml
- rest_definitions: testdata/conformance/rest/rest.yaml
- ollama_rest_definitions: testdata/conformance/rest/ollama-rest.yaml
- openapi_assets:
  1. testdata/conformance/rest/openapi/payments.yaml
  1. testdata/conformance/rest/openapi/ollama.yaml

<!-- field: expected | form: "tree" -->
#### Expected

- selected_tools:
  1. get_payment
  1. create_payment
  1. await_payment
  1. launch_payment_webhooks
  1. await_payment_webhook
  1. stop_payment_webhooks
- server_states_explicit:
  1. LaunchingWebhook
  1. WaitingWebhook
  1. StoppingWebhook
- rest_authority_owner: profile
- openapi_generates_go_code: false
- listener_released_at_terminal: true
- terminal_states:
  1. Succeeded
  1. Failed

<!-- field: traces -->
#### Traces

- rel07.1-uc003-rest-profile-boundary S1
- rel07.1-uc003-rest-profile-boundary S2
- rel07.1-uc003-rest-profile-boundary S3
- srd007-rest AC1
- srd007-rest AC3
- sm-agent-role-taxonomy rest
