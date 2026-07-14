# Go ADK v2.0 Code Examples

This reference guide contains templates for common ADK Go v2.0 tasks.

## Full Agent Bootstrap Template

```go
package main

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model/gemini"
	"google.golang.org/adk/v2/server/adkrest"
	"google.golang.org/adk/v2/session"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/skilltoolset"
	"google.golang.org/adk/v2/tool/skilltoolset/skill"
	"google.golang.org/genai"
)

//go:embed skills/*
var skillsFS embed.FS

func main() {
	ctx := context.Background()

	// Initialize model
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Model error: %v", err)
	}

	// Initialize skill source
	subFS, err := fs.Sub(skillsFS, "skills")
	if err != nil {
		log.Fatalf("SubFS error: %v", err)
	}
	fileSource := skill.NewFileSystemSource(subFS)
	preloaded, _, err := skill.WithCompletePreloadSource(ctx, fileSource)
	if err != nil {
		log.Fatalf("Preload error: %v", err)
	}

	skillToolset, err := skilltoolset.New(ctx, skilltoolset.Config{
		Source: preloaded,
	})
	if err != nil {
		log.Fatalf("Toolset error: %v", err)
	}

	// Build Agent
	skillsAgent, err := llmagent.New(llmagent.Config{
		Name:        "go_skills_agent",
		Model:       model,
		Description: "A helper agent using Go ADK v2.",
		Instruction: "I assist with generic tasks using my skills.",
		Toolsets: []tool.Toolset{
			skillToolset,
		},
	})
	if err != nil {
		log.Fatalf("Agent error: %v", err)
	}

	// Serve
	loader := agent.NewSingleLoader(skillsAgent)
	adkServer, err := adkrest.NewServer(adkrest.ServerConfig{
		AgentLoader:     loader,
		SessionService:  session.InMemoryService(),
		SSEWriteTimeout: 120 * time.Second,
	})
	if err != nil {
		log.Fatalf("REST Server error: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", adkServer))

	log.Info("Starting server on :8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
```

## Creating Custom Tools

Custom tools must adhere to the signature expected by `functiontool.New`:

```go
type FetchDataArgs struct {
	Endpoint string `json:"endpoint" description:"The URL endpoint to fetch data from"`
}

type FetchDataResult struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

func RegisterFetchTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:        "fetch_data",
		Description: "Fetches payload data from an API endpoint.",
	}, func(ctx tool.Context, args FetchDataArgs) (FetchDataResult, error) {
		if args.Endpoint == "" {
			return FetchDataResult{}, errors.New("endpoint cannot be empty")
		}
		// Custom HTTP Client request
		return FetchDataResult{Status: 200, Body: "Success"}, nil
	})
}
```

## Modern Agent Bootstrap Template with Launcher (Recommended)

In modern ADK applications, the standard way to run the agent is using `google.golang.org/adk/v2/cmd/launcher/full` which automatically handles both CLI and A2A API routes:

```go
package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"os"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/cmd/launcher"
	"google.golang.org/adk/v2/cmd/launcher/full"
	"google.golang.org/adk/v2/model/gemini"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/skilltoolset"
	"google.golang.org/adk/v2/tool/skilltoolset/skill"
	"google.golang.org/genai"
)

//go:embed skills/*
var skillsFS embed.FS

func main() {
	ctx := context.Background()

	// Initialize model
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Model error: %v", err)
	}

	// Initialize skill source from embedded filesystem
	subFS, err := fs.Sub(skillsFS, "skills")
	if err != nil {
		log.Fatalf("SubFS error: %v", err)
	}
	fileSource := skill.NewFileSystemSource(subFS)
	preloaded, _, err := skill.WithCompletePreloadSource(ctx, fileSource)
	if err != nil {
		log.Fatalf("Preload error: %v", err)
	}

	skillToolset, err := skilltoolset.New(ctx, skilltoolset.Config{
		Source: preloaded,
	})
	if err != nil {
		log.Fatalf("Toolset error: %v", err)
	}

	// Build Agent
	skillsAgent, err := llmagent.New(llmagent.Config{
		Name:        "go_skills_agent",
		Model:       model,
		Description: "A modern helper agent using Go ADK v2 with launcher.",
		Instruction: "I assist with generic tasks using my skills.",
		Toolsets: []tool.Toolset{
			skillToolset,
		},
	})
	if err != nil {
		log.Fatalf("Agent error: %v", err)
	}

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(skillsAgent),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Execution failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
```
