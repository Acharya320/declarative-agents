---
id: cfg-llm-config-format
title: LLM Configuration Format (llm/*.yaml)
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: purpose -->
## Purpose

Defines the per-model invoke_llm tool declaration format used under
each agent's llm/ directory. Each file configures a single model
variant as a builtin tool declaration with an invoke_llm config
block containing model, provider, endpoint, parser, and prompt settings.

<!-- field: schema -->
## Schema

<!-- field: description -->
### Description

Each LLM config file is a standard tool declaration file containing
a single tool entry. The tool always uses the same fixed identity
fields, with variation in the config block.

<!-- field: fields | form: "tree" -->
### Fields

1. - field: name
   - type: string
   - required: true
   - description:
     Always "invoke_llm". This is the fixed tool name that machine
     transitions reference to trigger LLM calls.
1. - field: type
   - type: string
   - required: true
   - description: Always "builtin". LLM invocation is a compiled-in capability.
1. - field: init
   - type: string
   - required: true
   - description:
     Always "invoke_llm". Maps to the invoke_llm factory in the
     BuiltinRegistry.
1. - field: visibility
   - type: string
   - required: true
   - description:
     Always "internal". The invoke_llm tool is called by machine
     transitions, not by the LLM itself. Hiding it from the LLM's
     tool list prevents recursive self-invocation.
1. - field: config
   - type: object
   - required: true
   - description:
     Model-specific configuration passed to the invoke_llm factory.
     This is where each file differs.
   - sub_fields:
     1. - field: model
        - type: string
        - required: true
        - description:
          Model name as recognized by the provider. Examples include
          "deepseek-r1:8b", "devstral:24b", and "gemma3:27b".
     1. - field: provider
        - type: string
        - required: false
        - description:
          Legacy provider name, ollama or cohere. It resolves to the shipped
          library's chat dialect under the agent-core install root,
          tools/providers/<provider>/chat-dialect.yaml, and omission with no
          dialect resolves to Ollama's. Superseded by dialect (srd058 R3.2);
          naming both is an error.
     1. - field: dialect
        - type: string
        - required: false
        - description:
          Path of the chat dialect the call reads its request body, reply
          selectors, failure map, auth profile, and parser profiles from
          (srd058 R2). It resolves like an import: relative to the
          declaring file, or under a library root, normally
          /opt/providers/chat-dialect.yaml with the providers root the
          profile declares. The loader makes the file a closure file.
          provider_url is required with a dialect.
     1. - field: provider_url
        - type: string
        - required: true
        - description:
          Base URL for the provider API. Local Ollama usually uses
          "http://localhost:11434"; Cohere commonly uses
          "${COHERE_API_URL:-https://api.cohere.com}". Ollama keeps its
          localhost default when provider is omitted. Explicit provider
          cohere requires this field and never inherits the Ollama default.
     1. - field: manifest_state
        - type: string
        - required: true
        - description:
          Machine state whose external tools are included in the
          invoke_llm manifest. For generator and planner machines this
          is currently "Composing". Declaring this in YAML keeps LLM
          manifest construction tied to the machine contract rather
          than hardcoded in cmd/agent.
     1. - field: response_profile
        - type: string
        - required: false
        - status: current
        - description:
          Name of the LLM response parser profile to use for model output
          extraction. The value selects one parser profile from
          internal/model/llm/profiles/*.yaml.
          This field is owned by the invoke_llm tool declaration that the
          selected agent profile loads. It is not a process-wide CLI
          override.
        - fallback:
          When omitted, invoke_llm resolves the parser from the model name
          by matching ProfileSpec.match_prefixes and falls back to the
          default parser profile when no prefix matches.
     1. - field: system_prompt
        - type: string
        - required: false
        - description:
          System prompt prepended to the conversation. Sets the
          agent's persona, constraints, and output format
          expectations.
     1. - field: tool_prompt
        - type: string
        - required: false
        - description:
          Additional prompt section describing available tools and
          how to format tool calls. Used by the planner agent to
          inject tool-use instructions into the LLM context.
     1. - field: max_time
        - type: int
        - required: false
        - description:
          Legacy per-call timeout in seconds. Used only when llm_timeout is
          absent. It never changes the machine run budget.
     1. - field: context_limit
        - type: int
        - required: false
        - default: 0
        - description:
          Hard local preflight ceiling over the assembled prompt's estimated
          tokens. A positive limit emits CommandError before provider Chat
          when the estimate exceeds it. Zero disables the estimate check;
          num_ctx remains the provider's actual context-window request.
     1. - field: num_ctx
        - type: int
        - required: false
        - description:
          Context window size passed only to the Ollama provider.
          Controls how many tokens the model can consider.
          The Cohere dialect does not reference this field, so its request omits it.
     1. - field: temperature
        - type: float
        - required: false
        - default: 0
        - description:
          Sampling temperature passed to the provider as a decoding
          parameter. Omit for the deterministic default of 0 (greedy
          decoding). Set a positive value for sampled decoding, for example
          when a downstream experiment needs run-to-run trajectory variance.
          The effective value is recorded on the invoke_llm span as
          gen_ai.request.temperature.
     1. - field: seed
        - type: int
        - required: false
        - default: 42
        - description:
          Random seed passed to the provider as a decoding parameter. Omit
          for the deterministic default of 42, which keeps sampled runs
          reproducible. The effective value is recorded on the invoke_llm
          span as gen_ai.request.seed.
     1. - field: llm_timeout
        - type: int
        - required: false
        - description:
          Per-LLM-call timeout in seconds. Individual calls that
          exceed this limit are cancelled and retried. When both llm_timeout
          and max_time are present, llm_timeout takes precedence.
     1. - field: max_tokens
        - type: int
        - required: false
        - description:
          Deprecated compatibility input. It does not change the run budget;
          declare budget.max_tokens on the machine instead.
     1. - field: ollama_url
        - type: string
        - required: false
        - description:
          Alias for provider_url. If both are set, provider_url
          takes precedence. Valid only for the Ollama compatibility path.

<!-- field: provider_behavior -->
## Provider behavior

<!-- field: status -->
### Status

implemented

<!-- field: ollama -->
### Ollama

invoke_llm sends non-streaming POST /api/chat. Request decoding options
live under the Ollama options object, including num_ctx. No credential is
required. An omitted provider selects this compatibility path and defaults
provider_url to http://localhost:11434.

<!-- field: cohere -->
### Cohere

invoke_llm sends non-streaming POST /v2/chat. model, messages,
temperature, and seed are top-level request fields. The adapter joins
ordered text blocks from message.content and maps usage.tokens into the
provider-neutral ChatResponse. num_ctx and Ollama request fields are not
sent.

<!-- field: credential -->
### Credential

Cohere Chat resolves COHERE_API_KEY from the process environment only when
Chat executes and sends it as Authorization Bearer. No API key or
credential-valued field is valid in this YAML format. Registration performs
no credential lookup.

<!-- field: response_profile_boundary -->
### Response profile boundary

response_profile operates after provider HTTP normalization. The cohere
profile extracts a [tool_call] envelope or bare JSON from the joined model
text; it does not parse the Cohere HTTP response object.

<!-- field: parse_response_config -->
## Parse response config

<!-- field: purpose -->
### Purpose

Defines the validation state and parser profile owned by one parse_response
ToolDef. These values do not come from mutable invoke_llm registration
state.

<!-- field: fields | form: "tree" -->
### Fields

1. - field: manifest_state
   - type: string
   - required: true when the parse word feeds $tool
   - description:
     The `$tool` target state whose external vocabulary validates the parsed
     tool name. Startup links the actual selector, parser, and dynamic
     transition and rejects a different value.
1. - field: response_profile
   - type: string
   - required: false
   - description:
     Explicit parser profile for this parse word. When omitted, the existing
     resolved model parser remains the compatibility fallback.

<!-- field: per_model_pattern -->
## Per model pattern

Each file under an agent's llm/ directory configures a different
model variant. The file name identifies the variant but has no semantic
meaning to the system. All configuration comes from the fields inside.

Parser behavior is part of the same invoke_llm config. A model variant can
set response_profile to choose the parser profile by name. Custom parser
behavior should be added as a parser profile YAML under
internal/model/llm/profiles/ and referenced from the agent's llm/*.yaml file.
Operators should not select parser behavior with --profiles-dir.

The evaluator uses these files to build a benchmark grid across samples and
models. By loading multiple llm/*.yaml files, the
evaluator can run the same coding task against different models and
compare their outputs. The generator and planner typically load a
single variant selected at invocation time.

<!-- field: profile_terms -->
## Profile terms

Agent profiles differ from response parser profiles.

Agent profiles are profile-owned assets. In normal deployments they live under
an external applications/catalog checkout or bundle, for example
$AGENT_CATALOG_ROOT/agents/*/profile.yaml. They select the machine,
tools.yaml files, declaration directories, and agent-local declaration files
that define one configured agent program.

Response parser profiles live under internal/model/llm/profiles/*.yaml. They
define output extraction for model families: match_prefixes, envelope,
strict_format, and extraction_pipeline. invoke_llm uses one response parser
profile to turn raw model text into a candidate tool call for parse_response.

<!-- field: config_files -->
## Config files

<!-- config_file: $AGENT_CATALOG_ROOT/agents/executor/llm/ -->
### $AGENT_CATALOG_ROOT/agents/executor/llm/

<!-- field: variant_count | form: "scalar" -->
#### Variant count

9

<!-- field: description -->
#### Description

Model variants for the generator agent: default, deepseek,
devstral, gemma4, glm, granite, laguna, mistral-small, and
nemotron. Each configures a different local model via Ollama
with generator-specific system prompts.

<!-- config_file: $AGENT_CATALOG_ROOT/agents/critic/llm/ -->
### $AGENT_CATALOG_ROOT/agents/critic/llm/

<!-- field: variant_count | form: "scalar" -->
#### Variant count

10

<!-- field: description -->
#### Description

Model variants for the evaluator's benchmark grid:
deepseek-16b, devstral, gemma4, glm, granite, laguna,
mistral-small, nemotron, qwen3-8b, and qwen3-35b. These
are harness profiles used by child generator processes, not
by the evaluator state machine itself.

<!-- config_file: $AGENT_CATALOG_ROOT/agents/planner/llm/ -->
### $AGENT_CATALOG_ROOT/agents/planner/llm/

<!-- field: variant_count | form: "scalar" -->
#### Variant count

1

<!-- field: description -->
#### Description

Single variant (default.yaml) for the planner agent. Includes
both system_prompt and tool_prompt fields since the planner
injects tool-use instructions into the LLM context.

<!-- field: related_documents -->
## Related documents

- [cfg-tool-declaration-format](tool-declaration-format.md)
- srd017
- [srd048-cohere-adapter](../software-requirements/srd048-cohere-adapter.md)
- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md)

<!-- field: implementation -->
## Implementation

- internal/tools/llm
- internal/tools/llm/dialect
