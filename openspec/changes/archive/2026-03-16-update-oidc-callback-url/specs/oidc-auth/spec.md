## MODIFIED Requirements

### Requirement: OIDC Callback Route
The system SHALL register the OIDC authorization code callback handler at `/auth/callback`.

#### Scenario: Successful authorization callback
- **WHEN** the OIDC provider redirects the user to `http://127.0.0.1:8080/auth/callback` with a valid `code` and `state` parameter
- **THEN** the application SHALL exchange the code for tokens, create a session, and redirect the user to the protected resource

#### Scenario: Default redirect URL construction
- **WHEN** `OIDC_REDIRECT_URL` environment variable is not set
- **THEN** the application SHALL use `http://127.0.0.1:{PORT}/auth/callback` as the redirect URL sent to the OIDC provider

#### Scenario: Explicit redirect URL override
- **WHEN** `OIDC_REDIRECT_URL` environment variable is set to a custom value
- **THEN** the application SHALL use that custom value unchanged as the redirect URL
