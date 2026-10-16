---
id: srd009-ollama-adapter
title: Ollama LLM Adapter
depends_on:
  - srd_id: srd007-tracer-interface
    symbols_used: [Tracer, Event, Push, SetAttributes]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Every agent built on agent-core needs to talk to an LLM backend. The
first supported backend is Ollama, a local inference server. Both
generator and planner independently implemented Ollama HTTP adapters
that satisfy the llm.Client interface. The duplicated code drifted:
generator's adapter had tracing and model checks; planner's had only a
minimal Chat stub.

A single, well-tested adapter in agent-core eliminates the duplication
and gives callers tracing without widening the llm.Client port beyond Chat.
Adapter construction performs no network I/O: profiles that require
preflight readiness declare a REST transition, while other profiles route
backend failure from the first Chat dispatch as CommandError (GH-1375).

<!-- field: goals -->
## Goals

- G1: Provide a production-ready Ollama adapter in internal/model/llm/ollama satisfying llm.Client.
- G2: Keep model availability checks in profile-declared workflow words.
- G3: Support optional tracing via the internal/observability/tracing.Tracer interface.
- G4: Replace generator's internal/llm/adapter.go and planner's cmd/planner/ollama.go.

<!-- req: R1 -->
## R1 - Adapter struct

<!-- req: R1.1 -->
### R1.1

Adapter holds the Ollama base URL, model name, an http.Client, and an optional tracing.Tracer.

<!-- req: R1.2 -->
### R1.2

NewAdapter accepts the base URL, model name, and variadic Option functions without performing network I/O.

<!-- req: R1.3 -->
### R1.3

The default http.Client has a 5-minute timeout. WithHTTPClient overrides it.

<!-- req: R1.4 -->
### R1.4

Adapter.Model() returns the model name the adapter was created for.

<!-- req: R2 -->
## R2 - Chat (POST /api/chat)

<!-- req: R2.1 -->
### R2.1

Chat converts []llm.Message to the Ollama JSON schema and sends a non-streaming POST to /api/chat.

<!-- req: R2.2 -->
### R2.2

The request body includes model, messages, stream=false, and options (temperature, seed, num_ctx) from ChatOptions. NumCtx is omitted when zero.

<!-- req: R2.3 -->
### R2.3

The response JSON is parsed into llm.ChatResponse with Content from message.content, TokensIn from prompt_eval_count, TokensOut from eval_count.

<!-- req: R2.4 -->
### R2.4

Missing token counts default to zero.

<!-- req: R4 -->
## R4 - Declarative model availability

<!-- req: R4.1 -->
### R4.1

NewAdapter SHALL NOT call GET /api/tags or perform any other hidden readiness probe.

<!-- req: R4.2 -->
### R4.2

A profile that requires model readiness SHALL declare the probe as a REST word and route its signals in MachineSpec.

<!-- req: R4.3 -->
### R4.3

Without a declared preflight, an unavailable backend SHALL surface from Chat as invoke_llm CommandError, unless the ToolDef opts into the provider failure signals of srd058 R2.4, in which case it surfaces as ProviderUnavailable.

<!-- req: R4.4 -->
### R4.4

The adapter SHALL NOT auto-pull models.

<!-- req: R5 -->
## R5 - Tracing

<!-- req: R5.1 -->
### R5.1

WithTracer injects a tracing.Tracer. When set, the adapter creates a semconv inference span (via chatSpan) for each Chat() call with genai attributes (model, provider, server address, port, token counts, response model).

<!-- req: R5.2 -->
### R5.2

Chat spans record HTTP and parsing errors with relevant attributes.

<!-- req: R5.3 -->
### R5.3

A nil tracer is safe — all tracing is no-op by default.

<!-- req: R6 -->
## R6 - Error handling

<!-- req: R6.1 -->
### R6.1

Connection failures, non-200 HTTP status, and malformed JSON each produce distinct error messages.

<!-- req: R6.2 -->
### R6.2

Chat propagates context cancellation via http.NewRequestWithContext.

<!-- req: R6.3 -->
### R6.3

The adapter does not retry. Retry logic belongs to the caller, which is why the failure a caller sees must say whether retrying is worth anything (srd058 R2.4).

<!-- field: non_goals -->
## Non-goals

- Streaming responses.
- Adapters for non-Ollama providers (OpenAI, vLLM).
- Automatic model pulling.
- Requiring every profile to run a model-availability preflight before Chat.
- Prompt assembly or conversation history (those are separate packages).
- Cost/pricing calculation (zero for local Ollama).

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

NewAdapter returns a working llm.Client. Chat round-trips messages (including NumCtx) and extracts tokens. Adapter.Model() returns the configured model name.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.4
- R2.1
- R2.2
- R2.3
- R2.4

<!-- ac: AC2 -->
### AC2

NewAdapter performs no model-availability network request; profiles may declare REST readiness and otherwise Chat reports backend failure.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
- R4.3
- R4.4

<!-- ac: AC3 -->
### AC3

WithTracer enables Chat inference spans and error recording. Nil tracer is safe.

<!-- field: traces -->
#### Traces

- R5.1
- R5.2
- R5.3

<!-- ac: AC4 -->
### AC4

Connection failures, bad HTTP status, and malformed JSON produce distinct errors. Context cancellation propagates. No retry.

<!-- field: traces -->
#### Traces

- R6.1
- R6.2
- R6.3

<!-- ac: AC5 -->
### AC5

WithHTTPClient overrides the default 5-minute-timeout client; invoke_llm registration performs no backend request and an unreachable backend emits CommandError when dispatched.

<!-- field: traces -->
#### Traces

- R1.2
- R1.3
- R4.1
- R4.3

<!-- ac: AC6 -->
### AC6

Generator imports internal/model/llm/ollama and deletes internal/llm/adapter.go. Planner imports internal/model/llm/ollama and deletes cmd/planner/ollama.go.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2

<!-- field: implementation_status -->
## Implementation status

```yaml
current: |
    Superseded by srd058-provider-libraries R5.2. The compiled adapter in
    internal/model/llm/ollama was removed after a capture-fixture parity run
    (GH-2240). Ollama chat is now the tools/providers/ollama library, read by
    the invoke_llm dialect client; provider: ollama and an omitted provider
    resolve to its dialect. Construction without network I/O, tracing, and
    failure reporting are held by internal/tools/llm tests.
```

<!-- field: diagrams -->
## Diagrams

````yaml
requirements_trace: |
    ```mermaid
    flowchart LR
      Problem["Problem: Every agent built on agent-core needs to talk to an L..."] --> Goals["Goals: Ollama LLM Adapter"]
      Goals --> R1["R1: Adapter struct"]
      Goals --> R2["R2: Chat (POST /api/chat)"]
      Goals --> R4["R4: Declarative model availability"]
      Goals --> R5["R5: Tracing"]
      Goals --> R6["R6: Error handling"]
      AC1["AC1: NewAdapter returns a working llm.Client. Chat round-trips..."]
      R1 --> AC1
      R2 --> AC1
      AC2["AC2: NewAdapter performs no hidden model-availability request..."]
      R4 --> AC2
      AC3["AC3: WithTracer enables Chat inference spans and error record..."]
      R5 --> AC3
      AC4["AC4: Connection failures, bad HTTP status, and malformed JSON..."]
      R6 --> AC4
      AC5["AC5: WithHTTPClient overrides the default 5-minute-timeout cli..."]
      R1 --> AC5
      AC6["AC6: Generator imports internal/model/llm/ollama and deletes i..."]
      R1 --> AC6
    ```
````
