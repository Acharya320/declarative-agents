---
id: test-rel01.0-demo-trace-collection
title: Test Suite for Release 01.0 Demo Trace Collection
release: "01.0"
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

This suite validates that the demo launch path wires the canonical collector
agent for trace collection. Deterministic Go tests inspect the tracing
toggle, collector command construction, and curator OTLP flag injection
without starting processes. Live spool and query evidence remains
dependency-gated.

<!-- field: traces -->
## Traces

- rel01.0-uc001-demo-trace-collection

<!-- field: preconditions -->
## Preconditions

- The application Go module compiles.
- Collector and curator profile constants reference canonical catalog paths.

<!-- field: test_cases -->
## Test cases

<!-- case: Tracing enabled by default | status: implemented -->
### Tracing enabled by default

The deterministic test verifies that tracingEnabled returns true for
empty, "true", and "1" environment values, and false only for
case-insensitive "false".

<!-- field: use_case -->
#### Use case

rel01.0-uc001-demo-trace-collection

<!-- field: id -->
#### Id

tracing_enabled

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestTracingEnabled'

<!-- field: inputs -->
#### Inputs

```yaml
env_values: ["", "true", "1", "false", "FALSE", "False"]
```

<!-- field: expected -->
#### Expected

```yaml
default_enabled: true
disabled_by: "false"
```

<!-- field: traces -->
#### Traces

- rel01.0-uc001-demo-trace-collection S3

<!-- case: Collector agent command construction | status: implemented -->
### Collector agent command construction

The deterministic test constructs the collector command and verifies
the profile path, working directory, and all COLLECTOR_* environment
variables including mode, bind host, receiver address, port
assignments, and spool path. Reservation tests hold the Chatbot collector
ports while allocating Agent Architecture endpoints and prove concurrent
reservations remain disjoint.

<!-- field: use_case -->
#### Use case

rel01.0-uc001-demo-trace-collection

<!-- field: id -->
#### Id

collector_command

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestCollector(Command|Endpoint)|^TestConcurrentCollector'

<!-- field: inputs -->
#### Inputs

```yaml
collector_profile: agents/collector/profile.yaml
mode: spool
ports: four distinct ephemeral loopback reservations
```

<!-- field: expected -->
#### Expected

```yaml
profile_from_catalog: true
working_directory_is_catalog: true
mode_is_spool: true
bind_host: "127.0.0.1"
```

<!-- field: traces -->
#### Traces

- rel01.0-uc001-demo-trace-collection S1
- rel01.0-uc001-demo-trace-collection S4

<!-- case: Curator receives OTLP export flags | status: implemented -->
### Curator receives OTLP export flags

The deterministic test constructs the curator run command plan and
appends the OTLP flags, verifying that --otel-otlp-endpoint and
--otel-service-name appear with the correct values.

<!-- field: use_case -->
#### Use case

rel01.0-uc001-demo-trace-collection

<!-- field: id -->
#### Id

traced_curator_command

<!-- field: go_test -->
#### Go test

go test ./magefiles -run '^TestTracedCurator'

<!-- field: inputs -->
#### Inputs

```yaml
otlp_endpoint: the reserved collector receiver address
service_name: knowledge-manager-curator
```

<!-- field: expected -->
#### Expected

```yaml
otel_flags_present: true
endpoint_matches_collector_receiver: true
```

<!-- field: traces -->
#### Traces

- rel01.0-uc001-demo-trace-collection S2
