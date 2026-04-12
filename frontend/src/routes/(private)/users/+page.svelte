<script lang="ts">
	import { resolve } from '$app/paths';
	import Table from '$lib/components/common/table/Table.svelte';
	import { Users } from '$lib/resources/users';
	import { dialogService } from '$lib/services/dialog.service.svelte';
	import { toastService } from '$lib/services/toast.service.svelte';
	import type { TableHeader } from '$lib/types/table.type';
	import type { User } from '$lib/types/user.type';
	import dayjs from 'dayjs';
	import { Badge, Button, TableBodyCell } from 'flowbite-svelte';
	import { Pencil, Plus, Trash } from 'lucide-svelte';
	import { onMount } from 'svelte';

	let loading = $state(true);
	let users = $state<User[]>([]);

	onMount(async () => {
		const [err, data] = await Users.getList();
		if (err) toastService.error('Failed to load users');
		else users = data;
		loading = false;
	});

	async function deleteUser(username: string) {
		const confirmed = await dialogService.confirm({
			title: 'Delete User',
			message: `Are you sure you want to delete user "${username}"? This action cannot be undone.`,
		});
		if (!confirmed) return;

		const [err] = await Users.delete(username);
		if (err) {
			toastService.error('Failed to delete user');
		} else {
			users = users.filter((u) => u.username !== username);
			toastService.success(`User "${username}" deleted`);
		}
	}

	const headers: TableHeader<User>[] = [
		{ key: 'fullName', label: 'Name' },
		{ key: 'username', label: 'Username' },
		{ key: 'email', label: 'Email' },
		{ key: 'roles', label: 'Roles' },
		{ key: 'DOB', label: 'Date of Birth' },
		{ key: 'id', class: 'w-fit' },
	];
</script>

<div class="space-y-4">
	<div
		class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
	>
		<div>
			<h1 class="text-3xl font-bold dark:text-white">Users</h1>
			<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
				Manage user accounts registered with the Identity Provider.
			</p>
		</div>
		<a
			href={resolve('/users/new')}
			class="flex items-center gap-2 rounded-lg bg-primary-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
		>
			<Plus class="size-4" />
			New User
		</a>
	</div>

	<Table canSelect data={users} {loading} {headers} trackBy="id">
		{#snippet children({ row })}
			<TableBodyCell class="font-medium dark:text-white">
				{row.fullName}
			</TableBodyCell>
			<TableBodyCell>{row.username}</TableBodyCell>
			<TableBodyCell>{row.email}</TableBodyCell>
			<TableBodyCell>
				<div class="flex flex-wrap gap-1">
					{#each row.roles ?? [] as role (role.id)}
						<Badge color="blue" class="text-xs">{role.name}</Badge>
					{/each}
				</div>
			</TableBodyCell>
			<TableBodyCell>
				{row.DOB ? dayjs(row.DOB).format('MMM D, YYYY') : '—'}
			</TableBodyCell>
			<TableBodyCell class="w-fit">
				<div class="flex gap-1">
					<a href={resolve(`/users/${row.username}`)}>
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
						onclick={() => deleteUser(row.username)}
					>
						<Trash class="size-4" />
					</Button>
				</div>
			</TableBodyCell>
		{/snippet}
	</Table>
</div>
