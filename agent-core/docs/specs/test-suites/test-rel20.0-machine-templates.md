---
id: test-rel20.0-machine-templates
title: Test Suite for Release 20.0 Machine Templates and Agent Blueprints
release: "20.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- [rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md)
- [rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md)

<!-- field: preconditions -->
## Preconditions

- srd054-machine-templates defines the machine template and instance formats, expansion, validation, and closure, identity, dump, and staging behavior.
- srd055-agent-blueprints defines the blueprint and instance formats, path resolution, the machine mapping, and identity, dump, staging, and audit behavior.
- srd052 fragments, stage splicing, and the expansions dump section; srd006 and srd038 machine validation; and pkg/profilestage and pkg/profileaudit are available.

<!-- field: test_cases -->
## Test cases

<!-- case: A template with params and a machine body loads through an instance | status: implemented -->
### A template with params and a machine body loads through an instance

Load control-template, whose machine.yaml expands a template holding the control machine as a machine body; then a template carrying a nested params.

<!-- field: use_case -->
#### Use case

[rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md)

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run 'TestMachineTemplateInstanceEqualsTheHandWrittenMachine|TestMachineTemplateBodyHoldsNoParamsUnitOrImports'

<!-- field: expected | form: "tree" -->
#### Expected

- loads: true
- nested_params_rejected: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md) [S1](../use-cases/rel20.0-uc001-machine-templates.md#s1)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [AC1](../software-requirements/srd054-machine-templates.md#ac1)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [R1.1](../software-requirements/srd054-machine-templates.md#r11)

<!-- case: An instance carries only unit, one expand, name, and purpose | status: implemented -->
### An instance carries only unit, one expand, name, and purpose

Add states to an instance, add an as to its entry, add a second expand entry, and override name and purpose.

<!-- field: use_case -->
#### Use case

[rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md)

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run 'TestMachineInstanceCarriesOnlyItsHeader'

<!-- field: expected | form: "tree" -->
#### Expected

- extra_field_error_names_field_and_instance: true
- as_rejected: true
- second_entry_rejected: true
- name_and_purpose_replace_template: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md) [S1](../use-cases/rel20.0-uc001-machine-templates.md#s1)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [AC1](../software-requirements/srd054-machine-templates.md#ac1)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [R1.2](../software-requirements/srd054-machine-templates.md#r12)

<!-- case: A template may splice stages and nesting stops at one level | status: implemented -->
### A template may splice stages and nesting stops at one level

Expand a template whose body expands a stage fragment; then a template expanding a template; then an instance expanding a stage beside its template.

<!-- field: use_case -->
#### Use case

[rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md)

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run 'TestMachineTemplateSplicesItsStagesRelativeToItself|TestMachineTemplateCannotNestOrBeLoadedAsAMachine|TestMachineInstanceCarriesOnlyItsHeader'

<!-- field: expected | form: "tree" -->
#### Expected

- stage_inside_template_spliced: true
- template_in_template_rejected: true
- stage_beside_template_rejected: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md) [S1](../use-cases/rel20.0-uc001-machine-templates.md#s1)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [AC1](../software-requirements/srd054-machine-templates.md#ac1)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [R1.3](../software-requirements/srd054-machine-templates.md#r13)

<!-- case: A template named where a machine is expected is rejected | status: implemented -->
### A template named where a machine is expected is rejected

Point a profile machine field at a template file, then a machine_request machine binding.

<!-- field: use_case -->
#### Use case

[rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md)

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core ./internal/load -run 'TestMachineTemplateCannotNestOrBeLoadedAsAMachine|TestLoadClosureRefusesAProfileNamingAMachineTemplate'

<!-- field: expected | form: "tree" -->
#### Expected

- error_names_referrer_and_template: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md) [S1](../use-cases/rel20.0-uc001-machine-templates.md#s1)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [AC1](../software-requirements/srd054-machine-templates.md#ac1)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [R1.4](../software-requirements/srd054-machine-templates.md#r14)

<!-- case: Template arguments follow srd052 and pass through to stages | status: implemented -->
### Template arguments follow srd052 and pass through to stages

Omit a required argument, supply an unknown one, leave a reference unfilled; then pass an argument through to a stage expansion in the body and to an action name.

<!-- field: use_case -->
#### Use case

[rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md)

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run 'TestMachineTemplateFaultsNameTemplateAndInstance|TestMachineTemplateSplicesItsStagesRelativeToItself'

<!-- field: expected | form: "tree" -->
#### Expected

- argument_faults_match_srd052: true
- reference_in_action_name_substituted: true
- reference_in_stage_args_substituted: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md) [S2](../use-cases/rel20.0-uc001-machine-templates.md#s2)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [AC2](../software-requirements/srd054-machine-templates.md#ac2)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [R2.1](../software-requirements/srd054-machine-templates.md#r21)

<!-- case: Template paths resolve against the template, instance paths against the instance | status: implemented -->
### Template paths resolve against the template, instance paths against the instance

Place the template and its stage in one directory and the instance in another; load the instance.

<!-- field: use_case -->
#### Use case

[rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md)

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core ./pkg/profilestage -run 'TestMachineTemplateSplicesItsStagesRelativeToItself|TestStageFollowsMachineTemplateEdges'

<!-- field: expected | form: "tree" -->
#### Expected

- instance_fragment_relative_to_instance: true
- stage_path_relative_to_template: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md) [S2](../use-cases/rel20.0-uc001-machine-templates.md#s2)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [AC2](../software-requirements/srd054-machine-templates.md#ac2)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [R2.2](../software-requirements/srd054-machine-templates.md#r22)

<!-- case: The target body kind selects replace or splice | status: implemented -->
### The target body kind selects replace or splice

Expand a machine body, a stage body, and a tools body from a machine file.

<!-- field: use_case -->
#### Use case

[rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md)

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run 'TestExpandTargetBodyKindSelectsReplaceOrSplice|TestMachineTemplateSplicesItsStagesRelativeToItself'

<!-- field: expected | form: "tree" -->
#### Expected

- machine_body_replaces: true
- stage_body_splices: true
- other_body_error_names_target_and_kind: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md) [S2](../use-cases/rel20.0-uc001-machine-templates.md#s2)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [AC2](../software-requirements/srd054-machine-templates.md#ac2)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [R2.3](../software-requirements/srd054-machine-templates.md#r23)

<!-- case: An expanded machine takes the hand-written validation path | status: implemented -->
### An expanded machine takes the hand-written validation path

Compare the parsed control-template machine with the parsed control machine; expand a template with an unreachable state and one missing a required policy.

<!-- field: use_case -->
#### Use case

[rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md)

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run 'TestMachineTemplateInstanceEqualsTheHandWrittenMachine|TestMachineTemplateFaultsNameTemplateAndInstance'

<!-- field: expected | form: "tree" -->
#### Expected

- parsed_specs_equal_except_name: true
- structural_checks_apply: true
- required_policy_applies: true
- no_generated_file_written: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md) [S2](../use-cases/rel20.0-uc001-machine-templates.md#s2)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [AC2](../software-requirements/srd054-machine-templates.md#ac2)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [R2.4](../software-requirements/srd054-machine-templates.md#r24)

<!-- case: Diagnostics in an expanded machine name template and instance | status: implemented -->
### Diagnostics in an expanded machine name template and instance

Misspell a label a transition in the template reads, then run --validate-config on the instance profile.

<!-- field: use_case -->
#### Use case

[rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md)

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core ./cmd/agent -run 'TestMachineTemplateFaultsNameTemplateAndInstance|TestStartupWiringFaultNamesTemplateAndInstance'

<!-- field: expected | form: "tree" -->
#### Expected

- error_names_template_path: true
- error_names_instance_path: true
- leftover_reference_line_points_into_template: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md) [S2](../use-cases/rel20.0-uc001-machine-templates.md#s2)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [AC2](../software-requirements/srd054-machine-templates.md#ac2)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [R2.5](../software-requirements/srd054-machine-templates.md#r25)

<!-- case: Templates and their stages are closure files and digest assets | status: implemented -->
### Templates and their stages are closure files and digest assets

Load a profile whose machine is an instance, read Files, then change one byte of the template and rebuild the program identity.

<!-- field: use_case -->
#### Use case

[rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md)

<!-- field: go_test -->
#### Go test

go test ./internal/load -run 'TestLoadClosureCarriesTheMachineTemplate'

<!-- field: expected | form: "tree" -->
#### Expected

- files_include_instance_template_stages: true
- digest_changes: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md) [S3](../use-cases/rel20.0-uc001-machine-templates.md#s3)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [AC3](../software-requirements/srd054-machine-templates.md#ac3)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [R3.1](../software-requirements/srd054-machine-templates.md#r31)

<!-- case: The dump renders the expanded machine and labels expansion kinds | status: implemented -->
### The dump renders the expanded machine and labels expansion kinds

Dump control-template twice and read its machine and expansions sections.

<!-- field: use_case -->
#### Use case

[rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md)

<!-- field: go_test -->
#### Go test

go test ./internal/load ./cmd/agent -run 'TestLoadClosureCarriesTheMachineTemplate|TestDumpGolden'

<!-- field: expected | form: "tree" -->
#### Expected

- machine_section: expanded machine
- expansion_kinds: machine, stage
- two_dumps_identical: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md) [S3](../use-cases/rel20.0-uc001-machine-templates.md#s3)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [AC3](../software-requirements/srd054-machine-templates.md#ac3)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [R3.2](../software-requirements/srd054-machine-templates.md#r32)

<!-- case: Staging follows instance, template, and stage edges | status: implemented -->
### Staging follows instance, template, and stage edges

Stage a profile whose machine is an instance of a template in another directory, remove the source, and load the staged profile.

<!-- field: use_case -->
#### Use case

[rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md)

<!-- field: go_test -->
#### Go test

go test ./pkg/profilestage -run 'TestStageFollowsMachineTemplateEdges'

<!-- field: expected | form: "tree" -->
#### Expected

- staged_tree_loads_without_checkout: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc001-machine-templates](../use-cases/rel20.0-uc001-machine-templates.md) [S3](../use-cases/rel20.0-uc001-machine-templates.md#s3)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [AC3](../software-requirements/srd054-machine-templates.md#ac3)
- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [R3.3](../software-requirements/srd054-machine-templates.md#r33)

<!-- case: A blueprint with params and a profile body loads through an instance | status: planned -->
### A blueprint with params and a profile body loads through an instance

Load control-blueprint, whose profile.yaml expands a blueprint holding every AgentProfile field.

<!-- field: use_case -->
#### Use case

[rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog -run 'TestBlueprintInstanceLoads'

<!-- field: expected | form: "tree" -->
#### Expected

- loads: true
- every_profile_field_accepted: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md) [S1](../use-cases/rel20.0-uc002-agent-blueprints.md#s1)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [AC1](../software-requirements/srd055-agent-blueprints.md#ac1)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [R1.1](../software-requirements/srd055-agent-blueprints.md#r11)

<!-- case: Instance fields replace blueprint fields wholesale | status: planned -->
### Instance fields replace blueprint fields wholesale

Override tools with a shorter list and directory with another path; add as to the entry.

<!-- field: use_case -->
#### Use case

[rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog -run 'TestBlueprintInstanceOverridesReplaceWholesale|TestBlueprintEntryRejectsAs'

<!-- field: expected | form: "tree" -->
#### Expected

- list_replaced_not_merged: true
- scalar_replaced: true
- as_rejected: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md) [S1](../use-cases/rel20.0-uc002-agent-blueprints.md#s1)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [AC1](../software-requirements/srd055-agent-blueprints.md#ac1)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [R1.2](../software-requirements/srd055-agent-blueprints.md#r12)

<!-- case: A blueprint cannot nest and is refused where a profile is expected | status: planned -->
### A blueprint cannot nest and is refused where a profile is expected

Expand a blueprint from a blueprint; pass a blueprint to --profile, to a child-agent profile reference, and to a machine_request profile binding.

<!-- field: use_case -->
#### Use case

[rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/tools/rest/definition ./cmd/agent -run 'TestBlueprintCannotNest|TestBlueprintNamedAsProfileFails|TestMachineRequestNamingBlueprintFails'

<!-- field: expected | form: "tree" -->
#### Expected

- nesting_rejected: true
- error_names_referrer_and_blueprint: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md) [S1](../use-cases/rel20.0-uc002-agent-blueprints.md#s1)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [AC1](../software-requirements/srd055-agent-blueprints.md#ac1)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [R1.3](../software-requirements/srd055-agent-blueprints.md#r13)

<!-- case: Blueprint arguments follow srd052 | status: planned -->
### Blueprint arguments follow srd052

Omit a required argument, supply an unknown one, pass a wrong type, and leave a reference unfilled.

<!-- field: use_case -->
#### Use case

[rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog -run 'TestBlueprintArgumentFaults'

<!-- field: expected | form: "tree" -->
#### Expected

- argument_faults_match_srd052: true
- body_decoded_after_substitution: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md) [S2](../use-cases/rel20.0-uc002-agent-blueprints.md#s2)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [AC2](../software-requirements/srd055-agent-blueprints.md#ac2)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [R2.1](../software-requirements/srd055-agent-blueprints.md#r21)

<!-- case: Blueprint paths resolve against the blueprint, instance paths against the instance | status: planned -->
### Blueprint paths resolve against the blueprint, instance paths against the instance

Place the blueprint and its units in one directory and the instance in another; override one path on the instance; use an install-prefix path in each.

<!-- field: use_case -->
#### Use case

[rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog -run 'TestBlueprintPathsResolveAgainstTheirFiles'

<!-- field: expected | form: "tree" -->
#### Expected

- blueprint_relative_to_blueprint: true
- override_relative_to_instance: true
- install_prefix_mapped_in_both: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md) [S2](../use-cases/rel20.0-uc002-agent-blueprints.md#s2)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [AC2](../software-requirements/srd055-agent-blueprints.md#ac2)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [R2.2](../software-requirements/srd055-agent-blueprints.md#r22)

<!-- case: A machine mapping expands a template directly | status: planned -->
### A machine mapping expands a template directly

Give the blueprint machine: {fragment: <template>, args: {server: $param(server)}} and load two instances with different server arguments.

<!-- field: use_case -->
#### Use case

[rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/load -run 'TestBlueprintMachineMappingExpandsTemplate'

<!-- field: expected | form: "tree" -->
#### Expected

- template_expanded_without_instance_file: true
- arguments_pass_through: true
- instances_differ_by_argument: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md) [S2](../use-cases/rel20.0-uc002-agent-blueprints.md#s2)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [AC2](../software-requirements/srd055-agent-blueprints.md#ac2)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [R2.3](../software-requirements/srd055-agent-blueprints.md#r23)

<!-- case: A resolved profile satisfies hand-written profile rules | status: planned -->
### A resolved profile satisfies hand-written profile rules

Resolve a blueprint whose body omits tools and whose instance does not supply them.

<!-- field: use_case -->
#### Use case

[rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog -run 'TestResolvedBlueprintProfileRequiresMachineAndTools'

<!-- field: expected | form: "tree" -->
#### Expected

- error_names_instance_and_blueprint: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md) [S2](../use-cases/rel20.0-uc002-agent-blueprints.md#s2)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [AC2](../software-requirements/srd055-agent-blueprints.md#ac2)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [R2.4](../software-requirements/srd055-agent-blueprints.md#r24)

<!-- case: The instance is the identity and the blueprint is in its digest | status: planned -->
### The instance is the identity and the blueprint is in its digest

Build the program identity for an instance, then change one byte of the blueprint and of the machine template it reaches.

<!-- field: use_case -->
#### Use case

[rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md)

<!-- field: go_test -->
#### Go test

go test ./internal/load -run 'TestBlueprintInstanceProgramIdentity'

<!-- field: expected | form: "tree" -->
#### Expected

- identity_profile: instance path
- digest_changes_with_blueprint: true
- digest_changes_with_template: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md) [S3](../use-cases/rel20.0-uc002-agent-blueprints.md#s3)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [AC3](../software-requirements/srd055-agent-blueprints.md#ac3)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [R3.1](../software-requirements/srd055-agent-blueprints.md#r31)

<!-- case: The dump renders the resolved profile and lists the blueprint | status: planned -->
### The dump renders the resolved profile and lists the blueprint

Dump control-blueprint twice and read its profile and expansions sections.

<!-- field: use_case -->
#### Use case

[rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md)

<!-- field: go_test -->
#### Go test

go test ./internal/load ./cmd/agent -run 'TestLoadClosureDumpsBlueprintExpansion|TestDumpGolden'

<!-- field: expected | form: "tree" -->
#### Expected

- profile_section: resolved instance
- expansion_kind: profile
- two_dumps_identical: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md) [S3](../use-cases/rel20.0-uc002-agent-blueprints.md#s3)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [AC3](../software-requirements/srd055-agent-blueprints.md#ac3)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [R3.2](../software-requirements/srd055-agent-blueprints.md#r32)

<!-- case: Staging follows blueprint edges | status: planned -->
### Staging follows blueprint edges

Stage an instance whose blueprint and machine template live in other directories, remove the source, and load the staged profile.

<!-- field: use_case -->
#### Use case

[rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md)

<!-- field: go_test -->
#### Go test

go test ./pkg/profilestage -run 'TestStageFollowsBlueprintEdges'

<!-- field: expected | form: "tree" -->
#### Expected

- staged_tree_loads_without_checkout: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md) [S3](../use-cases/rel20.0-uc002-agent-blueprints.md#s3)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [AC3](../software-requirements/srd055-agent-blueprints.md#ac3)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [R3.3](../software-requirements/srd055-agent-blueprints.md#r33)

<!-- case: profileaudit resolves a blueprint-based child profile | status: planned -->
### profileaudit resolves a blueprint-based child profile

Inspect a parent profile whose child-agent reference names a blueprint instance.

<!-- field: use_case -->
#### Use case

[rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md)

<!-- field: go_test -->
#### Go test

go test ./pkg/profileaudit -run 'TestInspectResolvesBlueprintChildProfile'

<!-- field: expected | form: "tree" -->
#### Expected

- child_resolved: true
- child_wiring_checked: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc002-agent-blueprints](../use-cases/rel20.0-uc002-agent-blueprints.md) [S3](../use-cases/rel20.0-uc002-agent-blueprints.md#s3)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [AC3](../software-requirements/srd055-agent-blueprints.md#ac3)
- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [R3.4](../software-requirements/srd055-agent-blueprints.md#r34)
