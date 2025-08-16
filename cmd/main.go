package main

import (
	"log"
	"net/http"
	"time"

	//
	"mcp-proxy/internal/globals"
	"mcp-proxy/internal/handlers"
	"mcp-proxy/internal/middlewares"
	"mcp-proxy/internal/proxy"
	"mcp-proxy/internal/tools"

	//
	"github.com/mark3labs/mcp-go/server"
)

func main() {

	// 0. Process the configuration
	appCtx, err := globals.NewApplicationContext()
	if err != nil {
		log.Fatalf("failed creating application context: %v", err.Error())
	}

	// 1. Initialize middlewares that need it
	accessLogsMw := middlewares.NewAccessLogsMiddleware(middlewares.AccessLogsMiddlewareDependencies{
		AppCtx: appCtx,
	})

	jwtValidationMw, err := middlewares.NewJWTValidationMiddleware(middlewares.JWTValidationMiddlewareDependencies{
		AppCtx: appCtx,
	})
	if err != nil {
		appCtx.Logger.Info("failed starting JWT validation middleware", "error", err.Error())
	}

	// 2. Create the proxy
	pxy := proxy.NewMCPProxy(proxy.MCPProxyDependencies{
		AppContext: appCtx,
	})

	// 3. Create the MCP server and client
	err = pxy.InitializeBackend(appCtx.Context)
	if err != nil {
		// TODO
	}

	pxy.McpServer = server.NewMCPServer(
		appCtx.Config.Server.Name,
		appCtx.Config.Server.Version,
		server.WithToolCapabilities(true),
	)

	// 4. Initialize extra handlers for later usage
	hm := handlers.NewHandlersManager(handlers.HandlersManagerDependencies{
		AppCtx: appCtx,
	})

	// 5. Add some useful magic in the form of tools to your MCP server
	// This is the most useful part
	tm := tools.NewToolsManager(tools.ToolsManagerDependencies{
		AppCtx: appCtx,
		Proxy:  pxy,
	})
	tm.AddTools()

	// 6. Wrap MCP server in a transport (stdio, HTTP, SSE)
	switch appCtx.Config.Server.Transport.Type {
	case "http":
		httpServer := server.NewStreamableHTTPServer(pxy.McpServer,
			server.WithHeartbeatInterval(30*time.Second),
			server.WithStateLess(false))

		// Register it under a path, then add custom endpoints.
		// Custom endpoints are needed as the library is not feature-complete according to MCP spec requirements (2025-06-16)
		// Ref: https://modelcontextprotocol.io/specification/2025-06-18/basic/authorization#overview
		mux := http.NewServeMux()
		mux.Handle("/mcp", accessLogsMw.Middleware(jwtValidationMw.Middleware(httpServer)))

		if appCtx.Config.OAuthAuthorizationServer.Enabled {
			mux.Handle("/.well-known/oauth-authorization-server", accessLogsMw.Middleware(http.HandlerFunc(hm.HandleOauthAuthorizationServer)))
		}

		if appCtx.Config.OAuthProtectedResource.Enabled {
			mux.Handle("/.well-known/oauth-protected-resource", accessLogsMw.Middleware(http.HandlerFunc(hm.HandleOauthProtectedResources)))
		}

		// Start StreamableHTTP server
		appCtx.Logger.Info("starting StreamableHTTP server", "host", appCtx.Config.Server.Transport.HTTP.Host)
		err := http.ListenAndServe(appCtx.Config.Server.Transport.HTTP.Host, mux)
		if err != nil {
			log.Fatal(err)
		}

	default:
		// Start stdio server
		appCtx.Logger.Info("starting stdio server")
		if err := server.ServeStdio(pxy.McpServer); err != nil {
			log.Fatal(err)
		}
	}
}
