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

<!-- field: tool_categories | form: "tree" -->
## Tool categories

- word:
  - meaning:
    A single deterministic verb. A word does one thing and can be composed
    freely by a state machine or, when externally visible, by LLM $tool
    dispatch.
  - examples:
    1. read
    1. write
    1. record_oracle_result
    1. collect_trace_tokens
- stateful_internal:
  - meaning:
    A single internal verb that mutates agent or domain state instead of
    filesystem state. It is hidden from the LLM and invoked by machine
    transitions.
  - examples:
    1. parse_response
    1. record_point_failure
    1. parse_suite_config
- boundary:
  - meaning:
    A tool that intentionally crosses a program, process, model, human,
    or nested-machine boundary. Boundary tools are allowed when the boundary
    is the vocabulary word, but they require especially explicit contracts.
  - examples:
    1. invoke_llm
    1. run_point
    1. self_invoke
    1. rest_await_event
- sentence:
  - meaning:
    A workflow hidden behind one tool name. Sentence tools are migration
    targets. New tools should not use this category except to mark legacy
    compatibility behavior while it is being decomposed.
  - examples:
    1. deprecated compatibility aggregators

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

<!-- field: required_contract_sections | form: "tree" -->
## Required contract sections

- problem:
  One paragraph explaining why the word exists and what vocabulary gap it
  fills. If the gap cannot be stated clearly, do not add the tool.
- goals:
  Numbered, measurable success statements. Goals describe what becomes
  possible when the tool is correct.
- requirements:
  - input: What the tool accepts, including formats, defaults, and preconditions.
  - output: What the tool produces, including schema, signals, and result fields.
  - side_effects: What the tool changes in files, memory, APIs, subprocesses, stderr/stdout, or human/model state.
  - undo: How effects are reversed or compensated, or why they are not.
  - errors: How the tool fails, which signal it emits, and what state remains.
- non_goals:
  Specific things the tool does not do. Non-goals prevent
  one word from being misused as a hidden workflow.
- output_schema:
  Use output.description and output.schema to describe the structured output
  contract. During migration, core.Result.Output may still be a string, but
  the schema is the target contract and should be precise.
- side_effects:
  Use structured side_effects entries. Include kind, target, state, paths
  when relevant, and a description. Use kind: none only when the tool truly
  has no effects beyond output.
- reversibility:
  classification is required for new tools: reversible, compensatable, or
  irreversible. requires_confirmation should be true for irreversible tools
  that affect external systems or user data. A reversible or compensatable
  state-mutating tool encodes an opaque receipt during Execute; a read-only
  tool returns an empty receipt.
- undo:
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
- errors:
  Each expected failure mode should list signal, condition, message_shape,
  and state_after_failure.
- relationships:
  Describe common neighbors. before lists likely preceding words, after lists
  likely consumers, and overlaps explains similar tools and differences.

<!-- field: machine_composition | form: "tree" -->
## Machine composition

- rule:
  machine.yaml owns sentences. If a behavior is a sequence, each step should
  be a tool action with an emitted signal and transition. This makes the
  program auditable, resumable, testable, and rollback-aware.
- allowed_tool_composition:
  A tool may call helpers internally to implement one verb. It must not hide
  program flow. A boundary tool may run a nested machine or child process only
  when that boundary is the declared word, and the contract must state that
  a nested program can mutate state or workspace.
- examples:
  - good:
    1. - sequence: create_point_dir -> copy_dir -> sample_docs -> conditional copy_dir -> git_init -> stage_all -> commit_workspace_baseline
       - reason: Workspace preparation is visible in point.yaml as separate words.
    1. - sequence: run_oracle_check -> record_oracle_result -> collect_trace_tokens -> check_agent_version -> summarize_point_results
       - reason: Correctness, metrics, version checks, and summary are distinct point words.
  - bad:
    1. - tool: prepare_workspace
       - reason: Creates directories, copies files, initializes git, stages, and commits in one command.
    1. - tool: check_results
       - reason: Runs oracle tests, parses trace tokens, checks version, and summarizes in one command.

<!-- field: rollback_model | form: "tree" -->
## Rollback model

- persistence:
  - description: Runtime state persists through the srd035 two-method checkpoint port (Save/Load), backed by the srd036 Dolt adapter by default and NoopCheckpoint when disabled. Each dispatched step is one ordered execution entry, and each entry carries the tool's opaque receipt.
  - tracked_by: checkpoint port + Dolt commit history
- receipt:
  - description: "A state-mutating tool encodes an opaque Receipt during Execute capturing exactly what its Undo needs: prior file content, external resource ids, or a workspace git commit hash the tool manages itself. Read-only tools return an empty receipt."
  - owned_by: the originating tool; the engine and the checkpoint adapters never interpret it
- rollback:
  - db_state: Dolt Revert(run_id, step_index) rewinds recorded state to a target step.
  - external_effects: The checkpoint_rollback lifecycle tool walks the reverted entries in reverse and hands each restored, receipt-bearing Result to the tool's receipt-consuming Undo.
- enforcement:
  - presence: The validator checks receipt presence for reversible state-mutating tools.
  - sufficiency: Receipt sufficiency is each tool's own round-trip test (Execute then Undo restores prior state).
- authoring_rules:
  1. State-mutating tools must encode a receipt during Execute and reverse the effect in a receipt-consuming Undo; read-only tools return an empty receipt and no-op undo.
  1. Filesystem-mutating tools encode prior content or a workspace commit hash the tool manages as their receipt, and undo via file_snapshot_restore or compensating_action.
  1. Boundary tools encode compensation metadata (child run id, resource id, artifact path) in the receipt when the effect can only be compensated, not erased.
  1. No-op undo with an empty receipt is acceptable only for truly read-only/stateless words or temporary migration steps explicitly tracked in an audit.
  1. Irreversible tools require requires_confirmation and declare why no receipt can reverse the effect.
  1. The engine and checkpoint adapters treat receipts as opaque; only the originating tool decodes its receipt.

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

<!-- field: declaration_placement | form: "tree" -->
## Declaration placement

- status: planned
- rule:
  Place a tool declaration in the smallest shared group that owns its
  vocabulary family. Keep behavior-shaping profile config in agent-local
  declarations. Keep tools.yaml files name-only.
- placement_steps:
  1. Classify the tool as builtin or exec.
  1. Classify the tool as shared reusable declaration, agent-local override, or compatibility bundle member.
  1. Choose the group that matches the tool's primary vocabulary family.
  1. Add profile-specific config only in agents/<agent>/builtin.yaml or agents/<agent>/llm/*.yaml.
  1. Add the tool name to the agent tools.yaml selection only when that agent should use it.
  1. Keep workflow order in machine.yaml transitions.
- shared_builtin_groups:
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
- shared_exec_groups:
  - go: tools/exec/go/
  - workspace: tools/exec/workspace/
  - git: tools/exec/git/
  - issues: tools/exec/issues/
- compatibility_bundles:
  1. tools/builtin.yaml
  1. tools/exec.yaml
  1. tools/builtin/all.yaml
  1. tools/exec/all.yaml
- override_rules:
  1. Shared declarations define the stable vocabulary word and reusable contract.
  1. Agent-local overrides define profile-bound config such as child profile, server address, prompt, suite paths, or model.
  1. Agent-local overrides keep the same name and init unless the change intentionally replaces the capability.
  1. Includes load first; local declarations override included declarations with the same tool name.
  1. Compatibility bundles include grouped files and do not become the source of new tool contracts.
- examples:
  1. - tool: read
     - placement: tools/builtin/filesystem/
     - reason: Shared file word with no agent-local config.
  1. - tool: invoke_llm
     - placement: agents/<agent>/llm/*.yaml
     - reason: Model, provider, prompt, and manifest state are profile-bound config.
  1. - tool: invoke_executor
     - placement: agents/planner/builtin.yaml as a self_invoke alias
     - reason: Shared child-agent boundary word with planner-local child profile config.
  1. - tool: build
     - placement: tools/exec/go/
     - reason: Shared Go CLI word.
  1. - tool: issue_create
     - placement: tools/exec/issues/
     - reason: Shared issue tracker CLI boundary word.

<!-- field: acceptance_tests_for_authors -->
## Acceptance tests for authors

- Could a developer implement the tool without asking what it should do?
- Could the agent decide when not to use this tool from problem, goals, and non_goals?
- Can another tool or machine transition consume the output without prose interpretation?
- If the tool fails mid-execution, does the contract say what state remains?
- Can rollback machinery find all state layers the tool changed?
- Can similar tools be distinguished from the relationships.overlaps entries?

<!-- field: example_tools | form: "tree" -->
## Example tools

- readonly_noop_word:
  - name: read
  - category: word
  - problem: The machine needs a deterministic word that returns file contents without changing agent, command, or workspace state.
  - emits:
    1. ToolDone
    1. CommandError
  - side_effects:
    1. - kind: none
       - description: Read-only filesystem access; no state changes beyond Result.Output.
  - reversibility:
    - classification: reversible
    - undo: noop
    - receipt: empty
  - undo:
    - strategy: noop
    - description: There is no effect to reverse because the word only observes state; it returns an empty receipt.
  - why_good:
    1. The verb is singular and deterministic.
    1. "No-op undo is justified by an explicit side_effects: none contract."
- reversible_state_word:
  - name: report_parse_error
  - category: stateful_internal
  - problem: The generator grammar needs a word that turns parse errors into feedback and tracks consecutive retry count.
  - emits:
    1. ToolDone
    1. BudgetExhausted
  - side_effects:
    1. - kind: command_state
       - target: parse_retry_counter
       - description: Increments the retry counter and may emit BudgetExhausted at the configured limit.
  - reversibility:
    - classification: reversible
    - undo: parse_retry_counter_restore
    - receipt: reversible_parse_retry_counter_snapshot
  - undo:
    - strategy: parse_retry_counter_restore
    - description: Decode the retry-counter snapshot from the receipt encoded before the word executed and restore it.
  - why_good:
    1. The retry policy is grammar-visible through ParseFailed -> report_parse_error.
    1. The command-owned state mutation has an exact undo strategy.
- compensatable_boundary:
  - name: create_tracker_issue
  - category: boundary
  - problem: The planner profile needs an exec word that creates an issue through its selected tracker CLI.
  - emits:
    1. ToolDone
    1. ToolFailed
    1. CommandError
  - side_effects:
    1. - kind: external_service_write
       - target: issue_tracker
       - description: Creates an issue outside the agent workspace.
  - reversibility:
    - classification: compensatable
    - undo: close_or_mark_superseded
    - receipt: compensatable_created_issue_id
  - undo:
    - strategy: compensating_action
    - description: Decode the created issue id from the receipt and close it or mark it superseded; external audit history cannot be erased.
  - why_good:
    1. The external profile owns the tracker binary, arguments, and common output contract.
    1. The contract does not pretend the external write is perfectly reversible.
- irreversible_confirmation:
  - name: publish_release
  - category: boundary
  - problem: A release machine may need a word that publishes an immutable artifact to external users.
  - emits:
    1. Published
    1. PublishFailed
    1. CommandError
  - side_effects:
    1. - kind: external_service_write
       - target: release_registry
       - description: Publishes an artifact that downstream users may consume immediately.
  - reversibility:
    - classification: irreversible
    - requires_confirmation: true
    - undo: none
    - receipt: empty
  - undo:
    - strategy: irreversible
    - description: The artifact may be deprecated later, but publication cannot be undone safely; no receipt can reverse it.
  - why_good:
    1. The irreversible effect is explicit before the tool is selected.
    1. Human or policy confirmation is required before affecting external users.
- good_word:
  - name: collect_trace_tokens
  - category: word
  - problem: The point machine needs a word that extracts token counts from trace.ndjson without also judging correctness or writing metadata.
  - emits:
    1. TraceTokensCollected
    1. CommandError
  - side_effects:
    1. - kind: state_mutation
       - target: point_context
       - state: trace_tokens_recorded
       - description: Sets PointContext.Tokens.
  - non_goals:
    1. Does not run oracle tests.
    1. Does not write meta.json.
    1. Does not classify convergence.
- good_boundary:
  - name: run_point
  - category: boundary
  - problem: The evaluator session machine needs a word that enters the point machine for one PointContext.
  - emits:
    1. PointDone
    1. CommandError
  - side_effects:
    1. - kind: nested_machine_execution
       - target: evaluator_point_machine
       - description: Runs agents/critic/point.yaml with point tools.
    1. - kind: state_mutation
       - target: eval_session_state
       - description: Records the completed point result.
  - non_goals:
    1. Does not choose the point sequence; point.yaml does.
- bad_sentence:
  - name: prepare_workspace
  - category: sentence
  - problem: Creates a point directory, copies workspace and docs, initializes git, stages files, and commits.
  - why_bad:
    1. The description contains several verbs.
    1. Intermediate failures leave ambiguous state unless each step is visible.
    1. The sequence belongs in the point machine.
  - replacement:
    1. create_point_dir
    1. copy_dir
    1. sample_docs
    1. conditional copy_dir
    1. git_init
    1. stage_all
    1. commit_workspace_baseline
  - why_replacement_is_better:
    1. Each intermediate filesystem or git side effect has its own emitted signal.
    1. Rollback can restore or compensate the exact word that failed.

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

- cfg-tool-declaration-format
- cfg-tool-selection-format
- cfg-machine-format
- srd001-core-types
- srd005-agentic-loop
- srd021-agent-cli-exec-tool
- srd022-builtin-tools
- srd023-cli-exec-tools
- srd024-blocking-input-tools
- srd025-rollback-lifecycle
- srd035-checkpoint-port
- srd036-dolt-state-persistence

<!-- field: implementation -->
## Implementation

- tools/builtin/
- tools/exec/
- internal/tools/exec/builder.go
- internal/tools/catalog/contracts.go
