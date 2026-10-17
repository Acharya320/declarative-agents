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

- [rel05.0-uc001-package-validated-chart](../use-cases/rel05.0-uc001-package-validated-chart.md)

<!-- field: preconditions -->
## Preconditions

- Helm is available.
- The

<!-- field: test_cases -->
## Test cases

<!-- case: Values schema matrix | status: implemented -->
### Values schema matrix

<!-- field: use_case -->
#### Use case

[rel05.0-uc001-package-validated-chart](../use-cases/rel05.0-uc001-package-validated-chart.md)

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestHelmSchemaFixtureMatrix$'

<!-- field: traces -->
#### Traces

- [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [AC1](../software-requirements/srd004-helm-package-validation.md#ac1)

<!-- case: Semantic guards and topology combinations | status: implemented -->
### Semantic guards and topology combinations

<!-- field: use_case -->
#### Use case

[rel05.0-uc001-package-validated-chart](../use-cases/rel05.0-uc001-package-validated-chart.md)

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestHelm(SemanticValidationSurvivesSchemaBypass|WorkspaceAndTelemetryCombinationsRender)$'

<!-- field: traces -->
#### Traces

- [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [R1.3](../software-requirements/srd004-helm-package-validation.md#r13)

<!-- case: Prepared package tamper rejection | status: implemented -->
### Prepared package tamper rejection

<!-- field: use_case -->
#### Use case

[rel05.0-uc001-package-validated-chart](../use-cases/rel05.0-uc001-package-validated-chart.md)

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestPreparedPackage'

<!-- field: traces -->
#### Traces

- [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [AC2](../software-requirements/srd004-helm-package-validation.md#ac2)

<!-- case: Self-contained chart archive | status: implemented -->
### Self-contained chart archive

<!-- field: use_case -->
#### Use case

[rel05.0-uc001-package-validated-chart](../use-cases/rel05.0-uc001-package-validated-chart.md)

<!-- field: integration_target -->
#### Integration target

mage helm:package

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^Test(HelmPackageArchiveIsSelfContained|ChartArchiveValidationRejectsOmissions)$'

<!-- field: traces -->
#### Traces

- [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [AC3](../software-requirements/srd004-helm-package-validation.md#ac3)
