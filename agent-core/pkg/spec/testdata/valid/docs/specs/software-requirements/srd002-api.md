---
id: srd002-api
title: REST API
depends_on:
  - srd_id: srd001-auth
    symbols_used: [SessionToken]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

The system needs an HTTP API.

<!-- field: goals -->
## Goals

- G1: RESTful endpoints for CRUD.

<!-- req: R1 -->
## R1 - Endpoints

<!-- req: R1.1 -->
### R1.1

GET /users returns a list of users.

<!-- req: R1.2 -->
### R1.2

POST /users creates a new user.

<!-- field: non_goals -->
## Non-goals

- Anything beyond this fixture.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

GET /users returns 200 with a JSON array.

<!-- field: traces -->
#### Traces

- R1.1
