---
id: srd022-applier
title: Canonical Deployment Applier Family
depends_on:
  - srd_id: srd001-agent-functional-blocks
    symbols_used: [R1, R2, R3]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: context -->
## Context

The applier family is the reusable Executor / Actuation Agent and Change
Manager realization shared by chatbot-mesh, coding-agent, and
agent-architecture. The catalog owns lifecycle hosting, request-machine
sequencing, name-only selections, and shared ToolDefs. Each application owns
its wrapper, REST authority, chart and schema, release and namespace,
Deployment targets, Helm and kubectl commands, and rollout policy.

<!-- field: problem -->
## Problem

Three applications copied the same lifecycle machine, request machines,
selections, and ToolDefs. Their application-specific deployment coordinates
differ, so copying obscured the stable family contract while promoting an
entire application implementation would transfer excessive authority.

<!-- field: goals -->
## Goals

- G1: Establish one reusable applier machine and ToolDef closure.
- G2: Preserve application ownership of deployment coordinates and authority.
- G3: Preserve validation, explicit verification, compensation, and least authority.
- G4: Keep the family portable, independently bootable, and versioned through root v0 releases.

<!-- req: R1 -->
## R1 - Canonical lifecycle and request machines

<!-- req: R1.1 -->
### R1.1

The family must host request, control, and monitor listeners, await an explicit lifecycle exit, and stop owned listeners in order.

<!-- req: R1.2 -->
### R1.2

The apply machine must sequence write, validate, apply, verify, and rollback, with terminal outcomes Done, Rejected, RolledBack, and Failed.

<!-- req: R1.3 -->
### R1.3

The rollout-read machine must distinguish Complete, Progressing, and Unavailable by following the bounded status poll with a Deployment read.

<!-- req: R2 -->
## R2 - Shared vocabulary

<!-- req: R2.1 -->
### R2.1

tools.yaml, apply-tools.yaml, and rollout-tools.yaml must remain name-only selections.

<!-- req: R2.2 -->
### R2.2

declarations.yaml must own lifecycle host ToolDefs and apply-declarations.yaml must own the workspace write ToolDef.

<!-- req: R2.3 -->
### R2.3

ToolDefs must declare signals, side effects, reversibility, errors, and compensation contracts.

<!-- req: R3 -->
## R3 - Application-owned boundary

<!-- req: R3.1 -->
### R3.1

Chart name, schema path, release, namespace, Deployment names, Helm and kubectl arguments, endpoint paths and ports, and authorization policy must remain application-owned.

<!-- req: R3.2 -->
### R3.2

A wrapper must select canonical machines and shared declarations while supplying only local REST and rollout declarations.

<!-- req: R3.3 -->
### R3.3

Request input must not select a profile, chart, schema, release, namespace, Deployment, endpoint, credential, or rollout command.

<!-- req: R4 -->
## R4 - Validation and compensation

<!-- req: R4.1 -->
### R4.1

Validation must complete before apply and a validation failure must reach Rejected without mutation.

<!-- req: R4.2 -->
### R4.2

Apply must be followed by explicit verification of every application-selected rollout target.

<!-- req: R4.3 -->
### R4.3

A verification failure must invoke rollback; successful compensation reaches RolledBack and failed compensation reaches Failed.

<!-- req: R5 -->
## R5 - Portability, release, and ownership

<!-- req: R5.1 -->
### R5.1

The catalog profile and request variants must close every non-runtime reference and boot independently of an application chart.

<!-- req: R5.2 -->
### R5.2

Applications must package the canonical closure and local wrapper with source provenance and no copied promoted file.

<!-- req: R5.3 -->
### R5.3

The family path and machine, signal, terminal, request, and selection contracts are part of the current root v0.* compatibility surface.

<!-- req: R5.4 -->
### R5.4

Repository statistics must count applier once under applications/catalog and classify application appliers as composition wrappers.

<!-- req: R5.5 -->
### R5.5

Every profile variant the family ships must carry a role-realization binding, classified so it does not count toward role coverage. The repository discovers profiles by two rules that disagree, and the pair is what makes an unregistered variant look correct: the catalog family inventory reads only files named profile.yaml, while the applications-level realization model reads profile.yaml, profile-*.yaml, and *-profile.yaml alike.

<!-- req: R6 -->
## R6 - Host-side deploy and undeploy machines

<!-- req: R6.1 -->
### R6.1

The family must host a deploy machine that sequences write, validate, apply with install semantics, verify, and rollback -- the apply machine's shape -- with terminal outcomes Deployed, Rejected, Compensated, and Failed.

<!-- req: R6.2 -->
### R6.2

The deploy machine must declare no model word, so a deploy runs as declared transitions alone, costs no tokens, and follows one fixed sequence.

<!-- req: R6.3 -->
### R6.3

The apply words must carry --atomic, so Helm deletes a failed first install and rolls back a failed upgrade before the machine routes the failure; a verify stall must compensate by rollback, successful compensation reaches Compensated, and failed compensation, including a first install that has no prior revision, reaches Failed with the release left in place.

<!-- req: R6.4 -->
### R6.4

The machine must not branch on release existence; day-0 and day-2 must differ in the declared words alone, never in the transition graph.

<!-- req: R6.5 -->
### R6.5

An undeploy path must uninstall the release and reach a terminal state that separates an absent release from a failed uninstall, and its caller must confirm the cluster is running before it trusts an absence probe.

<!-- req: R6.6 -->
### R6.6

Chart path, release, namespace, kubeconfig, and value overlays must reach the deploy words as rendered declarations the caller produces, never from process environment, and a rendered argv must carry no placeholder coordinate.

<!-- req: R6.7 -->
### R6.7

The deploy words must not carry --reuse-values, so we reproduce a deploy from the chart, the checked-in values overlay, and the request document alone.

<!-- field: non_goals -->
## Non-goals

- Own any consuming application's chart, endpoint, deployment topology, credentials, RBAC, or operator workflow.
- Implement Helm, kubectl, REST, lifecycle, or filesystem runtime behavior.
- Require one application to adopt another application's state-read extension.
- Provision a cluster; the deploy machine owns the Helm step alone, never kindrig or gcprig.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

Catalog conformance boots the standalone profile and checks lifecycle, apply, rollout-read, terminal, and name-only selection contracts.

<!-- field: traces -->
#### Traces

- [R1.1](#r11)
- [R1.2](#r12)
- [R1.3](#r13)
- [R2.1](#r21)
- [R2.2](#r22)
- [R2.3](#r23)
- [R5.1](#r51)

<!-- ac: AC2 -->
### AC2

All three consumers select the canonical machines and shared ToolDefs while retaining local REST, rollout declarations, and authority.

<!-- field: traces -->
#### Traces

- [R3.1](#r31)
- [R3.2](#r32)
- [R3.3](#r33)

<!-- ac: AC3 -->
### AC3

Fake and applicable live tests prove validation precedes mutation and verification failure is compensated.

<!-- field: traces -->
#### Traces

- [R4.1](#r41)
- [R4.2](#r42)
- [R4.3](#r43)

<!-- ac: AC4 -->
### AC4

Deterministic package, duplicate, release, role-binding, and statistics audits prove one implementation owner and current root v0 compatibility, and the role-realization audit finds a binding for every shipped profile variant under either discovery rule.

<!-- field: traces -->
#### Traces

- [R5.2](#r52)
- [R5.3](#r53)
- [R5.4](#r54)
- [R5.5](#r55)

<!-- ac: AC5 -->
### AC5

Catalog conformance drives the deploy machine over standalone exec doubles to Deployed, to Rejected on a validation failure, to Compensated on a verify failure whose rollback succeeds, and to Failed when compensation itself fails, and drives the undeploy machine to Removed and to Absent.

<!-- field: traces -->
#### Traces

- [R6.1](#r61)
- [R6.2](#r62)
- [R6.3](#r63)
- [R6.4](#r64)
- [R6.5](#r65)

<!-- ac: AC6 -->
### AC6

A guard test scans every rendered deploy argv and fails when a placeholder chart path, release, namespace, or kubeconfig survives rendering, and a second guard fails when --reuse-values appears.

<!-- field: traces -->
#### Traces

- [R6.6](#r66)
- [R6.7](#r67)

<!-- field: platform_dependencies | form: "typed table" -->
## Platform dependencies

| srd_id | provider | symbols_used |
| --- | --- | --- |
| srd003-application-consistency | applications | [R2, R4, R5, R7] |
