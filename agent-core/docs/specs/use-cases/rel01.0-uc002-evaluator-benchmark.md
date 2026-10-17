---
id: rel01.0-uc002-evaluator-benchmark
title: Evaluator runs generator profiles and collects results
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

The external evaluator profile configures the agent binary to load a suite
YAML with generator profiles, samples, optional grid parameters, and
repetitions, then a session-level state machine walks the
profile/grid/sample/rep combinations. For each point, a nested core.Loop
prepares the workspace, dumps config, runs the generator as a subprocess,
checks results, and collects metrics. Result files are written to a structured
directory tree for the bench UI and CLI reporting tools. Artifacts remain
inspectable after the run.

applications/coding-agent mage integration:criticGate exercises the application
gate with accepted and rejected workspaces. The benchmark path supports:
generator/profile-qwen35b.yaml (qwen3.6:35b-mlx) and
generator/profile-qwen27b.yaml (qwen3.6:27b-mlx). It reuses the same
sample workspace as the generator integration target.

<!-- field: actor -->
## Actor

Developer benchmarking model performance

<!-- field: trigger -->
## Trigger

A developer running the benchmark profile runs:
  agent --profile "$AGENT_CATALOG_ROOT/evaluator/profile.yaml" \
        --request testdata/integration/uc002-evaluator-benchmark/suite.yaml \
        --output eval-results

<!-- field: preconditions -->
## Preconditions

- Ollama is running at localhost:11434
- The qwen3.6:35b-mlx and qwen3.6:27b-mlx artifacts are pulled
- The agent binary is built (bin/agent)
- The sample workspace from uc001 is available

<!-- field: flow -->
## Flow

- F1: The session machine runs parse_suite_config, discover_suite_samples, expand_eval_grid, init_eval_session, and report_suite_summary to prepare suite input and create the session output directory (srd019).
- F2: parse_suite_config reads the suite's profiles list, resolves each profile relative to the suite file, and rejects obsolete suite fields that bypass profile-selected runtime configuration (srd019).
- F3: materialize_eval_points emits every (profile, gridpoint, sample, rep) combination in deterministic order. For the integration suite this means generator-qwen35b, then generator-qwen27b, over the discovered samples (srd019).
- F4: The session machine dispatches run_point sequentially with for_each; each item launches a nested core.Loop with the point machine (evaluator/point.yaml relative to the profile root) (srd005, evaluator/machine.yaml, evaluator/point.yaml).
- F5: The point machine runs create_point_dir, shared copy_dir, sample_docs with a conditional shared copy_dir, shared git_init and stage_all, and the profile-configured commit_workspace_baseline exec alias to prepare an isolated workspace (srd019).
- F6: dump_config serializes the resolved profile path, model, harness label, sample, timeout, and repetition into experiment.yaml for reproducibility (srd019).
- F7: run_agent spawns the generator binary as a subprocess with OTel traceparent propagation and passes --profile for profile-based points (srd016, srd018, srd019).
- F8: The profile-configured run_oracle_check exec word, record_oracle_result adapter, collect_trace_tokens, check_agent_version, and summarize_point_results classify the point outcome and collect trace-derived result state (srd019).
- F9: collect_metrics writes meta.json with point metadata, duration, and outcome classification; early point failures first pass through record_point_failure so the same writer persists terminal failure stage and cause (srd019).
- F10: Each nested loop completes as a separately traced and receipted for_each item.
- F11: When the iterator joins with PointsCompleted, report_session prints the summary.
- F12: Results are written to eval-results/<suite>/<timestamp>/<point-id>/ with meta.json, trace.ndjson, and experiment.yaml (srd019).

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd005-agentic-loop R4.1, R4.2, R4.3, R4.4, R5.1, R5.2, R5.3, R5.4, R5.5, R5.6, R5.7, R5.8, R5.9, R5.10, R5.11, R5.12 | LoopParams and Loop drive both session and point state machines |
| T2 | srd006-yaml-machine-config R1.1, R1.2, R1.3, R1.4, R1.5, R1.6, R1.7, R1.8, R1.9, R2.1, R2.2, R2.3, R2.4, R2.5, R2.6, R2.7, R2.8, R2.9, R2.10, R2.11, R3.1, R3.2, R3.3, R3.4, R3.5, R4.1, R4.2, R4.3, R4.4 | session and point machine.yaml loading and validation |
| T3 | srd016-traceparent R2.1, R2.2, R3.1, R3.2, R3.3 | W3C traceparent formatting and parent-span context propagation to subprocess |
| T4 | srd019-eval-harness R0.1, R0.2, R0.3, R0.4, R0.5, R0.6, R1.1, R1.2, R1.3, R1.4, R1.5, R1.6, R1.7, R2.1, R2.2, R2.3, R2.4, R2.5, R3.1, R3.1.1, R3.1.2, R3.1.3, R3.1.4, R3.1.5, R3.1.6, R3.2, R3.3, R3.4, R3.5, R3.6, R3.7, R3.8, R3.9, R4.1, R4.2, R4.3, R4.4, R4.5, R4.6, R5.1, R5.2, R5.3, R5.4 | suite parsing, grid eval, all eval tools, result analysis |
| T5 | evaluator/*.yaml under AGENT_CATALOG_ROOT | two-level state machine: session states, point states, signals, lifecycle |
| T6 | testdata/integration/uc002-evaluator-benchmark/suite.yaml | profile-based qwen35b/qwen27b integration suite |
| T7 | generator/profile-qwen35b.yaml and profile-qwen27b.yaml | generator profiles under evaluation |
| T8 | srd018-cli-flag-contract R1.1, R1.2, R1.3, R1.4, R1.5, R1.6, R1.7, R1.8, R1.9, R1.10, R2.1, R2.2, R2.3, R2.4, R2.5, R3.1, R3.2 | standard flags, traceparent propagation, model/directory/runtime data, remaining budget |

<!-- field: integration_test -->
## Integration test

applications/coding-agent mage integration:criticGate

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Suite YAML is parsed and expanded into a profile/gridpoint/sample/repetition grid.

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC1

<!-- sc: S2 -->
### S2

Each grid point runs the per-point pipeline to completion via nested core.Loop.

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC2

<!-- sc: S3 -->
### S3

Generator subprocess receives OTel traceparent for trace correlation.

<!-- field: traces -->
#### Traces

- srd016-traceparent AC1

<!-- sc: S4 -->
### S4

Results include experiment.yaml with the resolved profile path and model for reproducibility.

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC2

<!-- sc: S5 -->
### S5

Pass/fail/timeout classification is correct based on exit code and test output.

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC4

<!-- sc: S6 -->
### S6

Metrics are collected even when the generator fails or times out.

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC2

<!-- sc: S7 -->
### S7

Session summary reports total, passed, failed, and timed-out counts.

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC1

<!-- sc: S8 -->
### S8

Result directory contains meta.json and experiment.yaml for each point.

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC2

<!-- sc: S9 -->
### S9

Every attempted point has machine-readable failure evidence, and the release gate reports each failed point while rejecting any failed benchmark.

<!-- field: traces -->
#### Traces

- srd019-eval-harness AC2
- srd019-eval-harness AC3

<!-- field: test_data -->
## Test data

testdata/integration/uc002-evaluator-benchmark

<!-- field: out_of_scope -->
## Out of scope

- Real-time progress streaming during evaluation.
- Parallel grid point execution (points run sequentially).
- Result visualization (handled by bench, see rel01.0-uc003).
- Performance regression analysis.

<!-- field: test_suite_trace -->
## Test suite trace

test-rel01.0

<!-- field: machine_configs -->
## Machine configs

- evaluator/machine.yaml
- evaluator/point.yaml
