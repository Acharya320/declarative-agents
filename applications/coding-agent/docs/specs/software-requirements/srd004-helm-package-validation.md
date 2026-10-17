---
id: srd004-helm-package-validation
title: Helm Values and Package Validation
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

A rendered source chart is insufficient release evidence. Values must reject
unsafe or incoherent topology, prepared profile packages must resist stale or
tampered metadata and content, and the installable archive must include every
role asset and render independently of the repository.

<!-- req: R1 -->
## R1 - Strict values

<!-- req: R1.1 -->
### R1.1

values.schema.json shall reject unknown keys and malformed image, replica, port, resource, storage, URL, mount, workspace, model, security, and telemetry values.

<!-- req: R1.2 -->
### R1.2

Schema conditions shall reject external model mode without an HTTP endpoint.

<!-- req: R1.3 -->
### R1.3

Template guards shall retain mount, fixed-role-port, telemetry, and model invariants when schema validation is bypassed.

<!-- req: R2 -->
## R2 - Package integrity

<!-- req: R2.1 -->
### R2.1

Role manifests shall carry a deterministic checksum over ordered path and content bytes.

<!-- req: R2.2 -->
### R2.2

Preparation shall strictly parse manifests and reject stale roles, unknown fields, missing or extra files, checksum mismatch, duplicate coverage, invalid partition sizes, and encoded-key conflicts.

<!-- req: R2.3 -->
### R2.3

Prepared package validation shall run before chart staging or packaging.

<!-- req: R3 -->
## R3 - Self-contained Helm archive

<!-- req: R3.1 -->
### R3.1

mage helm:package shall regenerate role packages, stage the chart, run default, CI, and valid-fixture renders, and create one chart archive.

<!-- req: R3.2 -->
### R3.2

Archive validation shall require Chart.yaml, values, schema, core templates, and every generated profile package file and shall reject links.

<!-- req: R3.3 -->
### R3.3

The resulting .tgz shall lint and render from a directory with no source checkout.

<!-- field: goals -->
## Goals

- G1: Strict values.
- G2: Package integrity.
- G3: Self-contained Helm archive.

<!-- field: non_goals -->
## Non-goals

- Live Kubernetes installation.
- kind cluster setup and teardown.
- Operator deployment and teardown guide.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

Every valid fixture lints and renders, and every invalid fixture fails with its intended field.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3

<!-- ac: AC2 -->
### AC2

Deterministic tamper tests reject stale manifests, files, partitions, checksums, conflicts, and omissions.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3

<!-- ac: AC3 -->
### AC3

A packaged .tgz contains all role assets and independently passes helm lint and helm template.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3
