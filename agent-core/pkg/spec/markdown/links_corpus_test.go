// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// brokenLinks returns every reference link in a document whose file does not
// exist under root or whose anchor no heading of that file takes. Only links
// around a reference (a document id or an item) are checked: those are the
// links the converter writes (GH-2671).
func brokenLinks(root, docPath string, data []byte, anchorsOf func(string) map[string]bool) []string {
	var broken []string
	for _, m := range markdownLink.FindAllStringSubmatch(withoutCode(string(data)), -1) {
		label, target := m[1], m[2]
		file, anchor, _ := strings.Cut(target, "#")
		reference := idToken.MatchString(label) && anchor == "" && strings.HasSuffix(file, ".md") ||
			itemToken.MatchString(label) && anchor != ""
		if !reference || strings.Contains(target, "://") {
			continue
		}
		resolved := docPath
		if file != "" {
			resolved = path.Join(path.Dir(docPath), file)
		}
		if _, err := os.Stat(filepath.Join(root, resolved)); err != nil {
			broken = append(broken, fmt.Sprintf("%s: [%s](%s): no file %s", docPath, label, target, resolved))
			continue
		}
		if anchor != "" && !anchorsOf(resolved)[anchor] {
			broken = append(broken, fmt.Sprintf("%s: [%s](%s): no heading anchor #%s in %s", docPath, label, target, anchor, resolved))
		}
	}
	return broken
}

// withoutCode drops fenced blocks and code spans, which GitHub renders as
// text rather than links.
func withoutCode(text string) string {
	var kept []string
	fenced := false
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "```") {
			fenced = !fenced
			continue
		}
		if !fenced {
			kept = append(kept, codeSpan.ReplaceAllString(line, ""))
		}
	}
	return strings.Join(kept, "\n")
}

var codeSpan = regexp.MustCompile("`[^`]*`")

// headingAnchors is every anchor GitHub gives a page's headings, repeats
// included.
func headingAnchors(data []byte) map[string]bool {
	anchors, seen, fenced := map[string]bool{}, map[string]int{}, false
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "```") {
			fenced = !fenced
			continue
		}
		level := len(line) - len(strings.TrimLeft(line, "#"))
		if fenced || level < 1 || level > maxHeadingLevel || !strings.HasPrefix(line[level:], " ") {
			continue
		}
		base := slug(strings.TrimSpace(line[level:]))
		anchor := base
		if n := seen[base]; n > 0 {
			anchor = base + "-" + strconv.Itoa(n)
		}
		seen[base]++
		anchors[anchor] = true
	}
	return anchors
}

func anchorReader(root string) func(string) map[string]bool {
	cache := map[string]map[string]bool{}
	return func(p string) map[string]bool {
		if anchors, ok := cache[p]; ok {
			return anchors
		}
		data, _ := os.ReadFile(filepath.Join(root, p))
		cache[p] = headingAnchors(data)
		return cache[p]
	}
}

// TestEveryCorpusLinkResolves follows every reference link in the tracked
// corpus to an existing file and heading anchor.
func TestEveryCorpusLinkResolves(t *testing.T) {
	root := "../../../.."
	anchors := anchorReader(root)
	checked := 0
	for _, docPath := range trackedCorpusPaths(t, root) {
		data, err := os.ReadFile(filepath.Join(root, docPath))
		if err != nil {
			t.Fatal(err)
		}
		checked += strings.Count(string(data), "](")
		for _, problem := range brokenLinks(root, docPath, data, anchors) {
			t.Error(problem)
		}
	}
	if checked == 0 {
		t.Fatal("the corpus holds no links; the check checked nothing")
	}
}

// TestBrokenLinksAreFound plants a link to a missing anchor and one to a
// missing file, and requires both reported.
func TestBrokenLinksAreFound(t *testing.T) {
	root := t.TempDir()
	write := func(p, body string) {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(root, p)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, p), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("docs/specs/srd/srd001-a.md", "## R1 - Budget\n\n### R1.1\n")
	doc := "- [srd001-a](../srd/srd001-a.md) [R1.1](../srd/srd001-a.md#r11), [R1.2](../srd/srd001-a.md#r12)\n" +
		"- [srd009-gone](../srd/srd009-gone.md)\n- [R1](#r1---budget)\n- [see](#nowhere)\n"
	write("docs/specs/tests/t.md", doc+"\n## R1 - Budget\n")
	broken := brokenLinks(root, "docs/specs/tests/t.md", []byte(doc), anchorReader(root))
	if len(broken) != 2 || !strings.Contains(broken[0], "#r12") || !strings.Contains(broken[1], "srd009-gone") {
		t.Fatalf("broken = %v", broken)
	}
}

// trackedCorpusPaths lists the tracked Markdown corpus documents outside
// testdata.
func trackedCorpusPaths(t *testing.T, root string) []string {
	t.Helper()
	out, err := exec.Command("git", "-C", root, "ls-files", "*.md").Output()
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	for _, p := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if _, ok := repositoryDeclarations(t).ForMarkdown(p); ok && !strings.Contains(p, "/testdata/") {
			paths = append(paths, p)
		}
	}
	return paths
}
