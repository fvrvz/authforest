---
title: React
---

# React Integration

Use [`react-oidc-context`](https://github.com/authts/react-oidc-context) — a React wrapper around `oidc-client-ts`.

## Install

```bash
npm install oidc-client-ts react-oidc-context
```

## Setup Provider

Wrap your app with the `AuthProvider`:

```tsx
// src/main.tsx
import { AuthProvider } from "react-oidc-context";

const oidcConfig = {
  authority: "http://localhost:8080",
  client_id: "your-client-id",
  redirect_uri: "http://localhost:3000/callback",
  scope: "openid profile email",
  response_type: "code",
};

ReactDOM.createRoot(document.getElementById("root")!).render(
  <AuthProvider {...oidcConfig}>
    <App />
  </AuthProvider>,
);
```

## Login & Logout

```tsx
import { useAuth } from "react-oidc-context";

function LoginButton() {
  const auth = useAuth();

  if (auth.isAuthenticated) {
    return (
      <div>
        <p>Welcome, {auth.user?.profile.name}</p>
        <button onClick={() => auth.removeUser()}>Log out</button>
      </div>
    );
  }

  return <button onClick={() => auth.signinRedirect()}>Log in</button>;
}
```

## Handle Callback

`react-oidc-context` handles the callback automatically. Just make sure your `/callback` route renders a component inside the `AuthProvider`.

```tsx
function CallbackPage() {
  const auth = useAuth();

  if (auth.isLoading) return <p>Loading...</p>;
  if (auth.error) return <p>Error: {auth.error.message}</p>;
  if (auth.isAuthenticated) return <Navigate to="/dashboard" />;

  return null;
}
```

## Access User Info

```tsx
const auth = useAuth();

// User claims from the id_token
const email = auth.user?.profile.email;
const roles = auth.user?.profile.roles;

// Access token for API calls
const token = auth.user?.access_token;
```

## Protected Routes

```tsx
function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const auth = useAuth();

  if (auth.isLoading) return <p>Loading...</p>;
  if (!auth.isAuthenticated) {
    auth.signinRedirect();
    return null;
  }

  return <>{children}</>;
}
```
