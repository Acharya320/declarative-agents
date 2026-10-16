---
id: rel10.0-uc003-conformance-fixture-relocation
title: Conformance Grammar Fixtures Relocated to testdata/conformance
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A maintainer moves the three profile directories that are conformance grammar
fixtures rather than roles -- rest, control, and lifecycle -- out of agents/ to
testdata/conformance, beside testdata/integration. After the move agents/ lists
only role-shaped families. The move covers the fixture directories, the
conformance test references, the validate discovery (extended so profile-wiring
validation still reaches the relocated fixtures), and Integration.MonitorControl,
which reads the control fixture. mage conformance and Integration.MonitorControl
pass from the new paths. Migration issue #328.

<!-- field: actor -->
## Actor

A maintainer aligning the role catalog with the reference architecture.

<!-- field: trigger -->
## Trigger

The alignment review found rest, control, and lifecycle are conformance grammar
fixtures, not reference roles, sitting in the role directory.

<!-- field: preconditions -->
## Preconditions

- The rest, control, and lifecycle profile directories exist under agents/ and are referenced by the conformance tests.
- Integration.MonitorControl reads the control fixture from its current path.
- testdata/integration already holds integration fixtures beside agents/.

<!-- field: flow -->
## Flow

- F1: agents/rest, agents/control, and agents/lifecycle move to testdata/conformance/rest, control, and lifecycle.
- F2: The conformance tests and Integration.MonitorControl are updated to the new paths, and the mage validate discovery is extended so profile-wiring validation reaches the relocated fixtures.
- F3: agents/ lists only role-shaped families (critic, executor, planner, specification-critic, runtime-state-reader, knowledge-manager, bench, scenario-critic, and mock); the roadmap statuses are updated.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | sm-agent-role-taxonomy | The taxonomy's role families no longer include the conformance fixtures rest, control, and lifecycle. |
| T2 | test-rel10.0-role-catalog fixture_relocation | The suite validates the relocation with mage conformance, Integration.MonitorControl, and a zero-hit grep and an agents/ directory listing. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The rest, control, and lifecycle fixtures live under testdata/conformance, agents/ lists only role families, and grep for agents/rest, agents/control, and agents/lifecycle returns no hit outside retained historical release records.

<!-- field: traces -->
#### Traces

- sm-agent-role-taxonomy

<!-- sc: S2 -->
### S2

mage conformance and Integration.MonitorControl pass from the new paths, and mage validate discovery reaches the relocated fixtures.

<!-- field: traces -->
#### Traces

- sm-agent-role-taxonomy

<!-- field: out_of_scope -->
## Out of scope

- The critic and executor renames (uc001) and the chroma rescope (uc002).
- Changing the fixture grammar content; only the location changes.

<!-- field: test_suite_trace -->
## Test suite trace

test-rel10.0-role-catalog
