---
id: test-rel05.0-helm-package-validation
title: Helm Values and Package Validation Suite
release: "05.0"
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

Schema fixtures, schema-bypass guards, manifest tamper cases, topology
combinations, and archive inventory/render tests form the deterministic
release gate. No test installs Kubernetes resources.

<!-- field: traces -->
## Traces

- rel05.0-uc001-package-validated-chart

<!-- field: test_cases -->
## Test cases

<!-- case: Values schema matrix | status: implemented -->
### Values schema matrix

<!-- field: use_case -->
#### Use case

rel05.0-uc001-package-validated-chart

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestHelmSchemaFixtureMatrix$'

<!-- field: traces -->
#### Traces

- srd004-helm-package-validation AC1

<!-- case: Semantic guards and topology combinations | status: implemented -->
### Semantic guards and topology combinations

<!-- field: use_case -->
#### Use case

rel05.0-uc001-package-validated-chart

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestHelm(SemanticValidationSurvivesSchemaBypass|WorkspaceAndTelemetryCombinationsRender)$'

<!-- field: traces -->
#### Traces

- srd004-helm-package-validation R1.3

<!-- case: Prepared package tamper rejection | status: implemented -->
### Prepared package tamper rejection

<!-- field: use_case -->
#### Use case

rel05.0-uc001-package-validated-chart

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestPreparedPackage'

<!-- field: traces -->
#### Traces

- srd004-helm-package-validation AC2

<!-- case: Self-contained chart archive | status: implemented -->
### Self-contained chart archive

<!-- field: use_case -->
#### Use case

rel05.0-uc001-package-validated-chart

<!-- field: integration_target -->
#### Integration target

mage helm:package

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^Test(HelmPackageArchiveIsSelfContained|ChartArchiveValidationRejectsOmissions)$'

<!-- field: traces -->
#### Traces

- srd004-helm-package-validation AC3

<!-- field: preconditions -->
## Preconditions

- Helm is available.
- The
