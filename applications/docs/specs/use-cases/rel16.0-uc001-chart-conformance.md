---
id: rel16.0-uc001-chart-conformance
title: Add a Workload and Learn at Audit Time That Its Image Floats
status: planned
test_suite: test-rel16.0-chart-conformance
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: release -->
## Release

16.0

<!-- field: summary -->
## Summary

An engineer adds a workload to an application chart and gives its image a
convenient tag. The repository audit renders every chart under every
checked-in values overlay, checks the rendered manifests against the ENG01
chart rules, and fails naming the chart, the overlay, the resource, the
image, and the rule. The engineer pins the image and the audit passes. A
violation that must stand goes in the exception baseline with the issue
that will remove it, and the baseline shrinks from there.

<!-- field: actor -->
## Actor

An engineer changing an application chart or its values.

<!-- field: trigger -->
## Trigger

A chart gains a workload, an image reference changes, a values overlay is
added or edited, or a kind configuration is written.

<!-- field: preconditions -->
## Preconditions

- ENG01 carries the numbered chart conformance rules and srd005 owns them.
- Helm is on the host; the audit states which rules did not run when it is absent.
- Each application chart carries its values overlays under helm/ci/.

<!-- field: flow -->
## Flow

- F1: The engineer adds a workload to an application chart, or edits an
  image reference or a values overlay, and runs the repository audit.
  - traces:
    - srd005-chart-conformance R1
- F2: The audit vendors the shared library chart into each application
  chart, then renders every chart once under its own defaults and once
  per checked-in values overlay.
  - traces:
    - srd005-chart-conformance R1.2
    - srd005-chart-conformance R2.2
- F3: Each rendered manifest is checked against the rules: tags do not
  float, third-party images carry digests, image references name
  registry, role, and typed identity, pull policies keep a loaded
  node off the registry, no Service asks kind for a load balancer, and
  every selected workload declares readiness, and an applier that
  upgrades can read and write every kind its chart renders.
  - traces:
    - srd005-chart-conformance R1.1
    - srd005-chart-conformance R2.1
    - srd005-chart-conformance R3.1
    - srd005-chart-conformance R4.1
    - srd005-chart-conformance R5.1
    - srd005-chart-conformance R10.1
    - srd005-chart-conformance R10.2
    - srd005-chart-conformance R10.3
    - srd005-chart-conformance R11.1
- F4: Every checked-in kind configuration is checked for a node image pinned
  by version and digest.
  - traces:
    - srd005-chart-conformance R6.1
- F5: A finding absent from the exception baseline fails the audit, naming
  the chart, the overlay, the resource, the offending value, and the
  rule. The engineer reads it without re-rendering anything.
  - traces:
    - srd005-chart-conformance R7.2
    - srd005-chart-conformance R8.1
- F6: A violation that must stand is added to the exception baseline with
  the issue that will remove it and the reason it stands, and the audit
  passes.
  - traces:
    - srd005-chart-conformance R7.1
- F7: When the issue lands and the violation stops reproducing, the audit
  fails on the stale baseline entry until it is removed, so the file
  records only exceptions that still apply.
  - traces:
    - srd005-chart-conformance R7.3

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | docs/engineering/eng01-kind-test-demo-rig.md | States the numbered chart conformance rules the gate cites. |
| T2 | applications/<application>/helm/ | The chart and the values overlays rendered and checked. |
| T3 | applications/charts/agent-services/ | Shared library chart vendored into each application chart before rendering. |
| T4 | magefiles/audit.go | Repository audit that dispatches the conformance check. |
| T5 | magefiles/chartconf/baseline.yaml | Checked-in exception baseline the gate reads and never writes. |
| T6 | [srd005-chart-conformance](../software-requirements/srd005-chart-conformance.md) [AC1](../software-requirements/srd005-chart-conformance.md#ac1) | Rendering coverage across charts and overlays. |
| T7 | [srd005-chart-conformance](../software-requirements/srd005-chart-conformance.md) [AC2](../software-requirements/srd005-chart-conformance.md#ac2) | Digest and commit-tag classification. |
| T8 | [srd005-chart-conformance](../software-requirements/srd005-chart-conformance.md) [AC3](../software-requirements/srd005-chart-conformance.md#ac3) | Kind node image pinning. |
| T9 | [srd005-chart-conformance](../software-requirements/srd005-chart-conformance.md) [AC4](../software-requirements/srd005-chart-conformance.md#ac4) | Baseline suspension, unlisted failure, and stale-entry failure. |
| T10 | [srd005-chart-conformance](../software-requirements/srd005-chart-conformance.md) [AC5](../software-requirements/srd005-chart-conformance.md#ac5) | Finding names chart, overlay, resource, value, and rule. |
| T11 | [srd005-chart-conformance](../software-requirements/srd005-chart-conformance.md) [AC7](../software-requirements/srd005-chart-conformance.md#ac7) | Image-reference grammar and prohibited aliases. |
| T13 | [srd005-chart-conformance](../software-requirements/srd005-chart-conformance.md) [AC8](../software-requirements/srd005-chart-conformance.md#ac8) | Applier Role coverage of every rendered kind. |
| T12 | [test-rel16.0-chart-conformance](../test-suites/test-rel16.0-chart-conformance.md) | Structural and executable evidence for this flow. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Every application chart and every checked-in values overlay is
rendered and checked; no chart is covered by its defaults alone.

<!-- field: traces -->
#### Traces

- [srd005-chart-conformance](../software-requirements/srd005-chart-conformance.md) [AC1](../software-requirements/srd005-chart-conformance.md#ac1)

<!-- sc: S2 -->
### S2

No image in any rendered manifest floats, every third-party image
carries a digest, every repository-built image carries a
non-floating tag, and every image matches the reference grammar.

<!-- field: traces -->
#### Traces

- [srd005-chart-conformance](../software-requirements/srd005-chart-conformance.md) [AC1](../software-requirements/srd005-chart-conformance.md#ac1)
- [srd005-chart-conformance](../software-requirements/srd005-chart-conformance.md) [AC2](../software-requirements/srd005-chart-conformance.md#ac2)
- [srd005-chart-conformance](../software-requirements/srd005-chart-conformance.md) [AC7](../software-requirements/srd005-chart-conformance.md#ac7)

<!-- sc: S3 -->
### S3

Every checked-in kind configuration pins its node image by version and
digest.

<!-- field: traces -->
#### Traces

- [srd005-chart-conformance](../software-requirements/srd005-chart-conformance.md) [AC3](../software-requirements/srd005-chart-conformance.md#ac3)

<!-- sc: S4 -->
### S4

The exception baseline shrinks: an unlisted violation fails, and a
listed one that no longer reproduces fails until its entry is removed.

<!-- field: traces -->
#### Traces

- [srd005-chart-conformance](../software-requirements/srd005-chart-conformance.md) [AC4](../software-requirements/srd005-chart-conformance.md#ac4)

<!-- sc: S5 -->
### S5

A failure names the chart, the overlay, the resource, the value, and
the rule, and an engineer acts on it without re-rendering.

<!-- field: traces -->
#### Traces

- [srd005-chart-conformance](../software-requirements/srd005-chart-conformance.md) [AC5](../software-requirements/srd005-chart-conformance.md#ac5)

<!-- sc: S6 -->
### S6

An applier that may write can read and write every object kind its
chart renders, so an in-cluster upgrade is never refused by its own
Role.

<!-- field: traces -->
#### Traces

- [srd005-chart-conformance](../software-requirements/srd005-chart-conformance.md) [AC8](../software-requirements/srd005-chart-conformance.md#ac8)

<!-- field: out_of_scope -->
## Out of scope

- Deciding whether an application template should call the shared library chart instead of rendering its own workload.
- Refreshing a pinned digest when an upstream release appears.
- Running any chart against a live cluster; the gate reads manifests only.
- Charts owned by consumer repositories.
