---
id: test-rel05.0-control-plane
title: Test Suite for Release 05.0 Control Plane
release: "05.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: overview -->
## Overview

This suite covers the mesh control plane (srd004 provisioning-workflow-orchestrator, srd005 creator):
a provisioning intent flows from the UX agent through the provisioning-workflow-orchestrator to the
creator's instance operations and back as aggregated status, with the
deployment-API authority held by the creator alone.

The integration:controlPlane tracer is a fake-backed transport-and-authority
proof. It runs the real provisioning-workflow-orchestrator and creator agents against a fake deployment
API and proves the intent flow (chatbot -> provisioning-workflow-orchestrator -> creator -> deployment
API), the request/response sequencing, and the authority boundary (no endpoint or
credential crosses, only the creator calls the deployment API). It does not prove
the behavioral lifecycle: a real corpus-ingest with count verification, a real
request-draining rollout, distinct create/restart/terminate workflows, or a
grounded turn on a newly added source. Those need the real applier and profiles
(the applier behavioral tracer is integration:applier) and are not asserted
here; the trace links below relate the case to the srd004/srd005 acceptance
criteria structurally, not as behavioral proof.

<!-- field: traces -->
## Traces

- [rel05.0-uc001-control-plane-provisioning](../use-cases/rel05.0-uc001-control-plane-provisioning.md)

<!-- field: preconditions -->
## Preconditions

- agent-core is available for the request runtime and REST client transport.
- The chatbot, provisioning-workflow-orchestrator, and creator profiles are deployed with declared clients between them; the creator alone holds the deployment-API credentials.
- A deployment API and the chart's values-to-config co-generation are reachable.

<!-- field: test_cases -->
## Test cases

<!-- case: Provisioning intent tracer across chatbot, provisioning-workflow-orchestrator, and creator | status: implemented -->
### Provisioning intent tracer across chatbot, provisioning-workflow-orchestrator, and creator

Drives a provisioning intent through the real provisioning-workflow-orchestrator and creator agents
against a fake deployment API: the chatbot delegates an add-a-source intent to
the provisioning-workflow-orchestrator through its declared client; the provisioning-workflow-orchestrator sequences an
ingest request then a rollout request to the creator, forwarding the desired
values; the creator drives the deployment API's apply and reads the rollout;
the provisioning-workflow-orchestrator aggregates status. Asserts the authority boundary: the chatbot
and provisioning-workflow-orchestrator submit no endpoint or credential, and only the creator calls
the deployment API. The deployment API is a fake, so this proves transport
sequencing and authority, not the behavioral lifecycle -- a real corpus-ingest
with count verification, a real request-draining rollout, or a grounded turn
on the new source are not asserted (see overview).

<!-- field: expected_note -->
#### Expected note

Not proven by this fake-backed tracer (see overview): real corpus-ingest count verification, a real request-draining rollout, and a grounded turn on the newly added source.

<!-- field: use_case -->
#### Use case

[rel05.0-uc001-control-plane-provisioning](../use-cases/rel05.0-uc001-control-plane-provisioning.md)

<!-- field: integration_target -->
#### Integration target

mage integration:controlPlane

<!-- field: inputs | form: "tree" -->
#### Inputs

- intent: add a directory as a knowledge source
- directories:
  1. one seeded directory
- agents:
  1. chatbot
  1. provisioning-workflow-orchestrator
  1. creator

<!-- field: expected | form: "tree" -->
#### Expected

- chatbot_delegates_intent_no_credential: true
- provisioning_workflow_orchestrator_sequences_ingest_then_rollout: true
- provisioning_workflow_orchestrator_forwards_values_and_requests_rollout: true
- creator_drives_apply_and_reads_rollout_via_deployment_api: true
- only_creator_calls_deployment_api: true

<!-- field: traces -->
#### Traces

- [rel05.0-uc001-control-plane-provisioning](../use-cases/rel05.0-uc001-control-plane-provisioning.md) [S1](../use-cases/rel05.0-uc001-control-plane-provisioning.md#s1)
- [rel05.0-uc001-control-plane-provisioning](../use-cases/rel05.0-uc001-control-plane-provisioning.md) [S2](../use-cases/rel05.0-uc001-control-plane-provisioning.md#s2)
- [rel05.0-uc001-control-plane-provisioning](../use-cases/rel05.0-uc001-control-plane-provisioning.md) [S3](../use-cases/rel05.0-uc001-control-plane-provisioning.md#s3)
- [rel05.0-uc001-control-plane-provisioning](../use-cases/rel05.0-uc001-control-plane-provisioning.md) [S4](../use-cases/rel05.0-uc001-control-plane-provisioning.md#s4)
- [rel05.0-uc001-control-plane-provisioning](../use-cases/rel05.0-uc001-control-plane-provisioning.md) [S5](../use-cases/rel05.0-uc001-control-plane-provisioning.md#s5)
- [srd004-provisioning-workflow-orchestrator](../software-requirements/srd004-provisioning-workflow-orchestrator.md) [AC1](../software-requirements/srd004-provisioning-workflow-orchestrator.md#ac1)
- [srd004-provisioning-workflow-orchestrator](../software-requirements/srd004-provisioning-workflow-orchestrator.md) [AC2](../software-requirements/srd004-provisioning-workflow-orchestrator.md#ac2)
- [srd004-provisioning-workflow-orchestrator](../software-requirements/srd004-provisioning-workflow-orchestrator.md) [AC3](../software-requirements/srd004-provisioning-workflow-orchestrator.md#ac3)
- [srd004-provisioning-workflow-orchestrator](../software-requirements/srd004-provisioning-workflow-orchestrator.md) [AC4](../software-requirements/srd004-provisioning-workflow-orchestrator.md#ac4)
- [srd004-provisioning-workflow-orchestrator](../software-requirements/srd004-provisioning-workflow-orchestrator.md) [AC5](../software-requirements/srd004-provisioning-workflow-orchestrator.md#ac5)
- [srd005-creator](../software-requirements/srd005-creator.md) [AC1](../software-requirements/srd005-creator.md#ac1)
- [srd005-creator](../software-requirements/srd005-creator.md) [AC2](../software-requirements/srd005-creator.md#ac2)
- [srd005-creator](../software-requirements/srd005-creator.md) [AC3](../software-requirements/srd005-creator.md#ac3)
- [srd005-creator](../software-requirements/srd005-creator.md) [AC4](../software-requirements/srd005-creator.md#ac4)
- [srd005-creator](../software-requirements/srd005-creator.md) [AC5](../software-requirements/srd005-creator.md#ac5)
- [srd003-chatbot-deployment](../software-requirements/srd003-chatbot-deployment.md) [AC4](../software-requirements/srd003-chatbot-deployment.md#ac4)

<!-- case: A mismatched embedding family is refused before the ingest leg reads a count | status: implemented -->
### A mismatched embedding family is refused before the ingest leg reads a count

Drives the creator's ingest endpoint twice against a fake Chroma, with the
family this deployment ingests at fixed for the run. An intent expecting
another family answers 422 embedding_family_mismatch and reports the
deployment's family rather than the one it asked for; an intent expecting
the deployed family is admitted past the judge.

The claim is routing, and the collection count read is what makes it
observable. The count sits between the family judge and the child run, so
a refused intent that never reads it is a refused intent that never
reached the child -- which is what the refusal exists to guarantee, since
a family the serving rag unit does not report yields documents excluded
from every turn (srd002 R3.3). Asserting the status alone would pass on a
refusal that arrived after the leg had started toward the child.

The admitted half stops at that same observation rather than asserting a
200. Past the judge the leg runs a real corpus-ingest child, which this
tracer does not stage, so the run ends IngestFailed. Driving only the
refusal would pass against a judge that refused everything, which is why
both halves are here.

<!-- field: expected_note -->
#### Expected note

Not proven here: a corpus-ingest child running, documents written, or a turn over them.

<!-- field: use_case -->
#### Use case

[rel05.0-uc001-control-plane-provisioning](../use-cases/rel05.0-uc001-control-plane-provisioning.md)

<!-- field: integration_target -->
#### Integration target

mage integration:controlPlane

<!-- field: inputs | form: "tree" -->
#### Inputs

- deployed_family: the family fixed in the creator's environment for the run
- mismatched_intent: an ingest intent expecting a different family
- matched_intent: an ingest intent expecting the deployed family

<!-- field: expected | form: "tree" -->
#### Expected

- mismatched_intent_answers_422_embedding_family_mismatch: true
- refusal_reports_the_deployed_family_not_the_requested_one: true
- refused_intent_never_reads_a_collection_count: true
- matched_intent_is_admitted_past_the_judge: true

<!-- field: traces -->
#### Traces

- [rel05.0-uc001-control-plane-provisioning](../use-cases/rel05.0-uc001-control-plane-provisioning.md) [S2](../use-cases/rel05.0-uc001-control-plane-provisioning.md#s2)
- [srd004-provisioning-workflow-orchestrator](../software-requirements/srd004-provisioning-workflow-orchestrator.md) [AC2](../software-requirements/srd004-provisioning-workflow-orchestrator.md#ac2)
- [srd005-creator](../software-requirements/srd005-creator.md) [AC3](../software-requirements/srd005-creator.md#ac3)
