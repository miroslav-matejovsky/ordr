package knowledge

import (
	"fmt"
	"slices"
)

// Register is the validated set of opportunities of one workspace.
//
// Invariant: identifiers are unique. Iteration order is by ID.
type Register struct {
	byID map[ID]Opportunity
	ids  []ID
}

// NewRegister returns a Register or an error naming both sources of a
// duplicate identifier.
func NewRegister(opportunities []Opportunity) (Register, error) {
	r := Register{byID: make(map[ID]Opportunity, len(opportunities))}
	for _, o := range opportunities {
		if prev, dup := r.byID[o.ID]; dup {
			return Register{}, fmt.Errorf("duplicate id %q in %s and %s", o.ID, prev.Source, o.Source)
		}
		r.byID[o.ID] = o
		r.ids = append(r.ids, o.ID)
	}
	slices.Sort(r.ids)
	return r, nil
}

// IDs returns all identifiers sorted.
func (r Register) IDs() []ID {
	return slices.Clone(r.ids)
}

// Opportunity returns the opportunity with the given id.
func (r Register) Opportunity(id ID) (Opportunity, bool) {
	o, ok := r.byID[id]
	return o, ok
}
