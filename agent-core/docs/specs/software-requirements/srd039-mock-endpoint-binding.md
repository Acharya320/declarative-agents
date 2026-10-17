---
id: srd039-mock-endpoint-binding
title: Mock REST Endpoint Binding
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Scenario-driven agent tests need dependency mocks. A mock profile
stands in for a real upstream service by serving canned REST responses, so a
subject agent can be exercised hermetically: no live Chroma, Ollama, or
Kubernetes, and every dependency response is controlled by the test scenario.

The REST server runtime already declares endpoints with bindings such as
static_assets, redirect, and machine_request. None of those covers the mock
role declaratively: static_assets cannot vary status codes or script a
failure on the second call, and machine_request requires authoring a machine
for what is a fixed lookup. Tests also need to verify that the subject
actually called its dependency with the right request, which requires the
mock to record what it received and serve that record back to a validator.

<!-- field: goals -->
## Goals

- G1: Define a mock endpoint binding that serves canned responses from fixture data.
- G2: Define the fixture format, with routes keyed by method and path and ordered response scripts.
- G3: Define a request log the mock records and serves so validators can assert the subject's outbound calls.
- G4: Define validation so malformed fixtures fail configuration validation and server startup, not the first request.
- G5: Keep the binding concurrency-safe and deterministic.

<!-- req: R1 -->
## R1 - Binding scope

<!-- req: R1.1 -->
### R1.1

REST server endpoints SHALL support a mock binding that serves responses from fixture data without dispatching a machine transition.

<!-- req: R1.2 -->
### R1.2

The mock binding SHALL follow the shared REST server runtime rules of srd029-rest-server-tools for listener lifecycle, ports, and telemetry.

<!-- req: R1.3 -->
### R1.3

Mock fixture data SHALL come from a fixture file path declared in the endpoint configuration, with environment expansion applied per the shared REST definition loading rules, or from inline routes in the endpoint configuration.

<!-- req: R1.4 -->
### R1.4

Fixture data SHALL load when the server starts, not per request.

<!-- req: R1.5 -->
### R1.5

The binding SHALL remain part of the shared agent REST server runtime so the shipped mock profile uses the same agent binary and lifecycle. A third-party mock CLI MAY replace it only after a compatibility spike proves ordered-last-repeats scripts, unmatched-request journaling, per-instance isolation, startup validation, concurrent sequencing, declared lifecycle/telemetry, fixture migration, and runtime-image provisioning (GH-1386).

<!-- req: R1.6 -->
### R1.6

The shipped mock profile SHALL inherit the REST runtime's public-listener guard, default its listener to loopback, and require an explicit default-false authorization before binding a wildcard or public address.

<!-- req: R2 -->
## R2 - Fixture format

<!-- req: R2.1 -->
### R2.1

A fixture SHALL declare routes keyed by HTTP method and literal path.

<!-- req: R2.2 -->
### R2.2

Each route SHALL declare an ordered list of one or more responses, each with a status code and optional headers and body.

<!-- req: R2.3 -->
### R2.3

The Nth request to a route SHALL receive the Nth response; after the list is exhausted, the last response SHALL repeat.

<!-- req: R2.4 -->
### R2.4

A request that matches no route SHALL receive a 404 response and SHALL be recorded in the request log as a miss.

<!-- req: R2.5 -->
### R2.5

Response bodies SHALL be literal content; the binding SHALL NOT template or compute bodies from the request.

<!-- req: R3 -->
## R3 - Request log

<!-- req: R3.1 -->
### R3.1

The mock SHALL record every received request in order, capturing method, path, headers, and body.

<!-- req: R3.2 -->
### R3.2

A declared log endpoint SHALL serve the recorded requests as structured JSON, including whether each request matched a route.

<!-- req: R3.3 -->
### R3.3

The request log SHALL be scoped to the server instance and reset when the server starts.

<!-- req: R4 -->
## R4 - Validation

<!-- req: R4.1 -->
### R4.1

A fixture with an unknown HTTP method, an empty response list, or a duplicate method-and-path route SHALL fail validation with a named error.

<!-- req: R4.2 -->
### R4.2

Fixture validation SHALL run under --validate-config and at server startup; a malformed fixture SHALL prevent the server from serving.

<!-- req: R5 -->
## R5 - Concurrency and determinism

<!-- req: R5.1 -->
### R5.1

Per-route response sequencing and the request log SHALL be safe under concurrent requests.

<!-- req: R5.2 -->
### R5.2

For a serial request sequence, responses and log order SHALL be deterministic.

<!-- field: non_goals -->
## Non-goals

- Request matching beyond method and literal path; body and header matchers are validator-side assertions over the request log.
- Response templating, latency injection, or proxying to a real upstream.
- A client-side mock; the binding serves inbound requests only.
- Persistence of the request log across server restarts.
- Requiring Prism, Mockoon, WireMock, or another dedicated mock executable solely because it can serve some canned HTTP responses.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

A server endpoint declared with the same-binary mock binding serves the fixture's literal responses for a matching route in declared order over the shared REST server runtime, repeats the last response after the script is exhausted, requires compatibility plus provisioning before any third-party replacement, and accepts public listeners only through explicit authority.

<!-- field: traces -->
#### Traces

- [R1.1](#r11)
- [R1.2](#r12)
- [R1.3](#r13)
- [R1.4](#r14)
- [R1.5](#r15)
- [R1.6](#r16)
- [R2.1](#r21)
- [R2.2](#r22)
- [R2.3](#r23)
- [R2.5](#r25)

<!-- ac: AC2 -->
### AC2

A request matching no fixture route receives a 404 and appears in the request log as a miss.

<!-- field: traces -->
#### Traces

- [R2.4](#r24)
- [R3.1](#r31)

<!-- ac: AC3 -->
### AC3

The log endpoint returns the received requests in order with method, path, headers, and body, sufficient for a validator to assert what the subject sent.

<!-- field: traces -->
#### Traces

- [R3.1](#r31)
- [R3.2](#r32)
- [R3.3](#r33)

<!-- ac: AC4 -->
### AC4

A malformed fixture fails --validate-config and server startup with a named error before any request is served.

<!-- field: traces -->
#### Traces

- [R4.1](#r41)
- [R4.2](#r42)

<!-- ac: AC5 -->
### AC5

Concurrent requests to scripted routes neither skip nor duplicate script positions, and a serial sequence is deterministic.

<!-- field: traces -->
#### Traces

- [R5.1](#r51)
- [R5.2](#r52)
