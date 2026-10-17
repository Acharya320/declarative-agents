---
id: test-rel08.0-corpus-agents
title: Test Suite for Release 08.0 Chroma Corpus Agents
release: "08.0"
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

This suite proves the canonical corpus-ingest trust ordering and configurable
document-embedding provider binding. Deterministic ingest fixtures prove that
flat and single-row provider responses both reach Chroma as one vector and
that invalid rows fail before a write. The suite then executes the shipped corpus-reader
wrapper through its complete retrieval and answer path. The reader cases use
deterministic in-process Chroma and Ollama protocol fixtures. They patch only
transport addresses; the shipped profile, machine, declarations, REST
operations, prompt, and tool selection remain in control.

<!-- field: traces -->
## Traces

- rel08.0-uc001-corpus-ingest
- rel08.0-uc002-corpus-query

<!-- field: preconditions -->
## Preconditions

- A compatible agent-core checkout is available as the sibling ../../agent-core directory.
- No live Chroma or Ollama service is required for the deterministic reader conformance cases.

<!-- field: test_cases -->
## Test cases

<!-- case: Trusted corpus discovery precedes model control | status: implemented -->
### Trusted corpus discovery precedes model control

<!-- field: use_case -->
#### Use case

rel08.0-uc001-corpus-ingest

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCorpusIngestListsTrustedCorpusBeforeModelControl

<!-- field: expected | form: "tree" -->
#### Expected

- listing_corpus_precedes_model_control: true
- model_cannot_author_discovery_paths: true

<!-- field: traces -->
#### Traces

- srd012-chroma-corpus-agents R1.2
- srd012-chroma-corpus-agents R4.4

<!-- case: Consuming wrapper selects its document-embedding provider | status: implemented -->
### Consuming wrapper selects its document-embedding provider

<!-- field: use_case -->
#### Use case

rel08.0-uc001-corpus-ingest

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCorpusIngestSelectsConfiguredEmbeddingProvider

<!-- field: expected | form: "tree" -->
#### Expected

- ollama_embedding_operation_absent: true
- configured_provider_operation_selected: true
- provider_operation_names_globally_unique: true
- query_and_document_input_types_independent: true
- model_family_shared: true

<!-- field: traces -->
#### Traces

- rel08.0-uc001-corpus-ingest S4
- srd012-chroma-corpus-agents R2.6
- srd012-chroma-corpus-agents R4.5

<!-- case: Ingest normalizes flat and single-row provider vectors | status: implemented -->
### Ingest normalizes flat and single-row provider vectors

<!-- field: use_case -->
#### Use case

rel08.0-uc001-corpus-ingest

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCorpusIngestNormalizesProviderEmbeddings

<!-- field: expected | form: "tree" -->
#### Expected

- ollama_flat_vector_preserved: true
- cohere_single_row_unwrapped: true
- chroma_add_wraps_exactly_once: true
- document_and_id_preserved: true
- invalid_row_fails_before_chroma_add: true
- chat_inference_uses_protocol_fixture: true
- live_ollama_required: false

<!-- field: traces -->
#### Traces

- rel08.0-uc001-corpus-ingest S1
- rel08.0-uc001-corpus-ingest S4
- srd012-chroma-corpus-agents R3.5
- srd012-chroma-corpus-agents R6.2

<!-- case: Shipped reader embeds, retrieves, and answers from known chunks | status: implemented -->
### Shipped reader embeds, retrieves, and answers from known chunks

<!-- field: use_case -->
#### Use case

rel08.0-uc002-corpus-query

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCorpusReaderConformance

<!-- field: expected | form: "tree" -->
#### Expected

- shipped_wrapper_executes: true
- provider_vector_reaches_chroma_query: true
- retrieved_identity_reaches_model: true
- grounded_answer_reaches_succeeded: true

<!-- field: traces -->
#### Traces

- rel08.0-uc002-corpus-query S1
- rel08.0-uc002-corpus-query S2
- rel08.0-uc002-corpus-query S3

<!-- case: Reader REST failure reaches Failed before embedding | status: implemented -->
### Reader REST failure reaches Failed before embedding

<!-- field: use_case -->
#### Use case

rel08.0-uc002-corpus-query

<!-- field: go_test -->
#### Go test

go test ./conformance -run TestCorpusReaderRESTFailureConformance

<!-- field: expected | form: "tree" -->
#### Expected

- readiness_failure_reaches_failed: true
- embedding_does_not_run: true

<!-- field: traces -->
#### Traces

- srd012-chroma-corpus-agents R2.1
- srd012-chroma-corpus-agents R6.2
