# Architecture

`metatools-a2a` layers an A2A protocol surface on top of the ApertureStack
discovery + execution stack.

```
Client (A2A) → toolprotocol/a2a → metatools-a2a/agent → tooldiscovery + toolexec
```

## Execution Flow

1. JSON-RPC request arrives (`agent/invoke`)
2. Task created and streamed via SSE
3. Tool executed via `toolexec/run`
4. Task completed with result payload

## Discovery Flow

1. `skills` list is derived from `tooldiscovery` summaries
2. Agent card is generated from canonical provider metadata
