---
id: test-rel19.0-signed-boundary-contracts
title: Test Suite for Release 19.0 Signed Boundary Contracts
release: "19.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- [rel19.0-uc001-signed-boundary-contracts](../use-cases/rel19.0-uc001-signed-boundary-contracts.md)

<!-- field: preconditions -->
## Preconditions

- srd051-declaration-types R6 defines signatures and contract defaults, with R6.14 and R6.15 for boundary, stateful_internal, and uncategorized tools.
- The corpus audit in pkg/spec and the declaration legacy gate in internal/declstyle are available.

<!-- field: test_cases -->
## Test cases

<!-- case: An emits-only signature completes a boundary tool with its effect blocks | status: implemented -->
### An emits-only signature completes a boundary tool with its effect blocks

Check a boundary tool with a signature naming only its signals, an output schema, and side_effects, reversibility, and undo; then drop each effect block in turn.

<!-- field: use_case -->
#### Use case

[rel19.0-uc001-signed-boundary-contracts](../use-cases/rel19.0-uc001-signed-boundary-contracts.md)

<!-- field: go_test -->
#### Go test

go test ./pkg/spec -run 'TestEmitsOnlySignatureCompletesABoundaryToolWithItsEffectBlocks|TestValidate_BoundaryWordMustDeclareASignature|TestSignatureNeverWaivesBoundarySideEffects|TestSignatureLeavesStatefulInternalSideEffectsRequired'

<!-- field: expected | form: "tree" -->
#### Expected

- complete_with_effect_blocks: true
- each_missing_effect_block_named: true

<!-- field: traces -->
#### Traces

- [rel19.0-uc001-signed-boundary-contracts](../use-cases/rel19.0-uc001-signed-boundary-contracts.md) [S1](../use-cases/rel19.0-uc001-signed-boundary-contracts.md#s1)
- [srd051-declaration-types](../software-requirements/srd051-declaration-types.md) [AC8](../software-requirements/srd051-declaration-types.md#ac8)
- [srd051-declaration-types](../software-requirements/srd051-declaration-types.md) [R6.14](../software-requirements/srd051-declaration-types.md#r614)

<!-- case: The legacy gate lists unsigned boundary and stateful tools and uncategorized tools | status: implemented -->
### The legacy gate lists unsigned boundary and stateful tools and uncategorized tools

Classify a declaration holding an unsigned boundary tool, an unsigned stateful_internal tool, a signed boundary tool, and a tool with no category; the repository matches its seeded baseline.

<!-- field: use_case -->
#### Use case

[rel19.0-uc001-signed-boundary-contracts](../use-cases/rel19.0-uc001-signed-boundary-contracts.md)

<!-- field: go_test -->
#### Go test

go test ./internal/declstyle -run 'TestUntypedAndUncategorizedToolClassification|TestDeclarationLegacyBaseline'

<!-- field: expected | form: "tree" -->
#### Expected

- unsigned_boundary_and_stateful_untyped: true
- no_category_uncategorized: true
- signed_boundary_listed: false

<!-- field: traces -->
#### Traces

- [rel19.0-uc001-signed-boundary-contracts](../use-cases/rel19.0-uc001-signed-boundary-contracts.md) [S2](../use-cases/rel19.0-uc001-signed-boundary-contracts.md#s2)
- [srd051-declaration-types](../software-requirements/srd051-declaration-types.md) [AC8](../software-requirements/srd051-declaration-types.md#ac8)
- [srd051-declaration-types](../software-requirements/srd051-declaration-types.md) [R6.14](../software-requirements/srd051-declaration-types.md#r614)
- [srd051-declaration-types](../software-requirements/srd051-declaration-types.md) [R6.15](../software-requirements/srd051-declaration-types.md#r615)
