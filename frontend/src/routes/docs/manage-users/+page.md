---
title: Manage Users
---

# Manage Users

AuthForest supports multiple ways to create and manage users.

## User Registration

Users can self-register via the **Create Account** page or the registration API:

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

## Admin User Creation

Admins can create users directly from the admin panel:

1. Go to **Users** from the sidebar
2. Click **Add User**
3. Fill in the user's name, email, and password
4. Optionally assign roles during creation

Admin-created users can sign in immediately.

## Edit a User

From the users list, click on a user to edit:

- **First Name** and **Last Name**
- **Email Address**
- **Roles** — Assign or remove roles

> Admins can edit all users. Regular users can only edit their own profile.

## Password Reset

### Via Email

If SMTP is configured, users can request a password reset via email:

```bash
curl -X POST http://localhost:8080/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com"}'
```

The user will receive an email with a reset link.

### Admin Reset

Admins can reset any user's password directly from the user edit page.

## User Roles

Each user can have multiple roles. Roles are included in the JWT `id_token` as a `roles` claim. See [Manage Roles](/docs/manage-roles) for details.

## Delete a User

Click the delete button on the user detail page. This is irreversible.
