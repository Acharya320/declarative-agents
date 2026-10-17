---
id: agent-profiles-specifications
title: Agent Profiles Specification Index
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

This index starts the specification corpus for profile-owned agent assets.
Foundation documents come first because the migration needs a stable
ownership boundary before assets move.

With the first migration complete, the SRD layer defines migrated profile
families as functional blocks instead of runtime packages. srd001 now holds
the cross-cutting profile ownership boundary and dependency-direction
invariants, while its family_srds inventory names every current shipped
family contract and separates fixture-only REST/control/lifecycle contracts.
The semantic models record role boundaries, dependency directions, and
validation views for every shipped profile.

Architecture decision 4 defines this repository as a versioned reusable-agent
library. Membership requires independent usefulness or multiple consumers,
one canonical home, a family SRD, conformance, portable closure, and
application-safe parameterization. Applications own composition and package
referenced closures against exact canonical repository v0.* releases
(GH-1373). Earlier applications/catalog/v0.* releases and matching
agent-profiles/v0.* names are v0 compatibility aliases,
and immutable historical legacy tags remain valid. scenario-critic and mock are
supported test-time library members under the same release contract.
Rig-subject and REST/control/lifecycle conformance fixtures remain internal in
their current paths; no relocation is scheduled or tracked for those fixtures.

Release 11.0 adds two test-time families over the agent-core release 09.0
words: srd018-scenario-critic is the reusable test rig that composes a scenario
from profile-owned data and reduces validator terminals to a verdict, and
srd019-mock is the one generic mock that stands in for any REST
dependency. Scenarios live beside the agent they test, so one subject carries
several tests and one rig serves both the agent families and the examples
tree.

The pre-#1072 agents/assembler/profile.yaml path remains a non-forking
compatibility wrapper through the remainder of v0 and is removed at
v1. New consumers use
agents/scenario-critic/profile.yaml. Historical Release 11 SRD, use-case, and
test-suite IDs are retained only in explicit historical_aliases fields.

Release 12.0 completes the role-realization naming migration. The assembler,
monitor, and jurist paths are non-forking v0 wrappers over scenario-critic,
runtime-state-reader, and specification-critic. Application-local validator
and chatbot-mesh workload identities migrated atomically. The migration note
pins applications/catalog/v0.20260727.0 and its matching
agent-profiles/v0.20260727.0 alias, while deferring immutable tag publication
until the reviewed commit reaches main.

Release 06.0 now has formal use cases and a test suite for the checked-in
profile repository workflow. Those entries connect profile validation,
mounted container smoke execution, profile bundle tagging, and agent-core
integration consumption back to the SRD and semantic models. The release
records the historical migration into this catalog; presentation-specific
behavior now lives in applications/agent-architecture. The traceability
gap is closed.

Release 07.0 implements profile-owned integration tracer bullets for the
documentation-curator UX, bench-to-critic launch, and monitor/control
boundaries. The former evaluator-to-generator and planner-to-generator
library tracers are retired because their shell children fabricated the
claimed behavior. Their real critic/executor and planner/executor replacement
evidence is application-owned under applications/coding-agent and remains
separate from this release's implemented library targets.

Release 08.0 adds the Chroma corpus agents: an ingest profile and a reader
profile that use a local Chroma vector database as a retrieval boundary. Both
compute embeddings at a local Ollama provider and thread the vector into Chroma
add and query through the agent-core previous-Result parameter capability
(srd028 R12), because Chroma does not embed server-side over its raw HTTP API.
srd012 owns the family requirements. The mesh's integration:chroma and
integration:ragServer targets consume the canonical corpus-ingest program by
reference and prove ingest/count and served retrieval behavior with
application-owned REST parameterization. A deterministic profile-owned
conformance tracer executes the shipped reader through embed, query, grounded
answer, and failure paths against controlled protocol fixtures.

A reader can now start at the SRD, follow the semantic model role entries,
reach the two mounted-profile use cases, and land on named commands such as
mage validate, mage containerSmoke, mage tag, and the agent-core integration
targets that set catalog_root in demo.yaml.

<!-- field: roadmap_summary -->
## Roadmap summary

<!-- roadmap_summary_entry: 00.1 -->
### 00.1

<!-- field: name -->
#### Name

Repository bootstrap

<!-- field: status -->
#### Status

done

<!-- field: scope -->
#### Scope

Replace the template README, add the foundation docs needed for future
structured issues, migrate profile assets, add release tagging, and
validate mounted-profile references.

<!-- roadmap_summary_entry: 06.0 -->
### 06.0

<!-- field: name -->
#### Name

External agent profiles and mounted runtime assets

<!-- field: status -->
#### Status

done

<!-- field: scope -->
#### Scope

Release 06.0 moved profile-owned agent assets, demos, and integration
fixtures into this repository. The catalog still owns the canonical
profiles, profile-local UI, and integration fixtures consumed from
the demo.yaml catalog_root, /profiles, or release bundles. Presentation-specific
assets now live in applications/agent-architecture. Formal use cases
and test suites name executable proof for the completed migration path.

<!-- roadmap_summary_entry: 07.0 -->
### 07.0

<!-- field: name -->
#### Name

Profile-owned integration tracer bullets

<!-- field: status -->
#### Status

done

<!-- field: scope -->
#### Scope

Implemented profile-owned tracer bullets for documentation-curator,
bench-to-critic, and monitor/control boundaries. The retired
critic-to-executor and planner-to-executor library tracers point to
application-owned replacements under applications/coding-agent and are not
counted as implemented Release 07 proof.

<!-- roadmap_summary_entry: 07.1 -->
### 07.1

<!-- field: name -->
#### Name

Per-family SRDs and profile coverage

<!-- field: status -->
#### Status

in-progress

<!-- field: scope -->
#### Scope

Per-family SRDs replace the shared functional-block SRD, and per-family
boundary use cases plus a test suite cover the specification-critic, lifecycle, and rest
families that release 07.0 left to the generic mounted-profile and
profile-validation use cases. Profile-owned Mage integration targets for
the three families remain a deferred follow-up.

<!-- roadmap_summary_entry: 08.0 -->
### 08.0

<!-- field: name -->
#### Name

Chroma corpus agents

<!-- field: status -->
#### Status

done

<!-- field: scope -->
#### Scope

Two profile-owned agents use a local Chroma vector database as a retrieval
boundary: an ingest profile and a reader profile, built as pure
applications/catalog configuration over existing agent-core words. Embeddings are
computed at a local Ollama provider and threaded into Chroma add and query
(srd028 R12), because Chroma does not embed server-side over its raw HTTP
API. The ingest machine lists trusted corpus paths before the model selects
among them, then sequences deterministic embed and add steps. The mesh
integration:chroma target consumes that canonical program by reference and
proves ingest count; deterministic conformance executes the shipped reader
through provider embedding, Chroma query, grounded answer, and failure.

<!-- roadmap_summary_entry: 09.0 -->
### 09.0

<!-- field: name -->
#### Name

RAG server and single-RAG chatbot turn

<!-- field: use_cases_done | form: "scalar" -->
#### Use cases done

2

<!-- field: use_cases_total | form: "scalar" -->
#### Use cases total

2

<!-- field: status -->
#### Status

done

<!-- roadmap_summary_entry: 09.1 -->
### 09.1

<!-- field: name -->
#### Name

Tier selection and multi-RAG fan-out

<!-- field: use_cases_done | form: "scalar" -->
#### Use cases done

1

<!-- field: use_cases_total | form: "scalar" -->
#### Use cases total

1

<!-- field: status -->
#### Status

done

<!-- roadmap_summary_entry: 09.2 -->
### 09.2

<!-- field: name -->
#### Name

Observability

<!-- field: use_cases_done | form: "scalar" -->
#### Use cases done

1

<!-- field: use_cases_total | form: "scalar" -->
#### Use cases total

1

<!-- field: status -->
#### Status

done

<!-- roadmap_summary_entry: 09.3 -->
### 09.3

<!-- field: name -->
#### Name

Mesh deployment and tiered swap

<!-- field: use_cases_done | form: "scalar" -->
#### Use cases done

1

<!-- field: use_cases_total | form: "scalar" -->
#### Use cases total

1

<!-- field: status -->
#### Status

done

<!-- roadmap_summary_entry: 09.4 -->
### 09.4

<!-- field: name -->
#### Name

Control plane

<!-- field: use_cases_done | form: "scalar" -->
#### Use cases done

0

<!-- field: use_cases_total | form: "scalar" -->
#### Use cases total

1

<!-- field: status -->
#### Status

planned

<!-- roadmap_summary_entry: 09.5 -->
### 09.5

<!-- field: name -->
#### Name

In-cluster LLM tier

<!-- field: use_cases_done | form: "scalar" -->
#### Use cases done

1

<!-- field: use_cases_total | form: "scalar" -->
#### Use cases total

1

<!-- field: status -->
#### Status

done

<!-- roadmap_summary_entry: 10.0 -->
### 10.0

<!-- field: name -->
#### Name

Reference-architecture role alignment

<!-- field: use_cases_done | form: "scalar" -->
#### Use cases done

3

<!-- field: use_cases_total | form: "scalar" -->
#### Use cases total

3

<!-- field: status -->
#### Status

done

<!-- roadmap_summary_entry: 11.0 -->
### 11.0

<!-- field: name -->
#### Name

Scenario Critic test rig

<!-- field: status -->
#### Status

in-progress

<!-- field: scope -->
#### Scope

Two test-time role families over the agent-core release 09.0 words. The
scenario-critic (srd018) discovers scenarios under configurable roots covering
the agent families and the applications tree, and for each starts the
scenario's mocks, starts the subject with its REST dependencies
redirected through base_url environment expansion, runs the scenario's
validators against a fresh subject instance, reduces their terminal
states to a verdict naming any failing validator, and tears every child
down. The mock (srd019) is one generic mock serving the canned,
scripted responses a scenario fixture declares and recording a
per-instance request log validators read to verify a subject's outbound
calls. Tests are profile-owned data in tests/<scenario>/ beside the agent,
so one subject carries several tests, adding one never edits the rig, and
no scenario needs a live service or a per-test container image.

The implementation issues are complete and the profiles, fixtures, tests,
and mage target exist. Evidence is still partial: current tests inspect
request logs directly, validate malformed fixtures only with
--validate-config, and infer cleanup from stop spans plus repeatability
without checking process survival or mock-start-failure teardown.

<!-- roadmap_summary_entry: 12.0 -->
### 12.0

<!-- field: name -->
#### Name

Agent role-realization naming alignment

<!-- field: status -->
#### Status

release-ready

<!-- field: scope -->
#### Scope

Implemented migration aliases, exact consumer pins, package closures,
generated Helm inventory checks, release-tag semantics, and scoped
zero-hit evidence. Publication of the immutable root, canonical catalog,
and legacy compatibility tags is a post-merge action from clean main.

<!-- roadmap_summary_entry: 13.0 -->
### 13.0

<!-- field: name -->
#### Name

Collector family promotion

<!-- field: status -->
#### Status

implemented

<!-- field: scope -->
#### Scope

Promoted the chatbot-mesh collector agent to a canonical catalog family
with default spool mode and no relay endpoint default. Added srd020,
conformance coverage exercising agent-core OTLP builtins, and a
bindings.catalog role-realization entry.

<!-- roadmap_summary_entry: 14.0 -->
### 14.0

<!-- field: name -->
#### Name

Served collector trace UI

<!-- field: status -->
#### Status

implemented

<!-- field: scope -->
#### Scope

Serve the shipped collector ui/dist SPA from the collector_query server
via static_assets, same-origin with /query/*, with SPA deep-link
fallback and an explicit error on a missing UI root. The UI root is a
literal profile-declared path, not a new environment variable.

<!-- roadmap_summary_entry: 15.0 -->
### 15.0

<!-- field: name -->
#### Name

Collector Explore panel

<!-- field: status -->
#### Status

implemented

<!-- field: scope -->
#### Scope

Add a Honeycomb-style Explore panel: /query/spans/stats and
/query/spans/breakdown routes bind the agent-core span-analytics tools
(release 13.0), and the Explore UI page renders a duration heatmap,
drag-selects a region, shows the ranked attribute divergence, and links
exemplar spans into the per-trace waterfall.

<!-- roadmap_summary_entry: 16.0 -->
### 16.0

<!-- field: name -->
#### Name

Canonical applier family

<!-- field: use_cases_done | form: "scalar" -->
#### Use cases done

1

<!-- field: use_cases_total | form: "scalar" -->
#### Use cases total

1

<!-- field: status -->
#### Status

done

<!-- roadmap_summary_entry: 17.0 -->
### 17.0

<!-- field: name -->
#### Name

Rig doctor root-cause analyst

<!-- field: status -->
#### Status

implemented

<!-- field: scope -->
#### Scope

Add the rig-doctor family (srd023) on the root-cause-analyst sub-role.
Rig capture gains events, rollout counts, a trace spool tail, and a
manifest; mage diagnose captures a snapshot and runs the rig-doctor
over it. The diagnosis loop reads files only and cites evidence for
every claim.

<!-- roadmap_summary_entry: 18.0 -->
### 18.0

<!-- field: name -->
#### Name

Host-side deployer

<!-- field: status -->
#### Status

implemented

<!-- field: scope -->
#### Scope

Add deploy and undeploy machines to the applier family (srd022 R6) and
run them from the developer's machine. The mage target renders the
declarations with this application's chart, release, namespace, and
kubeconfig, seeds the machine with the values document, and gates on
the terminal state. Day-0 and day-2 differ in the rendered words alone.

<!-- roadmap_summary_entry: 99.0 -->
### 99.0

<!-- field: name -->
#### Name

Observer family promotion

<!-- field: status -->
#### Status

planned

<!-- field: scope -->
#### Scope

Promote the chatbot-mesh fleet observer to a catalog-owned
monitor/fleet-observer family (srd021), parameterized by namespace and
label selector so coding-agent, agent-architecture, and the chatbot-mesh
adopt it by reference. Specification only; the profile, UX, and
agent-role-realizations rows are the adoption implementation, deferred
because the inventory gate requires existing profile paths.

<!-- field: foundation_document_index -->
## Foundation document index

| id | title | summary | path |
| --- | --- | --- | --- |
| [agent-profiles](VISION.md) | Agent Profiles Repository Vision | Goals, boundaries, ownership, and success criteria for external profile assets. | docs/VISION.md |
| architecture-agent-profiles | Agent Profiles Architecture | Repository layout, runtime mount interfaces, ownership decisions, and implementation status. | docs/ARCHITECTURE.md |
| [agent-profiles-roadmap](road-map.md) | Agent Profiles Release Road Map | Bootstrap and migration release plan. | docs/road-map.md |
| [issue-format](constitutions/issue-format.md) | Issue Format Constitution | Required issue fields and scope rules for future work. | docs/constitutions/issue-format.md |
| [migration-v0.20260727.0-agent-role-realization-alignment](migrations/v0.20260727.0-agent-role-realization-alignment.md) | Agent Role-Realization Naming Migration | Old-to-new names and paths, compatibility durations and mechanisms, exact release pins, consumer provenance, and post-merge tag policy. | docs/migrations/v0.20260727.0-agent-role-realization-alignment.md |

<!-- field: srd_index | form: "typed table" -->
## SRD index

| id | title | historical_aliases | summary | status | support | path | canonical_path |
| --- | --- | --- | --- | --- | --- | --- | --- |
| [srd001-agent-functional-blocks](specs/software-requirements/srd001-agent-functional-blocks.md) | Agent Functional Blocks |  | Cross-cutting profile ownership boundary and cross-profile dependency-direction invariants shared by every profile family, with machine-readable shipped-family and fixture-contract inventories. | implemented |  | docs/specs/software-requirements/srd001-agent-functional-blocks.md |  |
| [srd002-executor](specs/software-requirements/srd002-executor.md) | Executor Functional Block |  | Executor coding-agent boundary, inputs, tool selections, runtime dependencies, outputs, and terminal outcomes. | implemented |  | docs/specs/software-requirements/srd002-executor.md |  |
| [srd003-critic](specs/software-requirements/srd003-critic.md) | Critic Functional Block |  | Critic session and point pipeline boundary, inputs, tool selections, runtime dependencies, outputs, and terminal outcomes. | implemented |  | docs/specs/software-requirements/srd003-critic.md |  |
| [srd004-planner](specs/software-requirements/srd004-planner.md) | Planner Functional Block |  | Planner requirement-graph boundary, inputs, tool selections, runtime dependencies, outputs, and terminal outcomes. | implemented |  | docs/specs/software-requirements/srd004-planner.md |  |
| [srd005-specification-critic](specs/software-requirements/srd005-specification-critic.md) | Specification Critic Functional Block | [srd005-jurist] | Specification-critic deterministic validator boundary, inputs, tool selection, responsibilities, provenance-rich findings, and terminal outcomes. | implemented |  | docs/specs/software-requirements/srd005-specification-critic.md |  |
| [srd006-bench](specs/software-requirements/srd006-bench.md) | Bench Functional Block |  | Bench human-operated launcher boundary, inputs, tool selections, runtime dependencies, outputs, and terminal outcomes. | implemented |  | docs/specs/software-requirements/srd006-bench.md |  |
| [srd007-rest](specs/software-requirements/srd007-rest.md) | REST Functional Block |  | REST client and server boundary words, inputs, tool selections, profile-owned authority, outcomes, and terminal outcomes. | implemented |  | docs/specs/software-requirements/srd007-rest.md |  |
| [srd008-runtime-state-reader](specs/software-requirements/srd008-runtime-state-reader.md) | Runtime State Reader Functional Block |  | Runtime-state-reader infrastructure-adapter boundary, cached monitor inputs, tool selections, runtime dependencies, outputs, and terminal outcomes. | implemented |  | docs/specs/software-requirements/srd008-runtime-state-reader.md |  |
| [srd009-lifecycle](specs/software-requirements/srd009-lifecycle.md) | Lifecycle Functional Blocks |  | Lifecycle history, rollback, and approval machine-action boundaries, inputs, tool selections, runtime dependencies, outputs, and terminal outcomes. | implemented |  | docs/specs/software-requirements/srd009-lifecycle.md |  |
| [srd010-control](specs/software-requirements/srd010-control.md) | Control Functional Block |  | Control signal-intake boundary, inputs, tool selections, runtime dependencies, outputs, and terminal outcomes. | implemented |  | docs/specs/software-requirements/srd010-control.md |  |
| [srd011-knowledge-manager](specs/software-requirements/srd011-knowledge-manager.md) | Knowledge Manager Functional Block |  | Documentation-curator UX and REST workflow boundary, inputs, tool selections, runtime dependencies, outputs, and terminal outcomes. | implemented |  | docs/specs/software-requirements/srd011-knowledge-manager.md |  |
| [srd012-chroma-corpus-agents](specs/software-requirements/srd012-chroma-corpus-agents.md) | Chroma Corpus Agents Functional Block |  | Ingest and reader Chroma retrieval boundary, REST vocabulary, provider-embedding threading (srd028 R12), inputs, tool selections, server lifecycle precondition, and terminal outcomes. | implemented |  | docs/specs/software-requirements/srd012-chroma-corpus-agents.md |  |
| [srd013-rag-server-agent](specs/software-requirements/srd013-rag-server-agent.md) | RAG Server Agent Functional Block (Relocated) |  | Pointer-only historical identifier; canonical RAG server requirements and evidence live with applications/chatbot-mesh. | relocated |  | docs/specs/software-requirements/srd013-rag-server-agent.md | applications/chatbot-mesh/docs/specs/software-requirements/srd001-rag-server-agent.md |
| [srd014-chatbot-agent](specs/software-requirements/srd014-chatbot-agent.md) | Chatbot Agent Functional Block (Relocated) |  | Pointer-only historical identifier; canonical chatbot requirements and evidence live with applications/chatbot-mesh. | relocated |  | docs/specs/software-requirements/srd014-chatbot-agent.md | applications/chatbot-mesh/docs/specs/software-requirements/srd002-chatbot-agent.md |
| [srd015-chatbot-deployment](specs/software-requirements/srd015-chatbot-deployment.md) | Chatbot Mesh Deployment Model (Relocated) |  | Pointer-only historical identifier; canonical deployment requirements and evidence live with applications/chatbot-mesh. | relocated |  | docs/specs/software-requirements/srd015-chatbot-deployment.md | applications/chatbot-mesh/docs/specs/software-requirements/srd003-chatbot-deployment.md |
| [srd016-coordinator](specs/software-requirements/srd016-coordinator.md) | Coordinator Agent Functional Block (Relocated) |  | Pointer-only historical identifier; canonical coordinator requirements and evidence live with applications/chatbot-mesh. | relocated |  | docs/specs/software-requirements/srd016-coordinator.md | applications/chatbot-mesh/docs/specs/software-requirements/srd004-provisioning-workflow-orchestrator.md |
| [srd017-creator](specs/software-requirements/srd017-creator.md) | Creator Agent Functional Block (Relocated) |  | Pointer-only historical identifier; canonical creator requirements and evidence live with applications/chatbot-mesh. | relocated |  | docs/specs/software-requirements/srd017-creator.md | applications/chatbot-mesh/docs/specs/software-requirements/srd005-creator.md |
| [srd018-scenario-critic](specs/software-requirements/srd018-scenario-critic.md) | Scenario Critic Functional Block | [srd018-assembler] | The repository's Critic-bound test rig. It composes scenarios as an actor function without claiming the canonical Assembler responsibility. It discovers scenarios under configurable roots covering the agent families and the examples tree, and for each starts the scenario's mocks, starts the subject with its REST dependencies redirected through base_url environment expansion, runs the scenario's validators against a fresh subject instance, reduces their terminal states to a verdict that names a failing validator, and tears every child down. Scenarios are profile-owned data in tests/<scenario>/ beside the agent, so one subject carries several tests and adding one never edits the rig. | in-progress | supported-test-time | docs/specs/software-requirements/srd018-scenario-critic.md |  |
| [srd019-mock](specs/software-requirements/srd019-mock.md) | Mock Functional Block |  | One generic mock profile standing in for any REST dependency a subject calls. It serves the canned, scripted responses a scenario fixture declares through the agent-core mock endpoint binding, selects that fixture by environment so a new dependency needs a fixture rather than a profile or an image, exposes a health route, and records a per-instance request log validators read to verify a subject's outbound calls. | in-progress | supported-test-time | docs/specs/software-requirements/srd019-mock.md |  |
| [srd020-collector](specs/software-requirements/srd020-collector.md) | Collector Agent Family |  | Canonical OTLP trace collector promoted from the chatbot-mesh, defaulting to spool mode with no relay endpoint. Exercises agent-core OTLP builtins (receive, spool, relay) in a reusable catalog family, serves its own trace UI same-origin with the query surface (R7), and exposes an Explore panel over span analytics (R8-R9). | implemented |  | docs/specs/software-requirements/srd020-collector.md |  |
| [srd021-observer](specs/software-requirements/srd021-observer.md) | Observer Agent Family |  | Generalized fleet observer promoted from the chatbot-mesh (srd008) to a catalog-owned monitor/fleet-observer family, parameterized by namespace and label selector for per-application adoption. Specification only; the profile, UX, and realization rows are the adoption implementation. | planned |  | docs/specs/software-requirements/srd021-observer.md |  |
| [srd022-applier](specs/software-requirements/srd022-applier.md) | Canonical Deployment Applier Family |  | Shared lifecycle host, request-machine sequencing, name-only selections, and ToolDefs used by three application-owned wrappers whose endpoints, charts, rollout commands, and authority remain local, plus the host-side deploy and undeploy machines that install a release from the developer's machine through the same declared sequence. | implemented |  | docs/specs/software-requirements/srd022-applier.md |  |
| [srd023-rig-doctor](specs/software-requirements/srd023-rig-doctor.md) | Rig Doctor Root-Cause Analyst |  | Catalog root-cause-analyst that reads one captured kind evidence directory and the trace spool tail, never the live cluster, and writes a diagnosis whose every claim cites a captured file and line. | implemented |  | docs/specs/software-requirements/srd023-rig-doctor.md |  |

<!-- field: config_format_index -->
## Config format index

<!-- field: semantic_model_index -->
## Semantic model index

<!-- semantic_model_index_entry: sm-agent-role-realizations -->
### sm-agent-role-realizations

<!-- field: title -->
#### Title

Agent Role Realizations and Bindings

<!-- field: summary -->
#### Summary

Shared authority for canonical role IDs, concrete actor/profile bindings,
naming grammar, classifications, inheritance, and migration aliases.

<!-- field: status -->
#### Status

implemented

<!-- field: path -->
#### Path

../docs/specs/semantic-models/agent-role-realizations.md

<!-- field: reciprocal_index -->
#### Reciprocal index

../docs/SPECIFICATIONS.md#semantic_model_index

<!-- semantic_model_index_entry: sm-multi-agent-profile-system -->
### sm-multi-agent-profile-system

<!-- field: title -->
#### Title

Multi-Agent Profile System

<!-- field: summary -->
#### Summary

Role groups, actor boundaries, dependency directions, and validation views for profile-owned agent programs.

<!-- field: status -->
#### Status

implemented

<!-- field: path -->
#### Path

docs/specs/semantic-models/multi-agent-profile-system.md

<!-- semantic_model_index_entry: sm-agent-role-taxonomy -->
### sm-agent-role-taxonomy

<!-- field: title -->
#### Title

Agent Role Taxonomy

<!-- field: summary -->
#### Summary

Per-profile role, inputs, outputs, selected tool families, runtime dependencies, outcomes, and checks.

<!-- field: status -->
#### Status

implemented

<!-- field: path -->
#### Path

docs/specs/semantic-models/agent-role-taxonomy.md

<!-- field: use_case_index_notes -->
## Use case index notes

Release 09 is staged one release per epic (GH-332): 09.0 core loop, 09.1
routing/fan-out, 09.2 observability, 09.3 deployment, 09.4 control plane, plus
the 10.0 role-alignment release. The rel09 use-case and test-suite files were
renumbered to their staged ids and the observability (09.2) and alignment
(10.0) tracers authored by GH-333; the 09.4 control-plane tracers land with
srd016/srd017 (GH-329).

Relocated (GH-418): the mesh data-plane specs (srd013/014/015 and the rel09.0,
09.1, 09.2, 09.3, and 09.5 use cases and test suites) now live with the
standalone example at applications/chatbot-mesh/docs/specs, renumbered to
srd001-003 and rel00.0-04.0. The entries below remain for history; the
example is the canonical copy.

Relocated (GH-430): the 09.4 control plane (srd016/srd017 and the rel09.4 use
case and test suite) now lives with the example too, renumbered to
srd004/srd005 and rel05.0. The entries below remain for history; the example
is canonical.

<!-- field: use_case_index | form: "typed table" -->
## Use case index

| id | title | historical_aliases | summary | status | path | canonical_path |
| --- | --- | --- | --- | --- | --- | --- |
| [rel06.0-uc001-mounted-profile-execution](specs/use-cases/rel06.0-uc001-mounted-profile-execution.md) | Mounted Profile Execution |  | Mounted execution path for running agent-core with canonical profiles, profile-local assets, fixtures, and workspaces supplied from applications/catalog; applications own presentation composition. | implemented | docs/specs/use-cases/rel06.0-uc001-mounted-profile-execution.md |  |
| [rel06.0-uc002-profile-validation](specs/use-cases/rel06.0-uc002-profile-validation.md) | Profile Validation and Release Bundle Readiness |  | Validation path for profile reference checks, mounted container smoke execution, and profile bundle tag readiness. | implemented | docs/specs/use-cases/rel06.0-uc002-profile-validation.md |  |
| [rel07.0-uc001-documentation-curator-profile-ux](specs/use-cases/rel07.0-uc001-documentation-curator-profile-ux.md) | Documentation Curator Profile UX Tracer |  | Profile-owned Knowledge Manager browser, machine_request, trace, and lifecycle proof. | implemented | docs/specs/use-cases/rel07.0-uc001-documentation-curator-profile-ux.md |  |
| [rel07.0-uc002-evaluator-generator-profile-boundary](specs/use-cases/rel07.0-uc002-evaluator-generator-profile-boundary.md) | Evaluator to Generator Profile Boundary |  | Retired library tracer; applications/coding-agent owns the real critic-to-executor replacement target. | planned | docs/specs/use-cases/rel07.0-uc002-evaluator-generator-profile-boundary.md |  |
| [rel07.0-uc003-planner-generator-profile-boundary](specs/use-cases/rel07.0-uc003-planner-generator-profile-boundary.md) | Planner to Generator Profile Boundary |  | Retired library tracer; applications/coding-agent owns the real planner-to-executor replacement target. | planned | docs/specs/use-cases/rel07.0-uc003-planner-generator-profile-boundary.md |  |
| [rel07.0-uc004-bench-evaluator-profile-boundary](specs/use-cases/rel07.0-uc004-bench-evaluator-profile-boundary.md) | Bench to Critic Profile Boundary |  | Bench launches the critic role from a human/operator boundary action; the historical use-case id retains the former evaluator name. | implemented | docs/specs/use-cases/rel07.0-uc004-bench-evaluator-profile-boundary.md |  |
| [rel07.0-uc005-monitor-control-profile-boundary](specs/use-cases/rel07.0-uc005-monitor-control-profile-boundary.md) | Runtime State Reader and Control Profile Boundary |  | Monitor/control profile boundary and lifecycle tracer evidence. | implemented | docs/specs/use-cases/rel07.0-uc005-monitor-control-profile-boundary.md |  |
| [rel07.1-uc001-specification-critic-profile-boundary](specs/use-cases/rel07.1-uc001-specification-critic-profile-boundary.md) | Specification Critic Profile Boundary | [rel07.1-uc001-jurist-profile-boundary] | Deterministic, LLM-free charter validation from load through report to a Passed or Failed terminal state. | specified | docs/specs/use-cases/rel07.1-uc001-specification-critic-profile-boundary.md |  |
| [rel07.1-uc002-lifecycle-profile-boundary](specs/use-cases/rel07.1-uc002-lifecycle-profile-boundary.md) | Lifecycle Profile Boundary |  | History, rollback, and approval as machine actions with explicit Approved or Rejected resume signals. | specified | docs/specs/use-cases/rel07.1-uc002-lifecycle-profile-boundary.md |  |
| [rel07.1-uc003-rest-profile-boundary](specs/use-cases/rel07.1-uc003-rest-profile-boundary.md) | REST Profile Boundary |  | Configured client and server boundary words, profile-owned REST authority, OpenAPI imports, and the ollama-* variant. | specified | docs/specs/use-cases/rel07.1-uc003-rest-profile-boundary.md |  |
| [rel08.0-uc001-corpus-ingest](specs/use-cases/rel08.0-uc001-corpus-ingest.md) | Corpus Ingest into Chroma |  | Ingest profile reads a corpus, embeds each document at Ollama, and threads the vector and carried document and id into Chroma add, then verifies the collection count. | implemented | docs/specs/use-cases/rel08.0-uc001-corpus-ingest.md |  |
| [rel08.0-uc002-corpus-query](specs/use-cases/rel08.0-uc002-corpus-query.md) | Corpus Query from Chroma |  | Reader profile embeds a question at Ollama, threads the vector into Chroma query, and answers grounded in the retrieved chunks. | implemented | docs/specs/use-cases/rel08.0-uc002-corpus-query.md |  |
| [rel09.0-uc001-rag-server-query](specs/use-cases/rel09.0-uc001-rag-server-query.md) | Persistent RAG Server Answers a Vector-In Query (Relocated) |  | Pointer-only historical identifier; the canonical use case and traces live with applications/chatbot-mesh release 00.0. | relocated | docs/specs/use-cases/rel09.0-uc001-rag-server-query.md | applications/chatbot-mesh/docs/specs/use-cases/rel00.0-uc001-rag-server-query.md |
| [rel09.0-uc002-chatbot-turn](specs/use-cases/rel09.0-uc002-chatbot-turn.md) | Chatbot Turn Embeds, Retrieves From One RAG, and Answers (Relocated) |  | Pointer-only historical identifier; the canonical use case and traces live with applications/chatbot-mesh release 00.0. | relocated | docs/specs/use-cases/rel09.0-uc002-chatbot-turn.md | applications/chatbot-mesh/docs/specs/use-cases/rel00.0-uc002-chatbot-turn.md |
| [rel09.1-uc001-routed-fanout](specs/use-cases/rel09.1-uc001-routed-fanout.md) | Chatbot Turn Selects a Model Tier and Fans Out Over Two RAGs (Relocated) |  | Pointer-only historical identifier; the canonical use case and traces live with applications/chatbot-mesh release 01.0. | relocated | docs/specs/use-cases/rel09.1-uc001-routed-fanout.md | applications/chatbot-mesh/docs/specs/use-cases/rel01.0-uc001-routed-fanout.md |
| [rel09.2-uc001-observability](specs/use-cases/rel09.2-uc001-observability.md) | Observability Panel and Cross-Agent Trace Waterfall (Relocated) |  | Pointer-only historical identifier; the canonical use case and traces live with applications/chatbot-mesh release 02.0. | relocated | docs/specs/use-cases/rel09.2-uc001-observability.md | applications/chatbot-mesh/docs/specs/use-cases/rel02.0-uc001-observability.md |
| [rel09.3-uc001-mesh-deployment](specs/use-cases/rel09.3-uc001-mesh-deployment.md) | Mesh Deployment and Tiered Swap (Relocated) |  | Pointer-only historical identifier; the canonical use case and traces live with applications/chatbot-mesh release 03.0. | relocated | docs/specs/use-cases/rel09.3-uc001-mesh-deployment.md | applications/chatbot-mesh/docs/specs/use-cases/rel03.0-uc001-mesh-deployment.md |
| [rel09.4-uc001-control-plane-provisioning](specs/use-cases/rel09.4-uc001-control-plane-provisioning.md) | Control-Plane Provisioning (Relocated) |  | Pointer-only historical identifier; the canonical use case and traces live with applications/chatbot-mesh release 05.0. | relocated | docs/specs/use-cases/rel09.4-uc001-control-plane-provisioning.md | applications/chatbot-mesh/docs/specs/use-cases/rel05.0-uc001-control-plane-provisioning.md |
| [rel09.5-uc001-in-cluster-llm-tier](specs/use-cases/rel09.5-uc001-in-cluster-llm-tier.md) | In-Cluster LLM Tier (Relocated) |  | Pointer-only historical identifier; the canonical use case and traces live with applications/chatbot-mesh release 04.0. | relocated | docs/specs/use-cases/rel09.5-uc001-in-cluster-llm-tier.md | applications/chatbot-mesh/docs/specs/use-cases/rel04.0-uc001-in-cluster-llm-tier.md |
| [rel10.0-uc001-role-catalog-renames](specs/use-cases/rel10.0-uc001-role-catalog-renames.md) | Evaluator and Generator Families Renamed to Critic and Executor |  | The evaluator family is renamed to critic and generator to executor (reference role names) across agents/, the SRD ids, the taxonomy, and the agent-core builtin tool dirs, with profile wiring and conformance intact. | implemented | docs/specs/use-cases/rel10.0-uc001-role-catalog-renames.md |  |
| [rel10.0-uc002-knowledge-family-rescope](specs/use-cases/rel10.0-uc002-knowledge-family-rescope.md) | Chroma Profiles Rescoped Under Knowledge Management as Corpus Ingest and Reader |  | The vendor-named chroma family moves under knowledge-manager as corpus-ingest and corpus-reader; Integration.Chroma passes from the new paths and agents/chroma is removed. | implemented | docs/specs/use-cases/rel10.0-uc002-knowledge-family-rescope.md |  |
| [rel10.0-uc003-conformance-fixture-relocation](specs/use-cases/rel10.0-uc003-conformance-fixture-relocation.md) | Conformance Grammar Fixtures Relocated to testdata/conformance |  | The rest, control, and lifecycle conformance fixtures move out of agents/ to testdata/conformance so agents/ lists only role families; conformance and Integration.MonitorControl pass from the new paths. | implemented | docs/specs/use-cases/rel10.0-uc003-conformance-fixture-relocation.md |  |
| [rel11.0-uc001-scenario-critic-scenario-run](specs/use-cases/rel11.0-uc001-scenario-critic-scenario-run.md) | The Scenario Critic Runs an Agent's Scenarios to a Verdict Without Live Services | [rel11.0-uc001-assembler-scenario-run] | The scenario critic discovers scenarios beside an agent, starts mocks and a redirected subject, dispatches validators, reports ordered pass/fail verdicts, exposes child-stop spans, and repeats without observed port or state leakage. Named validator causes and process non-survival are not yet asserted. | in-progress | docs/specs/use-cases/rel11.0-uc001-scenario-critic-scenario-run.md |  |
| [rel11.0-uc002-mock-serves-fixtures](specs/use-cases/rel11.0-uc002-mock-serves-fixtures.md) | One Mock Profile Stands In for Any REST Dependency and Records What It Received |  | One shipped mock serves fixture routes, scripts response sequences, records unmatched calls, and exposes a per-instance request log inspected directly by conformance tests; malformed fixtures are rejected by --validate-config, while serving-startup rejection remains unproven. | in-progress | docs/specs/use-cases/rel11.0-uc002-mock-serves-fixtures.md |  |
| [rel12.0-uc001-agent-role-realization-alignment](specs/use-cases/rel12.0-uc001-agent-role-realization-alignment.md) | Migrate Concrete Agent Names as One Versioned Compatibility Release |  | Coordinate migration aliases, exact catalog pins, package closures, tag semantics, generated artifacts, and scoped zero-hit evidence. | implemented | docs/specs/use-cases/rel12.0-uc001-agent-role-realization-alignment.md |  |
| [rel13.0-uc001-collector-family](specs/use-cases/rel13.0-uc001-collector-family.md) | A standalone collector receives, spools, and shuts down in spool mode |  | Canonical catalog collector launches in default spool mode, spools received traces, skips relay, and shuts down cleanly on exit signal. | implemented | docs/specs/use-cases/rel13.0-uc001-collector-family.md |  |
| [rel14.0-uc001-served-trace-ui](specs/use-cases/rel14.0-uc001-served-trace-ui.md) | The collector serves its own trace UI at the query origin |  | The collector_query server serves the shipped ui/dist SPA same-origin with /query/*, resolving deep links to the index and leaving query responses unchanged. | implemented | docs/specs/use-cases/rel14.0-uc001-served-trace-ui.md |  |
| [rel15.0-uc001-explore-panel](specs/use-cases/rel15.0-uc001-explore-panel.md) | An operator narrows a latency band to a responsible attribute and lands on a trace |  | The Explore page renders a duration heatmap from /query/spans/stats, drag-selects a region, shows the ranked attribute divergence from /query/spans/breakdown, and links exemplar spans into the waterfall. | implemented | docs/specs/use-cases/rel15.0-uc001-explore-panel.md |  |
| [rel16.0-uc001-applier-family](specs/use-cases/rel16.0-uc001-applier-family.md) | Three applications compose one canonical applier family |  | Catalog-owned lifecycle and request machines are packaged with thin application wrappers that retain local deployment authority. | implemented | docs/specs/use-cases/rel16.0-uc001-applier-family.md |  |
| [rel17.0-uc001-rig-doctor](specs/use-cases/rel17.0-uc001-rig-doctor.md) | An operator runs mage diagnose after a failed scenario and reads a cited diagnosis |  | The diagnose target captures a scenario snapshot with a manifest and runs the rig-doctor over it; the operator reads a cited diagnosis and runs the suggested read-only command. | implemented | docs/specs/use-cases/rel17.0-uc001-rig-doctor.md |  |
| [rel18.0-uc001-host-side-deploy](specs/use-cases/rel18.0-uc001-host-side-deploy.md) | An operator runs mage deploy and the applier's deploy machine installs the release |  | The deploy target renders the applier's declarations with this application's coordinates, seeds the deploy machine with the values document, and gates the build on the terminal state. | implemented | docs/specs/use-cases/rel18.0-uc001-host-side-deploy.md |  |
| [rel99.0-uc001-observer-family](specs/use-cases/rel99.0-uc001-observer-family.md) | The observer is a catalog-owned family adopted across applications |  | The fleet observer becomes a catalog-owned family generalized from the chatbot-mesh observer, parameterized per application; coding-agent, agent-architecture, and the chatbot-mesh adopt it by reference. | planned | docs/specs/use-cases/rel99.0-uc001-observer-family.md |  |

<!-- field: test_suite_index | form: "typed table" -->
## Test suite index

| id | title | historical_aliases | summary | status | path | canonical_path |
| --- | --- | --- | --- | --- | --- | --- |
| [test-rel06.0-agent-profiles](specs/test-suites/test-rel06.0-agent-profiles.md) | Test Suite for Release 06.0 Agent Profiles |  | Executable proof index for mage validate, mage containerSmoke, mage tag, and agent-core integration targets that consume the demo.yaml catalog_root. | implemented | docs/specs/test-suites/test-rel06.0-agent-profiles.md |  |
| [test-rel07.0-profile-integrations](specs/test-suites/test-rel07.0-profile-integrations.md) | Test Suite for Release 07.0 Profile Integrations |  | Implemented profile-owned tracer bullets for documentation-curator UX, bench-to-critic, and monitor/control boundaries, plus retired critic-to-executor and planner-to-executor entries that point to application-owned replacement targets. | implemented | docs/specs/test-suites/test-rel07.0-profile-integrations.md |  |
| [test-rel07.1-profile-boundaries](specs/test-suites/test-rel07.1-profile-boundaries.md) | Test Suite for Release 07.1 Profile Boundaries |  | Boundary use-case validation for the specification-critic, lifecycle, and rest families against their per-family SRD acceptance criteria. | specified | docs/specs/test-suites/test-rel07.1-profile-boundaries.md |  |
| [test-rel07.1-profile-conformance](specs/test-suites/test-rel07.1-profile-conformance.md) | Test Suite for Release 07.1 Profile Conformance |  | Executable per-family OpenTelemetry conformance for the ten Release 07 contracts, traced to srd002-srd011 and sm-agent-role-taxonomy; integration target mage conformance. Nine families run the shipped profile.yaml wrapper or a named shipped variant (copy-and-patch only harness-controlled fields), so their coverage is shipped-behavioral rather than synthesized; the rest sample server boundary is bounded-behavioral (shipped REST config and declarations over a synthesized launch/await/stop machine, with the client probe/delay grammar covered by shipped wiring). The executor and rest families also run their shipped variants and carry an ungated load-and-wiring companion. Executor (with qwen27b/qwen35b variants), the planner behavioral run, and the rest ollama-profile variant are Ollama-gated; lifecycle suspend/resume and history/rollback are Dolt-gated; rest proves its inbound webhook server boundary and explicitly declares asynchronous readiness probe/delay polling in MachineSpec. | implemented | docs/specs/test-suites/test-rel07.1-profile-conformance.md |  |
| [test-rel08.0-corpus-agents](specs/test-suites/test-rel08.0-corpus-agents.md) | Test Suite for Release 08.0 Chroma Corpus Agents |  | Deterministic profile-owned evidence for corpus-ingest trust ordering and the shipped corpus-reader embed, query, grounded-answer, and readiness-failure paths. | implemented | docs/specs/test-suites/test-rel08.0-corpus-agents.md |  |
| [test-rel09.0-chatbot-mesh](specs/test-suites/test-rel09.0-chatbot-mesh.md) | Test Suite for Release 09.0 RAG Server and Chatbot Turn (Relocated) |  | Pointer-only historical identifier; the canonical suite and proof commands live with applications/chatbot-mesh release 00.0. | relocated | docs/specs/test-suites/test-rel09.0-chatbot-mesh.md | applications/chatbot-mesh/docs/specs/test-suites/test-rel00.0-chatbot-mesh.md |
| [test-rel09.1-routing-fanout](specs/test-suites/test-rel09.1-routing-fanout.md) | Test Suite for Release 09.1 Tier Selection and Fan-Out (Relocated) |  | Pointer-only historical identifier; the canonical suite and proof commands live with applications/chatbot-mesh release 01.0. | relocated | docs/specs/test-suites/test-rel09.1-routing-fanout.md | applications/chatbot-mesh/docs/specs/test-suites/test-rel01.0-routing-fanout.md |
| [test-rel09.2-observability](specs/test-suites/test-rel09.2-observability.md) | Test Suite for Release 09.2 Observability (Relocated) |  | Pointer-only historical identifier; the canonical suite and proof commands live with applications/chatbot-mesh release 02.0. | relocated | docs/specs/test-suites/test-rel09.2-observability.md | applications/chatbot-mesh/docs/specs/test-suites/test-rel02.0-observability.md |
| [test-rel09.3-mesh-deployment](specs/test-suites/test-rel09.3-mesh-deployment.md) | Test Suite for Release 09.3 Mesh Deployment and Tiered Swap (Relocated) |  | Pointer-only historical identifier; the canonical suite and proof commands live with applications/chatbot-mesh release 03.0. | relocated | docs/specs/test-suites/test-rel09.3-mesh-deployment.md | applications/chatbot-mesh/docs/specs/test-suites/test-rel03.0-mesh-deployment.md |
| [test-rel09.4-control-plane](specs/test-suites/test-rel09.4-control-plane.md) | Test Suite for Release 09.4 Control Plane (Relocated) |  | Pointer-only historical identifier; the canonical suite and proof commands live with applications/chatbot-mesh release 05.0. | relocated | docs/specs/test-suites/test-rel09.4-control-plane.md | applications/chatbot-mesh/docs/specs/test-suites/test-rel05.0-control-plane.md |
| [test-rel09.5-llm-tier](specs/test-suites/test-rel09.5-llm-tier.md) | Test Suite for Release 09.5 In-Cluster LLM Tier (Relocated) |  | Pointer-only historical identifier; the canonical suite and proof commands live with applications/chatbot-mesh release 04.0. | relocated | docs/specs/test-suites/test-rel09.5-llm-tier.md | applications/chatbot-mesh/docs/specs/test-suites/test-rel04.0-llm-tier.md |
| [test-rel10.0-role-catalog](specs/test-suites/test-rel10.0-role-catalog.md) | Test Suite for Release 10.0 Reference-Architecture Role Alignment |  | Implemented cases validate the critic/executor renames, the Chroma rescope under knowledge-manager, and the conformance-fixture relocation with mage validate, mage conformance, go test, and zero-hit greps. | implemented | docs/specs/test-suites/test-rel10.0-role-catalog.md |  |
| [test-rel11.0-scenario-critic-rig](specs/test-suites/test-rel11.0-scenario-critic-rig.md) | Test Suite for Release 11.0 Scenario Critic Test Rig | [test-rel11.0-assembler-rig] | Implemented cases cover fixture serving and sequencing, two mocks from one profile, direct request-log inspection, malformed-fixture --validate-config rejection, ordered scenario verdicts, visible child-stop dispatches, repeatable runs, cross-root discovery, and shipped-profile validation. Validator pass/fail assertions, serving-startup rejection, process survival, named causes, and mock-start-failure cleanup remain evidence gaps. | in-progress | docs/specs/test-suites/test-rel11.0-scenario-critic-rig.md |  |
| [test-rel12.0-agent-role-realization-alignment](specs/test-suites/test-rel12.0-agent-role-realization-alignment.md) | Test Suite for Release 12.0 Agent Role-Realization Alignment |  | Executable aliases, exact release metadata, package closures, Helm render/drift, tag semantics, and zero-hit evidence. | implemented | docs/specs/test-suites/test-rel12.0-agent-role-realization-alignment.md |  |
| [test-rel13.0-collector](specs/test-suites/test-rel13.0-collector.md) | Test Suite for Release 13.0 Collector Family |  | Conformance test launching the canonical collector in spool mode and verifying clean shutdown with expected tool spans. | implemented | docs/specs/test-suites/test-rel13.0-collector.md |  |
| [test-rel14.0-served-ui](specs/test-suites/test-rel14.0-served-ui.md) | Test Suite for Release 14.0 Served Collector Trace UI |  | Conformance cases for the served SPA index at root, deep-link fallback, query API unaffected by the catch-all, and an explicit error on a missing UI root. | implemented | docs/specs/test-suites/test-rel14.0-served-ui.md |  |
| [test-rel15.0-explore](specs/test-suites/test-rel15.0-explore.md) | Test Suite for Release 15.0 Collector Explore Panel |  | Conformance cases for the /query/spans/stats and /query/spans/breakdown response contracts, config-driven caps, empty-spool responses, the Explore UI route declaration, and uiDist reproducibility. | implemented | docs/specs/test-suites/test-rel15.0-explore.md |  |
| [test-rel16.0-applier](specs/test-suites/test-rel16.0-applier.md) | Canonical Applier Family Conformance |  | Boot, machine, wrapper, duplicate, package, application, and statistics evidence for the promoted applier family. | implemented | docs/specs/test-suites/test-rel16.0-applier.md |  |
| [test-rel17.0-rig-doctor](specs/test-suites/test-rel17.0-rig-doctor.md) | Rig Doctor Root-Cause Analyst Conformance |  | Closure, mock-provider diagnosis, empty-evidence, capture-manifest, registration, and dependency-gated live cases for the rig-doctor. | implemented | docs/specs/test-suites/test-rel17.0-rig-doctor.md |  |
| [test-rel18.0-host-side-deploy](specs/test-suites/test-rel18.0-host-side-deploy.md) | Host-Side Deployer Conformance |  | Closure, terminal-state, compensation, rendering-guard, and build-gating cases for the deploy and undeploy machines, with one live case for the atomic install contract. | implemented | docs/specs/test-suites/test-rel18.0-host-side-deploy.md |  |
| [test-rel99.0-observer](specs/test-suites/test-rel99.0-observer.md) | Test Suite for Release 99.0 Observer Family Promotion |  | Specification-graph cases for the generalized observer-family SRD, the carried capabilities, and the named adopters with deferred realization rows. | planned | docs/specs/test-suites/test-rel99.0-observer.md |  |

<!-- field: coverage_gaps -->
## Coverage gaps

The SRD and semantic models describe ownership, role boundaries, and
functional-block requirements for migrated profile families. Profile assets,
profile-local UI, integration fixtures, release tagging, profile validation,
and mounted-container smoke checks are implemented. Application presentation
behavior is owned and evidenced by applications/agent-architecture.

Release 06.0 traceability now covers mounted profile execution and profile
validation. It also covers mounted container smoke, profile bundle tagging,
and agent-core integration targets that consume this repository through
the demo.yaml catalog_root.
No missing proof remains for this release.

Future catalog gaps should track new profile families, profile-owned fixtures,
or validation commands as they are added. Application demo gaps belong to the
owning application. Existing release coverage no longer depends on a deferred
follow-up issue.

Release 07.0 implements the documentation-curator profile UX tracer,
bench-to-critic launch tracer, and monitor/control boundary tracer. The
evaluator-to-generator and planner-to-generator library tracers are retired:
their shell children fabricated the claimed workspace and trace evidence.
applications/coding-agent owns the real critic-to-executor and
planner-to-executor replacements, which run the canonical profiles through
the built agent binary and remain outside this release's implemented library
proof. The monitor/control target records profile-owned configuration
evidence while agent-core retains runtime ownership for REST and lifecycle
builtin behavior.

Release 07.1 adds per-family boundary use cases for the specification-critic, lifecycle, and
rest families, which release 07.0 covered only through the generic rel06.0
mounted-profile and profile-validation use cases. Every Release 07 family has
a dedicated boundary use case: critic, executor, planner, bench,
monitor/control, and documentation-curator through release 07.0, and specification-critic,
lifecycle, and rest through release 07.1. test-rel07.1-profile-boundaries
validates the three new use cases against their per-family SRD acceptance
criteria. The remaining gap is profile-owned Mage integration targets that
execute the specification-critic, lifecycle, and rest boundary fixtures; until those land the
three use cases and their test cases stay status specified rather than
implemented.

Release 07.1 also adds test-rel07.1-profile-conformance, an executable
per-family conformance suite run by mage conformance. It drives each of the ten
families through the agent-core CLI and asserts on the captured OpenTelemetry
trace (span vocabulary and terminal state), tracing every Release 07 family to its
srd002-srd011 acceptance criteria and sm-agent-role-taxonomy role. Coverage is
shipped-behavioral, not synthesized, for nine families: each runs the shipped
profile.yaml wrapper (or a named shipped variant), copying it only to patch
harness-controlled fields — a listen address/port, a corpus directory, or the
iteration budget — never rebuilding the machine, tools, or wrapper. The rest
sample server boundary is the one bounded-behavioral case: it runs the shipped
REST config and declarations but a synthesized bounded launch/await/stop
machine, because the shipped machine's client head cannot run in a profile
machine (that head is proven at the agent-core package level).
Representative-vs-variant is
recorded explicitly: the executor ships and runs profile.yaml plus the
profile-qwen35b and profile-qwen27b variants, and the rest family runs the
sample inbound webhook boundary plus the ollama-profile variant; the executor,
planner, and rest families each also carry an ungated, model-free
load-and-wiring companion that asserts shipped-profile load and machine wiring
where a gate would otherwise skip the run. All ten Release 07 cases are implemented
against the discovered sibling agent-core checkout; executor and its qwen variants, the planner behavioral
run, and the rest ollama-profile variant are Ollama-gated, and lifecycle
suspend/resume and history/rollback are Dolt-gated. Planner conformance landed
once agent-core gained a load_graph pipeline word that seeds the requirement
graph, and lifecycle history/rollback landed once agent-core routed the
--request checkpoint id and to_iteration into a target-pinned checkpoint backend
and the tool configs. The rest family proves its inbound webhook server boundary
to Succeeded over the shipped REST config; the shipped MachineSpec declares
one-shot readiness probes and delay retries while legacy await_operation is
rejected. The rest sample client create/probe
head is not run in a profile machine because REST client tools enforce a
declared-only runtime input contract that the loop's verbatim output threading
does not satisfy.

Release 08.0 adds the Chroma corpus agents. srd012 owns the family and the two
use cases (rel08.0-uc001-corpus-ingest, rel08.0-uc002-corpus-query) trace to it
through touchpoints. The provider-embedding threading
that release 07.1 recorded as missing for the rest client head -- REST client
tools enforcing a declared-only runtime input contract that verbatim output
threading does not satisfy -- is now available: agent-core #245 added
previous-Result parameter threading (srd028 R12, body_source previous_result,
input_mapping, carry_forward, single-token vector reshape), which the Chroma
add and query words use to carry the Ollama-computed vector. The implementation
landed: the profile-owned Chroma and Ollama REST asset and the ingest and
reader profiles (#242-#244). The mesh integration:chroma target now consumes
the canonical ingest program by reference while retaining its application
REST configuration and corpus fixture. Running the agents surfaced runtime and reader defects,
fixed under #250, and revealed that trusted corpus discovery must precede
model file selection; #251 introduced that loop and #849 made it canonical.
The mesh extraction owns the live Chroma harness and application REST values,
while its wrapper selects the library ingest program. That target now proves
the canonical ingest/count path. The profile-owned reader conformance case
executes the shipped wrapper against deterministic Chroma and Ollama protocol
fixtures, asserts that the provider vector reaches the query, verifies the
retrieved chunk identity reaches invoke_llm, and observes Succeeded; a
companion readiness-failure case observes Failed before embedding.
