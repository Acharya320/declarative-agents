// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec"
	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec/markdown"
	"gopkg.in/yaml.v3"
)

// typedModels parses a document of a type that has a typed model; the
// round-trip laws require the same model from either front end.
var typedModels = map[string]func(path string) (any, error){
	"use_case":       func(path string) (any, error) { return spec.ParseUseCase(path) },
	"test_suite":     func(path string) (any, error) { return spec.ParseTestSuite(path) },
	"roadmap":        func(path string) (any, error) { return spec.ParseRoadmap(path) },
	"specifications": func(path string) (any, error) { return spec.ParseSpecIndex(path) },
}

func repositoryDeclarations(t *testing.T) markdown.Declarations {
	t.Helper()
	data, err := os.ReadFile(designPath)
	if err != nil {
		t.Fatal(err)
	}
	declarations, err := markdown.LoadDeclarations(data)
	if err != nil {
		t.Fatal(err)
	}
	return declarations
}

// trackedDocuments lists every tracked YAML document of a declared type
// other than the SRD, which TestEveryCorpusSRDRoundTrips covers.
func trackedDocuments(t *testing.T, declarations markdown.Declarations) map[string]*markdown.Declaration {
	t.Helper()
	out, err := exec.Command("git", "-C", repositoryRoot, "ls-files", "*.yaml").Output()
	if err != nil {
		t.Fatal(err)
	}
	documents := map[string]*markdown.Declaration{}
	for _, path := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if strings.Contains(path, "/testdata/") {
			continue
		}
		if d, ok := declarations.ForSource(path); ok && d.Type() != "srd" {
			documents[path] = d
		}
	}
	return documents
}

// TestEveryCorpusDocumentRoundTrips states the round-trip laws on every
// tracked document of every structured type besides the SRD (eng10).
func TestEveryCorpusDocumentRoundTrips(t *testing.T) {
	documents := trackedDocuments(t, repositoryDeclarations(t))
	seen := map[string]int{}
	for path, d := range documents {
		seen[d.Type()]++
		t.Run(path, func(t *testing.T) {
			documentRoundTrip(t, d, filepath.Join(repositoryRoot, path), path)
		})
	}
	for _, name := range []string{"vision", "architecture", "specifications", "roadmap", "pattern_language",
		"engineering_guideline", "use_case", "test_suite", "audit_register", "semantic_model",
		"config_format", "constitution", "migration"} {
		if seen[name] == 0 {
			t.Errorf("no tracked document of type %s", name)
		}
	}
}

func documentRoundTrip(t *testing.T, d *markdown.Declaration, file, path string) {
	canonical, _ := d.Canonicalize(readYAML(t, file))
	name := strings.TrimSuffix(path, ".yaml") + ".md"
	md, err := d.Print(canonical)
	if err != nil {
		t.Fatal(err)
	}
	back, findings, err := d.Parse(md, name)
	if err != nil || back == nil {
		t.Fatalf("parse: %v %v", err, findings)
	}
	for _, f := range findings {
		if f.Name != "prose-containment" {
			t.Errorf("printed form has %s", f)
		}
	}
	if diff := markdown.Diff(canonical, back); diff != "" {
		t.Errorf("yaml -> md -> yaml changed %s", diff)
	}
	if formatted, _, err := d.Format(md, name); err != nil || string(formatted) != string(md) {
		t.Errorf("printed Markdown is not a fixed point: %v", err)
	}
	if parse, ok := typedModels[d.Type()]; ok {
		compareModels(t, parse, file, back)
	}
}

// compareModels parses the original YAML and the parsed tree with the type's
// typed parser and requires equal models, strings compared without their
// trailing newlines. A document the typed parser rejects must be rejected
// the same way from either form.
func compareModels(t *testing.T, parse func(string) (any, error), file string, back *yaml.Node) {
	t.Helper()
	original, originalErr := parse(file)
	data, err := yaml.Marshal(back)
	if err != nil {
		t.Fatal(err)
	}
	rewritten := filepath.Join(t.TempDir(), filepath.Base(file))
	if err := os.WriteFile(rewritten, data, 0o600); err != nil {
		t.Fatal(err)
	}
	fromMarkdown, err := parse(rewritten)
	if originalErr != nil || err != nil {
		if withoutPath(originalErr, file) != withoutPath(err, rewritten) {
			t.Errorf("typed parser disagrees: %v; from Markdown: %v", originalErr, err)
		}
		return
	}
	if !reflect.DeepEqual(trimmedJSON(t, original), trimmedJSON(t, fromMarkdown)) {
		t.Errorf("typed model changed:\n%v\n%v", trimmedJSON(t, original), trimmedJSON(t, fromMarkdown))
	}
}

func trimmedJSON(t *testing.T, model any) any {
	t.Helper()
	data, err := json.Marshal(model)
	if err != nil {
		t.Fatal(err)
	}
	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatal(err)
	}
	return trimStrings(value)
}

func trimStrings(value any) any {
	switch v := value.(type) {
	case string:
		return strings.TrimRight(v, "\n")
	case []any:
		for i := range v {
			v[i] = trimStrings(v[i])
		}
	case map[string]any:
		for k := range v {
			v[k] = trimStrings(v[k])
		}
	}
	return value
}

func withoutPath(err error, path string) string {
	if err == nil {
		return ""
	}
	return strings.ReplaceAll(err.Error(), path, "<path>")
}
