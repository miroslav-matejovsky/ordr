// Package workspace isolates the filesystem from the ORDR domain.
//
// A workspace is a directory holding an ordr.yaml configuration file that
// names the knowledge, graph and projections folders explicitly:
//
//	knowledge: knowledge
//	graph: graph
//	projections: projections
//
// All three keys are required and non-empty. Unknown keys are rejected.
//
// The package discovers a workspace, reads its source files, hands their
// content to the knowledge and graph packages, and writes generated
// projections. Sources are never modified.
//
// Source files are read non-recursively and in lexical order: *.md from the
// knowledge folder and *.cue from the graph folder. Paths in error messages
// and statement origins are workspace relative with forward slashes.
package workspace

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"go.yaml.in/yaml/v3"

	"github.com/miroslav-matejovsky/ordr/graph"
	"github.com/miroslav-matejovsky/ordr/graph/cuegraph"
	"github.com/miroslav-matejovsky/ordr/knowledge"
)

// ConfigFile is the name of the file that marks a workspace root.
const ConfigFile = "ordr.yaml"

// config mirrors ConfigFile. Every field is required.
type config struct {
	Knowledge   string `yaml:"knowledge"`
	Graph       string `yaml:"graph"`
	Projections string `yaml:"projections"`
}

// Workspace is a discovered ORDR workspace. Folder fields are workspace
// relative, slash separated and guaranteed to stay inside the root.
type Workspace struct {
	root        string
	knowledge   string
	graph       string
	projections string
}

// Discover finds the workspace containing start by walking up the directory
// tree to the first folder holding ConfigFile, then loads its configuration.
func Discover(start string) (Workspace, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return Workspace{}, fmt.Errorf("resolve %q: %w", start, err)
	}
	for {
		_, err := os.Stat(filepath.Join(dir, ConfigFile))
		if err == nil {
			return open(dir)
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return Workspace{}, err
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return Workspace{}, fmt.Errorf("no %s found in %q or any parent folder", ConfigFile, start)
		}
		dir = parent
	}
}

func open(root string) (Workspace, error) {
	data, err := os.ReadFile(filepath.Join(root, ConfigFile))
	if err != nil {
		return Workspace{}, err
	}
	var cfg config
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&cfg); err != nil && !errors.Is(err, io.EOF) {
		return Workspace{}, fmt.Errorf("%s: %w", ConfigFile, err)
	}
	w := Workspace{root: root, knowledge: cfg.Knowledge, graph: cfg.Graph, projections: cfg.Projections}
	folders := []struct {
		key, dir  string
		mustExist bool // projections are created on first write
	}{
		{"knowledge", w.knowledge, true},
		{"graph", w.graph, true},
		{"projections", w.projections, false},
	}
	for _, f := range folders {
		if f.dir == "" {
			return Workspace{}, fmt.Errorf("%s: %s: missing", ConfigFile, f.key)
		}
		if !filepath.IsLocal(f.dir) {
			return Workspace{}, fmt.Errorf("%s: %s: %q must be a relative path inside the workspace", ConfigFile, f.key, f.dir)
		}
	}
	for _, f := range folders {
		if !f.mustExist {
			continue
		}
		info, err := os.Stat(filepath.Join(root, f.dir))
		if err != nil || !info.IsDir() {
			return Workspace{}, fmt.Errorf("%s: %s: folder %q not found", ConfigFile, f.key, f.dir)
		}
	}
	return w, nil
}

// Load reads all sources and returns the validated knowledge register and graph.
func (w Workspace) Load() (knowledge.Register, graph.Graph, error) {
	var opportunities []knowledge.Opportunity
	err := w.eachFile(w.knowledge, "*.md", func(name string, data []byte) error {
		o, err := knowledge.ParseOpportunity(name, data)
		if err != nil {
			return err
		}
		opportunities = append(opportunities, o)
		return nil
	})
	if err != nil {
		return knowledge.Register{}, graph.Graph{}, err
	}
	register, err := knowledge.NewRegister(opportunities)
	if err != nil {
		return knowledge.Register{}, graph.Graph{}, err
	}

	var comparisons []graph.Comparison
	var relations []graph.Relation
	err = w.eachFile(w.graph, "*.cue", func(name string, data []byte) error {
		c, r, err := cuegraph.Parse(name, data)
		if err != nil {
			return err
		}
		comparisons = append(comparisons, c...)
		relations = append(relations, r...)
		return nil
	})
	if err != nil {
		return knowledge.Register{}, graph.Graph{}, err
	}
	g, err := graph.New(register.IDs(), comparisons, relations)
	if err != nil {
		return knowledge.Register{}, graph.Graph{}, err
	}
	return register, g, nil
}

// eachFile calls fn with the workspace relative name and content of every
// file in dir matching pattern, in lexical order.
func (w Workspace) eachFile(dir, pattern string, fn func(name string, data []byte) error) error {
	matches, err := filepath.Glob(filepath.Join(w.root, dir, pattern))
	if err != nil {
		return err
	}
	for _, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			return err
		}
		if err := fn(path.Join(filepath.ToSlash(dir), filepath.Base(m)), data); err != nil {
			return err
		}
	}
	return nil
}

// WriteProjection writes content to the projections folder, creating it if
// needed, and returns the workspace relative path written.
func (w Workspace) WriteProjection(name, content string) (string, error) {
	if name != filepath.Base(name) || !filepath.IsLocal(name) {
		return "", fmt.Errorf("projection name %q must be a plain file name", name)
	}
	dir := filepath.Join(w.root, w.projections)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		return "", err
	}
	return path.Join(filepath.ToSlash(w.projections), name), nil
}
