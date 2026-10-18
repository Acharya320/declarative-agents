// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package controller

import (
	"context"
	"os"
	"path/filepath"
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

// catalogAgents is the catalog's agents directory, which holds the collector
// profile and the units it imports.
var catalogAgents = filepath.Join("..", "..", "..", "applications", "catalog", "agents")

// coreRoot is the agent-core checkout, for the profile's /opt/agent-core
// imports.
var coreRoot = filepath.Join("..", "..", "..", "agent-core")

// collector is the Agent the kind proof applies (srd009 R5.1): the catalog
// collector's profile, the image, and a work volume. Ports, the readiness
// probe, and the bind host come from the profile (R1.2).
func collector() *v1alpha1.Agent {
	return &v1alpha1.Agent{
		ObjectMeta: metav1.ObjectMeta{Name: "collector", Namespace: "demo", Generation: 1},
		Spec: v1alpha1.AgentSpec{
			Profile:     v1alpha1.ProfileSource{ConfigMap: "collector-profile", Path: "collector/profile.yaml"},
			Image:       "localhost/declarative-agents/runtime/agent-core:git-abc123def456-linux-arm64",
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
	return &AgentReconciler{Client: c, Scheme: scheme, CoreRoot: coreRoot}
}

// profileConfigMap is the real collector profile as the proof's ConfigMap:
// the profile directory's top-level files under collector__ and the units
// file it imports under units__.
func profileConfigMap() *corev1.ConfigMap {
	data := map[string]string{}
	entries, err := os.ReadDir(filepath.Join(catalogAgents, "collector"))
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if entry.Type().IsRegular() {
			data["collector__"+entry.Name()] = readFixture(filepath.Join(catalogAgents, "collector", entry.Name()))
		}
	}
	data["units__types-collector.yaml"] = readFixture(filepath.Join(catalogAgents, "units", "types-collector.yaml"))
	return &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "collector-profile", Namespace: "demo"}, Data: data}
}

func readFixture(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(data)
}

// editedProfile is profileConfigMap with one replacement in collector/rest.yaml.
func editedProfile(t *testing.T, old, replacement string) *corev1.ConfigMap {
	t.Helper()
	cm := profileConfigMap()
	rest := cm.Data["collector__rest.yaml"]
	if !strings.Contains(rest, old) {
		t.Fatalf("collector rest.yaml lacks %q", old)
	}
	cm.Data["collector__rest.yaml"] = strings.ReplaceAll(rest, old, replacement)
	return cm
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
	if c.Image != agent.Spec.Image || strings.Join(c.Args, " ") != "--profile /profiles/collector/profile.yaml" {
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

// TestReconcileMountsNestedProfileKeys covers the directory convention in
// srd009 R2.1: a ConfigMap whose keys spell directories with "__" mounts
// every key at the path it spells, so a profile's ../units import resolves.
func TestReconcileMountsNestedProfileKeys(t *testing.T) {
	agent := collector()
	profile := profileConfigMap()
	r := newReconciler(t, agent, profile)
	reconcile(t, r, agent)
	dep := get(t, r, &appsv1.Deployment{}, "collector")
	items := dep.Spec.Template.Spec.Volumes[0].ConfigMap.Items
	if len(items) != len(profile.Data) {
		t.Fatalf("%d items for %d keys", len(items), len(profile.Data))
	}
	paths := map[string]string{}
	for _, item := range items {
		paths[item.Key] = item.Path
	}
	if paths["collector__profile.yaml"] != "collector/profile.yaml" || paths["units__types-collector.yaml"] != "units/types-collector.yaml" {
		t.Errorf("items = %+v", items)
	}
	flat := &corev1.ConfigMap{Data: map[string]string{"profile.yaml": "name: x\n"}}
	if got := profileItems(flat); got != nil {
		t.Errorf("flat keys render items %+v, want none so the mount follows the ConfigMap", got)
	}
}

// TestReconcileDerivesPortsFromTheProfile covers srd009 R1.2, R2.1, R2.2:
// the container and Service ports, the readiness probe, the monitorPort
// label, and the bind-host environment all come from the collector's REST
// definitions, with nothing declared on the Agent.
func TestReconcileDerivesPortsFromTheProfile(t *testing.T) {
	agent := collector()
	r := newReconciler(t, agent, profileConfigMap())
	reconcile(t, r, agent)
	dep := get(t, r, &appsv1.Deployment{}, "collector")
	c := dep.Spec.Template.Spec.Containers[0]
	ports := map[string]int32{}
	for _, p := range c.Ports {
		ports[p.Name] = p.ContainerPort
	}
	if len(ports) != 3 || ports["control"] != 18191 || ports["monitor"] != 18192 || ports["query"] != 18193 {
		t.Errorf("container ports = %v", ports)
	}
	if c.ReadinessProbe == nil || c.ReadinessProbe.HTTPGet.Port.IntValue() != 18191 {
		t.Errorf("readiness probe = %+v, want the control server's port", c.ReadinessProbe)
	}
	if dep.Spec.Template.Labels[monitorPortLabel] != "18192" {
		t.Errorf("pod labels = %v, want monitorPort 18192", dep.Spec.Template.Labels)
	}
	if len(c.Env) != 1 || c.Env[0].Name != "COLLECTOR_BIND_HOST" || c.Env[0].Value != "0.0.0.0" {
		t.Errorf("env = %v, want the bind host on every interface", c.Env)
	}
	svc := get(t, r, &corev1.Service{}, "collector")
	if len(svc.Spec.Ports) != 3 || svc.Spec.Ports[0].Name != "control" || svc.Spec.Ports[0].Port != 18191 || svc.Spec.Ports[0].TargetPort.String() != "control" {
		t.Errorf("service ports = %+v", svc.Spec.Ports)
	}
}

// TestReconcileKeepsAnExplicitBindHost covers srd009 R2.1: an Agent that sets
// a server's bind variable itself is not overridden or duplicated.
func TestReconcileKeepsAnExplicitBindHost(t *testing.T) {
	agent := collector()
	agent.Spec.Env = []corev1.EnvVar{{Name: "COLLECTOR_BIND_HOST", Value: "127.0.0.1"}}
	r := newReconciler(t, agent, profileConfigMap())
	reconcile(t, r, agent)
	env := get(t, r, &appsv1.Deployment{}, "collector").Spec.Template.Spec.Containers[0].Env
	if len(env) != 1 || env[0].Value != "127.0.0.1" {
		t.Errorf("env = %v, want the explicit bind host alone", env)
	}
}

// TestReconcileReportsProfileInvalid covers srd009 R2.4: a profile with no
// health server, and one that fails agent-core's validation, report Ready
// False with reason ProfileInvalid and the loader's message, requeue, and
// create no Deployment.
func TestReconcileReportsProfileInvalid(t *testing.T) {
	noMachine := profileConfigMap()
	delete(noMachine.Data, "collector__machine.yaml")
	for name, profile := range map[string]*corev1.ConfigMap{
		"no health server": editedProfile(t, "path: /api/lifecycle/health", "path: /api/lifecycle/healthz"),
		"missing machine":  noMachine,
	} {
		agent := collector()
		r := newReconciler(t, agent, profile)
		if result := reconcile(t, r, agent); result.RequeueAfter != requeueAfter {
			t.Errorf("%s: result = %+v, want a requeue", name, result)
		}
		cond := ready(get(t, r, &v1alpha1.Agent{}, "collector"))
		if cond == nil || cond.Status != metav1.ConditionFalse || cond.Reason != ReasonProfileInvalid || cond.Message == "" {
			t.Errorf("%s: Ready = %+v", name, cond)
		}
		if err := r.Get(context.Background(), client.ObjectKey{Name: "collector", Namespace: "demo"}, &appsv1.Deployment{}); err == nil {
			t.Errorf("%s: a Deployment exists for an invalid profile", name)
		}
	}
}

// TestReconcileFollowsAProfileChange covers srd009 R2.5: a ConfigMap maps to
// the Agents naming it, and an edited port default reaches the Deployment
// and Service on the next reconcile.
func TestReconcileFollowsAProfileChange(t *testing.T) {
	agent := collector()
	profile := profileConfigMap()
	r := newReconciler(t, agent, profile)
	reconcile(t, r, agent)
	requests := r.agentsNamingConfigMap(context.Background(), profile)
	if len(requests) != 1 || requests[0].Name != "collector" || requests[0].Namespace != "demo" {
		t.Errorf("requests for the profile ConfigMap = %v", requests)
	}
	if other := r.agentsNamingConfigMap(context.Background(), &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "unrelated", Namespace: "demo"}}); len(other) != 0 {
		t.Errorf("requests for an unrelated ConfigMap = %v", other)
	}
	live := get(t, r, &corev1.ConfigMap{}, "collector-profile")
	live.Data["collector__rest.yaml"] = strings.ReplaceAll(live.Data["collector__rest.yaml"], "COLLECTOR_CONTROL_PORT:-18191", "COLLECTOR_CONTROL_PORT:-19191")
	if err := r.Update(context.Background(), live); err != nil {
		t.Fatal(err)
	}
	reconcile(t, r, agent)
	c := get(t, r, &appsv1.Deployment{}, "collector").Spec.Template.Spec.Containers[0]
	if c.ReadinessProbe.HTTPGet.Port.IntValue() != 19191 {
		t.Errorf("probe port = %d, want the edited default", c.ReadinessProbe.HTTPGet.Port.IntValue())
	}
	if svc := get(t, r, &corev1.Service{}, "collector"); svc.Spec.Ports[0].Port != 19191 {
		t.Errorf("service control port = %d, want 19191", svc.Spec.Ports[0].Port)
	}
}

// TestReconcilePassesArgsServiceAccountAndClaim covers srd009 R1.1 and R2.1:
// declared arguments follow --profile, the service account and pod labels
// reach the pod template with the operator's labels winning a collision, and
// a work volume naming a claim mounts that claim instead of an emptyDir.
func TestReconcilePassesArgsServiceAccountAndClaim(t *testing.T) {
	agent := collector()
	agent.Spec.Args = []string{"--directory", "/work", "--otel-service-name", "collector"}
	agent.Spec.ServiceAccountName = "collector-reader"
	agent.Spec.PodLabels = map[string]string{"app.kubernetes.io/component": "collector", "app.kubernetes.io/instance": "spoofed"}
	agent.Spec.WorkVolumes = []v1alpha1.WorkVolume{{Name: "work", MountPath: "/work", ClaimName: "collector-work"}}
	r := newReconciler(t, agent, profileConfigMap())
	reconcile(t, r, agent)
	dep := get(t, r, &appsv1.Deployment{}, "collector")
	pod := dep.Spec.Template
	if args := strings.Join(pod.Spec.Containers[0].Args, " "); args != "--profile /profiles/collector/profile.yaml --directory /work --otel-service-name collector" {
		t.Errorf("args = %s", args)
	}
	if pod.Spec.ServiceAccountName != "collector-reader" {
		t.Errorf("service account = %q", pod.Spec.ServiceAccountName)
	}
	if pod.Labels["app.kubernetes.io/component"] != "collector" || pod.Labels["app.kubernetes.io/instance"] != "collector" {
		t.Errorf("pod labels = %v, want the component label kept and the operator's instance label winning", pod.Labels)
	}
	if work := pod.Spec.Volumes[1]; work.PersistentVolumeClaim == nil || work.PersistentVolumeClaim.ClaimName != "collector-work" || work.EmptyDir != nil {
		t.Errorf("work volume = %+v, want the named claim", work.VolumeSource)
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
	live.Spec.Replicas = ptr(int32(2))
	if err := r.Update(context.Background(), live); err != nil {
		t.Fatal(err)
	}
	reconcile(t, r, agent)
	dep := get(t, r, &appsv1.Deployment{}, "collector")
	if img := dep.Spec.Template.Spec.Containers[0].Image; img != live.Spec.Image {
		t.Errorf("image = %s, want the new image", img)
	}
	if *dep.Spec.Replicas != 2 {
		t.Errorf("replicas = %d, want the new count", *dep.Spec.Replicas)
	}
}

func ptr[T any](v T) *T { return &v }

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
	bad.Name, bad.Spec.Image = "bad", ""
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
