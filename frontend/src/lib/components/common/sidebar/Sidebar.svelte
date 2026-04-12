<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import type { Pathname } from '$app/types';
	import {
		AppWindow,
		ChevronRight,
		LayoutDashboard,
		Settings,
		Shield,
		TreePine,
		UserCircle,
		Users,
	} from 'lucide-svelte';

	let isOpen = $state(false);

	const navItems = [
		{ label: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
		{ label: 'Applications', href: '/applications', icon: AppWindow },
		{ label: 'Users', href: '/users', icon: Users },
		{ label: 'Roles', href: '/roles', icon: Shield },
		{ label: 'Profile', href: '/profile', icon: UserCircle },
		{ label: 'Settings', href: '/settings', icon: Settings },
	];

	function isActive(href: string): boolean {
		const current = page.url?.pathname ?? '';
		const resolved = resolve(href as Pathname);
		if (href === '/dashboard') return current === resolved;
		return current.startsWith(resolved);
	}
</script>

<aside
	class="sticky top-0 z-1 h-dvh bg-gray-900 text-gray-50 transition-all max-sm:absolute max-sm:shadow-lg dark:bg-gray-950 dark:text-gray-200 {isOpen
		? 'w-full sm:w-72'
		: 'w-0 sm:w-14'}"
>
	<button
		type="button"
		onclick={() => (isOpen = !isOpen)}
		class="absolute top-5 left-full -translate-x-1/2 cursor-pointer rounded-full bg-white p-1 shadow-md transition-all hover:bg-gray-50 sm:top-10 dark:bg-gray-700"
	>
		<ChevronRight class="size-4 text-gray-800 {isOpen && 'rotate-180'}" />
	</button>
	<nav class="flex flex-col gap-1 {isOpen ? 'px-3' : 'px-0'}">
		<div
			class="mt-2 mb-4 flex items-center justify-center gap-2 select-none"
		>
			<TreePine
				class="size-6 shrink-0 text-primary-400"
				strokeWidth={1.8}
			/>
			{#if isOpen}
				<span class="text-xl font-bold tracking-tight">AuthForest</span>
			{/if}
		</div>
		{#each navItems as item (item.href)}
			{@const active = isActive(item.href)}
			<a
				href={resolve(item.href as Pathname)}
				class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors
					{active
					? 'bg-primary-600 text-white'
					: 'text-gray-400 hover:bg-gray-800 hover:text-white dark:text-gray-400 dark:hover:bg-gray-800'}
					{isOpen ? '' : 'justify-center sm:px-0'}"
				onclick={() => {
					if (window.innerWidth < 640) isOpen = false;
				}}
			>
				<item.icon class="size-5 shrink-0" />
				{#if isOpen}
					<span>{item.label}</span>
				{/if}
			</a>
		{/each}
	</nav>
</aside>
