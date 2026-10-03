<script lang="ts">
	import { resolve } from '$app/paths';
	import { oidcService } from '$lib/services/oidc.service';
	import { toastService } from '$lib/services/toast.service.svelte';
	import { A, Button, Spinner } from 'flowbite-svelte';
	import { TreePine } from 'lucide-svelte';
	import { onMount } from 'svelte';

	let loading = $state(true);

	async function handleLogin() {
		loading = true;
		try {
			await oidcService.login();
		} catch (e) {
			console.error('OIDC redirect failed:', e);
			toastService.error('Failed to initiate login. Please try again.');
			loading = false;
		}
	}

	onMount(() => {
		handleLogin();
	});
</script>

<div
	class="m-auto w-full max-w-md space-y-6 rounded-xl border border-gray-200 bg-white p-8 shadow-sm dark:border-gray-700 dark:bg-gray-800"
>
	<div class="flex flex-col items-center gap-2 text-center">
		<TreePine
			class="size-10 text-primary-600 dark:text-primary-400"
			strokeWidth={1.8}
		/>
		<h1 class="text-2xl font-bold dark:text-white">Redirecting to Sign In...</h1>
		<p class="text-sm text-gray-500 dark:text-gray-400">
			Connecting to AuthForest Identity Provider
		</p>
	</div>

	{#if loading}
		<div class="flex justify-center py-2">
			<Spinner size="8" class="text-primary-600 dark:text-primary-400" />
		</div>
	{/if}

	<Button onclick={handleLogin} class="w-full cursor-pointer" {loading}>
		{loading ? 'Redirecting...' : 'Continue with OIDC'}
	</Button>

	<div class="text-center">
		<A class="text-sm text-gray-500 dark:text-gray-400" href={resolve('/register')}>Create an account</A>
	</div>
</div>
