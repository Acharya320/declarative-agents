---
id: srd055-agent-blueprints
title: Agent Blueprints
depends_on:
  - srd_id: srd052-declaration-fragments
    symbols_used: [fragment format, params, expand, $param substitution, argument checks, expansions dump section]
  - srd_id: srd054-machine-templates
    symbols_used: [machine template, machine instance, template body kind]
  - srd_id: srd050-declaration-imports
    symbols_used: [closure files, relative path rule]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

An agent is a profile: a machine, a tool selection, tool declarations, REST
definitions, and a few directories. Several agents in the repository are the
same agent with different words wired in. The four chatbot-mesh serve
agents share one ten-state machine and one profile shape and differ only in
their request machines, their REST surface, and the words those reference.
Machine templates (srd054) remove the copied machine; the copied profile
and the relationship between it and its machine remain.

A blueprint is the profile-level template. It declares an agent once with
typed parameters, names its machine template and units by paths relative to
itself, and is expanded by a profile file that carries a name, the
arguments, and any fields that genuinely differ. Several agents become one
agent with arguments, and the profile file that `--profile` names stays the
unit of deployment.

<!-- field: goals -->
## Goals

- G1: Let a whole agent profile be declared once as a fragment with typed parameters and expanded per agent.
- G2: Let a blueprint carry its machine template and units with it, so an instance does not need to know where they live.
- G3: Keep the instance the program's identity, while every blueprint and template it reaches is part of that identity's digest.
- G4: Resolve, dump, stage, and audit a blueprint instance as a hand-written profile is.

<!-- req: R1 -->
## R1 - Blueprint and instance format

<!-- req: R1.1 -->
### R1.1

A blueprint SHALL be a fragment file (srd052 R1.1) whose body kind is `profile`. The `profile` body SHALL hold AgentProfile fields: `name`, `machine`, `tools`, `tool_declarations`, `tool_config_dirs`, `rest_definitions`, `rest_config_dirs`, and `directory`, plus any field a later profile format adds.

<!-- req: R1.2 -->
### R1.2

A profile instance SHALL carry `name` and exactly one `expand` entry whose target is a blueprint, and MAY carry any AgentProfile field. A field present on the instance SHALL replace the blueprint's field wholesale; a list is replaced, never merged. The entry SHALL NOT carry `as`. An instance MAY carry `unit`.

<!-- req: R1.3 -->
### R1.3

A blueprint SHALL NOT expand another blueprint. A blueprint named by `--profile`, by a child-agent profile reference, or by a REST `machine_request` profile binding SHALL be a load error naming the referring file and the blueprint, because a blueprint holds unfilled references and is not a profile until its arguments arrive.

<!-- req: R2 -->
## R2 - Resolution

<!-- req: R2.1 -->
### R2.1

Argument checking, `$param(name)` substitution, and leftover detection SHALL be srd052 R2.2 and R2.3 unchanged, applied to the blueprint's text before the `profile` body is decoded strictly.

<!-- req: R2.2 -->
### R2.2

Relative paths in the blueprint body SHALL resolve against the blueprint file's directory. Paths in fields present on the instance SHALL resolve against the instance file's directory. Install-prefix paths SHALL map through the agent-core install root in both, as profile paths do today, and under any library root srd056 declares.

<!-- req: R2.3 -->
### R2.3

A blueprint body's `machine` MAY be a path to a machine file or a machine template instance (srd054 R1.2), or a mapping carrying `fragment` and `args` that names a machine template directly. The mapping SHALL be expanded as if an instance file with those arguments stood beside the blueprint, and its argument values MAY be `$param` references to the blueprint's own parameters.

<!-- req: R2.4 -->
### R2.4

The resolved profile SHALL satisfy the rules a hand-written profile does, including a required `machine` and at least one `tools` entry, and a violation SHALL name both the instance and the blueprint.

<!-- req: R3 -->
## R3 - Identity, provenance, staging, and audit

<!-- req: R3.1 -->
### R3.1

The program identity's `profile` SHALL be the instance path. The closure's files, and therefore the identity digest, SHALL include the blueprint and every machine template and stage fragment the resolved profile reaches.

<!-- req: R3.2 -->
### R3.2

The canonical dump's `profile` section SHALL render the resolved instance, and `expansions` SHALL list the blueprint with `kind: profile`, the blueprint path, the arguments, and the instance name, beside any machine, stage, tool, and REST rows.

<!-- req: R3.3 -->
### R3.3

Staging a declaration tree (`pkg/profilestage`) SHALL follow the edges from an instance to its blueprint and from the blueprint to everything it names, so a staged instance loads with no source checkout present.

<!-- req: R3.4 -->
### R3.4

Profile inspection (`pkg/profileaudit`) SHALL resolve a blueprint-based child profile as it resolves a hand-written one, so reachability and wiring checks see the resolved child.

<!-- field: non_goals -->
## Non-goals

- Merging list fields between blueprint and instance. An instance field replaces the blueprint's; an instance that needs a merged list states the whole list.
- Blueprints expanding blueprints.
- Deriving request machines or REST surfaces from a blueprint. Those stay per-agent until a measured duplication calls for them.
- Machines invoked as tools by other machines, and loader words.
- Committing any resolved profile.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

A blueprint with params and a profile body loads through an instance carrying name and one expand entry; an instance field replaces the blueprint field wholesale; an as on the entry, a nested blueprint, and a blueprint named by --profile, a child-agent reference, or a machine_request binding each fail naming the referring file and the blueprint.

<!-- field: traces -->
#### Traces

- [R1.1](#r11)
- [R1.2](#r12)
- [R1.3](#r13)

<!-- ac: AC2 -->
### AC2

Argument faults fail as they do for srd052 fragments; blueprint paths resolve against the blueprint and instance paths against the instance; a machine mapping expands a template directly with passed-through arguments; and a resolved profile missing a required field fails naming instance and blueprint.

<!-- field: traces -->
#### Traces

- [R2.1](#r21)
- [R2.2](#r22)
- [R2.3](#r23)
- [R2.4](#r24)

<!-- ac: AC3 -->
### AC3

The identity profile is the instance and its digest changes when the blueprint or a reached template changes; the dump renders the resolved profile and lists the blueprint row with kind profile; a staged instance loads with no checkout present; and profileaudit resolves a blueprint-based child profile.

<!-- field: traces -->
#### Traces

- [R3.1](#r31)
- [R3.2](#r32)
- [R3.3](#r33)
- [R3.4](#r34)

<!-- field: related_documents -->
## Related documents

- cfg-agent-profile-format
- cfg-runtime-contract
- srd052-declaration-fragments
- srd054-machine-templates

<!-- field: implementation_status -->
## Implementation status

<!-- field: planned -->
### Planned

- GH-2123 implements blueprints in the profile loader, the closure, the dump, profilestage, and profileaudit, once the corpus holds a profile group that is one agent with arguments. GH-2115 measured the four chatbot-mesh serve agents and found their duplication was the serve machine, which srd054's template removed; their request machines differ, so they were not that group.
