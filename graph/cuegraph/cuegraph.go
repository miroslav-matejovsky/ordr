// Package cuegraph reads graph statements written in CUE.
//
// CUE is the textual notation of the graph context, not its model. This
// package translates a CUE document into [graph.Comparison] and
// [graph.Relation] values; validation of references and ordering happens in
// [graph.New].
//
// Notation: every top-level field is a knowledge id and holds the statements
// made from that record's point of view. Every statement field is optional,
// unknown fields are errors. Ids containing hyphens must be quoted.
//
//	alpha: {
//		supports:    ["beta"]  // alpha supports beta
//		blocks:      ["gamma"] // alpha blocks gamma
//		enables:     ["delta"] // delta depends on alpha
//		invalidates: [...]
//		contains:    [...]     // children of alpha
//		relates:     [...]     // symmetric
//
//		moreValuableThan:  ["beta"] // alpha > beta on value
//		lessValuableThan:  [...]    // ... > alpha on value
//		moreUncertainThan: [...]
//		lessUncertainThan: [...]
//		moreComplexThan:   [...]
//		lessComplexThan:   [...]
//	}
//
// "alpha: lessComplexThan: [beta]" and "beta: moreComplexThan: [alpha]" are
// the same statement; stating both is reported as a duplicate by graph.New.
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
#IDs: [...string]
#Record: close({
	supports?:    #IDs
	blocks?:      #IDs
	enables?:     #IDs
	invalidates?: #IDs
	contains?:    #IDs
	relates?:     #IDs

	moreValuableThan?:  #IDs
	lessValuableThan?:  #IDs
	moreUncertainThan?: #IDs
	lessUncertainThan?: #IDs
	moreComplexThan?:   #IDs
	lessComplexThan?:   #IDs
})
#Graph: [string]: #Record
`

// comparisonField maps a comparison field to its dimension. more tells
// whether the record the field belongs to ranks above the listed ids.
type comparisonField struct {
	name string
	dim  graph.Dimension
	more bool
}

var comparisonFields = []comparisonField{
	{"moreValuableThan", graph.Value, true},
	{"lessValuableThan", graph.Value, false},
	{"moreUncertainThan", graph.Uncertainty, true},
	{"lessUncertainThan", graph.Uncertainty, false},
	{"moreComplexThan", graph.Complexity, true},
	{"lessComplexThan", graph.Complexity, false},
}

// Parse reads one CUE graph file. filename is used for positions in errors
// and in the Origin of each statement. Statements are returned in source
// order of records, then in the field order of the notation.
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
	if err := def.Unify(doc).Validate(cue.Concrete(true)); err != nil {
		return nil, nil, cueError(err)
	}

	records, err := doc.Fields()
	if err != nil {
		return nil, nil, cueError(err)
	}
	var comparisons []graph.Comparison
	var relations []graph.Relation
	for records.Next() {
		record := records.Value()
		self, err := knowledge.ParseID(records.Selector().Unquoted())
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %w", record.Pos(), err)
		}
		for _, k := range graph.RelationKinds() {
			err := eachID(record, string(k), func(pos string, other knowledge.ID) {
				relations = append(relations, graph.Relation{Kind: k, From: self, To: other, Origin: pos})
			})
			if err != nil {
				return nil, nil, err
			}
		}
		for _, f := range comparisonFields {
			err := eachID(record, f.name, func(pos string, other knowledge.ID) {
				c := graph.Comparison{Dimension: f.dim, More: self, Less: other, Origin: pos}
				if !f.more {
					c.More, c.Less = other, self
				}
				comparisons = append(comparisons, c)
			})
			if err != nil {
				return nil, nil, err
			}
		}
	}
	return comparisons, relations, nil
}

// eachID calls fn for every id in the list field of record, if present.
// Errors are prefixed with the element position.
func eachID(record cue.Value, field string, fn func(pos string, id knowledge.ID)) error {
	list := record.LookupPath(cue.MakePath(cue.Str(field)))
	if !list.Exists() {
		return nil
	}
	it, err := list.List()
	if err != nil {
		return cueError(err)
	}
	for it.Next() {
		pos := it.Value().Pos().String()
		s, err := it.Value().String()
		if err != nil {
			return fmt.Errorf("%s: %w", pos, err)
		}
		id, err := knowledge.ParseID(s)
		if err != nil {
			return fmt.Errorf("%s: %w", pos, err)
		}
		fn(pos, id)
	}
	return nil
}

// cueError flattens a CUE error list into one error with positions.
func cueError(err error) error {
	return fmt.Errorf("%s", cueerrors.Details(err, nil))
}
