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

<!-- field: shared_authority | form: "tree" -->
## Shared authority

- id: sm-agent-role-realizations
- path: ../../../../docs/specs/semantic-models/agent-role-realizations.md
- owns:
  1. pinned canonical role and sub-role IDs
  1. actor-to-role and actor-to-sub-role bindings
  1. naming grammar
  1. realization classifications
  1. wrapper and variant inheritance
  1. migration aliases
- reciprocal_index: ../../SPECIFICATIONS.md#semantic_model_index
- compatibility_note:
  The `roles` sequence below is retained until executable governance migrates
  existing readers. Its IDs identify catalog profile entries, not a second
  canonical role authority.

<!-- field: roles | form: "tree" -->
## Roles

1. - id: executor
   - title: Coding agent
   - profile_path: agents/executor/profile.yaml
   - actor_boundary: model plus workspace tools
   - inputs:
     1. task prompt
     1. mounted workspace
     1. selected model config
   - selected_tool_families:
     1. filesystem
     1. llm
     1. go exec validation
   - runtime_dependencies:
     1. agent-core loop runtime
     1. agent-core builtin filesystem and LLM tools
     1. agent-core exec tool declarations for Go validation
   - outputs:
     1. workspace changes
     1. validation results
     1. terminal status
     1. task summary
   - outcomes:
     1. Succeeded
     1. Failed
     1. BudgetExceeded
   - checks:
     1. Load profile from the demo.yaml catalog_root.
     1. Run mounted profile command with --profile and --directory.
     1. Confirm build, lint, and test validation are machine-declared.
1. - id: critic
   - title: Evaluation session and point runner
   - profile_path: agents/critic/profile.yaml
   - actor_boundary: suite fixtures plus child agent process
   - inputs:
     1. suite YAML
     1. samples
     1. model or profile grid
     1. oracle command
   - selected_tool_families:
     1. critic session
     1. critic point
     1. child profile execution
   - runtime_dependencies:
     1. agent-core nested loop execution
     1. profile-owned suites and samples
   - outputs:
     1. point directories
     1. oracle results
     1. trace token counts
     1. metrics
     1. session report
   - outcomes:
     1. Done
     1. Failed
   - checks:
     1. Resolve suite fixtures through the demo.yaml catalog_root.
     1. Confirm point machine runs in isolated workspaces.
     1. Confirm failed child runs still produce point summaries and metrics.
     1. Trace critic to executor with rel07.0-uc002-evaluator-generator-profile-boundary.
1. - id: planner
   - title: Requirement planner and issue materializer
   - profile_path: agents/planner/profile.yaml
   - actor_boundary: requirement graph, model, issue queue, and child executor
   - inputs:
     1. requirement graph
     1. SRD context
     1. planning model config
   - selected_tool_families:
     1. planner tools
     1. llm
     1. Go validation
     1. child profile execution
   - runtime_dependencies:
     1. agent-core planning tool implementations
     1. agent-core model boundary
     1. profile-selected issue tracker CLI declaration
   - outputs:
     1. implementation plan
     1. created issue
     1. graph status changes
     1. validation results
   - outcomes:
     1. Completed
     1. Stalled
     1. Failed
   - checks:
     1. Confirm planner invokes executor as a child profile.
     1. Confirm issues are materialized through format, write, profile-selected tracker exec, and result-recording words.
     1. Confirm graph terminal states are explicit.
     1. Trace planner to executor with rel07.0-uc003-planner-generator-profile-boundary.
1. - id: specification-critic
   - title: Specification validator
   - profile_path: agents/specification-critic/profile.yaml
   - canonical_binding:
     - role: critic
     - sub_roles:
       1. output-evaluator
   - excluded_binding:
     - role: guardian
     - sub_role: policy-enforcer
     - reason: The profile evaluates and reports specification findings but has no policy enforcement authority.
   - actor_boundary: local target directory and charter suites
   - inputs:
     1. target directory
     1. charter suite path or paths
     1. docs corpus for the built-in spec_corpus suite
   - selected_tool_families:
     1. specification validation
     1. ripgrep search and finding reduction
   - runtime_dependencies:
     1. agent-core spec validation tools
     1. rg, shared with the filesystem find capability
   - outputs:
     1. validation report with suite and check provenance
     1. pass or fail status
   - outcomes:
     1. Passed
     1. Failed
   - checks:
     1. Confirm no LLM boundary is selected.
     1. Confirm report formatting runs for pass and fail outcomes.
     1. Confirm charter suites run deterministically over the target directory.
     1. Confirm findings identify suite, check, file, and line when available.
1. - id: bench
   - title: Human-operated evaluation launcher
   - profile_path: agents/bench/profile.yaml
   - actor_boundary: human HTTP operator and critic process
   - inputs:
     1. human UI action
     1. experiment configuration
   - selected_tool_families:
     1. bench UI
     1. critic launch
   - runtime_dependencies:
     1. agent-core serve lifecycle and service-child words (start_service, list_services)
     1. critic profile
   - outputs:
     1. served UI
     1. detached critic run handle
     1. run listing with outcomes
   - outcomes:
     1. Done
     1. Failed
   - checks:
     1. Confirm the control server's lifecycle exit is the only host await; launches run through the experiments capability.
     1. Confirm start_service delegates session execution to critic as a detached child bounded by max_running.
     1. Trace bench to critic with rel07.0-uc004-bench-evaluator-profile-boundary.
1. - id: runtime-state-reader
   - title: Runtime state reader
   - profile_path: agents/runtime-state-reader/profile.yaml
   - actor_boundary: supervisor and runtime monitor store
   - classification: infrastructure_adapter
   - adjacent_role: monitor
   - inputs:
     1. cached monitor state
     1. control exit request
   - selected_tool_families:
     1. monitor REST
     1. lifecycle control
   - runtime_dependencies:
     1. agent-core monitor store
     1. agent-core REST runtime
   - outputs:
     1. monitor REST responses
     1. listener lifecycle status
   - outcomes:
     1. Done
     1. Failed
   - checks:
     1. Confirm REST reads use cached monitor state.
     1. Confirm exit is routed through configured control.
     1. Confirm the adapter does not claim telemetry collection, anomaly detection, or assurance reporting.
     1. Trace runtime-state-reader/control status with rel07.0-uc005-monitor-control-profile-boundary.
1. - id: lifecycle-exit
   - title: Lifecycle exit client
   - profile_path: agents/lifecycle-exit/profile.yaml
   - actor_boundary: supervisor requesting a running agent's shutdown
   - classification: infrastructure_adapter
   - inputs:
     1. fixed exit request body
   - selected_tool_families:
     1. lifecycle control
   - runtime_dependencies:
     1. agent-core REST runtime
   - outputs:
     1. control-server exit acceptance
   - outcomes:
     1. Done
     1. Failed
   - checks:
     1. Confirm the exit request body is fixed in the REST operation, not runtime input.
     1. Confirm the target control-server address is a declared client base.
     1. Confirm the client posts one request through a declared boundary word rather than an embedded HTTP client.
1. - id: application-purge
   - title: Application retained-data purge
   - profile_path: agents/application-purge/profile.yaml
   - actor_boundary: retention authority deleting one apprig-rendered application prefix
   - classification: infrastructure_adapter
   - inputs:
     1. exact purge:<application> confirmation token
   - selected_tool_families:
     1. objectstore delete-prefix
     1. value predicates
   - runtime_dependencies:
     1. agent-core objectstore runtime
     1. apprig rendered bucket, endpoint, prefix, application, and audit path
   - outputs:
     1. bounded purge outcome
     1. external non-secret audit record
   - outcomes:
     1. Purged
     1. Refused
     1. Incomplete
     1. Failed
   - checks:
     1. Confirm request data selects no storage coordinate or credential.
     1. Confirm exact application confirmation precedes objectstore mutation.
     1. Confirm Purged requires verified zero remaining and a durable external audit.
     1. Confirm collector query declarations remain read-only.
1. - id: documentation-curator
   - title: Knowledge Manager documentation curator
   - profile_path: agents/knowledge-manager/documentation-curator/profile.yaml
   - actor_boundary: human browser, documentation corpus, and local REST API
   - inputs:
     1. browser request
     1. documentation corpus path
     1. control request
   - selected_tool_families:
     1. documentation serving
     1. document resources
     1. documentation validation
     1. patch review
     1. lifecycle control
   - runtime_dependencies:
     1. agent-core documentation host
     1. agent-core REST runtime
     1. profile-owned UI config and REST definitions
   - outputs:
     1. served documentation UI
     1. document responses
     1. validation or suggestion results
     1. patch review outcome
     1. controlled listener shutdown
   - outcomes:
     1. Done
     1. Failed
   - checks:
     1. Confirm UI placement stays profile-local.
     1. Confirm document resource behavior is configured by profile-owned YAML.
     1. Confirm exit releases the documentation listener.
     1. Trace profile UX with rel07.0-uc001-documentation-curator-profile-ux.
1. - id: corpus-ingest
   - title: Trusted corpus ingestion
   - profile_path: agents/knowledge-manager/corpus-ingest/profile.yaml
   - actor_boundary: trusted corpus, embedding provider, and Chroma
   - inputs:
     1. trusted corpus root
     1. configured embedding model
     1. Chroma collection
   - selected_tool_families:
     1. filesystem resource discovery and read
     1. REST client
     1. llm
   - runtime_dependencies:
     1. agent-core filesystem, REST, and model boundaries
     1. profile-owned corpus REST operations
   - outputs:
     1. embedded documents and identities
     1. Chroma collection count
   - outcomes:
     1. Succeeded
     1. Failed
   - checks:
     1. Confirm trusted path listing precedes model file selection.
     1. Confirm provider vectors, documents, and ids reach Chroma add.
     1. Trace ingest behavior with rel08.0-uc001-corpus-ingest.
1. - id: corpus-reader
   - title: Grounded corpus retrieval
   - profile_path: agents/knowledge-manager/corpus-reader/profile.yaml
   - actor_boundary: questioner, embedding provider, Chroma, and answer model
   - inputs:
     1. configured question
     1. configured embedding and answer models
     1. Chroma collection
   - selected_tool_families:
     1. REST client
     1. llm
   - runtime_dependencies:
     1. agent-core REST and model boundaries
     1. profile-owned corpus REST operations and grounding prompt
   - outputs:
     1. retrieved chunk identities
     1. grounded answer
   - outcomes:
     1. Succeeded
     1. Failed
   - checks:
     1. Confirm the provider vector reaches Chroma query.
     1. Confirm retrieved identities and chunks reach the answer model.
     1. Trace reader behavior with rel08.0-uc002-corpus-query.
1. - id: scenario-critic
   - title: Test rig and scenario composer
   - profile_path: agents/scenario-critic/profile.yaml
   - actor_boundary: child agent processes plus profile-owned scenario data
   - inputs:
     1. discovery roots
     1. discovered scenario directories (validator profiles and mock fixtures)
   - selected_tool_families:
     1. service-child (start, await healthy, stop)
     1. validator execution
     1. scenario discovery
   - runtime_dependencies:
     1. agent-core service-child, validator-execution, and discovery words (srd040)
     1. agent-core mock endpoint binding by way of the mocks it starts (srd039)
   - outputs:
     1. per-scenario verdicts naming a failing validator and its cause
     1. aggregate verdict
   - outcomes:
     1. Passed
     1. Failed
   - checks:
     1. Confirm the rig composes scenarios from profile-owned data, with no rig edit needed to add a test.
     1. Confirm each scenario gets a fresh subject and leaves no surviving child process.
     1. Confirm discovery covers the applications tree without special-casing.
1. - id: mock
   - title: Mock for a REST dependency
   - profile_path: agents/mock/profile.yaml
   - actor_boundary: inbound HTTP from a subject agent under test
   - inputs:
     1. fixture file path (selected by environment)
     1. listen address
   - selected_tool_families:
     1. REST server
     1. lifecycle control
   - runtime_dependencies:
     1. agent-core mock endpoint binding and REST server runtime (srd039)
   - outputs:
     1. canned and scripted responses
     1. per-instance request log
   - outcomes:
     1. Done
     1. Failed
   - checks:
     1. Confirm one profile serves any dependency, with response content living in scenario fixtures.
     1. Confirm the request log is readable, per-instance, and shows unmatched calls.
     1. Confirm a malformed fixture blocks serving instead of failing mid-scenario.
1. - id: collector
   - title: OTLP trace collector with spool and query
   - profile_path: agents/collector/profile.yaml
   - actor_boundary: OTLP trace sources, spool filesystem, and query HTTP clients
   - inputs:
     1. OTLP trace batches
     1. query HTTP requests
   - selected_tool_families:
     1. OTLP receiver
     1. spool
     1. query
     1. REST server
     1. lifecycle control
   - runtime_dependencies:
     1. agent-core OTLP receiver tools (srd042)
     1. agent-core spool query tools (srd044)
     1. agent-core REST runtime
   - outputs:
     1. spooled trace NDJSON
     1. paginated trace list and detail responses
     1. lifecycle status
   - outcomes:
     1. Done
     1. Failed
   - checks:
     1. Confirm the receiver spools by default with no relay endpoint.
     1. Confirm the query surface serves bounded reads from the spool.
     1. Confirm shutdown ordering stops query before monitor before control.
1. - id: rig-doctor
   - title: Rig doctor root-cause analyst
   - profile_path: agents/rig-doctor/profile.yaml
   - actor_boundary: model plus one captured evidence directory
   - inputs:
     1. evidence manifest naming every captured file
     1. kubectl events, describe output, rollout counts, and pod logs
     1. trace spool tail captured beside them
   - selected_tool_families:
     1. filesystem read
     1. filesystem write
     1. spool query
     1. llm
   - runtime_dependencies:
     1. agent-core filesystem tools
     1. agent-core spool query tools (srd044)
     1. a model provider library (srd058)
   - outputs:
     1. diagnosis naming a probable cause with cited evidence
     1. one suggested read-only next command
   - outcomes:
     1. Diagnosed
     1. NoEvidence
     1. Failed
   - checks:
     1. Confirm no exec, kubectl, or helm word is reachable from the profile.
     1. Confirm the only write targets the diagnosis beside the evidence.
     1. Confirm an empty evidence directory reaches NoEvidence without a diagnosis.
     1. Confirm an undecodable or incomplete diagnosis returns to the model instead of reaching Diagnosed.
1. - id: deployment-applier
   - title: Deployment applier
   - profile_path: agents/applier/profile.yaml
   - actor_boundary: decided values document, application-owned rollout commands, and lifecycle callers
   - inputs:
     1. decided values content
     1. fixed workspace path from the trusted request binding
     1. application-owned REST and rollout configuration
   - selected_tool_families:
     1. lifecycle REST host
     1. workspace file write
     1. validate, apply, verify, and rollback boundaries
     1. rollout status reads
   - runtime_dependencies:
     1. agent-core REST and lifecycle tools
     1. application-owned Helm and kubectl ToolDefs
   - outputs:
     1. applied, rejected, rolled-back, or failed apply outcome
     1. complete, progressing, or unavailable rollout read
     1. controlled listener shutdown
   - outcomes:
     1. Done
     1. Failed
   - checks:
     1. Confirm the catalog closure boots with chart-independent standalone boundaries.
     1. Confirm consumers reference canonical machines, selections, and shared ToolDefs.
     1. Confirm chart, release, namespace, Deployment, endpoint, and authority choices remain application-owned.
     1. Confirm every verification failure reaches explicit rollback compensation.

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
