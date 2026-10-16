// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package specmd

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
)

// parentAt names the context a heading at level opens under.
func (d *diagnosis) parentAt(level int) string {
	switch level {
	case 2:
		return "document"
	case 3:
		switch {
		case d.l2 == nil:
			return "none"
		case d.l2.kind == "requirement_group":
			return "requirement_group"
		case d.l2.kind == "document_field" && d.l2.id == keyAcceptanceCriteria:
			return "acceptance_criteria"
		}
		return d.l2.kind
	case 4:
		if d.l3 == nil {
			return "none"
		}
		return d.l3.kind
	}
	return fmt.Sprintf("level %d", level)
}

func (d *diagnosis) heading(b block, level int, heading string) {
	d.closeOpen()
	line := d.fileLine(b)
	if d.pending == nil {
		d.add("heading-untagged", 2, line, "a marker on the line above the heading", "\""+heading+"\" with no marker",
			"nodes", "add the marker that says what the heading opens, e.g. <!-- field: <key> -->")
		d.open = nil
		return
	}
	m, markerLine := *d.pending, d.pendLine
	d.pending = nil
	parent := d.parentAt(level)
	kind, ok := d.rules.nodeFor(m.Tag, level, parent)
	if !ok {
		d.add("heading-misplaced", 2, line, strings.Join(d.rules.allowedHere(parent), " or ")+" under "+parent,
			fmt.Sprintf("%s at level %d", m.Tag, level), "nodes",
			"change the heading depth or move the section under the node that may hold it")
		d.open = nil
		return
	}
	node := &dnode{kind: kind, id: m.Value, line: line, markerLine: markerLine, fields: map[string]bool{}}
	d.checkNode(node, m, heading)
	d.place(node, level)
}

// place makes node the open node at its depth.
func (d *diagnosis) place(node *dnode, level int) {
	switch level {
	case 2:
		d.l2, d.l3 = node, nil
	case 3:
		d.l3 = node
	}
	d.open = node
}

func (d *diagnosis) checkNode(node *dnode, m marker, heading string) {
	rule := d.rules.Nodes[node.kind]
	if pattern, ok := d.rules.idPatterns[node.kind]; ok && !pattern.MatchString(m.Value) {
		d.add("id-malformed", 2, node.line, "an id matching "+rule.IDPattern, m.Value,
			"nodes."+node.kind+".id_pattern", "renumber to the declared form; ids are never reused (design.yaml D6)")
	}
	if rule.IDExtendsParent && d.l2 != nil && !strings.HasPrefix(m.Value, d.l2.id+".") {
		d.add("item-wrong-group", 2, node.line, "an id under group "+d.l2.id, m.Value,
			"nodes."+node.kind+".id_extends_parent", "move the item under its own group, or renumber it as a retirement")
	}
	if rule.HeadingIsID && heading != m.Value {
		d.add("heading-id-mismatch", 2, node.line, "heading \""+m.Value+"\"", "\""+heading+"\"",
			"nodes."+node.kind+".heading_is_id", "make the heading say exactly the id in its marker")
	}
	if node.kind == "requirement_group" {
		d.groups++
		if heading != m.Value && !strings.HasPrefix(heading, m.Value+" - ") {
			d.add("heading-id-mismatch", 2, node.line, "heading \""+m.Value+" - <title>\"", "\""+heading+"\"",
				"nodes.requirement_group", "start the group heading with its id and \" - \"")
		}
	}
	if m.Tag == "req" || m.Tag == "ac" {
		if first, seen := d.ids[m.Value]; seen {
			d.add("duplicate-id", 2, node.line, "an id used once", fmt.Sprintf("%s, first used on line %d", m.Value, first),
				"nodes", "give the second node a new id; an id is never reused (design.yaml D6)")
		} else {
			d.ids[m.Value] = node.line
		}
	}
	if m.Tag == "field" {
		d.checkField(node)
	}
	d.checkAttributes(node, m)
}

// checkField enforces where a typed field may appear and that a node holds
// each field once.
func (d *diagnosis) checkField(node *dnode) {
	parentNode := d.l2
	if node.kind == "item_field" || node.kind == "criterion_field" {
		parentNode = d.l3
	}
	if node.kind == "document_field" {
		if d.docKeys[node.id] {
			d.add("field-duplicate", 2, node.line, "field "+node.id+" once per node", "a second "+node.id,
				"fields", "merge the two sections")
		}
		d.docKeys[node.id] = true
	} else if parentNode != nil {
		if parentNode.fields[node.id] {
			d.add("field-duplicate", 2, node.line, "field "+node.id+" once per node", "a second "+node.id,
				"fields", "merge the two sections")
		}
		parentNode.fields[node.id] = true
	}
	if field, typed := d.rules.Fields[node.id]; typed && !contains(field.In, node.kind) {
		d.add("field-disallowed", 2, node.line, node.id+" only as "+strings.Join(field.In, " or "), node.kind,
			"fields."+node.id+".in", "move the section to where the field belongs")
	}
}

func (d *diagnosis) checkAttributes(node *dnode, m marker) {
	allowed := d.rules.Nodes[node.kind].Attributes
	for _, a := range m.Attrs {
		if !contains(allowed, a.Key) {
			d.add("marker-unknown-attribute", 1, node.markerLine, "an attribute "+node.kind+" declares", a.Key,
				"nodes."+node.kind+".attributes", suggest(a.Key, allowed))
			continue
		}
		if !attributeShapeOK(d.rules.Attributes[a.Key], a.Value) {
			d.add("attribute-shape", 3, node.markerLine, a.Key+" as "+d.rules.Attributes[a.Key], a.Value,
				"attributes."+a.Key, "write the value in its declared shape (JSON for lists and references)")
		}
	}
	if node.kind == "acceptance_criterion" {
		d.checkTraces(node, m)
	}
}

func attributeShapeOK(shape, raw string) bool {
	switch shape {
	case "integer":
		var n int
		return decodeAttr(raw, &n) == nil
	case "list":
		var values []string
		return decodeAttr(raw, &values) == nil
	case "references":
		var refs []refJSON
		return decodeAttr(raw, &refs) == nil
	}
	var s string
	return decodeAttr(raw, &s) == nil
}

// checkTraces enforces the acceptance comment's arity and kind: a criterion
// traces at least one requirement, by id.
func (d *diagnosis) checkTraces(node *dnode, m marker) {
	raw, _ := markerAttr(m, "traces")
	var traces []string
	if raw != "" && decodeAttr(raw, &traces) != nil {
		return // attribute-shape already reports it
	}
	if len(traces) < d.rules.CriterionTraces.Min {
		d.add("criterion-traces-missing", 3, node.markerLine, fmt.Sprintf("at least %d trace", d.rules.CriterionTraces.Min),
			fmt.Sprintf("%d", len(traces)), "criterion_traces.min",
			"add | traces: [\"R1.1\"] to the marker naming the requirements this criterion verifies")
	}
	for _, trace := range traces {
		if !d.rules.tracePattern.MatchString(trace) {
			d.add("criterion-trace-invalid", 3, node.markerLine, "a requirement id matching "+d.rules.CriterionTraces.IDPattern,
				trace, "criterion_traces.id_pattern", "trace a requirement id such as R1 or R1.1")
		}
	}
}

func markerAttr(m marker, key string) (string, bool) {
	for _, a := range m.Attrs {
		if a.Key == key {
			return a.Value, true
		}
	}
	return "", false
}

// closeOpen checks the content the open node collected.
func (d *diagnosis) closeOpen() {
	node := d.open
	d.open = nil
	if node == nil {
		return
	}
	rule := d.rules.Nodes[node.kind]
	switch {
	case node.kind == "requirement_group" && len(node.content) > 0:
		d.add("content-untagged", 2, d.fileLine(node.content[0]), "content under a tagged field or item",
			"untagged content under group "+node.id, "nodes.requirement_group",
			"tag it with <!-- field: <key> --> and a level-3 heading, or move it into an item")
	case rule.RequiresContent && len(node.content) == 0:
		d.add("missing-content", 2, node.line, node.kind+" text", "an empty "+node.id,
			"nodes."+node.kind+".requires_content", "write the "+node.kind+" text under its heading")
	case rule.RequiresContent:
		if b, kind, found := containsStructure(node.content); found {
			d.add("untagged-structure", 2, d.fileLine(b), "prose only in "+node.kind+" text", "an untagged "+kind,
				"nodes."+node.kind, "move the "+kind+" under its own <!-- field: <key> --> with a level-4 heading")
		}
	}
	if strings.HasSuffix(node.kind, "field") {
		d.checkShape(node)
	}
	if contains(d.rules.ProseContainment.Outside, node.kind) {
		d.proseContainment(node)
	}
}

// checkShape enforces a typed field's content shape, list prefixes, and
// table columns.
func (d *diagnosis) checkShape(node *dnode) {
	field, typed := d.rules.Fields[node.id]
	if !typed {
		return
	}
	content := d.contentText(node)
	shape := shapeEmpty
	if content != "" {
		shape = classifyContent(content)
	}
	line := node.line
	if len(node.content) > 0 {
		line = d.fileLine(node.content[0])
	}
	if shape != field.Shape && (field.Shape != shapeProse || shape != shapeYAML) {
		d.add("field-shape-mismatch", 3, line, node.id+" as "+field.Shape, shape, "fields."+node.id+".shape",
			"write "+node.id+" as "+field.Shape)
		return
	}
	switch field.Shape {
	case shapeList:
		d.checkListPrefix(node, content, line)
	case shapeTable:
		if header := splitRow(strings.SplitN(content, "\n", 2)[0]); !sameStrings(header, field.Columns) {
			d.add("table-columns-mismatch", 3, line, "columns "+strings.Join(field.Columns, " | "),
				strings.Join(header, " | "), "fields."+node.id+".columns", "use exactly the declared columns, in order")
		}
	}
}

const shapeEmpty = "empty"

func (d *diagnosis) checkListPrefix(node *dnode, content string, line int) {
	prefix, ok := d.rules.itemPrefixes[node.id]
	if !ok {
		return
	}
	for index, item := range parseList(content).Content {
		if !prefix.MatchString(item.Value) {
			d.add("list-item-prefix", 3, line, "each item starting "+d.rules.Fields[node.id].ItemPrefix,
				fmt.Sprintf("item %d: %q", index+1, firstLine(item.Value)), "fields."+node.id+".item_prefix",
				"number the item, e.g. \"- G1: ...\"")
		}
	}
}

func firstLine(s string) string { return strings.SplitN(s, "\n", 2)[0] }

func (d *diagnosis) contentText(node *dnode) string {
	if len(node.content) == 0 {
		return ""
	}
	return spanText(d.lines, node.content[0].start, node.content[len(node.content)-1].end)
}

// proseContainment reports RFC 2119 keywords in paragraphs and lists outside
// any req-tagged block; code blocks are not prose.
func (d *diagnosis) proseContainment(node *dnode) {
	for _, b := range node.content {
		switch b.node.(type) {
		case *ast.Paragraph, *ast.List:
		default:
			continue
		}
		for i := b.start; i <= b.end; i++ {
			if keyword := d.rules.keywordRegexp.FindString(d.lines[i]); keyword != "" {
				d.add("prose-containment", 4, i+d.offset+1, "no RFC 2119 keyword outside a req-tagged block",
					"\""+keyword+"\" in "+node.kind+" "+node.id, "prose_containment",
					"move the obligation into a requirement item, or reword it as description")
			}
		}
	}
}
