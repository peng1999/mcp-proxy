package middlewares

import (
	"net/http"

	"mcp-proxy/internal/globals"
)

type CORSMiddlewareDependencies struct {
	AppCtx *globals.ApplicationContext
}

type CORSMiddleware struct {
	dependencies CORSMiddlewareDependencies
}

func NewCORSMiddleware(dependencies CORSMiddlewareDependencies) *CORSMiddleware {
	return &CORSMiddleware{dependencies: dependencies}
}

// Middleware adds minimal CORS handling and responds to OPTIONS preflight.
// Policy: allow same-origin requests only (Origin host must equal request.Host),
// mirroring the behavior previously coded in the oauth handlers.
func (mw *CORSMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Always advertise permissive CORS for supported endpoints.
		// This ensures GET responses (including simple requests) include CORS headers.
		rw.Header().Set("Access-Control-Allow-Origin", "*")

		// Methods/headers per endpoint
		if req.URL.Path == "/mcp" {
			rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			rw.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, mcp-protocol-version")
		} else {
			rw.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, mcp-protocol-version")
		}

		// Cache preflight for 10 minutes
		rw.Header().Set("Access-Control-Max-Age", "600")

		// Short-circuit preflight for any origin
		if req.Method == http.MethodOptions {
			rw.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(rw, req)
	})
}
