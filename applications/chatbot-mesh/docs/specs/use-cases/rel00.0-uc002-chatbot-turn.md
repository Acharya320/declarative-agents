---
id: rel00.0-uc002-chatbot-turn
title: Chatbot Turn Embeds, Retrieves From One RAG, and Answers
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A browser posts a chat turn to the chatbot agent (srd002). The turn runs a
request-scoped machine: embed_query embeds the query once, that vector queries a
single RAG server (srd001, vector-in), compose builds a answer prompt from the
returned chunks and the original query, and an invoke_llm word answers; the
terminal signal maps to the HTTP chat response. Non-adjacent data reaches compose
and answer through command-state $from() rather than carry_forward chains. When
the RAG fails, the turn degrades and answers from the message alone as a mapped
200 rather than a 500. This is the core loop (epic #293); the $tool tier selector and
the multi-RAG fan-out are release 01.0.

<!-- field: actor -->
## Actor

A browser user sending a chat message.

<!-- field: trigger -->
## Trigger

The browser posts a chat turn (message plus conversation history) to the
chatbot agent's chat endpoint with one RAG server declared.

<!-- field: preconditions -->
## Preconditions

- A RAG server agent (rel00.0-uc001) is reachable as a declared rest.yaml client.
- A chat LLM and an embedding provider are reachable as declared clients.
- The chatbot profile validates and owns the turn machine and the ui.yaml panel placement.

<!-- field: flow -->
## Flow

- F1: The chat endpoint binds the turn to a request-scoped machine seeded with the request body.
- F2: embed_query embeds the incoming query once; the single vector is the retrieval input.
- F3: The turn queries the declared RAG client with the query vector and selects the returned chunks.
- F4: compose reads the original query and the retrieved chunks through command-state $from() and builds the answer prompt.
- F5: The invoke_llm word answers over the composed prompt and the terminal signal maps to the HTTP chat response; a RAG failure degrades the turn to a mapped 200 answered from the message alone.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd002-chatbot-agent AC1 | The turn is a request-scoped machine with $from() data flow and one embedding computed once. |
| T2 | srd002-chatbot-agent AC4 | The turn is stateless; the browser sends conversation history with each turn. |
| T3 | test-rel00.0-chatbot-mesh chatbot_turn | The suite records the turn inputs and the expected embed -> retrieve -> compose -> answer trace, including the degraded-RAG case. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

One query embedding is computed once, and the original query and the retrieved chunks reach compose and answer through command-state $from() rather than carry_forward chains.

<!-- field: traces -->
#### Traces

- srd002-chatbot-agent R1.1
- srd002-chatbot-agent R1.2
- srd002-chatbot-agent R1.3
- srd002-chatbot-agent AC1

<!-- sc: S2 -->
### S2

A RAG failure degrades the turn so it answers from the message alone as a mapped 200 rather than a 500.

<!-- field: traces -->
#### Traces

- srd002-chatbot-agent R3.2
- srd002-chatbot-agent AC3

<!-- sc: S3 -->
### S3

The turn is stateless; conversation history arrives from the browser with each turn.

<!-- field: traces -->
#### Traces

- srd002-chatbot-agent R4.1
- srd002-chatbot-agent AC4

<!-- field: out_of_scope -->
## Out of scope

- The $tool tier selector and multi-RAG fan-out, which are release 01.0.
- Server-side session storage; the turn is stateless in v1.
- Panel content and SPA implementation, owned by their own issues.

<!-- field: test_suite_trace -->
## Test suite trace

test-rel00.0-chatbot-mesh
