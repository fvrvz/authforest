<script lang="ts">
	import { authStore } from '$lib/state/auth.svelte';
	import { DarkMode } from 'flowbite-svelte';
	import type { ClassValue } from 'svelte/elements';
	import UserControls from '../user-controls/UserControls.svelte';

	interface Props {
		class?: ClassValue;
	}

	let { class: className }: Props = $props();
</script>

<header class="flex items-center justify-between px-4 py-2 {className}">
	<section class="ml-auto flex items-center gap-3">
		<DarkMode class="rounded-full" />
		{#if authStore.accessToken}
			<UserControls />
		{/if}
	</section>
</header>
