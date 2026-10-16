---
id: srd031-declarative-document-resources
title: Declarative Document Resources
depends_on:
  - srd_id: srd027-rest-tool-runtime
    symbols_used: [REST configuration, request validation, redaction, response mapping]
  - srd_id: srd030-machine-request-rest-binding
    symbols_used: [machine_request, RequestEvent, response words]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Knowledge Manager needs configured REST and machine_request words that read
documentation files from `docs/`. Route wiring and UI workflow policy belong
to the external profile, not a documentation-specific runtime host.

Document access belongs behind a resource definition. Trusted profile config
should name corpus rules and UI placement, leaving Go code to provide safe
filesystem reading and response rendering without hardcoding Knowledge Manager
routes. That split keeps the resource reusable.

Release 03.0 keeps the generic REST document resource loader deferred. The
accepted runtime path uses profile-selected filesystem resource ToolDefs such
as doc_list_resource and doc_read_resource. Top-level rest.document_resources
and machine_request.document_resources remain typed target-format fields that
validation rejects before launch.

<!-- field: goals -->
## Goals

- G1: Define declarative filesystem document resources for REST and machine_request words.
- G2: Keep reusable file reading capability separate from Knowledge Manager UI configuration.
- G3: Define safe path resolution, include patterns, extension allowlists, and traversal denial.
- G4: Define raw YAML, parsed YAML, raw Markdown, and parsed Markdown response modes.
- G5: Define where Knowledge Manager UI assets and UI route/action configuration may live.

<!-- req: R1 -->
## R1 - Resource binding scope

<!-- req: R1.1 -->
### R1.1

REST configuration SHALL support declarative document resource definitions for filesystem-backed documentation corpora.

<!-- req: R1.2 -->
### R1.2

Document resources SHALL be usable by REST server endpoints, REST client words, and machine_request request machines when selected by trusted profile config.

<!-- req: R1.3 -->
### R1.3

A document resource SHALL define read-only corpus access unless a separate mutating resource contract declares write behavior.

<!-- req: R1.4 -->
### R1.4

Direct Go documentation handlers SHALL NOT be the composition mechanism for Agent Architecture document behavior.

<!-- req: R1.5 -->
### R1.5

Generic document resource code SHALL be reusable by profiles other than Knowledge Manager.

<!-- req: R2 -->
## R2 - Corpus authority

<!-- req: R2.1 -->
### R2.1

Trusted configuration SHALL name each resource root, such as docs, configs, or a profile-local content directory.

<!-- req: R2.2 -->
### R2.2

Configuration SHALL declare include patterns for listable and readable files.

<!-- req: R2.3 -->
### R2.3

Configuration SHALL declare extension allowlists before a file can be returned.

<!-- req: R2.4 -->
### R2.4

Runtime input MAY supply only declared path variables or request fields used to select a document inside the configured root.

<!-- req: R2.5 -->
### R2.5

Runtime input SHALL NOT choose arbitrary roots, absolute paths, glob rules, extension allowlists, or parser modes.

<!-- req: R3 -->
## R3 - Path safety

<!-- req: R3.1 -->
### R3.1

Document resource paths SHALL be slash-normalized before resolution.

<!-- req: R3.2 -->
### R3.2

Empty paths, absolute paths, parent-directory segments, path traversal encodings, and resolved paths outside the configured root SHALL be rejected.

<!-- req: R3.3 -->
### R3.3

Directory reads SHALL be rejected unless the selected operation is an explicit list operation.

<!-- req: R3.4 -->
### R3.4

Symlink handling SHALL be explicit in configuration and SHALL default to denying root escapes after symlink resolution.

<!-- req: R3.5 -->
### R3.5

File reads SHALL enforce configured maximum bytes and binary detection policy before returning body content.

<!-- req: R4 -->
## R4 - Resource operations

<!-- req: R4.1 -->
### R4.1

Document resources SHALL support list and get operations.

<!-- req: R4.2 -->
### R4.2

Document resources MAY support search, validate, and suggest operations only when those operations are declared as separate resource operations or machine words.

<!-- req: R4.3 -->
### R4.3

List output SHALL include path, name, category, extension, and optional metadata fields configured for the resource.

<!-- req: R4.4 -->
### R4.4

Get output SHALL include path, raw content when enabled, parsed content when enabled, content type, extension, and size.

<!-- req: R4.5 -->
### R4.5

Operation errors SHALL distinguish path_required, invalid_path, extension_denied, not_found, read_failed, parse_failed, and size_limit_exceeded.

<!-- req: R4.6 -->
### R4.6

Presentation-specific category rules SHALL be ordered trusted resource configuration; generic filesystem code SHALL use only deterministic root or first-directory fallback categories.

<!-- req: R4.7 -->
### R4.7

Recursive list traversal SHALL execute as a declared external filesystem CLI word; a separate resource word SHALL validate and shape discovered paths without recursively walking directories in Go.

<!-- req: R5 -->
## R5 - Response modes

<!-- req: R5.1 -->
### R5.1

Resource configuration SHALL declare one or more response modes for each operation.

<!-- req: R5.2 -->
### R5.2

raw_yaml mode SHALL return the YAML file text without parsing.

<!-- req: R5.3 -->
### R5.3

parsed_yaml mode SHALL parse YAML into structured Result output and MAY also include raw content when configured.

<!-- req: R5.4 -->
### R5.4

raw_markdown mode SHALL return Markdown text with content type metadata.

<!-- req: R5.5 -->
### R5.5

parsed_markdown mode MAY return a configured Markdown AST, front matter, headings, or rendered HTML only when the parser is explicitly configured.

<!-- req: R5.6 -->
### R5.6

Parse failures SHALL map to configured signals and SHALL NOT hide the raw file when raw fallback is enabled.

<!-- req: R6 -->
## R6 - UI actor configuration

<!-- req: R6.1 -->
### R6.1

Knowledge Manager UI assets MAY live under agents/knowledge-manager/documentation-curator or another profile-owned asset path.

<!-- req: R6.2 -->
### R6.2

UI route labels, sidebar grouping, action names, resource references, and response presentation rules MAY live in profile-local configuration.

<!-- req: R6.3 -->
### R6.3

The generic server SHALL serve configured UI assets and route/action config without embedding Knowledge Manager-specific workflow policy.

<!-- req: R6.4 -->
### R6.4

Human-facing UI configuration SHALL be treated as actor configuration, analogous to prompt configuration for model actors.

<!-- req: R6.5 -->
### R6.5

Generic document resource capability SHALL remain reusable without the Knowledge Manager UI.

<!-- field: non_goals -->
## Non-goals

- Defining mutating document writes or applying patches to files.
- Letting browser input choose arbitrary filesystem roots or parser behavior.
- Defining a general CMS or web framework.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

The SRD defines declarative document resources that connect REST or machine_request words to a filesystem documentation corpus.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R1.5

<!-- ac: AC2 -->
### AC2

The SRD states that direct Go documentation handlers are not the composition mechanism for Agent Architecture document behavior.

<!-- field: traces -->
#### Traces

- R1.4
- R6.3

<!-- ac: AC3 -->
### AC3

Corpus authority requirements cover roots, include patterns, extension allowlists, runtime input limits, and trusted parser mode selection.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3
- R2.4
- R2.5

<!-- ac: AC4 -->
### AC4

Path safety requirements cover normalization, traversal denial, directory denial, symlink policy, size limits, and binary policy.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3
- R3.4
- R3.5

<!-- ac: AC5 -->
### AC5

Response mode requirements cover raw YAML, parsed YAML, raw Markdown, parsed Markdown, parse failures, and raw fallback.

<!-- field: traces -->
#### Traces

- R5.1
- R5.2
- R5.3
- R5.4
- R5.5
- R5.6

<!-- ac: AC6 -->
### AC6

UI actor configuration requirements distinguish profile-owned Knowledge Manager assets and route/action config from reusable generic server and resource capability.

<!-- field: traces -->
#### Traces

- R6.1
- R6.2
- R6.3
- R6.4
- R6.5

<!-- ac: AC7 -->
### AC7

Resource operation requirements define list and get behavior, optional higher-level operations, output fields, and operation error classes.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
- R4.3
- R4.4
- R4.5
- R4.6
- R4.7

<!-- field: related_documents -->
## Related documents

- docs/specs/semantic-models/document-resource-language.md
- docs/specs/semantic-models/machine-request-http.md
- docs/specs/config-formats/rest-tool-format.md
- docs/VISION.md
- docs/ARCHITECTURE.md

<!-- field: implementation_status -->
## Implementation status

```yaml
implemented:
    - agents/knowledge-manager/documentation-curator/request-machine.yaml selects reusable filesystem list_resource and read_resource words for the documentation corpus.
    - Profile-declared find words expose recursive traversal; internal/tools/filesystem list_resource validates and shapes their paths while read_resource performs bounded reads.
    - agents/knowledge-manager/documentation-curator owns the Knowledge Manager UI assets, request machine, REST route configuration, and action placement.
    - Generic Go tests verify document resource config and machine_request trace metadata; applications/catalog owns the concrete browser proof for index, document detail, raw YAML, rendered HTML, and nested document path evidence.
reserved:
    - Top-level rest.document_resources and machine_request.document_resources are typed target-format fields. ValidateDefinition rejects both with explicit reserved-field errors. Current request machines use profile-selected filesystem resource ToolDefs until generic REST document resource loading is implemented.
planned_follow_up: []
```

<!-- field: diagrams -->
## Diagrams

````yaml
requirements_trace: |
    ```mermaid
    flowchart LR
      Problem["Problem: Knowledge Manager needs configured REST and machine_r..."] --> Goals["Goals: Declarative Document Resources"]
      Goals --> R1["R1: Resource binding scope"]
      Goals --> R2["R2: Corpus authority"]
      Goals --> R3["R3: Path safety"]
      Goals --> R4["R4: Resource operations"]
      Goals --> R5["R5: Response modes"]
      Goals --> R6["R6: UI actor configuration"]
      AC1["AC1: The SRD defines declarative document resources that conne..."]
      R1 --> AC1
      AC2["AC2: The SRD states that direct Go documentation handlers are..."]
      R1 --> AC2
      R6 --> AC2
      AC3["AC3: Corpus authority requirements cover roots, include patter..."]
      R2 --> AC3
      AC4["AC4: Path safety requirements cover normalization, traversal d..."]
      R3 --> AC4
      AC5["AC5: Response mode requirements cover raw YAML, parsed YAML, r..."]
      R5 --> AC5
      AC6["AC6: UI actor configuration requirements distinguish profile-o..."]
      R6 --> AC6
      AC7["AC7: Resource operation requirements define list and get behav..."]
      R4 --> AC7
    ```
````
