---
id: rel01.0-uc001-routed-fanout
title: Chatbot Turn Selects Sources and a Model Tier Before Trusted RAG Fan-Out
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A browser posts a chat turn to the chatbot agent (srd002). Beyond the core loop
(rel00.0-uc002), the turn adds separate classifiers for declared RAG sources and
one chat-LLM word through the $tool indirection. Structured parsing and subset
selection resolve model names to trusted topology entries before exactly one
sequential for_each queries that selected array with the one query embedding. The ordered join
feeds generic query-signal and embedding-model partitions; render_each formats
only compatible chunks under topology-name headers. QueryRejected, CommandError,
declared-but-not-selected, QueryRejected, CommandError, and reported-model
mismatch outcomes remain distinct in response metadata.

<!-- field: actor -->
## Actor

A browser user sending a chat message across two corpora.

<!-- field: trigger -->
## Trigger

The browser posts a chat turn to the chatbot agent's chat endpoint with two RAG
servers over disjoint corpora declared.

<!-- field: preconditions -->
## Preconditions

- Two RAG server agents (rel00.0-uc001) are reachable as declared rest.yaml clients over disjoint corpora.
- A small classifier model and two chat-LLM models are reachable as declared clients.
- Each topology item has a required description and trusted target/config fields.
- The chatbot profile validates and owns the turn machine, separate source and tier prompts, and fan-out declarations.

<!-- field: flow -->
## Flow

- F1: select_sources classifies the original question against the names-and-descriptions-only catalog; capture_source_selection JSON-escapes the adjacent raw reply into a fixed response field, parse_structured validates `{names:[string]}`, and select_subset resolves declared names to trusted topology entries. Classifier/capture error, malformed, empty, and only-undeclared outcomes select all declared entries; mixed picks discard undeclared names.
- F2: The separate tier selector preserves its existing classifier, $tool vocabulary, and default-word fallback contract.
- F3: embed_query embeds the query once; exactly one sequential MachineSpec for_each fans the vector over only the trusted selected array and emits one structured join.
- F4: partition_query_results separates QueryResponded entries from QueryRejected vector rejections and CommandError degradation entries; collect_all continues through every selected source.
- F5: partition_embedding_models excludes answered sources whose reported model differs from the query model, then render_grounding formats only the matched entries under topology-name headers.
- F6: The $tool-dispatched invoke_llm word answers over the composed prompt, citing chunks from both corpora when both RAGs contribute.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd002-chatbot-agent](../software-requirements/srd002-chatbot-agent.md) [AC2](../software-requirements/srd002-chatbot-agent.md#ac2) | Separate source and tier classifiers retain atomic prompts, trusted subset selection, and normative fallbacks. |
| T2 | [srd002-chatbot-agent](../software-requirements/srd002-chatbot-agent.md) [AC3](../software-requirements/srd002-chatbot-agent.md#ac3) | One sequential iterator plus generic partitions degrades query failures and excludes embedding-model mismatches without indexed words or states. |
| T3 | test-rel01.0-routing-fanout routed_fanout | The suite records tier selection, the two-RAG fan-out over disjoint corpora, and the degraded and excluded paths. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The source classifier selects a scoped declared subset, spanning questions may select all, and every classifier/malformed/empty/only-undeclared fallback queries all declared sources.

<!-- field: traces -->
#### Traces

- [srd002-chatbot-agent](../software-requirements/srd002-chatbot-agent.md) [R2.4](../software-requirements/srd002-chatbot-agent.md#r24)
- [srd002-chatbot-agent](../software-requirements/srd002-chatbot-agent.md) [R2.6](../software-requirements/srd002-chatbot-agent.md#r26)
- [srd002-chatbot-agent](../software-requirements/srd002-chatbot-agent.md) [AC2](../software-requirements/srd002-chatbot-agent.md#ac2)

<!-- sc: S2 -->
### S2

One query embedding fans out sequentially over trusted selected entries through one for_each; undeclared names never reach query target and not-selected names are distinct metadata.

<!-- field: traces -->
#### Traces

- [srd002-chatbot-agent](../software-requirements/srd002-chatbot-agent.md) [R3.1](../software-requirements/srd002-chatbot-agent.md#r31)
- [srd002-chatbot-agent](../software-requirements/srd002-chatbot-agent.md) [R3.4](../software-requirements/srd002-chatbot-agent.md#r34)
- [srd002-chatbot-agent](../software-requirements/srd002-chatbot-agent.md) [R6.2](../software-requirements/srd002-chatbot-agent.md#r62)
- [srd002-chatbot-agent](../software-requirements/srd002-chatbot-agent.md) [AC3](../software-requirements/srd002-chatbot-agent.md#ac3)

<!-- sc: S3 -->
### S3

A per-RAG failure is retained in the join's query_failed metadata set and the turn returns a mapped 200 composed from remaining compatible chunks.

<!-- field: traces -->
#### Traces

- [srd002-chatbot-agent](../software-requirements/srd002-chatbot-agent.md) [R3.2](../software-requirements/srd002-chatbot-agent.md#r32)
- [srd002-chatbot-agent](../software-requirements/srd002-chatbot-agent.md) [AC3](../software-requirements/srd002-chatbot-agent.md#ac3)

<!-- sc: S4 -->
### S4

A 200 response whose reported embedding model differs is excluded from render_each grounding and reported separately from a 400 QueryRejected vector rejection.

<!-- field: traces -->
#### Traces

- [srd002-chatbot-agent](../software-requirements/srd002-chatbot-agent.md) [R3.3](../software-requirements/srd002-chatbot-agent.md#r33)
- [srd002-chatbot-agent](../software-requirements/srd002-chatbot-agent.md) [AC3](../software-requirements/srd002-chatbot-agent.md#ac3)

<!-- field: out_of_scope -->
## Out of scope

- Parallel fan-out; v1 fan-out is sequential.
- The single-RAG core loop, which is rel00.0-uc002.
- Panel content and SPA implementation, owned by their own issues.

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel01.0-routing-fanout](../test-suites/test-rel01.0-routing-fanout.md)
