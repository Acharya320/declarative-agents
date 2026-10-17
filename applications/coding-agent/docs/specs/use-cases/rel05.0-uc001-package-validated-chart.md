---
id: rel05.0-uc001-package-validated-chart
title: Produce a Validated Self-Contained Chart
status: implemented
test_suite: test-rel05.0-helm-package-validation
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A release builder validates chart values and prepared role closures, renders
the supported configuration matrix, packages one Helm archive, checks its
complete file inventory, and renders that archive without source-tree files.

<!-- field: actor -->
## Actor

A release builder producing the coding-agent Helm artifact.

<!-- field: trigger -->
## Trigger

Run mage helm:package.

<!-- field: preconditions -->
## Preconditions

- Helm is available.
- The #875 role package can be regenerated.

<!-- field: flow -->
## Flow

- F1: Regenerate role packages with content checksums and validate manifest/file/partition consistency.
- F2: Stage those exact profiles into a temporary chart.
- F3: Lint and render defaults, CI values, and all valid schema fixtures.
- F4: Package the temporary chart and inspect its gzip/tar entries for required source and profile files.
- F5: Lint and render the resulting archive independently.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [R1.1](../software-requirements/srd004-helm-package-validation.md#r11), [R1.2](../software-requirements/srd004-helm-package-validation.md#r12), [R1.3](../software-requirements/srd004-helm-package-validation.md#r13) | S1 traces this requirement in the existing success criteria. |
| T2 | [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [AC1](../software-requirements/srd004-helm-package-validation.md#ac1) | S1 traces this requirement in the existing success criteria. |
| T3 | [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [R2.1](../software-requirements/srd004-helm-package-validation.md#r21), [R2.2](../software-requirements/srd004-helm-package-validation.md#r22), [R2.3](../software-requirements/srd004-helm-package-validation.md#r23) | S2 traces this requirement in the existing success criteria. |
| T4 | [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [AC2](../software-requirements/srd004-helm-package-validation.md#ac2) | S2 traces this requirement in the existing success criteria. |
| T5 | [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [R3.1](../software-requirements/srd004-helm-package-validation.md#r31), [R3.2](../software-requirements/srd004-helm-package-validation.md#r32), [R3.3](../software-requirements/srd004-helm-package-validation.md#r33) | S3 traces this requirement in the existing success criteria. |
| T6 | [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [AC3](../software-requirements/srd004-helm-package-validation.md#ac3) | S3 traces this requirement in the existing success criteria. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Invalid values fail before Kubernetes output is produced.

<!-- field: traces -->
#### Traces

- [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [R1.1](../software-requirements/srd004-helm-package-validation.md#r11)
- [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [R1.2](../software-requirements/srd004-helm-package-validation.md#r12)
- [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [R1.3](../software-requirements/srd004-helm-package-validation.md#r13)
- [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [AC1](../software-requirements/srd004-helm-package-validation.md#ac1)

<!-- sc: S2 -->
### S2

Package tampering is rejected before chart staging.

<!-- field: traces -->
#### Traces

- [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [R2.1](../software-requirements/srd004-helm-package-validation.md#r21)
- [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [R2.2](../software-requirements/srd004-helm-package-validation.md#r22)
- [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [R2.3](../software-requirements/srd004-helm-package-validation.md#r23)
- [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [AC2](../software-requirements/srd004-helm-package-validation.md#ac2)

<!-- sc: S3 -->
### S3

The archive includes every generated role asset and renders independently.

<!-- field: traces -->
#### Traces

- [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [R3.1](../software-requirements/srd004-helm-package-validation.md#r31)
- [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [R3.2](../software-requirements/srd004-helm-package-validation.md#r32)
- [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [R3.3](../software-requirements/srd004-helm-package-validation.md#r33)
- [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [AC3](../software-requirements/srd004-helm-package-validation.md#ac3)

<!-- field: out_of_scope -->
## Out of scope

- Live Kubernetes installation.
- kind cluster setup and teardown.
- Operator deployment and teardown guide.
