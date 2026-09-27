# ORDR

**Opportunity Research & Decision Register**

ORDR is a text-first system for exploring opportunities, reducing uncertainty, and preparing decisions.

It helps individuals organize strategic thinking around ideas, technologies, initiatives, research topics, and emerging opportunities without turning them into projects prematurely.

## Philosophy

ORDR is not a task manager.

ORDR is not a roadmap tool.

ORDR is not a project portfolio system.

ORDR exists to answer a simple question:

> What do I need to learn next to make a better decision?

Progress is measured by uncertainty removed, assumptions validated, and decisions clarified.

## Principles

ORDR is built around three principles:

| Principle | Description |
|-----------|-------------|
| [Decisions over Information](dev/principles/1-decisions-over-information.md) | Information exists to improve decisions. |
| [Knowledge over Activity](dev/principles/2-knowledge-over-activity.md) | Progress is measured by learning rather than effort. |
| [Comparison over Estimation](dev/principles/3-comparison-over-estimation.md) | Relative judgment is preferred over false precision. |

## Core Architecture

ORDR separates knowledge, relationships, and projections.

```text
Knowledge + Graph ──> Projections
```

### Knowledge

Human-readable records that capture understanding.

Examples:

- Opportunities
- Decisions
- Evidence
- Reviews
- Notes
- Wardley Maps

Knowledge answers:

- What do we know?
- Why do we believe it?
- What changed?
- What decisions were made?

### Graph

Explicit relationships between knowledge.

Examples:

- supports
- blocks
- enables
- invalidates
- relates
- greater-than (`>`)

Graph answers:

- What depends on what?
- What supports what?
- What blocks what?
- What is more important than what?
- What changed because of a decision?

### Projections

Generated views derived from knowledge and graph relationships.

Examples:

- Opportunity Rankings
- Strategic Opportunity Matrix
- Decision Readiness Board
- Review Summaries
- Dependency Views
- Wardley Projections

Projections are temporary.

Knowledge and relationships are the source of truth.

## Core Workflow

```text
Opportunity -> Framing -> Investigation -> Evidence -> Options -> Decision Ready
```

## Key Concepts

### Opportunities

Anything worth exploring:

- Business ideas
- Product initiatives
- Technology trends
- Market signals
- Research topics
- Strategic investments

### Relative Ordering

ORDR prefers comparison over scoring.

Instead of:

```text
Playback = 87
AI Assistant = 72
Platform = 91
```

ORDR stores:

```text
Platform > Playback
Playback > AI Assistant
```

When information is incomplete, relative comparison is often more reliable than estimation.

### Uncertainty Reduction

The primary measure of progress.

```text
Before:
- Unknown feasibility
- Unknown constraints
- Unknown cost

After:
- Constraints identified
- Options narrowed
- Key assumptions validated
```

### Decision Readiness

The goal is not certainty.

The goal is to reach a state where a decision can be made with acceptable confidence.

## Proof of Concept

One vertical slice: read knowledge and graph sources from a workspace, validate them, derive relative orderings and write Markdown projections.

### Build and run

```text
go build -o bin/ordr .
go install github.com/miroslav-matejovsky/ordr@latest
go run . project --workspace examples/basic
go run . project value --workspace examples/basic --stdout
go test ./...
```

`ordr project [PROJECTION] [--workspace PATH] [--stdout]`

- `PROJECTION`: `all` (default), `value`, `uncertainty`, `complexity`, `relationships`, `readiness`.
- `--workspace`: a folder inside the workspace. `ordr.cue` is searched upward from it.
- `--stdout`: print projections instead of writing them.
- Exit codes: `0` success, `1` invalid sources or I/O failure, `2` invalid command line.

### Source layout

| Folder | Responsibility |
| --- | --- |
| `knowledge/` | Opportunity record, Markdown parser, register of records. |
| `graph/` | Comparisons, relations, validation, partial ordering. `graph/cuegraph` reads CUE. |
| `projections/` | Markdown views derived from knowledge and graph. |
| `workspace/` | Discovery, configuration, reading sources, writing projections. |
| `cli/` | Argument parsing and output. `main.go` in the root is the entry point. |
| `examples/` | Sample workspaces. |

Dependencies point toward the domain: `cli -> workspace, projections -> graph -> knowledge`.

### Workspace

```text
workspace/
├── ordr.cue        # explicit folder configuration, marks the root
├── knowledge/*.md  # opportunity records, source of truth
├── graph/*.cue     # comparisons and relations, source of truth
└── projections/    # generated, disposable
```

`ordr.cue`:

```cue
knowledge:   "knowledge"
graph:       "graph"
projections: "projections"
```

### Knowledge format

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

### Graph notation (CUE)

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

Loading fails with every issue listed, each with file, line and ids:
unknown references, self-references, duplicates, contradictory direct comparisons (`a > b` and `b > a`)
and cycles (`a > b > c > a`).

### Projections

| File | Content |
| --- | --- |
| `value.md` | Strategic value ordering. |
| `uncertainty.md` | Uncertainty ordering. |
| `complexity.md` | Complexity ordering. |
| `relationships.md` | Relations grouped by kind. |
| `readiness.md` | Opportunities by workflow state with decision, hypothesis, evidence, open uncertainties and incoming relations. |

Orderings are layered partial orders. Level 1 ranks highest. Items on one level are not ordered against each other.
From `alpha > beta` and `alpha > gamma` ORDR derives `1. alpha`, `2. beta; gamma` and never invents `beta > gamma`.
Items without any comparison are listed as not compared.

### Known limitations

- Opportunity is the only record type.
- Markdown parsing is line based: no nested lists, no multi-line list items, text sections become one line.
- A lower level does not prove that a comparison exists between two items; the stated comparisons list is authoritative.
- Relations are not checked for cycles; `blocks` and `enables` cycles are accepted.
- Comparisons carry no rationale.
- Stale projection files are not removed.

## Inspirations

- Wardley Mapping
- WSJF
- Decision Journaling
- Systems Thinking
- Evidence-Based Management
- Architecture Decision Records (ADR)

## Non-Goals

ORDR is not intended for:

- Project management
- Sprint planning
- Resource allocation
- Work item tracking
- Predictive estimation
- Reporting activity

## Mission

Bring order to opportunities before they become commitments.

Transform uncertainty into decision readiness.