---
title: NextAuth / Auth.js
---

# NextAuth / Auth.js Integration

Use the built-in OIDC provider in [NextAuth.js](https://next-auth.js.org/) (or [Auth.js](https://authjs.dev/)) to connect to AuthForest.

## Install

```bash
npm install next-auth
```

## Configuration

### App Router (Next.js 13+)

```typescript
// app/api/auth/[...nextauth]/route.ts
import NextAuth from "next-auth";

const handler = NextAuth({
  providers: [
    {
      id: "authforest",
      name: "AuthForest",
      type: "oidc",
      issuer: "http://localhost:8080",
      clientId: process.env.AUTHFOREST_CLIENT_ID,
      clientSecret: process.env.AUTHFOREST_CLIENT_SECRET,
    },
  ],
  callbacks: {
    async jwt({ token, account, profile }) {
      if (account) {
        token.accessToken = account.access_token;
        token.roles = (profile as any)?.roles ?? [];
      }
      return token;
    },
    async session({ session, token }) {
      session.accessToken = token.accessToken as string;
      session.roles = token.roles as string[];
      return session;
    },
  },
});

export { handler as GET, handler as POST };
```

### Pages Router

```typescript
// pages/api/auth/[...nextauth].ts
import NextAuth from "next-auth";

export default NextAuth({
  providers: [
    {
      id: "authforest",
      name: "AuthForest",
      type: "oidc",
      issuer: "http://localhost:8080",
      clientId: process.env.AUTHFOREST_CLIENT_ID!,
      clientSecret: process.env.AUTHFOREST_CLIENT_SECRET!,
    },
  ],
});
```

## Environment Variables

```ini
AUTHFOREST_CLIENT_ID=your-client-id
AUTHFOREST_CLIENT_SECRET=your-client-secret
```

## Usage in Components

### App Router

```tsx
import { getServerSession } from "next-auth";

export default async function Page() {
  const session = await getServerSession();

  if (!session) {
    return <p>Not signed in</p>;
  }

  return <p>Welcome, {session.user?.name}</p>;
}
```

### Client Component

```tsx
"use client";
import { useSession, signIn, signOut } from "next-auth/react";

export function AuthButton() {
  const { data: session } = useSession();

  if (session) {
    return (
      <>
        <p>Signed in as {session.user?.email}</p>
        <button onClick={() => signOut()}>Sign out</button>
      </>
    );
  }

  return <button onClick={() => signIn("authforest")}>Sign in</button>;
}
```

## Registered Application Settings

When creating the application in AuthForest's admin panel, set:

| Field               | Value                                                |
| ------------------- | ---------------------------------------------------- |
| **Redirect URIs**   | `http://localhost:3000/api/auth/callback/authforest` |
| **Allowed Origins** | `http://localhost:3000`                              |
