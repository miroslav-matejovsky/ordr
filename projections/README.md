# Projections

The projections domain transforms knowledge and graph relationships into consumable views.

Responsibilities:

- Rankings
- Matrices
- Review summaries
- Decision boards
- Reports

Examples:

- Strategic Opportunity Matrix
- Value Ranking
- Uncertainty Ranking
- Decision Readiness Board

Questions answered:

- What should be visible right now?
- What decisions require attention?

Implementation (proof of concept):

- `value`, `uncertainty`, `complexity`: layered ordering, items never compared, stated comparisons.
- `relationships`: all relations grouped by kind.
- `readiness`: opportunities grouped by workflow state with decision, hypothesis, evidence,
  open uncertainties and incoming relations.
- Pure functions of a knowledge register and a graph. Deterministic Markdown. No file access.
