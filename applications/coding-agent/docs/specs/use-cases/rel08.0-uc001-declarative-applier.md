---
id: rel08.0-uc001-declarative-applier
title: The Applier Agent Validates, Applies, and Verifies a Coding-Chart Values Change
status: implemented
test_suite: test-rel08.0-applier
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An authorized operator or CI caller submits a decided, values-plane patch to the
applier (srd006) -- the standalone declarative actuator for the coding chart
(srd003). The applier validates the patch against the chart's values.schema.json
(srd004) with a helm dry-run, rejecting a non-conforming patch without a rollout.
For a valid patch it materializes the values as an overrides file and applies helm
upgrade --reuse-values, returning without waiting, re-rendering the prepared role
topology, then verifies the planner, executor, and critic Deployments via kubectl
rollout status and reports the phase, rolling the release back on a verify stall
and reporting an apply-command failure as failed. Throughout, the applier edits
deployment values and triggers rollouts
only; the patch carries no host, URL, method, or credential for a running agent,
the apply surface is reached only by the operator or CI caller, and authority stays
set by declared, validated config (agent-core srd028 V30). No agent-core change is
required: the actuation is a state machine binding the helm and kubectl CLIs in the
applier's container.

The coding application runs no control plane, so there is no in-cluster creator,
orchestrator, or provisioning panel in this flow. The applier is a standalone
deployment-plane agent whose apply surface is network-gated off the planner,
executor, critic, and collector pods and reachable only by the authorized caller.

<!-- field: actor -->
## Actor

An authorized operator or CI process actuating the coding chart through the applier's apply surface.

<!-- field: trigger -->
## Trigger

An operator or CI pipeline decides a coding-chart values change (adjust a role
replica count, switch the model endpoint, reconfigure a parameter) and submits it
to the applier's apply endpoint.

<!-- field: preconditions -->
## Preconditions

- The applier profile is deployed with its apply, control, and monitor servers, and its apply surface is network-gated off the serving and collector pods.
- The applier runs the agent-core image with the helm and kubectl CLIs delivered by the pinned CLI donor, and its ServiceAccount can upgrade the release and read the coding role Deployments.
- The coding chart carries a values.schema.json encoding the deployment values constraints (srd004).
- The apply surface is reachable only by the authorized operator or CI caller.

<!-- field: flow -->
## Flow

- F1: The operator or CI caller posts the decided values patch to the applier's apply endpoint; it submits no host, URL, method, or credential for a running agent.
- F2: The applier validates the patch against the chart's values.schema.json with a helm dry-run and rejects a non-conforming patch with a client error and no rollout.
- F3: The applier materializes the validated patch as an overrides file and applies helm upgrade --reuse-values, returning without waiting, re-rendering the prepared role topology, rather than constructing per-field --set arguments.
- F4: The applier verifies the planner, executor, and critic Deployments via kubectl rollout status and reports the phase, rolling the release back on a verify stall and reporting an apply-command failure as failed.
- F5: The applier reports the terminal phase to the caller; a subsequent coding request is served by the applied topology.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd006-applier AC1 | The applier exposes the request-scoped apply endpoint with control and monitor servers and no panel state view. |
| T2 | srd006-applier AC2 | The applier validates the patch against the chart schema by a dry-run and applies it as a values-file rollout, values-plane only. |
| T3 | srd006-applier AC3 | The applier verifies the role Deployments' rollout, rolls back on a verify stall, and reports an apply-command failure as failed. |
| T4 | srd006-applier AC4 | The applier alone holds helm and kubectl, the apply surface is network-gated off the serving pods, and authority stays set by declared validated config. |
| T5 | srd006-applier AC5 | Deciding the values and owning lifecycle are non-goals; the applier receives helm and kubectl from the pinned CLI donor, so no image diverges. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The applier exposes the apply surface bound as a request-scoped machine_request with control and monitor servers, serves no panel state view, and depends on no in-cluster caller.

<!-- field: traces -->
#### Traces

- srd006-applier R1.1
- srd006-applier R1.4
- srd006-applier AC1

<!-- sc: S2 -->
### S2

A non-conforming values patch is rejected by the chart-schema dry-run with a client error and no rollout.

<!-- field: traces -->
#### Traces

- srd006-applier R2.1
- srd004-helm-package-validation R1.1
- srd004-helm-package-validation R1.2
- srd004-helm-package-validation R1.3
- srd006-applier AC2

<!-- sc: S3 -->
### S3

A valid patch is materialized as an overrides file and applied with helm upgrade --reuse-values, returning without waiting, re-rendering the prepared role topology, with no per-field --set construction.

<!-- field: traces -->
#### Traces

- srd006-applier R2.2
- srd003-helm-topology R1.1
- srd003-helm-topology R1.2
- srd003-helm-topology R1.3
- srd003-helm-topology R1.4
- srd003-helm-topology R1.5
- srd006-applier AC2

<!-- sc: S4 -->
### S4

The planner, executor, and critic Deployments' rollout is verified via kubectl rollout status and its phase reported, and a verify stall rolls the release back and restores the prior computed release while an apply-command failure is reported as failed, proven on the kind rig.

<!-- field: traces -->
#### Traces

- srd006-applier R3.1
- srd006-applier R3.2
- srd005-live-helm-smoke R3.1
- srd005-live-helm-smoke R3.2
- srd005-live-helm-smoke R3.3
- srd006-applier AC3

<!-- sc: S5 -->
### S5

The patch submits no endpoint or credential to a running agent, the apply surface is reached only by the operator or CI caller and network-gated off the serving pods, and authority stays set by declared validated config.

<!-- field: traces -->
#### Traces

- srd006-applier R2.3
- srd006-applier R4.2
- srd006-applier R4.3
- srd006-applier AC4

<!-- sc: S6 -->
### S6

The applier alone holds helm and kubectl bound as exec tools, the serving roles carry none, and the applier receives the CLIs from the pinned CLI donor on a read-only /opt/tools, so no image diverges from the binary-only image contract.

<!-- field: traces -->
#### Traces

- srd006-applier R4.1
- srd006-applier R5.3
- srd006-applier AC4
- srd006-applier AC5

<!-- sc: S7 -->
### S7

A rollout status the cluster cannot serve -- an absent Deployment, a wrong kubeconfig, a denied RBAC read, an unreachable API server -- returns a gateway error rather than a phase, so a broken read is never reported as a progressing rollout.

<!-- field: traces -->
#### Traces

- srd006-applier R3.3
- srd006-applier AC3

<!-- field: out_of_scope -->
## Out of scope

- The decision of what to change, which the operator or CI caller owns (srd006 R5.1).
- Role profile selection and packaging, which are the closure and chart's (srd002, srd003).
- The coding request path itself (plan, delegate, gate), which is srd001 and the canonical library agents.
