---
id: rel04.0-uc001-configurable-tool-metrics
title: Profile configures tool metrics for embedded monitoring
status: done
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A profile author enables embedded monitor metrics for a MachineSpec run.
ToolDef metrics declarations describe instruments, units, value sources, and
safe attributes. MachineSpec metric_labels add workflow context.

Runtime dispatch supplies standard context. Tools emit specialized samples
through ToolMetricsRecorder. The monitor store keeps bounded snapshots for a
live run, while OpenTelemetry export reads the same normalized samples without
changing command results.

Disabled mode stays quiet.

<!-- field: actor -->
## Actor

Agent profile author configuring tool metrics

<!-- field: trigger -->
## Trigger

The author wants progress and usage metrics for one configured agent without
changing tool workflows or giving model output control over metric labels.

<!-- field: flow -->
## Flow

- F1: The author adds ToolDef metrics declarations for selected tools.
- F2: The author adds MachineSpec metric_labels for workflow phases.
- F3: Startup validation checks instrument names, units, value sources, and attributes.
- F4: Runtime dispatch passes standard context to ToolMetricsRecorder.
- F5: Tools emit specialized metric samples through the recorder.
- F6: The monitor store caches current state, per-tool usage, accepted samples, diagnostics, and recent errors.
- F7: OpenTelemetry metric export reads normalized samples without changing command results.
- F8: With monitor config absent, runtime uses disabled-mode recorder behavior and no monitor store.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd032-tool-metrics-monitor-store R1.1, R1.2, R1.3, R1.4, R1.5, R1.6 | monitor store scope and bounded snapshots |
| T2 | srd032-tool-metrics-monitor-store R2.1, R2.2, R2.3, R2.4, R2.5, R2.6, R2.7 | recorder ownership boundaries |
| T3 | srd032-tool-metrics-monitor-store R3.1, R3.2, R3.3, R3.4, R3.5, R3.6, R3.7, R3.8 | metric validation |
| T4 | srd032-tool-metrics-monitor-store R4.1, R4.2, R4.3, R4.4, R4.5 | OpenTelemetry metric export |
| T5 | srd032-tool-metrics-monitor-store R5.1, R5.2, R5.3, R5.4, R5.5 | disabled-mode behavior |
| T6 | sm-tool-metrics-monitoring ownership_rules | tools own specialized metrics and runtime owns standard context |
| T7 | cfg-tool-declaration-format metric_declaration_fields | ToolDef metrics declarations |
| T8 | cfg-machine-format workflow_metric_labels | MachineSpec workflow labels |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Monitor store snapshots expose bounded current and recent state without executing tools or reading workspace files.

<!-- field: traces -->
#### Traces

- srd032-tool-metrics-monitor-store AC1

<!-- sc: S2 -->
### S2

Recorder use separates runtime-owned standard context from tool-owned specialized metric meaning.

<!-- field: traces -->
#### Traces

- srd032-tool-metrics-monitor-store AC2

<!-- sc: S3 -->
### S3

Metric config validation rejects unsafe names, source paths, attributes, and label values before recording starts.

<!-- field: traces -->
#### Traces

- srd032-tool-metrics-monitor-store AC3

<!-- sc: S4 -->
### S4

OpenTelemetry metric export reads normalized samples and reports export failure as diagnostics.

<!-- field: traces -->
#### Traces

- srd032-tool-metrics-monitor-store AC4

<!-- sc: S5 -->
### S5

Disabled mode preserves current loop, dispatch, tracing, and tool behavior without allocating a monitor store.

<!-- field: traces -->
#### Traces

- srd032-tool-metrics-monitor-store AC5

<!-- field: out_of_scope -->
## Out of scope

- Defining REST or OpenAPI monitor read endpoints.
- Implementing monitor packages.
- Creating durable metric history.
- Replacing traces, RunResult, checkpoints, or rollback state.

<!-- field: semantic_model -->
## Semantic model

sm-tool-metrics-monitoring

<!-- field: config_formats -->
## Config formats

- cfg-tool-declaration-format
- cfg-machine-format
