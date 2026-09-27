# Core Architecture

ORDR separates knowledge, relationships, and projections.

```text
Knowledge + Graph -> Projections
```

## Knowledge

The source of truth. Human-readable records that capture understanding.

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

## Graph

The relationship model. Explicit relationships between knowledge.

Examples:

- supports
- blocks
- enables
- invalidates
- relates
- greater-than (`>`), ranks higher than

Graph answers:

- What depends on what?
- What supports what?
- What blocks what?
- What is more important than what?
- What changed because of a decision?

## Projections

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

## Implementation

The current implementation is a proof of concept: one vertical slice that reads knowledge
and graph sources from a workspace, validates them, derives relative orderings and writes
Markdown projections.

### Source layout

| Folder | Responsibility |
| --- | --- |
| `knowledge/` | Opportunity record, Markdown parser, register of records. |
| `graph/` | Comparisons, relations, validation, partial ordering. `graph/cuegraph` reads CUE. |
| `projections/` | Markdown views derived from knowledge and graph. |
| `workspace/` | Discovery, configuration, reading sources, writing projections. |
| `cli/` | Argument parsing and output. `main.go` in the root is the entry point. |
| `examples/` | Sample workspaces. |

### Dependencies

Dependencies point toward the domain:

```text
main -> cli -> workspace, projections -> graph -> knowledge
```

`graph/cuegraph` is the CUE adapter of the graph context. Only `cuegraph` and `workspace`
may use CUE. The rules are enforced by `.go-arch-lint.yml`.
