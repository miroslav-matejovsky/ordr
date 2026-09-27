# Examples

Sample ORDR workspaces consumed by the `ordr` command.

| Folder | Purpose |
| --- | --- |
| `basic/` | Valid workspace with three opportunities, comparisons on all dimensions, relations and generated projections. |
| `invalid/` | Workspace with a value-ordering cycle and an unknown reference. `ordr project` must fail on it. |

`basic/projections/` is generated. The CLI integration test compares fresh output with it.
Regenerate after changing sources or rendering:

```text
go run ./cmd/ordr project --workspace examples/basic
```
