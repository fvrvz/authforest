---
title: JavaScript / TypeScript
---

# JavaScript / TypeScript Integration

Use [`oidc-client-ts`](https://github.com/authts/oidc-client-ts) — the most popular OIDC client library for JavaScript and TypeScript.

## Install

```bash
npm install oidc-client-ts
```

## Configuration

```typescript
import { UserManager, WebStorageStateStore } from 'oidc-client-ts';

const userManager = new UserManager({
	authority: 'http://localhost:8080',
	client_id: 'your-client-id',
	redirect_uri: 'http://localhost:3000/callback',
	post_logout_redirect_uri: 'http://localhost:3000',
	scope: 'openid profile email',
	response_type: 'code',
	userStore: new WebStorageStateStore({ store: window.localStorage }),
});
```

## Login

Redirect the user to AuthForest's authorization endpoint:

```typescript
await userManager.signinRedirect();
```

## Handle Callback

On your callback page (`/callback`):

```typescript
const user = await userManager.signinRedirectCallback();
console.log('Logged in:', user.profile);
// user.profile contains: sub, email, name, roles
```

## Get Current User

```typescript
const user = await userManager.getUser();
if (user && !user.expired) {
	console.log('Access token:', user.access_token);
	console.log('User claims:', user.profile);
}
```

## Logout

```typescript
await userManager.signoutRedirect();
```

## Making API Calls

Include the access token in your API requests:

```typescript
const user = await userManager.getUser();

const response = await fetch('https://your-api.com/data', {
	headers: {
		Authorization: `Bearer ${user.access_token}`,
	},
});
```

## Silent Token Renewal

`oidc-client-ts` supports automatic silent token renewal using refresh tokens:

```typescript
const userManager = new UserManager({
	// ...other config
	automaticSilentRenew: true,
});

userManager.events.addAccessTokenExpired(() => {
	console.log('Token expired, redirecting to login...');
	userManager.signinRedirect();
});
```
