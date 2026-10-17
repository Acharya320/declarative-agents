---
id: srd023-rig-doctor
title: Rig Doctor Root-Cause Analyst
depends_on:
  - srd_id: srd001-agent-functional-blocks
    symbols_used: [R1, R2, R3]
  - srd_id: srd020-collector
    symbols_used: [R2]
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: context -->
## Context

A failed kind scenario leaves an evidence directory behind: kind logs,
namespace describe output, and pod logs (eng01 "Evidence and failure
capture"). The host OTLP ingress keeps a trace spool that outlives the run.
Reading that material to find why a rollout never became ready is
correlation work across events, logs, rollout counts, and traces. The
rig-doctor is the catalog realization of the root-cause-analyst sub-role:
it reads one captured evidence directory and writes a structured diagnosis
into it. The repository-root rig capture and each application's
`mage diagnose` target produce the directory; this family consumes it.

<!-- field: design_notes -->
## Design notes

invoke_llm reports every provider failure as CommandError (srd009 R4.3), so
a hung endpoint and a misconfigured one reach the machine as one signal and
it cannot retry the first without retrying the second. The vocabulary for
the distinction already exists: the provider-library words emit
ProviderUnavailable, ProviderThrottled, and ProviderUnauthorized, and the
rerank and embed stages route them separately (srd058 R4.1). invoke_llm is
the outlier, and teaching it those signals is the fix that would serve
every model-driven family rather than this one.

That is now what happens: the word emits the taxonomy when a ToolDef opts in
(srd058 R2.4), and this family opts in. The retry chain that GH-2297 drove
from CommandError is driven by ProviderUnavailable and ProviderThrottled,
and a rejected credential no longer costs three attempts to report.

<!-- field: problem -->
## Problem

Today a person reconstructs the failure by hand from kubectl events,
rollout status, pod logs, and the NDJSON trace spool. The rig captures no
events, records no manifest of what it captured, and nothing ranks probable
causes or cites the lines that support them. A diagnosis that the reader
cannot verify quickly costs more time than reading the logs directly.

<!-- field: goals -->
## Goals

- G1: Diagnose a failed rig from a captured snapshot, never from live cluster access.
- G2: Cite evidence for every claim so a person verifies the diagnosis in seconds.
- G3: Stay read-only by declaration; the only write is the diagnosis file.
- G4: Keep the family portable, independently bootable, and testable without a model provider.

<!-- req: R1 -->
## R1 - Evidence directory contract

<!-- req: R1.1 -->
### R1.1

The input is one directory whose manifest.yaml names the cluster, the captured namespaces, the capture time, the repository revision, every captured file with its kind, and every capture error.

<!-- req: R1.2 -->
### R1.2

File kinds are events, cluster-events, describe, pods, rollout, logs, traces, and kind-logs; a logs file also names its namespace and pod.

<!-- req: R1.3 -->
### R1.3

Per-namespace files follow the existing rig names, namespace-<ns>-<kind>.txt and namespace-<ns>-<pod>-logs.txt; cluster-wide events are cluster-events.txt; the trace spool tail is traces/collector.ndjson; kind export logs stay under kind/.

<!-- req: R1.4 -->
### R1.4

A capture error is recorded in the manifest and does not abort capture of the remaining files.

<!-- req: R2 -->
## R2 - Read-only diagnosis loop

<!-- req: R2.1 -->
### R2.1

The profile runs one model-driven loop whose tool manifest holds only filesystem reads (list_files, read, find), the captured-spool queries (query_list_traces, query_get_trace), the diagnosis write, and done.

<!-- req: R2.5 -->
### R2.5

The spool queries are declared external, because dynamic dispatch resolves only external words and machine exhaustiveness unions the emitted signals of exactly that set; the loop state routes TracesListed and TraceRetrieved back through the model.

<!-- req: R2.2 -->
### R2.2

The profile declares no exec word, no kubectl or helm access, and no network boundary other than the model provider.

<!-- req: R2.3 -->
### R2.3

The profile declares one write word, whose workspace is the evidence directory, and a run writes diagnosis.yaml and nothing else.

<!-- req: R2.4 -->
### R2.4

Every path a word can reach lies inside the evidence directory the run is pointed at, the spool queries through a literal declared traces/collector.ndjson; no path comes from an environment variable.

<!-- req: R3 -->
## R3 - Diagnosis output

<!-- req: R3.1 -->
### R3.1

diagnosis.yaml carries probable_cause, a confidence of high, medium, or low, a findings list, one next_command, and an unresolved list of questions the evidence cannot answer.

<!-- req: R3.2 -->
### R3.2

Every finding carries a claim and at least one evidence cite naming a file path relative to the evidence directory and, where the file is line-oriented, a line number.

<!-- req: R3.3 -->
### R3.3

next_command is a single read-only command a person runs to confirm or refute the probable cause.

<!-- req: R4 -->
## R4 - Terminal states

<!-- req: R4.1 -->
### R4.1

A run reaches Diagnosed only after the diagnosis it wrote is read back, decodes as YAML, names a probable cause, and carries at least one finding.

<!-- req: R4.4 -->
### R4.4

A diagnosis that is missing, undecodable, or incomplete returns to the model with the reason so it repairs and calls done again; the gate reads back only the declared diagnosis path, never a file the model names.

<!-- req: R4.2 -->
### R4.2

An evidence directory that is empty, or that cannot be listed, reaches NoEvidence before any model call and writes no diagnosis.

<!-- req: R4.3 -->
### R4.3

Budget exhaustion or an infrastructure error reaches Failed, from the read-back gate as much as from the loop, so a run that stops mid-gate still reports a terminal state.

<!-- req: R4.5 -->
### R4.5

A model call that reports ProviderUnavailable or ProviderThrottled is retried twice before the run gives up, so a busy provider does not discard evidence the run already paid to read; a third failure reaches Failed and writes no diagnosis. ProviderUnauthorized and CommandError reach Failed at once, because a rejected credential and an infrastructure fault do not improve on a second ask.

<!-- req: R5 -->
## R5 - Portability and registration

<!-- req: R5.1 -->
### R5.1

The profile closes every non-runtime reference and boots independently of any application chart or cluster.

<!-- req: R5.2 -->
### R5.2

The family binds the root-cause-analyst sub-role (parent analyst) in the role-realization, taxonomy, and profile-system semantic models.

<!-- req: R5.3 -->
### R5.3

Conformance exercises the machine against a checked-in fixture evidence directory with a scripted provider, so no test depends on a live model or cluster, and checks the written diagnosis against the R3 schema.

<!-- field: non_goals -->
## Non-goals

- Remediate, restart, roll back, or otherwise change any cluster, release, or file outside diagnosis.yaml.
- Query the cluster, the Kubernetes API, or the host ingress live.
- Capture evidence; the repository-root rig and the application diagnose targets own capture.
- Gate a build or release; a diagnosis reports and never fails its caller.

<!-- field: acceptance_criteria -->
## Acceptance criteria

<!-- ac: AC1 -->
### AC1

Catalog boot smoke validates the standalone profile, and conformance checks that its tool manifest holds only the R2.1 words and no kubectl or helm access.

<!-- field: traces -->
#### Traces

- R2.1
- R2.2
- R2.4
- R2.5
- R5.1

<!-- ac: AC2 -->
### AC2

A fixture evidence directory with a failing rollout drives the scripted-provider run to Diagnosed, and the written diagnosis.yaml satisfies the R3 schema with resolvable evidence cites.

<!-- field: traces -->
#### Traces

- R1.1
- R1.2
- R1.3
- R2.3
- R3.1
- R3.2
- R3.3
- R4.1
- R5.3

<!-- ac: AC3 -->
### AC3

An empty evidence directory reaches NoEvidence without calling the model and without writing a diagnosis.

<!-- field: traces -->
#### Traces

- R4.2
- R4.3
- R2.3

<!-- ac: AC4 -->
### AC4

The family inventory test finds the root-cause-analyst binding in all three semantic models.

<!-- field: traces -->
#### Traces

- R5.2

<!-- ac: AC8 -->
### AC8

A provider that fails twice and then answers reaches Diagnosed; one that never answers costs exactly three attempts and reaches Failed with no diagnosis written; one that rejects the credential costs one.

<!-- field: traces -->
#### Traces

- R4.3
- R4.5

<!-- ac: AC7 -->
### AC7

A scripted run that writes undecodable YAML, and one that calls done having written nothing, each repair and reach Diagnosed; a model that never repairs exhausts its budget and reaches Failed rather than stopping mid-gate.

<!-- field: traces -->
#### Traces

- R4.1
- R4.3
- R4.4

<!-- ac: AC6 -->
### AC6

A scripted run dispatches both spool queries by model choice and the spans they return reach the model in the next prompt.

<!-- field: traces -->
#### Traces

- R2.1
- R2.4
- R2.5

<!-- ac: AC5 -->
### AC5

Rig capture writes a manifest that lists every captured file and records capture errors without aborting.

<!-- field: traces -->
#### Traces

- R1.1
- R1.4

<!-- field: platform_dependencies | form: "typed table" -->
## Platform dependencies

| srd_id | provider | symbols_used |
| --- | --- | --- |
| srd003-application-consistency | applications | [R2, R4, R5, R7] |
