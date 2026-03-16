# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go web application demonstrating OIDC (OpenID Connect) authentication with the Gin framework. The module path is `git.woa.com/lideding/gin-tai-login`. Supports multiple OIDC providers (Google, Azure AD, Keycloak, Okta, TAI).

## Build & Run Commands

```bash
# Install/sync dependencies
go mod tidy

# Run the application (requires OIDC env vars)
source .env && go run cmd/main.go

# Build binary
go build -o gin-demo cmd/main.go
```

No test files exist yet. No CI/CD configuration.

## Required Environment Variables

- `OIDC_ISSUER_URL` — OIDC Provider issuer URL
- `OIDC_CLIENT_ID` — Client ID
- `OIDC_CLIENT_SECRET` — Client secret
- `OIDC_REDIRECT_URL` (optional, defaults to `http://127.0.0.1:{PORT}/auth/callback`)
- `OIDC_SCOPES` (optional, comma-separated, defaults to `openid,profile`)
- `PORT` (optional, defaults to `8080`)
- `GIN_MODE` (optional, defaults to `debug`)

## Architecture

```
cmd/main.go              → Entry point: loads config, inits OIDC middleware, starts server with graceful shutdown
internal/
  config/config.go       → Loads all config from environment variables, validates required OIDC fields
  middleware/oidc.go      → Core OIDC logic: provider init, session management, login/callback/logout handlers
  handler/oidc.go         → Thin handler layer that delegates to OIDCMiddleware methods
  handler/health.go       → Simple health check handlers (/hi, /ping)
  router/router.go        → Route registration, splits public vs protected (OIDC-guarded) route groups
  service/                → Empty service layer (placeholder)
```

**Key data flow:** `main.go` creates `OIDCMiddleware` → passes it to `router.SetupRouter()` → router creates `OIDCHandler` wrapping the middleware → registers public routes (`/hi`, `/oidc/login`, `/auth/callback`, `/oidc/logout`) and protected routes (`/ping`, `/oidc/userinfo`) guarded by `RequireOIDC()`.

**Session management:** In-memory `map[string]*OIDCSession` inside `OIDCMiddleware`. Sessions are keyed by random base64 IDs stored in `session_id` cookies. CSRF protection uses `oauth_state` cookies.

**TAI-specific field mapping:** `normalizeUserInfo()` in `middleware/oidc.go` maps TAI's `user_name` field to the standard `username` field, with fallback to `preferred_username` then `sub`.

## Adding Protected Routes

```go
// Single route
r.GET("/protected", oidcMiddleware.RequireOIDC(), yourHandler)

// Route group
protected := r.Group("/api")
protected.Use(oidcMiddleware.RequireOIDC())
```

Access user info in handlers via `c.Get("user_info")` (returns `map[string]interface{}`).
