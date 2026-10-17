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
| T1 | [srd049-pipeline-word](../software-requirements/srd049-pipeline-word.md) [R1.1](../software-requirements/srd049-pipeline-word.md#r11), [R1.2](../software-requirements/srd049-pipeline-word.md#r12), [R1.3](../software-requirements/srd049-pipeline-word.md#r13), [R1.4](../software-requirements/srd049-pipeline-word.md#r14), [R1.5](../software-requirements/srd049-pipeline-word.md#r15), [R1.6](../software-requirements/srd049-pipeline-word.md#r16) | declaration, stage config, eager registration |
| T2 | [srd049-pipeline-word](../software-requirements/srd049-pipeline-word.md) [R2.1](../software-requirements/srd049-pipeline-word.md#r21), [R2.2](../software-requirements/srd049-pipeline-word.md#r22), [R2.3](../software-requirements/srd049-pipeline-word.md#r23), [R2.4](../software-requirements/srd049-pipeline-word.md#r24) | closed stage vocabulary and family co-selection |
| T3 | [srd049-pipeline-word](../software-requirements/srd049-pipeline-word.md) [R3.1](../software-requirements/srd049-pipeline-word.md#r31), [R3.2](../software-requirements/srd049-pipeline-word.md#r32), [R3.3](../software-requirements/srd049-pipeline-word.md#r33), [R3.4](../software-requirements/srd049-pipeline-word.md#r34) | chaining, selectors, and injected dependencies |
| T4 | [srd049-pipeline-word](../software-requirements/srd049-pipeline-word.md) [R4.1](../software-requirements/srd049-pipeline-word.md#r41), [R4.2](../software-requirements/srd049-pipeline-word.md#r42), [R4.3](../software-requirements/srd049-pipeline-word.md#r43) | failure separation |
| T5 | [srd038-command-state-store](../software-requirements/srd038-command-state-store.md) [R1.1](../software-requirements/srd038-command-state-store.md#r11), [R1.2](../software-requirements/srd038-command-state-store.md#r12), [R1.3](../software-requirements/srd038-command-state-store.md#r13), [R1.4](../software-requirements/srd038-command-state-store.md#r14), [R1.5](../software-requirements/srd038-command-state-store.md#r15), [R1.6](../software-requirements/srd038-command-state-store.md#r16), [R1.7](../software-requirements/srd038-command-state-store.md#r17), [R1.8](../software-requirements/srd038-command-state-store.md#r18), [R1.9](../software-requirements/srd038-command-state-store.md#r19), [R2.1](../software-requirements/srd038-command-state-store.md#r21), [R2.2](../software-requirements/srd038-command-state-store.md#r22), [R2.3](../software-requirements/srd038-command-state-store.md#r23), [R2.4](../software-requirements/srd038-command-state-store.md#r24), [R2.5](../software-requirements/srd038-command-state-store.md#r25), [R2.6](../software-requirements/srd038-command-state-store.md#r26), [R2.7](../software-requirements/srd038-command-state-store.md#r27), [R2.8](../software-requirements/srd038-command-state-store.md#r28), [R2.9](../software-requirements/srd038-command-state-store.md#r29), [R2.10](../software-requirements/srd038-command-state-store.md#r210), [R2.11](../software-requirements/srd038-command-state-store.md#r211), [R2.12](../software-requirements/srd038-command-state-store.md#r212), [R2.13](../software-requirements/srd038-command-state-store.md#r213), [R2.14](../software-requirements/srd038-command-state-store.md#r214), [R2.15](../software-requirements/srd038-command-state-store.md#r215), [R2.16](../software-requirements/srd038-command-state-store.md#r216), [R2.17](../software-requirements/srd038-command-state-store.md#r217), [R2.18](../software-requirements/srd038-command-state-store.md#r218), [R2.19](../software-requirements/srd038-command-state-store.md#r219), [R2.20](../software-requirements/srd038-command-state-store.md#r220) | selector semantics inside stages |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

A two-stage render/compose pipeline emits its configured signal with the final stage's output, chaining through the current-value selector.

<!-- field: traces -->
#### Traces

- [srd049-pipeline-word](../software-requirements/srd049-pipeline-word.md) [AC1](../software-requirements/srd049-pipeline-word.md#ac1)

<!-- sc: S2 -->
### S2

A stage resolves a $from selector against the same view the standalone word would see.

<!-- field: traces -->
#### Traces

- [srd049-pipeline-word](../software-requirements/srd049-pipeline-word.md) [AC2](../software-requirements/srd049-pipeline-word.md#ac2)

<!-- sc: S3 -->
### S3

A stage init outside the closed set, a refused stage config, an empty stage list, and a missing signal each fail registration naming the fault.

<!-- field: traces -->
#### Traces

- [srd049-pipeline-word](../software-requirements/srd049-pipeline-word.md) [AC3](../software-requirements/srd049-pipeline-word.md#ac3)

<!-- sc: S4 -->
### S4

A failing stage yields CommandError naming the stage while a degraded compose flows through under the configured signal.

<!-- field: traces -->
#### Traces

- [srd049-pipeline-word](../software-requirements/srd049-pipeline-word.md) [AC4](../software-requirements/srd049-pipeline-word.md#ac4)

<!-- sc: S5 -->
### S5

A pipeline declaration alone selects its stage families.

<!-- field: traces -->
#### Traces

- [srd049-pipeline-word](../software-requirements/srd049-pipeline-word.md) [AC5](../software-requirements/srd049-pipeline-word.md#ac5)

<!-- field: out_of_scope -->
## Out of scope

- Branching, boundary calls, and parallel execution inside a pipeline.
- Changes to the standalone stage words or selector grammar.
