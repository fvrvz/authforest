<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { authStore } from '$lib/state/auth.svelte';
	import { onMount } from 'svelte';

	onMount(() => {
		const username = authStore.user?.username;
		if (username) {
			goto(resolve(`/users/${username}`), { replaceState: true });
		} else {
			goto(resolve('/dashboard'), { replaceState: true });
		}
	});
</script>
