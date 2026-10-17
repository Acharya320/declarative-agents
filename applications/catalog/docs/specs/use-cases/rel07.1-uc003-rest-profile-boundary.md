---
id: rel07.1-uc003-rest-profile-boundary
title: REST Profile Boundary
status: specified
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

The REST profile proves that configured client and server boundary words run
as explicit MachineSpec steps. A maintainer runs the sample REST flow, which
reads a payment resource, submits an async create, awaits completion, launches
a webhook server, waits for an inbound event, and stops the server. The run
reaches a Succeeded or Failed terminal state.

The proof belongs in applications/catalog because the REST authority in
testdata/conformance/rest/rest.yaml, the OpenAPI import assets under testdata/conformance/rest/openapi,
and the word selection in testdata/conformance/rest/machine.yaml and testdata/conformance/rest/tools.yaml
are profile-owned. The ollama-* profile variant shows the same boundary words
bound to a different REST definition set. agent-core supplies the reusable HTTP
runtime, but it should not encode the REST definitions or OpenAPI imports as
core-owned assets.

<!-- field: actor -->
## Actor

Profile maintainer verifying REST client and server boundary behavior.

<!-- field: trigger -->
## Trigger

A maintainer runs the REST sample profile from applications/catalog:

  agent --profile "testdata/conformance/rest/profile.yaml"

The ollama variant uses the same contract:

  agent --profile "testdata/conformance/rest/ollama-profile.yaml"

<!-- field: preconditions -->
## Preconditions

- testdata/conformance/rest/profile.yaml selects testdata/conformance/rest/machine.yaml, testdata/conformance/rest/tools.yaml, testdata/conformance/rest/declarations.yaml, and testdata/conformance/rest/rest.yaml.
- testdata/conformance/rest/tools.yaml selects get_payment, create_payment, await_payment, launch_payment_webhooks, await_payment_webhook, and stop_payment_webhooks.
- testdata/conformance/rest/openapi/payments.yaml and testdata/conformance/rest/openapi/ollama.yaml supply OpenAPI import data without generating Go code.

<!-- field: flow -->
## Flow

- F1: The machine reads a payment resource via get_payment in the ReadingPayment state and emits RESTResourceRead, RESTMissing, or RESTDomainFailed.
- F2: The machine submits an async create via create_payment and awaits completion via await_payment, emitting RESTAccepted, RESTResponded, or RESTAwaitTimedOut.
- F3: The machine launches the webhook server via launch_payment_webhooks and enters explicit LaunchingWebhook, WaitingWebhook, and StoppingWebhook states.
- F4: The machine waits for the inbound PaymentWebhookReceived signal and stops the server via stop_payment_webhooks.
- F5: The machine reaches the Succeeded terminal state on ServerStopped after StoppingWebhook, or the Failed terminal state on any missing, domain-failure, timeout, or CommandError path.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd007-rest](../software-requirements/srd007-rest.md) [R1](../software-requirements/srd007-rest.md#r1---behavior-model) | Configured client and server boundary words, including server launch, await, and stop, are explicit MachineSpec states. |
| T2 | [srd007-rest](../software-requirements/srd007-rest.md) [R2](../software-requirements/srd007-rest.md#r2---inputs-tools-and-runtime-dependencies) | REST inputs, the REST client, async client, and server tool families, and profile-owned rest.yaml authority sit over the agent-core HTTP runtime. |
| T3 | [srd007-rest](../software-requirements/srd007-rest.md) [R3](../software-requirements/srd007-rest.md#r3---outputs-and-outcomes) | REST outcomes include response signals, async completion, inbound event signals, listener lifecycle, and Succeeded or Failed terminal outcomes. |
| T4 | sm-agent-role-taxonomy rest | The role entry names the REST actor boundary, inputs, selected tool families, outputs, and checks. |
| T5 | test-rel07.1-profile-boundaries rest_profile_boundary | The suite records the base and ollama-* profiles, REST definitions, OpenAPI assets, and expected terminal states. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

REST authority lives in testdata/conformance/rest/rest.yaml and the ollama variant's ollama-rest.yaml while MachineSpec owns sequencing between boundary words.

<!-- field: traces -->
#### Traces

- [srd007-rest](../software-requirements/srd007-rest.md) [R2.3](../software-requirements/srd007-rest.md#r23)
- sm-agent-role-taxonomy rest checks

<!-- sc: S2 -->
### S2

Server launch, await, and stop are explicit machine states with a released listener at terminal state.

<!-- field: traces -->
#### Traces

- [srd007-rest](../software-requirements/srd007-rest.md) [R1.1](../software-requirements/srd007-rest.md#r11)
- [srd007-rest](../software-requirements/srd007-rest.md) [R1.2](../software-requirements/srd007-rest.md#r12)
- [srd007-rest](../software-requirements/srd007-rest.md) [R3.1](../software-requirements/srd007-rest.md#r31)

<!-- sc: S3 -->
### S3

The sample flow reaches Succeeded on the configured happy path and Failed on any missing, domain-failure, timeout, or infrastructure path.

<!-- field: traces -->
#### Traces

- [srd007-rest](../software-requirements/srd007-rest.md) [R3.2](../software-requirements/srd007-rest.md#r32)

<!-- field: out_of_scope -->
## Out of scope

- Reimplementing HTTP client or server transport in applications/catalog.
- Defining Go package APIs for the agent-core HTTP runtime.
- Generating Go code from the OpenAPI import assets.

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel07.1-profile-boundaries](../test-suites/test-rel07.1-profile-boundaries.md)
