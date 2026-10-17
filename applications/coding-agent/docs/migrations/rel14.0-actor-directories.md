---
id: rel14.0-coding-agent-actor-directories
title: Coding Agent Actor Directory Normalization
release: "14.0"
status: implemented
issue: GH-1417
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: summary -->
## Summary

Release 14 moves every application-owned coding actor directly below
agents/<actor>/ and makes agents/application.yaml the package and deployment
authority. The shared lifecycle wrapper is explicitly named role-server.

<!-- field: historical_source_layout | form: "tree" -->
## Historical source layout

- status: removed
- path: applications/coding-agent/agents/serving/
- note:
  This path is historical documentation only. No profile, package, Helm
  template, fixture, or runtime lookup resolves it.

<!-- field: current_source_layout | form: "tree" -->
## Current source layout

- shared_wrapper: applications/coding-agent/agents/role-server/
- actors:
  1. applications/coding-agent/agents/planner/
  1. applications/coding-agent/agents/executor/
  1. applications/coding-agent/agents/critic/
  1. applications/coding-agent/agents/applier/

<!-- field: runtime_paths | form: "tree" -->
## Runtime paths

- planner: applications/coding-agent/planner/profile.yaml
- executor: applications/coding-agent/executor/profile.yaml
- critic: applications/coding-agent/critic/profile.yaml
- applier: applications/coding-agent/applier/profile.yaml
- shared_wrapper: applications/coding-agent/role-server/

<!-- field: compatibility | form: "tree" -->
## Compatibility

- source_paths: breaking
- runtime_contract: preserved
- aliases: none
- guidance:
  Source consumers must use the normalized actor paths. Released profile
  packages remain self-contained; no source-tree compatibility alias is
  shipped, preventing duplicate active actor trees.
