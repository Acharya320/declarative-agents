---
id: rel23.0-uc001-provider-failure-signals
title: A machine retries a busy model provider and gives up on a misconfigured one
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An agent's machine tells a provider that is briefly unavailable apart from
one it cannot use at all, because invoke_llm reports them as different
signals. The profile opts into the failure taxonomy its dialect already
declares, routes ProviderUnavailable and ProviderThrottled to a bounded
retry, and leaves ProviderUnauthorized and CommandError terminal. A profile
that does not opt in runs unchanged.

<!-- field: actor -->
## Actor

A profile author deciding how an agent should answer a failing provider

<!-- field: trigger -->
## Trigger

A model call fails. The provider is loading a model, is rate limiting, has
rejected the credential, or was never reachable, and the machine must decide
whether another attempt is worth anything.

<!-- field: flow -->
## Flow

- F1: The profile's invoke_llm ToolDef opts into the failure taxonomy and declares ProviderUnavailable, ProviderThrottled, and ProviderUnauthorized among its emitted signals.
- F2: The provider answers 503. invoke_llm maps the status through the dialect's failure map and emits ProviderUnavailable; the machine routes it to a retry.
- F3: The provider does not answer at all and the call reaches its deadline. The failure carries no status, and invoke_llm emits ProviderUnavailable, because a provider that never answered is unavailable whatever the transport said.
- F4: The provider answers 401. invoke_llm emits ProviderUnauthorized and the machine terminates, since a rejected credential does not improve on a second ask.
- F5: The provider answers an error status its dialect does not map. invoke_llm emits CommandError rather than guessing which failure it was.
- F6: A credential cannot be resolved, or a successful answer cannot be decoded. Both emit CommandError, because neither describes the provider's state.
- F7: A profile that never opts in sees exactly the signals it saw before, and its machine needs no edit.
- F8: A profile that opts in without declaring the signals fails at load, naming the word and what it did not declare.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [R1.4](../software-requirements/srd058-provider-libraries.md#r14) | the shared failure taxonomy the dialects already map to |
| T2 | [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [R2.4](../software-requirements/srd058-provider-libraries.md#r24) | the opt-in, the transport rule, and what stays CommandError |
| T3 | [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [R2.5](../software-requirements/srd058-provider-libraries.md#r25) | declared emits remain the authority on what a machine routes |
| T4 | [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [R3.3](../software-requirements/srd058-provider-libraries.md#r33) | an unopted machine runs unmodified |
| T5 | [srd009-ollama-adapter](../software-requirements/srd009-ollama-adapter.md) [R4.3](../software-requirements/srd009-ollama-adapter.md#r43) | what an unavailable backend surfaces as |
| T6 | [srd009-ollama-adapter](../software-requirements/srd009-ollama-adapter.md) [R6.3](../software-requirements/srd009-ollama-adapter.md#r63) | the adapter does not retry, so the caller needs to know whether to |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

With the opt-in, a mapped error status reaches the machine as the signal the dialect names for it.

<!-- field: traces -->
#### Traces

- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [AC8](../software-requirements/srd058-provider-libraries.md#ac8)

<!-- sc: S2 -->
### S2

With the opt-in, a transport failure carrying no status reaches the machine as ProviderUnavailable, including a call that exceeds its deadline.

<!-- field: traces -->
#### Traces

- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [AC8](../software-requirements/srd058-provider-libraries.md#ac8)

<!-- sc: S3 -->
### S3

An unmapped error status, a pre-exchange fault, and a decode fault each remain CommandError.

<!-- field: traces -->
#### Traces

- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [AC8](../software-requirements/srd058-provider-libraries.md#ac8)

<!-- sc: S4 -->
### S4

Without the opt-in every path emits what it emitted before, and the shipped profiles that bind invoke_llm run unmodified.

<!-- field: traces -->
#### Traces

- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [AC8](../software-requirements/srd058-provider-libraries.md#ac8)

<!-- sc: S5 -->
### S5

A ToolDef whose opt-in and declared emits disagree fails at load, naming the word and the missing signals.

<!-- field: traces -->
#### Traces

- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [AC8](../software-requirements/srd058-provider-libraries.md#ac8)

<!-- field: out_of_scope -->
## Out of scope

- Retrying inside the adapter or inside invoke_llm; a retry is a machine's decision (srd009 R6.3).
- Migrating the shipped profiles that bind invoke_llm; each adopts the opt-in when its own behavior calls for it.
- Classifying provider failures for embed and rerank, which already emit the taxonomy.
