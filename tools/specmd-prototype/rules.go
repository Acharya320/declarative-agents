// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package specmd

import (
	_ "embed"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed rules/srd.yaml
var srdRulesYAML []byte

// rules is one document type's declaration: heading grammar, typed fields
// and their shapes, attribute shapes, and prose containment. Diagnostics read
// every constraint from here, so a new document type is a new declaration.
type rules struct {
	DocumentType string `yaml:"document_type"`
	Frontmatter  struct {
		Required          []string `yaml:"required"`
		Allowed           []string `yaml:"allowed"`
		IDMatchesFilename bool     `yaml:"id_matches_filename"`
	} `yaml:"frontmatter"`
	Nodes                      map[string]nodeRule  `yaml:"nodes"`
	Fields                     map[string]fieldRule `yaml:"fields"`
	RequiredFields             []string             `yaml:"required_fields"`
	RequiresRequirementGroups  bool                 `yaml:"requires_requirement_groups"`
	RequiredFieldsExemptStatus []string             `yaml:"required_fields_exempt_status"`
	Attributes                 map[string]string    `yaml:"attributes"`
	CriterionTraces            struct {
		Min       int    `yaml:"min"`
		IDPattern string `yaml:"id_pattern"`
	} `yaml:"criterion_traces"`
	ProseContainment struct {
		Keywords []string `yaml:"keywords"`
		Outside  []string `yaml:"outside"`
	} `yaml:"prose_containment"`

	idPatterns    map[string]*regexp.Regexp
	itemPrefixes  map[string]*regexp.Regexp
	tracePattern  *regexp.Regexp
	keywordRegexp *regexp.Regexp
}

type nodeRule struct {
	Tag             string   `yaml:"tag"`
	Level           int      `yaml:"level"`
	Parent          string   `yaml:"parent"`
	IDPattern       string   `yaml:"id_pattern"`
	IDExtendsParent bool     `yaml:"id_extends_parent"`
	HeadingIsID     bool     `yaml:"heading_is_id"`
	RequiresContent bool     `yaml:"requires_content"`
	Attributes      []string `yaml:"attributes"`
}

type fieldRule struct {
	In         []string `yaml:"in"`
	Shape      string   `yaml:"shape"`
	ItemPrefix string   `yaml:"item_prefix"`
	Columns    []string `yaml:"columns"`
}

// loadRules parses a declaration and compiles its patterns once.
func loadRules(data []byte) (*rules, error) {
	var r rules
	if err := yaml.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("rules: %w", err)
	}
	r.idPatterns = map[string]*regexp.Regexp{}
	for name, node := range r.Nodes {
		if node.IDPattern != "" {
			pattern, err := regexp.Compile(node.IDPattern)
			if err != nil {
				return nil, fmt.Errorf("rules nodes.%s.id_pattern: %w", name, err)
			}
			r.idPatterns[name] = pattern
		}
	}
	r.itemPrefixes = map[string]*regexp.Regexp{}
	for name, field := range r.Fields {
		if field.ItemPrefix != "" {
			pattern, err := regexp.Compile(field.ItemPrefix)
			if err != nil {
				return nil, fmt.Errorf("rules fields.%s.item_prefix: %w", name, err)
			}
			r.itemPrefixes[name] = pattern
		}
	}
	var err error
	if r.tracePattern, err = regexp.Compile(r.CriterionTraces.IDPattern); err != nil {
		return nil, fmt.Errorf("rules criterion_traces.id_pattern: %w", err)
	}
	r.keywordRegexp = keywordPattern(r.ProseContainment.Keywords)
	return &r, nil
}

// keywordPattern matches any keyword as a whole word, longest first so
// "SHALL NOT" reports as itself rather than as "SHALL".
func keywordPattern(keywords []string) *regexp.Regexp {
	sorted := append([]string(nil), keywords...)
	sort.Slice(sorted, func(i, j int) bool { return len(sorted[i]) > len(sorted[j]) })
	alternatives := make([]string, len(sorted))
	for i, keyword := range sorted {
		alternatives[i] = regexp.QuoteMeta(keyword)
	}
	return regexp.MustCompile(`\b(` + strings.Join(alternatives, "|") + `)\b`)
}

// nodeFor finds the node a marker tag opens at a heading depth under parent.
func (r *rules) nodeFor(tag string, level int, parent string) (string, bool) {
	for name, node := range r.Nodes {
		if node.Tag == tag && node.Level == level && node.Parent == parent {
			return name, true
		}
	}
	return "", false
}

// allowedHere lists the "<tag> at level <n>" headings legal under parent,
// for messages.
func (r *rules) allowedHere(parent string) []string {
	var allowed []string
	for _, node := range r.Nodes {
		if node.Parent == parent {
			allowed = append(allowed, fmt.Sprintf("%s at level %d", node.Tag, node.Level))
		}
	}
	sort.Strings(allowed)
	return allowed
}

var srdRules = func() *rules {
	r, err := loadRules(srdRulesYAML)
	if err != nil {
		panic(err)
	}
	return r
}()
