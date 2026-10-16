---
id: rel00.0-uc001-guided-agent-architecture
title: Operator Runs the Guided Agent Architecture
test_suite: test-rel00.0-guided-agent-architecture
status: partial
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An operator serves agent-architecture.slide with the application's pinned Go
present tool. From the standalone application root, the operator launches the
canonical catalog documentation-curator through agent-core, opens its
catalog-owned browser interface, follows the deck, and invokes the
catalog-owned declarative lifecycle-exit client. The curator reaches
Done and releases its listeners when the dependency-gated live flow succeeds.

<!-- field: actor -->
## Actor

A developer or presenter demonstrating the canonical Knowledge Manager profile.

<!-- field: trigger -->
## Trigger

The actor checks out applications/agent-architecture,
applications/catalog, and agent-core, then starts the guided Release 00.0
presentation.

<!-- field: preconditions -->
## Preconditions

- The application Go module has downloaded its version-pinned Go present tool dependency.
- The demo.yaml catalog_root names a catalog checkout or defaults to the monorepo ../catalog.
- The demo.yaml core_root names an agent-core checkout or defaults to the monorepo ../../agent-core.
- The canonical documentation-curator profile exists at the catalog_root's agents/knowledge-manager/documentation-curator/profile.yaml.
- Local presentation, documentation, control, and monitor ports are available.

<!-- field: evidence -->
## Evidence

```yaml
deterministic_status: implemented
deterministic_command: go test ./magefiles
proves: |
    The checked-in tests prove exact command construction, portable root
    resolution, lifecycle-exit configuration, reciprocal traces, and
    composition-only accounting.
live_observation_status: dependency_gated
live_steps: [F1, F2, F3, F5, F6, F7, F8]
limitation: |
    Serving the deck, starting the curator, rendering the catalog-owned browser
    interface, posting to a running control endpoint, observing Done, and
    observing readiness, liveness, representative failure, diagnostics,
    telemetry, shutdown, in-flight handling, and listener release require
    available ports, a buildable agent-core checkout, and operator execution.
    Deterministic command construction does not prove those live observations
    or srd002 AC1 through AC3.
```

<!-- field: flow -->
## Flow

- F1: The operator runs mage presentation from the application root; Mage invokes the pinned Go present tool for agent-architecture.slide with the application root as its working directory.
- F2: The operator opens the presentation and follows its Knowledge Manager launch step.
- F3: The operator runs mage run; Mage resolves catalog_root and core_root from demo.yaml and builds or locates the agent-core binary.
- F4: The run target launches the catalog_root's agents/knowledge-manager/documentation-curator/profile.yaml with explicit --profile, --directory, and development --core-root arguments, using the catalog root as the documentation workspace.
- F5: The operator opens the documentation interface on localhost:18081 and may inspect the monitor interface on localhost:18084.
- F6: The operator invokes the catalog-owned ../catalog/agents/lifecycle-exit/profile.yaml through the same agent-core interpreter.
- F7: The lifecycle-exit machine posts its fixed reason to 127.0.0.1:18082/api/lifecycle/exit, accepts HTTP 202, and reaches its Done terminal.
- F8: The canonical curator reaches Done and releases its documentation, control, and monitor listeners.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | architecture-agent-architecture interfaces | Mage presentation and run expose the operator entry points and explicit checkout roots. |
| T2 | architecture-agent-architecture components documentation-curator | The application launches the only canonical profile and UI from applications/catalog. |
| T3 | srd011-knowledge-manager AC1, AC2, AC3 | The catalog requirement defines the configured UX, runtime dependencies, and lifecycle outcomes this application presents. |
| T4 | rel07.0-uc001-documentation-curator-profile-ux S1, S3, S4 | Catalog integration evidence owns profile loading, UI serving, and controlled listener shutdown. |
| T5 | test-rel00.0-guided-agent-architecture | The reciprocal suite checks presentation, launch, lifecycle exit, audit, and stats contracts. |
| T6 | srd001-application-module-contract AC1, AC2, AC3, AC4, AC5, AC6 | The guided flow is the local composition and evidence surface for the shared runnable-module contract. |
| T7 | srd002-managed-application-services AC1, AC2, AC3 | The documentation-curator and Go present servers are the two managed services in this local presentation. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The tested mage presentation command invokes the version-pinned Go present tool for agent-architecture.slide from the application root; serving the deck is dependency-gated.

<!-- field: traces -->
#### Traces

- architecture-agent-architecture design_decisions 3
- test-rel00.0-guided-agent-architecture presentation_command

<!-- sc: S2 -->
### S2

The tested mage run command targets the canonical catalog documentation-curator through agent-core with explicit profile, directory, and development core-root arguments; live launch is dependency-gated.

<!-- field: traces -->
#### Traces

- architecture-agent-architecture design_decisions 1
- architecture-agent-architecture design_decisions 4
- srd011-knowledge-manager AC1
- test-rel00.0-guided-agent-architecture canonical_curator_command

<!-- sc: S3 -->
### S3

The application contains no profile or UI copy; observing the catalog-owned Knowledge Manager interface in a browser is dependency-gated.

<!-- field: traces -->
#### Traces

- srd011-knowledge-manager AC1
- rel07.0-uc001-documentation-curator-profile-ux S3
- test-rel00.0-guided-agent-architecture canonical_curator_command

<!-- sc: S4 -->
### S4

The catalog-owned declarative lifecycle-exit client defines the fixed request and accepted HTTP status; observing the live post, curator Done state, and listener release is dependency-gated.

<!-- field: traces -->
#### Traces

- architecture-agent-architecture design_decisions 2
- rel07.0-uc001-documentation-curator-profile-ux S4
- test-rel00.0-guided-agent-architecture lifecycle_exit_contract

<!-- sc: S5 -->
### S5

Application audit and stats preserve reciprocal traces, portable roots, agents_contributed as zero, and no agents section.

<!-- field: traces -->
#### Traces

- architecture-agent-architecture design_decisions 5
- test-rel00.0-guided-agent-architecture application_contract
- srd001-application-module-contract AC3
- srd001-application-module-contract AC5

<!-- sc: S6 -->
### S6

The curator and Go present servers require bounded startup, role-readiness, liveness, failure, diagnostics, telemetry, configuration, and shutdown evidence appropriate to their local roles; those dependency-gated observations must run before srd002 AC1 through AC3 are satisfied.

<!-- field: traces -->
#### Traces

- architecture-agent-architecture capability_classification managed_services
- srd002-managed-application-services AC1
- srd002-managed-application-services AC2
- srd002-managed-application-services AC3
- test-rel00.0-guided-agent-architecture managed_service_contract

<!-- field: out_of_scope -->
## Out of scope

- Changing the canonical documentation-curator profile, UI, or catalog tests.
- Re-testing generic agent-core REST route matching or documentation hosting.
- Packaging a copy of the catalog profile in the application source tree.
- Automating browser interaction or visual slide review.
- Defining production deployment or authentication for the local demo.
