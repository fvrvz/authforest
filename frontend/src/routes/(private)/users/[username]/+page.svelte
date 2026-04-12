<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Roles } from '$lib/resources/roles';
	import { Users } from '$lib/resources/users';
	import { dialogService } from '$lib/services/dialog.service.svelte';
	import { toastService } from '$lib/services/toast.service.svelte';
	import type { Role } from '$lib/types/role.type';
	import type { User } from '$lib/types/user.type';
	import dayjs from 'dayjs';
	import {
	  Badge,
	  Button,
	  Checkbox,
	  Datepicker,
	  Input,
	  Label,
	  Spinner,
	} from 'flowbite-svelte';
	import { KeyRound } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { SvelteSet } from 'svelte/reactivity';

	const username = page.params.username!;
	let loading = $state(true);
	let saving = $state(false);
	let resettingPassword = $state(false);
	let user = $state<User | null>(null);
	let allRoles = $state<Role[]>([]);
	const selectedRoleIds = new SvelteSet<string>();

	let firstName = $state('');
	let lastName = $state('');
	let email = $state('');
	let dob = $state('');

	onMount(async () => {
		const [userErr, userData] = await Users.get(username);
		const [rolesErr, rolesData] = await Roles.getList();

		if (userErr) {
			toastService.error('Failed to load user');
			loading = false;
			return;
		}

		user = userData;
		firstName = user.firstName ?? '';
		lastName = user.lastName ?? '';
		email = user.email ?? '';
		dob = user.DOB ?? '';
		selectedRoleIds.clear();
		(user.roles ?? []).forEach((r) => selectedRoleIds.add(r.id));

		if (!rolesErr) allRoles = rolesData;
		loading = false;
	});

	function toggleRole(roleId: string) {
		if (selectedRoleIds.has(roleId)) selectedRoleIds.delete(roleId);
		else selectedRoleIds.add(roleId);
	}

	async function save() {
		saving = true;
		const [err] = await Users.adminUpdate(username, {
			email,
			firstName,
			lastName,
			DOB: dob,
			roleIds: Array.from(selectedRoleIds),
		});
		saving = false;

		if (err) {
			toastService.error('Failed to update user');
			return;
		}

		toastService.success('User updated successfully');
		goto(resolve('/users'));
	}

	async function resetPassword() {
		if (!user?.email) {
			toastService.error('User has no email address');
			return;
		}

		const confirmed = await dialogService.confirm({
			title: 'Send Password Reset',
			message: `Send a password reset email to ${user.email}?`,
			confirmColor: 'primary',
			confirmText: 'Send',
		});
		if (!confirmed) return;

		resettingPassword = true;
		const [err] = await Users.requestPasswordReset(user.email);
		resettingPassword = false;

		if (err) {
			toastService.error('Failed to send reset email');
			return;
		}

		toastService.success('Password reset email sent');
	}
</script>

{#if loading}
	<div class="flex justify-center py-12">
		<Spinner size="10" />
	</div>
{:else if user}
	<div class="mx-auto max-w-2xl space-y-6">
		<div>
			<h1 class="text-3xl font-bold dark:text-white">Edit User</h1>
			<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
				Update details and roles for
				<span class="font-medium">{user.username}</span>.
			</p>
		</div>

		<form
			onsubmit={(e) => {
				e.preventDefault();
				save();
			}}
			class="space-y-5 rounded-xl border border-gray-200 bg-white p-6 shadow-sm dark:border-gray-600 dark:bg-gray-800"
		>
			<div class="grid gap-4 sm:grid-cols-2">
				<div class="space-y-2">
					<Label for="firstName">First Name</Label>
					<Input
						type="text"
						bind:value={firstName}
						id="firstName"
						placeholder="John"
						required
					/>
				</div>

				<div class="space-y-2">
					<Label for="lastName">Last Name</Label>
					<Input
						type="text"
						bind:value={lastName}
						id="lastName"
						placeholder="Doe"
					/>
				</div>

				<div class="space-y-2">
					<Label for="email">Email</Label>
					<Input
						type="email"
						bind:value={email}
						id="email"
						placeholder="john@example.com"
						required
					/>
				</div>

				<div class="space-y-2">
					<Label for="dob">Date of Birth</Label>
					<Datepicker
						id="dob"
						value={dob ? dayjs(dob).toDate() : undefined}
						onselect={(date) =>
							(dob = dayjs(date as Date).format('YYYY-MM-DD'))}
					/>
				</div>
			</div>

			{#if allRoles.length > 0}
				<div class="space-y-2">
					<Label>Roles</Label>
					<div class="flex flex-wrap gap-3">
						{#each allRoles as role (role.id)}
							<Checkbox
								checked={selectedRoleIds.has(role.id)}
								onchange={() => toggleRole(role.id)}
							>
								<Badge color={selectedRoleIds.has(role.id) ? 'blue' : 'gray'}>
									{role.name}
								</Badge>
							</Checkbox>
						{/each}
					</div>
				</div>
			{/if}

			<div
				class="flex flex-col items-center gap-3 border-t border-gray-200 pt-4 sm:flex-row sm:flex-wrap dark:border-gray-600"
			>
				<Button type="submit" class="w-full cursor-pointer sm:w-auto" loading={saving}>
					Save Changes
				</Button>
				<Button
					color="alternative"
					class="flex w-full cursor-pointer items-center justify-center gap-2 sm:w-auto"
					loading={resettingPassword}
					onclick={resetPassword}
				>
					<KeyRound class="size-4" />
					Reset Password
				</Button>
				<a
					href={resolve('/users')}
					class="text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
				>
					Cancel
				</a>
			</div>
		</form>
	</div>
{/if}
