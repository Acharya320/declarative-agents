---
id: sm-rest-tool-language
title: REST Tool Language Semantic Model
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: purpose -->
## Purpose

Defines the vocabulary for declarative REST tools. REST client words and REST
server words are boundary tools in the shared tool language. MachineSpec
decides the order. A REST word may run synchronously, split outbound work
into send and await phases, launch server state, stop server state, or handle
a signal from an inbound request.

REST definitions provide service specificity. Go code provides generic HTTP
capability through net/http and internal/tools/rest. That split keeps
transport authority in trusted profile and tool configuration while preserving
visible grammar flow.

<!-- field: requirements_source -->
## Requirements source

```yaml
canonical:
    - docs/specs/software-requirements/srd027-rest-tool-runtime.md
    - docs/specs/software-requirements/srd028-rest-client-tools.md
    - docs/specs/software-requirements/srd029-rest-server-tools.md
    - docs/specs/software-requirements/srd030-machine-request-rest-binding.md
supporting:
    - docs/specs/semantic-models/tool-language.md
    - docs/specs/semantic-models/rollback-lifecycle.md
    - docs/guides/declarative-programming-pattern.md
    - docs/guides/boundary-tools-and-async.md
```

<!-- field: primary_entities | form: "typed table" -->
## Primary entities

| name | kind | responsibility | owns | not_responsible_for | references |
| --- | --- | --- | --- | --- | --- |
| RestDefinition | declarative_boundary_program | Defines shared REST clients and servers, including their resources, operations, endpoints, mappings, limits, auth references, and side effects. | [rest_authority, rest_vocabulary] | [MachineSpec sequencing, model-provider semantics, arbitrary runtime URLs] |  |
| RestClient | outbound_actor_config | Names an upstream service and binds auth, base endpoint policy, network policy, and operation definitions for outbound HTTP calls. | [outbound_authority] |  | [Operation, Resource, OpenAPIImport] |
| RestServer | inbound_actor_config | Names a server address, route set, endpoint bindings, queue policy, stream policy, and shutdown behavior for inbound HTTP traffic. | [listener_authority, server_state] |  | [Endpoint, InboundEvent, RestToolState] |
| Resource | rest_domain_object | Groups resource-shaped operations that read, create, set, or delete stored upstream state. | [resource_identity, resource_operation_names] |  |  |
| Operation | outbound_rest_action | Defines one configured outbound HTTP action with method, path template, request binding, response mapping, status mapping, retry policy, and side-effect contract. | [operation_id, outbound_status_mapping] |  |  |
| Endpoint | inbound_rest_route | Defines one configured route with method, path template, request validation, endpoint binding kind, response mapping, and allowed signal behavior. | [route_id, inbound_validation] |  |  |
| RequestBinding | parameter_mapping | Maps runtime values from explicit tool input, previous Result output, or command state into declared path, query, header, and body fields. | [runtime_parameter_contract] |  |  |
| ResponseMapping | output_mapping | Maps HTTP status, headers, and body into Result output, domain signals, error output, resource metadata, request IDs, and redacted diagnostics. | [result_shape, signal_mapping] |  |  |
| OpenAPIImport | source_adapter | Loads OpenAPI 3 operations, paths, request schemas, and response schemas into the same internal RestDefinition model as hand-authored YAML. | [source_document_reference] |  |  |
| AsyncRequest | outbound_pending_state | Records an accepted outbound submission created by rest_client_send and consumed by rest_client_await. | [request_id, correlation_token, idempotency_token, pending_response_state] |  |  |
| InboundEvent | inbound_signal_message | Carries a validated server request, configured signal, source route, payload, and request metadata from HTTP handler state to rest_server_await. | [event_payload, event_signal] |  |  |
| RestEventSource | fan_in_source | Names one server queue or route-filtered event stream that rest_await_event can read with configured signal, route, timeout, and stopped-source policies. | [source_filter, stopped_source_policy] |  | [RestServer, Endpoint, InboundEvent] |
| LifecycleControlEndpoint | inbound_control_route | Validates fixed or allowlisted signal enqueue requests and publishes declared lifecycle events for MachineSpec. | [control_request_validation, control_signal_mapping] |  | [Endpoint, InboundEvent, MachineSpec] |
| MachineRequestBinding | inbound_machine_route | Validates one inbound HTTP request and starts one short-lived MachineSpec execution with configured request input and response mapping. | [request_to_machine_mapping, terminal_response_selection] |  | [Endpoint, MachineSpec, ToolDef, RequestEvent] |
| RequestEvent | grammar_input | Carries validated request fields and request metadata into a request-scoped machine sentence. | [validated_payload, request_metadata, correlation_id] |  |  |
| ResponseWord | response_boundary_word | Produces structured output that the machine_request binding maps to HTTP status, headers, content type, and body. | [response_output_contract, response_signal] |  | [ToolDef, MachineSpec] |
| LifecycleControlWord | profile_selected_lifecycle_vocabulary | Executes lifecycle behavior such as exit_agent after MachineSpec routes a validated REST event to that word. | [lifecycle_side_effects, lifecycle_undo_contract] |  | [ToolDef, MachineSpec] |
| RestToolState | command_domain_state | Tracks shared REST clients, shared REST servers, async requests, inbound queues, active streams, route ownership, and shutdown status. | [rest_command_state] |  |  |
| ToolDef | declaration_contract | Declares each REST tool word, its runtime contract, and its lifecycle contract. |  |  | [RestDefinition] |
| MachineSpec | declarative_program | Composes REST words into visible sync and async workflows through states, signals, and transition actions. |  |  | [ToolDef] |

<!-- field: tool_words -->
## Tool words

```yaml
client_resource:
    words: [rest_client_get, rest_client_set, rest_client_create, rest_client_delete]
    meaning: Resource words operate on configured stored state and expose resource-specific signals, effects, and compensation.
    preferred_when: The endpoint maps to read, write, create, delete, or deactivate behavior for a named resource.
client_invoke:
    words: [rest_client_invoke]
    meaning: Invoke calls a configured RPC-shaped or action-shaped operation that does not fit the resource model.
    preferred_when: The endpoint represents approve, reject, trigger, exchange, calculate, search, or another action that is not simple resource state access.
client_async:
    words: [rest_client_send, rest_client_await]
    meaning: Send completes one HTTP submission transaction and records accepted AsyncRequest state. Await receives the stored completion and emits the configured result signal.
server:
    words: [rest_server_launch, rest_server_await, rest_server_stop]
    meaning: Launch activates configured routes. Await receives validated InboundEvent messages. Stop shuts down shared server state.
event_aggregation:
    words: [rest_await_event]
    meaning: Await receives one event from multiple configured REST server event sources and returns the selected signal to MachineSpec.
lifecycle_control:
    words: [exit_agent]
    meaning: Lifecycle control words consume validated REST event payloads or profile config and perform explicit MachineSpec-selected lifecycle actions.
machine_request:
    words: [machine_request]
    meaning: Endpoint binding that turns one accepted HTTP request into one short-lived MachineSpec run and maps terminal response output back to HTTP.
```

<!-- field: relationships -->
## Relationships

```yaml
RestDefinition_to_ToolDef:
    meaning: REST ToolDefs point to configured REST definitions by named client, server, resource, operation, or endpoint reference.
    invariant: ToolDef owns the vocabulary word contract; RestDefinition owns the HTTP-specific authority.
RestDefinition_to_OpenAPIImport:
    meaning: OpenAPIImport may populate resources, operations, endpoints, request schemas, and response schemas in RestDefinition.
    invariant: Imported definitions use the same internal model as hand-authored YAML and do not generate Go clients or servers.
RestClient_to_Operation:
    meaning: A RestClient owns outbound Operation definitions and resource groupings.
    invariant: Runtime input can fill declared fields but cannot choose method, destination, or auth authority.
Resource_to_client_words:
    meaning: Resource operations bind to rest_client_get, rest_client_set, rest_client_create, and rest_client_delete.
    invariant: Resource words keep state access and rollback expectations visible in the word contract.
Operation_to_invoke:
    meaning: RPC-shaped operations bind to rest_client_invoke when no resource word describes the action.
    invariant: Invoke remains one configured boundary word and must not hide a service workflow.
rest_client_send_to_AsyncRequest:
    meaning: rest_client_send creates AsyncRequest state after submitting configured outbound work.
    invariant: The returned RESTAccepted signal means submission happened; final service outcome belongs to rest_client_await.
rest_client_await_to_AsyncRequest:
    meaning: rest_client_await reads AsyncRequest state and maps completion, timeout, missing state, or domain failure into Result.Signal.
    invariant: Await words declare every signal they can emit in the machine state where they run.
RestServer_to_Endpoint:
    meaning: A RestServer owns configured Endpoint routes and shared listener state.
    invariant: HTTP handlers validate requests and enqueue configured events; they do not choose hidden MachineSpec actions.
Endpoint_to_InboundEvent:
    meaning: Signal-producing endpoints translate accepted HTTP requests into InboundEvent messages.
    invariant: Invalid requests do not enqueue grammar events.
rest_server_await_to_InboundEvent:
    meaning: rest_server_await consumes queued InboundEvent messages and returns the event signal to MachineSpec.
    invariant: The HTTP request becomes grammar-visible only through the await word result.
rest_await_event_to_RestEventSource:
    meaning: rest_await_event waits over one or more configured RestEventSource entries and returns the first matching event, timeout, stopped-source, or error signal.
    invariant: Source filters are declarative and include server names, route names, allowed signals, timeout, and stopped-source behavior.
LifecycleControlEndpoint_to_InboundEvent:
    meaning: A lifecycle control endpoint validates a control request and enqueues an event such as AgentExitRequested, PauseRequested, RollbackRequested, ResumeRequested, or an allowlisted injected signal.
    invariant: The HTTP handler does not execute exit, pause, rollback, resume, or signal injection directly.
MachineRequestBinding_to_RequestEvent:
    meaning: A machine_request endpoint validates HTTP input and builds one RequestEvent for the configured request machine.
    invariant: Invalid HTTP input does not enqueue events or start a machine run.
MachineRequestBinding_to_MachineSpec:
    meaning: The binding resolves a trusted profile, MachineSpec, selected ToolDefs, and response policy for one short-lived run.
    invariant: Runtime HTTP callers cannot choose profile, machine, tool declaration roots, model config, credentials, URLs, or workspace roots.
MachineSpec_to_ResponseWord:
    meaning: Request machines route to response vocabulary when they need to shape an HTTP response.
    invariant: Direct Go handlers do not compose Agent Architecture request behavior.
ResponseWord_to_HTTPResponse:
    meaning: Response word output and terminal signal mappings produce HTTP status, headers, content type, and body.
    invariant: HTTP responses use configured terminal response selection and declared response output fields only.
InboundEvent_to_LifecycleControlWord:
    meaning: MachineSpec routes a lifecycle event from rest_server_await or rest_await_event to profile-selected vocabulary such as exit_agent or checkpoint_rollback.
    invariant: Lifecycle behavior remains visible as declared ToolDef words and MachineSpec transitions.
RestToolState_to_lifecycle:
    meaning: RestToolState is command/domain state for async requests, server queues, active streams, and route ownership.
    invariant: Lifecycle storage captures REST command state separately from workspace files and agent loop position.
```

<!-- field: state_layers -->
## State layers

```yaml
agent_state:
    owner: Loop
    rest_examples:
        - Current MachineSpec state after RESTAccepted, RESTResponded, ServerLaunched, or configured inbound signals.
        - Iteration and budget counters around blocking or async REST words.
    rollback_rule: Rest tools do not own agent state directly; Loop snapshots it when lifecycle policy requests a checkpoint.
command_domain_state:
    owner: RestToolState
    rest_examples:
        - AsyncRequest records created by rest_client_send.
        - Shared server registry entries created by rest_server_launch.
        - InboundEvent queues consumed by rest_server_await.
        - Active stream metadata and shutdown status.
    rollback_rule: REST commands undo or compensate command-owned state according to ToolDef undo contracts.
workspace_environment_state:
    owner: Workspace
    rest_examples:
        - REST tools normally do not mutate workspace files.
        - Future REST tools may write downloaded artifacts only when ToolDef declares filesystem side effects.
    rollback_rule: Workspace restore applies only when a REST word declares workspace effects.
external_service_state:
    owner: remote_http_service
    rest_examples:
        - Resources created, updated, deleted, approved, triggered, or cancelled through outbound REST calls.
        - Network listeners exposed by REST server tools.
    rollback_rule: External effects use noop, configured compensation, cancellation, or confirmation rules. The lifecycle model cannot erase remote audit history.
```

<!-- field: invariants -->
## Invariants

- REST tools are boundary words because they cross HTTP client, HTTP server, external service, human, or agent-to-agent boundaries.
- REST definitions provide specificity; Go code provides generic capability through internal/tools/rest and net/http.
- MachineSpec owns sync and async sequencing for REST words.
- ToolDef owns the word contract for input, output, emitted signals, side effects, errors, and undo or compensation.
- Runtime input never supplies arbitrary methods, URLs, credentials, redirect policy, or network destinations.
- Resource words are preferred for stored-state operations.
- rest_client_invoke is reserved for RPC-shaped or action-shaped operations.
- rest_client_send and rest_client_await make outbound async state visible in grammar words.
- rest_server_launch and rest_server_await make inbound HTTP work visible in grammar words.
- rest_await_event makes multi-source inbound HTTP work visible as one grammar await word.
- HTTP handlers validate and translate requests; they do not run hidden workflows.
- Lifecycle control endpoints enqueue events; lifecycle behavior runs only through MachineSpec-selected vocabulary.
- exit_agent is generic lifecycle vocabulary and is not owned by Knowledge Manager-specific HTTP handlers.
- machine_request endpoints validate one HTTP request and run one short-lived MachineSpec sentence.
- machine_request request behavior is composed by MachineSpec transitions, not direct Go handlers.
- Response words are the only bridge from request-machine output back to HTTP status, headers, content type, and body.
- OpenAPI imports compile into internal REST definitions and do not generate Go code in release 03.0.
- REST command/domain state is separate from agent state and workspace/environment state.
- External service effects are reversible only when the configured operation declares a valid compensation path.
- All traces, errors, metrics, and Result output redact configured secret material.

<!-- field: non_goals -->
## Non goals

- Defining the full REST tool config YAML schema; see cfg-rest-tool-format.
- Implementing modeled lifecycle_control endpoint config in this semantic model update.
- Adding generated Go clients or generated Go servers.
- Replacing MachineSpec control flow with REST handler code.
- Executing lifecycle control directly inside HTTP handlers or rest_await_event.
- Executing Agent Architecture request behavior directly inside Go HTTP handlers.
- Allowing model text or inbound callers to choose transport authority.
- Modeling every provider-specific REST API.

<!-- field: diagrams -->
## Diagrams

````yaml
outbound_sync: |
    ```mermaid
    sequenceDiagram
      participant M as MachineSpec
      participant T as REST client word
      participant D as RestDefinition
      participant H as net/http
      participant R as Result
      M->>T: dispatch rest_client_get or invoke
      T->>D: resolve client, resource, operation
      T->>D: bind declared request fields
      T->>H: send configured HTTP request
      H-->>T: status, headers, body
      T->>D: apply ResponseMapping
      T-->>R: structured output and signal
      R-->>M: route next transition
    ```
outbound_async_send_await: |
    ```mermaid
    sequenceDiagram
      participant M as MachineSpec
      participant S as rest_client_send
      participant A as AsyncRequest
      participant H as net/http
      participant W as rest_client_await
      M->>S: dispatch send
      S->>H: submit configured outbound work
      H-->>S: accepted status and response
      S->>A: store accepted request state
      S-->>M: RESTAccepted
      M->>W: dispatch await later
      W->>A: read request state
      W-->>M: RESTResponded, RESTDomainFailed, RESTMissing, RESTAwaitTimedOut, or CommandError
    ```
inbound_server_launch_await: |
    ```mermaid
    sequenceDiagram
      participant M as MachineSpec
      participant L as rest_server_launch
      participant S as RestServer
      participant E as Endpoint handler
      participant Q as InboundEvent queue
      participant A as rest_server_await
      M->>L: dispatch launch
      L->>S: bind address and register routes
      L-->>M: ServerLaunched
      E->>E: validate inbound HTTP request
      E->>Q: enqueue configured signal and payload
      M->>A: dispatch await
      A->>Q: receive event
      A-->>M: configured event signal
    ```
multi_source_rest_await_event: |
    ```mermaid
    sequenceDiagram
      participant H1 as Control endpoint
      participant H2 as Domain endpoint
      participant Q as Server event queues
      participant A as rest_await_event
      participant M as MachineSpec
      H1->>Q: enqueue AgentExitRequested
      H2->>Q: enqueue DomainEventReceived
      M->>A: dispatch with source and signal filters
      A->>Q: select first matching event
      A-->>M: selected event signal and payload
      M->>M: route to visible lifecycle or domain action
    ```
machine_request_binding: |
    ```mermaid
    sequenceDiagram
      participant H as HTTP handler
      participant B as machine_request binding
      participant M as MachineSpec
      participant W as request words
      participant R as response word
      H->>H: validate request
      H->>B: accepted request fields
      B->>M: seed RequestEvent
      M->>W: compose request sentence
      W-->>M: declared signals
      M->>R: prepare response
      R-->>B: response output
      B-->>H: mapped HTTP response
    ```
openapi_import_to_internal_definition: |
    ```mermaid
    flowchart LR
      O[OpenAPI 3 document] --> I[OpenAPIImport]
      I --> V[openapi3 validation]
      V --> D[RestDefinition]
      D --> C[RestClient operations]
      D --> S[RestServer endpoints]
      C --> T[REST client ToolDefs]
      S --> U[REST server ToolDefs]
    ```
````

<!-- field: related_documents -->
## Related documents

- sm-tool-language
- sm-rollback-lifecycle
- srd027-rest-tool-runtime
- srd028-rest-client-tools
- srd029-rest-server-tools
- srd030-machine-request-rest-binding
- sm-machine-request-http
- srd026-lifecycle-tools
- cfg-tool-declaration-format
- cfg-tool-authoring-standard

<!-- field: implementation_status -->
## Implementation status

```yaml
implemented:
    - internal/tools/rest and tools/builtin/rest.
    - docs/specs/config-formats/rest-tool-format.md.
    - REST client resource, invoke, send, and await words.
    - REST server launch, await, stop, and rest_await_event fan-in words.
    - exit_agent lifecycle control word selected by MachineSpec after REST control events.
    - machine_request endpoint binding for short-lived request-machine runs.
    - OpenAPIImport preserves configured endpoint behavior while importing transport and schema fields.
    - ResponseMapping schema validation is enforced for client responses and machine_request terminal responses.
    - External service compensation metadata produces REST client undo mementos for mutating operations and drives rollback-time compensation execution through receipt rollback.
    - Top-level REST document_resources and machine_request.document_resources are typed target-format fields and are rejected explicitly until generic loading exists.
partial:
    - LifecycleControlEndpoint remains a target semantic entity. Current release 03.0 control routes use emit_signal endpoint bindings and MachineSpec transitions; modeled lifecycle_control endpoint config remains reserved.
    - RestEventSource non-default read_policy behavior remains reserved; unsupported values are rejected during ToolDef factory validation.
```
