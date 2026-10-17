---
id: rel03.0-uc002-rest-client-async
title: Grammar sends and awaits an asynchronous REST client operation
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An agent profile needs to submit outbound REST work without hiding its
acceptance boundary. The profile uses rest_client_send for one submission
transaction, then reads that stored result once or authors remote readiness
as repeated one-shot REST probes separated by declared delay states.

<!-- field: actor -->
## Actor

Agent profile author configuring asynchronous outbound REST work

<!-- field: trigger -->
## Trigger

Some external work outlives one machine step. The agent submits it, keeps
moving, and waits only when the grammar reaches the await state.

<!-- field: flow -->
## Flow

- F1: The author defines a REST operation with async.request_id, async.timeout, async.state_retention, and optional idempotency_token.
- F2: The author declares a rest_client_send ToolDef for the operation and a separate rest_client_await ToolDef for completion.
- F3: MachineSpec dispatches rest_client_send from SubmitRemoteWork.
- F4: rest_client_send validates input, submits the configured request, records AsyncRequest state, and emits RESTAccepted or CommandError.
- F5: MachineSpec routes RESTAccepted to PrepareLocalContext and routes CommandError to retry or failure handling.
- F6: The grammar performs local work while the remote request is pending.
- F7: To consume the stored submission result, MachineSpec dispatches rest_client_await with the request ID or correlation token.
- F8: For remote readiness, MachineSpec dispatches one configured REST read probe, routes not-ready to a lifecycle delay, and retries through explicit transitions and budget.
- F9: MachineSpec routes completion, missing, timeout or exhaustion, and infrastructure signals explicitly.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [R5.3](../software-requirements/srd027-rest-tool-runtime.md#r53), [R5.5](../software-requirements/srd027-rest-tool-runtime.md#r55), [R5.6](../software-requirements/srd027-rest-tool-runtime.md#r56), [R5.7](../software-requirements/srd027-rest-tool-runtime.md#r57) | accepted send and await use request state and declared signals |
| T2 | [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [R5.1](../software-requirements/srd028-rest-client-tools.md#r51), [R5.2](../software-requirements/srd028-rest-client-tools.md#r52), [R5.3](../software-requirements/srd028-rest-client-tools.md#r53), [R5.4](../software-requirements/srd028-rest-client-tools.md#r54), [R5.5](../software-requirements/srd028-rest-client-tools.md#r55), [R5.6](../software-requirements/srd028-rest-client-tools.md#r56), [R5.7](../software-requirements/srd028-rest-client-tools.md#r57) | async client words define send, await, request ID, signals, and state retention |
| T3 | [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [R8.1](../software-requirements/srd028-rest-client-tools.md#r81), [R8.2](../software-requirements/srd028-rest-client-tools.md#r82), [R8.3](../software-requirements/srd028-rest-client-tools.md#r83), [R8.4](../software-requirements/srd028-rest-client-tools.md#r84), [R8.5](../software-requirements/srd028-rest-client-tools.md#r85) | idempotency and retry metadata are configured |
| T4 | [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [R10.6](../software-requirements/srd028-rest-client-tools.md#r106), [R10.7](../software-requirements/srd028-rest-client-tools.md#r107) | undo expectations for send and await are explicit |
| T5 | cfg-rest-tool-format async_client_config | async request_id, correlation, idempotency_token, timeout, and state_retention are documented; await_operation is rejected |
| T6 | sm-rest-tool-language AsyncRequest and rest_client_send_to_AsyncRequest | outbound pending state is command/domain state |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Send and await are separate MachineSpec actions with an intervening local-work state.

<!-- field: traces -->
#### Traces

- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC5](../software-requirements/srd027-rest-tool-runtime.md#ac5)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC5](../software-requirements/srd028-rest-client-tools.md#ac5)

<!-- sc: S2 -->
### S2

The send result contains request ID, operation ID, accepted status, idempotency metadata, correlation metadata, and any configured compensation receipt.

<!-- field: traces -->
#### Traces

- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC5](../software-requirements/srd028-rest-client-tools.md#ac5)

<!-- sc: S3 -->
### S3

Await routes RESTResponded, RESTDomainFailed, RESTMissing, RESTAwaitTimedOut, and CommandError explicitly.

<!-- field: traces -->
#### Traces

- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC6](../software-requirements/srd028-rest-client-tools.md#ac6)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC9](../software-requirements/srd028-rest-client-tools.md#ac9)

<!-- sc: S4 -->
### S4

Transport retry and idempotency policy are declared in REST config; readiness probe interval and exhaustion are declared in MachineSpec.

<!-- field: traces -->
#### Traces

- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC8](../software-requirements/srd028-rest-client-tools.md#ac8)

<!-- field: out_of_scope -->
## Out of scope

- Running parallel goroutine scheduling tests.
- Defining every polling protocol shape.
- Implementing cancellation against each external service.
- Creating the release 03.0 test suite.

<!-- field: semantic_model -->
## Semantic model

sm-rest-tool-language

<!-- field: config_format -->
## Config format

cfg-rest-tool-format
