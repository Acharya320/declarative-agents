---
id: architecture-application-catalog
title: Application Catalog Architecture
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview | form: "tree" -->
## Overview

- summary:
  applications/catalog is this repository's versioned catalog of reusable
  declarative tool and agent blocks. It is not a universal agent library. At
  runtime, agent-core loads external agent programs from it. Applications
  compose canonical catalog members by reference and package their transitive
  closure. Profile directories, catalog tool declarations, profile-local
  overrides, conformance, integration fixtures, profile-local UIs, and
  structured docs live here. Application presentations and composition live
  in their owning modules under applications/.

  Runtime ownership remains in agent-core. That repository owns the Go
  engine, cmd/agent, builtin tool implementations, shared runtime assets, Mage
  targets, and release image. Operators supply this checkout through explicit
  --profile paths (often after bind-mounting it at a stable container path such
  as /profiles).

  Bootstrap authority for argv, mounts, and file-only configuration lives in
  agent-core docs and constitutions. Track the umbrella work there as
  epic agent-core-tj96; the runtime contract text is anchored in
  docs/specs/config-formats/runtime-contract.md (see agent-core-tj96.1).
- runtime_mount_diagram:
  ```mermaid
  flowchart LR
    OP[Developer or operator] --> RUNTIME[agent-core runtime]
    OP --> PROFILES[applications/catalog checkout]
    OP --> WORK[workspace]
    PROFILES -->|--profile path| RUNTIME
    WORK -->|--directory| RUNTIME
    RUNTIME --> CORE[agent binary and core assets]
    RUNTIME --> PROGRAM[profile-owned agent program]
  ```

<!-- field: interfaces -->
## Interfaces

<!-- interface: Catalog checkout path -->
### Catalog checkout path

Filesystem location of this repository-specific catalog on the host or
inside a container.
cmd/agent receives a profile file path through --profile; mounts such
as /profiles are a packaging convention, not a separate env-based contract
defined here.

<!-- interface: /profiles mount -->
### /profiles mount

Container-side path for this catalog. Runtime commands pass
--profile /profiles/agents/<name>/profile.yaml.

<!-- interface: --profile -->
### --profile

agent-core CLI argument that points at a profile.yaml supplied from this
catalog, a release bundle, or another external profile path.

<!-- interface: --directory -->
### --directory

agent-core CLI argument that points at the workspace the selected agent
can read, edit, test, or inspect according to its selected tools.

<!-- field: components -->
## Components

<!-- component: agents -->
### agents

Own repository-specific reusable profile families that satisfy the
membership rule, including profile.yaml files, machines, tool selections,
profile-local declarations, boundary configs, lifecycle profiles, monitor
profiles, REST profiles, control profiles, and Knowledge Manager assets.
Profiles may reference core-owned shared tool declarations through
/opt/agent-core/tools.

<!-- field: path -->
#### Path

agents/

<!-- component: tool taxonomy -->
### tool taxonomy

Define catalog tool-block membership and distinguish tools.yaml selection
lists, profile-local overrides, reusable catalog ToolDef declarations,
core declarations, core implementations, and application-local
exclusions. Catalog declarations stay with their owning agents/<family>/
closure; this path is the taxonomy, not a second declaration registry.

<!-- field: path -->
#### Path

tools/README.md

<!-- component: application presentations -->
### application presentations

Runnable applications own presentation assets, launch composition, and
application-specific end-to-end evidence while referencing canonical
catalog profiles. The Knowledge Manager presentation lives at
../agent-architecture/.

<!-- field: path -->
#### Path

../agent-architecture/

<!-- component: integration fixtures -->
### integration fixtures

Own catalog conformance and integration suites plus fixture data for
profile-owned behavior.

<!-- field: path -->
#### Path

testdata/integration/

<!-- component: documentation -->
### documentation

Own catalog vision, architecture, specification indexes, roadmap, and
issue format rules for future catalog work.

<!-- field: path -->
#### Path

docs/

<!-- field: design_decisions -->
## Design decisions

<!-- decision: 1 -->
### 1 - Profile assets are external to the runtime image

agent-core release images do not copy profile directories from this
catalog. Operators mount, check out, or unpack profile assets and pass
an explicit profile path at runtime.

<!-- field: benefits -->
#### Benefits

- Runtime images stay focused on core-owned binaries and assets.
- Profile changes can ship without rebuilding the runtime image.
- Ownership of agent behavior remains visible in repository layout.

<!-- decision: 2 -->
### 2 - Runtime code stays in agent-core

No Go runtime packages or builtin tool implementations live in the
application catalog. YAML programs select and configure capabilities that
agent-core provides.

<!-- field: benefits -->
#### Benefits

- The catalog stays reviewable as configuration and fixtures.
- Runtime changes stay covered by agent-core build and test gates.
- Profile validation can focus on paths, declarations, and contracts.

<!-- decision: 3 -->
### 3 - Fixtures use portable paths

Catalog integration fixtures and application-owned demos must use portable
roots such as catalog mount points (/profiles in runtime examples),
argv-resolved profile paths, or explicit workspace mounts.
Developer-specific absolute paths are not valid fixture contracts.

<!-- field: benefits -->
#### Benefits

- Integration targets work across machines and containers.
- Container examples match the mounted-profile contract.
- Release bundles can preserve the same directory contract.

<!-- decision: 4 -->
### 4 - Catalog membership and one canonical home

An agent is a catalog member only when it is independently useful or has
multiple consumers in this repository. A member owns a profile family,
SRD, deterministic conformance, portable closure, and sufficient
parameterization for applications to configure it without edits.
Application-only behavior stays with the owning applications/ module.
Applications reference members and assemble closures at package time;
they do not fork reusable programs.

The canonical compatibility surface is versioned by exact repository
v0.* release tags (GH-1373). Earlier applications/catalog/v0.* tags and
matching agent-profiles/v0.* names are v0 compatibility aliases, and
immutable historical legacy tags remain valid.
Profile paths, machines, selected tool/signal/terminal contracts,
request/response shapes, configuration names/defaults, and closure
membership are compatibility-sensitive. Breaking changes require an
explicit migration note and coordinated consumer update.

scenario-critic and mock are supported test-time catalog members; the existing
conformance vocabulary continues to classify them as supported test-time
library members. They are reusable across agent families and applications,
live under agents/, and carry the same SRD, conformance, portability,
release, and compatibility obligations as other members. Rig-subject and
testdata/conformance REST/control/lifecycle fixtures test agent-core and
remain internal in their current paths; no relocation is scheduled or
tracked for them.

<!-- field: benefits -->
#### Benefits

- Reusable fixes flow upstream to every consumer.
- Application composition and catalog behavior have visible owners.
- Release compatibility and migration evidence are reviewable.
- The runtime remains an interpreter rather than a profile registry.

<!-- decision: 5 -->
### 5 - Tool declarations are first-class without moving core ownership

Reusable repository-specific ToolDef declarations are catalog tool blocks
under their owning agents/<family>/ closure. tools.yaml files only select
names, and profile-local declarations may override configuration for one
profile or variant. Generic Go implementations remain under
agent-core/internal/tools, while stable shared declarations remain under
agent-core/tools and the /opt/agent-core/tools runtime mount.

Applications may select and configure catalog blocks, but application-only
declarations, topology, endpoint choices, wrappers, deployment, and UX
remain in the owning application. MachineSpec owns sequencing; ToolDef
defines one vocabulary word.

<!-- field: benefits -->
#### Benefits

- Tool reuse and ownership are reviewable from repository paths.
- Catalog declarations can evolve without duplicating core implementations.
- Selection, override, declaration, and application composition are not conflated.

<!-- field: technology_choices -->
## Technology choices

| component | technology | purpose |
| --- | --- | --- |
| Profile format | YAML | agent-core loads profile, machine, tool selection, and boundary config files from YAML. |
| Issue tracking | issue tracker | Issues track migration, documentation, validation, and fixture work. |
| Runtime | agent-core | agent-core supplies the executable runtime and reusable tool implementations. |

<!-- field: project_structure | form: "tree" -->
## Project structure

- configuration:
  1. - path: agents/
     - role: Reusable catalog agent blocks, catalog tool declarations, selection lists, overrides, and profile-local assets.
  1. - path: tools/README.md
     - role: Catalog tool-block taxonomy and membership contract.
  1. - path: ../agent-architecture/
     - role: Application-owned Knowledge Manager presentation and composition that references the canonical curator.
  1. - path: testdata/integration/
     - role: Profile-owned fixtures consumed by agent-core integration targets.
- documentation:
  1. - path: docs/VISION.md
     - role: Goals, boundaries, and success criteria.
  1. - path: docs/ARCHITECTURE.md
     - role: Repository ownership, interfaces, components, and decisions.
  1. - path: docs/SPECIFICATIONS.md
     - role: Specification index and coverage notes.
  1. - path: docs/road-map.md
     - role: Migration release plan.
  1. - path: docs/constitutions/issue-format.md
     - role: Issue shape for future work.

<!-- field: implementation_status | form: "tree" -->
## Implementation status

- current_focus: Repository-specific reusable application catalog, validation, release tagging, conformance, and consumer reference contracts.
- progress:
  1. - done: Initial issue queue exists.
  1. - done: Repository purpose and ownership docs are defined.
  1. - done: Agent functional block requirements and semantic models are defined.
  1. - done: Agent profiles and Knowledge Manager UX assets migrated from agent-core into agents/.
  1. - done: Release 06 historically migrated demos and integration fixtures into demo/ and testdata/integration/.
  1. - done: Presentation-specific Knowledge Manager assets were subsequently relocated to applications/agent-architecture; catalog-owned documentation-curator profile UX and integration evidence remain here.
  1. - done: Profile bundle release tagging is implemented through mage tag.
  1. - done: Profile validation and container usage checks verify the mounted profile contract.
  1. - done: Variant profile discovery covers profile.yaml, profile-*.yaml, and *-profile.yaml files.
  1. - done: Formal migration use cases and test suites document the implemented profile repository workflow.
  1. - done: Profile-owned integration tracer bullets cover documentation-curator, bench-to-critic, and monitor/control boundaries. Retired critic-to-executor and planner-to-executor catalog tracers point to application-owned replacements under applications/coding-agent.
  1. - done: Catalog membership, member obligations, one-canonical-home, application reference, exact v0 release compatibility, and migration contracts are recorded.
  1. - done: Tool-first taxonomy separates catalog declaration blocks, selection lists, profile-local overrides, core declarations and implementations, and application-local exclusions.
  1. - done: Release 12 role-realization alignment provides non-forking v0 compatibility wrappers for renamed catalog paths, exact consumer pins, package-closure and generated-Helm evidence, and atomic canonical/legacy release-tag tests; immutable tag publication remains a post-merge action.
  1. - done: scenario-critic and mock are classified as supported test-time catalog members; rig-subject and runtime conformance fixtures remain internal, with no relocation scheduled.
  1. - in-progress: Release 11 scenario critic and mock implementation is present with profile, fixture, conformance, and rig evidence; validator semantics, serving-startup rejection, process survival, named failure causes, and mock-start-failure cleanup remain unproven.

<!-- field: related_documents -->
## Related documents

| doc | purpose |
| --- | --- |
| ../../agent-core/docs/specs/config-formats/runtime-contract.md | Canonical cmd/agent startup, argv flags, and file-only operator path contract. |
| ../../agent-core/docs/specs/software-requirements/srd034-external-agent-profiles.md | External profile ownership and mounted-profile requirements owned by agent-core. |
| docs/VISION.md | Goals, boundaries, and success criteria. |
| docs/SPECIFICATIONS.md | Specification index for future requirements and semantic models. |
| docs/road-map.md | Release plan for profile extraction and mounted runtime usage. |
| docs/constitutions/issue-format.md | Required shape for future issues. |
| docs/specs/software-requirements/srd001-agent-functional-blocks.md | Cross-cutting profile ownership boundary and dependency-direction invariants shared by every profile family. |
| docs/specs/software-requirements/ | Per-family functional block requirements in srd002-executor through srd011-knowledge-manager. |
| docs/specs/semantic-models/multi-agent-profile-system.md | Profile role groups, actor boundaries, and dependency directions. |
| docs/specs/semantic-models/agent-role-taxonomy.md | Per-profile role taxonomy for migration and validation work. |
| ../chatbot-mesh/docs/ARCHITECTURE.md | The deployable chatbot mesh application (a chatbot fanning one embedding out to N Chroma-backed RAG servers, plus a control plane of a coordinator and creator that turn a provisioning intent into an ingest and a rollout), located at applications/chatbot-mesh/ as a standalone, copyable program with its own docs/specs, Helm chart, and ux SPA. Its mesh agents derive from the chatbot and knowledge-manager profile families here. |
