---
id: srd002-managed-application-services
title: Capability-Dependent Application Services
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Long-running services, distributable packages, Helm charts, and kind demos
need additional contracts, but those contracts are not valid baseline
obligations for every runnable application. Each requirement needs an
explicit capability and observable evidence.

<!-- field: goals -->
## Goals

- G1: Define lifecycle, health, control, telemetry, configuration, and shutdown requirements for application-owned managed services.
- G2: Define package and Helm surfaces only for applications that ship them.
- G3: Bind kind demos to ENG01 by reference without widening its scope.
- G4: Give the Agent Architecture a contract it can cite selectively.

<!-- field: capability_classes -->
## Capability classes

```yaml
managed_service:
    definition: |
        The application owns one or more long-running processes or listeners whose
        lifecycle an operator or another application component manages.
packaged:
    definition: |
        The application produces a deterministic distributable artifact, such as
        a profile bundle, archive, image, or chart package.
helm_managed:
    definition: |
        The application ships a Helm chart as an installation and lifecycle
        surface. This capability implies packaged.
kind_demo:
    definition: |
        The application ships a kind-based integration proof or operator demo.
        This capability normally accompanies helm_managed and is governed by
        ENG01.
```

<!-- req: R1 -->
## R1 - Managed-service lifecycle

<!-- field: applicability -->
### Applicability

managed_service

<!-- req: R1.1 -->
### R1.1

Each managed service shall have an application-owned start path and a
finite stop path. Startup shall fail when required configuration or
authority cannot be validated.

<!-- field: evidence -->
#### Evidence

A profile, command, or workload manifest names the start path; a test
exercises successful startup, invalid-start rejection, and stop.

<!-- req: R1.2 -->
### R1.2

Readiness shall report that the service can perform its declared
application role. Process existence, an idle sidecar, or a sleep loop
shall not satisfy readiness.

<!-- field: evidence -->
#### Evidence

The readiness probe reaches the application-owned listener or role
dependency, and a test distinguishes live-but-not-ready from ready.

<!-- req: R1.3 -->
### R1.3

Liveness shall report that the service process and owned control loop
can continue. Health endpoints shall not conceal a failed role.

<!-- field: evidence -->
#### Evidence

A bounded health test observes success and a representative failed
dependency or failed control loop.

<!-- req: R1.4 -->
### R1.4

Graceful shutdown shall stop intake, allow bounded in-flight work to
finish or report cancellation, release owned listeners and resources,
and terminate within a declared deadline.

<!-- field: evidence -->
#### Evidence

A lifecycle test holds work in flight, requests stop, observes the
declared outcome, and proves listener/resource release.

<!-- req: R2 -->
## R2 - Control and operator verbs

<!-- field: applicability -->
### Applicability

managed_service

<!-- req: R2.1 -->
### R2.1

Control operations shall use declared, validated authority. Request
input shall not introduce filesystem roots, endpoints, credentials,
profiles, or deployment authority.

<!-- field: evidence -->
#### Evidence

Config or profile declarations own authority and negative tests reject
authority-bearing request fields.

<!-- req: R2.2 -->
### R2.2

Operator verbs shall have stable names, inputs, side effects, success
evidence, failure evidence, and idempotency or compensation behavior.

<!-- field: evidence -->
#### Evidence

Application docs and mage -l name each verb; tests assert side effects
and repeat or rollback behavior.

<!-- req: R2.3 -->
### R2.3

A managed service shall expose only operator verbs its architecture
implements. package, image, Helm, doctor, demo, and kind integration
verbs shall not be inferred from managed_service alone.

<!-- field: evidence -->
#### Evidence

Capability declaration, application architecture, and mage -l list
the same surfaces.

<!-- req: R3 -->
## R3 - Telemetry and configuration

<!-- field: applicability -->
### Applicability

managed_service

<!-- req: R3.1 -->
### R3.1

Service telemetry shall identify application, role, instance, target
or scenario, source revision, and run when those dimensions exist.
Remote boundaries shall preserve the runtime's trace context contract.

<!-- field: evidence -->
#### Evidence

Integration tests inspect emitted resource attributes and connected
traces or explicitly document a signal limitation.

<!-- req: R3.2 -->
### R3.2

Health and lifecycle failures shall emit bounded diagnostic evidence
that identifies the failed role and operation without exposing secret
configuration.

<!-- field: evidence -->
#### Evidence

Failure tests assert role-specific diagnostics and secret redaction.

<!-- req: R3.3 -->
### R3.3

Configuration shall have one declared authority for each value.
Startup shall validate required shape and cross-field invariants before
accepting work.

<!-- field: evidence -->
#### Evidence

Schema, profile validation, or startup tests reject missing, malformed,
conflicting, and authority-bearing runtime input.

<!-- req: R4 -->
## R4 - Deterministic package

<!-- field: applicability -->
### Applicability

packaged

<!-- req: R4.1 -->
### R4.1

Package output shall be deterministic for identical inputs and shall
record source revision, compatibility claims, and included inventory.

<!-- field: evidence -->
#### Evidence

Repeated-package tests compare output or manifests and inspect
provenance and sorted inventory.

<!-- req: R4.2 -->
### R4.2

Packaging shall reject unsafe, missing, conflicting, or out-of-bound
inputs before publication and shall not leave a partial artifact as
the current output.

<!-- field: evidence -->
#### Evidence

Negative tests cover the application's reference and size boundaries;
publication is atomic or staged then renamed.

<!-- req: R4.3 -->
### R4.3

The application shall expose a package or equivalent artifact target
and an executable validation target for the produced artifact.

<!-- field: evidence -->
#### Evidence

mage -l lists the targets and the local audit or release gate invokes
artifact validation.

<!-- req: R5 -->
## R5 - Helm-managed application

<!-- field: applicability -->
### Applicability

helm_managed

<!-- req: R5.1 -->
### R5.1

The chart shall own a values schema or equivalent structural
validation and shall render deterministically from declared package
inputs.

<!-- field: evidence -->
#### Evidence

Helm lint/template tests cover defaults, supported variants, and
invalid values.

<!-- req: R5.2 -->
### R5.2

Workload readiness, liveness, configuration mounts, authority, and
shutdown behavior shall implement R1 through R3 for each managed role.

<!-- field: evidence -->
#### Evidence

Render tests inspect probes, arguments, mounts, security boundaries,
and termination settings; integration tests exercise role health.

<!-- req: R5.3 -->
### R5.3

The application shall expose Helm render or package validation and
shall test the packaged chart independently of source-only files.

<!-- field: evidence -->
#### Evidence

A Mage Helm target lints and renders the archive in an isolated
directory and rejects missing inventory or stale generated assets.

<!-- req: R6 -->
## R6 - kind integration and demo

<!-- field: applicability -->
### Applicability

kind_demo

<!-- req: R6.1 -->
### R6.1

The application shall conform to ENG01 for kind toolchain preflight,
cluster lifecycle, image delivery, network access, standard verbs,
persistent integration observability, and failure evidence.

<!-- field: evidence -->
#### Evidence

Application-local targets and tests cite ENG01 and exercise its
required behavior. This SRD does not reproduce the ENG01 checklist.

<!-- req: R6.2 -->
### R6.2

Integration targets shall skip with a named missing prerequisite when
ENG01 permits a skip. Requested demo targets shall fail with actionable
installation guidance when prerequisites are missing.

<!-- field: evidence -->
#### Evidence

Unit tests drive missing-tool cases for integration and demo preflight.

<!-- req: R6.3 -->
### R6.3

A kind proof shall install the artifact produced by the packaged or
Helm-managed surface rather than an unvalidated source-tree surrogate.

<!-- field: evidence -->
#### Evidence

The integration target resolves the package path and tests assert that
the installed chart or bundle is the validated artifact.

<!-- req: R6.4 -->
### R6.4

Applications without kind_demo shall have no ENG01 conformance
obligation and shall not be required to expose doctor, demo:up,
demo:down, or integration:* kind verbs.

<!-- field: evidence -->
#### Evidence

Capability classification omits kind_demo and baseline conformance
accepts the absent verbs.

<!-- field: non_goals -->
## Non-goals

- Require managed services, packaging, Helm, Kubernetes, or kind for every runnable application.
- Restate agent-core lifecycle, telemetry, or runtime configuration SRDs.
- Copy the ENG01 rig contract into this SRD.
- Define application-specific endpoints, deployment topology, or business outcomes.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

A managed service has executable start, truthful readiness and liveness,
and bounded graceful shutdown with resource release.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R1.4

<!-- field: applicability -->
#### Applicability

managed_service

<!-- ac: AC2 -->
### AC2

Managed control and operator verbs use declared authority, document their
effects, and expose no optional capability surface by implication.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3

<!-- field: applicability -->
#### Applicability

managed_service

<!-- ac: AC3 -->
### AC3

Managed telemetry carries available identity, failures produce bounded
diagnostics, and one validated configuration authority exists per value.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3

<!-- field: applicability -->
#### Applicability

managed_service

<!-- ac: AC4 -->
### AC4

A packaged application produces a deterministic, provenance-bearing,
validated artifact and rejects unsafe or partial publication.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
- R4.3

<!-- field: applicability -->
#### Applicability

packaged

<!-- ac: AC5 -->
### AC5

A Helm-managed application validates deterministic chart output, realizes
managed-role lifecycle and authority, and validates the packaged archive
independently.

<!-- field: traces -->
#### Traces

- R5.1
- R5.2
- R5.3

<!-- field: applicability -->
#### Applicability

helm_managed

<!-- ac: AC6 -->
### AC6

A kind-demo application cites and executes ENG01, distinguishes integration
skips from requested-demo failures, and installs validated artifacts;
applications without kind_demo have no such obligation.

<!-- field: traces -->
#### Traces

- R6.1
- R6.2
- R6.3
- R6.4

<!-- field: applicability -->
#### Applicability

kind_demo, with R6.4 proving the exclusion

<!-- field: external_contracts -->
## External contracts

| id | path | usage |
| --- | --- | --- |
| eng01-kind-test-demo-rig | ../../../../docs/engineering/eng01-kind-test-demo-rig.md | Governs only the kind_demo capability. |

<!-- field: use_case_traces -->
## Use case traces

- rel00.0-uc001-application-contract-conformance

<!-- field: test_suite_traces -->
## Test suite traces

- test-rel00.0-application-contract-conformance
