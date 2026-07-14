# Go ADK v2.0 Graph-Based Workflows

The `google.golang.org/adk/v2/workflow` package provides a graph-based engine for building complex multi-agent orchestrations. A workflow is represented as a graph of **nodes** connected by **edges**.

---

## 1. Node Types

### Function Node
Wraps a standard Go function. Generics automatically infer input and output schemas:
```go
node := workflow.NewFunctionNode("classify", func(ctx agent.Context, in string) (Category, error) {
	// Processing
	return Category{Type: "question"}, nil
}, cfg)
```

### Emitting Function Node
Allows streaming intermediate events or pausing for a human via an `emit` callback:
```go
node := workflow.NewEmittingFunctionNode("process", func(ctx agent.Context, in Job, emit func(*session.Event) error) (Result, error) {
	emit(session.NewEvent(ctx.Context(), "processing"))
	return Result{Done: true}, nil
}, cfg)
```

### Agent & Tool Nodes
Wraps an `agent.Agent` or a `tool.Tool` directly as a graph step:
```go
agentNode := workflow.NewAgentNode("sub_agent", subAgent, cfg)
toolNode  := workflow.NewToolNode("my_tool", myTool, cfg)
```

### Join Node
Fan-in barrier that blocks execution until all predecessor nodes finish, returning a map of their outputs:
```go
joinNode := workflow.NewJoinNode("join_step", func(ctx agent.Context, inputs map[string]any) (Result, error) {
	// Access inputs by predecessor node name
	data := inputs["research_step_a"].(string)
	return Result{MergedData: data}, nil
}, cfg)
```

### Dynamic Node
Allows runtime orchestration of graph execution (e.g. loops, dynamic fan-out) using plain Go code:
```go
dynamicNode := workflow.NewDynamicNode("orchestrator", func(nc agent.Context, in string, emit func(*session.Event) error) (string, error) {
	res1, _ := workflow.RunNode[string](nc, stepOneNode, in)
	if res1 == "retry" {
		return workflow.RunNode[string](nc, stepOneNode, in)
	}
	return res1, nil
}, workflow.NodeConfig{})
```

### State-Bound Node
Pulls selected session-state values straight into a typed Params struct using struct tags:
```go
type StateParams struct {
	UserID string `state:"user_id"`
	Email  string `state:"user_email"`
}

node := workflow.NewFunctionNodeFromState("profile", func(ctx agent.Context, p StateParams) (Result, error) {
	// Access p.UserID and p.Email directly
	return Result{}, nil
}, cfg)
```

---

## 2. Edges and Routing

Edges connect nodes and support routing based on return values. Supported route engines:
*   `StringRoute`, `IntRoute`, `BoolRoute`, `MultiRoute`
*   `Default` (fallback route)

Example of LLM-steered conditional routing:
```go
b := workflow.NewEdgeBuilder()

// Routing depending on the string value returned by the router node
b.AddRoutes(routerNode, map[string]workflow.Node{
	"question":    answerNode,
	"statement":   commentNode,
	"exclamation": reactNode,
})

// Parallel execution (Fan-out)
b.AddFanOut(plannerNode, researchNodeA, researchNodeB)

// Merging parallel execution (Fan-in)
b.AddFanIn(joinNode, researchNodeA, researchNodeB)

edges := b.Build()
```
