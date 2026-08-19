# Project instructions

Project: __PROJECT_NAME__.

Use a modular monolith by default. Keep domain modules explicit and share
cross-cutting infrastructure through the Go API/worker packages.

## Boundaries

- The API owns `/api/v1` and is the source of truth for API contracts.
- PostgreSQL is the source of truth for durable business state and is accessed
  from Go with `github.com/jackc/pgx/v5`.
- The worker uses `github.com/adexaja/shoebox v0.1.6` with PostgreSQL storage;
  no external queue broker is required by the template.
- Background jobs must be retryable and idempotent; messages should contain
  resource IDs rather than large domain payloads.
- Keep HTTP handlers thin and put business logic in application/domain code.
- Every API request should retain its `X-Request-ID` correlation value in
  logs. Use `Idempotency-Key` for replay-safe mutating operations; it is
  distinct from the request ID and must be backed by durable storage when the
  API is deployed across instances.
- The web app uses TanStack Start and shared shadcn-style components from
  `packages/ui`; add components with the shadcn conventions and keep them
  composable.

Before completing work, run relevant tests, lint/typechecks, and migrations
checks. Update documentation and tests when behavior changes.
