// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package controller

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/Nokia-Bell-Labs/declarative-agents/operator/api/v1alpha1"
)

// desiredService renders the one Service an Agent implies (srd009 R2.2):
// the declared ports by name, selecting the Deployment's pods.
func desiredService(agent *v1alpha1.Agent, svc *corev1.Service) {
	svc.Labels = labels(agent)
	svc.Spec.Selector = labels(agent)
	svc.Spec.Type = corev1.ServiceTypeClusterIP
	svc.Spec.Ports = nil
	for _, p := range agent.Spec.Ports {
		svc.Spec.Ports = append(svc.Spec.Ports, corev1.ServicePort{Name: p.Name, Port: p.Port, TargetPort: intstr.FromString(p.Name)})
	}
}
