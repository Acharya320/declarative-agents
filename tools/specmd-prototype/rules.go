// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package specmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// designPath is where the document-type declarations live, relative to the
// repository root.
const designPath = "agent-core/docs/constitutions/design.yaml"

// rules is one document type's Markdown declaration, read from the markdown
// block of its entry in design.yaml's document_types: heading grammar, typed
// fields and their shapes, attribute shapes, and prose containment.
// Diagnostics read every constraint from here, so a new document type is a
// new declaration, not new code.
type rules struct {
	Location    string `yaml:"location"`
	Frontmatter struct {
		Required          []string `yaml:"required"`
		Allowed           []string `yaml:"allowed"`
		IDMatchesFilename bool     `yaml:"id_matches_filename"`
	} `yaml:"frontmatter"`
	Nodes            map[string]nodeRule  `yaml:"nodes"`
	Fields           map[string]fieldRule `yaml:"fields"`
	RequiredFields   []string             `yaml:"required_fields"`
	RequiredNodes    map[string]int       `yaml:"required_nodes"`
	ExemptStatus     []string             `yaml:"exempt_status"`
	Attributes       map[string]string    `yaml:"attributes"`
	ProseContainment struct {
		Keywords []string `yaml:"keywords"`
		Outside  []string `yaml:"outside"`
	} `yaml:"prose_containment"`

	documentType  string
	filePattern   *regexp.Regexp
	idPatterns    map[string]*regexp.Regexp
	itemPatterns  map[string]*regexp.Regexp
	keywordRegexp *regexp.Regexp
}

// Node heading forms and content kinds a declaration may name.
const (
	headingID      = "id"       // the heading is exactly the marker's id
	headingIDTitle = "id_title" // "<id>" or "<id> - <title>"
	contentNone    = "none"     // the node holds only tagged children
	contentProse   = "prose"    // required text, no untagged table or list
)

type nodeRule struct {
	Tag             string   `yaml:"tag"`
	Level           int      `yaml:"level"`
	Parent          string   `yaml:"parent"`
	IDPattern       string   `yaml:"id_pattern"`
	IDExtendsParent bool     `yaml:"id_extends_parent"`
	Heading         string   `yaml:"heading"`
	Content         string   `yaml:"content"`
	UniqueID        bool     `yaml:"unique_id"`
	Attributes      []string `yaml:"attributes"`
	RequiredFields  []string `yaml:"required_fields"`
}

type fieldRule struct {
	In          []string `yaml:"in"`
	Shape       string   `yaml:"shape"`
	ItemPattern string   `yaml:"item_pattern"`
	Columns     []string `yaml:"columns"`
}

// loadDeclarations reads every document type in design.yaml whose markdown
// entry is a declaration (a mapping); other markdown values, such as a note
// that a type stays YAML, declare no Markdown rules.
func loadDeclarations(data []byte) (map[string]*rules, error) {
	var design struct {
		DocumentTypes map[string]struct {
			Markdown yaml.Node `yaml:"markdown"`
		} `yaml:"document_types"`
	}
	if err := yaml.Unmarshal(data, &design); err != nil {
		return nil, fmt.Errorf("design.yaml: %w", err)
	}
	declarations := map[string]*rules{}
	for name, entry := range design.DocumentTypes {
		if entry.Markdown.Kind != yaml.MappingNode {
			continue
		}
		var r rules
		if err := entry.Markdown.Decode(&r); err != nil {
			return nil, fmt.Errorf("document_types.%s.markdown: %w", name, err)
		}
		r.documentType = name
		if err := r.compile(); err != nil {
			return nil, fmt.Errorf("document_types.%s.markdown.%w", name, err)
		}
		declarations[name] = &r
	}
	return declarations, nil
}

// compile prepares a declaration's patterns once.
func (r *rules) compile() error {
	pattern, err := locationFilePattern(filepath.Base(r.Location))
	if err != nil {
		return fmt.Errorf("location: %w", err)
	}
	r.filePattern = pattern
	r.idPatterns = map[string]*regexp.Regexp{}
	for name, node := range r.Nodes {
		if node.IDPattern == "" {
			continue
		}
		if r.idPatterns[name], err = regexp.Compile(node.IDPattern); err != nil {
			return fmt.Errorf("nodes.%s.id_pattern: %w", name, err)
		}
	}
	r.itemPatterns = map[string]*regexp.Regexp{}
	for name, field := range r.Fields {
		if field.ItemPattern == "" {
			continue
		}
		if r.itemPatterns[name], err = regexp.Compile(field.ItemPattern); err != nil {
			return fmt.Errorf("fields.%s.item_pattern: %w", name, err)
		}
	}
	r.keywordRegexp = keywordPattern(r.ProseContainment.Keywords)
	return nil
}

// locationFilePattern turns a declared file name such as
// "srdNNN-short-name.md" into a pattern over real file names: a run of N is
// that many digits, short-name is kebab-case, and name is any lowercase name.
func locationFilePattern(base string) (*regexp.Regexp, error) {
	if base == "" || base == "." {
		return nil, errors.New("no file name")
	}
	var pattern strings.Builder
	pattern.WriteString("^")
	for i := 0; i < len(base); {
		switch {
		case strings.HasPrefix(base[i:], "-short-name"):
			pattern.WriteString(`(?:-[a-z0-9]+(?:-[a-z0-9]+)*)?`)
			i += len("-short-name")
		case strings.HasPrefix(base[i:], "short-name"):
			pattern.WriteString(`[a-z0-9]+(?:-[a-z0-9]+)*`)
			i += len("short-name")
		case strings.HasPrefix(base[i:], "name"):
			pattern.WriteString(`[a-z0-9][a-z0-9._-]*`)
			i += len("name")
		case base[i] == 'N':
			end := i
			for end < len(base) && base[end] == 'N' {
				end++
			}
			fmt.Fprintf(&pattern, `[0-9]{%d}`, end-i)
			i = end
		default:
			pattern.WriteString(regexp.QuoteMeta(base[i : i+1]))
			i++
		}
	}
	pattern.WriteString("$")
	return regexp.Compile(pattern.String())
}

// keywordPattern matches any keyword as a whole word, longest first so
// "SHALL NOT" reports as itself rather than as "SHALL".
func keywordPattern(keywords []string) *regexp.Regexp {
	if len(keywords) == 0 {
		return nil
	}
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
	if len(allowed) == 0 {
		return []string{"no tagged heading"}
	}
	return allowed
}

// tags lists the marker tags a declaration uses, sorted.
func (r *rules) tags() []string {
	seen := map[string]bool{}
	var tags []string
	for _, node := range r.Nodes {
		if !seen[node.Tag] {
			seen[node.Tag] = true
			tags = append(tags, node.Tag)
		}
	}
	sort.Strings(tags)
	return tags
}

// declarationFor picks the declaration whose location names this file.
func declarationFor(declarations map[string]*rules, source string) (*rules, bool) {
	base := filepath.Base(source)
	names := make([]string, 0, len(declarations))
	for name := range declarations {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if declarations[name].filePattern.MatchString(base) {
			return declarations[name], true
		}
	}
	return nil, false
}

// findDesign walks up from dir to the repository directory holding
// design.yaml.
func findDesign(dir string) (string, error) {
	for {
		candidate := filepath.Join(dir, designPath)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("%s not found above the working directory", designPath)
		}
		dir = parent
	}
}

var (
	defaultOnce         sync.Once
	defaultDeclarations map[string]*rules
	defaultErr          error
)

// declarationsFromRepository loads the declarations once from the design.yaml
// of the repository enclosing the working directory.
func declarationsFromRepository() (map[string]*rules, error) {
	defaultOnce.Do(func() {
		wd, err := os.Getwd()
		if err != nil {
			defaultErr = err
			return
		}
		path, err := findDesign(wd)
		if err != nil {
			defaultErr = err
			return
		}
		data, err := os.ReadFile(path)
		if err != nil {
			defaultErr = err
			return
		}
		defaultDeclarations, defaultErr = loadDeclarations(data)
	})
	return defaultDeclarations, defaultErr
}
