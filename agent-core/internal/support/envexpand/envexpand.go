// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

// Package envexpand expands environment references in mounted configuration
// before it is parsed. REST definitions and tool declarations both use it, so a
// single mounted profile parameterizes per pod by the same rules whichever form
// carries the address (srd013 R5.6, srd028 R2).
package envexpand

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// refPattern matches ${NAME} and ${NAME:-default} references. It is
// brace-delimited on purpose: the pervasive $.jsonpath and $from(label)
// selectors in a configuration file carry no brace and never match, so
// expansion leaves them byte-identical.
var refPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)(:-([^}]*))?\}`)

// Expand replaces ${VAR} and ${VAR:-default} references with environment
// values. A set variable wins; otherwise the default is used, or empty when the
// variable is unset and no default is given.
func Expand(data []byte) []byte {
	//nolint:forbidigo // srd013 R5.6/R5.7 requires ${VAR:-default} expansion at the declaration boundary.
	return ExpandWith(data, os.LookupEnv)
}

// ExpandWith is Expand with an explicit variable source, for a caller that
// expands on behalf of another process: an operator rendering a pod's profile
// with the environment that pod will see reads nothing from its own
// (srd009 R1.2).
func ExpandWith(data []byte, lookup func(string) (string, bool)) []byte {
	return refPattern.ReplaceAllFunc(data, func(match []byte) []byte {
		groups := refPattern.FindSubmatch(match)
		if v, ok := lookup(string(groups[1])); ok {
			return []byte(v)
		}
		if len(groups[2]) > 0 { // the ":-default" group is present
			return groups[3]
		}
		return nil
	})
}

// Reference is one ${NAME} or ${NAME:-default} occurrence in a declaration.
type Reference struct {
	Name       string
	Default    string
	HasDefault bool
}

// Mark replaces every reference in data with what mark returns for it, so an
// unexpanded declaration can be parsed as YAML: a ${NAME:-default} is a flow
// mapping to a YAML parser, a plain token in its place is a scalar.
func Mark(data []byte, mark func(Reference) string) []byte {
	return refPattern.ReplaceAllFunc(data, func(match []byte) []byte {
		groups := refPattern.FindSubmatch(match)
		return []byte(mark(Reference{Name: string(groups[1]), Default: string(groups[3]), HasDefault: len(groups[2]) > 0}))
	})
}

// References lists the references in s in order of appearance, so a reader
// of an unexpanded address learns which variables a deployment may set.
func References(s string) []Reference {
	var refs []Reference
	for _, match := range refPattern.FindAllStringSubmatch(s, -1) {
		refs = append(refs, Reference{Name: match[1], Default: match[3], HasDefault: match[2] != ""})
	}
	return refs
}

// Templated reports whether s carries a ${...} reference, well-formed or not.
func Templated(s string) bool {
	return strings.Contains(s, "${")
}

// Variant is a path that selects among files by environment: every reference
// is ${NAME:-default}, so the default names the file a deployment gets when
// NAME is unset, and Pattern is the path with each reference as a wildcard.
type Variant struct {
	Default string
	Pattern string
}

// SelectVariant reads a path whose references select a variant (srd052 R4.3).
// A reference without a default, or a ${ outside the reference grammar, is an
// error: the default is what names the variant every deployment can rely on.
func SelectVariant(path string) (Variant, error) {
	matches := refPattern.FindAllStringSubmatch(path, -1)
	for _, match := range matches {
		if match[2] == "" {
			return Variant{}, fmt.Errorf("reference %s needs the form ${NAME:-default}", match[0])
		}
	}
	if strings.Contains(refPattern.ReplaceAllString(path, ""), "${") {
		return Variant{}, fmt.Errorf("reference in %q is not ${NAME:-default}", path)
	}
	return Variant{
		Default: refPattern.ReplaceAllString(path, "${3}"),
		Pattern: refPattern.ReplaceAllString(path, "*"),
	}, nil
}

// ExpandString is Expand for one string.
func ExpandString(s string) string {
	return string(Expand([]byte(s)))
}
