---
id: test-rel15.0-cohere-invoke
title: Test Suite for Release 15.0 Cohere invoke_llm
release: "15.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel15.0-uc001-cohere-invoke

<!-- field: preconditions -->
## Preconditions

- srd048-cohere-adapter defines the Cohere v2 wire, credential, normalization, provider-selection, error, and telemetry contracts.
- Tests use httptest and an injected API-key lookup; no live Cohere credential or service is required.
- The existing Ollama adapter tests remain the compatibility baseline.

<!-- field: test_cases -->
## Test cases

<!-- case: Cohere constructor is inert and configuration selects the adapter | status: implemented -->
### Cohere constructor is inert and configuration selects the adapter

Build provider cohere against an unreachable URL without supplying an
API key. Construction succeeds and performs no network or credential
lookup.

<!-- field: use_case -->
#### Use case

rel15.0-uc001-cohere-invoke

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm -run 'Test(NewInvokeLLMBuilder_CohereNoProbe|InvokeLLMBuildsWithoutNetworkIO)'

<!-- field: expected -->
#### Expected

```yaml
builder_created: true
network_requests: 0
credential_lookups: 0
```

<!-- field: traces -->
#### Traces

- rel15.0-uc001-cohere-invoke S1
- srd048-cohere-adapter AC1
- srd048-cohere-adapter AC4

<!-- case: Cohere Chat sends the declared v2 request | status: implemented -->
### Cohere Chat sends the declared v2 request

Execute Chat against httptest with an injected key. Inspect method, path,
bearer header, model, chronological messages, temperature, and seed.

<!-- field: use_case -->
#### Use case

rel15.0-uc001-cohere-invoke

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm -run TestProviderDialectParity/cohere

<!-- field: expected -->
#### Expected

```yaml
method_path: POST /v2/chat
bearer_present: true
request_fields: [model, messages, temperature, seed]
forbidden_fields: [options, num_ctx]
request_count: 1
```

<!-- field: traces -->
#### Traces

- rel15.0-uc001-cohere-invoke S2
- srd048-cohere-adapter AC1
- srd048-cohere-adapter AC2

<!-- case: Cohere response joins text and maps usage | status: implemented -->
### Cohere response joins text and maps usage

Return multiple text blocks, a non-text block, and usage.tokens. The
adapter preserves text order, excludes non-text data, and maps usage.

<!-- field: use_case -->
#### Use case

rel15.0-uc001-cohere-invoke

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm ./internal/tools/llm/dialect -run 'Test(ProviderDialectParity|ExtractReplyJoinsSelectedTextAndReadsUsage)'

<!-- field: expected -->
#### Expected

```yaml
content: joined ordered text blocks
non_text_exposed: false
token_fields: [TokensIn, TokensOut]
missing_tokens_default_zero: true
```

<!-- field: traces -->
#### Traces

- rel15.0-uc001-cohere-invoke S3
- srd048-cohere-adapter AC3

<!-- case: Credential and provider failures remain safe and single-attempt | status: implemented -->
### Credential and provider failures remain safe and single-attempt

Exercise missing key, cancellation, connection failure, non-2xx,
malformed JSON, missing text, and an unreachable invoke_llm dispatch.

<!-- field: use_case -->
#### Use case

rel15.0-uc001-cohere-invoke

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm -run 'Test(InvokeLLM_CohereCommandError|DialectMissingCredentialFailsBeforeNetwork|DialectCredentialAbsentFromDiagnostics|ProviderDialectFailureKeepsCommandError)'

<!-- field: expected -->
#### Expected

```yaml
errors_distinct: true
credential_value_exposed: false
retries: 0
invoke_signal: CommandError
```

<!-- field: traces -->
#### Traces

- rel15.0-uc001-cohere-invoke S5
- rel15.0-uc001-cohere-invoke S6
- srd048-cohere-adapter AC2
- srd048-cohere-adapter AC4
- srd048-cohere-adapter AC5

<!-- case: Provider validation preserves Ollama defaults | status: implemented -->
### Provider validation preserves Ollama defaults

Table-drive omitted, ollama, cohere, and unknown providers plus missing
URLs. Omitted provider remains Ollama; Cohere requires its own URL.

<!-- field: use_case -->
#### Use case

rel15.0-uc001-cohere-invoke

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm -run TestInvokeLLMProviderSelection

<!-- field: expected -->
#### Expected

```yaml
omitted_provider: ollama
omitted_provider_url: http://localhost:11434
cohere_without_url_rejected: true
unknown_provider_rejected: true
```

<!-- field: traces -->
#### Traces

- rel15.0-uc001-cohere-invoke S1
- rel15.0-uc001-cohere-invoke S6
- srd048-cohere-adapter AC4

<!-- case: Cohere parser profile handles joined classifier output | status: implemented -->
### Cohere parser profile handles joined classifier output

Resolve the profile by name and command- prefix, then extract both a
tool_call envelope and bare JSON. Existing profile resolutions remain
unchanged.

<!-- field: use_case -->
#### Use case

rel15.0-uc001-cohere-invoke

<!-- field: go_test -->
#### Go test

go test ./internal/model/llm -run TestCohereResponseProfile

<!-- field: expected -->
#### Expected

```yaml
named_resolution: cohere
prefix_resolution: command-r7b-12-2024
envelope_and_bare_json_extracted: true
```

<!-- field: traces -->
#### Traces

- rel15.0-uc001-cohere-invoke S4
- srd048-cohere-adapter AC3
- srd048-cohere-adapter AC4

<!-- case: Cohere transport emits GenAI attributes with caller parenting | status: implemented -->
### Cohere transport emits GenAI attributes with caller parenting

Execute success and error calls under a caller span and inspect the
transport child and dispatch attributes.

<!-- field: use_case -->
#### Use case

rel15.0-uc001-cohere-invoke

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm ./internal/observability/tracing ./internal/observability/telemetry/genai -run 'Test(InvokeLLMRunsTheDialectPipeline|ContextChildParentsOnTheCallerSpan|ProviderCohere)'

<!-- field: expected -->
#### Expected

```yaml
provider: cohere
attributes: [gen_ai.request.model, gen_ai.response.model, server.address, server.port, gen_ai.usage.input_tokens, gen_ai.usage.output_tokens]
caller_parent_preserved: true
error_type_recorded: true
```

<!-- field: traces -->
#### Traces

- rel15.0-uc001-cohere-invoke S5
- srd048-cohere-adapter AC5

<!-- case: Existing Ollama provider behavior remains compatible | status: implemented -->
### Existing Ollama provider behavior remains compatible

Run the existing Ollama adapter and invoke_llm tests unchanged after the
provider switch is added.

<!-- field: use_case -->
#### Use case

rel15.0-uc001-cohere-invoke

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm

<!-- field: expected -->
#### Expected

```yaml
existing_tests_pass: true
ollama_wire_unchanged: true
```

<!-- field: traces -->
#### Traces

- rel15.0-uc001-cohere-invoke S1
- srd048-cohere-adapter AC4
