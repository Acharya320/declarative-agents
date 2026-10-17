---
id: issue-format
title: Issue Description Format
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: purpose -->
## Purpose

We use this schema for issue descriptions. Agents parse the
description as YAML before they read, edit, test, close, and commit work.
The format keeps each issue small, traceable, and tied to listed files.

<!-- field: schema | form: "tree" -->
## Schema

- required_fields:
  1. deliverable_type
  1. required_reading
  1. files
  1. requirements
  1. acceptance_criteria
- optional_fields:
  1. format_rule
  1. required_sections
  1. design_decisions
  1. dependencies
  1. notes

<!-- field: field_specs -->
## Field specs

<!-- field: deliverable_type -->
### Deliverable type

<!-- field: type -->
#### Type

string

<!-- field: required | form: "scalar" -->
#### Required

true

<!-- field: values -->
#### Values

- code
- documentation

<!-- field: description -->
#### Description

The kind of deliverable. Use code when the issue touches cmd/,
internal/, pkg/, tools that require Go changes, or tests. Use
documentation when the issue only changes docs, YAML specs, or guides.

<!-- field: required_reading -->
### Required reading

<!-- field: type -->
#### Type

list of strings

<!-- field: required | form: "scalar" -->
#### Required

true

<!-- field: description -->
#### Description

Repo-relative paths the worker must read before changing files. Include
upstream specs, architecture sections, source files, tests, and related
issue context.

<!-- field: files -->
### Files

<!-- field: type -->
#### Type

list

<!-- field: required | form: "scalar" -->
#### Required

true

<!-- field: accepted_item_shapes -->
#### Accepted item shapes

- string
- mapping

<!-- field: mapping_fields | form: "keyed table" -->
#### Mapping fields

| id | type | values | required |
| --- | --- | --- | --- |
| path | string |  | true |
| action | string | [create, modify, delete] | false |
| note | string |  | false |

<!-- field: description -->
#### Description

Repo-relative paths the worker may change. Current agent-core issues
often use a simple list of path strings. More complex issues may use
mappings with path, action, and note.

<!-- field: requirements -->
### Requirements

<!-- field: type -->
#### Type

list

<!-- field: required | form: "scalar" -->
#### Required

true

<!-- field: accepted_item_shapes -->
#### Accepted item shapes

- string
- mapping

<!-- field: mapping_fields | form: "keyed table" -->
#### Mapping fields

| id | type | required | pattern |
| --- | --- | --- | --- |
| id | string | false | R[0-9]+ |
| text | string | false |  |

<!-- field: description -->
#### Description

Specific work to complete. Use numbered mappings for new issues when
practical. String items remain valid for compatibility with the current
queue. An item may carry the literal marker [NEEDS CLARIFICATION: question]
for a decision nobody has made; an issue carrying one is not ready for a
worker, and resolving it means replacing the marker with the decision in
the same item. The inline marker is the tracker-issue convention only;
corpus YAML documents record unmade decisions in their structured
needs_clarification field (design.yaml clarification_marker).

<!-- field: acceptance_criteria -->
### Acceptance criteria

<!-- field: type -->
#### Type

list

<!-- field: required | form: "scalar" -->
#### Required

true

<!-- field: accepted_item_shapes -->
#### Accepted item shapes

- string
- mapping

<!-- field: mapping_fields | form: "keyed table" -->
#### Mapping fields

| id | type | required | pattern |
| --- | --- | --- | --- |
| id | string | false | AC[0-9]+ |
| text | string | false |  |

<!-- field: description -->
#### Description

Checkable outcomes. Each item must be verifiable by inspection, tests,
mage audit, or another named command. An item may carry the literal marker
[NEEDS CLARIFICATION: question]; an issue carrying one is not ready for a
worker, and resolving it means replacing the marker with the decision in
the same item. The inline marker is the tracker-issue convention only;
corpus YAML documents use the needs_clarification field instead
(design.yaml clarification_marker).

<!-- field: format_rule -->
### Format rule

<!-- field: type -->
#### Type

string

<!-- field: required | form: "scalar" -->
#### Required

false

<!-- field: description -->
#### Description

Documentation issues may name the governing document type, such as
srd-format, use-case-format, test-suite-format, or constitution-format.

<!-- field: required_sections -->
### Required sections

<!-- field: type -->
#### Type

list of strings

<!-- field: required | form: "scalar" -->
#### Required

false

<!-- field: description -->
#### Description

Documentation issues may list headings or YAML fields that must appear in
the output artifact.

<!-- field: design_decisions -->
### Design decisions

<!-- field: type -->
#### Type

list

<!-- field: required | form: "scalar" -->
#### Required

false

<!-- field: description -->
#### Description

Decisions made while shaping the issue. Prefer mappings with id and text
when a decision has lasting design value.

<!-- field: yaml_rules -->
## YAML rules

<!-- yaml_rule: Y1 -->
### Y1

<!-- field: rule -->
#### Rule

Use valid YAML.

<!-- field: check -->
#### Check

python3 - <<'PY'
import sys, yaml
for p in sys.argv[1:]: yaml.safe_load(open(p))
PY <files>

<!-- yaml_rule: Y2 -->
### Y2

<!-- field: rule -->
#### Rule

Use block literals for prose longer than one sentence.

<!-- yaml_rule: Y3 -->
### Y3

<!-- field: rule -->
#### Rule

Use ASCII punctuation. Use "--" instead of em dashes.

<!-- yaml_rule: Y4 -->
### Y4

<!-- field: rule -->
#### Rule

Quote strings that contain colons when YAML would parse them as mappings.

<!-- yaml_rule: Y5 -->
### Y5

<!-- field: rule -->
#### Rule

Keep paths repo-relative unless the issue explicitly targets another repository.

<!-- field: examples -->
## Examples

<!-- field: documentation_issue -->
### Documentation issue

```yaml
deliverable_type: documentation
required_reading:
  - docs/VISION.md
  - docs/ARCHITECTURE.md
files:
  - docs/specs/use-cases/rel02.0-uc001-example.yaml
requirements:
  - Add the use case in the existing schema.
  - Link every success criterion to SRD acceptance criteria.
acceptance_criteria:
  - The YAML parses.
  - mage audit reports no new warnings for the file.
```

<!-- field: documentation_issue_with_clarification -->
### Documentation issue with clarification

An unresolved item reads
"requirements: [Retry the fetch N times, where N is [NEEDS CLARIFICATION:
how many retries before giving up?]]". The author who knows the answer
edits the item in place, and the resolved form reads "requirements: [Retry
the fetch 3 times]". No marker survives into a worked issue.

<!-- field: code_issue -->
### Code issue

```yaml
deliverable_type: code
required_reading:
  - docs/specs/software-requirements/srd014-spec-graph.md
  - pkg/spec/validate.go
files:
  - pkg/spec/validate.go
  - pkg/spec/validate_test.go
requirements:
  - id: R1
    text: Add the validator behavior described by the SRD.
acceptance_criteria:
  - id: AC1
    text: go test ./pkg/spec passes.
```

<!-- field: references -->
## References

- docs/constitutions/design.yaml
- docs/constitutions/execution.md
- AGENTS.md
