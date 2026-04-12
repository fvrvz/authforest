<script lang="ts">
	import { resolve } from '$app/paths';
	import type { Pathname } from '$app/types';
	import { OIDC } from '$lib/resources/oidc';
	import { toastService } from '$lib/services/toast.service.svelte';
	import { authStore } from '$lib/state/auth.svelte';
	import type { DashboardStats } from '$lib/types/oidc.type';
	import { Spinner } from 'flowbite-svelte';
	import { AppWindow, Shield, Users } from 'lucide-svelte';
	import { onMount } from 'svelte';

	let loading = $state(true);
	let stats = $state<DashboardStats>({ total_users: 0, total_clients: 0 });

	onMount(async () => {
		const [err, data] = await OIDC.getStats();
		if (err) toastService.error('Failed to load dashboard stats');
		else stats = data;
		loading = false;
	});

	const cards = $derived([
		{
			label: 'Total Users',
			value: stats.total_users,
			icon: Users,
			href: '/users',
			color: 'bg-blue-500',
		},
		{
			label: 'Applications',
			value: stats.total_clients,
			icon: AppWindow,
			href: '/applications',
			color: 'bg-emerald-500',
		},
		{
			label: 'Supported Scopes',
			value: 4,
			icon: Shield,
			href: '/settings',
			color: 'bg-purple-500',
		},
	]);
</script>

<div class="space-y-8">
	<div>
		<h1 class="text-3xl font-bold dark:text-white">
			Welcome back, {authStore.user?.firstName ?? 'Admin'}
		</h1>
		<p class="mt-1 text-gray-500 dark:text-gray-400">
			Here's an overview of your Identity Provider.
		</p>
	</div>

	{#if loading}
		<div class="flex justify-center py-12">
			<Spinner size="10" />
		</div>
	{:else}
		<div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
			{#each cards as card (card.label)}
				<a
					href={resolve(card.href as Pathname)}
					class="group flex items-center gap-4 rounded-xl border border-gray-200 bg-white p-6 shadow-sm transition-shadow hover:shadow-md dark:border-gray-600 dark:bg-gray-800"
				>
					<div class="rounded-lg {card.color} p-3 text-white">
						<card.icon class="size-6" />
					</div>
					<div>
						<p class="text-sm text-gray-500 dark:text-gray-400">
							{card.label}
						</p>
						<p class="text-2xl font-bold dark:text-white">{card.value}</p>
					</div>
				</a>
			{/each}
		</div>

		<div
			class="rounded-xl border border-gray-200 bg-white p-6 dark:border-gray-600 dark:bg-gray-800"
		>
			<h2 class="mb-4 text-xl font-semibold dark:text-white">Quick Actions</h2>
			<div class="flex flex-wrap gap-3">
				<a
					href={resolve('/applications/new')}
					class="rounded-lg bg-primary-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
				>
					Register New Application
				</a>
				<a
					href={resolve('/users')}
					class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-gray-500 dark:text-gray-300 dark:hover:bg-gray-700"
				>
					Manage Users
				</a>
			</div>
		</div>

		<div
			class="rounded-xl border border-gray-200 bg-white p-6 dark:border-gray-600 dark:bg-gray-800"
		>
			<h2 class="mb-4 text-xl font-semibold dark:text-white">
				OIDC Configuration
			</h2>
			<div class="space-y-3 text-sm">
				<div
					class="flex flex-col gap-1 border-b border-gray-100 pb-2 sm:flex-row sm:justify-between dark:border-gray-700"
				>
					<span class="text-gray-500 dark:text-gray-400"
						>Discovery Endpoint</span
					>
					<code class="text-xs text-gray-700 dark:text-gray-300"
						>/.well-known/openid-configuration</code
					>
				</div>
				<div
					class="flex flex-col gap-1 border-b border-gray-100 pb-2 sm:flex-row sm:justify-between dark:border-gray-700"
				>
					<span class="text-gray-500 dark:text-gray-400">JWKS Endpoint</span>
					<code class="text-xs text-gray-700 dark:text-gray-300"
						>/.well-known/jwks.json</code
					>
				</div>
				<div
					class="flex flex-col gap-1 border-b border-gray-100 pb-2 sm:flex-row sm:justify-between dark:border-gray-700"
				>
					<span class="text-gray-500 dark:text-gray-400">Token Signing</span>
					<span class="text-gray-700 dark:text-gray-300">RS256</span>
				</div>
				<div
					class="flex flex-col gap-1 border-b border-gray-100 pb-2 sm:flex-row sm:justify-between dark:border-gray-700"
				>
					<span class="text-gray-500 dark:text-gray-400">Supported Flows</span>
					<span class="text-gray-700 dark:text-gray-300"
						>Authorization Code + PKCE</span
					>
				</div>
				<div class="flex flex-col gap-1 sm:flex-row sm:justify-between">
					<span class="text-gray-500 dark:text-gray-400">Scopes</span>
					<span class="text-gray-700 dark:text-gray-300"
						>openid, profile, email, offline_access</span
					>
				</div>
			</div>
		</div>
	{/if}
</div>
