<script lang="ts">
	import { resolve } from '$app/paths';
	import { Users } from '$lib/resources/users';
	import { toastService } from '$lib/services/toast.service.svelte';
	import { authStore } from '$lib/state/auth.svelte';
	import type { User } from '$lib/types/user.type';
	import dayjs from 'dayjs';
	import { Button, Datepicker, Input, Label, Spinner } from 'flowbite-svelte';
	import { onMount } from 'svelte';

	const username = authStore.user?.username ?? '';
	let loading = $state(true);
	let saving = $state(false);
	let user = $state<User | null>(null);

	let firstName = $state('');
	let lastName = $state('');
	let email = $state('');
	let dob = $state('');

	onMount(async () => {
		if (!username) {
			toastService.error('Not logged in');
			loading = false;
			return;
		}

		const [err, userData] = await Users.get(username);
		if (err) {
			toastService.error('Failed to load profile');
			loading = false;
			return;
		}

		user = userData;
		firstName = user.firstName ?? '';
		lastName = user.lastName ?? '';
		email = user.email ?? '';
		dob = user.DOB ?? '';
		loading = false;
	});

	async function save() {
		saving = true;
		const [err] = await Users.update(username, {
			email,
			firstName,
			lastName,
			DOB: dob,
		});
		saving = false;

		if (err) {
			toastService.error('Failed to update profile');
			return;
		}

		toastService.success('Profile updated successfully');
	}
</script>

{#if loading}
	<div class="flex justify-center py-12">
		<Spinner size="10" />
	</div>
{:else if user}
	<div class="mx-auto max-w-2xl space-y-6">
		<div>
			<h1 class="text-3xl font-bold dark:text-white">Profile</h1>
			<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
				Update your personal details.
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

			<div class="flex flex-col items-center gap-3 pt-2 sm:flex-row">
				<Button type="submit" class="w-full cursor-pointer sm:w-auto" loading={saving}>
					Save Changes
				</Button>
				<a
					href={resolve('/dashboard')}
					class="text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
				>
					Cancel
				</a>
			</div>
		</form>
	</div>
{/if}
