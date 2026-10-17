---
id: rel08.0-uc002-corpus-query
title: Corpus Query from Chroma
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

The reader profile answers a question grounded in a Chroma collection. It
embeds the question at a local Ollama provider, queries Chroma with the
agent-supplied query vector for the nearest chunks, and invokes the model over
the retrieved chunks to produce a grounded answer. The proof belongs in
applications/catalog because the Chroma and Ollama REST definitions, the collection
identity, and the reader prompt are profile-owned assets. The deterministic
profile-owned conformance tracer executes the shipped wrapper against
controlled Chroma and Ollama protocol fixtures, proving the grounded-answer
path without treating an unavailable external service as passing evidence.

As in ingest, Chroma does not embed the query server-side over its raw HTTP
API, so the query vector is computed by the embed word and threaded into the
Chroma query word using the agent-core previous-Result parameter threading
(srd028 R12). The retrieved chunks, not a model-authored vector, ground the
answer.

<!-- field: actor -->
## Actor

Profile maintainer or operator asking a question against an ingested corpus.

<!-- field: trigger -->
## Trigger

A maintainer runs the profile-owned reader conformance cases:

  go test ./conformance -run TestCorpusReader -count=1

<!-- field: preconditions -->
## Preconditions

- A compatible agent-core checkout exists at the sibling ../../agent-core path.
- The deterministic harness serves known Chroma heartbeat, collection, and query responses.
- The deterministic harness serves Ollama tags, embedding, and chat responses for the configured answer model.
- agents/knowledge-manager/corpus-reader profile validates and owns the Chroma and Ollama REST definitions and the reader prompt.

<!-- field: flow -->
## Flow

- F1: The reader machine runs chroma_ready, failing fast when Chroma is unreachable, then composes the question and runs the embed-query stage (srd058 R4.1).
- F2: embed_query calls Ollama with the configured model for the request question, maps the returned vector into its Result output, and carries it forward.
- F3: The create-or-get collection word (immediately before query) resolves the existing collection id and carries the query vector forward so single-hop threading reaches the query body.
- F4: chroma_query threads the resolved collection id and the query vector into the query request body via body_source previous_result and input_mapping, placing the vector in a single-token array element, and requests the configured number of nearest chunks.
- F5: invoke_llm answers the question grounded in the retrieved chunks, recording the retrieved chunk identities.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd012-chroma-corpus-agents](../software-requirements/srd012-chroma-corpus-agents.md) [AC1](../software-requirements/srd012-chroma-corpus-agents.md#ac1) | The reader profile is a configured Chroma retrieval boundary machine, not a Go package. |
| T2 | [srd012-chroma-corpus-agents](../software-requirements/srd012-chroma-corpus-agents.md) [AC2](../software-requirements/srd012-chroma-corpus-agents.md#ac2) | chroma_ready, chroma_collection, embed_query, and chroma_query are the named REST vocabulary. |
| T3 | [srd012-chroma-corpus-agents](../software-requirements/srd012-chroma-corpus-agents.md) [AC3](../software-requirements/srd012-chroma-corpus-agents.md#ac3) | The embed-to-query vector step uses srd028 R12 body_source previous_result, input_mapping, and single-token reshape, with no model-authored vector. |
| T4 | [srd012-chroma-corpus-agents](../software-requirements/srd012-chroma-corpus-agents.md) [AC6](../software-requirements/srd012-chroma-corpus-agents.md#ac6) | The retrieved chunk identities and grounded answer are the checkable reader outputs. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The Chroma query carries an agent-supplied query vector threaded from the embed word, not a model-authored or server-side vector.

<!-- field: traces -->
#### Traces

- [srd012-chroma-corpus-agents](../software-requirements/srd012-chroma-corpus-agents.md) [R3.2](../software-requirements/srd012-chroma-corpus-agents.md#r32)
- [srd012-chroma-corpus-agents](../software-requirements/srd012-chroma-corpus-agents.md) [R3.4](../software-requirements/srd012-chroma-corpus-agents.md#r34)

<!-- sc: S2 -->
### S2

The reader retrieves the configured number of nearest chunks from the resolved collection before answering.

<!-- field: traces -->
#### Traces

- [srd012-chroma-corpus-agents](../software-requirements/srd012-chroma-corpus-agents.md) [R2.4](../software-requirements/srd012-chroma-corpus-agents.md#r24)

<!-- sc: S3 -->
### S3

invoke_llm produces an answer grounded in the retrieved chunks and records their identities.

<!-- field: traces -->
#### Traces

- [srd012-chroma-corpus-agents](../software-requirements/srd012-chroma-corpus-agents.md) [R1.3](../software-requirements/srd012-chroma-corpus-agents.md#r13)
- [srd012-chroma-corpus-agents](../software-requirements/srd012-chroma-corpus-agents.md) [R6.1](../software-requirements/srd012-chroma-corpus-agents.md#r61)

<!-- field: out_of_scope -->
## Out of scope

- Ingesting documents, which belongs to the ingest use case.
- Starting or stopping the Chroma or Ollama servers from the agent.
- Evaluating answer quality beyond the grounded-retrieval trace order.

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel08.0-corpus-agents](../test-suites/test-rel08.0-corpus-agents.md)
