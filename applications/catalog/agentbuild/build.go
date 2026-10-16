// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

// Package agentbuild holds the one agent-CLI build recipe shared by the
// catalog magefile and the conformance harness. Both previously duplicated the
// production build tag, the ./cmd/agent target, and the go build invocation;
// magefiles is package main, so the recipe lives here as an importable
// non-main package alongside catalogroot (GH-1390).
package agentbuild

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

const (
	// BuildTag selects the production tool catalog for the agent CLI.
	BuildTag = "production"
	// target is the agent CLI main package relative to the agent-core root.
	target = "./cmd/agent"
)

// Build compiles the agent CLI from coreRoot to outputPath with the shared
// production recipe. Combined build output is written to logOut when it is
// non-nil (callers stream it live or buffer it to report only on failure).
// outputPath is caller-chosen so the magefile and the harness can keep
// distinct binaries. On success it returns outputPath.
//
// The binary is compiled under a name unique to this process and renamed into
// place, so a concurrent run of the same binary never executes a partly
// written file.
func Build(coreRoot, outputPath string, logOut io.Writer) (string, error) {
	staging := outputPath + ".build-" + strconv.Itoa(os.Getpid())
	cmd := exec.Command("go", "build", "-tags", BuildTag, "-o", staging, target)
	cmd.Dir = coreRoot
	if logOut != nil {
		cmd.Stdout = logOut
		cmd.Stderr = logOut
	}
	if err := cmd.Run(); err != nil {
		_ = os.Remove(staging)
		return "", fmt.Errorf("build agent from %s: %w", coreRoot, err)
	}
	if err := os.Rename(staging, outputPath); err != nil {
		_ = os.Remove(staging)
		return "", fmt.Errorf("install agent binary %s: %w", outputPath, err)
	}
	return outputPath, nil
}

// OutputPath names the temporary agent binary for name built from coreRoot:
// os.TempDir() is shared by every checkout on the machine, so the name carries
// a digest of the absolute agent-core root. Without it, two checkouts audited
// at the same time overwrite one binary and one runs the other's agent
// (GH-2593).
func OutputPath(name, coreRoot string) string {
	root, err := filepath.Abs(coreRoot)
	if err != nil {
		root = coreRoot
	}
	digest := sha256.Sum256([]byte(filepath.Clean(root)))
	return filepath.Join(os.TempDir(), name+"-"+hex.EncodeToString(digest[:])[:12])
}
