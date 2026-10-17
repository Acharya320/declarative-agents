---
id: test-rel17.0-rig-doctor
title: Rig Doctor Root-Cause Analyst Conformance
status: implemented
release: "17.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel17.0-uc001-rig-doctor

<!-- field: preconditions -->
## Preconditions

- The catalog rig-doctor profile closes against agent-core shared tools.
- The read-back gate reads only the declared diagnosis path in the evidence directory.
- conformance/testdata/rig-doctor holds a failing-rollout fixture and an empty fixture.
- A scripted provider serves a fixed diagnosis for the failing-rollout fixture.

<!-- field: test_cases -->
## Test cases

<!-- case: Rig doctor closure | status: implemented -->
### Rig doctor closure

<!-- field: use_case -->
#### Use case

rel17.0-uc001-rig-doctor

<!-- field: id -->
#### ID

rig-doctor-closure

<!-- field: verifies -->
#### Verifies

- S2

<!-- field: traces -->
#### Traces

- rel17.0-uc001-rig-doctor S2

<!-- field: go_test -->
#### Go test

mage validate

<!-- field: assertions -->
#### Assertions

- The rig-doctor profile passes agent-core strict configuration validation.
- The tool manifest holds only list_files, read, find, the two captured-spool queries, the diagnosis write, done, and the loop words.
- No exec, kubectl, or helm declaration is reachable from the profile.

<!-- case: Rig doctor diagnosed | status: implemented -->
### Rig doctor diagnosed

<!-- field: use_case -->
#### Use case

rel17.0-uc001-rig-doctor

<!-- field: id -->
#### ID

rig-doctor-diagnosed

<!-- field: verifies -->
#### Verifies

- S3

<!-- field: traces -->
#### Traces

- rel17.0-uc001-rig-doctor S3

<!-- field: go_test -->
#### Go test

go test ./conformance -run '^TestRigDoctorDiagnosed'

<!-- field: inputs | form: "tree" -->
#### Inputs

- evidence: conformance/testdata/rig-doctor/failing-rollout

<!-- field: assertions -->
#### Assertions

- The run reaches Diagnosed.
- diagnosis.yaml carries probable_cause, confidence, findings, next_command, and unresolved.
- Every evidence cite names a file that exists in the fixture and a line within its length.
- The probable cause names the Deployment whose rollout counts show zero available replicas.

<!-- case: Rig doctor spool queries | status: implemented -->
### Rig doctor spool queries

<!-- field: use_case -->
#### Use case

rel17.0-uc001-rig-doctor

<!-- field: id -->
#### ID

rig-doctor-spool-queries

<!-- field: verifies -->
#### Verifies

- S2
- S3

<!-- field: traces -->
#### Traces

- rel17.0-uc001-rig-doctor S2
- rel17.0-uc001-rig-doctor S3

<!-- field: go_test -->
#### Go test

go test ./conformance -run '^TestRigDoctorQueriesTheCapturedSpool'

<!-- field: inputs | form: "tree" -->
#### Inputs

- evidence: testdata/conformance/rig-doctor/failing-rollout

<!-- field: assertions -->
#### Assertions

- The model dispatches query_list_traces and query_get_trace itself; both produce tool spans.
- The trace id and span names the queries return appear in the next prompt sent to the model.
- The run reaches Diagnosed.

<!-- case: Rig doctor repairs its diagnosis | status: implemented -->
### Rig doctor repairs its diagnosis

<!-- field: use_case -->
#### Use case

rel17.0-uc001-rig-doctor

<!-- field: id -->
#### ID

rig-doctor-repairs-its-diagnosis

<!-- field: verifies -->
#### Verifies

- S3

<!-- field: traces -->
#### Traces

- rel17.0-uc001-rig-doctor S3

<!-- field: go_test -->
#### Go test

go test ./conformance -run '^TestRigDoctorRepairs'

<!-- field: inputs | form: "tree" -->
#### Inputs

- evidence: testdata/conformance/rig-doctor/failing-rollout

<!-- field: assertions -->
#### Assertions

- Undecodable YAML is read back, returned to the model with the reason, and the repaired write reaches Diagnosed.
- A done call with nothing written returns for repair rather than reporting success.

<!-- case: Rig doctor refuses an incomplete diagnosis | status: implemented -->
### Rig doctor refuses an incomplete diagnosis

<!-- field: use_case -->
#### Use case

rel17.0-uc001-rig-doctor

<!-- field: id -->
#### ID

rig-doctor-refuses-an-incomplete-diagnosis

<!-- field: verifies -->
#### Verifies

- S3

<!-- field: traces -->
#### Traces

- rel17.0-uc001-rig-doctor S3

<!-- field: go_test -->
#### Go test

go test ./conformance -run '^TestRigDoctorRefusesADiagnosisWithoutFindings'

<!-- field: assertions -->
#### Assertions

- A decodable diagnosis carrying no findings never reaches Diagnosed.
- A model that never repairs exhausts its budget and reaches the Failed terminal.

<!-- case: Rig doctor retries a failed model call | status: implemented -->
### Rig doctor retries a failed model call

<!-- field: use_case -->
#### Use case

rel17.0-uc001-rig-doctor

<!-- field: id -->
#### ID

rig-doctor-retries-a-failed-model-call

<!-- field: verifies -->
#### Verifies

- S3
- S5

<!-- field: traces -->
#### Traces

- rel17.0-uc001-rig-doctor S3
- rel17.0-uc001-rig-doctor S5

<!-- field: go_test -->
#### Go test

go test ./conformance -run '^TestRigDoctorRetriesAFailedModelCall'

<!-- field: assertions -->
#### Assertions

```yaml
- A provider that fails twice with a mapped status and then answers reaches Diagnosed, driven by ProviderUnavailable.
- The run spends four provider calls: two failures, then the write and the done.
```

<!-- case: Rig doctor does not retry a rejected credential | status: implemented -->
### Rig doctor does not retry a rejected credential

<!-- field: use_case -->
#### Use case

rel17.0-uc001-rig-doctor

<!-- field: id -->
#### ID

rig-doctor-does-not-retry-a-rejected-credential

<!-- field: verifies -->
#### Verifies

- S5

<!-- field: traces -->
#### Traces

- rel17.0-uc001-rig-doctor S5

<!-- field: go_test -->
#### Go test

go test ./conformance -run '^TestRigDoctorDoesNotRetryARejectedCredential'

<!-- field: assertions -->
#### Assertions

- A provider answering 401 costs exactly one attempt.
- The run reaches Failed.

<!-- case: Rig doctor gives up on an unreachable provider | status: implemented -->
### Rig doctor gives up on an unreachable provider

<!-- field: use_case -->
#### Use case

rel17.0-uc001-rig-doctor

<!-- field: id -->
#### ID

rig-doctor-gives-up-on-an-unreachable-provider

<!-- field: verifies -->
#### Verifies

- S5

<!-- field: traces -->
#### Traces

- rel17.0-uc001-rig-doctor S5

<!-- field: go_test -->
#### Go test

go test ./conformance -run '^TestRigDoctorGivesUpOnAnUnreachableProvider'

<!-- field: assertions -->
#### Assertions

- A provider that never answers costs exactly three attempts.
- The run reaches Failed and writes no diagnosis.

<!-- case: Rig doctor no evidence | status: implemented -->
### Rig doctor no evidence

<!-- field: use_case -->
#### Use case

rel17.0-uc001-rig-doctor

<!-- field: id -->
#### ID

rig-doctor-no-evidence

<!-- field: verifies -->
#### Verifies

- S4

<!-- field: traces -->
#### Traces

- rel17.0-uc001-rig-doctor S4

<!-- field: go_test -->
#### Go test

go test ./conformance -run '^TestRigDoctorNoEvidence'

<!-- field: inputs | form: "tree" -->
#### Inputs

- evidence: conformance/testdata/rig-doctor/empty

<!-- field: assertions -->
#### Assertions

- The run reaches NoEvidence.
- No diagnosis.yaml is written.

<!-- case: Rig capture manifest | status: implemented -->
### Rig capture manifest

<!-- field: use_case -->
#### Use case

rel17.0-uc001-rig-doctor

<!-- field: id -->
#### ID

rig-capture-manifest

<!-- field: verifies -->
#### Verifies

- S1

<!-- field: traces -->
#### Traces

- rel17.0-uc001-rig-doctor S1

<!-- field: go_test -->
#### Go test

cd ../.. && go test ./magefiles/kindrig -run '^TestCaptureEvidence'

<!-- field: assertions -->
#### Assertions

- The manifest lists every captured file with its kind.
- A failing kubectl command is recorded as a capture error and later files are still captured.

<!-- case: Rig doctor registration | status: implemented -->
### Rig doctor registration

<!-- field: use_case -->
#### Use case

rel17.0-uc001-rig-doctor

<!-- field: id -->
#### ID

rig-doctor-registration

<!-- field: verifies -->
#### Verifies

- S6

<!-- field: traces -->
#### Traces

- rel17.0-uc001-rig-doctor S6

<!-- field: go_test -->
#### Go test

go test ./conformance -run '^TestFamily'

<!-- field: assertions -->
#### Assertions

- The root-cause-analyst binding appears in the role-realization, taxonomy, and profile-system models.

<!-- case: Rig doctor live | status: implemented -->
### Rig doctor live

<!-- field: use_case -->
#### Use case

rel17.0-uc001-rig-doctor

<!-- field: id -->
#### ID

rig-doctor-live

<!-- field: verifies -->
#### Verifies

- S3
- S5

<!-- field: traces -->
#### Traces

- rel17.0-uc001-rig-doctor S3
- rel17.0-uc001-rig-doctor S5

<!-- field: go_test -->
#### Go test

cd ../chatbot-mesh && mage diagnose helm-smoke

<!-- field: assertions -->
#### Assertions

- With a Deployment broken by an unpullable image tag, the target exits zero and prints the evidence and diagnosis paths.
- The diagnosis names the broken Deployment and cites the ErrImagePull or ImagePullBackOff event.
- The live case is dependency gated on kind, kubectl, and a model provider.
