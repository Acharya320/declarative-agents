---
id: rel17.0-uc001-pin-currency
title: Read Which Pins Moved and Which Aged, Then Decide
status: planned
test_suite: test-rel17.0-pin-currency
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: release -->
## Release

17.0

<!-- field: summary -->
## Summary

An operator runs one target and learns two things about every pinned
upstream dependency: whether its recorded digest is still what its tag
resolves to, and whether a newer release exists. A moved digest is reported
first because it means something was re-pushed underneath a pin. A pin that
has merely aged is reported with the version it is behind. Nothing is
edited; which bumps to take is the operator's call, and the report exists
so the call is informed rather than overdue.

<!-- field: actor -->
## Actor

An operator or maintainer deciding which upstream versions to take.

<!-- field: trigger -->
## Trigger

A release is being prepared, a dependency is suspected stale, or the
operator is doing periodic maintenance.

<!-- field: preconditions -->
## Preconditions

- ENG01 carries the pin currency section and srd006 owns its requirements.
- srd005 and its gate already guarantee every image is pinned.
- The host can reach the registries; one that cannot be reached is reported, not assumed.

<!-- field: flow -->
## Flow

- F1: The operator runs the bump target at the repository root.
  - traces:
    - srd006-pin-currency R4.2
- F2: The survey builds its inventory: chart pins read out of every chart
  rendered under every checked-in values overlay, and the remaining pins
  from the checked-in list naming the file and symbol that holds each.
  - traces:
    - srd006-pin-currency R1.1
    - srd006-pin-currency R1.2
- F3: Each pin is resolved against its registry. The tag's current digest
  answers integrity; the registry's tag list answers currency.
  - traces:
    - srd006-pin-currency R2.1
    - srd006-pin-currency R3.1
- F4: A recorded digest that no longer matches its tag is reported as moved,
  ahead of everything else, with both digests.
  - traces:
    - srd006-pin-currency R2.1
    - srd006-pin-currency R6.1
- F5: A pin behind a newer release is reported with the pinned version and
  the newest one. A tag that is not a semantic version is compared on
  integrity alone.
  - traces:
    - srd006-pin-currency R3.1
    - srd006-pin-currency R3.2
- F6: A registry that cannot be reached, or a tag that is absent, is reported
  as unknown. The operator can tell the difference between a pin that is
  current and one nobody could check.
  - traces:
    - srd006-pin-currency R5.1
- F7: The target exits zero and the working tree is unchanged. The operator
  edits the pins they choose to take, at the locations the report named.
  - traces:
    - srd006-pin-currency R4.1
    - srd006-pin-currency R6.1

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | docs/engineering/eng01-kind-test-demo-rig.md | States the two pin properties and the report-never-mutate contract for the verb. |
| T2 | magefiles/chartconf/ | Renders the charts whose images become the derived half of the inventory. |
| T3 | magefiles/pinsurvey/ | Registry client, inventory, verdicts, and report. |
| T4 | magefiles/kindrig/ | Holds the add-on digests and the node image the declared half enumerates. |
| T5 | [srd006-pin-currency](../software-requirements/srd006-pin-currency.md) [AC1](../software-requirements/srd006-pin-currency.md#ac1) | Inventory coverage, derived and declared. |
| T6 | [srd006-pin-currency](../software-requirements/srd006-pin-currency.md) [AC2](../software-requirements/srd006-pin-currency.md#ac2) | Integrity reported first. |
| T7 | [srd006-pin-currency](../software-requirements/srd006-pin-currency.md) [AC3](../software-requirements/srd006-pin-currency.md#ac3) | Currency and version comparison. |
| T8 | [srd006-pin-currency](../software-requirements/srd006-pin-currency.md) [AC4](../software-requirements/srd006-pin-currency.md#ac4) | Reports, never writes, and stays out of the audit. |
| T9 | [srd006-pin-currency](../software-requirements/srd006-pin-currency.md) [AC5](../software-requirements/srd006-pin-currency.md#ac5) | Unknown is not current. |
| T10 | [srd006-pin-currency](../software-requirements/srd006-pin-currency.md) [AC6](../software-requirements/srd006-pin-currency.md#ac6) | A finding names where to act. |
| T11 | [test-rel17.0-pin-currency](../test-suites/test-rel17.0-pin-currency.md) | Structural and executable evidence for this flow. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Every pinned upstream dependency appears in the report, chart pins
derived from the charts and the rest from a list verified against the
locations it names.

<!-- field: traces -->
#### Traces

- [srd006-pin-currency](../software-requirements/srd006-pin-currency.md) [AC1](../software-requirements/srd006-pin-currency.md#ac1)

<!-- sc: S2 -->
### S2

A re-pushed tag is reported as moved, ahead of any staleness, with the
recorded and current digests.

<!-- field: traces -->
#### Traces

- [srd006-pin-currency](../software-requirements/srd006-pin-currency.md) [AC2](../software-requirements/srd006-pin-currency.md#ac2)

<!-- sc: S3 -->
### S3

A pin behind upstream is reported with both versions, and a tag that
names no version is not compared.

<!-- field: traces -->
#### Traces

- [srd006-pin-currency](../software-requirements/srd006-pin-currency.md) [AC3](../software-requirements/srd006-pin-currency.md#ac3)

<!-- sc: S4 -->
### S4

The run leaves the working tree unchanged, exits zero, and never runs
inside the audit.

<!-- field: traces -->
#### Traces

- [srd006-pin-currency](../software-requirements/srd006-pin-currency.md) [AC4](../software-requirements/srd006-pin-currency.md#ac4)

<!-- sc: S5 -->
### S5

A pin nobody could check reads as unknown rather than current, and the
report names where every pin lives.

<!-- field: traces -->
#### Traces

- [srd006-pin-currency](../software-requirements/srd006-pin-currency.md) [AC5](../software-requirements/srd006-pin-currency.md#ac5)
- [srd006-pin-currency](../software-requirements/srd006-pin-currency.md) [AC6](../software-requirements/srd006-pin-currency.md#ac6)

<!-- field: out_of_scope -->
## Out of scope

- Taking any bump, editing any pin, or opening a branch or pull request.
- Classifying an update as safe, noteworthy, or breaking from its release notes.
- Go module versions, which the language's own tooling reports.
- Failing a build on a stale pin.
