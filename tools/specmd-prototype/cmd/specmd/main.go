// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

// Command specmd converts an SRD between YAML and Markdown and formats
// Markdown SRDs in place:
//
//	specmd md <srd.yaml>     print the Markdown form
//	specmd yaml <srd.md>     print the YAML form
//	specmd fmt <srd.md>...   rewrite each file in canonical Markdown
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Nokia-Bell-Labs/declarative-agents/agent-core/pkg/spec"
	specmd "github.com/Nokia-Bell-Labs/declarative-agents/tools/specmd-prototype"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "specmd:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) < 2 {
		return errors.New("usage: specmd md <srd.yaml> | yaml <srd.md> | fmt <srd.md> [<srd.md> ...]")
	}
	switch args[0] {
	case "md":
		srd, err := spec.ParseSRD(args[1])
		if err != nil {
			return err
		}
		return write(specmd.PrintMarkdown(&srd))
	case "yaml":
		srd, err := readMarkdown(args[1])
		if err != nil {
			return err
		}
		return write(specmd.PrintYAML(&srd))
	case "fmt":
		for _, path := range args[1:] {
			if err := formatInPlace(path); err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("unknown command %q", args[0])
}

func readMarkdown(path string) (spec.SRD, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return spec.SRD{}, err
	}
	return specmd.ParseMarkdown(data, path)
}

func formatInPlace(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	formatted, err := specmd.FormatMarkdown(data, path)
	if err != nil {
		return err
	}
	return os.WriteFile(path, formatted, 0o644)
}

func write(data []byte, err error) error {
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(data)
	return err
}
