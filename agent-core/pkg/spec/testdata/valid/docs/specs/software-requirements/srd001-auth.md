---
id: srd001-auth
title: Authentication
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Users need to log in.

<!-- field: goals -->
## Goals

- G1: Secure authentication.

<!-- req: R1 -->
## R1 - Login

<!-- req: R1.1 -->
### R1.1

Accept username and password.

<!-- req: R1.2 | weight: 2 -->
### R1.2

Validate credentials against the database.

<!-- req: R1.3 -->
### R1.3

Return a session token on success.

<!-- req: R2 -->
## R2 - Logout

<!-- req: R2.1 -->
### R2.1

Invalidate the session token.

<!-- field: non_goals -->
## Non-goals

- Anything beyond this fixture.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

Login returns a token for valid credentials.

<!-- field: traces -->
#### Traces

- [R1.1](#r11)
- [R1.2](#r12)
- [R1.3](#r13)

<!-- ac: AC2 -->
### AC2

Logout invalidates the session.

<!-- field: traces -->
#### Traces

- [R2.1](#r21)
