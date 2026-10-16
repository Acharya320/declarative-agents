---
id: srd008-persistent-application-platform
title: Persistent Platform and Namespaced Application Contract
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

The repository has two incomplete operating models. Namespaced integration
scenarios share `da-platform`, but interactive demos still own one cluster
per application and persistent integration telemetry lives on the host.
Neither model defines the desired durable platform: one ingress and object
service, stable application namespaces, application-owned collector
gateways, retained telemetry after compute teardown, and explicit destructive
lifecycle boundaries.

Without one contract, implementation work must improvise bucket tenancy,
local durability, object identity, query availability, platform reset, and
the package boundary downstream repositories import. Those choices affect
isolation and data loss and therefore belong in the specification.

<!-- field: goals -->
## Goals

- G1: Separate platform, application-compute, and retained-data ownership.
- G2: Run concurrent applications on one persistent local or cloud platform without changing application declarations.
- G3: Retain immutable application telemetry through collector, application, and platform compute lifecycles.
- G4: Make every destructive lifecycle action explicit and separately authorized.
- G5: Publish one reusable application-rig package consumed through thin Mage verbs.

<!-- field: applicability -->
## Applicability

```yaml
baseline: |
    R1 through R9 apply to applications that declare `helm_managed` and use the
    shared application platform. Cluster-mutating integration scenarios remain
    governed by ENG01's owned-cluster exception.
capability_gates:
    local_platform: R2, R4, R5, R7, and R8 local clauses.
    cloud_platform: R2, R4, R5, R7, and R8 cloud clauses.
    durable_collector: R3, R6, and R7 collector clauses.
```

<!-- req: R1 -->
## R1 - Ownership and authority are separate

<!-- field: applicability -->
### Applicability

helm_managed

<!-- req: R1.1 -->
### R1.1

The platform authority shall create, upgrade, diagnose, stop, and
reset shared ingress, object-storage endpoint, metrics-server, and
storage provisioning. It shall not own application releases or
application data-retention policy.

<!-- field: evidence -->
#### Evidence

The ownership matrix names one platform authority and excludes
application release and bucket-purge authority from it.

<!-- req: R1.2 -->
### R1.2

Each application authority shall create, upgrade, diagnose, and stop
one manifest-driven Helm release in one stable namespace, including
its agents, services, routes, collector gateway, and bucket binding.

<!-- field: evidence -->
#### Evidence

Two applications can be operated independently without either
lifecycle command selecting the other's namespace or release.

<!-- req: R1.3 -->
### R1.3

Retained data shall have an application-scoped authority independent
of compute ownership. Only `app:purge` may delete an application's
bucket data; neither application nor platform shutdown may do so.

<!-- field: evidence -->
#### Evidence

Teardown tests preserve objects, while a separately authorized purge
removes only the selected application's objects.

<!-- req: R2 -->
## R2 - One persistent platform hosts isolated application namespaces

<!-- field: applicability -->
### Applicability

helm_managed

<!-- req: R2.1 -->
### R2.1

`da-platform` shall provide one shared Traefik ingress, one
object-storage endpoint, metrics-server, and storage provisioning.
Platform readiness shall be observed before any application install.

<!-- field: evidence -->
#### Evidence

Platform conformance proves all four services before application
lifecycle operations proceed.

<!-- req: R2.2 -->
### R2.2

A runnable application shall occupy one stable namespace and one
manifest-derived Helm release. Two applications shall run
concurrently with distinct release, namespace, collector, route, and
bucket identities.

<!-- field: evidence -->
#### Evidence

A two-application scenario routes to both applications and observes
distinct namespace, collector, and bucket identities.

<!-- req: R2.3 -->
### R2.3

Shared platform components shall not become undeclared application
dependencies. An application manifest shall declare its platform
service bindings and retain authority over its own workloads.

<!-- field: evidence -->
#### Evidence

Application audit derives required bindings from
`agents/application.yaml` and finds no second lifecycle inventory.

<!-- req: R3 -->
## R3 - Each application retains its collector gateway

<!-- field: applicability -->
### Applicability

durable_collector

<!-- req: R3.1 -->
### R3.1

Agent Core shall instrument REST and export OTLP directly to the
collector gateway in the same application namespace. Internal REST
and OTLP traffic shall use Service DNS and shall not traverse
Traefik.

<!-- field: evidence -->
#### Evidence

Rendered routes and a live trace show external HTTP at Traefik,
internal REST over Service DNS, and OTLP at the application
collector Service.

<!-- req: R3.2 -->
### R3.2

The collector shall accept OTLP traces and metrics, write a local WAL
before remote upload, persist immutable objects, serve
application-scoped queries, and optionally relay OTLP when declared.

<!-- field: evidence -->
#### Evidence

Collector tests prove WAL-first upload, immutable object persistence,
scoped query results, and relay as an independent optional path.

<!-- req: R3.3 -->
### R3.3

A collector shall query only the application bucket bound to its
release. No permanent platform query service or cross-application
telemetry index is part of this release.

<!-- field: evidence -->
#### Evidence

Query tests refuse or return no records from another application,
and platform manifests contain no global query workload.

<!-- req: R4 -->
## R4 - Lifecycle verbs preserve data unless they say purge or reset

<!-- field: applicability -->
### Applicability

helm_managed

<!-- req: R4.1 -->
### R4.1

`app:up`, `app:status`, `app:down`, and `app:diagnose` shall operate
on one manifest-selected application. `app:down` shall remove its
compute and namespace while preserving its bucket data.

<!-- field: evidence -->
#### Evidence

Down removes the selected release and namespace, leaves both
applications' objects intact, and leaves the other application
serving.

<!-- req: R4.2 -->
### R4.2

`app:purge` shall be a separate explicit destructive operation. It
shall delete only the selected application's retained data and shall
not be an implicit part of install, upgrade, down, or platform
lifecycle.

<!-- field: evidence -->
#### Evidence

Purge requires explicit application identity and leaves another
application's bucket unchanged.

<!-- req: R4.3 -->
### R4.3

`platform:down` shall stop platform compute while preserving local
object data. `platform:reset` shall be the only local platform verb
that deletes the stable object-service host path.

<!-- field: evidence -->
#### Evidence

Platform recreation after down retains objects; recreation after
reset starts with an empty local object service.

<!-- req: R4.4 -->
### R4.4

Cloud buckets shall outlive cluster teardown. A cloud platform
teardown or recreation shall not delete application buckets unless a
separately authorized application purge selects one.

<!-- field: evidence -->
#### Evidence

Cloud lifecycle plans contain no bucket delete in ordinary down and
retain data across cluster recreation.

<!-- req: R5 -->
## R5 - Local and cloud platforms expose one application contract

<!-- field: applicability -->
### Applicability

helm_managed

<!-- req: R5.1 -->
### R5.1

Local kind shall run fake GCS with a filesystem backend mounted from
one stable host path. The path shall survive pod and
`da-platform` recreation.

<!-- field: evidence -->
#### Evidence

A platform down/up cycle recreates fake GCS against the same host
path and reads previously written objects.

<!-- req: R5.2 -->
### R5.2

Cloud platforms shall provide GCS or S3 through ambient workload
identity and application-scoped bucket permissions. Local and cloud
differences shall be confined to bootstrap, identity, and values
overlays.

<!-- field: evidence -->
#### Evidence

The same application and collector declarations render for local and
cloud; only platform and values inputs differ.

<!-- req: R5.3 -->
### R5.3

Cloud isolation shall be enforced by workload identity and bucket
IAM. Local fake GCS shall claim configuration and lifecycle isolation
only, not resistance to a malicious pod, because it has no equivalent
tenant authorization.

<!-- field: evidence -->
#### Evidence

The threat model and tests distinguish enforced cloud denial from
non-adversarial local separation.

<!-- req: R6 -->
## R6 - Immutable objects have stable identity before upload

<!-- field: applicability -->
### Applicability

durable_collector

<!-- req: R6.1 -->
### R6.1

Each application shall receive its own bucket and may organize
objects with prefixes inside that bucket. A shared endpoint shall
not imply one undifferentiated bucket.

<!-- field: evidence -->
#### Evidence

Bootstrap creates or binds distinct application buckets and the
collector configuration names exactly one.

<!-- req: R6.2 -->
### R6.2

Before the first remote attempt, the local WAL shall record the final
object key, stable batch identity, checksum, and upload state.
Retries shall reuse that key and identity.

<!-- field: evidence -->
#### Evidence

Crash and retry tests produce one immutable object identity rather
than a new key for each attempt.

<!-- req: R6.3 -->
### R6.3

The versioned object envelope shall minimally name schema version,
signal, application, namespace, run, collector instance, stable
batch identity, receive time, payload format, and checksum.

<!-- field: evidence -->
#### Evidence

Schema validation rejects an object missing any required field and
accepts trace and metric envelopes.

<!-- req: R6.4 -->
### R6.4

An uploaded key shall be immutable. A checksum or payload conflict
at an existing key shall be a collector fault and shall never
overwrite retained evidence silently.

<!-- field: evidence -->
#### Evidence

Duplicate identical upload is idempotent; conflicting content at the
same key is reported and remains unchanged.

<!-- req: R7 -->
## R7 - Query availability follows the application collector

<!-- field: applicability -->
### Applicability

durable_collector

<!-- req: R7.1 -->
### R7.1

While the application collector is running, its REST query surface
shall combine retained object data with any committed local state and
shall scope every response to that application.

<!-- field: evidence -->
#### Evidence

The application query API returns its pre-restart telemetry and no
other application's records.

<!-- req: R7.2 -->
### R7.2

After `app:down`, retained bucket data shall not be described as a
live query service. With the collector absent, its REST query API is
unavailable.

<!-- field: evidence -->
#### Evidence

The lifecycle scenario observes retained objects and an unavailable
application query endpoint after down.

<!-- req: R7.3 -->
### R7.3

Redeploying the collector against the same application bucket shall
restore query access to retained history without changing object
identity.

<!-- field: evidence -->
#### Evidence

Down/up returns the same historical batch identities and checksums
through the redeployed collector.

<!-- req: R8 -->
## R8 - Network paths and failures remain attributable

<!-- field: applicability -->
### Applicability

helm_managed

<!-- req: R8.1 -->
### R8.1

External HTTP shall enter through shared Traefik and route by
application host or path to an application Service. Internal REST,
OTLP, object-storage, and vector-query traffic shall use their
declared Service or provider endpoints directly.

<!-- field: evidence -->
#### Evidence

The architecture diagram and rendered configuration distinguish all
five paths and no internal path depends on Traefik.

<!-- req: R8.2 -->
### R8.2

Failure of one application release, collector, route, bucket binding,
or purge shall not stop another application or delete its data.
Failure of a shared platform service shall be reported as a platform
failure affecting all dependent applications.

<!-- field: evidence -->
#### Evidence

Fault scenarios classify application-local and shared-platform
failures and preserve the unaffected application's compute and data.

<!-- req: R8.3 -->
### R8.3

Persistent namespaced `da-platform` shall be the canonical
application and interactive-demo mode. ENG01 owned clusters shall
remain only for cluster-mutating scenarios. Host OTLP ingress shall
remain an optional hermetic integration-evidence mode and shall not
be the durable application query authority.

<!-- field: evidence -->
#### Evidence

ENG01 names the modes and routes each scenario without presenting
dedicated demo clusters or host collection as competing canonical
platform designs.

<!-- req: R9 -->
## R9 - One reusable application-rig package owns lifecycle mechanics

<!-- field: applicability -->
### Applicability

helm_managed

<!-- req: R9.1 -->
### R9.1

The root `magefiles` module shall publish a package such as
`magefiles/apprig` that accepts an application manifest, platform
binding, command runner, and evidence destination and implements the
shared application lifecycle.

<!-- field: evidence -->
#### Evidence

At least two in-repository applications call the package without
copying lifecycle implementation.

<!-- req: R9.2 -->
### R9.2

Downstream repositories shall import a released root Mage module and
expose thin `app:up|status|down|diagnose|purge` verbs. They shall not
copy platform, collector, bucket, or namespace lifecycle code.

<!-- field: evidence -->
#### Evidence

A downstream migration pins a release, contains only thin verb
wiring, and passes the same lifecycle contract.

<!-- req: R9.3 -->
### R9.3

Application-specific assertions and values shall remain with the
application. The reusable package shall own orchestration mechanics,
not application behavior or a second composition inventory.

<!-- field: evidence -->
#### Evidence

Package APIs receive manifest-derived inputs and callbacks for local
assertions without naming first-party applications.

<!-- field: non_goals -->
## Non-goals

- Implementing platform services, collector storage, query code, charts, or the application runner in this release.
- A permanent global telemetry query service or cross-application telemetry index.
- Authenticated adversarial multi-tenancy for local fake GCS.
- Shared vector storage before collection authorization and embedding-identity isolation are specified.
- Replacing application collectors with one global collector.
- Routing internal REST or OTLP through Traefik.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

Structured documentation assigns platform, application-compute, and
retained-data authority and traces every requirement to an executable or
planned test case.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R2.3
- R9.3

<!-- ac: AC2 -->
### AC2

A diagram and normative text distinguish external HTTP ingress, internal
REST, OTLP, object storage, and vector-query paths.

<!-- field: traces -->
#### Traces

- R3.1
- R8.1

<!-- ac: AC3 -->
### AC3

One unambiguous lifecycle covers two concurrent applications,
application down/redeploy/purge, platform down/recreation/reset, and
cloud cluster recreation without implicit data deletion.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R4.1
- R4.2
- R4.3
- R4.4
- R8.2

<!-- ac: AC4 -->
### AC4

Local and cloud platforms preserve one declaration contract while
stating the different durability, identity, authorization, and threat
boundaries exactly.

<!-- field: traces -->
#### Traces

- R5.1
- R5.2
- R5.3

<!-- ac: AC5 -->
### AC5

The collector contract fixes per-application buckets, WAL-before-upload
ordering, immutable object identity, the minimum versioned envelope, and
conflict behavior.

<!-- field: traces -->
#### Traces

- R3.2
- R3.3
- R6.1
- R6.2
- R6.3
- R6.4

<!-- ac: AC6 -->
### AC6

Retained history is explicitly unavailable as a live query service after
application down and becomes queryable through the application collector
after redeployment against the same bucket.

<!-- field: traces -->
#### Traces

- R7.1
- R7.2
- R7.3

<!-- ac: AC7 -->
### AC7

ENG01 and ENG08 name the canonical platform, the owned-cluster and host
evidence exceptions, cloud realization, failure boundaries, and the
reusable application-rig integration contract.

<!-- field: traces -->
#### Traces

- R8.3
- R9.1
- R9.2

<!-- field: use_case_traces -->
## Use case traces

- rel19.0-uc001-persistent-application-platform

<!-- field: test_suite_traces -->
## Test suite traces

- test-rel19.0-persistent-application-platform

<!-- field: related_documents -->
## Related documents

- docs/engineering/eng01-kind-test-demo-rig.md
- docs/engineering/eng08-gcp-demo-rig.md
- agent-core/docs/specs/software-requirements/srd008-telemetry.md
- agent-core/docs/specs/software-requirements/srd042-otlp-receiver-tools.md
- agent-core/docs/specs/software-requirements/srd044-spool-query-tools.md
- agent-core/docs/specs/software-requirements/srd059-objectstore-tools.md
