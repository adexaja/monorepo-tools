# Go + TypeScript Monorepo Template

Reusable starter based on the structure of `Workspace/examapp`:

- Moon for workspace orchestration
- Bun workspaces for TypeScript applications and packages
- Go `net/http` API with `pgx/v5` PostgreSQL access
- PostgreSQL-backed Shoebox `v0.1.6` background worker
- TanStack Start frontend and shadcn-style shared UI components
- PostgreSQL migrations and local Docker Compose infrastructure
- Modular-monolith documentation and CI-ready task conventions

## Create a project

From this directory:

```sh
./create.sh ../my-project "My Project" github.com/example/my-project
cd ../my-project
bun install
moon run :test
```

Run `bun install` in the generated project before invoking Moon. Start local
infrastructure with `docker compose -f infra/docker/docker-compose.yml up -d`
and copy `infra/.env.example` to your environment when running the API or
worker.

Arguments are destination, display name, and Go module path. The destination
must not already exist. The generated project contains no git history, build
artifacts, or dependency directories.

## Generated layout

```text
apps/api       Go net/http API (`/healthz`, `/readyz`, `/api/v1/health`)
apps/worker    Shoebox v0.1.6 worker using PostgreSQL storage
apps/web       TanStack Start React application
packages/ui    Shared shadcn-style React components
packages/contracts  Shared API contract types
packages/config     Shared runtime/config constants
infra           Docker Compose, env example, and migrations
docs            Architecture, API, and ADR documentation
```

Replace the package names, database defaults, and starter UI as the product
takes shape. Keep domain-specific policy in the generated project's
`AGENTS.md` rather than in this generic template.
