package workspace_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/miroslav-matejovsky/ordr/graph"
	"github.com/miroslav-matejovsky/ordr/knowledge"
	"github.com/miroslav-matejovsky/ordr/workspace"
)

const config = `knowledge: k
graph: g
projections: out
`

const record = `---
id: %s
title: T
state: framing
---
## Summary
S
## Hypothesis
H
## Evidence
## Uncertainties
## Decision
D
`

// writeTree creates files under a temporary root and returns it.
func writeTree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, content := range files {
		p := filepath.Join(root, filepath.FromSlash(name))
		require.NoError(t, os.MkdirAll(filepath.Dir(p), 0o755))
		require.NoError(t, os.WriteFile(p, []byte(content), 0o644))
	}
	return root
}

func rec(id string) string {
	return fmt.Sprintf(record, id)
}

func TestDiscoverWalksUp(t *testing.T) {
	root := writeTree(t, map[string]string{
		"ordr.yaml":    config,
		"k/a.md":       rec("a"),
		"g/x.cue":      "",
		"k/deep/.keep": "",
	})
	ws, err := workspace.Discover(filepath.Join(root, "k", "deep"))
	require.NoError(t, err)
	r, _, err := ws.Load()
	require.NoError(t, err)
	require.Equal(t, []knowledge.ID{"a"}, r.IDs())
}

func TestDiscoverErrors(t *testing.T) {
	tests := []struct {
		name  string
		files map[string]string
		want  string
	}{
		{"no config", map[string]string{"k/a.md": ""}, "no ordr.yaml found"},
		{"missing key", map[string]string{"ordr.yaml": "knowledge: k\ngraph: g\n"}, "ordr.yaml: projections: missing"},
		{"empty value", map[string]string{"ordr.yaml": "knowledge: k\ngraph: \"\"\nprojections: out\n"}, "ordr.yaml: graph: missing"},
		{"not a mapping", map[string]string{"ordr.yaml": "- k\n"}, "ordr.yaml"},
		{"unknown key", map[string]string{"ordr.yaml": config + "extra: x\n"}, "extra"},
		{"escaping path", map[string]string{"ordr.yaml": "knowledge: ../k\ngraph: g\nprojections: out\n"},
			`knowledge: "../k" must be a relative path`},
		{"missing folder", map[string]string{"ordr.yaml": config, "g/x.cue": ""}, `knowledge: folder "k" not found`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := workspace.Discover(writeTree(t, tt.files))
			require.ErrorContains(t, err, tt.want)
		})
	}
}

func TestLoadCombinesGraphFiles(t *testing.T) {
	root := writeTree(t, map[string]string{
		"ordr.yaml":  config,
		"k/a.md":     rec("a"),
		"k/b.md":     rec("b"),
		"k/c.md":     rec("c"),
		"k/note.txt": "ignored",
		"g/1.cue":    `a: moreValuableThan: ["b"]`,
		"g/2.cue":    `b: moreValuableThan: ["c"]` + "\nc: blocks: [\"a\"]",
	})
	ws, err := workspace.Discover(root)
	require.NoError(t, err)
	_, g, err := ws.Load()
	require.NoError(t, err)
	require.Equal(t, [][]knowledge.ID{{"a"}, {"b"}, {"c"}}, g.Order(graph.Value).Levels)
	require.Equal(t, []graph.Relation{{Kind: graph.Blocks, From: "c", To: "a", Origin: "g/2.cue:2:13"}}, g.Relations())
}

func TestLoadErrors(t *testing.T) {
	tests := []struct {
		name  string
		files map[string]string
		want  string
	}{
		{"bad record", map[string]string{"k/a.md": "no front matter"}, "k/a.md: front matter"},
		{"duplicate id", map[string]string{"k/a.md": rec("a"), "k/b.md": rec("a")}, `duplicate id "a" in k/a.md and k/b.md`},
		{"bad cue", map[string]string{"k/a.md": rec("a"), "g/x.cue": "a: {"}, "g/x.cue"},
		{"unknown reference", map[string]string{"k/a.md": rec("a"), "g/x.cue": `a: moreValuableThan: ["z"]`},
			`unknown knowledge id "z"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.files["ordr.yaml"] = config
			if _, ok := tt.files["g/x.cue"]; !ok {
				tt.files["g/x.cue"] = ""
			}
			ws, err := workspace.Discover(writeTree(t, tt.files))
			require.NoError(t, err)
			_, _, err = ws.Load()
			require.ErrorContains(t, err, tt.want)
		})
	}
}

func TestLoadReturnsValidationError(t *testing.T) {
	root := writeTree(t, map[string]string{
		"ordr.yaml": config,
		"k/a.md":    rec("a"),
		"g/x.cue":   `a: moreValuableThan: ["a"]`,
	})
	ws, err := workspace.Discover(root)
	require.NoError(t, err)
	_, _, err = ws.Load()
	var verr *graph.ValidationError
	require.True(t, errors.As(err, &verr))
	require.Equal(t, graph.SelfReference, verr.Issues[0].Kind)
}

func TestWriteProjection(t *testing.T) {
	root := writeTree(t, map[string]string{"ordr.yaml": config, "k/.keep": "", "g/.keep": ""})
	ws, err := workspace.Discover(root)
	require.NoError(t, err)

	written, err := ws.WriteProjection("value.md", "content\n")
	require.NoError(t, err)
	require.Equal(t, "out/value.md", written)
	data, err := os.ReadFile(filepath.Join(root, "out", "value.md"))
	require.NoError(t, err)
	require.Equal(t, "content\n", string(data))

	_, err = ws.WriteProjection("../k/a.md", "x")
	require.ErrorContains(t, err, "plain file name")
}
