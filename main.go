package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/model/gemini"
	"google.golang.org/adk/v2/runner"
	"google.golang.org/adk/v2/session"
	"google.golang.org/genai"

	"agentgenerator/generator"
)

func main() {
	_ = godotenv.Load()
	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n========================================================")
		fmt.Println("🤖  Go ADK Agent Generator & Manager Combine  🤖")
		fmt.Println("========================================================")
		fmt.Println("Select operation mode:")
		fmt.Println("[1] Classic Menu (Direct CLI prompts without AI)")
		fmt.Println("[2] AI Agent Assistant (Interactive Chat with LLM Agent)")
		fmt.Println("[3] Exit")

		fmt.Print("\nEnter choice (1-3): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > 3 {
			fmt.Println("❌ Invalid choice, please enter 1, 2, or 3.")
			continue
		}

		switch choice {
		case 1:
			runClassicMenu(ctx, reader)
		case 2:
			runAIMode(ctx, reader)
		case 3:
			fmt.Println("Goodbye!")
			return
		}
	}
}

func runClassicMenu(ctx context.Context, reader *bufio.Reader) {
	for {
		fmt.Println("\n--------------------------------------------------------")
		fmt.Println("📋  Classic Mode (Direct CLI Prompts)")
		fmt.Println("--------------------------------------------------------")
		fmt.Println("Select an action:")
		fmt.Println("[1] Scaffold a new Go ADK v2.0 Agent project")
		fmt.Println("[2] Generate OKF Documentation Bundle")
		fmt.Println("[3] Install or Copy Agent Skills")
		fmt.Println("[4] Update agents-cli skills")
		fmt.Println("[5] Back to main mode selection")

		fmt.Print("\nEnter choice (1-5): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > 5 {
			fmt.Println("❌ Invalid choice, please enter a number between 1 and 5.")
			continue
		}

		if choice == 5 {
			return
		}

		switch choice {
		case 1:
			promptScaffoldAgent(ctx, reader)
		case 2:
			promptGenerateOKF(ctx, reader)
		case 3:
			promptInstallSkills(ctx, reader)
		case 4:
			promptUpdateAgentsCLI(ctx)
		}
	}
}

func promptScaffoldAgent(ctx context.Context, reader *bufio.Reader) {
	fmt.Println("\n🛠️  Scaffolding New Go ADK Agent...")
	targetPath := readInputOrDefault(reader, "Target directory", "./my-agent")
	moduleName := readInputOrDefault(reader, "Go module name", "my-agent")
	agentName := readInputOrDefault(reader, "Agent name", "simple_agent")
	modelName := readInputOrDefault(reader, "Gemini model name", "gemini-3.5-flash-lite")
	description := readInputOrDefault(reader, "Short description", "A Go ADK v2.0 Agent")
	instruction := readInputOrDefault(reader, "Instruction prompt", "You are a helpful assistant.")
	authType := readInputOrDefault(reader, "Authentication type (api_key/vertex_ai/oauth2_token)", "api_key")
	stateType := readInputOrDefault(reader, "State management (none/in_memory/persistent)", "none")
	agentType := readInputOrDefault(reader, "Agent architecture (simple/sequential/graph)", "simple")
	withSkillsStr := readInputOrDefault(reader, "Include bundled skills? (y/N)", "N")

	withSkills := strings.ToLower(withSkillsStr) == "y" || strings.ToLower(withSkillsStr) == "yes"

	args := generator.ScaffoldAgentArgs{
		TargetPath:  targetPath,
		ModuleName:  moduleName,
		AgentName:   agentName,
		ModelName:   modelName,
		Description: description,
		Instruction: instruction,
		AuthType:    authType,
		StateType:   stateType,
		AgentType:   agentType,
		WithSkills:  withSkills,
	}

	res, err := generator.ScaffoldGoAgent(ctx, args)
	if err != nil {
		fmt.Printf("❌ Failed to scaffold agent: %v\n", err)
	} else {
		fmt.Printf("✅ %s\n", res.Message)
	}
}

func promptGenerateOKF(ctx context.Context, reader *bufio.Reader) {
	fmt.Println("\n📚 Generating OKF Documentation Bundle...")
	targetPath := readInputOrDefault(reader, "Target directory", "./docs/okf")
	title := readInputOrDefault(reader, "Asset title", "Data Asset")
	category := readInputOrDefault(reader, "Category (tables/apis/concepts/metrics)", "concepts")
	description := readInputOrDefault(reader, "Description", "Documentation for data asset.")
	fieldsJSON := readInputOrDefault(reader, "Fields JSON array (optional)", "")

	args := generator.GenerateOKFArgs{
		TargetPath:  targetPath,
		Title:       title,
		Category:    category,
		Description: description,
		FieldsJSON:  fieldsJSON,
	}

	res, err := generator.GenerateOKFWiki(ctx, args)
	if err != nil {
		fmt.Printf("❌ Failed to generate OKF bundle: %v\n", err)
	} else {
		fmt.Printf("✅ Generated OKF Wiki at %s with files:\n", targetPath)
		for _, p := range res.Paths {
			fmt.Printf("   - %s\n", p)
		}
	}
}

func promptInstallSkills(ctx context.Context, reader *bufio.Reader) {
	fmt.Println("\n📦 Installing Agent Skills...")
	homeDir, _ := os.UserHomeDir()
	defaultPath := "./skills"
	if homeDir != "" {
		defaultPath = homeDir + "/.agents/skills"
	}

	targetPath := readInputOrDefault(reader, "Target directory", defaultPath)
	skillsStr := readInputOrDefault(reader, "Specific skills (comma separated, empty for all)", "")

	var skillNames []string
	if skillsStr != "" {
		parts := strings.Split(skillsStr, ",")
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				skillNames = append(skillNames, trimmed)
			}
		}
	}

	args := generator.InstallSkillsArgs{
		TargetPath: targetPath,
		SkillNames: skillNames,
	}

	res, err := generator.InstallSkills(ctx, args)
	if err != nil {
		fmt.Printf("❌ Failed to install skills: %v\n", err)
	} else {
		fmt.Printf("✅ Installed %d skills to %s:\n", len(res.Paths), targetPath)
		for _, p := range res.Paths {
			fmt.Printf("   - %s\n", p)
		}
	}
}

func promptUpdateAgentsCLI(ctx context.Context) {
	fmt.Println("\n🔄 Updating agents-cli...")
	res, err := generator.UpdateAgentsCLI(ctx, generator.UpdateAgentsCLIArgs{})
	if err != nil {
		fmt.Printf("❌ Error updating agents-cli: %v\n", err)
	} else {
		fmt.Printf("✅ Update status: %s\n", res.Status)
		if res.Stdout != "" {
			fmt.Println(res.Stdout)
		}
	}
}

func runAIMode(ctx context.Context, reader *bufio.Reader) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}
	if apiKey == "" {
		apiKey = os.Getenv("GENAI_API_KEY")
	}

	if apiKey == "" {
		fmt.Println("❌ Missing Gemini API Key. Please set GOOGLE_API_KEY, GEMINI_API_KEY or GENAI_API_KEY to use the AI Agent Assistant.")
		return
	}

	clientConfig := &genai.ClientConfig{
		APIKey: apiKey,
	}
	if os.Getenv("GOOGLE_GENAI_USE_VERTEXAI") == "1" {
		clientConfig.Backend = genai.BackendVertexAI
		clientConfig.Project = os.Getenv("GOOGLE_CLOUD_PROJECT")
		clientConfig.Location = os.Getenv("GOOGLE_CLOUD_LOCATION")
	}

	modelName := os.Getenv("GEMINI_MODEL")
	if modelName == "" {
		modelName = "gemini-3.5-flash-lite"
	}

	model, err := gemini.NewModel(ctx, modelName, clientConfig)
	if err != nil {
		fmt.Printf("❌ Failed to initialize model: %v\n", err)
		return
	}

	genAgent, err := generator.NewGeneratorAgent(ctx, model)
	if err != nil {
		fmt.Printf("❌ Failed to create generator agent: %v\n", err)
		return
	}

	for {
		fmt.Println("\n--------------------------------------------------------")
		fmt.Println("💬  AI Agent Assistant (Interactive Chat)")
		fmt.Println("--------------------------------------------------------")
		fmt.Println("Select a starting prompt or conversation topic:")
		fmt.Println("[1] Scaffold a new Go ADK v2.0 Agent project")
		fmt.Println("[2] Generate OKF Documentation Bundle")
		fmt.Println("[3] Install or Copy Agent Skills")
		fmt.Println("[4] Update agents-cli skills and adapt for Go ADK v2.0")
		fmt.Println("[5] Talk with the Assistant (Free Chat)")
		fmt.Println("[6] Back to main mode selection")

		fmt.Print("\nEnter choice (1-6): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > 6 {
			fmt.Println("❌ Invalid choice, please enter a number between 1 and 6.")
			continue
		}

		if choice == 6 {
			return
		}

		var startPrompt string
		switch choice {
		case 1:
			startPrompt = "I want to scaffold a new Go ADK v2.0 agent project. Please ask me questions to gather the required parameters and create it."
		case 2:
			startPrompt = "I want to generate an Open Knowledge Format (OKF) documentation bundle. Please gather the metadata and create it."
		case 3:
			startPrompt = "I want to install or copy some agent skills into a directory (like ~/.agents/skills or a custom project path). Please help me do this."
		case 4:
			startPrompt = "I want to run agents-cli update to update skills, and if they changed, rewrite/adapt them for Go ADK v2.0."
		case 5:
			startPrompt = "Hello! I want to talk with you about Go ADK v2.0 or general code generation."
		}

		runChatSession(ctx, genAgent, startPrompt, reader)
	}
}

func runChatSession(ctx context.Context, genAgent agent.Agent, startPrompt string, reader *bufio.Reader) {
	sessionService := session.InMemoryService()
	r, err := runner.New(runner.Config{
		AppName:           "agent_generator",
		Agent:             genAgent,
		SessionService:    sessionService,
		AutoCreateSession: true,
	})
	if err != nil {
		fmt.Printf("❌ Failed to create session runner: %v\n", err)
		return
	}

	sessionID := "interactive_session"
	userID := "user"

	userMsg := &genai.Content{
		Role:  "user",
		Parts: []*genai.Part{{Text: startPrompt}},
	}

	fmt.Println("\n🤖 Starting conversation...")
	fmt.Println(strings.Repeat("-", 60))

	for {
		events := r.Run(ctx, userID, sessionID, userMsg, agent.RunConfig{})
		for event, err := range events {
			if err != nil {
				fmt.Printf("\n❌ Agent Error: %v\n", err)
				break
			}

			if event.Content != nil {
				for _, part := range event.Content.Parts {
					if part.FunctionCall != nil {
						fmt.Printf("\n🛠️  [Executing tool %s...]\n", part.FunctionCall.Name)
					}
					if part.Text != "" {
						fmt.Print(part.Text)
					}
				}
			}
		}
		fmt.Println()

		fmt.Print("\nYou: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		input = strings.TrimSpace(input)
		if strings.ToLower(input) == "exit" || strings.ToLower(input) == "quit" || input == "" {
			fmt.Println("\nReturning to AI menu...")
			break
		}

		userMsg = &genai.Content{
			Role:  "user",
			Parts: []*genai.Part{{Text: input}},
		}
	}
}

func readInputOrDefault(reader *bufio.Reader, promptText, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", promptText, defaultValue)
	} else {
		fmt.Printf("%s: ", promptText)
	}
	input, err := reader.ReadString('\n')
	if err != nil {
		return defaultValue
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}
