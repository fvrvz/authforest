<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { Roles } from '$lib/resources/roles';
	import { Users } from '$lib/resources/users';
	import { toastService } from '$lib/services/toast.service.svelte';
	import type { Role } from '$lib/types/role.type';
	import dayjs from 'dayjs';
	import {
		Badge,
		Button,
		Checkbox,
		Datepicker,
		Helper,
		Input,
		Label,
	} from 'flowbite-svelte';
	import { onMount } from 'svelte';
	import { SvelteSet } from 'svelte/reactivity';

	let submitting = $state(false);
	let allRoles = $state<Role[]>([]);
	const selectedRoleIds = new SvelteSet<string>();

	let username = $state('');
	let email = $state('');
	let password = $state('');
	let firstName = $state('');
	let lastName = $state('');
	let dob = $state('');

	onMount(async () => {
		const [err, data] = await Roles.getList();
		if (!err) allRoles = data;
	});

	function toggleRole(roleId: string) {
		if (selectedRoleIds.has(roleId)) selectedRoleIds.delete(roleId);
		else selectedRoleIds.add(roleId);
	}

	async function submit() {
		submitting = true;
		const [err] = await Users.adminCreate({
			username,
			email,
			password,
			firstName,
			lastName,
			DOB: dob,
			roleIds: Array.from(selectedRoleIds),
		});
		submitting = false;

		if (err) {
			toastService.error('Failed to create user');
			return;
		}

		toastService.success('User created successfully');
		goto(resolve('/users'));
	}
</script>

<div class="mx-auto max-w-2xl space-y-6">
	<div>
		<h1 class="text-3xl font-bold dark:text-white">Add User</h1>
		<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
			Create a new user account with role assignments.
		</p>
	</div>

	<form
		onsubmit={(e) => {
			e.preventDefault();
			submit();
		}}
		class="space-y-5 rounded-xl border border-gray-200 bg-white p-6 shadow-sm dark:border-gray-600 dark:bg-gray-800"
	>
		<div class="grid gap-4 sm:grid-cols-2">
			<div class="space-y-2">
				<Label for="username">Username</Label>
				<Input
					type="text"
					bind:value={username}
					id="username"
					placeholder="john.doe"
					required
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
				<Label for="password">Password</Label>
				<Input
					type="password"
					bind:value={password}
					id="password"
					placeholder="••••••••"
					required
				/>
			</div>

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
				<Label for="dob">Date of Birth</Label>
				<Datepicker
					id="dob"
					value={dob ? dayjs(dob).toDate() : undefined}
					onselect={(date) => (dob = dayjs(date as Date).format('YYYY-MM-DD'))}
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
				<Helper>Select the roles to assign to this user.</Helper>
			</div>
		{/if}

		<div class="flex items-center gap-3 pt-2">
			<Button type="submit" class="cursor-pointer" loading={submitting}>
				Create User
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
