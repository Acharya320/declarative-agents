---
id: issue-format
title: Issue Format Constitution
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: purpose -->
## Purpose

We use this constitution when writing issues for applications/catalog.
A valid issue names the asset boundary, required files, requirements, and
acceptance criteria. That shape lets agents complete one issue without
changing unrelated runtime code.

<!-- field: articles -->
## Articles

<!-- article: I1 -->
### I1 - Name the deliverable type

Set deliverable_type to documentation, code, asset, fixture, validation,
or migration. When one issue changes docs alongside moved assets, name
the type migration.

<!-- article: I2 -->
### I2 - List required reading

Include required_reading paths. Prefer paths in this repository. Use
../../agent-core paths only when the issue depends on source material or
runtime contracts that still live there.

<!-- article: I3 -->
### I3 - List files in scope

List files or directories that may be changed. Runtime code, builtin tool
implementations, and container build files belong in agent-core issues,
not applications/catalog issues.

<!-- article: I4 -->
### I4 - Use numbered requirements

Requirements use stable IDs such as R1, R2, and R3. Acceptance criteria
use AC1, AC2, and AC3. Each criterion checks a requirement, a migrated
path, a validation command, or a documented runtime usage.

<!-- article: I5 -->
### I5 - Preserve the ownership boundary

Issues must distinguish profile-owned behavior from agent-core-owned
runtime capability. Profiles configure behavior. agent-core executes it.

<!-- field: required_fields -->
## Required fields

- deliverable_type
- required_reading
- files
- requirements
- acceptance_criteria

<!-- field: file_entry_fields -->
## File entry fields

- path
- action

<!-- field: allowed_file_actions -->
## Allowed file actions

- create
- modify
- move
- delete

<!-- field: requirement_numbering -->
## Requirement numbering

```yaml
requirements: R1, R2, R3
acceptance_criteria: AC1, AC2, AC3
```

<!-- field: issue_templates -->
## Issue templates

```yaml
documentation: |
    deliverable_type: documentation
    required_reading:
      - docs/VISION.yaml
      - docs/ARCHITECTURE.yaml
    files:
      - path: docs/path.yaml
        action: create
    requirements:
      - id: R1
        text: Define the document purpose and ownership boundary.
    acceptance_criteria:
      - id: AC1
        text: The document parses as YAML and is indexed when applicable.
migration: |
    deliverable_type: migration
    required_reading:
      - docs/ARCHITECTURE.yaml
      - ../../agent-core/path/to/source
    files:
      - path: agents/name/
        action: create
    requirements:
      - id: R1
        text: Move profile-owned files into the target path.
    acceptance_criteria:
      - id: AC1
        text: The migrated profile loads from the demo.yaml catalog_root without
          developer-specific absolute paths.
```

<!-- field: validation_notes -->
## Validation notes

- Documentation issues run YAML parsing and de-AI checks on touched docs.
- Profile or fixture migration issues run the repository audit command when available.
- Container usage checks must use mounted profile paths such as /profiles.
- Follow-up issues must be created for unrelated validation failures.
