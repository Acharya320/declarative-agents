---
id: rel04.0-uc001-declarative-applier
title: The Applier Agent Validates, Applies, and Verifies an Agent Architecture Values Change
status: planned
test_suite: test-rel04.0-declarative-applier
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An authorized operator or CI caller submits a decided, values-plane patch to
the applier (srd002) -- the standalone declarative actuator for the Agent
Architecture chart (srd001). The applier validates the patch against the
chart's values.schema.json with a helm dry-run, rejecting a non-conforming
patch without a rollout. For a valid patch it materializes the values as an
overrides file and applies helm upgrade --reuse-values, returning without
waiting, re-rendering the prepared topology, then verifies the curator and
collector Deployments via kubectl rollout status and reports the phase,
rolling the release back on a verify stall and reporting an apply-command
failure as failed. Throughout, the applier edits deployment values and triggers
rollouts only; the patch carries no host, URL, method, or credential for a
running agent, the apply surface is reached only by the operator or CI caller,
and authority stays set by declared, validated config (agent-core srd028 V30).
No agent-core change is required: the actuation is a state machine binding the
helm and kubectl CLIs in the applier's container.

The Agent Architecture runs no control plane, so there is no in-cluster
creator, orchestrator, or provisioning panel in this flow. The applier is a
standalone deployment-plane agent whose apply surface is network-gated off the
curator and collector pods and reachable only by the authorized caller.

<!-- field: actor -->
## Actor

An authorized operator or CI process actuating the Agent Architecture chart through the applier's apply surface.

<!-- field: trigger -->
## Trigger

An operator or CI pipeline decides an Agent Architecture values change (adjust
a replica count, switch the model or telemetry endpoint, reconfigure a
parameter) and submits it to the applier's apply endpoint.

<!-- field: preconditions -->
## Preconditions

- The applier profile is deployed with its apply, control, and monitor servers, and its apply surface is network-gated off the curator and collector pods.
- The applier runs the agent-core image with the helm and kubectl CLIs delivered by the pinned CLI donor, and its ServiceAccount can upgrade the release and read the curator and collector Deployments.
- The chart carries a values.schema.json encoding the deployment values constraints (srd001 R2).
- The apply surface is reachable only by the authorized operator or CI caller.

<!-- field: flow -->
## Flow

- F1: The operator or CI caller posts the decided values patch to the applier's apply endpoint; it submits no host, URL, method, or credential for a running agent.
- F2: The applier validates the patch against the chart's values.schema.json with a helm dry-run and rejects a non-conforming patch with a client error and no rollout.
- F3: The applier materializes the validated patch as an overrides file and applies helm upgrade --reuse-values, returning without waiting, re-rendering the prepared topology, rather than constructing per-field --set arguments.
- F4: The applier verifies the curator and collector Deployments via kubectl rollout status and reports the phase, rolling the release back on a verify stall and reporting an apply-command failure as failed.
- F5: The applier reports the terminal phase to the caller; a subsequent guided run is served by the applied topology.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd002-applier](../software-requirements/srd002-applier.md) [AC1](../software-requirements/srd002-applier.md#ac1) | The applier exposes the request-scoped apply endpoint with control and monitor servers and no panel state view. |
| T2 | [srd002-applier](../software-requirements/srd002-applier.md) [AC2](../software-requirements/srd002-applier.md#ac2) | The applier validates the patch against the chart schema by a dry-run and applies it as a values-file rollout, values-plane only. |
| T3 | [srd002-applier](../software-requirements/srd002-applier.md) [AC3](../software-requirements/srd002-applier.md#ac3) | The applier verifies the workload rollout, rolls back on a verify stall, and reports an apply-command failure as failed. |
| T4 | [srd002-applier](../software-requirements/srd002-applier.md) [AC4](../software-requirements/srd002-applier.md#ac4) | The applier alone holds helm and kubectl, the apply surface is network-gated off the serving pods, and authority stays set by declared validated config. |
| T5 | [srd002-applier](../software-requirements/srd002-applier.md) [AC5](../software-requirements/srd002-applier.md#ac5) | Deciding the values and owning lifecycle are non-goals; the applier receives helm and kubectl from the pinned CLI donor, so no image diverges. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The applier exposes the apply surface bound as a request-scoped machine_request with control and monitor servers, serves no panel state view, and depends on no in-cluster caller.

<!-- field: traces -->
#### Traces

- [srd002-applier](../software-requirements/srd002-applier.md) [R1.1](../software-requirements/srd002-applier.md#r11)
- [srd002-applier](../software-requirements/srd002-applier.md) [R1.4](../software-requirements/srd002-applier.md#r14)
- [srd002-applier](../software-requirements/srd002-applier.md) [AC1](../software-requirements/srd002-applier.md#ac1)

<!-- sc: S2 -->
### S2

A non-conforming values patch is rejected by the chart-schema dry-run with a client error and no rollout.

<!-- field: traces -->
#### Traces

- [srd002-applier](../software-requirements/srd002-applier.md) [R2.1](../software-requirements/srd002-applier.md#r21)
- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [R2.3](../software-requirements/srd001-helm-deployment.md#r23)
- [srd002-applier](../software-requirements/srd002-applier.md) [AC2](../software-requirements/srd002-applier.md#ac2)

<!-- sc: S3 -->
### S3

A valid patch is materialized as an overrides file and applied with helm upgrade --reuse-values, returning without waiting, re-rendering the prepared topology, with no per-field --set construction.

<!-- field: traces -->
#### Traces

- [srd002-applier](../software-requirements/srd002-applier.md) [R2.2](../software-requirements/srd002-applier.md#r22)
- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [R1](../software-requirements/srd001-helm-deployment.md#r1---deployed-topology)
- [srd002-applier](../software-requirements/srd002-applier.md) [AC2](../software-requirements/srd002-applier.md#ac2)

<!-- sc: S4 -->
### S4

The curator and collector Deployments' rollout is verified via kubectl rollout status and its phase reported, and a verify stall rolls the release back and restores the prior computed release while an apply-command failure is reported as failed, proven on the kind rig.

<!-- field: traces -->
#### Traces

- [srd002-applier](../software-requirements/srd002-applier.md) [R3.1](../software-requirements/srd002-applier.md#r31)
- [srd002-applier](../software-requirements/srd002-applier.md) [R3.2](../software-requirements/srd002-applier.md#r32)
- [srd001-helm-deployment](../software-requirements/srd001-helm-deployment.md) [R3](../software-requirements/srd001-helm-deployment.md#r3---live-kind-smoke)
- [srd002-applier](../software-requirements/srd002-applier.md) [AC3](../software-requirements/srd002-applier.md#ac3)

<!-- sc: S5 -->
### S5

The patch submits no endpoint or credential to a running agent, the apply surface is reached only by the operator or CI caller and network-gated off the serving pods, and authority stays set by declared validated config.

<!-- field: traces -->
#### Traces

- [srd002-applier](../software-requirements/srd002-applier.md) [R2.3](../software-requirements/srd002-applier.md#r23)
- [srd002-applier](../software-requirements/srd002-applier.md) [R4.2](../software-requirements/srd002-applier.md#r42)
- [srd002-applier](../software-requirements/srd002-applier.md) [R4.3](../software-requirements/srd002-applier.md#r43)
- [srd002-applier](../software-requirements/srd002-applier.md) [AC4](../software-requirements/srd002-applier.md#ac4)

<!-- sc: S6 -->
### S6

The applier alone holds helm and kubectl bound as exec tools, the curator and collector carry none, and the applier receives the CLIs from the pinned CLI donor on a read-only /opt/tools, so no image diverges from the binary-only image contract.

<!-- field: traces -->
#### Traces

- [srd002-applier](../software-requirements/srd002-applier.md) [R4.1](../software-requirements/srd002-applier.md#r41)
- [srd002-applier](../software-requirements/srd002-applier.md) [R5.3](../software-requirements/srd002-applier.md#r53)
- [srd002-applier](../software-requirements/srd002-applier.md) [AC4](../software-requirements/srd002-applier.md#ac4)
- [srd002-applier](../software-requirements/srd002-applier.md) [AC5](../software-requirements/srd002-applier.md#ac5)

<!-- sc: S7 -->
### S7

A rollout status the cluster cannot serve -- an absent Deployment, a wrong kubeconfig, a denied RBAC read, an unreachable API server -- returns a gateway error rather than a phase, so a broken read is never reported as a progressing rollout.

<!-- field: traces -->
#### Traces

- [srd002-applier](../software-requirements/srd002-applier.md) [R3.3](../software-requirements/srd002-applier.md#r33)
- [srd002-applier](../software-requirements/srd002-applier.md) [AC3](../software-requirements/srd002-applier.md#ac3)

<!-- field: out_of_scope -->
## Out of scope

- The decision of what to change, which the operator or CI caller owns (srd002 R5.1).
- Workload profile selection and packaging, which are the closure and chart's (srd001).
- The guided demo request path itself, which is Release 00.0 and the canonical curator.
