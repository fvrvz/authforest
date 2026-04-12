<script lang="ts">
	import type { InputProps } from '$lib/types/input.type';

	let { label, value = $bindable(), ...rest }: InputProps = $props();
</script>

<div class="inline-flex items-center justify-between gap-3">
	<label for={rest.name ?? rest.id}>{label}</label>
	<input
		bind:value
		{...rest}
		class="rounded border-2 border-blue-400 hover:border-blue-500"
	/>
</div>
