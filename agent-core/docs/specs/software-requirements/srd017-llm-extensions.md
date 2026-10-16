---
id: srd017-llm-extensions
title: LLM Package Extensions (parseutil, profile, interfaces)
depends_on:
  - srd_id: srd010-conversation
    symbols_used: [Conversation, Message]
  - srd_id: srd004-tool-registry
    symbols_used: [Registry, ToolSpec]
  - srd_id: srd012-system-prompt-template
    symbols_used: [Envelope]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

The LLM package provides infrastructure beyond the core Client and
Conversation types: response parsing utilities for JSON repair and
validation, a YAML-driven model profile system for per-model extraction
pipelines, and domain interfaces (ResponseParser, PromptAssembler) for
extensible prompt and response handling. These capabilities need formal
specification to ensure stable contracts.

<!-- field: goals -->
## Goals

- G1: Parse utilities for JSON repair, flat parameter extraction, token estimation, and field validation.
- G2: YAML-driven model profile system for per-model response extraction pipelines.
- G3: Domain interfaces (ResponseParser, PromptAssembler) for extensible prompt and response handling.

<!-- req: R1 -->
## R1 - Parse utilities (parseutil)

<!-- req: R1.1 -->
### R1.1

FixNewlinesInStrings replaces literal newlines and tabs inside JSON string values with escape sequences.

<!-- req: R1.2 -->
### R1.2

ExtractFlatParams handles models that put parameters at the top level alongside "tool".

<!-- req: R1.3 -->
### R1.3

CountToolCallBlocks counts [tool_call] occurrences in raw LLM output.

<!-- req: R1.4 -->
### R1.4

EstimateTokens provides a rough token count by dividing content length by 4.

<!-- req: R1.5 -->
### R1.5

ExtractDoneSummary extracts the "summary" field from done-tool parameters.

<!-- req: R1.6 -->
### R1.6

CheckRequiredFields validates required JSON Schema fields are present in parameters.

<!-- req: R1.7 -->
### R1.7

ClassifyParseError returns a short category for parse error text.

<!-- req: R1.8 -->
### R1.8

Truncate shortens a string to a maximum length for logging.

<!-- req: R2 -->
## R2 - Model profiles

<!-- req: R2.1 -->
### R2.1

ProfileSpec defines response-parser metadata with name, match_prefixes, envelope, strict_format, and extraction_pipeline.

<!-- req: R2.2 -->
### R2.2

LoadProfilesFromFS loads parser profiles from the embedded profile filesystem used by DefaultProfileRegistry.

<!-- req: R2.3 -->
### R2.3

ProfileRegistry resolves profiles by model name using prefix matching, falling back to a default.

<!-- req: R2.4 -->
### R2.4

Extraction pipeline steps: strip_code_fences, strip_thinking_blocks, extract_envelope, extract_native_token, extract_braces.

<!-- req: R2.5 -->
### R2.5

DefaultProfile creates a fallback ResponseParser with envelope-based extraction.

<!-- req: R2.6 -->
### R2.6

A response-parser profile carrying a machine key SHALL be rejected because MachineSpec owns program selection.

<!-- req: R2.7 -->
### R2.7

invoke_llm tool config SHALL declare provider, provider URL, and manifest_state; supported providers are ollama and cohere, while startup SHALL reject unknown providers and missing required runtime fields. Omitted provider retains the Ollama compatibility default, but explicit Cohere SHALL NOT inherit the Ollama localhost URL.

<!-- req: R2.8 -->
### R2.8

invoke_llm tool config SHALL own response parser selection through a response_profile field that names an embedded model parser profile. When response_profile is absent, runtime parser resolution SHALL fall back to model-prefix matching and then the default parser profile.

<!-- req: R2.9 -->
### R2.9

The runtime CLI SHALL NOT expose --profiles-dir as the supported parser selection mechanism. Custom parser behavior SHALL be modeled as a parser profile YAML referenced by invoke_llm config.

<!-- req: R3 -->
## R3 - Domain interfaces and types

<!-- req: R3.1 -->
### R3.1

ResponseParser interface has ExtractToolCall(raw) and EnvelopeConfig() methods.

<!-- req: R3.2 -->
### R3.2

PromptAssembler interface has AssembleMessages(conv, registry, state) returning []Message.

<!-- req: R3.3 -->
### R3.3

ToolRequest struct represents a parsed tool call with ToolName (string, json:"tool") and Params (json.RawMessage, json:"parameters"). Used as the intermediate representation between parse_response and command dispatch.

<!-- req: R4 -->
## R4 - Extraction functions

<!-- req: R4.1 -->
### R4.1

StripCodeFences removes markdown code fences, keeping content.

<!-- req: R4.2 -->
### R4.2

ExtractWithEnvelope extracts JSON from between open/close tags, falling back to brace extraction.

<!-- req: R4.3 -->
### R4.3

StripThinkingBlocks removes <think>...</think> blocks.

<!-- req: R4.4 -->
### R4.4

ExtractBraces isolates JSON by finding the first '{' and last '}'.

<!-- req: R4.5 -->
### R4.5

MakeNativeTokenExtractor builds an extraction step for model-native tokens.

<!-- field: non_goals -->
## Non-goals

- Specific ResponseParser or PromptAssembler implementations (those live in domain agents or STL).
- LLM client implementations (srd009).
- Conversation management (srd010).

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

Parse utilities handle JSON repair, flat params, token estimation, required field checking, and error classification.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R1.4
- R1.5
- R1.6
- R1.7
- R1.8

<!-- ac: AC2 -->
### AC2

ProfileRegistry loads YAML profiles, resolves by model prefix, and provides a default. Pipeline steps are configurable. invoke_llm config owns parser selection through response_profile, with model-prefix/default fallback when omitted.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3
- R2.4
- R2.5
- R2.6
- R2.7
- R2.8
- R2.9

<!-- ac: AC3 -->
### AC3

ResponseParser and PromptAssembler are domain interfaces. yamlProfile satisfies ResponseParser at compile time. ToolRequest is the parsed tool call representation.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3

<!-- ac: AC4 -->
### AC4

Extraction functions handle code fences, envelopes, thinking blocks, braces, and native tokens.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
- R4.3
- R4.4
- R4.5

<!-- field: diagrams -->
## Diagrams

````yaml
requirements_trace: |
    ```mermaid
    flowchart LR
      Problem["Problem: The LLM package provides infrastructure beyond the co..."] --> Goals["Goals: LLM Package Extensions (parseutil, profile, int..."]
      Goals --> R1["R1: Parse utilities (parseutil)"]
      Goals --> R2["R2: Model profiles"]
      Goals --> R3["R3: Domain interfaces and types"]
      Goals --> R4["R4: Extraction functions"]
      AC1["AC1: Parse utilities handle JSON repair, flat params, token es..."]
      R1 --> AC1
      AC2["AC2: ProfileRegistry loads YAML profiles, resolves by model pr..."]
      R2 --> AC2
      AC3["AC3: ResponseParser and PromptAssembler are domain interfaces...."]
      R3 --> AC3
      AC4["AC4: Extraction functions handle code fences, envelopes, think..."]
      R4 --> AC4
    ```
````
