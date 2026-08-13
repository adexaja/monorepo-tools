# Go + TypeScript Monorepo Template

Reusable starter based on the structure of `Workspace/examapp`:

- Moon for workspace orchestration
- Bun workspaces for TypeScript applications and packages
- Go HTTP API and background worker
- PostgreSQL, Redis, migrations, and sqlc configuration
- Modular-monolith documentation and CI-ready task conventions

## Create a project

From this directory:

```sh
./create.sh ../my-project "My Project" github.com/example/my-project
cd ../my-project
bun install
moon run :test
```

Arguments are destination, display name, and Go module path. The destination
must not already exist. The generated project contains no git history, build
artifacts, or dependency directories.

## Generated layout

```text
apps/api       Go HTTP API (`/healthz`, `/api/v1/health`)
apps/worker    Go background worker entrypoint
apps/web       TypeScript web application placeholder
packages/ui    Shared UI package placeholder
packages/contracts  Shared API contract types
packages/config     Shared runtime/config constants
infra           Docker Compose, env example, migrations, sqlc config
docs            Architecture, API, and ADR documentation
```

Replace the package names, database defaults, and placeholder application code
as the product takes shape. Keep domain-specific policy in the generated
project's `AGENTS.md` rather than in this generic template.
