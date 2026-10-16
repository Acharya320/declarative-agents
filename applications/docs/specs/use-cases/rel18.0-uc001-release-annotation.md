---
id: rel18.0-uc001-release-annotation
title: A Degenerate Release Annotation Is Refused Before the Push
status: planned
test_suite: test-rel18.0-release-annotation
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: release -->
## Release

18.0

<!-- field: summary -->
## Summary

An operator finishes a release: the gates ran, the tag exists on main, and
the session annotated it with a summary of what shipped. Before pushing,
the operator runs one check against the tag. A real summary passes and the
push proceeds. An annotation that says nothing — the tag's own name, an
empty subject, a one-word message — is refused with the rule it broke, and
the operator re-annotates before anything leaves the machine. The tag's
annotation is the only record the release line has, so what ships is a
record that says something.

<!-- field: actor -->
## Actor

An operator pushing a release, following the runbook or the release skill.

<!-- field: trigger -->
## Trigger

A release tag exists and is annotated, and the operator is about to push
it.

<!-- field: preconditions -->
## Preconditions

- ENG01 carries the release-annotation section and the tag:verify verb; srd007 owns the rules.
- mage tag created the tag on a gated main commit; the session annotated it.
- The runbook and the release skill place the check between annotation and push.

<!-- field: flow -->
## Flow

- F1: The operator runs the verify target with the tag's name.
  - traces:
    - srd007-release-annotation R4.2
- F2: The check resolves the tag. A missing tag, or a lightweight ref where
  an annotated tag should be, is refused by name.
  - traces:
    - srd007-release-annotation R1.1
- F3: The annotation's subject and body are read. A subject equal to the
  tag name, an empty subject, or a message under the substance bar is
  refused, quoting the rule.
  - traces:
    - srd007-release-annotation R2.1
    - srd007-release-annotation R2.2
- F4: Reachability from main is checked, so a tag on a side branch cannot
  ship as a release.
  - traces:
    - srd007-release-annotation R3.1
- F5: On refusal the operator re-annotates the tag with a real summary and
  runs the check again. Nothing was created or deleted by the failed
  run, so the retry starts clean.
  - traces:
    - srd007-release-annotation R4.1
- F6: On success the operator pushes, and the annotation that ships is a
  record of what the release contained.
  - traces:
    - srd007-release-annotation R4.2

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | docs/engineering/eng01-kind-test-demo-rig.md | States that the annotation is the release record and the verb's contract. |
| T2 | magefiles/tag_verify.go | The check applying the four rules with one named reason per failure. |
| T3 | magefiles/tag.go | Creates the lightweight ref the session annotates; unchanged by this release. |
| T4 | srd007-release-annotation AC1 | The degenerate-annotation refusals. |
| T5 | srd007-release-annotation AC2 | Reachability from main. |
| T6 | srd007-release-annotation AC3 | A run mutates nothing. |
| T7 | srd007-release-annotation AC4 | Placement in the runbook and the skill. |
| T8 | test-rel18.0-release-annotation | Structural and executable evidence for this flow. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Every degenerate shape — missing tag, lightweight ref, subject equal
to the tag name, empty subject, under-bar summary — is refused with
its own named reason.

<!-- field: traces -->
#### Traces

- srd007-release-annotation AC1

<!-- sc: S2 -->
### S2

A tag off main is refused, and an absent main fails by name rather
than passing.

<!-- field: traces -->
#### Traces

- srd007-release-annotation AC2

<!-- sc: S3 -->
### S3

A run leaves the ref set and the working tree unchanged, passing or
refusing.

<!-- field: traces -->
#### Traces

- srd007-release-annotation AC3

<!-- sc: S4 -->
### S4

The runbook and the skill call the check between annotation and push,
so no release path skips it silently.

<!-- field: traces -->
#### Traces

- srd007-release-annotation AC4

<!-- field: out_of_scope -->
## Out of scope

- Writing the summary; the session that releases owns the words.
- Retroactively fixing the degenerate annotations already pushed.
- Gating mage tag, the audit, or any build on annotation quality.
