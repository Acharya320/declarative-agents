// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"path/filepath"
	"testing"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec"
)

// TestApplicationsCorpusPassesTheTypedChecks loads the applications corpus
// (applications/docs) through the typed model, as module audits load their
// own corpora, and fails on any error-level finding. Until GH-2611 no gate
// loaded it this way, so a use case the parser rejected escaped every typed
// check without failing anything.
func TestApplicationsCorpusPassesTheTypedChecks(t *testing.T) {
	root := filepath.Join(repositoryRoot(t), "applications")
	corpus, err := spec.LoadCorpus(root)
	if err != nil {
		t.Fatalf("load applications corpus: %v", err)
	}
	documents, err := filepath.Glob(filepath.Join(root, spec.UCSubdir, spec.UCGlob))
	if err != nil {
		t.Fatal(err)
	}
	if len(corpus.UseCases) != len(documents) {
		t.Fatalf("loaded %d of %d use cases", len(corpus.UseCases), len(documents))
	}
	graph, err := spec.BuildGraph(corpus)
	if err != nil {
		t.Fatalf("build applications graph: %v", err)
	}
	for _, finding := range spec.Validate(graph, corpus) {
		if finding.Level == "error" {
			t.Errorf("%s %s: %s", finding.Check, finding.File, finding.Message)
		}
	}
}
