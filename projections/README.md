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

- `focus`: the single opportunity to work on now. Starts at the top of the value ordering
  (ties: most advanced state, then id) and follows `blocks` and `enables` to the highest ranked
  constraint until none is left. Shows the path and the relationships of the result as ASCII trees.
- `value`, `uncertainty`, `complexity`: layered ordering, items never compared, stated comparisons.
- `relationships`: all relations grouped by kind.
- `readiness`: opportunities grouped by workflow state with decision, hypothesis, evidence,
  open uncertainties and incoming relations.
- Pure functions of a knowledge register and a graph. Deterministic Markdown. No file access.
