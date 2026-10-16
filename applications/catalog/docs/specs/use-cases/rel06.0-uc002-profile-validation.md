---
id: rel06.0-uc002-profile-validation
title: Profile Validation and Release Bundle Readiness
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

Profile authors validate that every profile-shaped YAML file can resolve its
declared machine, tool selection, declarations, REST definitions, and config
directories after migration. Validation runs against an external agent-core
checkout because agent-core still owns shared tool declarations and runtime
execution.

Release bundle readiness also needs a local tag target. The profile bundle
tag records a reviewed set of YAML programs, profile-local UI assets,
documentation, and fixtures without claiming ownership of application
presentations or agent-core runtime image tags.

Validation is intentionally boring. It should fail on missing files, stale
core-agent references, or unsupported profile names. It should not execute
profile behavior or mutate a workspace.

<!-- field: actor -->
## Actor

Profile repository maintainer.

<!-- field: trigger -->
## Trigger

A maintainer prepares profile changes for review, release bundle tagging, or
mounted execution:

  mage validate

After review, the maintainer can create a local profile bundle tag:

  mage tag

<!-- field: preconditions -->
## Preconditions

- agent-core is available as a sibling ../../agent-core checkout or through demo.yaml core_root.
- Profile assets live under agents/.
- The repository is on main before mage tag creates a release tag.

<!-- field: flow -->
## Flow

- F1: mage validate discovers profile.yaml, profile-*.yaml, and *-profile.yaml files under agents/.
- F2: Each profile is parsed and must declare a machine path and at least one tools.yaml selection.
- F3: Profile-local references resolve from the profile directory.
- F4: /opt/agent-core/tools references resolve through the selected agent-core checkout.
- F5: References to /opt/agent-core/agents are rejected because migrated profiles must not depend on copied core agent assets.
- F6: mage containerSmoke verifies the installed image lacks /opt/agent-core/agents before running a mounted specification-critic profile.
- F7: mage tag creates the next v0.YYYYMMDD.N tag from local tags when the repository is on main.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd001-agent-functional-blocks R1 | Validation protects the profile ownership boundary. |
| T2 | sm-multi-agent-profile-system profile_asset_rules | Profile-local and core-owned references follow separate ownership rules. |
| T3 | sm-multi-agent-profile-system validation_views | Validation views name profile reference checks and mounted container execution. |
| T4 | sm-agent-role-taxonomy taxonomy_rules | Checks name validation, mounted execution, integration proof, or documentation review paths. |
| T5 | mage validate | Local validation checks every profile-shaped YAML file. |
| T6 | mage tag | Release bundle tags version profile-owned assets. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

mage validate checks all supported profile filename shapes.

<!-- field: traces -->
#### Traces

- srd001-agent-functional-blocks R1.1
- sm-agent-role-taxonomy taxonomy_rules

<!-- sc: S2 -->
### S2

Validation resolves profile-local paths from the selected profile directory.

<!-- field: traces -->
#### Traces

- srd001-agent-functional-blocks R1.1
- sm-multi-agent-profile-system profile_asset_rules

<!-- sc: S3 -->
### S3

Validation maps /opt/agent-core/tools to the selected agent-core checkout without accepting /opt/agent-core/agents.

<!-- field: traces -->
#### Traces

- srd001-agent-functional-blocks R1.2
- srd001-agent-functional-blocks R1.3

<!-- sc: S4 -->
### S4

Container smoke checks image exclusion and mounted-profile execution.

<!-- field: traces -->
#### Traces

- sm-multi-agent-profile-system validation_views

<!-- sc: S5 -->
### S5

mage tag creates profile bundle tags without using agent-core image tags.

<!-- field: traces -->
#### Traces

- srd001-agent-functional-blocks R1.3

<!-- field: out_of_scope -->
## Out of scope

- Publishing tags to a remote.
- Building or tagging agent-core images.
- Adding new profile families.

<!-- field: test_suite_trace -->
## Test suite trace

test-rel06.0-agent-profiles
