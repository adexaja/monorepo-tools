# Project instructions

Project: __PROJECT_NAME__.

Use a modular monolith by default. Keep domain modules explicit and share
cross-cutting infrastructure through the Go API/worker packages.

## Boundaries

- The API owns `/api/v1` and is the source of truth for API contracts.
- PostgreSQL is the source of truth for durable business state.
- Redis is for cache, rate limiting, and short-lived coordination.
- Background jobs must be retryable and idempotent; messages should contain
  resource IDs rather than large domain payloads.
- Keep HTTP handlers thin and put business logic in application/domain code.

Before completing work, run relevant tests, lint/typechecks, and migrations
checks. Update documentation and tests when behavior changes.
