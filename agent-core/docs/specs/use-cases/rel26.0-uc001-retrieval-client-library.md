---
id: rel26.0-uc001-retrieval-client-library
title: A chatbot queries its retrieval servers through the shared agent-core client
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A chatbot author gives a chatbot retrieval by expanding the agent-core
retrieval client fragment with the application's route prefix, the
selector for the current server's base URL, the selector for the query
vector, and a result count. The chatbot declares no query operation of its
own. Each server answers with the shared QueryResult, including the
metadatas a citation resolves against. An application whose server adds
clearance or facets writes that operation itself and imports only the
shared result type.

<!-- field: actor -->
## Actor

A chatbot author wiring retrieval into a chatbot, and the maintainer of a retrieval server.

<!-- field: trigger -->
## Trigger

A chatbot needs to query one or more retrieval servers, or a retrieval
server's response shape changes.

<!-- field: flow -->
## Flow

- F1: The chatbot's REST unit expands /opt/agent-core/tools/rest/units/retrieval-query-client-fragment.yaml with route_prefix, base_url_selector, input_selector, n_results, and limits_ref (srd060 R2.1, R3.1).
- F2: The expansion produces the query operation, and the chatbot keeps its fan-out, base URLs, limits, and allowlist (srd060 R3.2, R4.1).
- F3: A server answers 200 with ids, documents, distances, metadatas, and embedding_model, and the operation emits QueryResponded (srd060 R1.2, R3.3).
- F4: A server holding another embedding space answers 400, the operation emits QueryRejected, and the machine excludes that source (srd060 R1.3).
- F5: An application whose server requires clearance writes its own query operation and imports types-retrieval.yaml for the result (srd060 R3.4).

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd060-retrieval-client-contract R1.1, R1.2, R1.3 | the request, the QueryResult, and the rejection mapping |
| T2 | srd060-retrieval-client-contract R2.1, R2.2 | agent-core placement reached through /opt/agent-core |
| T3 | srd060-retrieval-client-contract R3.1, R3.2, R3.3, R3.4 | the fragment's parameters, operation, and fixed shape |
| T4 | srd060-retrieval-client-contract R4.1, R4.2 | consumers declare no query operation; two importers |
| T5 | srd056-unit-libraries R4.1, R4.2 | the service-library convention this contract instantiates |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Expanding the fragment yields the expected query operation and output, and a missing argument fails naming the parameter.

<!-- field: traces -->
#### Traces

- srd060-retrieval-client-contract AC1

<!-- sc: S2 -->
### S2

The operation maps 200, 400, and other failures to QueryResponded, QueryRejected, and CommandError with a read-only side effect.

<!-- field: traces -->
#### Traces

- srd060-retrieval-client-contract AC2

<!-- sc: S3 -->
### S3

The contract is reached through /opt/agent-core with no declared root and carries no field beyond the core request.

<!-- field: traces -->
#### Traces

- srd060-retrieval-client-contract AC3

<!-- sc: S4 -->
### S4

chatbot-mesh queries through the fragment with unchanged behavior and receives metadatas.

<!-- field: traces -->
#### Traces

- srd060-retrieval-client-contract AC4

<!-- field: out_of_scope -->
## Out of scope

- Server promotion, clearance, facets, and document reads.
- The consumer repositories' adoption, which lands at their next pin bump.
