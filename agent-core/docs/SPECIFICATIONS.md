---
id: agent-core-specifications
title: Agent Core Specification Index
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

Specification index for the agent-core framework. Fifty SRDs define the
generic agentic loop engine and its supporting packages.
Releases 00.0 through 08.0 have implemented proof in their indexed use cases
and test suites. Release 03.0 includes named executable REST conformance
evidence and rollback-time compensation execution. Release 00.0 includes
executable AC6 profile-closure timeout inspection and cmd/agent startup
enforcement tests over core-owned synthetic fixtures. Replacing every
application owner's inventory check remains the GH-1679 follow-up.
Thirty-six use cases exercise the SRDs across nineteen releases. Twenty-one
formal YAML test suites cover the acceptance cases across releases 00.0
through 17.0 and release 99.0; additional
coverage comes from Go test functions across all packages (run
`go test ./... -v` for current counts).

The monitor release adds embedded runtime monitoring. Current code defines
the monitor store, recorder, metric validation, disabled mode, OpenTelemetry
metric export, REST/OpenAPI reads from cached monitor state, live cmd/agent
startup wiring, workflow labels, and the Release 04.0 proof.

SRD release assignments are range based. Release 00.0 covers srd001
through srd018. Release 01.0 covers srd019 through srd024. Release
02.0 covers srd025, srd026, srd035, and srd036. Release 03.0 covers srd027
through srd031. Release 04.0 covers srd032 and srd033. Release 06.0 covers
srd034. Release 07.0 covers srd038 and the amended srd036 and srd028 for the
command-state store and two-plane persistence (implemented; the end-to-end
process-restart lifecycle tracer is component-level, see test-rel07.0). Release 08.0 covers
the amended srd016 for HTTP trace propagation. Release 09.0 covers srd039 and
srd040 for the scenario rig runtime: the mock REST endpoint binding and the
service-child, validator, and scenario-discovery tools. Release 10.0 covers
srd041 for the value predicate word. Release 11.0 covers srd042 for the OTLP
receiver tool family. Release 12.0 covers srd044 for the spool query tools.
Release 13.0 covers srd045 for the span-analytics tools.
Release 14.0 covers srd046 for request-driven signal admission and srd047 for
configured Dolt provision, query, and write boundary words. Their formal
suite specifies the GH-1732 real-Dolt integration obligation and the GH-1733
model-free and model-path compatibility obligations.
Release 15.0 covers srd048 for Cohere v2 Chat behind invoke_llm, including
lazy process credentials, content-block normalization, response-profile
parsing, safe errors, and GenAI provider identity.
Release 16.0 covers srd049 for the pure marshalling pipeline word.
Release 17.0 covers srd050 for named declaration units, explicit tool and
REST imports, deterministic merge diagnostics, and closure usedness.
Release 18.0 covers srd052 for parameterized declaration fragments: typed
parameters, hygienic expansion before strict decoding, provenance in
the dump, and machine stage splicing.
Release 19.0 covers srd051 R6.14 and R6.15: boundary and stateful_internal
tools declare a signature naming their signals, which discharges their
descriptive prose while their effect blocks stay explicit, and every tool
states a category.
Release 20.0 covers srd054 for machine templates, whole machines declared
once as fragments and expanded by machine files, and srd055 for agent
blueprints, whole profiles expanded per agent with their machine
templates and units resolved relative to the blueprint. It also covers
srd056 for unit libraries: imports under the agent-core install root or a
root the profile declares, read-only libraries, and rooted staging, which
amends srd050 R1.2 and srd051 R1.5.
Release 22.0 covers srd058 for provider libraries: a providers library
root bound per run, a chat dialect that makes invoke_llm a template method,
embed and rerank stage fragments, and Ollama and Cohere as the first
conforming libraries replacing the compiled adapters.
Release 99.0 covers srd037 dispatch planning and srd043
application/catalog source-layout migration.

Release 02.0 storage revision.
- srd035-checkpoint-port replaces the key/blob StateStore with a two-method
  typed checkpoint port, and srd036-dolt-state-persistence defines the Dolt
  versioned-SQL adapter behind a database/sql seam, adding a git-style Revert
  for DB-persisted rollback. srd025 and srd026 now rest on the checkpoint port
  and receipt-driven undo rather than the retired engine undo-walk.

Source-of-truth note.
- docs/guides/roll-backs.md and docs/guides/tools-as-dsl.md are historical design notes. Their
  requirements have been consolidated into srd025, cfg-tool-authoring-standard,
  sm-rollback-lifecycle, and sm-tool-language. Release 02.0 use cases and
  formal test suites carry the scenario coverage. Future requirements changes
  should update those canonical documents first.

<!-- field: roadmap_summary | form: "typed table" -->
## Roadmap summary

| version | name | use_cases_done | use_cases_total | status |
| --- | --- | --- | --- | --- |
| "00.0" | Core engine and supporting packages | 2 | 2 | done |
| "01.0" | Unified agent binary and configuration-driven architecture | 3 | 3 | done |
| "02.0" | Rollback lifecycle and tool language documentation | 4 | 4 | done |
| "03.0" | Declarative REST client and server tools | 8 | 8 | done |
| "04.0" | Embedded monitor metrics and REST API | 2 | 2 | done |
| "05.0" | Monitored live model integrations | 1 | 1 | done |
| "06.0" | External agent profiles and mounted runtime assets | 1 | 1 | done |
| "07.0" | Command-state store and two-plane persistence | 2 | 2 | implemented |
| "08.0" | HTTP trace propagation | 1 | 1 | implemented |
| "09.0" | Scenario rig runtime | 0 | 1 | done |
| "10.0" | Value predicate word | 1 | 1 | implemented |
| "11.0" | OTLP receiver tools | 1 | 1 | implemented |
| "12.0" | Spool query tools | 1 | 1 | implemented |
| "13.0" | Span analytics tools | 1 | 1 | implemented |
| "14.0" | Request-driven signals and configured Dolt words | 2 | 2 | implemented |
| "15.0" | Cohere invoke_llm provider adapter | 1 | 1 | done |
| "16.0" | Pipeline word for marshalling blocks | 1 | 1 | done |
| "17.0" | Named declaration units and explicit imports | 1 | 1 | implemented |
| "18.0" | Parameterized declaration fragments | 1 | 1 | implemented |
| "19.0" | Signed boundary contracts | 1 | 1 | implemented |
| "20.0" | Machine templates, agent blueprints, and unit libraries | 2 | 3 | in_progress |
| "21.0" | Capability profiles | 0 | 1 | in_progress |
| "22.0" | Provider libraries | 1 | 1 | done |
| "23.0" | Provider failure signals | 1 | 1 | done |
| "24.0" | Objectstore tools | 1 | 1 | done |
| "25.0" | Monitor drain outcome | 1 | 1 | done |
| "26.0" | Retrieval client library | 1 | 1 | done |
| "99.0" | Architectural migrations and cleanup | 2 | 2 | done |

<!-- field: foundation_document_index -->
## Foundation document index

| id | title | summary | path |
| --- | --- | --- | --- |
| agent-core | Agent Core: Framework for Declarative Tool-Language Agents | Vision, problem statement, success criteria, boundaries, and canonical links for declarative tool-language agents. | docs/VISION.md |
| architecture-agent-core | Agent Core Architecture | Runtime, package, configuration, tool-language, lifecycle, telemetry, and specification architecture with Mermaid diagrams. | docs/ARCHITECTURE.md |

<!-- field: srd_index | form: "typed table" -->
## SRD index

| id | title | summary | status | path |
| --- | --- | --- | --- | --- |
| srd001-core-types | Core Types | State, Signal, Command, Result, Cost, Builder, and CommandResolver types. Three generic Signal constants. |  | docs/specs/software-requirements/srd001-core-types.md |
| srd002-state-machine | State Machine | TransitionTable mapping (State, Signal) to (NextState, ActionFunc). StateMachine.Step, TerminalFunc, NewStateMachine. |  | docs/specs/software-requirements/srd002-state-machine.md |
| srd003-command-dispatch | Command Dispatch | Dispatch wrapping Command.Execute with tracing, timeout, panic recovery. SafeExecute, FillDuration, ForceErrorSignal. |  | docs/specs/software-requirements/srd003-command-dispatch.md |
| srd004-tool-registry | Tool Registry | ToolSpec with Visibility and Phases. Registry pairing specs with Builders. Manifest filtering, freeze semantics. |  | docs/specs/software-requirements/srd004-tool-registry.md |
| srd005-agentic-loop | Agentic Loop | Loop function driving state machine with budget enforcement, event recording, LoopHooks, RunResult, and OTel tracing. |  | docs/specs/software-requirements/srd005-agentic-loop.md |
| srd006-yaml-machine-config | YAML State Machine Configuration | MachineSpec YAML schema, validation, BuildTransitionTable, profile-closure timeout-envelope inspection, LoadMachineSpec, and ParseMachineSpec. |  | docs/specs/software-requirements/srd006-yaml-machine-config.md |
| srd007-tracer-interface | Tracer Port Interface | "Tracer interface with five methods: Push, Event, SetAttributes, RecordError, Context. Includes NoopTracer." |  | docs/specs/software-requirements/srd007-tracer-interface.md |
| srd008-telemetry | Telemetry Implementation | OTel setup with Trace, TraceAdapter, NewRoot, ExporterConfig, and live file/OTLP exporters; offline replay is a declared OTLP-tool machine. |  | docs/specs/software-requirements/srd008-telemetry.md |
| srd009-ollama-adapter | Ollama LLM Adapter | Ollama chat behind llm.Client with model validation and tracing; the compiled adapter was retired by srd058 R5.2 and Ollama now ships as the tools/providers/ollama library. |  | docs/specs/software-requirements/srd009-ollama-adapter.md |
| srd010-conversation | Conversation Management | Single Conversation type supporting managed (Send) and manual (Append) modes with thread-safe history and system prompt management. |  | docs/specs/software-requirements/srd010-conversation.md |
| srd011-prompt-schema | Prompt Schema and Loading | Prompt struct with Role, Task, Constraints, OutputFormat. YAML loading with mode detection, defaults, validation, and assembly. |  | docs/specs/software-requirements/srd011-prompt-schema.md |
| srd012-system-prompt-template | System Prompt Template and Manifest Serialization | Embedded Go text/template for system prompts, PromptData struct, SerializeManifest, and custom template support. |  | docs/specs/software-requirements/srd012-system-prompt-template.md |
| srd013-standard-tool-library | Standard Tool Library (STL) | Overview of shared standard tools, focused internal/tools package ownership, grouped declaration layout, and retired compatibility facade history. |  | docs/specs/software-requirements/srd013-standard-tool-library.md |
| srd014-spec-graph | Specification Graph Loader and Validator | Parses SRDs, use cases, test suites into typed structs. Builds a labeled property graph for cross-artifact validation. |  | docs/specs/software-requirements/srd014-spec-graph.md |
| srd015-genai-semconv | GenAI Semantic Convention Helpers | OTel GenAI semantic convention constants, span name builders, and attribute helpers in internal/observability/telemetry/genai. |  | docs/specs/software-requirements/srd015-genai-semconv.md |
| srd016-traceparent | W3C Traceparent Parse and Format | W3C traceparent version-00 parser, formatter, and context builder in internal/observability/telemetry, plus HTTP trace-context propagation (planned) that injects traceparent on outbound REST client requests and extracts it on inbound REST server endpoints so a cross-agent trace continues. |  | docs/specs/software-requirements/srd016-traceparent.md |
| srd017-llm-extensions | LLM Package Extensions | Parse utilities, response parser profiles, invoke_llm parser selection, and domain interfaces in internal/model/llm. |  | docs/specs/software-requirements/srd017-llm-extensions.md |
| srd018-cli-flag-contract | CLI Standard Flag Contract | Canonical flag name constants in internal/support/cli, PropagateArgs for child-process argument building, RemainingBudget, and FormatFlag. |  | docs/specs/software-requirements/srd018-cli-flag-contract.md |
| srd019-eval-harness | Evaluation Harness Infrastructure | Suite configuration, grid evaluation, session execution, atomic point workspace tools, dump_config, run_agent, oracle/trace/version point result tools, collect_metrics, result analysis, and convergence tracking in internal/evaluation. |  | docs/specs/software-requirements/srd019-eval-harness.md |
| srd020-bench-web-ui | Bench Web UI | Profile-owned React frontend, generic REST and machine_request routes, evaluation artifact queries, operator event routing, and self-invoke evaluator launch. |  | docs/specs/software-requirements/srd020-bench-web-ui.md |
| srd021-agent-cli-exec-tool | Agent CLI as an Exec Tool | Requirements for declaring the agent binary as a configured exec tool for child program invocation instead of bespoke Go subprocess paths. |  | docs/specs/software-requirements/srd021-agent-cli-exec-tool.md |
| srd022-builtin-tools | Builtin Tool Declarations and Factories | Requirements for builtin declarations, grouped builtin declaration paths, target catalog and registry package ownership, typed config validation, state-machine integration, mode-local overrides, and collection words including ordered field projection. |  | docs/specs/software-requirements/srd022-builtin-tools.md |
| srd023-cli-exec-tools | CLI and Exec Tool Declarations | Requirements for declaring external CLIs as exec tools, grouped exec declaration paths, JSON-to-argv mapping, subprocess execution, command-state stdin and env overlay, output capture, declaration selection, and target exec package ownership. |  | docs/specs/software-requirements/srd023-cli-exec-tools.md |
| srd024-blocking-input-tools | Blocking Input Tools | Requirements for tools such as invoke_llm and generic REST event waits that intentionally block the loop while waiting for external model, human, or approval input. |  | docs/specs/software-requirements/srd024-blocking-input-tools.md |
| srd025-rollback-lifecycle | Rollback, Checkpoint, Suspend, and Resume Lifecycle | Requirements for checkpoint-port-backed persistence, ordered execution history, receipt-driven rollback, suspend, resume, approval gates, and profile-selected lifecycle tools. | done | docs/specs/software-requirements/srd025-rollback-lifecycle.md |
| srd026-lifecycle-tools | Lifecycle Tool Adapters | Requirements for checkpoint_history, checkpoint_rollback, and generic lifecycle control words such as exit_agent, selected by ToolDef declarations and MachineSpec transitions. Under revision for the Dolt-backed checkpoint port. | done | docs/specs/software-requirements/srd026-lifecycle-tools.md |
| srd035-checkpoint-port | Typed Checkpoint Port | Requirements for a two-method typed checkpoint port (Save/Load) over typed Position and Execution values, an opaque tool-owned Receipt on each Entry, a folded conversation snapshot, Noop and InMemory reference adapters, save timing, and a single resume contract. Replaces the key/blob StateStore. | done | docs/specs/software-requirements/srd035-checkpoint-port.md |
| srd036-dolt-state-persistence | Dolt State Persistence | Requirements for a DoltCheckpoint adapter implementing the srd035 port behind a database/sql seam, a generic five-table schema splitting tool_results into a tool_outputs forward plane and a receipts reverse plane, commit-per-step and branch-per-run versioning with merge-to-main on terminal state, Load reconstruction, and a DB-only git-style Revert for rollback. | done | docs/specs/software-requirements/srd036-dolt-state-persistence.md |
| srd027-rest-tool-runtime | Declarative REST Tool Runtime | Shared REST runtime requirements for declarative client and server boundary tools, OpenAPI imports, net/http execution, event fan-in, lifecycle control routing, safety controls, observability, and rollback. |  | docs/specs/software-requirements/srd027-rest-tool-runtime.md |
| srd028-rest-client-tools | REST Client Tools | Requirements for declarative outbound REST client words, resource operations, RPC-shaped invoke, async send/await, response shaping, retries, redaction, and compensation. |  | docs/specs/software-requirements/srd028-rest-client-tools.md |
| srd029-rest-server-tools | REST Server Tools | Requirements for declarative inbound REST server words, route launch, await, stop, generic event sources, lifecycle control endpoints, queues, OpenAPI server binding, and request validation. |  | docs/specs/software-requirements/srd029-rest-server-tools.md |
| srd030-machine-request-rest-binding | Machine Request REST Binding | Requirements for binding one inbound HTTP request to one short-lived MachineSpec execution with request mapping, response extraction, errors, timeouts, and tracing. |  | docs/specs/software-requirements/srd030-machine-request-rest-binding.md |
| srd031-declarative-document-resources | Declarative Document Resources | Requirements for filesystem-backed documentation resources with safe path resolution, include patterns, extension allowlists, YAML and Markdown response modes, and profile-owned UI placement. |  | docs/specs/software-requirements/srd031-declarative-document-resources.md |
| srd032-tool-metrics-monitor-store | Tool Metrics Monitor Store | Requirements for a bounded in-memory monitor store, ToolMetricsRecorder, metric validation, OpenTelemetry metric export, disabled-mode behavior, profile startup wiring, workflow labels, and release proof. | done | docs/specs/software-requirements/srd032-tool-metrics-monitor-store.md |
| srd033-monitor-rest-api | Monitor REST API | Requirements for read-only REST and OpenAPI monitor endpoints over cached monitor state and trusted profile declarations, including root and declared-machine, tool, metric, event, stream views, presentation metadata, and live-store wiring. | done | docs/specs/software-requirements/srd033-monitor-rest-api.md |
| srd034-external-agent-profiles | External Agent Profiles and Mounted Runtime Assets | Requirements for separating agent-core runtime assets from profile-owned programs and running containers with externally supplied profiles. | implemented | docs/specs/software-requirements/srd034-external-agent-profiles.md |
| srd037-tool-dispatch-simplification | Tool Dispatch Wiring Simplification | Planning requirements for registry error returns, unified builtin registration, and shared external-tool availability checks while preserving Command, Builder, Dispatch, and MachineSpec semantics. | in_progress | docs/specs/software-requirements/srd037-tool-dispatch-simplification.md |
| srd038-command-state-store | Command-State Store | Requirements for a core-level read-only forward view over the run's execution log, letting any word select declared params from an earlier labeled step through a $from selector, rehydrated from the srd036 tool_outputs plane after restart, with a declared-param-only authority boundary that never exposes transport or receipts, branch-scoped lifetime tied to Revert and terminal reaping, and redaction before persistence. | implemented | docs/specs/software-requirements/srd038-command-state-store.md |
| srd039-mock-endpoint-binding | Mock REST Endpoint Binding | A mock binding for inbound REST endpoints that serves canned responses from fixture data, with routes keyed by method and literal path, ordered response scripts whose last entry repeats, a 404-and-record miss path, a request log served for validator assertions, fixture validation at --validate-config and startup, and concurrency-safe sequencing. | implemented | docs/specs/software-requirements/srd039-mock-endpoint-binding.md |
| srd040-service-child-tools | Service-Child and Scenario Tools | Words for a rig machine that composes other machines - start_scenario_mock, start_scenario_subject, and stop_service manage background serve-mode children with injected environment, trusted runtime targets, and process-group reaping; profile-visible REST probes and delay transitions own readiness retries and timeout; one-validator child execution composes through MachineSpec iteration and join; init_scenario_session with next_scenario discovers and visits deterministic work under declared roots; and start_service with list_services starts any profile as a detached child under a declared count bound and reports its outcome. | implemented | docs/specs/software-requirements/srd040-service-child-tools.md |
| srd041-value-predicate-tool | Value Predicate Tool | A pure builtin word that resolves declared operands, compares them under a closed operator and type set, and emits configured satisfied or unsatisfied signals while treating selector and coercion failures as CommandError. | implemented | docs/specs/software-requirements/srd041-value-predicate-tool.md |
| srd042-otlp-receiver-tools | OTLP Receiver Tools | Trace-only boundary words for OTLP/gRPC launch, bounded batch await, stdouttrace-compatible spool, unchanged upstream relay, lifecycle stop, and ingress self-loop prevention. | implemented | docs/specs/software-requirements/srd042-otlp-receiver-tools.md |
| srd043-application-catalog-layout | Application and Catalog Source Layout | Requirements for applications/{catalog,coding-agent,chatbot-mesh}, core-owned integration profiles, owner-root execution, AGENT_CATALOG_ROOT and its temporary alias, path-class distinctions, v0 tag compatibility, and migration sequencing. | implemented | docs/specs/software-requirements/srd043-application-catalog-layout.md |
| srd044-spool-query-tools | Spool Query Tools | Read-only builtin words for paginated trace listing and per-trace span detail from the srd042 NDJSON spool, with rotation discovery and malformed-line tolerance. | implemented | docs/specs/software-requirements/srd044-spool-query-tools.md |
| srd045-span-analytics-tools | Span Analytics Tools | Read-only builtin words for aggregate reads over the srd042 NDJSON spool -- a filtered duration-over-time heatmap, group-by counts with a reported dropped remainder, and attribute divergence between a selected region and its complement. | implemented | docs/specs/software-requirements/srd045-span-analytics-tools.md |
| srd046-request-driven-signal-source | Request-Driven Signal Source | A transport-neutral SignalSource maps validated requests to declared signal and payload envelopes, atomically distinguishes accepted, refused_undeclared, and refused_conflict, preserves checkpoint and resume behavior, and allows model-free profiles without changing the model path. | implemented | docs/specs/software-requirements/srd046-request-driven-signal-source.md |
| srd047-dolt-boundary-word | Configured Dolt Boundary Words | General-purpose provision, query, and write words whose trusted configuration owns connection, database, SQL, parameter schemas, schema lifecycle, and commit messages, with bound caller params, commit-per-write, checkpoint-database separation, and real-Dolt integration proof. | implemented | docs/specs/software-requirements/srd047-dolt-boundary-word.md |
| srd048-cohere-adapter | Cohere v2 LLM Adapter | Cohere v2 Chat adapter behind llm.Client with lazy bearer credentials, ordered content-block and token normalization, declarative invoke_llm provider selection, parser profile, safe errors, and GenAI telemetry. | implemented | docs/specs/software-requirements/srd048-cohere-adapter.md |
| srd049-pipeline-word | Pipeline Word | One declared word for an ordered sequence of pure marshalling stages, with eager stage registration, previous-result chaining, selector propagation, and one machine-visible success signal. | implemented | docs/specs/software-requirements/srd049-pipeline-word.md |
| srd050-declaration-imports | Named Declaration Units and Explicit Imports | Explicit relative imports for named tool and REST declaration units, with deterministic graph resolution, duplicate and override rules, closure usedness, legacy includes compatibility, and canonical dump evidence. | specified | docs/specs/software-requirements/srd050-declaration-imports.md |
| srd052-declaration-fragments | Parameterized Declaration Fragments | Declaration units with typed parameters, expanded by importers with arguments; scalar-only substitution before strict decoding, the hand-written validation path, closure usedness and dump provenance for expansions, and machine stage fragments spliced under a prefix. | specified | docs/specs/software-requirements/srd052-declaration-fragments.md |
| srd054-machine-templates | Machine Templates | Whole machines as fragments with a machine body, expanded by machine files carrying unit, one expansion, name, and purpose; stage splicing inside a template, validation as a hand-written machine with template and instance named, and closure, digest, dump, and staging provenance. | specified | docs/specs/software-requirements/srd054-machine-templates.md |
| srd055-agent-blueprints | Agent Blueprints | Whole agent profiles as fragments with a profile body, expanded per agent with wholesale field overrides; blueprint-relative path resolution, direct machine template expansion, the instance as program identity with blueprint and templates in its digest, and dump, staging, and audit provenance. | specified | docs/specs/software-requirements/srd055-agent-blueprints.md |
| srd056-unit-libraries | Unit Libraries | Import paths relative to the declaring file or absolute under a library root -- the implicit /opt/agent-core root or a root the profile declares under libraries -- with a --library override, closure-scoped roots, read-only libraries, closure, digest, and dump provenance, rooted staging, and the service-library convention. | specified | docs/specs/software-requirements/srd056-unit-libraries.md |
| srd057-capability-profiles | Capability Profiles | A profile that does one job, terminates, and declares no serve lifecycle, control or monitor server, or static assets; the child-process, service-endpoint, and workload forms that consume it unchanged, the evaluation-scoped in-process boundary it may not assume, and the binding discipline that makes shared behaviour a named profile with findings for the machines more than one profile reaches. | specified | docs/specs/software-requirements/srd057-capability-profiles.md |
| srd058-provider-libraries | Provider Libraries and Template-Method Words | A provider as a directory of YAML reached through a declared providers root and rebound with --library; a chat dialect that invoke_llm reads as a template method in place of compiled adapters; embed and rerank stage fragments over library REST fragments under a shared signal vocabulary; Ollama and Cohere as the first libraries, and a synthetic third provider that needs no Go. | implemented | docs/specs/software-requirements/srd058-provider-libraries.md |
| srd059-objectstore-tools | Objectstore Builtin Tool Family | One provider-neutral object-storage family with read/write/list plus an internal declaration-bound delete-prefix word; compiled-in drivers are selected by the declared bucket URL, an explicit endpoint binds the local rig with no environment read, object_write undoes through a prior-object snapshot receipt, and destructive purge is bounded, verified, and externally audited. | implemented | docs/specs/software-requirements/srd059-objectstore-tools.md |
| srd060-retrieval-client-contract | Retrieval Query Client Contract | The retrieval query contract shared by applications that serve retrieval over a Chroma collection; a client-computed vector and result count in, QueryResult with ids, documents, distances, metadatas, and embedding model out, 400 as QueryRejected; declared once in agent-core as a type unit and a client fragment, the first service library of srd056 R4. | implemented | docs/specs/software-requirements/srd060-retrieval-client-contract.md |

<!-- field: config_format_index -->
## Config format index

| id | title | summary | path |
| --- | --- | --- | --- |
| cfg-machine-format | Machine Configuration Format | MachineSpec YAML schema for agent state machines, states, signals, transitions, actions, budgets, and validation. | docs/specs/config-formats/machine-format.md |
| cfg-runtime-contract | Machine + Tools Runtime Contract | Universal profile-first startup contract for runtime data flags, tool declarations, LLM boundary config, lifecycle ports, and nested machines. | docs/specs/config-formats/runtime-contract.md |
| cfg-tool-authoring-standard | Tool Definition Authoring Standard | Source-of-truth guidance for authoring tools as vocabulary words with typed contracts, declaration placement, side effects, reversibility, undo, errors, and machine-composed workflows. | docs/specs/config-formats/tool-authoring-standard.md |
| cfg-tool-declaration-format | Tool Declaration Format | YAML schema for builtin and exec tool declarations, grouped declaration files, include ordering, local overrides, and compatibility bundles. | docs/specs/config-formats/tool-declaration-format.md |
| cfg-tool-selection-format | Tool Selection Format | Name-only tools.yaml selection format used to choose available vocabulary from shared grouped declaration files and local overrides. | docs/specs/config-formats/tool-selection-format.md |
| cfg-llm-config-format | LLM Configuration Format | Per-agent invoke_llm declaration format for model/provider config, prompts, manifest state, response_profile parser selection, and time/token limits. | docs/specs/config-formats/llm-config-format.md |
| cfg-rest-tool-format | REST Tool Configuration Format | YAML format for declarative REST definitions, clients, servers, OpenAPI imports, auth profiles, limits, retry policies, response mappings, machine_request bindings, and target-format event source, lifecycle control, and document resource fields. Release 03.0 support status is recorded in the spec. | docs/specs/config-formats/rest-tool-format.md |
| cfg-agent-profile-format | Agent Profile Format | YAML format of the profile file --profile names -- machine, tool selections, tool declarations, REST definitions, config directories, workspace directory -- with path resolution, blueprints, and blueprint instances. | docs/specs/config-formats/agent-profile-format.md |
| cfg-jurist-charter-format | Jurist Charter Format | YAML format for deterministic jurist check suites, including grep_check, ref_check, consistency_check, spec_corpus, target globs, severity, and finding provenance. | docs/specs/config-formats/jurist-charter-format.md |

<!-- field: semantic_model_index -->
## Semantic model index

| id | title | summary | path |
| --- | --- | --- | --- |
| sm-command-undo-audit | Command Undo Audit | Migration checklist for commands that need real undo implementations versus long-term no-op undo behavior. | docs/specs/semantic-models/command-undo-audit.md |
| sm-tool-vocabulary-audit | Tool Vocabulary Audit | Audit classifying current tools as words, stateful internal tools, boundaries, or sentence-shaped migration candidates. | docs/specs/semantic-models/tool-vocabulary-audit.md |
| sm-rollback-lifecycle | Rollback Lifecycle Semantic Model | Semantic model for three-layer state ownership, lifecycle stores, checkpoints, suspend/resume, and rollback traversal. | docs/specs/semantic-models/rollback-lifecycle.md |
| sm-tool-language | Tool Language Semantic Model | Semantic model for tools as vocabulary words, MachineSpec composition, boundary tools, contracts, and legacy sentence migration. | docs/specs/semantic-models/tool-language.md |
| sm-rest-tool-language | REST Tool Language Semantic Model | Semantic model for declarative REST client and server boundary words, REST definitions, async request state, inbound event fan-in, lifecycle control routing, and OpenAPI import relationships. | docs/specs/semantic-models/rest-tool-language.md |
| sm-machine-request-http | Machine Request HTTP Semantic Model | Semantic model for HTTP requests as grammar input, request-machine sentences as handling behavior, and response words as the bridge back to HTTP responses. | docs/specs/semantic-models/machine-request-http.md |
| sm-document-resource-language | Document Resource Language Semantic Model | Semantic model for configured filesystem documentation corpora, safe path policy, response modes, document resource words, and human UI actor configuration. | docs/specs/semantic-models/document-resource-language.md |
| sm-tool-metrics-monitoring | Tool Metrics Monitoring Semantic Model | Semantic model for standard tool metrics, specialized metrics, recorder ownership, workflow labels, low-cardinality attributes, and monitor snapshots. | docs/specs/semantic-models/tool-metrics-monitoring.md |

<!-- field: use_case_index -->
## Use case index

| id | title | release | status | test_suite | path |
| --- | --- | --- | --- | --- | --- |
| rel00.0-uc001-domain-agent-integration | Domain agent imports agent-core and runs an agentic loop | 00.0 | done | test-rel00.0 | docs/specs/use-cases/rel00.0-uc001-domain-agent-integration.md |
| rel00.0-uc002-spec-corpus-validation | Maintainer validates the specification corpus graph | 00.0 | done | test-rel00.0 | docs/specs/use-cases/rel00.0-uc002-spec-corpus-validation.md |
| rel01.0-uc001-generator-coding | Generator agent solves Go coding tasks using YAML config | 01.0 | done | test-rel01.0 | docs/specs/use-cases/rel01.0-uc001-generator-coding.md |
| rel01.0-uc002-evaluator-benchmark | Evaluator runs generator profiles and collects results | 01.0 | done | test-rel01.0 | docs/specs/use-cases/rel01.0-uc002-evaluator-benchmark.md |
| rel01.0-uc003-bench-visualization | Bench serves web UI for evaluation result exploration | 01.0 | done | test-rel01.0 | docs/specs/use-cases/rel01.0-uc003-bench-visualization.md |
| rel02.0-uc001-approval-suspend-resume | Agent suspends at approval gate and resumes from the checkpoint port | 02.0 | done | test-rel02.0-lifecycle | docs/specs/use-cases/rel02.0-uc001-approval-suspend-resume.md |
| rel02.0-uc002-history-rollback | User inspects history and rolls back after a bad tool action | 02.0 | done | test-rel02.0-lifecycle | docs/specs/use-cases/rel02.0-uc002-history-rollback.md |
| rel02.0-uc003-tool-authoring | Developer authors a new tool as a vocabulary word | 02.0 | done | test-rel02.0-tool-contracts | docs/specs/use-cases/rel02.0-uc003-tool-authoring.md |
| rel02.0-uc004-sentence-tool-migration | Developer migrates a sentence-shaped tool into machine-composed words | 02.0 | done | test-rel02.0-tool-contracts | docs/specs/use-cases/rel02.0-uc004-sentence-tool-migration.md |
| rel03.0-uc001-rest-client-sync | Grammar invokes synchronous REST client resource words | 03.0 | done | test-rel03.0-rest-tools | docs/specs/use-cases/rel03.0-uc001-rest-client-sync.md |
| rel03.0-uc002-rest-client-async | Grammar sends and awaits an asynchronous REST client operation | 03.0 | done | test-rel03.0-rest-tools | docs/specs/use-cases/rel03.0-uc002-rest-client-async.md |
| rel03.0-uc003-rest-server-signals | Grammar receives inbound REST server signals | 03.0 | done | test-rel03.0-rest-tools | docs/specs/use-cases/rel03.0-uc003-rest-server-signals.md |
| rel03.0-uc004-openapi-import | OpenAPI import defines REST operations and endpoint bindings | 03.0 | done | test-rel03.0-rest-tools | docs/specs/use-cases/rel03.0-uc004-openapi-import.md |
| rel03.0-uc005-ollama-openapi-model-listing | Qwen agent lists Ollama models through an OpenAPI REST tool | 03.0 | done | test-rel03.0-rest-tools | docs/specs/use-cases/rel03.0-uc005-ollama-openapi-model-listing.md |
| rel03.0-uc006-documentation-curator | Knowledge Manager serves documentation curator UX and REST tools | 03.0 | done | test-rel03.0-rest-tools | docs/specs/use-cases/rel03.0-uc006-documentation-curator.md |
| rel03.0-uc007-machine-request-documentation-ux | Knowledge Manager UX uses machine_request document resources | 03.0 | done | test-rel03.0-rest-tools | docs/specs/use-cases/rel03.0-uc007-machine-request-documentation-ux.md |
| rel03.0-uc008-rest-await-any-fan-in | Staged REST await-any fan-in adds event sources | 03.0 | done | test-rel03.0-rest-tools | docs/specs/use-cases/rel03.0-uc008-rest-await-any-fan-in.md |
| rel04.0-uc001-configurable-tool-metrics | Profile configures tool metrics for embedded monitoring | 04.0 | done | test-rel04.0-monitor | docs/specs/use-cases/rel04.0-uc001-configurable-tool-metrics.md |
| rel04.0-uc002-monitor-rest-api | Supervisor reads monitor state through REST | 04.0 | done | test-rel04.0-monitor | docs/specs/use-cases/rel04.0-uc002-monitor-rest-api.md |
| rel05.0-uc001-qwen-monitor-token-progress | Qwen run exposes increasing token metrics through embedded monitor | 05.0 | implemented | test-rel05.0-monitor-qwen | docs/specs/use-cases/rel05.0-uc001-qwen-monitor-token-progress.md |
| rel06.0-uc001-mounted-profile-container-runtime | Runtime container runs an externally mounted agent profile | 06.0 | implemented | test-rel06.0-profile-migration | docs/specs/use-cases/rel06.0-uc001-mounted-profile-container-runtime.md |
| rel07.0-uc001-command-state-threading | A REST word selects a non-adjacent prior word's output by label and resumes mid-chain | 07.0 | implemented | test-rel07.0-command-state | docs/specs/use-cases/rel07.0-uc001-command-state-threading.md |
| rel07.0-uc002-two-plane-revert | Reverting to an earlier step reaps later two-plane rows while the receipt walk reverses the surviving step | 07.0 | implemented | test-rel07.0-command-state | docs/specs/use-cases/rel07.0-uc002-two-plane-revert.md |
| rel08.0-uc001-http-trace-propagation | Two agents linked by a REST call produce one connected trace | 08.0 | implemented | test-rel08.0-trace-propagation | docs/specs/use-cases/rel08.0-uc001-http-trace-propagation.md |
| rel09.0-uc001-scenario-rig-run | A rig machine runs one scenario hermetically to a verdict | 09.0 | specified | test-rel09.0-scenario-rig | docs/specs/use-cases/rel09.0-uc001-scenario-rig-run.md |
| rel10.0-uc001-value-predicate-branch | A machine branches on a value it read from a service | 10.0 | implemented | test-rel10.0-value-predicate | docs/specs/use-cases/rel10.0-uc001-value-predicate-branch.md |
| rel11.0-uc001-otlp-receiver-await | A persistent machine receives, records, and relays OTLP trace and metric batches | 11.0 | implemented | test-rel11.0-otlp-receiver | docs/specs/use-cases/rel11.0-uc001-otlp-receiver-await.md |
| rel12.0-uc001-spool-query | A collector reads trace summaries and span detail from its NDJSON spool | 12.0 | implemented | test-rel12.0-spool-query | docs/specs/use-cases/rel12.0-uc001-spool-query.md |
| rel13.0-uc001-span-analytics | Aggregate reads turn a spool into a latency heatmap and a divergence ranking | 13.0 | implemented | test-rel13.0-span-analytics | docs/specs/use-cases/rel13.0-uc001-span-analytics.md |
| rel14.0-uc001-request-driven-signal | A model-free profile advances from a validated HTTP request | 14.0 | implemented | test-rel14.0-agent-input-and-dolt | docs/specs/use-cases/rel14.0-uc001-request-driven-signal.md |
| rel14.0-uc002-dolt-boundary-word | Declared Dolt words provision and update a versioned store | 14.0 | implemented | test-rel14.0-agent-input-and-dolt | docs/specs/use-cases/rel14.0-uc002-dolt-boundary-word.md |
| rel15.0-uc001-cohere-invoke | A profile selects Cohere for an invoke_llm classifier | 15.0 | implemented | test-rel15.0-cohere-invoke | docs/specs/use-cases/rel15.0-uc001-cohere-invoke.md |
| rel16.0-uc001-pipeline-word | A profile folds a marshalling block into one pipeline word | 16.0 | implemented | test-rel16.0-pipeline-word | docs/specs/use-cases/rel16.0-uc001-pipeline-word.md |
| rel17.0-uc001-declaration-imports | A profile composes named tool and REST declaration units | 17.0 | specified | test-rel17.0-declaration-imports | docs/specs/use-cases/rel17.0-uc001-declaration-imports.md |
| rel18.0-uc001-declaration-fragments | A profile expands parameterized declaration fragments | 18.0 | specified | test-rel18.0-declaration-fragments | docs/specs/use-cases/rel18.0-uc001-declaration-fragments.md |
| rel20.0-uc001-machine-templates | A machine file expands a whole-machine template | 20.0 | implemented | test-rel20.0-machine-templates | docs/specs/use-cases/rel20.0-uc001-machine-templates.md |
| rel20.0-uc002-agent-blueprints | A profile expands an agent blueprint | 20.0 | specified | test-rel20.0-machine-templates | docs/specs/use-cases/rel20.0-uc002-agent-blueprints.md |
| rel20.0-uc003-unit-libraries | A unit imports from a library root | 20.0 | implemented | test-rel20.0-unit-libraries | docs/specs/use-cases/rel20.0-uc003-unit-libraries.md |
| rel21.0-uc001-capability-profiles | One capability profile is consumed by more than one form | 21.0 | specified | test-rel21.0-capability-profiles | docs/specs/use-cases/rel21.0-uc001-capability-profiles.md |
| rel22.0-uc001-provider-libraries | A deployment changes and adds a model provider without editing an agent | 22.0 | implemented | test-rel22.0-provider-libraries | docs/specs/use-cases/rel22.0-uc001-provider-libraries.md |
| rel23.0-uc001-provider-failure-signals | A machine retries a busy model provider and gives up on a misconfigured one | 23.0 | implemented | test-rel23.0-provider-failure-signals | docs/specs/use-cases/rel23.0-uc001-provider-failure-signals.md |
| rel24.0-uc001-objectstore-tools | One profile writes objects on kind and on GCS through a values overlay | 24.0 | implemented | test-rel24.0-objectstore-tools | docs/specs/use-cases/rel24.0-uc001-objectstore-tools.md |
| rel25.0-uc001-drain-outcome | A monitor exit on a loaded host succeeds and names its unfinished drain | 25.0 | implemented | test-rel25.0-drain-outcome | docs/specs/use-cases/rel25.0-uc001-drain-outcome.md |
| rel26.0-uc001-retrieval-client-library | A chatbot queries its retrieval servers through the shared agent-core client | 26.0 | implemented | test-rel26.0-retrieval-client-library | docs/specs/use-cases/rel26.0-uc001-retrieval-client-library.md |
| rel99.0-uc001-tool-dispatch-simplification | Maintainer simplifies tool dispatch wiring without changing core dispatch | 99.0 | done | test-rel99.0-tool-dispatch-simplification | docs/specs/use-cases/rel99.0-uc001-tool-dispatch-simplification.md |
| rel99.0-uc002-application-catalog-migration | Maintainers migrate applications and reusable catalog blocks to explicit owner roots | 99.0 | done | test-rel99.0-application-catalog-migration | docs/specs/use-cases/rel99.0-uc002-application-catalog-migration.md |

<!-- field: test_suite_index | form: "typed table" -->
## Test suite index

| id | title | release | traces | formal_test_cases | status | path |
| --- | --- | --- | --- | --- | --- | --- |
| test-rel00.0 | Test Suite for Release 00.0 (Core Engine Import) | "00.0" | [rel00.0-uc001-domain-agent-integration, rel00.0-uc002-spec-corpus-validation] | 99 |  | docs/specs/test-suites/test-rel00.0.md |
| test-rel01.0 | Test Suite for Release 01.0 Unified Agent Runtime | "01.0" | [rel01.0-uc001-generator-coding, rel01.0-uc002-evaluator-benchmark, rel01.0-uc003-bench-visualization] | 38 |  | docs/specs/test-suites/test-rel01.0.md |
| test-rel02.0-lifecycle | Test Suite for Release 02.0 Checkpoint Port and Dolt Persistence | "02.0" | [rel02.0-uc001-approval-suspend-resume, rel02.0-uc002-history-rollback] | 15 |  | docs/specs/test-suites/test-rel02.0-lifecycle.md |
| test-rel02.0-tool-contracts | Test Suite for Release 02.0 Tool Contracts | "02.0" | [rel02.0-uc003-tool-authoring, rel02.0-uc004-sentence-tool-migration] | 16 |  | docs/specs/test-suites/test-rel02.0-tool-contracts.md |
| test-rel03.0-rest-tools | Test Suite for Release 03.0 REST Tools | "03.0" | [rel03.0-uc001-rest-client-sync, rel03.0-uc002-rest-client-async, rel03.0-uc003-rest-server-signals, rel03.0-uc004-openapi-import, rel03.0-uc005-ollama-openapi-model-listing, rel03.0-uc006-documentation-curator, rel03.0-uc007-machine-request-documentation-ux, rel03.0-uc008-rest-await-any-fan-in] | 50 | implemented | docs/specs/test-suites/test-rel03.0-rest-tools.md |
| test-rel04.0-monitor | Test Suite for Release 04.0 Embedded Monitor | "04.0" | [rel04.0-uc001-configurable-tool-metrics, rel04.0-uc002-monitor-rest-api] | 12 |  | docs/specs/test-suites/test-rel04.0-monitor.md |
| test-rel05.0-monitor-qwen | Test Suite for Release 05.0 Qwen Monitor Integration | "05.0" | [rel05.0-uc001-qwen-monitor-token-progress] | 1 |  | docs/specs/test-suites/test-rel05.0-monitor-qwen.md |
| test-rel06.0-profile-migration | Test Suite for Release 06.0 Profile Migration | "06.0" | [rel06.0-uc001-mounted-profile-container-runtime] | 4 | implemented | docs/specs/test-suites/test-rel06.0-profile-migration.md |
| test-rel07.0-command-state | Test Suite for Release 07.0 Command-State Store and Two-Plane Persistence | "07.0" | [rel07.0-uc001-command-state-threading, rel07.0-uc002-two-plane-revert] | 21 | implemented | docs/specs/test-suites/test-rel07.0-command-state.md |
| test-rel08.0-trace-propagation | Test Suite for Release 08.0 HTTP Trace Propagation | "08.0" | [rel08.0-uc001-http-trace-propagation] | 6 | implemented | docs/specs/test-suites/test-rel08.0-trace-propagation.md |
| test-rel09.0-scenario-rig | Test Suite for Release 09.0 Scenario Rig Runtime | "09.0" | [rel09.0-uc001-scenario-rig-run] | 15 | specified | docs/specs/test-suites/test-rel09.0-scenario-rig.md |
| test-rel10.0-value-predicate | Test Suite for Release 10.0 Value Predicate Tool | "10.0" | [rel10.0-uc001-value-predicate-branch] | 8 | implemented | docs/specs/test-suites/test-rel10.0-value-predicate.md |
| test-rel11.0-otlp-receiver | Test Suite for Release 11.0 OTLP Receiver Tools | "11.0" | [rel11.0-uc001-otlp-receiver-await] | 7 | implemented | docs/specs/test-suites/test-rel11.0-otlp-receiver.md |
| test-rel12.0-spool-query | Test Suite for Release 12.0 Spool Query Tools | "12.0" | [rel12.0-uc001-spool-query] | 10 | implemented | docs/specs/test-suites/test-rel12.0-spool-query.md |
| test-rel13.0-span-analytics | Test Suite for Release 13.0 Span Analytics Tools | "13.0" | [rel13.0-uc001-span-analytics] | 8 | implemented | docs/specs/test-suites/test-rel13.0-span-analytics.md |
| test-rel14.0-agent-input-and-dolt | Test Suite for Release 14.0 Request Signals and Dolt Words | "14.0" | [rel14.0-uc001-request-driven-signal, rel14.0-uc002-dolt-boundary-word] | 13 | implemented | docs/specs/test-suites/test-rel14.0-agent-input-and-dolt.md |
| test-rel15.0-cohere-invoke | Test Suite for Release 15.0 Cohere invoke_llm | "15.0" | [rel15.0-uc001-cohere-invoke] | 8 | implemented | docs/specs/test-suites/test-rel15.0-cohere-invoke.md |
| test-rel16.0-pipeline-word | Test Suite for Release 16.0 Pipeline Word | "16.0" | [rel16.0-uc001-pipeline-word] | 5 | implemented | docs/specs/test-suites/test-rel16.0-pipeline-word.md |
| test-rel17.0-declaration-imports | Test Suite for Release 17.0 Declaration Imports | "17.0" | [rel17.0-uc001-declaration-imports] | 7 | specified | docs/specs/test-suites/test-rel17.0-declaration-imports.md |
| test-rel18.0-declaration-fragments | Test Suite for Release 18.0 Declaration Fragments | "18.0" | [rel18.0-uc001-declaration-fragments] | 14 | specified | docs/specs/test-suites/test-rel18.0-declaration-fragments.md |
| test-rel20.0-machine-templates | Test Suite for Release 20.0 Machine Templates and Agent Blueprints | "20.0" | [rel20.0-uc001-machine-templates, rel20.0-uc002-agent-blueprints] | 23 | specified | docs/specs/test-suites/test-rel20.0-machine-templates.md |
| test-rel20.0-unit-libraries | Test Suite for Release 20.0 Unit Libraries | "20.0" | [rel20.0-uc003-unit-libraries] | 12 | implemented | docs/specs/test-suites/test-rel20.0-unit-libraries.md |
| test-rel21.0-capability-profiles | Test Suite for Release 21.0 Capability Profiles | "21.0" | [rel21.0-uc001-capability-profiles] | 5 | specified | docs/specs/test-suites/test-rel21.0-capability-profiles.md |
| test-rel22.0-provider-libraries | Test Suite for Release 22.0 Provider Libraries | "22.0" | [rel22.0-uc001-provider-libraries] | 10 | implemented | docs/specs/test-suites/test-rel22.0-provider-libraries.md |
| test-rel23.0-provider-failure-signals | Test Suite for Release 23.0 Provider Failure Signals | "23.0" | [rel23.0-uc001-provider-failure-signals] | 6 | implemented | docs/specs/test-suites/test-rel23.0-provider-failure-signals.md |
| test-rel24.0-objectstore-tools | Test Suite for Release 24.0 Objectstore Tools | "24.0" | [rel24.0-uc001-objectstore-tools] | 7 | implemented | docs/specs/test-suites/test-rel24.0-objectstore-tools.md |
| test-rel25.0-drain-outcome | Test Suite for Release 25.0 Drain Outcome | "25.0" | [rel25.0-uc001-drain-outcome] | 5 | implemented | docs/specs/test-suites/test-rel25.0-drain-outcome.md |
| test-rel26.0-retrieval-client-library | Test Suite for Release 26.0 Retrieval Client Library | "26.0" | [rel26.0-uc001-retrieval-client-library] | 4 | implemented | docs/specs/test-suites/test-rel26.0-retrieval-client-library.md |
| test-rel99.0-tool-dispatch-simplification | Test Suite for Release 99.0 Tool Dispatch Simplification | "99.0" | [rel99.0-uc001-tool-dispatch-simplification] | 5 | implemented | docs/specs/test-suites/test-rel99.0-tool-dispatch-simplification.md |
| test-rel99.0-application-catalog-migration | Test Suite for Release 99.0 Application and Catalog Migration | "99.0" | [rel99.0-uc002-application-catalog-migration] | 14 | implemented | docs/specs/test-suites/test-rel99.0-application-catalog-migration.md |

<!-- field: go_test_coverage | form: "tree" -->
## Go test coverage

- note: Go test functions across all packages provide coverage beyond the formal YAML test suites. Release 00.0 YAML coverage (99 cases) covers srd001-srd018 where Go tests already prove behavior. For release 01.0, the YAML suite (37 cases) indexes generator and evaluator integration evidence now lives with applications/coding-agent and applications/catalog. The suite also covers bench integration targets plus Go tests for profile and machine behavior. Tool behavior and child-agent behavior are covered there too. Builtin, exec, Ollama, and LLM extraction contracts have focused cases. Bench API and result behavior now have direct cases. The benchEvaluator target provides the automated bench proof. Release 02.0 YAML coverage (30 cases) covers rollback lifecycle, lifecycle tool adapter, and tool contract requirements. The 03.0 YAML coverage (42 cases) names implemented Go tests for REST config loading, client execution, async send/await, server signal handling, OpenAPI import, credential resolution, redirect and size policies, server validation, redaction behavior, live Ollama OpenAPI model-listing, and documentation host workflows, generic REST fan-in, lifecycle-control routing, payload propagation, terminal-state routing, listener shutdown, machine_request conformance checks, and profile-owned browser coverage for machine_request documentation UX. Package tests cover srd009-srd018 and domain packages. Release 04.0 has package tests for the monitor store, metric export, REST API SRDs, and the cmd/agent release proof. mage integration:monitor runs the release proof. Release 05.0 starts with mage integration:ollamaMonitor for monitored Qwen token progress. Run `go test ./... -v | grep -c '=== RUN'` for the current test function count.

<!-- field: srd_to_use_case_mapping -->
## SRD to use case mapping

| use_case | srd | why_required | coverage |
| --- | --- | --- | --- |
| rel26.0-uc001-retrieval-client-library | srd060-retrieval-client-contract | The contract is R1, the agent-core placement R2, the fragment R3, and the consumer rules R4. | Full. Covers AC1 through AC4 via S1 through S4. |
| rel25.0-uc001-drain-outcome | srd033-monitor-rest-api | R6.4 releases the listener, R6.5 routes a released listener to the success terminal and keeps other faults failing, and R6.10 makes the drain outcome a structured field. | Full. Covers AC6 via S1 through S4. |
| rel24.0-uc001-objectstore-tools | srd059-objectstore-tools | The family, driver selection, explicit endpoint, receipt undo, and gate obligations are R1 through R4; the overlay-only provider switch is R2.3. | Full. Covers AC1 through AC4 via S1 through S4. |
| rel23.0-uc001-provider-failure-signals | srd058-provider-libraries | The opt-in, the transport rule, the faults that stay CommandError, and the declared-emits agreement are R2.4 and R2.5 over the R1.4 taxonomy. | Full. Covers AC8 via S1 through S5. |
| rel23.0-uc001-provider-failure-signals | srd009-ollama-adapter | R4.3 states what an unavailable backend surfaces as, and R6.3 leaves retry to the caller, which is why the caller must be told whether retrying is worth anything. | Partial. Covers R4.3 and R6.3; the adapter itself is unchanged. |
| rel00.0-uc001-domain-agent-integration | srd001-core-types | Domain agent uses State, Signal, Command, Result, Cost, Builder types from internal/runtime/core. | Full. Covers AC1 through AC5 via SC1 through SC8. |
| rel00.0-uc001-domain-agent-integration | srd002-state-machine | Domain agent builds a TransitionTable and StateMachine to drive the loop. | Full. Covers AC1 through AC4 via SC2 and test cases. |
| rel00.0-uc001-domain-agent-integration | srd003-command-dispatch | Dispatch wraps every command execution with tracing, timeout, and recovery. | Full. Covers AC1 through AC3 via SC6 and test cases. |
| rel00.0-uc001-domain-agent-integration | srd004-tool-registry | Domain agent registers tools in a Registry and generates tool lists. | Full. Covers AC1 through AC5 via SC1 and test cases. |
| rel00.0-uc001-domain-agent-integration | srd005-agentic-loop | Domain agent calls Loop with LoopParams to run the agentic loop. | Full. Covers AC1 through AC6 via SC3, SC4, SC5 and test cases. |
| rel00.0-uc001-domain-agent-integration | srd006-yaml-machine-config | Domain agent can define state machines via YAML using ParseMachineSpec and BuildTransitionTable. | AC1 through AC6 have executable agent-core evidence. GH-1678 covers reusable closure-derived authority inspection, environment-expanded rejection, and cmd/agent startup enforcement; replacing all owner-module inventory checks remains GH-1679 follow-up work. |
| rel00.0-uc001-domain-agent-integration | srd007-tracer-interface | Domain agent supplies a Tracer implementation satisfying the port interface. | Full. Covers AC1 and AC2 via SC8. |
| rel00.0-uc001-domain-agent-integration | srd008-telemetry | Domain agent uses internal/observability/telemetry to set up OTel providers, create root spans, and adapt Trace to the Tracer interface. | Full. Covers AC1 through AC9 via SC8 and telemetry initialization. |

<!-- field: coverage_gaps -->
## Coverage gaps

Release 15.0 implements srd048 and rel15.0-uc001. Seven Cohere-specific cases
and the existing Ollama compatibility case in test-rel15.0-cohere-invoke are
implemented. No live Cohere credential is a release prerequisite because
the provider wire and tracing proof is hermetic.

Release 14.0 implements srd046 request-driven signal admission and srd047
configured Dolt boundary words. rel14.0-uc001 and rel14.0-uc002 trace every
acceptance criterion into implemented cases in
test-rel14.0-agent-input-and-dolt. The request-signal proof covers model-free
profile serving, HTTP mapping, ordinary Loop dispatch, process-local and Dolt
continuation, Noop refusal, lifecycle ownership, and the unchanged mock-model
runOrResume path. The Dolt-word proof retains the real-server gate for schema
creation, bound reads and writes, commit history, database separation,
authority refusal, and failed-write rollback.

Release 99.0 srd043, rel99.0-uc002, and
test-rel99.0-application-catalog-migration define the application/catalog
migration. Fourteen implemented formal cases map one-to-one to GH-1036 child
issues; the source, compatibility, taxonomy, and owner-root gate criteria are
complete.

Original core SRDs srd001-srd008 have use case touchpoints and
formal YAML test suite coverage in release 00.0. Release 00.0 now also has
formal YAML cases for implemented conversation and prompt behavior. Spec
graph behavior is covered too. GenAI, traceparent, CLI helper, and STL
behavior now have cases. Core span override behavior, ToolMetrics JSON
structure, and the no-op tracer contract now have focused formal YAML cases
tied to dedicated Go tests in the owning packages.
Ollama options and Ollama migration proof now have focused formal YAML cases.
LLM extraction functions have a focused case as well. SRDs srd019-srd024 (01.0) now have formal YAML
suite coverage tied to integration targets and package tests.
The release 00.0 warning set is now narrower, without reopening broad adapter
or parser questions that the current tests already answer.
srd025-srd026, srd035, and srd036 (02.0) have formal YAML coverage in the
release 02.0 lifecycle and tool-contract test suites. Core checkpoint,
history, resume, rollback behavior, lifecycle tool adapters, approval-gate
reference behavior, tool authoring checks, and sentence-tool migration proofs
are implemented and tested.

For 03.0, srd027-srd031 have formal YAML coverage in
test-rel03.0-rest-tools. rel03.0-uc007 covers the Knowledge Manager
machine-request demo as a profile-owned browser tracer in applications/catalog.
Generic core tests fetch machine_request document responses through the
generic REST server routes
documentation_curator.documents and
documentation_curator.document, checks raw YAML and rendered
content, and records request trace evidence.
internal/tools/rest, tools/builtin/rest declarations,
OpenAPI import compilation, REST runtime tests, testdata/conformance/rest
profile examples, and the profile-owned documentation application are implemented.
docs/specs/audits/rel03-rest-design-audit.md records the design
audit result. REST side-effect vocabulary now accepts listener and listener
shutdown effects in the contract audit. Credential resolution, redirect
allowlists, request and response size limits, server query and header
validation, and redaction now have executable coverage. invoke_handler and
stream_events bindings have executable server coverage. Shutdown policy
validation, queue draining, await unblocking, and stream shutdown unblocking
have executable coverage. The Qwen Ollama OpenAPI E2E flow,
profile-owned documentation-curator proof, and machine_request browser proof
  now have formal release 03.0 coverage. Generic REST event fan-in, route and
  signal filters, timeout handling, machine_request fixture conformance,
  documentation route conformance, and REST-triggered lifecycle exit now
  have formal cases tied to Go tests and profile-owned uc006/uc007 evidence.
  rel03.0-uc008 now has executable staged
  await-any fan-in coverage that launches a later REST source after the first
  selected event routes through MachineSpec.
  Closed REST conformance work now covers client safety, async request
  metadata, OpenAPI endpoint preservation, REST compensation mementos and
  rollback execution, machine_request terminal response schemas, and
  document_resources validation rejection. The formal suite names executable
  tests for those contracts. Modeled lifecycle_control endpoint config and
  non-default read_policy behavior stay reserved or rejected until separate
  implementation work.
Release 04.0 package tests cover the monitor store, recorder ownership,
metric config validation, disabled mode, OpenTelemetry metric export,
REST/OpenAPI cached-state readers, and the cmd/agent live-store release
proof.
Lifecycle exit coverage asserts event payload propagation, structured
exit_agent output, terminal Done or Succeeded routing, absence of cancelled
status, and documentation/control listener release.

Release 06.0 coverage now names implemented proof points for mounted profile
execution. mage audit indexes srd034, rel06.0-uc001, and
test-rel06.0-profile-migration. pkg/spec tests cover installed
/opt/agent-core/tools mapping through AGENT_CORE_HOME. magefiles Docker tests
cover profile exclusion, the required --profile entrypoint, mounted profile
command output, and the Dockerfile integration target. Application-level
executor and critic proofs live under applications/coding-agent; profile demos,
fixtures, and Knowledge Manager UX assets remain profile-owned. The matching applications/catalog
containerSmoke target runs the installed image with the profile repository
mounted at /profiles.
Current examples no longer imply that agent-core ships standard agent
profiles under its source tree or under /opt/agent-core. Profile examples use
AGENT_CATALOG_ROOT, /profiles/agents, or profile-root-relative paths.

Current `go list ./...` output defines package coverage. Direct SRD
coverage is clear for runtime core through srd001-srd006 and srd025.
Observability and model-boundary packages map to srd007-srd018 plus
srd022-srd024 and srd026. Prompt, tool adapter, CLI, and spec graph packages
map there too. Those direct mappings cover runtime support and model
boundaries. They also cover tools, CLI support, and spec validation.

Uncovered release 01.0 and 02.0 criteria are deferred for specific
follow-up. Core type and tracer coverage now has focused YAML cases.
Tool declaration coverage now has focused YAML cases for STL, builtin, and
exec criteria. Ollama and LLM extraction coverage has focused cases too. CLI
helper and evaluator statistics now have focused formal YAML cases. Bench API
and child-agent outcome criteria have focused cases too. Blocking-input
criteria remain covered by the focused cases added earlier.
Keeping those follow-up issue IDs in this index lets a reader separate
planned documentation cleanup from missing runtime behavior that would change
the product contract.
A future maintainer can follow each issue ID back to the exact missing proof,
which keeps shipped behavior separate from work that still waits for release
planning during a later cleanup pass.

Use cases and formal suites cover package groups that implement
configured agent behavior. The cmd/agent package is covered by
cfg-runtime-contract, release 01.0 use cases, and test-rel01.0. Lifecycle use
cases and Go tests cover the remaining cmd/agent behavior. Evaluation
packages are covered by srd019, srd020, rel01.0-uc002, and rel01.0-uc003.
test-rel01.0 and package tests add the executable proof. Agent subprocess
execution is covered by srd021 and srd022 AC6. rel01.0-uc002, test-rel01.0,
and Go tests provide the scenarios.

Semantic models cover cross-cutting packages that do not need separate SRDs
today. The internal/planning packages are covered by planner machine configs
and sm-tool-language. sm-tool-vocabulary-audit and planning tests add
validation coverage. internal/support/subprocess is covered by srd013 STL
delegation and srd021 child-agent execution. Execute and STL integration
tests provide the package proof.

Root packages and assets are not specification gaps. internal/model,
internal/observability, and internal/planning are package roots or
documentation stubs with no standalone runtime behavior. internal/runtime,
internal/support, and internal/tools have the same status.
External bench and documentation-curator profiles own their React source,
built assets, route configuration, request MachineSpecs, and workflow ToolDefs.
srd020 and rel03.0-uc006 cover that ownership boundary, while rel03.0-uc007
covers the machine_request-backed human actor path through the generic REST
server.

Residual documentation cleanup is tracked as follow-up work. Historical
docs/guides/roll-backs.md and docs/guides/tools-as-dsl.md references should be normalized to canonical
spec IDs. docs/guides/lifecycle-rollback.md needs the same treatment. Remaining
uncovered-AC warnings track missing formal YAML test cases for packages that
already have Go tests.
This split matters because YAML suites provide traceability evidence, while
package tests remain the executable proof for implemented behavior.
Browser-level bench UI proof can become an integration target when a browser
harness joins the quality gate. Planner package SRDs are optional unless
planner behavior becomes a stable product surface outside profile configs.
