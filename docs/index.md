# metatools-a2a

`metatools-a2a` is the **A2A reference server** for the ApertureStack tool stack.
It exposes tools as A2A skills, supports JSON-RPC + REST endpoints, and streams
task updates over SSE.

## Highlights

- A2A JSON-RPC (`agent/invoke`, `agent/status`)
- Agent card generation from canonical provider metadata
- REST discovery endpoints for skills and tasks
- SSE task updates for long-running executions

## Key Packages

- `internal/agent` — maps ApertureStack discovery/execution into A2A semantics
- `internal/server` — HTTP server and routing
- `toolprotocol/a2a` — protocol binding shared across servers
