<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { authStore } from '$lib/state/auth.svelte';
	import { DarkMode } from 'flowbite-svelte';
	import { ChevronRight } from 'lucide-svelte';
	import type { ClassValue } from 'svelte/elements';
	import UserControls from '../user-controls/UserControls.svelte';

	interface Props {
		class?: ClassValue;
	}

	let { class: className }: Props = $props();

	const labelMap: Record<string, string> = {
		dashboard: 'Dashboard',
		applications: 'Applications',
		new: 'New',
		users: 'Users',
		roles: 'Roles',
		profile: 'Profile',
		settings: 'Settings',
	};

	const breadcrumbs = $derived.by(() => {
		const pathname = page.url?.pathname ?? '';
		const base = resolve('/');
		const relative = base && pathname.startsWith(base) ? pathname.slice(base.length) : pathname;
		const segments = relative.split('/').filter(Boolean);

		return segments.map((seg, i) => {
			const path = '/' + segments.slice(0, i + 1).join('/');
			const href = `${base}${path}`.replace('//', '/');
			const label = labelMap[seg] ?? decodeURIComponent(seg);
			const isLast = i === segments.length - 1;
			return { label, href, isLast };
		});
	});
</script>

<header class="flex items-center justify-between px-4 py-2 {className}">
	<nav class="flex items-center gap-1 text-sm">
		{#each breadcrumbs as crumb (crumb.href)}
			{#if crumb.isLast}
				<span class="font-medium text-gray-900 dark:text-white">{crumb.label}</span>
			{:else}
				<a
					href={crumb.href}
					class="text-gray-500 transition-colors hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
				>
					{crumb.label}
				</a>
				<ChevronRight class="size-3.5 text-gray-400 dark:text-gray-500" />
			{/if}
		{/each}
	</nav>
	<section class="ml-auto flex items-center gap-3">
		<DarkMode class="rounded-full" />
		{#if authStore.accessToken}
			<UserControls />
		{/if}
	</section>
</header>
