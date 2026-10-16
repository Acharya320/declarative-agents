---
id: test-rel24.0-objectstore-tools
title: Test Suite for Release 24.0 Objectstore Tools
release: "24.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel24.0-uc001-objectstore-tools

<!-- field: preconditions -->
## Preconditions

- srd059-objectstore-tools defines the words, the driver selection, and the receipt contract.
- Unit cases run on mem:// and file:// buckets with no container and no network.
- The integration case needs docker on the host and skips naming it otherwise.

<!-- field: test_cases -->
## Test cases

<!-- case: Read, write, and list round-trip on memory and file buckets | status: implemented -->
### Read, write, and list round-trip on memory and file buckets

The three commands against mem:// and file:// connections declared through DecodeToolConfig.

<!-- field: use_case -->
#### Use case

rel24.0-uc001-objectstore-tools

<!-- field: go_test -->
#### Go test

go test ./internal/tools/objectstore/... -run 'TestObjectRoundTrip'

<!-- field: expected -->
#### Expected

```yaml
write_then_read_returns_the_text: true
list_by_prefix_names_the_keys: true
both_schemes_behave_identically: true
```

<!-- field: traces -->
#### Traces

- rel24.0-uc001-objectstore-tools S2
- srd059-objectstore-tools AC2
- srd059-objectstore-tools R2.1

<!-- case: Undo restores the prior object and deletes a created one | status: implemented -->
### Undo restores the prior object and deletes a created one

A write over existing bytes and a write to an absent key, each undone through the receipt alone.

<!-- field: use_case -->
#### Use case

rel24.0-uc001-objectstore-tools

<!-- field: go_test -->
#### Go test

go test ./internal/tools/objectstore/... -run 'TestObjectWriteUndo'

<!-- field: expected -->
#### Expected

```yaml
prior_bytes_restored: true
created_key_deleted: true
receipt_survives_redecode: true
```

<!-- field: traces -->
#### Traces

- rel24.0-uc001-objectstore-tools S3
- srd059-objectstore-tools AC3
- srd059-objectstore-tools R3.1

<!-- case: Driver selection follows the declared URL and fails unknown schemes by name | status: implemented -->
### Driver selection follows the declared URL and fails unknown schemes by name

mem, file, gs, and s3 URLs open their drivers; ftp:// and an absent bucket_url fail at configuration.

<!-- field: use_case -->
#### Use case

rel24.0-uc001-objectstore-tools

<!-- field: go_test -->
#### Go test

go test ./internal/tools/objectstore/... -run 'TestBucketSelection'

<!-- field: expected -->
#### Expected

```yaml
scheme_selects_driver: true
unknown_scheme_named: true
missing_bucket_url_named: true
```

<!-- field: traces -->
#### Traces

- rel24.0-uc001-objectstore-tools S2
- srd059-objectstore-tools AC2
- srd059-objectstore-tools R2.1

<!-- case: The gs driver takes an explicit endpoint and reads no environment | status: implemented -->
### The gs driver takes an explicit endpoint and reads no environment

The endpoint configuration field points the gs driver at a local HTTP server; the process environment carries a decoy STORAGE_EMULATOR_HOST that must be ignored.

<!-- field: use_case -->
#### Use case

rel24.0-uc001-objectstore-tools

<!-- field: go_test -->
#### Go test

go test ./internal/tools/objectstore/... -run 'TestExplicitEndpoint'

<!-- field: expected -->
#### Expected

```yaml
endpoint_reaches_the_configured_server: true
decoy_environment_ignored: true
```

<!-- field: traces -->
#### Traces

- rel24.0-uc001-objectstore-tools S2
- srd059-objectstore-tools AC2
- srd059-objectstore-tools R2.2

<!-- case: A fixture profile boots the four words and the default surface is unchanged | status: implemented -->
### A fixture profile boots the four words and the default surface is unchanged

The bootstrap catalog covers the objectstore inits, and tools/builtin/all.yaml does not import the group.

<!-- field: use_case -->
#### Use case

rel24.0-uc001-objectstore-tools

<!-- field: go_test -->
#### Go test

go test ./cmd/agent/... -run 'TestBuiltinFactoryCatalogCoversSelectedActiveInits'

<!-- field: expected -->
#### Expected

```yaml
inits_covered: true
default_surface_unchanged: true
```

<!-- field: traces -->
#### Traces

- rel24.0-uc001-objectstore-tools S1
- srd059-objectstore-tools AC1
- srd059-objectstore-tools R1.1

<!-- case: The declarations pass every hand-maintained gate | status: implemented -->
### The declarations pass every hand-maintained gate

Contract completeness stays empty, the reach test sees the words, output contract and side-effect kinds agree, and the rollback policy accepts the undo strategy.

<!-- field: use_case -->
#### Use case

rel24.0-uc001-objectstore-tools

<!-- field: go_test -->
#### Go test

mage test

<!-- field: expected -->
#### Expected

```yaml
completeness_baseline_empty: true
reach_test_sees_the_words: true
rollback_policy_accepts_strategy: true
```

<!-- field: traces -->
#### Traces

- rel24.0-uc001-objectstore-tools S1
- srd059-objectstore-tools AC1
- srd059-objectstore-tools R3.3

<!-- case: The family runs end to end against fake-gcs-server | status: implemented -->
### The family runs end to end against fake-gcs-server

docker run starts a pinned fake-gcs-server; write, read, list, and both undo directions run over gs:// with the explicit endpoint; the test skips naming docker where it is absent and cleans up on every path.

<!-- field: use_case -->
#### Use case

rel24.0-uc001-objectstore-tools

<!-- field: go_test -->
#### Go test

go test ./internal/tools/objectstore/... -run 'TestObjectstoreAgainstFakeGCS'

<!-- field: expected -->
#### Expected

```yaml
end_to_end_roundtrip: true
undo_both_directions: true
skip_names_docker: true
image_pinned_by_digest: true
```

<!-- field: traces -->
#### Traces

- rel24.0-uc001-objectstore-tools S4
- srd059-objectstore-tools AC4
- srd059-objectstore-tools R4.2

<!-- case: Declaration-bound delete-prefix is bounded, isolated, verified, and audited | status: implemented -->
### Declaration-bound delete-prefix is bounded, isolated, verified, and audited

The request supplies a malicious peer prefix but the builder deletes only its configured application prefix; a second case exceeds max_objects and refuses without mutation.

<!-- field: use_case -->
#### Use case

rel24.0-uc001-objectstore-tools

<!-- field: go_test -->
#### Go test

go test ./internal/tools/objectstore/... -run 'TestDeletePrefix'

<!-- field: expected -->
#### Expected

```yaml
request_storage_authority_ignored: true
peer_prefix_unchanged: true
over_budget_refused_without_delete: true
post_delete_remaining_zero: true
external_audit_written: true
```

<!-- field: traces -->
#### Traces

- rel24.0-uc001-objectstore-tools S5
- srd059-objectstore-tools AC5
- srd059-objectstore-tools R5.1
- srd059-objectstore-tools R5.2
- srd059-objectstore-tools R5.3
