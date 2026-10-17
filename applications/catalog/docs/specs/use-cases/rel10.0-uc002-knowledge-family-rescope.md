---
id: rel10.0-uc002-knowledge-family-rescope
title: Chroma Profiles Rescoped Under Knowledge Management as Corpus Ingest and Reader
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A maintainer moves the chroma profile family, named after a vendor database, to
its reference role under Knowledge Management: agents/knowledge-manager/corpus-ingest becomes
agents/knowledge-manager/corpus-ingest and agents/knowledge-manager/corpus-reader becomes
agents/knowledge-manager/corpus-reader (the rag-server stays a RAG boundary but
moves with the family per the migration issue). The move covers the profile
directories, the profile/machine ids, the shared chroma REST assets, the
integration and conformance tests, and the roadmap statuses. Integration.Chroma
passes from the new paths and agents/chroma is removed. Migration issue #327.

<!-- field: actor -->
## Actor

A maintainer aligning the role catalog with the reference architecture.

<!-- field: trigger -->
## Trigger

The alignment review found the chroma family named after a vendor database
instead of its Knowledge Management role.

<!-- field: preconditions -->
## Preconditions

- The chroma ingest and reader profiles exist under agents/chroma with shared REST assets.
- Integration.Chroma and the knowledge-manager conformance test reference the current paths.
- The knowledge-manager family already hosts the documentation-curator profile.

<!-- field: flow -->
## Flow

- F1: agents/knowledge-manager/corpus-ingest and agents/knowledge-manager/corpus-reader move to agents/knowledge-manager/corpus-ingest and corpus-reader, with the profile and machine ids renamed to match.
- F2: The shared chroma REST assets move with the family and the profiles' relative references are updated; agents/chroma is removed.
- F3: Integration.Chroma and the knowledge-manager conformance test are updated to the new paths and pass; the roadmap statuses are updated.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [sm-agent-role-taxonomy](../semantic-models/agent-role-taxonomy.md) | The knowledge-manager family gains corpus-ingest and corpus-reader, replacing the vendor-named chroma family. |
| T2 | test-rel10.0-role-catalog knowledge_rescope | The suite validates the moved tree with Integration.Chroma and a zero-hit grep for agents/chroma. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The corpus-ingest and corpus-reader profiles live under agents/knowledge-manager, agents/chroma is removed, and grep for agents/chroma returns no hit outside retained historical release records.

<!-- field: traces -->
#### Traces

- [sm-agent-role-taxonomy](../semantic-models/agent-role-taxonomy.md)

<!-- sc: S2 -->
### S2

Integration.Chroma passes from the new paths and mage validate passes over the moved tree.

<!-- field: traces -->
#### Traces

- [sm-agent-role-taxonomy](../semantic-models/agent-role-taxonomy.md)

<!-- field: out_of_scope -->
## Out of scope

- The critic and executor renames (uc001) and the conformance-fixture relocation (uc003).
- Changing the Chroma REST word behavior; only the profile family location changes.

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel10.0-role-catalog](../test-suites/test-rel10.0-role-catalog.md)
