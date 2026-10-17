---
id: cfg-tool-selection-format
title: Tool Selection Format (tools.yaml)
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: purpose -->
## Purpose

Defines the simple name-list format used by each agent's tools.yaml
to select which declared tools are available to that agent. This
separates tool declaration semantics from per-agent availability, enabling
tool reuse across agents with different toolsets.

<!-- field: schema -->
## Schema

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

List of tool names to make available to the agent. Each name
must match a tool declared by the selected profile's
tool_config_dirs and tool_declarations. Unknown names cause a
startup error.

<!-- field: cli_mapping -->
## CLI mapping

<!-- field: description -->
### Description

Startup uses an agent profile. It points at one or more tools.yaml
selection files, tool_config_dirs containing individual declaration YAML
files, and optional explicit tool_declarations such as agent-local LLM or
override files. The normal CLI exposes only --profile for agent identity.
Tool selection paths and declaration paths belong to the profile.

<!-- field: flags -->
### Flags

<!-- flag: --profile -->
#### --profile

<!-- field: repeatable | form: "scalar" -->
##### Repeatable

false

<!-- field: description -->
##### Description

Preferred agent entry point. The profile supplies tool selection
files, declaration directories, and explicit declaration files.
Relative paths are resolved from the profile file.

<!-- field: removed_flags -->
### Removed flags

cmd/agent no longer accepts --tools, --tools-declaration, or
--tool-config-dir startup flags. Profile YAML owns those paths.

<!-- field: selection_files -->
## Selection files

<!-- selection_file: <profile-root>/generator/tools.yaml -->
### <profile-root>/generator/tools.yaml

<!-- field: tool_count | form: "scalar" -->
#### Tool count

15

<!-- field: description -->
#### Description

Selects file operations (read, write, edit, find, list_files),
build/test tools (build, vet, lint, test), and internal pipeline
tools for the LLM-driven coding loop (invoke_llm, parse_response,
report_parse_error, reset_history, nudge_reread, done). Final
validation order is declared directly in machine.yaml as
build → lint → test.

<!-- selection_file: <profile-root>/evaluator/tools.yaml -->
### <profile-root>/evaluator/tools.yaml

<!-- field: tool_count | form: "scalar" -->
#### Tool count

8

<!-- field: description -->
#### Description

Selects session-level tools for <profile-root>/evaluator/machine.yaml:
parse_suite_config, discover_suite_samples, expand_eval_grid,
init_eval_session, report_suite_summary, materialize_eval_points, run_point,
and report_session.

<!-- selection_file: <profile-root>/evaluator/tools-point.yaml -->
### <profile-root>/evaluator/tools-point.yaml

<!-- field: tool_count | form: "scalar" -->
#### Tool count

17

<!-- field: description -->
#### Description

Selects per-point tools for <profile-root>/evaluator/point.yaml:
create_point_dir, copy_dir, sample_docs,
git_init, stage_all,
commit_workspace_baseline, rev_parse, record_agent_commit,
dump_config, run_agent, run_oracle_check, record_oracle_result,
collect_trace_tokens, check_agent_version,
summarize_point_results, record_point_failure, and collect_metrics. Loaded by run_point
to build the nested point registry.

<!-- selection_file: <profile-root>/planner/tools.yaml -->
### <profile-root>/planner/tools.yaml

<!-- field: tool_count | form: "scalar" -->
#### Tool count

29

<!-- field: description -->
#### Description

Selects planning tools (extract_task, select_all_ready,
seed_passthrough_plan, mark_nodes_planning,
project_planner_context, capture_planner_failure,
reset_failure_context, compose_planner_prompt,
parse_plan, format_issue, write,
create_tracker_issue, record_tracker_issue, mark_nodes_executing,
format_task_file, invoke_executor,
reset_retry_count, increment_retry, check_retry_limit,
mark_task_done, mark_task_failed, remaining_work),
LLM tools (invoke_llm, reset_history),
and direct build/test transition tools (build, vet, test).

<!-- selection_file: <profile-root>/jurist/tools.yaml -->
### <profile-root>/jurist/tools.yaml

<!-- field: tool_count | form: "scalar" -->
#### Tool count

3

<!-- field: description -->
#### Description

Name-only selection of 3 tools (load_corpus, validate_specs,
format_report). Declarations are loaded from the profile's
shared builtin declaration directory.

<!-- selection_file: <profile-root>/bench/tools.yaml -->
### <profile-root>/bench/tools.yaml

<!-- field: tool_count | form: "scalar" -->
#### Tool count

5

<!-- field: description -->
#### Description

Selects profile aliases for generic REST launch/await/stop,
value_predicate, and self_invoke. Bench route and request-machine
declarations remain beside the profile.

<!-- field: grouped_declaration_loading -->
## Grouped declaration loading

<!-- field: status -->
### Status

planned

<!-- field: rule -->
### Rule

Tool selection files stay name-only during and after the declaration
grouping migration. A selection file never names tools/builtin/*,
tools/exec/*, or bundle files. Profiles and CLI flags choose declaration
sources; tools.yaml chooses tool names from those sources.

<!-- field: compatibility -->
### Compatibility

- Existing profiles must keep loading while declarations move into grouped files.
- tools/builtin.yaml and tools/exec.yaml remain compatibility bundles for direct CLI startup.
- Profile-local tool_declarations continue to load after shared groups so local config can override shared identity.
- A selected name must resolve after all shared groups, bundles, and local overrides are merged.

<!-- field: agent_selection_intent -->
### Agent selection intent

Generator profiles load builtin/filesystem, builtin/llm, and exec/go shared
groups, then load <profile-root>/generator/llm/*.yaml for model-boundary config.
The build, lint, and test order remains in machine.yaml.

Planner profiles load builtin/planner, builtin/llm, and exec/go groups.
They add <profile-root>/planner/builtin.yaml for child-agent config and
<profile-root>/planner/llm/*.yaml for model config.

Evaluator profiles load builtin/critic-session and
builtin/critic-point declarations. The agent-local
<profile-root>/evaluator/builtin.yaml file supplies suite and point-machine config,
and run_point loads <profile-root>/evaluator/tools-point.yaml for the nested point
machine.

Bench profiles load builtin/bench declarations and keep evaluator child
profile config in <profile-root>/bench/builtin.yaml. Specification-critic profiles load only
builtin/spec-validation. Lifecycle history and rollback profiles load
builtin/lifecycle; approval profiles add suspend and done through their own
tools.yaml.

<!-- field: related_documents -->
## Related documents

- [cfg-tool-declaration-format](tool-declaration-format.md)
- srd013
- srd022
- srd023

<!-- field: implementation -->
## Implementation

- internal/tools/catalog/loading.go
- internal/tools/catalog/profile.go
