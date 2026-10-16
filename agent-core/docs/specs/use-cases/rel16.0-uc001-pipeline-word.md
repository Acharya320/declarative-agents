---
id: rel16.0-uc001-pipeline-word
title: A profile folds a marshalling block into one pipeline word
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A profile replaces a render/compose/parse chain of machine states -- the
marshalling around one boundary call -- with a single pipeline word. The
word declares its stages from the existing pure-init vocabulary, each stage
carrying the config its standalone word takes. Stages chain through the
previous-Result contract, the current-value selector inside a stage means
the previous stage's output, and a stage still reads any prior machine step
through $from. The machine emits one declared signal for the whole block
and keeps every state that decides a path.

<!-- field: actor -->
## Actor

A profile author collapsing marshalling geometry around a boundary call

<!-- field: trigger -->
## Trigger

A transition reaches a ToolDef with init pipeline whose config names a
signal and an ordered stage list.

<!-- field: flow -->
## Flow

- F1: The profile loads a pipeline ToolDef whose stages name pure inits (render_each, compose) with their standalone configs.
- F2: Selection sees the pipeline declaration and co-selects the stage inits' factory families.
- F3: Registration builds every stage through its own factory; a stage config the standalone word would refuse fails registration naming the pipeline and the stage.
- F4: Dispatch builds the pipeline command; the engine injects the command-state view, tracer, and monitor recorder once.
- F5: The first stage reads the pipeline's previous Result; each later stage reads its predecessor through the current-value selector; $from selectors resolve against the injected view.
- F6: The final stage's output becomes the word's output under the configured signal; interior signals never reach the machine.
- F7: A failing stage ends the run as CommandError naming the stage; later stages do not run.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd049-pipeline-word R1.1, R1.2, R1.3, R1.4, R1.5, R1.6 | declaration, stage config, eager registration |
| T2 | srd049-pipeline-word R2.1, R2.2, R2.3, R2.4 | closed stage vocabulary and family co-selection |
| T3 | srd049-pipeline-word R3.1, R3.2, R3.3, R3.4 | chaining, selectors, and injected dependencies |
| T4 | srd049-pipeline-word R4.1, R4.2, R4.3 | failure separation |
| T5 | srd038-command-state-store R1.1, R1.2, R1.3, R1.4, R1.5, R1.6, R1.7, R1.8, R1.9, R2.1, R2.2, R2.3, R2.4, R2.5, R2.6, R2.7, R2.8, R2.9, R2.10, R2.11, R2.12, R2.13, R2.14, R2.15, R2.16, R2.17, R2.18, R2.19, R2.20 | selector semantics inside stages |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

A two-stage render/compose pipeline emits its configured signal with the final stage's output, chaining through the current-value selector.

<!-- field: traces -->
#### Traces

- srd049-pipeline-word AC1

<!-- sc: S2 -->
### S2

A stage resolves a $from selector against the same view the standalone word would see.

<!-- field: traces -->
#### Traces

- srd049-pipeline-word AC2

<!-- sc: S3 -->
### S3

A stage init outside the closed set, a refused stage config, an empty stage list, and a missing signal each fail registration naming the fault.

<!-- field: traces -->
#### Traces

- srd049-pipeline-word AC3

<!-- sc: S4 -->
### S4

A failing stage yields CommandError naming the stage while a degraded compose flows through under the configured signal.

<!-- field: traces -->
#### Traces

- srd049-pipeline-word AC4

<!-- sc: S5 -->
### S5

A pipeline declaration alone selects its stage families.

<!-- field: traces -->
#### Traces

- srd049-pipeline-word AC5

<!-- field: out_of_scope -->
## Out of scope

- Branching, boundary calls, and parallel execution inside a pipeline.
- Changes to the standalone stage words or selector grammar.
