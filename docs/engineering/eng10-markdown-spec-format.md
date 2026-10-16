---
id: eng10-markdown-spec-format
title: Markdown Specification Format
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: introduction -->
## Introduction

This guideline specifies the Markdown format that replaces the typed YAML specification corpus (epic GH-2578). A document in this format renders as a document on GitHub and GitLab, and it parses into the same document model the YAML front end feeds, so the traceability graph and every audit check keep working unchanged. Each heading carries a one-line HTML comment marker that says what the heading opens. Facts a reader needs stay in visible Markdown, and machine-only facts live in the markers.

We wrote the format from evidence, not from argument. The prototype in tools/specmd-prototype parses, prints, and diagnoses the format, and tools/specmd-prototype/FINDINGS.md records what it established: round-trip laws over the 106 tracked SRDs, the inventory of shape rules, and the quality of the diagnostics. The rules this guideline describes are data. Each document type declares them in the markdown block of its entry in agent-core/docs/constitutions/design.yaml, and the prototype reads its rules from nowhere else.

This guideline is itself written in the format. Running specmd check on this file reports zero findings, which makes the file its first conformance test.

<!-- section: one-layer-rule -->
## The one-layer rule

A fact lives in exactly one layer. Reader-relevant facts live only in visible Markdown: headings, paragraphs, lists, tables, and the frontmatter. Machine-only facts live only in marker comments: a requirement's weight, the lifecycle status and retirement pointers of a node, and provenance references. No fact appears in both layers, so the two cannot drift apart.

The frontmatter counts as visible. GitHub renders it as a table above the document, and GitLab shows it as a code box, so identity and lifecycle fields such as id, title, status, priority, and release belong there. Everything else a reader needs is a tagged field under a heading.

We applied the rule to criterion traces (GH-2581 decision 3). A reader must see which requirements a criterion verifies, so traces are a visible field under each criterion, not a marker attribute. The same holds for use case success criteria and for test cases.

<!-- section: file-layout -->
## File layout

A document is one Markdown file at the location its type declares in markdown.location. It opens with the frontmatter, then the two license comment lines, then its nodes in order. The skeleton of an SRD follows.

```markdown
---
id: srd004-declarative-ux
title: Declarative UX
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Prose under the heading.

<!-- req: R1 -->
## R1 - Contract

<!-- req: R1.1 | weight: 3 -->
### R1.1

The requirement text.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

The criterion text.

<!-- field: traces -->
#### Traces

- R1.1
```

The frontmatter lists the fields frontmatter.required names and admits only those frontmatter.allowed names. Where a type sets id_matches_filename, the id equals the file name without its extension. The license comments are plain comments, not markers, and the copyright gate checks them as it does for any Markdown file.

During the migration window a document exists as .yaml or as .md, never both; GH-2583 adds the audit finding that enforces this. After the sunset (GH-2585) only the Markdown locations remain.

<!-- section: marker-grammar -->
## Marker grammar

A marker is a single-line HTML comment directly above the heading it tags. It holds a tag and a value, then any attributes, each segment separated by a space, a vertical bar, and a space. Attribute values are JSON: strings are quoted, lists are arrays, and provenance references are arrays of objects. Table 1 gives the grammar.

Table 1: Marker grammar

| Part | Form | Example |
| --- | --- | --- |
| Marker | `<!-- <tag>: <value>[ \| <attr>: <json>]... -->` | `<!-- req: R1.1 \| weight: 3 -->` |
| Tag | one of the tags the document type declares | field, req, ac |
| Value | the node's id, or the field name for a field marker | R1.1, acceptance_criteria |
| Attribute | a name the node declares in its attributes list | status, retired_in, replaced_by |
| Attribute value | JSON in the shape declared under attributes | `["srd010 R2.1"]` |

Binding is positional. A marker binds to the next heading, and only to it; a marker followed by another marker or by the end of the document is an orphan. The tag set is closed per type: a tag the type does not declare is a marker-unknown-tag finding with the nearest declared tag as its hint. Across all declared types the tags are field, req, ac, sc, case, section, release, pattern, interface, component, decision, and article.

The issue that commissioned this guideline named the tags req, field, acceptance, and copyright. The prototype settled on ac for acceptance criteria, matching the AC1 ids, and it treats the copyright and SPDX lines as license comments rather than markers, so no tag exists for them.

<!-- section: heading-grammar -->
## Heading grammar

A document type declares its nodes. Each node declaration names a tag, a heading level, a parent, an optional id pattern, a heading form, a content kind, whether ids are unique, the attributes it admits, and the fields it requires. Table 2 lists the parts.

Table 2: Node declaration

| Part | Meaning |
| --- | --- |
| tag | the marker tag that opens the node |
| level | the heading depth, 2 to 6 |
| parent | document, another node kind, or field:<name> for a node under a field |
| id_pattern | a regular expression over the marker value |
| id_extends_parent | the id starts with the parent's id and a dot, as R1.1 under R1 |
| heading | id (the heading is the id), id_title (the id, or the id, a spaced hyphen, and a title), or free |
| content | none (only tagged children), prose (text, no untagged list or table), or any |
| unique_id | no two nodes of the kind share an id |
| attributes | the marker attributes the node admits |
| required_fields | the fields every node of the kind carries |

A node's parent is the node open one level up, or field:<name> when that node is a field. A heading at the wrong depth, or under a parent its declaration does not name, is heading-misplaced. The document as a whole declares the fields it requires (required_fields) and the minimum count of each node kind (required_nodes); a status listed in exempt_status, such as relocated or retired, lifts the content requirements from a stub.

<!-- section: field-shapes -->
## Field shapes

A field is a node opened by a field marker; its name is the marker value. A typed field declares where it may appear and the shape of its content. Untyped fields are passthrough: they hold any content, so a document keeps every field its YAML form had. Table 3 lists the shapes.

Table 3: Field shapes

| Shape | Content | Constraints |
| --- | --- | --- |
| prose | paragraphs, or one fenced yaml block | none beyond the shape |
| list | one unordered list with hyphen markers | item_pattern on each item, when declared |
| table | one GFM table | columns, when declared, exactly and in order |
| empty | nothing; the field only holds child nodes | none |
| yaml | one fenced yaml block | none |

A list or table field with no entries still has its shape, so an empty index is not a finding. A table field that declares no columns accepts any header, which the index documents need because their columns differ by module.

How each YAML value becomes one of these shapes is the job of the YAML mapping, in the next section. The fenced block is the fallback, and the SRD declaration needs it nowhere.

<!-- section: yaml-mapping -->
## YAML mapping

One converter serves every document type (GH-2582, split on 2026-10-16). It maps a Markdown document to the yaml.Node tree its YAML form parses into, and back, driven only by the type's declaration. Every consumer that reads the YAML tree today reads the same tree from a Markdown document: the typed parsers (ParseSRD, ParseUseCase, ParseTestSuite, ParseRoadmap, ParseSpecIndex), the corpus charter's YAML-path checks, and the reference checks in magefiles. The typed models stay the one document model; Markdown is a second way to reach them.

The rendering rules above say how a document looks. The YAML mapping says where each piece lives in the tree, and it sits beside them in the declaration. Frontmatter fields are top-level scalar keys of the same name. A field maps to the record key of its name unless its declaration names another, and its value converts by Table 5. A node kind maps to records by the parts in Table 4.

Table 4: YAML mapping of a node kind

| Part | Meaning |
| --- | --- |
| from | the key, in the parent's record or at the top level, that holds the records |
| records | list (a sequence of records) or map (records keyed by id, as SRD requirement groups are) |
| id | the record key that holds the id; for map records the id is the map key |
| title | the record key whose value follows the id in an id_title heading |
| content | the record key whose value is the node's prose |
| keyed | the record may be written as a one-entry map from id to content |
| fields | the child field kind that carries every other record key |

Table 5: Value conversions

| YAML value | Markdown | Used for |
| --- | --- | --- |
| string | prose | problem, summary, purpose |
| list of strings | list | non_goals, references |
| map of strings | keyed list: each item reads name, colon, text | vision success criteria |
| list of keyed records | outline: each item reads id, colon, content, with every other record key as a nested item | goals, flow steps |
| list of flat records | table, one column per record key | touchpoints, needs_clarification, index tables |
| list of records with prose | nodes, one per record | requirement items, acceptance criteria, test cases |
| anything else | fenced yaml block | nested passthrough values |

A passthrough value, one no declaration types, takes the first conversion in the table that parses back to the same node: a string becomes prose, a list of plain strings becomes a list, and anything else becomes a fenced yaml block. Record keys keep their order, so a document's key order survives as its field order.

<!-- section: round-trip-laws -->
## Round-trip laws

The laws are stated on the node tree, with comments removed and after the one-time form normalization of the next section. For every tracked document of every structured type, yaml to node to Markdown to node is the identity, Markdown is a fixed point of parse then print on canonical form, and the converter is idempotent. For the types that have a typed model, the model parsed from the Markdown front end equals the model parsed from the YAML one. GH-2605 proves the laws on the real corpus, and GH-2606 adds the batch converter and its inverse.

<!-- section: comments-and-forms -->
## Comments and two-form fields

YAML comments do not survive conversion. The tracked corpus holds 270 non-license comment lines in 35 files of the structured types, and they fall into two groups. Some carry facts a reader needs, such as the note in four catalog SRDs that the agent-core SRDs they cite are platform dependencies. Others restate structure: 10 repeat the document title on its first line, and 26 are blank or divider lines. Neither group belongs in a marker, since the one-layer rule keeps markers for machine-only facts, and placing a free note inside a typed field would break its shape. The converter therefore reports every comment it does not carry, with file and line, and the corpus migration (GH-2584) turns each comment that states a fact into visible text in the same pull request, called out as a hand edit.

Some fields are written in more than one YAML form, and one Markdown form cannot parse back to two. Before the round trips run, each such field is normalized to one canonical form in its own commit that changes no value: a record holding only its id and content takes the keyed form, any other record keeps the record form, and an unnumbered entry receives the next free number. Table 6 counts the entries the normalization touches.

Table 6: Fields normalized before conversion

| Field | Forms in the corpus | Entries normalized |
| --- | --- | --- |
| SRD goals | 404 keyed, 33 records with id and text, 36 numbered strings | 69 become keyed |
| SRD requirement groups | 104 SRDs as a map, 2 as a list | 2 SRDs become a map |
| SRD requirement items | 2467 keyed, 160 records with more keys | none |
| use case flow | 403 keyed, 262 records with id and text, 47 records with traces, 1 unnumbered string | 262 become keyed, 1 is numbered |
| use case touchpoints | 309 keyed strings, 272 records, 17 unnumbered strings | 309 split into records, 17 numbered |

The keyed touchpoint strings split at their double hyphen into target and reason, and every one of the 309 has one. Lists inside passthrough values keep whatever forms they have, because the fenced yaml block holds any node exactly.

<!-- section: identifiers -->
## Identifiers

Identifiers follow D6 in design.yaml: an id a release contains is never renumbered or reused, and a retired node keeps its id with retired_in and replaced_by attributes pointing at its successor. The format only adds where the id sits: in the marker value, and for heading forms id and id_title also in the heading text, where heading-id-mismatch keeps the two equal.

Requirement item ids match the pattern below (GH-2581 decision 1). The two-level form R1.1 is canonical for new items. The wider forms exist because 18 ids in 3 SRDs already use them (R3.1b, R10.0.1, R3.1.1), and D6 forbids renumbering them.

```text
^R[0-9]+(\.[0-9]+){1,2}[a-z]?$
```

A converter assigns new numbers only where the YAML had none: 17 unnumbered use case touchpoints in 3 files and 1 unnumbered flow step. Those are new ids, so no existing id changes.

<!-- section: diagnostics -->
## Diagnostics

An agent authoring a document iterates write, audit, fix, so every failure is a named finding that says how to fix it. Diagnostics come in tiers, and the run collects every finding rather than stopping at the first. Table 7 gives the tiers.

Table 7: Diagnostic tiers

| Tier | Covers | Handling |
| --- | --- | --- |
| 0 | canonical form: whitespace, spacing, marker spelling, table alignment | the formatter fixes it; never reported |
| 1 | grammar: document type, frontmatter, markers | named finding |
| 2 | structure: headings, nodes, ids, fields, required content | named finding |
| 3 | shapes: field shapes, list item patterns, table columns, attribute shapes | named finding |
| 4 | prose containment | named finding |

Each finding carries its name, tier, file, line, what the rule expected, what the document has, the rule path inside design.yaml (document_types.<type>.markdown.<rule>), and a fix hint. The names are stable: tools and agents match on them, so renaming one is a breaking change. Table 8 is the complete catalog.

Table 8: Finding catalog

| Tier | Name | Reported when |
| --- | --- | --- |
| 1 | document-type-unknown | no declared Markdown location matches the file's path |
| 1 | frontmatter-missing | the file does not open with a closed frontmatter block |
| 1 | frontmatter-invalid | the frontmatter is not a YAML mapping |
| 1 | frontmatter-missing-key | a field frontmatter.required names is absent |
| 1 | frontmatter-unknown-key | a frontmatter field is not in frontmatter.allowed |
| 1 | frontmatter-id-mismatch | the id differs from the file name where the type requires them equal |
| 1 | marker-malformed | a marker lacks the tag-colon-value form or an attribute is malformed |
| 1 | marker-unknown-tag | the tag is not one the type declares |
| 1 | marker-unknown-attribute | the attribute is not one the node declares |
| 1 | marker-orphan | no heading follows the marker |
| 2 | heading-untagged | a heading has no marker above it |
| 2 | heading-misplaced | the heading depth or parent does not fit any node declaration |
| 2 | id-malformed | the id does not match the node's id pattern |
| 2 | id-outside-parent | the id does not extend its parent's id |
| 2 | heading-id-mismatch | the heading text disagrees with the marker id |
| 2 | duplicate-id | an id the type declares unique appears twice |
| 2 | field-disallowed | a typed field appears under a node its declaration does not list |
| 2 | field-duplicate | a node holds the same field twice |
| 2 | content-untagged | content sits under a node whose content kind is none, or before any heading |
| 2 | untagged-structure | a prose node holds a list or table that belongs in its own field |
| 2 | missing-content | a prose node is empty |
| 2 | missing-required-field | a document or node lacks a field it requires |
| 2 | missing-required-node | the document holds fewer nodes of a kind than required |
| 3 | field-shape-mismatch | a typed field's content is not its declared shape |
| 3 | list-item-pattern | a list item does not match the field's item pattern |
| 3 | table-columns-mismatch | a table's header differs from the declared columns |
| 3 | attribute-shape | an attribute value is not JSON of its declared shape |
| 4 | prose-containment | a normative keyword appears in prose outside the normative nodes |

The graph and semantic checks the audit already runs (retired-id-missing, uncovered-req-item, unresolved-clarification, invalid-priority, and the rest) are unchanged. They run over the same model, so their names and coverage carry over to the Markdown corpus.

<!-- section: prose-containment -->
## Prose containment

Nothing normative exists only in prose. A normative keyword in a paragraph or list outside the nodes a type declares normative is a prose-containment finding: the obligation either moves into a requirement item, or the text is reworded as description. Code blocks and tables are not prose, so they do not trip the check. The keyword list is data, declared per type in prose_containment.keywords, and it holds the RFC 2119 keywords in capitals plus the lowercase forms below.

```text
MUST, MUST NOT, REQUIRED, SHALL, SHALL NOT, SHOULD, SHOULD NOT,
RECOMMENDED, MAY, OPTIONAL, shall, shall not
```

The normative nodes are the SRD requirement items and the constitution articles; everywhere else the check applies. Non-goals stay inside it (GH-2581 decision 2): a non-goal phrased as an obligation is a requirement nobody traces, or a scope statement in the wrong voice.

The check reports what the corpus migration (GH-2584) rewrites. On the printed SRDs it finds 15 hits in 5 files, 14 of them in non-goals. A scan of the YAML sources of the other types finds 2 hits in one engineering guideline (eng02 quoting srd020 R7.4), 2 in one constitution (go-style), and 61 in 4 config formats, 56 of them in rest-tool-format.

<!-- section: rules-as-data -->
## Rules as data

Every entry in document_types carries a markdown entry, and a test in the prototype holds design.yaml to that. A structured type declares its rules there; a type that takes no rules says why in one line. Table 9 lists the types.

Table 9: Markdown declarations by document type

| Type | Markdown location | Declaration |
| --- | --- | --- |
| vision | docs/VISION.md | fields; success criteria as a keyed list |
| architecture | docs/ARCHITECTURE.md | interface, component, and decision nodes; two tables |
| specifications | docs/SPECIFICATIONS.md | index tables with module-specific columns |
| roadmap | docs/road-map.md | release nodes with a status attribute |
| pattern_language | docs/pattern-language.md | pattern nodes requiring intent, problem, solution |
| engineering_guideline | docs/engineering/engNN-short-name.md | titled section nodes |
| srd | docs/specs/software-requirements/srdNNN-short-name.md | requirement groups, items, criteria |
| use_case | docs/specs/use-cases/relNN.N-ucNNN-short-name.md | success criterion nodes; flow, touchpoint, and scenario fields |
| test_suite | docs/specs/test-suites/test-relNN.N-short-name.md | test case nodes with visible traces |
| audit_register | docs/specs/audits/name-register.md | required record fields |
| audit_report | docs/specs/audits/name.md | none: native Markdown already |
| semantic_model | docs/specs/semantic-models/name.md | purpose and list fields |
| config_format | docs/specs/config-formats/name.md | purpose and list fields |
| constitution | docs/constitutions/name.md | article nodes, the normative blocks |
| migration | docs/migrations/name.md | lifecycle frontmatter; summary and evidence fields |
| guide | docs/guides/name.md | none: native Markdown already |
| stats_baseline | docs/stats/name.json | none: generated JSON |
| corpus_charter | docs/corpus-charter.yaml | none: stays YAML, an executable check plan |

A file selects its declaration by its whole path, matched at a path boundary, because several types share the name placeholder and differ only by directory. A run of N is a digit placeholder only after a lowercase letter or a dot, so VISION stays a word; the placement gate reads locations by the same rule. Audit registers take a -register suffix so they do not collide with the audit reports beside them.

Two types stay outside the format by design. The corpus charter is configuration the audit executes (spec.Charter), not a document a person reads, and the stats baseline is generated output. The design constitution converts with the other constitutions, but its machine-read sections, document_types among them, stay fenced yaml blocks so the loaders keep reading them; GH-2602 moves the loader into agent-core, and GH-2604 makes it read the block from design.md.

<!-- section: corpus-fit -->
## Corpus fit

We checked each new declaration against the corpus before any converter exists, comparing every YAML value with the shape its field declares. Table 10 gives the result; the SRD row is the prototype's measurement on printed Markdown, the others are estimates that GH-2605's round-trip tests replace.

Table 10: Corpus fit of the declarations

| Type | Files | Fit as declared | Migration work |
| --- | --- | --- | --- |
| vision | 6 | 6 | none |
| architecture | 6 | 6 | none |
| specifications | 6 | 6 | none |
| roadmap | 6 | 6 | none |
| pattern_language | 1 | 1 | none |
| engineering_guideline | 9 | 8 | eng02: 2 prose-containment hits |
| srd | 106 | 101 | 15 prose-containment hits in 5 files |
| use_case | 110 | 110 | conversion rules below |
| test_suite | 82 | 82 | none |
| audit_register | 2 | 2 | rename with the -register suffix |
| semantic_model | 11 | 10 | one non_goals block becomes a list |
| config_format | 9 | 5 | 61 prose-containment hits in 4 files; one implementation value becomes a list |
| constitution | 11 | 10 | go-style: 2 prose-containment hits |
| migration | 4 | 4 | none |

Use cases fit once their flow and touchpoints are normalized (Table 6). Flow then prints as an outline, where the 47 steps that carry traces show them as a nested traces item, and touchpoints print as one table with id, target, and reason columns.

<!-- section: rendering -->
## Rendering

GitHub hides HTML comments and renders frontmatter as a table, so a converted document shows its headings, prose, lists, and tables and none of its markers. For GitLab we rely on its documentation, which states that comments are hidden and that frontmatter is shown in a box above the document (GH-2581 decision 4). That is an unverified assumption, and the verification step is named: render one converted document on a GitLab instance. The first GitLab mirror of the repository closes it.

<!-- section: authoring-loop -->
## Authoring loop

An author, human or agent, writes a document, runs the audit, and fixes what it reports. The formatter (mage fmt, GH-2583) applies Tier 0 before anything is reported, so findings only ever concern meaning. Until GH-2583 lands, the prototype's check command runs the same diagnostics on any file:

```text
cd tools/specmd-prototype
go run ./cmd/specmd check ../../docs/engineering/eng10-markdown-spec-format.md
```

Each finding names the rule in design.yaml it enforces, so an author who disagrees with a finding knows which declaration to change, through an issue that amends this guideline.

<!-- field: references -->
## References

- agent-core/docs/constitutions/design.yaml
- docs/engineering/eng09-spec-kit-correspondence.yaml
- docs/engineering/eng03-declaration-standard-library.yaml
- tools/specmd-prototype/FINDINGS.md
- https://www.rfc-editor.org/rfc/rfc2119
- https://docs.gitlab.com/user/markdown/
- https://github.com/Nokia-Bell-Labs/declarative-agents/issues/2578
- https://github.com/Nokia-Bell-Labs/declarative-agents/issues/2581
- https://github.com/Nokia-Bell-Labs/declarative-agents/issues/2597
