---
id: cfg-agent-profile-format
title: Agent Profile Format (profile.yaml)
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: purpose -->
## Purpose

Defines the YAML format of an agent profile, the file `--profile` names. A
profile bundles everything one agent program needs: its machine, its tool
selection, the tool declarations and REST definitions that supply that
vocabulary, and an optional workspace directory. The same file is the unit a
child agent, a REST `machine_request` binding, staging, and the program
identity refer to.

Until this document the fields were recorded only by the AgentProfile type in
internal/tools/catalog/profile.go and the startup narrative in
cfg-runtime-contract. Release 20.0 adds blueprints and their instances
(srd055), which change the format, so the format is written down here.

<!-- field: schema -->
## Schema

<!-- schema_entry: name -->
### name

<!-- field: type -->
#### Type

string

<!-- field: required | form: "scalar" -->
#### Required

false

<!-- field: description -->
#### Description

The agent's name. Required on a blueprint instance (srd055 R1.2).

<!-- schema_entry: machine -->
### machine

<!-- field: type -->
#### Type

string or mapping

<!-- field: required | form: "scalar" -->
#### Required

true

<!-- field: description -->
#### Description

Path to the agent's machine file, which may be a machine template
instance (srd054 R1.2). Inside a blueprint body it can instead be a
mapping carrying `fragment` and `args` that names a machine template
directly (srd055 R2.3). A path to a machine template itself is a load
error (srd054 R1.4).

<!-- schema_entry: tools -->
### tools

<!-- field: type -->
#### Type

[]string

<!-- field: required | form: "scalar" -->
#### Required

true

<!-- field: description -->
#### Description

Tool selection files (cfg-tool-selection-format). At least one entry.

<!-- schema_entry: tool_declarations -->
### tool_declarations

<!-- field: type -->
#### Type

[]string

<!-- field: required | form: "scalar" -->
#### Required

false

<!-- field: description -->
#### Description

Tool declaration files and units (cfg-tool-declaration-format).

<!-- schema_entry: tool_config_dirs -->
### tool_config_dirs

<!-- field: type -->
#### Type

[]string

<!-- field: required | form: "scalar" -->
#### Required

false

<!-- field: description -->
#### Description

Directories whose declaration files are loaded as tool declarations.

<!-- schema_entry: rest_definitions -->
### rest_definitions

<!-- field: type -->
#### Type

[]string

<!-- field: required | form: "scalar" -->
#### Required

false

<!-- field: description -->
#### Description

REST definition files and units (cfg-rest-tool-format).

<!-- schema_entry: rest_config_dirs -->
### rest_config_dirs

<!-- field: type -->
#### Type

[]string

<!-- field: required | form: "scalar" -->
#### Required

false

<!-- field: description -->
#### Description

Directories whose files are loaded as REST definitions.

<!-- schema_entry: directory -->
### directory

<!-- field: type -->
#### Type

string

<!-- field: required | form: "scalar" -->
#### Required

false

<!-- field: description -->
#### Description

Default workspace directory; `--directory` overrides it.

<!-- schema_entry: libraries -->
### libraries

<!-- field: type -->
#### Type

map[string]string

<!-- field: required | form: "scalar" -->
#### Required

false

<!-- field: description -->
#### Description

Named library roots (srd056-unit-libraries R2). Each name is a lowercase
kebab root name other than `agent-core`; each value is a directory,
resolved against this file's directory. An import under `/opt/<name>/`
anywhere in this profile's closure maps to that directory, and
`--library <name>=<path>` overrides it for one run. Roots are scoped to
this profile's closure; a child-agent profile declares its own.

<!-- schema_entry: unit -->
### unit

<!-- field: type -->
#### Type

string

<!-- field: required | form: "scalar" -->
#### Required

false

<!-- field: description -->
#### Description

Unit name of a blueprint, or optionally of a blueprint instance.

<!-- schema_entry: params -->
### params

<!-- field: type -->
#### Type

[]param

<!-- field: required | form: "scalar" -->
#### Required

false

<!-- field: description -->
#### Description

Present only on a blueprint, with the entry format of srd052 R1.1. A
file carrying `params` is a blueprint and is not itself a profile.

<!-- schema_entry: profile -->
### profile

<!-- field: type -->
#### Type

mapping

<!-- field: required | form: "scalar" -->
#### Required

false

<!-- field: description -->
#### Description

A blueprint's body: the fields above, with `$param(name)` references in
scalar values (srd055 R1.1).

<!-- schema_entry: expand -->
### expand

<!-- field: type -->
#### Type

[]expansion

<!-- field: required | form: "scalar" -->
#### Required

false

<!-- field: description -->
#### Description

On a blueprint instance, exactly one entry carrying `fragment` and
`args` whose target has a `profile` body. `as` is not permitted
(srd055 R1.2).

<!-- field: path_resolution -->
## Path resolution

Relative paths resolve against the directory of the file that holds them.
Paths under the agent-core install prefix (/opt/agent-core) map through the
install root that `--core-root` sets. For a blueprint instance, paths in the
blueprint body resolve against the blueprint's directory, and paths in
fields present on the instance resolve against the instance's directory
(srd055 R2.2). Paths under `/opt/<name>` map to the directory the profile
declares for that root under `libraries` (srd056 R2.1), and any other
absolute path is used as written.

<!-- field: blueprints -->
## Blueprints

A blueprint declares an agent once with typed parameters:

  unit: serve-agent-blueprint
  params:
    - {name: agent, type: string}
    - {name: server, type: string}
  profile:
    name: $param(agent)
    machine:
      fragment: /opt/agent-core/tools/machines/serve-machine-template.yaml
      args: {server: $param(server)}
    tools:
      - ../$param(agent)/tools.yaml
    rest_definitions:
      - ../$param(agent)/rest.yaml

and each agent's profile.yaml expands it, carrying a name, the
arguments, and any field that differs:

  name: rag-server
  expand:
    - fragment: ../units/serve-agent-blueprint.yaml
      args: {agent: rag-server, server: rag}
  tool_declarations:
    - declarations.yaml

A field present on the instance replaces the blueprint's field wholesale; a
list is replaced, never merged. The resolved profile is validated as a
hand-written profile is, its program identity is the instance path, and its
digest covers the blueprint and every template it reaches (srd055 R3). A
blueprint named by `--profile`, a child-agent reference, or a
`machine_request` binding is a load error (srd055 R1.3). The fragment path
in the machine mapping above is an install-prefix path, which maps through
the install root as every profile path does (srd055 R2.2).

<!-- field: related_documents -->
## Related documents

- cfg-runtime-contract
- cfg-machine-format
- cfg-tool-selection-format
- cfg-tool-declaration-format
- cfg-rest-tool-format
- srd052-declaration-fragments
- srd054-machine-templates
- srd055-agent-blueprints
- srd056-unit-libraries

<!-- field: implementation -->
## Implementation

- internal/tools/catalog/profile.go
- internal/load/loader.go
