package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
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

	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}
	if apiKey == "" {
		apiKey = os.Getenv("GENAI_API_KEY")
	}

	if apiKey == "" {
		log.Fatalf("Missing API Key (ensure GOOGLE_API_KEY, GEMINI_API_KEY or GENAI_API_KEY is set)")
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
		modelName = "gemini-2.5-flash"
	}

	model, err := gemini.NewModel(ctx, modelName, clientConfig)
	if err != nil {
		log.Fatalf("Failed to initialize model: %v", err)
	}

	genAgent, err := generator.NewGeneratorAgent(ctx, model)
	if err != nil {
		log.Fatalf("Failed to create generator agent: %v", err)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n========================================================")
		fmt.Println("🤖  Go ADK Agent Generator & Manager Combine  🤖")
		fmt.Println("========================================================")
		fmt.Println("Select an option:")
		fmt.Println("[1] Scaffold a new Go ADK v2.0 Agent project")
		fmt.Println("[2] Generate OKF Documentation Bundle")
		fmt.Println("[3] Install or Copy Agent Skills (Global/Local)")
		fmt.Println("[4] Update agents-cli skills and rewrite them for Go ADK v2.0")
		fmt.Println("[5] Talk with the Assistant (Free Chat)")
		fmt.Println("[6] Exit")

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
			fmt.Println("Goodbye!")
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
			fmt.Println("\nReturning to main menu...")
			break
		}

		userMsg = &genai.Content{
			Role:  "user",
			Parts: []*genai.Part{{Text: input}},
		}
	}
}
