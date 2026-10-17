---
id: test-rel23.0-provider-failure-signals
title: Test Suite for Release 23.0 Provider Failure Signals
release: "23.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- [rel23.0-uc001-provider-failure-signals](../use-cases/rel23.0-uc001-provider-failure-signals.md)

<!-- field: preconditions -->
## Preconditions

- srd058-provider-libraries R2.4 and R2.5 define the opt-in, the transport rule, what stays CommandError, and the declared-emits agreement.
- The shipped dialects already map statuses to the R1.4 taxonomy; no dialect file changes in this release.
- Fake HTTP servers answer chat with chosen statuses, and one accepts a connection and never answers, so the deadline path is exercised without a real provider.

<!-- field: test_cases -->
## Test cases

<!-- case: A mapped error status reaches the machine as its declared signal | status: implemented -->
### A mapped error status reaches the machine as its declared signal

With the opt-in, a fake provider answers 503, 429, and 401 in turn.

<!-- field: use_case -->
#### Use case

[rel23.0-uc001-provider-failure-signals](../use-cases/rel23.0-uc001-provider-failure-signals.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm/... -run 'TestInvokeLLMEmitsMappedProviderFailure'

<!-- field: expected | form: "tree" -->
#### Expected

- unavailable_signal: ProviderUnavailable
- throttled_signal: ProviderThrottled
- unauthorized_signal: ProviderUnauthorized

<!-- field: traces -->
#### Traces

- [rel23.0-uc001-provider-failure-signals](../use-cases/rel23.0-uc001-provider-failure-signals.md) [S1](../use-cases/rel23.0-uc001-provider-failure-signals.md#s1)
- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [AC8](../software-requirements/srd058-provider-libraries.md#ac8)
- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [R2.4](../software-requirements/srd058-provider-libraries.md#r24)

<!-- case: A transport failure with no status is ProviderUnavailable | status: implemented -->
### A transport failure with no status is ProviderUnavailable

With the opt-in, one case refuses the connection and one accepts it and never answers, so the call reaches its deadline. This is the failure that motivated the release; it carries no HTTP status.

<!-- field: use_case -->
#### Use case

[rel23.0-uc001-provider-failure-signals](../use-cases/rel23.0-uc001-provider-failure-signals.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm/... -run 'TestInvokeLLMEmitsUnavailableOnTransportFailure'

<!-- field: expected | form: "tree" -->
#### Expected

- refused_connection_signal: ProviderUnavailable
- deadline_exceeded_signal: ProviderUnavailable

<!-- field: traces -->
#### Traces

- [rel23.0-uc001-provider-failure-signals](../use-cases/rel23.0-uc001-provider-failure-signals.md) [S2](../use-cases/rel23.0-uc001-provider-failure-signals.md#s2)
- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [AC8](../software-requirements/srd058-provider-libraries.md#ac8)
- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [R2.4](../software-requirements/srd058-provider-libraries.md#r24)

<!-- case: Unmapped, pre-exchange, and decode faults stay CommandError | status: implemented -->
### Unmapped, pre-exchange, and decode faults stay CommandError

With the opt-in, a status the dialect does not map, a credential that cannot be resolved, and a 200 whose body does not decode.

<!-- field: use_case -->
#### Use case

[rel23.0-uc001-provider-failure-signals](../use-cases/rel23.0-uc001-provider-failure-signals.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm/... -run 'TestInvokeLLMKeepsCommandErrorForUnclassifiedFailures'

<!-- field: expected | form: "tree" -->
#### Expected

- unmapped_status_signal: CommandError
- credential_fault_signal: CommandError
- decode_fault_signal: CommandError

<!-- field: traces -->
#### Traces

- [rel23.0-uc001-provider-failure-signals](../use-cases/rel23.0-uc001-provider-failure-signals.md) [S3](../use-cases/rel23.0-uc001-provider-failure-signals.md#s3)
- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [AC8](../software-requirements/srd058-provider-libraries.md#ac8)
- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [R2.4](../software-requirements/srd058-provider-libraries.md#r24)

<!-- case: Without the opt-in every failure stays CommandError | status: implemented -->
### Without the opt-in every failure stays CommandError

The same 503, 401, and deadline cases run against a ToolDef that does not opt in.

<!-- field: use_case -->
#### Use case

[rel23.0-uc001-provider-failure-signals](../use-cases/rel23.0-uc001-provider-failure-signals.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/llm/... -run 'TestInvokeLLMDefaultsToCommandError'

<!-- field: expected | form: "tree" -->
#### Expected

- every_signal: CommandError

<!-- field: traces -->
#### Traces

- [rel23.0-uc001-provider-failure-signals](../use-cases/rel23.0-uc001-provider-failure-signals.md) [S4](../use-cases/rel23.0-uc001-provider-failure-signals.md#s4)
- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [AC8](../software-requirements/srd058-provider-libraries.md#ac8)
- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [R3.3](../software-requirements/srd058-provider-libraries.md#r33)

<!-- case: Shipped profiles that bind invoke_llm are unaffected | status: implemented -->
### Shipped profiles that bind invoke_llm are unaffected

Every shipped catalog profile validates unchanged, which is the evidence that the default is genuinely off rather than merely documented.

<!-- field: use_case -->
#### Use case

[rel23.0-uc001-provider-failure-signals](../use-cases/rel23.0-uc001-provider-failure-signals.md)

<!-- field: go_test -->
#### Go test

cd ../applications/catalog && mage bootSmoke

<!-- field: expected | form: "tree" -->
#### Expected

- profiles_boot: true

<!-- field: traces -->
#### Traces

- [rel23.0-uc001-provider-failure-signals](../use-cases/rel23.0-uc001-provider-failure-signals.md) [S4](../use-cases/rel23.0-uc001-provider-failure-signals.md#s4)
- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [R3.3](../software-requirements/srd058-provider-libraries.md#r33)

<!-- case: An opt-in that does not declare its signals fails at load | status: implemented -->
### An opt-in that does not declare its signals fails at load

A ToolDef opts in while declaring only LLMResponded and CommandError.

<!-- field: use_case -->
#### Use case

[rel23.0-uc001-provider-failure-signals](../use-cases/rel23.0-uc001-provider-failure-signals.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog/... -run 'TestProviderFailureOptInRequiresDeclaredEmits'

<!-- field: expected | form: "tree" -->
#### Expected

- load_fails: true
- message_names_word: true
- message_names_missing_signals: true

<!-- field: traces -->
#### Traces

- [rel23.0-uc001-provider-failure-signals](../use-cases/rel23.0-uc001-provider-failure-signals.md) [S5](../use-cases/rel23.0-uc001-provider-failure-signals.md#s5)
- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [AC8](../software-requirements/srd058-provider-libraries.md#ac8)
- [srd058-provider-libraries](../software-requirements/srd058-provider-libraries.md) [R2.5](../software-requirements/srd058-provider-libraries.md#r25)
