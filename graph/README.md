# Graph

The graph domain manages relationships between knowledge.

Responsibilities:

- Dependency management
- Relative ordering
- Relationship modeling
- Graph traversal and inference

Examples:

- A supports B
- A blocks B
- A enables B
- A contains B (B is part of A)
- A > B

Questions answered:

- What is connected?
- What depends on what?
- What is more important than what?
- What should be re-evaluated after a change?

Implementation (proof of concept):

- `graph.New` validates comparisons and relations: unknown references, self-references,
  duplicates, contradictory direct comparisons, ordering cycles, records with more than one parent
  and containment cycles. All issues are reported together.
- `Graph.Order` derives a layered partial order per dimension (value, uncertainty, complexity).
  Items without a stated or implied comparison are never ordered against each other.
- Relation kinds: supports, blocks, enables, invalidates, contains, relates.
- `contains` forms a forest: each record has at most one parent.
- `graph/cuegraph` reads the CUE notation. The domain model does not depend on CUE.
