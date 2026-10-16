---
id: test-rel99.0-tool-dispatch-simplification
title: Test Suite for Release 99.0 Tool Dispatch Simplification
release: "99.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel99.0-uc001-tool-dispatch-simplification

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

rel99.0-uc001-tool-dispatch-simplification

<!-- field: go_test -->
#### Go test

mage audit

<!-- field: inputs -->
#### Inputs

```yaml
design_doc: srd037-tool-dispatch-simplification
```

<!-- field: expected -->
#### Expected

```yaml
dispatch_model_unchanged: true
extension_points_named: true
```

<!-- field: traces -->
#### Traces

- srd037-tool-dispatch-simplification AC1

<!-- case: Registry startup errors are scoped to configuration-originated failures | status: done -->
### Registry startup errors are scoped to configuration-originated failures

Focused registry tests prove empty names, duplicate names, and
post-freeze registration can be returned as errors during startup
without removing compatibility wrappers.

<!-- field: use_case -->
#### Use case

rel99.0-uc001-tool-dispatch-simplification

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run 'TestRegistry(CheckedRegistrationReturnsConfigurationErrors|CompatibilityWrappersStillPanic)'

<!-- field: inputs -->
#### Inputs

```yaml
registry: core.Registry
startup_path: profile tool declaration registration
```

<!-- field: expected -->
#### Expected

```yaml
configuration_errors_returned: true
compatibility_wrappers_available: true
```

<!-- field: traces -->
#### Traces

- srd037-tool-dispatch-simplification AC2

<!-- case: Builtin registration paths share one helper | status: done -->
### Builtin registration paths share one helper

Registry package tests prove RegisterSingleBuiltin and RegisterUnifiedTools
use the same builtin declaration helper while exec declarations still use
ExecBuilderFactory and binary validation.

<!-- field: use_case -->
#### Use case

rel99.0-uc001-tool-dispatch-simplification

<!-- field: go_test -->
#### Go test

go test ./internal/tools/registry -run 'Test(RegisterUnifiedToolsDuplicateReturnsRegistryError|RegisterSingleBuiltinUnknownInitMatchesUnifiedPath|RegisterSingleBuiltinOverrides)'

<!-- field: inputs -->
#### Inputs

```yaml
builtin_paths: [RegisterSingleBuiltin, RegisterUnifiedTools]
exec_path: registerExec
```

<!-- field: expected -->
#### Expected

```yaml
builtin_helper_shared: true
exec_registration_separate: true
```

<!-- field: traces -->
#### Traces

- srd037-tool-dispatch-simplification AC3

<!-- case: Parse-time and dispatch-time availability use the same rule | status: done -->
### Parse-time and dispatch-time availability use the same rule

The regression test proves parse_response and dynamic $tool dispatch agree
with the manifest on state-unavailable external tools through
core.Registry.ResolveExternalTool.

<!-- field: use_case -->
#### Use case

rel99.0-uc001-tool-dispatch-simplification

<!-- field: go_test -->
#### Go test

go test ./internal/tools/registry -run TestDynamicToolActionAndParseResponseShareAvailabilityRule

<!-- field: inputs -->
#### Inputs

```yaml
paths: [Registry.Manifest, parse_response, dynamic_$tool_dispatch]
states: [available, unavailable]
```

<!-- field: expected -->
#### Expected

```yaml
shared_availability_rule: true
existing_error_contracts_preserved: true
```

<!-- field: traces -->
#### Traces

- srd037-tool-dispatch-simplification AC4

<!-- case: Dynamic phase contradictions fail at startup | status: done -->
### Dynamic phase contradictions fail at startup

Phase validation traces the actual invoke selector and parse word that
feed a $tool transition. It rejects their mismatches and every cause that
leaves an external word scoped to no state, while unrelated invoke words
and declaration order cannot change parse validation state.

<!-- field: use_case -->
#### Use case

rel99.0-uc001-tool-dispatch-simplification

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/tools/llm -run 'Test(ValidateToolPhases(RejectsActualSelectorStateDifferentFromDynamicTarget|RejectsActualParserStateDifferentFromDynamicTarget|AcceptsManifestStateMatchingDynamicTarget|RejectsWordScopedToNoState|IgnoresUnrelatedInvokeManifestStatesAndOrder)|ParseResponseOwnsManifestStateIndependentOfInvokeRegistration|PhaseScopedChapterExamplesLoadAndDeriveAvailability)'

<!-- field: inputs -->
#### Inputs

```yaml
manifest_states: [matching_target, mismatched_selector, mismatched_parser, unrelated_invoke]
zero_phase_causes: [empty_explicit_intersection, missing_emits, unroutable_emit]
```

<!-- field: expected -->
#### Expected

```yaml
mismatched_validating_state_rejected: true
unrelated_invoke_ignored: true
declaration_order_independent: true
zero_phase_word_rejected_with_cause: true
matching_target_state_accepted: true
chapter_examples_pass_complete_phase_validation: true
```

<!-- field: traces -->
#### Traces

- srd037-tool-dispatch-simplification R4.6
- srd037-tool-dispatch-simplification R4.7
- srd037-tool-dispatch-simplification AC4

<!-- case: Code task boundaries exclude unrelated dispatch rewrites | status: done -->
### Code task boundaries exclude unrelated dispatch rewrites

Review the code task and resulting diff to ensure it does not change
MachineSpec parsing, BuildTransitionTable semantics, command execution,
receipt behavior, or LLM response envelope syntax.

<!-- field: use_case -->
#### Use case

rel99.0-uc001-tool-dispatch-simplification

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core ./internal/tools/registry ./internal/tools/llm

<!-- field: inputs -->
#### Inputs

```yaml
code_task: GH-171#173
```

<!-- field: expected -->
#### Expected

```yaml
machinespec_unchanged: true
command_execution_unchanged: true
parser_syntax_unchanged: true
```

<!-- field: traces -->
#### Traces

- srd037-tool-dispatch-simplification AC5
