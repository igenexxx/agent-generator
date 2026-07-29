package generator

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golang.org/x/oauth2"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/mcptoolset"
)

// NewGeneratorAgent creates the manager/generator agent.
func NewGeneratorAgent(ctx context.Context, m model.LLM) (agent.Agent, error) {
	scaffoldTool, err := NewScaffoldTool()
	if err != nil {
		return nil, fmt.Errorf("failed to create scaffold tool: %w", err)
	}

	okfTool, err := NewGenerateOKFTool()
	if err != nil {
		return nil, fmt.Errorf("failed to create okf tool: %w", err)
	}

	installSkillsTool, err := NewInstallSkillsTool()
	if err != nil {
		return nil, fmt.Errorf("failed to create install skills tool: %w", err)
	}

	updateCLITool, err := NewUpdateAgentsCLITool()
	if err != nil {
		return nil, fmt.Errorf("failed to create update_agents_cli tool: %w", err)
	}

	fetchTool, err := NewFetchURLTool()
	if err != nil {
		return nil, fmt.Errorf("failed to create fetch_url tool: %w", err)
	}

	searchTool, err := NewWebSearchTool()
	if err != nil {
		return nil, fmt.Errorf("failed to create web_search tool: %w", err)
	}

	instruction := `You are the Go ADK Agent Generator and Manager.
Your job is to help the user build Go-based ADK agents, generate Open Knowledge Format (OKF) documentation, manage/install agent skills, and update/rewrite python-centric skills for Go.

Workflow:
1. Ask the user what they would like to do:
   - Build/scaffold a new Go agent project.
   - Generate an OKF documentation bundle.
   - Copy or install agent skills (like the Go ADK skills we created or general ones) into ~/.agents/skills or a custom path.
   - Update python agents-cli skills and rewrite updated ones for Go ADK v2.0.
2. Have a brief, clear conversation to gather the necessary details:
   - For a Go Agent: TargetPath (directory), ModuleName, AgentName, ModelName (default gemini-2.5-flash), Description, Instruction, AuthType ('api_key', 'vertex_ai', or 'oauth2_token'), StateType ('none', 'in_memory', or 'persistent'), AgentType ('simple', 'sequential', or 'graph'), and whether to include Skills (WithSkills).
   - For an OKF Wiki: TargetPath, Title, Category (e.g., tables, apis), Description, and schema fields.
   - For Skill Installation: TargetPath (default ~/.agents/skills) and SkillNames (if installing subset, e.g. ["go-adk-code", "go-adk-deploy", "go-adk-workflow"]).
   - For Skills Update: ask if they want to run the update now.
3. Confirm the collected details with the user.
4. Run the appropriate tool:
   - 'scaffold_go_agent' to generate a Go ADK project.
   - 'generate_okf_wiki' to generate an OKF wiki.
   - 'install_skills' to copy embedded skills to a directory.
   - 'update_agents_cli' to update CLI skills, and then if skills are updated, use 'web_search' or 'fetch_url' to search documentation and write updated skills matching Go ADK v2.0 guidelines.
   - If GitHub MCP is available, you will have access to high-level GitHub tools (like listing repositories, reading files, etc.). Proactively use them when searching or updating skills on GitHub!
5. Provide a summary of the action once complete.

Be conversational, supportive, and precise.`

	toolsets := []tool.Toolset{
		&remoteGatherToolset{
			tools: []tool.Tool{scaffoldTool, okfTool, installSkillsTool, updateCLITool, fetchTool, searchTool},
		},
	}

	// Dynamically register GitHub MCP if available
	if mcpSet, ok := getGithubMCP(ctx); ok {
		toolsets = append(toolsets, mcpSet)
	}

	return llmagent.New(llmagent.Config{
		Name:        "agent_generator_manager",
		Model:       m,
		Description: "Helps users scaffold Go ADK agents, manage skills, update agents-cli, and generate OKF documentation.",
		Instruction: instruction,
		Toolsets:    toolsets,
	})
}

// remoteGatherToolset wraps the generator tools.
type remoteGatherToolset struct {
	tools []tool.Tool
}

func (r *remoteGatherToolset) Name() string { return "generator_tools" }

func (r *remoteGatherToolset) Tools(ctx agent.ReadonlyContext) ([]tool.Tool, error) {
	return r.tools, nil
}

func getGithubMCP(ctx context.Context) (tool.Toolset, bool) {
	githubPAT := os.Getenv("GITHUB_PAT")
	if githubPAT == "" {
		githubPAT = os.Getenv("GITHUB_TOKEN")
	}
	if githubPAT == "" {
		return nil, false
	}

	// Check if GitHub MCP endpoint is reachable
	client := &http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.githubcopilot.com/mcp/", nil)
	if err != nil {
		return nil, false
	}
	req.Header.Set("Authorization", "Bearer "+githubPAT)
	resp, err := client.Do(req)
	if err != nil {
		return nil, false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, false
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubPAT})
	transport := &mcp.StreamableClientTransport{
		Endpoint:   "https://api.githubcopilot.com/mcp/",
		HTTPClient: oauth2.NewClient(ctx, ts),
	}

	mcpToolSet, err := mcptoolset.New(mcptoolset.Config{
		Transport: transport,
	})
	if err != nil {
		return nil, false
	}

	return mcpToolSet, true
}
