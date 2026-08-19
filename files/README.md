# Template source tree

This directory is the project skeleton used by `../create.sh`. It contains a
Go `net/http` API, a Shoebox `v0.1.6` worker backed by PostgreSQL, a pgx-based
database boundary, and a TanStack Start frontend with shared shadcn-style UI.
Install the Bun workspace dependencies before running Moon:

```sh
bun install
moon run :test
moon run :lint
moon run :build
```

For a new project, use `../create.sh` so placeholders are rendered first.

The generated app expects `DATABASE_URL` for the API readiness check and the
worker queue. Local PostgreSQL is available through
`infra/docker/docker-compose.yml`.
