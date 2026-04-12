<script lang="ts" generics="T = Record<string, unknown>">
	import { type TableHeader } from '$lib/types/table.type';
	import {
		Checkbox,
		Spinner,
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
	} from 'flowbite-svelte';
	import type { Snippet } from 'svelte';
	import { SvelteMap } from 'svelte/reactivity';

	interface Props {
		data: T[];
		trackBy: keyof T;
		headers: TableHeader<T>[];
		loading?: boolean;
		canSelect?: boolean;
		children?: Snippet<[{ row: T; index: number }]>;
	}

	let {
		data,
		trackBy,
		headers,
		loading = false,
		canSelect,
		children,
	}: Props = $props();

	const selected = new SvelteMap<string, T>();

	const selectAll = $derived.by<0 | 1 | 2>(() => {
		const count = selected.size;
		const total = data.length;

		if (!count) return 0;
		if (count === total) return 1;
		return 2;
	});

	function getKey(row: T, index: number) {
		return String(row[trackBy] ?? index);
	}

	function onSelectAll(e: Event & { currentTarget: HTMLInputElement }) {
		if (!e.currentTarget.checked) {
			selected.clear();
			return;
		}

		data.forEach((row, index) => {
			selected.set(getKey(row, index), row);
		});
	}

	function onSingleSelect(
		e: Event & { currentTarget: HTMLInputElement },
		row: T,
		index: number
	) {
		const key = getKey(row, index);

		if (e.currentTarget.checked) {
			selected.set(key, row);
		} else {
			selected.delete(key);
		}
	}
</script>

<div class="overflow-x-auto">
	<Table hoverable>
		<TableHead>
			{#if canSelect}
				<TableHeadCell>
					<Checkbox
						checked={selectAll === 1}
						indeterminate={selectAll === 2}
						onchange={(e) => onSelectAll(e)}
						disabled={loading}
					/>
				</TableHeadCell>
			{/if}
			{#each headers as header, i (String(header.key) + i)}
				<TableHeadCell class="normal-case {header.class}">
					{header.label}
				</TableHeadCell>
			{/each}
		</TableHead>

		<TableBody>
			{#if loading}
				<Spinner type="dots" />
			{:else if !data.length}
				<p class="text-center text-gray-500">No item available</p>
			{:else}
				{#each data as row, index (row[trackBy] ?? index)}
					<TableBodyRow>
						{#if canSelect}
							<TableBodyCell>
								<Checkbox
									checked={selected.has(getKey(row, index))}
									onchange={(e) => onSingleSelect(e, row, index)}
								/>
							</TableBodyCell>
						{/if}
						{#if children}
							{@render children({ row, index })}
						{:else}
							{#each headers as col, i (String(col.key) + i)}
								<TableBodyCell>{row[col.key]}</TableBodyCell>
							{/each}
						{/if}
					</TableBodyRow>
				{/each}
			{/if}
		</TableBody>
	</Table>
</div>
