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

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Each shard contains only one serving role and its reachable canonical assets.

<!-- field: traces -->
#### Traces

- srd002-deployment-profile-package R1.1
- srd002-deployment-profile-package R1.2
- srd002-deployment-profile-package R1.3
- srd002-deployment-profile-package AC1

<!-- sc: S2 -->
### S2

Repeated packaging from the same inputs is byte-identical.

<!-- field: traces -->
#### Traces

- srd002-deployment-profile-package R2.1
- srd002-deployment-profile-package R2.2
- srd002-deployment-profile-package R2.3
- srd002-deployment-profile-package R2.4
- srd002-deployment-profile-package AC1

<!-- sc: S3 -->
### S3

ConfigMap partition metadata covers every role file once within the payload bound.

<!-- field: traces -->
#### Traces

- srd002-deployment-profile-package R3.1
- srd002-deployment-profile-package R3.2
- srd002-deployment-profile-package R3.3
- srd002-deployment-profile-package AC2

<!-- sc: S4 -->
### S4

All packaged role entry profiles validate with the real binary.

<!-- field: traces -->
#### Traces

- srd002-deployment-profile-package R4.1
- srd002-deployment-profile-package R4.2
- srd002-deployment-profile-package AC3

<!-- field: out_of_scope -->
## Out of scope

- Helm templates and values schema.
- Live Kubernetes installation or smoke testing.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd002-deployment-profile-package R1.1, R1.2, R1.3 | S1 traces this requirement in the existing success criteria. |
| T2 | srd002-deployment-profile-package R2.1, R2.2, R2.3, R2.4 | S2 traces this requirement in the existing success criteria. |
| T3 | srd002-deployment-profile-package R3.1, R3.2, R3.3 | S3 traces this requirement in the existing success criteria. |
| T4 | srd002-deployment-profile-package R4.1, R4.2 | S4 traces this requirement in the existing success criteria. |
