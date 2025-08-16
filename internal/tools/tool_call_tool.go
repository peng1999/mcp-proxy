package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	//
	"github.com/mark3labs/mcp-go/mcp"

	//
	"mcp-proxy/internal/cache"
)

// handleToolCallTool execute the tool from the backend
func (tm *ToolsManager) handleToolCallTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if err := tm.dependencies.Proxy.InitializeBackend(ctx); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Backend connection failed: %v", err)), nil
	}

	// Extract params
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError("name parameter is required"), nil
	}

	argsJsonStr, err := request.RequireString("args_json")
	if err != nil {
		return mcp.NewToolResultError("args_json parameter is required"), nil
	}

	// Parse JSON arguments
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(argsJsonStr), &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid JSON in args_json: %v", err)), nil
	}

	// Delete tool prefix when existing
	backendToolName := mapToolName(name)

	// Craft and execute backend request
	backendRequest := mcp.CallToolRequest{}
	backendRequest.Params.Name = backendToolName
	backendRequest.Params.Arguments = args

	result, err := tm.dependencies.Proxy.McpClient.CallTool(ctx, backendRequest)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Backend tool execution failed: %v", err)), nil
	}

	// For big response, use cache to store it
	resultJson, _ := json.Marshal(result)
	if len(resultJson) > tm.dependencies.AppCtx.Config.Server.Options.CacheThresholdBytes {
		cacheKey := cache.GenerateCacheKey()
		tm.dependencies.Proxy.Cache.Mu.Lock()
		tm.dependencies.Proxy.Cache.Registry[cacheKey] = cache.CacheEntry{
			Data:      result,
			Timestamp: time.Now().Unix(),
		}
		tm.dependencies.Proxy.Cache.Mu.Unlock()

		// Return reference cache key
		cacheResponse := map[string]interface{}{
			"cached":    true,
			"cache_key": cacheKey,
			"message":   fmt.Sprintf("Response cached due to size. Use read_cache with key: %s", cacheKey),
		}
		cacheBytes, _ := json.Marshal(cacheResponse)
		return mcp.NewToolResultText(string(cacheBytes)), nil
	}

	return result, nil
}
