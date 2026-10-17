---
id: srd006-pin-currency
title: Pin Integrity and Currency
status: planned
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

srd005 makes every image pinned and the audit enforces it. Nothing checks
that a pin is still true. Two properties go unwatched, and they fail in
opposite directions. A tag re-pushed upstream leaves a recorded digest
describing an image nobody chose, which is a supply-chain event. A pin left
alone while upstream moves on ages silently until something forces the
update; at the time of writing the ingress controller is three patch
releases behind and the cluster node image a minor behind, and no output
anywhere says so. The pins also live in two kinds of place — declared in
chart values, and compiled into rig source or repeated across cluster
configurations — so any answer has to cover both without maintaining a list
that drifts from what it describes.

<!-- field: goals -->
## Goals

- G1: Report whether each recorded digest still resolves from its tag.
- G2: Report whether a newer upstream release exists for each pin.
- G3: Derive the chart half of the inventory from the charts, and declare only what no render exposes.
- G4: Report and change nothing, so a version bump stays a decision a person takes.
- G5: Distinguish an unreachable registry from a current pin.

<!-- field: applicability | form: "tree" -->
## Applicability

- baseline:
  R1 through R6 apply to every pinned upstream dependency this repository
  resolves from a registry. Go module versions are outside this scope; the
  language's own tooling reports them.
- capability_gates:
  - helm_managed: The derived half of the inventory applies to an application that ships a chart.
  - kind_demo: The declared half covers the rig add-ons and cluster configurations a kind demo uses.

<!-- req: R1 -->
## R1 - The inventory covers both kinds of pin

<!-- field: applicability -->
### Applicability

runnable_module

<!-- req: R1.1 -->
### R1.1

Chart image pins shall be read from the rendered charts rather than
restated in any list, so the inventory cannot disagree with the
charts it describes.

<!-- field: evidence -->
#### Evidence

The survey renders every chart under every checked-in values
overlay and takes its chart pins from the resulting manifests,
reusing the rendering the conformance gate already performs.

<!-- req: R1.2 -->
### R1.2

Pins no render exposes shall be enumerated in one checked-in list.
Each entry shall name the file and the symbol or field holding the
pin, the image, the tag, and the recorded digest.

<!-- field: evidence -->
#### Evidence

The list covers the rig add-on digests compiled into the kindrig
package, the node image in each checked-in cluster configuration,
and the smoke dependency registry.

<!-- req: R1.3 -->
### R1.3

A checked-in entry whose named location no longer holds the pin it
records shall fail a test, so the list cannot rot into fiction
while the report keeps citing it.

<!-- field: evidence -->
#### Evidence

Moving or renaming a pinned constant without updating the list
fails, naming the entry and the location it no longer describes.

<!-- req: R2 -->
## R2 - Integrity is reported first and on its own

<!-- field: applicability -->
### Applicability

runnable_module

<!-- req: R2.1 -->
### R2.1

A recorded digest that is not what its tag resolves to upstream
shall be reported as moved, separately from and ahead of any
staleness.

<!-- field: evidence -->
#### Evidence

A tag re-pushed upstream produces a moved verdict naming the
recorded digest and the one the tag resolves to now, and that
verdict appears before currency findings in the report.

<!-- req: R3 -->
## R3 - Currency is reported with the version behind

<!-- field: applicability -->
### Applicability

runnable_module

<!-- req: R3.1 -->
### R3.1

A pin with a newer upstream release shall be reported as behind,
naming the pinned version and the newest available one.

<!-- field: evidence -->
#### Evidence

A registry offering a higher version than the pin produces a behind
verdict carrying both versions.

<!-- req: R3.2 -->
### R3.2

Version comparison shall apply where the tag is a semantic version,
with or without a leading v, and shall be skipped otherwise. A
pre-release or release-candidate tag shall not count as newer than
a release.

<!-- field: evidence -->
#### Evidence

A tag that names no version, such as a commit revision or a
fail-closed sentinel, reports on integrity alone, and a
release-candidate tag upstream does not produce a behind verdict.

<!-- req: R4 -->
## R4 - The report changes nothing

<!-- field: applicability -->
### Applicability

runnable_module

<!-- req: R4.1 -->
### R4.1

The target shall edit no file, create no branch, and make no
commit. It shall exit zero whatever it finds.

<!-- field: evidence -->
#### Evidence

The working tree is byte-identical across a run, and a run finding
moved and behind pins still exits zero.

<!-- req: R4.2 -->
### R4.2

The target shall reach the network and shall therefore not run as
part of the repository audit, which stays hermetic.

<!-- field: evidence -->
#### Evidence

The audit does not invoke the target, and the target's own tests
reach no network.

<!-- req: R5 -->
## R5 - An unreachable registry is unknown, never current

<!-- field: applicability -->
### Applicability

runnable_module

<!-- req: R5.1 -->
### R5.1

A pin whose registry cannot be reached, or whose tag is absent
upstream, shall be reported as unknown and shall never be reported
as current.

<!-- field: evidence -->
#### Evidence

A registry returning an error produces an unknown verdict naming
the reason, and no unknown pin is counted among the current ones.

<!-- req: R6 -->
## R6 - A finding names where to act

<!-- field: applicability -->
### Applicability

runnable_module

<!-- req: R6.1 -->
### R6.1

Every reported pin shall name its location, the value recorded
there, and the upstream value it was compared against.

<!-- field: evidence -->
#### Evidence

A reader edits the right line from the report alone, without
searching the repository for where the pin lives.

<!-- field: non_goals -->
## Non-goals

- Taking a bump. The report informs a decision; a person makes it.
- Reading release notes or classifying an update as safe, noteworthy, or breaking. That judgment needs a word class this repository does not have.
- Go module versions. The language's own tooling reports those.
- Gating. A stale pin is a fact to surface, not a build failure.
- Rewriting any pin, opening any branch, or drafting any pull request.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

The inventory covers chart pins derived from rendered charts and
non-chart pins from one checked-in list whose entries are verified
against the locations they name.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3

<!-- ac: AC2 -->
### AC2

A recorded digest that no longer resolves from its tag is reported as
moved, ahead of any currency finding.

<!-- field: traces -->
#### Traces

- R2.1

<!-- ac: AC3 -->
### AC3

A pin with a newer upstream release is reported as behind with both
versions, semantic versions are compared and other tags are not, and a
pre-release does not count as newer.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2

<!-- ac: AC4 -->
### AC4

The target leaves the working tree unmodified, exits zero whatever it
finds, and is not part of the repository audit.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2

<!-- ac: AC5 -->
### AC5

An unreachable registry or absent tag reports unknown and is never
counted as current.

<!-- field: traces -->
#### Traces

- R5.1

<!-- ac: AC6 -->
### AC6

Every reported pin names its location, its recorded value, and the
upstream value.

<!-- field: traces -->
#### Traces

- R6.1

<!-- field: use_case_traces -->
## Use case traces

- rel17.0-uc001-pin-currency

<!-- field: test_suite_traces -->
## Test suite traces

- test-rel17.0-pin-currency
