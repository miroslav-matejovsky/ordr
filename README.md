# ORDR

**Opportunity Research & Decision Register**

ORDR is a text-first tool for exploring opportunities, reducing uncertainty, and preparing decisions.

You write opportunities as Markdown and compare them pairwise in CUE files. ORDR validates
both and generates Markdown views: relative orderings, relationships and decision readiness.
It never scores anything: `alpha > beta` instead of `alpha = 87`.

Concepts and principles: [`docs/README.md`](docs/README.md). Architecture: [`dev/architecture.md`](dev/architecture.md).

## Install

Requires Go 1.27.1 or newer.

```text
go install github.com/miroslav-matejovsky/ordr@latest
```

## Quick start

```text
ordr project --workspace examples/basic
```

Run it from a clone of this repository. It writes the projections into `examples/basic/projections/`.

## Usage

```text
ordr project [PROJECTION] [--workspace PATH] [--stdout]
```

- `PROJECTION`: `all` (default), `value`, `uncertainty`, `complexity`, `relationships`, `readiness`.
- `--workspace`: a folder inside the workspace. `ordr.cue` is searched upward from it. Default `.`.
- `--stdout`: print projections instead of writing them.
- Exit codes: `0` success, `1` invalid sources or I/O failure, `2` invalid command line.

## Workspace

```text
my-workspace/
├── ordr.cue        # folder configuration, marks the root
├── knowledge/*.md  # opportunity records, source of truth
├── graph/*.cue     # comparisons and relations, source of truth
└── projections/    # generated, safe to delete
```

`ordr.cue` names every folder explicitly. There are no defaults.

```cue
knowledge:   "knowledge"
graph:       "graph"
projections: "projections"
```

## Knowledge format

One opportunity per Markdown file. All front matter keys and all five sections are required.
Evidence and Uncertainties are `- ` lists and may be empty. Text sections are joined into one line.

```markdown
---
id: alpha
title: Neighbourhood tool library
state: options
---

## Summary
A shared library where residents borrow rarely used tools.

## Hypothesis
Residents borrow tools they rarely need rather than buy them.

## Evidence
- A street survey found most respondents rarely use their drill.

## Uncertainties
- Whether residents return tools on time without deposits.

## Decision
Run a one-season pilot, or drop the idea.
```

- `id`: lowercase letters, digits, single hyphens. Unique within the workspace.
- `state`: `framing`, `investigation`, `evidence`, `options` or `decision-ready`.

## Graph notation (CUE)

Every `*.cue` file in the graph folder is read on its own. All fields are optional; unknown fields are errors.

```cue
comparisons: {
	value:       [{more: "alpha", than: "beta"}]  // alpha has more strategic value than beta
	uncertainty: [{more: "gamma", than: "beta"}]  // gamma is more uncertain than beta
	complexity:  [{more: "gamma", than: "alpha"}] // gamma is more complex than alpha
}
relations: {
	supports:    [{from: "beta", to: "alpha"}]
	blocks:      []
	enables:     [{from: "alpha", to: "gamma"}]    // gamma depends on alpha
	invalidates: []
	relates:     [{from: "beta", to: "gamma"}]     // symmetric
}
```

ORDR refuses invalid input and lists every issue with file, line and ids:
unknown references, self-references, duplicates, contradictory comparisons (`a > b` and `b > a`)
and cycles (`a > b > c > a`).

## Projections

| File | Content |
| --- | --- |
| `value.md` | Strategic value ordering. |
| `uncertainty.md` | Uncertainty ordering. |
| `complexity.md` | Complexity ordering. |
| `relationships.md` | Relations grouped by kind. |
| `readiness.md` | Opportunities by workflow state with decision, hypothesis, evidence, open uncertainties and incoming relations. |

Orderings are layered. Level 1 ranks highest. Items on one level are not ordered against each other.
From `alpha > beta` and `alpha > gamma` ORDR derives `1. alpha`, `2. beta; gamma` and never invents `beta > gamma`.
Items without any comparison are listed as not compared.

## Known limitations

- Opportunity is the only record type.
- Markdown parsing is line based: no nested lists, no multi-line list items, text sections become one line.
- A lower level does not prove that a comparison exists between two items; the stated comparisons list is authoritative.
- Relations are not checked for cycles; `blocks` and `enables` cycles are accepted.
- Comparisons carry no rationale.
- Stale projection files are not removed.
