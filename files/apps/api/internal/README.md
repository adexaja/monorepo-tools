# API internals

Place domain modules, application services, persistence queries, and platform
adapters here. Keep transport code in `cmd/api` thin.

The API uses `pgxpool.Pool` for PostgreSQL connections. Construct one pool at
startup from `DATABASE_URL`, pass it to application services, and close it on
shutdown. Keep SQL close to the repository that owns the query and use
context-aware calls so requests can be cancelled.

HTTP cross-cutting behavior lives in `internal/httpx`: request IDs are
propagated through context and `X-Request-ID`, requests are emitted as
structured logs, and `Idempotency-Key` can replay completed mutating requests.
The default idempotency store is intentionally in-memory for the starter;
provide a PostgreSQL-backed `httpx.IdempotencyStore` for production.
