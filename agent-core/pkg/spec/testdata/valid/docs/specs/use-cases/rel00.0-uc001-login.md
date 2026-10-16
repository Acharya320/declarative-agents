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
| T1 | srd001-auth R1, R2 | authentication flow |
| T2 | srd002-api R1 | REST endpoints |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

Valid credentials return a token.

<!-- field: traces -->
#### Traces

- srd001-auth AC1

<!-- sc: S2 -->
### S2

Invalid credentials return 401.

<!-- field: traces -->
#### Traces

- srd001-auth AC2

<!-- field: out_of_scope -->
## Out of scope

- Password reset flow.
