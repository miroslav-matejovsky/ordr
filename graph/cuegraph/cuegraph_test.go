package cuegraph_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/miroslav-matejovsky/ordr/graph"
	"github.com/miroslav-matejovsky/ordr/graph/cuegraph"
)

func TestParse(t *testing.T) {
	src := `// graph
alpha: {
	moreValuableThan: ["beta", "gamma"]
	lessComplexThan: ["gamma"]
	enables: ["gamma"]
	contains: ["beta"]
}
beta: supports: ["alpha"]
"tool-library": lessUncertainThan: ["beta"]
`
	comparisons, relations, err := cuegraph.Parse("graph/g.cue", []byte(src))
	require.NoError(t, err)
	require.Equal(t, []graph.Comparison{
		{Dimension: graph.Value, More: "alpha", Less: "beta", Origin: "graph/g.cue:3:21"},
		{Dimension: graph.Value, More: "alpha", Less: "gamma", Origin: "graph/g.cue:3:29"},
		{Dimension: graph.Complexity, More: "gamma", Less: "alpha", Origin: "graph/g.cue:4:20"},
		{Dimension: graph.Uncertainty, More: "beta", Less: "tool-library", Origin: "graph/g.cue:9:37"},
	}, comparisons)
	require.Equal(t, []graph.Relation{
		{Kind: graph.Enables, From: "alpha", To: "gamma", Origin: "graph/g.cue:5:12"},
		{Kind: graph.Contains, From: "alpha", To: "beta", Origin: "graph/g.cue:6:13"},
		{Kind: graph.Supports, From: "beta", To: "alpha", Origin: "graph/g.cue:8:18"},
	}, relations)
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{"syntax", `alpha: {`, "graph/g.cue"},
		{"unknown statement", `alpha: dependsOn: ["b"]`, "dependsOn"},
		{"record is not a struct", `alpha: ["b"]`, "alpha"},
		{"not a list", `alpha: supports: "b"`, "supports"},
		{"numeric id", `alpha: supports: [1]`, "supports"},
		{"invalid target id", `alpha: supports: ["Beta"]`, `invalid id "Beta"`},
		{"invalid record id", `Alpha: supports: ["b"]`, `invalid id "Alpha"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := cuegraph.Parse("graph/g.cue", []byte(tt.src))
			require.ErrorContains(t, err, tt.want)
		})
	}
}
