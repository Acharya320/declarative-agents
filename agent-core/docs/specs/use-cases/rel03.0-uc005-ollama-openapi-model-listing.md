---
id: rel03.0-uc005-ollama-openapi-model-listing
title: Qwen agent lists Ollama models through an OpenAPI REST tool
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

In the live proof, a configured agent runs with a Qwen 3 family Ollama model.
Only ollama_list_models is visible. Backed by the Ollama OpenAPI fixture, that
word calls GET /api/tags. When the user asks which local Ollama models are
available on the server, the model must call the tool instead of guessing from
prior context or model training. The final answer names live response models.

<!-- field: actor -->
## Actor

Release maintainer proving live Ollama REST tool use

<!-- field: trigger -->
## Trigger

Package-level REST tests are not enough for release closure. Before closure, a
maintainer needs proof that the tool manifest, model call, REST runtime, and
final answer line up in one run. The profile exposes a narrow model-listing
tool and keeps Ollama transport authority in configuration. MachineSpec routes
the model and REST tool signals.

<!-- field: flow -->
## Flow

- F1: The maintainer starts a local Ollama server on http://127.0.0.1:11434 with at least one Qwen 3 family model available.
- F2: The profile loads rest/openapi/ollama.yaml under AGENT_CATALOG_ROOT and exposes only operationId listOllamaModels.
- F3: The profile declares ollama_list_models as a boundary ToolDef backed by rest_client_invoke.
- F4: The selected tools manifest gives the model ollama_list_models and does not expose mutating Ollama REST operations.
- F5: The user asks the agent to list local Ollama models.
- F6: invoke_llm sends the prompt and tool manifest to the configured Qwen 3 family model.
- F7: The model emits a tool call for ollama_list_models instead of inventing model names.
- F8: rest_client_invoke calls the configured GET /api/tags operation through net/http and maps the models array into Result output.
- F9: MachineSpec returns the REST result to the model boundary for final answer generation.
- F10: The final answer names at least one model from /api/tags and does not claim models that were absent from the REST response.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd009-ollama-adapter](../software-requirements/srd009-ollama-adapter.md) [R4.1](../software-requirements/srd009-ollama-adapter.md#r41), [R4.2](../software-requirements/srd009-ollama-adapter.md#r42), [R2.1](../software-requirements/srd009-ollama-adapter.md#r21), [R2.3](../software-requirements/srd009-ollama-adapter.md#r23) | invoke_llm owns only the chat boundary; declared REST/OpenAPI words own model discovery |
| T2 | [srd024-blocking-input-tools](../software-requirements/srd024-blocking-input-tools.md) [R0.2](../software-requirements/srd024-blocking-input-tools.md#r02), [R1.1](../software-requirements/srd024-blocking-input-tools.md#r11), [R1.3](../software-requirements/srd024-blocking-input-tools.md#r13), [R2.1](../software-requirements/srd024-blocking-input-tools.md#r21), [R2.2](../software-requirements/srd024-blocking-input-tools.md#r22) | invoke_llm blocks at the model boundary and returns one Result for MachineSpec routing |
| T3 | [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [R2.1](../software-requirements/srd027-rest-tool-runtime.md#r21), [R2.4](../software-requirements/srd027-rest-tool-runtime.md#r24), [R2.5](../software-requirements/srd027-rest-tool-runtime.md#r25), [R3.1](../software-requirements/srd027-rest-tool-runtime.md#r31), [R3.3](../software-requirements/srd027-rest-tool-runtime.md#r33) | REST definitions and OpenAPI imports stay declarative and runtime input cannot choose transport authority |
| T4 | [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [R1.2](../software-requirements/srd028-rest-client-tools.md#r12), [R1.3](../software-requirements/srd028-rest-client-tools.md#r13), [R2.1](../software-requirements/srd028-rest-client-tools.md#r21), [R2.3](../software-requirements/srd028-rest-client-tools.md#r23), [R2.6](../software-requirements/srd028-rest-client-tools.md#r26), [R4.1](../software-requirements/srd028-rest-client-tools.md#r41), [R4.3](../software-requirements/srd028-rest-client-tools.md#r43) | rest_client_invoke calls a configured OpenAPI-backed operation |
| T5 | rest/openapi/ollama.yaml under AGENT_CATALOG_ROOT listOllamaModels | GET /api/tags is the only Ollama OpenAPI operation required for this proof |
| T6 | test-rel03.0-rest-tools mage integration:ollamaRest | live proof shows tool invocation and final-answer grounding |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The selected tool manifest exposes ollama_list_models and no mutating Ollama REST operation.

<!-- field: traces -->
#### Traces

- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC2](../software-requirements/srd027-rest-tool-runtime.md#ac2)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC1](../software-requirements/srd028-rest-client-tools.md#ac1)

<!-- sc: S2 -->
### S2

The live run proves the model invoked ollama_list_models before producing the final answer.

<!-- field: traces -->
#### Traces

- [srd024-blocking-input-tools](../software-requirements/srd024-blocking-input-tools.md) [AC1](../software-requirements/srd024-blocking-input-tools.md#ac1)
- [srd024-blocking-input-tools](../software-requirements/srd024-blocking-input-tools.md) [AC2](../software-requirements/srd024-blocking-input-tools.md#ac2)

<!-- sc: S3 -->
### S3

The REST call reaches GET /api/tags through the configured OpenAPI-backed operation.

<!-- field: traces -->
#### Traces

- [srd009-ollama-adapter](../software-requirements/srd009-ollama-adapter.md) [AC1](../software-requirements/srd009-ollama-adapter.md#ac1)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC3](../software-requirements/srd027-rest-tool-runtime.md#ac3)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC4](../software-requirements/srd028-rest-client-tools.md#ac4)

<!-- sc: S4 -->
### S4

The final answer includes model names returned by /api/tags and excludes names not present in the response.

<!-- field: traces -->
#### Traces

- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC7](../software-requirements/srd028-rest-client-tools.md#ac7)

<!-- field: out_of_scope -->
## Out of scope

- Implementing the Ollama REST profile files.
- Implementing the live integration target.
- Pulling or installing Qwen models.
- Exposing Ollama generation or mutation endpoints as REST tools.

<!-- field: semantic_model -->
## Semantic model

sm-rest-tool-language

<!-- field: config_format -->
## Config format

cfg-rest-tool-format
