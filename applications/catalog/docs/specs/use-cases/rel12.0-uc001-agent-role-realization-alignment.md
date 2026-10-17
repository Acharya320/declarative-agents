---
id: rel12.0-uc001-agent-role-realization-alignment
title: Migrate Concrete Agent Names as One Versioned Compatibility Release
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: release -->
## Release

12.0

<!-- field: summary -->
## Summary

A maintainer completes the repository-wide role-realization naming migration
after the semantic authority, binding validator, and individual renames land.
Reusable catalog paths retain non-forking v0 wrappers; application-local
identities migrate atomically. Consumers pin the exact planned catalog
release, packages record compatibility separately from source provenance, and
release tests prove the canonical and legacy tag pair without publishing tags
from a feature branch.

<!-- field: actor -->
## Actor

A maintainer preparing the naming-alignment release.

<!-- field: trigger -->
## Trigger

A maintainer starts the coordinated compatibility release.

<!-- field: preconditions -->
## Preconditions

- Issues 1068 through 1070 and 1072 through 1077 are complete.
- The shared role-realization semantic model and executable validator are present.
- Catalog releases 07.1 and 11.0 retain their independent in-progress evidence gaps.

<!-- field: flow -->
## Flow

- F1: The maintainer verifies every old-to-new profile or identity in the migration table.
- F2: Catalog aliases select canonical files through one-file wrappers supported through v0 and removed at v1.
- F3: Coding-agent and chatbot-mesh package metadata pin applications/catalog/v0.20260727.0 and retain checkout-versus-release provenance.
- F4: Catalog, coding-agent, and chatbot-mesh package closures are regenerated and validated from their source inventories.
- F5: Root release tests derive one atomic v0.20260727.0, applications/catalog/v0.20260727.0, and agent-profiles/v0.20260727.0 tag set.
- F6: After merge, a maintainer runs mage tag from a clean main checkout; this use case does not publish immutable tags.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [migration-v0.20260727.0-agent-role-realization-alignment](../../migrations/v0.20260727.0-agent-role-realization-alignment.md) | Authoritative old-to-new table, compatibility duration, consumer pins, and publication boundary. |
| T2 | sm-agent-role-realizations | Binding, classification, inheritance, and alias authority. |
| T3 | [test-rel12.0-agent-role-realization-alignment](../test-suites/test-rel12.0-agent-role-realization-alignment.md) | Executable release, closure, render, drift, and zero-hit evidence. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Every renamed profile or runtime identity has an old-to-new migration entry and explicit compatibility policy.

<!-- sc: S2 -->
### S2

Catalog compatibility wrappers contain only profile.yaml and reference canonical implementation files without copies.

<!-- sc: S3 -->
### S3

Coding-agent closure and chatbot-mesh staged Helm inventory contain the canonical paths and exact compatibility pin.

<!-- sc: S4 -->
### S4

Scoped checks find no stale active path or runtime identity outside wrappers, migration history, and explicit historical aliases.

<!-- sc: S5 -->
### S5

Root and catalog tests prove exact atomic tag semantics while historical tags remain untouched.

<!-- sc: S6 -->
### S6

Root audit, test, and stats plus catalog validation/conformance/integration and application package/Helm checks pass.

<!-- field: out_of_scope -->
## Out of scope

- Closing the independent Release 07.1 or 11.0 behavioral evidence gaps.
- Publishing immutable release tags before this branch is merged to main.
- Renaming stable capability words, historical IDs, or runtime monitor infrastructure.

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel12.0-agent-role-realization-alignment](../test-suites/test-rel12.0-agent-role-realization-alignment.md)
