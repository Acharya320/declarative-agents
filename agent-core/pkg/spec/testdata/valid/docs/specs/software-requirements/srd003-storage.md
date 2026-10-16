---
id: srd003-storage
title: Storage Layer
depends_on:
  - srd_id: srd001-auth
    symbols_used: [User]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Data needs to be persisted.

<!-- field: goals -->
## Goals

- G1: Persist data reliably.

<!-- req: R1 -->
## R1 - Database

<!-- req: R1.1 | weight: 3 -->
### R1.1

Connect to PostgreSQL.

<!-- req: R1.2 -->
### R1.2

Run migrations on startup.

<!-- field: non_goals -->
## Non-goals

- N1: NoSQL support.
- N2: Distributed transactions.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

Database connection works.

<!-- field: traces -->
#### Traces

- R1.1
