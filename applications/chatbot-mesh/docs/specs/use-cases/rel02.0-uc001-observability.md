---
id: rel02.0-uc001-observability
title: Observability Panel Shows Per-Agent State and a Cross-Agent Trace Waterfall
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

An operator opens the chatbot's observability panel (srd002 R5.1) and watches
the mesh serve a turn. The panel renders one live sub-panel per monitored agent,
each showing the agent's run state, a metric-sample counter, and a feed of
run-event state transitions, fed by that agent's monitor read_state and SSE
stream read through the chatbot's same-origin monitor_proxy so no read crosses
origin. Selecting a turn renders its connected cross-agent trace as a span
waterfall grouped and colored by service.name, fetched from the trace backend
(Jaeger query API) same-origin through the monitor_proxy. The connected trace is
real because the request-scoped machine exports its spans under the
machine_request span and propagates traceparent to the RAG servers; the panel
degrades to a backend-missing notice when the trace backend is absent. This is
the observability epic (#295). Integration runs also send attributable traces
and metrics through a shared host ingress so a failed disposable-cluster run
remains diagnosable after teardown.

<!-- field: actor -->
## Actor

An operator watching the mesh serve turns.

<!-- field: trigger -->
## Trigger

The operator opens the observability route and selects a chat turn.

<!-- field: preconditions -->
## Preconditions

- The chatbot serves the SPA and the monitor_proxy binding to each monitored agent's monitor server.
- Each agent exports OTLP and a trace backend is reachable for the waterfall (else the panel degrades).
- The chat response carries the turn's trace_id.
- Integration runs have a healthy shared OTLP ingress with Jaeger and Prometheus-compatible backends.

<!-- field: flow -->
## Flow

- F1: The panel opens one sub-panel per monitored agent (chatbot and each RAG) from ui.yaml monitored_agents, reading each agent's monitor state and SSE stream through the same-origin monitor_proxy.
- F2: Each sub-panel shows the agent's run state, a metric-sample counter, and a feed of run-event state transitions.
- F3: A turn-correlation bar records each turn's time window and highlights the events observed during a selected turn across every sub-panel.
- F4: Selecting a turn fetches its connected trace by trace_id from the trace backend through the monitor_proxy and renders a span waterfall grouped and colored by service.name.
- F5: When the trace backend is absent the waterfall degrades to a backend-missing notice, leaving the per-agent monitor panels live.
- F6: For an integration run, the launcher records the run identity and queries retained agent, Chroma, and Kubernetes traces plus agent OTLP metrics after the scenario processes and kind cluster stop; the in-cluster Dolt Prometheus scrape is retired in agent mode (GH-1366).

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd002-chatbot-agent AC5 | The observability panel is one of the three declared panels served over static_assets. |
| T2 | test-rel02.0-observability observability_panels | The suite records the per-agent monitor panels and the cross-agent trace waterfall with expected outcomes. |
| T3 | srd003-chatbot-deployment AC5 | The deployment contract defines the integration-only shared sink, service identity, backend signals, and retained evidence. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The panel renders one live sub-panel per monitored agent, each showing run state, a metric-sample counter, and a run-event feed, read through the same-origin monitor_proxy.

<!-- field: traces -->
#### Traces

- srd002-chatbot-agent R5.1
- srd002-chatbot-agent AC5

<!-- sc: S2 -->
### S2

Selecting a turn renders its connected cross-agent trace as a span waterfall grouped by service.name, fetched from the trace backend through the monitor_proxy.

<!-- field: traces -->
#### Traces

- srd002-chatbot-agent R5.1
- srd002-chatbot-agent AC5

<!-- sc: S3 -->
### S3

Cross-agent trace continuity relies on traceparent propagation; when the trace backend is absent the waterfall degrades to a notice and the monitor panels stay live.

<!-- field: traces -->
#### Traces

- srd002-chatbot-agent R5.3
- srd002-chatbot-agent AC5

<!-- sc: S4 -->
### S4

An integration run retains filterable traces and metrics after scenario teardown, including stable agent and Chroma services, separately correlated control-plane traces, and agent OTLP metrics; Dolt metrics apply only under the opt-in contrib collector (GH-1366).

<!-- field: traces -->
#### Traces

- srd003-chatbot-deployment R5.5
- srd003-chatbot-deployment R5.6
- srd003-chatbot-deployment R5.7
- srd003-chatbot-deployment R5.8
- srd003-chatbot-deployment R5.9
- srd003-chatbot-deployment R5.10
- srd003-chatbot-deployment AC5

<!-- field: out_of_scope -->
## Out of scope

- Request-id-tagged turn correlation, a follow-on that needs the monitor run events to carry the request id.
- The provisioning panel, part of the control-plane surface (provisioning-workflow-orchestrator, creator, applier) and out of scope for this data-plane use case.
- OTel Logs and native Ollama server tracing.

<!-- field: test_suite_trace -->
## Test suite trace

test-rel02.0-observability
