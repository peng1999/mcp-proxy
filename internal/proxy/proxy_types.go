package proxy

import (
	"sync"

	//
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/server"

	//
	"mcp-proxy/internal/cache"
	"mcp-proxy/internal/globals"
)

type MCPProxyDependencies struct {
	AppContext *globals.ApplicationContext
}

type MCPProxy struct {
	//
	Dependencies MCPProxyDependencies

	//
	Mu sync.RWMutex

	//
	McpServer *server.MCPServer
	McpClient *client.Client
	Cache     *cache.Cache

	//
	BackendURL  string
	Initialized bool
}
