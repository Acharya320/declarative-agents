---
id: test-rel22.0-provider-libraries
title: Test Suite for Release 22.0 Provider Libraries
release: "22.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel22.0-uc001-provider-libraries

<!-- field: preconditions -->
## Preconditions

- srd058-provider-libraries defines the provider library, the chat dialect, invoke_llm as a template method, the embed and rerank stage fragments, and the conformance suite.
- srd056 declared roots and the --library override, srd052 REST and stage fragments, and the REST request-body renderer are available and unchanged by this release.
- Fake HTTP servers answer the Ollama, Cohere, and OpenAI-shaped chat, embed, and rerank endpoints with recorded bodies and chosen statuses; no test reaches a real provider.

<!-- field: test_cases -->
## Test cases

<!-- case: A strict dialect loads and a bad one fails with its path | status: implemented -->
### A strict dialect loads and a bad one fails with its path

Load a fixture dialect; then copies with an unknown field, a credential value in place of token_ref, and messages referenced inside a larger string.

<!-- field: use_case -->
#### Use case

rel22.0-uc001-provider-libraries

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm/... -run 'TestChatDialectLoadsStrictly|TestChatDialectRejectsUnknownField|TestChatDialectRejectsCredentialLiteral|TestChatDialectRejectsMessagesInsideString'

<!-- field: expected | form: "tree" -->
#### Expected

- valid_dialect_loads: true
- unknown_field_rejected: true
- credential_literal_rejected: true
- embedded_messages_reference_rejected: true

<!-- field: traces -->
#### Traces

- rel22.0-uc001-provider-libraries S1
- srd058-provider-libraries AC1
- srd058-provider-libraries R1.5
- srd058-provider-libraries R2.1
- srd058-provider-libraries R2.2

<!-- case: A missing library file fails at load naming path and directory | status: implemented -->
### A missing library file fails at load naming path and directory

Bind a providers root that lacks chat-dialect.yaml, then one that lacks rerank-fragment.yaml under a profile that expands the rerank stage.

<!-- field: use_case -->
#### Use case

rel22.0-uc001-provider-libraries

<!-- field: go_test -->
#### Go test

go test ./internal/load -run 'TestProviderLibraryMissingDialectNamesPathAndDirectory|TestProviderLibraryMissingFragmentNamesPathAndDirectory'

<!-- field: expected | form: "tree" -->
#### Expected

- load_fails: true
- error_names_missing_path: true
- error_names_bound_directory: true

<!-- field: traces -->
#### Traces

- rel22.0-uc001-provider-libraries S1
- srd058-provider-libraries AC1
- srd058-provider-libraries R1.3

<!-- case: Rebinding the providers root switches provider with no declaration edit | status: implemented -->
### Rebinding the providers root switches provider with no declaration edit

Load one profile with its declared Ollama root and again with --library providers set to the Cohere library; compare the machine, tool, REST, and fragment bytes, the resolved dialect, the dump, and the program identity.

<!-- field: use_case -->
#### Use case

rel22.0-uc001-provider-libraries

<!-- field: go_test -->
#### Go test

go test ./internal/load ./internal/tools/llm/... -run 'TestLibraryRebindSwitchesProvider|TestDumpRecordsDialectRoot'

<!-- field: expected | form: "tree" -->
#### Expected

- declaration_bytes_identical: true
- resolved_dialect_differs: true
- dump_records_dialect_path_and_root: true
- program_identity_differs: true

<!-- field: traces -->
#### Traces

- rel22.0-uc001-provider-libraries S1
- srd058-provider-libraries AC1
- srd058-provider-libraries R1.1
- srd058-provider-libraries R1.2
- srd058-provider-libraries R2.3

<!-- case: The dialect and provider fields do not combine | status: implemented -->
### The dialect and provider fields do not combine

Configure invoke_llm with provider and dialect together; then with provider ollama and provider cohere alone.

<!-- field: use_case -->
#### Use case

rel22.0-uc001-provider-libraries

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm/... -run 'TestInvokeLLMRejectsProviderAndDialect|TestInvokeLLMMapsLegacyProvider'

<!-- field: expected | form: "tree" -->
#### Expected

- both_fields_rejected: true
- legacy_provider_resolves_to_shipped_dialect: true

<!-- field: traces -->
#### Traces

- rel22.0-uc001-provider-libraries S1
- srd058-provider-libraries AC1
- srd058-provider-libraries R3.2

<!-- case: The dialect path matches the compiled adapters | status: implemented -->
### The dialect path matches the compiled adapters

Run recorded multi-turn conversations through the Ollama and Cohere adapters and through their dialect files against the same fake servers; compare request JSON documents, histories, emitted signals, and span attributes.

<!-- field: use_case -->
#### Use case

rel22.0-uc001-provider-libraries

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm/... -run 'TestProviderDialectParity'

<!-- field: expected | form: "tree" -->
#### Expected

- request_documents_equal: true
- conversation_placed_as_array: true
- histories_equal: true
- signals_equal: true
- span_attributes_equal: true

<!-- field: traces -->
#### Traces

- rel22.0-uc001-provider-libraries S2
- srd058-provider-libraries AC2
- srd058-provider-libraries R2.2
- srd058-provider-libraries R3.1
- srd058-provider-libraries R3.3
- srd058-provider-libraries R5.1
- srd058-provider-libraries R5.2

<!-- case: Credentials resolve at call time and never leak | status: implemented -->
### Credentials resolve at call time and never leak

Call through a dialect whose token_ref names a variable set only through the injected resolver; then with the variable empty; then inspect output, errors, and spans.

<!-- field: use_case -->
#### Use case

rel22.0-uc001-provider-libraries

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm/... -run 'TestDialectCredentialResolvedAtCallTime|TestDialectMissingCredentialFailsBeforeNetwork|TestDialectCredentialAbsentFromDiagnostics'

<!-- field: expected | form: "tree" -->
#### Expected

- bearer_sent_from_resolver: true
- empty_credential_fails_before_network: true
- credential_absent_from_output_errors_and_spans: true

<!-- field: traces -->
#### Traces

- rel22.0-uc001-provider-libraries S2
- srd058-provider-libraries AC2
- srd058-provider-libraries R1.5
- srd058-provider-libraries R3.5

<!-- case: A library parser profile wins over the embedded registry | status: implemented -->
### A library parser profile wins over the embedded registry

Bind a library whose parser profile shares a name with an embedded profile and changes one extraction rule.

<!-- field: use_case -->
#### Use case

rel22.0-uc001-provider-libraries

<!-- field: go_test -->
#### Go test

go test ./internal/model/llm ./internal/tools/llm/... -run 'TestLibraryParserProfileOverridesEmbedded|TestDialectLibraryParserProfileResolves'

<!-- field: expected | form: "tree" -->
#### Expected

- library_profile_used: true

<!-- field: traces -->
#### Traces

- rel22.0-uc001-provider-libraries S2
- srd058-provider-libraries AC2
- srd058-provider-libraries R3.4

<!-- case: Stage fragments splice with each library and route every failure | status: implemented -->
### Stage fragments splice with each library and route every failure

Expand the embed-query, embed-document, and rerank stages in a fixture machine beside each library's REST fragments; answer 200, 401, 429, and 503 from the fake server.

<!-- field: use_case -->
#### Use case

rel22.0-uc001-provider-libraries

<!-- field: go_test -->
#### Go test

go test ./internal/load ./internal/tools/rest -run 'TestEmbedStageFragmentsSpliceWithEachLibrary|TestProviderFailureSignalsRouted'

<!-- field: expected | form: "tree" -->
#### Expected

- label_and_exhaustiveness_checks_pass: true
- success_signals_emitted:
  1. QueryEmbedded
  1. DocumentEmbedded
  1. Reranked
- failure_signals_emitted:
  1. ProviderUnauthorized
  1. ProviderThrottled
  1. ProviderUnavailable

<!-- field: traces -->
#### Traces

- rel22.0-uc001-provider-libraries S3
- srd058-provider-libraries AC3
- srd058-provider-libraries R1.4
- srd058-provider-libraries R4.1

<!-- case: No hand-written provider embedding remains | status: implemented -->
### No hand-written provider embedding remains

Scan the monorepo's REST declarations for provider embed and rerank operations declared outside a provider library; record mage stats:reuse before and after the adoption.

<!-- field: use_case -->
#### Use case

rel22.0-uc001-provider-libraries

<!-- field: go_test -->
#### Go test

go test ./internal/declstyle -run 'TestNoHandWrittenProviderEmbedding|TestHandWrittenProviderEmbeddingIsFound'

<!-- field: expected | form: "tree" -->
#### Expected

- hand_written_embed_blocks: 0
- embedding_openings_are_stage_instances: true
- ollama_path_declares_failures: true

<!-- field: traces -->
#### Traces

- rel22.0-uc001-provider-libraries S3
- srd058-provider-libraries AC3
- srd058-provider-libraries R4.2

<!-- case: A synthetic third provider needs no Go and no agent edit | status: implemented -->
### A synthetic third provider needs no Go and no agent edit

Bind the OpenAI-shaped testdata library, a chat dialect with its parser profile and two embed fragments, to the unchanged invoke_llm word, the unchanged embed stage words, and the unchanged ollama-rest profile; run chat, embedding, and every failure status against a fake OpenAI-shaped server; and check that the library is YAML only and that no declaration outside it names the provider. The embedding is read out of the data array by an indexing selector (srd038 R2.20).

<!-- field: use_case -->
#### Use case

rel22.0-uc001-provider-libraries

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm ./internal/tools/rest ./internal/load -run 'TestSyntheticProviderConformance|TestSyntheticProviderFailuresUseTheSharedTaxonomy|TestSyntheticProviderEmbeddingIsIndexedOutOfItsArray|TestSyntheticProviderTouchesOnlyItsLibrary|TestEmbedStageFragmentsSpliceWithEachLibrary'

<!-- field: expected | form: "tree" -->
#### Expected

- chat_passes: true
- embedding_indexed_out_of_data_array: true
- library_parser_profile_used: true
- shared_vocabulary_emitted: true
- go_changes_required: 0
- agent_declaration_changes_required: 0

<!-- field: traces -->
#### Traces

- rel22.0-uc001-provider-libraries S4
- srd058-provider-libraries AC4
- srd058-provider-libraries R1.3
- srd058-provider-libraries R1.4
- srd058-provider-libraries R5.3
