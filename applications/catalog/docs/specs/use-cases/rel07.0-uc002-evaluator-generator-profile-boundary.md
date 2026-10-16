---
id: rel07.0-uc002-evaluator-generator-profile-boundary
title: Evaluator to Generator Profile Boundary
status: planned
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

This library-owned tracer is retired. Its shell child fabricated the generated
workspace and trace, so it did not prove the configured executor profile.
The replacement belongs to the coding application and runs the canonical
critic and executor with the real built agent binary.

The proof belongs in applications/catalog because suite YAML, subject profile
selection, samples, oracle expectations, and expected artifacts are
profile-owned fixtures. agent-core supplies nested loop execution and the
binary, but it should not encode generator as a concrete evaluator subject.

<!-- field: actor -->
## Actor

Profile maintainer verifying evaluator profile behavior.

<!-- field: trigger -->
## Trigger

A maintainer runs the application-owned critic boundary target:

  cd applications/coding-agent && mage integration:criticGate

<!-- field: preconditions -->
## Preconditions

- agents/critic/profile.yaml and agents/executor/profile.yaml validate.
- testdata/integration/rel07-evaluator-generator contains a minimal suite, sample, oracle, and expected result shape.
- The target can create a temporary output directory for evaluator session results.

<!-- field: flow -->
## Flow

- F1: The target resolves agent-core and applications/catalog roots without developer-specific absolute paths.
- F2: The evaluator profile loads a small suite whose subject profile is generator.
- F3: The evaluator creates an isolated point workspace and invokes generator as the child agent.
- F4: The point records child run status, oracle status, trace evidence, and summary files.
- F5: The session summary records success and failure accounting for the benchmark boundary.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | sm-agent-role-taxonomy evaluator | Evaluator owns suite fixtures, point runs, child profile execution, and session reports. |
| T2 | sm-agent-role-taxonomy generator | Generator is the configured benchmark subject profile. |
| T3 | sm-multi-agent-profile-system evaluator->generator | The dependency direction is benchmark subject execution. |
| T4 | test-rel07.0-profile-integrations evaluator_generator_boundary | The suite records target, fixture paths, and expected artifacts. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The evaluator suite resolves through the demo.yaml catalog_root or the repository checkout.

<!-- field: traces -->
#### Traces

- sm-agent-role-taxonomy evaluator checks

<!-- sc: S2 -->
### S2

The evaluator invokes generator as a child profile rather than relying on core-owned profile names.

<!-- field: traces -->
#### Traces

- sm-multi-agent-profile-system dependency_directions evaluator->generator

<!-- sc: S3 -->
### S3

Point and session output prove child run status and oracle accounting.

<!-- field: traces -->
#### Traces

- sm-agent-role-taxonomy evaluator outputs

<!-- field: out_of_scope -->
## Out of scope

- Running a large model benchmark matrix.
- Changing evaluator scoring semantics.
- Teaching agent-core about the generator profile path.

<!-- field: test_suite_trace -->
## Test suite trace

test-rel07.0-profile-integrations
