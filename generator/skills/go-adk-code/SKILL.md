---
name: go-adk-code
description: >
  Go ADK v2.0 API references, interfaces, custom tools, and adapters.
  Provides patterns for model initialization, agent configuration, programmatic execution,
  and custom OpenAI/Copilot LLM adapters.
metadata:
  author: Antigravity
  license: Apache-2.0
  version: 2.0.0
---

# Go ADK v2.0 Code Reference

This skill provides a quick reference for building agents, registering custom tools, and executing agents programmatically using **Go ADK v2.0**.

---

## 1. Core Go ADK Interfaces

Go ADK v2.0 hides concrete struct definitions (e.g. `geminiModel`) behind interface types. Always use the interface types in your helper signatures:

* **Model Brain:** Use `model.LLM` interface type (instead of `gemini.Model`).
* **Agent:** Use `agent.Agent` interface type (instead of `llmagent.Agent`).

---

## 2. Initializing the Gemini Model

Initialize the model client with dynamic API key checks and automatic Vertex AI support:

```go
import (
	"context"
	"os"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/model/gemini"
	"google.golang.org/genai"
)

func InitModel(ctx context.Context, modelName string) (model.LLM, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}
	if apiKey == "" {
		apiKey = os.Getenv("GENAI_API_KEY")
	}

	clientConfig := &genai.ClientConfig{
		APIKey: apiKey,
	}

	// Dynamic Vertex AI fallback config
	if os.Getenv("GOOGLE_GENAI_USE_VERTEXAI") == "1" {
		clientConfig.Backend = genai.BackendVertexAI
		clientConfig.Project = os.Getenv("GOOGLE_CLOUD_PROJECT")
		clientConfig.Location = os.Getenv("GOOGLE_CLOUD_LOCATION")
	}

	return gemini.NewModel(ctx, modelName, clientConfig)
}
```

---

## 3. Instantiating a Simple Agent

Construct the agent using `llmagent.New` and `llmagent.Config`:

```go
import (
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model"
)

func NewSimpleAgent(m model.LLM, name, description, instruction string) (agent.Agent, error) {
	return llmagent.New(llmagent.Config{
		Name:        name,
		Model:       m,
		Description: description,
		Instruction: instruction,
	})
}
```

---

## 4. Programmatic Execution (`runner.Runner`)

Use the `runner.Runner` pattern and Go 1.23+ sequence iterators to execute agents programmatically:

```go
import (
	"context"
	"fmt"
	"log"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/runner"
	"google.golang.org/adk/v2/session"
	"google.golang.org/genai"
)

func ExecuteAgent(ctx context.Context, rootAgent agent.Agent, prompt string) {
	sessionService := session.InMemoryService()
	r, err := runner.New(runner.Config{
		AppName:           "my_app",
		Agent:             rootAgent,
		SessionService:    sessionService,
		AutoCreateSession: true,
	})
	if err != nil {
		log.Fatalf("Failed to create runner: %v", err)
	}

	userMessage := &genai.Content{
		Role: "user",
		Parts: []*genai.Part{
			{Text: prompt},
		},
	}

	events := r.Run(ctx, "user_1", "session_1", userMessage, agent.RunConfig{})
	for event, err := range events {
		if err != nil {
			log.Fatalf("Execution failed: %v", err)
		}
		if event.Content != nil {
			for _, part := range event.Content.Parts {
				if part.Text != "" {
					fmt.Print(part.Text)
				}
			}
		}
	}
	fmt.Println()
}
```

---

## 5. Custom OpenAI/Copilot Adapter Pattern

To use non-Gemini backends (such as GitHub Copilot or GitHub Models), write a custom class implementing `model.LLM`:

```go
type CopilotLLM struct {
	modelName string
	githubPAT string
}

func (c *CopilotLLM) Name() string { return c.modelName }

func (c *CopilotLLM) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		// 1. Translate req.Contents to OpenAI messages (plain text + base64 InlineData images)
		// 2. Translate req.Config.Tools to OpenAI tools array
		// 3. Dispatch POST HTTP request with Authorization Bearer header
		// 4. Parse response tool_calls / text content
		// 5. Yield translated model.LLMResponse
	}
}
```
For the complete implementation, refer to the [go_adk_copilot_advanced.md](file:///home/zhenya/.gemini/antigravity-cli/brain/3193f44d-b408-4c55-b83a-af24771783b2/go_adk_copilot_advanced.md) guide.

---

## 6. Separating Config (go:embed)

Keep long instruction prompts clean and organized in markdown files:

```go
import _ "embed"

//go:embed prompt.md
var systemInstruction string
```
