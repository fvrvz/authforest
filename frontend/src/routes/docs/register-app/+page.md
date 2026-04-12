---
title: Register an Application
---

# Register an Application

Applications (OAuth clients) are registered through the admin panel. Each application gets a unique **Client ID** and **Client Secret** for OIDC flows.

## Create a New Application

1. Sign in to the admin panel
2. Go to **Applications** from the sidebar
3. Click **New Application**
4. Fill in the form:

| Field                | Description                                                |
| -------------------- | ---------------------------------------------------------- |
| **Application Name** | A friendly name (e.g., "My SPA")                           |
| **Redirect URIs**    | Comma-separated list of allowed callback URLs              |
| **Allowed Origins**  | Origins permitted for CORS (e.g., `http://localhost:3000`) |

5. Click **Save**

## Client Credentials

After creation, you'll see:

- **Client ID** — A unique identifier. Safe to include in frontend code.
- **Client Secret** — Used for confidential clients (backend apps). Keep this secure.

> For single-page apps (SPAs), use Authorization Code + PKCE. You won't need the client secret.

## Edit an Application

From the applications list, click on any app to edit its:

- Name
- Redirect URIs
- Allowed Origins
- **Access Token Expiry** — How long access tokens are valid (in seconds)
- **Refresh Token Expiry** — How long refresh tokens are valid (in seconds)

## Token Expiry Defaults

| Token Type    | Default           |
| ------------- | ----------------- |
| Access Token  | 3600s (1 hour)    |
| Refresh Token | 86400s (24 hours) |

You can configure these per application for fine-grained control.

## Delete an Application

Click the delete button on the application detail page. This action is irreversible and will invalidate all tokens issued for this client.
