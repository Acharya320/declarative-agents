---
id: srd028-rest-client-tools
title: REST Client Tools
depends_on:
  - srd_id: srd001-core-types
    symbols_used: [Command, Result, Signal, Builder]
  - srd_id: srd005-agentic-loop
    symbols_used: [Loop, MachineSpec]
  - srd_id: srd022-builtin-tools
    symbols_used: [BuiltinRegistry, ToolDef]
  - srd_id: srd027-rest-tool-runtime
    symbols_used: [internal/tools/rest, net/http, OpenAPI import model, REST request state]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Agent profiles need outbound REST tools that are declared as vocabulary
words, not coded as one adapter per service. That matters. When an endpoint
behaves like stored state, the configured program should choose a resource
word so the operation's effect and rollback story stay visible in ToolDef.
Invoke is reserved for RPC-shaped endpoints that do not fit that resource
model.

REST client definitions must keep transport authority in trusted
configuration. Tool config fixes the upstream and operation. Runtime input
supplies only declared request parameters, usually from explicit tool input
or the previous Result. The generic client code sends the HTTP request
through net/http; then it maps the response into structured output and a
signal that MachineSpec can route.

<!-- field: goals -->
## Goals

- G1: Define outbound REST client tools as first-class ToolDef vocabulary words.
- G2: Define resource-oriented get, set, create, and delete words.
- G3: Define rest_client_invoke for RPC-shaped REST operations.
- G4: Define asynchronous rest_client_send and rest_client_await behavior.
- G5: Define result signals for reads, writes, missing resources, async sends, responses, domain failures, and infrastructure failures.
- G6: Define config-bound authority, runtime parameters, response shaping, tracing, metrics, idempotency, and undo expectations.
- G7: Define previous-Result parameter threading so a client word can select declared params from the prior Result and carry declared inputs forward.

<!-- req: R1 -->
## R1 - Client tool scope

<!-- req: R1.1 -->
### R1.1

REST client tools SHALL be declared as type builtin ToolDefs with category boundary.

<!-- req: R1.2 -->
### R1.2

REST client tools SHALL call outbound HTTP operations through net/http.

<!-- req: R1.3 -->
### R1.3

REST client tools SHALL use configured REST operation definitions loaded from hand-authored YAML or OpenAPI imports.

<!-- req: R1.4 -->
### R1.4

REST client tools SHALL use shared runtime rules from srd027-rest-tool-runtime.

<!-- req: R1.5 -->
### R1.5

REST client tools SHALL keep workflow order in MachineSpec transitions, not inside service-specific Go adapters.

<!-- req: R1.6 -->
### R1.6

The generic net/http request execution, policy enforcement, response reduction, observability, and receipt capture SHALL remain one boundary implementation. Non-CIDR operations SHALL NOT silently switch to a weaker subprocess path merely because curl is available (GH-1385).

<!-- req: R2 -->
## R2 - Config-bound operation authority

<!-- req: R2.1 -->
### R2.1

Tool config SHALL fix the REST client or upstream identity.

<!-- req: R2.2 -->
### R2.2

Tool config SHALL fix the resource and operation when a tool is operation-specific.

<!-- req: R2.3 -->
### R2.3

Tool config SHALL define method, path template, allowed query parameters, allowed header parameters, request schema, response schema, success status mapping, and failure status mapping.

<!-- req: R2.4 -->
### R2.4

Runtime input SHALL supply only declared path, query, header, and body parameters.

<!-- req: R2.5 -->
### R2.5

Runtime input MAY come from explicit tool parameters, previous Result output, or command state when declared by the tool.

<!-- req: R2.6 -->
### R2.6

Runtime input SHALL NOT supply arbitrary method, URL, authentication, redirect policy, or destination host.

<!-- req: R2.7 -->
### R2.7

Every declared path parameter SHALL be present and nonempty before URL rendering or transport; repeated placeholders SHALL use the same escaped value.

<!-- req: R3 -->
## R3 - Resource-oriented words

<!-- req: R3.1 -->
### R3.1

rest_client_get SHALL read one configured resource or collection without mutating the upstream service.

<!-- req: R3.2 -->
### R3.2

rest_client_set SHALL replace or update an existing configured resource when the REST definition declares update semantics.

<!-- req: R3.3 -->
### R3.3

rest_client_create SHALL create a configured resource when the REST definition declares create semantics.

<!-- req: R3.4 -->
### R3.4

rest_client_delete SHALL delete or deactivate a configured resource when the REST definition declares delete semantics.

<!-- req: R3.5 -->
### R3.5

Resource-oriented words SHOULD be preferred over rest_client_invoke when the endpoint maps to read, update, create, or delete behavior.

<!-- req: R3.6 -->
### R3.6

Resource-oriented words SHALL declare side effects and undo expectations based on the resource operation.

<!-- req: R4 -->
## R4 - rest_client_invoke

<!-- req: R4.1 -->
### R4.1

rest_client_invoke SHALL call a configured REST operation whose semantics are action-like or RPC-shaped.

<!-- req: R4.2 -->
### R4.2

rest_client_invoke SHALL be used for operations such as approve, reject, trigger, search, calculate, exchange, or provider-specific actions that are not simple resource reads or writes.

<!-- req: R4.3 -->
### R4.3

rest_client_invoke SHALL still resolve an operation name from trusted config, not from model-supplied transport data.

<!-- req: R4.4 -->
### R4.4

rest_client_invoke SHALL declare the domain signal mapping for success, accepted async work, domain failure, and infrastructure failure.

<!-- req: R4.5 -->
### R4.5

rest_client_invoke SHALL declare if the action is reversible, compensatable, or irreversible.

<!-- req: R5 -->
## R5 - Asynchronous client words

<!-- req: R5.1 -->
### R5.1

rest_client_send SHALL complete the configured HTTP submission transaction before returning, while leaving the remote operation's later domain completion to rest_client_await when the operation is declared asynchronous.

<!-- req: R5.2 -->
### R5.2

rest_client_send SHALL emit RESTAccepted or CommandError.

<!-- req: R5.3 -->
### R5.3

rest_client_send output SHALL include request ID, operation ID, accepted HTTP status, idempotency token when configured, submitted payload summary, and correlation metadata.

<!-- req: R5.4 -->
### R5.4

rest_client_send SHALL store async operation state in command state or shared request state indexed by request ID only after the submission status matches the configured success mapping.

<!-- req: R5.5 -->
### R5.5

rest_client_await SHALL wait for completion of a previously sent REST operation by request ID, correlation token, or configured operation state.

<!-- req: R5.6 -->
### R5.6

rest_client_await SHALL emit RESTResponded, RESTDomainFailed, RESTMissing, RESTAwaitTimedOut, or CommandError.

<!-- req: R5.7 -->
### R5.7

rest_client_await SHALL consume or update async operation state according to configured retention policy.

<!-- req: R5.8 -->
### R5.8

One REST client command dispatch SHALL perform at most one outbound HTTP request, except retries explicitly declared by the selected transport retry policy. Readiness polling SHALL use separate probe and delay states in MachineSpec.

<!-- req: R6 -->
## R6 - Result signals

<!-- req: R6.1 -->
### R6.1

RESTResourceRead SHALL signal a successful read of a configured resource.

<!-- req: R6.2 -->
### R6.2

RESTResourceWritten SHALL signal a successful create, update, set, delete, or deactivation operation.

<!-- req: R6.3 -->
### R6.3

RESTMissing SHALL signal a configured not-found or already-absent response.

<!-- req: R6.4 -->
### R6.4

RESTAccepted SHALL signal that an asynchronous request was accepted and can be awaited later.

<!-- req: R6.5 -->
### R6.5

RESTResponded SHALL signal that a configured non-resource or awaited operation returned a successful response.

<!-- req: R6.6 -->
### R6.6

RESTDomainFailed SHALL signal a service-level failure that is part of the operation contract.

<!-- req: R6.7 -->
### R6.7

CommandError SHALL signal infrastructure failure, config failure, schema failure, auth failure, network I/O failure, timeout, redaction failure, or response mapping failure.

<!-- req: R6.8 -->
### R6.8

Registration SHALL reject a bound operation whose configured success, failure, accepted, or await-timeout signal is absent from the ToolDef emits contract; CommandError remains the runtime infrastructure signal.

<!-- req: R7 -->
## R7 - Response shaping and output

<!-- req: R7.1 -->
### R7.1

REST client tools SHALL map HTTP response status, headers, and body into structured Result output according to configured response mappings.

<!-- req: R7.2 -->
### R7.2

REST client tools SHALL preserve raw status code and operation identity in output or diagnostics.

<!-- req: R7.3 -->
### R7.3

REST client tools SHALL validate response body against configured response schema when one is declared.

<!-- req: R7.4 -->
### R7.4

REST client tools SHALL support configured extraction of resource ID, version, ETag, request ID, domain error code, and retry-after metadata.

<!-- req: R7.5 -->
### R7.5

REST client tools SHALL redact configured secret fields before producing Result output.

<!-- req: R7.6 -->
### R7.6

REST client tools SHOULD expose output schemas in ToolDef output fields so downstream machine actions can consume structured results.

<!-- req: R8 -->
## R8 - Idempotency and retries

<!-- req: R8.1 -->
### R8.1

REST client tools SHOULD support configured idempotency tokens for create, set, delete, invoke, and send operations.

<!-- req: R8.2 -->
### R8.2

Idempotency tokens MAY come from config, runtime input, previous Result output, or generated request state when declared.

<!-- req: R8.3 -->
### R8.3

Retry policy SHALL be configured per operation and SHALL respect srd027 timeout, redirect, and size-limit rules.

<!-- req: R8.4 -->
### R8.4

Retries SHALL not be enabled for non-idempotent operations unless config declares an idempotency token or safe retry policy.

<!-- req: R8.5 -->
### R8.5

Tool output and traces SHOULD include retry count and idempotency metadata when available.

<!-- req: R8.6 -->
### R8.6

A retry policy's declared backoff and max_delay SHALL be honored between attempts—none waits zero, fixed uses initial_delay, and exponential doubles initial_delay with overflow-safe saturation at max_delay—and the delay SHALL be cancellation-aware so a cancelled dispatch stops waiting. An unsupported backoff, non-positive attempt count, missing retry_ref target, or invalid/negative delay SHALL be rejected at load.

<!-- req: R8.7 -->
### R8.7

The dispatch context SHALL govern both the in-flight HTTP request and retry waits. Cancellation of a non-idempotent request without an idempotency token SHALL report an indeterminate remote outcome.

<!-- req: R8.8 -->
### R8.8

Before dispatch, timeout-envelope inspection SHALL conservatively bound a selected retrying operation as attempts multiplied by the client request timeout plus the delay before each retry. Fixed and exponential delay terms SHALL match R8.6, checked duration arithmetic SHALL reject overflow instead of saturating, deterministic evidence SHALL name the client, retry_ref, attempts, backoff, aggregate duration, and enclosing command timeout, and the aggregate SHALL remain strictly below that command timeout. rest_client_await SHALL retain its separate async timeout and SHALL NOT inherit the prior submission's retry aggregate.

<!-- req: R9 -->
## R9 - Metrics, tracing, and redaction

<!-- req: R9.1 -->
### R9.1

REST client tools SHOULD trace operation ID, resource name, method, configured host, status code, elapsed time, retry count, request byte count, response byte count, and resulting signal. Outbound requests carry W3C trace context (traceparent, and tracestate when present) injected uniformly per srd016 R4 as a transport concern with no per-operation configuration.

<!-- req: R9.2 -->
### R9.2

ToolMetrics SHOULD include status code, retry count, response size, operation ID, and domain outcome when available.

<!-- req: R9.3 -->
### R9.3

REST client tools SHALL redact configured headers, query fields, body fields, credentials, and tokens from traces, metrics, errors, and Result output.

<!-- req: R9.4 -->
### R9.4

Verbose request or response body capture SHALL require explicit configuration and redaction rules.

<!-- req: R9.5 -->
### R9.5

At telemetry capture full, REST client dispatch spans SHALL carry the rendered request body under http.request.body and the response body under http.response.body. At capture off and capture delta they SHALL carry neither.

<!-- req: R9.6 -->
### R9.6

Captured bodies SHALL pass the operation's declared redaction selectors before reaching the span, so R9.3 holds for captured content and the redaction rules R9.4 requires are the operation's own. A body that does not parse as JSON SHALL NOT be captured, because field-wise redaction cannot apply to it.

<!-- req: R9.7 -->
### R9.7

Each captured body SHALL be capped at 16384 bytes, and a capped body SHALL set http.request.body.truncated or http.response.body.truncated to true. These four keys are repository-owned attributes, not OpenTelemetry semantic-convention keys, per srd015 R6.7.

<!-- req: R10 -->
## R10 - Undo and compensation

<!-- req: R10.1 -->
### R10.1

rest_client_get undo SHALL be noop because read operations do not mutate upstream state.

<!-- req: R10.2 -->
### R10.2

rest_client_set undo SHOULD use a configured previous-state restore or compensating action when the prior state is available.

<!-- req: R10.3 -->
### R10.3

rest_client_create undo SHOULD use a configured delete or deactivate compensation when the created resource ID is available.

<!-- req: R10.4 -->
### R10.4

rest_client_delete undo SHOULD use a configured restore or recreate compensation when the prior state is available.

<!-- req: R10.5 -->
### R10.5

rest_client_invoke undo SHALL be noop unless the operation declares a compensating action.

<!-- req: R10.6 -->
### R10.6

rest_client_send SHALL return a rollback receipt containing async request metadata and compensation state when the submitted operation can be cancelled or compensated; rest_client_await SHALL NOT duplicate that remote-effect receipt.

<!-- req: R10.7 -->
### R10.7

rest_client_await undo SHALL be noop unless await mutates command state that can be restored.

<!-- req: R10.8 -->
### R10.8

Irreversible REST client operations SHOULD require confirmation before affecting user data, external services, or shared infrastructure.

<!-- req: R11 -->
## R11 - Errors and recovery

<!-- req: R11.1 -->
### R11.1

REST client errors SHALL distinguish config_decode, operation_lookup, schema_invalid, request_rendering, auth_resolution, network_policy, network_io, timeout, status_mapping, response_mapping, async_state_missing, and compensation_missing stages.

<!-- req: R11.2 -->
### R11.2

Error output SHALL include operation ID, resource name when known, failure stage, status code when present, redacted diagnostic text, and resulting signal.

<!-- req: R11.3 -->
### R11.3

Domain failure status codes SHALL map to RESTDomainFailed or RESTMissing when declared by the operation.

<!-- req: R11.4 -->
### R11.4

Infrastructure failures SHALL map to CommandError and SHALL not be confused with service-domain failures.

<!-- req: R11.5 -->
### R11.5

Cancellation after a non-idempotent request may have reached the remote service SHALL report outcome indeterminate rather than claiming the mutation failed.

<!-- req: R11.6 -->
### R11.6

A configured scheme, host, CIDR, or port policy refusal SHALL report failure_stage network_policy during initial request construction and redirect-time revalidation. Selector and credential-scope failures remain target_resolution, malformed request templates remain request_rendering, and DNS or dial failures remain network_io.

<!-- req: R12 -->
## R12 - Previous-Result parameter threading

<!-- req: R12.1 -->
### R12.1

When an operation declares body_source previous_result, REST client tools SHALL populate declared runtime params by selecting fields from the previous Result output through an input_mapping map of declared-param-name to a $.-style selector, and SHALL present only declared params to the request builder so the previous Result's fixed output shape does not violate the declared-only contract.

<!-- req: R12.2 -->
### R12.2

input_mapping selectors SHALL resolve against the previous Result output shape, for example $.mapped.embedding, $.body.id, or $.carried.document, and SHALL target only declared path, query, header, or body params.

<!-- req: R12.3 -->
### R12.3

REST client operations SHALL support input carry-forward, so an operation MAY copy selected declared runtime inputs into its Result output under a carried key that a later word can select, without introducing a shared command-state store.

<!-- req: R12.4 -->
### R12.4

REST client tools SHALL reject input_mapping targets that are not declared params, and any input_mapping or carry-forward entry that names method, URL, authentication, redirect policy, or destination host. Accepted body_source values are params, previous_result, none, and command_state; command_state activation is governed by R13, which rejects it only when no command-state store is configured.

<!-- req: R12.5 -->
### R12.5

REST client body rendering SHALL preserve the single-token raw passthrough, so a body value equal to one "{{ params.NAME }}" token renders the param's raw typed value; a flat numeric vector param placed inside a nested array element, for example ["{{ params.query_embeddings }}"], SHALL render as [[...]].

<!-- req: R13 -->
## R13 - Command-state parameter threading

<!-- req: R13.1 -->
### R13.1

When an operation declares body_source command_state, REST client tools SHALL populate declared runtime params through an input_mapping map whose selectors take the $from(<label>).<dotted-path> form, selecting fields from a labeled prior step's output in the srd038 command-state store, and SHALL present only declared params to the request builder, identically to R12.1.

<!-- req: R13.2 -->
### R13.2

$from(<label>) label resolution SHALL follow srd038 R2, so the label names a prior step by command name or label, duplicate labels resolve most-recent-wins, resolution is scoped to the current run, and absolute step indices are not part of the authoring surface.

<!-- req: R13.3 -->
### R13.3

The R12.4 rejection list SHALL apply identically to command_state threading, so undeclared input_mapping targets are rejected and no selector or carry-forward entry may name method, URL, authentication, redirect policy, or destination host.

<!-- req: R13.4 -->
### R13.4

body_source previous_result SHALL remain supported and unchanged; command_state widens the addressable source from the immediately previous Result to any labeled prior step without altering the previous_result path.

<!-- req: R13.5 -->
### R13.5

The R12.4 rejection of command_state SHALL be lifted only when a command-state store is configured; without a configured store, body_source command_state SHALL continue to be rejected with the existing clear message.

<!-- req: R14 -->
## R14 - Trusted command-state operation targets

<!-- req: R14.1 -->
### R14.1

A configured operation MAY declare base_url_source command_state and a base_url_selector in $from(<label>).<path> form to select an absolute HTTP or HTTPS base URL from a labeled prior command output; omission SHALL preserve the configured client base_url path unchanged.

<!-- req: R14.2 -->
### R14.2

Operation target selection SHALL remain distinct from runtime input_mapping, which SHALL continue to reject method, URL, authentication, redirect policy, and destination-host targets under R12.4 and R13.3.

<!-- req: R14.3 -->
### R14.3

A selected target SHALL be rejected before network I/O when the selector is missing, unresolved, empty, non-string, malformed, relative, carries user information, uses a non-HTTP scheme, or violates configured network or redirect policy.

<!-- req: R14.4 -->
### R14.4

Configured credentials SHALL NOT be sent to a selected authority different from the configured client authority unless the trusted operation explicitly sets allow_auth_on_selected_authority.

<!-- req: R14.5 -->
### R14.5

REST results and traces SHALL identify the selected scheme, host, and port without recording user information, query values, or credentials.

<!-- req: R14.6 -->
### R14.6

A configured operation MAY instead declare a base_url_host_selector in $from(<label>).<path> form that resolves to a bare host or IP, composed with the operation's declared base_url_scheme (http or https, default http) and a port into the request base URL. This addresses a target discovered per item, for example a pod address a prior discovery step listed, which is published as an address rather than a whole URL. The two selector forms SHALL be mutually exclusive. The resolved host SHALL be rejected when it carries a scheme, credential, port, path, query, fragment, whitespace, or control character, so composition cannot widen transport authority. The scheme SHALL come only from trusted operation config; an unsupported scheme SHALL be rejected at load, and the composed URL SHALL pass the same R14.3 and R14.4 validation as a whole selected URL.

<!-- req: R14.7 -->
### R14.7

The composed port SHALL come from either a declared base_url_port literal or a base_url_port_selector in $from(<label>).<path> form, never both, and either form requires base_url_host_selector. A port selector serves a fleet whose targets do not share one port, where the port is published per item alongside the address. A literal outside 1-65535 SHALL be rejected at load. A resolved port SHALL be accepted as a string or a whole JSON number and SHALL be rejected before network I/O when it is another type, fractional, non-numeric, or outside 1-65535. A resolved port SHALL NOT widen destination authority; the composed authority SHALL still satisfy the client's configured network ports allowlist, so configuration bounds the reachable ports and the selector only chooses among them (R2.6).

<!-- field: non_goals -->
## Non-goals

- Defining inbound REST server tools; see srd029-rest-server-tools.
- Auto-wiring every lifecycle profile to a REST compensation executor; callers opt in through the core rollback options.
- Allowing runtime input to choose arbitrary URLs, methods, authentication, or network destinations.
- Generating Go clients from OpenAPI documents.
- Replacing MachineSpec sequencing with hidden REST client workflows.
- Requiring a curl executable for REST client words without full srd027/srd028 parity and runtime-image provisioning.

<!-- field: tool_words | form: "typed table" -->
## Tool words

| name | input | output | signals | side_effects | errors | undo |
| --- | --- | --- | --- | --- | --- | --- |
| rest_client_get | declared path, query, header parameters, and optional previous Result source. | resource payload, status code, headers, resource ID, version metadata, operation ID. | [RESTResourceRead, RESTMissing, RESTDomainFailed, CommandError] | [external_api_read] | [operation_lookup, schema_invalid, auth_resolution, network_io, timeout, response_mapping] | noop. |
| rest_client_set | declared path, query, header, and body parameters for an existing resource. | written resource payload or summary, status code, resource ID, version metadata, operation ID. | [RESTResourceWritten, RESTMissing, RESTDomainFailed, CommandError] | [external_api_write] | [operation_lookup, schema_invalid, request_rendering, auth_resolution, network_io, timeout, status_mapping] | previous-state restore or compensating action when configured; otherwise noop with irreversible or compensatable contract. |
| rest_client_create | declared path, query, header, and body parameters for a new resource. | created resource payload or summary, status code, resource ID, location header, operation ID. | [RESTResourceWritten, RESTDomainFailed, CommandError] | [external_api_write] | [operation_lookup, schema_invalid, request_rendering, auth_resolution, network_io, timeout, status_mapping] | configured delete or deactivate compensation when resource ID is available. |
| rest_client_delete | declared path, query, header, and optional body parameters for a resource delete or deactivate operation. | delete summary, status code, resource ID, operation ID. | [RESTResourceWritten, RESTMissing, RESTDomainFailed, CommandError] | [external_api_write] | [operation_lookup, schema_invalid, auth_resolution, network_io, timeout, status_mapping] | configured restore or recreate compensation when prior state is available. |
| rest_client_invoke | declared operation parameters for an RPC-shaped REST operation. | mapped response payload, status code, operation ID, domain outcome. | [RESTResponded, RESTAccepted, RESTDomainFailed, RESTMissing, CommandError] | [external_api_read, external_api_write] | [operation_lookup, schema_invalid, request_rendering, auth_resolution, network_io, timeout, response_mapping] | noop unless a compensating action is configured. |
| rest_client_send | declared operation parameters for an async REST operation. | request ID, operation ID, idempotency token, submitted payload summary, correlation metadata. | [RESTAccepted, CommandError] | [external_api_call, command_state] | [operation_lookup, schema_invalid, request_rendering, auth_resolution, network_io, timeout] | cancel or compensate when configured; otherwise record request metadata for audit. |
| rest_client_await | request ID, correlation token, operation state selector, timeout, and response mapping. | final response payload, status code, request ID, operation ID, retry and timing metadata. | [RESTResponded, RESTDomainFailed, RESTMissing, RESTAwaitTimedOut, CommandError] | [command_state] | [async_state_missing, timeout, network_io, response_mapping, status_mapping] | noop unless await mutates command state that can be restored. |

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

REST client tools are documented as first-class boundary ToolDef words backed by declarative client, resource, and operation definitions, with one policy-preserving generic transport implementation.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R1.4
- R1.5
- R1.6
- R2.1
- R2.2

<!-- ac: AC2 -->
### AC2

Tool config fixes client, resource, operation, transport policy, and schemas while runtime input supplies only declared parameters and every path parameter is present, nonempty, and escaped before transport.

<!-- field: traces -->
#### Traces

- R2.3
- R2.4
- R2.5
- R2.6
- R2.7

<!-- ac: AC3 -->
### AC3

Resource-oriented get, set, create, and delete words are preferred for resource-shaped endpoints and declare resource-specific effects and undo expectations.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3
- R3.4
- R3.5
- R3.6
- R10.1
- R10.2
- R10.3
- R10.4

<!-- ac: AC4 -->
### AC4

rest_client_invoke is reserved for RPC-shaped endpoints and still resolves operation authority from trusted config.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
- R4.3
- R4.4
- R4.5
- R10.5

<!-- ac: AC5 -->
### AC5

Async send and await requirements define request state, request IDs, idempotency metadata, accepted signals, final response signals, timeout signals, and state retention; readiness polling is a visible MachineSpec probe/delay loop with at most one HTTP request per REST dispatch.

<!-- field: traces -->
#### Traces

- R5.1
- R5.2
- R5.3
- R5.4
- R5.5
- R5.6
- R5.7
- R5.8
- R10.6
- R10.7

<!-- ac: AC6 -->
### AC6

Result signal requirements distinguish resource reads, writes, missing resources, accepted async sends, successful REST responses, REST domain failures, and CommandError infrastructure failures, and registration proves the bound operation is a subset of ToolDef emits.

<!-- field: traces -->
#### Traces

- R6.1
- R6.2
- R6.3
- R6.4
- R6.5
- R6.6
- R6.7
- R6.8
- R11.3
- R11.4

<!-- ac: AC7 -->
### AC7

Response shaping maps status, headers, body, resource metadata, request IDs, domain error codes, retry metadata, and redacted fields into structured output.

<!-- field: traces -->
#### Traces

- R7.1
- R7.2
- R7.3
- R7.4
- R7.5
- R7.6

<!-- ac: AC8 -->
### AC8

Idempotency, retry, metrics, tracing, and redaction requirements are explicit for external API calls.

<!-- field: traces -->
#### Traces

- R8.1
- R8.2
- R8.3
- R8.4
- R8.5
- R8.6
- R8.7
- R8.8
- R9.1
- R9.2
- R9.3
- R9.4
- R9.5
- R9.6
- R9.7
- R10.8

<!-- ac: AC9 -->
### AC9

Error and recovery requirements distinguish domain failures from infrastructure failures and report operation identity, failure stage, status, diagnostics, and resulting signal.

<!-- field: traces -->
#### Traces

- R11.1
- R11.2
- R11.3
- R11.4
- R11.5
- R11.6

<!-- ac: AC10 -->
### AC10

Each planned client tool word lists inputs, outputs, signals, side effects, errors, and undo strategy expectations.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3
- R3.4
- R4.1
- R5.1
- R5.5
- R10.1
- R10.2
- R10.3
- R10.4
- R10.5
- R10.6
- R10.7

<!-- ac: AC11 -->
### AC11

A REST client operation with body_source previous_result and input_mapping builds its request from a prior REST client Result without a not-declared error, threads a vector from one word's output into the next word's request body as a nested array, carries selected declared inputs forward under a carried key for a downstream word, and rejects undeclared input_mapping targets, command_state body_source, and transport-authority overrides.

<!-- field: traces -->
#### Traces

- R12.1
- R12.2
- R12.3
- R12.4
- R12.5
- R2.4
- R2.5
- R2.6

<!-- ac: AC12 -->
### AC12

A REST client operation with body_source command_state selects declared params through input_mapping $from(<label>).<path> selectors resolved against the srd038 command-state store (most-recent-wins, run-scoped), applies the R12.4 rejection list identically, leaves previous_result threading unchanged, and is rejected when no command-state store is configured.

<!-- field: traces -->
#### Traces

- R13.1
- R13.2
- R13.3
- R13.4
- R13.5

<!-- ac: AC13 -->
### AC13

A trusted REST operation selects a loopback child base URL from a labeled prior output, performs one request through normal REST dispatch, reports the selected authority, rejects unsafe target values and credential-scope changes before network I/O, preserves static configured targets, and leaves runtime authority parameters forbidden.

<!-- field: traces -->
#### Traces

- R14.1
- R14.2
- R14.3
- R14.4
- R14.5
- R2.4
- R2.5
- R2.6

<!-- ac: AC14 -->
### AC14

A trusted REST operation composes its base URL from a base_url_host_selector resolving to a bare host or IP plus a declared scheme and port, reaches the composed authority through one dispatch, rejects a resolved value carrying a scheme, credential, port, path, or query before network I/O, rejects both selector forms on one operation, and rejects an unsupported scheme or out-of-range port at load.

<!-- field: traces -->
#### Traces

- R14.6
- R14.3
- R14.4
- R2.4
- R2.6

<!-- ac: AC15 -->
### AC15

A trusted REST operation resolves the composed authority's port per item through base_url_port_selector, accepting a string or whole JSON number and reaching each target's own port; a malformed, fractional, wrongly typed, or out-of-range port fails before network I/O, a port outside the client's configured ports allowlist is refused, and both port forms or a port selector without a host selector are rejected at load.

<!-- field: traces -->
#### Traces

- R14.7
- R14.6
- R2.6

<!-- field: related_documents -->
## Related documents

- docs/specs/config-formats/tool-declaration-format.md
- docs/specs/config-formats/tool-authoring-standard.md
- docs/specs/semantic-models/tool-language.md
- docs/guides/rest-tool.md

<!-- field: implementation_status | form: "tree" -->
## Implementation status

- implemented:
  1. internal/tools/rest client commands for synchronous resource words, RPC-shaped invoke, async send, and async await.
  1. tools/builtin/rest client tool declarations and release 03.0 fixtures.
  1. Resource operation definitions loaded from YAML and OpenAPI imports.
  1. Config-bound method, path, host, auth, redirect, size-limit, request schema, response mapping, redaction, and status-signal behavior.
  1. Async request state indexed by request_id with idempotency and correlation metadata recorded in accepted output.
  1. CIDR allowlists are enforced while scheme, host, and port validation remains intact.
  1. Response bodies are validated against configured response schemas before success signals, and domain_error_code is included in structured output when configured.
  1. Async await can select stored completion by correlation token. Readiness polling is declared as one-shot REST probes separated by lifecycle delay states; await_operation is rejected with migration guidance.
  1. Idempotency metadata is emitted on configured outbound async requests, and retry validation requires idempotency for mutating operations when configured.
  1. Mutating REST client commands capture boundary compensation undo mementos with request, resource, idempotency, and compensation metadata. Read-only operations remain noop.
  1. Persisted checkpoint rollback can execute boundary compensation when a caller supplies the core CompensationExecutor. internal/tools/rest provides a CompensationExecutor that decodes REST boundary mementos, resolves the configured compensation operation from trusted definitions, filters runtime parameters to declared inputs, and executes through the REST client path.
  1. TestRESTClient_CompensationExecutorRunsConfiguredOperation, TestRESTClient_CompensationExecutorReportsMissingOperation, TestRollbackCheckpointExecutesBoundaryCompensation, and TestRollbackCheckpointReportsMissingBoundaryCompensationExecutor cover REST compensation execution and bounded failure reporting.
  1. - "Previous-Result parameter threading (R12)": body_source previous_result consumption through request_binding.input_mapping selects declared params from the prior Result output; request_binding.carry_forward copies declared inputs into the Result output under a carried key; validateRequestBinding enforces the R12.4 rejections (undeclared targets, transport-authority overrides). Covered by TestRESTClient_ThreadsPreviousResultParams, TestRESTClient_ChainThreadsVectorIntoQueryBody, TestRESTClient_CarryForwardEchoesDeclaredInputs, and TestRESTClient_RejectsPreviousResultThreadingMisuse.
  1. "Composed command-state targets (R14.6): an operation may declare base_url_host_selector with base_url_scheme and base_url_port, and client_target.go composes scheme://host[:port] from the resolved bare host before the unchanged R14.3/R14.4 validation. validateSelectedBaseURLTarget keeps the two selector forms mutually exclusive and rejects an unsupported scheme or an out-of-range port at load. IPv6 literals are bracketed. Covered by TestRESTClient_HostSelectorComposesSelectedAuthority, TestRESTClient_HostSelectorRejectsUnsafeValuesBeforeIO, TestResolveOperationBaseURLComposesAuthority, and TestRESTClient_HostSelectorRejectionTable."
  1. "Per-item composed ports (R14.7): base_url_port_selector resolves the composed authority's port from command state for a fleet whose targets do not share one port, accepting a label string or a whole JSON number. resolvedBaseURLPort and selectedPortString range-check the value before composition, validateComposedPort keeps the literal and selector forms mutually exclusive at load, and the composed authority still passes the client's network ports allowlist, so configuration bounds the reachable ports. Covered by TestRESTClient_PortSelectorResolvesPerItemPort and TestRESTClient_PortSelectorRejectsUnusablePorts."
  1. - "Command-state parameter threading (R13)": body_source command_state consumption through request_binding.input_mapping $from(<label>).<path> selectors resolved against the injected core CommandStateView (clientCmd implements core.CommandStateAware). selectCommandStateParams looks up the labeled step's output, decodes it, and walks the dotted path; a miss returns a typed error naming the label. validateBodySource accepts command_state; validateRequestBinding enforces V32 (a $from selector only under command_state, a $. selector only under previous_result) plus the R12.4 declared-target and transport-authority rejections; the no-store rejection is enforced at runtime (R13.5). previous_result threading is unchanged. Covered by TestRESTClient_CommandStateThreadsNonAdjacentStep, TestRESTClient_CommandStateDuplicateLabelMostRecentWins, TestRESTClient_CommandStateMissNamesLabel, TestRESTClient_CommandStateNoViewConfiguredRejected, and TestRESTClient_CommandStateRejectionTable.
- partial:
  1. Lifecycle profile wiring must still choose when to inject a REST CompensationExecutor into checkpoint rollback. Profile injection stays outside the REST client tool contract.

<!-- field: diagrams | form: "tree" -->
## Diagrams

- requirements_trace:
  ```mermaid
  flowchart LR
    Problem["Problem: Agent profiles need outbound REST tools that are decl..."] --> Goals["Goals: REST Client Tools"]
    Goals --> R1["R1: Client tool scope"]
    Goals --> R2["R2: Config-bound operation authority"]
    Goals --> R3["R3: Resource-oriented words"]
    Goals --> R4["R4: rest_client_invoke"]
    Goals --> R5["R5: Asynchronous client words"]
    Goals --> R6["R6: Result signals"]
    Goals --> R7["R7: Response shaping and output"]
    Goals --> R8["R8: Idempotency and retries"]
    Goals --> R9["R9: Metrics, tracing, and redaction"]
    Goals --> R10["R10: Undo and compensation"]
    Goals --> R11["R11: Errors and recovery"]
    Goals --> R12["R12: Previous-Result parameter threading"]
    Goals --> R13["R13: Command-state parameter threading"]
    AC1["AC1: REST client tools are documented as first-class boundary..."]
    R1 --> AC1
    R2 --> AC1
    AC2["AC2: Tool config fixes client, resource, operation, transport..."]
    R2 --> AC2
    AC3["AC3: Resource-oriented get, set, create, and delete words are..."]
    R3 --> AC3
    R10 --> AC3
    AC4["AC4: rest_client_invoke is reserved for RPC-shaped endpoints a..."]
    R4 --> AC4
    R10 --> AC4
    AC5["AC5: Async send and await requirements define request state, r..."]
    R5 --> AC5
    R10 --> AC5
    AC6["AC6: Result signal requirements distinguish resource reads, wr..."]
    R6 --> AC6
    R11 --> AC6
    AC7["AC7: Response shaping maps status, headers, body, resource met..."]
    R7 --> AC7
    AC8["AC8: Idempotency, retry, metrics, tracing, and redaction requi..."]
    R8 --> AC8
    R9 --> AC8
    R10 --> AC8
    AC9["AC9: Error and recovery requirements distinguish domain failur..."]
    R11 --> AC9
    AC10["AC10: Each planned client tool word lists inputs, outputs, sign..."]
    R3 --> AC10
    R4 --> AC10
    R5 --> AC10
    R10 --> AC10
    AC11["AC11: A REST client operation with body_source previous_result..."]
    R12 --> AC11
    R2 --> AC11
    AC12["AC12: A REST client operation with body_source command_state..."]
    R13 --> AC12
  ```
