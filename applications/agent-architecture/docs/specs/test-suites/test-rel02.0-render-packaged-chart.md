---
id: test-rel02.0-render-packaged-chart
title: Test Suite for Release 02.0 Packaged Chart Render and Validation
release: "02.0"
status: planned
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

This suite validates the packaged chart (srd001 R1 and R2): the rendered
curator and collector topology, the deterministic staging and checksum of the
catalog-owned closures, the values schema and template guards, and the
self-contained archive. The chart and its preparation and package targets land
in later sub-issues, so each case names a planned integration target rather
than an existing Go test.

<!-- field: traces -->
## Traces

- rel02.0-uc001-render-packaged-chart

<!-- field: preconditions -->
## Preconditions

- The demo.yaml catalog_root resolves to a catalog checkout with the curator and collector profiles.
- Helm is available for lint and template renders.

<!-- field: test_cases -->
## Test cases

<!-- case: Curator and collector topology renders truthfully | status: planned -->
### Curator and collector topology renders truthfully

Renders the chart and asserts exactly one curator and one collector
Deployment and ClusterIP Service, the shared profile-free image, mounted
canonical closures, fixed request, control, and monitor or query ports,
lifecycle-health readiness and liveness probes, and the non-root,
no-privilege-escalation, dropped-capability, read-only-root, and
disabled-token security defaults.

<!-- field: use_case -->
#### Use case

rel02.0-uc001-render-packaged-chart

<!-- field: id -->
#### Id

render_topology

<!-- field: integration_target -->
#### Integration target

mage helm:template (default and CI values)

<!-- field: inputs -->
#### Inputs

```yaml
values: [default, ci]
```

<!-- field: expected -->
#### Expected

```yaml
one_curator_and_one_collector_workload: true
profile_free_image_with_mounted_closures: true
truthful_lifecycle_probes: true
hardened_pod_security: true
```

<!-- field: traces -->
#### Traces

- rel02.0-uc001-render-packaged-chart S1
- srd001-helm-deployment AC1

<!-- case: Prepared closures are the chart source | status: planned -->
### Prepared closures are the chart source

Runs the preparation target and asserts it regenerates and atomically
stages the catalog-owned curator and collector closures as the chart's
only profile source, and that strict validation rejects stale roles,
unknown fields, missing or extra files, checksum mismatch, and out-of-bound
ConfigMap partitions before staging.

<!-- field: use_case -->
#### Use case

rel02.0-uc001-render-packaged-chart

<!-- field: id -->
#### Id

stage_and_checksum

<!-- field: integration_target -->
#### Integration target

mage helm:prepare with tamper fixtures

<!-- field: inputs -->
#### Inputs

```yaml
tamper_cases: [stale_role, unknown_field, missing_file, extra_file, checksum_mismatch, oversized_partition]
```

<!-- field: expected -->
#### Expected

```yaml
closures_staged_from_catalog: true
tamper_cases_rejected_before_staging: true
```

<!-- field: traces -->
#### Traces

- rel02.0-uc001-render-packaged-chart S2
- srd001-helm-deployment AC2

<!-- case: Values schema, guards, and self-contained archive | status: planned -->
### Values schema, guards, and self-contained archive

Asserts the values.schema.json rejects unknown keys and malformed image,
replica, port, resource, storage, URL, mount, security, and telemetry
values by their offending field, that template guards retain mount,
fixed-port, and telemetry invariants when the schema is bypassed, and that
the packaged archive lints and renders from a directory with no source
checkout present.

<!-- field: use_case -->
#### Use case

rel02.0-uc001-render-packaged-chart

<!-- field: id -->
#### Id

schema_and_archive

<!-- field: integration_target -->
#### Integration target

mage helm:package with schema fixtures

<!-- field: inputs -->
#### Inputs

```yaml
valid_fixtures: [default, ci]
invalid_fixtures: [unknown_key, bad_port, bad_replica, missing_endpoint]
```

<!-- field: expected -->
#### Expected

```yaml
valid_fixtures_render: true
invalid_fixtures_rejected_by_path: true
guards_hold_under_schema_bypass: true
archive_renders_without_source_checkout: true
```

<!-- field: traces -->
#### Traces

- rel02.0-uc001-render-packaged-chart S3
- srd001-helm-deployment AC3

<!-- case: Curator-to-collector trace export renders coherently | status: planned -->
### Curator-to-collector trace export renders coherently

Asserts the curator receives a values-selected OTLP endpoint and service
name pointing at the in-cluster collector Service, and that the collector
renders its receiver and query surface, so the in-cluster trace path is
coherent under default and CI values.

<!-- field: use_case -->
#### Use case

rel02.0-uc001-render-packaged-chart

<!-- field: id -->
#### Id

trace_wiring

<!-- field: integration_target -->
#### Integration target

mage helm:template (telemetry assertions)

<!-- field: inputs -->
#### Inputs

```yaml
values: [default, ci]
```

<!-- field: expected -->
#### Expected

```yaml
curator_exports_to_collector_service: true
collector_receiver_and_query_render: true
```

<!-- field: traces -->
#### Traces

- rel02.0-uc001-render-packaged-chart S4
- srd001-helm-deployment AC1
