---
id: rel15.0-uc001-declarative-ux
title: Compose an Application UI from Declared Panels
status: implemented
test_suite: test-rel15.0-declarative-ux
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: release -->
## Release

15.0

<!-- field: summary -->
## Summary

A maintainer builds an application UI from ui.yaml, kit panels, and the
application's own domain panels, and deploys an observer whose UI the REST
tool serves from an embedded bundle. Shared panels read only the
presentation contract through one client, so the same panels run in every
application and under any base URL. Chatbot Mesh is the in-repository proof;
cohere-demo and agentic-wiki-mesh follow as external consumers.

<!-- field: actor -->
## Actor

A maintainer building or migrating an application user interface.

<!-- field: trigger -->
## Trigger

An application adds a panel, adopts the UI kit, or deploys an observer; or a
contract owner changes a response field a kit panel reads.

<!-- field: preconditions -->
## Preconditions

- srd033, srd020, srd021, and chatbot-mesh srd008 define the contract responses.
- srd003 governs the application manifest, UI assets, and shared-token policy.
- agent-core srd029 governs REST server bindings, including static_assets.

<!-- field: endpoint_usage_review -->
## Endpoint usage review

```yaml
reviewed_at:
    declarative-agents: bd71b931
    cohere-demo: c67d082
    agentic-wiki-mesh: 629a799
urls:
    - url: /monitor-proxy/{agent}/monitor/state
      used_by: [chatbot-mesh chatbot, cohere-demo chatbot, agentic-wiki-mesh chatbot]
      class: contract
    - url: /monitor-proxy/{agent}/monitor/events/stream
      used_by: [chatbot-mesh chatbot, cohere-demo chatbot turnActivity, agentic-wiki-mesh chatbot turnActivity]
      class: contract
    - url: /monitor-proxy/{trace_backend}/query/traces/{trace_id}
      used_by: [chatbot-mesh chatbot, cohere-demo chatbot, agentic-wiki-mesh chatbot]
      class: contract
    - url: /monitor-proxy/{trace_backend}/query/traces
      used_by: [agentic-wiki-mesh chatbot]
      class: contract
    - url: /monitor/fleet and /monitor/state on the observer origin
      used_by: [chatbot-mesh observer, cohere-demo observer, agentic-wiki-mesh observer]
      class: contract
    - url: /monitor-proxy/{agent}/monitor/machines and /monitor/tools/declared
      used_by: [agentic-wiki-mesh chatbot machineViews]
      class: contract
    - url: /trace-proxy/{agent}/monitor/machines, /monitor/tools/declared, /query/traces
      used_by: [cohere-demo observer declaredViews and api]
      class: contract under a retired proxy name; migrates to /monitor-proxy per srd004 R2.3
    - url: /api/v1/chat
      used_by: [all three chatbots]
      class: domain
    - url: /api/v1/audit, /api/v1/export, /api/v1/harvest, /api/v1/search, /api/v1/summarize, /api/v1/translate
      used_by: [cohere-demo chatbot]
      class: domain
    - url: /api/v1/documents, /api/v1/documents/query
      used_by: [cohere-demo ProvenancePanel through documentsApi]
      class: domain; ProvenancePanel leaves the kit panel set per srd004 R4.4
    - url: /provisioning/api
      used_by: [chatbot-mesh chatbot ProvisioningPanel]
      class: domain
```

<!-- field: flow -->
## Flow

- F1: The maintainer lists the panels the application needs and classifies
  each as a kit panel or a domain panel. A panel is a kit panel only when
  every endpoint it reads is in the presentation contract.
  - traces:
    - srd004-declarative-ux R1
    - srd004-declarative-ux R4.4
- F2: The maintainer declares the panels, sidebar groups, monitored agents,
  trace backend, and branding in ui.yaml version 2 and validates it.
  - traces:
    - srd004-declarative-ux R7.1
    - srd004-declarative-ux R7.2
    - srd004-declarative-ux R7.3
- F3: The application entrypoint supplies a registry of kit panels and domain
  panels. The kit Vite plugin reads ui.yaml during the build, and the
  generic shell renders routes and sidebar from it with the existing URL
  scheme.
  - traces:
    - srd004-declarative-ux R5
    - srd004-declarative-ux R7.4
- F4: At run time every panel obtains the client from shell context. Reads of
  other agents go through /monitor-proxy on the serving origin; a proxy 404
  hides the panels that depend on that agent.
  - traces:
    - srd004-declarative-ux R2
    - srd004-declarative-ux R6
- F5: For an observer, the maintainer selects the embedded observer bundle in
  rest.yaml, declares its config, and adds no UI files. The REST tool
  serves the bundle and its ui-config.json.
  - traces:
    - srd004-declarative-ux R9
- F6: The UI imports design tokens from the kit package. In this repository
  the kit is a file dependency; external repositories pin the release
  tarball.
  - traces:
    - srd004-declarative-ux R8
- F7: When an owner SRD changes a field a kit panel reads, the fixture changes
  in the same change and kit contract tests confirm every client and panel
  still parses it.
  - traces:
    - srd004-declarative-ux R3

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | applications/ui-kit/ | Package owning tokens, client, fixtures, shell, hooks, and kit panels. |
| T2 | applications/<application>/agents/<actor>/ui/ui.yaml | Composition input of the generic shell. |
| T3 | applications/<application>/agents/<actor>/rest.yaml | Declares the /monitor-proxy route and the static_assets binding or embedded bundle. |
| T4 | agent-core/internal/tools/rest/ | Serves embedded bundles and their declared config. |
| T5 | applications/<application>/agents/application.yaml | Declares each UI asset, including embedded sources, and its shared-token policy. |
| T6 | srd004-declarative-ux AC1 | Closed contract and endpoint classification. |
| T7 | srd004-declarative-ux AC2 | One proxy convention and not-deployed handling. |
| T8 | srd004-declarative-ux AC3 | Fixture-based consumer contract tests. |
| T9 | srd004-declarative-ux AC4 | Panel manifest and kit isolation. |
| T10 | srd004-declarative-ux AC5 | Build-time composition from ui.yaml. |
| T11 | srd004-declarative-ux AC6 | Injectable client. |
| T12 | srd004-declarative-ux AC7 | ui.yaml version 2 and retirement of the cross-check. |
| T13 | srd004-declarative-ux AC8 | Kit ownership, tokens, and distribution. |
| T14 | srd004-declarative-ux AC9 | Embedded observer bundle. |
| T15 | test-rel15.0-declarative-ux | Structural and executable evidence for this flow. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Every kit panel reads only contract endpoints and owner-defined fields,
and every URL the current UIs fetch has a contract or domain
classification.

<!-- field: traces -->
#### Traces

- srd004-declarative-ux AC1
- srd004-declarative-ux AC4

<!-- sc: S2 -->
### S2

The Chatbot Mesh chatbot takes its routes and sidebar from ui.yaml alone,
and adding a panel takes one ui.yaml entry and one registry line.

<!-- field: traces -->
#### Traces

- srd004-declarative-ux AC5
- srd004-declarative-ux AC7

<!-- sc: S3 -->
### S3

The Chatbot Mesh observer serves the embedded bundle and carries no UI
files.

<!-- field: traces -->
#### Traces

- srd004-declarative-ux AC9

<!-- sc: S4 -->
### S4

No shell, client, token, or observability panel code is duplicated under
applications/, and all of it is imported from the kit.

<!-- field: traces -->
#### Traces

- srd004-declarative-ux AC6
- srd004-declarative-ux AC8

<!-- sc: S5 -->
### S5

Cross-agent reads use one proxy name, and fixtures gate contract changes.

<!-- field: traces -->
#### Traces

- srd004-declarative-ux AC2
- srd004-declarative-ux AC3

<!-- field: out_of_scope -->
## Out of scope

- Migrating cohere-demo and agentic-wiki-mesh; their repositories own those issues.
- Changing monitor, fleet, or trace response schemas.
- Desktop shell and runtime panel loading.
