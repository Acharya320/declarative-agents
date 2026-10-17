---
id: test-rel03.0-rest-tools
title: Test Suite for Release 03.0 REST Tools
release: "03.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- [rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md)
- [rel03.0-uc002-rest-client-async](../use-cases/rel03.0-uc002-rest-client-async.md)
- [rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md)
- [rel03.0-uc004-openapi-import](../use-cases/rel03.0-uc004-openapi-import.md)
- [rel03.0-uc005-ollama-openapi-model-listing](../use-cases/rel03.0-uc005-ollama-openapi-model-listing.md)
- [rel03.0-uc006-documentation-curator](../use-cases/rel03.0-uc006-documentation-curator.md)
- [rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md)
- [rel03.0-uc008-rest-await-any-fan-in](../use-cases/rel03.0-uc008-rest-await-any-fan-in.md)

<!-- field: preconditions -->
## Preconditions

- internal/tools/rest contains the REST runtime and executable Go coverage for this suite.
- Test fixtures can provide hand-authored REST config files and OpenAPI 3 documents.
- HTTP client tests can use httptest.Server with local-only addresses.
- HTTP server tests can bind listener port 0 and use per-test queues.
- Redaction tests can use synthetic credentials and payload fields.
- Live Ollama E2E coverage requires Ollama on 127.0.0.1:11434 and a Qwen 3 family model.
- Documentation host package tests use local fixture directories and do not require a browser.
- Concrete documentation-curator and Knowledge Manager browser tracer bullets belong to applications/catalog suites with the profile assets.
- Generic machine_request trace metadata is covered by internal/tools/rest Go tests without relying on profile-specific proxy enrichment.
- Generic REST route conformance runs under plain go test without browser prerequisites.

<!-- field: test_cases -->
## Test cases

<!-- case: REST definitions load into boundary ToolDefs | status: done -->
### REST definitions load into boundary ToolDefs

After loading a REST config file and its tool declaration file, inspect
each REST ToolDef. It should name one boundary word and reference
configured REST definitions without embedding arbitrary transport
authority fields.

<!-- field: use_case -->
#### Use case

[rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md)

<!-- field: go_test -->
#### Go test

TestRESTConfig_LoadsToolDefinitions

<!-- field: inputs | form: "tree" -->
#### Inputs

- config: rest_client_and_server_minimal
- declarations:
  1. rest_client_get
  1. rest_client_set
  1. rest_server_launch

<!-- field: expected | form: "tree" -->
#### Expected

- definitions_loaded: true
- tooldefs_category_boundary: true
- tooldefs_reference_named_rest_config: true
- arbitrary_url_from_tool_input_allowed: false
- package_layout: internal/tools/rest

<!-- field: traces -->
#### Traces

- [rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md) [S1](../use-cases/rel03.0-uc001-rest-client-sync.md#s1)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC1](../software-requirements/srd027-rest-tool-runtime.md#ac1)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC2](../software-requirements/srd027-rest-tool-runtime.md#ac2)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC9](../software-requirements/srd027-rest-tool-runtime.md#ac9)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC1](../software-requirements/srd028-rest-client-tools.md#ac1)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC2](../software-requirements/srd028-rest-client-tools.md#ac2)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC10](../software-requirements/srd028-rest-client-tools.md#ac10)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC1](../software-requirements/srd029-rest-server-tools.md#ac1)

<!-- case: Synchronous resource words build configured requests | status: done -->
### Synchronous resource words build configured requests

Against an httptest.Server, run rest_client_get and rest_client_set with
configured operations. net/http should receive the expected method and
path. Response mapping then produces structured output for MachineSpec.
The case checks the two resource paths together because their signals
prove that read, write, missing, domain-failure, and infrastructure
outcomes stay separate at the grammar boundary.
This matters when a remote service returns a valid HTTP response that
still represents a business failure, because the machine must route that
case through RESTDomainFailed while reserving CommandError for transport,
schema, auth, timeout, and mapping failures.

<!-- field: use_case -->
#### Use case

[rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md)

<!-- field: go_test -->
#### Go test

TestRESTClient_SyncResourceWords

<!-- field: inputs | form: "tree" -->
#### Inputs

- client: github_like_issue_client
- operations:
  1. get_issue
  1. set_issue
- server_responses:
  1. read_ok
  1. update_ok
  1. missing
  1. domain_failed

<!-- field: expected | form: "tree" -->
#### Expected

- get_signal: RESTResourceRead
- set_signal: RESTResourceWritten
- missing_signal: RESTMissing
- domain_failure_signal: RESTDomainFailed
- command_error_signal: CommandError
- output_contains_status_and_operation: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md) [S2](../use-cases/rel03.0-uc001-rest-client-sync.md#s2)
- [rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md) [S3](../use-cases/rel03.0-uc001-rest-client-sync.md#s3)
- [rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md) [S4](../use-cases/rel03.0-uc001-rest-client-sync.md#s4)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC4](../software-requirements/srd027-rest-tool-runtime.md#ac4)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC3](../software-requirements/srd028-rest-client-tools.md#ac3)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC6](../software-requirements/srd028-rest-client-tools.md#ac6)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC7](../software-requirements/srd028-rest-client-tools.md#ac7)

<!-- case: Runtime input cannot override configured authority | status: done -->
### Runtime input cannot override configured authority

A malicious input tries to replace transport authority at runtime. The
client rejects it before any HTTP request leaves the process.

<!-- field: use_case -->
#### Use case

[rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md)

<!-- field: go_test -->
#### Go test

TestRESTClient_RejectsAuthorityOverride

<!-- field: inputs | form: "tree" -->
#### Inputs

- operation: get_issue
- forbidden_runtime_fields:
  1. method
  1. url
  1. auth_ref
  1. redirect_policy
  1. host

<!-- field: expected | form: "tree" -->
#### Expected

- signal: CommandError
- failure_stage: schema_validation
- request_sent: false
- diagnostic_redacted: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md) [S1](../use-cases/rel03.0-uc001-rest-client-sync.md#s1)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC2](../software-requirements/srd027-rest-tool-runtime.md#ac2)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC6](../software-requirements/srd027-rest-tool-runtime.md#ac6)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC2](../software-requirements/srd028-rest-client-tools.md#ac2)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC9](../software-requirements/srd028-rest-client-tools.md#ac9)

<!-- case: Trusted operation selects a child service target | status: done -->
### Trusted operation selects a child service target

A service-start word records a loopback base_url under a labeled command
output. A configured health operation selects that value through
base_url_source command_state and performs one normal REST dispatch.
Unsafe values, network-policy violations, and credential-scope changes
fail before network I/O; runtime input still cannot set authority.

<!-- field: use_case -->
#### Use case

[rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md)

<!-- field: go_test -->
#### Go test

TestRESTClient_CommandStateSelectsTrustedBaseURL

<!-- field: inputs | form: "tree" -->
#### Inputs

- configured_base_url: http://127.0.0.1:1
- base_url_source: command_state
- base_url_selector: $from(start_subject).base_url
- operation_path: /healthz

<!-- field: expected | form: "tree" -->
#### Expected

- selected_authority_reported: true
- static_target_default_unchanged: true
- unsafe_target_rejected_before_io: true
- credential_scope_default_deny: true
- runtime_authority_override_rejected: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md) [S1](../use-cases/rel03.0-uc001-rest-client-sync.md#s1)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC2](../software-requirements/srd027-rest-tool-runtime.md#ac2)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC13](../software-requirements/srd028-rest-client-tools.md#ac13)

<!-- case: Trusted operation composes a target from a discovered host | status: done -->
### Trusted operation composes a target from a discovered host

A discovery word records a bare host or IP under a labeled command output,
the shape a listing endpoint returns for a per-item target. A configured
read operation selects that address through base_url_host_selector and
composes it with the declared base_url_scheme and base_url_port into the
request base URL, then performs one normal REST dispatch. Scheme and port
come only from trusted config: a resolved value carrying a scheme,
credential, port, path, or query is rejected before network I/O, both
selector forms on one operation are rejected, and an unsupported scheme or
out-of-range port is rejected at load.

<!-- field: use_case -->
#### Use case

[rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md)

<!-- field: go_test -->
#### Go test

TestRESTClient_HostSelectorComposesSelectedAuthority

<!-- field: inputs | form: "tree" -->
#### Inputs

- configured_base_url: http://127.0.0.1:1
- base_url_source: command_state
- base_url_host_selector: $from(discover_pods).ip
- base_url_scheme: http
- operation_path: /monitor/state

<!-- field: expected | form: "tree" -->
#### Expected

- composed_authority_reported: true
- unsafe_host_rejected_before_io: true
- both_selector_forms_rejected: true
- invalid_scheme_or_port_rejected_at_load: true
- ipv6_literal_bracketed: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md) [S1](../use-cases/rel03.0-uc001-rest-client-sync.md#s1)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC2](../software-requirements/srd027-rest-tool-runtime.md#ac2)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC14](../software-requirements/srd028-rest-client-tools.md#ac14)

<!-- case: Trusted operation resolves a composed target port per item | status: done -->
### Trusted operation resolves a composed target port per item

A fleet whose targets do not share one port publishes each target's port
beside its address, as a Kubernetes label string or a JSON number. A
configured read operation resolves that port through
base_url_port_selector and reaches each target's own port. Configuration
still bounds the reachable ports: a resolved port outside the client's
limits.network.ports allowlist is refused, and a malformed, fractional,
wrongly typed, or out-of-range port fails before network I/O. Declaring
both port forms, or a port selector without a host selector, is rejected
at load.

<!-- field: use_case -->
#### Use case

[rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md)

<!-- field: go_test -->
#### Go test

TestRESTClient_PortSelectorResolvesPerItemPort

<!-- field: inputs | form: "tree" -->
#### Inputs

- base_url_source: command_state
- base_url_host_selector: $from(discover_pods).ip
- base_url_port_selector: $from(discover_pods).port
- operation_path: /monitor/state

<!-- field: expected | form: "tree" -->
#### Expected

- string_and_number_ports_resolve: true
- per_item_authority_reported: true
- port_outside_allowlist_refused: true
- malformed_port_rejected_before_io: true
- both_port_forms_rejected_at_load: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md) [S1](../use-cases/rel03.0-uc001-rest-client-sync.md#s1)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC2](../software-requirements/srd027-rest-tool-runtime.md#ac2)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC15](../software-requirements/srd028-rest-client-tools.md#ac15)

<!-- case: REST client threads declared params from the previous Result | status: done -->
### REST client threads declared params from the previous Result

An operation declares body_source previous_result and an input_mapping
that selects declared params from a prior REST client Result output (for
example $.mapped.embedding into query_embeddings and $.carried.input into
a path param). After selection only declared params reach the request
builder, so the prior Result's fixed 11-key output no longer trips the
declared-only contract. The same case checks input carry-forward: a
producing operation copies selected declared inputs into its Result output
under a carried key, and the downstream operation selects them. Validation
rejects undeclared input_mapping targets, command_state body_source, and
any input_mapping or carry-forward entry that names transport authority.

<!-- field: use_case -->
#### Use case

[rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md)

<!-- field: go_test -->
#### Go test

TestRESTClient_ThreadsPreviousResultParams

<!-- field: inputs | form: "tree" -->
#### Inputs

- producer_operation: embed_query
- consumer_operation: query_collection
- body_source: previous_result
- input_mapping:
  - query_embeddings: $.mapped.embedding
  - collection: $.carried.input
- carry_forward:
  1. input
- rejection_cases:
  1. undeclared_input_mapping_target
  1. command_state_body_source
  1. transport_authority_override

<!-- field: expected | form: "tree" -->
#### Expected

- selection_presents_only_declared_params: true
- not_declared_error: false
- carry_forward_under_carried_key: true
- undeclared_target_rejected: true
- command_state_body_source_rejected: true
- transport_authority_override_rejected: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md) [S1](../use-cases/rel03.0-uc001-rest-client-sync.md#s1)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC11](../software-requirements/srd028-rest-client-tools.md#ac11)

<!-- case: REST client chain threads an embedding vector into a query body | status: done -->
### REST client chain threads an embedding vector into a query body

A two-word chain modeled on embed then query threads a numeric vector from
the first word's output into the second word's request body. The body
value ["{{ params.query_embeddings }}"] uses the single-token raw
passthrough so the flat vector renders as a nested array [[...]] rather
than a stringified scalar.

<!-- field: use_case -->
#### Use case

[rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md)

<!-- field: go_test -->
#### Go test

TestRESTClient_ChainThreadsVectorIntoQueryBody

<!-- field: inputs | form: "tree" -->
#### Inputs

- producer_operation: embed_query
- consumer_operation: query_collection
- body_template:
  - query_embeddings:
    1. "{{ params.query_embeddings }}"

<!-- field: expected | form: "tree" -->
#### Expected

- vector_threaded_into_body: true
- query_embeddings_rendered_as_nested_array: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md) [S1](../use-cases/rel03.0-uc001-rest-client-sync.md#s1)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC11](../software-requirements/srd028-rest-client-tools.md#ac11)

<!-- case: REST auth references resolve through trusted runtime credentials | status: done -->
### REST auth references resolve through trusted runtime credentials

Configure basic, bearer, header_token, and query_token auth profiles with
credential reference names. The runtime resolves synthetic credentials
through an injected resolver, sends resolved values to net/http, and does
not treat reference names as secret values.

<!-- field: use_case -->
#### Use case

[rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md)

<!-- field: go_test -->
#### Go test

TestRESTClient_ResolvesAuthCredentialRefs

<!-- field: inputs | form: "tree" -->
#### Inputs

- auth_profiles:
  1. basic
  1. bearer
  1. header_token
  1. query_token
- credential_refs:
  1. github_token
  1. user_ref
  1. password_ref

<!-- field: expected | form: "tree" -->
#### Expected

- resolved_credentials_sent: true
- literal_refs_sent: false
- result_output_contains_secret: false
- missing_ref_go_test: TestRESTClient_MissingCredentialReferenceFailsAuthResolution
- missing_ref_failure_stage: auth_resolution

<!-- field: traces -->
#### Traces

- [rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md) [S1](../use-cases/rel03.0-uc001-rest-client-sync.md#s1)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC6](../software-requirements/srd027-rest-tool-runtime.md#ac6)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC7](../software-requirements/srd027-rest-tool-runtime.md#ac7)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC2](../software-requirements/srd028-rest-client-tools.md#ac2)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC9](../software-requirements/srd028-rest-client-tools.md#ac9)

<!-- case: REST redirect and body size policies fail closed | status: done -->
### REST redirect and body size policies fail closed

Exercise outbound redirect and body-size policies with local httptest
servers. Allowlisted redirect hosts succeed. Unlisted redirect hosts fail.
Oversized request bodies fail before send, and oversized responses return
a size_limit CommandError without returning response content.

<!-- field: use_case -->
#### Use case

[rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md)

<!-- field: go_test -->
#### Go test

TestRESTClient_RedirectAllowlistPolicy

<!-- field: inputs | form: "tree" -->
#### Inputs

- redirect_mode: allowlist
- request_limit: max_request_bytes
- response_limit: max_response_bytes
- size_limit_go_test: TestRESTClient_RequestAndResponseSizeLimits

<!-- field: expected | form: "tree" -->
#### Expected

- allowed_redirect_succeeds: true
- blocked_redirect_signal: CommandError
- oversized_request_sent: false
- oversized_response_failure_stage: size_limit
- oversized_response_content_returned: false

<!-- field: traces -->
#### Traces

- [rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md) [S1](../use-cases/rel03.0-uc001-rest-client-sync.md#s1)
- [rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md) [S2](../use-cases/rel03.0-uc001-rest-client-sync.md#s2)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC6](../software-requirements/srd027-rest-tool-runtime.md#ac6)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC7](../software-requirements/srd027-rest-tool-runtime.md#ac7)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC2](../software-requirements/srd028-rest-client-tools.md#ac2)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC8](../software-requirements/srd028-rest-client-tools.md#ac8)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC9](../software-requirements/srd028-rest-client-tools.md#ac9)

<!-- case: Mutating REST operations declare effects and compensation | status: done -->
### Mutating REST operations declare effects and compensation

Review mutating resource operations beside an RPC-shaped invoke
operation. Each one must declare its side-effect contract and
reversibility before selection; compensation or confirmation metadata
completes the contract.

<!-- field: use_case -->
#### Use case

[rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md)

<!-- field: go_test -->
#### Go test

TestRESTClient_MutatingOperationsRequireEffects

<!-- field: inputs | form: "tree" -->
#### Inputs

- operations:
  1. rest_client_create
  1. rest_client_set
  1. rest_client_delete
  1. rest_client_invoke
- invalid_case: missing_reversibility

<!-- field: expected | form: "tree" -->
#### Expected

- valid_mutations_accepted: true
- missing_reversibility_rejected: true
- irreversible_requires_confirmation: true
- compensation_metadata_captured: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md) [S3](../use-cases/rel03.0-uc001-rest-client-sync.md#s3)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC8](../software-requirements/srd027-rest-tool-runtime.md#ac8)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC3](../software-requirements/srd028-rest-client-tools.md#ac3)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC4](../software-requirements/srd028-rest-client-tools.md#ac4)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC8](../software-requirements/srd028-rest-client-tools.md#ac8)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC10](../software-requirements/srd028-rest-client-tools.md#ac10)

<!-- case: Async send records request state and returns RESTAccepted | status: done -->
### Async send records request state and returns RESTAccepted

Dispatch rest_client_send for an async operation and confirm the command
waits for the submission response. AsyncRequest state is stored by request
ID only after an accepted status, then RESTAccepted carries status and
correlation metadata plus the remote compensation receipt back to the grammar.

<!-- field: use_case -->
#### Use case

[rel03.0-uc002-rest-client-async](../use-cases/rel03.0-uc002-rest-client-async.md)

<!-- field: go_test -->
#### Go test

TestRESTClient_SendRecordsAsyncRequest

<!-- field: inputs | form: "tree" -->
#### Inputs

- operation: create_payment
- async_config:
  - request_id: $.id
  - idempotency_token: "{{ params.order_id }}"
  - state_retention: consume

<!-- field: expected | form: "tree" -->
#### Expected

- signal: RESTAccepted
- request_state_written: true
- request_id_present: true
- idempotency_metadata_present: true
- submission_response_waited_for: true
- compensation_receipt_present: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc002-rest-client-async](../use-cases/rel03.0-uc002-rest-client-async.md) [S1](../use-cases/rel03.0-uc002-rest-client-async.md#s1)
- [rel03.0-uc002-rest-client-async](../use-cases/rel03.0-uc002-rest-client-async.md) [S2](../use-cases/rel03.0-uc002-rest-client-async.md#s2)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC5](../software-requirements/srd027-rest-tool-runtime.md#ac5)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC5](../software-requirements/srd028-rest-client-tools.md#ac5)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC8](../software-requirements/srd028-rest-client-tools.md#ac8)

<!-- case: Async await maps completion and timeout signals | status: done -->
### Async await maps completion and timeout signals

Begin with stored AsyncRequest state. Completion, domain failure,
absence, timeout, and missing state each map to a declared signal.
The table uses separate outcomes so RESTDomainFailed is not confused
with CommandError or the timeout path.
Await behavior also proves that stored async state is consumed according
to policy after the final signal is chosen.

<!-- field: use_case -->
#### Use case

[rel03.0-uc002-rest-client-async](../use-cases/rel03.0-uc002-rest-client-async.md)

<!-- field: go_test -->
#### Go test

TestRESTClient_AwaitAsyncRequest

<!-- field: inputs | form: "tree" -->
#### Inputs

- request_state: pending_payment
- outcomes:
  1. responded
  1. domain_failed
  1. missing
  1. timed_out
  1. async_state_missing

<!-- field: expected | form: "tree" -->
#### Expected

- responded_signal: RESTResponded
- domain_failed_signal: RESTDomainFailed
- missing_signal: RESTMissing
- timeout_signal: RESTAwaitTimedOut
- missing_state_signal: CommandError
- state_retention_applied: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc002-rest-client-async](../use-cases/rel03.0-uc002-rest-client-async.md) [S1](../use-cases/rel03.0-uc002-rest-client-async.md#s1)
- [rel03.0-uc002-rest-client-async](../use-cases/rel03.0-uc002-rest-client-async.md) [S3](../use-cases/rel03.0-uc002-rest-client-async.md#s3)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC5](../software-requirements/srd027-rest-tool-runtime.md#ac5)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC5](../software-requirements/srd028-rest-client-tools.md#ac5)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC6](../software-requirements/srd028-rest-client-tools.md#ac6)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC9](../software-requirements/srd028-rest-client-tools.md#ac9)

<!-- case: Async retry policy requires idempotency metadata | status: done -->
### Async retry policy requires idempotency metadata

For async operations, check retry policy at config validation time.
Non-idempotent requests reject retry settings unless config declares
safe duplicate suppression.

<!-- field: use_case -->
#### Use case

[rel03.0-uc002-rest-client-async](../use-cases/rel03.0-uc002-rest-client-async.md)

<!-- field: go_test -->
#### Go test

TestRESTClient_AsyncRetryPolicyValidation

<!-- field: inputs | form: "tree" -->
#### Inputs

- operation: create_payment
- retry_policy: retry_on_429
- invalid_case: retry_without_idempotency

<!-- field: expected | form: "tree" -->
#### Expected

- safe_retry_accepted: true
- unsafe_retry_rejected: true
- retry_count_traced: true
- idempotency_token_redacted: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc002-rest-client-async](../use-cases/rel03.0-uc002-rest-client-async.md) [S4](../use-cases/rel03.0-uc002-rest-client-async.md#s4)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC6](../software-requirements/srd027-rest-tool-runtime.md#ac6)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC8](../software-requirements/srd028-rest-client-tools.md#ac8)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC9](../software-requirements/srd028-rest-client-tools.md#ac9)

<!-- case: REST retry aggregate stays within the machine command envelope | status: done -->
### REST retry aggregate stays within the machine command envelope

Inspect selected REST operations without sending requests. The timeout
authority adds every bounded request attempt and the fixed or capped
exponential delay before each retry. Equality and overrun fail before
dispatch, as do a missing retry_ref target, a non-positive attempt count,
malformed delay, and checked-arithmetic overflow. A one-attempt operation
remains the per-attempt request bound, and rest_client_await keeps only its
separate async timeout.

<!-- field: use_case -->
#### Use case

[rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md)

<!-- field: go_test -->
#### Go test

go test ./pkg/profileaudit -run TestInspectRESTRetryAggregate

<!-- field: inputs | form: "tree" -->
#### Inputs

- retry_cases:
  1. one_attempt
  1. fixed_below_envelope
  1. exponential_capped_below_envelope
  1. aggregate_equal_to_envelope
  1. aggregate_above_envelope
  1. missing_retry_policy
  1. non_positive_attempts
  1. malformed_delay
  1. duration_overflow

<!-- field: expected | form: "tree" -->
#### Expected

- network_requests_sent: 0
- aggregate_matches_runtime_backoff: true
- deterministic_client_policy_evidence: true
- equal_and_larger_rejected: true
- invalid_and_overflowing_policy_rejected: true
- await_submission_aggregate_applied: false

<!-- field: traces -->
#### Traces

- [rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md) [S1](../use-cases/rel03.0-uc001-rest-client-sync.md#s1)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC6](../software-requirements/srd027-rest-tool-runtime.md#ac6)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC2](../software-requirements/srd028-rest-client-tools.md#ac2)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC8](../software-requirements/srd028-rest-client-tools.md#ac8)

<!-- case: REST dispatch cancellation reaches request and retry transport | status: done -->
### REST dispatch cancellation reaches request and retry transport

Cancel dispatch during an in-flight HTTP request and a declared retry
wait. Both stop promptly, while the legacy Execute entry point retains a
background context and an uncertain non-idempotent mutation reports an
indeterminate remote outcome.

<!-- field: use_case -->
#### Use case

[rel03.0-uc002-rest-client-async](../use-cases/rel03.0-uc002-rest-client-async.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest/client -run 'TestRESTClient(DispatchCancellation|CancelledMutation|ExecuteRetains)'

<!-- field: expected | form: "tree" -->
#### Expected

- in_flight_request_aborted: true
- retry_wait_aborted: true
- context_worker_joined: true
- non_idempotent_outcome: indeterminate
- legacy_execute_compatible: true

<!-- field: traces -->
#### Traces

- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC6](../software-requirements/srd027-rest-tool-runtime.md#ac6)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC8](../software-requirements/srd028-rest-client-tools.md#ac8)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC9](../software-requirements/srd028-rest-client-tools.md#ac9)

<!-- case: Server launch registers routes without blocking | status: done -->
### Server launch registers routes without blocking

Dispatch rest_server_launch with a listener on port 0. The command binds
the server, registers configured routes, records shared state, and
returns ServerLaunched without waiting for inbound requests.

<!-- field: use_case -->
#### Use case

[rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md)

<!-- field: go_test -->
#### Go test

TestRESTServer_LaunchRegistersRoutes

<!-- field: inputs | form: "tree" -->
#### Inputs

- server: control_api
- address: 127.0.0.1:0
- endpoints:
  1. approve
  1. reject
  1. pause

<!-- field: expected | form: "tree" -->
#### Expected

- signal: ServerLaunched
- route_count: 3
- shared_server_state_written: true
- network_listen_side_effect_declared: true
- blocks_for_inbound_request: false

<!-- field: traces -->
#### Traces

- [rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md) [S1](../use-cases/rel03.0-uc003-rest-server-signals.md#s1)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC5](../software-requirements/srd027-rest-tool-runtime.md#ac5)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC2](../software-requirements/srd029-rest-server-tools.md#ac2)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC7](../software-requirements/srd029-rest-server-tools.md#ac7)

<!-- case: Server await emits queued inbound signals | status: done -->
### Server await emits queued inbound signals

Send validated inbound HTTP requests to signal-producing endpoints.
rest_server_await consumes queued events and returns the configured
signal to MachineSpec.

<!-- field: use_case -->
#### Use case

[rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md)

<!-- field: go_test -->
#### Go test

TestRESTServer_AwaitInboundSignals

<!-- field: inputs | form: "tree" -->
#### Inputs

- server: control_api
- endpoint_bindings:
  1. emit_signal
  1. emit_dynamic_signal
- queued_events:
  1. Approved
  1. Rejected
  1. PauseRequested
  1. DomainEventReceived

<!-- field: expected | form: "tree" -->
#### Expected

- approved_route_visible: true
- rejected_route_visible: true
- pause_route_visible: true
- domain_event_route_visible: true
- invalid_request_enqueued: false

<!-- field: traces -->
#### Traces

- [rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md) [S2](../use-cases/rel03.0-uc003-rest-server-signals.md#s2)
- [rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md) [S3](../use-cases/rel03.0-uc003-rest-server-signals.md#s3)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC3](../software-requirements/srd029-rest-server-tools.md#ac3)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC5](../software-requirements/srd029-rest-server-tools.md#ac5)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC6](../software-requirements/srd029-rest-server-tools.md#ac6)

<!-- case: Generic REST await fan-in selects one event across servers | status: done -->
### Generic REST await fan-in selects one event across servers

Launch two configured REST servers with separate queues. Enqueue accepted
events on each server, then dispatch rest_await_event with both sources.
One Result carries the selected event. MachineSpec owns the workflow
route after that signal.

<!-- field: use_case -->
#### Use case

[rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md)

<!-- field: go_test -->
#### Go test

TestRESTAwaitEvent_MultiSourceFanIn

<!-- field: inputs | form: "tree" -->
#### Inputs

- servers:
  1. control_api
  1. domain_webhooks
- queued_events:
  1. - server: control_api
     - signal: Approved
  1. - server: domain_webhooks
     - signal: DomainEventReceived

<!-- field: expected | form: "tree" -->
#### Expected

- selected_event_returned_once: true
- output_contains_source_and_route: true
- hidden_workflow_action_selected: false

<!-- field: traces -->
#### Traces

- [rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md) [S2](../use-cases/rel03.0-uc003-rest-server-signals.md#s2)
- [rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md) [S3](../use-cases/rel03.0-uc003-rest-server-signals.md#s3)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC10](../software-requirements/srd027-rest-tool-runtime.md#ac10)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC10](../software-requirements/srd029-rest-server-tools.md#ac10)

<!-- case: Staged REST await-any fan-in adds a later source | status: done -->
### Staged REST await-any fan-in adds a later source

Build a MachineSpec-shaped factory path with three configured REST server
sources. Launch the first two servers, dispatch rest_await_event across
those two sources, and route the selected event to a visible third launch.
Then dispatch rest_await_event across all three launched sources. The
test should prove that each await returns one event and that the third
source is included only after the first selected event routes through the
machine.

<!-- field: use_case -->
#### Use case

[rel03.0-uc008-rest-await-any-fan-in](../use-cases/rel03.0-uc008-rest-await-any-fan-in.md)

<!-- field: go_test -->
#### Go test

TestRESTAwaitEvent_FactoryBuildsStagedFanIn

<!-- field: inputs | form: "tree" -->
#### Inputs

- initial_servers:
  1. control_api
  1. domain_webhooks
- later_server: audit_webhooks
- first_await_sources:
  1. control_api
  1. domain_webhooks
- second_await_sources:
  1. control_api
  1. domain_webhooks
  1. audit_webhooks
- factory_init: rest_await_event

<!-- field: expected | form: "tree" -->
#### Expected

- first_selected_event_returned_once: true
- third_server_launch_is_visible_transition: true
- second_selected_event_returned_once: true
- second_await_includes_later_source: true
- hidden_workflow_action_selected: false

<!-- field: traces -->
#### Traces

- [rel03.0-uc008-rest-await-any-fan-in](../use-cases/rel03.0-uc008-rest-await-any-fan-in.md) [S1](../use-cases/rel03.0-uc008-rest-await-any-fan-in.md#s1)
- [rel03.0-uc008-rest-await-any-fan-in](../use-cases/rel03.0-uc008-rest-await-any-fan-in.md) [S2](../use-cases/rel03.0-uc008-rest-await-any-fan-in.md#s2)
- [rel03.0-uc008-rest-await-any-fan-in](../use-cases/rel03.0-uc008-rest-await-any-fan-in.md) [S3](../use-cases/rel03.0-uc008-rest-await-any-fan-in.md#s3)
- [rel03.0-uc008-rest-await-any-fan-in](../use-cases/rel03.0-uc008-rest-await-any-fan-in.md) [S4](../use-cases/rel03.0-uc008-rest-await-any-fan-in.md#s4)
- [rel03.0-uc008-rest-await-any-fan-in](../use-cases/rel03.0-uc008-rest-await-any-fan-in.md) [S5](../use-cases/rel03.0-uc008-rest-await-any-fan-in.md#s5)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC10](../software-requirements/srd027-rest-tool-runtime.md#ac10)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC10](../software-requirements/srd029-rest-server-tools.md#ac10)

<!-- case: Generic REST await applies route and signal filters | status: done -->
### Generic REST await applies route and signal filters

Configure rest_await_event with server, route, and signal filters. Events
outside the route or signal allowlist remain queued or are skipped
according to policy. Matching input drives the emitted signal.

<!-- field: use_case -->
#### Use case

[rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md)

<!-- field: go_test -->
#### Go test

TestRESTAwaitEvent_SourceFiltersPreserveUnrelatedEvents

<!-- field: inputs | form: "tree" -->
#### Inputs

- sources:
  1. - server: control_api
     - routes:
       1. approve
       1. reject
     - signals:
       1. Approved
       1. Rejected
- queued_events:
  1. - route: pause
     - signal: PauseRequested
  1. - route: approve
     - signal: Approved

<!-- field: expected | form: "tree" -->
#### Expected

- disallowed_route_selected: false
- allowed_signal_selected: Approved
- filter_diagnostics_reported: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md) [S3](../use-cases/rel03.0-uc003-rest-server-signals.md#s3)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC10](../software-requirements/srd027-rest-tool-runtime.md#ac10)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC10](../software-requirements/srd029-rest-server-tools.md#ac10)

<!-- case: Generic REST await reports timeout without consuming events | status: done -->
### Generic REST await reports timeout without consuming events

Dispatch rest_await_event with no matching events before the configured
timeout. AwaitTimedOut reports the active source filters
without inventing a lifecycle or domain signal.

<!-- field: use_case -->
#### Use case

[rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md)

<!-- field: go_test -->
#### Go test

TestRESTAwaitEvent_Timeout

<!-- field: inputs | form: "tree" -->
#### Inputs

- sources:
  1. - server: control_api
     - routes:
       1. approve
     - signals:
       1. Approved
- timeout: 10ms

<!-- field: expected | form: "tree" -->
#### Expected

- signal: AwaitTimedOut
- matching_event_consumed: false
- hidden_workflow_action_selected: false

<!-- field: traces -->
#### Traces

- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC10](../software-requirements/srd027-rest-tool-runtime.md#ac10)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC10](../software-requirements/srd029-rest-server-tools.md#ac10)

<!-- case: Generic REST await handles stopped sources by policy | status: done -->
### Generic REST await handles stopped sources by policy

Stop one configured source while rest_await_event is waiting. Run the
ignore, emit_server_stopped, and command_error policies and verify each
policy returns the declared signal or continues waiting as configured.

<!-- field: use_case -->
#### Use case

[rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md)

<!-- field: go_test -->
#### Go test

TestRESTAwaitEvent_ServerStopped

<!-- field: inputs | form: "tree" -->
#### Inputs

- stopped_behaviors:
  1. ignore
  1. emit_server_stopped
  1. command_error

<!-- field: expected | form: "tree" -->
#### Expected

- ignore_continues_to_other_sources: true
- emit_server_stopped_signal: ServerStopped
- command_error_signal: CommandError
- command_error_go_test: TestRESTAwaitEvent_StoppedSourceCommandError

<!-- field: traces -->
#### Traces

- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC10](../software-requirements/srd027-rest-tool-runtime.md#ac10)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC10](../software-requirements/srd029-rest-server-tools.md#ac10)

<!-- case: Lifecycle exit through REST routes to exit_agent | status: done -->
### Lifecycle exit through REST routes to exit_agent

A configured lifecycle_control exit endpoint receives an operator
request. Handler validation enqueues ExitRequested. rest_await_event
returns that event, and MachineSpec routes it to exit_agent as explicit
lifecycle vocabulary. Command-level coverage also proves that runtime
payload fields reach structured exit_agent output, shutdown is requested
by exit_agent, and cancellation is deferred until after the loop records
AgentExited and resolves the terminal MachineSpec state.

<!-- field: use_case -->
#### Use case

[rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md)

<!-- field: go_test -->
#### Go test

TestRESTLifecycleControl_ExitAgentSignal

<!-- field: inputs | form: "tree" -->
#### Inputs

- endpoint: POST /api/lifecycle/exit
- payload:
  - reason: operator requested shutdown
  - status: success
- machine_route:
  - signal: ExitRequested
  - action: exit_agent

<!-- field: expected | form: "tree" -->
#### Expected

- handler_executed_exit_directly: false
- await_signal: ExitRequested
- exit_agent_action_visible: true
- exit_result_signal: AgentExited
- exit_output_json: true
- exit_output_reason_from_event: true
- exit_output_status_from_event: true
- control_profile_terminal_state: Succeeded
- run_status_cancelled: false
- deferred_shutdown_after_loop_return: true
- control_profile_go_test: TestControlProfileExitReachesSucceededBeforeDeferredShutdown
- payload_go_test: TestExitAgentUsesRuntimeEventPayload
- profile_owned_exit_tracer: applications/catalog documentation-curator suite

<!-- field: traces -->
#### Traces

- [rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md) [S3](../use-cases/rel03.0-uc003-rest-server-signals.md#s3)
- [srd026-lifecycle-tools](../software-requirements/srd026-lifecycle-tools.md) [AC7](../software-requirements/srd026-lifecycle-tools.md#ac7)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC10](../software-requirements/srd027-rest-tool-runtime.md#ac10)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC5](../software-requirements/srd029-rest-server-tools.md#ac5)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC10](../software-requirements/srd029-rest-server-tools.md#ac10)

<!-- case: Server stop drains queues and unblocks await | status: done -->
### Server stop drains queues and unblocks await

Stop a launched server while events and await calls exist. Shutdown
follows the accepted drain policies, rejects unsupported shutdown fields,
closes listeners, and unblocks await and stream_events with ServerStopped.

<!-- field: use_case -->
#### Use case

[rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md)

<!-- field: go_test -->
#### Go test

TestRESTServer_StopDrainsAndUnblocks, TestRESTServer_ShutdownConfigValidation, TestRESTServer_StreamEventsUnblocksOnStop

<!-- field: inputs | form: "tree" -->
#### Inputs

- server: control_api
- drain_policies:
  1. drain
  1. drain_then_stop
- rejected_drain_policies:
  1. reject_new
  1. drop_queued
  1. fail_queued
- pending_events: 2
- awaiting_call: true
- active_stream: true

<!-- field: expected | form: "tree" -->
#### Expected

- signal: ServerStopped
- drained_event_count: 2
- listener_closed: true
- await_unblocked: true
- stream_unblocked: true
- unsupported_shutdown_fields_rejected: true
- dropped_events_reported: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md) [S4](../use-cases/rel03.0-uc003-rest-server-signals.md#s4)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC4](../software-requirements/srd029-rest-server-tools.md#ac4)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC7](../software-requirements/srd029-rest-server-tools.md#ac7)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC9](../software-requirements/srd029-rest-server-tools.md#ac9)

<!-- case: Terminal state mapping separates client and server errors | status: done -->
### Terminal state mapping separates client and server errors

A request machine built from exec words emits only ToolDone and ToolFailed,
so two different failures reach two different terminal states through the
same signal. Keyed by terminal state the endpoint returns a distinct
response for each: a validate reject is the caller's fault, an apply
failure is the server's. State resolves before signal when both maps are
configured, and a config with only terminal_signals behaves as before.

<!-- field: use_case -->
#### Use case

[rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md)

<!-- field: go_test -->
#### Go test

TestMachineRequestTerminalStateSeparatesClientAndServerErrors, TestMachineRequestTerminalStateOutranksSignal, TestMachineRequestFallsBackToSignalMapping

<!-- field: inputs | form: "tree" -->
#### Inputs

- machine: apply-shaped, terminal states Done, Rejected, Failed
- words: exec-shaped, emitting only ToolDone or ToolFailed
- response_map:
  - Done: 200
  - Rejected: 400
  - Failed: 500

<!-- field: expected | form: "tree" -->
#### Expected

- validate_reject_status: 400
- apply_failure_status: 500
- success_status: 200
- failures_share_terminal_signal: true
- terminal_state_outranks_terminal_signal: true
- signal_only_config_unchanged: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md) [S2](../use-cases/rel03.0-uc003-rest-server-signals.md#s2)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC8](../software-requirements/srd030-machine-request-rest-binding.md#ac8)

<!-- case: Terminal response configuration fails closed | status: done -->
### Terminal response configuration fails closed

A mapping onto a state the machine never terminates in can never fire, so
it fails at load rather than surfacing as a per-request response_missing.
One of the two response maps is required, either alone suffices, and an
unmapped terminal names both the state and the signal that were looked up.

<!-- field: use_case -->
#### Use case

[rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md)

<!-- field: go_test -->
#### Go test

TestValidateMachineResponsesRejectsNonTerminalState, TestMachineRequestEndpointAcceptsEitherResponseMap, TestMachineRequestUnmappedTerminalNamesBothKeys

<!-- field: inputs | form: "tree" -->
#### Inputs

- non_terminal_state_mapping: Applying
- empty_response_map: true
- unmapped_terminal: Rejected

<!-- field: expected | form: "tree" -->
#### Expected

- non_terminal_state_rejected_at_load: true
- either_map_accepted: true
- neither_map_rejected: true
- diagnostic_names_state_and_signal: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md) [S2](../use-cases/rel03.0-uc003-rest-server-signals.md#s2)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC8](../software-requirements/srd030-machine-request-rest-binding.md#ac8)

<!-- case: Server validation rejects unsafe inbound traffic | status: done -->
### Server validation rejects unsafe inbound traffic

Bad inbound traffic reaches the configured routes. The server returns
configured HTTP errors and avoids hidden MachineSpec actions.

<!-- field: use_case -->
#### Use case

[rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md)

<!-- field: go_test -->
#### Go test

TestRESTServer_RequestValidationFailures

<!-- field: inputs | form: "tree" -->
#### Inputs

- invalid_requests:
  1. wrong_method
  1. schema_invalid
  1. body_too_large
  1. queue_full
  1. handler_failed

<!-- field: expected | form: "tree" -->
#### Expected

- grammar_event_enqueued: false
- configured_http_status_returned: true
- queue_overflow_reported: true
- handler_failure_signal_or_response_declared: true
- diagnostics_redacted: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md) [S2](../use-cases/rel03.0-uc003-rest-server-signals.md#s2)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC6](../software-requirements/srd027-rest-tool-runtime.md#ac6)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC6](../software-requirements/srd029-rest-server-tools.md#ac6)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC9](../software-requirements/srd029-rest-server-tools.md#ac9)

<!-- case: Endpoints sharing a path are selected by method | status: done -->
### Endpoints sharing a path are selected by method

One server declares GET and POST on one path, a POST-only path, and a GET
catch-all. Each method reaches its own endpoint; an undeclared method on
the shared path answers 405 with Allow naming both; a GET on the POST-only
path answers 405 rather than falling to the catch-all; other GETs reach
the catch-all. Two endpoints declaring one method on one path shape fail
definition validation, and a parameter beside a catch-all does not.

<!-- field: use_case -->
#### Use case

[rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest -run 'TestRESTServer_SelectsSamePathEndpointByMethod|TestRESTDefinition_RejectsSameMethodOnSamePath'

<!-- field: inputs | form: "tree" -->
#### Inputs

- endpoints:
  1. GET /api/items
  1. POST /api/items
  1. POST /api/chat
  1. GET /{path...}
- requests:
  1. GET /api/items
  1. POST /api/items
  1. PUT /api/items
  1. GET /api/chat
  1. GET /anything/else

<!-- field: expected | form: "tree" -->
#### Expected

- get_items_status: 200
- post_items_signal: ItemCreated
- put_items:
  - status: 405
  - allow: GET, POST
- get_chat:
  - status: 405
  - allow: POST
- catch_all_status: 200
- duplicate_method_and_shape_rejected: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md) [S2](../use-cases/rel03.0-uc003-rest-server-signals.md#s2)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC11](../software-requirements/srd029-rest-server-tools.md#ac11)

<!-- case: Server query header and output redaction policies fail closed | status: done -->
### Server query header and output redaction policies fail closed

Send unsafe inbound requests with undeclared query and header fields, then
send schema-invalid path and header values. The server rejects them before
enqueueing grammar events. Separate redaction checks prove await output,
SSE output, and handler responses redact configured query, header, and
body secrets.

<!-- field: use_case -->
#### Use case

[rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md)

<!-- field: go_test -->
#### Go test

TestRESTServer_RejectsUndeclaredQueryAndHeader

<!-- field: inputs | form: "tree" -->
#### Inputs

- invalid_requests:
  1. undeclared_query
  1. undeclared_header
  1. schema_invalid_path
  1. schema_invalid_header
- redaction_go_tests:
  1. TestRESTServer_RedactsAwaitAndStreamOutput
  1. TestRESTServer_RedactsHandlerResponses

<!-- field: expected | form: "tree" -->
#### Expected

- invalid_request_status: 400
- invalid_request_enqueued: false
- await_output_redacted: true
- stream_output_redacted: true
- handler_response_redacted: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md) [S2](../use-cases/rel03.0-uc003-rest-server-signals.md#s2)
- [rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md) [S3](../use-cases/rel03.0-uc003-rest-server-signals.md#s3)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC6](../software-requirements/srd027-rest-tool-runtime.md#ac6)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC7](../software-requirements/srd027-rest-tool-runtime.md#ac7)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC6](../software-requirements/srd029-rest-server-tools.md#ac6)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC9](../software-requirements/srd029-rest-server-tools.md#ac9)

<!-- case: Server invoke handler returns configured responses | status: done -->
### Server invoke handler returns configured responses

Send requests to configured invoke_handler endpoints. The server validates
requests, returns mapped JSON, and enqueues only a declared signal when
the endpoint configuration names one.

<!-- field: use_case -->
#### Use case

[rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md)

<!-- field: go_test -->
#### Go test

TestRESTServer_InvokeHandlerBindings

<!-- field: inputs | form: "tree" -->
#### Inputs

- endpoint_bindings:
  1. invoke_handler
- valid_request: handler_response
- signal_request: handler_signal
- invalid_request: handler_not_configured

<!-- field: expected | form: "tree" -->
#### Expected

- mapped_response_returned: true
- declared_signal_enqueued: true
- unconfigured_handler_fails: true
- hidden_workflow_action_selected: false

<!-- field: traces -->
#### Traces

- [rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md) [S2](../use-cases/rel03.0-uc003-rest-server-signals.md#s2)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC5](../software-requirements/srd029-rest-server-tools.md#ac5)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC6](../software-requirements/srd029-rest-server-tools.md#ac6)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC9](../software-requirements/srd029-rest-server-tools.md#ac9)

<!-- case: Server streams queued events as SSE | status: done -->
### Server streams queued events as SSE

Queue a validated inbound event, then call a stream_events endpoint. The
endpoint writes one server-sent event containing the configured signal and
route metadata without choosing the next MachineSpec action.

<!-- field: use_case -->
#### Use case

[rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md)

<!-- field: go_test -->
#### Go test

TestRESTServer_StreamEvents

<!-- field: inputs | form: "tree" -->
#### Inputs

- endpoint_bindings:
  1. emit_signal
  1. stream_events
- queued_event: Approved

<!-- field: expected | form: "tree" -->
#### Expected

- content_type: text/event-stream
- streamed_signal: Approved
- streamed_route: approve
- hidden_workflow_action_selected: false

<!-- field: traces -->
#### Traces

- [rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md) [S2](../use-cases/rel03.0-uc003-rest-server-signals.md#s2)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC5](../software-requirements/srd029-rest-server-tools.md#ac5)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC7](../software-requirements/srd029-rest-server-tools.md#ac7)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC9](../software-requirements/srd029-rest-server-tools.md#ac9)

<!-- case: OpenAPI import exposes allowlisted operations | status: done -->
### OpenAPI import exposes allowlisted operations

Load an OpenAPI 3 document with several operations, but import only the
allowlisted operation IDs. The import compiles them into RestDefinition
while ToolDef declarations stay explicit.

<!-- field: use_case -->
#### Use case

[rel03.0-uc004-openapi-import](../use-cases/rel03.0-uc004-openapi-import.md)

<!-- field: go_test -->
#### Go test

TestRESTOpenAPI_ImportAllowlist

<!-- field: inputs | form: "tree" -->
#### Inputs

- openapi_document: payments_openapi
- expose:
  1. getPayment
  1. createPayment
- tooldefs:
  1. get_payment
  1. create_payment

<!-- field: expected | form: "tree" -->
#### Expected

- imported_operations:
  1. getPayment
  1. createPayment
- unlisted_operation_imported: false
- tooldefs_explicit: true
- generated_go_code: false

<!-- field: traces -->
#### Traces

- [rel03.0-uc004-openapi-import](../use-cases/rel03.0-uc004-openapi-import.md) [S1](../use-cases/rel03.0-uc004-openapi-import.md#s1)
- [rel03.0-uc004-openapi-import](../use-cases/rel03.0-uc004-openapi-import.md) [S3](../use-cases/rel03.0-uc004-openapi-import.md#s3)
- [rel03.0-uc004-openapi-import](../use-cases/rel03.0-uc004-openapi-import.md) [S4](../use-cases/rel03.0-uc004-openapi-import.md#s4)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC3](../software-requirements/srd027-rest-tool-runtime.md#ac3)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC9](../software-requirements/srd027-rest-tool-runtime.md#ac9)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC1](../software-requirements/srd028-rest-client-tools.md#ac1)

<!-- case: OpenAPI server bind maps operation IDs to endpoints | status: done -->
### OpenAPI server bind maps operation IDs to endpoints

Bind selected OpenAPI operation IDs to configured server endpoint names.
Imported paths and schemas attach to endpoint binding kinds while server
ToolDefs remain explicit vocabulary words.

<!-- field: use_case -->
#### Use case

[rel03.0-uc004-openapi-import](../use-cases/rel03.0-uc004-openapi-import.md)

<!-- field: go_test -->
#### Go test

TestRESTOpenAPI_ServerBindMap

<!-- field: inputs | form: "tree" -->
#### Inputs

- openapi_document: payments_openapi
- bind:
  - receivePaymentWebhook: payment_webhook
- endpoint_binding: emit_signal

<!-- field: expected | form: "tree" -->
#### Expected

- endpoint_bound: payment_webhook
- request_schema_compiled: true
- response_schema_compiled: true
- server_tooldef_explicit: true
- generated_server_code: false

<!-- field: traces -->
#### Traces

- [rel03.0-uc004-openapi-import](../use-cases/rel03.0-uc004-openapi-import.md) [S2](../use-cases/rel03.0-uc004-openapi-import.md#s2)
- [rel03.0-uc004-openapi-import](../use-cases/rel03.0-uc004-openapi-import.md) [S3](../use-cases/rel03.0-uc004-openapi-import.md#s3)
- [rel03.0-uc004-openapi-import](../use-cases/rel03.0-uc004-openapi-import.md) [S4](../use-cases/rel03.0-uc004-openapi-import.md#s4)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC8](../software-requirements/srd029-rest-server-tools.md#ac8)

<!-- case: OpenAPI import reports invalid operation IDs | status: done -->
### OpenAPI import reports invalid operation IDs

When an import names a missing or duplicated operation ID, validation
fails before ToolDefs register. An incompatible server binding fails the
same way, with document and operation context in the error.

<!-- field: use_case -->
#### Use case

[rel03.0-uc004-openapi-import](../use-cases/rel03.0-uc004-openapi-import.md)

<!-- field: go_test -->
#### Go test

TestRESTOpenAPI_InvalidOperationIDs

<!-- field: inputs | form: "tree" -->
#### Inputs

- invalid_imports:
  1. missing_operation_id
  1. duplicate_operation_id
  1. incompatible_server_binding

<!-- field: expected | form: "tree" -->
#### Expected

- import_rejected: true
- source_document_reported: true
- operation_identity_reported: true
- schema_location_reported_when_available: true
- no_partial_tooldefs_registered: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc004-openapi-import](../use-cases/rel03.0-uc004-openapi-import.md) [S1](../use-cases/rel03.0-uc004-openapi-import.md#s1)
- [rel03.0-uc004-openapi-import](../use-cases/rel03.0-uc004-openapi-import.md) [S2](../use-cases/rel03.0-uc004-openapi-import.md#s2)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC3](../software-requirements/srd027-rest-tool-runtime.md#ac3)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC7](../software-requirements/srd027-rest-tool-runtime.md#ac7)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC8](../software-requirements/srd029-rest-server-tools.md#ac8)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC9](../software-requirements/srd029-rest-server-tools.md#ac9)

<!-- case: Qwen agent lists Ollama models through OpenAPI REST tool | status: done -->
### Qwen agent lists Ollama models through OpenAPI REST tool

Run the Ollama REST profile with a Qwen 3 family model and the
ollama_list_models tool selected. Ask the agent to list local Ollama
models. The integration target must prove that the model invoked the REST
tool, that the REST tool called GET /api/tags through the imported
listOllamaModels operation, and that the final answer contains model
names returned by /api/tags.

<!-- field: use_case -->
#### Use case

[rel03.0-uc005-ollama-openapi-model-listing](../use-cases/rel03.0-uc005-ollama-openapi-model-listing.md)

<!-- field: go_test -->
#### Go test

mage integration:ollamaRest

<!-- field: inputs | form: "tree" -->
#### Inputs

- profile: ${AGENT_CATALOG_ROOT}/rest/ollama-profile.yaml
- model_family: qwen3
- openapi_document: rest/openapi/ollama.yaml
- operation_id: listOllamaModels
- tool_name: ollama_list_models
- live_endpoint: http://127.0.0.1:11434/api/tags
- prompt: List the local Ollama models available on this machine.

<!-- field: expected | form: "tree" -->
#### Expected

- selected_tool_manifest_contains:
  1. ollama_list_models
- selected_tool_manifest_excludes_mutating_ollama_tools: true
- rest_tool_invoked: true
- rest_operation_called: listOllamaModels
- final_answer_contains_api_model_names: true
- final_answer_excludes_absent_model_names: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc005-ollama-openapi-model-listing](../use-cases/rel03.0-uc005-ollama-openapi-model-listing.md) [S1](../use-cases/rel03.0-uc005-ollama-openapi-model-listing.md#s1)
- [rel03.0-uc005-ollama-openapi-model-listing](../use-cases/rel03.0-uc005-ollama-openapi-model-listing.md) [S2](../use-cases/rel03.0-uc005-ollama-openapi-model-listing.md#s2)
- [rel03.0-uc005-ollama-openapi-model-listing](../use-cases/rel03.0-uc005-ollama-openapi-model-listing.md) [S3](../use-cases/rel03.0-uc005-ollama-openapi-model-listing.md#s3)
- [rel03.0-uc005-ollama-openapi-model-listing](../use-cases/rel03.0-uc005-ollama-openapi-model-listing.md) [S4](../use-cases/rel03.0-uc005-ollama-openapi-model-listing.md#s4)
- [srd009-ollama-adapter](../software-requirements/srd009-ollama-adapter.md) [AC1](../software-requirements/srd009-ollama-adapter.md#ac1)
- [srd024-blocking-input-tools](../software-requirements/srd024-blocking-input-tools.md) [AC1](../software-requirements/srd024-blocking-input-tools.md#ac1)
- [srd024-blocking-input-tools](../software-requirements/srd024-blocking-input-tools.md) [AC2](../software-requirements/srd024-blocking-input-tools.md#ac2)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC2](../software-requirements/srd027-rest-tool-runtime.md#ac2)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC3](../software-requirements/srd027-rest-tool-runtime.md#ac3)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC4](../software-requirements/srd028-rest-client-tools.md#ac4)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC7](../software-requirements/srd028-rest-client-tools.md#ac7)

<!-- case: Profile-owned documentation curator loads generic REST tooling | status: done -->
### Profile-owned documentation curator loads generic REST tooling

Load the documentation-curator profile through the profile-owned
integration suite. Agent Core owns the reusable startup, builtin
registration, and REST runtime behavior; applications/catalog owns the concrete
profile path, selected doc_* REST tools, and REST definitions.

<!-- field: use_case -->
#### Use case

[rel03.0-uc006-documentation-curator](../use-cases/rel03.0-uc006-documentation-curator.md)

<!-- field: go_test -->
#### Go test

mage integration:documentationCurator

<!-- field: inputs | form: "tree" -->
#### Inputs

- profile_owner: applications/catalog
- profile: agents/knowledge-manager/documentation-curator/profile.yaml
- builtin_tool: launch_curator_http
- rest_tools:
  1. doc_list
  1. doc_get
  1. doc_search
  1. doc_validate
  1. doc_suggest_changes
  1. doc_patch_approve
  1. doc_patch_reject

<!-- field: expected | form: "tree" -->
#### Expected

- profile_loads: true
- documentation_factory_registered: true
- rest_tooldefs_selected: true
- core_concrete_profile_required: false

<!-- field: traces -->
#### Traces

- [rel03.0-uc006-documentation-curator](../use-cases/rel03.0-uc006-documentation-curator.md) [S1](../use-cases/rel03.0-uc006-documentation-curator.md#s1)
- [rel03.0-uc006-documentation-curator](../use-cases/rel03.0-uc006-documentation-curator.md) [S4](../use-cases/rel03.0-uc006-documentation-curator.md#s4)
- [srd022-builtin-tools](../software-requirements/srd022-builtin-tools.md) [AC1](../software-requirements/srd022-builtin-tools.md#ac1)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC2](../software-requirements/srd027-rest-tool-runtime.md#ac2)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC1](../software-requirements/srd028-rest-client-tools.md#ac1)

<!-- case: Generic profile server serves docs API and SPA | status: done -->
### Generic profile server serves docs API and SPA

Launch the shipped profile with free loopback addresses and verify the
generic server becomes healthy, then exits through declared lifecycle
transitions. Detailed route behavior is covered by the profile integration.

<!-- field: use_case -->
#### Use case

[rel03.0-uc006-documentation-curator](../use-cases/rel03.0-uc006-documentation-curator.md)

<!-- field: go_test -->
#### Go test

mage integration:documentationCurator

<!-- field: inputs | form: "tree" -->
#### Inputs

- server_package: internal/tools/rest
- docs_fixture: temporary_docs_root
- route_prefix: /api/v1/docs

<!-- field: expected | form: "tree" -->
#### Expected

- docs_list_returned: true
- docs_detail_returned: true
- spa_fallback_returned: true
- bench_route_delegation: false
- concrete_profile_required: false

<!-- field: traces -->
#### Traces

- [rel03.0-uc006-documentation-curator](../use-cases/rel03.0-uc006-documentation-curator.md) [S2](../use-cases/rel03.0-uc006-documentation-curator.md#s2)
- [rel03.0-uc006-documentation-curator](../use-cases/rel03.0-uc006-documentation-curator.md) [S3](../use-cases/rel03.0-uc006-documentation-curator.md#s3)
- [srd020-bench-web-ui](../software-requirements/srd020-bench-web-ui.md) [AC4](../software-requirements/srd020-bench-web-ui.md#ac4)
- [srd024-blocking-input-tools](../software-requirements/srd024-blocking-input-tools.md) [AC1](../software-requirements/srd024-blocking-input-tools.md#ac1)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC2](../software-requirements/srd027-rest-tool-runtime.md#ac2)

<!-- case: Documentation curator exposes validation and patch proposal workflow | status: done -->
### Documentation curator exposes validation and patch proposal workflow

Call the profile-owned generic server APIs for search, validation,
suggestion generation, and patch review. Request MachineSpecs return
structured findings and a persisted pending review. Approval records the review
decision but does not write documentation files.

<!-- field: use_case -->
#### Use case

[rel03.0-uc006-documentation-curator](../use-cases/rel03.0-uc006-documentation-curator.md)

<!-- field: go_test -->
#### Go test

mage integration:documentationCurator

<!-- field: inputs | form: "tree" -->
#### Inputs

- routes:
  1. POST /api/v1/docs/search
  1. POST /api/v1/docs/validate
  1. POST /api/v1/docs/suggestions
  1. POST /api/v1/docs/patches/{patch_id}/approve
- docs_fixture: temporary_docs_root
- strict_validation: true

<!-- field: expected | form: "tree" -->
#### Expected

- search_results_returned: true
- validation_findings_structured: true
- suggestion_status: pending_review
- approval_recorded: true
- documentation_write_applied: false
- concrete_profile_required: false

<!-- field: traces -->
#### Traces

- [rel03.0-uc006-documentation-curator](../use-cases/rel03.0-uc006-documentation-curator.md) [S4](../use-cases/rel03.0-uc006-documentation-curator.md#s4)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC2](../software-requirements/srd027-rest-tool-runtime.md#ac2)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC4](../software-requirements/srd028-rest-client-tools.md#ac4)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC7](../software-requirements/srd028-rest-client-tools.md#ac7)

<!-- case: Profile-owned documentation curator lifecycle exit smoke test | status: done -->
### Profile-owned documentation curator lifecycle exit smoke test

In the profile-owned suite, build or locate the agent binary, start the
external documentation-curator profile, probe documentation health, docs
index, validation, suggestion, and SPA routes, post a generic lifecycle
exit request, and require clean shutdown through the declared Done path.
Agent Core owns the reusable lifecycle-control words and listener
shutdown behavior; applications/catalog owns this concrete tracer bullet.

<!-- field: use_case -->
#### Use case

[rel03.0-uc006-documentation-curator](../use-cases/rel03.0-uc006-documentation-curator.md)

<!-- field: go_test -->
#### Go test

mage integration:documentationCurator

<!-- field: inputs | form: "tree" -->
#### Inputs

- profile_owner: applications/catalog
- profile: agents/knowledge-manager/documentation-curator/profile.yaml
- listener: 127.0.0.1:18081
- frontend_builds:
  1. ../applications/catalog/agents/bench/ui
  1. ../applications/catalog/agents/knowledge-manager/documentation-curator/ui/docs
- bench_route_test: TestDocsAPIRoutesAreNotMounted
- lifecycle_exit_endpoint: POST /api/lifecycle/exit

<!-- field: expected | form: "tree" -->
#### Expected

- documentation_health_ok: true
- documentation_index_returned: true
- validation_endpoint_returned: true
- suggestion_endpoint_returned: true
- spa_route_returned_html: true
- bench_docs_routes_absent: true
- lifecycle_exit_requested: true
- terminal_state_output_succeeded: true
- run_result_output_succeeded: true
- cancelled_output_absent: true
- documentation_listener_released: true
- control_listener_released: true
- documentation_host_stop_go_test: TestKnowledgeManagerConformance
- launch_curator_http_undo: rest_server_stop
- clean_shutdown: true
- core_concrete_profile_required: false

<!-- field: traces -->
#### Traces

- [rel03.0-uc006-documentation-curator](../use-cases/rel03.0-uc006-documentation-curator.md) [S1](../use-cases/rel03.0-uc006-documentation-curator.md#s1)
- [rel03.0-uc006-documentation-curator](../use-cases/rel03.0-uc006-documentation-curator.md) [S2](../use-cases/rel03.0-uc006-documentation-curator.md#s2)
- [rel03.0-uc006-documentation-curator](../use-cases/rel03.0-uc006-documentation-curator.md) [S3](../use-cases/rel03.0-uc006-documentation-curator.md#s3)
- [rel03.0-uc006-documentation-curator](../use-cases/rel03.0-uc006-documentation-curator.md) [S4](../use-cases/rel03.0-uc006-documentation-curator.md#s4)
- [srd020-bench-web-ui](../software-requirements/srd020-bench-web-ui.md) [AC4](../software-requirements/srd020-bench-web-ui.md#ac4)
- [srd024-blocking-input-tools](../software-requirements/srd024-blocking-input-tools.md) [AC1](../software-requirements/srd024-blocking-input-tools.md#ac1)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC10](../software-requirements/srd027-rest-tool-runtime.md#ac10)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC10](../software-requirements/srd029-rest-server-tools.md#ac10)

<!-- case: Profile-owned browser loads machine-request documentation UX | status: done -->
### Profile-owned browser loads machine-request documentation UX

In the profile-owned suite, build or locate the agent binary and
documentation UI assets, start the documentation-curator profile with
machine_request document endpoints, and launch Puppeteer against the local
UX. The page must load from profile-owned UI assets without using bench
routes. Core package tests cover the reusable machine_request and static
asset behavior without naming the concrete profile.

<!-- field: use_case -->
#### Use case

[rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md)

<!-- field: go_test -->
#### Go test

mage integration:documentationCurator

<!-- field: inputs | form: "tree" -->
#### Inputs

- profile_owner: applications/catalog
- profile: agents/knowledge-manager/documentation-curator/profile.yaml
- browser_harness: Puppeteer
- entry_point: Integration.DocumentationCurator
- integration_target: mage integration:documentationCurator
- page_url: http://127.0.0.1:18081/

<!-- field: expected | form: "tree" -->
#### Expected

- page_loaded: true
- profile_ui_assets_served: true
- bench_route_delegation: false
- screenshot_on_failure: true
- console_log_captured_on_failure: true
- core_concrete_profile_required: false

<!-- field: traces -->
#### Traces

- [rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md) [S1](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md#s1)
- [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md) [AC6](../software-requirements/srd031-declarative-document-resources.md#ac6)

<!-- case: Browser verifies document index machine_request sentence | status: done -->
### Browser verifies document index machine_request sentence

Puppeteer waits for the document index network request. The server should
validate the browser request, run one short-lived MachineSpec sentence,
return a document list, and render sidebar categories in the page.

<!-- field: use_case -->
#### Use case

[rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md)

<!-- field: go_test -->
#### Go test

mage integration:documentationCurator

<!-- field: inputs | form: "tree" -->
#### Inputs

- browser_request: GET or POST document index endpoint configured as machine_request
- document_resource: documentation_corpus
- expected_document: SPECIFICATIONS.md

<!-- field: expected | form: "tree" -->
#### Expected

- index_response_ok: true
- sidebar_contains_specifications: true
- request_run_count: 1
- terminal_signal: DocumentIndexReady
- machine_request_trace_required: true
- generic_rest_route_used: documentation_curator.documents
- direct_docs_handler_restored_fails: true
- document_resource_root: docs

<!-- field: traces -->
#### Traces

- [rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md) [S2](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md#s2)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC1](../software-requirements/srd030-machine-request-rest-binding.md#ac1)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC4](../software-requirements/srd030-machine-request-rest-binding.md#ac4)
- [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md) [AC1](../software-requirements/srd031-declarative-document-resources.md#ac1)
- [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md) [AC7](../software-requirements/srd031-declarative-document-resources.md#ac7)

<!-- case: Browser opens raw YAML document through document resource | status: done -->
### Browser opens raw YAML document through document resource

After the sidebar appears, Puppeteer selects a YAML document. The browser
request must run one request machine, pass document resource path policy,
and return raw YAML plus configured parsed content.

<!-- field: use_case -->
#### Use case

[rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md)

<!-- field: go_test -->
#### Go test

mage integration:documentationCurator

<!-- field: inputs | form: "tree" -->
#### Inputs

- selected_document: SPECIFICATIONS.md
- nested_index_document: specs/use-cases/rel03.0-uc007-machine-request-documentation-ux.md
- response_modes:
  1. raw_yaml
  1. parsed_yaml
- blocked_paths:
  1. ../go.mod
  1. cmd/agent/main.go

<!-- field: expected | form: "tree" -->
#### Expected

- raw_yaml_returned: true
- nested_document_path_listed: true
- parsed_content_returned: true
- rendered_title_visible: true
- terminal_signal: DocumentDetailReady
- generic_rest_route_used: documentation_curator.document
- traversal_request_denied: true
- disallowed_extension_denied: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md) [S3](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md#s3)
- [rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md) [S4](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md#s4)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC5](../software-requirements/srd030-machine-request-rest-binding.md#ac5)
- [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md) [AC3](../software-requirements/srd031-declarative-document-resources.md#ac3)
- [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md) [AC4](../software-requirements/srd031-declarative-document-resources.md#ac4)
- [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md) [AC5](../software-requirements/srd031-declarative-document-resources.md#ac5)

<!-- case: Browser verifies rendered HTML and raw YAML view | status: done -->
### Browser verifies rendered HTML and raw YAML view

Puppeteer checks the document viewer after navigation. The rendered HTML
should show parsed document fields, and the raw YAML view should expose
the exact YAML text returned by the document resource response.

<!-- field: use_case -->
#### Use case

[rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md)

<!-- field: go_test -->
#### Go test

mage integration:documentationCurator

<!-- field: inputs | form: "tree" -->
#### Inputs

- selected_document: SPECIFICATIONS.md
- rendered_selectors:
  1. .doc-viewer
  1. .doc-raw-section
- raw_text_contains: "id: agent-core-specifications"

<!-- field: expected | form: "tree" -->
#### Expected

- rendered_html_visible: true
- raw_yaml_visible: true
- parsed_fields_match_raw_document: true
- browser_observed_detail_response_trace: true
- browser_console_errors: 0

<!-- field: traces -->
#### Traces

- [rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md) [S3](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md#s3)
- [rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md) [S5](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md#s5)
- [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md) [AC5](../software-requirements/srd031-declarative-document-resources.md#ac5)
- [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md) [AC6](../software-requirements/srd031-declarative-document-resources.md#ac6)

<!-- case: Browser captures auditable machine_request trace evidence | status: done -->
### Browser captures auditable machine_request trace evidence

The profile-owned browser proof collects trace evidence for the index and
document requests. A failure report must include screenshot path, page
URL, browser console output, network log, request ID, machine name,
terminal signal, and RunResult summary for each machine_request run.

<!-- field: use_case -->
#### Use case

[rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md)

<!-- field: go_test -->
#### Go test

mage integration:documentationCurator

<!-- field: inputs | form: "tree" -->
#### Inputs

- requests:
  1. document_index
  1. document_detail
- failure_artifacts_dir: testdata/integration/rel03-rest-tools/artifacts/uc007

<!-- field: expected | form: "tree" -->
#### Expected

- request_ids_present: true
- machine_names_present: true
- terminal_signals_present: true
- run_result_summaries_present: true
- generic_server_route_present: true
- captured_responses_recorded_on_failure: true
- network_log_captured: true
- screenshot_path_reported_on_failure: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md) [S5](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md#s5)
- [rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md) [S6](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md#s6)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC6](../software-requirements/srd030-machine-request-rest-binding.md#ac6)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC7](../software-requirements/srd030-machine-request-rest-binding.md#ac7)

<!-- case: Machine_request conformance loads configured request machine | status: done -->
### Machine_request conformance loads configured request machine

Run the opt-in conformance tests for machine_request endpoints. The tests
require a configured request-machine fixture, not an injected MachineSpec,
and cover nested catch-all document paths. The generic
internal/tools/rest trace envelope now carries server, route, and machine
identity.

<!-- field: use_case -->
#### Use case

[rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md)

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest -run 'TestRESTServerMachineRequestConformanceLoadsConfiguredMachineFile|TestProfileMachineRequestRunnerUsesDeclaredMachinePolicy|TestProfileMachineRequestRunnerEnforcesDeclaredCommandTimeout|TestProfileMachineRequestRunnerRejectsInvalidConfig'

<!-- field: inputs | form: "tree" -->
#### Inputs

- conformance_env: none
- fixture_owner: agent-core
- machine: generic machine_request fixture
- endpoint_path: /api/v1/docs/{path...}

<!-- field: expected | form: "tree" -->
#### Expected

- configured_request_machine_loaded: true
- injected_machine_spec_required: false
- nested_path_captured: true
- terminal_signal_mapped: DocumentDetailReady
- fixture_declaration_roots_used: true
- request_machine_actions_select_tools: true
- declared_iteration_budget_preserved: true
- declared_command_timeout_preserved_and_enforced: true
- omitted_request_machine_policy_rejected: true
- generic_trace_envelope_go_test: TestRESTServerMachineRequestTraceEnvelope
- generic_trace_contains_server_route_machine: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md) [S2](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md#s2)
- [rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md) [S3](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md#s3)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC1](../software-requirements/srd030-machine-request-rest-binding.md#ac1)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC3](../software-requirements/srd030-machine-request-rest-binding.md#ac3)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC5](../software-requirements/srd030-machine-request-rest-binding.md#ac5)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC6](../software-requirements/srd030-machine-request-rest-binding.md#ac6)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC7](../software-requirements/srd030-machine-request-rest-binding.md#ac7)

<!-- case: Documentation profile routes docs API through generic machine_request | status: done -->
### Documentation profile routes docs API through generic machine_request

Run the profile-owned documentation integration check. The docs API
response trace must identify the generic REST machine_request server
route, request machine, terminal signal, and request-scoped run summary.
Browser evidence for a concrete profile belongs to applications/catalog.

<!-- field: use_case -->
#### Use case

[rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md)

<!-- field: go_test -->
#### Go test

mage integration:documentationCurator

<!-- field: inputs | form: "tree" -->
#### Inputs

- conformance_env: none
- routes:
  1. GET /api/v1/docs
  1. GET /api/v1/docs/{path...}
- request_machine: documentation-curator-request

<!-- field: expected | form: "tree" -->
#### Expected

- generic_rest_route_used: true
- docs_specific_runner_used: false
- route_trace_present: true
- machine_name_present: true
- terminal_signal_present: true
- run_summary_present: true
- current_conformance_gate: none
- always_on_go_test: TestStandaloneServerConformanceUsesRESTMachineRequestRoutes
- browser_profile_tracer_owner: applications/catalog

<!-- field: traces -->
#### Traces

- [rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md) [S2](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md#s2)
- [rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md) [S5](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md#s5)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC1](../software-requirements/srd030-machine-request-rest-binding.md#ac1)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC2](../software-requirements/srd030-machine-request-rest-binding.md#ac2)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC3](../software-requirements/srd030-machine-request-rest-binding.md#ac3)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC5](../software-requirements/srd030-machine-request-rest-binding.md#ac5)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC7](../software-requirements/srd030-machine-request-rest-binding.md#ac7)
- [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md) [AC2](../software-requirements/srd031-declarative-document-resources.md#ac2)

<!-- case: Documentation browser actions use configured request-machine routes | status: done -->
### Documentation browser actions use configured request-machine routes

Post browser-shaped validation, suggestion, approval, and rejection
requests through the profile-launched generic REST server. Each route
returns configured request-machine trace metadata, and filesystem review
artifacts retain suggestions across short-lived action machines. Unknown
patch IDs route to the configured 404 response, while invalid review
transitions return a domain rejection.

<!-- field: use_case -->
#### Use case

[rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md)

<!-- field: go_test -->
#### Go test

mage integration:documentationCurator

<!-- field: inputs | form: "tree" -->
#### Inputs

- configured_routes:
  1. POST /api/v1/actions/validate
  1. POST /api/v1/actions/suggest
  1. POST /api/v1/actions/patches/{patch_id}/approve
  1. POST /api/v1/actions/patches/{patch_id}/reject
- removed_route: POST /api/v1/actions
- request_machine: documentation-curator-request

<!-- field: expected | form: "tree" -->
#### Expected

- private_action_registry_used: false
- direct_command_execute_used: false
- legacy_action_envelope_status: 404
- validation_terminal_signal: RESTResponded
- suggestion_terminal_signal: RESTAccepted
- approval_terminal_signal: RESTResponded
- rejection_terminal_signal: RESTResponded
- missing_patch_terminal_signal: RESTMissing
- missing_patch_status: 404
- missing_patch_error: patch_missing
- patch_store_survives_action_runs: true
- action_trace_present: true
- profile_wiring_test: TestDocumentationActionEndpointsUseMachineRequestBinding
- profile_ux_test: TestDocumentationUXActionsUseRequestMachineRoutes
- browser_validation_trace: validate_action

<!-- field: traces -->
#### Traces

- [rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md) [S7](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md#s7)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC1](../software-requirements/srd030-machine-request-rest-binding.md#ac1)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC2](../software-requirements/srd030-machine-request-rest-binding.md#ac2)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC3](../software-requirements/srd030-machine-request-rest-binding.md#ac3)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC7](../software-requirements/srd030-machine-request-rest-binding.md#ac7)

<!-- case: REST observability redacts secrets and reports failure stages | status: done -->
### REST observability redacts secrets and reports failure stages

Exercise client and server error paths with synthetic secret material.
Traces and Result output report identity, status, timing data, byte
counts, failure stage, and resulting signal. Configured secrets stay out.

<!-- field: use_case -->
#### Use case

[rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md)

<!-- field: go_test -->
#### Go test

TestRESTTools_TracingRedactionAndErrors

<!-- field: inputs | form: "tree" -->
#### Inputs

- failures:
  1. auth_resolution
  1. network_io
  1. timeout
  1. response_mapping
  1. route_validation
  1. shutdown_failed
- redaction_selectors:
  1. headers.authorization
  1. query.token
  1. body.secret

<!-- field: expected | form: "tree" -->
#### Expected

- traces_include_operation_or_route: true
- traces_include_status_and_timing: true
- failure_stage_reported: true
- secrets_redacted: true
- verbose_body_capture_requires_config: true

<!-- field: traces -->
#### Traces

- [rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md) [S2](../use-cases/rel03.0-uc001-rest-client-sync.md#s2)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC7](../software-requirements/srd027-rest-tool-runtime.md#ac7)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC8](../software-requirements/srd028-rest-client-tools.md#ac8)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC9](../software-requirements/srd028-rest-client-tools.md#ac9)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC9](../software-requirements/srd029-rest-server-tools.md#ac9)

<!-- case: REST server control and queue policies are enforced; unsupported read policies rejected | status: done -->
### REST server control and queue policies are enforced; unsupported read policies rejected

Validate modeled lifecycle_control enqueueing and endpoint queue overrides,
and that read_policy validation supports only first_available (the default):
other read policies such as priority_order and round_robin are rejected at
validation, not honored.

<!-- field: use_case -->
#### Use case

[rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md)

<!-- field: go_test -->
#### Go test

TestRESTServer_ControlQueueAndReadPolicyConformance

<!-- field: inputs | form: "tree" -->
#### Inputs

- endpoint_bindings:
  1. lifecycle_control
  1. emit_signal
- queue_policies:
  1. reject
  1. drop_oldest
  1. drop_newest
- supported_read_policies:
  1. first_available
- rejected_read_policies:
  1. priority_order
  1. round_robin

<!-- field: expected | form: "tree" -->
#### Expected

- lifecycle_handlers_execute_directly: false
- lifecycle_events_enqueued: true
- configured_queue_policy_honored: true
- unsupported_queue_drain_and_read_policies_rejected: true

<!-- field: traces -->
#### Traces

- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC10](../software-requirements/srd027-rest-tool-runtime.md#ac10)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC5](../software-requirements/srd029-rest-server-tools.md#ac5)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC7](../software-requirements/srd029-rest-server-tools.md#ac7)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC10](../software-requirements/srd029-rest-server-tools.md#ac10)

<!-- case: REST client safety and async selectors are enforced | status: done -->
### REST client safety and async selectors are enforced

Validate CIDR allowlists, response schema checks, domain_error_code
output, async correlation selection, idempotency metadata on outbound
requests, and rejection of legacy await_operation with guidance to
declare REST probe and lifecycle delay states in MachineSpec.

<!-- field: use_case -->
#### Use case

[rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md)

<!-- field: go_test -->
#### Go test

TestRESTClient_CIDRAllowlistPolicy, TestRESTClient_ResponseSchemaAndDomainErrorOutput, TestRESTClient_AsyncCorrelationAndIdempotencyHeader, TestRESTClient_AwaitOperationReferenceValidation

<!-- field: inputs | form: "tree" -->
#### Inputs

- evidence: TestRESTClient_SafetyAndAsyncConformance
- network_policy: cidr_allowlist
- response_schema: configured
- async_selectors:
  1. request_id
  1. correlation

<!-- field: expected | form: "tree" -->
#### Expected

- cidr_allowlist_enforced: true
- invalid_response_schema_success_signal: false
- domain_error_code_in_output: true
- await_operation_rejected: true

<!-- field: traces -->
#### Traces

- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC6](../software-requirements/srd027-rest-tool-runtime.md#ac6)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC5](../software-requirements/srd028-rest-client-tools.md#ac5)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC7](../software-requirements/srd028-rest-client-tools.md#ac7)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC8](../software-requirements/srd028-rest-client-tools.md#ac8)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC9](../software-requirements/srd028-rest-client-tools.md#ac9)

<!-- case: OpenAPI server binding preserves configured endpoint behavior | status: done -->
### OpenAPI server binding preserves configured endpoint behavior

Bind OpenAPI operations to endpoints that already declare response
mapping, machine_request, monitor view, queue, signal, and
allowed_signals. Import should add transport and schema data without
dropping configured endpoint behavior.

<!-- field: use_case -->
#### Use case

[rel03.0-uc004-openapi-import](../use-cases/rel03.0-uc004-openapi-import.md)

<!-- field: go_test -->
#### Go test

TestRESTOpenAPI_ServerBindPreservesEndpointConfig

<!-- field: inputs | form: "tree" -->
#### Inputs

- evidence: TestRESTOpenAPI_ServerBindPreservesEndpointConfig
- endpoint_fields:
  1. response
  1. machine_request
  1. monitor_view
  1. queue
  1. signal
  1. allowed_signals

<!-- field: expected | form: "tree" -->
#### Expected

- configured_fields_preserved: true
- openapi_method_path_schema_imported: true
- incompatible_operation_ids_rejected: true

<!-- field: traces -->
#### Traces

- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC3](../software-requirements/srd027-rest-tool-runtime.md#ac3)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC8](../software-requirements/srd029-rest-server-tools.md#ac8)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC3](../software-requirements/srd030-machine-request-rest-binding.md#ac3)

<!-- case: REST client compensation creates and executes rollback mementos | status: done -->
### REST client compensation creates and executes rollback mementos

Distinguish read-only noop undo from mutating REST operations that
declare configured compensation. The mutating path captures request and
response metadata needed by rollback. The configured compensation path
executes from the captured receipt.

<!-- field: use_case -->
#### Use case

[rel03.0-uc001-rest-client-sync](../use-cases/rel03.0-uc001-rest-client-sync.md)

<!-- field: go_test -->
#### Go test

TestRESTClient_CompensationUndoMemento, TestRESTClient_CompensationExecutorRunsFromReceipt

<!-- field: inputs | form: "tree" -->
#### Inputs

- evidence: TestRESTClient_CompensationUndoMemento
- execution_evidence: TestRESTClient_CompensationExecutorRunsFromReceipt
- operations:
  1. rest_client_get
  1. rest_client_create
  1. rest_client_set
  1. rest_client_delete

<!-- field: expected | form: "tree" -->
#### Expected

- read_undo_noop: true
- mutating_undo_noop: false
- compensation_payload_valid: true
- compensation_execution_succeeds: true
- metadata_contains_operation_resource_and_idempotency: true

<!-- field: traces -->
#### Traces

- [srd025-rollback-lifecycle](../software-requirements/srd025-rollback-lifecycle.md) [AC5](../software-requirements/srd025-rollback-lifecycle.md#ac5)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md) [AC8](../software-requirements/srd027-rest-tool-runtime.md#ac8)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC3](../software-requirements/srd028-rest-client-tools.md#ac3)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC10](../software-requirements/srd028-rest-client-tools.md#ac10)

<!-- case: Machine_request validates terminal response schemas | status: done -->
### Machine_request validates terminal response schemas

Configure a machine_request terminal response schema and return output
that violates it. The binding should reject the response or the public
config schema field should be removed.

<!-- field: use_case -->
#### Use case

[rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md)

<!-- field: go_test -->
#### Go test

TestRESTServerMachineRequestTerminalResponseSchema

<!-- field: inputs | form: "tree" -->
#### Inputs

- evidence: TestRESTServerMachineRequestTerminalResponseSchema
- terminal_signal: DocumentationReady
- response_schema: configured

<!-- field: expected | form: "tree" -->
#### Expected

- invalid_response_schema_success: false
- failure_stage: response_invalid
- diagnostic_redacted: true

<!-- field: traces -->
#### Traces

- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC5](../software-requirements/srd030-machine-request-rest-binding.md#ac5)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md) [AC6](../software-requirements/srd030-machine-request-rest-binding.md#ac6)

<!-- case: REST document_resources config is loaded or rejected | status: done -->
### REST document_resources config is loaded or rejected

Load a REST definition that declares top-level document_resources and a
machine_request document_resources reference. The runtime should either
make those resources available through trusted config or reject the fields
with a clear validation error.

<!-- field: use_case -->
#### Use case

[rel03.0-uc007-machine-request-documentation-ux](../use-cases/rel03.0-uc007-machine-request-documentation-ux.md)

<!-- field: go_test -->
#### Go test

TestRESTDocumentResourcesConfigConformance

<!-- field: inputs | form: "tree" -->
#### Inputs

- evidence: TestRESTDocumentResourcesConfigConformance
- document_resource: documentation_corpus
- request_machine: generic machine_request fixture

<!-- field: expected | form: "tree" -->
#### Expected

- unknown_document_resources_silently_ignored: false
- current_profile_resource_words_still_work: true
- validation_error_names_field_when_unsupported: true

<!-- field: traces -->
#### Traces

- [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md) [AC1](../software-requirements/srd031-declarative-document-resources.md#ac1)
- [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md) [AC3](../software-requirements/srd031-declarative-document-resources.md#ac3)
- [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md) [AC7](../software-requirements/srd031-declarative-document-resources.md#ac7)

<!-- case: static_assets serves an embedded bundle and its declared config | status: done -->
### static_assets serves an embedded bundle and its declared config

Launch a REST server whose static_assets endpoint selects the observer
bundle compiled into agent-core, with SPA fallback and a declared config
map. The bundle index answers the mount and unknown routes, and
ui-config.json returns the declared map. Validation rejects a definition
with both root and bundle, with neither, or with an unknown bundle name.

<!-- field: use_case -->
#### Use case

[rel03.0-uc003-rest-server-signals](../use-cases/rel03.0-uc003-rest-server-signals.md)

<!-- field: go_test -->
#### Go test

TestStaticAssets_BundleServesEmbeddedIndexWithSPAFallback, TestStaticAssets_ServesDeclaredConfigAsUIConfigJSON, TestStaticAssets_UIConfigWithoutDeclarationFallsThrough, TestValidateDefinition_staticAssetsRejectsInvalidConfigs

<!-- field: inputs | form: "tree" -->
#### Inputs

- bundle: observer
- config:
  - title: Chatbot Mesh observer
  - fleet_endpoint: /monitor/fleet

<!-- field: expected | form: "tree" -->
#### Expected

- bundle_index_served: true
- spa_fallback_to_bundle_index: true
- ui_config_json_equals_declared_map: true
- root_and_bundle_together_rejected: true
- unknown_bundle_rejected_with_known_names: true

<!-- field: traces -->
#### Traces

- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md) [AC5](../software-requirements/srd029-rest-server-tools.md#ac5)
