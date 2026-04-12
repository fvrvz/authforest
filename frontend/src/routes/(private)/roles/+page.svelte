<script lang="ts">
	import Table from '$lib/components/common/table/Table.svelte';
	import { Roles } from '$lib/resources/roles';
	import { dialogService } from '$lib/services/dialog.service.svelte';
	import { toastService } from '$lib/services/toast.service.svelte';
	import type { Role } from '$lib/types/role.type';
	import type { TableHeader } from '$lib/types/table.type';
	import {
	  Button,
	  Input,
	  Label,
	  Modal,
	  TableBodyCell,
	  Textarea,
	} from 'flowbite-svelte';
	import { Pencil, Plus, Trash } from 'lucide-svelte';
	import { onMount } from 'svelte';

	let loading = $state(true);
	let roles = $state<Role[]>([]);
	let showModal = $state(false);
	let submitting = $state(false);
	let editingRole = $state<Role | null>(null);
	let roleName = $state('');
	let roleDescription = $state('');

	onMount(async () => {
		await loadRoles();
	});

	async function loadRoles() {
		loading = true;
		const [err, data] = await Roles.getList();
		if (err) toastService.error('Failed to load roles');
		else roles = data;
		loading = false;
	}

	function openCreate() {
		editingRole = null;
		roleName = '';
		roleDescription = '';
		showModal = true;
	}

	function openEdit(role: Role) {
		editingRole = role;
		roleName = role.name;
		roleDescription = role.description;
		showModal = true;
	}

	async function save() {
		submitting = true;

		if (editingRole) {
			const [err, updated] = await Roles.update(editingRole.id, {
				name: roleName,
				description: roleDescription,
			});
			if (err) {
				toastService.error('Failed to update role');
			} else {
				roles = roles.map((r) => (r.id === editingRole!.id ? updated : r));
				toastService.success('Role updated');
				showModal = false;
			}
		} else {
			const [err, created] = await Roles.create({
				name: roleName,
				description: roleDescription,
			});
			if (err) {
				toastService.error('Failed to create role');
			} else {
				roles = [...roles, created];
				toastService.success('Role created');
				showModal = false;
			}
		}

		submitting = false;
	}

	async function deleteRole(role: Role) {
		const confirmed = await dialogService.confirm({
			title: 'Delete Role',
			message: `Are you sure you want to delete the "${role.name}" role? This will remove it from all users.`,
		});
		if (!confirmed) return;

		const [err] = await Roles.delete(role.id);
		if (err) {
			toastService.error('Failed to delete role');
		} else {
			roles = roles.filter((r) => r.id !== role.id);
			toastService.success('Role deleted');
		}
	}

	const headers: TableHeader<Role>[] = [
		{ key: 'name', label: 'Name' },
		{ key: 'description', label: 'Description' },
		{ key: 'id', class: 'w-fit' },
	];
</script>

<div class="space-y-4">
	<div
		class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
	>
		<div>
			<h1 class="text-3xl font-bold dark:text-white">Roles</h1>
			<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
				Manage roles that can be assigned to users.
			</p>
		</div>
		<Button class="flex cursor-pointer items-center gap-2 self-end sm:self-auto" onclick={openCreate}>
			<Plus class="size-4" />
			<span class="hidden sm:inline">New Role</span>
		</Button>
	</div>

	<Table data={roles} {loading} {headers} trackBy="id">
		{#snippet children({ row })}
			<TableBodyCell class="font-medium dark:text-white">
				{row.name}
			</TableBodyCell>
			<TableBodyCell>{row.description}</TableBodyCell>
			<TableBodyCell class="w-fit">
				<div class="flex gap-1">
					<button
						type="button"
						class="cursor-pointer rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-primary-600 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-primary-400"
						onclick={() => openEdit(row)}
					>
						<Pencil class="size-4" />
					</button>
					<button
						type="button"
						class="cursor-pointer rounded-lg p-2 text-gray-500 transition-colors hover:bg-red-50 hover:text-red-600 dark:text-gray-400 dark:hover:bg-red-950/30 dark:hover:text-red-400"
						onclick={() => deleteRole(row)}
					>
						<Trash class="size-4" />
					</button>
				</div>
			</TableBodyCell>
		{/snippet}
	</Table>
</div>

<Modal
	bind:open={showModal}
	title={editingRole ? 'Edit Role' : 'Create Role'}
	size="sm"
>
	<form
		onsubmit={(e) => {
			e.preventDefault();
			save();
		}}
		class="space-y-4"
	>
		<div class="space-y-2">
			<Label for="role_name">Name</Label>
			<Input
				type="text"
				bind:value={roleName}
				id="role_name"
				placeholder="e.g. editor"
				required
			/>
		</div>

		<div class="space-y-2">
			<Label for="role_description">Description</Label>
			<Textarea
				bind:value={roleDescription}
				id="role_description"
				placeholder="What this role allows..."
				rows={3}
			/>
		</div>

		<div class="flex justify-end gap-3">
			<Button
				color="alternative"
				class="cursor-pointer"
				onclick={() => (showModal = false)}
			>
				Cancel
			</Button>
			<Button type="submit" class="cursor-pointer" loading={submitting}>
				{editingRole ? 'Update' : 'Create'}
			</Button>
		</div>
	</form>
</Modal>
