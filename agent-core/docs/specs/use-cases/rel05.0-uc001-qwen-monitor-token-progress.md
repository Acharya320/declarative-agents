---
id: rel05.0-uc001-qwen-monitor-token-progress
title: Qwen run exposes increasing token metrics through embedded monitor
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A live Qwen/Ollama integration run starts the embedded monitor from the
MachineSpec before invoking the model. The monitor launch is non-blocking:
after ServerLaunched, the state machine continues through the normal Qwen REST
tool flow. The integration proof polls the monitor metrics endpoint while the
run is active and verifies that provider token metrics become visible and
increase over time.

<!-- field: actor -->
## Actor

Release maintainer proving live monitor observability on Qwen flows

<!-- field: trigger -->
## Trigger

UC005 already shows REST tool selection and final answers backed by model
data. Release 05.0 needs a live proof that the same model-boundary path can
expose monitor state while the state machine continues running, without a
sidecar or a blocking monitor action.

<!-- field: flow -->
## Flow

- F1: The maintainer starts local Ollama with the configured Qwen 3 family model available.
- F2: The integration target generates an isolated monitored Qwen profile with a free loopback monitor address.
- F3: The generated MachineSpec handles Seed by running launch_monitor_rest.
- F4: launch_monitor_rest emits ServerLaunched without waiting for monitor readers.
- F5: The state machine continues to invoke_llm and runs the Qwen Ollama REST model-listing flow.
- F6: invoke_llm records declared llm.prompt_tokens and llm.completion_tokens samples when provider usage is returned.
- F7: The integration target polls /monitor/metrics while the agent process remains active.
- F8: The observed token metric aggregate increases from its initial value.
- F9: The integration target posts /monitor/control/exit so the profile stops the monitor listener and exits succeeded.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | mage integration:ollamaMonitor | live monitored Ollama token progress proof |
| T2 | [rel03.0-uc005-ollama-openapi-model-listing](rel03.0-uc005-ollama-openapi-model-listing.md) | source Qwen REST tool flow |
| T3 | [rel04.0-uc002-monitor-rest-api](rel04.0-uc002-monitor-rest-api.md) | embedded monitor REST snapshot behavior |
| T4 | [srd032-tool-metrics-monitor-store](../software-requirements/srd032-tool-metrics-monitor-store.md) [R2.1](../software-requirements/srd032-tool-metrics-monitor-store.md#r21), [R2.2](../software-requirements/srd032-tool-metrics-monitor-store.md#r22), [R2.3](../software-requirements/srd032-tool-metrics-monitor-store.md#r23), [R2.4](../software-requirements/srd032-tool-metrics-monitor-store.md#r24), [R2.5](../software-requirements/srd032-tool-metrics-monitor-store.md#r25), [R2.6](../software-requirements/srd032-tool-metrics-monitor-store.md#r26), [R2.7](../software-requirements/srd032-tool-metrics-monitor-store.md#r27), [AC2](../software-requirements/srd032-tool-metrics-monitor-store.md#ac2) | tool-owned specialized metrics |
| T5 | [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [R1.1](../software-requirements/srd033-monitor-rest-api.md#r11), [R1.2](../software-requirements/srd033-monitor-rest-api.md#r12), [R1.3](../software-requirements/srd033-monitor-rest-api.md#r13), [R1.4](../software-requirements/srd033-monitor-rest-api.md#r14), [R1.5](../software-requirements/srd033-monitor-rest-api.md#r15), [R3.1](../software-requirements/srd033-monitor-rest-api.md#r31), [R3.2](../software-requirements/srd033-monitor-rest-api.md#r32), [R3.3](../software-requirements/srd033-monitor-rest-api.md#r33), [R3.4](../software-requirements/srd033-monitor-rest-api.md#r34), [R3.5](../software-requirements/srd033-monitor-rest-api.md#r35), [R3.6](../software-requirements/srd033-monitor-rest-api.md#r36), [R3.7](../software-requirements/srd033-monitor-rest-api.md#r37), [R3.8](../software-requirements/srd033-monitor-rest-api.md#r38), [R3.9](../software-requirements/srd033-monitor-rest-api.md#r39), [R3.10](../software-requirements/srd033-monitor-rest-api.md#r310), [R3.11](../software-requirements/srd033-monitor-rest-api.md#r311), [R3.12](../software-requirements/srd033-monitor-rest-api.md#r312), [AC1](../software-requirements/srd033-monitor-rest-api.md#ac1), [AC3](../software-requirements/srd033-monitor-rest-api.md#ac3) | cached monitor REST reads |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The monitor REST listener is launched by MachineSpec and becomes queryable before the Qwen run completes.

<!-- field: traces -->
#### Traces

- [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [AC1](../software-requirements/srd033-monitor-rest-api.md#ac1)

<!-- sc: S2 -->
### S2

The Qwen state machine continues into invoke_llm after ServerLaunched instead of blocking on monitor reads.

<!-- field: traces -->
#### Traces

- [srd024-blocking-input-tools](../software-requirements/srd024-blocking-input-tools.md) [AC1](../software-requirements/srd024-blocking-input-tools.md#ac1)
- [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [AC6](../software-requirements/srd033-monitor-rest-api.md#ac6)

<!-- sc: S3 -->
### S3

/monitor/metrics exposes llm.prompt_tokens or llm.completion_tokens from the live monitor store.

<!-- field: traces -->
#### Traces

- [srd032-tool-metrics-monitor-store](../software-requirements/srd032-tool-metrics-monitor-store.md) [AC2](../software-requirements/srd032-tool-metrics-monitor-store.md#ac2)
- [srd033-monitor-rest-api](../software-requirements/srd033-monitor-rest-api.md) [AC3](../software-requirements/srd033-monitor-rest-api.md#ac3)

<!-- sc: S4 -->
### S4

The observed token metric total increases while the agent process remains active.

<!-- field: traces -->
#### Traces

- [srd032-tool-metrics-monitor-store](../software-requirements/srd032-tool-metrics-monitor-store.md) [AC2](../software-requirements/srd032-tool-metrics-monitor-store.md#ac2)

<!-- field: out_of_scope -->
## Out of scope

- Pulling or installing Qwen models.
- Requiring streaming token telemetry from a provider before a model response completes.
- Changing the shared rel03 Qwen/Ollama profile.

<!-- field: semantic_model -->
## Semantic model

sm-tool-metrics-monitoring

<!-- field: config_format -->
## Config format

cfg-machine-format
