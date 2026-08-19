# ADR 0001: Modular monolith and selected runtime stack

## Decision

Start with one Go API and one Go worker with explicit domain modules. The API
uses the standard library `net/http`; PostgreSQL access uses `pgx/v5`; the
worker uses `github.com/adexaja/shoebox v0.1.6` with its PostgreSQL backend;
and the frontend uses TanStack Start with a shared shadcn-style UI package.

## Rationale

This keeps transaction boundaries simple while allowing slow work to run
asynchronously. `net/http` and pgx keep the Go runtime small and explicit.
Shoebox provides durable queue semantics without requiring another broker.
TanStack Start gives the frontend a typed, file-based routing model while
shadcn-style components remain owned by the repository. Split services or
replace infrastructure only when a concrete operational boundary justifies
it.

## Consequences

Shoebox’s PostgreSQL tables and application tables share a database but must
remain independently migrated. Queue handlers must tolerate retries and
duplicate delivery. The frontend build requires Bun and the TanStack Start
Vite plugin.
