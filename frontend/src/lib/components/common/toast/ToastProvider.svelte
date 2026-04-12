<script lang="ts">
	import { toastService } from '$lib/services/toast.service.svelte';
	import { Toast, ToastContainer } from 'flowbite-svelte';
	import { AlertTriangle, CheckCircle, Info, XCircle } from 'lucide-svelte';
	import { fly } from 'svelte/transition';

	const iconMap = {
		green: CheckCircle,
		red: XCircle,
		yellow: AlertTriangle,
		blue: Info,
	};
</script>

<ToastContainer position="top-right">
	{#each toastService.items as toast (toast.id)}
		<Toast
			color={toast.color}
			dismissable
			transition={fly}
			params={{ x: 200, duration: 300 }}
			class="mb-2 w-72"
			onclose={() => toastService.dismiss(toast.id)}
		>
			{#snippet icon()}
				{@const Icon = iconMap[toast.color]}
				<Icon class="size-5" />
			{/snippet}
			{toast.message}
		</Toast>
	{/each}
</ToastContainer>
