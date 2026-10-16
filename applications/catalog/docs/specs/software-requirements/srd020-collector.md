---
id: srd020-collector
title: Collector Agent Family
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

The chatbot-mesh owns a collector agent that receives, spools, and optionally
relays OTLP trace batches. This profile is independently useful outside the
mesh — any deployment that needs an OTLP trace collector with a declarative
state machine and spool-first evidence model can use it. Promoting it to the
catalog gives it a canonical home, a generalized default configuration, and
conformance coverage.

<!-- field: goals -->
## Goals

- G1: Provide a canonical collector family in the catalog generalized from the mesh profile.
- G2: Default to spool mode so a standalone collector works without a relay endpoint.
- G3: Cover the family with deterministic conformance tests exercising the agent-core OTLP builtins.
- G4: Serve the shipped trace UI from the collector itself, so every collector deployment browses its own spool without a separate UI workload.
- G5: Expose an Explore surface -- a duration heatmap, group-by counts, and attribute divergence -- so an operator drills from a latency region to the responsible attribute and into a trace waterfall.

<!-- req: R1 -->
## R1 - Generalized defaults

<!-- req: R1.1 -->
### R1.1

The collector family SHALL default COLLECTOR_MODE to spool, not relay.

<!-- req: R1.2 -->
### R1.2

The relay endpoint SHALL have no default; relay mode requires COLLECTOR_RELAY_ENDPOINT set.

<!-- req: R1.3 -->
### R1.3

The decide_collector_relay predicate SHALL compare COLLECTOR_MODE against relay, defaulting to spool when unset.

<!-- req: R2 -->
## R2 - Surface separation

<!-- req: R2.1 -->
### R2.1

The monitor surface SHALL serve lifecycle and machine state reads only; no monitor view SHALL carry span payloads.

<!-- req: R2.2 -->
### R2.2

A separate query surface, collector_query, SHALL serve bounded spool reads through spool_list_traces and spool_get_trace.

<!-- req: R3 -->
## R3 - Conformance coverage

<!-- req: R3.1 -->
### R3.1

The family SHALL enter the membership contract with a bindings.catalog entry in agent-role-realizations.md.

<!-- req: R3.2 -->
### R3.2

A per-family conformance test SHALL exercise otlp_receiver_launch, await_spans, spool_collector_spans, and the shutdown sequence.

<!-- req: R4 -->
## R4 - Role-realization binding

<!-- req: R4.1 -->
### R4.1

The bindings.catalog entry SHALL have classification role_realization, primary_role monitor, sub_roles [telemetry-collector], domain observability.

<!-- req: R4.2 -->
### R4.2

The mesh binding SHALL remain until the mesh swap issue replaces it with a catalog reference.

<!-- req: R5 -->
## R5 - Query REST surface

<!-- req: R5.1 -->
### R5.1

A collector_query server SHALL serve GET /query/traces for paginated trace summaries via spool_list_traces.

<!-- req: R5.2 -->
### R5.2

A collector_query server SHALL serve GET /query/traces/{trace_id} for single-trace span detail via spool_get_trace.

<!-- req: R5.3 -->
### R5.3

Page size SHALL be capped in tool config, not trusted from the HTTP request.

<!-- req: R5.4 -->
### R5.4

The query server lifecycle SHALL follow the same launch/stop pattern as control and monitor, with shutdown ordering after receiver stop and before monitor stop.

<!-- req: R6 -->
## R6 - Trace query UI

<!-- req: R6.1 -->
### R6.1

The collector family SHALL ship a ui/ directory with a React SPA rendering trace list and per-trace waterfall from the collector_query surface, using the UI kit TraceView read same-origin (applications/docs srd004 R2.4) with /traces/{trace_id} as the deep link of an open trace.

<!-- req: R6.2 -->
### R6.2

The UI SHALL import the canonical design tokens from the UI kit (applications/docs srd004 R8.2) and join the root uiDist reproducibility gate.

<!-- req: R7 -->
## R7 - Served UI

<!-- req: R7.1 -->
### R7.1

The collector_query server SHALL serve the shipped ui/dist tree through a static_assets endpoint with spa true, so the UI is reachable at the same origin as /query/*.

<!-- req: R7.2 -->
### R7.2

The UI index SHALL answer GET / and SPA deep links (for example /traces/{trace_id}) SHALL resolve to the index document so the client router renders the requested view.

<!-- req: R7.3 -->
### R7.3

The /query/* API endpoints SHALL take precedence over the static_assets catch-all; serving the UI SHALL NOT alter any /query/* response body or status.

<!-- req: R7.4 -->
### R7.4

The static_assets root SHALL be a literal path declared in the collector REST configuration, resolved relative to the profile directory; it SHALL NOT be introduced as a new environment variable.

<!-- req: R7.5 -->
### R7.5

When the configured UI root is absent, the server SHALL answer with an explicit HTTP error status, not a silent empty body.

<!-- req: R8 -->
## R8 - Explore query surface

<!-- req: R8.1 -->
### R8.1

The collector_query server SHALL serve GET /query/spans/stats through declared query-machine actions for agent-core spool_span_heatmap and spool_span_group_by followed by a compose word, returning a duration-over-time heatmap, group-by counts, and a capped list of exemplar trace IDs for the matched span set.

<!-- req: R8.2 -->
### R8.2

The collector_query server SHALL serve GET /query/spans/breakdown, bound to spool_span_breakdown, returning the ranked attribute divergence between a selected span region and its complement, and a capped list of exemplar trace IDs for the selection.

<!-- req: R8.3 -->
### R8.3

Both routes SHALL accept the filter, bucket, group-by, selection, and top-N inputs as request query parameters passed to the tool config; page and result caps SHALL come from tool config, not be trusted from the request. No route behavior SHALL depend on an environment variable.

<!-- req: R8.4 -->
### R8.4

The stats and breakdown JSON response contracts SHALL be pinned by conformance tests so UI and tool drift fail there first, following the R5 query-contract precedent.

<!-- req: R8.5 -->
### R8.5

An empty or absent spool SHALL yield a well-formed empty stats or breakdown response, not an error.

<!-- req: R9 -->
## R9 - Explore UI

<!-- req: R9.1 -->
### R9.1

The shipped UI SHALL add an Explore page that renders the duration heatmap from /query/spans/stats and lets the operator drag-select a region of it.

<!-- req: R9.2 -->
### R9.2

The Explore page SHALL show the ranked attribute divergence for the selected region from /query/spans/breakdown, and list the selection-scoped exemplar traces the breakdown (or, before any selection, the stats) response returns, each linking into the existing per-trace waterfall (R6).

<!-- req: R9.3 -->
### R9.3

The Explore page SHALL join the root uiDist reproducibility gate and reach the server through same-origin /query/* only, adding no new runtime dependency for the heatmap.

<!-- req: R10 -->
## R10 - Durable object-backed query

<!-- req: R10.1 -->
### R10.1

The collector spool and query words SHALL select a storage backend through declared configuration -- the filesystem NDJSON spool (default) or the object backend, which reads the collector's configured application bucket and in-bucket prefix (srd042 R10, R11). The bucket, endpoint, and prefix SHALL be declared configuration only; no request field SHALL redirect a query to another bucket.

<!-- req: R10.2 -->
### R10.2

A query response SHALL report storage status as complete, partial, or unavailable. A bucket outage SHALL NOT present locally staged (WAL) data as complete history.

<!-- req: R10.3 -->
### R10.3

After app:down the collector's query REST surface is unavailable while its bucket data remains; a collector redeployed against the same bucket SHALL restore query access to the retained history without changing object identity.

<!-- req: R10.4 -->
### R10.4

Collector readiness SHALL reflect OTLP intake and lifecycle health only. A remote object-store outage SHALL NOT flip a WAL-backed collector unready; instead the query surface SHALL expose the degraded storage_status (R10.2) so the collector never presents WAL-only history as complete. The chart readiness probe target SHALL be identical for the filesystem and object backends (GH-2484 R6).

<!-- req: R10.5 -->
### R10.5

The trace UI SHALL display complete, partial, unavailable, or unknown storage status beside the trace view. Partial and unavailable states SHALL state that displayed results are not complete retained history; a UI SHALL NOT silently collapse them into an empty or healthy trace list.

<!-- req: R10.6 -->
### R10.6

Object-backed collector declarations SHALL stamp a retention class, optional retain-until derived from non-negative retention days, and fixed app:purge authority on every immutable envelope. Retention metadata SHALL be descriptive only; app:down, collector restart, and platform lifecycle SHALL NOT delete an application retained set.

<!-- req: R11 -->
## R11 - Explicit application-data purge

<!-- req: R11.1 -->
### R11.1

A separate model-free application-purge profile SHALL require the exact token purge:<application> before its first mutating boundary call. Its request SHALL NOT select bucket URL, endpoint, prefix, credential, or application identity.

<!-- req: R11.2 -->
### R11.2

The purge boundary SHALL use only apprig-rendered storage coordinates, list a bounded retained set, refuse an over-budget set without mutation, delete only that application's configured prefix, and verify zero remaining before reporting Purged.

<!-- req: R11.3 -->
### R11.3

Purged, Refused, Incomplete, and Failed SHALL be distinct terminal outcomes. Every bounded destructive attempt SHALL append a non-secret fsynced audit record outside the target bucket.

<!-- req: R11.4 -->
### R11.4

Collector query profiles SHALL remain read-only. app:down, collector restart, and platform lifecycle SHALL never invoke purge.

<!-- field: non_goals -->
## Non-goals

- Removing the mesh collector or changing the mesh deployment to reference the catalog family.
- Metrics ingress or Dolt scraping.
- A permanent platform-wide query service or cross-application telemetry aggregation; each collector queries only its own application bucket.
- Mutation through the collector query profile; explicit application-data deletion belongs only to the approved application-purge profile run by apprig.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

The catalog conformance suite passes with the collector family; the profile launches standalone in spool mode and reaches the AwaitingSpans loop.

<!-- field: traces -->
#### Traces

- R1.1
- R1.3
- R3.2

<!-- ac: AC2 -->
### AC2

The role-realization audit passes with both the catalog and mesh bindings present, which is the membership contract the family enters.

<!-- field: traces -->
#### Traces

- R3.1
- R4.1
- R4.2

<!-- ac: AC3 -->
### AC3

The relay endpoint config has no default; the decide predicate defaults COLLECTOR_MODE to spool.

<!-- field: traces -->
#### Traces

- R1.2
- R1.3

<!-- ac: AC4 -->
### AC4

The collector_query server serves the UI index at / and resolves SPA deep links to that index, at the same origin as the query API, and no /query/* response changes.

<!-- field: traces -->
#### Traces

- R7.1
- R7.2
- R7.3

<!-- ac: AC5 -->
### AC5

The UI root is a literal profile-declared path with no new environment variable, and a missing UI root yields an explicit error status rather than a blank page.

<!-- field: traces -->
#### Traces

- R7.4
- R7.5

<!-- ac: AC6 -->
### AC6

GET /query/spans/stats returns the heatmap, group-by, and exemplar-trace-ID contract and GET /query/spans/breakdown returns the ranked divergence and exemplar-trace-ID contract, both pinned by conformance, with caps taken from tool config not the request.

<!-- field: traces -->
#### Traces

- R8.1
- R8.2
- R8.3
- R8.4

<!-- ac: AC7 -->
### AC7

An empty or absent spool yields a well-formed empty stats and breakdown response rather than an error.

<!-- field: traces -->
#### Traces

- R8.5

<!-- ac: AC8 -->
### AC8

The Explore page renders the heatmap, drag-selects a region, shows the ranked divergence and the selection-scoped exemplar traces linking into the waterfall, and reproduces from a clean uiDist build.

<!-- field: traces -->
#### Traces

- R9.1
- R9.2
- R9.3

<!-- ac: AC9 -->
### AC9

Every monitor view answers without carrying span content, while the query surface serves the same spool as paginated summaries and per-trace detail whose page size is clamped to the configured cap rather than the request.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R5.1
- R5.2
- R5.3
- R5.4

<!-- ac: AC10 -->
### AC10

The shipped ui/ tree renders the trace list and per-trace waterfall from the query surface, imports the canonical design tokens, and rebuilds byte-identically under the root uiDist gate.

<!-- field: traces -->
#### Traces

- R6.1
- R6.2

<!-- ac: AC11 -->
### AC11

With the object backend selected, the collector spool persists immutable objects to the configured application bucket with declared retention class, optional retain-until, and fixed app:purge authority, and the query words read that bucket through the srd042 R11 engine, reporting storage status and never reading a request-supplied bucket; with the filesystem default the render and behavior are unchanged. The collector profile boot-smokes with the object-backed declarations, app:down preserves the objects, and a redeployed collector queries retained history.

<!-- field: traces -->
#### Traces

- R10.1
- R10.2
- R10.3
- R10.6

<!-- ac: AC12 -->
### AC12

Every query REST response body carries storage_status from the srd042 R11 engine, so a bucket outage reports partial or unavailable rather than presenting WAL-only history as complete; the trace UI renders that status and names incomplete history; the chart readiness and liveness probes target the control lifecycle health surface identically for the filesystem and object backends, so a transient object-store outage never flips a WAL-backed collector unready. Conformance pins the REST/UI status contract, render tests pin both probe targets, and the boot-smoke covers response mapping.

<!-- field: traces -->
#### Traces

- R10.2
- R10.4
- R10.5

<!-- ac: AC13 -->
### AC13

A mismatched confirmation reaches Refused without a storage call; exact confirmation runs the model-free purge profile over declaration-bound coordinates, deletes and verifies only the selected application's bounded retained prefix, reaches Purged, and writes an external audit with zero remaining. Live fake-GCS evidence proves peer application keys remain, and app:down never invokes this path.

<!-- field: traces -->
#### Traces

- R11.1
- R11.2
- R11.3
- R11.4

<!-- field: related_documents -->
## Related documents

- srd007-collector-agent
- cfg-tool-authoring-standard

<!-- field: implementation_status -->
## Implementation status

```yaml
implemented:
    - "applications/catalog/agents/collector/ owns the canonical collector family."
    - "applications/catalog/conformance/collector_test.go proves receive-spool-exit in spool mode."
    - "applications/docs/specs/semantic-models/agent-role-realizations.md carries the catalog binding."
    - "R10: the collector spool and query declarations select an object backend through env-templated storage config (filesystem default), so the spool persists to and the query reads the configured application bucket via the srd042 R10/R11 engine; the profile boot-smokes with these declarations (GH-2484, GH-2489, GH-2490)."
    - "R10.2/R10.4: every /query/* REST response surfaces storage_status, and the shared agent-services chart readiness/liveness probes target control lifecycle health identically for both backends, pinned by chatbot-mesh render tests, so a bucket outage never flips a WAL-backed collector unready (GH-2484 R6, GH-2491 R6)."
    - "R10.5: the shipped trace UI performs a bounded storage-status read and renders complete, partial, unavailable, or unknown with explicit incomplete-history language; catalog conformance pins the source/API contract (GH-2490)."
    - "R10.6: object envelopes stamp retention_class, optional retain_until, and purge_authority=app:purge from the versioned chart/declaration contract; no teardown path consumes that metadata (GH-2490)."
    - "R11: the approved model-free application-purge profile requires exact application confirmation, invokes internal object_delete_prefix over rendered authority, verifies zero remaining, and writes an external audit. Unit, profile conformance, and live fake-GCS isolation proofs cover refusal and success (GH-2488, GH-2490, GH-2493)."
    - "Cloud application-bucket IAM is pinned by the gcprig three-application policy fixture; authenticated GKE denial remains dependency-gated where gcloud credentials are unavailable (GH-2505)."
```

<!-- field: platform_dependencies_note -->
## Platform dependencies note

These SRDs live in the agent-core corpus, not this one; they are platform dependencies rather than in-corpus depends_on (GH-2310).

<!-- field: platform_dependencies -->
## Platform dependencies

```yaml
- srd_id: srd042-otlp-receiver-tools
  provider: agent-core
  symbols_used: [otlp_receiver_launch, await_spans, spool_spans, relay_spans, otlp_receiver_stop, durable object-bucket spool storage (R10), object-backed query read with storage status (R11)]
- srd_id: srd044-spool-query-tools
  provider: agent-core
  symbols_used: [spool_list_traces, spool_get_trace]
- srd_id: srd059-objectstore-tools
  provider: agent-core
  symbols_used: [declaration-bound object_delete_prefix with bounded verification and external audit (R5)]
- srd_id: srd045-span-analytics-tools
  provider: agent-core
  symbols_used: [spool_span_heatmap, spool_span_group_by, spool_span_breakdown]
```
