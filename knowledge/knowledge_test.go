package knowledge_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/miroslav-matejovsky/ordr/knowledge"
)

const validRecord = `---
id: alpha
title: Tool lending library
state: options
---

## Summary

Lend tools
to neighbours.

## Hypothesis

People borrow rather than buy.

## Evidence

- Survey shows interest.
- Similar library exists elsewhere.

## Uncertainties

## Decision

Open a pilot or drop the idea.
`

func TestParseOpportunity(t *testing.T) {
	o, err := knowledge.ParseOpportunity("knowledge/alpha.md", []byte(strings.ReplaceAll(validRecord, "\n", "\r\n")))
	require.NoError(t, err)
	require.Equal(t, knowledge.Opportunity{
		ID:            "alpha",
		Title:         "Tool lending library",
		Summary:       "Lend tools to neighbours.",
		State:         knowledge.StateOptions,
		Hypothesis:    "People borrow rather than buy.",
		Evidence:      []string{"Survey shows interest.", "Similar library exists elsewhere."},
		Uncertainties: []string{},
		Decision:      "Open a pilot or drop the idea.",
		Source:        "knowledge/alpha.md",
	}, o)
}

func TestParseOpportunityErrors(t *testing.T) {
	tests := []struct {
		name    string
		old     string
		new     string
		wantErr string
	}{
		{"no front matter", "---\nid: alpha", "id: alpha", "front matter must start"},
		{"missing key", "state: options\n", "", `missing key "state"`},
		{"unknown key", "state: options", "state: options\nscore: 8", `unknown key "score"`},
		{"invalid id", "id: alpha", "id: Alpha One", `invalid id "Alpha One"`},
		{"unknown state", "state: options", "state: done", `unknown state "done"`},
		{"unknown section", "## Hypothesis", "## Priority", `unknown section "Priority"`},
		{"missing section", "## Decision\n\nOpen a pilot or drop the idea.\n", "", `missing section "Decision"`},
		{"duplicate section", "## Uncertainties", "## Evidence", `duplicate section "Evidence"`},
		{"empty summary", "Lend tools\nto neighbours.", "", "summary must not be empty"},
		{"non list evidence", "- Survey shows interest.", "Survey shows interest.", `expected "- " list item`},
		{"text before sections", "---\n\n## Summary", "---\nintro\n## Summary", "before first section"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := strings.Replace(validRecord, tt.old, tt.new, 1)
			require.NotEqual(t, validRecord, src, "test case did not change the record")
			_, err := knowledge.ParseOpportunity("knowledge/alpha.md", []byte(src))
			require.ErrorContains(t, err, "knowledge/alpha.md: ")
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func TestNewRegister(t *testing.T) {
	a := knowledge.Opportunity{ID: "beta", Source: "b.md"}
	b := knowledge.Opportunity{ID: "alpha", Source: "a.md"}

	r, err := knowledge.NewRegister([]knowledge.Opportunity{a, b})
	require.NoError(t, err)
	require.Equal(t, []knowledge.ID{"alpha", "beta"}, r.IDs())
	got, ok := r.Opportunity("beta")
	require.True(t, ok)
	require.Equal(t, a, got)

	_, err = knowledge.NewRegister([]knowledge.Opportunity{a, {ID: "beta", Source: "c.md"}})
	require.EqualError(t, err, `duplicate id "beta" in b.md and c.md`)
}
