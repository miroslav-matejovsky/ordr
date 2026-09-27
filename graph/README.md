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
- A > B

Questions answered:

- What is connected?
- What depends on what?
- What is more important than what?
- What should be re-evaluated after a change?

Implementation (proof of concept):

- `graph.New` validates comparisons and relations: unknown references, self-references,
  duplicates, contradictory direct comparisons and ordering cycles. All issues are reported together.
- `Graph.Order` derives a layered partial order per dimension (value, uncertainty, complexity).
  Items without a stated or implied comparison are never ordered against each other.
- Relation kinds: supports, blocks, enables, invalidates, relates.
- `graph/cuegraph` reads the CUE notation. The domain model does not depend on CUE.
