---
id: srd003-application-consistency
title: Application Composition Consistency and Conformance
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

The showcase applications express equivalent composition, actor, package,
deployment, user-interface, and status concepts with different paths and
inventories. A reviewer cannot determine from one authority which profiles
are direct roots, which paths are mounted at runtime, which assets are
application-owned, or which optional capabilities apply. Separate package
and Helm inventories can also drift from the declared composition.

<!-- field: goals -->
## Goals

- G1: Define one application-owned composition authority for runnable and audit-only applications.
- G2: Normalize actor, profile, machine, tool, declaration, REST, request, and UI ownership without requiring empty optional directories.
- G3: Derive package closure, deployment roots, and provenance from the composition manifest.
- G4: Make canonical promotion, shared UI-token consumption, README shape, and conformance status evidence-backed.
- G5: Preserve explicit boundaries among applications, applications/catalog, and agent-core.

<!-- field: applicability | form: "tree" -->
## Applicability

- baseline:
  R1, R2, R3, R5, R6, and R7 apply to every runnable application and to an
  application entering the root audit-only transition. Audit-only means the
  documentation and manifest gates are executable; it does not claim a
  runnable composition.
- capability_gates:
  - managed_service: Service lifecycle artifacts and evidence apply only when declared.
  - packaged: R4 applies only when a distributable artifact is produced.
  - helm_managed: Helm deployment entries apply only when a chart is shipped and imply packaged.
  - kind_demo: kind evidence applies only when declared and remains governed by ENG01.
  - ui: UI entries and shared-token checks apply only when an application or actor ships UI assets.

<!-- req: R1 -->
## R1 - Composition manifest authority

<!-- field: applicability -->
### Applicability

runnable_module or audit_only

<!-- req: R1.1 -->
### R1.1

Each runnable or audit-only application shall own exactly one
agents/application.yaml. It shall be the composition authority for
direct profile roots and application capability claims.

<!-- field: evidence -->
#### Evidence

The file parses, its application field is the declared identity
without requiring the checkout directory's base name to match, and
local audit reads this path rather than reconstructing composition
from prose, build code, or Helm templates.

<!-- req: R1.2 -->
### R1.2

The manifest shall declare schema_version, application, ownership,
module_status, and a capability map. Each capability shall use the
common status vocabulary and shall name evidence when its status is
partial or implemented.

<!-- field: evidence -->
#### Evidence

Schema validation rejects missing identity fields, unknown status
values, an implemented claim without evidence, and an audit_only
module that claims runnable behavior.

<!-- req: R1.3 -->
### R1.3

The manifest shall list every direct catalog or local profile root
once. Each entry shall have a stable actor or root identifier,
ownership, source, and runtime_path. Catalog entries shall also name
compatible_release; local entries shall remain below the application
root.

<!-- field: evidence -->
#### Evidence

Structural tests reject duplicate identifiers or runtime paths,
missing sources, unsafe traversal, absolute source paths, and catalog
references without compatibility.

<!-- req: R1.4 -->
### R1.4

Runtime mount paths and profile runtime paths shall be explicit,
normalized, non-overlapping, and independent of a developer checkout.
Request input shall not select a source root, runtime root, profile,
compatibility release, or credential.

<!-- field: evidence -->
#### Evidence

Manifest validation and negative tests reject absolute source paths,
traversal, collisions, and authority-bearing request fields.

<!-- req: R1.5 -->
### R1.5

Deployment entries shall identify only workloads or finite entry
profiles the application actually deploys. Each entry shall refer to
a declared direct root and may add workload identity, profile path,
mount, and capability-specific settings without becoming a second
composition inventory.

<!-- field: evidence -->
#### Evidence

Render and structural tests resolve every deployment entry to the
manifest and reject a chart-selected undeclared root.

<!-- req: R1.6 -->
### R1.6

An application that ships a user interface shall declare each UI
asset root, its owning actor or application surface, its serving REST
definition, runtime destination, and shared-token policy. Applications
without UI assets shall omit the UI section.

<!-- field: evidence -->
#### Evidence

UI entries resolve to files and a static_assets binding. Baseline
conformance accepts an absent UI section when no UI is shipped.

<!-- req: R2 -->
## R2 - Actor ownership and directory grammar

<!-- field: applicability -->
### Applicability

runnable_module or audit_only

<!-- req: R2.1 -->
### R2.1

Application-owned production actors and wrappers shall live under
agents/<actor>/. The lower-kebab actor name shall describe the runtime
identity or purpose and shall retain role bindings in the shared
agent-role realization model.

<!-- field: evidence -->
#### Evidence

Manifest local roots resolve below agents/<actor>/ and the semantic
model classifies each shipped production profile or its inheritance.

<!-- req: R2.2 -->
### R2.2

Canonical reusable implementations shall live under
applications/catalog/agents/. An application shall consume them as
catalog-owned roots or thin local wrappers and shall not copy their
machines, declarations, tools, family SRDs, or canonical UI.

<!-- field: evidence -->
#### Evidence

Ownership and closure tests distinguish catalog from local roots and
reject copied canonical assets.

<!-- req: R2.3 -->
### R2.3

Composition-only applications may contain agents/application.yaml and
no local actor directory. Audit-only applications may declare planned
roots without creating them only when each root is explicitly planned
and excluded from executable closure. Optional actor, service,
package, Helm, kind, and UI directories shall not be required empty.

<!-- field: evidence -->
#### Evidence

Capability-aware tests accept absent optional directories and reject a
planned or missing root used as executable evidence.

<!-- req: R2.4 -->
### R2.4

Test profiles and fixtures shall remain under an owning actor's tests/
tree or the application's testdata/ tree and shall not enter production
deployment, role-coverage, package-root, or agent-count inventories.

<!-- field: evidence -->
#### Evidence

Discovery and manifest checks classify fixtures separately and reject
fixture paths in production deployment entries.

<!-- req: R3 -->
## R3 - Agent program and UI naming

<!-- field: applicability -->
### Applicability

each declared local or catalog profile root

<!-- req: R3.1 -->
### R3.1

profile.yaml shall be the primary loadable entry in an actor root. A
distinct operation variant shall use
<operation>-profile.yaml and the same lower-kebab operation prefix
across its related files.

<!-- field: evidence -->
#### Evidence

Profile discovery resolves primary and operation entries
deterministically and rejects unrelated prefixes for one variant.

<!-- req: R3.2 -->
### R3.2

machine.yaml shall be the primary MachineSpec. Operation or request
variants shall use <operation>-machine.yaml, including the reserved
request-machine.yaml form when the variant serves a request.

<!-- field: evidence -->
#### Evidence

Every profile machine reference resolves within its declared closure
and follows the primary or operation-prefixed form.

<!-- req: R3.3 -->
### R3.3

tools.yaml and <operation>-tools.yaml shall contain name-only tool
selections. ToolDef documents shall use declarations.yaml or
<operation>-declarations.yaml. A selection file shall not become a
second ToolDef authority.

<!-- field: evidence -->
#### Evidence

Parsing tests distinguish name-only selections from declarations and
resolve every selected word to agent-core or profile-owned ToolDefs.

<!-- req: R3.4 -->
### R3.4

rest.yaml or <operation>-rest.yaml shall own REST definitions.
Request serving variants shall use the aligned request-profile.yaml,
request-machine.yaml, request-tools.yaml, and
request-declarations.yaml names only for files that exist.

<!-- field: evidence -->
#### Evidence

Profile and REST validation resolves every referenced file and does
not require an empty request file for an actor without that surface.

<!-- req: R3.5 -->
### R3.5

Actor-served UI shall live under agents/<actor>/ui/. Application-wide
UI that has no serving actor may live under ui/ and shall name its
application owner in the manifest. Build source and generated dist
paths shall remain distinguishable.

<!-- field: evidence -->
#### Evidence

Manifest UI assets resolve below the declared owner and REST
static_assets roots select the declared runtime or generated path.

<!-- req: R4 -->
## R4 - Manifest-derived package closure and provenance

<!-- field: applicability -->
### Applicability

packaged

<!-- req: R4.1 -->
### R4.1

A package resolver shall begin only from executable direct roots and
deployment entries in agents/application.yaml. Build code, chart
values, and handwritten role lists shall not add another root.

<!-- field: evidence -->
#### Evidence

Package tests change a manifest root and observe the generated closure
change without editing a second inventory; undeclared roots fail.

<!-- req: R4.2 -->
### R4.2

The resolver shall traverse every reachable profile, machine, tool
selection, ToolDef declaration, REST definition, child profile,
configuration, and declared UI asset. It shall preserve runtime paths
and reject missing, escaping, conflicting, or unbounded members.

<!-- field: evidence -->
#### Evidence

Closure tests cover each reference-bearing asset class and package
validation runs against the assembled runtime tree.

<!-- req: R4.3 -->
### R4.3

Generated package and Helm inventories shall be derived outputs with
checksums, not maintained composition authorities. Helm preparation
and packaging shall consume the generated output associated with the
same manifest.

<!-- field: evidence -->
#### Evidence

Tests reject stale generated inventories and chart-selected profile
roots that cannot be traced to the manifest.

<!-- req: R4.4 -->
### R4.4

Package provenance shall record the manifest identity, compatible
catalog release, actual catalog revision and dirty state when
applicable, application revision, sorted inventory, and content
checksums.

<!-- field: evidence -->
#### Evidence

Repeated packaging is deterministic and distinguishes compatibility
from the actual source revision without claiming a release for a dirty
or unmatched checkout.

<!-- req: R5 -->
## R5 - Canonical promotion and shared UI tokens

<!-- field: applicability -->
### Applicability

applications that own reusable candidates or UI assets

<!-- req: R5.1 -->
### R5.1

An application-owned actor or asset is a canonical-promotion candidate
only when it is independently useful outside the application or has
at least two real consumers. Similar names or speculative reuse shall
not transfer ownership.

<!-- field: evidence -->
#### Evidence

A promotion record names consumers, responsibility, source path,
canonical destination, compatibility, and migration evidence.

<!-- req: R5.2 -->
### R5.2

Promotion shall move reusable behavior and its conformance evidence to
applications/catalog, migrate every consumer to a canonical reference,
and remove application copies. Application-specific wrappers,
endpoints, topology, credentials, and UX shall remain local.

<!-- field: evidence -->
#### Evidence

Migration and ownership tests prove one canonical implementation and
no stale copied source while preserving application-owned bindings.

<!-- req: R5.3 -->
### R5.3

A UI that claims shared design tokens shall consume
applications/ui-kit/src/tokens.css through the build import
@declarative-agents/ui-kit/tokens.css or a deterministic generated
copy (srd004 R8.2 moved the canonical source from
applications/catalog/ui/design-tokens.css). A checked-in generated token block shall
carry canonical-source markers and a byte-for-byte drift test. A UI
that does not consume the canonical values shall not claim shared-token
conformance.

<!-- field: evidence -->
#### Evidence

UI build or drift tests compare the consumed token set with the
canonical source and fail on a local value change.

<!-- req: R6 -->
## R6 - README and status vocabulary

<!-- field: applicability -->
### Applicability

runnable_module or audit_only

<!-- req: R6.1 -->
### R6.1

Each application README shall contain Purpose, Status, Composition,
Capabilities, Ownership Boundaries, Run or Planned Entry Points,
Verification, and Documentation sections. A section may state not
applicable instead of pointing to an empty directory.

<!-- field: evidence -->
#### Evidence

Documentation audit finds each heading and verifies that referenced
paths and commands exist or are explicitly marked planned.

<!-- req: R6.2 -->
### R6.2

Module and capability status shall use planned, audit_only, partial,
dependency_gated, implemented, or not_applicable. audit_only means
manifest and documentation audit participation only. partial and
dependency_gated shall name missing evidence; implemented shall name
executable evidence.

<!-- field: evidence -->
#### Evidence

Manifest, README, architecture, roadmap, root orchestration, and Mage
targets use compatible statuses and reject unsupported synonyms in
conformance fields.

<!-- req: R6.3 -->
### R6.3

Status shall be capability-specific. An implemented baseline shall not
imply implemented managed_service, packaged, helm_managed, kind_demo,
or UI evidence, and a partial optional capability shall not downgrade
an otherwise implemented baseline.

<!-- field: evidence -->
#### Evidence

Cross-document tests compare each capability independently and require
evidence or a named gap for every non-not_applicable status.

<!-- req: R6.4 -->
### R6.4

Root classification shall agree with module status. Runnable modules
enter applicationModules; documentation-first modules enter
auditOnlyApplicationModules only after they own agents/application.yaml
and a local audit target. A prose-only unregistered directory shall not
call itself audit_only until that transition is executable.

<!-- field: evidence -->
#### Evidence

Root tests compare the manifest and README with applicationModules and
auditOnlyApplicationModules and verify every classified path exists
with the evidence required by its declared status.

<!-- req: R7 -->
## R7 - Ownership boundary closure

<!-- field: applicability -->
### Applicability

all applications

<!-- req: R7.1 -->
### R7.1

applications/catalog shall own canonical reusable profiles, machines,
declarations, tools, family requirements, and canonical agent UI.
Shared UI tokens and shared presentation code belong to
applications/ui-kit (srd004 R8). It shall not own a consuming application's composition,
deployment, credentials, scenario, or operator workflow.

<!-- field: evidence -->
#### Evidence

Manifest ownership and source paths distinguish catalog roots from
application roots and repository statistics count the canonical owner
once.

<!-- req: R7.2 -->
### R7.2

Each application shall own agents/application.yaml, local actors and
wrappers, topology, deployment entries, app-specific UI, capability
claims, README, local specifications, and end-to-end evidence.

<!-- field: evidence -->
#### Evidence

Every app-owned path remains beneath applications/<application>/ and
reciprocal specifications point to the shared contract by ID.

<!-- req: R7.3 -->
### R7.3

agent-core shall own profile, machine, tool, REST, lifecycle,
telemetry, checkpoint, and execution semantics. This contract shall
normalize application layout and references without redefining those
runtime formats or requiring application-specific runtime code.

<!-- field: evidence -->
#### Evidence

Application documents cite agent-core authorities; conformance checks
load application profiles through the existing runtime validators.

<!-- field: non_goals -->
## Non-goals

- Require identical optional capabilities across the three applications.
- Require a managed service, package, chart, kind demo, UI, request variant, local actor, or generated directory when the corresponding capability is absent.
- Change agent-core profile, machine, ToolDef, REST, lifecycle, or telemetry semantics.
- Promote an application asset without demonstrated reuse and a coordinated migration.
- Treat a planned path, prose claim, or empty directory as executable evidence.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

Every runnable or audit-only application has one validated
agents/application.yaml that identifies direct roots, runtime paths,
compatibility, applicable deployment and UI entries, capabilities, and
evidence-backed module status.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R1.4
- R1.5
- R1.6

<!-- field: applicability -->
#### Applicability

runnable_module or audit_only; deployment and UI fields only when present

<!-- ac: AC2 -->
### AC2

Application actors, catalog implementations, fixtures, and optional
directories follow the ownership grammar without empty capability stubs.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3
- R2.4

<!-- field: applicability -->
#### Applicability

runnable_module or audit_only

<!-- ac: AC3 -->
### AC3

Profile, machine, tools, ToolDef declarations, REST, request variants, and
UI assets use consistent primary or operation-prefixed names and every
reference resolves.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3
- R3.4
- R3.5

<!-- field: applicability -->
#### Applicability

each shipped surface; absent optional surfaces are not applicable

<!-- ac: AC4 -->
### AC4

A packaged application derives roots, closure, generated inventories, and
compatibility-separated provenance from agents/application.yaml only.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
- R4.3
- R4.4

<!-- field: applicability -->
#### Applicability

packaged

<!-- ac: AC5 -->
### AC5

Promotion occurs only for demonstrated reusable assets through a complete
consumer migration, and shared-token claims prove actual canonical token
consumption.

<!-- field: traces -->
#### Traces

- R5.1
- R5.2
- R5.3

<!-- field: applicability -->
#### Applicability

promotion candidates or shipped UI

<!-- ac: AC6 -->
### AC6

README sections, module status, capability status, evidence, and root
orchestration agree, including the executable meaning of audit_only,
partial, and dependency_gated.

<!-- field: traces -->
#### Traces

- R6.1
- R6.2
- R6.3
- R6.4

<!-- field: applicability -->
#### Applicability

runnable_module or audit_only

<!-- ac: AC7 -->
### AC7

Catalog, application, and agent-core ownership remain explicit and no
normalized layout claim transfers runtime semantics or duplicates a
canonical implementation.

<!-- field: traces -->
#### Traces

- R7.1
- R7.2
- R7.3

<!-- field: applicability -->
#### Applicability

all applications

<!-- field: use_case_traces -->
## Use case traces

- rel14.0-uc001-application-consistency

<!-- field: test_suite_traces -->
## Test suite traces

- test-rel14.0-application-consistency
