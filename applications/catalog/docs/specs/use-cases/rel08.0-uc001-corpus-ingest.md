---
id: rel08.0-uc001-corpus-ingest
title: Corpus Ingest into Chroma
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

The ingest profile loads a repository corpus into a local Chroma collection.
For each document it reads the file, computes an embedding at the configured
provider, and writes the document, id, and agent-supplied vector to Chroma,
then verifies the collection item count. The canonical profile defaults to
Ollama. Reusable machine, declarations, and tool selection belong in
applications/catalog. A consuming wrapper supplies its REST endpoints,
provider binding, collection/model values, and corpus fixture;
integration:chroma executes the canonical library program by reference.

Chroma does not embed text server-side over its raw HTTP API, so the vector is
computed by the embed word and threaded into the Chroma add word using the
agent-core previous-Result parameter threading (srd028 R12). No embedding
vector is authored by a model; it rides between boundary words.

<!-- field: actor -->
## Actor

Profile maintainer or operator ingesting a repository corpus.

<!-- field: trigger -->
## Trigger

A maintainer runs the mesh-owned Chroma integration with a running Chroma
server and reachable Ollama; its wrapper selects this canonical profile:

  cd applications/chatbot-mesh && mage integration:chroma

<!-- field: preconditions -->
## Preconditions

- A Chroma server is running against a persistence directory (chroma run --path <dir>), started by the harness or operator, not the agent.
- The chat and embedding provider is reachable; the canonical profile binds the Ollama library (srd058).
- The configured embedding provider is reachable and serves the selected model family; the canonical default is Ollama.
- agents/knowledge-manager/corpus-ingest profile validates and owns the Chroma and default Ollama REST definitions and the corpus selection.

<!-- field: flow -->
## Flow

- F1: The ingest machine runs chroma_ready, failing fast when Chroma is unreachable; a provider failure surfaces on the first embed.
- F2: For each corpus document the machine reads the file content and identity.
- F3: embed_document calls the configured REST client and globally unique document-embedding operation, maps the returned flat vector or one-row matrix into its Result output, and carries the document text and id forward.
- F4: normalize_document_embedding preserves a non-empty flat numeric vector or unwraps exactly one non-empty numeric row while retaining the carried document and id. Invalid shapes fail before the create-or-get collection word, which then create-or-gets a plain collection with no server-side embedding function, resolves its id, and carries the normalized embedding, document, and id forward so single-hop threading reaches the add body.
- F5: chroma_add threads the resolved collection id, the vector, and the carried document and id into the add request body via body_source previous_result and input_mapping, placing the vector in a single-token array element so it renders as the nested array Chroma expects.
- F6: When the corpus is exhausted the machine re-resolves the collection id and chroma_count reads the collection item count to verify the added documents.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd012-chroma-corpus-agents AC1 | The ingest profile is a configured Chroma retrieval boundary machine, not a Go package. |
| T2 | srd012-chroma-corpus-agents AC2 | chroma_ready, chroma_collection, library-bound embed_document, chroma_add, and chroma_count are the named REST vocabulary with create-or-get identity resolution. |
| T3 | srd012-chroma-corpus-agents AC3 | The embed-to-add vector step uses srd028 R12 body_source previous_result, input_mapping, carry_forward, and single-token reshape, with no model-authored vector. |
| T4 | srd012-chroma-corpus-agents AC5 | The running Chroma server and persistence directory are external preconditions surfaced through chroma_ready. |
| T5 | applications/chatbot-mesh integration:chroma | The application integration runs the canonical library machine and declarations with mesh REST parameterization and verifies ingest count. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Every Chroma add carries one flat agent-supplied embedding vector normalized from a flat or single-row provider response and threaded from the embed word, not a model-authored or server-side vector.

<!-- field: traces -->
#### Traces

- srd012-chroma-corpus-agents R3.2
- srd012-chroma-corpus-agents R3.4
- srd012-chroma-corpus-agents R3.5

<!-- sc: S2 -->
### S2

The collection is create-or-got without a server-side embedding function and its id is resolved into the add and count paths.

<!-- field: traces -->
#### Traces

- srd012-chroma-corpus-agents R2.2

<!-- sc: S3 -->
### S3

The post-ingest chroma_count matches the number of documents added.

<!-- field: traces -->
#### Traces

- srd012-chroma-corpus-agents R2.5
- srd012-chroma-corpus-agents R6.1

<!-- sc: S4 -->
### S4

A consuming wrapper selects a non-Ollama document-embedding client and operation whose one-row response reaches Chroma as exactly one vector, without copying the machine or overriding the catalog ToolDef.

<!-- field: traces -->
#### Traces

- srd012-chroma-corpus-agents R2.6
- srd012-chroma-corpus-agents R4.5

<!-- field: out_of_scope -->
## Out of scope

- Starting or stopping the Chroma or Ollama servers from the agent.
- Server-side embedding or a new embedding builtin word.
- Ranking or answer generation, which belong to the query use case.

<!-- field: test_evidence -->
## Test evidence

```yaml
target: applications/chatbot-mesh mage integration:chroma
conformance: go test ./conformance -run 'TestCorpusIngest(ListsTrustedCorpusBeforeModelControl|SelectsConfiguredEmbeddingProvider)'
```
