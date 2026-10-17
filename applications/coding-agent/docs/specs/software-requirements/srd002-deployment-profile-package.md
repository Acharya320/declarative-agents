---
id: srd002-deployment-profile-package
title: Role-Specific Deployment Profile Package
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Separate planner, executor, and critic containers must not receive the full
application profile tree. Deployment packaging needs one deterministic
resolver-backed closure per role, with provenance and bounded ConfigMap
partitions, without creating a second reference walker or baking profiles
into the shared runtime image.

<!-- req: R1 -->
## R1 - Role-specific closures

<!-- req: R1.1 -->
### R1.1

The application manifest shall name exactly planner, executor, and critic serving entry profiles and their /profiles runtime paths.

<!-- req: R1.2 -->
### R1.2

Packaging shall use the existing profileClosure resolver to follow application and canonical library references.

<!-- req: R1.3 -->
### R1.3

Each role shard shall contain its application entry profile and exactly the reachable application and canonical assets required by that profile.

<!-- req: R2 -->
## R2 - Safe reproducible package

<!-- req: R2.1 -->
### R2.1

Role paths, file lists, shard order, ConfigMap partitions, and YAML manifests shall be deterministic.

<!-- req: R2.2 -->
### R2.2

Traversal, absolute non-core references, globs, dangling references, symlinks, non-regular files, and conflicting destinations shall fail before publication.

<!-- req: R2.3 -->
### R2.3

Publication shall atomically replace the prior output only after every role succeeds.

<!-- req: R2.4 -->
### R2.4

Deployment and role manifests shall retain compatible library release, library and application checkout revisions and dirty state, mount path, and profile-free image metadata.

<!-- req: R3 -->
## R3 - ConfigMap partitions

<!-- req: R3.1 -->
### R3.1

Packaging shall partition each sorted role closure into deterministic chunks no larger than 900 KiB of encoded key and value payload.

<!-- req: R3.2 -->
### R3.2

One asset that exceeds the payload limit shall fail as unshardable rather than producing an invalid install artifact.

<!-- req: R3.3 -->
### R3.3

Role manifests shall list every partition, its encoded size, and its ordered file membership exactly once.

<!-- req: R4 -->
## R4 - Runtime proof

<!-- req: R4.1 -->
### R4.1

Every packaged serving entry profile shall pass agent-core --validate-config from its own role root.

<!-- req: R4.2 -->
### R4.2

Local combined-tree serving integrations shall remain unchanged.

<!-- field: goals -->
## Goals

- G1: Role-specific closures.
- G2: Safe reproducible package.
- G3: ConfigMap partitions.
- G4: Runtime proof.

<!-- field: non_goals -->
## Non-goals

- Helm templates and values schema.
- Live Kubernetes installation or smoke testing.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

mage package emits deterministic planner, executor, and critic role roots plus deployment and role manifests.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R2.1
- R2.2
- R2.3
- R2.4

<!-- ac: AC2 -->
### AC2

Boundary tests prove deterministic partitioning and reject dangling, symlinked, conflicting, and oversized inputs.

<!-- field: traces -->
#### Traces

- R2.2
- R3.1
- R3.2
- R3.3

<!-- ac: AC3 -->
### AC3

mage audit boot-validates every produced role profile with the real agent binary.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
