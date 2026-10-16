---
id: test-rel99.0-observer
title: Test Suite for Release 99.0 Observer Family Promotion
release: "99.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel99.0-uc001-observer-family

<!-- field: preconditions -->
## Preconditions

- srd021-observer defines the generalized catalog observer family.
- srd008-observer-agent defines the chatbot-mesh observer the family generalizes.
- The catalog observer profile, per-application workloads, and realization rows are the adoption implementation, tracked separately.

<!-- field: test_cases -->
## Test cases

<!-- case: Generalized observer-family SRD with parameterized defaults | status: planned -->
### Generalized observer-family SRD with parameterized defaults

srd021-observer specifies parameterized namespace and label selector
declaration defaults and states the no-new-environment-variable rule; the
spec graph validates and the SRD is not orphaned.

<!-- field: use_case -->
#### Use case

rel99.0-uc001-observer-family

<!-- field: integration_target -->
#### Integration target

cd ../.. && mage audit

<!-- field: inputs -->
#### Inputs

```yaml
srd: srd021-observer
```

<!-- field: expected -->
#### Expected

```yaml
parameterized_defaults: true
no_new_environment_variable: true
```

<!-- field: traces -->
#### Traces

- rel99.0-uc001-observer-family S1
- srd021-observer AC1

<!-- case: Carried capabilities from the mesh observer | status: planned -->
### Carried capabilities from the mesh observer

srd021-observer carries topology discovery, monitor fan-in, resource
usage, and the served fleet UX forward from srd008 as generalized
requirements, with related_documents pointing at srd008.

<!-- field: use_case -->
#### Use case

rel99.0-uc001-observer-family

<!-- field: integration_target -->
#### Integration target

cd ../.. && mage audit

<!-- field: inputs -->
#### Inputs

```yaml
generalized_from: srd008-observer-agent
```

<!-- field: expected -->
#### Expected

```yaml
topology_discovery: true
monitor_fan_in: true
resource_usage: true
served_fleet_ux: true
```

<!-- field: traces -->
#### Traces

- rel99.0-uc001-observer-family S2
- srd021-observer AC2

<!-- case: Adopters named and realization rows deferred | status: planned -->
### Adopters named and realization rows deferred

srd021-observer names coding-agent, agent-architecture, and the
chatbot-mesh as adopters binding monitor/fleet-observer, and defers the
bindings.catalog and deployment_realizations rows to the adoption
implementation because the actor-role inventory gate rejects a binding
whose profile path does not yet exist.

<!-- field: use_case -->
#### Use case

rel99.0-uc001-observer-family

<!-- field: integration_target -->
#### Integration target

cd ../.. && mage audit

<!-- field: inputs -->
#### Inputs

```yaml
adopters: [coding-agent, agent-architecture, chatbot-mesh]
```

<!-- field: expected -->
#### Expected

```yaml
adopters_named: true
realization_rows_deferred: true
inventory_gate_passes: true
```

<!-- field: traces -->
#### Traces

- rel99.0-uc001-observer-family S3
- srd021-observer AC3
