---
id: rel19.0-uc001-persistent-application-platform
title: Operate Two Namespaced Applications Without Losing Retained Telemetry
status: implemented
test_suite: test-rel19.0-persistent-application-platform
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: release -->
## Release

19.0

<!-- field: summary -->
## Summary

An operator creates or reuses one persistent `da-platform`, then deploys two
manifest-driven applications into separate namespaces. Shared Traefik routes
external HTTP while each application's agents use Service DNS and export OTLP
to that application's collector. Each collector records stable object
identity in its WAL before uploading immutable trace and metric batches to
its application bucket.

The operator removes one application. Its compute and query API disappear,
but its bucket remains and the second application is unaffected. Redeploying
the first collector against the same bucket restores its historical query
view. Platform down and recreation also retain local objects; only explicit
application purge or platform reset deletes the selected data.

<!-- field: actor -->
## Actor

An application operator using the root Mage platform and application verbs.

<!-- field: trigger -->
## Trigger

The operator needs to run two applications concurrently on one local or
cloud platform and preserve their telemetry independently of compute
teardown.

<!-- field: preconditions -->
## Preconditions

- The platform implementation satisfies ENG01 locally or ENG08 in cloud.
- Both applications have validated `agents/application.yaml` manifests and Helm releases.
- Each application has a distinct namespace, collector identity, and bucket binding.
- Local fake GCS has a stable host-backed filesystem path, or cloud bucket identity and IAM are ready.

<!-- field: flow -->
## Flow

- F1: The operator runs `platform:up`. The runner creates or reuses
  `da-platform`, observes Traefik, object storage, metrics-server, storage
  provisioning, and data-path conformance, and leaves the platform running.
  - traces:
    - srd008-persistent-application-platform R1.1
    - srd008-persistent-application-platform R2.1
    - srd008-persistent-application-platform R5.1
- F2: The operator runs `app:up` for applications A and B. The reusable runner
  derives one release and stable namespace per manifest and binds each
  application collector to its own bucket.
  - traces:
    - srd008-persistent-application-platform R1.2
    - srd008-persistent-application-platform R2.2
    - srd008-persistent-application-platform R2.3
    - srd008-persistent-application-platform R6.1
    - srd008-persistent-application-platform R9.1
- F3: External requests enter shared Traefik and route to the selected
  application. Internal REST uses Service DNS. Agent Core exports OTLP
  directly to the application collector; object and vector clients use
  their declared endpoints directly.
  - traces:
    - srd008-persistent-application-platform R3.1
    - srd008-persistent-application-platform R8.1
- F4: Each collector records the final key, stable batch identity, checksum,
  and pending state in its local WAL, then uploads a versioned immutable
  envelope to its application bucket. Retry reuses the identity and a key
  conflict fails visibly.
  - traces:
    - srd008-persistent-application-platform R3.2
    - srd008-persistent-application-platform R6.2
    - srd008-persistent-application-platform R6.3
    - srd008-persistent-application-platform R6.4
- F5: The operator queries each application's collector and receives only that
  application's retained trace and metric history. No platform query
  service or cross-application index participates.
  - traces:
    - srd008-persistent-application-platform R3.3
    - srd008-persistent-application-platform R7.1
- F6: The operator runs `app:down` for application A. A's release and namespace
  disappear, its query API is unavailable, and its bucket remains.
  Application B continues serving and querying its own bucket.
  - traces:
    - srd008-persistent-application-platform R1.3
    - srd008-persistent-application-platform R4.1
    - srd008-persistent-application-platform R7.2
    - srd008-persistent-application-platform R8.2
- F7: The operator runs `app:up` for application A again. Its collector binds
  the same bucket and restores query access to the same historical batch
  identities and checksums.
  - traces:
    - srd008-persistent-application-platform R7.3
- F8: The operator runs `platform:down` and later `platform:up`. Local fake GCS
  remounts the stable host path, or the cloud cluster reconnects to
  surviving buckets. Both applications can recover their retained data.
  - traces:
    - srd008-persistent-application-platform R4.3
    - srd008-persistent-application-platform R4.4
    - srd008-persistent-application-platform R5.2
- F9: When deletion is intended, the operator selects either `app:purge` for
  one application's data or `platform:reset` for the local platform path.
  Neither action is implied by ordinary down.
  - traces:
    - srd008-persistent-application-platform R4.2
    - srd008-persistent-application-platform R4.3
- F10: Cloud workload identity and bucket IAM enforce application isolation.
  The local scenario reports configuration and lifecycle separation without
  claiming adversarial tenant authorization.
  - traces:
    - srd008-persistent-application-platform R5.3
- F11: First-party and downstream applications expose the same thin application
  verbs over the released application-rig package while retaining their own
  values and assertions.
  - traces:
    - srd008-persistent-application-platform R9.2
    - srd008-persistent-application-platform R9.3
- F12: A cluster-mutating integration scenario selects ENG01's owned-cluster
  exception, while a hermetic integration run may dual-export to host OTLP
  evidence. Neither exception changes the persistent platform contract or
  becomes its query authority.
  - traces:
    - srd008-persistent-application-platform R8.3

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | docs/engineering/eng01-kind-test-demo-rig.yaml | Local platform lifecycle, canonical modes, persistence, routing, and verbs. |
| T2 | docs/engineering/eng08-gcp-demo-rig.yaml | Cloud realization, bucket IAM, workload identity, and durable cluster-independent data. |
| T3 | applications/docs/ARCHITECTURE.yaml | Ownership, protocol paths, failure boundaries, and the reusable runner component. |
| T4 | applications/docs/pattern-language.yaml | Persistent platform and retained-data lifecycle construction patterns. |
| T5 | srd008-persistent-application-platform AC1 | Ownership and complete trace closure. |
| T6 | srd008-persistent-application-platform AC2 | Distinct protocol paths. |
| T7 | srd008-persistent-application-platform AC3 | Concurrent application and platform lifecycle. |
| T8 | srd008-persistent-application-platform AC4 | Local and cloud equivalence and threat boundaries. |
| T9 | srd008-persistent-application-platform AC5 | Durable object and WAL contract. |
| T10 | srd008-persistent-application-platform AC6 | Query unavailability and recovery. |
| T11 | srd008-persistent-application-platform AC7 | Rig modes and reusable package. |
| T12 | test-rel19.0-persistent-application-platform | Structural and planned executable evidence for the flow. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Platform, application-compute, and retained-data owners have disjoint
destructive authority.

<!-- field: traces -->
#### Traces

- srd008-persistent-application-platform AC1

<!-- sc: S2 -->
### S2

External HTTP, internal REST, OTLP, object storage, and vector queries
follow their declared paths without routing internal traffic through
Traefik.

<!-- field: traces -->
#### Traces

- srd008-persistent-application-platform AC2

<!-- sc: S3 -->
### S3

Two applications run concurrently and every down, redeploy, purge, and
reset transition preserves or deletes exactly the documented resources.

<!-- field: traces -->
#### Traces

- srd008-persistent-application-platform AC3

<!-- sc: S4 -->
### S4

Local and cloud use the same application declarations while reporting
their different identity and authorization guarantees truthfully.

<!-- field: traces -->
#### Traces

- srd008-persistent-application-platform AC4

<!-- sc: S5 -->
### S5

WAL state fixes immutable object identity before upload and every retained
object carries the required versioned envelope.

<!-- field: traces -->
#### Traces

- srd008-persistent-application-platform AC5

<!-- sc: S6 -->
### S6

Application down leaves history retained but not live-queryable, and
collector redeployment restores scoped query access.

<!-- field: traces -->
#### Traces

- srd008-persistent-application-platform AC6

<!-- sc: S7 -->
### S7

Local and cloud guidelines expose one persistent platform contract and
applications consume one reusable lifecycle package through thin verbs.

<!-- field: traces -->
#### Traces

- srd008-persistent-application-platform AC7

<!-- field: out_of_scope -->
## Out of scope

- Implementing the platform, collector persistence, object query, charts, or lifecycle runner.
- A global query service, authenticated local multi-tenancy, or shared vector storage.
- Replacing application collectors or routing their internal traffic through Traefik.
