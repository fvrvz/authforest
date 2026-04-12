<script lang="ts">
	import { resolve } from '$app/paths';
	import {
		Coffee,
		Database,
		FileCode2,
		Github,
		KeyRound,
		Linkedin,
		Server,
		Shield,
		Users,
	} from 'lucide-svelte';

	const cards = [
		{
			title: 'Docker Setup',
			desc: 'Get AuthForest running with Docker Compose in under 2 minutes.',
			href: '/docs/docker-setup',
			icon: Server,
		},
		{
			title: 'Register an App',
			desc: 'Create OAuth clients and configure redirect URIs, token expiry & more.',
			href: '/docs/register-app',
			icon: FileCode2,
		},
		{
			title: 'Manage Users',
			desc: 'Self-registration, admin user creation, and password resets.',
			href: '/docs/manage-users',
			icon: Users,
		},
		{
			title: 'Manage Roles',
			desc: 'Create roles and assign them to users. Roles appear in JWT claims.',
			href: '/docs/manage-roles',
			icon: Shield,
		},
		{
			title: 'OIDC Integration',
			desc: 'Connect any OIDC-compatible library: JS, React, Python, Go & more.',
			href: '/docs/oidc-integration',
			icon: KeyRound,
		},
		{
			title: 'API Reference',
			desc: 'Complete list of all REST endpoints: Auth, Users, Roles, OAuth, OIDC.',
			href: '/docs/api-reference',
			icon: Database,
		},
	];
</script>

# AuthForest Documentation

A self-hosted **OpenID Connect Identity Provider** built with Go. Manage OAuth 2.0 applications, users, and roles from one place.

<div class="not-prose my-8 grid gap-4 sm:grid-cols-2">
	{#each cards as card (card.href)}
		<a
			href={resolve(card.href)}
			class="group flex items-start gap-4 rounded-xl border border-gray-200 bg-white p-5 shadow-sm transition-all hover:border-primary-300 hover:shadow-md dark:border-gray-700 dark:bg-gray-800 dark:hover:border-primary-700"
		>
			<div class="rounded-lg bg-primary-50 p-2.5 dark:bg-primary-950">
				<card.icon class="size-5 text-primary-600 dark:text-primary-400" />
			</div>
			<div>
				<h3 class="font-semibold text-gray-900 group-hover:text-primary-600 dark:text-white dark:group-hover:text-primary-400">{card.title}</h3>
				<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{card.desc}</p>
			</div>
		</a>
	{/each}
</div>

## Features

- **OIDC Compliant** — Standard OpenID Connect discovery, authorization, token & userinfo endpoints
- **Multi-App** — Register unlimited OAuth clients with per-app token expiry configuration
- **Role-Based Access** — Create roles and assign them to users. Roles are included in JWT claims
- **Self-Hosted** — Single Docker image. Your data stays on your infrastructure
- **RS256 + HS256** — OIDC tokens signed with RS256 (RSA). Internal API tokens with HS256
- **Admin Panel** — Built-in SvelteKit dashboard to manage apps, users, and roles

---

<section class="not-prose border-t border-gray-200 dark:border-gray-700 pt-10 mt-10">
	<div class="text-center">
		<h2 class="text-sm font-semibold uppercase tracking-wider text-gray-400 dark:text-gray-500">Built by</h2>
		<p class="mt-3 text-xl font-bold dark:text-white">Faraaz Khan</p>
		<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Full-stack developer &amp; open-source enthusiast</p>
		<div class="mt-5 flex items-center justify-center gap-4">
			<a href="https://github.com/fvrvz" target="_blank" rel="noopener noreferrer" class="rounded-lg border border-gray-200 p-2.5 text-gray-500 transition-colors hover:border-gray-300 hover:text-gray-900 dark:border-gray-600 dark:text-gray-400 dark:hover:border-gray-500 dark:hover:text-white" aria-label="GitHub">
				<Github class="size-5" />
			</a>
			<a href="https://linkedin.com/in/FaraazKhhan" target="_blank" rel="noopener noreferrer" class="rounded-lg border border-gray-200 p-2.5 text-gray-500 transition-colors hover:border-gray-300 hover:text-gray-900 dark:border-gray-600 dark:text-gray-400 dark:hover:border-gray-500 dark:hover:text-white" aria-label="LinkedIn">
				<Linkedin class="size-5" />
			</a>
			<a href="https://ko-fi.com/fayso" target="_blank" rel="noopener noreferrer" class="inline-flex items-center gap-2 rounded-lg border border-yellow-300 bg-yellow-50 px-4 py-2.5 text-sm font-medium text-yellow-800 transition-colors hover:bg-yellow-100 dark:border-yellow-600 dark:bg-yellow-900/30 dark:text-yellow-300 dark:hover:bg-yellow-900/50">
				<Coffee class="size-4" />
				Buy me a coffee
			</a>
		</div>
		<p class="mt-6 flex items-center justify-center gap-1.5 text-sm text-gray-400 dark:text-gray-500">Powered by <img src="/fayso-logo.svg" alt="FaySo" class="h-5 dark:brightness-200" /></p>
	</div>
</section>
