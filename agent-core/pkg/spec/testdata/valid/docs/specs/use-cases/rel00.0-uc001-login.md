---
id: rel00.0-uc001-login
title: Login flow
test_suite: test-rel00.0
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A user logs in via the REST API.

<!-- field: actor -->
## Actor

End user

<!-- field: trigger -->
## Trigger

POST /login with credentials

<!-- field: flow -->
## Flow

- F1: User submits credentials.
- F2: API validates against auth module.
- F3: Token is returned.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | [srd001-auth](../software-requirements/srd001-auth.md) [R1](../software-requirements/srd001-auth.md#r1---login), [R2](../software-requirements/srd001-auth.md#r2---logout) | authentication flow |
| T2 | [srd002-api](../software-requirements/srd002-api.md) [R1](../software-requirements/srd002-api.md#r1---endpoints) | REST endpoints |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Valid credentials return a token.

<!-- field: traces -->
#### Traces

- [srd001-auth](../software-requirements/srd001-auth.md) [AC1](../software-requirements/srd001-auth.md#ac1)

<!-- sc: S2 -->
### S2

Invalid credentials return 401.

<!-- field: traces -->
#### Traces

- [srd001-auth](../software-requirements/srd001-auth.md) [AC2](../software-requirements/srd001-auth.md#ac2)

<!-- field: out_of_scope -->
## Out of scope

- Password reset flow.
