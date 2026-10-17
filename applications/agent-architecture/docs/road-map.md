---
id: agent-architecture-roadmap
title: Agent Architecture Application Road Map
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

Release 00.0 establishes the standalone guided demo. Release 01.0 adds
trace collection through the canonical collector agent so the presenter can
show the trace waterfall as part of the demo story. Both releases reference
the canonical catalog profiles, pin the Go present tool, and provide
deterministic Mage evidence, extending shared srd001 and srd002 R1 through R3
by reference. Releases 02.0 and 03.0 provide the implemented Kubernetes deployment: a
chart running the documentation-curator and the collector from mounted
catalog-owned closures, a values schema, and a bounded kind smoke, citing
srd002 R4 through R6 and the application-owned srd001-helm-deployment. Release
04.0 provides the implemented standalone declarative applier
(srd002-applier). The package, render, bounded kind, and applier integration
targets are checked in; live execution remains prerequisite-gated.

<!-- release: 00.0 | status: partial -->
## 00.0 - Guided Knowledge Manager presentation

An operator serves the slide deck, launches the canonical
documentation-curator from an explicit catalog root, opens its
catalog-owned browser interface, and requests clean shutdown through the
catalog-owned declarative lifecycle-exit client. The application
reports composition only and contributes no canonical agents.

<!-- field: capabilities | form: "typed table" -->
### Capabilities

| id | status | services | live_observation_status | realized_in |
| --- | --- | --- | --- | --- |
| runnable_module | implemented |  |  |  |
| managed_service | partial | [documentation-curator, Go present] | dependency_gated |  |
| packaged | implemented |  |  | "02.0" |
| helm_managed | implemented |  |  | "02.0" |
| kind_demo | implemented |  |  | "03.0" |

<!-- field: shared_requirements | form: "typed table" -->
### Shared requirements

| owner | id | requirements |
| --- | --- | --- |
| applications | srd001-application-module-contract | [R1, R2, R3, R4, R5, R6] |
| applications | srd002-managed-application-services | [R1, R2, R3] |

<!-- field: external_requirements -->
### External requirements

| owner | id |
| --- | --- |
| applications/catalog | srd011-knowledge-manager |
| applications/catalog | rel07.0-uc001-documentation-curator-profile-ux |

<!-- field: use_cases -->
### Use cases

| id | summary | status | live_observation_status |
| --- | --- | --- | --- |
| [rel00.0-uc001-guided-agent-architecture](specs/use-cases/rel00.0-uc001-guided-agent-architecture.md) | An operator serves the deck, starts the canonical curator, opens its UI, and exits it declaratively. | partial | dependency_gated |

<!-- field: test_suites -->
### Test suites

| id | summary | status | live_observation_status |
| --- | --- | --- | --- |
| [test-rel00.0-guided-agent-architecture](specs/test-suites/test-rel00.0-guided-agent-architecture.md) | Deterministic command and contract tests plus dependency-gated launch and exit evidence. | partial | dependency_gated |

<!-- field: evidence | form: "typed table" -->
### Evidence

| command | status | proves_when_run | working_directory | evidence_path | proves |
| --- | --- | --- | --- | --- | --- |
| go test ./magefiles |  |  |  |  | Exact run and presentation commands, working directories, root resolution, trace reciprocity, and failure messages. |
| mage run | dependency_gated | The canonical catalog profile starts through agent-core with explicit profile, directory, and development core-root arguments. |  |  |  |
| mage presentation | dependency_gated | The pinned Go present tool serves agent-architecture.slide from the application root. |  |  |  |
| mage audit |  |  |  |  | Required documents, reciprocal traces, portable paths, canonical references, and stale-path exclusions satisfy the application contract. |
| mage stats |  |  |  |  | Valid composition-only JSON reports agents_contributed as zero and omits an agents section. |
| mage audit |  |  | repository root | magefiles/build.go auditParticipants and magefiles/example_modules_test.go TestApplicationModulesParticipateInAudit | Root audit dispatches applications/agent-architecture. |
| mage test |  |  | repository root | magefiles/example_modules_test.go TestTestSubModulesDispatchesApplicationModules | Root tests dispatch the demo's Go tests. |
| mage stats |  |  | repository root | magefiles/build.go statsParticipants and magefiles/stats_test.go TestSumAgentsTotals | Root stats accepts the composition-only application section without double-counting the curator. |
| go test ./magefiles -run '^Test(ApplicationModulesParticipateInAudit\|KnowledgeManager\|Orchestration\|ReleaseTags\|StatsParticipants\|SumAgentsTotals\|TestSubModulesDispatchesApplicationModules)' |  |  | repository root | magefiles/example_modules_test.go, magefiles/stats_test.go, and magefiles/tag_test.go | Stable root audit, test, stats, and release participation. |

<!-- field: evidence_boundary -->
### Evidence boundary

Checked-in deterministic tests prove command construction, configuration,
traceability, accounting, and root dispatch. They do not prove live
startup, browser rendering, role readiness, liveness, representative
failure, diagnostics, telemetry, shutdown, in-flight handling, or listener
release. Those observations remain dependency_gated on available ports, a
buildable agent-core checkout, and operator execution. Until they run,
srd002 AC1 through AC3 and the aggregate managed-service, use-case, suite,
and release statuses remain partial.

<!-- release: 01.0 | status: implemented -->
## 01.0 - Demo trace collection

The demo launch path starts the canonical collector agent in spool mode
before the documentation-curator, exports OTLP traces from the curator
to the local collector, and stops the collector cleanly via its control
exit route after the curator terminates. The presenter opens the collector
query surface to show the trace waterfall.

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| [rel01.0-uc001-demo-trace-collection](specs/use-cases/rel01.0-uc001-demo-trace-collection.md) | The demo run starts the collector, exports curator traces, and stops the collector cleanly. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| [test-rel01.0-demo-trace-collection](specs/test-suites/test-rel01.0-demo-trace-collection.md) | Deterministic tracing toggle, collector command, and curator OTLP flag tests. | implemented |

<!-- field: evidence -->
### Evidence

| command | proves |
| --- | --- |
| go test ./magefiles -run '^TestTracingEnabled' | Tracing is on by default and disabled by tracing:false in demo.yaml. |
| go test ./magefiles -run '^TestCollectorCommand' | Collector command targets the canonical profile with correct env and ports. |
| go test ./magefiles -run '^TestTracedCurator' | Curator receives OTLP endpoint and service name flags. |

<!-- release: 02.0 | status: implemented -->
## 02.0 - Packaged chart render and validation

A chart runs the documentation-curator and the collector as persistent
Kubernetes workloads from mounted catalog-owned profile closures. A mage
preparation target stages and checksum-validates the closures, a
values.schema.json and template guards own the deployment values, and a
package target emits one self-contained archive that lints and renders
without a source checkout. The chart deploys canonical profiles by
reference and copies neither canonical implementation. Oversized curator UI
assets are provisioned out-of-release from manifest-declared catalog roots;
the collector UI remains in its bounded role package.

<!-- field: capabilities -->
### Capabilities

| id | status |
| --- | --- |
| packaged | implemented |
| helm_managed | implemented |

<!-- field: application_srds -->
### Application SRDs

- srd001-helm-deployment

<!-- field: shared_requirements | form: "typed table" -->
### Shared requirements

| owner | id | requirements |
| --- | --- | --- |
| applications | srd002-managed-application-services | [R4, R5] |

<!-- field: external_requirements -->
### External requirements

| owner | id |
| --- | --- |
| applications/catalog | srd011-knowledge-manager |

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| [rel02.0-uc001-render-packaged-chart](specs/use-cases/rel02.0-uc001-render-packaged-chart.md) | An operator prepares the catalog-owned closures and renders and validates the packaged chart. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| [test-rel02.0-render-packaged-chart](specs/test-suites/test-rel02.0-render-packaged-chart.md) | Topology render, manifest-derived closure and checksum, values schema and guards, and self-contained archive evidence. | implemented |

<!-- field: evidence -->
### Evidence

| command | proves |
| --- | --- |
| go test ./magefiles | Manifest mutations control package output, undeclared roots fail, prepared provenance validates, and package/render/UI/applier guards pass. |
| mage helm:package | The manifest-derived self-contained chart archive lints and renders. |

<!-- release: 03.0 | status: implemented -->
## 03.0 - Live Helm smoke

A bounded kind target installs only the validated archive, builds and loads
the runtime image the chart references, and mounts the packaged
curator and collector closures. The curator and collector Deployments and
their real lifecycle health endpoints become ready, the curator serves its
documentation interface, a declarative lifecycle-exit request stops the
curator cleanly, and the collector retains a trace the curator exported. The
target skips on missing prerequisites and cleans up on every terminal path.
The target is checked in and skips truthfully when a required live
prerequisite is unavailable.

<!-- field: capabilities -->
### Capabilities

| id | status |
| --- | --- |
| kind_demo | implemented |

<!-- field: application_srds -->
### Application SRDs

- srd001-helm-deployment

<!-- field: shared_requirements | form: "typed table" -->
### Shared requirements

| owner | id | requirements |
| --- | --- | --- |
| applications | srd002-managed-application-services | [R6] |

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| [rel03.0-uc001-live-helm-smoke](specs/use-cases/rel03.0-uc001-live-helm-smoke.md) | A developer proves the packaged topology on a disposable kind cluster. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| [test-rel03.0-live-helm-smoke](specs/test-suites/test-rel03.0-live-helm-smoke.md) | Packaged-archive install, real readiness and clean exit, retained trace, and skip and cleanup evidence. | implemented |

<!-- field: evidence -->
### Evidence

| command | proves |
| --- | --- |
| mage integration:helmSmoke | The bounded prerequisite-gated kind workflow installs the archive, observes workloads and lifecycle exit, and cleans up. |
| go test ./magefiles -run '^Test(SelectCuratorPod\|CuratorShutdownObserved)$' | Deterministic pod selection and clean lifecycle termination classification. |

<!-- release: 04.0 | status: implemented -->
## 04.0 - Declarative applier agent

A standalone declarative applier actuates the chart. An authorized operator
or CI caller submits a decided values-plane patch to the applier's apply
endpoint; the applier validates it against the chart's values.schema.json by
a helm dry-run, applies it with helm upgrade --reuse-values returning
without waiting, verifies the curator and collector Deployments via kubectl
rollout status, rolls back on a verify stall, and reports an apply-command
failure as failed. The applier alone holds helm and kubectl bound as exec
tools, the curator and collector carry none, and its apply surface is
network-gated off the serving pods. The Agent Architecture runs no control
plane, so there is no creator, orchestrator, or panel; the applier runs the
agent-core image with helm and kubectl delivered by a pinned CLI donor
(GH-2222), so no image diverges from the profile-free image contract. The profile, chart resources, recording-fake tracer, schema
tests, and prerequisite-gated live target are checked in. The
application-local applier wrapper composes the canonical catalog applier
while retaining chart-specific REST, command, schema, and rollout bindings.

<!-- field: application_srds -->
### Application SRDs

- srd002-applier

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| [rel04.0-uc001-declarative-applier](specs/use-cases/rel04.0-uc001-declarative-applier.md) | An operator or CI caller submits a decided values patch that the applier validates, applies, verifies, and rolls back on failure. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| [test-rel04.0-declarative-applier](specs/test-suites/test-rel04.0-declarative-applier.md) | Recording-fake terminal coverage, real chart-schema rejection, applier-only CLI authority, broken-read classification, and a live on-cluster rollout and rollback target. | implemented |

<!-- field: evidence -->
### Evidence

| command | proves |
| --- | --- |
| go test ./magefiles -run '^TestApplier' | Local applier machine, schema, render, authority, rollout, and bounded infrastructure contracts. |
| mage integration:applier | Recording fakes drive every declared terminal and enforce values-plane authority. |
| mage integration:applierLive | The prerequisite-gated kind target exercises real upgrade, verification, rollback, and schema rejection. |

<!-- release: 99.0 | status: planned -->
## 99.0 - Fleet observer adoption

Adopt the catalog-owned observer family (catalog srd021-observer,
monitor/fleet-observer) as an agent-architecture workload pointed at this
application's release label selector and namespace, giving the curator and
collector topology a fleet view without forking the profile. Unscheduled;
specification of the catalog family lands first (GH-1247). No new
environment variables.

<!-- field: planned_work -->
### Planned work

| id | summary | status |
| --- | --- | --- |
| gh-1247 | Specify the cross-application observer family this application adopts. | planned |
