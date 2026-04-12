<script lang="ts">
	import type { AvatarProps } from '$lib/types/avatar.type';
	import { Avatar } from 'flowbite-svelte';
	import { User } from 'lucide-svelte';

	const { initials, image }: AvatarProps = $props();
</script>

<Avatar class="relative" src={image}>
	{#if initials && !image}
		<span class="font-medium text-gray-600 uppercase dark:text-gray-300">
			{initials}
		</span>
	{:else}
		<div
			class="relative size-10 overflow-hidden rounded-full bg-gray-100 dark:bg-gray-600"
		>
			<User class="absolute -left-1 size-12 text-gray-400" />
		</div>
	{/if}
</Avatar>
