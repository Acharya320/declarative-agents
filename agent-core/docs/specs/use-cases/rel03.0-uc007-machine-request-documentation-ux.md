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

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd030-machine-request-rest-binding R1.1, R1.2, R1.3, R1.4, R1.5, R1.6, R2.1, R2.2, R2.3, R2.4, R2.5, R2.6, R2.7, R2.8 | one accepted HTTP request runs one configured MachineSpec sentence |
| T2 | srd030-machine-request-rest-binding R3.1, R3.2, R3.3, R3.4, R3.5, R3.6, R3.7, R3.8, R3.9, R3.10, R3.11, R4.1, R4.2, R4.3, R4.6, R4.7, R4.8, R4.9, R4.4, R4.5, R6.1, R6.2, R6.3, R6.4, R6.5, R6.6 | request input, terminal response extraction, and trace recording |
| T3 | srd031-declarative-document-resources R1.1, R1.2, R1.3, R1.4, R1.5, R2.1, R2.2, R2.3, R2.4, R2.5 | document resources connect request machines to docs/ |
| T4 | srd031-declarative-document-resources R3.1, R3.2, R3.3, R3.4, R3.5, R5.1, R5.2, R5.3, R5.4, R5.5, R5.6 | path safety and raw or parsed YAML response modes |
| T5 | srd031-declarative-document-resources R6.1, R6.2, R6.3, R6.4, R6.5 | Knowledge Manager UI assets and action placement are human actor configuration |
| T6 | rel03.0-uc006-documentation-curator | direct-host smoke behavior retained beside the machine_request successor proof |
| T7 | applications/catalog profile integration suite | concrete browser orchestration for the profile-owned tracer bullet |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The browser loads the Knowledge Manager UX from profile-owned UI assets without using bench routes.

<!-- field: traces -->
#### Traces

- srd031-declarative-document-resources AC6

<!-- sc: S2 -->
### S2

The document index request is fulfilled by one generic REST machine_request run and returns a sidebar-ready document list.

<!-- field: traces -->
#### Traces

- srd030-machine-request-rest-binding AC1
- srd031-declarative-document-resources AC1

<!-- sc: S3 -->
### S3

The raw YAML document request is fulfilled by one generic REST machine_request run and returns raw YAML plus configured parsed content.

<!-- field: traces -->
#### Traces

- srd030-machine-request-rest-binding AC5
- srd031-declarative-document-resources AC5

<!-- sc: S4 -->
### S4

Path traversal and disallowed extension attempts are rejected before a document resource reads a file.

<!-- field: traces -->
#### Traces

- srd031-declarative-document-resources AC4

<!-- sc: S5 -->
### S5

Puppeteer verifies rendered HTML, raw YAML visibility, and request trace evidence for both browser-driven machine_request runs.

<!-- field: traces -->
#### Traces

- srd030-machine-request-rest-binding AC7
- srd031-declarative-document-resources AC6

<!-- sc: S6 -->
### S6

Failures report the page URL, screenshot path, browser console output, network log, and machine_request trace metadata.

<!-- field: traces -->
#### Traces

- srd030-machine-request-rest-binding AC6

<!-- sc: S7 -->
### S7

Browser validation, suggestion, approval, and rejection use profile-owned machine_request routes and inherit standard dispatch trace, timeout, and metrics behavior.

<!-- field: traces -->
#### Traces

- srd030-machine-request-rest-binding AC1
- srd030-machine-request-rest-binding AC7

<!-- field: out_of_scope -->
## Out of scope

- Applying an approved patch to source files; approval records review state only.
- Replacing the rel03.0-uc006 direct-host smoke scenario.
- Running Puppeteer against non-local public listener addresses.
- Letting browser input choose document roots, parser modes, or machine files.

<!-- field: successor_to -->
## Successor to

rel03.0-uc006-documentation-curator

<!-- field: proof_owner -->
## Proof owner

applications/catalog

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

test-rel03.0-rest-tools
