package projections

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/miroslav-matejovsky/ordr/graph"
	"github.com/miroslav-matejovsky/ordr/knowledge"
)

// renderFocus names the single opportunity most worth working on now.
//
// The start is the highest ranked opportunity by priority. While the current
// opportunity is blocked by or depends on others, the walk continues with the
// highest ranked of them. The walk stops at an unconstrained opportunity, or
// when the next step would revisit one, since blocks and enables may form
// cycles. No scores are computed; only the stated orderings are used.
func renderFocus(b *strings.Builder, r knowledge.Register, g graph.Graph) {
	fmt.Fprintf(b, "# Focus\n\n%s\n", generatedNote)
	b.WriteString("Which opportunity is the most valuable to work on now?\n\n")
	b.WriteString("Start at the opportunity with the highest strategic value. While it is blocked by or depends on\n")
	b.WriteString("other opportunities, continue with the highest ranked of them. Ties within a value level go to the\n")
	b.WriteString("most advanced state, then to the id.\n\n")

	order, top := priority(r, g)
	if len(order) == 0 {
		b.WriteString("## Next\n\nNo opportunities.\n")
		return
	}
	rank := map[knowledge.ID]int{}
	for i, id := range order {
		rank[id] = i
	}
	rels := g.Relations()

	path := []node{{text: stateLabel(r, order[0])}}
	visited := map[knowledge.ID]bool{order[0]: true}
	current := order[0]
	for {
		var next *graph.Relation
		for _, rel := range rels {
			if rel.To != current || (rel.Kind != graph.Blocks && rel.Kind != graph.Enables) {
				continue
			}
			if next == nil || rank[rel.From] < rank[next.From] {
				next = &rel
			}
		}
		if next == nil || visited[next.From] {
			break
		}
		step := "blocked by"
		if next.Kind == graph.Enables {
			step = "depends on"
		}
		path = append(path, node{text: step + ": " + stateLabel(r, next.From)})
		visited[next.From] = true
		current = next.From
	}

	o, _ := r.Opportunity(current)
	fmt.Fprintf(b, "## Next\n\n%s\n\n", stateLabel(r, current))
	fmt.Fprintf(b, "- Decision to prepare: %s\n", o.Decision)
	fmt.Fprintf(b, "- Current hypothesis: %s\n", o.Hypothesis)
	writeNested(b, "Evidence", o.Evidence, "None recorded.")
	writeNested(b, "Open uncertainties", o.Uncertainties, "None recorded.")
	b.WriteString("\n")

	for i := len(path) - 1; i > 0; i-- {
		path[i-1].children = []node{path[i]}
	}
	b.WriteString("## Path\n\n")
	writeTree(b, path[0])

	b.WriteString("## Relationships and comparisons\n\n")
	writeTree(b, neighbourhood(r, g, rels, current))

	b.WriteString("## Also on the top value level\n\n")
	var peers []string
	for _, id := range top[1:] {
		peers = append(peers, stateLabel(r, id))
	}
	writeList(b, peers, "None.")
	if len(peers) > 0 {
		fmt.Fprintf(b, "Compare these with %s on value to choose the starting point explicitly.\n", label(r, order[0]))
	}
}

// priority orders all opportunities by value level, unranked last. Within a
// level the most advanced state comes first, then the id. It also returns the
// top group in the same order: the first value level, or all unranked
// opportunities when no value comparison exists.
func priority(r knowledge.Register, g graph.Graph) (order, top []knowledge.ID) {
	states := knowledge.States()
	byReadiness := func(a, b knowledge.ID) int {
		oa, _ := r.Opportunity(a)
		ob, _ := r.Opportunity(b)
		return cmp.Or(
			cmp.Compare(slices.Index(states, ob.State), slices.Index(states, oa.State)),
			cmp.Compare(a, b))
	}
	ord := g.Order(graph.Value)
	groups := append(slices.Clone(ord.Levels), ord.Unranked)
	for _, group := range groups {
		group = slices.Clone(group)
		slices.SortFunc(group, byReadiness)
		if top == nil && len(group) > 0 {
			top = group
		}
		order = append(order, group...)
	}
	if len(order) != len(r.IDs()) {
		panic("projections: graph and register disagree on known ids")
	}
	return order, top
}

// neighbourhood groups the direct relations and stated comparisons of id by
// how they read from id. Implied comparisons are left out.
func neighbourhood(r knowledge.Register, g graph.Graph, rels []graph.Relation, id knowledge.ID) node {
	root := node{text: stateLabel(r, id)}
	add := func(title string, match func(graph.Relation) (knowledge.ID, bool)) {
		group := node{text: title}
		for _, rel := range rels {
			if other, ok := match(rel); ok {
				group.children = append(group.children, node{text: stateLabel(r, other)})
			}
		}
		if len(group.children) > 0 {
			root.children = append(root.children, group)
		}
	}
	for _, in := range incoming {
		add(strings.ToLower(in.label), func(rel graph.Relation) (knowledge.ID, bool) {
			return rel.From, rel.Kind == in.kind && rel.To == id
		})
	}
	for _, k := range []graph.RelationKind{graph.Blocks, graph.Invalidates, graph.Enables, graph.Supports, graph.Contains} {
		add(string(k), func(rel graph.Relation) (knowledge.ID, bool) {
			return rel.To, rel.Kind == k && rel.From == id
		})
	}
	add("related to", func(rel graph.Relation) (knowledge.ID, bool) {
		if rel.Kind != graph.Relates {
			return "", false
		}
		switch id {
		case rel.From:
			return rel.To, true
		case rel.To:
			return rel.From, true
		}
		return "", false
	})
	for _, t := range focusComparisons {
		cs := g.Comparisons(t.dim)
		group := node{text: t.more}
		for _, c := range cs {
			if c.More == id {
				group.children = append(group.children, node{text: stateLabel(r, c.Less)})
			}
		}
		if len(group.children) > 0 {
			root.children = append(root.children, group)
		}
		group = node{text: t.less}
		for _, c := range cs {
			if c.Less == id {
				group.children = append(group.children, node{text: stateLabel(r, c.More)})
			}
		}
		if len(group.children) > 0 {
			root.children = append(root.children, group)
		}
	}
	return root
}

// focusComparisons names stated comparisons as read from the focus item.
var focusComparisons = []struct {
	dim        graph.Dimension
	more, less string
}{
	{graph.Value, "more valuable than", "less valuable than"},
	{graph.Uncertainty, "more uncertain than", "less uncertain than"},
	{graph.Complexity, "more complex than", "less complex than"},
}

// node is a line of an ASCII tree.
type node struct {
	text     string
	children []node
}

// writeTree renders n as a fenced ASCII tree.
func writeTree(b *strings.Builder, n node) {
	b.WriteString("```text\n")
	b.WriteString(n.text + "\n")
	writeChildren(b, n.children, "")
	b.WriteString("```\n\n")
}

func writeChildren(b *strings.Builder, children []node, indent string) {
	for i, c := range children {
		branch, next := "├── ", "│   "
		if i == len(children)-1 {
			branch, next = "└── ", "    "
		}
		b.WriteString(indent + branch + c.text + "\n")
		writeChildren(b, c.children, indent+next)
	}
}

func stateLabel(r knowledge.Register, id knowledge.ID) string {
	o, _ := r.Opportunity(id)
	return fmt.Sprintf("%s [%s]", label(r, id), o.State)
}
