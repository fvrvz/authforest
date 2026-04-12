<script lang="ts">
	import { resolve } from '$app/paths';
	import { oidcService } from '$lib/services/oidc.service';
	import { toastService } from '$lib/services/toast.service.svelte';
	import { A, Button } from 'flowbite-svelte';
	import { TreePine } from 'lucide-svelte';

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

	<Button onclick={handleLogin} class="w-full cursor-pointer" {loading}>
		Continue with OIDC
	</Button>

	<div class="text-center">
		<A class="text-sm" href={resolve('/register')}>Create an account</A>
	</div>
</div>
