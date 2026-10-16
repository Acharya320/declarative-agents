---
id: rel07.1-uc002-lifecycle-profile-boundary
title: Lifecycle Profile Boundary
status: specified
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

The lifecycle profiles prove that checkpoint history inspection, rollback, and
approval are expressed as machine actions and visible signals rather than
hidden CLI branches. A maintainer runs the history, rollback, and approval
programs under testdata/conformance/lifecycle and observes explicit terminal
states from their tool signals.

The proof belongs in applications/catalog because the machine wiring under
testdata/conformance/lifecycle and the tool selections in testdata/conformance/lifecycle/tools.yaml and
testdata/conformance/lifecycle/approval/tools.yaml are profile-owned. agent-core supplies the
lifecycle ports, the checkpoint store, and receipt-backed rollback (opaque
persisted receipts consumed by a receipt-walk Undo, replacing the earlier
workspace-restore model), but it should not encode the approval gate or
history and rollback flows as core-owned CLI behavior.

<!-- field: actor -->
## Actor

Profile maintainer verifying lifecycle history, rollback, and approval behavior.

<!-- field: trigger -->
## Trigger

A maintainer runs each lifecycle program from applications/catalog:

  agent --profile "testdata/conformance/lifecycle/history/profile.yaml"
  agent --profile "testdata/conformance/lifecycle/rollback/profile.yaml"
  agent --profile "testdata/conformance/lifecycle/approval/profile.yaml"

<!-- field: preconditions -->
## Preconditions

- testdata/conformance/lifecycle/history.yaml and testdata/conformance/lifecycle/rollback.yaml select checkpoint_history and checkpoint_rollback from testdata/conformance/lifecycle/tools.yaml.
- testdata/conformance/lifecycle/approval/machine.yaml selects suspend and done from testdata/conformance/lifecycle/approval/tools.yaml.
- A checkpoint ID, a target iteration, and an approval signal are available as inputs.

<!-- field: flow -->
## Flow

- F1: The history machine reads checkpoint state via checkpoint_history and exits Done on ToolDone or Failed on ToolFailed or CommandError.
- F2: The rollback machine rewrites a checkpoint via checkpoint_rollback with explicit checkpoint and target inputs and exits Done or Failed.
- F3: The approval machine suspends before the risky action via suspend and waits in AwaitingApproval.
- F4: A later resume with the Approved signal runs done and reaches the Succeeded terminal state, while the Rejected signal reaches the Rejected terminal state without running the approved path.
- F5: Suspend or resume infrastructure failures route through CommandError to the Failed terminal state.

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd009-lifecycle R1 | History, rollback, and approval are modeled as machine actions and visible signals. |
| T2 | srd009-lifecycle R2 | Lifecycle inputs and the checkpoint_history, checkpoint_rollback, suspend, and done tool families are profile-owned selections over agent-core lifecycle ports. |
| T3 | srd009-lifecycle R3 | Lifecycle outputs and Done, Succeeded, Rejected, or Failed terminal outcomes are checkable. |
| T4 | sm-agent-role-taxonomy lifecycle | The role entry names the lifecycle actor boundary, inputs, selected tool families, outputs, and checks. |
| T5 | test-rel07.1-profile-boundaries lifecycle_profile_boundary | The suite records the history, rollback, and approval machines, tool selections, and expected terminal states. |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

History and rollback run as single declared machine actions that read or rewrite checkpoint state and reach Done or Failed.

<!-- field: traces -->
#### Traces

- srd009-lifecycle R1.1
- srd009-lifecycle R1.2
- srd009-lifecycle R3.1

<!-- sc: S2 -->
### S2

Approval suspends as a machine action and resumes through explicit Approved or Rejected signals rather than hidden CLI branches.

<!-- field: traces -->
#### Traces

- srd009-lifecycle R1.3
- srd009-lifecycle R2.2
- sm-agent-role-taxonomy lifecycle checks

<!-- sc: S3 -->
### S3

Lifecycle terminal outcomes are Done, Succeeded, Rejected, or Failed.

<!-- field: traces -->
#### Traces

- srd009-lifecycle R3.2

<!-- field: out_of_scope -->
## Out of scope

- Reimplementing checkpoint or rollback storage in applications/catalog.
- Defining Go package APIs for the agent-core lifecycle ports or checkpoint store.
- Changing suspend, resume, or rollback semantics during this documentation update.

<!-- field: test_suite_trace -->
## Test suite trace

test-rel07.1-profile-boundaries
