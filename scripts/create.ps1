$folders = @(
    "knowledge",
    "graph",
    "projections",
    "workspace",
    "cli"
)

foreach ($folder in $folders) {
    New-Item -ItemType Directory -Path $folder -Force | Out-Null
}

@"
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
"@ | Set-Content "knowledge\README.md"

@"
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
"@ | Set-Content "graph\README.md"

@"
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
"@ | Set-Content "projections\README.md"

@"
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
"@ | Set-Content "workspace\README.md"

@"
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
"@ | Set-Content "cli\README.md"

Write-Host "ORDR domain structure created."