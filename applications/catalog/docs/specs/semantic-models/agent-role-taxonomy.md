---
id: sm-agent-role-taxonomy
title: Agent Role Taxonomy
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: purpose -->
## Purpose

Describe catalog profile realizations and non-role classifications by actor
boundary, inputs, outputs, selected tool families, runtime dependencies, and
checkable outcomes. Canonical role identity and actor-role bindings are owned
by the shared realization model.

<!-- field: overview -->
## Overview

Entries give migration issues a compact checklist for each profile family.
They describe configured realizations, wrappers, human boundaries, test
harnesses, mocks, and fixtures; they do not define canonical roles. The
historical `roles` key remains for compatibility with existing conformance
readers, but each entry is interpreted as a profile realization or
classification under sm-agent-role-realizations.

Executor, critic, planner, specification-critic, and Knowledge Manager profiles realize
canonical responsibilities. Bench is a human interface; scenario-critic is a test
harness with Critic bindings; mock is a test double; runtime-state-reader is
an infrastructure adapter adjacent to, but not implementing, canonical
Monitor. Those classifications describe configured behavior, not package
names, binary names, or additions to the reference role catalog.

Tool fields follow the catalog taxonomy. A role's tools.yaml is a name-only
selection list. ToolDef files under the role's agents/<family>/ closure are
reusable catalog declaration blocks or profile-local overrides. Shared core
declarations remain under agent-core/tools, and their Go implementations
remain under agent-core/internal/tools. Application-only declarations,
topology, wrappers, deployment, and UX are excluded from catalog membership.

<!-- field: shared_authority -->
## Shared authority

<!-- field: id -->
### ID

sm-agent-role-realizations

<!-- field: path -->
### Path

../../../../docs/specs/semantic-models/agent-role-realizations.md

<!-- field: owns -->
### Owns

- pinned canonical role and sub-role IDs
- actor-to-role and actor-to-sub-role bindings
- naming grammar
- realization classifications
- wrapper and variant inheritance
- migration aliases

<!-- field: reciprocal_index -->
### Reciprocal index

../../SPECIFICATIONS.md#semantic_model_index

<!-- field: compatibility_note -->
### Compatibility note

The `roles` sequence below is retained until executable governance migrates
existing readers. Its IDs identify catalog profile entries, not a second
canonical role authority.

<!-- field: roles -->
## Roles

<!-- role: executor -->
### executor

<!-- field: title -->
#### Title

Coding agent

<!-- field: profile_path -->
#### Profile path

agents/executor/profile.yaml

<!-- field: actor_boundary -->
#### Actor boundary

model plus workspace tools

<!-- field: inputs -->
#### Inputs

- task prompt
- mounted workspace
- selected model config

<!-- field: selected_tool_families -->
#### Selected tool families

- filesystem
- llm
- go exec validation

<!-- field: runtime_dependencies -->
#### Runtime dependencies

- agent-core loop runtime
- agent-core builtin filesystem and LLM tools
- agent-core exec tool declarations for Go validation

<!-- field: outputs -->
#### Outputs

- workspace changes
- validation results
- terminal status
- task summary

<!-- field: outcomes -->
#### Outcomes

- Succeeded
- Failed
- BudgetExceeded

<!-- field: checks -->
#### Checks

- Load profile from the demo.yaml catalog_root.
- Run mounted profile command with --profile and --directory.
- Confirm build, lint, and test validation are machine-declared.

<!-- role: critic -->
### critic

<!-- field: title -->
#### Title

Evaluation session and point runner

<!-- field: profile_path -->
#### Profile path

agents/critic/profile.yaml

<!-- field: actor_boundary -->
#### Actor boundary

suite fixtures plus child agent process

<!-- field: inputs -->
#### Inputs

- suite YAML
- samples
- model or profile grid
- oracle command

<!-- field: selected_tool_families -->
#### Selected tool families

- critic session
- critic point
- child profile execution

<!-- field: runtime_dependencies -->
#### Runtime dependencies

- agent-core nested loop execution
- profile-owned suites and samples

<!-- field: outputs -->
#### Outputs

- point directories
- oracle results
- trace token counts
- metrics
- session report

<!-- field: outcomes -->
#### Outcomes

- Done
- Failed

<!-- field: checks -->
#### Checks

- Resolve suite fixtures through the demo.yaml catalog_root.
- Confirm point machine runs in isolated workspaces.
- Confirm failed child runs still produce point summaries and metrics.
- Trace critic to executor with rel07.0-uc002-evaluator-generator-profile-boundary.

<!-- role: planner -->
### planner

<!-- field: title -->
#### Title

Requirement planner and issue materializer

<!-- field: profile_path -->
#### Profile path

agents/planner/profile.yaml

<!-- field: actor_boundary -->
#### Actor boundary

requirement graph, model, issue queue, and child executor

<!-- field: inputs -->
#### Inputs

- requirement graph
- SRD context
- planning model config

<!-- field: selected_tool_families -->
#### Selected tool families

- planner tools
- llm
- Go validation
- child profile execution

<!-- field: runtime_dependencies -->
#### Runtime dependencies

- agent-core planning tool implementations
- agent-core model boundary
- profile-selected issue tracker CLI declaration

<!-- field: outputs -->
#### Outputs

- implementation plan
- created issue
- graph status changes
- validation results

<!-- field: outcomes -->
#### Outcomes

- Completed
- Stalled
- Failed

<!-- field: checks -->
#### Checks

- Confirm planner invokes executor as a child profile.
- Confirm issues are materialized through format, write, profile-selected tracker exec, and result-recording words.
- Confirm graph terminal states are explicit.
- Trace planner to executor with rel07.0-uc003-planner-generator-profile-boundary.

<!-- role: specification-critic -->
### specification-critic

<!-- field: title -->
#### Title

Specification validator

<!-- field: profile_path -->
#### Profile path

agents/specification-critic/profile.yaml

<!-- field: canonical_binding | form: "tree" -->
#### Canonical binding

- role: critic
- sub_roles:
  1. output-evaluator

<!-- field: excluded_binding -->
#### Excluded binding

<!-- field: role -->
##### Role

guardian

<!-- field: sub_role -->
##### Sub role

policy-enforcer

<!-- field: reason -->
##### Reason

The profile evaluates and reports specification findings but has no policy enforcement authority.

<!-- field: actor_boundary -->
#### Actor boundary

local target directory and charter suites

<!-- field: inputs -->
#### Inputs

- target directory
- charter suite path or paths
- docs corpus for the built-in spec_corpus suite

<!-- field: selected_tool_families -->
#### Selected tool families

- specification validation
- ripgrep search and finding reduction

<!-- field: runtime_dependencies -->
#### Runtime dependencies

- agent-core spec validation tools
- rg, shared with the filesystem find capability

<!-- field: outputs -->
#### Outputs

- validation report with suite and check provenance
- pass or fail status

<!-- field: outcomes -->
#### Outcomes

- Passed
- Failed

<!-- field: checks -->
#### Checks

- Confirm no LLM boundary is selected.
- Confirm report formatting runs for pass and fail outcomes.
- Confirm charter suites run deterministically over the target directory.
- Confirm findings identify suite, check, file, and line when available.

<!-- role: bench -->
### bench

<!-- field: title -->
#### Title

Human-operated evaluation launcher

<!-- field: profile_path -->
#### Profile path

agents/bench/profile.yaml

<!-- field: actor_boundary -->
#### Actor boundary

human HTTP operator and critic process

<!-- field: inputs -->
#### Inputs

- human UI action
- experiment configuration

<!-- field: selected_tool_families -->
#### Selected tool families

- bench UI
- critic launch

<!-- field: runtime_dependencies -->
#### Runtime dependencies

- agent-core serve lifecycle and service-child words (start_service, list_services)
- critic profile

<!-- field: outputs -->
#### Outputs

- served UI
- detached critic run handle
- run listing with outcomes

<!-- field: outcomes -->
#### Outcomes

- Done
- Failed

<!-- field: checks -->
#### Checks

- Confirm the control server's lifecycle exit is the only host await; launches run through the experiments capability.
- Confirm start_service delegates session execution to critic as a detached child bounded by max_running.
- Trace bench to critic with rel07.0-uc004-bench-evaluator-profile-boundary.

<!-- role: runtime-state-reader -->
### runtime-state-reader

<!-- field: title -->
#### Title

Runtime state reader

<!-- field: profile_path -->
#### Profile path

agents/runtime-state-reader/profile.yaml

<!-- field: actor_boundary -->
#### Actor boundary

supervisor and runtime monitor store

<!-- field: classification -->
#### Classification

infrastructure_adapter

<!-- field: adjacent_role -->
#### Adjacent role

monitor

<!-- field: inputs -->
#### Inputs

- cached monitor state
- control exit request

<!-- field: selected_tool_families -->
#### Selected tool families

- monitor REST
- lifecycle control

<!-- field: runtime_dependencies -->
#### Runtime dependencies

- agent-core monitor store
- agent-core REST runtime

<!-- field: outputs -->
#### Outputs

- monitor REST responses
- listener lifecycle status

<!-- field: outcomes -->
#### Outcomes

- Done
- Failed

<!-- field: checks -->
#### Checks

- Confirm REST reads use cached monitor state.
- Confirm exit is routed through configured control.
- Confirm the adapter does not claim telemetry collection, anomaly detection, or assurance reporting.
- Trace runtime-state-reader/control status with rel07.0-uc005-monitor-control-profile-boundary.

<!-- role: lifecycle-exit -->
### lifecycle-exit

<!-- field: title -->
#### Title

Lifecycle exit client

<!-- field: profile_path -->
#### Profile path

agents/lifecycle-exit/profile.yaml

<!-- field: actor_boundary -->
#### Actor boundary

supervisor requesting a running agent's shutdown

<!-- field: classification -->
#### Classification

infrastructure_adapter

<!-- field: inputs -->
#### Inputs

- fixed exit request body

<!-- field: selected_tool_families -->
#### Selected tool families

- lifecycle control

<!-- field: runtime_dependencies -->
#### Runtime dependencies

- agent-core REST runtime

<!-- field: outputs -->
#### Outputs

- control-server exit acceptance

<!-- field: outcomes -->
#### Outcomes

- Done
- Failed

<!-- field: checks -->
#### Checks

- Confirm the exit request body is fixed in the REST operation, not runtime input.
- Confirm the target control-server address is a declared client base.
- Confirm the client posts one request through a declared boundary word rather than an embedded HTTP client.

<!-- role: application-purge -->
### application-purge

<!-- field: title -->
#### Title

Application retained-data purge

<!-- field: profile_path -->
#### Profile path

agents/application-purge/profile.yaml

<!-- field: actor_boundary -->
#### Actor boundary

retention authority deleting one apprig-rendered application prefix

<!-- field: classification -->
#### Classification

infrastructure_adapter

<!-- field: inputs -->
#### Inputs

- exact purge:<application> confirmation token

<!-- field: selected_tool_families -->
#### Selected tool families

- objectstore delete-prefix
- value predicates

<!-- field: runtime_dependencies -->
#### Runtime dependencies

- agent-core objectstore runtime
- apprig rendered bucket, endpoint, prefix, application, and audit path

<!-- field: outputs -->
#### Outputs

- bounded purge outcome
- external non-secret audit record

<!-- field: outcomes -->
#### Outcomes

- Purged
- Refused
- Incomplete
- Failed

<!-- field: checks -->
#### Checks

- Confirm request data selects no storage coordinate or credential.
- Confirm exact application confirmation precedes objectstore mutation.
- Confirm Purged requires verified zero remaining and a durable external audit.
- Confirm collector query declarations remain read-only.

<!-- role: documentation-curator -->
### documentation-curator

<!-- field: title -->
#### Title

Knowledge Manager documentation curator

<!-- field: profile_path -->
#### Profile path

agents/knowledge-manager/documentation-curator/profile.yaml

<!-- field: actor_boundary -->
#### Actor boundary

human browser, documentation corpus, and local REST API

<!-- field: inputs -->
#### Inputs

- browser request
- documentation corpus path
- control request

<!-- field: selected_tool_families -->
#### Selected tool families

- documentation serving
- document resources
- documentation validation
- patch review
- lifecycle control

<!-- field: runtime_dependencies -->
#### Runtime dependencies

- agent-core documentation host
- agent-core REST runtime
- profile-owned UI config and REST definitions

<!-- field: outputs -->
#### Outputs

- served documentation UI
- document responses
- validation or suggestion results
- patch review outcome
- controlled listener shutdown

<!-- field: outcomes -->
#### Outcomes

- Done
- Failed

<!-- field: checks -->
#### Checks

- Confirm UI placement stays profile-local.
- Confirm document resource behavior is configured by profile-owned YAML.
- Confirm exit releases the documentation listener.
- Trace profile UX with rel07.0-uc001-documentation-curator-profile-ux.

<!-- role: corpus-ingest -->
### corpus-ingest

<!-- field: title -->
#### Title

Trusted corpus ingestion

<!-- field: profile_path -->
#### Profile path

agents/knowledge-manager/corpus-ingest/profile.yaml

<!-- field: actor_boundary -->
#### Actor boundary

trusted corpus, embedding provider, and Chroma

<!-- field: inputs -->
#### Inputs

- trusted corpus root
- configured embedding model
- Chroma collection

<!-- field: selected_tool_families -->
#### Selected tool families

- filesystem resource discovery and read
- REST client
- llm

<!-- field: runtime_dependencies -->
#### Runtime dependencies

- agent-core filesystem, REST, and model boundaries
- profile-owned corpus REST operations

<!-- field: outputs -->
#### Outputs

- embedded documents and identities
- Chroma collection count

<!-- field: outcomes -->
#### Outcomes

- Succeeded
- Failed

<!-- field: checks -->
#### Checks

- Confirm trusted path listing precedes model file selection.
- Confirm provider vectors, documents, and ids reach Chroma add.
- Trace ingest behavior with rel08.0-uc001-corpus-ingest.

<!-- role: corpus-reader -->
### corpus-reader

<!-- field: title -->
#### Title

Grounded corpus retrieval

<!-- field: profile_path -->
#### Profile path

agents/knowledge-manager/corpus-reader/profile.yaml

<!-- field: actor_boundary -->
#### Actor boundary

questioner, embedding provider, Chroma, and answer model

<!-- field: inputs -->
#### Inputs

- configured question
- configured embedding and answer models
- Chroma collection

<!-- field: selected_tool_families -->
#### Selected tool families

- REST client
- llm

<!-- field: runtime_dependencies -->
#### Runtime dependencies

- agent-core REST and model boundaries
- profile-owned corpus REST operations and grounding prompt

<!-- field: outputs -->
#### Outputs

- retrieved chunk identities
- grounded answer

<!-- field: outcomes -->
#### Outcomes

- Succeeded
- Failed

<!-- field: checks -->
#### Checks

- Confirm the provider vector reaches Chroma query.
- Confirm retrieved identities and chunks reach the answer model.
- Trace reader behavior with rel08.0-uc002-corpus-query.

<!-- role: scenario-critic -->
### scenario-critic

<!-- field: title -->
#### Title

Test rig and scenario composer

<!-- field: profile_path -->
#### Profile path

agents/scenario-critic/profile.yaml

<!-- field: actor_boundary -->
#### Actor boundary

child agent processes plus profile-owned scenario data

<!-- field: inputs -->
#### Inputs

- discovery roots
- discovered scenario directories (validator profiles and mock fixtures)

<!-- field: selected_tool_families -->
#### Selected tool families

- service-child (start, await healthy, stop)
- validator execution
- scenario discovery

<!-- field: runtime_dependencies -->
#### Runtime dependencies

- agent-core service-child, validator-execution, and discovery words (srd040)
- agent-core mock endpoint binding by way of the mocks it starts (srd039)

<!-- field: outputs -->
#### Outputs

- per-scenario verdicts naming a failing validator and its cause
- aggregate verdict

<!-- field: outcomes -->
#### Outcomes

- Passed
- Failed

<!-- field: checks -->
#### Checks

- Confirm the rig composes scenarios from profile-owned data, with no rig edit needed to add a test.
- Confirm each scenario gets a fresh subject and leaves no surviving child process.
- Confirm discovery covers the applications tree without special-casing.

<!-- role: mock -->
### mock

<!-- field: title -->
#### Title

Mock for a REST dependency

<!-- field: profile_path -->
#### Profile path

agents/mock/profile.yaml

<!-- field: actor_boundary -->
#### Actor boundary

inbound HTTP from a subject agent under test

<!-- field: inputs -->
#### Inputs

- fixture file path (selected by environment)
- listen address

<!-- field: selected_tool_families -->
#### Selected tool families

- REST server
- lifecycle control

<!-- field: runtime_dependencies -->
#### Runtime dependencies

- agent-core mock endpoint binding and REST server runtime (srd039)

<!-- field: outputs -->
#### Outputs

- canned and scripted responses
- per-instance request log

<!-- field: outcomes -->
#### Outcomes

- Done
- Failed

<!-- field: checks -->
#### Checks

- Confirm one profile serves any dependency, with response content living in scenario fixtures.
- Confirm the request log is readable, per-instance, and shows unmatched calls.
- Confirm a malformed fixture blocks serving instead of failing mid-scenario.

<!-- role: collector -->
### collector

<!-- field: title -->
#### Title

OTLP trace collector with spool and query

<!-- field: profile_path -->
#### Profile path

agents/collector/profile.yaml

<!-- field: actor_boundary -->
#### Actor boundary

OTLP trace sources, spool filesystem, and query HTTP clients

<!-- field: inputs -->
#### Inputs

- OTLP trace batches
- query HTTP requests

<!-- field: selected_tool_families -->
#### Selected tool families

- OTLP receiver
- spool
- query
- REST server
- lifecycle control

<!-- field: runtime_dependencies -->
#### Runtime dependencies

- agent-core OTLP receiver tools (srd042)
- agent-core spool query tools (srd044)
- agent-core REST runtime

<!-- field: outputs -->
#### Outputs

- spooled trace NDJSON
- paginated trace list and detail responses
- lifecycle status

<!-- field: outcomes -->
#### Outcomes

- Done
- Failed

<!-- field: checks -->
#### Checks

- Confirm the receiver spools by default with no relay endpoint.
- Confirm the query surface serves bounded reads from the spool.
- Confirm shutdown ordering stops query before monitor before control.

<!-- role: rig-doctor -->
### rig-doctor

<!-- field: title -->
#### Title

Rig doctor root-cause analyst

<!-- field: profile_path -->
#### Profile path

agents/rig-doctor/profile.yaml

<!-- field: actor_boundary -->
#### Actor boundary

model plus one captured evidence directory

<!-- field: inputs -->
#### Inputs

- evidence manifest naming every captured file
- kubectl events, describe output, rollout counts, and pod logs
- trace spool tail captured beside them

<!-- field: selected_tool_families -->
#### Selected tool families

- filesystem read
- filesystem write
- spool query
- llm

<!-- field: runtime_dependencies -->
#### Runtime dependencies

- agent-core filesystem tools
- agent-core spool query tools (srd044)
- a model provider library (srd058)

<!-- field: outputs -->
#### Outputs

- diagnosis naming a probable cause with cited evidence
- one suggested read-only next command

<!-- field: outcomes -->
#### Outcomes

- Diagnosed
- NoEvidence
- Failed

<!-- field: checks -->
#### Checks

- Confirm no exec, kubectl, or helm word is reachable from the profile.
- Confirm the only write targets the diagnosis beside the evidence.
- Confirm an empty evidence directory reaches NoEvidence without a diagnosis.
- Confirm an undecodable or incomplete diagnosis returns to the model instead of reaching Diagnosed.

<!-- role: deployment-applier -->
### deployment-applier

<!-- field: title -->
#### Title

Deployment applier

<!-- field: profile_path -->
#### Profile path

agents/applier/profile.yaml

<!-- field: actor_boundary -->
#### Actor boundary

decided values document, application-owned rollout commands, and lifecycle callers

<!-- field: inputs -->
#### Inputs

- decided values content
- fixed workspace path from the trusted request binding
- application-owned REST and rollout configuration

<!-- field: selected_tool_families -->
#### Selected tool families

- lifecycle REST host
- workspace file write
- validate, apply, verify, and rollback boundaries
- rollout status reads

<!-- field: runtime_dependencies -->
#### Runtime dependencies

- agent-core REST and lifecycle tools
- application-owned Helm and kubectl ToolDefs

<!-- field: outputs -->
#### Outputs

- applied, rejected, rolled-back, or failed apply outcome
- complete, progressing, or unavailable rollout read
- controlled listener shutdown

<!-- field: outcomes -->
#### Outcomes

- Done
- Failed

<!-- field: checks -->
#### Checks

- Confirm the catalog closure boots with chart-independent standalone boundaries.
- Confirm consumers reference canonical machines, selections, and shared ToolDefs.
- Confirm chart, release, namespace, Deployment, endpoint, and authority choices remain application-owned.
- Confirm every verification failure reaches explicit rollback compensation.

<!-- field: relocated_application_roles_note -->
## Relocated application roles note

Pointer-only records: these application-specific roles are not reusable applications/catalog library members. Their requirements, checks, and status live in the canonical chatbot-mesh example.

<!-- field: relocated_application_roles -->
## Relocated application roles

| id | status | canonical_profile_path |
| --- | --- | --- |
| provisioning-workflow-orchestrator | relocated | applications/chatbot-mesh/agents/provisioning-workflow-orchestrator/profile.yaml |
| creator | relocated | applications/chatbot-mesh/agents/creator/profile.yaml |

<!-- field: conformance_fixtures_note -->
## Conformance fixtures note

Conformance grammar fixtures, not roles (GH-328). Their only executors are the conformance tests and the monitor/control mage integration target; they prove REST boundary words, control signals, and checkpoint/suspend/resume vocabulary at runtime. They live under testdata/conformance/, outside the agents/ role catalog. Application profiles may embed the same runtime vocabulary without promoting these standalone fixtures to library roles.

<!-- field: conformance_fixtures | form: "typed table" -->
## Conformance fixtures

| id | title | profile_path | actor_boundary | inputs | selected_tool_families | runtime_dependencies | outputs | outcomes | checks |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| rest | REST boundary sample | testdata/conformance/rest/profile.yaml | configured HTTP clients and servers | [REST definitions, OpenAPI import data, inbound webhook event] | [REST client, REST async client, REST server] | [agent-core REST runtime, profile-owned rest.yaml] | [response signals, async completion, webhook event signal, listener lifecycle status] | [Succeeded, Failed] | [Confirm REST authority lives in rest.yaml., "Confirm server launch, await, and stop are explicit machine states."] |
| control | REST lifecycle control profile | testdata/conformance/control/profile.yaml | lifecycle supervisor | [REST control request] | [control REST, exit_agent] | [agent-core REST runtime, agent-core lifecycle vocabulary] | [selected control signal, AgentExited signal, terminal status] | [Succeeded, Failed] | [Confirm HTTP handlers enqueue signals only., Confirm MachineSpec routes ExitRequested to exit_agent., Trace monitor/control status with rel07.0-uc005-monitor-control-profile-boundary.] |
| lifecycle | Checkpoint, rollback, and approval profiles | testdata/conformance/lifecycle/*/profile.yaml | operator, checkpoint store, and workspace | [checkpoint ID, target iteration, approval signal] | [checkpoint_history, checkpoint_rollback, suspend, done] | [agent-core lifecycle ports, checkpoint store, "receipt-backed rollback (opaque persisted receipts, receipt-walk Undo)"] | [history report, rewritten checkpoint, suspended checkpoint, approved or rejected terminal state] | [Done, Succeeded, Rejected, Failed] | [Confirm lifecycle behavior is expressed as machine actions., Confirm approval uses resume signals instead of hidden CLI branches.] |

<!-- field: taxonomy_rules -->
## Taxonomy rules

- Each entry under the compatibility `roles` key is a configured realization or classification, not a canonical role definition or Go package.
- Canonical role and sub-role identity, repeated bindings, and aliases resolve through sm-agent-role-realizations.
- tools.yaml files select vocabulary by name; they do not declare tools.
- Reusable repository-specific ToolDef files are catalog tool blocks under one canonical agents/<family>/ closure.
- Profile-local overrides configure one profile or variant without replacing core declarations for other consumers.
- Generic implementations remain under agent-core/internal/tools and core shared declarations remain under agent-core/tools.
- Application-local topology, wrappers, deployment, packaging, UX, and single-application declarations are not catalog roles or tool blocks.
- Conformance grammar fixtures (rest, control, lifecycle) are not roles; they live under testdata/conformance/ and prove runtime vocabulary, keeping agents/ a role catalog (GH-328).
- The scenario-critic is a test harness with Critic / Sandbox Validator and Output Evaluator bindings; it is not a canonical Assembler realization. Mock is a non-role test double. The environment the harness composes is a sandbox in the reference architecture's sense, and the per-scenario validators are Critic realizations. The reference has no term for a scripted test double, so mock retains the runtime binding name rather than borrowing digital twin, which the reference reserves for a replica derived from a real referent (GH-704, superseding GH-665).
- The Jurist charter format and jurist_rg tool word are stable capability vocabulary, not profile identities; their names remain unchanged by the specification-critic migration.
- Actor boundaries name who or what crosses the profile boundary.
- Runtime dependencies name agent-core capabilities without copying them into applications/catalog.
- Checks name validation, mounted execution, integration proof, or documentation review paths.
- Rel07 integration checks point at implemented cases in test-rel07.0-profile-integrations.
