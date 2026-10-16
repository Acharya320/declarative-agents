---
id: srd004-declarative-ux
title: Presentation Contract and Panel Composition
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Agents, tools, and REST routes are declared, but application user interfaces
are hand-written single-page applications with no shared unit smaller than
the whole bundle. Ten React applications across this repository,
cohere-demo, and agentic-wiki-mesh share only one 39-line token file. The two
demo repositories forked the chatbot-mesh chatbot and observer and diverged:
about twenty same-named files exist in both, none byte-identical. The token
file reaches cohere-demo through a sibling-checkout relative import and
agentic-wiki-mesh through two vendored copies.

Nothing pins the interface between a panel and its backend. A panel reads
whatever the application's backend happens to serve, so shared code drifts
with each backend. The same proxy exists under two names: the chatbot
declares /monitor-proxy/{agent}/{path...} while the cohere-demo observer
declares /trace-proxy/{agent}/{path...}. ui.yaml describes routes that a Go
test cross-checks against a hand-written routes.ts, so the declaration is a
second copy rather than the input.

We fix this contract first: a closed endpoint surface that shared panels may
read, a panel unit, a build-time composition model driven by ui.yaml, one
injectable client, and one shared package that owns tokens, clients, shell,
and observability panels. The observer UI becomes a bundle the REST tool
serves by declaration.

<!-- field: goals -->
## Goals

- G1: Pin the endpoint surface shared panels may read, citing the SRD that owns each response.
- G2: Define the panel as the reuse unit and the rules for composing panels into one application bundle.
- G3: Make ui.yaml the composition input of a generic shell instead of a test-checked description.
- G4: Give shared presentation code, including design tokens, one owner and one distribution path for in-repository and external consumers.
- G5: Let an observer agent carry no UI files by serving an embedded bundle configured only by declaration.

<!-- req: R1 -->
## R1 - Presentation contract

<!-- req: R1.1 -->
### R1.1

The presentation contract shall be the closed set of endpoints below.
This SRD names each endpoint and its owner; the owning SRD defines
the response schema, and this SRD does not restate it.

<!-- field: evidence -->
#### Evidence

The UI kit exports the list as one constant, and a kit test fails when
a published panel names an endpoint outside it.

<!-- field: endpoints -->
#### Endpoints

| path | owner |
| --- | --- |
| GET /monitor/state | agent-core srd033-monitor-rest-api R3.9 |
| GET /monitor/machines | agent-core srd033-monitor-rest-api R8 |
| GET /monitor/tools/declared | agent-core srd033-monitor-rest-api R9 |
| GET /monitor/events/stream | agent-core srd033-monitor-rest-api R4, R2.12 |
| GET /monitor/fleet | catalog srd021-observer R3.1 over chatbot-mesh srd008-observer-agent R3 |
| GET /query/traces?page_size=&offset= | catalog srd020-collector R5.1 |
| GET /query/traces/{trace_id} | catalog srd020-collector R5.2 |

<!-- req: R1.2 -->
### R1.2

A shared panel shall read only fields that the owning SRD's response
schema defines. A field a panel needs and no owner defines shall be
added to the owner first.

<!-- field: evidence -->
#### Evidence

Kit fixtures carry only owner-defined fields, and panel render tests
run against those fixtures alone.

<!-- req: R1.3 -->
### R1.3

Domain endpoints such as /api/v1/chat, /api/v1/documents, or
/provisioning/api are outside the contract. A panel that needs one is
a domain panel and stays in its application.

<!-- field: evidence -->
#### Evidence

The endpoint-usage review in the rel15.0 use case classifies every
URL fetched by the current UIs as contract or domain.

<!-- req: R2 -->
## R2 - Agent proxy convention

<!-- req: R2.1 -->
### R2.1

A shell shall reach another agent's contract endpoints through
GET /monitor-proxy/{agent}/{path...} on its own origin, where path is
a contract path without its leading slash. The serving REST
definition declares the route; the shell never addresses an agent's
host directly.

<!-- field: evidence -->
#### Evidence

The chatbot-mesh chatbot rest.yaml declares the route, and kit client
tests build proxy URLs only through one function.

<!-- req: R2.2 -->
### R2.2

A proxy response of 404 shall mean the agent is not deployed. The
client shall return a typed not-deployed result rather than throw,
and the shell shall hide every panel whose monitored agents are all
not deployed.

<!-- field: evidence -->
#### Evidence

Kit client and shell tests cover the not-deployed result and the
hidden panel.

<!-- req: R2.3 -->
### R2.3

/trace-proxy/{agent}/{path...} is retired. A REST definition that
serves a kit-built UI shall declare /monitor-proxy and no other proxy
name. Existing /trace-proxy routes are removed in the consumer
migrations.

<!-- field: evidence -->
#### Evidence

No REST definition under applications/ declares /trace-proxy, and the
embedded observer wrapper declares /monitor-proxy.

<!-- req: R2.4 -->
### R2.4

Trace endpoints shall be reached through the trace backend ui.yaml
declares, not through a constant compiled into a panel. The backend
is an agent name, reached through that agent's proxy, or a
same-origin path prefix starting with "/" under which /query/traces
lives; "/" is the origin root, and "/monitor-proxy/collector" reads
the same URLs as the name collector. The shell derives the backend
from trace_backend.query_path when it is set, taking the prefix
before /query/traces/{trace_id} ("/" when the prefix is empty), and
from trace_backend.name otherwise. A UI served beside the query
surface, such as the collector's, uses the same-origin form.

<!-- field: evidence -->
#### Evidence

The kit trace client takes the backend from shell configuration;
kit tests cover both forms and the shell's derivation, and the
collector UI reads its traces at the origin root.

<!-- req: R3 -->
## R3 - Compatibility policy

<!-- req: R3.1 -->
### R3.1

The contract shall carry a major version, starting at 1, exported by
the UI kit. Within a major version, owners may add endpoints or
response fields and shall not remove, rename, or retype a field a
published panel reads.

<!-- field: evidence -->
#### Evidence

The kit exports the contract version, and the kit package version
changes major whenever the contract major changes.

<!-- req: R3.2 -->
### R3.2

The UI kit shall keep one recorded response fixture per contract
endpoint and shall run each client and each panel against those
fixtures. The fixtures are the consumer contract tests.

<!-- field: evidence -->
#### Evidence

Kit tests parse every fixture through its client and render every
published panel from fixtures.

<!-- req: R3.3 -->
### R3.3

A change to an owner SRD that alters a field a fixture carries shall
update the fixture in the same change.

<!-- field: evidence -->
#### Evidence

Review of owner SRD changes; kit tests fail when a fixture and a
client type disagree.

<!-- req: R4 -->
## R4 - Panel unit

<!-- req: R4.1 -->
### R4.1

A panel shall be a React component exported together with a manifest
of id, title, route, required_endpoints, and monitored_agents.

<!-- field: evidence -->
#### Evidence

Every published panel exports its manifest, and a kit test checks the
manifest fields.

<!-- req: R4.2 -->
### R4.2

required_endpoints shall be a subset of R1.1. monitored_agents shall
be either a list of agent names or the literal declared, meaning the
agents ui.yaml lists.

<!-- field: evidence -->
#### Evidence

A kit test rejects a manifest that names an endpoint outside R1.1.

<!-- req: R4.3 -->
### R4.3

Kit panel code shall import nothing from an application directory and
shall perform no I/O except through the client of R6.

<!-- field: evidence -->
#### Evidence

A kit source scan finds no application import and no direct fetch or
EventSource construction outside the client module.

<!-- req: R4.4 -->
### R4.4

The kit panels are StatusBar, FleetPanel, TracePanel, MachineView with
AgentPanel, Topology, and AgentCard. ProvenancePanel is a
cohere-demo domain panel because it reads /api/v1/documents.

<!-- field: evidence -->
#### Evidence

The kit panel registry exports exactly these panels.

<!-- req: R5 -->
## R5 - Build-time composition

<!-- req: R5.1 -->
### R5.1

An application bundle shall be composed at build time from its
ui.yaml, a registry of kit panels, and the application's domain
panels. Runtime module loading is not supported.

<!-- field: evidence -->
#### Evidence

The kit Vite plugin reads ui.yaml during the build, and the shell
mounts panels only from the registry the entrypoint supplies.

<!-- req: R5.2 -->
### R5.2

Routes and sidebar shall derive from ui.yaml alone. A panel id in
ui.yaml with no registry entry shall fail the build.

<!-- field: evidence -->
#### Evidence

Shell tests render routes and sidebar from a fixture ui.yaml, and a
plugin test fails the build for a missing registry entry.

<!-- req: R5.3 -->
### R5.3

The shell shall keep the existing panel URL scheme so current deep
links resolve to the same panel.

<!-- field: evidence -->
#### Evidence

A table test replays the splitPanelPath cases of every current copy.

<!-- req: R6 -->
## R6 - Injectable client and base URL

<!-- req: R6.1 -->
### R6.1

Every fetch and EventSource in kit code and in kit-shell
applications shall go through one client created with an injectable
base URL, fetch function, and EventSource constructor. The default
base URL is same-origin.

<!-- field: evidence -->
#### Evidence

A kit test runs the monitor client against a remote base URL and
asserts absolute URLs; the source scan of R4.3 applies.

<!-- req: R6.2 -->
### R6.2

Panels shall obtain the client from shell context, so a desktop
shell can point the same bundle at remote base URLs without a
rebuild.

<!-- field: evidence -->
#### Evidence

Panels read the client only through the kit context hook.

<!-- req: R7 -->
## R7 - ui.yaml version 2

<!-- req: R7.1 -->
### R7.1

ui.yaml version 2 shall add version, panels, and branding to the
version 1 fields, which already carry monitored_agents,
trace_backend, and the presentation flags. Each panels entry carries
id, package, route, and optionally export, label, sidebar_group,
hidden, and config. For package @declarative-agents/ui-kit, export
names the kit panel manifest id and the shell mounts it without an
application registry entry; any other package requires a registry
entry keyed by the panel id. routes remains valid for routes that
are not panels. A file without version is version 1 and stays valid.
trace_backend.query_path, when set, shall be an absolute path ending
in /query/traces/{trace_id}; the prefix before that suffix is the
same-origin trace backend of R2.4.

<!-- field: evidence -->
#### Evidence

magefiles/uiyaml accepts every current ui.yaml unchanged and a
version 2 fixture with every new field, and the kit JSON Schema
applications/ui-kit/schema/ui.v2.schema.json agrees with it.

<!-- req: R7.2 -->
### R7.2

The validator shall reject duplicate panel ids, route collisions, and
a panel whose sidebar_group is not declared under sidebar.groups.

<!-- field: evidence -->
#### Evidence

Negative validator tests cover each rejection.

<!-- req: R7.3 -->
### R7.3

Version 2 shall remain renderable by plain Helm templating: no YAML
anchors, no multi-document files, and lists of flat maps only.

<!-- field: evidence -->
#### Evidence

A test renders a version 2 ui.yaml from chart values and validates
the result.

<!-- req: R7.4 -->
### R7.4

The routes.ts cross-check test shall be retired once an application
runs on the kit shell. It is replaced by a check that every ui.yaml
panel id has a registry entry.

<!-- field: evidence -->
#### Evidence

chatbot-mesh has no routes.ts and its magefile tests assert the
registry coverage property.

<!-- req: R8 -->
## R8 - UI kit ownership and distribution

<!-- req: R8.1 -->
### R8.1

applications/ui-kit shall own the npm package
@declarative-agents/ui-kit: design tokens, client, contract fixtures,
shell, hooks, and kit panels. The kit is a library with React 19 as a
peer dependency and ships no committed dist.

<!-- field: evidence -->
#### Evidence

The package builds and tests standalone under its mage targets.

<!-- req: R8.2 -->
### R8.2

The kit shall be the canonical home of the shared design tokens. This
transfers the token ownership that srd003 R5.3 and R7.1 assigned to
applications/catalog/ui/design-tokens.css; srd003 carries the amended
text.

<!-- field: evidence -->
#### Evidence

Every token-claiming UI imports @declarative-agents/ui-kit/tokens.css,
and the drift test finds no other token block.

<!-- req: R8.3 -->
### R8.3

UIs inside this repository shall depend on the kit through a file:
dependency. External repositories shall depend on the npm pack
tarball attached to the GitHub release tagged ui-kit/vX.Y.Z. No
package registry, credential, or environment variable is involved.

<!-- field: evidence -->
#### Evidence

The pack target produces the tarball, and the release target checks
the tag against package.json before printing the release command.

<!-- req: R9 -->
## R9 - Embedded bundles

<!-- req: R9.1 -->
### R9.1

A REST static_assets binding may select a named bundle compiled into
agent-core in place of a filesystem root. root and bundle are mutually
exclusive, and an unknown bundle name fails at load. The binding
semantics are owned by agent-core srd029-rest-server-tools.

<!-- field: evidence -->
#### Evidence

REST server tests cover bundle serving, SPA fallback, and both load
failures.

<!-- req: R9.2 -->
### R9.2

The binding's config map shall be served as ui-config.json beneath
the mount. The bundle reads its title, fleet endpoint, and monitored
agents from that file, so a declaration configures it without a
rebuild or environment variable.

<!-- field: evidence -->
#### Evidence

A REST server test reads ui-config.json and compares it with the
declared map.

<!-- req: R9.3 -->
### R9.3

An embedded UI remains a manifest UI asset. Its source is
embedded:<bundle> and its shared-token policy is the kit's.

<!-- field: evidence -->
#### Evidence

Manifest validation accepts the embedded source for the chatbot-mesh
observer.

<!-- req: R9.4 -->
### R9.4

The first bundle is observer. An observer agent that serves it shall
carry no UI files in its directory.

<!-- field: evidence -->
#### Evidence

applications/chatbot-mesh/agents/observer contains declaration files
only.

<!-- field: non_goals -->
## Non-goals

- Runtime panel loading or module federation.
- Unifying domain panels such as cohere-demo intsum, audit, export, and provenance or agentic-wiki-mesh corpus, document, and staleness.
- A desktop or Electron shell; R6 keeps that path open without building it.
- Changing monitor, fleet, or trace response schemas; their owners keep them.
- Publishing the kit to a package registry.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 | traces: ["R1.1","R1.2","R1.3","R3.1"] -->
### AC1

The kit exports the closed endpoint list and contract version, every
published panel names only listed endpoints and reads only owner-defined
fields, and every current UI URL is classified as contract or domain.

<!-- ac: AC2 | traces: ["R2.1","R2.2","R2.3","R2.4"] -->
### AC2

Cross-agent reads use /monitor-proxy only, a proxy 404 hides the dependent
panel, trace reads follow the declared trace backend, and no REST
definition under applications/ declares /trace-proxy.

<!-- ac: AC3 | traces: ["R3.2","R3.3"] -->
### AC3

Recorded fixtures exist for every contract endpoint, every client and
panel is tested against them, and owner changes update fixtures.

<!-- ac: AC4 | traces: ["R4.1","R4.2","R4.3","R4.4"] -->
### AC4

Each kit panel exports a manifest within the contract, imports nothing
from an application, and performs I/O only through the client; the panel
set matches R4.4.

<!-- ac: AC5 | traces: ["R5.1","R5.2","R5.3"] -->
### AC5

An application bundle derives routes and sidebar from ui.yaml at build
time, fails the build for an unregistered panel, and keeps existing deep
links.

<!-- ac: AC6 | traces: ["R6.1","R6.2"] -->
### AC6

Every fetch and EventSource goes through the injectable client obtained
from shell context.

<!-- ac: AC7 | traces: ["R7.1","R7.2","R7.3","R7.4"] -->
### AC7

ui.yaml version 2 validates with its rejections, stays Helm-renderable,
accepts version 1 files unchanged, and replaces the routes.ts
cross-check.

<!-- ac: AC8 | traces: ["R8.1","R8.2","R8.3"] -->
### AC8

applications/ui-kit owns the package and the design tokens, in-repository
UIs consume it by file dependency, and the release tarball is the only
external distribution.

<!-- ac: AC9 | traces: ["R9.1","R9.2","R9.3","R9.4"] -->
### AC9

The REST tool serves a named embedded bundle and its declared config,
the embedded observer is a manifest UI asset, and the chatbot-mesh
observer directory holds declarations only.

<!-- field: applicability -->
## Applicability

```yaml
baseline: |
    R1 through R6 apply to every panel published by the UI kit and to every
    application shell built from it. Domain panels an application keeps to
    itself follow R6 and need not satisfy R1 or R4.
capability_gates:
    ui: The whole SRD applies only to applications or actors that ship UI assets.
    embedded_bundle: R9 applies only to REST definitions that select a named bundle.
```

<!-- field: use_case_traces -->
## Use case traces

- rel15.0-uc001-declarative-ux

<!-- field: test_suite_traces -->
## Test suite traces

- test-rel15.0-declarative-ux
