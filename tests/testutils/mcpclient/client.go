// Copyright 2025-2026 The MathWorks, Inc.

package mcpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// syncBuffer is a thread-safe wrapper around bytes.Buffer
type syncBuffer struct {
	lock        sync.Mutex
	innerBuffer bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (n int, err error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.innerBuffer.Write(p)
}

func (b *syncBuffer) String() string {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.innerBuffer.String()
}

// MCPClient wraps MCP session with helper methods
// Lifecycle: NewClient -> CreateSession -> [operations] -> Close
type MCPClient struct {
	client     *mcp.Client
	clientImpl *mcp.Implementation
	transport  *mcp.CommandTransport
	stderr     *syncBuffer
}

// MCPClientSession represents an active MCP session
type MCPClientSession struct {
	session    *mcp.ClientSession
	clientImpl *mcp.Implementation
	stderr     *syncBuffer
}

func GetMCPClientImplementation() *mcp.Implementation {
	return &mcp.Implementation{
		Name:       "test-client-for-matlab-mcp-core-server",
		Title:      "Test Client",
		Version:    "v.0.1-some-version",
		WebsiteURL: "https://test-client.example.com",
	}
}

// NewClient creates a new MCP client with roots capability enabled.
func NewClient(ctx context.Context, serverPath string, env []string, args ...string) *MCPClient {
	return newClient(ctx, serverPath, env, &mcp.ClientOptions{
		Capabilities: &mcp.ClientCapabilities{
			RootsV2: &mcp.RootCapabilities{ListChanged: true},
		},
	}, args...)
}

// NewClientWithoutRootsCapability creates a new MCP client that does not
// advertise roots capability. The server will not request roots from this client.
func NewClientWithoutRootsCapability(ctx context.Context, serverPath string, env []string, args ...string) *MCPClient {
	return newClient(ctx, serverPath, env, &mcp.ClientOptions{
		Capabilities: &mcp.ClientCapabilities{},
	}, args...)
}

func newClient(ctx context.Context, serverPath string, env []string, clientOpts *mcp.ClientOptions, args ...string) *MCPClient {
	cmd := exec.CommandContext(ctx, serverPath, args...)
	stderr := &syncBuffer{}
	cmd.Stderr = stderr
	if env != nil {
		cmd.Env = env
	}
	// Use a longer terminate duration to allow the server to gracefully shut down
	// MATLAB sessions, which can take some time.
	transport := &mcp.CommandTransport{
		Command:           cmd,
		TerminateDuration: 3 * time.Minute,
	}
	impl := GetMCPClientImplementation()
	client := mcp.NewClient(impl, clientOpts)

	return &MCPClient{
		client:     client,
		clientImpl: impl,
		transport:  transport,
		stderr:     stderr,
	}
}

// CreateSessionOption configures the client before connecting to the server.
type CreateSessionOption func(*MCPClient)

// WithRoots returns a CreateSessionOption that adds roots to the client before
// the MCP handshake. This ensures the server receives roots during initialization.
func WithRoots(roots ...*mcp.Root) CreateSessionOption {
	return func(c *MCPClient) {
		c.client.AddRoots(roots...)
	}
}

// AddRoots adds roots to the MCP client. Roots can be added before or after
// creating a session. If a session is already connected, the server is notified.
func (c *MCPClient) AddRoots(roots ...*mcp.Root) {
	c.client.AddRoots(roots...)
}

func (c *MCPClient) CreateSession(ctx context.Context, opts ...CreateSessionOption) (*MCPClientSession, error) {
	for _, opt := range opts {
		opt(c)
	}
	session, err := c.client.Connect(ctx, c.transport, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP client session: %w", err)
	}
	return &MCPClientSession{
		session:    session,
		clientImpl: c.clientImpl,
		stderr:     c.stderr,
	}, nil
}

func (s *MCPClientSession) ClientName() string {
	return s.clientImpl.Name
}

func (s *MCPClientSession) ClientVersion() string {
	return s.clientImpl.Version
}

func (s *MCPClientSession) ClientTitle() string {
	return s.clientImpl.Title
}

func (s *MCPClientSession) ClientWebsiteURL() string {
	return s.clientImpl.WebsiteURL
}

func (s *MCPClientSession) Close() error {
	err := s.session.Close()
	if err != nil {
		return fmt.Errorf("session close failed: %w\n\nServer stderr:\n%s", err, s.stderr.String())
	}
	return nil
}

// Stderr returns the captured stderr output from the MCP server process.
// This is useful for debugging test failures.
func (s *MCPClientSession) Stderr() string {
	return s.stderr.String()
}

// InitializeResult returns the result of the Initialize call.
func (s *MCPClientSession) InitializeResult() *mcp.InitializeResult {
	return s.session.InitializeResult()
}

func (s *MCPClientSession) ListTools(ctx context.Context, params *mcp.ListToolsParams) (*mcp.ListToolsResult, error) {
	return s.session.ListTools(ctx, params)
}

func (s *MCPClientSession) ListResources(ctx context.Context, params *mcp.ListResourcesParams) (*mcp.ListResourcesResult, error) {
	return s.session.ListResources(ctx, params)
}

// CallTool calls an MCP tool and asserts it doesn't error
func (s *MCPClientSession) CallTool(ctx context.Context, name string, args map[string]any) (*mcp.CallToolResult, error) {
	result, err := s.session.CallTool(ctx, &mcp.CallToolParams{
		Name:      name,
		Arguments: args,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call tool %s: %w", name, err)
	}

	if result.IsError {
		// Log the error details for better debugging
		errorMsg := "unknown error"
		if textContent, err := s.GetTextContent(result); err == nil {
			errorMsg = textContent
		}
		return nil, fmt.Errorf("tool %s returned an error: %s", name, errorMsg)
	}
	return result, nil
}

// GetTextContent extracts text content from a tool result
func (*MCPClientSession) GetTextContent(result *mcp.CallToolResult) (string, error) {
	if result.Content == nil {
		return "", fmt.Errorf("result content is nil")
	}
	if len(result.Content) == 0 {
		return "", fmt.Errorf("result should have content")
	}
	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		return "", fmt.Errorf("content should be TextContent")
	}
	return textContent.Text, nil
}

// UnmarshalStructuredContent unmarshals structured content into a target
func (*MCPClientSession) UnmarshalStructuredContent(result *mcp.CallToolResult, target interface{}) error {
	jsonBytes, err := json.Marshal(result.StructuredContent)
	if err != nil {
		return fmt.Errorf("failed to marshal structured content: %w", err)
	}
	err = json.Unmarshal(jsonBytes, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal structured content: %w", err)
	}
	return nil
}

// EvaluateCode evaluates MATLAB code
func (s *MCPClientSession) EvaluateCode(ctx context.Context, code string, projectPath ...string) (string, error) {
	args := map[string]any{"code": code}
	if len(projectPath) > 0 {
		args["project_path"] = projectPath[0]
	}
	result, err := s.CallTool(ctx, "evaluate_matlab_code", args)
	if err != nil {
		return "", err
	}
	return s.GetTextContent(result)
}

// CodeIssue represents a structured code issue returned by check_matlab_code.
type CodeIssue struct {
	Description string `json:"description"`
	Line        int    `json:"line"`
	StartColumn int    `json:"start_column"`
	EndColumn   int    `json:"end_column"`
	Severity    string `json:"severity"`
	Fixable     bool   `json:"fixable"`
}

// CheckCode checks MATLAB code
func (s *MCPClientSession) CheckCode(ctx context.Context, scriptPath string) ([]CodeIssue, error) {
	result, err := s.CallTool(ctx, "check_matlab_code", map[string]any{
		"script_path": scriptPath,
	})
	if err != nil {
		return nil, err
	}
	var output struct {
		CodeIssues []CodeIssue `json:"code_issues"`
	}
	err = s.UnmarshalStructuredContent(result, &output)
	if err != nil {
		return nil, err
	}
	return output.CodeIssues, nil
}

// RunFile runs a MATLAB file
func (s *MCPClientSession) RunFile(ctx context.Context, scriptPath string) (string, error) {
	result, err := s.CallTool(ctx, "run_matlab_file", map[string]any{
		"script_path": scriptPath,
	})
	if err != nil {
		return "", err
	}
	return s.GetTextContent(result)
}

// RunTestFile runs a MATLAB test file
func (s *MCPClientSession) RunTestFile(ctx context.Context, scriptPath string) (string, error) {
	result, err := s.CallTool(ctx, "run_matlab_test_file", map[string]any{
		"script_path": scriptPath,
	})
	if err != nil {
		return "", err
	}
	return s.GetTextContent(result)
}

// DetectToolboxes detects installed MATLAB toolboxes
func (s *MCPClientSession) DetectToolboxes(ctx context.Context) (string, error) {
	result, err := s.CallTool(ctx, "detect_matlab_toolboxes", map[string]any{})
	if err != nil {
		return "", err
	}
	var output struct {
		InstallationInfo string `json:"installation_info"`
	}
	err = s.UnmarshalStructuredContent(result, &output)
	if err != nil {
		return "", err
	}
	return output.InstallationInfo, nil
}

// NewSessionManager creates a new session manager for multi-session workflows
func (s *MCPClientSession) NewSessionManager() *SessionManager {
	return &SessionManager{
		session: s,
	}
}

// ReadResource reads an MCP resource by URI and returns its text content
func (s *MCPClientSession) ReadResource(ctx context.Context, uri string) (string, error) {
	result, err := s.session.ReadResource(ctx, &mcp.ReadResourceParams{
		URI: uri,
	})
	if err != nil {
		return "", fmt.Errorf("failed to read resource %s: %w", uri, err)
	}
	if len(result.Contents) == 0 {
		return "", fmt.Errorf("resource %s returned no contents", uri)
	}
	if strings.TrimSpace(result.Contents[0].Text) == "" {
		return "", fmt.Errorf("resource %s returned empty text content", uri)
	}
	return result.Contents[0].Text, nil
}
