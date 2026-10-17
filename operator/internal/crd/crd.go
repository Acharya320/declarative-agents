// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

// Package crd carries the Agent CustomResourceDefinition inside the operator
// binary and applies it on boot, so the schema and the binary ship as one
// unit (applications srd009 R1.3).
package crd

import (
	"context"
	_ "embed"
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

//go:embed agent-crd.yaml
var agentCRD []byte

// FieldOwner names the operator as the server-side-apply field manager.
const FieldOwner = "agent-operator"

// Definition decodes the embedded CRD.
func Definition() (*apiextensionsv1.CustomResourceDefinition, error) {
	var def apiextensionsv1.CustomResourceDefinition
	if err := yaml.Unmarshal(agentCRD, &def); err != nil {
		return nil, fmt.Errorf("decode embedded Agent CRD: %w", err)
	}
	return &def, nil
}

// Apply creates the CRD or updates it to the embedded definition, keeping
// the live object's resource version.
func Apply(ctx context.Context, c client.Client) error {
	desired, err := Definition()
	if err != nil {
		return err
	}
	var live apiextensionsv1.CustomResourceDefinition
	err = c.Get(ctx, client.ObjectKey{Name: desired.Name}, &live)
	switch {
	case client.IgnoreNotFound(err) != nil:
		return fmt.Errorf("read CRD %s: %w", desired.Name, err)
	case err != nil:
		if err := c.Create(ctx, desired); err != nil {
			return fmt.Errorf("create CRD %s: %w", desired.Name, err)
		}
		return nil
	}
	live.Spec = desired.Spec
	live.Labels, live.Annotations = desired.Labels, desired.Annotations
	if err := c.Update(ctx, &live); err != nil {
		return fmt.Errorf("update CRD %s: %w", desired.Name, err)
	}
	return nil
}
