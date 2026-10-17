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
| T1 | [srd032-tool-metrics-monitor-store](../software-requirements/srd032-tool-metrics-monitor-store.md) [R1.1](../software-requirements/srd032-tool-metrics-monitor-store.md#r11), [R1.2](../software-requirements/srd032-tool-metrics-monitor-store.md#r12), [R1.3](../software-requirements/srd032-tool-metrics-monitor-store.md#r13), [R1.4](../software-requirements/srd032-tool-metrics-monitor-store.md#r14), [R1.5](../software-requirements/srd032-tool-metrics-monitor-store.md#r15), [R1.6](../software-requirements/srd032-tool-metrics-monitor-store.md#r16) | monitor store scope and bounded snapshots |
| T2 | [srd032-tool-metrics-monitor-store](../software-requirements/srd032-tool-metrics-monitor-store.md) [R2.1](../software-requirements/srd032-tool-metrics-monitor-store.md#r21), [R2.2](../software-requirements/srd032-tool-metrics-monitor-store.md#r22), [R2.3](../software-requirements/srd032-tool-metrics-monitor-store.md#r23), [R2.4](../software-requirements/srd032-tool-metrics-monitor-store.md#r24), [R2.5](../software-requirements/srd032-tool-metrics-monitor-store.md#r25), [R2.6](../software-requirements/srd032-tool-metrics-monitor-store.md#r26), [R2.7](../software-requirements/srd032-tool-metrics-monitor-store.md#r27) | recorder ownership boundaries |
| T3 | [srd032-tool-metrics-monitor-store](../software-requirements/srd032-tool-metrics-monitor-store.md) [R3.1](../software-requirements/srd032-tool-metrics-monitor-store.md#r31), [R3.2](../software-requirements/srd032-tool-metrics-monitor-store.md#r32), [R3.3](../software-requirements/srd032-tool-metrics-monitor-store.md#r33), [R3.4](../software-requirements/srd032-tool-metrics-monitor-store.md#r34), [R3.5](../software-requirements/srd032-tool-metrics-monitor-store.md#r35), [R3.6](../software-requirements/srd032-tool-metrics-monitor-store.md#r36), [R3.7](../software-requirements/srd032-tool-metrics-monitor-store.md#r37), [R3.8](../software-requirements/srd032-tool-metrics-monitor-store.md#r38) | metric validation |
| T4 | [srd032-tool-metrics-monitor-store](../software-requirements/srd032-tool-metrics-monitor-store.md) [R4.1](../software-requirements/srd032-tool-metrics-monitor-store.md#r41), [R4.2](../software-requirements/srd032-tool-metrics-monitor-store.md#r42), [R4.3](../software-requirements/srd032-tool-metrics-monitor-store.md#r43), [R4.4](../software-requirements/srd032-tool-metrics-monitor-store.md#r44), [R4.5](../software-requirements/srd032-tool-metrics-monitor-store.md#r45) | OpenTelemetry metric export |
| T5 | [srd032-tool-metrics-monitor-store](../software-requirements/srd032-tool-metrics-monitor-store.md) [R5.1](../software-requirements/srd032-tool-metrics-monitor-store.md#r51), [R5.2](../software-requirements/srd032-tool-metrics-monitor-store.md#r52), [R5.3](../software-requirements/srd032-tool-metrics-monitor-store.md#r53), [R5.4](../software-requirements/srd032-tool-metrics-monitor-store.md#r54), [R5.5](../software-requirements/srd032-tool-metrics-monitor-store.md#r55) | disabled-mode behavior |
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

- [srd032-tool-metrics-monitor-store](../software-requirements/srd032-tool-metrics-monitor-store.md) [AC1](../software-requirements/srd032-tool-metrics-monitor-store.md#ac1)

<!-- sc: S2 -->
### S2

Recorder use separates runtime-owned standard context from tool-owned specialized metric meaning.

<!-- field: traces -->
#### Traces

- [srd032-tool-metrics-monitor-store](../software-requirements/srd032-tool-metrics-monitor-store.md) [AC2](../software-requirements/srd032-tool-metrics-monitor-store.md#ac2)

<!-- sc: S3 -->
### S3

Metric config validation rejects unsafe names, source paths, attributes, and label values before recording starts.

<!-- field: traces -->
#### Traces

- [srd032-tool-metrics-monitor-store](../software-requirements/srd032-tool-metrics-monitor-store.md) [AC3](../software-requirements/srd032-tool-metrics-monitor-store.md#ac3)

<!-- sc: S4 -->
### S4

OpenTelemetry metric export reads normalized samples and reports export failure as diagnostics.

<!-- field: traces -->
#### Traces

- [srd032-tool-metrics-monitor-store](../software-requirements/srd032-tool-metrics-monitor-store.md) [AC4](../software-requirements/srd032-tool-metrics-monitor-store.md#ac4)

<!-- sc: S5 -->
### S5

Disabled mode preserves current loop, dispatch, tracing, and tool behavior without allocating a monitor store.

<!-- field: traces -->
#### Traces

- [srd032-tool-metrics-monitor-store](../software-requirements/srd032-tool-metrics-monitor-store.md) [AC5](../software-requirements/srd032-tool-metrics-monitor-store.md#ac5)

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
