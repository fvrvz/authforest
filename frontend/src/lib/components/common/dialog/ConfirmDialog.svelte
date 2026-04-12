<script lang="ts">
	import { dialogService } from '$lib/services/dialog.service.svelte';
	import { Button, Modal } from 'flowbite-svelte';
	import { AlertTriangle } from 'lucide-svelte';
	import { slide } from 'svelte/transition';
</script>

<Modal
	form
	bind:open={dialogService.open}
	size="xs"
	transition={slide}
	permanent
	onaction={({ action }) => dialogService.handleAction(action === 'yes')}
>
	<div class="text-center">
		<AlertTriangle
			class="mx-auto mb-4 size-12 text-gray-400 dark:text-gray-200"
		/>
		{#if dialogService.options.title}
			<h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">
				{dialogService.options.title}
			</h3>
		{/if}
		<p class="mb-5 text-sm font-normal text-gray-500 dark:text-gray-400">
			{dialogService.options.message}
		</p>
		<div class="flex justify-center gap-3">
			<Button
				type="submit"
				value="yes"
				color={dialogService.options.confirmColor ?? 'red'}
				class="cursor-pointer"
			>
				{dialogService.options.confirmText ?? "Yes, I'm sure"}
			</Button>
			<Button
				type="submit"
				value="no"
				color="alternative"
				class="cursor-pointer"
			>
				{dialogService.options.cancelText ?? 'No, cancel'}
			</Button>
		</div>
	</div>
</Modal>
