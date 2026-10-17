---
id: rel03.0-uc001-package-serving-profiles
title: Package Serving Profiles for Separate Containers
status: implemented
test_suite: test-rel03.0-deployment-profile-package
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A release builder resolves each application serving entry profile through the
canonical #847 closure walker and publishes three independent profile roots.
Each root is ready to mount at /profiles in one role container and carries
deterministic ConfigMap partition metadata for the later Helm chart.

<!-- field: actor -->
## Actor

A release builder producing coding-agent deployment artifacts.

<!-- field: trigger -->
## Trigger

Run mage package from applications/coding-agent.

<!-- field: preconditions -->
## Preconditions

- agents/application.yaml names the three serving entry profiles.
- The compatible applications/catalog checkout is available.
- Application and canonical source trees contain no symlinks.

<!-- field: flow -->
## Flow

- F1: Validate canonical and serving references, runtime paths, and the profile-free image invariant.
- F2: Stage application and canonical sources into one temporary virtual profile namespace.
- F3: Resolve planner, executor, and critic independently with the existing transitive closure walker.
- F4: Sort each file set and partition encoded ConfigMap payload without splitting an individual entry.
- F5: Write role manifests and the deployment layout manifest, then atomically publish the complete output.
- F6: Boot-validate each role's packaged entry profile from its own root.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd002-deployment-profile-package](../software-requirements/srd002-deployment-profile-package.md) [R1.1](../software-requirements/srd002-deployment-profile-package.md#r11), [R1.2](../software-requirements/srd002-deployment-profile-package.md#r12), [R1.3](../software-requirements/srd002-deployment-profile-package.md#r13) | S1 traces this requirement in the existing success criteria. |
| T2 | [srd002-deployment-profile-package](../software-requirements/srd002-deployment-profile-package.md) [R2.1](../software-requirements/srd002-deployment-profile-package.md#r21), [R2.2](../software-requirements/srd002-deployment-profile-package.md#r22), [R2.3](../software-requirements/srd002-deployment-profile-package.md#r23), [R2.4](../software-requirements/srd002-deployment-profile-package.md#r24) | S2 traces this requirement in the existing success criteria. |
| T3 | [srd002-deployment-profile-package](../software-requirements/srd002-deployment-profile-package.md) [R3.1](../software-requirements/srd002-deployment-profile-package.md#r31), [R3.2](../software-requirements/srd002-deployment-profile-package.md#r32), [R3.3](../software-requirements/srd002-deployment-profile-package.md#r33) | S3 traces this requirement in the existing success criteria. |
| T4 | [srd002-deployment-profile-package](../software-requirements/srd002-deployment-profile-package.md) [R4.1](../software-requirements/srd002-deployment-profile-package.md#r41), [R4.2](../software-requirements/srd002-deployment-profile-package.md#r42) | S4 traces this requirement in the existing success criteria. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Each shard contains only one serving role and its reachable canonical assets.

<!-- field: traces -->
#### Traces

- [srd002-deployment-profile-package](../software-requirements/srd002-deployment-profile-package.md) [R1.1](../software-requirements/srd002-deployment-profile-package.md#r11)
- [srd002-deployment-profile-package](../software-requirements/srd002-deployment-profile-package.md) [R1.2](../software-requirements/srd002-deployment-profile-package.md#r12)
- [srd002-deployment-profile-package](../software-requirements/srd002-deployment-profile-package.md) [R1.3](../software-requirements/srd002-deployment-profile-package.md#r13)
- [srd002-deployment-profile-package](../software-requirements/srd002-deployment-profile-package.md) [AC1](../software-requirements/srd002-deployment-profile-package.md#ac1)

<!-- sc: S2 -->
### S2

Repeated packaging from the same inputs is byte-identical.

<!-- field: traces -->
#### Traces

- [srd002-deployment-profile-package](../software-requirements/srd002-deployment-profile-package.md) [R2.1](../software-requirements/srd002-deployment-profile-package.md#r21)
- [srd002-deployment-profile-package](../software-requirements/srd002-deployment-profile-package.md) [R2.2](../software-requirements/srd002-deployment-profile-package.md#r22)
- [srd002-deployment-profile-package](../software-requirements/srd002-deployment-profile-package.md) [R2.3](../software-requirements/srd002-deployment-profile-package.md#r23)
- [srd002-deployment-profile-package](../software-requirements/srd002-deployment-profile-package.md) [R2.4](../software-requirements/srd002-deployment-profile-package.md#r24)
- [srd002-deployment-profile-package](../software-requirements/srd002-deployment-profile-package.md) [AC1](../software-requirements/srd002-deployment-profile-package.md#ac1)

<!-- sc: S3 -->
### S3

ConfigMap partition metadata covers every role file once within the payload bound.

<!-- field: traces -->
#### Traces

- [srd002-deployment-profile-package](../software-requirements/srd002-deployment-profile-package.md) [R3.1](../software-requirements/srd002-deployment-profile-package.md#r31)
- [srd002-deployment-profile-package](../software-requirements/srd002-deployment-profile-package.md) [R3.2](../software-requirements/srd002-deployment-profile-package.md#r32)
- [srd002-deployment-profile-package](../software-requirements/srd002-deployment-profile-package.md) [R3.3](../software-requirements/srd002-deployment-profile-package.md#r33)
- [srd002-deployment-profile-package](../software-requirements/srd002-deployment-profile-package.md) [AC2](../software-requirements/srd002-deployment-profile-package.md#ac2)

<!-- sc: S4 -->
### S4

All packaged role entry profiles validate with the real binary.

<!-- field: traces -->
#### Traces

- [srd002-deployment-profile-package](../software-requirements/srd002-deployment-profile-package.md) [R4.1](../software-requirements/srd002-deployment-profile-package.md#r41)
- [srd002-deployment-profile-package](../software-requirements/srd002-deployment-profile-package.md) [R4.2](../software-requirements/srd002-deployment-profile-package.md#r42)
- [srd002-deployment-profile-package](../software-requirements/srd002-deployment-profile-package.md) [AC3](../software-requirements/srd002-deployment-profile-package.md#ac3)

<!-- field: out_of_scope -->
## Out of scope

- Helm templates and values schema.
- Live Kubernetes installation or smoke testing.
