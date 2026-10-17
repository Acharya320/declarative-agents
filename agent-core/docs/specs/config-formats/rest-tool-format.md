---
id: cfg-rest-tool-format
title: REST Tool Configuration Format
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: purpose -->
## Purpose

Defines the YAML format for declarative REST definitions under a top-level
rest mapping. Agent profiles and tool declaration files use this
format to configure REST client and REST server boundary words.

The format keeps HTTP authority in trusted configuration. ToolDef config
references named clients, servers, resources, operations, and endpoints.
Runtime input supplies only declared request values.

<!-- field: release_03_0_support -->
## Release 03 0 support

<!-- field: implemented | form: "tree" -->
### Implemented

1. rest.version, clients, servers, openapi, auth, limits, retry_policies, and response_mappings.
1. REST client resource operations and RPC-shaped operations.
1. REST client async send and await by request_id.
1. REST server emit_signal, emit_dynamic_signal, read_state, invoke_handler, stream_events, health, static_metadata, machine_request, signal_source, and monitor endpoints.
1. OpenAPI import from local files into client operations and server endpoints.
1. network.cidrs parsing and CIDR enforcement.
1. async correlation selectors and wire-level idempotency metadata; readiness polling is authored as explicit REST probe and lifecycle delay states in MachineSpec.
1. response_mapping.schema validation before client success signals.
1. OpenAPI-bound endpoint preservation for configured endpoint behavior.
1. REST client compensation mementos and rollback-time execution for mutating operations with configured compensation.
1. machine_request terminal response schema validation before success responses.
1. - "lifecycle_control endpoints": ValidateDefinition accepts only fixed enqueue_signal and allowlisted inject_signal modes; both enqueue events for MachineSpec routing rather than running lifecycle operations directly (TestRESTServer_LifecycleControlEnqueuesSignals).
1. - "canonical lifecycle-exit injection": agent-core synthesizes a POST /api/lifecycle/exit lifecycle_control endpoint (action enqueue_signal, signal ExitRequested) for every served agent at launch, so a profile need not declare it. Injection is skipped when the server opts out (lifecycle_exit.disabled), already declares the canonical path, or already uses the reserved route name, and never mutates the resolved endpoint map. A declared auth_ref is enforced before enqueue; auth:none is loopback-only even when listener policy permits a public bind (TestRESTServer_InjectsCanonicalLifecycleExit, TestInjectedLifecycleExitEnforcesDeclaredBearerAuth). Hand-declared emit_signal control routes remain supported.
1. - "declared_machines monitor view": a read_state endpoint returns the loaded root and distinct machine_request MachineSpecs as declarations, with authored presentation metadata and no request-machine run state.

<!-- field: partial -->
### Partial

- rest_await_event read_policy values outside the implemented first-available behavior are rejected during ToolDef factory validation.

<!-- field: reserved -->
### Reserved

- document_resources top-level config and machine_request.document_resources are typed target-format fields that ValidateDefinition rejects until generic REST document resource loading is implemented. Current request machines use profile-selected filesystem resource ToolDefs.

<!-- field: timeout_envelope -->
## Timeout envelope

<!-- field: authority -->
### Authority

A selected REST client operation derives its finite authority from the
referenced client limit profile and its async timeout where applicable. A
machine_request endpoint derives its outer request authority from
machine_request.timeout and separately inspects the referenced request
profile or machine. A signal_source endpoint derives its finite host-run
authority from signal_source.timeout. REST await words derive their
authority from the selected ToolDef config and referenced queue or source
timeout.

<!-- field: environment_expansion -->
### Environment expansion

REST definitions and ToolDef declarations expand ${VAR} and
${VAR:-default} before timeout-envelope validation. The MachineSpec that
dispatches the selected REST word supplies the finite maximum: every
expanded operation timeout must remain strictly below
budget.command_timeout. Equal, malformed, non-positive, or larger expanded
values fail startup before network or queue dispatch. Deployment schemas
may enforce tighter bounds but cannot replace this runtime comparison.

<!-- field: retry_aggregate -->
### Retry aggregate

A client with retry_ref derives one conservative dispatch authority from
its selected request timeout and retry policy. For A attempts with
per-attempt timeout T, the aggregate is A*T plus the delay before each
retry: sum(delay(i), i=1..A-1). attempts is a positive integer.
A missing retry_ref target, malformed duration, non-positive attempt count,
or duration-arithmetic overflow fails startup.

none contributes zero delay. fixed and the legacy empty backoff contribute
min(initial_delay, max_delay) before every retry when max_delay is positive.
exponential contributes initial_delay, then doubles for each later retry
with overflow-safe saturation at a positive max_delay when present. The
inspector uses checked arithmetic rather than treating time.Duration
saturation as a finite authority.

The aggregate applies to synchronous client words and rest_client_send.
rest_client_await reads previously stored async state in a later dispatch,
so its async.timeout remains a separate authority and does not acquire the
submission retry aggregate.

<!-- field: retry_evidence -->
### Retry evidence

Deterministic inspection evidence names the client, retry_ref, attempts,
normalized backoff, aggregate duration, and enclosing command_timeout.
The ordinary client limit entries remain visible as per-attempt evidence.

<!-- field: closure_resolution -->
### Closure resolution

Validation follows profile REST definitions, declaration includes and
profile-local override precedence, machine_request profile or machine
references, and compatibility wrappers. Diagnostics name the selected
action, REST authority path, operation duration, and enclosing machine
command timeout.

<!-- field: schema -->
## Schema

<!-- field: description -->
### Description

A REST definition file can live beside an agent profile, a tool declaration
file, or a future REST-specific config directory. Each file contains one
rest mapping. The runtime merges selected files through profile-owned
declaration paths.

<!-- field: fields -->
### Fields

<!-- field_entry: rest -->
#### rest

<!-- field: type -->
##### Type

object

<!-- field: required | form: "scalar" -->
##### Required

true

<!-- field: description -->
##### Description

Root mapping for all REST definitions in the file.

<!-- field: sub_fields | form: "typed table" -->
##### Sub fields

| field | type | required | description |
| --- | --- | --- | --- |
| version | string | true | Format version. Initial value is "v1". |
| clients | map[string]rest_client | false | Named outbound REST clients. |
| servers | map[string]rest_server | false | Named inbound REST servers. |
| openapi | map[string]openapi_import | false | Named OpenAPI 3 imports. |
| auth | map[string]auth_profile | false | Named authentication profiles referenced by clients, operations, or servers. |
| limits | map[string]limit_profile | false | Named timeout, size, redirect, and network limit profiles. |
| retry_policies | map[string]retry_policy | false | Named retry policies for outbound operations. |
| response_mappings | map[string]response_mapping | false | Reusable response mapping profiles. |
| document_resources | map[string]document_resource | false | Reserved target-format field for named filesystem-backed document resources. ValidateDefinition rejects this field. Knowledge Manager behavior uses profile-selected filesystem resource ToolDefs instead. |

<!-- field: rest_client -->
### REST client

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| base_url | string | false | Base URL for relative operation paths. Runtime input cannot replace this value. |
| auth_ref | string | false | Name of an auth profile. |
| limits_ref | string | false | Name of a limit profile. |
| retry_ref | string | false | Default retry policy for operations on this client. |
| resources | map[string]resource | false | Resource-oriented operation groups. |
| operations | map[string]operation | false | RPC-shaped or standalone outbound operations. |

<!-- field: resource -->
### Resource

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| path | string | true | Resource path template relative to the client base_url. |
| operations | map[string]resource_operation | true | Allowed resource verbs. Valid names are get, set, create, and delete. |
| id_field | string | false | Reserved field rejected at load until resource-level response defaults are implemented. |
| version_field | string | false | Reserved field rejected at load until resource-level response defaults are implemented. |

<!-- field: operation -->
### Operation

<!-- field: fields -->
#### Fields

<!-- field_entry_deep: openapi_operation_id -->
##### openapi_operation_id

<!-- field: type -->
###### Type

string

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

Imported OpenAPI operation ID used to populate method, path, and schemas.

<!-- field_entry_deep: method -->
##### method

<!-- field: type -->
###### Type

string

<!-- field: required -->
###### Required

unless_openapi_operation_id

<!-- field: description -->
###### Description

HTTP method fixed by configuration or imported from OpenAPI.

<!-- field_entry_deep: path -->
##### path

<!-- field: type -->
###### Type

string

<!-- field: required -->
###### Required

unless_openapi_operation_id_or_resource_path

<!-- field: description -->
###### Description

Path template. Resource operations may inherit the resource path. Absolute URLs are rejected unless network policy allows them.

<!-- field_entry_deep: base_url_source -->
##### base_url_source

<!-- field: type -->
###### Type

string

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

Optional command_state value that lets this trusted operation select
its base URL from a labeled prior command output. Omission preserves
the client's configured base_url.

<!-- field_entry_deep: base_url_selector -->
##### base_url_selector

<!-- field: type -->
###### Type

string

<!-- field: required -->
###### Required

when_base_url_source_is_command_state

<!-- field: description -->
###### Description

A $from(<label>).<path> selector resolving to an absolute HTTP or
HTTPS URL. This is operation configuration, not a runtime parameter;
input_mapping remains forbidden from setting transport authority.
Mutually exclusive with base_url_host_selector.

<!-- field_entry_deep: base_url_host_selector -->
##### base_url_host_selector

<!-- field: type -->
###### Type

string

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

A $from(<label>).<path> selector resolving to a bare host or IP, for a
target a prior step discovered as an address rather than a whole URL
(for example a pod address). The resolved value is composed with
base_url_scheme and base_url_port. A value carrying a scheme,
credential, port, path, query, fragment, whitespace, or control
character is rejected before network I/O. Mutually exclusive with
base_url_selector; requires base_url_source command_state.

<!-- field_entry_deep: base_url_scheme -->
##### base_url_scheme

<!-- field: type -->
###### Type

string

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

Scheme for a host composed through base_url_host_selector. Accepts
http or https and defaults to http. Rejected without
base_url_host_selector.

<!-- field_entry_deep: base_url_port -->
##### base_url_port

<!-- field: type -->
###### Type

string

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

Literal port for a host composed through base_url_host_selector. Must
be 1-65535 and is omitted when unset, leaving the scheme default.
Mutually exclusive with base_url_port_selector; rejected without
base_url_host_selector.

<!-- field_entry_deep: base_url_port_selector -->
##### base_url_port_selector

<!-- field: type -->
###### Type

string

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

A $from(<label>).<path> selector resolving the composed port per item,
for a fleet whose targets do not share one port (for example each
agent's own monitor port, published beside its address). Accepts a
string or a whole JSON number in 1-65535; another type, a fractional
or out-of-range value is rejected before network I/O. The composed
authority must still satisfy the client's limits.network.ports
allowlist, so configuration bounds the reachable ports and the
selector only chooses among them. Mutually exclusive with
base_url_port; rejected without base_url_host_selector.

<!-- field_entry_deep: allow_auth_on_selected_authority -->
##### allow_auth_on_selected_authority

<!-- field: type -->
###### Type

bool

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

Explicitly permits configured credentials on an authority selected
through base_url_source. The default is false: authenticated
operations reject a selected authority that differs from the
configured client authority before network I/O.

<!-- field_entry_deep: params -->
##### params

<!-- field: type -->
###### Type

request_binding

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

Declared path, query, header, and body inputs.

<!-- field_entry_deep: body -->
##### body

<!-- field: type -->
###### Type

object

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

Request body template built from declared params.

<!-- field_entry_deep: success -->
##### success

<!-- field: type -->
###### Type

status_mapping

<!-- field: required | form: "scalar" -->
###### Required

true

<!-- field: description -->
###### Description

Successful HTTP status mapping and emitted signal.

<!-- field_entry_deep: failures -->
##### failures

<!-- field: type -->
###### Type

[]status_mapping

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

Domain failure status mappings.

<!-- field_entry_deep: response_ref -->
##### response_ref

<!-- field: type -->
###### Type

string

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

Reusable response mapping profile.

<!-- field_entry_deep: response -->
##### response

<!-- field: type -->
###### Type

response_mapping

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

Inline response mapping.

<!-- field_entry_deep: side_effects -->
##### side_effects

<!-- field: type -->
###### Type

[]object

<!-- field: required | form: "scalar" -->
###### Required

true

<!-- field: description -->
###### Description

ToolDef-style side-effect declarations for the operation.

<!-- field_entry_deep: reversibility -->
##### reversibility

<!-- field: type -->
###### Type

object

<!-- field: required | form: "scalar" -->
###### Required

true

<!-- field: description -->
###### Description

Reversible, compensatable, or irreversible classification.

<!-- field_entry_deep: compensation -->
##### compensation

<!-- field: type -->
###### Type

object

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

Configured compensating operation or cancellation path.

<!-- field_entry_deep: async -->
##### async

<!-- field: type -->
###### Type

async_client_config

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

Enables rest_client_send and rest_client_await for this operation.

<!-- field: resource_operation -->
### Resource operation

<!-- field: alias_of -->
#### Alias of

operation

<!-- field: description -->
#### Description

A resource_operation uses the same fields as operation. It appears only
under resource.operations, where the operation name is get, set, create, or
delete. Resource operations omit path when the resource path is
sufficient, and set path when a verb needs a relative subpath.

<!-- field: status_mapping -->
### Status mapping

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| status | "[]int" | true | HTTP status codes that select this mapping. |
| signal | string | true | Result signal emitted when one of the status codes matches. |
| domain_error_code | string | false | Optional service-domain error code placed in Result output. |
| response_ref | string | false | Name of a reusable response_mapping applied for this status. |
| response | response_mapping | false | Inline response mapping applied for this status. |

<!-- field: auth_profile -->
### Auth profile

<!-- field: description -->
#### Description

Defines authentication material by reference. Secret values live outside
the REST config file. Release 03.0 supports none, basic, bearer,
header_token, and query_token.

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| type | string | true | One of none, basic, bearer, header_token, or query_token. |
| username_ref | string | by_type | Required for basic. Names the credential containing the username. |
| password_ref | string | by_type | Required for basic. Names the credential containing the password. |
| token_ref | string | by_type | Required for bearer, header_token, and query_token. |
| header | string | by_type | Required for header_token. Names the outbound header. |
| query | string | by_type | Required for query_token. Names the outbound query parameter. |
| scheme | string | false | Optional bearer scheme override. Default is Bearer. |

<!-- field: limit_profile -->
### Limit profile

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| timeout | duration string | false | Overall request or route handling timeout. |
| connect_timeout | duration string | false | Outbound connection timeout. |
| read_timeout | duration string | false | Maximum duration for reading a response or inbound body. |
| max_request_bytes | int | false | Maximum rendered outbound request body or accepted inbound body size. |
| max_response_bytes | int | false | Maximum response body bytes read before CommandError. |
| max_header_bytes | int | false | Maximum inbound header bytes for server listeners. |
| redirect | redirect_policy | false | Outbound redirect behavior. |
| network | network_policy | false | Allowed schemes, hosts, CIDRs, ports, and listener bindings. |

<!-- field: redirect_policy -->
### Redirect policy

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| mode | string | true | One of none, same_host, or allowlist. |
| allow_hosts | "[]string" | false | Host allowlist used when mode is allowlist. |
| max_redirects | int | false | Maximum redirects followed by the client. |

<!-- field: network_policy -->
### Network policy

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| schemes | "[]string" | false | Allowed URL schemes. |
| hosts | "[]string" | false | Allowed destination or listener hostnames. |
| cidrs | "[]string" | false | Allowed destination or listener CIDR ranges. Parsed and enforced by the network policy. |
| ports | "[]int" | false | Allowed destination or listener ports. |
| allow_public_listener | bool | false | Permits server addresses such as 0.0.0.0 when true. |

<!-- field: retry_policy -->
### Retry policy

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| attempts | int | true | Strictly positive maximum attempts including the first try. |
| backoff | string | false | One of none (zero delay), fixed (initial_delay each retry), or exponential (doubling initial_delay). |
| initial_delay | duration string | false | Delay before the first retry. |
| max_delay | duration string | false | Non-negative maximum delay between retries; exponential growth saturates at this value without overflow. |
| retry_status | "[]int" | false | HTTP status codes that permit retry. |
| retry_network_errors | bool | false | Allows retry for temporary network I/O failures. |
| require_idempotency | bool | false | Requires idempotency metadata before retrying mutating operations. |

<!-- field: rest_server -->
### REST server

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| address | string | true | Listener address. Unsafe public binding is rejected unless policy allows it. |
| limits_ref | string | false | Name of a limit profile. |
| queue | queue_config | false | Default inbound queue behavior. |
| event_sources | map[string]event_source | false | Target-format source filters for rest_await_event. Current ToolDef configs can provide sources directly; server-level event_sources need follow-up validation and wiring. |
| endpoints | map[string]endpoint | true | Configured inbound routes. |
| shutdown | shutdown_config | false | Graceful shutdown timeout and drain behavior. |
| lifecycle_exit | lifecycle_exit_injection_config | false | Tunes agent-core's synthesis of the canonical lifecycle-control exit endpoint for this server. Injection is on by default; omit this field to accept it. Set disabled to opt a server out (a pure fixture, or a monitor server that must carry observability only). |

<!-- field: lifecycle_exit_injection_config -->
### Lifecycle exit injection config

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| disabled | bool | false | When true, agent-core injects no canonical exit endpoint for this server. Default false injects POST /api/lifecycle/exit. |
| auth_ref | string | false | Optional auth profile reference enforced before the injected endpoint validates or enqueues a request. Empty is an explicit loopback-only mode; non-loopback lifecycle control requires a declared auth profile. |

<!-- field: endpoint -->
### Endpoint

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| openapi_operation_id | string | false | Imported OpenAPI operation ID used to populate method, path, and schemas. |
| method | string | unless_openapi_operation_id | HTTP method accepted by the endpoint or imported from OpenAPI. |
| path | string | unless_openapi_operation_id | Route path template or imported OpenAPI path. |
| binding | string | true | Endpoint binding kind from srd029-rest-server-tools R5, srd030-machine-request-rest-binding, or srd046-request-driven-signal-source. signal_source is independent from machine_request and drives the injected host-profile SignalSource without an inbound queue. |
| monitor_view | string | false | Selects a cached monitor view for read_state, static_metadata, and stream_events endpoints, or a loaded declaration view for read_state. One of machine_spec, declared_machines, current_state, tools, metrics, events, openapi, or command_state (srd033-monitor-rest-api R3, R7, R8). declared_machines returns the root plus distinct machine_request machines from the trusted profile closure. command_state additionally requires labels. |
| labels | "[]string" | when_monitor_view_command_state | Non-empty allowlist of authored transition labels or command names whose redaction-cleared declared output a command_state endpoint exposes through the srd038 command-state store. Required and non-empty when monitor_view is command_state; rejected on any other view (srd033-monitor-rest-api R7.1, R7.6). |
| signal | string | false | Static signal emitted by emit_signal endpoints. |
| allowed_signals | "[]string" | false | Signal allowlist for emit_dynamic_signal endpoints. |
| signal_field | string | false | Optional dotted validated-payload path whose string value selects an emit_dynamic_signal outcome. |
| signal_mapping | map[string]string | false | Optional profile-owned mapping from signal_field values to allowed_signals; requires signal_field and every mapped result must be allowlisted. |
| lifecycle_control | lifecycle_control_config | false | Lifecycle control request mapping. ValidateDefinition accepts fixed enqueue_signal and allowlisted inject_signal modes plus an optional target selector; the endpoint enqueues the lifecycle request event for MachineSpec routing rather than running a lifecycle operation directly. agent-core injects a canonical lifecycle_control exit endpoint per served agent (see rest_server.lifecycle_exit), so most profiles need not declare one; hand-declared emit_signal control endpoints remain a supported alternative. |
| machine_request | machine_request_config | false | Request-to-machine mapping used when binding is machine_request. |
| signal_source | signal_source_config | false | Trusted request-to-SignalEnvelope mapping used only when binding is signal_source. |
| request | request_binding | false | Declared path, query, header, and body inputs. |
| response | response_mapping | false | HTTP response mapping after validation or event enqueue. |
| queue | queue_config | false | Endpoint-specific queue behavior. |

<!-- field: event_source -->
### Event source

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| server | string | false | Server name for this source. Defaults to the containing server when nested under rest.servers. |
| routes | "[]string" | false | Route or endpoint names accepted by this source. |
| signals | "[]string" | false | Event signals accepted by this source. |
| queue | string | false | Queue name to read when it differs from the server default. |
| timeout | duration string | false | Source-specific await timeout. |
| stopped_behavior | string | false | One of ignore, emit_server_stopped, or command_error. |

<!-- field: rest_await_event_config -->
### REST await event config

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| sources | "[]event_source" | true | Ordered or fair-read set of REST server event sources. |
| allowed_signals | "[]string" | false | Final signal allowlist applied after source selection. |
| timeout | duration string | false | Maximum time to wait before emitting AwaitTimedOut. |
| read_policy | string | false | Optional read policy. Omit this field for the implemented first-available behavior. priority_order and round_robin are rejected until implemented. |
| stopped_behavior | string | false | Default stopped-source behavior when a source does not override it. |

<!-- field: lifecycle_control_config -->
### Lifecycle control config

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| action | string | true | One of enqueue_signal or inject_signal. enqueue_signal requires a configured signal; inject_signal reads an allowlisted signal from the request. |
| signal | string | false | Static lifecycle signal enqueued for MachineSpec routing. |
| allowed_signals | "[]string" | false | Allowlist required for inject_signal endpoints. |
| target_schema | json_schema | false | Request body schema for rollback targets, exit reason, injected signal payload, or operator metadata. |
| require_auth_ref | string | false | Auth profile enforced before validation or enqueue. Omit only for loopback-only lifecycle control. |

<!-- field: machine_request_config -->
### Machine request config

<!-- field: note -->
#### Note

REST definitions load with strict field checking; any field not listed here is rejected at load rather than silently ignored (GH-486). Fields that were previously documented but never implemented (tools, tool_declarations, initial_state, error_responses, retain_trace, trace) are reserved for a future release and are not accepted today.

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| profile | string | true | Trusted profile or profile-equivalent config used for the short-lived request machine run. |
| machine | string | false | Optional MachineSpec override inside the trusted profile scope. |
| initial_signal | string | false | Optional initial signal seeded into the request machine's first transition. Defaults to Seed. |
| request | machine_request_input_mapping | true | Declared mapping from validated HTTP request fields to RequestEvent input. |
| response | machine_request_response_mapping | true | Mapping from terminal machine output to HTTP status, headers, body, and content type. |
| timeout | duration string | false | Maximum duration for the whole request machine run. |
| document_resources | "[]string" | false | Reserved target-format references to document resources. ValidateDefinition rejects this field. Current request machines select filesystem resource words through trusted profile ToolDefs. |

<!-- field: machine_request_input_mapping -->
### Machine request input mapping

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| path | map[string]string | false | Selector mapping from route path variables to RequestEvent fields. |
| query | map[string]string | false | Selector mapping from query parameters to RequestEvent fields. |
| headers | map[string]string | false | Selector mapping from accepted headers to RequestEvent fields. |
| body | map[string]string | false | Selector mapping from validated body fields to RequestEvent fields. |
| sensitive | "[]string" | false | Mapped request field names whose values core removes from the synthetic seed entry before live selection or persistence. Every name must be a key declared by path, query, headers, or body; metadata contains names only, never request values. |

<!-- field: machine_request_response_mapping -->
### Machine request response mapping

<!-- field: note -->
#### Note

default_signal, response_word, and missing_response_stage were previously documented but never implemented; they are reserved and rejected at load (GH-486).

<!-- field: fields -->
#### Fields

<!-- field_entry_deep: terminal_signals -->
##### terminal_signals

<!-- field: type -->
###### Type

map[string]machine_request_terminal_response

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

HTTP response mapping by allowed terminal signal. Either this or terminal_states is required.

<!-- field_entry_deep: terminal_states -->
##### terminal_states

<!-- field: type -->
###### Type

map[string]machine_request_terminal_response

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

HTTP response mapping by terminal machine state. Resolved before
terminal_signals when both are present, because many signals may reach
one terminal state and one signal may reach many. A request machine
built from exec words emits only ToolDone and ToolFailed, so state
keying is what lets it return a distinct response per outcome
(GH-615). Each key must be a terminal state of the request machine, or
the endpoint fails configuration validation.

<!-- field: machine_request_terminal_response -->
### Machine request terminal response

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| status | int | true | HTTP status code written for this terminal signal. |
| headers | map[string]string | false | Header values or selectors copied from response word output. |
| content_type | string | false | Response content type. Default is application/json. |
| body | map[string]string | false | Selector mapping from response word output to HTTP response body fields. |
| schema | json_schema | false | Response body schema that terminal output should satisfy. Validation is enforced before machine_request writes a success response. |

<!-- field: signal_source_config -->
### Signal source config

<!-- field: note -->
#### Note

REST definitions load this mapping with strict field checking. Runtime input supplies only values selected by these trusted selectors and cannot provide a raw signal or program authority.

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| source | string | true | Configured source identity copied to SignalEnvelope.source and bounded trace attributes. |
| discriminator_field | string | true | Declared $.body, $.query, $.path, or $.headers selector whose string value is looked up only in signal_mapping. |
| signal_mapping | map[string]string | true | Non-empty closed discriminator mapping to trusted signal names. Unknown keys return refused_undeclared without invoking the host runner. |
| run_id_field | string | true | Declared request selector yielding the non-empty host run id. |
| expected_state_field | string | false | Optional declared request selector yielding the expected-state admission guard. |
| payload | map[string]string | true | Non-empty output-field to declared-request-selector mapping used to construct SignalEnvelope.payload. |
| sensitive | "[]string" | false | Payload output fields converted to typed core.OutputRedactionPath values before admission. Every entry must name a payload mapping key. |
| timeout | duration string | true | Positive maximum time for admission and the accepted host run. |
| responses | signal_source_response_mappings | true | HTTP policy selected after source validation or admission; it does not rewrite accepted when the later run fails. |

<!-- field: signal_source_response_mappings -->
### Signal source response mappings

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| accepted | signal_source_response | true | Response for an accepted admission whose run did not fail. |
| refused_undeclared | signal_source_response | true | Response for an unknown discriminator or core refused_undeclared outcome. |
| refused_conflict | signal_source_response | false | Response for refused_conflict. Status defaults to HTTP 409 when omitted. |
| source_validation | signal_source_response | true | Response for transport, schema, mapping, or caller-authority validation failure. The runner is not called. |
| machine_run_failed | signal_source_response | true | Response policy for runner failure, timeout, cancellation, or an accepted admission whose run later failed. |

<!-- field: signal_source_response -->
### Signal source response

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| status | int | true | HTTP status in the range 100-599. |
| include_diagnostic | bool | false | Permits a redacted stage or diagnostic in addition to the fixed outcome fields. |
| include_output | bool | false | Permits the redacted terminal run summary in addition to the fixed outcome fields. |

<!-- field: document_resource -->
### Document resource

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| root | string | true | Trusted filesystem root for the corpus, relative to the profile directory or agent working directory. |
| include | "[]string" | false | Glob-style include patterns for list and get operations. |
| extensions | "[]string" | true | Allowed file extensions such as .yaml, .yml, or .md. |
| response_modes | "[]string" | true | Allowed modes. Values are raw_yaml, parsed_yaml, raw_markdown, and parsed_markdown. |
| default_response_mode | string | false | Response mode used when an operation does not override it. |
| category_rules | "[]document_category_rule" | false | Path prefix rules used to derive UI or index categories. |
| max_bytes | int | false | Maximum file size returned by one get operation. |
| symlinks | string | false | Symlink policy. Default is deny_escape. |
| binary_policy | string | false | Binary file policy. Default is deny. |
| operations | map[string]document_resource_operation | false | Operation-specific response modes, signals, and output mapping. |
| ui | document_resource_ui | false | Human-facing labels, grouping, and presentation hints for profile-owned UI configuration. |

<!-- field: document_category_rule -->
### Document category rule

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| prefix | string | true | Resource-relative path prefix. |
| category | string | true | Category label assigned to matching documents. |

<!-- field: document_resource_operation -->
### Document resource operation

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| type | string | true | One of list, get, search, validate, or suggest. |
| response_mode | string | false | Operation-specific response mode chosen from the resource response_modes list. |
| success_signal | string | true | Signal emitted when the operation succeeds. |
| not_found_signal | string | false | Signal emitted when a requested document is missing. |
| denied_signal | string | false | Signal emitted when path policy rejects the selector. |
| output | map[string]string | false | Field selectors copied into Result output or HTTP response data. |

<!-- field: document_resource_ui -->
### Document resource UI

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| label | string | false | Human-facing corpus label. |
| sidebar_grouping | string | false | Grouping strategy such as category, prefix, or flat. |
| actions | "[]string" | false | Profile-local UI action names that may use this resource. |
| asset_ref | string | false | Profile-owned UI asset reference used to present the resource. |

<!-- field: shutdown_config -->
### Shutdown config

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| timeout | duration string | false | Graceful shutdown timeout. |
| drain_policy | string | false | Optional drain_then_stop value naming the implemented close-then-drain behavior. The inert drain spelling and target-format reject_new, drop_queued, and fail_queued values are rejected. |
| drain_timeout | duration string | false | Reserved target-format field. Release 03.0 rejects this field until timed queue and stream draining is implemented. |
| stop_listeners | bool | false | Optional true value documenting the current listener shutdown behavior. Explicit false is rejected because rest_server_stop closes owned listeners. |
| queue_on_shutdown | string | false | Reserved target-format field. Release 03.0 rejects keep, drop, and fail until distinct queue shutdown behavior is implemented. |
| unblock_await_signal | string | false | Signal emitted to a waiting word when the server stops. Defaults to ServerStopped; custom values are included in ToolDef emits validation. |

<!-- field: request_binding -->
### Request binding

<!-- field: fields -->
#### Fields

<!-- field_entry_deep: path -->
##### path

<!-- field: type -->
###### Type

map[string]json_schema

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

Path variables accepted from runtime input.

<!-- field_entry_deep: query -->
##### query

<!-- field: type -->
###### Type

map[string]json_schema

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

Query parameters accepted from runtime input.

<!-- field_entry_deep: headers -->
##### headers

<!-- field: type -->
###### Type

map[string]json_schema

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

Header parameters accepted from runtime input.

<!-- field_entry_deep: body_schema -->
##### body_schema

<!-- field: type -->
###### Type

json_schema

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

Schema for runtime body input.

<!-- field_entry_deep: body_source -->
##### body_source

<!-- field: type -->
###### Type

string

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

Selects where declared runtime params come from. One of params
(default, explicit tool input), previous_result (select declared
params from the prior Result output through input_mapping),
command_state (select declared params from a labeled prior step in the
srd038 command-state store through input_mapping $from(...) selectors),
or none. body_source command_state is rejected only when no
command-state store is configured.

<!-- field_entry_deep: input_mapping -->
##### input_mapping

<!-- field: type -->
###### Type

map[string]string

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

Map of declared-param-name to a selector. Under body_source
previous_result the selector is a $.-style path into the previous
Result output (for example embedding: $.mapped.embedding,
id: $.body.id, document: $.carried.doc). Under body_source
command_state the selector takes the $from(<label>).<dotted-path> form,
naming a prior step by command name or label (for example
embedding: $from(embed_query).mapped.embedding). After selection only
declared params reach the request builder, so the source output shape
does not violate the declared-only contract. Every target is a
declared path, query, header, or body param; selectors do not name
method, URL, auth, redirect, or host.

<!-- field_entry_deep: carry_forward -->
##### carry_forward

<!-- field: type -->
###### Type

[]string

<!-- field: required | form: "scalar" -->
###### Required

false

<!-- field: description -->
###### Description

Declared param names copied into this operation's Result output under a
carried key, so a later word can select them (for example
$.carried.document). Carry-forward threads data past an intermediate
word without a shared command-state store. Every entry is a
declared param and does not name transport authority.

<!-- field: response_mapping -->
### Response mapping

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| schema | json_schema | false | Expected response body schema. |
| output | map[string]string | false | Field selectors copied into structured Result output. |
| redact | "[]string" | false | Selectors removed from output, traces, logs, and errors. |
| resource_id | string | false | Selector for created or updated resource identity. |
| request_id | string | false | Selector for async request identity. |

<!-- field: openapi_import -->
### OpenAPI import

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| path | string | true | Local OpenAPI 3 file path or trusted URL. |
| base_url | string | false | Override for OpenAPI server URL. |
| expose | "[]string" | false | Operation IDs imported for REST client operations. |
| bind | map[string]string | false | Operation ID to server endpoint name mapping. |
| side_effects | map[string][]object | false | Side-effect metadata by operation ID. |
| reversibility | map[string]object | false | Reversibility metadata by operation ID. |

<!-- field: async_client_config -->
### Async client config

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| request_id | string | true | Selector or template that yields the async request ID. |
| correlation | string | false | Selector or template for correlation data. |
| idempotency_token | string | false | Template used for safe retries or duplicate suppression. |
| await_operation | string | false | Reserved legacy field. Validation rejects it; profile authors declare one-shot REST probe and lifecycle delay states in MachineSpec. rest_client_await only reads stored request completion by request_id or correlation token. |
| timeout | duration string | false | Maximum await duration. |
| state_retention | string | false | One of consume, keep, or expire. |

<!-- field: queue_config -->
### Queue config

<!-- field: fields | form: "typed table" -->
#### Fields

| field | type | required | description |
| --- | --- | --- | --- |
| name | string | false | Queue or inbox name included in each accepted InboundEvent as queue. |
| capacity | int | false | Maximum queued inbound events. |
| overflow | string | false | One of reject, drop_oldest, or drop_newest. |
| timeout | duration string | false | Default await timeout. |
| payload_shape | json_schema | false | JSON schema validated against the complete event payload before redaction and enqueue. |

<!-- field: json_schema -->
### JSON schema

<!-- field: description -->
#### Description

Embedded JSON Schema-compatible content. The value is a YAML mapping that
follows JSON Schema object structure. The runtime preserves it as raw
schema data for validation and for OpenAPI schema interop.

<!-- field: allowed_fields -->
#### Allowed fields

- type
- properties
- required
- items
- enum
- format
- minimum
- maximum
- min_length
- max_length
- pattern
- additional_properties
- description

<!-- field: fragment_files | form: "tree" -->
## Fragment files

- description:
  A REST definition file MAY be a fragment (srd052-declaration-fragments):
  `params` beside the `rest` root declares typed parameters, and an
  importing REST file's `expand` list applies the fragment with
  arguments. The two fields sit beside `unit` and `imports` (srd050) and
  follow the same rules as in cfg-tool-declaration-format: scalar-only
  substitution of `$param(name)`, and the expanded definition compiled
  through ValidateDefinition like a hand-written one. When `as` is given,
  every name the fragment's `rest` body produces -- the clients, servers,
  openapi, auth, limits, retry_policies, response_mappings, and
  document_resources entries, and the operation and endpoint names nested
  under clients and servers, which the closure holds unique -- is renamed
  `<as>_<name>`, and the fragment's own references to those names
  (`limits_ref`, `auth_ref`, `retry_ref`, `response_ref`,
  `require_auth_ref`, `document_resources`, `bind`) follow, so each
  expansion resolves on its own. A tool that reaches a produced endpoint
  names the prefixed form, `routes: [alpha_exit]`. A reference to a name the fragment does
  not declare is left for the importer's closure to supply. Landed with
  GH-2082.
- example:
  # units/http-json-api.yaml -- the fragment
  unit: http-json-api
  params:
  - name: service
    type: string
  - name: base_url
    type: string
  - name: timeout
    type: string
    default: 5s
  rest:
    version: v1
    clients:
      $param(service)_api:
        base_url: $param(base_url)
        limits: default
    limits:
      default:
        connect_timeout: $param(timeout)

  # agents/observer/rest.yaml -- the importer
  unit: observer-rest
  expand:
  - fragment: ../units/http-json-api.yaml
    as: metrics
    args: {service: metrics, base_url: http://metrics:9090}
  - fragment: ../units/http-json-api.yaml
    as: traces
    args: {service: traces, base_url: http://traces:4318, timeout: 2s}
  rest:
    version: v1
  # produces clients metrics_metrics_api and traces_traces_api, each with
  # its own limits profile metrics_default and traces_default

<!-- field: library_imports -->
## Library imports

REST `imports` and `expand` fragment paths follow the import path forms
of cfg-tool-declaration-format `library_imports` (srd056-unit-libraries R1):
relative to the declaring file, or absolute under `/opt/agent-core` or a root
the profile declares. An `openapi.path` keeps its existing base, the REST file
that declares it (srd050 R4.5); a library-rooted REST unit's OpenAPI document
therefore travels with the unit. A REST file under a library root cannot
import outside its root except by another root's prefix.

<!-- field: service_libraries -->
## Service libraries

An agent that serves a REST surface can publish a service library
(srd056-unit-libraries R4). Inside one application it is a directory beside
the agents, holding

  lib/<service>/
    rest.yaml     # unit with the client operations for the served surface
    types.yaml    # type unit for the request and response shapes
    words.yaml    # the client words consumers share

Consumers declare the application's library directory as a root on their
profile and import `/opt/<root>/<service>/rest.yaml`, and none of them
redeclares the client. Per-consumer values such as base URLs stay in the
consumer, because a library is read-only. A service library earns its place
with two or more importing profiles, a fragment instantiation counting as
one; with one, the client stays in its only consumer.

A contract served by more than one application lives in agent-core instead:
its type unit under tools/units and its client fragment under
tools/rest/units, reached through `/opt/agent-core` with no declared root.
The first instance is the retrieval query client
(srd060-retrieval-client-contract): consumers expand
`/opt/agent-core/tools/rest/units/retrieval-query-client-fragment.yaml` with
their route prefix and selectors, and import
`/opt/agent-core/tools/units/types-retrieval.yaml` for QueryResult.

<!-- field: format_examples | form: "tree" -->
## Format examples

- hand_authored:
  rest:
    version: v1
    auth:
      github_app:
        type: bearer
        token_ref: github_token
    limits:
      public_api:
        timeout: 30s
        max_request_bytes: 65536
        max_response_bytes: 1048576
        redirect:
          mode: same_host
    retry_policies:
      safe_reads:
        attempts: 3
        backoff: exponential
        retry_status: [429, 500, 502, 503]
    clients:
      github:
        base_url: https://api.github.com
        auth_ref: github_app
        limits_ref: public_api
        resources:
          issue:
            path: /repos/{owner}/{repo}/issues/{number}
            operations:
              get:
                method: GET
                params:
                  path:
                    owner: {type: string}
                    repo: {type: string}
                    number: {type: integer}
                success: {status: [200], signal: RESTResourceRead}
                failures:
                  - status: [404]
                    signal: RESTMissing
                response:
                  output:
                    title: $.title
                    state: $.state
                side_effects:
                  - kind: external_api
                    target: github.issue
                    state: read_only
                reversibility:
                  classification: reversible
                  undo: noop
              set:
                method: PATCH
                params:
                  path:
                    owner: {type: string}
                    repo: {type: string}
                    number: {type: integer}
                  body_schema:
                    type: object
                    properties:
                      title: {type: string}
                      state: {type: string}
                success: {status: [200], signal: RESTResourceWritten}
                side_effects:
                  - kind: external_api
                    target: github.issue
                    state: issue_updated
                reversibility:
                  classification: compensatable
                  undo: restore_previous_issue_fields
        operations:
          search_issues:
            method: GET
            path: /search/issues
            params:
              query:
                q: {type: string}
            success: {status: [200], signal: RESTResponded}
            side_effects:
              - kind: external_api
                target: github.search
                state: read_only
            reversibility:
              classification: reversible
              undo: noop
- openapi_backed:
  rest:
    version: v1
    openapi:
      payments:
        path: ./openapi/payments.yaml
        base_url: https://payments.internal
        expose:
          - getPayment
          - createPayment
          - cancelPayment
        bind:
          receivePaymentWebhook: payment_webhook
        side_effects:
          createPayment:
            - kind: external_api
              target: payments.payment
              state: payment_created
        reversibility:
          createPayment:
            classification: compensatable
            undo: cancelPayment
    clients:
      payments:
        base_url: https://payments.internal
        operations:
          create_payment:
            openapi_operation_id: createPayment
            success: {status: [202], signal: RESTAccepted}
            async:
              request_id: $.id
              correlation: $.correlation_id
              idempotency_token: "{{ params.order_id }}"
              timeout: 2m
              state_retention: consume
            side_effects:
              - kind: external_api
                target: payments.payment
                state: payment_submitted
            reversibility:
              classification: compensatable
              undo: cancelPayment
    servers:
      payments_webhooks:
        address: 127.0.0.1:9090
        queue:
          name: payments
          capacity: 128
          overflow: reject
        endpoints:
          payment_webhook:
            openapi_operation_id: receivePaymentWebhook
            binding: emit_signal
            signal: PaymentWebhookReceived
            queue:
              name: payments
            response:
              output:
                accepted: "true"
- previous_result_threading:
  # An embedding word runs before a Chroma query word. The embed word carries
  # the query text forward and yields a flat vector under mapped.embedding; the
  # Chroma word selects that vector from the previous Result and renders it as
  # a nested array through the single-token raw passthrough.
  rest:
    version: v1
    clients:
      ollama:
        base_url: http://127.0.0.1:11434
        operations:
          embed_query:
            method: POST
            path: /api/embed
            params:
              body_schema:
                type: object
                required: [model, input]
                properties:
                  model: {type: string}
                  input: {type: string}
              carry_forward: [input]
            success: {status: [200], signal: RESTResponded}
            response:
              output:
                embedding: $.embedding
            side_effects:
              - kind: external_api
                target: ollama.embed
                state: read_only
            reversibility:
              classification: reversible
              undo: noop
      chroma:
        base_url: http://127.0.0.1:8000
        operations:
          query_collection:
            method: POST
            path: /api/v2/tenants/default/databases/default/collections/{collection}/query
            params:
              path:
                collection: {type: string}
              body_schema:
                type: object
                required: [query_embeddings]
                properties:
                  query_embeddings: {type: array}
                  n_results: {type: integer}
              body_source: previous_result
              input_mapping:
                query_embeddings: $.mapped.embedding
                collection: $.carried.input
            body:
              query_embeddings: ["{{ params.query_embeddings }}"]
              n_results: 5
            success: {status: [200], signal: RESTResponded}
            side_effects:
              - kind: external_api
                target: chroma.collection
                state: read_only
            reversibility:
              classification: reversible
              undo: noop
- command_state_threading:
  # A three-word chain. embed_query (word 1) yields a flat vector under
  # embedding and carries its input text forward. record_event (word 2) runs
  # next. query_collection (word 3) reaches back past word 2 to word 1 by
  # label through body_source command_state and $from(embed_query) selectors,
  # so a non-adjacent output threads into its request body as a nested array.
  rest:
    version: v1
    clients:
      ollama:
        base_url: http://127.0.0.1:11434
        operations:
          embed_query:
            method: POST
            path: /api/embed
            params:
              body_schema:
                type: object
                required: [model, input]
                properties:
                  model: {type: string}
                  input: {type: string}
              carry_forward: [input]
            success: {status: [200], signal: RESTResponded}
            response:
              output:
                embedding: $.embedding
            side_effects:
              - kind: external_api
                target: ollama.embed
                state: read_only
            reversibility:
              classification: reversible
              undo: noop
      audit:
        base_url: http://127.0.0.1:9000
        operations:
          record_event:
            method: POST
            path: /api/events
            params:
              body_schema:
                type: object
                required: [event]
                properties:
                  event: {type: string}
            body:
              event: query_prepared
            success: {status: [200], signal: RESTResponded}
            side_effects:
              - kind: external_api
                target: audit.events
                state: append_only
            reversibility:
              classification: reversible
              undo: noop
      chroma:
        base_url: http://127.0.0.1:8000
        operations:
          query_collection:
            method: POST
            path: /api/v2/tenants/default/databases/default/collections/{collection}/query
            params:
              path:
                collection: {type: string}
              body_schema:
                type: object
                required: [query_embeddings]
                properties:
                  query_embeddings: {type: array}
                  n_results: {type: integer}
              body_source: command_state
              input_mapping:
                query_embeddings: $from(embed_query).mapped.embedding
                collection: $from(embed_query).carried.input
            body:
              query_embeddings: ["{{ params.query_embeddings }}"]
              n_results: 5
            success: {status: [200], signal: RESTResponded}
            side_effects:
              - kind: external_api
                target: chroma.collection
                state: read_only
            reversibility:
              classification: reversible
              undo: noop
- lifecycle_control_fan_in:
  rest:
    version: v1
    servers:
      agent_control:
        address: 127.0.0.1:19090
        queue:
          name: control
          capacity: 64
          overflow: reject
        event_sources:
          control_events:
            routes: [exit, pause, rollback, inject]
            signals: [AgentExitRequested, PauseRequested, RollbackRequested, InjectSignalRequested]
            stopped_behavior: emit_server_stopped
        endpoints:
          exit:
            method: POST
            path: /api/lifecycle/exit
            binding: lifecycle_control
            lifecycle_control:
              action: enqueue_signal
              signal: AgentExitRequested
              target_schema:
                type: object
                properties:
                  reason: {type: string}
                  status: {type: string}
          pause:
            method: POST
            path: /api/lifecycle/pause
            binding: lifecycle_control
            lifecycle_control:
              action: enqueue_signal
              signal: PauseRequested
          rollback:
            method: POST
            path: /api/lifecycle/rollback
            binding: lifecycle_control
            lifecycle_control:
              action: enqueue_signal
              signal: RollbackRequested
              target_schema:
                type: object
                properties:
                  target_iteration: {type: integer}
          inject:
            method: POST
            path: /api/machine/inject
            binding: lifecycle_control
            lifecycle_control:
              action: inject_signal
              allowed_signals: [Approved, Rejected, PauseRequested]
- tool_declarations:
  tools:
    - name: get_issue
      type: builtin
      init: rest_client_get
      category: boundary
      emits: [RESTResourceRead, RESTMissing, RESTDomainFailed, CommandError]
      config:
        rest_ref: github
        resource: issue
        operation: get

    - name: create_payment
      type: builtin
      init: rest_client_send
      category: boundary
      emits: [RESTAccepted, CommandError]
      config:
        rest_ref: payments
        operation: create_payment

    - name: await_payment
      type: builtin
      init: rest_client_await
      category: boundary
      emits: [RESTResponded, RESTDomainFailed, RESTMissing, RESTAwaitTimedOut, CommandError]
      config:
        rest_ref: payments
        operation: create_payment

    - name: launch_payment_webhooks
      type: builtin
      init: rest_server_launch
      category: boundary
      emits: [ServerLaunched, CommandError]
      config:
        rest_ref: payments_webhooks

    - name: await_any_control_event
      type: builtin
      init: rest_await_event
      category: boundary
      emits: [AgentExitRequested, PauseRequested, RollbackRequested, InjectSignalRequested, AwaitTimedOut, ServerStopped, CommandError]
      config:
        sources:
          - server: agent_control
            routes: [exit, pause, rollback, inject]
            signals: [AgentExitRequested, PauseRequested, RollbackRequested, InjectSignalRequested]
        timeout: 30s
        read_policy: first_available
        stopped_behavior: emit_server_stopped

    - name: exit_agent
      type: builtin
      init: exit_agent
      category: boundary
      emits: [AgentExited, CommandError]
      config:
        drain_policy: stop_servers
- reserved_document_resource_target_format:
  # Reserved target-format sketch. ValidateDefinition rejects
  # rest.document_resources and machine_request.document_resources in release
  # 03.0. Use profile-selected filesystem resource ToolDefs for current
  # Knowledge Manager request machines.
  rest:
    version: v1
    document_resources:
      documentation_corpus:
        root: docs
        include: ["**/*.yaml", "**/*.yml", "**/*.md"]
        extensions: [.yaml, .yml, .md]
        response_modes: [raw_yaml, parsed_yaml, raw_markdown]
        default_response_mode: parsed_yaml
        max_bytes: 262144
        symlinks: deny_escape
        binary_policy: deny
        category_rules:
          - prefix: specs/software-requirements/
            category: srd
          - prefix: specs/semantic-models/
            category: semantic-model
          - prefix: specs/config-formats/
            category: config-format
        operations:
          list:
            type: list
            success_signal: DocumentListReady
          get:
            type: get
            success_signal: DocumentReady
            not_found_signal: DocumentMissing
            denied_signal: DocumentResourceDenied
        ui:
          label: Documentation
          sidebar_grouping: category
          actions: [doc_list, doc_get, doc_search]
          asset_ref: ${AGENT_CATALOG_ROOT}/knowledge-manager/documentation-curator/ui
    servers:
      knowledge_manager:
        address: 127.0.0.1:19191
        endpoints:
          get_documentation:
            method: POST
            path: /api/docs/query
            binding: machine_request
            request:
              body_schema:
                type: object
                required: [query]
                properties:
                  query: {type: string}
            machine_request:
              profile: ${AGENT_CATALOG_ROOT}/knowledge-manager/documentation-curator/profile.yaml
              document_resources: [documentation_corpus]
              request:
                body:
                  query: $.query
              response:
                terminal_signals:
                  DocumentationReady:
                    status: 200
                    content_type: application/json
                    body:
                      title: $.title
                      body: $.body
                  DocumentationMissing:
                    status: 404
                    body:
                      error: $.message
              timeout: 10s

<!-- field: validation_rules -->
## Validation rules

<!-- validation_rule: V1 -->
### V1

<!-- field: rule -->
#### Rule

A REST config file MUST contain rest.version.

<!-- validation_rule: V2 -->
### V2

<!-- field: rule -->
#### Rule

Client operation names MUST be unique after OpenAPI imports and hand-authored definitions are merged.

<!-- validation_rule: V3 -->
### V3

<!-- field: rule -->
#### Rule

Server endpoint names MUST be unique after OpenAPI bind mappings are applied.

<!-- validation_rule: V4 -->
### V4

<!-- field: rule -->
#### Rule

Runtime input fields MUST be declared under request_binding before use in path, query, header, or body templates.

<!-- validation_rule: V5 -->
### V5

<!-- field: rule -->
#### Rule

Resource operations MUST use only get, set, create, or delete names.

<!-- validation_rule: V6 -->
### V6

<!-- field: rule -->
#### Rule

RPC-shaped operations MUST live under operations rather than resource operation maps.

<!-- validation_rule: V7 -->
### V7

<!-- field: rule -->
#### Rule

Mutating operations MUST declare side_effects and reversibility.

<!-- validation_rule: V8 -->
### V8

<!-- field: rule -->
#### Rule

Async operations MUST declare request_id and timeout when selected by rest_client_send.

<!-- validation_rule: V9 -->
### V9

<!-- field: rule -->
#### Rule

Server endpoints with emit_dynamic_signal MUST declare allowed_signals.

<!-- validation_rule: V10 -->
### V10

<!-- field: rule -->
#### Rule

Server endpoints that bind public addresses MUST declare an explicit network policy.

<!-- validation_rule: V11 -->
### V11

<!-- field: rule -->
#### Rule

Auth type values MUST be one of none, basic, bearer, header_token, or query_token.

<!-- validation_rule: V12 -->
### V12

<!-- field: rule -->
#### Rule

OpenAPI operation IDs in expose and bind MUST exist exactly once in the imported document.

<!-- validation_rule: V13 -->
### V13

<!-- field: rule -->
#### Rule

Redirect mode MUST be none, same_host, or allowlist.

<!-- validation_rule: V14 -->
### V14

<!-- field: rule -->
#### Rule

Redaction selectors MUST be valid for the configured response or request mapping syntax.

<!-- validation_rule: V15 -->
### V15

<!-- field: rule -->
#### Rule

rest_await_event configs MUST declare at least one source.

<!-- validation_rule: V16 -->
### V16

<!-- field: rule -->
#### Rule

Event source stopped_behavior MUST be one of ignore, emit_server_stopped, or command_error.

<!-- validation_rule: V17 -->
### V17

<!-- field: rule -->
#### Rule

lifecycle_control enqueue_signal endpoints MUST declare signal; inject_signal endpoints MUST declare allowed_signals.

<!-- validation_rule: V18 -->
### V18

<!-- field: rule -->
#### Rule

lifecycle_control endpoints MUST enqueue declared signals and MUST NOT execute lifecycle operations inside HTTP handlers.

<!-- validation_rule: V19 -->
### V19

<!-- field: rule -->
#### Rule

machine_request endpoints MUST declare machine_request.profile, request mapping, response.terminal_signals, and timeout or inherited route timeout.

<!-- validation_rule: V20 -->
### V20

<!-- field: rule -->
#### Rule

machine_request endpoints MUST NOT let runtime HTTP input choose profile, machine, tools, tool_declarations, credentials, URLs, or workspace roots.

<!-- validation_rule: V21 -->
### V21

<!-- field: rule -->
#### Rule

machine_request responses MUST be selected from configured terminal signals or response word output and MUST NOT be assembled by direct Go demo handlers.

<!-- validation_rule: V22 -->
### V22

<!-- field: rule -->
#### Rule

Release 03.0 MUST reject rest.document_resources and machine_request.document_resources as reserved target-format fields until generic REST document resource loading is implemented.

<!-- validation_rule: V23 -->
### V23

<!-- field: rule -->
#### Rule

Current request machines MUST use profile-selected filesystem resource ToolDefs, not machine_request.document_resources, to keep roots, include patterns, extension allowlists, response modes, symlink policy, and binary policy under trusted profile control.

<!-- validation_rule: V24 -->
### V24

<!-- field: rule -->
#### Rule

Profile-selected filesystem resource ToolDefs MUST deny empty paths, absolute paths, parent-directory traversal, and resolved paths outside the configured root.

<!-- validation_rule: V25 -->
### V25

<!-- field: rule -->
#### Rule

machine_request profile-backed tool selection MUST use trusted profile declaration roots and MUST fail when any request MachineSpec action is missing from the selected ToolDefs.

<!-- validation_rule: V26 -->
### V26

<!-- field: rule -->
#### Rule

machine_request response trace output MUST identify generic route identity and run outcome when exposed.

<!-- validation_rule: V27 -->
### V27

<!-- field: rule -->
#### Rule

Server shutdown config MUST reject unsupported drain_policy values and reserved shutdown fields before launch.

<!-- validation_rule: V28 -->
### V28

<!-- field: rule -->
#### Rule

input_mapping and carry_forward entries MUST target params declared under request_binding path, query, headers, or body_schema; undeclared targets are rejected.

<!-- validation_rule: V29 -->
### V29

<!-- field: rule -->
#### Rule

input_mapping MUST only be set when body_source is previous_result or command_state, and body_source MUST be one of params, previous_result, command_state, or none; command_state is rejected only when no command-state store is configured.

<!-- validation_rule: V30 -->
### V30

<!-- field: rule -->
#### Rule

input_mapping and carry_forward MUST NOT name method, URL, authentication, redirect policy, or destination host, so runtime input cannot acquire transport authority.

<!-- validation_rule: V31 -->
### V31

<!-- field: rule -->
#### Rule

A body value equal to a single "{{ params.NAME }}" token MUST render the param's raw typed value, so a flat vector param placed inside a nested array element renders as a nested array rather than a stringified scalar.

<!-- validation_rule: V32 -->
### V32

<!-- field: rule -->
#### Rule

$from(<label>).<path> selectors are valid only under body_source command_state, and $.-style selectors are valid only under body_source previous_result; an input_mapping selector whose form does not match its body_source is rejected.

<!-- validation_rule: V33 -->
### V33

<!-- field: rule -->
#### Rule

Hand-authored REST YAML MUST reject fields unknown to the typed definition model at every nesting level.

<!-- validation_rule: V34 -->
### V34

<!-- field: rule -->
#### Rule

Inbound JSON bodies MUST contain exactly one object and satisfy declared string, number, integral integer, object, array, boolean, or null field types.

<!-- validation_rule: V35 -->
### V35

<!-- field: rule -->
#### Rule

$.path and $from(label).path selectors MUST use the shared dotted-path grammar; empty components, whitespace or control characters, malformed parentheses, and empty labels or paths are invalid.

<!-- validation_rule: V36 -->
### V36

<!-- field: rule -->
#### Rule

base_url_source, when present, MUST be command_state and MUST carry
either a base_url_selector or a base_url_host_selector in
$from(<label>).<path> form, never both. A base_url_selector value MUST be
an absolute HTTP or HTTPS URL without user information. A
base_url_host_selector value MUST be a bare host or IP, and is composed
with base_url_scheme (http or https, default http) and a port from either
base_url_port (a 1-65535 literal) or base_url_port_selector (resolved per
item as a string or whole JSON number in 1-65535), never both.
base_url_scheme, base_url_port, and base_url_port_selector are invalid
without base_url_host_selector. Either form MUST pass the operation's
network, redirect, and credential-scope policies before network I/O, so a
resolved port outside the client's ports allowlist is refused.
Runtime input_mapping authority fields remain forbidden.

<!-- validation_rule: V37 -->
### V37

<!-- field: rule -->
#### Rule

A non-empty client retry_ref MUST resolve to one retry policy whose
attempts is positive. Timeout-envelope inspection MUST add the selected
request timeout once per attempt plus every fixed or capped exponential
delay before attempts two through N using checked duration arithmetic.
Missing references, malformed delays, non-positive attempts, and
aggregate overflow MUST fail before dispatch. rest_client_await keeps its
separate async timeout and MUST NOT inherit the submission aggregate.

<!-- validation_rule: V38 -->
### V38

<!-- field: rule -->
#### Rule

signal_source endpoints MUST declare a bounded source identity, declared discriminator and run-id selectors, a non-empty closed signal mapping, a non-empty payload mapping, a positive timeout, and statuses for accepted, refused_undeclared, source_validation, and machine_run_failed; refused_conflict defaults to HTTP 409.

<!-- validation_rule: V39 -->
### V39

<!-- field: rule -->
#### Rule

signal_source selectors MUST reference fields declared by endpoint.request and MUST NOT select signal, profile, machine, tools, tool declarations, model configuration, checkpoint, URL, method, host, or other caller program authority.

<!-- validation_rule: V40 -->
### V40

<!-- field: rule -->
#### Rule

signal_source sensitive entries MUST name payload mapping keys and become typed one-segment Result output-redaction paths before the injected host runner is called.

<!-- validation_rule: V41 -->
### V41

<!-- field: rule -->
#### Rule

signal_source validates and maps the complete request before invoking its injected SignalSource runner; unknown discriminators and invalid requests do not invoke the runner, enter emit_signal queues, or use machine_request.

<!-- validation_rule: V42 -->
### V42

<!-- field: rule -->
#### Rule

monitor_view declared_machines is valid only on read_state and exposes the root plus distinct machine_request declarations from the trusted profile closure without request-machine run state.

<!-- field: invalid_examples | form: "tree" -->
## Invalid examples

- ambiguous_resource_operation:
  # Resource operations cannot define RPC-style verbs.
  rest:
    version: v1
    clients:
      crm:
        base_url: https://crm.example
        resources:
          account:
            path: /accounts/{id}
            operations:
              approve:
                method: POST
                path: /accounts/{id}/approve
                success: {status: [200], signal: RESTResponded}
                side_effects:
                  - kind: external_api
                    target: crm.account
                    state: approval_requested
                reversibility:
                  classification: compensatable
                  undo: revokeApproval

  invalid_because:
    - approve is RPC-shaped and belongs under operations.
    - Resource operation names are limited to get, set, create, and delete.
- missing_schema_for_body:
  # Body templates must reference declared runtime input.
  rest:
    version: v1
    clients:
      billing:
        base_url: https://billing.example
        operations:
          create_invoice:
            method: POST
            path: /invoices
            body:
              customer: "{{ params.customer }}"
            success: {status: [201], signal: RESTResourceWritten}
            side_effects:
              - kind: external_api
                target: billing.invoice
                state: invoice_created
            reversibility:
              classification: compensatable
              undo: voidInvoice

  invalid_because:
    - body references params.customer without a declared body_schema or request binding.
- unsupported_auth:
  # Auth types are limited to the supported set.
  rest:
    version: v1
    auth:
      cloud:
        type: magic_signature
        token_ref: cloud_token

  invalid_because:
    - magic_signature is not a supported auth type.
- duplicate_openapi_operation:
  # Imported operation IDs must remain unique in one REST definition.
  rest:
    version: v1
    openapi:
      service_a:
        path: ./service-a.yaml
        expose: [createThing]
      service_b:
        path: ./service-b.yaml
        expose: [createThing]

  invalid_because:
    - createThing maps to two imported operations in the same REST definition scope.
- unsafe_server_binding:
  # Public listeners need an explicit network policy.
  rest:
    version: v1
    servers:
      public_control:
        address: 0.0.0.0:8080
        endpoints:
          inject:
            method: POST
            path: /inject
            binding: emit_dynamic_signal
            allowed_signals: [Approved, Rejected]

  invalid_because:
    - Public listener binding needs explicit network policy.
    - Dynamic signal endpoints need authentication or another configured boundary control.

<!-- field: related_documents -->
## Related documents

- [srd052-declaration-fragments](../software-requirements/srd052-declaration-fragments.md)
- [srd056-unit-libraries](../software-requirements/srd056-unit-libraries.md)
- [cfg-tool-declaration-format](tool-declaration-format.md)
- [cfg-tool-authoring-standard](tool-authoring-standard.md)
- [sm-rest-tool-language](../semantic-models/rest-tool-language.md)
- [srd027-rest-tool-runtime](../software-requirements/srd027-rest-tool-runtime.md)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md)
- [srd029-rest-server-tools](../software-requirements/srd029-rest-server-tools.md)
- [srd030-machine-request-rest-binding](../software-requirements/srd030-machine-request-rest-binding.md)
- [srd031-declarative-document-resources](../software-requirements/srd031-declarative-document-resources.md)
- [sm-machine-request-http](../semantic-models/machine-request-http.md)
- [sm-document-resource-language](../semantic-models/document-resource-language.md)
- [srd026-lifecycle-tools](../software-requirements/srd026-lifecycle-tools.md)

<!-- field: implementation_status -->
## Implementation status

<!-- field: implemented | form: "tree" -->
### Implemented

1. internal/tools/rest typed config structs.
1. internal/tools/rest YAML loader.
1. OpenAPI import compilation through github.com/getkin/kin-openapi/openapi3.
1. rest_await_event source-filter config through ToolDef config.
1. emit_signal lifecycle exit routes and exit_agent ToolDef config.
1. machine_request endpoint binding config.
1. signal_source endpoint binding with trusted closed request mapping, typed sensitive paths, host-runner injection, configured outcome responses, bounded tracing, and no inbound queue.
1. filesystem resource words used by the Knowledge Manager request machine.
1. generic machine_request response trace fields for server, route, and request machine emitted directly by internal/tools/rest.
1. always-on Knowledge Manager conformance coverage for generic machine_request routes under plain go test.
1. - "request_binding.input_mapping and request_binding.carry_forward with body_source previous_result consumption (V28-V31)": a client operation selects declared params from the prior Result output and copies declared inputs forward under a carried key. Validation rejects undeclared targets and transport-authority overrides.
1. - "request_binding body_source command_state consumption (srd028 R13, V29/V32)": a client operation selects declared params from a labeled prior step in the srd038 command-state store through input_mapping $from(<label>).<path> selectors resolved against the injected core CommandStateView, with the R12.4 rejection list applied identically and previous_result threading unchanged. command_state is structurally valid; it is rejected only at runtime when no store view is configured. V32 requires $from(...) selectors under command_state and $. selectors under previous_result. Delivered by GH-278 (REST consumption and V32 validation) over GH-277 (CommandStateView).

<!-- field: partial -->
### Partial

- read_policy values outside first-available are rejected until implemented.

<!-- field: reserved -->
### Reserved

- top-level document_resources config. ValidateDefinition rejects this field.
- machine_request.document_resources config. ValidateDefinition rejects this field.
- lifecycle_control endpoint config.
