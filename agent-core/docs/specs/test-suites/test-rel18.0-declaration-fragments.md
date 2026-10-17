---
id: test-rel18.0-declaration-fragments
title: Test Suite for Release 18.0 Declaration Fragments
release: "18.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel18.0-uc001-declaration-fragments

<!-- field: preconditions -->
## Preconditions

- srd052-declaration-fragments defines the fragment format, expansion, hygienic substitution, usedness, provenance, and stage splicing.
- srd050 imports and closure usedness, srd051 type units and the schema subset, and the canonical dump are available.
- Machine stage cases additionally require the srd038 label checks and srd006 exhaustiveness checks.

<!-- field: test_cases -->
## Test cases

<!-- case: A fragment with params and one body loads | status: implemented -->
### A fragment with params and one body loads

Load a fragment carrying unit, params of every allowed type, and a tools body; then one with a rest body.

<!-- field: use_case -->
#### Use case

rel18.0-uc001-declaration-fragments

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/tools/rest/definition -run 'TestExpandProducesPrefixedTypedTools|TestExpandRESTFragmentTwiceWithPrefixes|TestFragmentBodyIsDecodedOnlyAfterSubstitution'

<!-- field: expected | form: "tree" -->
#### Expected

- loads: true
- param_types: string, integer, number, boolean

<!-- field: traces -->
#### Traces

- rel18.0-uc001-declaration-fragments S1
- srd052-declaration-fragments AC1
- srd052-declaration-fragments R1.1

<!-- case: A fragment may import units but may not expand a fragment | status: implemented -->
### A fragment may import units but may not expand a fragment

Load a fragment importing a type unit, then one whose expand list names another fragment.

<!-- field: use_case -->
#### Use case

rel18.0-uc001-declaration-fragments

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog -run 'TestFragmentImportsResolveFromTheFragment|TestFragmentCannotNest'

<!-- field: expected | form: "tree" -->
#### Expected

- type_import_accepted: true
- nested_expand_rejected: true

<!-- field: traces -->
#### Traces

- rel18.0-uc001-declaration-fragments S1
- srd052-declaration-fragments AC1
- srd052-declaration-fragments R1.2

<!-- case: Half a fragment is a load error naming the file | status: implemented -->
### Half a fragment is a load error naming the file

Load a file with params and no body, then a file with a body and no params but a $param reference.

<!-- field: use_case -->
#### Use case

rel18.0-uc001-declaration-fragments

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog -run 'TestHalfDeclaredFragmentFails'

<!-- field: expected | form: "tree" -->
#### Expected

- error_names_file: true

<!-- field: traces -->
#### Traces

- rel18.0-uc001-declaration-fragments S1
- srd052-declaration-fragments AC1
- srd052-declaration-fragments R1.3

<!-- case: An expand entry names fragment, args, and an optional prefix | status: implemented -->
### An expand entry names fragment, args, and an optional prefix

Expand with and without as; resolve fragment relative to the importing file, including from a re-rooted directory.

<!-- field: use_case -->
#### Use case

rel18.0-uc001-declaration-fragments

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog -run 'TestExpandProducesPrefixedTypedTools|TestPlainUnitCannotBeExpanded'

<!-- field: expected | form: "tree" -->
#### Expected

- fragment_resolved_relative_to_importer: true
- as_optional: true

<!-- field: traces -->
#### Traces

- rel18.0-uc001-declaration-fragments S2
- srd052-declaration-fragments AC2
- srd052-declaration-fragments R2.1

<!-- case: Argument faults name importer, fragment, and parameter | status: implemented -->
### Argument faults name importer, fragment, and parameter

Omit a required argument, supply an unknown one, pass a string for an integer, and pass a value outside an enum; a parameter with a default may be omitted.

<!-- field: use_case -->
#### Use case

rel18.0-uc001-declaration-fragments

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/tools/rest/definition -run 'TestExpandArgumentFaultsNameImporterFragmentAndParameter|TestRESTArgumentFaultsNameImporterFragmentAndParameter'

<!-- field: expected | form: "tree" -->
#### Expected

- every_error_names_importer_fragment_parameter: true
- default_fills_omitted: true

<!-- field: traces -->
#### Traces

- rel18.0-uc001-declaration-fragments S2
- srd052-declaration-fragments AC2
- srd052-declaration-fragments R2.2

<!-- case: Substitution is hygienic and leaves no reference behind | status: implemented -->
### Substitution is hygienic and leaves no reference behind

Reference a parameter inside a scalar, inside a mapping name, and inside a list; reference an undeclared parameter.

<!-- field: use_case -->
#### Use case

rel18.0-uc001-declaration-fragments

<!-- field: go_test -->
#### Go test

go test ./internal/fragments ./internal/tools/catalog -run 'TestSubstitute|TestReferenceInMappingNameIsLeftAndReported'

<!-- field: expected | form: "tree" -->
#### Expected

- scalar_substituted: true
- mapping_name_not_substituted_and_rejected: true
- undeclared_reference_names_fragment_and_line: true
- env_and_selector_forms_untouched: true

<!-- field: traces -->
#### Traces

- rel18.0-uc001-declaration-fragments S2
- srd052-declaration-fragments AC2
- srd052-declaration-fragments R2.3

<!-- case: The expanded document takes the hand-written validation path | status: implemented -->
### The expanded document takes the hand-written validation path

Expand a fragment producing a tool with a misspelled field, then one whose per-tool hook rejects its config.

<!-- field: use_case -->
#### Use case

rel18.0-uc001-declaration-fragments

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog -run 'TestExpandedUnitTakesTheHandWrittenPath'

<!-- field: expected | form: "tree" -->
#### Expected

- strict_decode_applies: true
- per_tool_hooks_apply: true
- no_generated_file_written: true

<!-- field: traces -->
#### Traces

- rel18.0-uc001-declaration-fragments S2
- srd052-declaration-fragments AC2
- srd052-declaration-fragments R2.4

<!-- case: A prefix makes repeated expansion possible | status: implemented -->
### A prefix makes repeated expansion possible

Expand one fragment twice with as: cohere and as: openai, then read both produced names.

<!-- field: use_case -->
#### Use case

rel18.0-uc001-declaration-fragments

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/tools/rest/definition -run 'TestExpandProducesPrefixedTypedTools|TestExpandRESTFragmentTwiceWithPrefixes'

<!-- field: expected | form: "tree" -->
#### Expected

- names_prefixed: <as>_<name>
- both_selectable: true

<!-- field: traces -->
#### Traces

- rel18.0-uc001-declaration-fragments S3
- srd052-declaration-fragments AC2
- srd052-declaration-fragments R2.5

<!-- case: A fragment under imports is rejected | status: implemented -->
### A fragment under imports is rejected

List a fragment file under imports rather than expand.

<!-- field: use_case -->
#### Use case

rel18.0-uc001-declaration-fragments

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/tools/rest/definition -run 'TestFragmentUnderImportsIsRejected|TestRESTFragmentUnderImportsIsRejected'

<!-- field: expected | form: "tree" -->
#### Expected

- error_names_importer_and_fragment: true

<!-- field: traces -->
#### Traces

- rel18.0-uc001-declaration-fragments S3
- srd052-declaration-fragments AC2
- srd052-declaration-fragments R2.6

<!-- case: Repeated expansion without a prefix is the existing duplicate error | status: implemented -->
### Repeated expansion without a prefix is the existing duplicate error

Expand one fragment twice with identical args and no as.

<!-- field: use_case -->
#### Use case

rel18.0-uc001-declaration-fragments

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/tools/rest/definition -run 'TestRepeatedExpansionWithoutPrefixIsADuplicate|TestRepeatedRESTExpansionWithoutPrefixIsTheConflictError'

<!-- field: expected | form: "tree" -->
#### Expected

- error: srd050 duplicate-name error, not a fragment-specific one

<!-- field: traces -->
#### Traces

- rel18.0-uc001-declaration-fragments S3
- srd052-declaration-fragments AC2
- srd052-declaration-fragments R2.7

<!-- case: An unselected expansion is an unused import | status: implemented -->
### An unselected expansion is an unused import

Expand a fragment whose produced tools tools.yaml never selects; then select one of them.

<!-- field: use_case -->
#### Use case

rel18.0-uc001-declaration-fragments

<!-- field: go_test -->
#### Go test

go test ./internal/load -run 'TestLoadClosureReportsAnUnselectedExpansionAsUnused|TestLoadClosureReportsAnUnreferencedRESTExpansionAsUnused'

<!-- field: expected | form: "tree" -->
#### Expected

- unselected_fails_as_unused_import: true
- evidence: importing file, fragment path, arguments
- selected_passes: true

<!-- field: traces -->
#### Traces

- rel18.0-uc001-declaration-fragments S4
- srd052-declaration-fragments AC3
- srd052-declaration-fragments R3.1

<!-- case: The dump lists every expansion | status: implemented -->
### The dump lists every expansion

Dump a profile with two expansions and compare the expansions section across two runs.

<!-- field: use_case -->
#### Use case

rel18.0-uc001-declaration-fragments

<!-- field: go_test -->
#### Go test

go test ./internal/load ./cmd/agent -run 'TestLoadClosureDumpsEveryExpansion|TestLoadClosureDumpsRESTExpansionsBesideToolOnes|TestDumpGolden'

<!-- field: expected | form: "tree" -->
#### Expected

- section: expansions
- entries_carry: fragment, args, produced names
- two_dumps_identical: true

<!-- field: traces -->
#### Traces

- rel18.0-uc001-declaration-fragments S4
- srd052-declaration-fragments AC3
- srd052-declaration-fragments R3.2

<!-- case: A stage fragment splices under a prefix with parameterized entry and exit | status: implemented -->
### A stage fragment splices under a prefix with parameterized entry and exit

Expand a stage fragment into a machine with entry and exit signals passed as arguments and read the resulting states and transitions.

<!-- field: use_case -->
#### Use case

rel18.0-uc001-declaration-fragments

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run 'TestLoadMachineClosureSplicesAStageTwice|TestStageArgumentFaultsNameFragmentAndParameter'

<!-- field: expected | form: "tree" -->
#### Expected

- states_prefixed: true
- entry_exit_from_args: true

<!-- field: traces -->
#### Traces

- rel18.0-uc001-declaration-fragments S5
- srd052-declaration-fragments AC4
- srd052-declaration-fragments R4.1

<!-- case: A spliced machine is checked as one machine | status: implemented -->
### A spliced machine is checked as one machine

Splice a stage whose transition reads a label the importing machine never publishes, then one with an unhandled signal, then a correct one.

<!-- field: use_case -->
#### Use case

rel18.0-uc001-declaration-fragments

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core ./internal/load ./cmd/agent -run 'TestSplicedStageIsValidatedAsOneMachine|TestSplicingTheSamePrefixTwiceIsTheDuplicateStateError|TestLoadClosureDumpsMachineStageExpansions|TestDumpGolden'

<!-- field: expected | form: "tree" -->
#### Expected

- label_check_applies: true
- exhaustiveness_applies: true
- correct_splice_loads: true

<!-- field: traces -->
#### Traces

- rel18.0-uc001-declaration-fragments S5
- srd052-declaration-fragments AC4
- srd052-declaration-fragments R4.2

<!-- case: A stage fragment path selects its variant by environment | status: implemented -->
### A stage fragment path selects its variant by environment

Expand stages/s-${X:-a}.yaml beside s-a.yaml and s-b.yaml, load with X unset and with X=b, then select a variant whose splice leaves a transition to an undeclared state, then use a path without a default and one whose default file is missing. Stage the machine's tree and package it with appmanifest (TestResolveStagesEveryVariantOfATemplatedMachineStage in magefiles).

<!-- field: use_case -->
#### Use case

rel18.0-uc001-declaration-fragments

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core ./internal/support/envexpand ./pkg/profilestage -run 'TestStageFragmentVariantIsSelectedByEnvironment|TestSelectedStageVariantIsValidatedAsOneMachine|TestTemplatedStagePathsThatCannotSelectAVariantFail|TestSelectVariant|TestStageCarriesEveryVariantOfAnEnvironmentSelectedStage|TestStageRefusesAVariantPathItCannotStage'

<!-- field: expected | form: "tree" -->
#### Expected

- unset_selects_default: true
- bound_selects_variant: true
- one_variant_spliced: true
- dump_names_resolved_path: true
- broken_variant_fails_naming_resolved_fragment: true
- missing_default_or_no_default_fails_naming_path: true
- staging_and_packaging_carry_every_variant: true

<!-- field: traces -->
#### Traces

- rel18.0-uc001-declaration-fragments S5
- srd052-declaration-fragments AC5
- srd052-declaration-fragments R4.3
