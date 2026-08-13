# ADR 0001: Modular monolith

## Decision

Start with one Go API and one Go worker with explicit domain modules.

## Rationale

This keeps transaction boundaries simple while allowing slow work to run
asynchronously. Split services only when a concrete operational boundary
justifies it.
