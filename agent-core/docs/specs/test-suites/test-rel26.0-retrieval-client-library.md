---
id: test-rel26.0-retrieval-client-library
title: Test Suite for Release 26.0 Retrieval Client Library
release: "26.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- [rel26.0-uc001-retrieval-client-library](../use-cases/rel26.0-uc001-retrieval-client-library.md)

<!-- field: preconditions -->
## Preconditions

- srd060-retrieval-client-contract defines the contract, the placement, and the fragment.
- Fragment cases run from a checkout with --core-root and need no network.

<!-- field: test_cases -->
## Test cases

<!-- case: The fragment expands to the query operation for each route prefix | status: implemented -->
### The fragment expands to the query operation for each route prefix

Expand the fragment with the chatbot-mesh arguments (rag, $from(rag_unit).base_url, 5) and the agentic-wiki-mesh arguments (knowledge, $from(knowledge_unit).base_url).

<!-- field: use_case -->
#### Use case

[rel26.0-uc001-retrieval-client-library](../use-cases/rel26.0-uc001-retrieval-client-library.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest/... -run 'TestRetrievalQueryClientFragment'

<!-- field: expected | form: "tree" -->
#### Expected

- path_follows_route_prefix: true
- output_carries_five_fields: true
- missing_argument_names_parameter: true

<!-- field: traces -->
#### Traces

- [rel26.0-uc001-retrieval-client-library](../use-cases/rel26.0-uc001-retrieval-client-library.md) [S1](../use-cases/rel26.0-uc001-retrieval-client-library.md#s1)
- [srd060-retrieval-client-contract](../software-requirements/srd060-retrieval-client-contract.md) [AC1](../software-requirements/srd060-retrieval-client-contract.md#ac1)
- [srd060-retrieval-client-contract](../software-requirements/srd060-retrieval-client-contract.md) [R3.1](../software-requirements/srd060-retrieval-client-contract.md#r31)

<!-- case: Status mapping and side effect of the produced operation | status: implemented -->
### Status mapping and side effect of the produced operation

Inspect the expanded operation's success, failure, and side-effect declarations.

<!-- field: use_case -->
#### Use case

[rel26.0-uc001-retrieval-client-library](../use-cases/rel26.0-uc001-retrieval-client-library.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest/... -run 'TestRetrievalQueryClientFragment'

<!-- field: expected | form: "tree" -->
#### Expected

- status_200_signal: QueryResponded
- status_400_signal: QueryRejected
- side_effect: external_api read_only

<!-- field: traces -->
#### Traces

- [rel26.0-uc001-retrieval-client-library](../use-cases/rel26.0-uc001-retrieval-client-library.md) [S2](../use-cases/rel26.0-uc001-retrieval-client-library.md#s2)
- [srd060-retrieval-client-contract](../software-requirements/srd060-retrieval-client-contract.md) [AC2](../software-requirements/srd060-retrieval-client-contract.md#ac2)
- [srd060-retrieval-client-contract](../software-requirements/srd060-retrieval-client-contract.md) [R3.3](../software-requirements/srd060-retrieval-client-contract.md#r33)

<!-- case: The contract loads through the agent-core root with no declared root | status: implemented -->
### The contract loads through the agent-core root with no declared root

A unit importing /opt/agent-core/tools/units/types-retrieval.yaml and expanding the fragment loads with --core-root and no library root.

<!-- field: use_case -->
#### Use case

[rel26.0-uc001-retrieval-client-library](../use-cases/rel26.0-uc001-retrieval-client-library.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest/... ./internal/load/... -run 'TestRetrievalQueryClientFragment|TestLoadClosureResolvesTheShippedQueryResult'

<!-- field: expected | form: "tree" -->
#### Expected

- loads_without_declared_root: true
- no_field_beyond_core_request: true

<!-- field: traces -->
#### Traces

- [rel26.0-uc001-retrieval-client-library](../use-cases/rel26.0-uc001-retrieval-client-library.md) [S3](../use-cases/rel26.0-uc001-retrieval-client-library.md#s3)
- [srd060-retrieval-client-contract](../software-requirements/srd060-retrieval-client-contract.md) [AC3](../software-requirements/srd060-retrieval-client-contract.md#ac3)
- [srd060-retrieval-client-contract](../software-requirements/srd060-retrieval-client-contract.md) [R2.1](../software-requirements/srd060-retrieval-client-contract.md#r21)

<!-- case: chatbot-mesh queries through the fragment and receives metadatas | status: implemented -->
### chatbot-mesh queries through the fragment and receives metadatas

chatbot-mesh's chatbot declares no local query operation, its chat and degraded-RAG tests pass, and rag-server's query response carries metadatas.

<!-- field: use_case -->
#### Use case

[rel26.0-uc001-retrieval-client-library](../use-cases/rel26.0-uc001-retrieval-client-library.md)

<!-- field: go_test -->
#### Go test

mage audit

<!-- field: expected | form: "tree" -->
#### Expected

- local_query_operations: 0
- chat_behavior_unchanged: true
- metadatas_returned: true

<!-- field: traces -->
#### Traces

- [rel26.0-uc001-retrieval-client-library](../use-cases/rel26.0-uc001-retrieval-client-library.md) [S4](../use-cases/rel26.0-uc001-retrieval-client-library.md#s4)
- [srd060-retrieval-client-contract](../software-requirements/srd060-retrieval-client-contract.md) [AC4](../software-requirements/srd060-retrieval-client-contract.md#ac4)
- [srd060-retrieval-client-contract](../software-requirements/srd060-retrieval-client-contract.md) [R4.1](../software-requirements/srd060-retrieval-client-contract.md#r41)
