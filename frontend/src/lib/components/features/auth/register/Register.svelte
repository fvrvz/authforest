<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { Users } from '$lib/resources/users';
	import { userSchema } from '$lib/schemas/user.schema';
	import { toastService } from '$lib/services/toast.service.svelte';
	import dayjs from 'dayjs';
	import { A, Button, Datepicker, Input, Label } from 'flowbite-svelte';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { zod4, zod4Client } from 'sveltekit-superforms/adapters';

	const superform = superForm(defaults(zod4(userSchema)), {
		validators: zod4Client(userSchema),
		SPA: true,
		onUpdate: async ({ form }) => {
			if (!form.valid || !$tainted) return;

			const [err] = await Users.register(form.data);

			if (err) {
				toastService.error('Registration failed. Please try again.');
				return;
			}

			toastService.success('Account created! Please login.');
			goto(resolve('/login'));
		},
	});

	const { form, tainted, enhance, formId, submitting } = superform;
</script>

<form
	use:enhance
	class="m-auto max-w-lg rounded-lg border border-gray-300 bg-gray-50 p-6 space-y-4"
	id={$formId}
>
	<h1 class="mb-5 text-4xl dark:text-white">Register</h1>
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

		<div class="space-y-2">
			<Label for="password">Password</Label>
			<Input
				type="password"
				bind:value={$form.password}
				id="password"
				autocomplete="new-password"
				placeholder="Enter password"
				required
			/>
		</div>
	</article>

	<div class="flex items-center justify-between gap-2 max-sm:flex-col">
		<span class="text-gray-600">
			Already have an account?
			<A class="text-sm sm:text-base" href={resolve('/login')}>Signin</A>
		</span>
		<Button
			type="submit"
			form={$formId}
			class="w-full cursor-pointer sm:w-auto"
			loading={$submitting}
		>
			Create
		</Button>
	</div>
</form>
