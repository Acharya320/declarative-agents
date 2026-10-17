// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

// Package controller reconciles an Agent into one Deployment and one
// Service (applications srd009 R2, R3).
package controller

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/Nokia-Bell-Labs/declarative-agents/operator/api/v1alpha1"
)

// requeueAfter is how long the operator waits before it looks again at an
// Agent whose profile ConfigMap is missing (srd009 R2.4).
const requeueAfter = 30 * time.Second

// AgentReconciler owns the Deployment and Service of each Agent.
type AgentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// SetupWithManager watches Agents and the objects they own.
func (r *AgentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Agent{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}

// Reconcile brings the Deployment and Service of one Agent to the state its
// spec implies and writes what it observed into the status.
func (r *AgentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var agent v1alpha1.Agent
	if err := r.Get(ctx, req.NamespacedName, &agent); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err) // deleted: the owned objects follow it
	}
	if reason, problem := validate(&agent); problem != "" {
		return r.fail(ctx, &agent, reason, problem, ctrl.Result{})
	}
	var profile corev1.ConfigMap
	if err := r.Get(ctx, client.ObjectKey{Namespace: agent.Namespace, Name: agent.Spec.Profile.ConfigMap}, &profile); err != nil {
		if client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, fmt.Errorf("read profile ConfigMap: %w", err)
		}
		return r.fail(ctx, &agent, ReasonProfileMissing,
			fmt.Sprintf("ConfigMap %s/%s does not exist", agent.Namespace, agent.Spec.Profile.ConfigMap), ctrl.Result{RequeueAfter: requeueAfter})
	}
	dep := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: agent.Name, Namespace: agent.Namespace}}
	if _, err := controllerutil.CreateOrUpdate(ctx, r.Client, dep, func() error {
		desiredDeployment(&agent, dep)
		return controllerutil.SetControllerReference(&agent, dep, r.Scheme)
	}); err != nil {
		return ctrl.Result{}, fmt.Errorf("reconcile Deployment: %w", err)
	}
	svc := &corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: agent.Name, Namespace: agent.Namespace}}
	if _, err := controllerutil.CreateOrUpdate(ctx, r.Client, svc, func() error {
		desiredService(&agent, svc)
		return controllerutil.SetControllerReference(&agent, svc, r.Scheme)
	}); err != nil {
		return ctrl.Result{}, fmt.Errorf("reconcile Service: %w", err)
	}
	observe(&agent, dep)
	return ctrl.Result{}, r.Status().Update(ctx, &agent)
}

// fail records a False Ready condition with its reason and changes nothing
// else (srd009 R2.4).
func (r *AgentReconciler) fail(ctx context.Context, agent *v1alpha1.Agent, reason, message string, result ctrl.Result) (ctrl.Result, error) {
	setReady(agent, metav1.ConditionFalse, reason, message)
	return result, r.Status().Update(ctx, agent)
}

// validate reports a spec the operator cannot render.
func validate(agent *v1alpha1.Agent) (reason, problem string) {
	switch {
	case agent.Spec.Image == "":
		return ReasonInvalidSpec, "spec.image is empty"
	case agent.Spec.Profile.ConfigMap == "" || agent.Spec.Profile.Path == "":
		return ReasonInvalidSpec, "spec.profile needs configMap and path"
	case agent.ControlPortNumber() == 0:
		return ReasonInvalidSpec, fmt.Sprintf("spec.controlPort %q names no declared port", agent.Spec.ControlPort)
	}
	return "", ""
}
