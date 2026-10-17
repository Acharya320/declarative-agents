---
id: sm-multi-agent-profile-system
title: Multi-Agent Profile System
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: purpose -->
## Purpose

Define how migrated profile families relate to each other as configured
functional blocks. The model describes roles, boundaries, and dependency
directions without treating profiles as Go packages.

<!-- field: overview -->
## Overview

Think of the profile system as a directed graph of configured behaviors.
Work-producing agents sit beside critics and operational boundary agents,
but their edges carry different meaning. Migration has to preserve those
edges after the files leave agent-core; otherwise profile ownership becomes
visible while behavior quietly changes.

When planner delegates task execution, the edge points to executor. A
benchmark point uses a similar edge when critic runs a subject profile.
From the human side, bench starts critic through an operator action in the
web UI. Boundary-oriented profiles face outward toward REST callers,
supervisors, checkpoint stores, document corpora, and browser users.

These links do not turn applications/catalog into a runtime implementation
repository. They name configuration, actor contracts, and dependency
direction. MachineSpec interpretation, tool dispatch, nested loop execution,
and reusable runtime ports remain in agent-core.

<!-- field: system_boundary | form: "tree" -->
## System boundary

- profile_repository: applications/catalog
- runtime_repository: agent-core
- rule:
  applications/catalog owns configured agent programs and profile-local assets.
  agent-core owns runtime execution, builtin implementations, shared runtime
  assets, and container build mechanics.

<!-- field: role_groups | form: "tree" -->
## Role groups

1. - id: production_work
   - profiles:
     1. executor
     1. planner
     1. deployment-applier
   - summary:
     Production work covers task planning, coding, and deployment actuation.
     Planner decomposes requirements and invokes executor. Executor edits a
     target workspace, while applier realizes a caller-decided deployment
     change through application-owned rollout boundaries.
1. - id: evaluation_work
   - profiles:
     1. critic
     1. bench
   - summary:
     Evaluation work measures agent behavior. Bench accepts human launch
     requests, while critic expands suites and runs point pipelines.
     Configured child profile execution stays inside the critic point flow.
1. - id: validation_work
   - profiles:
     1. specification-critic
   - summary:
     Validation work reads documentation and specification corpora through a
     deterministic tool pipeline. No model boundary participates in this role.
1. - id: knowledge_work
   - profiles:
     1. documentation-curator
     1. corpus-ingest
     1. corpus-reader
   - summary:
     Knowledge work curates documentation and maintains a reusable retrieval
     corpus. Corpus profiles compute provider embeddings and sequence explicit
     Chroma boundaries; documentation-curator serves human-facing workflows.
1. - id: test_infrastructure
   - profiles:
     1. scenario-critic
     1. mock
   - summary:
     Supported test-time members compose profile-owned scenarios and replace
     external REST dependencies with deterministic fixtures.
1. - id: diagnostic_work
   - profiles:
     1. rig-doctor
   - summary:
     Diagnostic work reads one captured rig evidence directory and the trace
     spool tail it carries, correlates events, rollout counts, logs, and
     spans, and writes a diagnosis whose claims cite the captured files. No
     cluster boundary participates; the snapshot is the only input.
1. - id: boundary_runtime_work
   - profiles:
     1. rest
     1. runtime-state-reader
     1. collector
     1. control
     1. lifecycle
     1. lifecycle-exit
     1. application-purge
   - summary:
     Boundary runtime work exercises REST, monitor, control, and lifecycle
     vocabulary. Each boundary appears as explicit MachineSpec steps rather
     than hidden runtime branching.

<!-- field: dependency_directions | form: "typed table" -->
## Dependency directions

| from | to | kind | meaning | evidence |
| --- | --- | --- | --- | --- |
| planner | executor | child_profile_execution | Planner materializes and executes tasks through the executor profile. Executor does not depend on planner. | rel07.0-uc003-planner-generator-profile-boundary |
| critic | executor | benchmark_subject_execution | Critic point runs can invoke executor or another configured profile as the benchmark subject. | rel07.0-uc002-evaluator-generator-profile-boundary |
| bench | critic | human_launched_evaluation | Bench serves the human UI and launches critic runs. | rel07.0-uc004-bench-evaluator-profile-boundary |
| runtime-state-reader | agent-core monitor store | runtime_state_read | Runtime-state-reader exposes cached monitor state supplied by agent-core. | rel07.0-uc005-monitor-control-profile-boundary |
| control | lifecycle vocabulary | signal_to_lifecycle_action | Control routes REST events to visible lifecycle words such as exit_agent. | rel07.0-uc005-monitor-control-profile-boundary |
| documentation-curator | control | embedded_control_flow | The Knowledge Manager documentation profile launches control routes and waits for explicit exit events. | rel07.0-uc001-documentation-curator-profile-ux |
| corpus-ingest | provider embedding and Chroma | configured_corpus_write | Corpus ingest discovers trusted files, computes provider vectors, and writes documents and ids through profile-owned REST operations. | rel08.0-uc001-corpus-ingest |
| corpus-reader | provider embedding, Chroma, and model | grounded_retrieval | Corpus reader threads a provider vector into Chroma query and gives the retrieved identities and chunks to the answer model. | rel08.0-uc002-corpus-query |
| scenario-critic | mock, subject, and validator profiles | scenario_composition | Scenario-critic discovers scenario data, starts deterministic mocks and a fresh subject, runs validators, reduces a verdict, and stops all children. | rel11.0-uc001-scenario-critic-scenario-run |
| mock | subject REST dependency | deterministic_test_boundary | Mock serves scenario-owned responses and records subject requests. | rel11.0-uc002-mock-serves-fixtures |
| lifecycle profiles | agent-core lifecycle ports | checkpoint_operation | Lifecycle profiles select checkpoint, rollback, suspend, and resume vocabulary implemented by agent-core. |  |
| rest | external HTTP systems | configured_boundary | REST profile definitions describe trusted outbound and inbound HTTP boundaries. |  |
| deployment-applier | application-owned rollout boundaries | configured_deployment_actuation | Applier supplies lifecycle and compensation sequencing while each application supplies chart, release, namespace, Deployment, endpoint, and authority configuration. |  |

<!-- field: actor_boundaries | form: "typed table" -->
## Actor boundaries

| profile | actors | boundary |
| --- | --- | --- |
| executor | [model, workspace, validation commands] | LLM responses choose selected tools; filesystem and validation commands mutate or inspect the workspace. |
| rig-doctor | [model, captured evidence directory] | LLM responses choose filesystem reads and spool queries over one captured directory; the only write is the diagnosis beside the evidence. |
| deployment-applier | [authorized apply caller, application rollout boundary, lifecycle caller] | The catalog machine sequences request hosting, validation, apply, verification, and compensation. Applications retain every deployment coordinate, endpoint, and authorization decision. |
| critic | [suite author, child agent process, oracle command] | Suites define points; nested point runs execute child profiles and collect result artifacts. |
| planner | [requirement graph, model, issue queue, child executor process] | Planner turns requirements into issues and invokes executor for execution. |
| specification-critic | [specification corpus] | Specification-critic reads local docs and emits validation reports. |
| bench | [human operator, critic process] | Human HTTP actions launch critic runs. |
| rest | [configured upstream HTTP service, inbound webhook sender] | REST definitions constrain client calls and server events. |
| runtime-state-reader | [supervisor, runtime monitor store] | The infrastructure adapter lets a supervisor read cached monitor state and request controlled exit; it does not implement canonical Monitor. |
| collector | [OTLP trace sources, spool filesystem, query HTTP clients] | The collector receives OTLP trace batches, spools them to NDJSON, and serves bounded reads through query endpoints. |
| control | [lifecycle supervisor] | REST events become MachineSpec signals. |
| lifecycle | [operator, checkpoint store, workspace] | Lifecycle tools read or rewrite checkpoints and route approval decisions. |
| lifecycle-exit | [supervisor] | A fixed exit request becomes one control-server POST that emits ExitRequested. |
| application-purge | [application retention authority, rendered object bucket prefix, external audit file] | An exact application confirmation permits one bounded, verified deletion of declaration-fixed retained data; request content cannot select storage authority. |
| documentation-curator | [human browser, documentation corpus, local REST API] | Browser requests read or manage documentation through configured document resource and control routes. |
| corpus-ingest | [trusted corpus, embedding provider, Chroma] | The profile lists trusted paths before model selection, computes embeddings at the provider, and writes identified documents to Chroma. |
| corpus-reader | [questioner, embedding provider, Chroma, answer model] | The profile computes a query embedding, retrieves identified chunks, and constrains the model answer to those chunks. |
| scenario-critic | [scenario author, child subject, validator profiles, mock profiles] | Scenario data controls a fresh test topology whose terminal results are reduced to one verdict. |
| mock | [subject agent under test, scenario fixture] | Fixture routes provide deterministic REST responses and a per-instance request log. |

<!-- field: profile_asset_rules -->
## Profile asset rules

- Profiles may depend on agent-core builtin declarations through documented core tool roots.
- Profile-local declarations stay beside the profile that owns their behavior.
- REST definitions stay profile-owned when they define trusted external API authority.
- UI placement and human actor config stay profile-owned.
- MachineSpec owns sequencing between tools, child profiles, REST events, monitor events, and lifecycle actions.

<!-- field: validation_views -->
## Validation views

- Profile validation checks profile.yaml references, machine paths, tool selections, declarations, and REST definitions.
- Mounted container execution checks that /profiles paths can run selected profiles without bundled agent directories in the runtime image.
- Integration tests check profile relationships such as bench to critic, critic to executor, planner to executor, and Knowledge Manager control exit.
- Corpus conformance checks provider-vector threading, retrieved identities, grounded answers, and readiness failures.
- Rig conformance checks scenario-critic-to-mock/subject/validator composition and deterministic verdict reduction.
- Release 07.0 records implemented documentation-curator, bench-to-critic, and monitor/control relationship checks in test-rel07.0-profile-integrations; its critic-to-executor and planner-to-executor entries are retired library tracers that point to application-owned replacements.
- Documentation review checks that requirements describe profiles as configured blocks rather than Go packages.
