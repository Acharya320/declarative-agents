// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package controller

import (
	"path"

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

// desiredDeployment renders the one Deployment an Agent implies (srd009
// R2.1): the image with --profile at the mounted path, the ConfigMap mounted
// read-only, a readiness probe on the control port, and the declared
// resources, environment, secret references, and work volumes. It writes
// into dep, so CreateOrUpdate keeps what the server owns.
func desiredDeployment(agent *v1alpha1.Agent, dep *appsv1.Deployment) {
	selector := labels(agent)
	replicas := agent.DesiredReplicas()
	container := corev1.Container{
		Name:  containerName,
		Image: agent.Spec.Image,
		Args:  []string{"--profile", path.Join(profileMountPath, agent.Spec.Profile.Path)},
		VolumeMounts: []corev1.VolumeMount{
			{Name: "profile", MountPath: profileMountPath, ReadOnly: true},
		},
		Resources: agent.Spec.Resources,
		Env:       agent.Spec.Env,
		EnvFrom:   agent.Spec.EnvFrom,
	}
	for _, p := range agent.Spec.Ports {
		container.Ports = append(container.Ports, corev1.ContainerPort{Name: p.Name, ContainerPort: p.Port})
	}
	if port := agent.ControlPortNumber(); port > 0 {
		container.ReadinessProbe = &corev1.Probe{
			ProbeHandler:  corev1.ProbeHandler{HTTPGet: &corev1.HTTPGetAction{Path: healthPath, Port: intstr.FromInt32(port)}},
			PeriodSeconds: 5,
		}
	}
	volumes := []corev1.Volume{{
		Name: "profile",
		VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{Name: agent.Spec.Profile.ConfigMap},
		}},
	}}
	for _, w := range agent.Spec.WorkVolumes {
		volumes = append(volumes, corev1.Volume{Name: w.Name, VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}})
		container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{Name: w.Name, MountPath: w.MountPath})
	}
	dep.Labels = selector
	dep.Spec.Replicas = &replicas
	dep.Spec.Selector = &metav1.LabelSelector{MatchLabels: selector}
	dep.Spec.Template.Labels = selector
	dep.Spec.Template.Spec.Containers = []corev1.Container{container}
	dep.Spec.Template.Spec.Volumes = volumes
}
