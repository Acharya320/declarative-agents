---
id: rel07.1-uc001-specification-critic-profile-boundary
title: Specification Critic Profile Boundary
status: specified
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: historical_aliases -->
## Historical aliases

- rel07.1-uc001-jurist-profile-boundary

<!-- field: summary -->
## Summary

The specification-critic profile proves a deterministic, LLM-free specification validator
boundary. A maintainer points the profile at a target directory and one or
more charter suites, and the machine runs explicit load, validate, rg search,
finding reduction, and report steps. The run produces a provenance-rich
validation report and a Passed or Failed terminal state.

The proof belongs in applications/catalog because the charter suites, the word
selection in agents/specification-critic/machine.yaml, and the tool selection in
agents/specification-critic/tools.yaml are profile-owned. agent-core supplies the
specification validation tools, but it should not encode the charter suites or
the target directory as core-owned assets.

<!-- field: actor -->
## Actor

Profile maintainer verifying specification-critic charter validation behavior.

<!-- field: trigger -->
## Trigger

```yaml
|
    A maintainer runs the specification-critic profile against a target directory and a charter
    suite:

      agent --profile "agents/specification-critic/profile.yaml" \
            --directory <target-directory>
```

<!-- field: preconditions -->
## Preconditions

- agents/specification-critic/profile.yaml validates and selects agents/specification-critic/machine.yaml and agents/specification-critic/tools.yaml.
- agents/specification-critic/tools.yaml selects load_corpus, validate_specs, jurist_rg, reduce_grep_checks, and format_report with no LLM tool family.
- rg is available on PATH, as required by the existing filesystem find capability.
- agents/specification-critic/suites/demo-charter.yaml supplies a deterministic charter suite over the target directory.

<!-- field: flow -->
## Flow

- F1: The maintainer resolves the specification-critic profile from the demo.yaml catalog_root or the repository checkout without developer-specific absolute paths.
- F2: The machine loads the specification corpus from the target directory in the Loading state via load_corpus.
- F3: The machine validates non-grep charter checks, iterates grep_check plans through the declared jurist_rg executable word, and visibly reduces rg events into findings.
- F4: The machine formats a report in the Reporting state via format_report for both pass and fail outcomes.
- F5: The machine reaches the Passed or Failed terminal state, and findings identify suite, check, severity, message, and any file and line provenance.
- F6: Audit gates select audit-profile.yaml, whose machine declares Go module, package, test-list, and full test-run exec steps, resolves formal claims, reduces pass, fail, skip, and did-not-run events, and reports through format_report. The catalog gate stages a catalog-local executable binding that preserves the same Go JSON inventory and execution streams while running conformance from a unique stable test-binary path; the generic profile remains reusable by other modules.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd005-specification-critic R1 | The specification-critic profile models deterministic, LLM-free, read-only corpus validation. |
| T2 | srd005-specification-critic R2 | Specification-critic inputs, the specification validation tool family, and corpus, graph, consistency, and report responsibilities are profile-owned. |
| T3 | srd005-specification-critic R3 | Specification-critic outputs are a formatted report, provenance-rich findings, and a Passed or Failed status. |
| T4 | sm-agent-role-taxonomy specification-critic | The role entry names the specification-critic actor boundary, inputs, selected tool family, outputs, and checks. |
| T5 | test-rel07.1-profile-boundaries specification_critic_profile_boundary | The suite records the profile, machine word selection, charter suite, and expected report shape. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The specification-critic profile selects load_corpus, validate_specs, jurist_rg, reduce_grep_checks, and format_report and selects no LLM tool family.

<!-- field: traces -->
#### Traces

- srd005-specification-critic R1.1
- srd005-specification-critic R2.2
- sm-agent-role-taxonomy specification-critic checks

<!-- sc: S2 -->
### S2

Charter checks run deterministically and read-only over the target directory across repeated runs, with grep_check search and reduction visible in the machine trace.

<!-- field: traces -->
#### Traces

- srd005-specification-critic R1.2
- srd005-specification-critic R2.3

<!-- sc: S3 -->
### S3

Both pass and fail outcomes produce a formatted report whose findings identify suite, check, severity, message, and available file and line provenance.

<!-- field: traces -->
#### Traces

- srd005-specification-critic R3.1
- srd005-specification-critic R3.2
- srd005-specification-critic R3.3

<!-- sc: S4 -->
### S4

Core, profile, and example audit gates invoke the specification-critic audit profile; the agent CLI has no test-evidence mode and pkg/spec performs only deterministic reduction.

<!-- field: traces -->
#### Traces

- srd005-specification-critic R2.5
- srd005-specification-critic R3.1

<!-- field: out_of_scope -->
## Out of scope

- Adding an LLM boundary or network access to charter validation.
- Defining Go package APIs for the agent-core specification validation tools.
- Changing charter check semantics during this documentation update.
- Renaming the Jurist charter format or jurist_rg tool word, which remain stable capability vocabulary.

<!-- field: test_suite_trace -->
## Test suite trace

test-rel07.1-profile-boundaries
