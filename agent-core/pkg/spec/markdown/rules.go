// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Declaration is one document type's Markdown declaration, read from the markdown
// block of its entry in design.yaml's document_types: heading grammar, typed
// fields and their shapes, attribute shapes, and prose containment. Its
// location is the type's location, and the type's converts_from names the
// former YAML location mage convert reads.
// Diagnostics read every constraint from here, so a new document type is a
// new declaration, not new code.
type Declaration struct {
	Location    string   `yaml:"-"`       // the type's location
	Exclude     []string `yaml:"exclude"` // source documents that stay in their format
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

	documentType   string
	pathPattern    *regexp.Regexp
	sourcePattern  *regexp.Regexp
	sourceLocation string
	idPatterns     map[string]*regexp.Regexp
	itemPatterns   map[string]*regexp.Regexp
	keywordRegexp  *regexp.Regexp
}

// Node heading forms and content kinds a declaration may name.
const (
	headingID      = "id"       // the heading is exactly the marker's id
	headingIDTitle = "id_title" // "<id>" or "<id> - <title>"
	contentNone    = "none"     // the node holds only tagged children
	contentProse   = "prose"    // required text, no untagged table or list
)

// Declarations holds every declared document type's Markdown declaration,
// keyed by document type name.
type Declarations map[string]*Declaration

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
	YAML            nodeYAML `yaml:"yaml"`
	// Fields lists, in declared order, every field a field slot holds; a slot
	// without the list admits any field in any order (GH-2654).
	Fields []string `yaml:"fields"`
}

type fieldRule struct {
	In          []string `yaml:"in"`
	Shape       string   `yaml:"shape"`
	ItemPattern string   `yaml:"item_pattern"`
	Columns     []string `yaml:"columns"`
	Heading     string   `yaml:"heading"`      // printed heading; default from the name
	Value       string   `yaml:"value"`        // keyed or outline; default by value
	ID          string   `yaml:"id"`           // outline: the record key holding the id
	ContentKey  string   `yaml:"content"`      // outline: the record key holding the text
	Split       string   `yaml:"split"`        // table: separator of a keyed "<id>: <a><split><b>" row
	Code        string   `yaml:"code"`         // a string printed in one fence of this language
	Links       bool     `yaml:"links"`        // its values are references, rendered as links
	LinkColumns []string `yaml:"link_columns"` // the table columns whose cells are references
}

// LoadDeclarations reads every document type in design.yaml whose markdown
// entry is a declaration (a mapping); other markdown values, such as a note
// that a type stays YAML, declare no Markdown rules.
func LoadDeclarations(data []byte) (Declarations, error) {
	var design struct {
		DocumentTypes map[string]struct {
			Location     string    `yaml:"location"`
			ConvertsFrom string    `yaml:"converts_from"`
			Markdown     yaml.Node `yaml:"markdown"`
		} `yaml:"document_types"`
	}
	if err := yaml.Unmarshal(data, &design); err != nil {
		return nil, fmt.Errorf("design.yaml: %w", err)
	}
	declarations := Declarations{}
	for name, entry := range design.DocumentTypes {
		if entry.Markdown.Kind != yaml.MappingNode {
			continue
		}
		var r Declaration
		if err := entry.Markdown.Decode(&r); err != nil {
			return nil, fmt.Errorf("document_types.%s.markdown: %w", name, err)
		}
		r.documentType = name
		r.Location = entry.Location
		if err := r.compile(); err != nil {
			return nil, fmt.Errorf("document_types.%s.markdown.%w", name, err)
		}
		if entry.ConvertsFrom != "" {
			source, err := locationPattern(entry.ConvertsFrom)
			if err != nil {
				return nil, fmt.Errorf("document_types.%s.converts_from: %w", name, err)
			}
			r.sourcePattern = source
			r.sourceLocation = entry.ConvertsFrom
		}
		declarations[name] = &r
	}
	return declarations, nil
}

// compile prepares a declaration's patterns once.
func (r *Declaration) compile() error {
	pattern, err := locationPattern(r.Location)
	if err != nil {
		return fmt.Errorf("location: %w", err)
	}
	r.pathPattern = pattern
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
	return r.checkMapping()
}

// locationPattern turns a declared location such as
// "docs/specs/software-requirements/srdNNN-short-name.md" into a pattern over
// real paths, matched at a path boundary so any module's docs directory
// qualifies: a run of N after a lowercase letter or a dot is that many
// digits, short-name is kebab-case, and name is any lowercase name. The
// directory matters because several types share the name placeholder.
func locationPattern(location string) (*regexp.Regexp, error) {
	if !strings.HasPrefix(location, "docs/") || strings.HasSuffix(location, "/") {
		return nil, errors.New("a location is a file path under docs/")
	}
	var pattern strings.Builder
	pattern.WriteString(`(?:^|/)`)
	for i := 0; i < len(location); {
		switch {
		case strings.HasPrefix(location[i:], "-short-name"):
			pattern.WriteString(`(?:-[a-z0-9]+(?:-[a-z0-9]+)*)?`)
			i += len("-short-name")
		case strings.HasPrefix(location[i:], "short-name"):
			pattern.WriteString(`[a-z0-9]+(?:-[a-z0-9]+)*`)
			i += len("short-name")
		case strings.HasPrefix(location[i:], "name"):
			pattern.WriteString(`[a-z0-9][a-z0-9._-]*`)
			i += len("name")
		case location[i] == 'N' && isNumericPlaceholder(location, i):
			end := i
			for end < len(location) && location[end] == 'N' {
				end++
			}
			fmt.Fprintf(&pattern, `[0-9]{%d}`, end-i)
			i = end
		default:
			pattern.WriteString(regexp.QuoteMeta(location[i : i+1]))
			i++
		}
	}
	pattern.WriteString("$")
	return regexp.Compile(pattern.String())
}

// isNumericPlaceholder reports whether the N at index starts a digit
// placeholder: at the start, or after a dot or a lowercase letter, so
// "srdNNN" and "relNN.N" are placeholders and "VISION" is a word. The
// placement gate in magefiles reads locations by the same rule.
func isNumericPlaceholder(location string, index int) bool {
	if index == 0 {
		return true
	}
	previous := location[index-1]
	return previous == '.' || (previous >= 'a' && previous <= 'z')
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
func (r *Declaration) nodeFor(tag string, level int, parent string) (string, bool) {
	for name, node := range r.Nodes {
		if node.Tag == tag && node.Level == level && node.Parent == parent {
			return name, true
		}
	}
	return "", false
}

// allowedHere lists the "<tag> at level <n>" headings legal under parent,
// for messages.
func (r *Declaration) allowedHere(parent string) []string {
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
func (r *Declaration) tags() []string {
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

// declarationFor picks the declaration whose location names this file: by
// the path as given, else by its absolute path, so a relative name such as
// "eng10-x.md" run inside docs/engineering still resolves.
func declarationFor(declarations Declarations, source string) (*Declaration, bool) {
	if r, ok := declarationForPath(declarations, filepath.ToSlash(source)); ok {
		return r, true
	}
	abs, err := filepath.Abs(source)
	if err != nil {
		return nil, false
	}
	return declarationForPath(declarations, filepath.ToSlash(abs))
}

func declarationForPath(declarations Declarations, path string) (*Declaration, bool) {
	names := make([]string, 0, len(declarations))
	for name := range declarations {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if declarations[name].pathPattern.MatchString(path) {
			return declarations[name], true
		}
	}
	return nil, false
}

// Type names the document type the declaration belongs to.
func (r *Declaration) Type() string { return r.documentType }

// ForSource finds the declaration of the type whose converts_from location
// names the path, the YAML document a conversion starts from.
func (declarations Declarations) ForSource(path string) (*Declaration, bool) {
	slashed := filepath.ToSlash(path)
	for _, name := range sortedKeys(declarations) {
		d := declarations[name]
		if d.sourcePattern != nil && d.sourcePattern.MatchString(slashed) && !d.excluded(slashed) {
			return d, true
		}
	}
	return nil, false
}

// excluded reports a source document the declaration leaves in its format,
// named by its path from a docs directory up (docs/constitutions/design.yaml).
func (r *Declaration) excluded(path string) bool {
	for _, exclude := range r.Exclude {
		if path == exclude || strings.HasSuffix(path, "/"+exclude) {
			return true
		}
	}
	return false
}

// ForMarkdown finds the declaration whose Markdown location names the path.
func (declarations Declarations) ForMarkdown(path string) (*Declaration, bool) {
	return declarationFor(declarations, path)
}
