# Architecture

The default shape is a modular monolith with four clear runtime boundaries:

- `apps/api` is a Go `net/http` server. Handlers expose `/api/v1` and delegate
  application work to internal packages.
- `apps/worker` is a Go process using Shoebox `v0.1.6` for in-process job
  dispatch backed by PostgreSQL. Jobs are retryable and idempotent.
- PostgreSQL is the durable source of truth. Go code uses `pgx/v5`; the API
  uses a shared `pgxpool.Pool`, and Shoebox uses its PostgreSQL store.
- `apps/web` is a TanStack Start React application. Reusable shadcn-style
  components live in `packages/ui` and are composed with Radix primitives.

The browser talks to the API through `/api/v1`. Keep transport, domain, and
persistence concerns separate, and put slow or retryable work in the worker.
