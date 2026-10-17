---
id: rel17.0-uc001-rig-doctor
title: An operator runs mage diagnose after a failed scenario and reads a cited diagnosis
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

After a kind scenario fails, an operator runs the application's diagnose
target. The target captures a snapshot of the scenario namespace and the
trace spool tail into an evidence directory, then runs the catalog
rig-doctor over it. The operator reads a diagnosis whose every claim cites
a captured file and line, and runs the suggested read-only command.

<!-- field: actor -->
## Actor

A developer or release operator working on the kind rig.

<!-- field: trigger -->
## Trigger

An integration or demo target fails, or a rollout stays unready.

<!-- field: preconditions -->
## Preconditions

- The da-platform cluster or the scenario's own cluster is still running.
- The catalog rig-doctor profile closes against agent-core shared tools.
- A model provider is reachable for the diagnosis loop.

<!-- field: flow -->
## Flow

- F1: The operator runs mage diagnose with the scenario name as an explicit argument.
- F2: The target resolves the cluster, the da-<scenario> namespace, and a private kubeconfig without reading environment variables.
- F3: Capture writes events, describe, pods, rollout counts, pod logs, cluster events, and the trace spool tail with a manifest listing each file and each capture error.
- F4: The target runs the rig-doctor profile with the evidence directory as its working directory.
- F5: The rig-doctor reads the manifest, correlates events, rollout counts, logs, and traces, and writes diagnosis.yaml.
- F6: The target prints the evidence directory and the diagnosis path and exits zero whatever the diagnosis outcome.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd023-rig-doctor](../software-requirements/srd023-rig-doctor.md) [R1](../software-requirements/srd023-rig-doctor.md#r1---evidence-directory-contract) | Defines the evidence directory the capture writes and the agent reads. |
| T2 | [srd023-rig-doctor](../software-requirements/srd023-rig-doctor.md) [R2](../software-requirements/srd023-rig-doctor.md#r2---read-only-diagnosis-loop) | Keeps the diagnosis loop read-only by declaration. |
| T3 | [srd023-rig-doctor](../software-requirements/srd023-rig-doctor.md) [R3](../software-requirements/srd023-rig-doctor.md#r3---diagnosis-output) | Defines the cited diagnosis the operator reads. |
| T4 | [srd023-rig-doctor](../software-requirements/srd023-rig-doctor.md) [R4](../software-requirements/srd023-rig-doctor.md#r4---terminal-states) | Guarantees a terminal state on empty or invalid input. |
| T5 | [srd023-rig-doctor](../software-requirements/srd023-rig-doctor.md) [R5](../software-requirements/srd023-rig-doctor.md#r5---portability-and-registration) | Registers the family and keeps its tests offline. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Capture produces a manifest listing events, describe, pods, rollout, logs, cluster-events, and traces files, and records capture errors without aborting.

<!-- sc: S2 -->
### S2

The rig-doctor tool manifest holds no exec, kubectl, or helm word, and its only write targets diagnosis.yaml.

<!-- sc: S3 -->
### S3

Given a fixture with a failing rollout, the diagnosis names the failing Deployment and cites the event or log line that shows the cause.

<!-- sc: S4 -->
### S4

An empty evidence directory reaches NoEvidence and writes no diagnosis.

<!-- sc: S5 -->
### S5

mage diagnose exits zero and prints the evidence and diagnosis paths whether or not diagnosis succeeds.

<!-- sc: S6 -->
### S6

The root-cause-analyst binding appears in all three semantic models.

<!-- field: out_of_scope -->
## Out of scope

- Remediation or any cluster mutation.
- Live queries against the cluster or the host ingress during diagnosis.
- Diagnose targets in coding-agent and agent-architecture; they adopt the verb after chatbot-mesh.
