package generator

import (
	"fmt"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/tool"
)

// NewGeneratorAgent creates the manager/generator agent.
func NewGeneratorAgent(m model.LLM) (agent.Agent, error) {
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

	instruction := `You are the Go ADK Agent Generator and Manager.
Your job is to help the user build Go-based ADK agents, generate Open Knowledge Format (OKF) documentation, or manage/install agent skills.

Workflow:
1. Ask the user what they would like to do:
   - Build/scaffold a new Go agent project.
   - Generate an OKF documentation bundle.
   - Copy or install agent skills (like the Go ADK skills we created or general ones) into ~/.agents/skills or a custom path.
2. Have a brief, clear conversation to gather the necessary details:
   - For a Go Agent: TargetPath (directory), ModuleName, AgentName, ModelName (default gemini-2.5-flash), Description, Instruction, and whether to include Skills (WithSkills).
   - For an OKF Wiki: TargetPath, Title, Category (e.g., tables, apis), Description, and schema fields.
   - For Skill Installation: TargetPath (default ~/.agents/skills) and SkillNames (if installing subset, e.g. ["go-adk-code", "go-adk-deploy", "go-adk-workflow"]).
3. Confirm the collected details with the user.
4. Run the appropriate tool:
   - 'scaffold_go_agent' to generate a Go ADK project.
   - 'generate_okf_wiki' to generate an OKF wiki.
   - 'install_skills' to copy embedded skills to a directory.
5. Provide a summary of the action once complete.

Be conversational, supportive, and precise.`

	return llmagent.New(llmagent.Config{
		Name:        "agent_generator_manager",
		Model:       m,
		Description: "Helps users scaffold Go ADK agents, manage skills, and generate OKF documentation.",
		Instruction: instruction,
		Toolsets: []tool.Toolset{
			&remoteGatherToolset{
				tools: []tool.Tool{scaffoldTool, okfTool, installSkillsTool},
			},
		},
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
