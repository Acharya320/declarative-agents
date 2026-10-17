---
id: rel18.0-uc001-host-side-deploy
title: An operator runs mage deploy and the applier's deploy machine installs the release
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An operator brings an application up on a running cluster with one target.
The target renders the applier's deploy declarations with this
application's chart, release, namespace, and kubeconfig, seeds the deploy
machine with the values document, and gates on the terminal state. The same
machine that owns day-2 rollouts inside the cluster owns the first install
from the developer's machine, so we keep one deployment mechanism rather
than one per application and one per environment.

<!-- field: actor -->
## Actor

A developer or release operator bringing up a demo or a cloud environment.

<!-- field: trigger -->
## Trigger

The operator runs mage deploy, or a demo target that delegates to it, against a running cluster.

<!-- field: preconditions -->
## Preconditions

- The target cluster is running, and the target confirms it before any Helm word executes.
- The catalog applier deploy and undeploy profiles close against agent-core shared tools.
- The local Helm major version matches the spellings the rendered words carry.
- No model provider is needed; the deploy machine declares no model word.

<!-- field: flow -->
## Flow

- F1: The operator runs mage deploy with no arguments, and the target reads its coordinates from checked-in configuration rather than from process environment.
- F2: The target confirms the cluster is running and resolves a private kubeconfig for it.
- F3: The target packages the chart and resolves the release, namespace, values overlay, and timeout, failing when any coordinate is empty.
- F4: The target renders the deploy declarations into build, replacing every placeholder coordinate, and fails when a placeholder survives.
- F5: The target writes the values document as a request file and runs the deploy profile with the rendered declarations.
- F6: The machine writes the overrides, validates them with a dry run, applies them with install semantics, and verifies every Deployment in the release's namespace.
- F7: A validation failure reaches Rejected before any mutation, and a verify stall compensates by rollback.
- F8: The target reports the terminal state and exits non-zero when that state is a failure, so the build fails rather than continuing against a release that did not come up.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd022-applier](../software-requirements/srd022-applier.md) [R1](../software-requirements/srd022-applier.md#r1---canonical-lifecycle-and-request-machines) | The deploy machine reuses the canonical request-machine sequencing this requirement fixes. |
| T2 | [srd022-applier](../software-requirements/srd022-applier.md) [R3](../software-requirements/srd022-applier.md#r3---application-owned-boundary) | Chart, release, namespace, and rollout coordinates stay application-owned, now as rendered declarations. |
| T3 | [srd022-applier](../software-requirements/srd022-applier.md) [R4](../software-requirements/srd022-applier.md#r4---validation-and-compensation) | Validation precedes mutation and a verification failure reaches compensation. |
| T4 | [srd022-applier](../software-requirements/srd022-applier.md) [R6](../software-requirements/srd022-applier.md#r6---host-side-deploy-and-undeploy-machines) | Defines the deploy and undeploy machines, the atomic contract, and the rendering rules this flow exercises. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The deploy and undeploy tool selections hold no model word, so a deploy consumes no tokens.

<!-- sc: S2 -->
### S2

Every rendered deploy argv carries a resolved chart path, release, namespace, and kubeconfig, and no placeholder token survives rendering.

<!-- sc: S3 -->
### S3

No rendered deploy argv carries --reuse-values, so a deploy restates its values from the chart, the overlay, and the request document.

<!-- sc: S4 -->
### S4

A first install that fails leaves no release behind, because the apply word carries --atomic.

<!-- sc: S5 -->
### S5

A verify stall after an upgrade reaches Compensated through rollback, and one after a first install reaches Failed with the release left in place for inspection.

<!-- sc: S6 -->
### S6

mage deploy exits non-zero on a failure terminal, unlike mage diagnose, which reports and never gates.

<!-- sc: S7 -->
### S7

mage undeploy removes a present release and reports Absent for a release that is already gone, exiting zero in both cases.

<!-- field: out_of_scope -->
## Out of scope

- Cluster provisioning; the deploy machine owns the Helm step alone, and kindrig and gcprig keep the bootstrap layer.
- The chatbot-mesh and agent-architecture demo migrations; they adopt the target after coding-agent proves it.
- GKE, which adopts the same target once the cloud bootstrap exists.
- Retiring the helm install path in the chatbot-mesh integration smoke, which carries release-budget and telemetry concerns of its own.
