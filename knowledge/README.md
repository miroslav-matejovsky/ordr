# Knowledge

The knowledge domain manages all information handled by ORDR.

Responsibilities:

- Loading and parsing knowledge records
- Managing schemas and validation
- Representing opportunities, decisions, evidence, reviews and maps
- Providing a consistent domain model to other domains

Questions answered:

- What do we know?
- Why do we believe it?
- What changed?

Implementation (proof of concept):

- `Opportunity` record: id, title, state, summary, hypothesis, evidence, uncertainties, decision.
- `ParseOpportunity` reads the Markdown record format described in the root `README.md`.
- `Register` holds all opportunities of a workspace and rejects duplicate ids.
- Records hold no rankings; ordering lives in `graph/`.
- No filesystem access; `workspace/` supplies file content.
