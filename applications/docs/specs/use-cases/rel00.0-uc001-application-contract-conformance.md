---
id: rel00.0-uc001-application-contract-conformance
title: Verify a Runnable Application Against the Shared Contract
test_suite: test-rel00.0-application-contract-conformance
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A maintainer classifies a module as a runnable application, records only the
optional capabilities it implements, and runs structural conformance. The
checks preserve the catalog boundary, validate local documentation and
reciprocal traces, verify root Mage participation, and connect each optional
claim to executable application evidence.

<!-- field: actor -->
## Actor

A maintainer adding or reviewing a runnable application module.

<!-- field: trigger -->
## Trigger

A module is added under applications/, changes its capability set, or enters a
root release.

<!-- field: preconditions -->
## Preconditions

- The candidate has one applications/<application>/ root.
- applications/catalog is available as the canonical reusable catalog.
- Root Mage orchestration and Go tests are available.
- The candidate declares runnable_module and any implemented optional capabilities in its architecture.

<!-- field: current_candidates -->
## Current candidates

```yaml
- path: applications/chatbot-mesh
  ownership: agent-owning
- path: applications/coding-agent
  ownership: composition-only
- path: applications/agent-architecture
  ownership: composition-only
  runnable_module_status: implemented
  managed_service_status: partial
  live_managed_service_evidence: dependency_gated
  packaged_status: planned
  helm_managed_status: planned
  kind_demo_status: planned
- path: applications/catalog
  classification: reusable catalog
  runnable: false
```

<!-- field: flow -->
## Flow

- F1: The maintainer identifies an executable application-owned composition and
  classifies ownership as agent-owning or composition-only. The reusable
  catalog remains outside the candidate set.
- F2: Structural checks validate canonical catalog references, app-owned
  composition paths, and ownership-aware stats with no duplicate canonical
  agents. For agent-architecture, local stats report one canonical
  reference, zero contributed agents, and no agents section, so the root
  total remains 15 canonical agents.
- F3: YAML checks parse the application vision, architecture, roadmap,
  specification index, SRDs, use cases, and test suites; indexed paths and
  reciprocal use-case/test-suite traces resolve.
- F4: Root tests verify stable application paths and audit, test, stats, and
  release participation. The application-local audit verifies its corpus
  and executable evidence. Root orchestration registers
  applications/agent-architecture alongside chatbot-mesh and
  coding-agent while applications/catalog remains a reusable submodule.
- F5: If managed_service is declared, the maintainer runs lifecycle, truthful
  health, authority, telemetry, configuration, and graceful-shutdown checks.
  A deterministic command or configuration check may be implemented while
  the aggregate capability remains partial; dependency-gated live evidence
  does not satisfy srd002 AC1 through AC3.
- F6: If packaged or helm_managed is declared, the maintainer runs deterministic
  package, safety, provenance, schema, render, archive, and role-lifecycle
  checks. Undeclared capabilities are skipped as not applicable, not passed.
- F7: If kind_demo is declared, the maintainer runs the application's ENG01
  evidence, including preflight classification, validated-artifact install,
  cluster cleanup, observability, and failure capture. No ENG01 check is
  required without kind_demo.
- F8: Root mage audit and mage stats complete. The maintainer records which
  evidence was executed and which capability groups were not applicable.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd001-application-module-contract AC1 | Classifies the runnable module and preserves the catalog boundary. |
| T2 | srd001-application-module-contract AC2 | Verifies canonical references, composition ownership, and package provenance when applicable. |
| T3 | srd001-application-module-contract AC3 | Parses local documents and verifies reciprocal traceability. |
| T4 | srd001-application-module-contract AC4 | Verifies root audit, test, stats, and release participation. |
| T5 | srd001-application-module-contract AC5 | Verifies zero duplicate agent accounting. |
| T6 | srd001-application-module-contract AC6 | Distinguishes baseline audit from conditional Mage surfaces. |
| T7 | srd002-managed-application-services AC1 | Applies lifecycle and health checks only to managed services. |
| T8 | srd002-managed-application-services AC2 | Applies authority and operator-verb checks only to managed services. |
| T9 | srd002-managed-application-services AC3 | Applies telemetry and configuration checks only to managed services. |
| T10 | srd002-managed-application-services AC4 | Applies deterministic package checks only to packaged applications. |
| T11 | srd002-managed-application-services AC5 | Applies chart checks only to Helm-managed applications. |
| T12 | srd002-managed-application-services AC6 | Applies ENG01 only to kind-demo applications and proves the exclusion for others. |
| T13 | test-rel00.0-application-contract-conformance | The suite executes baseline and capability-dependent structural evidence. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The candidate is a runnable application with explicit ownership, while
applications/catalog remains a reusable non-runnable module.

<!-- field: traces -->
#### Traces

- srd001-application-module-contract R1
- srd001-application-module-contract AC1

<!-- sc: S2 -->
### S2

Canonical references and local ownership are structurally visible, and
root stats count each canonical agent exactly once.

<!-- field: traces -->
#### Traces

- srd001-application-module-contract R2
- srd001-application-module-contract R5
- srd001-application-module-contract AC2
- srd001-application-module-contract AC5

<!-- sc: S3 -->
### S3

Required YAML parses, index paths resolve, reciprocal traces hold, and
local documents cite shared requirements rather than copying them.

<!-- field: traces -->
#### Traces

- srd001-application-module-contract R3
- srd001-application-module-contract AC3

<!-- sc: S4 -->
### S4

Root Mage dispatch and the local Mage surface match the application's
actual tests, metrics, release state, and optional capabilities.

<!-- field: traces -->
#### Traces

- srd001-application-module-contract R4
- srd001-application-module-contract R6
- srd001-application-module-contract AC4
- srd001-application-module-contract AC6

<!-- sc: S5 -->
### S5

Each managed-service claim has executable lifecycle, health, authority,
telemetry, configuration, diagnostic, and shutdown evidence.

<!-- field: traces -->
#### Traces

- srd002-managed-application-services R1
- srd002-managed-application-services R2
- srd002-managed-application-services R3
- srd002-managed-application-services AC1
- srd002-managed-application-services AC2
- srd002-managed-application-services AC3

<!-- sc: S6 -->
### S6

Each package or Helm claim has deterministic artifact, safety, provenance,
schema, render, archive, and role-lifecycle evidence.

<!-- field: traces -->
#### Traces

- srd002-managed-application-services R4
- srd002-managed-application-services R5
- srd002-managed-application-services AC4
- srd002-managed-application-services AC5

<!-- sc: S7 -->
### S7

A kind-demo claim executes ENG01 evidence; a candidate without kind_demo
records that group as not applicable.

<!-- field: traces -->
#### Traces

- srd002-managed-application-services R6
- srd002-managed-application-services AC6

<!-- field: out_of_scope -->
## Out of scope

- Validate application-specific business outcomes.
- Promote a catalog member or change catalog conformance.
- Require optional capabilities that the candidate does not declare.
- Reproduce ENG01 checks in this use case.
- Treat dependency-gated Knowledge Manager live observations as completed managed-service conformance.
