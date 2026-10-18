// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

// Package startupcheck is the startup boundary agent-core applies to a loaded
// profile closure: the machine's runtime policy, the wiring between the
// machine and its selected tools, and the same check on every machine the
// closure's walk reaches. cmd/agent runs it before it serves, and
// pkg/profileservers runs it before an operator schedules a pod (srd009
// R2.4), so both reject exactly the same profiles.
package startupcheck

import (
	"fmt"

	internalload "github.com/Nokia-Bell-Labs/declarative-agents/agent-core/internal/load"
	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/internal/runtime/core"
	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/internal/tools/catalog"
	toolrest "github.com/Nokia-Bell-Labs/declarative-agents/agent-core/internal/tools/rest"
	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/internal/typesys"
	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/profileaudit"
)

// ValidateClosure applies the boundary to a loaded closure and returns its
// validated machine. Every machine the walk reaches is checked here and not
// by the caller: the profile's own machine is the one this function holds,
// and the machines its endpoints dispatch, its child profiles run, and its
// evaluator points loop over are reachable only through the walk (GH-2059,
// GH-2060).
func ValidateClosure(closure *internalload.Closure) (core.MachineSpec, error) {
	machineSpec := closure.Machine
	if err := core.ValidateRequiredMachinePolicy(machineSpec); err != nil {
		return core.MachineSpec{}, fmt.Errorf("load machine runtime policy: %w", err)
	}
	if err := ValidateToolWiring(
		machineSpec, closure.Selected, closure.Types, closure.Rest,
		catalog.ExhaustivenessInputs{External: RequestSourceSignals(closure.Rest)},
	); err != nil {
		if machineSpec.TemplatePath() != "" {
			err = fmt.Errorf("machine %s: %w", core.DescribeMachine(closure.Profile.Machine, machineSpec), err)
		}
		return core.MachineSpec{}, err
	}
	if err := profileaudit.ValidateClosureWithOptions(closure, profileaudit.Options{
		OnMachine: ValidateReachedMachine,
	}); err != nil {
		return core.MachineSpec{}, fmt.Errorf("inspect profile timeout closure: %w", err)
	}
	return machineSpec, nil
}

// ValidateReachedMachine applies the boundary to one machine the closure walk
// reaches. A request-scoped machine may consume the signals its request
// source injects and the seed label the runtime supplies.
func ValidateReachedMachine(reached profileaudit.ReachedMachine) error {
	inputs := catalog.ExhaustivenessInputs{External: RequestSourceSignals(reached.Rest)}
	var runtimeLabels []string
	kind := "machine"
	if reached.RequestScoped {
		inputs.External = append(inputs.External, reached.InitialSignals...)
		runtimeLabels = []string{catalog.RequestSeedLabel}
		kind = "request machine"
	}
	err := ValidateToolWiring(
		reached.Machine, reached.Selected, reached.Types, reached.Rest,
		inputs, runtimeLabels...,
	)
	if err != nil {
		return fmt.Errorf("%s %s: %w", kind, core.DescribeMachine(reached.MachinePath, reached.Machine), err)
	}
	return nil
}

// ValidateToolWiring is the ordinary startup boundary. It rejects unselected
// named actions, machine/tool signal mismatches, incomplete parse-retry
// routes, and reversible effects without receipt-consuming undo. Full
// six-section contract completeness remains an authoring and audit concern.
func ValidateToolWiring(
	machine core.MachineSpec, defs []catalog.ToolDef, types *typesys.Registry,
	operations catalog.RESTOperations,
	exhaustiveness catalog.ExhaustivenessInputs, runtimeLabels ...string,
) error {
	if err := catalog.ValidateMachineActions(machine, defs); err != nil {
		return err
	}
	if err := catalog.ValidateToolPhases(machine, defs); err != nil {
		return err
	}
	if err := catalog.ValidateParseRetryWiring(machine, defs); err != nil {
		return err
	}
	if err := catalog.ValidateToolEmits(machine, defs); err != nil {
		return err
	}
	if err := catalog.ValidateProviderFailureOptIn(defs); err != nil {
		return err
	}
	if err := catalog.ValidateSelectorLabelsStrict(machine, defs, runtimeLabels...); err != nil {
		return err
	}
	if err := catalog.ValidateSelectorPathsStrict(machine, defs, types, operations); err != nil {
		return err
	}
	if err := catalog.ValidateMachineExhaustivenessStrict(machine, defs, exhaustiveness); err != nil {
		return err
	}
	return catalog.ValidateReceiptContracts(defs)
}

// RequestSourceSignals collects the signals a request signal source can
// inject, which reach the machine without any transition producing them
// (srd046).
func RequestSourceSignals(defs toolrest.Collection) []string {
	var signals []string
	for _, server := range defs.Servers {
		for _, endpoint := range server.Endpoints {
			signals = append(signals, endpoint.Signal)
			signals = append(signals, endpoint.AllowedSignals...)
			for _, signal := range endpoint.SignalSource.SignalMapping {
				signals = append(signals, signal)
			}
		}
	}
	return signals
}
