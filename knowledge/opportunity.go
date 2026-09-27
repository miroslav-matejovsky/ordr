package knowledge

import (
	"errors"
	"fmt"
	"regexp"
)

// ID is the stable identifier of a knowledge record.
//
// Invariant: lowercase ASCII letters and digits in segments joined by single
// hyphens, for example "alpha" or "tool-library".
type ID string

var idPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// ParseID validates s and returns it as an ID.
func ParseID(s string) (ID, error) {
	if !idPattern.MatchString(s) {
		return "", fmt.Errorf("invalid id %q: use lowercase letters, digits and single hyphens", s)
	}
	return ID(s), nil
}

// State is the position of an opportunity in the ORDR workflow.
//
// The workflow is: framing -> investigation -> evidence -> options -> decision-ready.
// A state describes understanding reached, not work performed.
type State string

// Workflow states in order.
const (
	StateFraming       State = "framing"
	StateInvestigation State = "investigation"
	StateEvidence      State = "evidence"
	StateOptions       State = "options"
	StateDecisionReady State = "decision-ready"
)

// States returns all workflow states in workflow order.
func States() []State {
	return []State{StateFraming, StateInvestigation, StateEvidence, StateOptions, StateDecisionReady}
}

// ParseState validates s and returns it as a State.
func ParseState(s string) (State, error) {
	for _, st := range States() {
		if string(st) == s {
			return st, nil
		}
	}
	return "", fmt.Errorf("unknown state %q: expected one of %v", s, States())
}

// Opportunity is anything worth exploring toward a decision.
//
// Evidence and Uncertainties keep the author's order. They may be empty;
// Summary, Hypothesis and Decision may not.
type Opportunity struct {
	ID            ID
	Title         string
	Summary       string
	State         State
	Hypothesis    string   // current belief being tested
	Evidence      []string // what is known and why it is believed
	Uncertainties []string // what is still open
	Decision      string   // the decision this opportunity prepares
	Source        string   // where the record was read from, for error messages
}

// NewOpportunity returns a validated Opportunity.
func NewOpportunity(o Opportunity) (Opportunity, error) {
	if _, err := ParseID(string(o.ID)); err != nil {
		return Opportunity{}, err
	}
	if _, err := ParseState(string(o.State)); err != nil {
		return Opportunity{}, err
	}
	required := []struct{ field, value string }{
		{"title", o.Title},
		{"summary", o.Summary},
		{"hypothesis", o.Hypothesis},
		{"decision", o.Decision},
	}
	var errs []error
	for _, r := range required {
		if r.value == "" {
			errs = append(errs, fmt.Errorf("%s must not be empty", r.field))
		}
	}
	if len(errs) > 0 {
		return Opportunity{}, errors.Join(errs...)
	}
	return o, nil
}
