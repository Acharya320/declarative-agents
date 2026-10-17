---
id: rel03.0-uc007-machine-request-documentation-ux
title: Knowledge Manager UX uses machine_request document resources
status: done
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A developer opens the Knowledge Manager documentation UX. Browser requests for
document reads and curator actions enter through machine_request endpoints.
Each accepted HTTP request runs one short-lived MachineSpec sentence that
dispatches configured resource or REST words and returns configured terminal
output. Puppeteer drives the browser proof and captures request trace evidence.

<!-- field: actor -->
## Actor

Developer browsing Agent Core documentation in a browser

<!-- field: trigger -->
## Trigger

The developer opens the Knowledge Manager UX after the profile has launched a
machine_request-backed documentation server with configured document resources.

<!-- field: flow -->
## Flow

- F1: The agent loads the external documentation-curator profile with machine_request endpoints and the documentation_corpus document resource selected.
- F2: The browser opens the Knowledge Manager documentation URL served from profile-owned UI assets.
- F3: The UI requests the document index through the generic REST documentation_curator.documents machine_request route. Request validation runs before one short-lived request machine starts.
- F4: The request machine dispatches declared find against the fixed docs/ corpus, then a resource word validates and shapes those paths before emitting DocumentListReady.
- F5: The response word maps DocumentListReady output to the HTTP index response. The UI shows sidebar categories from the returned document list.
- F6: The operator selects SPECIFICATIONS.md or another YAML document from the sidebar.
- F7: The browser requests that raw YAML document through the generic REST documentation_curator.document machine_request route. Path validation runs before a second short-lived request machine starts.
- F8: The request machine dispatches a document get word. Path policy denies traversal and permits only configured YAML or Markdown extensions.
- F9: The response word returns raw YAML and parsed content according to the document resource response mode.
- F10: The UI shows the document as HTML and exposes the raw YAML view.
- F11: The browser resolves profile-owned action routes and posts validation, suggestion, approval, and rejection input directly to configured machine_request endpoints.
- F12: Each action request dispatches its selected REST word through the request machine; the legacy generic action envelope and private direct command runner are absent.
- F13: Puppeteer verifies visible content, raw YAML retrieval, rendered HTML, action dispatch, nested document path evidence, captured request IDs, machine names, terminal signals, and RunResult summaries.
- F14: Browser, server, or machine failures produce a structured Puppeteer failure report with screenshot, console output, network log, and captured trace metadata.

<!-- field: successor_to -->
## Successor to

rel03.0-uc006-documentation-curator

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [R1.1](../software-requirements/srd030-machine-request-rest-binding.md#r11), [R1.2](../software-requirements/srd030-machine-request-rest-binding.md#r12), [R1.3](../software-requirements/srd030-machine-request-rest-binding.md#r13), [R1.4](../software-requirements/srd030-machine-request-rest-binding.md#r14), [R1.5](../software-requirements/srd030-machine-request-rest-binding.md#r15), [R1.6](../software-requirements/srd030-machine-request-rest-binding.md#r16), [R2.1](../software-requirements/srd030-machine-request-rest-binding.md#r21), [R2.2](../software-requirements/srd030-machine-request-rest-binding.md#r22), [R2.3](../software-requirements/srd030-machine-request-rest-binding.md#r23), [R2.4](../software-requirements/srd030-machine-request-rest-binding.md#r24), [R2.5](../software-requirements/srd030-machine-request-rest-binding.md#r25), [R2.6](../software-requirements/srd030-machine-request-rest-binding.md#r26), [R2.7](../software-requirements/srd030-machine-request-rest-binding.md#r27), [R2.8](../software-requirements/srd030-machine-request-rest-binding.md#r28) | one accepted HTTP request runs one configured MachineSpec sentence |
| T2 | [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [R3.1](../software-requirements/srd030-machine-request-rest-binding.md#r31), [R3.2](../software-requirements/srd030-machine-request-rest-binding.md#r32), [R3.3](../software-requirements/srd030-machine-request-rest-binding.md#r33), [R3.4](../software-requirements/srd030-machine-request-rest-binding.md#r34), [R3.5](../software-requirements/srd030-machine-request-rest-binding.md#r35), [R3.6](../software-requirements/srd030-machine-request-rest-binding.md#r36), [R3.7](../software-requirements/srd030-machine-request-rest-binding.md#r37), [R3.8](../software-requirements/srd030-machine-request-rest-binding.md#r38), [R3.9](../software-requirements/srd030-machine-request-rest-binding.md#r39), [R3.10](../software-requirements/srd030-machine-request-rest-binding.md#r310), [R3.11](../software-requirements/srd030-machine-request-rest-binding.md#r311), [R4.1](../software-requirements/srd030-machine-request-rest-binding.md#r41), [R4.2](../software-requirements/srd030-machine-request-rest-binding.md#r42), [R4.3](../software-requirements/srd030-machine-request-rest-binding.md#r43), [R4.6](../software-requirements/srd030-machine-request-rest-binding.md#r46), [R4.7](../software-requirements/srd030-machine-request-rest-binding.md#r47), [R4.8](../software-requirements/srd030-machine-request-rest-binding.md#r48), [R4.9](../software-requirements/srd030-machine-request-rest-binding.md#r49), [R4.4](../software-requirements/srd030-machine-request-rest-binding.md#r44), [R4.5](../software-requirements/srd030-machine-request-rest-binding.md#r45), [R6.1](../software-requirements/srd030-machine-request-rest-binding.md#r61), [R6.2](../software-requirements/srd030-machine-request-rest-binding.md#r62), [R6.3](../software-requirements/srd030-machine-request-rest-binding.md#r63), [R6.4](../software-requirements/srd030-machine-request-rest-binding.md#r64), [R6.5](../software-requirements/srd030-machine-request-rest-binding.md#r65), [R6.6](../software-requirements/srd030-machine-request-rest-binding.md#r66) | request input, terminal response extraction, and trace recording |
| T3 | [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md) [R1.1](../software-requirements/srd031-declarative-document-resources.md#r11), [R1.2](../software-requirements/srd031-declarative-document-resources.md#r12), [R1.3](../software-requirements/srd031-declarative-document-resources.md#r13), [R1.4](../software-requirements/srd031-declarative-document-resources.md#r14), [R1.5](../software-requirements/srd031-declarative-document-resources.md#r15), [R2.1](../software-requirements/srd031-declarative-document-resources.md#r21), [R2.2](../software-requirements/srd031-declarative-document-resources.md#r22), [R2.3](../software-requirements/srd031-declarative-document-resources.md#r23), [R2.4](../software-requirements/srd031-declarative-document-resources.md#r24), [R2.5](../software-requirements/srd031-declarative-document-resources.md#r25) | document resources connect request machines to docs/ |
| T4 | [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md) [R3.1](../software-requirements/srd031-declarative-document-resources.md#r31), [R3.2](../software-requirements/srd031-declarative-document-resources.md#r32), [R3.3](../software-requirements/srd031-declarative-document-resources.md#r33), [R3.4](../software-requirements/srd031-declarative-document-resources.md#r34), [R3.5](../software-requirements/srd031-declarative-document-resources.md#r35), [R5.1](../software-requirements/srd031-declarative-document-resources.md#r51), [R5.2](../software-requirements/srd031-declarative-document-resources.md#r52), [R5.3](../software-requirements/srd031-declarative-document-resources.md#r53), [R5.4](../software-requirements/srd031-declarative-document-resources.md#r54), [R5.5](../software-requirements/srd031-declarative-document-resources.md#r55), [R5.6](../software-requirements/srd031-declarative-document-resources.md#r56) | path safety and raw or parsed YAML response modes |
| T5 | [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md) [R6.1](../software-requirements/srd031-declarative-document-resources.md#r61), [R6.2](../software-requirements/srd031-declarative-document-resources.md#r62), [R6.3](../software-requirements/srd031-declarative-document-resources.md#r63), [R6.4](../software-requirements/srd031-declarative-document-resources.md#r64), [R6.5](../software-requirements/srd031-declarative-document-resources.md#r65) | Knowledge Manager UI assets and action placement are human actor configuration |
| T6 | [rel03.0-uc006-documentation-curator](rel03.0-uc006-documentation-curator.md) | direct-host smoke behavior retained beside the machine_request successor proof |
| T7 | applications/catalog profile integration suite | concrete browser orchestration for the profile-owned tracer bullet |

<!-- field: proof_owner -->
## Proof owner

applications/catalog

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The browser loads the Knowledge Manager UX from profile-owned UI assets without using bench routes.

<!-- field: traces -->
#### Traces

- [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md) [AC6](../software-requirements/srd031-declarative-document-resources.md#ac6)

<!-- sc: S2 -->
### S2

The document index request is fulfilled by one generic REST machine_request run and returns a sidebar-ready document list.

<!-- field: traces -->
#### Traces

- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC1](../software-requirements/srd030-machine-request-rest-binding.md#ac1)
- [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md) [AC1](../software-requirements/srd031-declarative-document-resources.md#ac1)

<!-- sc: S3 -->
### S3

The raw YAML document request is fulfilled by one generic REST machine_request run and returns raw YAML plus configured parsed content.

<!-- field: traces -->
#### Traces

- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC5](../software-requirements/srd030-machine-request-rest-binding.md#ac5)
- [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md) [AC5](../software-requirements/srd031-declarative-document-resources.md#ac5)

<!-- sc: S4 -->
### S4

Path traversal and disallowed extension attempts are rejected before a document resource reads a file.

<!-- field: traces -->
#### Traces

- [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md) [AC4](../software-requirements/srd031-declarative-document-resources.md#ac4)

<!-- sc: S5 -->
### S5

Puppeteer verifies rendered HTML, raw YAML visibility, and request trace evidence for both browser-driven machine_request runs.

<!-- field: traces -->
#### Traces

- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC7](../software-requirements/srd030-machine-request-rest-binding.md#ac7)
- [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md) [AC6](../software-requirements/srd031-declarative-document-resources.md#ac6)

<!-- sc: S6 -->
### S6

Failures report the page URL, screenshot path, browser console output, network log, and machine_request trace metadata.

<!-- field: traces -->
#### Traces

- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC6](../software-requirements/srd030-machine-request-rest-binding.md#ac6)

<!-- sc: S7 -->
### S7

Browser validation, suggestion, approval, and rejection use profile-owned machine_request routes and inherit standard dispatch trace, timeout, and metrics behavior.

<!-- field: traces -->
#### Traces

- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC1](../software-requirements/srd030-machine-request-rest-binding.md#ac1)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC7](../software-requirements/srd030-machine-request-rest-binding.md#ac7)

<!-- field: out_of_scope -->
## Out of scope

- Applying an approved patch to source files; approval records review state only.
- Replacing the rel03.0-uc006 direct-host smoke scenario.
- Running Puppeteer against non-local public listener addresses.
- Letting browser input choose document roots, parser modes, or machine files.

<!-- field: machine_config -->
## Machine config

agents/knowledge-manager/documentation-curator/machine.yaml

<!-- field: rest_definitions -->
## REST definitions

agents/knowledge-manager/documentation-curator/rest.yaml

<!-- field: document_resource -->
## Document resource

documentation_corpus

<!-- field: browser_harness -->
## Browser harness

Puppeteer

<!-- field: integration_target -->
## Integration target

applications/catalog documentation-curator browser target

<!-- field: entry_point -->
## Entry point

profile-owned browser integration entry point

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel03.0-rest-tools](../test-suites/test-rel03.0-rest-tools.md)
