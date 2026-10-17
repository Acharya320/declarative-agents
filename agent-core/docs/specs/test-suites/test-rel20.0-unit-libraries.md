---
id: test-rel20.0-unit-libraries
title: Test Suite for Release 20.0 Unit Libraries
release: "20.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel20.0-uc003-unit-libraries

<!-- field: preconditions -->
## Preconditions

- srd056-unit-libraries defines import path forms, declared library roots, read-only libraries, closure, dump, and staging behavior, and the service-library convention.
- srd050 imports as amended by srd056, srd051 type units, srd052 expansion, srd054 machine template instances, srd055 blueprint machine mappings, and pkg/profilestage are available.

<!-- field: test_cases -->
## Test cases

<!-- case: An import is relative or under a library root and nothing else | status: implemented -->
### An import is relative or under a library root and nothing else

Import a unit by a relative path, by an /opt/agent-core path, and by /home/user/unit.yaml, in a tool file and a REST file.

<!-- field: use_case -->
#### Use case

rel20.0-uc003-unit-libraries

<!-- field: go_test -->
#### Go test

go test ./internal/load ./internal/tools/catalog ./internal/support/corepath -run 'TestLoadClosureRejectsAnAbsoluteImportOutsideALibraryRoot|TestToolFragmentOutsideALibraryRootIsRejected|TestImportTargetResolvesRelativeAndLibraryRootedPaths'

<!-- field: expected | form: "tree" -->
#### Expected

- relative_accepted: true
- rooted_accepted: true
- other_absolute_rejected_with_srd050_wording: true

<!-- field: traces -->
#### Traces

- rel20.0-uc003-unit-libraries S1
- srd056-unit-libraries AC1
- srd056-unit-libraries R1.1

<!-- case: Rooted paths work in every form srd050 governs | status: implemented -->
### Rooted paths work in every form srd050 governs

Name a rooted path in an expand fragment, a machine template instance, a blueprint machine mapping, and a profile tool_declarations list; then put an unrooted absolute path in the profile list and in an import.

<!-- field: use_case -->
#### Use case

rel20.0-uc003-unit-libraries

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/tools/rest/definition ./internal/runtime/core ./internal/load -run 'TestToolImportsAndFragmentsResolveUnderTheAgentCoreLibrary|TestRESTImportsAndFragmentsResolveUnderTheAgentCoreLibrary|TestLoadMachineClosureSplicesAStageFromTheAgentCoreLibrary|TestProfileDeclaresLibraryRootsAndMapsItsOwnPaths|TestLoadClosureRejectsAnAbsoluteImportOutsideALibraryRoot'

<!-- field: expected | form: "tree" -->
#### Expected

- each_form_resolves: true
- profile_absolute_path_accepted: true
- import_absolute_path_rejected: true

<!-- field: traces -->
#### Traces

- rel20.0-uc003-unit-libraries S1
- srd056-unit-libraries AC1
- srd056-unit-libraries R1.2

<!-- case: The agent-core root maps through --core-root or is used as written | status: implemented -->
### The agent-core root maps through --core-root or is used as written

Import /opt/agent-core/tools/units/types-core.yaml from an application unit with an install root set, then with none set against a tree laid out as the image.

<!-- field: use_case -->
#### Use case

rel20.0-uc003-unit-libraries

<!-- field: go_test -->
#### Go test

go test ./internal/support/corepath ./internal/load ./internal/tools/catalog ./internal/tools/rest/definition ./internal/runtime/core -run 'TestImportTargetResolvesRelativeAndLibraryRootedPaths|TestLoadClosureImportsAUnitFromTheAgentCoreLibrary|TestToolImportsAndFragmentsResolveUnderTheAgentCoreLibrary|TestRESTImportsAndFragmentsResolveUnderTheAgentCoreLibrary|TestLoadMachineClosureSplicesAStageFromTheAgentCoreLibrary'

<!-- field: expected | form: "tree" -->
#### Expected

- mapped_with_install_root: true
- used_as_written_without: true

<!-- field: traces -->
#### Traces

- rel20.0-uc003-unit-libraries S1
- srd056-unit-libraries AC1
- srd056-unit-libraries R1.3

<!-- case: A profile declares a named root for a directory | status: implemented -->
### A profile declares a named root for a directory

Declare libraries: {mesh-lib: ../lib} and import /opt/mesh-lib/units/types.yaml; then declare agent-core and a name with an underscore.

<!-- field: use_case -->
#### Use case

rel20.0-uc003-unit-libraries

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/load -run 'TestProfileDeclaresLibraryRootsAndMapsItsOwnPaths|TestLoadClosureResolvesADeclaredLibraryRoot|TestLoadClosureRejectsAReservedOrMalformedRootName'

<!-- field: expected | form: "tree" -->
#### Expected

- import_maps_to_declared_directory: true
- agent_core_name_rejected: true
- malformed_name_rejected: true

<!-- field: traces -->
#### Traces

- rel20.0-uc003-unit-libraries S2
- srd056-unit-libraries AC2
- srd056-unit-libraries R2.1

<!-- case: --library overrides a declared root for one run | status: implemented -->
### --library overrides a declared root for one run

Run with --library mesh-lib=<other dir>, then with --library absent-lib=<dir>.

<!-- field: use_case -->
#### Use case

rel20.0-uc003-unit-libraries

<!-- field: go_test -->
#### Go test

go test ./cmd/agent ./internal/load -run 'TestLibraryFlagOverridesADeclaredRootOnly|TestLoadClosureAppliesALibraryOverride|TestRootCommandExposesEveryRuntimeFlag'

<!-- field: expected | form: "tree" -->
#### Expected

- override_applies: true
- undeclared_name_startup_error: true

<!-- field: traces -->
#### Traces

- rel20.0-uc003-unit-libraries S2
- srd056-unit-libraries AC2
- srd056-unit-libraries R2.2

<!-- case: Roots are closure-scoped and must agree | status: implemented -->
### Roots are closure-scoped and must agree

Import under /opt/unknown; declare one root name with two directories on a profile and its blueprint; load a child profile relying on a parent root.

<!-- field: use_case -->
#### Use case

rel20.0-uc003-unit-libraries

<!-- field: go_test -->
#### Go test

go test ./internal/load ./internal/support/corepath -run 'TestLoadClosureReportsAnUndeclaredLibraryPrefix|TestLoadClosureResolvesADeclaredLibraryRoot|TestDeclaredLibraryRootsResolveScopedAndReadOnly'

<!-- field: expected | form: "tree" -->
#### Expected

- undeclared_prefix_names_prefix_and_profile: true
- conflict_names_both: deferred with agent blueprints (GH-2123); a profile has one libraries mapping
- child_does_not_inherit: true

<!-- field: traces -->
#### Traces

- rel20.0-uc003-unit-libraries S2
- srd056-unit-libraries AC2
- srd056-unit-libraries R2.3

<!-- case: A library is read-only to its importers | status: implemented -->
### A library is read-only to its importers

Put override: true on a tool in a rooted file; give a rooted file a ../ import that leaves the root; give one an import of another root by prefix.

<!-- field: use_case -->
#### Use case

rel20.0-uc003-unit-libraries

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/support/corepath -run 'TestLibraryFileCannotOverrideAnImport|TestLibraryRelativeImportCannotLeaveItsRoot|TestDeclaredLibraryRootsResolveScopedAndReadOnly'

<!-- field: expected | form: "tree" -->
#### Expected

- override_rejected_names_root_and_file: true
- escape_rejected_names_root_and_file: true
- cross_root_prefix_import_accepted: true

<!-- field: traces -->
#### Traces

- rel20.0-uc003-unit-libraries S2
- srd056-unit-libraries AC2
- srd056-unit-libraries R2.4

<!-- case: Rooted files are closure files and digest assets | status: implemented -->
### Rooted files are closure files and digest assets

Load a profile importing through both roots, read Files, change one byte of a rooted file, and rebuild the program identity.

<!-- field: use_case -->
#### Use case

rel20.0-uc003-unit-libraries

<!-- field: go_test -->
#### Go test

go test ./internal/load -run 'TestLoadClosureImportsAUnitFromTheAgentCoreLibrary|TestLoadClosureReportsAnUnusedLibraryImport|TestLoadClosureResolvesADeclaredLibraryRoot'

<!-- field: expected | form: "tree" -->
#### Expected

- files_include_rooted_canonical_paths: true
- digest_changes: true

<!-- field: traces -->
#### Traces

- rel20.0-uc003-unit-libraries S3
- srd056-unit-libraries AC3
- srd056-unit-libraries R3.1

<!-- case: The dump marks rooted files with their root | status: implemented -->
### The dump marks rooted files with their root

Dump a profile importing through both roots.

<!-- field: use_case -->
#### Use case

rel20.0-uc003-unit-libraries

<!-- field: go_test -->
#### Go test

go test ./internal/load ./cmd/agent -run 'TestLoadClosureImportsAUnitFromTheAgentCoreLibrary|TestLoadClosureResolvesADeclaredLibraryRoot|TestDumpGolden'

<!-- field: expected | form: "tree" -->
#### Expected

- root_marker: agent-core, declared name
- two_dumps_identical: true

<!-- field: traces -->
#### Traces

- rel20.0-uc003-unit-libraries S3
- srd056-unit-libraries AC3
- srd056-unit-libraries R3.2

<!-- case: Staging copies declared-root files and leaves agent-core to the image | status: implemented -->
### Staging copies declared-root files and leaves agent-core to the image

Stage a profile importing under a declared root and under /opt/agent-core, remove the source checkout, and load the staged profile with only --core-root.

<!-- field: use_case -->
#### Use case

rel20.0-uc003-unit-libraries

<!-- field: go_test -->
#### Go test

go test ./pkg/profilestage -run 'TestStageCopiesADeclaredLibraryRootToItsStagedPosition|TestStageLeavesAgentCoreLibraryImportsToTheImage|TestStageReportsAnUndeclaredLibraryRoot'

<!-- field: expected | form: "tree" -->
#### Expected

- declared_root_files_copied: true
- agent_core_files_not_copied: true
- staged_loads_without_library_flag: true

<!-- field: traces -->
#### Traces

- rel20.0-uc003-unit-libraries S3
- srd056-unit-libraries AC3
- srd056-unit-libraries R3.3

<!-- case: The service-library convention is documented | status: implemented -->
### The service-library convention is documented

Read the service library section of cfg-rest-tool-format.

<!-- field: use_case -->
#### Use case

rel20.0-uc003-unit-libraries

<!-- field: go_test -->
#### Go test

mage audit

<!-- field: expected | form: "tree" -->
#### Expected

- layout_documented: REST unit, type unit, client words
- consumers_import_through_root: true

<!-- field: traces -->
#### Traces

- rel20.0-uc003-unit-libraries S4
- srd056-unit-libraries AC4
- srd056-unit-libraries R4.1

<!-- case: No service library has fewer than two importing profiles | status: implemented -->
### No service library has fewer than two importing profiles

Count importing profiles for every service-library unit, an instantiation counting as an importer (srd056 R4.2). The retrieval client has three: chatbot-mesh and agentic-wiki-mesh expand the fragment, and cohere-demo imports the type unit.

<!-- field: use_case -->
#### Use case

rel20.0-uc003-unit-libraries

<!-- field: go_test -->
#### Go test

mage stats:reuse

<!-- field: expected | form: "tree" -->
#### Expected

- min_importing_profiles: 2

<!-- field: traces -->
#### Traces

- rel20.0-uc003-unit-libraries S4
- srd056-unit-libraries AC4
- srd056-unit-libraries R4.2
