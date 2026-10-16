---
id: test-rel00.0
title: Release 00.0 Tests
release: "00.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel00.0-uc001-login

<!-- field: preconditions -->
## Preconditions

- Fixture project directory available.

<!-- field: test_cases -->
## Test cases

<!-- case: TestLogin_ValidCredentials -->
### TestLogin_ValidCredentials

Valid credentials return a session token.

<!-- field: use_case -->
#### Use case

rel00.0-uc001-login

<!-- field: traces -->
#### Traces

- srd001-auth AC1

<!-- case: TestLogin_InvalidCredentials -->
### TestLogin_InvalidCredentials

Invalid credentials return 401.

<!-- field: use_case -->
#### Use case

rel00.0-uc001-login

<!-- field: traces -->
#### Traces

- srd001-auth AC2
