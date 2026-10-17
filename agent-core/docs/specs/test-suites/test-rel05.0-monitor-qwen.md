---
id: test-rel05.0-monitor-qwen
title: Test Suite for Release 05.0 Qwen Monitor Integration
release: "05.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel05.0-uc001-qwen-monitor-token-progress

<!-- field: preconditions -->
## Preconditions

- Ollama is reachable on http://localhost:11434.
- The configured Qwen 3 family model is present in Ollama.
- The embedded monitor runtime is enabled by monitor REST views or declared metrics.
- The integration target can bind an isolated loopback monitor address.

<!-- field: test_cases -->
## Test cases

<!-- case: Qwen monitor token metrics increase during live run | status: implemented -->
### Qwen monitor token metrics increase during live run

Generate a temporary Qwen/Ollama REST profile that starts the embedded
monitor through MachineSpec, continues into the existing Qwen REST
model-listing flow, and keeps the monitor listener alive until the test
sends a control-exit request. Poll /monitor/metrics during the live
process and require the combined llm.prompt_tokens and
llm.completion_tokens aggregate to increase from the initial value.

<!-- field: use_case -->
#### Use case

rel05.0-uc001-qwen-monitor-token-progress

<!-- field: go_test -->
#### Go test

mage integration:ollamaMonitor

<!-- field: inputs | form: "tree" -->
#### Inputs

- model_family: qwen3
- prompt: List the local Ollama models available on this machine.
- monitor_routes:
  1. /monitor/state
  1. /monitor/metrics
  1. /monitor/control/exit
- token_metrics:
  1. llm.prompt_tokens
  1. llm.completion_tokens
- binding: isolated_loopback_address

<!-- field: expected | form: "tree" -->
#### Expected

- request_prompt_passed_as_file: true
- startup_exit_preempts_readiness_timeout: true
- monitor_launch_non_blocking: true
- qwen_flow_continues_after_server_launched: true
- token_metric_present: true
- token_metric_total_increases: true
- process_running_while_polled: true
- control_exit_stops_monitor: true
- terminal_status_succeeded: true

<!-- field: traces -->
#### Traces

- rel05.0-uc001-qwen-monitor-token-progress S1
- rel05.0-uc001-qwen-monitor-token-progress S2
- rel05.0-uc001-qwen-monitor-token-progress S3
- rel05.0-uc001-qwen-monitor-token-progress S4
- srd032-tool-metrics-monitor-store AC2
- srd033-monitor-rest-api AC1
- srd033-monitor-rest-api AC3
