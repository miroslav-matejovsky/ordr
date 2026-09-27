# ORDR Documentation

Project overview, concepts and architecture. Usage lives in the root [`README.md`](../README.md).

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
| [Decisions over Information](../dev/principles/1-decisions-over-information.md) | Information exists to improve decisions. |
| [Knowledge over Activity](../dev/principles/2-knowledge-over-activity.md) | Progress is measured by learning rather than effort. |
| [Comparison over Estimation](../dev/principles/3-comparison-over-estimation.md) | Relative judgment is preferred over false precision. |

## Architecture

Knowledge + Graph -> Projections. See [`dev/architecture.md`](../dev/architecture.md).

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

## Implementation

The current implementation is a proof of concept. Source layout and dependency rules:
[`dev/architecture.md`](../dev/architecture.md).

### Development

```text
task all      # tidy, fmt, vet, deadcode, arch-lint, lint, test
go test ./...
go run . project --workspace examples/basic
```

`examples/basic/projections/` is compared with fresh output by the CLI integration test.
Regenerate it after changing example sources or rendering.

### Assessments

- [Initial dependency assessment](assessments/deps/2026-09-27-initial-dependency-assessment.md)

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
