// Package cuegraph reads graph statements written in CUE.
//
// CUE is the textual notation of the graph context, not its model. This
// package translates a CUE document into [graph.Comparison] and
// [graph.Relation] values; validation of references and ordering happens in
// [graph.New].
//
// Notation (every field is optional, unknown fields are errors):
//
//	comparisons: {
//		value:       [{more: "alpha", than: "beta"}]
//		uncertainty: [{more: "gamma", than: "beta"}]
//		complexity:  [{more: "gamma", than: "alpha"}]
//	}
//	relations: {
//		supports:    [{from: "beta", to: "alpha"}]
//		blocks:      [...]
//		enables:     [...]
//		invalidates: [...]
//		relates:     [...]
//	}
//
// Each file is read on its own; statements from several files are simply
// concatenated by the caller. Lists do not unify across files.
package cuegraph

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	cueerrors "cuelang.org/go/cue/errors"

	"github.com/miroslav-matejovsky/ordr/graph"
	"github.com/miroslav-matejovsky/ordr/knowledge"
)

const schema = `
#Comparison: {more: string, than: string}
#Relation: {from: string, to: string}
#Graph: {
	comparisons?: {
		value?:       [...#Comparison]
		uncertainty?: [...#Comparison]
		complexity?:  [...#Comparison]
	}
	relations?: {
		supports?:    [...#Relation]
		blocks?:      [...#Relation]
		enables?:     [...#Relation]
		invalidates?: [...#Relation]
		relates?:     [...#Relation]
	}
}
`

var dimensions = []graph.Dimension{graph.Value, graph.Uncertainty, graph.Complexity}

// Parse reads one CUE graph file. filename is used for positions in errors
// and in the Origin of each statement.
func Parse(filename string, data []byte) ([]graph.Comparison, []graph.Relation, error) {
	ctx := cuecontext.New()
	def := ctx.CompileString(schema, cue.Filename("ordr-graph-schema.cue")).LookupPath(cue.ParsePath("#Graph"))
	if err := def.Err(); err != nil {
		return nil, nil, fmt.Errorf("graph schema: %w", err)
	}
	doc := ctx.CompileBytes(data, cue.Filename(filename))
	if err := doc.Err(); err != nil {
		return nil, nil, cueError(err)
	}
	// Validate against the schema, but read statements from doc: values of the
	// unified result carry schema positions instead of source positions.
	v := def.Unify(doc)
	if err := v.Validate(cue.Concrete(true)); err != nil {
		return nil, nil, cueError(err)
	}

	var comparisons []graph.Comparison
	for _, d := range dimensions {
		err := eachElement(doc, "comparisons."+string(d), func(pos string, e cue.Value) error {
			var raw struct {
				More string `json:"more"`
				Than string `json:"than"`
			}
			if err := e.Decode(&raw); err != nil {
				return err
			}
			more, err := knowledge.ParseID(raw.More)
			if err != nil {
				return err
			}
			less, err := knowledge.ParseID(raw.Than)
			if err != nil {
				return err
			}
			comparisons = append(comparisons, graph.Comparison{Dimension: d, More: more, Less: less, Origin: pos})
			return nil
		})
		if err != nil {
			return nil, nil, err
		}
	}

	var relations []graph.Relation
	for _, k := range graph.RelationKinds() {
		err := eachElement(doc, "relations."+string(k), func(pos string, e cue.Value) error {
			var raw struct {
				From string `json:"from"`
				To   string `json:"to"`
			}
			if err := e.Decode(&raw); err != nil {
				return err
			}
			from, err := knowledge.ParseID(raw.From)
			if err != nil {
				return err
			}
			to, err := knowledge.ParseID(raw.To)
			if err != nil {
				return err
			}
			relations = append(relations, graph.Relation{Kind: k, From: from, To: to, Origin: pos})
			return nil
		})
		if err != nil {
			return nil, nil, err
		}
	}
	return comparisons, relations, nil
}

// eachElement calls fn for every element of the list at path, if present.
// Errors from fn are prefixed with the element position.
func eachElement(v cue.Value, path string, fn func(pos string, e cue.Value) error) error {
	list := v.LookupPath(cue.ParsePath(path))
	if !list.Exists() {
		return nil
	}
	it, err := list.List()
	if err != nil {
		return cueError(err)
	}
	for it.Next() {
		pos := it.Value().Pos().String()
		if err := fn(pos, it.Value()); err != nil {
			return fmt.Errorf("%s: %w", pos, err)
		}
	}
	return nil
}

// cueError flattens a CUE error list into one error with positions.
func cueError(err error) error {
	return fmt.Errorf("%s", cueerrors.Details(err, nil))
}
