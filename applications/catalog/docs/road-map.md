---
id: agent-profiles-roadmap
title: Agent Profiles Release Road Map
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: numbering_policy -->
## Numbering policy

Release numbers share the number space defined by agent-core. A release exists
in this road map only when this repository has deliverables for the matching
agent-core release, so the sequence skips numbers with no profile-owned work
(hence the jump from 00.1 to 06.0). Where a release carries
related_agent_core_release, that field names the agent-core release it tracks,
and any proof command that runs from an agent-core checkout says so explicitly.

<!-- release: 00.1 | status: done -->
## 00.1 - Repository bootstrap

Define the repository purpose and ownership boundary before assets move.
The first release also adds a documentation index and migration roadmap.

Issue format rules matter because the next issues touch profile assets,
application presentations, and fixtures. With those rules in place,
future issues can move assets without re-opening the runtime ownership
decision.

<!-- field: deliverables -->
### Deliverables

- README.md documents local usage, container usage, and ownership.
- docs/VISION.yaml defines goals and boundaries.
- docs/ARCHITECTURE.yaml defines repository structure and mount interfaces.
- docs/SPECIFICATIONS.yaml indexes the initial docs.
- docs/constitutions/issue-format.yaml defines future issue fields.
- docs/specs/software-requirements/srd001-agent-functional-blocks.yaml defines migrated profile families as functional blocks.
- docs/specs/semantic-models/multi-agent-profile-system.yaml defines profile role groups and dependency directions.
- docs/specs/semantic-models/agent-role-taxonomy.yaml classifies each profile family by boundary, inputs, outputs, tools, and checks.
- magefiles/release.go adds profile bundle release tagging.
- magefiles/validation.go validates migrated profiles and mounted container usage against an external agent-core runtime.

<!-- release: 06.0 | status: done -->
## 06.0 - External agent profiles and mounted runtime assets

Profile-owned programs moved from agent-core into this repository.
At Release 06.0, demos and integration fixtures moved with the profiles
they exercised. Presentation-specific Knowledge Manager assets were later
relocated to applications/agent-architecture; this release record
preserves the historical migration boundary.

Runtime code, builtin tool implementations, release image build logic,
and core-owned assets stay in agent-core. At execution time, commands
load profiles through the demo.yaml catalog_root, /profiles, or another explicit
external path.

Release 06.0 is complete for this repository when the profile assets can
be validated, tagged, mounted into a runtime container, and consumed by
agent-core integration targets. The remaining ownership boundary belongs
to agent-core release work, not to this profile bundle.
Done here means the profile bundle has its own evidence.

Operators still need a compatible agent-core runtime image or checkout,
but this repository now supplies the profile programs, fixtures, UI
assets, validation target, smoke target, release tag target, and formal
traceability records that the runtime consumes.

<!-- field: related_agent_core_release -->
### Related agent core release

06.0

<!-- field: planned_work -->
### Planned work

| id | summary | status |
| --- | --- | --- |
| agent-profiles-dte | Specify agent functional blocks and multi-agent semantics. | done |
| agent-profiles-sqe | Migrate agent profiles and UX assets from agent-core. | done |
| agent-profiles-6ar | Historically moved demos and agent integration fixtures into applications/catalog for Release 06.0; presentation assets were later relocated to owning applications. | done |
| agent-profiles-73s | Add Mage release tag target for profile bundles. | done |
| agent-profiles-8t8 | Add profile validation and container usage checks. | done |
| agent-profiles-dqm | Validate all profile variants in applications/catalog. | done |
| agent-profiles-0g7 | Add formal profile migration use cases and test suites. | done |

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| rel06.0-uc001-mounted-profile-execution | Runtime executes profile-owned programs mounted from applications/catalog. | implemented |
| rel06.0-uc002-profile-validation | Profile validation and release bundle tagging prove repository readiness. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| test-rel06.0-agent-profiles | Formal proof suite for profile validation, mounted container execution, tagging, and agent-core integrations. | implemented |

<!-- field: proof -->
### Proof

- mage validate
- mage containerSmoke
- mage tag
- go test ./magefiles
- mage integration:codingLoop with catalog_root set to /path/to/applications/catalog in demo.yaml (run from applications/coding-agent)
- mage integration:monitor with catalog_root set to /path/to/applications/catalog in demo.yaml (run from an agent-core checkout)

<!-- field: success_markers -->
### Success markers

- Migrated profile.yaml files load through agent-core from this repository.
- Knowledge Manager UX config remains profile-local beside its profile.
- At Release 06.0, demos and fixtures used repository-relative paths, the demo.yaml catalog_root, /profiles, or explicit workspace mounts; current presentation composition lives in its owning application.
- agent-core runtime images do not copy migrated profile directories.

<!-- release: 07.0 | status: done -->
## 07.0 - Profile-owned integration tracer bullets

Profile-owned integration tests should prove concrete profile behavior
where the profile assets live. Release 07.0 defines tracer bullets for
Knowledge Manager/documentation-curator UX, evaluator to generator,
planner to generator, bench to evaluator, and monitor/control boundaries.

The release includes executable profile-owned tracer bullets for
Knowledge Manager/documentation-curator UX, bench to critic, and
monitor/control boundaries. The former evaluator-to-generator and
planner-to-generator library tracers are retired because they used shell
children that fabricated the claimed behavior. Their real
critic-to-executor and planner-to-executor replacements belong to
applications/coding-agent and are not counted as implemented library proof.

<!-- field: related_agent_core_release -->
### Related agent core release

07.0

<!-- field: planned_work -->
### Planned work

| id | summary | status |
| --- | --- | --- |
| gh-106 | Define profile integration tracer-bullet use cases and suite. | done |
| gh-107 | Move documentation-curator integration proof to applications/catalog. | done |
| gh-108 | Add evaluator-to-generator profile tracer bullet. | done |
| gh-109 | Add planner-to-generator profile tracer bullet. | done |
| gh-110 | Add bench-to-evaluator profile tracer bullet. | done |
| gh-111 | Add or explicitly scope monitor-control tracer evidence. | done |

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| rel07.0-uc001-documentation-curator-profile-ux | Knowledge Manager profile UX, document machine_request, trace, and lifecycle proof. | implemented |
| rel07.0-uc002-evaluator-generator-profile-boundary | Retired library tracer; application critic boundary is the replacement. | planned |
| rel07.0-uc003-planner-generator-profile-boundary | Retired library tracer; application planner delegation is the replacement. | planned |
| rel07.0-uc004-bench-evaluator-profile-boundary | Bench launches evaluator from a human/operator action. | implemented |
| rel07.0-uc005-monitor-control-profile-boundary | Monitor/control tracer decision and lifecycle boundary proof. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| test-rel07.0-profile-integrations | Implemented profile-owned tracer bullet evidence. | implemented |

<!-- field: proof -->
### Proof

- mage integration:all
- mage integration:documentationCurator
- mage integration:benchEvaluator
- mage integration:monitorControl

<!-- field: success_markers -->
### Success markers

- Each implemented tracer target lives in applications/catalog.
- Fixtures resolve through the demo.yaml catalog_root or repository-relative paths.
- Cross-profile boundaries record concrete output, trace, or terminal state evidence.
- agent-core remains unaware of concrete profile names in its own integration targets.

<!-- release: 07.1 | status: in-progress -->
## 07.1 - Per-family SRDs and profile coverage

Release 07.1 gives each profile family its own requirements contract and
validates the completed 07.0 tracer bullets without renumbering existing
use cases. The requirements layer moves from one shared SRD holding ten
blocks to one SRD per family, with srd001 keeping only the cross-cutting
ownership boundary and dependency-direction invariants.

The release also adds per-family boundary use cases and a test suite for
the specification-critic, lifecycle, and rest families, which were previously validated
only by the generic rel06.0 mounted-profile and profile-validation use
cases. Each new use case names the actor boundary, inputs, outputs, and
terminal outcomes for its family and traces to the family SRD.

<!-- field: planned_work -->
### Planned work

| id | summary | status |
| --- | --- | --- |
| gh-193 | Split srd001-agent-functional-blocks into per-family SRDs (srd002-executor through srd011-knowledge-manager) and retarget existing traces. | done |
| gh-194 | Add per-family boundary use cases and a test suite for the specification-critic, lifecycle, and rest profiles. | done |
| gh-195 | Add an executable per-family OpenTelemetry conformance suite (mage conformance) covering all ten profile families against srd002-srd011. | done |
| gh-229 | Exercise the shipped profile.yaml wrappers and their variants behaviorally in conformance instead of synthesized reconstructions — copy-and-patch only harness-controlled fields (addresses, ports, corpus dirs, iteration budget) — and record per-family and per-variant coverage (shipped-behavioral vs load-and-wiring) with reasons. | done |

<!-- field: deliverables -->
### Deliverables

- docs/specs/software-requirements/srd002-executor.yaml through srd011-knowledge-manager.yaml define per-family functional blocks.
- docs/specs/software-requirements/srd001-agent-functional-blocks.yaml keeps the profile ownership boundary and cross-profile dependency-direction invariants.
- docs/specs/use-cases/rel07.1-uc001-specification-critic-profile-boundary.yaml, rel07.1-uc002-lifecycle-profile-boundary.yaml, and rel07.1-uc003-rest-profile-boundary.yaml add per-family boundary use cases for the previously uncovered families.
- docs/specs/test-suites/test-rel07.1-profile-boundaries.yaml validates the three boundary use cases against the family SRD acceptance criteria.
- docs/specs/test-suites/test-rel07.1-profile-conformance.yaml records the executable per-family conformance suite that runs each family through the agent-core CLI and asserts on captured OpenTelemetry traces. Every case runs the shipped profile.yaml wrapper (or a named shipped variant), copying it only to patch harness-controlled fields; the generator and rest families additionally run their shipped variants and carry an ungated load-and-wiring companion. Coverage is annotated per family and per variant (shipped-behavioral vs load-and-wiring).

<!-- field: use_cases -->
### Use cases

```yaml
- id: rel07.1-uc001-specification-critic-profile-boundary
  historical_aliases: [rel07.1-uc001-jurist-profile-boundary]
  summary: Specification-critic deterministic, LLM-free charter validation boundary.
  status: specified
- id: rel07.1-uc002-lifecycle-profile-boundary
  summary: Lifecycle history, rollback, and approval machine-action boundary.
  status: specified
- id: rel07.1-uc003-rest-profile-boundary
  summary: REST client and server boundary words and the ollama-* variant.
  status: specified
```

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| test-rel07.1-profile-boundaries | Boundary use-case validation for the specification-critic, lifecycle, and rest families. | specified |
| test-rel07.1-profile-conformance | Executable per-family OpenTelemetry conformance for all ten families. Nine are implemented shipped-behavioral (runs the shipped profile.yaml/variant, copy-and-patch only harness-controlled fields); the rest sample server boundary is bounded-behavioral (shipped REST config and declarations over a synthesized launch/await/stop machine). Generator (with qwen27b/qwen35b variants), the planner behavioral run, and the rest ollama-profile variant are Ollama-gated; lifecycle suspend/resume and history/rollback are Dolt-gated; rest proves its inbound webhook server boundary over the shipped REST config (client head proven at the agent-core package level). Generator, planner, and rest carry ungated load-and-wiring companions. | implemented |

<!-- field: proof -->
### Proof

- mage validate
- mage conformance

<!-- field: deferred -->
### Deferred

- Profile-owned Mage integration targets that execute the specification-critic, lifecycle, and rest boundary fixtures, mirroring magefiles/integration_*.go for the covered families, and that trace to rel07.1-uc001/uc002/uc003 S1-S3. These are distinct from the delivered mage conformance suite, which asserts SRD acceptance criteria per family (srd005-specification-critic, srd009-lifecycle, srd007-rest) plus the role taxonomy, so the conformance cases are executable evidence at the SRD layer but do not by themselves satisfy the boundary use cases. The gap is concrete for rest, whose client read/create/probe head is not run in a profile machine (conformance proves only the inbound webhook server boundary; the shipped probe/delay retry grammar is proven by wiring conformance), and the lifecycle conformance evidence is Dolt-gated. Until the boundary integration targets that trace uc001/uc002/uc003 S1-S3 land, the three use cases and their test-rel07.1-profile-boundaries cases stay status specified rather than implemented; promoting them on conformance evidence alone would overstate rest and hide the Dolt gate.
- Running the full rest client create/probe head inside a profile machine (rel07.1-uc003), which needs inter-tool input shaping so REST client steps satisfy the declared-only runtime input contract. The rest conformance case proves the inbound webhook server boundary to Succeeded; readiness polling is explicit in the shipped MachineSpec and legacy await_operation is rejected. Planner conformance (#211) is delivered — agent-core gained a load_graph pipeline word that seeds the requirement graph before extraction. Lifecycle history/rollback conformance (#209) is delivered — agent-core routes the --request checkpoint id and to_iteration into a target-pinned checkpoint backend and the tool configs.

<!-- release: 08.0 | status: done -->
## 08.0 - Chroma corpus agents

Two profile-owned agents use a local Chroma vector database as a retrieval
boundary: an ingest profile that loads a repository corpus into a Chroma
collection, and a reader profile that answers a question grounded in the
retrieved chunks. Both are pure applications/catalog configuration over existing
agent-core words -- REST client words, list_resource, read_resource, and
invoke_llm. The ingest machine lists trusted corpus paths before model
control; the model selects listed files, and the machine sequences the
deterministic embed and add steps after each read (#251, #849).

Chroma does not embed text server-side over its raw HTTP API (spike #239),
so the agents compute embeddings at a local Ollama provider (an
/api/embeddings REST word) and thread the resulting vector into the Chroma
add and query bodies. Provider embedding, not server-side embedding, is the
design. The
threading uses the agent-core previous-Result parameter capability
delivered by #245 (srd028-rest-client-tools R12): body_source
previous_result selects the vector through input_mapping, carry_forward
echoes the source text and id, and a single-token array element reshapes
the flat vector into the nested array Chroma expects. No embedding vector is
authored by a model.

<!-- field: planned_work -->
### Planned work

| id | summary | status |
| --- | --- | --- |
| gh-230 | Chroma corpus agents epic -- srd012, release 08.0 docs, profiles, and integration test. | done |
| gh-240 | Add srd012-chroma-corpus-agents SRD. | done |
| gh-241 | Add release 08.0 roadmap, use cases, test suite, and SPECIFICATIONS entries. | done |
| gh-242 | Add the profile-owned Chroma and Ollama REST assets and the ingest profile. | done |
| gh-243 | Add the Chroma reader profile. | done |
| gh-244 | Add the mage integration:chroma target and tracer test. | done |
| gh-250 | Fix the runtime and reader defects found by running the agents end-to-end. | done |
| gh-251 | Redesign the ingest as a model-driven loop that actually loads the corpus. | done |

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| rel08.0-uc001-corpus-ingest | Ingest profile lists and reads a corpus under model control, embeds at Ollama, and threads vectors into Chroma add. | implemented |
| rel08.0-uc002-corpus-query | Reader profile embeds a question, threads the vector into Chroma query, and answers grounded in retrieved chunks. | implemented |

<!-- field: proof -->
### Proof

- mage validate
- cd applications/chatbot-mesh && mage integration:chroma
- go test ./conformance -run TestCorpusIngestListsTrustedCorpusBeforeModelControl
- go test ./conformance -run TestCorpusReader

<!-- field: success_markers -->
### Success markers

- The ingest and reader profiles are pure applications/catalog configuration over existing agent-core words.
- Every Chroma add and query carries an agent-supplied embedding vector threaded from the Ollama embed word, never a model-authored or server-side vector.
- The reader answer is grounded in chunks retrieved through embed_query -> chroma_query -> invoke_llm.
- The Chroma server is an external precondition surfaced through chroma_ready, and the embedding provider is the library the profile binds (srd058); neither is started by the agent.

<!-- release: 09.0 | status: implemented -->
## 09.0 - RAG server and single-RAG chatbot turn

Release 09.0 is the core loop of the chatbot mesh: a persistent RAG server
agent (srd013, vector-in) that answers a query_embeddings request with
chunks and embedding-model metadata, and a browser-facing chatbot agent
(srd014) whose request-scoped turn embeds the message once, queries a single
RAG, composes a grounded prompt through command-state $from(), and answers,
degrading to a mapped 200 when the RAG fails. Epics #292 (RAG server) and
#293 (chatbot core loop).

<!-- field: numbering_note -->
### Numbering note

Renumber decision (GH-332): the original single "09.0 Chatbot mesh" release
carried three use cases across five epics (#292-#296), and the observability
epic (#295) had no use case. Release 09 is broken into one release per epic so
implementation traces one-to-one, and the too-broad rel09.0-uc002-chatbot-turn
is split: its single-RAG reduction stays in 09.0 (the core-loop epic #293),
and its tier-selection/fan-out/degradation criteria move to 09.1 (the
tier-selection epic #294). The rel09.0 use cases were status planned when this
decision was made; the use-case and test-suite files are renamed to their
staged ids by the follow-on spec issue GH-333, which also authors the
observability (09.2), control-plane (09.4), and alignment (10.0) tracers. Each
release below maps epics to use cases to a test suite one-to-one-to-one.

<!-- field: relocated_to -->
### Relocated to

Relocated to applications/chatbot-mesh (rel00.0). The mesh data-plane specs (SRDs, use cases, test suites) now live with the standalone example; this entry remains for history. The control plane (09.4) is relocated too (rel05.0, GH-430).

<!-- field: srds -->
### Srds

- srd013-rag-server-agent
- srd014-chatbot-agent

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| rel09.0-uc001-rag-server-query | Persistent RAG server answers a vector-in query_embeddings request and returns chunks with embedding-model metadata. | implemented |
| rel09.0-uc002-chatbot-turn | Single-RAG chat turn embeds once, queries one RAG, composes via $from(), answers, and degrades to a 200 when the RAG fails. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| test-rel09.0-chatbot-mesh | Profile-owned tracers for the RAG server query and the single-RAG chatbot turn with a degraded RAG. | implemented |

<!-- release: 09.1 | status: implemented -->
## 09.1 - Tier selection and multi-RAG fan-out

Release 09.1 extends the chatbot turn (srd014 R2, R3) with the LLM $tool
tier selector that classifies a question to one chat-LLM word, and the sequential
multi-RAG fan-out that embeds once and queries N RAGs, each degraded
(CommandError) or embedding-mismatch-excluded (QueryRejected) source routing
on as a visible machine transition, compose assembling the surviving
per-source chunks. Epic #294.

<!-- field: relocated_to -->
### Relocated to

Relocated to applications/chatbot-mesh (rel01.0). The mesh data-plane specs (SRDs, use cases, test suites) now live with the standalone example; this entry remains for history. The control plane (09.4) is relocated too (rel05.0, GH-430).

<!-- field: srds -->
### Srds

- srd014-chatbot-agent

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| rel09.1-uc001-routed-fanout | A tier-selected chat turn dispatches one chat-LLM word via $tool and fans the one embedding out to two RAGs, citing both corpora and degrading when one fails. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| test-rel09.1-routing-fanout | Tracers for the $tool tier-selector dispatch, the two-RAG fan-out over disjoint corpora, and the degraded and embedding-mismatch paths. | implemented |

<!-- release: 09.2 | status: implemented -->
## 09.2 - Observability

Release 09.2 is the observability surface (srd014 R5.1, R5.3): a per-agent
monitor panel fed by each agent's monitor read_state and SSE stream through
the chatbot's same-origin monitor_proxy, and a cross-agent trace waterfall
that renders a turn's connected trace from the Jaeger query API, real
because the request-scoped machine exports its spans and propagates
traceparent to the RAG servers. Epic #295.

<!-- field: relocated_to -->
### Relocated to

Relocated to applications/chatbot-mesh (rel02.0). The mesh data-plane specs (SRDs, use cases, test suites) now live with the standalone example; this entry remains for history. The control plane (09.4) is relocated too (rel05.0, GH-430).

<!-- field: srds -->
### Srds

- srd014-chatbot-agent

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| rel09.2-uc001-observability | The observability panel shows each agent's run state and event stream and renders a selected turn's connected cross-agent trace as a span waterfall. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| test-rel09.2-observability | Tracers for the per-agent monitor panels over monitor_proxy and the cross-agent trace waterfall from the trace backend. | implemented |

<!-- release: 09.3 | status: implemented -->
## 09.3 - Mesh deployment and tiered swap

Release 09.3 deploys the mesh as one Helm chart (srd015) with values-driven
RAG pairs, co-generated chatbot client config and N-RAG fan-out from one
source, a tiered hot swap (repoint is a Service change; add-RAG re-renders
and rolls the chatbot with a Dolt checkpoint resume), OTLP observability,
and a provisioning deployment API. Epic #296.

<!-- field: relocated_to -->
### Relocated to

Relocated to applications/chatbot-mesh (rel03.0). The mesh data-plane specs (SRDs, use cases, test suites) now live with the standalone example; this entry remains for history. The control plane (09.4) is relocated too (rel05.0, GH-430).

<!-- field: srds -->
### Srds

- srd015-chatbot-deployment

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| rel09.3-uc001-mesh-deployment | One Helm chart deploys the mesh on a kind cluster with values-to-config co-generation, a Service repoint, and an add-RAG rollout with checkpoint resume, provisioned through a separately-authenticated deployment API. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| test-rel09.3-mesh-deployment | Tracers for the kind deploy smoke, the co-generation drift lock, the repoint and add-RAG swaps, and the provisioning deployment API. | implemented |

<!-- release: 09.4 | status: planned -->
## 09.4 - Control plane

Release 09.4 adds the mesh control plane: a Coordinator agent (srd016) that
sequences directory ingestion and reconfiguration and to which the UX agent
delegates provisioning intents, and a Creator agent (srd017) for agent
lifecycle management (template selection, instance creation, checkpoint-resume
restarts, health monitoring) acting on the coordinator's requests. The
provisioning and deployment swap surfaces of #296 land behind these roles.
Authored by GH-329.

<!-- field: relocated_to -->
### Relocated to

Relocated to applications/chatbot-mesh (rel05.0) as srd004-provisioning-workflow-orchestrator and srd005-creator, with the rel05.0 control-plane use case and test suite; this entry remains for history and the example is the canonical copy (GH-430, GH-1077).

<!-- field: srds -->
### Srds

- srd016-coordinator
- srd017-creator

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| rel09.4-uc001-control-plane-provisioning | The coordinator sequences a reconfiguration the UX agent delegates, and the creator materializes and health-monitors the agent instances, resuming from checkpoints on restart. | specified |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| test-rel09.4-control-plane | Tracers for the coordinator's reconfiguration sequencing and the creator's instance lifecycle and checkpoint-resume restart. | specified |

<!-- release: 09.5 | status: implemented -->
## 09.5 - In-cluster LLM tier

Release 09.5 promotes Ollama from an external values-driven endpoint to a
first-class, chart-owned LLM tier (srd015 R6): an optional Ollama subchart,
enabled by default, with PVC-backed model storage, a model-preload Job, agent
readiness gated on model presence, and GPU scheduling values. A fresh-cluster
install is self-contained; setting ollama.enabled=false with an endpoint
override reproduces the external-endpoint behavior, and the co-generation
contract (R2) is unchanged either way. This is a minor release validating the
completed rel09.0-09.3 mesh. Epic #335.

<!-- field: numbering_note -->
### Numbering note

Release 09.5 validates the completed rel09.0-09.3 mesh with an in-cluster LLM
tier (a minor release, no renumbering of the deployment releases). The epic
proposed rel09.1 for this, but GH-332 assigned rel09.1 to tier
selection/fan-out; 09.5 is the next free number in the chatbot-mesh 09.x line.

<!-- field: relocated_to -->
### Relocated to

Relocated to applications/chatbot-mesh (rel04.0). The mesh data-plane specs (SRDs, use cases, test suites) now live with the standalone example; this entry remains for history. The control plane (09.4) is relocated too (rel05.0, GH-430).

<!-- field: srds -->
### Srds

- srd015-chatbot-deployment

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| rel09.5-uc001-in-cluster-llm-tier | A fresh-cluster install with defaults answers a chat turn with zero external LLM dependency; the disabled-tier override reproduces the external endpoint, and agents are not ready until the models are pulled. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| test-rel09.5-llm-tier | Tracers for the default-install turn, the disabled-tier endpoint override, the preload-gate ordering, and the CPU-only kind smoke. | implemented |

<!-- release: 10.0 | status: implemented -->
## 10.0 - Reference-architecture role alignment

Release 10.0 aligns the agents/ role catalog with the autogenic-systems
reference architecture (roles per ITU-T Y.3502 / TM Forum IG1251): the
evaluator family is renamed to critic and the generator family to executor
(the reference role names they implement), the chroma family moves under
knowledge-manager as corpus-ingest and corpus-reader (a role, not a vendor
database), and the rest, control, and lifecycle profile directories -- which
are conformance grammar fixtures rather than roles -- move to
testdata/conformance beside testdata/integration. After this release agents/
reads as a role catalog. Epic #323 (issues #325-#328, #330).

<!-- field: proof_commands -->
### Proof commands

- mage validate  # profile wiring holds over the renamed tree
- mage conformance  # conformance fixtures pass under testdata/conformance
- go test ./...  # agent-core critic-session/critic-point builtin tool dirs
- grep -r 'agents/critic\|agents/executor\|agents/chroma\|testdata/conformance/rest/\|testdata/conformance/control/\|testdata/conformance/lifecycle/' -- returns no stale path outside retained historical release records

<!-- field: srds -->
### Srds

- srd013-rag-server-agent
- srd014-chatbot-agent

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| rel10.0-uc001-role-catalog-renames | The evaluator family is renamed to critic and generator to executor (reference role names). | implemented |
| rel10.0-uc002-knowledge-family-rescope | The vendor-named chroma family moves under knowledge-manager as corpus-ingest and corpus-reader. | implemented |
| rel10.0-uc003-conformance-fixture-relocation | The rest, control, and lifecycle conformance fixtures move to testdata/conformance so agents/ lists only role families. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| test-rel10.0-role-catalog | Cases validating the critic/executor renames, the chroma rescope, and the conformance-fixture relocation with mage validate, mage conformance, go test, and zero-hit greps. | implemented |

<!-- release: 11.0 | status: in-progress -->
## 11.0 - Scenario Critic test rig

Release 11.0 adds the repository's reusable test rig as two supported
test-time families. The scenario-critic (srd018) is the rig itself: a deterministic
machine that discovers scenarios under configurable roots, and for each
one starts the mocks its subject depends on, starts the
subject with its REST dependencies redirected at those mocks, runs the
scenario's validator machines, reduces their terminal states to a
verdict, and tears every child down. The mock (srd019) is one generic
profile that stands in for any REST dependency, serving the canned,
scripted responses a scenario fixture declares and recording a request
log validators assert against.

Both profiles are supported test-time library members. Their paths and
behavior are part of the canonical v0.* public release
surface even while the release remains in-progress on the evidence gaps
below. Earlier applications/catalog/v0.* tags and matching
agent-profiles/v0.* names are v0 compatibility aliases,
and immutable historical legacy tags remain valid. Rig-subject and other
conformance fixtures remain internal test data, not members.

GH-1072 renames the concrete rig from assembler to scenario-critic because
the actor binds Critic with Sandbox Validator and Output Evaluator; its
scenario composition is an actor function, not the canonical Assembler
responsibility. New consumers use agents/scenario-critic/profile.yaml.
The old agents/assembler/profile.yaml path is a non-forking wrapper through
the remainder of v0 and is removed at v1.

Tests are profile-owned data. A scenario is a tests/<scenario>/
directory beside the agent it tests, holding validator machines and
fixtures, so one subject carries several tests and adding one never
edits the rig. Discovery roots cover both the agent families here and
the agents under applications/, so one rig serves the whole repository. No
scenario needs a live service and none needs a container image per
test: every child is the same agent binary over a mounted profile.

The runtime words the rig calls are agent-core release 09.0 (srd039
mock endpoint binding, srd040 service-child tools). This release records
the profile-owned documentation; the mock lands in GH-668, the scenario critic
in GH-669, and the reference subject with its scenarios and the mage
target in GH-670. Epic #662.

The implementation issues are complete and the shipped profiles, fixtures,
conformance tests, and mage rig target exist. Release status remains
in-progress because current evidence does not execute validator assertions
over right and wrong request logs, exercise malformed-fixture serving
startup, inspect child-process survival, or induce mock-start-failure
teardown.

<!-- field: related_agent_core_release -->
### Related agent core release

09.0

<!-- field: planned_work -->
### Planned work

| id | summary | status |
| --- | --- | --- |
| gh-665 | Specify the scenario critic and mock profile families (this release's documentation). | done |
| gh-668 | Add the mock profile and fixture conformance. | done |
| gh-669 | Add the scenario critic profile (the rig machine). | done |
| gh-670 | Add the reference subject, self-proving scenarios, and the mage rig target. | done |
| gh-678 | Add the scenario-session layer in agent-core (split from gh-669). | done |
| gh-1072 | Rename the concrete test rig to scenario-critic and retain a time-bounded compatibility wrapper. | done |

<!-- field: historical_aliases -->
### Historical aliases

```yaml
srds: [srd018-assembler]
use_cases: [rel11.0-uc001-assembler-scenario-run]
test_suites: [test-rel11.0-assembler-rig]
```

<!-- field: srds -->
### Srds

- srd018-scenario-critic
- srd019-mock

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| rel11.0-uc001-scenario-critic-scenario-run | The scenario critic runs an agent's scenarios with no live services, reports ordered pass/fail verdicts, exposes machine-visible child-stop dispatches, and repeats without observed port or state leakage. | in-progress |
| rel11.0-uc002-mock-serves-fixtures | One mock profile stands in for any REST dependency, scripts response sequences, and exposes a per-instance request log whose contents are inspected by conformance tests. | in-progress |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| test-rel11.0-scenario-critic-rig | Cases for fixture serving and sequencing, two mocks from one profile, direct request-log inspection, malformed-fixture --validate-config rejection, scenario verdicts, visible stop dispatches, repeatable runs, cross-root discovery, and shipped-profile validation. | in-progress |

<!-- field: proof -->
### Proof

- mage validate
- mage rig

<!-- field: success_markers -->
### Success markers

```yaml
- Adding a test means adding files under a tests/ directory, never editing the scenario critic.
- One mock profile serves every dependency; response content lives in scenario fixtures.
- Scenarios run with no live Chroma, Ollama, Docker, or Kubernetes.
- The rig can fail: a wrong expectation reaches Failed and names the failing validator.
```

<!-- release: 12.0 | status: release-ready -->
## 12.0 - Agent role-realization naming alignment

Release 12.0 completes the naming migration without rewriting the
incomplete Release 07.1 or 11.0 histories. Canonical roles remain shared
interfaces; catalog profiles and application actors use concrete
realization names, explicit bindings, and one canonical implementation
home.

Reusable old catalog paths are one-file, non-forking wrappers through
the remainder of v0 and are removed at v1. Application-local test
identities and the chatbot-mesh coordinator workload migrated atomically
and retain historical entries without runtime aliases. The exact planned
release is applications/catalog/v0.20260727.0 with matching
agent-profiles/v0.20260727.0 compatibility tag.

Executable evidence is implemented and this release is release-ready.
Immutable tags are intentionally not published from this feature branch;
after merge, mage tag must run from a clean main checkout and create the
root, canonical catalog, and legacy alias tags atomically.

<!-- field: dependencies -->
### Dependencies

| id | status | effect |
| --- | --- | --- |
| catalog-release-07.1 | in-progress | Historical boundary evidence remains independently incomplete; no status is promoted here. |
| catalog-release-11.0 | in-progress | Historical rig evidence remains independently incomplete; renamed IDs and aliases preserve traceability. |

<!-- field: planned_work -->
### Planned work

| id | summary | status |
| --- | --- | --- |
| gh-1076 | Define role-realization semantic authority | done |
| gh-1069 | Publish the application composition pattern language | done |
| gh-1074 | Enforce binding and classification governance | done |
| gh-1072 | Rename assembler to scenario-critic | done |
| gh-1075 | Rename monitor to runtime-state-reader | done |
| gh-1068 | Rename jurist to specification-critic | done |
| gh-1070 | Rename mesh validators as Critic realizations | done |
| gh-1073 | Classify chatbot-mesh actors | done |
| gh-1077 | Rename mesh coordinator realization | done |
| gh-1071 | Coordinate aliases, pins, package closures, and release evidence | done |

<!-- field: migration_note -->
### Migration note

docs/migrations/v0.20260727.0-agent-role-realization-alignment.yaml

<!-- field: use_cases -->
### Use cases

| id | summary | status |
| --- | --- | --- |
| rel12.0-uc001-agent-role-realization-alignment | Coordinate all realization renames as one versioned compatibility migration. | implemented |

<!-- field: test_suites -->
### Test suites

| id | summary | status |
| --- | --- | --- |
| test-rel12.0-agent-role-realization-alignment | Release metadata, aliases, closures, package renders, tag semantics, and zero-hit evidence. | implemented |

<!-- field: proof -->
### Proof

- mage validate
- mage conformance
- mage integration:all
- cd ../.. && mage audit
- cd ../.. && mage test
- cd ../.. && mage stats
- cd ../coding-agent && mage audit
- cd ../chatbot-mesh && mage helm:package

<!-- field: post_merge -->
### Post merge

- Run mage tag from a clean main checkout after all release gates pass.
- Push v0.20260727.0, applications/catalog/v0.20260727.0, and agent-profiles/v0.20260727.0 together.

<!-- release: 13.0 | status: implemented -->
## 13.0

<!-- field: title -->
### Title

Collector family promotion

<!-- field: dependencies -->
### Dependencies

- agent-core release 12.0 (spool query builtins)

<!-- field: planned_work -->
### Planned work

```yaml
- {gh: 1145, title: "Promote the collector agent to a canonical catalog family", estimated_loc: "300 YAML, 250 test"}
```

<!-- field: specification -->
### Specification

```yaml
srd: srd020-collector
status: implemented
```

<!-- field: use_cases -->
### Use cases

| id | title | status |
| --- | --- | --- |
| rel13.0-uc001-collector-family | A standalone collector receives, spools, and shuts down in spool mode | implemented |

<!-- field: test_suites -->
### Test suites

```yaml
- id: test-rel13.0-collector
  title: Collector family conformance
  cases: 5
  status: implemented
```

<!-- field: proof -->
### Proof

- mage validate
- mage conformance
- cd ../.. && mage audit

<!-- release: 14.0 | status: implemented -->
## 14.0

<!-- field: title -->
### Title

Served collector trace UI

<!-- field: dependencies -->
### Dependencies

- catalog release 13.0 (collector family and shipped ui/)

<!-- field: planned_work -->
### Planned work

```yaml
- {gh: 1236, title: "Specify the served collector trace UI", estimated_loc: "0 (documentation)"}
- {gh: 1237, title: "Collector serves ui/dist via static_assets on the query listener", estimated_loc: "50 YAML, 150 test"}
- {gh: 1254, title: "Serve the collector UI from the coding-agent chart", estimated_loc: "120"}
- {gh: 1255, title: "Serve the collector UI from the agent-architecture chart", estimated_loc: "120"}
- {gh: 1256, title: "Serve the collector UI from the chatbot-mesh chart (dedicated ConfigMap, writable metrics spool)", estimated_loc: "160"}
```

<!-- field: specification -->
### Specification

```yaml
srd: srd020-collector
status: implemented
```

<!-- field: use_cases -->
### Use cases

| id | title | status |
| --- | --- | --- |
| rel14.0-uc001-served-trace-ui | The collector serves its own trace UI at the query origin | implemented |

<!-- field: test_suites -->
### Test suites

```yaml
- id: test-rel14.0-served-ui
  title: Served collector trace UI
  cases: 5
  status: implemented
```

<!-- field: proof -->
### Proof

- mage validate
- mage conformance
- cd ../.. && mage audit

<!-- release: 15.0 | status: implemented -->
## 15.0

<!-- field: title -->
### Title

Collector Explore panel

<!-- field: related_agent_core_release -->
### Related agent core release

13.0

<!-- field: dependencies -->
### Dependencies

- agent-core release 13.0 (span-analytics builtins)
- catalog release 14.0 (served collector UI)

<!-- field: planned_work -->
### Planned work

```yaml
- {gh: 1239, title: "Specify span-analytics tools and the collector Explore surface", estimated_loc: "0 (documentation)"}
- {gh: 1241, title: "Collector Explore query routes: /query/spans/stats and /query/spans/breakdown", estimated_loc: "200 YAML, 250 test"}
- {gh: 1242, title: "Collector UI Explore page with heatmap drag-select and TraceDetail links", estimated_loc: "500 TS/TSX"}
```

<!-- field: specification -->
### Specification

```yaml
srd: srd020-collector
status: implemented
```

<!-- field: use_cases -->
### Use cases

| id | title | status |
| --- | --- | --- |
| rel15.0-uc001-explore-panel | An operator narrows a latency band to a responsible attribute and lands on a trace | implemented |

<!-- field: test_suites -->
### Test suites

```yaml
- id: test-rel15.0-explore
  title: Collector Explore panel
  cases: 6
  status: implemented
```

<!-- field: proof -->
### Proof

- mage validate
- mage conformance
- cd ../.. && mage uiDist
- cd ../.. && mage audit

<!-- release: 16.0 | status: implemented -->
## 16.0

<!-- field: title -->
### Title

Canonical applier family

<!-- field: dependencies -->
### Dependencies

- applications Release 14.0 application consistency contract
- manifest-derived closure and normalized application actor trees

<!-- field: planned_work -->
### Planned work

```yaml
- {gh: 1420, title: "Promote one canonical applier family and migrate application wrappers", estimated_loc: "600 canonical production, 500 test"}
```

<!-- field: specification -->
### Specification

```yaml
srd: srd022-applier
status: implemented
```

<!-- field: migration -->
### Migration

docs/migrations/v0.20260804.0-applier-family.yaml

<!-- field: use_cases -->
### Use cases

| id | title | status |
| --- | --- | --- |
| rel16.0-uc001-applier-family | Three applications compose one canonical applier family | implemented |

<!-- field: test_suites -->
### Test suites

| id | title | status |
| --- | --- | --- |
| test-rel16.0-applier | Canonical applier family conformance | implemented |

<!-- field: proof -->
### Proof

- mage validate
- mage conformance
- go test ./conformance -run Applier
- cd ../.. && mage audit

<!-- field: success_markers -->
### Success markers

- Lifecycle, apply, and rollout-read machines have one catalog owner.
- Application wrappers retain REST and deployment authority.
- Package provenance distinguishes canonical and local closure members.
- Duplicate and statistics audits prevent copied forks and double counting.

<!-- release: 17.0 | status: implemented -->
## 17.0

<!-- field: title -->
### Title

Rig doctor root-cause analyst

<!-- field: dependencies -->
### Dependencies

- Release 13.0 collector family (spool trace query words)
- repository-root kind rig failure-evidence capture (eng01)

<!-- field: planned_work -->
### Planned work

```yaml
- {gh: 2286, title: "Specify the rig-doctor: srd023, use case, test suite, eng01 diagnose verb", estimated_loc: "0 (documentation)"}
- {gh: 2287, title: "Rig evidence capture with manifest and mage diagnose (capture only)", estimated_loc: "180 production"}
- {gh: 2288, title: "rig-doctor profile on root-cause-analyst, registration, and diagnose wiring", estimated_loc: "350 production and test"}
```

<!-- field: specification -->
### Specification

```yaml
srd: srd023-rig-doctor
status: implemented
```

<!-- field: use_cases -->
### Use cases

| id | title | status |
| --- | --- | --- |
| rel17.0-uc001-rig-doctor | An operator runs mage diagnose after a failed scenario and reads a cited diagnosis | implemented |

<!-- field: test_suites -->
### Test suites

```yaml
- id: test-rel17.0-rig-doctor
  title: Rig doctor root-cause analyst conformance
  cases: 12
  status: implemented
```

<!-- field: proof -->
### Proof

- mage validate
- mage conformance
- go test ./conformance -run RigDoctor
- cd ../.. && mage test
- cd ../.. && mage audit

<!-- field: success_markers -->
### Success markers

- Capture records a manifest of every file and every capture error.
- The diagnosis loop declares no exec, kubectl, or helm word.
- Every diagnosis claim cites a captured file and line.
- A diagnosis that cannot be read back is repaired rather than reported.

<!-- release: 18.0 | status: implemented -->
## 18.0

<!-- field: title -->
### Title

Host-side deployer

<!-- field: dependencies -->
### Dependencies

- Release 16.0 canonical applier family
- Release 17.0 rig-doctor mage-target agent harness

<!-- field: planned_work -->
### Planned work

```yaml
- {gh: 2300, title: "Specify the host-side deployer: srd022 R6, rel18.0 use case, test suite, road map", estimated_loc: "0 (documentation)"}
- {gh: 2301, title: "Catalog applier deploy and undeploy machines with standalone conformance", estimated_loc: "200 YAML and 260 test"}
- {gh: 2302, title: "Host-side deploy declaration rendering and the kindrig deploy harness", estimated_loc: "250 production and 240 test"}
- {gh: 2304, title: "Migrate coding-agent demo:up to mage deploy and add mage undeploy", estimated_loc: "180 production and 140 test"}
```

<!-- field: specification -->
### Specification

```yaml
srd: srd022-applier
status: implemented
```

<!-- field: use_cases -->
### Use cases

| id | title | status |
| --- | --- | --- |
| rel18.0-uc001-host-side-deploy | An operator runs mage deploy and the applier's deploy machine installs the release | implemented |

<!-- field: test_suites -->
### Test suites

```yaml
- id: test-rel18.0-host-side-deploy
  title: Host-side deployer conformance
  cases: 11
  status: implemented
```

<!-- field: proof -->
### Proof

- mage validate
- mage conformance
- go test ./conformance -run Applier
- cd ../.. && mage test
- cd ../.. && mage audit

<!-- field: success_markers -->
### Success markers

- The deploy path declares no model word, so a deploy costs no tokens.
- Day-0 and day-2 differ in the rendered words alone, never in the transition graph.
- No rendered argv carries a placeholder coordinate or a --reuse-values flag.
- A first install that fails leaves no release behind.
- The deploy target gates its build on the machine's terminal state.
- A busy provider costs two retries, not the evidence already read.

<!-- release: 99.0 | status: planned -->
## 99.0

<!-- field: title -->
### Title

Observer family promotion

<!-- field: dependencies -->
### Dependencies

- chatbot-mesh release 08.0 (the fleet observer this family generalizes)

<!-- field: planned_work -->
### Planned work

```yaml
- {gh: 1247, title: "Specify the cross-application fleet observer (catalog promotion + per-application adoption)", estimated_loc: "0 (documentation)"}
```

<!-- field: specification -->
### Specification

```yaml
srd: srd021-observer
status: planned
```

<!-- field: use_cases -->
### Use cases

| id | title | status |
| --- | --- | --- |
| rel99.0-uc001-observer-family | The observer is a catalog-owned family adopted across applications | planned |

<!-- field: test_suites -->
### Test suites

```yaml
- id: test-rel99.0-observer
  title: Observer family promotion
  cases: 3
  status: planned
```

<!-- field: proof -->
### Proof

- cd ../.. && mage audit
