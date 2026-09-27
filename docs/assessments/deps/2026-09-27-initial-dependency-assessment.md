# Initial Dependency Assessment

- Date: 2026-09-27
- Source revision: `9eb2516` (branch `mirmat-wip`)
- Scope: first-party Go packages and `.go-arch-lint.yml`
- Out of scope: `vendor/`, `examples/`, test code, Taskfile tooling
- Nature: analysis only; no recommendations or change plan

## Executive Summary

The dependency graph is a strict, acyclic, layered structure with `knowledge` as its
only fully stable package. The `.go-arch-lint.yml` rules describe the actual imports
exactly: every allowed edge is used and no used edge is missing. The model is sound
for a proof of concept.

The main findings are:

1. The dependency direction supports the architecture diagram `Knowledge + Graph -> Projections`
   and keeps the domain (`knowledge`, `graph`) free of CUE and the filesystem.
2. `graph` depends on `knowledge` only for the identifier type `knowledge.ID`. It never
   sees knowledge content. The coupling is real but narrow.
3. `workspace` carries more than its stated purpose. Beside discovery, configuration
   and filesystem access, it acts as the application service that assembles and
   validates the domain. It is the least stable non-entry package, yet `cli` depends on it.
4. The vocabulary of the graph context (dimensions, relation kinds) is enumerated in
   three packages: `graph`, `cuegraph` and `projections`. The dependency rules allow this,
   and nothing detects drift between the copies.
5. The linter rules constrain first-party and CUE imports only. The standard library
   is unconstrained, so "domain has no filesystem access" is a convention, not a rule.
6. "Comparison over Estimation" is structurally enforced. "Decisions over Information"
   is partly enforced. "Knowledge over Activity" cannot be observed yet: there is no
   model of change over time.

## Current Dependency Model

### Declared rules

| Component | Package | May depend on | Vendor |
| --- | --- | --- | --- |
| `main` | `.` | `cli` | - |
| `cli` | `cli` | `workspace`, `projections` | - |
| `workspace` | `workspace` | `knowledge`, `graph`, `cuegraph` | `cue` |
| `projections` | `projections` | `knowledge`, `graph` | - |
| `cuegraph` | `graph/cuegraph` | `knowledge`, `graph` | `cue` |
| `graph` | `graph` | `knowledge` | - |
| `knowledge` | `knowledge` | (nothing) | - |

`allow.depOnAnyVendor: false` and `allow.deepScan: true` are set. Test files are excluded.

### Actual imports

```text
main --> cli --> workspace --> cuegraph --> graph --> knowledge
          |          |  +-----------------> graph
          |          +--------------------------------> knowledge
          +----> projections --> graph
                     +------------------------------> knowledge

cue vendor: workspace, cuegraph
```

Every declared edge is used. No undeclared edge exists. The configuration is a precise
snapshot of the code, so any new edge requires an explicit configuration change.

### Coupling metrics

Afferent coupling (Ca, incoming), efferent coupling (Ce, outgoing, first-party only),
instability I = Ce / (Ca + Ce).

| Component | Ca | Ce | I | Abstract types |
| --- | --- | --- | --- | --- |
| `knowledge` | 4 | 0 | 0.00 | none |
| `graph` | 3 | 1 | 0.25 | none |
| `cuegraph` | 1 | 2 | 0.67 | none |
| `projections` | 1 | 2 | 0.67 | none |
| `cli` | 1 | 2 | 0.67 | none |
| `workspace` | 1 | 3 | 0.75 | none |
| `main` | 0 | 1 | 1.00 | none |

All edges point from higher to lower or equal instability except one:
`cli (0.67) -> workspace (0.75)`. The package is concrete throughout. No interfaces
exist, so stability comes from having few dependencies, not from abstraction.

## Principle Alignment Assessment

### Decisions over Information

Reinforced:

- `knowledge.Opportunity` requires a non-empty `Decision`. A record cannot exist without
  a decision context. This is a domain invariant in the most stable package, so every
  consumer inherits it.
- The `readiness` projection is organised around the decision to prepare, not around
  record content volume.

Limited:

- The decision is a string attribute of an opportunity, not a knowledge record in its own
  right. The README lists Decisions, Evidence and Reviews as knowledge. The dependency model
  currently has no place where a decision is a first-class concept that the graph can
  reference.
- Relations such as `supports` and `invalidates` point between opportunities, not from
  evidence to a hypothesis or decision. The graph can say "beta supports alpha" but not
  "this evidence strengthens this option".

### Knowledge over Activity

Reinforced:

- No package depends on time, effort, assignee or task concepts. Nothing in the
  dependency graph could host activity tracking without a new component.
- Workflow states (`framing` ... `decision-ready`) describe understanding reached.

Not yet observable:

- The principle measures progress as uncertainty removed. That is a change between two
  points in time. `workspace.Load` reads one snapshot. No component owns history or
  change detection, so the projections cannot show uncertainty reduced, only uncertainty
  present.
- The Knowledge README question "What changed?" has no owner in the current model.

### Comparison over Estimation

Structurally enforced:

- `graph.Comparison` has no numeric field. `graph.Ordering` exposes levels, not ranks or
  scores.
- The CUE schema in `cuegraph` is closed. A `score` field is a load error, covered by
  a test.
- `knowledge` has no numeric attributes and does not depend on `graph`, so it cannot
  carry a calculated ranking.
- Unranked items are reported separately. The ordering does not invent comparisons.

Subtle point:

- The `readiness` projection sorts opportunities within a state by the value ordering.
  Items on the same value level are then shown in identifier order. A reader could
  perceive this list order as a ranking. The presentation, not the dependency model,
  carries this risk.

## Bounded Context Assessment

| Context | Stated responsibility | Observed responsibility | Boundary quality |
| --- | --- | --- | --- |
| `knowledge` | Records, parsing, schemas, validation | Opportunity type, Markdown parser, register, id grammar, workflow states | Clear. No dependencies. Owns identity for the whole system. |
| `graph` | Relationships, ordering, traversal | Comparisons, relations, validation, partial order | Clear. Depends only on `knowledge.ID`. |
| `cuegraph` | (not stated in README) | Translate CUE text into graph statements | Clear as an adapter. Duplicates graph vocabulary. |
| `projections` | Views derived from knowledge and graph | Rendering, plus interpretation of relation semantics and workflow order | Mostly clear. Holds some graph semantics. |
| `workspace` | Discovery, configuration, filesystem | Discovery, configuration, filesystem, parser selection, domain assembly, graph validation trigger | Broader than stated. |
| `cli` | Arguments, rendering, interaction | Arguments, projection selection, load-render-write sequence | Clear, with a small orchestration role. |

Observations:

- `knowledge` is the shared kernel of the system. Its `ID` type and id grammar are used
  by every other domain-facing package. This concentration is appropriate for identity,
  and it makes `knowledge` changes the most far-reaching.
- `projections` decides what relations mean to a reader: `enables` becomes "Depends on",
  `blocks` becomes "Blocked by" (`projections/projections.go:120`, `:146`). The graph
  defines the relation kinds, but their directional reading is defined downstream.
- `projections` assumes that the graph was built from exactly the register's
  identifiers and panics otherwise (`projections/projections.go:230`, `:238`). This
  pairing is established in `workspace.Load` (`workspace/workspace.go:155`). No type
  expresses it.

## Dependency Direction Analysis

### Inward flow

Dependencies flow from delivery (`main`, `cli`) through assembly (`workspace`) and
presentation (`projections`) toward the domain (`graph`, `knowledge`). The domain
packages import only the standard library. This matches the intended direction.

### Stability

The stable dependencies principle holds for every edge except `cli -> workspace`.
`workspace` has the highest efferent coupling of any library package because it wires
three domain-facing packages together. Changes in `knowledge`, `graph` or `cuegraph`
signatures reach `cli` through `workspace`.

### Dependency inversion

There are no interfaces. The composition is direct:

- `workspace` calls `knowledge.ParseOpportunity`, `cuegraph.Parse` and `graph.New` directly.
- `cli` calls `workspace` and `projections` directly.

Consequences, stated without preference:

- The code path is short, explicit and easy to follow, which fits the project rules
  (no premature interfaces, accept concrete types).
- Load behaviour can only be exercised through a real directory tree and real CUE
  evaluation.
- `cli` cannot distinguish a validation failure from an I/O failure without importing
  `graph` for `*graph.ValidationError`. The current rules forbid `cli -> graph`, and both
  map to exit code `1` (`cli/cli.go:21`). The project rule "distinguish business
  failures from system failures" meets a dependency rule here.

### Cyclic risk

The graph is a DAG, and Go rejects import cycles anyway. The latent risk is conceptual:

- If `knowledge` ever needs to refer to graph concepts (for example a record that lists
  its own relations), the only legal direction is already taken by `graph -> knowledge`.
- If `projections` ever needs workspace metadata (paths, configuration), it would need
  `projections -> workspace`, while `cli` already joins both.

### Infrastructure leaking into the domain

- Source locations enter domain types: `graph.Comparison.Origin`, `graph.Relation.Origin`
  and `knowledge.Opportunity.Source` hold file paths and `file:line:col` strings
  produced by adapters. They are opaque strings, so no import leaks, but the values
  travel into projection output. Projections therefore change when a comment is added
  to a CUE file.
- The linter does not restrict the standard library. `knowledge`, `graph` and
  `projections` could import `os` or `path/filepath` without a lint failure.

### Domain leaking into infrastructure

- `workspace` knows that `*.md` files are opportunities and `*.cue` files are graph
  statements (`workspace/workspace.go:125`, `:143`). The mapping from file type to
  record type lives in the filesystem layer. Each new record type touches `workspace`.
- `workspace` decides the assembly order (register first, then graph validated against
  register ids). That ordering is a domain rule held in infrastructure.

## Knowledge/Graph Analysis

### What the graph knows today

`graph` imports `knowledge` for one type: `knowledge.ID`. It receives known identifiers
as `[]knowledge.ID` in `graph.New` (`graph/graph.go:85`), not as a `knowledge.Register`.
It never reads titles, states, evidence or any other content. The graph knows identity,
not knowledge.

### Graph knows knowledge (current)

- Benefit: one identifier type and one identifier grammar. A graph statement cannot
  hold an id that is not a valid knowledge id; `cuegraph` parses ids with
  `knowledge.ParseID`.
- Benefit: reference validation (unknown ids) is expressed in domain terms.
- Cost: `graph` is only as stable as `knowledge.ID`. Any change to identity (namespaces,
  record types in ids, cross-workspace references) is a graph change as well.
- Cost: the graph can only relate things that `knowledge` defines ids for. Relating
  a comparison, a relation or an external source is outside its reach.

### Knowledge knows graph

- Would allow records to carry or validate their own relationships.
- Would place derived structure next to authoritative content, which conflicts with
  "Knowledge must not contain calculated rankings" and with the README separation
  `Knowledge + Graph -> Projections`.
- Would make `knowledge` unstable and invert the current direction.

### Both independent

- Would require identity to live outside both, or be a plain string in the graph.
- Graph validation of references would then need a set of identifiers from somewhere,
  which moves the join to the assembler (today `workspace`).
- Makes the graph reusable for any kind of node, at the cost of weaker typing and
  a join that must be kept correct by convention.

### Assessment

The two concepts are independent in content and dependent in identity. The README
describes them as separate sources of truth joined by projections. The code matches
that description, with identity as the single shared element. Whether identity belongs
to knowledge or is a concept of its own is the open design question in this area.

## CUE Integration Analysis

### Classification

`cuegraph` is an implementation adapter. It translates text into `graph.Comparison` and
`graph.Relation`, performs no validation of graph semantics and holds no state. It does
not form its own bounded context.

However, its CUE schema (`graph/cuegraph/cuegraph.go:38`) is the notation that users
write by hand. For users, the schema is the published language of the graph context.
It is a contract with every existing workspace, independent of Go types.

### CUE use outside the graph

`workspace` also imports CUE, for `ordr.cue` configuration (`workspace/workspace.go:39`).
CUE is therefore a workspace-wide notation, not a graph-only one. The single `cue` vendor
entry does not separate these two uses. The notation choice for configuration and for
graph statements can no longer change independently without touching both components.

### Vocabulary duplication

The set of dimensions and relation kinds is enumerated in:

| Place | Dimensions | Relation kinds |
| --- | --- | --- |
| `graph` constants and validation | yes (`graph/validate.go:105`) | yes |
| `cuegraph` schema text | yes | yes |
| `cuegraph` iteration | yes (`graph/cuegraph/cuegraph.go:57`) | uses `graph.RelationKinds()` |
| `projections` | yes (`projections/projections.go:84`) | yes (`:120`, `:146`) |

Adding a dimension touches three packages. A dimension known to `graph` but missing
from the CUE schema is silently unreachable from source files. A dimension missing from
`projections` is a runtime panic on lookup, not a compile error.

### Maintainability implications

- The CUE API surface used is small (`cuecontext`, `cue.Value`, `cue/errors`), which limits
  upgrade risk to two packages.
- Source positions depend on reading the original document rather than the value
  unified with the schema. That is an API behaviour, not a documented contract, and is
  covered only by test expectations of exact positions.
- The domain would survive replacing CUE. Workspaces written in CUE would not.

## Workspace Analysis

Stated purpose: discovery, configuration loading, filesystem interaction.

| Responsibility | Present | Matches purpose |
| --- | --- | --- |
| Find root by walking up to `ordr.cue` | yes | yes |
| Load and validate configuration | yes | yes |
| List and read source files | yes | yes |
| Write projection files | yes | yes |
| Choose the parser per folder | yes | partly; this is format knowledge |
| Build the register from parsed records | yes | no; domain assembly |
| Validate the graph against the register | yes | no; domain assembly |
| Return domain aggregates to the caller | yes | no; application service |

`workspace.Load` (`workspace/workspace.go:123`) is effectively the application service of
the system. This explains its dependency fan-out and its position as the least stable
library package.

Notable decoupling: `WriteProjection` accepts a name and a string. `workspace` does not
import `projections`. The output side of the workspace is already format-agnostic.

Notable coupling: the only way to obtain a validated register and graph is through the
filesystem. Any other source (a single in-memory document, standard input, a remote
store) would need a second assembler or a change to `workspace`.

## Strengths

- **Acyclic, layered graph** with a single, obvious entry point.
- **Domain free of infrastructure imports**: `knowledge` and `graph` use only the
  standard library. Both are fully unit-testable without files.
- **Most stable package holds identity and records.** The sources of truth sit at the
  bottom of the graph, as the README intends.
- **Notation isolated behind an adapter.** No CUE type appears in any domain signature.
- **Projections are pure.** They depend on domain types only, never on `workspace`
  or CUE, and can be rendered from any source of a register and a graph.
- **Output decoupled from rendering.** `workspace` writes strings, not projection types.
- **Configuration is exact.** The lint rules mirror the imports one to one.
  `depOnAnyVendor: false` makes every external library an explicit decision.
- **Comparison over Estimation is encoded in types and schema**, not only in documentation.

## Risks

### Over-coupling

- `workspace` couples filesystem access to domain assembly and to every record format.
- `cli -> workspace` points toward a less stable package.

### Hidden coupling

- The register and graph must be built from the same identifier set. Only
  `workspace.Load` guarantees it; `projections` panics if it is violated.
- Graph vocabulary is enumerated in three packages with no shared source.
- Relation semantics for readers live in `projections`, not in `graph`.
- Source locations flow from adapters through domain types into rendered output.

### Future evolution constraints

- New record types (Decision, Evidence, Review) must fit `knowledge.ID` and will change
  `workspace` parser selection.
- Change over time ("What changed?", uncertainty reduced) has no owning component.
- Cross-workspace references would change identity, which reaches every package.
- CUE is used by two components for two purposes, so replacing or versioning the notation
  is a two-front change.

### Maintainability concerns

- No rule prevents standard library infrastructure imports in the domain.
- Test code is excluded from the linter, so tests may cross any boundary unnoticed.
- The CUE source-position behaviour is relied upon implicitly.

### Testing concerns

- `workspace.Load` can only be tested with a real directory tree and CUE evaluation.
- The integration test compares against committed projections that contain file line
  numbers, so unrelated source edits break it until projections are regenerated.
- `cli` tests cannot assert business versus system failure, because both share exit
  code `1`.

## Open Questions

1. Is identity a knowledge concept, or a concept shared by knowledge and graph?
2. Is `workspace` meant to be the application service, or should that role be named?
3. Is the CUE notation part of the graph context's published contract, and how is it
   versioned relative to the Go model?
4. Should CUE for configuration and CUE for graph statements be treated as one decision
   or two?
5. Who owns the meaning of relation kinds for readers: `graph` or `projections`?
6. Should source locations be part of domain statements, or metadata attached by adapters?
7. Where will change over time live, given that uncertainty reduction is the primary
   measure of progress?
8. Should the dependency rules also cover standard library infrastructure packages for
   domain components?
9. Should the dependency rules apply to test code?
10. The repository documents assessments under `dev/assessments/<name>/README.md`.
    This document is in `docs/assessments/deps/` as requested. Which location is canonical?

## Conclusion

The dependency model is small, acyclic and directed toward the domain. The two sources
of truth, `knowledge` and `graph`, are the most stable packages, free of infrastructure
and notation concerns, and joined only by identity. The CUE adapter is well contained
as far as Go types are concerned, while its notation is a user-facing contract.

The tension sits in the middle layer. `workspace` is both the filesystem boundary and the
system's assembler, which widens its dependencies beyond its stated purpose and makes
`cli` depend on the least stable library package. Smaller forms of hidden coupling
exist around shared vocabulary, the register and graph pairing, and source locations
in domain values.

On principles, Comparison over Estimation is enforced by structure, Decisions over
Information by a required field, and Knowledge over Activity is not yet expressible
because no component models change over time.
