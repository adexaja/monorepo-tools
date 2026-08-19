# Worker internals

Place job definitions, consumers, and retry/idempotency logic here. Workers
must not become the source of truth for durable business state.

The default queue is Shoebox backed by PostgreSQL. Register handlers with
explicit retry limits and timeouts, pass compact resource-oriented payloads,
and make every handler safe to run more than once. Use `DATABASE_URL` for the
Shoebox store; the library uses pgx under the hood for its PostgreSQL backend.
