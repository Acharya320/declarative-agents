<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

# Messaging bus analysis for distributed agents

This guide records the analysis for a messaging tool family that would let
declarative agents coordinate over a message bus: inference workers claiming
tasks they have capacity for, and process-manager agents scattering ordered
tasks and gathering their completions. It is analysis only. Nothing is
implemented, no SRD exists, and the repository's maintenance mode means the
work needs an explicit exception in the GH-2580 pattern before any is
written. NATS claims were checked against the NATS documentation on
2026-10-17.

## The coordination pattern

The target architecture is a blackboard-style work queue with
admission-controlled competing consumers, coordinated by process-manager
agents:

- **Workers** (for example an inference agent owning a node's vLLM engine —
  see [vllm-engine-control](vllm-engine-control.md)) watch a subject for
  tasks matching their capability. A worker claims a task only when it has
  capacity, performs it, stores the result, and publishes a completion.
- **Process managers** create tasks, publish them, hold the state of which
  completions they await, and on each arrival decide: emit more tasks, or
  finish. Scatter-gather sequencing is centralized; placement is not — the
  manager never chooses which worker takes a task.
- **Payloads carry references, not artifacts.** Messages hold a workflow id,
  task id, and objectstore references (srd059); bodies stay small and
  replayable.

Pull beats push here because only the worker knows its engine's capacity; a
dispatcher choosing workers would need that knowledge remotely and staleness
makes it wrong. The queue is the load balancer.

## Required semantics

| Semantic | Why |
|---|---|
| Claim with lease | A claimed task is invisible to other workers for a bounded time, not forever |
| Ack / nack | Finish releases the task; failure or rollback redelivers it |
| Redelivery on lease expiry | A dead worker's tasks return to the pool |
| Queue groups | Many workers, one subject, each task delivered to exactly one claimant |
| Deduplication | At-least-once delivery demands idempotent task ids |
| Correlation subjects | Gather is "claim replies carrying my correlation id" |
| Blocking receive with timeout | The claim word must wait without spinning |

## Fit with existing contracts

The REST tool family is the template; every piece of the messaging design
has a shipped precedent:

| Precedent | Messaging analog |
|---|---|
| srd028 rest-client-tools | Client words: publish, claim, ack, nack |
| srd030 machine-request-rest-binding | A claimed message starts a machine run, payload as request |
| srd046 request-driven-signal-source | Bus messages mapped to signals for a running machine |
| srd024 blocking-input-tools | Blocking claim semantics |
| srd025 rollback-lifecycle + undo | **Claim's compensation is nack**: rollback redelivers the task |
| srd058 provider-libraries | Declaration grammar reuses the REST format's grammar |
| srd059 objectstore-tools | Pass-by-reference payload convention |

The compensation row matters most: at-least-once delivery needs no new
mechanism, because registering nack as the claim word's compensation makes
the existing rollback machinery the redelivery trigger.

### Open fork: how a worker receives tasks

Two shapes, both with precedent, deliberately left undecided here:

1. **Blocking claim word** (srd024 shape): the worker profile is a visible
   loop — gate on capacity, claim, work, ack, repeat. The admission decision
   sits between seeing a message and claiming it, declared as a transition.
2. **Message-driven machine binding** (srd030 shape): the runtime claims and
   spawns a machine run per task; ack/nack tie to the run's terminal state.
   Cleaner profiles, implicit concurrency, but the capacity gate moves into
   binding configuration.

The claim word keeps the original contract — *the agent sees capacity,
claims, works, acks* — literally readable in YAML; the binding is the
convenience layer that can follow once the words exist.

## Backend analysis: NATS JetStream

[NATS](https://docs.nats.io/) with JetStream persistence is the recommended
first backend. The decisive properties:

- **Pure Go, embeddable.** `nats-server` embeds in-process, so test suites
  run a real broker inside the test binary — no container, no network, in
  line with how this repository tests everything else. This alone separates
  NATS from every alternative below.
- **Semantics map one-to-one.** Queue groups are competing consumers;
  JetStream ack-wait is the lease; unacked messages redeliver;
  `Nats-Msg-Id` gives publish deduplication; pull consumers give blocking
  claim with timeout; subject hierarchies give per-task-type and
  per-correlation subjects without administration.
- **Operationally one binary** — a serious consideration for sovereign,
  self-managed deployments.

Alternatives considered:

| Backend | Verdict |
|---|---|
| Redis Streams | Consumer groups fit, but leases are manual (XAUTOCLAIM bookkeeping) and embedding for tests means miniredis, a reimplementation rather than the real engine |
| RabbitMQ | Mature semantics, but Erlang runtime: heavier operations, no in-process testing |
| Kafka | Partition-ordered log, wrong shape for a work queue with per-message acks; heavy operations |
| Google Pub/Sub | Managed equivalent with the same claim/ack/lease shape; the right second backend when deployments move to GCP, not a first backend for sovereign or local use |

The tool declarations should stay neutral enough (subjects, queue group,
lease, payload template, failure-to-signal map) that a Pub/Sub provider can
be added in the srd058 provider-library pattern without touching profiles.

## What implementation would entail

Recorded for scoping, not as a plan: one `internal/tools/messaging` package
and its `tools/units` declarations; an SRD in the software-requirements
corpus with a road-map release entry; a use case and a YAML test suite over
the embedded broker; and the maintenance-mode exception that authorizes all
of it. The client words are the whole first release — the srd030-style
binding and the srd046-style signal source are follow-ups, not
prerequisites.

## Sources

- [NATS JetStream documentation](https://docs.nats.io/nats-concepts/jetstream)
  — acks, redelivery, pull consumers, deduplication.
- [Embedding NATS in Go](https://docs.nats.io/running-a-nats-service/clients#embedding-nats)
  — in-process broker for tests.
- Enterprise Integration Patterns (Hohpe & Woolf) — Competing Consumers,
  Process Manager, Correlation Identifier: the vocabulary this guide reuses.
