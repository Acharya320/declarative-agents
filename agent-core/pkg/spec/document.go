// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec/markdown"
	"gopkg.in/yaml.v3"
)

// Corpus documents are YAML or Markdown during the migration window
// (docs/engineering/eng10-markdown-spec-format.md). A Markdown document
// reaches every parser as the YAML tree its declaration parses it into, so
// the typed models and every check stay format-blind.
const (
	yamlExt     = ".yaml"
	markdownExt = ".md"
)

// Audit checks the document format raises.
const (
	checkMarkdownFormat = "markdown-format"
	checkDualFormat     = "dual-format-document"
	checkNotCanonical   = "not-canonical"
)

// readDocument returns a corpus document as YAML: the file itself, or the
// YAML tree a Markdown document parses into.
func readDocument(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil || !strings.HasSuffix(path, markdownExt) {
		return data, err
	}
	return markdownAsYAML(data, path)
}

// markdownAsYAML parses a Markdown corpus document under its declaration and
// encodes the tree as YAML. A document with structural findings has no tree;
// the format scan reports those findings.
func markdownAsYAML(data []byte, source string) ([]byte, error) {
	d, err := declarationFor(source)
	if err != nil {
		return nil, err
	}
	doc, findings, err := d.Parse(data, source)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, fmt.Errorf("%s: %d Markdown format finding(s), first %s", source, len(findings), findings[0])
	}
	return yaml.Marshal(doc)
}

func declarationFor(source string) (*markdown.Declaration, error) {
	declarations, err := markdown.Default()
	if err != nil {
		return nil, fmt.Errorf("design constitution: %w", err)
	}
	d, ok := declarations.ForMarkdown(source)
	if !ok {
		return nil, fmt.Errorf("%s: no document type declares this Markdown location", source)
	}
	return d, nil
}

// formatScan picks each corpus document's file and records the findings its
// format raises: the same document in both formats, a Markdown document out
// of canonical form, and every finding the Markdown diagnostics report.
// Findings are collected, never fatal: a document whose structure cannot be
// read is left out and the rest of the corpus loads.
type formatScan struct {
	rootDir  string
	files    map[string]string
	findings []Finding
}

func newFormatScan(rootDir string) *formatScan {
	return &formatScan{rootDir: rootDir, files: map[string]string{}}
}

// record notes the file a document of kind and id was read from.
func (s *formatScan) record(kind, id, path string) {
	s.files[kind+"/"+id] = s.relative(path)
}

// documents lists the corpus documents in dir whose names start with prefix,
// one file per document: the YAML file when both formats exist (and a
// dual-format finding), and only Markdown documents that parse.
func (s *formatScan) documents(dir, prefix string) ([]string, error) {
	stems := map[string]map[string]bool{}
	for _, ext := range []string{yamlExt, markdownExt} {
		matches, err := filepath.Glob(filepath.Join(dir, prefix+"*"+ext))
		if err != nil {
			return nil, fmt.Errorf("glob %s: %w", dir, err)
		}
		for _, match := range matches {
			stem := strings.TrimSuffix(match, ext)
			if stems[stem] == nil {
				stems[stem] = map[string]bool{}
			}
			stems[stem][ext] = true
		}
	}
	var paths []string
	for _, stem := range sortedStems(stems) {
		if path, ok := s.choose(stem, stems[stem]); ok {
			paths = append(paths, path)
		}
	}
	return paths, nil
}

// single picks the file of a document with one fixed name, such as
// docs/road-map, or "" when neither format exists.
func (s *formatScan) single(stem string) string {
	present := map[string]bool{}
	for _, ext := range []string{yamlExt, markdownExt} {
		if _, err := os.Stat(stem + ext); err == nil {
			present[ext] = true
		}
	}
	if len(present) == 0 {
		return ""
	}
	path, ok := s.choose(stem, present)
	if !ok {
		return ""
	}
	return path
}

func (s *formatScan) choose(stem string, present map[string]bool) (string, bool) {
	if present[yamlExt] && present[markdownExt] {
		s.add(Finding{Check: checkDualFormat, Level: "error", File: s.relative(stem + markdownExt), Line: 1,
			Message: "document exists as both " + s.relative(stem+yamlExt) + " and " +
				s.relative(stem+markdownExt) + "; keep one (mage convert or mage convertBack)"})
		return stem + yamlExt, true
	}
	if present[yamlExt] {
		return stem + yamlExt, true
	}
	return stem + markdownExt, s.checkMarkdown(stem + markdownExt)
}

// checkMarkdown diagnoses one Markdown document and reports whether its
// structure lets it load.
func (s *formatScan) checkMarkdown(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		s.add(Finding{Check: checkMarkdownFormat, Level: "error", File: s.relative(path), Message: err.Error()})
		return false
	}
	d, err := declarationFor(path)
	if err != nil {
		s.add(Finding{Check: checkMarkdownFormat, Level: "error", File: s.relative(path), Message: err.Error()})
		return false
	}
	formatted, findings, err := d.Format(data, path)
	loads := err == nil && formatted != nil
	for _, f := range findings {
		s.add(Finding{Check: checkMarkdownFormat, CheckID: f.Name, Level: "error", File: s.relative(path), Line: f.Line,
			Message: fmt.Sprintf("%s: expected %s, found %s [design.yaml %s]; %s", f.Name, f.Expected, f.Found, f.Rule, f.Hint)})
	}
	if loads && string(formatted) != string(data) {
		s.add(Finding{Check: checkNotCanonical, Level: "error", File: s.relative(path), Line: 1,
			Message: "document is not in canonical form; run mage fmt"})
	}
	return loads
}

func (s *formatScan) add(f Finding) { s.findings = append(s.findings, f) }

func (s *formatScan) relative(path string) string {
	if rel, err := filepath.Rel(s.rootDir, path); err == nil {
		return filepath.ToSlash(rel)
	}
	return filepath.ToSlash(path)
}

func sortedStems(stems map[string]map[string]bool) []string {
	keys := make([]string, 0, len(stems))
	for stem := range stems {
		keys = append(keys, stem)
	}
	sort.Strings(keys)
	return keys
}

// checkDocumentFormat returns the findings the format scan recorded while
// the corpus loaded.
func checkDocumentFormat(corpus *Corpus) []Finding {
	return corpus.FormatFindings
}

// documentFile names the file a document was read from, or the YAML path its
// id implies when the corpus did not record one.
func (c *Corpus) documentFile(kind, subdir, id string) string {
	if file, ok := c.Files[kind+"/"+id]; ok {
		return file
	}
	return filepath.ToSlash(filepath.Join(subdir, id+yamlExt))
}

// isCorpusDocument reports a file name in either corpus format.
func isCorpusDocument(file string) bool {
	return strings.HasSuffix(file, yamlExt) || strings.HasSuffix(file, markdownExt)
}

// documentYAML returns a document's content as YAML, converting a Markdown
// document read from a git ref.
func documentYAML(content, file string) ([]byte, error) {
	if strings.HasSuffix(file, markdownExt) {
		return markdownAsYAML([]byte(content), file)
	}
	return []byte(content), nil
}

// documentExists reports whether a path a document names exists. During the
// migration window a corpus document named by its .yaml path may already be
// Markdown, so a missing YAML document resolves to its Markdown form, which
// an audit register names with its -register suffix.
func documentExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	if !strings.HasSuffix(path, yamlExt) {
		return false
	}
	declarations, err := markdown.Default()
	if err != nil {
		return false
	}
	if d, ok := declarations.ForSource(filepath.ToSlash(path)); ok {
		_, err := os.Stat(d.MarkdownPath(path))
		return err == nil
	}
	return false
}

// explain adds to a corpus validation error the format findings of the
// documents the scan left out, which are often why a reference dangles.
func (s *formatScan) explain(err error) error {
	var lines []string
	for _, f := range s.findings {
		if f.Check == checkMarkdownFormat {
			lines = append(lines, fmt.Sprintf("  %s:%d: %s", f.File, f.Line, f.Message))
		}
	}
	if len(lines) == 0 {
		return err
	}
	return fmt.Errorf("%w\ndocuments left out for Markdown format findings:\n%s", err, strings.Join(lines, "\n"))
}
