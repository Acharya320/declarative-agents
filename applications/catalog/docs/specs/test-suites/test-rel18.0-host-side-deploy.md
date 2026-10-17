---
id: test-rel18.0-host-side-deploy
title: Host-Side Deployer Conformance
status: implemented
release: "18.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- [rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md)

<!-- field: preconditions -->
## Preconditions

- The catalog applier deploy and undeploy profiles close against agent-core shared tools.
- Standalone exec doubles stand in for every Helm and kubectl word, so no case needs a cluster.
- Fixture profiles redeclare one failing word each, so every terminal is reachable offline.
- No case needs a model provider; the deploy machine declares no model word.

<!-- field: test_cases -->
## Test cases

<!-- case: Deploy closure | status: planned -->
### Deploy closure

<!-- field: id -->
#### ID

deploy-closure

<!-- field: use_case -->
#### Use case

[rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md)

<!-- field: go_test -->
#### Go test

mage validate

<!-- field: verifies -->
#### Verifies

- S1

<!-- field: traces -->
#### Traces

- [rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md) [S1](../use-cases/rel18.0-uc001-host-side-deploy.md#s1)

<!-- field: assertions -->
#### Assertions

- The deploy and undeploy profiles pass agent-core strict configuration validation.
- Neither selection holds invoke_llm, parse_response, report_parse_error, or done.
- Both profiles boot with no chart, no cluster, and no provider.

<!-- case: Deploy reaches deployed | status: planned -->
### Deploy reaches deployed

<!-- field: id -->
#### ID

deploy-reaches-deployed

<!-- field: use_case -->
#### Use case

[rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run '^TestApplierDeployReachesDeployed'

<!-- field: verifies -->
#### Verifies

- S1

<!-- field: traces -->
#### Traces

- [rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md) [S1](../use-cases/rel18.0-uc001-host-side-deploy.md#s1)

<!-- field: assertions -->
#### Assertions

- The run reaches the Deployed terminal.
- write_overrides, helm_dry_run, helm_upgrade, and verify_rollout each produce one tool span, in that order.
- The run records no model call.

<!-- case: Deploy rejects a failed dry run | status: planned -->
### Deploy rejects a failed dry run

<!-- field: id -->
#### ID

deploy-rejects-a-failed-dry-run

<!-- field: use_case -->
#### Use case

[rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run '^TestApplierDeployRejectsAFailedDryRun'

<!-- field: verifies -->
#### Verifies

- S5

<!-- field: traces -->
#### Traces

- [rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md) [S5](../use-cases/rel18.0-uc001-host-side-deploy.md#s5)

<!-- field: assertions -->
#### Assertions

- A failing dry run reaches the Rejected terminal.
- helm_upgrade produces no span, so validation precedes mutation.

<!-- case: Deploy rolls back a stalled verify | status: planned -->
### Deploy rolls back a stalled verify

<!-- field: id -->
#### ID

deploy-rolls-back-a-stalled-verify

<!-- field: use_case -->
#### Use case

[rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run '^TestApplierDeployRollsBackAStalledVerify'

<!-- field: verifies -->
#### Verifies

- S5

<!-- field: traces -->
#### Traces

- [rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md) [S5](../use-cases/rel18.0-uc001-host-side-deploy.md#s5)

<!-- field: assertions -->
#### Assertions

- A failing verify_rollout reaches the Compensated terminal through helm_rollback.
- helm_uninstall produces no span, so the machine removes nothing it did not create.

<!-- case: Deploy fails when rollback fails | status: planned -->
### Deploy fails when rollback fails

<!-- field: id -->
#### ID

deploy-fails-when-rollback-fails

<!-- field: use_case -->
#### Use case

[rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run '^TestApplierDeployFailsWhenRollbackFails'

<!-- field: verifies -->
#### Verifies

- S5

<!-- field: traces -->
#### Traces

- [rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md) [S5](../use-cases/rel18.0-uc001-host-side-deploy.md#s5)

<!-- field: assertions -->
#### Assertions

- A failing verify_rollout whose rollback also fails reaches the Failed terminal.
- The run reports a terminal state rather than stopping mid-compensation.

<!-- case: Undeploy removes and reports absent | status: planned -->
### Undeploy removes and reports absent

<!-- field: id -->
#### ID

undeploy-removes-and-reports-absent

<!-- field: use_case -->
#### Use case

[rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run '^TestApplierUndeployRemovesAndReportsAbsent'

<!-- field: verifies -->
#### Verifies

- S7

<!-- field: traces -->
#### Traces

- [rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md) [S7](../use-cases/rel18.0-uc001-host-side-deploy.md#s7)

<!-- field: assertions -->
#### Assertions

- A present release reaches Removed through helm_uninstall.
- A failing helm_history probe reaches Absent, and helm_uninstall produces no span.
- Both terminals carry a succeeded run status, so a teardown of nothing exits zero.

<!-- case: Deploy declares no model word | status: planned -->
### Deploy declares no model word

<!-- field: id -->
#### ID

deploy-declares-no-model-word

<!-- field: use_case -->
#### Use case

[rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md)

<!-- field: go_test -->
#### Go test

go test ./conformance -run '^TestApplierDeployDeclaresNoModelWord'

<!-- field: verifies -->
#### Verifies

- S1

<!-- field: traces -->
#### Traces

- [rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md) [S1](../use-cases/rel18.0-uc001-host-side-deploy.md#s1)

<!-- field: assertions -->
#### Assertions

- Reading the deploy and undeploy selections finds no model word by name.
- The guard fails when a later change adds one.

<!-- case: Deploy rendering resolves every coordinate | status: planned -->
### Deploy rendering resolves every coordinate

<!-- field: id -->
#### ID

deploy-rendering-resolves-every-coordinate

<!-- field: use_case -->
#### Use case

[rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md)

<!-- field: go_test -->
#### Go test

cd ../.. && go test ./magefiles/kindrig -run '^TestRenderDeployDeclarations'

<!-- field: verifies -->
#### Verifies

- S2

<!-- field: traces -->
#### Traces

- [rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md) [S2](../use-cases/rel18.0-uc001-host-side-deploy.md#s2)

<!-- field: assertions -->
#### Assertions

- Every placeholder token is replaced, and each resolved coordinate appears in its expected argv position.
- A coordinate the caller does not supply returns an error naming the token it left behind.

<!-- case: Deploy rendered argv guards | status: planned -->
### Deploy rendered argv guards

<!-- field: id -->
#### ID

deploy-rendered-argv-guards

<!-- field: use_case -->
#### Use case

[rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md)

<!-- field: go_test -->
#### Go test

cd ../.. && go test ./magefiles/kindrig -run '^TestRenderedDeployArgv'

<!-- field: verifies -->
#### Verifies

- S2
- S3

<!-- field: traces -->
#### Traces

- [rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md) [S2](../use-cases/rel18.0-uc001-host-side-deploy.md#s2)
- [rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md) [S3](../use-cases/rel18.0-uc001-host-side-deploy.md#s3)

<!-- field: assertions -->
#### Assertions

- Walking every rendered ToolDef argv finds no placeholder token.
- Walking every rendered ToolDef argv finds no --reuse-values flag.

<!-- case: Deploy gates on a failure terminal | status: planned -->
### Deploy gates on a failure terminal

<!-- field: id -->
#### ID

deploy-gates-on-a-failure-terminal

<!-- field: use_case -->
#### Use case

[rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md)

<!-- field: go_test -->
#### Go test

cd ../.. && go test ./magefiles/kindrig -run '^TestDeployGatesOnAFailureTerminal'

<!-- field: verifies -->
#### Verifies

- S6

<!-- field: traces -->
#### Traces

- [rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md) [S6](../use-cases/rel18.0-uc001-host-side-deploy.md#s6)

<!-- field: assertions -->
#### Assertions

- An agent run that exits non-zero makes the deploy target return an error.
- An agent run that exits zero returns no error, so only a failure terminal fails the build.

<!-- case: Deploy atomic contract live | status: planned -->
### Deploy atomic contract live

<!-- field: id -->
#### ID

deploy-atomic-contract-live

<!-- field: use_case -->
#### Use case

[rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md)

<!-- field: go_test -->
#### Go test

GH-2304 adds the coding-agent deploy target; executable evidence lands with it

<!-- field: verifies -->
#### Verifies

- S4
- S5

<!-- field: traces -->
#### Traces

- [rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md) [S4](../use-cases/rel18.0-uc001-host-side-deploy.md#s4)
- [rel18.0-uc001-host-side-deploy](../use-cases/rel18.0-uc001-host-side-deploy.md) [S5](../use-cases/rel18.0-uc001-host-side-deploy.md#s5)

<!-- field: assertions -->
#### Assertions

- A first install that fails leaves no release behind, because the apply word carries --atomic.
- A failed upgrade against a healthy release leaves the prior revision serving.
- This case needs a running cluster and is the live proof that closes the release.
