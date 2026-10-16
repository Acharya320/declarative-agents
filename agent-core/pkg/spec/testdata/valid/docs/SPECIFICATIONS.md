---
id: test-specifications
title: Test Specification Index
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

Three SRDs covering authentication, REST API, and storage.

SRD release assignments.
- 00.0: srd001-auth (authentication), srd002-api (REST API)
- 00.1: srd003-storage (data persistence)

<!-- field: foundation_document_index -->
## Foundation document index

```yaml
[]
```

<!-- field: config_format_index -->
## Config format index

```yaml
[]
```

<!-- field: semantic_model_index -->
## Semantic model index

```yaml
[]
```

<!-- field: coverage_gaps -->
## Coverage gaps

None known.

<!-- field: roadmap_summary -->
## Roadmap summary

```yaml
- version: "00.0"
  name: Foundation
  use_cases_done: 0
  use_cases_total: 1
  status: open
```

<!-- field: srd_index -->
## Srd index

| id | title | summary | path |
| --- | --- | --- | --- |
| srd001-auth | Authentication | Secure authentication. | docs/specs/software-requirements/srd001-auth.md |
| srd002-api | REST API | RESTful endpoints. | docs/specs/software-requirements/srd002-api.md |
| srd003-storage | Storage Layer | Data persistence. | docs/specs/software-requirements/srd003-storage.md |

<!-- field: use_case_index -->
## Use case index

| id | title | release | status | test_suite | path |
| --- | --- | --- | --- | --- | --- |
| rel00.0-uc001-login | Login flow | 00.0 | open | test-rel00.0 | docs/specs/use-cases/rel00.0-uc001-login.md |

<!-- field: test_suite_index -->
## Test suite index

```yaml
- id: test-rel00.0
  title: Release 00.0 Tests
  release: "00.0"
  test_case_count: 2
  traces:
    - rel00.0-uc001-login
  path: docs/specs/test-suites/test-rel00.0.md
```
