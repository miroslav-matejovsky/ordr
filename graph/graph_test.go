package graph_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/miroslav-matejovsky/ordr/graph"
	"github.com/miroslav-matejovsky/ordr/knowledge"
)

var known = []knowledge.ID{"alpha", "beta", "delta", "gamma"}

func cmp(d graph.Dimension, more, less knowledge.ID, origin string) graph.Comparison {
	return graph.Comparison{Dimension: d, More: more, Less: less, Origin: origin}
}

func rel(k graph.RelationKind, from, to knowledge.ID, origin string) graph.Relation {
	return graph.Relation{Kind: k, From: from, To: to, Origin: origin}
}

func TestOrder(t *testing.T) {
	tests := []struct {
		name        string
		comparisons []graph.Comparison
		want        graph.Ordering
	}{
		{
			name: "transitive chain",
			comparisons: []graph.Comparison{
				cmp(graph.Value, "beta", "gamma", "f:2"),
				cmp(graph.Value, "alpha", "beta", "f:1"),
			},
			want: graph.Ordering{
				Levels:   [][]knowledge.ID{{"alpha"}, {"beta"}, {"gamma"}},
				Unranked: []knowledge.ID{"delta"},
			},
		},
		{
			name: "incomparable items share a level",
			comparisons: []graph.Comparison{
				cmp(graph.Value, "alpha", "gamma", "f:1"),
				cmp(graph.Value, "alpha", "beta", "f:2"),
			},
			want: graph.Ordering{
				Levels:   [][]knowledge.ID{{"alpha"}, {"beta", "gamma"}},
				Unranked: []knowledge.ID{"delta"},
			},
		},
		{
			name: "longest chain decides the level",
			comparisons: []graph.Comparison{
				cmp(graph.Value, "alpha", "gamma", "f:1"),
				cmp(graph.Value, "alpha", "beta", "f:2"),
				cmp(graph.Value, "beta", "gamma", "f:3"),
				cmp(graph.Value, "delta", "gamma", "f:4"),
			},
			want: graph.Ordering{
				Levels: [][]knowledge.ID{{"alpha", "delta"}, {"beta"}, {"gamma"}},
			},
		},
		{
			name:        "other dimensions are ignored",
			comparisons: []graph.Comparison{cmp(graph.Complexity, "alpha", "beta", "f:1")},
			want: graph.Ordering{
				Unranked: []knowledge.ID{"alpha", "beta", "delta", "gamma"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := graph.New(known, tt.comparisons, nil)
			require.NoError(t, err)
			require.Equal(t, tt.want, g.Order(graph.Value))
		})
	}
}

func TestNewReportsIssues(t *testing.T) {
	tests := []struct {
		name        string
		comparisons []graph.Comparison
		relations   []graph.Relation
		want        []graph.Issue
	}{
		{
			name:        "unknown reference in comparison",
			comparisons: []graph.Comparison{cmp(graph.Value, "alpha", "omega", "f:1")},
			want: []graph.Issue{{Kind: graph.UnknownReference,
				Message: `f:1: "value alpha > omega" references unknown knowledge id "omega"`}},
		},
		{
			name:      "unknown reference in relation",
			relations: []graph.Relation{rel(graph.Supports, "omega", "alpha", "f:1")},
			want: []graph.Issue{{Kind: graph.UnknownReference,
				Message: `f:1: "omega supports alpha" references unknown knowledge id "omega"`}},
		},
		{
			name:        "self reference in comparison",
			comparisons: []graph.Comparison{cmp(graph.Value, "alpha", "alpha", "f:1")},
			want: []graph.Issue{{Kind: graph.SelfReference,
				Message: `f:1: "value alpha > alpha" references "alpha" on both sides`}},
		},
		{
			name:      "self reference in relation",
			relations: []graph.Relation{rel(graph.Blocks, "beta", "beta", "f:1")},
			want: []graph.Issue{{Kind: graph.SelfReference,
				Message: `f:1: "beta blocks beta" references "beta" on both sides`}},
		},
		{
			name: "duplicate comparison",
			comparisons: []graph.Comparison{
				cmp(graph.Value, "alpha", "beta", "f:1"),
				cmp(graph.Value, "alpha", "beta", "f:2"),
			},
			want: []graph.Issue{{Kind: graph.DuplicateStatement,
				Message: `f:2: "value alpha > beta" already stated at f:1`}},
		},
		{
			name: "same pair on different dimensions is not a duplicate",
			comparisons: []graph.Comparison{
				cmp(graph.Value, "alpha", "beta", "f:1"),
				cmp(graph.Complexity, "alpha", "beta", "f:2"),
			},
		},
		{
			name: "duplicate relation",
			relations: []graph.Relation{
				rel(graph.Enables, "alpha", "beta", "f:1"),
				rel(graph.Enables, "alpha", "beta", "f:2"),
			},
			want: []graph.Issue{{Kind: graph.DuplicateStatement,
				Message: `f:2: "alpha enables beta" already stated at f:1`}},
		},
		{
			name: "relates is symmetric",
			relations: []graph.Relation{
				rel(graph.Relates, "alpha", "beta", "f:1"),
				rel(graph.Relates, "beta", "alpha", "f:2"),
			},
			want: []graph.Issue{{Kind: graph.DuplicateStatement,
				Message: `f:2: "beta relates alpha" already stated at f:1`}},
		},
		{
			name: "contains forms a forest",
			relations: []graph.Relation{
				rel(graph.Contains, "alpha", "beta", "f:1"),
				rel(graph.Contains, "alpha", "gamma", "f:2"),
				rel(graph.Contains, "beta", "delta", "f:3"),
			},
		},
		{
			name: "second parent",
			relations: []graph.Relation{
				rel(graph.Contains, "alpha", "gamma", "f:1"),
				rel(graph.Contains, "beta", "gamma", "f:2"),
			},
			want: []graph.Issue{{Kind: graph.MultipleParents,
				Message: `f:2: "beta contains gamma" gives "gamma" a second parent; "alpha contains gamma" already states its parent at f:1`}},
		},
		{
			name: "containment cycle",
			relations: []graph.Relation{
				rel(graph.Contains, "alpha", "beta", "f:1"),
				rel(graph.Contains, "beta", "gamma", "f:2"),
				rel(graph.Contains, "gamma", "alpha", "f:3"),
				rel(graph.Contains, "gamma", "delta", "f:4"),
			},
			want: []graph.Issue{{Kind: graph.ContainmentCycle,
				Message: "contains relations form a cycle among alpha, beta, gamma; " +
					"remove one of: gamma contains alpha (f:3), alpha contains beta (f:1), beta contains gamma (f:2)"}},
		},
		{
			name: "two record containment cycle",
			relations: []graph.Relation{
				rel(graph.Contains, "alpha", "beta", "f:1"),
				rel(graph.Contains, "beta", "alpha", "f:2"),
			},
			want: []graph.Issue{{Kind: graph.ContainmentCycle,
				Message: "contains relations form a cycle among alpha, beta; " +
					"remove one of: beta contains alpha (f:2), alpha contains beta (f:1)"}},
		},
		{
			name: "contradictory comparisons",
			comparisons: []graph.Comparison{
				cmp(graph.Uncertainty, "beta", "alpha", "f:1"),
				cmp(graph.Uncertainty, "alpha", "beta", "f:2"),
			},
			want: []graph.Issue{{Kind: graph.ContradictoryComparison,
				Message: `"uncertainty alpha > beta" at f:2 contradicts "uncertainty beta > alpha" at f:1`}},
		},
		{
			name: "cycle",
			comparisons: []graph.Comparison{
				cmp(graph.Value, "alpha", "beta", "f:1"),
				cmp(graph.Value, "beta", "gamma", "f:2"),
				cmp(graph.Value, "gamma", "alpha", "f:3"),
				cmp(graph.Value, "delta", "alpha", "f:4"),
			},
			want: []graph.Issue{{Kind: graph.OrderingCycle,
				Message: "value comparisons form a cycle among alpha, beta, gamma; " +
					"remove or reverse one of: alpha > beta (f:1), beta > gamma (f:2), gamma > alpha (f:3)"}},
		},
		{
			name: "several issues are reported together",
			comparisons: []graph.Comparison{
				cmp(graph.Value, "alpha", "omega", "f:1"),
				cmp(graph.Value, "beta", "beta", "f:2"),
			},
			want: []graph.Issue{
				{Kind: graph.UnknownReference, Message: `f:1: "value alpha > omega" references unknown knowledge id "omega"`},
				{Kind: graph.SelfReference, Message: `f:2: "value beta > beta" references "beta" on both sides`},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := graph.New(known, tt.comparisons, tt.relations)
			if tt.want == nil {
				require.NoError(t, err)
				return
			}
			var verr *graph.ValidationError
			require.True(t, errors.As(err, &verr), "want ValidationError, got %v", err)
			require.Equal(t, tt.want, verr.Issues)
		})
	}
}

func TestRelationsAndComparisonsAreSorted(t *testing.T) {
	g, err := graph.New(known,
		[]graph.Comparison{
			cmp(graph.Value, "gamma", "delta", "f:1"),
			cmp(graph.Value, "alpha", "beta", "f:2"),
		},
		[]graph.Relation{
			rel(graph.Supports, "gamma", "alpha", "f:3"),
			rel(graph.Blocks, "alpha", "beta", "f:4"),
		})
	require.NoError(t, err)
	require.Equal(t, []graph.Comparison{
		cmp(graph.Value, "alpha", "beta", "f:2"),
		cmp(graph.Value, "gamma", "delta", "f:1"),
	}, g.Comparisons(graph.Value))
	require.Equal(t, []graph.Relation{
		rel(graph.Blocks, "alpha", "beta", "f:4"),
		rel(graph.Supports, "gamma", "alpha", "f:3"),
	}, g.Relations())
}
