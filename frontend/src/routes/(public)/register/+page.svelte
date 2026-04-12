<script>
	import { resolve } from '$app/paths';
	import Register from '$lib/components/features/auth/register/Register.svelte';
</script>

<div class="flex min-h-dvh items-center justify-center bg-gray-50 px-4 dark:bg-gray-900">
	<div class="flex flex-col items-center">
		<Register />
		<a
			href={resolve('/')}
			class="mt-4 text-sm text-gray-500 transition-colors hover:text-primary-600 dark:text-gray-400 dark:hover:text-primary-400"
		>
			&larr; Back to home
		</a>
	</div>
</div>
