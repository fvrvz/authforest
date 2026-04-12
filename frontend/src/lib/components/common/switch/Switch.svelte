<script lang="ts">
	import { SizeMap, type SwitchProp } from '$lib/types/switch.type';

	let {
		label,
		isChecked = $bindable(),
		onChange,
		Icon,
		size = 'default',
	}: SwitchProp = $props();
</script>

<label class="relative inline-flex cursor-pointer items-center">
	<input
		type="checkbox"
		bind:checked={isChecked}
		class="peer sr-only"
		onchange={onChange}
	/>
	<div
		class="relative rounded-full bg-gray-200 transition-all peer-checked:bg-blue-600 after:absolute
		       after:top-0.5 after:left-0.5 after:rounded-full
		       hover:border hover:border-red-700 {SizeMap[size].knob}
		       after:border after:border-gray-300 after:bg-white after:transition-all after:content-[''] peer-checked:after:translate-x-full
		       peer-checked:after:border-white dark:bg-gray-700 dark:peer-checked:bg-blue-600"
	></div>
	{#if Icon}
		<div
			class="pointer-events-none absolute top-0.5 left-0.5 flex items-center
			       justify-center transition-transform
			       peer-checked:translate-x-full {SizeMap[size].iconWrapper}"
		>
			<Icon class="text-gray-800 peer-checked:text-white {SizeMap[size].icon}" />
		</div>
	{/if}

	{#if label}
		<span class="ms-3 text-sm font-medium text-gray-900 dark:text-gray-300"
			>{label}</span
		>
	{/if}
</label>
