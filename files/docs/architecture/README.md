# Architecture

The default shape is a modular monolith: one Go API and one Go worker with
explicit domain modules. PostgreSQL owns durable state. Redis is limited to
cache, rate limiting, and short-lived coordination. Slow or retryable work
belongs in the worker.
