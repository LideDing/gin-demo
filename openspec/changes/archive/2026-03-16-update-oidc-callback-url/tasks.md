## 1. Update Default Redirect URL in Config

- [x] 1.1 In `internal/config/config.go`, change the default `OIDC_REDIRECT_URL` fallback from `http://127.0.0.1:%s/oidc/callback` to `http://127.0.0.1:%s/auth/callback`

## 2. Update Route Registration in Router

- [x] 2.1 In `internal/router/router.go`, move the callback route from the `/oidc` group to a new `/auth` group (or register it directly) so it is served at `/auth/callback`

## 3. Verify No Stale References

- [x] 3.1 Search the codebase for any remaining references to `/oidc/callback` and update or remove them
