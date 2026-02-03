# Usage

## Run the Server

```bash
export METATOOLS_A2A_HOST=0.0.0.0
export METATOOLS_A2A_PORT=8091
export METATOOLS_A2A_BASE_PATH=/a2a
export METATOOLS_A2A_TOOLS_FILE=./tools.yaml

go run ./cmd/metatools-a2a
```

## Agent Card

```bash
curl http://localhost:8091/a2a/agent-card
```

## Invoke a Skill (JSON-RPC)

```bash
curl -X POST http://localhost:8091/a2a \\
  -H "Content-Type: application/json" \\
  -d '{
    "jsonrpc": "2.0",
    "id": "task-1",
    "method": "agent/invoke",
    "params": {
      "skillId": "example:echo:1.0.0",
      "arguments": {"message": "hello"}
    }
  }'
```

## Task Events (SSE)

```bash
curl http://localhost:8091/a2a/tasks/task-1/events
```
