---
id: test-rel20.0-machine-templates
title: Test Suite for Release 20.0 Machine Templates and Agent Blueprints
release: "20.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel20.0-uc001-machine-templates
- rel20.0-uc002-agent-blueprints

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

rel20.0-uc001-machine-templates

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run 'TestMachineTemplateInstanceEqualsTheHandWrittenMachine|TestMachineTemplateBodyHoldsNoParamsUnitOrImports'

<!-- field: expected -->
#### Expected

```yaml
loads: true
nested_params_rejected: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc001-machine-templates S1
- srd054-machine-templates AC1
- srd054-machine-templates R1.1

<!-- case: An instance carries only unit, one expand, name, and purpose | status: implemented -->
### An instance carries only unit, one expand, name, and purpose

Add states to an instance, add an as to its entry, add a second expand entry, and override name and purpose.

<!-- field: use_case -->
#### Use case

rel20.0-uc001-machine-templates

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run 'TestMachineInstanceCarriesOnlyItsHeader'

<!-- field: expected -->
#### Expected

```yaml
extra_field_error_names_field_and_instance: true
as_rejected: true
second_entry_rejected: true
name_and_purpose_replace_template: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc001-machine-templates S1
- srd054-machine-templates AC1
- srd054-machine-templates R1.2

<!-- case: A template may splice stages and nesting stops at one level | status: implemented -->
### A template may splice stages and nesting stops at one level

Expand a template whose body expands a stage fragment; then a template expanding a template; then an instance expanding a stage beside its template.

<!-- field: use_case -->
#### Use case

rel20.0-uc001-machine-templates

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run 'TestMachineTemplateSplicesItsStagesRelativeToItself|TestMachineTemplateCannotNestOrBeLoadedAsAMachine|TestMachineInstanceCarriesOnlyItsHeader'

<!-- field: expected -->
#### Expected

```yaml
stage_inside_template_spliced: true
template_in_template_rejected: true
stage_beside_template_rejected: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc001-machine-templates S1
- srd054-machine-templates AC1
- srd054-machine-templates R1.3

<!-- case: A template named where a machine is expected is rejected | status: implemented -->
### A template named where a machine is expected is rejected

Point a profile machine field at a template file, then a machine_request machine binding.

<!-- field: use_case -->
#### Use case

rel20.0-uc001-machine-templates

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core ./internal/load -run 'TestMachineTemplateCannotNestOrBeLoadedAsAMachine|TestLoadClosureRefusesAProfileNamingAMachineTemplate'

<!-- field: expected -->
#### Expected

```yaml
error_names_referrer_and_template: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc001-machine-templates S1
- srd054-machine-templates AC1
- srd054-machine-templates R1.4

<!-- case: Template arguments follow srd052 and pass through to stages | status: implemented -->
### Template arguments follow srd052 and pass through to stages

Omit a required argument, supply an unknown one, leave a reference unfilled; then pass an argument through to a stage expansion in the body and to an action name.

<!-- field: use_case -->
#### Use case

rel20.0-uc001-machine-templates

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run 'TestMachineTemplateFaultsNameTemplateAndInstance|TestMachineTemplateSplicesItsStagesRelativeToItself'

<!-- field: expected -->
#### Expected

```yaml
argument_faults_match_srd052: true
reference_in_action_name_substituted: true
reference_in_stage_args_substituted: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc001-machine-templates S2
- srd054-machine-templates AC2
- srd054-machine-templates R2.1

<!-- case: Template paths resolve against the template, instance paths against the instance | status: implemented -->
### Template paths resolve against the template, instance paths against the instance

Place the template and its stage in one directory and the instance in another; load the instance.

<!-- field: use_case -->
#### Use case

rel20.0-uc001-machine-templates

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core ./pkg/profilestage -run 'TestMachineTemplateSplicesItsStagesRelativeToItself|TestStageFollowsMachineTemplateEdges'

<!-- field: expected -->
#### Expected

```yaml
instance_fragment_relative_to_instance: true
stage_path_relative_to_template: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc001-machine-templates S2
- srd054-machine-templates AC2
- srd054-machine-templates R2.2

<!-- case: The target body kind selects replace or splice | status: implemented -->
### The target body kind selects replace or splice

Expand a machine body, a stage body, and a tools body from a machine file.

<!-- field: use_case -->
#### Use case

rel20.0-uc001-machine-templates

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run 'TestExpandTargetBodyKindSelectsReplaceOrSplice|TestMachineTemplateSplicesItsStagesRelativeToItself'

<!-- field: expected -->
#### Expected

```yaml
machine_body_replaces: true
stage_body_splices: true
other_body_error_names_target_and_kind: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc001-machine-templates S2
- srd054-machine-templates AC2
- srd054-machine-templates R2.3

<!-- case: An expanded machine takes the hand-written validation path | status: implemented -->
### An expanded machine takes the hand-written validation path

Compare the parsed control-template machine with the parsed control machine; expand a template with an unreachable state and one missing a required policy.

<!-- field: use_case -->
#### Use case

rel20.0-uc001-machine-templates

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core -run 'TestMachineTemplateInstanceEqualsTheHandWrittenMachine|TestMachineTemplateFaultsNameTemplateAndInstance'

<!-- field: expected -->
#### Expected

```yaml
parsed_specs_equal_except_name: true
structural_checks_apply: true
required_policy_applies: true
no_generated_file_written: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc001-machine-templates S2
- srd054-machine-templates AC2
- srd054-machine-templates R2.4

<!-- case: Diagnostics in an expanded machine name template and instance | status: implemented -->
### Diagnostics in an expanded machine name template and instance

Misspell a label a transition in the template reads, then run --validate-config on the instance profile.

<!-- field: use_case -->
#### Use case

rel20.0-uc001-machine-templates

<!-- field: go_test -->
#### Go test

go test ./internal/runtime/core ./cmd/agent -run 'TestMachineTemplateFaultsNameTemplateAndInstance|TestStartupWiringFaultNamesTemplateAndInstance'

<!-- field: expected -->
#### Expected

```yaml
error_names_template_path: true
error_names_instance_path: true
leftover_reference_line_points_into_template: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc001-machine-templates S2
- srd054-machine-templates AC2
- srd054-machine-templates R2.5

<!-- case: Templates and their stages are closure files and digest assets | status: implemented -->
### Templates and their stages are closure files and digest assets

Load a profile whose machine is an instance, read Files, then change one byte of the template and rebuild the program identity.

<!-- field: use_case -->
#### Use case

rel20.0-uc001-machine-templates

<!-- field: go_test -->
#### Go test

go test ./internal/load -run 'TestLoadClosureCarriesTheMachineTemplate'

<!-- field: expected -->
#### Expected

```yaml
files_include_instance_template_stages: true
digest_changes: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc001-machine-templates S3
- srd054-machine-templates AC3
- srd054-machine-templates R3.1

<!-- case: The dump renders the expanded machine and labels expansion kinds | status: implemented -->
### The dump renders the expanded machine and labels expansion kinds

Dump control-template twice and read its machine and expansions sections.

<!-- field: use_case -->
#### Use case

rel20.0-uc001-machine-templates

<!-- field: go_test -->
#### Go test

go test ./internal/load ./cmd/agent -run 'TestLoadClosureCarriesTheMachineTemplate|TestDumpGolden'

<!-- field: expected -->
#### Expected

```yaml
machine_section: expanded machine
expansion_kinds: "machine, stage"
two_dumps_identical: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc001-machine-templates S3
- srd054-machine-templates AC3
- srd054-machine-templates R3.2

<!-- case: Staging follows instance, template, and stage edges | status: implemented -->
### Staging follows instance, template, and stage edges

Stage a profile whose machine is an instance of a template in another directory, remove the source, and load the staged profile.

<!-- field: use_case -->
#### Use case

rel20.0-uc001-machine-templates

<!-- field: go_test -->
#### Go test

go test ./pkg/profilestage -run 'TestStageFollowsMachineTemplateEdges'

<!-- field: expected -->
#### Expected

```yaml
staged_tree_loads_without_checkout: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc001-machine-templates S3
- srd054-machine-templates AC3
- srd054-machine-templates R3.3

<!-- case: A blueprint with params and a profile body loads through an instance | status: planned -->
### A blueprint with params and a profile body loads through an instance

Load control-blueprint, whose profile.yaml expands a blueprint holding every AgentProfile field.

<!-- field: use_case -->
#### Use case

rel20.0-uc002-agent-blueprints

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog -run 'TestBlueprintInstanceLoads'

<!-- field: expected -->
#### Expected

```yaml
loads: true
every_profile_field_accepted: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc002-agent-blueprints S1
- srd055-agent-blueprints AC1
- srd055-agent-blueprints R1.1

<!-- case: Instance fields replace blueprint fields wholesale | status: planned -->
### Instance fields replace blueprint fields wholesale

Override tools with a shorter list and directory with another path; add as to the entry.

<!-- field: use_case -->
#### Use case

rel20.0-uc002-agent-blueprints

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog -run 'TestBlueprintInstanceOverridesReplaceWholesale|TestBlueprintEntryRejectsAs'

<!-- field: expected -->
#### Expected

```yaml
list_replaced_not_merged: true
scalar_replaced: true
as_rejected: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc002-agent-blueprints S1
- srd055-agent-blueprints AC1
- srd055-agent-blueprints R1.2

<!-- case: A blueprint cannot nest and is refused where a profile is expected | status: planned -->
### A blueprint cannot nest and is refused where a profile is expected

Expand a blueprint from a blueprint; pass a blueprint to --profile, to a child-agent profile reference, and to a machine_request profile binding.

<!-- field: use_case -->
#### Use case

rel20.0-uc002-agent-blueprints

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/tools/rest/definition ./cmd/agent -run 'TestBlueprintCannotNest|TestBlueprintNamedAsProfileFails|TestMachineRequestNamingBlueprintFails'

<!-- field: expected -->
#### Expected

```yaml
nesting_rejected: true
error_names_referrer_and_blueprint: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc002-agent-blueprints S1
- srd055-agent-blueprints AC1
- srd055-agent-blueprints R1.3

<!-- case: Blueprint arguments follow srd052 | status: planned -->
### Blueprint arguments follow srd052

Omit a required argument, supply an unknown one, pass a wrong type, and leave a reference unfilled.

<!-- field: use_case -->
#### Use case

rel20.0-uc002-agent-blueprints

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog -run 'TestBlueprintArgumentFaults'

<!-- field: expected -->
#### Expected

```yaml
argument_faults_match_srd052: true
body_decoded_after_substitution: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc002-agent-blueprints S2
- srd055-agent-blueprints AC2
- srd055-agent-blueprints R2.1

<!-- case: Blueprint paths resolve against the blueprint, instance paths against the instance | status: planned -->
### Blueprint paths resolve against the blueprint, instance paths against the instance

Place the blueprint and its units in one directory and the instance in another; override one path on the instance; use an install-prefix path in each.

<!-- field: use_case -->
#### Use case

rel20.0-uc002-agent-blueprints

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog -run 'TestBlueprintPathsResolveAgainstTheirFiles'

<!-- field: expected -->
#### Expected

```yaml
blueprint_relative_to_blueprint: true
override_relative_to_instance: true
install_prefix_mapped_in_both: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc002-agent-blueprints S2
- srd055-agent-blueprints AC2
- srd055-agent-blueprints R2.2

<!-- case: A machine mapping expands a template directly | status: planned -->
### A machine mapping expands a template directly

Give the blueprint machine: {fragment: <template>, args: {server: $param(server)}} and load two instances with different server arguments.

<!-- field: use_case -->
#### Use case

rel20.0-uc002-agent-blueprints

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/load -run 'TestBlueprintMachineMappingExpandsTemplate'

<!-- field: expected -->
#### Expected

```yaml
template_expanded_without_instance_file: true
arguments_pass_through: true
instances_differ_by_argument: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc002-agent-blueprints S2
- srd055-agent-blueprints AC2
- srd055-agent-blueprints R2.3

<!-- case: A resolved profile satisfies hand-written profile rules | status: planned -->
### A resolved profile satisfies hand-written profile rules

Resolve a blueprint whose body omits tools and whose instance does not supply them.

<!-- field: use_case -->
#### Use case

rel20.0-uc002-agent-blueprints

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog -run 'TestResolvedBlueprintProfileRequiresMachineAndTools'

<!-- field: expected -->
#### Expected

```yaml
error_names_instance_and_blueprint: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc002-agent-blueprints S2
- srd055-agent-blueprints AC2
- srd055-agent-blueprints R2.4

<!-- case: The instance is the identity and the blueprint is in its digest | status: planned -->
### The instance is the identity and the blueprint is in its digest

Build the program identity for an instance, then change one byte of the blueprint and of the machine template it reaches.

<!-- field: use_case -->
#### Use case

rel20.0-uc002-agent-blueprints

<!-- field: go_test -->
#### Go test

go test ./internal/load -run 'TestBlueprintInstanceProgramIdentity'

<!-- field: expected -->
#### Expected

```yaml
identity_profile: instance path
digest_changes_with_blueprint: true
digest_changes_with_template: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc002-agent-blueprints S3
- srd055-agent-blueprints AC3
- srd055-agent-blueprints R3.1

<!-- case: The dump renders the resolved profile and lists the blueprint | status: planned -->
### The dump renders the resolved profile and lists the blueprint

Dump control-blueprint twice and read its profile and expansions sections.

<!-- field: use_case -->
#### Use case

rel20.0-uc002-agent-blueprints

<!-- field: go_test -->
#### Go test

go test ./internal/load ./cmd/agent -run 'TestLoadClosureDumpsBlueprintExpansion|TestDumpGolden'

<!-- field: expected -->
#### Expected

```yaml
profile_section: resolved instance
expansion_kind: profile
two_dumps_identical: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc002-agent-blueprints S3
- srd055-agent-blueprints AC3
- srd055-agent-blueprints R3.2

<!-- case: Staging follows blueprint edges | status: planned -->
### Staging follows blueprint edges

Stage an instance whose blueprint and machine template live in other directories, remove the source, and load the staged profile.

<!-- field: use_case -->
#### Use case

rel20.0-uc002-agent-blueprints

<!-- field: go_test -->
#### Go test

go test ./pkg/profilestage -run 'TestStageFollowsBlueprintEdges'

<!-- field: expected -->
#### Expected

```yaml
staged_tree_loads_without_checkout: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc002-agent-blueprints S3
- srd055-agent-blueprints AC3
- srd055-agent-blueprints R3.3

<!-- case: profileaudit resolves a blueprint-based child profile | status: planned -->
### profileaudit resolves a blueprint-based child profile

Inspect a parent profile whose child-agent reference names a blueprint instance.

<!-- field: use_case -->
#### Use case

rel20.0-uc002-agent-blueprints

<!-- field: go_test -->
#### Go test

go test ./pkg/profileaudit -run 'TestInspectResolvesBlueprintChildProfile'

<!-- field: expected -->
#### Expected

```yaml
child_resolved: true
child_wiring_checked: true
```

<!-- field: traces -->
#### Traces

- rel20.0-uc002-agent-blueprints S3
- srd055-agent-blueprints AC3
- srd055-agent-blueprints R3.4
