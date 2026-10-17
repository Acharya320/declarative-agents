// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

// Package v1alpha1 holds the Agent resource (applications srd009 R1): one
// namespaced object that declares an agent from a profile ConfigMap and the
// agent-core image, which the operator reconciles into a Deployment and a
// Service.
package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Group, Version, and the kinds this package serves.
const (
	Group   = "declarative-agents.nokia.com"
	Version = "v1alpha1"
	Kind    = "Agent"
	Plural  = "agents"
)

// GroupVersion is the API group and version of the Agent resource.
var GroupVersion = schema.GroupVersion{Group: Group, Version: Version}

// AddToScheme registers Agent and AgentList with a scheme.
func AddToScheme(s *runtime.Scheme) error {
	s.AddKnownTypes(GroupVersion, &Agent{}, &AgentList{})
	metav1.AddToGroupVersion(s, GroupVersion)
	return nil
}

// ProfileSource names the ConfigMap holding the agent's profile directory
// and the path of the profile file inside it (srd009 R1.1). A key holding
// "__" mounts at the directory path it spells, so a ConfigMap can carry the
// sibling directories a profile imports (R2.1).
type ProfileSource struct {
	ConfigMap string `json:"configMap"`
	Path      string `json:"path"`
}

// AgentPort is one port the profile's REST servers listen on.
type AgentPort struct {
	Name string `json:"name"`
	Port int32  `json:"port"`
}

// WorkVolume is an emptyDir mounted for the agent's working files.
type WorkVolume struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
}

// AgentSpec declares the agent (srd009 R1.1).
type AgentSpec struct {
	Profile ProfileSource `json:"profile"`
	Image   string        `json:"image"`
	Ports   []AgentPort   `json:"ports,omitempty"`
	// ControlPort names the port in Ports that answers GET
	// /api/lifecycle/health, the readiness signal (srd009 R1.2).
	ControlPort string                      `json:"controlPort"`
	Replicas    *int32                      `json:"replicas,omitempty"`
	Resources   corev1.ResourceRequirements `json:"resources,omitempty"`
	Env         []corev1.EnvVar             `json:"env,omitempty"`
	EnvFrom     []corev1.EnvFromSource      `json:"envFrom,omitempty"`
	WorkVolumes []WorkVolume                `json:"workVolumes,omitempty"`
}

// ConditionReady is the one condition the operator keeps (srd009 R3.1).
const ConditionReady = "Ready"

// AgentStatus reports what the operator observed (srd009 R3.1).
type AgentStatus struct {
	ObservedGeneration int64              `json:"observedGeneration,omitempty"`
	ReadyReplicas      int32              `json:"readyReplicas,omitempty"`
	Conditions         []metav1.Condition `json:"conditions,omitempty"`
}

// Agent is the declared agent.
type Agent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AgentSpec   `json:"spec"`
	Status            AgentStatus `json:"status,omitempty"`
}

// AgentList is a list of Agents.
type AgentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Agent `json:"items"`
}

// DesiredReplicas is the replica count, one when unset.
func (a *Agent) DesiredReplicas() int32 {
	if a.Spec.Replicas == nil {
		return 1
	}
	return *a.Spec.Replicas
}

// ControlPortNumber returns the number of the control port, or 0 when no
// declared port carries that name.
func (a *Agent) ControlPortNumber() int32 {
	for _, p := range a.Spec.Ports {
		if p.Name == a.Spec.ControlPort {
			return p.Port
		}
	}
	return 0
}
