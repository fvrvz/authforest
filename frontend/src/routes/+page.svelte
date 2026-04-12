<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { oidcService } from '$lib/services/oidc.service';
	import { toastService } from '$lib/services/toast.service.svelte';
	import { authStore } from '$lib/state/auth.svelte';
	import { Button, DarkMode } from 'flowbite-svelte';
	import { ArrowRight, BookOpen, TreePine } from 'lucide-svelte';
	import { onMount } from 'svelte';

	let checking = $state(true);
	let loggingIn = $state(false);

	onMount(async () => {
		const user = await oidcService.getUser();
		if (user && !user.expired) {
			authStore.setUser(user);
			goto(resolve('/dashboard'));
			return;
		}
		checking = false;
	});

	async function handleLogin() {
		loggingIn = true;
		try {
			await oidcService.login();
		} catch {
			toastService.error('Failed to initiate login. Please try again.');
			loggingIn = false;
		}
	}
</script>

{#if !checking}
	<div class="flex min-h-dvh flex-col bg-gray-50 dark:bg-gray-900">
		<header class="flex items-center justify-between px-6 py-4">
			<div class="flex items-center gap-2.5">
				<TreePine
					class="size-8 text-primary-600 dark:text-primary-400"
					strokeWidth={1.8}
				/>
				<span class="text-xl font-bold tracking-tight dark:text-white"
					>AuthForest</span
				>
			</div>
			<DarkMode class="rounded-full" />
		</header>

		<main class="flex flex-1 items-center justify-center px-4">
			<div class="max-w-lg space-y-6 text-center">
				<h1 class="text-5xl font-extrabold tracking-tight dark:text-white">
					Identity, simplified.
				</h1>
				<p class="text-lg text-gray-600 dark:text-gray-400">
					Manage your OAuth 2.0 & OpenID Connect applications, users, and roles
					from a single control panel.
				</p>
				<div class="flex justify-center gap-4">
					<Button
						onclick={handleLogin}
						class="cursor-pointer gap-2"
						size="lg"
						loading={loggingIn}
					>
						Sign in <ArrowRight class="size-4" />
					</Button>
					<Button
						href={resolve('/register')}
						outline
						class="cursor-pointer"
						size="lg"
					>
						Create account
					</Button>
				</div>
				<a
					href={resolve('/docs')}
					class="inline-flex items-center gap-1.5 text-sm text-gray-500 transition-colors hover:text-primary-600 dark:text-gray-400 dark:hover:text-primary-400"
				>
					<BookOpen class="size-4" />
					Getting Started Guide
				</a>
			</div>
		</main>

		<footer
			class="flex items-center justify-center gap-1.5 py-4 text-xs text-gray-400 dark:text-gray-500"
		>
			Powered by
			<img src="/fayso-logo.svg" alt="FaySo" class="h-5 dark:brightness-200" />
		</footer>
	</div>
{/if}
