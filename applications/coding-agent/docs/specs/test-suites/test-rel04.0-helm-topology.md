---
id: test-rel04.0-helm-topology
title: Core Coding-Agent Helm Topology Suite
release: "04.0"
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

Focused deterministic renders prove role topology, package-driven ConfigMaps,
scoped checksums, shared storage, lifecycle probes, fixed URLs, security,
model modes, and OTLP routing. Schema matrices and live installation remain
follow-up suites.

<!-- field: traces -->
## Traces

- rel04.0-uc001-render-helm-topology

<!-- field: preconditions -->
## Preconditions

- The
- Helm is available.

<!-- field: test_cases -->
## Test cases

<!-- case: Prepared package is chart source | status: implemented -->
### Prepared package is chart source

<!-- field: use_case -->
#### Use case

rel04.0-uc001-render-helm-topology

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestHelmPrepare'

<!-- field: traces -->
#### Traces

- srd003-helm-topology R1.3
- srd003-helm-topology R4.1

<!-- case: Core role topology | status: implemented -->
### Core role topology

<!-- field: use_case -->
#### Use case

rel04.0-uc001-render-helm-topology

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestHelmCoreTopologyRendersRoleContract$'

<!-- field: traces -->
#### Traces

- srd003-helm-topology AC1

<!-- case: Role-scoped profile rollout | status: implemented -->
### Role-scoped profile rollout

<!-- field: use_case -->
#### Use case

rel04.0-uc001-render-helm-topology

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestHelmProfileChangeRollsOnlyAffectedRole$'

<!-- field: traces -->
#### Traces

- srd003-helm-topology AC2

<!-- case: Small and in-cluster model renders | status: implemented -->
### Small and in-cluster model renders

<!-- field: use_case -->
#### Use case

rel04.0-uc001-render-helm-topology

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestHelmSmallAndOllamaValuesRender$'

<!-- field: traces -->
#### Traces

- srd003-helm-topology AC3
- srd003-helm-topology AC4
