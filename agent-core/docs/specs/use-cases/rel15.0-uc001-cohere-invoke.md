---
id: rel15.0-uc001-cohere-invoke
title: A profile selects Cohere for an invoke_llm classifier
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A profile changes an invoke_llm ToolDef from Ollama to Cohere by selecting
provider cohere, a Cohere model, provider_url, and response_profile cohere.
MachineSpec, prompt assembly, conversation, capture, and parse words remain
unchanged. The adapter sends the assembled messages to Cohere v2 Chat,
normalizes ordered text blocks and usage, and returns LLMResponded.

The bearer credential remains process-owned. Registration does not read it
or contact Cohere. Chat resolves COHERE_API_KEY immediately before one HTTP
request, and a missing credential or provider failure becomes CommandError.

<!-- field: actor -->
## Actor

A profile author selecting a Cohere-backed builtin classifier

<!-- field: trigger -->
## Trigger

An invoke_llm transition reaches a ToolDef configured with provider cohere
and a trusted Cohere provider URL.

<!-- field: flow -->
## Flow

- F1: The profile loads an invoke_llm ToolDef whose config names provider cohere, model command-r7b-12-2024, provider_url, manifest_state, and response_profile cohere.
- F2: Builder construction resolves the Cohere adapter and parser profile without network I/O or credential lookup.
- F3: invoke_llm assembles provider-neutral chronological messages from its prompt, conversation, and state-scoped manifest.
- F4: Chat resolves COHERE_API_KEY and creates one bearer-authenticated POST to provider_url/v2/chat with model, messages, temperature, and seed.
- F5: The Cohere adapter joins ordered text blocks, maps usage tokens, and returns one llm.ChatResponse.
- F6: invoke_llm appends the normalized assistant text, emits LLMResponded, and records provider, model, server, usage, and response-model trace attributes.
- F7: The existing capture and parse words consume joined model text through response_profile cohere without seeing Cohere HTTP JSON.
- F8: Missing credential, cancellation, transport, status, decoding, or content failure produces a safe CommandError and no retry.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd048-cohere-adapter R1.1, R1.2, R1.3, R1.4, R1.5 | adapter behind the unchanged llm.Client boundary |
| T2 | srd048-cohere-adapter R2.1, R2.2, R2.3, R2.4, R2.5 | Cohere v2 request and cancellation |
| T3 | srd048-cohere-adapter R3.1, R3.2, R3.3, R3.4, R3.5 | lazy process credential and secret boundary |
| T4 | srd048-cohere-adapter R4.1, R4.2, R4.3, R4.4, R4.5, R4.6 | content-block and usage normalization |
| T5 | srd048-cohere-adapter R5.1, R5.2, R5.3, R5.4, R5.5, R5.6 | declarative provider and response-profile selection |
| T6 | srd048-cohere-adapter R6.1, R6.2, R6.3, R6.4, R6.5, R6.6, R6.7 | errors and GenAI telemetry |
| T7 | srd017-llm-extensions R2.1, R2.2, R2.3, R2.4, R2.5, R2.6, R2.7, R2.8, R2.9 | supported provider and parser-profile selection |
| T8 | srd015-genai-semconv R3.1, R3.2, R3.3 | Cohere provider identity |
| T9 | srd015-genai-semconv R5.1, R5.2, R5.3, R5.4, R5.5, R5.6 | inference and usage attributes |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

A synthetic provider cohere ToolDef builds without contacting the provider or reading COHERE_API_KEY.

<!-- field: traces -->
#### Traces

- srd048-cohere-adapter AC1
- srd048-cohere-adapter AC4

<!-- sc: S2 -->
### S2

One authenticated v2 request carries chronological messages and supported decoding options but no Ollama-specific fields.

<!-- field: traces -->
#### Traces

- srd048-cohere-adapter AC1
- srd048-cohere-adapter AC2

<!-- sc: S3 -->
### S3

Ordered text blocks and usage tokens become the same ChatResponse and LLMResponded path used by Ollama.

<!-- field: traces -->
#### Traces

- srd048-cohere-adapter AC3

<!-- sc: S4 -->
### S4

The Cohere response profile extracts envelope-wrapped or bare JSON from joined model text without parsing provider response JSON.

<!-- field: traces -->
#### Traces

- srd048-cohere-adapter AC3
- srd048-cohere-adapter AC4

<!-- sc: S5 -->
### S5

Cohere transport and dispatch spans identify Cohere, model, server, usage, response model, and errors without exposing the credential.

<!-- field: traces -->
#### Traces

- srd048-cohere-adapter AC2
- srd048-cohere-adapter AC5

<!-- sc: S6 -->
### S6

Missing credentials and provider failures emit CommandError at dispatch, while unknown providers and missing Cohere URLs fail builder construction.

<!-- field: traces -->
#### Traces

- srd048-cohere-adapter AC4
- srd048-cohere-adapter AC5

<!-- field: failure_scenarios -->
## Failure scenarios

| id | scenario | expected |
| --- | --- | --- |
| E1 | COHERE_API_KEY is missing when Chat begins. | Chat fails before network I/O and invoke_llm emits CommandError without naming a secret value. |
| E2 | Cohere returns non-2xx, malformed JSON, or no text blocks. | The adapter returns a distinct bounded error, records error.type, and does not retry. |
| E3 | A ToolDef selects cohere without provider_url or names an unknown provider. | Builder construction rejects the configuration before dispatch. |

<!-- field: out_of_scope -->
## Out of scope

- Chatbot-mesh declaration, model-default, Helm, or REST-word changes.
- Cohere embeddings, reranking, native tool calls, streaming, or live-service CI.
- Provider failover and retry policy.

<!-- field: test_suite_trace -->
## Test suite trace

test-rel15.0-cohere-invoke
