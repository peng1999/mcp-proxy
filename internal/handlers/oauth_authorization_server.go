package handlers

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

// HandleOauthAuthorizationServer process requests for endpoint: /.well-known/oauth-authorization-server
func (h *HandlersManager) HandleOauthAuthorizationServer(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		response.Header().Set("Allow", http.MethodGet)
		http.Error(response, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	remoteUrl := h.dependencies.AppCtx.Config.OAuthAuthorizationServer.IssuerUri + "/.well-known/openid-configuration"
	client := &http.Client{Timeout: 10 * time.Second}
	remoteResponse, err := client.Get(remoteUrl)
	if err != nil {
		h.dependencies.AppCtx.Logger.Error("error getting content from /.well-known/openid-configuration", "error", err.Error())
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer remoteResponse.Body.Close()

	if remoteResponse.StatusCode != http.StatusOK {
		h.dependencies.AppCtx.Logger.Error("unexpected status code from /.well-known/openid-configuration", "status", remoteResponse.Status)
		http.Error(response, "Bad Gateway", http.StatusBadGateway)
		return
	}

	remoteResponseBytes, err := io.ReadAll(remoteResponse.Body)
	if err != nil {
		h.dependencies.AppCtx.Logger.Error("error reading bytes from remote response", "error", err.Error())
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Cache-Control", "max-age=3600")
	if origin := request.Header.Get("Origin"); origin != "" {
		if o, err := url.Parse(origin); err == nil && o.Host == request.Host {
			response.Header().Set("Access-Control-Allow-Origin", origin)
			response.Header().Set("Vary", "Origin")
			response.Header().Set("Access-Control-Allow-Methods", http.MethodGet)
			response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}
	}

	_, err = response.Write(remoteResponseBytes)
	if err != nil {
		h.dependencies.AppCtx.Logger.Error("error sending response to client", "error", err.Error())
		return
	}
}
