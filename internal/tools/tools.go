package tools

import (
	"mcp-proxy/internal/globals"
	"mcp-proxy/internal/proxy"

	//
	"github.com/mark3labs/mcp-go/mcp"
)

type ToolsManagerDependencies struct {
	AppCtx *globals.ApplicationContext
	Proxy  *proxy.MCPProxy
}

type ToolsManager struct {
	dependencies ToolsManagerDependencies
}

func NewToolsManager(deps ToolsManagerDependencies) *ToolsManager {
	return &ToolsManager{
		dependencies: deps,
	}
}

func (tm *ToolsManager) AddTools() {

	// Tool 1: retrieve_tools
	retrieveToolsTool := mcp.NewTool(
		"retrieve_tools",
		mcp.WithDescription("Discover and search for available tools from the backend MCP server"),
		mcp.WithString("query",
			mcp.Description("Search query to find relevant tools"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of tools to return (default: 20, max: 100)"),
		),
		mcp.WithBoolean("debug",
			mcp.Description("Enable debug mode with detailed scoring (default: false)"),
		),
	)
	tm.dependencies.Proxy.McpServer.AddTool(retrieveToolsTool, tm.handleToolRetrieveTools)

	// Tool 2: call_tool
	callToolTool := mcp.NewTool(
		"call_tool",
		mcp.WithDescription("Execute a tool on the backend MCP server"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Tool name in format 'server:tool' (e.g., 'github:create_repository')"),
		),
		mcp.WithString("args_json",
			mcp.Description("Arguments to pass to the tool as JSON string"),
		),
	)
	tm.dependencies.Proxy.McpServer.AddTool(callToolTool, tm.handleToolCallTool)

	// Tool 3: read_cache
	readCacheTool := mcp.NewTool(
		"read_cache",
		mcp.WithDescription("Retrieve paginated data from proxy cache"),
		mcp.WithString("key",
			mcp.Required(),
			mcp.Description("Cache key provided when a response was truncated"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of records to return per page (default: 50, max: 1000)"),
		),
		mcp.WithNumber("offset",
			mcp.Description("Starting record offset for pagination (default: 0)"),
		),
	)
	tm.dependencies.Proxy.McpServer.AddTool(readCacheTool, tm.handleToolReadCache)
}
