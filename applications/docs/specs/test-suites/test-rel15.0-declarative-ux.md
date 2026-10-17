---
id: test-rel15.0-declarative-ux
title: Release 15.0 Declarative UX Test Suite
release: "15.0"
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: test_case_count | form: "scalar" -->
## Test case count

9

<!-- field: overview -->
## Overview

TC1 parses and traces the Release 15.0 documents and is executable now. TC2
through TC9 name the tests each implementation issue under GH-2154 adds; each
case turns implemented when its command exists and passes. Kit cases run
under the ui-kit mage targets; REST and manifest cases run under go test.

<!-- field: traces -->
## Traces

- rel15.0-uc001-declarative-ux

<!-- field: srds -->
## Srds

- srd004-declarative-ux

<!-- field: preconditions -->
## Preconditions

- Run commands from the repository root.
- Python 3 with PyYAML, Go, Mage, Node, and npm are available.
- A planned case is not passing evidence until its command exists and passes.

<!-- field: test_cases -->
## Test cases

<!-- case: Release 15.0 documents parse and trace reciprocally | status: implemented -->
### Release 15.0 documents parse and trace reciprocally

<!-- field: id -->
#### ID

TC1

<!-- field: use_case -->
#### Use case

rel15.0-uc001-declarative-ux

<!-- field: evidence_type -->
#### Evidence type

executable_structural

<!-- field: command -->
#### Command

python3 - <<'PY'
from pathlib import Path
import subprocess
import yaml

root = Path("applications/docs")
paths = {
    "srd": root / "specs/software-requirements/srd004-declarative-ux.md",
    "use_case": root / "specs/use-cases/rel15.0-uc001-declarative-ux.md",
    "suite": root / "specs/test-suites/test-rel15.0-declarative-ux.md",
    "index": root / "SPECIFICATIONS.md",
    "roadmap": root / "road-map.md",
}
def load(path):
    # A corpus document is Markdown; mage specYAML prints the YAML tree
    # its declaration parses it into (GH-2627).
    out = subprocess.run(["mage", "specYAML", str(path)], cwd="agent-core",
                         check=True, capture_output=True, text=True).stdout
    return yaml.safe_load(out)

docs = {name: load(path) for name, path in paths.items()}
required = {
    "srd": {"id", "title", "problem", "goals", "requirements",
            "non_goals", "acceptance_criteria"},
    "use_case": {"id", "title", "summary", "actor", "trigger", "flow",
                 "touchpoints", "success_criteria", "out_of_scope"},
    "suite": {"id", "title", "release", "traces", "preconditions",
              "test_cases"},
}
for name, fields in required.items():
    assert not fields - set(docs[name]), (name, fields - set(docs[name]))

srd, use_case, suite = docs["srd"], docs["use_case"], docs["suite"]
assert use_case["test_suite"] == suite["id"]
assert use_case["id"] in suite["traces"]
assert use_case["id"] in srd["use_case_traces"]
assert suite["id"] in srd["test_suite_traces"]

item_ids = {
    item["id"]
    for group in srd["requirements"].values()
    for item in group["items"]
}
ac_traces = {t for ac in srd["acceptance_criteria"] for t in ac["traces"]}
assert item_ids <= ac_traces, sorted(item_ids - ac_traces)
ac_ids = {ac["id"] for ac in srd["acceptance_criteria"]}
case_traces = {
    t.rsplit(" ", 1)[-1]
    for case in suite["test_cases"]
    for t in case.get("traces", [])
    if t.startswith(srd["id"] + " ")
}
assert ac_ids <= case_traces, sorted(ac_ids - case_traces)
assert len(suite["test_cases"]) == suite["test_case_count"]

index = docs["index"]
assert any(e["id"] == srd["id"] for e in index["srd_index"])
assert any(e["id"] == use_case["id"] for e in index["use_case_index"])
assert any(e["id"] == suite["id"] for e in index["test_suite_index"])
release = next(r for r in docs["roadmap"]["releases"] if r["version"] == "15.0")
assert any(e["id"] == use_case["id"] for e in release["use_cases"])
assert any(e["id"] == suite["id"] for e in release["test_suites"])
print("validated Release 15.0 declarative UX trace closure")
PY

<!-- field: expected | form: "tree" -->
#### Expected

- yaml_parses: true
- design_constitution_fields_present: true
- requirement_to_acceptance_criterion_closure: true
- reciprocal_srd_use_case_suite_traces: true
- specification_and_roadmap_indexes_present: true

<!-- field: traces -->
#### Traces

- rel15.0-uc001-declarative-ux S1
- srd004-declarative-ux AC1

<!-- case: Kit client honours the base URL, proxy convention, and not-deployed result | status: implemented -->
### Kit client honours the base URL, proxy convention, and not-deployed result

<!-- field: id -->
#### ID

TC2

<!-- field: use_case -->
#### Use case

rel15.0-uc001-declarative-ux

<!-- field: evidence_type -->
#### Evidence type

kit_test

<!-- field: implemented_by -->
#### Implemented by

GH-2261

<!-- field: command -->
#### Command

mage uikit:test

<!-- field: expected | form: "tree" -->
#### Expected

- remote_base_url_yields_absolute_urls: true
- proxy_urls_built_by_one_function: true
- proxy_404_returns_not_deployed: true
- trace_backend_taken_from_configuration: true
- no_fetch_or_eventsource_outside_client: true

<!-- field: traces -->
#### Traces

- rel15.0-uc001-declarative-ux S5
- srd004-declarative-ux AC2
- srd004-declarative-ux AC6

<!-- case: Contract fixtures parse through every client | status: implemented -->
### Contract fixtures parse through every client

<!-- field: id -->
#### ID

TC3

<!-- field: use_case -->
#### Use case

rel15.0-uc001-declarative-ux

<!-- field: evidence_type -->
#### Evidence type

kit_test

<!-- field: implemented_by -->
#### Implemented by

GH-2261

<!-- field: command -->
#### Command

mage uikit:test

<!-- field: expected | form: "tree" -->
#### Expected

- one_fixture_per_contract_endpoint: true
- every_client_parses_its_fixture: true
- contract_version_exported: true

<!-- field: traces -->
#### Traces

- rel15.0-uc001-declarative-ux S5
- srd004-declarative-ux AC1
- srd004-declarative-ux AC3

<!-- case: Kit panels export manifests within the contract and render from fixtures | status: implemented -->
### Kit panels export manifests within the contract and render from fixtures

<!-- field: id -->
#### ID

TC4

<!-- field: use_case -->
#### Use case

rel15.0-uc001-declarative-ux

<!-- field: evidence_type -->
#### Evidence type

kit_test

<!-- field: implemented_by -->
#### Implemented by

GH-2263, GH-2264, GH-2265, GH-2266

<!-- field: command -->
#### Command

mage uikit:test

<!-- field: expected | form: "tree" -->
#### Expected

- panel_set_matches_r4_4: true
- required_endpoints_subset_of_contract: true
- no_application_imports: true
- every_panel_renders_from_fixtures: true

<!-- field: traces -->
#### Traces

- rel15.0-uc001-declarative-ux S1
- srd004-declarative-ux AC4

<!-- case: ui.yaml version 2 validates, rejects, and renders from Helm values | status: implemented -->
### ui.yaml version 2 validates, rejects, and renders from Helm values

<!-- field: id -->
#### ID

TC5

<!-- field: use_case -->
#### Use case

rel15.0-uc001-declarative-ux

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2270

<!-- field: command -->
#### Command

(cd magefiles && go test ./uiyaml/)
(cd applications/chatbot-mesh && go test ./magefiles/ -run TestChatbotUIRendersAsValidUIYAML)

<!-- field: expected | form: "tree" -->
#### Expected

- every_current_ui_yaml_valid_unchanged: true
- version_2_fixture_valid: true
- duplicate_panel_route_collision_and_unknown_group_rejected: true
- helm_rendered_version_2_valid: true

<!-- field: traces -->
#### Traces

- rel15.0-uc001-declarative-ux S2
- srd004-declarative-ux AC7

<!-- case: Generic shell composes routes and sidebar from ui.yaml | status: implemented -->
### Generic shell composes routes and sidebar from ui.yaml

<!-- field: id -->
#### ID

TC6

<!-- field: use_case -->
#### Use case

rel15.0-uc001-declarative-ux

<!-- field: evidence_type -->
#### Evidence type

kit_test

<!-- field: implemented_by -->
#### Implemented by

GH-2271, GH-2272

<!-- field: command -->
#### Command

mage uikit:test
(cd applications/chatbot-mesh && go test ./magefiles/...)

<!-- field: expected | form: "tree" -->
#### Expected

- routes_and_sidebar_from_ui_yaml: true
- unregistered_panel_fails_build: true
- deep_link_cases_preserved: true
- chatbot_mesh_has_no_routes_ts: true
- registry_covers_every_panel_id: true

<!-- field: traces -->
#### Traces

- rel15.0-uc001-declarative-ux S2
- srd004-declarative-ux AC5
- srd004-declarative-ux AC7

<!-- case: Kit owns tokens and ships as a file dependency and release tarball | status: implemented -->
### Kit owns tokens and ships as a file dependency and release tarball

<!-- field: id -->
#### ID

TC7

<!-- field: use_case -->
#### Use case

rel15.0-uc001-declarative-ux

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2259, GH-2260

<!-- field: command -->
#### Command

mage uikit:build uikit:test uikit:pack
(cd applications/catalog && go test ./conformance -run DesignTokens)

<!-- field: expected | form: "tree" -->
#### Expected

- tarball_produced: true
- every_token_claiming_ui_imports_kit_tokens: true
- no_other_root_token_block: true

<!-- field: traces -->
#### Traces

- rel15.0-uc001-declarative-ux S4
- srd004-declarative-ux AC8

<!-- case: REST tool serves an embedded bundle and its declared config | status: implemented -->
### REST tool serves an embedded bundle and its declared config

<!-- field: id -->
#### ID

TC8

<!-- field: use_case -->
#### Use case

rel15.0-uc001-declarative-ux

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2267, GH-2268, GH-2269

<!-- field: command -->
#### Command

(cd agent-core && go test ./internal/tools/rest/...)
(cd magefiles && go test ./appmanifest/...)
(cd applications/chatbot-mesh && go test ./magefiles/...)

<!-- field: expected | form: "tree" -->
#### Expected

- bundle_served_with_spa_fallback: true
- ui_config_json_matches_declaration: true
- root_and_bundle_together_rejected: true
- unknown_bundle_rejected: true
- observer_directory_holds_declarations_only: true
- embedded_source_accepted_by_manifest: true

<!-- field: traces -->
#### Traces

- rel15.0-uc001-declarative-ux S3
- srd004-declarative-ux AC9

<!-- case: No shared UI code is duplicated under applications/ | status: implemented -->
### No shared UI code is duplicated under applications/

<!-- field: id -->
#### ID

TC9

<!-- field: use_case -->
#### Use case

rel15.0-uc001-declarative-ux

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2275

<!-- field: command -->
#### Command

(cd magefiles && go test -run 'TestNoUIKitDuplicationInApplications|TestSweepUIDuplicationReportsEachShape' ./)

<!-- field: expected | form: "tree" -->
#### Expected

- no_splitpanelpath_copy_outside_kit: true
- no_monitor_or_trace_client_copy_outside_kit: true
- no_trace_proxy_route_under_applications: true

<!-- field: traces -->
#### Traces

- rel15.0-uc001-declarative-ux S4
- srd004-declarative-ux AC2
- srd004-declarative-ux AC8
