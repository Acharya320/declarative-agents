---
id: rel06.0-uc001-mounted-profile-container-runtime
title: Runtime container runs an externally mounted agent profile
test_suite: test-rel06.0-profile-migration
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A release maintainer builds an Agent Core runtime image. It carries the agent
binary plus core-owned assets, not a bundled catalog of agent profiles.

At run time, the selected program comes from a mounted profile repository.
Commands point --profile at that mounted file and --directory at the
workspace mounted for the task.

<!-- field: actor -->
## Actor

Release maintainer proving mounted-profile container execution

<!-- field: trigger -->
## Trigger

Agent profile assets have moved to applications/catalog. Container users still need
a simple command that runs configured agents. That command depends on an
explicit profile mount, not on profile directories copied into the Agent Core
image.

<!-- field: runtime_notes -->
## Runtime notes

The proof is intentionally ordinary.

Operators should see the command they would paste into a release note:
mount the profile repository, mount the workspace, name the profile file, and
let the existing runtime load the program from that path. No hidden image
fallback should make the proof pass.

Failure matters too. Before core.Loop runs, the runtime image entrypoint
rejects profile-less startup. In Mage and package tests, missing external
roots fail with diagnostics that name AGENT_CATALOG_ROOT.

On purpose, the container remains boring. It owns the executable runtime,
while the mounted profile owns the program under test.

<!-- field: flow -->
## Flow

- F1: The maintainer builds the Agent Core runtime image from agent-core.
- F2: The build keeps the agent binary and core-owned shared tool declarations under /opt/agent-core.
- F3: The image build does not copy agent profile directories from agent-core.
- F4: The maintainer checks out or unpacks applications/catalog outside the runtime image.
- F5: The container command mounts applications/catalog at /profiles and the target workspace at /work.
- F6: The container command starts agent with --profile /profiles/agents/<agent>/profile.yaml and --directory /work.
- F7: profile.yaml resolves its machine, tools, tool declarations, REST definitions, request machines, and UX config relative to the mounted profile directory.
- F8: The runtime loads core-owned shared declarations from /opt/agent-core when the mounted profile references them.
- F9: Missing mounted profile paths fail before the loop starts with a diagnostic that names the missing path.

<!-- field: proof -->
## Proof

- TestResolveProfilePathMapsInstalledCoreHome maps /opt/agent-core/tools to AGENT_CORE_HOME for installed runtime assets.
- TestDockerfileRuntimeExcludesAgentProfiles verifies that the runtime image copies /opt/agent-core/tools and does not copy migrated profile directories.
- TestContainerBuildSummaryForDocker documents the mounted /profiles/agents runtime command emitted by mage docker.
- TestDockerfileDefinesIntegrationTarget verifies the Dockerfile integration target used for Mage integration runs inside the image.
- applications/coding-agent mage integration:executorLive packages and consumes the canonical executor profile through AGENT_CATALOG_ROOT.
- Profile-owned documentation-curator and browser UX tracer bullets belong to applications/catalog, where the concrete profile and profile-local UX config live.
- applications/catalog mage containerSmoke runs the installed runtime image with the profile repository mounted at /profiles.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | Dockerfile | runtime image contents |
| T2 | README.md | mounted-profile container command |
| T3 | mage docker | release-image build wrapper |
| T4 | [srd034-external-agent-profiles](../software-requirements/srd034-external-agent-profiles.md) [R1.1](../software-requirements/srd034-external-agent-profiles.md#r11), [R1.2](../software-requirements/srd034-external-agent-profiles.md#r12), [R1.3](../software-requirements/srd034-external-agent-profiles.md#r13), [R1.4](../software-requirements/srd034-external-agent-profiles.md#r14), [R2.1](../software-requirements/srd034-external-agent-profiles.md#r21), [R2.2](../software-requirements/srd034-external-agent-profiles.md#r22), [R2.3](../software-requirements/srd034-external-agent-profiles.md#r23), [R2.4](../software-requirements/srd034-external-agent-profiles.md#r24), [R3.1](../software-requirements/srd034-external-agent-profiles.md#r31), [R3.2](../software-requirements/srd034-external-agent-profiles.md#r32), [R3.3](../software-requirements/srd034-external-agent-profiles.md#r33), [R3.4](../software-requirements/srd034-external-agent-profiles.md#r34), [R4.1](../software-requirements/srd034-external-agent-profiles.md#r41), [R4.2](../software-requirements/srd034-external-agent-profiles.md#r42), [R4.3](../software-requirements/srd034-external-agent-profiles.md#r43), [R4.4](../software-requirements/srd034-external-agent-profiles.md#r44), [R4.5](../software-requirements/srd034-external-agent-profiles.md#r45) | ownership and mounted runtime contract |
| T5 | knowledge-manager/documentation-curator/profile.yaml | profile-relative asset example under /profiles/agents |
| T6 | knowledge-manager/documentation-curator/ui/ux.yaml | profile-local UX config example under /profiles/agents |

<!-- field: implementation_status -->
## Implementation status

implemented

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The runtime image contains the agent binary and core-owned shared assets, but it does not contain copied agent profile directories.

<!-- field: traces -->
#### Traces

- [srd034-external-agent-profiles](../software-requirements/srd034-external-agent-profiles.md) [AC1](../software-requirements/srd034-external-agent-profiles.md#ac1)
- [srd034-external-agent-profiles](../software-requirements/srd034-external-agent-profiles.md) [AC3](../software-requirements/srd034-external-agent-profiles.md#ac3)

<!-- sc: S2 -->
### S2

A container command can mount applications/catalog and run agent with --profile /profiles/agents/... plus --directory /work.

<!-- field: traces -->
#### Traces

- [srd034-external-agent-profiles](../software-requirements/srd034-external-agent-profiles.md) [AC2](../software-requirements/srd034-external-agent-profiles.md#ac2)
- [srd034-external-agent-profiles](../software-requirements/srd034-external-agent-profiles.md) [AC4](../software-requirements/srd034-external-agent-profiles.md#ac4)

<!-- sc: S3 -->
### S3

Profile-local paths continue to resolve from the mounted profile directory.

<!-- field: traces -->
#### Traces

- [srd034-external-agent-profiles](../software-requirements/srd034-external-agent-profiles.md) [AC2](../software-requirements/srd034-external-agent-profiles.md#ac2)

<!-- sc: S4 -->
### S4

Missing mounted profile paths fail with a clear pre-loop diagnostic.

<!-- field: traces -->
#### Traces

- [srd034-external-agent-profiles](../software-requirements/srd034-external-agent-profiles.md) [AC2](../software-requirements/srd034-external-agent-profiles.md#ac2)

<!-- field: out_of_scope -->
## Out of scope

- Building or publishing applications/catalog release artifacts.
- Choosing between bind mounts, submodules, release tarballs, or derived images as the only profile distribution model.
- Reintroducing in-repo profile assets after equivalent files exist in applications/catalog.
- Rewriting Agent Core runtime packages during the profile migration.

<!-- field: srd -->
## SRD

srd034-external-agent-profiles
