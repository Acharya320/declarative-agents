---
id: test-rel16.0-pipeline-word
title: Test Suite for Release 16.0 Pipeline Word
release: "16.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel16.0-uc001-pipeline-word

<!-- field: preconditions -->
## Preconditions

- srd049-pipeline-word defines the stage vocabulary, chaining, failure, and co-selection contracts.
- Tests run hermetically against in-process registries and views; no network or credential is required.

<!-- field: test_cases -->
## Test cases

<!-- case: A marshalling block runs as one word under the configured signal | status: implemented -->
### A marshalling block runs as one word under the configured signal

Build a two-stage pipeline (render_each then compose) and execute it over
a previous Result carrying chunks. The first stage reads the pipeline's
previous Result through the current-value selector, the second reads the
first, and the machine sees only the configured signal with the final
stage's output.

<!-- field: use_case -->
#### Use case

rel16.0-uc001-pipeline-word

<!-- field: go_test -->
#### Go test

go test ./internal/tools/pipeline -run TestPipelineChainsStagesAndEmitsTheConfiguredSignal

<!-- field: expected -->
#### Expected

```yaml
signal: CandidatesComposed
output: final stage's rendered JSON document
interior_signals_visible: false
```

<!-- field: traces -->
#### Traces

- rel16.0-uc001-pipeline-word S1
- srd049-pipeline-word AC1

<!-- case: A stage reads prior machine steps through the injected view | status: implemented -->
### A stage reads prior machine steps through the injected view

A compose stage mixes a $from selector over a prior machine step with a
current-value selector over the pipeline's previous Result; both resolve
exactly as the standalone word would resolve them.

<!-- field: use_case -->
#### Use case

rel16.0-uc001-pipeline-word

<!-- field: go_test -->
#### Go test

go test ./internal/tools/pipeline -run TestPipelineStagesReadPriorStepsThroughTheInjectedView

<!-- field: expected -->
#### Expected

```yaml
from_selector_resolved: true
current_value_resolved: true
```

<!-- field: traces -->
#### Traces

- rel16.0-uc001-pipeline-word S2
- srd049-pipeline-word AC2
- srd049-pipeline-word AC6

<!-- case: Misconfiguration fails registration naming the pipeline and stage | status: implemented -->
### Misconfiguration fails registration naming the pipeline and stage

A boundary init as a stage, a stage config its standalone factory
refuses, an empty stage list, and a missing signal each fail at
registration with the pipeline word and the stage named; none reaches a
run.

<!-- field: use_case -->
#### Use case

rel16.0-uc001-pipeline-word

<!-- field: go_test -->
#### Go test

go test ./internal/tools/pipeline -run TestPipelineMisconfigurationFailsRegistration

<!-- field: expected -->
#### Expected

```yaml
registration_errors: 4
errors_name_stage: true
```

<!-- field: traces -->
#### Traces

- rel16.0-uc001-pipeline-word S3
- srd049-pipeline-word AC3

<!-- case: A failing stage is a named fault and degradation flows through | status: implemented -->
### A failing stage is a named fault and degradation flows through

A render stage over unparsable input ends the pipeline as CommandError
naming the stage and the compose stage never runs; a compose stage whose
selector resolves empty stays a degraded success under the configured
signal (srd038 R1.5 preserved).

<!-- field: use_case -->
#### Use case

rel16.0-uc001-pipeline-word

<!-- field: go_test -->
#### Go test

go test ./internal/tools/pipeline -run TestPipelineStageFailureNamesTheStageAndDegradationFlowsThrough

<!-- field: expected -->
#### Expected

```yaml
failure_signal: CommandError
failure_names_stage: true
later_stages_ran: false
degraded_signal: configured signal
```

<!-- field: traces -->
#### Traces

- rel16.0-uc001-pipeline-word S4
- srd049-pipeline-word AC4

<!-- case: A pipeline declaration co-selects its stage families | status: implemented -->
### A pipeline declaration co-selects its stage families

SelectedBuiltinInits over a profile declaring only a pipeline word
selects the pipeline init and every stage init it names, so the stage
factory families register without a standalone word.

<!-- field: use_case -->
#### Use case

rel16.0-uc001-pipeline-word

<!-- field: go_test -->
#### Go test

go test ./internal/tools/registry ./internal/tools/pipeline -run 'Test(SelectedBuiltinInitsCoSelectPipelineStages|StandardFactoryCatalogSelectsEveryRegisteredInit)'

<!-- field: expected -->
#### Expected

```yaml
pipeline_selected: true
stage_inits_selected: [render_each, parse_structured]
```

<!-- field: traces -->
#### Traces

- rel16.0-uc001-pipeline-word S5
- srd049-pipeline-word AC5
