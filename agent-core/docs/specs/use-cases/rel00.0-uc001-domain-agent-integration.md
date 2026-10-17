---
id: rel00.0-uc001-domain-agent-integration
title: Domain agent imports agent-core and runs an agentic loop
test_suite: test-rel00.0
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A domain agent (e.g., Go coding agent, research agent) imports the
agent-core library to build and run an agentic loop. The agent creates
a Registry, registers tools, defines a TransitionTable (programmatically
or from YAML), supplies a Tracer, and calls Loop. The generic engine
drives the state machine, dispatches commands, enforces budget, and
returns a structured RunResult.

<!-- field: actor -->
## Actor

Domain agent developer

<!-- field: trigger -->
## Trigger

The domain agent calls core.Loop(params, ctx) with a configured
LoopParams struct containing the registry, transition table, budget,
and hooks.

<!-- field: flow -->
## Flow

- F1: Domain agent imports internal/runtime/core for types, state machine, dispatch, registry, and loop. Imports internal/observability/tracing for the Tracer interface (srd001, srd007).
- F2: Agent creates a Registry via NewRegistry and registers ToolSpec-Builder pairs for its domain tools (srd004).
- F3: Agent defines a TransitionTable either programmatically or by loading a YAML MachineSpec via ParseMachineSpec and BuildTransitionTable (srd002, srd006).
- F4: Agent sets up telemetry via telemetry.NewRoot with an ExporterConfig, receiving a Trace and shutdown function. Wraps Trace in a TraceAdapter to satisfy the Tracer interface (srd008, srd007).
- F5: Agent can use shared conversation and prompt packages to assemble model input with conversation history, prompt sections, system templates, and registered tool text (srd010, srd011, srd012).
- F6: Agent constructs LoopParams with InitialState, Registry, Table, IsTerminal, Trace, Budget, and optional LoopHooks (srd005).
- F7: Agent calls Loop(params, ctx). Loop freezes the registry, creates OTel spans, and enters the iteration loop (srd005, srd015).
- F8: Each iteration checks context cancellation and budget, steps the state machine, dispatches the command via Dispatch with tracing and timeout (srd003, srd005).
- F9: Loop reaches a terminal state, resolves RunStatus, and returns RunResult with status, iterations, tokens, cost, events, and summary (srd005).
- F10: Agent calls the shutdown function to flush exporters, end the root span, and close files. Offline OTLP protobuf-JSON batches are loaded and relayed by a dedicated MachineSpec (srd008).
- F11: An integration launcher composes file and OTLP export, assigns bounded run identity and a stable service name, and propagates the shared endpoint with a child-specific service name (srd008, srd018).

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd001-core-types](../software-requirements/srd001-core-types.md) [R1.1](../software-requirements/srd001-core-types.md#r11), [R1.2](../software-requirements/srd001-core-types.md#r12), [R1.3](../software-requirements/srd001-core-types.md#r13), [R1.4](../software-requirements/srd001-core-types.md#r14), [R1.5](../software-requirements/srd001-core-types.md#r15), [R2.1](../software-requirements/srd001-core-types.md#r21), [R2.2](../software-requirements/srd001-core-types.md#r22), [R2.3](../software-requirements/srd001-core-types.md#r23), [R3.1](../software-requirements/srd001-core-types.md#r31), [R3.2](../software-requirements/srd001-core-types.md#r32), [R4.1](../software-requirements/srd001-core-types.md#r41), [R4.2](../software-requirements/srd001-core-types.md#r42), [R4.3](../software-requirements/srd001-core-types.md#r43), [R5.1](../software-requirements/srd001-core-types.md#r51), [R5.2](../software-requirements/srd001-core-types.md#r52) | State, Signal, Command, Result, Cost, Builder, CommandResolver types |
| T2 | [srd002-state-machine](../software-requirements/srd002-state-machine.md) [R1.1](../software-requirements/srd002-state-machine.md#r11), [R1.2](../software-requirements/srd002-state-machine.md#r12), [R1.3](../software-requirements/srd002-state-machine.md#r13), [R1.4](../software-requirements/srd002-state-machine.md#r14), [R2.1](../software-requirements/srd002-state-machine.md#r21), [R2.2](../software-requirements/srd002-state-machine.md#r22), [R3.1](../software-requirements/srd002-state-machine.md#r31), [R3.2](../software-requirements/srd002-state-machine.md#r32), [R3.3](../software-requirements/srd002-state-machine.md#r33), [R3.4](../software-requirements/srd002-state-machine.md#r34), [R3.5](../software-requirements/srd002-state-machine.md#r35) | TransitionTable, StateMachine, Step, TerminalFunc |
| T3 | [srd003-command-dispatch](../software-requirements/srd003-command-dispatch.md) [R1.1](../software-requirements/srd003-command-dispatch.md#r11), [R1.2](../software-requirements/srd003-command-dispatch.md#r12), [R1.3](../software-requirements/srd003-command-dispatch.md#r13), [R1.4](../software-requirements/srd003-command-dispatch.md#r14), [R1.5](../software-requirements/srd003-command-dispatch.md#r15), [R2.1](../software-requirements/srd003-command-dispatch.md#r21), [R2.2](../software-requirements/srd003-command-dispatch.md#r22), [R2.3](../software-requirements/srd003-command-dispatch.md#r23), [R2.4](../software-requirements/srd003-command-dispatch.md#r24), [R2.5](../software-requirements/srd003-command-dispatch.md#r25), [R2.6](../software-requirements/srd003-command-dispatch.md#r26), [R3.1](../software-requirements/srd003-command-dispatch.md#r31), [R3.2](../software-requirements/srd003-command-dispatch.md#r32) | Dispatch, SafeExecute, FillDuration, ForceErrorSignal |
| T4 | [srd004-tool-registry](../software-requirements/srd004-tool-registry.md) [R1.1](../software-requirements/srd004-tool-registry.md#r11), [R1.2](../software-requirements/srd004-tool-registry.md#r12), [R2.1](../software-requirements/srd004-tool-registry.md#r21), [R2.2](../software-requirements/srd004-tool-registry.md#r22), [R2.3](../software-requirements/srd004-tool-registry.md#r23), [R2.4](../software-requirements/srd004-tool-registry.md#r24), [R2.5](../software-requirements/srd004-tool-registry.md#r25), [R3.1](../software-requirements/srd004-tool-registry.md#r31), [R3.2](../software-requirements/srd004-tool-registry.md#r32), [R4.1](../software-requirements/srd004-tool-registry.md#r41), [R4.2](../software-requirements/srd004-tool-registry.md#r42), [R5.1](../software-requirements/srd004-tool-registry.md#r51) | ToolSpec, Registry, Manifest, PhaseMatch |
| T5 | [srd005-agentic-loop](../software-requirements/srd005-agentic-loop.md) [R1.1](../software-requirements/srd005-agentic-loop.md#r11), [R1.2](../software-requirements/srd005-agentic-loop.md#r12), [R1.3](../software-requirements/srd005-agentic-loop.md#r13), [R1.4](../software-requirements/srd005-agentic-loop.md#r14), [R2.1](../software-requirements/srd005-agentic-loop.md#r21), [R2.2](../software-requirements/srd005-agentic-loop.md#r22), [R2.3](../software-requirements/srd005-agentic-loop.md#r23), [R2.4](../software-requirements/srd005-agentic-loop.md#r24), [R3.1](../software-requirements/srd005-agentic-loop.md#r31), [R3.2](../software-requirements/srd005-agentic-loop.md#r32), [R3.3](../software-requirements/srd005-agentic-loop.md#r33), [R3.4](../software-requirements/srd005-agentic-loop.md#r34), [R3.5](../software-requirements/srd005-agentic-loop.md#r35), [R3.6](../software-requirements/srd005-agentic-loop.md#r36), [R4.1](../software-requirements/srd005-agentic-loop.md#r41), [R4.2](../software-requirements/srd005-agentic-loop.md#r42), [R4.3](../software-requirements/srd005-agentic-loop.md#r43), [R4.4](../software-requirements/srd005-agentic-loop.md#r44), [R5.1](../software-requirements/srd005-agentic-loop.md#r51), [R5.2](../software-requirements/srd005-agentic-loop.md#r52), [R5.3](../software-requirements/srd005-agentic-loop.md#r53), [R5.4](../software-requirements/srd005-agentic-loop.md#r54), [R5.5](../software-requirements/srd005-agentic-loop.md#r55), [R5.6](../software-requirements/srd005-agentic-loop.md#r56), [R5.7](../software-requirements/srd005-agentic-loop.md#r57), [R5.8](../software-requirements/srd005-agentic-loop.md#r58), [R5.9](../software-requirements/srd005-agentic-loop.md#r59), [R5.10](../software-requirements/srd005-agentic-loop.md#r510), [R5.11](../software-requirements/srd005-agentic-loop.md#r511), [R5.12](../software-requirements/srd005-agentic-loop.md#r512), [R6.1](../software-requirements/srd005-agentic-loop.md#r61) | Budget, RunResult, LoopHooks, LoopParams, Loop, ValidateBuilders |
| T6 | [srd006-yaml-machine-config](../software-requirements/srd006-yaml-machine-config.md) [R1.1](../software-requirements/srd006-yaml-machine-config.md#r11), [R1.2](../software-requirements/srd006-yaml-machine-config.md#r12), [R1.3](../software-requirements/srd006-yaml-machine-config.md#r13), [R1.4](../software-requirements/srd006-yaml-machine-config.md#r14), [R1.5](../software-requirements/srd006-yaml-machine-config.md#r15), [R1.6](../software-requirements/srd006-yaml-machine-config.md#r16), [R1.7](../software-requirements/srd006-yaml-machine-config.md#r17), [R1.8](../software-requirements/srd006-yaml-machine-config.md#r18), [R1.9](../software-requirements/srd006-yaml-machine-config.md#r19), [R2.1](../software-requirements/srd006-yaml-machine-config.md#r21), [R2.2](../software-requirements/srd006-yaml-machine-config.md#r22), [R2.3](../software-requirements/srd006-yaml-machine-config.md#r23), [R2.4](../software-requirements/srd006-yaml-machine-config.md#r24), [R2.5](../software-requirements/srd006-yaml-machine-config.md#r25), [R2.6](../software-requirements/srd006-yaml-machine-config.md#r26), [R2.7](../software-requirements/srd006-yaml-machine-config.md#r27), [R2.8](../software-requirements/srd006-yaml-machine-config.md#r28), [R2.9](../software-requirements/srd006-yaml-machine-config.md#r29), [R2.10](../software-requirements/srd006-yaml-machine-config.md#r210), [R2.11](../software-requirements/srd006-yaml-machine-config.md#r211), [R3.1](../software-requirements/srd006-yaml-machine-config.md#r31), [R3.2](../software-requirements/srd006-yaml-machine-config.md#r32), [R3.3](../software-requirements/srd006-yaml-machine-config.md#r33), [R3.4](../software-requirements/srd006-yaml-machine-config.md#r34), [R3.5](../software-requirements/srd006-yaml-machine-config.md#r35), [R4.1](../software-requirements/srd006-yaml-machine-config.md#r41), [R4.2](../software-requirements/srd006-yaml-machine-config.md#r42), [R4.3](../software-requirements/srd006-yaml-machine-config.md#r43), [R4.4](../software-requirements/srd006-yaml-machine-config.md#r44) | MachineSpec, validation, BuildTransitionTable, LoadMachineSpec |
| T7 | [srd007-tracer-interface](../software-requirements/srd007-tracer-interface.md) [R1.1](../software-requirements/srd007-tracer-interface.md#r11), [R1.2](../software-requirements/srd007-tracer-interface.md#r12), [R1.3](../software-requirements/srd007-tracer-interface.md#r13), [R1.4](../software-requirements/srd007-tracer-interface.md#r14), [R1.5](../software-requirements/srd007-tracer-interface.md#r15), [R1.6](../software-requirements/srd007-tracer-interface.md#r16), [R2.1](../software-requirements/srd007-tracer-interface.md#r21), [R2.2](../software-requirements/srd007-tracer-interface.md#r22) | Tracer interface, minimal dependencies |
| T8 | [srd008-telemetry](../software-requirements/srd008-telemetry.md) [R1.1](../software-requirements/srd008-telemetry.md#r11), [R1.2](../software-requirements/srd008-telemetry.md#r12), [R2.1](../software-requirements/srd008-telemetry.md#r21), [R2.2](../software-requirements/srd008-telemetry.md#r22), [R2.3](../software-requirements/srd008-telemetry.md#r23), [R2.4](../software-requirements/srd008-telemetry.md#r24), [R2.5](../software-requirements/srd008-telemetry.md#r25), [R2.6](../software-requirements/srd008-telemetry.md#r26), [R2.7](../software-requirements/srd008-telemetry.md#r27), [R3.1](../software-requirements/srd008-telemetry.md#r31), [R3.1b](../software-requirements/srd008-telemetry.md#r31b), [R3.2](../software-requirements/srd008-telemetry.md#r32), [R3.3](../software-requirements/srd008-telemetry.md#r33), [R3.4](../software-requirements/srd008-telemetry.md#r34), [R3.5](../software-requirements/srd008-telemetry.md#r35), [R4.1](../software-requirements/srd008-telemetry.md#r41), [R4.2](../software-requirements/srd008-telemetry.md#r42), [R5.1](../software-requirements/srd008-telemetry.md#r51), [R5.2](../software-requirements/srd008-telemetry.md#r52), [R5.3](../software-requirements/srd008-telemetry.md#r53), [R6.1](../software-requirements/srd008-telemetry.md#r61), [R6.2](../software-requirements/srd008-telemetry.md#r62), [R6.3](../software-requirements/srd008-telemetry.md#r63), [R7.1](../software-requirements/srd008-telemetry.md#r71), [R7.2](../software-requirements/srd008-telemetry.md#r72), [R7.3](../software-requirements/srd008-telemetry.md#r73), [R7.4](../software-requirements/srd008-telemetry.md#r74), [R9.1](../software-requirements/srd008-telemetry.md#r91), [R9.2](../software-requirements/srd008-telemetry.md#r92), [R9.3](../software-requirements/srd008-telemetry.md#r93), [R9.4](../software-requirements/srd008-telemetry.md#r94), [R9.5](../software-requirements/srd008-telemetry.md#r95), [R9.6](../software-requirements/srd008-telemetry.md#r96), [R9.7](../software-requirements/srd008-telemetry.md#r97) | ExporterConfig, Trace, NewRoot, file/OTLP exporters, TraceAdapter, declared OTLP replay, retained integration telemetry |
| T9 | [srd010-conversation](../software-requirements/srd010-conversation.md) [R1.1](../software-requirements/srd010-conversation.md#r11), [R1.2](../software-requirements/srd010-conversation.md#r12), [R1.3](../software-requirements/srd010-conversation.md#r13), [R1.4](../software-requirements/srd010-conversation.md#r14), [R2.1](../software-requirements/srd010-conversation.md#r21), [R2.2](../software-requirements/srd010-conversation.md#r22), [R2.3](../software-requirements/srd010-conversation.md#r23), [R4.1](../software-requirements/srd010-conversation.md#r41), [R4.2](../software-requirements/srd010-conversation.md#r42), [R4.3](../software-requirements/srd010-conversation.md#r43), [R4.4](../software-requirements/srd010-conversation.md#r44), [R5.1](../software-requirements/srd010-conversation.md#r51), [R5.2](../software-requirements/srd010-conversation.md#r52), [R5.3](../software-requirements/srd010-conversation.md#r53), [R6.1](../software-requirements/srd010-conversation.md#r61), [R6.2](../software-requirements/srd010-conversation.md#r62), [R6.3](../software-requirements/srd010-conversation.md#r63), [R7.1](../software-requirements/srd010-conversation.md#r71), [R7.2](../software-requirements/srd010-conversation.md#r72), [R7.3](../software-requirements/srd010-conversation.md#r73) | Message, Conversation, append-only history, assembly, system prompt, reset |
| T10 | [srd011-prompt-schema](../software-requirements/srd011-prompt-schema.md) [R1.1](../software-requirements/srd011-prompt-schema.md#r11), [R1.2](../software-requirements/srd011-prompt-schema.md#r12), [R1.3](../software-requirements/srd011-prompt-schema.md#r13), [R1.4](../software-requirements/srd011-prompt-schema.md#r14), [R1.5](../software-requirements/srd011-prompt-schema.md#r15), [R2.1](../software-requirements/srd011-prompt-schema.md#r21), [R2.2](../software-requirements/srd011-prompt-schema.md#r22), [R2.3](../software-requirements/srd011-prompt-schema.md#r23), [R2.4](../software-requirements/srd011-prompt-schema.md#r24), [R3.1](../software-requirements/srd011-prompt-schema.md#r31), [R3.2](../software-requirements/srd011-prompt-schema.md#r32), [R3.3](../software-requirements/srd011-prompt-schema.md#r33), [R3.4](../software-requirements/srd011-prompt-schema.md#r34), [R3.5](../software-requirements/srd011-prompt-schema.md#r35), [R4.1](../software-requirements/srd011-prompt-schema.md#r41), [R4.2](../software-requirements/srd011-prompt-schema.md#r42), [R4.3](../software-requirements/srd011-prompt-schema.md#r43), [R5.1](../software-requirements/srd011-prompt-schema.md#r51), [R5.2](../software-requirements/srd011-prompt-schema.md#r52), [R5.3](../software-requirements/srd011-prompt-schema.md#r53), [R5.4](../software-requirements/srd011-prompt-schema.md#r54), [R6.1](../software-requirements/srd011-prompt-schema.md#r61), [R6.2](../software-requirements/srd011-prompt-schema.md#r62), [R6.3](../software-requirements/srd011-prompt-schema.md#r63), [R6.4](../software-requirements/srd011-prompt-schema.md#r64), [R7.1](../software-requirements/srd011-prompt-schema.md#r71), [R7.2](../software-requirements/srd011-prompt-schema.md#r72) | Prompt, YAML loading, defaults, validation, assembly |
| T11 | [srd012-system-prompt-template](../software-requirements/srd012-system-prompt-template.md) [R1.1](../software-requirements/srd012-system-prompt-template.md#r11), [R1.2](../software-requirements/srd012-system-prompt-template.md#r12), [R1.3](../software-requirements/srd012-system-prompt-template.md#r13), [R2.1](../software-requirements/srd012-system-prompt-template.md#r21), [R2.2](../software-requirements/srd012-system-prompt-template.md#r22), [R2.3](../software-requirements/srd012-system-prompt-template.md#r23), [R2.4](../software-requirements/srd012-system-prompt-template.md#r24), [R3.1](../software-requirements/srd012-system-prompt-template.md#r31), [R4.1](../software-requirements/srd012-system-prompt-template.md#r41), [R4.2](../software-requirements/srd012-system-prompt-template.md#r42), [R4.3](../software-requirements/srd012-system-prompt-template.md#r43), [R4.4](../software-requirements/srd012-system-prompt-template.md#r44), [R4.5](../software-requirements/srd012-system-prompt-template.md#r45), [R4.6](../software-requirements/srd012-system-prompt-template.md#r46), [R5.1](../software-requirements/srd012-system-prompt-template.md#r51), [R5.2](../software-requirements/srd012-system-prompt-template.md#r52), [R6.1](../software-requirements/srd012-system-prompt-template.md#r61), [R6.2](../software-requirements/srd012-system-prompt-template.md#r62), [R6.3](../software-requirements/srd012-system-prompt-template.md#r63), [R6.4](../software-requirements/srd012-system-prompt-template.md#r64), [R6.5](../software-requirements/srd012-system-prompt-template.md#r65), [R6.6](../software-requirements/srd012-system-prompt-template.md#r66), [R6.7](../software-requirements/srd012-system-prompt-template.md#r67) | embedded system template, PromptData, Envelope, rendering, manifest serialization |
| T12 | [srd015-genai-semconv](../software-requirements/srd015-genai-semconv.md) [R1.1](../software-requirements/srd015-genai-semconv.md#r11), [R1.2](../software-requirements/srd015-genai-semconv.md#r12), [R1.3](../software-requirements/srd015-genai-semconv.md#r13), [R1.4](../software-requirements/srd015-genai-semconv.md#r14), [R1.5](../software-requirements/srd015-genai-semconv.md#r15), [R1.6](../software-requirements/srd015-genai-semconv.md#r16), [R1.7](../software-requirements/srd015-genai-semconv.md#r17), [R1.8](../software-requirements/srd015-genai-semconv.md#r18), [R2.1](../software-requirements/srd015-genai-semconv.md#r21), [R3.1](../software-requirements/srd015-genai-semconv.md#r31), [R3.2](../software-requirements/srd015-genai-semconv.md#r32), [R3.3](../software-requirements/srd015-genai-semconv.md#r33), [R4.1](../software-requirements/srd015-genai-semconv.md#r41), [R4.2](../software-requirements/srd015-genai-semconv.md#r42), [R4.3](../software-requirements/srd015-genai-semconv.md#r43), [R4.4](../software-requirements/srd015-genai-semconv.md#r44), [R5.1](../software-requirements/srd015-genai-semconv.md#r51), [R5.2](../software-requirements/srd015-genai-semconv.md#r52), [R5.3](../software-requirements/srd015-genai-semconv.md#r53), [R5.4](../software-requirements/srd015-genai-semconv.md#r54), [R5.5](../software-requirements/srd015-genai-semconv.md#r55), [R5.6](../software-requirements/srd015-genai-semconv.md#r56) | GenAI attribute constants, operation names, providers, span names, attribute builders |
| T13 | [srd018-cli-flag-contract](../software-requirements/srd018-cli-flag-contract.md) [R1.1](../software-requirements/srd018-cli-flag-contract.md#r11), [R1.2](../software-requirements/srd018-cli-flag-contract.md#r12), [R1.3](../software-requirements/srd018-cli-flag-contract.md#r13), [R1.4](../software-requirements/srd018-cli-flag-contract.md#r14), [R1.5](../software-requirements/srd018-cli-flag-contract.md#r15), [R1.6](../software-requirements/srd018-cli-flag-contract.md#r16), [R1.7](../software-requirements/srd018-cli-flag-contract.md#r17), [R1.8](../software-requirements/srd018-cli-flag-contract.md#r18), [R1.9](../software-requirements/srd018-cli-flag-contract.md#r19), [R1.10](../software-requirements/srd018-cli-flag-contract.md#r110), [R2.1](../software-requirements/srd018-cli-flag-contract.md#r21), [R2.2](../software-requirements/srd018-cli-flag-contract.md#r22), [R2.3](../software-requirements/srd018-cli-flag-contract.md#r23), [R2.4](../software-requirements/srd018-cli-flag-contract.md#r24), [R2.5](../software-requirements/srd018-cli-flag-contract.md#r25), [R7.1](../software-requirements/srd018-cli-flag-contract.md#r71), [R7.2](../software-requirements/srd018-cli-flag-contract.md#r72), [R7.3](../software-requirements/srd018-cli-flag-contract.md#r73), [R7.4](../software-requirements/srd018-cli-flag-contract.md#r74) | canonical live telemetry flags and child propagation |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Domain agent can create a Registry, register ToolSpec-Builder pairs, and freeze it.

<!-- field: traces -->
#### Traces

- [srd004-tool-registry](../software-requirements/srd004-tool-registry.md) [AC2](../software-requirements/srd004-tool-registry.md#ac2)
- [srd004-tool-registry](../software-requirements/srd004-tool-registry.md) [AC5](../software-requirements/srd004-tool-registry.md#ac5)

<!-- sc: S2 -->
### S2

Domain agent can build a TransitionTable from a YAML MachineSpec or programmatically.

<!-- field: traces -->
#### Traces

- [srd002-state-machine](../software-requirements/srd002-state-machine.md) [AC1](../software-requirements/srd002-state-machine.md#ac1)
- [srd006-yaml-machine-config](../software-requirements/srd006-yaml-machine-config.md) [AC3](../software-requirements/srd006-yaml-machine-config.md#ac3)

<!-- sc: S3 -->
### S3

Domain agent can call Loop with LoopParams and receive a RunResult.

<!-- field: traces -->
#### Traces

- [srd005-agentic-loop](../software-requirements/srd005-agentic-loop.md) [AC4](../software-requirements/srd005-agentic-loop.md#ac4)
- [srd005-agentic-loop](../software-requirements/srd005-agentic-loop.md) [AC5](../software-requirements/srd005-agentic-loop.md#ac5)

<!-- sc: S4 -->
### S4

Budget enforcement stops the loop when limits are exceeded.

<!-- field: traces -->
#### Traces

- [srd005-agentic-loop](../software-requirements/srd005-agentic-loop.md) [AC1](../software-requirements/srd005-agentic-loop.md#ac1)
- [srd005-agentic-loop](../software-requirements/srd005-agentic-loop.md) [AC5](../software-requirements/srd005-agentic-loop.md#ac5)

<!-- sc: S5 -->
### S5

LoopHooks inject domain-specific behavior (validation, budget, status, result processing).

<!-- field: traces -->
#### Traces

- [srd005-agentic-loop](../software-requirements/srd005-agentic-loop.md) [AC3](../software-requirements/srd005-agentic-loop.md#ac3)

<!-- sc: S6 -->
### S6

Dispatch wraps command execution with tracing, timeout, and panic recovery.

<!-- field: traces -->
#### Traces

- [srd003-command-dispatch](../software-requirements/srd003-command-dispatch.md) [AC1](../software-requirements/srd003-command-dispatch.md#ac1)
- [srd003-command-dispatch](../software-requirements/srd003-command-dispatch.md) [AC2](../software-requirements/srd003-command-dispatch.md#ac2)

<!-- sc: S7 -->
### S7

MachineSpec YAML is validated for referential integrity before building the table.

<!-- field: traces -->
#### Traces

- [srd006-yaml-machine-config](../software-requirements/srd006-yaml-machine-config.md) [AC2](../software-requirements/srd006-yaml-machine-config.md#ac2)
- [srd006-yaml-machine-config](../software-requirements/srd006-yaml-machine-config.md) [AC4](../software-requirements/srd006-yaml-machine-config.md#ac4)

<!-- sc: S8 -->
### S8

Tracer interface is satisfied by TraceAdapter from internal/observability/telemetry.

<!-- field: traces -->
#### Traces

- [srd007-tracer-interface](../software-requirements/srd007-tracer-interface.md) [AC1](../software-requirements/srd007-tracer-interface.md#ac1)
- [srd008-telemetry](../software-requirements/srd008-telemetry.md) [AC8](../software-requirements/srd008-telemetry.md#ac8)

<!-- sc: S9 -->
### S9

NewRoot creates providers from ExporterConfig and returns a Trace and shutdown function.

<!-- field: traces -->
#### Traces

- [srd008-telemetry](../software-requirements/srd008-telemetry.md) [AC1](../software-requirements/srd008-telemetry.md#ac1)
- [srd008-telemetry](../software-requirements/srd008-telemetry.md) [AC4](../software-requirements/srd008-telemetry.md#ac4)
- [srd008-telemetry](../software-requirements/srd008-telemetry.md) [AC5](../software-requirements/srd008-telemetry.md#ac5)

<!-- sc: S10 -->
### S10

Trace supports Push, Event, SetAttributes, Context, and Meter methods.

<!-- field: traces -->
#### Traces

- [srd008-telemetry](../software-requirements/srd008-telemetry.md) [AC2](../software-requirements/srd008-telemetry.md#ac2)
- [srd008-telemetry](../software-requirements/srd008-telemetry.md) [AC3](../software-requirements/srd008-telemetry.md#ac3)

<!-- sc: S11 -->
### S11

File exporter writes atomically via temp file rename. OTLP exporter sends via gRPC.

<!-- field: traces -->
#### Traces

- [srd008-telemetry](../software-requirements/srd008-telemetry.md) [AC6](../software-requirements/srd008-telemetry.md#ac6)
- [srd008-telemetry](../software-requirements/srd008-telemetry.md) [AC7](../software-requirements/srd008-telemetry.md#ac7)

<!-- sc: S12 -->
### S12

A replay MachineSpec loads a trusted OTLP protobuf-JSON batch and sends it unchanged through relay_spans.

<!-- field: traces -->
#### Traces

- [srd008-telemetry](../software-requirements/srd008-telemetry.md) [AC9](../software-requirements/srd008-telemetry.md#ac9)

<!-- sc: S13 -->
### S13

An integration launcher can dual-export attributable file and OTLP evidence, propagate the endpoint and child service identity, and query retained trace and metric backends after the run.

<!-- field: traces -->
#### Traces

- [srd008-telemetry](../software-requirements/srd008-telemetry.md) [AC10](../software-requirements/srd008-telemetry.md#ac10)
- [srd018-cli-flag-contract](../software-requirements/srd018-cli-flag-contract.md) [AC2](../software-requirements/srd018-cli-flag-contract.md#ac2)
- [srd018-cli-flag-contract](../software-requirements/srd018-cli-flag-contract.md) [AC7](../software-requirements/srd018-cli-flag-contract.md#ac7)

<!-- field: out_of_scope -->
## Out of scope

- Domain-specific tool implementations.
- TLS or authenticated OTLP connections.
- CLI integration (domain agent responsibility).
