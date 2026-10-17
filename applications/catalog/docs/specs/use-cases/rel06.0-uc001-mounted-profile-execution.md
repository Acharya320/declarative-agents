---
id: rel06.0-uc001-mounted-profile-execution
title: Mounted Profile Execution
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

Agent Core runs profile-owned programs from this repository without copying
those programs into the runtime image. Operators provide applications/catalog as a
checkout, release bundle, or /profiles container mount. The runtime image
provides the agent binary and core-owned assets.

This use case covers the execution path after profile migration. The selected
profile, profile-local assets, and catalog integration fixtures stay under
applications/catalog. Agent Core loads the profile through --profile and uses
--directory for the target workspace. Application presentations compose that
profile from their own roots; the Knowledge Manager presentation now lives in
applications/agent-architecture.

The important failure mode is quiet fallback to bundled agent directories.
A correct run cannot require /opt/agent-core/agents. If the profile root is
missing or mounted at the wrong path, startup should fail before a MachineSpec
run begins.

<!-- field: actor -->
## Actor

Operator running an agent-core runtime with external profiles.

<!-- field: trigger -->
## Trigger

An operator runs an installed agent-core image or local agent binary with a
profile path supplied from applications/catalog:

  agent --profile "agents/executor/profile.yaml" \
        --directory <workspace>

Container callers use the same contract with /profiles:

  agent --profile /profiles/agents/specification-critic/profile.yaml --directory /work

<!-- field: preconditions -->
## Preconditions

- applications/catalog is available as the demo.yaml catalog_root, /profiles, or a release bundle.
- agent-core supplies the agent binary and core-owned shared tool roots.
- The selected profile's workspace is mounted or available through --directory.

<!-- field: flow -->
## Flow

- F1: The operator selects a profile path from applications/catalog.
- F2: The runtime starts with --profile pointing at that external profile.
- F3: The profile resolves machine.yaml, tools.yaml, declarations, REST definitions, and UI assets relative to its profile directory.
- F4: Shared tool declarations resolve through documented core-owned roots such as /opt/agent-core/tools.
- F5: Catalog-owned profiles and integration fixtures are read from applications/catalog when integration targets set catalog_root in demo.yaml; application presentations remain under their owning application roots.
- F6: A mounted container run uses /profiles for agent programs and /work for the workspace.
- F7: Execution fails early if the runtime image requires copied profile directories under /opt/agent-core/agents.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd001-agent-functional-blocks](../software-requirements/srd001-agent-functional-blocks.md) [R1](../software-requirements/srd001-agent-functional-blocks.md#r1---profile-ownership-boundary) | Profile-owned files stay in applications/catalog while runtime assets stay in agent-core. |
| T2 | per-family SRDs srd002-executor through srd011-knowledge-manager | Mounted execution must preserve each functional block's actor boundary. |
| T3 | sm-multi-agent-profile-system system_boundary | applications/catalog owns configured programs and agent-core owns execution. |
| T4 | sm-agent-role-taxonomy checks | Role checks name the demo.yaml catalog_root and mounted profile execution. |
| T5 | mage containerSmoke | Container proof mounts this repository at /profiles and runs the specification-critic profile. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Runtime commands point --profile at the demo.yaml catalog_root or /profiles.

<!-- field: traces -->
#### Traces

- [srd001-agent-functional-blocks](../software-requirements/srd001-agent-functional-blocks.md) [R1.1](../software-requirements/srd001-agent-functional-blocks.md#r11)
- sm-multi-agent-profile-system profile_asset_rules

<!-- sc: S2 -->
### S2

Profile-local machine, tool, declaration, REST, and UI references resolve from the selected profile directory.

<!-- field: traces -->
#### Traces

- [srd001-agent-functional-blocks](../software-requirements/srd001-agent-functional-blocks.md) [R1.1](../software-requirements/srd001-agent-functional-blocks.md#r11)
- sm-agent-role-taxonomy checks

<!-- sc: S3 -->
### S3

Core-owned shared tool roots resolve without copying agent profiles into the runtime image.

<!-- field: traces -->
#### Traces

- [srd001-agent-functional-blocks](../software-requirements/srd001-agent-functional-blocks.md) [R1.2](../software-requirements/srd001-agent-functional-blocks.md#r12)
- [srd001-agent-functional-blocks](../software-requirements/srd001-agent-functional-blocks.md) [R1.3](../software-requirements/srd001-agent-functional-blocks.md#r13)

<!-- sc: S4 -->
### S4

Mounted container smoke execution uses /profiles and /work.

<!-- field: traces -->
#### Traces

- sm-multi-agent-profile-system validation_views

<!-- sc: S5 -->
### S5

Integration targets can consume canonical profiles and catalog-owned fixtures through the demo.yaml catalog_root without making application presentations catalog-owned.

<!-- field: traces -->
#### Traces

- [srd003-critic](../software-requirements/srd003-critic.md) [R2.3](../software-requirements/srd003-critic.md#r23)
- [srd008-runtime-state-reader](../software-requirements/srd008-runtime-state-reader.md) [R2.1](../software-requirements/srd008-runtime-state-reader.md#r21)
- [srd011-knowledge-manager](../software-requirements/srd011-knowledge-manager.md) [R1.2](../software-requirements/srd011-knowledge-manager.md#r12)

<!-- field: out_of_scope -->
## Out of scope

- Building the agent-core runtime image.
- Implementing builtin tool behavior.
- Changing profile YAML behavior during this documentation update.

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel06.0-agent-profiles](../test-suites/test-rel06.0-agent-profiles.md)
