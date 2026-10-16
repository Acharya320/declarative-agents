---
id: rel24.0-uc001-objectstore-tools
title: One profile writes objects on kind and on GCS through a values overlay
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

A profile author gives an agent object storage by listing the read/write/list
words, while an explicitly destructive lifecycle profile selects the
internal delete-prefix boundary with declaration-fixed authority. Both use
tools.yaml and declaring a connection whose bucket_url is templated. On the
kind rig the overlay supplies a gs:// bucket with an explicit endpoint at
the local fake-gcs service; on GCS the overlay supplies the real bucket and
no endpoint, and ambient workload identity signs the calls. The agent
writes an object and receives a receipt carrying what was there before; an
undo restores the prior object, or deletes the key the write created. No
declaration file and no Go code changes between the two deployments.

<!-- field: actor -->
## Actor

A profile author giving an agent durable object storage, and the operator moving that agent between the rig and the cloud.

<!-- field: trigger -->
## Trigger

An agent needs to persist an artifact beyond its pod, or a deployed mesh
moves between the kind rig and a cloud project.

<!-- field: flow -->
## Flow

- F1: The profile lists object_read, object_write, and object_list in tools.yaml and declares a named connection with a templated bucket_url (srd059 R1.1, R1.2).
- F2: The kind overlay renders bucket_url gs://agents with endpoint at the in-cluster fake-gcs service; the GCS overlay renders the real bucket and omits endpoint (srd059 R2.2, R2.3).
- F3: The driver opens by scheme; a misspelled scheme or a missing bucket_url fails at configuration with the fault named (srd059 R2.1).
- F4: object_write stores the text and returns a receipt carrying the prior bytes or the absent marker (srd059 R3.1).
- F5: A rollback walks the receipt, restoring the prior object or deleting the created key, with no further input (srd059 R3.2).
- F6: The same profile boots unchanged in both deployments, and only the values overlay differs (srd059 R2.3).
- F7: After exact application confirmation, a model-free lifecycle profile invokes object_delete_prefix with bucket/prefix fixed in declarations; the word bounds, deletes, verifies, and externally audits the selected set (srd059 R5.1, R5.2, R5.3).

<!-- field: touchpoints -->
## Touchpoints

| id | target | reason |
| --- | --- | --- |
| T1 | srd059-objectstore-tools R1.1 | four words, selected explicitly like the dolt family; delete-prefix remains internal |
| T2 | srd059-objectstore-tools R2.1 | one implementation opening the bucket by declared URL scheme |
| T3 | srd059-objectstore-tools R2.2 | the explicit endpoint as the only emulator path, no environment read |
| T4 | srd059-objectstore-tools R2.3 | identical declarations across the rig, GCS, and S3 |
| T5 | srd059-objectstore-tools R3.1 | the prior-object snapshot receipt behind object_write |
| T6 | srd059-objectstore-tools R3.2 | undo restores or deletes through the registered strategy |
| T7 | srd059-objectstore-tools R5.1 | declaration-bound irreversible prefix deletion; R5.2 and R5.3 add bounded verification and external audit |

<!-- field: success_criteria -->
## Success criteria

<!-- sc: S1 -->
### S1

The four words boot from a name-only tools.yaml and the default builtin surface does not grow.

<!-- field: traces -->
#### Traces

- srd059-objectstore-tools AC1

<!-- sc: S2 -->
### S2

Backend selection follows the declared URL and explicit endpoint alone, across mem, file, gs, and s3, with unknown schemes failing by name and no environment read.

<!-- field: traces -->
#### Traces

- srd059-objectstore-tools AC2

<!-- sc: S3 -->
### S3

Write-then-undo restores or deletes correctly in both directions, through the registered strategy.

<!-- field: traces -->
#### Traces

- srd059-objectstore-tools AC3

<!-- sc: S4 -->
### S4

Unit evidence runs with no container; the end-to-end proof runs against fake-gcs under docker and skips naming docker without it.

<!-- field: traces -->
#### Traces

- srd059-objectstore-tools AC4

<!-- sc: S5 -->
### S5

Delete-prefix accepts no request storage authority, refuses an over-budget set without mutation, deletes only its declared prefix, verifies zero remaining, and writes an external audit.

<!-- field: traces -->
#### Traces

- srd059-objectstore-tools AC5

<!-- field: out_of_scope -->
## Out of scope

- The kindrig fake-gcs add-on (GH-2043) and any AWS environment.
- Migrating existing filesystem consumers to this family.
