// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package controller

import (
	"fmt"
	"sort"
	"strings"
	"testing/fstest"

	corev1 "k8s.io/api/core/v1"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/profileservers"
	"github.com/Nokia-Bell-Labs/declarative-agents/operator/api/v1alpha1"
)

// maxPortName is what Kubernetes allows for a port name (an IANA service
// name): 15 characters.
const maxPortName = 15

// profileFS lays a ConfigMap's keys out as the file tree the pod mounts:
// a key spelling a directory with the separator becomes that path.
func profileFS(profile *corev1.ConfigMap) fstest.MapFS {
	tree := fstest.MapFS{}
	for key, data := range profile.Data {
		tree[strings.ReplaceAll(key, keyDirSeparator, "/")] = &fstest.MapFile{Data: []byte(data)}
	}
	for key, data := range profile.BinaryData {
		tree[strings.ReplaceAll(key, keyDirSeparator, "/")] = &fstest.MapFile{Data: data}
	}
	return tree
}

// literalEnv is the Agent's environment as the profile loader sees it: the
// values set inline. A value from a ConfigMap or Secret reference is not
// known here, so a port variable set only that way keeps its declaration
// default in the derived ports.
func literalEnv(agent *v1alpha1.Agent) map[string]string {
	env := make(map[string]string, len(agent.Spec.Env))
	for _, e := range agent.Spec.Env {
		if e.ValueFrom == nil {
			env[e.Name] = e.Value
		}
	}
	return env
}

// inspection is what the reconciler derives from a profile: the servers with
// the port names they get, the control and monitor servers (srd009 R1.2).
type inspection struct {
	profileservers.Inspection
	portNames map[string]string // server name to container port name
}

// port is the port of the named server.
func (i inspection) port(name string) (int, bool) {
	for _, server := range i.Servers {
		if server.Name == name && name != "" {
			return server.Port, true
		}
	}
	return 0, false
}

// inspectProfile validates the profile with agent-core's startup boundary
// and reads its servers (srd009 R1.2, R2.4). coreRoot is where agent-core's
// tools live in the operator's own filesystem, for /opt/agent-core imports.
func inspectProfile(agent *v1alpha1.Agent, profile *corev1.ConfigMap, coreRoot string) (inspection, string) {
	tree := profileFS(profile)
	options := profileservers.Options{Env: literalEnv(agent), CoreRoot: coreRoot}
	if err := profileservers.Validate(tree, agent.Spec.Profile.Path, options); err != nil {
		return inspection{}, err.Error()
	}
	result, err := profileservers.Inspect(tree, agent.Spec.Profile.Path, options)
	if err != nil {
		return inspection{}, err.Error()
	}
	names, problem := portNames(result.Servers)
	if problem != "" {
		return inspection{}, problem
	}
	return inspection{Inspection: result, portNames: names}, ""
}

// portNames derives a container port name from each server name: the part
// after the last underscore, with hyphens for underscores, which is what the
// catalog's profiles already read as the port's role (collector_control ->
// control). Two servers that collapse to one name, or a name Kubernetes
// cannot hold, are a problem the author resolves by naming the servers.
func portNames(servers []profileservers.Server) (map[string]string, string) {
	names := make(map[string]string, len(servers))
	owners := map[string]string{}
	for _, server := range servers {
		name := server.Name
		if i := strings.LastIndex(name, "_"); i >= 0 {
			name = name[i+1:]
		}
		name = strings.ReplaceAll(name, "_", "-")
		if len(name) > maxPortName || name == "" {
			return nil, fmt.Sprintf("server %s: port name %q must be 1 to %d characters", server.Name, name, maxPortName)
		}
		if other, taken := owners[name]; taken {
			return nil, fmt.Sprintf("servers %s and %s both take port name %q", other, server.Name, name)
		}
		owners[name] = server.Name
		names[server.Name] = name
	}
	return names, ""
}

// bindEnv is the environment the pod gains so every server listens on every
// interface: <bind variable>=0.0.0.0 for each server whose bind variable the
// Agent does not set itself (srd009 R2.1).
func bindEnv(agent *v1alpha1.Agent, servers []profileservers.Server) []corev1.EnvVar {
	set := map[string]bool{}
	for _, e := range agent.Spec.Env {
		set[e.Name] = true
	}
	var added []corev1.EnvVar
	seen := map[string]bool{}
	for _, server := range servers {
		if server.BindVar == "" || set[server.BindVar] || seen[server.BindVar] {
			continue
		}
		seen[server.BindVar] = true
		added = append(added, corev1.EnvVar{Name: server.BindVar, Value: "0.0.0.0"})
	}
	sort.Slice(added, func(i, j int) bool { return added[i].Name < added[j].Name })
	return added
}
