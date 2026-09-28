# Examples

Sample ORDR workspaces consumed by the `ordr` command.

| Folder | Purpose |
| --- | --- |
| `basic/` | Valid workspace with three opportunities, comparisons on all dimensions, relations and generated projections. |
| `invoicing/` | Realistic workspace: a small invoicing SaaS planning its next quarter. Ten opportunities in all workflow states, every relation kind including a `contains` hierarchy, and the graph split into several files by product area. |
| `invalid/` | Workspace with a value-ordering cycle and an unknown reference. `ordr project` must fail on it. |

`basic/projections/` and `invoicing/projections/` are generated. The CLI integration test compares fresh output with them.
Regenerate after changing sources or rendering:

```text
go run . project --workspace examples/basic
go run . project --workspace examples/invoicing
```
