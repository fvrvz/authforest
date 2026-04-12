<script lang="ts">
	import { afterNavigate, goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { type Pathname } from '$app/types';
	import { oidcService } from '$lib/services/oidc.service';
	import { authStore } from '$lib/state/auth.svelte';
	import type { MenuItem } from '$lib/types/common..type';
	import { onDestroy, onMount } from 'svelte';
	import Avatar from '../avatar/Avatar.svelte';

	const menuItems: MenuItem[] = [
		{ label: 'Profile', link: '/profile' },
		{ label: 'Settings', link: '/settings' },
	];

	let isOpen = $state(false);
	let container: HTMLDivElement;

	async function signout() {
		isOpen = false;
		await oidcService.logout();
		authStore.setUser(null);
		goto(resolve('/login'));
	}

	afterNavigate(() => (isOpen = false));

	function handleClickOutside(event: MouseEvent) {
		if (!container.contains(event.target as Node)) {
			isOpen = false;
		}
	}

	onMount(() => {
		document.addEventListener('click', handleClickOutside);
	});

	onDestroy(() => {
		document.removeEventListener('click', handleClickOutside);
	});
</script>

<div class="relative" bind:this={container}>
	<button
		onclick={(e) => {
			e.stopPropagation();
			isOpen = !isOpen;
		}}
		class="cursor-pointer"
	>
		<Avatar initials={authStore.userInitials} />
	</button>

	{#if isOpen}
		<nav
			id="userDropdown"
			class="absolute right-0 z-10 mt-2 w-44 divide-y divide-gray-100 rounded-lg bg-white shadow-sm dark:divide-gray-600 dark:bg-gray-700"
		>
			<section class="px-4 py-3 text-sm text-gray-900 dark:text-white">
				<p>{authStore.user?.fullName}</p>
				<p class="truncate font-medium">{authStore.user?.email}</p>
			</section>
			<ul
				class="py-2 text-sm text-gray-700 dark:text-gray-200"
				aria-labelledby="avatarButton"
			>
				{#each menuItems as item (item.link)}
					<li>
						<a
							href={resolve(item.link as Pathname)}
							class="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white"
						>
							{item.label}
						</a>
					</li>
				{/each}
			</ul>
			<section class="py-1">
				<button
					type="button"
					onclick={signout}
					class="w-full cursor-pointer px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-600 dark:hover:text-white"
				>
					Sign out
				</button>
			</section>
		</nav>
	{/if}
</div>
