---
id: rel03.0-uc001-rest-client-sync
title: Grammar invokes synchronous REST client resource words
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An agent profile needs to read and update an external resource through REST.
The profile declares narrow REST client words that point to configured
resources and operations. MachineSpec routes the resulting signals. Nothing
else decides the workflow. Inside the boundary word, generic runtime code
builds the configured request, sends it with net/http, maps the response, and
returns Result.Signal with redacted diagnostics.

<!-- field: actor -->
## Actor

Agent profile author configuring outbound resource access

<!-- field: trigger -->
## Trigger

Before the grammar can choose its next action, it needs fresh state from a
remote service. A later state may update the same resource after local
checks finish.

<!-- field: flow -->
## Flow

- F1: The author defines a rest.clients entry with base_url, auth_ref, limits_ref, and one resource with get and set operations.
- F2: The author declares rest_client_get and rest_client_set ToolDefs that reference the client, resource, and operation by name.
- F3: MachineSpec dispatches rest_client_get from a ReadRemoteResource state.
- F4: rest_client_get validates declared path/query/header input, sends the configured GET request, and emits RESTResourceRead, RESTMissing, RESTDomainFailed, or CommandError.
- F5: MachineSpec routes RESTResourceRead to a DecideUpdate state and routes RESTMissing, RESTDomainFailed, and CommandError to explicit recovery states.
- F6: After a local decision, MachineSpec dispatches rest_client_set from UpdateRemoteResource.
- F7: rest_client_set validates the declared body input, sends the configured update request, and emits RESTResourceWritten, RESTDomainFailed, or CommandError.
- F8: MachineSpec routes RESTResourceWritten to ContinueWorkflow and routes failure signals to compensation or operator review.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [R2.1](../software-requirements/srd027-rest-tool-runtime.md#r21), [R2.2](../software-requirements/srd027-rest-tool-runtime.md#r22), [R2.3](../software-requirements/srd027-rest-tool-runtime.md#r23), [R2.4](../software-requirements/srd027-rest-tool-runtime.md#r24), [R2.5](../software-requirements/srd027-rest-tool-runtime.md#r25) | REST definitions are declarative and MachineSpec owns sequencing |
| T2 | [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [R4.1](../software-requirements/srd027-rest-tool-runtime.md#r41), [R4.4](../software-requirements/srd027-rest-tool-runtime.md#r44) | net/http executes REST work and core.Result drives transitions |
| T3 | [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [R2.1](../software-requirements/srd028-rest-client-tools.md#r21), [R2.2](../software-requirements/srd028-rest-client-tools.md#r22), [R2.3](../software-requirements/srd028-rest-client-tools.md#r23), [R2.4](../software-requirements/srd028-rest-client-tools.md#r24), [R2.6](../software-requirements/srd028-rest-client-tools.md#r26) | config fixes operation authority and runtime input supplies declared parameters |
| T4 | [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [R3.1](../software-requirements/srd028-rest-client-tools.md#r31), [R3.2](../software-requirements/srd028-rest-client-tools.md#r32), [R3.5](../software-requirements/srd028-rest-client-tools.md#r35), [R3.6](../software-requirements/srd028-rest-client-tools.md#r36) | resource words cover read and update behavior |
| T5 | [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [R6.1](../software-requirements/srd028-rest-client-tools.md#r61), [R6.2](../software-requirements/srd028-rest-client-tools.md#r62), [R6.3](../software-requirements/srd028-rest-client-tools.md#r63), [R6.6](../software-requirements/srd028-rest-client-tools.md#r66), [R6.7](../software-requirements/srd028-rest-client-tools.md#r67) | result signals distinguish reads, writes, missing resources, domain failure, and CommandError |
| T6 | cfg-rest-tool-format validation rules V4, V5, V7 | request bindings, resource verbs, and side effects are validated |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The REST client ToolDefs reference named REST definitions instead of embedding full HTTP requests.

<!-- field: traces -->
#### Traces

- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC1](../software-requirements/srd028-rest-client-tools.md#ac1)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC2](../software-requirements/srd028-rest-client-tools.md#ac2)

<!-- sc: S2 -->
### S2

MachineSpec visibly routes RESTResourceRead, RESTResourceWritten, RESTDomainFailed, RESTMissing, and CommandError signals.

<!-- field: traces -->
#### Traces

- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC2](../software-requirements/srd027-rest-tool-runtime.md#ac2)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC6](../software-requirements/srd028-rest-client-tools.md#ac6)

<!-- sc: S3 -->
### S3

Resource get and set operations declare side effects and undo expectations.

<!-- field: traces -->
#### Traces

- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC3](../software-requirements/srd028-rest-client-tools.md#ac3)

<!-- sc: S4 -->
### S4

The implementation uses net/http and configured operation definitions.

<!-- field: traces -->
#### Traces

- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC4](../software-requirements/srd027-rest-tool-runtime.md#ac4)

<!-- field: out_of_scope -->
## Out of scope

- Implementing internal/tools/rest.
- Defining provider-specific client code.
- Allowing model-selected URLs or methods.
- Creating the release 03.0 test suite.

<!-- field: semantic_model -->
## Semantic model

sm-rest-tool-language

<!-- field: config_format -->
## Config format

cfg-rest-tool-format
