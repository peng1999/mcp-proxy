package tools

import (
	"context"
	"encoding/json"
	"fmt"

	//
	"github.com/mark3labs/mcp-go/mcp"
)

// handleToolReadCache read cached data with pagination
func (tm *ToolsManager) handleToolReadCache(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract params
	key, err := request.RequireString("key")
	if err != nil {
		return mcp.NewToolResultError("key parameter is required"), nil
	}

	// Optional params for pagination
	args := request.GetArguments()
	limit := tm.dependencies.AppCtx.Config.Server.Options.PaginationDefaultPageSize // default
	offset := 0                                                                     // default

	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
		if limit > tm.dependencies.AppCtx.Config.Server.Options.PaginationMaxPageSize {
			limit = tm.dependencies.AppCtx.Config.Server.Options.PaginationMaxPageSize // max limit
		}
	}

	if o, ok := args["offset"].(float64); ok {
		offset = int(o)
	}

	// Look for key in cache
	tm.dependencies.Proxy.Cache.Mu.RLock()
	entry, exists := tm.dependencies.Proxy.Cache.Registry[key]
	tm.dependencies.Proxy.Cache.Mu.RUnlock()

	if !exists {
		return mcp.NewToolResultError(fmt.Sprintf("Cache key not found: %s", key)), nil
	}

	// Paginate data when needed
	paginatedData := paginateData(entry.Data, offset, limit)

	// Format response
	response := map[string]interface{}{
		"key":    key,
		"data":   paginatedData,
		"offset": offset,
		"limit":  limit,
	}

	responseBytes, _ := json.Marshal(response)
	return mcp.NewToolResultText(string(responseBytes)), nil
}
