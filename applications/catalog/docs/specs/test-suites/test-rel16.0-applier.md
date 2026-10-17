---
id: test-rel16.0-applier
title: Canonical Applier Family Conformance
status: implemented
release: "16.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel16.0-uc001-applier-family

<!-- field: preconditions -->
## Preconditions

- The catalog applier profile closes against agent-core shared tools.
- Each application manifest declares the canonical applier and one local wrapper.
- Each wrapper owns its REST and rollout declarations.

<!-- field: test_cases -->
## Test cases

<!-- case: Applier catalog closure | status: implemented -->
### Applier catalog closure

<!-- field: use_case -->
#### Use case

rel16.0-uc001-applier-family

<!-- field: id -->
#### ID

applier-catalog-closure

<!-- field: verifies -->
#### Verifies

- S1

<!-- field: traces -->
#### Traces

- rel16.0-uc001-applier-family S1

<!-- field: go_test -->
#### Go test

mage validate

<!-- field: assertions -->
#### Assertions

- The catalog applier profile and request variants pass agent-core strict configuration validation.
- The apply and rollout selections are name-only.

<!-- case: Applier machine contract | status: implemented -->
### Applier machine contract

<!-- field: use_case -->
#### Use case

rel16.0-uc001-applier-family

<!-- field: id -->
#### ID

applier-machine-contract

<!-- field: verifies -->
#### Verifies

- S1
- S3
- S4
- S5

<!-- field: traces -->
#### Traces

- rel16.0-uc001-applier-family S1
- rel16.0-uc001-applier-family S3
- rel16.0-uc001-applier-family S4
- rel16.0-uc001-applier-family S5

<!-- field: go_test -->
#### Go test

go test ./conformance -run '^TestCanonicalApplier'

<!-- field: assertions -->
#### Assertions

- Validate precedes apply and verify.
- Verification failure reaches rollback.
- Terminal outcomes remain distinct.

<!-- case: Applier wrapper ownership | status: implemented -->
### Applier wrapper ownership

<!-- field: use_case -->
#### Use case

rel16.0-uc001-applier-family

<!-- field: id -->
#### ID

applier-wrapper-ownership

<!-- field: verifies -->
#### Verifies

- S2
- S6

<!-- field: traces -->
#### Traces

- rel16.0-uc001-applier-family S2
- rel16.0-uc001-applier-family S6

<!-- field: go_test -->
#### Go test

go test ./conformance -run '^TestApplierPromotion'

<!-- field: assertions -->
#### Assertions

- Every wrapper references canonical machines and shared ToolDefs.
- REST and exec declarations remain local.
- No promoted asset has an application-local byte copy.

<!-- case: Applier consumers | status: implemented -->
### Applier consumers

<!-- field: use_case -->
#### Use case

rel16.0-uc001-applier-family

<!-- field: id -->
#### ID

applier-consumers

<!-- field: verifies -->
#### Verifies

- S2
- S3
- S4
- S5
- S6

<!-- field: traces -->
#### Traces

- rel16.0-uc001-applier-family S2
- rel16.0-uc001-applier-family S3
- rel16.0-uc001-applier-family S4
- rel16.0-uc001-applier-family S5
- rel16.0-uc001-applier-family S6

<!-- field: go_test -->
#### Go test

mage audit

<!-- field: assertions -->
#### Assertions

- Fake, render, authority, package, rollout, and rollback tests pass in all consumers.
- Applicable live kind tests remain dependency gated.

<!-- case: Applier statistics | status: implemented -->
### Applier statistics

<!-- field: use_case -->
#### Use case

rel16.0-uc001-applier-family

<!-- field: id -->
#### ID

applier-statistics

<!-- field: verifies -->
#### Verifies

- S6

<!-- field: traces -->
#### Traces

- rel16.0-uc001-applier-family S6

<!-- field: go_test -->
#### Go test

mage stats

<!-- field: assertions -->
#### Assertions

- Catalog counts applier once.
- Application appliers are composition wrappers and contribute no additional canonical agent.
