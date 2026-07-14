---
name: go-adk-workflow
description: >
  Go ADK v2.0 Developer Lifecycle Workflow. Covers project scoping, scaffolding,
  Go compilation, unit testing, and systematic Go debugging. Use this as the entrypoint
  when building Go agents.
metadata:
  author: Antigravity
  license: Apache-2.0
  version: 2.0.0
---

# Go ADK v2.0 Developer Lifecycle Workflow

This skill defines the complete lifecycle workflow for building, testing, and verifying AI agents in Go using the **Agent Development Kit (ADK) v2.0**.

---

## Session Phases

### Phase 0: Scoping & Scenarios
Before writing code, clarify the agent's requirements:
1. **Core Task:** What problem does the agent solve?
2. **Backends:** Will it use Google AI Studio (Gemini Developer API) or Vertex AI?
3. **Tools & Integrations:** What external APIs or local databases will the agent query?
4. **Environment Constraints:** Port preferences, credentials setup.

### Phase 1: Study Reference Samples
Study existing Go ADK v2.0 layouts before building:
* Inspect `/home/zhenya/projects/ai/ADK_examples/adk-go/examples/` for reference structures (workflows, basic tool registration, sessions).

### Phase 2: Scaffold
Because the python-based `agents-cli scaffold` is not suited for Go, use the dedicated Go scaffolding script:
```bash
/home/zhenya/scripts/scaffold-agent.sh <project_name>
```
This generates the Go structure:
```
project/
├── agent/
│   ├── agent.go      # Agent configuration (llmagent.Config)
│   └── agent_test.go # Table-driven unit tests
├── main.go           # Entrypoint (Godotenv, gemini model, launcher)
├── .env & .env.example
├── .gitignore
├── Taskfile.yml
└── README.md
```

### Phase 3: Build & Implement
* Write your agent logic inside `agent/agent.go`.
* Ensure compile-time checks by building frequently:
  ```bash
  go build -o agent-bin main.go
  ```
* For custom Go ADK API patterns, interfaces, and custom adapters, load the `/go-adk-code` skill.

### Phase 4: Test & Verify
Go testing relies on the standard `testing` library (no external assertions unless requested).
* **Verify code correctness:** Use table-driven unit tests checking behavior and boundary inputs (empty/nil values).
* Run tests with the race detector enabled:
  ```bash
  go test -race -v ./...
  ```
* **Smoke testing:** Recompile the binary and run the CLI console launcher mode:
  ```bash
  go build -o agent-bin main.go
  ./agent-bin console
  ```

### Phase 5: Deploy & Observe
* Containerize the Go agent using a multi-stage distroless build.
* Use `/go-adk-deploy` for detailed instructions on Docker packaging, Cloud Run deployments, and port bindings.

---

## Systematic Debugging in Go

When the agent fails to compile or respond:

1. **Verify Compilation:**
   Read the Go compiler output carefully. Pay attention to interface mismatch errors (e.g. concrete types vs `model.LLM` or `agent.Agent`).
2. **Print Stack & Errors:**
   Always log returned errors in handlers using format verbs like `%w` to preserve the underlying error context.
3. **Environment Caching Issues:**
   If changes to `.env` seem ignored, remember that `godotenv` does not override active shell variables. Instruct the user to run:
   ```bash
   unset GOOGLE_API_KEY GEMINI_API_KEY GENAI_API_KEY
   ```
4. **401 Unauthorized / Access Token Errors:**
   AI Studio API keys starting with `AQ.` are restricted. In case of `ACCESS_TOKEN_TYPE_UNSUPPORTED` errors, switch to a clean `AIzaSy...` key or authenticate Vertex AI via Application Default Credentials (ADC).
