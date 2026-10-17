---
id: rel20.0-uc001-machine-templates
title: A machine file expands a whole-machine template
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A profile author finds the same machine written in several agents with only
the words wired into it changed. They move the machine into a template with
typed parameters for those words, and replace each copy with a machine file
that carries a unit name, one expansion, and optionally its own name and
purpose. The loader substitutes the arguments, splices any stages the
template expands, validates the result as one hand-written machine, and
records the template in the closure, the program digest, and the dump.

<!-- field: actor -->
## Actor

A profile author replacing copied machines with one template

<!-- field: trigger -->
## Trigger

A machine file lists one `expand` entry whose target is a fragment with
a `machine` body.

<!-- field: flow -->
## Flow

- F1: The loader reads the instance header and rejects any top-level field other than unit, expand, name, and purpose, an as on the entry, and a second entry.
- F2: It reads the target, selects replace or splice from its body kind, and rejects a template that expands another template.
- F3: It checks the arguments and substitutes them into the template's text under srd052's rules, then splices any stage fragments the template body expands, resolving their paths against the template's directory.
- F4: It applies the instance's name and purpose and validates the result by the hand-written machine path, naming both template and instance in any diagnostic.
- F5: It visits the instance, the template, and each stage file, so they enter the closure, the program digest, the dump's expansions with their kinds, and a staged tree.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [R1.1](../software-requirements/srd054-machine-templates.md#r11), [R1.2](../software-requirements/srd054-machine-templates.md#r12), [R1.3](../software-requirements/srd054-machine-templates.md#r13), [R1.4](../software-requirements/srd054-machine-templates.md#r14) | template and instance format |
| T2 | [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [R2.1](../software-requirements/srd054-machine-templates.md#r21), [R2.2](../software-requirements/srd054-machine-templates.md#r22), [R2.3](../software-requirements/srd054-machine-templates.md#r23), [R2.4](../software-requirements/srd054-machine-templates.md#r24), [R2.5](../software-requirements/srd054-machine-templates.md#r25) | expansion and validation |
| T3 | [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [R3.1](../software-requirements/srd054-machine-templates.md#r31), [R3.2](../software-requirements/srd054-machine-templates.md#r32), [R3.3](../software-requirements/srd054-machine-templates.md#r33) | closure, identity, provenance, and staging |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

An instance loads from a template; every forbidden instance field, a prefix, nesting, and a template named where a machine is expected each fail naming the file.

<!-- field: traces -->
#### Traces

- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [AC1](../software-requirements/srd054-machine-templates.md#ac1)

<!-- sc: S2 -->
### S2

The expanded machine equals the equivalent hand-written machine apart from name and purpose, and a fault inside the template names template and instance.

<!-- field: traces -->
#### Traces

- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [AC2](../software-requirements/srd054-machine-templates.md#ac2)

<!-- sc: S3 -->
### S3

A template change moves the program digest, the dump lists machine and stage expansions with their kinds, and a staged instance loads with no checkout.

<!-- field: traces -->
#### Traces

- [srd054-machine-templates](../software-requirements/srd054-machine-templates.md) [AC3](../software-requirements/srd054-machine-templates.md#ac3)

<!-- field: out_of_scope -->
## Out of scope

- Instance overrides other than name and purpose.
- Templates expanding templates.
- Machines invoked as tools by other machines.
- Committing expanded machines.
