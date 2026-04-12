<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { oidcService } from '$lib/services/oidc.service';
	import { toastService } from '$lib/services/toast.service.svelte';
	import { authStore } from '$lib/state/auth.svelte';
	import { Spinner } from 'flowbite-svelte';
	import { onMount } from 'svelte';

	let error = $state('');

	onMount(async () => {
		try {
			const user = await oidcService.handleCallback();
			authStore.setUser(user);
			goto(resolve('/dashboard'));
		} catch (e) {
			console.error('OIDC callback error:', e);
			error = 'Authentication failed. Please try again.';
			toastService.error(error);
			setTimeout(() => goto(resolve('/login')), 2000);
		}
	});
</script>

<div class="flex min-h-[60vh] items-center justify-center">
	{#if error}
		<p class="text-red-600 dark:text-red-400">{error}</p>
	{:else}
		<div class="text-center space-y-4">
			<Spinner size="10" />
			<p class="text-gray-500 dark:text-gray-400">Completing sign in...</p>
		</div>
	{/if}
</div>
