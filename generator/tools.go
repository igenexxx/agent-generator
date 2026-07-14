package generator

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

//go:embed skills/*
var embeddedSkillsFS embed.FS

// ScaffoldAgentArgs describes arguments for scaffolding a Go agent.
type ScaffoldAgentArgs struct {
	TargetPath  string `json:"targetPath" description:"The absolute path to initialize the Go agent project."`
	ModuleName  string `json:"moduleName" description:"The Go module name (e.g., github.com/user/myagent)."`
	AgentName   string `json:"agentName" description:"The internal name of the agent (e.g., my_agent)."`
	ModelName   string `json:"modelName" description:"The model name to use (e.g., gemini-2.5-flash)."`
	Description string `json:"description" description:"A short description of what the agent does."`
	Instruction string `json:"instruction" description:"Behavioral instruction prompts guiding the agent."`
	WithSkills  bool   `json:"withSkills" description:"Whether to bundle the agent with dynamic skilltoolsets."`
}

// ScaffoldAgentResult is the output of the scaffold tool.
type ScaffoldAgentResult struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// GenerateOKFArgs describes arguments for generating OKF documentation.
type GenerateOKFArgs struct {
	TargetPath  string `json:"targetPath" description:"The folder to write the OKF bundle to."`
	Title       string `json:"title" description:"The display title of the data asset (dataset, table, or API)."`
	Category    string `json:"category" description:"Logical category folder (e.g., tables, metrics, apis)."`
	Description string `json:"description" description:"Prose description of the asset."`
	FieldsJSON  string `json:"fieldsJson" description:"A JSON array representing fields (name, type, description)."`
}

// GenerateOKFResult is the output of the OKF generator.
type GenerateOKFResult struct {
	Status string `json:"status"`
	Paths  []string `json:"paths"`
}

// NewScaffoldTool creates the scaffold_go_agent tool.
func NewScaffoldTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:        "scaffold_go_agent",
		Description: "Generates a fully working Go ADK v2.0 Agent project at the specified path.",
	}, scaffoldGoAgent)
}

// NewGenerateOKFTool creates the generate_okf_wiki tool.
func NewGenerateOKFTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:        "generate_okf_wiki",
		Description: "Generates an Open Knowledge Format (OKF) documentation bundle for datasets/APIs.",
	}, generateOKFWiki)
}

func scaffoldGoAgent(ctx agent.Context, args ScaffoldAgentArgs) (ScaffoldAgentResult, error) {
	if args.TargetPath == "" {
		return ScaffoldAgentResult{}, fmt.Errorf("targetPath is required")
	}
	if args.ModuleName == "" {
		args.ModuleName = "my-agent-module"
	}
	if args.AgentName == "" {
		args.AgentName = "simple_agent"
	}
	if args.ModelName == "" {
		args.ModelName = "gemini-2.5-flash"
	}
	if args.Instruction == "" {
		args.Instruction = "You are a helpful assistant."
	}

	// Create directories
	agentDir := filepath.Join(args.TargetPath, "agent")
	if err := os.MkdirAll(agentDir, 0755); err != nil {
		return ScaffoldAgentResult{}, fmt.Errorf("failed to create directory tree: %w", err)
	}

	data := map[string]any{
		"ModuleName":  args.ModuleName,
		"AgentName":   args.AgentName,
		"ModelName":   args.ModelName,
		"Description": args.Description,
		"Instruction": args.Instruction,
		"APIKey":      "your_gemini_api_key_here",
	}

	// Write main.go
	if err := writeTemplate(filepath.Join(args.TargetPath, "main.go"), MainTemplate, data); err != nil {
		return ScaffoldAgentResult{}, err
	}

	// Write agent/agent.go
	agentTemplate := AgentTemplate
	if args.WithSkills {
		agentTemplate = AgentWithSkillsTemplate
	}
	if err := writeTemplate(filepath.Join(agentDir, "agent.go"), agentTemplate, data); err != nil {
		return ScaffoldAgentResult{}, err
	}

	// Write agent/agent_test.go
	if err := writeTemplate(filepath.Join(agentDir, "agent_test.go"), AgentTestTemplate, data); err != nil {
		return ScaffoldAgentResult{}, err
	}

	// Write go.mod
	goModTemplate := "module {{.ModuleName}}\n\ngo 1.23\n\nrequire (\n\tgithub.com/joho/godotenv v1.5.1\n\tgoogle.golang.org/adk/v2 v2.0.0\n\tgoogle.golang.org/genai v1.62.0\n)\n"
	if err := writeTemplate(filepath.Join(args.TargetPath, "go.mod"), goModTemplate, data); err != nil {
		return ScaffoldAgentResult{}, err
	}

	// Write .env.example & .env
	if err := writeTemplate(filepath.Join(args.TargetPath, ".env.example"), EnvTemplate, data); err != nil {
		return ScaffoldAgentResult{}, err
	}
	_ = writeTemplate(filepath.Join(args.TargetPath, ".env"), EnvTemplate, data)

	// Write .gitignore
	if err := writeTemplate(filepath.Join(args.TargetPath, ".gitignore"), GitignoreTemplate, data); err != nil {
		return ScaffoldAgentResult{}, err
	}

	// Write Taskfile.yml
	if err := writeTemplate(filepath.Join(args.TargetPath, "Taskfile.yml"), TaskfileTemplate, data); err != nil {
		return ScaffoldAgentResult{}, err
	}

	// Write README.md
	if err := writeTemplate(filepath.Join(args.TargetPath, "README.md"), ReadmeTemplate, data); err != nil {
		return ScaffoldAgentResult{}, err
	}

	// Copy skills folder if enabled
	if args.WithSkills {
		destSkills := filepath.Join(agentDir, "skills")
		_ = os.MkdirAll(destSkills, 0755)
		_ = copyEmbeddedFS(embeddedSkillsFS, "skills", destSkills)
	}

	// Run go mod tidy in target
	cmd := exec.CommandContext(ctx, "go", "mod", "tidy")
	cmd.Dir = args.TargetPath
	_ = cmd.Run()

	return ScaffoldAgentResult{
		Status:  "SUCCESS",
		Message: fmt.Sprintf("Go agent project scaffolded successfully at: %s", args.TargetPath),
	}, nil
}

func generateOKFWiki(ctx agent.Context, args GenerateOKFArgs) (GenerateOKFResult, error) {
	if args.TargetPath == "" {
		return GenerateOKFResult{}, fmt.Errorf("targetPath is required")
	}
	if args.Title == "" {
		return GenerateOKFResult{}, fmt.Errorf("title is required")
	}
	if args.Category == "" {
		args.Category = "concepts"
	}

	categoryDir := filepath.Join(args.TargetPath, args.Category)
	if err := os.MkdirAll(categoryDir, 0755); err != nil {
		return GenerateOKFResult{}, fmt.Errorf("failed to create category directory: %w", err)
	}

	timestamp := time.Now().Format(time.RFC3339)
	conceptName := strings.ToLower(strings.ReplaceAll(args.Title, " ", "_"))
	conceptPath := filepath.Join(categoryDir, conceptName+".md")

	// Parse fields
	type Field struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
	}
	var fields []Field
	if args.FieldsJSON != "" {
		_ = json.Unmarshal([]byte(args.FieldsJSON), &fields)
	}

	// Build concept.md
	var buf strings.Builder
	buf.WriteString("---\n")
	buf.WriteString(fmt.Sprintf("type: %s\n", args.Category))
	buf.WriteString(fmt.Sprintf("title: %s\n", args.Title))
	buf.WriteString(fmt.Sprintf("description: %s\n", args.Description))
	buf.WriteString(fmt.Sprintf("timestamp: %s\n", timestamp))
	buf.WriteString("---\n\n")

	buf.WriteString(fmt.Sprintf("# %s\n\n", args.Title))
	buf.WriteString(args.Description)
	buf.WriteString("\n\n")

	if len(fields) > 0 {
		buf.WriteString("## Schema\n\n")
		buf.WriteString("| Field Name | Data Type | Description |\n")
		buf.WriteString("| --- | --- | --- |\n")
		for _, f := range fields {
			buf.WriteString(fmt.Sprintf("| `%s` | `%s` | %s |\n", f.Name, f.Type, f.Description))
		}
		buf.WriteString("\n")
	}

	if err := os.WriteFile(conceptPath, []byte(buf.String()), 0644); err != nil {
		return GenerateOKFResult{}, fmt.Errorf("failed to write concept document: %w", err)
	}

	// Build index.md in category
	categoryIndex := filepath.Join(categoryDir, "index.md")
	catIdxContent := fmt.Sprintf("---\ntitle: %s Category\ndescription: List of all %s documents\n---\n\n# %s Documents\n\n- [%s](%s.md): %s\n",
		args.Category, args.Category, args.Category, args.Title, conceptName, args.Description)
	_ = os.WriteFile(categoryIndex, []byte(catIdxContent), 0644)

	// Build root index.md
	rootIndex := filepath.Join(args.TargetPath, "index.md")
	rootIdxContent := fmt.Sprintf("---\ntitle: Knowledge Wiki\ndescription: Root index catalog\n---\n\n# Welcome to the Knowledge Wiki\n\n- [%s Directory](%s/index.md): Access all %s schemas.\n",
		args.Category, args.Category, args.Category)
	_ = os.WriteFile(rootIndex, []byte(rootIdxContent), 0644)

	// Build root log.md
	rootLog := filepath.Join(args.TargetPath, "log.md")
	logContent := fmt.Sprintf("# Knowledge Bundle Changelog\n\n## [%s] - Bundle Initialization\n- Created OKF Wiki for %q\n",
		timestamp, args.Title)
	_ = os.WriteFile(rootLog, []byte(logContent), 0644)

	return GenerateOKFResult{
		Status: "SUCCESS",
		Paths: []string{
			conceptPath,
			categoryIndex,
			rootIndex,
			rootLog,
		},
	}, nil
}

func writeTemplate(path, tmpl string, data any) error {
	tmpl = strings.ReplaceAll(tmpl, "~", "`")
	t, err := template.New("tmpl").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()
	return t.Execute(f, data)
}

// InstallSkillsArgs describes arguments for installing embedded skills.
type InstallSkillsArgs struct {
	TargetPath string   `json:"targetPath" description:"Destination folder to install skills (defaults to ~/.agents/skills)."`
	SkillNames []string `json:"skillNames" description:"Specific skills to install (or all if empty). Available: adk-go-coder, blog-writer, content-research-writer, go-adk-code, go-adk-deploy, go-adk-workflow, okf-creator, okf-server-creator, seo-checklist, skill-creator."`
}

// InstallSkillsResult is the output of the install skills tool.
type InstallSkillsResult struct {
	Status string   `json:"status"`
	Paths  []string `json:"paths"`
}

// NewInstallSkillsTool creates the install_skills tool.
func NewInstallSkillsTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:        "install_skills",
		Description: "Installs selected or all embedded agent skills to a specified directory (like ~/.agents/skills or project workspace).",
	}, installSkills)
}

func installSkills(ctx agent.Context, args InstallSkillsArgs) (InstallSkillsResult, error) {
	if args.TargetPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return InstallSkillsResult{}, fmt.Errorf("failed to get user home dir: %w", err)
		}
		args.TargetPath = filepath.Join(homeDir, ".agents", "skills")
	}

	if err := os.MkdirAll(args.TargetPath, 0755); err != nil {
		return InstallSkillsResult{}, fmt.Errorf("failed to create target skills directory: %w", err)
	}

	var copiedPaths []string
	
	// If SkillNames is empty, we install all skills
	if len(args.SkillNames) == 0 {
		entries, err := embeddedSkillsFS.ReadDir("skills")
		if err != nil {
			return InstallSkillsResult{}, fmt.Errorf("failed to read embedded skills: %w", err)
		}
		for _, e := range entries {
			if e.IsDir() {
				args.SkillNames = append(args.SkillNames, e.Name())
			}
		}
	}

	for _, name := range args.SkillNames {
		srcSkillDir := "skills/" + name
		entries, err := embeddedSkillsFS.ReadDir(srcSkillDir)
		if err != nil {
			return InstallSkillsResult{}, fmt.Errorf("skill %s does not exist in embedded resources: %w", name, err)
		}
		if len(entries) == 0 {
			continue
		}

		destSkillDir := filepath.Join(args.TargetPath, name)
		if err := os.MkdirAll(destSkillDir, 0755); err != nil {
			return InstallSkillsResult{}, fmt.Errorf("failed to create destination skill folder %s: %w", name, err)
		}

		err = copyEmbeddedFS(embeddedSkillsFS, srcSkillDir, destSkillDir)
		if err != nil {
			return InstallSkillsResult{}, fmt.Errorf("failed to copy skill %s: %w", name, err)
		}
		copiedPaths = append(copiedPaths, destSkillDir)
	}

	return InstallSkillsResult{
		Status: "SUCCESS",
		Paths:  copiedPaths,
	}, nil
}

func copyEmbeddedFS(srcFS embed.FS, srcDir, destDir string) error {
	return fs.WalkDir(srcFS, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(destDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		data, err := srcFS.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(targetPath, data, 0644)
	})
}
