---
id: srd014-spec-graph
title: Specification Graph Loader and Validator
depends_on: []
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

As the number of SRDs, use cases, and test suites grows, manually
verifying cross-references becomes error-prone. Broken touchpoints,
orphaned SRDs, uncovered requirements, and stale citations can hide
silently. A specification graph loads all artifacts into a labeled
property graph and runs automated consistency checks.

<!-- field: goals -->
## Goals

- G1: Parse SRDs, use cases, test suites, roadmaps, and SPECIFICATIONS.yaml into typed structs.
- G2: Build a labeled property graph connecting all artifacts by their relationships.
- G3: Validate cross-artifact consistency and report findings.
- G4: Corpus-level loading from a project directory.
- G5: Run deterministic jurist charter suites over a target directory.

<!-- req: R1 -->
## R1 - Parsing

<!-- req: R1.1 -->
### R1.1

ParseSRD reads a single SRD YAML file into an SRD struct with requirements, acceptance criteria, and dependencies.

<!-- req: R1.2 -->
### R1.2

ParseUseCase reads a use case YAML file into a UseCase struct with flow, touchpoints, and success criteria.

<!-- req: R1.3 -->
### R1.3

ParseTestSuite reads a test suite YAML file into a TestSuite struct with test cases and traces.

<!-- req: R1.4 -->
### R1.4

ParseRoadmap reads road-map.yaml into a Roadmap struct with releases and use case references.

<!-- req: R1.5 -->
### R1.5

ParseSpecIndex reads SPECIFICATIONS.yaml into a SpecIndex struct.

<!-- req: R1.6 -->
### R1.6

Tagged goals, non-goals, flows, out-of-scope items, and touchpoints SHALL reject unsupported or ambiguous YAML shapes with artifact and field diagnostics; every accepted list element SHALL produce exactly one ordered logical item.

<!-- req: R1.7 -->
### R1.7

ParseSRD reads the optional assumptions tagged list; ParseUseCase reads the optional assumptions and edge_cases tagged lists and the optional priority field (design.yaml use_case priority).

<!-- req: R2 -->
## R2 - Graph construction

<!-- req: R2.1 -->
### R2.1

Graph uses a directed labeled property graph with node kinds: release, srd, req-group, req-item, ac, use-case, test-suite, test-case.

<!-- req: R2.2 -->
### R2.2

Edge relationships include: contains, depends-on, traces, covers, tests.

<!-- req: R2.3 -->
### R2.3

BuildGraph populates the graph from a Corpus of all parsed artifacts.

<!-- req: R3 -->
## R3 - Validation checks

<!-- req: R3.1 -->
### R3.1

Detect orphaned SRDs not referenced by any use case.

<!-- req: R3.2 -->
### R3.2

Detect broken touchpoints referencing non-existent SRDs.

<!-- req: R3.3 -->
### R3.3

Detect broken citations in test cases referencing non-existent acceptance criteria.

<!-- req: R3.4 -->
### R3.4

Detect uncovered requirement items without test case traces.

<!-- req: R3.5 -->
### R3.5

Detect uncovered acceptance criteria without test case traces.

<!-- req: R3.6 -->
### R3.6

Detect depends_on violations referencing non-existent SRDs.

<!-- req: R3.7 -->
### R3.7

Detect releases without test suites.

<!-- req: R3.8 -->
### R3.8

Findings carry check identifier, severity level, human-readable message, and available provenance.

<!-- req: R3.9 -->
### R3.9

Detect releases in the road map that have no roadmap_summary row in SPECIFICATIONS.yaml.

<!-- req: R3.10 -->
### R3.10

Detect index entries whose title, status, or release disagrees with the document they name. Silence on either side is not disagreement: an index may summarize less than a document states, and a document may omit a field the index carries.

<!-- req: R3.12 -->
### R3.12

Detect ids the previous release tag contains that are absent at HEAD without a retirement entry, and retired entries without retired_in (retired-id-missing, design.yaml D6). The previous tag's corpus is read through the local checkout with read-only git calls: the latest release tag reachable from HEAD^ names the baseline, and a module without a repository, a release tag, or tagged specs yields nothing.

<!-- req: R3.13 -->
### R3.13

Warn when a touchpoint, citation, or trace names a retired id (citation-of-retired-id); the message names the replaced_by ids so the finding carries the fix.

<!-- req: R3.14 -->
### R3.14

Detect needs_clarification entries on SRDs, use cases, and test suites (unresolved-clarification, design.yaml clarification_marker). One error per entry, naming the document, the blocked id, and the question; an entry whose id names nothing in the document still fires, saying the target is unknown. An absent or empty list is silent.

<!-- req: R3.15 -->
### R3.15

Detect a use-case priority outside the P1/P2/P3 vocabulary (invalid-priority, design.yaml use_case priority). An absent priority means unranked and is silent.

<!-- req: R3.16 -->
### R3.16

Resolve every derived_from entry against the parent's sibling clone cache through read-only local git: a missing cache or a validated release absent from it is unverifiable-external-reference (error, message carries the clone or fetch command), and a cited SRD or id absent at the validated release is broken-external-reference (error). docs/external-checkouts.yaml overrides a cache location (design.yaml external_references).

<!-- req: R3.17 -->
### R3.17

Warn when a cited id is retired at the cache's newest release tag (retired-external-reference, message carries replaced_by), and when the cited item's text changed between the entry's validated release and the tracked pin in docs/external-pins.yaml (stale-external-reference). An unchanged item is silent -- a new parent release alone is not staleness -- and a validated release ahead of the tracked pin is also reported.

<!-- req: R3.18 -->
### R3.18

Warn on a derived_from entry without a release (unpinned-external-reference): legal during co-development, resolved against the cache HEAD, and refused by the release gate together with retired citations.

<!-- req: R4 -->
## R4 - Corpus loading

<!-- req: R4.1 -->
### R4.1

LoadCorpus reads all specification artifacts from a project directory tree.

<!-- req: R4.2 -->
### R4.2

SRD files are discovered under docs/specs/software-requirements/.

<!-- req: R4.3 -->
### R4.3

Use case files under docs/specs/use-cases/, test suites under docs/specs/test-suites/.

<!-- req: R4.4 -->
### R4.4

Road-map from docs/road-map.yaml, spec index from docs/SPECIFICATIONS.yaml.

<!-- req: R5 -->
## R5 - Jurist charter suites

<!-- req: R5.1 -->
### R5.1

Jurist charter YAML declares suite identity, check IDs, check kinds, severity, target globs, and check-specific configuration.

<!-- req: R5.2 -->
### R5.2

Supported charter check kinds include grep_check, ref_check, consistency_check, and spec_corpus.

<!-- req: R5.3 -->
### R5.3

Charter execution is deterministic, local, read-only, and does not call an LLM or the network.

<!-- req: R5.4 -->
### R5.4

The jurist target directory comes from runtime data such as --directory, while profile or request data names one or more charter suite paths.

<!-- req: R5.5 -->
### R5.5

Charter findings record suite ID, check ID, check kind, severity, message, and any available file and line provenance.

<!-- req: R5.6 -->
### R5.6

The built-in spec_corpus charter preserves the current specification corpus validation behavior.

<!-- req: R5.7 -->
### R5.7

grep_check target search runs through a declared rg exec word, followed by a machine-visible reduction step; pkg/spec does not scan target files for grep matches.

<!-- req: R5.8 -->
### R5.8

ref_check target extraction and file/directory inventory loading run through a declared external scan word, followed by a machine-visible reduction step; pkg/spec retains only plan lowering and finding reduction.

<!-- req: R5.9 -->
### R5.9

consistency_check YAML loading and required-path inventory traversal run through a declared external scan word; the reducer retains YAML-path semantics and line provenance without opening or statting target files.

<!-- field: non_goals -->
## Non-goals

- Generating specifications from code.
- Modifying specification files.
- UI or interactive reporting.
- Using LLMs, network APIs, or non-deterministic inputs during validation.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

ParseSRD, ParseUseCase, ParseTestSuite, ParseRoadmap, and ParseSpecIndex produce typed structs from YAML files without silently dropping malformed tagged-list elements.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R1.4
- R1.5
- R1.6
- R1.7

<!-- ac: AC2 -->
### AC2

BuildGraph constructs a directed labeled property graph from a Corpus with correct node kinds and edge relationships.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.3

<!-- ac: AC3 -->
### AC3

Validate detects orphaned SRDs, broken touchpoints and citations, uncovered items and ACs, depends_on violations, releases without test suites, releases missing a roadmap_summary row, index entries disagreeing with their documents, ids dropped without retirement entries, citations of retired ids, unresolved needs_clarification entries, invalid use-case priorities, and external references that are unverifiable, broken, retired, stale, or unpinned.

<!-- field: traces -->
#### Traces

- R3.1
- R3.2
- R3.3
- R3.4
- R3.5
- R3.6
- R3.7
- R3.8
- R3.9
- R3.10
- R3.12
- R3.13
- R3.14
- R3.15
- R3.16
- R3.17
- R3.18

<!-- ac: AC4 -->
### AC4

LoadCorpus discovers and parses all artifacts from the standard directory layout.

<!-- field: traces -->
#### Traces

- R4.1
- R4.2
- R4.3
- R4.4

<!-- ac: AC5 -->
### AC5

Jurist charter suites load from YAML, run deterministic check kinds, expose grep_check and ref_check as declared external scans plus finding reduction steps, and report provenance-rich findings.

<!-- field: traces -->
#### Traces

- R5.1
- R5.2
- R5.3
- R5.4
- R5.5
- R5.7
- R5.8
- R5.9

<!-- ac: AC6 -->
### AC6

The built-in spec_corpus charter reproduces the current specification corpus validation behavior used by mage audit.

<!-- field: traces -->
#### Traces

- R5.2
- R5.6

<!-- field: diagrams -->
## Diagrams

````yaml
requirements_trace: |
    ```mermaid
    flowchart LR
      Problem["Problem: As the number of SRDs, use cases, and test suites gro..."] --> Goals["Goals: Specification Graph Loader and Validator"]
      Goals --> R1["R1: Parsing"]
      Goals --> R2["R2: Graph construction"]
      Goals --> R3["R3: Validation checks"]
      Goals --> R4["R4: Corpus loading"]
      AC1["AC1: ParseSRD, ParseUseCase, ParseTestSuite, ParseRoadmap, Par..."]
      R1 --> AC1
      AC2["AC2: BuildGraph constructs a directed labeled property graph f..."]
      R2 --> AC2
      AC3["AC3: Validate detects orphaned SRDs, broken touchpoints and ci..."]
      R3 --> AC3
      AC4["AC4: LoadCorpus discovers and parses all artifacts from the st..."]
      R4 --> AC4
    ```
````
