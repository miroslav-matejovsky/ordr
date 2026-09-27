package cuegraph_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/miroslav-matejovsky/ordr/graph"
	"github.com/miroslav-matejovsky/ordr/graph/cuegraph"
)

func TestParse(t *testing.T) {
	src := `// ordering
comparisons: {
	value: [
		{more: "alpha", than: "beta"},
		{more: "alpha", than: "gamma"},
	]
	complexity: [{more: "gamma", than: "alpha"}]
}
relations: {
	supports: [{from: "beta", to: "alpha"}]
	enables: [{from: "alpha", to: "gamma"}]
}
`
	comparisons, relations, err := cuegraph.Parse("graph/g.cue", []byte(src))
	require.NoError(t, err)
	require.Equal(t, []graph.Comparison{
		{Dimension: graph.Value, More: "alpha", Less: "beta", Origin: "graph/g.cue:4:3"},
		{Dimension: graph.Value, More: "alpha", Less: "gamma", Origin: "graph/g.cue:5:3"},
		{Dimension: graph.Complexity, More: "gamma", Less: "alpha", Origin: "graph/g.cue:7:15"},
	}, comparisons)
	require.Equal(t, []graph.Relation{
		{Kind: graph.Supports, From: "beta", To: "alpha", Origin: "graph/g.cue:10:13"},
		{Kind: graph.Enables, From: "alpha", To: "gamma", Origin: "graph/g.cue:11:12"},
	}, relations)
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{"syntax", `comparisons: {`, "graph/g.cue"},
		{"unknown dimension", `comparisons: urgency: [{more: "a", than: "b"}]`, "urgency"},
		{"unknown relation", `relations: depends: [{from: "a", to: "b"}]`, "depends"},
		{"unknown field", `comparisons: value: [{more: "a", than: "b", score: 3}]`, "score"},
		{"missing field", `comparisons: value: [{more: "a"}]`, "than"},
		{"numeric id", `comparisons: value: [{more: 1, than: "b"}]`, "more"},
		{"invalid id", `comparisons: value: [{more: "Alpha", than: "b"}]`, `invalid id "Alpha"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := cuegraph.Parse("graph/g.cue", []byte(tt.src))
			require.ErrorContains(t, err, tt.want)
		})
	}
}
