---
id: rel16.0-uc001-applier-family
title: Three applications compose one canonical applier family
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

Applications invoke one catalog-owned applier family through local wrappers.
The catalog owns reusable validation, apply, verification, and compensation
behavior while each application retains endpoint and rollout authority.

<!-- field: actor -->
## Actor

Catalog maintainers, application wrapper owners, and authorized application callers.

<!-- field: trigger -->
## Trigger

An application submits a decided values document to its local applier wrapper.

<!-- field: preconditions -->
## Preconditions

- The catalog applier profile closes against agent-core shared tools.
- Each application manifest declares the canonical applier and one local wrapper.
- Each wrapper owns its REST and rollout declarations.

<!-- field: flow -->
## Flow

- F1: The catalog boots the chart-independent applier profile and validates the host, apply, and rollout-read machines.
- F2: An application wrapper hosts its local endpoint and binds a decided values document to the canonical apply machine.
- F3: The machine writes and validates the document before invoking the application-owned apply command.
- F4: The application-owned verify word checks every local rollout target.
- F5: A failed verification invokes the application-owned rollback and reaches RolledBack only when compensation succeeds.
- F6: Packaging records canonical and local sources, and audits reject a copied promoted asset.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd022-applier R1 | Defines the canonical lifecycle and request-machine sequence. |
| T2 | srd022-applier R2 | Defines the shared applier vocabulary. |
| T3 | srd022-applier R3 | Keeps deployment authority in each application wrapper. |
| T4 | srd022-applier R4 | Defines validation, verification, and compensation behavior. |
| T5 | srd022-applier R5 | Defines portability, release, and ownership evidence. |
| T6 | srd022-applier R5 | Defines canonical-reference and composition-wrapper consistency, which the applications corpus states as srd003-application-consistency R5 and this SRD carries as a platform dependency. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The catalog applier profile boots independently with canonical host, apply, and rollout-read machines.

<!-- sc: S2 -->
### S2

Every application wrapper binds its local endpoint and decided values document to the canonical apply machine.

<!-- sc: S3 -->
### S3

Validation completes before the application-owned apply command runs.

<!-- sc: S4 -->
### S4

Verification checks every rollout target selected by the application.

<!-- sc: S5 -->
### S5

Verification failure invokes rollback and reports compensation success or failure explicitly.

<!-- sc: S6 -->
### S6

Packaging records canonical and local sources without copying a promoted asset.

<!-- field: out_of_scope -->
## Out of scope

- Owning any consuming application's chart, endpoint, deployment topology, credentials, RBAC, or operator workflow.
- Implementing Helm, kubectl, REST, lifecycle, or filesystem runtime behavior.
- Requiring one application to adopt another application's state-read extension.
