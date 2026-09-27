# CLI

The command-line interface domain.

Responsibilities:

- Command execution
- Argument parsing
- Terminal rendering
- User interaction

The CLI should contain no business logic.

Questions answered:

- How does a user interact with ORDR?
- How are projections presented?

Implementation (proof of concept):

- `ordr project [PROJECTION] [--workspace PATH] [--stdout]`
- Exit codes: `0` success, `1` invalid sources or I/O failure, `2` invalid command line.
