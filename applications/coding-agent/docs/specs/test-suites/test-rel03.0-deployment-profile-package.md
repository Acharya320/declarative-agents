---
id: test-rel03.0-deployment-profile-package
title: Role-Specific Deployment Profile Package Suite
release: "03.0"
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

Deterministic unit tests prove exact role closures, reproducible manifests,
ConfigMap boundaries, and strict source failures. The application audit
packages the same artifacts and validates every role entry profile with the
real production agent binary.

<!-- field: traces -->
## Traces

- rel03.0-uc001-package-serving-profiles

<!-- field: preconditions -->
## Preconditions

- agents/application.yaml names the three serving entry profiles.
- The compatible applications/catalog checkout is available.
- Application and canonical source trees contain no symlinks.

<!-- field: test_cases -->
## Test cases

<!-- case: Exact deterministic role closures | status: implemented -->
### Exact deterministic role closures

<!-- field: use_case -->
#### Use case

rel03.0-uc001-package-serving-profiles

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestDeploymentPackage(ContainsExactRoleClosures|IsDeterministic)$'

<!-- field: traces -->
#### Traces

- srd002-deployment-profile-package AC1

<!-- case: ConfigMap partition boundaries | status: implemented -->
### ConfigMap partition boundaries

<!-- field: use_case -->
#### Use case

rel03.0-uc001-package-serving-profiles

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestConfigMapPartition'

<!-- field: traces -->
#### Traces

- srd002-deployment-profile-package AC2

<!-- case: Unsafe deployment sources fail | status: implemented -->
### Unsafe deployment sources fail

<!-- field: use_case -->
#### Use case

rel03.0-uc001-package-serving-profiles

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestDeployment(SourceRejectsSymlink|ServingReference)'

<!-- field: traces -->
#### Traces

- srd002-deployment-profile-package AC2

<!-- case: Profile-free packaged runtime validates | status: implemented -->
### Profile-free packaged runtime validates

<!-- field: use_case -->
#### Use case

rel03.0-uc001-package-serving-profiles

<!-- field: integration_target -->
#### Integration target

mage packageValidate

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestDeploymentManifestPreservesProfileFreeRuntime$'

<!-- field: traces -->
#### Traces

- srd002-deployment-profile-package AC3
