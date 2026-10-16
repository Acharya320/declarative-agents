---
id: rel10.0-uc001-role-catalog-renames
title: Evaluator and Generator Families Renamed to Critic and Executor
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A maintainer aligns two profile families with the reference role catalog: the
evaluator family is renamed to critic and the generator family to executor, the
reference role names they implement (autogenic-systems reference architecture,
roles per ITU-T Y.3502 / TM Forum IG1251). The rename covers the agents/
directories, the profile and machine names, the srd003 (evaluator) and srd002
(generator) SRD ids, the agent-role-taxonomy role ids, and the agent-core
builtin tool directories tools/builtin/critic-session and critic-point
(renamed to critic-session and critic-point, whose configs the critic profile's
tool_config_dirs point at). The dispatched tool words stay activity-named
(self_invoke, run_point, run_agent) and evaluation stays the critic's reference
activity, so no word is renamed. Historical release-record ids are retained
unchanged. Migration issues #325, #326, #330.

<!-- field: actor -->
## Actor

A maintainer aligning the role catalog with the reference architecture.

<!-- field: trigger -->
## Trigger

The reference-architecture alignment review found evaluator carries the Critic
role and generator carries the Executor role under the wrong names.

<!-- field: preconditions -->
## Preconditions

- The agent-role-taxonomy (sm-agent-role-taxonomy) lists the evaluator and generator roles.
- The evaluator and generator profile families and their integration and conformance tests exist under their current names.
- The agent-core builtin tool dirs critic-session and critic-point exist.

<!-- field: flow -->
## Flow

- F1: agents/critic is renamed to agents/critic and agents/executor to agents/executor, with the profile, machine, and internal ids renamed to match.
- F2: srd003 (evaluator) and srd002 (generator) are renamed to critic and executor, and the agent-role-taxonomy role ids follow.
- F3: The agent-core builtin tool dirs tools/builtin/critic-session and critic-point are renamed to critic-session and critic-point, and the critic profile's tool_config_dirs point at the new dirs.
- F4: The integration and conformance tests, cross-repo path references, and roadmap statuses are updated; generic child invocation uses self_invoke while run_point, run_agent, and internal/evaluation retain their roles.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | sm-agent-role-taxonomy | The taxonomy role ids for the two families become critic and executor, the reference role names. |
| T2 | test-rel10.0-role-catalog role_renames | The suite validates the renamed tree with mage validate, conformance, go test, and zero-hit greps. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

agents/ contains critic and executor families and no evaluator or generator family; the profile, machine, and SRD ids match, and grep for agents/critic and agents/executor returns no hit outside retained historical release records.

<!-- field: traces -->
#### Traces

- sm-agent-role-taxonomy

<!-- sc: S2 -->
### S2

The agent-core builtin tool dirs are critic-session and critic-point, the critic profile's tool_config_dirs point at them, and go test passes over agent-core.

<!-- field: traces -->
#### Traces

- sm-agent-role-taxonomy

<!-- sc: S3 -->
### S3

mage validate and mage conformance pass over the renamed tree, and the activity-named tool words are unchanged.

<!-- field: traces -->
#### Traces

- sm-agent-role-taxonomy

<!-- field: out_of_scope -->
## Out of scope

- Renaming the agent-core Go package internal/evaluation or the activity-named tool words; evaluation is the critic's reference activity.
- The chroma and conformance-fixture moves, which are uc002 and uc003.

<!-- field: test_suite_trace -->
## Test suite trace

test-rel10.0-role-catalog
