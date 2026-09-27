package graph

import (
	"fmt"
	"slices"
	"strings"

	"github.com/miroslav-matejovsky/ordr/knowledge"
)

// IssueKind classifies a graph validation problem.
type IssueKind string

// Validation problem kinds.
const (
	UnknownReference        IssueKind = "unknown-reference"
	SelfReference           IssueKind = "self-reference"
	DuplicateStatement      IssueKind = "duplicate"
	ContradictoryComparison IssueKind = "contradiction"
	OrderingCycle           IssueKind = "cycle"
	InvalidStatement        IssueKind = "invalid"
)

// Issue is one validation problem. Message names the statement, its origin
// and the ids involved so the source can be corrected.
type Issue struct {
	Kind    IssueKind
	Message string
}

// ValidationError lists every problem found by New.
type ValidationError struct {
	Issues []Issue
}

func (e *ValidationError) Error() string {
	lines := make([]string, 0, len(e.Issues)+1)
	lines = append(lines, fmt.Sprintf("invalid graph: %d issue(s)", len(e.Issues)))
	for _, i := range e.Issues {
		lines = append(lines, fmt.Sprintf("  [%s] %s", i.Kind, i.Message))
	}
	return strings.Join(lines, "\n")
}

type edgeKey struct {
	dim  Dimension
	more knowledge.ID
	less knowledge.ID
}

func validate(known []knowledge.ID, comparisons []Comparison, relations []Relation) []Issue {
	var issues []Issue
	add := func(kind IssueKind, format string, args ...any) {
		issues = append(issues, Issue{Kind: kind, Message: fmt.Sprintf(format, args...)})
	}
	isKnown := func(id knowledge.ID) bool {
		_, found := slices.BinarySearch(known, id)
		return found
	}
	// checkRefs reports unknown and self references and tells whether the
	// statement may take part in further checks.
	checkRefs := func(origin, stmt string, a, b knowledge.ID) bool {
		ok := true
		for _, id := range []knowledge.ID{a, b} {
			if !isKnown(id) {
				add(UnknownReference, "%s: %q references unknown knowledge id %q", origin, stmt, id)
				ok = false
			}
		}
		if a == b {
			add(SelfReference, "%s: %q references %q on both sides", origin, stmt, a)
			ok = false
		}
		return ok
	}

	edges := map[edgeKey]Comparison{}
	var order []edgeKey
	for _, c := range comparisons {
		if !c.Dimension.valid() {
			add(InvalidStatement, "%s: unknown dimension %q", c.Origin, c.Dimension)
			continue
		}
		if !checkRefs(c.Origin, c.String(), c.More, c.Less) {
			continue
		}
		k := edgeKey{c.Dimension, c.More, c.Less}
		if prev, dup := edges[k]; dup {
			add(DuplicateStatement, "%s: %q already stated at %s", c.Origin, c.String(), prev.Origin)
			continue
		}
		edges[k] = c
		order = append(order, k)
	}

	for _, k := range order {
		reverse, found := edges[edgeKey{k.dim, k.less, k.more}]
		if found && k.more < k.less {
			c := edges[k]
			add(ContradictoryComparison, "%q at %s contradicts %q at %s",
				c.String(), c.Origin, reverse.String(), reverse.Origin)
		}
	}

	for _, d := range []Dimension{Value, Uncertainty, Complexity} {
		issues = append(issues, findCycles(d, edges)...)
	}

	seen := map[Relation]Relation{}
	for _, r := range relations {
		if !slices.Contains(RelationKinds(), r.Kind) {
			add(InvalidStatement, "%s: unknown relation kind %q", r.Origin, r.Kind)
			continue
		}
		if !checkRefs(r.Origin, r.String(), r.From, r.To) {
			continue
		}
		k := Relation{Kind: r.Kind, From: r.From, To: r.To}
		if r.Kind == Relates && k.To < k.From {
			k.From, k.To = k.To, k.From
		}
		if prev, dup := seen[k]; dup {
			add(DuplicateStatement, "%s: %q already stated at %s", r.Origin, r.String(), prev.Origin)
			continue
		}
		seen[k] = r
	}
	return issues
}

// findCycles reports every strongly connected component of the comparisons
// on d with more than two members, listing the comparisons inside it. Two
// member components are direct contradictions and reported separately.
func findCycles(d Dimension, edges map[edgeKey]Comparison) []Issue {
	succ := map[knowledge.ID][]knowledge.ID{}
	var nodes []knowledge.ID
	addNode := func(id knowledge.ID) {
		if _, ok := succ[id]; !ok {
			succ[id] = nil
			nodes = append(nodes, id)
		}
	}
	for k := range edges {
		if k.dim != d {
			continue
		}
		addNode(k.more)
		addNode(k.less)
		succ[k.more] = append(succ[k.more], k.less)
	}
	slices.Sort(nodes)
	for id := range succ {
		slices.Sort(succ[id])
	}

	var issues []Issue
	for _, scc := range stronglyConnected(nodes, succ) {
		if len(scc) <= 2 {
			continue
		}
		var inside []string
		for _, from := range scc {
			for _, to := range succ[from] {
				if slices.Contains(scc, to) {
					c := edges[edgeKey{d, from, to}]
					inside = append(inside, fmt.Sprintf("%s > %s (%s)", c.More, c.Less, c.Origin))
				}
			}
		}
		issues = append(issues, Issue{
			Kind: OrderingCycle,
			Message: fmt.Sprintf("%s comparisons form a cycle among %s; remove or reverse one of: %s",
				d, joinIDs(scc), strings.Join(inside, ", ")),
		})
	}
	return issues
}

// stronglyConnected returns the strongly connected components using Tarjan's
// algorithm. Each component is sorted; components are sorted by first id.
func stronglyConnected(nodes []knowledge.ID, succ map[knowledge.ID][]knowledge.ID) [][]knowledge.ID {
	index := map[knowledge.ID]int{}
	low := map[knowledge.ID]int{}
	onStack := map[knowledge.ID]bool{}
	var stack []knowledge.ID
	var out [][]knowledge.ID
	next := 0

	var visit func(v knowledge.ID)
	visit = func(v knowledge.ID) {
		index[v], low[v] = next, next
		next++
		stack = append(stack, v)
		onStack[v] = true
		for _, w := range succ[v] {
			if _, seen := index[w]; !seen {
				visit(w)
				low[v] = min(low[v], low[w])
			} else if onStack[w] {
				low[v] = min(low[v], index[w])
			}
		}
		if low[v] != index[v] {
			return
		}
		var scc []knowledge.ID
		for {
			w := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			onStack[w] = false
			scc = append(scc, w)
			if w == v {
				break
			}
		}
		slices.Sort(scc)
		out = append(out, scc)
	}
	for _, v := range nodes {
		if _, seen := index[v]; !seen {
			visit(v)
		}
	}
	slices.SortFunc(out, func(a, b []knowledge.ID) int { return strings.Compare(string(a[0]), string(b[0])) })
	return out
}

func joinIDs(ids []knowledge.ID) string {
	s := make([]string, len(ids))
	for i, id := range ids {
		s[i] = string(id)
	}
	return strings.Join(s, ", ")
}
