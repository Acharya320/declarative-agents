---
id: test-rel08.0-observer
title: Test Suite for Release 08.0 Observer Agent
release: "08.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel08.0-uc001-mesh-observer

<!-- field: preconditions -->
## Preconditions

- agent-core REST client get, REST server launch, await, and stop words are available.
- The observer profile mounts a ServiceAccount with namespace-scoped pod list/get RBAC.
- At least the chatbot and one RAG server are running in the demo cluster.

<!-- field: test_cases -->
## Test cases

<!-- case: Observer profile validates | status: implemented -->
### Observer profile validates

The observer profile loads and every action, signal, declaration, REST client, REST server route, and monitor route validates. Boot smoke passes.

<!-- field: use_case -->
#### Use case

rel08.0-uc001-mesh-observer

<!-- field: integration_target -->
#### Integration target

mage audit

<!-- field: expected | form: "tree" -->
#### Expected

- profile_loads: true
- machine_validates: true
- rest_definitions_validate: true
- boot_smoke_passes: true

<!-- field: traces -->
#### Traces

- rel08.0-uc001-mesh-observer S1
- srd008-observer-agent AC1

<!-- case: Observer discovers mesh agent pods on the demo cluster | status: done -->
### Observer discovers mesh agent pods on the demo cluster

On the running demo cluster the observer discovers its agent pods through
the kubectl-proxy sidecar, which holds the pod's service-account token. The
discovery result carries pod names, addresses, and each pod's monitor-port
label. Only agent pods are selected: the release selector is composed with
the monitor-port label so the Dolt and Chroma StatefulSets, which serve no
monitor surface, are excluded. Discovery follows membership, tracking pod
replacement and a Deployment scaled to zero and back.

<!-- field: use_case -->
#### Use case

rel08.0-uc001-mesh-observer

<!-- field: integration_target -->
#### Integration target

mage demo:up

<!-- field: expected | form: "tree" -->
#### Expected

- chatbot_discovered: true
- rag_server_discovered: true
- discovery_by_label_selector: true
- infrastructure_pods_excluded: true
- membership_follows_pod_churn: true

<!-- field: traces -->
#### Traces

- rel08.0-uc001-mesh-observer S1
- srd008-observer-agent AC1

<!-- case: Observer serves four monitor fan-in joins and survives degraded discovery | status: done -->
### Observer serves four monitor fan-in joins and survives degraded discovery

The observer's monitor fan-in is four MachineSpec for_each passes, one per
monitor endpoint, each expanding a single-request read word over the
discovered pods. The profile validates that grammar at load, and the served
/monitor/fleet view carries all four join labels beside the discovery
labels. With the kube client pointed at a non-routable address the pod list
never resolves, so the fan-in records nothing and the cycle still proceeds
to metrics, waits, and exits cleanly rather than failing.

<!-- field: use_case -->
#### Use case

rel08.0-uc001-mesh-observer

<!-- field: integration_target -->
#### Integration target

mage integration:observer

<!-- field: expected | form: "tree" -->
#### Expected

- four_join_labels_served: true
- one_request_per_dispatch: true
- degraded_discovery_does_not_fail_cycle: true
- per_item_counts_consistent: true
- clean_lifecycle_exit: true

<!-- field: traces -->
#### Traces

- rel08.0-uc001-mesh-observer S2
- srd008-observer-agent AC2

<!-- case: Observer polls agent monitor endpoints and aggregates fleet data | status: implemented -->
### Observer polls agent monitor endpoints and aggregates fleet data

Helm smoke reads the in-cluster observer and proves each of the four joins
contains every discovered agent exactly once with AgentMonitorRead, HTTP
200, and the pod's selected authority. It then gives only the observer an
authorized but unserved monitor port and proves every other member remains
reachable while a complete later cycle advances through metrics.

<!-- field: use_case -->
#### Use case

rel08.0-uc001-mesh-observer

<!-- field: integration_target -->
#### Integration target

mage integration:helmSmoke

<!-- field: expected | form: "tree" -->
#### Expected

- machine_spec_returned: true
- current_state_returned: true
- tools_returned: true
- events_returned: true
- unreachable_agent_degraded: true
- one_item_per_discovered_pod: true
- per_pod_monitor_port_resolved: true

<!-- field: traces -->
#### Traces

- rel08.0-uc001-mesh-observer S2
- srd008-observer-agent AC2

<!-- case: Fleet UX renders a card per pod from the four fan-in joins | status: done -->
### Fleet UX renders a card per pod from the four fan-in joins

fleetData zips the four monitor fan-in joins into one card per discovered
pod, joining them on the pod name each item carries, and maps each
endpoint's payload to the shape the card reads. A pod whose state read
failed renders unreachable while the others render normally, even though
the join counts every dispatched item as succeeded. Absent or unavailable
joins yield the empty fleet view rather than an exception.

<!-- field: use_case -->
#### Use case

rel08.0-uc001-mesh-observer

<!-- field: go_test -->
#### Go test

cd ../ui-kit && npm test

<!-- field: expected | form: "tree" -->
#### Expected

- card_per_pod: true
- machine_tools_state_events_mapped: true
- pod_name_keys_resource_lookup: true
- failed_state_read_renders_unreachable: true
- absent_joins_degrade_to_empty: true

<!-- field: traces -->
#### Traces

- rel08.0-uc001-mesh-observer S2
- rel08.0-uc001-mesh-observer S3
- srd008-observer-agent AC3

<!-- case: Fleet UX serves and displays agent cards | status: done -->
### Fleet UX serves and displays agent cards

The observer UX is reachable at http://observer.localhost/, which redirects
to /ui/ and serves the observer bundle compiled into agent-core, configured
by the static_assets config map in agents/observer/monitor-rest.yaml
(applications srd004 R9). The fleet payload behind it carries the four
fan-in joins the cards read, one item per discovered agent. Rendering those
joins into cards is covered by the ui-kit tests; TestObserverMonitorEndpoints
boots the observer and reads /ui/ and ui-config.json; this case covers the
served surface.

<!-- field: use_case -->
#### Use case

rel08.0-uc001-mesh-observer

<!-- field: integration_target -->
#### Integration target

mage demo:up

<!-- field: expected | form: "tree" -->
#### Expected

- ux_reachable: true
- auto_refresh: true
- served_from_embedded_bundle: true

<!-- field: traces -->
#### Traces

- rel08.0-uc001-mesh-observer S2
- rel08.0-uc001-mesh-observer S3
- srd008-observer-agent AC3

<!-- case: Observer serves the embedded bundle and proxies the agents its UI names | status: implemented -->
### Observer serves the embedded bundle and proxies the agents its UI names

agents/observer holds declaration files only. monitor-rest.yaml binds
/ui/{path...} to the observer bundle compiled into agent-core with a config
map naming the title, the chatbot's monitored agents, and the collector
trace backend, and declares GET /monitor-proxy/{agent}/{path...} with an
upstream for every agent that config names and no /trace-proxy route. The
chart co-generates the same file with upstreams at the in-cluster Services
derived from ragUnits and the collector toggle, and renders no observer UI
ConfigMap, archive, volume, or root. The booted observer serves /ui/ and
ui-config.json and answers 404 for an undeclared proxy agent.

<!-- field: use_case -->
#### Use case

rel08.0-uc001-mesh-observer

<!-- field: go_test -->
#### Go test

go test ./magefiles -run 'TestObserver(UIRouteAndPackageContract|MonitorRestCoGeneratedForCluster|ChartDeliversNoUIFiles|MonitorEndpoints)'

<!-- field: expected | form: "tree" -->
#### Expected

- declaration_files_only: true
- bundle_binding_with_config: true
- monitor_proxy_covers_ui_agents: true
- no_trace_proxy: true
- cluster_upstreams_are_services: true
- no_chart_ui_delivery: true

<!-- field: traces -->
#### Traces

- rel08.0-uc001-mesh-observer S2
- srd008-observer-agent AC3

<!-- case: Observer RBAC is namespace-scoped and read-only | status: implemented -->
### Observer RBAC is namespace-scoped and read-only

The rendered Helm template creates a Role (not ClusterRole) with only get and list verbs on pods, services, deployments, and metrics.k8s.io pod metrics in the release namespace.

<!-- field: use_case -->
#### Use case

rel08.0-uc001-mesh-observer

<!-- field: go_test -->
#### Go test

go test ./magefiles -run TestObserverRBACRender

<!-- field: expected | form: "tree" -->
#### Expected

- role_not_cluster_role: true
- verbs_get_list_only: true
- resources_pods_services_deployments_podmetrics: true
- namespace_scoped: true

<!-- field: traces -->
#### Traces

- rel08.0-uc001-mesh-observer S5
- srd008-observer-agent AC5

<!-- case: Observer discovers topology and the UX groups pods by Deployment and Service | status: done -->
### Observer discovers topology and the UX groups pods by Deployment and Service

The retained demo observation proves the observer lists the release Deployments and Services alongside pods; the ui-kit Topology test proves the fleet view groups each pod under its owning Deployment and fronting Service. Configuration uses declarations and flags with no new environment variable. This case stays manual rather than attributing topology equality to Helm smoke, whose automated observer proof currently covers pods and fan-in.

<!-- field: use_case -->
#### Use case

rel08.0-uc001-mesh-observer

<!-- field: integration_target -->
#### Integration target

mage demo:up

<!-- field: go_test -->
#### Go test

cd ../ui-kit && npm test

<!-- field: expected | form: "tree" -->
#### Expected

- deployments_discovered: true
- services_discovered: true
- pods_grouped_by_deployment: true
- no_new_environment_variable: true

<!-- field: traces -->
#### Traces

- rel08.0-uc001-mesh-observer S6
- srd008-observer-agent AC6

<!-- case: Observer shows per-pod resource usage and degrades when metrics are absent | status: implemented -->
### Observer shows per-pod resource usage and degrades when metrics are absent

Helm smoke installs pinned test-only metrics-server and proves every
discovered agent has populated CPU and memory. The UI sums all container
quantities into pod totals and renders them. The persistent demo and UI
tests retain the absent-metrics behavior: topology and fan-in render while
cards show "metrics unavailable".

<!-- field: use_case -->
#### Use case

rel08.0-uc001-mesh-observer

<!-- field: integration_target -->
#### Integration target

mage integration:helmSmoke

<!-- field: go_test -->
#### Go test

cd ../ui-kit && npm test

<!-- field: expected | form: "tree" -->
#### Expected

- cpu_and_memory_per_pod: true
- metrics_unavailable_state_when_absent: true
- discovery_unaffected_by_metrics_failure: true

<!-- field: traces -->
#### Traces

- rel08.0-uc001-mesh-observer S7
- srd008-observer-agent AC7

<!-- case: Observer lifecycle exit stops cleanly | status: implemented -->
### Observer lifecycle exit stops cleanly

A lifecycle exit request stops the discovery loop, monitor, and control servers. The observer reaches terminal Done.

<!-- field: use_case -->
#### Use case

rel08.0-uc001-mesh-observer

<!-- field: integration_target -->
#### Integration target

mage integration:observer

<!-- field: expected | form: "tree" -->
#### Expected

- discovery_stopped: true
- monitor_stopped: true
- terminal_state_done: true

<!-- field: traces -->
#### Traces

- rel08.0-uc001-mesh-observer S4
- srd008-observer-agent AC4

<!-- case: Fleet updates during a deterministic live chat turn | status: implemented -->
### Fleet updates during a deterministic live chat turn

Helm swap holds the actual answer-model request open and proves newer
in-cluster observer state and event entries before rollout, while the
chatbot and rag0 remain reachable. The chatbot answer-path tuple must
remain in its fleet entry. Rag activity accepts that query tuple when it
remains in the bounded recent-event window; after eviction, the shared
collector must retain a post-baseline chatbot trace carrying the completed
`execute_tool rag_query` span. The held answer boundary follows that tool,
so reaching both proves the selected RAG request returned.
The observer fan-in iteration and reachability still have to advance, so
durable trace evidence supplements rather than replaces the fleet proof.

<!-- field: use_case -->
#### Use case

rel08.0-uc001-mesh-observer

<!-- field: integration_target -->
#### Integration target

mage integration:helmSwap

<!-- field: expected | form: "tree" -->
#### Expected

- chatbot_state_advanced: true
- chatbot_answer_path_event: true
- rag_query_proven: true
- rag_query_event_or_trace_after_eviction: true
- both_agents_reachable: true

<!-- field: traces -->
#### Traces

- rel08.0-uc001-mesh-observer S3
- srd008-observer-agent AC3
