## Context

The application uses Gin to serve an OIDC authentication flow. The callback route is currently registered at `/oidc/callback`, and the default redirect URL is constructed with this path. The OIDC provider has been reconfigured to use `http://127.0.0.1:8080/auth/callback` as the allowed redirect URI, so the application must be updated to match.

Three files are involved:
- `internal/config/config.go` — builds the default `OIDC_REDIRECT_URL` from host/port
- `internal/router/router.go` — registers the HTTP route for the callback handler
- `internal/middleware/oidc.go` — may reference the callback path when constructing the redirect URL

## Goals / Non-Goals

**Goals:**
- Change the OIDC callback route path from `/oidc/callback` to `/auth/callback`
- Ensure the default redirect URL fallback uses the new path
- Ensure no stale references to the old path remain in the codebase

**Non-Goals:**
- Changing any other route paths (login, logout, userinfo)
- Modifying OIDC provider configuration or credentials
- Introducing a runtime-configurable path

## Decisions

**Single path constant**
Use a simple string replacement in each affected file. The path is already driven by config (`OIDC_REDIRECT_URL`), so no new abstraction is needed. If users supply the env var explicitly, the code is unaffected.

Alternatives considered:
- Extract the path to a constant shared across files — unnecessary complexity for a one-liner change.

## Risks / Trade-offs

- [Stale OIDC provider config] If the OIDC provider (e.g., Keycloak, Google) has not been updated to allow `http://127.0.0.1:8080/auth/callback`, logins will fail with a redirect_uri mismatch. → Mitigation: Update allowed redirect URIs in the provider console before or alongside this code change.
- [Existing sessions] Any in-flight OAuth2 state cookies created before the deploy will reference the old path. → Mitigation: Existing sessions expire naturally; users will need to re-login.
