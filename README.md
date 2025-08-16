# MCP Proxy

![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/achetronic/mcp-proxy)
![GitHub](https://img.shields.io/github/license/achetronic/mcp-proxy)

![YouTube Channel Subscribers](https://img.shields.io/youtube/channel/subscribers/UCeSb3yfsPNNVr13YsYNvCAw?label=achetronic&link=http%3A%2F%2Fyoutube.com%2Fachetronic)
![GitHub followers](https://img.shields.io/github/followers/achetronic?label=achetronic&link=http%3A%2F%2Fgithub.com%2Fachetronic)
![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/achetronic?style=flat&logo=twitter&link=https%3A%2F%2Ftwitter.com%2Fachetronic)

A proxy for Model Context Protocol (MCP) servers that adds authentication, authorization, and enterprise features to any MCP backend. Transform your local MCP servers into secure, scalable services ready for remote access.

## Motivation

MCP servers are often developed for local use with stdio transport, making them unsuitable for production environments where remote access, authentication, and scalability are required. MCP Proxy bridges this gap by providing a layer in front of any MCP server, enabling enterprise deployment while preserving the original MCP functionality.

## Features


- 🔐 **OAuth RFC 8414 and RFC 9728 compliant**
- Support for `.well-known/oauth-protected-resource` and `.well-known/oauth-authorization-server` endpoints
- Both endpoints are configurable

- 🛡️ **Several JWT validation methods**
- Delegated to external systems like Istio
- Locally validated based on JWKS URI and CEL expressions for claims

- 🔄 **Transport bridging**
- Accept StreamableHTTP and Stdio requests and forward to HTTP or stdio MCP backends
- Enable remote access to local MCP servers

- 📋 Access logs can exclude or redact fields
- 🚀 Production-ready: Included full examples, Dockerfile, Helm Chart and GitHub Actions for CI
- ⚡ Super easy to extend: Production vitamins added to a good juice: [mcp-go](https://github.com/mark3labs/mcp-go)

## Deployment

### Production 🚀
Deploy to Kubernetes using the Helm chart located in the `chart/` directory.

## How to develop your MCP

### Prerequisites
- Go 1.24+

### Modify the code and run

Modify the entire codebase and execute `make run`

**Note:** Default YAML config executing the previous command start the server as an HTTP server forwarding 
to a Stdio server. Other needs? just modify the Makefile to use other YAML provided in examples.

### Configuration Examples

Several configuration examples are available [here](./docs)

## 🤝 Contributing

All contributions are welcome! Whether you're reporting bugs, suggesting features, or submitting code — thank you! Here’s how to get involved:

▸ [Open an issue](https://github.com/achetronic/mcp-proxy/issues/new) to report bugs or request features

▸ [Submit a pull request](https://github.com/achetronic/mcp-proxy/pulls) to contribute improvements


## 📄 License

MCP Proxy is licensed under the [Apache 2.0 License](./LICENSE).