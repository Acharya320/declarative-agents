// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec/markdown"
	"gopkg.in/yaml.v3"
)

// Audit validates the application documentation, canonical profile boot
// closure, and every formal Go-test evidence claim.
func Audit() error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := auditDocumentation(root); err != nil {
		return err
	}
	roots, err := resolveIntegrationRoots()
	if err != nil {
		return err
	}
	manifest, err := readApplicationProfileManifest(filepath.Join(root, filepath.FromSlash(profileManifestPath)))
	if err != nil {
		return err
	}
	packagedRoot, err := os.MkdirTemp("", "coding-agent-profiles-*")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(packagedRoot) }()
	packagedRoot = filepath.Join(packagedRoot, "profiles")
	source, err := inspectPackageSource(roots.Profiles, manifest.Catalog.CompatibleRelease)
	if err != nil {
		return err
	}
	if _, err := assembleProfileClosure(manifest, roots.Profiles, packagedRoot, source); err != nil {
		return fmt.Errorf("assemble canonical profile closure: %w", err)
	}
	deploymentRoot := filepath.Join(filepath.Dir(packagedRoot), "deployment")
	shards, err := packageServingDeployment(root, roots.Profiles, deploymentRoot, manifest, source)
	if err != nil {
		return fmt.Errorf("package serving deployment: %w", err)
	}
	binary, cleanup, err := buildAgent(roots.Core)
	if err != nil {
		return err
	}
	defer cleanup()
	profiles := make([]string, 0, len(manifest.Catalog.References))
	for _, ref := range manifest.Catalog.References {
		profiles = append(profiles, filepath.Join(packagedRoot, filepath.FromSlash(ref.RuntimePath)))
	}
	for _, shard := range shards {
		profiles = append(profiles,
			filepath.Join(deploymentRoot, shard.Path, filepath.FromSlash(shard.Profile)))
	}
	if err := bootSmokeProfiles(binary, roots.Core, profiles); err != nil {
		return err
	}
	if err := runCorpusAudit(binary, root, roots.Core, roots.Profiles); err != nil {
		return err
	}
	return runTestEvidenceAudit(binary, root, roots.Core, roots.Profiles)
}

// codingAgentRequiredFields lists the top-level fields each foundation
// document carries.
var codingAgentRequiredFields = map[string][]string{
	"docs/VISION.md": {"id", "title", "executive_summary", "problem", "what_this_does",
		"why_we_build_this", "success_criteria", "not"},
	"docs/ARCHITECTURE.md": {"id", "title", "overview", "interfaces", "components",
		"design_decisions", "technology_choices", "project_structure", "implementation_status", "related_documents"},
	"docs/road-map.md": {"id", "title", "releases"},
	"docs/SPECIFICATIONS.md": {"id", "title", "overview", "roadmap_summary", "foundation_document_index",
		"srd_index", "config_format_index", "semantic_model_index", "use_case_index", "test_suite_index", "coverage_gaps"},
}

// auditDocumentation checks the coding-agent corpus: the foundation
// documents carry their fields, every indexed path exists, every use case
// names a test suite that traces it, and the coding-loop suite defines its
// stages in order. Corpus documents are Markdown, read as their YAML tree
// (eng10).
func auditDocumentation(root string) error {
	parsed, errs := parseEveryDocument(root)
	for path, fields := range codingAgentRequiredFields {
		document, err := readCodingAgentDocument(root, path)
		if err != nil {
			return err
		}
		if missing := missingFields(document, fields); len(missing) > 0 {
			errs = append(errs, fmt.Sprintf("%s missing %v", path, missing))
		}
	}
	index, err := readCodingAgentDocument(root, "docs/SPECIFICATIONS.md")
	if err != nil {
		return err
	}
	errs = append(errs, indexPathErrors(root, index)...)
	traceErrs, err := traceabilityErrors(root, index)
	if err != nil {
		return err
	}
	errs = append(errs, traceErrs...)
	if len(errs) > 0 {
		return fmt.Errorf("audit coding-agent docs:\n%s", strings.Join(errs, "\n"))
	}
	fmt.Printf("audit: parsed and validated %d coding-agent documents\n", parsed)
	return nil
}

// parseEveryDocument parses every Markdown document at a declared corpus
// location and every YAML file under docs (the corpus charter), and returns
// how many parsed and the errors of those that did not.
func parseEveryDocument(root string) (int, []string) {
	parsed := 0
	var errs []string
	err := filepath.WalkDir(filepath.Join(root, "docs"), func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !(strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".md")) {
			return err
		}
		if strings.HasSuffix(path, ".md") && !declaredMarkdown(path) {
			return nil // native Markdown such as a guide
		}
		data, readErr := readDocsFile(path)
		var document any
		if readErr == nil {
			readErr = yaml.Unmarshal(data, &document)
		}
		if readErr != nil {
			errs = append(errs, fmt.Sprintf("%s does not parse: %v", path, readErr))
			return nil
		}
		parsed++
		return nil
	})
	if err != nil {
		errs = append(errs, err.Error())
	}
	return parsed, errs
}

func readCodingAgentDocument(root, path string) (map[string]any, error) {
	data, err := markdown.ReadAsYAML(filepath.Join(root, filepath.FromSlash(path)))
	if err != nil {
		return nil, fmt.Errorf("audit coding-agent docs: read %s: %w", path, err)
	}
	var document map[string]any
	if err := yaml.Unmarshal(data, &document); err != nil {
		return nil, fmt.Errorf("audit coding-agent docs: parse %s: %w", path, err)
	}
	return document, nil
}

func missingFields(document map[string]any, fields []string) []string {
	var missing []string
	for _, field := range fields {
		if _, ok := document[field]; !ok {
			missing = append(missing, field)
		}
	}
	sort.Strings(missing)
	return missing
}

// indexEntries returns the entries of one specification-index section.
func indexEntries(index map[string]any, section string) []map[string]any {
	list, _ := index[section].([]any)
	entries := make([]map[string]any, 0, len(list))
	for _, item := range list {
		if entry, ok := item.(map[string]any); ok {
			entries = append(entries, entry)
		}
	}
	return entries
}

func indexPathErrors(root string, index map[string]any) []string {
	var errs []string
	for _, section := range []string{"foundation_document_index", "srd_index", "use_case_index", "test_suite_index"} {
		for _, entry := range indexEntries(index, section) {
			path, _ := entry["path"].(string)
			if !indexedPathExists(filepath.Join(root, filepath.FromSlash(path))) {
				errs = append(errs, fmt.Sprintf("%s path does not exist: %s", section, path))
			}
		}
	}
	return errs
}

// indexedPathExists reports an indexed file.
func indexedPathExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// readDocsFile reads a Markdown corpus document as its YAML tree and a YAML
// file as written.
func readDocsFile(path string) ([]byte, error) {
	if strings.HasSuffix(path, ".md") {
		return markdown.ReadAsYAML(path)
	}
	return os.ReadFile(path)
}

// traceabilityErrors requires every indexed use case to name an indexed test
// suite that traces it back, and the coding-loop suite to stage A, B, C.
func traceabilityErrors(root string, index map[string]any) ([]string, error) {
	suites := map[string]map[string]any{}
	for _, entry := range indexEntries(index, "test_suite_index") {
		suite, err := readCodingAgentDocument(root, fmt.Sprint(entry["path"]))
		if err != nil {
			return nil, err
		}
		suites[fmt.Sprint(entry["id"])] = suite
	}
	var errs []string
	for _, entry := range indexEntries(index, "use_case_index") {
		useCase, err := readCodingAgentDocument(root, fmt.Sprint(entry["path"]))
		if err != nil {
			return nil, err
		}
		suiteID := fmt.Sprint(useCase["test_suite"])
		suite, ok := suites[suiteID]
		switch {
		case !ok:
			errs = append(errs, fmt.Sprintf("use case %v names unknown test suite %s", entry["id"], suiteID))
		case !containsValue(suite["traces"], fmt.Sprint(entry["id"])):
			errs = append(errs, fmt.Sprintf("test suite %s does not trace use case %v", suiteID, entry["id"]))
		}
	}
	stageErrs, err := codingLoopStageErrors(root)
	return append(errs, stageErrs...), err
}

func codingLoopStageErrors(root string) ([]string, error) {
	useCase, err := readCodingAgentDocument(root, "docs/specs/use-cases/rel01.0-uc001-coding-loop.md")
	if err != nil {
		return nil, err
	}
	suite, err := readCodingAgentDocument(root, "docs/specs/test-suites/test-rel01.0-coding-loop.md")
	if err != nil {
		return nil, err
	}
	var errs []string
	if fmt.Sprint(useCase["test_suite"]) != fmt.Sprint(suite["id"]) {
		errs = append(errs, "use case does not name the coding-loop test suite")
	}
	if !containsValue(suite["traces"], fmt.Sprint(useCase["id"])) {
		errs = append(errs, "test suite does not trace the coding-loop use case")
	}
	var stages []string
	cases, _ := suite["test_cases"].([]any)
	for _, item := range cases {
		if testCase, ok := item.(map[string]any); ok {
			stages = append(stages, fmt.Sprint(testCase["stage"]))
		}
	}
	if strings.Join(stages, ",") != "A,B,C" {
		errs = append(errs, "test suite must define ordered stages A, B, and C")
	}
	return errs, nil
}

func containsValue(list any, want string) bool {
	items, _ := list.([]any)
	for _, item := range items {
		if fmt.Sprint(item) == want {
			return true
		}
	}
	return false
}

func bootSmokeProfiles(binary, coreRoot string, profiles []string) error {
	var failures []string
	for _, profile := range profiles {
		if err := runAgentPreflight(binary, "--validate-config", "--profile", profile, "--core-root", coreRoot); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", profile, err))
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("canonical profile boot smoke failed:\n%s", strings.Join(failures, "\n"))
	}
	fmt.Printf("boot smoke passed for %d coding profiles\n", len(profiles))
	return nil
}

func runAgentPreflight(binary string, args ...string) error {
	cmd := exec.Command(binary, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		detail := strings.TrimSpace(string(output))
		if detail == "" {
			detail = err.Error()
		}
		return fmt.Errorf("%s", detail)
	}
	return nil
}

// declaredMarkdown reports a Markdown file at a declared corpus location.
func declaredMarkdown(path string) bool {
	declarations, err := markdown.Default()
	if err != nil {
		return false
	}
	_, ok := declarations.ForMarkdown(filepath.ToSlash(path))
	return ok
}
