## Why

The OIDC callback URL has been updated from `http://127.0.0.1:8080/oidc/callback` to `http://127.0.0.1:8080/auth/callback`. The codebase must be updated to register the new route path and use it as the default redirect URL so that the OIDC flow continues to work correctly.

## What Changes

- Update the OIDC callback route from `/oidc/callback` to `/auth/callback`
- Update the default `OIDC_REDIRECT_URL` fallback value to use `/auth/callback`
- Update the login handler to build the redirect URL with `/auth/callback`
- Update any hardcoded references to the old callback path

## Capabilities

### New Capabilities
<!-- None -->

### Modified Capabilities
- `oidc-auth`: The callback route path changes from `/oidc/callback` to `/auth/callback`, affecting the redirect URL used in the OAuth2 authorization flow.

## Impact

- `internal/config/config.go`: Default redirect URL fallback uses `/oidc/callback` — must change to `/auth/callback`
- `internal/middleware/oidc.go`: Callback path may be referenced when building the redirect URL
- `internal/router/router.go`: Route registration for the OIDC callback handler
- No external API or dependency changes; only internal route path update
