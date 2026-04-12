<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import Login from '$lib/components/features/auth/login/Login.svelte';
	import { oidcService } from '$lib/services/oidc.service';
	import { authStore } from '$lib/state/auth.svelte';
	import { onMount } from 'svelte';

	onMount(async () => {
		// If already authenticated, go to dashboard
		const user = await oidcService.getUser();
		if (user && !user.expired) {
			authStore.setUser(user);
			goto(resolve('/dashboard'));
		}
	});
</script>

<div class="flex min-h-dvh items-center justify-center bg-gray-50 px-4 dark:bg-gray-900">
	<div class="flex flex-col items-center">
		<Login />
		<a
			href={resolve('/')}
			class="mt-4 text-sm text-gray-500 transition-colors hover:text-primary-600 dark:text-gray-400 dark:hover:text-primary-400"
		>
			&larr; Back to home
		</a>
	</div>
</div>
