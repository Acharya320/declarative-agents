---
id: test-rel20.1-profile-derived-agent
title: Release 20.1 Profile-Derived Agent Test Suite
release: "20.1"
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: test_case_count | form: "scalar" -->
## Test case count

4

<!-- field: overview -->
## Overview

TC1 is the structural gate for the release 20.1 documents. TC2 tests
agent-core's public profile inspection: the servers, ports, control and
monitor servers, bind variables, and validation read from a profile tree.
TC3 tests the operator deriving the Deployment and Service from that
inspection, reporting an invalid profile, and following a ConfigMap
change. TC4 is the collector proof with an Agent of profile, image, and
work volume only.

<!-- field: traces -->
## Traces

- [rel20.1-uc001-profile-derived-agent](../use-cases/rel20.1-uc001-profile-derived-agent.md)

<!-- field: srds -->
## SRDs

- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md)

<!-- field: preconditions -->
## Preconditions

- Run commands from the repository root.
- Go, Mage, Helm, kubectl, kind, and Docker are available for the platform case.

<!-- field: test_cases -->
## Test cases

<!-- case: Release 20.1 documents parse and trace reciprocally | status: implemented -->
### Release 20.1 documents parse and trace reciprocally

<!-- field: id -->
#### ID

TC1

<!-- field: use_case -->
#### Use case

[rel20.1-uc001-profile-derived-agent](../use-cases/rel20.1-uc001-profile-derived-agent.md)

<!-- field: evidence_type -->
#### Evidence type

executable_structural

<!-- field: implemented_by -->
#### Implemented by

GH-2691

<!-- field: command -->
#### Command

```bash
cd agent-core && mage audit
```

<!-- field: expected | form: "tree" -->
#### Expected

- audit_passes: true
- srd009_use_case_and_suite_trace_each_other: true

<!-- field: traces -->
#### Traces

- [rel20.1-uc001-profile-derived-agent](../use-cases/rel20.1-uc001-profile-derived-agent.md) [S1](../use-cases/rel20.1-uc001-profile-derived-agent.md#s1)
- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC5](../software-requirements/srd009-agent-operator.md#ac5)

<!-- case: agent-core inspects and validates a profile's servers from a file tree | status: implemented -->
### agent-core inspects and validates a profile's servers from a file tree

<!-- field: id -->
#### ID

TC2

<!-- field: use_case -->
#### Use case

[rel20.1-uc001-profile-derived-agent](../use-cases/rel20.1-uc001-profile-derived-agent.md)

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2692

<!-- field: command -->
#### Command

```bash
cd agent-core && go test ./pkg/profileservers
```

<!-- field: expected | form: "tree" -->
#### Expected

- collector_profile_yields_three_servers_with_control_and_monitor: true
- agent_env_moves_a_port_over_its_default: true
- missing_default_two_health_servers_and_no_health_server_each_error: true
- validate_agrees_with_the_cli: true

<!-- field: traces -->
#### Traces

- [rel20.1-uc001-profile-derived-agent](../use-cases/rel20.1-uc001-profile-derived-agent.md) [S1](../use-cases/rel20.1-uc001-profile-derived-agent.md#s1)
- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC5](../software-requirements/srd009-agent-operator.md#ac5)

<!-- case: The reconciler derives the Deployment and Service from the profile, reports an invalid profile, and follows a ConfigMap change | status: implemented -->
### The reconciler derives the Deployment and Service from the profile, reports an invalid profile, and follows a ConfigMap change

<!-- field: id -->
#### ID

TC3

<!-- field: use_case -->
#### Use case

[rel20.1-uc001-profile-derived-agent](../use-cases/rel20.1-uc001-profile-derived-agent.md)

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2694

<!-- field: command -->
#### Command

```bash
cd operator && go test ./internal/controller -run 'TestReconcile(Derives|Keeps|Passes|ReportsProfileInvalid|Follows)'
```

<!-- field: expected | form: "tree" -->
#### Expected

- ports_probe_monitor_label_and_bind_env_come_from_the_profile: true
- an_explicit_bind_host_in_env_is_kept: true
- args_service_account_and_claim_reach_the_pod: true
- invalid_profile_reports_profile_invalid_and_creates_nothing: true
- configmap_change_enqueues_the_agents_naming_it: true

<!-- field: traces -->
#### Traces

- [rel20.1-uc001-profile-derived-agent](../use-cases/rel20.1-uc001-profile-derived-agent.md) [S2](../use-cases/rel20.1-uc001-profile-derived-agent.md#s2)
- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC5](../software-requirements/srd009-agent-operator.md#ac5)
- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC6](../software-requirements/srd009-agent-operator.md#ac6)

<!-- case: The collector Agent of profile, image, and work volume becomes Ready and answers its health endpoint | status: implemented -->
### The collector Agent of profile, image, and work volume becomes Ready and answers its health endpoint

<!-- field: id -->
#### ID

TC4

<!-- field: use_case -->
#### Use case

[rel20.1-uc001-profile-derived-agent](../use-cases/rel20.1-uc001-profile-derived-agent.md)

<!-- field: evidence_type -->
#### Evidence type

integration

<!-- field: implemented_by -->
#### Implemented by

GH-2694

<!-- field: command -->
#### Command

```bash
cd operator && mage integration:all
```

<!-- field: expected | form: "tree" -->
#### Expected

- agent_manifest_declares_no_ports_control_port_or_bind_host: true
- agent_ready_condition_true: true
- health_answers_through_the_service_on_the_profile_port: true
- delete_removes_deployment_and_service: true

<!-- field: traces -->
#### Traces

- [rel20.1-uc001-profile-derived-agent](../use-cases/rel20.1-uc001-profile-derived-agent.md) [S3](../use-cases/rel20.1-uc001-profile-derived-agent.md#s3)
- [srd009-agent-operator](../software-requirements/srd009-agent-operator.md) [AC5](../software-requirements/srd009-agent-operator.md#ac5)
