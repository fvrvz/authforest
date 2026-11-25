<script lang="ts">
	import { browser } from '$app/environment';
	import { asset, resolve } from '$app/paths';
	import { oidcService } from '$lib/services/oidc.service';
	import { toastService } from '$lib/services/toast.service.svelte';
	import { authStore } from '$lib/state/auth.svelte';
	import { Button, DarkMode } from 'flowbite-svelte';
	import gsap from 'gsap';
	import { ScrollTrigger } from 'gsap/ScrollTrigger';
	import Lenis from 'lenis';
	import {
		ArrowRight,
		BookOpen,
		Check,
		Code2,
		Copy,
		Cpu,
		Database,
		ExternalLink,
		Github,
		KeyRound,
		Layers,
		Lock,
		Server,
		Shield,
		Terminal,
		TreePine
	} from 'lucide-svelte';
	import { onDestroy, onMount } from 'svelte';

	let loggingIn = $state(false);
	let isLoggedIn = $state(false);
	let copiedSnippet = $state<string | null>(null);

	// Interactive Token Studio State
	let tokenRole = $state<'admin' | 'security_architect' | 'developer' | 'viewer'>('security_architect');
	let tokenUsername = $state('faraaz@authforest.local');
	let tokenScopes = $state<string[]>(['openid', 'profile', 'email', 'roles']);

	// Interactive Flow Simulator State
	let isSimulating = $state(false);
	let simulationStep = $state<number>(0);

	// Interactive Workbench State
	let activeSdk = $state<'react' | 'nextjs' | 'go' | 'python' | 'curl'>('nextjs');
	let activeTerminal = $state<'discovery' | 'jwks' | 'health'>('discovery');

	let lenisInstance: Lenis | null = null;

	onMount(() => {
		const handlePageShow = () => {
			loggingIn = false;
		};
		window.addEventListener('pageshow', handlePageShow);

		// Check auth status
		oidcService.getUser().then((user) => {
			if (user && !user.expired) {
				authStore.setUser(user);
				isLoggedIn = true;
			}
		}).catch(() => {
			isLoggedIn = false;
		});

		// Lenis smooth scroll & GSAP Animations (Browser only)
		if (browser) {
			try {
				gsap.registerPlugin(ScrollTrigger);

				lenisInstance = new Lenis({
					duration: 1.2,
					easing: (t: number) => Math.min(1, 1.001 - Math.pow(2, -10 * t)),
					smoothWheel: true
				});

				const raf = (time: number) => {
					lenisInstance?.raf(time);
					requestAnimationFrame(raf);
				};
				requestAnimationFrame(raf);

				lenisInstance.on('scroll', ScrollTrigger.update);
				gsap.ticker.add((time: number) => {
					lenisInstance?.raf(time * 1000);
				});
				gsap.ticker.lagSmoothing(0);

				// Hero Reveal Animation
				const heroTl = gsap.timeline({ defaults: { ease: 'power3.out' } });
				heroTl
					.from('.gsap-hero-badge', { opacity: 0, y: -16, duration: 0.6 })
					.from('.gsap-hero-title', { opacity: 0, y: 24, duration: 0.8 }, '-=0.3')
					.from('.gsap-hero-desc', { opacity: 0, y: 16, duration: 0.7 }, '-=0.4')
					.from('.gsap-hero-cta', { opacity: 0, y: 16, duration: 0.6 }, '-=0.4')
					.from('.gsap-hero-card', { opacity: 0, scale: 0.98, y: 30, duration: 0.9 }, '-=0.3');

				// ScrollTrigger reveals on features
				gsap.utils.toArray<HTMLElement>('.gsap-reveal').forEach((el) => {
					gsap.from(el, {
						scrollTrigger: {
							trigger: el,
							start: 'top 85%',
							toggleActions: 'play none none none'
						},
						opacity: 0,
						y: 30,
						duration: 0.8,
						ease: 'power2.out'
					});
				});
			} catch (e) {
				console.warn('Animation initialization skipped:', e);
			}
		}

		return () => {
			window.removeEventListener('pageshow', handlePageShow);
			lenisInstance?.destroy();
		};
	});

	onDestroy(() => {
		lenisInstance?.destroy();
	});

	async function handleLogin() {
		loggingIn = true;
		try {
			await oidcService.login();
		} catch {
			toastService.error('Failed to initiate OIDC authorization flow. Please try again.');
			loggingIn = false;
		}
	}

	function copyToClipboard(text: string, id: string) {
		navigator.clipboard.writeText(text);
		copiedSnippet = id;
		setTimeout(() => {
			if (copiedSnippet === id) copiedSnippet = null;
		}, 2000);
	}

	function toggleScope(scope: string) {
		if (tokenScopes.includes(scope)) {
			if (tokenScopes.length > 1) {
				tokenScopes = tokenScopes.filter((s) => s !== scope);
			}
		} else {
			tokenScopes = [...tokenScopes, scope];
		}
	}

	function runFlowSimulation() {
		if (isSimulating) return;
		isSimulating = true;
		simulationStep = 1;

		setTimeout(() => {
			simulationStep = 2;
		}, 1200);

		setTimeout(() => {
			simulationStep = 3;
		}, 2400);

		setTimeout(() => {
			simulationStep = 4;
		}, 3600);

		setTimeout(() => {
			isSimulating = false;
		}, 5200);
	}

	// Dynamic interactive token payload based on controls
	let tokenPayload = $derived(
		JSON.stringify(
			{
				iss: 'https://authforest.byfaraaz.com',
				sub: 'usr_89f02c91',
				preferred_username: tokenUsername.split('@')[0],
				email: tokenUsername,
				email_verified: true,
				roles: [tokenRole],
				scope: tokenScopes.join(' '),
				client_id: 'prod_app_client_01',
				iat: 1732541800,
				exp: 1732545400,
				jti: 'jwt_d83e1074a3f5'
			},
			null,
			2
		)
	);

	const tokenHeader = `{
  "alg": "RS256",
  "typ": "JWT",
  "kid": "authforest-2026-rsa-key-1"
}`;

	const discoveryJson = `{
  "issuer": "https://authforest.byfaraaz.com",
  "authorization_endpoint": "https://authforest.byfaraaz.com/auth",
  "token_endpoint": "https://authforest.byfaraaz.com/oauth/token",
  "userinfo_endpoint": "https://authforest.byfaraaz.com/userinfo",
  "jwks_uri": "https://authforest.byfaraaz.com/.well-known/jwks.json",
  "response_types_supported": ["code"],
  "subject_types_supported": ["public"],
  "id_token_signing_alg_values_supported": ["RS256"],
  "scopes_supported": ["openid", "profile", "email", "roles"],
  "token_endpoint_auth_methods_supported": [
    "client_secret_basic",
    "client_secret_post",
    "none"
  ],
  "code_challenge_methods_supported": ["S256"]
}`;

	const jwksJson = `{
  "keys": [
    {
      "kty": "RSA",
      "use": "sig",
      "kid": "authforest-2026-rsa-key-1",
      "alg": "RS256",
      "n": "u1w...[2048-bit modulus]...9kQ",
      "e": "AQAB"
    }
  ]
}`;

	const healthJson = `{
  "status": "pass",
  "version": "2.0.1",
  "engine": "go1.24",
  "database": "postgresql-16 (connected)",
  "uptime_seconds": 184920,
  "jwt_key_rotation": "active (RS256)"
}`;

	const sdkSnippets = {
		nextjs: `// Next.js App Router (Auth.js / NextAuth v5)
` + ['import', 'NextAuth', 'from', '"next-auth";'].join(' ') + `

export const { handlers, auth, signIn, signOut } = NextAuth({
  providers: [
    {
      id: "authforest",
      name: "AuthForest",
      type: "oidc",
      issuer: "https://authforest.byfaraaz.com",
      clientId: process.env.AUTHFOREST_CLIENT_ID,
      clientSecret: process.env.AUTHFOREST_CLIENT_SECRET,
      authorization: { params: { scope: "openid profile email roles" } }
    }
  ]
});`,
		react: `// React SPA with PKCE (oidc-client-ts)
import { UserManager } from "oidc-client-ts";

export const userManager = new UserManager({
  authority: "https://authforest.byfaraaz.com",
  client_id: "spa-client-id",
  redirect_uri: window.location.origin + "/callback",
  response_type: "code",
  scope: "openid profile email roles",
  code_challenge_method: "S256"
});

// Redirect user to AuthForest login
export const login = () => userManager.signinRedirect();`,
		go: `// Go backend: zero-roundtrip offline RS256 token verification
package main

import (
    "context"
    "github.com/coreos/go-oidc/v3/oidc"
)

func verifyToken(ctx context.Context, rawToken string) (*oidc.IDToken, error) {
    provider, err := oidc.NewProvider(ctx, "https://authforest.byfaraaz.com")
    if err != nil {
        return nil, err
    }

    verifier := provider.Verifier(&oidc.Config{ClientID: "api-service"})
    return verifier.Verify(ctx, rawToken)
}`,
		python: `# Python FastAPI with Authlib
from fastapi import FastAPI, Depends, Security
from authlib.integrations.starlette_client import OAuth

app = FastAPI()
oauth = OAuth()
oauth.register(
    name="authforest",
    server_metadata_url="https://authforest.byfaraaz.com/.well-known/openid-configuration",
    client_id="fastapi-client-id",
    client_secret="fastapi-secret",
    client_kwargs={"scope": "openid profile email roles"}
)`,
		curl: `# 1. Initiate PKCE Authorization Code Flow in Browser
https://authforest.byfaraaz.com/auth?
  response_type=code
  &client_id=my-app
  &redirect_uri=https://my-app.com/callback
  &scope=openid%20profile%20email%20roles
  &code_challenge=E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM
  &code_challenge_method=S256

# 2. Exchange authorization code for signed JWT tokens
curl -X POST https://authforest.byfaraaz.com/oauth/token \\
  -H "Content-Type: application/x-www-form-urlencoded" \\
  -d "grant_type=authorization_code" \\
  -d "code=AUTH_CODE_RECEIVED" \\
  -d "redirect_uri=https://my-app.com/callback" \\
  -d "client_id=my-app" \\
  -d "code_verifier=dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"`
	};

	const guides = [
		{
			title: 'Quick Start Guide',
			desc: 'Learn the architectural foundations and integrate your first application in 5 minutes.',
			href: '/docs/getting-started',
			icon: BookOpen,
			tag: '5 min'
		},
		{
			title: 'Docker & Compose Setup',
			desc: 'Deploy AuthForest with PostgreSQL, Redis, and health checks with a single command.',
			href: '/docs/docker-setup',
			icon: Server,
			tag: 'DevOps'
		},
		{
			title: 'Next.js & NextAuth Integration',
			desc: 'Seamless OIDC provider configuration for App Router, server components, and session handling.',
			href: '/docs/oidc-nextauth',
			icon: Layers,
			tag: 'Fullstack'
		},
		{
			title: 'React & SPA (PKCE)',
			desc: 'Single Page App authentication with oidc-client-ts, secure tokens, and silent refresh.',
			href: '/docs/oidc-react',
			icon: Code2,
			tag: 'Frontend'
		},
		{
			title: 'Go Microservice Verification',
			desc: 'Fast, offline RS256 JWT validation against cached JWKS keys with zero DB roundtrips.',
			href: '/docs/oidc-go',
			icon: KeyRound,
			tag: 'Backend'
		},
		{
			title: 'REST API & Specs',
			desc: 'Complete technical reference for OIDC Discovery, Token Exchange, Introspection & Users.',
			href: '/docs/api-reference',
			icon: Database,
			tag: 'OpenAPI'
		}
	];
</script>

<svelte:head>
	<title>AuthForest — Open Source OAuth 2.0 & OpenID Connect Identity Engine</title>
	<meta
		name="description"
		content="AuthForest is a high-performance, developer-first OAuth 2.0 and OpenID Connect (OIDC) identity provider written in Go. Standard PKCE, RS256 JWTs, and fine-grained RBAC."
	/>
</svelte:head>

<div class="relative min-h-dvh overflow-x-hidden bg-[#090b0e] text-gray-100 selection:bg-emerald-500/25 selection:text-emerald-300">
	<!-- Background Technical Grid & Ambient Lighting -->
	<div class="pointer-events-none fixed inset-0 bg-grid-pattern opacity-60"></div>
	<div class="pointer-events-none fixed -top-40 left-1/2 -z-10 h-[500px] w-[800px] -translate-x-1/2 rounded-full bg-gradient-to-b from-emerald-500/10 via-cyan-500/5 to-transparent blur-3xl"></div>

	<!-- Floating Glass Navbar -->
	<div class="fixed top-4 inset-x-0 z-50 flex justify-center px-4 sm:px-6">
		<header class="flex w-full max-w-5xl items-center justify-between rounded-full border border-white/10 bg-gray-950/80 px-4 py-2.5 shadow-2xl backdrop-blur-xl sm:px-6">
			<!-- Brand -->
			<a href={resolve('/')} class="group flex items-center gap-2.5">
				<div class="flex size-8 items-center justify-center rounded-lg bg-emerald-500/15 text-emerald-400 border border-emerald-500/30 transition-transform group-hover:scale-105">
					<TreePine class="size-4.5" strokeWidth={2.2} />
				</div>
				<div class="flex items-center gap-2">
					<span class="text-sm font-bold tracking-tight text-white font-mono">AuthForest</span>
					<span class="rounded-full bg-emerald-500/10 border border-emerald-500/20 px-2 py-0.5 text-2xs font-mono font-medium text-emerald-400">
						OIDC 1.0
					</span>
				</div>
			</a>

			<!-- Center Nav Links -->
			<nav class="hidden items-center gap-6 text-xs font-mono text-gray-300 md:flex">
				<a href="#token-studio" class="transition-colors hover:text-emerald-400">Token Lab</a>
				<a href="#flow-architecture" class="transition-colors hover:text-emerald-400">Protocol Flow</a>
				<a href="#code-matrix" class="transition-colors hover:text-emerald-400">SDKs</a>
				<a href={resolve('/docs')} class="transition-colors hover:text-emerald-400">Documentation</a>
				<a href={resolve('/docs/api-reference')} class="transition-colors hover:text-emerald-400">API</a>
			</nav>

			<!-- Right Actions -->
			<div class="flex items-center gap-2.5">
				<a
					href="https://github.com/fvrvz/authforest"
					target="_blank"
					rel="noreferrer"
					class="flex size-8 items-center justify-center rounded-full border border-white/10 text-gray-300 transition-colors hover:border-emerald-500/40 hover:text-white"
					aria-label="GitHub Repository"
				>
					<Github class="size-4" />
				</a>

				<DarkMode class="rounded-full p-1.5 text-gray-400 hover:text-white" />

				{#if isLoggedIn}
					<Button
						href={resolve('/dashboard')}
						size="xs"
						class="rounded-full bg-emerald-500 font-mono text-xs text-black font-semibold hover:bg-emerald-400 shadow-md shadow-emerald-500/20"
					>
						Console
						<ArrowRight class="ml-1 size-3" />
					</Button>
				{:else}
					<Button
						onclick={handleLogin}
						size="xs"
						class="rounded-full bg-white font-mono text-xs text-black font-semibold hover:bg-gray-200 shadow-md"
						loading={loggingIn}
					>
						Sign In
						<ArrowRight class="ml-1 size-3" />
					</Button>
				{/if}
			</div>
		</header>
	</div>

	<!-- Hero Section -->
	<main class="relative z-10 pt-32 pb-24 lg:pt-40 lg:pb-32">
		<div class="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
			<div class="flex flex-col items-center text-center">
				<!-- Technical Badge -->
				<div class="gsap-hero-badge inline-flex items-center gap-2 rounded-full border border-emerald-500/30 bg-emerald-950/40 px-3.5 py-1 text-2xs font-mono text-emerald-300 backdrop-blur-md shadow-inner">
					<span class="size-1.5 rounded-full bg-emerald-400 animate-ping"></span>
					<span>RFC 6749 • RFC 7636 PKCE • ZERO-ROUNDTRIP RS256 VERIFICATION</span>
				</div>

				<!-- Main Headline -->
				<h1 class="gsap-hero-title mt-6 max-w-4xl text-4xl font-extrabold tracking-tight text-white sm:text-6xl lg:text-7xl font-sans leading-[1.08]">
					The Open Source Identity Engine <br class="hidden sm:inline" />
					<span class="bg-gradient-to-r from-emerald-400 via-teal-300 to-cyan-400 bg-clip-text text-transparent">
						Built for Developers
					</span>
				</h1>

				<!-- Subhead -->
				<p class="gsap-hero-desc mt-6 max-w-2xl text-base text-gray-400 sm:text-lg sm:leading-relaxed font-sans">
					AuthForest is a lightweight, self-hostable OpenID Connect & OAuth 2.1 Identity Provider written in Go.
					Issue signed RS256 tokens, enforce fine-grained RBAC claims, and protect your cloud stack with zero external dependencies.
				</p>

				<!-- Hero CTAs -->
				<div class="gsap-hero-cta mt-8 flex flex-wrap items-center justify-center gap-3.5">
					<Button
						href={resolve('/docs/getting-started')}
						class="rounded-xl bg-emerald-500 px-6 py-3 font-mono text-xs font-semibold text-black transition-all hover:bg-emerald-400 shadow-lg shadow-emerald-500/20"
						size="lg"
					>
						<BookOpen class="mr-2 size-4" />
						Explore Documentation
					</Button>

					{#if isLoggedIn}
						<Button
							href={resolve('/dashboard')}
							outline
							class="rounded-xl border-white/20 bg-white/5 px-6 py-3 font-mono text-xs font-semibold text-white hover:bg-white/10"
							size="lg"
						>
							<Lock class="mr-2 size-4 text-emerald-400" />
							Admin Console
						</Button>
					{:else}
						<Button
							onclick={handleLogin}
							outline
							class="rounded-xl border-white/20 bg-white/5 px-6 py-3 font-mono text-xs font-semibold text-white hover:bg-white/10"
							size="lg"
							loading={loggingIn}
						>
							<Lock class="mr-2 size-4 text-emerald-400" />
							Console Sign In
						</Button>
					{/if}

					<div class="flex items-center gap-2 rounded-xl border border-white/10 bg-gray-900/90 px-4 py-3 font-mono text-xs text-gray-300 shadow-inner">
						<Terminal class="size-4 text-emerald-400" />
						<span>docker compose up -d</span>
						<button
							type="button"
							onclick={() => copyToClipboard('docker compose up -d', 'hero-docker')}
							class="ml-2 cursor-pointer text-gray-400 hover:text-white"
							title="Copy command"
						>
							{#if copiedSnippet === 'hero-docker'}
								<Check class="size-3.5 text-emerald-400" />
							{:else}
								<Copy class="size-3.5" />
							{/if}
						</button>
					</div>
				</div>

				<!-- Architecture Benchmarks Metric Bar -->
				<div class="mt-12 grid w-full max-w-4xl grid-cols-2 gap-3 border-y border-white/10 py-4 font-mono text-xs text-gray-400 sm:grid-cols-4">
					<div class="flex flex-col items-center">
						<span class="text-sm font-bold text-white">&lt;0.8ms</span>
						<span class="text-2xs text-gray-500">Token Issuance Latency</span>
					</div>
					<div class="flex flex-col items-center">
						<span class="text-sm font-bold text-emerald-400">RS256</span>
						<span class="text-2xs text-gray-500">2048-bit Asymmetric Keys</span>
					</div>
					<div class="flex flex-col items-center">
						<span class="text-sm font-bold text-white">RFC 7636</span>
						<span class="text-2xs text-gray-500">PKCE S256 Mandated</span>
					</div>
					<div class="flex flex-col items-center">
						<span class="text-sm font-bold text-cyan-400">0 ms DB</span>
						<span class="text-2xs text-gray-500">Offline API Verification</span>
					</div>
				</div>
			</div>

			<!-- Live Interactive JWT & Claims Studio -->
			<div id="token-studio" class="gsap-hero-card mt-16 overflow-hidden rounded-2xl border border-white/10 bg-gray-950/90 shadow-2xl backdrop-blur-2xl">
				<!-- Studio Header -->
				<div class="flex flex-wrap items-center justify-between border-b border-white/10 bg-gray-900/50 px-5 py-3.5">
					<div class="flex items-center gap-3">
						<div class="flex gap-1.5">
							<span class="size-3 rounded-full bg-red-500/70"></span>
							<span class="size-3 rounded-full bg-yellow-500/70"></span>
							<span class="size-3 rounded-full bg-green-500/70"></span>
						</div>
						<span class="font-mono text-xs font-semibold text-gray-300">Live Cryptographic Claims Studio</span>
					</div>

					<div class="flex items-center gap-2 font-mono text-2xs text-emerald-400">
						<span class="size-2 rounded-full bg-emerald-400 animate-pulse"></span>
						<span>Realtime RS256 JWT Generation</span>
					</div>
				</div>

				<!-- Studio Body -->
				<div class="grid grid-cols-1 lg:grid-cols-12 divide-y lg:divide-y-0 lg:divide-x divide-white/10">
					<!-- Left Control Panel -->
					<div class="p-5 lg:col-span-5 space-y-5 bg-gray-950/40 font-mono">
						<div>
							<label for="token-role-select" class="text-2xs uppercase tracking-wider text-gray-400">Active User Role</label>
							<div id="token-role-select" class="mt-2 grid grid-cols-2 gap-2">
								{#each (['admin', 'security_architect', 'developer', 'viewer'] as const) as role (role)}
									<button
										type="button"
										onclick={() => (tokenRole = role)}
										class="cursor-pointer rounded-lg border px-3 py-2 text-left text-xs transition-all {tokenRole === role ? 'border-emerald-500/60 bg-emerald-500/15 text-emerald-300 font-semibold' : 'border-white/10 bg-gray-900/60 text-gray-400 hover:border-white/20'}"
									>
										{role}
									</button>
								{/each}
							</div>
						</div>

						<div>
							<label for="identity-input" class="text-2xs uppercase tracking-wider text-gray-400">User Identity (Subject)</label>
							<input
								id="identity-input"
								type="text"
								bind:value={tokenUsername}
								class="mt-1.5 w-full rounded-lg border border-white/10 bg-gray-900 px-3 py-1.5 text-xs text-white placeholder-gray-500 focus:border-emerald-500 focus:outline-none"
							/>
						</div>

						<div>
							<label for="scope-badges" class="text-2xs uppercase tracking-wider text-gray-400">Granted Scopes</label>
							<div id="scope-badges" class="mt-2 flex flex-wrap gap-1.5">
								{#each ['openid', 'profile', 'email', 'roles', 'vault:read', 'api:deploy'] as scope (scope)}
									<button
										type="button"
										onclick={() => toggleScope(scope)}
										class="cursor-pointer rounded-md border px-2.5 py-1 text-2xs transition-colors {tokenScopes.includes(scope) ? 'border-cyan-500/50 bg-cyan-500/15 text-cyan-300' : 'border-white/10 bg-gray-900 text-gray-500 hover:text-gray-300'}"
									>
										{scope}
									</button>
								{/each}
							</div>
						</div>

						<div class="rounded-xl border border-emerald-500/20 bg-emerald-500/5 p-3 text-2xs text-emerald-400">
							<div class="flex items-center gap-1.5 font-bold">
								<Shield class="size-3.5 text-emerald-400" />
								<span>Zero-Database Verification</span>
							</div>
							<p class="mt-1 text-gray-400 leading-normal">
								Resource servers download public keys from <code class="text-emerald-300">/.well-known/jwks.json</code> and verify token authenticity locally in &lt;100 microseconds.
							</p>
						</div>
					</div>

					<!-- Right Live Token Inspector -->
					<div class="p-5 lg:col-span-7 bg-gray-900/30 font-mono text-xs space-y-4">
						<div class="flex items-center justify-between text-2xs text-gray-400 border-b border-white/10 pb-2">
							<div class="flex items-center gap-3">
								<span class="text-red-400 font-bold">HEADER</span>
								<span class="text-purple-400 font-bold">PAYLOAD</span>
								<span class="text-cyan-400 font-bold">SIGNATURE</span>
							</div>
							<button
								type="button"
								onclick={() => copyToClipboard(tokenPayload, 'copy-token-payload')}
								class="flex cursor-pointer items-center gap-1 rounded bg-white/10 px-2 py-1 text-gray-300 hover:bg-white/20 transition-colors"
							>
								{#if copiedSnippet === 'copy-token-payload'}
									<Check class="size-3 text-emerald-400" />
									<span>Copied JSON</span>
								{:else}
									<Copy class="size-3" />
									<span>Copy Payload</span>
								{/if}
							</button>
						</div>

						<div class="space-y-3">
							<!-- Header -->
							<div>
								<span class="text-2xs text-gray-500">// JOSE Header (Algorithm & Key ID)</span>
								<pre class="overflow-x-auto rounded-lg bg-gray-950 p-2.5 text-red-300 text-2xs leading-relaxed"><code>{tokenHeader}</code></pre>
							</div>

							<!-- Payload -->
							<div>
								<span class="text-2xs text-gray-500">// Decoded Claims Payload (Auto-generated)</span>
								<pre class="overflow-x-auto rounded-lg bg-gray-950 p-2.5 text-purple-300 text-2xs leading-relaxed max-h-56"><code>{tokenPayload}</code></pre>
							</div>

							<!-- Verified Badge -->
							<div class="flex items-center justify-between rounded-lg border border-emerald-500/30 bg-emerald-950/30 px-3 py-2 text-2xs text-emerald-400">
								<div class="flex items-center gap-2">
									<Check class="size-3.5 text-emerald-400" />
									<span class="font-bold">Cryptographically Validated via AuthForest JWKS</span>
								</div>
								<span class="font-mono text-gray-400">kid: authforest-2026-rsa-key-1</span>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</main>

	<!-- Interactive PKCE Protocol Flow Simulator -->
	<section id="flow-architecture" class="gsap-reveal border-t border-white/10 bg-[#07090c] py-20 lg:py-28">
		<div class="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
			<div class="flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-end">
				<div>
					<span class="font-mono text-xs uppercase tracking-widest text-emerald-400">[RFC 7636 PROTOCOL FLOW]</span>
					<h2 class="mt-2 text-3xl font-extrabold tracking-tight text-white sm:text-4xl font-sans">
						Interactive PKCE Authorization Pipeline
					</h2>
					<p class="mt-2 text-sm text-gray-400 max-w-xl">
						Observe how proof-key-for-code-exchange prevents authorization interception attacks without exposing confidential secrets to client applications.
					</p>
				</div>

				<Button
					onclick={runFlowSimulation}
					class="rounded-xl bg-emerald-500/15 border border-emerald-500/40 font-mono text-xs font-semibold text-emerald-300 hover:bg-emerald-500/25"
					size="sm"
				>
					{#if isSimulating}
						<span class="size-2 rounded-full bg-emerald-400 animate-ping mr-2"></span>
						Simulating Flow...
					{:else}
						<Cpu class="size-3.5 mr-1.5" />
						Simulate OIDC Flow
					{/if}
				</Button>
			</div>

			<!-- Pipeline Cards -->
			<div class="mt-12 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4 font-mono">
				<!-- Step 1 -->
				<div class="relative rounded-2xl border p-5 transition-all {simulationStep === 1 ? 'border-emerald-500 bg-emerald-950/20 shadow-lg shadow-emerald-500/10' : 'border-white/10 bg-gray-950/60'}">
					<div class="flex items-center justify-between">
						<span class="flex size-7 items-center justify-center rounded-lg bg-white/10 text-xs font-bold text-white">01</span>
						<span class="text-2xs text-gray-500">Client / Browser</span>
					</div>
					<h3 class="mt-4 text-sm font-bold text-white">PKCE Challenge</h3>
					<p class="mt-2 text-2xs text-gray-400 leading-normal">
						App generates cryptographically random <code class="text-emerald-400">code_verifier</code> and computes SHA-256 <code class="text-emerald-400">code_challenge</code>.
					</p>
					<div class="mt-4 text-3xs rounded bg-black/60 p-2 text-gray-400">
						S256(verifier) &rarr; challenge
					</div>
				</div>

				<!-- Step 2 -->
				<div class="relative rounded-2xl border p-5 transition-all {simulationStep === 2 ? 'border-emerald-500 bg-emerald-950/20 shadow-lg shadow-emerald-500/10' : 'border-white/10 bg-gray-950/60'}">
					<div class="flex items-center justify-between">
						<span class="flex size-7 items-center justify-center rounded-lg bg-white/10 text-xs font-bold text-white">02</span>
						<span class="text-2xs text-gray-500">AuthForest IdP</span>
					</div>
					<h3 class="mt-4 text-sm font-bold text-white">Consent & Login</h3>
					<p class="mt-2 text-2xs text-gray-400 leading-normal">
						User authenticates via AuthForest <code class="text-emerald-400">/auth</code>. Session validated and single-use authorization code generated.
					</p>
					<div class="mt-4 text-3xs rounded bg-black/60 p-2 text-emerald-400 font-bold">
						302 Redirect &rarr; /callback?code=...
					</div>
				</div>

				<!-- Step 3 -->
				<div class="relative rounded-2xl border p-5 transition-all {simulationStep === 3 ? 'border-emerald-500 bg-emerald-950/20 shadow-lg shadow-emerald-500/10' : 'border-white/10 bg-gray-950/60'}">
					<div class="flex items-center justify-between">
						<span class="flex size-7 items-center justify-center rounded-lg bg-white/10 text-xs font-bold text-white">03</span>
						<span class="text-2xs text-gray-500">Backchannel Token</span>
					</div>
					<h3 class="mt-4 text-sm font-bold text-white">Code Exchange</h3>
					<p class="mt-2 text-2xs text-gray-400 leading-normal">
						Client exchanges code + original <code class="text-emerald-400">code_verifier</code> via POST <code class="text-emerald-400">/oauth/token</code>.
					</p>
					<div class="mt-4 text-3xs rounded bg-black/60 p-2 text-cyan-400">
						200 OK &rarr; Access, ID & Refresh Tokens
					</div>
				</div>

				<!-- Step 4 -->
				<div class="relative rounded-2xl border p-5 transition-all {simulationStep === 4 ? 'border-emerald-500 bg-emerald-950/20 shadow-lg shadow-emerald-500/10' : 'border-white/10 bg-gray-950/60'}">
					<div class="flex items-center justify-between">
						<span class="flex size-7 items-center justify-center rounded-lg bg-white/10 text-xs font-bold text-white">04</span>
						<span class="text-2xs text-gray-500">Resource API</span>
					</div>
					<h3 class="mt-4 text-sm font-bold text-white">Offline Verification</h3>
					<p class="mt-2 text-2xs text-gray-400 leading-normal">
						APIs verify JWT RS256 signature against cached JWKS public keys with zero network overhead.
					</p>
					<div class="mt-4 text-3xs rounded bg-black/60 p-2 text-emerald-300 font-bold">
						Claims Validated in &lt;100&mu;s
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Developer Integration Matrix -->
	<section id="code-matrix" class="gsap-reveal border-t border-white/10 py-20 lg:py-28">
		<div class="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
			<div class="text-center">
				<span class="font-mono text-xs uppercase tracking-widest text-emerald-400">[ZERO-CONFIG COMPATIBILITY]</span>
				<h2 class="mt-2 text-3xl font-extrabold tracking-tight text-white sm:text-4xl font-sans">
					Integrate Any Language in 60 Seconds
				</h2>
				<p class="mx-auto mt-3 max-w-2xl text-sm text-gray-400">
					Because AuthForest strictly adheres to OpenID Connect Core 1.0, you don&apos;t need proprietary SDKs. Use standard libraries.
				</p>
			</div>

			<div class="mt-12 overflow-hidden rounded-2xl border border-white/10 bg-gray-950 shadow-2xl">
				<!-- Workbench Nav Tabs -->
				<div class="flex flex-wrap items-center justify-between border-b border-white/10 bg-gray-900/40 px-4 py-3 sm:px-6">
					<div class="flex flex-wrap gap-2 font-mono text-xs">
						{#each ([
							{ id: 'nextjs', label: 'Next.js (Auth.js)' },
							{ id: 'react', label: 'React / SPA (PKCE)' },
							{ id: 'go', label: 'Go (go-oidc)' },
							{ id: 'python', label: 'Python (Authlib)' },
							{ id: 'curl', label: 'cURL / HTTP' }
						] as const) as item (item.id)}
							<button
								type="button"
								onclick={() => (activeSdk = item.id)}
								class="cursor-pointer rounded-lg px-3 py-1.5 transition-all {activeSdk === item.id ? 'bg-emerald-500/15 text-emerald-400 border border-emerald-500/30 font-semibold' : 'text-gray-400 hover:text-white'}"
							>
								{item.label}
							</button>
						{/each}
					</div>

					<button
						type="button"
						onclick={() => copyToClipboard(sdkSnippets[activeSdk], `sdk-copy-${activeSdk}`)}
						class="flex cursor-pointer items-center gap-1.5 rounded-lg border border-white/10 bg-white/5 px-2.5 py-1 text-xs font-mono text-gray-300 hover:bg-white/10"
					>
						{#if copiedSnippet === `sdk-copy-${activeSdk}`}
							<Check class="size-3.5 text-emerald-400" />
							<span>Copied!</span>
						{:else}
							<Copy class="size-3.5" />
							<span>Copy Code</span>
						{/if}
					</button>
				</div>

				<!-- Code Block -->
				<div class="p-4 sm:p-6 bg-black/60 font-mono text-xs sm:text-sm">
					<pre class="overflow-x-auto leading-relaxed text-gray-300"><code>{sdkSnippets[activeSdk]}</code></pre>
				</div>
			</div>
		</div>
	</section>

	<!-- Live Technical Terminal: Discovery & Health -->
	<section class="gsap-reveal border-t border-white/10 bg-[#07090c] py-20 lg:py-28">
		<div class="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
			<div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-center">
				<div class="lg:col-span-5 space-y-4">
					<span class="font-mono text-xs uppercase tracking-widest text-emerald-400">[AUTOMATIC DISCOVERY]</span>
					<h2 class="text-3xl font-extrabold tracking-tight text-white font-sans">
						Standards-Compliant OIDC Discovery
					</h2>
					<p class="text-sm text-gray-400 leading-relaxed font-sans">
						No manual endpoint configuration required. Client SDKs query AuthForest&apos;s well-known discovery document to automatically configure token endpoints, userinfo routes, and public JWKS verification keys.
					</p>

					<div class="space-y-2 pt-2 font-mono text-xs">
						<div class="flex items-center gap-2 text-gray-300">
							<Check class="size-4 text-emerald-400" />
							<span>RFC 8414 OAuth 2.0 Authorization Server Metadata</span>
						</div>
						<div class="flex items-center gap-2 text-gray-300">
							<Check class="size-4 text-emerald-400" />
							<span>OpenID Connect Discovery 1.0</span>
						</div>
						<div class="flex items-center gap-2 text-gray-300">
							<Check class="size-4 text-emerald-400" />
							<span>High-availability health check endpoint</span>
						</div>
					</div>

					<div class="pt-4">
						<Button
							href={resolve('/docs/api-reference')}
							class="rounded-xl border border-white/20 bg-white/5 font-mono text-xs text-white hover:bg-white/10"
							size="sm"
						>
							<Database class="mr-2 size-3.5 text-emerald-400" />
							View OpenAPI Spec
						</Button>
					</div>
				</div>

				<!-- Terminal Card -->
				<div class="lg:col-span-7 rounded-2xl border border-white/10 bg-gray-950 shadow-2xl overflow-hidden font-mono">
					<div class="flex items-center justify-between border-b border-white/10 bg-gray-900/50 px-4 py-3">
						<div class="flex gap-2 text-xs">
							<button
								type="button"
								onclick={() => (activeTerminal = 'discovery')}
								class="cursor-pointer px-2 py-0.5 rounded {activeTerminal === 'discovery' ? 'bg-emerald-500/20 text-emerald-300 font-bold' : 'text-gray-400'}"
							>
								/.well-known/openid-configuration
							</button>
							<button
								type="button"
								onclick={() => (activeTerminal = 'jwks')}
								class="cursor-pointer px-2 py-0.5 rounded {activeTerminal === 'jwks' ? 'bg-emerald-500/20 text-emerald-300 font-bold' : 'text-gray-400'}"
							>
								/.well-known/jwks.json
							</button>
							<button
								type="button"
								onclick={() => (activeTerminal = 'health')}
								class="cursor-pointer px-2 py-0.5 rounded {activeTerminal === 'health' ? 'bg-emerald-500/20 text-emerald-300 font-bold' : 'text-gray-400'}"
							>
								/health
							</button>
						</div>

						<span class="text-3xs text-gray-500">HTTP/2 200 OK</span>
					</div>

					<div class="p-4 bg-black/80 text-xs text-gray-300 max-h-80 overflow-y-auto">
						{#if activeTerminal === 'discovery'}
							<pre><code>{discoveryJson}</code></pre>
						{:else if activeTerminal === 'jwks'}
							<pre><code>{jwksJson}</code></pre>
						{:else}
							<pre><code>{healthJson}</code></pre>
						{/if}
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Documentation Hub -->
	<section class="gsap-reveal border-t border-white/10 py-20 lg:py-28">
		<div class="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
			<div class="flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-end">
				<div>
					<span class="font-mono text-xs uppercase tracking-widest text-emerald-400">[ENGINEERING DOCUMENTATION]</span>
					<h2 class="mt-2 text-3xl font-extrabold tracking-tight text-white sm:text-4xl font-sans">
						Comprehensive Integration Manuals
					</h2>
					<p class="mt-2 text-sm text-gray-400 max-w-xl font-sans">
						Step-by-step production guides, API specifications, and Docker deployment blueprints.
					</p>
				</div>

				<Button
					href={resolve('/docs')}
					class="rounded-xl bg-emerald-500/15 border border-emerald-500/30 font-mono text-xs text-emerald-300 hover:bg-emerald-500/25"
					size="sm"
				>
					Browse All Docs
					<ArrowRight class="ml-1.5 size-3.5" />
				</Button>
			</div>

			<div class="mt-12 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{#each guides as guide (guide.title)}
					<a
						href={resolve(guide.href as any)}
						class="group flex flex-col justify-between rounded-2xl border border-white/10 bg-gray-950/60 p-6 transition-all hover:border-emerald-500/40 hover:bg-gray-950/90 hover:shadow-xl hover:shadow-emerald-500/5"
					>
						<div>
							<div class="flex items-center justify-between font-mono">
								<div class="flex size-9 items-center justify-center rounded-lg bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 group-hover:scale-105 transition-transform">
									<guide.icon class="size-4.5" />
								</div>
								<span class="text-3xs text-gray-500 border border-white/10 rounded-full px-2 py-0.5">
									{guide.tag}
								</span>
							</div>

							<h3 class="mt-4 text-base font-bold text-white group-hover:text-emerald-300 transition-colors">
								{guide.title}
							</h3>
							<p class="mt-1.5 text-xs text-gray-400 leading-normal">
								{guide.desc}
							</p>
						</div>

						<div class="mt-5 flex items-center gap-1 font-mono text-xs text-emerald-400 font-semibold group-hover:translate-x-1 transition-transform">
							<span>Read manual</span>
							<ArrowRight class="size-3" />
						</div>
					</a>
				{/each}
			</div>
		</div>
	</section>

	<!-- Footer -->
	<footer class="border-t border-white/10 bg-[#060709] py-16 text-xs text-gray-500 font-mono">
		<div class="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
			<div class="grid grid-cols-2 gap-8 md:grid-cols-4">
				<div class="col-span-2 md:col-span-1 space-y-3 font-sans">
					<div class="flex items-center gap-2">
						<TreePine class="size-5 text-emerald-400" />
						<span class="text-base font-bold text-white font-mono">AuthForest</span>
					</div>
					<p class="text-xs text-gray-400 leading-normal">
						Self-hosted, high-performance OAuth 2.0 & OpenID Connect Identity Provider written in Go.
					</p>
					<div class="flex items-center gap-1.5 text-gray-400 text-xs">
						<span>Engineered by</span>
						<img src={asset('/fayso-logo.svg')} alt="FaySo" class="h-4 invert" />
					</div>
				</div>

				<div>
					<h4 class="font-mono text-2xs uppercase tracking-wider text-white">Protocol Core</h4>
					<ul class="mt-3 space-y-2">
						<li><a href="#token-studio" class="hover:text-emerald-400">Claims Studio</a></li>
						<li><a href="#flow-architecture" class="hover:text-emerald-400">PKCE Architecture</a></li>
						<li><a href="#code-matrix" class="hover:text-emerald-400">SDK Matrix</a></li>
						<li><a href={resolve('/login')} class="hover:text-emerald-400">Admin Console</a></li>
					</ul>
				</div>

				<div>
					<h4 class="font-mono text-2xs uppercase tracking-wider text-white">Documentation</h4>
					<ul class="mt-3 space-y-2">
						<li><a href={resolve('/docs/getting-started')} class="hover:text-emerald-400">Quick Start</a></li>
						<li><a href={resolve('/docs/docker-setup')} class="hover:text-emerald-400">Docker Setup</a></li>
						<li><a href={resolve('/docs/register-app')} class="hover:text-emerald-400">Register an App</a></li>
						<li><a href={resolve('/docs/api-reference')} class="hover:text-emerald-400">REST API Specs</a></li>
					</ul>
				</div>

				<div>
					<h4 class="font-mono text-2xs uppercase tracking-wider text-white">Source & Community</h4>
					<ul class="mt-3 space-y-2">
						<li>
							<a
								href="https://github.com/fvrvz/authforest"
								target="_blank"
								rel="noreferrer"
								class="inline-flex items-center gap-1 hover:text-emerald-400"
							>
								GitHub Repo
								<ExternalLink class="size-3" />
							</a>
						</li>
						<li>
							<a
								href="https://github.com/fvrvz/authforest/issues"
								target="_blank"
								rel="noreferrer"
								class="inline-flex items-center gap-1 hover:text-emerald-400"
							>
								Issues & Roadmap
								<ExternalLink class="size-3" />
							</a>
						</li>
						<li>
							<a
								href="https://github.com/fvrvz/authforest/blob/main/LICENSE"
								target="_blank"
								rel="noreferrer"
								class="inline-flex items-center gap-1 hover:text-emerald-400"
							>
								MIT License
								<ExternalLink class="size-3" />
							</a>
						</li>
					</ul>
				</div>
			</div>

			<div class="mt-12 border-t border-white/5 pt-6 text-center text-3xs text-gray-600">
				&copy; {new Date().getFullYear()} AuthForest. Open-source under MIT License.
			</div>
		</div>
	</footer>
</div>
