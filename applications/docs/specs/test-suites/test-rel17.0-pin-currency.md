---
id: test-rel17.0-pin-currency
title: Release 17.0 Pin Currency Test Suite
release: "17.0"
status: implemented
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: test_case_count -->
## Test case count

```yaml
9
```

<!-- field: overview -->
## Overview

TC1 parses and traces the Release 17.0 documents. TC2 through TC9 run the
pin survey GH-2416 added. Every case runs without a network: the registry
cases drive an httptest server, so the suite proves the client's behaviour
rather than any upstream project's release cadence, which would make the
suite fail on somebody else's release day.

<!-- field: traces -->
## Traces

- rel17.0-uc001-pin-currency

<!-- field: srds -->
## Srds

- srd006-pin-currency

<!-- field: preconditions -->
## Preconditions

- Run commands from the repository root.
- Python 3 with PyYAML, Go, Mage, and Helm are available.
- No case reaches the network; a planned case is not passing evidence until its command exists and passes.

<!-- field: test_cases -->
## Test cases

<!-- case: Release 17.0 documents parse and trace reciprocally | status: implemented -->
### Release 17.0 documents parse and trace reciprocally

<!-- field: id -->
#### Id

TC1

<!-- field: use_case -->
#### Use case

rel17.0-uc001-pin-currency

<!-- field: evidence_type -->
#### Evidence type

executable_structural

<!-- field: command -->
#### Command

```yaml
|
    python3 - <<'PY'
    from pathlib import Path
    import subprocess
    import yaml

    root = Path("applications/docs")
    paths = {
        "srd": root / "specs/software-requirements/srd006-pin-currency.md",
        "use_case": root / "specs/use-cases/rel17.0-uc001-pin-currency.md",
        "suite": root / "specs/test-suites/test-rel17.0-pin-currency.md",
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
    assert "Pin currency" in titles, titles
    verbs = next(s for s in docs["guideline"]["sections"]
                 if s["title"] == "Standard mage verbs")["content"]
    assert "| bump |" in verbs, "the root verb table does not carry bump"
    currency = next(s for s in docs["guideline"]["sections"]
                    if s["title"] == "Pin currency")["content"]
    assert srd["id"] in currency

    index = docs["index"]
    assert any(e["id"] == srd["id"] for e in index["srd_index"])
    assert any(e["id"] == use_case["id"] for e in index["use_case_index"])
    assert any(e["id"] == suite["id"] for e in index["test_suite_index"])
    assert any(e["version"] == "17.0" for e in index["roadmap_summary"])
    release = next(r for r in docs["roadmap"]["releases"] if r["version"] == "17.0")
    assert any(e["id"] == use_case["id"] for e in release["use_cases"])
    assert any(e["id"] == suite["id"] for e in release["test_suites"])
    print("validated Release 17.0 pin currency trace closure")
    PY
```

<!-- field: expected -->
#### Expected

```yaml
yaml_parses: true
design_constitution_fields_present: true
requirement_to_acceptance_criterion_closure: true
reciprocal_srd_use_case_suite_traces: true
guideline_carries_the_section_and_the_verb: true
specification_and_roadmap_indexes_present: true
```

<!-- field: traces -->
#### Traces

- rel17.0-uc001-pin-currency S1
- srd006-pin-currency AC1

<!-- case: The inventory derives chart pins and declares the rest | status: implemented -->
### The inventory derives chart pins and declares the rest

<!-- field: id -->
#### Id

TC2

<!-- field: use_case -->
#### Use case

rel17.0-uc001-pin-currency

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2416

<!-- field: command -->
#### Command

cd magefiles && go test ./pinsurvey/... .

<!-- field: expected -->
#### Expected

```yaml
chart_pins_read_from_rendered_charts: true
chart_pins_not_restated_in_any_list: true
declared_list_covers_kindrig_node_image_and_smoke_registry: true
every_declared_entry_names_file_and_symbol: true
```

<!-- field: traces -->
#### Traces

- rel17.0-uc001-pin-currency S1
- srd006-pin-currency AC1

<!-- case: A declared entry that no longer describes its location fails | status: implemented -->
### A declared entry that no longer describes its location fails

<!-- field: id -->
#### Id

TC3

<!-- field: use_case -->
#### Use case

rel17.0-uc001-pin-currency

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2416

<!-- field: command -->
#### Command

cd magefiles && go test . -run TestDeclaredPins

<!-- field: expected -->
#### Expected

```yaml
moved_or_renamed_pin_fails_the_test: true
failure_names_the_entry_and_the_location: true
```

<!-- field: traces -->
#### Traces

- rel17.0-uc001-pin-currency S1
- srd006-pin-currency AC1

<!-- case: A re-pushed tag is reported as moved, ahead of staleness | status: implemented -->
### A re-pushed tag is reported as moved, ahead of staleness

<!-- field: id -->
#### Id

TC4

<!-- field: use_case -->
#### Use case

rel17.0-uc001-pin-currency

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2416

<!-- field: command -->
#### Command

cd magefiles && go test ./pinsurvey/...

<!-- field: expected -->
#### Expected

```yaml
digest_mismatch_reported_as_moved: true
report_lists_moved_before_behind: true
both_digests_named: true
```

<!-- field: traces -->
#### Traces

- rel17.0-uc001-pin-currency S2
- srd006-pin-currency AC2

<!-- case: A newer release is reported with both versions | status: implemented -->
### A newer release is reported with both versions

<!-- field: id -->
#### Id

TC5

<!-- field: use_case -->
#### Use case

rel17.0-uc001-pin-currency

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2416

<!-- field: command -->
#### Command

cd magefiles && go test ./pinsurvey/...

<!-- field: expected -->
#### Expected

```yaml
newer_tag_reported_as_behind: true
pinned_and_newest_versions_named: true
tag_list_pagination_followed: true
```

<!-- field: traces -->
#### Traces

- rel17.0-uc001-pin-currency S3
- srd006-pin-currency AC3

<!-- case: Version comparison applies to semantic versions only | status: implemented -->
### Version comparison applies to semantic versions only

<!-- field: id -->
#### Id

TC6

<!-- field: use_case -->
#### Use case

rel17.0-uc001-pin-currency

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2416

<!-- field: command -->
#### Command

cd magefiles && go test ./pinsurvey/...

<!-- field: expected -->
#### Expected

```yaml
v_prefixed_and_bare_semver_compared: true
non_version_tag_compared_on_integrity_alone: true
prerelease_does_not_count_as_newer: true
```

<!-- field: traces -->
#### Traces

- rel17.0-uc001-pin-currency S3
- srd006-pin-currency AC3

<!-- case: An unreachable registry reports unknown, never current | status: implemented -->
### An unreachable registry reports unknown, never current

<!-- field: id -->
#### Id

TC7

<!-- field: use_case -->
#### Use case

rel17.0-uc001-pin-currency

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2416

<!-- field: command -->
#### Command

cd magefiles && go test ./pinsurvey/...

<!-- field: expected -->
#### Expected

```yaml
registry_error_reported_as_unknown: true
absent_tag_reported_as_unknown: true
unknown_not_counted_among_current: true
reason_named_in_the_finding: true
```

<!-- field: traces -->
#### Traces

- rel17.0-uc001-pin-currency S5
- srd006-pin-currency AC5

<!-- case: The target reports and writes nothing | status: implemented -->
### The target reports and writes nothing

<!-- field: id -->
#### Id

TC8

<!-- field: use_case -->
#### Use case

rel17.0-uc001-pin-currency

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2416

<!-- field: command -->
#### Command

cd magefiles && go test . -run TestBump

<!-- field: expected -->
#### Expected

```yaml
working_tree_unchanged_across_a_run: true
exit_zero_with_moved_and_behind_pins: true
audit_does_not_invoke_the_target: true
no_test_reaches_the_network: true
```

<!-- field: traces -->
#### Traces

- rel17.0-uc001-pin-currency S4
- srd006-pin-currency AC4

<!-- case: Authentication and reference parsing reach the real registries | status: implemented -->
### Authentication and reference parsing reach the real registries

<!-- field: id -->
#### Id

TC9

<!-- field: use_case -->
#### Use case

rel17.0-uc001-pin-currency

<!-- field: evidence_type -->
#### Evidence type

go_test

<!-- field: implemented_by -->
#### Implemented by

GH-2416

<!-- field: command -->
#### Command

cd magefiles && go test ./pinsurvey/...

<!-- field: expected -->
#### Expected

```yaml
unauthenticated_request_succeeds_where_allowed: true
challenge_then_token_then_retry: true
bare_repository_resolves_under_library: true
finding_names_location_recorded_and_upstream_values: true
```

<!-- field: traces -->
#### Traces

- rel17.0-uc001-pin-currency S5
- srd006-pin-currency AC6
