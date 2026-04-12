---
title: Manage Roles
---

# Manage Roles

Roles let you control access in your applications. When users authenticate, their assigned roles are included in the JWT `id_token` as a `roles` claim array.

## Create a Role

1. Go to **Roles** from the sidebar
2. Click **Add Role**
3. Enter the role name (e.g., `admin`, `editor`, `viewer`)
4. Save

## Assign Roles to Users

1. Go to **Users** and select a user
2. In the edit panel, select roles from the available list
3. Save changes

A user can have multiple roles.

## Roles in JWT Claims

When a user authenticates via OIDC, their roles appear in the `id_token`:

```json
{
  "sub": "uuid-of-the-user",
  "email": "john@example.com",
  "name": "John Doe",
  "roles": ["admin", "editor"],
  "iat": 1700000000,
  "exp": 1700003600
}
```

You can use this in your application to implement authorization logic:

```javascript
const claims = parseIdToken(token);
if (claims.roles.includes("admin")) {
  // Allow access to admin features
}
```

## Edit & Delete Roles

From the roles list, you can:

- **Edit** — rename a role
- **Delete** — remove a role. Users who had this role will lose it.

> Deleting a role does not invalidate existing tokens. The change takes effect on the next token issuance.
