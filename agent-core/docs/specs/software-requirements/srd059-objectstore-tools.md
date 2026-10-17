---
id: srd059-objectstore-tools
title: Objectstore Builtin Tool Family
depends_on:
  - srd_id: srd022-builtin-tools
    symbols_used: [factory registrar model, StandardFactoryDeps, declaration placement gates]
  - srd_id: srd035-checkpoint-port
    symbols_used: [receipt-driven undo, core.Reverser]
  - srd_id: srd023-cli-exec-tools
    symbols_used: [boundary category discipline for external side effects]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: problem -->
## Problem

Agents persist artifacts through the filesystem family, which binds every
object to the pod that wrote it: a workspace, an evaluation tree, or a
spool vanishes with its volume, and nothing an agent writes is reachable
from outside its pod. The platform doctrine keeps one local setup and one
cloud setup with differences confined to bootstrap and values overlays, so
object storage cannot arrive as a per-provider tool family: a gcs_write
beside an s3_write would put the provider choice into agent declarations,
which is exactly where it must never live. What is missing is one family
whose words are provider-neutral, whose backend is chosen by a declared
bucket URL flowing through the existing templating pipeline, and whose
local backend needs no cloud at all.

<!-- field: goals -->
## Goals

- G1: One objectstore family with four words, object_read, object_write, object_list, and the internal destructive object_delete_prefix boundary, whose compiled-in drivers are selected by the declared bucket URL scheme.
- G2: Provider switching is configuration. Moving a workload between the kind rig, GCS, and S3 changes a values overlay and never a declaration file or Go code.
- G3: The local rig binds through an explicit endpoint configuration field, and no code path reads an environment variable for endpoint or backend selection.
- G4: object_write is reversible through a prior-object snapshot receipt, so an undo restores what was there or deletes what was not.
- G5: The family follows the dolt model, explicitly selected by profiles, configured through named connections, absent from the default builtin surface.

<!-- req: R1 -->
## R1 - The family and its words

<!-- req: R1.1 -->
### R1.1

The family SHALL declare exactly four words under tools/builtin/objectstore/, object_read taking a key and producing the object as text output, object_write taking text and a key, object_list taking a request prefix and producing a listing, and internal object_delete_prefix taking no request authority and deleting only its declaration-configured prefix. The group SHALL NOT be imported into tools/builtin/all.yaml; a profile selects it by name, as the dolt family is selected.

<!-- req: R1.2 -->
### R1.2

Configuration SHALL be decoded through catalog.DecodeToolConfig as named connections. A connection SHALL carry a required bucket_url whose value flows through the existing declaration templating pipeline, and MAY carry an explicit endpoint override. No other mechanism SHALL influence backend selection.

<!-- req: R1.3 -->
### R1.3

Every word SHALL carry the full authoring-standard contract, and the external_api side effect SHALL place each word in the boundary category.

<!-- req: R2 -->
## R2 - Drivers by declared URL

<!-- req: R2.1 -->
### R2.1

One Go implementation SHALL open the bucket by the URL scheme with compiled-in drivers for gs://, s3://, file://, and mem://. An unknown scheme and an absent bucket_url SHALL fail with named errors at configuration time, not at first use.

<!-- req: R2.2 -->
### R2.2

The endpoint override SHALL be the only way to point the gs:// driver at a local emulator, and the implementation SHALL NOT read STORAGE_EMULATOR_HOST or any other environment variable for endpoint, credential, or backend selection. Cloud credentials are the platform's ambient identity, resolved by the SDK without repository code touching the environment.

<!-- req: R2.3 -->
### R2.3

A profile's declarations SHALL be identical across the kind rig, GCS, and S3 deployments of the same workload; the difference SHALL be carried entirely by the values overlay that supplies bucket_url and, for the local rig, endpoint.

<!-- req: R3 -->
## R3 - Reversibility and side effects

<!-- req: R3.1 -->
### R3.1

object_write SHALL produce a receipt carrying the prior object bytes, or an absent marker when the key did not exist, and its builder SHALL implement core.Reverser so the rollback receipt walk needs no further input.

<!-- req: R3.2 -->
### R3.2

Undo SHALL restore the prior object where the receipt carries bytes and SHALL delete the object where it carries the absent marker. The undo strategy SHALL be registered in the rollback policy under the reversible tier, and checkToolUndoConsistency SHALL accept the declarations.

<!-- req: R3.3 -->
### R3.3

object_read and object_list SHALL declare non-mutating side effects, and the receipt contract's non-mutating kinds SHALL agree with the declared kinds, so the completeness gates hold without baseline entries.

<!-- req: R4 -->
## R4 - Tests and the local proof

<!-- req: R4.1 -->
### R4.1

Unit tests SHALL cover read, write, list, both undo directions, and driver selection against mem:// and file:// buckets, with no container and no network.

<!-- req: R4.2 -->
### R4.2

An integration test SHALL drive the family end to end against fake-gcs-server started with docker run, skipping with the dependency named when docker is absent, with the image pinned by version and digest per ENG01 C2, and with cleanup on every exit path.

<!-- req: R5 -->
## R5 - Declaration-bound destructive purge boundary

<!-- req: R5.1 -->
### R5.1

object_delete_prefix SHALL accept no request field for bucket URL, endpoint, prefix, credential, or application. Those values SHALL be fixed in declaration configuration before the profile runs.

<!-- req: R5.2 -->
### R5.2

The word SHALL list a bounded pre-delete set, refuse without mutation when max_objects is exceeded, delete only the configured prefix, list again to verify absence, and return purged, partial, or refused with matched/deleted/remaining counts and a hash of the selected keys.

<!-- req: R5.3 -->
### R5.3

Every attempt that produces a bounded outcome SHALL append a non-secret fsynced audit record outside the target bucket. The word SHALL be irreversible, internal-only, and absent from the default builtin surface.

<!-- field: non_goals -->
## Non-goals

- The kind-rig fake-gcs add-on; GH-2043 owns it and reuses this family's endpoint shape.
- S3 integration testing and an AWS values overlay; the s3 driver ships compiled in and unit-tested until an AWS environment exists.
- Migrating filesystem usage, the workspace tools, the evaluation tree, or the OTLP spool; that follow-up begins only after this family exists.
- Credentials in YAML. Ambient platform identity only; a credential field is not part of this family.
- Any new substitution or templating mechanism; bucket_url rides the pipeline that exists.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

A profile listing the four words by name in tools.yaml boots, and the words register, dispatch, and validate; object_delete_prefix remains internal and the default builtin surface does not grow.

<!-- field: traces -->
#### Traces

- [R1.1](#r11)
- [R1.2](#r12)
- [R1.3](#r13)

<!-- ac: AC2 -->
### AC2

The backend is chosen solely by the declared bucket_url plus the optional explicit endpoint; the same declarations open mem, file, gs, and s3 buckets as the configuration changes, an unknown scheme fails by name, and no environment variable is read on any path.

<!-- field: traces -->
#### Traces

- [R2.1](#r21)
- [R2.2](#r22)
- [R2.3](#r23)

<!-- ac: AC3 -->
### AC3

A write over an existing object undoes to the prior bytes, a write to an absent key undoes to absence, and the strategy passes the rollback-policy and receipt-contract gates without baseline entries.

<!-- field: traces -->
#### Traces

- [R3.1](#r31)
- [R3.2](#r32)
- [R3.3](#r33)

<!-- ac: AC4 -->
### AC4

Unit tests pass on mem and file buckets with no container, and the fake-gcs integration passes end to end where docker is present and skips naming docker where it is not.

<!-- field: traces -->
#### Traces

- [R4.1](#r41)
- [R4.2](#r42)

<!-- ac: AC5 -->
### AC5

The internal delete-prefix word ignores request-supplied storage selectors, deletes one declared prefix without crossing into a peer prefix, refuses an over-budget set without mutation, verifies remaining objects, writes an external audit record, and is selected only by an explicitly destructive profile.

<!-- field: traces -->
#### Traces

- [R5.1](#r51)
- [R5.2](#r52)
- [R5.3](#r53)

<!-- field: related_documents -->
## Related documents

- cfg-tool-authoring-standard
- srd036-dolt-state-persistence
