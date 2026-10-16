---
id: test-rel17.0-declaration-imports
title: Test Suite for Release 17.0 Declaration Imports
release: "17.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel17.0-uc001-declaration-imports

<!-- field: preconditions -->
## Preconditions

- srd050-declaration-imports defines unit syntax, graph order, merge rules, usedness, and legacy compatibility.
- The single closure loader, strict YAML decoder, and canonical dump are available.

<!-- field: test_cases -->
## Test cases

<!-- case: Tool and REST import chains resolve once in deterministic order | status: planned -->
### Tool and REST import chains resolve once in deterministic order

Resolve two-level and diamond imports from their declaring files and compare two closure inventories and dumps byte for byte.

<!-- field: use_case -->
#### Use case

rel17.0-uc001-declaration-imports

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/tools/rest/definition -run TestDeclarationImportChain

<!-- field: expected -->
#### Expected

```yaml
canonical_sources_read_once: true
traversal: first-discovery depth-first postorder
```

<!-- field: traces -->
#### Traces

- rel17.0-uc001-declaration-imports S1
- srd050-declaration-imports AC1

<!-- case: Tool unit collisions cycles and invalid overrides fail with provenance | status: planned -->
### Tool unit collisions cycles and invalid overrides fail with provenance

Exercise duplicate unit names, duplicate imported and local tools, cycles, unmarked collisions, missing override targets, and imported override markers.

<!-- field: use_case -->
#### Use case

rel17.0-uc001-declaration-imports

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog -run TestToolDeclarationImportFailures

<!-- field: expected -->
#### Expected

```yaml
every_error_names_unit_and_path: true
imported_collision_precedes_override: true
```

<!-- field: traces -->
#### Traces

- rel17.0-uc001-declaration-imports S2
- srd050-declaration-imports AC2

<!-- case: One local tool override replaces one imported target in place | status: planned -->
### One local tool override replaces one imported target in place

Apply one marked local override and verify value, provenance, and deterministic position.

<!-- field: use_case -->
#### Use case

rel17.0-uc001-declaration-imports

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog -run TestToolDeclarationLocalOverride

<!-- field: expected -->
#### Expected

```yaml
target_count: 1
provenance: local unit
```

<!-- field: traces -->
#### Traces

- rel17.0-uc001-declaration-imports S2
- srd050-declaration-imports AC3

<!-- case: REST imports merge every family and compile source-relative OpenAPI | status: planned -->
### REST imports merge every family and compile source-relative OpenAPI

Merge every REST family, retain entry provenance, compile imported OpenAPI into a local binding, and reject duplicate names and conflicting versions.

<!-- field: use_case -->
#### Use case

rel17.0-uc001-declaration-imports

<!-- field: go_test -->
#### Go test

go test ./internal/tools/rest/definition -run TestRESTDeclarationImports

<!-- field: expected -->
#### Expected

```yaml
all_families_merged: true
openapi_base: declaring source
```

<!-- field: traces -->
#### Traces

- rel17.0-uc001-declaration-imports S3
- srd050-declaration-imports AC4

<!-- case: Closure usedness accepts contributing edges and rejects unused edges | status: planned -->
### Closure usedness accepts contributing edges and rejects unused edges

Traverse selected-tool and fully qualified REST nodes, including override and OpenAPI attribution and both arms of a used diamond.

<!-- field: use_case -->
#### Use case

rel17.0-uc001-declaration-imports

<!-- field: go_test -->
#### Go test

go test ./internal/load -run TestDeclarationImportUsedness

<!-- field: expected -->
#### Expected

```yaml
used_edges_accepted: true
unused_error_order: importer then imported unit then path
```

<!-- field: traces -->
#### Traces

- rel17.0-uc001-declaration-imports S4
- srd050-declaration-imports AC5

<!-- case: Legacy includes warn once and mixed closures fail | status: planned -->
### Legacy includes warn once and mixed closures fail

Load a legacy-only include diamond and a closure containing both import and include edges.

<!-- field: use_case -->
#### Use case

rel17.0-uc001-declaration-imports

<!-- field: go_test -->
#### Go test

go test ./internal/tools/catalog ./internal/load -run TestLegacyIncludesCompatibility

<!-- field: expected -->
#### Expected

```yaml
warnings_per_canonical_includer: 1
mixed_closure_rejected: true
```

<!-- field: traces -->
#### Traces

- rel17.0-uc001-declaration-imports S5
- srd050-declaration-imports AC6

<!-- case: Every entry path shares one import failure and control proof | status: planned -->
### Every entry path shares one import failure and control proof

Compare validation, dump, startup, and profile-audit diagnostics, then compare control-imports with control after removing identity and files.

<!-- field: use_case -->
#### Use case

rel17.0-uc001-declaration-imports

<!-- field: go_test -->
#### Go test

go test ./cmd/agent ./internal/load -run 'TestDeclarationImportEntryPaths|TestControlImportsDump'

<!-- field: expected -->
#### Expected

```yaml
diagnostic_implementations: 1
normalized_control_dump_equal: true
```

<!-- field: traces -->
#### Traces

- rel17.0-uc001-declaration-imports S6
- srd050-declaration-imports AC7
- srd050-declaration-imports AC8
