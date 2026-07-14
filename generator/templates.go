package generator

import _ "embed"

// MainTemplate is the template for main.go in scaffolded agents.
const MainTemplate = `package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/cmd/launcher"
	"google.golang.org/adk/v2/cmd/launcher/full"
	"google.golang.org/adk/v2/model/gemini"

	"{{.ModuleName}}/agent"
)

func main() {
	// Load environment variables from .env file
	_ = godotenv.Load()

	ctx := context.Background()

	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}
	if apiKey == "" {
		apiKey = os.Getenv("GENAI_API_KEY")
	}

	modelName := os.Getenv("GEMINI_MODEL")
	if modelName == "" {
		modelName = "{{.ModelName}}"
	}

	clientConfig := &genai.ClientConfig{
		APIKey: apiKey,
	}

	// Support Vertex AI dynamically
	if os.Getenv("GOOGLE_GENAI_USE_VERTEXAI") == "1" {
		clientConfig.Backend = genai.BackendVertexAI
		clientConfig.Project = os.Getenv("GOOGLE_CLOUD_PROJECT")
		clientConfig.Location = os.Getenv("GOOGLE_CLOUD_LOCATION")
	}

	model, err := gemini.NewModel(ctx, modelName, clientConfig)
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	rootAgent, err := agent.NewAgent(model)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(rootAgent),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
`

// AgentTemplate is the template for agent/agent.go without skills.
const AgentTemplate = `package agent

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model/gemini"
	"google.golang.org/genai"
)

// Config holds the configuration options for NewAgentWithConfig.
type Config struct {
	ModelName   string
	APIKey      string
	AgentName   string
	Instruction string
	VertexAI    bool
	Project     string
	Location    string
}

// NewAgent creates a new ADK agent with defaults from environment variables.
func NewAgent(ctx context.Context) (agent.Agent, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}
	if apiKey == "" {
		apiKey = os.Getenv("GENAI_API_KEY")
	}

	modelName := os.Getenv("GEMINI_MODEL")
	if modelName == "" {
		modelName = "{{.ModelName}}"
	}

	var isVertex bool
	if os.Getenv("GOOGLE_GENAI_USE_VERTEXAI") == "1" {
		isVertex = true
	}

	return NewAgentWithConfig(ctx, Config{
		ModelName:   modelName,
		APIKey:      apiKey,
		AgentName:   "{{.AgentName}}",
		Instruction: "{{.Instruction}}",
		VertexAI:    isVertex,
		Project:     os.Getenv("GOOGLE_CLOUD_PROJECT"),
		Location:    os.Getenv("GOOGLE_CLOUD_LOCATION"),
	})
}

// NewAgentWithConfig creates a new ADK agent using explicit config parameters.
func NewAgentWithConfig(ctx context.Context, cfg Config) (agent.Agent, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("missing Gemini API Key (ensure GOOGLE_API_KEY, GEMINI_API_KEY or GENAI_API_KEY is set)")
	}
	if cfg.ModelName == "" {
		cfg.ModelName = "{{.ModelName}}"
	}
	if cfg.AgentName == "" {
		cfg.AgentName = "simple-agent"
	}

	clientConfig := &genai.ClientConfig{
		APIKey: cfg.APIKey,
	}
	if cfg.VertexAI {
		clientConfig.Backend = genai.BackendVertexAI
		clientConfig.Project = cfg.Project
		clientConfig.Location = cfg.Location
	}

	model, err := gemini.NewModel(ctx, cfg.ModelName, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize gemini model: %w", err)
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        cfg.AgentName,
		Model:       model,
		Description: "{{.Description}}",
		Instruction: cfg.Instruction,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create llm agent: %w", err)
	}

	return a, nil
}
`

// AgentWithSkillsTemplate is the template for agent/agent.go with skills.
const AgentWithSkillsTemplate = `package agent

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model/gemini"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/skilltoolset"
	"google.golang.org/adk/v2/tool/skilltoolset/skill"
	"google.golang.org/genai"
)

//go:embed skills/*
var skillsFS embed.FS

// Config holds the configuration options for NewAgentWithConfig.
type Config struct {
	ModelName   string
	APIKey      string
	AgentName   string
	Instruction string
	VertexAI    bool
	Project     string
	Location    string
}

// NewAgent creates a new ADK agent with defaults from environment variables.
func NewAgent(ctx context.Context) (agent.Agent, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}
	if apiKey == "" {
		apiKey = os.Getenv("GENAI_API_KEY")
	}

	modelName := os.Getenv("GEMINI_MODEL")
	if modelName == "" {
		modelName = "{{.ModelName}}"
	}

	var isVertex bool
	if os.Getenv("GOOGLE_GENAI_USE_VERTEXAI") == "1" {
		isVertex = true
	}

	return NewAgentWithConfig(ctx, Config{
		ModelName:   modelName,
		APIKey:      apiKey,
		AgentName:   "{{.AgentName}}",
		Instruction: "{{.Instruction}}",
		VertexAI:    isVertex,
		Project:     os.Getenv("GOOGLE_CLOUD_PROJECT"),
		Location:    os.Getenv("GOOGLE_CLOUD_LOCATION"),
	})
}

// NewAgentWithConfig creates a new ADK agent using explicit config parameters.
func NewAgentWithConfig(ctx context.Context, cfg Config) (agent.Agent, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("missing Gemini API Key (ensure GOOGLE_API_KEY, GEMINI_API_KEY or GENAI_API_KEY is set)")
	}
	if cfg.ModelName == "" {
		cfg.ModelName = "{{.ModelName}}"
	}
	if cfg.AgentName == "" {
		cfg.AgentName = "simple-agent"
	}

	clientConfig := &genai.ClientConfig{
		APIKey: cfg.APIKey,
	}
	if cfg.VertexAI {
		clientConfig.Backend = genai.BackendVertexAI
		clientConfig.Project = cfg.Project
		clientConfig.Location = cfg.Location
	}

	model, err := gemini.NewModel(ctx, cfg.ModelName, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize gemini model: %w", err)
	}

	// Setup skills sources from embedded filesystem
	subFS, err := fs.Sub(skillsFS, "skills")
	if err != nil {
		return nil, fmt.Errorf("failed to load skills filesystem: %w", err)
	}
	fileSource := skill.NewFileSystemSource(subFS)
	preloaded, _, err := skill.WithCompletePreloadSource(ctx, fileSource)
	if err != nil {
		return nil, fmt.Errorf("failed to preload skills: %w", err)
	}

	skillToolset, err := skilltoolset.New(ctx, skilltoolset.Config{
		Source: preloaded,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create skill toolset: %w", err)
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        cfg.AgentName,
		Model:       model,
		Description: "{{.Description}}",
		Instruction: cfg.Instruction,
		Toolsets: []tool.Toolset{
			skillToolset,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create llm agent: %w", err)
	}

	return a, nil
}
`

// AgentTestTemplate is the template for agent/agent_test.go.
const AgentTestTemplate = `package agent

import (
	"context"
	"strings"
	"testing"
)

func TestNewAgentWithConfig(t *testing.T) {
	tests := []struct {
		name          string
		cfg           Config
		wantErr       bool
		expectedError string
	}{
		{
			name: "missing API key (negative test)",
			cfg: Config{
				ModelName: "{{.ModelName}}",
				APIKey:    "",
			},
			wantErr:       true,
			expectedError: "missing Gemini API Key",
		},
		{
			name: "zero value config (negative test)",
			cfg:  Config{},
			wantErr:       true,
			expectedError: "missing Gemini API Key",
		},
		{
			name: "valid minimal config (happy path)",
			cfg: Config{
				ModelName: "{{.ModelName}}",
				APIKey:    "fake-api-key-for-testing",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			a, err := NewAgentWithConfig(ctx, tt.cfg)

			if (err != nil) != tt.wantErr {
				t.Fatalf("NewAgentWithConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("expected error containing %q, got %q", tt.expectedError, err.Error())
				}
			} else {
				if a == nil {
					t.Fatal("expected non-nil agent but got nil")
				}
			}
		})
	}
}
`

// EnvTemplate is the template for .env and .env.example.
const EnvTemplate = `# Gemini API Key configuration
GOOGLE_API_KEY={{.APIKey}}

# Alternatively, the SDK supports GEMINI_API_KEY or GENAI_API_KEY
GEMINI_API_KEY={{.APIKey}}
GENAI_API_KEY={{.APIKey}}

# Optional configuration
GEMINI_MODEL={{.ModelName}}
`

// GitignoreTemplate is the template for .gitignore.
const GitignoreTemplate = `# Binaries
bin/
*.exe
*.test
*.prof

# Environment variables file
.env

# IDEs
.idea/
.vscode/
*.suo
*.ntvs*
*.njsproj
*.sln
*.sw?
`

// TaskfileTemplate is the template for Taskfile.yml.
const TaskfileTemplate = `version: '3'

tasks:
  tidy:
    desc: Download and tidy Go modules dependencies
    cmds:
      - go mod tidy

  run:
    desc: Run the agent in CLI mode
    cmds:
      - go run main.go console

  run-backend:
    desc: Run the agent in local backend mode (API & Web UI)
    cmds:
      - go run main.go web --port 8000 api a2a webui

  test:
    desc: Run tests
    cmds:
      - go test -v ./...

  build:
    desc: Build the agent binary
    cmds:
      - go build -o bin/agent main.go

  clean:
    desc: Clean up build artifacts
    cmds:
      - rm -rf bin/
`

// ReadmeTemplate is the template for README.md.
const ReadmeTemplate = `# {{.AgentName}}

A Go agent built using Google ADK (Agent Development Kit) v2.0 and the Gemini API.

## Getting Started

### Prerequisites

- Go installed (version 1.23+ recommended)
- A Gemini API Key from Google AI Studio
- Task task runner installed: https://taskfile.dev/

### Setup

1. Copy ~\.env.example~ to ~\.env~ and fill in your API key:
   ~~~bash
   cp .env.example .env
   ~~~
2. Initialize dependencies:
   ~~~bash
   task tidy
   ~~~

### Running

To run the agent locally via CLI:
~~~bash
task run
~~~

To launch the local web server with API, Web UI, and A2A support:
~~~bash
task run-backend
~~~

### Testing

Run the test suite:
~~~bash
task test
~~~
`
