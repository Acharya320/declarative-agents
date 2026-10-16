---
id: test-rel19.0-persistent-application-platform
title: Release 19.0 Persistent Application Platform Test Suite
release: "19.0"
status: planned
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: test_case_count -->
## Test case count

```yaml
8
```

<!-- field: overview -->
## Overview

TC1 is the executable structural gate for the Release 19.0 specification.
TC2 through TC8 are the implemented structural and live evidence delivered
under GH-2474. Together the cases cover ownership, the five network paths,
concurrent application lifecycle, local/cloud declaration equivalence, WAL
and immutable objects, query restoration, failure isolation, explicit purge,
and the reusable application-rig boundary. Authenticated GKE execution remains
dependency-gated where gcloud/project credentials are unavailable; the exact
IAM policy fixture is executable without cloud access.

<!-- field: traces -->
## Traces

- rel19.0-uc001-persistent-application-platform

<!-- field: srds -->
## Srds

- srd008-persistent-application-platform

<!-- field: preconditions -->
## Preconditions

- Run commands from the repository root.
- Python 3 with PyYAML, Go, Mage, Helm, kubectl, kind, and Docker are available for applicable cases.
- Cloud cases require an authenticated project and the ENG08 workload-identity prerequisites.
- Cloud live-denial evidence requires gcloud and an authenticated project; the policy fixture is the portable gate.

<!-- field: test_cases -->
## Test cases

<!-- case: Release 19.0 documents parse and trace reciprocally | status: implemented -->
### Release 19.0 documents parse and trace reciprocally

<!-- field: id -->
#### Id

TC1

<!-- field: use_case -->
#### Use case

rel19.0-uc001-persistent-application-platform

<!-- field: evidence_type -->
#### Evidence type

executable_structural

<!-- field: implemented_by -->
#### Implemented by

GH-2476

<!-- field: command -->
#### Command

```yaml
|
    python3 - <<'PY'
    from pathlib import Path
    import subprocess
    import yaml

    root = Path("applications/docs")
    paths = {
        "srd": root / "specs/software-requirements/srd008-persistent-application-platform.md",
        "use_case": root / "specs/use-cases/rel19.0-uc001-persistent-application-platform.md",
        "suite": root / "specs/test-suites/test-rel19.0-persistent-application-platform.md",
        "index": root / "SPECIFICATIONS.md",
        "roadmap": root / "road-map.md",
        "architecture": root / "ARCHITECTURE.md",
        "pattern": root / "pattern-language.md",
        "eng01": Path("docs/engineering/eng01-kind-test-demo-rig.md"),
        "eng08": Path("docs/engineering/eng08-gcp-demo-rig.md"),
    }
    def load(path):
        # A corpus document is Markdown; mage specYAML prints the YAML tree
        # its declaration parses it into (GH-2627).
        out = subprocess.run(["mage", "specYAML", str(path)], cwd="agent-core",
                             check=True, capture_output=True, text=True).stdout
        return yaml.safe_load(out)

    docs = {name: load(path) for name, path in paths.items()}
    required = {
        "srd": {"id", "title", "problem", "goals", "requirements",
                "non_goals", "acceptance_criteria"},
        "use_case": {"id", "title", "summary", "actor", "trigger", "flow",
                     "touchpoints", "success_criteria", "out_of_scope"},
        "suite": {"id", "title", "release", "traces", "preconditions",
                  "test_cases"},
        "architecture": {"id", "title", "overview", "interfaces",
                         "components", "design_decisions",
                         "technology_choices", "project_structure",
                         "implementation_status", "related_documents"},
        "eng01": {"id", "title", "introduction", "sections", "references"},
        "eng08": {"id", "title", "introduction", "sections", "references"},
    }
    for name, fields in required.items():
        assert not fields - set(docs[name]), (name, fields - set(docs[name]))

    srd, use_case, suite = docs["srd"], docs["use_case"], docs["suite"]
    assert use_case["test_suite"] == suite["id"]
    assert use_case["id"] in suite["traces"]
    assert use_case["id"] in srd["use_case_traces"]
    assert suite["id"] in srd["test_suite_traces"]

    item_ids = {
        item["id"]
        for group in srd["requirements"].values()
        for item in group["items"]
    }
    ac_traces = {trace for ac in srd["acceptance_criteria"]
                 for trace in ac["traces"]}
    assert item_ids <= ac_traces, sorted(item_ids - ac_traces)
    ac_ids = {ac["id"] for ac in srd["acceptance_criteria"]}
    case_traces = {
        trace.rsplit(" ", 1)[-1]
        for case in suite["test_cases"]
        for trace in case.get("traces", [])
        if trace.startswith(srd["id"] + " ")
    }
    assert ac_ids <= case_traces, sorted(ac_ids - case_traces)
    assert len(suite["test_cases"]) == suite["test_case_count"]

    index = docs["index"]
    assert any(e["id"] == srd["id"] for e in index["srd_index"])
    assert any(e["id"] == use_case["id"] for e in index["use_case_index"])
    assert any(e["id"] == suite["id"] for e in index["test_suite_index"])
    assert any(e["version"] == "19.0" for e in index["roadmap_summary"])
    release = next(r for r in docs["roadmap"]["releases"]
                   if r["version"] == "19.0")
    assert any(e["id"] == use_case["id"] for e in release["use_cases"])
    assert any(e["id"] == suite["id"] for e in release["test_suites"])

    eng01 = "\n".join(s["content"] for s in docs["eng01"]["sections"])
    eng08 = "\n".join(s["content"] for s in docs["eng08"]["sections"])
    for term in ("app:down", "app:purge", "platform:reset",
                 "application-specific bucket", "apprig"):
        assert term in eng01, term
    for term in ("workload identity", "bucket IAM", "application-specific bucket"):
        assert term in eng08, term
    pattern_ids = {p["id"] for p in docs["pattern"]["patterns"]}
    assert "persistent-application-platform" in pattern_ids
    print("validated Release 19.0 persistent platform trace closure")
    PY
```

<!-- field: expected -->
#### Expected

```yaml
yaml_parses: true
design_constitution_fields_present: true
requirement_to_acceptance_criterion_closure: true
reciprocal_srd_use_case_suite_traces: true
guideline_contract_terms_present: true
specification_and_roadmap_indexes_present: true
```

<!-- field: traces -->
#### Traces

- rel19.0-uc001-persistent-application-platform S1
- srd008-persistent-application-platform AC1

<!-- case: Two applications use distinct external, internal, OTLP, object, and vector paths | status: implemented -->
### Two applications use distinct external, internal, OTLP, object, and vector paths

<!-- field: id -->
#### Id

TC2

<!-- field: use_case -->
#### Use case

rel19.0-uc001-persistent-application-platform

<!-- field: evidence_type -->
#### Evidence type

integration

<!-- field: implemented_by -->
#### Implemented by

- GH-2477
- GH-2484
- GH-2482
- GH-2483
- GH-2480

<!-- field: command -->
#### Command

mage platform:up
mage app:up agent-architecture
mage app:up coding-agent
mage app:up chatbot-mesh
kubectl --context kind-da-platform get ns -l app.kubernetes.io/managed-by=apprig

<!-- field: expected -->
#### Expected

```yaml
shared_traefik_routes_both_applications: true
internal_rest_uses_service_dns: true
otlp_targets_each_application_collector: true
object_clients_use_distinct_application_buckets: true
vector_clients_do_not_traverse_traefik: true
```

<!-- field: traces -->
#### Traces

- rel19.0-uc001-persistent-application-platform S2
- srd008-persistent-application-platform AC2

<!-- case: Application and platform lifecycle preserve data until explicit deletion | status: implemented -->
### Application and platform lifecycle preserve data until explicit deletion

<!-- field: id -->
#### Id

TC3

<!-- field: use_case -->
#### Use case

rel19.0-uc001-persistent-application-platform

<!-- field: evidence_type -->
#### Evidence type

integration

<!-- field: implemented_by -->
#### Implemented by

- GH-2477
- GH-2479
- GH-2482
- GH-2483

<!-- field: command -->
#### Command

mage integration:applicationLifecycle
mage app:down chatbot-mesh
mage app:purge chatbot-mesh purge:chatbot-mesh

<!-- field: expected -->
#### Expected

```yaml
two_applications_run_concurrently: true
app_down_removes_only_selected_compute: true
app_purge_deletes_only_selected_bucket_data: true
platform_down_preserves_local_host_path: true
platform_reset_deletes_local_host_path: true
cloud_cluster_recreation_preserves_buckets: true
```

<!-- field: traces -->
#### Traces

- rel19.0-uc001-persistent-application-platform S3
- srd008-persistent-application-platform AC3

<!-- case: Local and cloud overlays preserve one application declaration contract | status: implemented -->
### Local and cloud overlays preserve one application declaration contract

<!-- field: id -->
#### Id

TC4

<!-- field: use_case -->
#### Use case

rel19.0-uc001-persistent-application-platform

<!-- field: evidence_type -->
#### Evidence type

structural_and_cloud_integration

<!-- field: implemented_by -->
#### Implemented by

- GH-2477
- GH-2484

<!-- field: command -->
#### Command

cd magefiles && go test ./gcprig -run TestFirstPartyApplicationBucketIAMIsolation
cd ../applications/chatbot-mesh && go test ./magefiles -run TestGcpDeployOverridesShape

<!-- field: expected -->
#### Expected

```yaml
application_declarations_are_identical: true
local_fake_gcs_uses_stable_filesystem_backend: true
cloud_uses_ambient_workload_identity: true
cloud_bucket_iam_denies_cross_application_access: true
local_claims_no_adversarial_tenant_authorization: true
```

<!-- field: traces -->
#### Traces

- rel19.0-uc001-persistent-application-platform S4
- srd008-persistent-application-platform AC4

<!-- case: Collector WAL fixes immutable object identity before upload | status: implemented -->
### Collector WAL fixes immutable object identity before upload

<!-- field: id -->
#### Id

TC5

<!-- field: use_case -->
#### Use case

rel19.0-uc001-persistent-application-platform

<!-- field: evidence_type -->
#### Evidence type

go_test_and_integration

<!-- field: implemented_by -->
#### Implemented by

- GH-2475
- GH-2484

<!-- field: command -->
#### Command

cd agent-core && go test ./internal/tools/otlp -run 'Test(Object|WAL)'
cd ../applications/agent-architecture && mage integration:collectorWALRestart

<!-- field: expected -->
#### Expected

```yaml
wal_precedes_first_upload_attempt: true
retries_reuse_key_and_batch_identity: true
envelope_contains_every_required_field: true
identical_duplicate_is_idempotent: true
conflicting_existing_key_is_not_overwritten: true
```

<!-- field: traces -->
#### Traces

- rel19.0-uc001-persistent-application-platform S5
- srd008-persistent-application-platform AC5

<!-- case: Query access disappears with compute and returns after collector redeployment | status: implemented -->
### Query access disappears with compute and returns after collector redeployment

<!-- field: id -->
#### Id

TC6

<!-- field: use_case -->
#### Use case

rel19.0-uc001-persistent-application-platform

<!-- field: evidence_type -->
#### Evidence type

integration

<!-- field: implemented_by -->
#### Implemented by

GH-2478

<!-- field: command -->
#### Command

mage app:down agent-architecture
mage app:up agent-architecture
curl -fsS 'http://telemetry.agent-architecture.localhost/query/traces?page_size=1'

<!-- field: expected -->
#### Expected

```yaml
running_collector_queries_only_its_application: true
retained_objects_survive_app_down: true
query_api_is_unavailable_after_app_down: true
redeployed_collector_restores_history: true
batch_identity_and_checksum_are_unchanged: true
```

<!-- field: traces -->
#### Traces

- rel19.0-uc001-persistent-application-platform S6
- srd008-persistent-application-platform AC6

<!-- case: Rig modes and the reusable application runner have one owner | status: implemented -->
### Rig modes and the reusable application runner have one owner

<!-- field: id -->
#### Id

TC7

<!-- field: use_case -->
#### Use case

rel19.0-uc001-persistent-application-platform

<!-- field: evidence_type -->
#### Evidence type

go_test_and_structural

<!-- field: implemented_by -->
#### Implemented by

GH-2479

<!-- field: command -->
#### Command

cd magefiles && go test ./apprig ./kindrig
cd .. && mage integration:applicationLifecycle

<!-- field: expected -->
#### Expected

```yaml
persistent_platform_is_canonical_application_mode: true
owned_cluster_is_limited_to_cluster_mutation: true
host_otlp_is_integration_evidence_not_query_authority: true
two_first_party_applications_use_apprig: true
downstream_verbs_are_thin: true
no_second_composition_inventory_exists: true
```

<!-- field: traces -->
#### Traces

- rel19.0-uc001-persistent-application-platform S7
- srd008-persistent-application-platform AC7

<!-- case: Application failures remain isolated and shared failures identify the platform | status: implemented -->
### Application failures remain isolated and shared failures identify the platform

<!-- field: id -->
#### Id

TC8

<!-- field: use_case -->
#### Use case

rel19.0-uc001-persistent-application-platform

<!-- field: evidence_type -->
#### Evidence type

integration

<!-- field: implemented_by -->
#### Implemented by

- GH-2477
- GH-2478
- GH-2479

<!-- field: command -->
#### Command

cd magefiles && go test ./kindrig -run 'TestPlatformHealth|TestDownPlatform'
cd ../agent-core && go test ./internal/tools/otlp -run 'TestObjectSpool'
cd ../applications/catalog && go test ./conformance -run TestApplicationPurge

<!-- field: expected -->
#### Expected

```yaml
application_failure_leaves_peer_running: true
collector_failure_leaves_peer_queryable: true
purge_failure_never_selects_peer_bucket: true
platform_service_failure_is_classified_as_shared: true
diagnosis_names_owner_and_retained_data_state: true
```

<!-- field: traces -->
#### Traces

- rel19.0-uc001-persistent-application-platform S1
- rel19.0-uc001-persistent-application-platform S3
- srd008-persistent-application-platform AC1
- srd008-persistent-application-platform AC3
