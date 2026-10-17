---
id: srd001-greet
title: Greeting Function
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

The fixture needs a deterministic greeting implementation.

<!-- field: goals -->
## Goals

- G1: Return the specified greeting for every supplied name.

<!-- req: R1 -->
## R1 - Greeting behavior

<!-- req: R1.1 -->
### R1.1

Hello(name) must return "Hello, <name>!" where <name> is the argument.

<!-- field: non_goals -->
## Non-goals

- Anything beyond this fixture.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

All tests in greet_test.go pass.

<!-- field: traces -->
#### Traces

- R1.1
