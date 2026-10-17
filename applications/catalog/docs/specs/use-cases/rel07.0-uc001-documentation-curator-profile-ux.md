---
id: rel07.0-uc001-documentation-curator-profile-ux
title: Documentation Curator Profile UX Tracer
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

The documentation-curator profile owns the Knowledge Manager browser and
document-resource proof. The profile-owned integration target loads the
concrete profile from applications/catalog, starts its documentation host and control
routes through agent-core, exercises document and curator-action machine_request
endpoints, and shuts the profile down through declared lifecycle vocabulary.

This tracer belongs in applications/catalog because the profile, UX config, request
machine, REST definitions, and UI assets live beside the profile. agent-core
proves generic REST and machine_request behavior; this use case proves the
concrete profile wiring across that runtime boundary.

<!-- field: actor -->
## Actor

Profile maintainer verifying the Knowledge Manager profile bundle.

<!-- field: trigger -->
## Trigger

A maintainer runs the profile-owned documentation-curator integration target:

  mage integration:documentationCurator

<!-- field: preconditions -->
## Preconditions

- agent-core is available as a sibling ../../agent-core checkout or through demo.yaml core_root.
- agents/knowledge-manager/documentation-curator/profile.yaml resolves profile-local assets.
- The test can bind local loopback ports for documentation, machine_request, and control routes.

<!-- field: flow -->
## Flow

- F1: The target builds or locates an agent-core binary.
- F2: The target launches agents/knowledge-manager/documentation-curator/profile.yaml with a fixture documentation root.
- F3: The browser/API host serves profile-owned UI assets and UX config.
- F4: The document index and detail routes dispatch through the profile-owned request machine.
- F5: The browser posts validation and patch-review inputs to profile-owned action routes; the request machine dispatches the configured REST words.
- F6: Document and action response traces record server, route, request machine, terminal signal, and succeeded status.
- F7: The control route emits ExitRequested and the profile reaches Done after releasing listeners.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | sm-agent-role-taxonomy documentation-curator | The role owns document serving, document resources, profile-local UI, and controlled shutdown. |
| T2 | sm-multi-agent-profile-system documentation-curator->control | The profile embeds control flow and waits for explicit exit events. |
| T3 | test-rel07.0-profile-integrations documentation_curator_profile_ux | The test suite records the target and expected trace evidence. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The integration loads the documentation-curator profile from applications/catalog.

<!-- field: traces -->
#### Traces

- sm-agent-role-taxonomy documentation-curator
- sm-multi-agent-profile-system profile_asset_rules

<!-- sc: S2 -->
### S2

Document index and detail requests return terminal machine_request responses.

<!-- field: traces -->
#### Traces

- sm-multi-agent-profile-system actor_boundaries documentation-curator

<!-- sc: S3 -->
### S3

Profile-owned UI/static assets and UX config are served from the profile directory.

<!-- field: traces -->
#### Traces

- sm-agent-role-taxonomy documentation-curator checks

<!-- sc: S4 -->
### S4

Lifecycle exit reaches Done and releases documentation/control listeners.

<!-- field: traces -->
#### Traces

- sm-multi-agent-profile-system dependency_directions documentation-curator->control

<!-- sc: S5 -->
### S5

Validation, suggestion, approval, and rejection browser calls use the configured action machine_request endpoints, not a generic action envelope.

<!-- field: traces -->
#### Traces

- [srd011-knowledge-manager](../software-requirements/srd011-knowledge-manager.md) [AC1](../software-requirements/srd011-knowledge-manager.md#ac1)

<!-- field: out_of_scope -->
## Out of scope

- Re-testing generic REST route matching in agent-core.
- Applying approved patch records to source files.
- Requiring agent-core to know the documentation-curator profile path.

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel07.0-profile-integrations](../test-suites/test-rel07.0-profile-integrations.md)
