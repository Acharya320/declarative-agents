// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

// Package profileservers reads what a profile declares about the process that
// will serve it: every REST server with its host and port, the environment
// variables its address reads, which server answers the lifecycle health
// endpoint, and which serves the monitor views. An operator derives a pod's
// ports, readiness probe, monitor label, and bind-host environment from this
// instead of asking an Agent to repeat them (srd009 R1.2, R2.1), and
// Validate applies agent-core's own startup boundary before any pod exists
// (R2.4).
//
// Both read the profile from an fs.FS - a ConfigMap's keys laid out as
// paths, or a directory - with an explicit environment, so nothing they
// report depends on the caller's process. agent-core's loaders read files
// and expand ${VAR:-default} from the process environment, so the tree is
// first written to a temporary directory with every YAML file expanded
// against the given environment; the loaders then find nothing left to
// expand.
package profileservers

import (
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"

	internalload "github.com/Nokia-Bell-Labs/declarative-agents/agent-core/internal/load"
	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/internal/startupcheck"
	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/internal/support/corepath"
	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/internal/support/envexpand"
	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/internal/tools/lifecycle"
	restdef "github.com/Nokia-Bell-Labs/declarative-agents/agent-core/internal/tools/rest/definition"
)

// HealthPath is the lifecycle endpoint whose server is the agent's readiness
// signal (srd009 R1.2).
const HealthPath = "/api/lifecycle/health"

// monitorPrefix is where a monitor server's views live.
const monitorPrefix = "/monitor/"

// Server is one REST server a profile declares, as the given environment
// resolves it.
type Server struct {
	// Name is the server's key in the REST definitions.
	Name string
	// Host and Port are the expanded address.
	Host string
	Port int
	// BindVar and PortVar name the variables the unexpanded address reads
	// for its host and port; empty when that part is literal or the
	// declaring file is outside the tree.
	BindVar string
	PortVar string
	// Health reports an endpoint at HealthPath; Monitor one under /monitor/.
	Health  bool
	Monitor bool
}

// Options is what the caller supplies about the pod a profile will serve.
type Options struct {
	// Env is the environment the pod will see; it wins over the declaration
	// defaults in every ${VAR:-default} reference, and nothing is read from
	// the calling process.
	Env map[string]string
	// CoreRoot is the agent-core installation the profile's /opt/agent-core
	// imports resolve against: the module checkout in development, the
	// directory an image copies agent-core's tools into in a pod.
	CoreRoot string
}

// loadMu serializes loads: the /opt/agent-core mapping is process-wide.
var loadMu sync.Mutex

// Inspection is what a profile declares about its servers.
type Inspection struct {
	// Servers is sorted by name.
	Servers []Server
	// ControlServer is the one server with a health endpoint.
	ControlServer string
	// MonitorServer is the one server with monitor endpoints, or empty.
	MonitorServer string
}

// Inspect loads the profile at profilePath inside fsys with env as the
// process environment and reports its servers. It fails when the profile
// does not load, when a server's port variable has neither a default nor a
// value in env, when two servers answer the health endpoint or serve monitor
// views, or when no server answers the health endpoint.
func Inspect(fsys fs.FS, profilePath string, options Options) (Inspection, error) {
	var result Inspection
	err := load(fsys, profilePath, options, func(dir string, closure *internalload.Closure) error {
		var err error
		result, err = inspect(closure.Rest.Servers, addressTemplates(fsys, dir, closure.Files))
		return err
	})
	return result, err
}

// Validate loads the profile the way Inspect does and applies agent-core's
// startup boundary: machine policy, tool wiring, every reached machine, and
// the declared request sources. It is what `agent --validate-config` runs,
// so the operator and the CLI reject the same profiles.
func Validate(fsys fs.FS, profilePath string, options Options) error {
	return load(fsys, profilePath, options, func(_ string, closure *internalload.Closure) error {
		if _, err := startupcheck.ValidateClosure(closure); err != nil {
			return err
		}
		return lifecycle.ValidateDeclaredRequestSources("", closure.Selected)
	})
}

// load materializes the tree, maps /opt/agent-core for the duration, loads
// the closure, and hands it to use with the directory it was loaded from.
func load(fsys fs.FS, profilePath string, options Options, use func(string, *internalload.Closure) error) error {
	loadMu.Lock()
	defer loadMu.Unlock()
	if options.CoreRoot != "" {
		previous := corepath.InstallRoot()
		corepath.SetInstallRoot(options.CoreRoot)
		defer corepath.SetInstallRoot(previous)
	}
	if err := checkRootAddresses(fsys, profilePath, options.Env); err != nil {
		return err
	}
	dir, cleanup, err := materialize(fsys, options.Env)
	if err != nil {
		return err
	}
	defer cleanup()
	closure, err := internalload.LoadClosure(filepath.Join(dir, filepath.FromSlash(profilePath)), internalload.Options{CoreRoot: options.CoreRoot})
	if err != nil {
		return fmt.Errorf("load profile %s: %w", profilePath, err)
	}
	return use(dir, closure)
}

// materialize writes fsys to a temporary directory, expanding every YAML
// file's ${VAR:-default} references against env, and returns the directory
// with its cleanup.
func materialize(fsys fs.FS, env map[string]string) (string, func(), error) {
	dir, err := os.MkdirTemp("", "profileservers-*")
	if err != nil {
		return "", nil, fmt.Errorf("materialize profile tree: %w", err)
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	lookup := func(name string) (string, bool) {
		value, ok := env[name]
		return value, ok
	}
	err = fs.WalkDir(fsys, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}
		if isYAML(path) {
			data = envexpand.ExpandWith(data, lookup)
		}
		target := filepath.Join(dir, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
	if err != nil {
		cleanup()
		return "", nil, fmt.Errorf("materialize profile tree: %w", err)
	}
	return dir, cleanup, nil
}

func isYAML(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".yaml", ".yml":
		return true
	}
	return false
}

// checkRootAddresses reads the unexpanded addresses in the profile's REST
// definition roots before anything is expanded: a port variable with neither
// a default nor a value would expand to an address the loader cannot parse,
// and the author needs the variable's name, not a YAML error (srd009 R1.2).
func checkRootAddresses(fsys fs.FS, profilePath string, env map[string]string) error {
	raw, err := fs.ReadFile(fsys, profilePath)
	if err != nil {
		return fmt.Errorf("read profile %s: %w", profilePath, err)
	}
	var profile struct {
		RestDefinitions []string `yaml:"rest_definitions"`
	}
	if yaml.Unmarshal(raw, &profile) != nil {
		return nil // the loader reports the shape error with its own context
	}
	for _, rest := range profile.RestDefinitions {
		if strings.HasPrefix(rest, "/") {
			continue // a library root the loader maps
		}
		for name, template := range serverAddresses(fsys, path.Join(path.Dir(profilePath), rest)) {
			_, portPart := splitTemplate(template)
			for _, ref := range markedReferences(portPart) {
				if _, set := env[ref.Name]; !ref.HasDefault && !set {
					return fmt.Errorf("server %s: port variable %s has no default and no value", name, ref.Name)
				}
			}
		}
	}
	return nil
}

// addressTemplates reads the unexpanded address of every server the closure
// declares, from the files of the closure that lie inside fsys. A file that
// is not a REST definition decodes to no servers and is skipped, as is one
// outside the tree (a library root the pod mounts elsewhere).
func addressTemplates(fsys fs.FS, dir string, files []string) map[string]string {
	templates := map[string]string{}
	roots := []string{dir}
	if resolved, err := filepath.EvalSymlinks(dir); err == nil && resolved != dir {
		roots = append(roots, resolved)
	}
	for _, file := range files {
		rel, ok := relativeTo(roots, file)
		if !ok || !isYAML(rel) {
			continue
		}
		for name, address := range serverAddresses(fsys, filepath.ToSlash(rel)) {
			templates[name] = address
		}
	}
	return templates
}

// relativeTo finds file under one of roots, trying the file's resolved path
// too: a temporary directory and the loader's canonical paths may differ by
// a symlink (/var against /private/var).
func relativeTo(roots []string, file string) (string, bool) {
	candidates := []string{file}
	if resolved, err := filepath.EvalSymlinks(file); err == nil && resolved != file {
		candidates = append(candidates, resolved)
	}
	for _, root := range roots {
		for _, candidate := range candidates {
			rel, err := filepath.Rel(root, candidate)
			if err == nil && !strings.HasPrefix(rel, "..") {
				return rel, true
			}
		}
	}
	return "", false
}

// markedRef is the scalar a ${NAME} or ${NAME:-default} reference becomes
// before an unexpanded declaration is decoded: the parser needs a plain
// token where the reference's braces would open a flow mapping.
var markedRef = regexp.MustCompile(`REF(D?)__([A-Za-z_][A-Za-z0-9_]*)__`)

func markReference(ref envexpand.Reference) string {
	if ref.HasDefault {
		return "REFD__" + ref.Name + "__"
	}
	return "REF__" + ref.Name + "__"
}

// markedReferences reads the references markReference wrote into s.
func markedReferences(s string) []envexpand.Reference {
	var refs []envexpand.Reference
	for _, match := range markedRef.FindAllStringSubmatch(s, -1) {
		refs = append(refs, envexpand.Reference{Name: match[2], HasDefault: match[1] == "D"})
	}
	return refs
}

// serverAddresses decodes the servers' unexpanded addresses from one file in
// fsys, loosely: anything that is not a REST definition yields none. The
// addresses come back with their references marked, for markedReferences.
func serverAddresses(fsys fs.FS, name string) map[string]string {
	raw, err := fs.ReadFile(fsys, name)
	if err != nil {
		return nil
	}
	raw = envexpand.Mark(raw, markReference)
	type servers struct {
		Servers map[string]struct {
			Address string `yaml:"address"`
		} `yaml:"servers"`
	}
	var decoded struct {
		Rest    servers `yaml:"rest"` // a declaration unit wraps its definition in rest:
		servers `yaml:",inline"`
	}
	if yaml.Unmarshal(raw, &decoded) != nil {
		return nil
	}
	addresses := map[string]string{}
	for _, declared := range []servers{decoded.Rest, decoded.servers} {
		for server, address := range declared.Servers {
			if address.Address != "" {
				addresses[server] = address.Address
			}
		}
	}
	return addresses
}

// inspect reads the expanded servers beside their unexpanded addresses.
func inspect(servers map[string]restdef.Server, templates map[string]string) (Inspection, error) {
	if len(servers) == 0 {
		return Inspection{}, ErrNoServers
	}
	names := make([]string, 0, len(servers))
	for name := range servers {
		names = append(names, name)
	}
	sort.Strings(names)
	var result Inspection
	for _, name := range names {
		server, err := resolveServer(name, servers[name], templates[name])
		if err != nil {
			return Inspection{}, err
		}
		if server.Health {
			if result.ControlServer != "" {
				return Inspection{}, fmt.Errorf("servers %s and %s both answer GET %s; one is the readiness signal", result.ControlServer, name, HealthPath)
			}
			result.ControlServer = name
		}
		if server.Monitor {
			if result.MonitorServer != "" {
				return Inspection{}, fmt.Errorf("servers %s and %s both serve %s*", result.MonitorServer, name, monitorPrefix)
			}
			result.MonitorServer = name
		}
		result.Servers = append(result.Servers, server)
	}
	if result.ControlServer == "" {
		return Inspection{}, fmt.Errorf("no server answers GET %s, the readiness signal", HealthPath)
	}
	return result, nil
}

func resolveServer(name string, server restdef.Server, template string) (Server, error) {
	bindVar, portVar := addressVars(template)
	host, portText, err := net.SplitHostPort(server.Address)
	if err != nil || portText == "" {
		if portVar != "" {
			return Server{}, fmt.Errorf("server %s: port variable %s has no default and no value", name, portVar)
		}
		return Server{}, fmt.Errorf("server %s: address %q is not host:port", name, server.Address)
	}
	port, err := strconv.Atoi(portText)
	if err != nil || port <= 0 || port > 65535 {
		return Server{}, fmt.Errorf("server %s: port %q is not a port number", name, portText)
	}
	resolved := Server{Name: name, Host: host, Port: port, BindVar: bindVar, PortVar: portVar}
	for _, endpoint := range server.Endpoints {
		if endpoint.Path == HealthPath {
			resolved.Health = true
		}
		if strings.HasPrefix(endpoint.Path, monitorPrefix) {
			resolved.Monitor = true
		}
	}
	return resolved, nil
}

// addressVars names the variables the host and port parts of an unexpanded
// address read; the port is what follows the last colon.
func addressVars(template string) (bindVar, portVar string) {
	hostPart, portPart := splitTemplate(template)
	return firstReference(hostPart), firstReference(portPart)
}

// splitTemplate splits an unexpanded address at its last colon.
func splitTemplate(template string) (hostPart, portPart string) {
	if i := strings.LastIndex(template, ":"); i >= 0 {
		return template[:i], template[i+1:]
	}
	return template, ""
}

func firstReference(s string) string {
	refs := markedReferences(s)
	if len(refs) == 0 {
		return ""
	}
	return refs[0].Name
}

// ErrNoServers reports a profile that declares no REST servers at all.
var ErrNoServers = errors.New("profile declares no REST servers")
