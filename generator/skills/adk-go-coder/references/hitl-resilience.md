# Human-in-the-Loop (HITL) and Resilience in Go ADK v2.0

ADK v2.0 introduces first-class support for Human-in-the-Loop (HITL) execution, along with declarative resilience features.

---

## 1. Human-in-the-Loop (HITL)

Any node can pause workflow execution to ask a human for approval or information. This pause is durable and persists across process restarts.

### Triggering a Pause
To ask for input, emit a `RequestInputEvent`:
```go
event := workflow.NewRequestInputEvent(ctx, session.RequestInput{
	InterruptID:    "approve_payment",
	Message:        "Do you approve the $50 charge? (yes/no)",
	ResponseSchema: schema,
})
```

### Resuming Workflows
The scheduler reconstructs the execution state from the session history once a response is submitted.
*   **Handoff Mode**: The human's response flows directly as the input to the next node.
*   **Re-entry Mode**: The paused node re-runs, and can access the human's response via the context:
```go
if ctx.ResumedInput("approve_payment") != nil {
	response := ctx.ResumedInput("approve_payment").(string)
	// Process approval
}
```

---

## 2. Declarative Resilience

### Node Retry Policy
Configure exponential backoff and jitter directly on nodes without boilerplate:
```go
nodeCfg := workflow.NodeConfig{
	RetryConfig: &workflow.RetryConfig{
		MaxAttempts: 5,
		InitialDelay: 1 * time.Second,
		MaxDelay:     60 * time.Second,
		BackoffFactor: 2.0,
		Jitter:        true,
	},
}
```

### Timeouts and Concurrency limits
*   **Timeout**: Cap node execution duration using `NodeConfig.Timeout`.
*   **Concurrency**: Cap overall graph concurrency:
```go
wf, err := workflowagent.New(workflowagent.Config{
	Name:           "refund_flow",
	Edges:          edges,
	MaxConcurrency: 5,
})
```

### Isolation Scopes
Isolate parallel graph branches so that their chat and model call histories do not cross-contaminate:
```go
// Isolated child node executions
workflow.RunNode[string](nc, stepNode, in, workflow.WithIsolationScope("branch_A"))
```
