# API

The API contract is rooted at `/api/v1` and is owned by the backend. The
transport is Go’s standard `net/http` package; `/healthz` is liveness and
`/readyz` checks PostgreSQL when `DATABASE_URL` is configured. Generate
TypeScript client types from the backend contract when endpoint modules exist.

Every response includes an `X-Request-ID` header. Clients may send a request
ID for tracing; empty or oversized values are replaced with a generated ID.
The API emits one structured request log with method, path, status, duration,
and request ID.

Mutating requests may send an `Idempotency-Key`. The starter middleware
replays completed results for 24 hours in its in-memory store. Replace that
store with a PostgreSQL implementation before running multiple API instances
or requiring restart-safe idempotency.
