<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import type { Pathname } from '$app/types';
	import {
	  AppWindow,
	  LayoutDashboard,
	  Settings,
	  Shield,
	  TreePine,
	  UserCircle,
	  Users,
	} from 'lucide-svelte';

	interface Props {
		open?: boolean;
	}

	let { open = $bindable(false) }: Props = $props();

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

<!-- svelte-ignore a11y_no_static_element_interactions -->
<!-- svelte-ignore a11y_click_events_have_key_events -->

<!-- Mobile backdrop -->
<div
	class="fixed inset-0 z-30 bg-black/50 transition-opacity sm:hidden {open
		? 'opacity-100'
		: 'pointer-events-none opacity-0'}"
	onclick={() => (open = false)}
></div>

<!-- Sidebar -->
<aside
	class="fixed top-0 left-0 z-40 h-dvh w-64 bg-gray-900 text-gray-50 transition-transform sm:sticky sm:z-1 sm:translate-x-0 dark:bg-gray-950 dark:text-gray-200 {open
		? 'translate-x-0 sm:w-64'
		: '-translate-x-full sm:translate-x-0 sm:w-14'}"
>
	<nav class="flex flex-col gap-1 pt-5.5 {open ? 'px-3' : 'sm:px-0'}">
		<div
			class="mb-4 flex items-center justify-center gap-2 select-none {open
				? ''
				: 'max-sm:hidden'}"
		>
			<TreePine
				class="size-6 shrink-0 text-primary-400"
				strokeWidth={1.8}
			/>
			{#if open}
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
					: 'text-gray-300 hover:bg-gray-800 hover:text-white dark:text-gray-300 dark:hover:bg-gray-800'}
					{open ? '' : 'max-sm:hidden justify-center sm:px-0'}"
				onclick={() => {
					if (window.innerWidth < 640) open = false;
				}}
			>
				<item.icon class="size-5 shrink-0" />
				{#if open}
					<span>{item.label}</span>
				{/if}
			</a>
		{/each}
	</nav>
</aside>
