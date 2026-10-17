---
id: sm-tool-metrics-monitoring
title: Tool Metrics Monitoring Semantic Model
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: purpose -->
## Purpose

Defines the vocabulary for tool metrics and embedded monitor snapshots.
Tools own metric meaning. Runtime code supplies dispatch context and a
shared recorder path. Monitor readers consume cached state from the embedded
store.

This model defines implemented Release 04.0 language for ToolDef metric
fields, MachineSpec workflow labels, OpenTelemetry metric export, and REST or
OpenAPI monitor reads. Recorder ports, stores, exporters, and REST readers are
implemented by the packages named in the related SRDs.

The value is shared wording: SRDs, use cases, and test suites can talk about
metric samples, recorder ports, workflow labels, and snapshots without
redefining those terms.

<!-- field: requirements_source | form: "tree" -->
## Requirements source

- canonical:
  1. docs/ARCHITECTURE.md
  1. docs/specs/semantic-models/tool-language.md
  1. docs/specs/semantic-models/rest-tool-language.md
  1. docs/specs/semantic-models/rollback-lifecycle.md
  1. docs/specs/config-formats/tool-declaration-format.md
  1. docs/specs/config-formats/machine-format.md
  1. docs/specs/software-requirements/srd032-tool-metrics-monitor-store.md
  1. docs/specs/software-requirements/srd033-monitor-rest-api.md

<!-- field: vocabulary -->
## Vocabulary

<!-- field: standard_tool_metric -->
### Standard tool metric

A metric every tool dispatch may report with the same meaning across tool
families. Common examples are duration, dispatch counts, safe byte counts,
and retry counts.

<!-- field: specialized_tool_metric -->
### Specialized tool metric

A metric whose meaning belongs to one tool family or one declared tool.
Examples include files changed for filesystem tools, process exit code for
exec tools, provider tokens for model tools, and HTTP status code for REST
tools.

<!-- field: metric_source -->
### Metric source

The actor that produces a metric sample. Tool commands produce specialized
samples. Runtime dispatch code may add run identity, state, signal, tool
name, status, and duration.

<!-- field: recorder_port -->
### Recorder port

Optional runtime port that accepts normalized metric samples. Tools emit
through this port instead of writing directly to a store, exporter, or REST
handler.

<!-- field: metric_instrument -->
### Metric instrument

Declared metric stream. Each instrument has a name, kind, unit, value
source, and allowed attributes. Instrument kinds include counters,
histograms, and gauges.

<!-- field: metric_attribute -->
### Metric attribute

Low-cardinality name-value label attached to a metric sample. Attributes
identify stable dimensions. Tool family and state are examples. Signal,
status, and workflow labels are also valid when their value sets stay
bounded.

<!-- field: workflow_label -->
### Workflow label

MachineSpec or transition metadata that groups metric samples by declared
workflow meaning. Workflow labels come from trusted machine YAML, not tool
output or model text.

<!-- field: monitor_snapshot -->
### Monitor snapshot

Cached view of recent metric state held by the embedded monitor store. A
snapshot can include run summary, tool aggregates, recent samples, schema
metadata, and monitor diagnostics.

<!-- field: monitor_reader -->
### Monitor reader

Consumer that reads cached monitor state. REST handlers and OpenTelemetry
export paths are monitor readers. Local UIs and supervisors may read the
same snapshots.

<!-- field: metric_schema -->
### Metric schema

Declared shape for one metric instrument, including value type, unit,
allowed attributes, source, and redaction policy. Schema validation belongs
to configuration loading or monitor setup, not to the LLM.

<!-- field: ownership_rules -->
## Ownership rules

<!-- field: tools -->
### Tools

Tools own specialized metric meaning and choose which tool-local values become samples.

<!-- field: runtime_dispatch -->
### Runtime dispatch

Runtime dispatch code owns standard context that is independent of tool family.

<!-- field: tooldef_metrics -->
### Tooldef metrics

ToolDef metrics config declares instrument names, units, value sources, and allowed attributes.

<!-- field: machine_labels -->
### Machine labels

MachineSpec config may add workflow labels for states or transitions.

<!-- field: recorder_port -->
### Recorder port

The recorder port owns sample normalization and rejects samples that violate metric schema.

<!-- field: monitor_store -->
### Monitor store

The embedded monitor store owns cached recent state, not durable history.

<!-- field: otel_export -->
### OTel export

OpenTelemetry export maps normalized samples to OTel instruments and attributes.

<!-- field: rest_readers -->
### REST readers

REST and OpenAPI monitor readers read cached snapshots only.

<!-- field: reader_limits -->
### Reader limits

Monitor readers do not execute tools, inspect prompts, or scrape workspace files.

<!-- field: standard_metric_set | form: "keyed table" -->
## Standard metric set

| id | kind | unit | source | meaning |
| --- | --- | --- | --- | --- |
| dispatch_duration | histogram | ms | runtime_dispatch | Wall-clock duration for one command dispatch. |
| dispatch_count | counter | "{dispatch}" | runtime_dispatch | Count of tool dispatch attempts. |
| dispatch_success | counter | "{dispatch}" | runtime_dispatch | Count of dispatches whose Result.Signal indicates success for the declared tool. |
| dispatch_failure | counter | "{dispatch}" | runtime_dispatch | Count of command errors, failed signals, panics, or timeouts. |
| input_bytes | histogram | By | tool_or_runtime | Size of structured or serialized input when the tool can report it safely. |
| output_bytes | histogram | By | tool_or_runtime | Size of safe structured output or redacted output preview. |
| retry_count | counter | "{retry}" | tool | Count of retries attempted inside one declared tool word. |

<!-- field: relationships | form: "keyed table" -->
## Relationships

| id | meaning | invariant |
| --- | --- | --- |
| ToolDef_to_metric_schema | ToolDef metrics config declares instruments and allowed attributes for one vocabulary word. | ToolDef config cannot make the LLM choose metric names, labels, or secret material. |
| Command_to_recorder_port | A command emits standard or specialized samples through the shared recorder port. | Commands do not write directly to REST handlers, OpenTelemetry exporters, or monitor store internals. |
| Runtime_dispatch_to_standard_metrics | Dispatch can add duration and status context. It can also add state, signal, and tool identity. | Runtime code does not invent specialized tool-family metric meaning. |
| MachineSpec_to_workflow_labels | Machine or transition metadata can attach stable workflow labels to samples. | Workflow labels are trusted config values and are not derived from model text or tool output. |
| Recorder_port_to_monitor_store | The recorder validates and normalizes samples before the store caches them. | Invalid samples are rejected or recorded as diagnostics without failing the agent loop. |
| Monitor_store_to_monitor_snapshot | The store builds snapshots from bounded recent samples and aggregate state. | Snapshots are cached monitor views, not durable audit logs. |
| Monitor_store_to_OpenTelemetry | Export reads normalized samples or aggregates and maps them to OpenTelemetry metric instruments. | Export failure does not change tool execution results. |
| Monitor_store_to_REST_API | REST and OpenAPI handlers return cached snapshots, schema, and recent sample windows. | Read handlers do not execute tools or expose prohibited labels. |
| Metrics_to_rollback_lifecycle | Metric samples are observability records and do not replace checkpoints, history, UndoMemento data, or workspace refs. | Rollback uses lifecycle state. Monitor snapshots may help operators inspect behavior but do not restore state. |

<!-- field: specialized_metric_families | form: "tree" -->
## Specialized metric families

- filesystem:
  - examples: files_read, files_written, bytes_changed, path_policy_denials.
  - owner: internal/tools/filesystem declarations and commands.
- exec:
  - examples: process_exit_code, stdout_bytes, stderr_bytes, elapsed_process_ms.
  - owner: internal/tools/exec declarations and commands.
- llm:
  - examples: prompt_tokens, completion_tokens, model_latency_ms, parse_error_count.
  - owner: internal/tools/llm and model-boundary declarations.
- rest:
  - examples:
    - status: http_status_code
    - request_size: request_bytes
    - response_size: response_bytes
    - redirect: redirect_count
    - retry: retry_count
  - owner: internal/tools/rest declarations and commands.
- lifecycle:
  - examples: checkpoint_count, checkpoint_bytes, rollback_steps, restore_duration_ms.
  - owner: internal/tools/lifecycle declarations and commands.

<!-- field: attribute_rules -->
## Attribute rules

<!-- field: allowed_sources -->
### Allowed sources

Attribute values may come from trusted ToolDef config, trusted MachineSpec
config, bounded runtime enums, or normalized operation names declared in
REST or LLM config.

<!-- field: prohibited_sources -->
### Prohibited sources

Attribute values must not come from raw prompts, full model responses,
full tool output, or secret material. They also must not come from
unbounded paths or arbitrary URLs. User free text, stack traces, and
command output are also prohibited.

<!-- field: path_policy -->
### Path policy

Path labels must use bounded names from config, path class values, or
redacted path hashes when an SRD allows them. Raw workspace paths are not
metric attributes.

<!-- field: cardinality_policy -->
### Cardinality policy

Attribute value sets should be small and predictable. Tool names, states,
signals, statuses, configured operation IDs, and fixed workflow labels are
allowed. Request IDs and timestamps are not allowed as metric labels.
Filenames, prompt fragments, and output fragments are also not allowed.

<!-- field: redaction_policy -->
### Redaction policy

Redaction happens before a sample reaches the recorder. If a value cannot
be safely classified, the tool should omit the attribute and emit a monitor
diagnostic instead.

<!-- field: monitor_snapshot_shape -->
## Monitor snapshot shape

<!-- field: fields -->
### Fields

<!-- field: run_id -->
#### Run ID

Runtime run identifier or configured run label.

<!-- field: status -->
#### Status

Current or final run status known to the monitor.

<!-- field: updated_at -->
#### Updated at

Timestamp for the latest accepted sample.

<!-- field: instruments -->
#### Instruments

Known metric schemas selected for the run.

<!-- field: aggregates -->
#### Aggregates

Bounded counters, histograms, gauges, and summaries by allowed attributes.

<!-- field: recent_samples -->
#### Recent samples

Optional recent sample window for debugging and UI display.

<!-- field: diagnostics -->
#### Diagnostics

Recorder, store, export, or schema warnings.

<!-- field: non_goals -->
## Non goals

- This model does not make metrics a durable audit log or replace RunResult, History, checkpoints, or OpenTelemetry traces.
- Metric labels must not be defined by model text or tool output. Labels must not contain raw prompts, full outputs, secrets, or unbounded paths.

<!-- field: related_documents -->
## Related documents

- sm-tool-language
- sm-rest-tool-language
- sm-rollback-lifecycle
- cfg-tool-declaration-format
- cfg-machine-format
