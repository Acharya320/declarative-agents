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
| T1 | [srd048-cohere-adapter](../software-requirements/srd048-cohere-adapter.md) [R1.1](../software-requirements/srd048-cohere-adapter.md#r11), [R1.2](../software-requirements/srd048-cohere-adapter.md#r12), [R1.3](../software-requirements/srd048-cohere-adapter.md#r13), [R1.4](../software-requirements/srd048-cohere-adapter.md#r14), [R1.5](../software-requirements/srd048-cohere-adapter.md#r15) | adapter behind the unchanged llm.Client boundary |
| T2 | [srd048-cohere-adapter](../software-requirements/srd048-cohere-adapter.md) [R2.1](../software-requirements/srd048-cohere-adapter.md#r21), [R2.2](../software-requirements/srd048-cohere-adapter.md#r22), [R2.3](../software-requirements/srd048-cohere-adapter.md#r23), [R2.4](../software-requirements/srd048-cohere-adapter.md#r24), [R2.5](../software-requirements/srd048-cohere-adapter.md#r25) | Cohere v2 request and cancellation |
| T3 | [srd048-cohere-adapter](../software-requirements/srd048-cohere-adapter.md) [R3.1](../software-requirements/srd048-cohere-adapter.md#r31), [R3.2](../software-requirements/srd048-cohere-adapter.md#r32), [R3.3](../software-requirements/srd048-cohere-adapter.md#r33), [R3.4](../software-requirements/srd048-cohere-adapter.md#r34), [R3.5](../software-requirements/srd048-cohere-adapter.md#r35) | lazy process credential and secret boundary |
| T4 | [srd048-cohere-adapter](../software-requirements/srd048-cohere-adapter.md) [R4.1](../software-requirements/srd048-cohere-adapter.md#r41), [R4.2](../software-requirements/srd048-cohere-adapter.md#r42), [R4.3](../software-requirements/srd048-cohere-adapter.md#r43), [R4.4](../software-requirements/srd048-cohere-adapter.md#r44), [R4.5](../software-requirements/srd048-cohere-adapter.md#r45), [R4.6](../software-requirements/srd048-cohere-adapter.md#r46) | content-block and usage normalization |
| T5 | [srd048-cohere-adapter](../software-requirements/srd048-cohere-adapter.md) [R5.1](../software-requirements/srd048-cohere-adapter.md#r51), [R5.2](../software-requirements/srd048-cohere-adapter.md#r52), [R5.3](../software-requirements/srd048-cohere-adapter.md#r53), [R5.4](../software-requirements/srd048-cohere-adapter.md#r54), [R5.5](../software-requirements/srd048-cohere-adapter.md#r55), [R5.6](../software-requirements/srd048-cohere-adapter.md#r56) | declarative provider and response-profile selection |
| T6 | [srd048-cohere-adapter](../software-requirements/srd048-cohere-adapter.md) [R6.1](../software-requirements/srd048-cohere-adapter.md#r61), [R6.2](../software-requirements/srd048-cohere-adapter.md#r62), [R6.3](../software-requirements/srd048-cohere-adapter.md#r63), [R6.4](../software-requirements/srd048-cohere-adapter.md#r64), [R6.5](../software-requirements/srd048-cohere-adapter.md#r65), [R6.6](../software-requirements/srd048-cohere-adapter.md#r66), [R6.7](../software-requirements/srd048-cohere-adapter.md#r67) | errors and GenAI telemetry |
| T7 | [srd017-llm-extensions](../software-requirements/srd017-llm-extensions.md) [R2.1](../software-requirements/srd017-llm-extensions.md#r21), [R2.2](../software-requirements/srd017-llm-extensions.md#r22), [R2.3](../software-requirements/srd017-llm-extensions.md#r23), [R2.4](../software-requirements/srd017-llm-extensions.md#r24), [R2.5](../software-requirements/srd017-llm-extensions.md#r25), [R2.6](../software-requirements/srd017-llm-extensions.md#r26), [R2.7](../software-requirements/srd017-llm-extensions.md#r27), [R2.8](../software-requirements/srd017-llm-extensions.md#r28), [R2.9](../software-requirements/srd017-llm-extensions.md#r29) | supported provider and parser-profile selection |
| T8 | [srd015-genai-semconv](../software-requirements/srd015-genai-semconv.md) [R3.1](../software-requirements/srd015-genai-semconv.md#r31), [R3.2](../software-requirements/srd015-genai-semconv.md#r32), [R3.3](../software-requirements/srd015-genai-semconv.md#r33) | Cohere provider identity |
| T9 | [srd015-genai-semconv](../software-requirements/srd015-genai-semconv.md) [R5.1](../software-requirements/srd015-genai-semconv.md#r51), [R5.2](../software-requirements/srd015-genai-semconv.md#r52), [R5.3](../software-requirements/srd015-genai-semconv.md#r53), [R5.4](../software-requirements/srd015-genai-semconv.md#r54), [R5.5](../software-requirements/srd015-genai-semconv.md#r55), [R5.6](../software-requirements/srd015-genai-semconv.md#r56) | inference and usage attributes |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

A synthetic provider cohere ToolDef builds without contacting the provider or reading COHERE_API_KEY.

<!-- field: traces -->
#### Traces

- [srd048-cohere-adapter](../software-requirements/srd048-cohere-adapter.md) [AC1](../software-requirements/srd048-cohere-adapter.md#ac1)
- [srd048-cohere-adapter](../software-requirements/srd048-cohere-adapter.md) [AC4](../software-requirements/srd048-cohere-adapter.md#ac4)

<!-- sc: S2 -->
### S2

One authenticated v2 request carries chronological messages and supported decoding options but no Ollama-specific fields.

<!-- field: traces -->
#### Traces

- [srd048-cohere-adapter](../software-requirements/srd048-cohere-adapter.md) [AC1](../software-requirements/srd048-cohere-adapter.md#ac1)
- [srd048-cohere-adapter](../software-requirements/srd048-cohere-adapter.md) [AC2](../software-requirements/srd048-cohere-adapter.md#ac2)

<!-- sc: S3 -->
### S3

Ordered text blocks and usage tokens become the same ChatResponse and LLMResponded path used by Ollama.

<!-- field: traces -->
#### Traces

- [srd048-cohere-adapter](../software-requirements/srd048-cohere-adapter.md) [AC3](../software-requirements/srd048-cohere-adapter.md#ac3)

<!-- sc: S4 -->
### S4

The Cohere response profile extracts envelope-wrapped or bare JSON from joined model text without parsing provider response JSON.

<!-- field: traces -->
#### Traces

- [srd048-cohere-adapter](../software-requirements/srd048-cohere-adapter.md) [AC3](../software-requirements/srd048-cohere-adapter.md#ac3)
- [srd048-cohere-adapter](../software-requirements/srd048-cohere-adapter.md) [AC4](../software-requirements/srd048-cohere-adapter.md#ac4)

<!-- sc: S5 -->
### S5

Cohere transport and dispatch spans identify Cohere, model, server, usage, response model, and errors without exposing the credential.

<!-- field: traces -->
#### Traces

- [srd048-cohere-adapter](../software-requirements/srd048-cohere-adapter.md) [AC2](../software-requirements/srd048-cohere-adapter.md#ac2)
- [srd048-cohere-adapter](../software-requirements/srd048-cohere-adapter.md) [AC5](../software-requirements/srd048-cohere-adapter.md#ac5)

<!-- sc: S6 -->
### S6

Missing credentials and provider failures emit CommandError at dispatch, while unknown providers and missing Cohere URLs fail builder construction.

<!-- field: traces -->
#### Traces

- [srd048-cohere-adapter](../software-requirements/srd048-cohere-adapter.md) [AC4](../software-requirements/srd048-cohere-adapter.md#ac4)
- [srd048-cohere-adapter](../software-requirements/srd048-cohere-adapter.md) [AC5](../software-requirements/srd048-cohere-adapter.md#ac5)

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

[test-rel15.0-cohere-invoke](../test-suites/test-rel15.0-cohere-invoke.md)
