---
id: rel22.0-uc001-provider-libraries
title: A deployment changes and adds a model provider without editing an agent
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An operator runs an agent that chats through invoke_llm and embeds through
a stage fragment, with its profile declaring the providers root at the
shipped Ollama library. The operator reruns the same profile with
--library providers bound to the shipped Cohere library, and the agent
chats, embeds, and reports provider failures in the same signal vocabulary
without an edit to any machine, tool, REST, or fragment file. A library
author then adds a third provider as a directory of YAML: a chat dialect, an
embed fragment, and a parser profile. The operator binds it the same way,
and no Go change or agent declaration change follows.

<!-- field: actor -->
## Actor

An operator choosing a model provider per deployment, and a library author adding one

<!-- field: trigger -->
## Trigger

A deployment must run the same agents against a different model provider,
or against a provider the runtime has never shipped.

<!-- field: flow -->
## Flow

- F1: The profile declares libraries.providers at the shipped Ollama library; the agent's invoke_llm config names dialect /opt/providers/chat-dialect.yaml, and its machine expands the embed-query stage fragment with the library's embed_query REST fragment beside it.
- F2: invoke_llm renders the message history, builds the request from the dialect's body template with the conversation placed as a structured array, calls the provider, extracts text and usage by the dialect's selectors, parses the reply with the library's parser profile, and appends history.
- F3: The operator reruns with --library providers=<cohere library>. The token_ref the Cohere dialect names is resolved at call time, and the same machine emits LLMResponded and QueryEmbedded, or ProviderUnauthorized, ProviderThrottled, or ProviderUnavailable on a failing status.
- F4: The dump records the dialect and fragments by resolved path and root, and the program identity changes with the bound library.
- F5: A library author writes an OpenAI-shaped library directory with a chat dialect and its parser profile; the operator binds it, and the conformance suite passes with no Go change and no agent declaration change.
- F6: A binding that lacks a file the agent needs fails at load, naming the missing path and the bound directory.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [R1.1](../software-requirements/srd058-provider-libraries.md#r11), [R1.2](../software-requirements/srd058-provider-libraries.md#r12), [R1.3](../software-requirements/srd058-provider-libraries.md#r13), [R1.4](../software-requirements/srd058-provider-libraries.md#r14), [R1.5](../software-requirements/srd058-provider-libraries.md#r15) | the provider library, its root, its vocabulary, and its auth |
| T2 | [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [R2.1](../software-requirements/srd058-provider-libraries.md#r21), [R2.4](../software-requirements/srd058-provider-libraries.md#r24), [R2.5](../software-requirements/srd058-provider-libraries.md#r25), [R2.2](../software-requirements/srd058-provider-libraries.md#r22), [R2.3](../software-requirements/srd058-provider-libraries.md#r23) | the chat dialect declaration |
| T3 | [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [R3.1](../software-requirements/srd058-provider-libraries.md#r31), [R3.2](../software-requirements/srd058-provider-libraries.md#r32), [R3.3](../software-requirements/srd058-provider-libraries.md#r33), [R3.4](../software-requirements/srd058-provider-libraries.md#r34), [R3.5](../software-requirements/srd058-provider-libraries.md#r35) | invoke_llm as a template method |
| T4 | [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [R4.1](../software-requirements/srd058-provider-libraries.md#r41), [R4.2](../software-requirements/srd058-provider-libraries.md#r42) | embed and rerank stage fragments |
| T5 | [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [R5.1](../software-requirements/srd058-provider-libraries.md#r51), [R5.2](../software-requirements/srd058-provider-libraries.md#r52), [R5.3](../software-requirements/srd058-provider-libraries.md#r53) | first conforming libraries and the conformance suite |
| T6 | [srd056-unit-libraries](../software-requirements/srd056-unit-libraries.md) [R2.1](../software-requirements/srd056-unit-libraries.md#r21), [R2.2](../software-requirements/srd056-unit-libraries.md#r22), [R2.3](../software-requirements/srd056-unit-libraries.md#r23), [R2.4](../software-requirements/srd056-unit-libraries.md#r24) | declared roots and the --library override |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Rebinding --library providers switches a profile from Ollama to Cohere with no edit to any machine, tool, REST, or fragment file, and a missing dialect or fragment fails at load naming the path and the bound directory.

<!-- field: traces -->
#### Traces

- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [AC1](../software-requirements/srd058-provider-libraries.md#ac1)

<!-- sc: S2 -->
### S2

The dialect path produces the same request JSON documents and equivalent histories, signals, and spans to the compiled adapters for both shipped providers, and resolves credentials only at call time.

<!-- field: traces -->
#### Traces

- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [AC2](../software-requirements/srd058-provider-libraries.md#ac2)

<!-- sc: S3 -->
### S3

The embed stage fragments and library REST fragments replace the hand-written embedding stages and Ollama embed blocks, and every provider operation emits the shared vocabulary, including failures on the Ollama path.

<!-- field: traces -->
#### Traces

- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [AC3](../software-requirements/srd058-provider-libraries.md#ac3)

<!-- sc: S4 -->
### S4

A synthetic OpenAI-shaped library passes the conformance suite with zero Go and zero agent declaration changes.

<!-- field: traces -->
#### Traces

- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [AC4](../software-requirements/srd058-provider-libraries.md#ac4)

<!-- field: out_of_scope -->
## Out of scope

- Streaming, provider tool-call APIs, batching, caching, failover, retries, and rate-limit policy.
- Non-bearer request signing; such providers bind through an OpenAI-compatible gateway.
- A shared reasoning-loop machine template.
- Migrating consumer repositories.
