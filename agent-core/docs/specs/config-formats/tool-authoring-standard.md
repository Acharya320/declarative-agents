---
id: cfg-tool-authoring-standard
title: Tool Definition Authoring Standard
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: purpose -->
## Purpose

Defines how to author tools as vocabulary words in the agent language.
A program is a declarative state machine wired to tools. The machine
composes behavior; tools provide the smallest useful verbs with typed
input, typed output contracts, declared side effects, explicit failure
semantics, and a rollback story.

<!-- field: principle -->
## Principle

A tool is a word, not a sentence. "read_file", "record_oracle_result", and
"collect_trace_tokens" are words. "prepare a workspace, initialize git,
run tests, parse traces, and write metrics" is a sentence and belongs in
machine.yaml as a sequence of tool actions.

<!-- field: tool_categories -->
## Tool categories

<!-- field: word -->
### Word

<!-- field: meaning -->
#### Meaning

A single deterministic verb. A word does one thing and can be composed
freely by a state machine or, when externally visible, by LLM $tool
dispatch.

<!-- field: examples -->
#### Examples

- read
- write
- record_oracle_result
- collect_trace_tokens

<!-- field: stateful_internal -->
### Stateful internal

<!-- field: meaning -->
#### Meaning

A single internal verb that mutates agent or domain state instead of
filesystem state. It is hidden from the LLM and invoked by machine
transitions.

<!-- field: examples -->
#### Examples

- parse_response
- record_point_failure
- parse_suite_config

<!-- field: boundary -->
### Boundary

<!-- field: meaning -->
#### Meaning

A tool that intentionally crosses a program, process, model, human,
or nested-machine boundary. Boundary tools are allowed when the boundary
is the vocabulary word, but they require especially explicit contracts.

<!-- field: examples -->
#### Examples

- invoke_llm
- run_point
- self_invoke
- rest_await_event

<!-- field: sentence -->
### Sentence

<!-- field: meaning -->
#### Meaning

A workflow hidden behind one tool name. Sentence tools are migration
targets. New tools should not use this category except to mark legacy
compatibility behavior while it is being decomposed.

<!-- field: examples -->
#### Examples

- deprecated compatibility aggregators

<!-- field: good_tool_rules -->
## Good tool rules

| id | rule | test |
| --- | --- | --- |
| GT1 | A tool does one thing. | If the description needs 'and' between observable behaviors, split it. |
| GT2 | A tool is deterministic at its boundary. | The same input and external state produce the same declared result; no hidden LLM judgment lives inside ordinary words. |
| GT3 | A tool has clear structured input and output contracts. | A machine or later tool can know what was produced without interpreting prose. |
| GT4 | A tool declares every side effect beyond its return value. | Filesystem writes, state mutations, subprocesses, external APIs, stderr output, and human/model calls are listed. |
| GT5 | A tool states its rollback story. | The declaration classifies effects as reversible, compensatable, irreversible, or no-op. |
| GT6 | A tool fails loudly and specifically. | Failures use declared signals and explain the state left behind. |
| GT7 | A tool name describes the verb, not the implementation. | Prefer parse_suite_config over run_yaml_loader; prefer collect_trace_tokens over scan_file. |

<!-- field: anti_patterns -->
## Anti patterns

| id | name | description | fix |
| --- | --- | --- | --- |
| AP1 | workflow_tool | A tool that performs a multi-step workflow that could be expressed as machine transitions. | Split the workflow into words and wire them in machine.yaml. |
| AP2 | flag_trenchcoat | One tool whose flags choose unrelated behaviors. | Declare separate tools for each verb. |
| AP3 | hidden_judgment | A tool whose output requires an LLM or human to decide if it succeeded. | Make the tool compute and emit the decision itself, or mark it as a boundary tool. |
| AP4 | undeclared_side_effect | A tool writes files, mutates state, calls APIs, or emits warnings without declaring it. | Add structured side_effects entries and undo/reversibility metadata. |
| AP5 | implementation_name | A tool named for how it is implemented rather than what it means. | Rename around the vocabulary verb. |

<!-- field: required_contract_sections -->
## Required contract sections

<!-- field: problem -->
### Problem

One paragraph explaining why the word exists and what vocabulary gap it
fills. If the gap cannot be stated clearly, do not add the tool.

<!-- field: goals -->
### Goals

Numbered, measurable success statements. Goals describe what becomes
possible when the tool is correct.

<!-- field: requirements -->
### Requirements

<!-- field: input -->
#### Input

What the tool accepts, including formats, defaults, and preconditions.

<!-- field: output -->
#### Output

What the tool produces, including schema, signals, and result fields.

<!-- field: side_effects -->
#### Side effects

What the tool changes in files, memory, APIs, subprocesses, stderr/stdout, or human/model state.

<!-- field: undo -->
#### Undo

How effects are reversed or compensated, or why they are not.

<!-- field: errors -->
#### Errors

How the tool fails, which signal it emits, and what state remains.

<!-- field: non_goals -->
### Non goals

Specific things the tool does not do. Non-goals prevent
one word from being misused as a hidden workflow.

<!-- field: output_schema -->
### Output schema

Use output.description and output.schema to describe the structured output
contract. During migration, core.Result.Output may still be a string, but
the schema is the target contract and should be precise.

<!-- field: side_effects -->
### Side effects

Use structured side_effects entries. Include kind, target, state, paths
when relevant, and a description. Use kind: none only when the tool truly
has no effects beyond output.

<!-- field: reversibility -->
### Reversibility

classification is required for new tools: reversible, compensatable, or
irreversible. requires_confirmation should be true for irreversible tools
that affect external systems or user data. A reversible or compensatable
state-mutating tool encodes an opaque receipt during Execute; a read-only
tool returns an empty receipt.

<!-- field: undo -->
### Undo

Name the exact strategy. Common strategies are noop,
conversation_truncate, session_state_restore, point_context_restore,
file_snapshot_restore, child_command_undo, or compensating_action. Undo is
receipt-consuming: it decodes the receipt the tool encoded during Execute
(prior content, resource ids, or a workspace commit hash) from the restored
Result and reverses the effect. The engine never interprets the receipt.
Builtin strategy names describe the typed Undo implementation registered by
that word. Exec declarations use the closed generic-exec set: noop,
workspace_restore, compensating_action, or irreversible. Startup rejects an
exec strategy the shared interpreter cannot honor.

<!-- field: errors -->
### Errors

Each expected failure mode should list signal, condition, message_shape,
and state_after_failure.

<!-- field: relationships -->
### Relationships

Describe common neighbors. before lists likely preceding words, after lists
likely consumers, and overlaps explains similar tools and differences.

<!-- field: machine_composition -->
## Machine composition

<!-- field: rule -->
### Rule

machine.yaml owns sentences. If a behavior is a sequence, each step should
be a tool action with an emitted signal and transition. This makes the
program auditable, resumable, testable, and rollback-aware.

<!-- field: allowed_tool_composition -->
### Allowed tool composition

A tool may call helpers internally to implement one verb. It must not hide
program flow. A boundary tool may run a nested machine or child process only
when that boundary is the declared word, and the contract must state that
a nested program can mutate state or workspace.

<!-- field: examples -->
### Examples

<!-- field: good -->
#### Good

| sequence | reason |
| --- | --- |
| create_point_dir -> copy_dir -> sample_docs -> conditional copy_dir -> git_init -> stage_all -> commit_workspace_baseline | Workspace preparation is visible in point.yaml as separate words. |
| run_oracle_check -> record_oracle_result -> collect_trace_tokens -> check_agent_version -> summarize_point_results | Correctness, metrics, version checks, and summary are distinct point words. |

<!-- field: bad -->
#### Bad

| tool | reason |
| --- | --- |
| prepare_workspace | Creates directories, copies files, initializes git, stages, and commits in one command. |
| check_results | Runs oracle tests, parses trace tokens, checks version, and summarizes in one command. |

<!-- field: rollback_model -->
## Rollback model

<!-- field: persistence -->
### Persistence

<!-- field: description -->
#### Description

Runtime state persists through the srd035 two-method checkpoint port (Save/Load), backed by the srd036 Dolt adapter by default and NoopCheckpoint when disabled. Each dispatched step is one ordered execution entry, and each entry carries the tool's opaque receipt.

<!-- field: tracked_by -->
#### Tracked by

checkpoint port + Dolt commit history

<!-- field: receipt -->
### Receipt

<!-- field: description -->
#### Description

A state-mutating tool encodes an opaque Receipt during Execute capturing exactly what its Undo needs: prior file content, external resource ids, or a workspace git commit hash the tool manages itself. Read-only tools return an empty receipt.

<!-- field: owned_by -->
#### Owned by

the originating tool; the engine and the checkpoint adapters never interpret it

<!-- field: rollback -->
### Rollback

<!-- field: db_state -->
#### Db state

Dolt Revert(run_id, step_index) rewinds recorded state to a target step.

<!-- field: external_effects -->
#### External effects

The checkpoint_rollback lifecycle tool walks the reverted entries in reverse and hands each restored, receipt-bearing Result to the tool's receipt-consuming Undo.

<!-- field: enforcement -->
### Enforcement

<!-- field: presence -->
#### Presence

The validator checks receipt presence for reversible state-mutating tools.

<!-- field: sufficiency -->
#### Sufficiency

Receipt sufficiency is each tool's own round-trip test (Execute then Undo restores prior state).

<!-- field: authoring_rules -->
### Authoring rules

- State-mutating tools must encode a receipt during Execute and reverse the effect in a receipt-consuming Undo; read-only tools return an empty receipt and no-op undo.
- Filesystem-mutating tools encode prior content or a workspace commit hash the tool manages as their receipt, and undo via file_snapshot_restore or compensating_action.
- Boundary tools encode compensation metadata (child run id, resource id, artifact path) in the receipt when the effect can only be compensated, not erased.
- No-op undo with an empty receipt is acceptable only for truly read-only/stateless words or temporary migration steps explicitly tracked in an audit.
- Irreversible tools require requires_confirmation and declare why no receipt can reverse the effect.
- The engine and checkpoint adapters treat receipts as opaque; only the originating tool decodes its receipt.

<!-- field: authoring_workflow | form: "typed table" -->
## Authoring workflow

| step | name | instruction |
| --- | --- | --- |
| 1 | choose_category | Classify the tool as word, stateful_internal, boundary, or legacy sentence. |
| 2 | write_problem_and_non_goals | Explain why the word exists and what it intentionally does not do. |
| 3 | define_io | Declare parameters and output.schema before implementation. |
| 4 | declare_signals | List emits and make sure named machine actions have transitions for those signals. |
| 5 | declare_effects_and_undo | List side effects and reversibility, and define the receipt the tool encodes on Execute plus the receipt-consuming Undo strategy that reverses the effect. |
| 6 | document_relationships | List before/after neighbors and overlapping tools. |
| 7 | add_tests | Test success, domain failure, command error, side effects, and important undo/reversibility behavior. |

<!-- field: declaration_placement -->
## Declaration placement

<!-- field: status -->
### Status

planned

<!-- field: rule -->
### Rule

Place a tool declaration in the smallest shared group that owns its
vocabulary family. Keep behavior-shaping profile config in agent-local
declarations. Keep tools.yaml files name-only.

<!-- field: placement_steps -->
### Placement steps

- Classify the tool as builtin or exec.
- Classify the tool as shared reusable declaration, agent-local override, or compatibility bundle member.
- Choose the group that matches the tool's primary vocabulary family.
- Add profile-specific config only in agents/<agent>/builtin.yaml or agents/<agent>/llm/*.yaml.
- Add the tool name to the agent tools.yaml selection only when that agent should use it.
- Keep workflow order in machine.yaml transitions.

<!-- field: shared_builtin_groups | form: "tree" -->
### Shared builtin groups

- filesystem: tools/builtin/filesystem/
- llm: tools/builtin/llm/
- lifecycle: tools/builtin/lifecycle/
- control: tools/builtin/control/
- planner: tools/builtin/planner/
- evaluator_session: tools/builtin/critic-session/
- evaluator_point: tools/builtin/critic-point/
- bench: tools/builtin/bench/
- spec_validation: tools/builtin/spec-validation/
- objectstore: tools/builtin/objectstore/

<!-- field: shared_exec_groups | form: "tree" -->
### Shared exec groups

- go: tools/exec/go/
- workspace: tools/exec/workspace/
- git: tools/exec/git/
- issues: tools/exec/issues/

<!-- field: compatibility_bundles -->
### Compatibility bundles

- tools/builtin.yaml
- tools/exec.yaml
- tools/builtin/all.yaml
- tools/exec/all.yaml

<!-- field: override_rules -->
### Override rules

- Shared declarations define the stable vocabulary word and reusable contract.
- Agent-local overrides define profile-bound config such as child profile, server address, prompt, suite paths, or model.
- Agent-local overrides keep the same name and init unless the change intentionally replaces the capability.
- Includes load first; local declarations override included declarations with the same tool name.
- Compatibility bundles include grouped files and do not become the source of new tool contracts.

<!-- field: examples -->
### Examples

| tool | placement | reason |
| --- | --- | --- |
| read | tools/builtin/filesystem/ | Shared file word with no agent-local config. |
| invoke_llm | agents/<agent>/llm/*.yaml | Model, provider, prompt, and manifest state are profile-bound config. |
| invoke_executor | agents/planner/builtin.yaml as a self_invoke alias | Shared child-agent boundary word with planner-local child profile config. |
| build | tools/exec/go/ | Shared Go CLI word. |
| issue_create | tools/exec/issues/ | Shared issue tracker CLI boundary word. |

<!-- field: acceptance_tests_for_authors -->
## Acceptance tests for authors

- Could a developer implement the tool without asking what it should do?
- Could the agent decide when not to use this tool from problem, goals, and non_goals?
- Can another tool or machine transition consume the output without prose interpretation?
- If the tool fails mid-execution, does the contract say what state remains?
- Can rollback machinery find all state layers the tool changed?
- Can similar tools be distinguished from the relationships.overlaps entries?

<!-- field: example_tools -->
## Example tools

<!-- field: readonly_noop_word -->
### Readonly noop word

<!-- field: name -->
#### Name

read

<!-- field: category -->
#### Category

word

<!-- field: problem -->
#### Problem

The machine needs a deterministic word that returns file contents without changing agent, command, or workspace state.

<!-- field: emits -->
#### Emits

- ToolDone
- CommandError

<!-- field: side_effects -->
#### Side effects

| kind | description |
| --- | --- |
| none | Read-only filesystem access; no state changes beyond Result.Output. |

<!-- field: reversibility | form: "tree" -->
#### Reversibility

- classification: reversible
- undo: noop
- receipt: empty

<!-- field: undo -->
#### Undo

<!-- field: strategy -->
##### Strategy

noop

<!-- field: description -->
##### Description

There is no effect to reverse because the word only observes state; it returns an empty receipt.

<!-- field: why_good -->
#### Why good

- The verb is singular and deterministic.
- No-op undo is justified by an explicit side_effects: none contract.

<!-- field: reversible_state_word -->
### Reversible state word

<!-- field: name -->
#### Name

report_parse_error

<!-- field: category -->
#### Category

stateful_internal

<!-- field: problem -->
#### Problem

The generator grammar needs a word that turns parse errors into feedback and tracks consecutive retry count.

<!-- field: emits -->
#### Emits

- ToolDone
- BudgetExhausted

<!-- field: side_effects -->
#### Side effects

| kind | target | description |
| --- | --- | --- |
| command_state | parse_retry_counter | Increments the retry counter and may emit BudgetExhausted at the configured limit. |

<!-- field: reversibility | form: "tree" -->
#### Reversibility

- classification: reversible
- undo: parse_retry_counter_restore
- receipt: reversible_parse_retry_counter_snapshot

<!-- field: undo -->
#### Undo

<!-- field: strategy -->
##### Strategy

parse_retry_counter_restore

<!-- field: description -->
##### Description

Decode the retry-counter snapshot from the receipt encoded before the word executed and restore it.

<!-- field: why_good -->
#### Why good

- The retry policy is grammar-visible through ParseFailed -> report_parse_error.
- The command-owned state mutation has an exact undo strategy.

<!-- field: compensatable_boundary -->
### Compensatable boundary

<!-- field: name -->
#### Name

create_tracker_issue

<!-- field: category -->
#### Category

boundary

<!-- field: problem -->
#### Problem

The planner profile needs an exec word that creates an issue through its selected tracker CLI.

<!-- field: emits -->
#### Emits

- ToolDone
- ToolFailed
- CommandError

<!-- field: side_effects -->
#### Side effects

| kind | target | description |
| --- | --- | --- |
| external_service_write | issue_tracker | Creates an issue outside the agent workspace. |

<!-- field: reversibility | form: "tree" -->
#### Reversibility

- classification: compensatable
- undo: close_or_mark_superseded
- receipt: compensatable_created_issue_id

<!-- field: undo -->
#### Undo

<!-- field: strategy -->
##### Strategy

compensating_action

<!-- field: description -->
##### Description

Decode the created issue id from the receipt and close it or mark it superseded; external audit history cannot be erased.

<!-- field: why_good -->
#### Why good

- The external profile owns the tracker binary, arguments, and common output contract.
- The contract does not pretend the external write is perfectly reversible.

<!-- field: irreversible_confirmation -->
### Irreversible confirmation

<!-- field: name -->
#### Name

publish_release

<!-- field: category -->
#### Category

boundary

<!-- field: problem -->
#### Problem

A release machine may need a word that publishes an immutable artifact to external users.

<!-- field: emits -->
#### Emits

- Published
- PublishFailed
- CommandError

<!-- field: side_effects -->
#### Side effects

| kind | target | description |
| --- | --- | --- |
| external_service_write | release_registry | Publishes an artifact that downstream users may consume immediately. |

<!-- field: reversibility | form: "tree" -->
#### Reversibility

- classification: irreversible
- requires_confirmation: true
- undo: none
- receipt: empty

<!-- field: undo -->
#### Undo

<!-- field: strategy -->
##### Strategy

irreversible

<!-- field: description -->
##### Description

The artifact may be deprecated later, but publication cannot be undone safely; no receipt can reverse it.

<!-- field: why_good -->
#### Why good

- The irreversible effect is explicit before the tool is selected.
- Human or policy confirmation is required before affecting external users.

<!-- field: good_word -->
### Good word

<!-- field: name -->
#### Name

collect_trace_tokens

<!-- field: category -->
#### Category

word

<!-- field: problem -->
#### Problem

The point machine needs a word that extracts token counts from trace.ndjson without also judging correctness or writing metadata.

<!-- field: emits -->
#### Emits

- TraceTokensCollected
- CommandError

<!-- field: side_effects -->
#### Side effects

| kind | target | state | description |
| --- | --- | --- | --- |
| state_mutation | point_context | trace_tokens_recorded | Sets PointContext.Tokens. |

<!-- field: non_goals -->
#### Non goals

- Does not run oracle tests.
- Does not write meta.json.
- Does not classify convergence.

<!-- field: good_boundary -->
### Good boundary

<!-- field: name -->
#### Name

run_point

<!-- field: category -->
#### Category

boundary

<!-- field: problem -->
#### Problem

The evaluator session machine needs a word that enters the point machine for one PointContext.

<!-- field: emits -->
#### Emits

- PointDone
- CommandError

<!-- field: side_effects -->
#### Side effects

| kind | target | description |
| --- | --- | --- |
| nested_machine_execution | evaluator_point_machine | Runs agents/critic/point.yaml with point tools. |
| state_mutation | eval_session_state | Records the completed point result. |

<!-- field: non_goals -->
#### Non goals

- Does not choose the point sequence; point.yaml does.

<!-- field: bad_sentence -->
### Bad sentence

<!-- field: name -->
#### Name

prepare_workspace

<!-- field: category -->
#### Category

sentence

<!-- field: problem -->
#### Problem

Creates a point directory, copies workspace and docs, initializes git, stages files, and commits.

<!-- field: why_bad -->
#### Why bad

- The description contains several verbs.
- Intermediate failures leave ambiguous state unless each step is visible.
- The sequence belongs in the point machine.

<!-- field: replacement -->
#### Replacement

- create_point_dir
- copy_dir
- sample_docs
- conditional copy_dir
- git_init
- stage_all
- commit_workspace_baseline

<!-- field: why_replacement_is_better -->
#### Why replacement is better

- Each intermediate filesystem or git side effect has its own emitted signal.
- Rollback can restore or compensate the exact word that failed.

<!-- field: migration_guidance -->
## Migration guidance

- When a legacy tool is a sentence, first classify it as sentence or deprecated boundary.
- Split the highest-risk side effects first: filesystem writes, child processes, and domain state mutation.
- Move sequencing into machine.yaml and update emits on each new word.
- Keep compatibility shims only when a shipped machine still selects the old tool; mark them internal and deprecated.
- Update semantic-model audits and SRDs in the same change as the declaration or machine change.

<!-- field: historical_background -->
## Historical background

- docs/guides/lifecycle-rollback.md
- docs/guides/tools-as-dsl.md

<!-- field: related_documents -->
## Related documents

- [cfg-tool-declaration-format](tool-declaration-format.md)
- [cfg-tool-selection-format](tool-selection-format.md)
- [cfg-machine-format](machine-format.md)
- [srd001-core-types](../software-requirements/srd001-core-types.md)
- [srd005-agentic-loop](../software-requirements/srd005-agentic-loop.md)
- [srd021-agent-cli-exec-tool](../software-requirements/srd021-agent-cli-exec-tool.md)
- [srd022-builtin-tools](../software-requirements/srd022-builtin-tools.md)
- [srd023-cli-exec-tools](../software-requirements/srd023-cli-exec-tools.md)
- [srd024-blocking-input-tools](../software-requirements/srd024-blocking-input-tools.md)
- [srd025-rollback-lifecycle](../software-requirements/srd025-rollback-lifecycle.md)
- [srd035-checkpoint-port](../software-requirements/srd035-checkpoint-port.md)
- [srd036-dolt-state-persistence](../software-requirements/srd036-dolt-state-persistence.md)

<!-- field: implementation -->
## Implementation

- tools/builtin/
- tools/exec/
- internal/tools/exec/builder.go
- internal/tools/catalog/contracts.go
