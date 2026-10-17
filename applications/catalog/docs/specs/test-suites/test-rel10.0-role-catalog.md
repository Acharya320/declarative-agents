---
id: test-rel10.0-role-catalog
title: Test Suite for Release 10.0 Reference-Architecture Role Alignment
release: "10.0"
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

This suite validates the completed agents/ tree alignment with the reference
role catalog while preserving profile wiring and conformance. Migration
issues #325-#328 and #330 delivered the three implemented cases below:
critic/executor renames, knowledge-family rescope, and conformance-fixture
relocation. Each case names its proof commands and expected outputs.

<!-- field: traces -->
## Traces

- [rel10.0-uc001-role-catalog-renames](../use-cases/rel10.0-uc001-role-catalog-renames.md)
- [rel10.0-uc002-knowledge-family-rescope](../use-cases/rel10.0-uc002-knowledge-family-rescope.md)
- [rel10.0-uc003-conformance-fixture-relocation](../use-cases/rel10.0-uc003-conformance-fixture-relocation.md)

<!-- field: preconditions -->
## Preconditions

- agent-core is available as a sibling ../../agent-core checkout for builtin tool-dir tests; mage validate may instead use demo.yaml core_root.
- Docker and Ollama are available for the Integration.Chroma and Integration.MonitorControl cases; they skip cleanly otherwise.

<!-- field: test_cases -->
## Test cases

<!-- case: Critic and executor renames validate over the tree | status: implemented -->
### Critic and executor renames validate over the tree

After the rename, agents/critic and agents/executor exist and no evaluator or
generator family remains; the agent-core builtin tool dirs are critic-session
and critic-point and the critic profile's tool_config_dirs point at them. The
case validates mage validate over the renamed profile tree, go test over
agent-core, and zero-hit greps.

<!-- field: use_case -->
#### Use case

[rel10.0-uc001-role-catalog-renames](../use-cases/rel10.0-uc001-role-catalog-renames.md)

<!-- field: integration_target -->
#### Integration target

mage validate

<!-- field: go_test -->
#### Go test

cd ../../agent-core && go test ./...

<!-- field: inputs | form: "tree" -->
#### Inputs

- renamed_families:
  1. evaluator->critic
  1. generator->executor
- renamed_tool_dirs:
  1. critic-session->critic-session
  1. critic-point->critic-point

<!-- field: expected | form: "tree" -->
#### Expected

- agents_lists_critic_and_executor: true
- no_evaluator_or_generator_family: true
- grep_agents_evaluator_generator_zero_hits: true
- agent_core_builtin_dirs_renamed: true
- mage_validate_passes: true
- activity_tool_words_unchanged: true

<!-- field: traces -->
#### Traces

- [rel10.0-uc001-role-catalog-renames](../use-cases/rel10.0-uc001-role-catalog-renames.md) [S1](../use-cases/rel10.0-uc001-role-catalog-renames.md#s1)
- [rel10.0-uc001-role-catalog-renames](../use-cases/rel10.0-uc001-role-catalog-renames.md) [S2](../use-cases/rel10.0-uc001-role-catalog-renames.md#s2)
- [rel10.0-uc001-role-catalog-renames](../use-cases/rel10.0-uc001-role-catalog-renames.md) [S3](../use-cases/rel10.0-uc001-role-catalog-renames.md#s3)

<!-- case: Corpus ingest and reader relocate under knowledge-manager (rescope) | status: implemented -->
### Corpus ingest and reader relocate under knowledge-manager (rescope)

After the move, agents/knowledge-manager hosts corpus-ingest and
corpus-reader and agents/chroma is removed. The case validates the rescope:
the profiles resolve from the new paths (mage validate), the
documentation-curator conformance test passes from the knowledge-manager
family, and a zero-hit grep for agents/chroma. It does not validate
The canonical ListingCorpus/trusted-path contract is covered by
TestCorpusIngestListsTrustedCorpusBeforeModelControl, and the mesh
integration:chroma target consumes that library program by reference.
The standalone corpus-reader grounded-answer tracer remains separate.

<!-- field: use_case -->
#### Use case

[rel10.0-uc002-knowledge-family-rescope](../use-cases/rel10.0-uc002-knowledge-family-rescope.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run 'Test(KnowledgeManager|CorpusIngest)'

<!-- field: inputs | form: "tree" -->
#### Inputs

- moved_profiles:
  1. agents/knowledge-manager/corpus-ingest->agents/knowledge-manager/corpus-ingest
  1. agents/knowledge-manager/corpus-reader->agents/knowledge-manager/corpus-reader

<!-- field: expected | form: "tree" -->
#### Expected

- knowledge_manager_hosts_corpus_profiles: true
- agents_chroma_removed: true
- grep_agents_chroma_zero_hits: true
- mage_validate_passes: true
- documentation_curator_conformance_passes: true

<!-- field: traces -->
#### Traces

- [rel10.0-uc002-knowledge-family-rescope](../use-cases/rel10.0-uc002-knowledge-family-rescope.md) [S1](../use-cases/rel10.0-uc002-knowledge-family-rescope.md#s1)
- [rel10.0-uc002-knowledge-family-rescope](../use-cases/rel10.0-uc002-knowledge-family-rescope.md) [S2](../use-cases/rel10.0-uc002-knowledge-family-rescope.md#s2)

<!-- case: Conformance fixtures relocate and validate from testdata/conformance | status: implemented -->
### Conformance fixtures relocate and validate from testdata/conformance

After the move, testdata/conformance holds the rest, control, and lifecycle
fixtures and agents/ lists only role families. The formal Go evidence runs
the deterministic rest, control, and lifecycle fixture cases from their new
paths; the separately documented REST Ollama variant remains behind the
explicit live-conformance opt-in. The case also validates mage conformance
and Integration.MonitorControl, that mage validate discovery reaches the
relocated fixtures, and an agents/ directory listing plus a zero-hit grep.

<!-- field: use_case -->
#### Use case

[rel10.0-uc003-conformance-fixture-relocation](../use-cases/rel10.0-uc003-conformance-fixture-relocation.md)

<!-- field: integration_target -->
#### Integration target

mage conformance

<!-- field: go_test -->
#### Go test

go test ./conformance -run 'Test(RestConformance|RestShippedProfileWiring|ControlConformance|LifecycleApprovalSuspendResume|LifecycleHistoryRollback)'

<!-- field: inputs | form: "tree" -->
#### Inputs

- moved_fixtures:
  1. agents/rest->testdata/conformance/rest
  1. agents/control->testdata/conformance/control
  1. agents/lifecycle->testdata/conformance/lifecycle

<!-- field: expected | form: "tree" -->
#### Expected

- fixtures_under_testdata_conformance: true
- agents_lists_only_role_families: true
- grep_agents_rest_control_lifecycle_zero_hits: true
- mage_conformance_and_monitor_control_pass: true

<!-- field: traces -->
#### Traces

- [rel10.0-uc003-conformance-fixture-relocation](../use-cases/rel10.0-uc003-conformance-fixture-relocation.md) [S1](../use-cases/rel10.0-uc003-conformance-fixture-relocation.md#s1)
- [rel10.0-uc003-conformance-fixture-relocation](../use-cases/rel10.0-uc003-conformance-fixture-relocation.md) [S2](../use-cases/rel10.0-uc003-conformance-fixture-relocation.md#s2)
