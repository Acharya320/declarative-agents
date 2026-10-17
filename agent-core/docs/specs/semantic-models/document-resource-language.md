---
id: sm-document-resource-language
title: Document Resource Language Semantic Model
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: purpose -->
## Purpose

Models filesystem documentation corpora as declarative resources. A profile
names the corpus root, allowed paths, response modes, and UI placement.
Generic document resource code reads files safely. MachineSpec and REST
definitions decide when the resource is used.

<!-- field: requirements_source | form: "tree" -->
## Requirements source

- canonical:
  1. docs/specs/software-requirements/srd031-declarative-document-resources.md
- supporting:
  1. docs/specs/software-requirements/srd027-rest-tool-runtime.md
  1. docs/specs/software-requirements/srd030-machine-request-rest-binding.md
  1. docs/specs/semantic-models/machine-request-http.md
  1. docs/specs/semantic-models/rest-tool-language.md
  1. docs/guides/declarative-programming-pattern.md

<!-- field: primary_entities | form: "typed table" -->
## Primary entities

| name | kind | responsibility | owns | references |
| --- | --- | --- | --- | --- |
| DocumentResource | declarative_filesystem_resource | Defines a named read-only documentation corpus with root, include rules, extension allowlist, response modes, parser policy, and ordered presentation category rules. | [corpus_authority, response_modes, category_rules] |  |
| CorpusRoot | trusted_filesystem_root | Names the base directory for document reads and list operations. | [root_path, root_escape_policy] |  |
| DocumentSelector | runtime_resource_input | Carries declared path variables, query fields, or machine request input used to select a document. | [relative_path, query_values] |  |
| PathPolicy | safety_contract | Defines normalization, traversal denial, symlink handling, size limits, binary policy, and extension allowlists. | [path_safety] |  |
| DocumentResponseMode | output_shape | Defines configured output shape for raw YAML, parsed YAML, raw Markdown, parsed Markdown, or combined raw and parsed content. | [content_shape, parser_policy] |  |
| DocumentResourceWord | vocabulary_word | Executes one declared document operation such as list, get, search, validate, or suggest. | [operation_signal, result_output] | [DocumentResource, ToolDef, MachineSpec] |
| UIActorConfig | human_actor_configuration | Defines Knowledge Manager UI assets, labels, route/action mapping, grouping, and presentation rules for a human actor. | [ui_surface, human_action_contract] | [DocumentResource] |

<!-- field: language_model | form: "tree" -->
## Language model

- corpus_as_resource:
  A documentation corpus is a configured resource, not a Go route table. The
  resource definition supplies authority. Runtime input supplies only declared
  selectors inside that authority.
- operations_as_words:
  Document list, get, search, validate, and suggest behavior should be
  visible as resource operations or MachineSpec words. A host may transport
  requests, but it should not choose a hidden document workflow.
- ui_as_actor_config:
  The Knowledge Manager UI is human actor configuration. It shapes what the
  human sees and which actions are available, just as a prompt shapes a model
  actor. Generic server code only serves configured assets and dispatches
  declared actions.
- response_modes:
  Response modes make output shape explicit. YAML and Markdown can be returned
  raw, parsed, or both when configuration allows it.

<!-- field: relationships | form: "keyed table" -->
## Relationships

| id | meaning | invariant |
| --- | --- | --- |
| DocumentResource_to_CorpusRoot | A resource names the trusted filesystem root for every operation. | Runtime selectors cannot replace the configured root. |
| DocumentResource_to_PathPolicy | Include patterns, extension allowlists, symlink policy, size limits, and binary policy constrain every read. | A selected path must pass path policy before any file content is read. |
| DocumentSelector_to_PathPolicy | Runtime request data is normalized into one relative path or query inside the configured policy. | Empty paths, absolute paths, parent segments, encoded traversal, and root escapes fail before dispatch output is produced. |
| DocumentResourceWord_to_DocumentResponseMode | Each operation returns only the configured output shape. | Parser selection is trusted config, not browser or model input. |
| DocumentResourceWord_to_MachineSpec | A machine_request run may dispatch document resource words while handling one HTTP request. | The machine sentence owns request-specific sequencing. |
| UIActorConfig_to_DocumentResource | The UI references named document resources and actions for sidebar grouping, labels, and display rules. | UI config can select known resources but cannot grant new filesystem authority. |

<!-- field: operations | form: "keyed table" -->
## Operations

| id | input | output | signals |
| --- | --- | --- | --- |
| list | Newline-delimited resource-relative paths from a preceding declared filesystem CLI dispatch. | path, name, category, extension, size, and configured metadata. | [DocumentListReady, DocumentResourceDenied, CommandError] |
| get | Resource-relative path from a declared selector. | path, raw content, parsed content, content type, extension, size, and parse diagnostics according to response mode. | [DocumentReady, DocumentMissing, DocumentResourceDenied, DocumentParseFailed, CommandError] |
| search | Query text and optional configured category or prefix filter. | query, matches, count, and redacted excerpts when configured. | [DocumentSearchReady, DocumentResourceDenied, CommandError] |
| validate | Resource-relative paths or configured validation scope. | status, findings, checked paths, and parser diagnostics. | [DocumentValidationPassed, DocumentValidationFailed, DocumentResourceDenied, CommandError] |
| suggest | Resource-relative path and instruction payload from a trusted machine or UI action. | suggestion ID, path, suggestions, and redacted patch preview. | [DocumentSuggestionReady, DocumentResourceDenied, CommandError] |

<!-- field: diagrams | form: "tree" -->
## Diagrams

- profile_configures_resource_and_ui:
  ```mermaid
  flowchart LR
    P[Profile] --> R[DocumentResource]
    P --> U[UIActorConfig]
    R --> ROOT[CorpusRoot docs/]
    R --> POLICY[PathPolicy]
    R --> MODE[Response modes]
    U --> ACTIONS[UI route and action config]
    ACTIONS --> R
  ```
- request_machine_reads_document:
  ```mermaid
  sequenceDiagram
    participant H as HTTP request
    participant M as MachineSpec
    participant W as document resource word
    participant R as DocumentResource
    participant P as PathPolicy
    H->>M: RequestEvent
    M->>W: dispatch get_document
    W->>R: resolve configured resource
    R->>P: validate selector
    P-->>R: accepted relative path
    R-->>W: raw or parsed response output
    W-->>M: DocumentReady signal
  ```

<!-- field: invariants -->
## Invariants

- Document resources are configured filesystem resources, not direct Go route policy.
- Corpus roots, include patterns, extension allowlists, parser modes, presentation category rules, and UI assets come from trusted config.
- Recursive discovery is a declared external filesystem CLI dispatch; generic Go validates and shapes candidates but does not walk the tree.
- Runtime request data can select only declared resource fields inside the configured root.
- Path traversal and root escape attempts fail before file reads.
- YAML and Markdown response modes are explicit per operation.
- UI actor config may live under agents/knowledge-manager/documentation-curator.
- Generic document resource code remains reusable without Knowledge Manager UI assets.
- MachineSpec owns request-handling order when document resources are used by machine_request.

<!-- field: non_goals -->
## Non goals

- Applying file edits or patch approvals.
- Defining a browser UI framework.
- Allowing runtime input to choose arbitrary roots, parser modes, or extension policies.
- Defining application route policy outside profile-owned MachineSpecs and REST YAML.

<!-- field: related_documents -->
## Related documents

- srd031-declarative-document-resources
- srd030-machine-request-rest-binding
- sm-machine-request-http
- sm-rest-tool-language
- cfg-rest-tool-format
