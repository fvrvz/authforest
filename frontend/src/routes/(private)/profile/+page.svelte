<script lang="ts">
	import { Users } from '$lib/resources/users';
	import { userSchema } from '$lib/schemas/user.schema';
	import { toastService } from '$lib/services/toast.service.svelte';
	import { authStore } from '$lib/state/auth.svelte';
	import dayjs from 'dayjs';
	import { Button, Datepicker, Input, Label, Spinner } from 'flowbite-svelte';
	import { onMount } from 'svelte';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { zod4, zod4Client } from 'sveltekit-superforms/adapters';

	let loading = $state(true);

	const superform = superForm(defaults(zod4(userSchema)), {
		validators: zod4Client(userSchema),
		SPA: true,
		onUpdate: async ({ form: f }) => {
			if (!f.valid || !$tainted || !authStore.user?.username) return;

			const [err] = await Users.update(authStore.user.username, f.data);

			if (err) {
				toastService.error('Failed to update profile');
				return;
			}

			toastService.success('Profile updated successfully');
		},
	});

	const { form, enhance, formId, tainted, submitting } = superform;

	onMount(async () => {
		if (!authStore.user?.username) return;
		const [err, user] = await Users.get(authStore.user.username);
		if (err) {
			toastService.error('Failed to load profile');
			loading = false;
			return;
		}
		$form.firstName = user.firstName ?? '';
		$form.lastName = user.lastName ?? '';
		$form.email = user.email ?? '';
		$form.username = user.username ?? '';
		$form.DOB = user.DOB ?? '';
		loading = false;
	});
</script>

{#if loading}
	<div class="flex justify-center py-12">
		<Spinner size="10" />
	</div>
{:else}
	<form
		use:enhance
		class="rounded-lg border border-gray-300 bg-gray-50 p-6 space-y-4"
		id={$formId}
	>
		<h1 class="mb-5 text-4xl dark:text-white">Profile</h1>
		<article class="grid sm:grid-cols-2 items-center gap-4">
			<div class="space-y-2">
				<Label for="fname">First Name</Label>
				<Input
					type="text"
					bind:value={$form.firstName}
					id="fname"
					placeholder="ex: John"
					required
				/>
			</div>

			<div class="space-y-2">
				<Label for="lname">Last Name</Label>
				<Input
					type="text"
					bind:value={$form.lastName}
					id="lname"
					placeholder="ex: Doe"
				/>
			</div>

			<div class="space-y-2">
				<Label for="email">Email</Label>
				<Input
					type="email"
					bind:value={$form.email}
					id="email"
					placeholder="ex: john.doe@example.com"
					required
				/>
			</div>

			<div class="space-y-2">
				<Label for="dob">Date of Birth</Label>
				<Datepicker
					id="dob"
					required
					value={$form.DOB ? dayjs($form.DOB).toDate() : undefined}
					onselect={(date) => ($form.DOB = date.toString())}
				/>
			</div>

			<div class="space-y-2">
				<Label for="username">Username</Label>
				<Input
					type="text"
					bind:value={$form.username}
					id="username"
					placeholder="ex: john.doe"
					required
				/>
			</div>
		</article>

		<Button
			type="submit"
			form={$formId}
			class="w-full cursor-pointer sm:w-auto"
			loading={$submitting}
		>
			Save
		</Button>
	</form>
{/if}
