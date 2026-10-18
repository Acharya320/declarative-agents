// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package controller

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Nokia-Bell-Labs/declarative-agents/operator/api/v1alpha1"
)

// Condition reasons the operator writes (srd009 R2.4, R3.1).
const (
	ReasonAvailable      = "Available"
	ReasonProgressing    = "Progressing"
	ReasonProfileMissing = "ProfileMissing"
	ReasonInvalidSpec    = "InvalidSpec"
	ReasonProfileInvalid = "ProfileInvalid"
)

// setReady records the Ready condition and the observed generation.
func setReady(agent *v1alpha1.Agent, status metav1.ConditionStatus, reason, message string) {
	agent.Status.ObservedGeneration = agent.Generation
	meta.SetStatusCondition(&agent.Status.Conditions, metav1.Condition{
		Type: v1alpha1.ConditionReady, Status: status, Reason: reason, Message: message,
		ObservedGeneration: agent.Generation,
	})
}

// observe fills the status from the Deployment: Ready is True only when the
// available replicas equal the desired count (srd009 R3.1).
func observe(agent *v1alpha1.Agent, dep *appsv1.Deployment) {
	agent.Status.ReadyReplicas = dep.Status.AvailableReplicas
	desired := agent.DesiredReplicas()
	if dep.Status.ObservedGeneration >= dep.Generation && dep.Status.AvailableReplicas == desired {
		setReady(agent, metav1.ConditionTrue, ReasonAvailable, "every replica is available")
		return
	}
	setReady(agent, metav1.ConditionFalse, ReasonProgressing, "the Deployment is rolling out")
}
