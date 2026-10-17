---
id: tenant-isolated-vector-storage
title: Tenant-Isolated Shared Vector Storage Decision Audit
audit_date: "2026-09-22"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: scope -->
## Scope

<!-- field: issue -->
### Issue

GH-2481

<!-- field: decision -->
### Decision

Retain application-owned vector stores; do not add a shared vector service to da-platform.

<!-- field: evidence_boundary -->
### Evidence boundary

The audit covers every in-repository first-party application and searches
the accessible Nokia-Bell-Labs GitHub corpus for downstream CHROMA_URL and
chromadb consumers. Prose Editor is documented as downstream-owned, but its
source and a second downstream application are not present or discoverable
from this checkout; no claim is made about their private implementation.

<!-- field: examined -->
### Examined

- applications/chatbot-mesh/helm/templates/rag-units.yaml
- applications/chatbot-mesh/helm/values.yaml
- applications/chatbot-mesh/agents/rag-server/
- applications/catalog/agents/knowledge-manager/
- applications/coding-agent/agents/application.yaml
- applications/agent-architecture/agents/application.yaml
- applications/catalog/docs/specs/software-requirements/srd012-chroma-corpus-agents.md
- docs/engineering/eng01-kind-test-demo-rig.md

<!-- field: inventory -->
## Inventory

<!-- field: first_party -->
### First party

<!-- first_party_entry: chatbot-mesh -->
#### chatbot-mesh

<!-- field: vector_workloads -->
##### Vector workloads

One Chroma StatefulSet and PVC per declared ragUnits entry.

<!-- field: current_local_instance -->
##### Current local instance

chatbot-mesh-chatbot-mesh-rag0-chroma

<!-- field: namespace -->
##### Namespace

app-chatbot-mesh

<!-- field: collection -->
##### Collection

corpus

<!-- field: embedding_identity | form: "tree" -->
##### Embedding identity

- provider: Ollama
- model: qwen3-embedding:8b
- declared_separately_from_collection: true

<!-- field: persistence -->
##### Persistence

256Mi ReadWriteOnce PVC mounted at /data

<!-- field: credentials -->
##### Credentials

none

<!-- field: authorization -->
##### Authorization

none in Chroma; local separation is workload/Service/PVC topology, not authenticated tenancy

<!-- field: lifecycle -->
##### Lifecycle

<!-- field: create -->
###### Create

application Helm release

<!-- field: ingest -->
###### Ingest

catalog corpus-ingest profile through application-owned app:seed

<!-- field: query -->
###### Query

application rag-server

<!-- field: backup -->
###### Backup

no declared production backup authority; benchmark used a read-only tar stream

<!-- field: retention -->
###### Retention

PVC follows application namespace deletion today

<!-- field: purge -->
###### Purge

no vector-specific portable purge contract

<!-- first_party_entry: coding-agent -->
#### coding-agent

<!-- field: vector_workloads -->
##### Vector workloads

none

<!-- field: retrieval_policy -->
##### Retrieval policy

none

<!-- first_party_entry: agent-architecture -->
#### agent-architecture

<!-- field: vector_workloads -->
##### Vector workloads

none

<!-- field: retrieval_policy -->
##### Retrieval policy

Documentation curator reads files, not a vector engine.

<!-- field: downstream | form: "typed table" -->
### Downstream

| application | availability | evidence | disposition |
| --- | --- | --- | --- |
| prose-editor | documented downstream owner; source absent from this repository and inaccessible to this audit | applications/catalog/docs/specs/software-requirements/srd021-observer.md R4.3 |  |
| second-downstream-consumer | not named in repository documentation and not discoverable through accessible organization code search |  | No implementation facts inferred. |

<!-- field: measurements -->
## Measurements

<!-- field: environment | form: "tree" -->
### Environment

- platform: da-platform kind, kindest/node v1.36.1
- architecture: darwin arm64 host, linux/arm64 node
- workload: seeded chatbot-mesh rag0 Chroma 1.5.3
- sample_time: 2026-09-22T04:53Z

<!-- field: steady_state -->
### Steady state

<!-- field: cpu_millicores | form: "scalar" -->
#### CPU millicores

3

<!-- field: memory_mib | form: "scalar" -->
#### Memory mib

11

<!-- field: pvc_capacity_mib | form: "scalar" -->
#### PVC capacity mib

256

<!-- field: data_bytes | form: "scalar" -->
#### Data bytes

188416

<!-- field: command -->
#### Command

kubectl top pod --containers and kubectl exec -- du -sb /data

<!-- field: heartbeat_latency -->
### Heartbeat latency

<!-- field: requests | form: "scalar" -->
#### Requests

50

<!-- field: min_ms | form: "scalar" -->
#### Min ms

2.437

<!-- field: p50_ms | form: "scalar" -->
#### P50 ms

3.393

<!-- field: p95_ms | form: "scalar" -->
#### P95 ms

6.420

<!-- field: max_ms | form: "scalar" -->
#### Max ms

36.088

<!-- field: command -->
#### Command

curl /api/v2/heartbeat through a local Service port-forward

<!-- field: application_query_latency -->
### Application query latency

<!-- field: requests | form: "scalar" -->
#### Requests

3

<!-- field: min_ms | form: "scalar" -->
#### Min ms

401.8

<!-- field: median_ms | form: "scalar" -->
#### Median ms

465.2

<!-- field: max_ms | form: "scalar" -->
#### Max ms

2515.7

<!-- field: command -->
#### Command

POST /api/v1/chat through shared Traefik after corpus seed

<!-- field: note -->
#### Note

End-to-end latency includes embedding, vector query, routing, and answer generation; it is not a Chroma-only measure.

<!-- field: ingest -->
### Ingest

<!-- field: corpus -->
#### Corpus

checked-in chatbot-mesh integration corpus

<!-- field: elapsed_seconds | form: "scalar" -->
#### Elapsed seconds

13.717

<!-- field: result -->
#### Result

corpus-ingest reached success and app:verify served a non-empty answer

<!-- field: restart -->
### Restart

<!-- field: delete_to_ready_seconds | form: "scalar" -->
#### Delete to ready seconds

12.202

<!-- field: persistence_check -->
#### Persistence check

app:verify passed after the StatefulSet pod replacement

<!-- field: backup_stream -->
### Backup stream

<!-- field: bytes | form: "scalar" -->
#### Bytes

194560

<!-- field: elapsed_seconds | form: "scalar" -->
#### Elapsed seconds

0.091

<!-- field: method -->
#### Method

tar stream of /data through kubectl exec; no backup authority or schedule is implied

<!-- field: application_removal -->
### Application removal

<!-- field: app_down_seconds | form: "scalar" -->
#### App down seconds

112.562

<!-- field: result -->
#### Result

namespace and Chroma PVC removed; application telemetry bucket remained separate and retained

<!-- field: implication -->
#### Implication

Vector retention across app:down is not specified today and must not be silently changed by a shared service.

<!-- field: candidates -->
## Candidates

| id | design | authentication | authorization | embedding_identity | quotas | retention | backup | local_cloud_parity | migration_cost | result |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| per-application | Keep Chroma (or an application-chosen adapter) inside each application namespace and lifecycle. | No shared tenant endpoint; cloud deployments can add namespace NetworkPolicy and application-specific credentials without inventing cross-tenant Chroma auth. | Workload and storage ownership are distinct per application; local kind claims isolation of configuration/lifecycle, not adversarial security. | Application declarations keep provider, model, collection, and corpus independent. | Kubernetes requests/limits and PVC size are application-scoped. | Remains application-owned and must be specified before changing app:down behavior. | Application-owned; current Chroma path has no declared production backup authority. | Same application topology; cloud hardening remains application-specific. | none | recommended |
| shared-chroma-direct | One platform Chroma endpoint and collections named by application. | Current deployment has none. | Collection or tenant naming is not an access-control boundary. | Compatible dimensions/models can still be mixed by a guessed or misrouted collection. | Shared process/PVC has no application enforcement in the current contract. | Platform and application authorities become conflated. | One failure/restore domain for every application. | Cannot honestly claim cloud tenant isolation. | medium | rejected |
| authenticated-chroma-gateway | A platform gateway authenticates workload identity, maps immutable application/corpus/embedding coordinates, enforces quotas, and hides Chroma. | feasible at gateway | feasible only if every request and administrative path is forced through the gateway | Gateway schema can enforce dimensions/provider/model version. | Requires new accounting and admission controls. | Requires new application-data authority and purge/backup contracts. | Gateway metadata and Chroma data need one consistent snapshot protocol. | Local can prove routing/identity but not workload-identity denial. | high; new platform service, catalog API, credentials, and all consumer adapters | rejected for this release; reconsider only with measured multi-application demand |
| postgres-pgvector | Shared PostgreSQL/pgvector with database roles and row-level security. | native database credentials or workload-identity proxy | roles and RLS can enforce application/corpus rows | explicit relational keys and dimension-constrained columns | database roles plus provider resource controls | row/schema deletion can be application-scoped but needs the same purge authority decision | mature database snapshot/restore, with a larger shared failure domain | PostgreSQL is portable; cloud identity/proxy remains provider-specific | high; rewrite Chroma adapters and revalidate ranking/query behavior | rejected; capability exceeds current demand |

<!-- field: findings | form: "typed table" -->
## Findings

| id | severity | evidence | finding |
| --- | --- | --- | --- |
| F1 | decision | No auth/token/tenant configuration exists in rag-units.yaml; rag-server receives only CHROMA_URL, CHROMA_HOST, collection, and embedding model. | Current Chroma cannot provide authenticated application isolation in the shipped topology. |
| F2 | decision | A direct shared Chroma candidate has no request identity to compare against the selected collection. | Collection names and Kubernetes namespaces are not sufficient authorization for one shared endpoint. |
| F3 | measured |  | The current single seeded Chroma workload is small (3m CPU, 11Mi memory, 188416 data bytes), so consolidation has no demonstrated capacity benefit. |
| F4 | lifecycle |  | Sharing would introduce unresolved backup, retention, purge, quota, and noisy-neighbor authorities that do not exist in the current per-application contract. |
| F5 | compatibility |  | Retrieval selection, ranking, prompting, and answer policy can remain unchanged only behind an adapter; a direct engine migration is not behavior-neutral. |
| F6 | scope |  | Only Chatbot Mesh currently deploys a vector engine among the three first-party applications; there is no measured cross-application consolidation target. |

<!-- field: threat_model | form: "typed table" -->
## Threat model

| threat | per_application_disposition | shared_disposition | disposition |
| --- | --- | --- | --- |
| cross-application read/write or collection-name guessing | No shared endpoint; enforce cloud network/credential isolation separately. | Direct Chroma rejected; a gateway or RLS-capable engine would be mandatory. |  |
| credential theft |  |  | Scope credentials to one application and rotate at the platform binding; no shared administrator credential in application pods. |
| incompatible vectors |  |  | Keep corpus, collection, provider, model version, and dimensions as independent immutable coordinates. |
| purge mistake |  |  | Vector purge remains application authority and is not inferred from namespace or collection naming. |
| noisy-neighbor exhaustion |  |  | Per-application requests/limits and PVCs preserve attribution; shared candidates require explicit quotas before adoption. |

<!-- field: authority -->
## Authority

<!-- field: create -->
### Create

application lifecycle and chart

<!-- field: ingest -->
### Ingest

application-owned workflow using catalog corpus-ingest capability

<!-- field: query -->
### Query

application rag-server; platform never chooses retrieval/ranking/prompt policy

<!-- field: backup -->
### Backup

application data authority; currently an identified gap

<!-- field: migration -->
### Migration

application owner with embedding-identity compatibility proof

<!-- field: retention -->
### Retention

application data authority

<!-- field: purge -->
### Purge

separately approved application purge authority; no vector deletion fallback

<!-- field: recommendation -->
## Recommendation

<!-- field: decision -->
### Decision

reject-shared-vector-service

<!-- field: rationale -->
### Rationale

There is one measured first-party vector workload, no authenticated Chroma
tenant boundary, and no demonstrated resource pressure. A direct shared
Chroma leaks authority; a gateway or pgvector migration adds substantial
control-plane and migration cost before demand exists. Retaining
application-owned stores preserves behavior and failure isolation.

<!-- field: revisit_when -->
### Revisit when

- At least two applications have measured concurrent vector workloads.
- A stable workload identity and authorization contract exists.
- Embedding identity, quotas, retention, backup, migration, and purge are specified and executable.
- A benchmark demonstrates a material cost or operability advantage over per-application stores.

<!-- field: epic_effect -->
### Epic effect

Record per-application vector storage as canonical for GH-2474; file no shared-service implementation issues.

<!-- field: rejected_candidates -->
## Rejected candidates

- shared Chroma selected by collection naming alone
- namespace naming presented as authenticated tenant isolation
- platform-wide vector service without application-scoped retention and purge authority
- engine migration that changes retrieval policy as an incidental platform concern

<!-- field: references -->
## References

- applications/chatbot-mesh/helm/templates/rag-units.yaml
- applications/chatbot-mesh/agents/rag-server/
- applications/catalog/agents/knowledge-manager/
- applications/chatbot-mesh/docs/specs/software-requirements/
- docs/engineering/eng01-kind-test-demo-rig.md
- applications/docs/specs/software-requirements/srd008-persistent-application-platform.md
- GH-2481
