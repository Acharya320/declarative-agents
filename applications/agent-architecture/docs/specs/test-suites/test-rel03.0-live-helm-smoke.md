---
id: test-rel03.0-live-helm-smoke
title: Test Suite for Release 03.0 Live Helm Smoke
release: "03.0"
status: planned
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

This suite validates the bounded kind smoke (srd001 R3): the packaged archive
install, real curator and collector readiness, a browsable curator surface,
clean lifecycle exit, an in-cluster trace the collector retains, prerequisite
skips, and cleanup with bounded diagnostics. The kind target lands in a later
sub-issue, so each case names a planned integration target rather than an
existing Go test.

<!-- field: traces -->
## Traces

- rel03.0-uc001-live-helm-smoke

<!-- field: preconditions -->
## Preconditions

- Docker, kind, Helm, and kubectl are available and Docker is healthy for the live case.
- The validated chart archive and the catalog and agent-core checkouts are available.

<!-- field: test_cases -->
## Test cases

<!-- case: The validated archive installs on kind | status: planned -->
### The validated archive installs on kind

Creates a kind cluster through a temporary kubeconfig, builds and loads the
shared profile-free agent-core image the chart references, and installs only the
validated archive. Asserts the install source is the validated archive, not
a source-tree surrogate, and that the packaged curator and collector
closures are the mounted profile source.

<!-- field: use_case -->
#### Use case

rel03.0-uc001-live-helm-smoke

<!-- field: id -->
#### ID

install_packaged_chart

<!-- field: integration_target -->
#### Integration target

mage integration:helmSmoke

<!-- field: inputs | form: "tree" -->
#### Inputs

- cluster: kind with the runtime image built and loaded
- install_source: validated chart archive

<!-- field: expected | form: "tree" -->
#### Expected

- installs_validated_archive_only: true
- profile_free_image_loaded: true
- packaged_closures_mounted: true

<!-- field: traces -->
#### Traces

- rel03.0-uc001-live-helm-smoke S1
- srd001-helm-deployment AC4

<!-- case: Curator and collector reach readiness and exit cleanly | status: planned -->
### Curator and collector reach readiness and exit cleanly

Waits for the curator and collector Deployments and their real lifecycle
health endpoints to become ready, reaches the curator documentation
interface, and issues a declarative lifecycle-exit request. Asserts the
curator stops cleanly and releases its listeners rather than being killed.

<!-- field: use_case -->
#### Use case

rel03.0-uc001-live-helm-smoke

<!-- field: id -->
#### ID

readiness_and_exit

<!-- field: integration_target -->
#### Integration target

mage integration:helmSmoke

<!-- field: inputs | form: "tree" -->
#### Inputs

- workloads:
  1. curator
  1. collector

<!-- field: expected | form: "tree" -->
#### Expected

- real_lifecycle_health_ready: true
- curator_surface_served: true
- declarative_exit_stops_curator_cleanly: true

<!-- field: traces -->
#### Traces

- rel03.0-uc001-live-helm-smoke S2
- srd001-helm-deployment AC4

<!-- case: The collector retains an exported curator trace | status: planned -->
### The collector retains an exported curator trace

Drives the curator to export a trace under its service name and queries the
collector's spool surface for it. Asserts the collector retains the
curator trace, proving the in-cluster OTLP path end to end rather than by
render inspection alone.

<!-- field: use_case -->
#### Use case

rel03.0-uc001-live-helm-smoke

<!-- field: id -->
#### ID

trace_retained

<!-- field: integration_target -->
#### Integration target

mage integration:helmSmoke

<!-- field: inputs | form: "tree" -->
#### Inputs

- service_name: knowledge-manager-curator

<!-- field: expected | form: "tree" -->
#### Expected

- collector_retains_curator_trace: true

<!-- field: traces -->
#### Traces

- rel03.0-uc001-live-helm-smoke S3
- srd001-helm-deployment AC4

<!-- case: Prerequisite skips and terminal cleanup | status: planned -->
### Prerequisite skips and terminal cleanup

Drives missing-tool preflight cases and asserts the target prints a
specific skip without mutating a cluster. Drives terminal paths and asserts
the namespace, release, forwards, temporary files, and any owned cluster
are cleaned, and that failure diagnostics share one deadline and include
Helm status, workloads, events, and targeted curator and collector logs.

<!-- field: use_case -->
#### Use case

rel03.0-uc001-live-helm-smoke

<!-- field: id -->
#### ID

skip_and_cleanup

<!-- field: integration_target -->
#### Integration target

mage test (smoke preflight and cleanup)

<!-- field: inputs | form: "tree" -->
#### Inputs

- missing_tools:
  1. docker
  1. kind
  1. helm
  1. kubectl
- terminal_paths:
  1. success
  1. semantic_failure
  1. infrastructure_failure

<!-- field: expected | form: "tree" -->
#### Expected

- missing_prerequisite_skips_without_cluster_mutation: true
- owned_resources_cleaned_on_every_terminal_path: true
- failure_diagnostics_bounded_and_targeted: true

<!-- field: traces -->
#### Traces

- rel03.0-uc001-live-helm-smoke S4
- srd001-helm-deployment AC4
