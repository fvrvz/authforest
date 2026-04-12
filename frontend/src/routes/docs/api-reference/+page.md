---
title: API Reference
---

# API Reference

Full list of all REST endpoints exposed by the AuthForest backend.

## Authentication

| Method | Endpoint                       | Description                    |
| ------ | ------------------------------ | ------------------------------ |
| `POST` | `/api/v1/auth/login`           | Login with email & password    |
| `POST` | `/api/v1/auth/refresh-token`   | Refresh an access token        |
| `POST` | `/api/v1/auth/forgot-password` | Request a password reset email |
| `POST` | `/api/v1/auth/reset-password`  | Reset password with token      |

## Users

| Method   | Endpoint                           | Description                       |
| -------- | ---------------------------------- | --------------------------------- |
| `POST`   | `/api/v1/users/register`           | Self-register a new user          |
| `GET`    | `/api/v1/users`                    | List all users (admin)            |
| `GET`    | `/api/v1/users/profile`            | Get current user profile          |
| `GET`    | `/api/v1/users/:id`                | Get a user by ID (admin)          |
| `PUT`    | `/api/v1/users/:id`                | Update a user                     |
| `DELETE` | `/api/v1/users/:id`                | Delete a user (admin)             |
| `POST`   | `/api/v1/users/admin-create`       | Admin-create a user (admin)       |
| `PUT`    | `/api/v1/users/:id/reset-password` | Admin-reset user password (admin) |

## Roles

| Method   | Endpoint            | Description           |
| -------- | ------------------- | --------------------- |
| `GET`    | `/api/v1/roles`     | List all roles        |
| `POST`   | `/api/v1/roles`     | Create a role (admin) |
| `GET`    | `/api/v1/roles/:id` | Get a role by ID      |
| `PUT`    | `/api/v1/roles/:id` | Update a role (admin) |
| `DELETE` | `/api/v1/roles/:id` | Delete a role (admin) |

## OAuth 2.0 Clients

| Method   | Endpoint                     | Description                 |
| -------- | ---------------------------- | --------------------------- |
| `POST`   | `/api/v1/oauth2/register`    | Register a new OAuth client |
| `GET`    | `/api/v1/oauth2/clients`     | List all OAuth clients      |
| `GET`    | `/api/v1/oauth2/clients/:id` | Get an OAuth client by ID   |
| `PUT`    | `/api/v1/oauth2/clients/:id` | Update an OAuth client      |
| `DELETE` | `/api/v1/oauth2/clients/:id` | Delete an OAuth client      |

## OpenID Connect

| Method | Endpoint                            | Description                    |
| ------ | ----------------------------------- | ------------------------------ |
| `GET`  | `/.well-known/openid-configuration` | OIDC Discovery document        |
| `GET`  | `/.well-known/jwks.json`            | JSON Web Key Set (public keys) |
| `GET`  | `/oauth2/authorize`                 | Authorization endpoint         |
| `POST` | `/oauth2/token`                     | Token endpoint                 |
| `GET`  | `/oauth2/userinfo`                  | UserInfo endpoint              |

## Authentication Headers

Most API endpoints require authentication. Include the access token:

```
Authorization: Bearer <access_token>
```

Endpoints marked **(admin)** require the user to have an admin role.

## Error Responses

All errors follow a consistent format:

```json
{
	"error": "error_code",
	"message": "Human-readable error description"
}
```

Common HTTP status codes:

| Status | Meaning                                 |
| ------ | --------------------------------------- |
| `400`  | Bad Request — invalid input             |
| `401`  | Unauthorized — missing or invalid token |
| `403`  | Forbidden — insufficient permissions    |
| `404`  | Not Found                               |
| `409`  | Conflict — duplicate resource           |
| `500`  | Internal Server Error                   |
