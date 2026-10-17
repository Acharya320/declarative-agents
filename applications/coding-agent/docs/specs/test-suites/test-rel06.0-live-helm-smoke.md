---
id: test-rel06.0-live-helm-smoke
title: Packaged-Chart Live Helm Smoke Suite
release: "06.0"
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

Deterministic unit tests pin prerequisite skips, infrastructure probes,
semantic classification, bounded diagnostics, canonical profile-free image
construction, and kind fixtures. The live target supplies the Kubernetes
evidence: one canonical agent image ID across every agent and mock pod, the
executor's tools mounted from donors, and the model boundary served by the
canonical mock profile. The one-image construction and mock-profile evidence
are delivered by the smoke refactor (GH-2496) and the model-mock deployment
(GH-2501); this suite states the contract those units satisfy.

<!-- field: traces -->
## Traces

- [rel06.0-uc001-live-helm-smoke](../use-cases/rel06.0-uc001-live-helm-smoke.md)

<!-- field: preconditions -->
## Preconditions

- Docker, kind, Helm, kubectl, and the required local images are available.
- The agent-core and applications/catalog checkouts are present.

<!-- field: test_cases -->
## Test cases

<!-- case: Missing prerequisite skip | status: implemented -->
### Missing prerequisite skip

<!-- field: use_case -->
#### Use case

[rel06.0-uc001-live-helm-smoke](../use-cases/rel06.0-uc001-live-helm-smoke.md)

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestCodingHelmSmokeSkipsMissingExternalPrerequisite$'

<!-- field: traces -->
#### Traces

- [srd005-live-helm-smoke](../software-requirements/srd005-live-helm-smoke.md) [AC2](../software-requirements/srd005-live-helm-smoke.md#ac2)

<!-- case: Failure classification and diagnostics | status: implemented -->
### Failure classification and diagnostics

<!-- field: use_case -->
#### Use case

[rel06.0-uc001-live-helm-smoke](../use-cases/rel06.0-uc001-live-helm-smoke.md)

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestCodingHelm(InfrastructureProbeClassification|FailureDistinguishesSemanticFromInfrastructure|DiagnosticsUseOneBoundedContext)$'

<!-- field: traces -->
#### Traces

- [srd005-live-helm-smoke](../software-requirements/srd005-live-helm-smoke.md) [AC3](../software-requirements/srd005-live-helm-smoke.md#ac3)

<!-- case: Canonical image and kind fixtures | status: implemented -->
### Canonical image and kind fixtures

<!-- field: use_case -->
#### Use case

[rel06.0-uc001-live-helm-smoke](../use-cases/rel06.0-uc001-live-helm-smoke.md)

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestCodingHelm(FixturesMatchServingAndStorageContract|SmokeUsesProductionProfileFreeImageRecipe)$'

<!-- field: traces -->
#### Traces

- [srd005-live-helm-smoke](../software-requirements/srd005-live-helm-smoke.md) [R2](../software-requirements/srd005-live-helm-smoke.md#r2---real-packaged-topology)

<!-- case: Full packaged topology | status: implemented -->
### Full packaged topology

<!-- field: use_case -->
#### Use case

[rel06.0-uc001-live-helm-smoke](../use-cases/rel06.0-uc001-live-helm-smoke.md)

<!-- field: integration_target -->
#### Integration target

mage integration:helmSmoke

<!-- field: traces -->
#### Traces

- [srd005-live-helm-smoke](../software-requirements/srd005-live-helm-smoke.md) [AC1](../software-requirements/srd005-live-helm-smoke.md#ac1)
