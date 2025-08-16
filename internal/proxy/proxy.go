package proxy

import (
	"context"
	"fmt"
	"log"

	//
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"

	//
	"mcp-proxy/internal/cache"
)

func NewMCPProxy(deps MCPProxyDependencies) *MCPProxy {

	tmpCache := cache.NewCache()
	return &MCPProxy{
		Dependencies: deps,
		Cache:        tmpCache,
	}
}

// InitializeBackend init the connection with the backend MCP
func (p *MCPProxy) InitializeBackend(ctx context.Context) (err error) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	if p.Initialized {
		return nil
	}

	var mcpClient *client.Client
	switch p.Dependencies.AppContext.Config.Backend.Transport.Type {
	case "http":
		mcpClient, err = client.NewStreamableHttpClient(p.Dependencies.AppContext.Config.Backend.Transport.HTTP.URL,
			[]transport.StreamableHTTPCOption{
				transport.WithHTTPHeaders(p.Dependencies.AppContext.Config.Backend.Transport.HTTP.Headers),
				//transport.WithSession("custom_session"),
			}...)
	default:
		mcpClient, err = client.NewStdioMCPClient(p.Dependencies.AppContext.Config.Backend.Transport.Stdio.Command,
			p.Dependencies.AppContext.Config.Backend.Transport.Stdio.Env,
			p.Dependencies.AppContext.Config.Backend.Transport.Stdio.Args...)
	}

	if err != nil {
		return fmt.Errorf("failed creating backend MCP client: %s", err.Error())
	}

	// Init connection
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "backend",
		Version: "1.0.0",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{
		Experimental: make(map[string]interface{}),
		Roots:        nil,
		Sampling:     nil,
	}

	_, err = mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		return fmt.Errorf("failed to initialize backend connection: %w", err)
	}

	p.McpClient = mcpClient
	p.Initialized = true

	log.Printf("Successfully connected to backend MCP server")
	return nil
}
