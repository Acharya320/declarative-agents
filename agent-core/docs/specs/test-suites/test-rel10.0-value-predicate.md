---
id: test-rel10.0-value-predicate
title: Test Suite for Release 10.0 Value Predicate Tool
release: "10.0"
---
<!-- Copyright (c) 2026 Nokia -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->

<!-- field: traces -->
## Traces

- rel10.0-uc001-value-predicate-branch

<!-- field: preconditions -->
## Preconditions

- srd041-value-predicate-tool defines the declaration, the closed operator set, operand typing and coercion, the signal contract, and purity.
- srd022-builtin-tools defines the BuiltinRegistry factory lookup and typed config decoding the word registers through.
- srd038-command-state-store defines the selector primitive operands resolve through.
- GH-770 implemented the word; every case below carries the go_test that proves it. The cases were authored with planned_go_test while the tests did not yet exist, because the formal evidence validator runs every go_test string against the real inventory and citing an unwritten test would have recorded false evidence.

<!-- field: test_cases -->
## Test cases

<!-- case: Each operator decides correctly at its boundary | status: implemented -->
### Each operator decides correctly at its boundary

Every operator in the closed set decides correctly on both sides of its
boundary, so no operator is off by one. The equal case is what separates
lt from lte and gt from gte, which is why each is exercised at the
boundary rather than only well clear of it.

<!-- field: use_case -->
#### Use case

rel10.0-uc001-value-predicate-branch

<!-- field: go_test -->
#### Go test

go test ./internal/tools/control -run TestValuePredicateOperators

<!-- field: inputs | form: "tree" -->
#### Inputs

- comparisons:
  1. - left: "5"
     - op: eq
     - right: "5"
     - expect: satisfied
  1. - left: "5"
     - op: eq
     - right: "6"
     - expect: unsatisfied
  1. - left: "5"
     - op: ne
     - right: "6"
     - expect: satisfied
  1. - left: "5"
     - op: lt
     - right: "5"
     - expect: unsatisfied
  1. - left: "4"
     - op: lt
     - right: "5"
     - expect: satisfied
  1. - left: "5"
     - op: lte
     - right: "5"
     - expect: satisfied
  1. - left: "5"
     - op: gt
     - right: "5"
     - expect: unsatisfied
  1. - left: "6"
     - op: gt
     - right: "5"
     - expect: satisfied
  1. - left: "5"
     - op: gte
     - right: "5"
     - expect: satisfied

<!-- field: expected | form: "tree" -->
#### Expected

- every_comparison: emits the configured signal for its expected outcome

<!-- field: traces -->
#### Traces

- rel10.0-uc001-value-predicate-branch S2
- srd041-value-predicate-tool AC2

<!-- case: Emptiness operators test the left operand alone | status: implemented -->
### Emptiness operators test the left operand alone

empty and non_empty judge the left operand and ignore any right operand,
and an absent right operand does not fail registration for them. An empty
string, an empty list, an empty object, and a null are all empty; every
other resolved value is not.

<!-- field: use_case -->
#### Use case

rel10.0-uc001-value-predicate-branch

<!-- field: go_test -->
#### Go test

go test ./internal/tools/control -run TestValuePredicateEmptiness

<!-- field: inputs | form: "tree" -->
#### Inputs

- empty_values:
  1. ""
  1. []
  1. {}
  1. null
- non_empty_values:
  1. "0"
  1. "false"
  1. 1. 0
  1. - a: 1
- right_operand: absent

<!-- field: expected | form: "tree" -->
#### Expected

- empty_values: emit the satisfied signal under empty
- non_empty_values: emit the satisfied signal under non_empty
- registration: succeeds with no right operand

<!-- field: traces -->
#### Traces

- rel10.0-uc001-value-predicate-branch S2
- srd041-value-predicate-tool AC2

<!-- case: A string-typed count compares numerically under the default operand type | status: implemented -->
### A string-typed count compares numerically under the default operand type

A REST read of a scalar body yields a string, so a count read back from a
store arrives as "10" rather than 10. Under the default number operand
type it orders above "6"; a lexicographic comparison would put it below,
which is the bug this case exists to catch. The string operand type
orders lexicographically on purpose.

<!-- field: use_case -->
#### Use case

rel10.0-uc001-value-predicate-branch

<!-- field: go_test -->
#### Go test

go test ./internal/tools/control -run TestValuePredicateNumericCoercion

<!-- field: inputs | form: "keyed table" -->
#### Inputs

| id | left | op | right | operand_type |
| --- | --- | --- | --- | --- |
| numeric | "10" | gt | "6" | default |
| lexicographic | "10" | gt | "6" | string |

<!-- field: expected | form: "tree" -->
#### Expected

- numeric: satisfied
- lexicographic: unsatisfied

<!-- field: traces -->
#### Traces

- rel10.0-uc001-value-predicate-branch S3
- srd041-value-predicate-tool AC3

<!-- case: An unresolvable operand emits CommandError and never the unsatisfied signal | status: implemented -->
### An unresolvable operand emits CommandError and never the unsatisfied signal

This is the assertion that matters most, because the failure it prevents
is silent. A mistyped or absent selector must not read as a legitimate
negative: a machine would take the unsatisfied branch on a comparison
that never ran, and nothing in the trace would say so. The same rule
applies to a coercion that fails.

<!-- field: use_case -->
#### Use case

rel10.0-uc001-value-predicate-branch

<!-- field: go_test -->
#### Go test

go test ./internal/tools/control -run TestValuePredicateUnresolvedOperandIsAFault

<!-- field: inputs | form: "tree" -->
#### Inputs

- unresolved_selector: a selector naming a path absent from the previous Result
- uncoercible:
  - left: not-a-number
  - op: gt
  - right: "0"
  - operand_type: number

<!-- field: expected | form: "tree" -->
#### Expected

- signal: CommandError
- never: the configured unsatisfied signal
- output: names the operand that could not be resolved or coerced

<!-- field: traces -->
#### Traces

- rel10.0-uc001-value-predicate-branch S4
- srd041-value-predicate-tool AC4

<!-- case: Malformed config fails registration with the tool name | status: implemented -->
### Malformed config fails registration with the tool name

An operator outside the closed set and a missing signal name each fail at
registration rather than at the first run, so a misconfigured predicate
never reaches a machine. The error names the tool, per srd022 R3.3.

<!-- field: use_case -->
#### Use case

rel10.0-uc001-value-predicate-branch

<!-- field: go_test -->
#### Go test

go test ./internal/tools/control -run TestValuePredicateRegistrationRejectsMalformedConfig

<!-- field: inputs | form: "tree" -->
#### Inputs

- unknown_operator:
  - op: approximately
- missing_signal:
  - satisfied: Ok
  - unsatisfied: absent

<!-- field: expected | form: "tree" -->
#### Expected

- registration: fails
- error: names the tool and the offending field

<!-- field: traces -->
#### Traces

- rel10.0-uc001-value-predicate-branch S5
- srd041-value-predicate-tool AC5

<!-- case: A machine branches on a value through the declared word | status: implemented -->
### A machine branches on a value through the declared word

The end-to-end shape: a machine dispatches the predicate from one state
and reaches two different next states depending on the value, with both
outcomes as visible transitions. This is what makes the word a branch in
the grammar rather than a decision hidden in Go.

<!-- field: use_case -->
#### Use case

rel10.0-uc001-value-predicate-branch

<!-- field: go_test -->
#### Go test

go test ./internal/tools/control -run TestValuePredicateDrivesTwoTransitions

<!-- field: inputs | form: "tree" -->
#### Inputs

- machine: one state dispatching the predicate, routing each configured signal to its own next state
- values:
  1. a value that satisfies
  1. a value that does not

<!-- field: expected | form: "tree" -->
#### Expected

- transitions: the two values reach the two different next states
- side_effects: none recorded; undo is a noop

<!-- field: traces -->
#### Traces

- rel10.0-uc001-value-predicate-branch S1
- srd041-value-predicate-tool AC1
- srd041-value-predicate-tool AC6

<!-- case: A predicate compares across an intervening word | status: implemented -->
### A predicate compares across an intervening word

A predicate compares a value from an earlier labelled step against one
from the previous Result, so a machine can measure a change across an
action rather than only a level after it. The words are not adjacent, and
the word between them is an exec word whose Result carries nothing
forward, so the earlier value is reachable only through command state
(GH-774).

<!-- field: use_case -->
#### Use case

rel10.0-uc001-value-predicate-branch

<!-- field: go_test -->
#### Go test

go test ./internal/tools/control -run TestValuePredicateComparesAcrossAnInterveningWord

<!-- field: inputs | form: "tree" -->
#### Inputs

- before_step: count_before, output {mapped:{count:'4'}}
- after_previous_result: "{mapped:{count:'6'}} and {mapped:{count:'4'}}"
- comparison: $.mapped.count gt $from(count_before).mapped.count

<!-- field: expected | form: "tree" -->
#### Expected

- grew: satisfied
- unchanged: unsatisfied
- both_operands_from_earlier_steps: satisfied

<!-- field: traces -->
#### Traces

- rel10.0-uc001-value-predicate-branch S1
- srd041-value-predicate-tool AC1

<!-- case: An unreachable label is a fault, not a shortfall | status: implemented -->
### An unreachable label is a fault, not a shortfall

The failure separation holds for the command-state operand form too. A
label that never ran, a path absent from the step it names, and a missing
view each emit CommandError rather than the unsatisfied signal, so a
machine never reads an unreachable operand as a legitimate negative.

<!-- field: use_case -->
#### Use case

rel10.0-uc001-value-predicate-branch

<!-- field: go_test -->
#### Go test

go test ./internal/tools/control -run TestValuePredicateUnreachableLabelIsAFault

<!-- field: inputs | form: "tree" -->
#### Inputs

- no_such_label: $from(count_before).mapped.count with only other_step recorded
- missing_path: $from(count_before).mapped.count with count_before output {mapped:{}}
- no_view: $from(count_before).mapped.count with no command-state view

<!-- field: expected | form: "tree" -->
#### Expected

- signal: CommandError
- never: the configured unsatisfied signal
- output: explains which operand could not be reached

<!-- field: traces -->
#### Traces

- rel10.0-uc001-value-predicate-branch S4
- srd041-value-predicate-tool AC4
