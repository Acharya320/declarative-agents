---
id: go-style
title: Go Style Constitution
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: purpose -->
## Purpose

We use this constitution for Go source, tests, and magefiles in Agent Core.
It mirrors the local Cursor rule and adds repository-specific package,
lifecycle, and configuration guidance.

<!-- field: copyright_header -->
## Copyright header

Every Go file starts with this header before the package declaration:

  // Copyright (c) 2026 Nokia
  // SPDX-License-Identifier: BSD-3-Clause

<!-- field: implementation_strategy -->
## Implementation strategy

Write all source files before running build or test commands. Plan the full
change, edit files in dependency order, then compile and test once. Avoid
compile-fix loops that change one missing import or symbol at a time.

<!-- field: package_layout | form: "tree" -->
## Package layout

- cmd: Entry points and adapter wiring. Keep main packages small.
- internal: Private implementation packages for runtime, model, tools, evaluation, planning, observability, and support.
- pkg: Public API. In this repo, pkg/spec is the public specification parser and validator surface.
- magefiles: Build tooling. Keep targets small and named by concern.
- agents: YAML profiles, machines, tool selections, and boundary configs.
- tools: Shared builtin and exec tool declarations.

<!-- field: interfaces | form: "tree" -->
## Interfaces

- rules:
  1. Define package boundaries with small interfaces when a parent package depends on a sub-package behavior.
  1. Accept interfaces as parameters.
  1. Return structs.
  1. Keep interfaces to one to three methods.
  1. Do not use interface{} or any in parameters, returns, or struct fields except at standard library serialization boundaries.

<!-- field: design_patterns -->
## Design patterns

| name | use |
| --- | --- |
| Strategy | Multiple implementations of the same behavior selected by config. |
| Command | Runtime actions that execute, report a Result, and may undo or compensate. |
| Facade | Thin adapter over a subsystem when callers should not know internal steps. |
| Factory | NewXxx constructors and builtin factory registration. |
| Decorator | Tracing, timing, retries, and other cross-cutting behavior. |
| Builder | Optional construction parameters or core.Builder command creation. |
| Adapter | Wrapping external APIs, subprocesses, telemetry, or model clients. |
| Observer | Channels or callbacks for lifecycle and async notification when needed. |

<!-- field: errors | form: "tree" -->
## Errors

- rules:
  1. Use guard clauses.
  1. "Wrap errors with context using fmt.Errorf(\"doing X: %w\", err)."
  1. Return errors for runtime conditions.
  1. Panic only for programming contract violations.
  1. Do not panic on I/O, user input, YAML data, subprocess failure, or model responses.
  1. Do not silently discard errors except best-effort cleanup with a comment.

<!-- field: naming | form: "tree" -->
## Naming

- exported: PascalCase
- unexported: camelCase
- cli_flags: kebab-case
- binary_constants: binXxx
- factories: NewXxx
- receivers: one or two lowercase letters, consistent for the type
- packages: short domain nouns, not util, common, helpers, or misc

<!-- field: file_organization | form: "tree" -->
## File organization

- file_names:
  - rule: Name files by concern, not by kind.
  - forbidden:
    1. types.go
    1. helpers.go
    1. utils.go
- declaration_order:
  1. constants and package-level variables
  1. types
  1. constructors
  1. exported functions and methods
  1. unexported helpers
  1. main, for cmd packages
- imports:
  - groups:
    1. standard library
    1. external modules
    1. internal packages
  - rule: Separate groups with blank lines and let gofmt or goimports order within groups.

<!-- field: functions | form: "tree" -->
## Functions

- max_lines: 40
- parameter_rule: Group more than three related parameters into a config struct.
- responsibility_rule: Each function does one thing.
- split_points:
  1. validation
  1. parsing
  1. formatting
  1. I/O
  1. branch-specific behavior

<!-- field: files | form: "tree" -->
## Files

- max_production_lines: 500
- split_rule: Split large files by responsibility before adding unrelated behavior.

<!-- field: enforcement | form: "tree" -->
## Enforcement

- gate: internal/gostyle/size_test.go checks both limits across every production Go file in the module (tests, generated, Mage, and vendored code excluded). Files and functions over the limits today are grandfathered in internal/gostyle/size_baseline.txt; the gate fails on any new oversized file or function, and on any stale baseline entry once its code is split. Delete a baseline line as its file or function drops under the limit.

<!-- field: configuration | form: "tree" -->
## Configuration

- rules:
  1. Runtime agent behavior belongs in YAML profiles, machines, tools, and tool declarations.
  1. cmd/agent SHALL take its program shape from argv flags (for example --profile) and from files those flags reference. Do not use os.Getenv for runtime configuration channels such as profile roots, core install prefixes, or workspace defaults.
  1. Do not add mode-specific branches to cmd/agent when a profile, machine, or ToolDef can express the behavior.
  1. Use typed config structs when decoding ToolDef config.
  1. Keep compatibility paths explicit and tested.
  1. The forbidigo gate enforces the process-environment read/write ban in every non-fixture Go module, including standalone magefiles modules. Production and Mage code SHALL NOT call os.Getenv, os.LookupEnv, or os.Setenv; _test.go files are excluded because tests may stage environment for child and re-exec helpers.
  1. "Every annotated contract site carries an inline //nolint:forbidigo naming the contract that forces it. The gate reported nothing until GH-1479 made mage lint runnable, so this list is the complete set as of that run, in three groups. Declaration expansion mirrored outside the runtime (srd013 R5.6/R5.7): internal/support/envexpand/envexpand.go expands ${VAR:-default} in mounted declarations; applications/chatbot-mesh/magefiles/integration_chroma.go mirrors that expansion in its preflight; applications/catalog/conformance/ollama.go mirrors ${OLLAMA_URL:-...} so a live gate probes the endpoint the profile resolves rather than one a flag chose. The srd048 provider credential, now the token_ref of an srd058 chat dialect, which internal/tools/rest/credentials resolves only when the call executes, behind an injectable resolver, because the credential must never enter ToolDef configuration. A third-party convention: applications/chatbot-mesh/magefiles/integration_helm.go reads SSL_CERT_FILE, the OS and Go name for the host CA bundle. Values a parent sets on a child it owns: magefiles/uidist.go reads the release-gate variable magefiles/tag.go sets on the gate it launches."
  1. os.Environ is allowed for child-process inheritance. A caller may add environment entries on an exec.Cmd when a declared or third-party child contract requires them, including PATH prepends, KUBECONFIG, DOLT_ROOT_PATH, CGO_ENABLED, GOOS, DOCKER_BUILDKIT, and BD_NON_INTERACTIVE; those values do not authorize reading process environment as an application configuration channel. Reading one of those values back is sanctioned only where it is the child half of a contract a parent in this repository owns, or a third-party convention, and only with the annotation above.

<!-- field: tool_language | form: "tree" -->
## Tool language

- rules:
  1. MachineSpec owns workflow sequencing.
  1. ToolDef owns the contract for one vocabulary word.
  1. Boundary tools must declare side effects, emitted signals, and undo or compensation.
  1. New workflow behavior should prefer machine.yaml transitions over compound Go tools.
  1. Deployment-plane authority lives in the executor agent (srd006), a declarative profile that binds helm and kubectl as exec words behind an authenticated apply API. The creator agent (srd005) drives that API; the executor alone wields the cluster credentials and helm/kubectl binaries, keeping that authority out of the general declarative runtime (which carries no credential resolver). Because the executor is itself an agent profile rather than a bespoke Go service, no executable-arm exception to executables-are-agents is needed; the recurring declarative-decomposition audit checks that deployment logic stays in the executor profile.
  1. "Sanctioned home for LLM transport (tool-contract pattern, design-patterns 04): internal/model/llm/ is the model-client adapter layer, and the http.Client construction in internal/tools/llm/invoke.go that wires it is exempt from the rest-tool-carries-transport rule. LLM inference is a distinct transport concern -- streaming responses, per-request model options (options.go), and provider-specific request and response shapes -- that the declarative rest tool (fixed operations, no streaming, no dynamic options) does not model. internal/model/llm/ is a deliberate pluggable-provider abstraction (Ollama and Cohere shipped, with further providers behind the same interfaces). invoke_llm's response parsing is already a separate tool, the right shape; only the inference transport lives in the adapter. The recurring declarative-decomposition audit cites this instead of re-flagging invoke.go and the ollama adapter."

<!-- field: testing | form: "tree" -->
## Testing

- rules:
  1. Use table-driven tests with named rows and t.Run.
  1. Call t.Parallel() first in tests that are safe to parallelize.
  1. Do not call t.Parallel() in tests that call os.Chdir or mutate global state.
  1. Call t.Helper() first in every test helper.
  1. Cover contracts and observable behavior, not implementation details.

<!-- field: comments | form: "tree" -->
## Comments

- exported_symbols: Godoc comment starts with the symbol name.
- unexported_symbols: Comment only when the reason is not obvious from names.
- requirement_ids: Reference SRD or requirement IDs when a comment explains a specific contract.

<!-- field: review_checklist -->
## Review checklist

- No new duplicated logic.
- No scattered binary names, paths, URLs, or repeated labels.
- Every error is handled or intentionally discarded with a comment.
- New functions stay at or below 40 lines.
- New production files stay near or below 500 lines.
- Exported symbols have Godoc comments.
- Tests cover observable behavior.
- Profile-driven behavior stays in YAML unless Go code is the required runtime boundary.

<!-- field: references -->
## References

- docs/ARCHITECTURE.md
- docs/specs/semantic-models/tool-language.md
- docs/specs/config-formats/runtime-contract.md
