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