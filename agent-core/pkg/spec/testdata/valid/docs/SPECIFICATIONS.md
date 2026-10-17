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

<!-- field: foundation_document_index | form: "scalar" -->
## Foundation document index

[]

<!-- field: config_format_index | form: "scalar" -->
## Config format index

[]

<!-- field: semantic_model_index | form: "scalar" -->
## Semantic model index

[]

<!-- field: coverage_gaps -->
## Coverage gaps

None known.

<!-- field: roadmap_summary | form: "typed table" -->
## Roadmap summary

| version | name | use_cases_done | use_cases_total | status |
| --- | --- | --- | --- | --- |
| "00.0" | Foundation | 0 | 1 | open |

<!-- field: srd_index -->
## SRD index

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

<!-- field: test_suite_index | form: "typed table" -->
## Test suite index

| id | title | release | test_case_count | traces | path |
| --- | --- | --- | --- | --- | --- |
| test-rel00.0 | Release 00.0 Tests | "00.0" | 2 | [rel00.0-uc001-login] | docs/specs/test-suites/test-rel00.0.md |
