---
id: rel00.0-uc001-rag-server-query
title: Persistent RAG Server Answers a Vector-In Query
status: specified
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A persistent RAG server agent (srd001) serves one Chroma collection behind a
machine_request query endpoint. A caller posts a query_embeddings vector and a
result count; the agent runs a request-scoped retrieval machine
(create-or-get collection, chroma_query) and returns the nearest chunks with
their record ids, distances, and metadata, including the collection's embedding
model name. The RAG server does not embed the query (the caller supplies the
vector) and does not answer (it returns chunks).

<!-- field: actor -->
## Actor

A caller agent or operator issuing a vector-in retrieval request.

<!-- field: trigger -->
## Trigger

A caller posts query_embeddings and n_results to the RAG server's query
endpoint while a Chroma server holds an ingested collection.

<!-- field: preconditions -->
## Preconditions

- A Chroma server is running with a collection ingested by the corpus-ingest seed agent.
- The RAG server profile validates and owns the Chroma REST definition, the served collection identity, and the query endpoint mapping.
- The query, monitor, and control server bind addresses are loopback-configurable inputs.

<!-- field: flow -->
## Flow

- F1: The caller posts query_embeddings and n_results to the machine_request query endpoint, which binds the request to one request-scoped run of the retrieval machine.
- F2: chroma_ready confirms the served Chroma database is reachable; the create-or-get collection word resolves the served collection id.
- F3: chroma_query threads the query_embeddings vector into the query body via the single-token reshape and requests n_results nearest chunks.
- F4: The endpoint maps the terminal signal to an HTTP response returning the chunks with record ids, distances, and metadata that names the collection's embedding model.
- F5: The monitor server exposes the run through the srd033 views and SSE, and the control server exposes exit and health.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd001-rag-server-agent](../software-requirements/srd001-rag-server-agent.md) [AC1](../software-requirements/srd001-rag-server-agent.md#ac1) | The RAG server hosts the reader retrieval chain behind a machine_request query endpoint with monitor and control servers. |
| T2 | [srd001-rag-server-agent](../software-requirements/srd001-rag-server-agent.md) [AC2](../software-requirements/srd001-rag-server-agent.md#ac2) | The vector-in contract accepts query_embeddings and returns chunks with ids, distances, and the embedding-model metadata. |
| T3 | [srd001-rag-server-agent](../software-requirements/srd001-rag-server-agent.md) [AC3](../software-requirements/srd001-rag-server-agent.md#ac3) | A single agent serves one collection; the run does not multiplex corpora. |
| T4 | test-rel00.0-chatbot-mesh rag_server_query | The suite records the query inputs and the expected chunk-and-metadata response. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The query endpoint accepts a caller-supplied query_embeddings vector and n_results and does not require question text.

<!-- field: traces -->
#### Traces

- [srd001-rag-server-agent](../software-requirements/srd001-rag-server-agent.md) [R2.1](../software-requirements/srd001-rag-server-agent.md#r21)
- [srd001-rag-server-agent](../software-requirements/srd001-rag-server-agent.md) [AC2](../software-requirements/srd001-rag-server-agent.md#ac2)

<!-- sc: S2 -->
### S2

The response returns chunks with record ids, distances, and metadata that names the served collection's embedding model.

<!-- field: traces -->
#### Traces

- [srd001-rag-server-agent](../software-requirements/srd001-rag-server-agent.md) [R2.2](../software-requirements/srd001-rag-server-agent.md#r22)
- [srd001-rag-server-agent](../software-requirements/srd001-rag-server-agent.md) [R2.3](../software-requirements/srd001-rag-server-agent.md#r23)
- [srd001-rag-server-agent](../software-requirements/srd001-rag-server-agent.md) [AC2](../software-requirements/srd001-rag-server-agent.md#ac2)

<!-- sc: S3 -->
### S3

The retrieval runs as a request-scoped machine over one served collection, with no LLM answering step in the RAG server.

<!-- field: traces -->
#### Traces

- [srd001-rag-server-agent](../software-requirements/srd001-rag-server-agent.md) [R1.1](../software-requirements/srd001-rag-server-agent.md#r11)
- [srd001-rag-server-agent](../software-requirements/srd001-rag-server-agent.md) [R1.5](../software-requirements/srd001-rag-server-agent.md#r15)
- [srd001-rag-server-agent](../software-requirements/srd001-rag-server-agent.md) [AC1](../software-requirements/srd001-rag-server-agent.md#ac1)

<!-- field: out_of_scope -->
## Out of scope

- Embedding the query, which the caller does.
- Composing or answering over the chunks, which the chatbot agent does.
- Ingest, which remains the corpus-ingest loading path.

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel00.0-chatbot-mesh](../test-suites/test-rel00.0-chatbot-mesh.md)
