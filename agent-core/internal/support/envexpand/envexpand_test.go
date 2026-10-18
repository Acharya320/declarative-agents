// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package envexpand

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpandReplacesEnvironmentReference(t *testing.T) {
	t.Setenv("ENVEXPAND_VALUE", "configured")

	require.Equal(t, []byte("value=configured"), Expand([]byte("value=${ENVEXPAND_VALUE}")))
}

func TestExpandUsesDefaultForUndefinedVariable(t *testing.T) {
	unsetForTest(t, "ENVEXPAND_DEFAULT")

	require.Equal(t, []byte("value=fallback"), Expand([]byte("value=${ENVEXPAND_DEFAULT:-fallback}")))
}

func TestExpandRemovesUndefinedVariableWithoutDefault(t *testing.T) {
	unsetForTest(t, "ENVEXPAND_UNDEFINED")

	require.Equal(t, []byte("value="), Expand([]byte("value=${ENVEXPAND_UNDEFINED}")))
}

func TestExpandPreservesSelectorsExactly(t *testing.T) {
	input := []byte("$.path\n$from(label).field")

	require.Equal(t, input, Expand(input))
}

func unsetForTest(t *testing.T, name string) {
	t.Helper()
	t.Setenv(name, "temporary")
	require.NoError(t, os.Unsetenv(name))
}

func TestExpandWithReadsOnlyTheGivenLookup(t *testing.T) {
	t.Setenv("ONLY_IN_PROCESS", "leaked")
	lookup := func(name string) (string, bool) {
		if name == "GIVEN" {
			return "value", true
		}
		return "", false
	}
	got := string(ExpandWith([]byte("${GIVEN}/${ONLY_IN_PROCESS:-default}/${UNSET}"), lookup))
	if got != "value/default/" {
		t.Fatalf("ExpandWith = %q", got)
	}
}

func TestReferencesListsNamesAndDefaults(t *testing.T) {
	refs := References("${BIND:-127.0.0.1}:${PORT} $.jsonpath")
	if len(refs) != 2 {
		t.Fatalf("refs = %+v", refs)
	}
	if refs[0] != (Reference{Name: "BIND", Default: "127.0.0.1", HasDefault: true}) || refs[1] != (Reference{Name: "PORT"}) {
		t.Errorf("refs = %+v", refs)
	}
}

func TestMarkReplacesReferencesWithScalars(t *testing.T) {
	got := string(Mark([]byte("ports: [${PORT:-18191}] ${HOST}"), func(ref Reference) string {
		if ref.HasDefault {
			return "D_" + ref.Name
		}
		return "N_" + ref.Name
	}))
	if got != "ports: [D_PORT] N_HOST" {
		t.Fatalf("Mark = %q", got)
	}
}
