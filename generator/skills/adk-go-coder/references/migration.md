# Migrating from Go ADK 1.x to 2.0

ADK Go 2.0 introduces a unified execution runtime that impacts package paths, function signatures, and context management.

---

## 1. Import Path Changes
All ADK package paths must include the `/v2` version suffix:
```diff
-import "google.golang.org/adk/agent"
-import "google.golang.org/adk/server/adkrest"
+import "google.golang.org/adk/v2/agent"
+import "google.golang.org/adk/v2/server/adkrest"
```

---

## 2. Unified Context
In ADK v1, context types were split (e.g. `ToolContext`, `CallbackContext`, `InvocationContext`). ADK v2 replaces them with a single `agent.Context` type.

### Tool Definition Signatures
```diff
-func MyTool(ctx tool.Context, args Args) (Result, error)
+func MyTool(ctx agent.Context, args Args) (Result, error)
```

### Graph Node Functions
Change the first parameter to `agent.Context`:
```diff
-func Process(ctx agent.InvocationContext, in string) (string, error)
+func Process(ctx agent.Context,           in string) (string, error)
```

---

## 3. Session Event Signatures
`session.NewEvent` now requires the active Go `context.Context`:
```diff
-event := session.NewEvent("session_id", "my-event-name")
+event := session.NewEvent(ctx.Context(), "my-event-name")
```

---

## 4. Telemetry and Event Assertions in Tests
Go ADK v2.0 events carry richer metadata, including node execution details (`NodeInfo`, `IsolationScope`, `Routes`, etc.). 
If unit tests assert on exact equality of `session.Event` objects, ensure that either the new metadata fields are cleared out during comparison or the test doubles populate them correctly.
*   Use `agent/context_mock.go`'s `StrictContextMock` for mocking context arguments in tests.
