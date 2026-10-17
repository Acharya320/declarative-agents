---
id: rel07.0-uc001-command-state-threading
title: A REST word selects a non-adjacent prior word's output by label and resumes mid-chain
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A profile author composes a three-word REST chain: an embed-shaped word, a
transform word, and a query-shaped word. The query word reaches back past the
transform word to the embed word's output through body_source command_state and
$from(...) selectors. The run suspends mid-chain and resumes in a fresh process
with the command-state store rehydrated from persisted output, so the query
word resolves the same label after restart.

<!-- field: actor -->
## Actor

Agent profile author composing a multi-word REST workflow

<!-- field: trigger -->
## Trigger

A workflow needs a later word to consume an earlier, non-adjacent word's output
without threading it through every intervening word.

<!-- field: flow -->
## Flow

- F1: Word 1 (embed-shaped) runs and records its output in the execution log; the command-state store presents that output as a read-only forward view keyed by (run_id, step_index) (srd038 R1).
- F2: Word 2 (transform) runs next; its presence does not change how word 3 addresses word 1.
- F3: Word 3 (query-shaped) declares body_source command_state and selects word 1's fields through input_mapping $from(<label>).<path> selectors, resolving by label with most-recent-wins and no absolute indices (srd038 R2, srd028 R13).
- F4: The store hands word 3 only the fields it selected into declared params; no transport authority or receipt is addressable (srd038 R3).
- F5: The run suspends after word 1 and resumes in a fresh process; the store rehydrates from the srd036 tool_outputs plane so word 3 resolves the same label (srd038 R1).
- F6: If word 1 is a REST command, it translates configured body and header redaction selectors into typed paths through its emitted Result shape while preserving its immediate marker response; core removes those paths before the live Entry becomes selectable and the checkpoint adapter reapplies removal before Save, so word 3 cannot select them live or after resume (srd038 R5).

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd038-command-state-store](../software-requirements/srd038-command-state-store.md) [R1.1](../software-requirements/srd038-command-state-store.md#r11), [R1.2](../software-requirements/srd038-command-state-store.md#r12), [R1.3](../software-requirements/srd038-command-state-store.md#r13), [R1.4](../software-requirements/srd038-command-state-store.md#r14), [R1.5](../software-requirements/srd038-command-state-store.md#r15), [R1.6](../software-requirements/srd038-command-state-store.md#r16), [R1.7](../software-requirements/srd038-command-state-store.md#r17), [R1.8](../software-requirements/srd038-command-state-store.md#r18), [R1.9](../software-requirements/srd038-command-state-store.md#r19) | read-only forward view over the execution log, rehydrated from tool_outputs after restart |
| T2 | [srd038-command-state-store](../software-requirements/srd038-command-state-store.md) [R2.1](../software-requirements/srd038-command-state-store.md#r21), [R2.2](../software-requirements/srd038-command-state-store.md#r22), [R2.3](../software-requirements/srd038-command-state-store.md#r23), [R2.4](../software-requirements/srd038-command-state-store.md#r24), [R2.5](../software-requirements/srd038-command-state-store.md#r25), [R2.6](../software-requirements/srd038-command-state-store.md#r26), [R2.7](../software-requirements/srd038-command-state-store.md#r27), [R2.8](../software-requirements/srd038-command-state-store.md#r28), [R2.9](../software-requirements/srd038-command-state-store.md#r29), [R2.10](../software-requirements/srd038-command-state-store.md#r210), [R2.11](../software-requirements/srd038-command-state-store.md#r211), [R2.12](../software-requirements/srd038-command-state-store.md#r212), [R2.13](../software-requirements/srd038-command-state-store.md#r213), [R2.14](../software-requirements/srd038-command-state-store.md#r214), [R2.15](../software-requirements/srd038-command-state-store.md#r215), [R2.16](../software-requirements/srd038-command-state-store.md#r216), [R2.17](../software-requirements/srd038-command-state-store.md#r217), [R2.18](../software-requirements/srd038-command-state-store.md#r218), [R2.19](../software-requirements/srd038-command-state-store.md#r219), [R2.20](../software-requirements/srd038-command-state-store.md#r220) | $from label addressing, most-recent-wins, run-scoped, no absolute indices |
| T3 | [srd038-command-state-store](../software-requirements/srd038-command-state-store.md) [R3.1](../software-requirements/srd038-command-state-store.md#r31), [R3.2](../software-requirements/srd038-command-state-store.md#r32), [R3.3](../software-requirements/srd038-command-state-store.md#r33), [R3.4](../software-requirements/srd038-command-state-store.md#r34) | declared-param-only authority boundary; no transport, no receipts |
| T4 | [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [R13.1](../software-requirements/srd028-rest-client-tools.md#r131), [R13.2](../software-requirements/srd028-rest-client-tools.md#r132), [R13.3](../software-requirements/srd028-rest-client-tools.md#r133), [R13.4](../software-requirements/srd028-rest-client-tools.md#r134), [R13.5](../software-requirements/srd028-rest-client-tools.md#r135) | REST body_source command_state consumes the store through $from selectors |
| T5 | [srd038-command-state-store](../software-requirements/srd038-command-state-store.md) [R5.1](../software-requirements/srd038-command-state-store.md#r51), [R5.2](../software-requirements/srd038-command-state-store.md#r52), [R5.3](../software-requirements/srd038-command-state-store.md#r53), [R5.4](../software-requirements/srd038-command-state-store.md#r54), [R5.5](../software-requirements/srd038-command-state-store.md#r55), [R5.6](../software-requirements/srd038-command-state-store.md#r56), [R5.7](../software-requirements/srd038-command-state-store.md#r57), [R5.8](../software-requirements/srd038-command-state-store.md#r58), [R5.9](../software-requirements/srd038-command-state-store.md#r59), [R5.10](../software-requirements/srd038-command-state-store.md#r510) | versioned typed redaction before live Entry construction and defense-in-depth enforcement before persistence |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The third word's request body carries the first word's vector as a nested array selected by label past the intervening word.

<!-- field: traces -->
#### Traces

- [srd038-command-state-store](../software-requirements/srd038-command-state-store.md) [AC2](../software-requirements/srd038-command-state-store.md#ac2)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC12](../software-requirements/srd028-rest-client-tools.md#ac12)

<!-- sc: S2 -->
### S2

The store exposes only selected declared params; a selector naming transport authority, an undeclared target, or a receipt is rejected.

<!-- field: traces -->
#### Traces

- [srd038-command-state-store](../software-requirements/srd038-command-state-store.md) [AC3](../software-requirements/srd038-command-state-store.md#ac3)
- [srd028-rest-client-tools](../software-requirements/srd028-rest-client-tools.md) [AC12](../software-requirements/srd028-rest-client-tools.md#ac12)

<!-- sc: S3 -->
### S3

A run suspended after the first word resumes in a fresh process and the third word resolves the same label from the rehydrated store.

<!-- field: traces -->
#### Traces

- [srd038-command-state-store](../software-requirements/srd038-command-state-store.md) [AC1](../software-requirements/srd038-command-state-store.md#ac1)

<!-- sc: S4 -->
### S4

Duplicate labels within the run resolve to the most recent matching step.

<!-- field: traces -->
#### Traces

- [srd038-command-state-store](../software-requirements/srd038-command-state-store.md) [AC2](../software-requirements/srd038-command-state-store.md#ac2)

<!-- sc: S5 -->
### S5

REST body and header selectors produce typed sensitive-output paths without secret values; each path is absent from the live Entry and persisted tool_outputs row, and a later word cannot resolve it before or after resume. Malformed metadata omits the whole output, and an unversioned legacy output returns a typed unavailable error.

<!-- field: traces -->
#### Traces

- [srd038-command-state-store](../software-requirements/srd038-command-state-store.md) [AC5](../software-requirements/srd038-command-state-store.md#ac5)

<!-- field: failure_scenarios -->
## Failure scenarios

| id | scenario | expected |
| --- | --- | --- |
| E1 | Word 3 references a label that no step produced. | A typed unresolved-label error names the label; no empty value is yielded (srd038 R1.5). |
| E2 | After restart the store cannot rehydrate the label's step. | Word 3 fails to resolve and the run reports the error rather than reading an empty value. |
| E3 | A selector names an authentication header. | The selector is rejected before the request is built (srd028 R13.3). |
| E4 | A command supplies an invalid sensitive-output path or marks a path on non-object output. | Core omits the whole output before the Entry becomes live, and Save persists no raw output (srd038 R5.4-R5.5). |
| E5 | A resumed run loads an Entry whose output has no recognized redaction version. | Lifecycle accounting and receipt rollback remain available, but command-state lookup returns a typed unavailable error and exposes no output (srd038 R1.8, R5.6). |

<!-- field: out_of_scope -->
## Out of scope

- Cross-run addressing.
- Non-REST consumers of the store.
- The persistence schema split and revert reaping (covered by rel07.0-uc002).

<!-- field: test_suite_trace -->
## Test suite trace

[test-rel07.0-command-state](../test-suites/test-rel07.0-command-state.md)
