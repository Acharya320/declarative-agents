---
id: rel20.0-uc002-agent-blueprints
title: A profile expands an agent blueprint
status: specified
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A profile author has several agents that are one agent with different words
wired in. They declare the agent once as a blueprint whose profile body names
its machine template and units relative to itself, and replace each agent's
profile with an instance carrying a name, the arguments, and only the fields
that differ. The loader resolves the blueprint, expands the machine
template it names, overlays the instance's fields, and treats the result as
a hand-written profile for validation, identity, dump, staging, and audit.

<!-- field: actor -->
## Actor

A profile author collapsing several agents into one agent with arguments

<!-- field: trigger -->
## Trigger

A profile file lists one `expand` entry whose target is a fragment with
a `profile` body.

<!-- field: flow -->
## Flow

- F1: The loader reads the instance header, rejects an as on the entry and a blueprint that expands a blueprint, and refuses a blueprint named directly where a profile is expected.
- F2: It checks the arguments and substitutes them into the blueprint's text under srd052's rules, then decodes the profile body strictly.
- F3: It resolves the blueprint's paths against the blueprint's directory, expanding a machine template named by a machine mapping, and overlays each field present on the instance, resolved against the instance's directory.
- F4: It validates the resolved profile by the hand-written profile rules, naming instance and blueprint on failure.
- F5: It records the instance as the program identity with the blueprint and templates in its digest, lists the blueprint in the dump, stages the blueprint edges, and lets profileaudit resolve a blueprint-based child.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [R1.1](../software-requirements/srd055-agent-blueprints.md#r11), [R1.2](../software-requirements/srd055-agent-blueprints.md#r12), [R1.3](../software-requirements/srd055-agent-blueprints.md#r13) | blueprint and instance format |
| T2 | [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [R2.1](../software-requirements/srd055-agent-blueprints.md#r21), [R2.2](../software-requirements/srd055-agent-blueprints.md#r22), [R2.3](../software-requirements/srd055-agent-blueprints.md#r23), [R2.4](../software-requirements/srd055-agent-blueprints.md#r24) | resolution |
| T3 | [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [R3.1](../software-requirements/srd055-agent-blueprints.md#r31), [R3.2](../software-requirements/srd055-agent-blueprints.md#r32), [R3.3](../software-requirements/srd055-agent-blueprints.md#r33), [R3.4](../software-requirements/srd055-agent-blueprints.md#r34) | identity, provenance, staging, and audit |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

An instance loads from a blueprint with its fields replacing the blueprint's wholesale; a prefix, nesting, and a blueprint named where a profile is expected each fail naming the referring file.

<!-- field: traces -->
#### Traces

- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [AC1](../software-requirements/srd055-agent-blueprints.md#ac1)

<!-- sc: S2 -->
### S2

Blueprint and instance paths resolve against their own files, a machine mapping expands a template with passed-through arguments, and a resolved profile missing a required field names instance and blueprint.

<!-- field: traces -->
#### Traces

- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [AC2](../software-requirements/srd055-agent-blueprints.md#ac2)

<!-- sc: S3 -->
### S3

The instance is the identity and the blueprint moves its digest, the dump lists the blueprint row, a staged instance loads, and profileaudit resolves a blueprint-based child.

<!-- field: traces -->
#### Traces

- [srd055-agent-blueprints](../software-requirements/srd055-agent-blueprints.md) [AC3](../software-requirements/srd055-agent-blueprints.md#ac3)

<!-- field: out_of_scope -->
## Out of scope

- Merging list fields between blueprint and instance.
- Blueprints expanding blueprints.
- Deriving request machines or REST surfaces from a blueprint.
- Committing resolved profiles.
