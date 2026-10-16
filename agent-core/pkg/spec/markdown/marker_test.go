// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"reflect"
	"testing"
)

func TestMarkerAttributesRoundTrip(t *testing.T) {
	replaced := []string{"R1.5", "R1.6"}
	note := "progressing | complete"
	m := Marker{Tag: "req", Value: "R1.2"}
	for _, v := range []struct {
		key   string
		value any
	}{{"status", "retired"}, {"replaced_by", replaced}, {"note", note}, {"weight", 3}} {
		encoded, err := EncodeAttr(v.value)
		if err != nil {
			t.Fatal(err)
		}
		m.Attrs = append(m.Attrs, Attr{Key: v.key, Value: encoded})
	}
	parsed, ok, err := ParseMarker(m.String())
	if err != nil || !ok {
		t.Fatalf("parse %q: ok=%v err=%v", m.String(), ok, err)
	}
	var status, gotNote string
	var gotReplaced []string
	var weight int
	for _, a := range parsed.Attrs {
		var err error
		switch a.Key {
		case "status":
			err = DecodeAttr(a.Value, &status)
		case "replaced_by":
			err = DecodeAttr(a.Value, &gotReplaced)
		case "note":
			err = DecodeAttr(a.Value, &gotNote)
		case "weight":
			err = DecodeAttr(a.Value, &weight)
		}
		if err != nil {
			t.Fatalf("decode %s: %v", a.Key, err)
		}
	}
	if status != "retired" || !reflect.DeepEqual(gotReplaced, replaced) || gotNote != note || weight != 3 {
		t.Fatalf("round trip lost a value: %q %v %q %d", status, gotReplaced, gotNote, weight)
	}
}

func TestMarkerRejectsMalformedAndSkipsLicense(t *testing.T) {
	if _, ok, err := ParseMarker("<!-- req R1.1 -->"); !ok || err == nil {
		t.Fatalf("missing colon: ok=%v err=%v, want a marker error", ok, err)
	}
	if _, ok, err := ParseMarker("<!-- Copyright (c) 2026 Nokia -->"); ok || err != nil {
		t.Fatalf("license comment: ok=%v err=%v, want not a marker", ok, err)
	}
	if _, ok, _ := ParseMarker("plain text"); ok {
		t.Fatal("plain text read as a marker")
	}
}
