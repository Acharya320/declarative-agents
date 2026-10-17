---
id: test-rel16.0-chart-conformance
title: Release 16.0 Chart Conformance Test Suite
release: "16.0"
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: test_case_count | form: "scalar" -->
## Test case count

12

<!-- field: overview -->
## Overview

TC1 parses and traces the Release 16.0 documents. TC2 through TC9 run the
conformance module GH-2408 added. TC10 proves GH-2409 emptied the baseline:
every upstream image resolves through a digest, every container states a
pull policy, and both affected applications installed on a kind cluster on
the pinned images. Every case but TC1 runs under the repository root
audit.

<!-- field: traces -->
## Traces

- rel16.0-uc001-chart-conformance

<!-- field: srds -->
## Srds

- srd005-chart-conformance

<!-- field: preconditions -->
## Preconditions

- Run commands from the repository root.
- Python 3 with PyYAML, Go, Mage, and Helm are available.
- A planned case is not passing evidence until its command exists and passes.

<!-- field: test_cases -->
## Test cases

<!-- case: Release 16.0 documents parse and trace reciprocally | status: implemented -->
### Release 16.0 documents parse and trace reciprocally

<!-- field: id -->
#### Id

TC1

<!-- field: use_case -->
#### Use case

rel16.0-uc001-chart-conformance

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
    "srd": root / "specs/software-requirements/srd005-chart-conformance.md",
    "use_case": root / "specs/use-cases/rel16.0-uc001-chart-conformance.md",
    "suite": root / "specs/test-suites/test-rel16.0-chart-conformance.md",
    "index": root / "SPECIFICATIONS.md",
    "roadmap": root / "road-map.md",
    "guideline": Path("docs/engineering/eng01-kind-test-demo-rig.md"),
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
    "guideline": {"id", "title", "introduction", "sections", "references"},
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

titles = [section["title"] for section in docs["guideline"]["sections"]]
assert "Chart conformance rules" in titles, titles
rules = next(s for s in docs["guideline"]["sections"]
             if s["title"] == "Chart conformance rules")["content"]
for rule in ("C1", "C2", "C3", "C4", "C5", "C6", "C7"):
    assert f"| {rule} |" in rules, rule
assert srd["id"] in rules

index = docs["index"]
assert any(e["id"] == srd["id"] for e in index["srd_index"])
assert any(e["id"] == use_case["id"] for e in index["use_case_index"])
assert any(e["id"] == suite["id"] for e in index["test_suite_index"])
assert any(e["version"] == "16.0" for e in index["roadmap_summary"])
release = next(r for r in docs["roadmap"]["releases"] if r["version"] == "16.0")
assert any(e["id"] == use_case["id"] for e in release["use_cases"])
assert any(e["id"] == suite["id"] for e in release["test_suites"])
print("validated Release 16.0 chart conformance trace closure")
PY

<!-- field: expected | form: "tree" -->
#### Expected

- yaml_parses: true
- design_constitution_fields_present: true
- requirement_to_acceptance_criterion_closure: true
- reciprocal_srd_use_case_suite_traces: true
- guideline_carries_every_numbered_rule: true
- specification_and_roadmap_indexes_present: true

<!-- field: traces -->
#### Traces

- rel16.0-uc001-chart-conformance S1
- srd005-chart-conformance AC1

<!-- case: Every chart renders under its defaults and every checked-in overlay | status: implemented -->
### Every chart renders under its defaults and every checked-in overlay

<!-- field: id -->
#### Id

TC2

<!-- field: use_case -->
#### Use case

rel16.0-uc001-chart-conformance

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2408

<!-- field: command -->
#### Command

cd magefiles && go test . -run TestChartConformance

<!-- field: expected | form: "tree" -->
#### Expected

- every_application_chart_rendered: true
- library_chart_vendored_before_render: true
- overlays_discovered_not_hardcoded: true
- defaults_render_checked_alongside_overlays: true
- non_values_files_in_ci_not_treated_as_overlays: true

<!-- field: traces -->
#### Traces

- rel16.0-uc001-chart-conformance S1
- srd005-chart-conformance AC1

<!-- case: A floating tag fails the gate | status: implemented -->
### A floating tag fails the gate

<!-- field: id -->
#### Id

TC3

<!-- field: use_case -->
#### Use case

rel16.0-uc001-chart-conformance

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2408

<!-- field: command -->
#### Command

cd magefiles && go test . -run TestChartConformance

<!-- field: expected | form: "tree" -->
#### Expected

- latest_tag_reported: true
- absent_tag_reported: true
- floating_tag_in_defaults_reported: true

<!-- field: traces -->
#### Traces

- rel16.0-uc001-chart-conformance S2
- srd005-chart-conformance AC1

<!-- case: Third-party images need digests and repository images need commit tags | status: implemented -->
### Third-party images need digests and repository images need commit tags

<!-- field: id -->
#### Id

TC4

<!-- field: use_case -->
#### Use case

rel16.0-uc001-chart-conformance

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2408

<!-- field: command -->
#### Command

cd magefiles && go test . -run TestChartConformance

<!-- field: expected | form: "tree" -->
#### Expected

- third_party_image_without_digest_reported: true
- repository_image_classified_by_prefix: true
- repository_image_not_required_to_carry_digest: true

<!-- field: traces -->
#### Traces

- rel16.0-uc001-chart-conformance S2
- srd005-chart-conformance AC2

<!-- case: Pull policy and service type stay inside what kind provides | status: implemented -->
### Pull policy and service type stay inside what kind provides

<!-- field: id -->
#### Id

TC5

<!-- field: use_case -->
#### Use case

rel16.0-uc001-chart-conformance

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2408

<!-- field: command -->
#### Command

cd magefiles && go test . -run TestChartConformance

<!-- field: expected | form: "tree" -->
#### Expected

- image_pull_policy_always_reported: true
- load_balancer_service_reported: true
- switched_off_load_balancer_not_reported: true

<!-- field: traces -->
#### Traces

- rel16.0-uc001-chart-conformance S1
- srd005-chart-conformance AC1

<!-- case: A selected workload declares readiness | status: implemented -->
### A selected workload declares readiness

<!-- field: id -->
#### Id

TC6

<!-- field: use_case -->
#### Use case

rel16.0-uc001-chart-conformance

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2408

<!-- field: command -->
#### Command

cd magefiles && go test . -run TestChartConformance

<!-- field: expected | form: "tree" -->
#### Expected

- service_selector_resolved_to_workload: true
- workload_without_readiness_probe_reported: true
- unselected_workload_not_required_to_declare_one: true

<!-- field: traces -->
#### Traces

- rel16.0-uc001-chart-conformance S1
- srd005-chart-conformance AC1

<!-- case: Kind configurations pin the node image by version and digest | status: implemented -->
### Kind configurations pin the node image by version and digest

<!-- field: id -->
#### Id

TC7

<!-- field: use_case -->
#### Use case

rel16.0-uc001-chart-conformance

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2408

<!-- field: command -->
#### Command

cd magefiles && go test . -run TestChartConformance

<!-- field: expected | form: "tree" -->
#### Expected

- every_checked_in_kind_configuration_read: true
- node_image_without_digest_reported: true
- node_image_without_version_reported: true

<!-- field: traces -->
#### Traces

- rel16.0-uc001-chart-conformance S3
- srd005-chart-conformance AC3

<!-- case: The exception baseline suspends, refuses, and expires | status: implemented -->
### The exception baseline suspends, refuses, and expires

<!-- field: id -->
#### Id

TC8

<!-- field: use_case -->
#### Use case

rel16.0-uc001-chart-conformance

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2408

<!-- field: command -->
#### Command

cd magefiles && go test . -run TestChartConformance

<!-- field: expected | form: "tree" -->
#### Expected

- listed_violation_does_not_fail: true
- unlisted_violation_fails: true
- entry_whose_violation_no_longer_reproduces_fails: true
- gate_writes_no_edit_to_the_baseline: true
- entry_missing_issue_or_reason_rejected: true

<!-- field: traces -->
#### Traces

- rel16.0-uc001-chart-conformance S4
- srd005-chart-conformance AC4

<!-- case: A finding names chart, overlay, resource, value, and rule | status: implemented -->
### A finding names chart, overlay, resource, value, and rule

<!-- field: id -->
#### Id

TC9

<!-- field: use_case -->
#### Use case

rel16.0-uc001-chart-conformance

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2408

<!-- field: command -->
#### Command

cd magefiles && go test . -run TestChartConformance

<!-- field: expected | form: "tree" -->
#### Expected

- message_names_chart: true
- message_names_overlay: true
- message_names_resource_kind_and_name: true
- message_names_offending_value: true
- message_names_requirement_id: true

<!-- field: traces -->
#### Traces

- rel16.0-uc001-chart-conformance S5
- srd005-chart-conformance AC5

<!-- case: The seeded violations are pinned and the baseline holds no floating tag | status: implemented -->
### The seeded violations are pinned and the baseline holds no floating tag

<!-- field: id -->
#### Id

TC10

<!-- field: use_case -->
#### Use case

rel16.0-uc001-chart-conformance

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2409

<!-- field: command -->
#### Command

mage audit

<!-- field: expected | form: "tree" -->
#### Expected

- ollama_pinned_in_both_charts: true
- dolt_pinned: true
- remaining_third_party_images_carry_digests: true
- baseline_holds_no_floating_tag_entry: true
- demo_comes_up_on_the_pinned_images: true

<!-- field: traces -->
#### Traces

- rel16.0-uc001-chart-conformance S2
- srd005-chart-conformance AC1
- srd005-chart-conformance AC2
- srd005-chart-conformance AC6

<!-- case: Image references match the grammar and prohibited aliases fail | status: implemented -->
### Image references match the grammar and prohibited aliases fail

<!-- field: id -->
#### Id

TC11

<!-- field: use_case -->
#### Use case

rel16.0-uc001-chart-conformance

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2512

<!-- field: command -->
#### Command

cd magefiles && go test ./chartconf -run 'TestImageReferenceGrammar|TestConformingRenderProducesNoFinding'

<!-- field: expected | form: "tree" -->
#### Expected

- canonical_local_roles_accepted: true
- fully_qualified_upstream_accepted: true
- kindrig_alias_reported: true
- local_tag_reported: true
- untyped_hex_local_tag_reported: true
- smoke_repository_reported: true
- short_third_party_name_reported: true

<!-- field: traces -->
#### Traces

- rel16.0-uc001-chart-conformance S2
- srd005-chart-conformance AC7
- srd005-chart-conformance AC5

<!-- case: The applier Role covers every kind its chart renders | status: implemented -->
### The applier Role covers every kind its chart renders

<!-- field: id -->
#### Id

TC12

<!-- field: use_case -->
#### Use case

rel16.0-uc001-chart-conformance

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2533

<!-- field: command -->
#### Command

cd magefiles && go test ./chartconf -run 'TestApplierRole|TestReadOnlyApplier|TestPluralResource' && go test . -run TestChartConformance

<!-- field: expected | form: "tree" -->
#### Expected

- missing_rendered_kind_reported: true
- write_without_read_reported: true
- read_only_applier_not_checked: true
- shipped_charts_have_no_uncovered_kind: true

<!-- field: traces -->
#### Traces

- rel16.0-uc001-chart-conformance S6
- srd005-chart-conformance AC8
- srd005-chart-conformance AC5
