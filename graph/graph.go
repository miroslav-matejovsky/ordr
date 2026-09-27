package graph

import (
	"fmt"
	"slices"
	"strings"

	"github.com/miroslav-matejovsky/ordr/knowledge"
)

// Dimension is an aspect on which opportunities are compared.
type Dimension string

// Supported comparison dimensions.
const (
	Value       Dimension = "value"       // relative strategic value
	Uncertainty Dimension = "uncertainty" // relative uncertainty
	Complexity  Dimension = "complexity"  // relative complexity
)

func (d Dimension) valid() bool {
	return d == Value || d == Uncertainty || d == Complexity
}

// Comparison states that More ranks above Less on Dimension.
//
// Origin locates the statement in its source, for example
// "graph/ordering.cue:4:3". It is used only for error messages and traceability.
type Comparison struct {
	Dimension Dimension
	More      knowledge.ID
	Less      knowledge.ID
	Origin    string
}

func (c Comparison) String() string {
	return fmt.Sprintf("%s %s > %s", c.Dimension, c.More, c.Less)
}

// RelationKind is the type of a relation between two records.
type RelationKind string

// Supported relation kinds. All are directed from From to To except Relates,
// which is symmetric.
const (
	Supports    RelationKind = "supports"    // From provides evidence for To
	Blocks      RelationKind = "blocks"      // From prevents a decision on To
	Enables     RelationKind = "enables"     // To depends on From
	Invalidates RelationKind = "invalidates" // From contradicts the hypothesis of To
	Relates     RelationKind = "relates"     // untyped, symmetric association
)

// RelationKinds returns all relation kinds in a fixed order.
func RelationKinds() []RelationKind {
	return []RelationKind{Supports, Blocks, Enables, Invalidates, Relates}
}

// Relation links From to To with Kind. Origin is as in Comparison.
type Relation struct {
	Kind   RelationKind
	From   knowledge.ID
	To     knowledge.ID
	Origin string
}

func (r Relation) String() string {
	return fmt.Sprintf("%s %s %s", r.From, r.Kind, r.To)
}

// Graph is a validated set of comparisons and relations over known records.
//
// Invariants: every referenced id is known, no statement references itself,
// no statement is duplicated, and comparisons form no cycle per dimension.
type Graph struct {
	known       []knowledge.ID
	comparisons []Comparison
	relations   []Relation
}

// New validates the statements and returns a Graph.
//
// All problems are reported together in a *ValidationError so the source can
// be corrected in one pass. A statement with an unknown dimension or relation
// kind is a programming error of the caller and is reported as well.
func New(known []knowledge.ID, comparisons []Comparison, relations []Relation) (Graph, error) {
	known = slices.Clone(known)
	slices.Sort(known)
	issues := validate(known, comparisons, relations)
	if len(issues) > 0 {
		return Graph{}, &ValidationError{Issues: issues}
	}
	g := Graph{
		known:       known,
		comparisons: slices.Clone(comparisons),
		relations:   slices.Clone(relations),
	}
	slices.SortFunc(g.comparisons, func(a, b Comparison) int {
		return strings.Compare(a.String(), b.String())
	})
	slices.SortFunc(g.relations, func(a, b Relation) int {
		return strings.Compare(a.String(), b.String())
	})
	return g, nil
}

// Comparisons returns the comparisons on d sorted by More, then Less.
func (g Graph) Comparisons(d Dimension) []Comparison {
	var out []Comparison
	for _, c := range g.comparisons {
		if c.Dimension == d {
			out = append(out, c)
		}
	}
	return out
}

// Relations returns all relations sorted by From, Kind, then To.
func (g Graph) Relations() []Relation {
	return slices.Clone(g.relations)
}

// Ordering is a partial order of records on one dimension, as layers.
//
// Levels[0] holds records nothing ranks above. A record sits one level below
// the lowest level of any record stated or implied above it. Records on the
// same level are incomparable. A lower level alone does not imply that a
// comparison exists between two records on different levels.
//
// Unranked holds known records that appear in no comparison on the dimension.
type Ordering struct {
	Levels   [][]knowledge.ID
	Unranked []knowledge.ID
}

// Order derives the layered partial order on d. Output is deterministic:
// ids within a level and Unranked are sorted.
func (g Graph) Order(d Dimension) Ordering {
	edges := g.Comparisons(d)
	succ := map[knowledge.ID][]knowledge.ID{}
	indegree := map[knowledge.ID]int{}
	for _, c := range edges {
		succ[c.More] = append(succ[c.More], c.Less)
		indegree[c.Less]++
		if _, ok := indegree[c.More]; !ok {
			indegree[c.More] = 0
		}
	}

	var ord Ordering
	for _, id := range g.known {
		if _, ranked := indegree[id]; !ranked {
			ord.Unranked = append(ord.Unranked, id)
		}
	}

	// Kahn's algorithm, removing a whole layer at a time. Terminates because
	// New guarantees the comparisons are acyclic.
	var level []knowledge.ID
	for id, n := range indegree {
		if n == 0 {
			level = append(level, id)
		}
	}
	for len(level) > 0 {
		slices.Sort(level)
		ord.Levels = append(ord.Levels, level)
		var next []knowledge.ID
		for _, id := range level {
			for _, s := range succ[id] {
				indegree[s]--
				if indegree[s] == 0 {
					next = append(next, s)
				}
			}
		}
		level = next
	}
	return ord
}
