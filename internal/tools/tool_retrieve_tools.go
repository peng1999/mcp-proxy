package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	//
	"github.com/mark3labs/mcp-go/mcp"
)

// handleToolRetrieveTools look for available tools in backend MCP server
func (tm *ToolsManager) handleToolRetrieveTools(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if err := tm.dependencies.Proxy.InitializeBackend(ctx); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Backend connection failed: %v", err)), nil
	}

	// Extract the query from request
	query := request.GetString("query", "")

	// Get the list of tools from backend
	listRequest := mcp.ListToolsRequest{}
	listResult, err := tm.dependencies.Proxy.McpClient.ListTools(ctx, listRequest)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list tools from backend: %v", err)), nil
	}

	// Filter tools based on the query. This should speed up the retrieval from the user POV
	var relevantTools []mcp.Tool
	for _, tool := range listResult.Tools {
		if isRelevantTool(tool, query) {
			relevantTools = append(relevantTools, tool)
		}
	}

	if strings.EqualFold(query, "") || strings.EqualFold(query, "*") {
		relevantTools = listResult.Tools
	}

	// Craft the response
	response := map[string]interface{}{
		"query": query,
		"tools": relevantTools,
		"total": len(relevantTools),
	}

	responseBytes, _ := json.Marshal(response)
	return mcp.NewToolResultText(string(responseBytes)), nil
}
