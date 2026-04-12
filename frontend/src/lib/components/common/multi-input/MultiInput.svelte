<script lang="ts">
	import { Button, Input } from 'flowbite-svelte';
	import { Plus, Trash2 } from 'lucide-svelte';

	interface Props {
		values: string[];
		placeholder?: string;
		id?: string;
	}

	let { values = $bindable(), placeholder = '', id = '' }: Props = $props();

	function add() {
		values = [...values, ''];
	}

	function remove(index: number) {
		values = values.filter((_, i) => i !== index);
	}

	function update(index: number, value: string) {
		values = values.map((v, i) => (i === index ? value : v));
	}
</script>

<div class="space-y-2">
	{#each values as value, i (i)}
		<div class="flex items-center gap-2">
			<Input
				type="text"
				{placeholder}
				id="{id}-{i}"
				value={value}
				oninput={(e) => update(i, (e.currentTarget as HTMLInputElement).value)}
			/>
			{#if values.length > 1}
				<button
					type="button"
					class="shrink-0 cursor-pointer rounded-lg p-2 text-red-500 transition-colors hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-950/30"
					onclick={() => remove(i)}
				>
					<Trash2 class="size-4" />
				</button>
			{/if}
		</div>
	{/each}
	<Button
		type="button"
		size="xs"
		color="alternative"
		class="cursor-pointer"
		onclick={add}
	>
		<Plus class="mr-1 size-3.5" />
		Add
	</Button>
</div>
