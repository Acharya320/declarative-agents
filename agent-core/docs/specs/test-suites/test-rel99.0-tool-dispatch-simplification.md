---
id: test-rel99.0-tool-dispatch-simplification
title: Test Suite for Release 99.0 Tool Dispatch Simplification
release: "99.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- [rel99.0-uc001-tool-dispatch-simplification](../use-cases/rel99.0-uc001-tool-dispatch-simplification.md)

<!-- field: preconditions -->
## Preconditions

- srd037-tool-dispatch-simplification defines the design scope before code changes.
- The sibling code task owns executable registry and availability regression tests.
- MachineSpec, Command, Builder, Result, and Dispatch semantics remain unchanged.

<!-- field: test_cases -->
## Test cases

<!-- case: Design preserves the current core dispatch model | status: done -->
### Design preserves the current core dispatch model

Review srd037 and confirm it keeps MachineSpec sequencing,
core.Registry lookup, core.Builder, core.Command, core.Result, and
Dispatch semantics intact while naming only startup and availability
extension points for the follow-up code task.

<!-- field: use_case -->
#### Use case

[rel99.0-uc001-tool-dispatch-simplification](../use-cases/rel99.0-uc001-tool-dispatch-simplification.md)

<!-- field: go_test -->
#### Go test

mage audit

<!-- field: inputs | form: "tree" -->
#### Inputs

- design_doc: srd037-tool-dispatch-simplification

<!-- field: expected | form: "tree" -->
#### Expected

- dispatch_model_unchanged: true
- extension_points_named: true

<!-- field: traces -->
#### Traces

- [srd037-tool-dispatch-simplification](../software-requirements/srd037-tool-dispatch-simplification.md) [AC1](../software-requirements/srd037-tool-dispatch-simplification.md#ac1)

<!-- case: Registry startup errors are scoped to configuration-originated failures | status: done -->
### Registry startup errors are scoped to configuration-originated failures

Focused registry tests prove empty names, duplicate names, and
post-freeze registration can be returned as errors during startup
without removing compatibility wrappers.

<!-- field: use_case -->
#### Use case

[rel99.0-uc001-tool-dispatch-simplification](../use-cases/rel99.0-uc001-tool-dispatch-simplification.md)

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run 'TestRegistry(CheckedRegistrationReturnsConfigurationErrors|CompatibilityWrappersStillPanic)'

<!-- field: inputs | form: "tree" -->
#### Inputs

- registry: core.Registry
- startup_path: profile tool declaration registration

<!-- field: expected | form: "tree" -->
#### Expected

- configuration_errors_returned: true
- compatibility_wrappers_available: true

<!-- field: traces -->
#### Traces

- [srd037-tool-dispatch-simplification](../software-requirements/srd037-tool-dispatch-simplification.md) [AC2](../software-requirements/srd037-tool-dispatch-simplification.md#ac2)

<!-- case: Builtin registration paths share one helper | status: done -->
### Builtin registration paths share one helper

Registry package tests prove RegisterSingleBuiltin and RegisterUnifiedTools
use the same builtin declaration helper while exec declarations still use
ExecBuilderFactory and binary validation.

<!-- field: use_case -->
#### Use case

[rel99.0-uc001-tool-dispatch-simplification](../use-cases/rel99.0-uc001-tool-dispatch-simplification.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/registry -run 'Test(RegisterUnifiedToolsDuplicateReturnsRegistryError|RegisterSingleBuiltinUnknownInitMatchesUnifiedPath|RegisterSingleBuiltinOverrides)'

<!-- field: inputs | form: "tree" -->
#### Inputs

- builtin_paths:
  1. RegisterSingleBuiltin
  1. RegisterUnifiedTools
- exec_path: registerExec

<!-- field: expected | form: "tree" -->
#### Expected

- builtin_helper_shared: true
- exec_registration_separate: true

<!-- field: traces -->
#### Traces

- [srd037-tool-dispatch-simplification](../software-requirements/srd037-tool-dispatch-simplification.md) [AC3](../software-requirements/srd037-tool-dispatch-simplification.md#ac3)

<!-- case: Parse-time and dispatch-time availability use the same rule | status: done -->
### Parse-time and dispatch-time availability use the same rule

The regression test proves parse_response and dynamic $tool dispatch agree
with the manifest on state-unavailable external tools through
core.Registry.ResolveExternalTool.

<!-- field: use_case -->
#### Use case

[rel99.0-uc001-tool-dispatch-simplification](../use-cases/rel99.0-uc001-tool-dispatch-simplification.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/registry -run TestDynamicToolActionAndParseResponseShareAvailabilityRule

<!-- field: inputs | form: "tree" -->
#### Inputs

- paths:
  1. Registry.Manifest
  1. parse_response
  1. dynamic_$tool_dispatch
- states:
  1. available
  1. unavailable

<!-- field: expected | form: "tree" -->
#### Expected

- shared_availability_rule: true
- existing_error_contracts_preserved: true

<!-- field: traces -->
#### Traces

- [srd037-tool-dispatch-simplification](../software-requirements/srd037-tool-dispatch-simplification.md) [AC4](../software-requirements/srd037-tool-dispatch-simplification.md#ac4)

<!-- case: Dynamic phase contradictions fail at startup | status: done -->
### Dynamic phase contradictions fail at startup

Phase validation traces the actual invoke selector and parse word that
feed a $tool transition. It rejects their mismatches and every cause that
leaves an external word scoped to no state, while unrelated invoke words
and declaration order cannot change parse validation state.

<!-- field: use_case -->
#### Use case

[rel99.0-uc001-tool-dispatch-simplification](../use-cases/rel99.0-uc001-tool-dispatch-simplification.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/tools/llm -run 'Test(ValidateToolPhases(RejectsActualSelectorStateDifferentFromDynamicTarget|RejectsActualParserStateDifferentFromDynamicTarget|AcceptsManifestStateMatchingDynamicTarget|RejectsWordScopedToNoState|IgnoresUnrelatedInvokeManifestStatesAndOrder)|ParseResponseOwnsManifestStateIndependentOfInvokeRegistration|PhaseScopedChapterExamplesLoadAndDeriveAvailability)'

<!-- field: inputs | form: "tree" -->
#### Inputs

- manifest_states:
  1. matching_target
  1. mismatched_selector
  1. mismatched_parser
  1. unrelated_invoke
- zero_phase_causes:
  1. empty_explicit_intersection
  1. missing_emits
  1. unroutable_emit

<!-- field: expected | form: "tree" -->
#### Expected

- mismatched_validating_state_rejected: true
- unrelated_invoke_ignored: true
- declaration_order_independent: true
- zero_phase_word_rejected_with_cause: true
- matching_target_state_accepted: true
- chapter_examples_pass_complete_phase_validation: true

<!-- field: traces -->
#### Traces

- [srd037-tool-dispatch-simplification](../software-requirements/srd037-tool-dispatch-simplification.md) [R4.6](../software-requirements/srd037-tool-dispatch-simplification.md#r46)
- [srd037-tool-dispatch-simplification](../software-requirements/srd037-tool-dispatch-simplification.md) [R4.7](../software-requirements/srd037-tool-dispatch-simplification.md#r47)
- [srd037-tool-dispatch-simplification](../software-requirements/srd037-tool-dispatch-simplification.md) [AC4](../software-requirements/srd037-tool-dispatch-simplification.md#ac4)

<!-- case: Code task boundaries exclude unrelated dispatch rewrites | status: done -->
### Code task boundaries exclude unrelated dispatch rewrites

Review the code task and resulting diff to ensure it does not change
MachineSpec parsing, BuildTransitionTable semantics, command execution,
receipt behavior, or LLM response envelope syntax.

<!-- field: use_case -->
#### Use case

[rel99.0-uc001-tool-dispatch-simplification](../use-cases/rel99.0-uc001-tool-dispatch-simplification.md)

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core ./internal/tools/registry ./internal/tools/llm

<!-- field: inputs | form: "tree" -->
#### Inputs

- code_task: GH-171#173

<!-- field: expected | form: "tree" -->
#### Expected

- machinespec_unchanged: true
- command_execution_unchanged: true
- parser_syntax_unchanged: true

<!-- field: traces -->
#### Traces

- [srd037-tool-dispatch-simplification](../software-requirements/srd037-tool-dispatch-simplification.md) [AC5](../software-requirements/srd037-tool-dispatch-simplification.md#ac5)
