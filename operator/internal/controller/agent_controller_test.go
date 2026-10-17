// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package controller

import (
	"context"
	"strings"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/Nokia-Bell-Labs/declarative-agents/operator/api/v1alpha1"
)

// collector is the Agent the kind proof applies (srd009 R5.1): the catalog
// collector's profile, three ports, the control port as the readiness
// signal, and the bind host set to every interface.
func collector() *v1alpha1.Agent {
	return &v1alpha1.Agent{
		ObjectMeta: metav1.ObjectMeta{Name: "collector", Namespace: "demo", Generation: 1},
		Spec: v1alpha1.AgentSpec{
			Profile:     v1alpha1.ProfileSource{ConfigMap: "collector-profile", Path: "profile.yaml"},
			Image:       "localhost/declarative-agents/runtime/agent-core:git-abc123def456-linux-arm64",
			Ports:       []v1alpha1.AgentPort{{Name: "control", Port: 18191}, {Name: "monitor", Port: 18192}, {Name: "query", Port: 18193}},
			ControlPort: "control",
			Env:         []corev1.EnvVar{{Name: "COLLECTOR_BIND_HOST", Value: "0.0.0.0"}},
			WorkVolumes: []v1alpha1.WorkVolume{{Name: "work", MountPath: "/work"}},
		},
	}
}

func newReconciler(t *testing.T, objects ...client.Object) *AgentReconciler {
	t.Helper()
	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	if err := v1alpha1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).WithStatusSubresource(&v1alpha1.Agent{}).Build()
	return &AgentReconciler{Client: c, Scheme: scheme}
}

func profileConfigMap() *corev1.ConfigMap {
	return &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "collector-profile", Namespace: "demo"},
		Data: map[string]string{"profile.yaml": "name: collector\n"}}
}

func reconcile(t *testing.T, r *AgentReconciler, agent *v1alpha1.Agent) ctrl.Result {
	t.Helper()
	result, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Name: agent.Name, Namespace: agent.Namespace}})
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func get[T client.Object](t *testing.T, r *AgentReconciler, obj T, name string) T {
	t.Helper()
	if err := r.Get(context.Background(), client.ObjectKey{Name: name, Namespace: "demo"}, obj); err != nil {
		t.Fatalf("get %s: %v", name, err)
	}
	return obj
}

func ready(agent *v1alpha1.Agent) *metav1.Condition {
	return meta.FindStatusCondition(agent.Status.Conditions, v1alpha1.ConditionReady)
}

// TestReconcileCreatesDeploymentAndService covers srd009 R2.1, R2.2, R2.3:
// the Deployment runs the image with --profile at the mounted path, probes
// the control port, mounts the profile and work volumes; the Service
// exposes the ports; both carry the Agent as their controlling owner.
func TestReconcileCreatesDeploymentAndService(t *testing.T) {
	agent := collector()
	r := newReconciler(t, agent, profileConfigMap())
	reconcile(t, r, agent)
	dep := get(t, r, &appsv1.Deployment{}, "collector")
	c := dep.Spec.Template.Spec.Containers[0]
	if c.Image != agent.Spec.Image || strings.Join(c.Args, " ") != "--profile /profiles/profile.yaml" {
		t.Errorf("container = %s %v", c.Image, c.Args)
	}
	if c.ReadinessProbe == nil || c.ReadinessProbe.HTTPGet.Path != healthPath || c.ReadinessProbe.HTTPGet.Port.IntValue() != 18191 {
		t.Errorf("readiness probe = %+v", c.ReadinessProbe)
	}
	if len(c.Ports) != 3 || len(c.VolumeMounts) != 2 || c.Env[0].Name != "COLLECTOR_BIND_HOST" {
		t.Errorf("ports %d, mounts %d, env %v", len(c.Ports), len(c.VolumeMounts), c.Env)
	}
	if v := dep.Spec.Template.Spec.Volumes; len(v) != 2 || v[0].ConfigMap.Name != "collector-profile" || v[1].EmptyDir == nil {
		t.Errorf("volumes = %+v", v)
	}
	if *dep.Spec.Replicas != 1 {
		t.Errorf("replicas = %d, want 1 by default", *dep.Spec.Replicas)
	}
	svc := get(t, r, &corev1.Service{}, "collector")
	if len(svc.Spec.Ports) != 3 || svc.Spec.Ports[0].Name != "control" || svc.Spec.Ports[0].Port != 18191 || svc.Spec.Selector["app.kubernetes.io/instance"] != "collector" {
		t.Errorf("service = %+v", svc.Spec)
	}
	for _, owned := range []metav1.Object{dep, svc} {
		owner := metav1.GetControllerOf(owned)
		if owner == nil || owner.Kind != v1alpha1.Kind || owner.Name != "collector" {
			t.Errorf("%s owner = %+v", owned.GetName(), owner)
		}
	}
	if cond := ready(get(t, r, &v1alpha1.Agent{}, "collector")); cond == nil || cond.Status != metav1.ConditionFalse || cond.Reason != ReasonProgressing {
		t.Errorf("Ready before rollout = %+v, want False/Progressing", cond)
	}
}

// TestReconcileRollsOnSpecChange covers srd009 R2.3: a new image reaches the
// Deployment's pod template on the next reconcile.
func TestReconcileRollsOnSpecChange(t *testing.T) {
	agent := collector()
	r := newReconciler(t, agent, profileConfigMap())
	reconcile(t, r, agent)
	live := get(t, r, &v1alpha1.Agent{}, "collector")
	live.Spec.Image = "localhost/declarative-agents/runtime/agent-core:git-fedcba987654-linux-arm64"
	live.Spec.Ports = append(live.Spec.Ports, v1alpha1.AgentPort{Name: "extra", Port: 18194})
	if err := r.Update(context.Background(), live); err != nil {
		t.Fatal(err)
	}
	reconcile(t, r, agent)
	dep := get(t, r, &appsv1.Deployment{}, "collector")
	if img := dep.Spec.Template.Spec.Containers[0].Image; img != live.Spec.Image {
		t.Errorf("image = %s, want the new image", img)
	}
	if svc := get(t, r, &corev1.Service{}, "collector"); len(svc.Spec.Ports) != 4 {
		t.Errorf("service ports = %d, want 4", len(svc.Spec.Ports))
	}
}

// TestReconcileReportsMissingProfile covers srd009 R2.4: no ConfigMap means a
// False Ready condition naming it, a requeue, and no Deployment.
func TestReconcileReportsMissingProfile(t *testing.T) {
	agent := collector()
	r := newReconciler(t, agent)
	if result := reconcile(t, r, agent); result.RequeueAfter != requeueAfter {
		t.Errorf("result = %+v, want requeue after %s", result, requeueAfter)
	}
	cond := ready(get(t, r, &v1alpha1.Agent{}, "collector"))
	if cond == nil || cond.Status != metav1.ConditionFalse || cond.Reason != ReasonProfileMissing || !strings.Contains(cond.Message, "demo/collector-profile") {
		t.Errorf("Ready = %+v", cond)
	}
	if err := r.Get(context.Background(), client.ObjectKey{Name: "collector", Namespace: "demo"}, &appsv1.Deployment{}); client.IgnoreNotFound(err) != nil || err == nil {
		t.Errorf("a Deployment exists without a profile: %v", err)
	}
	bad := collector()
	bad.Name, bad.Spec.ControlPort = "bad", "nowhere"
	r2 := newReconciler(t, bad, profileConfigMap())
	reconcile(t, r2, bad)
	if cond := ready(get(t, r2, &v1alpha1.Agent{}, "bad")); cond == nil || cond.Reason != ReasonInvalidSpec {
		t.Errorf("invalid spec Ready = %+v", cond)
	}
}

// TestReconcileReportsReady covers srd009 R3.1: once the Deployment reports
// every replica available, Ready is True with the count.
func TestReconcileReportsReady(t *testing.T) {
	agent := collector()
	r := newReconciler(t, agent, profileConfigMap())
	reconcile(t, r, agent)
	dep := get(t, r, &appsv1.Deployment{}, "collector")
	dep.Status.AvailableReplicas, dep.Status.ObservedGeneration = 1, dep.Generation
	if err := r.Status().Update(context.Background(), dep); err != nil {
		t.Fatal(err)
	}
	reconcile(t, r, agent)
	live := get(t, r, &v1alpha1.Agent{}, "collector")
	cond := ready(live)
	if cond == nil || cond.Status != metav1.ConditionTrue || cond.Reason != ReasonAvailable || live.Status.ReadyReplicas != 1 || live.Status.ObservedGeneration != live.Generation {
		t.Errorf("status = %+v", live.Status)
	}
}
