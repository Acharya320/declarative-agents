// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package controller

import (
	"path"
	"sort"
	"strconv"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/Nokia-Bell-Labs/declarative-agents/operator/api/v1alpha1"
)

// Where the profile ConfigMap mounts in the agent container, and the health
// path every agent-core profile answers on its control port (srd009 R1.2).
const (
	profileMountPath = "/profiles"
	healthPath       = "/api/lifecycle/health"
	containerName    = "agent"
)

// labels are the selector labels an Agent's Deployment and Service share.
func labels(agent *v1alpha1.Agent) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "agent",
		"app.kubernetes.io/instance":   agent.Name,
		"app.kubernetes.io/managed-by": "agent-operator",
	}
}

// keyDirSeparator stands for "/" in a profile ConfigMap key: keys cannot
// hold a slash, while profiles import sibling directories
// (../units/types.yaml), so units__types.yaml mounts at units/types.yaml.
const keyDirSeparator = "__"

// monitorPortLabel is the pod label an observer's monitor fan-in reads to
// find the monitor server's port (srd009 R2.1).
const monitorPortLabel = "monitorPort"

// profileItems maps a ConfigMap's keys onto mount paths when any key spells a
// directory with the separator. A ConfigMap of flat keys returns nil, which
// mounts every key at its own name and keeps following key additions
// without a reconcile.
func profileItems(profile *corev1.ConfigMap) []corev1.KeyToPath {
	nested := false
	keys := make([]string, 0, len(profile.Data)+len(profile.BinaryData))
	for key := range profile.Data {
		keys = append(keys, key)
	}
	for key := range profile.BinaryData {
		keys = append(keys, key)
	}
	for _, key := range keys {
		nested = nested || strings.Contains(key, keyDirSeparator)
	}
	if !nested {
		return nil
	}
	sort.Strings(keys)
	items := make([]corev1.KeyToPath, 0, len(keys))
	for _, key := range keys {
		items = append(items, corev1.KeyToPath{Key: key, Path: strings.ReplaceAll(key, keyDirSeparator, "/")})
	}
	return items
}

// podLabels joins the Agent's declared pod labels with the operator's own,
// which win on a collision: the selector must keep matching the pods
// whatever a chart declares (srd009 R2.1).
func podLabels(agent *v1alpha1.Agent) map[string]string {
	merged := make(map[string]string, len(agent.Spec.PodLabels)+3)
	for k, v := range agent.Spec.PodLabels {
		merged[k] = v
	}
	for k, v := range labels(agent) {
		merged[k] = v
	}
	return merged
}

// workVolume renders one declared work volume: the named claim when one is
// given, an emptyDir otherwise.
func workVolume(w v1alpha1.WorkVolume) corev1.Volume {
	if w.ClaimName != "" {
		return corev1.Volume{Name: w.Name, VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: w.ClaimName},
		}}
	}
	return corev1.Volume{Name: w.Name, VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}}
}

// desiredDeployment renders the one Deployment an Agent implies (srd009
// R2.1): the image with --profile at the mounted path followed by the
// declared arguments, the ConfigMap mounted read-only with its keys laid out
// as profileItems says, a readiness probe on the control port, and the
// declared resources, environment, secret references, service account, pod
// labels, and work volumes. It writes into dep, so CreateOrUpdate keeps
// what the server owns.
func desiredDeployment(agent *v1alpha1.Agent, profile *corev1.ConfigMap, insp inspection, dep *appsv1.Deployment) {
	selector := labels(agent)
	replicas := agent.DesiredReplicas()
	container := corev1.Container{
		Name:  containerName,
		Image: agent.Spec.Image,
		Args:  append([]string{"--profile", path.Join(profileMountPath, agent.Spec.Profile.Path)}, agent.Spec.Args...),
		VolumeMounts: []corev1.VolumeMount{
			{Name: "profile", MountPath: profileMountPath, ReadOnly: true},
		},
		Resources: agent.Spec.Resources,
		Env:       agent.Spec.Env,
		EnvFrom:   agent.Spec.EnvFrom,
	}
	for _, server := range insp.Servers {
		port := int32(server.Port)
		container.Ports = append(container.Ports, corev1.ContainerPort{Name: insp.portNames[server.Name], ContainerPort: port})
		if server.Name == insp.ControlServer {
			container.ReadinessProbe = &corev1.Probe{
				ProbeHandler:  corev1.ProbeHandler{HTTPGet: &corev1.HTTPGetAction{Path: healthPath, Port: intstr.FromInt32(port)}},
				PeriodSeconds: 5,
			}
		}
	}
	container.Env = append(append([]corev1.EnvVar(nil), agent.Spec.Env...), bindEnv(agent, insp.Servers)...)
	pod := podLabels(agent)
	if port, ok := insp.port(insp.MonitorServer); ok {
		pod[monitorPortLabel] = strconv.Itoa(port)
	}
	volumes := []corev1.Volume{{
		Name: "profile",
		VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{Name: agent.Spec.Profile.ConfigMap},
			Items:                profileItems(profile),
		}},
	}}
	for _, w := range agent.Spec.WorkVolumes {
		volumes = append(volumes, workVolume(w))
		container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{Name: w.Name, MountPath: w.MountPath})
	}
	dep.Labels = selector
	dep.Spec.Replicas = &replicas
	dep.Spec.Selector = &metav1.LabelSelector{MatchLabels: selector}
	dep.Spec.Template.Labels = pod
	dep.Spec.Template.Spec.ServiceAccountName = agent.Spec.ServiceAccountName
	dep.Spec.Template.Spec.Containers = []corev1.Container{container}
	dep.Spec.Template.Spec.Volumes = volumes
}
