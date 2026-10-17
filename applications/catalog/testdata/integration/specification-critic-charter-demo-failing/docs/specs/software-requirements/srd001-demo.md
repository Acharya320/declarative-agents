---
id: srd001-demo
title: Specification Critic Charter Demo
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

The failing demo fixture needs a minimal specification corpus that carries one
error-level violation so the specification critic reaches the Failed terminal.

<!-- field: goals -->
## Goals

- G1: Load a small corpus before charter-only checks execute.

<!-- req: R1 -->
## R1 - Demo corpus

<!-- req: R1.1 -->
### R1.1

Provide one SRD file for graph construction.

<!-- req: R1.2 -->
### R1.2

Provide an uncovered requirement item so a charter error fires.

<!-- field: non_goals -->
## Non-goals

- Anything beyond this fixture.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

The fixture can be loaded before generic charter checks run.

<!-- field: traces -->
#### Traces

- R1.1
