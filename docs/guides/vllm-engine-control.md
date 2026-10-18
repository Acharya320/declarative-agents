<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

# vLLM engine control for inference agents

This guide records the options for integrating declarative agents with a
co-located [vLLM](https://docs.vllm.ai/) inference engine. It is an options
record, not a requirements document: nothing here is implemented, and no SRD
exists for any surface it names. If a surface graduates to real work, it gets
its own SRD citing this guide. vLLM flags and endpoints were surveyed on
2026-10-17 from the public vLLM documentation; verify each against the
deployed vLLM version before relying on it, because several are opt-in flags
whose names have moved between releases.

## The inference agent pattern

The deployment this guide serves: one vLLM engine per GPU node, owned by an
inference agent on that node. The agent — not the engine — is the unit other
agents see. It claims inference tasks from a message bus (see
[messaging-bus-analysis](messaging-bus-analysis.md)), decides admission
against the engine's capacity, executes the inference, and publishes the
result. vLLM is an implementation detail inside the agent's profile, the way
a shell command is an implementation detail of an exec word.

Two distinct uses of the engine fall out, and they take different tools:

1. **Inference as the workload.** The task payload carries a prepared
   request; the agent executes it with no conversation state. This is a plain
   HTTP call against vLLM's OpenAI-compatible server
   (`POST /v1/chat/completions`), and the existing REST client tools
   (srd028) cover it: retries, failure-to-signal mapping, auth, and response
   field extraction (`$.usage.*` for accounting) are already declared
   grammar. **Decision: no new tool for the completion call.** Wrapping one
   OpenAI-compatible API in another API gains nothing.
2. **The engine as an agent's own model.** A profile that reasons with the
   local model — history, prompt assembly, structured parse policies — must
   go through `invoke_llm`, which needs a vLLM provider library
   (srd058): one `chat-dialect.yaml` in the ollama/cohere pattern, with the
   reply at `$.choices[0].message.content` and usage at
   `$.usage.prompt_tokens` / `$.usage.completion_tokens`. Deferred until an
   agent actually needs the local model as its brain.

## Engine-control surfaces

The value of an inference agent is not proxying the completion endpoint but
exploiting the control surfaces a generic OpenAI client never touches. Each
option below names the surface, what it would let a profile declare, and its
stability. All are deferred.

### Priority scheduling (stable HTTP)

With `--scheduling-policy priority`, each request carries a priority and the
scheduler preempts lower-priority work rather than queueing behind it. An
agent would map bus task classes to engine priorities: an escalation
genuinely jumps ahead of batch work already in flight. This is the only
surface that changes *scheduling* behavior rather than observing it.

### Tokenize-based admission (stable HTTP)

`POST /tokenize` returns exact token counts without running inference. An
agent that tokenizes a claimed task's prompt before accepting it can compute
the KV-cache footprint arithmetically instead of estimating, turning the
admission gate from a heuristic into a ledger. vLLM has no reservation API;
an agent that is the engine's only client *is* the reservation mechanism,
and this endpoint makes its ledger exact.

### Capacity metrics (stable HTTP)

`GET /metrics` exposes Prometheus gauges — `vllm:num_requests_waiting`,
`vllm:gpu_cache_usage_perc` — that re-anchor the agent's admission ledger
after restarts and drift. Prometheus text format is not JSON, so the REST
tools cannot extract from it today; a metrics-scrape word (gauge names in,
values out) composed with the value predicate (srd041) would make the
capacity gate a declared transition condition. This is the smallest
plausible new tool.

### Dynamic LoRA adapters (stable HTTP, opt-in flag)

With runtime LoRA updating enabled, `POST /v1/load_lora_adapter` and its
unload twin mount task-specific adapters on the running base model. One base
model per node, but the agent hot-loads an adapter when it claims a task of
that type: specialized workflow stages without separate models or nodes. The
capability set in an agent's profile becomes a list of adapters it may
mount.

### Sleep and wake (stable HTTP, opt-in)

The sleep endpoints offload weights and free GPU memory without killing the
server; wake restores them. Enables node-level scale-to-zero without pod
churn, or an agent swapping which model its node serves based on backlog
shape. Wake latency is model-load latency, so this is a capacity-planning
lever, not a per-task one.

### Prefix-cache control and events (semi-public)

`POST /reset_prefix_cache` gives cache hygiene between tenants or
workflows. More interesting: with KV-event publishing configured, vLLM
emits which prefix blocks are cached over ZMQ, the feed that external
routers (llm-d, NVIDIA Dynamo) use for prefix-aware routing. Tasks from the
same workflow share long prompt prefixes, so routing siblings to the node
already holding their prefix is a large real-world throughput win. A
consumer for this feed is the most speculative option here and only pays at
multi-node scale.

### EngineCore ZMQ protocol (internal — ruled out)

vLLM V1 splits the OpenAI front end from the EngineCore process; they speak
msgpack over ZMQ. Speaking that protocol directly from Go would bypass the
HTTP layer entirely — and chain this repository to an unversioned internal
protocol that moves between vLLM releases. **Ruled out.** The stable
surfaces above cover every behavior the inference agent pattern needs.

## Summary of decisions

| Surface | Status |
|---|---|
| Completion calls via REST tools (srd028) | Decided — no new tool |
| vLLM provider dialect for `invoke_llm` (srd058) | Deferred until an agent needs the local model |
| Metrics-scrape word + value predicate gate | Candidate — smallest useful addition |
| Priority, tokenize, LoRA, sleep/wake | Options on record, no work planned |
| KV-event (ZMQ) consumer | Speculative, multi-node scale only |
| EngineCore protocol | Ruled out |

The repository is in maintenance mode; any item graduating from this guide
needs an explicit exception in the GH-2580 pattern before an SRD is written.

## Sources

- [vLLM OpenAI-compatible server](https://docs.vllm.ai/en/latest/serving/openai_compatible_server.html)
  — endpoints, including tokenize and LoRA management.
- [vLLM production metrics](https://docs.vllm.ai/en/latest/usage/metrics.html)
  — gauge names used above.
- [llm-d](https://llm-d.ai/) and
  [NVIDIA Dynamo](https://developer.nvidia.com/dynamo) — external consumers
  of the KV-event feed, cited as prior art for prefix-aware routing.
