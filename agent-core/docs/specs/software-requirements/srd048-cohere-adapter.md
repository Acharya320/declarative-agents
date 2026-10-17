---
id: srd048-cohere-adapter
title: Cohere v2 LLM Adapter
depends_on:
  - srd_id: srd009-ollama-adapter
    symbols_used: [llm.Client compatibility, no-constructor-I/O precedent]
  - srd_id: srd015-genai-semconv
    symbols_used: [ProviderCohere, InferenceAttrs, UsageAttrs, ErrorAttrs]
  - srd_id: srd017-llm-extensions
    symbols_used: [ResponseParser, ProfileRegistry, response_profile]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

invoke_llm supports only Ollama, so a profile that otherwise uses Cohere
still sends builtin classifier calls to a local Ollama model. The generic
llm.Client port already separates inference transport from MachineSpec,
prompt assembly, conversation history, capture, and response parsing. A
Cohere adapter can implement that port without changing those owners.

Cohere v2 Chat differs from Ollama on endpoint, authentication, request
options, response content, and usage fields. The adapter must normalize those
differences without putting a bearer credential in ToolDef configuration or
treating provider HTTP JSON as model output.

<!-- field: goals -->
## Goals

- G1: Implement Cohere v2 Chat behind the existing llm.Client port.
- G2: Make provider selection declarative through invoke_llm configuration.
- G3: Keep credentials outside profile YAML and resolve them only for a Chat call.
- G4: Emit the same GenAI identity and usage attributes as other invoke_llm providers.

<!-- req: R1 -->
## R1 - Adapter boundary

<!-- req: R1.1 -->
### R1.1

The Cohere adapter SHALL implement llm.Client without changing Client, Message, ChatOptions, ChatResponse, Conversation, or PromptAssembler.

<!-- req: R1.2 -->
### R1.2

NewAdapter SHALL accept a base URL, model, and testable options without performing network I/O or credential lookup.

<!-- req: R1.3 -->
### R1.3

The adapter SHALL use an injected HTTP client and tracing port when supplied and safe defaults otherwise.

<!-- req: R1.4 -->
### R1.4

Adapter.Model SHALL return the configured model.

<!-- req: R1.5 -->
### R1.5

The adapter SHALL NOT implement retries, streaming, readiness probes, model discovery, embeddings, reranking, or REST-declared provider words.

<!-- req: R2 -->
## R2 - Cohere v2 request

<!-- req: R2.1 -->
### R2.1

Chat SHALL send one non-streaming POST to the normalized provider URL plus /v2/chat.

<!-- req: R2.2 -->
### R2.2

The request SHALL contain model and chronological messages whose role and content come from the provider-neutral llm.Message values.

<!-- req: R2.3 -->
### R2.3

Temperature and seed SHALL be top-level request fields derived from ChatOptions.

<!-- req: R2.4 -->
### R2.4

The request SHALL NOT send Ollama options, num_ctx, /api/chat fields, or /api/tags probes.

<!-- req: R2.5 -->
### R2.5

Context cancellation and deadlines SHALL propagate through the HTTP request.

<!-- req: R3 -->
## R3 - Credential boundary

<!-- req: R3.1 -->
### R3.1

Chat SHALL resolve the fixed process credential COHERE_API_KEY through an injectable lookup immediately before request construction.

<!-- req: R3.2 -->
### R3.2

A missing or empty credential SHALL fail Chat before network I/O.

<!-- req: R3.3 -->
### R3.3

A present credential SHALL be sent only as Authorization Bearer on the Cohere request.

<!-- req: R3.4 -->
### R3.4

ToolDef config, command output, traces, logs, and errors SHALL NOT contain the credential or a credential-valued config field.

<!-- req: R3.5 -->
### R3.5

The adapter's process-environment lookup is the narrow srd048 exception to the go-style environment ban; tests SHALL replace it through an option.

<!-- req: R4 -->
## R4 - Response normalization

<!-- req: R4.1 -->
### R4.1

Chat SHALL decode the Cohere v2 assistant message and concatenate text from message.content blocks in authored order.

<!-- req: R4.2 -->
### R4.2

Non-text content blocks SHALL NOT become model text; a response with no usable text SHALL return a distinct error.

<!-- req: R4.3 -->
### R4.3

Chat SHALL map usage.tokens.input_tokens and usage.tokens.output_tokens to ChatResponse TokensIn and TokensOut.

<!-- req: R4.4 -->
### R4.4

Missing token counts SHALL map to zero.

<!-- req: R4.5 -->
### R4.5

The joined text SHALL become ChatResponse.Content and the existing invoke_llm command SHALL append it to conversation history and emit LLMResponded.

<!-- req: R4.6 -->
### R4.6

The embedded cohere response_profile SHALL parse joined model text, not Cohere HTTP response JSON, using the [tool_call] envelope with bare-JSON fallback.

<!-- req: R5 -->
## R5 - Provider selection and defaults

<!-- req: R5.1 -->
### R5.1

invoke_llm SHALL accept provider values ollama and cohere and SHALL reject unknown providers during builder construction.

<!-- req: R5.2 -->
### R5.2

An omitted provider SHALL retain the Ollama compatibility default and its localhost provider URL.

<!-- req: R5.3 -->
### R5.3

provider cohere SHALL require provider_url and SHALL NOT inherit the Ollama localhost default.

<!-- req: R5.4 -->
### R5.4

A Cohere builder SHALL be constructible without contacting the provider or resolving COHERE_API_KEY.

<!-- req: R5.5 -->
### R5.5

Missing credentials, unreachable service, and provider response failures SHALL become invoke_llm CommandError at dispatch.

<!-- req: R5.6 -->
### R5.6

response_profile cohere SHALL resolve by name, and command- model prefixes SHALL select it when response_profile is omitted.

<!-- req: R6 -->
## R6 - Errors and observability

<!-- req: R6.1 -->
### R6.1

Missing credential, request encoding, request creation, connection, non-2xx status, body read, malformed JSON, and missing text SHALL produce distinct safe errors.

<!-- req: R6.2 -->
### R6.2

A non-2xx error MAY include a bounded response diagnostic but SHALL NOT include the Authorization value.

<!-- req: R6.3 -->
### R6.3

Each Cohere transport span SHALL carry gen_ai.operation.name chat, gen_ai.provider.name cohere, gen_ai.request.model, and server.address with server.port when present.

<!-- req: R6.4 -->
### R6.4

A successful span SHALL record gen_ai.response.model and gen_ai.usage.input_tokens and gen_ai.usage.output_tokens.

<!-- req: R6.5 -->
### R6.5

HTTP and decoding failures SHALL record error.type and the error on the transport span.

<!-- req: R6.6 -->
### R6.6

A caller-context span SHALL parent the Cohere transport span; an adapter tracer SHALL remain a safe fallback.

<!-- req: R6.7 -->
### R6.7

The invoke_llm dispatch span and monitor metrics SHALL continue to derive provider and model from the selected ToolDef config.

<!-- field: non_goals -->
## Non-goals

- Change chatbot-mesh classifier, answer-model, Helm, or REST configuration.
- Change REST-declared Cohere command, embedding, or rerank words.
- Add Cohere streaming, tool-call APIs, response_format, safety modes, citations, or billing.
- Refactor or widen the Ollama-specific srd009 contract.
- Add provider failover, retries, rate-limit policy, or live Cohere CI credentials.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

A Cohere adapter implements llm.Client, performs no constructor I/O, sends exactly one v2 Chat request with provider-neutral messages and supported decoding options, and propagates cancellation.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R1.4
- R1.5
- R2.1
- R2.2
- R2.3
- R2.4
- R2.5

<!-- ac: AC2 -->
### AC2

COHERE_API_KEY is resolved only at Chat time, sent only as bearer authentication, replaceable in tests, absent from YAML and diagnostics, and required before network I/O.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3
- R3.4
- R3.5

<!-- ac: AC3 -->
### AC3

Ordered Cohere text blocks and usage tokens normalize into ChatResponse, while the cohere response profile parses only the joined raw model text.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
- R4.3
- R4.4
- R4.5
- R4.6

<!-- ac: AC4 -->
### AC4

invoke_llm selects Ollama or Cohere from ToolDef config, preserves omitted-provider compatibility, requires an explicit Cohere URL, resolves the Cohere parser profile, rejects unknown providers, and routes Chat-time failures as CommandError.

<!-- field: traces -->
#### Traces

- R5.1
- R5.2
- R5.3
- R5.4
- R5.5
- R5.6

<!-- ac: AC5 -->
### AC5

Cohere spans expose provider, request and response model, server authority, usage, parenting, and safe typed failures equivalent to the existing inference boundary.

<!-- field: traces -->
#### Traces

- R6.1
- R6.2
- R6.3
- R6.4
- R6.5
- R6.6
- R6.7

<!-- field: implementation_status | form: "tree" -->
## Implementation status

- current:
  Superseded by srd058-provider-libraries R5.2. The compiled adapter in
  internal/model/llm/cohere was removed after a capture-fixture parity run
  (GH-2240). Cohere v2 chat is now the tools/providers/cohere library: its
  chat-dialect.yaml declares the request, the text-block join, usage, the
  failure map, and the COHERE_API_KEY bearer reference that R3 requires, and
  invoke_llm reads it through the dialect client. provider: cohere resolves
  to that dialect. R1.5 holds: embed and rerank stay REST-declared, now as the
  library's fragments. The wire, credential, normalization, failure, and
  tracing contracts are held by internal/tools/llm tests against the
  captured adapter requests in testdata/parity.

<!-- field: related_documents -->
## Related documents

- cfg-llm-config-format
- cfg-runtime-contract
- rel15.0-uc001-cohere-invoke
- test-rel15.0-cohere-invoke

<!-- field: diagrams | form: "tree" -->
## Diagrams

- requirements_trace:
  ```mermaid
  flowchart LR
    P[Ollama-only invoke_llm] --> G[Cohere adapter goals]
    G --> R1[Adapter boundary]
    G --> R2[v2 request]
    G --> R3[Credential boundary]
    G --> R4[Response normalization]
    G --> R5[Provider selection]
    G --> R6[Errors and observability]
    R1 --> AC1
    R2 --> AC1
    R3 --> AC2
    R4 --> AC3
    R5 --> AC4
    R6 --> AC5
  ```
