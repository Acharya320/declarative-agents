---
id: test-rel20.0-agent-operator
title: Release 20.0 Agent Operator Test Suite
release: "20.0"
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: test_case_count | form: "scalar" -->
## Test case count

5

<!-- field: overview -->
## Overview

TC1 is the structural gate for the release 20.0 documents. TC2 and TC3 are
the operator's unit tests: reconciliation of an Agent into its Deployment
and Service, with status, and the CRD applied on boot. TC4 proves the
platform boots with the operator ready, and TC5 is the collector proof on
the kind platform, which the release gates run.

<!-- field: traces -->
## Traces

- [rel20.0-uc001-agent-operator](../use-cases/rel20.0-uc001-agent-operator.md)

<!-- field: srds -->
## SRDs

- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md)

<!-- field: preconditions -->
## Preconditions

- Run commands from the repository root.
- Go, Mage, Helm, kubectl, kind, and Docker are available for the platform cases.

<!-- field: test_cases -->
## Test cases

<!-- case: Release 20.0 documents parse and trace reciprocally | status: implemented -->
### Release 20.0 documents parse and trace reciprocally

<!-- field: id -->
#### ID

TC1

<!-- field: use_case -->
#### Use case

[rel20.0-uc001-agent-operator](../use-cases/rel20.0-uc001-agent-operator.md)

<!-- field: evidence_type -->
#### Evidence type

executable_structural

<!-- field: implemented_by -->
#### Implemented by

GH-2682

<!-- field: command -->
#### Command

```bash
cd agent-core && mage audit
```

<!-- field: expected | form: "tree" -->
#### Expected

- audit_passes: true
- srd009_use_case_and_suite_trace_each_other: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc001-agent-operator](../use-cases/rel20.0-uc001-agent-operator.md) [S1](../use-cases/rel20.0-uc001-agent-operator.md#s1)
- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC1](../software-requirements/srd009-agent-operator.md#ac1)

<!-- case: The reconciler creates, rolls, and removes the Deployment and Service and reports status | status: implemented -->
### The reconciler creates, rolls, and removes the Deployment and Service and reports status

<!-- field: id -->
#### ID

TC2

<!-- field: use_case -->
#### Use case

[rel20.0-uc001-agent-operator](../use-cases/rel20.0-uc001-agent-operator.md)

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2683

<!-- field: command -->
#### Command

```bash
cd operator && go test ./internal/controller -run 'TestReconcile'
```

<!-- field: expected | form: "tree" -->
#### Expected

- deployment_and_service_created_with_owner_references: true
- nested_configmap_keys_mount_as_directories: true
- image_change_rolls_the_deployment: true
- missing_configmap_sets_ready_false_with_reason: true
- available_deployment_sets_ready_true: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc001-agent-operator](../use-cases/rel20.0-uc001-agent-operator.md) [S2](../use-cases/rel20.0-uc001-agent-operator.md#s2)
- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC2](../software-requirements/srd009-agent-operator.md#ac2)

<!-- case: The operator applies its CRD on boot | status: implemented -->
### The operator applies its CRD on boot

<!-- field: id -->
#### ID

TC3

<!-- field: use_case -->
#### Use case

[rel20.0-uc001-agent-operator](../use-cases/rel20.0-uc001-agent-operator.md)

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2683

<!-- field: command -->
#### Command

```bash
cd operator && go test ./internal/crd -run 'TestCRD'
```

<!-- field: expected | form: "tree" -->
#### Expected

- crd_generated_from_the_agent_types: true
- apply_creates_then_updates: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc001-agent-operator](../use-cases/rel20.0-uc001-agent-operator.md) [S1](../use-cases/rel20.0-uc001-agent-operator.md#s1)
- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC1](../software-requirements/srd009-agent-operator.md#ac1)

<!-- case: The kind platform boots with the operator ready | status: implemented -->
### The kind platform boots with the operator ready

<!-- field: id -->
#### ID

TC4

<!-- field: use_case -->
#### Use case

[rel20.0-uc001-agent-operator](../use-cases/rel20.0-uc001-agent-operator.md)

<!-- field: evidence_type -->
#### Evidence type

integration

<!-- field: implemented_by -->
#### Implemented by

GH-2684

<!-- field: command -->
#### Command

```bash
mage platform:up && kubectl --context kind-da-platform rollout status deployment/agent-operator --namespace agent-operator --timeout=180s
```

<!-- field: expected | form: "tree" -->
#### Expected

- operator_image_built_from_the_checkout: true
- chart_installed_after_fake_gcs: true
- rollout_complete: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc001-agent-operator](../use-cases/rel20.0-uc001-agent-operator.md) [S3](../use-cases/rel20.0-uc001-agent-operator.md#s3)
- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC3](../software-requirements/srd009-agent-operator.md#ac3)

<!-- case: The catalog collector Agent becomes Ready and answers its health endpoint | status: implemented -->
### The catalog collector Agent becomes Ready and answers its health endpoint

<!-- field: id -->
#### ID

TC5

<!-- field: use_case -->
#### Use case

[rel20.0-uc001-agent-operator](../use-cases/rel20.0-uc001-agent-operator.md)

<!-- field: evidence_type -->
#### Evidence type

integration

<!-- field: implemented_by -->
#### Implemented by

GH-2685

<!-- field: command -->
#### Command

```bash
cd operator && mage integration:all
```

<!-- field: expected | form: "tree" -->
#### Expected

- configmap_from_collector_profile_directory: true
- agent_ready_condition_true: true
- health_answers_through_the_service: true
- delete_removes_deployment_and_service: true

<!-- field: traces -->
#### Traces

- [rel20.0-uc001-agent-operator](../use-cases/rel20.0-uc001-agent-operator.md) [S4](../use-cases/rel20.0-uc001-agent-operator.md#s4)
- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC4](../software-requirements/srd009-agent-operator.md#ac4)
