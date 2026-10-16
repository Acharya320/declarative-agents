---
id: test-rel07.0-collector-agent-tracing
title: Canonical Collector Agent Tracing Suite
release: "07.0"
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

The release re-proves the existing topology, schema, and smoke contracts
with the canonical collector agent as the only trace workload. Deterministic
unit tests pin the collector render, the schema fixture matrix, catalog-root
discovery, and the install command; the live smoke target supplies the
in-cluster trace evidence through the collector query surface.

<!-- field: traces -->
## Traces

- rel07.0-uc001-collector-agent-tracing

<!-- field: test_cases -->
## Test cases

<!-- case: Collector agent topology render | status: implemented -->
### Collector agent topology render

<!-- field: use_case -->
#### Use case

rel07.0-uc001-collector-agent-tracing

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestHelm(CoreTopologyRendersRoleContract|ProfileChangeRollsOnlyAffectedRole|SmallAndOllamaValuesRender)$'

<!-- field: traces -->
#### Traces

- srd003-helm-topology R3.3
- srd003-helm-topology R3.4

<!-- case: Schema and semantic validation without Jaeger keys | status: implemented -->
### Schema and semantic validation without Jaeger keys

<!-- field: use_case -->
#### Use case

rel07.0-uc001-collector-agent-tracing

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestHelm(SchemaFixtureMatrix|SemanticValidationSurvivesSchemaBypass|WorkspaceAndTelemetryCombinationsRender)$'

<!-- field: traces -->
#### Traces

- srd004-helm-package-validation R1

<!-- case: Catalog-root discovery from any startup directory | status: implemented -->
### Catalog-root discovery from any startup directory

<!-- field: use_case -->
#### Use case

rel07.0-uc001-collector-agent-tracing

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestCodingAgentCatalogRoot(FromDemoConfig|DiscoversSiblingCatalog|DiscoversFromRelativeStartupDir)$'

<!-- field: traces -->
#### Traces

- srd004-helm-package-validation R1

<!-- case: Install command without the removed utilityImage knob | status: implemented -->
### Install command without the removed utilityImage knob

<!-- field: use_case -->
#### Use case

rel07.0-uc001-collector-agent-tracing

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestCodingHelmCommitImagePropagatesToBuildManifestAndDeploy$'

<!-- field: traces -->
#### Traces

- srd005-live-helm-smoke R2

<!-- case: Live packaged topology with collector query trace proof | status: implemented -->
### Live packaged topology with collector query trace proof

<!-- field: use_case -->
#### Use case

rel07.0-uc001-collector-agent-tracing

<!-- field: integration_target -->
#### Integration target

mage integration:helmSmoke

<!-- field: traces -->
#### Traces

- srd005-live-helm-smoke R3.3

<!-- field: preconditions -->
## Preconditions

- The applications/catalog checkout is discoverable or the demo.yaml catalog_root names it.
- The prepared role packages and the staged collector profile are present in the chart.
