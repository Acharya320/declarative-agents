---
id: test-rel08.0-applier
title: Test Suite for Release 08.0 Declarative Applier Agent
release: "08.0"
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

This suite validates the applier agent (srd006): the standalone declarative
actuator for the coding chart. It drives a decided values patch through validate,
apply, and verify, asserting that a non-conforming patch is rejected by the
chart-schema dry-run, that a valid patch is applied as a values-file rollout and
verified via kubectl rollout status of the planner, executor, and critic
Deployments, that an apply or verify failure rolls back, and that the authority
boundary holds. A recording-fake tracer covers every declared terminal, real helm
validates the coding chart's values.schema.json, and a live kind target proves
apply and rollback on cluster. The coding application has no control plane, so no
creator, orchestrator, or panel hop appears in this suite.

<!-- field: traces -->
## Traces

- [rel08.0-uc001-declarative-applier](../use-cases/rel08.0-uc001-declarative-applier.md)

<!-- field: preconditions -->
## Preconditions

- agent-core is available for the request runtime, REST server, monitor, and exec tools.
- The applier profile is deployed with its apply, control, and monitor servers; it runs the agent-core image with helm and kubectl on a read-only /opt/tools from the pinned CLI donor.
- The coding chart carries a values.schema.json, and the apply surface is network-gated off the serving and collector pods.

<!-- field: test_cases -->
## Test cases

<!-- case: Applier validates, applies, and verifies a coding-chart values change | status: implemented -->
### Applier validates, applies, and verifies a coding-chart values change

Drives a decided values patch against the applier profile: the apply endpoint
receives the patch, the applier validates it against the coding chart's
values.schema.json by a helm dry-run, rejecting a non-conforming patch with a
client error and no rollout; for a valid patch it materializes the values as an
overrides file and applies helm upgrade --reuse-values, returning without
waiting; it verifies the planner, executor, and critic Deployments via kubectl
rollout status and reports the phase, rolling the release back on a verify stall.
Asserts the authority boundary: no invocation carries an endpoint, a
credential, or per-field --set construction.

helm and kubectl are recording fakes on the applier's PATH, taking their exit
codes from the scenario, so one run walks every terminal of the apply and
rollout machines. It is evidence about the machine, the arguments it
constructs, and the responses it maps, not that real helm and kubectl behave as
the declarations assume; the chart schema is covered by its own case below
against real helm, and the live kind target supplies the on-cluster evidence.

<!-- field: use_case -->
#### Use case

[rel08.0-uc001-declarative-applier](../use-cases/rel08.0-uc001-declarative-applier.md)

<!-- field: integration_target -->
#### Integration target

mage integration:applier

<!-- field: inputs | form: "tree" -->
#### Inputs

- valid_patch: adjust a coding role parameter through a values-plane patch
- invalid_patch: a patch violating the coding chart values schema
- agents:
  1. applier

<!-- field: expected | form: "tree" -->
#### Expected

- invalid_patch_rejected_by_schema_dry_run_no_rollout: true
- valid_patch_applied_as_values_file_reuse_values_no_wait: true
- no_per_field_set_construction: true
- rollout_verified_via_kubectl_rollout_status: true
- verify_stall_rolls_back_apply_failure_reports_failed: true
- patch_carries_no_endpoint_or_credential: true

<!-- field: traces -->
#### Traces

- [rel08.0-uc001-declarative-applier](../use-cases/rel08.0-uc001-declarative-applier.md) [S1](../use-cases/rel08.0-uc001-declarative-applier.md#s1)
- [rel08.0-uc001-declarative-applier](../use-cases/rel08.0-uc001-declarative-applier.md) [S2](../use-cases/rel08.0-uc001-declarative-applier.md#s2)
- [rel08.0-uc001-declarative-applier](../use-cases/rel08.0-uc001-declarative-applier.md) [S3](../use-cases/rel08.0-uc001-declarative-applier.md#s3)
- [rel08.0-uc001-declarative-applier](../use-cases/rel08.0-uc001-declarative-applier.md) [S4](../use-cases/rel08.0-uc001-declarative-applier.md#s4)
- [rel08.0-uc001-declarative-applier](../use-cases/rel08.0-uc001-declarative-applier.md) [S5](../use-cases/rel08.0-uc001-declarative-applier.md#s5)
- [srd006-applier](../software-requirements/srd006-applier.md) [AC1](../software-requirements/srd006-applier.md#ac1)
- [srd006-applier](../software-requirements/srd006-applier.md) [AC2](../software-requirements/srd006-applier.md#ac2)
- [srd006-applier](../software-requirements/srd006-applier.md) [AC3](../software-requirements/srd006-applier.md#ac3)
- [srd006-applier](../software-requirements/srd006-applier.md) [AC4](../software-requirements/srd006-applier.md#ac4)

<!-- case: The coding chart schema rejects a non-conforming values patch | status: implemented -->
### The coding chart schema rejects a non-conforming values patch

The applier's validate step is a helm dry-run against the coding chart's
values.schema.json, so the rejection is only as real as the schema. The case
above cannot prove it: its fake helm returns whatever exit code the scenario
sets, and a schema that accepted every document would pass it unchanged.

Real helm renders the real coding chart with two fixtures. The rejection must
name the offending path and the constraint it violated, not merely exit
non-zero -- a chart failing to render for an unrelated reason exits non-zero
too. A third assertion pins the constraint the fixture rests on, so loosening
the schema reads as the guard being removed rather than as a broken fixture.

<!-- field: use_case -->
#### Use case

[rel08.0-uc001-declarative-applier](../use-cases/rel08.0-uc001-declarative-applier.md)

<!-- field: integration_target -->
#### Integration target

mage test (coding chart values.schema.json against real helm)

<!-- field: inputs | form: "tree" -->
#### Inputs

- conforming: a values patch within the chart schema
- non_conforming: a values patch violating the chart schema

<!-- field: expected | form: "tree" -->
#### Expected

- conforming_patch_renders: true
- non_conforming_patch_rejected_by_path_and_constraint: true
- constrained_value_pinned_by_a_third_assertion: true

<!-- field: traces -->
#### Traces

- [rel08.0-uc001-declarative-applier](../use-cases/rel08.0-uc001-declarative-applier.md) [S2](../use-cases/rel08.0-uc001-declarative-applier.md#s2)
- [srd004-helm-package-validation](../software-requirements/srd004-helm-package-validation.md) [R1](../software-requirements/srd004-helm-package-validation.md#r1---strict-values)
- [srd006-applier](../software-requirements/srd006-applier.md) [AC2](../software-requirements/srd006-applier.md#ac2)

<!-- case: The applier holds the deployment CLIs and the serving roles do not | status: implemented -->
### The applier holds the deployment CLIs and the serving roles do not

Asserts over the rendered chart and the packaged role profiles that only the
applier carries helm and kubectl bound as exec tools and that the planner,
executor, and critic profiles carry no deployment CLI. Asserts the applier
Service targets the pods the applier Deployment creates and that its apply
surface is admitted only from the authorized caller, not from the serving or
collector pods, so route and NetworkPolicy agree on one authority path.

<!-- field: use_case -->
#### Use case

[rel08.0-uc001-declarative-applier](../use-cases/rel08.0-uc001-declarative-applier.md)

<!-- field: integration_target -->
#### Integration target

mage test (rendered chart and role profiles)

<!-- field: inputs | form: "tree" -->
#### Inputs

- rendered_chart: staged chart with the applier enabled and networkPolicy enabled

<!-- field: expected | form: "tree" -->
#### Expected

- applier_alone_binds_helm_and_kubectl: true
- serving_roles_carry_no_deployment_cli: true
- apply_surface_admits_only_the_authorized_caller: true
- apply_surface_not_reachable_from_serving_or_collector_pods: true

<!-- field: traces -->
#### Traces

- [rel08.0-uc001-declarative-applier](../use-cases/rel08.0-uc001-declarative-applier.md) [S6](../use-cases/rel08.0-uc001-declarative-applier.md#s6)
- [srd006-applier](../software-requirements/srd006-applier.md) [AC4](../software-requirements/srd006-applier.md#ac4)
- [srd006-applier](../software-requirements/srd006-applier.md) [AC5](../software-requirements/srd006-applier.md#ac5)

<!-- case: A broken rollout read is distinguishable from a progressing one | status: implemented -->
### A broken rollout read is distinguishable from a progressing one

Asserts over the declared profile that the rollout read reaches a broken-read
terminal rather than folding it into progressing. kubectl rollout status exits
non-zero for a rollout that is still progressing and for one it could not read
at all, so a phase taken from that exit code alone would report an absent
Deployment, a wrong kubeconfig, a denied RBAC read, or an unreachable API
server as ongoing progress. A read the cluster cannot serve returns a gateway
error rather than a phase.

<!-- field: use_case -->
#### Use case

[rel08.0-uc001-declarative-applier](../use-cases/rel08.0-uc001-declarative-applier.md)

<!-- field: integration_target -->
#### Integration target

mage test (applier rollout machine)

<!-- field: inputs | form: "tree" -->
#### Inputs

- deployment_shapes: rolled out, no replica ready yet, and unreadable

<!-- field: expected | form: "tree" -->
#### Expected

- unreadable_rollout_returns_gateway_error_not_a_phase: true
- complete_and_progressing_remain_distinguishable: true

<!-- field: traces -->
#### Traces

- [rel08.0-uc001-declarative-applier](../use-cases/rel08.0-uc001-declarative-applier.md) [S7](../use-cases/rel08.0-uc001-declarative-applier.md#s7)
- [srd006-applier](../software-requirements/srd006-applier.md) [AC3](../software-requirements/srd006-applier.md#ac3)

<!-- case: A live on-cluster rollout and rollback through the applier | status: implemented -->
### A live on-cluster rollout and rollback through the applier

What the fake CLIs cannot stand in for. The applier runs the agent-core
runtime with helm and kubectl from the pinned CLI donor, loaded into kind,
and the coding chart installed with the applier enabled; inside the pod, helm
reports the donor's version and a write to /opt/tools fails; a values patch then goes through the apply endpoint
and the release moves to a new revision, which is the evidence a 200 is not --
the fake tier returns one of those without helm existing.

The compensating leg is live too. A reserved patch forces a deterministic
post-upgrade verification failure; the applier's declared kubectl rollout
status times out, reaches the rollback terminal, and invokes real helm
rollback, which restores the prior computed release values and ready planner,
executor, and critic Deployments. The same non-conforming document the schema
case validates locally is rejected here against the real coding chart schema,
leaving the release untouched. Bounded host and in-pod probes classify a
Kubernetes API or pod-network loss apart from the semantic proof.

<!-- field: use_case -->
#### Use case

[rel08.0-uc001-declarative-applier](../use-cases/rel08.0-uc001-declarative-applier.md)

<!-- field: integration_target -->
#### Integration target

mage integration:applierLive

<!-- field: inputs | form: "tree" -->
#### Inputs

- cluster: kind with the agent-core runtime and the pinned CLI donor loaded
- patch: a conforming coding values patch
- rollback_trigger: a patch that forces a post-upgrade verification failure
- rejected: a non-conforming coding values patch

<!-- field: expected | form: "tree" -->
#### Expected

- real_helm_upgrade_changes_the_release: true
- post_upgrade_verify_failure_reaches_rolled_back: true
- real_helm_rollback_restores_prior_computed_values: true
- coding_role_deployments_recover_after_rollback: true
- non_conforming_patch_rejected_on_the_cluster: true

<!-- field: traces -->
#### Traces

- [rel08.0-uc001-declarative-applier](../use-cases/rel08.0-uc001-declarative-applier.md) [S3](../use-cases/rel08.0-uc001-declarative-applier.md#s3)
- [rel08.0-uc001-declarative-applier](../use-cases/rel08.0-uc001-declarative-applier.md) [S4](../use-cases/rel08.0-uc001-declarative-applier.md#s4)
- [srd005-live-helm-smoke](../software-requirements/srd005-live-helm-smoke.md) [R3](../software-requirements/srd005-live-helm-smoke.md#r3---end-to-end-evidence)
- [srd006-applier](../software-requirements/srd006-applier.md) [AC3](../software-requirements/srd006-applier.md#ac3)
