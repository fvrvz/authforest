<script lang="ts">
	import { resolve } from '$app/paths';
	import { oidcService } from '$lib/services/oidc.service';
	import { toastService } from '$lib/services/toast.service.svelte';
	import { A, Button } from 'flowbite-svelte';
	import { ArrowRight, TreePine } from 'lucide-svelte';
	import { onMount } from 'svelte';

	let loading = $state(false);

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
		// Reset loading state if user returns via browser Back button (bfcache)
		const handlePageShow = () => {
			loading = false;
		};
		window.addEventListener('pageshow', handlePageShow);
		return () => {
			window.removeEventListener('pageshow', handlePageShow);
		};
	});
</script>

<div
	class="m-auto max-w-lg space-y-5 rounded-xl border border-gray-200 bg-white p-8 shadow-sm dark:border-gray-600 dark:bg-gray-800"
>
	<div class="flex flex-col items-center gap-2">
		<TreePine
			class="size-10 text-primary-600 dark:text-primary-400"
			strokeWidth={1.8}
		/>
		<h1 class="text-2xl font-bold dark:text-white">Sign in to AuthForest</h1>
		<p class="text-sm text-gray-500 dark:text-gray-400">
			Continue with your account
		</p>
	</div>

	<Button onclick={handleLogin} class="w-full cursor-pointer gap-2" {loading}>
		Continue with OIDC <ArrowRight class="size-4" />
	</Button>

	<div class="text-center">
		<span class="text-sm text-gray-500 dark:text-gray-400">
			Don't have an account?
		</span>
		<A class="text-sm" href={resolve('/register')}>Create an account</A>
	</div>
</div>
