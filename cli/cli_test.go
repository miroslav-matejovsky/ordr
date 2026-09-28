package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/miroslav-matejovsky/ordr/cli"
)

const basic = "../examples/basic"

// copyWorkspace copies the knowledge and graph sources of an example
// workspace into a temporary folder so tests never write into the repository.
func copyWorkspace(t *testing.T, src string) string {
	t.Helper()
	dst := t.TempDir()
	for _, pattern := range []string{"ordr.yaml", "knowledge/*.md", "graph/*.cue"} {
		matches, err := filepath.Glob(filepath.Join(src, pattern))
		require.NoError(t, err)
		require.NotEmpty(t, matches, pattern)
		for _, m := range matches {
			rel, err := filepath.Rel(src, m)
			require.NoError(t, err)
			data, err := os.ReadFile(m)
			require.NoError(t, err)
			require.NoError(t, os.MkdirAll(filepath.Join(dst, filepath.Dir(rel)), 0o755))
			require.NoError(t, os.WriteFile(filepath.Join(dst, rel), data, 0o644))
		}
	}
	return dst
}

func readNormalized(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return strings.ReplaceAll(string(data), "\r\n", "\n")
}

// TestProjectValidExamples generates all projections for each valid example
// workspace and compares them with the committed projections of the example.
func TestProjectValidExamples(t *testing.T) {
	for _, example := range []string{basic, "../examples/invoicing"} {
		t.Run(filepath.Base(example), func(t *testing.T) {
			ws := copyWorkspace(t, example)
			var stdout, stderr bytes.Buffer

			code := cli.Run([]string{"project", "--workspace", ws}, &stdout, &stderr)

			require.Equal(t, cli.ExitOK, code, stderr.String())
			require.Equal(t, "wrote projections/value.md\n"+
				"wrote projections/uncertainty.md\n"+
				"wrote projections/complexity.md\n"+
				"wrote projections/relationships.md\n"+
				"wrote projections/readiness.md\n", stdout.String())
			for _, name := range []string{"value.md", "uncertainty.md", "complexity.md", "relationships.md", "readiness.md"} {
				want := readNormalized(t, filepath.Join(example, "projections", name))
				got := readNormalized(t, filepath.Join(ws, "projections", name))
				require.Equal(t, want, got, "%s differs from %s; regenerate with: go run . project --workspace %s",
					name, example, strings.TrimPrefix(example, "../"))
			}
		})
	}
}

func TestProjectSingleToStdout(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := cli.Run([]string{"project", "value", "--workspace", basic, "--stdout"}, &stdout, &stderr)

	require.Equal(t, cli.ExitOK, code, stderr.String())
	require.Equal(t, readNormalized(t, filepath.Join(basic, "projections", "value.md")), stdout.String())
}

func TestProjectInvalidExample(t *testing.T) {
	ws := copyWorkspace(t, "../examples/invalid")
	var stdout, stderr bytes.Buffer

	code := cli.Run([]string{"project", "--workspace", ws}, &stdout, &stderr)

	require.Equal(t, cli.ExitError, code)
	require.Empty(t, stdout.String())
	require.Contains(t, stderr.String(), `[unknown-reference] graph/graph.cue:2:35: "value alpha > delta" references unknown knowledge id "delta"`)
	require.Contains(t, stderr.String(), "[cycle] value comparisons form a cycle among alpha, beta, gamma")
	require.NoDirExists(t, filepath.Join(ws, "projections"))
}

func TestUsageErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"no command", nil, "missing command"},
		{"unknown command", []string{"rank"}, `unknown command "rank"`},
		{"unknown projection", []string{"project", "score"}, `unknown projection "score"`},
		{"unknown flag", []string{"project", "--verbose"}, "verbose"},
		{"extra argument", []string{"project", "value", "extra"}, "unexpected arguments [extra]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := cli.Run(tt.args, &stdout, &stderr)
			require.Equal(t, cli.ExitUsage, code)
			require.Contains(t, stderr.String(), tt.want)
			require.Contains(t, stderr.String(), "usage: ordr project")
		})
	}
}
