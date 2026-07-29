package generator

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"google.golang.org/adk/v2/agent"
)

type mockContext struct {
	agent.Context
	ctx context.Context
}

func (m *mockContext) Deadline() (deadline time.Time, ok bool) { return m.ctx.Deadline() }
func (m *mockContext) Done() <-chan struct{}                   { return m.ctx.Done() }
func (m *mockContext) Err() error                              { return m.ctx.Err() }
func (m *mockContext) Value(key any) any                       { return m.ctx.Value(key) }

func TestScaffoldGoAgent(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "scaffold-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ctx := &mockContext{ctx: context.Background()}

	tests := []struct {
		name    string
		args    ScaffoldAgentArgs
		wantErr bool
		errMsg  string
	}{
		{
			name: "missing target path (negative test)",
			args: ScaffoldAgentArgs{
				TargetPath: "",
			},
			wantErr: true,
			errMsg:  "targetPath is required",
		},
		{
			name: "happy path minimal config",
			args: ScaffoldAgentArgs{
				TargetPath: filepath.Join(tempDir, "my_agent"),
				ModuleName: "my_agent_module",
				AgentName:  "my_agent",
			},
			wantErr: false,
		},
		{
			name: "happy path sequential agent with vertex_ai auth",
			args: ScaffoldAgentArgs{
				TargetPath: filepath.Join(tempDir, "seq_agent"),
				ModuleName: "seq_agent_module",
				AgentName:  "seq_agent",
				AgentType:  "sequential",
				AuthType:   "vertex_ai",
				StateType:  "in_memory",
			},
			wantErr: false,
		},
		{
			name: "happy path graph agent with state and oauth2 auth",
			args: ScaffoldAgentArgs{
				TargetPath: filepath.Join(tempDir, "graph_agent"),
				ModuleName: "graph_agent_module",
				AgentName:  "graph_agent",
				AgentType:  "graph",
				AuthType:   "oauth2_token",
				StateType:  "persistent",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := scaffoldGoAgent(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("scaffoldGoAgent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if res.Status != "SUCCESS" {
					t.Errorf("expected status SUCCESS, got %q", res.Status)
				}
				// Verify directories and files generated
				goModPath := filepath.Join(tt.args.TargetPath, "go.mod")
				if _, err := os.Stat(goModPath); os.IsNotExist(err) {
					t.Errorf("go.mod was not generated at %s", goModPath)
				}
				mainPath := filepath.Join(tt.args.TargetPath, "main.go")
				if data, err := os.ReadFile(mainPath); err == nil {
					content := string(data)
					if strings.Contains(content, "agent.NewAgent(model)") {
						t.Errorf("main.go contains invalid agent.NewAgent(model) call")
					}
					if !strings.Contains(content, `adkagent "google.golang.org/adk/v2/agent"`) {
						t.Errorf("main.go missing aliased adkagent import")
					}
				}
			}
		})
	}
}

func TestGenerateOKFWiki(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "okf-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ctx := &mockContext{ctx: context.Background()}

	tests := []struct {
		name    string
		args    GenerateOKFArgs
		wantErr bool
		errMsg  string
	}{
		{
			name: "missing target path (negative test)",
			args: GenerateOKFArgs{
				TargetPath: "",
				Title:      "My API",
			},
			wantErr: true,
			errMsg:  "targetPath is required",
		},
		{
			name: "missing title (negative test)",
			args: GenerateOKFArgs{
				TargetPath: filepath.Join(tempDir, "my_wiki"),
				Title:      "",
			},
			wantErr: true,
			errMsg:  "title is required",
		},
		{
			name: "happy path wiki generation",
			args: GenerateOKFArgs{
				TargetPath:  filepath.Join(tempDir, "my_wiki"),
				Title:       "Transactions Table",
				Category:    "tables",
				Description: "Stores customer transaction data.",
				FieldsJSON:  `[{"name":"id","type":"INT","description":"Transaction primary key"}]`,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := generateOKFWiki(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("generateOKFWiki() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if res.Status != "SUCCESS" {
					t.Errorf("expected status SUCCESS, got %q", res.Status)
				}
				if len(res.Paths) != 4 {
					t.Errorf("expected 4 paths generated, got %d", len(res.Paths))
				}
			}
		})
	}
}

func TestInstallSkills(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "skills-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ctx := &mockContext{ctx: context.Background()}

	tests := []struct {
		name    string
		args    InstallSkillsArgs
		wantErr bool
	}{
		{
			name: "install all skills (happy path)",
			args: InstallSkillsArgs{
				TargetPath: tempDir,
			},
			wantErr: false,
		},
		{
			name: "install specific skill (happy path)",
			args: InstallSkillsArgs{
				TargetPath: tempDir,
				SkillNames: []string{"go-adk-code"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := installSkills(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("installSkills() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if res.Status != "SUCCESS" {
					t.Errorf("expected status SUCCESS, got %q", res.Status)
				}
				if len(res.Paths) == 0 {
					t.Error("expected at least one path copied")
				}
			}
		})
	}
}

func TestFetchURL(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<html><body><h1>Hello ADK Go</h1></body></html>"))
	}))
	defer ts.Close()

	ctx := &mockContext{ctx: context.Background()}

	tests := []struct {
		name    string
		args    FetchURLArgs
		wantErr bool
		contain string
	}{
		{
			name: "empty url (negative test)",
			args: FetchURLArgs{
				URL: "",
			},
			wantErr: true,
		},
		{
			name: "invalid scheme (negative test)",
			args: FetchURLArgs{
				URL: "ftp://example.com",
			},
			wantErr: true,
		},
		{
			name: "happy path fetch HTML",
			args: FetchURLArgs{
				URL: ts.URL,
			},
			wantErr: false,
			contain: "Hello ADK Go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := fetchURLHandler(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("fetchURLHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if !strings.Contains(res.Content, tt.contain) {
					t.Errorf("expected content containing %q, got %q", tt.contain, res.Content)
				}
			}
		})
	}
}

func TestWebSearch(t *testing.T) {
	ctx := &mockContext{ctx: context.Background()}

	tests := []struct {
		name    string
		args    WebSearchArgs
		wantErr bool
	}{
		{
			name: "empty query (negative test)",
			args: WebSearchArgs{
				Query: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := webSearchHandler(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("webSearchHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateAgentsCLI(t *testing.T) {
	ctx := &mockContext{ctx: context.Background()}
	// Test invocation runs without panic
	res, err := updateAgentsCLI(ctx, UpdateAgentsCLIArgs{})
	if err != nil {
		t.Logf("updateAgentsCLI returned error (might be expected if tools missing): %v", err)
	} else {
		t.Logf("updateAgentsCLI status: %s", res.Status)
	}
}
