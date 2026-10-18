// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package v1alpha1

import "k8s.io/apimachinery/pkg/runtime"

// DeepCopyObject implements runtime.Object.
func (a *Agent) DeepCopyObject() runtime.Object { return a.DeepCopy() }

// DeepCopy returns a deep copy of the Agent.
func (a *Agent) DeepCopy() *Agent {
	if a == nil {
		return nil
	}
	out := *a
	a.DeepCopyInto(&out.ObjectMeta)
	out.Spec = a.Spec.DeepCopy()
	out.Status = a.Status.DeepCopy()
	return &out
}

// DeepCopy returns a deep copy of the spec.
func (s AgentSpec) DeepCopy() AgentSpec {
	out := s
	if s.Replicas != nil {
		r := *s.Replicas
		out.Replicas = &r
	}
	out.WorkVolumes = append([]WorkVolume(nil), s.WorkVolumes...)
	out.Args = append([]string(nil), s.Args...)
	out.PodLabels = nil
	if s.PodLabels != nil {
		out.PodLabels = make(map[string]string, len(s.PodLabels))
		for k, v := range s.PodLabels {
			out.PodLabels[k] = v
		}
	}
	out.Resources = *s.Resources.DeepCopy()
	out.Env = nil
	for _, e := range s.Env {
		out.Env = append(out.Env, *e.DeepCopy())
	}
	out.EnvFrom = nil
	for _, e := range s.EnvFrom {
		out.EnvFrom = append(out.EnvFrom, *e.DeepCopy())
	}
	return out
}

// DeepCopy returns a deep copy of the status.
func (s AgentStatus) DeepCopy() AgentStatus {
	out := s
	out.Conditions = nil
	for _, c := range s.Conditions {
		out.Conditions = append(out.Conditions, *c.DeepCopy())
	}
	return out
}

// DeepCopyObject implements runtime.Object.
func (l *AgentList) DeepCopyObject() runtime.Object {
	if l == nil {
		return nil
	}
	out := *l
	l.DeepCopyInto(&out.ListMeta)
	out.Items = nil
	for i := range l.Items {
		out.Items = append(out.Items, *l.Items[i].DeepCopy())
	}
	return &out
}
