# Workspace

The workspace domain manages a repository containing ORDR data.

Responsibilities:

- Workspace discovery
- Configuration loading
- Filesystem interaction
- Repository structure management

Questions answered:

- Where is the active workspace?
- Which files belong to the workspace?
- How is the workspace configured?

Implementation (proof of concept):

- `Discover` walks up from a folder to the first `ordr.yaml` and validates it.
- `ordr.yaml` names the `knowledge`, `graph` and `projections` folders explicitly. No defaults.
- `Load` reads `knowledge/*.md` and `graph/*.cue` in lexical order and returns a validated register and graph.
- `WriteProjection` writes generated files into the projections folder. Sources are never written.
