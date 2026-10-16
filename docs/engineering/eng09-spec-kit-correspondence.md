---
id: eng09-spec-kit-correspondence
title: "Spec Kit Correspondence"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: introduction -->
## Introduction

GitHub's Spec Kit (https://github.com/github/spec-kit) packages
Spec-Driven Development as a sequence of agent commands over markdown
artifacts: constitution, specify, clarify, plan, checklist, tasks,
analyze, implement, converge. This repository practices the same
discipline over a typed YAML corpus with deterministic gates. The two
are functionally equivalent, and this guideline states the
correspondence so the claim can be checked row by row rather than
asserted.

The correspondence is functional, not structural. We do not mirror
Spec Kit's files; we show that every step of its workflow and every
element of its spec artifact has a mechanism here with the same job.
Where our mechanism is machine-checked and Spec Kit's is a markdown
convention, the difference favors the corpus. One step, converge, is
closed on its cross-repository axis and open on its intra-repository
one; the last section bounds what remains.

<!-- section: unit-of-correspondence -->
## Unit of correspondence

One Spec Kit spec.md maps to one use case together with the SRD
requirement slices its touchpoints cite and the test suite that
traces it. Spec Kit bundles user stories, scenarios, success
criteria, assumptions, and open questions into a single file per
feature; the corpus distributes the same information across typed
documents joined by the traceability graph, and mage audit checks
the joins.

A user story and a use case are not the same thing. A story is a
small, prioritized, independently shippable unit of scheduling; a
use case is the classical interaction contract with actor,
trigger, flow, touchpoints, and success criteria, and one use case
typically covers several stories. We put priority (P1, P2, P3) on
the use case rather than introduce a user-story document type:
releases schedule use cases, so that is the grain the roadmap
ranks. A story-level grain, if a release ever ships half a use
case, would be priorities on individual success criteria; nothing
in the formats blocks that extension.

<!-- section: workflow-correspondence -->
## Workflow correspondence

Table 1 maps every Spec Kit command to the mechanism here. The
optional quality gates (clarify, checklist, analyze) correspond to
mechanisms that are not optional in this repository: the audit
runs them on every gate.

Table 1: Spec Kit commands and their mechanisms here

| Spec Kit command | Mechanism here |
|------------------|----------------|
| speckit.constitution | docs/constitutions/*.yaml; the audit enforces articles rather than restating them |
| speckit.specify | SRD and use-case formats with numbered requirements (design.yaml document_types) |
| speckit.clarify | needs_clarification entries plus mage clarify, the pre-planning worklist; the unresolved-clarification check refuses an unresolved corpus |
| speckit.plan | ARCHITECTURE documents plus the issue body: required_reading, files, requirements, design_decisions |
| speckit.checklist | Completeness checks C1 through C7 in design.yaml, enforced by the audit |
| speckit.tasks | Issue push/pop with sub-issues on an epic, or beads in bd-mode repositories |
| speckit.analyze | mage audit: the spec_corpus checks and the traceability graph |
| speckit.implement | The do-work loop: one unit at a time, acceptance criteria verified at finish |
| speckit.converge | Cross-repository: the external-reference upgrade loop (GH-2557, GH-2558); intra-repository code-vs-spec: no automated equivalent, see the last section |
| speckit.taskstoissues | Native: tasks are issues from the start |

<!-- section: artifact-correspondence -->
## Artifact correspondence

Table 2 maps the elements of a Spec Kit spec to corpus fields and
standards. Every row lands on a parsed field or a stated standard,
so the audit can see it; Spec Kit checks the same elements by
prompt convention.

Table 2: Spec artifact elements and their corpus forms

| Spec Kit element | Corpus form |
|------------------|-------------|
| User stories with P1/P2/P3 | Use cases with the optional priority field; invalid-priority gates the vocabulary |
| Given/When/Then scenarios | Test cases with traces to acceptance criteria; checkable rather than narrative |
| Success criteria | success_criteria entries under the documentation standard: measurable, technology-agnostic |
| Assumptions | assumptions tagged lists (A1, A2) on SRDs and use cases |
| Edge cases | edge_cases tagged lists (E1, E2) on use cases |
| NEEDS CLARIFICATION markers | needs_clarification entries ({id, question}) on corpus documents; tracker issues keep the inline literal marker |
| Review checklist | Completeness checks C1 through C7 |
| Stable requirement anchors | D6 immutable identifiers with retirement entries; Spec Kit has no equivalent rule |

<!-- section: markdown-format-re-check -->
## Markdown format re-check

The corpus is migrating from YAML to Markdown with marker
comments (eng10). Spec Kit's artifacts are Markdown checked by
convention, so the migration could have turned every row of
Table 2 back into a convention. It does not: Table 3 maps each
row to the parsed construct that carries it in the new format and
the finding that reports it when it breaks. Every row lands on a
construct; none degrades to convention.

Table 3: Spec artifact elements in the Markdown format

| Spec Kit element | Parsed construct | Finding when broken |
|------------------|------------------|---------------------|
| User stories with P1/P2/P3 | use case frontmatter priority | frontmatter-unknown-key; invalid-priority for the value |
| Given/When/Then scenarios | test case nodes (case marker) with a required visible traces field | missing-required-field, field-shape-mismatch |
| Success criteria | success criterion nodes (sc marker, S1 ids) with prose and a required traces field | id-malformed, missing-content, missing-required-field |
| Assumptions | assumptions list field with item pattern A1, A2 | list-item-pattern, field-shape-mismatch |
| Edge cases | edge_cases list field with item pattern E1, E2 | list-item-pattern, field-shape-mismatch |
| NEEDS CLARIFICATION markers | needs_clarification table field with columns id, question | table-columns-mismatch; unresolved-clarification |
| Review checklist | completeness checks C1 through C7 over the same model | the existing completeness findings, unchanged |
| Stable requirement anchors | req and ac markers with declared id patterns, unique ids, and retirement attributes | id-malformed, id-outside-parent, duplicate-id, attribute-shape; retired-id-missing |

The graph findings named after the semicolons are the audit's
existing checks; the rest are the format's own catalog in eng10.

<!-- section: differences-and-the-converge-gap -->
## Differences and the converge gap

Three differences remain after the mapping, and only one is a gap.

Enforcement. Spec Kit's artifacts are markdown read by prompts; a
drifting spec is found when an agent happens to notice. Here the
same properties are parsed fields and graph edges, and mage audit
refuses the corpus when they break: renumbered ids
(retired-id-missing), citations of retired ids, unresolved
clarifications, invalid priorities, untraced requirements. The
correspondence holds element by element, and the corpus side is
the stronger enforcement of the two.

The claim survives the move to Markdown (eng10, GH-2578). A
Markdown corpus document parses into the same document model, so
the graph checks above run unchanged, and the parse itself is
enforced: every element is a marker-tagged node or field whose
grammar, structure, and shape design.yaml declares, and a
document that breaks a declaration is a named finding with a
file, a line, and a fix hint. Prose containment closes the last
route around the parse, because a normative keyword outside a
requirement item is itself a finding. Table 3 re-checks Table 2
against the new format.

Task store. Spec Kit appends tasks to a repo-local tasks.md, so
the queue is versioned with the branch. Our backlog store is
pluggable: GitHub issues with sub-issue epics, or the beads
database in bd-mode repositories, which is the closer analogue to
tasks.md. The durable record does not depend on the tracker
either way: pop and work commits embed the issue bodies verbatim,
so the planning record survives a tracker migration in git.

Converge. Spec Kit's closing step re-reads spec, plan, and tasks
against the code and either reports Converged or appends the gaps
as new tasks; it is append-only and touches nothing but tasks.md.
We hold the deterministic approximations -- uncovered-req-item and
uncovered-ac prove every requirement traces to a test, and do-work
verifies acceptance criteria before a unit closes -- but the
judgment "does the code realize this requirement" needs a model,
and srd014 R5.3 keeps model calls out of the audit by design. An
equivalent here is therefore a skill or an agent profile beside
the gate, filing its gaps as sub-issues the way do-work already
files follow-up work: the model proposes, the audit enforces.

One axis of converge is closed. The external-reference upgrade
loop (GH-2557, GH-2558) is converge across the repository
boundary: stale-external-reference findings are deterministic
gap detection between a child corpus and its parent's releases,
mage upgrade performs the mechanical re-pins, upgrade:status
prints the remaining gaps as a worklist that becomes tracker
issues, and the loop terminates when the audit is clean at the
tracked pin -- implement until converged, Spec Kit's exact
shape, with the tracker in place of tasks.md. What stays open is
the intra-repository half: the judgment that this repository's
own code realizes its own requirements, which srd014 R5.3 keeps
outside the deterministic audit.

<!-- field: references -->
## References

- agent-core/docs/constitutions/design.yaml
- docs/engineering/eng10-markdown-spec-format.md
- agent-core/docs/constitutions/issue-format.yaml
- agent-core/docs/specs/software-requirements/srd014-spec-graph.yaml
- https://github.com/github/spec-kit
- https://github.github.io/spec-kit/reference/agentic-sdd.html
- https://github.com/Nokia-Bell-Labs/declarative-agents/issues/2545
- https://github.com/Nokia-Bell-Labs/declarative-agents/issues/2549
- https://github.com/Nokia-Bell-Labs/declarative-agents/issues/2551
- https://github.com/Nokia-Bell-Labs/declarative-agents/issues/2552
- https://github.com/Nokia-Bell-Labs/declarative-agents/issues/2557
- https://github.com/Nokia-Bell-Labs/declarative-agents/issues/2558
- https://github.com/Nokia-Bell-Labs/declarative-agents/issues/2581
