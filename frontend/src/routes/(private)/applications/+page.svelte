<script lang="ts">
	import { resolve } from '$app/paths';
	import Table from '$lib/components/common/table/Table.svelte';
	import { OIDC } from '$lib/resources/oidc';
	import { dialogService } from '$lib/services/dialog.service.svelte';
	import { toastService } from '$lib/services/toast.service.svelte';
	import type { OAuthClient } from '$lib/types/oidc.type';
	import type { TableHeader } from '$lib/types/table.type';
	import dayjs from 'dayjs';
	import { Badge, Button, TableBodyCell } from 'flowbite-svelte';
	import { Check, Copy, Pencil, Plus, Trash } from 'lucide-svelte';
	import { onMount } from 'svelte';

	let loading = $state(true);
	let clients = $state<OAuthClient[]>([]);
	let copiedId = $state<string | null>(null);

	onMount(async () => {
		await loadClients();
	});

	async function loadClients() {
		loading = true;
		const [err, data] = await OIDC.listClients();
		if (err) toastService.error('Failed to load applications');
		else clients = data;
		loading = false;
	}

	async function deleteClient(clientId: string) {
		const confirmed = await dialogService.confirm({
			title: 'Delete Application',
			message:
				'Are you sure you want to delete this application? This action cannot be undone.',
		});
		if (!confirmed) return;

		const [err] = await OIDC.deleteClient(clientId);
		if (err) {
			toastService.error('Failed to delete application');
		} else {
			clients = clients.filter((c) => c.client_id !== clientId);
			toastService.success('Application deleted');
		}
	}

	async function copyClientId(clientId: string) {
		await navigator.clipboard.writeText(clientId);
		copiedId = clientId;
		toastService.success('Client ID copied to clipboard');
		setTimeout(() => (copiedId = null), 2000);
	}

	const headers: TableHeader<OAuthClient>[] = [
		{ key: 'client_name', label: 'Application Name' },
		{ key: 'client_id', label: 'Client ID' },
		{ key: 'client_type', label: 'Type' },
		{ key: 'scopes', label: 'Scopes' },
		{ key: 'created_at', label: 'Created' },
		{ key: 'client_id', class: 'w-fit' },
	];
</script>

<div class="space-y-4">
	<div
		class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
	>
		<div>
			<h1 class="text-3xl font-bold dark:text-white">Applications</h1>
			<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
				Manage your registered OAuth 2.0 / OIDC client applications.
			</p>
		</div>
		<a
			href={resolve('/applications/new')}
			class="flex items-center gap-2 rounded-lg bg-primary-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
		>
			<Plus class="size-4" />
			New Application
		</a>
	</div>

	<Table data={clients} {loading} {headers} trackBy="client_id">
		{#snippet children({ row })}
			<TableBodyCell class="font-medium dark:text-white">
				{row.client_name}
			</TableBodyCell>
			<TableBodyCell>
				<div class="flex items-center gap-1.5">
					<code class="text-xs text-gray-600 dark:text-gray-300">
						{row.client_id.slice(0, 8)}...
					</code>
					<button
						type="button"
						class="cursor-pointer text-gray-400 hover:text-gray-600 dark:hover:text-gray-200"
						onclick={() => copyClientId(row.client_id)}
					>
						{#if copiedId === row.client_id}
							<Check class="size-3.5 text-green-500" />
						{:else}
							<Copy class="size-3.5" />
						{/if}
					</button>
				</div>
			</TableBodyCell>
			<TableBodyCell>
				<Badge color={row.client_type === 'public' ? 'green' : 'blue'}>
					{row.client_type}
				</Badge>
			</TableBodyCell>
			<TableBodyCell>
				<div class="flex flex-wrap gap-1">
					{#each row.scopes.split(' ') as scope (scope)}
						<Badge color="gray" class="text-xs">{scope}</Badge>
					{/each}
				</div>
			</TableBodyCell>
			<TableBodyCell class="text-sm text-gray-500">
				{row.created_at ? dayjs(row.created_at).format('MMM D, YYYY') : '—'}
			</TableBodyCell>
			<TableBodyCell class="w-fit">
				<div class="flex gap-1">
					<a href={resolve(`/applications/${row.client_id}`)}>
						<Button
							pill
							size="sm"
							outline
							color="primary"
							class="cursor-pointer"
						>
							<Pencil class="size-4" />
						</Button>
					</a>
					<Button
						pill
						size="sm"
						outline
						color="red"
						class="cursor-pointer"
						onclick={() => deleteClient(row.client_id)}
					>
						<Trash class="size-4" />
					</Button>
				</div>
			</TableBodyCell>
		{/snippet}
	</Table>
</div>
