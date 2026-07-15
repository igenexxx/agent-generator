# 🤖 Go ADK Agent Generator & Manager Combine

[![Go Version](https://img.shields.io/badge/Go-1.23%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)
[![ADK Version](https://img.shields.io/badge/ADK-v2.0-blue?style=for-the-badge)](https://github.com/google/agents-bar)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=for-the-badge)](#)

An all-in-one interactive combine utility for creating, scaffolding, and managing Google Agent Development Kit (ADK) v2.0 projects, managing custom agent skills, and generating Open Knowledge Format (OKF) documentation in Go.

This application combines a **structured CLI menu** with a **conversational AI agent assistant** that conducts an interview, gathers parameters, and runs embedded local tools to generate production-ready configurations.

---

## 🛠️ Key Features

*   **Self-Contained Portability (`go:embed`):** The entire application compiles into a single static binary. All file templates and the skills catalog are embedded at compile time.
*   **Go ADK v2.0 Scaffolding:** Generates Go directories, client initialization boilerplate, and table-driven unit tests following Google's Go engineering conventions. Includes dynamic Vertex AI fallback configs.
*   **Decoupled Skills Installer:** Copies or installs embedded modular skills into the global directory (`~/.agents/skills/`) or any custom local workspace directory.
*   **OKF Documentation Creator:** Generates structured metadata, indices (`index.md`), and update histories (`log.md`) following the Open Knowledge Format spec.
*   **Conversational Interface:** Integrated with an ADK runner session (`runner.Runner` with `InMemoryService`) allowing the agent to ask clarification questions and call tools on-demand.

---

## 📂 Project Structure

```
.
├── generator/
│   ├── skills/          # Embedded directory containing skills metadata (.md)
│   ├── agent.go         # Manager agent & prompt instruction configuration
│   ├── templates.go     # Go templates for main.go, agent.go, Taskfile, etc.
│   ├── tools.go         # Core tool handlers (scaffold, OKF, skills installer)
│   └── tools_test.go    # Robust, adversarial table-driven unit tests
├── main.go              # CLI interactive menu loop & ADK runner execution session
├── go.mod               # Module dependencies
├── Taskfile.yml         # Automation tasks
└── .env                 # API Credentials (ignored by git)
```

---

## 📖 Embedded Skills Catalog

The generator embeds 10 reusable agent skills inside `/generator/skills/`, which can be injected into new agents or installed globally:

| Skill Directory | Purpose |
| --- | --- |
| `go-adk-workflow` | Go ADK v2.0 developer lifecycle (scaffolding, testing, console mode, debugging). |
| `go-adk-code` | Code references for Go ADK interfaces, clients, and custom adapters. |
| `go-adk-deploy` | Packaging Go agents using multi-stage Docker builds and deploying to Cloud Run. |
| `adk-go-coder` | Guidelines for Go coding adapters, models, and subcommands. |
| `okf-creator` | Rules for building structured wikis following the OKF v0.1 draft specification. |
| `okf-server-creator` | Guides for building a Go GraphQL/REST Server to read and query OKF wikis. |
| `skill-creator` | Framework for generating new custom `.md` skills for the agent skills catalog. |
| `blog-writer` | Layout rules for writing rich markdown blogs with Tailwind layouts. |
| `content-research-writer` | Internet research methodology and text extraction parsing. |
| `seo-checklist` | Checklist items for automating SEO optimization checks. |

---

## 🚀 Getting Started

### Prerequisites

*   Go 1.23+ installed.
*   A Gemini API Key (e.g. `GOOGLE_API_KEY`) from Google AI Studio.
*   Task runner installed (optional, you can run `go run` directly).

### Setup

1. Configure `.env` with your API credentials:
   ```bash
   GOOGLE_API_KEY="AIzaSy..."
   ```
2. Download dependencies:
   ```bash
   go mod tidy
   ```

---

## 🤖 Usage Reference

To start the interactive combine console, run:
```bash
go run main.go
```
Or with Task:
```bash
task run
```

### Menu Options

```
Select an option:
[1] Scaffold a new Go ADK v2.0 Agent project
[2] Generate OKF Documentation Bundle
[3] Install or Copy Agent Skills (Global/Local)
[4] Update agents-cli skills and rewrite them for Go ADK v2.0
[5] Talk with the Assistant (Free Chat)
[6] Exit
```

#### Option 1: Scaffold a Go Agent
Starts a conversation with the generator. The agent will ask you for:
*   **TargetPath:** Absolute target folder.
*   **ModuleName:** Name of the Go module (e.g. `github.com/username/project`).
*   **AgentName:** Structural name for the agent client.
*   **ModelName:** Gemini model to load (defaults to `gemini-2.5-flash`).
*   **WithSkills:** Yes/No, copies all embedded skills to the target agent folder.

Once confirmed, the agent calls the `scaffold_go_agent` tool, generates the codebase, and triggers `go mod tidy` in the target directory.

#### Option 2: Generate OKF Documentation
Converts an API, dataset, or BigQuery table into an OKF catalog. The agent asks for:
*   **TargetPath:** Target directory for the bundle.
*   **Title:** Display name.
*   **Category:** Subdirectory category (e.g. `tables`, `apis`).
*   **Fields JSON:** Table columns with types and descriptions.

The agent calls `generate_okf_wiki` and writes `index.md`, `log.md`, and the formatted table markdown file.

#### Option 3: Install Agent Skills
Copies the embedded `.md` skills to your global workspace directory or a local project.
The agent asks for the destination path (defaults to global path `~/.agents/skills`) and which skills to copy (or all if left empty), then invokes `install_skills`.

#### Option 4: Update & Rewrite Skills
Triggers the `update_agents_cli` tool, which runs `agents-cli update` (installing `uv`/`uvx` and `agents-cli` as fallbacks if missing). If skills are updated, the agent is directed to use `web_search` and `fetch_url` to research the new/updated skills from GitHub or docs, and rewrite/adapt them for Go ADK v2.0.

---

## 🛠️ Tool Catalog

In addition to scaffolding, the agent is equipped with network gathering tools:
*   `update_agents_cli`: Automates skill updates using `agents-cli` and `uv`.
*   `fetch_url`: Downloads web pages (e.g. GitHub documentation) and parses HTML into clean text.
*   `web_search`: Queries Google Custom Search (if API key/CX are set) or falls back to DuckDuckGo and Wikipedia search APIs.

---

## 🔌 MCP Integration

If a GitHub Personal Access Token (`GITHUB_PAT` or `GITHUB_TOKEN`) is set in the environment, the manager agent automatically checks the availability of the GitHub MCP server (`https://api.githubcopilot.com/mcp/`).

If reachable, it dynamically registers the GitHub MCP toolset, granting the agent native capability to read repositories, search GitHub, examine commit histories, and execute GitHub operations directly.

---

## 🧪 Testing

The unit tests follow strict Go standards, validating boundary conditions, negative scenarios (missing folders, invalid schemas), web mocks, and context cancellations.

Run the test suite:
```bash
go test -v -race ./...
```
