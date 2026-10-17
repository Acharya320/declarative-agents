---
id: srd007-collector-agent
title: Collector Agent Functional Block
depends_on:
  - srd_id: srd003-chatbot-deployment
    symbols_used: [R5]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: role_realization_decision | form: "tree" -->
## Role realization decision

- current_actor: collector
- implemented_responsibility: Receive bounded OTLP trace batches, apply declared policy, spool evidence, optionally relay unchanged traces, and expose lifecycle state.
- canonical_role_sub_role: Monitor / Telemetry Collector.
- retained_or_proposed_name: Retain `collector` as a compatible descriptive actor.
- authority_boundary: Owns trace and OTLP metric intake (GH-1207); the in-cluster Dolt Prometheus scrape is retired (GH-1366), contrib is an opt-in trace-path option only, and file-only self-telemetry prevents a loop.
- evidence:
  1. agents/collector/machine.yaml
  1. agents/collector/declarations.yaml
  1. helm/templates/collector.yaml
- migration_impact: Agent mode is the default and deploys no contrib collector or collector-metrics infrastructure; contrib remains selectable for the trace path only.

<!-- field: context -->
## Context

The `collector` actor assumes Monitor / Telemetry Collector and implements the mesh's
OTLP trace-ingress contract with agent-core srd042 words. It replaces the
contrib collector on the trace path while MachineSpec owns receive, inspect,
spool, relay, retry, control, and monitor behavior.

The migration is complete. The collector agent now owns both OTLP trace and
OTLP metric intake (GH-1207) and exposes the metric query surface the
observability evidence reads. Agent mode is the default and deploys no contrib
collector; the pull-based in-cluster Dolt Prometheus scrape it once provided is
retired (GH-1366). Contrib remains selectable for the trace path only.

<!-- field: problem -->
## Problem

The chart's third-party collector hides trace intake and routing outside the
agent grammar. The mesh needs its own trace ingress as a declarative agent,
without losing persistent run identity, deterministic file evidence,
upstream Jaeger traces, or the metrics evidence the current collector owns.

<!-- field: goals -->
## Goals

- G1: Run OTLP trace intake, spool, and relay as a persistent profile over agent-core srd042.
- G2: Preserve repository, module, target, commit, run, service, trace, span, and parent identity.
- G3: Expose control and monitor surfaces using the standard persistent-agent pattern.
- G4: Support deterministic test mode and upstream relay mode from one profile.
- G5: Prevent collector telemetry from feeding its own receiver.

<!-- req: R1 -->
## R1 - Persistent trace intake

<!-- req: R1.1 -->
### R1.1

The collector SHALL launch one named OTLP/gRPC receiver on a values-configured address with a bounded queue and declared overflow policy.

<!-- req: R1.2 -->
### R1.2

Its persistent machine SHALL loop through await, policy, spool, optional relay, and await again; these SHALL remain separate visible words.

<!-- req: R1.3 -->
### R1.3

The profile SHALL accept traces only; metrics and logs SHALL NOT be registered as collector-agent receiver words.

<!-- req: R1.4 -->
### R1.4

Queue overflow, await timeout, stop, spool failure, and relay failure SHALL route through explicit machine transitions.

<!-- req: R2 -->
## R2 - Identity and deterministic spool

<!-- req: R2.1 -->
### R2.1

The collector SHALL preserve every received resource span, scope span, trace id, span id, parent id, timestamp, event, link, status, and attribute.

<!-- req: R2.2 -->
### R2.2

Repository, module, target or scenario, git commit, unique run id, and stable service.name identity from srd003 R5.6 SHALL remain addressable after await and unchanged after relay.

<!-- req: R2.3 -->
### R2.3

Accepted batches SHALL append to a values-configured stdouttrace-compatible NDJSON spool before relay.

<!-- req: R2.4 -->
### R2.4

Spool rotation SHALL be bounded by values and SHALL occur only between complete batches.

<!-- req: R3 -->
## R3 - Modes and policy

<!-- req: R3.1 -->
### R3.1

Test mode SHALL receive and spool without requiring Jaeger; its success evidence SHALL be the deterministic spool plus receiver monitor state.

<!-- req: R3.2 -->
### R3.2

Relay mode SHALL spool first and then export the unchanged selected batch to the configured Jaeger OTLP endpoint.

<!-- req: R3.3 -->
### R3.3

Filtering and sampling SHALL use declared machine words over await metadata; transport handlers SHALL NOT hide those decisions.

<!-- req: R3.4 -->
### R3.4

Retry after relay failure SHALL be a declared transition and SHALL retain the spooled batch as evidence.

<!-- req: R4 -->
## R4 - Control and monitor

<!-- req: R4.1 -->
### R4.1

The profile SHALL run standard control and monitor servers with health, lifecycle exit, state, event, and stream surfaces.

<!-- req: R4.2 -->
### R4.2

Lifecycle exit SHALL stop the OTLP receiver, unblock await, release the listener, and then stop monitor and control servers.

<!-- req: R4.3 -->
### R4.3

Monitor state SHALL report launch, received, spooled, relayed, rejected, dropped, timed-out, and stopped outcomes without exposing raw span payloads.

<!-- req: R5 -->
## R5 - Self-loop exclusion

<!-- req: R5.1 -->
### R5.1

The collector process SHALL configure its own agent telemetry with a file exporter only and SHALL have no OTLP exporter endpoint.

<!-- req: R5.2 -->
### R5.2

Relay registration SHALL reject an endpoint equivalent to the receiver address after wildcard and loopback normalization.

<!-- req: R5.3 -->
### R5.3

Proof SHALL export a client-agent batch and observe it at the receiver while collector launch, await, spool, relay, monitor, and stop spans remain only in the collector trace file.

<!-- req: R5.4 -->
### R5.4

A received resource whose service.name equals collector SHALL NOT be dropped solely by name; self-loop prevention is configuration and process ownership, not service filtering.

<!-- req: R6 -->
## R6 - Staged Helm migration

<!-- req: R6.1 -->
### R6.1

Values SHALL select `contrib` or `agent` ingress; `agent` SHALL be the default.

<!-- req: R6.2 -->
### R6.2

In agent mode, mesh agents and Chroma SHALL send traces and OTLP metrics to the collector agent, which SHALL relay traces to Jaeger.

<!-- req: R6.3 -->
### R6.3

Agent mode SHALL deploy no contrib collector and no separate collector-metrics gateway; the pull-based in-cluster Dolt Prometheus scrape SHALL be retired in agent mode (GH-1366).

<!-- req: R6.4 -->
### R6.4

Agent-mode health SHALL require the collector-agent trace and OTLP metric path; the opt-in `contrib` option SHALL remain available for the trace path only.

<!-- req: R6.5 -->
### R6.5

Teardown SHALL NOT delete persistent host trace or metric evidence.

<!-- field: non_goals -->
## Non-goals

- Reimplementing OTLP transport in the profile.
- Receiving OTLP logs in the collector agent.
- Adding a Prometheus scrape capability to the collector agent (the retired Dolt scrape is not replaced).
- General-purpose collector throughput or multi-tenant isolation.
- Querying Jaeger from the collector profile.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

A persistent profile launches bounded OTLP trace intake and visibly composes await, policy, spool, relay, retry, and stop.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R1.4

<!-- ac: AC2 -->
### AC2

Spool and relay preserve complete span and integration identity, with bounded rotation between complete batches.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3
- R2.4

<!-- ac: AC3 -->
### AC3

One profile supports deterministic spool-only test mode and spool-first relay mode with declared filtering and retry.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3
- R3.4

<!-- ac: AC4 -->
### AC4

Control and monitor surfaces expose bounded lifecycle outcomes and lifecycle exit releases every listener and waiter.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
- R4.3

<!-- ac: AC5 -->
### AC5

File-only process telemetry and normalized endpoint validation prevent a self-loop, proven by client spans entering while collector spans do not.

<!-- field: traces -->
#### Traces

- R5.1
- R5.2
- R5.3
- R5.4

<!-- ac: AC6 -->
### AC6

Agent mode is the default and deploys no contrib collector or collector-metrics gateway; agents send traces and OTLP metrics to the collector agent, the in-cluster Dolt scrape is retired, contrib stays selectable for the trace path only, and teardown preserves persistent host evidence.

<!-- field: traces -->
#### Traces

- R6.1
- R6.2
- R6.3
- R6.4
- R6.5

<!-- field: platform_dependencies | form: "typed table" -->
## Platform dependencies

| srd_id | provider | symbols_used |
| --- | --- | --- |
| srd042-otlp-receiver-tools | agent-core | [R1, R2, R3, R4, R5, R6, R7, R8] |
| srd008-telemetry | agent-core | [R9] |

<!-- field: implementation_status | form: "tree" -->
## Implementation status

- implemented:
  1. The collector profile, bounded OTLP intake, deterministic spool and relay modes, self-loop exclusion, default agent-mode trace and OTLP metric intake with the opt-in contrib trace path, and persistent evidence checks. mage integration:rig drives a real client batch through the shipped filter-to-spool machine leg, while Helm smoke and observability evidence cover live routing and identity.
- planned:
  1. One lifecycle integration that exits while await_spans is blocked, proves the receiver, control, and monitor addresses all rebind, and records the bounded final monitor outcome required by R4 and AC4.
