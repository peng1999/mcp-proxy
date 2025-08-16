package tools

import (
	"strings"

	//
	"github.com/mark3labs/mcp-go/mcp"
)

// isRelevantTool decides whether a tool is relevant based
// on a query that must be contained in name/description
func isRelevantTool(tool mcp.Tool, query string) bool {
	// Simple text matching en nombre y descripción
	queryLower := strings.ToLower(query)
	nameLower := strings.ToLower(tool.Name)
	descLower := strings.ToLower(tool.Description)

	return strings.Contains(nameLower, queryLower) ||
		strings.Contains(descLower, queryLower)
}

// mapToolName return the name of a tool without the prefix
// server:tool -> tool
func mapToolName(frontendName string) string {
	parts := strings.Split(frontendName, ":")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return frontendName
}

// paginateData return a page of data results based on passed offset and limit
func paginateData(data interface{}, offset, limit int) interface{} {

	// Simple pagination for arrays/slices
	switch v := data.(type) {
	case []interface{}:
		total := len(v)
		start := offset
		if start > total {
			return []interface{}{}
		}
		end := start + limit
		if end > total {
			end = total
		}
		return v[start:end]
	default:
		// Pagination is not possible
		return data
	}
}
