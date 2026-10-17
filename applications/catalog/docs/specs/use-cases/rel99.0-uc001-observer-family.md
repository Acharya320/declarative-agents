---
id: rel99.0-uc001-observer-family
title: The observer is a catalog-owned family adopted across applications
status: planned
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

The fleet observer becomes a catalog-owned family generalized from the
chatbot-mesh observer: its namespace and label selector are parameterized so
each application binds it against its own release. Coding-agent and
agent-architecture adopt the family by reference, and the chatbot-mesh
observer rebinds to it, following the collector's mesh-to-catalog promotion.

<!-- field: actor -->
## Actor

Application author adopting the catalog observer family

<!-- field: trigger -->
## Trigger

An application needs a fleet-level view of its mesh and binds the catalog
observer family against its release selector instead of forking the profile.

<!-- field: flow -->
## Flow

- F1: The catalog carries the generalized observer-family SRD (srd021) with parameterized namespace and label selector defaults.
- F2: The family carries the mesh observer's topology discovery, monitor fan-in, resource usage, and served fleet UX forward as generalized requirements.
- F3: An adopting application deploys the family as a workload pointed at its own release selector and namespace, without forking the profile.
- F4: The chatbot-mesh observer rebinds to the catalog family, retaining the descriptive observer actor name.
- F5: The adoption implementation adds the bindings.catalog and deployment_realizations rows once the catalog observer profile and per-application workloads exist.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd021-observer](../software-requirements/srd021-observer.md) [R1](../software-requirements/srd021-observer.md#r1---generalized-defaults) | generalized namespace and label selector defaults, no new environment variable |
| T2 | [srd021-observer](../software-requirements/srd021-observer.md) [R3](../software-requirements/srd021-observer.md#r3---carried-capabilities) | carried topology, fan-in, resource, and served-UX capabilities |
| T3 | [srd021-observer](../software-requirements/srd021-observer.md) [R4](../software-requirements/srd021-observer.md#r4---per-application-adoption) | coding-agent, agent-architecture, and chatbot-mesh adoption |
| T4 | [srd021-observer](../software-requirements/srd021-observer.md) [R1](../software-requirements/srd021-observer.md#r1---generalized-defaults) | the chatbot-mesh observer this family generalizes, whose own SRD is srd008-observer-agent in that application corpus |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The catalog carries a generalized observer-family SRD with parameterized namespace and label selector and no new environment variable.

<!-- field: traces -->
#### Traces

- [srd021-observer](../software-requirements/srd021-observer.md) [AC1](../software-requirements/srd021-observer.md#ac1)

<!-- sc: S2 -->
### S2

The family carries topology discovery, monitor fan-in, resource usage, and the served fleet UX forward from the mesh observer.

<!-- field: traces -->
#### Traces

- [srd021-observer](../software-requirements/srd021-observer.md) [AC2](../software-requirements/srd021-observer.md#ac2)

<!-- sc: S3 -->
### S3

The specification names the adopters and defers realization rows to the adoption implementation because the inventory gate requires existing profile paths.

<!-- field: traces -->
#### Traces

- [srd021-observer](../software-requirements/srd021-observer.md) [AC3](../software-requirements/srd021-observer.md#ac3)

<!-- field: out_of_scope -->
## Out of scope

- Shipping the catalog observer profile, tools, and UX.
- Prose Editor adoption; it is downstream-owned and not an in-repository consumer.
- Trace intake or metrics scraping (collector, srd020).

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel99.0-observer](../test-suites/test-rel99.0-observer.md)
