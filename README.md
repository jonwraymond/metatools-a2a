# metatools-a2a

Reference **A2A (Agent-to-Agent)** server for the ApertureStack tool stack. It
bridges **tooldiscovery** and **toolexec** to provide a JSON-RPC + REST + SSE
surface for A2A agents.

## Features

- A2A JSON-RPC (`agent/invoke`, `agent/status`)
- REST endpoints for `agent-card`, `skills`, and task metadata
- SSE streaming for task updates
- Bootstraps tools from a YAML file

## Quick Start

```bash
export METATOOLS_A2A_HOST=0.0.0.0
export METATOOLS_A2A_PORT=8091
export METATOOLS_A2A_BASE_PATH=/a2a
export METATOOLS_A2A_TOOLS_FILE=./tools.yaml

go run ./cmd/metatools-a2a
```

## Tool Bootstrap File

```yaml
tools:
  - tool:
      name: echo
      description: Echo a message
      inputSchema:
        type: object
        properties:
          message:
            type: string
        required: [message]
    backend:
      kind: local
      local:
        name: echo
```

The backend definition follows `toolfoundation/model.ToolBackend`.

## Endpoints

- `POST /a2a` JSON-RPC
- `GET /a2a/agent-card`
- `GET /a2a/skills`
- `GET /a2a/tasks`
- `GET /a2a/tasks/{id}`
- `GET /a2a/tasks/{id}/events` (SSE)

## License

See `LICENSE` in the repository root.
