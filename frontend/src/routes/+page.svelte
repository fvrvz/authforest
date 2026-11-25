<script lang="ts">
	import { asset, resolve } from '$app/paths';
	import { oidcService } from '$lib/services/oidc.service';
	import { toastService } from '$lib/services/toast.service.svelte';
	import { authStore } from '$lib/state/auth.svelte';
	import { Button, DarkMode } from 'flowbite-svelte';
	import {
		ArrowRight,
		BookOpen,
		Check,
		Code2,
		Copy,
		Cpu,
		Database,
		ExternalLink,
		FileCode2,
		Github,
		Globe,
		KeyRound,
		Layers,
		Lock,
		Server,
		Shield,
		Terminal,
		TreePine,
		Users
	} from 'lucide-svelte';
	import { onMount } from 'svelte';

	let loggingIn = $state(false);
	let isLoggedIn = $state(false);
	let activeTab = $state<'discovery' | 'quickstart' | 'console'>('discovery');
	let activeSdk = $state<'react' | 'nextjs' | 'go' | 'python' | 'curl'>('react');
	let copiedSnippet = $state<string | null>(null);

	onMount(() => {
		const handlePageShow = () => {
			loggingIn = false;
		};
		window.addEventListener('pageshow', handlePageShow);
		return () => {
			window.removeEventListener('pageshow', handlePageShow);
		};
	});

	onMount(async () => {
		try {
			const user = await oidcService.getUser();
			if (user && !user.expired) {
				authStore.setUser(user);
				isLoggedIn = true;
			}
		} catch {
			isLoggedIn = false;
		}
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

	const discoveryJson = `{
  "issuer": "https://authforest.byfaraaz.com",
  "authorization_endpoint": "https://authforest.byfaraaz.com/auth",
  "token_endpoint": "https://authforest.byfaraaz.com/oauth/token",
  "userinfo_endpoint": "https://authforest.byfaraaz.com/userinfo",
  "jwks_uri": "https://authforest.byfaraaz.com/.well-known/jwks.json",
  "response_types_supported": ["code"],
  "subject_types_supported": ["public"],
  "id_token_signing_alg_values_supported": ["RS256"],
  "scopes_supported": [
    "openid",
    "profile",
    "email",
    "roles"
  ],
  "token_endpoint_auth_methods_supported": [
    "client_secret_basic",
    "client_secret_post",
    "none"
  ],
  "code_challenge_methods_supported": ["S256"]
}`;

	const sdkSnippets = {
		react: `// React SPA with oidc-client-ts
import { UserManager } from "oidc-client-ts";

export const userManager = new UserManager({
  authority: "https://authforest.byfaraaz.com",
  client_id: "your-client-id",
  redirect_uri: window.location.origin + "/callback",
  response_type: "code",
  scope: "openid profile email roles",
  code_challenge_method: "S256"
});

export const signIn = () => userManager.signinRedirect();`,
		nextjs: `// Next.js (Auth.js / NextAuth)
` + ['import', 'NextAuth', 'from', '"next-auth";'].join(' ') + `

export const { handlers, auth } = NextAuth({
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
		go: `// Go backend token verification with coreos/go-oidc
ctx := context.Background()
provider, err := oidc.NewProvider(ctx, "https://authforest.byfaraaz.com")
if err != nil {
    log.Fatalf("OIDC discovery failed: %v", err)
}

verifier := provider.Verifier(&oidc.Config{
    ClientID: "your-client-id",
})

// Validates RS256 signature against AuthForest JWKS automatically
idToken, err := verifier.Verify(ctx, rawToken)`,
		python: `# Python FastAPI with Authlib
from authlib.integrations.starlette_client import OAuth

oauth = OAuth()
oauth.register(
    name="authforest",
    server_metadata_url="https://authforest.byfaraaz.com/.well-known/openid-configuration",
    client_id="your-client-id",
    client_secret="your-client-secret",
    client_kwargs={"scope": "openid profile email roles"}
)`,
		curl: `# 1. Start Authorization Code Flow with PKCE
https://authforest.byfaraaz.com/auth?response_type=code
  &client_id=your-client-id
  &redirect_uri=https://your-app.com/callback
  &scope=openid%20profile%20email%20roles
  &code_challenge=E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM
  &code_challenge_method=S256

# 2. Exchange authorization code for tokens
curl -X POST https://authforest.byfaraaz.com/oauth/token \\
  -H "Content-Type: application/x-www-form-urlencoded" \\
  -d "grant_type=authorization_code" \\
  -d "code=AUTHORIZATION_CODE" \\
  -d "redirect_uri=https://your-app.com/callback" \\
  -d "client_id=your-client-id" \\
  -d "code_verifier=dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"`
	};

	const features = [
		{
			title: 'OAuth 2.0 & OpenID Connect 1.0',
			desc: 'Complete implementation of Authorization Code Flow with PKCE (RFC 7636), Refresh Tokens, ID Tokens, and UserInfo.',
			icon: Shield,
			badge: 'Protocol Core'
		},
		{
			title: 'Zero-Config OIDC Discovery',
			desc: 'Standard /.well-known/openid-configuration endpoint ensures instant plug-and-play compatibility with any OIDC client library.',
			icon: Globe,
			badge: 'Standard'
		},
		{
			title: 'Centralized Application Registry',
			desc: 'Register and manage client applications, configure redirect URIs, toggle PKCE enforcement, and issue secure client credentials.',
			icon: FileCode2,
			badge: 'Console'
		},
		{
			title: 'Fine-Grained RBAC & Claims',
			desc: 'Define custom roles and permissions. Assigned roles are automatically embedded into signed JWT access and ID tokens.',
			icon: Users,
			badge: 'Security'
		},
		{
			title: 'Built with Go & PostgreSQL',
			desc: 'Sub-millisecond token issuance, lightweight memory footprint, single-binary distribution, and battle-tested SQL persistence.',
			icon: Cpu,
			badge: 'High Performance'
		},
		{
			title: 'Developer SDKs & Documentation',
			desc: 'Comprehensive guides, code samples, and production recipes for React, Next.js, Go, Python, and Docker container deployment.',
			icon: BookOpen,
			badge: 'Developer First'
		}
	];

	const guides = [
		{
			title: 'Quick Start Guide',
			desc: 'Learn how AuthForest works and integrate your first app in 5 minutes.',
			href: '/docs/getting-started',
			icon: BookOpen,
			tag: 'Intro'
		},
		{
			title: 'Docker Setup',
			desc: 'Spin up AuthForest with PostgreSQL in Docker Compose with one command.',
			href: '/docs/docker-setup',
			icon: Server,
			tag: 'DevOps'
		},
		{
			title: 'React & SPA Integration',
			desc: 'Connect Single Page Applications using oidc-client-ts with PKCE.',
			href: '/docs/oidc-react',
			icon: Code2,
			tag: 'Frontend'
		},
		{
			title: 'Next.js & NextAuth',
			desc: 'Drop-in OIDC provider setup for NextAuth / Auth.js in Next.js 14/15.',
			href: '/docs/oidc-nextauth',
			icon: Layers,
			tag: 'Fullstack'
		},
		{
			title: 'Go Backend Verification',
			desc: 'Verify JWT access tokens and parse user claims using coreos/go-oidc.',
			href: '/docs/oidc-go',
			icon: KeyRound,
			tag: 'Backend'
		},
		{
			title: 'API Reference',
			desc: 'Full specification of all REST endpoints: Auth, OAuth, OIDC, and Users.',
			href: '/docs/api-reference',
			icon: Database,
			tag: 'Reference'
		}
	];

	const standards = [
		'OpenID Connect Core 1.0',
		'OAuth 2.1 & RFC 6749',
		'RFC 7636 PKCE (S256)',
		'RFC 7519 JSON Web Tokens (JWT)',
		'RFC 7517 JWKS Key Rotation',
		'Role-Based Access Control (RBAC)'
	];
</script>

<svelte:head>
	<title>AuthForest — Modern Open Source OAuth 2.0 & OIDC Identity Provider</title>
	<meta
		name="description"
		content="AuthForest is a lightweight, self-hostable OAuth 2.0 and OpenID Connect (OIDC) identity provider written in Go. Secure your web apps, SPAs, and APIs with standards-compliant authentication."
	/>
</svelte:head>

<div class="min-h-dvh bg-white text-gray-900 transition-colors selection:bg-primary-500/20 selection:text-primary-700 dark:bg-gray-950 dark:text-gray-100 dark:selection:text-primary-300">
	<!-- Top Navigation -->
	<header class="sticky top-0 z-50 border-b border-gray-200/80 bg-white/80 backdrop-blur-md dark:border-gray-800/80 dark:bg-gray-950/80">
		<div class="mx-auto flex max-w-7xl items-center justify-between px-4 py-3.5 sm:px-6 lg:px-8">
			<!-- Logo -->
			<a href={resolve('/')} class="flex items-center gap-2.5 transition-transform hover:scale-[1.02]">
				<div class="flex size-9 items-center justify-center rounded-xl bg-primary-600/10 text-primary-600 shadow-xs dark:bg-primary-500/15 dark:text-primary-400">
					<TreePine class="size-5.5" strokeWidth={2.2} />
				</div>
				<div class="flex items-center gap-2">
					<span class="text-xl font-bold tracking-tight text-gray-900 dark:text-white">AuthForest</span>
					<span class="rounded-full bg-primary-100 px-2 py-0.5 text-xs font-semibold text-primary-700 dark:bg-primary-950/80 dark:text-primary-300">
						v2.0
					</span>
				</div>
			</a>

			<!-- Nav items -->
			<nav class="hidden items-center gap-7 text-sm font-medium text-gray-600 md:flex dark:text-gray-300">
				<a href="#features" class="transition-colors hover:text-primary-600 dark:hover:text-primary-400">Features</a>
				<a href="#architecture" class="transition-colors hover:text-primary-600 dark:hover:text-primary-400">Architecture</a>
				<a href="#discovery" class="transition-colors hover:text-primary-600 dark:hover:text-primary-400">OIDC Discovery</a>
				<a href={resolve('/docs')} class="transition-colors hover:text-primary-600 dark:hover:text-primary-400">Documentation</a>
				<a href={resolve('/docs/api-reference')} class="transition-colors hover:text-primary-600 dark:hover:text-primary-400">API</a>
			</nav>

			<!-- Right actions -->
			<div class="flex items-center gap-3">
				<a
					href="https://github.com/fvrvz/authforest"
					target="_blank"
					rel="noreferrer"
					class="flex size-9 items-center justify-center rounded-lg border border-gray-200 text-gray-700 transition-colors hover:bg-gray-100 dark:border-gray-800 dark:text-gray-300 dark:hover:bg-gray-800"
					aria-label="GitHub Repository"
				>
					<Github class="size-4.5" />
				</a>

				<DarkMode class="rounded-lg p-2" />

				{#if isLoggedIn}
					<Button
						href={resolve('/dashboard')}
						class="gap-1.5 font-medium shadow-xs"
						size="sm"
					>
						Open Console
						<ArrowRight class="size-4" />
					</Button>
				{:else}
					<Button
						onclick={handleLogin}
						class="gap-1.5 font-medium shadow-xs"
						size="sm"
						loading={loggingIn}
					>
						Admin Console
						<ArrowRight class="size-4" />
					</Button>
				{/if}
			</div>
		</div>
	</header>

	<!-- Hero Section -->
	<section class="relative overflow-hidden pt-12 pb-20 lg:pt-20 lg:pb-28">
		<!-- Background Glow Accents -->
		<div class="pointer-events-none absolute -top-40 left-1/2 -z-10 h-96 w-[700px] -translate-x-1/2 rounded-full bg-gradient-to-tr from-primary-500/20 via-sky-500/15 to-transparent blur-3xl dark:from-primary-600/25 dark:via-sky-500/10"></div>
		
		<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			<div class="flex flex-col items-center text-center">
				<!-- Pill Tag -->
				<div class="inline-flex items-center gap-2 rounded-full border border-primary-200 bg-primary-50/70 px-3.5 py-1 text-xs font-semibold text-primary-800 shadow-2xs backdrop-blur-xs dark:border-primary-800/80 dark:bg-primary-950/60 dark:text-primary-300">
					<span class="size-2 rounded-full bg-primary-500 animate-pulse"></span>
					<span>Open Source OAuth 2.0 & OpenID Connect 1.0 Identity Provider</span>
				</div>

				<!-- Main Headline -->
				<h1 class="mt-6 max-w-4xl text-4xl font-extrabold tracking-tight text-gray-900 sm:text-6xl lg:text-7xl dark:text-white">
					Modern, Developer-First <br class="hidden sm:inline" />
					<span class="bg-gradient-to-r from-primary-600 via-indigo-500 to-sky-500 bg-clip-text text-transparent dark:from-primary-400 dark:via-indigo-300 dark:to-sky-400">
						Identity Engine
					</span>
					for Cloud & Edge
				</h1>

				<!-- Subtitle -->
				<p class="mt-6 max-w-2xl text-base text-gray-600 sm:text-lg sm:leading-8 dark:text-gray-300">
					AuthForest is a lightweight, self-hostable OIDC Identity Provider built with Go.
					Secure your web apps, SPAs, and microservices with standards-compliant OAuth 2.1,
					PKCE, fine-grained RBAC, and zero-config discovery.
				</p>

				<!-- Hero CTAs -->
				<div class="mt-8 flex flex-wrap items-center justify-center gap-3.5">
					<Button
						href={resolve('/docs/getting-started')}
						class="gap-2 px-5 py-3 text-sm font-semibold shadow-md shadow-primary-500/20 hover:shadow-primary-500/30"
						size="lg"
					>
						<BookOpen class="size-4.5" />
						Explore Documentation
					</Button>

					{#if isLoggedIn}
						<Button
							href={resolve('/dashboard')}
							outline
							class="gap-2 px-5 py-3 text-sm font-semibold"
							size="lg"
						>
							<Lock class="size-4" />
							Admin Console
						</Button>
					{:else}
						<Button
							onclick={handleLogin}
							outline
							class="gap-2 px-5 py-3 text-sm font-semibold"
							size="lg"
							loading={loggingIn}
						>
							<Lock class="size-4" />
							Sign In / Console
						</Button>
					{/if}

					<Button
						href={resolve('/docs/api-reference')}
						color="alternative"
						class="gap-2 px-5 py-3 text-sm font-semibold"
						size="lg"
					>
						<Database class="size-4" />
						REST API Specs
					</Button>
				</div>

				<!-- Quick Docker CLI Pill -->
				<div class="mt-8 flex items-center gap-2 rounded-xl border border-gray-200 bg-gray-50/90 px-3.5 py-2 text-xs font-mono text-gray-700 shadow-2xs dark:border-gray-800 dark:bg-gray-900/90 dark:text-gray-300">
					<Terminal class="size-3.5 text-primary-600 dark:text-primary-400" />
					<span>docker compose up -d</span>
					<button
						type="button"
						onclick={() => copyToClipboard('docker compose up -d', 'docker-cmd')}
						class="ml-2 flex size-5 cursor-pointer items-center justify-center rounded text-gray-400 hover:text-gray-700 dark:hover:text-gray-200"
						title="Copy command"
					>
						{#if copiedSnippet === 'docker-cmd'}
							<Check class="size-3 text-emerald-500" />
						{:else}
							<Copy class="size-3" />
						{/if}
					</button>
				</div>
			</div>

			<!-- Interactive Showcase Tabbed Card -->
			<div id="discovery" class="mt-14 overflow-hidden rounded-2xl border border-gray-200 bg-gray-900 shadow-2xl dark:border-gray-800 dark:bg-gray-950">
				<!-- Card Header / Tabs -->
				<div class="flex flex-wrap items-center justify-between border-b border-gray-800 bg-gray-900/90 px-4 py-3 sm:px-6">
					<div class="flex items-center gap-2">
						<span class="size-3 rounded-full bg-red-500/80"></span>
						<span class="size-3 rounded-full bg-yellow-500/80"></span>
						<span class="size-3 rounded-full bg-green-500/80"></span>
						<span class="ml-2 text-xs font-mono text-gray-400">authforest.byfaraaz.com</span>
					</div>

					<!-- Tab Switcher -->
					<div class="flex items-center gap-1.5 rounded-lg bg-gray-950/80 p-1">
						<button
							type="button"
							onclick={() => (activeTab = 'discovery')}
							class="cursor-pointer rounded-md px-3 py-1 text-xs font-medium transition-colors {activeTab === 'discovery' ? 'bg-primary-600 text-white' : 'text-gray-400 hover:text-white'}"
						>
							OIDC Discovery
						</button>
						<button
							type="button"
							onclick={() => (activeTab = 'quickstart')}
							class="cursor-pointer rounded-md px-3 py-1 text-xs font-medium transition-colors {activeTab === 'quickstart' ? 'bg-primary-600 text-white' : 'text-gray-400 hover:text-white'}"
						>
							SDK Integration
						</button>
						<button
							type="button"
							onclick={() => (activeTab = 'console')}
							class="cursor-pointer rounded-md px-3 py-1 text-xs font-medium transition-colors {activeTab === 'console' ? 'bg-primary-600 text-white' : 'text-gray-400 hover:text-white'}"
						>
							Console Preview
						</button>
					</div>
				</div>

				<!-- Tab 1: OIDC Discovery -->
				{#if activeTab === 'discovery'}
					<div class="p-4 sm:p-6">
						<div class="mb-3 flex items-center justify-between text-xs text-gray-400">
							<div class="flex items-center gap-2 font-mono">
								<span class="rounded bg-emerald-500/20 px-1.5 py-0.5 font-bold text-emerald-400">GET</span>
								<span>/.well-known/openid-configuration</span>
							</div>
							<button
								type="button"
								onclick={() => copyToClipboard(discoveryJson, 'discovery-json')}
								class="flex cursor-pointer items-center gap-1 rounded bg-gray-800 px-2 py-1 text-gray-300 hover:bg-gray-700"
							>
								{#if copiedSnippet === 'discovery-json'}
									<Check class="size-3 text-emerald-400" />
									<span>Copied</span>
								{:else}
									<Copy class="size-3" />
									<span>Copy JSON</span>
								{/if}
							</button>
						</div>
						<pre class="overflow-x-auto rounded-xl bg-gray-950 p-4 font-mono text-xs text-gray-300 sm:text-sm leading-relaxed"><code>{discoveryJson}</code></pre>
					</div>
				{/if}

				<!-- Tab 2: SDK Quickstart -->
				{#if activeTab === 'quickstart'}
					<div class="p-4 sm:p-6">
						<!-- Language pill selector -->
						<div class="mb-4 flex flex-wrap items-center justify-between gap-2 border-b border-gray-800 pb-3">
							<div class="flex flex-wrap gap-1.5">
								<button
									type="button"
									onclick={() => (activeSdk = 'react')}
									class="cursor-pointer rounded-md px-2.5 py-1 text-xs font-medium {activeSdk === 'react' ? 'bg-sky-500/20 text-sky-400 border border-sky-500/40' : 'text-gray-400 hover:text-white'}"
								>
									React / SPA
								</button>
								<button
									type="button"
									onclick={() => (activeSdk = 'nextjs')}
									class="cursor-pointer rounded-md px-2.5 py-1 text-xs font-medium {activeSdk === 'nextjs' ? 'bg-primary-500/20 text-primary-400 border border-primary-500/40' : 'text-gray-400 hover:text-white'}"
								>
									Next.js (Auth.js)
								</button>
								<button
									type="button"
									onclick={() => (activeSdk = 'go')}
									class="cursor-pointer rounded-md px-2.5 py-1 text-xs font-medium {activeSdk === 'go' ? 'bg-cyan-500/20 text-cyan-400 border border-cyan-500/40' : 'text-gray-400 hover:text-white'}"
								>
									Go
								</button>
								<button
									type="button"
									onclick={() => (activeSdk = 'python')}
									class="cursor-pointer rounded-md px-2.5 py-1 text-xs font-medium {activeSdk === 'python' ? 'bg-amber-500/20 text-amber-400 border border-amber-500/40' : 'text-gray-400 hover:text-white'}"
								>
									Python
								</button>
								<button
									type="button"
									onclick={() => (activeSdk = 'curl')}
									class="cursor-pointer rounded-md px-2.5 py-1 text-xs font-medium {activeSdk === 'curl' ? 'bg-purple-500/20 text-purple-400 border border-purple-500/40' : 'text-gray-400 hover:text-white'}"
								>
									cURL / HTTP
								</button>
							</div>

							<button
								type="button"
								onclick={() => copyToClipboard(sdkSnippets[activeSdk], `sdk-${activeSdk}`)}
								class="flex cursor-pointer items-center gap-1 rounded bg-gray-800 px-2 py-1 text-xs text-gray-300 hover:bg-gray-700"
							>
								{#if copiedSnippet === `sdk-${activeSdk}`}
									<Check class="size-3 text-emerald-400" />
									<span>Copied</span>
								{:else}
									<Copy class="size-3" />
									<span>Copy Code</span>
								{/if}
							</button>
						</div>

						<pre class="overflow-x-auto rounded-xl bg-gray-950 p-4 font-mono text-xs text-gray-200 sm:text-sm leading-relaxed"><code>{sdkSnippets[activeSdk]}</code></pre>
					</div>
				{/if}

				<!-- Tab 3: Console Preview -->
				{#if activeTab === 'console'}
					<div class="p-4 sm:p-6 text-gray-200">
						<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
							<!-- App Card -->
							<div class="rounded-xl border border-gray-800 bg-gray-950 p-4">
								<div class="flex items-center justify-between">
									<span class="text-xs font-semibold text-primary-400">SPA Application</span>
									<span class="rounded bg-emerald-500/20 px-1.5 py-0.5 text-2xs text-emerald-400">Active</span>
								</div>
								<h4 class="mt-2 text-sm font-bold text-white">Frontend Dashboard</h4>
								<div class="mt-2 space-y-1 font-mono text-2xs text-gray-400">
									<div>Client ID: <span class="text-gray-300">authforest_admin_panel</span></div>
									<div>Flow: <span class="text-gray-300">Authorization Code + PKCE</span></div>
									<div>Redirect: <span class="text-gray-300">/callback</span></div>
								</div>
							</div>

							<!-- Roles Card -->
							<div class="rounded-xl border border-gray-800 bg-gray-950 p-4">
								<div class="flex items-center justify-between">
									<span class="text-xs font-semibold text-indigo-400">RBAC Configuration</span>
									<span class="text-2xs text-gray-400">3 Roles</span>
								</div>
								<h4 class="mt-2 text-sm font-bold text-white">Role Mapping & Scopes</h4>
								<div class="mt-2 flex flex-wrap gap-1 font-mono text-2xs">
									<span class="rounded bg-primary-950 px-1.5 py-0.5 text-primary-300 border border-primary-800">admin</span>
									<span class="rounded bg-indigo-950 px-1.5 py-0.5 text-indigo-300 border border-indigo-800">developer</span>
									<span class="rounded bg-gray-800 px-1.5 py-0.5 text-gray-300">user</span>
								</div>
								<p class="mt-2 text-2xs text-gray-400">Embedded in JWT token claims automatically.</p>
							</div>

							<!-- OIDC Keys Card -->
							<div class="rounded-xl border border-gray-800 bg-gray-950 p-4">
								<div class="flex items-center justify-between">
									<span class="text-xs font-semibold text-sky-400">Key Management</span>
									<span class="rounded bg-sky-500/20 px-1.5 py-0.5 text-2xs text-sky-400">RS256</span>
								</div>
								<h4 class="mt-2 text-sm font-bold text-white">JWKS Key Rotation</h4>
								<div class="mt-2 font-mono text-2xs text-gray-400 space-y-1">
									<div>Algorithm: <span class="text-gray-300">RSA 2048-bit</span></div>
									<div>Public Endpoint: <span class="text-gray-300">/.well-known/jwks.json</span></div>
									<div>Auto-rotated & signed by IdP</div>
								</div>
							</div>
						</div>

						<div class="mt-4 flex items-center justify-between rounded-xl bg-gray-950 p-4 border border-gray-800">
							<div class="flex items-center gap-3">
								<div class="flex size-8 items-center justify-center rounded-lg bg-primary-600/20 text-primary-400">
									<Lock class="size-4" />
								</div>
								<div>
									<h5 class="text-xs font-bold text-white">Full-Featured Admin Management</h5>
									<p class="text-2xs text-gray-400">Manage client credentials, users, roles, and security audit logs.</p>
								</div>
							</div>
							{#if isLoggedIn}
								<Button
									href={resolve('/dashboard')}
									size="xs"
									class="gap-1"
								>
									Open Dashboard
									<ArrowRight class="size-3" />
								</Button>
							{:else}
								<Button
									onclick={handleLogin}
									size="xs"
									class="gap-1"
								>
									Open Dashboard
									<ArrowRight class="size-3" />
								</Button>
							{/if}
						</div>
					</div>
				{/if}
			</div>
		</div>
	</section>

	<!-- Protocol Standards Banner -->
	<section class="border-y border-gray-200/80 bg-gray-50/70 py-6 dark:border-gray-800/80 dark:bg-gray-900/40">
		<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			<div class="flex flex-wrap items-center justify-center gap-x-8 gap-y-3 text-xs font-semibold uppercase tracking-wider text-gray-600 dark:text-gray-400">
				{#each standards as item (item)}
					<div class="flex items-center gap-2">
						<span class="size-1.5 rounded-full bg-primary-500"></span>
						<span>{item}</span>
					</div>
				{/each}
			</div>
		</div>
	</section>

	<!-- Core Features Grid -->
	<section id="features" class="py-20 lg:py-28">
		<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			<div class="text-center">
				<h2 class="text-xs font-bold uppercase tracking-widest text-primary-600 dark:text-primary-400">Features</h2>
				<h3 class="mt-2 text-3xl font-extrabold tracking-tight text-gray-900 sm:text-4xl dark:text-white">
					Everything You Need from an Identity Engine
				</h3>
				<p class="mx-auto mt-4 max-w-2xl text-base text-gray-600 dark:text-gray-400">
					Engineered for developer happiness, compliance, and enterprise-level reliability.
				</p>
			</div>

			<div class="mt-14 grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3">
				{#each features as feat (feat.title)}
					<div class="relative flex flex-col justify-between rounded-2xl border border-gray-200 bg-white p-7 shadow-xs transition-all hover:border-primary-400/60 hover:shadow-md dark:border-gray-800 dark:bg-gray-900 dark:hover:border-primary-600/50">
						<div>
							<div class="flex items-center justify-between">
								<div class="flex size-11 items-center justify-center rounded-xl bg-primary-100 text-primary-600 dark:bg-primary-950 dark:text-primary-400">
									<feat.icon class="size-5.5" strokeWidth={1.9} />
								</div>
								<span class="rounded-full bg-gray-100 px-2.5 py-0.5 text-2xs font-semibold text-gray-600 dark:bg-gray-800 dark:text-gray-300">
									{feat.badge}
								</span>
							</div>
							<h4 class="mt-5 text-lg font-bold text-gray-900 dark:text-white">{feat.title}</h4>
							<p class="mt-2 text-sm leading-relaxed text-gray-600 dark:text-gray-400">{feat.desc}</p>
						</div>
					</div>
				{/each}
			</div>
		</div>
	</section>

	<!-- Architecture & Flow Sequence -->
	<section id="architecture" class="border-t border-gray-200 bg-gray-50/50 py-20 lg:py-28 dark:border-gray-800 dark:bg-gray-900/30">
		<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			<div class="text-center">
				<h2 class="text-xs font-bold uppercase tracking-widest text-primary-600 dark:text-primary-400">Architecture</h2>
				<h3 class="mt-2 text-3xl font-extrabold tracking-tight text-gray-900 sm:text-4xl dark:text-white">
					How AuthForest Authentication Works
				</h3>
				<p class="mx-auto mt-4 max-w-2xl text-base text-gray-600 dark:text-gray-400">
					Standard OAuth 2.0 Authorization Code Flow with PKCE ensures secure authentication without exposing client secrets in browser environments.
				</p>
			</div>

			<div class="mt-14 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
				<!-- Step 1 -->
				<div class="rounded-2xl border border-gray-200 bg-white p-6 shadow-xs dark:border-gray-800 dark:bg-gray-900">
					<div class="flex size-9 items-center justify-center rounded-lg bg-primary-600 text-sm font-bold text-white">
						1
					</div>
					<h4 class="mt-4 text-base font-bold text-gray-900 dark:text-white">Client Initiation</h4>
					<p class="mt-2 text-xs leading-relaxed text-gray-600 dark:text-gray-400">
						Client app redirects user to AuthForest <code class="font-mono text-2xs">/auth</code> with PKCE <code class="font-mono text-2xs">code_challenge</code> and requested scopes.
					</p>
				</div>

				<!-- Step 2 -->
				<div class="rounded-2xl border border-gray-200 bg-white p-6 shadow-xs dark:border-gray-800 dark:bg-gray-900">
					<div class="flex size-9 items-center justify-center rounded-lg bg-primary-600 text-sm font-bold text-white">
						2
					</div>
					<h4 class="mt-4 text-base font-bold text-gray-900 dark:text-white">User Authentication</h4>
					<p class="mt-2 text-xs leading-relaxed text-gray-600 dark:text-gray-400">
						AuthForest presents login screen, verifies user credentials, checks role assignments, and records user consent.
					</p>
				</div>

				<!-- Step 3 -->
				<div class="rounded-2xl border border-gray-200 bg-white p-6 shadow-xs dark:border-gray-800 dark:bg-gray-900">
					<div class="flex size-9 items-center justify-center rounded-lg bg-primary-600 text-sm font-bold text-white">
						3
					</div>
					<h4 class="mt-4 text-base font-bold text-gray-900 dark:text-white">Code & Token Exchange</h4>
					<p class="mt-2 text-xs leading-relaxed text-gray-600 dark:text-gray-400">
						AuthForest issues one-time authorization code. Client exchanges code with <code class="font-mono text-2xs">code_verifier</code> for signed JWT Access, ID & Refresh tokens.
					</p>
				</div>

				<!-- Step 4 -->
				<div class="rounded-2xl border border-gray-200 bg-white p-6 shadow-xs dark:border-gray-800 dark:bg-gray-900">
					<div class="flex size-9 items-center justify-center rounded-lg bg-primary-600 text-sm font-bold text-white">
						4
					</div>
					<h4 class="mt-4 text-base font-bold text-gray-900 dark:text-white">Offline Verification</h4>
					<p class="mt-2 text-xs leading-relaxed text-gray-600 dark:text-gray-400">
						Backend APIs verify RS256 token signatures offline against cached public keys from <code class="font-mono text-2xs">/.well-known/jwks.json</code> with zero database round-trips.
					</p>
				</div>
			</div>
		</div>
	</section>

	<!-- Documentation & Guides Section -->
	<section class="py-20 lg:py-28">
		<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			<div class="flex flex-col items-center justify-between gap-4 sm:flex-row">
				<div>
					<h2 class="text-xs font-bold uppercase tracking-widest text-primary-600 dark:text-primary-400">Documentation</h2>
					<h3 class="mt-1 text-2xl font-extrabold tracking-tight text-gray-900 sm:text-3xl dark:text-white">
						Developer Guides & Integration Recipes
					</h3>
				</div>
				<Button
					href={resolve('/docs')}
					class="gap-1.5 self-start sm:self-auto"
					size="sm"
				>
					View All Guides
					<ArrowRight class="size-4" />
				</Button>
			</div>

			<div class="mt-10 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
				{#each guides as guide (guide.title)}
					<a
						href={resolve(guide.href as any)}
						class="group flex flex-col justify-between rounded-2xl border border-gray-200 bg-white p-6 shadow-xs transition-all hover:border-primary-500/50 hover:shadow-md dark:border-gray-800 dark:bg-gray-900 dark:hover:border-primary-500/50"
					>
						<div>
							<div class="flex items-center justify-between">
								<div class="flex size-10 items-center justify-center rounded-lg bg-gray-100 text-gray-700 transition-colors group-hover:bg-primary-100 group-hover:text-primary-600 dark:bg-gray-800 dark:text-gray-300 dark:group-hover:bg-primary-950 dark:group-hover:text-primary-400">
									<guide.icon class="size-5" />
								</div>
								<span class="rounded bg-gray-100 px-2 py-0.5 text-2xs font-semibold text-gray-600 dark:bg-gray-800 dark:text-gray-400">
									{guide.tag}
								</span>
							</div>
							<h4 class="mt-4 text-base font-bold text-gray-900 transition-colors group-hover:text-primary-600 dark:text-white dark:group-hover:text-primary-400">
								{guide.title}
							</h4>
							<p class="mt-1.5 text-xs text-gray-600 dark:text-gray-400">{guide.desc}</p>
						</div>
						<div class="mt-4 flex items-center gap-1 text-xs font-semibold text-primary-600 transition-transform group-hover:translate-x-1 dark:text-primary-400">
							<span>Read guide</span>
							<ArrowRight class="size-3.5" />
						</div>
					</a>
				{/each}
			</div>
		</div>
	</section>

	<!-- Self-Host Callout Banner -->
	<section class="mx-auto max-w-7xl px-4 pb-20 sm:px-6 lg:px-8">
		<div class="relative overflow-hidden rounded-3xl border border-primary-500/30 bg-gradient-to-r from-primary-900 via-indigo-950 to-gray-950 p-8 sm:p-12 text-white shadow-2xl">
			<div class="max-w-2xl">
				<h3 class="text-2xl font-extrabold sm:text-3xl">Self-Host AuthForest in Under 2 Minutes</h3>
				<p class="mt-3 text-sm text-gray-300 sm:text-base">
					AuthForest is distributed as a lightweight container image with zero external runtime dependencies beyond PostgreSQL.
				</p>

				<div class="mt-6 flex flex-wrap gap-3">
					<Button
						href={resolve('/docs/docker-setup')}
						class="gap-1.5 bg-white text-gray-900 hover:bg-gray-100 font-semibold"
						size="sm"
					>
						<Server class="size-4" />
						Docker Setup Instructions
					</Button>
					<Button
						href="https://github.com/fvrvz/authforest"
						target="_blank"
						outline
						class="gap-1.5 text-white border-white/30 hover:bg-white/10"
						size="sm"
					>
						<Github class="size-4" />
						Star on GitHub
					</Button>
				</div>
			</div>
		</div>
	</section>

	<!-- Footer -->
	<footer class="border-t border-gray-200 bg-gray-50 py-12 text-xs text-gray-500 dark:border-gray-800 dark:bg-gray-950 dark:text-gray-400">
		<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			<div class="grid grid-cols-2 gap-8 md:grid-cols-4">
				<div class="col-span-2 md:col-span-1">
					<div class="flex items-center gap-2">
						<TreePine class="size-5 text-primary-600 dark:text-primary-400" />
						<span class="text-base font-bold text-gray-900 dark:text-white">AuthForest</span>
					</div>
					<p class="mt-2 text-xs text-gray-500 dark:text-gray-400">
						Lightweight, self-hosted OAuth 2.0 & OpenID Connect 1.0 Identity Provider.
					</p>
					<div class="mt-4 flex items-center gap-1.5 text-gray-400">
						<span>Powered by</span>
						<img src={asset('/fayso-logo.svg')} alt="FaySo" class="h-4 dark:invert" />
					</div>
				</div>

				<div>
					<h5 class="font-bold text-gray-900 dark:text-white">Product</h5>
					<ul class="mt-3 space-y-2">
						<li><a href="#features" class="hover:text-primary-600 dark:hover:text-primary-400">Features</a></li>
						<li><a href="#architecture" class="hover:text-primary-600 dark:hover:text-primary-400">Architecture</a></li>
						<li><a href="#discovery" class="hover:text-primary-600 dark:hover:text-primary-400">OIDC Discovery</a></li>
						<li><a href={resolve('/login')} class="hover:text-primary-600 dark:hover:text-primary-400">Admin Console</a></li>
					</ul>
				</div>

				<div>
					<h5 class="font-bold text-gray-900 dark:text-white">Documentation</h5>
					<ul class="mt-3 space-y-2">
						<li><a href={resolve('/docs/getting-started')} class="hover:text-primary-600 dark:hover:text-primary-400">Quick Start</a></li>
						<li><a href={resolve('/docs/docker-setup')} class="hover:text-primary-600 dark:hover:text-primary-400">Docker Setup</a></li>
						<li><a href={resolve('/docs/register-app')} class="hover:text-primary-600 dark:hover:text-primary-400">Register an App</a></li>
						<li><a href={resolve('/docs/api-reference')} class="hover:text-primary-600 dark:hover:text-primary-400">API Reference</a></li>
					</ul>
				</div>

				<div>
					<h5 class="font-bold text-gray-900 dark:text-white">Community & Source</h5>
					<ul class="mt-3 space-y-2">
						<li>
							<a
								href="https://github.com/fvrvz/authforest"
								target="_blank"
								rel="noreferrer"
								class="inline-flex items-center gap-1 hover:text-primary-600 dark:hover:text-primary-400"
							>
								GitHub Repository
								<ExternalLink class="size-3" />
							</a>
						</li>
						<li>
							<a
								href="https://github.com/fvrvz/authforest/issues"
								target="_blank"
								rel="noreferrer"
								class="inline-flex items-center gap-1 hover:text-primary-600 dark:hover:text-primary-400"
							>
								Issue Tracker
								<ExternalLink class="size-3" />
							</a>
						</li>
						<li>
							<a
								href="https://github.com/fvrvz/authforest/blob/main/LICENSE"
								target="_blank"
								rel="noreferrer"
								class="inline-flex items-center gap-1 hover:text-primary-600 dark:hover:text-primary-400"
							>
								MIT License
								<ExternalLink class="size-3" />
							</a>
						</li>
					</ul>
				</div>
			</div>

			<div class="mt-10 border-t border-gray-200 pt-6 text-center text-2xs text-gray-400 dark:border-gray-800 dark:text-gray-500">
				&copy; {new Date().getFullYear()} AuthForest. Open-source under MIT License.
			</div>
		</div>
	</footer>
</div>
