// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package markdown

import (
	"encoding/json"
	"fmt"
	"strings"
)

// A marker is the one-line HTML comment that tags the visible block below it:
//
//	<!-- req: R1.1 | status: retired | replaced_by: ["R1.5"] -->
//
// The first segment names the tag and its value; later segments are
// attributes carrying machine-only data. Segments are separated by " | ".
// Attribute values are raw text when that is unambiguous and JSON otherwise,
// with "|" escaped so a value never splits a segment.
type Marker struct {
	Tag   string
	Value string
	Attrs []Attr
}

type Attr struct {
	Key   string
	Value string
}

// Reference is the JSON form of one provenance reference in a references
// attribute, such as derived_from.
type Reference struct {
	Repo    string `json:"repo"`
	Release string `json:"release,omitempty"`
	SRD     string `json:"srd"`
	Item    string `json:"item"`
	Reason  string `json:"reason,omitempty"`
}

const (
	markerOpen  = "<!--"
	markerClose = "-->"
	segmentSep  = " | "
)

// parseMarker reads one marker line. It reports false for a comment that is
// not a marker, such as the license header.
func parseMarker(line string) (Marker, bool, error) {
	text := strings.TrimSpace(line)
	if !strings.HasPrefix(text, markerOpen) || !strings.HasSuffix(text, markerClose) {
		return Marker{}, false, nil
	}
	body := strings.TrimSpace(text[len(markerOpen) : len(text)-len(markerClose)])
	if isLicenseComment(body) {
		return Marker{}, false, nil
	}
	segments := strings.Split(body, segmentSep)
	tag, value, ok := strings.Cut(segments[0], ":")
	if !ok {
		return Marker{}, true, fmt.Errorf("marker %q: expected \"<tag>: <value>\"", text)
	}
	m := Marker{Tag: strings.TrimSpace(tag), Value: strings.TrimSpace(value)}
	for _, segment := range segments[1:] {
		key, raw, ok := strings.Cut(segment, ":")
		if !ok {
			return Marker{}, true, fmt.Errorf("marker %q: attribute %q: expected \"<key>: <value>\"", text, segment)
		}
		m.Attrs = append(m.Attrs, Attr{Key: strings.TrimSpace(key), Value: strings.TrimSpace(raw)})
	}
	return m, true, nil
}

func isLicenseComment(body string) bool {
	return strings.HasPrefix(body, "Copyright ") || strings.HasPrefix(body, "SPDX-License-Identifier:")
}

// String prints the marker in canonical form.
func (m Marker) String() string {
	var b strings.Builder
	b.WriteString(markerOpen + " " + m.Tag + ": " + m.Value)
	for _, a := range m.Attrs {
		b.WriteString(segmentSep + a.Key + ": " + a.Value)
	}
	b.WriteString(" " + markerClose)
	return b.String()
}

// encodeAttr renders an attribute value: a plain string stays raw when
// nothing in it could be misread; everything else is JSON.
func encodeAttr(value any) (string, error) {
	if s, ok := value.(string); ok && rawAttrSafe(s) {
		return s, nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(data), "|", "\\u007c"), nil
}

func rawAttrSafe(s string) bool {
	if s == "" || strings.TrimSpace(s) != s || strings.ContainsAny(s, "|\n") {
		return false
	}
	return !strings.ContainsAny(s[:1], `[{"`) && !strings.Contains(s, markerClose)
}

// decodeAttr reads an attribute value into target. A string target takes
// the raw text unless the value is a JSON string; every other target is JSON.
func decodeAttr(raw string, target any) error {
	if s, ok := target.(*string); ok && !strings.HasPrefix(raw, `"`) {
		*s = raw
		return nil
	}
	if err := json.Unmarshal([]byte(raw), target); err != nil {
		return fmt.Errorf("attribute value %q: %w", raw, err)
	}
	return nil
}
